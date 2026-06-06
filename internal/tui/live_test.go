package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/steer"
	"parley-deck-cli/internal/store"
)

func TestProjectEventsDerivesAgentAndRoundState(t *testing.T) {
	base := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		{Time: base, Type: "agent.started", Data: map[string]any{"agent": "codex", "stdout": "/tmp/codex.out", "stderr": "/tmp/codex.err"}},
		{Time: base.Add(time.Second), Type: "agent.finished", Data: map[string]any{"agent": "codex", "duration_ms": float64(1500)}},
		{Time: base.Add(2 * time.Second), Type: "agent.failed", Data: map[string]any{"agent": "claude", "error": "exit status 1", "duration_ms": float64(2000)}},
		{Time: base.Add(3 * time.Second), Type: "agent.started", Data: map[string]any{"agent": "agy"}},
		{Time: base.Add(4 * time.Second), Type: "agent.finished", Data: map[string]any{"agent": "agy", "duration_ms": float64(50)}},
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
	if got := agents["opus"].State; got != statePending {
		t.Fatalf("opus state=%s, want %s", got, statePending)
	}
	if got := agents["agy"].State; got != stateUnknown {
		t.Fatalf("agy state=%s, want %s", got, stateUnknown)
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
	appendString(t, path, complete)
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
	if strings.Contains(got, "unfinished") || strings.Contains(got, "\x1b") {
		t.Fatalf("tail leaked partial/ansi: %q", got)
	}
	if got != "line 2\nline 3" {
		t.Fatalf("tail=%q, want last two whole lines", got)
	}
}

func TestSummarizeHITLEvents(t *testing.T) {
	question := summarizeEvent(store.Event{Type: "hitl.question", Data: map[string]any{"agent": "agy", "question_id": "q1", "risk": "normal"}})
	if question.Text != "agy question q1 normal" {
		t.Fatalf("question text=%q", question.Text)
	}
	answered := summarizeEvent(store.Event{Type: "hitl.answered", Data: map[string]any{"agent": "agy", "question_id": "q1", "status": "answered"}})
	if answered.Text != "agy answered q1 answered" {
		t.Fatalf("answered text=%q", answered.Text)
	}
}

// --- focus-read pipeline (reused by the per-agent transcript buffers) ---

func TestLoadFocusTailBoundsLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stdout.log")
	var sb strings.Builder
	total := maxFocusLines + 50
	for i := 0; i < total; i++ {
		fmt.Fprintf(&sb, "L%d\n", i)
	}
	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	lines, _, truncated := loadFocusTail(path)
	if len(lines) != maxFocusLines || !truncated {
		t.Fatalf("len=%d truncated=%v, want %d/true", len(lines), truncated, maxFocusLines)
	}
	if want := fmt.Sprintf("L%d", total-1); lines[len(lines)-1] != want {
		t.Fatalf("last line=%q, want %q", lines[len(lines)-1], want)
	}
}

func TestReadAppendedLinesIncremental(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stdout.log")
	if err := os.WriteFile(path, []byte("a\nb\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, offset, _ := loadFocusTail(path)
	if offset != 4 {
		t.Fatalf("offset=%d, want 4", offset)
	}
	appendString(t, path, "c\npartial")
	lines, newOffset := readAppendedLines(path, offset)
	if len(lines) != 1 || lines[0] != "c" || newOffset != offset+2 {
		t.Fatalf("lines=%v offset=%d, want [c]/%d", lines, newOffset, offset+2)
	}
	appendString(t, path, " line\n")
	lines, _ = readAppendedLines(path, newOffset)
	if len(lines) != 1 || lines[0] != "partial line" {
		t.Fatalf("lines=%v, want [partial line]", lines)
	}
}

func TestFocusTailDropsOversizedLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stdout.log")
	big := strings.Repeat("x", maxFocusBytes+4096)
	if err := os.WriteFile(path, []byte(big+"\nshort tail line\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lines, _, truncated := loadFocusTail(path)
	total := 0
	for _, l := range lines {
		total += len(l) + 1
	}
	if total > maxFocusBytes || !truncated {
		t.Fatalf("total=%d truncated=%v, want <=%d/true", total, truncated, maxFocusBytes)
	}
	if len(lines) == 0 || lines[len(lines)-1] != "short tail line" {
		t.Fatalf("lines=%v, want short line retained", lines)
	}
}

func TestCapFocusLinesByteBudget(t *testing.T) {
	line := strings.Repeat("y", 1000)
	n := (maxFocusBytes / 1001) + 100
	lines := make([]string, n)
	for i := range lines {
		lines[i] = line
	}
	capped, truncated := capFocusLines(lines)
	total := 0
	for _, l := range capped {
		total += len(l) + 1
	}
	if !truncated || total > maxFocusBytes {
		t.Fatalf("truncated=%v total=%d, want true/<=%d", truncated, total, maxFocusBytes)
	}
}

func TestCapFocusLinesTruncatesSingleOversizedLine(t *testing.T) {
	capped, truncated := capFocusLines([]string{strings.Repeat("z", maxFocusBytes+500)})
	total := 0
	for _, l := range capped {
		total += len(l) + 1
	}
	if !truncated || total > maxFocusBytes {
		t.Fatalf("single oversized line not capped: truncated=%v total=%d", truncated, total)
	}
}

func appendString(t *testing.T, path, s string) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString(s); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

// --- Claude-CLI tabbed layout ---

// liveModelWithLog builds a live model with one running agent "codex" whose
// stdout log holds content, with the active transcript buffer loaded.
func liveModelWithLog(t *testing.T, content string) liveModel {
	t.Helper()
	dir := t.TempDir()
	logPath := filepath.Join(dir, "stdout.log")
	if err := os.WriteFile(logPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	model := newLiveModel(LiveOptions{
		Idea:         testIdea("tui-claude-cli-layout"),
		Participants: []string{"codex"},
		RunID:        "run-1",
		RunDir:       dir,
	})
	model.height, model.width = 30, 100
	model.events = []store.Event{
		{Time: base, Type: "agent.started", Data: map[string]any{"agent": "codex", "stdout": logPath, "segment_id": "segment-0001"}},
	}
	model.state = ProjectEvents([]string{"codex"}, model.events, model.now)
	model.ensureActiveBuffer()
	return model
}

func TestTabbedDefaultShowsTranscriptAndInput(t *testing.T) {
	model := liveModelWithLog(t, "hello from codex\nworking on the parser\nDONE\n")
	if got := model.activeTabResolved(); got != "agent:codex" {
		t.Fatalf("default active tab=%q, want agent:codex (running agent transcript)", got)
	}
	view := model.View()
	for _, want := range []string{"codex", "working on the parser", "steer codex ›", "↑/↓ tabs"} {
		if !strings.Contains(view, want) {
			t.Fatalf("default view missing %q\n%s", want, view)
		}
	}
}

func TestTabSwitchWithUpDownArrows(t *testing.T) {
	model := liveModelWithLog(t, "x\n") // tabs: [agent:codex, status]
	if got := model.activeTabResolved(); got != "agent:codex" {
		t.Fatalf("default=%q", got)
	}
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(liveModel)
	if got := model.activeTabResolved(); got != statusTabID {
		t.Fatalf("after ↓ = %q, want status", got)
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updated.(liveModel)
	if got := model.activeTabResolved(); got != "agent:codex" {
		t.Fatalf("after ↑ = %q, want agent:codex", got)
	}
}

func TestTranscriptScrollAndFollow(t *testing.T) {
	var sb strings.Builder
	for i := 1; i <= 100; i++ {
		fmt.Fprintf(&sb, "line %d\n", i)
	}
	model := liveModelWithLog(t, sb.String())
	if b := model.buffers["codex"]; b == nil || !b.follow {
		t.Fatal("active buffer should be loaded and following")
	}
	if !strings.Contains(model.View(), "line 100") {
		t.Fatalf("follow should show the last line\n%s", model.View())
	}
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	model = updated.(liveModel)
	if model.buffers["codex"].follow {
		t.Fatal("PgUp must drop follow")
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnd})
	model = updated.(liveModel)
	if !model.buffers["codex"].follow {
		t.Fatal("End must re-enable follow")
	}
	if !strings.Contains(model.View(), "line 100") {
		t.Fatalf("End should show the last line\n%s", model.View())
	}
}

func TestInputSteersActiveAgent(t *testing.T) {
	model := liveModelWithLog(t, "x\n")
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("focus on the parser")})
	model = updated.(liveModel)
	if model.inputText != "focus on the parser" {
		t.Fatalf("inputText=%q", model.inputText)
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)
	// No SubmitSteer seam here → record-only path; the queued list below verifies
	// the steer was persisted for the active agent.
	if !strings.Contains(model.statusMsg, "recorded steer steer-0001-") {
		t.Fatalf("statusMsg=%q", model.statusMsg)
	}
	queued, err := steer.List(model.opts.RunDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(queued) != 1 || queued[0].Agent != "codex" || queued[0].Text != "focus on the parser" || queued[0].SegmentID != "segment-0001" {
		t.Fatalf("queued=%+v", queued)
	}
}

func TestAnswerViaInputRowWhenQuestionOpen(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	runDir := t.TempDir()
	logPath := filepath.Join(runDir, "stdout.log")
	if err := os.WriteFile(logPath, []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	q, err := hitl.New(runDir).Create(hitl.Question{Agent: "codex", Prompt: "Which branch?", Risk: hitl.RiskNormal})
	if err != nil {
		t.Fatal(err)
	}
	model := newLiveModel(LiveOptions{Idea: testIdea("hitl"), Participants: []string{"codex"}, RunID: "run-1", RunDir: runDir})
	model.height, model.width = 30, 100
	model.events = []store.Event{{Time: base, Type: "agent.started", Data: map[string]any{"agent": "codex", "stdout": logPath, "segment_id": "segment-0001"}}}
	model.state = ProjectEvents([]string{"codex"}, model.events, model.now)
	model.ensureActiveBuffer()

	updated, _ := model.Update(questionsMsg{questions: []hitl.Question{q}})
	model = updated.(liveModel)
	if !strings.Contains(model.View(), "answer codex/"+q.ID) {
		t.Fatalf("input row not in answer mode\n%s", model.View())
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("main")})
	model = updated.(liveModel)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)

	questions, err := hitl.New(runDir).List()
	if err != nil {
		t.Fatal(err)
	}
	if questions[0].Answer != "main" || questions[0].Status != hitl.StatusAnswered {
		t.Fatalf("question=%+v, want answered=main", questions[0])
	}
}

// D9 routing table: printable keys append to input and never trigger legacy
// single-letter actions or change mode.
func TestKeyRoutingPrintableAppendsNotHotkey(t *testing.T) {
	model := liveModelWithLog(t, "x\n")
	for _, r := range []rune{'q', '?', 'a', 'f', 'j', 'k', 'i'} {
		updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = updated.(liveModel)
		if cmd != nil {
			t.Fatalf("printable %q must not trigger a command/quit", r)
		}
	}
	if model.inputText != "q?afjki" {
		t.Fatalf("inputText=%q, want all chars appended", model.inputText)
	}
	if model.mode != modeOverview {
		t.Fatalf("printable keys must not change mode, got %d", model.mode)
	}
}

func TestSlashHelpOpensOverlay(t *testing.T) {
	model := liveModelWithLog(t, "x\n")
	for _, r := range "/help" {
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = updated.(liveModel)
	}
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)
	if model.mode != modeHelp {
		t.Fatalf("/help should open the help overlay, mode=%d", model.mode)
	}
	if !strings.Contains(model.View(), "Help") {
		t.Fatalf("help overlay missing 'Help'\n%s", model.View())
	}
}

func TestEscClearsInputThenDetaches(t *testing.T) {
	model := liveModelWithLog(t, "x\n")
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hi")})
	model = updated.(liveModel)
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(liveModel)
	if model.inputText != "" || cmd != nil {
		t.Fatalf("esc with text should clear input and not quit (input=%q cmd=%v)", model.inputText, cmd)
	}
	if _, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc}); cmd == nil {
		t.Fatal("esc with empty input should detach (quit)")
	}
}

func TestInputBackspaceRemovesWholeRune(t *testing.T) {
	model := liveModelWithLog(t, "x\n")
	model.inputText = "á"
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = updated.(liveModel)
	if model.inputText != "" {
		t.Fatalf("inputText=%q, want empty", model.inputText)
	}
}

func TestResumeStatusAndEscDetach(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	model := newLiveModel(LiveOptions{Idea: testIdea("resume"), Participants: []string{"codex"}, RunID: "run-1", RunDir: t.TempDir(), Resume: true})
	model.height, model.width = 30, 100
	model.events = []store.Event{{Time: base, Type: "agent.started", Data: map[string]any{"agent": "codex"}}}
	model.state = ProjectEvents([]string{"codex"}, model.events, base.Add(time.Minute))
	view := model.View()
	if strings.Contains(view, "round=running") {
		t.Fatalf("resume view must not imply a live process\n%s", view)
	}
	if !strings.Contains(view, "unverified") {
		t.Fatalf("resume view should show unverified status\n%s", view)
	}
	if _, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc}); cmd == nil {
		t.Fatal("esc on empty input should detach")
	}
}

// AF1: the active tab stays visible (with a …+N marker) when the strip overflows.
func TestTabStripKeepsActiveTabVisible(t *testing.T) {
	parts := make([]string, 0, 8)
	for i := 0; i < 8; i++ {
		parts = append(parts, fmt.Sprintf("agent%d", i))
	}
	model := newLiveModel(LiveOptions{Idea: testIdea("x"), Participants: parts, RunID: "r", RunDir: t.TempDir()})
	model.width = 40
	model.state = ProjectEvents(parts, nil, model.now)
	model.activeTab = "agent:agent7" // far right
	strip := stripANSI(model.renderTabStrip(40))
	if !strings.Contains(strip, "agent7") {
		t.Fatalf("active tab agent7 not visible in narrow strip: %q", strip)
	}
	if !strings.Contains(strip, "…+") {
		t.Fatalf("expected an overflow marker in: %q", strip)
	}
}

// AF4: shift+↑ scrolls one line and drops follow (plain ↑/↓ stay tab switches).
func TestShiftArrowLineScroll(t *testing.T) {
	var sb strings.Builder
	for i := 1; i <= 100; i++ {
		fmt.Fprintf(&sb, "line %d\n", i)
	}
	model := liveModelWithLog(t, sb.String())
	if !model.buffers["codex"].follow {
		t.Fatal("buffer should start following")
	}
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyShiftUp})
	model = updated.(liveModel)
	if model.buffers["codex"].follow {
		t.Fatal("shift+↑ must drop follow (line scroll)")
	}
}

// AF5: slash commands route — /status switches tab, /deck records a deck steer.
func TestSlashDeckAndStatus(t *testing.T) {
	model := liveModelWithLog(t, "x\n")
	model.inputText = "/status"
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)
	if model.activeTabResolved() != statusTabID {
		t.Fatalf("/status should switch to the Status tab, got %q", model.activeTabResolved())
	}
	model.inputText = "/deck wrap it up"
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)
	queued, err := steer.List(model.opts.RunDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(queued) != 1 || queued[0].Target != "deck" || queued[0].Text != "wrap it up" {
		t.Fatalf("/deck queued=%+v", queued)
	}
}

// AF2: a replaced log file (different inode, grown past the old offset) reloads.
func TestBufferReloadsOnFileReplace(t *testing.T) {
	model := liveModelWithLog(t, "old line\n")
	b := model.buffers["codex"]
	if b == nil || !b.loaded {
		t.Fatal("active buffer should be loaded")
	}
	path := b.path
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("fresh line one\nfresh line two\nfresh line three\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	model.refreshBuffers()
	joined := strings.Join(model.buffers["codex"].lines, "\n")
	if !strings.Contains(joined, "fresh line three") || strings.Contains(joined, "old line") {
		t.Fatalf("buffer did not reload from the replaced file: %q", joined)
	}
}

// unified-tui-home: Home is the default surface when no run is attached.
func TestHomeDefaultWhenNoRunAndTabOrder(t *testing.T) {
	model := newLiveModel(LiveOptions{Home: true, Status: protocol.WorkspaceStatus{}, Root: t.TempDir()})
	model.height, model.width = 30, 100
	if got := model.activeTabResolved(); got != homeTabID {
		t.Fatalf("no-run default tab = %q, want home", got)
	}
	if ids := model.tabIDs(); len(ids) != 1 || ids[0] != homeTabID {
		t.Fatalf("no-run tabIDs = %v, want [home]", ids)
	}
	view := model.View()
	for _, want := range []string{"Home", "Ideas", "Recent runs", "N new idea"} {
		if !strings.Contains(view, want) {
			t.Fatalf("Home view missing %q\n%s", want, view)
		}
	}
}

func TestActiveRunTabOrderHomeFirst(t *testing.T) {
	model := liveModelWithLog(t, "x\n") // running codex, hasRun
	ids := model.tabIDs()
	if len(ids) < 2 || ids[0] != homeTabID || ids[len(ids)-1] != statusTabID {
		t.Fatalf("active tabIDs = %v, want home first + status last", ids)
	}
	if got := model.activeTabResolved(); got != "agent:codex" {
		t.Fatalf("parley-run default = %q, want agent:codex", got)
	}
}

func TestDoneDoesNotQuit(t *testing.T) {
	model := liveModelWithLog(t, "x\n")
	updated, _ := model.Update(doneMsg{token: model.runToken})
	model = updated.(liveModel)
	if !model.done {
		t.Fatal("doneMsg should set done")
	}
	if _, cmd := model.Update(eventsMsg{token: model.runToken}); cmd != nil {
		t.Fatal("after done, an eventsMsg must not quit (cmd should be nil)")
	}
}

func TestNLaunchesViaStartFunc(t *testing.T) {
	runDir := t.TempDir()
	called := false
	model := newLiveModel(LiveOptions{
		Home: true, Root: t.TempDir(),
		Start: func(req LaunchRequest) (LaunchResult, error) {
			called = true
			return LaunchResult{Idea: protocol.IdeaStatus{Slug: "x"}, Participants: []string{"codex"}, RunID: "run-x", RunDir: runDir}, nil
		},
	})
	model.height, model.width = 30, 100
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
	model = updated.(liveModel)
	if !model.composing {
		t.Fatal("N should open the new-idea composer")
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("redesign the parser")})
	model = updated.(liveModel)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)
	if !called {
		t.Fatal("Enter in compose should call Start")
	}
	if !model.hasRun() || model.opts.RunID != "run-x" {
		t.Fatalf("after launch the model should attach to run-x, got RunDir=%q", model.opts.RunDir)
	}
	if got := model.activeTabResolved(); got != "agent:codex" {
		t.Fatalf("after launch active tab = %q, want agent:codex", got)
	}
}

func TestRunTokenDropsStaleEvents(t *testing.T) {
	model := liveModelWithLog(t, "x\n")
	stale := model.runToken
	model.activateRun(LaunchResult{Participants: []string{"claude"}, RunID: "r2", RunDir: t.TempDir()})
	before := len(model.events)
	updated, _ := model.Update(eventsMsg{token: stale, events: []store.Event{{Type: "agent.started", Data: map[string]any{"agent": "codex"}}}})
	model = updated.(liveModel)
	if len(model.events) != before {
		t.Fatal("a stale-token eventsMsg must be ignored after a run swap")
	}
}

// AF3 / owner #3: a real on-disk run (events.jsonl with agent.started carrying a
// stdout path + a stdout.log) must show that agent's live transcript.
func TestTranscriptPopulatesFromOnDiskRun(t *testing.T) {
	runDir := t.TempDir()
	stdoutPath := filepath.Join(runDir, "codex-stdout.log")
	if err := os.WriteFile(stdoutPath, []byte("reading the parser\nwriting changes\nDONE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	s := store.New(runDir)
	for _, e := range []store.Event{
		{Time: base, Type: "run.created", Data: map[string]any{"idea": "x", "participants": []string{"codex"}}},
		{Time: base.Add(time.Second), Type: "run.segment_started", Data: map[string]any{"segment_id": "segment-0001", "reason": "initial", "targets": []string{"codex"}}},
		{Time: base.Add(2 * time.Second), Type: "agent.started", Data: map[string]any{"agent": "codex", "stdout": stdoutPath, "segment_id": "segment-0001"}},
	} {
		if err := s.Append(e); err != nil {
			t.Fatal(err)
		}
	}

	model := newLiveModel(LiveOptions{Idea: testIdea("x"), Participants: []string{"codex"}, RunID: "run-x", RunDir: runDir})
	model.height, model.width = 30, 100
	events, off, err := readEventsFromOffset(filepath.Join(runDir, "events.jsonl"), 0)
	if err != nil {
		t.Fatal(err)
	}
	updated, _ := model.Update(eventsMsg{events: events, offset: off, token: model.runToken})
	model = updated.(liveModel)

	if got := model.activeTabResolved(); got != "agent:codex" {
		t.Fatalf("active tab = %q, want agent:codex", got)
	}
	if b := model.buffers["codex"]; b == nil || len(b.lines) == 0 {
		t.Fatalf("transcript buffer is empty: %+v", b)
	}
	if view := model.View(); !strings.Contains(view, "writing changes") {
		t.Fatalf("transcript not shown in the view:\n%s", view)
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

// --- tui-command-picker tests ---

func pressRunes(t *testing.T, m liveModel, s string) liveModel {
	t.Helper()
	u, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)})
	return u.(liveModel)
}

func pressKey(t *testing.T, m liveModel, k tea.KeyType) (liveModel, tea.Cmd) {
	t.Helper()
	u, c := m.Update(tea.KeyMsg{Type: k})
	return u.(liveModel), c
}

// runCmd types a slash command into the input row and presses Enter.
func runCmd(t *testing.T, m liveModel, cmd string) (liveModel, tea.Cmd) {
	t.Helper()
	m = pressRunes(t, m, cmd)
	return pressKey(t, m, tea.KeyEnter)
}

// homeModelWithIdeas is a no-run Home model carrying the given idea slugs.
func homeModelWithIdeas(t *testing.T, slugs ...string) liveModel {
	t.Helper()
	ideas := make([]protocol.IdeaStatus, 0, len(slugs))
	for _, s := range slugs {
		ideas = append(ideas, protocol.IdeaStatus{Slug: s, Status: "final"})
	}
	m := newLiveModel(LiveOptions{Home: true, Root: t.TempDir(), Status: protocol.WorkspaceStatus{Ideas: ideas}})
	m.height, m.width = 30, 100
	return m
}

// runModelWithQuestion is a model attached to a run with one open question.
func runModelWithQuestion(t *testing.T) (liveModel, string, hitl.Question) {
	t.Helper()
	runDir := t.TempDir()
	logPath := filepath.Join(runDir, "stdout.log")
	if err := os.WriteFile(logPath, []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	q, err := hitl.New(runDir).Create(hitl.Question{Agent: "codex", Prompt: "Which branch?", Risk: hitl.RiskNormal})
	if err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	m := newLiveModel(LiveOptions{Idea: testIdea("hitl"), Participants: []string{"codex"}, RunID: "run-1", RunDir: runDir})
	m.height, m.width = 30, 100
	m.events = []store.Event{{Time: base, Type: "agent.started", Data: map[string]any{"agent": "codex", "stdout": logPath, "segment_id": "segment-0001"}}}
	m.state = ProjectEvents([]string{"codex"}, m.events, m.now)
	m.ensureActiveBuffer()
	u, _ := m.Update(questionsMsg{questions: []hitl.Question{q}})
	return u.(liveModel), runDir, q
}

// 1. Bare /open opens the picker; explicit /open <target> does not.
func TestPickerOpenBareVsExplicit(t *testing.T) {
	m := homeModelWithIdeas(t, "alpha", "beta")
	m, _ = runCmd(t, m, "/open")
	if !m.picker.Active || m.picker.Kind != pickerOpen {
		t.Fatalf("bare /open should open the open-picker, got active=%v kind=%q", m.picker.Active, m.picker.Kind)
	}
	if len(m.picker.filtered()) != 2 {
		t.Fatalf("picker should list 2 ideas, got %d", len(m.picker.filtered()))
	}

	m2 := homeModelWithIdeas(t, "alpha", "beta")
	m2, _ = runCmd(t, m2, "/open alpha")
	if m2.picker.Active {
		t.Fatal("explicit /open <target> must NOT open the picker")
	}
}

// 11. Empty candidates set inputErr and do not open the picker.
func TestPickerOpenEmptyCandidates(t *testing.T) {
	m := homeModelWithIdeas(t) // no ideas, no runs
	m, _ = runCmd(t, m, "/open")
	if m.picker.Active {
		t.Fatal("/open with no candidates must not open a picker")
	}
	if m.inputErr != "nothing to open yet" {
		t.Fatalf("inputErr=%q, want 'nothing to open yet'", m.inputErr)
	}

	mq := liveModelWithLog(t, "x\n") // run, but no open questions
	mq, _ = runCmd(t, mq, "/answer")
	if mq.picker.Active {
		t.Fatal("/answer with no open questions must not open a picker")
	}
	if mq.inputErr != "no open questions" {
		t.Fatalf("inputErr=%q, want 'no open questions'", mq.inputErr)
	}
}

// 3. Picker ↑/↓ moves the selection, not the active tab.
func TestPickerArrowsMoveSelectionNotTab(t *testing.T) {
	m := homeModelWithIdeas(t, "alpha", "beta")
	tabBefore := m.activeTabResolved()
	m, _ = runCmd(t, m, "/open")
	if m.picker.Index != 0 {
		t.Fatalf("picker starts at index 0, got %d", m.picker.Index)
	}
	m, _ = pressKey(t, m, tea.KeyDown)
	if m.picker.Index != 1 {
		t.Fatalf("down should move selection to 1, got %d", m.picker.Index)
	}
	if got := m.activeTabResolved(); got != tabBefore {
		t.Fatalf("picker arrows changed the active tab: %q -> %q", tabBefore, got)
	}
	m, _ = pressKey(t, m, tea.KeyUp)
	if m.picker.Index != 0 {
		t.Fatalf("up should move selection back to 0, got %d", m.picker.Index)
	}
}

// 4. Printable runes (incl. N and /) filter the picker, never touch inputText.
func TestPickerPrintableFiltersNotInput(t *testing.T) {
	m := homeModelWithIdeas(t, "alpha", "beta")
	m, _ = runCmd(t, m, "/open")
	m = pressRunes(t, m, "N")
	m = pressRunes(t, m, "/")
	if m.picker.Filter != "N/" {
		t.Fatalf("picker.Filter=%q, want 'N/'", m.picker.Filter)
	}
	if m.inputText != "" {
		t.Fatalf("inputText=%q, want empty while picker filters", m.inputText)
	}
	if m.composing {
		t.Fatal("typing N while picking must not open the new-idea composer")
	}
}

// 5. reclamp keeps Index inside the filtered list.
func TestPickerReclampClampsIndex(t *testing.T) {
	p := pickerState{Items: []pickerItem{{Label: "a"}, {Label: "b"}, {Label: "c"}}, Index: 5}
	p.reclamp(8)
	if p.Index != 2 {
		t.Fatalf("reclamp should clamp Index to 2 (len-1), got %d", p.Index)
	}
	p.Filter = "a"
	p.Index = 2
	p.reclamp(8)
	if p.Index != 0 {
		t.Fatalf("after filtering to 1 match, Index should clamp to 0, got %d", p.Index)
	}
}

// 6. Picker Enter for /open dispatches through the same openRun path as explicit /open <value>.
func TestPickerOpenEnterUsesOpenRun(t *testing.T) {
	mPick := homeModelWithIdeas(t, "zzz")
	mPick, _ = runCmd(t, mPick, "/open")
	mPick, _ = pressKey(t, mPick, tea.KeyEnter) // select "zzz"
	if mPick.picker.Active {
		t.Fatal("Enter on a candidate should close the picker")
	}
	mExpl := homeModelWithIdeas(t, "zzz")
	mExpl, _ = runCmd(t, mExpl, "/open zzz")
	if mPick.inputErr != mExpl.inputErr {
		t.Fatalf("picker-select and explicit /open differ: %q vs %q", mPick.inputErr, mExpl.inputErr)
	}
	if !strings.HasPrefix(mPick.inputErr, "open failed:") {
		t.Fatalf("expected an open-failed error from openRun, got %q", mPick.inputErr)
	}
}

// 2 + 9. Bare /answer opens the answer-picker; selecting a question enters answer
// composition; submitting answers the question and clears compose state.
func TestPickerAnswerTwoStepSubmit(t *testing.T) {
	m, runDir, q := runModelWithQuestion(t)
	m, _ = runCmd(t, m, "/answer")
	if !m.picker.Active || m.picker.Kind != pickerAnswer {
		t.Fatalf("bare /answer should open the answer-picker, got active=%v kind=%q", m.picker.Active, m.picker.Kind)
	}
	m, _ = pressKey(t, m, tea.KeyEnter) // select the one question
	if m.picker.Active {
		t.Fatal("selecting a question should close the picker")
	}
	if !m.composing || m.answerQID != q.ID || m.inputText != "" {
		t.Fatalf("after select: composing=%v answerQID=%q inputText=%q, want composing=true qid=%s empty", m.composing, m.answerQID, m.inputText, q.ID)
	}
	if !strings.Contains(m.View(), "answer "+q.ID) {
		t.Fatalf("input row should show the answer label\n%s", m.View())
	}
	m = pressRunes(t, m, "main")
	m, _ = pressKey(t, m, tea.KeyEnter)
	if m.composing || m.answerQID != "" || m.inputText != "" {
		t.Fatalf("after submit: composing=%v answerQID=%q inputText=%q, want all cleared", m.composing, m.answerQID, m.inputText)
	}
	qs, err := hitl.New(runDir).List()
	if err != nil {
		t.Fatal(err)
	}
	if qs[0].Answer != "main" || qs[0].Status != hitl.StatusAnswered {
		t.Fatalf("question=%+v, want answered=main", qs[0])
	}
}

// 7. esc cancels only the picker and preserves the attached run / active tab.
func TestPickerEscPreservesRunContext(t *testing.T) {
	m, runDir, _ := runModelWithQuestion(t)
	tabBefore := m.activeTabResolved()
	m, _ = runCmd(t, m, "/answer")
	if !m.picker.Active {
		t.Fatal("/answer should have opened the picker")
	}
	m, cmd := pressKey(t, m, tea.KeyEsc)
	if cmd != nil {
		t.Fatal("esc on an open picker must not quit")
	}
	if m.picker.Active {
		t.Fatal("esc should cancel the picker")
	}
	if !m.hasRun() || m.opts.RunDir != runDir {
		t.Fatalf("esc on the picker must not detach the run (RunDir=%q)", m.opts.RunDir)
	}
	if got := m.activeTabResolved(); got != tabBefore {
		t.Fatalf("esc on the picker changed the active tab: %q -> %q", tabBefore, got)
	}
}

// 8. Empty filtered results show an empty state; Enter selects nothing.
func TestPickerEmptyFilterShowsStateAndEnterNoops(t *testing.T) {
	m := homeModelWithIdeas(t, "alpha")
	m, _ = runCmd(t, m, "/open")
	m = pressRunes(t, m, "zzzz") // matches nothing
	if len(m.picker.filtered()) != 0 {
		t.Fatalf("filter should match nothing, got %d", len(m.picker.filtered()))
	}
	if !strings.Contains(m.View(), "no matches") {
		t.Fatalf("empty filter should show a 'no matches' state\n%s", m.View())
	}
	before := m
	m, _ = pressKey(t, m, tea.KeyEnter)
	if !m.picker.Active {
		t.Fatal("Enter on an empty filtered list should keep the picker open")
	}
	if m.composing != before.composing || m.opts.RunDir != before.opts.RunDir {
		t.Fatal("Enter on empty results must not mutate run/answer state")
	}
}

// 10. Cancelling an answer clears answerQID so a later composition can't answer it.
func TestPickerAnswerCancelClearsQID(t *testing.T) {
	m, _, q := runModelWithQuestion(t)
	m, _ = runCmd(t, m, "/answer")
	m, _ = pressKey(t, m, tea.KeyEnter) // enter answer composition
	if m.answerQID != q.ID {
		t.Fatalf("expected answerQID=%s, got %q", q.ID, m.answerQID)
	}
	m, _ = pressKey(t, m, tea.KeyEsc) // cancel composition
	if m.composing || m.answerQID != "" {
		t.Fatalf("esc should clear answer composition, got composing=%v answerQID=%q", m.composing, m.answerQID)
	}
}

// FINAL §8: a background questionsMsg while the /answer picker is open rebuilds
// the items without resetting the filter or the cursor.
func TestPickerAnswerRefreshesOnBackgroundUpdate(t *testing.T) {
	m, runDir, q1 := runModelWithQuestion(t)
	q2, err := hitl.New(runDir).Create(hitl.Question{Agent: "codex", Prompt: "Deploy now?", Risk: hitl.RiskNormal})
	if err != nil {
		t.Fatal(err)
	}
	m, _ = runCmd(t, m, "/answer")
	if len(m.picker.filtered()) != 1 {
		t.Fatalf("picker should start with 1 known question, got %d", len(m.picker.filtered()))
	}
	m = pressRunes(t, m, "Deploy") // matches q2's prompt once it arrives
	u, _ := m.Update(questionsMsg{questions: []hitl.Question{q1, q2}, token: m.runToken})
	m = u.(liveModel)
	if !m.picker.Active {
		t.Fatal("picker should stay open across a background refresh")
	}
	if m.picker.Filter != "Deploy" {
		t.Fatalf("filter must be preserved across refresh, got %q", m.picker.Filter)
	}
	got := m.picker.filtered()
	if len(got) != 1 || got[0].Value != q2.ID {
		t.Fatalf("rebuilt+filtered list should be [q2], got %+v", got)
	}
	m, _ = pressKey(t, m, tea.KeyEnter)
	if m.answerQID != q2.ID {
		t.Fatalf("Enter should select q2 from the rebuilt list, answerQID=%q", m.answerQID)
	}
}

// agy #1: a failed answer write keeps the user in answer composition to retry,
// rather than stranding the qid.
func TestPickerAnswerFailureKeepsCompose(t *testing.T) {
	m, _, _ := runModelWithQuestion(t)
	m, _ = runCmd(t, m, "/answer")
	m, _ = pressKey(t, m, tea.KeyEnter) // enter answer composition
	m.answerQID = "nope-does-not-exist" // force the write to fail
	m = pressRunes(t, m, "x")
	m, _ = pressKey(t, m, tea.KeyEnter)
	if !m.composing || m.answerQID != "nope-does-not-exist" {
		t.Fatalf("a failed answer must keep compose state: composing=%v answerQID=%q", m.composing, m.answerQID)
	}
	if m.inputErr == "" {
		t.Fatal("a failed answer should surface an error")
	}
}

// --- tui-live-steering tests ---

// Tab completes a single-match slash command (no-arg → exact name).
func TestSuggestTabCompletesSingleMatch(t *testing.T) {
	m := homeModelWithIdeas(t)
	m = pressRunes(t, m, "/q")
	if !m.suggest || len(m.suggestItems) != 1 || m.suggestItems[0].Name != "/quit" {
		t.Fatalf("expected suggest=[/quit], got suggest=%v items=%v", m.suggest, m.suggestItems)
	}
	m, _ = pressKey(t, m, tea.KeyTab)
	if m.inputText != "/quit" {
		t.Fatalf("Tab should complete to /quit, got %q", m.inputText)
	}
}

// Tab completes the longest common prefix when several commands match.
func TestSuggestTabCompletesCommonPrefix(t *testing.T) {
	m := homeModelWithIdeas(t)
	m = pressRunes(t, m, "/h") // matches /help and /home → LCP "/h"
	if len(m.suggestItems) != 2 {
		t.Fatalf("expected 2 matches for /h, got %v", m.suggestItems)
	}
	m, _ = pressKey(t, m, tea.KeyTab)
	if m.inputText != "/h" {
		t.Fatalf("Tab should leave the common prefix /h, got %q", m.inputText)
	}
	m = pressRunes(t, m, "o") // → /home only
	m, _ = pressKey(t, m, tea.KeyTab)
	if m.inputText != "/home" {
		t.Fatalf("Tab should complete /home, got %q", m.inputText)
	}
}

// The suggest menu clears once the input is no longer a bare slash command.
func TestSuggestClearsOnSpace(t *testing.T) {
	m := homeModelWithIdeas(t)
	m = pressRunes(t, m, "/open")
	if !m.suggest {
		t.Fatal("/open should show a suggestion")
	}
	m = pressRunes(t, m, " ")
	if m.suggest {
		t.Fatal("a space after the command must close the suggest menu")
	}
}

// Conditional Tab: with non-slash input, Tab switches tabs (legacy behavior).
func TestConditionalTabSwitchesTabsWhenNotSlash(t *testing.T) {
	m := liveModelWithLog(t, "x\n") // run attached → Home + agent:codex + Status tabs
	before := m.activeTabResolved()
	m, _ = pressKey(t, m, tea.KeyTab)
	if m.activeTabResolved() == before {
		t.Fatalf("Tab on empty input should switch tabs (was %q)", before)
	}
}

// ctrl+k on a running agent tab opens a modal confirm; y kills via the seam.
func TestConfirmKillModal(t *testing.T) {
	m := liveModelWithLog(t, "x\n") // codex running, active tab = agent:codex
	var killed string
	m.opts.KillAgent = func(agentID string) error { killed = agentID; return nil }

	m, _ = pressKey(t, m, tea.KeyCtrlK)
	if m.confirmKillAgent != "codex" {
		t.Fatalf("ctrl+k should open confirm for codex, got %q", m.confirmKillAgent)
	}
	// n cancels without calling the seam.
	m = pressRunes(t, m, "n")
	if m.confirmKillAgent != "" || killed != "" {
		t.Fatalf("n should cancel the confirm (confirm=%q killed=%q)", m.confirmKillAgent, killed)
	}
	// ctrl+k again, then y confirms.
	m, _ = pressKey(t, m, tea.KeyCtrlK)
	m = pressRunes(t, m, "y")
	if m.confirmKillAgent != "" {
		t.Fatal("y should close the confirm")
	}
	if killed != "codex" {
		t.Fatalf("y should kill codex via the seam, got %q", killed)
	}
}

// A steer on an agent tab routes through the SubmitSteer seam and registers the
// reply for inline display.
func TestSubmitInputSteerViaSeam(t *testing.T) {
	m := liveModelWithLog(t, "x\n") // active tab = agent:codex, no open question
	replyLog := filepath.Join(t.TempDir(), "steer-stdout.log")
	if err := os.WriteFile(replyLog, []byte("codex: here is my reply\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var got SteerRequest
	m.opts.SubmitSteer = func(req SteerRequest) (SteerResult, error) {
		got = req
		return SteerResult{ID: "steer-9", Status: "running", StdoutPath: replyLog}, nil
	}
	m = pressRunes(t, m, "hello there")
	m, _ = pressKey(t, m, tea.KeyEnter)

	if got.Target != steer.TargetAgent || got.AgentID != "codex" || got.Text != "hello there" {
		t.Fatalf("steer request not routed correctly: %+v", got)
	}
	sr := m.steerReplies["codex"]
	if sr == nil || sr.stdoutPath != replyLog || sr.query != "hello there" {
		t.Fatalf("steer reply not registered for inline display: %+v", sr)
	}
}

// A steer.replied / steer.reply_failed event flips the inline reply's state so
// the "replying…" indicator clears.
func TestSteerReplyEventTransitionsState(t *testing.T) {
	m := liveModelWithLog(t, "x\n")
	m.steerReplies = map[string]*steerReply{"codex": {id: "s1", query: "hi"}}
	updated, _ := m.Update(eventsMsg{
		token:  m.runToken,
		events: []store.Event{{Type: "steer.replied", Data: map[string]any{"agent": "codex", "id": "s1"}}},
	})
	m = updated.(liveModel)
	if sr := m.steerReplies["codex"]; sr == nil || !sr.done || sr.failed {
		t.Fatalf("steer.replied should mark done (not failed): %+v", m.steerReplies["codex"])
	}

	m2 := liveModelWithLog(t, "x\n")
	m2.steerReplies = map[string]*steerReply{"codex": {id: "s2", query: "hi"}}
	updated, _ = m2.Update(eventsMsg{
		token:  m2.runToken,
		events: []store.Event{{Type: "steer.reply_failed", Data: map[string]any{"agent": "codex", "id": "s2"}}},
	})
	m2 = updated.(liveModel)
	if sr := m2.steerReplies["codex"]; sr == nil || !sr.done || !sr.failed {
		t.Fatalf("steer.reply_failed should mark done+failed: %+v", m2.steerReplies["codex"])
	}
}

// A killed agent shows the KILL badge, not ERR.
func TestKilledAgentShortState(t *testing.T) {
	if got := shortState(stateKilled); got != "KILL" {
		t.Fatalf("killed badge = %q, want KILL", got)
	}
}

// activateRun must copy the SubmitSteer/KillAgent seams onto the active run, or
// runs launched from Home silently lose steering/kill.
func TestActivateRunCopiesSteerKillSeams(t *testing.T) {
	m := newLiveModel(LiveOptions{Home: true, Root: t.TempDir()})
	m.height, m.width = 30, 100
	m.activateRun(LaunchResult{
		Idea:         testIdea("x"),
		Participants: []string{"codex"},
		RunID:        "r",
		RunDir:       t.TempDir(),
		SubmitSteer:  func(SteerRequest) (SteerResult, error) { return SteerResult{}, nil },
		KillAgent:    func(string) error { return nil },
	})
	if m.opts.SubmitSteer == nil || m.opts.KillAgent == nil {
		t.Fatal("activateRun must copy SubmitSteer and KillAgent onto opts")
	}
}
