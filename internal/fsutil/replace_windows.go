package fsutil

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// ReplaceSyncedFile replaces path with an already-synchronized sibling file.
// Windows directory handles opened for reading cannot supply a FlushFileBuffers
// barrier. Request write-through on the replacement operation instead. Do not
// enable COPY_ALLOWED: copying would expose a different publication contract.
func ReplaceSyncedFile(staged, path string) error {
	if filepath.Clean(filepath.Dir(staged)) != filepath.Clean(filepath.Dir(path)) {
		return fmt.Errorf("replacement must use a sibling staging file")
	}
	from, err := windows.UTF16PtrFromString(staged)
	if err != nil {
		return err
	}
	to, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(from, to, windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH)
}
