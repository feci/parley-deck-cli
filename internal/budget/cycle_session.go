package budget

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
)

type cycleSessionKey struct{ Kind Kind }
type cycleRefusalKey struct{ Kind Kind }
type cycleSession struct {
	input             ActionInput
	mu                sync.Mutex
	binding           CycleBinding
	receipt           reservationReceipt
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
		if err := checkSessionActionInput(ctx, old.input); err != nil {
			return ctx, noop, err
		}
		if !old.active {
			return ctx, noop, errors.New("cycle session has ended")
		}
		if !sameCycleAuthority(old.binding.Policy, b.Policy) || old.binding.Store.Scope != b.Store.Scope || old.binding.Store.Dir != b.Store.Dir {
			return ctx, noop, errors.New("nested cycle scope mismatch")
		}
		return ctx, noop, nil
	}
	ctx = InheritActionInput(ctx, ActionInput{Operation: string(b.Policy.Kind), Basis: "policy-only", InputSHA256: CyclePolicyDigest(b.Policy)})
	copy := *b
	copy.Policy = cloneCyclePolicy(b.Policy)
	s := &cycleSession{binding: copy, active: true, input: sessionActionInput(ctx)}
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
	if err := checkSessionActionInput(ctx, s.input); err != nil {
		s.attempted = true
		s.err = err
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
			if err == nil {
				err = s.receipt.check(state)
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
	s.ordinal, s.receipt, s.err = s.binding.reserveWithReceipt(ctx, "cycle:"+hex.EncodeToString(id[:]))
	return s.ordinal, s.err
}

// PreflightCycleCharge is a read-only, non-mutating preflight (MINOR-1 known
// cross-resource exhaustion): it reports whether the cycle charge this context
// would attempt at the given binding is already KNOWN to refuse from
// NEW-RESERVATION exhaustion, so a caller about to spend a DIFFERENT budget
// first (a driver step) can refuse for free. It charges nothing, persists
// nothing, takes no resource guard and never refunds; reserveWithReceipt under
// the guard remains the only charging authority.
//
// It mirrors ChargeCycle's reuse semantics rather than a naive
// current-count-only check: a live session of this kind that already charged
// successfully for this logical action stays reusable at the cap (nil), and a
// session with a cached refusal replays it. Only a charge that would be a NEW
// reservation compares the binding's current count with its frozen maximum
// (a zero maximum forbids the operation, exactly as cycleLimits denies it).
//
// The coverage is deliberately no wider than that. A cached SUCCESSFUL session
// is replayed as nil WITHOUT ChargeCycle's attempted-branch revalidation
// (checkProtocolMigrationCharges, receipt.check, the lost-reserved-charge
// guard), and a deferred refusal under cycleRefusalKey — which ChargeCycle
// checks first — is not consulted (unreachable from the current call sites).
// A refusal that only materializes at the real ChargeCycle stays spent: a
// different budget charged between this preflight and the refusal is not
// refunded. Concurrent exhaustion after this preflight still refuses at
// ChargeCycle, and an actual failed attempt stays spent. A nil binding
// charges nothing, so there is nothing to refuse.
func PreflightCycleCharge(ctx context.Context, b *CycleBinding) error {
	if b == nil {
		return nil
	}
	if s, ok := ctx.Value(cycleSessionKey{b.Policy.Kind}).(*cycleSession); ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.active {
			return errors.New("cycle session has ended")
		}
		if s.attempted {
			return s.err
		}
	}
	state, err := b.Store.Inspect(ctx)
	if err != nil {
		return err
	}
	if b.Count(state) >= b.Policy.Maximum {
		return fmt.Errorf("%w: %s protocol cycles", ErrLimit, b.Policy.Kind)
	}
	return nil
}
