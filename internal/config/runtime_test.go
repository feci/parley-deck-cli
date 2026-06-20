package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

func TestLoadAgentSpecsLayersAndTracksSources(t *testing.T) {
	root := t.TempDir()
	t.Setenv(EnvParleyHome, filepath.Join(root, "no-central")) // isolate from real ~/.parley
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	project := filepath.Join(root, protocol.DeckDir, "agents.toml")
	local := filepath.Join(root, protocol.DeckDir, "agents.local.toml")
	env := filepath.Join(root, "env-agents.toml")

	if err := os.WriteFile(project, []byte(`
[agents.codex]
model = "project-model"
timeout_ms = 1000
sandbox_mode = "workspace-write"

[agents.extra]
command = "{root}/bin/extra"
headless_args = ["--prompt", "{deck}/prompt.md"]
external_backend = "local"

[agents.extra.isolated_home_env]
EXTRA_HOME = "{tempdir}/extra"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(local, []byte(`
[agents.codex]
model = "local-model"
approval_policy = "on-failure"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(env, []byte(`
[agents.codex]
model = "env-model"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvAgentConfig, env)

	specs, err := LoadAgentSpecs(root)
	if err != nil {
		t.Fatal(err)
	}
	codex := findSpec(t, specs, "codex")
	if codex.Model != "env-model" {
		t.Fatalf("codex model=%q, want env-model", codex.Model)
	}
	if codex.TimeoutMS != 1000 {
		t.Fatalf("timeout=%d, want 1000", codex.TimeoutMS)
	}
	if got := codex.Sources["model"]; !strings.HasPrefix(got, EnvAgentConfig+":") {
		t.Fatalf("model source=%q, want env source", got)
	}
	if got := codex.Sources["timeout_ms"]; got != "parley-deck/agents.toml" {
		t.Fatalf("timeout source=%q, want project source", got)
	}
	if got := codex.Sources["approval_policy"]; got != "parley-deck/agents.local.toml" {
		t.Fatalf("approval source=%q, want local source", got)
	}

	extra := findSpec(t, specs, "extra")
	if got, want := extra.Commands[0], filepath.Join(root, "bin", "extra"); got != want {
		t.Fatalf("extra command=%q, want %q", got, want)
	}
	if got, want := extra.HeadlessArgs[1], filepath.Join(root, protocol.DeckDir, "prompt.md"); got != want {
		t.Fatalf("headless arg=%q, want %q", got, want)
	}
	if extra.ExternalBackend != agents.ExternalLocal {
		t.Fatalf("external backend=%q, want local", extra.ExternalBackend)
	}
	if got := extra.Sources["model"]; got == agents.SourceDiscovered || got == "" {
		t.Fatalf("extra model source=%q, want config default source", got)
	}
	if got, want := extra.IsolatedHomeEnv["EXTRA_HOME"], "{tempdir}/extra"; got != want {
		t.Fatalf("isolated env=%q, want %q", got, want)
	}
}

func TestCodexBuiltInRuntimeDefaults(t *testing.T) {
	specs, err := LoadAgentSpecs(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	codex := findSpec(t, specs, "codex")
	if codex.SandboxMode != "workspace-write" {
		t.Fatalf("sandbox=%q", codex.SandboxMode)
	}
	if codex.ApprovalPolicy != "on-failure" {
		t.Fatalf("approval=%q", codex.ApprovalPolicy)
	}
	if strings.Join(codex.HeadlessArgs, " ") != "exec --skip-git-repo-check --cd {root} --sandbox workspace-write -c approval_policy=\"on-failure\" -" {
		t.Fatalf("headless args=%v", codex.HeadlessArgs)
	}
	if codex.Model != agents.CLIDefault || codex.Reasoning != agents.CLIDefault || codex.Profile != agents.CLIDefault {
		t.Fatalf("unexpected model/reasoning/profile: %+v", codex)
	}
}

func TestLoadAgentSpecsInteractiveLaunchFields(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	local := filepath.Join(root, protocol.DeckDir, "agents.local.toml")
	if err := os.WriteFile(local, []byte(`
[agents.claude]
launch_mode = "interactive"
interactive_mode = "claude tty"
interactive_command = "claude"
interactive_args = ["--resume", "{prompt_path}"]
interactive_prompt_mode = "file"
interactive_invoke = "print-only"
interactive_timeout_ms = 120000
interactive_poll_ms = 500
interactive_notes = "local note"
`), 0o644); err != nil {
		t.Fatal(err)
	}

	specs, err := LoadAgentSpecs(root)
	if err != nil {
		t.Fatal(err)
	}
	claude := findSpec(t, specs, "claude")
	if claude.LaunchMode != agents.LaunchInteractive {
		t.Fatalf("launch=%q", claude.LaunchMode)
	}
	if claude.InteractiveMode != "claude tty" || claude.InteractiveCommand != "claude" {
		t.Fatalf("interactive command fields: %+v", claude)
	}
	if got := strings.Join(claude.InteractiveArgs, " "); got != "--resume {prompt_path}" {
		t.Fatalf("interactive args=%q", got)
	}
	if claude.InteractivePromptMode != agents.InteractivePromptFile || claude.InteractiveInvoke != agents.InteractiveInvokePrintOnly {
		t.Fatalf("interactive mode/invoke: %+v", claude)
	}
	if claude.InteractiveTimeoutMS != 120000 || claude.InteractivePollMS != 500 {
		t.Fatalf("interactive timeouts: %+v", claude)
	}
	if claude.InteractiveNotes != "local note" {
		t.Fatalf("notes=%q", claude.InteractiveNotes)
	}
}

func TestLoadAgentSpecsACPArgs(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	local := filepath.Join(root, protocol.DeckDir, "agents.local.toml")
	if err := os.WriteFile(local, []byte(`
[agents.codex]
launch_mode = "acp"
acp_args = ["acp", "--stdio"]

[agents.vibe]
command = "vibe-acp"
launch_mode = "acp"
acp_args = []
`), 0o644); err != nil {
		t.Fatal(err)
	}

	specs, err := LoadAgentSpecs(root)
	if err != nil {
		t.Fatal(err)
	}
	codex := findSpec(t, specs, "codex")
	if codex.LaunchMode != agents.LaunchACP {
		t.Fatalf("codex launch=%q", codex.LaunchMode)
	}
	if got := strings.Join(codex.ACPArgs, " "); got != "acp --stdio" {
		t.Fatalf("codex acp args=%q", got)
	}
	if got := codex.Sources["acp_args"]; got != "parley-deck/agents.local.toml" {
		t.Fatalf("codex acp_args source=%q", got)
	}
	vibe := findSpec(t, specs, "vibe")
	if vibe.ACPArgs == nil || len(vibe.ACPArgs) != 0 {
		t.Fatalf("vibe ACPArgs=%v, want configured empty list", vibe.ACPArgs)
	}
}

func TestExpandPlaceholders(t *testing.T) {
	root := filepath.Join("tmp", "repo")
	temp := filepath.Join("tmp", "agent")
	got := ExpandPlaceholders("{root}/x:{deck}/y:{tempdir}/z:{prompt}", root, temp)
	want := filepath.Join("tmp", "repo", "x") + ":" +
		filepath.Join("tmp", "repo", protocol.DeckDir, "y") + ":" +
		filepath.Join("tmp", "agent", "z") + ":{prompt}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func findSpec(t *testing.T, specs []agents.Spec, id string) agents.Spec {
	t.Helper()
	for _, spec := range specs {
		if spec.ID == id {
			return spec
		}
	}
	t.Fatalf("missing spec %s", id)
	return agents.Spec{}
}

// Supervision knobs are tri-state: absent leaves the spec default, an explicit
// 0 maps to -1 (disabled), a positive value passes through.
func TestSupervisionKnobOverrides(t *testing.T) {
	root := t.TempDir()
	deck := filepath.Join(root, "parley-deck")
	if err := os.MkdirAll(deck, 0o755); err != nil {
		t.Fatal(err)
	}
	toml := `[agents.codex]
first_event_timeout_ms = 0
stall_timeout_ms = 900000
`
	if err := os.WriteFile(filepath.Join(deck, "agents.toml"), []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}
	specs, err := LoadAgentSpecs(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, spec := range specs {
		if spec.ID != "codex" {
			continue
		}
		if spec.FirstEventTimeoutMS != -1 {
			t.Fatalf("explicit 0 must map to -1 (disabled), got %d", spec.FirstEventTimeoutMS)
		}
		if spec.StallTimeoutMS != 900000 {
			t.Fatalf("stall_timeout_ms=%d, want 900000", spec.StallTimeoutMS)
		}
		if spec.HeartbeatMS != 0 {
			t.Fatalf("untouched heartbeat must stay 0 (default), got %d", spec.HeartbeatMS)
		}
		return
	}
	t.Fatal("codex spec not found")
}

func TestLoadAgentSpecsCentralDefault(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvParleyHome, home)
	if err := os.WriteFile(filepath.Join(home, "agents.toml"), []byte(`
[agents.claude]
model = "central-model"
reasoning = "central-reasoning"
`), 0o644); err != nil {
		t.Fatal(err)
	}

	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}

	// Central default applies when the deck does not override it.
	specs, err := LoadAgentSpecs(root)
	if err != nil {
		t.Fatal(err)
	}
	claude := findSpec(t, specs, "claude")
	if claude.Model != "central-model" {
		t.Fatalf("claude model=%q, want central-model", claude.Model)
	}
	if claude.Reasoning != "central-reasoning" {
		t.Fatalf("claude reasoning=%q, want central-reasoning", claude.Reasoning)
	}
	if got := claude.Sources["model"]; got != "~/.parley/agents.toml" {
		t.Fatalf("claude model source=%q, want ~/.parley/agents.toml", got)
	}

	// A deck override beats the central default; fields the deck leaves unset
	// fall through to the central value.
	project := filepath.Join(root, protocol.DeckDir, "agents.toml")
	if err := os.WriteFile(project, []byte(`
[agents.claude]
model = "project-model"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	specs, err = LoadAgentSpecs(root)
	if err != nil {
		t.Fatal(err)
	}
	claude = findSpec(t, specs, "claude")
	if claude.Model != "project-model" {
		t.Fatalf("claude model=%q, want project-model (deck overrides central)", claude.Model)
	}
	if claude.Reasoning != "central-reasoning" {
		t.Fatalf("claude reasoning=%q, want central-reasoning to survive partial deck override", claude.Reasoning)
	}
}

func TestLoadDefaultsMergesLayers(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvParleyHome, home)
	if err := os.WriteFile(filepath.Join(home, "agents.toml"), []byte(`
[defaults]
speed = "deep"
ping_tier = "hosted-pong"
preferred_transport = "github-pr"
roster_change_policy = "confirm-breaking"
[defaults.timeouts]
signoff_ms = 900000
round_ms = 1800000
`), 0o644); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}

	defs, err := LoadDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if defs.Speed != "deep" || defs.PingTier != "hosted-pong" ||
		defs.PreferredTransport != "github-pr" || defs.RosterChangePolicy != "confirm-breaking" {
		t.Fatalf("central defaults not loaded: %+v", defs)
	}
	if defs.SignoffMS != 900000 || defs.RoundMS != 1800000 {
		t.Fatalf("timeouts not loaded: %+v", defs)
	}

	// A deck overrides a subset; unset deck fields fall through to central.
	if err := os.WriteFile(filepath.Join(root, protocol.DeckDir, "agents.toml"), []byte(`
[defaults]
ping_tier = "none"
preferred_transport = "local-dir"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	defs, err = LoadDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if defs.PingTier != "none" || defs.PreferredTransport != "local-dir" {
		t.Fatalf("deck override failed: %+v", defs)
	}
	if defs.Speed != "deep" || defs.RosterChangePolicy != "confirm-breaking" {
		t.Fatalf("unset deck fields should fall through to central: %+v", defs)
	}
}

func TestGlobalSpeedAppliesToSpecs(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvParleyHome, home)
	if err := os.WriteFile(filepath.Join(home, "agents.toml"), []byte(`
[defaults]
speed = "deep"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	specs, err := LoadAgentSpecs(root)
	if err != nil {
		t.Fatal(err)
	}
	c := findSpec(t, specs, "claude")
	if c.Speed != "deep" {
		t.Fatalf("claude speed=%q, want deep (global default)", c.Speed)
	}
	if got := c.Sources["speed"]; got != "~/.parley/agents.toml:defaults" {
		t.Fatalf("speed source=%q, want ~/.parley/agents.toml:defaults", got)
	}
}

func TestEnsureCentralDefaultSeedsAndPreserves(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvParleyHome, home)

	path, err := EnsureCentralDefault()
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(home, "agents.toml"); path != want {
		t.Fatalf("path=%q, want %q", path, want)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"[agents.claude]", "model =", "reasoning ="} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("template missing %q:\n%s", want, data)
		}
	}

	// The generated template must be valid TOML that LoadAgentSpecs accepts.
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadAgentSpecs(root); err != nil {
		t.Fatalf("central template not loadable: %v", err)
	}

	// An existing central file is preserved, never overwritten.
	if err := os.WriteFile(path, []byte("[agents.claude]\nmodel = \"sentinel\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureCentralDefault(); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(path)
	if !strings.Contains(string(data), "sentinel") {
		t.Fatalf("EnsureCentralDefault overwrote an existing central file")
	}
}
