package app

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/tui"
)

const appName = "parley"

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stdout)
		return 0
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	switch args[0] {
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	case "version", "--version":
		fmt.Fprintln(stdout, "parley dev")
		return 0
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "agents":
		return runAgents(ctx, args[1:], stdout, stderr)
	case "status":
		return runStatus(args[1:], stdout, stderr)
	case "run":
		return runTask(ctx, args[1:], stdout, stderr)
	case "resume":
		fmt.Fprintln(stderr, "parley resume is not implemented yet; durable run loading exists, process reattach is next")
		return 1
	case "answer":
		fmt.Fprintln(stderr, "parley answer is not implemented yet; HITL question storage is next")
		return 1
	case "tui":
		return runTUI(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, `%s orchestrates Parley Deck multi-agent workflows.

Usage:
  %s init [--dir DIR]
  %s agents discover|probe
  %s status [--dir DIR]
  %s run [--no-tui] [--auto] [--participants AGENTS] [--yes] TASK
  %s resume RUN_OR_IDEA
  %s answer QUESTION_ID
  %s tui [--dir DIR]
  %s help
  %s version

`, appName, appName, appName, appName, appName, appName, appName, appName, appName, appName)
}

func runInit(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if err := protocol.InitWorkspace(*root); err != nil {
		fmt.Fprintf(stderr, "init failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Initialized Parley Deck workspace at %s\n", filepath.Join(*root, protocol.DeckDir))
	return 0
}

func runAgents(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || (args[0] != "discover" && args[0] != "probe") {
		fmt.Fprintln(stderr, "usage: parley agents discover|probe")
		return 2
	}

	results := agents.Discover(ctx, agents.DefaultSpecs())
	agents.PrintDiscovery(stdout, results)
	return 0
}

func runStatus(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	status, err := protocol.ReadWorkspaceStatus(*root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(stderr, "no parley-deck workspace found; run `parley init` first")
			return 1
		}
		fmt.Fprintf(stderr, "status failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Transport: %s\n", status.Transport)
	if len(status.Ideas) == 0 {
		fmt.Fprintln(stdout, "Ideas: none")
		return 0
	}

	fmt.Fprintln(stdout, "Ideas:")
	for _, idea := range status.Ideas {
		fmt.Fprintf(stdout, "  %s  status=%s  participants=%s\n", idea.Slug, idea.Status, strings.Join(idea.Participants, ","))
	}
	return 0
}

func runTask(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	noTUI := fs.Bool("no-tui", false, "run without opening the TUI after the round")
	auto := fs.Bool("auto", false, "enable automatic low-risk progression policy")
	participantsFlag := fs.String("participants", "", "comma-separated agent IDs to run")
	yes := fs.Bool("yes", false, "launch selected agents without interactive confirmation")
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	task := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if task == "" {
		fmt.Fprintln(stderr, "usage: parley run [--no-tui] [--auto] [--participants AGENTS] [--yes] TASK")
		return 2
	}

	if _, err := protocol.ReadWorkspaceStatus(*root); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(stderr, "no parley-deck workspace found; run `parley init` first")
			return 1
		}
		fmt.Fprintf(stderr, "workspace read failed: %v\n", err)
		return 1
	}

	discovered := agents.Discover(ctx, agents.DefaultSpecs())
	participants, err := selectedParticipantIDs(discovered, *participantsFlag)
	if err != nil {
		fmt.Fprintf(stderr, "participant selection failed: %v\n", err)
		return 1
	}
	if len(participants) == 0 {
		fmt.Fprintln(stderr, "no installed headless agents found; run `parley agents discover` to inspect configuration")
		return 1
	}
	if !*auto && !*yes && !confirmLaunch(os.Stdin, stdout, participants) {
		fmt.Fprintln(stdout, "No run started. Use `--yes` or `--auto` to launch without an interactive confirmation prompt.")
		return 0
	}
	idea, err := protocol.CreateIdea(*root, task, participants)
	if err != nil {
		fmt.Fprintf(stderr, "idea create failed: %v\n", err)
		return 1
	}

	runID := store.NewRunID(time.Now())
	runStore := store.New(filepath.Join(*root, protocol.DeckDir, "runs", runID))
	if err := runStore.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "run.created",
		Data: map[string]any{
			"task":         task,
			"mode":         modeName(*auto),
			"idea":         idea.Slug,
			"participants": participants,
		},
	}); err != nil {
		fmt.Fprintf(stderr, "run create failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Created idea %s and run %s\n", idea.Slug, runID)
	fmt.Fprintf(stdout, "Starting round-01 with participants: %s\n", strings.Join(participants, ", "))
	results := runner.RunRoundOne(ctx, runner.Options{
		Root:   *root,
		RunID:  runID,
		Idea:   idea,
		Task:   task,
		Agents: discovered,
		Store:  runStore,
	})
	failed := printRunResults(stdout, results)
	if *noTUI {
		if failed {
			return 1
		}
		return 0
	}
	tuiCode := runTUIViewWithDiscovery(ctx, *root, discovered, stdout, stderr)
	if tuiCode != 0 {
		return tuiCode
	}
	if failed {
		return 1
	}
	return 0
}

func runTUI(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("tui", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	return runTUIView(ctx, *root, stdout, stderr)
}

func runTUIView(ctx context.Context, root string, stdout, stderr io.Writer) int {
	results := agents.Discover(ctx, agents.DefaultSpecs())
	return runTUIViewWithDiscovery(ctx, root, results, stdout, stderr)
}

func runTUIViewWithDiscovery(ctx context.Context, root string, results []agents.Discovery, stdout, stderr io.Writer) int {
	status, err := protocol.ReadWorkspaceStatus(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(stderr, "no parley-deck workspace found; run `parley init` first")
			return 1
		}
		fmt.Fprintf(stderr, "tui failed: %v\n", err)
		return 1
	}

	if err := tui.Run(status, results); err != nil {
		fmt.Fprintf(stderr, "tui failed: %v\n", err)
		return 1
	}
	return 0
}

func modeName(auto bool) string {
	if auto {
		return "auto"
	}
	return "hitl"
}

func installedAgentIDs(discovered []agents.Discovery) []string {
	ids := make([]string, 0, len(discovered))
	for _, result := range discovered {
		if result.Found {
			ids = append(ids, result.ID)
		}
	}
	return ids
}

func selectedParticipantIDs(discovered []agents.Discovery, raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return installedAgentIDs(discovered), nil
	}

	installed := map[string]bool{}
	for _, result := range discovered {
		if result.Found {
			installed[result.ID] = true
		}
	}

	var selected []string
	seen := map[string]bool{}
	for _, part := range strings.Split(raw, ",") {
		id := strings.TrimSpace(part)
		if id == "" || seen[id] {
			continue
		}
		if !installed[id] {
			return nil, fmt.Errorf("%s is not an installed agent", id)
		}
		selected = append(selected, id)
		seen[id] = true
	}
	return selected, nil
}

func confirmLaunch(stdin *os.File, stdout io.Writer, participants []string) bool {
	info, err := stdin.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice == 0 {
		fmt.Fprintf(stdout, "HITL confirmation required before launching agents: %s\n", strings.Join(participants, ", "))
		return false
	}

	fmt.Fprintf(stdout, "Launch round-01 with participants %s? [y/N] ", strings.Join(participants, ", "))
	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil {
		return false
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes"
}

func printRunResults(stdout io.Writer, results []runner.Result) bool {
	if len(results) == 0 {
		fmt.Fprintln(stdout, "No participants were selected for round-01.")
		return true
	}

	failed := false
	for _, result := range results {
		switch {
		case result.Skipped && result.ExitError == "":
			fmt.Fprintf(stdout, "  %-8s skipped: %s\n", result.AgentID, result.SkipReason)
		case result.ArtifactOK && result.ExitError == "":
			fmt.Fprintf(stdout, "  %-8s wrote %s in %s\n", result.AgentID, result.OutputPath, result.Duration.Round(time.Millisecond))
		default:
			failed = true
			reason := result.ExitError
			if reason == "" {
				reason = "artifact was not created"
			}
			fmt.Fprintf(stdout, "  %-8s failed: %s\n", result.AgentID, reason)
			if result.StdoutPath != "" || result.StderrPath != "" {
				fmt.Fprintf(stdout, "           logs: %s %s\n", result.StdoutPath, result.StderrPath)
			}
		}
	}
	return failed
}
