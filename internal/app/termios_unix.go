//go:build darwin || freebsd || netbsd || openbsd

package app

import "golang.org/x/sys/unix"

// termiosGet is the ioctl that reads terminal attributes. It succeeds only on a real tty, which
// is what makes it a usable controlling-terminal test; the constant differs between BSD-derived
// systems (TIOCGETA) and Linux (TCGETS).
const termiosGet = unix.TIOCGETA
