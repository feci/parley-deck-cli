//go:build windows

package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

// §D.6 read-only-attribute trap: replacing a FILE_ATTRIBUTE_READONLY target
// must succeed after the write bit is restored — POSIX rename is not blocked
// by the target's own permissions, and the Windows path must match.
func TestReplaceSyncedFileOverReadOnlyTarget(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "state.json")
	if err := os.WriteFile(target, []byte("old"), 0o400); err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(dir, "state.json.staged")
	if err := os.WriteFile(staged, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ReplaceSyncedFile(staged, target); err != nil {
		t.Fatalf("replace over read-only target: %v", err)
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "new" {
		t.Fatalf("target content wrong: %q err=%v", data, err)
	}
	// The ordinary writable-target replace keeps working.
	if err := os.WriteFile(staged, []byte("newer"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ReplaceSyncedFile(staged, target); err != nil {
		t.Fatalf("plain replace: %v", err)
	}
	if data, _ := os.ReadFile(target); string(data) != "newer" {
		t.Fatalf("plain replace content wrong: %q", data)
	}
}
