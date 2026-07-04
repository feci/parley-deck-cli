package driver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRoundFile(t *testing.T, dir, agent, body string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, agent+".md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBuildRoundDigestWellFormed(t *testing.T) {
	idea := t.TempDir()
	rd := filepath.Join(idea, "round-01")
	writeRoundFile(t, rd, "claude-1", "---\nagent: claude-1\n---\n\n## Summary\n\nI accept the v2 framing and agree with the split. More detail follows.\n\n## Risks\n")
	writeRoundFile(t, rd, "codex-1", "---\nagent: codex-1\n---\n\n## Summary\n\nThis has a blocker: the lock races. We need a counter-proposal here.\n")

	d := BuildRoundDigest(idea, "demo", 1, []string{"claude-1", "codex-1"}, "opening round-02")
	if d.Total != 2 || d.Completed != 2 {
		t.Fatalf("completeness = %d/%d", d.Completed, d.Total)
	}
	if d.Lines[0].Position == "" || d.Lines[0].Fell {
		t.Fatalf("claude line = %+v", d.Lines[0])
	}
	if !strings.Contains(d.Lines[0].Position, "accept") {
		t.Fatalf("position not extracted: %q", d.Lines[0].Position)
	}
	if d.FlagBlock < 1 || d.FlagCounter < 1 || d.FlagAccept < 1 {
		t.Fatalf("flags = block:%d counter:%d accept:%d", d.FlagBlock, d.FlagCounter, d.FlagAccept)
	}
	if d.Next != "opening round-02" {
		t.Fatalf("next = %q", d.Next)
	}
}

func TestBuildRoundDigestMissingSummaryFallback(t *testing.T) {
	idea := t.TempDir()
	rd := filepath.Join(idea, "round-01")
	writeRoundFile(t, rd, "hermes-1", "---\nagent: hermes-1\n---\n\nNo summary heading here, just a leading paragraph of prose.\n")
	d := BuildRoundDigest(idea, "demo", 1, []string{"hermes-1"}, "")
	if !d.Lines[0].Fell {
		t.Fatal("expected fell-back extraction")
	}
	if !strings.HasPrefix(d.Lines[0].Position, "No summary heading") {
		t.Fatalf("fallback position = %q", d.Lines[0].Position)
	}
}

func TestBuildRoundDigestMissingAgent(t *testing.T) {
	idea := t.TempDir()
	rd := filepath.Join(idea, "round-01")
	writeRoundFile(t, rd, "claude-1", "## Summary\n\nPresent.\n")
	// codex-1 file absent → not present, does not error, counts toward Total not Completed.
	d := BuildRoundDigest(idea, "demo", 1, []string{"claude-1", "codex-1"}, "")
	if d.Total != 2 || d.Completed != 1 {
		t.Fatalf("completeness = %d/%d", d.Completed, d.Total)
	}
	if d.Lines[1].Present {
		t.Fatal("codex-1 should be not-present")
	}
}

func TestExtractPositionCapsLongLine(t *testing.T) {
	long := strings.Repeat("word ", 60) // ~300 chars, no sentence boundary
	pos, _ := extractPosition("## Summary\n\n" + long + "\n")
	if len([]rune(pos)) > digestPositionCap+1 {
		t.Fatalf("position not capped: %d chars", len([]rune(pos)))
	}
	if !strings.HasSuffix(pos, "…") {
		t.Fatalf("capped position should end with ellipsis: %q", pos)
	}
}

func TestStanceFlagsNoDoubleCountBlocker(t *testing.T) {
	// "blocker" is a superstring of "block": must count once, not twice.
	b, _, _, _ := stanceFlags("this is a blocker")
	if b != 1 {
		t.Fatalf("blocker should count as 1 block mention, got %d", b)
	}
}
