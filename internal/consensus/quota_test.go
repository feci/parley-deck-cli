package consensus

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
)

func quotaConsensusFixture(t *testing.T) (string, protocol.IdeaStatus, context.Context) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	p := quota.NewPolicy(nil, nil)
	idea, k, err := protocol.CreateIdeaWithQuota(root, "quota consensus", []string{"a", "b", "c", "d"}, nil, "deliberation", "", "test-run", &p, nil)
	if err != nil {
		t.Fatal(err)
	}
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: "test-run", IdeaSlug: idea.Slug, Participants: idea.Participants, QuotaKickoff: k})
	if err = runmanifest.Write(root, "test-run", m); err != nil {
		t.Fatal(err)
	}
	ctx, release, err := membership.Acquire(context.Background(), idea.Path, "test-run")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(release)
	return root, idea, ctx
}
func quotaConsensusMembers() []quota.Member {
	now := time.Date(2026, 10, 4, 0, 0, 0, 0, time.UTC)
	reset := now.Add(2 * time.Hour)
	c := quota.Evidence{Eligible: true, InvocationID: "inv-c", Adapter: "synthetic-test", Provenance: "synthetic-test", RuleID: "test", ObservedAt: now, ResetAt: &reset, Excerpt: "weekly limit exhausted"}
	d := c
	d.InvocationID = "inv-d"
	return []quota.Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, {ID: "c", Evidence: &c}, {ID: "d", Evidence: &d}}
}
func TestQuotaExcludedVetoRemainsBlockedKnownAndAppendGateNarrow(t *testing.T) {
	root, idea, ctx := quotaConsensusFixture(t)
	path := filepath.Join(idea.Path, "consensus.md")
	raw := "---\nidea: " + idea.Slug + "\ndrafted-by: a\n---\n\n## Signoffs\n" + signoffBlock("a", "2026-10-04", StatusAccept, "", "") + signoffBlock("b", "2026-10-04", StatusAccept, "", "") + signoffBlock("c", "2026-10-04", StatusBlock, "The design loses data.", "Preserve all records.")
	writeFile(t, path, raw)
	writeFile(t, filepath.Join(idea.Path, "round-01", "c.md"), "---\nagent: c\n---\n### [ major ] Durable records are missing\nClaim C1 is DISPUTED\nNo DISPUTED claims in C2.\n")
	before, _ := os.ReadFile(path)
	b, err := membership.Settle(ctx, root, idea.Path, "test-run", "consensus", idea.Participants, quotaConsensusMembers())
	if err != nil || b == nil {
		t.Fatal(b, err)
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatal("signoffs changed")
	}
	s, err := Status(root, idea.Slug, false)
	if err != nil || s.Triage != TriageBlocked || len(s.Errors) > 0 || strings.Join(s.Participants, ",") != "a,b" || len(s.Retained) != 3 {
		t.Fatalf("%+v %v", s, err)
	}
	if _, err = AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "d", Status: "accept"}); err == nil {
		t.Fatal("excluded signer appended")
	}
	// A rewritten synthesis cannot silently dispose of the frozen veto/findings.
	writeFile(t, path, "---\nidea: "+idea.Slug+"\ndrafted-by: a\n---\n## Signoffs\n"+signoffBlock("a", "2026-10-04", StatusAccept, "", "")+signoffBlock("b", "2026-10-04", StatusAccept, "", ""))
	s, err = Status(root, idea.Slug, false)
	if err != nil || s.Triage != TriageBlocked || len(s.Retained) != 3 {
		t.Fatalf("%+v %v", s, err)
	}
}

func TestQuotaPendingBlocksSignoffAndCloseWithoutRewritingArtifacts(t *testing.T) {
	root, idea, ctx := quotaConsensusFixture(t)
	path := filepath.Join(idea.Path, "consensus.md")
	writeFile(t, path, "---\nidea: "+idea.Slug+"\ndrafted-by: a\n---\n## Signoffs\n")
	h, err := quota.ReadHistory(idea.Path)
	if err != nil {
		t.Fatal(err)
	}
	d := quota.Evaluate(h.Kickoff.Policy, h.Current, quotaConsensusMembers(), quota.Roles{})
	b := quota.NewBatch(h, "test-run", "round-01", idea.Participants, d, time.Now())
	if err = quota.CommitBatch(idea.Path, b); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	if _, err = Status(root, idea.Slug, false); err == nil {
		t.Fatal("pending evaluated signoffs")
	}
	if _, err = AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "a", Status: "accept"}); err == nil {
		t.Fatal("pending appended")
	}
	if _, _, err = Finalize(root, idea.Slug, FinalizeOptions{By: "a"}); err == nil {
		t.Fatal("competing close bypassed lifetime lock")
	}
	// Recovery also gates on incomplete survivors; it cannot mint a completed round.
	if _, err = membership.Before(ctx, root, idea.Path, "test-run"); err == nil {
		t.Fatal("invalid survivor accepted")
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatal("pending changed signoffs")
	}
}

func TestQuotaKickoffExcludedIdentityNeverKnown(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	p := quota.NewPolicy(nil, nil)
	idea, _, err := protocol.CreateIdeaWithQuota(root, "kickoff known", []string{"a", "b"}, []string{"c — unavailable — confirmed 2026-10-04"}, "deliberation", "", "kickoff", &p, nil)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(idea.Path, "consensus.md"), "---\nidea: "+idea.Slug+"\ndrafted-by: a\n---\n## Signoffs\n"+signoffBlock("c", "2026-10-04", StatusBlock, "Veto.", "Alternative."))
	s, err := Status(root, idea.Slug, false)
	if err != nil || s.Triage != TriageMalformed || !strings.Contains(strings.Join(s.Errors, ";"), "unknown participant c") {
		t.Fatalf("%+v %v", s, err)
	}
}

func TestQuotaRetainedVetoOnlyOwnerRulingCanDispose(t *testing.T) {
	root, idea, ctx := quotaConsensusFixture(t)
	path := filepath.Join(idea.Path, "consensus.md")
	raw := "---\nidea: " + idea.Slug + "\ndrafted-by: a\n---\n## Signoffs\n" + signoffBlock("a", "2026-10-04", StatusAccept, "", "") + signoffBlock("b", "2026-10-04", StatusAccept, "", "") + signoffBlock("c", "2026-10-04", StatusBlock, "Data loss.", "Retain records.")
	writeFile(t, path, raw)
	batch, err := membership.Settle(ctx, root, idea.Path, "test-run", "consensus", idea.Participants, quotaConsensusMembers())
	if err != nil {
		t.Fatal(err)
	}
	if len(batch.Retained) != 1 {
		t.Fatal(batch.Retained)
	}
	ob := batch.Retained[0]
	disposition := "\n\n" + ob.ID + "\nDisposition: operator-ruling\nRationale: Owner selected the documented alternative.\nAuthority: owner\nEvidence: round-02/a.md\n"
	writeFile(t, path, raw+disposition)
	s, err := Status(root, idea.Slug, false)
	if err != nil || s.Triage != TriageBlocked {
		t.Fatal(s, err)
	}
	writeFile(t, filepath.Join(idea.Path, "round-02", "a.md"), "---\nagent: a\n---\n## User direction\nOwner ruling for "+ob.ID+": retain the original evidence and adopt the alternative.\n")
	s, err = Status(root, idea.Slug, false)
	if err != nil || s.Triage != TriageReady || len(s.Signoffs) != 3 {
		t.Fatal(s, err)
	}
	after, _ := os.ReadFile(path)
	if !bytes.HasPrefix(after, []byte(raw)) {
		t.Fatal("filed veto rewritten")
	}
}
