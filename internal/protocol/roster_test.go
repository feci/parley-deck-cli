package protocol

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadRosterIDs(t *testing.T) {
	root := t.TempDir()
	deck := filepath.Join(root, DeckDir)
	if err := os.MkdirAll(deck, 0o755); err != nil {
		t.Fatal(err)
	}
	coop := "## 2. Active agents (roster)\n\n" +
		"| Agent ID       | Workspace dir | Role |\n" +
		"| -------------- | ------------- | ---- |\n" +
		"| `claude-1`      | `../claude/`  | `facilitator+participant` |\n" +
		"| `codex-1`       | `../codex/`   | `participant` |\n" +
		"| `gemini-1`      | `../gemini/`  | `participant` (inactive — retired) |\n" +
		"\nsome prose after the table\n\n" +
		"| Agent ID | Host handle |\n" +
		"| --- | --- |\n" +
		"| `claude-1` | `feci` |\n"
	if err := os.WriteFile(filepath.Join(deck, "COOPERATION.md"), []byte(coop), 0o644); err != nil {
		t.Fatal(err)
	}

	active, inactive, ok := ReadRosterIDs(root)
	if !ok {
		t.Fatal("expected ok")
	}
	for _, id := range []string{"claude-1", "codex-1", "gemini-1"} {
		if !active[id] {
			t.Fatalf("missing active id %q", id)
		}
	}
	// The host-handle table must NOT add rows (only the first roster table is read),
	// and it must not have leaked a second parse — claude-1 already present is fine.
	if len(active) != 3 {
		t.Fatalf("want 3 roster IDs, got %d: %v", len(active), active)
	}
	if !inactive["gemini-1"] {
		t.Fatal("gemini-1 should be inactive")
	}
	if inactive["claude-1"] {
		t.Fatal("claude-1 should not be inactive")
	}
}

func TestReadRosterIDsMissing(t *testing.T) {
	if _, _, ok := ReadRosterIDs(t.TempDir()); ok {
		t.Fatal("missing COOPERATION.md should return ok=false")
	}
}
