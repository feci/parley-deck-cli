package protocol

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// liveDeckPath is the in-repo canonical project deck, relative to this package
// directory (`go test` runs with the package directory as the working dir).
const liveDeckPath = "../../parley-deck/COOPERATION.md"

// The embedded bootstrap template (internal/protocol/defaults/COOPERATION.md)
// and the live project deck (parley-deck/COOPERATION.md) must stay identical
// EXCEPT for five deliberately project-specific zones. The drift guard below
// normalizes ONLY these zones and compares everything else byte-for-byte.
//
// If this test failed: a protocol edit landed in only one copy. Apply the SAME
// change to BOTH files. If the change is genuinely project-specific, extend the
// allowlist here deliberately (and explain why) rather than widening the
// normalizer.
const (
	workspacePrefix    = "**Workspace:** "
	createdPrefix      = "**Created:** "
	protocolSyncPrefix = "**Protocol synced:**"
	rosterSectionLine  = "## 2. Active agents (roster)"
	rosterHeaderLine   = "| Agent ID       | Workspace dir                       | Role          |"
	handleHeaderLine   = "| Agent ID       | Host handle    |"
)

// The exact bootstrap-shape lines the embedded default must hold (D2/D3). The
// drift guard ASSERTS these directly so that an in-allowlist edit — an
// illustrative roster row, a placeholder value changed while keeping the key
// prefix, or a stray `**Protocol synced:**` line — cannot be silently
// normalized away and pass.
const (
	embWorkspaceLine = "**Workspace:** `<workspace-name>`"
	embCreatedLine   = "**Created:** `<date> — created by parley init`"
	embTransportLine = "**Transport:** `github-pr`"
)

// TestEmbeddedDefaultMatchesLiveDeck is the anti-drift guard (idea
// embedded-default-protocol-resync, D5). It fails closed: a missing deck file or
// a missing/duplicated structural anchor is a hard failure, never a silent skip.
func TestEmbeddedDefaultMatchesLiveDeck(t *testing.T) {
	deckBytes, err := os.ReadFile(liveDeckPath)
	if err != nil {
		t.Fatalf("cannot read live deck %s — the drift guard fails closed rather than skip: %v", liveDeckPath, err)
	}
	deck := string(deckBytes)

	// Fail closed: each anchor must appear exactly once in each file so the
	// normalizer can never swallow a wider region than the allowlist intends.
	for _, f := range []struct{ name, text string }{
		{"embedded default", defaultCooperation},
		{"live deck", deck},
	} {
		assertExactLineOnce(t, f.name, f.text, rosterSectionLine)
		assertExactLineOnce(t, f.name, f.text, rosterHeaderLine)
		assertExactLineOnce(t, f.name, f.text, handleHeaderLine)
		assertLinePrefixOnce(t, f.name, f.text, workspacePrefix)
		assertLinePrefixOnce(t, f.name, f.text, createdPrefix)
	}

	// D2/D3: assert the embedded default actually holds the generic bootstrap
	// shape, turning the allowlisted (normalized-away) zones into asserted
	// invariants on the embedded side.
	assertEmbeddedBootstrapShape(t, defaultCooperation)
	// Provenance line: the live deck carries exactly one; the embedded default
	// (bootstrap template) must carry none.
	if n := countLinePrefix(deck, protocolSyncPrefix); n != 1 {
		t.Fatalf("live deck: %q* appears %d times, want exactly 1", protocolSyncPrefix, n)
	}
	if n := countLinePrefix(defaultCooperation, protocolSyncPrefix); n != 0 {
		t.Fatalf("embedded default: %q* appears %d times, want exactly 0 (D2 — bootstrap output carries no sync provenance)", protocolSyncPrefix, n)
	}

	embNorm := normalizeProtocol(t, "embedded default", defaultCooperation)
	deckNorm := normalizeProtocol(t, "live deck", deck)

	if len(embNorm) != len(deckNorm) {
		t.Fatalf("embedded default and live deck differ after normalizing the allowlisted zones "+
			"(embedded=%d lines, deck=%d lines).\nA protocol edit likely landed in only one copy — apply it to BOTH "+
			"internal/protocol/defaults/COOPERATION.md and parley-deck/COOPERATION.md.\n%s",
			len(embNorm), len(deckNorm), firstLineDiff(embNorm, deckNorm))
	}
	for i := range embNorm {
		if embNorm[i] != deckNorm[i] {
			t.Fatalf("embedded default and live deck diverge at normalized line %d, outside the allowlisted zones:\n"+
				"  embedded: %q\n  deck:     %q\n"+
				"Apply the change to BOTH copies, or extend the drift allowlist deliberately.",
				i+1, embNorm[i], deckNorm[i])
		}
	}
}

// normalizeProtocol canonicalizes exactly the five allowlisted zones: it drops
// any `**Protocol synced:**` line, flattens the Workspace/Created header values
// to a constant, and empties both §2 table bodies while keeping each table's
// header and separator rows. Everything else is returned verbatim.
func normalizeProtocol(t *testing.T, name, text string) []string {
	t.Helper()
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		switch {
		case strings.HasPrefix(line, protocolSyncPrefix):
			continue // deck-only sync provenance
		case strings.HasPrefix(line, workspacePrefix):
			out = append(out, workspacePrefix+"<normalized>")
		case strings.HasPrefix(line, createdPrefix):
			out = append(out, createdPrefix+"<normalized>")
		case line == rosterHeaderLine || line == handleHeaderLine:
			out = append(out, line)
			if i+1 >= len(lines) || !strings.HasPrefix(lines[i+1], "| -") {
				t.Fatalf("%s: table header %q is not immediately followed by a separator row (fail-closed)", name, line)
			}
			i++
			out = append(out, lines[i]) // separator row
			for i+1 < len(lines) && strings.HasPrefix(lines[i+1], "|") {
				i++ // empty the table body
			}
		default:
			out = append(out, line)
		}
	}
	return out
}

func assertExactLineOnce(t *testing.T, name, text, anchor string) {
	t.Helper()
	if n := countExactLines(text, anchor); n != 1 {
		t.Fatalf("%s: anchor %q appears %d times, want exactly 1 (drift guard fails closed)", name, anchor, n)
	}
}

func assertLinePrefixOnce(t *testing.T, name, text, prefix string) {
	t.Helper()
	if n := countLinePrefix(text, prefix); n != 1 {
		t.Fatalf("%s: header line %q* appears %d times, want exactly 1 (drift guard fails closed)", name, prefix, n)
	}
}

// assertEmbeddedBootstrapShape verifies the embedded default holds the exact
// D2/D3 bootstrap shape: the three genericized header lines verbatim, and empty
// §2 table bodies (header + separator retained).
func assertEmbeddedBootstrapShape(t *testing.T, text string) {
	t.Helper()
	for _, want := range []string{embWorkspaceLine, embCreatedLine, embTransportLine} {
		if n := countExactLines(text, want); n != 1 {
			t.Fatalf("embedded default must contain exactly the bootstrap line %q once (D2), found %d", want, n)
		}
	}
	assertEmptyTableBody(t, text, rosterHeaderLine)
	assertEmptyTableBody(t, text, handleHeaderLine)
}

// assertEmptyTableBody confirms the §2 table under the given header has its
// separator row and an empty body (the line after the separator is not a table
// row), so no illustrative roster member can ship in the bootstrap template.
func assertEmptyTableBody(t *testing.T, text, header string) {
	t.Helper()
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != header {
			continue
		}
		if i+1 >= len(lines) || !strings.HasPrefix(lines[i+1], "| -") {
			t.Fatalf("embedded default: table %q is not immediately followed by a separator row", header)
		}
		if i+2 < len(lines) && strings.HasPrefix(lines[i+2], "|") {
			t.Fatalf("embedded default: §2 table %q must have an empty body (D3), found row: %q", header, lines[i+2])
		}
		return
	}
	t.Fatalf("embedded default: table header %q not found", header)
}

func countLinePrefix(text, prefix string) int {
	n := 0
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, prefix) {
			n++
		}
	}
	return n
}

func countExactLines(text, anchor string) int {
	n := 0
	for _, line := range strings.Split(text, "\n") {
		if line == anchor {
			n++
		}
	}
	return n
}

func firstLineDiff(a, b []string) string {
	min := len(a)
	if len(b) < min {
		min = len(b)
	}
	for i := 0; i < min; i++ {
		if a[i] != b[i] {
			return fmt.Sprintf("first divergence at normalized line %d:\n  embedded: %q\n  deck:     %q", i+1, a[i], b[i])
		}
	}
	return fmt.Sprintf("the shorter file is a prefix of the longer; first extra line at %d", min+1)
}

// TestDefaultCooperationForInit checks the bootstrap output that `parley init`
// writes (idea embedded-default-protocol-resync, D6).
func TestDefaultCooperationForInit(t *testing.T) {
	out := defaultCooperationForInit()

	mustContain := []string{
		"**Transport:** `local-dir`",                                  // the swap happened
		"## 12. Pipeline blocks & action stages",                      // §12 carried
		"ratified by idea `meta-protocol-change-end-to-end-pipeline`", // §12 provenance kept
		"**Workspace:** `<workspace-name>`",                           // static placeholder
		"created by parley init",                                      // created-date placeholder
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("defaultCooperationForInit() output is missing %q", s)
		}
	}

	mustNotContain := []string{
		"**Transport:** `github-pr`", // source transport must be swapped out
		protocolSyncPrefix,           // sync provenance never belongs in bootstrap output
		"`codex`",                    // no project roster members
		"`claude`",
		"`agy`",
		"`hermes`",
	}
	for _, s := range mustNotContain {
		if strings.Contains(out, s) {
			t.Errorf("defaultCooperationForInit() output must not contain %q", s)
		}
	}
}

// TestNoSection2AsAStoreInstructions pins the authority cutover in TEXT.
//
// G3 required the cutover to be atomic: §2 became "a generated, non-authoritative view"
// in every protocol copy, but the bootstrap section and Appendix A still told a
// facilitator to record roster choices *in* §2 and to "Fill in §2 roster" by hand. Both
// statements shipped in the same document, so an agent following the instructions would
// hand-edit the file the same document calls generated. Phrases, not prose style, are
// what regressed — so phrases are what this test forbids.
func TestNoSection2AsAStoreInstructions(t *testing.T) {
	banned := []string{
		"mirrored in the §2 roster",
		"Fill in §2 roster",
		"Modify the active roster (§2).",
		"and the §2 roster)",
	}
	for _, tc := range []struct{ name, text string }{
		{"embedded default", defaultCooperation},
		{"live deck", readLiveDeck(t)},
	} {
		for _, phrase := range banned {
			if strings.Contains(tc.text, phrase) {
				t.Errorf("%s still instructs the reader to treat §2 as a store: %q\n"+
					"§2 is a generated view; membership lives in parley-deck/agents.toml. "+
					"Update the instruction to name `parley roster set` / `parley roster render`.",
					tc.name, phrase)
			}
		}
	}
}

// readLiveDeck loads the canonical project deck, failing closed like the guard above.
func readLiveDeck(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile(liveDeckPath)
	if err != nil {
		t.Fatalf("cannot read live deck %s: %v", liveDeckPath, err)
	}
	return string(b)
}
