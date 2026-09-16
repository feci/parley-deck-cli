package driver

import (
	"context"

	"parley-deck-cli/internal/runner"
)

// The precheck seam (N1): a production adapter MAY expose a read-only
// "would this launch be refused before it starts?" check that mirrors the
// real launch's agent selection, idea, phase and artifact. The step/cycle
// wrappers in budget.go consult it BEFORE any ChargeStep/ChargeCycle, so a
// known invalid-protocol refusal costs no reservation. The interfaces the
// driver core already depends on are unchanged; adapters (and test fakes)
// that do not implement a seam keep today's charge-then-launch behavior.
// A precheck is never authority: the real launch re-renders at
// beginProtocolLaunch and may still refuse after charging (late failure
// stays spent).
//
// Fresh-scope order (N1 MAJOR-1): on a never-bound cross-review or fix-up
// scope, the legitimate cycle binding is admitted charge-free BEFORE the
// precheck runs (stepRoundRunner.RunRound via admitCrossReviewCycle; the
// fix-up branch in impl.go before reserveFixupCycle). A free precheck
// refusal then lands on an already-admitted scope instead of writing a
// refusal receipt that poisons the scope's first binding. Existing
// unmigrated history still refuses at admission — the hoist changes order,
// not verdicts.

// RoundLaunchPrechecker is the optional precheck for cross-review round
// launches (RoundRunner). round is the round the wrapper is about to open.
type RoundLaunchPrechecker interface {
	PrecheckRound(ctx context.Context, round int) error
}

// ConsensusLaunchPrechecker is the optional precheck for the agent-launching
// ConsensusOps methods (Status and Reopen stay pure and unprechecked).
type ConsensusLaunchPrechecker interface {
	PrecheckConsensusDraft(ctx context.Context) error
	PrecheckConsensusFinal(ctx context.Context) error
	PrecheckConsensusSignoffs(ctx context.Context, missing []string) error
}

// ImplLaunchPrechecker is the optional precheck for the agent-launching
// ImplOps methods beyond Fixup (which already has PrecheckFixup in the
// interface, invoked before reserveFixupCycle). RunChecks and Complete are
// deterministic local operations, not agent launches, and stay unprechecked.
type ImplLaunchPrechecker interface {
	PrecheckImplementation(ctx context.Context) error
	PrecheckReviewLaunch(ctx context.Context, round int) error
	PrecheckReviewConsensus(ctx context.Context, round int) error
	PrecheckReviewSignoffs(ctx context.Context, missing []string) error
	PrecheckGoalCheck(ctx context.Context) error
}

// PrecheckRound mirrors RunRound's own option normalization (round label,
// no-overwrite) and reuses the runner's read-only protocol precheck.
func (a roundRunnerAdapter) PrecheckRound(ctx context.Context, round int) error {
	opts := a.base
	opts.Round = round
	opts.RoundLabel = roundLabel(round)
	opts.Overwrite = false
	return runner.PrecheckRound(ctx, opts)
}

// precheckRoundLaunch resolves the precheck target through the private
// cycle-wrapping copy the step wrapper holds, so the check runs once at the
// outermost wrapper — before the step charge and the nested cycle charge.
func precheckRoundLaunch(ctx context.Context, r RoundRunner, round int) error {
	if c, ok := r.(cycleRoundRunner); ok {
		r = c.RoundRunner
	}
	if p, ok := r.(RoundLaunchPrechecker); ok {
		return p.PrecheckRound(ctx, round)
	}
	return nil
}

func precheckConsensusDraft(ctx context.Context, ops ConsensusOps) error {
	if p, ok := ops.(ConsensusLaunchPrechecker); ok {
		return p.PrecheckConsensusDraft(ctx)
	}
	return nil
}

func precheckConsensusFinal(ctx context.Context, ops ConsensusOps) error {
	if p, ok := ops.(ConsensusLaunchPrechecker); ok {
		return p.PrecheckConsensusFinal(ctx)
	}
	return nil
}

func precheckConsensusSignoffs(ctx context.Context, ops ConsensusOps, missing []string) error {
	if p, ok := ops.(ConsensusLaunchPrechecker); ok {
		return p.PrecheckConsensusSignoffs(ctx, missing)
	}
	return nil
}

func precheckImplementation(ctx context.Context, ops ImplOps) error {
	if p, ok := ops.(ImplLaunchPrechecker); ok {
		return p.PrecheckImplementation(ctx)
	}
	return nil
}

func precheckReviewLaunch(ctx context.Context, ops ImplOps, round int) error {
	if p, ok := ops.(ImplLaunchPrechecker); ok {
		return p.PrecheckReviewLaunch(ctx, round)
	}
	return nil
}

func precheckReviewConsensus(ctx context.Context, ops ImplOps, round int) error {
	if p, ok := ops.(ImplLaunchPrechecker); ok {
		return p.PrecheckReviewConsensus(ctx, round)
	}
	return nil
}

func precheckReviewSignoffs(ctx context.Context, ops ImplOps, missing []string) error {
	if p, ok := ops.(ImplLaunchPrechecker); ok {
		return p.PrecheckReviewSignoffs(ctx, missing)
	}
	return nil
}

func precheckGoalCheck(ctx context.Context, ops ImplOps) error {
	if p, ok := ops.(ImplLaunchPrechecker); ok {
		return p.PrecheckGoalCheck(ctx)
	}
	return nil
}
