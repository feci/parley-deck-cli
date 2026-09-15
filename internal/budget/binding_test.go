package budget

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestLaunchBindingFreezesPolicyAndPreservesHistory(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	policy := LaunchPolicy{MaxLaunches: 1}
	b, err := ConfigureLaunchBudget(ctx, root, "idea", policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Store.Reserve(ctx, Request{ID: "spent", Kind: Launch}, b.Policy.Limits()); err != nil {
		t.Fatal(err)
	}
	prior, _ := b.Store.Inspect(ctx)
	replay, err := ConfigureLaunchBudget(ctx, root, "idea", policy)
	if err != nil {
		t.Fatal(err)
	}
	now, _ := replay.Store.Inspect(ctx)
	if len(now.Entries) != 1 || !now.StartedAt.Equal(prior.StartedAt) {
		t.Fatal("configure replay reset count or clock")
	}
	policy.MaxLaunches = 2
	if _, err := ConfigureLaunchBudget(ctx, root, "idea", policy); err == nil {
		t.Fatal("configure extended a frozen policy")
	}
	loaded, err := LoadLaunchBinding(ctx, root, "idea")
	if err != nil || loaded == nil {
		t.Fatalf("load: %v", err)
	}
	if _, err := loaded.Store.Reserve(ctx, Request{ID: "other-run", Kind: Launch}, loaded.Policy.Limits()); !errors.Is(err, ErrLimit) {
		t.Fatalf("new run reset policy: %v", err)
	}
	path := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadLaunchBinding(ctx, root, "idea"); err == nil {
		t.Fatal("missing configured policy became unbudgeted")
	}
	if _, err := ConfigureLaunchBudget(ctx, root, "idea", policy); err == nil {
		t.Fatal("missing policy allowed an extension around retained ledger")
	}
}

func TestLaunchBindingRefusesIncompleteHistory(t *testing.T) {
	for _, mode := range []string{"missing-request", "history-is-file", "aliased-history"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			history := filepath.Join(root, ".parley-runtime", "invocations")
			if err := os.MkdirAll(filepath.Dir(history), 0700); err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "missing-request":
				if err := os.MkdirAll(filepath.Join(history, "old"), 0700); err != nil {
					t.Fatal(err)
				}
			case "history-is-file":
				if err := os.WriteFile(history, []byte("lost history"), 0600); err != nil {
					t.Fatal(err)
				}
			case "aliased-history":
				if err := os.Symlink(t.TempDir(), history); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := ConfigureLaunchBudget(context.Background(), root, "idea", LaunchPolicy{MaxLaunches: 1}); err == nil {
				t.Fatal("incomplete or aliased history allowed a fresh policy")
			}
			if b, err := LoadLaunchBinding(context.Background(), root, "idea"); err != nil || b != nil {
				t.Fatalf("legacy preflight left partial activation: %+v %v", b, err)
			}
		})
	}
}

func TestLaunchBindingConcurrentConfigurationHasOneFrozenWinner(t *testing.T) {
	root := t.TempDir()
	start := make(chan struct{})
	results := make(chan *LaunchBinding, 2)
	var wg sync.WaitGroup
	for _, count := range []int{1, 2} {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			<-start
			b, err := ConfigureLaunchBudget(context.Background(), root, "idea", LaunchPolicy{MaxLaunches: count})
			if err != nil {
				results <- nil
				return
			}
			results <- b
		}(count)
	}
	close(start)
	wg.Wait()
	close(results)
	var winner *LaunchBinding
	for result := range results {
		if result != nil {
			if winner != nil {
				t.Fatal("conflicting policies both succeeded")
			}
			winner = result
		}
	}
	if winner == nil {
		t.Fatal("neither first policy succeeded")
	}
	loaded, err := LoadLaunchBinding(context.Background(), root, "idea")
	if err != nil || loaded == nil || loaded.Policy.MaxLaunches != winner.Policy.MaxLaunches {
		t.Fatalf("persisted policy differs from winner: %+v %v", loaded, err)
	}
}

func TestLaunchBindingSharedAcrossGitWorktrees(t *testing.T) {
	root := t.TempDir()
	git := func(dir string, args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git: %v %s", err, out)
		}
	}
	git(root, "init", "-q")
	git(root, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "--allow-empty", "-qm", "fixture")
	other := filepath.Join(t.TempDir(), "other")
	git(root, "worktree", "add", "-q", "-b", "other", other)
	b, err := ConfigureLaunchBudget(context.Background(), root, "idea", LaunchPolicy{MaxLaunches: 1})
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadLaunchBinding(context.Background(), other, "idea")
	if err != nil || loaded == nil {
		t.Fatalf("worktree has no policy: %v", err)
	}
	if loaded.Store.Dir != b.Store.Dir || loaded.Store.Scope != b.Store.Scope {
		t.Fatalf("worktrees use different ledgers: %+v %+v", b.Store, loaded.Store)
	}
	if _, err := b.Store.Reserve(context.Background(), Request{ID: "root-one", Kind: Launch}, b.Policy.Limits()); err != nil {
		t.Fatal(err)
	}
	if _, err := loaded.Store.Reserve(context.Background(), Request{ID: "root-two", Kind: Launch}, loaded.Policy.Limits()); !errors.Is(err, ErrLimit) {
		t.Fatalf("worktree bypassed cap: %v", err)
	}
	// A deck may exist in only one branch. Its absence in another accessible
	// worktree is not missing history, and its prefix must not borrow root cap.
	nested := filepath.Join(root, "projects", "nested")
	if err := os.MkdirAll(nested, 0700); err != nil {
		t.Fatal(err)
	}
	n, err := ConfigureLaunchBudget(context.Background(), nested, "idea", LaunchPolicy{MaxLaunches: 2})
	if err != nil || n == nil || n.Store.Scope == b.Store.Scope {
		t.Fatalf("nested-only deck configuration: %+v %v", n, err)
	}
	otherNested := filepath.Join(other, "projects", "nested")
	if err := os.MkdirAll(otherNested, 0700); err != nil {
		t.Fatal(err)
	}
	inherited, err := LoadLaunchBinding(context.Background(), otherNested, "idea")
	if err != nil || inherited == nil || inherited.Store.Dir != n.Store.Dir {
		t.Fatalf("later nested worktree did not inherit policy: %+v %v", inherited, err)
	}
	if err := os.RemoveAll(other); err != nil {
		t.Fatal(err)
	}
	if existing, err := LoadLaunchBinding(context.Background(), root, "idea"); err != nil || existing == nil {
		t.Fatalf("existing policy unnecessarily depends on other worktree availability: %+v %v", existing, err)
	}
	if _, err := ConfigureLaunchBudget(context.Background(), root, "new-idea", LaunchPolicy{MaxLaunches: 1}); err == nil {
		t.Fatal("first activation ignored an unavailable historical worktree")
	}
}

func TestLaunchBindingRefusesUnmigratedAndMalformedPolicy(t *testing.T) {
	root := t.TempDir()
	p := filepath.Join(root, ".parley-runtime/invocations/old/requested.json")
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(`{"metadata":{"idea":"old-idea"}}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfigureLaunchBudget(context.Background(), root, "old-idea", LaunchPolicy{MaxLaunches: 1}); err == nil || !strings.Contains(err.Error(), "migration") {
		t.Fatalf("legacy reset: %v", err)
	}
	b, err := ConfigureLaunchBudget(context.Background(), root, "fresh-idea", LaunchPolicy{MaxLaunches: 1})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{strings.Replace(string(original), `"max_launches": 1,`, "", 1), strings.Replace(string(original), `"max_launches": 1`, `"Max_Launches": 1`, 1), strings.Replace(string(original), `"max_launches": 1`, `"max_launches": 1, "max_launches": 0`, 1)} {
		if err := os.WriteFile(path, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadLaunchBinding(context.Background(), root, "fresh-idea"); err == nil {
			t.Fatal("malformed policy became an unlimited grant")
		}
	}
}
