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
	return syncWithFallback(file.Sync, func() error {
		raw, err := file.SyscallConn()
		if err != nil {
			return err
		}
		var syncErr error
		if err = raw.Control(func(fd uintptr) { syncErr = syscall.Fsync(int(fd)) }); err != nil {
			return err
		}
		return syncErr
	})
}

func syncWithFallback(full, plain func() error) error {
	err := full()
	if errors.Is(err, syscall.ENOTTY) {
		return plain()
	}
	return err
}
