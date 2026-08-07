//go:build linux

package app

import (
	"os"

	"golang.org/x/sys/unix"
)

const hasTTYSupported = true

func platformHasTTY() bool {
	for _, f := range []*os.File{os.Stdin, os.Stdout} {
		if _, err := unix.IoctlGetTermios(int(f.Fd()), unix.TCGETS); err == nil {
			return true
		}
	}
	return false
}
