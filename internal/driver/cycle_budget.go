package driver

import (
	"context"
	"fmt"
	"path/filepath"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/trajectory"
)

type cycleRoundRunner struct {
	RoundRunner
	cfg Config
}

func crossReviewBase(cfg Config) int {
	cap := cfg.MaxRounds
	if cfg.HardCrossReviewCap > 0 && cfg.HardCrossReviewCap < cap {
		cap = cfg.HardCrossReviewCap
	}
	if cfg.Track == "fast" {
		cap = 0
	}
	return cap
}

// A recorded finite grant changes only this cycle ceiling. The original track
// configuration, skipped phases and independent review gates remain intact.
func (d *Driver) cycleMaximum(ctx context.Context, kind budget.Kind, base int) (int, error) {
	b, err := budget.LoadCycleBinding(ctx, d.cfg.Root, d.cfg.IdeaSlug, kind)
	if err != nil {
		return 0, err
	}
	if b == nil {
		return base, nil
	}
	rel, err := filepath.Rel(d.cfg.Root, d.cfg.IdeaDir)
	if err != nil || filepath.ToSlash(rel) != b.Policy.IdeaPath || base != b.Policy.InitialMaximum() {
		return 0, fmt.Errorf("%s original cycle policy differs from the current configuration", kind)
	}
	if _, err := b.Inspect(ctx); err != nil {
		return 0, err
	}
	return b.Policy.Maximum, nil
}

// admitCrossReviewCycle performs the legitimate, charge-free first binding of
// the cross-review scope BEFORE the step wrapper's protocol precheck (N1
// MAJOR-1): on a fresh, never-bound scope a free precheck refusal would
// otherwise write a refusal receipt that poisons the scope's first binding.
// Existing unmigrated history still refuses here — admission is hoisted, not
// forgiven: no receipt is ignored, nothing is auto-migrated, and neither the
// cap nor any classification changes. It mirrors the exact parameters
// cycleRoundRunner.RunRound binds with; the nested RunRound re-ensures the
// same binding and charges it as before.
func admitCrossReviewCycle(ctx context.Context, r RoundRunner, round int) error {
	c, ok := r.(cycleRoundRunner)
	if !ok || round <= 1 {
		return nil
	}
	cap := crossReviewBase(c.cfg)
	floor, err := budget.LegacyCycleFloor(c.cfg.IdeaDir, budget.CrossReview)
	if err != nil {
		return err
	}
	if _, err := budget.EnsureCycleBinding(ctx, c.cfg.Root, c.cfg.IdeaSlug, budget.CrossReview, cap, floor, c.cfg.RunDir, c.cfg.IdeaDir); err != nil {
		return fmt.Errorf("cross-review accounting: %w", err)
	}
	return nil
}

// preflightRunnerCycle is the MINOR-1 read-only known-exhaustion check for the
// cross-review cycle this RunRound is about to charge INSIDE the step wrapper:
// it runs after admission (which binds the scope charge-free) and before
// ChargeStep, so a cycle already at its frozen maximum refuses the transition
// without spending a step. It mirrors the exact binding cycleRoundRunner
// charges; reuse of an already charged session and late races keep
// ChargeCycle's exact semantics — nothing is migrated, refunded or re-capped.
func preflightRunnerCycle(ctx context.Context, r RoundRunner, round int) error {
	c, ok := r.(cycleRoundRunner)
	if !ok || round <= 1 {
		return nil
	}
	b, err := budget.LoadCycleBinding(ctx, c.cfg.Root, c.cfg.IdeaSlug, budget.CrossReview)
	if err != nil {
		return fmt.Errorf("cross-review accounting: %w", err)
	}
	if err := budget.PreflightCycleCharge(ctx, b); err != nil {
		return fmt.Errorf("cross-review budget: %w", err)
	}
	return nil
}

func (r cycleRoundRunner) RunRound(ctx context.Context, round int) error {
	if round <= 1 {
		return r.RoundRunner.RunRound(ctx, round)
	}
	cap := crossReviewBase(r.cfg)
	floor, err := budget.LegacyCycleFloor(r.cfg.IdeaDir, budget.CrossReview)
	if err != nil {
		return err
	}
	b, err := budget.EnsureCycleBinding(ctx, r.cfg.Root, r.cfg.IdeaSlug, budget.CrossReview, cap, floor, r.cfg.RunDir, r.cfg.IdeaDir)
	if err != nil {
		return fmt.Errorf("cross-review accounting: %w", err)
	}
	ctx, finish, err := budget.OpenCycleSession(ctx, b)
	if err != nil {
		return err
	}
	defer finish()
	if _, err := budget.ChargeCycle(ctx, budget.CrossReview); err != nil {
		return fmt.Errorf("cross-review budget: %w", err)
	}
	return r.RoundRunner.RunRound(ctx, round)
}

// Reserve before the compatibility cursor is updated. Its failed publication
// cannot refund the durable charge; the nested real runner reuses this session.
func (d *Driver) reserveFixupCycle(ctx context.Context, charged int) (context.Context, int, func(), error) {
	ctx = budget.WithCycleObserver(ctx, &trajectory.Observer{Root: d.cfg.Root})
	b, err := budget.EnsureCycleBinding(ctx, d.cfg.Root, d.cfg.IdeaSlug, budget.Fixup, d.cfg.MaxFixupCycles, charged, d.cfg.RunDir, d.cfg.IdeaDir)
	if err != nil {
		return ctx, 0, func() {}, err
	}
	ctx, finish, err := budget.OpenCycleSession(ctx, b)
	if err != nil {
		return ctx, 0, finish, err
	}
	// MINOR-1: a KNOWN exhausted fix-up cycle refuses for free BEFORE the step
	// charge. An already charged live session stays reusable at the cap; a
	// refusal the read-only preflight cannot know (a late race) still surfaces
	// at ChargeCycle below and stays spent.
	if err := budget.PreflightCycleCharge(ctx, b); err != nil {
		return ctx, 0, finish, err
	}
	if err := budget.ChargeStep(ctx); err != nil {
		return ctx, 0, finish, err
	}
	n, err := budget.ChargeCycle(ctx, budget.Fixup)
	return ctx, n, finish, err
}
