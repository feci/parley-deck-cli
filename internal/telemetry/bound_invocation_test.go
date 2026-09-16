package telemetry

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// BeginBound exists for one trusted bridge: an immutable launch authority pins
// the invocation ID before the launch exists, and the real process runner must
// allocate exactly that invocation, exactly once, or fail closed.

func TestBeginBoundAllocatesPinnedInvocationWithOrdinaryLifecycle(t *testing.T) {
	store := t.TempDir()
	inv, err := BeginBound(store, "recovered-verifier-invocation", Metadata{Agent: "reviewer", Adapter: "claude", AttemptOrdinal: 1})
	if err != nil {
		t.Fatal(err)
	}
	if inv.ID != "recovered-verifier-invocation" || inv.Dir != filepath.Join(store, "recovered-verifier-invocation") {
		t.Fatalf("bound invocation identity drifted: %q %q", inv.ID, inv.Dir)
	}
	requested := readRecord(t, filepath.Join(inv.Dir, "requested.json"))
	if requested.InvocationID != inv.ID || requested.Type != "invocation.requested" || requested.StartedAt != nil || requested.Outcome != nil {
		t.Fatalf("bound request record is not an ordinary request: %+v", requested)
	}
	// The started/terminal records are written by the ordinary lifecycle the
	// real process runner already uses — nothing bound-specific after Begin.
	if err := inv.Started(4242); err != nil {
		t.Fatal(err)
	}
	exit := 0
	if err := inv.Finish(Outcome{Status: "process-exited", ExitCode: &exit}); err != nil {
		t.Fatal(err)
	}
	terminal := readRecord(t, filepath.Join(inv.Dir, "terminal.json"))
	if terminal.StartedAt == nil || terminal.PID == nil || *terminal.PID != 4242 ||
		terminal.Outcome == nil || terminal.Outcome.ExitCode == nil || *terminal.Outcome.ExitCode != 0 {
		t.Fatalf("bound terminal record lost its observed process: %+v", terminal)
	}
	info, err := os.Stat(filepath.Join(inv.Dir, "terminal.json"))
	if err != nil || info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("bound record readable by others: %v %o", err, info.Mode().Perm())
	}
}

func TestValidBoundInvocationID(t *testing.T) {
	valid := []string{"recovered-verifier-invocation", "a", "invocation.2026-09-16T10:00:00Z", "recovered-a-1"}
	for _, id := range valid {
		if !ValidBoundInvocationID(id) {
			t.Fatalf("valid bound ID refused: %q", id)
		}
	}
	// Separators would escape or alias inside the store; dot segments would
	// alias the store itself; over-long or secret-shaped labels are not safe
	// telemetry labels at all.
	invalid := []string{"", ".", "..", "a/b", `a\b`, "sub/dir", "/abs", "rel/",
		"spaces in id", "nl\nid", strings.Repeat("a", 200), strings.Repeat("ab", 16)}
	for _, id := range invalid {
		if ValidBoundInvocationID(id) {
			t.Fatalf("invalid bound ID accepted: %q", id)
		}
	}
}

func TestBeginBoundRefusesEveryReuseAndCollision(t *testing.T) {
	for _, tc := range []struct {
		name  string
		plant func(t *testing.T, store, id string)
	}{
		{name: "completed-invocation", plant: func(t *testing.T, store, id string) {
			inv, err := BeginBound(store, id, Metadata{Agent: "reviewer"})
			if err != nil {
				t.Fatal(err)
			}
			if err = inv.Started(9); err != nil {
				t.Fatal(err)
			}
			if err = inv.Finish(Outcome{Status: "failed", FailureClass: String("timeout")}); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "plain-directory", plant: func(t *testing.T, store, id string) {
			if err := os.Mkdir(filepath.Join(store, id), 0o700); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "regular-file", plant: func(t *testing.T, store, id string) {
			if err := os.WriteFile(filepath.Join(store, id), []byte("{}\n"), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "symlink-alias", plant: func(t *testing.T, store, id string) {
			if err := os.Symlink(t.TempDir(), filepath.Join(store, id)); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := t.TempDir()
			const id = "recovered-verifier-invocation"
			tc.plant(t, store, id)
			if _, err := BeginBound(store, id, Metadata{Agent: "reviewer"}); err == nil {
				t.Fatalf("%s: bound invocation reused a non-fresh name", tc.name)
			}
			entries, err := os.ReadDir(store)
			if err != nil || len(entries) != 1 || entries[0].Name() != id {
				t.Fatalf("%s: refusal changed the store inventory: %v %v", tc.name, entries, err)
			}
			if link, lerr := os.Readlink(filepath.Join(store, id)); lerr == nil && link != "" {
				if target, terr := os.ReadDir(link); terr != nil || len(target) != 0 {
					t.Fatalf("%s: refusal wrote through the symlink alias: %v %v", tc.name, target, terr)
				}
			}
		})
	}
}

func TestBeginBoundRefusesUnsafeIDsAndAliasedStores(t *testing.T) {
	store := t.TempDir()
	for _, id := range []string{"", ".", "..", "a/b", `a\b`, strings.Repeat("a", 200)} {
		if _, err := BeginBound(store, id, Metadata{Agent: "reviewer"}); err == nil {
			t.Fatalf("unsafe bound ID accepted: %q", id)
		}
	}
	entries, err := os.ReadDir(store)
	if err != nil || len(entries) != 0 {
		t.Fatalf("unsafe IDs created invocation directories: %v %v", entries, err)
	}
	// The shared store itself must be a real directory, not a symlink alias.
	real := t.TempDir()
	aliased := filepath.Join(t.TempDir(), "invocations")
	if err = os.Symlink(real, aliased); err != nil {
		t.Fatal(err)
	}
	if _, err = BeginBound(aliased, "recovered-verifier-invocation", Metadata{Agent: "reviewer"}); err == nil {
		t.Fatal("aliased telemetry store accepted")
	}
	if entries, err = os.ReadDir(real); err != nil || len(entries) != 0 {
		t.Fatalf("aliased store wrote through the symlink: %v %v", entries, err)
	}
}

func TestConcurrentBeginBoundAllocatesExactlyOnce(t *testing.T) {
	store := t.TempDir()
	const id = "recovered-verifier-invocation"
	var won atomic.Int32
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			inv, err := BeginBound(store, id, Metadata{Agent: "reviewer"})
			if err != nil {
				return
			}
			won.Add(1)
			if err := inv.Finish(Outcome{Status: "failed"}); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if won.Load() != 1 {
		t.Fatalf("bound invocation allocated %d times", won.Load())
	}
	if requested := readRecord(t, filepath.Join(store, id, "requested.json")); requested.InvocationID != id {
		t.Fatalf("surviving invocation lost its pinned identity: %+v", requested)
	}
}
