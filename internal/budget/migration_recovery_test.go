package budget

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMigrationRecoveryProcessChild(t *testing.T) {
	root := os.Getenv("PARLEY_TEST_RECOVERY_PROCESS_ROOT")
	if root == "" {
		return
	}
	raw, err := os.ReadFile(os.Getenv("PARLEY_TEST_RECOVERY_PROCESS_REQUEST"))
	if err != nil {
		t.Fatal(err)
	}
	var r MigrationRecoveryRequest
	if err := json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	fmt.Println("ready")
	if !bufio.NewScanner(os.Stdin).Scan() {
		t.Fatal("parent did not release process barrier")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if _, err := RecoverMigration(ctx, root, "idea", Launch, r); err != nil {
		t.Fatal(err)
	}
}

func TestMigrationRecoverySeparateProcesses(t *testing.T) {
	for _, conflict := range []bool{false, true} {
		t.Run(fmt.Sprintf("conflict-%v", conflict), func(t *testing.T) {
			root := t.TempDir()
			b := inactiveRecoveryFixture(t, root, Launch, "migration-active")
			addRecoveryHistory(t, root, Launch)
			r := recoveryRequestFixture(t, root, Launch)
			requests := t.TempDir()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			t.Cleanup(cancel)
			type child struct {
				cmd    *exec.Cmd
				input  io.WriteCloser
				waited bool
				output chan string
			}
			var children []*child
			for n := 0; n < 3; n++ {
				q := r
				if conflict {
					q.DecisionID = fmt.Sprintf("competing-fixture-%d", n)
				}
				raw, _ := json.Marshal(q)
				path := filepath.Join(requests, fmt.Sprintf("request-%d.json", n))
				if err := os.WriteFile(path, raw, 0600); err != nil {
					t.Fatal(err)
				}
				cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestMigrationRecoveryProcessChild$")
				cmd.Env = append(os.Environ(), "PARLEY_TEST_RECOVERY_PROCESS_ROOT="+root, "PARLEY_TEST_RECOVERY_PROCESS_REQUEST="+path)
				input, err := cmd.StdinPipe()
				if err != nil {
					t.Fatal(err)
				}
				output, err := cmd.StdoutPipe()
				if err != nil {
					t.Fatal(err)
				}
				c := &child{cmd: cmd, input: input, output: make(chan string, 1)}
				if err := cmd.Start(); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if !c.waited {
						_ = c.cmd.Process.Kill()
						_ = c.cmd.Wait()
					}
					_ = c.input.Close()
				})
				ready := make(chan string, 1)
				go func() {
					scanner := bufio.NewScanner(output)
					if scanner.Scan() {
						ready <- scanner.Text()
					} else {
						ready <- ""
					}
					var text strings.Builder
					for scanner.Scan() {
						text.WriteString(scanner.Text())
						text.WriteByte('\n')
					}
					c.output <- text.String()
				}()
				select {
				case value := <-ready:
					if value != "ready" {
						t.Fatal("process did not reach barrier", value)
					}
				case <-ctx.Done():
					t.Fatal("process barrier timed out")
				}
				children = append(children, c)
			}
			for _, c := range children {
				if _, err := fmt.Fprintln(c.input, "go"); err != nil {
					t.Fatal(err)
				}
				_ = c.input.Close()
			}
			passed := 0
			for _, c := range children {
				output := <-c.output
				err := c.cmd.Wait()
				c.waited = true
				if err == nil {
					passed++
				} else {
					var exit *exec.ExitError
					if !conflict || !errors.As(err, &exit) || exit.ExitCode() != 1 || !strings.Contains(output, "active migration cannot be replaced by recovery") {
						t.Fatalf("unexpected child failure: %v: %s", err, output)
					}
				}
			}
			want := 3
			if conflict {
				want = 1
			}
			if passed != want {
				t.Fatalf("successful independent processes: got %d want %d", passed, want)
			}
			j, err := readRecoveryJournal(b)
			if err != nil || len(j.Decisions) != 1 {
				t.Fatalf("processes duplicated or damaged recovery: %v", err)
			}
			status, err := inspectRecoveredBinding(context.Background(), root, "idea", Launch)
			if err != nil || status.Spent != 3 {
				t.Fatalf("process recovery lost/doubled charge: %+v %v", status, err)
			}
		})
	}
}

func inactiveRecoveryFixture(t *testing.T, root string, kind Kind, stopAt string) recoveryBase {
	t.Helper()
	failure := errors.New("fixture interrupted original import")
	persist := func(path string, data []byte) error {
		if filepath.Base(path) == stopAt {
			return failure
		}
		return writeSynced(path, data)
	}
	var err error
	if kind == Launch {
		historicalMigrationCall(t, root, "process-exited", "", true, nil)
		r := migrationFixtureRequest(t, root)
		r.AdditionalLaunches, r.Policy.MaxLaunches = 1, 6
		_, err = migrateLaunchBudget(context.Background(), root, "idea", r, persist)
	} else {
		historicalProtocolCall(t, root, "fixup", true)
		r := protocolMigrationFixture(t, root, kind)
		r.TotalActions, r.Maximum = 1, 6
		_, err = migrateProtocolBudget(context.Background(), root, "idea", kind, r, persist)
	}
	if !errors.Is(err, failure) {
		t.Fatalf("fixture failed before intended boundary: %v", err)
	}
	b, err := loadRecoveryBase(context.Background(), root, "idea", kind)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func recoveryRequestFixture(t *testing.T, root string, kind Kind) MigrationRecoveryRequest {
	t.Helper()
	p, err := InspectMigrationRecovery(context.Background(), root, "idea", kind)
	if err != nil {
		t.Fatal(err)
	}
	if p.Active {
		t.Fatal("fixture unexpectedly active")
	}
	return MigrationRecoveryRequest{ExpectedImportSHA256: p.ImportSHA256, ExpectedHistorySHA256: p.HistorySHA256,
		DecisionID: "fixture-recovery", Reason: "Explicit fixture-only recovery; no real operator decision",
		StartedAt: p.Retained.StartedAt, AccountedActions: p.MinimumAccountedActions, WritersStopped: true}
}

func addRecoveryHistory(t *testing.T, root string, kind Kind) {
	t.Helper()
	if kind == Launch {
		cost := 0.125
		historicalMigrationCall(t, root, "process-exited", "", true, &cost)
	} else {
		historicalProtocolCall(t, root, "fixup", true)
	}
}

func TestMigrationRecoveryAllKindsPreserveHistoryLimitsAndReplay(t *testing.T) {
	for _, kind := range []Kind{Launch, DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			b := inactiveRecoveryFixture(t, root, kind, "migration-active")
			original, _ := os.ReadFile(filepath.Join(b.dir, "migration.json"))
			policy, _ := os.ReadFile(filepath.Join(b.dir, "policy.json"))
			addRecoveryHistory(t, root, kind)
			r := recoveryRequestFixture(t, root, kind)
			r.AccountedActions = 2
			r.StartedAt = r.StartedAt.Add(-time.Hour)
			status, err := RecoverMigration(ctx, root, "idea", kind, r)
			if err != nil {
				t.Fatal(err)
			}
			want := 2
			if kind == Launch {
				want = 4
			}
			if status.Spent != want || !status.StartedAt.Equal(r.StartedAt) {
				t.Fatalf("lost count/epoch: %+v", status)
			}
			if kind == Launch && status.ExposureMicros != nil {
				t.Fatal("unknown imported cost became known")
			}
			ledger := Store{Dir: status.LedgerDir, Scope: b.scope}
			state, err := ledger.Inspect(ctx)
			if err != nil {
				t.Fatal(err)
			}
			for id, old := range b.initial.Entries {
				if !reflect.DeepEqual(state.Entries[id], old) {
					t.Fatal("old observation rewritten")
				}
			}
			for name, old := range map[string][]byte{"migration.json": original, "policy.json": policy} {
				current, err := os.ReadFile(filepath.Join(b.dir, name))
				if err != nil || !bytes.Equal(current, old) {
					t.Fatalf("recovery rewrote %s", name)
				}
			}
			marker, _ := os.ReadFile(filepath.Join(b.dir, "migration-active"))
			if !strings.HasPrefix(string(marker), "parley-budget-migration-recovery-active/v1\n") {
				t.Fatal("old readers would accept recovered state")
			}
			if kind == Launch {
				binding, err := LoadLaunchBinding(ctx, root, "idea")
				if err != nil {
					t.Fatal(err)
				}
				_, err = binding.Reserve(ctx, "actual-later-attempt")
				if err != nil {
					t.Fatal(err)
				}
			} else if err := reserveProtocolFixture(ctx, root, kind, "later-action"); err != nil {
				t.Fatal(err)
			}
			status, err = RecoverMigration(ctx, root, "idea", kind, r)
			if err != nil || status.Spent != want+1 {
				t.Fatalf("replay lost later charge: %+v %v", status, err)
			}
			if kind == Launch || kind == DriverStep {
				_, err = ExtendRuntimeBudget(ctx, root, "idea", kind, PolicyExtensionRequest{DecisionID: "later-fixture-extension", ExpectedPolicySHA256: status.PolicySHA256, Reason: "Separate fixture extension after recovery", Ceilings: PolicyCeilings{Actions: 7}})
			} else {
				_, err = ExtendCycleBudget(ctx, root, "idea", kind, CycleExtensionRequest{DecisionID: "later-fixture-extension", ExpectedPolicySHA256: status.PolicySHA256, Reason: "Separate fixture extension after recovery", Maximum: 7})
			}
			if err != nil {
				t.Fatal("separate later extension failed", err)
			}
			status, err = RecoverMigration(ctx, root, "idea", kind, r)
			if err != nil || status.Spent != want+1 {
				t.Fatalf("replay lost later extension: %+v %v", status, err)
			}
			if kind == Launch {
				q := b.launch.Request
				q.Policy.Version, q.Policy.Scope, q.Policy.Idea = 0, "", ""
				status, err = MigrateLaunchBudget(ctx, root, "idea", q)
			} else {
				status, err = MigrateProtocolBudget(ctx, root, "idea", kind, b.protocol.Request)
			}
			if err != nil || status.Spent != want+1 {
				t.Fatalf("original replay lost recovered state: %+v %v", status, err)
			}
		})
	}
}

func TestMigrationRecoveryInterruptedPublication(t *testing.T) {
	for _, kind := range []Kind{Launch, DriverStep, Fixup, CrossReview} {
		for _, name := range []string{migrationRecoveryFile, "ledger.json", "policy.json", "migration-active"} {
			for _, after := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s-%s-after-%v", kind, name, after), func(t *testing.T) {
					ctx := context.Background()
					root := t.TempDir()
					b := inactiveRecoveryFixture(t, root, kind, "ledger.json")
					addRecoveryHistory(t, root, kind)
					r := recoveryRequestFixture(t, root, kind)
					r.AccountedActions = 2
					failure := errors.New("injected recovery publication error")
					persist := func(path string, data []byte) error {
						if filepath.Base(path) == name && !after {
							return failure
						}
						if err := writeSynced(path, data); err != nil {
							return err
						}
						if filepath.Base(path) == name {
							return failure
						}
						return nil
					}
					if _, err := recoverMigration(ctx, root, "idea", kind, r, persist); !errors.Is(err, failure) {
						t.Fatalf("wrong failure: %v", err)
					}
					if name != "migration-active" || !after {
						if _, err := inspectRecoveredBinding(ctx, root, "idea", kind); err == nil {
							t.Fatal("partial recovery authorized work")
						}
					}
					status, err := RecoverMigration(ctx, root, "idea", kind, r)
					want := 2
					if kind == Launch {
						want = 4
					}
					if err != nil || status.Spent != want {
						t.Fatalf("exact replay: %+v %v", status, err)
					}
					j, err := readRecoveryJournal(b)
					if err != nil || len(j.Decisions) != 1 {
						t.Fatalf("duplicate recovery: %+v %v", j, err)
					}
				})
			}
		}
	}
}

func TestMigrationRecoveryChangedAgainRetainsEveryDecision(t *testing.T) {
	for _, kind := range []Kind{Launch, DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			b := inactiveRecoveryFixture(t, root, kind, "ledger.json")
			addRecoveryHistory(t, root, kind)
			r := recoveryRequestFixture(t, root, kind)
			r.AccountedActions = 2
			persist := func(path string, data []byte) error {
				if err := writeSynced(path, data); err != nil {
					return err
				}
				if filepath.Base(path) == "policy.json" {
					addRecoveryHistory(t, root, kind)
				}
				return nil
			}
			if _, err := recoverMigration(ctx, root, "idea", kind, r, persist); err == nil {
				t.Fatal("changed history activated import")
			}
			if _, err := RecoverMigration(ctx, root, "idea", kind, r); err == nil {
				t.Fatal("stale decision replay accepted")
			}
			if kind == Launch {
				q := b.launch.Request
				q.Policy.Version, q.Policy.Scope, q.Policy.Idea = 0, "", ""
				if _, err := MigrateLaunchBudget(ctx, root, "idea", q); err == nil {
					t.Fatal("original replay bypassed pending recovery")
				}
			} else if _, err := MigrateProtocolBudget(ctx, root, "idea", kind, b.protocol.Request); err == nil {
				t.Fatal("original replay bypassed recovery")
			}
			next := recoveryRequestFixture(t, root, kind)
			next.DecisionID, next.AccountedActions = "second-fixture-recovery", 3
			status, err := RecoverMigration(ctx, root, "idea", kind, next)
			want := 3
			if kind == Launch {
				want = 6
			}
			if err != nil || status.Spent != want {
				t.Fatalf("new decision: %+v %v", status, err)
			}
			j, err := readRecoveryJournal(b)
			if err != nil || len(j.Decisions) != 2 || !reflect.DeepEqual(j.Decisions[0].Request, r) {
				t.Fatal("old recovery decision lost", err)
			}
		})
	}
}

func TestMigrationRecoveryRejectsStaleReducedOrConflictingDecisions(t *testing.T) {
	for _, scenario := range []string{"import-hash", "history-hash", "reduce-count", "later-epoch", "not-stopped", "reuse-id", "conflicting-ledger", "changed-policy", "lost-ledger", "changed-observation"} {
		t.Run(scenario, func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			b := inactiveRecoveryFixture(t, root, Launch, "migration-active")
			r := recoveryRequestFixture(t, root, Launch)
			switch scenario {
			case "import-hash":
				r.ExpectedImportSHA256 = strings.Repeat("0", 64)
			case "history-hash":
				r.ExpectedHistorySHA256 = strings.Repeat("0", 64)
			case "reduce-count":
				r.AccountedActions = 0
			case "later-epoch":
				r.StartedAt = time.Now().UTC()
			case "not-stopped":
				r.WritersStopped = false
			case "reuse-id":
				r.DecisionID = b.decisionID
			case "conflicting-ledger":
				migrationRewrite(t, filepath.Join(b.dir, "ledger", "ledger.json"), func(v map[string]any) { v["entries"] = map[string]any{} })
			case "changed-policy":
				migrationRewrite(t, filepath.Join(b.dir, "policy.json"), func(v map[string]any) { v["max_launches"] = 100 })
			case "lost-ledger":
				if err := os.Remove(filepath.Join(b.dir, "ledger", "ledger.json")); err != nil {
					t.Fatal(err)
				}
			case "changed-observation":
				files, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
				if err != nil || len(files) != 1 {
					t.Fatal("fixture invocation path", err)
				}
				migrationRewrite(t, files[0], func(v map[string]any) { v["outcome"].(map[string]any)["usage"].(map[string]any)["cost_usd"] = 9.0 })
				r = recoveryRequestFixture(t, root, Launch)
			}
			before, err := inspectImportFiles(b.dir)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := RecoverMigration(ctx, root, "idea", Launch, r); err == nil {
				t.Fatal("invalid recovery accepted")
			}
			after, err := inspectImportFiles(b.dir)
			if err != nil || !reflect.DeepEqual(before, after) {
				t.Fatal("refusal changed import state", err)
			}
		})
	}
}

func TestMigrationRecoveryConcurrentReplayAndJournalTampering(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b := inactiveRecoveryFixture(t, root, Launch, "migration-active")
	addRecoveryHistory(t, root, Launch)
	r := recoveryRequestFixture(t, root, Launch)
	var wg sync.WaitGroup
	results := make(chan error, 6)
	for n := 0; n < 6; n++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := RecoverMigration(ctx, root, "idea", Launch, r); results <- err }()
	}
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatal(err)
		}
	}
	j, err := readRecoveryJournal(b)
	if err != nil || len(j.Decisions) != 1 {
		t.Fatal("concurrent replay duplicated decisions", err)
	}
	path := filepath.Join(b.dir, migrationRecoveryFile)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(map[string]any){
		func(v map[string]any) { delete(v, "decisions") },
		func(v map[string]any) { v["base_migration_sha256"] = strings.Repeat("0", 64) },
		func(v map[string]any) {
			v["decisions"].([]any)[0].(map[string]any)["initial"].(map[string]any)["entries"] = map[string]any{}
		},
		func(v map[string]any) {
			v["decisions"].([]any)[0].(map[string]any)["previous_sha256"] = strings.Repeat("0", 64)
		},
		func(v map[string]any) {
			delete(v["decisions"].([]any)[0].(map[string]any)["request"].(map[string]any), "accounted_actions")
		},
	} {
		var fields map[string]any
		if err := json.Unmarshal(raw, &fields); err != nil {
			t.Fatal(err)
		}
		mutate(fields)
		data, _ := json.Marshal(fields)
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadLaunchBinding(ctx, root, "idea"); err == nil {
			t.Fatal("malformed recovered authority accepted")
		}
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadLaunchBinding(ctx, root, "idea"); err == nil {
		t.Fatal("deleted recovery authority accepted")
	}
}
