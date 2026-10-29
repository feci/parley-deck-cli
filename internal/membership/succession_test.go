package membership

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/runmanifest"
)

func TestOwnerWideningFreezesRolesWithoutKickoffRewrite(t *testing.T) {
	root, dir, run := fixture(t, quota.NewParticipantPolicy(nil, nil))
	before, _ := os.ReadFile(filepath.Join(dir, quota.KickoffFile))
	s := &quota.RoleSnapshot{Version: 1, Order: []string{"a", "b", "c", "d"}, Organizer: "a"}
	p := quota.NewExhaustionPolicy(nil, nil)
	directive := quota.RevisionDirective{Participants: s.Order, Policy: p, RoleSnapshot: s}
	a := quotatest.Authority(t, root, filepath.Base(dir), "widen", directive.Text())
	ctx := leaseFixture(t, dir, run)
	if _, err := Revise(ctx, root, dir, run, s.Order, p, a, nil); err == nil {
		t.Fatal("widening without roles accepted")
	}
	if _, err := Revise(ctx, root, dir, run, s.Order, p, a, nil, s); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(filepath.Join(dir, quota.KickoffFile))
	if string(after) != string(before) {
		t.Fatal("old kickoff rewritten")
	}
	h, err := quota.ReadHistory(dir)
	if err != nil || !h.Policy().Exhaustion() || !reflect.DeepEqual(h.RoleSnapshot(), s) {
		t.Fatalf("%+v %v", h, err)
	}
}

func TestSuccessionSettlesWholeBatchAndRetainsFindings(t *testing.T) {
	root, dir, run := fixture(t, quota.NewExhaustionPolicy(nil, nil))
	k, _ := quota.ReadKickoff(dir)
	k.Roles.Organizer = "a"
	os.Remove(filepath.Join(dir, quota.KickoffFile))
	if err := quota.WriteKickoff(dir, *k); err != nil {
		t.Fatal(err)
	}
	m, _ := runmanifest.Load(root, run)
	m.QuotaKickoff = k
	m.RosterSnapshot = []runmanifest.RosterSnapshotEntry{{Agent: "a", Model: "a-model"}, {Agent: "b", Model: "b-model"}, {Agent: "c", Model: "c-model"}, {Agent: "d", Model: "d-model"}}
	if err := runmanifest.Write(root, run, m); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	os.WriteFile(path, []byte(strings.Replace(string(raw), "status:", "organizer: a\nimplementer: a\nauto_implement: true\nstatus:", 1)), 0600)
	validRound(t, dir, filepath.Base(dir), "c")
	reviewDir := filepath.Join(dir, "review", "round-01")
	os.MkdirAll(reviewDir, 0700)
	finding := "---\nagent: b\nidea: " + filepath.Base(dir) + "\nreview-round: 1\n---\n\n## Findings\n### [MAJOR] Preserve this finding\nConcrete defect remains binding.\n"
	os.WriteFile(filepath.Join(reviewDir, "b.md"), []byte(finding), 0600)
	partial := "---\nidea: " + filepath.Base(dir) + "\ndrafted-by: a\ndate: 2026-10-29\n---\n\n## Agreed decisions\nPartial unsigned draft.\n"
	os.WriteFile(filepath.Join(dir, "consensus.md"), []byte(partial), 0600)
	e1, e2 := quotatest.Exhaustion(filepath.Base(dir), "a", "consensus"), quotatest.Exhaustion(filepath.Base(dir), "d", "round-02")
	ctx := leaseFixture(t, dir, run)
	b, err := Settle(ctx, root, dir, run, "consensus", []string{"a", "d"}, []quota.Member{{ID: "a", Evidence: &e1}, {ID: "d", Evidence: &e2}})
	if err != nil || b == nil {
		t.Fatalf("settlement: %+v %v", b, err)
	}
	if !reflect.DeepEqual(b.Decision.After, []string{"b", "c"}) || b.Decision.RoleTransition.Successor != "b" || len(b.Decision.DraftArchives) != 1 {
		t.Fatalf("%+v", b.Decision)
	}
	archive, _ := os.ReadFile(filepath.Join(dir, b.Decision.DraftArchives[0].Archive))
	if string(archive) != partial {
		t.Fatal("archive changed bytes")
	}
	if _, err := os.Stat(filepath.Join(dir, "consensus.md")); !os.IsNotExist(err) {
		t.Fatal("partial canonical draft remains")
	}
	if len(b.Retained) == 0 {
		t.Fatal("successor's earlier findings were lost")
	}
	r := b.Decision.RoleTransition
	handover := fmt.Sprintf("---\nidea: %s\ntransition: %s\nauthor: b\norganizer: b\nimplementer: b\nphase: consensus\nbase: none\nbase-sha256:\ntakeover-branch: %s\ntakeover-head: %s\ndate: 2026-10-29\nstatus: in-progress\n---\n\n## Takeover\nBoth roles transfer.\n\n## Resume state\nRead review/round-01/b.md and retain its finding.\n", filepath.Base(dir), r.ID, r.Branch, r.Head)
	os.MkdirAll(filepath.Join(dir, "role-handover"), 0700)
	os.WriteFile(filepath.Join(dir, r.Path()), []byte(handover), 0600)
	if ok, err := SingleReviewerAfterDropout(root, dir, run, "b", []string{"b", "c"}, nil); !ok || err != nil {
		t.Fatalf("review cause: %v %v", ok, err)
	}
	if _, err := Before(ctx, root, dir, run); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(filepath.Join(reviewDir, "b.md"))
	if string(after) != finding {
		t.Fatal("successor's filed review changed")
	}
	archive[0] = 'x'
	os.WriteFile(filepath.Join(dir, b.Decision.DraftArchives[0].Archive), archive, 0600)
	if _, err := Before(ctx, root, dir, run); err == nil {
		t.Fatal("corrupt archive accepted")
	}
}

func TestDraftArchiveReplayPreservesNewCanonical(t *testing.T) {
	dir := t.TempDir()
	raw := []byte("unsigned partial bytes\n")
	hash := fmt.Sprintf("%x", sha256.Sum256(raw))
	b := quota.Batch{Decision: quota.Decision{DraftArchives: []quota.DraftArchive{{Path: "consensus.md", Author: "a", SHA256: hash, Archive: "draft-archive/partial.md"}}}}
	os.WriteFile(filepath.Join(dir, "consensus.md"), raw, 0600)
	if err := reconcileDrafts(dir, b, false); err != nil {
		t.Fatal(err)
	}
	if err := reconcileDrafts(dir, b, false); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "consensus.md"), []byte("successor draft\n"), 0600)
	if err := reconcileDrafts(dir, b, true); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "consensus.md"))
	if string(got) != "successor draft\n" {
		t.Fatal("replay changed successor draft")
	}
}
