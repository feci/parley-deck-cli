package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

// newSteerHandle builds a Handle wired for steer/kill tests with one fake agent.
func newSteerHandle(t *testing.T, agentID string, headlessArgs []string) (*Handle, string) {
	t.Helper()
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Steer task", []string{agentID})
	if err != nil {
		t.Fatal(err)
	}
	runDir := filepath.Join(root, protocol.DeckDir, "runs", "steer-run")
	opts := Options{
		Root:    root,
		RunID:   "steer-run",
		Idea:    idea,
		Agents:  []agents.Discovery{{Spec: agents.Spec{ID: agentID, HeadlessArgs: headlessArgs, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true}},
		Timeout: 5 * time.Second,
		Store:   store.New(runDir),
	}
	h := &Handle{
		RunID:     "steer-run",
		RunDir:    runDir,
		rootCtx:   context.Background(),
		active:    map[string]*attempt{},
		steerBusy: map[string]bool{},
		opts:      opts,
	}
	return h, runDir
}

func TestRunSteerAttemptCapturesReply(t *testing.T) {
	h, runDir := newSteerHandle(t, "agy", []string{"-test.run=TestFakeSteerReplyHelper", "--", "parley-fake-steer"})
	res, err := h.RunSteerAttempt(context.Background(), SteerAttemptRequest{AgentID: "agy", Text: "what is the plan?", SteerID: "steer-0001"})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Accepted || res.Status != "running" {
		t.Fatalf("expected accepted+running, got %+v", res)
	}
	// Wait for the async attempt to finish (steerBusy clears on completion).
	waitFor(t, 5*time.Second, func() bool { return !h.steerBusyFor("agy") })

	reply := readSteerReply(res.ReplyPath, res.StdoutPath)
	if !strings.Contains(reply, "STEER-REPLY-OK") {
		t.Fatalf("reply not captured: %q", reply)
	}
	events, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if !hasEvent(events, "steer.reply_started", "steer-0001") || !hasEvent(events, "steer.replied", "steer-0001") {
		t.Fatalf("expected steer.reply_started + steer.replied for steer-0001, got %v", eventTypes(events))
	}
	// The steer must write under its own per-steer dir, never the round stdout.log.
	if !strings.Contains(res.StdoutPath, filepath.Join("agents", "agy", "steers", "steer-0001")) {
		t.Fatalf("steer stdout not isolated: %s", res.StdoutPath)
	}
}

func TestRunSteerAttemptRejectsSecond(t *testing.T) {
	h, _ := newSteerHandle(t, "agy", []string{"-test.run=TestFakeSteerReplyHelper", "--", "parley-fake-steer"})
	h.mu.Lock()
	h.steerBusy["agy"] = true // simulate an in-flight steer
	h.mu.Unlock()
	res, err := h.RunSteerAttempt(context.Background(), SteerAttemptRequest{AgentID: "agy", Text: "again", SteerID: "steer-0002"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Accepted || res.Status != "rejected" || !strings.Contains(res.Message, "already replying") {
		t.Fatalf("second steer must be rejected, got %+v", res)
	}
}

func TestRunSteerAttemptRejectsNonParticipant(t *testing.T) {
	h, _ := newSteerHandle(t, "agy", []string{"-test.run=TestFakeSteerReplyHelper", "--", "parley-fake-steer"})
	// A discovered agent that is NOT a participant in the idea.
	h.opts.Agents = append(h.opts.Agents, agents.Discovery{Spec: agents.Spec{ID: "ghost"}, Path: os.Args[0], Found: true})
	res, err := h.RunSteerAttempt(context.Background(), SteerAttemptRequest{AgentID: "ghost", Text: "hi", SteerID: "steer-x"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Accepted || !strings.Contains(res.Message, "not a participant") {
		t.Fatalf("non-participant steer must be rejected, got %+v", res)
	}
}

func TestRunSteerAttemptQueuedThenRuns(t *testing.T) {
	h, runDir := newSteerHandle(t, "agy", []string{"-test.run=TestFakeSteerReplyHelper", "--", "parley-fake-steer"})
	// Simulate an in-flight round attempt for agy.
	_, cancel := context.WithCancel(context.Background())
	h.register("agy", "segment-0001", "round", "", cancel)

	res, err := h.RunSteerAttempt(context.Background(), SteerAttemptRequest{AgentID: "agy", Text: "later", SteerID: "steer-q"})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Accepted || res.Status != "queued" {
		t.Fatalf("steer behind a running attempt should be queued, got %+v", res)
	}
	// Release the round attempt → the queued steer proceeds and finishes.
	h.finish("agy")
	cancel()
	waitFor(t, 5*time.Second, func() bool { return !h.steerBusyFor("agy") })
	events, _ := store.New(runDir).Load()
	if !hasEvent(events, "steer.replied", "steer-q") {
		t.Fatalf("queued steer should run and reply, got %v", eventTypes(events))
	}
}

// A queued steer must NOT emit run.segment_started — that would reset/reorder the
// still-running round agent's projected state (regression guard for round-02).
func TestQueuedSteerEmitsNoSegmentBoundary(t *testing.T) {
	h, runDir := newSteerHandle(t, "agy", []string{"-test.run=TestFakeSteerReplyHelper", "--", "parley-fake-steer"})
	rootCtx, cancelRoot := context.WithCancel(context.Background())
	h.rootCtx = rootCtx
	// agy is mid-round.
	_ = store.New(runDir).Append(store.Event{Type: "agent.started", Data: map[string]any{"agent": "agy", "segment_id": "segment-0001"}})
	_, cancelRound := context.WithCancel(context.Background())
	h.register("agy", "segment-0001", "round", "", cancelRound)

	res, err := h.RunSteerAttempt(rootCtx, SteerAttemptRequest{AgentID: "agy", Text: "q", SteerID: "steer-q2"})
	if err != nil || res.Status != "queued" {
		t.Fatalf("expected queued, got %+v err=%v", res, err)
	}
	if res.SegmentID != "steer/steer-q2" {
		t.Fatalf("steer segment label should be steer-scoped, got %q", res.SegmentID)
	}
	events, _ := store.New(runDir).Load()
	for _, e := range events {
		if e.Type == "run.segment_started" {
			t.Fatalf("a steer must not emit run.segment_started (would reset the running agent): %+v", e.Data)
		}
	}
	// Drain the queued goroutine cleanly.
	cancelRoot()
	cancelRound()
	waitFor(t, 5*time.Second, func() bool { return !h.steerBusyFor("agy") })
}

func TestRunSteerAttemptClearsBusyOnFailure(t *testing.T) {
	h, runDir := newSteerHandle(t, "agy", []string{"-test.run=TestFakeSteerSilentHelper", "--", "parley-fake-steer-silent"})
	res, err := h.RunSteerAttempt(context.Background(), SteerAttemptRequest{AgentID: "agy", Text: "hi", SteerID: "steer-f"})
	if err != nil || !res.Accepted {
		t.Fatalf("expected accepted, got %+v err=%v", res, err)
	}
	waitFor(t, 5*time.Second, func() bool { return !h.steerBusyFor("agy") })
	events, _ := store.New(runDir).Load()
	if !hasEvent(events, "steer.reply_failed", "steer-f") {
		t.Fatalf("a silent agent should produce steer.reply_failed, got %v", eventTypes(events))
	}
	if h.steerBusyFor("agy") {
		t.Fatal("steerBusy must be cleared after a failed steer")
	}
}

func TestKillAgentIdempotent(t *testing.T) {
	h, runDir := newSteerHandle(t, "agy", nil)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	h.register("agy", "segment-0001", "round", "", cancel)
	if r := h.KillAgent("agy"); !r.Killed {
		t.Fatalf("first kill should succeed, got %+v", r)
	}
	if r := h.KillAgent("agy"); r.Killed {
		t.Fatalf("second kill must be a no-op while the first is in flight, got %+v", r)
	}
	events, _ := store.New(runDir).Load()
	n := 0
	for _, e := range events {
		if e.Type == "agent.killed" {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("exactly one agent.killed event expected, got %d", n)
	}
}

func TestKillAgentTargetsOnlyOneAgent(t *testing.T) {
	h, runDir := newSteerHandle(t, "agy", nil)
	ctxA, cancelA := context.WithCancel(context.Background())
	ctxB, cancelB := context.WithCancel(context.Background())
	defer cancelA()
	defer cancelB()
	h.register("agy", "segment-0001", "round", "", cancelA)
	h.register("hermes", "segment-0001", "round", "", cancelB)

	res := h.KillAgent("agy")
	if !res.Killed {
		t.Fatalf("KillAgent should report killed, got %+v", res)
	}
	if ctxA.Err() == nil {
		t.Fatal("killed agent's context must be canceled")
	}
	if ctxB.Err() != nil {
		t.Fatal("the other agent's context must NOT be canceled (run continues)")
	}
	if !h.finish("agy") {
		t.Fatal("finish must report the agent was killed")
	}
	events, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if !hasEventAgent(events, "agent.killed", "agy") {
		t.Fatalf("expected agent.killed for agy, got %v", eventTypes(events))
	}
}

func TestKillAgentRaceWithCompletion(t *testing.T) {
	h, runDir := newSteerHandle(t, "agy", nil)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	h.register("agy", "segment-0001", "round", "", cancel)
	// Normal completion wins the race and deregisters first.
	h.finish("agy")
	res := h.KillAgent("agy")
	if res.Killed {
		t.Fatal("KillAgent must not claim a kill after normal completion")
	}
	events, _ := store.New(runDir).Load()
	if hasEventAgent(events, "agent.killed", "agy") {
		t.Fatal("no agent.killed event should be emitted when the kill loses the race")
	}
}

func (h *Handle) steerBusyFor(agentID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.steerBusy[agentID]
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}

func hasEvent(events []store.Event, typ, id string) bool {
	for _, e := range events {
		if e.Type == typ {
			if got, _ := e.Data["id"].(string); got == id {
				return true
			}
		}
	}
	return false
}

func hasEventAgent(events []store.Event, typ, agent string) bool {
	for _, e := range events {
		if e.Type == typ {
			if got, _ := e.Data["agent"].(string); got == agent {
				return true
			}
		}
	}
	return false
}

// TestFakeSteerReplyHelper is the fake agent: it prints a recognizable reply to
// stdout (the steer round-trip reads stdout as the reply).
func TestFakeSteerReplyHelper(t *testing.T) {
	if !hasArg("parley-fake-steer") {
		return
	}
	_, _ = os.Stdin.Read(make([]byte, 1)) // tolerate a stdin prompt
	os.Stdout.WriteString("STEER-REPLY-OK: here is my answer.\n")
	os.Exit(0)
}

// TestFakeSteerSilentHelper is a fake agent that produces no reply (→ failure).
func TestFakeSteerSilentHelper(t *testing.T) {
	if !hasArg("parley-fake-steer-silent") {
		return
	}
	_, _ = os.Stdin.Read(make([]byte, 1))
	os.Exit(0) // no stdout, no reply.md → steer.reply_failed
}
