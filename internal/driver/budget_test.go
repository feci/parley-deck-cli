package driver

import (
	"context"
	"errors"
	"parley-deck-cli/internal/budget"
	"path/filepath"
	"testing"
	"time"
)

func TestDriverConsumesRecordedStepAndMonetaryExtensions(t *testing.T) {
	ctx := context.Background()
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{writeOnRun: func(n int) { writeAll(t, ideaDir, n, parts) }}
	d := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
	d.cfg.MaxDriverSteps = 1
	if _, _, err := d.Advance(ctx); err != nil {
		t.Fatal(err)
	}
	s, err := budget.InspectRuntimeBudget(ctx, d.cfg.Root, "demo", budget.DriverStep)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := budget.ExtendRuntimeBudget(ctx, d.cfg.Root, "demo", budget.DriverStep, budget.PolicyExtensionRequest{DecisionID: "step", ExpectedPolicySHA256: s.PolicySHA256, Reason: "Explicit fixture step extension", Ceilings: budget.PolicyCeilings{Actions: 2}}); err != nil {
		t.Fatal(err)
	}
	// Run reconstructs state while still carrying the original one-step config.
	if err := d.Run(ctx); err != nil || len(fr.calls) != 2 {
		t.Fatalf("recorded step extension ignored: calls=%v err=%v", fr.calls, err)
	}
	// Monetary loop enforcement must use current exposure and the recorded grant,
	// even when runtime configuration still names the original dollar ceiling.
	root := t.TempDir()
	reserve := int64(1000000)
	b, err := budget.ConfigureLaunchBudget(ctx, root, "money", budget.LaunchPolicy{MaxCostMicros: reserve, ReserveMicros: &reserve})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Reserve(ctx, "first"); err != nil {
		t.Fatal(err)
	}
	monetary := New(Config{Root: root, IdeaSlug: "money", MaxCostUSD: 1}, &fakeRunner{})
	if reason, _ := monetary.runtimeLoopBudgetBreach(ctx, 0, time.Now()); reason == "" {
		t.Fatal("spent original ceiling passed")
	}
	s, err = budget.InspectRuntimeBudget(ctx, root, "money", budget.Launch)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := budget.ExtendRuntimeBudget(ctx, root, "money", budget.Launch, budget.PolicyExtensionRequest{DecisionID: "money", ExpectedPolicySHA256: s.PolicySHA256, Reason: "Explicit fixture cost extension", Ceilings: budget.PolicyCeilings{CostMicros: 2 * reserve}}); err != nil {
		t.Fatal(err)
	}
	if reason, limit := monetary.runtimeLoopBudgetBreach(ctx, 0, time.Now()); reason != "" || limit == nil || *limit != 2*reserve {
		t.Fatalf("original config ignored monetary grant: %s %v", reason, limit)
	}
}

func TestDriverMonetaryLimitRequiresPersistentReservationPolicy(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{}
	d := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
	d.cfg.MaxCostUSD = 1
	action, _, err := d.Advance(context.Background())
	if !errors.Is(err, budget.ErrUnknownCost) || action != ActionEscalated || len(fr.calls) != 0 {
		t.Fatalf("unreserved monetary dispatch: %s %v calls=%v", action, err, fr.calls)
	}
}

func TestDriverLifetimeStepCapSurvivesRunAndDirectResume(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{writeOnRun: func(round int) { writeAll(t, ideaDir, round, parts) }}
	d := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
	d.cfg.MaxDriverSteps = 1
	if err := d.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(fr.calls) != 1 {
		t.Fatalf("first run: %v", fr.calls)
	}
	// Re-entry omits the old flags and uses another run directory.
	resumed := newTestDriver(ideaDir, filepath.Join(filepath.Dir(runDir), "second-run"), parts, 3, true, fr)
	action, _, err := resumed.Advance(context.Background())
	if !errors.Is(err, budget.ErrLimit) || action != ActionEscalated || len(fr.calls) != 1 {
		t.Fatalf("resume granted work: %s %v %v", action, err, fr.calls)
	}
	if err := d.Run(context.Background()); err != nil || len(fr.calls) != 1 {
		t.Fatalf("same run resume: %v calls=%v", err, fr.calls)
	}
}

func TestDriverResumedRunReportsLifetimeStepsAndOrigin(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{writeOnRun: func(n int) { writeAll(t, ideaDir, n, parts) }}
	d := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
	d.cfg.MaxDriverSteps = 2
	if _, _, err := d.Advance(context.Background()); err != nil {
		t.Fatal(err)
	}
	b, err := budget.LoadStepBinding(context.Background(), d.cfg.Root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	before, err := b.Store.Inspect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resumed := newTestDriver(ideaDir, filepath.Join(filepath.Dir(runDir), "resumed"), parts, 3, true, fr)
	if err := resumed.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	events, err := resumed.cfg.Events.Load()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range events {
		if e.Type != "loop.budget" {
			continue
		}
		found = true
		if e.Data["steps"] != float64(2) || e.Data["max_driver_steps"] != float64(2) {
			t.Fatalf("resumed budget report reset: %+v", e.Data)
		}
		elapsed, ok := e.Data["elapsed_ms"].(float64)
		if !ok || time.Duration(elapsed)*time.Millisecond < e.Time.Sub(before.StartedAt)-2*time.Millisecond {
			t.Fatalf("resumed clock origin reset: %+v", e)
		}
	}
	if !found || len(fr.calls) != 2 {
		t.Fatalf("found=%v calls=%v", found, fr.calls)
	}
}
func TestDriverAwaitDoesNotChargeAndFailedRoundStaysSpent(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	fr := &fakeRunner{err: errors.New("actual dispatch failure")}
	d := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
	d.cfg.MaxDriverSteps = 1
	for i := 0; i < 2; i++ {
		a, _, err := d.Advance(context.Background())
		if err != nil || a != ActionAwait {
			t.Fatalf("await: %s %v", a, err)
		}
	}
	b, err := budget.LoadStepBinding(context.Background(), d.cfg.Root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Store.Inspect(context.Background())
	if err != nil || len(s.Entries) != 0 {
		t.Fatalf("await charged: %+v %v", s, err)
	}
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	if _, _, err = d.Advance(context.Background()); err == nil {
		t.Fatal("failed round passed")
	}
	fr.err = nil
	if _, _, err = d.Advance(context.Background()); !errors.Is(err, budget.ErrLimit) || len(fr.calls) != 1 {
		t.Fatalf("failed step refunded: %v %v", err, fr.calls)
	}
}
