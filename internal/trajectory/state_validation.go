package trajectory

import (
	"context"
	"errors"

	"parley-deck-cli/internal/budget"
)

// Full historical reads release the cycle guard needed by live terminal
// publication and ticket stop. Only an unchanged, fully validated authority may
// reach the caller callback, which still holds the guard. One authority drift
// may restart all reads before the callback; evidence and callback errors are
// never retried. This preserves exact concurrent replay without repeating work
// that can charge, launch, consume a ticket or publish a caller result.
func withValidatedState(ctx context.Context, root, idea string, check func(context.Context, budget.CycleBinding, State) error, fn func(budget.CycleBinding, budget.Snapshot, State) error) error {
	for attempt := 0; ; attempt++ {
		var originalBinding budget.CycleBinding
		var originalLedger budget.Snapshot
		var originalState State
		found := false
		err := withStateControl(ctx, root, idea, func(b budget.CycleBinding, ledger budget.Snapshot, s State) error {
			found = true
			originalBinding, originalLedger, originalState = b, ledger, s
			return nil
		})
		if err != nil {
			return err
		}
		if !found {
			if attempt == 0 {
				return nil
			}
			return errors.New("trajectory authority disappeared during full evidence validation")
		}
		if err = checkSourceSnapshots(ctx, originalBinding, originalState); err != nil {
			return err
		}
		if err = check(ctx, originalBinding, originalState); err != nil {
			return err
		}
		rechecked, changed := false, false
		err = withStateControl(ctx, root, idea, func(b budget.CycleBinding, ledger budget.Snapshot, s State) error {
			rechecked = true
			if b.Store.Dir != originalBinding.Store.Dir || b.Store.Scope != originalBinding.Store.Scope || !sameJSON(b.Policy, originalBinding.Policy) || !sameJSON(ledger, originalLedger) || !sameJSON(s, originalState) {
				changed = true
				return errors.New("trajectory authority changed during full evidence validation")
			}
			return fn(b, ledger, s)
		})
		if changed && attempt == 0 {
			continue
		}
		if err == nil && !rechecked {
			return errors.New("trajectory authority disappeared during full evidence validation")
		}
		return err
	}
}
