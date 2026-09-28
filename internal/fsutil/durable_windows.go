//go:build windows

package fsutil

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

// PublishFileDurable is the §B A1 conversion (dormant behind the app's
// trajectoryRuntime POSIX gate today): stage create O_EXCL + mandatory file
// fsync, then MoveFileEx(MOVEFILE_WRITE_THROUGH) WITHOUT
// MOVEFILE_REPLACE_EXISTING. The write-through move is the create-entry
// barrier — no directory fsync is attempted. A completed publication is never
// clobbered; concurrent publishers race on the stage O_EXCL exactly as they
// raced on the final-name O_EXCL. Every error path leaves the stage in place
// as the anti-replay blocker.
func PublishFileDurable(path string, data []byte, perm os.FileMode) error {
	stage := path + StageSuffix
	f, err := os.OpenFile(stage, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	if _, writeErr := f.Write(data); writeErr != nil {
		f.Close()
		return writeErr
	}
	if err = SyncFile(f); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	from, err := windows.UTF16PtrFromString(stage)
	if err != nil {
		return err
	}
	to, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(from, to, windows.MOVEFILE_WRITE_THROUGH)
}

// PublishDirDurable is the §B B1/B2 conversion (dormant behind the same
// gate): stage-Mkdir + MoveFileEx(MOVEFILE_WRITE_THROUGH) of the staged
// directory, with the EEXIST-tolerant recheck that preserves B1's idempotent
// ensure-exists semantics. A stale stage with no satisfying final is a loud,
// visible blocker — the deliberate interrupted-publication signal — never
// silently cleared.
func PublishDirDurable(path string, perm os.FileMode) error {
	if info, err := os.Lstat(path); err == nil {
		if info.IsDir() {
			return nil // EEXIST-tolerant recheck: the final already exists
		}
		return fmt.Errorf("%s exists and is not a directory", path)
	} else if !os.IsNotExist(err) {
		return err
	}
	stage := path + StageSuffix
	if err := os.Mkdir(stage, perm); err != nil {
		if os.IsExist(err) {
			// Another publisher holds the stage, or a prior attempt was
			// interrupted. Tolerate only when the final already satisfies the
			// ensure-exists contract; otherwise the stage is the visible
			// blocker and the operation fails loudly.
			if info, serr := os.Lstat(path); serr == nil && info.IsDir() {
				_ = os.RemoveAll(stage)
				return nil
			}
			return fmt.Errorf("staging directory %s blocks the publication; resolve it before retrying", stage)
		}
		return err
	}
	from, err := windows.UTF16PtrFromString(stage)
	if err != nil {
		return err
	}
	to, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	if err := windows.MoveFileEx(from, to, windows.MOVEFILE_WRITE_THROUGH); err != nil {
		if errors.Is(err, windows.ERROR_ALREADY_EXISTS) || errors.Is(err, windows.ERROR_FILE_EXISTS) {
			// A concurrent publisher completed the final name: tolerant recheck.
			if info, serr := os.Lstat(path); serr == nil && info.IsDir() {
				_ = os.Remove(stage)
				return nil
			}
		}
		return err // stage left in place: the visible blocker
	}
	return nil
}
