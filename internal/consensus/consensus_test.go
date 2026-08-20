package consensus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/protocol"
)

func TestDraftAndAppendSignoffTriage(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex", "claude"})
	writeRoundFiles(t, root, "sample", false, "round-01", []string{"codex", "claude"})
	now := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)

	summary, err := Draft(root, "sample", DraftOptions{By: "codex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != TriagePartial || strings.Join(summary.Missing, ",") != "codex,claude" {
		t.Fatalf("summary=%+v", summary)
	}
	assertPromptStatus(t, root, "sample", "consensus")

	summary, err = AppendSignoff(root, "sample", SignoffOptions{
		Agent:  "codex",
		Status: "accept",
		Notes:  "Accept.",
		Now:    now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != TriagePartial || strings.Join(summary.Missing, ",") != "claude" {
		t.Fatalf("summary=%+v", summary)
	}

	summary, err = AppendSignoff(root, "sample", SignoffOptions{
		Agent:  "claude",
		Status: "reserve",
		Notes:  "Reservation is logged.",
		Now:    now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != TriageReserved {
		t.Fatalf("triage=%s, want %s", summary.Triage, TriageReserved)
	}

	_, err = AppendSignoff(root, "sample", SignoffOptions{Agent: "codex", Status: "accept", Now: now})
	if err == nil || !strings.Contains(err.Error(), "already signed") {
		t.Fatalf("duplicate err=%v", err)
	}
}

func TestMalformedSignoffs(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	path := filepath.Join(root, protocol.DeckDir, "ideas", "sample", "consensus.md")
	writeFile(t, path, `---
idea: sample
drafted-by: codex
date: 2026-05-12
---

## Signoffs

### Signoff: codex — 2026-05-12
Status: ❌ BLOCK
Notes: Missing counter.

### Signoff: ghost — 2026-05-12
Status: ✅ ACCEPT
`)

	summary, err := Status(root, "sample", false)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != TriageMalformed {
		t.Fatalf("triage=%s errors=%v", summary.Triage, summary.Errors)
	}
	if len(summary.Errors) != 2 {
		t.Fatalf("errors=%v", summary.Errors)
	}
}

func TestManualStatusAliasesAffectTriage(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	path := filepath.Join(root, protocol.DeckDir, "ideas", "sample", "consensus.md")
	writeFile(t, path, `---
idea: sample
drafted-by: codex
date: 2026-05-12
---

## Signoffs

### Signoff: codex — 2026-05-12
Status: block
Notes: Manual block.
Counter-proposal: Run another round.
`)

	summary, err := Status(root, "sample", false)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != TriageBlocked || len(summary.Errors) != 0 {
		t.Fatalf("summary=%+v", summary)
	}
}

func TestFinalizeCreatesFinalAndUpdatesStatus(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	writeRoundFiles(t, root, "sample", false, "round-01", []string{"codex"})
	now := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
	if _, err := Draft(root, "sample", DraftOptions{By: "codex", Now: now}); err != nil {
		t.Fatal(err)
	}
	if _, err := AppendSignoff(root, "sample", SignoffOptions{Agent: "codex", Status: "accept", Notes: "Accept.", Now: now}); err != nil {
		t.Fatal(err)
	}

	// Step 1: the scaffold is written and the idea stays OPEN. Closing an idea around an empty
	// outline is what codex-1/F5 caught.
	finalPath, summary, err := Finalize(root, "sample", FinalizeOptions{By: "codex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if !summary.Scaffolded {
		t.Fatalf("first finalize should report a scaffold, got %+v", summary)
	}
	if summary.Triage != TriageReady {
		t.Fatalf("summary=%+v", summary)
	}
	finalData, err := os.ReadFile(finalPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range protocol.RequiredFinalSections {
		if !strings.Contains(string(finalData), want) {
			t.Fatalf("scaffold missing the protocol section %q:\n%s", want, finalData)
		}
	}
	assertPromptStatus(t, root, "sample", "consensus")

	// Step 1': re-running while it is still a scaffold must refuse and say why.
	if _, _, err := Finalize(root, "sample", FinalizeOptions{By: "codex", Now: now}); err == nil {
		t.Fatal("finalize closed the idea around an unwritten scaffold")
	} else if !strings.Contains(err.Error(), "scaffold") {
		t.Fatalf("refusal does not name the cause: %v", err)
	}
	assertPromptStatus(t, root, "sample", "consensus")

	// Step 2: once written, finalize closes the idea.
	writeFile(t, finalPath, writtenFinal("sample"))
	if _, summary, err = Finalize(root, "sample", FinalizeOptions{By: "codex", Now: now}); err != nil {
		t.Fatal(err)
	}
	if summary.Scaffolded {
		t.Fatal("a written FINAL.md must close the idea, not report a scaffold")
	}
	assertPromptStatus(t, root, "sample", "final")
}

// writtenFinal is a FINAL.md with every protocol section filled in.
func writtenFinal(slug string) string {
	var b strings.Builder
	b.WriteString("---\nidea: " + slug + "\nstatus: final\nauthor: codex\n---\n\n")
	for _, section := range protocol.RequiredFinalSections {
		b.WriteString(section + "\n\nReal content for this section.\nA second line of detail.\nA third line of detail.\n\n")
	}
	return b.String()
}

func TestReservedFinalizeRequiresOpenItems(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	writeRoundFiles(t, root, "sample", false, "round-01", []string{"codex"})
	now := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
	if _, err := Draft(root, "sample", DraftOptions{By: "codex", Now: now}); err != nil {
		t.Fatal(err)
	}
	if _, err := AppendSignoff(root, "sample", SignoffOptions{Agent: "codex", Status: "reserve", Notes: "Needs follow-up.", Now: now}); err != nil {
		t.Fatal(err)
	}

	if _, _, err := Finalize(root, "sample", FinalizeOptions{By: "codex", Now: now}); err == nil {
		t.Fatal("expected reserved consensus without open items to fail")
	}
}

func TestReservedFinalizeSucceedsWithOpenItems(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	writeRoundFiles(t, root, "sample", false, "round-01", []string{"codex"})
	now := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
	if _, err := Draft(root, "sample", DraftOptions{By: "codex", Now: now}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, protocol.DeckDir, "ideas", "sample", "consensus.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(data), "## Open items deferred to implementation\n\n## Signoffs", "## Open items deferred to implementation\n\n- codex: carry this reservation into implementation.\n\n## Signoffs", 1)
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := AppendSignoff(root, "sample", SignoffOptions{Agent: "codex", Status: "reserve", Notes: "Needs follow-up.", Now: now}); err != nil {
		t.Fatal(err)
	}

	// Two steps since codex-1/F5: scaffold first, close once it is written.
	finalPath, summary, err := Finalize(root, "sample", FinalizeOptions{By: "codex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if !summary.Scaffolded {
		t.Fatalf("first finalize should report a scaffold, got %+v", summary)
	}
	assertPromptStatus(t, root, "sample", "consensus")

	writeFile(t, finalPath, writtenFinal("sample"))
	if _, _, err := Finalize(root, "sample", FinalizeOptions{By: "codex", Now: now}); err != nil {
		t.Fatal(err)
	}
	assertPromptStatus(t, root, "sample", "final")
}

func TestReopenBlockedConsensus(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	writeRoundFiles(t, root, "sample", false, "round-01", []string{"codex"})
	path := filepath.Join(root, protocol.DeckDir, "ideas", "sample", "consensus.md")
	writeFile(t, path, `---
idea: sample
drafted-by: codex
date: 2026-05-12
---

## Signoffs

### Signoff: codex — 2026-05-12
Status: ❌ BLOCK
Notes: Missing implementation guard.
Counter-proposal (required if ❌): Add the guard and run another round.
`)

	aborted, err := Reopen(root, "sample", ReopenOptions{Reason: "Blocked by codex."})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("consensus.md should be moved, err=%v", err)
	}
	data, err := os.ReadFile(aborted)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Blocked by codex.") {
		t.Fatalf("aborted file missing reason:\n%s", data)
	}
	if !strings.HasSuffix(aborted, "round-01-consensus-aborted-01.md") {
		t.Fatalf("aborted=%q", aborted)
	}
	assertPromptStatus(t, root, "sample", "round-01")
}

func TestReviewDraftUsesReviewPath(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	writeRoundFiles(t, root, "sample", true, "round-01", []string{"codex"})
	now := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)

	summary, err := Draft(root, "sample", DraftOptions{Review: true, By: "codex", ReviewedCommit: "abc123", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if !summary.Review || !strings.HasSuffix(summary.Path, filepath.Join("review", "consensus.md")) {
		t.Fatalf("summary=%+v", summary)
	}
	data, err := os.ReadFile(summary.Path)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"cycle: 1", "reviewed-commit: abc123"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("review consensus missing %q:\n%s", want, data)
		}
	}
	assertPromptStatus(t, root, "sample", "round-01")
}

func TestDraftSelectsLatestRoundNumerically(t *testing.T) {
	root := setupIdea(t, "sample", []string{"codex"})
	writeRoundFiles(t, root, "sample", false, "round-2", []string{"codex"})
	writeRoundFiles(t, root, "sample", false, "round-10", []string{"codex"})

	summary, err := Draft(root, "sample", DraftOptions{By: "codex", Now: time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(summary.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Review round-10") {
		t.Fatalf("draft did not use numeric latest round:\n%s", data)
	}
}

func setupIdea(t *testing.T, slug string, participants []string) string {
	t.Helper()
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(ideaDir, "00-prompt.md"), `---
idea: `+slug+`
author: codex
created: 2026-05-12
participants: [`+strings.Join(participants, ", ")+`]
status: round-01
---
`)
	return root
}

func writeRoundFiles(t *testing.T, root, slug string, review bool, round string, participants []string) {
	t.Helper()
	base := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	if review {
		base = filepath.Join(base, "review")
	}
	roundDir := filepath.Join(base, round)
	if err := os.MkdirAll(roundDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, participant := range participants {
		writeFile(t, filepath.Join(roundDir, participant+".md"), "# "+participant+"\n")
	}
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertPromptStatus(t *testing.T, root, slug, want string) {
	t.Helper()
	meta, err := protocol.ReadFrontmatter(filepath.Join(root, protocol.DeckDir, "ideas", slug, "00-prompt.md"))
	if err != nil {
		t.Fatal(err)
	}
	if got := meta["status"]; got != want {
		t.Fatalf("status=%q, want %q", got, want)
	}
}
