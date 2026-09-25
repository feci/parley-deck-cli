package protocol

import (
	"os"
	"path/filepath"
	"testing"
)

// T-1 (AC-5): the four-state parse — absent / present-empty / none / id — plus the
// malformed inline-comment value (R4: fails closed, never repaired) and the quoted
// forms (R5). Each case fails if the corresponding rule is removed.
func TestImplementerDesignationFourStates(t *testing.T) {
	for _, tc := range []struct {
		name   string
		meta   map[string]string
		want   DesignationState
		wantID string
	}{
		{"key absent", map[string]string{}, DesignationAbsent, ""},
		{"present empty", map[string]string{"implementer": ""}, DesignationEmpty, ""},
		{"present whitespace", map[string]string{"implementer": "   "}, DesignationEmpty, ""},
		{"none", map[string]string{"implementer": "none"}, DesignationNone, ""},
		{"none any casing", map[string]string{"implementer": "NoNe"}, DesignationNone, ""},
		{"none padded", map[string]string{"implementer": "  none  "}, DesignationNone, ""},
		{"id", map[string]string{"implementer": "kimi-1"}, DesignationSet, "kimi-1"},
		{"double-quoted id", map[string]string{"implementer": `"kimi-1"`}, DesignationSet, "kimi-1"},
		{"single-quoted id", map[string]string{"implementer": `'kimi-1'`}, DesignationSet, "kimi-1"},
		// R4: ReadFrontmatter does no comment stripping, so the value stays literal
		// and malformed — it must NOT be repaired into a bare id.
		{"malformed inline comment never repaired", map[string]string{"implementer": "kimi-1  # from the global default"}, DesignationSet, "kimi-1  # from the global default"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ImplementerDesignationFromMeta(tc.meta)
			if got.State != tc.want || got.ID != tc.wantID {
				t.Fatalf("got %+v, want state=%v id=%q", got, tc.want, tc.wantID)
			}
		})
	}
}

// The bare `implementer:` line must survive ReadFrontmatter as present-with-empty —
// the property the four-state parse stands on (R2).
func TestReadImplementerDesignationBareKeyIsPresentEmpty(t *testing.T) {
	dir := t.TempDir()
	prompt := "---\nidea: x\nparticipants: [a, b]\nimplementer:\nstatus: round-01\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte(prompt), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ReadImplementerDesignation(dir); got.State != DesignationEmpty {
		t.Fatalf("bare `implementer:` line: got %+v, want DesignationEmpty", got)
	}
	if got := ReadImplementerDesignation(filepath.Join(dir, "missing")); got.State != DesignationAbsent {
		t.Fatalf("missing 00-prompt.md: got %+v, want DesignationAbsent", got)
	}
}

// The shared chain: ordered sources, explicit eligibility list, (id, source, ok).
func TestResolveImplementerChain(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("FINAL.md", "---\nidea: x\ndrafted-by: drafter-1\n---\n")
	eligible := []string{"impl-1", "drafter-1", "rev-1"}

	// FINAL metadata resolves when no pin exists.
	id, src, ok := ResolveImplementerChain(dir, LegacyImplementerCandidates(), eligible)
	if !ok || id != "drafter-1" || src != SourceImplementerLegacy {
		t.Fatalf("FINAL fallback: got (%q,%q,%v)", id, src, ok)
	}

	// The pin outranks FINAL metadata.
	write("IMPLEMENTATION.md", "---\nidea: x\nimplementer: impl-1\n---\n")
	id, src, ok = ResolveImplementerChain(dir, LegacyImplementerCandidates(), eligible)
	if !ok || id != "impl-1" || src != SourceImplementerPin {
		t.Fatalf("pin precedence: got (%q,%q,%v)", id, src, ok)
	}

	// An ineligible id is skipped, not honored.
	write("IMPLEMENTATION.md", "---\nidea: x\nimplementer: ghost-9\n---\n")
	id, _, ok = ResolveImplementerChain(dir, LegacyImplementerCandidates(), eligible)
	if !ok || id != "drafter-1" {
		t.Fatalf("ineligible pin must be skipped: got (%q,%v)", id, ok)
	}
	if _, _, ok := ResolveImplementerChain(dir, PinImplementerCandidates(), eligible); ok {
		t.Fatal("pin-only chain must not read FINAL.md")
	}

	// Nothing eligible anywhere → ok=false (the caller's tail decides).
	write("FINAL.md", "---\nidea: x\ndrafted-by: ghost-9\n---\n")
	if _, _, ok := ResolveImplementerChain(dir, LegacyImplementerCandidates(), eligible); ok {
		t.Fatal("no eligible source must report ok=false")
	}

	// Explicit pre-resolved candidates (the designation tiers' shape) honor order.
	cands := []ImplementerCandidate{
		{Source: SourceImplementerDesignation, ID: "ghost-9"},
		{Source: SourceImplementerGlobalDefault, ID: "impl-1"},
	}
	id, src, ok = ResolveImplementerChain(dir, cands, eligible)
	if !ok || id != "impl-1" || src != SourceImplementerGlobalDefault {
		t.Fatalf("explicit candidates: got (%q,%q,%v)", id, src, ok)
	}
}

// The source enum keeps the R47 distinctions: pin / designation / global-default /
// none / fall-through-by-cause are pairwise distinct, so no consumer can confuse an
// opt-out with an unset deck or a fall-through.
func TestImplementerSourceDistinctions(t *testing.T) {
	seen := map[ImplementerSource]bool{}
	for _, s := range []ImplementerSource{
		SourceImplementerPin, SourceImplementerDesignation, SourceImplementerGlobalDefault,
		SourceImplementerNone, SourceImplementerFallThroughUnavailable,
		SourceImplementerFallThroughInapplicable, SourceImplementerLegacy,
	} {
		if s == "" {
			t.Fatal("source constants must be non-empty")
		}
		if seen[s] {
			t.Fatalf("duplicate source value %q", s)
		}
		seen[s] = true
	}
}

// The waiver and reassignment records (R18/R28): the §9.0-shaped confirmation is
// required; a bare line is not a record.
func TestDesignationRecordsRequireConfirmation(t *testing.T) {
	if ImplementerWaived(map[string]string{"implementer_waived": "kimi-1 — offline"}, "kimi-1") {
		t.Fatal("waiver without confirmation must not clear the gate")
	}
	if !ImplementerWaived(map[string]string{"implementer_waived": "kimi-1 — offline — confirmed 2026-09-25"}, "kimi-1") {
		t.Fatal("confirmed waiver must clear the gate")
	}
	if ImplementerWaived(map[string]string{"implementer_waived": "zcode-1 — offline — confirmed 2026-09-25"}, "kimi-1") {
		t.Fatal("waiver naming a different agent must not clear this designee's gate")
	}
	meta := map[string]string{"implementer_reassigned": "claude-1 to kimi-1 — drafter rotation — confirmed 2026-09-25"}
	if !ParseImplementerReassignment(meta, "claude-1", "kimi-1") {
		t.Fatal("confirmed reassignment of the disputed pair must parse")
	}
	if ParseImplementerReassignment(map[string]string{"implementer_reassigned": "claude-1 to kimi-1 — rotation"}, "claude-1", "kimi-1") {
		t.Fatal("reassignment without confirmation must not re-pin")
	}
	if ParseImplementerReassignment(meta, "kimi-1", "claude-1") {
		t.Fatal("reassignment of a different pair must not re-pin this dispute")
	}
	if ParseImplementerReassignment(meta, "claude-1", "zcode-1") {
		t.Fatal("reassignment naming a different new id must not re-pin this dispute")
	}
}
