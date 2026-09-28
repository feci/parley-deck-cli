//go:build windows

package budget

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

// foreignExclusiveHandle opens path with NO sharing, the way a third party
// (antivirus, indexer) transiently does — the §D.6 self-held/foreign-held
// classification: this process's own handles share read|write, so only a
// foreign holder can produce a sharing violation.
func foreignExclusiveHandle(t *testing.T, path string) windows.Handle {
	t.Helper()
	p16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatal(err)
	}
	h, err := windows.CreateFile(p16, windows.GENERIC_READ, 0, nil,
		windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatalf("foreign open: %v", err)
	}
	return h
}

// The §D.6 open-site diagnostic cycle: reproduce the foreign sharing holder,
// log the ACTUAL error class the product open sees, and record the bounded
// retry outcome — recovery after the holder leaves, and loud exhaustion when
// it stays. The product's own opens cannot self-conflict (share read|write).
func TestOpenLockFileSharingDiagnostic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "899a.lock")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Class 1 — transient foreign holder: the bounded retry recovers.
	h := foreignExclusiveHandle(t, path)
	recovered := make(chan error, 1)
	go func() {
		start := time.Now()
		f, err := openLockFile(path)
		if err == nil {
			f.Close()
		}
		t.Logf("diagnostic: transient foreign holder — openLockFile err=%v elapsed=%s (retry budget 250ms)", err, time.Since(start).Round(10*time.Millisecond))
		recovered <- err
	}()
	time.Sleep(80 * time.Millisecond) // holder present for one+ retry slot
	windows.CloseHandle(h)
	if err := <-recovered; err != nil {
		t.Fatalf("bounded retry did not recover after the foreign holder left: %v", err)
	}

	// Class 2 — persistent foreign holder: loud exhaustion within the budget,
	// never an unbounded wait, naming the retry budget.
	h = foreignExclusiveHandle(t, path)
	defer windows.CloseHandle(h)
	start := time.Now()
	_, err := openLockFile(path)
	elapsed := time.Since(start)
	t.Logf("diagnostic: persistent foreign holder — openLockFile err=%v elapsed=%s (must be the named budget error)", err, elapsed.Round(10*time.Millisecond))
	if err == nil {
		t.Fatal("persistent foreign holder must not be waited out silently")
	}
	if !strings.Contains(err.Error(), "retry budget") {
		t.Fatalf("exhaustion must be loud and name the budget: %v", err)
	}
	if !errors.Is(err, windows.ERROR_SHARING_VIOLATION) {
		t.Fatalf("the ACTUAL error class must be ERROR_SHARING_VIOLATION: %v", err)
	}
	if elapsed < 200*time.Millisecond || elapsed > 2*time.Second {
		t.Fatalf("retry window out of bounds: %s", elapsed)
	}
}
