//go:build windows

package procctl

import (
	"os"
	"os/exec"
)

// Windows children inherit the console; the live process-tree kill remains
// owned by KillGroup. No Unix foreground process-group transition is needed.
func AttachTerminal(cmd *exec.Cmd, _ *os.File) (func() error, error) {
	SetNewProcessGroup(cmd)
	return func() error { return nil }, nil
}
