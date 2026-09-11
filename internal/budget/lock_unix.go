//go:build !windows

package budget

import (
	"errors"
	"golang.org/x/sys/unix"
	"os"
	"parley-deck-cli/internal/fsutil"
	"path/filepath"
)

func tryLock(f *os.File) (bool, error) {
	err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if errors.Is(err, unix.EWOULDBLOCK) {
		return false, nil
	}
	return err == nil, err
}
func unlock(f *os.File) { _ = unix.Flock(int(f.Fd()), unix.LOCK_UN) }

func publishOrigin(staged, path string) error {
	if err := os.Link(staged, path); err != nil {
		return err
	}
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return fsutil.SyncFile(dir)
}
