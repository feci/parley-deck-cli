package budget

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
)

// A session groups the nested work of ONE synchronous transition. It is runtime
// state, not an authorization token for participant artifacts or another call.
type stepSession struct {
	mu                sync.Mutex
	binding           *StepBinding
	active, attempted bool
	err               error
}
type stepSessionKey struct{}
type stepRefusalKey struct{}

// OpenStepSession does not charge reads/awaits. Close invalidates even a retained
// child context. Cooperating nested runners share the same one-step charge.
func OpenStepSession(ctx context.Context, b *StepBinding) (context.Context, func(), error) {
	noop := func() {}
	if b == nil {
		if _, exists := ctx.Value(stepSessionKey{}).(*stepSession); exists {
			return ctx, noop, errors.New("nested step binding is unavailable")
		}
		return ctx, noop, nil
	}
	if old, ok := ctx.Value(stepSessionKey{}).(*stepSession); ok {
		old.mu.Lock()
		defer old.mu.Unlock()
		if !old.active {
			return ctx, noop, errors.New("driver step session has ended")
		}
		if old.binding.Store.Dir != b.Store.Dir || old.binding.Store.Scope != b.Store.Scope {
			return ctx, noop, errors.New("nested step scope mismatch")
		}
		return ctx, noop, nil
	}
	if _, err := b.Check(ctx); err != nil {
		return ctx, noop, err
	}
	binding := *b
	s := &stepSession{binding: &binding, active: true}
	close := func() { s.mu.Lock(); s.active = false; s.mu.Unlock() }
	return context.WithValue(ctx, stepSessionKey{}, s), close, nil
}

// ChargeStep reserves durably before the first mutation/process of a transition.
// The cost here is exactly zero monetary spend: process costs are separately
// charged by the launch policy. Any failure stays cached within this session.
func ChargeStep(ctx context.Context) error {
	s, ok := ctx.Value(stepSessionKey{}).(*stepSession)
	if !ok {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		return errors.New("driver step session has ended")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.attempted {
		if s.err == nil {
			// The inclusive step cap permits the rest of this charged transition.
			// It does not permit a later launch after lifetime time has expired,
			// or after required accounting state has disappeared/corrupted.
			state, err := s.binding.Store.Inspect(ctx)
			if err == nil {
				err = s.binding.checkTime(state)
			}
			s.err = err
		}
		return s.err
	}
	s.attempted = true
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		s.err = err
		return err
	}
	zero := int64(0)
	_, s.err = s.binding.Store.Reserve(ctx, Request{ID: "step:" + hex.EncodeToString(id[:]), Kind: DriverStep, ReserveMicros: &zero}, s.binding.Policy.limits())
	if s.err != nil {
		s.err = fmt.Errorf("driver step reservation refused: %w", s.err)
	}
	return s.err
}

// JoinStepSession is used by manual protocol runners. It loads an existing
// binding; it does not invent a policy or overwrite runtime/operator limits.
func JoinStepSession(ctx context.Context, root, idea string) (context.Context, func(), error) {
	if err, ok := ctx.Value(stepRefusalKey{}).(error); ok {
		return ctx, func() {}, err
	}
	b, err := LoadStepBinding(ctx, root, idea)
	if err != nil {
		return ctx, func() {}, err
	}
	return OpenStepSession(ctx, b)
}

// GroupStepSession groups a manual runner without suppressing invocation
// evidence. A refusal is retained in the context until the common launch
// boundary has written its requested record; that boundary records the failed
// terminal without starting a child. An all-artifacts-present no-op stays free.
func GroupStepSession(ctx context.Context, root, idea string) (context.Context, func()) {
	scoped, finish, err := JoinStepSession(ctx, root, idea)
	if err != nil {
		return context.WithValue(ctx, stepRefusalKey{}, err), func() {}
	}
	return scoped, finish
}
