//go:build !windows

package evidence

import "syscall"

// mkfifo is the Unix arm of the unsupported-entry probe fixture: FIFOs are
// the natural first choice, and this host arm creates one directly.
func mkfifo(path string, mode uint32) error {
	return syscall.Mkfifo(path, mode)
}
