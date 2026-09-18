package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
)

// failedAfterObserver keeps the published intent and original charge but loses
// the charged row, as an interrupted AfterCycle does. The raw request ID stays
// known so the fixture can add real ledger observations later.
type failedAfterObserver struct{ *Observer }

func (failedAfterObserver) AfterCycle(context.Context, budget.CycleBinding, budget.Snapshot, string) error {
	return errors.New("fixture lost the charged row after its original charge")
}

func missingRecoveryRow(t *testing.T) (string, *budget.CycleBinding, string) {
	t.Helper()
	ctx := context.Background()
	root, b, _ := accountingFixture(t)
	if _, err := b.Reserve(budget.WithCycleObserver(ctx, failedAfterObserver{&Observer{Root: root}}), "recovery-row"); err == nil {
		t.Fatal("fixture published the charged row")
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil || len(ledger.Entries) != 1 {
		t.Fatalf("original charge unavailable: %v %v", ledger.Entries, err)
	}
	for entry := range ledger.Entries {
		return root, b, entry
	}
	return "", nil, ""
}

func TestReservationRecoveryContentCheckDoesNotBlockLiveFinish(t *testing.T) {
	root, b, _ := accountingFixture(t)
	charged := chargeFixture(t, root, b)
	run, err := Begin(charged, root, "fixture", "builder", "recovery-live")
	if err != nil {
		t.Fatal(err)
	}
	state, _, err := readState(statePath(*b))
	if err != nil || len(state.Attempts) != 1 || state.Attempts[0].ReservationIntentSHA256 == "" || state.Attempts[0].Launch == nil {
		t.Fatal("fixture lacks a live charged row with its original intent", err)
	}
	entry := state.Attempts[0].Charge.EntryKey
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	p, err := PreviewReservationRecovery(ctx, root, "fixture", entry)
	if err != nil || p.Status != "charged-observation" {
		t.Fatal("live charged row has no exact preview", p.Status, err)
	}
	cmd := exec.Command("sh", "-c", "printf 'changed\\n' > source; exit 7")
	cmd.Dir = root
	var exit *exec.ExitError
	if err = cmd.Run(); !errors.As(err, &exit) || exit.ExitCode() != 7 {
		t.Fatal("actual fixture exit missing", err)
	}
	code := 7
	type result struct {
		p   ReservationRecoveryPreview
		err error
	}
	entered, release := make(chan struct{}), make(chan struct{})
	done := make(chan result, 1)
	checks, terminal := 0, false
	go func() {
		q, err := recoverReservationChecked(ctx, root, "fixture", entry, "", false, func(ctx context.Context, b budget.CycleBinding, s State) error {
			if err := checkStateSnapshots(ctx, b, s); err != nil {
				return err
			}
			checks++
			if checks > 1 {
				terminal = s.Attempts[0].Terminal != nil
				return nil
			}
			close(entered)
			select {
			case <-release:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}, writeState)
		done <- result{q, err}
	}()
	select {
	case <-entered:
	case r := <-done:
		t.Fatal("recovery preview never reached its content check", r.err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	finishCtx, stop := context.WithTimeout(ctx, 2*time.Second)
	finishErr := run.Finish(finishCtx, "failed", &code)
	stop()
	close(release)
	r := <-done
	if finishErr != nil {
		t.Fatal("recovery preview content check blocked live terminal publication", finishErr)
	}
	if r.err != nil || checks != 2 || !terminal || r.p.SHA256() != p.SHA256() {
		t.Fatal("changed live row was not freshly revalidated before the preview returned", r.err, checks, terminal)
	}
	state, _, err = readState(statePath(*b))
	if err != nil || state.Attempts[0].Terminal == nil || state.Attempts[0].After == nil || state.Attempts[0].Terminal.ExitCode == nil || *state.Attempts[0].Terminal.ExitCode != 7 {
		t.Fatal("actual terminal was not retained", err)
	}
	published := snapshotRead(t, statePath(*b))
	if q, err := RecoverReservation(ctx, root, "fixture", entry, p.SHA256()); err != nil || q.SHA256() != p.SHA256() || !bytes.Equal(published, snapshotRead(t, statePath(*b))) {
		t.Fatal("exact replay rewrote the finished charged row", err)
	}
}

func TestReservationRecoveryConcurrentApplyReplaysAfterContentCheck(t *testing.T) {
	root, b, entry := missingRecoveryRow(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	p, err := PreviewReservationRecovery(ctx, root, "fixture", entry)
	if err != nil || p.Attempt == nil {
		t.Fatal("missing charged row has no exact preview", err)
	}
	entered, release := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	checks, calls := 0, 0
	go func() {
		_, err := recoverReservationChecked(ctx, root, "fixture", entry, p.SHA256(), true, func(ctx context.Context, b budget.CycleBinding, s State) error {
			if err := checkStateSnapshots(ctx, b, s); err != nil {
				return err
			}
			checks++
			if checks > 1 {
				return nil
			}
			close(entered)
			select {
			case <-release:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(path string, s State) error { calls++; return writeState(path, s) })
		done <- err
	}()
	select {
	case <-entered:
	case err := <-done:
		t.Fatal("recovery apply never reached its content check", err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	published := 0
	applyCtx, stop := context.WithTimeout(ctx, 2*time.Second)
	q, applyErr := recoverReservationChecked(applyCtx, root, "fixture", entry, p.SHA256(), true, checkStateSnapshots, func(path string, s State) error {
		wait, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
		defer cancel()
		if release, err := budget.AcquireResourceGuard(wait, filepath.Dir(b.Store.Dir)); err == nil {
			release()
			return errors.New("missing-row publication lost the cycle guard")
		}
		published++
		return writeState(path, s)
	})
	stop()
	close(release)
	firstErr := <-done
	if applyErr != nil || published != 1 || q.SHA256() != p.SHA256() {
		t.Fatal("recovery content check blocked a concurrent exact apply", applyErr, published)
	}
	recovered := snapshotRead(t, statePath(*b))
	if firstErr != nil || checks != 2 || calls != 0 || !bytes.Equal(recovered, snapshotRead(t, statePath(*b))) {
		t.Fatal("paused exact apply did not replay the concurrently published row", firstErr, checks, calls)
	}
	s, err := Inspect(ctx, root, "fixture")
	if err != nil || len(s.Attempts) != 1 || s.Attempts[0].Launch != nil || s.Attempts[0].Terminal != nil {
		t.Fatal("recovered row invented execution or failed full validation", err)
	}
}

func TestReservationRecoveryRechecksAuthorityAfterUnguardedContent(t *testing.T) {
	for _, tc := range []struct {
		mode   string
		repeat bool
	}{{"policy", false}, {"policy", true}, {"ledger", false}, {"ledger", true}, {"intent", true}, {"state-removed", false}, {"intent-removed", false}, {"binding-removed", false}} {
		t.Run(fmt.Sprintf("%s/repeat=%t", tc.mode, tc.repeat), func(t *testing.T) {
			ctx := context.Background()
			root, b, entry := missingRecoveryRow(t)
			p, err := PreviewReservationRecovery(ctx, root, "fixture", entry)
			if err != nil {
				t.Fatal(err)
			}
			original := snapshotRead(t, statePath(*b))
			intent := filepath.Join(filepath.Dir(b.Store.Dir), intentName(entry))
			checks, calls := 0, 0
			q, err := recoverReservationChecked(ctx, root, "fixture", entry, p.SHA256(), true, func(ctx context.Context, current budget.CycleBinding, s State) error {
				if err := checkStateSnapshots(ctx, current, s); err != nil {
					return err
				}
				checks++
				wait, cancel := context.WithTimeout(ctx, time.Second)
				defer cancel()
				release, err := budget.AcquireResourceGuard(wait, filepath.Dir(current.Store.Dir))
				if err != nil {
					return errors.New("recovery content check holds the cycle guard")
				}
				release()
				if checks > 1 && !tc.repeat {
					return nil
				}
				switch tc.mode {
				case "policy":
					preview, err := budget.InspectCycleBudget(ctx, root, "fixture", budget.Fixup)
					if err != nil {
						return err
					}
					_, err = budget.ExtendCycleBudget(ctx, root, "fixture", budget.Fixup, budget.CycleExtensionRequest{
						DecisionID: fmt.Sprintf("recovery-drift-%d", checks), ExpectedPolicySHA256: preview.PolicySHA256,
						Maximum: 5 + checks, Reason: "Deterministic recovery authority drift fixture",
					})
					return err
				case "ledger":
					if checks == 1 {
						_, err := b.Store.Settle(ctx, "recovery-row", nil)
						return err
					}
					_, err := b.Store.ReconcileUnknown(ctx, "recovery-row", "recovery-drift", 4, "Deterministic recovery ledger drift fixture")
					return err
				case "intent":
					raw, err := os.ReadFile(intent)
					if err != nil {
						return err
					}
					var i reservationIntent
					if err = json.Unmarshal(raw, &i); err != nil {
						return err
					}
					i.PreparedAt = i.PreparedAt.Add(-time.Nanosecond)
					if raw, err = canonical(i); err != nil {
						return err
					}
					return os.WriteFile(intent, raw, 0600)
				case "state-removed":
					return os.Remove(statePath(current))
				case "intent-removed":
					return os.Remove(intent)
				default:
					return os.Rename(filepath.Dir(current.Store.Dir), filepath.Dir(current.Store.Dir)+"-moved")
				}
			}, func(path string, s State) error { calls++; return writeState(path, s) })
			switch {
			case strings.HasSuffix(tc.mode, "-removed"):
				if err == nil || strings.Contains(err.Error(), "authority changed") || checks != 1 || calls != 0 {
					t.Fatal("disappearing recovery authority was retried or published", err, checks, calls)
				}
				if _, statErr := os.Lstat(statePath(*b)); statErr == nil && !bytes.Equal(original, snapshotRead(t, statePath(*b))) {
					t.Fatal("refusal changed trajectory state")
				}
			case tc.repeat:
				if err == nil || !strings.Contains(err.Error(), "authority changed during full evidence validation") || checks != 2 || calls != 0 || !bytes.Equal(original, snapshotRead(t, statePath(*b))) {
					t.Fatal("repeated "+tc.mode+" drift was not bounded before publication", err, checks, calls)
				}
			default:
				if err != nil || checks != 2 || calls != 1 || q.SHA256() != p.SHA256() {
					t.Fatal("changed original "+tc.mode+" was not freshly revalidated before publication", err, checks, calls)
				}
				if s, err := Inspect(ctx, root, "fixture"); err != nil || len(s.Attempts) != 1 {
					t.Fatal("published row failed full validation", err)
				}
			}
		})
	}
}
