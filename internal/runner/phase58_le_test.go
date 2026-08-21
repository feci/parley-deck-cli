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

// codex-1/F23, MAJOR: the implementation gate accepted ANY non-empty status and matched
// "## Summary of work" as a SUBSTRING, so an artifact reading `status: banana` with a bare heading
// validated and the auto-driver could treat a non-reviewable scaffold as a finished Phase 5
// artifact. @zcode-1 CONFIRMED it; the consensus mis-recorded that verdict as PARTIAL-only and the
// finding fell off the fix list with no disposition anywhere.
//
// Measured before enforcing: of 72 IMPLEMENTATION.md files in this deck, the number newly rejected
// by this change is ZERO. The 31 that fail lack the heading entirely and failed the old gate too.
func TestValidateImplementationArtifactRejectsAnUnknownStatus(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "IMPLEMENTATION.md")
	write := func(body string) {
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("---\nidea: demo\nstatus: banana\n---\n\n## Summary of work\n\nReal work happened here.\n")
	err := ValidateImplementationArtifact(p, "demo")
	if err == nil {
		t.Fatal("`status: banana` validated as a finished Phase 5 artifact")
	}
	if !strings.Contains(err.Error(), "vocabulary") {
		t.Fatalf("refusal does not name the cause: %v", err)
	}

	// Every documented and in-use status must still pass — the gate must not reject live work.
	for _, status := range []string{"implemented", "complete", "fix-up-cycle-1", "fix-up-cycle-12", "ready-for-review", "in-progress"} {
		write("---\nidea: demo\nstatus: " + status + "\n---\n\n## Summary of work\n\nReal work happened here.\n")
		if err := ValidateImplementationArtifact(p, "demo"); err != nil {
			t.Errorf("status %q rejected: %v", status, err)
		}
	}
}

// The other half of F23: a heading is not a summary, and a prose mention is not a heading.
func TestValidateImplementationArtifactRejectsAnEmptyOrMentionedSummary(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "IMPLEMENTATION.md")

	empty := "---\nidea: demo\nstatus: complete\n---\n\n## Summary of work\n\n## Verification\n\nok\n"
	if err := os.WriteFile(p, []byte(empty), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateImplementationArtifact(p, "demo"); err == nil {
		t.Fatal("an empty '## Summary of work' validated")
	} else if !strings.Contains(err.Error(), "empty") {
		t.Fatalf("refusal does not name the cause: %v", err)
	}

	mention := "---\nidea: demo\nstatus: complete\n---\n\nI will add a ## Summary of work later.\n"
	if err := os.WriteFile(p, []byte(mention), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateImplementationArtifact(p, "demo"); err == nil {
		t.Fatal("a prose mention of the heading validated as the heading")
	}
}
