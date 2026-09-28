package pipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func mustGateJSON(t *testing.T, g Gate) []byte {
	t.Helper()
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(data, '\n')
}

// AC-NAME-1: the FINAL example must hold on every OS — PathEscape escapes the
// edge's '>' and the encoding is applied inside GatePath. Unsafe edge IDs are
// refused by the GatePath backstop (next test), but the ENCODER beneath it is
// total: every input still maps to one safe component (defense in depth, and
// the property that makes the .tmp staging name safe by construction).
func TestGatePathUniversalEncoding(t *testing.T) {
	for _, tc := range []struct {
		edgeID string
		file   string
	}{
		{"block-a->block-b", "block-a-%3Eblock-b.gate.json"},
		{"spec->deploy", "spec-%3Edeploy.gate.json"},
		{"a b->c d", "a%20b-%3Ec%20d.gate.json"}, // spaces are cosmetic-safe: encodable
		{"CON->x", "CON-%3Ex.gate.json"},         // compound never a bare device name
	} {
		path, err := GatePath("deck", "slug", tc.edgeID)
		if err != nil {
			t.Fatalf("GatePath(%q) refused a safe edge ID: %v", tc.edgeID, err)
		}
		got := filepath.Base(path)
		if got != tc.file {
			t.Errorf("GatePath(%q) base = %q, want %q", tc.edgeID, got, tc.file)
		}
	}
	// The encoder is TOTAL — these inputs are refused by the GatePath
	// backstop (ADS colon, reserved names, trailing dot, separators), yet
	// still each encode to one safe, injective component.
	for _, tc := range []struct {
		edgeID string
		file   string
	}{
		{"a:b->c", "a%3Ab-%3Ec"},             // ADS colon escapes
		{"CON", "%43ON"},                     // bare reserved device: first char escaped
		{"con.x", "%63on.x"},                 // reserved even with extension, case-insensitive
		{"COM1", "%43OM1"},                   // COM1-9 reserved
		{"a->b.", "a-%3Eb%2E"},               // trailing dot guard
		{"../escape->x", "..%2Fescape-%3Ex"}, // separator never survives: no traversal
		{`..\escape->x`, "..%5Cescape-%3Ex"}, // backslash neither
		{"C:\\evil->x", "C%3A%5Cevil-%3Ex"},  // drive path neither
		{"../../other-idea->x", "..%2F..%2Fother-idea-%3Ex"},
	} {
		if got := encodeEdgeID(tc.edgeID); got != tc.file {
			t.Errorf("encodeEdgeID(%q) = %q, want %q", tc.edgeID, got, tc.file)
		}
	}
	// The .tmp staging name SaveGate derives from GatePath inherits the
	// encoding by construction (the hosted ERROR_INVALID_NAME failures were on
	// the staging name).
	safe, err := GatePath("d", "s", "a->b")
	if err != nil {
		t.Fatalf("GatePath refused a->b: %v", err)
	}
	if strings.ContainsAny(filepath.Base(safe), `><:/\`) {
		t.Fatalf("encoded gate name still carries an unsafe character: %q", filepath.Base(safe))
	}
}

// AC-NAME-3: both interpolation backstops refuse the unsafe classes on the
// raw slug/edge/block ID on every OS — no legacy exemption.
func TestBackstopsRefuseUnsafeIDs(t *testing.T) {
	for _, id := range []string{"../escape", `a\b`, "a:b", "CON", "trailing."} {
		if _, err := GatePath("deck", "slug", id); err == nil {
			t.Errorf("GatePath accepted unsafe edge ID %q", id)
		}
		if _, err := GatePath("deck", id, "a->b"); err == nil {
			t.Errorf("GatePath accepted unsafe slug %q", id)
		}
		if _, err := BlockWorkspace("deck", "slug", id); err == nil {
			t.Errorf("BlockWorkspace accepted unsafe block ID %q", id)
		}
		if _, err := BlockWorkspace("deck", id, "block"); err == nil {
			t.Errorf("BlockWorkspace accepted unsafe slug %q", id)
		}
	}
	if _, err := BlockWorkspace("deck", "slug", "block-1"); err != nil {
		t.Errorf("BlockWorkspace refused a safe ID: %v", err)
	}
}

// AC-NAME-3/4: unsafe classes are refused on the RAW ID — the executed IsLocal
// table (N4) showed the concatenated-path check wrongly accepts four of these
// five traversals; all are refused here, including ../../other-idea.
func TestUnsafeBlockIDsRefusedOnRawID(t *testing.T) {
	for _, id := range []string{
		"..", "../x", "x/..", "../../other-idea", // traversal table (N4)
		`..\x`, "a/b", `a\b`, "/abs", "C:\\abs", "C:abs", "a:b", // separators, absolute, drive, ADS
		"CON", "con", "COM1", "lpt9.gate", // reserved device names, even with extension
		"trailing.", "trailing ", "", // trailing dot/space, empty
	} {
		if err := checkBlockID(id); err == nil {
			t.Errorf("checkBlockID(%q) accepted an unsafe ID", id)
		}
	}
	// Cosmetic grandfathering: safe-but-non-allowlisted IDs stay loadable.
	for _, id := range []string{"my block", "café", "a~b", "block-1"} {
		if err := checkBlockID(id); err != nil {
			t.Errorf("checkBlockID(%q) refused a safe cosmetic ID: %v", id, err)
		}
	}
	// The strict creation allowlist refuses the cosmetic remainder too.
	for _, id := range []string{"my block", "café", "a~b", ".hidden", strings.Repeat("x", maxBlockIDLen+1)} {
		if err := ValidNewBlockID(id); err == nil {
			t.Errorf("validNewBlockID(%q) accepted a non-allowlisted new ID", id)
		}
	}
	for _, id := range []string{"block-1", "a.b_c-d", "COM10"} { // COM10 is not reserved
		if err := ValidNewBlockID(id); err != nil {
			t.Errorf("validNewBlockID(%q) refused an allowlisted ID: %v", id, err)
		}
	}
}

// AC-NAME-2: a pre-§D.4 deck holds gates under the legacy raw name; LoadGate
// falls back to it, the first rewrite retires the shadow (no stale HITL
// answer survives), and a mixed-name deck resolves both. The legacy name may
// carry '>' — representable on POSIX, where the full flow runs; on Windows a
// raw '>' name cannot exist (the hosted ERROR_INVALID_NAME class the encoding
// exists to fix), so that platform instead pins the premise: creating the
// legacy name must fail, an absent gate reads as not-found (never a hard
// error), and the encoded save/load round-trip works.
func TestLegacyGateFallbackShadowRetirementAndMixedDeck(t *testing.T) {
	deck := t.TempDir()
	newGate := Gate{ID: "a->b", PipelineSlug: "p", Edge: "a->b", Status: GateOpen, CreatedAt: time.Unix(0, 0)}
	legacyGate := Gate{ID: "c->d", PipelineSlug: "p", Edge: "c->d", Status: GateApproved, ApprovedBy: "user", CreatedAt: time.Unix(0, 0)}

	if runtime.GOOS == "windows" {
		if err := os.MkdirAll(filepath.Dir(legacyGatePath(deck, "p", "a->b")), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(legacyGatePath(deck, "p", "a->b"), []byte("x"), 0o644); err == nil {
			t.Fatal("legacy raw-'>' gate name was creatable on Windows — the W2 encoding premise is broken")
		}
		if _, ok, err := LoadGate(deck, "p", "a->b"); err != nil || ok {
			t.Fatalf("absent gate on Windows: ok=%v err=%v (must be not-found, never ERROR_INVALID_NAME)", ok, err)
		}
		resolved := newGate
		resolved.Resolve(true, "user", time.Unix(1, 0))
		if err := SaveGate(deck, resolved); err != nil {
			t.Fatalf("SaveGate: %v", err)
		}
		loaded, ok, err := LoadGate(deck, "p", "a->b")
		if err != nil || !ok || loaded.Status != GateApproved || loaded.ApprovedBy != "user" {
			t.Fatalf("encoded gate round-trip: loaded=%v ok=%v err=%v", loaded, ok, err)
		}
		return
	}

	// Write both under their LEGACY names (a pre-encoding deck).
	for _, g := range []Gate{newGate, legacyGate} {
		p := legacyGatePath(deck, "p", g.Edge)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, mustGateJSON(t, g), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Read fallback: both resolve through the encoded API.
	for _, g := range []Gate{newGate, legacyGate} {
		loaded, ok, err := LoadGate(deck, "p", g.Edge)
		if err != nil || !ok || loaded.Edge != g.Edge {
			t.Fatalf("legacy read fallback for %q: loaded=%v ok=%v err=%v", g.Edge, loaded, ok, err)
		}
	}

	// First rewrite retires the shadow: after SaveGate the legacy file is gone
	// and the encoded file answers with the new answer.
	resolved := newGate
	resolved.Resolve(true, "user", time.Unix(1, 0))
	if err := SaveGate(deck, resolved); err != nil {
		t.Fatalf("SaveGate: %v", err)
	}
	if _, err := os.Lstat(legacyGatePath(deck, "p", "a->b")); !os.IsNotExist(err) {
		t.Fatalf("legacy shadow survived its first rewrite: %v", err)
	}
	loaded, ok, err := LoadGate(deck, "p", "a->b")
	if err != nil || !ok || loaded.Status != GateApproved || loaded.ApprovedBy != "user" {
		t.Fatalf("encoded gate after rewrite: loaded=%v ok=%v err=%v", loaded, ok, err)
	}

	// Mixed-name deck: the untouched legacy gate still resolves; a fresh save
	// of it also retires its shadow.
	fresh := legacyGate
	fresh.Resolve(false, "other", time.Unix(2, 0))
	if err := SaveGate(deck, fresh); err != nil {
		t.Fatalf("SaveGate mixed deck: %v", err)
	}
	if _, err := os.Lstat(legacyGatePath(deck, "p", "c->d")); !os.IsNotExist(err) {
		t.Fatalf("mixed-deck legacy shadow survived: %v", err)
	}
	loaded, ok, err = LoadGate(deck, "p", "c->d")
	if err != nil || !ok || loaded.Status != GateRejected {
		t.Fatalf("mixed-deck reloaded gate: loaded=%v ok=%v err=%v", loaded, ok, err)
	}
}

// The §D.4 write backstop: unsafe edge IDs never reach the disk.
func TestSaveGateRefusesUnsafeEdgeID(t *testing.T) {
	deck := t.TempDir()
	g := Gate{ID: "a->../escape", PipelineSlug: "p", Edge: "a->../escape", Status: GateOpen, CreatedAt: time.Unix(0, 0)}
	if err := SaveGate(deck, g); err == nil {
		t.Fatal("SaveGate accepted a traversal edge ID")
	}
	g.Edge = "CON"
	g.ID = "CON"
	if err := SaveGate(deck, g); err == nil {
		t.Fatal("SaveGate accepted a reserved-name edge ID")
	}
}
