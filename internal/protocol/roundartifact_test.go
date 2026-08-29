package protocol

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRoundOne(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "claude-1.md")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

const roundOneHead = "---\nagent: claude-1\nidea: demo\nround: 1\n---\n\n## Summary\nok\n\n## Proposed approach\nbuild it\n\n"

func TestValidateRoundOneArtifactRequiresExistingAlternatives(t *testing.T) {
	// Missing entirely.
	p := writeRoundOne(t, roundOneHead+"## Risks\nnone\n")
	if err := ValidateRoundOneArtifact(p); err == nil {
		t.Fatal("a round-1 artifact without an Existing alternatives section must be rejected")
	} else if !strings.Contains(err.Error(), RoundOneRequiredSection) {
		t.Fatalf("error should name the missing section, got: %v", err)
	}

	// Present but empty -- a heading is not content.
	p = writeRoundOne(t, roundOneHead+"## Existing alternatives\n\n## Risks\nnone\n")
	if err := ValidateRoundOneArtifact(p); err == nil {
		t.Fatal("an empty Existing alternatives section must be rejected")
	}

	// Substring mention in prose is not the section.
	p = writeRoundOne(t, roundOneHead+"I reviewed the ## Existing alternatives inline.\n\n## Risks\nnone\n")
	if err := ValidateRoundOneArtifact(p); err == nil {
		t.Fatal("a substring mention must not satisfy the gate")
	}

	// Non-empty section validates, including the ratified null-result form.
	p = writeRoundOne(t, roundOneHead+"## Existing alternatives\nSearched `pnpm --help`, the Node stdlib and the lockfile; nothing ships this. The hand-built route is correct.\n\n## Risks\nnone\n")
	if err := ValidateRoundOneArtifact(p); err != nil {
		t.Fatalf("a non-empty section must validate, including a scoped null: %v", err)
	}
}

func TestValidateRoundOneArtifactUnreadable(t *testing.T) {
	if err := ValidateRoundOneArtifact(filepath.Join(t.TempDir(), "absent.md")); err == nil {
		t.Fatal("an unreadable artifact must fail closed")
	}
}
