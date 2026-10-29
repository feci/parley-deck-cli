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
	for _, raw := range []string{`{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":null}`, `{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":"participant-failure-v3"}`, `{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":"participant-failure-v2","trigger":"participant-failure-v1"}`} {
		if json.Unmarshal([]byte(raw), &p) == nil {
			t.Fatal(raw)
		}
	}
}
