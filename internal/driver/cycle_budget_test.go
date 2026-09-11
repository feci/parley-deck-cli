package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

func TestCycleFixupFailedAttemptSurvivesNewRunAndCursorDeletion(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImplWithCycles(t, ideaDir, "implemented", 0)
	if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-01"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	fi := &fakeImpl{roundComplete: true, checksOK: true, fixupErr: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	makeDriver := func(dir string) *Driver {
		return New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: dir, Participants: parts, Events: store.New(dir), Auto: true, AutoImplement: true, MaxFixupCycles: 1, Impl: fi}, &fakeRunner{})
	}
	d := makeDriver(runDir)
	if _, _, err := d.Advance(context.Background()); err == nil || !contains(fi.calls, "fixup") {
		t.Fatalf("first failed attempt: %v %v", err, fi.calls)
	}
	if err := os.Remove(d.cursorPath()); err != nil {
		t.Fatal(err)
	}
	fi.calls = nil
	if _, _, err := makeDriver(filepath.Join(filepath.Dir(runDir), "new-run")).Advance(context.Background()); err == nil || contains(fi.calls, "fixup") {
		t.Fatalf("new run refunded the failed fixup: %v %v", err, fi.calls)
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, "demo", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	state, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(state) != 1 {
		t.Fatalf("shared fixup count: %+v %v", state, err)
	}
}

type cycleProcessImpl struct {
	*fakeImpl
	opts runner.Options
}

func (f *cycleProcessImpl) Fixup(ctx context.Context, cycle int) error {
	f.calls = append(f.calls, "fixup")
	r := runner.RunFixup(ctx, f.opts)
	if r.ExitError == "" {
		return errors.New("expected local child fixture to fail")
	}
	return fmt.Errorf("actual fixup child: %s", r.ExitError)
}

func TestCycleDriverAndManualRunnerShareExactlyOneFailedCharge(t *testing.T) {
	parts := []string{"builder", "reviewer"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImplWithCycles(t, ideaDir, "implemented", 0)
	if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-01"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	// Declare the test workspace as a live protocol source for the actual child
	// context renderer, just as the runner's process fixtures do.
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "parley-deck", "meta"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "meta", "version.json"), []byte(`{"protocolRole":"source"}`), 0600); err != nil {
		t.Fatal(err)
	}
	agent := agents.Discovery{Spec: agents.Spec{ID: "builder", HeadlessArgs: []string{"-test.run=TestCycleDriverChildHelper", "--", "parley-cycle-child"}, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true}
	opts := runner.Options{Root: root, RunID: filepath.Base(runDir), Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: []string{"builder"}}, Agents: []agents.Discovery{agent}, Timeout: time.Second, Store: store.New(runDir)}
	fi := &cycleProcessImpl{fakeImpl: &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}, opts: opts}
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), Auto: true, AutoImplement: true, MaxFixupCycles: 1, Impl: fi}, &fakeRunner{})
	if _, _, err := d.Advance(context.Background()); err == nil {
		t.Fatal("failed child was accepted")
	} else {
		t.Logf("first actual driver attempt: %v", err)
	}
	opts.RunID = "manual-after-driver"
	opts.Store = store.New(filepath.Join(filepath.Dir(runDir), opts.RunID))
	if r := runner.RunFixup(context.Background(), opts); r.ExitError == "" {
		t.Fatal("manual runner reset driver's failed charge")
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, "demo", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	state, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(state) != 1 {
		t.Fatalf("nested charge was duplicated or refunded: %+v %v", state, err)
	}
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	started, denied := 0, 0
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var r telemetry.Record
		if err := json.Unmarshal(data, &r); err != nil {
			t.Fatal(err)
		}
		if r.StartedAt != nil {
			started++
		}
		if r.Outcome.FailureClass != nil && *r.Outcome.FailureClass == "budget_refused" {
			denied++
		}
	}
	if len(paths) != 2 || started != 1 || denied != 1 {
		t.Fatalf("actual driver/manual boundary: terminals=%d started=%d refused=%d", len(paths), started, denied)
	}
}

func TestCycleDriverChildHelper(t *testing.T) {
	for _, arg := range os.Args {
		if arg == "parley-cycle-child" {
			os.Exit(7)
		}
	}
}

func TestCycleFailedCrossReviewsSpendTheHardCap(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{err: errors.New("failed before round publication")}
	d := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
	d.cfg.HardCrossReviewCap = 3
	for i := 0; i < 4; i++ {
		if _, _, err := d.Advance(context.Background()); err == nil {
			t.Fatal("failed attempt passed")
		}
	}
	if len(fr.calls) != 3 {
		t.Fatalf("hard cap granted %d failed attempts, want 3", len(fr.calls))
	}
}

func TestCycleBlockedReopenUsesOrdinaryCrossReviewCharges(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{err: errors.New("failed round dispatch")}
	d := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
	d.cfg.HardCrossReviewCap = 3
	if _, _, err := d.Advance(context.Background()); err == nil {
		t.Fatal("first failure passed")
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "consensus.md"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	d.cfg.Consensus = &fakeConsensus{statusSeq: []string{consensus.TriageBlocked}}
	for i := 0; i < 3; i++ {
		if _, _, err := d.Advance(context.Background()); err == nil {
			t.Fatal("blocked failure passed")
		}
	}
	if len(fr.calls) != 3 {
		t.Fatalf("BLOCK reset the ordinary round budget: %v", fr.calls)
	}
	if _, err := os.Stat(filepath.Join(ideaDir, "consensus.md")); err != nil {
		t.Fatalf("failed BLOCK backedge lost consensus: %v", err)
	}
}

func TestCycleFinalAllowedFixupCanCloseButOverCapCannot(t *testing.T) {
	for _, tc := range []struct{ spent, frozen int }{{5, 5}, {6, 5}, {5, 4}} {
		spent := tc.spent
		t.Run(fmt.Sprintf("spent-%d/frozen-%d", spent, tc.frozen), func(t *testing.T) {
			parts := []string{"builder", "reviewer"}
			ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
			writeFinalValid(t, ideaDir)
			writeImplWithCycles(t, ideaDir, "implemented", spent)
			if err := os.MkdirAll(filepath.Join(ideaDir, "review", roundLabel(spent+1)), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0644); err != nil {
				t.Fatal(err)
			}
			fi := &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, ReviewerCount: 2}}
			d := New(Config{Root: filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))), IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), Auto: true, AutoImplement: true, MaxFixupCycles: 5, Impl: fi}, &fakeRunner{})
			if _, err := budget.EnsureCycleBinding(context.Background(), d.cfg.Root, "demo", budget.Fixup, tc.frozen, spent, runDir, ideaDir); err != nil {
				t.Fatal(err)
			}
			a, _, err := d.Advance(context.Background())
			allowed := spent == 5 && tc.frozen == 5
			if allowed && (err != nil || a != ActionComplete || !contains(fi.calls, "complete")) {
				t.Fatalf("final allowed attempt stranded: %s %v %v", a, err, fi.calls)
			}
			if !allowed && (err == nil || a != ActionEscalated || contains(fi.calls, "complete")) {
				t.Fatalf("over-cap state closed: %s %v %v", a, err, fi.calls)
			}
		})
	}
}
