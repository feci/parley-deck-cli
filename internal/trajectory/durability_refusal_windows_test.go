//go:build windows

package trajectory

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/fsutil"
)

// AC-DUR-1/§B: the fsutil.SyncDir named-type contract refuses fail-closed on
// Windows, and no directory handle can pass through SyncFile silently.
func TestSyncDirContractRefusesFailClosed(t *testing.T) {
	dir := t.TempDir()
	f, err := os.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	handle, err := fsutil.DirHandleOf(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := fsutil.SyncDir(handle); !errors.Is(err, fsutil.ErrDirEntryDurabilityUnsupported) {
		t.Fatalf("SyncDir must return the named refusal, got: %v", err)
	}
	f.Close()
	// The mechanical backstop: a directory handle reaching SyncFile gets the
	// same named refusal, never a raw FlushFileBuffers failure.
	d, err := os.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	if err := fsutil.SyncFile(d); !errors.Is(err, fsutil.ErrDirEntryDurabilityUnsupported) {
		t.Fatalf("SyncFile on a directory must refuse with the named error, got: %v", err)
	}
	// A regular FILE still syncs normally through SyncFile.
	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	rw, err := os.OpenFile(file, os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := fsutil.SyncFile(rw); err != nil {
		t.Fatalf("SyncFile on a writable regular file must work: %v", err)
	}
	rw.Close()
	// DirHandleOf fails closed on a non-directory handle.
	nf, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fsutil.DirHandleOf(nf); err == nil {
		t.Fatal("DirHandleOf accepted a regular file")
	}
	nf.Close()
}

// AC-DUR-2 hosted pins: each Windows-reachable rooted row refuses
// pre-mutation with the named, blocking refusal and publishes nothing.
func TestRootedRowsRefusePreMutationOnWindows(t *testing.T) {
	base := t.TempDir()
	root, err := os.OpenRoot(base)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()

	// A3 — publishReservationIntent refuses before the O_EXCL create.
	intent := reservationIntent{Version: 1}
	if _, err := publishReservationIntent(root, intent); err == nil ||
		!strings.Contains(err.Error(), "precharge reservation-intent publication") ||
		!strings.Contains(err.Error(), "nothing was published") {
		t.Fatalf("A3 refusal wrong: %v", err)
	}
	if entries, err := os.ReadDir(filepath.Join(base, "reservation-intents")); err == nil && len(entries) > 0 {
		t.Fatalf("A3 refusal published something: %v", entries)
	}
	if _, err := os.Lstat(filepath.Join(base, "reservation-intents")); !os.IsNotExist(err) {
		t.Fatalf("A3 refusal must fire before the intent root exists: %v", err)
	}

	// B5 — openIntentRoot(create=true) refuses before the mkdir.
	if _, err := openIntentRoot(budget.CycleBinding{}, true); err == nil ||
		!strings.Contains(err.Error(), "precharge reservation-intent directory publication") {
		t.Fatalf("B5 refusal wrong: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(base, "reservation-intents")); !os.IsNotExist(err) {
		t.Fatalf("B5 refusal must fire before the mkdir: %v", err)
	}

	// C1 — publishParentRecoveryWithSync refuses before the stage or rename.
	if err := publishParentRecoveryWithSync(root, ".parley-runtime/x", ParentRecovery{}, func(*os.Root, string) error {
		t.Fatal("sync barrier must never run after a refusal")
		return nil
	}); err == nil || !strings.Contains(err.Error(), "parent-recovery publication") {
		t.Fatalf("C1 refusal wrong: %v", err)
	}

	// C2 — publishRecoveredParent refuses before the stage or rename.
	if err := publishRecoveredParent(root, ".parley-runtime/x", RecoveredParentRecord{}); err == nil ||
		!strings.Contains(err.Error(), "recovered-parent publication") {
		t.Fatalf("C2 refusal wrong: %v", err)
	}
}
