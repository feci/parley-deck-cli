package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestStepBudgetGroupedRunnerCountsOnceAndRetainsRefusals(t *testing.T) {
	for _, missingPolicy := range []bool{false, true} {
		t.Run(fmt.Sprintf("missing-policy=%v", missingPolicy), func(t *testing.T) {
			root := t.TempDir()
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			declareTestLaunchSource(t, root)
			idea, err := protocol.CreateIdea(root, "Grouped step fixture", []string{"alpha", "beta"})
			if err != nil {
				t.Fatal(err)
			}
			binding, err := budget.EnsureStepBinding(context.Background(), root, idea.Slug, 1, 0)
			if err != nil {
				t.Fatal(err)
			}
			var discovered []agents.Discovery
			for _, id := range idea.Participants {
				discovered = append(discovered, agents.Discovery{Spec: agents.Spec{ID: id,
					HeadlessArgs: []string{"-test.run=TestStepBudgetAgentHelper", "--", "parley-step-agent"},
					PromptMode:   agents.PromptStdin}, Path: os.Args[0], Found: true})
			}
			for attempt := 0; attempt < 2; attempt++ {
				if attempt > 0 {
					for _, id := range idea.Participants {
						if err := os.Remove(filepath.Join(idea.Path, "round-01", id+".md")); err != nil {
							t.Fatal(err)
						}
					}
					if missingPolicy {
						if err := os.Remove(filepath.Join(filepath.Dir(binding.Store.Dir), "policy.json")); err != nil {
							t.Fatal(err)
						}
					}
				}
				runID := fmt.Sprintf("run-%d", attempt)
				results := RunRoundOne(context.Background(), Options{Root: root, RunID: runID, Idea: idea,
					Task: "Grouped step fixture", Agents: discovered, Timeout: 5 * time.Second,
					Store: store.New(filepath.Join(root, protocol.DeckDir, "runs", runID))})
				if len(results) != 2 {
					t.Fatalf("group results: %+v", results)
				}
				for _, r := range results {
					if attempt == 0 && (!r.ArtifactOK || r.ExitError != "") {
						t.Fatalf("allowed group: %+v", r)
					}
					if attempt == 1 && r.ExitError == "" {
						t.Fatalf("new group escaped cap: %+v", r)
					}
				}
			}
			state, err := binding.Store.Inspect(context.Background())
			if err != nil || budget.StepCount(state) != 1 {
				t.Fatalf("group charged per participant: %+v %v", state, err)
			}
			records := terminalRecords(t, root)
			started, denied := 0, 0
			for _, r := range records {
				if r.StartedAt != nil {
					started++
				}
				if r.Outcome.FailureClass != nil && *r.Outcome.FailureClass == "budget_refused" {
					denied++
					if r.StartedAt != nil || r.PID != nil {
						t.Fatal("refused group spawned")
					}
				}
			}
			if len(records) != 4 || started != 2 || denied != 2 {
				t.Fatalf("records=%d started=%d denied=%d", len(records), started, denied)
			}
		})
	}
}

func TestStepBudgetAgentHelper(t *testing.T) {
	if !hasArg("parley-step-agent") {
		return
	}
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		t.Fatal(err)
	}
	match := regexp.MustCompile(`(?m)^- Create exactly this file and no other protocol artifact: (.+)$`).FindStringSubmatch(string(input))
	idea := regexp.MustCompile(`(?m)^idea: (.+)$`).FindStringSubmatch(string(input))
	if len(match) != 2 || len(idea) != 2 {
		t.Fatal("fixture prompt lacks owned artifact or idea")
	}
	id := strings.TrimSuffix(filepath.Base(match[1]), ".md")
	body := fmt.Sprintf("---\nagent: %s\nidea: %s\nround: 1\ndate: 2026-09-11\n---\n\n## Summary\nLocal process fixture.\n\n## Proposed approach\nExercise the actual grouped runner.\n\n## Existing alternatives\nThe existing runner is the required boundary.\n\n## Concerns / open questions\nNone.\n\n## Risks\nNone.\n", id, idea[1])
	if err := os.WriteFile(match[1], []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

func TestStepBudgetConcurrentIndependentGroupsShareCap(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell fixture")
	}
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	const cap = 2
	binding, err := budget.EnsureStepBinding(context.Background(), root, "idea", cap, 0)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _ = RunMeasured(context.Background(), ExecOptions{Root: root, Agent: telemetryShell("", false),
				Info: LaunchInfo{Idea: "idea", Phase: "review"}, Prompt: "fixture", Timeout: time.Second})
		}()
	}
	close(start)
	wg.Wait()
	started, denied := 0, 0
	for _, r := range terminalRecords(t, root) {
		if r.StartedAt != nil {
			started++
		}
		if r.Outcome.FailureClass != nil && *r.Outcome.FailureClass == "budget_refused" {
			denied++
		}
	}
	if started != cap || denied != 8-cap {
		t.Fatalf("started=%d refused=%d", started, denied)
	}
	state, err := binding.Store.Inspect(context.Background())
	if err != nil || budget.StepCount(state) != cap {
		t.Fatalf("charges=%+v err=%v", state, err)
	}
}

func TestStepBudgetReviewSnapshotKeepsLiveOrigin(t *testing.T) {
	for _, policy := range []string{"step", "launch"} {
		t.Run(policy, func(t *testing.T) {
			root := t.TempDir()
			git := func(args ...string) {
				t.Helper()
				cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("git %v: %v %s", args, err, out)
				}
			}
			git("init", "-q")
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			declareTestLaunchSource(t, root)
			idea, err := protocol.CreateIdea(root, "Snapshot budget fixture", []string{"rev1"})
			if err != nil {
				t.Fatal(err)
			}
			var ledger budget.Store
			if policy == "step" {
				binding, err := budget.EnsureStepBinding(context.Background(), root, idea.Slug, 1, 0)
				if err != nil {
					t.Fatal(err)
				}
				ledger = binding.Store
			} else {
				binding, err := budget.ConfigureLaunchBudget(context.Background(), root, idea.Slug, budget.LaunchPolicy{MaxLaunches: 1})
				if err != nil {
					t.Fatal(err)
				}
				ledger = binding.Store
			}
			mustWrite(t, filepath.Join(idea.Path, "IMPLEMENTATION.md"), "---\nidea: "+idea.Slug+"\nstatus: implemented\n---\n\n## Summary of work\nfixture\n")
			git("add", ".")
			git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "commit", "-qm", "Fixture")
			for attempt := 0; attempt < 2; attempt++ {
				if attempt > 0 {
					if err := os.Remove(filepath.Join(idea.Path, "review", "round-01", "rev1.md")); err != nil {
						t.Fatal(err)
					}
				}
				runID := fmt.Sprintf("review-%d", attempt)
				events := store.New(filepath.Join(root, protocol.DeckDir, "runs", runID))
				results := RunReviewRound(context.Background(), Options{Root: root, RunID: runID, Idea: idea, Round: 1,
					Agents:  []agents.Discovery{{Spec: agents.Spec{ID: "rev1", HeadlessArgs: []string{"-test.run=TestFakeReviewHelper", "--", "parley-fake-review"}, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true}},
					Timeout: 5 * time.Second, Store: events})
				if len(results) != 1 {
					t.Fatalf("review results: %+v", results)
				}
				if attempt == 0 && (!results[0].ArtifactOK || results[0].ExitError != "") {
					t.Fatalf("allowed snapshot review failed: %+v", results)
				}
				if attempt == 1 && results[0].ExitError == "" {
					t.Fatal("snapshot granted a second launch")
				}
				rows, err := events.Load()
				if err != nil {
					t.Fatal(err)
				}
				created := false
				for _, e := range rows {
					if e.Type == "review.snapshot_fallback" {
						t.Fatalf("fixture did not use a snapshot: %+v", e.Data)
					}
					if e.Type == "review.snapshot_created" {
						created = true
						dir, _ := e.Data["dir"].(string)
						if _, err := os.Stat(dir); !os.IsNotExist(err) {
							t.Fatalf("snapshot was not removed: %s %v", dir, err)
						}
					}
				}
				if !created {
					t.Fatal("no actual snapshot created")
				}
			}
			records := terminalRecords(t, root)
			started, denied := 0, 0
			for _, r := range records {
				if r.StartedAt != nil {
					started++
				}
				if r.Outcome.FailureClass != nil && *r.Outcome.FailureClass == "budget_refused" {
					denied++
				}
			}
			if len(records) != 2 || started != 1 || denied != 1 {
				t.Fatalf("snapshot cleanup lost evidence: records=%d started=%d refused=%d", len(records), started, denied)
			}
			state, err := ledger.Inspect(context.Background())
			if err != nil || len(state.Entries) != 1 {
				t.Fatalf("origin charges: %+v %v", state, err)
			}
		})
	}
}
