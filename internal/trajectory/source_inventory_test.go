package trajectory

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/evidence"
)

func TestSourceInventoryHiddenGoInputRefusesObservationAndArchive(t *testing.T) {
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for _, kind := range []string{"local", "global"} {
		t.Run(kind, func(t *testing.T) {
			root, store := snapshotFixture(t), snapshotStoreFixture(t)
			snapshotWrite(t, root, "go.mod", []byte("module source-inventory-fixture\n\ngo 1.22\n"), 0600)
			snapshotWrite(t, root, "main.go", []byte("package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"baseline\")}\n"), 0600)
			gitFixture(t, root, "add", ".")
			gitFixture(t, root, "commit", "-qm", "Go source fixture")
			before, err := Observe(context.Background(), root)
			if err != nil {
				t.Fatal(err)
			}
			runGo := func() string {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				cmd := exec.CommandContext(ctx, "go", "run", ".")
				cmd.Dir = root
				cmd.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off")
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("actual Go fixture failed: %v: %s", err, out)
				}
				return string(out)
			}
			if got := runGo(); got != "baseline\n" {
				t.Fatalf("baseline fixture mismatch: %q", got)
			}
			snapshotWrite(t, root, "extra.go", []byte("package main\nimport \"fmt\"\nfunc init(){fmt.Println(\"hidden\")}\n"), 0600)
			ignore := filepath.Join(root, ".git/info/exclude")
			if kind == "global" {
				ignore = filepath.Join(t.TempDir(), "ignore")
				config := filepath.Join(t.TempDir(), "gitconfig")
				gitFixture(t, root, "config", "--file", config, "core.excludesFile", ignore)
				t.Setenv("GIT_CONFIG_GLOBAL", config)
			}
			if err := os.WriteFile(ignore, []byte("extra.go\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if got := runGo(); got != "hidden\nbaseline\n" {
				t.Fatalf("hidden file did not affect real execution: %q", got)
			}
			if got := gitFixture(t, root, "status", "--porcelain=v1", "--untracked-files=all"); got != "" {
				t.Fatalf("counterexample was not Git-clean: %q", got)
			}
			if _, _, err := snapshotInventory(context.Background(), root); !errors.Is(err, evidence.ErrLocalSourceExcludes) {
				t.Fatalf("snapshot inventory certified hidden Go input: %v", err)
			}
			if got, err := evidence.TreeDigest(root); got != "" || !errors.Is(err, evidence.ErrLocalSourceExcludes) {
				t.Fatalf("digest certified hidden Go input: %q %v", got, err)
			}
			if got, err := Observe(context.Background(), root); got != (Source{}) || !errors.Is(err, evidence.ErrLocalSourceExcludes) {
				t.Fatalf("Observe certified hidden Go input: %+v %v", got, err)
			}
			if got, err := CaptureSnapshot(context.Background(), root, store, before); got != (SnapshotRef{}) || !errors.Is(err, evidence.ErrLocalSourceExcludes) {
				t.Fatalf("capture certified hidden Go input: %+v %v", got, err)
			}
			entries, err := os.ReadDir(store)
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), ".tar") {
					t.Fatal("refused source published an archive")
				}
			}
			// Deliberately making this an explicit project exclusion is a visible
			// scope change. It is permitted but is not proof the build ignores it.
			snapshotWrite(t, root, ".gitignore", []byte("ignored/\nextra.go\n"), 0600)
			after, err := Observe(context.Background(), root)
			if err != nil || after.Tree.SHA256 == before.Tree.SHA256 || after.Clean {
				t.Fatalf("project scope change not retained: %+v %v", after, err)
			}
		})
	}
}

func TestSourceInventoryRetainsExplicitScopeOnRestore(t *testing.T) {
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	root, store := snapshotFixture(t), snapshotStoreFixture(t)
	snapshotWrite(t, root, "nested/.gitignore", []byte("*.local\n!keep.local\n"), 0600)
	snapshotWrite(t, root, "nested/hide.local", []byte("omitted\n"), 0600)
	snapshotWrite(t, root, "nested/keep.local", []byte("retained\n"), 0600)
	snapshotWrite(t, root, ".git/info/exclude", []byte("ignored/\n*.local\n"), 0600)
	if err := os.Remove(filepath.Join(root, "removed")); err != nil {
		t.Fatal(err)
	}
	source, ref := captureFixture(t, root, store)
	restored, err := RestoreSnapshot(context.Background(), store, ref, source, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"nested/hide.local", "removed"} {
		if _, err := os.Stat(filepath.Join(restored, name)); !os.IsNotExist(err) {
			t.Fatalf("excluded/absent source restored: %s %v", name, err)
		}
	}
	if got := string(snapshotRead(t, filepath.Join(restored, "nested/keep.local"))); got != "retained\n" {
		t.Fatal("negated project file lost")
	}
	if got, err := evidence.TreeDigest(restored); err != nil || got != source.Tree.SHA256 {
		t.Fatalf("restored digest drifted: %s %v", got, err)
	}
}
