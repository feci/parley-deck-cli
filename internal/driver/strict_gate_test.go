package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/store"
)

// newStrictDriver builds an impl-phase driver with strict_gate enabled (LE-2).
func newStrictDriver(ideaDir, runDir string, parts []string, maxFixup int, fi ImplOps) *Driver {
	return New(Config{
		IdeaDir: ideaDir, IdeaSlug: "demo", Participants: parts, RunDir: runDir,
		Root: filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))), Events: store.New(runDir),
		Auto: true, AutoImplement: true, StrictGate: true, MaxFixupCycles: maxFixup, Impl: fi,
	}, &fakeRunner{})
}

func reviewReady(fixes int) ReviewStatus {
	return ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: fixes}
}

// LE-2: 0 outstanding fixes but the round is NOT certified clean → strict_gate must
// open one more fresh closing review round instead of completing.
func TestStrictGateOpensFreshClosingRound(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\nstrict_gate: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	rs := reviewReady(0)
	rs.StrictGateClean = false
	fi := &fakeImpl{
		roundComplete: true,
		review:        rs,
		onOpenReview:  func(round int) { os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(round)), 0o755) },
	}
	d := newStrictDriver(ideaDir, runDir, parts, 3, fi)
	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionReviewOpened || c.Phase != PhaseReview {
		t.Fatalf("action=%s phase=%s want review-opened/review (fresh strict closing round)", action, c.Phase)
	}
	if contains(fi.calls, "complete") {
		t.Fatal("must NOT complete before a certified-clean closing round")
	}
	if !contains(fi.calls, "open-review") {
		t.Fatalf("expected open-review, calls=%v", fi.calls)
	}
	if fileExists(filepath.Join(ideaDir, "review", "consensus.md")) {
		t.Fatal("root review/consensus.md should be archived before the closing round")
	}
}

// LE-2: 0 outstanding fixes, certified clean for THIS round, and the scan finds no
// findings → complete.
func TestStrictGateCompletesWhenCertifiedClean(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\nstrict_gate: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	rs := reviewReady(0)
	rs.StrictGateClean = true
	rs.ClosingReviewRound = 1
	fi := &fakeImpl{roundComplete: true, review: rs}
	d := newStrictDriver(ideaDir, runDir, parts, 3, fi)
	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionComplete || c.Phase != PhaseDone {
		t.Fatalf("action=%s phase=%s want complete/done", action, c.Phase)
	}
	if !contains(fi.calls, "complete") {
		t.Fatalf("expected complete, calls=%v", fi.calls)
	}
}

// LE-2: the deterministic scan vetoes a clean certification that is contradicted by a
// real finding in a review file (fail closed).
func TestStrictGateVetoesUncleanCertification(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\nstrict_gate: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "review", "round-01", "codex.md"),
		[]byte("---\nagent: codex\nidea: demo\nreview-round: 1\n---\n\n## Refutation attempts\ntried\n\n## Findings\n### [MAJOR] real bug\nbroken\n"), 0o644)
	rs := reviewReady(0)
	rs.StrictGateClean = true
	rs.ClosingReviewRound = 1
	fi := &fakeImpl{roundComplete: true, review: rs}
	d := newStrictDriver(ideaDir, runDir, parts, 3, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("action=%s err=%v want escalated (scan veto)", action, err)
	}
	if contains(fi.calls, "complete") {
		t.Fatal("must NOT complete when the finding scan vetoes the clean claim")
	}
}

// LE-2: the strict-close loop is bounded by MaxFixupCycles so it always terminates.
func TestStrictGateBoundedByMaxFixupCycles(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\nstrict_gate: true\n")
	writeFinalValid(t, ideaDir)
	writeImpl(t, ideaDir, "implemented")
	for r := 1; r <= 3; r++ {
		os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(r)), 0o755)
	}
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	rs := reviewReady(0)
	rs.StrictGateClean = false
	fi := &fakeImpl{roundComplete: true, review: rs}
	d := newStrictDriver(ideaDir, runDir, parts, 3, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("action=%s err=%v want escalated (strict-close bound)", action, err)
	}
	if contains(fi.calls, "complete") {
		t.Fatal("must not complete past the strict-close bound")
	}
}

// scanHasRealFinding ignores the template placeholder but catches a concrete finding.
func TestScanHasRealFinding(t *testing.T) {
	placeholder := "## Findings\n### [CRITICAL] <title>\n### [NIT] <title>\n"
	if scanHasRealFinding(placeholder) {
		t.Fatal("template placeholders must not count as findings")
	}
	real := "## Findings\n### [MINOR] off-by-one in loop bound\nfix it\n"
	if !scanHasRealFinding(real) {
		t.Fatal("a concrete finding heading must count")
	}
}
