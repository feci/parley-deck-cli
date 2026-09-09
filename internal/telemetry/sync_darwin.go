package telemetry

import (
	"errors"
	"os"
	"syscall"
)

// File.Sync requests F_FULLFSYNC on Darwin. Some shared filesystems reject
// that device-specific operation with ENOTTY but support ordinary fsync.
// Require the latter to succeed; never ignore an I/O or persistence failure.
func syncRecord(file *os.File) error {
	err := file.Sync()
	if errors.Is(err, syscall.ENOTTY) {
		return syscall.Fsync(int(file.Fd()))
	}
	return err
}
