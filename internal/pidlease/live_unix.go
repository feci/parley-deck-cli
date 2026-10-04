//go:build !windows

package pidlease

import (
	"errors"
	"syscall"
)

func definitelyDead(pid int) bool { return pid > 0 && errors.Is(syscall.Kill(pid, 0), syscall.ESRCH) }
