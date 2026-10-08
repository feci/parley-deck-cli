package quota

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func dropoutEvidence(id string) Evidence {
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	code := 1
	a := FailedAttempt{InvocationID: id + "-first", Status: "failed", FailureClass: "process_failure", ExitCode: &code, ObservedAt: now, Excerpt: `"exit_code":1`}
	b := a
	b.InvocationID, b.RetryOf, b.ObservedAt = id+"-second", a.InvocationID, now.Add(5*time.Second)
	return Evidence{InvocationID: b.InvocationID, Adapter: "test", RuleID: ParticipantFailureRule, Provenance: "supervisor-terminal", Eligible: true, Reason: "two settled child failures", Excerpt: b.Excerpt, ObservedAt: b.ObservedAt,
		Failure: &FailedStep{Idea: "idea", Agent: id, Step: "round-01/" + id + ".md", Attempts: []FailedAttempt{a, b}}}
}

func TestDropoutPolicyLegacyBytesAndStrictExtension(t *testing.T) {
	legacy := `{"enabled":true,"scope":"kickoff-and-mid-idea"}`
	var p Policy
	if err := json.Unmarshal([]byte(legacy), &p); err != nil || p.Dropout() {
		t.Fatal(p, err)
	}
	got, err := json.Marshal(p)
	if err != nil || string(got) != legacy {
		t.Fatalf("legacy hash input changed: %s %v", got, err)
	}
	p = NewParticipantPolicy(nil, nil)
	encoded, _ := json.Marshal(p)
	var decoded Policy
	if err := json.Unmarshal(encoded, &decoded); err != nil || decoded != p || !decoded.Dropout() {
		t.Fatal(decoded, err)
	}
	for _, value := range []string{`null`, `""`, `false`, `"future"`} {
		raw := strings.TrimSuffix(legacy, "}") + `,"trigger":` + value + `}`
		if json.Unmarshal([]byte(raw), &decoded) == nil {
			t.Fatalf("accepted invalid trigger %s", raw)
		}
	}
	for _, raw := range []string{
		`{"enabled":true,"scope":"kickoff-and-mid-idea","other":"participant-failure-v1"}`,
		`{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":"participant-failure-v1","trigger":"participant-failure-v1"}`,
	} {
		if json.Unmarshal([]byte(raw), &decoded) == nil {
			t.Fatalf("accepted ambiguous policy %s", raw)
		}
	}
	off, on := false, true
	if NewParticipantPolicy(&on, &off).Enabled || NewParticipantPolicy(&off, nil).Enabled {
		t.Fatal("old opt-out no longer disables new trigger")
	}
}

func TestDropoutProofAndEvidenceBackedImplementerFloor(t *testing.T) {
	p := NewParticipantPolicy(nil, nil)
	e := dropoutEvidence("c")
	ms := []Member{{ID: "a", Usable: true}, {ID: "b", ValidArtifact: true}, {ID: "c", Evidence: &e}}
	d := Evaluate(p, []string{"a", "b", "c"}, ms, Roles{Designee: "a"})
	if !d.Applied || len(d.After) != 2 {
		t.Fatal(d)
	}
	// Even two unrelated usable seats cannot replace an unresolved implementer.
	d = Evaluate(p, []string{"a", "b", "c", "impl"}, ms, Roles{PinnedImplementer: "impl"})
	if d.Applied || !strings.Contains(d.Block, "positive usability") {
		t.Fatal(d)
	}
	ms[2].ValidArtifact = true // A valid BLOCK/dispute is equally valid evidence.
	if d = Evaluate(p, []string{"a", "b", "c"}, ms, Roles{}); d.Applied || len(d.Candidates) != 0 {
		t.Fatal(d)
	}
	for _, mutate := range []func(*Evidence){
		func(e *Evidence) { e.Failure.Attempts = e.Failure.Attempts[:1] },
		func(e *Evidence) { e.Failure.Attempts[1].RetryOf = "unrelated" },
		func(e *Evidence) { e.Failure.Attempts[1].FailureClass = "budget_refused" },
		func(e *Evidence) { e.Failure.Attempts[0].FailureClass = "cancelled" },
		func(e *Evidence) { e.RuleID = "quota.named-reset-ge-60m.v1" },
	} {
		bad := dropoutEvidence("c")
		mutate(&bad)
		if ValidateFailureEvidence(bad, "idea", "c") == nil {
			t.Fatalf("accepted bad proof: %+v", bad)
		}
	}
}

func TestDropoutPermanentAcrossPolicyDowngradeAndKickoff(t *testing.T) {
	p := NewParticipantPolicy(nil, nil)
	e := dropoutEvidence("c")
	e.Failure.Step = "readiness"
	d := Evaluate(p, []string{"a", "b", "c"}, []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, {ID: "c", Evidence: &e}}, Roles{})
	k := NewKickoff("idea", "run", p, d.After, &d, time.Now())
	h := &History{Kickoff: &k, Current: d.After, Known: d.After}
	if err := k.Validate(); err != nil {
		t.Fatal(err)
	}
	if !h.Dropped("c") || h.CheckReturn([]string{"a", "c"}) == nil {
		t.Fatal("kickoff-dropped non-historical identity can rejoin")
	}
	h.Batches = []Batch{{Owner: &Revision{}, Policy: Policy{Enabled: false, Scope: KickoffAndMidIdea}}}
	if h.Policy().Enabled || h.CheckReturn([]string{"a", "c"}) == nil {
		t.Fatal("opt-out erased terminal dropout")
	}
	b := NewRevision(h, "later-run", []string{"a", "b", "c"}, h.Policy(), Revision{}, time.Now())
	if err := b.validateRevision(h); err == nil || !strings.Contains(err.Error(), "permanently dropped") {
		t.Fatal(err)
	}
	k.Transition.Decision.Candidates[0].Evidence.RuleID = "quota.account-exhausted-no-reset.v1"
	if h.Dropped("c") || h.CheckReturn([]string{"c"}) != nil {
		t.Fatal("legacy quota return unexpectedly banned")
	}
}
