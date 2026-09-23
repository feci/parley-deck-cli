package driver

// lean-organizer D.1: the driver writes a schema-valid per-phase handoff record at
// each phase transition, on the same write discipline as the shipped handoff
// machinery, with the schema doc stating recomputation is authoritative.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	if rec.Idea != "handoff-idea" || rec.Phase != "consensus" || rec.RunID == "" {
		t.Errorf("record fields incomplete: %+v", rec)
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

func TestPhaseHandoffRoundTrip(t *testing.T) {
	dir := t.TempDir()
	rec := PhaseHandoffRecord{RunID: "r1", Idea: "i", Phase: "round", Action: "promoted", Next: NextAwaitRoundArtifact}
	path := filepath.Join(dir, PhaseHandoffPath(PhaseRound))
	if err := storeHandoffRecord(path, rec); err != nil {
		t.Fatal(err)
	}
	back, err := LoadPhaseHandoffRecord(path)
	if err != nil {
		t.Fatal(err)
	}
	if back.RunID != rec.RunID || back.Idea != rec.Idea || back.Phase != rec.Phase || back.Next != rec.Next {
		t.Errorf("round-trip mismatch: %+v vs %+v", back, rec)
	}
}
