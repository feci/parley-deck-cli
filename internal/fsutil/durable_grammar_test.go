package fsutil

import (
	"regexp"
	"strings"
	"testing"
)

// AC-DUR-3 requirement (a): the deterministic stage name is a fixed
// derivation of the final name and must be provably non-colliding with every
// valid final name of the §B plain rows. Valid finals are exactly three
// grammars: A1's fixed artifact names, B1's fixed runtime directory names,
// and B2's telemetry IDs (UUID v4: lowercase hex groups joined by dashes —
// telemetry.NewID's format, which contains no dots). A dotted suffix can
// never be part of the UUID grammar, so final+StageSuffix is always outside
// every grammar — no stage name can equal a valid final, and no two distinct
// finals can share one stage name.
func TestStageSuffixOutsideEveryFinalGrammar(t *testing.T) {
	fixedFinals := []string{
		"request.json", "parent-result.json", // A1
		".parley-runtime", "trajectory-verification", // B1
	}
	uuidFinal := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`) // B2

	if !strings.Contains(StageSuffix, ".") {
		t.Fatalf("StageSuffix %q must contain a dot to sit outside the dot-free UUID grammar", StageSuffix)
	}
	isValidFinal := func(name string) bool {
		if uuidFinal.MatchString(name) {
			return true
		}
		for _, f := range fixedFinals {
			if name == f {
				return true
			}
		}
		return false
	}
	for _, f := range fixedFinals {
		if strings.HasSuffix(f, StageSuffix) {
			t.Fatalf("fixed final %q ends with the stage suffix", f)
		}
		staged := f + StageSuffix
		if isValidFinal(staged) {
			t.Fatalf("staged name %q collides with a valid final", staged)
		}
		// B2 finals cannot carry the suffix at all: no dots in the grammar.
		if uuidFinal.MatchString(staged) {
			t.Fatalf("staged name %q collides with the UUID grammar", staged)
		}
	}
	// The UUID grammar itself rejects any dotted suffix.
	for _, id := range []string{
		"0f4d2b3a-9c1e-4f6a-8b2d-7e5c1a0b3d9f",
		"ffffffff-ffff-4fff-8fff-ffffffffffff",
	} {
		if !uuidFinal.MatchString(id) {
			t.Fatalf("grammar sample %q is not a valid UUID final", id)
		}
		if staged := id + StageSuffix; isValidFinal(staged) {
			t.Fatalf("staged UUID %q collides with a valid final", staged)
		}
		if strings.HasSuffix(id, StageSuffix) {
			t.Fatalf("UUID final %q ends with the stage suffix", id)
		}
	}
}
