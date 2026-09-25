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

// TestEnsurePrivateStoreCreatesAndPasses is platform-neutral: own creation
// must pass verification everywhere. The mode-bit expression of privacy is
// pinned unix-side (fsacl_unix_test.go); on Windows the DACL round-trip is
// pinned in fsacl_windows_test.go (hosted fact 36172828285: dir perms
// synthesize 0777 there, so the perm assertion itself is Unix mechanics).
func TestEnsurePrivateStoreCreatesAndPasses(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	if err := EnsurePrivateStore(dir); err != nil {
		t.Fatalf("EnsurePrivateStore: %v", err)
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
