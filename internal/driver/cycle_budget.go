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
	if err := budget.ChargeStep(ctx); err != nil {
		return ctx, 0, finish, err
	}
	n, err := budget.ChargeCycle(ctx, budget.Fixup)
	return ctx, n, finish, err
}
