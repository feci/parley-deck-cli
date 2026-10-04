// Package membership implements quota transition projections and the driving lease.
// Immutable authority lives in quota; read-only protocol consumers never call Repair.
package membership

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"parley-deck-cli/internal/protocol"
)

type leaseKey struct{}
type lease struct {
	path, run string
	active    atomic.Bool
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
	path, err := filepath.Abs(filepath.Join(ideaDir, "quota-driver.lock"))
	if err != nil {
		return ctx, nil, err
	}
	if l, _ := ctx.Value(leaseKey{}).(*lease); l != nil && l.path == path && l.run == run && l.active.Load() {
		return ctx, func() {}, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return ctx, nil, err
	}
	ok, err := tryLock(f)
	if err != nil || !ok {
		f.Close()
		if err == nil {
			err = fmt.Errorf("idea driving lock held: %s (competing run %s)", filepath.Base(ideaDir), run)
		}
		return ctx, nil, err
	}
	l := &lease{path: path, run: run}
	l.active.Store(true)
	release := func() {
		if l.active.Swap(false) {
			unlock(f)
			f.Close()
		}
	}
	return context.WithValue(ctx, leaseKey{}, l), release, nil
}
func requireLease(ctx context.Context, ideaDir string) error {
	path, _ := filepath.Abs(filepath.Join(ideaDir, "quota-driver.lock"))
	l, _ := ctx.Value(leaseKey{}).(*lease)
	if l == nil || l.path != path || !l.active.Load() {
		return fmt.Errorf("quota mutation requires idea driving lease")
	}
	return nil
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
	path := filepath.Join(ideaDir, "quota-projection.lock")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	ok, err := tryLock(f)
	if err != nil || !ok {
		f.Close()
		if err == nil {
			err = fmt.Errorf("quota projection lock held")
		}
		return nil, err
	}
	return func() { unlock(f); f.Close() }, nil
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
