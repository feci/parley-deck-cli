package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

// zcodeDeck builds a deck whose roster is a single zcode-1, with a fake HOME carrying the
// zcode config files the adapter reads at launch.
func zcodeDeck(t *testing.T, deckModel string) string {
	t.Helper()
	root := setupRosterDeck(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	mk := func(rel, body string) {
		p := filepath.Join(home, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mk(".zcode/cli/config.json", `{"model":{"main":"zai/glm-5.3"}}`)
	mk(".zcode/v2/config.json",
		`{"provider":{"builtin:zai":{"models":{"GLM-5.3":{"reasoning":{"defaultVariant":"max"}}}}}}`)

	entry := "[roster.zcode-1]\nadapter = \"zcode\"\nspeed = \"deep\"\n"
	if deckModel != "" {
		entry += "model = \"" + deckModel + "\"\n"
	}
	if err := os.WriteFile(filepath.Join(root, protocol.DeckDir, "agents.toml"), []byte(entry), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func zcodeRow(t *testing.T, root string) []string {
	t.Helper()
	var out, errb bytes.Buffer
	if code := runRoster([]string{"show", "--dir", root}, &out, &errb); code != 0 {
		t.Fatalf("roster show exit=%d stderr=%s", code, errb.String())
	}
	for _, line := range strings.Split(out.String(), "\n") {
		if strings.HasPrefix(line, "zcode-1") {
			return strings.Fields(line)
		}
	}
	t.Fatalf("no zcode-1 row in:\n%s", out.String())
	return nil
}

// An adapter whose CLI has no model flag reports the model from the agent's OWN config, so the
// operator sees what will actually run instead of `unknown`.
func TestZcodeRowReportsModelAndEffortReadFromAgentConfig(t *testing.T) {
	fields := zcodeRow(t, zcodeDeck(t, ""))
	joined := strings.Join(fields, " ")
	if !strings.Contains(joined, "zai/glm-5.3") {
		t.Errorf("model not read from zcode config: %s", joined)
	}
	if !strings.Contains(joined, "model-from-config") || !strings.Contains(joined, "effort-from-config") {
		t.Errorf("status must record that the values were read, not bound: %s", joined)
	}
	// Provenance must never be reported as `ok` — `ok` means the launch carries the value.
	if strings.Contains(joined, " ok") {
		t.Errorf("a config-read row must not claim ok: %s", joined)
	}
}

// The reason the roster contract exists: a model declared in a parley config for an adapter that
// cannot carry one must NOT appear, because the launch will not pass it. Reading the agent's own
// config is allowed precisely because it is a different source; a deck declaration is not.
func TestDeckDeclaredModelNeverOverridesAgentConfigForUnbindableAdapter(t *testing.T) {
	joined := strings.Join(zcodeRow(t, zcodeDeck(t, "zai/glm-9-imaginary")), " ")
	if strings.Contains(joined, "glm-9-imaginary") {
		t.Fatalf("a deck declaration the launch cannot pass leaked into the row: %s", joined)
	}
	if !strings.Contains(joined, "zai/glm-5.3") {
		t.Fatalf("agent-config model missing: %s", joined)
	}
}

// With no agent config present there is nothing to read, and the row must fall back to the
// honest unknown rather than to any parley-side declaration.
func TestUnreadableAgentConfigFallsBackToUnknown(t *testing.T) {
	root := zcodeDeck(t, "zai/glm-9-imaginary")
	if err := os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".zcode")); err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(zcodeRow(t, root), " ")
	if !strings.Contains(joined, "model-unbound") || !strings.Contains(joined, "effort-unknown") {
		t.Fatalf("want unbound/unknown fallback, got: %s", joined)
	}
	if strings.Contains(joined, "glm-9-imaginary") {
		t.Fatalf("declaration leaked on fallback: %s", joined)
	}
}
