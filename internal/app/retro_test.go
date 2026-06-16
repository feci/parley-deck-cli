package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/retro"
)

func TestRetroProposeWritesOnly00PromptAndFailsClosed(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	coreset := []retro.IdeaSignals{{Slug: "x", FailureType: "fix-up-heavy", Reasons: []string{"fix-up churn"}}}

	var out, errb bytes.Buffer
	if code := retroPropose(root, "retro-test-idea", coreset, &out, &errb); code != 0 {
		t.Fatalf("propose code=%d err=%s", code, errb.String())
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", "retro-test-idea")
	promptPath := filepath.Join(ideaDir, "00-prompt.md")
	body, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("00-prompt.md not written: %v", err)
	}
	for _, want := range []string{"idea: retro-test-idea", "drafted-by: parley retro", "advisory", "Failure mode: fix-up-heavy"} {
		if !bytes.Contains(body, []byte(want)) {
			t.Fatalf("00-prompt missing %q:\n%s", want, body)
		}
	}
	// Only 00-prompt.md is written — no round/review/participant files.
	entries, _ := os.ReadDir(ideaDir)
	if len(entries) != 1 || entries[0].Name() != "00-prompt.md" {
		t.Fatalf("propose wrote more than 00-prompt.md: %v", entries)
	}

	// Fail-closed: a second propose at the same slug must refuse (no overwrite).
	out.Reset()
	errb.Reset()
	if code := retroPropose(root, "retro-test-idea", coreset, &out, &errb); code == 0 {
		t.Fatal("propose must fail when 00-prompt.md already exists")
	}

	// Invalid slugs are rejected (path sep, dotfile, space, uppercase, double-hyphen).
	for _, bad := range []string{"", "a/b", ".hidden", "has space", "UPPER", "a--b"} {
		if code := retroPropose(root, bad, coreset, &out, &errb); code == 0 {
			t.Fatalf("propose must reject invalid slug %q", bad)
		}
	}

	// Fail-closed when the slug directory already exists WITHOUT 00-prompt.md
	// (codex: do not write into a pre-existing idea dir).
	emptyDir := filepath.Join(root, protocol.DeckDir, "ideas", "preexisting-dir")
	if err := os.MkdirAll(emptyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if code := retroPropose(root, "preexisting-dir", coreset, &out, &errb); code == 0 {
		t.Fatal("propose must fail when ideas/<slug> already exists without 00-prompt.md")
	}
	if _, err := os.Stat(filepath.Join(emptyDir, "00-prompt.md")); !os.IsNotExist(err) {
		t.Fatal("propose must not write into a pre-existing slug dir")
	}

	// Fail-closed on a symlinked slug entry (Lstat must not follow the link).
	target := t.TempDir()
	linkPath := filepath.Join(root, protocol.DeckDir, "ideas", "symlinked")
	if err := os.Symlink(target, linkPath); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	if code := retroPropose(root, "symlinked", coreset, &out, &errb); code == 0 {
		t.Fatal("propose must fail on a symlinked slug entry")
	}
	if _, err := os.Stat(filepath.Join(target, "00-prompt.md")); !os.IsNotExist(err) {
		t.Fatal("propose must not write through a symlinked slug")
	}
}
