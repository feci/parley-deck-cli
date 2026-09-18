package budget

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"parley-deck-cli/internal/telemetry"
)

func historicalProtocolCall(t *testing.T, root, phase string, started bool) *telemetry.Invocation {
	t.Helper()
	i, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{Idea: "idea", Phase: phase, RunID: "legacy", SegmentID: "legacy", Agent: "fixture"})
	if err != nil {
		t.Fatal(err)
	}
	if started {
		if err := i.Started(12345); err != nil {
			t.Fatal(err)
		}
	}
	outcome := telemetry.Outcome{Status: "failed", FailureClass: telemetry.String("budget_refused")}
	if started {
		outcome.FailureClass = telemetry.String("process_failure")
	}
	if err := i.Finish(outcome); err != nil {
		t.Fatal(err)
	}
	return i
}

func protocolMigrationFixture(t *testing.T, root string, kind Kind) ProtocolMigrationRequest {
	t.Helper()
	i, err := InspectProtocolMigration(context.Background(), root, "idea", kind, "")
	if err != nil {
		t.Fatal(err)
	}
	return ProtocolMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "fixture-import", Reason: "Explicit fixture-only protocol accounting", StartedAt: time.Now().UTC().Add(-time.Hour), TotalActions: i.LowerBound, Maximum: 3, WritersStopped: true}
}

func reserveProtocolFixture(ctx context.Context, root string, kind Kind, id string) error {
	if kind == DriverStep {
		b, err := LoadStepBinding(ctx, root, "idea")
		if err != nil {
			return err
		}
		_, err = b.reserve(ctx, id)
		return err
	}
	b, err := LoadCycleBinding(ctx, root, "idea", kind)
	if err != nil {
		return err
	}
	_, err = b.Reserve(ctx, id)
	return err
}

func TestProtocolMigrationPreservesTotalsEpochAndLaterGrants(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			historicalProtocolCall(t, root, "fixup", false)
			r := protocolMigrationFixture(t, root, kind)
			r.TotalActions = 2
			s, err := MigrateProtocolBudget(ctx, root, "idea", kind, r)
			if err != nil || s.Spent != 2 || !s.StartedAt.Equal(r.StartedAt) {
				t.Fatalf("import: %+v %v", s, err)
			}
			if err := reserveProtocolFixture(ctx, root, kind, "third"); err != nil {
				t.Fatal(err)
			}
			if err := reserveProtocolFixture(ctx, root, kind, "fourth"); !errors.Is(err, ErrLimit) {
				t.Fatal("cap reset", err)
			}
			s, err = inspectProtocolMigrationBinding(ctx, root, "idea", kind)
			if err != nil {
				t.Fatal(err)
			}
			if kind == DriverStep {
				_, err = ExtendRuntimeBudget(ctx, root, "idea", kind, PolicyExtensionRequest{DecisionID: "grant", Reason: "Explicit finite fixture grant", ExpectedPolicySHA256: s.PolicySHA256, Ceilings: PolicyCeilings{Actions: 4}})
			} else {
				_, err = ExtendCycleBudget(ctx, root, "idea", kind, CycleExtensionRequest{DecisionID: "grant", Reason: "Explicit finite fixture grant", ExpectedPolicySHA256: s.PolicySHA256, Maximum: 4})
			}
			if err != nil {
				t.Fatal(err)
			}
			if err := reserveProtocolFixture(ctx, root, kind, "fourth"); err != nil {
				t.Fatal(err)
			}
			s, err = MigrateProtocolBudget(ctx, root, "idea", kind, r)
			if err != nil || s.Spent != 4 || !s.StartedAt.Equal(r.StartedAt) {
				t.Fatalf("replay lost grant/spend: %+v %v", s, err)
			}
			if err := reserveProtocolFixture(ctx, root, kind, "fifth"); !errors.Is(err, ErrLimit) {
				t.Fatal("grant became unlimited", err)
			}
			changed := r
			changed.TotalActions++
			if _, err := MigrateProtocolBudget(ctx, root, "idea", kind, changed); err == nil {
				t.Fatal("conflicting replay accepted")
			}
			if kind == DriverStep {
				_, err = EnsureStepBinding(ctx, root, "idea", 3, 0)
			} else {
				_, err = EnsureCycleBinding(ctx, root, "idea", kind, 3, 2, "", "")
			}
			if err != nil {
				t.Fatal("original runtime cannot resume", err)
			}
		})
	}
}

func TestProtocolMigrationRecoversPublicationWithoutActivatingPartialState(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		for _, file := range []string{"migration.json", "ledger.json", "policy.json", "migration-active"} {
			for _, after := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%s/after-%v", kind, file, after), func(t *testing.T) {
					ctx := context.Background()
					root := t.TempDir()
					r := protocolMigrationFixture(t, root, kind)
					r.TotalActions = 1
					fired := false
					persist := func(path string, raw []byte) error {
						if filepath.Base(path) == file && !fired {
							fired = true
							if after {
								if err := writeSynced(path, raw); err != nil {
									return err
								}
							}
							return errors.New("injected publication fault")
						}
						return writeSynced(path, raw)
					}
					if _, err := migrateProtocolBudget(ctx, root, "idea", kind, r, persist); err == nil || !fired {
						t.Fatal("failure not exercised", err)
					}
					if _, err := inspectProtocolMigrationBinding(ctx, root, "idea", kind); err == nil && !(file == "migration-active" && after) {
						t.Fatal("partial import became usable")
					}
					s, err := MigrateProtocolBudget(ctx, root, "idea", kind, r)
					if err != nil || s.Spent != 1 {
						t.Fatalf("replay failed: %+v %v", s, err)
					}
					if _, err := os.Stat(filepath.Join(s.LedgerDir, "ledger-established")); err != nil {
						t.Fatal(err)
					}
					if err := reserveProtocolFixture(ctx, root, kind, "later"); err != nil {
						t.Fatal(err)
					}
				})
			}
		}
	}
}

func TestProtocolMigrationRejectsStaleAndChangingSources(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	r := protocolMigrationFixture(t, root, Fixup)
	historicalProtocolCall(t, root, "fixup", false)
	if _, err := MigrateProtocolBudget(ctx, root, "idea", Fixup, r); err == nil {
		t.Fatal("stale inventory accepted")
	}
	dir, _, err := protocolMigrationScope(ctx, root, "idea", Fixup)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("stale import created a scope", err)
	}
	r = protocolMigrationFixture(t, root, Fixup)
	persist := func(path string, raw []byte) error {
		if err := writeSynced(path, raw); err != nil {
			return err
		}
		if filepath.Base(path) == "policy.json" {
			historicalProtocolCall(t, root, "fixup", true)
		}
		return nil
	}
	if _, err := migrateProtocolBudget(ctx, root, "idea", Fixup, r, persist); err == nil {
		t.Fatal("source change silently accepted")
	}
	if _, err := LoadCycleBinding(ctx, root, "idea", Fixup); err == nil {
		t.Fatal("changed-source import activated")
	}
	if _, err := MigrateProtocolBudget(ctx, root, "idea", Fixup, r); err == nil {
		t.Fatal("stale partial replay rewrote import")
	}
}

func TestProtocolMigrationDeletedAndNormalizedAuthorityRefuses(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		for _, mutation := range []string{"policy", "marker", "reference", "zero-total", "epoch", "entry"} {
			t.Run(string(kind)+"/"+mutation, func(t *testing.T) {
				ctx := context.Background()
				root := t.TempDir()
				r := protocolMigrationFixture(t, root, kind)
				if mutation == "entry" {
					r.TotalActions = 1
				}
				s, err := MigrateProtocolBudget(ctx, root, "idea", kind, r)
				if err != nil {
					t.Fatal(err)
				}
				dir := filepath.Dir(s.LedgerDir)
				switch mutation {
				case "policy":
					if err := os.Remove(filepath.Join(dir, "policy.json")); err != nil {
						t.Fatal(err)
					}
				case "marker":
					if err := os.Remove(filepath.Join(dir, "migration-active")); err != nil {
						t.Fatal(err)
					}
				case "reference":
					migrationRewrite(t, filepath.Join(dir, "policy.json"), func(v map[string]any) { v["migration_sha256"] = nil })
				case "zero-total":
					migrationRewrite(t, filepath.Join(dir, "migration.json"), func(v map[string]any) { delete(v["request"].(map[string]any), "total_actions") })
				case "epoch":
					migrationRewrite(t, filepath.Join(s.LedgerDir, "ledger.json"), func(v map[string]any) { v["started_at"] = time.Now().UTC().Format(time.RFC3339Nano) })
				case "entry":
					migrationRewrite(t, filepath.Join(s.LedgerDir, "ledger.json"), func(v map[string]any) { v["entries"] = map[string]any{} })
				}
				if _, err := inspectProtocolMigrationBinding(ctx, root, "idea", kind); err == nil {
					t.Fatal("corrupt import authorized accounting")
				}
				if mutation == "policy" {
					if _, err := MigrateProtocolBudget(ctx, root, "idea", kind, r); err == nil {
						t.Fatal("active missing policy recreated")
					}
				}
			})
		}
	}
}

func TestProtocolMigrationCachedSessionsCheckImportedIdentities(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			r := protocolMigrationFixture(t, root, kind)
			r.TotalActions = 1
			s, err := MigrateProtocolBudget(ctx, root, "idea", kind, r)
			if err != nil {
				t.Fatal(err)
			}
			var scoped context.Context
			var finish func()
			var charge func() error
			if kind == DriverStep {
				b, e := LoadStepBinding(ctx, root, "idea")
				if e != nil {
					t.Fatal(e)
				}
				scoped, finish, err = OpenStepSession(ctx, b)
				charge = func() error { return ChargeStep(scoped) }
			} else {
				b, e := LoadCycleBinding(ctx, root, "idea", kind)
				if e != nil {
					t.Fatal(e)
				}
				scoped, finish, err = OpenCycleSession(ctx, b)
				charge = func() error { _, e := ChargeCycle(scoped, kind); return e }
			}
			if err != nil {
				t.Fatal(err)
			}
			defer finish()
			if err := charge(); err != nil {
				t.Fatal(err)
			}
			migrationRewrite(t, filepath.Join(s.LedgerDir, "ledger.json"), func(v map[string]any) {
				entries := v["entries"].(map[string]any)
				id := key("protocol-migration:" + string(kind) + ":" + r.DecisionID + ":0")
				entries[key("substituted-charge")] = entries[id]
				delete(entries, id)
			})
			if err := charge(); err == nil {
				t.Fatal("aggregate count hid a deleted imported identity")
			}
		})
	}
}

func TestProtocolMigrationConcurrentDecisionsHaveOneWinner(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	r := protocolMigrationFixture(t, root, DriverStep)
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, id := range []string{"first", "second"} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			q := r
			q.DecisionID = id
			<-start
			_, err := MigrateProtocolBudget(ctx, root, "idea", DriverStep, q)
			results <- err
		}(id)
	}
	close(start)
	wg.Wait()
	close(results)
	winners := 0
	for err := range results {
		if err == nil {
			winners++
		}
	}
	if winners != 1 {
		t.Fatalf("migration decisions winning=%d", winners)
	}
}

func TestProtocolMigrationExhaustedImportDoesNotWaiveCaps(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			r := protocolMigrationFixture(t, root, kind)
			r.TotalActions = 4
			s, err := MigrateProtocolBudget(ctx, root, "idea", kind, r)
			if err != nil || s.Spent != 4 {
				t.Fatalf("cannot retain exhausted state: %+v %v", s, err)
			}
			if err := reserveProtocolFixture(ctx, root, kind, "over-cap"); !errors.Is(err, ErrLimit) {
				t.Fatal("exhausted import authorized work", err)
			}
			var policy map[string]any
			if err := json.Unmarshal(s.Policy, &policy); err != nil {
				t.Fatal(err)
			}
			if policy["migration_sha256"] == nil {
				t.Fatal("import provenance missing")
			}
		})
	}
}

func TestProtocolMigrationZeroCeilingsAndElapsedWallTime(t *testing.T) {
	ctx := context.Background()
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			root := t.TempDir()
			r := protocolMigrationFixture(t, root, kind)
			r.Maximum = 0
			if _, err := MigrateProtocolBudget(ctx, root, "idea", kind, r); err != nil {
				t.Fatal(err)
			}
			err := reserveProtocolFixture(ctx, root, kind, "after-zero-ceiling")
			if kind == DriverStep && err != nil || kind != DriverStep && !errors.Is(err, ErrLimit) {
				t.Fatal("zero ceiling semantics changed", err)
			}
		})
	}
	root := t.TempDir()
	r := protocolMigrationFixture(t, root, DriverStep)
	r.WallClockNS = int64(30 * time.Minute)
	if _, err := MigrateProtocolBudget(ctx, root, "idea", DriverStep, r); err != nil {
		t.Fatal(err)
	}
	if err := reserveProtocolFixture(ctx, root, DriverStep, "after-expired-import"); !errors.Is(err, ErrLimit) {
		t.Fatal("import restarted elapsed wall clock", err)
	}
}
