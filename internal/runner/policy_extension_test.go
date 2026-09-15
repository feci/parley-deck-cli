package runner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/config"
)

func TestRuntimePolicyExtensionResumesActualManualProcesses(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX process fixture")
	}
	t.Setenv(config.EnvParleyHome, t.TempDir())
	t.Setenv(config.EnvAgentConfig, "")
	ctx := context.Background()
	for _, kind := range []budget.Kind{budget.Launch, budget.DriverStep} {
		t.Run(string(kind), func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			ceilings := budget.PolicyCeilings{Actions: 2}
			if kind == budget.Launch {
				reserve := int64(1000000)
				binding, err := budget.ConfigureLaunchBudget(ctx, root, "idea", budget.LaunchPolicy{MaxLaunches: 1, MaxCostMicros: reserve, ReserveMicros: &reserve})
				if err != nil {
					t.Fatal(err)
				}
				ceilings.CostMicros = 2 * reserve
				if err := os.WriteFile(filepath.Join(root, "parley-deck", "agents.toml"), []byte("[defaults.loop]\nmax_cost_usd = 1\n"), 0600); err != nil {
					t.Fatal(err)
				}
				// The original trusted context also remains a valid reference.
				ctx = WithLaunchBudget(context.Background(), LaunchBudget{Store: binding.Store, Limits: binding.Policy.Limits(), ReserveMicros: binding.Policy.ReserveMicros})
			} else {
				ctx = context.Background()
				if _, err := budget.EnsureStepBinding(ctx, root, "idea", 1, 0); err != nil {
					t.Fatal(err)
				}
			}
			opts := ExecOptions{Root: root, Agent: telemetryShell("printf 'real opaque process\\n'", false), Prompt: "task", Timeout: time.Second, Info: LaunchInfo{Idea: "idea", RunID: "first", Phase: "review"}}
			if _, err := RunMeasured(ctx, opts); err != nil {
				t.Fatal(err)
			}
			if _, err := RunMeasured(ctx, opts); err == nil {
				t.Fatal("original cap did not refuse second launch")
			}
			s, err := budget.InspectRuntimeBudget(ctx, root, "idea", kind)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := budget.ExtendRuntimeBudget(ctx, root, "idea", kind, budget.PolicyExtensionRequest{DecisionID: "grant", ExpectedPolicySHA256: s.PolicySHA256, Reason: "Explicit finite fixture extension", Ceilings: ceilings}); err != nil {
				t.Fatal(err)
			}
			opts.Info.RunID = "resumed"
			if _, err := RunMeasured(ctx, opts); err != nil {
				t.Fatalf("granted resumed process failed: %v", err)
			}
			if _, err := RunMeasured(ctx, opts); err == nil {
				t.Fatal("extended cap became unlimited")
			}
			started, refused := 0, 0
			for _, r := range terminalRecords(t, root) {
				if r.StartedAt != nil {
					started++
				}
				if r.Outcome.FailureClass != nil && *r.Outcome.FailureClass == "budget_refused" {
					refused++
				}
			}
			if started != 2 || refused != 2 {
				t.Fatalf("started=%d refused=%d", started, refused)
			}
			after, err := budget.InspectRuntimeBudget(ctx, root, "idea", kind)
			if err != nil || after.Spent != 2 || !after.StartedAt.Equal(s.StartedAt) {
				t.Fatalf("lost ledger: %+v %v", after, err)
			}
		})
	}
}
