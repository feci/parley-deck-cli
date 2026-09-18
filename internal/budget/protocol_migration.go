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
	"strings"
	"time"
)

// DeclaredUnavailable is the operator's request-scoped statement that these
// registered worktrees cannot be read and their history is unknown.
// DeclaredUnscopedRuns is the separate statement that these run directories,
// identified by their exact recursive file set, prove no idea identity — their
// bytes are readable and retained; only their scope is unknown. Both are
// arguments to this one decision: they are retained inside the immutable record
// for replay and recovery, and no ordinary bootstrap ever consults either.
//
// Both are omitempty, so a request carrying neither marshals exactly as before
// and keeps its digest. An older binary rejects a non-empty new field.
type ProtocolMigrationRequest struct {
	ExpectedHistorySHA256 string    `json:"expected_history_sha256"`
	DecisionID            string    `json:"decision_id"`
	Reason                string    `json:"reason"`
	IdeaPath              string    `json:"idea_path"`
	StartedAt             time.Time `json:"started_at"`
	TotalActions          int       `json:"total_actions"`
	Maximum               int       `json:"maximum"`
	WallClockNS           int64     `json:"wall_clock_ns"`
	WritersStopped        bool      `json:"writers_stopped"`
	DeclaredUnavailable   []string  `json:"declared_unavailable,omitempty"`
	DeclaredUnscopedRuns  []string  `json:"declared_unscoped_runs,omitempty"`
}

// Both declarations agree with the retained rows in both directions, and the
// inventory's coverage statement names exactly the combined basis its rows
// require. An undeclared inventory may claim no incomplete coverage at all.
func checkDeclarations(i ProtocolMigrationInventory, declared, unscoped []string) error {
	if err := checkDeclaredUnavailable(i, declared); err != nil {
		return err
	}
	if err := checkDeclaredUnscopedRuns(i, unscoped); err != nil {
		return err
	}
	basis := lowerBoundBasis(i.History)
	if basis == "" {
		if i.History.HistoryCoverage != "" || i.LowerBoundBasis != "" {
			return errors.New("undeclared protocol inventory claims incomplete history coverage")
		}
		return nil
	}
	if i.History.HistoryCoverage != DeclaredIncomplete || i.LowerBoundBasis != basis {
		return errors.New("declared protocol inventory does not record its incomplete coverage")
	}
	return nil
}

// A declared unscoped import must state exactly what it could not identify. The
// declaration and the retained rows agree in both directions, every row binds
// its declared manifest digest to at least one actual root with unknown
// history, and — the opposite of a declared unavailable worktree — the run's
// readable bytes are still retained as sources. The exclusion is from
// identity-scoped counting, never structural.
func checkDeclaredUnscopedRuns(i ProtocolMigrationInventory, declared []string) error {
	canonical, err := normalizeDeclaredUnscopedRuns(declared)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(canonical, declared) && len(declared) != 0 {
		return errors.New("declared-unscoped runs are not in canonical sorted deduplicated form")
	}
	rows := i.History.UnscopedRuns
	if len(rows) != len(canonical) {
		return errors.New("declared-unscoped runs differ from the retained unknown-identity rows")
	}
	for n, row := range rows {
		if row.Path+"="+row.SHA256 != canonical[n] || row.History != UnknownHistory || row.Copies != len(row.Roots) || row.Copies == 0 {
			return errors.New("unknown-identity row lost its declared run, manifest digest, unknown history or observed copies")
		}
		for _, origin := range row.Roots {
			if !retainedMigrationSource(i.History.Sources, origin, row.Path+"/events.jsonl") {
				return errors.New("declared-unscoped run lost the readable source it was retained for")
			}
		}
	}
	return nil
}

func retainedMigrationSource(sources []MigrationSource, root, path string) bool {
	for _, source := range sources {
		if source.Root == root && source.Path == path {
			return true
		}
	}
	return false
}

// A declared import must state exactly what it could not read. The declaration
// and the retained rows agree in both directions, every row binds an admissible
// stat class with unknown history, declared paths supply no source or root, and
// the floor announces that it covers the surviving visible worktrees only.
func checkDeclaredUnavailable(i ProtocolMigrationInventory, declared []string) error {
	canonical, err := normalizeDeclaredUnavailable(declared)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(canonical, declared) && len(declared) != 0 {
		return errors.New("declared-unavailable worktrees are not in canonical sorted deduplicated form")
	}
	rows := i.History.UnavailableRoots
	if len(rows) != len(canonical) {
		return errors.New("declared-unavailable worktrees differ from the retained unknown-history rows")
	}
	for n, row := range rows {
		if row.Path != canonical[n] || row.History != UnknownHistory || row.Observation != UnavailableMissing && row.Observation != UnavailableNotDirectory {
			return errors.New("unknown-history row lost its declared path, stat class or unknown history")
		}
	}
	// The coverage statement is one combined fact over both declaration kinds
	// and is checked by checkDeclarations, which owns it for exactly that
	// reason: this clause can only see its own kind.
	if len(rows) == 0 {
		return nil
	}
	// No evidence may be claimed from a path declared unreadable.
	for _, row := range rows {
		inside := func(path string) bool {
			return path == row.Path || strings.HasPrefix(path, row.Path+string(filepath.Separator))
		}
		for _, origin := range i.History.Roots {
			if inside(origin) {
				return errors.New("declared-unavailable worktree supplied a readable history root")
			}
		}
		for _, source := range i.History.Sources {
			if inside(source.Root) {
				return errors.New("declared-unavailable worktree supplied a history source")
			}
		}
	}
	return nil
}

type protocolMigrationRecord struct {
	Version    int                        `json:"version"`
	RecordedAt time.Time                  `json:"recorded_at"`
	Request    ProtocolMigrationRequest   `json:"request"`
	Inventory  ProtocolMigrationInventory `json:"inventory"`
	Initial    Snapshot                   `json:"initial"`
}

func protocolMigrationInitial(i ProtocolMigrationInventory, r ProtocolMigrationRequest) (Snapshot, error) {
	s := Snapshot{Schema: 1, Scope: i.Scope, StartedAt: r.StartedAt.UTC(), Entries: map[string]Reservation{}}
	if i.Version != 1 || i.Scope == "" || i.Idea == "" || i.Kind != DriverStep && i.Kind != Fixup && i.Kind != CrossReview || i.HistorySHA256 != i.digest() || i.History.HistorySHA256 != i.History.digest() || i.History.Idea != i.Idea || i.Evidence == nil || i.UngroupedInvocations == nil || i.LowerBound < 0 || i.LowerBound > maxMigrationItems {
		return s, errors.New("invalid protocol migration inventory")
	}
	if !r.WritersStopped || !validCycleDecision(r.DecisionID, r.Reason, r.ExpectedHistorySHA256) || r.ExpectedHistorySHA256 != i.HistorySHA256 || r.IdeaPath != i.IdeaPath || r.TotalActions < i.LowerBound || r.TotalActions < 0 || r.TotalActions > maxMigrationItems || r.StartedAt.IsZero() || r.StartedAt.After(time.Now().UTC()) || i.Earliest != nil && r.StartedAt.After(*i.Earliest) || r.Maximum < 0 || r.WallClockNS < 0 || i.Kind != DriverStep && r.WallClockNS != 0 {
		return s, errors.New("invalid protocol accounting decision, epoch, total or ceiling")
	}
	if _, err := protocolMigrationPath(i.Idea, i.IdeaPath); err != nil {
		return s, err
	}
	if err := checkDeclarations(i, r.DeclaredUnavailable, r.DeclaredUnscopedRuns); err != nil {
		return s, err
	}
	if floor, err := protocolEvidenceFloor(i.Kind, i.Evidence); err != nil || floor != i.LowerBound {
		return s, errors.New("protocol import count differs from retained evidence")
	}
	if i.History.Earliest != nil && (i.Earliest == nil || i.Earliest.After(*i.History.Earliest)) {
		return s, errors.New("protocol import lost its earliest observed accounting time")
	}
	// These entries record the explicit reconciled protocol count. Their epoch
	// is declared accounting time, not an invented observed child start. Model
	// costs live in the separate launch ledger and remain unknown when absent.
	zero := int64(0)
	for n := 0; n < r.TotalActions; n++ {
		id := key("protocol-migration:" + string(i.Kind) + ":" + r.DecisionID + ":" + strconv.Itoa(n))
		s.Entries[id] = Reservation{Kind: i.Kind, ReservedAt: s.StartedAt, Settled: true, ReserveMicros: &zero, ActualMicros: &zero}
	}
	return s, nil
}

func protocolMigrationPolicy(r protocolMigrationRecord, digest string) any {
	i, q := r.Inventory, r.Request
	if i.Kind == DriverStep {
		return StepPolicy{Version: 1, Scope: i.Scope, Idea: i.Idea, MaxSteps: q.Maximum, WallClockNS: q.WallClockNS, MigrationSHA256: digest}
	}
	return CyclePolicy{Version: 1, Scope: i.Scope, Idea: i.Idea, IdeaPath: i.IdeaPath, Kind: i.Kind, Maximum: q.Maximum, Carried: 0, MigrationSHA256: digest}
}

func readProtocolMigration(dir string) (protocolMigrationRecord, error) {
	var r protocolMigrationRecord
	raw, err := readStepHistoryFile(filepath.Join(dir, "migration.json"), 8<<20)
	if err != nil {
		return r, err
	}
	if err := migrationJSON(raw, &r); err != nil {
		return r, err
	}
	if r.Version != 1 || r.RecordedAt.IsZero() || r.RecordedAt.Before(r.Request.StartedAt) || r.RecordedAt.After(time.Now().UTC()) {
		return r, errors.New("invalid protocol import version/time")
	}
	initial, err := protocolMigrationInitial(r.Inventory, r.Request)
	if err != nil {
		return r, err
	}
	if !reflect.DeepEqual(initial, r.Initial) {
		return r, errors.New("protocol import ledger differs from its retained decision")
	}
	return r, nil
}

func protocolMigrationActive(digest string) []byte {
	return []byte("parley-protocol-migration-active/v1\n" + digest + "\n")
}

func checkProtocolMigrationPolicy(dir, digest string, original any) error {
	if digest == "" {
		for _, name := range []string{"migration.json", "migration-active"} {
			if _, err := os.Lstat(filepath.Join(dir, name)); err == nil || !os.IsNotExist(err) {
				return errors.New("protocol policy lost its migration reference")
			}
		}
		return nil
	}
	if !validCycleDecision("migration", "migration", digest) {
		return errors.New("invalid protocol migration reference")
	}
	r, err := readProtocolMigration(dir)
	if err != nil {
		return err
	}
	if migrationDigest(r) != digest || !reflect.DeepEqual(original, protocolMigrationPolicy(r, digest)) {
		return errors.New("protocol policy differs from immutable migration decision")
	}
	_, marker, err := recoveredMigrationState(protocolRecoveryBase(dir, r))
	if err != nil {
		return err
	}
	active, err := readLockOrigin(filepath.Join(dir, "migration-active"))
	if err != nil || string(active) != string(marker) {
		return errors.New("protocol migration is not durably active; replay the exact decision before work")
	}
	return nil
}

func checkProtocolMigrationCharges(dir, digest string, s Snapshot) error {
	if digest == "" {
		return nil
	}
	r, err := readProtocolMigration(dir)
	if err != nil {
		return err
	}
	initial, _, err := recoveredMigrationState(protocolRecoveryBase(dir, r))
	if err != nil {
		return err
	}
	if migrationDigest(r) != digest || s.Scope != initial.Scope || !s.StartedAt.Equal(initial.StartedAt) {
		return errors.New("protocol import lost its accounting origin")
	}
	for id, old := range initial.Entries {
		if current, ok := s.Entries[id]; !ok || !reflect.DeepEqual(current, old) {
			return errors.New("imported protocol charge was removed or rewritten")
		}
	}
	return nil
}

func inspectProtocolMigrationBinding(ctx context.Context, root, idea string, kind Kind) (PolicyStatus, error) {
	if kind == DriverStep {
		return InspectRuntimeBudget(ctx, root, idea, kind)
	}
	s, err := InspectCycleBudget(ctx, root, idea, kind)
	if err != nil {
		return PolicyStatus{}, err
	}
	raw, err := json.Marshal(s.Policy)
	zero := int64(0)
	return PolicyStatus{Kind: kind, Policy: raw, PolicySHA256: s.PolicySHA256, Spent: s.Spent, ExposureMicros: &zero, StartedAt: s.StartedAt, LedgerDir: s.LedgerDir}, err
}

// MigrateProtocolBudget is the explicit import control for steps and cycles.
// The application independently establishes attendance; this API cannot prove
// who made a decision. It never replaces an already configured/charged scope.
func MigrateProtocolBudget(ctx context.Context, root, idea string, kind Kind, r ProtocolMigrationRequest) (PolicyStatus, error) {
	return migrateProtocolBudget(ctx, root, idea, kind, r, writeSynced)
}

func migrateProtocolBudget(ctx context.Context, root, idea string, kind Kind, r ProtocolMigrationRequest, persist func(string, []byte) error) (PolicyStatus, error) {
	dir, _, err := protocolMigrationScope(ctx, root, idea, kind)
	if err != nil {
		return PolicyStatus{}, err
	}
	r.IdeaPath, err = protocolMigrationPath(idea, r.IdeaPath)
	if err != nil {
		return PolicyStatus{}, err
	}
	r.StartedAt = r.StartedAt.UTC()
	// Canonicalize both declarations here, at the one boundary, so an unsorted,
	// duplicated, nil or empty-slice spelling of the same decision is one exact
	// replay, not a conflicting one. A new declaration kind is canonicalized
	// beside the old one and never on a separate path.
	if r.DeclaredUnavailable, err = normalizeDeclaredUnavailable(r.DeclaredUnavailable); err != nil {
		return PolicyStatus{}, err
	}
	if r.DeclaredUnscopedRuns, err = normalizeDeclaredUnscopedRuns(r.DeclaredUnscopedRuns); err != nil {
		return PolicyStatus{}, err
	}
	if !r.WritersStopped || !validCycleDecision(r.DecisionID, r.Reason, r.ExpectedHistorySHA256) {
		return PolicyStatus{}, errors.New("protocol migration requires exact history and an explicit stopped-writer decision")
	}
	prior, priorErr := readProtocolMigration(dir)
	if priorErr == nil {
		if prior.Inventory.Idea != idea || prior.Inventory.Kind != kind || !reflect.DeepEqual(prior.Request, r) {
			return PolicyStatus{}, errors.New("conflicting protocol migration replay")
		}
	} else if !os.IsNotExist(priorErr) {
		return PolicyStatus{}, priorErr
	} else {
		for _, name := range []string{"policy.json", "ledger/ledger.json", "migration-active"} {
			if _, err := os.Lstat(filepath.Join(dir, filepath.FromSlash(name))); err == nil || !os.IsNotExist(err) {
				return PolicyStatus{}, errors.New("existing protocol accounting requires recovery, not first import")
			}
		}
		i, err := InspectProtocolMigrationDeclarations(ctx, root, idea, kind, r.IdeaPath, r.DeclaredUnavailable, r.DeclaredUnscopedRuns)
		if err != nil {
			return PolicyStatus{}, err
		}
		if _, err := protocolMigrationInitial(i, r); err != nil {
			return PolicyStatus{}, err
		}
	}
	release, err := AcquireResourceGuard(ctx, dir)
	if err != nil {
		return PolicyStatus{}, err
	}
	defer release()
	prior, priorErr = readProtocolMigration(dir)
	if priorErr != nil && !os.IsNotExist(priorErr) {
		return PolicyStatus{}, priorErr
	}
	if priorErr == nil {
		if prior.Inventory.Idea != idea || prior.Inventory.Kind != kind || !reflect.DeepEqual(prior.Request, r) {
			return PolicyStatus{}, errors.New("another protocol migration decision won publication")
		}
		active, err := readLockOrigin(filepath.Join(dir, "migration-active"))
		if err == nil {
			_, marker, err := recoveredMigrationState(protocolRecoveryBase(dir, prior))
			if err != nil {
				return PolicyStatus{}, err
			}
			if string(active) != string(marker) {
				return PolicyStatus{}, errors.New("protocol activation witness differs")
			}
			return inspectProtocolMigrationBinding(ctx, root, idea, kind)
		}
		if !os.IsNotExist(err) {
			return PolicyStatus{}, err
		}
	}
	if err := refusePendingRecovery(dir); err != nil {
		return PolicyStatus{}, err
	}
	i, err := InspectProtocolMigrationDeclarations(ctx, root, idea, kind, r.IdeaPath, r.DeclaredUnavailable, r.DeclaredUnscopedRuns)
	if err != nil {
		return PolicyStatus{}, err
	}
	initial, err := protocolMigrationInitial(i, r)
	if err != nil {
		return PolicyStatus{}, err
	}
	if priorErr != nil {
		for _, name := range []string{"policy.json", "ledger/ledger.json", "migration-active"} {
			if _, err := os.Lstat(filepath.Join(dir, filepath.FromSlash(name))); err == nil || !os.IsNotExist(err) {
				return PolicyStatus{}, errors.New("concurrent accounting cannot be overwritten by protocol import")
			}
		}
		prior = protocolMigrationRecord{Version: 1, RecordedAt: time.Now().UTC(), Request: r, Inventory: i, Initial: initial}
		raw, err := json.MarshalIndent(prior, "", "  ")
		if err != nil {
			return PolicyStatus{}, err
		}
		if len(raw) > 8<<20 {
			return PolicyStatus{}, errors.New("protocol import exceeds record bound")
		}
		if err := persist(filepath.Join(dir, "migration.json"), append(raw, '\n')); err != nil {
			return PolicyStatus{}, err
		}
	}
	digest := migrationDigest(prior)
	ledger := Store{Dir: filepath.Join(dir, "ledger"), Scope: i.Scope, persist: persist}
	current, err := ledger.Inspect(ctx)
	if os.IsNotExist(err) {
		_, err = ledger.update(ctx, func(s *Snapshot, _ time.Time) error {
			if len(s.Entries) != 0 {
				return errors.New("concurrent charges prevent protocol import")
			}
			*s = prior.Initial
			return nil
		})
	} else if err == nil {
		if !reflect.DeepEqual(current, prior.Initial) {
			return PolicyStatus{}, errors.New("inactive protocol ledger differs from frozen import")
		}
		_, err = ledger.update(ctx, func(s *Snapshot, _ time.Time) error {
			if !reflect.DeepEqual(*s, prior.Initial) {
				return errors.New("inactive import changed during recovery")
			}
			return nil
		})
	}
	if err != nil {
		return PolicyStatus{}, err
	}
	policy := protocolMigrationPolicy(prior, digest)
	path := filepath.Join(dir, "policy.json")
	raw, err := readStepHistoryFile(path, 1<<20)
	if err == nil {
		var old any
		if kind == DriverStep {
			var p StepPolicy
			err = migrationJSON(raw, &p)
			old = p
		} else {
			var p CyclePolicy
			err = migrationJSON(raw, &p)
			old = p
		}
		if err != nil {
			return PolicyStatus{}, err
		}
		if !reflect.DeepEqual(old, policy) {
			return PolicyStatus{}, errors.New("inactive protocol policy differs")
		}
	} else if os.IsNotExist(err) {
		raw, err = json.MarshalIndent(policy, "", "  ")
		if err != nil {
			return PolicyStatus{}, err
		}
		if err := persist(path, append(raw, '\n')); err != nil {
			return PolicyStatus{}, err
		}
	} else {
		return PolicyStatus{}, err
	}
	// The pre-activation re-inspection carries both declarations: a declared
	// path that became readable, lost its registration or changed stat class,
	// and a declared run whose file set gained, lost or changed a byte in any
	// copy, recovered an identity or vanished from every root, refuse here,
	// before the activation witness is written.
	latest, err := InspectProtocolMigrationDeclarations(ctx, root, idea, kind, r.IdeaPath, r.DeclaredUnavailable, r.DeclaredUnscopedRuns)
	if err != nil {
		return PolicyStatus{}, err
	}
	if latest.HistorySHA256 != r.ExpectedHistorySHA256 {
		return PolicyStatus{}, fmt.Errorf("%w before protocol import activation", errHistoryChanged)
	}
	if err := persist(filepath.Join(dir, "migration-active"), protocolMigrationActive(digest)); err != nil {
		return PolicyStatus{}, err
	}
	return inspectProtocolMigrationBinding(ctx, root, idea, kind)
}
