package budget

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type activeReservationFixture struct {
	store  Store
	charge func() error
	finish func()
}

func newActiveReservationFixture(t *testing.T, kind Kind) activeReservationFixture {
	t.Helper()
	ctx := context.Background()
	root := t.TempDir()
	var f activeReservationFixture
	if kind == DriverStep {
		b, err := EnsureStepBinding(ctx, root, "idea", 5, 0)
		if err != nil {
			t.Fatal(err)
		}
		scoped, finish, err := OpenStepSession(ctx, b)
		if err != nil {
			t.Fatal(err)
		}
		f = activeReservationFixture{b.Store, func() error { return ChargeStep(scoped) }, finish}
	} else {
		b, err := EnsureCycleBinding(ctx, root, "idea", kind, 5, 0, "", "")
		if err != nil {
			t.Fatal(err)
		}
		scoped, finish, err := OpenCycleSession(ctx, b)
		if err != nil {
			t.Fatal(err)
		}
		f = activeReservationFixture{b.Store, func() error { _, err := ChargeCycle(scoped, kind); return err }, finish}
	}
	t.Cleanup(f.finish)
	if err := f.charge(); err != nil {
		t.Fatal(err)
	}
	return f
}

func publishReservationFixture(t *testing.T, store Store, state Snapshot) {
	t.Helper()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeSynced(filepath.Join(store.Dir, "ledger.json"), data); err != nil {
		t.Fatal(err)
	}
}

// Counts alone cannot prove that THIS still-active operation retains the
// reservation which granted it permission. All mutations remain valid ledger
// JSON and preserve the aggregate number of charged entries.
func TestActiveReservationRejectsChangedOriginalCharge(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		for _, mutation := range []string{"replace-id", "change-time", "unknown-reserve", "changed-reserve", "change-epoch"} {
			t.Run(string(kind)+"/"+mutation, func(t *testing.T) {
				f := newActiveReservationFixture(t, kind)
				state, err := f.store.Inspect(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				if len(state.Entries) != 1 {
					t.Fatal("expected exactly the active reservation")
				}
				original, err := os.ReadFile(filepath.Join(f.store.Dir, "ledger.json"))
				if err != nil {
					t.Fatal(err)
				}
				var id string
				var entry Reservation
				for id, entry = range state.Entries {
				}
				switch mutation {
				case "replace-id":
					delete(state.Entries, id)
					id = key("different-logical-action")
					// Preserve structural identity consistency so the live receipt,
					// rather than a malformed ledger, must reject the replacement.
					if entry.Action != nil {
						entry.Action.EntrySHA256 = id
					}
				case "change-time":
					entry.ReservedAt = entry.ReservedAt.Add(time.Nanosecond)
				case "unknown-reserve":
					entry.ReserveMicros = nil
				case "changed-reserve":
					value := int64(1)
					entry.ReserveMicros = &value
				case "change-epoch":
					state.StartedAt = state.StartedAt.Add(-time.Second)
				}
				if entry.Action != nil {
					entry.Action.ReservationSHA256 = originalChargeDigest(state.Scope, state.StartedAt, id, entry)
				}
				state.Entries[id] = entry
				publishReservationFixture(t, f.store, state)
				if _, err := f.store.Inspect(context.Background()); err != nil {
					t.Fatalf("fixture must remain structurally valid: %v", err)
				}
				before, err := os.ReadFile(filepath.Join(f.store.Dir, "ledger.json"))
				if err != nil {
					t.Fatal(err)
				}
				if err := f.charge(); err == nil {
					t.Fatal("nested work permitted despite changed original reservation")
				}
				after, err := os.ReadFile(filepath.Join(f.store.Dir, "ledger.json"))
				if err != nil {
					t.Fatal(err)
				}
				if string(before) != string(after) {
					t.Fatal("refusal modified or recreated accounting")
				}
				if err := writeSynced(filepath.Join(f.store.Dir, "ledger.json"), original); err != nil {
					t.Fatal(err)
				}
				if err := f.charge(); err == nil {
					t.Fatal("restoring bytes revived a session with an observed continuity refusal")
				}
			})
		}
	}
}

func TestActiveReservationAllowsLaterChargesAndSettlement(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			f := newActiveReservationFixture(t, kind)
			ctx := context.Background()
			state, err := f.store.Inspect(ctx)
			if err != nil {
				t.Fatal(err)
			}
			var activeID string
			for activeID = range state.Entries {
			}
			zero := int64(0)
			// Later, separately charged operations fill the actual five-action
			// cap; the still-active group may continue under its original charge.
			for _, id := range []string{"second", "third", "fourth", "fifth"} {
				if _, err := f.store.Reserve(ctx, Request{ID: id, Kind: kind, ReserveMicros: &zero}, Limits{Actions: map[Kind]int{kind: 5}}); err != nil {
					t.Fatal(err)
				}
			}
			// Settlement adds an observation; it does not alter the reservation
			// that granted the live group permission. Exercise the real update
			// transaction because session raw random IDs are intentionally private.
			if _, err := f.store.update(ctx, func(s *Snapshot, _ time.Time) error {
				entry := s.Entries[activeID]
				entry.Settled, entry.ActualMicros = true, &zero
				s.Entries[activeID] = entry
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(filepath.Join(f.store.Dir, "ledger.json"))
			if err != nil {
				t.Fatal(err)
			}
			for i := 0; i < 3; i++ {
				if err := f.charge(); err != nil {
					t.Fatalf("valid original reservation was refused: %v", err)
				}
			}
			after, err := os.ReadFile(filepath.Join(f.store.Dir, "ledger.json"))
			if err != nil {
				t.Fatal(err)
			}
			if string(before) != string(after) {
				t.Fatal("continuing the same group changed accounting")
			}
		})
	}
}

func TestReservationReceiptPinsImmutableFieldsOnly(t *testing.T) {
	ctx := context.Background()
	store := Store{Dir: t.TempDir(), Scope: "receipt-fixture"}
	zero := int64(0)
	state, err := store.Reserve(ctx, Request{ID: "known-fixture-action", Kind: DriverStep, ReserveMicros: &zero}, Limits{})
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := newReservationReceipt(state, "known-fixture-action", DriverStep)
	if err != nil {
		t.Fatal(err)
	}
	// The captured numeric value must not alias the snapshot's pointer.
	*state.Entries[key("known-fixture-action")].ReserveMicros = 9
	if err := receipt.check(state); err == nil {
		t.Fatal("receipt adopted mutation of its input snapshot")
	}
	state, err = store.Settle(ctx, "known-fixture-action", nil)
	if err != nil {
		t.Fatal(err)
	}
	state, err = store.ReconcileUnknown(ctx, "known-fixture-action", "fixture-observation", 4, "Fixture-only conservative unknown-cost observation")
	if err != nil {
		t.Fatal(err)
	}
	if err := receipt.check(state); err != nil {
		t.Fatalf("legitimate settlement/reconciliation invalidated the original charge: %v", err)
	}
	entry := state.Entries[key("known-fixture-action")]
	entry.Kind = Fixup
	state.Entries[key("known-fixture-action")] = entry
	if err := receipt.check(state); err == nil {
		t.Fatal("receipt accepted a different accounting kind")
	}
	if _, err := newReservationReceipt(state, "missing", Fixup); err == nil {
		t.Fatal("captured a receipt for a missing reservation")
	}
	if _, err := newReservationReceipt(state, "known-fixture-action", DriverStep); err == nil {
		t.Fatal("captured a receipt for the wrong kind")
	}
}
