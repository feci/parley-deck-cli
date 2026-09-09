//go:build !windows

package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReplaceSyncedFile replaces path with an already-synchronized sibling file.
// Callers must SyncFile and close staged before calling. The directory barrier
// is required; a failure after rename leaves the conservative replacement on
// disk and returns an error rather than claiming the transaction is durable.
func ReplaceSyncedFile(staged, path string) error {
	return replaceSyncedFile(staged, path, SyncFile)
}

func replaceSyncedFile(staged, path string, syncDir func(*os.File) error) error {
	if filepath.Clean(filepath.Dir(staged)) != filepath.Clean(filepath.Dir(path)) {
		return fmt.Errorf("replacement must use a sibling staging file")
	}
	if err := os.Rename(staged, path); err != nil {
		return err
	}
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return syncDir(dir)
}
