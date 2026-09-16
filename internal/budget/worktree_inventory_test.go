package budget

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// A real repository with a real linked worktree. Every availability witness in
// this file uses actual Git registration state, never fabricated wire bytes.
func worktreeRepoFixture(t *testing.T) (root, linked string, git func(dir string, args ...string) string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git unavailable: %v", err)
	}
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git = func(dir string, args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
		return string(out)
	}
	root = filepath.Join(base, "main")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	git(root, "init", "-q")
	git(root, "-c", "user.name=fixture", "-c", "user.email=fixture@example.invalid", "commit", "--allow-empty", "-qm", "fixture")
	linked = filepath.Join(base, "linked")
	git(root, "worktree", "add", "-q", "-b", "linked", linked)
	return root, linked, git
}

// Field names are part of the contract: an unexpected one could smuggle a
// count, cap or attestation into a record that must carry none.
func assertWorktreeFields(t *testing.T, value any, names ...string) {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	expected := map[string]bool{}
	for _, name := range names {
		expected[name] = true
		if _, ok := fields[name]; !ok {
			t.Fatalf("missing field %q: %s", name, raw)
		}
	}
	for name := range fields {
		if !expected[name] {
			t.Fatalf("unexpected field %q: %s", name, raw)
		}
	}
}

func worktreeRow(t *testing.T, i WorktreeInventory, path string) WorktreeRegistration {
	t.Helper()
	for _, r := range i.Registrations {
		if r.RegisteredPath == path {
			return r
		}
	}
	t.Fatalf("registration %q was dropped from the inventory: %+v", path, i.Registrations)
	return WorktreeRegistration{}
}

// The load-bearing witness: a registration whose working directory is gone is
// retained as unavailable with unknown history, is not pruned or recreated, and
// does not weaken the existing refusal that shares the same Git output.
func TestWorktreeInventoryRetainsRemovedLinkedWorktree(t *testing.T) {
	ctx := context.Background()
	root, linked, git := worktreeRepoFixture(t)
	adminDir := filepath.Join(root, ".git", "worktrees")
	before, err := os.ReadDir(adminDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	inventory, err := InspectWorktreeInventory(ctx, root)
	if err != nil {
		t.Fatalf("inventory refused a removed registration instead of retaining it: %v", err)
	}
	if !inventory.Git || !inventory.CoversRepository {
		t.Fatalf("repository inventory does not claim Git coverage: %+v", inventory)
	}
	removed := worktreeRow(t, inventory, linked)
	if removed.Availability != WorktreeMissing || removed.HistoryCoverage != WorktreeHistoryUnknown {
		t.Fatalf("removed registration is not retained as unavailable/unknown: %+v", removed)
	}
	if removed.Head == "" || removed.Branch != "refs/heads/linked" {
		t.Fatalf("removed registration lost its raw registered identity: %+v", removed)
	}
	if removed.ResolvedPath != "" {
		t.Fatalf("missing registration reports a resolved path: %+v", removed)
	}
	if present := worktreeRow(t, inventory, root); present.Availability != WorktreeAvailable || present.HistoryCoverage != WorktreeHistoryLocal {
		t.Fatalf("surviving registration is not available/local: %+v", present)
	}
	if len(inventory.Uncertainty) == 0 {
		t.Fatal("an inventory holding an unavailable registration reports no uncertainty")
	}
	// The emitted schema is an exact, closed key set. No count, cap, floor or
	// attestation field exists, so an unavailable worktree cannot be read as
	// zero history and nothing here infers an unauthorized quantity.
	assertWorktreeFields(t, inventory, "schema", "root", "git", "git_common_dir", "repo_prefix", "source_scope", "covers_repository", "registrations", "uncertainty", "observation_sha256")
	assertWorktreeFields(t, removed, "registered_path", "resolved_path", "head", "branch", "detached", "bare", "locked", "lock_reason", "prunable", "prunable_reason", "availability", "history_coverage")

	// Nothing was mutated: the path stays gone, the registration stays listed,
	// the admin directory is intact and no runtime/policy state was created.
	if _, err := os.Lstat(linked); !os.IsNotExist(err) {
		t.Fatalf("inspection recreated the missing worktree path: %v", err)
	}
	if listed := git(root, "worktree", "list", "--porcelain"); !strings.Contains(listed, linked) {
		t.Fatalf("inspection pruned the registration: %s", listed)
	}
	after, err := os.ReadDir(adminDir)
	if err != nil || len(after) != len(before) {
		t.Fatalf("worktree admin directory changed: %v %d -> %d", err, len(before), len(after))
	}
	for _, path := range []string{filepath.Join(root, ".parley-runtime"), filepath.Join(root, ".git", "parley-launch-budgets")} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("read-only inspection created %s: %v", path, err)
		}
	}

	// The existing refusal reads the same Git output and is unchanged: shared
	// scope resolution still stops on unavailable history.
	if _, err := ConfigureLaunchBudget(ctx, root, "idea", LaunchPolicy{MaxLaunches: 1}); err == nil || !strings.Contains(err.Error(), "historical worktree is unavailable") {
		t.Fatalf("launch scope history refusal was weakened: %v", err)
	}
}

// Identical previews must hash identically, so the digest carries no clock.
func TestWorktreeInventoryDigestIsStableAndClockFree(t *testing.T) {
	ctx := context.Background()
	root, linked, _ := worktreeRepoFixture(t)
	first, err := InspectWorktreeInventory(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	second, err := InspectWorktreeInventory(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	if first.ObservationSHA256 == "" || first.ObservationSHA256 != second.ObservationSHA256 {
		t.Fatalf("unchanged state produced different digests: %q %q", first.ObservationSHA256, second.ObservationSHA256)
	}
	a, err := first.OperatorRecord()
	if err != nil {
		t.Fatal(err)
	}
	b, err := second.OperatorRecord()
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatalf("operator record is not reproducible:\n%s\n%s", a, b)
	}
	// The digest input is an exact, closed key set: neither a clock nor the
	// caller's own location can enter it unnoticed.
	assertWorktreeFields(t, observedWorktrees{}, "schema", "git", "covers_repository", "source_scope", "git_common_dir", "registrations")
	// The same repository is one observation whether it is inspected from the
	// main worktree, from a linked worktree, or from a subdirectory.
	sub := filepath.Join(root, "nested")
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, from := range []string{linked, sub} {
		elsewhere, err := InspectWorktreeInventory(ctx, from)
		if err != nil {
			t.Fatal(err)
		}
		if elsewhere.ObservationSHA256 != first.ObservationSHA256 {
			t.Fatalf("digest depends on where it was inspected from (%s): %q %q", from, elsewhere.ObservationSHA256, first.ObservationSHA256)
		}
	}
}

func TestWorktreeInventoryDigestTracksRegistrationAndAvailability(t *testing.T) {
	ctx := context.Background()
	root, linked, git := worktreeRepoFixture(t)
	original, err := InspectWorktreeInventory(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	git(root, "worktree", "add", "-q", "-b", "third", filepath.Join(filepath.Dir(root), "third"))
	added, err := InspectWorktreeInventory(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	if added.ObservationSHA256 == original.ObservationSHA256 {
		t.Fatal("a new registration did not change the digest")
	}
	if len(added.Registrations) != len(original.Registrations)+1 {
		t.Fatalf("registration count did not grow: %d -> %d", len(original.Registrations), len(added.Registrations))
	}
	// Move the directory away and back: a pure availability change, with the
	// registration itself byte-identical either way.
	away := linked + ".away"
	if err := os.Rename(linked, away); err != nil {
		t.Fatal(err)
	}
	unavailable, err := InspectWorktreeInventory(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	if unavailable.ObservationSHA256 == added.ObservationSHA256 {
		t.Fatal("losing a registered working directory did not change the digest")
	}
	if row := worktreeRow(t, unavailable, linked); row.Availability != WorktreeMissing {
		t.Fatalf("moved-away registration is not missing: %+v", row)
	}
	if err := os.Rename(away, linked); err != nil {
		t.Fatal(err)
	}
	restored, err := InspectWorktreeInventory(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	if restored.ObservationSHA256 != added.ObservationSHA256 {
		t.Fatalf("restoring the exact prior state did not restore the digest: %q %q", restored.ObservationSHA256, added.ObservationSHA256)
	}
}

// A root with no Git authority anywhere above it is reported explicitly and
// claims no repository-wide coverage, rather than being silently empty.
func TestWorktreeInventoryNonGitRootIsExplicit(t *testing.T) {
	ctx := context.Background()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	plain := filepath.Join(dir, "plain")
	if err := os.MkdirAll(plain, 0o700); err != nil {
		t.Fatal(err)
	}
	inventory, err := InspectWorktreeInventory(ctx, plain)
	if err != nil {
		t.Skipf("temporary directory is inside a repository: %v", err)
	}
	if inventory.Git || inventory.CoversRepository {
		t.Fatalf("non-Git root claims repository coverage: %+v", inventory)
	}
	if len(inventory.Registrations) != 0 || inventory.SourceScope != "single-directory/non-git" {
		t.Fatalf("non-Git root did not report an explicit single-directory scope: %+v", inventory)
	}
	if len(inventory.Uncertainty) == 0 || inventory.ObservationSHA256 == "" {
		t.Fatalf("non-Git root is not explicit about its own limits: %+v", inventory)
	}
	root, _, _ := worktreeRepoFixture(t)
	repo, err := InspectWorktreeInventory(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	if repo.ObservationSHA256 == inventory.ObservationSHA256 {
		t.Fatal("a non-Git observation collides with a repository observation")
	}
}

// Broken Git authority must not degrade into a non-Git answer that silently
// drops every registration the repository still holds.
func TestWorktreeInventoryRefusesBrokenGitAuthority(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(dir, "broken")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git"), []byte("not a gitdir pointer\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	inventory, err := InspectWorktreeInventory(context.Background(), root)
	if err == nil {
		t.Fatalf("broken Git authority produced an inventory: %+v", inventory)
	}
	if !strings.Contains(err.Error(), "cannot resolve Git worktree authority") {
		t.Fatalf("unexpected refusal: %v", err)
	}
	if inventory.Git || len(inventory.Registrations) != 0 || inventory.ObservationSHA256 != "" {
		t.Fatalf("refusal returned a usable inventory: %+v", inventory)
	}
}

func worktreeWire(records ...[]string) []byte {
	var b strings.Builder
	for _, record := range records {
		for _, attribute := range record {
			b.WriteString(attribute)
			b.WriteByte(0)
		}
		b.WriteByte(0)
	}
	return []byte(b.String())
}

// Fabricated wire bytes only, for shapes real Git will not emit on demand.
// They supplement, and never replace, the real missing-worktree witness above.
func TestWorktreeRegistrationRefusesMalformedAuthority(t *testing.T) {
	sha := strings.Repeat("a", 40)
	for name, wire := range map[string][]byte{
		"unknown attribute":     worktreeWire([]string{"worktree /a", "HEAD " + sha, "surprise yes"}),
		"duplicate attribute":   worktreeWire([]string{"worktree /a", "HEAD " + sha, "HEAD " + sha}),
		"duplicate path":        worktreeWire([]string{"worktree /a", "HEAD " + sha, "detached"}, []string{"worktree /a", "HEAD " + sha, "detached"}),
		"relative path":         worktreeWire([]string{"worktree relative/a", "HEAD " + sha, "detached"}),
		"empty path":            worktreeWire([]string{"worktree ", "HEAD " + sha, "detached"}),
		"valueless worktree":    worktreeWire([]string{"worktree", "HEAD " + sha, "detached"}),
		"branch and detached":   worktreeWire([]string{"worktree /a", "HEAD " + sha, "branch refs/heads/x", "detached"}),
		"bare and checked out":  worktreeWire([]string{"worktree /a", "HEAD " + sha, "bare"}),
		"neither bare nor head": worktreeWire([]string{"worktree /a"}),
		"invalid head":          worktreeWire([]string{"worktree /a", "HEAD nothex", "detached"}),
		"short head":            worktreeWire([]string{"worktree /a", "HEAD " + sha[:39], "detached"}),
		"unqualified branch":    worktreeWire([]string{"worktree /a", "HEAD " + sha, "branch main"}),
		"orphan attribute":      worktreeWire([]string{"HEAD " + sha}),
		"valued flag":           worktreeWire([]string{"worktree /a", "HEAD " + sha, "detached yes"}),
		"nested registration":   worktreeWire([]string{"worktree /a", "worktree /b"}),
		"empty output":          {},
		"unterminated record":   []byte("worktree /a\x00HEAD " + sha),
	} {
		t.Run(name, func(t *testing.T) {
			if got, err := parseWorktreeRegistrations(wire); err == nil {
				t.Fatalf("malformed registration authority accepted: %+v", got)
			}
		})
	}
}

// The strict parser must still accept every shape real Git does emit.
func TestWorktreeRegistrationAcceptsStandardAttributes(t *testing.T) {
	sha := strings.Repeat("b", 40)
	got, err := parseWorktreeRegistrations(worktreeWire(
		[]string{"worktree /repo", "bare"},
		[]string{"worktree /a", "HEAD " + sha, "branch refs/heads/main"},
		[]string{"worktree /b", "HEAD " + sha, "detached", "locked on removable media"},
		[]string{"worktree /c", "HEAD " + sha, "detached", "prunable gitdir file points to non-existent location"},
		[]string{"worktree /d", "HEAD " + strings.Repeat("c", 64), "detached", "locked"},
	))
	if err != nil {
		t.Fatalf("standard Git registration output refused: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("expected five registrations: %+v", got)
	}
	if !got[0].Bare || got[0].Head != "" {
		t.Fatalf("bare registration mis-parsed: %+v", got[0])
	}
	if got[1].Branch != "refs/heads/main" || got[1].Detached {
		t.Fatalf("branch registration mis-parsed: %+v", got[1])
	}
	if !got[2].Locked || got[2].LockReason != "on removable media" {
		t.Fatalf("lock metadata lost: %+v", got[2])
	}
	if !got[3].Prunable || got[3].PrunableReason == "" {
		t.Fatalf("prunable metadata lost: %+v", got[3])
	}
	if !got[4].Locked || got[4].LockReason != "" {
		t.Fatalf("reasonless lock mis-parsed: %+v", got[4])
	}
}
