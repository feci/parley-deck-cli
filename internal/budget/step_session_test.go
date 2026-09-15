package budget

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestStepSessionOneChargeAndExpiredContext(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	b, err := EnsureStepBinding(ctx, root, "idea", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	scoped, close, err := OpenStepSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			child, finish, e := JoinStepSession(scoped, root, "idea")
			if e != nil {
				t.Error(e)
				return
			}
			defer finish()
			if e = ChargeStep(child); e != nil {
				t.Error(e)
			}
		}()
	}
	wg.Wait()
	close()
	s, err := b.Store.Inspect(ctx)
	if err != nil || StepCount(s) != 1 {
		t.Fatalf("nested charges: %+v %v", s, err)
	}
	if ChargeStep(scoped) == nil {
		t.Fatal("retained context grants work after session end")
	}
	if _, _, err := OpenStepSession(ctx, b); !errors.Is(err, ErrLimit) {
		t.Fatalf("new session bypasses cap: %v", err)
	}
}
func TestStepSessionAwaitAndFailedPersistence(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	b, err := EnsureStepBinding(ctx, root, "idea", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		_, close, e := OpenStepSession(ctx, b)
		if e != nil {
			t.Fatal(e)
		}
		close()
	}
	s, _ := b.Store.Inspect(ctx)
	if len(s.Entries) != 0 {
		t.Fatal("reads consumed steps")
	}
	b.Store.persist = func(path string, data []byte) error {
		if err := writeSynced(path, data); err != nil {
			return err
		}
		return errors.New("post-publication fault")
	}
	scoped, close, err := OpenStepSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer close()
	if ChargeStep(scoped) == nil {
		t.Fatal("persistence fault granted work")
	}
	s, err = b.Store.Inspect(ctx)
	if err != nil || StepCount(s) != 1 {
		t.Fatalf("charge was refunded: %+v %v", s, err)
	}
	if ChargeStep(scoped) == nil {
		t.Fatal("fault was cleared within the same action")
	}
}
func TestStepPolicyClockAndChargeContinuity(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	b, err := EnsureStepBinding(ctx, root, "idea", 9, 20*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Millisecond)
	resumed, err := EnsureStepBinding(ctx, root, "idea", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = resumed.Check(ctx); !errors.Is(err, ErrLimit) {
		t.Fatalf("resume reset lifetime clock: %v", err)
	}
	if _, err = EnsureStepBinding(ctx, root, "idea", 99, time.Hour); err == nil {
		t.Fatal("runtime configuration extended saved ceilings")
	}
	if err = os.Remove(filepath.Join(b.Store.Dir, "ledger.json")); err != nil {
		t.Fatal(err)
	}
	if _, err = resumed.Check(ctx); err == nil {
		t.Fatal("lost history became a new ledger")
	}
}

func TestStepSessionChargedTransitionCannotOutliveClockOrLoseScope(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	b, err := EnsureStepBinding(ctx, root, "idea", 1, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	scoped, close, err := OpenStepSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer close()
	if err := ChargeStep(scoped); err != nil {
		t.Fatal(err)
	}
	if _, _, err := JoinStepSession(scoped, root, "different"); err == nil {
		t.Fatal("charged session escaped into an unconfigured idea")
	}
	time.Sleep(1100 * time.Millisecond)
	if err := ChargeStep(scoped); !errors.Is(err, ErrLimit) {
		t.Fatalf("charged group gained a launch after expiry: %v", err)
	}
	state, err := b.Store.Inspect(ctx)
	if err != nil || StepCount(state) != 1 {
		t.Fatalf("expiry refunded/repeated the charge: %+v %v", state, err)
	}
}
