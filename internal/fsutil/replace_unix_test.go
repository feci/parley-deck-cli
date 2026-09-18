//go:build !windows

package fsutil

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestReplacementDirectoryBarrierFailurePreservesConservativeState(t *testing.T) {
	for _, failure := range []error{syscall.EINVAL, syscall.EOPNOTSUPP, syscall.EIO} {
		t.Run(failure.Error(), func(t *testing.T) {
			root := t.TempDir()
			target := filepath.Join(root, "cursor")
			staged := filepath.Join(root, "staged")
			if err := os.WriteFile(target, []byte("uncharged"), 0o600); err != nil {
				t.Fatal(err)
			}
			f, err := os.OpenFile(staged, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = f.WriteString("charged"); err != nil {
				t.Fatal(err)
			}
			if err = SyncFile(f); err != nil {
				t.Fatal(err)
			}
			if err = f.Close(); err != nil {
				t.Fatal(err)
			}
			err = replaceSyncedFile(staged, target, func(dir *os.File) error {
				info, e := dir.Stat()
				if e != nil || !info.IsDir() {
					t.Fatal("barrier is not the parent directory")
				}
				return failure
			})
			if !errors.Is(err, failure) {
				t.Fatalf("persistence failure hidden: %v", err)
			}
			body, err := os.ReadFile(target)
			if err != nil || string(body) != "charged" {
				t.Fatalf("conservative state lost: %q %v", body, err)
			}
		})
	}
}
