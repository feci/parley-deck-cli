//go:build !windows

package procctl

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// AttachTerminal gives a terminal child its own foreground process group.
// The caller must restore the previous foreground group after Wait, including
// failed Start. Redirected descriptors need only an isolated process group.
func AttachTerminal(cmd *exec.Cmd, input *os.File) (restore func() error, err error) {
	group, err := unix.IoctlGetInt(int(input.Fd()), unix.TIOCGPGRP)
	if errors.Is(err, syscall.ENOTTY) {
		SetNewProcessGroup(cmd)
		return func() error { return nil }, nil
	}
	if err != nil {
		return nil, err
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Foreground: true, Ctty: int(input.Fd())}
	return func() error {
		// os/exec's child setup masks SIGTTOU around the foreground ioctl.
		// Restore through that setup instead of temporarily changing the
		// parent's process-wide signal handlers or stopping its background group.
		// This constant no-op neither interprets task data nor runs an agent.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		restore := exec.CommandContext(ctx, "/bin/sh", "-c", "exit 0")
		restore.Stdin = input
		restore.SysProcAttr = &syscall.SysProcAttr{Foreground: true, Pgid: group, Ctty: int(input.Fd())}
		return restore.Run()
	}, nil
}
