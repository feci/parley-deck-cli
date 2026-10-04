package membership

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
)

func TestQuotaFixupProjectionWaitAndProcessExclusivity(t *testing.T) {
	if dir := os.Getenv("PARLEY_PROJECTION_CHILD"); dir != "" {
		release, e := ProjectionLock(dir)
		if e == nil {
			release()
			os.Exit(21)
		}
		if !strings.Contains(e.Error(), "projection lock held") {
			os.Exit(22)
		}
		return
	}
	_, dir, _ := fixture(t, quota.NewPolicy(nil, nil))
	release, e := ProjectionLock(dir)
	if e != nil {
		t.Fatal(e)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestQuotaFixupProjectionWaitAndProcessExclusivity$")
	cmd.Env = append(os.Environ(), "PARLEY_PROJECTION_CHILD="+dir)
	start := time.Now()
	raw, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatalf("projection child: %v %s", e, raw)
	}
	if time.Since(start) < projectionWait {
		t.Fatal("projection did not wait its bound")
	}
	release()
	held, e := ProjectionLock(dir)
	if e != nil {
		t.Fatal(e)
	}
	result := make(chan error, 1)
	go func() {
		r, e := ProjectionLock(dir)
		if e == nil {
			r()
		}
		result <- e
	}()
	time.Sleep(50 * time.Millisecond)
	held()
	if e := <-result; e != nil {
		t.Fatal("bounded waiter did not acquire after release", e)
	}
	path, _ := lockPath(dir, "projection")
	t.Logf("%s: different-process refusal after bounded wait; waiter acquired after release", path)
}
func TestQuotaFixupCorruptReceiptRecoveryChecksAllProjections(t *testing.T) {
	for _, bad := range []string{"", "wrong-transition\n"} {
		t.Run(bad, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			b, e := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
			if e != nil {
				t.Fatal(e)
			}
			receipt := filepath.Join(dir, "quota-applied", b.ID)
			os.WriteFile(receipt, []byte(bad), 0600)
			if _, _, e = protocol.QuotaMembers(dir, nil); e == nil {
				t.Fatal("broken receipt admitted signoff")
			}
			m, e := runmanifest.Load(root, run)
			if e != nil {
				t.Fatal(e)
			}
			good := m
			m.Participants = []string{"intruder"}
			runmanifest.Write(root, run, m)
			if _, e = Before(ctx, root, dir, run); e == nil {
				t.Fatal("contradictory manifest repaired blindly")
			}
			raw, _ := os.ReadFile(receipt)
			if string(raw) != bad {
				t.Fatal("bad receipt marked applied before manifest validation")
			}
			runmanifest.Write(root, run, good)
			for i := 0; i < 2; i++ {
				if _, e = Before(ctx, root, dir, run); e != nil {
					t.Fatal(e)
				}
			}
			if _, _, e = protocol.QuotaMembers(dir, nil); e != nil {
				t.Fatal(e)
			}
			events, e := store.New(filepath.Join(root, protocol.DeckDir, "runs", run)).Load()
			if e != nil {
				t.Fatal(e)
			}
			n := 0
			for _, ev := range events {
				if ev.Type == "round.completed" && ev.Data["quota_transition"] == b.ID {
					n++
				}
			}
			if n != 1 {
				t.Fatalf("duplicate evaluation: %d", n)
			}
			notices, _ := filepath.Glob(filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_batch-*.md"))
			if len(notices) != 1 {
				t.Fatal(notices)
			}
		})
	}
}
func TestQuotaFixupReleasedContextCannotBeReused(t *testing.T) {
	_, dir, run := fixture(t, quota.NewPolicy(nil, nil))
	ctx, r, e := Acquire(context.Background(), dir, run)
	if e != nil {
		t.Fatal(e)
	}
	r()
	if CheckLease(ctx, dir, run) == nil {
		t.Fatal("released context accepted")
	}
}

func TestQuotaFixupLegacyHasNoLifetimeOrProjectionSingleton(t *testing.T) {
	root := t.TempDir()
	if e := protocol.InitWorkspace(root); e != nil {
		t.Fatal(e)
	}
	idea, e := protocol.CreateIdea(root, "legacy fixture", []string{"a", "b"})
	if e != nil {
		t.Fatal(e)
	}
	for _, run := range []string{"legacy-a", "legacy-b"} {
		_, r, e := Acquire(context.Background(), idea.Path, run)
		if e != nil {
			t.Fatal(e)
		}
		defer r()
		p, e := ProjectionLock(idea.Path)
		if e != nil {
			t.Fatal(e)
		}
		defer p()
	}
	if _, e = os.Stat(filepath.Join(root, ".parley-runtime/membership")); !os.IsNotExist(e) {
		t.Fatal("legacy idea created singleton state", e)
	}
}
