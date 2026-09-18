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

// Exercise the real pre-start refusal -> migration -> actual child path. No
// provider is called, and the fixture decision authorizes no real migration.
func TestLaunchMigrationResumesActualProcessAfterMonetaryRefusal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX process fixture")
	}
	t.Setenv(config.EnvParleyHome, t.TempDir())
	t.Setenv(config.EnvAgentConfig, "")
	ctx := context.Background()
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "agents.toml"), []byte("[defaults.loop]\nmax_cost_usd = 1\n"), 0600); err != nil {
		t.Fatal(err)
	}
	opts := ExecOptions{Root: root, Agent: telemetryShell("printf 'real opaque process\\n'", false), Prompt: "fixture task", Timeout: time.Second, Info: LaunchInfo{Idea: "idea", RunID: "before-import", Phase: "review"}}
	if _, err := RunMeasured(ctx, opts); err == nil {
		t.Fatal("unreserved monetary default did not refuse")
	}
	old := terminalRecords(t, root)
	if len(old) != 1 || old[0].StartedAt != nil || old[0].Outcome.FailureClass == nil || *old[0].Outcome.FailureClass != "budget_refused" {
		t.Fatalf("want one retained pre-start refusal: %+v", old)
	}
	refusal := filepath.Join(root, ".parley-runtime", "invocations", old[0].InvocationID, "terminal.json")
	original, err := os.ReadFile(refusal)
	if err != nil {
		t.Fatal(err)
	}
	reserve := int64(1_000_000)
	p := budget.LaunchPolicy{MaxLaunches: 1, MaxCostMicros: reserve, ReserveMicros: &reserve, WallClockMS: int64((2 * time.Hour) / time.Millisecond)}
	if _, err := budget.ConfigureLaunchBudget(ctx, root, "idea", p); err == nil {
		t.Fatal("counterexample lost its historical bootstrap barrier")
	}
	i, err := budget.InspectLaunchMigration(ctx, root, "idea")
	if err != nil {
		t.Fatal(err)
	}
	r := budget.LaunchMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "fixture-import", Reason: "Explicit finite fixture accounting", StartedAt: time.Now().UTC().Add(-time.Hour), WritersStopped: true, Policy: p}
	initial, err := budget.MigrateLaunchBudget(ctx, root, "idea", r)
	if err != nil || initial.Spent != 0 {
		t.Fatalf("refusal import: %+v %v", initial, err)
	}
	opts.Info.RunID = "after-import"
	if _, err := RunMeasured(ctx, opts); err != nil {
		t.Fatal("funded actual child failed", err)
	}
	for _, run := range []string{"same-run", "new-run"} {
		opts.Info.RunID = run
		if _, err := RunMeasured(ctx, opts); err == nil {
			t.Fatal("persistent capacity reset", run)
		}
		status, err := budget.MigrateLaunchBudget(ctx, root, "idea", r)
		if err != nil || status.Spent != 1 || !status.StartedAt.Equal(initial.StartedAt) || status.ExposureMicros == nil || *status.ExposureMicros != reserve {
			t.Fatalf("replay changed lifetime exposure: %+v %v", status, err)
		}
	}
	started, refused := 0, 0
	for _, rec := range terminalRecords(t, root) {
		if rec.StartedAt != nil {
			started++
			if rec.Outcome.Usage.CostUSD != nil {
				t.Fatal("opaque child acquired invented price")
			}
		}
		if rec.Outcome.FailureClass != nil && *rec.Outcome.FailureClass == "budget_refused" {
			refused++
		}
	}
	if started != 1 || refused != 3 {
		t.Fatalf("started=%d refused=%d", started, refused)
	}
	current, err := os.ReadFile(refusal)
	if err != nil || string(current) != string(original) {
		t.Fatal("original refusal changed", err)
	}
}
