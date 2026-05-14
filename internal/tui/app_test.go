package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"parley-deck-cli/internal/agents"
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
	init := newInitModel(root, nil)

	updated, cmd := init.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	init = updated.(initModel)
	if cmd == nil {
		t.Fatal("initialize key did not return a command")
	}
	if !init.initializing {
		t.Fatal("model did not enter initializing state")
	}

	updated, _ = init.Update(cmd())
	dashboard, ok := updated.(model)
	if !ok {
		t.Fatalf("post-init model type=%T, want tui.model", updated)
	}
	if _, err := os.Stat(filepath.Join(root, protocol.DeckDir, "COOPERATION.md")); err != nil {
		t.Fatalf("workspace was not initialized: %v", err)
	}

	view := dashboard.View()
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

func TestDashboardRendersSelectedAgentDetails(t *testing.T) {
	m := newModel(testStatus(), testDiscoveries())
	m.width = 120

	view := m.View()
	for _, want := range []string{
		"> codex",
		"Agent details",
		"id: codex",
		"configured launch: headless",
		"effective: headless",
		"model: cli-default",
		"sandbox: workspace-write",
		"headless: codex exec -",
		"interactive: codex",
		"session-only preview",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("dashboard view missing %q\n%s", want, view)
		}
	}
}

func TestDashboardRendersFallbackCommandDetails(t *testing.T) {
	m := newModel(testStatus(), testDiscoveries())
	m.selectedAgent = 1
	m.width = 120

	view := m.View()
	for _, want := range []string{
		"id: claude",
		"configured launch: interactive",
		"backend: unknown",
		"headless: claude",
		"interactive: claude --resume {prompt_path}",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("dashboard view missing %q\n%s", want, view)
		}
	}
	if strings.Contains(view, "interactive: claude \n") {
		t.Fatalf("interactive command has trailing space\n%s", view)
	}
}

func TestDashboardAgentNavigationClamps(t *testing.T) {
	m := newModel(testStatus(), testDiscoveries())

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(model)
	if m.selectedAgent != 0 {
		t.Fatalf("selectedAgent=%d, want 0", m.selectedAgent)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.selectedAgent != 1 {
		t.Fatalf("selectedAgent=%d, want 1", m.selectedAgent)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(model)
	if m.selectedAgent != 0 {
		t.Fatalf("selectedAgent=%d, want 0 after up", m.selectedAgent)
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(model)
	if m.selectedAgent != 1 {
		t.Fatalf("selectedAgent=%d, want 1 after down", m.selectedAgent)
	}
}

func TestDashboardFocusSwitchPreservesSelection(t *testing.T) {
	m := newModel(testStatus(), testDiscoveries())

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.selectedAgent != 1 {
		t.Fatalf("selectedAgent=%d, want 1", m.selectedAgent)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.focus != focusIdeas {
		t.Fatalf("focus=%s, want %s", m.focus, focusIdeas)
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.selectedIdea != 1 {
		t.Fatalf("selectedIdea=%d, want 1", m.selectedIdea)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = updated.(model)
	if m.focus != focusAgents {
		t.Fatalf("focus=%s, want %s", m.focus, focusAgents)
	}
	if m.selectedAgent != 1 {
		t.Fatalf("selectedAgent=%d, want preserved 1", m.selectedAgent)
	}
}

func TestDashboardLaunchModeOverridesAreSessionOnly(t *testing.T) {
	m := newModel(testStatus(), testDiscoveries())

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	m = updated.(model)
	if got := m.launchOverrides["codex"]; got != agents.LaunchManual {
		t.Fatalf("override=%q, want manual", got)
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(model)
	if got := m.launchOverrides["codex"]; got != agents.LaunchHeadless {
		t.Fatalf("override=%q, want headless", got)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = updated.(model)
	if got := m.launchOverrides["codex"]; got != agents.LaunchInteractive {
		t.Fatalf("override=%q, want interactive", got)
	}
	if got := agents.LaunchModeOrDefault(m.agents[0].LaunchMode); got != agents.LaunchHeadless {
		t.Fatalf("underlying launch mode mutated to %q", got)
	}
	if !strings.Contains(m.View(), "effective: interactive (session only)") {
		t.Fatalf("view missing session override\n%s", m.View())
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	m = updated.(model)
	if _, ok := m.launchOverrides["codex"]; ok {
		t.Fatalf("override was not cleared: %+v", m.launchOverrides)
	}
	if !strings.Contains(m.View(), "effective: headless") {
		t.Fatalf("view did not return to configured mode\n%s", m.View())
	}
}

func TestDashboardModeKeysNoopOutsideAgentFocus(t *testing.T) {
	m := newModel(testStatus(), testDiscoveries())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.focus != focusIdeas {
		t.Fatalf("focus=%s, want ideas", m.focus)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = updated.(model)
	if len(m.launchOverrides) != 0 {
		t.Fatalf("override created while ideas pane focused: %+v", m.launchOverrides)
	}
}

func testStatus() protocol.WorkspaceStatus {
	return protocol.WorkspaceStatus{
		Root:      "/repo",
		Transport: "github-pr",
		Ideas: []protocol.IdeaStatus{
			{Slug: "first", Status: "final", Participants: []string{"codex"}, Path: "/repo/parley-deck/ideas/first"},
			{Slug: "second", Status: "round-01", Participants: []string{"claude"}, Path: "/repo/parley-deck/ideas/second"},
		},
	}
}

func testDiscoveries() []agents.Discovery {
	return []agents.Discovery{
		{
			Spec: agents.Spec{
				ID:                    "codex",
				Commands:              []string{"codex"},
				LaunchMode:            agents.LaunchHeadless,
				HeadlessMode:          "codex exec -",
				InteractivePromptMode: agents.InteractivePromptNone,
				InteractiveInvoke:     agents.InteractiveInvokePrintOnly,
				SandboxMode:           "workspace-write",
				ApprovalPolicy:        "on-failure",
				Model:                 agents.CLIDefault,
				TimeoutMS:             1800000,
				ExternalBackend:       agents.ExternalHosted,
			},
			Path:    "/usr/bin/codex",
			Found:   true,
			Version: "codex test",
		},
		{
			Spec: agents.Spec{
				ID:                    "claude",
				Commands:              []string{"claude"},
				LaunchMode:            agents.LaunchInteractive,
				InteractiveCommand:    "claude",
				InteractiveArgs:       []string{"--resume", "{prompt_path}"},
				InteractivePromptMode: agents.InteractivePromptFile,
				InteractiveInvoke:     agents.InteractiveInvokePrintOnly,
			},
			Path:    "/usr/bin/claude",
			Found:   true,
			Version: "claude test",
		},
	}
}
