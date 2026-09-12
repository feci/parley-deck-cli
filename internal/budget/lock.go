package budget

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
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
// This is same-origin coordination, not distributed cross-host locking. A
// same-host cache move uses journaled migration with both kernel locks held.
// Never delete a live lock inode or its origin file.
// Kernel ownership is released on exit, without PID files or stale-lock races.
func lock(ctx context.Context, path string) (func(), error) {
	return lockWithOps(ctx, path, tryLock, unlock)
}

func lockWithOps(ctx context.Context, path string, take func(*os.File) (bool, error), drop func(*os.File)) (func(), error) {
	return lockWithReady(ctx, path, take, drop, nil)
}

func lockWithReady(ctx context.Context, path string, take func(*os.File) (bool, error), drop func(*os.File), ready func(string) error) (func(), error) {
	canonical, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	canonical, err = filepath.EvalSymlinks(canonical)
	if err != nil {
		return nil, err
	}
	prior, readErr := readLockOrigin(filepath.Join(canonical, "lock-origin"))
	if readErr != nil && !os.IsNotExist(readErr) {
		return nil, readErr
	}
	if bytes.HasPrefix(prior, []byte("parley-budget-lock/v3\n")) {
		origin, err := resolveMigratedOrigin(canonical, prior)
		if err != nil {
			return nil, err
		}
		return acquirePinnedKernelLock(ctx, origin.Path, origin.Token, take, drop, func() error {
			if err := verifyExactLockOrigin(canonical, prior); err != nil {
				return err
			}
			_, err := resolveMigratedOrigin(canonical, prior)
			return err
		}, func() error {
			if ready != nil {
				return ready(canonical)
			}
			return nil
		})
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
	_, originErr := os.Lstat(filepath.Join(canonical, "lock-origin"))
	newOrigin := os.IsNotExist(originErr)
	if originErr != nil && !newOrigin {
		return nil, originErr
	}
	if newOrigin {
		for _, witness := range []string{"ledger-established", originMigrationDirectory} {
			if _, err := os.Lstat(filepath.Join(canonical, witness)); err == nil {
				return nil, errors.New("established budget history has no lock origin; preserve retained history; refusing recreation")
			} else if !os.IsNotExist(err) {
				return nil, err
			}
		}
		if _, err := os.Lstat(filepath.Join(canonical, resourceGuardWitness)); err == nil {
			return nil, errors.New("established resource guard has no lock origin; preserve its witness and stop writers; refusing recreation")
		} else if !os.IsNotExist(err) {
			return nil, err
		}
		if _, err := os.Lstat(filepath.Join(canonical, "ledger.json")); err == nil {
			return nil, errors.New("existing budget ledger has no lock origin; preserve charges and stop writers: migration requires the original accessible identity")
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	if !newOrigin {
		if err := checkLockOriginLocation(canonical, path); err != nil {
			return nil, err
		}
	}
	token, err := lockIdentity(path, newOrigin)
	if err != nil {
		return nil, err
	}
	if err := pinLockOrigin(canonical, path, token); err != nil {
		return nil, err
	}
	host, err := os.Hostname()
	if err != nil || host == "" {
		return nil, errors.New("cannot identify budget lock host")
	}
	expected := []byte("parley-budget-lock/v2\n" + host + "\n" + path + "\n" + token + "\n")
	return acquirePinnedKernelLock(ctx, path, token, take, drop, func() error {
		return verifyExactLockOrigin(canonical, expected)
	}, func() error {
		if ready != nil {
			return ready(canonical)
		}
		return nil
	})
}

// The check closure captures authority BEFORE any wait. A waiter cannot adopt
// an origin observed only after acquisition, including a completed migration.
func acquirePinnedKernelLock(ctx context.Context, path, token string, take func(*os.File) (bool, error), drop func(*os.File), check, ready func() error) (func(), error) {
	// Never recreate an established inode, including during migration recovery.
	f, err := os.OpenFile(path, os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	if err := verifyLockIdentity(f, path, token); err != nil {
		f.Close()
		return nil, err
	}
	waited := false
	interrupted := func(err error) error {
		if waited {
			return fmt.Errorf("%w at %s: %w", ErrLockContention, path, err)
		}
		return err
	}
	for {
		if err := ctx.Err(); err != nil {
			f.Close()
			return nil, interrupted(err)
		}
		held, err := take(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		if held {
			// Rebind the held descriptor to the pinned token AND current inode.
			// Bootstrap racers may have read identity before a cache replacement;
			// matching only a pathname would authorize a different lock here.
			if err := verifyLockIdentity(f, path, token); err != nil {
				drop(f)
				f.Close()
				return nil, err
			}
			if err := check(); err != nil {
				drop(f)
				f.Close()
				return nil, err
			}
			// Verify the actual lock filesystem, rather than trusting a successful
			// syscall on a filesystem that implements it as a no-op.
			probe, err := os.OpenFile(path, os.O_RDWR, 0o600)
			if err != nil {
				drop(f)
				f.Close()
				return nil, err
			}
			if err := verifyLockIdentity(probe, path, token); err != nil {
				probe.Close()
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
			if err := verifyLockIdentity(f, path, token); err != nil {
				drop(f)
				f.Close()
				return nil, err
			}
			if err := check(); err != nil {
				drop(f)
				f.Close()
				return nil, err
			}
			if ready != nil {
				if err := ready(); err != nil {
					drop(f)
					f.Close()
					return nil, err
				}
				// Continuity publication may perform filesystem work. It cannot
				// switch this acquisition to a newly observed origin or inode.
				if err := verifyLockIdentity(f, path, token); err != nil {
					drop(f)
					f.Close()
					return nil, err
				}
				if err := check(); err != nil {
					drop(f)
					f.Close()
					return nil, err
				}
			}
			if err := ctx.Err(); err != nil {
				drop(f)
				f.Close()
				return nil, interrupted(err)
			}
			var once sync.Once
			return func() { once.Do(func() { drop(f); f.Close() }) }, nil
		}
		waited = true
		timer := time.NewTimer(20 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			f.Close()
			return nil, interrupted(ctx.Err())
		case <-timer.C:
		}
	}
}

// An origin read before waiting is not authority after acquiring a kernel lock.
// In particular, never use pinLockOrigin here: a missing origin must not be
// recreated by an old waiter. This read-only check also follows the exclusion
// probe, which may have run after a concurrent origin change.
func verifyPinnedLockOrigin(dir, lockPath, token string) error {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return errors.New("cannot identify budget lock host")
	}
	data, err := readLockOrigin(filepath.Join(dir, "lock-origin"))
	if err != nil {
		return fmt.Errorf("budget lock origin unavailable after acquisition: %w", err)
	}
	expected := []byte("parley-budget-lock/v2\n" + host + "\n" + lockPath + "\n" + token + "\n")
	if !bytes.Equal(data, expected) {
		return errors.New("budget lock origin changed before exclusion was established")
	}
	return nil
}

// Verify the exact captured version, destination and journal reference.
func verifyExactLockOrigin(dir string, expected []byte) error {
	actual, err := readLockOrigin(filepath.Join(dir, "lock-origin"))
	if err != nil {
		return fmt.Errorf("budget lock origin unavailable after acquisition: %w", err)
	}
	if !bytes.Equal(actual, expected) {
		return errors.New("budget lock origin changed before exclusion was established")
	}
	return nil
}

// Read through the actual descriptor that will hold the kernel lock. ReadAt
// leaves its position unchanged and bounds the read even for a replaced file.
func verifyLockIdentity(f *os.File, path, token string) error {
	info, err := os.Lstat(path)
	opened, statErr := f.Stat()
	if err != nil || statErr != nil || !info.Mode().IsRegular() || !os.SameFile(info, opened) || opened.Size() != 65 {
		return errors.New("budget lock must be a stable regular file with its pinned identity")
	}
	var data [66]byte
	n, err := f.ReadAt(data[:], 0)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if n != 65 || string(data[:n]) != token+"\n" {
		return errors.New("budget lock identity changed before exclusion was established")
	}
	return nil
}

// Publish the entire origin atomically before any writer acquires a local lock.
// A competing bootstrap can either use that exact origin or fail; it cannot
// silently initialize an independent lock and overwrite the other writer.
func pinLockOrigin(dir, lockPath, token string) error {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return errors.New("cannot identify budget lock host")
	}
	data := []byte("parley-budget-lock/v2\n" + host + "\n" + lockPath + "\n" + token + "\n")
	if len(data) > 16<<10 {
		return errors.New("budget lock origin exceeds limit")
	}
	path := filepath.Join(dir, "lock-origin")
	check := func() error {
		prior, err := readLockOrigin(path)
		if err != nil {
			return err
		}
		if !bytes.Equal(prior, data) {
			reason := "unsupported or malformed origin version"
			parts := strings.Split(string(prior), "\n")
			if len(parts) == 5 && parts[0] == "parley-budget-lock/v2" {
				switch {
				case parts[1] != host:
					reason = "hostname changed"
				case parts[2] != lockPath:
					reason = "cache path or ledger location changed"
				case parts[3] != token:
					reason = "local lock identity changed or was copied from another environment"
				}
			}
			return fmt.Errorf("budget lock origin mismatch at %s: %s; preserve ledger and use its original compatible environment; same-host cache relocation requires budget origin inspect/apply with the original accessible identity", path, reason)
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
	if err := publishExclusive(f.Name(), path); err != nil && !os.IsExist(err) {
		return err
	}
	return check()
}

// Each permanent local lock gets its own random identity, atomically published
// before origin pinning. Hostnames and cache paths alone can repeat. A copied
// identity is not machine authentication; distributed/cloned writers remain
// outside the supported boundary. Missing established identities fail closed.
func lockIdentity(path string, create bool) (string, error) {
	read := func() (string, error) {
		info, err := os.Lstat(path)
		if err != nil {
			return "", err
		}
		if !info.Mode().IsRegular() || info.Size() != 65 {
			return "", errors.New("invalid local budget lock identity")
		}
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()
		opened, err := f.Stat()
		if err != nil || !os.SameFile(info, opened) {
			return "", errors.New("local budget lock changed during identity read")
		}
		data, err := io.ReadAll(io.LimitReader(f, 66))
		if err != nil {
			return "", err
		}
		if len(data) != 65 || data[64] != '\n' {
			return "", errors.New("invalid local budget lock identity")
		}
		decoded, err := hex.DecodeString(string(data[:64]))
		if err != nil || len(decoded) != 32 {
			return "", errors.New("invalid local budget lock identity")
		}
		return string(data[:64]), nil
	}
	token, err := read()
	if !os.IsNotExist(err) {
		return token, err
	}
	if !create {
		return "", fmt.Errorf("established local budget lock is missing at %s; refusing recreation", path)
	}
	var nonce [32]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".lock-identity-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.WriteString(hex.EncodeToString(nonce[:]) + "\n"); err != nil {
		return "", err
	}
	if err := fsutil.SyncFile(f); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	if err := publishExclusive(f.Name(), path); err != nil && !os.IsExist(err) {
		return "", err
	}
	return read()
}

// Diagnose the immutable origin before consulting a possibly absent local
// identity. This is a read-only header check; token equality is still enforced
// by pinLockOrigin and the held descriptor checks before any permission.
func checkLockOriginLocation(dir, lockPath string) error {
	path := filepath.Join(dir, "lock-origin")
	data, err := readLockOrigin(path)
	if err != nil {
		return err
	}
	parts := strings.Split(string(data), "\n")
	reason := ""
	host, err := os.Hostname()
	if err != nil || host == "" {
		return errors.New("cannot identify budget lock host")
	}
	switch {
	case len(parts) != 5 || parts[0] != "parley-budget-lock/v2" || parts[4] != "":
		reason = "unsupported or malformed origin version"
	case parts[1] != host:
		reason = "hostname changed"
	case parts[2] != lockPath:
		reason = "cache path or ledger location changed"
	}
	if reason != "" {
		return fmt.Errorf("budget lock origin mismatch at %s: %s; preserve charges; same-host cache relocation requires budget origin inspect/apply with the original accessible identity", path, reason)
	}
	return nil
}

func readLockOrigin(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<10 {
		return nil, errors.New("invalid budget lock origin")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return nil, errors.New("budget lock origin changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(f, (16<<10)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 16<<10 {
		return nil, errors.New("budget lock origin exceeds limit")
	}
	return data, nil
}
