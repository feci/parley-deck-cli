package budget

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestCycleBindingSharesCapAcrossGitWorktrees(t *testing.T) {
	root := t.TempDir()
	other := filepath.Join(t.TempDir(), "linked")
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root, "-c", "core.hooksPath=/dev/null"}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	git("init", "-q")
	if err := os.WriteFile(filepath.Join(root, "fixture"), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	git("add", "fixture")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "-c", "commit.gpgSign=false", "commit", "-qm", "Fixture")
	git("worktree", "add", "-qb", "cycle-other", other)
	ctx := context.Background()
	first, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 3, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	second, err := LoadCycleBinding(ctx, other, "idea", Fixup)
	if err != nil || second == nil {
		t.Fatalf("linked binding: %+v %v", second, err)
	}
	if first.Store.Dir != second.Store.Dir || first.Store.Scope != second.Store.Scope {
		t.Fatal("linked worktree acquired an independent scope")
	}
	for i, b := range []*CycleBinding{first, second, first, second} {
		_, err := b.Reserve(ctx, fmt.Sprint(i))
		if i < 3 && err != nil {
			t.Fatal(err)
		}
		if i == 3 && !errors.Is(err, ErrLimit) {
			t.Fatalf("fourth linked attempt: %v", err)
		}
	}
	status, err := InspectCycleBudget(ctx, other, "idea", Fixup)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ExtendCycleBudget(ctx, other, "idea", Fixup, CycleExtensionRequest{DecisionID: "linked-grant", ExpectedPolicySHA256: status.PolicySHA256, Maximum: 4, Reason: "One finite linked-worktree fixture cycle"}); err != nil {
		t.Fatal(err)
	}
	if n, err := first.Reserve(ctx, "linked-fourth"); err != nil || n != 4 {
		t.Fatalf("linked extension not visible to cached binding: %d %v", n, err)
	}
	if _, err := second.Reserve(ctx, "linked-fifth"); !errors.Is(err, ErrLimit) {
		t.Fatalf("linked extension reset count: %v", err)
	}
}

func TestCycleBindingPreservesCarriedCountInclusiveCapAndReplay(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 5, 4, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if n, err := b.Reserve(ctx, "fifth"); err != nil || n != 5 {
		t.Fatalf("last allowed cycle: %d %v", n, err)
	}
	resumed, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 5, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := resumed.Reserve(ctx, "sixth"); !errors.Is(err, ErrLimit) {
		t.Fatalf("sixth attempt: %v", err)
	}
	if _, err := resumed.Reserve(ctx, "fifth"); !errors.Is(err, ErrReserved) {
		t.Fatalf("same identity replay: %v", err)
	}
	if _, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 99, 0, "", ""); err == nil {
		t.Fatal("runtime changed the frozen ceiling")
	}
	state, err := resumed.Store.Inspect(ctx)
	if err != nil || resumed.Count(state) != 5 {
		t.Fatalf("carried charges reset: %+v %v", state, err)
	}
	if err := os.Remove(filepath.Join(resumed.Store.Dir, "ledger.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := resumed.Reserve(ctx, "lost-ledger"); err == nil {
		t.Fatal("missing ledger authorized another cycle")
	}
}

func TestCyclePolicyRejectsRemovedNullAndAliasedFields(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	b, err := EnsureCycleBinding(ctx, root, "idea", CrossReview, 3, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{
		strings.Replace(string(raw), `"maximum": 3,`, "", 1),
		strings.Replace(string(raw), `"maximum": 3`, `"maximum": null`, 1),
		strings.Replace(string(raw), `"carried": 0`, `"carried": null`, 1),
		strings.Replace(string(raw), `"maximum": 3`, `"Maximum": 99`, 1),
		strings.Replace(string(raw), `"maximum": 3`, `"maximum": 3, "maximum": 99`, 1),
		strings.Replace(string(raw), `"idea_path": "parley-deck/ideas/idea"`, `"idea_path": ".."`, 1),
	} {
		if err := os.WriteFile(path, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadCycleBinding(ctx, root, "idea", CrossReview); err == nil {
			t.Fatalf("malformed frozen policy accepted: %s", bad)
		}
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureCycleBinding(ctx, root, "idea", CrossReview, 3, 0, "", ""); err == nil {
		t.Fatal("missing policy became a new grant")
	}
}

func TestCycleSessionNestedCallsChargeOnceAndCloseRevokes(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	b, err := EnsureCycleBinding(ctx, root, "idea", CrossReview, 1, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	scoped, finish, err := OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			child, close, err := OpenCycleSession(scoped, b)
			if err != nil {
				t.Error(err)
				return
			}
			defer close()
			if n, err := ChargeCycle(child, CrossReview); err != nil || n != 1 {
				t.Errorf("nested charge: %d %v", n, err)
			}
		}()
	}
	wg.Wait()
	finish()
	if _, err := ChargeCycle(scoped, CrossReview); err == nil {
		t.Fatal("ended session grants another action")
	}
	next, close, err := OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer close()
	if _, err := ChargeCycle(next, CrossReview); !errors.Is(err, ErrLimit) {
		t.Fatalf("new operation reset cap: %v", err)
	}
}

func TestCycleReservationPersistenceFailureStaysSpent(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	b, err := EnsureCycleBinding(ctx, root, "idea", Fixup, 1, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	b.Store.persist = func(path string, data []byte) error {
		if err := writeSynced(path, data); err != nil {
			return err
		}
		return errors.New("publication acknowledged with injected fault")
	}
	scoped, close, err := OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer close()
	if _, err := ChargeCycle(scoped, Fixup); err == nil {
		t.Fatal("failed persistence granted work")
	}
	if _, err := ChargeCycle(scoped, Fixup); err == nil {
		t.Fatal("cached failure disappeared")
	}
	state, err := b.Store.Inspect(ctx)
	if err != nil || b.Count(state) != 1 {
		t.Fatalf("failed charged attempt disappeared: %+v %v", state, err)
	}
}

func TestCycleCancelledSessionDoesNotChargeOrRefund(t *testing.T) {
	ctx := context.Background()
	b, err := EnsureCycleBinding(ctx, t.TempDir(), "idea", Fixup, 1, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	scoped, finish, err := OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	cancelled, cancel := context.WithCancel(scoped)
	cancel()
	if _, err := ChargeCycle(cancelled, Fixup); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
	state, err := b.Store.Inspect(ctx)
	if err != nil || b.Count(state) != 0 {
		t.Fatalf("cancelled attempt was charged: %+v %v", state, err)
	}
	if n, err := ChargeCycle(scoped, Fixup); err != nil || n != 1 {
		t.Fatalf("live parent: %d %v", n, err)
	}
	if _, err := ChargeCycle(cancelled, Fixup); !errors.Is(err, context.Canceled) {
		t.Fatalf("reserved cancellation: %v", err)
	}
	state, err = b.Store.Inspect(ctx)
	if err != nil || b.Count(state) != 1 {
		t.Fatalf("cancellation refunded a charge: %+v %v", state, err)
	}
}

func TestCycleIndependentConcurrentSessionsShareCap(t *testing.T) {
	ctx := context.Background()
	b, err := EnsureCycleBinding(ctx, t.TempDir(), "idea", CrossReview, 3, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 8)
	for i := 0; i < 8; i++ {
		go func() {
			<-start
			scoped, finish, err := OpenCycleSession(ctx, b)
			defer finish()
			if err == nil {
				_, err = ChargeCycle(scoped, CrossReview)
			}
			results <- err
		}()
	}
	close(start)
	granted, refused := 0, 0
	for i := 0; i < 8; i++ {
		err := <-results
		if err == nil {
			granted++
		} else if errors.Is(err, ErrLimit) {
			refused++
		} else {
			t.Errorf("unexpected concurrent refusal: %v", err)
		}
	}
	state, err := b.Store.Inspect(ctx)
	if err != nil || granted != 3 || refused != 5 || b.Count(state) != 3 {
		t.Fatalf("concurrent cap: granted=%d refused=%d state=%+v err=%v", granted, refused, state, err)
	}
	active, finish, err := OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	changed := *b
	changed.Policy.Maximum++
	if CycleSessionMatches(active, &changed) {
		t.Fatal("altered nested policy matched")
	}
	if _, _, err := OpenCycleSession(active, &changed); err == nil {
		t.Fatal("altered nested policy joined")
	}
}
