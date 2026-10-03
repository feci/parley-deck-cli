package quota

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func candidate(id string) Member {
	return Member{ID: id, Evidence: &Evidence{InvocationID: id + "-inv", Adapter: "synthetic-test", Eligible: true, RuleID: "test", Provenance: "synthetic-test", ObservedAt: time.Now().UTC()}}
}
func TestQuotaBatchPermutations(t *testing.T) {
	p := NewPolicy(nil, nil)
	members := []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, candidate("c"), candidate("d")}
	expected := Evaluate(p, []string{"a", "b", "c", "d"}, members, Roles{})
	if !expected.Applied || !reflect.DeepEqual(expected.After, []string{"a", "b"}) {
		t.Fatal(expected)
	}
	var permute func(int)
	count := 0
	permute = func(n int) {
		if n == len(members) {
			got := Evaluate(p, []string{"a", "b", "c", "d"}, members, Roles{})
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("order-dependent: %+v", got)
			}
			count++
			return
		}
		for i := n; i < len(members); i++ {
			members[n], members[i] = members[i], members[n]
			permute(n + 1)
			members[n], members[i] = members[i], members[n]
		}
	}
	permute(0)
	if count != 24 {
		t.Fatal(count)
	}
}
func TestQuotaFloorRolesAndSuccess(t *testing.T) {
	p := NewPolicy(nil, nil)
	for _, tc := range []struct {
		name  string
		ids   []string
		ms    []Member
		roles Roles
		want  bool
	}{
		{"three-to-one", []string{"a", "c", "d"}, []Member{{ID: "a", Usable: true}, candidate("c"), candidate("d")}, Roles{}, false},
		{"two", []string{"a", "c"}, []Member{{ID: "a", Usable: true}, candidate("c")}, Roles{}, false},
		{"duplicate", []string{"a", "a", "c"}, []Member{{ID: "a", Usable: true}, {ID: "a", Usable: true}, candidate("c")}, Roles{}, false},
		{"unresolved", []string{"a", "b", "c"}, []Member{{ID: "a", Usable: true}, {ID: "b"}, candidate("c")}, Roles{}, false},
		{"facilitator", []string{"a", "b", "c"}, []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, candidate("c")}, Roles{Facilitator: "b"}, false},
		{"designee", []string{"a", "b", "c"}, []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, candidate("c")}, Roles{Designee: "c"}, false},
		{"pin", []string{"a", "b", "c"}, []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, candidate("c")}, Roles{PinnedImplementer: "c"}, false},
		{"drafter", []string{"a", "b", "c"}, []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, candidate("c")}, Roles{StartedDrafters: []string{"c"}}, false},
		{"global-default-unpinned", []string{"a", "b", "c"}, []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, candidate("c")}, Roles{}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := Evaluate(p, tc.ids, tc.ms, tc.roles)
			if d.Applied != tc.want {
				t.Fatal(d)
			}
			if !d.Applied && (d.Block == "" || !reflect.DeepEqual(d.After, d.Before)) {
				t.Fatal(d)
			}
		})
	}
	c := candidate("c")
	c.ValidArtifact = true
	d := Evaluate(p, []string{"a", "b", "c"}, []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, c}, Roles{})
	if d.Applied || len(d.Candidates) > 0 {
		t.Fatal(d)
	}
	c.ValidArtifact = false
	ms := []Member{c, {ID: "c", LaterSuccess: true}, {ID: "a", Usable: true}, {ID: "b", Usable: true}}
	d = Evaluate(p, []string{"a", "b", "c"}, ms, Roles{})
	if d.Applied || len(d.Candidates) > 0 {
		t.Fatal(d)
	}
}
func TestQuotaFrozenKickoffAndTruncation(t *testing.T) {
	dir := t.TempDir()
	p := NewPolicy(nil, nil)
	k := NewKickoff("idea", "run", p, []string{"a", "b"}, nil, time.Now())
	if err := WriteKickoff(dir, k); err != nil {
		t.Fatal(err)
	}
	if err := WriteKickoff(dir, k); err != nil {
		t.Fatal(err)
	}
	read, err := ReadKickoff(dir)
	if err != nil || read.Policy.Scope != KickoffOnly {
		t.Fatal(read, err)
	}
	k.Policy.Scope = KickoffAndMidIdea
	if WriteKickoff(dir, k) == nil {
		t.Fatal("upgrade widened scope")
	}
	if err = os.WriteFile(filepath.Join(dir, KickoffFile), []byte(`{"version":1`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = ReadKickoff(dir); err == nil {
		t.Fatal("truncated history accepted")
	}
}
func TestQuotaPolicyAndHints(t *testing.T) {
	f, tv := false, true
	if NewPolicy(&tv, &f).Enabled || NewPolicy(&f, nil).Enabled || !NewPolicy(nil, nil).Enabled {
		t.Fatal("presence merge")
	}
	e := Evidence{}
	if e.RelaunchHint() != "" {
		t.Fatal("unknown hint")
	}
	at := time.Date(2026, 10, 4, 22, 14, 57, 0, time.UTC)
	e.ResetAt = &at
	if e.RelaunchHint() != "owner-authorized relaunch after 2026-10-04T22:19:57Z (provider estimate)" {
		t.Fatal(e.RelaunchHint())
	}
}
