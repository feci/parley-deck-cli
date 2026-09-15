package budget

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestContinuityPublicationFailureRetainsChargeWithoutGrant(t *testing.T) {
	s := testStore(t)
	s.persist = func(path string, data []byte) error {
		if err := writeSynced(path, data); err != nil {
			return err
		}
		// Actual immutable-witness publication now fails after the charged
		// ledger is durable, without allowing a process to start.
		return os.Mkdir(s.continuityPath(), 0700)
	}
	limits := Limits{Actions: map[Kind]int{Launch: 1}}
	if _, err := s.Reserve(context.Background(), Request{ID: "spent", Kind: Launch}, limits); err == nil {
		t.Fatal("failed witness persistence granted work")
	}
	state, err := s.Inspect(context.Background())
	if err != nil || len(state.Entries) != 1 {
		t.Fatalf("lost conservative charge: %+v %v", state, err)
	}
	if err := os.Remove(s.continuityPath()); err != nil {
		t.Fatal(err)
	}
	s.persist = nil
	if _, err := s.Reserve(context.Background(), Request{ID: "spent", Kind: Launch}, limits); !errors.Is(err, ErrReserved) {
		t.Fatalf("replayed failed attempt: %v", err)
	}
	if _, err := s.Reserve(context.Background(), Request{ID: "new", Kind: Launch}, limits); !errors.Is(err, ErrLimit) {
		t.Fatalf("failed attempt was refunded: %v", err)
	}
}

func TestLegacyLedgerAcquiresWitnessWithoutLosingCharges(t *testing.T) {
	s := testStore(t)
	if _, err := s.Reserve(context.Background(), Request{ID: "one", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(s.continuityPath()); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Reserve(context.Background(), Request{ID: "two", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	state, err := s.Inspect(context.Background())
	if err != nil || len(state.Entries) != 2 {
		t.Fatalf("legacy charges lost: %+v %v", state, err)
	}
	if err := os.Remove(filepath.Join(s.Dir, "ledger.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Reserve(context.Background(), Request{ID: "three", Kind: Launch}, Limits{}); err == nil {
		t.Fatal("migrated ledger deletion reset charges")
	}
}

func TestDeletedLedgerCannotResetCharges(t *testing.T) {
	s := testStore(t)
	limits := Limits{Actions: map[Kind]int{Launch: 1}}
	if _, err := s.Reserve(context.Background(), Request{ID: "spent", Kind: Launch}, limits); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(s.Dir, "ledger.json")); err != nil {
		t.Fatal(err)
	}
	resumed := Store{Dir: s.Dir, Scope: s.Scope}
	if _, err := resumed.Reserve(context.Background(), Request{ID: "new-run", Kind: Launch}, limits); err == nil {
		t.Fatal("deleting only ledger.json granted another attempt above the frozen cap")
	}
	if _, err := os.Stat(filepath.Join(s.Dir, "ledger.json")); !os.IsNotExist(err) {
		t.Fatal("missing charged history was silently recreated")
	}
}
