package budget

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/fsutil"
)

// The lock inode is permanent: never unlink it or atomic-replace it. Use a
// host-local cache because some shared mounts report successful flock without
// exclusion. An immutable origin pins the first host/cache path for the ledger;
// another cache environment or hostname refuses instead of taking a second lock.
// This is same-origin coordination, not distributed cross-host locking. Cache
// migration/removal requires quiescent operator maintenance; never delete a
// live lock inode or its origin file.
// Kernel ownership is released on exit, without PID files or stale-lock races.
func lock(ctx context.Context, path string) (func(), error) {
	return lockWithOps(ctx, path, tryLock, unlock)
}

func lockWithOps(ctx context.Context, path string, take func(*os.File) (bool, error), drop func(*os.File)) (func(), error) {
	canonical, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	canonical, err = filepath.EvalSymlinks(canonical)
	if err != nil {
		return nil, err
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	localDir := filepath.Join(cache, "parley", "budget-locks")
	if err := fsutil.MkdirAllResilient(localDir, 0o700); err != nil {
		return nil, err
	}
	localDir, err = filepath.Abs(localDir)
	if err != nil {
		return nil, err
	}
	localDir, err = filepath.EvalSymlinks(localDir)
	if err != nil {
		return nil, err
	}
	// Lowercase on every platform: Linux can mount a case-insensitive filesystem
	// too. Overlocking distinct case-sensitive paths is safe; underlocking aliases
	// is not. Resolve symlinks before this conservative key normalization.
	path = filepath.Join(localDir, key(strings.ToLower(canonical))+".lock")
	if err := pinLockOrigin(canonical, path); err != nil {
		return nil, err
	}
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
		held, err := take(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		if held {
			// Verify the actual lock filesystem, rather than trusting a successful
			// syscall on a filesystem that implements it as a no-op.
			probe, err := os.OpenFile(path, os.O_RDWR, 0o600)
			if err != nil {
				drop(f)
				f.Close()
				return nil, err
			}
			second, probeErr := take(probe)
			if second {
				drop(probe)
			}
			probe.Close()
			if probeErr != nil || second {
				drop(f)
				f.Close()
				return nil, fmt.Errorf("budget lock filesystem does not provide verified exclusion at %s: %v", path, probeErr)
			}
			var once sync.Once
			return func() { once.Do(func() { drop(f); f.Close() }) }, nil
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

// Publish the entire origin atomically before any writer acquires a local lock.
// A competing bootstrap can either use that exact origin or fail; it cannot
// silently initialize an independent lock and overwrite the other writer.
func pinLockOrigin(dir, lockPath string) error {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return errors.New("cannot identify budget lock host")
	}
	data := []byte("parley-budget-lock/v1\n" + host + "\n" + lockPath + "\n")
	if len(data) > 16<<10 {
		return errors.New("budget lock origin exceeds limit")
	}
	path := filepath.Join(dir, "lock-origin")
	check := func() error {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Size() > 16<<10 {
			return errors.New("invalid budget lock origin")
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		opened, err := f.Stat()
		if err != nil || !os.SameFile(info, opened) {
			return errors.New("budget lock origin changed during open")
		}
		prior, err := io.ReadAll(io.LimitReader(f, (16<<10)+1))
		if err != nil {
			return err
		}
		if !bytes.Equal(prior, data) {
			return fmt.Errorf("budget lock origin mismatch at %s: different host or cache path; preserve ledger and use its original environment", path)
		}
		return nil
	}
	if err := check(); !os.IsNotExist(err) {
		return err
	}
	f, err := os.CreateTemp(dir, ".lock-origin-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := fsutil.SyncFile(f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := publishOrigin(f.Name(), path); err != nil && !os.IsExist(err) {
		return err
	}
	return check()
}
