package runner

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
)

func quotaBefore(ctx context.Context, opts Options) (context.Context, Options, func(), error) {
	ctx, release, err := membership.Acquire(ctx, opts.Idea.Path, opts.RunID)
	if err != nil {
		return ctx, opts, nil, err
	}
	h, err := membership.Before(ctx, opts.Root, opts.Idea.Path, opts.RunID)
	if err != nil {
		release()
		return ctx, opts, nil, err
	}
	if h != nil {
		opts.quotaMidIdea = h.MidIdea()
		opts.participantDropout = h.MidIdea() && h.Policy().Dropout()
		opts.Idea.Participants = membership.Intersect(opts.Idea.Participants, h.Current)
		if len(opts.Idea.Participants) == 0 {
			release()
			return ctx, opts, nil, fmt.Errorf("no current quota participants in requested dispatch")
		}
	}
	return ctx, opts, release, nil
}

func quotaSettle(ctx context.Context, opts Options, results []Result) []Result {
	h, err := quota.ReadHistory(opts.Idea.Path)
	if err != nil {
		return append(results, Result{AgentID: "runner/quota", ExitError: err.Error(), QuotaBlocked: membership.IsBlocked(err)})
	}
	if !h.MidIdea() {
		return results
	}
	members := []quota.Member{}
	for _, r := range results {
		members = append(members, quota.Member{ID: r.AgentID, Usable: r.Success(), ValidArtifact: r.ArtifactOK, Evidence: r.QuotaEvidence})
	}
	// A review batch does not launch the pinned implementer; a valid completed
	// implementation is its usability evidence, not a new provider observation.
	if opts.Phase == "review" {
		meta, _ := protocol.ReadFrontmatter(filepath.Join(opts.Idea.Path, "IMPLEMENTATION.md"))
		impl := meta["implementer"]
		if membership.Has(h.Current, impl) && !membership.Has(opts.Idea.Participants, impl) && ValidateImplementationArtifact(filepath.Join(opts.Idea.Path, "IMPLEMENTATION.md"), opts.Idea.Slug) == nil {
			members = append(members, quota.Member{ID: impl, Usable: true, ValidArtifact: true})
		}
	}
	for _, r := range results {
		if r.ExitError != "" && (r.AgentID == "runner" || strings.HasPrefix(r.AgentID, "runner/") || opts.participantDropout && r.QuotaEvidence == nil) {
			roles, _ := membership.Roles(opts.Idea.Path)
			decision := quota.Evaluate(h.Policy(), h.Current, members, roles)
			decision.Block = "integrity/recovery gate: " + r.ExitError
			err := membership.Block(opts.Root, opts.Idea.Path, opts.RunID, opts.RoundLabel, decision)
			return append(results, Result{AgentID: "runner/quota", ExitError: err.Error(), QuotaBlocked: membership.IsBlocked(err)})
		}
	}
	b, err := membership.Settle(ctx, opts.Root, opts.Idea.Path, opts.RunID, opts.RoundLabel, opts.Idea.Participants, members)
	if err != nil {
		return append(results, Result{AgentID: "runner/quota", ExitError: err.Error(), QuotaBlocked: membership.IsBlocked(err)})
	}
	if b != nil {
		for i := range results {
			if !membership.Has(b.Decision.After, results[i].AgentID) {
				results[i].QuotaExcluded = true
			}
		}
	}
	return results
}

func quotaSettleSingle(ctx context.Context, opts Options, r Result) Result {
	results := quotaSettle(ctx, opts, []Result{r})
	if len(results) > 1 {
		r.QuotaBlocked = results[len(results)-1].QuotaBlocked
		r.ExitError = combineError(r.ExitError, fmt.Errorf("%s", results[len(results)-1].ExitError))
		return r
	}
	return results[0]
}
