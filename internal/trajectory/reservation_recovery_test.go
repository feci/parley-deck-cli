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

	"parley-deck-cli/internal/budget"
)

type interruptedReservationObserver struct {
	*Observer
	boundary string
}

func (o *interruptedReservationObserver) PrepareCycleReservation(ctx context.Context, b budget.CycleBinding, s budget.Snapshot, i budget.CycleReservationIntent) error {
	if err := o.Observer.PrepareCycleReservation(ctx, b, s, i); err != nil {
		return err
	}
	if o.boundary == "before-charge" {
		os.Exit(24)
	}
	return nil
}

func (o *interruptedReservationObserver) AfterCycle(context.Context, budget.CycleBinding, budget.Snapshot, string) error {
	os.Exit(23)
	return nil
}

func TestReservationRecoveryProcessHelper(t *testing.T) {
	root := os.Getenv("PARLEY_RESERVATION_TEST_ROOT")
	if root == "" {
		return
	}
	ctx := context.Background()
	b, err := budget.LoadCycleBinding(ctx, root, "fixture", budget.Fixup)
	if err == nil {
		o := &interruptedReservationObserver{&Observer{Root: root}, os.Getenv("PARLEY_RESERVATION_TEST_BOUNDARY")}
		ctx = budget.WithCycleObserver(ctx, o)
		var finish func()
		ctx, finish, err = budget.OpenCycleSession(ctx, b)
		if err == nil {
			defer finish()
			_, err = budget.ChargeCycle(ctx, budget.Fixup)
		}
	}
	fmt.Fprintln(os.Stderr, "expected interrupted boundary was not reached:", err)
	os.Exit(25)
}

func interruptedReservationFixture(t *testing.T, boundary string) (string, *budget.CycleBinding, string) {
	t.Helper()
	root, b, _ := accountingFixture(t)
	cmd := exec.Command(os.Args[0], "-test.run=^TestReservationRecoveryProcessHelper$")
	cmd.Env = append(os.Environ(), "PARLEY_RESERVATION_TEST_ROOT="+root, "PARLEY_RESERVATION_TEST_BOUNDARY="+boundary)
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	want := 23
	if boundary == "before-charge" {
		want = 24
	}
	if !errors.As(err, &exit) || exit.ExitCode() != want {
		t.Fatalf("original process did not reach interrupted boundary: %v %s", err, out)
	}
	names, err := os.ReadDir(filepath.Join(filepath.Dir(b.Store.Dir), "reservation-intents"))
	if err != nil || len(names) != 1 {
		t.Fatalf("original intent unavailable: %v %v", names, err)
	}
	return root, b, strings.TrimSuffix(names[0].Name(), ".json")
}

func TestReservationRecoveryActualCrashPreservesChargeAndUnresolvedExecution(t *testing.T) {
	ctx := context.Background()
	for _, boundary := range []string{"before-charge", "after-charge"} {
		t.Run(boundary, func(t *testing.T) {
			root, b, entry := interruptedReservationFixture(t, boundary)
			ledgerBytes := snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))
			original := snapshotRead(t, statePath(*b))
			if boundary == "after-charge" {
				if _, err := Inspect(ctx, root, "fixture"); err == nil {
					t.Fatal("orphan charge admitted by normal reads")
				}
			}
			// Recovery must use the original retained archive, not today's worktree.
			snapshotWrite(t, root, "source", []byte("later unrelated dirty source\n"), 0600)
			p, err := PreviewReservationRecovery(ctx, root, "fixture", entry)
			if err != nil {
				t.Fatal(err)
			}
			if p.Permission != "none" || p.ExecutionStatus != "not-established" {
				t.Fatal("recovery invented execution authority")
			}
			if !bytes.Equal(original, snapshotRead(t, statePath(*b))) {
				t.Fatal("preview mutated state")
			}
			if boundary == "before-charge" {
				if p.Status != "intent-without-published-charge" || p.Attempt != nil {
					t.Fatal("unspent intent became a charge")
				}
				if _, err = RecoverReservation(ctx, root, "fixture", entry, p.SHA256()); err == nil {
					t.Fatal("absent charge permitted execution or mutation")
				}
				if !bytes.Equal(original, snapshotRead(t, statePath(*b))) {
					t.Fatal("uncharged intent changed state")
				}
			} else {
				if p.Attempt == nil || p.Attempt.Launch != nil || p.Attempt.Terminal != nil || p.Attempt.After != nil {
					t.Fatal("recovery guessed original execution")
				}
				if _, err = RecoverReservation(ctx, root, "fixture", entry, strings.Repeat("f", 64)); err == nil {
					t.Fatal("wrong preview accepted")
				}
				if _, err = RecoverReservation(ctx, root, "fixture", entry, p.SHA256()); err != nil {
					t.Fatal(err)
				}
				recovered := snapshotRead(t, statePath(*b))
				for j := 0; j < 2; j++ {
					q, err := RecoverReservation(ctx, root, "fixture", entry, p.SHA256())
					if err != nil || q.SHA256() != p.SHA256() || !bytes.Equal(recovered, snapshotRead(t, statePath(*b))) {
						t.Fatal("exact replay rewrote original history", err)
					}
				}
				s, err := Inspect(ctx, root, "fixture")
				if err != nil || len(s.Attempts) != 1 || s.Attempts[0].Before.Tree.SHA256 != s.Policy.Baseline.Tree.SHA256 {
					t.Fatal("original source was not preserved", err)
				}
				if err = RequireResolved(ctx, root, "fixture"); err == nil {
					t.Fatal("restored accounting closed unresolved execution")
				}
				if _, err = b.Reserve(budget.WithCycleObserver(ctx, &Observer{Root: root}), "fresh-after-recovery"); err == nil {
					t.Fatal("restored accounting authorized another attempt")
				}
				if _, err = Begin(ctx, root, "fixture", "builder", "invented-launch"); err == nil {
					t.Fatal("recovery created a live cycle session")
				}
			}
			if !bytes.Equal(ledgerBytes, snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))) {
				t.Fatal("recovery refunded or repeated the original charge")
			}
			if got := snapshotRead(t, filepath.Join(root, "source")); string(got) != "later unrelated dirty source\n" {
				t.Fatal("recovery rewrote current source")
			}
		})
	}
}

func TestReservationRecoveryRefusesMissingChangedAndPartialEvidence(t *testing.T) {
	ctx := context.Background()
	for _, change := range []string{"missing-intent", "partial-intent", "changed-root", "changed-before", "changed-limits", "changed-action", "missing-archive", "extra-charge", "changed-trajectory", "symlink-intent"} {
		t.Run(change, func(t *testing.T) {
			root, b, entry := interruptedReservationFixture(t, "after-charge")
			path := filepath.Join(filepath.Dir(b.Store.Dir), intentName(entry))
			raw := snapshotRead(t, path)
			var i reservationIntent
			if err := json.Unmarshal(raw, &i); err != nil {
				t.Fatal(err)
			}
			switch change {
			case "missing-intent":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "partial-intent":
				snapshotWrite(t, filepath.Dir(path), filepath.Base(path), []byte("{\n"), 0600)
			case "changed-root":
				i.Root += "-other"
			case "changed-before":
				i.Before.Policy.Baseline.Tree.SHA256 = digest([]byte("other"))
				bts, _ := canonical(i.Before)
				i.BeforeSHA256 = digest(bts)
			case "changed-limits":
				i.Accounting.Policy.Maximum++
			case "changed-action":
				i.Accounting.Action.Input.InputSHA256 = digest([]byte("other"))
			case "missing-archive":
				if err := os.RemoveAll(snapshotDirectory(*b)); err != nil {
					t.Fatal(err)
				}
			case "extra-charge":
				zero := int64(0)
				if _, err := b.Store.Reserve(ctx, budget.Request{ID: "unexpected", Kind: budget.Fixup, ReserveMicros: &zero}, budget.Limits{}); err != nil {
					t.Fatal(err)
				}
			case "changed-trajectory":
				s, _, err := readState(statePath(*b))
				if err != nil {
					t.Fatal(err)
				}
				s.Policy.Implementer = "other"
				if err = writeState(statePath(*b), s); err != nil {
					t.Fatal(err)
				}
			case "symlink-intent":
				target := filepath.Join(t.TempDir(), "retained.json")
				if err := os.WriteFile(target, raw, 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, path); err != nil {
					t.Skip(err)
				}
			}
			if strings.HasPrefix(change, "changed-") && change != "changed-trajectory" {
				data, _ := canonical(i)
				snapshotWrite(t, filepath.Dir(path), filepath.Base(path), data, 0600)
			}
			before := snapshotRead(t, statePath(*b))
			ledger := snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))
			if _, err := PreviewReservationRecovery(ctx, root, "fixture", entry); err == nil {
				t.Fatal("changed original evidence accepted")
			}
			if !bytes.Equal(before, snapshotRead(t, statePath(*b))) || !bytes.Equal(ledger, snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))) {
				t.Fatal("refusal changed history")
			}
		})
	}
}

func TestReservationRecoveryPublicationFailureAndLostOutputReplay(t *testing.T) {
	ctx := context.Background()
	for _, published := range []bool{false, true} {
		t.Run(fmt.Sprint(published), func(t *testing.T) {
			root, b, entry := interruptedReservationFixture(t, "after-charge")
			p, err := PreviewReservationRecovery(ctx, root, "fixture", entry)
			if err != nil {
				t.Fatal(err)
			}
			fault := errors.New("injected original state publication failure")
			_, err = recoverReservation(ctx, root, "fixture", entry, p.SHA256(), true, func(path string, s State) error {
				if published {
					if err := writeState(path, s); err != nil {
						return err
					}
				}
				return fault
			})
			if !errors.Is(err, fault) {
				t.Fatal("publication failure hidden", err)
			}
			before := snapshotRead(t, statePath(*b))
			calls := 0
			q, err := recoverReservation(ctx, root, "fixture", entry, p.SHA256(), true, func(path string, s State) error { calls++; return writeState(path, s) })
			if err != nil || q.SHA256() != p.SHA256() {
				t.Fatal("original preview could not replay", err)
			}
			if published && (calls != 0 || !bytes.Equal(before, snapshotRead(t, statePath(*b)))) {
				t.Fatal("published original state was rewritten")
			}
			if !published && calls != 1 {
				t.Fatal("missing state was not published exactly once")
			}
		})
	}
}

func TestReservationRecoveryRetainsRequiredIntentAndPolicyAfterExtension(t *testing.T) {
	ctx := context.Background()
	root, b, entry := interruptedReservationFixture(t, "after-charge")
	p, err := PreviewReservationRecovery(ctx, root, "fixture", entry)
	if err != nil {
		t.Fatal(err)
	}
	_, err = budget.ExtendCycleBudget(ctx, root, "fixture", budget.Fixup, budget.CycleExtensionRequest{DecisionID: "later-extension", ExpectedPolicySHA256: budget.CyclePolicyDigest(b.Policy), Reason: "fixture later budget", Maximum: 7})
	if err != nil {
		t.Fatal(err)
	}
	q, err := RecoverReservation(ctx, root, "fixture", entry, p.SHA256())
	if err != nil || q.SHA256() != p.SHA256() {
		t.Fatal("later extension rewrote original recovery", err)
	}
	s, err := Inspect(ctx, root, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	s.Attempts[0].ReservationIntentSHA256 = ""
	if err = writeState(statePath(*b), *s); err != nil {
		t.Fatal(err)
	}
	if _, err = Inspect(ctx, root, "fixture"); err == nil {
		t.Fatal("new reservation was silently downgraded to legacy")
	}
}
