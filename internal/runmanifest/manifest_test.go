package runmanifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteLoadAndDefaults(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	manifest := New(Options{
		Root:         root,
		RunID:        "run-1",
		IdeaSlug:     "idea",
		Task:         "Task",
		Mode:         "hitl",
		Transport:    "github-pr",
		Phase:        "round",
		IdeaStatus:   "round-01",
		CurrentRound: "round-01",
		ActiveSteps:  []Step{{ID: "round-01.codex", Kind: "round", AgentID: "codex", Status: "running"}},
		NextActions:  []NextAction{{ID: "inspect.run", Kind: "inspect", RunID: "run-1", IdeaSlug: "idea", Risk: "low", Summary: "Inspect run state"}},
		Participants: []string{"codex"},
		CreatedAt:    now,
	})
	if err := Write(root, manifest.RunID, manifest); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(root, "parley-deck", "runs", "run-1", "run.json"))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if raw["schema_version"] != float64(SchemaVersion) || raw["status"] != StatusRunning || raw["created_at"] != "2026-05-18T12:00:00Z" {
		t.Fatalf("raw manifest=%+v", raw)
	}

	loaded, err := Load(root, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.RunID != "run-1" || loaded.IdeaSlug != "idea" || loaded.Status != StatusRunning || loaded.Transport != "github-pr" {
		t.Fatalf("loaded=%+v", loaded)
	}
	if loaded.Phase != "round" || loaded.IdeaStatus != "round-01" || loaded.CurrentRound != "round-01" {
		t.Fatalf("loaded continuation fields=%+v", loaded)
	}
	if len(loaded.ActiveSteps) != 1 || loaded.ActiveSteps[0].ID != "round-01.codex" {
		t.Fatalf("loaded active steps=%+v", loaded.ActiveSteps)
	}
	if len(loaded.NextActions) != 1 || loaded.NextActions[0].Kind != "inspect" {
		t.Fatalf("loaded next actions=%+v", loaded.NextActions)
	}
	if loaded.WorkspaceRoot != root {
		t.Fatalf("workspace root=%q want %q", loaded.WorkspaceRoot, root)
	}
}

func TestPathNormalizesRoot(t *testing.T) {
	path := Path(".", "run-1")
	if !filepath.IsAbs(path) {
		t.Fatalf("path=%q is not absolute", path)
	}
}
