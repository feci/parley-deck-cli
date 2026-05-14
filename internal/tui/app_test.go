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

func TestInitModelInitializesAndShowsDashboard(t *testing.T) {
	root := t.TempDir()
	model := newInitModel(root, nil)

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updated.(initModel)
	if cmd == nil {
		t.Fatal("initialize key did not return a command")
	}
	if !model.initializing {
		t.Fatal("model did not enter initializing state")
	}

	updated, _ = model.Update(cmd())
	model = updated.(initModel)
	if model.status == nil {
		t.Fatalf("status was not loaded after init; err=%s", model.errText)
	}
	if _, err := os.Stat(filepath.Join(root, protocol.DeckDir, "COOPERATION.md")); err != nil {
		t.Fatalf("workspace was not initialized: %v", err)
	}

	view := model.View()
	for _, want := range []string{"Parley Deck  transport=", "Protocol", "Agents"} {
		if !strings.Contains(view, want) {
			t.Fatalf("dashboard view missing %q\n%s", want, view)
		}
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
