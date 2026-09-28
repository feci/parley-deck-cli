//go:build !windows

package fsutil

import (
	"fmt"
	"os"
	"runtime"
)

// PublishFileDurable and PublishDirDurable implement the §B plain-row
// conversions (A1 det-stage; B1/B2 stage-Mkdir + WT-move) for Windows only.
// POSIX keeps the byte-identical inline implementations at the call sites and
// never invokes these; reaching one here is a caller bug, so both fail closed
// with a named error rather than silently no-op.

func PublishFileDurable(path string, data []byte, perm os.FileMode) error {
	return fmt.Errorf("fsutil.PublishFileDurable is the Windows-only §B A1 dormant conversion; reached on %s — caller bug", runtime.GOOS)
}

func PublishDirDurable(path string, perm os.FileMode) error {
	return fmt.Errorf("fsutil.PublishDirDurable is the Windows-only §B B1/B2 dormant conversion; reached on %s — caller bug", runtime.GOOS)
}
