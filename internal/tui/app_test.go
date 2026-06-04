package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"parley-deck-cli/internal/protocol"
)

func TestInitModelViewPromptsForWorkspaceInit(t *testing.T) {
	root := t.TempDir()
	model := newInitModel(root, nil)
	model.width = 100

	view := model.View()
	for _, want := range []string{
		"Parley Deck setup",
		"Workspace is not initialized",
		"Path:",
		filepath.Base(root),
		"Press i or enter to initialize",
		"Keys: i/enter initialize  q/esc/ctrl+c quit",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("view missing %q\n%s", want, view)
		}
	}
}

// After a successful init the wizard quits so the caller can re-read the
// workspace and open the unified Home view (the old workspace dashboard that
// used to be shown inline was retired).
func TestInitModelInitializesThenQuitsForHandoff(t *testing.T) {
	root := t.TempDir()
	init := newInitModel(root, nil)

	updated, cmd := init.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	init = updated.(initModel)
	if cmd == nil {
		t.Fatal("initialize key did not return a command")
	}
	if !init.initializing {
		t.Fatal("model did not enter initializing state")
	}

	updated, quitCmd := init.Update(cmd())
	init, ok := updated.(initModel)
	if !ok {
		t.Fatalf("post-init model type=%T, want tui.initModel", updated)
	}
	if init.status == nil {
		t.Fatal("workspace status was not recorded after init")
	}
	if _, err := os.Stat(filepath.Join(root, protocol.DeckDir, "COOPERATION.md")); err != nil {
		t.Fatalf("workspace was not initialized: %v", err)
	}
	if quitCmd == nil {
		t.Fatal("init did not return a quit command for handoff")
	}
	if _, ok := quitCmd().(tea.QuitMsg); !ok {
		t.Fatalf("post-init command did not quit the wizard")
	}
}

func TestInitModelFailureStaysOnSetupScreen(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, protocol.DeckDir), []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	model := newInitModel(root, nil)

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updated.(initModel)
	if cmd == nil {
		t.Fatal("initialize key did not return a command")
	}

	updated, _ = model.Update(cmd())
	model = updated.(initModel)
	if model.status != nil {
		t.Fatal("status was loaded despite init failure")
	}
	if model.initializing {
		t.Fatal("model stayed in initializing state after failure")
	}
	if model.errText == "" {
		t.Fatal("init failure was not recorded")
	}
	if _, err := protocol.ReadWorkspaceStatus(root); err == nil {
		t.Fatalf("workspace unexpectedly became readable")
	}

	view := model.View()
	for _, want := range []string{"Workspace is not initialized", "Init failed:"} {
		if !strings.Contains(view, want) {
			t.Fatalf("failure view missing %q\n%s", want, view)
		}
	}
}
