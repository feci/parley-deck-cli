package app

import (
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

func TestOutputTextPolicyAcrossDriverPrompts(t *testing.T) {
	root := t.TempDir()
	idea := filepath.Join(root, "parley-deck", "ideas", "exact-100%")
	consensusPath := filepath.Join(idea, "consensus.md")
	finalPath := filepath.Join(idea, "FINAL.md")
	reviewPath := filepath.Join(idea, "review", "consensus.md")
	command := "'/tmp/a b/parley' trajectory verify-helper --request '/tmp/príklad.json' --request-sha256 'a1%2'"
	cases := []struct {
		name, prompt string
		preserve     []string
	}{
		{"consensus", buildConsensusDraftPrompt(idea, consensusPath), append([]string{consensusPath, "<!-- <agent-id> appends its signoff below -->"}, protocol.RequiredConsensusSections...)},
		{"final", buildFinalDraftPrompt(idea, finalPath), append([]string{finalPath, "idea: exact-100%", "status: final"}, protocol.RequiredFinalSections...)},
		{"design-signoff", buildConsensusSignoffPrompt(root, consensusPath, false, "writer"), []string{consensusPath, "### Signoff: writer —", "Status: ✅ ACCEPT", "❌ BLOCK", "🟡 ACCEPT-WITH-RESERVATIONS"}},
		{"review-signoff", buildConsensusSignoffPrompt(root, reviewPath, true, "writer"), []string{reviewPath, filepath.Join(idea, "00-prompt.md"), "### Signoff: writer —", "Status: ✅ ACCEPT"}},
		{"trajectory", trajectoryVerifierPrompt("writer", "exact-100%", command), []string{command, "PARLEY_RUN_ID, PARLEY_AGENT_ID and PARLEY_PROC_MARKER", "do not retry"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if strings.Count(tc.prompt, protocol.OutputTextInstruction) != 1 || !strings.HasSuffix(tc.prompt, protocol.OutputTextInstruction) {
				t.Fatal("the real builder must append the shared output policy exactly once")
			}
			for _, exact := range tc.preserve {
				if !strings.Contains(tc.prompt, exact) {
					t.Errorf("lost exact input or required syntax: %q", exact)
				}
			}
		})
	}
}
