package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	fixupErr      bool
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
	if f.fixupErr {
		return fmt.Errorf("synthetic fix-up failure")
	}
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

// writeImplWithCycles records n already-published fix-up cycles the way the DRIVER does:
// one `.fixup-done` marker per review round. The §4.0 budget counts driver-CHARGED
// attempts, deliberately not `## Fix-up cycle` headings in the implementer-owned
// IMPLEMENTATION.md — see chargedFixupAttempts for why.
func writeImplWithCycles(t *testing.T, ideaDir, status string, n int) {
	t.Helper()
	writeImpl(t, ideaDir, status)
	for i := 1; i <= n; i++ {
		rd := filepath.Join(ideaDir, "review", roundLabel(i))
		if err := os.MkdirAll(rd, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(rd, ".fixup-done"), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
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
	// PAST the ceiling: the cap is inclusive and counts PUBLISHED cycles, so with
	// MaxFixupCycles=3 and three cycles already published, the next one must escalate.
	writeImplWithCycles(t, ideaDir, "implemented", 3)
	os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(4)), 0o755)
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

// Boundary test for the inclusive fix-up cap (idea
// meta-protocol-change-phase-packet-and-fixup-budget). The §4.0 table prints "cap N
// cycles" for every track; before this idea the guard was `cycle >= cap`, which published
// only N-1 — so `fast` (cap 1) published none at all. The cap is inclusive: cycles 1..N
// run, and N+1 escalates.
func TestFixupCapIsInclusive(t *testing.T) {
	for _, tc := range []struct {
		name      string
		cap       int
		published int // fix-up cycles ALREADY published; the next one is published+1
		wantFixup bool
	}{
		{"cap 5: the 5th published cycle is allowed", 5, 4, true},
		{"cap 5: the 6th escalates", 5, 5, false},
		{"cap 2: the 2nd is allowed", 2, 1, true},
		{"cap 2: the 3rd escalates", 2, 2, false},
		{"cap 1: the 1st is allowed", 1, 0, true},
		{"cap 1: the 2nd escalates", 1, 1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parts := []string{"codex", "agy"}
			ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
			writeFinalValid(t, ideaDir)
			// Rounds 1..published each carry a driver `.fixup-done` marker; the CURRENT
			// round (published+1) has none, so this Advance is the one that decides
			// whether the next cycle may run.
			writeImplWithCycles(t, ideaDir, "implemented", tc.published)
			os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(tc.published+1)), 0o755)
			os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
			fi := &fakeImpl{roundComplete: true, review: ReviewStatus{
				Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}
			d := New(Config{IdeaDir: ideaDir, IdeaSlug: "demo", Participants: parts, RunDir: runDir,
				Root: filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))), Events: store.New(runDir),
				Auto: true, AutoImplement: true, MaxFixupCycles: tc.cap, Impl: fi}, &fakeRunner{})
			action, _, err := d.Advance(context.Background())
			ranFixup := contains(fi.calls, "fixup")
			if ranFixup != tc.wantFixup {
				t.Fatalf("cap=%d published=%d: ranFixup=%v want %v (action=%s err=%v)",
					tc.cap, tc.published, ranFixup, tc.wantFixup, action, err)
			}
			if !tc.wantFixup && action != ActionEscalated {
				t.Fatalf("cap=%d published=%d: action=%s want escalated", tc.cap, tc.published, action)
			}
		})
	}
}

// Review round-01 finding: strict-gate rounds that publish NO fix-up used to consume the
// budget, because the cycle number came from the review-round ordinal. The ratified unit
// is published cycles, so rounds that published nothing must not spend any.
func TestZeroFixRoundsDoNotSpendTheFixupBudget(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	// Four review rounds on disk, but NOTHING published yet — e.g. earlier rounds closed
	// with zero agreed fixes under strict_gate.
	writeImplWithCycles(t, ideaDir, "implemented", 0)
	for r := 1; r <= 4; r++ {
		os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(r)), 0o755)
	}
	// No `.fixup-done` anywhere: four rounds ran, none published a fix-up.
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{
		Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}
	d := New(Config{IdeaDir: ideaDir, IdeaSlug: "demo", Participants: parts, RunDir: runDir,
		Root: filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))), Events: store.New(runDir),
		Auto: true, AutoImplement: true, MaxFixupCycles: 2, Impl: fi}, &fakeRunner{})
	if _, _, err := d.Advance(context.Background()); err != nil && !contains(fi.calls, "fixup") {
		t.Fatalf("round 4 with cap 2 but zero published cycles must still run cycle 1: %v", err)
	}
	if !contains(fi.calls, "fixup") {
		t.Fatal("the budget was spent by rounds that published no fix-up")
	}
}

// Review round-03: every tamper direction on the fix-up budget must be fail-safe.
// Deleting driver state must not buy a cycle; forging it may only escalate sooner.
func TestFixupBudgetIsTamperFailSafe(t *testing.T) {
	setup := func(t *testing.T) (string, string, []string) {
		parts := []string{"codex", "agy"}
		ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
		writeFinalValid(t, ideaDir)
		writeImplWithCycles(t, ideaDir, "implemented", 2) // markers in rounds 1-2
		os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(3)), 0o755)
		os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
		return ideaDir, runDir, parts
	}
	newD := func(ideaDir, runDir string, parts []string, fi ImplOps, cap int) *Driver {
		return New(Config{IdeaDir: ideaDir, IdeaSlug: "demo", Participants: parts, RunDir: runDir,
			Root: filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))), Events: store.New(runDir),
			Auto: true, AutoImplement: true, MaxFixupCycles: cap, Impl: fi}, &fakeRunner{})
	}

	t.Run("deleting the markers does not lower the count when the cursor holds it", func(t *testing.T) {
		ideaDir, runDir, parts := setup(t)
		fi := &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{
			Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}
		d := newD(ideaDir, runDir, parts, fi, 2)
		// The driver's own count says 2 cycles are spent; an agent wipes the markers.
		for i := 1; i <= 2; i++ {
			os.Remove(filepath.Join(ideaDir, "review", roundLabel(i), ".fixup-done"))
		}
		c := Rebuild(ideaDir, 4)
		c.FixupCyclesPublished = 2
		if err := saveCursor(c, d.cursorPath()); err != nil {
			t.Fatal(err)
		}
		action, _, err := d.Advance(context.Background())
		if action != ActionEscalated || err == nil {
			t.Fatalf("wiping markers bought a cycle past the cap: action=%s err=%v", action, err)
		}
		if contains(fi.calls, "fixup") {
			t.Fatal("a fix-up ran past the cap after the markers were deleted")
		}
	})

	t.Run("a forged marker in the current round cannot drive AF2 past the cap", func(t *testing.T) {
		ideaDir, runDir, parts := setup(t)
		// Forge a marker in the CURRENT round: AF2 used to open the next round here
		// without ever consulting the budget.
		os.WriteFile(filepath.Join(ideaDir, "review", roundLabel(3), ".fixup-done"), []byte("x\n"), 0o644)
		fi := &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{
			Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}
		d := newD(ideaDir, runDir, parts, fi, 2)
		action, _, err := d.Advance(context.Background())
		if action != ActionEscalated || err == nil {
			t.Fatalf("AF2 opened another round past the cap: action=%s err=%v calls=%v", action, err, fi.calls)
		}
		if contains(fi.calls, "open-review") {
			t.Fatal("AF2 opened a review round with the budget exhausted")
		}
	})
}

// Review round-03: an oddly named directory must not be counted as a published cycle.
func TestOnlyExactRoundDirsCountAsPublishedCycles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"round-01", "round-02", "round-backup", "round-x", "rounds-03", "round-"} {
		rd := filepath.Join(dir, "review", name)
		os.MkdirAll(rd, 0o755)
		os.WriteFile(filepath.Join(rd, ".fixup-done"), []byte("x\n"), 0o644)
	}
	n, err := markedFixupCycles(dir)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("counted %d published cycles, want 2 — only round-01 and round-02 are real rounds", n)
	}
}

// Round-05 found that all three round-04 fixes could be reverted together without a
// single shipped test going red. These are the permanent equivalents.
func TestRound4FixesHavePermanentTests(t *testing.T) {
	base := func(t *testing.T, published int, fixupErr bool) (*Driver, *fakeImpl, string) {
		parts := []string{"codex", "agy"}
		ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
		writeFinalValid(t, ideaDir)
		writeImplWithCycles(t, ideaDir, "implemented", published)
		os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(published+1)), 0o755)
		os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
		fi := &fakeImpl{roundComplete: true, checksOK: true, fixupErr: fixupErr, review: ReviewStatus{
			Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}
		d := New(Config{IdeaDir: ideaDir, IdeaSlug: "demo", Participants: parts, RunDir: runDir,
			Root: filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))), Events: store.New(runDir),
			Auto: true, AutoImplement: true, MaxFixupCycles: 1, Impl: fi}, &fakeRunner{})
		return d, fi, ideaDir
	}

	// The reservation is taken BEFORE the code-writing call, so a fix-up that errors has
	// still spent its cycle. Saving after Fixup would hand the budget back on every error.
	t.Run("an errored fix-up keeps its reservation", func(t *testing.T) {
		d, fi, _ := base(t, 0, true)
		if _, _, err := d.Advance(context.Background()); err == nil {
			t.Fatal("an errored Fixup must escalate")
		}
		if !contains(fi.calls, "fixup") {
			t.Fatal("Fixup was never called")
		}
		c, err := LoadCursor(d.cursorPath())
		if err != nil {
			t.Fatalf("the reservation must be durable even when the fix-up failed: %v", err)
		}
		if c.FixupCyclesPublished != 1 {
			t.Fatalf("FixupCyclesPublished=%d after an errored cycle, want 1 (the cycle is spent when it runs)", c.FixupCyclesPublished)
		}
	})

	// AF2 completes the LAST allowed cycle at equality, and refuses beyond it.
	t.Run("AF2 finishes the last allowed cycle at equality", func(t *testing.T) {
		d, fi, ideaDir := base(t, 1, false)
		// The current round carries the marker: AF2's crash-recovery path.
		os.WriteFile(filepath.Join(ideaDir, "review", roundLabel(2), ".fixup-done"), []byte("x\n"), 0o644)
		// spent == 2 > cap 1 → refused. Now make spent == cap by using cap 2.
		d.cfg.MaxFixupCycles = 2
		if _, _, err := d.Advance(context.Background()); err != nil {
			t.Fatalf("AF2 at equality must be allowed to finish the spent cycle: %v", err)
		}
		if !contains(fi.calls, "open-review") {
			t.Fatal("AF2 did not complete the recovery")
		}
	})

	t.Run("AF2 refuses beyond the cap", func(t *testing.T) {
		d, fi, ideaDir := base(t, 1, false)
		os.WriteFile(filepath.Join(ideaDir, "review", roundLabel(2), ".fixup-done"), []byte("x\n"), 0o644)
		d.cfg.MaxFixupCycles = 1 // spent 2 > cap 1
		if _, _, err := d.Advance(context.Background()); err == nil {
			t.Fatal("AF2 past the cap must escalate")
		}
		if contains(fi.calls, "open-review") {
			t.Fatal("AF2 opened a round past the cap")
		}
	})

	// A cursor that exists but cannot be parsed is an UNKNOWN count and must escalate;
	// an absent cursor is a fresh run and must not.
	t.Run("a corrupt cursor escalates, an absent one does not", func(t *testing.T) {
		d, _, _ := base(t, 0, false)
		os.WriteFile(d.cursorPath(), []byte("{not json"), 0o644)
		if _, _, err := d.Advance(context.Background()); err == nil {
			t.Fatal("a corrupt cursor must escalate rather than be read as zero")
		}
		os.Remove(d.cursorPath())
		if _, _, err := d.Advance(context.Background()); err != nil {
			t.Fatalf("an absent cursor is a fresh run, not an error: %v", err)
		}
	})
}

// kimi-1 F5: a stat error on one round's marker must not read as "no marker" — that
// direction lowers the safety count. It must propagate so the caller escalates.
func TestUnreadableMarkerEscalatesInsteadOfLoweringTheCount(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "review", "round-01")
	bad := filepath.Join(dir, "review", "round-02")
	os.MkdirAll(good, 0o755)
	os.MkdirAll(bad, 0o755)
	os.WriteFile(filepath.Join(good, ".fixup-done"), []byte("x\n"), 0o644)
	os.WriteFile(filepath.Join(bad, ".fixup-done"), []byte("x\n"), 0o644)
	if err := os.Chmod(bad, 0o000); err != nil {
		t.Skipf("cannot make a directory unreadable here: %v", err)
	}
	t.Cleanup(func() { os.Chmod(bad, 0o755) })

	n, err := markedFixupCycles(dir)
	if err == nil {
		t.Fatalf("an unreadable round must escalate, got count=%d and no error", n)
	}
}

// Round-07: the escalation must report the count that was actually CHARGED, not the
// refused ordinal, and must not call a charged attempt a completed cycle. @codex-1's
// counterexample: one errored attempt, zero completed cycles, and the message said
// "after 1 cycle(s)".
func TestOverCapEscalationReportsChargedAttemptsNotCycles(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImpl(t, ideaDir, "implemented")
	os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(1)), 0o755)
	os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0o644)
	fi := &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{
		Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}
	d := New(Config{IdeaDir: ideaDir, IdeaSlug: "demo", Participants: parts, RunDir: runDir,
		Root: filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))), Events: store.New(runDir),
		Auto: true, AutoImplement: true, MaxFixupCycles: 1, Impl: fi}, &fakeRunner{})

	// One attempt charged, NO `.fixup-done` marker: the reservation-only state.
	c := Rebuild(ideaDir, 4)
	c.FixupCyclesPublished = 1
	if err := saveCursor(c, d.cursorPath()); err != nil {
		t.Fatal(err)
	}
	_, _, err := d.Advance(context.Background())
	if err == nil {
		t.Fatal("one charged attempt against cap 1 must refuse the next")
	}
	msg := err.Error()
	if !strings.Contains(msg, "after 1 charged attempt(s)") {
		t.Errorf("the escalation must report the CHARGED count (1), got: %s", msg)
	}
	if !strings.Contains(msg, "attempt 2 would exceed") {
		t.Errorf("the escalation must name the refused ordinal separately, got: %s", msg)
	}
	if strings.Contains(msg, "cycle(s)") {
		t.Errorf("a charged attempt is not a completed cycle; zero cycles completed here: %s", msg)
	}
}
