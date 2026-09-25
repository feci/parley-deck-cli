package fsacl

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// The refusal string is product UX (§D.1: sentence-stable, the hosted 71x
// signature). Unix must stay byte-identical to the previous inline guard; the
// constant is shared by both platform implementations.
func TestRefusalStringIsStable(t *testing.T) {
	if got := ErrNotPrivate.Error(); got != "snapshot store must be a private real directory" {
		t.Fatalf("refusal string drifted: %q", got)
	}
}

func TestEnsurePrivateStoreCreatesAndPasses(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	if err := EnsurePrivateStore(dir); err != nil {
		t.Fatalf("EnsurePrivateStore: %v", err)
	}
	info, err := os.Lstat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("created store perm=%#o, group/other bits set", info.Mode().Perm())
	}
	if err := VerifyPrivateStore(dir); err != nil {
		t.Fatalf("VerifyPrivateStore after own creation: %v", err)
	}
}

func TestVerifyRefusesPermissivePreExistingStore(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	err := VerifyPrivateStore(dir)
	if err == nil {
		t.Fatal("permissive store accepted")
	}
	if !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("refusal does not wrap ErrNotPrivate: %v", err)
	}
}

func TestVerifyRefusesSymlinkAndNonDirectory(t *testing.T) {
	base := t.TempDir()
	real := filepath.Join(base, "real")
	if err := os.Mkdir(real, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "link")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}
	if err := VerifyPrivateStore(link); !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("symlink store: err=%v", err)
	}
	file := filepath.Join(base, "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := VerifyPrivateStore(file); !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("non-directory store: err=%v", err)
	}
}

// DenyRead is the §D.9 shared helper; on Unix it must make the file
// unreadable so the unreadable-marker classes stay testable everywhere.
func TestDenyReadMakesFileUnreadable(t *testing.T) {
	file := filepath.Join(t.TempDir(), "marker")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := DenyRead(file); err != nil {
		t.Fatalf("DenyRead: %v", err)
	}
	if _, err := os.ReadFile(file); err == nil {
		t.Fatal("file still readable after DenyRead")
	}
}
