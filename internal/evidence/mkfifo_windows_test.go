//go:build windows

package evidence

import "errors"

// mkfifo is the Windows arm of the unsupported-entry probe fixture: FIFOs do
// not exist, so the probe reports that and the caller falls back to the
// file's existing makeUnsupportedSocket path (AC-BLD-1, FINAL §D.8).
func mkfifo(path string, mode uint32) error {
	_ = mode
	return errors.New("fifo not supported on windows: " + path)
}
