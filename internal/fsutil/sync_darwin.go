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
		raw, controlErr := file.SyscallConn()
		if controlErr != nil {
			return controlErr
		}
		var syncErr error
		if controlErr = raw.Control(func(fd uintptr) { syncErr = syscall.Fsync(int(fd)) }); controlErr != nil {
			return controlErr
		}
		return syncErr
	}
	return err
}
