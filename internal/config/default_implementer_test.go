package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

// T-3 (AC-7): default_implementer merges on the existing non-empty-string pattern
// across the REAL layered config files: a higher non-empty layer wins; an empty value
// at a higher layer does NOT clear a lower layer; "none" at a higher layer wins the
// merge (and resolves as suppression downstream).
func TestDefaultImplementerLayerPrecedence(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvParleyHome, home)
	t.Setenv(EnvAgentConfig, "") // keep the env layer out of this test
	root := t.TempDir()
	deck := filepath.Join(root, protocol.DeckDir)
	if err := os.MkdirAll(deck, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(path, body string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Two real layers: central (low) and deck (high) — the higher non-empty wins.
	write(filepath.Join(home, "agents.toml"), "[defaults]\ndefault_implementer = \"central-agent\"\n")
	write(filepath.Join(deck, "agents.toml"), "[defaults]\ndefault_implementer = \"deck-agent\"\n")
	defs, err := LoadDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if defs.DefaultImplementer != "deck-agent" {
		t.Fatalf("higher non-empty layer must win, got %q", defs.DefaultImplementer)
	}

	// An empty value at the higher layer does NOT clear the lower layer (R9).
	write(filepath.Join(deck, "agents.toml"), "[defaults]\ndefault_implementer = \"\"\n")
	defs, err = LoadDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if defs.DefaultImplementer != "central-agent" {
		t.Fatalf("empty higher layer must not clear the lower layer, got %q", defs.DefaultImplementer)
	}

	// The documented suppressor is non-empty, so it wins the merge (R9).
	write(filepath.Join(deck, "agents.toml"), "[defaults]\ndefault_implementer = \"none\"\n")
	defs, err = LoadDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if defs.DefaultImplementer != "none" {
		t.Fatalf("\"none\" must win the merge, got %q", defs.DefaultImplementer)
	}

	// Unset everywhere stays unset (the shipped default).
	write(filepath.Join(home, "agents.toml"), "[defaults]\nspeed = \"fast\"\n")
	if err := os.Remove(filepath.Join(deck, "agents.toml")); err != nil {
		t.Fatal(err)
	}
	defs, err = LoadDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if defs.DefaultImplementer != "" {
		t.Fatalf("unset layers must resolve to unset, got %q", defs.DefaultImplementer)
	}
}

// AC-19 (R12, K-2): the freshly generated central config carries default_implementer
// COMMENTED OUT — a deliberate deviation from the template's active-key shape, because
// an active value would ship the global default SET.
func TestCentralDefaultTemplateShipsImplementerUnset(t *testing.T) {
	out := centralDefaultTemplate()
	foundComment := false
	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "default_implementer") {
			t.Fatalf("default_implementer must ship commented out, found active line: %q", line)
		}
		if strings.HasPrefix(trimmed, "# default_implementer") {
			foundComment = true
		}
	}
	if !foundComment {
		t.Fatal("central default template must document default_implementer as a commented-out key")
	}
}

// AC-19 (R12), second half: there is no deck agents.toml generator — a deck created
// by the tooling contains the key NOT AT ALL. The assertion targets deck CONFIG files
// (agents.toml / agents.local.toml): the generated COOPERATION.md prose legitimately
// documents the key (§0), which is not the key being set.
func TestDeckCreatedByToolingHasNoDefaultImplementer(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	found := false
	err := filepath.Walk(filepath.Join(root, protocol.DeckDir), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".toml") {
			return err
		}
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		if strings.Contains(string(data), "default_implementer") {
			found = true
			t.Logf("unexpected default_implementer mention in deck config %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("a deck created by the tooling must not carry default_implementer in any config file")
	}
}
