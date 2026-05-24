//go:build !windows

package acp

import (
	"os/exec"
	"syscall"
)

// setSysProcAttr puts the child into its own process group via Setsid.
// AionUi does the equivalent with detached:true to prevent agents that
// touch /dev/tty (e.g. CodeBuddy) from suspending the parent on SIGTTOU.
func setSysProcAttr(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setsid = true
}
