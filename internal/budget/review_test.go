package budget

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNoOpKernelLockIsRefused(t *testing.T) {
	s := testStore(t)
	release, err := lockWithOps(context.Background(), filepath.Join(s.Dir, "ledger.lock"), func(*os.File) (bool, error) { return true, nil }, func(*os.File) {})
	if release != nil {
		release()
		t.Fatal("no-op locking authorized work")
	}
	if err == nil || !strings.Contains(err.Error(), "does not provide verified exclusion") {
		t.Fatalf("no-op lock: %v", err)
	}
}

func TestOriginBootstrapCannotSplitAcrossCachePaths(t *testing.T) {
	dir := t.TempDir()
	paths := []string{filepath.Join(t.TempDir(), "one.lock"), filepath.Join(t.TempDir(), "two.lock")}
	var wg sync.WaitGroup
	var mu sync.Mutex
	winners := map[int]int{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := pinLockOrigin(dir, paths[i%2])
			if err == nil {
				mu.Lock()
				winners[i%2]++
				mu.Unlock()
			} else if !strings.Contains(err.Error(), "origin mismatch") {
				t.Errorf("origin: %v", err)
			}
		}(i)
	}
	wg.Wait()
	if len(winners) != 1 {
		t.Fatalf("different cache origins admitted: %v", winners)
	}
}

func TestRelativeSymlinkAliasSharesLock(t *testing.T) {
	parent := t.TempDir()
	actual := filepath.Join(parent, "actual")
	alias := filepath.Join(parent, "alias")
	if err := os.Mkdir(actual, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(actual, alias); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	t.Chdir(alias)
	release, err := lock(context.Background(), "ledger.lock")
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if got, err := lock(ctx, filepath.Join(actual, "ledger.lock")); !errors.Is(err, context.DeadlineExceeded) {
		if got != nil {
			got()
		}
		t.Fatalf("alias got independent permission: %v", err)
	}
}

func TestSeparateProcessesCannotExceedActionOrCostCap(t *testing.T) {
	for _, cost := range []bool{false, true} {
		t.Run(fmt.Sprint(cost), func(t *testing.T) {
			s := testStore(t)
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			var children []*exec.Cmd
			for i := 0; i < 12; i++ {
				cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestBudgetReservationProcess$")
				cmd.Env = append(os.Environ(), "PARLEY_BUDGET_RESERVE_DIR="+s.Dir, "PARLEY_BUDGET_RESERVE_ID="+fmt.Sprint(i), "PARLEY_BUDGET_RESERVE_COST="+strconv.FormatBool(cost))
				if err := cmd.Start(); err != nil {
					t.Fatal(err)
				}
				children = append(children, cmd)
			}
			passed := 0
			for _, cmd := range children {
				err := cmd.Wait()
				if err == nil {
					passed++
				} else if cmd.ProcessState.ExitCode() != 23 {
					t.Errorf("child failure: %v", err)
				}
			}
			if passed != 5 {
				t.Fatalf("%d admissions against five-call cap", passed)
			}
			state, err := read(filepath.Join(s.Dir, "ledger.json"))
			if err != nil || len(state.Entries) != 5 {
				t.Fatalf("durable process charges: %d, %v", len(state.Entries), err)
			}
		})
	}
}

func TestBudgetReservationProcess(t *testing.T) {
	dir := os.Getenv("PARLEY_BUDGET_RESERVE_DIR")
	if dir == "" {
		t.Skip("subprocess fixture")
	}
	limits := Limits{Actions: map[Kind]int{Launch: 5}}
	if os.Getenv("PARLEY_BUDGET_RESERVE_COST") == "true" {
		limits = Limits{CostMicros: 15}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := (Store{Dir: dir, Scope: "idea-test"}).Reserve(ctx, Request{ID: os.Getenv("PARLEY_BUDGET_RESERVE_ID"), Kind: Launch, ReserveMicros: micros(3)}, limits)
	if errors.Is(err, ErrLimit) {
		os.Exit(23)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperatorReconciliationRetainsUnknownAndSpentCharge(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	if _, err := s.Reserve(ctx, Request{ID: "unknown", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Settle(ctx, "unknown", nil); err != nil {
		t.Fatal(err)
	}
	state, err := s.ReconcileUnknown(ctx, "unknown", "decision-1", 7, "Explicit operator conservative ceiling")
	if err != nil {
		t.Fatal(err)
	}
	entry := state.Entries[key("unknown")]
	if entry.ActualMicros != nil || entry.ReserveMicros != nil || len(entry.Reconciliations) != 1 {
		t.Fatalf("rewrote observation: %+v", entry)
	}
	if total, err := state.ExposureError(); total != 7 || err != nil {
		t.Fatalf("exposure: %d %v", total, err)
	}
	if _, err := s.ReconcileUnknown(ctx, "unknown", "decision-1", 7, "Explicit operator conservative ceiling"); err != nil {
		t.Fatalf("idempotence: %v", err)
	}
	if _, err := s.ReconcileUnknown(ctx, "unknown", "decision-1", 0, "change decision"); err == nil {
		t.Fatal("conflicting operator replay accepted")
	}
	if _, err := s.Reserve(ctx, Request{ID: "unknown", Kind: Launch}, Limits{}); !errors.Is(err, ErrReserved) {
		t.Fatalf("reconciliation refunded action: %v", err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "next", Kind: Launch, ReserveMicros: micros(3)}, Limits{CostMicros: 10}); err != nil {
		t.Fatalf("bounded recovery: %v", err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "over", Kind: Launch, ReserveMicros: micros(1)}, Limits{CostMicros: 10}); !errors.Is(err, ErrLimit) {
		t.Fatalf("recovered budget overspent: %v", err)
	}
	if _, err := s.ReconcileUnknown(ctx, "next", "decision-known", 0, "not allowed"); err == nil {
		t.Fatal("overwrote known observation")
	}
}

func TestOperatorReconciliationWriteFailureCannotLoseCharge(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	if _, err := s.Reserve(ctx, Request{ID: "x", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	s.persist = func(string, []byte) error { return errors.New("write failure") }
	if _, err := s.ReconcileUnknown(ctx, "x", "decision", 7, "explicit ceiling"); err == nil {
		t.Fatal("failed decision write passed")
	}
	s.persist = nil
	state, err := read(filepath.Join(s.Dir, "ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Entries) != 1 {
		t.Fatal("lost charge")
	}
	if _, err := state.ExposureError(); !errors.Is(err, ErrUnknownCost) {
		t.Fatalf("invented saved ceiling: %v", err)
	}
}

func TestDeniedKindDoesNotMeanUnlimited(t *testing.T) {
	s := testStore(t)
	if _, err := s.Reserve(context.Background(), Request{ID: "x", Kind: Fixup}, Limits{Denied: map[Kind]bool{Fixup: true}}); !errors.Is(err, ErrLimit) {
		t.Fatalf("forbidden action: %v", err)
	}
}

func TestClockAndOverflowAreExplicit(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	now := time.Now().UTC()
	s.now = func() time.Time { return now }
	if _, err := s.Reserve(ctx, Request{ID: "x", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	s.now = func() time.Time { return now.Add(-time.Second) }
	if _, err := s.Settle(ctx, "x", micros(1)); !errors.Is(err, ErrClockSkew) {
		t.Fatalf("clock error: %v", err)
	}
	state := Snapshot{Entries: map[string]Reservation{"a": {ReserveMicros: micros(math.MaxInt64)}, "b": {ReserveMicros: micros(1)}}}
	if _, err := state.ExposureError(); !errors.Is(err, ErrCostOverflow) {
		t.Fatalf("overflow: %v", err)
	}
}

func TestJSONDepthBound(t *testing.T) {
	data := strings.Repeat(`{"nested":`, 40) + "0" + strings.Repeat("}", 40)
	if err := checkJSON(json.NewDecoder(strings.NewReader(data))); err == nil || !strings.Contains(err.Error(), "nesting") {
		t.Fatalf("depth: %v", err)
	}
}
