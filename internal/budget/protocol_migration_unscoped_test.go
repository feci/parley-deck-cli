package budget

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// The fixture is the shape actually observed in this repository's own history:
// a single run.created line whose data carries no idea and no idea_slug key.
// Nothing here fabricates an identity, a cursor or a count.
const unscopedRunEvent = `{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","data":{"mode":"auto","task":"smoke implementation run"}}`

func writeRunDirectory(t *testing.T, root, name string, lines ...string) {
	t.Helper()
	dir := filepath.Join(root, "parley-deck", "runs", name)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func declareRun(t *testing.T, root, name string) string {
	t.Helper()
	digest, err := RunDirectoryManifestDigest(filepath.Join(root, "parley-deck", "runs", name))
	if err != nil {
		t.Fatalf("manifest digest: %v", err)
	}
	return unscopedRunPrefix + name + "=" + digest
}

func unscopedFixtureRequest(t *testing.T, root string, declaredRuns []string) ProtocolMigrationRequest {
	t.Helper()
	i, err := InspectProtocolMigrationDeclarations(context.Background(), root, "idea", CrossReview, "", nil, declaredRuns)
	if err != nil {
		t.Fatalf("declared inspect: %v", err)
	}
	// TotalActions is the observed floor: this fixture invents no count for the
	// excluded region and chooses no operator reconciliation.
	return ProtocolMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "unscoped-fixture",
		Reason: "Explicit fixture-only declared-unscoped protocol accounting", StartedAt: time.Now().UTC().Add(-time.Hour),
		TotalActions: i.LowerBound, Maximum: 3, WritersStopped: true, DeclaredUnscopedRuns: declaredRuns}
}

// The load-bearing witness: an undeclared unidentifiable run keeps its exact
// refusal, a declared one keeps its readable bytes as sources, and it produces
// no evidence row and no count — in particular no floor of zero.
func TestDeclaredUnscopedRunRetainsBytesAndCountsNothing(t *testing.T) {
	ctx := context.Background()
	root, _, _ := worktreeRepoFixture(t)
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent)
	if _, err := InspectProtocolMigration(ctx, root, "idea", CrossReview, ""); err == nil || !strings.Contains(err.Error(), "historical run identity is missing or conflicting") {
		t.Fatalf("undeclared unidentifiable run was not refused verbatim: %v", err)
	}
	declaration := declareRun(t, root, "20260510T194003Z")
	i, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration})
	if err != nil {
		t.Fatalf("declared inspect refused a real unidentifiable run: %v", err)
	}
	if len(i.History.UnscopedRuns) != 1 {
		t.Fatalf("declared run was dropped or duplicated: %+v", i.History.UnscopedRuns)
	}
	row := i.History.UnscopedRuns[0]
	if row.Path != unscopedRunPrefix+"20260510T194003Z" || row.History != UnknownHistory || row.Copies != 1 || len(row.Roots) != 1 || row.Roots[0] != root {
		t.Fatalf("unknown-identity row lost its path, unknown history or actual root: %+v", row)
	}
	if row.Path+"="+row.SHA256 != declaration {
		t.Fatalf("row is not bound to the declared manifest digest: %+v", row)
	}
	if i.History.HistoryCoverage != DeclaredIncomplete || i.LowerBoundBasis != ScopedRunsFloor {
		t.Fatalf("declared inventory does not record its incomplete coverage: %q %q", i.History.HistoryCoverage, i.LowerBoundBasis)
	}
	// The readable bytes stay evidence of themselves: the exclusion is from
	// scope-counted evidence only, never structural.
	if !retainedMigrationSource(i.History.Sources, root, row.Path+"/events.jsonl") {
		t.Fatalf("declared run lost its retained readable source: %+v", i.History.Sources)
	}
	// Neither counted nor zeroed: no evidence row mentions the run at all.
	for _, e := range i.Evidence {
		if strings.Contains(e.Path, "20260510T194003Z") {
			t.Fatalf("declared run supplied count evidence: %+v", e)
		}
	}
	if i.LowerBound != 0 || i.Earliest != nil {
		t.Fatalf("declared run supplied a floor or this idea's epoch: %d %v", i.LowerBound, i.Earliest)
	}
	if row.Earliest == nil {
		t.Fatal("the run's own earliest event time was not recorded on its row")
	}
	// An undeclared inventory keeps its exact previous shape and digest.
	j, err := InspectProtocolMigration(ctx, t.TempDir(), "idea", CrossReview, "")
	if err != nil {
		t.Fatal(err)
	}
	if j.LowerBoundBasis != "" || j.History.HistoryCoverage != "" || j.History.UnscopedRuns != nil {
		t.Fatalf("undeclared inventory claims declared fields: %+v", j)
	}
	raw, err := json.Marshal(j.History)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "unscoped_runs") {
		t.Fatalf("undeclared inventory marshals the new key: %s", raw)
	}
}

// One declaration covers every identical copy of one file set. A copy that
// gained, lost or changed a byte, and a declaration matching no visible copy,
// are refusals rather than a narrowed or widened admission.
func TestDeclaredUnscopedRunBindsEveryVisibleCopy(t *testing.T) {
	ctx := context.Background()
	root, linked, _ := worktreeRepoFixture(t)
	for _, origin := range []string{root, linked} {
		writeRunDirectory(t, origin, "20260510T194003Z", unscopedRunEvent)
	}
	declaration := declareRun(t, root, "20260510T194003Z")
	i, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration})
	if err != nil {
		t.Fatalf("identical copies were not grouped: %v", err)
	}
	row := i.History.UnscopedRuns[0]
	if row.Copies != 2 || !reflect.DeepEqual(row.Roots, []string{linked, root}) {
		t.Fatalf("copies of one file set were not bound to their actual roots: %+v", row)
	}
	// Copies counts copies of one file set. It is never an action count, so the
	// second copy changed no floor.
	if i.LowerBound != 0 {
		t.Fatalf("a second copy was read as more actions: %d", i.LowerBound)
	}
	mutated := filepath.Join(linked, "parley-deck", "runs", "20260510T194003Z", "events.jsonl")
	if err := os.WriteFile(mutated, []byte(unscopedRunEvent+" \n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration}); err == nil || !strings.Contains(err.Error(), "differs at") {
		t.Fatalf("a changed copy was admitted under the declaration: %v", err)
	}
	writeRunDirectory(t, linked, "20260510T194003Z", unscopedRunEvent)
	// The manifest is the whole directory, so a file the scanner never reads
	// still invalidates the declaration.
	added := filepath.Join(linked, "parley-deck", "runs", "20260510T194003Z", "agents")
	if err := os.MkdirAll(added, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(added, "agent.log"), []byte("unrelated\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration}); err == nil || !strings.Contains(err.Error(), "differs at") {
		t.Fatalf("an added file was admitted under the declaration: %v", err)
	}
	for _, origin := range []string{root, linked} {
		if err := os.RemoveAll(filepath.Join(origin, "parley-deck", "runs", "20260510T194003Z")); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration}); err == nil || !strings.Contains(err.Error(), "not present in any visible root") {
		t.Fatalf("a stale declaration was admitted: %v", err)
	}
}

// Only an actually absent identity is declarable. A recovered identity and a
// conflicting one are evidence of identity, so both refuse: a declaration
// admits an unknown, it never overrides or discards evidence.
func TestDeclaredUnscopedRunPermitsAbsentIdentityOnly(t *testing.T) {
	ctx := context.Background()
	root, _, _ := worktreeRepoFixture(t)
	recoverable := `{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","data":{"idea":"idea"}}`
	writeRunDirectory(t, root, "recoverable", recoverable)
	declaration := declareRun(t, root, "recoverable")
	if _, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration}); err == nil || !strings.Contains(err.Error(), "recoverable idea identity") {
		t.Fatalf("a run that proves its own identity was declared away: %v", err)
	}
	if err := os.RemoveAll(filepath.Join(root, "parley-deck", "runs", "recoverable")); err != nil {
		t.Fatal(err)
	}
	conflicting := `{"time":"2026-05-10T19:40:04.126637Z","type":"run.phase","data":{"idea":"other"}}`
	writeRunDirectory(t, root, "conflicting", recoverable, conflicting)
	conflictDeclaration := declareRun(t, root, "conflicting")
	if _, err := InspectProtocolMigration(ctx, root, "idea", CrossReview, ""); err == nil || !strings.Contains(err.Error(), "missing or conflicting") {
		t.Fatalf("a conflicting identity was not refused undeclared: %v", err)
	}
	if _, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{conflictDeclaration}); err == nil || !strings.Contains(err.Error(), "missing or conflicting") {
		t.Fatalf("a conflicting identity was declarable: %v", err)
	}
}

// A declared run contributes no process-start evidence either: its events are
// not harvested, so no legacy start floor is derived from a run whose scope is
// unknown, and no zero is asserted for it.
func TestDeclaredUnscopedRunHarvestsNoProcessStarts(t *testing.T) {
	ctx := context.Background()
	root, _, _ := worktreeRepoFixture(t)
	started := `{"time":"2026-05-10T19:40:05.126637Z","type":"agent.started","data":{"invocation_id":"legacy-1"}}`
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent, started)
	declaration := declareRun(t, root, "20260510T194003Z")
	i, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration})
	if err != nil {
		t.Fatalf("declared inspect: %v", err)
	}
	if i.History.LegacyStartFloor != 0 || i.LowerBound != 0 {
		t.Fatalf("a declared run supplied a start floor or count: %d %d", i.History.LegacyStartFloor, i.LowerBound)
	}
	if len(i.History.UnscopedRuns) != 1 || i.History.UnscopedRuns[0].Copies != 1 {
		t.Fatalf("the declared row was lost: %+v", i.History.UnscopedRuns)
	}
}

// Canonical form is one decision. Ordering, repetition and the nil/empty
// spelling collapse into one exact replay; a different digest is a conflicting
// decision, and every malformed spelling is refused rather than repaired.
func TestDeclaredUnscopedRunNormalizationAndReplay(t *testing.T) {
	ctx := context.Background()
	root, _, _ := worktreeRepoFixture(t)
	writeRunDirectory(t, root, "aaa", unscopedRunEvent)
	writeRunDirectory(t, root, "bbb", unscopedRunEvent)
	first, second := declareRun(t, root, "aaa"), declareRun(t, root, "bbb")
	digest := strings.SplitN(first, "=", 2)[1]
	for _, bad := range [][]string{
		{"/abs/parley-deck/runs/aaa=" + digest},
		{"parley-deck/runs/../runs/aaa=" + digest},
		{"parley-deck/ideas/aaa=" + digest},
		{"parley-deck/runs/aaa/nested=" + digest},
		{"parley-deck/runs/aaa=" + strings.ToUpper(digest)},
		{"parley-deck/runs/aaa=deadbeef"},
		{"parley-deck/runs/aaa"},
		{""},
		{"parley-deck/runs/aaa=" + digest, "parley-deck/runs/aaa=" + strings.Repeat("0", 64)},
	} {
		if _, err := normalizeDeclaredUnscopedRuns(bad); err == nil {
			t.Fatalf("malformed declaration accepted: %v", bad)
		}
	}
	canonical, err := normalizeDeclaredUnscopedRuns([]string{second, first, first})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(canonical, []string{first, second}) {
		t.Fatalf("declaration was not canonicalized: %+v", canonical)
	}
	if empty, err := normalizeDeclaredUnscopedRuns([]string{}); err != nil || empty != nil {
		t.Fatalf("an empty spelling is not the nil decision: %+v %v", empty, err)
	}
	r := unscopedFixtureRequest(t, root, canonical)
	s, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, r)
	if err != nil {
		t.Fatalf("declared import: %v", err)
	}
	repeat := r
	repeat.DeclaredUnscopedRuns = []string{second, first, second}
	replayed, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, repeat)
	if err != nil || replayed.Spent != s.Spent || !replayed.StartedAt.Equal(s.StartedAt) {
		t.Fatalf("reordered duplicate spelling was not the same decision: %+v %v", replayed, err)
	}
	narrower := r
	narrower.DeclaredUnscopedRuns = []string{first}
	if _, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, narrower); err == nil {
		t.Fatal("a narrower declaration replayed as the same decision")
	}
	dir, _, err := protocolMigrationScope(ctx, root, "idea", CrossReview)
	if err != nil {
		t.Fatal(err)
	}
	record, err := readProtocolMigration(dir)
	if err != nil {
		t.Fatalf("retained declared record: %v", err)
	}
	if !reflect.DeepEqual(record.Request.DeclaredUnscopedRuns, canonical) || record.Request.DeclaredUnavailable != nil {
		t.Fatalf("record lost the canonical declaration or gained another: %+v", record.Request)
	}
	// The declaration was an argument to one decision, not repository state:
	// another idea inspected with no declaration still refuses on the same runs,
	// and ordinary step bootstrap — which passes no declaration and consults no
	// persisted exclusion — still refuses with its own verbatim text.
	if _, err := InspectProtocolMigration(ctx, root, "other-idea", CrossReview, ""); err == nil || !strings.Contains(err.Error(), "missing or conflicting") {
		t.Fatalf("a declaration leaked into another idea's inspection: %v", err)
	}
	if _, err := EnsureStepBinding(ctx, root, "idea", 3, 0); err == nil || !strings.Contains(err.Error(), "historical driver idea identity is missing or malformed") {
		t.Fatalf("ordinary step bootstrap consulted a declaration: %v", err)
	}
}

// A publication fault before the activation witness leaves the import inactive;
// the recovery preview re-inspects with the declaration read back from the
// frozen request, and the exact replay completes the activation.
func TestDeclaredUnscopedRunRecoveryRetainsTheDeclaration(t *testing.T) {
	ctx := context.Background()
	root, _, _ := worktreeRepoFixture(t)
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent)
	r := unscopedFixtureRequest(t, root, []string{declareRun(t, root, "20260510T194003Z")})
	persist := func(path string, raw []byte) error {
		if filepath.Base(path) == "migration-active" {
			return errors.New("injected activation fault")
		}
		return writeSynced(path, raw)
	}
	if _, err := migrateProtocolBudget(ctx, root, "idea", CrossReview, r, persist); err == nil {
		t.Fatal("injected activation fault was not reported")
	}
	if _, err := LoadCycleBinding(ctx, root, "idea", CrossReview); err == nil || !strings.Contains(err.Error(), "not durably active") {
		t.Fatalf("an unactivated declared import was usable: %v", err)
	}
	preview, err := InspectMigrationRecovery(ctx, root, "idea", CrossReview)
	if err != nil {
		t.Fatalf("recovery preview lost the retained declaration: %v", err)
	}
	if preview.Active || preview.HistorySHA256 != r.ExpectedHistorySHA256 {
		t.Fatalf("recovery preview state or history changed: %+v", preview)
	}
	if preview.Protocol == nil || len(preview.Protocol.History.UnscopedRuns) != 1 || preview.Protocol.LowerBoundBasis != ScopedRunsFloor {
		t.Fatalf("recovery preview dropped the unknown-identity rows: %+v", preview.Protocol)
	}
	s, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, r)
	if err != nil || s.Spent != r.TotalActions {
		t.Fatalf("exact replay did not complete the activation: %+v %v", s, err)
	}
}

// Both declaration kinds compose: each keeps its own retention rule, and the
// inventory records the one combined basis its rows require rather than
// claiming either kind alone.
func TestBothDeclarationsComposeIntoOneCombinedBasis(t *testing.T) {
	ctx := context.Background()
	root, linked, _ := worktreeRepoFixture(t)
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent)
	declaration := declareRun(t, root, "20260510T194003Z")
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	i, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", []string{linked}, []string{declaration})
	if err != nil {
		t.Fatalf("composed declarations refused: %v", err)
	}
	if len(i.History.UnavailableRoots) != 1 || len(i.History.UnscopedRuns) != 1 {
		t.Fatalf("a declaration kind was dropped: %+v", i.History)
	}
	if i.LowerBoundBasis != SurvivingScopedFloor {
		t.Fatalf("composed coverage claimed one kind alone: %q", i.LowerBoundBasis)
	}
	assertOutsideDeclared(t, i, linked)
	if !retainedMigrationSource(i.History.Sources, root, unscopedRunPrefix+"20260510T194003Z/events.jsonl") {
		t.Fatalf("the readable declared run was excluded structurally: %+v", i.History.Sources)
	}
	// An inventory whose recorded basis does not match its rows is refused, in
	// both directions and for both kinds.
	for _, broken := range []ProtocolMigrationInventory{
		func() ProtocolMigrationInventory { j := i; j.LowerBoundBasis = SurvivingVisibleFloor; return j }(),
		func() ProtocolMigrationInventory { j := i; j.LowerBoundBasis = ScopedRunsFloor; return j }(),
		func() ProtocolMigrationInventory { j := i; j.LowerBoundBasis = ""; return j }(),
	} {
		if err := checkDeclarations(broken, []string{linked}, []string{declaration}); err == nil {
			t.Fatalf("a mismatched coverage basis was accepted: %q", broken.LowerBoundBasis)
		}
	}
	if err := checkDeclarations(i, []string{linked}, nil); err == nil {
		t.Fatal("retained unscoped rows were accepted without their declaration")
	}
	if err := checkDeclarations(i, nil, []string{declaration}); err == nil {
		t.Fatal("retained unavailable rows were accepted without their declaration")
	}
}

// The declaration exempts exactly one case: an idea key that is absent. Bytes
// that cannot be parsed, an event that does not say what it is, and a directory
// this walk cannot hash are all refusals still, because none of them is an
// established absence of identity — each is something unread.
func TestDeclaredUnscopedRunExemptsOnlyTheAbsentKey(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct{ name, line, refusal string }{
		// The decoder's own text is not asserted for the malformed line: only
		// that unparseable bytes are never read as an established absence.
		{"malformed", `{"time":"2026-05-10T19:40:03.126637Z","type":"run.created"`, ""},
		{"untyped", `{"time":"2026-05-10T19:40:03.126637Z","data":{"mode":"auto"}}`, "missing required historical field type"},
	} {
		root, _, _ := worktreeRepoFixture(t)
		writeRunDirectory(t, root, c.name, c.line)
		declaration := declareRun(t, root, c.name)
		_, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration})
		if err == nil || c.refusal != "" && !strings.Contains(err.Error(), c.refusal) {
			t.Fatalf("declared %s bytes were admitted as an absent identity: %v", c.name, err)
		}
	}
	// A declaration binds an exact recursive file set, so a directory holding an
	// entry the walk refuses to hash cannot be bound at all — the digest check
	// refuses before any of its bytes are read as evidence.
	root, _, _ := worktreeRepoFixture(t)
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent)
	declaration := declareRun(t, root, "20260510T194003Z")
	if err := os.Symlink(filepath.Join(root, "parley-deck"), filepath.Join(root, "parley-deck", "runs", "20260510T194003Z", "link")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration}); err == nil || !strings.Contains(err.Error(), "contains a symlink") {
		t.Fatalf("a declared run directory with an unhashable entry was bound: %v", err)
	}
}
