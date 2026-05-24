package agents

import "testing"

func TestACPCatalogCoversKnownBackends(t *testing.T) {
	wantBinaries := map[string]string{
		"claude-acp":  "claude",
		"qwen":        "qwen",
		"codex-acp":   "codex",
		"codebuddy":   "codebuddy",
		"goose":       "goose",
		"auggie":      "auggie",
		"kimi":        "kimi",
		"opencode":    "opencode",
		"droid":       "droid",
		"copilot":     "copilot",
		"qoder":       "qodercli",
		"vibe":        "vibe-acp",
		"cursor":      "agent",
		"kiro":        "kiro-cli",
		"hermes-acp":  "hermes",
		"snow":        "snow",
	}

	got := make(map[string]string, len(ACPCatalog()))
	for _, backend := range ACPCatalog() {
		got[backend.ID] = backend.Command
	}

	for id, cmd := range wantBinaries {
		if got[id] != cmd {
			t.Errorf("backend %q: want command %q, got %q", id, cmd, got[id])
		}
	}
	if len(got) != len(wantBinaries) {
		t.Errorf("catalog size mismatch: want %d, got %d (got=%v)", len(wantBinaries), len(got), got)
	}
}

func TestACPSpecsAreMarkedLaunchACP(t *testing.T) {
	specs := ACPSpecs()
	if len(specs) == 0 {
		t.Fatal("ACPSpecs returned no entries")
	}
	for _, spec := range specs {
		if spec.LaunchMode != LaunchACP {
			t.Errorf("spec %q has launch_mode %q, want %q", spec.ID, spec.LaunchMode, LaunchACP)
		}
		if len(spec.Commands) == 0 {
			t.Errorf("spec %q has no Commands", spec.ID)
		}
	}
}

func TestDefaultSpecsMergesACPCatalog(t *testing.T) {
	ids := map[string]bool{}
	for _, spec := range DefaultSpecs() {
		ids[spec.ID] = true
	}
	for _, required := range []string{"codex", "claude", "gemini", "hermes", "goose", "qwen", "opencode"} {
		if !ids[required] {
			t.Errorf("DefaultSpecs missing %q", required)
		}
	}
}
