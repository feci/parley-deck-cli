//go:build windows

package budget

import (
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// openLockFile opens a budget lock file with the §D.6 bounded
// ERROR_SHARING_VIOLATION retry. Classification: every handle this process
// takes on the lock file (kernel lock, probe, identity read) opens with Go's
// FILE_SHARE_READ|FILE_SHARE_WRITE and read/write access, so our own handles
// can never conflict with this open — a sharing violation here is by
// construction a FOREIGN holder (another process: antivirus, indexer,
// search). Only that class is retried, bounded to ~250ms; ERROR_ACCESS_DENIED
// (delete-pending, permissions) is a different class and never retried. On
// retry success the caller's identity/exclusion verification chain re-runs
// unchanged (the open is all that is retried); on exhaustion the failure is
// loud, naming the retry budget.
func openLockFile(path string) (*os.File, error) {
	const budget = 250 * time.Millisecond
	deadline := time.Now().Add(budget)
	for {
		f, err := os.OpenFile(path, os.O_RDWR, 0o600)
		if err == nil {
			return f, nil
		}
		if !errors.Is(err, windows.ERROR_SHARING_VIOLATION) {
			return nil, err
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("open %s: a foreign process (e.g. antivirus or indexer) held the lock file through the %s sharing-violation retry budget; refusing rather than extending the wait: %w", path, budget, err)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

// openLockFileRead is the read-only variant for the identity read: same
// §D.6 bounded retry and same structural classification — the process's own
// handles on the lock file (kernel lock, probe: read/write access, share
// read|write) are compatible with a read-only open, so a sharing violation
// here is by construction foreign. Hosted evidence (36438474793): the
// identity read was the open that surfaced the raw unretried failures.
func openLockFileRead(path string) (*os.File, error) {
	const budget = 250 * time.Millisecond
	deadline := time.Now().Add(budget)
	for {
		f, err := os.Open(path)
		if err == nil {
			return f, nil
		}
		if !errors.Is(err, windows.ERROR_SHARING_VIOLATION) {
			return nil, err
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("open %s: a foreign process (e.g. antivirus or indexer) held the lock file through the %s sharing-violation retry budget; refusing rather than extending the wait: %w", path, budget, err)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

// openOriginFileRead opens the lock-ORIGIN record with the same bounded
// §D.6 retry, with an honestly narrower classification than the lock-file
// opens: hosted evidence (36440560978, invocation 11) showed the origin read
// hitting sharing violations during CONCURRENT ORIGIN PUBLICATION by a
// sibling parley process (the transient delete-pending/rename window of the
// bootstrap writer), not only foreign scanners. Both holder classes are
// transient by nature; the retry is bounded and exhaustion is loud, and the
// origin/token verification chain re-runs unchanged after recovery. The
// writer-side discipline question (origin publication shape,
// close-before-rename) is recorded as follow-up for review.
func openOriginFileRead(path string) (*os.File, error) {
	const budget = 250 * time.Millisecond
	deadline := time.Now().Add(budget)
	for {
		f, err := os.Open(path)
		if err == nil {
			return f, nil
		}
		if !errors.Is(err, windows.ERROR_SHARING_VIOLATION) {
			return nil, err
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("open %s: the lock-origin record was held (a sibling parley process publishing it, or a foreign scanner) through the %s sharing-violation retry budget; refusing rather than extending the wait: %w", path, budget, err)
		}
		time.Sleep(25 * time.Millisecond)
	}
}
