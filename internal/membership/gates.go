package membership

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/track"
)

// CheckGates evaluates prospective membership with the same role chain and
// per-track counts. It only vetoes a reduction; it never certifies a close.
func CheckGates(root, ideaDir, runID string, ids []string) error {
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return err
	}
	auto := strings.EqualFold(meta["auto_implement"], "true")
	strict := strings.EqualFold(meta["strict_gate"], "true")
	_, implErr := os.Stat(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
	code := auto || strict || implErr == nil
	diversity := strings.EqualFold(meta["require_model_diversity"], "true") || strings.EqualFold(meta["track"], "fast")
	if !code && !diversity {
		return nil
	}
	role := protocol.ReadFacilitatorRole(ideaDir)
	eligible := []string{}
	for _, id := range ids {
		if !role.IneligibleForRoles(id) {
			eligible = append(eligible, id)
		}
	}
	des := protocol.ReadImplementerDesignation(ideaDir)
	if des.State == protocol.DesignationEmpty {
		return fmt.Errorf("incomplete per-idea implementer designation")
	}
	chain := protocol.PinImplementerCandidates()
	if des.State == protocol.DesignationSet {
		if !Has(eligible, des.ID) {
			return fmt.Errorf("per-idea designee absent from survivors")
		}
		chain = append(chain, protocol.ImplementerCandidate{ID: des.ID})
	} else if des.State == protocol.DesignationAbsent {
		defs, e := config.LoadDefaults(root)
		if e != nil {
			return fmt.Errorf("cannot resolve quota gate defaults: %w", e)
		}
		if defs.DefaultImplementer != "" && !strings.EqualFold(defs.DefaultImplementer, "none") {
			chain = append(chain, protocol.ImplementerCandidate{ID: defs.DefaultImplementer})
		}
	}
	chain = append(chain, protocol.LegacyImplementerCandidates()...)
	impl, _, ok := protocol.ResolveImplementerChain(ideaDir, chain, eligible)
	if !ok && len(eligible) > 0 {
		impl = eligible[0]
	}
	reviewers := []string{}
	for _, id := range eligible {
		if id != impl {
			reviewers = append(reviewers, id)
		}
	}
	t, present, err := track.NormalizeStrict(meta["track"])
	if err != nil {
		return err
	}
	policy, err := track.PolicyFor(t, present, len(reviewers), auto, strict)
	if err != nil {
		return err
	}
	if policy.MaxReviewers > 0 && len(reviewers) > policy.MaxReviewers {
		reviewers = reviewers[:policy.MaxReviewers]
	}
	minimum := policy.MinReviewers
	if minimum == 0 {
		minimum = 1
	}
	if auto && minimum < 2 {
		minimum = 2
	}
	if code && len(reviewers) < minimum {
		return fmt.Errorf("review/LE-7/LE-11 gate: %d independent reviewers remain; require %d", len(reviewers), minimum)
	}
	if (auto || strict) && len(reviewers) == 0 {
		return fmt.Errorf("no independent goal checker remains")
	}
	if diversity {
		m, err := runmanifest.Load(root, runID)
		if err != nil {
			return err
		}
		models := map[string]string{}
		for _, r := range m.RosterSnapshot {
			models[r.Agent] = r.Model
		}
		model := models[impl]
		different := false
		for _, id := range reviewers {
			if models[id] != "" && model != "" && models[id] != model {
				different = true
			}
		}
		if !different {
			return fmt.Errorf("require_model_diversity: survivors have no known model-diverse independent reviewer")
		}
	}
	return nil
}
