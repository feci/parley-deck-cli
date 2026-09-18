package runner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/config"
)

func TestMonetaryDefaultRefusesUnreservedManualLaunch(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX child fixture")
	}
	t.Setenv(config.EnvParleyHome, t.TempDir())
	t.Setenv(config.EnvAgentConfig, "")
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "agents.toml"), []byte("[defaults.loop]\nmax_cost_usd = 1.0\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := RunMeasured(context.Background(), ExecOptions{Root: root,
		Agent: telemetryShell("printf 'unreserved child started\\n'", false), Prompt: "task", Timeout: time.Second,
		Info: LaunchInfo{Idea: "monetary-fixture", Phase: "review"}})
	if err == nil || !strings.Contains(err.Error(), "budget") {
		t.Fatalf("configured monetary ceiling must refuse before an unreserved child starts; got %v", err)
	}
	records := terminalRecords(t, root)
	if len(records) != 1 || records[0].StartedAt != nil || records[0].PID != nil || records[0].Outcome.FailureClass == nil || *records[0].Outcome.FailureClass != "budget_refused" {
		t.Fatalf("expected one budget refusal without process execution: %+v", records)
	}
}

func TestMonetaryDefaultUsesSavedExposureAcrossNewRuns(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX child fixtures")
	}
	t.Setenv(config.EnvParleyHome, t.TempDir())
	t.Setenv(config.EnvAgentConfig, "")
	for _, mode := range []string{"known", "unknown", "default-removed", "explicit-zero"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			configPath := filepath.Join(root, "parley-deck", "agents.toml")
			if err := os.WriteFile(configPath, []byte("[defaults.loop]\nmax_cost_usd = 1\n"), 0600); err != nil {
				t.Fatal(err)
			}
			reserve := int64(600000)
			binding, err := budget.ConfigureLaunchBudget(context.Background(), root, "fixture", budget.LaunchPolicy{MaxCostMicros: 1000000, ReserveMicros: &reserve})
			if err != nil {
				t.Fatal(err)
			}
			before, err := binding.Store.Inspect(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			script := `printf '%s\n' '{"type":"result","usage":{"input_tokens":1},"total_cost_usd":0.6}'`
			if mode == "unknown" {
				script = "printf 'opaque child output\\n'"
			}
			opts := ExecOptions{Root: root, Agent: telemetryShell(script, mode != "unknown"), Prompt: "task", Timeout: time.Second,
				Info: LaunchInfo{Idea: "fixture", RunID: "first", Phase: "review"}}
			if _, err := RunMeasured(context.Background(), opts); err != nil {
				t.Fatalf("funded first launch: %v", err)
			}
			if mode == "default-removed" {
				if err := os.Remove(configPath); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "explicit-zero" {
				if err := os.WriteFile(configPath, []byte("[defaults.loop]\nmax_cost_usd = 0\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			opts.Info.RunID = "second"
			if _, err := RunMeasured(context.Background(), opts); err == nil {
				t.Fatal("new run lost the original monetary exposure")
			}
			after, err := binding.Store.Inspect(context.Background())
			if err != nil || len(after.Entries) != 1 || after.StartedAt != before.StartedAt {
				t.Fatalf("charges or clock changed: %+v %v", after, err)
			}
			exposure, err := after.ExposureError()
			if err != nil || exposure != reserve {
				t.Fatalf("exposure=%d err=%v", exposure, err)
			}
			started, refused := 0, 0
			for _, rec := range terminalRecords(t, root) {
				if rec.StartedAt != nil {
					started++
				}
				if rec.Outcome.FailureClass != nil && *rec.Outcome.FailureClass == "budget_refused" {
					refused++
				}
				if mode == "unknown" && rec.Outcome.Usage.CostUSD != nil {
					t.Fatal("conservative reservation was reported as observed cost")
				}
			}
			if started != 1 || refused != 1 {
				t.Fatalf("started=%d refused=%d", started, refused)
			}
		})
	}
}

func TestMonetaryDefaultInvalidConfigurationRefusesWithEvidence(t *testing.T) {
	t.Setenv(config.EnvParleyHome, t.TempDir())
	t.Setenv(config.EnvAgentConfig, "")
	for _, content := range []string{"max_cost_usd = -1", "max_cost_usd = nan", "max_cost_usd = inf", "max_cost_usd = 0.0000001", "max_cost_usd = 1e30", "max_cost_usd = 'not-money'"} {
		t.Run(content, func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			if err := os.WriteFile(filepath.Join(root, "parley-deck", "agents.toml"), []byte("[defaults.loop]\n"+content+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			l, err := beginLaunch(context.Background(), root, "fixture", telemetryShell("", false))
			if err == nil || l != nil {
				t.Fatal("invalid config passed the launch boundary")
			}
			records := terminalRecords(t, root)
			if len(records) != 1 || records[0].StartedAt != nil || records[0].Outcome.FailureClass == nil || *records[0].Outcome.FailureClass != "budget_refused" {
				t.Fatalf("missing refusal evidence: %+v", records)
			}
		})
	}
	root := t.TempDir()
	t.Setenv(config.EnvAgentConfig, "required-but-missing.toml")
	if _, err := beginLaunch(context.Background(), root, "missing-config", telemetryShell("", false)); err == nil {
		t.Fatal("explicit missing config became unlimited")
	}
}

func TestMonetaryDefaultUsesLiveOriginAndDoesNotActivateHandoff(t *testing.T) {
	t.Setenv(config.EnvParleyHome, t.TempDir())
	t.Setenv(config.EnvAgentConfig, "")
	origin, clone := t.TempDir(), t.TempDir()
	for _, root := range []string{origin, clone} {
		writeLaunchProtocol(t, root)
	}
	if err := os.WriteFile(filepath.Join(origin, "parley-deck", "agents.toml"), []byte("[defaults.loop]\nmax_cost_usd = 1\n"), 0600); err != nil {
		t.Fatal(err)
	}
	ctx := withLaunchOrigin(context.Background(), origin)
	if _, err := beginLaunch(ctx, clone, "live-origin", telemetryShell("", false)); err == nil {
		t.Fatal("disposable clone lost the live-origin monetary default")
	}
	if _, err := os.Stat(filepath.Join(clone, ".parley-runtime", "invocations")); !os.IsNotExist(err) {
		t.Fatal("invocation evidence was stored in the disposable clone")
	}
	l, err := beginLaunch(ctx, clone, "handoff", telemetryShell("", false), launchHandoff)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.finish(nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if b, err := budget.LoadLaunchBinding(context.Background(), origin, ""); err != nil || b != nil {
		t.Fatalf("default/handoff invented an operator grant: %+v %v", b, err)
	}
}
