//go:build windows

package fsutil

import "os"

// SyncFile flushes a FILE handle. Directory handles are refused with the
// named error instead of reaching FlushFileBuffers on a read-only handle
// (raw Access denied — the hosted row-25 family): on Windows a directory
// entry-durability barrier exists only through SyncDir, which refuses
// fail-closed (AC-DUR-1: no directory handle reaches SyncFile).
func SyncFile(file *os.File) error {
	if info, err := file.Stat(); err == nil && info.IsDir() {
		return ErrDirEntryDurabilityUnsupported
	}
	return file.Sync()
}
