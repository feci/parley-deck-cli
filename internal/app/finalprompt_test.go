package app

import (
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

// Review round 1, @codex-1 MAJOR: the driver instructed its own FINAL drafter to produce an
// artifact its own gate rejects — the prompt named one section while finalScaffoldReason requires
// all seven, and never mentioned the idea: frontmatter the slug check reads.
func TestFinalDraftPromptDescribesWhatTheGateRequires(t *testing.T) {
	prompt := buildFinalDraftPrompt("/deck/parley-deck/ideas/my-idea", "/deck/parley-deck/ideas/my-idea/FINAL.md")

	for _, section := range protocol.RequiredFinalSections {
		if !strings.Contains(prompt, section) {
			t.Errorf("prompt does not ask for %q, which the gate requires", section)
		}
	}
	if !strings.Contains(prompt, "idea: my-idea") {
		t.Error("prompt does not tell the drafter which idea slug to declare; the gate checks it")
	}
	if !strings.Contains(prompt, "status: final") {
		t.Error("prompt dropped the status the gate requires")
	}
	// The N/A affordance must be stated, or a drafter faced with seven headings on a trivial idea
	// has no honest way to fill them.
	if !strings.Contains(prompt, "N/A") {
		t.Error("prompt does not offer N/A for the sections the protocol allows it for")
	}
}
