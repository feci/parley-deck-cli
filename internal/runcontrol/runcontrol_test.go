package runcontrol

import (
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/sessionstore"
	"parley-deck-cli/internal/store"
)

func TestCreateWritesRunCreatedAndRegistersSession(t *testing.T) {
	root := t.TempDir()
	parleyHome := t.TempDir()
	t.Setenv("PARLEY_HOME", parleyHome)
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)

	created, err := Create(CreateOptions{
		Root:         root,
		Task:         "Session task",
		Participants: []string{"codex"},
		Discovered: []agents.Discovery{{
			Spec:  agents.Spec{ID: "codex", Model: "local-model"},
			Found: true,
		}},
		Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Idea.Slug == "" || created.RunID == "" || created.RunDir == "" {
		t.Fatalf("created=%+v", created)
	}

	events, err := store.New(created.RunDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Type != "run.created" {
		t.Fatalf("events=%+v", events)
	}
	if events[0].Data["task"] != "Session task" || events[0].Data["idea"] != created.Idea.Slug {
		t.Fatalf("run.created data=%+v", events[0].Data)
	}
	manifest, err := runmanifest.Load(root, created.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.RunID != created.RunID || manifest.IdeaSlug != created.Idea.Slug || manifest.Task != "Session task" || manifest.Mode != "hitl" || manifest.Status != runmanifest.StatusRunning {
		t.Fatalf("manifest=%+v", manifest)
	}

	sessions, err := sessionstore.New(filepath.Join(parleyHome, "sessions.json")).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("sessions=%+v", sessions)
	}
	if sessions[0].RunID != created.RunID || sessions[0].IdeaSlug != created.Idea.Slug || sessions[0].Task != "Session task" {
		t.Fatalf("session=%+v", sessions[0])
	}
}
