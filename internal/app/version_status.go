package app

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runVersion(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(stderr)
	all := fs.Bool("all", false, "include Parley Deck skill and project status")
	jsonOut := fs.Bool("json", false, "print JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "usage: parley version [--all] [--json]")
		return 2
	}

	if !*all {
		if *jsonOut {
			writeVersionJSON(stdout, nil, "")
			return 0
		}
		fmt.Fprintln(stdout, versionLine())
		return 0
	}

	skillStatus, skillError := parleyDeckSkillStatus(ctx)
	if *jsonOut {
		writeVersionJSON(stdout, skillStatus, skillError)
		return 0
	}

	fmt.Fprintln(stdout, versionLine())
	if skillError != "" {
		fmt.Fprintf(stdout, "parley-deck-skill: unavailable (%s)\n", skillError)
		return 0
	}
	printSkillStatusSummary(stdout, skillStatus)
	return 0
}

func writeVersionJSON(stdout io.Writer, skillStatus map[string]any, skillError string) {
	payload := map[string]any{
		"ok": skillError == "",
		"parley": map[string]any{
			"name":    appName,
			"version": version,
			"line":    versionLine(),
		},
	}
	if skillStatus != nil {
		payload["parley_deck_skill"] = skillStatus
	}
	if skillError != "" {
		payload["parley_deck_skill_error"] = skillError
	}
	_ = json.NewEncoder(stdout).Encode(payload)
}

func parleyDeckSkillStatus(ctx context.Context) (map[string]any, string) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	probeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(probeCtx, "parley-deck-skill", "status", "--target", "all", "--project", cwd, "--json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(errOut.String())
		if message == "" {
			message = err.Error()
		}
		return legacyParleyDeckSkillStatus(probeCtx, message)
	}

	var status map[string]any
	if err := json.Unmarshal(out.Bytes(), &status); err != nil {
		return nil, "parley-deck-skill returned invalid JSON"
	}
	return status, ""
}

func legacyParleyDeckSkillStatus(ctx context.Context, statusError string) (map[string]any, string) {
	cmd := exec.CommandContext(ctx, "parley-deck-skill", "--version")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(errOut.String())
		if message == "" {
			message = err.Error()
		}
		return nil, statusError + "; version probe failed: " + message
	}

	return map[string]any{
		"ok":              false,
		"statusSupported": false,
		"statusError":     statusError,
		"installer": map[string]any{
			"version": strings.TrimSpace(out.String()),
			"source":  "legacy-command",
		},
		"compatibility": map[string]any{
			"status":  "unknown",
			"reasons": []any{"status-command-unavailable"},
		},
	}, ""
}

func printSkillStatusSummary(stdout io.Writer, status map[string]any) {
	installer := mapValue(status, "installer")
	compatibility := mapValue(status, "compatibility")
	project := mapValue(status, "project")

	installerVersion := stringValue(installer, "version")
	installerSource := stringValue(installer, "source")
	if installerVersion == "" {
		installerVersion = "unknown"
	}
	if installerSource != "" {
		fmt.Fprintf(stdout, "parley-deck-skill %s (%s)\n", installerVersion, installerSource)
	} else {
		fmt.Fprintf(stdout, "parley-deck-skill %s\n", installerVersion)
	}

	if status := stringValue(compatibility, "status"); status != "" {
		fmt.Fprintf(stdout, "compatibility: %s\n", status)
	}
	if value, ok := status["statusSupported"].(bool); ok && !value {
		fmt.Fprintln(stdout, "parley-deck-skill status: unsupported by installed command")
	}
	if metadataStatus := stringValue(project, "metadataStatus"); metadataStatus != "" {
		fmt.Fprintf(stdout, "project metadata: %s\n", metadataStatus)
	}
}

func mapValue(parent map[string]any, key string) map[string]any {
	if value, ok := parent[key].(map[string]any); ok {
		return value
	}
	return nil
}

func stringValue(parent map[string]any, key string) string {
	if parent == nil {
		return ""
	}
	if value, ok := parent[key].(string); ok {
		return value
	}
	return ""
}
