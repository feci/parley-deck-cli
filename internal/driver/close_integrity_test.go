package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/consensus"
)

func closeReady(triage string, fixes, reviewers int) ReviewStatus {
	return ReviewStatus{Summary: consensus.Summary{Triage: triage}, OutstandingAgreedFixes: fixes, ReviewerCount: reviewers}
}

// LE-11: under auto_implement, ACCEPT-WITH-RESERVATIONS must not auto-complete — escalate.
func TestReservedEscalatesUnderAuto(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, review: closeReady(consensus.TriageReserved, 0, 2)}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("Reserved under auto must escalate; action=%s err=%v", action, err)
	}
	if contains(fi.calls, "complete") {
		t.Fatal("must not complete on reservations under auto")
	}
}

// LE-11: under auto_implement, fewer than 2 independent reviewers must not auto-complete.
func TestSingleReviewerEscalatesUnderAuto(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, review: closeReady(consensus.TriageReady, 0, 1)}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("<2 reviewers under auto must escalate; action=%s err=%v", action, err)
	}
	if contains(fi.calls, "complete") {
		t.Fatal("must not auto-complete with a single reviewer")
	}
}

// LE-7: a confident goal-check FAIL escalates before completion.
func TestGoalCheckFailEscalatesUnderAuto(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, review: closeReady(consensus.TriageReady, 0, 2), goalFail: true}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("goal-check FAIL must escalate; action=%s err=%v", action, err)
	}
	if !contains(fi.calls, "goal-check") {
		t.Fatal("expected the goal-check to run before completion")
	}
	if contains(fi.calls, "complete") {
		t.Fatal("must not complete when the goal-check fails")
	}
}

// Ready + ≥2 reviewers + goal-check pass → completes (the happy path through the gates).
func TestCompletesWhenGatesPassUnderAuto(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, review: closeReady(consensus.TriageReady, 0, 2)}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionComplete || c.Phase != PhaseDone {
		t.Fatalf("action=%s phase=%s want complete/done", action, c.Phase)
	}
	if !contains(fi.calls, "goal-check") || !contains(fi.calls, "complete") {
		t.Fatalf("expected goal-check then complete, calls=%v", fi.calls)
	}
}

// CF5 (hermes MINOR #2): a strict_gate design-only idea (StrictGate, AutoImplement=false)
// with a Reserved triage and a single reviewer still COMPLETES — LE-11 is auto-only — but
// the LE-7 goal-check DOES run (gated on AutoImplement || StrictGate). Documents the
// conditional-rigor boundary so a regression that wrongly gates LE-7 on auto, or fires
// LE-11 for design ideas, is caught.
func TestStrictDesignOnlyCompletesAndRunsGoalCheck(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "strict_gate: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	rs := closeReady(consensus.TriageReserved, 0, 1)
	rs.StrictGateClean = true
	rs.ClosingReviewRound = 1
	fi := &fakeImpl{roundComplete: true, review: rs}
	d := newStrictDesignDriver(ideaDir, runDir, parts, 3, fi)
	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionComplete || c.Phase != PhaseDone {
		t.Fatalf("action=%s phase=%s want complete/done (conditional-rigor boundary)", action, c.Phase)
	}
	if !contains(fi.calls, "goal-check") {
		t.Fatalf("LE-7 goal-check must run for a strict design idea even without auto, calls=%v", fi.calls)
	}
}

// CF5: the same strict design-only path with a confident goal-check FAIL still escalates —
// LE-7 fires on StrictGate independent of AutoImplement.
func TestStrictDesignOnlyGoalCheckFailEscalates(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "strict_gate: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	rs := closeReady(consensus.TriageReserved, 0, 1)
	rs.StrictGateClean = true
	rs.ClosingReviewRound = 1
	fi := &fakeImpl{roundComplete: true, review: rs, goalFail: true}
	d := newStrictDesignDriver(ideaDir, runDir, parts, 3, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("strict design goal-check FAIL must escalate; action=%s err=%v", action, err)
	}
	if contains(fi.calls, "complete") {
		t.Fatal("must not complete when the strict-design goal-check fails")
	}
}

// Conditional rigor: a non-auto (design-only) idea with a single reviewer still completes
// and does NOT run the goal-check (the LE-7/LE-11 gates are auto-only).
func TestNonAutoCompletesWithSingleReviewer(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, review: closeReady(consensus.TriageReady, 0, 1)}
	d := newImplDriver(ideaDir, runDir, parts, false, fi)
	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionComplete || c.Phase != PhaseDone {
		t.Fatalf("action=%s phase=%s want complete/done", action, c.Phase)
	}
	if contains(fi.calls, "goal-check") {
		t.Fatal("a non-auto idea must not run the goal-check gate")
	}
}
