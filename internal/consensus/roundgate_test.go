package consensus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeF(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

const goodArtifact = "---\nagent: a\nround: 1\n---\n\n## Summary\n\nSomething substantive.\n"

// codex-1/F1: a file containing only a newline used to satisfy a round, so `consensus draft`
// could advance an idea containing no participant analysis at all.
func TestBlankRoundArtifactDoesNotCountAsFiled(t *testing.T) {
	dir := t.TempDir()
	writeF(t, filepath.Join(dir, "alice.md"), goodArtifact)
	writeF(t, filepath.Join(dir, "bob.md"), "\n")

	got := missingRoundArtifacts(dir, []string{"alice", "bob"})
	if len(got) != 1 || !strings.HasPrefix(got[0], "bob.md") {
		t.Fatalf("want bob.md reported, got %v", got)
	}
	// The reason must be stated: "missing" and "blank" call for different actions.
	if !strings.Contains(got[0], "blank") {
		t.Errorf("reason not reported: %q", got[0])
	}
}

// Frontmatter alone is not an artifact either.
func TestFrontmatterOnlyArtifactDoesNotCountAsFiled(t *testing.T) {
	dir := t.TempDir()
	writeF(t, filepath.Join(dir, "alice.md"), "---\nagent: alice\n---\n")
	if got := missingRoundArtifacts(dir, []string{"alice"}); len(got) != 1 {
		t.Fatalf("frontmatter-only file must not count, got %v", got)
	}
}

// A body with no section heading is not a protocol artifact.
func TestHeadinglessArtifactDoesNotCountAsFiled(t *testing.T) {
	dir := t.TempDir()
	writeF(t, filepath.Join(dir, "alice.md"), "---\nagent: alice\n---\n\njust prose, no sections\n")
	got := missingRoundArtifacts(dir, []string{"alice"})
	if len(got) != 1 || !strings.Contains(got[0], "heading") {
		t.Fatalf("want a headings complaint, got %v", got)
	}
}

// An absent file still reports plainly, without a parenthetical reason.
func TestAbsentArtifactStillReportsAsMissing(t *testing.T) {
	dir := t.TempDir()
	got := missingRoundArtifacts(dir, []string{"ghost"})
	if len(got) != 1 || got[0] != "ghost.md" {
		t.Fatalf("want bare ghost.md, got %v", got)
	}
}

// A well-formed artifact passes; the gate must not reject real work.
func TestWellFormedArtifactCounts(t *testing.T) {
	dir := t.TempDir()
	writeF(t, filepath.Join(dir, "alice.md"), goodArtifact)
	if got := missingRoundArtifacts(dir, []string{"alice"}); len(got) != 0 {
		t.Fatalf("well-formed artifact rejected: %v", got)
	}
}

// codex-1/F2: §6 forbids the implementer reviewing its own work, so a review round must not
// expect its file. Requiring it made a compliant Phase 6 unreachable.
func TestReviewRoundDoesNotExpectTheImplementersFile(t *testing.T) {
	idea := t.TempDir()
	writeF(t, filepath.Join(idea, "IMPLEMENTATION.md"), "---\nimplementer: impl\n---\n\n## Work\n")
	participants := []string{"impl", "rev-a", "rev-b"}

	got := expectedRoundParticipants(idea, participants, true)
	for _, p := range got {
		if p == "impl" {
			t.Fatalf("review round still expects the implementer: %v", got)
		}
	}
	if len(got) != 2 {
		t.Fatalf("want the two reviewers, got %v", got)
	}

	// A DESIGN round expects everyone, implementer included.
	if len(expectedRoundParticipants(idea, participants, false)) != 3 {
		t.Fatal("design round must expect every participant")
	}
}

// Fail closed: if the implementer cannot be resolved, expect everyone rather than silently
// accepting a short round.
func TestUnresolvableImplementerExpectsEveryone(t *testing.T) {
	idea := t.TempDir()
	participants := []string{"a", "b"}
	if len(expectedRoundParticipants(idea, participants, true)) != 2 {
		t.Fatal("with no IMPLEMENTATION.md the full list must be expected")
	}
}

// The FINAL drafter is the fallback implementer when IMPLEMENTATION.md has no explicit field.
func TestFinalDrafterIsTheFallbackImplementer(t *testing.T) {
	idea := t.TempDir()
	writeF(t, filepath.Join(idea, "FINAL.md"), "---\ndrafted-by: dee\n---\n\n## Decision\n")
	got := expectedRoundParticipants(idea, []string{"dee", "eve"}, true)
	if len(got) != 1 || got[0] != "eve" {
		t.Fatalf("want only eve, got %v", got)
	}
}

// codex-1/F20: signoff-shaped headings anywhere in the document used to satisfy the consensus
// gate, so example or quoted text under an earlier section counted as real approval.
func TestSignoffsBeforeTheSignoffsHeadingDoNotCount(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "consensus.md")
	writeF(t, path, strings.Join([]string{
		"---", "idea: x", "---", "",
		"## Agreed decisions", "",
		"### Signoff: alice — 2026-01-01", "Status: accept", "",
		"### Signoff: bob — 2026-01-01", "Status: accept", "",
		"## Signoffs", "",
		"_nobody has signed here_", "",
	}, "\n"))

	doc, err := parseDocument(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Signoffs) != 0 {
		t.Fatalf("example signoffs above the heading were counted as approval: %+v", doc.Signoffs)
	}
}

// The 32 real signoffs that sit under later headings such as "Cycle 2 (…)" must keep counting;
// a stricter "inside the Signoffs section only" rule would have dropped them.
func TestSignoffsUnderLaterHeadingsStillCount(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "consensus.md")
	writeF(t, path, strings.Join([]string{
		"## Signoffs", "",
		"### Signoff: alice — 2026-01-01", "Status: accept", "",
		"## Cycle 2 (review/round-02 → complete)", "",
		"### Signoff: bob — 2026-01-02", "Status: accept", "",
	}, "\n"))

	doc, err := parseDocument(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Signoffs) != 2 {
		t.Fatalf("want both signoffs, got %d: %+v", len(doc.Signoffs), doc.Signoffs)
	}
}
