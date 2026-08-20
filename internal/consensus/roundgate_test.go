package consensus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

// codex-1/F4: --by was written into FINAL.md unchecked, so an idea could be closed in the name of
// an agent that never took part.
func TestFinalizeRejectsANonParticipantDrafter(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	writeRoundFiles(t, root, "sample", false, "round-01", []string{"codex"})
	now := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
	if _, err := Draft(root, "sample", DraftOptions{By: "codex", Now: now}); err != nil {
		t.Fatal(err)
	}
	if _, err := AppendSignoff(root, "sample", SignoffOptions{Agent: "codex", Status: "accept", Notes: "ok", Now: now}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Finalize(root, "sample", FinalizeOptions{By: "stranger", Now: now}); err == nil {
		t.Fatal("a non-participant was allowed to close the idea")
	} else if !strings.Contains(err.Error(), "not a participant") {
		t.Fatalf("refusal does not name the cause: %v", err)
	}
}

// codex-1/F9: any filler satisfied the deferred-items section, so a reservation could be carried
// past finalize without being written down.
func TestReservationMustBeNamedInDeferredItems(t *testing.T) {
	doc := document{
		Raw: "## Open items deferred to implementation\n\n- something generic\n\n## Signoffs\n",
		Signoffs: []Signoff{
			{Agent: "codex", Status: "🟡 RESERVE"},
			{Agent: "kimi", Status: "✅ ACCEPT"},
		},
	}
	got := unloggedReservations(doc)
	if len(got) != 1 || got[0] != "codex" {
		t.Fatalf("want codex reported as unlogged, got %v", got)
	}

	doc.Raw = "## Open items deferred to implementation\n\n- codex: carry the sandbox question\n\n## Signoffs\n"
	if got := unloggedReservations(doc); len(got) != 0 {
		t.Fatalf("a named reservation must count as logged, got %v", got)
	}
}

// codex-1/F21: a consensus.md whose frontmatter named a different idea was read as this idea's
// consensus, so a document copied between ideas carried its signoffs with it.
func TestConsensusDeclaringAnotherIdeaIsReported(t *testing.T) {
	doc := document{
		Path: "consensus.md",
		Raw:  "---\nidea: some-other-idea\n---\n\n## Signoffs\n",
	}
	summary := validateDocument("this-idea", []string{"a"}, false, doc)
	var found bool
	for _, e := range summary.Errors {
		if strings.Contains(e, "some-other-idea") && strings.Contains(e, "this-idea") {
			found = true
		}
	}
	if !found {
		t.Fatalf("slug mismatch not reported: %v", summary.Errors)
	}

	// The matching case must not complain.
	ok := document{Path: "consensus.md", Raw: "---\nidea: this-idea\n---\n\n## Signoffs\n"}
	for _, e := range validateDocument("this-idea", []string{"a"}, false, ok).Errors {
		if strings.Contains(e, "frontmatter idea") {
			t.Fatalf("matching slug reported as an error: %v", e)
		}
	}
}

// codex-1/F3: a review consensus demanded the implementer's signoff, which §6 forbids it from
// giving, so a standard-track review consensus stayed `partial` however many reviewers accepted.
func TestReviewConsensusDoesNotAwaitTheImplementersSignoff(t *testing.T) {
	root := setupIdea(t, "sample", []string{"impl", "rev-a"})
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "sample")
	writeF(t, filepath.Join(ideaDir, "IMPLEMENTATION.md"), "---\nidea: sample\nimplementer: impl\n---\n\n## Summary of work\nok\n")
	writeRoundFiles(t, root, "sample", true, "round-01", []string{"rev-a"})

	now := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
	if _, err := Draft(root, "sample", DraftOptions{By: "rev-a", Review: true, Now: now}); err != nil {
		t.Fatalf("review draft rejected: %v", err)
	}
	if _, err := AppendSignoff(root, "sample", SignoffOptions{Agent: "rev-a", Status: "accept", Notes: "ok", Review: true, Now: now}); err != nil {
		t.Fatal(err)
	}

	summary, err := Status(root, "sample", true)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range summary.Missing {
		if strings.Contains(m, "impl") {
			t.Fatalf("review consensus waits on the implementer: missing=%v", summary.Missing)
		}
	}
	if summary.Triage != TriageReady {
		t.Fatalf("triage=%s, want ready once every reviewer accepted (missing=%v errors=%v)", summary.Triage, summary.Missing, summary.Errors)
	}
}

// codex-1/F10: the review-consensus generator wrote `cycle:` where the protocol requires
// `review-cycle:`, and emitted neither `outstanding_agreed_fixes` nor `blocked` — which the
// auto-driver requires. The manual command produced an artifact its own automation rejects.
func TestReviewConsensusTemplateMatchesTheDriverSchema(t *testing.T) {
	root := setupIdea(t, "sample", []string{"impl", "rev-a"})
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "sample")
	writeF(t, filepath.Join(ideaDir, "IMPLEMENTATION.md"), "---\nidea: sample\nimplementer: impl\n---\n\n## Summary of work\nok\n")
	writeRoundFiles(t, root, "sample", true, "round-01", []string{"rev-a"})

	summary, err := Draft(root, "sample", DraftOptions{By: "rev-a", Review: true, ReviewedCommit: "abc1234", Now: time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(summary.Path)
	if err != nil {
		t.Fatal(err)
	}
	body := string(data)

	for _, want := range []string{"review-cycle:", "outstanding_agreed_fixes:", "blocked:"} {
		if !strings.Contains(body, want) {
			t.Errorf("review consensus template is missing %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "\ncycle:") {
		t.Errorf("template still writes the non-schema `cycle:` key:\n%s", body)
	}
	// The count must be a placeholder the drafter replaces, never a silent 0 — asserting
	// "nothing outstanding" on the drafter's behalf is the claim the gate exists to check.
	if strings.Contains(body, "outstanding_agreed_fixes: 0") {
		t.Error("template pre-answers outstanding_agreed_fixes as 0")
	}
}

// Review round 1, @kimi-1 MAJOR: excluding the implementer from the participant list made its OWN
// signoff an "unknown participant", and malformed outranks every other triage — nine review
// consensuses in this deck flipped to malformed, two of them in flight.
//
// Who is AWAITED and who may SIGN are different questions.
func TestTheImplementerMaySignAReviewConsensusEvenThoughItIsNotAwaited(t *testing.T) {
	root := setupIdea(t, "sample", []string{"impl", "rev-a"})
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "sample")
	writeF(t, filepath.Join(ideaDir, "IMPLEMENTATION.md"), "---\nidea: sample\nimplementer: impl\n---\n\n## Summary of work\nok\n")
	writeRoundFiles(t, root, "sample", true, "round-01", []string{"rev-a"})

	now := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
	if _, err := Draft(root, "sample", DraftOptions{By: "rev-a", Review: true, Now: now}); err != nil {
		t.Fatal(err)
	}
	for _, agent := range []string{"rev-a", "impl"} {
		if _, err := AppendSignoff(root, "sample", SignoffOptions{Agent: agent, Status: "accept", Notes: "ok", Review: true, Now: now}); err != nil {
			t.Fatalf("%s could not sign: %v", agent, err)
		}
	}

	summary, err := Status(root, "sample", true)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range summary.Errors {
		if strings.Contains(e, "unknown participant") {
			t.Fatalf("the implementer's own signoff was rejected: %v", summary.Errors)
		}
	}
	if summary.Triage != TriageReady {
		t.Fatalf("triage=%s, want ready; errors=%v missing=%v", summary.Triage, summary.Errors, summary.Missing)
	}
	// …and it is still not AWAITED: absent, the consensus is still ready.
	for _, m := range summary.Missing {
		if strings.Contains(m, "impl") {
			t.Fatalf("implementer is awaited again: %v", summary.Missing)
		}
	}
}
