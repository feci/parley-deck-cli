package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocolcore"
)

// overlayFile builds a valid overlay for a payload written against coreBody.
func overlayFile(id, rationale, payload, coreBody string) string {
	return "---\n" +
		"schema: " + protocolcore.OverlaySchema + "\n" +
		"operations:\n" +
		"  - id: " + id + "\n" +
		"    kind: extend\n" +
		"    rationale: " + rationale + "\n" +
		"    core-sha256: " + protocolcore.Hash(coreBody) + "\n" +
		"    markdown: |-\n" +
		indent(payload, "      ") + "\n" +
		"---\n"
}

func indent(s, pad string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		lines[i] = pad + l
	}
	return strings.Join(lines, "\n")
}

// overlayFixture publishes a core, writes a deck, and adopts both core and overlay in the lock.
func overlayFixture(t *testing.T, coreBody, deckBody, overlayRaw string) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("PARLEY_HOME", home)
	if _, err := protocolcore.StoreAt(home).Publish("1.0.0", coreBody); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	meta := filepath.Join(root, "parley-deck", "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "COOPERATION.md"), []byte(deckBody), 0o644); err != nil {
		t.Fatal(err)
	}
	overlayHash := protocolcore.OverlayNone
	if overlayRaw != "" {
		p := filepath.Join(root, "parley-deck", protocolcore.OverlayFileName)
		if err := os.WriteFile(p, []byte(overlayRaw), 0o644); err != nil {
			t.Fatal(err)
		}
		ov, err := protocolcore.ParseOverlay(overlayRaw)
		if err != nil {
			t.Fatalf("fixture overlay does not parse: %v", err)
		}
		overlayHash = ov.SHA256
	}
	if err := os.WriteFile(filepath.Join(meta, "protocol-lock.yaml"),
		lockV2("1.0.0", coreBody, overlayHash), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

const ovCore = `# T

**Workspace:** ` + "`x`" + `
**Created:** 2026-01-01
**Transport:** ` + "`local-dir`" + `

## 1. One

core one.

## 2. Two

core two.
`

const ovPayload = `## Project-local procedures

This deck keeps a local note.`

// G1/G3: an adopted overlay is rendered at the terminal boundary, and rendering twice is
// byte-identical with an empty report.
func TestOverlayRendersAtTheTerminalBoundaryAndIsIdempotent(t *testing.T) {
	raw := overlayFile("deck.local", "local procedures", ovPayload, ovCore)
	root := overlayFixture(t, ovCore, "", raw)
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")

	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("render exit=%d: %s", code, errb.String())
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), "## Project-local procedures") {
		t.Fatalf("overlay payload missing from render:\n%s", first)
	}
	// Terminal boundary: the payload is after ALL core content, not wedged between sections.
	if strings.Index(string(first), "## Project-local procedures") < strings.Index(string(first), "core two.") {
		t.Error("payload was not rendered at the terminal boundary")
	}

	out.Reset()
	errb.Reset()
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("second render exit=%d: %s", code, errb.String())
	}
	second, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Error("render is not idempotent with an overlay")
	}
	if !strings.Contains(out.String(), "nothing to do") {
		t.Errorf("second render should be a no-op, got: %s", out.String())
	}
}

// G1: each reconciliation failure BLOCKS. Table-driven so a newly-added state cannot quietly
// default to "render anyway".
func TestOverlayMismatchesBlockBeforeComposition(t *testing.T) {
	raw := overlayFile("deck.local", "local procedures", ovPayload, ovCore)

	cases := []struct {
		name string
		// mutate is applied to a fully-adopted fixture to create the failure state.
		mutate func(t *testing.T, root string)
		want   string
	}{
		{
			name: "declared but the file is absent",
			mutate: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "parley-deck", protocolcore.OverlayFileName)); err != nil {
					t.Fatal(err)
				}
			},
			want: "is absent",
		},
		{
			name: "file changed after the lock was written",
			mutate: func(t *testing.T, root string) {
				edited := overlayFile("deck.local", "local procedures", ovPayload+"\n\nedited.", ovCore)
				if err := os.WriteFile(filepath.Join(root, "parley-deck", protocolcore.OverlayFileName),
					[]byte(edited), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: "the overlay changed without the lock being updated",
		},
		{
			name: "file is empty",
			mutate: func(t *testing.T, root string) {
				if err := os.WriteFile(filepath.Join(root, "parley-deck", protocolcore.OverlayFileName),
					[]byte(""), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: "empty",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := overlayFixture(t, ovCore, "", raw)
			tc.mutate(t, root)
			var out, errb strings.Builder
			code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb)
			if code == 0 {
				t.Fatalf("rendered despite %s", tc.name)
			}
			if !strings.Contains(errb.String(), tc.want) {
				t.Errorf("blocked for the wrong reason: want %q, got %s", tc.want, errb.String())
			}
		})
	}
}

// G1: the reverse direction — an overlay on disk that the lock never adopted must also block,
// rather than being silently ignored (which would render the deck WITHOUT its local content).
func TestUnadoptedOverlayBlocks(t *testing.T) {
	raw := overlayFile("deck.local", "local procedures", ovPayload, ovCore)
	root := overlayFixture(t, ovCore, "", "")
	if err := os.WriteFile(filepath.Join(root, "parley-deck", protocolcore.OverlayFileName),
		[]byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code == 0 {
		t.Fatal("rendered with an overlay the lock never adopted")
	}
	if !strings.Contains(errb.String(), "adopt it explicitly") {
		t.Errorf("wrong reason: %s", errb.String())
	}
}

// G2: a pre-v2 flat lock is refused with a migration message, not silently honoured.
func TestFlatLockIsRefusedWithMigrationGuidance(t *testing.T) {
	root := overlayFixture(t, ovCore, "", "")
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "meta", "protocol-lock.yaml"),
		[]byte("core-version: 1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code == 0 {
		t.Fatal("accepted a pre-v2 flat lock")
	}
	if !strings.Contains(errb.String(), "pre-v2 flat lock") {
		t.Errorf("wrong reason: %s", errb.String())
	}
}

// The lock attests the core's BYTES. A release republished under the same version must not render.
func TestLockBodyHashMismatchBlocks(t *testing.T) {
	root := overlayFixture(t, ovCore, "", "")
	// Attest a body that is not the one installed.
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "meta", "protocol-lock.yaml"),
		lockV2("1.0.0", ovCore+"\ndrifted.\n", ""), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code == 0 {
		t.Fatal("rendered a core the lock does not vouch for")
	}
	if !strings.Contains(errb.String(), "does not vouch for") {
		t.Errorf("wrong reason: %s", errb.String())
	}
}

// G5: the relocation witness is a PROOF, not a similarity heuristic.
//
// Exact, complete and unique → relocated. Anything else stays removed. The near-miss half of this
// test is the important half: without it the witness could be a substring match and still pass.
func TestRelocationWitnessIsExactAndNearMissesStayRemoved(t *testing.T) {
	raw := overlayFile("deck.local", "local procedures", ovPayload, ovCore)

	// The deck already carries the payload, mid-document.
	deckWithPayload := "# T\n\n**Workspace:** `x`\n\n## 1. One\n\ncore one.\n\n" + ovPayload + "\n\n## 2. Two\n\ncore two.\n"
	root := overlayFixture(t, ovCore, deckWithPayload, raw)
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	if strings.Contains(out.String(), "Project-local procedures") &&
		strings.Contains(out.String(), "not carried forward") {
		t.Errorf("an exactly-relocated payload was reported as lost:\n%s", out.String())
	}

	// Near miss: the deck's copy differs by one character, so it is NOT the payload and must be
	// reported as removed rather than forgiven.
	edited := strings.Replace(deckWithPayload, "a local note", "a local NOTE", 1)
	root2 := overlayFixture(t, ovCore, edited, raw)
	out.Reset()
	errb.Reset()
	if code := runProtocol([]string{"render", "--dir", root2, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	if !strings.Contains(out.String(), "not carried forward") {
		t.Errorf("an edited near-miss was forgiven as a relocation:\n%s", out.String())
	}
}

// The overlay grammar is strict. Each case is a refusal the design requires; a permissive YAML
// reader would silently accept most of them.
func TestOverlayGrammarRefusals(t *testing.T) {
	valid := overlayFile("deck.local", "why", ovPayload, ovCore)
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"unknown key", strings.Replace(valid, "  - id:", "  - unexpected: x\n  - id:", 1), "field unexpected"},
		{"wrong schema", strings.Replace(valid, protocolcore.OverlaySchema, "parley.protocol-overlay/v9", 1), "schema"},
		{"replace is not a v1 kind", strings.Replace(valid, "kind: extend", "kind: replace", 1), "only \"extend\""},
		{"bad operation id", strings.Replace(valid, "id: deck.local", "id: local", 1), "deck.<slug>"},
		{"empty file", "", "empty"},
		{"free-form body after the fence", valid + "\ntrailing prose\n", "free-form body"},
		{"short core hash", strings.Replace(valid, protocolcore.Hash(ovCore), "abc123", 1), "64 lowercase hex"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := protocolcore.ParseOverlay(tc.raw)
			if err == nil {
				t.Fatalf("accepted %s", tc.name)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("rejected for the wrong reason: want %q, got %v", tc.want, err)
			}
		})
	}
}

// The lock hash covers the WHOLE overlay file, not just the payload. codex-1 blocked consensus over
// this: hashing only the payload would let a rationale or id change leave the lock unchanged, so
// the lock would attest an overlay that is no longer the one on disk.
func TestLockHashCoversMetadataNotJustThePayload(t *testing.T) {
	a := overlayFile("deck.local", "first reason", ovPayload, ovCore)
	b := overlayFile("deck.local", "second reason", ovPayload, ovCore)
	ovA, err := protocolcore.ParseOverlay(a)
	if err != nil {
		t.Fatal(err)
	}
	ovB, err := protocolcore.ParseOverlay(b)
	if err != nil {
		t.Fatal(err)
	}
	if ovA.Operations[0].PayloadSHA256 != ovB.Operations[0].PayloadSHA256 {
		t.Fatal("fixture is wrong: the payloads should be identical")
	}
	if ovA.SHA256 == ovB.SHA256 {
		t.Error("changing the rationale left the lock hash unchanged; the hash is not covering the whole file")
	}
}

// `protocol overlay` is the documented way to discover the overlay's state, so it must succeed on
// a deck that has none. The equivalent bug shipped once already: `protocol --help` exited 2.
func TestProtocolOverlayReportsAbsenceAsValid(t *testing.T) {
	root := overlayFixture(t, ovCore, "", "")
	var out, errb strings.Builder
	if code := runProtocol([]string{"overlay", "show", "--dir", root}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	if !strings.Contains(out.String(), "absent") {
		t.Errorf("did not report absence: %s", out.String())
	}
}
