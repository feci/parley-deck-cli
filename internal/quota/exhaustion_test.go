package quota

import (
	"encoding/json"
	"testing"
)

func TestExhaustionPolicyDoesNotWidenSavedBytes(t *testing.T) {
	for _, raw := range []string{`{"enabled":true,"scope":"kickoff-and-mid-idea"}`, `{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":"participant-failure-v1"}`} {
		var p Policy
		if err := json.Unmarshal([]byte(raw), &p); err != nil {
			t.Fatal(err)
		}
		out, _ := json.Marshal(p)
		if string(out) != raw || p.Exhaustion() {
			t.Fatalf("saved policy changed: %s", out)
		}
	}
	p := NewExhaustionPolicy(nil, nil)
	if !p.Exhaustion() || !p.Dropout() || p.Trigger != ParticipantFailureV2 {
		t.Fatal(p)
	}
	for _, raw := range []string{`{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":null}`, `{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":"participant-failure-v4"}`, `{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":"participant-failure-v2","trigger":"participant-failure-v1"}`} {
		if json.Unmarshal([]byte(raw), &p) == nil {
			t.Fatal(raw)
		}
	}
}

func TestOptionalParticipantPolicyIsNewAndExplicit(t *testing.T) {
	p := NewOptionalParticipantPolicy(nil, nil)
	if !p.Exhaustion() || !p.Dropout() || !p.OptionalParticipants() || p.Trigger != ParticipantFailureV3 || p.MinUsableSurvivors() != 0 || p.MinMembers() != 1 {
		t.Fatal(p)
	}
	raw, _ := json.Marshal(p)
	var back Policy
	if err := json.Unmarshal(raw, &back); err != nil || back != p {
		t.Fatalf("%s %v", raw, err)
	}
	v2 := NewExhaustionPolicy(nil, nil)
	if v2.OptionalParticipants() || v2.MinUsableSurvivors() != Floor || v2.MinMembers() != Floor {
		t.Fatal("v2 widened")
	}
	off := NewOptionalParticipantPolicy(nil, nil)
	off.Enabled = false
	if off.OptionalParticipants() || off.Exhaustion() {
		t.Fatal("disabled v3 policy acts")
	}
}
