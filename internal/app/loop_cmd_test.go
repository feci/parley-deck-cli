package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

// `loop tick` with no config is disabled: exits 0, says disabled, writes no idea.
func TestLoopTickDisabledByDefault(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	sig := filepath.Join(t.TempDir(), "signals.json")
	os.WriteFile(sig, []byte(`[{"source":"commit","id":"abc","title":"x"}]`), 0o644)

	var out, errb bytes.Buffer
	if code := runLoopTick([]string{"--dir", root, "--signals", sig}, &out, &errb); code != 0 {
		t.Fatalf("disabled tick should exit 0; code=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "disabled") {
		t.Fatalf("expected a disabled notice, got: %q", out.String())
	}
	ideas, _ := os.ReadDir(filepath.Join(root, protocol.DeckDir, "ideas"))
	if len(ideas) != 0 {
		t.Fatalf("disabled tick must write no idea; got %d", len(ideas))
	}
}

// `loop tick --enable` drafts a candidate idea (and only a 00-prompt.md, status: candidate).
func TestLoopTickEnableDraftsCandidateOnly(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	sig := filepath.Join(t.TempDir(), "signals.json")
	os.WriteFile(sig, []byte(`[{"source":"ci","id":"build-9","title":"flaky","detail":"retry"}]`), 0o644)

	var out, errb bytes.Buffer
	if code := runLoopTick([]string{"--dir", root, "--signals", sig, "--enable"}, &out, &errb); code != 0 {
		t.Fatalf("enabled tick code=%d err=%s", code, errb.String())
	}
	ideasDir := filepath.Join(root, protocol.DeckDir, "ideas")
	ideas, _ := os.ReadDir(ideasDir)
	if len(ideas) != 1 {
		t.Fatalf("enabled tick should draft exactly one candidate; got %d", len(ideas))
	}
	slug := ideas[0].Name()
	entries, _ := os.ReadDir(filepath.Join(ideasDir, slug))
	if len(entries) != 1 || entries[0].Name() != "00-prompt.md" {
		t.Fatalf("a candidate must be 00-prompt.md only (no round/quorum); got %v", entries)
	}
	body, _ := os.ReadFile(filepath.Join(ideasDir, slug, "00-prompt.md"))
	if !strings.Contains(string(body), "status: candidate") {
		t.Fatalf("candidate must be status: candidate:\n%s", body)
	}
	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "participants:") {
			t.Fatalf("candidate must NOT claim a participants: quorum (§14):\n%s", body)
		}
	}
}

// `loop tick --json` emits the TickResult shape with the enabled flag.
func TestLoopTickJSON(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	if code := runLoopTick([]string{"--dir", root, "--json"}, &out, &errb); code != 0 {
		t.Fatalf("json tick code=%d err=%s", code, errb.String())
	}
	var res struct {
		Enabled bool     `json:"enabled"`
		Created []string `json:"created"`
	}
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if res.Enabled {
		t.Fatal("default (no config) must report enabled=false")
	}
}
