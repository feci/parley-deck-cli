// Package fsutil holds small filesystem helpers that harden standard-library
// operations against the weak cache coherence of networked or virtualized mounts
// (virtio-fs, NFS, SMB), where a directory created on the host can be momentarily
// invisible to a guest process whose dentry/attribute cache is stale.
package fsutil

import (
	"errors"
	"io/fs"
	"os"
	"time"
)

// Seams, overridable in tests only.
var (
	mkdirAll = os.MkdirAll
	stat     = os.Stat
	sleep    = time.Sleep
)

// retryDelays are the sleeps before each retry after the initial attempt. The first
// retry is immediate (the retry itself often forces a cache revalidation); later
// retries give a weakly-coherent attribute cache a brief window to settle. This is 5
// total MkdirAll attempts with a worst-case added latency of 75ms — paid only on the
// error path, never on a healthy first success.
var retryDelays = []time.Duration{0 /* immediate first retry, no sleep */, 5 * time.Millisecond, 20 * time.Millisecond, 50 * time.Millisecond}

// MkdirAllResilient behaves like os.MkdirAll but tolerates spurious, transient failures
// from weakly-coherent filesystems whose dentry/attribute cache can be momentarily stale
// across processes. The healthy path is a single os.MkdirAll call with no Stat and no
// sleep. On error it returns nil whenever a fresh Stat shows the path is a directory,
// fails fast on permission errors, and otherwise retries a small bounded number of times
// before returning the last os.MkdirAll error.
//
// A pre-existing regular file at the target path is never reported as success: the
// success oracle is always a fresh Stat reporting a directory, so fs.ErrExist is not
// trusted blindly.
func MkdirAllResilient(path string, perm os.FileMode) error {
	err := mkdirAll(path, perm)
	if err == nil {
		return nil
	}
	for _, d := range retryDelays {
		if isDir(path) {
			return nil
		}
		if errors.Is(err, fs.ErrPermission) {
			return err
		}
		if d > 0 {
			sleep(d)
		}
		if err = mkdirAll(path, perm); err == nil {
			return nil
		}
	}
	if isDir(path) {
		return nil
	}
	return err
}

func isDir(path string) bool {
	info, e := stat(path)
	return e == nil && info.IsDir()
}
