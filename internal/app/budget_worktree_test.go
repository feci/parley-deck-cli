package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
)

// A real repository whose linked worktree directory has been removed, leaving
// the registration behind. The CLI must report it, not prune or repair it.
func worktreeCLIFixture(t *testing.T) (root, removed string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git unavailable: %v", err)
	}
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root = filepath.Join(base, "main")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	git("init", "-q")
	git("-c", "user.name=fixture", "-c", "user.email=fixture@example.invalid", "commit", "--allow-empty", "-qm", "fixture")
	removed = filepath.Join(base, "gone")
	git("worktree", "add", "-q", "-b", "gone", removed)
	if err := os.RemoveAll(removed); err != nil {
		t.Fatal(err)
	}
	return root, removed
}

func TestBudgetWorktreeInspectCLIReportsRetainedRegistration(t *testing.T) {
	root, removed := worktreeCLIFixture(t)
	var out, errout bytes.Buffer
	// Unsupported and unattended: a read-only inventory needs no attended
	// operator control, and none may have been introduced.
	if code := runBudgetPlatformControl(context.Background(), []string{"worktree", "inspect", "--dir", root}, &out, &errout, false, false); code != 0 {
		t.Fatalf("read-only inspect failed: %d %s", code, errout.String())
	}
	var inventory budget.WorktreeInventory
	dec := json.NewDecoder(bytes.NewReader(out.Bytes()))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&inventory); err != nil {
		t.Fatalf("inspect output is not a WorktreeInventory: %v %s", err, out.String())
	}
	if dec.More() {
		t.Fatalf("inspect emitted more than one record: %s", out.String())
	}
	if !inventory.Git || inventory.ObservationSHA256 == "" {
		t.Fatalf("inspect did not report Git authority and a digest: %+v", inventory)
	}
	var found bool
	for _, r := range inventory.Registrations {
		if r.RegisteredPath != removed {
			continue
		}
		found = true
		if r.Availability != budget.WorktreeMissing || r.HistoryCoverage != budget.WorktreeHistoryUnknown {
			t.Fatalf("removed registration is not unavailable/unknown: %+v", r)
		}
	}
	if !found {
		t.Fatalf("CLI dropped the unavailable registration: %+v", inventory.Registrations)
	}
	if len(inventory.Uncertainty) == 0 {
		t.Fatal("CLI hid the residual uncertainty")
	}
	if _, err := os.Lstat(removed); !os.IsNotExist(err) {
		t.Fatalf("CLI recreated the missing worktree path: %v", err)
	}
	// A second invocation of the unchanged state is byte-identical, so an
	// operator can compare two previews without a clock invalidating them.
	var again bytes.Buffer
	if code := runBudgetPlatformControl(context.Background(), []string{"worktree", "inspect", "--dir", root}, &again, io.Discard, false, false); code != 0 {
		t.Fatalf("repeated inspect failed: %d", code)
	}
	if again.String() != out.String() {
		t.Fatalf("repeated inspect differs:\n%s\n%s", out.String(), again.String())
	}
}

// No mutation, attestation or apply surface exists, and inspect cannot be
// handed an action to perform.
func TestBudgetWorktreeRefusesMutationAndExtraArguments(t *testing.T) {
	root, _ := worktreeCLIFixture(t)
	for _, args := range [][]string{
		{"worktree"},
		{"worktree", "apply", "--dir", root},
		{"worktree", "attest", "--dir", root},
		{"worktree", "prune", "--dir", root},
		{"worktree", "repair", "--dir", root},
		{"worktree", "recover", "--dir", root},
		{"worktree", "remove", "--dir", root},
		{"worktree", "inspect", "--dir", root, "apply"},
		{"worktree", "inspect", "--dir", root, "--yes"},
		{"worktree", "inspect", "--dir", root, "--apply"},
		{"worktree", "inspect", "--dir", root, "--attest"},
		{"worktree", "inspect", "--dir", root, "--decision-id", "x"},
		{"worktree", "inspect", "--dir", root, "--idea", "idea"},
		{"worktree", "inspect", "--dir", root, "--prune"},
		{"worktree", "inspect", "--dir", root, "--expected-observation-sha256", "x"},
	} {
		var out, errout bytes.Buffer
		// Attended and supported: even the strongest operator context must not
		// unlock a mutation verb that does not exist.
		if code := runBudgetPlatformControl(context.Background(), args, &out, &errout, true, true); code != 2 {
			t.Fatalf("%v returned %d (stderr %q)", args, code, errout.String())
		}
		if out.Len() != 0 {
			t.Fatalf("%v wrote a result: %s", args, out.String())
		}
	}
	for _, path := range []string{filepath.Join(root, ".parley-runtime"), filepath.Join(root, ".git", "parley-launch-budgets")} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("refused control created %s: %v", path, err)
		}
	}
	// The usage text must not advertise a verb this stage does not implement.
	var errout bytes.Buffer
	runBudgetPlatformControl(context.Background(), []string{"worktree"}, io.Discard, &errout, true, true)
	for _, banned := range []string{"apply", "attest", "--yes", "recover", "prune"} {
		if strings.Contains(errout.String(), banned) {
			t.Fatalf("usage advertises %q: %s", banned, errout.String())
		}
	}
}

func TestBudgetWorktreeInspectRefusesBrokenAuthority(t *testing.T) {
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(base, "broken")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git"), []byte("not a gitdir pointer\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errout bytes.Buffer
	if code := runBudgetPlatformControl(context.Background(), []string{"worktree", "inspect", "--dir", root}, &out, &errout, false, false); code != 1 {
		t.Fatalf("broken authority returned %d: %s", code, errout.String())
	}
	if out.Len() != 0 {
		t.Fatalf("broken authority produced a result: %s", out.String())
	}
}
