package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

const EnvAgentConfig = "PARLEY_HEADLESS_AGENT_CONFIG"

type fileConfig struct {
	Agents map[string]agentOverride `toml:"agents"`
}

type agentOverride struct {
	Command               string            `toml:"command"`
	Path                  string            `toml:"path"`
	Commands              []string          `toml:"commands"`
	VersionArgs           []string          `toml:"version_args"`
	LaunchMode            string            `toml:"launch_mode"`
	HeadlessMode          string            `toml:"headless_mode"`
	HeadlessArgs          []string          `toml:"headless_args"`
	ACPArgs               []string          `toml:"acp_args"`
	InteractiveMode       string            `toml:"interactive_mode"`
	InteractiveCommand    string            `toml:"interactive_command"`
	InteractiveArgs       []string          `toml:"interactive_args"`
	InteractivePromptMode string            `toml:"interactive_prompt_mode"`
	InteractiveInvoke     string            `toml:"interactive_invoke"`
	InteractiveTimeoutMS  int               `toml:"interactive_timeout_ms"`
	InteractivePollMS     int               `toml:"interactive_poll_ms"`
	InteractiveNotes      string            `toml:"interactive_notes"`
	PromptMode            string            `toml:"prompt_mode"`
	SandboxMode           string            `toml:"sandbox_mode"`
	ApprovalPolicy        string            `toml:"approval_policy"`
	Model                 string            `toml:"model"`
	Reasoning             string            `toml:"reasoning"`
	Profile               string            `toml:"profile"`
	Speed                 string            `toml:"speed"`
	TimeoutMS             int               `toml:"timeout_ms"`
	IsolateHome           *bool             `toml:"isolate_home"`
	BuffersStdout         *bool             `toml:"buffers_stdout"`
	IsolatedHomeEnv       map[string]string `toml:"isolated_home_env"`
	ExternalBackend       string            `toml:"external_backend"`
	Telemetry             string            `toml:"telemetry"`
	Notes                 string            `toml:"notes"`
}

func LoadAgentSpecs(root string) ([]agents.Spec, error) {
	specs := cloneSpecs(agents.DefaultSpecs())
	deck := filepath.Join(root, protocol.DeckDir)

	for _, item := range []struct {
		path   string
		source string
	}{
		{path: filepath.Join(deck, "agents.toml"), source: "parley-deck/agents.toml"},
		{path: filepath.Join(deck, "agents.local.toml"), source: "parley-deck/agents.local.toml"},
	} {
		var err error
		specs, err = applyFile(root, specs, item.path, item.source, true)
		if err != nil {
			return nil, err
		}
	}

	if envPath := strings.TrimSpace(os.Getenv(EnvAgentConfig)); envPath != "" {
		if !filepath.IsAbs(envPath) {
			envPath = filepath.Join(root, envPath)
		}
		next, err := applyFile(root, specs, envPath, EnvAgentConfig+":"+envPath, false)
		if err != nil {
			return nil, err
		}
		specs = next
	}

	return specs, nil
}

func ExpandPlaceholders(value, root, tempdir string) string {
	replacements := map[string]string{
		"{root}":    root,
		"{deck}":    filepath.Join(root, protocol.DeckDir),
		"{tempdir}": tempdir,
	}
	for from, to := range replacements {
		value = strings.ReplaceAll(value, from, to)
	}
	return value
}

func applyFile(root string, specs []agents.Spec, path, source string, optional bool) ([]agents.Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if optional && errors.Is(err, os.ErrNotExist) {
			return specs, nil
		}
		return nil, err
	}

	var cfg fileConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(cfg.Agents) == 0 {
		return specs, nil
	}

	byID := map[string]int{}
	for i, spec := range specs {
		byID[spec.ID] = i
	}
	for id, override := range cfg.Agents {
		index, ok := byID[id]
		if !ok {
			specs = append(specs, agents.Spec{
				ID:                    id,
				PromptMode:            agents.PromptStdin,
				LaunchMode:            agents.LaunchHeadless,
				InteractivePromptMode: agents.InteractivePromptNone,
				InteractiveInvoke:     agents.InteractiveInvokePrintOnly,
				InteractivePollMS:     agents.DefaultInteractivePollMS,
				Model:                 agents.CLIDefault,
				Reasoning:             agents.CLIDefault,
				Profile:               agents.CLIDefault,
				Speed:                 agents.DefaultSpeed,
				TimeoutMS:             agents.DefaultTimeoutMS,
				ExternalBackend:       agents.ExternalUnknown,
				Sources:               configDefaultSources(source),
			})
			index = len(specs) - 1
			byID[id] = index
		}
		spec := specs[index]
		spec = applyOverride(root, spec, override, source)
		specs[index] = spec
	}

	return specs, nil
}

func configDefaultSources(source string) map[string]string {
	sources := map[string]string{}
	for _, field := range []string{
		"id",
		"launch_mode",
		"prompt_mode",
		"interactive_prompt_mode",
		"interactive_invoke",
		"interactive_poll_ms",
		"model",
		"reasoning",
		"profile",
		"speed",
		"timeout_ms",
		"external_backend",
	} {
		sources[field] = source + ":default"
	}
	return sources
}

func applyOverride(root string, spec agents.Spec, override agentOverride, source string) agents.Spec {
	tempdir := os.TempDir()
	if spec.Sources == nil {
		spec.Sources = map[string]string{}
	}

	if override.Path != "" {
		spec.Commands = []string{ExpandPlaceholders(override.Path, root, tempdir)}
		spec.Sources["commands"] = source
	} else if override.Command != "" {
		spec.Commands = []string{ExpandPlaceholders(override.Command, root, tempdir)}
		spec.Sources["commands"] = source
	} else if len(override.Commands) > 0 {
		spec.Commands = expandSlice(override.Commands, root, tempdir)
		spec.Sources["commands"] = source
	}
	if len(override.VersionArgs) > 0 {
		spec.VersionArgs = expandSlice(override.VersionArgs, root, tempdir)
		spec.Sources["version_args"] = source
	}
	if override.LaunchMode != "" {
		spec.LaunchMode = override.LaunchMode
		spec.Sources["launch_mode"] = source
	}
	if override.HeadlessMode != "" {
		spec.HeadlessMode = override.HeadlessMode
		spec.Sources["headless_mode"] = source
	}
	if len(override.HeadlessArgs) > 0 {
		spec.HeadlessArgs = expandSlice(override.HeadlessArgs, root, tempdir)
		spec.Sources["headless_args"] = source
	}
	if override.ACPArgs != nil {
		spec.ACPArgs = expandSlice(override.ACPArgs, root, tempdir)
		spec.Sources["acp_args"] = source
	}
	if override.InteractiveMode != "" {
		spec.InteractiveMode = override.InteractiveMode
		spec.Sources["interactive_mode"] = source
	}
	if override.InteractiveCommand != "" {
		spec.InteractiveCommand = ExpandPlaceholders(override.InteractiveCommand, root, tempdir)
		spec.Sources["interactive_command"] = source
	}
	if len(override.InteractiveArgs) > 0 {
		spec.InteractiveArgs = expandSlice(override.InteractiveArgs, root, tempdir)
		spec.Sources["interactive_args"] = source
	}
	if override.InteractivePromptMode != "" {
		spec.InteractivePromptMode = override.InteractivePromptMode
		spec.Sources["interactive_prompt_mode"] = source
	}
	if override.InteractiveInvoke != "" {
		spec.InteractiveInvoke = override.InteractiveInvoke
		spec.Sources["interactive_invoke"] = source
	}
	if override.InteractiveTimeoutMS > 0 {
		spec.InteractiveTimeoutMS = override.InteractiveTimeoutMS
		spec.Sources["interactive_timeout_ms"] = source
	}
	if override.InteractivePollMS > 0 {
		spec.InteractivePollMS = override.InteractivePollMS
		spec.Sources["interactive_poll_ms"] = source
	}
	if override.InteractiveNotes != "" {
		spec.InteractiveNotes = override.InteractiveNotes
		spec.Sources["interactive_notes"] = source
	}
	if override.PromptMode != "" {
		spec.PromptMode = agents.PromptMode(override.PromptMode)
		spec.Sources["prompt_mode"] = source
	}
	if override.SandboxMode != "" {
		spec.SandboxMode = override.SandboxMode
		spec.Sources["sandbox_mode"] = source
	}
	if override.ApprovalPolicy != "" {
		spec.ApprovalPolicy = override.ApprovalPolicy
		spec.Sources["approval_policy"] = source
	}
	if override.Model != "" {
		spec.Model = override.Model
		spec.Sources["model"] = source
	}
	if override.Reasoning != "" {
		spec.Reasoning = override.Reasoning
		spec.Sources["reasoning"] = source
	}
	if override.Profile != "" {
		spec.Profile = override.Profile
		spec.Sources["profile"] = source
	}
	if override.Speed != "" {
		spec.Speed = override.Speed
		spec.Sources["speed"] = source
	}
	if override.TimeoutMS > 0 {
		spec.TimeoutMS = override.TimeoutMS
		spec.Sources["timeout_ms"] = source
	}
	if override.IsolateHome != nil {
		spec.IsolateHome = *override.IsolateHome
		spec.Sources["isolate_home"] = source
	}
	if override.BuffersStdout != nil {
		spec.BuffersStdout = *override.BuffersStdout
		spec.Sources["buffers_stdout"] = source
	}
	if len(override.IsolatedHomeEnv) > 0 {
		spec.IsolatedHomeEnv = map[string]string{}
		for key, value := range override.IsolatedHomeEnv {
			spec.IsolatedHomeEnv[key] = ExpandPlaceholders(value, root, "{tempdir}")
		}
		spec.Sources["isolated_home_env"] = source
	}
	if override.ExternalBackend != "" {
		spec.ExternalBackend = override.ExternalBackend
		spec.Sources["external_backend"] = source
	}
	if override.Telemetry != "" {
		spec.Telemetry = override.Telemetry
		spec.Sources["telemetry"] = source
	}
	if override.Notes != "" {
		spec.Notes = override.Notes
		spec.Sources["notes"] = source
	}

	return spec
}

func expandSlice(values []string, root, tempdir string) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = ExpandPlaceholders(value, root, tempdir)
	}
	return out
}

func cloneSpecs(specs []agents.Spec) []agents.Spec {
	out := make([]agents.Spec, len(specs))
	for i, spec := range specs {
		out[i] = spec
		out[i].Commands = append([]string(nil), spec.Commands...)
		out[i].VersionArgs = append([]string(nil), spec.VersionArgs...)
		out[i].HeadlessArgs = append([]string(nil), spec.HeadlessArgs...)
		out[i].ACPArgs = cloneOptionalStringSlice(spec.ACPArgs)
		out[i].InteractiveArgs = append([]string(nil), spec.InteractiveArgs...)
		if spec.IsolatedHomeEnv != nil {
			out[i].IsolatedHomeEnv = map[string]string{}
			for key, value := range spec.IsolatedHomeEnv {
				out[i].IsolatedHomeEnv[key] = value
			}
		}
		if spec.Sources != nil {
			out[i].Sources = map[string]string{}
			for key, value := range spec.Sources {
				out[i].Sources[key] = value
			}
		}
	}
	return out
}

func cloneOptionalStringSlice(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string{}, values...)
}
