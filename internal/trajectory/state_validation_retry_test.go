package trajectory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
)

func TestStateValidationRetryIsBoundedAndRevalidatesEvidence(t *testing.T) {
	for _, mode := range []string{"repeated-drift", "archive-removed", "checker-error", "callback-error", "cancelled"} {
		t.Run(mode, func(t *testing.T) {
			root, b := unchangedFixture(t, "true")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			checks, callbacks := 0, 0
			failure := errors.New("fixture failure")
			err := withStateResolutionCheck(ctx, root, "fixture", func(ctx context.Context, current budget.CycleBinding, s State) error {
				checks++
				switch mode {
				case "checker-error":
					return failure
				case "callback-error":
					return nil
				case "cancelled":
					cancel()
					return ctx.Err()
				}
				preview, err := budget.InspectCycleBudget(ctx, root, "fixture", budget.Fixup)
				if err != nil {
					return err
				}
				_, err = budget.ExtendCycleBudget(ctx, root, "fixture", budget.Fixup, budget.CycleExtensionRequest{
					DecisionID: fmt.Sprintf("n2-drift-%d", checks), ExpectedPolicySHA256: preview.PolicySHA256,
					Maximum: 5 + checks, Reason: "Deterministic repeated authority drift fixture",
				})
				if err != nil {
					return err
				}
				if mode == "archive-removed" {
					return os.Remove(snapshotPath(snapshotDirectory(*b), s.BaselineArchive))
				}
				return nil
			}, func(budget.CycleBinding, budget.Snapshot, State) error { callbacks++; return failure })
			switch mode {
			case "repeated-drift":
				if err == nil || !strings.Contains(err.Error(), "authority changed during full evidence validation") || checks != 2 || callbacks != 0 {
					t.Fatal("revalidation was not bounded before callbacks", err, checks, callbacks)
				}
			case "archive-removed":
				if err == nil || !strings.Contains(err.Error(), "archive unavailable") || checks != 1 || callbacks != 0 {
					t.Fatal("retry reused historical evidence", err, checks, callbacks)
				}
			case "checker-error":
				if !errors.Is(err, failure) || checks != 1 || callbacks != 0 {
					t.Fatal("checker error retried", err, checks, callbacks)
				}
			case "callback-error":
				if !errors.Is(err, failure) || checks != 1 || callbacks != 1 {
					t.Fatal("caller callback repeated", err, checks, callbacks)
				}
			case "cancelled":
				if !errors.Is(err, context.Canceled) || checks != 1 || callbacks != 0 {
					t.Fatal("cancelled validation continued", err, checks, callbacks)
				}
			}
		})
	}
}
