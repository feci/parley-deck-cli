package trajectory

import (
	"context"
	"errors"
	"slices"

	"gopkg.in/yaml.v3"

	"parley-deck-cli/internal/budget"
)

// Activation pins the earliest source this trajectory actually retained. It
// cannot establish an earlier Phase-0 membership absent from that archive.
func archivedQuorum(ctx context.Context, b budget.CycleBinding, s State, ref SnapshotRef, source Source) ([]string, error) {
	name := "parley-deck/ideas/" + s.Policy.Idea + "/00-prompt.md"
	raw, err := readSnapshotMember(ctx, snapshotDirectory(b), ref, source, name, 1<<20)
	if err != nil {
		return nil, err
	}
	frontmatter, err := reconciliationFrontmatter(raw)
	if err != nil {
		return nil, err
	}
	// The lower-level policy API already pins explicit criterion arguments.
	// Pin membership without imposing a new checks syntax on that old API.
	var fm struct {
		Participants []string `yaml:"participants"`
	}
	if err := yaml.Unmarshal(frontmatter, &fm); err != nil {
		return nil, err
	}
	members := fm.Participants
	if len(members) < 2 || !slices.Contains(members, s.Policy.Implementer) {
		return nil, errors.New("archived quorum lacks the original implementer or independent participant")
	}
	seen := map[string]bool{}
	for _, id := range members {
		if !safeLabel(id) || seen[id] {
			return nil, errors.New("archived quorum contains invalid or duplicate participant identities")
		}
		seen[id] = true
	}
	return members, nil
}

func activationQuorum(ctx context.Context, b budget.CycleBinding, s State) ([]string, error) {
	return archivedQuorum(ctx, b, s, s.BaselineArchive, s.Policy.Baseline)
}

// This is permission to prepare/execute/accept verification, not the structural
// ticket check used to retain evidence or stop an already issued old ticket.
// A bad patch must still retain its terminal, archives and spent reservation.
func checkCapturedActivationQuorum(ctx context.Context, b budget.CycleBinding, s State, r CapturedRequest) ([]string, error) {
	members, err := activationQuorum(ctx, b, s)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(members, r.Verifier) || r.Verifier == s.Policy.Implementer {
		return nil, errors.New("selected verifier is outside the original independent activation quorum")
	}
	for _, original := range []struct {
		ref    SnapshotRef
		source Source
	}{{r.BeforeArchive, r.Before}, {r.AfterArchive, r.After}} {
		current, err := archivedQuorum(ctx, b, s, original.ref, original.source)
		if err != nil {
			return nil, err
		}
		if !slices.Equal(current, members) {
			return nil, errors.New("captured patch changed the original activation quorum")
		}
	}
	return members, nil
}
