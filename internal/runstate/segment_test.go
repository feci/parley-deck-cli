package runstate

import (
	"testing"
	"time"

	"parley-deck-cli/internal/store"
)

func agentByID(state RunState, id string) (AgentState, bool) {
	for _, a := range state.Agents {
		if a.ID == id {
			return a, true
		}
	}
	return AgentState{}, false
}

func mustAgent(t *testing.T, state RunState, id string) AgentState {
	t.Helper()
	a, ok := agentByID(state, id)
	if !ok {
		t.Fatalf("agent %q not present in projection", id)
	}
	return a
}

// segStarted builds a run.segment_started event for the given segment + targets.
func segStarted(at time.Time, segment, reason string, targets ...string) store.Event {
	return store.Event{Time: at, Type: "run.segment_started", Data: map[string]any{
		"segment_id": segment, "reason": reason, "targets": targets,
	}}
}

func agentEvent(at time.Time, typ, id, segment string) store.Event {
	return store.Event{Time: at, Type: typ, Data: map[string]any{"agent": id, "segment_id": segment}}
}

// A killed agent projects to StateKilled, and a sticky kill survives the trailing
// agent.failed that the canceled process emits (tui-live-steering).
func TestProjectEventsKilledIsSticky(t *testing.T) {
	base := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		agentEvent(base, "agent.started", "agy", "segment-0001"),
		agentEvent(base.Add(time.Second), "agent.killed", "agy", "segment-0001"),
		agentEvent(base.Add(2*time.Second), "agent.failed", "agy", "segment-0001"),
	}
	state := ProjectEvents([]string{"agy"}, events, base.Add(3*time.Second))
	a := mustAgent(t, state, "agy")
	if a.State != StateKilled {
		t.Fatalf("killed agent should project to %q, got %q", StateKilled, a.State)
	}
	if !a.Killed {
		t.Fatal("the Killed flag should be set")
	}
	// A later segment resets the killed state (re-run).
	events = append(events, segStarted(base.Add(4*time.Second), "segment-0002", "continue", "agy"))
	state = ProjectEvents([]string{"agy"}, events, base.Add(5*time.Second))
	if a := mustAgent(t, state, "agy"); a.Killed || a.State != StatePending {
		t.Fatalf("a new segment should clear killed (got state=%q killed=%v)", a.State, a.Killed)
	}
}

// Old runs emit no segment events; the projection must behave exactly as before.
func TestProjectEventsLegacyUnsegmentedUnchanged(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		{Time: base, Type: "agent.started", Data: map[string]any{"agent": "a"}},
		{Time: base.Add(time.Second), Type: "agent.finished", Data: map[string]any{"agent": "a"}},
	}
	state := ProjectEvents([]string{"a"}, events, base.Add(2*time.Second))
	if got := mustAgent(t, state, "a"); got.State != StateFinished {
		t.Fatalf("legacy run: a.State = %q, want finished", got.State)
	}
}

// The core bug: a finished agent re-targeted by a new segment must NOT keep its
// stale terminal badge; it resets to pending until the new segment produces an
// event.
func TestSegmentResetUnsticksFinishedBadge(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		segStarted(base, "segment-0001", "initial", "a"),
		agentEvent(base.Add(time.Second), "agent.started", "a", "segment-0001"),
		agentEvent(base.Add(2*time.Second), "agent.finished", "a", "segment-0001"),
		segStarted(base.Add(3*time.Second), "segment-0002", "continue", "a"),
	}
	state := ProjectEvents([]string{"a"}, events, base.Add(4*time.Second))
	got := mustAgent(t, state, "a")
	if got.State != StatePending {
		t.Fatalf("after new segment: a.State = %q, want pending (badge must unstick)", got.State)
	}
	if got.Segment != "segment-0002" {
		t.Fatalf("a.Segment = %q, want segment-0002", got.Segment)
	}
	if !got.StartedAt.IsZero() || got.Duration != 0 {
		t.Fatalf("reset must clear StartedAt/Duration, got StartedAt=%v Duration=%v", got.StartedAt, got.Duration)
	}
}

// After the reset, the new segment's agent.started transitions back to running.
func TestSegmentContinueThenRunning(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		segStarted(base, "segment-0001", "initial", "a"),
		agentEvent(base.Add(time.Second), "agent.finished", "a", "segment-0001"),
		segStarted(base.Add(2*time.Second), "segment-0002", "continue", "a"),
		agentEvent(base.Add(3*time.Second), "agent.started", "a", "segment-0002"),
	}
	state := ProjectEvents([]string{"a"}, events, base.Add(4*time.Second))
	if got := mustAgent(t, state, "a"); got.State != StateRunning {
		t.Fatalf("a.State = %q, want running", got.State)
	}
}

// A new segment resets ONLY its targets; other agents keep their state.
func TestNonTargetedAgentKeepsState(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		segStarted(base, "segment-0001", "initial", "a", "b"),
		agentEvent(base.Add(time.Second), "agent.finished", "a", "segment-0001"),
		agentEvent(base.Add(time.Second), "agent.finished", "b", "segment-0001"),
		segStarted(base.Add(2*time.Second), "segment-0002", "continue", "a"),
	}
	state := ProjectEvents([]string{"a", "b"}, events, base.Add(3*time.Second))
	if got := mustAgent(t, state, "a"); got.State != StatePending {
		t.Fatalf("targeted a.State = %q, want pending", got.State)
	}
	if got := mustAgent(t, state, "b"); got.State != StateFinished {
		t.Fatalf("non-targeted b.State = %q, want finished (must be untouched)", got.State)
	}
}

// A targeted-but-skipped agent renders skipped for the segment, never inheriting
// a prior finished badge.
func TestSkipInSegmentRendersSkipped(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		segStarted(base, "segment-0001", "initial", "a"),
		agentEvent(base.Add(time.Second), "agent.finished", "a", "segment-0001"),
		segStarted(base.Add(2*time.Second), "segment-0002", "continue", "a"),
		{Time: base.Add(3 * time.Second), Type: "agent.skipped", Data: map[string]any{"agent": "a", "reason": "artifact already exists", "segment_id": "segment-0002"}},
	}
	state := ProjectEvents([]string{"a"}, events, base.Add(4*time.Second))
	if got := mustAgent(t, state, "a"); got.State != StateSkipped {
		t.Fatalf("a.State = %q, want skipped", got.State)
	}
}

// A failed agent retried in a new segment recovers to finished and clears the
// prior error.
func TestRetryAfterFailed(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		segStarted(base, "segment-0001", "initial", "a"),
		{Time: base.Add(time.Second), Type: "agent.failed", Data: map[string]any{"agent": "a", "error": "boom", "segment_id": "segment-0001"}},
		segStarted(base.Add(2*time.Second), "segment-0002", "retry", "a"),
		agentEvent(base.Add(3*time.Second), "agent.started", "a", "segment-0002"),
		agentEvent(base.Add(4*time.Second), "agent.finished", "a", "segment-0002"),
	}
	state := ProjectEvents([]string{"a"}, events, base.Add(5*time.Second))
	got := mustAgent(t, state, "a")
	if got.State != StateFinished {
		t.Fatalf("a.State = %q, want finished", got.State)
	}
	if got.Error != "" {
		t.Fatalf("retry must clear prior error, got %q", got.Error)
	}
}

// AF5: an untagged agent.* event after a segment boundary inherits the current
// segment (FINAL backward-compat rule).
func TestUntaggedAgentEventInheritsCurrentSegment(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		segStarted(base, "segment-0007", "continue", "a"),
		{Time: base.Add(time.Second), Type: "agent.started", Data: map[string]any{"agent": "a"}}, // no segment_id
	}
	state := ProjectEvents([]string{"a"}, events, base.Add(2*time.Second))
	got := mustAgent(t, state, "a")
	if got.State != StateRunning {
		t.Fatalf("a.State = %q, want running", got.State)
	}
	if got.Segment != "segment-0007" {
		t.Fatalf("untagged event should inherit current segment, got %q", got.Segment)
	}
}

// AF7: the sticky-badge fix holds through the real LoadRunAt path (events.jsonl
// -> projection), not just the direct ProjectEvents unit tests. After a run
// completes and the user continues (a new segment targeting the agent), the
// badge must read pending, not the prior finished.
func TestLoadRunSegmentUnsticksFinishedBadge(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "sample", []string{"codex"})
	runID := "20260604T130000.000000000Z"
	runDir := RunDir(root, runID)
	base := time.Date(2026, 6, 4, 13, 0, 0, 0, time.UTC)
	appendEvents(t, runDir,
		store.Event{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "participants": []string{"codex"}}},
		segStarted(base.Add(time.Second), "segment-0001", "initial", "codex"),
		store.Event{Time: base.Add(2 * time.Second), Type: "agent.started", Data: map[string]any{"agent": "codex", "segment_id": "segment-0001"}},
		store.Event{Time: base.Add(3 * time.Second), Type: "agent.finished", Data: map[string]any{"agent": "codex", "segment_id": "segment-0001"}},
		store.Event{Time: base.Add(4 * time.Second), Type: "round.completed", Data: map[string]any{"completed": float64(1), "total": float64(1)}},
		segStarted(base.Add(5*time.Second), "segment-0002", "continue", "codex"),
	)
	run, err := LoadRunAt(root, runID, base.Add(6*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	got, ok := agentByID(run.State, "codex")
	if !ok {
		t.Fatal("codex missing from projection")
	}
	if got.State != StatePending {
		t.Fatalf("after continue: codex state=%q, want pending (badge must unstick via LoadRunAt)", got.State)
	}
	if got.Segment != "segment-0002" {
		t.Fatalf("codex segment=%q, want segment-0002", got.Segment)
	}
}

// A segment targeting an agent absent from the projection is a no-op (no panic,
// no materialized ghost agent).
func TestSegmentTargetingUnknownAgentIsNoop(t *testing.T) {
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	events := []store.Event{
		segStarted(base, "segment-0001", "initial", "ghost"),
	}
	state := ProjectEvents([]string{"a"}, events, base.Add(time.Second))
	if got := mustAgent(t, state, "a"); got.State != StatePending {
		t.Fatalf("a.State = %q, want pending (untouched)", got.State)
	}
	if _, ok := agentByID(state, "ghost"); ok {
		t.Fatalf("segment targeting must not materialize an unknown agent")
	}
}
