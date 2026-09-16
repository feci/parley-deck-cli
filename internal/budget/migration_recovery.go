package budget

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

const migrationRecoveryFile = "migration-recovery.json"
const maxMigrationRecoveries = 32

// AccountedActions is the total pre-telemetry launch count for launch imports,
// or the total historical step/cycle count for protocol imports. It never means
// a new allowance. The original ceilings remain unchanged.
type MigrationRecoveryRequest struct {
	ExpectedImportSHA256  string    `json:"expected_import_sha256"`
	ExpectedHistorySHA256 string    `json:"expected_history_sha256"`
	DecisionID            string    `json:"decision_id"`
	Reason                string    `json:"reason"`
	StartedAt             time.Time `json:"started_at"`
	AccountedActions      int       `json:"accounted_actions"`
	WritersStopped        bool      `json:"writers_stopped"`
}

type migrationImportFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"` // empty is an explicitly absent file
}

var migrationImportPaths = []string{
	"migration.json", migrationRecoveryFile, "migration-active", "policy.json",
	"ledger/ledger.json", "ledger/ledger-established",
}

type migrationRecoveryDecision struct {
	RecordedAt     time.Time                   `json:"recorded_at"`
	PreviousSHA256 string                      `json:"previous_sha256"`
	Request        MigrationRecoveryRequest    `json:"request"`
	ImportFiles    []migrationImportFile       `json:"import_files"`
	Launch         *LaunchMigrationInventory   `json:"launch,omitempty"`
	Protocol       *ProtocolMigrationInventory `json:"protocol,omitempty"`
	Initial        Snapshot                    `json:"initial"`
}

type migrationRecoveryJournal struct {
	Version             int                         `json:"version"`
	BaseMigrationSHA256 string                      `json:"base_migration_sha256"`
	Scope               string                      `json:"scope"`
	Kind                Kind                        `json:"kind"`
	Decisions           []migrationRecoveryDecision `json:"decisions"`
}

type MigrationRecoveryPreview struct {
	ImportDirectory         string                      `json:"import_directory"`
	OriginalDecisionID      string                      `json:"original_decision_id"`
	LastDecision            *MigrationRecoveryRequest   `json:"last_decision,omitempty"`
	Kind                    Kind                        `json:"kind"`
	Scope                   string                      `json:"scope"`
	Idea                    string                      `json:"idea"`
	ImportSHA256            string                      `json:"import_sha256"`
	HistorySHA256           string                      `json:"history_sha256"`
	Active                  bool                        `json:"active"`
	CountBasis              string                      `json:"count_basis"`
	PriorAccountedActions   int                         `json:"prior_accounted_actions"`
	MinimumAccountedActions int                         `json:"minimum_accounted_actions"`
	Retained                Snapshot                    `json:"retained"`
	Policy                  json.RawMessage             `json:"policy"`
	Launch                  *LaunchMigrationInventory   `json:"launch,omitempty"`
	Protocol                *ProtocolMigrationInventory `json:"protocol,omitempty"`
}

type recoveryBase struct {
	dir, scope, idea, ideaPath, digest, decisionID string
	kind                                           Kind
	at                                             time.Time
	count                                          int
	initial                                        Snapshot
	policy                                         any
	launch                                         *launchMigrationRecord
	protocol                                       *protocolMigrationRecord
	// declared and unscoped are the original import's request-scoped
	// declarations, read back from the immutable record. Recovery re-inspects
	// the same history, so without both a declared import could never be
	// recovered — and a declaration silently dropped here would re-inspect a
	// different history and refuse, or worse, a narrower one.
	declared []string
	unscoped []string
}

func launchRecoveryBase(dir string, r launchMigrationRecord) recoveryBase {
	digest := migrationDigest(r)
	p := r.Request.Policy
	p.MigrationSHA256 = digest
	return recoveryBase{dir: dir, scope: r.Initial.Scope, idea: r.Inventory.Idea,
		digest: digest, decisionID: r.Request.DecisionID, kind: Launch,
		at: r.RecordedAt, count: r.Request.AdditionalLaunches, initial: r.Initial,
		policy: p, launch: &r}
}

func protocolRecoveryBase(dir string, r protocolMigrationRecord) recoveryBase {
	digest := migrationDigest(r)
	return recoveryBase{dir: dir, scope: r.Initial.Scope, idea: r.Inventory.Idea,
		ideaPath: r.Inventory.IdeaPath, digest: digest, decisionID: r.Request.DecisionID,
		kind: r.Inventory.Kind, at: r.RecordedAt, count: r.Request.TotalActions,
		initial: r.Initial, policy: protocolMigrationPolicy(r, digest), protocol: &r,
		declared: r.Request.DeclaredUnavailable, unscoped: r.Request.DeclaredUnscopedRuns}
}

func loadRecoveryBase(ctx context.Context, root, idea string, kind Kind) (recoveryBase, error) {
	if kind == Launch {
		dir, scope, _, err := launchScope(ctx, root, idea, false)
		if err != nil {
			return recoveryBase{}, err
		}
		r, err := readLaunchMigration(dir)
		if err != nil {
			return recoveryBase{}, err
		}
		if r.Initial.Scope != scope || r.Inventory.Idea != idea {
			return recoveryBase{}, errors.New("recovery import scope differs")
		}
		return launchRecoveryBase(dir, r), nil
	}
	dir, scope, err := protocolMigrationScope(ctx, root, idea, kind)
	if err != nil {
		return recoveryBase{}, err
	}
	r, err := readProtocolMigration(dir)
	if err != nil {
		return recoveryBase{}, err
	}
	if r.Initial.Scope != scope || r.Inventory.Idea != idea || r.Inventory.Kind != kind {
		return recoveryBase{}, errors.New("recovery protocol scope differs")
	}
	return protocolRecoveryBase(dir, r), nil
}

func recoveryPrefixDigest(j migrationRecoveryJournal) string {
	if len(j.Decisions) == 0 {
		return j.BaseMigrationSHA256
	}
	return migrationDigest(j)
}

func readRecoveryJournal(b recoveryBase) (migrationRecoveryJournal, error) {
	j := migrationRecoveryJournal{Version: 1, BaseMigrationSHA256: b.digest, Scope: b.scope, Kind: b.kind, Decisions: []migrationRecoveryDecision{}}
	raw, err := readStepHistoryFile(filepath.Join(b.dir, migrationRecoveryFile), 16<<20)
	if os.IsNotExist(err) {
		return j, nil
	}
	if err != nil {
		return j, err
	}
	if err := migrationJSON(raw, &j); err != nil {
		return j, err
	}
	if j.Version != 1 || j.BaseMigrationSHA256 != b.digest || j.Scope != b.scope || j.Kind != b.kind || len(j.Decisions) == 0 || len(j.Decisions) > maxMigrationRecoveries {
		return j, errors.New("invalid recovery journal binding or length")
	}
	previous, count, at := b.initial, b.count, b.at
	seen := map[string]bool{b.decisionID: true}
	for n, d := range j.Decisions {
		prefix := j
		prefix.Decisions = j.Decisions[:n]
		if d.PreviousSHA256 != recoveryPrefixDigest(prefix) || d.RecordedAt.Before(at) || d.RecordedAt.IsZero() || d.RecordedAt.After(time.Now().UTC()) || seen[d.Request.DecisionID] {
			return j, errors.New("invalid recovery decision chain or time")
		}
		if err := validateImportFiles(d.ImportFiles); err != nil || migrationDigest(d.ImportFiles) != d.Request.ExpectedImportSHA256 {
			return j, errors.New("recovery decision lost its observed import binding")
		}
		if err := checkRecoveryObservations(b, prefix, d.Launch, d.Protocol); err != nil {
			return j, err
		}
		expected, err := recoveredInitial(b, previous, count, d.Request, d.Launch, d.Protocol)
		if err != nil {
			return j, err
		}
		if !reflect.DeepEqual(expected, d.Initial) {
			return j, errors.New("recovery charges differ from retained decisions")
		}
		seen[d.Request.DecisionID] = true
		previous, count, at = expected, d.Request.AccountedActions, d.RecordedAt
	}
	return j, nil
}

// An immutable terminal observation cannot silently turn from a refused launch
// into a spent attempt, or change provenance while keeping the same amount.
// Every earlier observation remains available in the original or prior journal.
func checkRecoveryObservations(b recoveryBase, prefix migrationRecoveryJournal, launch *LaunchMigrationInventory, protocol *ProtocolMigrationInventory) error {
	seen := map[string]HistoricalLaunch{}
	add := func(rows []HistoricalLaunch) error {
		for _, row := range rows {
			if prior, ok := seen[row.ID]; ok && !reflect.DeepEqual(prior, row) {
				return errors.New("previously imported terminal observation changed; preserve both views for separate reconciliation")
			}
			seen[row.ID] = row
		}
		return nil
	}
	if b.kind == Launch {
		if err := add(b.launch.Inventory.Launches); err != nil {
			return err
		}
	} else if err := add(b.protocol.Inventory.History.Launches); err != nil {
		return err
	}
	for _, d := range prefix.Decisions {
		if d.Launch != nil {
			if err := add(d.Launch.Launches); err != nil {
				return err
			}
		}
		if d.Protocol != nil {
			if err := add(d.Protocol.History.Launches); err != nil {
				return err
			}
		}
	}
	if launch != nil {
		if err := add(launch.Launches); err != nil {
			return err
		}
	}
	if protocol != nil {
		return add(protocol.History.Launches)
	}
	return nil
}

func recoveredInitial(b recoveryBase, previous Snapshot, oldCount int, r MigrationRecoveryRequest, launch *LaunchMigrationInventory, protocol *ProtocolMigrationInventory) (Snapshot, error) {
	if !r.WritersStopped || !validCycleDecision(r.DecisionID, r.Reason, r.ExpectedHistorySHA256) || !validCycleDecision("import", "import", r.ExpectedImportSHA256) || r.AccountedActions < oldCount || r.AccountedActions > maxMigrationItems || r.StartedAt.IsZero() || r.StartedAt.After(previous.StartedAt) {
		return Snapshot{}, errors.New("recovery needs exact hashes, stopped writers, nondecreasing accounting and no later epoch")
	}
	var candidate Snapshot
	var err error
	legacyPrefix := ""
	if b.kind == Launch {
		if launch == nil || protocol != nil || launch.Scope != b.scope || launch.Idea != b.idea {
			return Snapshot{}, errors.New("recovery launch inventory scope differs")
		}
		q := b.launch.Request
		q.ExpectedHistorySHA256, q.StartedAt, q.AdditionalLaunches = r.ExpectedHistorySHA256, r.StartedAt, r.AccountedActions
		candidate, err = migrationInitial(*launch, q)
		legacyPrefix = "legacy-migration:" + b.decisionID + ":"
	} else {
		if protocol == nil || launch != nil || protocol.Scope != b.scope || protocol.Idea != b.idea || protocol.Kind != b.kind || protocol.IdeaPath != b.ideaPath {
			return Snapshot{}, errors.New("recovery protocol inventory scope differs")
		}
		q := b.protocol.Request
		q.ExpectedHistorySHA256, q.StartedAt, q.TotalActions = r.ExpectedHistorySHA256, r.StartedAt, r.AccountedActions
		candidate, err = protocolMigrationInitial(*protocol, q)
		legacyPrefix = "protocol-migration:" + string(b.kind) + ":" + b.decisionID + ":"
	}
	if err != nil {
		return Snapshot{}, err
	}
	// Retain exact prior observations, including disappeared historical sources.
	// An earlier accounting epoch does not rewrite the timestamps on old entries.
	state := Snapshot{Schema: previous.Schema, Scope: previous.Scope, StartedAt: r.StartedAt.UTC(), Entries: map[string]Reservation{}}
	for id, e := range previous.Entries {
		state.Entries[id] = e
	}
	accounted := map[string]bool{}
	for n := 0; n < oldCount; n++ {
		accounted[key(legacyPrefix+strconv.Itoa(n))] = true
	}
	for id, e := range candidate.Entries {
		if old, exists := state.Entries[id]; exists {
			if !accounted[id] && !reflect.DeepEqual(old, e) {
				return Snapshot{}, errors.New("previously imported invocation observation changed; preserve it for separate reconciliation")
			}
			continue
		}
		state.Entries[id] = e
	}
	if len(state.Entries) > maxMigrationItems {
		return Snapshot{}, errors.New("recovered import exceeds action bound")
	}
	return state, nil
}

// Legacy activation has a different shape so older runtimes cannot treat a
// recovered ledger as if it were only the original import.
func recoveredMigrationState(b recoveryBase) (Snapshot, []byte, error) {
	j, err := readRecoveryJournal(b)
	if err != nil {
		return Snapshot{}, nil, err
	}
	if len(j.Decisions) == 0 {
		if b.kind == Launch {
			return b.initial, migrationActiveBytes(b.digest), nil
		}
		return b.initial, protocolMigrationActive(b.digest), nil
	}
	return j.Decisions[len(j.Decisions)-1].Initial, []byte("parley-budget-migration-recovery-active/v1\n" + b.digest + "\n" + migrationDigest(j) + "\n"), nil
}

func refusePendingRecovery(dir string) error {
	if _, err := os.Lstat(filepath.Join(dir, migrationRecoveryFile)); err == nil {
		return errors.New("pending migration recovery requires its exact recovery decision, not original import replay")
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func validateImportFiles(files []migrationImportFile) error {
	if len(files) != len(migrationImportPaths) {
		return errors.New("incomplete import file binding")
	}
	for n, f := range files {
		if f.Path != migrationImportPaths[n] || f.SHA256 != "" && !validCycleDecision("file", "file", f.SHA256) {
			return errors.New("invalid import file binding")
		}
	}
	return nil
}

func inspectImportFiles(dir string) ([]migrationImportFile, error) {
	files := make([]migrationImportFile, 0, len(migrationImportPaths))
	for _, name := range migrationImportPaths {
		if err := migrationDirectory(dir, filepath.ToSlash(filepath.Dir(name))); err != nil {
			return nil, err
		}
		data, err := readStepHistoryFile(filepath.Join(dir, filepath.FromSlash(name)), 16<<20)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		f := migrationImportFile{Path: name}
		if err == nil {
			f.SHA256 = key(string(data))
		}
		files = append(files, f)
	}
	return files, nil
}

type recoveryObservation struct {
	base    recoveryBase
	journal migrationRecoveryJournal
	files   []migrationImportFile
	preview MigrationRecoveryPreview
}

func policyBeforeRecovery(b recoveryBase, active bool) error {
	raw, err := readStepHistoryFile(filepath.Join(b.dir, "policy.json"), 1<<20)
	if os.IsNotExist(err) && !active {
		return nil
	}
	if err != nil {
		return err
	}
	var policy any
	switch b.kind {
	case Launch:
		var p LaunchPolicy
		err = migrationJSON(raw, &p)
		policy = p
		if active {
			policy = p.originalPolicy()
		}
	case DriverStep:
		var p StepPolicy
		err = migrationJSON(raw, &p)
		policy = p
		if active {
			policy = p.originalPolicy()
		}
	default:
		var p CyclePolicy
		err = migrationJSON(raw, &p)
		policy = p
		if active {
			policy = originalCyclePolicy(p)
		}
	}
	if err != nil || !reflect.DeepEqual(policy, b.policy) {
		return errors.New("recovery cannot replace a conflicting policy or extend an inactive import")
	}
	return nil
}

func recoveryLedgerMatches(b recoveryBase, j migrationRecoveryJournal, s Snapshot) bool {
	if reflect.DeepEqual(s, b.initial) {
		return true
	}
	for _, d := range j.Decisions {
		if reflect.DeepEqual(s, d.Initial) {
			return true
		}
	}
	return false
}

func observeMigrationRecovery(ctx context.Context, root, idea string, kind Kind) (recoveryObservation, error) {
	var o recoveryObservation
	if err := ctx.Err(); err != nil {
		return o, err
	}
	b, err := loadRecoveryBase(ctx, root, idea, kind)
	if err != nil {
		return o, err
	}
	o.base = b
	o.files, err = inspectImportFiles(b.dir)
	if err != nil {
		return o, err
	}
	// Re-read the original after the file snapshot, then verify the snapshot
	// again. A read-only preview is observational, never a grant or a lock.
	b, err = loadRecoveryBase(ctx, root, idea, kind)
	if err != nil {
		return o, err
	}
	o.base = b
	o.journal, err = readRecoveryJournal(b)
	if err != nil {
		return o, err
	}
	initial, marker, err := recoveredMigrationState(b)
	if err != nil {
		return o, err
	}
	active, err := readLockOrigin(filepath.Join(b.dir, "migration-active"))
	if err != nil && !os.IsNotExist(err) {
		return o, err
	}
	isActive := err == nil
	if isActive && string(active) != string(marker) {
		return o, errors.New("migration activation differs from retained recovery")
	}
	if err := policyBeforeRecovery(b, isActive); err != nil {
		return o, err
	}
	ledger := Store{Dir: filepath.Join(b.dir, "ledger"), Scope: b.scope}
	current, err := ledger.Inspect(ctx)
	if os.IsNotExist(err) && !isActive {
		if err := ledger.checkContinuity(true); err != nil {
			return o, err
		}
	} else if err != nil {
		return o, err
	} else if !isActive && !recoveryLedgerMatches(b, o.journal, current) {
		return o, errors.New("inactive ledger differs from every retained import checkpoint")
	}
	count := b.count
	if len(o.journal.Decisions) > 0 {
		count = o.journal.Decisions[len(o.journal.Decisions)-1].Request.AccountedActions
	}
	policy, err := json.Marshal(b.policy)
	if err != nil {
		return o, err
	}
	o.preview = MigrationRecoveryPreview{Kind: kind, Scope: b.scope, Idea: idea, ImportSHA256: migrationDigest(o.files), Active: isActive, PriorAccountedActions: count, MinimumAccountedActions: count, Retained: initial, Policy: policy, CountBasis: "historical-protocol-total"}
	o.preview.ImportDirectory, o.preview.OriginalDecisionID = b.dir, b.decisionID
	if len(o.journal.Decisions) > 0 {
		last := o.journal.Decisions[len(o.journal.Decisions)-1].Request
		o.preview.LastDecision = &last
	}
	if kind == Launch {
		o.preview.CountBasis = "pre-telemetry-launch-total"
	}
	if !isActive {
		if kind == Launch {
			i, err := InspectLaunchMigration(ctx, root, idea)
			if err != nil {
				return o, err
			}
			o.preview.Launch, o.preview.HistorySHA256 = &i, i.HistorySHA256
			if i.LegacyStartFloor > count {
				o.preview.MinimumAccountedActions = i.LegacyStartFloor
			}
		} else {
			i, err := InspectProtocolMigrationDeclarations(ctx, root, idea, kind, b.ideaPath, b.declared, b.unscoped)
			if err != nil {
				return o, err
			}
			o.preview.Protocol, o.preview.HistorySHA256 = &i, i.HistorySHA256
			if i.LowerBound > count {
				o.preview.MinimumAccountedActions = i.LowerBound
			}
		}
	} else if _, err := inspectRecoveredBinding(ctx, root, idea, kind); err != nil {
		return o, err
	}
	after, err := inspectImportFiles(b.dir)
	if err != nil {
		return o, err
	}
	if !reflect.DeepEqual(after, o.files) {
		return o, errors.New("import state changed during recovery preview; inspect again")
	}
	return o, nil
}

func InspectMigrationRecovery(ctx context.Context, root, idea string, kind Kind) (MigrationRecoveryPreview, error) {
	o, err := observeMigrationRecovery(ctx, root, idea, kind)
	return o.preview, err
}

func RecoverMigration(ctx context.Context, root, idea string, kind Kind, r MigrationRecoveryRequest) (PolicyStatus, error) {
	return recoverMigration(ctx, root, idea, kind, r, writeSynced)
}

func inspectRecoveredBinding(ctx context.Context, root, idea string, kind Kind) (PolicyStatus, error) {
	if kind != Launch {
		return inspectProtocolMigrationBinding(ctx, root, idea, kind)
	}
	b, err := LoadLaunchBinding(ctx, root, idea)
	if err != nil {
		return PolicyStatus{}, err
	}
	if b == nil {
		return PolicyStatus{}, errors.New("recovered launch binding is missing")
	}
	return b.Inspect(ctx)
}

func recoverMigration(ctx context.Context, root, idea string, kind Kind, r MigrationRecoveryRequest, persist func(string, []byte) error) (PolicyStatus, error) {
	r.StartedAt = r.StartedAt.UTC()
	if !r.WritersStopped || !validCycleDecision(r.DecisionID, r.Reason, r.ExpectedHistorySHA256) || !validCycleDecision("import", "import", r.ExpectedImportSHA256) {
		return PolicyStatus{}, errors.New("recovery requires exact import/history hashes and a stopped-writer operator decision")
	}
	base, err := loadRecoveryBase(ctx, root, idea, kind)
	if err != nil {
		return PolicyStatus{}, err
	}
	// Serialize before observing mutable publication files. A competing exact
	// recovery can be between journal and ledger writes while we wait.
	release, err := AcquireResourceGuard(ctx, base.dir)
	if err != nil {
		return PolicyStatus{}, err
	}
	defer release()
	o, err := observeMigrationRecovery(ctx, root, idea, kind)
	if err != nil {
		return PolicyStatus{}, err
	}
	b, j := o.base, o.journal
	replay := false
	for n, d := range j.Decisions {
		if d.Request.DecisionID != r.DecisionID {
			continue
		}
		if n != len(j.Decisions)-1 || !reflect.DeepEqual(d.Request, r) {
			return PolicyStatus{}, errors.New("conflicting or superseded recovery decision")
		}
		replay = true
	}
	if o.preview.Active {
		if !replay {
			return PolicyStatus{}, errors.New("active migration cannot be replaced by recovery")
		}
		return inspectRecoveredBinding(ctx, root, idea, kind)
	}
	if o.preview.HistorySHA256 != r.ExpectedHistorySHA256 {
		return PolicyStatus{}, errors.New("history changed; inspect and issue a new recovery decision")
	}
	if !replay {
		if r.ExpectedImportSHA256 != o.preview.ImportSHA256 || r.DecisionID == b.decisionID || len(j.Decisions) >= maxMigrationRecoveries {
			return PolicyStatus{}, errors.New("stale import hash, reused decision or recovery journal limit")
		}
		if err := checkRecoveryObservations(b, j, o.preview.Launch, o.preview.Protocol); err != nil {
			return PolicyStatus{}, err
		}
		initial, err := recoveredInitial(b, o.preview.Retained, o.preview.PriorAccountedActions, r, o.preview.Launch, o.preview.Protocol)
		if err != nil {
			return PolicyStatus{}, err
		}
		d := migrationRecoveryDecision{RecordedAt: time.Now().UTC(), PreviousSHA256: recoveryPrefixDigest(j), Request: r, ImportFiles: o.files, Launch: o.preview.Launch, Protocol: o.preview.Protocol, Initial: initial}
		j.Decisions = append(j.Decisions, d)
		raw, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			return PolicyStatus{}, err
		}
		if len(raw) > 16<<20 {
			return PolicyStatus{}, errors.New("recovery journal exceeds size bound")
		}
		if err := persist(filepath.Join(b.dir, migrationRecoveryFile), append(raw, '\n')); err != nil {
			return PolicyStatus{}, err
		}
	}
	// Re-read persisted journal rather than relying on a successful writer return.
	persisted, err := readRecoveryJournal(b)
	if err != nil || !reflect.DeepEqual(persisted, j) {
		return PolicyStatus{}, errors.New("recovery journal was not retained exactly")
	}
	initial := j.Decisions[len(j.Decisions)-1].Initial
	ledger := Store{Dir: filepath.Join(b.dir, "ledger"), Scope: b.scope, persist: persist}
	_, err = ledger.update(ctx, func(s *Snapshot, _ time.Time) error {
		missing := len(s.Entries) == 0 && s.StartedAt.After(initial.StartedAt)
		// Only an actually absent ledger can use Store.update's new snapshot.
		if _, statErr := os.Lstat(filepath.Join(ledger.Dir, "ledger.json")); !os.IsNotExist(statErr) {
			missing = false
		}
		if !missing && !recoveryLedgerMatches(b, j, *s) {
			return errors.New("concurrent or unknown charges prevent recovery")
		}
		*s = initial
		return nil
	})
	if err != nil {
		return PolicyStatus{}, err
	}
	if err := policyBeforeRecovery(b, false); err != nil {
		return PolicyStatus{}, err
	}
	policyPath := filepath.Join(b.dir, "policy.json")
	if _, err := os.Lstat(policyPath); os.IsNotExist(err) {
		raw, err := json.MarshalIndent(b.policy, "", "  ")
		if err != nil {
			return PolicyStatus{}, err
		}
		if err := persist(policyPath, append(raw, '\n')); err != nil {
			return PolicyStatus{}, err
		}
	} else if err != nil {
		return PolicyStatus{}, err
	}
	latest, err := observeMigrationRecovery(ctx, root, idea, kind)
	if err != nil {
		return PolicyStatus{}, err
	}
	if latest.preview.HistorySHA256 != r.ExpectedHistorySHA256 {
		return PolicyStatus{}, errors.New("history changed during recovery; retained import stays inactive")
	}
	if !reflect.DeepEqual(latest.journal, j) {
		return PolicyStatus{}, errors.New("recovery journal changed before activation")
	}
	_, marker, err := recoveredMigrationState(b)
	if err != nil {
		return PolicyStatus{}, err
	}
	if err := ctx.Err(); err != nil {
		return PolicyStatus{}, err
	}
	if err := persist(filepath.Join(b.dir, "migration-active"), marker); err != nil {
		return PolicyStatus{}, err
	}
	return inspectRecoveredBinding(ctx, root, idea, kind)
}
