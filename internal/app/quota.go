package app

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"parley-deck-cli/internal/fsutil"
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
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if errors.Is(err, os.ErrExist) {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f, "---\nfrom: parley\nto: user\nphase: kickoff\nblocking: yes\ndate: %s\n---\n\nQuota batch blocked: %s\n\nCandidates: %v\n\nProposed: %v\n\nUsable non-facilitator survivors: %d; fixed floor: %d. No quota exclusions applied.\n\nEvidence:\n```json\n%s\n```\n", time.Now().UTC().Format("2006-01-02"), d.Block, quota.CandidateIDs(d.Candidates), d.Before, d.UsableSurvivors, quota.Floor, b)
	if err == nil {
		err = fsutil.SyncFile(f)
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
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
	lines = append(lines, "quota current participants: "+strings.Join(h.Current, ", "), "quota known participants: "+strings.Join(h.Known, ", "))
	if v.Pending != "" {
		lines = append(lines, "quota transition pending: "+v.Pending)
	}
	add := func(id string, cs []quota.Candidate) {
		for _, c := range cs {
			label := "automatic exclusion"
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
