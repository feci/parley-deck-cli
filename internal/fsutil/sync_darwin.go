package fsutil

import (
	"errors"
	"os"
	"syscall"
)

// SyncFile requests F_FULLFSYNC on Darwin. Some shared filesystems reject
// that device-specific operation with ENOTTY but support ordinary fsync.
// Require the latter to succeed; never ignore an I/O or persistence failure.
func SyncFile(file *os.File) error {
	err := file.Sync()
	if errors.Is(err, syscall.ENOTTY) {
		return syscall.Fsync(int(file.Fd()))
	}
	return err
}
