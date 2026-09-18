package budget

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// CycleObserver is trusted runtime integration, not a participant-supplied
// verdict. Both callbacks run under the cycle policy guard, before permission
// to execute is returned. A failed AfterCycle never refunds the charge.
// Callbacks must not reacquire the same policy guard or reserve another cycle.
type CycleObserver interface {
	BeforeCycle(context.Context, CycleBinding, Snapshot) error
	AfterCycle(context.Context, CycleBinding, Snapshot, string) error
}
type cycleObserverKey struct{}

func WithCycleObserver(ctx context.Context, observer CycleObserver) context.Context {
	return context.WithValue(ctx, cycleObserverKey{}, observer)
}
func requiredCycleObserver(ctx context.Context, b *CycleBinding) (CycleObserver, error) {
	if b.Policy.TrajectorySHA256 == "" {
		if _, err := os.Lstat(filepath.Join(filepath.Dir(b.Store.Dir), "trajectory.json")); !os.IsNotExist(err) {
			return nil, errors.New("trajectory initialization must be replayed before fixup")
		}
		return nil, nil
	}
	o, ok := ctx.Value(cycleObserverKey{}).(CycleObserver)
	if !ok || o == nil {
		return nil, errors.New("fixup requires the persisted trajectory observer; no work was authorized")
	}
	return o, nil
}

// CycleCharge identifies the immutable original published reservation. The
// already-hashed ledger entry key avoids exposing raw caller action strings.
type CycleCharge struct {
	ActionSHA256  string    `json:"action_sha256,omitempty"`
	Scope         string    `json:"scope"`
	StartedAt     time.Time `json:"started_at"`
	EntryKey      string    `json:"entry_key"`
	Kind          Kind      `json:"kind"`
	ReservedAt    time.Time `json:"reserved_at"`
	ReserveMicros *int64    `json:"reserve_micros"`
}

func PublishedCycleCharge(s Snapshot, entryKey string) (CycleCharge, error) {
	e, ok := s.Entries[entryKey]
	if !ok || !validCycleDecision("entry", "entry", entryKey) || e.Kind != Fixup || s.Scope == "" || s.StartedAt.IsZero() || e.ReservedAt.IsZero() {
		return CycleCharge{}, errors.New("missing exact fixup reservation")
	}
	return CycleCharge{Scope: s.Scope, StartedAt: s.StartedAt, EntryKey: entryKey, Kind: e.Kind, ReservedAt: e.ReservedAt, ReserveMicros: copyInt(e.ReserveMicros), ActionSHA256: actionIdentityDigest(e.Action)}, nil
}
func (c CycleCharge) Check(s Snapshot) error {
	r := reservationReceipt{scope: c.Scope, id: c.EntryKey, startedAt: c.StartedAt, kind: c.Kind, reservedAt: c.ReservedAt, reserveMicros: copyInt(c.ReserveMicros), actionSHA256: c.ActionSHA256}
	return r.check(s)
}

// ActiveCycleCharge never creates or charges a session. Consumers use it only
// after the common precharge boundary has granted this exact live session.
func ActiveCycleCharge(ctx context.Context) (CycleCharge, error) {
	s, ok := ctx.Value(cycleSessionKey{Fixup}).(*cycleSession)
	if !ok {
		return CycleCharge{}, errors.New("no live fixup accounting session")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active || !s.attempted || s.err != nil {
		return CycleCharge{}, errors.New("fixup reservation is unavailable")
	}
	r := s.receipt
	c := CycleCharge{Scope: r.scope, StartedAt: r.startedAt, EntryKey: r.id, Kind: r.kind, ReservedAt: r.reservedAt, ReserveMicros: copyInt(r.reserveMicros), ActionSHA256: r.actionSHA256}
	state, err := s.binding.Store.Inspect(ctx)
	if err != nil {
		return CycleCharge{}, err
	}
	if err := c.Check(state); err != nil {
		s.err = err
		return CycleCharge{}, err
	}
	return c, nil
}

// ActivateCycleTrajectory binds a fresh cycle policy to a durable observation
// policy. Existing charged/imported histories require explicit migration, not
// adoption of an empty trajectory. initialize publishes BEFORE the reference;
// interrupted initialization must be replayed exactly by its owning component.
// Activation shares the reservation guard and grants no execution by itself.
func ActivateCycleTrajectory(ctx context.Context, root, idea, expected, sha string, initialize func(CycleBinding, Snapshot) error) error {
	if !validCycleDecision("policy", "policy", sha) || !validCycleDecision("expected", "expected", expected) || initialize == nil {
		return errors.New("invalid trajectory activation")
	}
	b, err := LoadCycleBinding(ctx, root, idea, Fixup)
	if err != nil {
		return err
	}
	if b == nil {
		return errors.New("initialize the original fixup policy first")
	}
	release, err := AcquireResourceGuard(ctx, filepath.Dir(b.Store.Dir))
	if err != nil {
		return err
	}
	defer release()
	b, err = b.current()
	if err != nil {
		return err
	}
	if b.Policy.TrajectorySHA256 != "" && b.Policy.TrajectorySHA256 != sha {
		return errors.New("trajectory policy already frozen")
	}
	original := b.Policy
	original.TrajectorySHA256 = ""
	if CyclePolicyDigest(original) != expected {
		return errors.New("cycle policy changed since trajectory preview")
	}
	state, err := b.Store.Inspect(ctx)
	if err != nil {
		return err
	}
	if b.Policy.Carried != 0 || b.Policy.MigrationSHA256 != "" || len(b.Policy.Extensions) != 0 || len(state.Entries) != 0 {
		return errors.New("trajectory activation requires a fresh uncharged fixup policy")
	}
	if err := initialize(*b, state); err != nil {
		return err
	}
	b.Policy.TrajectorySHA256 = sha
	data, err := json.MarshalIndent(b.Policy, "", "  ")
	if err != nil {
		return err
	}
	return writeSynced(filepath.Join(filepath.Dir(b.Store.Dir), "policy.json"), append(data, '\n'))
}

// ActiveCycleTrajectory pins the requirement already accepted by this live
// session. Losing durable state later cannot reinterpret it as opt-out.
func ActiveCycleTrajectory(ctx context.Context) (string, error) {
	s, ok := ctx.Value(cycleSessionKey{Fixup}).(*cycleSession)
	if !ok {
		return "", errors.New("no live fixup accounting session")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active || !s.attempted || s.err != nil {
		return "", errors.New("fixup reservation is unavailable")
	}
	return s.binding.Policy.TrajectorySHA256, nil
}
