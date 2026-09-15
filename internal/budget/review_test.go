package budget

import (
	"bytes"
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

func testPinnedLockPath(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "lock-origin"))
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(string(data), "\n")
	if len(parts) != 5 || parts[0] != "parley-budget-lock/v2" {
		t.Fatalf("origin: %q", data)
	}
	return parts[2]
}

func TestUnknownSettlementRetainsConservativeReserve(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	limits := Limits{CostMicros: 10, RequireKnownCost: true}
	if _, err := s.Reserve(ctx, Request{ID: "first", Kind: Launch, ReserveMicros: micros(7)}, limits); err != nil {
		t.Fatal(err)
	}
	state, err := s.Settle(ctx, "first", nil)
	if err != nil {
		t.Fatal(err)
	}
	entry := state.Entries[key("first")]
	if !entry.Settled || entry.ActualMicros != nil || *entry.ReserveMicros != 7 {
		t.Fatalf("observation changed: %+v", entry)
	}
	if exposure, err := state.ExposureError(); exposure != 7 || err != nil {
		t.Fatalf("lost conservative reserve: %d %v", exposure, err)
	}
	before, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.ReconcileUnknown(ctx, "first", "lower", 0, "Attempt to lower an existing bound"); err == nil {
		t.Fatal("operator silently lowered retained reserve")
	}
	after, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("refusal changed charges: %v", err)
	}
	state, err = s.ReconcileUnknown(ctx, "first", "higher", 8, "Explicit larger conservative bound")
	if err != nil {
		t.Fatal(err)
	}
	if exposure, err := state.ExposureError(); exposure != 8 || err != nil {
		t.Fatalf("larger bound ignored: %d %v", exposure, err)
	}
	// Legacy persisted reconciliations cannot undercut the original reservation.
	legacy := state.Entries[key("first")]
	legacy.Reconciliations[len(legacy.Reconciliations)-1].CeilingMicros = 0
	state.Entries[key("first")] = legacy
	if exposure, err := state.ExposureError(); exposure != 7 || err != nil {
		t.Fatalf("legacy ceiling undercut reserve: %d %v", exposure, err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "remaining", Kind: Launch, ReserveMicros: micros(2)}, limits); err != nil {
		t.Fatalf("known remaining allowance refused: %v", err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "over", Kind: Launch, ReserveMicros: micros(1)}, limits); !errors.Is(err, ErrLimit) {
		t.Fatalf("unknown actual allowed overspend: %v", err)
	}
}

func TestOriginDiagnosticsPrecedeMissingLocalIdentity(t *testing.T) {
	for _, tc := range []struct {
		name        string
		field       int
		value, want string
	}{
		{"old-version", 0, "parley-budget-lock/v1", "unsupported or malformed origin version"},
		{"other-host", 1, "another-budget-test-host.invalid", "hostname changed"},
		{"relocated", 2, "/different/cache/or/ledger/location", "cache path or ledger location changed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := testStore(t)
			ctx := context.Background()
			if _, err := s.Reserve(ctx, Request{ID: "spent", Kind: Launch}, Limits{}); err != nil {
				t.Fatal(err)
			}
			lockPath := testPinnedLockPath(t, s.Dir)
			if err := os.Remove(lockPath); err != nil {
				t.Fatal(err)
			}
			origin := filepath.Join(s.Dir, "lock-origin")
			data, err := os.ReadFile(origin)
			if err != nil {
				t.Fatal(err)
			}
			parts := strings.Split(string(data), "\n")
			parts[tc.field] = tc.value
			if err := os.WriteFile(origin, []byte(strings.Join(parts, "\n")), 0o600); err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := s.Reserve(ctx, Request{ID: "next", Kind: Launch}, Limits{}); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("misdirected origin diagnostic: %v", err)
			}
			if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
				t.Fatalf("recreated local identity: %v", err)
			}
			after, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("refusal changed charges: %v", err)
			}
		})
	}
}

func TestSnapshotRetryIsBoundedSelectiveAndCancellable(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	if _, err := s.Reserve(ctx, Request{ID: "spent", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Dir, "ledger.json")
	calls := 0
	state, err := readStableSnapshot(ctx, path, func(p string) (Snapshot, error) {
		calls++
		if calls <= 2 {
			return Snapshot{}, fmt.Errorf("concurrent publish: %w", ErrSnapshotChanged)
		}
		return read(p)
	})
	if err != nil || calls != 3 || len(state.Entries) != 1 {
		t.Fatalf("retry: %d %+v %v", calls, state, err)
	}
	calls = 0
	_, err = readStableSnapshot(ctx, path, func(string) (Snapshot, error) { calls++; return Snapshot{}, ErrSnapshotChanged })
	if !errors.Is(err, ErrSnapshotChanged) || calls != 5 {
		t.Fatalf("unbounded read: %d %v", calls, err)
	}
	calls = 0
	corrupt := errors.New("malformed ledger")
	_, err = readStableSnapshot(ctx, path, func(string) (Snapshot, error) { calls++; return Snapshot{}, corrupt })
	if !errors.Is(err, corrupt) || calls != 1 {
		t.Fatalf("corruption retried: %d %v", calls, err)
	}
	calls = 0
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, err = readStableSnapshot(cancelCtx, path, func(string) (Snapshot, error) { calls++; cancel(); return Snapshot{}, ErrSnapshotChanged })
	if !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatalf("cancel ignored: %d %v", calls, err)
	}
}

func TestMissingEstablishedLockIsNotRecreatedWhileHeld(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	if _, err := s.Reserve(ctx, Request{ID: "spent", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	release, err := lock(ctx, filepath.Join(s.Dir, "ledger.lock"))
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	path := testPinnedLockPath(t, s.Dir)
	if err := os.Remove(path); err != nil {
		t.Skipf("platform already prevents removing the held lock: %v", err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "next", Kind: Launch}, Limits{}); err == nil || !strings.Contains(err.Error(), "refusing recreation") {
		t.Fatalf("deleted live lock authorized work: %v", err)
	}
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("lock recreated: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("charge changed: %v", err)
	}
}

func TestReplacedLockIdentityRefusesBeforeReservation(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	if _, err := s.Reserve(ctx, Request{ID: "spent", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	path := testPinnedLockPath(t, s.Dir)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err := lockIdentity(path, true); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "next", Kind: Launch}, Limits{}); err == nil || !strings.Contains(err.Error(), "identity changed") {
		t.Fatalf("foreign identity admitted: %v", err)
	}
	state, err := s.Inspect(ctx)
	if err != nil || len(state.Entries) != 1 {
		t.Fatalf("charges: %+v %v", state, err)
	}
}

func TestLockDescriptorCannotSwitchIdentityDuringAcquisition(t *testing.T) {
	s := testStore(t)
	changed := false
	take := func(f *os.File) (bool, error) {
		if !changed {
			changed = true
			// Preserve the inode while changing the token between identity read
			// and lock acquisition, so an inode-only check cannot catch it.
			if _, err := f.WriteAt([]byte(strings.Repeat("0", 64)+"\n"), 0); err != nil {
				return false, err
			}
		}
		return tryLock(f)
	}
	release, err := lockWithOps(context.Background(), filepath.Join(s.Dir, "ledger.lock"), take, unlock)
	if release != nil {
		release()
		t.Fatal("changed descriptor authorized work")
	}
	if err == nil || !strings.Contains(err.Error(), "identity changed") {
		t.Fatalf("descriptor identity: %v", err)
	}
}

func TestOriginlessLedgerInspectionNeverPinsOrLosesCharges(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	if _, err := s.Reserve(ctx, Request{ID: "spent", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(s.Dir, "lock-origin")); err != nil {
		t.Fatal(err)
	}
	for _, scope := range []string{s.Scope, "wrong-scope"} {
		reader := s
		reader.Scope = scope
		state, err := reader.Inspect(ctx)
		if scope == s.Scope && (err != nil || len(state.Entries) != 1) {
			t.Fatalf("inspection: %+v %v", state, err)
		}
		if scope != s.Scope && err == nil {
			t.Fatal("wrong scope accepted")
		}
		files, err := os.ReadDir(s.Dir)
		if err != nil || len(files) != 2 || files[0].Name() != "ledger-established" || files[1].Name() != "ledger.json" {
			t.Fatalf("read changed directory: %v %v", files, err)
		}
	}
	if _, err := s.Reserve(ctx, Request{ID: "next", Kind: Launch}, Limits{}); err == nil || !strings.Contains(err.Error(), "has no lock origin") {
		t.Fatalf("missing origin silently bootstrapped: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("lost historical charge: %v", err)
	}
}

func TestRequireKnownCostPreventsUnknownReservationWithoutCap(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	limits := Limits{RequireKnownCost: true}
	if _, err := s.Reserve(ctx, Request{ID: "unknown", Kind: Launch}, limits); !errors.Is(err, ErrUnknownCost) {
		t.Fatalf("unknown reservation: %v", err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "known", Kind: Launch, ReserveMicros: micros(7)}, limits); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Reserve(ctx, Request{ID: "known", Kind: Launch}, limits); !errors.Is(err, ErrReserved) {
		t.Fatalf("replay priority lost: %v", err)
	}
}

func TestContentionNamesBudgetAndPreservesDeadline(t *testing.T) {
	s := testStore(t)
	release, err := lock(context.Background(), filepath.Join(s.Dir, "ledger.lock"))
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = s.Reserve(ctx, Request{ID: "next", Kind: Launch}, Limits{})
	if !errors.Is(err, ErrLockContention) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("ambiguous contention: %v", err)
	}
}

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
			err := pinLockOrigin(dir, paths[i%2], strings.Repeat(strconv.Itoa(i%2), 64))
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
			outputs := make(map[*exec.Cmd]*bytes.Buffer)
			for i := 0; i < 12; i++ {
				cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestBudgetReservationProcess$")
				cmd.Env = append(os.Environ(), "PARLEY_BUDGET_RESERVE_DIR="+s.Dir, "PARLEY_BUDGET_RESERVE_ID="+fmt.Sprint(i), "PARLEY_BUDGET_RESERVE_COST="+strconv.FormatBool(cost))
				output := &bytes.Buffer{}
				cmd.Stdout, cmd.Stderr = output, output
				outputs[cmd] = output
				if err := cmd.Start(); err != nil {
					t.Fatal(err)
				}
				children = append(children, cmd)
			}
			passed := 0
			unexpected := 0
			for _, cmd := range children {
				err := cmd.Wait()
				if err == nil {
					passed++
				} else if cmd.ProcessState.ExitCode() != 23 {
					unexpected++
					t.Errorf("reservation fixture process failed (not a cap refusal): %v; output=%s", err, outputs[cmd].String())
				}
			}
			if unexpected > 0 {
				t.Fatalf("%d process/acquisition failures prevent evaluating the admission count", unexpected)
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
