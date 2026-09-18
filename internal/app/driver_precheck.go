package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/runner"
)

// Production precheck implementations for the driver's optional launch-precheck
// seams (N1). Each method mirrors its real launch's agent selection, idea,
// phase and artifact so a known invalid-protocol launch is refused before the
// driver spends a step or cycle reservation. A precheck never authorizes
// anything: the real launch re-renders the protocol at beginProtocolLaunch.

func (o driverImplOps) PrecheckImplementation(ctx context.Context) error {
	return runner.PrecheckImplementation(ctx, o.withParticipants(o.implementer))
}

func (o driverImplOps) PrecheckReviewLaunch(ctx context.Context, round int) error {
	opts := o.withParticipants(o.reviewers...)
	opts.Round = round
	return runner.PrecheckReviewRound(ctx, opts)
}

func (o driverImplOps) PrecheckReviewConsensus(ctx context.Context, round int) error {
	opts := o.withParticipants(o.drafter)
	opts.Round = round
	return runner.PrecheckReviewConsensus(ctx, opts)
}

func (o driverImplOps) PrecheckReviewSignoffs(ctx context.Context, missing []string) error {
	return precheckSignoffLaunch(ctx, o.root, o.ideaSlug, true, missing)
}

// PrecheckGoalCheck mirrors GoalCheck's checker resolution. The no-checker and
// unresolved-checker cases never launch an agent, so they are not known
// refusals and stay nil — the real GoalCheck fails on them as before.
func (o driverImplOps) PrecheckGoalCheck(ctx context.Context) error {
	checker := o.drafter
	if checker == "" || checker == o.implementer {
		return nil
	}
	agent, err := agents.ResolveParticipant(checker, o.base.Agents, rosterMappingFor(o.root))
	if err != nil {
		return nil
	}
	return runner.PrecheckProtocolLaunch(ctx, o.root, agent, runner.LaunchInfo{
		RunID: o.base.RunID, Idea: o.ideaSlug, Phase: "goal-check", Store: o.base.Store,
	})
}

func (o driverConsensusOps) PrecheckConsensusDraft(ctx context.Context) error {
	return precheckDrafterLaunch(ctx, o)
}

func (o driverConsensusOps) PrecheckConsensusFinal(ctx context.Context) error {
	return precheckDrafterLaunch(ctx, o)
}

func (o driverConsensusOps) PrecheckConsensusSignoffs(ctx context.Context, missing []string) error {
	return precheckSignoffLaunch(ctx, o.root, o.ideaSlug, false, missing)
}

// precheckDrafterLaunch mirrors runDrafter's selection (firstHeadlessAgent over
// the idea participants). The real drafter launch goes through
// runHeadlessSignoffAgent → CommandFor with no LaunchInfo on the context, so
// the matching precheck renders with the same empty phase/idea and the
// "one-shot" run identity that boundary assigns.
func precheckDrafterLaunch(ctx context.Context, o driverConsensusOps) error {
	drafter, ok := firstHeadlessAgent(o.discovered, o.participants, rosterMappingFor(o.root))
	if !ok {
		return nil // the real Draft/Final errors before launching; not a known protocol refusal
	}
	return runner.PrecheckProtocolLaunch(ctx, o.root, drafter, runner.LaunchInfo{RunID: "one-shot"})
}

// precheckSignoffLaunch mirrors requestConsensusSignoffs' triage gates and
// target and agent resolution. The malformed and blocked-consensus aborts run
// in the production order with the production reasons BEFORE any discovery,
// render or telemetry, so a blocked consensus can never record a
// protocol_context_refused terminal for a launch the real path would not
// attempt. Only the headless signoff launches render the protocol, so the
// precheck names the first headless selected agent with the same consensus
// phase and artifact the real runSignoffAgent would carry; interactive and
// manual handoffs launch no protocol render. Resolution failures here
// reproduce deterministically in the real call and are returned so they cost
// no reservation either.
func precheckSignoffLaunch(ctx context.Context, root, ideaSlug string, review bool, missing []string) error {
	summary, err := consensus.Status(root, ideaSlug, review)
	if err != nil {
		return err
	}
	if len(summary.Errors) > 0 {
		return fmt.Errorf("target consensus is malformed: %s", strings.Join(summary.Errors, "; "))
	}
	if summary.Triage == consensus.TriageBlocked {
		return errors.New("target consensus is blocked; resolve the BLOCK before requesting more signoffs")
	}
	targets, _, err := requestSignoffTargets(summary, strings.Join(missing, ","))
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil
	}
	discovered, err := discoverConfigured(ctx, root)
	if err != nil {
		return err
	}
	selected, err := requestSignoffAgents(targets, discovered, rosterMappingFor(root))
	if err != nil {
		return err
	}
	phase := "consensus"
	if review {
		phase = "review-consensus"
	}
	for _, agent := range selected {
		if agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchHeadless {
			continue
		}
		return runner.PrecheckProtocolLaunch(ctx, root, agent, runner.LaunchInfo{
			Idea: ideaSlug, Phase: phase, ArtifactPath: summary.Path,
		})
	}
	return nil
}
