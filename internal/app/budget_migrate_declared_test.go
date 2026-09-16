package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
)

// A real repository with a real registered worktree whose directory is then
// removed. The CLI surface is exercised against actual Git registration state.
func declaredWorktreeCLIFixture(t *testing.T) (root, linked string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git unavailable: %v", err)
	}
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git := func(dir string, args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	root = filepath.Join(base, "main")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	git(root, "init", "-q")
	git(root, "-c", "user.name=fixture", "-c", "user.email=fixture@example.invalid", "commit", "--allow-empty", "-qm", "fixture")
	linked = filepath.Join(base, "linked")
	git(root, "worktree", "add", "-q", "-b", "linked", linked)
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	return root, linked
}

// The declaration is a protocol inspect/apply control only: it is repeatable,
// it is rejected for launch accounting, for the recovery verb and for an
// unknown verb, and no refused or read-only command creates accounting state.
func TestBudgetMigrateDeclaredWorktreeIsProtocolInspectApplyOnly(t *testing.T) {
	ctx := context.Background()
	root, linked := declaredWorktreeCLIFixture(t)
	protocol := []string{"migrate", "inspect", "--dir", root, "--idea", "idea", "--kind", "cross-review"}
	if code := runBudgetPlatformControl(ctx, protocol, io.Discard, io.Discard, false, false); code != 1 {
		t.Fatalf("undeclared inspect did not refuse the unavailable registration: %d", code)
	}
	var out, errout bytes.Buffer
	declared := append(append([]string{}, protocol...), "--declare-unavailable-worktree", linked)
	if code := runBudgetPlatformControl(ctx, declared, &out, &errout, false, false); code != 0 {
		t.Fatalf("declared inspect: %d %s", code, errout.String())
	}
	var i budget.ProtocolMigrationInventory
	if err := json.Unmarshal(out.Bytes(), &i); err != nil {
		t.Fatal(err)
	}
	if len(i.History.UnavailableRoots) != 1 || i.History.UnavailableRoots[0].Path != linked {
		t.Fatalf("inspect output lost the declared row: %+v", i.History.UnavailableRoots)
	}
	if i.History.UnavailableRoots[0].History != budget.UnknownHistory || i.History.HistoryCoverage != budget.DeclaredIncomplete || i.LowerBoundBasis != budget.SurvivingVisibleFloor {
		t.Fatalf("inspect output does not report unknown, incomplete coverage: %+v", i)
	}
	refused := [][]string{
		{"migrate", "inspect", "--dir", root, "--idea", "idea", "--kind", "launch", "--declare-unavailable-worktree", linked},
		{"migrate", "apply", "--dir", root, "--idea", "idea", "--kind", "launch", "--declare-unavailable-worktree", linked},
		{"migrate", "recover", "inspect", "--dir", root, "--idea", "idea", "--kind", "cross-review", "--declare-unavailable-worktree", linked},
		{"migrate", "recover", "apply", "--dir", root, "--idea", "idea", "--kind", "cross-review", "--declare-unavailable-worktree", linked},
		{"migrate", "declare", "--dir", root, "--idea", "idea", "--kind", "cross-review", "--declare-unavailable-worktree", linked},
		{"worktree", "inspect", "--dir", root, "--declare-unavailable-worktree", linked},
	}
	for _, args := range refused {
		if code := runBudgetPlatformControl(ctx, args, io.Discard, io.Discard, true, true); code != 2 {
			t.Fatalf("declaration accepted outside protocol inspect/apply: %v %d", args, code)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "parley-launch-budgets")); !os.IsNotExist(err) {
		t.Fatalf("a read-only or refused command created accounting state: %v", err)
	}
}

// An attended apply carries the repeatable declaration into the immutable
// decision, replays exactly, and refuses a changed declaration.
func TestBudgetMigrateDeclaredWorktreeApplyReplaysExactly(t *testing.T) {
	ctx := context.Background()
	root, linked := declaredWorktreeCLIFixture(t)
	var out, errout bytes.Buffer
	inspect := []string{"migrate", "inspect", "--dir", root, "--idea", "idea", "--kind", "cross-review", "--declare-unavailable-worktree", linked}
	if code := runBudgetPlatformControl(ctx, inspect, &out, &errout, false, false); code != 0 {
		t.Fatalf("declared inspect: %d %s", code, errout.String())
	}
	var i budget.ProtocolMigrationInventory
	if err := json.Unmarshal(out.Bytes(), &i); err != nil {
		t.Fatal(err)
	}
	// The count is the observed floor over the surviving visible worktrees. No
	// operator reconciliation is chosen here and no ceiling is raised.
	apply := []string{"migrate", "apply", "--dir", root, "--idea", "idea", "--kind", "cross-review",
		"--declare-unavailable-worktree", linked, "--expected-history-sha256", i.HistorySHA256,
		"--decision-id", "declared-cli-fixture", "--reason", "Explicit fixture-only declared-unavailable protocol accounting",
		"--started-at", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano),
		"--total-actions", strconv.Itoa(i.LowerBound), "--max-cycles", "3", "--writers-stopped", "--yes"}
	if code := runBudgetPlatformControl(ctx, apply, io.Discard, &errout, true, true); code != 0 {
		t.Fatalf("declared apply: %s", errout.String())
	}
	if code := runBudgetPlatformControl(ctx, apply, io.Discard, &errout, true, true); code != 0 {
		t.Fatalf("exact declared replay: %s", errout.String())
	}
	b, err := budget.LoadCycleBinding(ctx, root, "idea", budget.CrossReview)
	if err != nil || b == nil {
		t.Fatalf("declared import is not loadable: %v", err)
	}
	if b.Policy.Maximum != 3 || b.Policy.Carried != 0 {
		t.Fatalf("declared import changed the frozen cap: %+v", b.Policy)
	}
	var undeclared []string
	for n := 0; n < len(apply); n++ {
		if apply[n] == "--declare-unavailable-worktree" {
			n++
			continue
		}
		undeclared = append(undeclared, apply[n])
	}
	if code := runBudgetPlatformControl(ctx, undeclared, io.Discard, io.Discard, true, true); code == 0 {
		t.Fatal("an apply without the declaration replayed over a declared import")
	}
}

// The shape actually observed in this repository's own history: one run.created
// line whose data carries no idea and no idea_slug key.
const unscopedRunCLIEvent = `{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","data":{"mode":"auto","task":"smoke implementation run"}}`

// Builds the run directory inside the throwaway fixture repository and returns
// the exact flag value for it, computed from the actual file set by the same
// exported helper the accounting boundary uses. No digest is fabricated here,
// and nothing in this file touches history outside the temporary repository.
func declaredUnscopedRunCLIFixture(t *testing.T, root, name string) string {
	t.Helper()
	dir := filepath.Join(root, "parley-deck", "runs", name)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), []byte(unscopedRunCLIEvent+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	digest, err := budget.RunDirectoryManifestDigest(dir)
	if err != nil {
		t.Fatalf("manifest digest: %v", err)
	}
	return "parley-deck/runs/" + name + "=" + digest
}

// The unscoped-run declaration is a protocol inspect/apply control only: it is
// repeatable, it composes with the worktree declaration, it is rejected for
// launch accounting, for the recovery and declare verbs, for an unknown verb
// and for the worktree surface, and no refused or read-only command creates
// accounting state.
func TestBudgetMigrateDeclaredUnscopedRunIsProtocolInspectApplyOnly(t *testing.T) {
	ctx := context.Background()
	root, linked := declaredWorktreeCLIFixture(t)
	run := declaredUnscopedRunCLIFixture(t, root, "20260510T194003Z")
	protocol := []string{"migrate", "inspect", "--dir", root, "--idea", "idea", "--kind", "cross-review"}
	worktreeOnly := append(append([]string{}, protocol...), "--declare-unavailable-worktree", linked)
	// The run declaration is load-bearing: declaring only the worktree still
	// refuses on the unidentifiable run.
	if code := runBudgetPlatformControl(ctx, worktreeOnly, io.Discard, io.Discard, false, false); code != 1 {
		t.Fatalf("an undeclared unidentifiable run was not refused: %d", code)
	}
	var out, errout bytes.Buffer
	declared := append(append([]string{}, worktreeOnly...), "--declare-unscoped-run", run)
	if code := runBudgetPlatformControl(ctx, declared, &out, &errout, false, false); code != 0 {
		t.Fatalf("declared inspect: %d %s", code, errout.String())
	}
	var i budget.ProtocolMigrationInventory
	if err := json.Unmarshal(out.Bytes(), &i); err != nil {
		t.Fatal(err)
	}
	if len(i.History.UnscopedRuns) != 1 {
		t.Fatalf("inspect output lost the declared run row: %+v", i.History.UnscopedRuns)
	}
	row := i.History.UnscopedRuns[0]
	if row.Path+"="+row.SHA256 != run || row.History != budget.UnknownHistory || row.Copies != 1 {
		t.Fatalf("the row is not the declared file set reported as unknown: %+v", row)
	}
	// Both declarations compose into the one combined basis their rows require,
	// and no field of the output is a total.
	if i.History.HistoryCoverage != budget.DeclaredIncomplete || i.LowerBoundBasis != budget.SurvivingScopedFloor {
		t.Fatalf("composed coverage claimed one declaration kind alone: %q %q", i.History.HistoryCoverage, i.LowerBoundBasis)
	}
	if len(i.History.UnavailableRoots) != 1 || i.History.UnavailableRoots[0].Path != linked {
		t.Fatalf("the worktree declaration was dropped by the second flag: %+v", i.History.UnavailableRoots)
	}
	refused := [][]string{
		{"migrate", "inspect", "--dir", root, "--idea", "idea", "--kind", "launch", "--declare-unscoped-run", run},
		{"migrate", "apply", "--dir", root, "--idea", "idea", "--kind", "launch", "--declare-unscoped-run", run},
		{"migrate", "recover", "inspect", "--dir", root, "--idea", "idea", "--kind", "cross-review", "--declare-unscoped-run", run},
		{"migrate", "recover", "apply", "--dir", root, "--idea", "idea", "--kind", "cross-review", "--declare-unscoped-run", run},
		{"migrate", "declare", "--dir", root, "--idea", "idea", "--kind", "cross-review", "--declare-unscoped-run", run},
		{"worktree", "inspect", "--dir", root, "--declare-unscoped-run", run},
	}
	for _, args := range refused {
		if code := runBudgetPlatformControl(ctx, args, io.Discard, io.Discard, true, true); code != 2 {
			t.Fatalf("declaration accepted outside protocol inspect/apply: %v %d", args, code)
		}
	}
	// A value the accounting boundary cannot parse is refused there, not
	// repaired at the flag: an empty occurrence is refused by the flag itself.
	for _, bad := range []string{"parley-deck/runs/20260510T194003Z", "/abs/parley-deck/runs/x=" + strings.Repeat("0", 64), ""} {
		args := append(append([]string{}, worktreeOnly...), "--declare-unscoped-run", bad)
		if code := runBudgetPlatformControl(ctx, args, io.Discard, io.Discard, false, false); code == 0 {
			t.Fatalf("a malformed declaration was admitted: %q", bad)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "parley-launch-budgets")); !os.IsNotExist(err) {
		t.Fatalf("a read-only or refused command created accounting state: %v", err)
	}
}

// An attended apply carries both repeatable declarations into the immutable
// decision, replays exactly, and refuses once a declaration is dropped. The
// apply is local to this temporary repository: it is not an amendment of any
// real history, and it chooses no count beyond the observed floor.
func TestBudgetMigrateDeclaredUnscopedRunApplyReplaysExactly(t *testing.T) {
	ctx := context.Background()
	root, linked := declaredWorktreeCLIFixture(t)
	run := declaredUnscopedRunCLIFixture(t, root, "20260510T194003Z")
	var out, errout bytes.Buffer
	inspect := []string{"migrate", "inspect", "--dir", root, "--idea", "idea", "--kind", "cross-review",
		"--declare-unavailable-worktree", linked, "--declare-unscoped-run", run}
	if code := runBudgetPlatformControl(ctx, inspect, &out, &errout, false, false); code != 0 {
		t.Fatalf("declared inspect: %d %s", code, errout.String())
	}
	var i budget.ProtocolMigrationInventory
	if err := json.Unmarshal(out.Bytes(), &i); err != nil {
		t.Fatal(err)
	}
	// The declared run supplies no evidence, so it raises no floor and asserts
	// no zero: the count applied is the observed floor and nothing more.
	apply := []string{"migrate", "apply", "--dir", root, "--idea", "idea", "--kind", "cross-review",
		"--declare-unavailable-worktree", linked, "--declare-unscoped-run", run,
		"--expected-history-sha256", i.HistorySHA256,
		"--decision-id", "declared-unscoped-cli-fixture", "--reason", "Explicit fixture-only declared-unscoped protocol accounting",
		"--started-at", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano),
		"--total-actions", strconv.Itoa(i.LowerBound), "--max-cycles", "3", "--writers-stopped", "--yes"}
	if code := runBudgetPlatformControl(ctx, apply, io.Discard, &errout, true, true); code != 0 {
		t.Fatalf("declared apply: %s", errout.String())
	}
	if code := runBudgetPlatformControl(ctx, apply, io.Discard, &errout, true, true); code != 0 {
		t.Fatalf("exact declared replay: %s", errout.String())
	}
	b, err := budget.LoadCycleBinding(ctx, root, "idea", budget.CrossReview)
	if err != nil || b == nil {
		t.Fatalf("declared import is not loadable: %v", err)
	}
	if b.Policy.Maximum != 3 || b.Policy.Carried != 0 {
		t.Fatalf("declared import changed the frozen cap: %+v", b.Policy)
	}
	// Dropping either declaration is a different decision, not the same one.
	for _, drop := range []string{"--declare-unscoped-run", "--declare-unavailable-worktree"} {
		var narrower []string
		for n := 0; n < len(apply); n++ {
			if apply[n] == drop {
				n++
				continue
			}
			narrower = append(narrower, apply[n])
		}
		if code := runBudgetPlatformControl(ctx, narrower, io.Discard, io.Discard, true, true); code == 0 {
			t.Fatalf("an apply without %s replayed over a declared import", drop)
		}
	}
}
