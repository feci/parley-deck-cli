//go:build !windows

package fsacl

import (
	"os"
	"path/filepath"
	"testing"
)

// Unix mechanics pin (AC-PRIV-6 byte-identical behaviour): own creation
// yields a real 0700 directory.
func TestUnixCreatedStoreHasOwnerOnlyModeBits(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	if err := EnsurePrivateStore(dir); err != nil {
		t.Fatalf("EnsurePrivateStore: %v", err)
	}
	info, err := os.Lstat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("created store mode=%v, want real directory with no group/other bits", info.Mode())
	}
}
