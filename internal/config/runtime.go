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

// EnvParleyHome overrides the user-global config dir (default ~/.parley). It
// lets tests point the central-default layer at a scratch directory.
const EnvParleyHome = "PARLEY_HOME"

type fileConfig struct {
	Defaults *globalDefaults           `toml:"defaults"`
	Agents   map[string]agentOverride  `toml:"agents"`
	Rosters  map[string]rosterOverride `toml:"rosters"`
}

// rosterOverride is one [rosters.<slug>] block — a named participant preset
// (named-roster-presets). Members are §2 canonical roster IDs (e.g. "claude-1").
type rosterOverride struct {
	Participants []string `toml:"participants"`
}

// globalDefaults is the optional [defaults] block of a layered config file —
// non-agent policy knobs that apply project-wide.
type globalDefaults struct {
	Speed              string            `toml:"speed"`
	PingTier           string            `toml:"ping_tier"`
	PreferredTransport string            `toml:"preferred_transport"`
	RosterChangePolicy string            `toml:"roster_change_policy"`
	Timeouts           *timeoutsBlock    `toml:"timeouts"`
	Loop               *loopBlock        `toml:"loop"`
	TrackRosters       map[string]string `toml:"track_rosters"`
}

type timeoutsBlock struct {
	SignoffMS       int `toml:"signoff_ms"`
	RoundMS         int `toml:"round_ms"`
	ReviewMS        int `toml:"review_ms"`
	DeepReasoningMS int `toml:"deep_reasoning_ms"`
}

// loopBlock is the [defaults.loop] policy: explicit auto-drive loop ceilings (LE-5).
// A breach escalates and halts; it never marks an idea complete. 0 = unlimited.
// Fields are pointers so a deliberate `= 0` at a higher layer overrides a lower
// layer's seeded value (presence-aware merge); absence (nil) falls through (review
// fix F-T2-1).
type loopBlock struct {
	MaxDriverSteps *int     `toml:"max_driver_steps"`
	MaxWallClockMS *int     `toml:"max_wall_clock_ms"`
	MaxCostUSD     *float64 `toml:"max_cost_usd"`
}

// CentralDefaults are the merged [defaults] policy knobs across the layered
// config files (central ~/.parley/agents.toml first, then the project deck).
// A zero-value field means "not set" — the consumer supplies its own fallback.
type CentralDefaults struct {
	Speed              string
	PingTier           string
	PreferredTransport string
	RosterChangePolicy string
	SignoffMS          int
	RoundMS            int
	ReviewMS           int
	DeepReasoningMS    int
	// Loop ceilings (LE-5); 0 = unlimited.
	MaxDriverSteps int
	MaxWallClockMS int
	MaxCostUSD     float64
}

type configLayer struct {
	path     string
	source   string
	optional bool
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
	ModelLabel            string            `toml:"model_label"`
	Reasoning             string            `toml:"reasoning"`
	Profile               string            `toml:"profile"`
	Speed                 string            `toml:"speed"`
	TimeoutMS             int               `toml:"timeout_ms"`
	FirstEventTimeoutMS   *int              `toml:"first_event_timeout_ms"`
	StallTimeoutMS        *int              `toml:"stall_timeout_ms"`
	HeartbeatMS           *int              `toml:"heartbeat_ms"`
	IsolateHome           *bool             `toml:"isolate_home"`
	BuffersStdout         *bool             `toml:"buffers_stdout"`
	IsolatedHomeEnv       map[string]string `toml:"isolated_home_env"`
	ExternalBackend       string            `toml:"external_backend"`
	Telemetry             string            `toml:"telemetry"`
	Notes                 string            `toml:"notes"`
}

// configLayers lists the agent/defaults config files in low-to-high precedence:
// the user-global central default (~/.parley/agents.toml) seeds every project,
// the project deck files override it, and $PARLEY_HEADLESS_AGENT_CONFIG wins.
func configLayers(root string) []configLayer {
	deck := filepath.Join(root, protocol.DeckDir)
	layers := []configLayer{}
	if central := CentralAgentsPath(); central != "" {
		layers = append(layers, configLayer{path: central, source: "~/.parley/agents.toml", optional: true})
	}
	layers = append(layers,
		configLayer{path: filepath.Join(deck, "agents.toml"), source: "parley-deck/agents.toml", optional: true},
		configLayer{path: filepath.Join(deck, "agents.local.toml"), source: "parley-deck/agents.local.toml", optional: true},
	)
	if envPath := strings.TrimSpace(os.Getenv(EnvAgentConfig)); envPath != "" {
		if !filepath.IsAbs(envPath) {
			envPath = filepath.Join(root, envPath)
		}
		layers = append(layers, configLayer{path: envPath, source: EnvAgentConfig + ":" + envPath, optional: false})
	}
	return layers
}

func LoadAgentSpecs(root string) ([]agents.Spec, error) {
	specs := cloneSpecs(agents.DefaultSpecs())
	for _, item := range configLayers(root) {
		var err error
		specs, err = applyFile(root, specs, item.path, item.source, item.optional)
		if err != nil {
			return nil, err
		}
	}
	return specs, nil
}

// LoadDefaults merges the [defaults] policy block across the layered config
// files (central first, deck overrides). Missing files are skipped; later
// non-empty values win. Consumers (preflight ping tier, init transport) read
// the result and fall back to their own default on a zero value.
func LoadDefaults(root string) (CentralDefaults, error) {
	var out CentralDefaults
	for _, item := range configLayers(root) {
		data, err := os.ReadFile(item.path)
		if err != nil {
			if item.optional && errors.Is(err, os.ErrNotExist) {
				continue
			}
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		var cfg fileConfig
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return out, fmt.Errorf("%s: %w", item.path, err)
		}
		if cfg.Defaults != nil {
			mergeDefaults(&out, cfg.Defaults)
		}
	}
	return out, nil
}

func mergeDefaults(out *CentralDefaults, gd *globalDefaults) {
	if s := strings.TrimSpace(gd.Speed); s != "" {
		out.Speed = s
	}
	if s := strings.TrimSpace(gd.PingTier); s != "" {
		out.PingTier = s
	}
	if s := strings.TrimSpace(gd.PreferredTransport); s != "" {
		out.PreferredTransport = s
	}
	if s := strings.TrimSpace(gd.RosterChangePolicy); s != "" {
		out.RosterChangePolicy = s
	}
	if gd.Timeouts != nil {
		if gd.Timeouts.SignoffMS > 0 {
			out.SignoffMS = gd.Timeouts.SignoffMS
		}
		if gd.Timeouts.RoundMS > 0 {
			out.RoundMS = gd.Timeouts.RoundMS
		}
		if gd.Timeouts.ReviewMS > 0 {
			out.ReviewMS = gd.Timeouts.ReviewMS
		}
		if gd.Timeouts.DeepReasoningMS > 0 {
			out.DeepReasoningMS = gd.Timeouts.DeepReasoningMS
		}
	}
	if gd.Loop != nil {
		// Presence-aware (F-T2-1): a deliberate `= 0` overrides a lower layer's seed
		// to unlimited; absence (nil) falls through.
		if gd.Loop.MaxDriverSteps != nil {
			out.MaxDriverSteps = *gd.Loop.MaxDriverSteps
		}
		if gd.Loop.MaxWallClockMS != nil {
			out.MaxWallClockMS = *gd.Loop.MaxWallClockMS
		}
		if gd.Loop.MaxCostUSD != nil {
			out.MaxCostUSD = *gd.Loop.MaxCostUSD
		}
	}
}

// CentralHome returns the user-global Parley config directory (~/.parley),
// honoring the PARLEY_HOME override. It returns "" when the home directory
// cannot be resolved, so the central-default layer is simply skipped.
func CentralHome() string {
	if v := strings.TrimSpace(os.Getenv(EnvParleyHome)); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return ""
	}
	return filepath.Join(home, ".parley")
}

// CentralAgentsPath returns ~/.parley/agents.toml, or "" if home is unresolved.
func CentralAgentsPath() string {
	home := CentralHome()
	if home == "" {
		return ""
	}
	return filepath.Join(home, "agents.toml")
}

// EnsureCentralDefault writes a starter ~/.parley/agents.toml that lists the
// built-in agents with their default model and reasoning, if the file does not
// already exist. It is a no-op when the file is present or home is unresolved.
// Returns the ensured path ("" when skipped).
func EnsureCentralDefault() (string, error) {
	path := CentralAgentsPath()
	if path == "" {
		return "", nil
	}
	if _, err := os.Stat(path); err == nil {
		return path, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(centralDefaultTemplate()), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// centralDefaultTemplate renders a starter central config from the built-in
// specs: one [agents.<id>] block per agent with model + reasoning. Default
// reasoning should be the strongest level each agent supports; edit to taste.
func centralDefaultTemplate() string {
	var b strings.Builder
	b.WriteString("# ~/.parley/agents.toml — central per-user Parley Deck defaults.\n")
	b.WriteString("# Lists each agent's default model and reasoning/effort level. Used by\n")
	b.WriteString("# every project unless a deck overrides it in parley-deck/agents.toml.\n")
	b.WriteString("# Prefer an exact model id over a vendor \"latest\" alias, and keep the\n")
	b.WriteString("# strongest (highest) reasoning level each agent supports.\n\n")
	b.WriteString("# Project-wide policy defaults; a deck's parley-deck/agents.toml overrides them.\n")
	b.WriteString("[defaults]\n")
	b.WriteString("speed = \"fast\"                             # fast output at the SAME model+effort (Claude Code /fast), NOT a downgrade; a separate axis from reasoning. Use \"deep\" per idea for heavy work.\n")
	b.WriteString("ping_tier = \"hosted-pong\"                 # §9.0 roster liveness ping before each idea (or \"none\")\n")
	b.WriteString("preferred_transport = \"local-dir\"          # parley init default transport (local-dir|github-pr|gitlab-mr)\n")
	b.WriteString("roster_change_policy = \"confirm-breaking\"  # auto-add new agents; user confirms drops/breaking changes\n\n")
	b.WriteString("[defaults.timeouts]\n")
	b.WriteString("signoff_ms = 600000          # 10 min\n")
	b.WriteString("round_ms = 1200000           # 20 min\n")
	b.WriteString("review_ms = 1200000          # 20 min\n")
	b.WriteString("deep_reasoning_ms = 1200000  # 20 min\n\n")
	b.WriteString("[defaults.loop]              # auto-drive loop ceilings (LE-5); breach escalates, never completes. 0 = unlimited.\n")
	b.WriteString("max_driver_steps = 200       # total progress steps before escalation (generous safety net)\n")
	b.WriteString("max_wall_clock_ms = 7200000  # 2 h total run budget (distinct from the per-tick 30 min round deadline)\n")
	b.WriteString("max_cost_usd = 0             # 0 = unlimited; best-effort, telemetry-gated (LE-6)\n\n")
	for _, spec := range agents.DefaultSpecs() {
		model := spec.Model
		if strings.TrimSpace(model) == "" {
			model = agents.CLIDefault
		}
		reasoning := spec.Reasoning
		if strings.TrimSpace(reasoning) == "" {
			reasoning = agents.CLIDefault
		}
		b.WriteString(fmt.Sprintf("[agents.%s]\n", spec.ID))
		b.WriteString(fmt.Sprintf("model = %q\n", model))
		b.WriteString(fmt.Sprintf("reasoning = %q\n\n", reasoning))
	}
	return b.String()
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

	// Global [defaults].speed applies to every spec (overriding the built-in
	// per-agent default); a per-agent override in this same file still wins below.
	if cfg.Defaults != nil && strings.TrimSpace(cfg.Defaults.Speed) != "" {
		for i := range specs {
			specs[i].Speed = cfg.Defaults.Speed
			if specs[i].Sources == nil {
				specs[i].Sources = map[string]string{}
			}
			specs[i].Sources["speed"] = source + ":defaults"
		}
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
	if override.ModelLabel != "" {
		spec.ModelLabel = override.ModelLabel
		spec.Sources["model_label"] = source
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
	// Supervision knobs are pointer-typed so an explicit `0` (disable) is
	// distinguishable from "not set"; 0 maps to -1 (disabled) on the Spec.
	applySupervisionMS := func(value *int, field *int, key string) {
		if value == nil {
			return
		}
		if *value == 0 {
			*field = -1
		} else {
			*field = *value
		}
		spec.Sources[key] = source
	}
	applySupervisionMS(override.FirstEventTimeoutMS, &spec.FirstEventTimeoutMS, "first_event_timeout_ms")
	applySupervisionMS(override.StallTimeoutMS, &spec.StallTimeoutMS, "stall_timeout_ms")
	applySupervisionMS(override.HeartbeatMS, &spec.HeartbeatMS, "heartbeat_ms")
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
