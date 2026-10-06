//go:build !windows

package pidlease

import (
	"errors"
	"syscall"
)

func definitelyDead(pid int) bool { return pid > 0 && errors.Is(syscall.Kill(pid, 0), syscall.ESRCH) }

// ProcessGroup records the group at launch; a missing identity cannot be recovered.
func ProcessGroup(pid int) int {
	pg, err := syscall.Getpgid(pid)
	if err != nil {
		return 0
	}
	return pg
}

// GroupStopped also checks descendants still in the supervised process group.
func GroupStopped(pgid int) bool {
	return pgid > 0 && pgid <= 2147483647 && errors.Is(syscall.Kill(-pgid, 0), syscall.ESRCH)
}
