package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

// Component D (composite-agent-naming-and-roster-reinit): `fast` is a SEPARATE
// axis from reasoning effort. It must never downgrade the pinned model or effort
// (unlike the banned legacy "fast profile = weaker model + lower thinking"). This
// is the regression tripwire for any future code path that maps speed -> model.
func TestSpeedFastDoesNotDowngradeModelOrEffort(t *testing.T) {
	root := t.TempDir()
	t.Setenv(EnvParleyHome, filepath.Join(root, "no-central")) // isolate from real ~/.parley
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	project := filepath.Join(root, protocol.DeckDir, "agents.toml")
	if err := os.WriteFile(project, []byte(`
[defaults]
speed = "fast"

[agents.codex]
model = "gpt-5.6-sol"
model_label = "GPT-5.6 Sol"
reasoning = "xhigh"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	specs, err := LoadAgentSpecs(root)
	if err != nil {
		t.Fatal(err)
	}
	codex := findSpec(t, specs, "codex")
	if codex.Speed != "fast" {
		t.Errorf("speed=%q, want fast", codex.Speed)
	}
	if codex.Model != "gpt-5.6-sol" {
		t.Errorf("fast downgraded the model to %q, want gpt-5.6-sol", codex.Model)
	}
	if codex.Reasoning != "xhigh" {
		t.Errorf("fast downgraded the effort to %q, want xhigh", codex.Reasoning)
	}
	if codex.ModelLabel != "GPT-5.6 Sol" {
		t.Errorf("model_label=%q, want 'GPT-5.6 Sol'", codex.ModelLabel)
	}
}

func TestCentralDefaultTemplateSpeedFast(t *testing.T) {
	tmpl := centralDefaultTemplate()
	if !strings.Contains(tmpl, `speed = "fast"`) {
		t.Error("central template should default speed to fast")
	}
	if strings.Contains(tmpl, `speed = "deep"`) {
		t.Error("central template must not seed the legacy speed = deep default")
	}
}
