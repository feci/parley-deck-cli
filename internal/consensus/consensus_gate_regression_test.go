package consensus

// lean-organizer fix-up F1 regression fixtures (review round-01 CRIT-1): the removed
// required-sections gate rejected the protocol's own published Phase-3 template and
// 79 of 80 existing deck consensus artifacts. These tests pin the restored behavior:
// (i) a consensus.md assembled verbatim from the protocol's own Phase-3 template
// reads `ready`; (ii) the live deck's malformed count never exceeds the pre-idea base
// count (9 of 80 at b37f7ef, measured independently by both reviewers).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// phaseThreeTemplateBlock extracts the indented consensus.md template block from the
// deck's COOPERATION.md §Phase 3 ("When discussion has converged, an agent creates
// `ideas/<slug>/consensus.md`:"), de-indented to the artifact form.
func phaseThreeTemplateBlock(t *testing.T) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "parley-deck", "COOPERATION.md"))
	if err != nil {
		t.Skipf("live deck protocol not available: %v", err)
	}
	lines := strings.Split(string(raw), "\n")
	start := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "When discussion has converged, an agent creates") {
			start = i + 1
			break
		}
	}
	if start < 0 {
		t.Fatal("Phase-3 template anchor not found in the deck protocol")
	}
	var out []string
	for _, line := range lines[start:] {
		if line == "" {
			out = append(out, "")
			continue
		}
		if !strings.HasPrefix(line, "    ") {
			break // first non-indented prose line ends the template block
		}
		out = append(out, strings.TrimPrefix(line, "    "))
	}
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	if len(out) == 0 {
		t.Fatal("Phase-3 template block is empty")
	}
	return strings.Join(out, "\n")
}

// TestProtocolPhaseThreeTemplateConsensusReady: the fixture claude-1 and kimi-1 both
// executed — a consensus.md verbatim from the protocol's own published Phase-3
// template, with every participant at ✅ ACCEPT — must triage ready.
func TestProtocolPhaseThreeTemplateConsensusReady(t *testing.T) {
	template := phaseThreeTemplateBlock(t)
	root := setupIdea(t, "tpl-idea", []string{"alpha-1", "beta-1"})
	// The template's frontmatter placeholders (`<slug>`, `<agent-id>`, `YYYY-MM-DD`)
	// are substituted with the fixture's real values; every template heading and body
	// line below them stays verbatim.
	body := strings.ReplaceAll(template, "<slug>", "tpl-idea")
	body += "\n\n### Signoff: alpha-1 — 2026-09-24\nStatus: ✅ ACCEPT\nNotes: ok\n\n### Signoff: beta-1 — 2026-09-24\nStatus: ✅ ACCEPT\nNotes: ok\n"
	writeFile(t, filepath.Join(root, "parley-deck", "ideas", "tpl-idea", "consensus.md"), body)
	summary, err := Status(root, "tpl-idea", false)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != TriageReady {
		t.Fatalf("the protocol's own Phase-3 template must read ready, got %q (errors=%v)", summary.Triage, summary.Errors)
	}
	// Document why the removed gate rejected it: the published template does not
	// carry the §15.5/§15.6 duty headings.
	for _, duty := range []string{"## Drafter position changes", "## Alternatives disposition"} {
		if strings.Contains(template, duty) {
			t.Errorf("unexpected: the published template now carries %q — revisit this regression test", duty)
		}
	}
}

// TestLiveDeckConsensusCorpusMalformedNotAboveBase sweeps every consensus.md in the
// live deck and pins the corpus malformed count at or below the pre-idea base count
// (9 of 80 at b37f7ef; the removed gate had driven it to 79 of 80).
func TestLiveDeckConsensusCorpusMalformedNotAboveBase(t *testing.T) {
	deckRoot := filepath.Join("..", "..")
	glob, err := filepath.Glob(filepath.Join(deckRoot, "parley-deck", "ideas", "*", "consensus.md"))
	if err != nil || len(glob) == 0 {
		t.Skipf("live deck corpus not available: %v", err)
	}
	const baseMalformed = 9
	malformed := 0
	checked := 0
	for _, path := range glob {
		slug := filepath.Base(filepath.Dir(path))
		summary, serr := Status(deckRoot, slug, false)
		if serr != nil {
			continue
		}
		checked++
		if summary.Triage == TriageMalformed {
			malformed++
			t.Logf("malformed: %s (%v)", slug, summary.Errors)
		}
	}
	if malformed > baseMalformed {
		t.Fatalf("deck consensus corpus malformed count %d exceeds the pre-idea base count %d (the F1 gate regression)", malformed, baseMalformed)
	}
	t.Logf("corpus: %d consensus.md checked, %d malformed (base %d)", checked, malformed, baseMalformed)
}
