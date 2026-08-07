package protocolcore

import (
	"strings"
	"testing"
)

// The diff must not be the thing that crashes on a large input. The earlier full DP table was
// O(n*m) MEMORY — a 100k-line document would have allocated tens of gigabytes. Hirschberg is
// linear-space; this exercises a size where the old approach would already be allocating ~1.6GB.
func TestLargeDocumentDoesNotExhaustMemory(t *testing.T) {
	var a []string
	for i := 0; i < 20000; i++ {
		a = append(a, "line "+strings.Repeat("x", i%7))
	}
	b := append([]string(nil), a...)
	b[10000] = "changed"

	mask := lcsMask(a, b)
	if len(mask) != len(a) {
		t.Fatalf("mask length %d, want %d", len(mask), len(a))
	}
	if !mask[0] || !mask[len(mask)-1] {
		t.Error("unchanged boundary lines were not matched")
	}
	if mask[10000] {
		t.Error("the changed line was reported as carried forward")
	}
}

// A moved line must show up as not-carried at its old position — this is what a multiset cannot do.
func TestLCSSeesOrder(t *testing.T) {
	a := []string{"alpha", "beta", "gamma"}
	b := []string{"gamma", "beta", "alpha"}
	mask := lcsMask(a, b)
	matched := 0
	for _, m := range mask {
		if m {
			matched++
		}
	}
	if matched >= len(a) {
		t.Errorf("reordering was treated as no change (matched %d of %d)", matched, len(a))
	}
}
