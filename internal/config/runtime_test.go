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
	if strings.Join(codex.HeadlessArgs, " ") != "exec --cd {root} --sandbox workspace-write -c approval_policy=\"on-failure\" -" {
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
