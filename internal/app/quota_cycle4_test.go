package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/store"
)

func cycle4Before(t *testing.T, root string, idea protocol.IdeaStatus, run string) {
	t.Helper()
	ctx, release, err := membership.Acquire(context.Background(), idea.Path, run)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if _, err := membership.Before(ctx, root, idea.Path, run); err != nil {
		t.Fatal(err)
	}
}

func cycle4Status(t *testing.T, root string) {
	t.Helper()
	code, text := cycle3Run(t, "status", "--dir", root)
	if code != 0 || strings.Contains(text, "integrity gate") || strings.Contains(text, "consensus=error") || strings.Contains(text, "quota transition pending") {
		t.Fatal(text)
	}
}

func TestQuotaCycle4CLINoticeRecoveryAndDispatch(t *testing.T) {
	for _, mode := range []string{"live-annotated", "archived-annotated", "live-replaced", "archived-replaced", "unsafe-file", "unsafe-inbox", "unsafe-archive"} {
		for _, beforeReceipt := range []bool{false, true} {
			label := "applied"
			if beforeReceipt {
				label = "before-receipt"
			}
			t.Run(mode+"/"+label, func(t *testing.T) {
				root, idea, run := fixupAppFixture(t, quota.NewPolicy(nil, nil))
				for _, id := range []string{"alpha", "beta"} {
					cycle3Write(t, filepath.Join(idea.Path, "round-01", id+".md"), cycle3Round(id, idea.Slug))
				}
				h, err := quota.ReadHistory(idea.Path)
				if err != nil {
					t.Fatal(err)
				}
				now := time.Now()
				reset := now.Add(2 * time.Hour)
				e := quota.Evidence{Eligible: true, InvocationID: "local-fixture-gamma", Adapter: "synthetic", RuleID: "synthetic", Provenance: "test", ObservedAt: now, ResetAt: &reset, Excerpt: "synthetic quota exhaustion"}
				d := quota.Evaluate(h.Policy(), h.Current, []quota.Member{{ID: "alpha", Usable: true}, {ID: "beta", Usable: true}, {ID: "gamma", Evidence: &e}}, quota.Roles{})
				b := quota.NewBatch(h, run, "round-01", h.Current, d, now)
				if err := quota.CommitBatch(idea.Path, b); err != nil {
					t.Fatal(err)
				}
				if !beforeReceipt {
					cycle4Before(t, root, idea, run)
				}
				inbox := filepath.Join(root, protocol.DeckDir, "inbox")
				path := filepath.Join(inbox, "parley-to-user_"+b.ID+".md")
				outside := t.TempDir()
				cycle3Write(t, filepath.Join(outside, "owner.txt"), "owner bytes")
				switch {
				case mode == "unsafe-file":
					if !beforeReceipt {
						if err := os.Remove(path); err != nil {
							t.Fatal(err)
						}
					}
					if err := os.Symlink(filepath.Join(outside, "owner.txt"), path); err != nil {
						t.Fatal(err)
					}
				case mode == "unsafe-inbox":
					if err := os.Rename(inbox, inbox+"-owner-saved"); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(outside, inbox); err != nil {
						t.Fatal(err)
					}
				case mode == "unsafe-archive":
					if !beforeReceipt {
						if err := os.Remove(path); err != nil {
							t.Fatal(err)
						}
					}
					if err := os.Symlink(outside, filepath.Join(inbox, "archived")); err != nil {
						t.Fatal(err)
					}
				default:
					if strings.HasPrefix(mode, "archived") {
						if !beforeReceipt {
							if err := os.Remove(path); err != nil {
								t.Fatal(err)
							}
						}
						path = filepath.Join(inbox, "archived", filepath.Base(path))
					}
					text := "Owner replacement\r\n"
					if strings.HasSuffix(mode, "annotated") {
						text = b.Notice() + "\n## Owner answer\nAcknowledged.\n"
					}
					cycle3Write(t, path, text)
				}
				historyPath := filepath.Join(idea.Path, quota.HistoryDir, "000001.json")
				history, err := os.ReadFile(historyPath)
				if err != nil {
					t.Fatal(err)
				}
				cycle3Consensus(t, idea)
				for i := 0; i < 2; i++ {
					cycle4Before(t, root, idea, run)
					cycle4Status(t, root)
				}
				if cycle3Sign(t, root, idea.Slug, "alpha") != 0 {
					t.Fatal("notice blocked actual signoff")
				}
				target := filepath.Join(idea.Path, "round-02/alpha.md")
				if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
					t.Fatal(err)
				}
				args := cycle3Stub(t, root, "alpha", target, strings.Replace(cycle3Round("alpha", idea.Slug), "round: 1", "round: 2", 1), 0)
				// Dispatch under the same authentic lifetime lease the driver carries.
				// Standalone policy-on exec without a lease intentionally remains gated.
				ctx, release, err := membership.Acquire(context.Background(), idea.Path, run)
				if err != nil {
					t.Fatal(err)
				}
				var out, errs bytes.Buffer
				code := runAgentsExec(ctx, append(args[2:], "--run-id", run), &out, &errs)
				release()
				t.Logf("leased agents exec: exit=%d stdout=%s stderr=%s", code, out.String(), errs.String())
				if code != 0 {
					t.Fatal("notice blocked leased local dispatch")
				}
				if err := protocol.ValidateParticipantRoundArtifact(target, "alpha", idea.Slug, 2); err != nil {
					t.Fatal(err)
				}
				after, err := os.ReadFile(historyPath)
				if err != nil || !bytes.Equal(history, after) {
					t.Fatal("history changed", err)
				}
				if ok, err := quota.ReadApplied(idea.Path, b.ID); !ok || err != nil {
					t.Fatal(ok, err)
				}
				events, err := store.New(filepath.Join(root, protocol.DeckDir, "runs", run)).Load()
				if err != nil {
					t.Fatal(err)
				}
				count := 0
				for _, e := range events {
					if e.Type == "round.completed" && e.Data["quota_transition"] == b.ID {
						count++
					}
				}
				if count != 1 {
					t.Fatal("terminal evaluations", count)
				}
				entries, err := os.ReadDir(outside)
				if err != nil || len(entries) != 1 {
					t.Fatal("unsafe target written", entries, err)
				}
			})
		}
	}
}

func TestQuotaCycle4CLIBoundOwnerAnswerAnnotations(t *testing.T) {
	for _, archived := range []bool{false, true} {
		label := "live"
		if archived {
			label = "archived"
		}
		t.Run(label, func(t *testing.T) {
			root, idea, run := fixupAppFixture(t, quota.NewPolicy(nil, nil))
			ids := []string{"alpha", "beta"}
			p := quota.NewPolicy(nil, nil)
			a := quotatest.Authority(t, root, idea.Slug, "bound-answer", (quota.RevisionDirective{Participants: ids, Policy: p}).Text())
			req, _ := json.Marshal(map[string]any{"participants": ids, "policy": p, "authority": a})
			request := filepath.Join(root, "request.json")
			cycle3Write(t, request, string(req))
			if code, _ := cycle3Run(t, "quota", "revise", "--dir", root, "--idea", idea.Slug, "--run", run, "--request", request); code != 0 {
				t.Fatal("revision failed")
			}
			historyPath := filepath.Join(idea.Path, quota.HistoryDir, "000001.json")
			history, err := os.ReadFile(historyPath)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, a.Path)
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if archived {
				dst := filepath.Join(filepath.Dir(path), "archived", filepath.Base(path))
				if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(path, dst); err != nil {
					t.Fatal(err)
				}
				path = dst
			}
			cycle3Write(t, path, string(raw)+"\n## Owner update\nKeep this answer in the audit trail.\n")
			cycle3Consensus(t, idea)
			cycle4Status(t, root)
			if cycle3Sign(t, root, idea.Slug, "alpha") != 0 {
				t.Fatal("owner append blocked signoff")
			}
			cycle4Before(t, root, idea, run)
			h, err := quota.ReadHistory(idea.Path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(filepath.Join(idea.Path, "quota-applied", h.Batches[0].ID)); err != nil {
				t.Fatal(err)
			}
			cycle4Before(t, root, idea, run)
			cycle4Status(t, root)
			if cycle3Sign(t, root, idea.Slug, "beta") != 0 {
				t.Fatal("owner append blocked second signoff")
			}
			after, err := os.ReadFile(historyPath)
			if err != nil || !bytes.Equal(history, after) {
				t.Fatal("bound history changed", err)
			}
			// Missing immutable evidence still fails actual status, signoff and Before.
			if err := os.Remove(filepath.Join(root, ".git/objects", a.Blob[:2], a.Blob[2:])); err != nil {
				t.Fatal(err)
			}
			_, text := cycle3Run(t, "status", "--dir", root)
			if !strings.Contains(text, "integrity gate") && !strings.Contains(text, "consensus=error") {
				t.Fatal(text)
			}
			if cycle3Sign(t, root, idea.Slug, "alpha") == 0 {
				t.Fatal("missing blob admitted signoff")
			}
			if _, err := membership.Before(context.Background(), root, idea.Path, run); err == nil {
				t.Fatal("missing blob admitted Before")
			}
		})
	}
}

func TestQuotaCycle4CLIRoundOnePlainEdits(t *testing.T) {
	for _, mode := range []string{"return-kept", "return-removed", "new"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			t.Setenv("PARLEY_HOME", t.TempDir())
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			declareAppTestSource(t, root)
			bin := t.TempDir()
			configs := []fakeAgentConfig{}
			for _, id := range []string{"alpha", "beta", "delta"} {
				writeFakeRoundAgentCLI(t, bin, id, "local-cycle4-stub")
				path := filepath.Join(bin, id)
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, bytes.Replace(raw, []byte("agent: codex"), []byte("agent: "+id), 1), 0700); err != nil {
					t.Fatal(err)
				}
				configs = append(configs, fakeAgentConfig{ID: id, Path: path})
			}
			writeAgentsLocalConfig(t, root, configs...)
			fixupIsolateDiscovery(t, root)
			withFakeProbe(t, func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) readinessObservation {
				if a.ID == "delta" {
					return procFailObs()
				}
				return readyObs()
			})
			participants := "alpha,beta,delta"
			if mode == "new" {
				participants = "alpha,beta"
			}
			if code, _ := cycle3Run(t, "run", "--dir", root, "--no-tui", "--no-auto", "--yes", "--quota-auto-exclude=false", "--participants", participants, "--track", "deliberation", "cycle4 round-one CLI kickoff"); code != 0 {
				t.Fatal("CLI kickoff failed")
			}
			ws, err := protocol.ReadWorkspaceStatus(root)
			if err != nil || len(ws.Ideas) != 1 {
				t.Fatal(ws, err)
			}
			idea := ws.Ideas[0]
			initial, err := quota.ReadHistory(idea.Path)
			if err != nil {
				t.Fatal(err)
			}
			run := initial.Kickoff.RunID

			kickoff, err := os.ReadFile(filepath.Join(idea.Path, quota.KickoffFile))
			if err != nil {
				t.Fatal(err)
			}
			h, err := quota.ReadHistory(idea.Path)
			if err != nil || membership.Has(h.Known, "delta") {
				t.Fatal("false kickoff membership", h, err)
			}
			cycle3Consensus(t, idea)
			prompt := filepath.Join(idea.Path, "00-prompt.md")
			cycle3Edit(t, prompt, "[alpha, beta]", "[alpha, beta, delta]")
			if mode == "return-removed" {
				raw, err := os.ReadFile(prompt)
				if err != nil {
					t.Fatal(err)
				}
				lines := strings.Split(string(raw), "\n")
				out := []string{}
				removed := false
				for _, line := range lines {
					if strings.HasPrefix(line, "excluded: delta — ") {
						removed = true
						continue
					}
					out = append(out, line)
				}
				if !removed {
					t.Fatal("kickoff exclusion marker missing")
				}
				cycle3Write(t, prompt, strings.Join(out, "\n"))
			}
			cycle4Status(t, root)
			if cycle3Sign(t, root, idea.Slug, "delta") != 0 {
				t.Fatal("round-one plain edit rejected")
			}
			cycle4Before(t, root, idea, run)
			h, err = quota.ReadHistory(idea.Path)
			if err != nil || h.Revision != 1 || !membership.Has(h.Known, "delta") || membership.Has(h.Kickoff.Participants, "delta") || h.Batches[0].Owner.Authority != nil || len(h.Batches[0].Owner.Catchup) != 0 {
				t.Fatal(h, err)
			}
			after, err := os.ReadFile(filepath.Join(idea.Path, quota.KickoffFile))
			if err != nil || !bytes.Equal(kickoff, after) {
				t.Fatal("kickoff changed", err)
			}
			// History must use the frozen round-one prompt after the live round advances.
			cycle3Edit(t, prompt, "status: round-01", "status: round-02")
			if _, err := quota.ReadHistory(idea.Path); err != nil {
				t.Fatal(err)
			}
			cycle4Before(t, root, idea, run)
		})
	}
}

func TestQuotaCycle4CLIPlainEditBoundaries(t *testing.T) {
	for _, mode := range []string{"policy-on", "final", "closed"} {
		t.Run(mode, func(t *testing.T) {
			root, idea, run := fixupAppFixture(t, quota.Policy{Enabled: mode == "policy-on", Scope: quota.KickoffAndMidIdea})
			cycle3Consensus(t, idea)
			prompt := filepath.Join(idea.Path, "00-prompt.md")
			if mode != "policy-on" {
				cycle3Edit(t, prompt, "status: round-01", "status: "+mode)
			}
			cycle3Edit(t, prompt, "[alpha, beta, gamma]", "[alpha, beta, gamma, delta]")
			if cycle3Sign(t, root, idea.Slug, "delta") == 0 {
				t.Fatal("forbidden plain edit signed")
			}
			if _, err := membership.Before(context.Background(), root, idea.Path, run); err == nil {
				t.Fatal("forbidden plain edit dispatched")
			}
			h, err := quota.ReadHistory(idea.Path)
			if err != nil || h.Revision != 0 || membership.Has(h.Known, "delta") {
				t.Fatal("boundary granted history", h, err)
			}
		})
	}
}
