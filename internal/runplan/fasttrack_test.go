package runplan

import (
	"os"
	"path/filepath"
	"testing"
)

func ideaWith(t *testing.T, frontmatter string) string {
	t.Helper()
	dir := t.TempDir()
	body := "---\nidea: demo\n" + frontmatter + "---\n\n## Prompt\n\nx\n"
	if err := os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// codex-1/F7: the planner consulted only `cross_review_rounds:` (defaulting it to 1) and never
// looked at `track:`, so it told a fast-track idea to open the cross-review round that §4.0's
// binding per-track table explicitly skips.
func TestFastTrackIsNotToldToOpenACrossReviewRound(t *testing.T) {
	if got := nextCrossReviewRound(ideaWith(t, "track: fast\n"), "round-01"); got != "" {
		t.Fatalf("fast track was told to open %q; §4.0 skips cross-review on fast", got)
	}
}

// standard and deliberation keep their rounds.
func TestNonFastTracksStillOpenTheirCrossReviewRound(t *testing.T) {
	for _, tr := range []string{"standard", "deliberation"} {
		if got := nextCrossReviewRound(ideaWith(t, "track: "+tr+"\n"), "round-01"); got != "round-02" {
			t.Errorf("track %s: got %q, want round-02", tr, got)
		}
	}
	// An idea declaring no track is unchanged.
	if got := nextCrossReviewRound(ideaWith(t, ""), "round-01"); got != "round-02" {
		t.Errorf("absent track: got %q, want round-02", got)
	}
}

// An explicit cross_review_rounds: 0 still wins for non-fast tracks.
func TestExplicitZeroRoundsStillHonoured(t *testing.T) {
	if got := nextCrossReviewRound(ideaWith(t, "track: standard\ncross_review_rounds: 0\n"), "round-01"); got != "" {
		t.Fatalf("got %q, want none", got)
	}
}
