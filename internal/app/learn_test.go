package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupClosedIdea(t *testing.T, status string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "parley-deck", "ideas", "demo-idea")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte("---\nidea: demo-idea\ntrack: deliberation\nparticipants: [claude-1, codex-1]\nstatus: complete\n---\n\n## Problem\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "IMPLEMENTATION.md"), []byte("---\nidea: demo-idea\nstatus: "+status+"\n---\n\n## Summary of work\n\ndone\n\n## Fix-up cycle 1\nstatus: complete\n"), 0o644)
	return root
}

func TestLearnWritesPlaybook(t *testing.T) {
	root := setupClosedIdea(t, "complete")
	var out, errb bytes.Buffer
	if code := runLearn([]string{"demo-idea", "--dir", root}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	pb := filepath.Join(root, "parley-deck", "playbooks", "demo-idea.md")
	body, err := os.ReadFile(pb)
	if err != nil {
		t.Fatalf("playbook not written: %v", err)
	}
	s := string(body)
	for _, want := range []string{"status: advisory", "distilled-from: ideas/demo-idea", "Track: deliberation", "## Step checklist"} {
		if !strings.Contains(s, want) {
			t.Errorf("playbook missing %q", want)
		}
	}
}

func TestLearnRejectsIncomplete(t *testing.T) {
	root := setupClosedIdea(t, "implemented") // not complete
	var out, errb bytes.Buffer
	if code := runLearn([]string{"demo-idea", "--dir", root}, &out, &errb); code == 0 {
		t.Fatal("should reject a non-complete idea")
	}
	if !strings.Contains(errb.String(), "not complete") {
		t.Fatalf("stderr = %q", errb.String())
	}
}

func TestLearnFailsClosedOnExisting(t *testing.T) {
	root := setupClosedIdea(t, "complete")
	var out, errb bytes.Buffer
	runLearn([]string{"demo-idea", "--dir", root}, &out, &errb)
	errb.Reset()
	if code := runLearn([]string{"demo-idea", "--dir", root}, &out, &errb); code == 0 {
		t.Fatal("second run should fail closed (target exists)")
	}
	if !strings.Contains(errb.String(), "already exists") {
		t.Fatalf("stderr = %q", errb.String())
	}
}

func TestLearnFlagAfterSlug(t *testing.T) {
	// `parley learn <slug> --topic X` must work (Go flag stops at the first positional).
	root := setupClosedIdea(t, "complete")
	var out, errb bytes.Buffer
	if code := runLearn([]string{"demo-idea", "--topic", "my-topic", "--dir", root}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	if _, err := os.Stat(filepath.Join(root, "parley-deck", "playbooks", "my-topic.md")); err != nil {
		t.Fatalf("--topic target not written: %v", err)
	}
}

func TestLearnRejectsSymlinkedPlaybooksDir(t *testing.T) {
	root := setupClosedIdea(t, "complete")
	outside := t.TempDir()
	deck := filepath.Join(root, "parley-deck")
	// Symlink parley-deck/playbooks → outside dir.
	if err := os.Symlink(outside, filepath.Join(deck, "playbooks")); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	if code := runLearn([]string{"demo-idea", "--dir", root}, &out, &errb); code == 0 {
		t.Fatal("symlinked playbooks/ must fail closed")
	}
	if _, err := os.Stat(filepath.Join(outside, "demo-idea.md")); err == nil {
		t.Fatal("must NOT write into the symlink target")
	}
	if !strings.Contains(errb.String(), "symlink") {
		t.Fatalf("stderr = %q", errb.String())
	}
}
