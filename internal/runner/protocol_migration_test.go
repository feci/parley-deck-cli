package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

func TestProtocolMigrationRunnerActualChildrenRespectRemainingCycles(t *testing.T) {
	for _, tc := range []struct {
		kind          budget.Kind
		phase         string
		cap, children int
	}{{budget.Fixup, "fixup", 5, 1}, {budget.CrossReview, "round-02", 3, 2}} {
		t.Run(string(tc.kind), func(t *testing.T) {
			t.Setenv(config.EnvParleyHome, t.TempDir())
			t.Setenv(config.EnvAgentConfig, "")
			ctx, root := context.Background(), t.TempDir()
			writeLaunchProtocol(t, root)
			ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea")
			mustWrite(t, filepath.Join(ideaDir, "00-prompt.md"), "---\nidea: idea\nparticipants: [alpha, beta]\ntrack: deliberation\nstatus: round-01\n---\nFixture\n")
			for _, id := range []string{"alpha", "beta"} {
				mustWrite(t, filepath.Join(ideaDir, "round-01", id+".md"), "---\nagent: "+id+"\nidea: idea\nround: 1\n---\n\n## Summary\nFixture\n")
			}
			legacy, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{Idea: "idea", Phase: tc.phase, Agent: "legacy", RunID: "legacy", SegmentID: "legacy"})
			if err != nil {
				t.Fatal(err)
			}
			if err := legacy.Finish(telemetry.Outcome{Status: "failed", FailureClass: telemetry.String("budget_refused")}); err != nil {
				t.Fatal(err)
			}
			oldPath := filepath.Join(legacy.Dir, "terminal.json")
			old, err := os.ReadFile(oldPath)
			if err != nil {
				t.Fatal(err)
			}
			before, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: failedCycleAgent("alpha"), Prompt: "fixture", Timeout: time.Second, Info: LaunchInfo{Idea: "idea", Phase: tc.phase, RunID: "before-import"}})
			if err == nil || before.StartedAt != nil || before.Outcome.FailureClass == nil || *before.Outcome.FailureClass != "budget_refused" {
				t.Fatalf("legacy accounting did not stop spawn: %+v %v", before, err)
			}
			i, err := budget.InspectProtocolMigration(ctx, root, "idea", tc.kind, "")
			if err != nil {
				t.Fatal(err)
			}
			r := budget.ProtocolMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "fixture", Reason: "Explicit fixture reconciliation including unobserved protocol work", StartedAt: time.Now().UTC().Add(-time.Hour), TotalActions: tc.cap - 1, Maximum: tc.cap, WritersStopped: true}
			if _, err := budget.MigrateProtocolBudget(ctx, root, "idea", tc.kind, r); err != nil {
				t.Fatal(err)
			}
			for _, run := range []string{"after-import", "new-run"} {
				opts := Options{Root: root, Idea: protocol.IdeaStatus{Slug: "idea", Path: ideaDir, Participants: []string{"alpha", "beta"}}, RunID: run, Round: 2, Agents: []agents.Discovery{failedCycleAgent("alpha"), failedCycleAgent("beta")}, Timeout: 5 * time.Second, Store: store.New(filepath.Join(root, "parley-deck", "runs", run))}
				if tc.kind == budget.Fixup {
					opts.Idea.Participants = []string{"alpha"}
					opts.Agents = opts.Agents[:1]
					if result := RunFixup(ctx, opts); result.ExitError == "" {
						t.Fatal("failed fixture passed")
					}
				} else {
					for _, result := range RunRound(ctx, opts) {
						if result.ExitError == "" {
							t.Fatal("failed child passed")
						}
					}
				}
				status, err := budget.MigrateProtocolBudget(ctx, root, "idea", tc.kind, r)
				if err != nil || status.Spent != tc.cap || !status.StartedAt.Equal(r.StartedAt) {
					t.Fatalf("replay lost charge: %+v %v", status, err)
				}
			}
			started, refused := 0, 0
			for _, rec := range terminalRecords(t, root) {
				if rec.StartedAt != nil {
					started++
				}
				if rec.Outcome.FailureClass != nil && *rec.Outcome.FailureClass == "budget_refused" {
					refused++
				}
			}
			if started != tc.children || refused != 2+tc.children {
				t.Fatalf("group charged per child or reset: started=%d refused=%d", started, refused)
			}
			// A real review child can still execute after the final allowed cycle.
			verify, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: failedCycleAgent("beta"), Prompt: "verification fixture", Timeout: time.Second, Info: LaunchInfo{Idea: "idea", Phase: "review"}})
			if err == nil || verify.StartedAt == nil {
				t.Fatalf("final-cycle verification was refused: %+v %v", verify, err)
			}
			status, err := budget.MigrateProtocolBudget(ctx, root, "idea", tc.kind, r)
			if err != nil || status.Spent != tc.cap {
				t.Fatal("verification charged another cycle", err)
			}
			current, err := os.ReadFile(oldPath)
			if err != nil || string(current) != string(old) {
				t.Fatal("old refusal changed", err)
			}
		})
	}
}

func TestProtocolMigrationRunnerCannotWeakenTrack(t *testing.T) {
	t.Setenv(config.EnvParleyHome, t.TempDir())
	t.Setenv(config.EnvAgentConfig, "")
	ctx, root := context.Background(), t.TempDir()
	writeLaunchProtocol(t, root)
	mustWrite(t, filepath.Join(root, "parley-deck", "ideas", "idea", "00-prompt.md"), "---\nidea: idea\nparticipants: [alpha, beta]\ntrack: fast\n---\nFixture\n")
	i, err := budget.InspectProtocolMigration(ctx, root, "idea", budget.CrossReview, "")
	if err != nil {
		t.Fatal(err)
	}
	_, err = budget.MigrateProtocolBudget(ctx, root, "idea", budget.CrossReview, budget.ProtocolMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "fixture", Reason: "Mismatched ceiling fixture", StartedAt: time.Now().UTC().Add(-time.Hour), Maximum: 3, WritersStopped: true})
	if err != nil {
		t.Fatal(err)
	}
	r, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: failedCycleAgent("alpha"), Prompt: "fixture", Timeout: time.Second, Info: LaunchInfo{Idea: "idea", Phase: "round-02"}})
	if err == nil || r.StartedAt != nil || r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "budget_refused" {
		t.Fatalf("import bypassed fast track: %+v %v", r, err)
	}
}
