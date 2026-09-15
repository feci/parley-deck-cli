package budget

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type rejectingCycleObserver struct {
	before, after int
	failBefore    bool
}

func (o *rejectingCycleObserver) BeforeCycle(context.Context, CycleBinding, Snapshot) error {
	o.before++
	if o.failBefore {
		return errors.New("fixture pre-observation failed")
	}
	return nil
}
func (o *rejectingCycleObserver) AfterCycle(context.Context, CycleBinding, Snapshot, string) error {
	o.after++
	return errors.New("fixture durable observation publication failed")
}
func TestCycleObserverPublicationFailureNeverGrantsOrRefunds(t *testing.T) {
	for _, before := range []bool{true, false} {
		t.Run(map[bool]string{true: "before-reservation", false: "after-reservation"}[before], func(t *testing.T) {
			root := t.TempDir()
			b, err := EnsureCycleBinding(context.Background(), root, "fixture", Fixup, 5, 0, "", "")
			if err != nil {
				t.Fatal(err)
			}
			err = ActivateCycleTrajectory(context.Background(), root, "fixture", CyclePolicyDigest(b.Policy), strings.Repeat("a", 64), func(CycleBinding, Snapshot) error { return nil })
			if err != nil {
				t.Fatal(err)
			}
			b, err = LoadCycleBinding(context.Background(), root, "fixture", Fixup)
			if err != nil {
				t.Fatal(err)
			}
			o := &rejectingCycleObserver{failBefore: before}
			ctx, finish, err := OpenCycleSession(WithCycleObserver(context.Background(), o), b)
			if err != nil {
				t.Fatal(err)
			}
			defer finish()
			for i := 0; i < 2; i++ {
				if _, err = ChargeCycle(ctx, Fixup); err == nil {
					t.Fatal("failed observation granted execution")
				}
			}
			state, err := b.Store.Inspect(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			want := 1
			if before {
				want = 0
			}
			if b.Count(state) != want || o.before != 1 || o.after != want {
				t.Fatalf("charge or sticky refusal mismatch: count=%d before=%d after=%d", b.Count(state), o.before, o.after)
			}
		})
	}
}
