package driver

import (
	"context"
	"fmt"

	"parley-deck-cli/internal/budget"
)

type cycleRoundRunner struct {
	RoundRunner
	cfg Config
}

func (r cycleRoundRunner) RunRound(ctx context.Context, round int) error {
	if round <= 1 {
		return r.RoundRunner.RunRound(ctx, round)
	}
	cap := r.cfg.MaxRounds
	if r.cfg.HardCrossReviewCap > 0 && r.cfg.HardCrossReviewCap < cap {
		cap = r.cfg.HardCrossReviewCap
	}
	if r.cfg.Track == "fast" {
		cap = 0
	}
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
