package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

func TestLaunchBudgetConcurrentRootsAndReservationFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX process fixtures")
	}
	for _, corrupt := range []bool{false, true} {
		name := "shared-cap"
		if corrupt {
			name = "corrupt-ledger"
		}
		t.Run(name, func(t *testing.T) {
			roots := []string{t.TempDir(), t.TempDir()}
			for _, root := range roots {
				writeLaunchProtocol(t, root)
			}
			ledger := budget.Store{Dir: filepath.Join(t.TempDir(), "ledger"), Scope: "same-idea-across-roots"}
			if corrupt {
				if err := os.MkdirAll(ledger.Dir, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(ledger.Dir, "ledger.json"), []byte("corrupt prior charges"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			ctx := WithLaunchBudget(context.Background(), LaunchBudget{Store: ledger, Limits: budget.Limits{Actions: map[budget.Kind]int{budget.Launch: 3}}})
			var wg sync.WaitGroup
			start := make(chan struct{})
			for i := 0; i < 8; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-start
					_, _ = RunMeasured(ctx, ExecOptions{Root: roots[i%2], Agent: telemetryShell("printf 'real process\\n'", false), Prompt: "task", Timeout: 5 * time.Second})
				}(i)
			}
			close(start)
			wg.Wait()
			started, denied, attempts := 0, 0, 0
			for _, root := range roots {
				for _, rec := range terminalRecords(t, root) {
					attempts++
					if rec.StartedAt != nil {
						started++
					}
					if rec.Outcome.FailureClass != nil && *rec.Outcome.FailureClass == "budget_refused" {
						denied++
					}
				}
			}
			wantStarted := 3
			if corrupt {
				wantStarted = 0
			}
			if attempts != 8 || started != wantStarted || denied != 8-wantStarted {
				t.Fatalf("attempts=%d started=%d denied=%d", attempts, started, denied)
			}
			if !corrupt {
				state, err := ledger.Inspect(context.Background())
				if err != nil || len(state.Entries) != 3 {
					t.Fatalf("shared charges: %+v %v", state, err)
				}
			} else {
				raw, err := os.ReadFile(filepath.Join(ledger.Dir, "ledger.json"))
				if err != nil || string(raw) != "corrupt prior charges" {
					t.Fatal("replaced corrupt history with a new grant")
				}
			}
		})
	}
}

func TestLaunchBudgetAcrossProcessBoundaries(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX fixture commands; Windows cross-build is separate")
	}
	for _, policy := range []string{"context", "launch", "step"} {
		configured := policy == "launch"
		steps := policy == "step"
		for _, surface := range []string{"manual", "round-process", "consult", "probe", "interactive", "acp"} {
			name := surface + "-" + policy
			t.Run(name, func(t *testing.T) {
				root := t.TempDir()
				writeLaunchProtocol(t, root)
				ledger := budget.Store{Dir: filepath.Join(t.TempDir(), "shared-ledger"), Scope: "idea:fixture"}
				limits := budget.Limits{Actions: map[budget.Kind]int{budget.Launch: 1}}
				ctx := WithLaunchBudget(context.Background(), LaunchBudget{Store: ledger, Limits: limits})
				// The injected copy must survive both caller mutation and metadata replacement.
				limits.Actions[budget.Launch] = 100
				if configured && surface != "acp" {
					binding, err := budget.ConfigureLaunchBudget(context.Background(), root, "", budget.LaunchPolicy{MaxLaunches: 1})
					if err != nil {
						t.Fatal(err)
					}
					ledger = binding.Store
					// Ordinary application calls attach no explicit budget context.
					ctx = context.Background()
				}
				ideaID := ""
				if steps && surface != "acp" {
					ideaID = "step-fixture"
					binding, err := budget.EnsureStepBinding(context.Background(), root, ideaID, 1, 0)
					if err != nil {
						t.Fatal(err)
					}
					ledger, ctx = binding.Store, context.Background()
				}
				ctx = WithLaunchInfo(ctx, LaunchInfo{RunID: "changed-run", Idea: ideaID, Phase: "review"})
				agent := telemetryShell("printf 'actual child\\n'", false)
				var launch func(int) error
				switch surface {
				case "manual":
					launch = func(int) error {
						_, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: agent, Prompt: "task", Timeout: time.Second, Info: LaunchInfo{Idea: ideaID, Phase: "review"}})
						return err
					}
				case "round-process":
					launch = func(int) error { return runTelemetryFixture(ctx, root, agent) }
				case "consult":
					launch = func(int) error {
						r := RunConsult(ctx, ConsultOptions{Root: root, Agent: agent, Prompt: "task", Timeout: time.Second,
							StdoutPath: filepath.Join(root, "consult-out"), StderrPath: filepath.Join(root, "consult-err")})
						if r.ExitError != "" {
							return errors.New(r.ExitError)
						}
						return nil
					}
				case "probe":
					launch = func(int) error {
						cmd, cleanup, err := ProbeCommandFor(WithLaunchInfo(ctx, LaunchInfo{Idea: ideaID, Phase: "preflight"}), root, agent, "PONG")
						if cleanup != nil {
							defer cleanup()
						}
						if err != nil {
							return err
						}
						return cmd.Run()
					}
				case "interactive":
					interactive, terminal := interactiveFixture(t, root, "printf 'actual terminal child\\n'")
					launch = func(int) error {
						return RunInteractive(ctx, root, interactive, "task", "", terminal, terminal, terminal)
					}
				case "acp":
					if err := protocol.InitWorkspace(root); err != nil {
						t.Fatal(err)
					}
					declareTestLaunchSource(t, root)
					idea, err := protocol.CreateIdea(root, "Budget ACP fixture", []string{"fake-acp"})
					if err != nil {
						t.Fatal(err)
					}
					if configured {
						binding, err := budget.ConfigureLaunchBudget(context.Background(), root, idea.Slug, budget.LaunchPolicy{MaxLaunches: 1})
						if err != nil {
							t.Fatal(err)
						}
						ledger = binding.Store
						ctx = context.Background()
					}
					if steps {
						binding, err := budget.EnsureStepBinding(context.Background(), root, idea.Slug, 1, 0)
						if err != nil {
							t.Fatal(err)
						}
						ledger, ctx = binding.Store, context.Background()
					}
					acpAgent := agents.Discovery{Spec: agents.Spec{ID: "fake-acp", LaunchMode: agents.LaunchACP,
						ACPArgs: []string{"-test.run=TestFakeACPAgentHelper", "--", "parley-fake-acp-agent"}, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true}
					launch = func(attempt int) error {
						if attempt > 0 {
							if err := os.Remove(filepath.Join(idea.Path, "round-01/fake-acp.md")); err != nil {
								return err
							}
						}
						runID := []string{"first", "second"}[attempt]
						results := RunRoundOne(ctx, Options{Root: root, RunID: runID, Idea: idea, Task: "Budget ACP fixture",
							Agents: []agents.Discovery{acpAgent}, Timeout: 10 * time.Second, Store: store.New(filepath.Join(root, protocol.DeckDir, "runs", runID))})
						if len(results) != 1 || results[0].ExitError != "" || !results[0].ArtifactOK {
							return errors.New("ACP did not complete with its own artifact")
						}
						return nil
					}
				}
				if err := launch(0); err != nil {
					t.Fatalf("first actual %s launch: %v", surface, err)
				}
				if err := launch(1); err == nil {
					t.Fatal("second launch exceeded one-call cap")
				}
				records := terminalRecords(t, root)
				if len(records) != 2 {
					t.Fatalf("two attempted requests required, got %d", len(records))
				}
				started, refused := 0, 0
				for _, rec := range records {
					if rec.StartedAt != nil && rec.PID != nil {
						started++
					}
					if rec.Outcome.FailureClass != nil && *rec.Outcome.FailureClass == "budget_refused" {
						refused++
						if rec.PID != nil || rec.StartedAt != nil {
							t.Fatal("denied request spawned a child")
						}
					}
				}
				if started != 1 || refused != 1 {
					t.Fatalf("actual lifecycle started=%d refused=%d", started, refused)
				}
				snapshot, err := ledger.Inspect(context.Background())
				if err != nil || len(snapshot.Entries) != 1 {
					t.Fatalf("spent entries=%+v err=%v", snapshot.Entries, err)
				}
				for _, entry := range snapshot.Entries {
					if !steps && !entry.Settled {
						t.Fatal("actual process did not settle")
					}
				}
			})
		}
	}
}

func TestLaunchBudgetHandoffDoesNotChargeProcess(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	ledger := budget.Store{Dir: filepath.Join(t.TempDir(), "ledger"), Scope: "handoff"}
	ctx := WithLaunchBudget(context.Background(), LaunchBudget{Store: ledger, Limits: budget.Limits{Denied: map[budget.Kind]bool{budget.Launch: true}}})
	packet, err := WriteHandoffPacket(HandoffOptions{Context: ctx, Root: root, RunID: "handoff", Agent: telemetryShell("", false), Prompt: "task"})
	if err != nil || packet.InvocationID == "" {
		t.Fatalf("handoff failed: %v", err)
	}
	if _, err := os.Stat(ledger.Dir); !os.IsNotExist(err) {
		t.Fatal("nonexecuting handoff touched the process ledger")
	}
}

func TestLaunchBudgetRetainsFailedAndUnreportedSpend(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX process fixtures")
	}
	for _, scenario := range []string{"failed-start", "timeout-known-cost", "terminal-write", "unknown-cost", "settlement-write"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			ledger := budget.Store{Dir: filepath.Join(t.TempDir(), "ledger"), Scope: "failed-spend"}
			reserve := int64(200000)
			ctx := WithLaunchBudget(context.Background(), LaunchBudget{Store: ledger, ReserveMicros: &reserve, Limits: budget.Limits{CostMicros: reserve}})
			script := `printf '%s\n' '{"type":"result","usage":{"input_tokens":1},"total_cost_usd":0.125}'`
			if scenario == "unknown-cost" {
				script = "printf 'opaque output\\n'"
			}
			if scenario == "timeout-known-cost" {
				script += "; exec sleep 5"
			}
			agent := telemetryShell(script, scenario != "unknown-cost")
			if scenario == "failed-start" {
				agent.Path = filepath.Join(root, "nonexistent")
			}
			info := LaunchInfo{Observe: func(rec telemetry.Record) {
				if scenario == "terminal-write" && rec.Type == "invocation.requested" {
					if err := os.Mkdir(filepath.Join(root, ".parley-runtime/invocations", rec.InvocationID, "terminal.json"), 0700); err != nil {
						t.Fatal(err)
					}
				}
				if scenario == "settlement-write" && rec.Type == "invocation.terminal" {
					if err := os.Rename(filepath.Join(ledger.Dir, "ledger.json"), filepath.Join(ledger.Dir, "saved-ledger.json")); err != nil {
						t.Fatal(err)
					}
					if err := os.Mkdir(filepath.Join(ledger.Dir, "ledger.json"), 0700); err != nil {
						t.Fatal(err)
					}
				}
			}}
			_, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: agent, Prompt: "task", Timeout: 300 * time.Millisecond, Info: info})
			if scenario == "unknown-cost" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil {
				t.Fatal("failed execution/accounting returned success")
			}
			if scenario == "settlement-write" {
				if !strings.Contains(err.Error(), "reservation retained") {
					t.Fatalf("lost settlement failure: %v", err)
				}
				if err := os.Remove(filepath.Join(ledger.Dir, "ledger.json")); err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(filepath.Join(ledger.Dir, "saved-ledger.json"), filepath.Join(ledger.Dir, "ledger.json")); err != nil {
					t.Fatal(err)
				}
			}
			snapshot, inspectErr := ledger.Inspect(context.Background())
			if inspectErr != nil || len(snapshot.Entries) != 1 {
				t.Fatalf("lost charge: %v %+v", inspectErr, snapshot)
			}
			for _, entry := range snapshot.Entries {
				if entry.Settled != (scenario != "settlement-write") {
					t.Fatalf("settlement state: %+v", entry)
				}
				if scenario == "timeout-known-cost" || scenario == "terminal-write" {
					if entry.ActualMicros == nil || *entry.ActualMicros != 125000 {
						t.Fatalf("known observed spend not settled: %+v", entry)
					}
				}
			}
			if scenario == "failed-start" || scenario == "unknown-cost" || scenario == "settlement-write" {
				exposure, err := snapshot.ExposureError()
				if err != nil || exposure != reserve {
					t.Fatalf("reservation refunded: %d %v", exposure, err)
				}
			}
		})
	}
}
