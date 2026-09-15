package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

func runAgentsExec(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("agents exec", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace root")
	agentID := fs.String("agent", "", "configured agent or roster ID")
	promptFile := fs.String("prompt-file", "", "prompt file; use - for stdin")
	artifact := fs.String("artifact", "", "optional expected new artifact, relative to the workspace")
	phase := fs.String("phase", "manual", "launch phase")
	idea := fs.String("idea", "", "idea identity for the audit record")
	runID := fs.String("run-id", "", "run identity; generated when omitted")
	segment := fs.String("segment-id", "manual", "run segment identity")
	timeout := fs.Duration("timeout", 0, "hard timeout; default is the configured agent timeout")
	yes := fs.Bool("yes", false, "confirm the configured agent invocation and external backend")
	jsonOut := fs.Bool("json", false, "print content-free terminal evidence")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 || *agentID == "" || *promptFile == "" || !*yes || *timeout < 0 {
		fmt.Fprintln(stderr, "usage: parley agents exec --agent ID --prompt-file FILE [--dir DIR] [--artifact PATH] [--timeout D] [--json] --yes")
		return 2
	}
	for _, value := range []string{*agentID, *phase, *idea, *runID, *segment} {
		if value != "" && telemetry.SafeLabel(value) == nil {
			fmt.Fprintln(stderr, "agents exec: invalid audit identifier")
			return 2
		}
	}
	if strings.ContainsAny(*runID, `/\`) || strings.Contains(*runID, "..") {
		fmt.Fprintln(stderr, "agents exec: run identity must be one safe path component")
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintln(stderr, "agents exec: invalid workspace")
		return 2
	}
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		fmt.Fprintln(stderr, "agents exec: workspace must exist")
		return 2
	}
	artifactPath := ""
	if *artifact != "" {
		artifactPath = *artifact
		if !filepath.IsAbs(artifactPath) {
			artifactPath = filepath.Join(root, artifactPath)
		}
		rel, err := filepath.Rel(root, artifactPath)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			fmt.Fprintln(stderr, "agents exec: artifact must be inside the workspace")
			return 2
		}
		if _, err := os.Lstat(artifactPath); !os.IsNotExist(err) {
			fmt.Fprintln(stderr, "agents exec: refusing an existing or inaccessible artifact")
			return 2
		}
	}
	var input io.Reader = os.Stdin
	if *promptFile != "-" {
		file, err := os.Open(*promptFile)
		if err != nil {
			fmt.Fprintln(stderr, "agents exec: cannot read prompt file")
			return 1
		}
		defer file.Close()
		input = file
	}
	prompt, err := io.ReadAll(io.LimitReader(input, (1<<20)+1))
	if err != nil || len(prompt) > 1<<20 || strings.TrimSpace(string(prompt)) == "" {
		fmt.Fprintln(stderr, "agents exec: prompt must be nonempty and at most 1 MiB")
		return 2
	}
	discovered, err := discoverConfigured(ctx, root)
	if err != nil {
		fmt.Fprintln(stderr, "agents exec: cannot load configured agents")
		return 1
	}
	mapping, err := config.LoadRosterAdapters(root)
	if err != nil {
		fmt.Fprintln(stderr, "agents exec: invalid roster mapping")
		return 1
	}
	agent, err := resolveMeasuredAgent(*agentID, discovered, mapping)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if *runID == "" {
		*runID, err = telemetry.NewID()
		if err != nil {
			fmt.Fprintln(stderr, "agents exec: cannot allocate run identity")
			return 1
		}
	}
	info := runner.LaunchInfo{RunID: *runID, SegmentID: *segment, Idea: *idea, Phase: *phase,
		ArtifactPath: artifactPath, Store: store.New(filepath.Join(root, ".parley-runtime", "manual-runs", *runID))}
	record, launchErr := runner.RunMeasured(ctx, runner.ExecOptions{Root: root, Agent: agent, Prompt: string(prompt), Timeout: *timeout, Info: info})
	if record.InvocationID == "" {
		fmt.Fprintln(stderr, "agents exec: launch setup failed before an invocation could be recorded; check headless mode and telemetry storage")
		return 1
	}
	if *jsonOut {
		if err := json.NewEncoder(stdout).Encode(record); err != nil {
			fmt.Fprintln(stderr, "agents exec: cannot write result")
			return 1
		}
	} else if record.InvocationID != "" {
		fmt.Fprintf(stdout, "Invocation: %s\nEvidence: %s\n", record.InvocationID,
			filepath.Join(root, ".parley-runtime", "invocations", record.InvocationID, "terminal.json"))
	}
	if launchErr != nil {
		fmt.Fprintln(stderr, "agents exec: invocation failed; inspect private local evidence")
		return 1
	}
	if artifactPath != "" && (record.Outcome == nil || record.Outcome.ArtifactSHA256 == nil) {
		fmt.Fprintln(stderr, "agents exec: expected artifact was not observed; process exit is not artifact acceptance")
		return 1
	}
	return 0
}

// Unlike roster readiness, a measured explicit attempt preserves a missing
// executable as a failed start rather than losing it before telemetry begins.
func resolveMeasuredAgent(id string, discovered []agents.Discovery, mapping map[string]string) (agents.Discovery, error) {
	candidates := append([]agents.Discovery(nil), discovered...)
	for i := range candidates {
		candidates[i].Found = true
		if candidates[i].Path == "" && len(candidates[i].Commands) > 0 {
			candidates[i].Path = candidates[i].Commands[0]
		}
	}
	return agents.ResolveParticipant(id, candidates, mapping)
}
