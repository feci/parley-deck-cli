package runner

import (
	"context"
	"os"
	"path/filepath"
)

// PrecheckRound mirrors RunRound's normalization and checks the first pending
// participant without opening a step/cycle session or writing a segment.
func PrecheckRound(ctx context.Context, opts Options) error {
	if opts.Round < 2 {
		opts.Round = 2
	}
	if opts.RoundLabel == "" {
		opts.RoundLabel = roundLabel(opts.Round)
	}
	return precheckSelectedLaunch(ctx, opts, false)
}

// PrecheckImplementation mirrors the first implementer's actual artifact skip.
func PrecheckImplementation(ctx context.Context, opts Options) error {
	opts.Phase = "implementation"
	opts.ArtifactName = "IMPLEMENTATION.md"
	if opts.RoundLabel == "" {
		opts.RoundLabel = "implementation"
	}
	return precheckSelectedLaunch(ctx, opts, true)
}

// PrecheckReviewRound mirrors the reviewer selection and output directory.
func PrecheckReviewRound(ctx context.Context, opts Options) error {
	if opts.Round < 1 {
		opts.Round = 1
	}
	opts.Phase = "review"
	opts.RoundLabel = filepath.Join("review", roundLabel(opts.Round))
	return precheckSelectedLaunch(ctx, opts, false)
}

// Review consensus always rewrites its prior draft, even when Overwrite was
// false on entry, exactly as RunReviewConsensus does.
func PrecheckReviewConsensus(ctx context.Context, opts Options) error {
	opts.Phase = "review-consensus"
	opts.ArtifactName = filepath.Join("review", "consensus.md")
	opts.Overwrite = true
	return precheckSelectedLaunch(ctx, opts, true)
}

// The renderer depends on phase/idea/root, so one pending participant suffices
// for this known-protocol-refusal check. Its retained refusal must name an agent
// that would actually reach a launch; an already-owned artifact is not a launch.
// Success grants no authority: every real launch still renders again.
func precheckSelectedLaunch(ctx context.Context, opts Options, firstOnly bool) error {
	selected, _ := selectedAgents(opts.Idea.Participants, opts.Agents, resolveMapping(opts))
	if firstOnly && len(selected) > 1 {
		selected = selected[:1]
	}
	for _, agent := range selected {
		output := filepath.Join(opts.Idea.Path, opts.RoundLabel, agent.ID+".md")
		if opts.ArtifactName != "" {
			output = filepath.Join(opts.Idea.Path, opts.ArtifactName)
		}
		// Match runAgent's exact skip predicate, including stat errors proceeding
		// to the real launch checks rather than becoming an inferred skip.
		if _, err := os.Stat(output); err == nil && !opts.Overwrite {
			continue
		}
		return PrecheckProtocolLaunch(ctx, opts.Root, agent, LaunchInfo{
			RunID: opts.RunID, Idea: opts.Idea.Slug, Phase: protocolLaunchPhase(opts),
			AttemptOrdinal: 1, Store: opts.Store, ArtifactPath: output,
		})
	}
	return nil
}
