package budget

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"parley-deck-cli/internal/fsutil"
)

const resourceGuardWitness = "guard-established"

// AcquireResourceGuard reuses the ledger's pinned, verified host-local kernel
// lock for a synchronization-only resource. dir must be a permanent, exclusive
// metadata directory for that resource. No ledger or charge is created. The
// same origin/cache migration restrictions as budget locks apply; removing the
// origin or established local lock while writers exist is unsupported.
func AcquireResourceGuard(ctx context.Context, dir string) (func(), error) {
	return acquireResourceGuard(ctx, dir, publishExclusive)
}

func acquireResourceGuard(ctx context.Context, dir string, publish func(string, string) error) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := fsutil.MkdirAllResilient(dir, 0700); err != nil {
		return nil, err
	}
	return lockWithReady(ctx, filepath.Join(dir, "resource"), tryLock, unlock, func(canonical string) error {
		if err := establishResourceGuard(canonical, publish); err != nil {
			return fmt.Errorf("resource guard continuity: %w", err)
		}
		return nil
	})
}

// A synchronization-only guard has no ledger whose existence can distinguish
// first use from lost lock state. Publish its exact origin before returning the
// first permission. Losing the origin/cache after that point must fail closed,
// even when an old unlinked lock descriptor is still held in another process.
// Erasing both the witness and the origin is outside this continuity contract.
func establishResourceGuard(dir string, publish func(string, string) error) error {
	origin, err := readLockOrigin(filepath.Join(dir, "lock-origin"))
	if err != nil {
		return err
	}
	path := filepath.Join(dir, resourceGuardWitness)
	check := func() error {
		prior, err := readLockOrigin(path)
		if err != nil {
			return err
		}
		if !bytes.Equal(prior, origin) {
			return errors.New("established resource guard differs from its pinned origin; preserve both records")
		}
		current, err := readLockOrigin(filepath.Join(dir, "lock-origin"))
		if err != nil {
			return err
		}
		if !bytes.Equal(current, origin) {
			return errors.New("resource guard origin changed during continuity publication")
		}
		return nil
	}
	if err := check(); !os.IsNotExist(err) {
		return err
	}
	f, err := os.CreateTemp(dir, ".guard-established-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.Write(origin); err != nil {
		return err
	}
	if err := fsutil.SyncFile(f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := publish(f.Name(), path); err != nil && !os.IsExist(err) {
		return err
	}
	return check()
}
