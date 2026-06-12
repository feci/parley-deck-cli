// Package fsutil holds small filesystem helpers that harden standard-library
// operations against the weak cache coherence of networked or virtualized mounts
// (virtio-fs, NFS, SMB), where a directory created on the host can be momentarily
// invisible to a guest process whose dentry/attribute cache is stale.
package fsutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// Seams, overridable in tests only.
var (
	mkdirAll = os.MkdirAll
	stat     = os.Stat
	sleep    = time.Sleep
)

// retryDelays are the sleeps before each retry after the initial attempt. The first
// retry is immediate (the retry itself often forces a cache revalidation); later retries
// give a weakly-coherent attribute cache time to settle. The window is sized to outlast
// a virtio-fs / NFS attribute+entry cache timeout (typically ~1s): 8 total MkdirAll
// attempts with a worst-case added latency of ~1.9s — paid ONLY on the error path, never
// on a healthy first success. A spurious stale-cache failure is recoverable within this
// window; a genuine error simply costs ~1.9s before surfacing (permission errors fail
// fast and never wait).
var retryDelays = []time.Duration{
	0, // immediate first retry, no sleep
	15 * time.Millisecond,
	35 * time.Millisecond,
	100 * time.Millisecond,
	250 * time.Millisecond,
	500 * time.Millisecond,
	1000 * time.Millisecond,
}

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

// AppendLine appends one line to a ledger file as a single O_APPEND write
// (mirroring the event store's append discipline). Because O_APPEND atomicity
// is weak on weakly-coherent mounts (virtio-fs), cross-process writers are
// serialized through an adjacent mkdir claim — mkdir IS atomic even there
// (runner-hardening-kindly D10). A wedged claim degrades to an unlocked append
// after the bounded wait rather than losing the record.
func AppendLine(path string, line []byte) error {
	if err := MkdirAllResilient(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	claim := path + ".lock"
	acquired := false
	for attempt := 0; attempt < 50; attempt++ { // ~5s bound; a ledger append is tiny
		if err := os.Mkdir(claim, 0o755); err == nil {
			acquired = true
			break
		}
		sleep(100 * time.Millisecond)
	}
	if acquired {
		defer os.Remove(claim)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	if len(line) == 0 || line[len(line)-1] != '\n' {
		line = append(append([]byte(nil), line...), '\n')
	}
	_, err = file.Write(line)
	return err
}
