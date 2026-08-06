package agents

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type Spec struct {
	ID string
	// AdapterID is the vendor/family adapter used for launch + discovery dispatch
	// (env cleaning, isolated home, per-CLI invocation quirks). It is distinct from
	// ID, which after participant resolution is the stable roster identity used for
	// artifact paths and signoffs (idea: composite-agent-naming-and-roster-reinit).
	// Empty means "same as ID" — specs that predate the split keep working.
	AdapterID             string
	Commands              []string
	VersionArgs           []string
	LaunchMode            string
	HeadlessMode          string
	HeadlessArgs          []string
	InteractiveMode       string
	InteractiveCommand    string
	InteractiveArgs       []string
	InteractivePromptMode string
	InteractiveInvoke     string
	InteractiveTimeoutMS  int
	InteractivePollMS     int
	InteractiveNotes      string
	PromptMode            PromptMode
	SandboxMode           string
	ApprovalPolicy        string
	Model                 string
	// ModelLabel is the human model name used to DERIVE the composite display
	// name's model section (e.g. "Opus 4.8 1m" -> "opus-4.8-1m"). Empty falls back
	// to Model. Never an identity; purely for rendering.
	ModelLabel string
	Reasoning  string
	Profile    string
	Speed      string
	TimeoutMS  int
	// Supervision windows (runner-hardening-kindly D2). 0 = use the default
	// (first-event 120s, stall 30m clamped under timeout_ms, heartbeat 60s);
	// negative = guard explicitly disabled (the TOML layer maps an explicit 0
	// override to -1).
	FirstEventTimeoutMS int
	StallTimeoutMS      int
	HeartbeatMS         int
	IsolateHome         bool
	IsolatedHomeEnv     map[string]string
	ExternalBackend     string
	Telemetry           string
	Notes               string
	// BuffersStdout declares that the CLI buffers ALL stdout until process exit
	// (e.g. agy --print), so a silent transcript is expected, not a hang. The TUI
	// uses it for the buffered-agent placeholder hint.
	BuffersStdout bool
	// AutonomousWrite declares the CLI's non-interactive auto-approve mode so a
	// participant can write its own artifact without a blocking permission prompt
	// (idea composite-agent-naming-and-roster-reinit, component C). There is no
	// common flag across vendors, so each spec names its own. Scope is "workspace"
	// ONLY where a real sandbox is enforced (honesty rule); Declared() reports the
	// mode, Confined() the sandbox. Secret redaction is orthogonal.
	AutonomousWrite AutonomousWrite
	Sources         map[string]string
	// ACPArgs are the launch flags that put an ACP-capable CLI into ACP mode
	// (e.g. ["--experimental-acp"] for claude, ["acp"] for goose, ["--acp"] for qwen).
	// When LaunchMode == LaunchACP, the runner spawns Commands[0] with ACPArgs
	// and speaks JSON-RPC 2.0 over NDJSON on stdio instead of a one-shot text run.
	ACPArgs []string
}

// Adapter returns the vendor/family adapter id for launch + discovery dispatch.
// It falls back to ID so specs that predate the roster-ID/adapter split (where ID
// already IS the family) keep working unchanged.
func (s Spec) Adapter() string {
	if s.AdapterID != "" {
		return s.AdapterID
	}
	return s.ID
}

// AutonomousWrite is the per-agent auto-approve declaration (component C). Mode is
// the vendor's mode name (e.g. "bypassPermissions", "yolo"); Args are the flags that
// enable it; Scope names the DEMONSTRABLE confinement and is set to "workspace" ONLY
// where the CLI enforces a real workspace sandbox (codex `--sandbox workspace-write`).
// For CLIs whose only confinement is a grant/bypass flag + cwd (claude `--add-dir`,
// hermes `--yolo`), Scope is left EMPTY — the honesty rule (consensus §C, review
// CRITICAL codex-1): never falsely assert workspace confinement that is not enforced.
type AutonomousWrite struct {
	Mode  string
	Args  []string
	Scope string
}

// Declared reports that an autonomous write mode is configured. It does NOT assert
// OS confinement — read Scope for that (only "workspace" is a demonstrated sandbox).
func (a AutonomousWrite) Declared() bool {
	return strings.TrimSpace(a.Mode) != ""
}

// Confined reports a demonstrated workspace sandbox (honesty signal, distinct from
// merely being autonomous).
func (a AutonomousWrite) Confined() bool {
	return a.Scope == "workspace"
}

type PromptMode string

const (
	PromptStdin PromptMode = "stdin"
	PromptArg   PromptMode = "arg"
)

const (
	CLIDefault               = "cli-default"
	DefaultSpeed             = "balanced"
	DefaultTimeoutMS         = 1_800_000
	DefaultInteractivePollMS = 2_000
	ExternalHosted           = "hosted"
	ExternalLocal            = "local"
	ExternalUnknown          = "unknown"
	SourceBuiltIn            = "built-in"
	SourceDiscovered         = "discovered"
)

const (
	LaunchHeadless    = "headless"
	LaunchInteractive = "interactive"
	LaunchManual      = "manual"
	LaunchACP         = "acp"

	InteractivePromptNone = "none"
	InteractivePromptFile = "file"
	InteractivePromptArg  = "arg"

	InteractiveInvokePrintOnly = "print-only"
	InteractiveInvokeSpawnTTY  = "spawn-tty"
)

type Discovery struct {
	Spec
	Path      string
	Found     bool
	Version   string
	Error     string
	ProbeTime time.Duration
}

func DefaultSpecs() []Spec {
	return mergeACPCatalog(defaultBuiltinSpecs(), ACPCatalog())
}

func defaultBuiltinSpecs() []Spec {
	return []Spec{
		withBuiltinSources(Spec{
			ID:                    "codex",
			Commands:              []string{"codex"},
			VersionArgs:           []string{"--version"},
			LaunchMode:            LaunchHeadless,
			HeadlessMode:          "codex exec --skip-git-repo-check -",
			HeadlessArgs:          []string{"exec", "--skip-git-repo-check", "--cd", "{root}", "--sandbox", "workspace-write", "-c", "approval_policy=\"never\"", "-"},
			InteractivePromptMode: InteractivePromptNone,
			InteractiveInvoke:     InteractiveInvokePrintOnly,
			InteractivePollMS:     DefaultInteractivePollMS,
			PromptMode:            PromptStdin,
			SandboxMode:           "workspace-write",
			ApprovalPolicy:        "never",
			Model:                 CLIDefault,
			Reasoning:             CLIDefault,
			Profile:               CLIDefault,
			Speed:                 DefaultSpeed,
			TimeoutMS:             DefaultTimeoutMS,
			ExternalBackend:       ExternalHosted,
			Telemetry:             "json events when --json is available",
			// Autonomous writes confined by the workspace-write sandbox (no full-fs escalation).
			AutonomousWrite: AutonomousWrite{Mode: "approval_policy=never", Args: []string{"-c", "approval_policy=\"never\""}, Scope: "workspace"},
		}),
		withBuiltinSources(Spec{
			ID:                    "claude",
			Commands:              []string{"claude"},
			VersionArgs:           []string{"--version"},
			LaunchMode:            LaunchHeadless,
			HeadlessMode:          "claude --print",
			HeadlessArgs:          []string{"-p", "--model", "claude-opus-4-8[1m]", "--effort", "max", "--output-format", "text", "--permission-mode", "bypassPermissions", "--add-dir", "{root}"},
			InteractivePromptMode: InteractivePromptNone,
			InteractiveInvoke:     InteractiveInvokePrintOnly,
			InteractivePollMS:     DefaultInteractivePollMS,
			PromptMode:            PromptStdin,
			SandboxMode:           CLIDefault,
			ApprovalPolicy:        "bypassPermissions",
			Model:                 "claude-opus-4-8[1m]",
			ModelLabel:            "Opus 4.8 1m",
			Reasoning:             "max",
			Profile:               CLIDefault,
			Speed:                 DefaultSpeed,
			TimeoutMS:             DefaultTimeoutMS,
			ExternalBackend:       ExternalHosted,
			Telemetry:             "stream-json or final text depending on flags",
			// Full autonomous writes, confined to the workspace via --add-dir {root}.
			AutonomousWrite: AutonomousWrite{Mode: "bypassPermissions", Args: []string{"--permission-mode", "bypassPermissions"}, Scope: ""},
		}),
		withBuiltinSources(Spec{
			ID:           "agy",
			Commands:     []string{"agy"},
			VersionArgs:  []string{"--version"},
			LaunchMode:   LaunchHeadless,
			HeadlessMode: "agy --print",
			// Keep {prompt} immediately after --print and last; agy treats --print as a value-taking flag.
			HeadlessArgs:          []string{"--print-timeout", "30m", "--dangerously-skip-permissions", "--model", "Gemini 3.6 Flash (High)", "--add-dir", "{root}", "--print", "{prompt}"},
			InteractivePromptMode: InteractivePromptNone,
			InteractiveInvoke:     InteractiveInvokePrintOnly,
			InteractivePollMS:     DefaultInteractivePollMS,
			PromptMode:            PromptArg,
			SandboxMode:           CLIDefault,
			ApprovalPolicy:        "dangerously-skip-permissions",
			Model:                 "Gemini 3.6 Flash (High)",
			Reasoning:             CLIDefault,
			Profile:               CLIDefault,
			Speed:                 DefaultSpeed,
			TimeoutMS:             DefaultTimeoutMS,
			ExternalBackend:       ExternalHosted,
			Telemetry:             "unknown",
			BuffersStdout:         true, // agy --print emits nothing until exit
			Notes:                 "Antigravity CLI (active Gemini-family participant); agy 1.1.8 exposes --model. Strongest Gemini model as of 2026-07-29: Gemini 3.6 Flash (High); see `agy models` (3.6/3.5 Flash Low/Med/High, 3.1 Pro Low/High)",
			AutonomousWrite:       AutonomousWrite{Mode: "dangerously-skip-permissions", Args: []string{"--dangerously-skip-permissions"}, Scope: ""},
		}),
		withBuiltinSources(Spec{
			ID:                    "gemini",
			Commands:              []string{"gemini"},
			VersionArgs:           []string{"--version"},
			LaunchMode:            LaunchHeadless,
			HeadlessMode:          "gemini --prompt ... --output-format json",
			HeadlessArgs:          []string{"--prompt", "Follow the Parley Deck participant instructions provided on stdin.", "--skip-trust", "--approval-mode", "auto_edit", "--output-format", "text"},
			InteractivePromptMode: InteractivePromptNone,
			InteractiveInvoke:     InteractiveInvokePrintOnly,
			InteractivePollMS:     DefaultInteractivePollMS,
			PromptMode:            PromptStdin,
			SandboxMode:           CLIDefault,
			ApprovalPolicy:        "auto_edit",
			Model:                 CLIDefault,
			Reasoning:             CLIDefault,
			Profile:               CLIDefault,
			Speed:                 DefaultSpeed,
			TimeoutMS:             DefaultTimeoutMS,
			IsolateHome:           true,
			IsolatedHomeEnv:       map[string]string{"GEMINI_CLI_HOME": "{tempdir}"},
			ExternalBackend:       ExternalHosted,
			Telemetry:             "json stats when output-format json succeeds",
			Notes:                 "DEPRECATED legacy Gemini CLI support; prefer agy. Uses isolated GEMINI_CLI_HOME for oauth-personal profiles that hang",
		}),
		withBuiltinSources(Spec{
			ID:                    "hermes",
			Commands:              []string{"hermes", "hermes-agent", "hermesagent"},
			VersionArgs:           []string{"--version"},
			LaunchMode:            LaunchHeadless,
			HeadlessMode:          "hermes --yolo --oneshot",
			HeadlessArgs:          []string{"--yolo", "--oneshot", "{prompt}", "--model", "xai/grok-4.3", "--accept-hooks"},
			InteractivePromptMode: InteractivePromptNone,
			InteractiveInvoke:     InteractiveInvokePrintOnly,
			InteractivePollMS:     DefaultInteractivePollMS,
			PromptMode:            PromptArg,
			SandboxMode:           CLIDefault,
			ApprovalPolicy:        "yolo",
			Model:                 "xai/grok-4.3",
			Reasoning:             CLIDefault,
			Profile:               CLIDefault,
			Speed:                 DefaultSpeed,
			TimeoutMS:             DefaultTimeoutMS,
			IsolateHome:           true,
			IsolatedHomeEnv:       map[string]string{"HERMES_HOME": "{tempdir}"},
			ExternalBackend:       ExternalHosted,
			Telemetry:             "unknown",
			Notes:                 "first supported command found on PATH is used; uses isolated HERMES_HOME for writable logs",
			AutonomousWrite:       AutonomousWrite{Mode: "yolo", Args: []string{"--yolo"}, Scope: ""},
		}),
		withBuiltinSources(Spec{
			ID:           "kimi",
			Commands:     []string{"kimi"},
			VersionArgs:  []string{"--version"},
			LaunchMode:   LaunchHeadless,
			HeadlessMode: "kimi -p",
			// `-p/--prompt` IS kimi's non-interactive mode AND its auto-approve path: in print
			// mode the default permission policy applies tool calls without asking (probed
			// 2026-08-06 on 0.33.0 — wrote its file, exit 0). kimi also ships `--auto` ("fully
			// autonomous") and `-y/--yolo`, but BOTH are rejected alongside `-p`:
			// `kimi --auto -p …` exits 1 with "Cannot combine --prompt with --auto." So `-p` is
			// the only autonomous headless shape, and it is what we declare.
			HeadlessArgs:          []string{"-p", "{prompt}"},
			InteractivePromptMode: InteractivePromptNone,
			InteractiveInvoke:     InteractiveInvokePrintOnly,
			InteractivePollMS:     DefaultInteractivePollMS,
			PromptMode:            PromptArg,
			SandboxMode:           CLIDefault,
			ApprovalPolicy:        "print-mode default",
			// Model and effort come from ~/.kimi-code/config.toml (default_model +
			// [thinking].effort). A per-invocation `-m/--model` exists; no effort flag does.
			Model:           CLIDefault,
			Reasoning:       CLIDefault,
			Profile:         CLIDefault,
			Speed:           DefaultSpeed,
			TimeoutMS:       DefaultTimeoutMS,
			ExternalBackend: ExternalHosted,
			Telemetry:       "final text on stdout",
			Notes:           "Kimi Code. Headless is `kimi -p <prompt>`; --auto/--yolo cannot combine with -p. ACP remains available as an alternative launch mode via `kimi acp`.",
			// Scope is deliberately EMPTY: print mode confines nothing at the OS level, it only
			// auto-approves. Only codex --sandbox workspace-write earns Scope "workspace".
			AutonomousWrite: AutonomousWrite{Mode: "prompt", Args: []string{"-p"}, Scope: ""},
		}),
		withBuiltinSources(Spec{
			ID:           "opencode",
			Commands:     []string{"opencode"},
			VersionArgs:  []string{"--version"},
			LaunchMode:   LaunchHeadless,
			HeadlessMode: "opencode run --auto",
			// `run` alone already writes in-workspace unattended (probed 2026-08-06 on 1.18.13,
			// exit 0, file written). `--auto` is nevertheless passed EXPLICITLY: every other
			// built-in adapter carries its auto-approve flag in HeadlessArgs, and an implicit
			// default is exactly the thing a vendor may change between versions. Declaring it
			// keeps AutonomousWrite.Args a subset of HeadlessArgs, as for claude/codex/hermes.
			HeadlessArgs:          []string{"run", "--auto", "{prompt}"},
			InteractivePromptMode: InteractivePromptNone,
			InteractiveInvoke:     InteractiveInvokePrintOnly,
			InteractivePollMS:     DefaultInteractivePollMS,
			PromptMode:            PromptArg,
			SandboxMode:           CLIDefault,
			ApprovalPolicy:        "auto",
			// Model comes from ~/.config/opencode/opencode.jsonc ("model"); `-m provider/model`
			// overrides per invocation. Reasoning effort is `--variant <level>`; the level set
			// cannot be enumerated from the CLI, so it is left to config.
			Model:           CLIDefault,
			Reasoning:       CLIDefault,
			Profile:         CLIDefault,
			Speed:           DefaultSpeed,
			TimeoutMS:       DefaultTimeoutMS,
			ExternalBackend: ExternalHosted,
			Telemetry:       "final text on stdout",
			Notes:           "OpenCode. Headless is `opencode run --auto <prompt>`; the message is argv, not stdin. Per-invocation model is `-m provider/model`, reasoning effort is `--variant`. ACP remains available via `opencode acp`.",
			// Scope EMPTY: --auto is a permission grant, not an enforced sandbox.
			AutonomousWrite: AutonomousWrite{Mode: "auto", Args: []string{"--auto"}, Scope: ""},
		}),
	}
}

func withBuiltinSources(spec Spec) Spec {
	spec.Sources = map[string]string{}
	for _, field := range []string{
		"id",
		"commands",
		"version_args",
		"launch_mode",
		"headless_mode",
		"headless_args",
		"interactive_mode",
		"interactive_command",
		"interactive_args",
		"interactive_prompt_mode",
		"interactive_invoke",
		"interactive_timeout_ms",
		"interactive_poll_ms",
		"interactive_notes",
		"prompt_mode",
		"sandbox_mode",
		"approval_policy",
		"model",
		"model_label",
		"reasoning",
		"profile",
		"speed",
		"timeout_ms",
		"isolate_home",
		"isolated_home_env",
		"external_backend",
		"telemetry",
		"notes",
		"autonomous_write",
		"acp_args",
	} {
		spec.Sources[field] = SourceBuiltIn
	}
	return spec
}

func Discover(ctx context.Context, specs []Spec) []Discovery {
	results := make([]Discovery, 0, len(specs))
	for _, spec := range specs {
		result := Discovery{Spec: spec}
		for _, command := range spec.Commands {
			path, err := exec.LookPath(command)
			if err == nil {
				result.Path = path
				result.Found = true
				break
			}
		}
		if result.Found && len(spec.VersionArgs) > 0 {
			start := time.Now()
			version, err := probeVersion(ctx, result.Path, spec.VersionArgs)
			result.ProbeTime = time.Since(start)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Version = version
			}
		}
		results = append(results, result)
	}
	return results
}

func PrintDiscovery(w io.Writer, results []Discovery) {
	fmt.Fprintln(w, "AGENT    INSTALLED  VERSION                 HEADLESS")
	for _, result := range results {
		installed := "no"
		version := "-"
		if result.Found {
			installed = "yes"
		}
		if result.Version != "" {
			version = result.Version
		} else if result.Error != "" {
			version = "error"
		}
		fmt.Fprintf(w, "%-8s %-9s %-23s %s\n", result.ID, installed, truncate(version, 23), result.HeadlessMode)
		if result.Notes != "" {
			fmt.Fprintf(w, "  note: %s\n", result.Notes)
		}
		if result.Telemetry != "" {
			fmt.Fprintf(w, "  telemetry: %s\n", result.Telemetry)
		}
		if result.Error != "" {
			fmt.Fprintf(w, "  probe error: %s\n", result.Error)
		}
	}
}

func PrintRuntimeMatrix(w io.Writer, results []Discovery) {
	fmt.Fprintln(w, "AGENT    INSTALLED  VERSION                 LAUNCH       HEADLESS  SANDBOX          APPROVAL     MODEL        TIMEOUT  HOME  AUTO  BACKEND")
	for _, result := range results {
		installed := "no"
		version := "-"
		headless := "missing"
		launchMode := LaunchModeOrDefault(result.LaunchMode)
		if result.Found {
			installed = "yes"
		}
		if len(result.HeadlessArgs) > 0 {
			headless = "configured"
		}
		if result.Version != "" {
			version = result.Version
		} else if result.Error != "" {
			version = "error"
		}
		home := "no"
		if result.IsolateHome {
			home = "yes"
		}
		auto := "no"
		if result.AutonomousWrite.Declared() {
			auto = "yes"
		}
		fmt.Fprintf(
			w,
			"%-8s %-9s %-23s %-12s %-9s %-16s %-12s %-12s %-8d %-5s %-5s %s\n",
			result.ID,
			installed,
			truncate(version, 23),
			launchMode,
			headless,
			valueOrDefault(result.SandboxMode),
			valueOrDefault(result.ApprovalPolicy),
			valueOrDefault(result.Model),
			timeoutOrDefault(result.TimeoutMS),
			home,
			auto,
			valueOrDefault(result.ExternalBackend),
		)
		if len(result.Sources) > 0 {
			fmt.Fprintf(w, "  sources: sandbox=%s approval=%s model=%s timeout=%s\n",
				sourceFor(result.Sources, "sandbox_mode"),
				sourceFor(result.Sources, "approval_policy"),
				sourceFor(result.Sources, "model"),
				sourceFor(result.Sources, "timeout_ms"),
			)
		}
		if result.HeadlessMode != "" {
			fmt.Fprintf(w, "  headless: %s\n", result.HeadlessMode)
		}
		if result.ACPArgs != nil {
			command := CLIDefault
			if len(result.Commands) > 0 {
				command = result.Commands[0]
			}
			fmt.Fprintf(w, "  acp: %s %s\n", command, quoteArgs(result.ACPArgs))
		}
		if launchMode == LaunchInteractive || launchMode == LaunchManual || result.InteractiveMode != "" || len(result.InteractiveArgs) > 0 {
			fmt.Fprintf(w, "  interactive: %s %s prompt=%s invoke=%s\n",
				valueOrDefault(InteractiveCommandOrDefault(result.Spec)),
				quoteArgs(result.InteractiveArgs),
				InteractivePromptModeOrDefault(result.InteractivePromptMode),
				InteractiveInvokeOrDefault(result.InteractiveInvoke),
			)
		}
		if result.Notes != "" {
			fmt.Fprintf(w, "  note: %s\n", result.Notes)
		}
		if result.Error != "" {
			fmt.Fprintf(w, "  probe error: %s\n", result.Error)
		}
	}
}

func LaunchModeOrDefault(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return LaunchHeadless
	}
	return value
}

func InteractivePromptModeOrDefault(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return InteractivePromptNone
	}
	return value
}

func InteractiveInvokeOrDefault(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return InteractiveInvokePrintOnly
	}
	return value
}

func InteractiveTimeoutMSOrDefault(value int) int {
	if value > 0 {
		return value
	}
	return DefaultTimeoutMS
}

func InteractivePollMSOrDefault(value int) int {
	if value > 0 {
		return value
	}
	return DefaultInteractivePollMS
}

func InteractiveCommandOrDefault(spec Spec) string {
	if strings.TrimSpace(spec.InteractiveCommand) != "" {
		return spec.InteractiveCommand
	}
	if len(spec.Commands) > 0 {
		return spec.Commands[0]
	}
	return ""
}

func quoteArgs(args []string) string {
	if len(args) == 0 {
		return "[]"
	}
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = fmt.Sprintf("%q", arg)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func probeVersion(parent context.Context, path string, args []string) (string, error) {
	ctx, cancel := context.WithTimeout(parent, 4*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, path, args...)
	output, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", err
	}
	return bestVersionLine(string(output)), nil
}

func truncate(value string, width int) string {
	if len(value) <= width {
		return value
	}
	if width <= 1 {
		return value[:width]
	}
	return value[:width-1] + "…"
}

func bestVersionLine(output string) string {
	lines := strings.Split(output, "\n")
	var fallback string
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if fallback == "" {
			fallback = line
		}
		if !strings.HasPrefix(strings.ToLower(line), "warning:") {
			return line
		}
	}
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			return line
		}
	}
	return ""
}

func valueOrDefault(value string) string {
	if strings.TrimSpace(value) == "" {
		return CLIDefault
	}
	return value
}

func timeoutOrDefault(timeoutMS int) int {
	if timeoutMS <= 0 {
		return DefaultTimeoutMS
	}
	return timeoutMS
}

func sourceFor(sources map[string]string, field string) string {
	if value := strings.TrimSpace(sources[field]); value != "" {
		return value
	}
	return SourceDiscovered
}
