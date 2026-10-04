// Package membership implements quota transition projections and the driving lease.
// Immutable authority lives in quota; read-only protocol consumers never call Repair.
package membership

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"parley-deck-cli/internal/pidlease"
	"path/filepath"
	"sync/atomic"
	"time"

	"parley-deck-cli/internal/protocol"
)

type leaseKey struct{}
type lease struct {
	path, run string
	active    atomic.Bool
	owner     *pidlease.Lease
}

// Acquire holds the idea lock for the entire driving operation. Nested calls may
// reuse only the very same active context lease; same-PID competing runs cannot.
func Acquire(ctx context.Context, ideaDir, run string) (context.Context, func(), error) {
	v, err := protocol.InspectQuota(ideaDir)
	if os.IsNotExist(err) && v.History == nil {
		return ctx, func() {}, nil
	}
	if err != nil {
		root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
		return ctx, nil, IntegrityBlock(root, ideaDir, run, v.History, err)
	}
	if !v.History.MidIdea() {
		return ctx, func() {}, nil
	}
	return acquireScoped(ctx, ideaDir, run)
}

func acquireScoped(ctx context.Context, ideaDir, run string) (context.Context, func(), error) {
	path, err := lockPath(ideaDir, "driver")
	if err != nil {
		return ctx, nil, err
	}
	if l, _ := ctx.Value(leaseKey{}).(*lease); l != nil && l.path == path && l.run == run && l.active.Load() {
		if err := l.owner.Check(); err != nil {
			return ctx, nil, err
		}
		return ctx, func() {}, nil
	}
	owner, err := pidlease.TryAcquire(path, filepath.Base(ideaDir)+"/"+run)
	if err != nil {
		return ctx, nil, fmt.Errorf("idea driving lock held: %s (competing run %s): %w", filepath.Base(ideaDir), run, err)
	}
	l := &lease{path: path, run: run, owner: owner}
	l.active.Store(true)
	release := func() {
		if l.active.Swap(false) {
			owner.Release()
		}
	}
	return context.WithValue(ctx, leaseKey{}, l), release, nil
}
func requireLease(ctx context.Context, ideaDir string) error {
	path, _ := lockPath(ideaDir, "driver")
	l, _ := ctx.Value(leaseKey{}).(*lease)
	if l == nil || l.path != path || !l.active.Load() {
		return fmt.Errorf("quota mutation requires idea driving lease")
	}
	return l.owner.Check()
}

func DrivingRunID(ctx context.Context) string {
	if l, _ := ctx.Value(leaseKey{}).(*lease); l != nil && l.active.Load() {
		return l.run
	}
	return ""
}

// ProjectionLock serializes the short append/commit critical section across CLI
// processes. It is separate from the driving lease: a supervised signer can
// append while its parent owns the driving lease, but never during a transition.
func ProjectionLock(ideaDir string) (func(), error) {
	h, err := protocol.InspectQuota(ideaDir)
	if err != nil {
		return nil, err
	}
	if !h.History.MidIdea() {
		return func() {}, nil
	}
	return projectionLock(ideaDir)
}

func projectionLock(ideaDir string) (func(), error) {
	path, err := lockPath(ideaDir, "projection")
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(projectionWait)
	for {
		owner, err := pidlease.TryAcquire(path, filepath.Base(ideaDir)+"/projection")
		if err == nil {
			return owner.Release, nil
		}
		if !errors.Is(err, pidlease.ErrHeld) || !time.Now().Before(deadline) {
			return nil, fmt.Errorf("quota projection lock held: %w", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

}

// CheckLease binds a launch to this exact idea and driving run, not just any lease.
func CheckLease(ctx context.Context, ideaDir, run string) error {
	if err := requireLease(ctx, ideaDir); err != nil {
		return err
	}
	if DrivingRunID(ctx) != run {
		return fmt.Errorf("quota driving run identity mismatch")
	}
	return nil
}

// Lock state stays on the deck filesystem, in the already ignored runtime tree.
// Relative identity is stable when different hosts mount the deck at other paths.
func lockPath(ideaDir, kind string) (string, error) {
	dir, err := filepath.EvalSymlinks(ideaDir)
	if err != nil {
		return "", err
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	deck := dir
	for filepath.Base(deck) != protocol.DeckDir {
		parent := filepath.Dir(deck)
		if parent == deck {
			return "", fmt.Errorf("idea is outside a deck")
		}
		deck = parent
	}
	rel, err := filepath.Rel(deck, dir)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(filepath.ToSlash(rel))))
	return filepath.Join(filepath.Dir(deck), ".parley-runtime", "membership", key, kind+".lease"), nil
}

var projectionWait = 2 * time.Second
