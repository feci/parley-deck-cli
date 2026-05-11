package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestProjectEventsDerivesAgentAndRoundState(t *testing.T) {
	base := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		{Time: base, Type: "agent.started", Data: map[string]any{"agent": "codex", "stdout": "/tmp/codex.out", "stderr": "/tmp/codex.err"}},
		{Time: base.Add(time.Second), Type: "agent.finished", Data: map[string]any{"agent": "codex", "duration_ms": float64(1500)}},
		{Time: base.Add(2 * time.Second), Type: "agent.failed", Data: map[string]any{"agent": "claude", "error": "exit status 1", "duration_ms": float64(2000)}},
		{Time: base.Add(3 * time.Second), Type: "agent.started", Data: map[string]any{"agent": "gemini"}},
		{Time: base.Add(4 * time.Second), Type: "agent.finished", Data: map[string]any{"agent": "gemini", "duration_ms": float64(50)}},
		{Time: base.Add(5 * time.Second), Type: "agent.skipped", Data: map[string]any{"agent": "hermes", "reason": "artifact already exists"}},
		{Time: base.Add(6 * time.Second), Type: "round.incomplete", Data: map[string]any{"completed": float64(1), "total": float64(4)}},
	}

	state := ProjectEvents([]string{"codex", "claude", "hermes", "opus"}, events, base.Add(7*time.Second))
	agents := mapByID(state.Agents)

	if got := agents["codex"].State; got != stateFinished {
		t.Fatalf("codex state=%s, want %s", got, stateFinished)
	}
	if got := agents["codex"].Duration; got != 1500*time.Millisecond {
		t.Fatalf("codex duration=%s, want 1.5s", got)
	}
	if got := agents["claude"].State; got != stateFailed {
		t.Fatalf("claude state=%s, want %s", got, stateFailed)
	}
	if got := agents["hermes"].State; got != stateSkipped {
		t.Fatalf("hermes state=%s, want %s", got, stateSkipped)
	}
	if got := agents["hermes"].Duration; got != 0 {
		t.Fatalf("hermes duration=%s, want zero for skip without start", got)
	}
	if got := agents["opus"].State; got != statePending {
		t.Fatalf("opus state=%s, want %s", got, statePending)
	}
	if got := agents["gemini"].State; got != stateUnknown {
		t.Fatalf("gemini state=%s, want %s", got, stateUnknown)
	}
	if got := state.RoundStatus; got != "incomplete" {
		t.Fatalf("round status=%s, want incomplete", got)
	}
}

func TestReadEventsFromOffsetKeepsPartialLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.jsonl")
	first := `{"time":"2026-05-10T12:00:00Z","type":"agent.started","data":{"agent":"codex"}}` + "\n"
	partial := `{"time":"2026-05-10T12:00:01Z","type":"agent.finished"`
	if err := os.WriteFile(path, []byte(first+partial), 0o644); err != nil {
		t.Fatal(err)
	}

	events, offset, err := readEventsFromOffset(path, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Type != "agent.started" {
		t.Fatalf("events=%v, want one started event", events)
	}
	if offset != int64(len(first)) {
		t.Fatalf("offset=%d, want %d", offset, len(first))
	}

	complete := `,"data":{"agent":"codex","duration_ms":1000}}` + "\n"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString(complete); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	events, _, err = readEventsFromOffset(path, offset)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Type != "agent.finished" {
		t.Fatalf("events=%v, want one finished event", events)
	}
}

func TestTailLogFileReturnsBoundedWholeLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stdout.log")
	content := "partial first line\nline 1\nline 2\n\x1b[31mline 3\x1b[0m\nunfinished"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := tailLogFile(path, 64, 2)
	if strings.Contains(got, "unfinished") {
		t.Fatalf("tail included partial trailing line: %q", got)
	}
	if strings.Contains(got, "\x1b") {
		t.Fatalf("tail included ansi escape: %q", got)
	}
	if got != "line 2\nline 3" {
		t.Fatalf("tail=%q, want last two whole lines", got)
	}
}

func TestLiveViewIncludesRequiredPanels(t *testing.T) {
	model := newLiveModel(LiveOptions{
		Idea:         testIdea("live-run-tui"),
		Participants: []string{"codex"},
		RunID:        "run-1",
		RunDir:       t.TempDir(),
	})
	model.width = 100
	model.events = []store.Event{
		{Time: time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC), Type: "agent.started", Data: map[string]any{"agent": "codex"}},
	}
	model.state = ProjectEvents([]string{"codex"}, model.events, model.now)

	view := model.View()
	for _, want := range []string{"idea=live-run-tui", "Agents", "Latest events", "Questions", "Log preview", "q/esc detach", "ctrl+c cancel"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view missing %q\n%s", want, view)
		}
	}
}

func TestLiveQuestionsPaneAndAnswerMode(t *testing.T) {
	runDir := t.TempDir()
	question, err := hitl.New(runDir).Create(hitl.Question{
		Agent:  "codex",
		Prompt: "Which branch should I target?",
		Risk:   hitl.RiskNormal,
	})
	if err != nil {
		t.Fatal(err)
	}
	model := newLiveModel(LiveOptions{
		Idea:         testIdea("hitl-tui-questions"),
		Participants: []string{"codex"},
		RunID:        "run-1",
		RunDir:       runDir,
	})

	updated, _ := model.Update(questionsMsg{questions: []hitl.Question{question}})
	model = updated.(liveModel)
	if !strings.Contains(model.View(), "Which branch should I target?") {
		t.Fatalf("view missing question\n%s", model.View())
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updated.(liveModel)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("main")})
	model = updated.(liveModel)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)

	questions, err := hitl.New(runDir).List()
	if err != nil {
		t.Fatal(err)
	}
	if got := questions[0].Answer; got != "main" {
		t.Fatalf("answer=%q, want main", got)
	}
	if questions[0].Status != hitl.StatusAnswered {
		t.Fatalf("status=%s, want answered", questions[0].Status)
	}
}

func mapByID(agents []AgentState) map[string]AgentState {
	mapped := map[string]AgentState{}
	for _, agent := range agents {
		mapped[agent.ID] = agent
	}
	return mapped
}

func testIdea(slug string) protocol.IdeaStatus {
	return protocol.IdeaStatus{Slug: slug, Status: "round-01"}
}
