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
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/runmanifest"
)

func quotaConsensusFixture(t *testing.T) (string, protocol.IdeaStatus, context.Context) {
	t.Helper()
	return quotaConsensusFixturePolicy(t, quota.NewPolicy(nil, nil))
}

func quotaConsensusFixturePolicy(t *testing.T, p quota.Policy) (string, protocol.IdeaStatus, context.Context) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
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
	quote := "Owner ruling for " + ob.ID + ": retain the original evidence and adopt the alternative."
	authority := quotatest.Authority(t, root, idea.Slug, "ruling", quote)
	writeFile(t, path, raw+disposition+quotatest.Fields(authority))
	writeFile(t, filepath.Join(idea.Path, "round-02", "a.md"), "---\nagent: a\nidea: "+idea.Slug+"\n---\n## User direction\n"+quote+"\n")
	s, err = Status(root, idea.Slug, false)
	if err != nil || s.Triage != TriageReady || len(s.Signoffs) != 3 {
		t.Fatal(s, err)
	}
	after, _ := os.ReadFile(path)
	if !bytes.HasPrefix(after, []byte(raw)) {
		t.Fatal("filed veto rewritten")
	}
}

func TestQuotaFixupReincludeWithdrawThenAppendSignoff(t *testing.T) {
	root, idea, ctx := quotaConsensusFixture(t)
	path := filepath.Join(idea.Path, "consensus.md")
	raw := "---\nidea: " + idea.Slug + "\ndrafted-by: a\n---\n## Signoffs\n" + signoffBlock("a", "2026-10-04", StatusAccept, "", "") + signoffBlock("b", "2026-10-04", StatusAccept, "", "") + signoffBlock("c", "2026-10-04", StatusBlock, "Data loss.", "Retain records.")
	writeFile(t, path, raw)
	batch, e := membership.Settle(ctx, root, idea.Path, "test-run", "consensus", idea.Participants, quotaConsensusMembers())
	if e != nil {
		t.Fatal(e)
	}
	ids := []string{"a", "b", "c"}
	p := quota.NewPolicy(nil, nil)
	a := quotatest.Authority(t, root, idea.Slug, "reinclusion", (quota.RevisionDirective{Participants: ids, Policy: p}).Text())
	revision, e := membership.Revise(ctx, root, idea.Path, "test-run", ids, p, a, nil)
	if e != nil {
		t.Fatal(e)
	}
	if s, e := Status(root, idea.Slug, false); e != nil || s.Triage != TriageBlocked {
		t.Fatal(s, e)
	}
	ob := batch.Retained[0]
	disposition := "\n\n" + ob.ID + "\nDisposition: withdrawn\nRationale: Author read the corrected design after owner-confirmed return.\nAuthority: c\nEvidence: round-02/c.md\n"
	writeFile(t, path, raw+disposition)
	for _, binding := range []string{"", "batch-invented", revision.ID} {
		writeFile(t, filepath.Join(idea.Path, "round-02/c.md"), "---\nagent: c\nidea: "+idea.Slug+"\nquota-revision: 2\nreinclusion: "+binding+"\n---\nWithdrawal: "+ob.ID+"\n")
		s, e := Status(root, idea.Slug, false)
		if e != nil {
			t.Fatal(e)
		}
		if binding != revision.ID && s.Triage != TriageBlocked {
			t.Fatal("unbound withdrawal released veto", s)
		}
		if binding == revision.ID && (s.Triage == TriageBlocked || s.Triage == TriageReady) {
			t.Fatal("withdrawal either kept veto or counted as fresh signoff", s)
		}
	}
	s, e := AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "c", Status: "accept"})
	if e != nil || s.Triage != TriageReady {
		t.Fatal(s, e)
	}
	after, _ := os.ReadFile(path)
	if !bytes.HasPrefix(after, []byte(raw)) {
		t.Fatal("historical veto bytes changed")
	}
	// Re-inclusion authority survives permitted deletion from the inbox.
	os.Remove(filepath.Join(root, a.Path))
	if s, e = Status(root, idea.Slug, false); e != nil || s.Triage != TriageReady {
		t.Fatal(s, e)
	}
	// A later new veto is never disposed merely because an earlier veto was.
	writeFile(t, path, string(after)+signoffBlock("c", "2026-10-05", StatusBlock, "A new defect.", "Fix new defect."))
	s, e = Status(root, idea.Slug, false)
	if e != nil || s.Triage == TriageReady {
		t.Fatal("later veto swallowed", s, e)
	}
}

func TestQuotaFixupOffModeReturnCannotAuthorizeWithdrawal(t *testing.T) {
	root, idea, _ := quotaConsensusFixturePolicy(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	path := filepath.Join(idea.Path, "consensus.md")
	raw := "---\nidea: " + idea.Slug + "\ndrafted-by: a\n---\n## Signoffs\n" + signoffBlock("a", "2026-10-04", StatusAccept, "", "") + signoffBlock("b", "2026-10-04", StatusAccept, "", "") + signoffBlock("c", "2026-10-04", StatusBlock, "Old defect.", "Retain records.")
	writeFile(t, path, raw)
	promptPath := filepath.Join(idea.Path, "00-prompt.md")
	prompt, _ := os.ReadFile(promptPath)
	prompt = bytes.Replace(prompt, []byte("participants: [a, b, c, d]"), []byte("participants: [a, b]\nexcluded: [c — unavailable — confirmed 2026-10-04]\nexcluded: [d — unavailable — confirmed 2026-10-04]"), 1)
	os.WriteFile(promptPath, prompt, 0600)
	if e := membership.RecordManual(root, idea.Path); e != nil {
		t.Fatal(e)
	}
	h, e := quota.ReadHistory(idea.Path)
	if e != nil || len(h.Batches[0].Retained) != 1 {
		t.Fatal(h, e)
	}
	ob := h.Batches[0].Retained[0]
	prompt, _ = os.ReadFile(promptPath)
	prompt = bytes.Replace(prompt, []byte("participants: [a, b]"), []byte("participants: [a, b, c]"), 1)
	os.WriteFile(promptPath, prompt, 0600)
	if e = membership.RecordManual(root, idea.Path); e != nil {
		t.Fatal(e)
	}
	h, e = quota.ReadHistory(idea.Path)
	if e != nil {
		t.Fatal(e)
	}
	writeFile(t, filepath.Join(idea.Path, "round-02/c.md"), "---\nagent: c\nidea: "+idea.Slug+"\nquota-revision: 2\nreinclusion: "+h.Batches[1].ID+"\n---\nWithdrawal: "+ob.ID+"\n")
	writeFile(t, path, raw+"\n\n"+ob.ID+"\nDisposition: withdrawn\nRationale: Checked correction on return.\nAuthority: c\nEvidence: round-02/c.md\n")
	s, e := Status(root, idea.Slug, false)
	if e != nil || s.Triage != TriageBlocked {
		t.Fatal("manual return released veto", s, e)
	}
	if _, e = AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "c", Status: "accept"}); e == nil {
		t.Fatal("manual return permitted replacing retained veto")
	}

	after, _ := os.ReadFile(path)
	if !bytes.HasPrefix(after, []byte(raw)) {
		t.Fatal("old veto changed")
	}
	t.Log("ordinary off-mode return imported without extra record; manual revision cannot authorize retained-veto withdrawal; original bytes preserved")
}

func TestQuotaFixupWithdrawalMatchesBlankLineSignoffAndNotNewVeto(t *testing.T) {
	root, idea, ctx := quotaConsensusFixture(t)
	path := filepath.Join(idea.Path, "consensus.md")
	veto := strings.Replace(signoffBlock("c", "2026-10-04", StatusBlock, "Data loss.", "Retain records."), "\nStatus:", "\n\nStatus:", 1)
	raw := "---\nidea: " + idea.Slug + "\ndrafted-by: a\n---\n## Signoffs\n" + signoffBlock("a", "2026-10-04", StatusAccept, "", "") + signoffBlock("b", "2026-10-04", StatusAccept, "", "") + veto
	writeFile(t, path, raw)
	batch, e := membership.Settle(ctx, root, idea.Path, "test-run", "consensus", idea.Participants, quotaConsensusMembers())
	if e != nil {
		t.Fatal(e)
	}
	ids := []string{"a", "b", "c"}
	p := quota.NewPolicy(nil, nil)
	a := quotatest.Authority(t, root, idea.Slug, "blank-veto", (quota.RevisionDirective{Participants: ids, Policy: p}).Text())
	rev, e := membership.Revise(ctx, root, idea.Path, "test-run", ids, p, a, nil)
	if e != nil {
		t.Fatal(e)
	}
	ob := batch.Retained[0]
	writeFile(t, filepath.Join(idea.Path, "round-02/c.md"), "---\nagent: c\nidea: "+idea.Slug+"\nquota-revision: 2\nreinclusion: "+rev.ID+"\n---\nWithdrawal: "+ob.ID+"\n")
	disposition := "\n\n" + ob.ID + "\nDisposition: withdrawn\nRationale: Returned author withdraws.\nAuthority: c\nEvidence: round-02/c.md\n"
	writeFile(t, path, raw+disposition)
	s, e := AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "c", Status: "accept"})
	if e != nil || s.Triage != TriageReady {
		t.Fatal(s, e)
	}
	after, _ := os.ReadFile(path)
	// An identical signoff filed again on the same day is a new veto, not the old one.
	writeFile(t, path, string(after)+veto)
	s, e = Status(root, idea.Slug, false)
	if e != nil || s.Triage == TriageReady {
		t.Fatal("new identical veto swallowed", s, e)
	}
}
