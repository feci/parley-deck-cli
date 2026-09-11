package app

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
)

func TestBudgetControlAttendanceScopeAndReplay(t *testing.T) {
	ctx := context.Background()
	s := budget.Store{Dir: t.TempDir(), Scope: "idea"}
	if _, err := s.Reserve(ctx, budget.Request{ID: "invocation-x", Kind: budget.Launch}, budget.Limits{}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"reconcile", "--ledger", s.Dir, "--scope", s.Scope, "--action-id", "invocation-x", "--decision-id", "decision-1", "--ceiling-micros", "7000000", "--reason", "Explicit conservative ceiling", "--yes"}
	var out, stderr bytes.Buffer
	if rc := runBudgetControl(ctx, args, &out, &stderr, false); rc != 2 {
		t.Fatalf("unattended accepted: %d %s", rc, stderr.String())
	}
	after, _ := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if !bytes.Equal(before, after) {
		t.Fatal("unattended command changed charge")
	}
	badScope := append([]string(nil), args...)
	badScope[4] = "different"
	if rc := runBudgetControl(ctx, badScope, &out, &stderr, true); rc != 1 {
		t.Fatalf("scope mismatch: %d", rc)
	}
	for i := 0; i < 2; i++ {
		if rc := runBudgetControl(ctx, args, &out, &stderr, true); rc != 0 {
			t.Fatalf("attended/replay failed: %d %s", rc, stderr.String())
		}
	}
	state, err := s.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Entries) != 1 {
		t.Fatal("adjustment changed action count")
	}
	for _, entry := range state.Entries {
		if entry.ActualMicros != nil || entry.ReserveMicros != nil || len(entry.Reconciliations) != 1 {
			t.Fatalf("rewrote unknown/history: %+v", entry)
		}
	}
	if _, err := s.Reserve(ctx, budget.Request{ID: "invocation-x", Kind: budget.Launch}, budget.Limits{}); !errors.Is(err, budget.ErrReserved) {
		t.Fatalf("refunded action: %v", err)
	}
}

func TestBudgetZeroCeilingExplicitlyRetainsUnknownObservation(t *testing.T) {
	s := budget.Store{Dir: t.TempDir(), Scope: "idea"}
	ctx := context.Background()
	if _, err := s.Reserve(ctx, budget.Request{ID: "x", Kind: budget.Launch}, budget.Limits{}); err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	args := []string{"reconcile", "--ledger", s.Dir, "--scope", s.Scope, "--action-id", "x", "--decision-id", "zero", "--ceiling-micros", "0", "--reason", "Operator established no charge", "--yes"}
	if rc := runBudgetControl(ctx, args, &out, &stderr, true); rc != 0 {
		t.Fatalf("reconcile: %d %s", rc, stderr.String())
	}
	if !strings.Contains(stderr.String(), "explicit zero ceiling") || !strings.Contains(stderr.String(), "observed cost remains unknown") {
		t.Fatalf("silent zero coercion: %s", stderr.String())
	}
	state, err := s.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range state.Entries {
		if r.ActualMicros != nil || r.ReserveMicros != nil || len(r.Reconciliations) != 1 || r.Reconciliations[0].CeilingMicros != 0 {
			t.Fatalf("observation rewritten: %+v", r)
		}
	}
}

func TestBudgetInspectDoesNotInitializeState(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "missing")
	var out, stderr bytes.Buffer
	if rc := runBudgetControl(context.Background(), []string{"inspect", "--ledger", dir, "--scope", "idea"}, &out, &stderr, false); rc != 1 {
		t.Fatalf("missing ledger: %d", rc)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("inspection initialized budget state")
	}
}
