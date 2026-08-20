package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

func disc(id string) agents.Discovery { return agents.Discovery{Spec: agents.Spec{ID: id}} }

// LE-1: the review prompt carries the adversarial posture + the Refutation attempts section.
func TestBuildReviewPromptRefutationDefault(t *testing.T) {
	p := BuildReviewPrompt(disc("rev"), protocol.IdeaStatus{Slug: "demo"}, 1, "/x/rev.md", "ctx")
	for _, want := range []string{"Refutation-default", "## Refutation attempts", "assume the implementation is WRONG"} {
		if !strings.Contains(p, want) {
			t.Fatalf("review prompt missing %q", want)
		}
	}
}

// LE-1: ValidateReviewArtifact rejects a review with no Refutation attempts section.
func TestValidateReviewArtifactRequiresRefutation(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "rev.md")
	if err := os.WriteFile(p, []byte("---\nagent: rev\nidea: demo\nreview-round: 1\n---\n\n## Findings\nnone\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateReviewArtifact(p, "rev", "demo", 1); err == nil {
		t.Fatal("a review without '## Refutation attempts' must be rejected")
	}
	// F5: a substring mention in prose (not a real heading) must be rejected.
	if err := os.WriteFile(p, []byte("---\nagent: rev\nidea: demo\nreview-round: 1\n---\n\nI did some ## Refutation attempts inline.\n## Findings\nnone\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateReviewArtifact(p, "rev", "demo", 1); err == nil {
		t.Fatal("a '## Refutation attempts' substring in prose must be rejected (F5)")
	}
	// F5: a real heading but an EMPTY section must be rejected.
	if err := os.WriteFile(p, []byte("---\nagent: rev\nidea: demo\nreview-round: 1\n---\n\n## Refutation attempts\n\n## Findings\nnone\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateReviewArtifact(p, "rev", "demo", 1); err == nil {
		t.Fatal("an empty '## Refutation attempts' section must be rejected (F5)")
	}
	// A real heading with content validates.
	if err := os.WriteFile(p, []byte("---\nagent: rev\nidea: demo\nreview-round: 1\nreviewed-commit: cafe123\n---\n\n## Refutation attempts\ntried X; held\n\n## Findings\nnone\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateReviewArtifact(p, "rev", "demo", 1); err != nil {
		t.Fatalf("a non-empty '## Refutation attempts' section must validate: %v", err)
	}
}

// LE-2: the strict_gate close fields appear in the consensus prompt only under strict_gate.
func TestReviewConsensusPromptStrictGateFields(t *testing.T) {
	on := BuildReviewConsensusPrompt(disc("d"), protocol.IdeaStatus{Slug: "demo"}, "/x/c.md", "ctx", true)
	for _, want := range []string{"closing_review_round", "strict_gate_clean"} {
		if !strings.Contains(on, want) {
			t.Fatalf("strict consensus prompt missing %q", want)
		}
	}
	off := BuildReviewConsensusPrompt(disc("d"), protocol.IdeaStatus{Slug: "demo"}, "/x/c.md", "ctx", false)
	if strings.Contains(off, "strict_gate_clean") {
		t.Fatal("non-strict consensus prompt must not include strict_gate_clean")
	}
}
