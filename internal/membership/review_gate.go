package membership

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
)

const (
	legacyResetRule     = "quota.named-reset-ge-60m.v1"
	legacyAccountRule   = "quota.account-exhausted-no-reset.v1"
	legacyAllowanceRule = "quota.allowance-ge-24h-no-reset.v1"
	legacyProvenance    = "zcode.stderr.owner-deviation.2026-10-04"
)

// SingleReviewerAfterDropout changes only the numeric reviewer minimum. Cause
// comes from immutable history, never display markers or writable gate metadata.
// prospective is used only by Settle after Evaluate, before CommitBatch: requiring
// an already-committed transition there would deadlock the first eligible loss.
func SingleReviewerAfterDropout(root, ideaDir, runID, implementer string, ids []string, prospective *quota.Decision) (bool, error) {
	v, err := protocol.InspectQuota(ideaDir)
	if err != nil {
		return false, err
	}
	if v.Pending != "" || v.Manual != nil {
		return false, fmt.Errorf("single-reviewer proof requires settled membership history")
	}
	h := v.History
	if h == nil {
		return false, nil
	}
	if !validMemberSet(ids) || implementer == "" || !Has(ids, implementer) {
		return false, nil
	}
	role := protocol.ReadFacilitatorRole(ideaDir)
	if frozen := h.RoleSnapshot(); frozen != nil {
		role.Facilitator = frozen.Facilitator
	}
	if implementer == role.Facilitator && h.RoleSnapshot() == nil {
		return false, nil
	}
	reviewers := func(members []string) []string {
		var out []string
		for _, id := range members {
			if id != implementer && id != role.Facilitator {
				out = append(out, id)
			}
		}
		return out
	}
	survivors := reviewers(ids)
	if len(survivors) != 1 {
		return false, nil
	}
	var decision *quota.Decision
	proofIdea := h.Kickoff.Idea
	kickoff := false
	if prospective != nil {
		if !sameMemberSet(prospective.Before, h.Current) {
			return false, fmt.Errorf("prospective reviewer loss does not extend current history")
		}
		decision = prospective
	} else {
		if !sameMemberSet(ids, h.Current) {
			return false, fmt.Errorf("reviewer membership differs from immutable history")
		}
		// Policy-only revisions preserve cause. A later manual membership edit
		// invalidates it even if that edit happens to restore an older member set.
		for i := len(h.Batches) - 1; i >= 0; i-- {
			b := &h.Batches[i]
			if sameMemberSet(b.Decision.Before, b.Decision.After) {
				continue
			}
			if b.Owner != nil {
				return false, nil
			}
			decision = &b.Decision
			break
		}
		if decision == nil && h.Kickoff.Transition != nil {
			decision = &h.Kickoff.Transition.Decision
			// Readiness evidence is bound to the proposed batch, before the final
			// slug exists. Kickoff.Validate checks its shared readiness identity.
			proofIdea, kickoff = "", true
		}
	}
	beforeImplementer := implementer
	if decision != nil && decision.RoleTransition != nil {
		if decision.RoleTransition.After.Implementer != implementer {
			return false, fmt.Errorf("reviewer proof uses a stale implementer")
		}
		beforeImplementer = decision.RoleTransition.Before.Implementer
	}
	beforeReviewers := func(members []string) []string {
		var out []string
		for _, id := range members {
			if id != beforeImplementer && id != role.Facilitator {
				out = append(out, id)
			}
		}
		return out
	}
	if decision == nil || !decision.Applied || decision.Block != "" || decision.UsableSurvivors < quota.Floor ||
		!validMemberSet(decision.Before) || !sameMemberSet(decision.After, ids) ||
		(!Has(decision.Before, implementer) && !(decision.RoleTransition != nil && h.RoleSnapshot() != nil && implementer == h.RoleSnapshot().Facilitator)) || !Has(beforeReviewers(decision.Before), survivors[0]) || len(beforeReviewers(decision.Before)) < 2 {
		return false, nil
	}
	roles, err := Roles(ideaDir)
	if err != nil {
		return false, err
	}
	removed := []string{}
	for _, id := range ids {
		if !Has(decision.Before, id) && !(decision.RoleTransition != nil && h.RoleSnapshot() != nil && id == h.RoleSnapshot().Facilitator && id == implementer) {
			return false, nil
		}
	}
	for _, id := range decision.Before {
		if !Has(ids, id) {
			if (id == role.Facilitator || roles.Protected(id)) && !exhaustionCandidate(decision.Candidates, id) {
				return false, nil
			}
			removed = append(removed, id)
		}
	}
	if decision.RoleTransition != nil {
		for _, c := range decision.Candidates {
			if !Has(decision.Before, c.Agent) && h.RoleSnapshot() != nil && c.Agent == h.RoleSnapshot().Facilitator && c.Agent == decision.RoleTransition.Before.Organizer {
				removed = append(removed, c.Agent)
			}
		}
	}
	if len(removed) == 0 || len(removed) != len(decision.Candidates) || !sameMemberSet(removed, quota.CandidateIDs(decision.Candidates)) {
		return false, nil
	}
	for _, c := range decision.Candidates {
		if !reviewFailureEvidence(c.Evidence, proofIdea, c.Agent, kickoff) {
			return false, nil
		}
	}
	m, err := runmanifest.Load(root, runID)
	if err != nil {
		return false, fmt.Errorf("single-reviewer roster snapshot: %w", err)
	}
	if m.RunID != runID || m.IdeaSlug != filepath.Base(ideaDir) {
		return false, fmt.Errorf("single-reviewer roster snapshot belongs to another run/idea")
	}
	models := map[string]string{}
	for _, entry := range m.RosterSnapshot {
		if _, duplicate := models[entry.Agent]; duplicate {
			return false, fmt.Errorf("single-reviewer roster snapshot has duplicate identity %s", entry.Agent)
		}
		models[entry.Agent] = knownModel(entry.Model)
	}
	implModel, reviewerModel := models[implementer], models[survivors[0]]
	if implModel == "" || reviewerModel == "" || implModel == reviewerModel {
		return false, fmt.Errorf("single-reviewer exception requires known distinct snapshot models for implementer and reviewer")
	}
	return true, nil
}

func validMemberSet(ids []string) bool {
	if len(ids) != len(quota.Unique(ids)) {
		return false
	}
	for _, id := range ids {
		if strings.TrimSpace(id) != id || id == "" {
			return false
		}
	}
	return true
}

func sameMemberSet(a, b []string) bool {
	return validMemberSet(a) && validMemberSet(b) && reflect.DeepEqual(quota.Unique(a), quota.Unique(b))
}

func knownModel(model string) string {
	model = strings.ToLower(strings.TrimSpace(model))
	if model == agents.Unknown || model == agents.CLIDefault {
		return ""
	}
	return model
}

func exhaustionCandidate(cs []quota.Candidate, id string) bool {
	for _, c := range cs {
		if c.Agent == id && quota.ValidateExhaustionEvidence(c.Evidence, "", id) == nil {
			return true
		}
	}
	return false
}

func reviewFailureEvidence(e quota.Evidence, idea, agent string, kickoff bool) bool {
	if e.RuleID == quota.ProviderExhaustionRule {
		return quota.ValidateExhaustionEvidence(e, idea, agent) == nil && (!kickoff || e.Exhaustion.Step == "readiness")
	}
	if e.RuleID == quota.ParticipantFailureRule {
		return quota.ValidateFailureEvidence(e, idea, agent) == nil && (!kickoff || e.Failure.Step == "readiness")
	}
	if !e.Eligible || e.InvocationID == "" || e.ObservedAt.IsZero() || e.Excerpt == "" ||
		e.Adapter != "zcode" || e.Provenance != legacyProvenance || e.Failure != nil {
		return false
	}
	switch e.RuleID {
	case legacyResetRule:
		return e.ResetAt != nil && e.RawReset != "" && e.ResetAt.Sub(e.ObservedAt) >= quota.MinimumReset
	case legacyAccountRule, legacyAllowanceRule:
		return e.ResetAt == nil && e.RawReset == ""
	}
	return false
}
