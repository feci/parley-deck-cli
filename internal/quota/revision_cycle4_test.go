package quota_test

import (
	"fmt"
	"testing"
	"time"

	"parley-deck-cli/internal/quota"
)

func TestQuotaCycle4ManualRoundOneHistoryUsesOnlySnapshot(t *testing.T) {
	p := quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea}
	k := quota.NewKickoff("idea", "run", p, []string{"a", "b"}, nil, time.Now())
	h := &quota.History{Kickoff: &k, Current: k.Participants, Known: k.Participants}
	for _, status := range []string{"round-01", "round-02", "consensus", "final", "closed", ""} {
		for _, marker := range []string{"", "excluded: new — unavailable — confirmed 2026-10-04\n"} {
			t.Run(status+"/"+marker, func(t *testing.T) {
				raw := fmt.Sprintf("---\nidea: idea\nparticipants: [a, b, new]\nquota_auto_exclude: false\nquota_auto_exclude_scope: kickoff-and-mid-idea\nstatus: %s\n%s---\n", status, marker)
				b := quota.NewRevision(h, "run", []string{"a", "b", "new"}, p, quota.Revision{ManualPrompt: raw}, time.Now())
				if err := b.Validate(h); (err == nil) != (status == "round-01") {
					t.Fatal("wrong immutable round predicate", err)
				}
			})
		}
	}
}
