package runner

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestRunnerActionIdentityActualFixupAndManualAttempt(t *testing.T) {
	root := t.TempDir()
	if e := protocol.InitWorkspace(root); e != nil {
		t.Fatal(e)
	}
	declareTestLaunchSource(t, root)
	idea, e := protocol.CreateIdea(root, "Action identity fixture", []string{"builder", "reviewer"})
	if e != nil {
		t.Fatal(e)
	}
	mustWrite(t, filepath.Join(idea.Path, "00-prompt.md"), "---\nidea: "+idea.Slug+"\nparticipants: [builder, reviewer]\ntrack: deliberation\nstatus: implemented\n---\nFixture\n")
	idea.Participants = []string{"builder"}
	step, e := budget.EnsureStepBinding(context.Background(), root, idea.Slug, 10, 0)
	if e != nil {
		t.Fatal(e)
	}
	runID := "action-fixup"
	launchStore := budget.Store{Dir: t.TempDir(), Scope: "actual-fixup-launch"}
	zero := int64(0)
	ctx := WithLaunchBudget(context.Background(), LaunchBudget{Store: launchStore, Limits: budget.Limits{Actions: map[budget.Kind]int{budget.Launch: 2}}, ReserveMicros: &zero})
	result := RunFixup(ctx, Options{Root: root, Idea: idea, RunID: runID, Task: "private task text", Agents: []agents.Discovery{failedCycleAgent("builder")}, Timeout: 5 * time.Second, Store: store.New(filepath.Join(root, "parley-deck", "runs", runID))})
	if result.ExitError == "" {
		t.Fatal("fixture child unexpectedly succeeded")
	}
	cycle, e := budget.LoadCycleBinding(context.Background(), root, idea.Slug, budget.Fixup)
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range []budget.Store{step.Store, cycle.Store, launchStore} {
		rows, e := s.InspectActions(context.Background())
		if e != nil || len(rows) != 1 {
			t.Fatalf("group charge: %v", e)
		}
		id := rows[0].Identity
		if id == nil || id.Input.Basis != "runtime-input" || id.Input.Operation != "fixup" || id.Input.RunID != runID {
			t.Fatalf("actual runner lost original recipe: %+v", id)
		}
		replay, e := s.ReplayAction(context.Background(), rows[0].EntryKey, rows[0].IdentitySHA256)
		if e != nil || replay.Permission != "none" {
			t.Fatal("failed child receipt implied a rerun")
		}
		if s.Scope == launchStore.Scope {
			if result.InvocationID == "" || rows[0].EntryKey != fmt.Sprintf("%x", sha256.Sum256([]byte(result.InvocationID))) {
				t.Fatal("launch charge lost its original invocation identity")
			}
			state, e := s.Inspect(context.Background())
			if e != nil || !state.Entries[rows[0].EntryKey].Settled {
				t.Fatal("failed launch did not retain its settled charge")
			}
		}
	}
	r, e := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: failedCycleAgent("reviewer"), Prompt: "manual verification", Timeout: 5 * time.Second, Info: LaunchInfo{RunID: "manual-action", Idea: idea.Slug, Phase: "review"}})
	if e == nil || r.StartedAt == nil {
		t.Fatalf("manual fixture failed to launch: %v", e)
	}
	rows, e := step.Store.InspectActions(context.Background())
	if e != nil || len(rows) != 2 {
		t.Fatalf("manual step identity: %v", e)
	}
	found := false
	for _, row := range rows {
		if row.Identity.Input.RunID == "manual-action" {
			found = true
			if row.Identity.Input.Basis != "launch-metadata" {
				t.Fatal("manual metadata claimed complete parent input")
			}
		}
	}
	if !found {
		t.Fatal("manual request identity was not retained")
	}
}
func TestRunnerActionInputChangesWithProtocolAndTask(t *testing.T) {
	root := t.TempDir()
	idea := protocol.IdeaStatus{Slug: "idea", Path: filepath.Join(root, "idea"), Participants: []string{"a"}}
	opts := Options{Root: root, Idea: idea, RunID: "same-run", Task: "one"}
	s := budget.Store{Dir: t.TempDir(), Scope: "recipe-scope"}
	request := budget.Request{ID: "original-attempt", Kind: budget.DriverStep}
	if _, e := s.Reserve(withRunnerActionInput(context.Background(), opts, "implementation"), request, budget.Limits{}); e != nil {
		t.Fatal(e)
	}
	opts.Task = "different"
	if _, e := s.Reserve(withRunnerActionInput(context.Background(), opts, "implementation"), request, budget.Limits{}); !errors.Is(e, budget.ErrActionConflict) {
		t.Fatalf("changed task reused identity: %v", e)
	}
	opts.Task = "one"
	mustWrite(t, filepath.Join(idea.Path, "FINAL.md"), "A changed binding input")
	if _, e := s.Reserve(withRunnerActionInput(context.Background(), opts, "implementation"), request, budget.Limits{}); !errors.Is(e, budget.ErrActionConflict) {
		t.Fatalf("changed FINAL reused identity: %v", e)
	}
}
