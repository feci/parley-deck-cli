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
	ID              string
	Commands        []string
	VersionArgs     []string
	HeadlessMode    string
	HeadlessArgs    []string
	PromptMode      PromptMode
	SandboxMode     string
	ApprovalPolicy  string
	Model           string
	Reasoning       string
	Profile         string
	Speed           string
	TimeoutMS       int
	IsolateHome     bool
	IsolatedHomeEnv map[string]string
	ExternalBackend string
	Telemetry       string
	Notes           string
	Sources         map[string]string
}

type PromptMode string

const (
	PromptStdin PromptMode = "stdin"
	PromptArg   PromptMode = "arg"
)

const (
	CLIDefault       = "cli-default"
	DefaultSpeed     = "balanced"
	DefaultTimeoutMS = 1_800_000
	ExternalHosted   = "hosted"
	ExternalLocal    = "local"
	ExternalUnknown  = "unknown"
	SourceBuiltIn    = "built-in"
	SourceDiscovered = "discovered"
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
	return []Spec{
		withBuiltinSources(Spec{
			ID:              "codex",
			Commands:        []string{"codex"},
			VersionArgs:     []string{"--version"},
			HeadlessMode:    "codex exec -",
			HeadlessArgs:    []string{"exec", "--cd", "{root}", "--sandbox", "workspace-write", "--ask-for-approval", "on-failure", "-"},
			PromptMode:      PromptStdin,
			SandboxMode:     "workspace-write",
			ApprovalPolicy:  "on-failure",
			Model:           CLIDefault,
			Reasoning:       CLIDefault,
			Profile:         CLIDefault,
			Speed:           DefaultSpeed,
			TimeoutMS:       DefaultTimeoutMS,
			ExternalBackend: ExternalHosted,
			Telemetry:       "json events when --json is available",
		}),
		withBuiltinSources(Spec{
			ID:              "claude",
			Commands:        []string{"claude"},
			VersionArgs:     []string{"--version"},
			HeadlessMode:    "claude --print",
			HeadlessArgs:    []string{"-p", "--output-format", "text", "--permission-mode", "acceptEdits", "--add-dir", "{root}"},
			PromptMode:      PromptStdin,
			SandboxMode:     CLIDefault,
			ApprovalPolicy:  "acceptEdits",
			Model:           CLIDefault,
			Reasoning:       CLIDefault,
			Profile:         CLIDefault,
			Speed:           DefaultSpeed,
			TimeoutMS:       DefaultTimeoutMS,
			ExternalBackend: ExternalHosted,
			Telemetry:       "stream-json or final text depending on flags",
		}),
		withBuiltinSources(Spec{
			ID:              "gemini",
			Commands:        []string{"gemini"},
			VersionArgs:     []string{"--version"},
			HeadlessMode:    "gemini --prompt ... --output-format json",
			HeadlessArgs:    []string{"--prompt", "Follow the Parley Deck participant instructions provided on stdin.", "--skip-trust", "--approval-mode", "auto_edit", "--output-format", "text"},
			PromptMode:      PromptStdin,
			SandboxMode:     CLIDefault,
			ApprovalPolicy:  "auto_edit",
			Model:           CLIDefault,
			Reasoning:       CLIDefault,
			Profile:         CLIDefault,
			Speed:           DefaultSpeed,
			TimeoutMS:       DefaultTimeoutMS,
			IsolateHome:     true,
			IsolatedHomeEnv: map[string]string{"GEMINI_CLI_HOME": "{tempdir}"},
			ExternalBackend: ExternalHosted,
			Telemetry:       "json stats when output-format json succeeds",
			Notes:           "uses isolated GEMINI_CLI_HOME for oauth-personal profiles that hang",
		}),
		withBuiltinSources(Spec{
			ID:              "hermes",
			Commands:        []string{"hermes", "hermes-agent", "hermesagent"},
			VersionArgs:     []string{"--version"},
			HeadlessMode:    "hermes --oneshot",
			HeadlessArgs:    []string{"--oneshot", "{prompt}", "--accept-hooks"},
			PromptMode:      PromptArg,
			SandboxMode:     CLIDefault,
			ApprovalPolicy:  "accept-hooks",
			Model:           CLIDefault,
			Reasoning:       CLIDefault,
			Profile:         CLIDefault,
			Speed:           DefaultSpeed,
			TimeoutMS:       DefaultTimeoutMS,
			IsolateHome:     true,
			IsolatedHomeEnv: map[string]string{"HERMES_HOME": "{tempdir}"},
			ExternalBackend: ExternalHosted,
			Telemetry:       "unknown",
			Notes:           "first supported command found on PATH is used; uses isolated HERMES_HOME for writable logs",
		}),
	}
}

func withBuiltinSources(spec Spec) Spec {
	spec.Sources = map[string]string{}
	for _, field := range []string{
		"id",
		"commands",
		"version_args",
		"headless_mode",
		"headless_args",
		"prompt_mode",
		"sandbox_mode",
		"approval_policy",
		"model",
		"reasoning",
		"profile",
		"speed",
		"timeout_ms",
		"isolate_home",
		"isolated_home_env",
		"external_backend",
		"telemetry",
		"notes",
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
	fmt.Fprintln(w, "AGENT    INSTALLED  VERSION                 HEADLESS  SANDBOX          APPROVAL     MODEL        TIMEOUT  HOME  BACKEND")
	for _, result := range results {
		installed := "no"
		version := "-"
		headless := "not-probed"
		if result.Found {
			installed = "yes"
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
		fmt.Fprintf(
			w,
			"%-8s %-9s %-23s %-9s %-16s %-12s %-12s %-8d %-5s %s\n",
			result.ID,
			installed,
			truncate(version, 23),
			headless,
			valueOrDefault(result.SandboxMode),
			valueOrDefault(result.ApprovalPolicy),
			valueOrDefault(result.Model),
			timeoutOrDefault(result.TimeoutMS),
			home,
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
		if result.Notes != "" {
			fmt.Fprintf(w, "  note: %s\n", result.Notes)
		}
		if result.Error != "" {
			fmt.Fprintf(w, "  probe error: %s\n", result.Error)
		}
	}
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
