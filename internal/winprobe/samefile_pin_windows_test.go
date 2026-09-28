//go:build windows

package winprobe

import (
	"os"
	"path/filepath"
	"testing"
)

// Discriminating probe for the row-30 mechanism question (claude-1 advisory
// consult, §6, 2026-09-28): on Windows, an os.Stat/os.Lstat FileInfo defers
// its identity lookup to the first os.SameFile call, re-opening BY PATH at
// comparison time — so a lazy "before" operand re-resolves to whatever now
// lives at the path and is structurally incapable of detecting a
// replace-at-same-path. A handle-derived FileInfo (File.Stat / Root.Stat /
// Root.Lstat) captures the NTFS file identity eagerly.
//
// Predictions (both must hold for the lazy-mechanism attribution):
//   SameFile(lazyPin,  after) == true   — the lazy idiom is a no-op here (the bug)
//   SameFile(eagerPin, after) == false  — eager identity detects the replacement
//
// If the second prediction FAILS (eager compares equal), NTFS identity reuse
// is in play after all and the consult's mechanism is wrong — the probe fails
// loudly either way, converting the source-derived claim into a hosted fact.
// Mechanics/diagnostic only: no product change rides on this probe without
// its hosted outcome.
func TestSameFileLazyVsEagerPin(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "material")
	if err := os.WriteFile(path, []byte("identical bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	lazyPin, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	h, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	eagerPin, err := h.Stat()
	if err != nil {
		t.Fatal(err)
	}
	h.Close()

	// The adversary: replace the file at the same path with identical bytes.
	if err := os.Rename(path, filepath.Join(dir, "moved-aside")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("identical bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	after, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}

	if !os.SameFile(lazyPin, after) {
		t.Fatal("PREDICTED-BUG ABSENT: SameFile(lazyPin, after) == false — the lazy Lstat pin captured real identity; the consult's mechanism does not hold on this platform/toolchain")
	}
	if os.SameFile(eagerPin, after) {
		t.Fatal("NTFS IDENTITY REUSE: SameFile(eagerPin, after) == true — the eager handle-derived pin ALSO compares equal; the problem is filesystem identity reuse, not the lazy FileInfo, and the consult's §1 is wrong")
	}
}
