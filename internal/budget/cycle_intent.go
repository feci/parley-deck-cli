package budget

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

// CycleReservationObserver adds a durable preparation boundary before Reserve.
// It runs under the same cycle guard as both ordinary observer callbacks. Its
// input describes the actual request; no reservation timestamp exists yet.
type CycleReservationObserver interface {
	PrepareCycleReservation(context.Context, CycleBinding, Snapshot, CycleReservationIntent) error
}

type CycleReservationIntent struct {
	Version       int             `json:"version"`
	Policy        CyclePolicy     `json:"policy"`
	StartedAt     time.Time       `json:"started_at"`
	EntryKey      string          `json:"entry_key"`
	ReserveMicros *int64          `json:"reserve_micros"`
	Action        *ActionIdentity `json:"action"`
}

func cycleLimits(p CyclePolicy) Limits {
	l := Limits{Actions: map[Kind]int{}, Denied: map[Kind]bool{}}
	if remaining := p.Maximum - p.Carried; remaining > 0 {
		l.Actions[p.Kind] = remaining
	} else {
		l.Denied[p.Kind] = true
	}
	return l
}

func newCycleReservationIntent(ctx context.Context, b CycleBinding, s Snapshot, req Request, limits Limits) (CycleReservationIntent, error) {
	a, err := actionIdentityFor(ctx, b.Store.Scope, req, limits)
	if err != nil {
		return CycleReservationIntent{}, err
	}
	i := CycleReservationIntent{1, cloneCyclePolicy(b.Policy), s.StartedAt, key(req.ID), copyInt(req.ReserveMicros), a}
	return i, i.CheckBefore(b.Policy, s)
}

func (i CycleReservationIntent) validate(s Snapshot) error {
	p := i.Policy
	if i.Version != 1 || (p.Version != 1 && p.Version != 2) || p.Kind != Fixup || p.Carried != 0 || p.Maximum <= 0 ||
		!originHash(p.TrajectorySHA256) || !originHash(i.EntryKey) || p.Scope != s.Scope ||
		i.StartedAt.IsZero() || !i.StartedAt.Equal(s.StartedAt) || i.ReserveMicros == nil || *i.ReserveMicros != 0 {
		return errors.New("invalid original cycle reservation intent")
	}
	raw, _ := json.Marshal(p)
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil || validateCycleExtensions(p, fields) != nil {
		return errors.New("precharge intent has an invalid original policy extension chain")
	}
	if i.Action != nil {
		a := i.Action
		if validateActionInput(a.Input) != nil || a.Version != 1 || a.EntrySHA256 != i.EntryKey ||
			a.LogicalSHA256 != actionLogicalDigest(p.Scope, a.Input) || a.LimitsSHA256 != migrationDigest(cycleLimits(p)) || a.ReservationSHA256 != "" {
			return errors.New("precharge intent changed original action identity or effective limits")
		}
	}
	return nil
}

// CheckPolicy permits only an exact original prefix of the current validated
// extension chain. A later grant cannot change the intent's effective ceilings.
func (i CycleReservationIntent) CheckPolicy(current CyclePolicy) error {
	n := len(i.Policy.Extensions)
	if len(current.Extensions) < n {
		return errors.New("original precharge policy history is missing")
	}
	prior := cloneCyclePolicy(current)
	if n == 0 {
		prior = originalCyclePolicy(prior)
	} else {
		prior.Extensions = prior.Extensions[:n]
		prior.Maximum = prior.Extensions[n-1].Maximum
	}
	if CyclePolicyDigest(prior) != CyclePolicyDigest(i.Policy) {
		return errors.New("original precharge policy differs from the retained extension history")
	}
	return nil
}

func (i CycleReservationIntent) CheckBefore(p CyclePolicy, s Snapshot) error {
	if err := i.validate(s); err != nil {
		return err
	}
	if CyclePolicyDigest(i.Policy) != CyclePolicyDigest(p) {
		return errors.New("precharge intent changed original cycle policy")
	}
	if _, exists := s.Entries[i.EntryKey]; exists {
		return ErrReserved
	}
	if len(s.Entries) >= i.Policy.Maximum-i.Policy.Carried {
		return ErrLimit
	}
	return nil
}

// CheckCharge compares the published charge to the retained original request,
// including its action input and effective original ceilings. Settlement and a
// later operator extension cannot rewrite that contract. Nil Action explicitly
// retains an unbound direct caller; it never acquires invented semantics.
func (i CycleReservationIntent) CheckCharge(s Snapshot) (CycleCharge, error) {
	if err := i.validate(s); err != nil {
		return CycleCharge{}, err
	}
	c, err := PublishedCycleCharge(s, i.EntryKey)
	if err != nil {
		return c, err
	}
	e := s.Entries[i.EntryKey]
	a := cloneActionIdentity(e.Action)
	if a != nil {
		if err := a.validate(s.Scope, i.EntryKey, s.StartedAt, e); err != nil {
			return CycleCharge{}, err
		}
		a.ReservationSHA256 = ""
	}
	if !reflect.DeepEqual(a, i.Action) || !reflect.DeepEqual(e.ReserveMicros, i.ReserveMicros) {
		return CycleCharge{}, errors.New("published cycle charge differs from its original precharge intent")
	}
	return c, nil
}
