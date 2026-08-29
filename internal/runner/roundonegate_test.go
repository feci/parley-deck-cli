package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

// These tests bind the §15.6(a) acquisition gate to the RUNTIME CALL SITE, not to the helper.
//
// The first implementation of this rule shipped protocol.ValidateRoundOneArtifact with a passing
// mutation test and zero non-test callers: the runtime path used this file's same-named function,
// which checked only the four legacy sections. A round-1 artifact with no "## Existing alternatives"
// validated clean. Two reviewers found it independently. A test that exercises the helper proves the
// helper rejects; it proves nothing about what runs.
const roundOneFixture = "---\nagent: hermes-1\nidea: demo\nround: 1\n---\n\n" +
	"## Summary\nok\n\n## Proposed approach\nbuild it by hand\n\n%s" +
	"## Concerns / open questions\nnone\n\n## Risks\nnone\n"

func writeRoundOneFixture(t *testing.T, alternatives string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "hermes-1.md")
	body := strings.Replace(roundOneFixture, "%s", alternatives, 1)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRuntimeValidateRoundOneArtifactEnforcesExistingAlternatives(t *testing.T) {
	// Absent: an otherwise fully compliant artifact must be REJECTED by the runtime validator.
	// This is the exact fixture that passed with err=nil before the fix.
	p := writeRoundOneFixture(t, "")
	if err := ValidateRoundOneArtifact(p, "hermes-1", "demo"); err == nil {
		t.Fatal("runtime validator accepted a round-1 artifact with no Existing alternatives section")
	} else if !strings.Contains(err.Error(), "Existing alternatives") {
		t.Fatalf("error must name the missing section, got: %v", err)
	}

	// A bare heading is a rubber-stamp, not an enumerated search.
	p = writeRoundOneFixture(t, "## Existing alternatives\n\n")
	if err := ValidateRoundOneArtifact(p, "hermes-1", "demo"); err == nil {
		t.Fatal("runtime validator accepted an empty Existing alternatives section")
	}

	// A substring mention in prose is not the section.
	p = writeRoundOneFixture(t, "I looked at ## Existing alternatives while writing.\n\n")
	if err := ValidateRoundOneArtifact(p, "hermes-1", "demo"); err == nil {
		t.Fatal("runtime validator accepted a substring mention as the section")
	}

	// Enumerated content validates.
	p = writeRoundOneFixture(t, "## Existing alternatives\nBuilt by hand: dependency copying. Ships already: `pnpm deploy` (pnpm CLI). Constraint-forced: no.\n\n")
	if err := ValidateRoundOneArtifact(p, "hermes-1", "demo"); err != nil {
		t.Fatalf("an enumerated section must validate: %v", err)
	}

	// The ratified null-result form validates too.
	p = writeRoundOneFixture(t, "## Existing alternatives\nSearched `pnpm --help`, the Node stdlib and the lockfile; nothing ships this. The hand-built route is correct.\n\n")
	if err := ValidateRoundOneArtifact(p, "hermes-1", "demo"); err != nil {
		t.Fatalf("a scoped null result must validate: %v", err)
	}
}

// The prompt that asks for the section and the gate that requires it must not drift apart: an agent
// told nothing about the duty, then failed for omitting it, is a trap rather than a rule.
func TestRoundOnePromptCarriesTheSectionItIsGatedOn(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte("---\nidea: demo\n---\n\n## Problem / idea\nx\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	prompt, err := BuildRoundOnePrompt(
		agents.Discovery{Spec: agents.Spec{ID: "hermes-1"}},
		protocol.IdeaStatus{Path: dir, Slug: "demo"},
		"task", filepath.Join(dir, "round-01", "hermes-1.md"), filepath.Join(dir, "questions"),
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"## Existing alternatives", "ALREADY SHIPS", "constraint-forced"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("round-1 prompt must carry %q so the gate is reachable by an agent that reads it", want)
		}
	}
}
