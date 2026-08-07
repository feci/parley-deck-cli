//go:build linux

package app

import "golang.org/x/sys/unix"

const termiosGet = unix.TCGETS
