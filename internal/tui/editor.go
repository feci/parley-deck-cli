package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// editor.go implements the /editor + ctrl+e composer (tui-editor-composer):
// suspend the TUI, open $VISUAL/$EDITOR/vi on a temp file, and on a clean exit
// return its contents so the caller can drop them into the composer. Nothing is
// sent directly — the existing Enter path keeps the steer/answer semantics.

// editorFinishedMsg is delivered after the external editor exits.
type editorFinishedMsg struct {
	content  string // trimmed file content (valid only when err == nil && !canceled)
	err      error  // exec/read failure
	canceled bool   // editor exited non-zero → keep prior composer content
}

// resolveEditor picks the editor command and its leading args, honoring $VISUAL
// then $EDITOR (split on fields so values like "code --wait" work), falling back
// to "vi". Returns (name, args-before-path).
func resolveEditor() (string, []string) {
	for _, env := range []string{os.Getenv("VISUAL"), os.Getenv("EDITOR")} {
		if fields := strings.Fields(env); len(fields) > 0 {
			return fields[0], fields[1:]
		}
	}
	return "vi", nil
}

// openEditorCmd writes seed into a 0600 temp file and returns a tea.Cmd that
// suspends the TUI, runs the editor on that file, then reads it back and removes
// it on every exit path (best-effort; a host SIGKILL mid-edit cannot be covered).
func openEditorCmd(seed string) tea.Cmd {
	f, err := os.CreateTemp(os.TempDir(), "parley-editor-*.md")
	if err != nil {
		return func() tea.Msg { return editorFinishedMsg{err: fmt.Errorf("editor temp file: %w", err)} }
	}
	path := f.Name()
	if err := f.Chmod(0o600); err != nil {
		f.Close()
		os.Remove(path)
		return func() tea.Msg { return editorFinishedMsg{err: fmt.Errorf("editor temp perms: %w", err)} }
	}
	if seed != "" {
		if _, err := f.WriteString(seed); err != nil {
			f.Close()
			os.Remove(path)
			return func() tea.Msg { return editorFinishedMsg{err: fmt.Errorf("editor seed write: %w", err)} }
		}
	}
	if err := f.Close(); err != nil {
		os.Remove(path)
		return func() tea.Msg { return editorFinishedMsg{err: fmt.Errorf("editor temp close: %w", err)} }
	}

	name, args := resolveEditor()
	cmd := exec.Command(name, append(args, path)...) //nolint:gosec // $EDITOR is user-controlled local config
	return tea.ExecProcess(cmd, func(runErr error) tea.Msg {
		defer os.Remove(path)
		if runErr != nil {
			// Non-zero editor exit (or launch failure): treat as cancel, keep prior input.
			return editorFinishedMsg{canceled: true, err: runErr}
		}
		b, readErr := os.ReadFile(path)
		if readErr != nil {
			return editorFinishedMsg{err: fmt.Errorf("editor read back: %w", readErr)}
		}
		return editorFinishedMsg{content: strings.TrimRight(string(b), "\n")}
	})
}

// editorPreview renders composer content for the single-line input row: a
// multi-line draft (from the editor) is flattened to its first non-empty line
// plus a "[N lines]" affordance, so it never breaks the one-line row. The raw
// value in m.inputText is still what gets submitted.
func editorPreview(s string) string {
	if !strings.ContainsAny(s, "\n") {
		return s
	}
	lines := strings.Split(s, "\n")
	first := ""
	for _, ln := range lines {
		if strings.TrimSpace(ln) != "" {
			first = ln
			break
		}
	}
	return fmt.Sprintf("%s [%d lines]", first, len(lines))
}
