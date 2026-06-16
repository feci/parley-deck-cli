package retro

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// writeIdea creates parley-deck/ideas/<slug> with the given files (relative path
// → content) and the given number of round-NN and review/round-NN dirs.
func writeIdea(t *testing.T, root, slug string, rounds, reviewRounds int, files map[string]string) {
	t.Helper()
	dir := filepath.Join(root, "parley-deck", "ideas", slug)
	for i := 1; i <= rounds; i++ {
		mustMkdir(t, filepath.Join(dir, fmt.Sprintf("round-%02d", i)))
	}
	for i := 1; i <= reviewRounds; i++ {
		mustMkdir(t, filepath.Join(dir, "review", fmt.Sprintf("round-%02d", i)))
	}
	for rel, content := range files {
		p := filepath.Join(dir, rel)
		mustMkdir(t, filepath.Dir(p))
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func mustMkdir(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestScanScoreAndClassify(t *testing.T) {
	root := t.TempDir()
	// fix-up-heavy: 2 design rounds, 2 review rounds, 2 fix-up cycles, a NOT-FIXED.
	writeIdea(t, root, "fixup-heavy", 2, 2, map[string]string{
		"IMPLEMENTATION.md":        "## Fix-up cycle 1\nx\n## Fix-up cycle 2\ny\n",
		"review/round-02/codex.md": "Verdict: NOT-FIXED something\n",
	})
	// blocked: a ❌ signoff.
	writeIdea(t, root, "blocked-one", 1, 1, map[string]string{
		"consensus.md": "### Signoff: codex\nStatus: ❌ BLOCKER\n",
	})
	// escalated: inbox to-user note referencing the slug.
	writeIdea(t, root, "escalated-one", 1, 0, nil)
	mustMkdir(t, filepath.Join(root, "parley-deck", "inbox"))
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "inbox", "claude-to-user_escalated-one_q.md"), []byte("blocking: yes\nidea: escalated-one\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// clean: minimal, low friction.
	writeIdea(t, root, "clean-one", 1, 1, nil)

	signals, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(signals) != 4 {
		t.Fatalf("want 4 ideas, got %d", len(signals))
	}
	// Results are score-sorted descending.
	for i := 1; i < len(signals); i++ {
		if signals[i-1].Score < signals[i].Score {
			t.Fatalf("not score-sorted: %+v", signals)
		}
	}
	by := map[string]IdeaSignals{}
	for _, s := range signals {
		by[s.Slug] = s
	}
	if got := by["fixup-heavy"]; got.FailureType != "fix-up-heavy" || got.FixupCycles != 2 || got.NotFixed != 1 {
		t.Fatalf("fixup-heavy classified wrong: %+v", got)
	}
	if got := by["blocked-one"]; got.FailureType != "blocked-or-abandoned" || !got.Blocked {
		t.Fatalf("blocked-one classified wrong: %+v", got)
	}
	if got := by["escalated-one"]; got.FailureType != "escalation" || got.Escalations != 1 {
		t.Fatalf("escalated-one classified wrong: %+v", got)
	}
	if got := by["clean-one"]; got.FailureType != "low-friction" || got.Score != 0 {
		t.Fatalf("clean-one should be low-friction score 0: %+v", got)
	}
}

func TestSelectIsTypeDiverseAndExcludesLowFriction(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "fixup-a", 2, 2, map[string]string{"IMPLEMENTATION.md": "## Fix-up cycle 1\n## Fix-up cycle 2\n"})
	writeIdea(t, root, "fixup-b", 2, 2, map[string]string{"IMPLEMENTATION.md": "## Fix-up cycle 1\n## Fix-up cycle 2\n"})
	writeIdea(t, root, "blocked-a", 1, 1, map[string]string{"consensus.md": "Status: ❌ BLOCKER\n"})
	writeIdea(t, root, "clean-a", 1, 1, nil)

	signals, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	coreset := Select(signals, 10)
	for _, s := range coreset {
		if s.FailureType == "low-friction" || s.Score <= 0 {
			t.Fatalf("coreset must exclude low-friction/zero-score: %+v", s)
		}
	}
	// Type diversity: blocked-a must be picked before the second fix-up idea even
	// though a fix-up idea may outscore it, because pass 1 covers each type once.
	pickedTypes := map[string]int{}
	for _, s := range coreset {
		pickedTypes[s.FailureType]++
	}
	if pickedTypes["blocked-or-abandoned"] == 0 {
		t.Fatalf("coreset missing the 'blocked-or-abandoned' type: %+v", coreset)
	}
	// k cap.
	if got := Select(signals, 1); len(got) != 1 {
		t.Fatalf("k=1 must yield 1, got %d", len(got))
	}
}

func TestD4SignalsAndDesignChurn(t *testing.T) {
	root := t.TempDir()
	// abandoned work (status frontmatter) → blocked-or-abandoned.
	writeIdea(t, root, "abandoned-one", 1, 0, map[string]string{
		"IMPLEMENTATION.md": "---\nidea: abandoned-one\nstatus: abandoned\n---\n",
	})
	// 2 design rounds only → must be design-churn (not low-friction), so it is
	// NOT dropped from the coreset (review fix 2).
	writeIdea(t, root, "design-two", 2, 0, nil)
	// run failures: a structured run event log with agent.failed attributed to
	// the idea via run.created (review fix 4).
	writeIdea(t, root, "runtime-fail", 1, 0, nil)
	mustMkdir(t, filepath.Join(root, "parley-deck", "runs", "r1"))
	events := `{"type":"run.created","data":{"idea":"runtime-fail"}}` + "\n" +
		`{"type":"agent.failed","data":{"agent":"codex"}}` + "\n" +
		`{"type":"agent.stalled","data":{"agent":"agy"}}` + "\n"
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "runs", "r1", "events.jsonl"), []byte(events), 0o644); err != nil {
		t.Fatal(err)
	}

	signals, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	by := map[string]IdeaSignals{}
	for _, s := range signals {
		by[s.Slug] = s
	}
	if got := by["abandoned-one"]; !got.Abandoned || got.FailureType != "blocked-or-abandoned" {
		t.Fatalf("abandoned-one wrong: %+v", got)
	}
	if got := by["design-two"]; got.FailureType != "design-churn" || got.Score <= 0 {
		t.Fatalf("design-two must be design-churn with score>0 (review fix 2): %+v", got)
	}
	if got := by["runtime-fail"]; got.RunFailures != 2 || got.FailureType != "runtime-failure" {
		t.Fatalf("runtime-fail must count 2 run failures and bucket runtime-failure: %+v", got)
	}
	// design-two must survive coreset selection (not excluded as low-friction).
	coreset := Select(signals, 10)
	found := false
	for _, s := range coreset {
		if s.Slug == "design-two" {
			found = true
		}
	}
	if !found {
		t.Fatalf("design-two (2 rounds) must be in the coreset, not dropped: %+v", coreset)
	}
}

// TestBlockerDetectionIgnoresProse is the regression test for the 1.28.0
// false-positive: a review consensus that QUOTES "Verdict: BLOCK" in prose must
// not flag the idea as blocked, while real BLOCK verdicts / ❌ signoffs must.
func TestBlockerDetectionIgnoresProse(t *testing.T) {
	root := t.TempDir()
	// prose-block: review consensus discusses "Verdict: BLOCK"; all real verdicts
	// are ACCEPT. Must NOT be blocked (the bug).
	writeIdea(t, root, "prose-block", 1, 1, map[string]string{
		"review/consensus.md":      "Fix: also match `Verdict: BLOCK` in reviewer files.\n### Signoff: codex\nStatus: ✅ ACCEPT\n",
		"review/round-01/codex.md": "## Verdict\nACCEPT\n",
	})
	// real-block: a genuine BLOCK verdict in a round file. MUST be blocked.
	writeIdea(t, root, "real-block", 1, 1, map[string]string{
		"review/round-01/codex.md": "## Findings\n### [CRITICAL] x\n## Verdict\nBLOCK\n",
	})
	// request-changes: a fixes-requested verdict using the word "block" as a verb.
	// MUST NOT be a hard BLOCK (it is captured by fix-up signals instead).
	writeIdea(t, root, "request-changes", 1, 1, map[string]string{
		"review/round-01/codex.md": "## Verdict\nREQUEST-CHANGES: block on the missing test, then I expect to ACCEPT.\n",
	})
	// sign-block: a ❌ BLOCKER signoff in consensus.md. MUST be blocked.
	writeIdea(t, root, "sign-block", 1, 1, map[string]string{
		"consensus.md": "### Signoff: agy\nStatus: ❌ BLOCKER\nNotes: no.\n",
	})

	signals, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	by := map[string]IdeaSignals{}
	for _, s := range signals {
		by[s.Slug] = s
	}
	if by["prose-block"].Blocked {
		t.Fatalf("prose-block must NOT be blocked (prose mention of Verdict: BLOCK): %+v", by["prose-block"])
	}
	if !by["real-block"].Blocked {
		t.Fatalf("real-block MUST be blocked (## Verdict\\nBLOCK): %+v", by["real-block"])
	}
	if by["request-changes"].Blocked {
		t.Fatalf("request-changes must NOT be a hard BLOCK: %+v", by["request-changes"])
	}
	if !by["sign-block"].Blocked {
		t.Fatalf("sign-block MUST be blocked (Status: ❌ signoff): %+v", by["sign-block"])
	}
}

func TestDiagnoseEmptyAndGrouped(t *testing.T) {
	if got := Diagnose(nil); got == "" || !contains(got, "No hard cases") {
		t.Fatalf("empty diagnose should say no hard cases: %q", got)
	}
	cs := []IdeaSignals{
		{Slug: "a", FailureType: "fix-up-heavy", Reasons: []string{"fix-up churn"}},
		{Slug: "b", FailureType: "blocked", Reasons: []string{"blocker signoff"}},
	}
	out := Diagnose(cs)
	for _, want := range []string{"Failure mode: fix-up-heavy", "Failure mode: blocked", "hypotheses, not findings"} {
		if !contains(out, want) {
			t.Fatalf("diagnose missing %q:\n%s", want, out)
		}
	}
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (indexOf(s, sub) >= 0) }
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
