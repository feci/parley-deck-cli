package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/store"
)

// fakeImpl is an injected ImplOps that records calls and returns scripted values,
// so the impl/review gate is tested without live agents (consensus D14).
type fakeImpl struct {
	calls         []string
	status        string // ImplementationStatus
	statusErr     error
	checksOK      bool
	roundComplete bool
	review        ReviewStatus
	reviewErr     error
	goalFail      bool
	onOpenReview  func(round int)
	onDraft       func()
	onComplete    func()
}

func (f *fakeImpl) GoalCheck(ctx context.Context) (bool, string) {
	f.calls = append(f.calls, "goal-check")
	if f.goalFail {
		return false, "criteria unmet"
	}
	return true, ""
}

func (f *fakeImpl) Implement(ctx context.Context) error {
	f.calls = append(f.calls, "implement")
	return nil
}
func (f *fakeImpl) ImplementationStatus() (string, error) {
	return f.status, f.statusErr
}
func (f *fakeImpl) RunChecks(ctx context.Context) (bool, string) {
	f.calls = append(f.calls, "checks")
	return f.checksOK, "detail"
}
func (f *fakeImpl) OpenReviewRound(ctx context.Context, round int) error {
	f.calls = append(f.calls, "open-review")
	if f.onOpenReview != nil {
		f.onOpenReview(round)
	}
	return nil
}
func (f *fakeImpl) ReviewRoundComplete(round int) (bool, error) { return f.roundComplete, nil }
func (f *fakeImpl) DraftReviewConsensus(ctx context.Context, round int) error {
	f.calls = append(f.calls, "draft-review")
	if f.onDraft != nil {
		f.onDraft()
	}
	return nil
}
func (f *fakeImpl) ReviewStatus() (ReviewStatus, error) { return f.review, f.reviewErr }
func (f *fakeImpl) RequestReviewSignoffs(ctx context.Context, missing []string) error {
	f.calls = append(f.calls, "request-signoffs")
	return nil
}
func (f *fakeImpl) Fixup(ctx context.Context, cycle int) error {
	f.calls = append(f.calls, "fixup")
	return nil
}
func (f *fakeImpl) Complete(ctx context.Context) error {
	f.calls = append(f.calls, "complete")
	if f.onComplete != nil {
		f.onComplete()
	}
	return nil
}

func writeFinalValid(t *testing.T, ideaDir string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(ideaDir, "FINAL.md"), []byte(validFinal), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeImpl(t *testing.T, ideaDir, status string) {
	t.Helper()
	body := "---\nidea: demo\nstatus: " + status + "\n---\n\n## Summary of work\nx\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func newImplDriver(ideaDir, runDir string, parts []string, autoImplement bool, fi ImplOps) *Driver {
	return New(Config{
		IdeaDir:       ideaDir,
		IdeaSlug:      "demo",
		Participants:  parts,
		RunDir:        runDir,
		Root:          filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))),
		Events:        store.New(runDir),
		Auto:          true,
		AutoImplement: autoImplement,
		Impl:          fi,
	}, &fakeRunner{})
}

func TestRebuildImplReviewDonePrecedence(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, _ := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	if c := Rebuild(ideaDir, 4); c.Phase != PhaseFinal {
		t.Fatalf("FINAL only: phase=%s want final", c.Phase)
	}
	writeImpl(t, ideaDir, "implemented")
	if c := Rebuild(ideaDir, 4); c.Phase != PhaseImpl {
		t.Fatalf("IMPLEMENTATION present: phase=%s want impl (must not be hidden by FINAL)", c.Phase)
	}
	if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-01"), 0o755); err != nil {
		t.Fatal(err)
	}
	if c := Rebuild(ideaDir, 4); c.Phase != PhaseReview {
		t.Fatalf("review round present: phase=%s want review", c.Phase)
	}
	writeImpl(t, ideaDir, "complete")
	if c := Rebuild(ideaDir, 4); c.Phase != PhaseDone {
		t.Fatalf("status complete: phase=%s want done", c.Phase)
	}
}

func TestPhaseFinalImplementsWhenOptedIn(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	fi := &fakeImpl{}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)

	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionImplemented {
		t.Fatalf("action=%s want implemented", action)
	}
	if c.Phase != PhaseImpl || len(fi.calls) != 1 || fi.calls[0] != "implement" {
		t.Fatalf("phase=%s calls=%v", c.Phase, fi.calls)
	}
	wantLastRunPhase(t, runDir, "implemented", PhaseImpl, PhaseFinal)
}

func TestPhaseFinalSurfaceWhenNotOptedIn(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "") // no auto_implement
	writeFinalValid(t, ideaDir)
	fi := &fakeImpl{}
	d := newImplDriver(ideaDir, runDir, parts, false, fi)

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionSurfaceOnly {
		t.Fatalf("action=%s want surface-only (no opt-in)", action)
	}
	if len(fi.calls) != 0 {
		t.Fatal("Implement must not run without auto_implement")
	}
}

func TestPhaseImplChecksGate(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImpl(t, ideaDir, "implemented")

	// Checks fail → escalate, no review opened.
	fiFail := &fakeImpl{status: "implemented", checksOK: false}
	d := newImplDriver(ideaDir, runDir, parts, true, fiFail)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("checks fail: action=%s err=%v want escalated", action, err)
	}
	for _, c := range fiFail.calls {
		if c == "open-review" {
			t.Fatal("must not open review when checks fail")
		}
	}

	// Checks pass → open review round 1.
	fiOK := &fakeImpl{status: "implemented", checksOK: true, onOpenReview: func(round int) {
		os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(round)), 0o755)
	}}
	d = newImplDriver(ideaDir, runDir, parts, true, fiOK)
	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionReviewOpened || c.Phase != PhaseReview {
		t.Fatalf("checks pass: action=%s phase=%s want review-opened/review", action, c.Phase)
	}
	wantLastRunPhase(t, runDir, "review-opened", PhaseReview, PhaseImpl)
}

func TestPhaseImplNotReadyEscalates(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImpl(t, ideaDir, "in-progress-nonsense")
	fi := &fakeImpl{status: "in-progress-nonsense"}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("non-ready status: action=%s err=%v want escalated", action, err)
	}
}

func setupReviewPhase(t *testing.T, autoImpl string) (ideaDir, runDir string, parts []string) {
	t.Helper()
	parts = []string{"codex", "agy"}
	ideaDir, runDir = setupIdea(t, parts, autoImpl)
	writeFinalValid(t, ideaDir)
	writeImpl(t, ideaDir, "implemented")
	if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-01"), 0o755); err != nil {
		t.Fatal(err)
	}
	return ideaDir, runDir, parts
}

func TestPhaseReviewDraftsThenCompletes(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	// round complete, no consensus yet → draft.
	fi := &fakeImpl{roundComplete: true, onDraft: func() {
		os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	}}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionReviewDrafted {
		t.Fatalf("action=%s want review-drafted", action)
	}

	// Now consensus exists, Ready + zero fixes + ≥2 reviewers → complete (LE-11).
	fi.review = ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 0, ReviewerCount: 2}
	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionComplete || c.Phase != PhaseDone {
		t.Fatalf("action=%s phase=%s want complete/done", action, c.Phase)
	}
	if !contains(fi.calls, "complete") {
		t.Fatalf("expected Complete call, calls=%v", fi.calls)
	}
	wantLastRunPhase(t, runDir, "complete", PhaseDone, PhaseReview)
}

func TestPhaseReviewFixupWhenFixesOutstanding(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{
		roundComplete: true,
		checksOK:      true, // AF1: RunChecks runs after Fixup before the next round
		review:        ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 2},
		onOpenReview:  func(round int) { os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(round)), 0o755) },
	}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionFixup {
		t.Fatalf("action=%s want fixup", action)
	}
	if !contains(fi.calls, "fixup") || !contains(fi.calls, "checks") || !contains(fi.calls, "open-review") {
		t.Fatalf("expected fixup + checks + open-review, calls=%v", fi.calls)
	}
	if !fileExists(filepath.Join(ideaDir, "review", "round-01", ".fixup-done")) {
		t.Fatal("AF2: expected a .fixup-done marker after the fix-up")
	}
	// consensus archived to review/round-01/consensus.md, root cleared.
	if fileExists(filepath.Join(ideaDir, "review", "consensus.md")) {
		t.Fatal("root review/consensus.md should be archived after fixup")
	}
	if !fileExists(filepath.Join(ideaDir, "review", "round-01", "consensus.md")) {
		t.Fatal("expected archived review/round-01/consensus.md")
	}
	wantLastRunPhase(t, runDir, "fixup", PhaseReview, PhaseReview)
}

func TestPhaseReviewFixupChecksFailEscalates(t *testing.T) {
	// AF1: a fix-up that breaks the build escalates and does NOT open the next round.
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	opened := false
	fi := &fakeImpl{
		roundComplete: true,
		checksOK:      false, // checks fail after fix-up
		review:        ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1},
		onOpenReview:  func(round int) { opened = true },
	}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("action=%s err=%v want escalated", action, err)
	}
	if !contains(fi.calls, "fixup") || !contains(fi.calls, "checks") {
		t.Fatalf("expected fixup + checks, calls=%v", fi.calls)
	}
	if opened {
		t.Fatal("must not open the next review round when post-fix-up checks fail")
	}
	if fileExists(filepath.Join(ideaDir, "review", "round-01", ".fixup-done")) {
		t.Fatal("must not write the fix-up marker when checks fail")
	}
}

func TestPhaseReviewFixupMarkerSkipsRefixup(t *testing.T) {
	// AF2: a present .fixup-done marker (crash before the next round opened) must
	// finish the transition without re-running Fixup.
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "review", "round-01", ".fixup-done"), []byte("done"), 0o644)
	fi := &fakeImpl{
		review:       ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1},
		onOpenReview: func(round int) { os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(round)), 0o755) },
	}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionFixup {
		t.Fatalf("action=%s want fixup", action)
	}
	if contains(fi.calls, "fixup") {
		t.Fatal("AF2: must NOT re-run Fixup when the marker is present")
	}
	if !contains(fi.calls, "open-review") {
		t.Fatalf("expected open-review, calls=%v", fi.calls)
	}
	wantLastRunPhase(t, runDir, "fixup", PhaseReview, PhaseReview)
}

func TestPhaseImplInProgressAwaits(t *testing.T) {
	// AF7: a known in-progress status awaits rather than escalating.
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImpl(t, ideaDir, "in-progress")
	fi := &fakeImpl{status: "in-progress"}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionAwait {
		t.Fatalf("action=%s want await for in-progress status", action)
	}
}

func TestPhaseReviewMaxFixupCyclesEscalates(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImpl(t, ideaDir, "implemented")
	// At the cycle ceiling: review/round-03 exists (cycle 3 == MaxFixupCycles).
	for r := 1; r <= 3; r++ {
		os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(r)), 0o755)
	}
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}
	d := New(Config{IdeaDir: ideaDir, IdeaSlug: "demo", Participants: parts, RunDir: runDir,
		Root: filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))), Events: store.New(runDir),
		Auto: true, AutoImplement: true, MaxFixupCycles: 3, Impl: fi}, &fakeRunner{})
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("max cycles: action=%s err=%v want escalated", action, err)
	}
	if contains(fi.calls, "fixup") {
		t.Fatal("must not run fixup past MaxFixupCycles")
	}
}

func TestPhaseReviewBlockedEscalates(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageBlocked}, Blocked: true}}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil || action != ActionEscalated {
		t.Fatalf("blocked: action=%s err=%v want escalated", action, err)
	}
}

func contains(s []string, want string) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

// TestPhaseReviewListChecksVetoCompletion: with list-form checks:, a review consensus
// that is Ready with zero agreed fixes must NOT complete when the completion contract
// fails at HEAD (review fix for the CRITICAL: the zero-fixes path was un-gated).
func TestPhaseReviewListChecksVetoCompletion(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\nchecks:\n  - name: unit\n    command: \"true\"\n")
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{
		roundComplete: true,
		checksOK:      false, // contract fails at HEAD
		review:        ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 0, ReviewerCount: 2},
	}
	d := newImplDriver(ideaDir, runDir, parts, true, fi)
	action, _, err := d.Advance(context.Background())
	if err == nil {
		t.Fatalf("expected completion-contract veto to escalate, got action=%s", action)
	}
	if action != ActionEscalated {
		t.Fatalf("action=%s want escalated", action)
	}
	if contains(fi.calls, "complete") {
		t.Fatalf("Complete must NOT be called when the contract fails; calls=%v", fi.calls)
	}
}
