package fsutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// R3 regression pin for the row-30 AFTER-operand rule (claude-1 consult):
// a PinLstat-derived FileInfo compares via os.SameFile WITHOUT re-opening
// the path — the comparison still works after the path is gone — while a
// lazy os.Lstat operand must re-open by path and therefore cannot. On
// Windows this also pins the side effect that matters: the lazy operand's
// implicit CreateFile uses dwShareMode=0, refusing concurrent atomic
// replacements of the file; the Root-derived pin performs no such open.
// Runs on every platform (hosted Windows execution is the replace-side
// evidence); no skips.
func TestPinLstatOperandComparesWithoutPathReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "material")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	eager, err := PinLstat(path)
	if err != nil {
		t.Fatal(err)
	}
	h, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	opened, err := h.Stat()
	if err != nil {
		t.Fatal(err)
	}
	h.Close()
	lazy, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	// Remove the name: an identity comparison must not need the path back.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(opened, eager) {
		t.Fatal("Root-derived (PinLstat) operand required a path re-open — the row-30 AFTER rule is violated")
	}
	// The lazy-operand contrast is a Windows property (POSIX os.Lstat is
	// eager and compares fine after removal); on Windows the lazy operand
	// must need the path back — and, off-test, its implicit dwShareMode=0
	// open is what refused concurrent replacements.
	if runtime.GOOS == "windows" && os.SameFile(opened, lazy) {
		t.Fatal("lazy os.Lstat operand compared without the path on Windows — the mechanism model is wrong and this pin needs revisiting")
	}
}
