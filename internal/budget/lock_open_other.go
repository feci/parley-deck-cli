//go:build !windows

package budget

import "os"

// openLockFile is the plain open on POSIX: flock semantics make sharing
// violations impossible, and the byte-range kernel lock serializes holders.
func openLockFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, 0o600)
}

func openLockFileRead(path string) (*os.File, error) {
	return os.Open(path)
}
