package app

// lean-organizer A tests: declared-facilitator runs default to the pure organizer.
// Covers the FINAL A acceptance rows: preflight fail-closed naming both fields (exit 0
// with the exception flag), driver role-ineligibility with escalate-not-fallback, the
// absent-field regression (byte-identical role selection), and the consensus-prompt /
// scaffold parity over protocol.RequiredConsensusSections (fix-up F1: no status gate
// reads the constant — append-only signoffs remain the only machine-validated gate).

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/consensus"
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

	gates, err := facilitatorConflictGates(dir)
	if err != nil {
		t.Fatalf("conflict gates: %v", err)
	}
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
	if gates, err := facilitatorConflictGates(dir); err != nil || len(gates) != 0 {
		t.Fatalf("facilitator_participates: true must clear the conflict, got %d gates", len(gates))
	}
}

func TestPreflightNoFacilitatorFieldUntouched(t *testing.T) {
	dir := t.TempDir()
	seedMinimalDeck(t, dir)
	writeIdeaPrompt(t, dir, "plain-run", "idea: plain-run\nauthor: user\ntrack: standard\nparticipants: [claude-1, kimi-1]")
	if gates, err := facilitatorConflictGates(dir); err != nil || len(gates) != 0 {
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

// TestPreflightFailsClosedWhenWorkspaceStatusUnreadable (fix-up G10, kimi-1 K2-F2):
// a tree whose parley-deck/ exists but whose COOPERATION.md is missing must exit 1
// naming the read failure — at the reviewed HEAD this shape printed "Ready: no
// pending gates" and exited 0 while the facilitator-conflict gate silently never
// ran (facilitatorConflictGates returned nil on the ReadWorkspaceStatus error).
func TestPreflightFailsClosedWhenWorkspaceStatusUnreadable(t *testing.T) {
	dir := t.TempDir()
	// The kimi-1 fixture shape: agents.toml + meta/version.json + a conflicting
	// idea prompt, but NO COOPERATION.md — workspaceExists passes, ReadWorkspaceStatus
	// cannot.
	deck := filepath.Join(dir, protocol.DeckDir)
	if err := os.MkdirAll(filepath.Join(deck, "meta"), 0o755); err != nil {
		t.Fatal(err)
	}
	meta := `{"protocolRole": "source", "protocolSha256": "` + strings.Repeat("a", 64) + `", "deckVersion": "2.12.0"}`
	if err := os.WriteFile(filepath.Join(deck, "meta", "version.json"), []byte(meta), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.toml"), []byte("[agents]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeIdeaPrompt(t, dir, "declared-run",
		"idea: declared-run\nauthor: user\ntrack: standard\nparticipants: [codex-1, claude-1]\nfacilitator: codex-1")

	opts := preflightOptions{Root: dir, NoPing: true}
	var out, errOut strings.Builder
	_, code, err := preflight(context.Background(), opts, nil, &out, &errOut)
	if err == nil {
		t.Fatalf("unreadable workspace status must be a hard preflight error, got nil (out=%q err=%q)", out.String(), errOut.String())
	}
	if code != 1 {
		t.Fatalf("unreadable workspace status must fail closed with exit 1, got %d", code)
	}
	joined := err.Error() + " " + errOut.String()
	if !strings.Contains(joined, "cannot read workspace status") || !strings.Contains(joined, "COOPERATION.md") {
		t.Errorf("the failure must name the read failure and COOPERATION.md, got err=%q errOut=%q", err, errOut.String())
	}
	if strings.Contains(out.String(), "Ready") {
		t.Errorf("no readiness summary may follow an unreadable workspace (the report is suppressed on hard error); out=%q", out.String())
	}
}

// TestConsensusDraftPromptScaffoldParity proves the fix-up-F1 parity property: the
// drafting prompt and the scaffold generator read ONE constant
// (protocol.RequiredConsensusSections + ConditionalConsensusSections), so the prompt
// can never instruct a drafter to produce an artifact the scaffold does not shape —
// while NO status gate requires the sections of an existing consensus document
// (append-only signoffs remain the only machine-validated gate).
func TestConsensusDraftPromptScaffoldParity(t *testing.T) {
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

	// The scaffold half: a fixture deck with a complete round-01, drafted through the
	// shipped Draft entry point, must carry every constant section as a heading — and
	// no `## ` heading outside the constant set (the scaffold cannot invent sections).
	dir := t.TempDir()
	if err := protocol.InitWorkspace(dir); err != nil {
		t.Fatal(err)
	}
	writeIdeaPrompt(t, dir, "parity-idea",
		"idea: parity-idea\nauthor: user\ntrack: standard\nparticipants: [alpha-1, beta-1]")
	ideaDir := filepath.Join(dir, protocol.DeckDir, "ideas", "parity-idea")
	roundDir := filepath.Join(ideaDir, "round-01")
	if err := os.MkdirAll(roundDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, agent := range []string{"alpha-1", "beta-1"} {
		body := "---\nagent: " + agent + "\nidea: parity-idea\nround: 1\ndate: 2026-09-24\n---\n\n## Summary\nPosition.\n\n## Existing alternatives\nNone.\n\n## Proposed approach\nExtend.\n\n## Concerns / open questions\nNone.\n\n## Risks\nLow.\n"
		if err := os.WriteFile(filepath.Join(roundDir, agent+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := consensus.Draft(dir, "parity-idea", consensus.DraftOptions{By: "alpha-1"}); err != nil {
		t.Fatalf("consensus.Draft: %v", err)
	}
	scaffoldBytes, err := os.ReadFile(filepath.Join(ideaDir, "consensus.md"))
	if err != nil {
		t.Fatal(err)
	}
	scaffold := string(scaffoldBytes)
	inConstant := map[string]bool{}
	for _, section := range protocol.RequiredConsensusSections {
		inConstant[section] = true
		if !headingLinePresent(scaffold, section) {
			t.Errorf("scaffold must carry canonical section heading %q", section)
		}
	}
	for line := range strings.SplitSeq(scaffold, "\n") {
		line = strings.TrimRight(line, "\r")
		if !strings.HasPrefix(line, "## ") {
			continue
		}
		if !inConstant[line] {
			t.Errorf("scaffold emits heading %q which is not in RequiredConsensusSections (drift)", line)
		}
	}
	// And the prompt half stays parity-bound: every constant section must appear in
	// the prompt with its heading form, so prompt, scaffold and constant cannot drift.
	for _, section := range protocol.RequiredConsensusSections {
		if !strings.Contains(prompt, "\n"+section+"\n") && !strings.HasPrefix(prompt, section) {
			t.Errorf("prompt must list canonical section %q as a heading line", section)
		}
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
