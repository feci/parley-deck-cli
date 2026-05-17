package sessionstore

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultDirUsesParleyHome(t *testing.T) {
	root := t.TempDir()
	t.Setenv(homeEnv, root)

	dir, err := DefaultDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != root {
		t.Fatalf("dir=%q, want %q", dir, root)
	}

	store, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	if got := store.Path(); got != filepath.Join(root, "sessions.json") {
		t.Fatalf("path=%q", got)
	}
}

func TestUpsertListAndDedupe(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sessions.json")
	store := New(path)
	base := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)

	if err := store.Upsert(Session{
		WorkspaceRoot: ".",
		RunID:         "run-1",
		IdeaSlug:      "first",
		Task:          "Task",
		Participants:  []string{"codex"},
		CreatedAt:     base,
		UpdatedAt:     base,
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Upsert(Session{
		WorkspaceRoot: ".",
		RunID:         "run-1",
		LastEventAt:   base.Add(time.Minute),
		UpdatedAt:     base.Add(time.Minute),
		Terminal:      true,
	}); err != nil {
		t.Fatal(err)
	}

	sessions, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("sessions=%d, want 1", len(sessions))
	}
	session := sessions[0]
	if session.IdeaSlug != "first" || session.Task != "Task" || len(session.Participants) != 1 {
		t.Fatalf("metadata was not preserved on update: %+v", session)
	}
	if !session.Terminal || !session.LastEventAt.Equal(base.Add(time.Minute)) {
		t.Fatalf("update fields missing: %+v", session)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("sessions file missing: %v", err)
	}
}

func TestListMissingFile(t *testing.T) {
	sessions, err := New(filepath.Join(t.TempDir(), "missing", "sessions.json")).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 0 {
		t.Fatalf("sessions=%+v, want empty", sessions)
	}
}
