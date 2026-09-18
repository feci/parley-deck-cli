package runner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

// Each grouped-launch precheck must refuse on a known-invalid protocol with
// exactly one retained unstarted terminal scoped to the launch it mirrors,
// and pass again once the exact source is restored.
func TestPrecheckLaunchesRefuseKnownInvalidProtocol(t *testing.T) {
	cases := []struct {
		name  string
		phase string
		call  func(ctx context.Context, opts Options) error
	}{
		{"round", "round-02", func(ctx context.Context, opts Options) error {
			opts.Round = 2
			return PrecheckRound(ctx, opts)
		}},
		{"implementation", "implementation", PrecheckImplementation},
		{"review", "review", func(ctx context.Context, opts Options) error {
			opts.Round = 1
			return PrecheckReviewRound(ctx, opts)
		}},
		{"review-consensus", "review-consensus", PrecheckReviewConsensus},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			agent := telemetryShell("exit 0", false)
			agent.Found = true // grouped selection requires an actually discovered agent
			opts := Options{Root: root, RunID: "precheck-" + tc.name,
				Idea:   protocol.IdeaStatus{Slug: "demo", Path: filepath.Join(root, protocol.DeckDir, "ideas", "demo"), Participants: []string{"test-1"}},
				Agents: []agents.Discovery{agent},
				Store:  store.New(t.TempDir())}
			restore := refusePreparedProtocolCache(t, root, LaunchInfo{Idea: "demo", Phase: tc.phase})
			if err := tc.call(context.Background(), opts); err == nil {
				t.Fatal("known invalid protocol passed the precheck")
			}
			records := terminalRecords(t, root)
			if len(records) != 1 || records[0].StartedAt != nil || records[0].Outcome.FailureClass == nil || *records[0].Outcome.FailureClass != "protocol_context_refused" {
				t.Fatalf("precheck lost the exact unstarted refusal terminal: %+v", records)
			}
			if records[0].Metadata.Phase != tc.phase || records[0].Metadata.Idea != "demo" || records[0].Metadata.Agent != "test-1" {
				t.Fatalf("precheck terminal scope mismatch: %+v", records[0].Metadata)
			}
			restore()
			if err := tc.call(context.Background(), opts); err != nil {
				t.Fatalf("restored exact source still refused: %v", err)
			}
		})
	}
}

// With no resolvable agent the real grouped launch renders nothing, so the
// precheck is a no-op rather than a refusal.
func TestPrecheckRoundWithoutSelectedAgentsIsNoop(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	opts := Options{Root: root, RunID: "precheck-unresolved",
		Idea:   protocol.IdeaStatus{Slug: "demo", Path: filepath.Join(root, protocol.DeckDir, "ideas", "demo"), Participants: []string{"ghost"}},
		Agents: []agents.Discovery{telemetryShell("exit 0", false)}, Store: store.New(t.TempDir())}
	if err := PrecheckRound(context.Background(), opts); err != nil {
		t.Fatalf("unresolved participants became a precheck refusal: %v", err)
	}
	if records := terminalRecords(t, root); len(records) != 0 {
		t.Fatalf("noop precheck retained terminals: %+v", records)
	}
}

// N1, manual-runner half: an invalid monetary binding or unreadable defaults
// is a KNOWN preflight failure and must refuse before the cycle/step charges.
// Restoring the exact configuration permits one grouped attempt with the
// original count semantics (one cycle + one step, one started process).
func TestLaunchBudgetPreflightRefusesBeforeCharges(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell fixture")
	}
	for _, scenario := range []string{"missing-monetary-binding", "unreadable-defaults"} {
		t.Run(scenario, func(t *testing.T) {
			t.Setenv(config.EnvParleyHome, t.TempDir())
			t.Setenv(config.EnvAgentConfig, "")
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", "demo")
			if err := os.MkdirAll(ideaDir, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte("---\nidea: demo\nparticipants: [test-1]\ntransport: local-dir\n---\n\n## Problem\nx\n"), 0644); err != nil {
				t.Fatal(err)
			}
			body := "[defaults.loop]\nmax_cost_usd = 1\n"
			if scenario == "unreadable-defaults" {
				body = "[defaults.loop\n"
			}
			if err := os.WriteFile(filepath.Join(root, protocol.DeckDir, "agents.toml"), []byte(body), 0600); err != nil {
				t.Fatal(err)
			}
			step, err := budget.EnsureStepBinding(context.Background(), root, "demo", 5, 0)
			if err != nil {
				t.Fatal(err)
			}
			launch := func() error {
				_, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: telemetryShell("printf 'child\\n'", false),
					Prompt: "task", Timeout: 5 * time.Second, Info: LaunchInfo{RunID: "preflight", Idea: "demo", Phase: "round-02"}})
				return err
			}
			if err := launch(); err == nil {
				t.Fatal("invalid monetary preflight passed")
			}
			stepState, err := step.Store.Inspect(context.Background())
			if err != nil || len(stepState.Entries) != 0 {
				t.Fatalf("invalid preflight spent a step reservation: %+v %v", stepState, err)
			}
			if b, err := budget.LoadCycleBinding(context.Background(), root, "demo", budget.CrossReview); err != nil {
				t.Fatal(err)
			} else if b != nil {
				if state, err := b.Store.Inspect(context.Background()); err != nil || b.Count(state) != 0 {
					t.Fatalf("invalid preflight spent a cross-review cycle: %+v %v", state, err)
				}
			}
			records := terminalRecords(t, root)
			if len(records) != 1 || records[0].StartedAt != nil || records[0].Outcome.FailureClass == nil || *records[0].Outcome.FailureClass != "budget_refused" {
				t.Fatalf("preflight refusal lost the unstarted budget terminal: %+v", records)
			}
			if err := os.Remove(filepath.Join(root, protocol.DeckDir, "agents.toml")); err != nil {
				t.Fatal(err)
			}
			if err := launch(); err != nil {
				t.Fatalf("restored configuration still refused: %v", err)
			}
			stepState, err = step.Store.Inspect(context.Background())
			if err != nil || len(stepState.Entries) != 1 {
				t.Fatalf("grouped attempt did not spend exactly one step: %+v %v", stepState, err)
			}
			cyc, err := budget.LoadCycleBinding(context.Background(), root, "demo", budget.CrossReview)
			if err != nil {
				t.Fatal(err)
			}
			cycState, err := cyc.Store.Inspect(context.Background())
			if err != nil || cyc.Count(cycState) != 1 {
				t.Fatalf("grouped attempt did not spend exactly one cycle: %+v %v", cycState, err)
			}
			records = terminalRecords(t, root)
			started := 0
			for _, rec := range records {
				if rec.StartedAt != nil {
					started++
				}
			}
			if len(records) != 2 || started != 1 {
				t.Fatalf("grouped attempt lifecycle: terminals=%d started=%d", len(records), started)
			}
		})
	}
}
