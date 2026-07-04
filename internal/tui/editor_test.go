package tui

import (
	"os"
	"strings"
	"testing"
)

func TestResolveEditor(t *testing.T) {
	// Save + restore env.
	oldV, hadV := os.LookupEnv("VISUAL")
	oldE, hadE := os.LookupEnv("EDITOR")
	t.Cleanup(func() {
		restore("VISUAL", oldV, hadV)
		restore("EDITOR", oldE, hadE)
	})

	os.Unsetenv("VISUAL")
	os.Unsetenv("EDITOR")
	if name, args := resolveEditor(); name != "vi" || len(args) != 0 {
		t.Fatalf("fallback: got %q %v, want vi []", name, args)
	}

	os.Setenv("EDITOR", "code --wait")
	if name, args := resolveEditor(); name != "code" || len(args) != 1 || args[0] != "--wait" {
		t.Fatalf("EDITOR split: got %q %v", name, args)
	}

	os.Setenv("VISUAL", "nvim")
	if name, _ := resolveEditor(); name != "nvim" {
		t.Fatalf("VISUAL precedence: got %q, want nvim", name)
	}
}

func restore(key, val string, had bool) {
	if had {
		os.Setenv(key, val)
	} else {
		os.Unsetenv(key)
	}
}

func TestEditorPreviewSingleLine(t *testing.T) {
	if got := editorPreview("just one line"); got != "just one line" {
		t.Fatalf("single line should pass through, got %q", got)
	}
}

func TestEditorPreviewMultiLine(t *testing.T) {
	got := editorPreview("\n\nfirst real line\nsecond\nthird")
	if !strings.HasPrefix(got, "first real line") {
		t.Fatalf("preview should start at first non-empty line, got %q", got)
	}
	if !strings.Contains(got, "[5 lines]") {
		t.Fatalf("preview should report line count, got %q", got)
	}
}

func TestOpenEditorCmdNonNil(t *testing.T) {
	// The returned command is opaque (tea.ExecProcess), so we only assert it is
	// produced without error for a normal seed; the readback/cancel branches are
	// covered via TestEditorFinishedMsgHandling.
	if cmd := openEditorCmd("seed"); cmd == nil {
		t.Fatal("openEditorCmd returned nil for a valid seed")
	}
}

func TestEditorFinishedMsgHandling(t *testing.T) {
	base := liveModel{inputText: "prior"}

	// Success replaces input, clears suggest.
	m, _ := base.Update(editorFinishedMsg{content: "new text"})
	lm := m.(liveModel)
	if lm.inputText != "new text" || lm.suggest {
		t.Fatalf("success: inputText=%q suggest=%v", lm.inputText, lm.suggest)
	}

	// Cancel keeps prior content.
	m2, _ := base.Update(editorFinishedMsg{canceled: true})
	lm2 := m2.(liveModel)
	if lm2.inputText != "prior" {
		t.Fatalf("cancel should preserve input, got %q", lm2.inputText)
	}
	if lm2.statusMsg == "" {
		t.Fatalf("cancel should set a status message")
	}

	// Error keeps prior content and surfaces the error.
	m3, _ := base.Update(editorFinishedMsg{err: os.ErrNotExist})
	lm3 := m3.(liveModel)
	if lm3.inputText != "prior" || lm3.inputErr == "" {
		t.Fatalf("error: inputText=%q inputErr=%q", lm3.inputText, lm3.inputErr)
	}
}

func TestEditorCommandInSpecs(t *testing.T) {
	found := false
	for _, s := range commandSpecs {
		if s.Name == "/editor" {
			found = true
		}
	}
	if !found {
		t.Fatal("/editor missing from commandSpecs")
	}
	if got := commandMatches("/edi"); len(got) == 0 || got[0].Name != "/editor" {
		t.Fatalf("/edi should suggest /editor, got %v", got)
	}
}
