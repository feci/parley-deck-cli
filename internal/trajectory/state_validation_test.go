package trajectory

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/telemetry"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStateValidationDoesNotHoldControlGuard(t *testing.T) {
	root, b := unchangedFixture(t, "true")
	charged := chargeFixture(t, root, b)
	inv, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{RunID: "n2-control", Idea: "fixture", Phase: "fixup", Agent: "builder", LaunchMode: "headless"})
	if err != nil {
		t.Fatal(err)
	}
	run, err := Begin(charged, root, "fixture", "builder", inv.ID)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", "-c", "printf 'changed\n' > source; exit 7")
	cmd.Dir = root
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err = inv.Started(cmd.Process.Pid); err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		t.Fatal(err)
	}
	var exit *exec.ExitError
	if err = cmd.Wait(); !errors.As(err, &exit) || exit.ExitCode() != 7 {
		t.Fatal("actual fixture exit missing", err)
	}
	code := 7
	if err = inv.Finish(telemetry.Outcome{Status: "failed", ExitCode: &code, FailureClass: telemetry.String("process_failure")}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	entered, release := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	called := false
	checks := 0
	go func() {
		done <- withStateResolutionCheck(ctx, root, "fixture", func(ctx context.Context, b budget.CycleBinding, s State) error {
			if err := checkResolutions(ctx, b, s); err != nil {
				return err
			}
			checks++
			if checks == 2 {
				return nil
			}
			close(entered)
			select {
			case <-release:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(_ budget.CycleBinding, _ budget.Snapshot, s State) error {
			called = true
			if s.Attempts[0].Terminal == nil || checks != 2 {
				return errors.New("stale authority reached callback")
			}
			return nil
		})
	}()
	select {
	case <-entered:
	case err := <-done:
		t.Fatal("reader never reached pause", err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	finishCtx, stop := context.WithTimeout(ctx, 2*time.Second)
	finishErr := run.Finish(finishCtx, "failed", &code)
	stop()
	close(release)
	readerErr := <-done
	if finishErr != nil {
		t.Fatal("full-history reader blocked terminal publication", finishErr)
	}
	if readerErr != nil || !called || checks != 2 {
		t.Fatal("changed launch state was not fully revalidated before callback", readerErr, called, checks)
	}
	state, _, err := readState(statePath(*b))
	if err != nil || state.Attempts[0].Terminal == nil || state.Attempts[0].After == nil || state.Attempts[0].Terminal.ExitCode == nil || *state.Attempts[0].Terminal.ExitCode != 7 {
		t.Fatal("actual terminal was not retained", state, err)
	}
}

func TestStateValidationRechecksPolicyLedgerAndExistence(t *testing.T) {
	for _, mode := range []string{"policy", "ledger", "disappeared"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			root, b := unchangedFixture(t, "true")
			if mode == "ledger" {
				if _, err := b.Reserve(budget.WithCycleObserver(ctx, &Observer{Root: root}), "n2-ledger"); err != nil {
					t.Fatal(err)
				}
			}
			called := false
			checks := 0
			err := withStateResolutionCheck(ctx, root, "fixture", func(ctx context.Context, original budget.CycleBinding, s State) error {
				checks++
				if checks > 1 {
					return checkResolutions(ctx, original, s)
				}
				switch mode {
				case "policy":
					preview, err := budget.InspectCycleBudget(ctx, root, "fixture", budget.Fixup)
					if err != nil {
						return err
					}
					_, err = budget.ExtendCycleBudget(ctx, root, "fixture", budget.Fixup, budget.CycleExtensionRequest{DecisionID: "n2-extension", ExpectedPolicySHA256: preview.PolicySHA256, Maximum: 6, Reason: "Synthetic authority drift fixture"})
					return err
				case "ledger":
					_, err := b.Store.Settle(ctx, "n2-ledger", nil)
					return err
				default:
					return withStateControl(ctx, root, "fixture", func(current budget.CycleBinding, _ budget.Snapshot, _ State) error {
						current.Policy.TrajectorySHA256 = ""
						raw, err := json.MarshalIndent(current.Policy, "", "  ")
						if err != nil {
							return err
						}
						if err = os.WriteFile(filepath.Join(filepath.Dir(current.Store.Dir), "policy.json"), raw, 0600); err != nil {
							return err
						}
						return os.Remove(statePath(current))
					})
				}
			}, func(budget.CycleBinding, budget.Snapshot, State) error { called = true; return nil })
			if mode == "disappeared" {
				if err == nil || !strings.Contains(err.Error(), "authority disappeared during full evidence validation") || called || checks != 1 {
					t.Fatal("disappearing authority reached callback or retried", err, called, checks)
				}
			} else if err != nil || !called || checks != 2 {
				t.Fatal("changed original "+mode+" was not fully revalidated", err, called, checks)
			}
		})
	}
}

func TestStateValidationPreservesFullChecksAndCallbackGuard(t *testing.T) {
	for _, mode := range []string{"valid", "missing-archive", "failed-check", "inactive"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			root, b := unchangedFixture(t, "true")
			if mode == "missing-archive" {
				s, _, err := readState(statePath(*b))
				if err != nil {
					t.Fatal(err)
				}
				if err = os.Remove(snapshotPath(snapshotDirectory(*b), s.BaselineArchive)); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "inactive" {
				root = t.TempDir()
			}
			checked, called := false, false
			checkFailure := errors.New("fixture full check failed")
			err := withStateResolutionCheck(ctx, root, "fixture", func(ctx context.Context, b budget.CycleBinding, s State) error {
				checked = true
				wait, cancel := context.WithTimeout(ctx, time.Second)
				defer cancel()
				release, err := budget.AcquireResourceGuard(wait, filepath.Dir(b.Store.Dir))
				if err != nil {
					return errors.New("full check still holds the control guard")
				}
				release()
				if mode == "failed-check" {
					return checkFailure
				}
				return checkResolutions(ctx, b, s)
			}, func(b budget.CycleBinding, _ budget.Snapshot, _ State) error {
				if called {
					return errors.New("callback repeated")
				}
				called = true
				wait, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
				defer cancel()
				release, err := budget.AcquireResourceGuard(wait, filepath.Dir(b.Store.Dir))
				if err == nil {
					release()
					return errors.New("callback lost required control guard")
				}
				return nil
			})
			switch mode {
			case "valid":
				if err != nil || !checked || !called {
					t.Fatal("valid full read/guard failed", err, checked, called)
				}
			case "inactive":
				if err != nil || checked || called {
					t.Fatal("inactive authority manufactured a callback", err, checked, called)
				}
			case "missing-archive":
				if err == nil || checked || called {
					t.Fatal("missing archive reached callback", err, checked, called)
				}
			case "failed-check":
				if !errors.Is(err, checkFailure) || !checked || called {
					t.Fatal("failed full evidence reached callback", err, checked, called)
				}
			}
		})
	}
}

func TestPendingCycleRefusesBeforeHistoricalContent(t *testing.T) {
	for _, pending := range []bool{true, false} {
		t.Run(map[bool]string{true: "pending", false: "positive"}[pending], func(t *testing.T) {
			ctx := context.Background()
			root, b := unchangedFixture(t, "true")
			if pending {
				chargeFixture(t, root, b)
			}
			state, _, err := readState(statePath(*b))
			if err != nil {
				t.Fatal(err)
			}
			before, err := b.Store.Inspect(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.Remove(snapshotPath(snapshotDirectory(*b), state.BaselineArchive)); err != nil {
				t.Fatal(err)
			}
			_, err = b.Reserve(budget.WithCycleObserver(ctx, &Observer{Root: root}), "must-refuse")
			marker := "archive unavailable"
			if pending {
				marker = "awaits independent reconciliation"
			}
			if err == nil || !strings.Contains(err.Error(), marker) {
				t.Fatal("pending/positive refusal ordering changed", pending, err)
			}
			after, err := b.Store.Inspect(ctx)
			if err != nil || !sameJSON(before, after) {
				t.Fatal("refusal consumed a new charge", after, err)
			}
		})
	}
}
