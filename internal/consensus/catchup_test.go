package consensus

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
)

func pendingDeclineFixture(t *testing.T, policyOn bool) (string, protocol.IdeaStatus) {
	t.Helper()
	root, idea, _ := quotaConsensusFixturePolicy(t, quota.Policy{Enabled: policyOn, Scope: quota.KickoffAndMidIdea})
	prompt := filepath.Join(idea.Path, "00-prompt.md")
	raw, err := os.ReadFile(prompt)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, prompt, strings.Replace(strings.Replace(string(raw), "status: round-01", "status: round-02", 1), "[a, b, c, d]", "[a, b, c, d, e]", 1))
	for _, rel := range []string{"consensus.md", "review/consensus.md"} {
		writeFile(t, filepath.Join(idea.Path, rel), "---\nidea: "+idea.Slug+"\ndrafted-by: a\n---\n## Signoffs\n")
	}
	return root, idea
}

func TestQuotaCycle3PendingDeclineDoesNotGrantMembership(t *testing.T) {
	root, idea := pendingDeclineFixture(t, false)
	path := filepath.Join(idea.Path, "consensus.md")
	kickoff, _ := os.ReadFile(filepath.Join(idea.Path, quota.KickoffFile))
	s, err := AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "e", Status: "block", Notes: "❌ NON-PARTICIPANT", CounterProposal: "Continue without me"})
	if err != nil || s.Triage != TriageBlocked || len(s.Errors) > 0 || !contains(s.Missing, "e") {
		t.Fatal(s, err)
	}
	declineBytes, _ := os.ReadFile(path)
	if !bytes.Contains(declineBytes, []byte("Status: ❌ BLOCK\nNotes: ❌ NON-PARTICIPANT\nCounter-proposal (required if ❌): Continue without me")) {
		t.Fatal("decline changed its baseline CLI representation", string(declineBytes))
	}
	for _, id := range []string{"a", "b", "c", "d"} {
		s, err = AppendSignoff(root, idea.Slug, SignoffOptions{Agent: id, Status: "accept"})
		if err != nil || s.Triage != TriageBlocked || len(s.Errors) > 0 {
			t.Fatal("existing signer cannot append after decline", s, err)
		}
	}
	if strings.Join(s.Missing, ",") != "e" {
		t.Fatal("decline counted as completed quorum", s)
	}
	if _, _, err = Finalize(root, idea.Slug, FinalizeOptions{By: "a"}); err == nil {
		t.Fatal("pending decline finalized")
	}
	h, err := quota.ReadHistory(idea.Path)
	if err != nil || h.Revision != 0 || contains(h.Known, "e") || contains(h.Current, "e") || len(h.Batches) != 0 {
		t.Fatal("decline granted membership or owner authority", h, err)
	}
	afterKickoff, _ := os.ReadFile(filepath.Join(idea.Path, quota.KickoffFile))
	if !bytes.Equal(kickoff, afterKickoff) {
		t.Fatal("decline changed kickoff")
	}
	before, _ := os.ReadFile(path)
	if _, err = AppendSignoff(root, idea.Slug, SignoffOptions{Agent: "e", Status: "block", Notes: "❌ NON-PARTICIPANT", CounterProposal: "Continue without me"}); err == nil {
		t.Fatal("duplicate decline accepted")
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) || !bytes.HasPrefix(after, declineBytes) {
		t.Fatal("append rejection changed existing signoffs")
	}
	// A hand-written duplicate must be malformed on the same read path.
	writeFile(t, path, string(after)+signoffBlock("e", "2026-10-06", StatusBlock, "❌ NON-PARTICIPANT", "Continue without me"))
	s, err = Status(root, idea.Slug, false)
	if err != nil || s.Triage != TriageMalformed || !strings.Contains(strings.Join(s.Errors, ";"), "duplicate signoff") {
		t.Fatal(s, err)
	}
}

func TestQuotaCycle3PendingDeclineBoundaries(t *testing.T) {
	for _, mode := range []string{"accept", "reserve", "ordinary-block", "quoted-note", "mentioned-note", "missing-counter", "literal-status", "outsider", "review", "policy-on"} {
		t.Run(mode, func(t *testing.T) {
			root, idea := pendingDeclineFixture(t, mode == "policy-on")
			opts := SignoffOptions{Agent: "e", Status: "block", Notes: "❌ NON-PARTICIPANT", CounterProposal: "Continue without me"}
			switch mode {
			case "accept", "reserve":
				opts.Status = mode
			case "ordinary-block":
				opts.Notes = "Fix a defect"
			case "quoted-note":
				opts.Notes = "`❌ NON-PARTICIPANT`"
			case "mentioned-note":
				opts.Notes = "Example: ❌ NON-PARTICIPANT"
			case "missing-counter":
				opts.CounterProposal = ""
			case "literal-status":
				opts.Status = "❌ NON-PARTICIPANT"
			case "outsider":
				opts.Agent = "stranger"
			case "review":
				opts.Review = true
			}
			path := consensusPath(idea.Path, schemaFor(opts.Review))
			before, _ := os.ReadFile(path)
			if _, err := AppendSignoff(root, idea.Slug, opts); err == nil {
				t.Fatal("unsupported pending signoff accepted")
			}
			after, _ := os.ReadFile(path)
			if !bytes.Equal(before, after) {
				t.Fatal("rejected API call changed canonical bytes")
			}
			// Reading a manually filed invalid signoff must also stay closed.
			writeFile(t, path, string(before)+signoffBlock(opts.Agent, "2026-10-06", opts.Status, opts.Notes, opts.CounterProposal))
			s, err := Status(root, idea.Slug, opts.Review)
			if err == nil && s.Triage != TriageMalformed {
				t.Fatal("unsupported filed signoff accepted", s)
			}
		})
	}
}
