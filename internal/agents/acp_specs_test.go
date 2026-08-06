package agents

import (
	"strings"
	"testing"
)

func TestACPCatalogCoversKnownBackends(t *testing.T) {
	wantBinaries := map[string]string{
		"claude":    "claude",
		"qwen":      "qwen",
		"codex":     "codex",
		"codebuddy": "codebuddy",
		"goose":     "goose",
		"auggie":    "auggie",
		"kimi":      "kimi",
		"opencode":  "opencode",
		"droid":     "droid",
		"copilot":   "copilot",
		"qoder":     "qodercli",
		"vibe":      "vibe-acp",
		"cursor":    "agent",
		"kiro":      "kiro-cli",
		"hermes":    "hermes",
		"snow":      "snow",
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
	vibe := findSpecForTest(specs, "vibe")
	if vibe.ACPArgs == nil || len(vibe.ACPArgs) != 0 {
		t.Fatalf("vibe ACPArgs=%v, want configured empty list", vibe.ACPArgs)
	}
}

func TestDefaultSpecsMergesACPCatalog(t *testing.T) {
	ids := map[string]int{}
	for _, spec := range DefaultSpecs() {
		ids[spec.ID]++
	}
	for _, required := range []string{"codex", "claude", "agy", "gemini", "hermes", "goose", "qwen", "opencode"} {
		if ids[required] == 0 {
			t.Errorf("DefaultSpecs missing %q", required)
		}
	}
	for _, duplicate := range []string{"codex-acp", "claude-acp", "hermes-acp"} {
		if ids[duplicate] != 0 {
			t.Errorf("DefaultSpecs should not expose duplicate ACP agent %q", duplicate)
		}
	}
	for id, count := range ids {
		if count != 1 {
			t.Errorf("DefaultSpecs has %d entries for %q", count, id)
		}
	}

	claude := findSpecForTest(DefaultSpecs(), "claude")
	if !sameStringsForTest(claude.ACPArgs, []string{"--experimental-acp"}) {
		t.Fatalf("claude ACPArgs=%v", claude.ACPArgs)
	}
	hermes := findSpecForTest(DefaultSpecs(), "hermes")
	if !sameStringsForTest(hermes.ACPArgs, []string{"acp"}) {
		t.Fatalf("hermes ACPArgs=%v", hermes.ACPArgs)
	}
	codex := findSpecForTest(DefaultSpecs(), "codex")
	if codex.ACPArgs != nil {
		t.Fatalf("codex ACPArgs=%v, want nil until ACP launch args are configured", codex.ACPArgs)
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
	for _, want := range []string{"--print-timeout", "30m", "--dangerously-skip-permissions", "--model", "Gemini 3.6 Flash (High)", "--add-dir", "{root}", "--print", "{prompt}"} {
		if !containsStringForTest(agy.HeadlessArgs, want) {
			t.Fatalf("agy HeadlessArgs missing %q: %v", want, agy.HeadlessArgs)
		}
	}
	if agy.PromptMode != PromptArg {
		t.Fatalf("agy PromptMode=%q, want %q", agy.PromptMode, PromptArg)
	}
	if agy.Model != "Gemini 3.6 Flash (High)" || agy.Reasoning != CLIDefault {
		t.Fatalf("agy model/reasoning=%q/%q, want Gemini 3.6 Flash (High)/cli-default", agy.Model, agy.Reasoning)
	}

	claude := findSpecForTest(specs, "claude")
	if claude.Model != "claude-opus-4-8[1m]" || claude.Reasoning != "max" {
		t.Fatalf("claude model/reasoning=%q/%q, want claude-opus-4-8[1m]/max", claude.Model, claude.Reasoning)
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

func sameStringsForTest(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// AF-3 (review MINOR, kimi-1): kimi and opencode now live in DefaultSpecs, so
// mergeACPCatalog takes the MERGE branch for them instead of the append branch.
// Nothing previously asserted that `kimi acp` / `opencode acp` survive as an
// alternative launch mode — the constraint 00-prompt names — nor that the merge
// leaves their headless launch intact.
func TestPromotedAdaptersKeepACPAsAlternative(t *testing.T) {
	byID := map[string]Spec{}
	for _, s := range DefaultSpecs() {
		byID[s.ID] = s
	}
	for _, id := range []string{"kimi", "opencode"} {
		s, ok := byID[id]
		if !ok {
			t.Fatalf("missing built-in %q", id)
		}
		if got := strings.Join(s.ACPArgs, " "); got != "acp" {
			t.Errorf("%s: ACPArgs=%v want [acp] — the ACP launch mode was lost in the merge", id, s.ACPArgs)
		}
		// The merge must not demote the headless promotion this idea made.
		if s.LaunchMode != LaunchHeadless {
			t.Errorf("%s: launch mode=%q want %q — merge overwrote the headless promotion", id, s.LaunchMode, LaunchHeadless)
		}
		if len(s.HeadlessArgs) == 0 {
			t.Errorf("%s: HeadlessArgs emptied by the merge", id)
		}
		if !s.AutonomousWrite.Declared() {
			t.Errorf("%s: AutonomousWrite cleared by the merge", id)
		}
	}
}
