package budget

import (
	"errors"
	"time"
)

// reservationReceipt pins the immutable part of the exact published charge
// that granted this live session permission. Aggregate counts cannot replace
// that identity. Settlements and operator reconciliations remain independent
// observations and may legitimately be appended after reservation.
//
// This is a live-session witness, not a cross-process execution/replay token.
type reservationReceipt struct {
	scope, id     string
	startedAt     time.Time
	kind          Kind
	reservedAt    time.Time
	reserveMicros *int64
}

func newReservationReceipt(state Snapshot, id string, kind Kind) (reservationReceipt, error) {
	entry, ok := state.Entries[key(id)]
	if id == "" || state.Scope == "" || state.StartedAt.IsZero() || !ok || entry.Kind != kind || !validKind(kind) {
		return reservationReceipt{}, errors.New("published accounting lacks the requested reservation")
	}
	return reservationReceipt{
		scope: state.Scope, id: key(id), startedAt: state.StartedAt,
		kind: kind, reservedAt: entry.ReservedAt, reserveMicros: copyInt(entry.ReserveMicros),
	}, nil
}

func (r reservationReceipt) check(state Snapshot) error {
	entry, ok := state.Entries[r.id]
	if r.id == "" || r.scope != state.Scope || !r.startedAt.Equal(state.StartedAt) || !ok ||
		entry.Kind != r.kind || !entry.ReservedAt.Equal(r.reservedAt) ||
		(entry.ReserveMicros == nil) != (r.reserveMicros == nil) ||
		(entry.ReserveMicros != nil && *entry.ReserveMicros != *r.reserveMicros) {
		return errors.New("active accounting session lost or changed its original reservation")
	}
	return nil
}
