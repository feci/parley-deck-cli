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
	if got := agents["hermes"].Duration; got != 0 {
		t.Fatalf("hermes duration=%s, want zero for skip without start", got)
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

func TestLiveCompactLayoutFitsShortTerminal(t *testing.T) {
	model := newLiveModel(LiveOptions{
		Idea:         testIdea("compact-live-run-tui"),
		Participants: []string{"codex"},
		RunID:        "run-1",
		RunDir:       t.TempDir(),
	})
	model.width = 100
	model.height = 18
	model.events = []store.Event{
		{Time: time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC), Type: "agent.started", Data: map[string]any{"agent": "codex"}},
	}
	model.state = ProjectEvents([]string{"codex"}, model.events, model.now)

	view := model.View()
	for _, want := range []string{"layout=compact", "Agents", "Latest events", "Questions", "Log preview", "q/esc detach"} {
		if !strings.Contains(view, want) {
			t.Fatalf("compact live view missing %q\n%s", want, view)
		}
	}
	if got := renderedLineCount(view); got > model.height {
		t.Fatalf("compact live view rendered %d lines, want <= %d\n%s", got, model.height, view)
	}
}

func TestLiveCompactThresholdBoundary(t *testing.T) {
	model := newLiveModel(LiveOptions{
		Idea:         testIdea("compact-live-boundary"),
		Participants: []string{"codex"},
		RunID:        "run-1",
		RunDir:       t.TempDir(),
	})
	model.width = 100
	model.height = compactLiveHeight
	if view := model.View(); strings.Contains(view, "layout=compact") {
		t.Fatalf("live view used compact layout at threshold height %d\n%s", compactLiveHeight, view)
	}

	model.height = compactLiveHeight - 1
	if view := model.View(); !strings.Contains(view, "layout=compact") {
		t.Fatalf("live view did not use compact layout below threshold height %d\n%s", compactLiveHeight-1, view)
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

func TestResumeViewHasExplicitExitPath(t *testing.T) {
	base := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	model := newLiveModel(LiveOptions{
		Idea:         testIdea("runtime-status-resume"),
		Participants: []string{"codex"},
		RunID:        "run-1",
		RunDir:       t.TempDir(),
		Resume:       true,
	})
	model.events = []store.Event{
		{Time: base, Type: "agent.started", Data: map[string]any{"agent": "codex"}},
	}
	model.state = ProjectEvents([]string{"codex"}, model.events, base.Add(time.Minute))

	view := model.View()
	if strings.Contains(view, "status=running") {
		t.Fatalf("resume view must not imply live process\n%s", view)
	}
	for _, want := range []string{"status=unverified", "q/esc/ctrl+c close resume view"} {
		if !strings.Contains(view, want) {
			t.Fatalf("resume view missing %q\n%s", want, view)
		}
	}
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("q did not return a quit command")
	}
}

func TestAnswerModeBackspaceRemovesWholeRune(t *testing.T) {
	model := newLiveModel(LiveOptions{RunDir: t.TempDir()})
	model.mode = modeAnswerQuestion
	model.answerText = "á"

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = updated.(liveModel)
	if model.answerText != "" {
		t.Fatalf("answerText=%q, want empty", model.answerText)
	}
}

func TestSummarizeHITLEvents(t *testing.T) {
	question := summarizeEvent(store.Event{
		Type: "hitl.question",
		Data: map[string]any{"agent": "agy", "question_id": "q1", "risk": "normal"},
	})
	if question.Text != "agy question q1 normal" {
		t.Fatalf("question text=%q", question.Text)
	}

	answered := summarizeEvent(store.Event{
		Type: "hitl.answered",
		Data: map[string]any{"agent": "agy", "question_id": "q1", "status": "answered"},
	})
	if answered.Text != "agy answered q1 answered" {
		t.Fatalf("answered text=%q", answered.Text)
	}
}

func focusModelWithLog(t *testing.T, content string) liveModel {
	t.Helper()
	dir := t.TempDir()
	logPath := filepath.Join(dir, "stdout.log")
	if err := os.WriteFile(logPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	model := newLiveModel(LiveOptions{
		Idea:         testIdea("tui-interactivity-overhaul"),
		Participants: []string{"codex"},
		RunID:        "run-1",
		RunDir:       dir,
	})
	model.height = 24
	model.width = 100
	model.events = []store.Event{
		{Time: base, Type: "agent.started", Data: map[string]any{"agent": "codex", "stdout": logPath, "segment_id": "segment-0001"}},
	}
	model.state = ProjectEvents([]string{"codex"}, model.events, model.now)
	return model
}

func TestFocusViewShowsAgentLogAndExits(t *testing.T) {
	model := focusModelWithLog(t, "hello from codex\nworking on round-01\nDONE\n")

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)
	if model.mode != modeAgentDetail {
		t.Fatal("enter did not open the focus view")
	}
	view := model.View()
	for _, want := range []string{"agent=codex", "segment=segment-0001", "working on round-01", "f follow(on)", "esc back"} {
		if !strings.Contains(view, want) {
			t.Fatalf("focus view missing %q\n%s", want, view)
		}
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(liveModel)
	if model.mode == modeAgentDetail {
		t.Fatal("esc did not exit the focus view")
	}
	if !strings.Contains(model.View(), "Log preview") {
		t.Fatal("after esc the overview should render again")
	}
}

func TestFocusFollowAndScroll(t *testing.T) {
	var sb strings.Builder
	for i := 1; i <= 100; i++ {
		fmt.Fprintf(&sb, "line %d\n", i)
	}
	model := focusModelWithLog(t, sb.String())

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)
	if !model.follow {
		t.Fatal("follow must default on in focus")
	}
	if !strings.Contains(model.View(), "line 100") {
		t.Fatalf("follow view should show the last line\n%s", model.View())
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(liveModel)
	if model.follow {
		t.Fatal("scrolling up (k) must disable follow")
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	model = updated.(liveModel)
	if model.follow {
		t.Fatal("g (top) must disable follow")
	}
	topView := model.View()
	if !strings.Contains(topView, "line 1\n") || strings.Contains(topView, "line 100") {
		t.Fatalf("g should show the top, not the bottom\n%s", topView)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	model = updated.(liveModel)
	if !model.follow {
		t.Fatal("G (bottom) must re-enable follow")
	}
	if !strings.Contains(model.View(), "line 100") {
		t.Fatal("G should show the last line again")
	}
}

func TestLoadFocusTailBoundsLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stdout.log")
	var sb strings.Builder
	total := maxFocusLines + 50
	for i := 0; i < total; i++ {
		fmt.Fprintf(&sb, "L%d\n", i)
	}
	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	lines, _, truncated := loadFocusTail(path)
	if len(lines) != maxFocusLines {
		t.Fatalf("len(lines)=%d, want %d (bounded scrollback)", len(lines), maxFocusLines)
	}
	if !truncated {
		t.Fatal("expected truncated=true when capping lines")
	}
	if want := fmt.Sprintf("L%d", total-1); lines[len(lines)-1] != want {
		t.Fatalf("last line=%q, want %q", lines[len(lines)-1], want)
	}
}

func TestReadAppendedLinesIncremental(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stdout.log")
	if err := os.WriteFile(path, []byte("a\nb\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, offset, _ := loadFocusTail(path)
	if offset != 4 {
		t.Fatalf("initial offset=%d, want 4", offset)
	}

	appendString(t, path, "c\npartial")
	lines, newOffset := readAppendedLines(path, offset)
	if len(lines) != 1 || lines[0] != "c" {
		t.Fatalf("lines=%v, want [c] (partial trailing line excluded)", lines)
	}
	if newOffset != offset+2 {
		t.Fatalf("offset=%d, want %d (advanced only past complete lines)", newOffset, offset+2)
	}

	appendString(t, path, " line\n")
	lines, _ = readAppendedLines(path, newOffset)
	if len(lines) != 1 || lines[0] != "partial line" {
		t.Fatalf("lines=%v, want [partial line] (completed partial)", lines)
	}
}

// AF4: a partial final line in the seeded tail is not fragmented — it is
// re-read and merged once it completes.
func TestFocusPartialLineNotFragmented(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stdout.log")
	if err := os.WriteFile(path, []byte("done line\nprogress"), 0o644); err != nil {
		t.Fatal(err)
	}
	lines, offset, _ := loadFocusTail(path)
	if len(lines) != 1 || lines[0] != "done line" {
		t.Fatalf("loadFocusTail lines=%v, want [done line] (trailing partial excluded)", lines)
	}
	appendString(t, path, " complete\n")
	newLines, _ := readAppendedLines(path, offset)
	if len(newLines) != 1 || newLines[0] != "progress complete" {
		t.Fatalf("readAppendedLines=%v, want [progress complete] (partial merged, not fragmented)", newLines)
	}
}

// AF2: an oversized newline-less prefix is dropped so the retained buffer stays
// within the byte cap.
func TestFocusTailDropsOversizedLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stdout.log")
	big := strings.Repeat("x", maxFocusBytes+4096)
	if err := os.WriteFile(path, []byte(big+"\nshort tail line\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lines, _, truncated := loadFocusTail(path)
	total := 0
	for _, l := range lines {
		total += len(l) + 1
	}
	if total > maxFocusBytes {
		t.Fatalf("retained %d bytes, want <= %d (byte cap)", total, maxFocusBytes)
	}
	if !truncated {
		t.Fatal("expected truncated=true when dropping the oversized line")
	}
	if len(lines) == 0 || lines[len(lines)-1] != "short tail line" {
		t.Fatalf("lines=%v, want the short line retained after dropping the oversized one", lines)
	}
}

// AF2: capFocusLines evicts oldest lines until the buffer is within both the
// line and byte budgets.
func TestCapFocusLinesByteBudget(t *testing.T) {
	line := strings.Repeat("y", 1000)
	n := (maxFocusBytes / 1001) + 100 // exceeds the byte budget, under the line cap
	lines := make([]string, n)
	for i := range lines {
		lines[i] = line
	}
	capped, truncated := capFocusLines(lines)
	if !truncated {
		t.Fatal("expected truncated=true")
	}
	total := 0
	for _, l := range capped {
		total += len(l) + 1
	}
	if total > maxFocusBytes {
		t.Fatalf("capped total=%d bytes, want <= %d", total, maxFocusBytes)
	}
}

// fix-up cycle 2: a single retained line larger than the byte cap is
// head-truncated so the buffer still honors the budget (codex round-02 finding).
func TestCapFocusLinesTruncatesSingleOversizedLine(t *testing.T) {
	capped, truncated := capFocusLines([]string{strings.Repeat("z", maxFocusBytes+500)})
	if !truncated {
		t.Fatal("expected truncated=true for a single oversized line")
	}
	total := 0
	for _, l := range capped {
		total += len(l) + 1
	}
	if total > maxFocusBytes {
		t.Fatalf("single oversized line not capped: total=%d > %d", total, maxFocusBytes)
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

func TestComposerQueuesAgentSteer(t *testing.T) {
	model := focusModelWithLog(t, "x\n")

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updated.(liveModel)
	if model.mode != modeCompose {
		t.Fatalf("i should open the composer, mode=%d", model.mode)
	}
	if !strings.Contains(model.View(), "Steer agent codex") {
		t.Fatalf("composer view missing target\n%s", model.View())
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("focus on the parser")})
	model = updated.(liveModel)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(liveModel)
	if model.mode != modeOverview {
		t.Fatalf("enter should submit and return to overview, mode=%d", model.mode)
	}
	if !strings.Contains(model.statusMsg, "recorded steer-0001-") || !strings.Contains(model.statusMsg, "codex") || !strings.Contains(model.statusMsg, "queued") {
		t.Fatalf("statusMsg=%q, want an honest recorded/queued confirmation", model.statusMsg)
	}

	queued, err := steer.List(model.opts.RunDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(queued) != 1 || queued[0].Agent != "codex" || queued[0].Text != "focus on the parser" {
		t.Fatalf("persisted steer = %+v", queued)
	}
	if queued[0].SegmentID != "segment-0001" {
		t.Fatalf("steer should capture the agent's segment, got %q", queued[0].SegmentID)
	}
}

func TestComposerEscCancels(t *testing.T) {
	model := focusModelWithLog(t, "x\n")
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updated.(liveModel)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("oops")})
	model = updated.(liveModel)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(liveModel)
	if model.mode != modeOverview {
		t.Fatal("esc should cancel the composer")
	}
	queued, _ := steer.List(model.opts.RunDir)
	if len(queued) != 0 {
		t.Fatalf("cancel must not persist a steer, got %d", len(queued))
	}
}

func TestHelpOverlayToggles(t *testing.T) {
	model := focusModelWithLog(t, "x\n")

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model = updated.(liveModel)
	if model.mode != modeHelp {
		t.Fatalf("? should open the help overlay, mode=%d", model.mode)
	}
	view := model.View()
	for _, want := range []string{"Help", "open agent focus view", "toggle follow"} {
		if !strings.Contains(view, want) {
			t.Fatalf("help overlay missing %q\n%s", want, view)
		}
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(liveModel)
	if model.mode != modeOverview {
		t.Fatal("esc should close the help overlay")
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
