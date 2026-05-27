package tui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runaction"
	"parley-deck-cli/internal/runstate"
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
	for _, want := range []string{"Parley Deck  transport=", "Sessions", "Agents"} {
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
	m := newTestModel(nil)
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
		"headless: codex exec --skip-git-",
		"repo-check -",
		"acp: codex acp",
		"interactive: codex",
		"Sessions",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("dashboard view missing %q\n%s", want, view)
		}
	}
}

func TestDashboardRendersFallbackCommandDetails(t *testing.T) {
	m := newTestModel(nil)
	m.selectedAgent = 1
	m.width = 180

	view := m.View()
	for _, want := range []string{
		"id: claude",
		"configured launch: interactive",
		"backend:",
		"unknown",
		"headless: claude",
		"acp: not configured",
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
	m := newTestModel(nil)

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
	m := newTestModel(testRuns())
	m.focus = focusAgents

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
	if m.selectedAgent != 0 {
		t.Fatalf("selectedAgent=%d, want reset 0 for selected run", m.selectedAgent)
	}
}

func TestDashboardLaunchModeOverridesAreSessionOnly(t *testing.T) {
	m := newTestModel(nil)
	m.width = 180

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
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(model)
	if got := m.launchOverrides["codex"]; got != agents.LaunchACP {
		t.Fatalf("override=%q, want acp", got)
	}
	if got := agents.LaunchModeOrDefault(m.agents[0].LaunchMode); got != agents.LaunchHeadless {
		t.Fatalf("underlying launch mode mutated to %q", got)
	}
	if !strings.Contains(m.View(), "effective: acp") || !strings.Contains(m.View(), "session only") {
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
	m := newTestModel(nil)
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

func TestDashboardACPModeRequiresACPConfig(t *testing.T) {
	m := newTestModel(nil)
	m.focus = focusAgents
	m.selectedAgent = 1

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(model)
	if _, ok := m.launchOverrides["claude"]; ok {
		t.Fatalf("unsupported ACP mode created override: %+v", m.launchOverrides)
	}
	if !strings.Contains(m.View(), "claude has no acp launch configuration") {
		t.Fatalf("missing unsupported ACP message\n%s", m.View())
	}
}

func TestWorkspaceRendersRunsEventsAndQuestions(t *testing.T) {
	m := newTestModel(testRuns())
	m.width = 140

	view := m.View()
	for _, want := range []string{
		"Sessions",
		"ACTION",
		"sample",
		"Event stream",
		"agent.started",
		"Run details",
		"Questions",
		"Which branch?",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("workspace view missing %q\n%s", want, view)
		}
	}
}

func TestDashboardCompactLayoutFitsShortTerminal(t *testing.T) {
	m := newTestModel(testRuns())
	m.width = 100
	m.height = 20

	view := m.View()
	for _, want := range []string{"layout=compact", "Sessions", "Run details", "Agents", "Actions", "Questions", "q/esc quit"} {
		if !strings.Contains(view, want) {
			t.Fatalf("compact dashboard missing %q\n%s", want, view)
		}
	}
	if got := renderedLineCount(view); got > m.height {
		t.Fatalf("compact dashboard rendered %d lines, want <= %d\n%s", got, m.height, view)
	}
}

func TestDashboardCompactThresholdBoundary(t *testing.T) {
	m := newTestModel(testRuns())
	m.width = 100
	m.height = compactDashboardHeight
	if view := m.View(); strings.Contains(view, "layout=compact") {
		t.Fatalf("dashboard used compact layout at threshold height %d\n%s", compactDashboardHeight, view)
	}

	m.height = compactDashboardHeight - 1
	if view := m.View(); !strings.Contains(view, "layout=compact") {
		t.Fatalf("dashboard did not use compact layout below threshold height %d\n%s", compactDashboardHeight-1, view)
	}
}

func TestDashboardActionFocusAndSelection(t *testing.T) {
	m := newTestModel(testRuns())

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.focus != focusActions {
		t.Fatalf("focus=%s, want actions", m.focus)
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.selectedAction != 1 {
		t.Fatalf("selectedAction=%d, want 1", m.selectedAction)
	}
	view := m.View()
	if !strings.Contains(view, "> answer-question") {
		t.Fatalf("view did not mark selected action\n%s", view)
	}
	if !strings.Contains(view, "enter action") {
		t.Fatalf("footer missing action key\n%s", view)
	}
}

func TestDashboardEnterRunsSelectedAction(t *testing.T) {
	var got ActionRequest
	refreshCalled := false
	m := newModel(WorkspaceOptions{
		Root:   "/repo",
		Status: testStatus(),
		Agents: testDiscoveries(),
		Runs:   testRuns(),
		RefreshRuns: func() ([]runstate.RunSummary, error) {
			refreshCalled = true
			return testRuns(), nil
		},
		ActionRunner: func(_ context.Context, request ActionRequest) (ActionResult, error) {
			got = request
			return ActionResult{Message: "action done", Refresh: true}, nil
		},
	})
	m.focus = focusActions
	m.selectedAction = 1

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if cmd == nil {
		t.Fatal("enter did not return action command")
	}
	if !m.actionRunning {
		t.Fatal("model did not enter action running state")
	}
	updated, refreshCmd := m.Update(cmd())
	m = updated.(model)
	if got.Action.ID != "answer-question.q1" {
		t.Fatalf("action=%+v", got.Action)
	}
	if m.actionRunning {
		t.Fatal("model stayed actionRunning after result")
	}
	if !strings.Contains(m.actionText, "action done") {
		t.Fatalf("actionText=%q", m.actionText)
	}
	if refreshCmd == nil {
		t.Fatal("successful action did not request refresh")
	}
	updated, _ = m.Update(refreshCmd())
	m = updated.(model)
	if !refreshCalled {
		t.Fatal("refresh callback was not called")
	}
}

func TestDashboardActionWithoutRunnerShowsCommand(t *testing.T) {
	m := newTestModel(testRuns())
	m.focus = focusActions
	m.selectedAction = 1

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if cmd != nil {
		t.Fatal("enter returned command without runner")
	}
	for _, want := range []string{
		"action execution is not available",
		"parley answer run-1 q1 <answer>",
	} {
		if !strings.Contains(m.actionText, want) {
			t.Fatalf("actionText missing %q: %q", want, m.actionText)
		}
	}
}

func TestDashboardRequiresConfirmationBeforeRunningAction(t *testing.T) {
	calls := 0
	m := newModel(WorkspaceOptions{
		Root:   "/repo",
		Status: testStatus(),
		Agents: testDiscoveries(),
		Runs:   testRuns(),
		ActionRunner: func(_ context.Context, request ActionRequest) (ActionResult, error) {
			calls++
			return ActionResult{Message: request.Action.ID}, nil
		},
	})
	m.focus = focusActions

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if cmd != nil {
		t.Fatal("first enter returned command for RequiresYes action")
	}
	if calls != 0 {
		t.Fatalf("runner called before confirmation: %d", calls)
	}
	if !strings.Contains(m.actionText, "press enter again") {
		t.Fatalf("actionText=%q", m.actionText)
	}

	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if cmd == nil {
		t.Fatal("second enter did not return action command")
	}
	updated, _ = m.Update(cmd())
	m = updated.(model)
	if calls != 1 {
		t.Fatalf("runner calls=%d, want 1", calls)
	}
	if strings.Contains(m.actionText, "press enter again") {
		t.Fatalf("confirmation text was not replaced: %q", m.actionText)
	}
}

func TestDashboardBlocksEnterWhileActionRunning(t *testing.T) {
	calls := 0
	m := newModel(WorkspaceOptions{
		Root:   "/repo",
		Status: testStatus(),
		Agents: testDiscoveries(),
		Runs:   testRuns(),
		ActionRunner: func(_ context.Context, request ActionRequest) (ActionResult, error) {
			calls++
			return ActionResult{Message: request.Action.ID}, nil
		},
	})
	m.focus = focusActions
	m.selectedAction = 1
	m.actionRunning = true

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if cmd != nil {
		t.Fatal("enter returned command while actionRunning")
	}
	if calls != 0 {
		t.Fatalf("runner calls=%d, want 0", calls)
	}
	if !strings.Contains(m.actionText, "already running") {
		t.Fatalf("actionText=%q", m.actionText)
	}
}

func TestWorkspaceStartModeUsesCallback(t *testing.T) {
	var got StartRequest
	m := newModel(WorkspaceOptions{
		Root:   "/repo",
		Status: testStatus(),
		Agents: testDiscoveries(),
		StartRun: func(_ context.Context, request StartRequest) (runstate.RunSummary, error) {
			got = request
			return runstate.RunSummary{RunID: "run-new", IdeaSlug: "new-idea", Attention: runstate.AttentionIdle}, nil
		},
	})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
	m = updated.(model)
	if !m.startMode {
		t.Fatal("N did not enter start mode")
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("New task")})
	m = updated.(model)
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if cmd == nil {
		t.Fatal("enter did not return start command")
	}
	updated, _ = m.Update(cmd())
	m = updated.(model)
	if got.Task != "New task" || len(got.Participants) != 1 || got.Participants[0] != "codex" {
		t.Fatalf("request=%+v", got)
	}
	if len(m.runs) != 1 || m.runs[0].RunID != "run-new" {
		t.Fatalf("runs=%+v", m.runs)
	}
}

func TestRefreshPreservesSelectedRunAndAgent(t *testing.T) {
	runs := testRuns()
	m := newTestModel(runs)
	m.selectedIdea = 0
	m.selectedAgent = 1

	newRun := runstate.RunSummary{RunID: "run-newer", IdeaSlug: "newer", Attention: runstate.AttentionRunning}
	updated, _ := m.Update(refreshRunsMsg{runs: append([]runstate.RunSummary{newRun}, runs...)})
	m = updated.(model)

	run, ok := m.selectedRun()
	if !ok || run.RunID != "run-1" {
		t.Fatalf("selected run=%+v ok=%v, want run-1", run, ok)
	}
	if m.selectedAgent != 1 {
		t.Fatalf("selectedAgent=%d, want preserved 1", m.selectedAgent)
	}
}

func newTestModel(runs []runstate.RunSummary) model {
	return newModel(WorkspaceOptions{Root: "/repo", Status: testStatus(), Agents: testDiscoveries(), Runs: runs})
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
				HeadlessMode:          "codex exec --skip-git-repo-check -",
				HeadlessArgs:          []string{"exec", "--skip-git-repo-check", "-"},
				ACPArgs:               []string{"acp"},
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

func testRuns() []runstate.RunSummary {
	base := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	return []runstate.RunSummary{
		{
			RunID:         "run-1",
			IdeaSlug:      "sample",
			Task:          "Sample task",
			Participants:  []string{"codex"},
			OpenQuestions: 1,
			Attention:     runstate.AttentionAction,
			LastEventAt:   base,
			LastEventAge:  time.Minute,
			State: runstate.RunState{
				Agents: []runstate.AgentState{
					{ID: "codex", State: runstate.StateRunning, StartedAt: base, LatestEvent: "agent.started"},
					{ID: "claude", State: runstate.StatePending},
				},
				Recent: []runstate.EventSummary{{Time: base, Type: "agent.started", Agent: "codex", Text: "codex"}},
			},
			NextActions: []runaction.NextAction{
				{ID: "draft-consensus.sample", Kind: runaction.KindDraftConsensus, IdeaSlug: "sample", Risk: runaction.RiskNormal, RequiresYes: true, Summary: "Draft consensus"},
				{ID: "answer-question.q1", Kind: runaction.KindAnswerQuestion, RunID: "run-1", IdeaSlug: "sample", Risk: runaction.RiskLow, Summary: "Answer HITL question"},
			},
			Questions: []hitl.Question{{ID: "q1", Agent: "codex", Prompt: "Which branch?", Risk: hitl.RiskNormal, Status: hitl.StatusOpen}},
		},
		{
			RunID:       "run-2",
			IdeaSlug:    "second",
			Attention:   runstate.AttentionIdle,
			LastEventAt: base.Add(-time.Hour),
		},
	}
}

func renderedLineCount(view string) int {
	view = strings.TrimRight(view, "\n")
	if view == "" {
		return 0
	}
	return strings.Count(view, "\n") + 1
}
