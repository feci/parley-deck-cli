package budget

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func cycleGrant(t *testing.T, b *CycleBinding, id string, maximum int) CycleExtensionRequest {
	t.Helper()
	status, err := b.Inspect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return CycleExtensionRequest{DecisionID: id, ExpectedPolicySHA256: status.PolicySHA256, Reason: "Explicit finite fixture grant", Maximum: maximum}
}

func TestCycleExtensionPreservesChargesClockAndDecisionReplay(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 5, 4, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Reserve(ctx, "fifth"); err != nil {
		t.Fatal(err)
	}
	ledger := filepath.Join(b.Store.Dir, "ledger.json")
	before, err := os.ReadFile(ledger)
	if err != nil {
		t.Fatal(err)
	}
	original, err := b.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	req := cycleGrant(t, b, "grant-six", 6)
	status, err := ExtendCycleBudget(ctx, root, "idea", Fixup, req)
	if err != nil {
		t.Fatal(err)
	}
	if status.Spent != 5 || !status.StartedAt.Equal(original.StartedAt) || status.Policy.Maximum != 6 || status.Policy.InitialMaximum() != 5 || status.Policy.Version != 2 || len(status.Policy.Extensions) != 1 {
		t.Fatalf("grant reset state: %+v", status)
	}
	after, err := os.ReadFile(ledger)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("extension rewrote charged ledger: %v", err)
	}
	policy := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	published, err := os.ReadFile(policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ExtendCycleBudget(ctx, root, "idea", Fixup, req); err != nil {
		t.Fatal(err)
	}
	replayed, _ := os.ReadFile(policy)
	if !bytes.Equal(published, replayed) {
		t.Fatal("exact replay rewrote policy")
	}
	conflict := req
	conflict.Maximum = 7
	if _, err := ExtendCycleBudget(ctx, root, "idea", Fixup, conflict); err == nil {
		t.Fatal("conflicting decision replay accepted")
	}
	conflict.DecisionID = "stale-preview"
	if _, err := ExtendCycleBudget(ctx, root, "idea", Fixup, conflict); err == nil {
		t.Fatal("stale preview became a new grant")
	}
	// Even a binding loaded before the grant observes the actual saved maximum.
	if n, err := b.Reserve(ctx, "sixth"); err != nil || n != 6 {
		t.Fatalf("cached binding ignored operator grant: %d %v", n, err)
	}
	if _, err := b.Reserve(ctx, "seventh"); !errors.Is(err, ErrLimit) {
		t.Fatalf("finite grant became unlimited: %v", err)
	}
	current, err := LoadCycleBinding(ctx, root, "idea", Fixup)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ExtendCycleBudget(ctx, root, "idea", Fixup, cycleGrant(t, current, "grant-seven", 7)); err != nil {
		t.Fatal(err)
	}
	status, err = ExtendCycleBudget(ctx, root, "idea", Fixup, req)
	if err != nil || status.Policy.Maximum != 7 || len(status.Policy.Extensions) != 2 || status.Spent != 6 {
		t.Fatalf("old exact replay changed later grant: %+v %v", status, err)
	}
	if _, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 5, 0, "", ""); err != nil {
		t.Fatalf("original runtime cannot resume grant: %v", err)
	}
	if _, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 7, 0, "", ""); err == nil {
		t.Fatal("changed runtime flags replaced original authority")
	}
}

func TestCycleExtensionConcurrentDecisionsRequireFreshPreview(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := EnsureCycleBinding(ctx, root, "idea", CrossReview, 3, 3, "", "")
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 8)
	req := cycleGrant(t, b, "base", 4)
	for i := 0; i < 8; i++ {
		go func(i int) {
			<-start
			r := req
			r.DecisionID = fmt.Sprint(i)
			r.Maximum += i
			_, err := ExtendCycleBudget(ctx, root, "idea", CrossReview, r)
			results <- err
		}(i)
	}
	close(start)
	passed := 0
	for i := 0; i < 8; i++ {
		if err := <-results; err == nil {
			passed++
		} else if !strings.Contains(err.Error(), "changed since inspection") && !errors.Is(err, errHistoryChanged) {
			t.Errorf("unexpected conflict: %v", err)
		}
	}
	// The initial bounded read can refuse a concurrently replaced policy.
	// Both that refusal and a stale hash must still leave exactly one grant.
	status, err := InspectCycleBudget(ctx, root, "idea", CrossReview)
	if err != nil || passed != 1 || len(status.Policy.Extensions) != 1 || status.Spent != 3 {
		t.Fatalf("concurrent decision lost history: passed=%d status=%+v err=%v", passed, status, err)
	}
}

func TestCycleExtensionPublicationFailureDoesNotRefundOrDuplicate(t *testing.T) {
	for _, publish := range []bool{false, true} {
		t.Run(fmt.Sprint(publish), func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			b, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 1, 1, "", "")
			if err != nil {
				t.Fatal(err)
			}
			req := cycleGrant(t, b, "operator-one", 2)
			_, err = b.extend(ctx, req, func(path string, data []byte) error {
				if publish {
					if err := writeSynced(path, data); err != nil {
						return err
					}
				}
				return errors.New("injected policy publication failure")
			})
			if err == nil {
				t.Fatal("injected persistence failure passed")
			}
			status, err := InspectCycleBudget(ctx, root, "idea", Fixup)
			if err != nil {
				t.Fatal(err)
			}
			want := 1
			if publish {
				want = 2
			}
			if status.Policy.Maximum != want || status.Spent != 1 {
				t.Fatalf("failure reset state: %+v", status)
			}
			status, err = ExtendCycleBudget(ctx, root, "idea", Fixup, req)
			if err != nil || len(status.Policy.Extensions) != 1 || status.Policy.Maximum != 2 || status.Spent != 1 {
				t.Fatalf("recovery duplicated grant: %+v %v", status, err)
			}
		})
	}
}

func TestCycleExtensionRejectsCorruptChainAndLostCharges(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 3, 2, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Reserve(ctx, "third"); err != nil {
		t.Fatal(err)
	}
	if _, err := ExtendCycleBudget(ctx, root, "idea", Fixup, cycleGrant(t, b, "grant", 4)); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{
		strings.Replace(string(raw), `"maximum": 4,`, `"maximum": 99,`, 1),
		strings.Replace(string(raw), `"original_maximum": 3,`, `"original_maximum": 1,`, 1),
		strings.Replace(string(raw), `"spent": 3,`, `"spent": null,`, 1),
		strings.Replace(string(raw), `"spent": 3,`, "", 1),
		strings.Replace(string(raw), `"maximum": 4,`, `"Maximum": 4,`, 1),
		strings.Replace(string(raw), `"maximum": 4,`, `"maximum": 4, "maximum": 9,`, 1),
		strings.Replace(string(raw), `"previous_sha256": "`, `"previous_sha256": "f`, 1),
	} {
		if err := os.WriteFile(path, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := b.Reserve(ctx, "bad"); err == nil {
			t.Fatal("corrupt grant authorized work")
		}
		if _, err := InspectCycleBudget(ctx, root, "idea", Fixup); err == nil {
			t.Fatal("corrupt grant inspected as valid")
		}
	}
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	// The ledger's own test seam simulates a retained but rolled-back snapshot.
	if _, err := b.Store.update(ctx, func(s *Snapshot, _ time.Time) error { delete(s.Entries, key("third")); return nil }); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Reserve(ctx, "after-lost-charge"); err == nil {
		t.Fatal("grant ignored missing spent history")
	}
}

func TestCycleExtensionCannotEnableSkippedPhaseOrUnboundedGrant(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := EnsureCycleBinding(ctx, root, "idea", CrossReview, 0, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ExtendCycleBudget(ctx, root, "idea", CrossReview, cycleGrant(t, b, "enable-fast", 1)); err == nil {
		t.Fatal("finite extension enabled a skipped phase")
	}
	for _, n := range []int{-1, 0} {
		if _, err := ExtendCycleBudget(ctx, root, "idea", CrossReview, cycleGrant(t, b, "invalid", n)); err == nil {
			t.Fatalf("invalid maximum %d accepted", n)
		}
	}
}

func TestCycleExtensionRejectsInvalidInputAndClockRegression(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 1, 1, "", "")
	if err != nil {
		t.Fatal(err)
	}
	valid := cycleGrant(t, b, "grant", 2)
	for _, modify := range []func(*CycleExtensionRequest){
		func(r *CycleExtensionRequest) { r.Maximum = 0 },
		func(r *CycleExtensionRequest) { r.Maximum = int(^uint(0) >> 1) },
		func(r *CycleExtensionRequest) { r.Maximum = 1 },
		func(r *CycleExtensionRequest) { r.DecisionID = "" },
		func(r *CycleExtensionRequest) { r.Reason = "" },
		func(r *CycleExtensionRequest) { r.Reason = strings.Repeat("x", 1025) },
		func(r *CycleExtensionRequest) { r.ExpectedPolicySHA256 = "not-a-digest" },
	} {
		r := valid
		modify(&r)
		if _, err := ExtendCycleBudget(ctx, root, "idea", Fixup, r); err == nil {
			t.Fatalf("invalid grant accepted: %+v", r)
		}
	}
	status, err := b.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	b.Store.now = func() time.Time { return status.StartedAt.Add(-time.Second) }
	if _, err := b.extend(ctx, valid, writeSynced); !errors.Is(err, ErrClockSkew) {
		t.Fatalf("backward clock: %v", err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := ExtendCycleBudget(cancelled, root, "idea", Fixup, valid); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled grant: %v", err)
	}
	current, err := InspectCycleBudget(ctx, root, "idea", Fixup)
	if err != nil || current.PolicySHA256 != status.PolicySHA256 || current.Spent != 1 {
		t.Fatalf("refusal changed state: %+v %v", current, err)
	}
}

func TestCycleExtensionNestedSessionKeepsSingleChargeAcrossGrant(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := EnsureCycleBinding(ctx, root, "idea", CrossReview, 1, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	scoped, finish, err := OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	if _, err := ChargeCycle(scoped, CrossReview); err != nil {
		t.Fatal(err)
	}
	if _, err := ExtendCycleBudget(ctx, root, "idea", CrossReview, cycleGrant(t, b, "another", 2)); err != nil {
		t.Fatal(err)
	}
	current, err := LoadCycleBinding(ctx, root, "idea", CrossReview)
	if err != nil {
		t.Fatal(err)
	}
	nested, done, err := OpenCycleSession(scoped, current)
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	if n, err := ChargeCycle(nested, CrossReview); err != nil || n != 1 {
		t.Fatalf("grant split the current operation: %d %v", n, err)
	}
	status, err := InspectCycleBudget(ctx, root, "idea", CrossReview)
	if err != nil || status.Spent != 1 {
		t.Fatalf("nested grant charged another operation: %+v %v", status, err)
	}
}
