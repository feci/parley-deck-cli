package protocol

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "00-prompt.md")
	content := "---\nidea: sample\nparticipants: [codex, claude]\nstatus: round-01\n---\n\nbody\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	meta, err := ReadFrontmatter(path)
	if err != nil {
		t.Fatal(err)
	}
	if meta["idea"] != "sample" {
		t.Fatalf("idea = %q", meta["idea"])
	}
	participants := parseList(meta["participants"])
	if len(participants) != 2 || participants[0] != "codex" || participants[1] != "claude" {
		t.Fatalf("participants = %#v", participants)
	}
}

func TestInitWorkspaceAndStatus(t *testing.T) {
	dir := t.TempDir()
	if err := InitWorkspace(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, DeckDir, "COOPERATION.md")); err != nil {
		t.Fatal(err)
	}

	status, err := ReadWorkspaceStatus(dir)
	if err != nil {
		t.Fatal(err)
	}
	if status.Transport != "local-dir" {
		t.Fatalf("transport = %q", status.Transport)
	}
}

func TestTimestampedSlugShortensLongTask(t *testing.T) {
	now := time.Date(2026, 5, 17, 21, 39, 35, 0, time.UTC)
	got := timestampedSlug("podme nastartovat novu ideu, skus vymysliet hru na roblox", now)
	want := "2026-05-17T21-39-35-podme-nastartova"
	if got != want {
		t.Fatalf("slug=%q want %q", got, want)
	}
}

func TestTimestampedSlugFallsBackForEmptyTask(t *testing.T) {
	now := time.Date(2026, 5, 17, 21, 39, 35, 0, time.UTC)
	got := timestampedSlug("!!!", now)
	want := "2026-05-17T21-39-35-task"
	if got != want {
		t.Fatalf("slug=%q want %q", got, want)
	}
}
