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
