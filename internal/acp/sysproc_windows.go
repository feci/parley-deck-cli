//go:build windows

package acp

import "os/exec"

// setSysProcAttr is a no-op on Windows; the /dev/tty issue does not apply.
func setSysProcAttr(_ *exec.Cmd) {}
