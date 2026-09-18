package budget

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"parley-deck-cli/internal/telemetry"
)

func historicalMigrationCall(t *testing.T, root, status, failure string, started bool, cost *float64) *telemetry.Invocation {
	t.Helper()
	i, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{Idea: "idea", Phase: "implementation", Agent: "fixture", RunID: "run", SegmentID: "segment"})
	if err != nil {
		t.Fatal(err)
	}
	if started {
		if err := i.Started(12345); err != nil {
			t.Fatal(err)
		}
	}
	if err := i.Finish(telemetry.Outcome{Status: status, FailureClass: telemetry.String(failure), Usage: telemetry.Usage{CostUSD: cost, CostBasis: "cli-estimate", Source: "fixture", Coverage: "partial"}}); err != nil {
		t.Fatal(err)
	}
	return i
}

func migrationFixtureRequest(t *testing.T, root string) LaunchMigrationRequest {
	t.Helper()
	i, err := InspectLaunchMigration(context.Background(), root, "idea")
	if err != nil {
		t.Fatal(err)
	}
	return LaunchMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "fixture-import", Reason: "Explicit fixture accounting, no real operator decision", StartedAt: time.Now().UTC().Add(-time.Hour), WritersStopped: true, Policy: LaunchPolicy{MaxLaunches: 10}}
}

func TestLaunchMigrationRetainsRefusalWithoutInventingLaunch(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	old := historicalMigrationCall(t, root, "failed", "budget_refused", false, nil)
	raw, err := os.ReadFile(filepath.Join(old.Dir, "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ConfigureLaunchBudget(ctx, root, "idea", LaunchPolicy{MaxLaunches: 1}); err == nil {
		t.Fatal("counterexample must require migration")
	}
	r := migrationFixtureRequest(t, root)
	r.Policy.MaxLaunches = 1
	status, err := MigrateLaunchBudget(ctx, root, "idea", r)
	if err != nil {
		t.Fatal(err)
	}
	if status.Spent != 0 || status.ExposureMicros == nil || *status.ExposureMicros != 0 || !status.StartedAt.Equal(r.StartedAt) {
		t.Fatalf("refusal invented work or reset epoch: %+v", status)
	}
	current, _ := os.ReadFile(filepath.Join(old.Dir, "terminal.json"))
	if string(current) != string(raw) {
		t.Fatal("migration changed original refusal")
	}
	b, err := LoadLaunchBinding(ctx, root, "idea")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Reserve(ctx, "new-actual-attempt"); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Reserve(ctx, "after-cap"); !errors.Is(err, ErrLimit) {
		t.Fatalf("imported policy was bypassed: %v", err)
	}
	status, err = MigrateLaunchBudget(ctx, root, "idea", r)
	if err != nil || status.Spent != 1 {
		t.Fatalf("replay reset later spend: %+v %v", status, err)
	}
	if _, err := ConfigureLaunchBudget(ctx, root, "idea", r.Policy); err != nil {
		t.Fatalf("exact original configuration replay failed: %v", err)
	}
	dir := filepath.Dir(b.Store.Dir)
	if _, err := os.Stat(filepath.Join(dir, "ledger", "ledger-established")); err != nil {
		t.Fatal("missing continuity witness", err)
	}
	rec, err := readLaunchMigration(dir)
	if err != nil || len(rec.Inventory.Launches) != 1 || rec.Inventory.Launches[0].Classification != "retained-pre-start-refusal" {
		t.Fatalf("refusal provenance missing: %+v %v", rec, err)
	}
}

func TestLaunchMigrationPreservesAttemptsUnknownCostAndExtensions(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	cost := 0.0000001
	known := historicalMigrationCall(t, root, "process-exited", "", true, &cost)
	unknown := historicalMigrationCall(t, root, "failed", "start_failure", false, nil)
	historicalMigrationCall(t, root, "unobserved-handoff", "", false, nil)
	r := migrationFixtureRequest(t, root)
	r.AdditionalLaunches = 2
	r.Policy.MaxLaunches = 5
	status, err := MigrateLaunchBudget(ctx, root, "idea", r)
	if err != nil {
		t.Fatal(err)
	}
	if status.Spent != 5 || status.ExposureMicros != nil {
		t.Fatalf("lost unknown/extra attempts: %+v", status)
	}
	b, err := LoadLaunchBinding(ctx, root, "idea")
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Store.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if s.Entries[key(known.ID)].ActualMicros == nil || *s.Entries[key(known.ID)].ActualMicros != 1 || s.Entries[key(unknown.ID)].ActualMicros != nil {
		t.Fatalf("monetary provenance lost: %+v", s.Entries)
	}
	if _, err := b.Reserve(ctx, "sixth"); !errors.Is(err, ErrLimit) {
		t.Fatal("count reset", err)
	}
	if _, err := b.Store.ReconcileUnknown(ctx, unknown.ID, "ceiling", 10, "Fixture conservative bound"); err != nil {
		t.Fatal(err)
	}
	grant := PolicyExtensionRequest{DecisionID: "extension", ExpectedPolicySHA256: status.PolicySHA256, Reason: "Fixture grant after import", Ceilings: PolicyCeilings{Actions: 6}}
	if _, err := ExtendRuntimeBudget(ctx, root, "idea", Launch, grant); err != nil {
		t.Fatal(err)
	}
	status, err = MigrateLaunchBudget(ctx, root, "idea", r)
	if err != nil || status.Spent != 5 {
		t.Fatalf("exact import replay lost grant: %+v %v", status, err)
	}
	if _, err := b.Reserve(ctx, "sixth"); err != nil {
		t.Fatal(err)
	}
	current, _ := b.Store.Inspect(ctx)
	delete(current.Entries, key(unknown.ID))
	data, _ := json.Marshal(current)
	if err := writeSynced(filepath.Join(b.Store.Dir, "ledger.json"), data); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Inspect(ctx); err == nil {
		t.Fatal("deleted imported charge was ignored")
	}
}

func TestLaunchMigrationRecoversPublicationBeforeActivation(t *testing.T) {
	for _, name := range []string{"migration.json", "ledger.json", "policy.json", "migration-active"} {
		for _, after := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "-before", true: "-after"}[after], func(t *testing.T) {
				ctx := context.Background()
				root := t.TempDir()
				historicalMigrationCall(t, root, "failed", "start_failure", false, nil)
				r := migrationFixtureRequest(t, root)
				fired := false
				persist := func(path string, data []byte) error {
					if filepath.Base(path) == name && !fired {
						fired = true
						if after {
							if err := writeSynced(path, data); err != nil {
								return err
							}
						}
						return errors.New("injected persistence failure")
					}
					return writeSynced(path, data)
				}
				if _, err := migrateLaunchBudget(ctx, root, "idea", r, persist); err == nil || !fired {
					t.Fatalf("failure not exercised: %v", err)
				}
				b, loadErr := LoadLaunchBinding(ctx, root, "idea")
				if !(name == "migration-active" && after) && loadErr == nil && b != nil {
					t.Fatal("partial import authorized work")
				}
				status, err := MigrateLaunchBudget(ctx, root, "idea", r)
				if err != nil || status.Spent != 1 {
					t.Fatalf("replay failed: %+v %v", status, err)
				}
				b, err = LoadLaunchBinding(ctx, root, "idea")
				if err != nil {
					t.Fatal(err)
				}
				if _, err := os.Stat(b.Store.continuityPath()); err != nil {
					t.Fatal("retry omitted continuity", err)
				}
				if _, err := b.Reserve(ctx, "after-import"); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestLaunchMigrationStaleConflictingAndActiveDeletionRefuse(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	historicalMigrationCall(t, root, "failed", "budget_refused", false, nil)
	r := migrationFixtureRequest(t, root)
	historicalMigrationCall(t, root, "process-exited", "", true, nil)
	if _, err := MigrateLaunchBudget(ctx, root, "idea", r); err == nil {
		t.Fatal("stale history accepted")
	}
	b, err := LoadLaunchBinding(ctx, root, "idea")
	if b != nil || err != nil {
		t.Fatal("stale preview initialized state")
	}
	r = migrationFixtureRequest(t, root)
	if _, err := MigrateLaunchBudget(ctx, root, "idea", r); err != nil {
		t.Fatal(err)
	}
	changed := r
	changed.AdditionalLaunches++
	if _, err := MigrateLaunchBudget(ctx, root, "idea", changed); err == nil {
		t.Fatal("conflicting replay accepted")
	}
	b, _ = LoadLaunchBinding(ctx, root, "idea")
	policyPath := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	if err := os.Remove(policyPath); err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateLaunchBudget(ctx, root, "idea", r); err == nil {
		t.Fatal("active migration recreated lost policy")
	}
}

func TestLaunchMigrationConcurrentDecisionsHaveOneWinner(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	historicalMigrationCall(t, root, "failed", "budget_refused", false, nil)
	r := migrationFixtureRequest(t, root)
	barrier := make(chan struct{})
	var wg sync.WaitGroup
	results := make(chan bool, 2)
	for _, id := range []string{"first", "second"} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			request := r
			request.DecisionID = id
			<-barrier
			_, err := MigrateLaunchBudget(ctx, root, "idea", request)
			results <- err == nil
		}(id)
	}
	close(barrier)
	wg.Wait()
	close(results)
	wins := 0
	for win := range results {
		if win {
			wins++
		}
	}
	if wins != 1 {
		t.Fatalf("want one decision, got %d", wins)
	}
}

func TestLaunchMigrationCorruptWitnessAndReferencesRefuse(t *testing.T) {
	for _, mutation := range []string{"null-reference", "empty-reference", "missing-reference", "missing-witness", "changed-witness"} {
		t.Run(mutation, func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			r := migrationFixtureRequest(t, root)
			if _, err := MigrateLaunchBudget(ctx, root, "idea", r); err != nil {
				t.Fatal(err)
			}
			b, _ := LoadLaunchBinding(ctx, root, "idea")
			dir := filepath.Dir(b.Store.Dir)
			policy, _ := os.ReadFile(filepath.Join(dir, "policy.json"))
			var fields map[string]any
			_ = json.Unmarshal(policy, &fields)
			switch mutation {
			case "null-reference":
				fields["migration_sha256"] = nil
			case "empty-reference":
				fields["migration_sha256"] = ""
			case "missing-reference":
				delete(fields, "migration_sha256")
			case "missing-witness":
				_ = os.Remove(filepath.Join(dir, "migration.json"))
			case "changed-witness":
				raw, _ := os.ReadFile(filepath.Join(dir, "migration.json"))
				_ = os.WriteFile(filepath.Join(dir, "migration.json"), []byte(strings.Replace(string(raw), "fixture accounting", "forged accounting", 1)), 0600)
			}
			raw, _ := json.Marshal(fields)
			_ = os.WriteFile(filepath.Join(dir, "policy.json"), raw, 0600)
			if _, err := LoadLaunchBinding(ctx, root, "idea"); err == nil {
				t.Fatal("corrupted migration accepted")
			}
		})
	}
}

func TestLaunchMigrationInitialEpochAndScopeValidation(t *testing.T) {
	root := t.TempDir()
	historicalMigrationCall(t, root, "process-exited", "", true, nil)
	r := migrationFixtureRequest(t, root)
	inventory, err := InspectLaunchMigration(context.Background(), root, "idea")
	if err != nil {
		t.Fatal(err)
	}
	r.Policy.Version, r.Policy.Scope, r.Policy.Idea = 1, inventory.Scope, "idea"
	initial, err := migrationInitial(inventory, r)
	if err != nil {
		t.Fatal(err)
	}
	r.StartedAt = time.Now().UTC()
	if _, err := migrationInitial(inventory, r); err == nil {
		t.Fatal("migration reset original accounting epoch")
	}
	if initial.Entries == nil || reflect.DeepEqual(initial, Snapshot{}) {
		t.Fatal("empty initial accounting")
	}
}

func TestLaunchMigrationChangedSourcesDuringPublicationStayInactive(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	historicalMigrationCall(t, root, "failed", "budget_refused", false, nil)
	r := migrationFixtureRequest(t, root)
	persist := func(path string, data []byte) error {
		if err := writeSynced(path, data); err != nil {
			return err
		}
		if filepath.Base(path) == "policy.json" {
			historicalMigrationCall(t, root, "process-exited", "", true, nil)
		}
		return nil
	}
	if _, err := migrateLaunchBudget(ctx, root, "idea", r, persist); err == nil {
		t.Fatal("concurrent historical work silently excluded")
	}
	if _, err := LoadLaunchBinding(ctx, root, "idea"); err == nil {
		t.Fatal("changed sources activated import")
	}
	if _, err := MigrateLaunchBudget(ctx, root, "idea", r); err == nil {
		t.Fatal("stale recovery changed frozen accounting")
	}
}
