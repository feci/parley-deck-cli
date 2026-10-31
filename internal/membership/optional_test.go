package membership

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/telemetry"
)

// failedMembers gives a usable implementer "a" and typed participant-failure
// evidence for every other member.
func failedMembers(t *testing.T, idea string, failed ...string) []quota.Member {
	t.Helper()
	ms := []quota.Member{{ID: "a", Usable: true}}
	for _, id := range failed {
		now := time.Now().UTC()
		code := 1
		first := quota.FailedAttempt{InvocationID: id + "-1", Status: "failed", FailureClass: "process_failure", ExitCode: &code, ObservedAt: now, Excerpt: "exit_code: 1"}
		second := first
		second.InvocationID = id + "-2"
		second.RetryOf = first.InvocationID
		second.ObservedAt = now.Add(5 * time.Second)
		e, err := telemetry.ParticipantEvidence(idea, id, "round-01/"+id+".md", "fixture", []quota.FailedAttempt{first, second})
		if err != nil {
			t.Fatal(err)
		}
		ms = append(ms, quota.Member{ID: id, Evidence: &e})
	}
	return ms
}

func optionalFixture(t *testing.T, p quota.Policy) (string, string, string) {
	t.Helper()
	root, dir, run := fixture(t, p)
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	os.WriteFile(path, []byte(strings.Replace(string(raw), "status:", "auto_implement: true\nimplementer: a\nrequire_model_diversity: true\nstatus:", 1)), 0600)
	reviewModels(t, root, run, "Model-A", "Model-B")
	return root, dir, run
}

func TestOptionalParticipantsKeepTheSoleImplementer(t *testing.T) {
	root, dir, run := optionalFixture(t, quota.NewOptionalParticipantPolicy(nil, nil))
	ctx := leaseFixture(t, dir, run)
	b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, failedMembers(t, filepath.Base(dir), "b", "c", "d"))
	if err != nil || b == nil {
		t.Fatalf("v3 reduction to the implementer: %+v %v", b, err)
	}
	if got := b.Decision.After; len(got) != 1 || got[0] != "a" {
		t.Fatalf("after: %v", got)
	}
	if err := CheckGates(root, dir, run, []string{"a"}); err != nil {
		t.Fatalf("v3 gate vetoed a proved zero-reviewer loss: %v", err)
	}
	if ok, err := NoReviewerAfterDropout(root, dir, run, "a", []string{"a"}, nil); !ok || err != nil {
		t.Fatalf("zero-reviewer proof: %v %v", ok, err)
	}
	if ok, _ := SingleReviewerAfterDropout(root, dir, run, "a", []string{"a"}, nil); ok {
		t.Fatal("zero reviewers passed as one reviewer")
	}
}

func TestOptionalParticipantsNeedV3AndProof(t *testing.T) {
	root, dir, run := optionalFixture(t, quota.NewExhaustionPolicy(nil, nil))
	ctx := leaseFixture(t, dir, run)
	if _, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, failedMembers(t, filepath.Base(dir), "b", "c", "d")); err == nil {
		t.Fatal("v2 accepted a reduction below two usable non-facilitators")
	}
	if ok, _ := NoReviewerAfterDropout(root, dir, run, "a", []string{"a"}, nil); ok {
		t.Fatal("v2 earned the v3 zero-reviewer exception")
	}

	root, dir, run = optionalFixture(t, quota.NewOptionalParticipantPolicy(nil, nil))
	if ok, _ := NoReviewerAfterDropout(root, dir, run, "a", []string{"a"}, nil); ok {
		t.Fatal("no automatic transition earned the zero-reviewer exception")
	}
	if err := CheckGates(root, dir, run, []string{"a"}); err == nil {
		t.Fatal("an unproved solo membership passed the gate")
	}
}

func TestOptionalParticipantsProofSurvivesPolicyOnlyRevision(t *testing.T) {
	root, dir, run := optionalFixture(t, quota.NewOptionalParticipantPolicy(nil, nil))
	ctx := leaseFixture(t, dir, run)
	if _, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, failedMembers(t, filepath.Base(dir), "b", "c", "d")); err != nil {
		t.Fatal(err)
	}
	off := quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea}
	a := quotatest.Authority(t, root, filepath.Base(dir), "off", (quota.RevisionDirective{Participants: []string{"a"}, Policy: off}).Text())
	if _, err := Revise(ctx, root, dir, run, []string{"a"}, off, a, nil); err != nil {
		t.Fatal(err)
	}
	if ok, err := NoReviewerAfterDropout(root, dir, run, "a", []string{"a"}, nil); !ok || err != nil {
		t.Fatalf("policy-only revision lost the v3 cause: %v %v", ok, err)
	}
}
