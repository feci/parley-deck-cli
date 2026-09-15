package budget

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

// AdditionalLaunches is an explicit operator accounting assertion for work
// preceding the invocation records. It is not inferred from artifact headings.
// Unknown costs for those attempts require the ordinary reconciliation control.
type LaunchMigrationRequest struct {
	ExpectedHistorySHA256 string       `json:"expected_history_sha256"`
	DecisionID            string       `json:"decision_id"`
	Reason                string       `json:"reason"`
	StartedAt             time.Time    `json:"started_at"`
	AdditionalLaunches    int          `json:"additional_launches"`
	WritersStopped        bool         `json:"writers_stopped"`
	Policy                LaunchPolicy `json:"policy"`
}

type launchMigrationRecord struct {
	Version    int                      `json:"version"`
	RecordedAt time.Time                `json:"recorded_at"`
	Request    LaunchMigrationRequest   `json:"request"`
	Inventory  LaunchMigrationInventory `json:"inventory"`
	Initial    Snapshot                 `json:"initial"`
}

func migrationInitial(i LaunchMigrationInventory, r LaunchMigrationRequest) (Snapshot, error) {
	state := Snapshot{Schema: 1, Scope: i.Scope, StartedAt: r.StartedAt.UTC(), Entries: map[string]Reservation{}}
	if i.Version != 1 || i.Scope == "" || len(i.Roots) == 0 || i.Sources == nil || i.Launches == nil || i.LegacyStartFloor < 0 || i.HistorySHA256 != i.digest() || !r.WritersStopped || !validCycleDecision(r.DecisionID, r.Reason, r.ExpectedHistorySHA256) || r.ExpectedHistorySHA256 != i.HistorySHA256 || r.StartedAt.IsZero() || r.StartedAt.After(time.Now().UTC()) || i.Earliest != nil && r.StartedAt.After(*i.Earliest) || r.AdditionalLaunches < i.LegacyStartFloor || r.AdditionalLaunches < 0 || r.AdditionalLaunches > maxMigrationItems {
		return state, errors.New("invalid migration decision, legacy floor, accounting epoch or history hash")
	}
	if r.Policy.Version != 1 || r.Policy.Scope != i.Scope || r.Policy.Idea != i.Idea || r.Policy.MigrationSHA256 != "" || r.Policy.Original != nil || r.Policy.Extensions != nil {
		return state, errors.New("migration must initialize an original launch policy, not change an existing grant")
	}
	if err := r.Policy.validate(); err != nil {
		return state, err
	}
	seen := map[string]bool{}
	for _, historical := range i.Launches {
		if historical.ID == "" || seen[historical.ID] || !validCycleDecision("source", "source", historical.TerminalSHA256) || historical.RequestedAt.Before(state.StartedAt) || historical.RequestedAt.After(time.Now().UTC()) || historical.ObservedMicros != nil && *historical.ObservedMicros < 0 {
			return state, errors.New("invalid imported invocation")
		}
		seen[historical.ID] = true
		switch historical.Classification {
		case "retained-pre-start-refusal":
			if historical.ObservedMicros != nil && *historical.ObservedMicros != 0 {
				return state, errors.New("refusal carries monetary work")
			}
			continue
		case "attempt", "unobserved-handoff":
		default:
			return state, errors.New("unsupported imported invocation classification")
		}
		state.Entries[key(historical.ID)] = Reservation{Kind: Launch, ReservedAt: historical.RequestedAt, Settled: true, ActualMicros: copyInt(historical.ObservedMicros)}
	}
	if len(state.Entries)+r.AdditionalLaunches > maxMigrationItems {
		return state, errors.New("imported action count exceeds migration bound")
	}
	for n := 0; n < r.AdditionalLaunches; n++ {
		id := key("legacy-migration:" + r.DecisionID + ":" + strconv.Itoa(n))
		if _, exists := state.Entries[id]; exists {
			return state, errors.New("import action identity collision")
		}
		state.Entries[id] = Reservation{Kind: Launch, ReservedAt: state.StartedAt, Settled: true}
	}
	return state, nil
}

func readLaunchMigration(dir string) (launchMigrationRecord, error) {
	var record launchMigrationRecord
	raw, err := readStepHistoryFile(filepath.Join(dir, "migration.json"), 8<<20)
	if err != nil {
		return record, err
	}
	if err := migrationJSON(raw, &record); err != nil {
		return record, err
	}
	if record.Version != 1 || record.RecordedAt.IsZero() || record.RecordedAt.Before(record.Request.StartedAt) || record.RecordedAt.After(time.Now().UTC()) {
		return record, errors.New("invalid migration record version or time")
	}
	initial, err := migrationInitial(record.Inventory, record.Request)
	if err != nil {
		return record, err
	}
	if !reflect.DeepEqual(initial, record.Initial) {
		return record, errors.New("migration ledger differs from its retained input")
	}
	return record, nil
}

func migrationActiveBytes(digest string) []byte {
	return []byte("parley-launch-migration-active/v1\n" + digest + "\n")
}

// The final activation marker makes pre-policy publication recoverable without
// treating a missing previously active policy as a first-time configuration.
func checkLaunchMigrationPolicy(dir string, policy LaunchPolicy) error {
	if policy.MigrationSHA256 == "" {
		return nil
	}
	r, err := readLaunchMigration(dir)
	if err != nil {
		return err
	}
	digest := migrationDigest(r)
	expected := r.Request.Policy
	expected.MigrationSHA256 = digest
	if policy.MigrationSHA256 != digest || !reflect.DeepEqual(policy.originalPolicy(), expected) {
		return errors.New("migrated policy differs from its immutable operator decision")
	}
	_, marker, err := recoveredMigrationState(launchRecoveryBase(dir, r))
	if err != nil {
		return err
	}
	raw, err := readLockOrigin(filepath.Join(dir, "migration-active"))
	if err != nil || string(raw) != string(marker) {
		return errors.New("launch migration is not durably active; replay the exact operator decision before work")
	}
	return nil
}

func checkLaunchMigrationCharges(dir string, p LaunchPolicy, state Snapshot) error {
	if p.MigrationSHA256 == "" {
		return nil
	}
	r, err := readLaunchMigration(dir)
	if err != nil {
		return err
	}
	initial, _, err := recoveredMigrationState(launchRecoveryBase(dir, r))
	if err != nil {
		return err
	}
	if migrationDigest(r) != p.MigrationSHA256 || !state.StartedAt.Equal(initial.StartedAt) {
		return errors.New("migration accounting epoch or witness differs")
	}
	for id, initial := range initial.Entries {
		current, ok := state.Entries[id]
		// Monetary reconciliation may be appended, but it cannot alter the
		// imported observation or erase the original spent action.
		current.Reconciliations = nil
		if !ok || !reflect.DeepEqual(current, initial) {
			return errors.New("imported launch charge was removed or rewritten")
		}
	}
	return nil
}

// MigrateLaunchBudget imports into an unconfigured scope only. The CLI must
// obtain a concrete operator decision and compute attendance independently.
func MigrateLaunchBudget(ctx context.Context, root, idea string, request LaunchMigrationRequest) (PolicyStatus, error) {
	return migrateLaunchBudget(ctx, root, idea, request, writeSynced)
}

func migrateLaunchBudget(ctx context.Context, root, idea string, request LaunchMigrationRequest, persist func(string, []byte) error) (PolicyStatus, error) {
	dir, scope, _, err := launchScope(ctx, root, idea, false)
	if err != nil {
		return PolicyStatus{}, err
	}
	if request.Policy.Version != 0 || request.Policy.Scope != "" || request.Policy.Idea != "" || request.Policy.MigrationSHA256 != "" || request.Policy.Original != nil || request.Policy.Extensions != nil {
		return PolicyStatus{}, errors.New("migration request cannot provide internal policy authority")
	}
	request.Policy.Version, request.Policy.Scope, request.Policy.Idea = 1, scope, idea
	request.StartedAt = request.StartedAt.UTC()
	if !request.WritersStopped || !validCycleDecision(request.DecisionID, request.Reason, request.ExpectedHistorySHA256) {
		return PolicyStatus{}, errors.New("migration requires the exact history hash, operator decision and stopped writers")
	}
	// Read and validate before creating any scope/lock state. Recovery skips
	// this first scan only for a retained, matching migration decision.
	prior, priorErr := readLaunchMigration(dir)
	var inventory LaunchMigrationInventory
	if priorErr == nil {
		if !reflect.DeepEqual(prior.Request, request) {
			return PolicyStatus{}, errors.New("conflicting migration decision; preserve the existing import")
		}
	} else if !os.IsNotExist(priorErr) {
		return PolicyStatus{}, priorErr
	} else {
		if _, err := os.Lstat(filepath.Join(dir, "policy.json")); err == nil || !os.IsNotExist(err) {
			return PolicyStatus{}, errors.New("configured launch policy cannot be replaced by migration")
		}
		if _, err := os.Lstat(filepath.Join(dir, "ledger", "ledger.json")); err == nil || !os.IsNotExist(err) {
			return PolicyStatus{}, errors.New("existing ledger requires recovery, not legacy import")
		}
		inventory, err = InspectLaunchMigration(ctx, root, idea)
		if err != nil {
			return PolicyStatus{}, err
		}
		if _, err := migrationInitial(inventory, request); err != nil {
			return PolicyStatus{}, err
		}
	}
	release, err := AcquireResourceGuard(ctx, dir)
	if err != nil {
		return PolicyStatus{}, err
	}
	defer release()
	prior, priorErr = readLaunchMigration(dir)
	if priorErr == nil && !reflect.DeepEqual(prior.Request, request) {
		return PolicyStatus{}, errors.New("a different migration won publication")
	}
	if priorErr != nil && !os.IsNotExist(priorErr) {
		return PolicyStatus{}, priorErr
	}
	if priorErr == nil {
		active, activeErr := readLockOrigin(filepath.Join(dir, "migration-active"))
		if activeErr == nil {
			_, marker, err := recoveredMigrationState(launchRecoveryBase(dir, prior))
			if err != nil {
				return PolicyStatus{}, err
			}
			if string(active) != string(marker) {
				return PolicyStatus{}, errors.New("migration activation witness differs")
			}
			binding, err := LoadLaunchBinding(ctx, root, idea)
			if err != nil || binding == nil {
				return PolicyStatus{}, fmt.Errorf("active import requires its retained policy and ledger: %w", err)
			}
			return binding.Inspect(ctx)
		}
		if !os.IsNotExist(activeErr) {
			return PolicyStatus{}, activeErr
		}
	}
	if err := refusePendingRecovery(dir); err != nil {
		return PolicyStatus{}, err
	}
	// Scope existence now prevents current readers from starting unbudgeted
	// work. Recheck all historical files after taking the publication guard.
	inventory, err = InspectLaunchMigration(ctx, root, idea)
	if err != nil {
		return PolicyStatus{}, err
	}
	initial, err := migrationInitial(inventory, request)
	if err != nil {
		return PolicyStatus{}, err
	}
	if priorErr != nil {
		for _, path := range []string{filepath.Join(dir, "policy.json"), filepath.Join(dir, "ledger", "ledger.json"), filepath.Join(dir, "migration-active")} {
			if _, err := os.Lstat(path); err == nil || !os.IsNotExist(err) {
				return PolicyStatus{}, errors.New("existing runtime state cannot be overwritten by a first import")
			}
		}
		prior = launchMigrationRecord{1, time.Now().UTC(), request, inventory, initial}
		raw, _ := json.MarshalIndent(prior, "", "  ")
		if len(raw) > 8<<20 {
			return PolicyStatus{}, errors.New("migration record exceeds size bound")
		}
		if err := persist(filepath.Join(dir, "migration.json"), append(raw, '\n')); err != nil {
			return PolicyStatus{}, err
		}
	}
	digest := migrationDigest(prior)
	policy := prior.Request.Policy
	policy.MigrationSHA256 = digest
	ledger := Store{Dir: filepath.Join(dir, "ledger"), Scope: scope, persist: persist}
	state, err := ledger.Inspect(ctx)
	if os.IsNotExist(err) {
		_, err = ledger.update(ctx, func(state *Snapshot, _ time.Time) error {
			if len(state.Entries) != 0 {
				return errors.New("concurrent charges prevent import initialization")
			}
			*state = prior.Initial
			return nil
		})
	} else if err == nil && !reflect.DeepEqual(state, prior.Initial) {
		return PolicyStatus{}, errors.New("inactive migration ledger differs from the frozen import")
	} else if err == nil {
		// A previous ledger replacement may have succeeded before its caller
		// received an error. Recheck under the ledger lock and establish the
		// continuity witness before activation can authorize any work.
		_, err = ledger.update(ctx, func(state *Snapshot, _ time.Time) error {
			if !reflect.DeepEqual(*state, prior.Initial) {
				return errors.New("inactive import changed during recovery")
			}
			return nil
		})
	}
	if err != nil {
		return PolicyStatus{}, err
	}
	// Once a policy exists before activation it must be this exact policy;
	// extensions cannot be issued through a still-inactive binding.
	old, err := readLaunchPolicyRaw(filepath.Join(dir, "policy.json"))
	if err == nil {
		if !reflect.DeepEqual(old, policy) {
			return PolicyStatus{}, errors.New("inactive migration has a conflicting policy")
		}
	} else if os.IsNotExist(err) {
		raw, _ := json.MarshalIndent(policy, "", "  ")
		if err := persist(filepath.Join(dir, "policy.json"), append(raw, '\n')); err != nil {
			return PolicyStatus{}, err
		}
	} else {
		return PolicyStatus{}, err
	}
	// The stopped-writers assertion is not a lock on historical artifacts.
	// Detect changes observed during publication and leave the import inactive.
	latest, err := InspectLaunchMigration(ctx, root, idea)
	if err != nil {
		return PolicyStatus{}, err
	}
	if latest.HistorySHA256 != request.ExpectedHistorySHA256 {
		return PolicyStatus{}, errors.New("historical sources changed during migration; preserve the inactive import for recovery")
	}
	if err := persist(filepath.Join(dir, "migration-active"), migrationActiveBytes(digest)); err != nil {
		return PolicyStatus{}, err
	}
	b := &LaunchBinding{Policy: policy, Store: ledger}
	return b.Inspect(ctx)
}
