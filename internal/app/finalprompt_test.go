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

// The contract @codex-1 asked for after demonstrating that the prompt test could pass while the
// production gate required something the prompt never mentioned: build a FINAL out of NOTHING BUT
// what the prompt tells the drafter to write, and hand it to the gate.
//
// If a section is ever added to the gate and not to the prompt, the constructed FINAL will lack it
// and this fails. Two authorities cannot both be right; there is now only one, and this is the
// test that notices if a second one comes back.
func TestAFinalBuiltFromThePromptSatisfiesTheProductionGate(t *testing.T) {
	prompt := buildFinalDraftPrompt("/deck/ideas/demo", "/deck/ideas/demo/FINAL.md")

	var headings []string
	for _, line := range strings.Split(prompt, "\n") {
		if trimmed := strings.TrimSpace(line); strings.HasPrefix(trimmed, "## ") {
			headings = append(headings, trimmed)
		}
	}
	if len(headings) == 0 {
		t.Fatal("the prompt names no sections at all; a drafter reading it cannot satisfy any gate")
	}

	// The FINAL a compliant drafter would produce from this prompt, and nothing more.
	var b strings.Builder
	b.WriteString("---\nidea: demo\nstatus: final\n---\n\n")
	for _, h := range headings {
		b.WriteString(h + "\n\nA concrete line of agreed design.\nA second concrete line.\nA third concrete line.\n\n")
	}

	if reason := protocol.ValidateFinal(b.String(), "demo"); reason != "" {
		t.Fatalf("the gate rejects a FINAL written exactly as the prompt instructs: %s\n\nprompt:\n%s", reason, prompt)
	}

	// And the reverse direction: the prompt must not omit a section the gate requires.
	for _, want := range protocol.RequiredFinalSections {
		var found bool
		for _, h := range headings {
			if h == want {
				found = true
			}
		}
		if !found {
			t.Errorf("the gate requires %q but the prompt never names it", want)
		}
	}
}
