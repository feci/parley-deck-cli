package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceSyncedFilePreservesOldFileOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state")
	if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ReplaceSyncedFile(filepath.Join(dir, "missing"), path); err == nil {
		t.Fatal("missing staging file passed")
	}
	if b, err := os.ReadFile(path); err != nil || string(b) != "old" {
		t.Fatalf("old state damaged: %q %v", b, err)
	}
	f, err := os.CreateTemp(dir, ".stage-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString("new"); err != nil {
		t.Fatal(err)
	}
	if err := SyncFile(f); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := ReplaceSyncedFile(f.Name(), filepath.Join(t.TempDir(), "elsewhere")); err == nil {
		t.Fatal("non-sibling staging allowed")
	}
	if err := ReplaceSyncedFile(f.Name(), path); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(path); err != nil || string(b) != "new" {
		t.Fatalf("replacement missing: %q %v", b, err)
	}
	if _, err := os.Stat(f.Name()); !os.IsNotExist(err) {
		t.Fatal("staging entry survived replacement")
	}
}
