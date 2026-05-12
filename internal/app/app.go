package app

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
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
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/runstate"
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
		return runResume(args[1:], stdout, stderr)
	case "answer":
		return runAnswer(args[1:], stdout, stderr)
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
  %s agents list|verify
  %s status [--dir DIR] [--run RUN_ID] [--idea SLUG] [--json]
  %s run [--no-tui] [--auto] [--participants AGENTS] [--yes] TASK
  %s resume [--dir DIR] [--no-tui] RUN_OR_IDEA
  %s answer [--dir DIR] RUN_ID QUESTION_ID ANSWER...
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
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: parley agents list|verify")
		return 2
	}

	switch args[0] {
	case "list", "discover":
		return runAgentsList(ctx, args[1:], stdout, stderr)
	case "verify", "probe":
		return runAgentsVerify(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintln(stderr, "usage: parley agents list|verify")
		return 2
	}
}

func runAgentsList(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("agents list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	results, err := discoverConfigured(ctx, *root)
	if err != nil {
		fmt.Fprintf(stderr, "agent config failed: %v\n", err)
		return 1
	}
	agents.PrintRuntimeMatrix(stdout, results)
	return 0
}

func runAgentsVerify(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("agents verify", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	agentID := fs.String("agent", "", "agent ID to verify")
	full := fs.Bool("full", false, "run behavioral headless probes")
	yes := fs.Bool("yes", false, "confirm hosted backend probes")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	results, err := discoverConfigured(ctx, *root)
	if err != nil {
		fmt.Fprintf(stderr, "agent config failed: %v\n", err)
		return 1
	}
	selected, err := selectDiscoveries(results, *agentID)
	if err != nil {
		fmt.Fprintf(stderr, "verify failed: %v\n", err)
		return 2
	}
	if *full && !*yes {
		for _, result := range selected {
			if result.ExternalBackend != agents.ExternalLocal {
				fmt.Fprintf(stderr, "verify --full for %s may use an external backend; rerun with --yes to confirm\n", result.ID)
				return 2
			}
		}
	}

	failed := false
	for _, result := range selected {
		switch {
		case !result.Found:
			failed = true
			fmt.Fprintf(stdout, "%s: not installed\n", result.ID)
			if strings.TrimSpace(*agentID) == "" {
				fmt.Fprintln(stdout, "  hint: use --agent ID to verify one configured agent without requiring every built-in agent")
			}
		case result.Error != "":
			failed = true
			fmt.Fprintf(stdout, "%s: version probe failed: %s\n", result.ID, result.Error)
		default:
			fmt.Fprintf(stdout, "%s: installed version=%s\n", result.ID, valueOr(result.Version, "unknown"))
		}
	}
	if failed {
		return 1
	}
	if !*full {
		return 0
	}

	if err := runFullVerification(ctx, *root, selected, stdout); err != nil {
		fmt.Fprintf(stderr, "verify --full failed: %v\n", err)
		return 1
	}
	return 0
}

func runStatus(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	runID := fs.String("run", "", "run ID to inspect")
	ideaSlug := fs.String("idea", "", "idea slug to inspect")
	jsonOut := fs.Bool("json", false, "print unstable JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *runID != "" && *ideaSlug != "" {
		fmt.Fprintln(stderr, "status accepts only one of --run or --idea")
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

	if *runID != "" || *ideaSlug != "" {
		target := *runID
		if target == "" {
			target = *ideaSlug
		}
		run, err := runstate.ResolveRun(*root, target)
		if err != nil {
			fmt.Fprintf(stderr, "status failed: %v\n", err)
			return 1
		}
		if *jsonOut {
			return printJSON(stdout, run, stderr)
		}
		printRunDetail(stdout, run)
		return 0
	}

	runs, err := runstate.ListRuns(*root)
	if err != nil {
		fmt.Fprintf(stderr, "status failed: %v\n", err)
		return 1
	}
	if *jsonOut {
		payload := struct {
			Workspace protocol.WorkspaceStatus `json:"workspace"`
			Runs      []runstate.RunSummary    `json:"runs"`
		}{Workspace: status, Runs: runs}
		return printJSON(stdout, payload, stderr)
	}

	fmt.Fprintf(stdout, "Transport: %s\n", status.Transport)
	if len(status.Ideas) == 0 {
		fmt.Fprintln(stdout, "Ideas: none")
	} else {
		fmt.Fprintln(stdout, "Ideas:")
		for _, idea := range status.Ideas {
			fmt.Fprintf(stdout, "  %s  status=%s  participants=%s\n", idea.Slug, idea.Status, strings.Join(idea.Participants, ","))
		}
	}
	printRunsOverview(stdout, runs, 10)
	return 0
}

func runResume(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("resume", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	noTUI := fs.Bool("no-tui", false, "print run detail without opening the TUI")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley resume [--dir DIR] [--no-tui] RUN_OR_IDEA")
		return 2
	}

	run, err := runstate.ResolveRun(*root, fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "resume failed: %v\n", err)
		return 1
	}
	if *noTUI {
		printRunDetail(stdout, run)
		return 0
	}

	status, err := protocol.ReadWorkspaceStatus(*root)
	if err != nil {
		fmt.Fprintf(stderr, "resume failed: %v\n", err)
		return 1
	}
	idea := ideaForRun(status, run)
	if err := tui.RunLive(tui.LiveOptions{
		Status:       status,
		Idea:         idea,
		Participants: run.Participants,
		RunID:        run.RunID,
		RunDir:       run.RunDir,
		Resume:       true,
	}); err != nil {
		fmt.Fprintf(stderr, "resume tui failed: %v\n", err)
		return 1
	}
	return 0
}

func printJSON(stdout io.Writer, value any, stderr io.Writer) int {
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		fmt.Fprintf(stderr, "json failed: %v\n", err)
		return 1
	}
	return 0
}

func printRunsOverview(stdout io.Writer, runs []runstate.RunSummary, limit int) {
	if len(runs) == 0 {
		fmt.Fprintln(stdout, "Runs: none")
		return
	}
	fmt.Fprintln(stdout, "Runs:")
	if limit <= 0 || limit > len(runs) {
		limit = len(runs)
	}
	for _, run := range runs[:limit] {
		if run.Error != "" {
			fmt.Fprintf(stdout, "  %s  error=%s\n", run.RunID, run.Error)
			continue
		}
		fmt.Fprintf(stdout, "  %s  idea=%s  state=%s  agents=%s  questions=%d open\n",
			run.RunID,
			valueOr(run.IdeaSlug, "unknown"),
			displayRunState(run),
			agentProgress(run),
			run.OpenQuestions,
		)
	}
	if len(runs) > limit {
		fmt.Fprintf(stdout, "  ... %d older run(s) hidden; use --run RUN_ID for details\n", len(runs)-limit)
	}
}

func printRunDetail(stdout io.Writer, run runstate.RunSummary) {
	fmt.Fprintf(stdout, "Run: %s\n", run.RunID)
	fmt.Fprintf(stdout, "Idea: %s\n", valueOr(run.IdeaSlug, "unknown"))
	if run.Task != "" {
		fmt.Fprintf(stdout, "Task: %s\n", run.Task)
	}
	if run.Mode != "" {
		fmt.Fprintf(stdout, "Mode: %s\n", run.Mode)
	}
	if run.Error != "" {
		fmt.Fprintf(stdout, "State: error: %s\n", run.Error)
		return
	}
	fmt.Fprintf(stdout, "State: %s\n", displayRunState(run))
	fmt.Fprintf(stdout, "Open questions: %d\n", run.OpenQuestions)

	if len(run.State.Agents) > 0 {
		fmt.Fprintln(stdout, "Agents:")
		for _, agent := range run.State.Agents {
			fmt.Fprintf(stdout, "  %-10s state=%-10s duration=%-8s latest=%s\n",
				agent.ID,
				agent.State,
				formatDuration(agentDuration(agent)),
				valueOr(agent.LatestEvent, "-"),
			)
			if agent.ArtifactPath != "" {
				fmt.Fprintf(stdout, "             artifact: %s\n", agent.ArtifactPath)
			}
			if agent.StdoutPath != "" || agent.StderrPath != "" {
				fmt.Fprintf(stdout, "             logs: stdout=%s stderr=%s\n", valueOr(agent.StdoutPath, "-"), valueOr(agent.StderrPath, "-"))
			}
			if agent.Error != "" {
				fmt.Fprintf(stdout, "             error: %s\n", agent.Error)
			}
			if agent.Reason != "" {
				fmt.Fprintf(stdout, "             reason: %s\n", agent.Reason)
			}
		}
	}

	if len(run.Questions) > 0 {
		fmt.Fprintln(stdout, "Open HITL questions:")
		printed := false
		for _, question := range run.Questions {
			if question.Status != hitl.StatusOpen {
				continue
			}
			printed = true
			fmt.Fprintf(stdout, "  %s  agent=%s  risk=%s  %s\n", question.ID, question.Agent, question.Risk, question.Prompt)
		}
		if !printed {
			fmt.Fprintln(stdout, "  none")
		}
	}

	if len(run.State.Recent) > 0 {
		fmt.Fprintln(stdout, "Recent events:")
		for _, event := range run.State.Recent {
			fmt.Fprintf(stdout, "  %s  %-16s %s\n", event.Time.Format(time.RFC3339), event.Type, event.Text)
		}
	}

	fmt.Fprintf(stdout, "Next: %s\n", nextRunAction(run))
}

func displayRunState(run runstate.RunSummary) string {
	if run.Terminal {
		return valueOr(run.Outcome, "unknown")
	}
	state := valueOr(run.Liveness, "idle")
	if !run.LastEventAt.IsZero() {
		return fmt.Sprintf("%s — last event %s ago", state, formatDuration(run.LastEventAge))
	}
	return state
}

func agentProgress(run runstate.RunSummary) string {
	total := len(run.State.Agents)
	if total == 0 {
		total = len(run.Participants)
	}
	if total == 0 {
		return "-"
	}
	done := 0
	for _, agent := range run.State.Agents {
		switch agent.State {
		case runstate.StateFinished, runstate.StateFailed, runstate.StateSkipped:
			done++
		}
	}
	return fmt.Sprintf("%d/%d", done, total)
}

func agentDuration(agent runstate.AgentState) time.Duration {
	if agent.Duration > 0 {
		return agent.Duration
	}
	if agent.State == runstate.StateRunning && !agent.StartedAt.IsZero() {
		if elapsed := time.Since(agent.StartedAt); elapsed > 0 {
			return elapsed
		}
	}
	return 0
}

func formatDuration(value time.Duration) string {
	if value <= 0 {
		return "-"
	}
	if value < time.Second {
		return value.Round(time.Millisecond).String()
	}
	return value.Round(time.Second).String()
}

func nextRunAction(run runstate.RunSummary) string {
	for _, question := range run.Questions {
		if question.Status == hitl.StatusOpen {
			return fmt.Sprintf("parley answer %s %s <answer>", run.RunID, question.ID)
		}
	}
	if !run.Terminal {
		return fmt.Sprintf("parley resume %s", run.RunID)
	}
	return "no recoverable action; inspect artifacts/logs"
}

func ideaForRun(status protocol.WorkspaceStatus, run runstate.RunSummary) protocol.IdeaStatus {
	for _, idea := range status.Ideas {
		if idea.Slug == run.IdeaSlug {
			return idea
		}
	}
	slug := valueOr(run.IdeaSlug, "unknown")
	idea := protocol.IdeaStatus{
		Slug:         slug,
		Status:       "unknown",
		Participants: run.Participants,
	}
	if slug != "unknown" {
		idea.Path = filepath.Join(status.Root, protocol.DeckDir, "ideas", slug)
	}
	return idea
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

	workspaceStatus, err := protocol.ReadWorkspaceStatus(*root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(stderr, "no parley-deck workspace found; run `parley init` first")
			return 1
		}
		fmt.Fprintf(stderr, "workspace read failed: %v\n", err)
		return 1
	}

	discovered, err := discoverConfigured(ctx, *root)
	if err != nil {
		fmt.Fprintf(stderr, "agent config failed: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, "Resolved agent runtime:")
	agents.PrintRuntimeMatrix(stdout, discovered)
	participants, err := selectedParticipantIDs(discovered, *participantsFlag)
	if err != nil {
		fmt.Fprintf(stderr, "participant selection failed: %v\n", err)
		return 1
	}
	if len(participants) == 0 {
		fmt.Fprintln(stderr, "no installed headless agents found; run `parley agents list` to inspect configuration")
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
			"runtime":      runtimeEventData(discovered),
		},
	}); err != nil {
		fmt.Fprintf(stderr, "run create failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Created idea %s and run %s\n", idea.Slug, runID)
	fmt.Fprintf(stdout, "Starting round-01 with participants: %s\n", strings.Join(participants, ", "))
	runOpts := runner.Options{
		Root:   *root,
		RunID:  runID,
		Idea:   idea,
		Task:   task,
		Agents: discovered,
		Store:  runStore,
	}
	if *noTUI {
		runCtx, cancelRun := context.WithCancel(ctx)
		defer cancelRun()
		if *auto {
			startAutoAnswerer(runCtx, handleRunDir(*root, runID))
		}
		results := runner.RunRoundOne(runCtx, runOpts)
		cancelRun()
		failed := printRunResults(stdout, results)
		if failed {
			return 1
		}
		return 0
	}

	runCtx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()
	if *auto {
		startAutoAnswerer(runCtx, handleRunDir(*root, runID))
	}
	handle := runner.RunRoundOneAsync(runCtx, runOpts)
	workspaceStatus.Ideas = []protocol.IdeaStatus{idea}
	if err := tui.RunLive(tui.LiveOptions{
		Status:       workspaceStatus,
		Idea:         idea,
		Participants: participants,
		RunID:        runID,
		RunDir:       handle.RunDir,
		Done:         handle.Done(),
		Cancel:       cancelRun,
	}); err != nil {
		cancelRun()
		results := handle.Wait()
		printRunResults(stdout, results)
		fmt.Fprintf(stderr, "tui failed: %v\n", err)
		return 1
	}
	results := handle.Wait()
	failed := printRunResults(stdout, results)
	if failed {
		return 1
	}
	return 0
}

func runAnswer(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("answer", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	rest := fs.Args()
	if len(rest) < 3 {
		fmt.Fprintln(stderr, "usage: parley answer [--dir DIR] RUN_ID QUESTION_ID ANSWER...")
		return 2
	}

	runID := rest[0]
	questionID := rest[1]
	answer := strings.TrimSpace(strings.Join(rest[2:], " "))
	runDir := handleRunDir(*root, runID)
	question, err := hitl.New(runDir).Answer(questionID, answer, false)
	if err != nil {
		fmt.Fprintf(stderr, "answer failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Answered %s for run %s\n", question.ID, runID)
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

func handleRunDir(root, runID string) string {
	return filepath.Join(root, protocol.DeckDir, "runs", runID)
}

func discoverConfigured(ctx context.Context, root string) ([]agents.Discovery, error) {
	specs, err := config.LoadAgentSpecs(root)
	if err != nil {
		return nil, err
	}
	return agents.Discover(ctx, specs), nil
}

func selectDiscoveries(results []agents.Discovery, agentID string) ([]agents.Discovery, error) {
	if strings.TrimSpace(agentID) == "" {
		return results, nil
	}
	for _, result := range results {
		if result.ID == agentID {
			return []agents.Discovery{result}, nil
		}
	}
	return nil, fmt.Errorf("unknown agent %s", agentID)
}

func runFullVerification(ctx context.Context, root string, selected []agents.Discovery, stdout io.Writer) error {
	runID := store.NewRunID(time.Now())
	probeDir := filepath.Join(root, protocol.DeckDir, "meta", "runtime-probes", runID)
	if err := os.MkdirAll(probeDir, 0o755); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "probe dir: %s\n", probeDir)

	var failures []string
	for _, result := range selected {
		if err := runHeadlessProbe(ctx, root, probeDir, runID, result, stdout); err != nil {
			fmt.Fprintf(stdout, "%s: full probe failed: %v\n", result.ID, err)
			failures = append(failures, result.ID)
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("%d full verification probe(s) failed: %s", len(failures), strings.Join(failures, ", "))
	}
	return nil
}

func runHeadlessProbe(ctx context.Context, root, probeDir, runID string, result agents.Discovery, stdout io.Writer) error {
	outputPath := filepath.Join(probeDir, result.ID+".md")
	sentinel := fmt.Sprintf("# parley-runtime-probe agent=%s run=%s", result.ID, runID)
	prompt := probePrompt(result, outputPath, sentinel)

	cmd, cleanup, err := runner.CommandFor(ctx, root, result, prompt)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return fmt.Errorf("%s: %w", result.ID, err)
	}
	cmd.Dir = root
	if result.PromptMode == agents.PromptStdin {
		cmd.Stdin = strings.NewReader(prompt)
	}
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s headless probe: %w: stdout=%s stderr=%s", result.ID, err, strings.TrimSpace(out.String()), strings.TrimSpace(errOut.String()))
	}
	data, err := os.ReadFile(outputPath)
	if err != nil {
		return fmt.Errorf("%s headless probe did not create %s: %w", result.ID, outputPath, err)
	}
	if !strings.HasPrefix(string(data), sentinel) {
		return fmt.Errorf("%s headless probe sentinel mismatch in %s", result.ID, outputPath)
	}
	fmt.Fprintf(stdout, "%s: headless probe passed\n", result.ID)
	return nil
}

func probePrompt(result agents.Discovery, outputPath, sentinel string) string {
	extra := ""
	if result.ID == "codex" {
		extra = `
Before writing the success sentinel, run these Git smoke commands in the current repository:

git status
git branch tmp-codex-git-test
git branch -D tmp-codex-git-test
printf test | git hash-object -w --stdin

If any command fails, do not write the success sentinel. Instead, write the failure and return a short blocker.
`
	}
	return fmt.Sprintf(`Create exactly this file and no other file:
%s

The first line must be exactly:
%s
%s

After the sentinel, write one short line confirming the headless probe.`, outputPath, sentinel, extra)
}

func startAutoAnswerer(ctx context.Context, runDir string) {
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := hitl.New(runDir).AutoAnswerOpen(); err != nil {
					continue
				}
			}
		}
	}()
}

func runTUIView(ctx context.Context, root string, stdout, stderr io.Writer) int {
	results, err := discoverConfigured(ctx, root)
	if err != nil {
		fmt.Fprintf(stderr, "agent config failed: %v\n", err)
		return 1
	}
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

func runtimeEventData(discovered []agents.Discovery) []map[string]any {
	data := make([]map[string]any, 0, len(discovered))
	for _, result := range discovered {
		data = append(data, map[string]any{
			"agent":            result.ID,
			"installed":        result.Found,
			"version":          result.Version,
			"sandbox_mode":     result.SandboxMode,
			"approval_policy":  result.ApprovalPolicy,
			"model":            result.Model,
			"reasoning":        result.Reasoning,
			"profile":          result.Profile,
			"speed":            result.Speed,
			"timeout_ms":       result.TimeoutMS,
			"isolate_home":     result.IsolateHome,
			"external_backend": result.ExternalBackend,
			"sources":          result.Sources,
		})
	}
	return data
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
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
