package budget

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// The lock inode is permanent: never unlink it or atomic-replace it. Use a
// host-local cache because some shared mounts report successful flock without
// exclusion. This coordinates runtimes on one host, not distributed writers.
// Kernel ownership is released on exit, without PID files or stale-lock races.
func lock(ctx context.Context, path string) (func(), error) {
	canonical, err := filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	canonical, err = filepath.Abs(canonical)
	if err != nil {
		return nil, err
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		canonical = strings.ToLower(canonical)
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	localDir := filepath.Join(cache, "parley", "budget-locks")
	if err := os.MkdirAll(localDir, 0o700); err != nil {
		return nil, err
	}
	path = filepath.Join(localDir, key(canonical)+".lock")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(path)
	opened, statErr := f.Stat()
	if err != nil || statErr != nil || !info.Mode().IsRegular() || !os.SameFile(info, opened) {
		f.Close()
		return nil, errors.New("budget lock must be a stable regular file")
	}
	for {
		if err := ctx.Err(); err != nil {
			f.Close()
			return nil, err
		}
		held, err := tryLock(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		if held {
			// Verify the actual lock filesystem, rather than trusting a successful
			// syscall on a filesystem that implements it as a no-op.
			probe, err := os.OpenFile(path, os.O_RDWR, 0o600)
			if err != nil {
				unlock(f)
				f.Close()
				return nil, err
			}
			second, probeErr := tryLock(probe)
			if second {
				unlock(probe)
			}
			probe.Close()
			if probeErr != nil || second {
				unlock(f)
				f.Close()
				return nil, errors.New("budget lock filesystem does not provide verified exclusion")
			}
			return func() { unlock(f); f.Close() }, nil
		}
		timer := time.NewTimer(20 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			f.Close()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}
