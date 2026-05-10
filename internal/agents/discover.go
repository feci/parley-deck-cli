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
	ID           string
	Commands     []string
	VersionArgs  []string
	HeadlessMode string
	HeadlessArgs []string
	PromptMode   PromptMode
	IsolateHome  bool
	Telemetry    string
	Notes        string
}

type PromptMode string

const (
	PromptStdin PromptMode = "stdin"
	PromptArg   PromptMode = "arg"
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
		{
			ID:           "codex",
			Commands:     []string{"codex"},
			VersionArgs:  []string{"--version"},
			HeadlessMode: "codex exec -",
			HeadlessArgs: []string{"exec", "--cd", "{root}", "--sandbox", "workspace-write", "--ask-for-approval", "never", "-"},
			PromptMode:   PromptStdin,
			Telemetry:    "json events when --json is available",
		},
		{
			ID:           "claude",
			Commands:     []string{"claude"},
			VersionArgs:  []string{"--version"},
			HeadlessMode: "claude --print",
			HeadlessArgs: []string{"-p", "--output-format", "text", "--permission-mode", "acceptEdits", "--add-dir", "{root}"},
			PromptMode:   PromptStdin,
			Telemetry:    "stream-json or final text depending on flags",
		},
		{
			ID:           "gemini",
			Commands:     []string{"gemini"},
			VersionArgs:  []string{"--version"},
			HeadlessMode: "gemini --prompt ... --output-format json",
			HeadlessArgs: []string{"--prompt", "Follow the Parley Deck participant instructions provided on stdin.", "--skip-trust", "--approval-mode", "auto_edit", "--output-format", "text"},
			PromptMode:   PromptStdin,
			IsolateHome:  true,
			Telemetry:    "json stats when output-format json succeeds",
			Notes:        "uses isolated GEMINI_CLI_HOME for oauth-personal profiles that hang",
		},
		{
			ID:           "hermes",
			Commands:     []string{"hermes", "hermes-agent", "hermesagent"},
			VersionArgs:  []string{"--version"},
			HeadlessMode: "hermes --oneshot",
			HeadlessArgs: []string{"--oneshot", "{prompt}", "--accept-hooks"},
			PromptMode:   PromptArg,
			IsolateHome:  true,
			Telemetry:    "unknown",
			Notes:        "first supported command found on PATH is used; uses isolated HERMES_HOME for writable logs",
		},
	}
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
