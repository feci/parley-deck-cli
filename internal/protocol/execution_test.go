package protocol

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
)

func handoverText(idea string, r *quota.RoleTransition) string {
	return fmt.Sprintf("---\nidea: %s\ntransition: %s\nauthor: %s\norganizer: %s\nimplementer: %s\nphase: %s\nbase: %s\nbase-sha256: %s\ntakeover-branch: %s\ntakeover-head: %s\ndate: 2026-10-29\nstatus: in-progress\n---\n\n## Takeover\nRoles: %v. Filed positions stay binding.\n\n## Resume state\nContinue from consensus.md after all signoffs.\n", idea, r.ID, r.Successor, r.After.Organizer, r.After.Implementer, r.Phase, r.Base, r.BaseSHA256, r.Branch, r.Head, r.Roles)
}

func TestActiveExecutionHistoryChainAndCorruption(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "idea")
	os.MkdirAll(dir, 0700)
	p := quota.NewExhaustionPolicy(nil, nil)
	s := &quota.RoleSnapshot{Version: 1, Order: []string{"a", "b", "c", "d"}, Organizer: "a"}
	k := quota.NewKickoff("idea", "run", p, s.Order, nil, time.Now(), s)
	if err := quota.WriteKickoff(dir, k); err != nil {
		t.Fatal(err)
	}
	before := quota.RoleAssignment{Organizer: "a", Implementer: "a", ImplementerSource: "designation"}
	base, hash := "none", ""
	var predecessor string
	for _, lost := range []string{"a", "b"} {
		h, err := quota.ReadHistory(dir)
		if err != nil {
			t.Fatal(err)
		}
		e := quotatest.Exhaustion("idea", lost, "round-02")
		var ms []quota.Member
		for _, id := range h.Current {
			m := quota.Member{ID: id, Usable: id != lost}
			if id == lost {
				m.Evidence = &e
			}
			ms = append(ms, m)
		}
		d := quota.Evaluate(p, h.Current, ms, quota.Roles{Snapshot: s, Assignment: before})
		d = d.BindHandover("round-02", base, hash, "feature/test", "abc123")
		if err := quota.CommitBatch(dir, quota.NewBatch(h, "run", "round-02", h.Current, d, time.Now())); err != nil {
			t.Fatal(err)
		}
		a, err := ResolveActiveExecution(dir)
		var pending *PendingHandover
		if !errors.As(err, &pending) || pending.Transition.Successor != d.RoleTransition.Successor {
			t.Fatalf("%+v %v", a, err)
		}
		os.MkdirAll(filepath.Dir(a.Path), 0700)
		os.WriteFile(filepath.Join(filepath.Dir(a.Path), "zz-fake.md"), []byte("status: complete"), 0600)
		raw := []byte(handoverText("idea", d.RoleTransition))
		os.WriteFile(a.Path, raw, 0600)
		resolved, err := ResolveActiveExecution(dir)
		if err != nil || resolved.Path != a.Path {
			t.Fatalf("%+v %v", resolved, err)
		}
		if started, err := ExecutionStarted(dir); started || err != nil {
			t.Fatalf("early takeover started Phase 5: %v %v", started, err)
		}
		if predecessor != "" {
			old, _ := os.ReadFile(predecessor)
			os.WriteFile(predecessor, append(old, 'x'), 0600)
			if _, err := ResolveActiveExecution(dir); err == nil {
				t.Fatal("changed predecessor accepted")
			}
			os.WriteFile(predecessor, old, 0600)
		}
		before, base, hash, predecessor = d.RoleTransition.After, d.RoleTransition.Path(), fmt.Sprintf("%x", sha256.Sum256(raw)), a.Path
	}
	context, err := ExecutionContext(dir)
	if err != nil || strings.Count(context, "## Takeover") != 2 {
		t.Fatalf("chain context: %v %s", err, context)
	}
	os.Remove(predecessor)
	os.Symlink("zz-fake.md", predecessor)
	if _, err := ResolveActiveExecution(dir); err == nil {
		t.Fatal("symlink authority accepted")
	}
}

func TestDraftRecoveryIntegrityAndExactBinding(t *testing.T) {
	for _, name := range []string{"consensus.md", "FINAL.md", "review/consensus.md"} {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "idea")
			path := filepath.Join(dir, filepath.FromSlash(name))
			os.MkdirAll(filepath.Dir(path), 0700)
			raw := "---\nidea: idea\nauthor: a\ndate: 2026-10-29\n---\n\n## Summary\nPartial valid draft.\n"
			cs := []quota.Candidate{{Agent: "a", Evidence: quotatest.Exhaustion("idea", "a", "draft")}}
			os.WriteFile(path, []byte(raw), 0600)
			plan, err := PlanDraftRecovery(dir, cs)
			if err != nil || len(plan) != 1 || plan[0].SHA256 != fmt.Sprintf("%x", sha256.Sum256([]byte(raw))) {
				t.Fatalf("%v %v", plan, err)
			}
			for _, bad := range []string{raw + "\n### Signoff a\nStatus: ACCEPT\n", raw + "\n```json\n{\n", strings.Replace(raw, "idea: idea", "idea: other", 1), strings.Replace(raw, "author: a", "author: a\nauthor: b", 1), strings.TrimSuffix(raw, "\n")} {
				os.WriteFile(path, []byte(bad), 0600)
				if _, err := PlanDraftRecovery(dir, cs); err == nil {
					t.Fatalf("invalid draft accepted: %s", bad)
				}
			}
		})
	}
	dir := filepath.Join(t.TempDir(), "idea")
	os.MkdirAll(dir, 0700)
	complete := "---\nidea: idea\ndrafted-by: a\ndate: 2026-10-29\n---\n\n"
	for _, h := range []string{"Agreed decisions", "Agreed trade-offs", "Open items deferred to implementation"} {
		complete += "## " + h + "\nNone.\n\n"
	}
	complete += "## Signoffs\n"
	os.WriteFile(filepath.Join(dir, "consensus.md"), []byte(complete), 0600)
	plan, err := PlanDraftRecovery(dir, []quota.Candidate{{Agent: "a", Evidence: quotatest.Exhaustion("idea", "a", "draft")}})
	if err != nil || len(plan) != 0 {
		t.Fatalf("complete draft moved: %v %v", plan, err)
	}
}
