package trajectory

import (
	"context"
	"errors"
	"os"
	"slices"

	"parley-deck-cli/internal/budget"
)

// CheckHelperScope refuses a known scope mismatch before a caller consumes a
// ticket or executes a helper. It reuses the reconciliation parser and original
// activation authority. It grants no execution and does not fence later edits;
// execution and final reconciliation must retain their own checks.
func CheckHelperScope(ctx context.Context, req HelperRequest) error {
	if req.Version != 1 {
		return errors.New("unsupported trajectory helper request")
	}
	if _, err := req.Ticket.SHA256(); err != nil {
		return err
	}
	root, err := canonicalRoot(req.Ticket.Root)
	if err != nil || root != req.Ticket.Root {
		return errors.New("trajectory helper origin changed")
	}
	found := false
	err = withState(ctx, root, req.Ticket.Request.Idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		found = true
		r, err := capturedRequestAt(s, req.Ticket.Request.Verifier, req.Ticket.Request.Sequence)
		if err != nil || !sameJSON(r, req.Ticket.Request) {
			return errors.New("trajectory helper differs from its original charged patch")
		}
		members, err := checkCapturedActivationQuorum(ctx, b, s, r)
		if err != nil {
			return err
		}
		if !slices.Equal(req.Participants, members) {
			return errors.New("helper request differs from the original activation quorum")
		}
		dir, err := os.OpenRoot(root)
		if err != nil {
			return err
		}
		defer dir.Close()
		return checkReconciliationScope(dir, req)
	})
	if err == nil && !found {
		return errors.New("trajectory helper authority disappeared")
	}
	return err
}
