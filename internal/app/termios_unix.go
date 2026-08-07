//go:build darwin || freebsd || netbsd || openbsd

package app

import (
	"os"

	"golang.org/x/sys/unix"
)

const hasTTYSupported = true

// platformHasTTY asks the kernel for the terminal attributes of the fd. It succeeds only for a
// real tty — unlike a character-device test, which also accepts /dev/null.
func platformHasTTY() bool {
	for _, f := range []*os.File{os.Stdin, os.Stdout} {
		if _, err := unix.IoctlGetTermios(int(f.Fd()), unix.TIOCGETA); err == nil {
			return true
		}
	}
	return false
}
