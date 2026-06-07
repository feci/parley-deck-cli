package runcontrol

import (
	"errors"
	"os"
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

// A transient manifest-write failure (e.g. a virtio-fs mkdir hiccup) must NOT orphan the
// run: Create still succeeds, the run is defined by events.jsonl, and the deferral is
// recorded for the audit trail.
func TestCreateBestEffortManifest(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	orig := writeManifest
	writeManifest = func(string, string, runmanifest.Manifest) error { return errors.New("simulated manifest failure") }
	t.Cleanup(func() { writeManifest = orig })

	created, err := Create(CreateOptions{
		Root:         root,
		Task:         "hardening task",
		Participants: []string{"codex"},
		Now:          time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create must not fail when the manifest write fails: %v", err)
	}
	if created.RunID == "" || created.RunDir == "" {
		t.Fatalf("run must still be created: %+v", created)
	}
	if _, err := os.Stat(runmanifest.Path(root, created.RunID)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("run.json must be absent (the write was forced to fail); stat err=%v", err)
	}
	events, err := store.New(created.RunDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	var hasCreated, hasDeferred bool
	for _, e := range events {
		switch e.Type {
		case "run.created":
			hasCreated = true
		case "run.manifest_deferred":
			hasDeferred = true
		}
	}
	if !hasCreated || !hasDeferred {
		t.Fatalf("want run.created + run.manifest_deferred audit events; got %+v", events)
	}
}
