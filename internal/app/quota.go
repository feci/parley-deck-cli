package app

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
)

// A failed kickoff has no idea identity yet. Keep one decision-scoped blocking
// record, with all candidates and the complete arithmetic; never one per agent.
func writeQuotaBlock(root string, d quota.Decision) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	id := fmt.Sprintf("%x", sha256.Sum256(b))
	path := filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_quota-kickoff-"+id+".md")
	text := fmt.Sprintf("---\nfrom: parley\nto: user\nphase: kickoff\nblocking: yes\ndate: %s\n---\n\nAutomatic exclusion batch blocked: %s\n\nCandidates: %v\n\nProposed: %v\n\nUsable non-facilitator survivors: %d; fixed floor: %d. No exclusions applied.\n\n%s\n\nEvidence:\n```json\n%s\n```\n", time.Now().UTC().Format("2006-01-02"), d.Block, quota.CandidateIDs(d.Candidates), d.Before, d.UsableSurvivors, quota.Floor, quota.OwnerOptions, b)
	return quota.PublishNotice(filepath.Dir(path), filepath.Base(path), text)
}
func quotaSurface(ideaDir string) []string {
	v, err := protocol.InspectQuota(ideaDir)
	if err != nil {
		return []string{"quota transition pending/integrity gate: " + err.Error()}
	}
	h := v.History
	if h == nil {
		return nil
	}
	lines := []string{fmt.Sprintf("quota policy: enabled=%t scope=%s revision=%d", h.Policy().Enabled, h.Policy().Scope, h.Revision)}
	if h.Policy().Trigger != "" {
		lines = append(lines, "automatic exclusion trigger: "+h.Policy().Trigger)
	}
	lines = append(lines, "quota current participants: "+strings.Join(h.Current, ", "), "quota known participants: "+strings.Join(h.Known, ", "))
	if v.Pending != "" {
		lines = append(lines, "quota transition pending: "+v.Pending)
	}
	add := func(id string, cs []quota.Candidate) {
		for _, c := range cs {
			label := "automatic exclusion"
			if h.Dropped(c.Agent) {
				label = "permanent participant dropout"
			}
			for _, member := range h.Current {
				if member == c.Agent {
					label = "historical automatic exclusion (now re-included)"
				}
			}
			lines = append(lines, fmt.Sprintf("%s: %s; reset=%s; transition=%s", label, c.Agent, c.Evidence.ResetHint(), id))
		}
	}
	if h.Kickoff.Transition != nil {
		add(h.Kickoff.Transition.ID, h.Kickoff.Transition.Decision.Candidates)
	}
	for _, b := range h.Batches {
		add(b.ID, b.Decision.Candidates)
	}
	return lines
}
func quotaSurfaceText(ideaDir string) string { return strings.Join(quotaSurface(ideaDir), "\n") }
