//go:build windows

package fsutil

import (
	"fmt"
	"io"
	"os"
)

// SyncFile flushes a FILE handle. Directory handles are refused with the
// named error instead of reaching FlushFileBuffers on a read-only handle
// (raw Access denied — the hosted row-25 family): on Windows a directory
// entry-durability barrier exists only through SyncDir, which refuses
// fail-closed (AC-DUR-1: no directory handle reaches SyncFile).
//
// A read-only FILE handle cannot supply FlushFileBuffers either (it needs
// write access — the D4a apply-path finding): rather than surfacing the
// raw errno, the caller's durability intent fails with the SAME named
// refusal — FINAL §C requires every refusal to be named and actionable,
// and §B's D4a/D4b follow A3/B5 into refusal. The write-access check is
// the handle's own stat: a handle our writer opened for writing passes.
func SyncFile(file *os.File) error {
	if info, err := file.Stat(); err == nil {
		if info.IsDir() {
			return ErrDirEntryDurabilityUnsupported
		}
	}
	if writable(file) {
		return file.Sync()
	}
	return fmt.Errorf("%w: syncing a read-only file handle cannot supply the write-through durability barrier", ErrDirEntryDurabilityUnsupported)
}

// writable reports whether the handle was opened with write access, via a
// zero-byte write probe at the current offset (harmless on every filesystem
// the product writes) — the most portable check without NtQueryInformationFile.
func writable(file *os.File) bool {
	off, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return false
	}
	_, werr := file.Write(nil)
	_, _ = file.Seek(off, io.SeekStart)
	return werr == nil
}
