package pipeline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReviewAgreedFixesParsesContract(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "consensus.md")
	os.WriteFile(p, []byte("---\nidea: x\noutstanding_agreed_fixes: 2\nblocked: false\n---\n\n## Agreed fixes\n"), 0o644)
	count, blocked, found, err := ReviewAgreedFixes(p)
	if err != nil || !found || count != 2 || blocked {
		t.Fatalf("got count=%d blocked=%v found=%v err=%v", count, blocked, found, err)
	}

	os.WriteFile(p, []byte("---\nidea: x\noutstanding_agreed_fixes: 0\nblocked: true\n---\n"), 0o644)
	count, blocked, found, _ = ReviewAgreedFixes(p)
	if count != 0 || !blocked || !found {
		t.Fatalf("blocked case: count=%d blocked=%v found=%v", count, blocked, found)
	}

	// Missing field -> found=false (auto fails closed).
	os.WriteFile(p, []byte("---\nidea: x\n---\n\n## Agreed fixes\n"), 0o644)
	_, _, found, _ = ReviewAgreedFixes(p)
	if found {
		t.Fatal("missing outstanding_agreed_fixes must report found=false")
	}
}

func TestPhase8Decision(t *testing.T) {
	cases := []struct {
		name                       string
		outstanding                int
		blocked, found             bool
		cycle, max                 int
		want                       Phase8Action
	}{
		{"missing contract fails closed", 0, false, false, 1, 3, Phase8NoData},
		{"blocked stops", 2, true, true, 1, 3, Phase8Blocked},
		{"zero completes", 0, false, true, 1, 3, Phase8Complete},
		{"fixes under cap -> fixup", 3, false, true, 1, 3, Phase8Fixup},
		{"fixes at cap -> maxed", 3, false, true, 3, 3, Phase8Maxed},
		{"blocked beats zero", 0, true, true, 1, 3, Phase8Blocked},
	}
	for _, c := range cases {
		if got := Phase8Decision(c.outstanding, c.blocked, c.found, c.cycle, c.max); got != c.want {
			t.Errorf("%s: got %q want %q", c.name, got, c.want)
		}
	}
}
