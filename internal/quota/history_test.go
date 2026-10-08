package quota

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestQuotaBatchCommitFaultsAndTruncation(t *testing.T) {
	for _, stage := range []string{"create", "write", "sync", "publish", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			dir := t.TempDir()
			p := NewPolicy(nil, nil)
			k := NewKickoff("idea", "run", p, []string{"a", "b", "c"}, nil, time.Now())
			if err := WriteKickoff(dir, k); err != nil {
				t.Fatal(err)
			}
			h, _ := ReadHistory(dir)
			e := Evidence{InvocationID: "c", Eligible: true, RuleID: "test", Provenance: "test", ObservedAt: time.Now(), Excerpt: "weekly limit exhausted"}
			d := Evaluate(p, h.Current, []Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, {ID: "c", Evidence: &e}}, Roles{})
			b := NewBatch(h, "run", "", nil, d, time.Now())
			writeFault = func(s, path string) error {
				if s == stage {
					return fmt.Errorf("injected %s", stage)
				}
				return nil
			}
			defer func() { writeFault = func(string, string) error { return nil } }()
			if err := CommitBatch(dir, b); err == nil {
				t.Fatal("missing injected failure")
			}
			h, err := ReadHistory(dir)
			if err != nil {
				t.Fatal(err)
			}
			if stage == "directory-sync" {
				if h.Revision != 1 {
					t.Fatal(h)
				}
			} else if h.Revision != 0 {
				t.Fatal(h)
			}
			writeFault = func(string, string) error { return nil }
			if err = CommitBatch(dir, b); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, HistoryDir, "000001.json")
			raw, _ := os.ReadFile(path)
			os.WriteFile(path, raw[:len(raw)/2], 0600)
			if _, err = ReadHistory(dir); err == nil {
				t.Fatal("truncation accepted")
			}
			if err = CommitBatch(dir, b); err == nil {
				t.Fatal("truncation overwritten")
			}
		})
	}
}

func TestQuotaPolicyJSONRequiresPresenceAndHistoryRejectsDuplicateKeys(t *testing.T) {
	for _, raw := range []string{`{"scope":"kickoff-only"}`, `{"enabled":null,"scope":"kickoff-only"}`, `{"enabled":false}`, `{"enabled":"false","scope":"kickoff-only"}`, `{"enabled":false,"scope":"kickoff-only","other":1}`} {
		var p Policy
		if json.Unmarshal([]byte(raw), &p) == nil {
			t.Fatal("ambiguous policy accepted", raw)
		}
	}
	dir := t.TempDir()
	k := NewKickoff("idea", "run", NewPolicy(nil, nil), []string{"a", "b"}, nil, time.Now())
	if err := WriteKickoff(dir, k); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, KickoffFile)
	raw, _ := os.ReadFile(path)
	os.WriteFile(path, []byte(strings.Replace(string(raw), `"enabled": true`, `"enabled": false, "enabled": true`, 1)), 0600)
	if _, err := ReadHistory(dir); err == nil {
		t.Fatal("duplicate policy key accepted")
	}
}
func TestQuotaNewScopeAndOldScopeImmutable(t *testing.T) {
	if NewPolicy(nil, nil).Scope != KickoffAndMidIdea {
		t.Fatal("new delivery not enabled")
	}
	for _, p := range []Policy{{Enabled: true, Scope: KickoffOnly}, {Enabled: false, Scope: KickoffAndMidIdea}} {
		dir := t.TempDir()
		k := NewKickoff("idea", "run", p, []string{"a", "b"}, nil, time.Now())
		if err := WriteKickoff(dir, k); err != nil {
			t.Fatal(err)
		}
		k.Policy = NewPolicy(nil, nil)
		if WriteKickoff(dir, k) == nil {
			t.Fatal("upgrade widened policy")
		}
	}
}
