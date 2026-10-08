package consensus

import (
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/telemetry"
)

func TestDropoutHistoricalBlockAndFindingsStayBinding(t *testing.T) {
	root, idea, ctx := quotaConsensusFixturePolicy(t, quota.NewParticipantPolicy(nil, nil))
	path := filepath.Join(idea.Path, "consensus.md")
	raw := "---\nidea: " + idea.Slug + "\ndrafted-by: a\n---\n## Signoffs\n" + signoffBlock("a", "2026-10-08", StatusAccept, "", "") + signoffBlock("b", "2026-10-08", StatusAccept, "", "") + signoffBlock("c", "2026-10-08", StatusBlock, "Historical data-loss veto.", "Preserve records.")
	writeFile(t, path, raw)
	writeFile(t, filepath.Join(idea.Path, "round-01", "c.md"), "---\nagent: c\n---\n### [MAJOR] Loss of history\nClaim C1 is DISPUTED\n")
	ms := []quota.Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}}
	for _, id := range []string{"c", "d"} {
		now := time.Now().UTC()
		code := 1
		a := quota.FailedAttempt{InvocationID: id + "-1", Status: "failed", ExitCode: &code, FailureClass: "process_failure", ObservedAt: now, Excerpt: "exit_code: 1"}
		b := a
		b.InvocationID = id + "-2"
		b.RetryOf = a.InvocationID
		b.ObservedAt = now.Add(5 * time.Second)
		e, err := telemetry.ParticipantEvidence(idea.Slug, id, "later-step", "fixture", []quota.FailedAttempt{a, b})
		if err != nil {
			t.Fatal(err)
		}
		ms = append(ms, quota.Member{ID: id, Evidence: &e})
	}
	if _, err := membership.Settle(ctx, root, idea.Path, "test-run", "later-step", idea.Participants, ms); err != nil {
		t.Fatal(err)
	}
	s, err := Status(root, idea.Slug, false)
	if err != nil || s.Triage != TriageBlocked || len(s.Errors) > 0 || len(s.Retained) != 3 {
		t.Fatalf("%+v %v", s, err)
	}
	if _, err := AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "c", Status: "accept"}); err == nil {
		t.Fatal("dropped author signed")
	}
	if _, err := AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "d", Status: "block", Notes: "❌ NON-PARTICIPANT", CounterProposal: "Continue"}); err == nil {
		t.Fatal("dropped author used decline")
	}
}
