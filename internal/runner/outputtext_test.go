package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

func TestOutputTextPolicyAcrossRunnerPrompts(t *testing.T) {
	root := t.TempDir()
	idea := protocol.IdeaStatus{Slug: "ste-example", Path: root, Participants: []string{"writer", "peer"}}
	agent := agents.Discovery{Spec: agents.Spec{ID: "writer"}}
	// Preserve exact evidence, including CRLF, percent signs, a non-English quote,
	// frontmatter, and code. The output policy must never transform its inputs.
	evidence := "---\r\nste_output: false\r\nreason: exact quotation\r\n---\r\n> pridaj všetko\n```sh\nprintf '%s' '/tmp/a b'\n```\nERROR: source_key=100%\n"
	if err := os.WriteFile(filepath.Join(root, "00-prompt.md"), []byte(evidence), 0600); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(root, "output with spaces.md")
	one, err := BuildRoundOnePrompt(agent, idea, "", out, root)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, prompt string
		preserve     []string
	}{
		{"round-one", one, []string{evidence, out, "round: 1", "## Existing alternatives"}},
		{"round-two", BuildRoundPrompt(agent, idea, 2, out, root, evidence), []string{evidence, out, "### @peer", "responding-to:"}},
		{"implementation", BuildImplementationPrompt(agent, idea, out), []string{out, filepath.Join(root, "FINAL.md"), "status: implemented", "## Summary of work"}},
		{"fix-up", BuildFixupPrompt(agent, idea), []string{filepath.Join(root, "review", "consensus.md"), filepath.Join(root, "IMPLEMENTATION.md"), "## Fix-up cycle"}},
		{"review", BuildReviewPrompt(agent, idea, 1, out, evidence), []string{evidence, out, "reviewed-commit:", "## Refutation attempts"}},
		{"review-consensus", BuildReviewConsensusPrompt(agent, idea, out, evidence, true), []string{evidence, out, "outstanding_agreed_fixes: 0", "strict_gate_clean: false", "## Signoffs"}},
		{"consult", BuildConsultPrompt(agent, evidence), []string{evidence, "Do NOT create, modify, or delete any files"}},
		{"goal-check", BuildGoalCheckPrompt(agent, idea), []string{filepath.Join(root, "FINAL.md"), "GOAL-CHECK: PASS", "GOAL-CHECK: FAIL —"}},
		{"steer", BuildSteerPrompt(agent, Options{Root: root, Idea: idea, RunID: "example"}, "steer-1", evidence, out), []string{evidence, out, "NOT a protocol round"}},
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
