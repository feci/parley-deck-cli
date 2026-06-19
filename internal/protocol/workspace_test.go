package protocol

import (
	"os"
	"path/filepath"
	"strings"
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
	cooperation, err := os.ReadFile(filepath.Join(dir, DeckDir, "COOPERATION.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(cooperation)
	if !strings.Contains(text, "## 11. Transport mechanics") {
		t.Fatalf("COOPERATION.md does not contain the full protocol")
	}
	if strings.Contains(text, "Replace this bootstrap file") {
		t.Fatalf("COOPERATION.md still contains the old bootstrap placeholder")
	}

	status, err := ReadWorkspaceStatus(dir)
	if err != nil {
		t.Fatal(err)
	}
	if status.Transport != "local-dir" {
		t.Fatalf("transport = %q", status.Transport)
	}
}

func TestInitWorkspaceWritesConsumerVersionMeta(t *testing.T) {
	dir := t.TempDir()
	if err := InitWorkspace(dir); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, DeckDir, "meta", "version.json"))
	if err != nil {
		t.Fatalf("init did not write meta/version.json: %v", err)
	}
	if !strings.Contains(string(data), `"protocolRole": "consumer"`) {
		t.Fatalf("version.json missing protocolRole=consumer:\n%s", data)
	}
}

func TestCreateIdeaWithExclusionsRecordsExcluded(t *testing.T) {
	dir := t.TempDir()
	if err := InitWorkspace(dir); err != nil {
		t.Fatal(err)
	}
	idea, err := CreateIdeaWithExclusions(dir, "task body", []string{"codex", "claude"},
		[]string{"hermes — unavailable:no-pong — confirmed 2026-06-19"})
	if err != nil {
		t.Fatal(err)
	}
	prompt, err := os.ReadFile(filepath.Join(idea.Path, "00-prompt.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(prompt)
	if !strings.Contains(text, "participants: [codex, claude]") {
		t.Fatalf("00-prompt.md missing participants:\n%s", text)
	}
	if !strings.Contains(text, "excluded: hermes — unavailable:no-pong — confirmed 2026-06-19") {
		t.Fatalf("00-prompt.md did not record the confirmed exclusion:\n%s", text)
	}
}

func TestCreateIdeaNoExclusionsOmitsExcludedLine(t *testing.T) {
	dir := t.TempDir()
	if err := InitWorkspace(dir); err != nil {
		t.Fatal(err)
	}
	idea, err := CreateIdea(dir, "task body", []string{"codex", "claude"})
	if err != nil {
		t.Fatal(err)
	}
	prompt, _ := os.ReadFile(filepath.Join(idea.Path, "00-prompt.md"))
	if strings.Contains(string(prompt), "excluded:") {
		t.Fatalf("00-prompt.md should carry no excluded: line when none confirmed:\n%s", prompt)
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
