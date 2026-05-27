package agents

import "testing"

func TestACPCatalogCoversKnownBackends(t *testing.T) {
	wantBinaries := map[string]string{
		"claude-acp": "claude",
		"qwen":       "qwen",
		"codex-acp":  "codex",
		"codebuddy":  "codebuddy",
		"goose":      "goose",
		"auggie":     "auggie",
		"kimi":       "kimi",
		"opencode":   "opencode",
		"droid":      "droid",
		"copilot":    "copilot",
		"qoder":      "qodercli",
		"vibe":       "vibe-acp",
		"cursor":     "agent",
		"kiro":       "kiro-cli",
		"hermes-acp": "hermes",
		"snow":       "snow",
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
	for _, required := range []string{"codex", "claude", "agy", "gemini", "hermes", "goose", "qwen", "opencode"} {
		if !ids[required] {
			t.Errorf("DefaultSpecs missing %q", required)
		}
	}
}

func TestDefaultSpecsPreferAntigravityAndStrongVerifiedDefaults(t *testing.T) {
	specs := DefaultSpecs()
	agy := findSpecForTest(specs, "agy")
	if agy.ID == "" {
		t.Fatal("DefaultSpecs missing agy")
	}
	if got := agy.Commands[0]; got != "agy" {
		t.Fatalf("agy command=%q, want agy", got)
	}
	for _, want := range []string{"--print-timeout", "30m", "--dangerously-skip-permissions", "--add-dir", "{root}", "--print", "{prompt}"} {
		if !containsStringForTest(agy.HeadlessArgs, want) {
			t.Fatalf("agy HeadlessArgs missing %q: %v", want, agy.HeadlessArgs)
		}
	}
	if agy.PromptMode != PromptArg {
		t.Fatalf("agy PromptMode=%q, want %q", agy.PromptMode, PromptArg)
	}
	if agy.Model != CLIDefault || agy.Reasoning != CLIDefault {
		t.Fatalf("agy model/reasoning=%q/%q, want cli-default", agy.Model, agy.Reasoning)
	}

	claude := findSpecForTest(specs, "claude")
	if claude.Model != "opus" || claude.Reasoning != "max" {
		t.Fatalf("claude model/reasoning=%q/%q, want opus/max", claude.Model, claude.Reasoning)
	}

	hermes := findSpecForTest(specs, "hermes")
	if hermes.Model != "xai/grok-4.3" {
		t.Fatalf("hermes model=%q, want xai/grok-4.3", hermes.Model)
	}

	gemini := findSpecForTest(specs, "gemini")
	if !containsStringForTest([]string{gemini.Notes}, "DEPRECATED legacy Gemini CLI support; prefer agy. Uses isolated GEMINI_CLI_HOME for oauth-personal profiles that hang") {
		t.Fatalf("gemini notes do not mark legacy: %q", gemini.Notes)
	}
}

func findSpecForTest(specs []Spec, id string) Spec {
	for _, spec := range specs {
		if spec.ID == id {
			return spec
		}
	}
	return Spec{}
}

func containsStringForTest(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
