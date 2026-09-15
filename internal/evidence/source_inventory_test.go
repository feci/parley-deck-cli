package evidence

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func isolateInventoryGit(t *testing.T) {
	t.Helper()
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
}

func inventoryWrite(t *testing.T, root, name, body string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestGitSourceInventoryRefusesPrivateExclusions(t *testing.T) {
	for _, tc := range []string{"info-exclude", "global-config", "global-default", "hidden-directory", "hidden-project-rules"} {
		t.Run(tc, func(t *testing.T) {
			isolateInventoryGit(t)
			root := scratchGitRepo(t, map[string]string{"source": "baseline\n"})
			name, rule := "extra.go", "extra.go\n"
			if tc == "hidden-directory" {
				name, rule = "hidden/extra.go", "hidden/\n"
			}
			if tc == "hidden-project-rules" {
				name, rule = ".gitignore", ".gitignore\n"
			}
			inventoryWrite(t, root, name, "private fixture bytes\n")
			switch tc {
			case "global-config":
				config := filepath.Join(t.TempDir(), "config")
				ignore := filepath.Join(t.TempDir(), "ignore")
				inventoryWrite(t, filepath.Dir(ignore), "ignore", rule)
				git(t, root, "config", "--file", config, "core.excludesFile", ignore)
				t.Setenv("GIT_CONFIG_GLOBAL", config)
			case "global-default":
				inventoryWrite(t, os.Getenv("XDG_CONFIG_HOME"), "git/ignore", rule)
			default:
				inventoryWrite(t, root, ".git/info/exclude", rule)
			}
			if got, err := GitSourceInventory(context.Background(), root); !errors.Is(err, ErrLocalSourceExcludes) || got != nil {
				t.Fatalf("private exclusion must refuse inventory: %v %v", got, err)
			}
			if got, err := TreeDigest(root, name); !errors.Is(err, ErrLocalSourceExcludes) || got != "" {
				t.Fatalf("explicit evidence exclusion must not bypass ambiguous source: %q %v", got, err)
			}
		})
	}
}

func TestGitSourceInventoryPreservesProjectScope(t *testing.T) {
	isolateInventoryGit(t)
	root := scratchGitRepo(t, map[string]string{
		".gitignore":  "ignored/\n*.tmp\n!kept.tmp\n",
		"tracked.tmp": "tracked even when ignored\n",
		"source":      "baseline\n",
	})
	// scratchGitRepo follows git add's ordinary exclusions: explicitly force the
	// tracked ignored fixture into this isolated repository's index.
	git(t, root, "add", "-f", "tracked.tmp")
	git(t, root, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-qm", "tracked ignored")
	inventoryWrite(t, root, "ignored/private", "excluded by project\n")
	inventoryWrite(t, root, "lost.tmp", "excluded by project\n")
	inventoryWrite(t, root, "kept.tmp", "project negation overrides private rule\n")
	inventoryWrite(t, root, "nested/.gitignore", "*.local\n!keep.local\n")
	inventoryWrite(t, root, "nested/hide.local", "excluded by untracked project rules\n")
	inventoryWrite(t, root, "nested/keep.local", "included by nested negation\n")
	inventoryWrite(t, root, ".git/info/exclude", "ignored/\n*.tmp\n*.local\nnonexistent\n")
	want := []string{".gitignore", "kept.tmp", "nested/.gitignore", "nested/keep.local", "source", "tracked.tmp"}
	got, err := GitSourceInventory(context.Background(), root)
	if err != nil || !slices.Equal(got, want) {
		t.Fatalf("project scope changed: %v %v", got, err)
	}
	before, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	inventoryWrite(t, root, "ignored/private", "still omitted\n")
	after, err := TreeDigest(root)
	if err != nil || before != after {
		t.Fatalf("project ignored content entered digest: %v", err)
	}
	inventoryWrite(t, root, "tracked.tmp", "tracked change\n")
	after, err = TreeDigest(root)
	if err != nil || before == after {
		t.Fatalf("tracked ignored file omitted: %v", err)
	}
	// An untracked project's ignore rules are themselves ordinary source bytes.
	before = after
	inventoryWrite(t, root, "nested/.gitignore", "*.local\n!keep.local\n# changed\n")
	after, err = TreeDigest(root)
	if err != nil || before == after {
		t.Fatalf("untracked project rules omitted: %v", err)
	}
}

func TestGitSourceInventoryWorktreeAndDeletedFile(t *testing.T) {
	isolateInventoryGit(t)
	root := scratchGitRepo(t, map[string]string{"source": "baseline\n", "deleted": "original\n"})
	wt := filepath.Join(t.TempDir(), "worktree")
	git(t, root, "worktree", "add", "--detach", wt, "HEAD")
	if err := os.Remove(filepath.Join(wt, "deleted")); err != nil {
		t.Fatal(err)
	}
	got, err := GitSourceInventory(context.Background(), wt)
	if err != nil || !slices.Equal(got, []string{"deleted", "source"}) {
		t.Fatalf("tracked absence lost: %v %v", got, err)
	}
	inventoryWrite(t, wt, "extra.go", "hidden\n")
	inventoryWrite(t, root, ".git/info/exclude", "extra.go\n")
	if _, err := GitSourceInventory(context.Background(), wt); !errors.Is(err, ErrLocalSourceExcludes) {
		t.Fatalf("worktree common exclusions bypassed: %v", err)
	}
}

func TestGitSourceInventoryPathsBoundsAndCancellation(t *testing.T) {
	isolateInventoryGit(t)
	root := scratchGitRepo(t, map[string]string{"source": "baseline\n"})
	for _, name := range []string{"space name", "žltý", "-option", "line\nbreak"} {
		inventoryWrite(t, root, name, "included\n")
	}
	want := []string{"-option", "line\nbreak", "source", "space name", "žltý"}
	got, err := GitSourceInventory(context.Background(), root)
	if err != nil || !slices.Equal(got, want) {
		t.Fatalf("NUL-delimited names changed: %v %v", got, err)
	}
	for _, tc := range []struct {
		name           string
		entries, bytes int
		message        string
	}{
		{"entries", 4, MaxSourceInventoryBytes, "entry count exceeds bound"},
		{"bytes", MaxSourceInventoryEntries, 1, "output exceeds bound"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := readGitSourceInventory(context.Background(), root, "--exclude-standard", tc.entries, tc.bytes); err == nil || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("inventory must enforce %s: %v", tc.name, err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := GitSourceInventory(ctx, root); !errors.Is(err, context.Canceled) {
		t.Fatalf("inventory ignored cancellation: %v", err)
	}
	if _, err := TreeDigestContext(ctx, root); !errors.Is(err, context.Canceled) {
		t.Fatalf("digest ignored cancellation: %v", err)
	}
	if _, err := GitSourceInventory(context.Background(), filepath.Join(root, "private-missing-path")); err == nil || strings.Contains(err.Error(), "private-missing-path") {
		t.Fatalf("Git failure leaked path or passed: %v", err)
	}
}

func TestGitSourceInventoryPreservesIndexAndConfiguration(t *testing.T) {
	isolateInventoryGit(t)
	root := scratchGitRepo(t, map[string]string{"source": "baseline\n"})
	inventoryWrite(t, root, "new", "untracked\n")
	paths := []string{".git/index", ".git/config", ".git/HEAD", ".git/info/exclude"}
	before := map[string]string{}
	for _, p := range paths {
		b, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			t.Fatal(err)
		}
		before[p] = string(b)
	}
	if _, err := GitSourceInventory(context.Background(), root); err != nil {
		t.Fatal(err)
	}
	for _, p := range paths {
		b, err := os.ReadFile(filepath.Join(root, p))
		if err != nil || string(b) != before[p] {
			t.Fatalf("read changed Git file %s: %v", p, err)
		}
	}
	// Verify the fixture really used an untracked path; no implicit staging.
	cmd := exec.Command("git", "-C", root, "ls-files", "--error-unmatch", "new")
	if err := cmd.Run(); err == nil {
		t.Fatal("inventory staged source")
	}
}
