package driver

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/store"
)

// runPhaseEvents loads the run's event log and returns the run.phase events.
func runPhaseEvents(t *testing.T, runDir string) []store.Event {
	t.Helper()
	events, err := store.New(runDir).Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // no log at all → no run.phase events
		}
		t.Fatal(err)
	}
	var out []store.Event
	for _, e := range events {
		if e.Type == "run.phase" {
			out = append(out, e)
		}
	}
	return out
}

// TestPromotedEmitsRunPhase: the round-promotion commit emits exactly one
// run.phase event carrying the action and the new round, after the cursor save.
func TestPromotedEmitsRunPhase(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{writeOnRun: func(round int) { writeAll(t, ideaDir, round, parts) }}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)

	if action, _, err := d.Advance(context.Background()); err != nil || action != ActionPromoted {
		t.Fatalf("action=%v err=%v, want promoted", action, err)
	}
	got := runPhaseEvents(t, runDir)
	if len(got) != 1 {
		t.Fatalf("run.phase events=%d, want 1", len(got))
	}
	data := got[0].Data
	if data["action"] != "promoted" || data["phase"] != "round" || data["previous_phase"] != "round" ||
		data["round_label"] != "round-02" || data["idea"] != "demo" || data["source"] != "driver" {
		t.Fatalf("run.phase data=%+v", data)
	}
	if data["run_id"] != filepath.Base(runDir) {
		t.Fatalf("run_id=%v, want %s", data["run_id"], filepath.Base(runDir))
	}
}

// TestConsensusDraftedEmitsRunPhase: the round→consensus commit emits the
// transition with the correct previous phase.
func TestConsensusDraftedEmitsRunPhase(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	writeAll(t, ideaDir, 2, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	appendEvent(t, runDir, "round.completed", "round-02")
	fc := &fakeConsensus{statusSeq: []string{consensus.TriagePartial}, onDraft: func() { writeConsensusDoc(t, ideaDir, "consensus") }}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, &fakeRunner{})

	if action, _, err := d.Advance(context.Background()); err != nil || action != ActionConsensusDrafted {
		t.Fatalf("action=%v err=%v, want consensus-drafted", action, err)
	}
	got := runPhaseEvents(t, runDir)
	if len(got) != 1 {
		t.Fatalf("run.phase events=%d, want 1", len(got))
	}
	if got[0].Data["action"] != "consensus-drafted" || got[0].Data["phase"] != "consensus" || got[0].Data["previous_phase"] != "round" {
		t.Fatalf("run.phase data=%+v", got[0].Data)
	}
}

// TestAwaitAndSurfaceEmitNoRunPhase: non-committing ticks must not narrate.
func TestAwaitAndSurfaceEmitNoRunPhase(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeArtifact(t, ideaDir, "demo", 1, "codex", parts) // claude missing → await
	d := newTestDriver(ideaDir, runDir, parts, 1, true, &fakeRunner{})
	if action, _, err := d.Advance(context.Background()); err != nil || action != ActionAwait {
		t.Fatalf("action=%v err=%v, want await", action, err)
	}

	off := newTestDriver(ideaDir, runDir, parts, 1, false, &fakeRunner{}) // auto off → surface-only
	if action, _, err := off.Advance(context.Background()); err != nil || action != ActionSurfaceOnly {
		t.Fatalf("action=%v err=%v, want surface-only", action, err)
	}
	if got := runPhaseEvents(t, runDir); len(got) != 0 {
		t.Fatalf("run.phase events=%d, want 0", len(got))
	}
}

// TestSaveFailureEscalatesAndEmitsNothing: a failed cursor save returns an error
// (the branch escalates) and the event log must not claim the phase.
func TestSaveFailureEscalatesAndEmitsNothing(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{writeOnRun: func(round int) { writeAll(t, ideaDir, round, parts) }}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)

	orig := saveCursor
	saveCursor = func(Cursor, string) error { return errors.New("simulated save failure") }
	t.Cleanup(func() { saveCursor = orig })

	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("action=%v err=%v, want escalated with error", action, err)
	}
	if got := runPhaseEvents(t, runDir); len(got) != 0 {
		t.Fatalf("run.phase events=%d, want 0 after save failure", len(got))
	}
}

// TestRebuildDetailReviewEvidence: RebuildDetail carries the display evidence
// that disambiguates review (6) vs review-consensus (7) vs fix-up (8).
func TestRebuildDetailReviewEvidence(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, _ := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	if err := os.WriteFile(filepath.Join(ideaDir, "FINAL.md"), []byte(validFinal), 0o644); err != nil {
		t.Fatal(err)
	}
	impl := "---\nidea: demo\nstatus: fix-up-cycle-1\n---\n\nbody\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"), []byte(impl), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-02"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	detail, err := RebuildDetail(ideaDir, 4)
	if err != nil {
		t.Fatal(err)
	}
	if detail.Cursor.Phase != PhaseReview {
		t.Fatalf("phase=%s, want review", detail.Cursor.Phase)
	}
	if detail.HighestReviewRound != 2 || !detail.ReviewConsensusExists {
		t.Fatalf("detail=%+v, want HighestReviewRound=2 ReviewConsensusExists=true", detail)
	}
	if detail.ImplementationStatus != "fix-up-cycle-1" {
		t.Fatalf("ImplementationStatus=%q", detail.ImplementationStatus)
	}
	if detail.FinalScaffoldReason != "" {
		t.Fatalf("FinalScaffoldReason=%q, want empty for a valid FINAL.md", detail.FinalScaffoldReason)
	}

	// Rebuild stays a faithful wrapper.
	if c := Rebuild(ideaDir, 4); c.Phase != PhaseReview {
		t.Fatalf("Rebuild phase=%s, want review", c.Phase)
	}
}
