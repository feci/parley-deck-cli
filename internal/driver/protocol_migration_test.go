package driver

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

type migrationProcessRound struct {
	opts  runner.Options
	calls int
}

func (r *migrationProcessRound) RunRound(ctx context.Context, n int) error {
	r.calls++
	r.opts.Round = n
	results := runner.RunRound(ctx, r.opts)
	for _, result := range results {
		if result.ExitError != "" {
			return errors.New(result.ExitError)
		}
	}
	return errors.New("expected actual child fixture failure")
}

func importDriverFixture(t *testing.T, root string, kind budget.Kind, total, maximum int) budget.ProtocolMigrationRequest {
	t.Helper()
	i, err := budget.InspectProtocolMigration(context.Background(), root, "demo", kind, "")
	if err != nil {
		t.Fatal(err)
	}
	r := budget.ProtocolMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "fixture", Reason: "Explicit historical driver fixture accounting", StartedAt: time.Now().UTC().Add(-time.Hour), TotalActions: total, Maximum: maximum, WritersStopped: true}
	if _, err := budget.MigrateProtocolBudget(context.Background(), root, "demo", kind, r); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestProtocolMigrationDriverNestedActualChildrenAndBlockResume(t *testing.T) {
	for _, stepLimit := range []int{3, 4} {
		t.Run(string(rune('0'+stepLimit)), func(t *testing.T) {
			t.Setenv(config.EnvParleyHome, t.TempDir())
			t.Setenv(config.EnvAgentConfig, "")
			ctx := context.Background()
			parts := []string{"alpha", "beta"}
			ideaDir, runDir := setupIdea(t, parts, "track: deliberation\n")
			root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(root, "parley-deck", "meta"), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "parley-deck", "meta", "version.json"), []byte(`{"protocolRole":"source"}`), 0600); err != nil {
				t.Fatal(err)
			}
			writeAll(t, ideaDir, 1, parts)
			appendEvent(t, runDir, "round.completed", "round-01")
			legacy, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{Idea: "demo", Phase: "round-02", Agent: "legacy", RunID: "legacy", SegmentID: "legacy"})
			if err != nil {
				t.Fatal(err)
			}
			if err := legacy.Finish(telemetry.Outcome{Status: "failed", FailureClass: telemetry.String("budget_refused")}); err != nil {
				t.Fatal(err)
			}
			var discoveries []agents.Discovery
			for _, id := range parts {
				discoveries = append(discoveries, agents.Discovery{Spec: agents.Spec{ID: id, HeadlessArgs: []string{"-test.run=TestCycleDriverChildHelper", "--", "parley-cycle-child"}, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true})
			}
			actual := &migrationProcessRound{opts: runner.Options{Root: root, Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: parts}, RunID: filepath.Base(runDir), Agents: discoveries, Timeout: 5 * time.Second, Store: store.New(runDir)}}
			d := newTestDriver(ideaDir, runDir, parts, 3, true, actual)
			d.cfg.MaxDriverSteps = stepLimit
			d.cfg.HardCrossReviewCap = 3
			d.cfg.Track = "deliberation"
			if _, _, err := d.Advance(ctx); err == nil || actual.calls != 0 {
				t.Fatal("unmigrated driver dispatched", err)
			}
			step := importDriverFixture(t, root, budget.DriverStep, 2, stepLimit)
			cycle := importDriverFixture(t, root, budget.CrossReview, 2, 3)
			if _, _, err := d.Advance(ctx); err == nil || actual.calls != 1 {
				t.Fatal("last allowed actual attempt did not execute", err)
			}
			initial, err := budget.InspectRuntimeBudget(ctx, root, "demo", budget.DriverStep)
			if err != nil || initial.Spent != 3 {
				t.Fatalf("nested children consumed more than one driver step: %+v %v", initial, err)
			}
			// Preserve a BLOCK while changing run identity. Neither path can buy
			// another cross-review or erase the failed nested operation.
			if err := os.WriteFile(filepath.Join(ideaDir, "consensus.md"), []byte("fixture BLOCK"), 0600); err != nil {
				t.Fatal(err)
			}
			resumed := newTestDriver(ideaDir, filepath.Join(filepath.Dir(runDir), "resumed"), parts, 3, true, actual)
			resumed.cfg.HardCrossReviewCap = 3
			resumed.cfg.Track = "deliberation"
			resumed.cfg.Consensus = &fakeConsensus{statusSeq: []string{consensus.TriageBlocked}}
			if _, _, err := resumed.Advance(ctx); err == nil || actual.calls != 1 {
				t.Fatal("BLOCK/new run reset imported accounting", err)
			}
			for _, tc := range []struct {
				kind budget.Kind
				r    budget.ProtocolMigrationRequest
			}{{budget.DriverStep, step}, {budget.CrossReview, cycle}} {
				s, err := budget.MigrateProtocolBudget(ctx, root, "demo", tc.kind, tc.r)
				want := 3
				if tc.kind == budget.DriverStep {
					// A new BLOCK attempt spends its step before the nested cycle
					// refusal. The refusal cannot refund that attempted work.
					want = stepLimit
				}
				if err != nil || s.Spent != want || !s.StartedAt.Equal(tc.r.StartedAt) {
					t.Fatalf("nested operation counted per child or refunded: %+v %v", s, err)
				}
			}
			paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
			if err != nil {
				t.Fatal(err)
			}
			started := 0
			for _, path := range paths {
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				var rec telemetry.Record
				if err := json.Unmarshal(raw, &rec); err != nil {
					t.Fatal(err)
				}
				if rec.StartedAt != nil {
					started++
				}
			}
			if started != 2 {
				t.Fatalf("actual helper starts=%d; want two children in one step/cycle", started)
			}
		})
	}
}

func TestProtocolMigrationDriverFinalAllowedFixupStillVerifies(t *testing.T) {
	t.Setenv(config.EnvParleyHome, t.TempDir())
	t.Setenv(config.EnvAgentConfig, "")
	ctx := context.Background()
	parts := []string{"builder", "reviewer"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImplWithCycles(t, ideaDir, "implemented", 5)
	if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-06"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	r := importDriverFixture(t, root, budget.Fixup, 5, 5)
	fi := &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, ReviewerCount: 2}}
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), Auto: true, AutoImplement: true, MaxFixupCycles: 5, Impl: fi}, &fakeRunner{})
	a, _, err := d.Advance(ctx)
	if err != nil || a != ActionComplete || !contains(fi.calls, "complete") {
		t.Fatalf("import stranded final-cycle verification path: %s %v %v", a, err, fi.calls)
	}
	s, err := budget.MigrateProtocolBudget(ctx, root, "demo", budget.Fixup, r)
	if err != nil || s.Spent != 5 {
		t.Fatal("verification spent another fixup", err)
	}
}
