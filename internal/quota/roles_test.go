package quota

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestOwnerSuccessionOrder(t *testing.T) {
	s := RoleSnapshot{Version: 1, Order: []string{"organizer", "implementer", "z", "a"}, Organizer: "organizer"}
	before := RoleAssignment{Organizer: "organizer", Implementer: "implementer", ImplementerSource: "designation"}
	for _, tc := range []struct {
		name                   string
		lost, usable           []string
		organizer, implementer string
		blocked                bool
	}{
		{"participant", []string{"a"}, []string{"organizer", "implementer", "z"}, "", "", false},
		{"organizer", []string{"organizer"}, []string{"implementer", "z", "a"}, "implementer", "implementer", false},
		{"implementer", []string{"implementer"}, []string{"organizer", "z", "a"}, "organizer", "organizer", false},
		{"both", []string{"implementer", "organizer"}, []string{"z", "a"}, "z", "z", false},
		{"simultaneous third loss", []string{"implementer", "z", "organizer"}, []string{"a"}, "a", "a", false},
		{"do not skip first survivor", []string{"implementer", "organizer"}, []string{"a"}, "", "", true},
		{"do not replace counterpart", []string{"organizer"}, []string{"z", "a"}, "", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, err := SelectSuccession(s, before, s.Order, tc.lost, tc.usable)
			if (err != nil) != tc.blocked {
				t.Fatalf("result=%+v error=%v", r, err)
			}
			if tc.blocked {
				return
			}
			if tc.organizer == "" {
				if r != nil {
					t.Fatal(r)
				}
				return
			}
			if r.After.Organizer != tc.organizer || r.After.Implementer != tc.implementer {
				t.Fatal(r)
			}
		})
	}
	before.Implementer = "organizer"
	r, err := SelectSuccession(s, before, s.Order, []string{"organizer"}, []string{"implementer", "z", "a"})
	if err != nil || r.After.Organizer != "implementer" || !reflect.DeepEqual(r.Roles, []string{"organizer", "implementer"}) {
		t.Fatalf("%+v %v", r, err)
	}
	before.Implementer = ""
	if _, err := SelectSuccession(s, before, s.Order, []string{"organizer"}, []string{"implementer", "z", "a"}); err == nil {
		t.Fatal("unknown counterpart accepted")
	}
}

func roleExhaustion(id string) Evidence {
	now := time.Date(2026, 10, 29, 20, 0, 0, 0, time.UTC)
	exit := 1
	message := "Your workspace is out of credits. Add credits to continue."
	return Evidence{InvocationID: id + "-inv", Adapter: "codex", RuleID: ProviderExhaustionRule, Provenance: CodexExhaustionProvenance,
		Eligible: true, Excerpt: message, ObservedAt: now,
		Exhaustion: &ExhaustedStep{Idea: "idea", Agent: id, Step: "round-02", CLIVersion: "0.165.0", ExplicitJSON: true,
			TerminalSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ValidationSHA256: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			Attempt: FailedAttempt{InvocationID: id + "-inv", Status: "failed", ExitCode: &exit, ObservedAt: now, Excerpt: message}}}
}

func TestExhaustionRoleBatchFloorAndExternalFacilitator(t *testing.T) {
	p := NewExhaustionPolicy(nil, nil)
	s := &RoleSnapshot{Version: 1, Order: []string{"impl", "r1", "r2"}, Organizer: "fac", Facilitator: "fac"}
	r := Roles{Snapshot: s, Facilitator: "fac", Designee: "impl", Assignment: RoleAssignment{Organizer: "fac", Implementer: "impl", ImplementerSource: "designation"}}
	e := roleExhaustion("impl")
	m := []Member{{ID: "impl", Evidence: &e}, {ID: "r1", Usable: true}, {ID: "r2", Usable: true}, {ID: "fac", Usable: true}}
	d := Evaluate(p, s.Order, m, r)
	if !d.Applied || d.UsableSurvivors != 2 || !reflect.DeepEqual(d.After, []string{"r1", "r2", "fac"}) || d.RoleTransition.After.Implementer != "fac" {
		t.Fatalf("%+v", d)
	}
	m[2].Usable = false
	d = Evaluate(p, s.Order, m, r)
	if d.Applied || d.Block == "" {
		t.Fatal("facilitator counted as floor seat")
	}

	external := roleExhaustion("fac")
	m = []Member{{ID: "impl", Usable: true}, {ID: "r1", Usable: true}, {ID: "r2", Usable: true}, {ID: "fac", Evidence: &external}}
	d = Evaluate(p, s.Order, m, r)
	if !d.Applied || !reflect.DeepEqual(d.After, s.Order) || d.RoleTransition.After.Organizer != "impl" {
		t.Fatalf("role-only: %+v", d)
	}
}

func TestRoleFieldsDoNotChangeLegacySerialization(t *testing.T) {
	p := NewParticipantPolicy(nil, nil)
	k := NewKickoff("idea", "run", p, []string{"a", "b"}, nil, time.Now())
	raw, _ := json.Marshal(k)
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(raw, &fields)
	if fields["roles"] != nil {
		t.Fatal(string(raw))
	}
	if err := k.Validate(); err != nil {
		t.Fatal(err)
	}
	k.Policy = NewExhaustionPolicy(nil, nil)
	if k.Validate() == nil {
		t.Fatal("v2 accepted missing frozen role order")
	}
}
