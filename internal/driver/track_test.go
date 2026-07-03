package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func writeTrackPrompt(t *testing.T, dir, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestNewAppliesFastTrack(t *testing.T) {
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\ntrack: fast\n---\n")
	d := New(Config{IdeaDir: dir, Participants: []string{"a", "b", "c"}}, nil)
	if d.cfg.Track != "fast" {
		t.Errorf("Track = %q, want fast", d.cfg.Track)
	}
	if d.cfg.MaxReviewers != 1 || d.cfg.MinReviewers != 1 {
		t.Errorf("fast reviewers: Max=%d Min=%d, want 1/1", d.cfg.MaxReviewers, d.cfg.MinReviewers)
	}
	if d.cfg.CrossReviewRounds != 0 || d.cfg.MaxFixupCycles != 1 {
		t.Errorf("fast loop: CrossReview=%d Fixup=%d, want 0/1", d.cfg.CrossReviewRounds, d.cfg.MaxFixupCycles)
	}
	if d.trackErr != nil {
		t.Errorf("unexpected trackErr: %v", d.trackErr)
	}
}

func TestNewExplicitStandardAppliesCaps(t *testing.T) {
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\ntrack: standard\n---\n")
	d := New(Config{IdeaDir: dir, Participants: []string{"a", "b", "c", "d"}}, nil)
	if d.cfg.MaxReviewers != 2 || d.cfg.MinReviewers != 2 || d.cfg.MaxFixupCycles != 2 {
		t.Errorf("explicit standard: Max=%d Min=%d Fixup=%d, want 2/2/2", d.cfg.MaxReviewers, d.cfg.MinReviewers, d.cfg.MaxFixupCycles)
	}
}

func TestNewDeliberationIsLegacy(t *testing.T) {
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\ntrack: deliberation\n---\n")
	// CrossReviewRounds:-1 exercises New's default (→1); deliberation must PRESERVE it.
	d := New(Config{IdeaDir: dir, Participants: []string{"a", "b"}, CrossReviewRounds: -1}, nil)
	// deliberation == today's full lifecycle: no reviewer cap, LE-11 min 2, fixup 3, cross-review 1.
	if d.cfg.MaxReviewers != 0 || d.cfg.MinReviewers != 2 || d.cfg.MaxFixupCycles != 3 || d.cfg.CrossReviewRounds != 1 {
		t.Errorf("deliberation not legacy: Max=%d Min=%d Fixup=%d Cross=%d", d.cfg.MaxReviewers, d.cfg.MinReviewers, d.cfg.MaxFixupCycles, d.cfg.CrossReviewRounds)
	}
}

func TestNewAbsentTrackIsLegacy(t *testing.T) {
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\n---\n")
	d := New(Config{IdeaDir: dir, Participants: []string{"a", "b"}}, nil)
	if d.cfg.MaxReviewers != 0 || d.cfg.MinReviewers != 2 || d.cfg.MaxFixupCycles != 3 {
		t.Errorf("absent track not legacy: Max=%d Min=%d Fixup=%d", d.cfg.MaxReviewers, d.cfg.MinReviewers, d.cfg.MaxFixupCycles)
	}
}

func TestFastContradictionEscalates(t *testing.T) {
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\ntrack: fast\nauto_implement: true\n---\n")
	d := New(Config{IdeaDir: dir, Participants: []string{"a", "b"}, AutoImplement: true, Auto: true}, nil)
	if d.trackErr == nil {
		t.Fatal("expected trackErr for fast + auto_implement contradiction")
	}
	act, _, err := d.Advance(context.Background())
	if act != ActionEscalated || err == nil {
		t.Errorf("Advance = (%v, %v); want ActionEscalated + error", act, err)
	}
}

func TestFastNonSoloEscalates(t *testing.T) {
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\ntrack: fast\n---\n")
	// Only one participant → 0 available reviewers → non-solo violation.
	d := New(Config{IdeaDir: dir, Participants: []string{"solo"}, Auto: true}, nil)
	if d.trackErr == nil {
		t.Fatal("expected trackErr for fast with no independent reviewer (non-solo)")
	}
}

func TestExplicitDeliberationNonSoloEscalates(t *testing.T) {
	// review-01 F1: explicit deliberation is subject to the non-solo floor too.
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\ntrack: deliberation\n---\n")
	d := New(Config{IdeaDir: dir, Participants: []string{"solo"}, Auto: true}, nil)
	if d.trackErr == nil {
		t.Fatal("explicit deliberation with a single participant must set trackErr (non-solo)")
	}
}

func TestAbsentTrackSoloDoesNotTrackError(t *testing.T) {
	// Absent track keeps today's driver behaviour; the non-solo floor is enforced
	// at the app/preflight layer, not this driver gate.
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\n---\n")
	d := New(Config{IdeaDir: dir, Participants: []string{"solo"}, Auto: true}, nil)
	if d.trackErr != nil {
		t.Errorf("absent track must not set trackErr for a solo roster: %v", d.trackErr)
	}
}

func TestExplicitStandardCapsCrossReview(t *testing.T) {
	// review-01 F5: explicit standard caps cross-review rounds at 2.
	dir := t.TempDir()
	writeTrackPrompt(t, dir, "---\nidea: x\ntrack: standard\n---\n")
	d := New(Config{IdeaDir: dir, Participants: []string{"a", "b", "c"}, CrossReviewRounds: 5}, nil)
	if d.cfg.CrossReviewRounds != 2 {
		t.Errorf("standard cross-review not capped to 2, got %d", d.cfg.CrossReviewRounds)
	}
}
