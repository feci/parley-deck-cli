package budget

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func micros(n int64) *int64        { return &n }
func testStore(t *testing.T) Store { t.Helper(); return Store{Dir: t.TempDir(), Scope: "idea-test"} }

func TestReservationsArePrechargedAndSurviveResume(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	limits := Limits{Actions: map[Kind]int{Fixup: 2}}
	if _, err := s.Reserve(ctx, Request{ID: "attempt-1", Kind: Fixup}, limits); err != nil {
		t.Fatal(err)
	}
	// The operation failed or the caller crashed. A new Store sees that charge.
	s = Store{Dir: s.Dir, Scope: s.Scope}
	if _, err := s.Reserve(ctx, Request{ID: "attempt-1", Kind: Fixup}, limits); !errors.Is(err, ErrReserved) {
		t.Fatalf("replay: %v", err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "attempt-2", Kind: Fixup}, limits); err != nil {
		t.Fatalf("inclusive final attempt: %v", err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "attempt-3", Kind: Fixup}, limits); !errors.Is(err, ErrLimit) {
		t.Fatalf("cap: %v", err)
	}
	// A status action must not be misclassified as a fixup by its caller. Other
	// independently limited kinds do not reset or consume that fixup count.
	if _, err := s.Reserve(ctx, Request{ID: "review-call", Kind: Launch}, limits); err != nil {
		t.Fatal(err)
	}
}

func TestConcurrentReservationsCannotOverspend(t *testing.T) {
	s := testStore(t)
	var wg sync.WaitGroup
	var passed atomic.Int32
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := s.Reserve(context.Background(), Request{ID: fmt.Sprint(i), Kind: Launch, ReserveMicros: micros(3)}, Limits{Actions: map[Kind]int{Launch: 8}, CostMicros: 15})
			if err == nil {
				passed.Add(1)
			} else if !errors.Is(err, ErrLimit) {
				t.Errorf("reserve: %v", err)
			}
		}(i)
	}
	wg.Wait()
	if passed.Load() != 5 {
		t.Fatalf("accepted %d against cost cap of five", passed.Load())
	}
	state, err := read(filepath.Join(s.Dir, "ledger.json"))
	if err != nil || len(state.Entries) != 5 {
		t.Fatalf("persisted: %d, %v", len(state.Entries), err)
	}
}

func TestUnknownExposureAndSettlement(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	if _, err := s.Reserve(ctx, Request{ID: "unknown", Kind: Launch}, Limits{CostMicros: 10}); !errors.Is(err, ErrUnknownCost) {
		t.Fatalf("unknown new call: %v", err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "known", Kind: Launch, ReserveMicros: micros(8)}, Limits{CostMicros: 10}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "would-overrun", Kind: Launch, ReserveMicros: micros(3)}, Limits{CostMicros: 10}); !errors.Is(err, ErrLimit) {
		t.Fatal(err)
	}
	if _, err := s.Settle(ctx, "known", micros(2)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Settle(ctx, "known", micros(2)); err != nil {
		t.Fatalf("same settlement: %v", err)
	}
	if _, err := s.Settle(ctx, "known", micros(1)); err == nil {
		t.Fatal("conflicting settlement accepted")
	}
	if _, err := s.Reserve(ctx, Request{ID: "next", Kind: Launch, ReserveMicros: micros(8)}, Limits{CostMicros: 10}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Settle(ctx, "next", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "after-unknown", Kind: Launch, ReserveMicros: micros(1)}, Limits{CostMicros: 10}); !errors.Is(err, ErrLimit) {
		t.Fatalf("retained reservation must exhaust the cap: %v", err)
	}
}

func TestPersistFailureNeverAuthorizesWork(t *testing.T) {
	for _, replaced := range []bool{false, true} {
		t.Run(fmt.Sprint(replaced), func(t *testing.T) {
			s := testStore(t)
			s.persist = func(path string, data []byte) error {
				if replaced {
					if err := writeSynced(path, data); err != nil {
						return err
					}
				}
				return errors.New("injected persistence failure")
			}
			if _, err := s.Reserve(context.Background(), Request{ID: "one", Kind: Fixup}, Limits{}); err == nil {
				t.Fatal("persistence failure authorized work")
			}
			s.persist = nil
			_, err := s.Reserve(context.Background(), Request{ID: "one", Kind: Fixup}, Limits{})
			if replaced && !errors.Is(err, ErrReserved) {
				t.Fatalf("replacement lost conservative charge: %v", err)
			}
			if !replaced && err != nil {
				t.Fatalf("unpublished reservation: %v", err)
			}
		})
	}
}

func TestWallClockDoesNotResetOnResume(t *testing.T) {
	s := testStore(t)
	start := time.Now().UTC()
	s.now = func() time.Time { return start }
	if _, err := s.Reserve(context.Background(), Request{ID: "one", Kind: DriverStep}, Limits{WallClock: time.Second}); err != nil {
		t.Fatal(err)
	}
	s = Store{Dir: s.Dir, Scope: s.Scope, now: func() time.Time { return start.Add(time.Second) }}
	if _, err := s.Reserve(context.Background(), Request{ID: "two", Kind: DriverStep}, Limits{WallClock: time.Second}); !errors.Is(err, ErrLimit) {
		t.Fatalf("resume reset wall clock: %v", err)
	}
}

func TestCorruptPersistedChargeFailsClosed(t *testing.T) {
	s := testStore(t)
	if _, err := s.Reserve(context.Background(), Request{ID: "one", Kind: Fixup, ReserveMicros: micros(2)}, Limits{}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Dir, "ledger.json")
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cases := []string{"null", "{}", string(original) + "{}", strings.Replace(string(original), `"schema": 1`, `"schema": 1, "schema": 1`, 1), strings.Replace(string(original), `"schema": 1`, `"Schema": 1`, 1), strings.Replace(string(original), `"reserve_micros": 2`, `"reserve_micros": -2`, 1), strings.Replace(string(original), `"settled": false`, `"settled": null`, 1), strings.Replace(string(original), `"schema": 1`, `"schema": 1, "unknown": 0`, 1)}
	for i, data := range cases {
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Reserve(context.Background(), Request{ID: "two", Kind: Fixup}, Limits{}); err == nil {
			t.Errorf("corruption %d accepted", i)
		}
	}
}

func TestKernelLockSurvivesProcessDeath(t *testing.T) {
	s := testStore(t)
	cmd := exec.Command(os.Args[0], "-test.run=^TestBudgetLockChild$")
	cmd.Env = append(os.Environ(), "PARLEY_BUDGET_LOCK_CHILD="+s.Dir)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(s.Dir, "ready")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("child did not acquire lock")
		}
		time.Sleep(10 * time.Millisecond)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if _, err := s.Reserve(ctx, Request{ID: "one", Kind: Launch}, Limits{}); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("live lock did not exclude parent: %v", err)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = cmd.Wait()
	if _, err := s.Reserve(context.Background(), Request{ID: "one", Kind: Launch}, Limits{}); err != nil {
		t.Fatalf("kernel did not release crashed owner: %v", err)
	}
}

func TestBudgetLockChild(t *testing.T) {
	dir := os.Getenv("PARLEY_BUDGET_LOCK_CHILD")
	if dir == "" {
		t.Skip("subprocess fixture")
	}
	release, err := lock(context.Background(), filepath.Join(dir, "ledger.lock"))
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if err := os.WriteFile(filepath.Join(dir, "ready"), []byte("ready"), 0o600); err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Second)
}
