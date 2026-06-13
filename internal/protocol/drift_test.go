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
	n := 0
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, prefix) {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("%s: header line %q* appears %d times, want exactly 1 (drift guard fails closed)", name, prefix, n)
	}
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
