package budget

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Every witness here uses a real repository with a real registered worktree
// whose directory is actually removed. No porcelain bytes are fabricated and no
// availability class is asserted on the operator's word alone.
//
// TotalActions below is taken from the observed floor, so these fixtures invent
// no count for the unknown region and choose no operator reconciliation.
func declaredFixtureRequest(t *testing.T, root string, kind Kind, declared []string) ProtocolMigrationRequest {
	t.Helper()
	i, err := InspectProtocolMigrationDeclared(context.Background(), root, "idea", kind, "", declared)
	if err != nil {
		t.Fatalf("declared inspect: %v", err)
	}
	return ProtocolMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "declared-fixture",
		Reason: "Explicit fixture-only declared-unavailable protocol accounting", StartedAt: time.Now().UTC().Add(-time.Hour),
		TotalActions: i.LowerBound, Maximum: 3, WritersStopped: true, DeclaredUnavailable: declared}
}

func assertOutsideDeclared(t *testing.T, i ProtocolMigrationInventory, path string) {
	t.Helper()
	inside := func(candidate string) bool {
		return candidate == path || strings.HasPrefix(candidate, path+string(filepath.Separator))
	}
	for _, origin := range i.History.Roots {
		if inside(origin) {
			t.Fatalf("declared-unavailable worktree %q supplied history root %q", path, origin)
		}
	}
	for _, source := range i.History.Sources {
		if inside(source.Root) {
			t.Fatalf("declared-unavailable worktree %q supplied history source %q", path, source.Path)
		}
	}
}

// The load-bearing witness: a declared registration is retained as explicitly
// unknown history, contributes no evidence and no count, and the undeclared
// refusal that shares the same Git output is untouched.
func TestDeclaredUnavailableWorktreeAdmitsUnknownHistoryOnly(t *testing.T) {
	ctx := context.Background()
	root, linked, _ := worktreeRepoFixture(t)
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	if _, err := InspectProtocolMigration(ctx, root, "idea", CrossReview, ""); err == nil || !strings.Contains(err.Error(), "historical worktree is unavailable") {
		t.Fatalf("undeclared missing registration was not refused: %v", err)
	}
	i, err := InspectProtocolMigrationDeclared(ctx, root, "idea", CrossReview, "", []string{linked})
	if err != nil {
		t.Fatalf("declared inspect refused a real missing registration: %v", err)
	}
	if len(i.History.UnavailableRoots) != 1 {
		t.Fatalf("declared registration was dropped or duplicated: %+v", i.History.UnavailableRoots)
	}
	row := i.History.UnavailableRoots[0]
	if row.Path != linked || row.Observation != UnavailableMissing || row.History != UnknownHistory {
		t.Fatalf("unknown-history row lost its verbatim path, stat class or unknown history: %+v", row)
	}
	if i.History.HistoryCoverage != DeclaredIncomplete || i.LowerBoundBasis != SurvivingVisibleFloor {
		t.Fatalf("declared inventory does not record incomplete coverage: %q %q", i.History.HistoryCoverage, i.LowerBoundBasis)
	}
	assertOutsideDeclared(t, i, linked)
	// The unknown region became neither a count nor a zero-filled floor: the
	// floor is exactly the surviving visible evidence, which here is none.
	if i.LowerBound != 0 || len(i.Evidence) != 1 || i.Evidence[0].Root != root {
		t.Fatalf("unknown history contributed evidence or a count: %d %+v", i.LowerBound, i.Evidence)
	}
	// An undeclared inventory keeps its exact previous shape, so an already
	// applied import still reads and digests identically in this binary.
	other := t.TempDir()
	j, err := InspectProtocolMigration(ctx, other, "idea", CrossReview, "")
	if err != nil {
		t.Fatal(err)
	}
	if j.LowerBoundBasis != "" || j.History.HistoryCoverage != "" || j.History.UnavailableRoots != nil {
		t.Fatalf("undeclared inventory claims declared fields: %+v", j)
	}
}

// Declarations bind the actual stat class and the retained registration. An
// available, unregistered, non-absolute or otherwise unreadable path is refused
// rather than admitted, and ordinary bootstrap never sees a declaration at all.
func TestDeclaredUnavailableRefusesAvailableUnregisteredAndRelativePaths(t *testing.T) {
	ctx := context.Background()
	root, linked, _ := worktreeRepoFixture(t)
	if _, err := InspectProtocolMigrationDeclared(ctx, root, "idea", CrossReview, "", []string{linked}); err == nil || !strings.Contains(err.Error(), "available again") {
		t.Fatalf("a present registration was declared unknown: %v", err)
	}
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	for _, declared := range [][]string{{linked + "-unregistered"}, {"relative/path"}, {""}} {
		if _, err := InspectProtocolMigrationDeclared(ctx, root, "idea", CrossReview, "", declared); err == nil {
			t.Fatalf("unregistered or non-verbatim declaration accepted: %v", declared)
		}
	}
	// A non-directory registration is a separate, separately bound stat class.
	if err := os.WriteFile(linked, []byte("not a worktree"), 0o600); err != nil {
		t.Fatal(err)
	}
	i, err := InspectProtocolMigrationDeclared(ctx, root, "idea", CrossReview, "", []string{linked})
	if err != nil {
		t.Fatalf("declared non-directory registration refused: %v", err)
	}
	if len(i.History.UnavailableRoots) != 1 || i.History.UnavailableRoots[0].Observation != UnavailableNotDirectory {
		t.Fatalf("stat class was not bound: %+v", i.History.UnavailableRoots)
	}
	// Ordinary bootstrap passes no declaration and stays fail-closed.
	if _, err := EnsureCycleBinding(ctx, root, "idea", CrossReview, 3, 0, "", ""); err == nil || !strings.Contains(err.Error(), "historical worktree is unavailable") {
		t.Fatalf("ordinary cycle bootstrap consulted a declaration: %v", err)
	}
	if _, err := ConfigureLaunchBudget(ctx, root, "idea", LaunchPolicy{MaxLaunches: 1}); err == nil || !strings.Contains(err.Error(), "historical worktree is unavailable") {
		t.Fatalf("ordinary launch configuration consulted a declaration: %v", err)
	}
}

// One declaration is one decision: ordering, repetition and the nil/empty
// spelling collapse into the same exact replay, and a different declared set is
// a conflicting decision rather than a silent second import.
func TestDeclaredMigrationReplayIsCanonicalAndConflictsOnChange(t *testing.T) {
	ctx := context.Background()
	root, linked, git := worktreeRepoFixture(t)
	second := filepath.Join(filepath.Dir(linked), "second")
	git(root, "worktree", "add", "-q", "-b", "second", second)
	for _, path := range []string{linked, second} {
		if err := os.RemoveAll(path); err != nil {
			t.Fatal(err)
		}
	}
	r := declaredFixtureRequest(t, root, CrossReview, []string{second, linked})
	s, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, r)
	if err != nil {
		t.Fatalf("declared import: %v", err)
	}
	repeat := r
	repeat.DeclaredUnavailable = []string{linked, second, linked}
	replayed, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, repeat)
	if err != nil || replayed.Spent != s.Spent || !replayed.StartedAt.Equal(s.StartedAt) {
		t.Fatalf("reordered duplicate spelling was not the same decision: %+v %v", replayed, err)
	}
	conflict := r
	conflict.DeclaredUnavailable = []string{linked}
	if _, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, conflict); err == nil {
		t.Fatal("a narrower declaration replayed as the same decision")
	}
	empty := r
	empty.DeclaredUnavailable = []string{}
	if _, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, empty); err == nil {
		t.Fatal("an empty declaration replayed over a declared import")
	}
	dir, _, err := protocolMigrationScope(ctx, root, "idea", CrossReview)
	if err != nil {
		t.Fatal(err)
	}
	record, err := readProtocolMigration(dir)
	if err != nil {
		t.Fatalf("retained declared record: %v", err)
	}
	canonical, err := normalizeDeclaredUnavailable([]string{second, linked})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(record.Request.DeclaredUnavailable, canonical) {
		t.Fatalf("record lost the canonical declaration: %+v", record.Request.DeclaredUnavailable)
	}
	if len(record.Inventory.History.UnavailableRoots) != 2 {
		t.Fatalf("record lost its unknown-history rows: %+v", record.Inventory.History.UnavailableRoots)
	}
	for _, row := range record.Inventory.History.UnavailableRoots {
		if row.History != UnknownHistory {
			t.Fatalf("retained row carries a history value other than unknown: %+v", row)
		}
	}
}

// A publication fault before the activation witness leaves the import inactive
// and ungrantable; the recovery preview re-inspects the same declared history,
// and the exact replay completes the activation without a second scan.
func TestDeclaredMigrationRecoveryPreviewAndReplayStayRequestScoped(t *testing.T) {
	ctx := context.Background()
	root, linked, _ := worktreeRepoFixture(t)
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	r := declaredFixtureRequest(t, root, CrossReview, []string{linked})
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
	if preview.Protocol == nil || len(preview.Protocol.History.UnavailableRoots) != 1 || preview.Protocol.LowerBoundBasis != SurvivingVisibleFloor {
		t.Fatalf("recovery preview dropped the unknown-history rows: %+v", preview.Protocol)
	}
	s, err := MigrateProtocolBudget(ctx, root, "idea", CrossReview, r)
	if err != nil || s.Spent != r.TotalActions {
		t.Fatalf("exact replay did not complete the activation: %+v %v", s, err)
	}
	b, err := LoadCycleBinding(ctx, root, "idea", CrossReview)
	if err != nil || b == nil {
		t.Fatalf("activated declared import is not loadable: %v", err)
	}
	if b.Policy.Maximum != r.Maximum || b.Policy.Carried != 0 {
		t.Fatalf("declared import changed the frozen cap: %+v", b.Policy)
	}
	// The declaration was an argument to one decision, not repository state: a
	// different idea on the same broken repository is still refused.
	if _, err := EnsureCycleBinding(ctx, root, "other-idea", CrossReview, 3, 0, "", ""); err == nil || !strings.Contains(err.Error(), "historical worktree is unavailable") {
		t.Fatalf("a declaration leaked into another idea's bootstrap: %v", err)
	}
}
