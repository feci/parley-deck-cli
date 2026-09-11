package budget

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"path/filepath"
	"sync"
)

type cycleSessionKey struct{ Kind Kind }
type cycleRefusalKey struct{ Kind Kind }
type cycleSession struct {
	mu                sync.Mutex
	binding           CycleBinding
	active, attempted bool
	ordinal           int
	err               error
}

// CycleSessionMatches recognizes only a live nested operation at the exact
// binding. Its output directories may already exist before its first child
// reaches the reservation boundary; those are not additional historical cycles.
func CycleSessionMatches(ctx context.Context, b *CycleBinding) bool {
	s, ok := ctx.Value(cycleSessionKey{b.Policy.Kind}).(*cycleSession)
	if !ok {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active && sameCycleAuthority(s.binding.Policy, b.Policy) && s.binding.Store.Dir == b.Store.Dir && s.binding.Store.Scope == b.Store.Scope
}

func OpenCycleSession(ctx context.Context, b *CycleBinding) (context.Context, func(), error) {
	noop := func() {}
	if b == nil {
		return ctx, noop, errors.New("cycle session requires a saved binding")
	}
	if err, ok := ctx.Value(cycleRefusalKey{b.Policy.Kind}).(error); ok {
		return ctx, noop, err
	}
	if old, ok := ctx.Value(cycleSessionKey{b.Policy.Kind}).(*cycleSession); ok {
		old.mu.Lock()
		defer old.mu.Unlock()
		if !old.active {
			return ctx, noop, errors.New("cycle session has ended")
		}
		if !sameCycleAuthority(old.binding.Policy, b.Policy) || old.binding.Store.Scope != b.Store.Scope || old.binding.Store.Dir != b.Store.Dir {
			return ctx, noop, errors.New("nested cycle scope mismatch")
		}
		return ctx, noop, nil
	}
	copy := *b
	copy.Policy = cloneCyclePolicy(b.Policy)
	s := &cycleSession{binding: copy, active: true}
	finish := func() { s.mu.Lock(); s.active = false; s.mu.Unlock() }
	return context.WithValue(ctx, cycleSessionKey{b.Policy.Kind}, s), finish, nil
}

// DeferCycleRefusal preserves the error until requested invocation evidence is
// written. It never converts failed configuration into an unbudgeted launch.
func DeferCycleRefusal(ctx context.Context, kind Kind, err error) context.Context {
	return context.WithValue(ctx, cycleRefusalKey{kind}, err)
}

func ChargeCycle(ctx context.Context, kind Kind) (int, error) {
	if err, ok := ctx.Value(cycleRefusalKey{kind}).(error); ok {
		return 0, err
	}
	s, ok := ctx.Value(cycleSessionKey{kind}).(*cycleSession)
	if !ok {
		return 0, errors.New("protocol cycle has no bound accounting session")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		return 0, errors.New("cycle session has ended")
	}
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if s.attempted {
		if s.err == nil {
			current, err := s.binding.current()
			var state Snapshot
			if err == nil {
				state, err = current.Store.Inspect(ctx)
			}
			if err == nil {
				err = checkProtocolMigrationCharges(filepath.Dir(current.Store.Dir), current.Policy.MigrationSHA256, state)
			}
			if err == nil && s.binding.Count(state) < s.ordinal {
				err = errors.New("cycle accounting lost a reserved charge")
			}
			s.err = err
		}
		return s.ordinal, s.err
	}
	s.attempted = true
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		s.err = err
		return 0, err
	}
	s.ordinal, s.err = s.binding.Reserve(ctx, "cycle:"+hex.EncodeToString(id[:]))
	return s.ordinal, s.err
}
