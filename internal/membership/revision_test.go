package membership

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/runmanifest"
)

func TestQuotaFixupOwnerRevisionReturnPolicyCatchupAndAuthority(t *testing.T) {
	root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
	ctx := leaseFixture(t, dir, run)
	if _, e := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates()); e != nil {
		t.Fatal(e)
	}
	history, _ := os.ReadFile(filepath.Join(dir, quota.HistoryDir, "000001.json"))
	ids := []string{"a", "b", "c"}
	p := quota.NewPolicy(nil, nil)
	a := quotatest.Authority(t, root, filepath.Base(dir), "return", (quota.RevisionDirective{Participants: ids, Policy: p}).Text())
	for _, bad := range []string{"missing", "wrong-idea", "wrong-quote", "changed-hash", "fabricated-blob", "uncommitted", "self-authored"} {
		t.Run(bad, func(t *testing.T) {
			invalid := a
			switch bad {
			case "missing":
				invalid.Path = "parley-deck/inbox/user-to-all_missing.md"
			case "wrong-idea":
				invalid = quotatest.Authority(t, root, "unrelated", "wrong-idea", a.Quote)
			case "wrong-quote":
				invalid.Quote = "I prefer something else"
			case "changed-hash":
				invalid.SHA256 = strings.Repeat("0", 64)
			case "fabricated-blob":
				invalid.Blob = strings.Repeat("0", 40)
			case "uncommitted":
				invalid.Commit = strings.Repeat("0", 40)
			case "self-authored":
				data, _ := os.ReadFile(filepath.Join(root, a.Path))
				data = bytes.Replace(data, []byte("from: user"), []byte("from: a"), 1)
				invalid.Path = "parley-deck/inbox/user-to-all_self-authored.md"
				if err := os.WriteFile(filepath.Join(root, invalid.Path), data, 0600); err != nil {
					t.Fatal(err)
				}
				quotatest.Git(t, root, "add", "--", invalid.Path)
				quotatest.Git(t, root, "commit", "-qm", "synthetic self-authored negative", "--", invalid.Path)
				invalid.Commit = quotatest.Git(t, root, "rev-parse", "HEAD")
				invalid.Blob = quotatest.Git(t, root, "rev-parse", invalid.Commit+":"+invalid.Path)
				invalid.SHA256 = fmt.Sprintf("%x", sha256.Sum256(data))
			}
			if _, e := Revise(ctx, root, dir, run, ids, p, invalid, nil); e == nil {
				t.Fatal("invalid authority accepted")
			}
		})
	}
	b, e := Revise(ctx, root, dir, run, ids, p, a, nil)
	if e != nil {
		t.Fatal(e)
	}
	if b.PriorRevision != 1 {
		t.Fatal(b)
	}
	// Permitted archive then deletion must not invalidate a committed decision.
	archived := filepath.Join(root, "parley-deck/inbox/archived", filepath.Base(a.Path))
	os.MkdirAll(filepath.Dir(archived), 0700)
	os.Rename(filepath.Join(root, a.Path), archived)
	for i := 0; i < 2; i++ {
		if i == 1 {
			os.Remove(archived)
		}
		h, e := Before(ctx, root, dir, run)
		if e != nil || strings.Join(h.Current, ",") != "a,b,c" || strings.Join(h.Known, ",") != "a,b,c,d" {
			t.Fatal(h, e)
		}
	}
	current, _ := os.ReadFile(filepath.Join(dir, quota.HistoryDir, "000001.json"))
	if !bytes.Equal(current, history) {
		t.Fatal("historical batch changed")
	}
	m, e := runmanifest.Load(root, run)
	if e != nil || m.QuotaRevision != 2 || strings.Join(m.Participants, ",") != "a,b,c" {
		t.Fatal(m, e)
	}
	// New identity is not re-inclusion. It requires the existing catch-up contract.
	join := []string{"a", "b", "c", "e"}
	ja := quotatest.Authority(t, root, filepath.Base(dir), "join", (quota.RevisionDirective{Participants: join, Policy: p}).Text())
	if _, e = Revise(ctx, root, dir, run, join, p, ja, nil); e == nil {
		t.Fatal("new identity joined without catch-up")
	}
	validRound(t, dir, filepath.Base(dir), "e")
	path := filepath.Join(dir, "round-01/e.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("---\n"), []byte("---\ncatch-up: true\nread-priors: [round-01/a.md, round-01/b.md]\njoin-from: round-02\n"), 1)
	os.WriteFile(path, raw, 0600)
	if _, e = Revise(ctx, root, dir, run, join, p, ja, map[string]string{"e": "round-01/e.md"}); e != nil {
		t.Fatal(e)
	}
	h, e := quota.ReadHistory(dir)
	if e != nil || !Has(h.Known, "e") {
		t.Fatal(h, e)
	}
}
func TestQuotaFixupScopeWideningIsExplicit(t *testing.T) {
	root, dir, run := fixture(t, quota.Policy{Enabled: true, Scope: quota.KickoffOnly})
	if _, e := Before(context.Background(), root, dir, run); e != nil {
		t.Fatal(e)
	}
	h, _ := quota.ReadHistory(dir)
	if h.MidIdea() {
		t.Fatal("resume widened")
	}
	p := quota.NewPolicy(nil, nil)
	a := quotatest.Authority(t, root, filepath.Base(dir), "widen", (quota.RevisionDirective{Participants: h.Current, Policy: p}).Text())
	if _, e := Revise(context.Background(), root, dir, run, h.Current, p, a, nil); e != nil {
		t.Fatal(e)
	}
	ctx, r, e := Acquire(context.Background(), dir, run)
	if e != nil {
		t.Fatal(e)
	}
	defer r()
	if _, _, e = Acquire(context.Background(), dir, "competing"); e == nil {
		t.Fatal("widening did not enable lease")
	}
	if _, e = Before(ctx, root, dir, run); e != nil {
		t.Fatal(e)
	}
}
func TestQuotaFixupKnobOffRecordedManualChanges(t *testing.T) {
	root, dir, run := fixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("participants: [a, b, c, d]"), []byte("participants: [a, b, c]\nexcluded: [d — unavailable — confirmed 2026-10-04]"), 1)
	os.WriteFile(path, raw, 0600)
	current, known, e := protocol.QuotaMembers(dir, nil)
	if e != nil || strings.Join(current, ",") != "a,b,c" || !Has(known, "d") {
		t.Fatal(current, known, e)
	}
	// Normal driver mutation imports the already-recorded confirmation: no new CLI.
	if _, e = Before(context.Background(), root, dir, run); e != nil {
		t.Fatal(e)
	}
	h, e := quota.ReadHistory(dir)
	if e != nil || h.Revision != 1 {
		t.Fatal(h, e)
	}
	raw, _ = os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("participants: [a, b, c]"), []byte("participants: [a, b, c, d]"), 1)
	os.WriteFile(path, raw, 0600)
	if _, e = Before(context.Background(), root, dir, run); e != nil {
		t.Fatal(e)
	}
	h, e = quota.ReadHistory(dir)
	if e != nil || h.Revision != 2 || len(h.Current) != 4 {
		t.Fatal(h, e)
	}
	_, r, e := Acquire(context.Background(), dir, "off-a")
	if e != nil {
		t.Fatal(e)
	}
	defer r()
	_, r2, e := Acquire(context.Background(), dir, "off-b")
	if e != nil {
		t.Fatal(e)
	}
	r2()
}

func TestQuotaFixupOffModeCatchupKeepsPreChangePath(t *testing.T) {
	root, dir, run := fixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("participants: [a, b, c, d]"), []byte("participants: [a, b, c, d, e]"), 1)
	raw = bytes.Replace(raw, []byte("status: round-01"), []byte("status: round-02"), 1)
	os.WriteFile(path, raw, 0600)
	validRound(t, dir, filepath.Base(dir), "e")
	late := filepath.Join(dir, "round-01/e.md")

	if _, e := Before(context.Background(), root, dir, run); e != nil {
		t.Fatal(e)
	}
	// The snapshot is immutable; permitted inbox cleanup does not revoke it.
	os.Remove(late)
	h, e := quota.ReadHistory(dir)
	if e != nil || !Has(h.Known, "e") || !Has(h.Current, "e") || h.Policy().Enabled {
		t.Fatal(h, e)
	}
	t.Log("off-mode catch-up: plain participants edit plus late round accepted through ordinary Before; immutable manual snapshot survives artifact cleanup")
}
func TestQuotaFixupClosedRevisionFrozen(t *testing.T) {
	root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("status: round-01"), []byte("status: final"), 1)
	os.WriteFile(path, raw, 0600)
	p := quota.NewPolicy(nil, nil)
	ids := []string{"a", "b"}
	a := quotatest.Authority(t, root, filepath.Base(dir), "closed", (quota.RevisionDirective{Participants: ids, Policy: p}).Text())
	if _, e := Revise(context.Background(), root, dir, run, ids, p, a, nil); e == nil {
		t.Fatal("closed revision accepted")
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(raw, after) {
		t.Fatal("closed prompt modified")
	}
}
