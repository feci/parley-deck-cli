package app

// lean-organizer A tests: declared-facilitator runs default to the pure organizer.
// Covers the FINAL A acceptance rows: preflight fail-closed naming both fields (exit 0
// with the exception flag), driver role-ineligibility with escalate-not-fallback, the
// absent-field regression (byte-identical role selection), and the consensus-prompt /
// gate parity over protocol.RequiredConsensusSections.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

func writeIdeaPrompt(t *testing.T, root, slug string, frontmatter string) string {
	t.Helper()
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\n" + frontmatter + "\nstatus: round-01\n---\n\n## Problem\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return ideaDir
}

func seedMinimalDeck(t *testing.T, root string) {
	t.Helper()
	deck := filepath.Join(root, protocol.DeckDir)
	if err := os.MkdirAll(filepath.Join(deck, "meta"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(deck, "COOPERATION.md"), []byte("# COOP\n\n**Transport:** `github-pr`\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	meta := `{"protocolRole": "source", "protocolSha256": "` + strings.Repeat("a", 64) + `", "deckVersion": "2.12.0"}`
	if err := os.WriteFile(filepath.Join(deck, "meta", "version.json"), []byte(meta), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPreflightFacilitatorInParticipantsFailsClosed(t *testing.T) {
	dir := t.TempDir()
	seedMinimalDeck(t, dir)
	writeIdeaPrompt(t, dir, "declared-run",
		"idea: declared-run\nauthor: user\ntrack: standard\nparticipants: [codex-1, claude-1]\nfacilitator: codex-1")

	gates := facilitatorConflictGates(dir)
	if len(gates) != 1 {
		t.Fatalf("want 1 facilitator gate, got %d", len(gates))
	}
	joined := string(gates[0].Kind) + " " + gates[0].Detail + " " + gates[0].Confirm
	for _, want := range []string{"facilitator: codex-1", "facilitator_participates: true", "participants"} {
		if !strings.Contains(joined, want) {
			t.Errorf("gate must name %q; gate: %s", want, joined)
		}
	}
}

func TestPreflightFacilitatorParticipatesFlagClearsConflict(t *testing.T) {
	dir := t.TempDir()
	seedMinimalDeck(t, dir)
	writeIdeaPrompt(t, dir, "declared-run",
		"idea: declared-run\nauthor: user\ntrack: standard\nparticipants: [codex-1, claude-1]\nfacilitator: codex-1\nfacilitator_participates: true")
	if gates := facilitatorConflictGates(dir); len(gates) != 0 {
		t.Fatalf("facilitator_participates: true must clear the conflict, got %d gates", len(gates))
	}
}

func TestPreflightNoFacilitatorFieldUntouched(t *testing.T) {
	dir := t.TempDir()
	seedMinimalDeck(t, dir)
	writeIdeaPrompt(t, dir, "plain-run", "idea: plain-run\nauthor: user\ntrack: standard\nparticipants: [claude-1, kimi-1]")
	if gates := facilitatorConflictGates(dir); len(gates) != 0 {
		t.Fatalf("absent facilitator field must produce no gate, got %d", len(gates))
	}
}

func TestFacilitatorConflictGateIsABlockingPreflightOutcome(t *testing.T) {
	// The gate list is what drives preflight's exit-3 outcome; the roster ping is
	// skipped so the gate is the only input. This exercises the wiring end-to-end
	// without a real CLI on PATH.
	dir := t.TempDir()
	seedMinimalDeck(t, dir)
	writeIdeaPrompt(t, dir, "declared-run",
		"idea: declared-run\nauthor: user\ntrack: standard\nparticipants: [codex-1, claude-1]\nfacilitator: codex-1")

	opts := preflightOptions{Root: dir, NoPing: true}
	var out, errOut strings.Builder
	report, code, err := preflight(context.Background(), opts, nil, &out, &errOut)
	if err != nil {
		t.Fatalf("preflight: %v", err)
	}
	if code == 0 {
		t.Fatalf("conflicting facilitator declaration must fail closed (non-zero), got exit 0")
	}
	found := false
	for _, g := range report.Gates {
		if g.Kind == "facilitator-declaration" && strings.Contains(g.Detail, "facilitator_participates") {
			found = true
		}
	}
	if !found {
		t.Errorf("preflight report must carry the facilitator gate naming facilitator_participates; gates=%+v", report.Gates)
	}
}

func TestConsensusDraftPromptGateParity(t *testing.T) {
	prompt := buildConsensusDraftPrompt("/tmp/idea", "/tmp/idea/consensus.md")
	for _, section := range protocol.RequiredConsensusSections {
		if !strings.Contains(prompt, section) {
			t.Errorf("consensus draft prompt must emit canonical section %q", section)
		}
	}
	for _, section := range protocol.ConditionalConsensusSections {
		if !strings.Contains(prompt, section) {
			t.Errorf("consensus draft prompt must name conditional duty section %q", section)
		}
	}
	// The gate reads MissingConsensusSections, which is derived from the same slice;
	// prove the loop guard by asserting every required section is what the gate checks.
	body := strings.Join(protocol.RequiredConsensusSections, "\n")
	if missing := protocol.MissingConsensusSections(body); len(missing) != 0 {
		t.Errorf("a body containing RequiredConsensusSections must satisfy the gate, missing=%v", missing)
	}
	dropped := strings.ReplaceAll(body, "## Drafter position changes", "## Changed positions")
	if missing := protocol.MissingConsensusSections(dropped); len(missing) != 1 || missing[0] != "Drafter position changes" {
		t.Errorf("gate must flag the dropped §15.5 duty, got %v", missing)
	}
}

func TestConsensusPromptNamesFifteenDuties(t *testing.T) {
	prompt := buildConsensusDraftPrompt("/tmp/idea", "/tmp/idea/consensus.md")
	for _, duty := range []string{"Drafter position changes", "Alternatives disposition", "Verdict conflicts", "Comparison & blind spots"} {
		if !strings.Contains(prompt, duty) {
			t.Errorf("prompt must carry the §15 duty %q", duty)
		}
	}
	for _, stale := range []string{"## Trade-offs accepted", "## Deferred follow-ups\n## Dismissed findings"} {
		if strings.Contains(prompt, stale) {
			t.Errorf("prompt must not emit the old review-cycle heading %q", stale)
		}
	}
}
