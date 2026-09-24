package driver

// lean-organizer D.1: the driver writes a schema-valid per-phase handoff record at
// each phase transition, on the same write discipline as the shipped handoff
// machinery, with the schema doc stating recomputation is authoritative.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
)

func TestCommitCursorWritesPhaseHandoffRecord(t *testing.T) {
	ideaDir := t.TempDir()
	runDir := t.TempDir()
	events := store.New(runDir)
	d := &Driver{cfg: Config{
		IdeaDir: ideaDir, IdeaSlug: "handoff-idea", Participants: []string{"a-1"},
		RunDir: runDir, Root: filepath.Dir(filepath.Dir(ideaDir)), Events: events,
	}}
	c := Cursor{Phase: PhaseConsensus, CurrentRound: 1, IdeaStatus: "consensus"}
	if err := d.commitCursor(c, ActionConsensusDrafted, PhaseRound); err != nil {
		t.Fatalf("commitCursor: %v", err)
	}
	path := filepath.Join(runDir, PhaseHandoffPath(PhaseConsensus))
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("phase handoff record not written: %v", err)
	}
	rec, err := LoadPhaseHandoffRecord(path)
	if err != nil {
		t.Fatalf("record is not schema-valid: %v", err)
	}
	runID := filepath.Base(runDir)
	if rec.Idea != "handoff-idea" || rec.Phase != "consensus" || rec.RunID != runID {
		t.Errorf("record fields incomplete (RunID must be the run dir's base, %q): %+v", runID, rec)
	}
	body, _ := os.ReadFile(path)
	if !strings.Contains(string(body), "recomputed view is authoritative") {
		t.Errorf("schema doc must state recomputation is authoritative")
	}
	if !strings.Contains(string(body), "PhaseDigest") {
		t.Errorf("record must embed the PhaseDigest computed at the transition")
	}
	// The transition is also evented.
	evs, _ := events.Load()
	saw := false
	for _, e := range evs {
		if e.Type == "driver.phase_handoff" {
			saw = true
		}
	}
	if !saw {
		t.Errorf("phase handoff write must be evented")
	}
}

// TestPhaseHandoffRoundTrip (fix-up F13, claude-1 MIN-5): ALL nine frontmatter
// fields survive the store→load round-trip, and the builder takes the run dir as a
// parameter (RunID derived from it, not from the idea path).
func TestPhaseHandoffRoundTrip(t *testing.T) {
	written := time.Date(2026, 9, 24, 1, 2, 3, 4, time.UTC)
	rec := PhaseHandoffRecord{
		RunID:         "20260924T000000.000000000Z",
		Idea:          "handoff-idea",
		Phase:         "impl",
		PreviousPhase: "review",
		Action:        "implementation-published",
		RoundLabel:    "round-02",
		IdeaStatus:    "fix-up-cycle-1",
		WrittenAt:     written,
		Next:          NextAwaitReviewArtifact,
	}
	dir := t.TempDir()
	path := filepath.Join(dir, PhaseHandoffPath(PhaseImpl))
	if err := storeHandoffRecord(path, rec); err != nil {
		t.Fatal(err)
	}
	back, err := LoadPhaseHandoffRecord(path)
	if err != nil {
		t.Fatal(err)
	}
	if back.RunID != rec.RunID || back.Idea != rec.Idea || back.Phase != rec.Phase ||
		back.PreviousPhase != rec.PreviousPhase || back.Action != rec.Action ||
		back.RoundLabel != rec.RoundLabel || back.IdeaStatus != rec.IdeaStatus ||
		back.Next != rec.Next || !back.WrittenAt.Equal(rec.WrittenAt) {
		t.Errorf("nine-field round-trip mismatch:\nback  %+v\nwant  %+v", back, rec)
	}
}

// TestBuildPhaseHandoffRecordTakesRunDir (fix-up F13): the exported builder derives
// RunID from the run-dir parameter — no "parley-deck" placeholder exists anymore.
func TestBuildPhaseHandoffRecordTakesRunDir(t *testing.T) {
	rec := BuildPhaseHandoffRecord("/root", "/root/parley-deck/runs/20260924T010101.000000000Z", "idea", "/root/parley-deck/ideas/idea", []string{"a-1"},
		Cursor{Phase: PhaseImpl, CurrentRound: 1, IdeaStatus: "implemented"}, ActionPromoted, PhaseReview)
	if rec.RunID != "20260924T010101.000000000Z" {
		t.Fatalf("RunID must be the run dir's base name, got %q", rec.RunID)
	}
}

// TestCommitCursorAdvancesRunManifestUpdatedAt (fix-up F7, claude-1 MAJ-7): the run
// record is created by the REAL manifest machinery (created_at == updated_at at
// creation — the only shape the runner produces), and commitCursor's transition
// touch must open a real window (updated_at > created_at).
func TestCommitCursorAdvancesRunManifestUpdatedAt(t *testing.T) {
	ideaDir := t.TempDir()
	runDir := t.TempDir()
	root := filepath.Dir(filepath.Dir(ideaDir))
	runID := filepath.Base(runDir)
	created := time.Now().UTC().Add(-1 * time.Hour)
	manifest := runmanifest.New(runmanifest.Options{RunID: runID, Root: root, IdeaSlug: "handoff-idea", CreatedAt: created, UpdatedAt: created})
	if err := runmanifest.Write(root, runID, manifest); err != nil {
		t.Fatal(err)
	}
	events := store.New(runDir)
	d := &Driver{cfg: Config{
		IdeaDir: ideaDir, IdeaSlug: "handoff-idea", Participants: []string{"a-1"},
		RunDir: runDir, Root: root, Events: events,
	}}
	c := Cursor{Phase: PhaseConsensus, CurrentRound: 1, IdeaStatus: "consensus"}
	if err := d.commitCursor(c, ActionConsensusDrafted, PhaseRound); err != nil {
		t.Fatalf("commitCursor: %v", err)
	}
	touched, err := runmanifest.Load(root, runID)
	if err != nil {
		t.Fatal(err)
	}
	if !touched.UpdatedAt.After(created) {
		t.Fatalf("commitCursor must advance run.json updated_at beyond created_at; created=%s updated=%s", created, touched.UpdatedAt)
	}
}
