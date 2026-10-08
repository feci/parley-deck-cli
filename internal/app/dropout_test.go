package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/telemetry"
)

func dropoutAppRecords(t *testing.T, root string) []telemetry.Record {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	var records []telemetry.Record
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var r telemetry.Record
		if err := json.Unmarshal(raw, &r); err != nil {
			t.Fatal(err)
		}
		if r.Metadata.ParticipantStep != "" {
			records = append(records, r)
		}
	}
	return records
}

// Delete only the last receipt to simulate terminal telemetry surviving a
// driver crash before validation, without fabricating a child's outcome.
func forgetLastDropoutValidation(t *testing.T, root string, records []telemetry.Record) telemetry.Record {
	t.Helper()
	if len(records) == 0 {
		t.Fatal("no actual participant invocation")
	}
	last := records[0]
	for _, r := range records[1:] {
		if r.RequestedAt.After(last.RequestedAt) {
			last = r
		}
	}
	if err := os.Remove(filepath.Join(root, ".parley-runtime", "invocations", last.InvocationID, "participant-result.json")); err != nil {
		t.Fatal(err)
	}
	return last
}

func dropoutGoalOps(t *testing.T, script string) driverImplOps {
	t.Helper()
	root, idea, run := fixupAppFixture(t, quota.NewParticipantPolicy(nil, nil))
	agent := agents.Discovery{Spec: agents.Spec{ID: "gamma", HeadlessArgs: []string{"-c", script}, PromptMode: agents.PromptStdin}, Found: true, Path: "/bin/sh"}
	o := newOpsFor(root, idea.Path, []agents.Discovery{agent}, "alpha", []string{"gamma"})
	o.ideaSlug, o.drafter = idea.Slug, "gamma"
	o.base = runner.Options{Root: root, RunID: run, Idea: idea, Agents: []agents.Discovery{agent}}
	return o
}

func dropoutGoalCheck(t *testing.T, o driverImplOps) (bool, string) {
	t.Helper()
	ctx, release, err := membership.Acquire(context.Background(), o.ideaDir, o.base.RunID)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	return o.GoalCheck(ctx)
}

func TestDropoutGoalCheckReplaysActualOutputAcrossRuns(t *testing.T) {
	for _, verdict := range []string{"PASS", "FAIL", "malformed"} {
		for _, beforePreservation := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/before-preservation=%t", verdict, beforePreservation), func(t *testing.T) {
				o := dropoutGoalOps(t, "cat >/dev/null; printf 'GOAL-CHECK: "+verdict+"\\n'")
				ok, detail := dropoutGoalCheck(t, o)
				if ok != (verdict == "PASS") || verdict == "FAIL" && !strings.Contains(detail, "GOAL-CHECK: FAIL") {
					t.Fatalf("first result=%v %q", ok, detail)
				}
				records := dropoutAppRecords(t, o.root)
				want := 1
				if verdict == "malformed" {
					want = 2
				}
				if len(records) != want {
					t.Fatalf("children=%d want=%d", len(records), want)
				}
				last := forgetLastDropoutValidation(t, o.root, records)
				if beforePreservation {
					if err := os.RemoveAll(filepath.Join(o.root, ".parley-runtime", "invocations", last.InvocationID, "participant-files")); err != nil {
						t.Fatal(err)
					}
				}
				o.base.RunID = "different-driver-run"
				// A new run's plausible but contradictory log is not this attempt's
				// evidence, even when the private copy never reached disk.
				newLog := filepath.Join(o.root, "parley-deck", "runs", o.base.RunID, "agents", "gamma", "goal-check.stdout.log")
				if err := os.MkdirAll(filepath.Dir(newLog), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(newLog, []byte("GOAL-CHECK: FAIL\n"), 0600); err != nil {
					t.Fatal(err)
				}
				again, replayDetail := dropoutGoalCheck(t, o)
				if again != ok || verdict == "FAIL" && !strings.Contains(replayDetail, "GOAL-CHECK: FAIL") || len(dropoutAppRecords(t, o.root)) != want {
					t.Fatalf("replay=%v %q children=%d", again, replayDetail, len(dropoutAppRecords(t, o.root)))
				}
			})
		}
	}
}

func TestDropoutGoalCheckValidEvidenceCannotChangeOrExcuseFailure(t *testing.T) {
	for _, kind := range []string{"changed-valid", "missing", "nonzero"} {
		t.Run(kind, func(t *testing.T) {
			script := "cat >/dev/null; printf 'GOAL-CHECK: PASS\\n'"
			if kind == "nonzero" {
				script += "; exit 7"
			}
			o := dropoutGoalOps(t, script)
			ok, detail := dropoutGoalCheck(t, o)
			if ok != (kind != "nonzero") {
				t.Fatalf("first=%v %q", ok, detail)
			}
			records := dropoutAppRecords(t, o.root)
			if len(records) != 1 {
				t.Fatal("valid output authorized retry", records)
			}
			retained := filepath.Join(o.root, ".parley-runtime", "invocations", records[0].InvocationID, "participant-files", "00-goal-check.stdout.log")
			switch kind {
			case "changed-valid":
				if err := os.WriteFile(retained, []byte("GOAL-CHECK: FAIL\n"), 0600); err != nil {
					t.Fatal(err)
				}
			case "missing":
				for _, path := range []string{retained, filepath.Join(o.root, "parley-deck", "runs", o.base.RunID, "agents", "gamma", "goal-check.stdout.log")} {
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				}
			}
			o.base.RunID = "new-run"
			if ok, detail := dropoutGoalCheck(t, o); ok || len(dropoutAppRecords(t, o.root)) != 1 {
				t.Fatalf("failure passed or replaced: %v %q", ok, detail)
			}
		})
	}
}

func TestDropoutReadinessSameBatchReplay(t *testing.T) {
	for _, kind := range []string{"PONG", "malformed", "changed-valid", "missing"} {
		t.Run(kind, func(t *testing.T) {
			root, _, _ := fixupAppFixture(t, quota.NewParticipantPolicy(nil, nil))
			p := quota.NewParticipantPolicy(nil, nil)
			text := "PONG"
			if kind == "malformed" {
				text = "I cannot return PONG"
			}
			a := agents.Discovery{Spec: agents.Spec{ID: "gamma", HeadlessArgs: []string{"-c", "cat >/dev/null; printf '" + text + "'"}, PromptMode: agents.PromptStdin}, Found: true, Path: "/bin/sh"}
			opts := preflightOptions{Root: root, ProbeID: "same-proposed-batch", QuotaPolicy: &p, PingTimeout: time.Second}
			first := checkRoster(context.Background(), opts, []agents.Discovery{a})[0]
			records := dropoutAppRecords(t, root)
			want := 1
			if kind == "malformed" {
				want = 2
			}
			if first.Available != (kind != "malformed") || len(records) != want {
				t.Fatalf("first=%+v children=%d", first, len(records))
			}
			if kind == "PONG" || kind == "malformed" {
				forgetLastDropoutValidation(t, root, records)
			} else {
				path := filepath.Join(root, ".parley-runtime", "invocations", records[0].InvocationID, "probe.stdout.log")
				if kind == "missing" {
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				} else if err := os.WriteFile(path, []byte("PONG\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			again := checkRoster(context.Background(), opts, []agents.Discovery{a})[0]
			if again.Available != (kind == "PONG") || len(dropoutAppRecords(t, root)) != want {
				t.Fatalf("replay=%+v children=%d", again, len(dropoutAppRecords(t, root)))
			}
			if kind == "malformed" && (again.QuotaEvidence == nil || len(again.QuotaEvidence.Failure.Attempts) != 2) {
				t.Fatal("lost paired malformed-output proof", again)
			}
			if (kind == "missing" || kind == "changed-valid") && again.QuotaEvidence != nil {
				t.Fatal("integrity stop became dropout", again)
			}
		})
	}
}

func TestDropoutSignoffInvalidOwnSuffixAndTampering(t *testing.T) {
	for _, kind := range []string{"invalid-own", "tampering", "block"} {
		t.Run(kind, func(t *testing.T) {
			root, idea, _ := fixupAppFixture(t, quota.NewParticipantPolicy(nil, nil))
			promptPath := filepath.Join(idea.Path, "00-prompt.md")
			prompt, _ := os.ReadFile(promptPath)
			writeConsensusIdea(t, root, idea.Slug, idea.Participants, false, nil)
			os.WriteFile(promptPath, prompt, 0600)
			path := filepath.Join(idea.Path, "consensus.md")
			raw, _ := os.ReadFile(path)
			raw = bytes.Replace(raw, []byte("---\n"), []byte("---\ndrafted-by: alpha\n"), 1)
			os.WriteFile(path, raw, 0600)
			bin := t.TempDir()
			alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
			beta := writeFakeSignoffCLI(t, bin, "beta", "accept", 0)
			gamma := filepath.Join(bin, "gamma")
			target := "'" + strings.ReplaceAll(path, "'", "'\\''") + "'"
			script := "#!/bin/sh\ncat >/dev/null\nprintf '\\n### Signoff: gamma — 2026-10-08\\n\\nStatus: nonsense\\n' >> " + target + "\nexit 7\n"
			if kind == "tampering" {
				script = "#!/bin/sh\ncat >/dev/null\nprintf 'overwrote all signoffs' > " + target + "\nexit 7\n"
			}
			if kind == "block" {
				gamma = writeFakeSignoffCLI(t, bin, "gamma", "block", 1)
			} else {
				os.WriteFile(gamma, []byte(script), 0700)
			}
			writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha}, fakeAgentConfig{ID: "beta", Path: beta}, fakeAgentConfig{ID: "gamma", Path: gamma})
			var out, errs bytes.Buffer
			err := requestConsensusSignoffs(context.Background(), requestSignoffsOptions{Root: root, IdeaSlug: idea.Slug, Yes: true}, &out, &errs)
			h, he := quota.ReadHistory(idea.Path)
			if he != nil {
				t.Fatal(he)
			}
			if kind == "invalid-own" {
				if err != nil || !h.Dropped("gamma") {
					t.Fatalf("%v history=%+v\n%s\n%s", err, h, out.String(), errs.String())
				}
				summary, e := consensus.Status(root, idea.Slug, false)
				if e != nil || summary.Triage != consensus.TriageReady {
					t.Fatalf("%+v %v", summary, e)
				}
				copies, _ := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "participant-files", "00-consensus.md"))
				invalid := 0
				for _, p := range copies {
					b, _ := os.ReadFile(p)
					if strings.Contains(string(b), "Status: nonsense") {
						invalid++
					}
				}
				if invalid != 2 {
					t.Fatal("invalid suffix copies", invalid)
				}
			} else {
				if err == nil || h.Revision != 0 {
					t.Fatalf("%s did not stop: %+v %v", kind, h, err)
				}
				if kind == "block" {
					summary, e := consensus.Status(root, idea.Slug, false)
					if e != nil || summary.Triage != consensus.TriageBlocked {
						t.Fatalf("BLOCK lost: %+v %v", summary, e)
					}
				}
			}
			paths, _ := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
			count := 0
			for _, p := range paths {
				b, _ := os.ReadFile(p)
				var r telemetry.Record
				json.Unmarshal(b, &r)
				if r.Metadata.Agent == "gamma" {
					count++
				}
			}
			want := 1
			if kind == "invalid-own" {
				want = 2
			}
			if count != want {
				t.Fatalf("attempts=%d want=%d", count, want)
			}
		})
	}
}

func TestDropoutReadinessRetriesBeforeMembership(t *testing.T) {
	root, _, _ := fixupAppFixture(t, quota.NewParticipantPolicy(nil, nil))
	p := quota.NewParticipantPolicy(nil, nil)
	var discovered []agents.Discovery
	for _, id := range []string{"alpha", "beta", "gamma"} {
		script := "cat >/dev/null; printf PONG"
		if id == "gamma" {
			script = "cat >/dev/null; printf 'HTTP 503' >&2; exit 9"
		}
		discovered = append(discovered, agents.Discovery{Spec: agents.Spec{ID: id, HeadlessArgs: []string{"-c", script}, PromptMode: agents.PromptStdin}, Found: true, Path: "/bin/sh"})
	}
	start := time.Now()
	entries := checkRoster(context.Background(), preflightOptions{Root: root, QuotaPolicy: &p, PingTimeout: time.Second}, discovered)
	if time.Since(start) < 5*time.Second {
		t.Fatal("retry delay was skipped")
	}
	ms := []quota.Member{}
	for _, e := range entries {
		ms = append(ms, quota.Member{ID: e.RosterID, Usable: e.Available, Evidence: e.QuotaEvidence})
	}
	d := quota.Evaluate(p, []string{"alpha", "beta", "gamma"}, ms, quota.Roles{})
	if !d.Applied || len(d.Candidates) != 1 || len(d.Candidates[0].Evidence.Failure.Attempts) != 2 {
		t.Fatalf("%+v entries=%+v", d, entries)
	}
	// Standalone reporting has no policy and therefore no automatic retry/reduction.
	standalone := checkRoster(context.Background(), preflightOptions{Root: root, PingTimeout: time.Second}, discovered)
	if standalone[2].QuotaEvidence != nil && standalone[2].QuotaEvidence.Failure != nil {
		t.Fatal("standalone authorized dropout")
	}
}

func TestDropoutKickoffBlockCreatesSafeInbox(t *testing.T) {
	root := t.TempDir()
	d := quota.Decision{Before: []string{"a", "b", "c"}, UsableSurvivors: 1, Block: "usable floor"}
	if err := writeQuotaBlock(root, d); err != nil {
		t.Fatal(err)
	}
	files, _ := filepath.Glob(filepath.Join(root, "parley-deck", "inbox", "*.md"))
	if len(files) != 1 {
		t.Fatal(files)
	}
	raw, _ := os.ReadFile(files[0])
	if !strings.Contains(string(raw), quota.OwnerOptions) || !strings.Contains(string(raw), "fixed floor: 2") {
		t.Fatal(string(raw))
	}
	if err := writeQuotaBlock(root, d); err != nil {
		t.Fatal(err)
	}
	os.Remove(files[0])
	os.Remove(filepath.Dir(files[0]))
	os.Symlink(t.TempDir(), filepath.Dir(files[0]))
	if err := writeQuotaBlock(root, d); err == nil {
		t.Fatal(fmt.Sprint("unsafe inbox accepted"))
	}
}
