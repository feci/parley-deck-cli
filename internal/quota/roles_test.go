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

func TestExhaustionSuccessionUsesNativeRoleLossAndWholeBatch(t *testing.T) {
	for _, source := range []string{"legacy", "global-default"} {
		for _, tc := range []struct {
			name, organizerRule, implementerRule, peerRule, successor string
			pinned, designee, blocked                                 bool
		}{
			{name: "ordinary default", implementerRule: ParticipantFailureRule},
			{name: "ordinary default and exhausted peer", implementerRule: ParticipantFailureRule, peerRule: ProviderExhaustionRule},
			{name: "exhausted organizer and ordinary peer", organizerRule: ProviderExhaustionRule, peerRule: ParticipantFailureRule, successor: "impl"},
			{name: "both roles and ordinary first survivor", organizerRule: ProviderExhaustionRule, implementerRule: ProviderExhaustionRule, peerRule: ParticipantFailureRule, successor: "z"},
			{name: "unusable ordinary counterpart", organizerRule: ProviderExhaustionRule, implementerRule: ParticipantFailureRule, blocked: true},
			{name: "ordinary pinned implementer", implementerRule: ParticipantFailureRule, pinned: true, blocked: true},
			{name: "ordinary designated implementer", implementerRule: ParticipantFailureRule, designee: true, blocked: true},
		} {
			t.Run(source+"/"+tc.name, func(t *testing.T) {
				p := NewExhaustionPolicy(nil, nil)
				s := &RoleSnapshot{Version: 1, Order: []string{"org", "impl", "peer", "z", "a"}, Organizer: "org"}
				r := Roles{Snapshot: s, Assignment: RoleAssignment{Organizer: "org", Implementer: "impl", ImplementerSource: source}}
				if tc.pinned {
					r.PinnedImplementer = "impl"
					r.Assignment.ImplementerSource = "pin"
				}
				if tc.designee {
					r.Designee = "impl"
					r.Assignment.ImplementerSource = "designation"
				}
				var members []Member
				for _, id := range s.Order {
					rule := map[string]string{"org": tc.organizerRule, "impl": tc.implementerRule, "peer": tc.peerRule}[id]
					m := Member{ID: id, Usable: rule == ""}
					if rule != "" {
						e := dropoutEvidence(id)
						e.Failure.Step = "readiness"
						if rule == ProviderExhaustionRule {
							e = roleExhaustion(id)
							e.Exhaustion.Step = "readiness"
						}
						m.Evidence = &e
					}
					members = append(members, m)
				}
				d := Evaluate(p, s.Order, members, r)
				if tc.blocked {
					if d.Applied || d.Block == "" || d.RoleTransition != nil {
						t.Fatalf("unusable/protected role did not block: %+v", d)
					}
					return
				}
				if !d.Applied || (d.RoleTransition != nil) != (tc.successor != "") {
					t.Fatalf("wrong transition class: %+v", d)
				}
				if tc.successor != "" && d.RoleTransition.Successor != tc.successor {
					t.Fatalf("wrong successor: %+v", d.RoleTransition)
				}
				k := NewKickoff("idea", "run", p, d.After, &d, time.Now(), s)
				if err := k.Validate(); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
