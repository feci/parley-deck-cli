package app

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
)

func cycle3Run(t *testing.T, args ...string) (int, string) {
	t.Helper()
	var out, errs bytes.Buffer
	code := Run(args, &out, &errs)
	text := out.String() + errs.String()
	t.Logf("parley %q: exit=%d\n%s", args, code, text)
	return code, text
}
func cycle3Write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
}
func cycle3Edit(t *testing.T, path, before, after string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Replace(string(raw), before, after, 1)
	if text == string(raw) {
		t.Fatalf("edit missed %q", before)
	}
	cycle3Write(t, path, text)
}
func cycle3Round(id, slug string) string {
	raw := fmt.Sprintf("---\nagent: %s\nidea: %s\nround: 1\n---\n", id, slug)
	for _, h := range []string{"Summary", "Proposed approach", "Concerns / open questions", "Risks", "Existing alternatives"} {
		raw += "\n## " + h + "\nLocal catch-up fixture: read priors and join from round 2.\n"
	}
	return raw
}
func cycle3Stub(t *testing.T, root, id, target, raw string, exit int) []string {
	t.Helper()
	script := filepath.Join(root, "stub-"+id)
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
	cycle3Write(t, script, "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo local-cycle3-stub; exit 0; fi\ncat >/dev/null\nprintf called > "+quote(filepath.Join(root, "called"))+"\ncat > "+quote(target)+" <<'CATCHUP'\n"+raw+"CATCHUP\nexit "+fmt.Sprint(exit)+"\n")
	if err := os.Chmod(script, 0700); err != nil {
		t.Fatal(err)
	}
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: id, Path: script})
	fixupIsolateDiscovery(t, root)
	prompt := filepath.Join(root, "catchup-prompt.txt")
	cycle3Write(t, prompt, "Read priors and complete your own late round-1 artifact.")
	return []string{"agents", "exec", "--dir", root, "--agent", id, "--prompt-file", prompt, "--artifact", target, "--yes"}
}
func cycle3Consensus(t *testing.T, idea protocol.IdeaStatus) {
	cycle3Write(t, filepath.Join(idea.Path, "consensus.md"), "---\nidea: "+idea.Slug+"\ndrafted-by: alpha\n---\n\n## Signoffs\n")
}
func cycle3Sign(t *testing.T, root, slug, id string) int {
	t.Helper()
	code, _ := cycle3Run(t, "consensus", "signoff", "--dir", root, "--agent", id, "--status", "accept", slug)
	return code
}

func TestQuotaCycle3CLICatchupBothOrdersAndIncompleteRetry(t *testing.T) {
	for _, order := range []string{"artifact-first", "edit-first", "retry-stub"} {
		t.Run(order, func(t *testing.T) {
			root, idea, _ := fixupAppFixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
			prompt := filepath.Join(idea.Path, "00-prompt.md")
			target := filepath.Join(idea.Path, "round-01/delta.md")
			cycle3Consensus(t, idea)
			if order != "artifact-first" {
				cycle3Edit(t, prompt, "[alpha, beta, gamma]", "[alpha, beta, gamma, delta]")
				for _, id := range idea.Participants {
					if cycle3Sign(t, root, idea.Slug, id) != 0 {
						t.Fatal("existing signer could not append")
					}
				}
				v, err := protocol.InspectQuota(idea.Path)
				if err != nil || len(v.Catchup) != 1 || !strings.Contains(v.Pending, "parley agents exec --agent delta") || membership.Has(v.History.Known, "delta") {
					t.Fatal(v, err)
				}
				if _, err := membership.Before(context.Background(), root, idea.Path, "recorded-run"); err == nil {
					t.Fatal("pending catch-up dispatched normal work")
				}
				s, err := consensus.Status(root, idea.Slug, false)
				if err != nil || s.Triage != consensus.TriagePartial || strings.Join(s.Missing, ",") != "delta" {
					t.Fatal(s, err)
				}
				before, _ := os.ReadFile(filepath.Join(idea.Path, "consensus.md"))
				if cycle3Sign(t, root, idea.Slug, "delta") == 0 {
					t.Fatal("unknown joiner signed before catch-up")
				}
				after, _ := os.ReadFile(filepath.Join(idea.Path, "consensus.md"))
				if !bytes.Equal(before, after) {
					t.Fatal("rejected CLI signoff wrote canonical bytes")
				}
				c, text := cycle3Run(t, "status", "--dir", root)
				if c != 0 || !strings.Contains(text, "pending catch-up") || strings.Contains(text, "integrity gate") {
					t.Fatal(text)
				}
			}
			args := cycle3Stub(t, root, "delta", target, cycle3Round("delta", idea.Slug), 0)
			if order == "retry-stub" {
				stub := "---\nagent: delta\nidea: " + idea.Slug + "\nround: 1\n---\n"
				args = cycle3Stub(t, root, "delta", target, stub, 1)
				if code, _ := cycle3Run(t, args...); code == 0 {
					t.Fatal("failed stub accepted")
				}
				raw, _ := os.ReadFile(target)
				if string(raw) != stub {
					t.Fatal("partial artifact lost")
				}
				h, err := quota.ReadHistory(idea.Path)
				if err != nil || h.Revision != 0 || membership.Has(h.Known, "delta") {
					t.Fatal(h, err)
				}
				args = cycle3Stub(t, root, "delta", target, cycle3Round("delta", idea.Slug), 0)
			}
			if code, _ := cycle3Run(t, args...); code != 0 {
				t.Fatal("ordinary CLI catch-up failed")
			}
			if err := protocol.ValidateParticipantRoundArtifact(target, "delta", idea.Slug, 1); err != nil {
				t.Fatal(err)
			}
			if order == "artifact-first" {
				h, err := quota.ReadHistory(idea.Path)
				if err != nil || h.Revision != 0 || membership.Has(h.Known, "delta") {
					t.Fatal("artifact alone imported membership", h, err)
				}
				cycle3Edit(t, prompt, "[alpha, beta, gamma]", "[alpha, beta, gamma, delta]")
				for _, id := range idea.Participants {
					if cycle3Sign(t, root, idea.Slug, id) != 0 {
						t.Fatal("import signoff failed")
					}
				}
			}
			h, err := quota.ReadHistory(idea.Path)
			if err != nil || h.Revision != 1 || !membership.Has(h.Known, "delta") || h.Batches[0].Owner.Authority != nil || membership.Has(h.Kickoff.Participants, "delta") {
				t.Fatal(h, err)
			}
			v, err := protocol.InspectQuota(idea.Path)
			if err != nil || v.Pending != "" {
				t.Fatal(v, err)
			}
			if cycle3Sign(t, root, idea.Slug, "delta") != 0 {
				t.Fatal("completed joiner could not sign")
			}
			s, err := consensus.Status(root, idea.Slug, false)
			if err != nil || s.Triage != consensus.TriageReady {
				t.Fatal(s, err)
			}
		})
	}
}

func TestQuotaCycle3CLIKickoffExcludedReturn(t *testing.T) {
	for _, status := range []string{"round-01", "round-02"} {
		t.Run(status, func(t *testing.T) {
			root := t.TempDir()
			t.Setenv("PARLEY_HOME", t.TempDir())
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			declareAppTestSource(t, root)
			p := quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea}
			idea, k, err := protocol.CreateIdeaWithQuota(root, "kickoff return", []string{"alpha", "beta"}, []string{"delta — unavailable — confirmed 2026-10-04"}, "deliberation", "", "run-1", &p, nil)
			if err != nil {
				t.Fatal(err)
			}
			if err = runmanifest.Write(root, "run-1", runmanifest.New(runmanifest.Options{Root: root, RunID: "run-1", IdeaSlug: idea.Slug, Participants: idea.Participants, QuotaKickoff: k})); err != nil {
				t.Fatal(err)
			}
			historyBefore, _ := os.ReadFile(filepath.Join(idea.Path, quota.KickoffFile))
			h, err := quota.ReadHistory(idea.Path)
			if err != nil || membership.Has(h.Known, "delta") {
				t.Fatal("kickoff exclusion known", h, err)
			}
			cycle3Consensus(t, idea)
			if status != "round-01" {
				cycle3Edit(t, filepath.Join(idea.Path, "00-prompt.md"), "status: round-01", "status: "+status)
			}
			cycle3Edit(t, filepath.Join(idea.Path, "00-prompt.md"), "[alpha, beta]", "[alpha, beta, delta]")
			if cycle3Sign(t, root, idea.Slug, "alpha") != 0 {
				t.Fatal("existing signer rejected")
			}
			if status == "round-02" {
				if cycle3Sign(t, root, idea.Slug, "delta") == 0 {
					t.Fatal("later return skipped catch-up")
				}
				args := cycle3Stub(t, root, "delta", filepath.Join(idea.Path, "round-01/delta.md"), cycle3Round("delta", idea.Slug), 0)
				if code, _ := cycle3Run(t, args...); code != 0 {
					t.Fatal("later return could not catch up")
				}
			}
			if cycle3Sign(t, root, idea.Slug, "delta") != 0 {
				t.Fatal("returned signer rejected")
			}
			h, err = quota.ReadHistory(idea.Path)
			if err != nil || h.Revision != 1 || !membership.Has(h.Known, "delta") || membership.Has(h.Kickoff.Participants, "delta") || h.Batches[0].Owner.Authority != nil {
				t.Fatal(h, err)
			}
			after, _ := os.ReadFile(filepath.Join(idea.Path, quota.KickoffFile))
			if !bytes.Equal(after, historyBefore) {
				t.Fatal("kickoff rewritten")
			}
		})
	}
}

func TestQuotaCycle3CLICatchupBoundaries(t *testing.T) {
	for _, mode := range []string{"policy-on", "policy-on-edit", "foreign-path", "later-round", "known-excluded", "kickoff-excluded", "closed", "wrong-idea", "symlink", "foreign-stub", "complete-artifact"} {
		t.Run(mode, func(t *testing.T) {
			p := quota.Policy{Enabled: strings.HasPrefix(mode, "policy-on"), Scope: quota.KickoffAndMidIdea}
			root, idea, _ := fixupAppFixture(t, p)
			id := "delta"
			target := filepath.Join(idea.Path, "round-01/delta.md")
			prompt := filepath.Join(idea.Path, "00-prompt.md")
			switch mode {
			case "policy-on-edit":
				cycle3Edit(t, prompt, "[alpha, beta, gamma]", "[alpha, beta, gamma, delta]")
			case "foreign-path":
				target = filepath.Join(idea.Path, "round-01/stranger.md")
			case "later-round":
				target = filepath.Join(idea.Path, "round-02/delta.md")
			case "known-excluded":
				id = "gamma"
				target = filepath.Join(idea.Path, "round-01/gamma.md")
				cycle3Edit(t, prompt, "participants: [alpha, beta, gamma]", "participants: [alpha, beta]\nexcluded: gamma — unavailable — confirmed 2026-10-04")
				if err := membership.RecordManual(root, idea.Path); err != nil {
					t.Fatal(err)
				}
			case "kickoff-excluded":
				cycle3Edit(t, prompt, "status:", "excluded: delta — unavailable — confirmed 2026-10-04\nstatus:")
			case "closed":
				cycle3Edit(t, prompt, "status: round-01", "status: final")
			case "symlink":
				if err := os.Remove(filepath.Join(idea.Path, "round-01")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(t.TempDir(), filepath.Join(idea.Path, "round-01")); err != nil {
					t.Fatal(err)
				}
			case "foreign-stub":
				cycle3Write(t, target, "---\nagent: alpha\nidea: "+idea.Slug+"\nround: 1\n---\n")
			case "complete-artifact":
				cycle3Write(t, target, cycle3Round(id, idea.Slug))
			}
			args := cycle3Stub(t, root, id, target, cycle3Round(id, idea.Slug), 0)
			if mode == "wrong-idea" {
				args = append(args, "--idea", "unrelated")
			}
			if code, _ := cycle3Run(t, args...); code == 0 {
				t.Fatal("forbidden launch accepted")
			}
			if _, err := os.Stat(filepath.Join(root, "called")); !os.IsNotExist(err) {
				t.Fatal("forbidden process dispatched", err)
			}
		})
	}
}

func TestQuotaCycle3MultiplePendingJoinersKeepAtomicHistory(t *testing.T) {
	root, idea, _ := fixupAppFixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	cycle3Consensus(t, idea)
	cycle3Edit(t, filepath.Join(idea.Path, "00-prompt.md"), "[alpha, beta, gamma]", "[alpha, beta, gamma, delta, echo]")
	for _, id := range []string{"delta", "echo"} {
		args := cycle3Stub(t, root, id, filepath.Join(idea.Path, "round-01", id+".md"), cycle3Round(id, idea.Slug), 0)
		if code, _ := cycle3Run(t, args...); code != 0 {
			t.Fatal("catch-up failed")
		}
		if id == "delta" {
			h, err := quota.ReadHistory(idea.Path)
			if err != nil || h.Revision != 0 || membership.Has(h.Known, "delta") {
				t.Fatal("partial batch imported", h, err)
			}
			before, _ := os.ReadFile(filepath.Join(idea.Path, "consensus.md"))
			if cycle3Sign(t, root, idea.Slug, "delta") == 0 {
				t.Fatal("uncommitted new signer accepted")
			}
			after, _ := os.ReadFile(filepath.Join(idea.Path, "consensus.md"))
			if !bytes.Equal(before, after) {
				t.Fatal("uncommitted signer wrote")
			}
			if cycle3Sign(t, root, idea.Slug, "alpha") != 0 {
				t.Fatal("existing signer blocked")
			}
		}
	}
	h, err := quota.ReadHistory(idea.Path)
	if err != nil || h.Revision != 1 || !membership.Has(h.Known, "delta") || !membership.Has(h.Known, "echo") {
		t.Fatal(h, err)
	}
}

func TestQuotaCycle3CLINonParticipantDecline(t *testing.T) {
	root, idea, _ := fixupAppFixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	cycle3Consensus(t, idea)
	cycle3Edit(t, filepath.Join(idea.Path, "00-prompt.md"), "[alpha, beta, gamma]", "[alpha, beta, gamma, delta]")
	code, text := cycle3Run(t, "consensus", "signoff", "--dir", root, "--agent", "delta", "--status", "block", "--notes", "❌ NON-PARTICIPANT", "--counter", "Continue without me", idea.Slug)
	if code != 0 || !strings.Contains(text, "Consensus: blocked") {
		t.Fatal("baseline decline route failed", text)
	}
	for _, id := range idea.Participants {
		if cycle3Sign(t, root, idea.Slug, id) != 0 {
			t.Fatal("decline blocked existing signer")
		}
	}
	s, err := consensus.Status(root, idea.Slug, false)
	if err != nil || s.Triage != consensus.TriageBlocked || len(s.Errors) > 0 || strings.Join(s.Missing, ",") != "delta" {
		t.Fatal(s, err)
	}
	h, err := quota.ReadHistory(idea.Path)
	if err != nil || h.Revision != 0 || membership.Has(h.Known, "delta") {
		t.Fatal(h, err)
	}
}
