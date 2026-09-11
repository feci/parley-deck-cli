package driver

import (
	"context"
	"errors"
	"parley-deck-cli/internal/budget"
)

// Each Advance uses a private adapter copy. Concurrent callers never mutate the
// shared Driver configuration while nested operations reuse one step session.
func (d *Driver) withStepBudget(ctx context.Context) (context.Context, *Driver, func(), error) {
	if d.cfg.MaxCostUSD != 0 {
		monetary, err := budget.LoadLaunchBinding(ctx, d.cfg.Root, d.cfg.IdeaSlug)
		if err == nil {
			err = budget.RequireMonetaryBinding(monetary, d.cfg.MaxCostUSD)
		}
		if err != nil {
			return ctx, d, func() {}, err
		}
	}
	b, err := budget.EnsureStepBinding(ctx, d.cfg.Root, d.cfg.IdeaSlug, d.cfg.MaxDriverSteps, d.cfg.MaxWallClock)
	if err != nil {
		return ctx, d, func() {}, err
	}
	ctx, finish, err := budget.OpenStepSession(ctx, b)
	if err != nil {
		return ctx, d, finish, err
	}
	copy := *d
	if d.runner != nil {
		copy.runner = stepRoundRunner{cycleRoundRunner{d.runner, d.cfg}}
	}
	if d.cfg.Consensus != nil {
		copy.cfg.Consensus = stepConsensusOps{d.cfg.Consensus}
	}
	if d.cfg.Impl != nil {
		copy.cfg.Impl = stepImplOps{d.cfg.Impl}
	}
	return ctx, &copy, finish, nil
}

type stepRoundRunner struct{ RoundRunner }

func (o stepRoundRunner) RunRound(ctx context.Context, n int) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.RoundRunner.RunRound(ctx, n)
}

type stepConsensusOps struct{ ConsensusOps }

func (o stepConsensusOps) Draft(ctx context.Context) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ConsensusOps.Draft(ctx)
}
func (o stepConsensusOps) DraftFinal(ctx context.Context) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ConsensusOps.DraftFinal(ctx)
}
func (o stepConsensusOps) RequestSignoffs(ctx context.Context, missing []string) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ConsensusOps.RequestSignoffs(ctx, missing)
}
func (o stepConsensusOps) Reopen(ctx context.Context, reason string) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ConsensusOps.Reopen(ctx, reason)
}

type stepImplOps struct{ ImplOps }

func (o stepImplOps) Implement(ctx context.Context) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ImplOps.Implement(ctx)
}
func (o stepImplOps) OpenReviewRound(ctx context.Context, n int) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ImplOps.OpenReviewRound(ctx, n)
}
func (o stepImplOps) DraftReviewConsensus(ctx context.Context, n int) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ImplOps.DraftReviewConsensus(ctx, n)
}
func (o stepImplOps) RequestReviewSignoffs(ctx context.Context, missing []string) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ImplOps.RequestReviewSignoffs(ctx, missing)
}
func (o stepImplOps) Fixup(ctx context.Context, n int) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ImplOps.Fixup(ctx, n)
}
func (o stepImplOps) RunChecks(ctx context.Context) (bool, string) {
	if err := budget.ChargeStep(ctx); err != nil {
		return false, err.Error()
	}
	return o.ImplOps.RunChecks(ctx)
}
func (o stepImplOps) GoalCheck(ctx context.Context) (bool, string) {
	if err := budget.ChargeStep(ctx); err != nil {
		return false, err.Error()
	}
	return o.ImplOps.GoalCheck(ctx)
}
func (o stepImplOps) VerifyCompletionEvidence(ctx context.Context) (bool, string) {
	verifier, ok := o.ImplOps.(CompletionEvidenceOps)
	if !ok {
		return false, errors.New("adapter does not provide independent verifier execution").Error()
	}
	if err := budget.ChargeStep(ctx); err != nil {
		return false, err.Error()
	}
	return verifier.VerifyCompletionEvidence(ctx)
}
func (o stepImplOps) Complete(ctx context.Context) error {
	if err := budget.ChargeStep(ctx); err != nil {
		return err
	}
	return o.ImplOps.Complete(ctx)
}
