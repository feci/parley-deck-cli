//go:build windows

package budget

import (
	"errors"
	"golang.org/x/sys/windows"
	"os"
)

// Keep the mandatory Windows lock beyond every bounded identity read. Reading
// the locked byte through the second handle would otherwise fail before the
// exclusion probe can run. Windows permits a lock range beyond current EOF.
const budgetLockOffset = 1 << 20

func tryLock(f *os.File) (bool, error) {
	err := windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &windows.Overlapped{Offset: budgetLockOffset})
	if errors.Is(err, windows.ERROR_LOCK_VIOLATION) {
		return false, nil
	}
	return err == nil, err
}
func unlock(f *os.File) {
	_ = windows.UnlockFileEx(windows.Handle(f.Fd()), 0, 1, 0, &windows.Overlapped{Offset: budgetLockOffset})
}

func publishExclusive(staged, path string) error {
	from, err := windows.UTF16PtrFromString(staged)
	if err != nil {
		return err
	}
	to, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	// No REPLACE_EXISTING: a competing complete file wins and must match.
	return windows.MoveFileEx(from, to, windows.MOVEFILE_WRITE_THROUGH)
}
