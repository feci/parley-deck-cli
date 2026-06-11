package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/store"
)

// --- pure cursor→9-step mapping (consensus D3) ---------------------------------

func TestDisplayStepMapping(t *testing.T) {
	cur := func(p driver.Phase, round int, status string) driver.Cursor {
		return driver.Cursor{Phase: p, CurrentRound: round, IdeaStatus: status}
	}
	cases := []struct {
		name   string
		detail driver.PhaseDetail
		step   int
	}{
		{"kickoff", driver.PhaseDetail{Cursor: cur(driver.PhaseRound, 1, "kickoff")}, 0},
		{"round-01", driver.PhaseDetail{Cursor: cur(driver.PhaseRound, 1, "round-01")}, 1},
		{"cross-review", driver.PhaseDetail{Cursor: cur(driver.PhaseRound, 2, "round-02")}, 2},
		{"consensus", driver.PhaseDetail{Cursor: cur(driver.PhaseConsensus, 2, "consensus")}, 3},
		{"final", driver.PhaseDetail{Cursor: cur(driver.PhaseFinal, 2, "final")}, 4},
		{"implement", driver.PhaseDetail{Cursor: cur(driver.PhaseImpl, 2, "implemented")}, 5},
		{"review", driver.PhaseDetail{Cursor: cur(driver.PhaseReview, 2, "implemented"), HighestReviewRound: 1}, 6},
		{"review-consensus", driver.PhaseDetail{Cursor: cur(driver.PhaseReview, 2, "implemented"), HighestReviewRound: 1, ReviewConsensusExists: true}, 7},
		{"fix-up", driver.PhaseDetail{Cursor: cur(driver.PhaseReview, 2, "implemented"), HighestReviewRound: 1, ReviewConsensusExists: true, ImplementationStatus: "fix-up-cycle-1"}, 8},
		{"complete", driver.PhaseDetail{Cursor: cur(driver.PhaseDone, 2, "complete"), ImplementationStatus: "complete"}, 8},
	}
	for _, tc := range cases {
		if step, _ := displayStep(tc.detail); step != tc.step {
			t.Errorf("%s: step=%d, want %d", tc.name, step, tc.step)
		}
	}
}

// --- snapshot builder over a real idea dir --------------------------------------

// setupSnapIdea creates root/parley-deck/ideas/demo at cross-review (round-02).
func setupSnapIdea(t *testing.T) (root, ideaDir string) {
	t.Helper()
	root = t.TempDir()
	ideaDir = filepath.Join(root, "parley-deck", "ideas", "demo")
	if err := os.MkdirAll(filepath.Join(ideaDir, "round-02"), 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\nidea: demo\nparticipants: [codex, agy, hermes]\ntransport: local-dir\ncross_review_rounds: 1\nstatus: round-02\n---\n\n## Problem\nx\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, ideaDir
}

func snapInput(root, ideaDir string, events []store.Event, prev *ProtocolSnapshot, now time.Time) ProtocolSnapshotInput {
	parts := []string{"codex", "agy", "hermes"}
	return ProtocolSnapshotInput{
		Root: root, RunID: "run1", RunDir: filepath.Join(root, "parley-deck", "runs", "run1"),
		IdeaSlug: "demo", IdeaDir: ideaDir,
		Participants: parts, MaxRounds: 4,
		Events: events, State: runstate.ProjectEvents(parts, events, now),
		Previous: prev, Now: now,
	}
}

// TestSnapshotDeliveryMergeRules: events primary, disk fallback, `?` on
// disagreement (consensus D5 merge table).
func TestSnapshotDeliveryMergeRules(t *testing.T) {
	root, ideaDir := setupSnapIdea(t)
	now := time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC)
	// codex: file on disk, NO event → delivered, unvalidated, disk fallback.
	if err := os.WriteFile(filepath.Join(ideaDir, "round-02", "codex.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	// agy: agent.finished event, file MISSING → delivered, unvalidated.
	events := []store.Event{
		{Time: now, Type: "agent.started", Data: map[string]any{"agent": "agy", "artifact": filepath.Join(ideaDir, "round-02", "agy.md")}},
		{Time: now.Add(time.Minute), Type: "agent.finished", Data: map[string]any{"agent": "agy", "artifact": filepath.Join(ideaDir, "round-02", "agy.md"), "artifact_ok": true}},
	}
	// hermes: nothing → pending + waiting.

	snap, err := BuildProtocolSnapshot(snapInput(root, ideaDir, events, nil, now))
	if err != nil {
		t.Fatal(err)
	}
	if snap.Step != 2 || snap.RoundLabel != "round-02" || snap.TotalRounds != 2 {
		t.Fatalf("snap step=%d round=%s total=%d, want 2/round-02/2", snap.Step, snap.RoundLabel, snap.TotalRounds)
	}
	rows := map[string]AgentDelivery{}
	for _, row := range snap.Delivery {
		rows[row.ID] = row
	}
	if r := rows["codex"]; r.State != "delivered" || !r.Unvalidated {
		t.Fatalf("codex row=%+v, want delivered+unvalidated (disk only)", r)
	}
	if r := rows["agy"]; r.State != "delivered" || !r.Unvalidated {
		t.Fatalf("agy row=%+v, want delivered+unvalidated (event only, file missing)", r)
	}
	if r := rows["hermes"]; r.State != "pending" {
		t.Fatalf("hermes row=%+v, want pending", r)
	}
	if !snap.DiskFallback {
		t.Fatal("DiskFallback must be true when a row came from disk")
	}
	if len(snap.Waiting) != 1 || snap.Waiting[0] != "hermes" {
		t.Fatalf("waiting=%v, want [hermes]", snap.Waiting)
	}
}

// TestSnapshotKeepLastOnError: an unreadable idea dir keeps the previous
// snapshot and surfaces the reconcile error (consensus D5 keep-last).
func TestSnapshotKeepLastOnError(t *testing.T) {
	root, ideaDir := setupSnapIdea(t)
	now := time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC)
	// Make review/ a regular FILE so the review ReadDir fails with ENOTDIR.
	if err := os.WriteFile(filepath.Join(ideaDir, "review"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	prev := &ProtocolSnapshot{Step: 2, StepName: "cross-review", ReconciledAt: now.Add(-time.Minute)}

	snap, err := BuildProtocolSnapshot(snapInput(root, ideaDir, nil, prev, now))
	if err == nil {
		t.Fatal("want a reconcile error for the unreadable review dir")
	}
	if snap.Step != 2 || snap.Err == "" {
		t.Fatalf("snap=%+v, want previous step retained with Err set", snap)
	}
	if !snap.ReconciledAt.Equal(now) {
		t.Fatalf("ReconciledAt=%v, want refreshed to now", snap.ReconciledAt)
	}
}

// TestSnapshotRegressionNeedsTwoAgreeingReconciles: a single lower-step read
// (virtio-fs stale dir) must NOT bounce the pipeline backwards.
func TestSnapshotRegressionNeedsTwoAgreeingReconciles(t *testing.T) {
	root, ideaDir := setupSnapIdea(t)
	now := time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC)
	prev := &ProtocolSnapshot{Step: 3, StepName: "consensus"} // previously at consensus

	first, err := BuildProtocolSnapshot(snapInput(root, ideaDir, nil, prev, now))
	if err != nil {
		t.Fatal(err)
	}
	if first.Step != 3 {
		t.Fatalf("first reconcile regressed immediately to %d; want previous step 3 kept", first.Step)
	}
	second, err := BuildProtocolSnapshot(snapInput(root, ideaDir, nil, &first, now.Add(15*time.Second)))
	if err != nil {
		t.Fatal(err)
	}
	if second.Step != 2 {
		t.Fatalf("second agreeing reconcile step=%d, want regression accepted to 2", second.Step)
	}
}

// --- model-side gating (consensus D5: budget + coalescing) ----------------------

// TestTickBudgetAndCoalescing: ticks alone schedule no snapshot; only trigger
// events do; a trigger while busy coalesces into exactly one follow-up build.
func TestTickBudgetAndCoalescing(t *testing.T) {
	m := newLiveModel(LiveOptions{Participants: []string{"codex"}, RunID: "r", RunDir: t.TempDir()})
	for i := 0; i < 60; i++ {
		next, _ := m.Update(eventTickMsg{token: m.runToken})
		m = next.(liveModel)
		next, _ = m.Update(elapsedTickMsg{t: m.now, token: m.runToken})
		m = next.(liveModel)
	}
	if m.protoSeq != 0 {
		t.Fatalf("protoSeq=%d after 60 event+elapsed ticks, want 0 (no snapshot work on ticks)", m.protoSeq)
	}

	// A non-trigger event batch does not schedule.
	next, _ := m.Update(eventsMsg{events: []store.Event{{Type: "steer.requested"}}, token: m.runToken})
	m = next.(liveModel)
	if m.protoSeq != 0 {
		t.Fatalf("protoSeq=%d after steer event, want 0", m.protoSeq)
	}

	// A trigger schedules one build; a second trigger while busy only marks dirty.
	next, cmd := m.Update(eventsMsg{events: []store.Event{{Type: "round.completed"}}, token: m.runToken})
	m = next.(liveModel)
	if m.protoSeq != 1 || cmd == nil || !m.protoBusy {
		t.Fatalf("first trigger: seq=%d busy=%v cmd=%v, want 1/true/non-nil", m.protoSeq, m.protoBusy, cmd)
	}
	next, cmd = m.Update(eventsMsg{events: []store.Event{{Type: "agent.finished", Data: map[string]any{"agent": "codex"}}}, token: m.runToken})
	m = next.(liveModel)
	if m.protoSeq != 1 || !m.protoDirty {
		t.Fatalf("trigger while busy: seq=%d dirty=%v, want 1/true", m.protoSeq, m.protoDirty)
	}
	// Completion of the in-flight build re-fires exactly once for the dirty flag.
	next, cmd = m.Update(protoMsg{snap: ProtocolSnapshot{Step: 1}, token: m.runToken, seq: 1})
	m = next.(liveModel)
	if m.protoSeq != 2 || m.protoDirty || cmd == nil {
		t.Fatalf("dirty completion: seq=%d dirty=%v cmd=%v, want 2/false/non-nil", m.protoSeq, m.protoDirty, cmd)
	}
	// A stale seq is dropped.
	next, _ = m.Update(protoMsg{snap: ProtocolSnapshot{Step: 7}, token: m.runToken, seq: 1})
	m = next.(liveModel)
	if m.proto != nil && m.proto.Step == 7 {
		t.Fatal("stale protoMsg (old seq) must be dropped")
	}
}

// TestStatusPhaseSegment: the compressed ph= grammar with waiting list (D10).
func TestStatusPhaseSegment(t *testing.T) {
	m := newLiveModel(LiveOptions{Participants: []string{"codex"}})
	if got := m.statusPhaseSegment(); !strings.HasPrefix(got, "round=") {
		t.Fatalf("fallback segment=%q, want legacy round= before the first snapshot", got)
	}
	m.proto = &ProtocolSnapshot{Step: 2, StepName: "cross-review", RoundLabel: "round-02", Waiting: []string{"agy", "hermes"}}
	if got := m.statusPhaseSegment(); got != "ph=2:xrev-r02 wait=agy,hermes" {
		t.Fatalf("segment=%q, want ph=2:xrev-r02 wait=agy,hermes", got)
	}
}

// --- narrator ring + glyphs -------------------------------------------------------

func TestNarratorWeaveAndReplay(t *testing.T) {
	m := newLiveModel(LiveOptions{Participants: []string{"codex", "agy"}})
	loaded := &agentBuffer{loaded: true, follow: true, partial: map[transcriptStream]string{}, crPending: map[transcriptStream]bool{}}
	m.buffers["codex"] = loaded

	m.appendProtocolEvents([]store.Event{
		{Type: "round.completed", Data: map[string]any{"completed": "3", "total": "3"}},
		{Type: "agent.acp.message_chunk", Data: map[string]any{"agent": "codex"}}, // excluded chatter
	})
	if len(loaded.lines) != 1 || loaded.lines[0].Stream != transcriptEvent {
		t.Fatalf("loaded buffer lines=%+v, want exactly one woven narrator line", loaded.lines)
	}
	if len(m.narratorRing) != 1 {
		t.Fatalf("ring=%d, want 1 (chatter excluded)", len(m.narratorRing))
	}

	// A buffer that loads later replays the ring exactly once.
	late := &agentBuffer{loaded: true, partial: map[transcriptStream]string{}, crPending: map[transcriptStream]bool{}}
	m.replayNarrator(late)
	m.replayNarrator(late) // second replay must be a no-op (seq dedup)
	if len(late.lines) != 1 {
		t.Fatalf("late buffer lines=%d, want 1 after dedup'd replay", len(late.lines))
	}

	// narrate off → nothing woven.
	m.narrateMode = narrateOff
	m.appendProtocolEvents([]store.Event{{Type: "round.completed"}})
	if len(m.narratorRing) != 1 {
		t.Fatalf("ring grew while narrator off: %d", len(m.narratorRing))
	}
}

func TestAgentGlyphSpinnerVsSilent(t *testing.T) {
	m := newLiveModel(LiveOptions{Participants: []string{"codex"}})
	m.state.Agents = []AgentState{{ID: "codex", State: stateRunning}}

	// Silent running agent → ·
	if got := m.agentGlyph("codex"); got != "·" {
		t.Fatalf("silent glyph=%q, want ·", got)
	}
	// Recent growth via loaded buffer → spinner frame.
	m.buffers["codex"] = &agentBuffer{loaded: true, lastGrowthAt: m.now.Add(-time.Second)}
	if got := m.agentGlyph("codex"); !strings.Contains(strings.Join(spinnerFrames[:], ""), got) {
		t.Fatalf("flowing glyph=%q, want a spinner frame", got)
	}
	// Growth cache covers unvisited tabs.
	delete(m.buffers, "codex")
	m.growth = map[string]growthInfo{"codex": {lastGrowth: m.now.Add(-2 * time.Second)}}
	if got := m.agentGlyph("codex"); !strings.Contains(strings.Join(spinnerFrames[:], ""), got) {
		t.Fatalf("cache-flowing glyph=%q, want a spinner frame", got)
	}
	// STALE overlay wins over running.
	m.opts.Liveness = func(string) string { return "stale" }
	if got := m.agentGlyph("codex"); got != "!" {
		t.Fatalf("stale glyph=%q, want !", got)
	}
}
