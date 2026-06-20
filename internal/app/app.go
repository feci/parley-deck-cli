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
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runaction"
	"parley-deck-cli/internal/runcontrol"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/runplan"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/sessionstore"
	"parley-deck-cli/internal/steer"
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
	case "--version":
		fmt.Fprintln(stdout, versionLine())
		return 0
	case "version":
		return runVersion(ctx, args[1:], stdout, stderr)
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "agents":
		return runAgents(ctx, args[1:], stdout, stderr)
	case "consensus":
		return runConsensus(ctx, args[1:], stdout, stderr)
	case "pipeline":
		return runPipeline(ctx, args[1:], stdout, stderr)
	case "context":
		return runContext(args[1:], stdout, stderr)
	case "status":
		return runStatus(args[1:], stdout, stderr)
	case "sessions":
		return runSessions(args[1:], stdout, stderr)
	case "preflight":
		return runPreflight(ctx, args[1:], stdout, stderr)
	case "run":
		return runTask(ctx, args[1:], stdout, stderr)
	case "continue":
		return runContinue(ctx, args[1:], stdout, stderr)
	case "steer":
		return runSteer(args[1:], stdout, stderr)
	case "resume":
		return runResume(args[1:], stdout, stderr)
	case "answer":
		return runAnswer(args[1:], stdout, stderr)
	case "consult":
		return runConsult(ctx, args[1:], stdout, stderr)
	case "consults":
		return runConsults(args[1:], stdout, stderr)
	case "retro":
		return runRetro(args[1:], stdout, stderr)
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

Parley Deck keeps the cooperation trail in a repository-local parley-deck/
directory. The CLI initializes that workspace, discovers local agent CLIs,
starts or resumes multi-agent rounds, tracks HITL questions, manages consensus
signoffs, and emits repository context that can be handed to agents.

Usage:
  %s init [--dir DIR]
  %s agents list [--dir DIR]
  %s agents verify [--dir DIR] [--agent ID] [--full] [--yes]
  %s consensus status [--dir DIR] [--review] [--json] IDEA
  %s consensus draft [--dir DIR] [--review] [--round N] [--by AGENT] IDEA
  %s consensus signoff [--dir DIR] [--review] --agent ID --status accept|reserve|reservations|block [--notes TEXT] [--counter TEXT] IDEA
  %s consensus request-signoffs [--dir DIR] [--review] [--participants IDS] [--yes] [--dry-run] IDEA
  %s consensus finalize [--dir DIR] [--by AGENT] IDEA
  %s consensus reopen [--dir DIR] [--review] --reason TEXT IDEA
  %s context repo-map [--dir DIR] [--format markdown|json] [--max-files N]
  %s status [--dir DIR] [--run RUN_ID] [--idea SLUG] [--json]
  %s sessions list [--json]
  %s sessions inspect [--dir DIR] [--json] RUN_ID
  %s consult [--dir DIR] [--timeout D] AGENT [QUESTION]
  %s consults list [--dir DIR]
  %s retro <scan|select|diagnose|propose> [--dir DIR] [--k N] [--json] [--slug SLUG]
  %s preflight [--dir DIR] [--json] [--yes] [--ping-timeout D] [--no-ping]
  %s run [--no-tui] [--no-auto] [--no-preflight] [--no-ping] [--participants AGENTS] [--yes] TASK
  %s continue [--dir DIR] [--json] RUN_OR_IDEA
  %s resume [--dir DIR] [--no-tui] RUN_OR_IDEA
  %s answer [--dir DIR] RUN_ID QUESTION_ID ANSWER...
  %s tui [--dir DIR]
  %s help
  %s version [--dir DIR] [--all] [--json]

Commands:
  init
      Create parley-deck/COOPERATION.md, ideas/, inbox/, meta/, and runs/.

  agents list
      Print the effective runtime matrix for configured agents: installed
      state, command, launch mode, sandbox/approval policy, model/profile,
      timeout, home isolation, backend class, and config sources.

  agents verify
      Check configured agent CLIs. By default this is a cheap version probe.
      Use --full --yes to run behavioral probes that may call hosted backends.

  preflight
      Readiness check before an idea: protocol freshness (advisory in a source
      repo; additive consumer auto-sync preserving project zones; breaking bump
      gated) plus a hosted-PONG roster liveness ping. Exit 0 ready, 3 pending
      gate (names the gate + confirm command), 1 hard failure (no workspace, or
      excluding unavailable agents would leave <2 participants), 2 usage.

  run
      Create a new idea from TASK and start round-01 with selected agents.
      Auto-drive is ON by default: after round-01 the protocol advances
      automatically (cross-review, consensus, finalize, and — if the idea opts
      into auto_implement — the implementation/fix-up phases), shown live in the
      TUI. Pass --no-auto to stop after round-01 and advance manually; with
      --no-auto and without --yes, Parley asks before launching hosted agents.

  resume
      Re-open a run by run ID or idea slug, validate pending handoffs, and
      continue the live TUI unless --no-tui is set.

  continue
      Inspect a run by run ID or idea slug and print planner-derived next
      actions for continuing the workflow. This first slice is read-only; use
      the printed commands for existing safe actions such as answering HITL
      questions or consensus operations.

  pipeline
      Compose the cooperation engine into §12 pipeline blocks via a
      pipeline.yaml manifest. Subcommands: validate MANIFEST, start MANIFEST,
      status SLUG, continue SLUG, gate approve|reject SLUG EDGE. Supervised-first:
      block boundaries gate by default; production mutations are non-bypassable.

  status
      Show workspace, idea, consensus, run, and HITL question state.

  sessions
      List and inspect the local Parley session index in ~/.parley-deck.
      Inspect resolves the stored workspace run and prints manifest details
      when parley-deck/runs/<run-id>/run.json exists.
      Use --dir to disambiguate or inspect an unindexed run from a workspace.
      Slice 1 is read-only and assumes a single active parley process per run.

  answer
      Answer a pending human-in-the-loop question for a run.

  consensus
      Draft, inspect, sign, request signoffs for, finalize, or reopen design
      and review consensus files under parley-deck/ideas/<idea>/.

  context repo-map
      Emit a deterministic local repository map for agent context. Markdown is
      optimized for humans and prompts; JSON is optimized for tooling.

  tui
      Open the project TUI for workspace status, run state, questions, and
      agent/runtime inspection.

  version
      Print the CLI version. With --all, also print parley-deck-skill and
      project metadata compatibility status when available.

Parameters and flags:
  --dir DIR
      Workspace root. Defaults to the current directory.

  --participants AGENTS
      Comma-separated agent IDs, for example claude,agy,hermes.

  --yes
      Confirm an operation that may launch configured hosted backends.

  --auto
      Enable automatic low-risk HITL handling during a run.

  --no-tui
      Run or resume without opening the live terminal UI.

  --json
      Print machine-readable JSON for commands that support it.

  --review
      Target review consensus at review/consensus.md instead of consensus.md.

  --agent ID
      Select one participant or signer by stable Parley agent ID.

  --status VALUE
      Signoff status: accept, reserve, reservations, or block.

  --notes TEXT
      Notes appended to a signoff.

  --counter TEXT
      Counter-proposal required for block signoffs.

  --round N
      Round number used when drafting consensus.

  --by AGENT
      Agent ID recorded as drafter/finalizer.

  --format markdown|json
      Output format for context repo-map. Defaults to markdown.

  --max-files N
      Maximum number of files in repo-map output. Defaults to 1000.

  TASK
      Free-form work request used to create a new idea.

  IDEA / SLUG
      Idea directory name under parley-deck/ideas/.

  RUN_ID
      Run directory ID under parley-deck/runs/.

  QUESTION_ID
      HITL question ID from a run's questions directory.

Examples:
  # Initialize a repository for Parley Deck.
  %s init --dir .

  # Inspect configured agents and runtime defaults.
  %s agents list --dir .
  %s agents verify --dir . --agent claude

  # Generate agent-friendly repository context.
  %s context repo-map --dir . --format markdown --max-files 50
  %s context repo-map --dir . --format json --max-files 10

  # Start a headless round with selected agents.
  %s run --dir . --no-tui --participants claude,agy --yes "Plan the next CLI slice"

  # Resume or inspect active work.
  %s status --dir .
  %s status --dir . --idea repo-map-mvp
  %s sessions list
  %s sessions inspect 20260517T120000.000000000Z
  %s continue --dir . 20260517T120000.000000000Z
  %s resume --dir . 20260517T120000.000000000Z

  # Consensus flow.
  %s consensus status --dir . repo-map-mvp
  %s consensus draft --dir . --round 1 --by codex repo-map-mvp
  %s consensus request-signoffs --dir . --participants claude,agy --yes repo-map-mvp

  # Answer a HITL question.
  %s answer --dir . 20260517T120000.000000000Z q1 "Use the conservative default"

  # Version and compatibility status.
  %s version
  %s version --all

Exit codes:
  0  Success.
  1  Runtime failure, failed probe, malformed workspace, or agent failure.
  2  Usage error or missing required argument.
  3  Pending manual/interactive handoff for consensus request-signoffs.

`, appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
		appName,
	)
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

	// Seed the user-global central default (~/.parley/agents.toml) so every
	// project inherits the same roster + model + reasoning unless it overrides
	// them in parley-deck/agents.toml. Non-fatal: a fresh deck still works
	// from built-in defaults if the home dir is unwritable.
	if path, err := config.EnsureCentralDefault(); err != nil {
		fmt.Fprintf(stderr, "warning: could not write central default %s: %v\n", config.CentralAgentsPath(), err)
	} else if path != "" {
		fmt.Fprintf(stdout, "Central agent defaults: %s (override per-project in %s)\n",
			path, filepath.Join(*root, protocol.DeckDir, "agents.toml"))
	}
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

func runConsensus(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printConsensusUsage(stderr)
		return 2
	}
	switch args[0] {
	case "status":
		return runConsensusStatus(args[1:], stdout, stderr)
	case "draft":
		return runConsensusDraft(args[1:], stdout, stderr)
	case "signoff":
		return runConsensusSignoff(args[1:], stdout, stderr)
	case "request-signoffs":
		return runConsensusRequestSignoffs(ctx, args[1:], stdout, stderr)
	case "finalize":
		return runConsensusFinalize(args[1:], stdout, stderr)
	case "reopen":
		return runConsensusReopen(args[1:], stdout, stderr)
	default:
		printConsensusUsage(stderr)
		return 2
	}
}

func printConsensusUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: parley consensus status|draft|signoff|request-signoffs|finalize|reopen")
}

func runConsensusStatus(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("consensus status", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	review := fs.Bool("review", false, "use review consensus")
	jsonOut := fs.Bool("json", false, "print unstable JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley consensus status [--dir DIR] [--review] [--json] IDEA")
		return 2
	}
	summary, err := consensus.Status(*root, fs.Arg(0), *review)
	if err != nil {
		fmt.Fprintf(stderr, "consensus status failed: %v\n", err)
		return 1
	}
	if *jsonOut {
		return printJSON(stdout, summary, stderr)
	}
	printConsensusSummary(stdout, summary)
	return 0
}

func runConsensusDraft(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("consensus draft", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	review := fs.Bool("review", false, "use review consensus")
	round := fs.Int("round", 0, "round number to draft from")
	by := fs.String("by", "user", "drafting agent ID")
	reviewedCommit := fs.String("reviewed-commit", "", "reviewed commit for review consensus")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley consensus draft [--dir DIR] [--review] [--round N] [--by AGENT] IDEA")
		return 2
	}
	summary, err := consensus.Draft(*root, fs.Arg(0), consensus.DraftOptions{
		Review:         *review,
		Round:          *round,
		By:             *by,
		ReviewedCommit: *reviewedCommit,
	})
	if err != nil {
		fmt.Fprintf(stderr, "consensus draft failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Drafted consensus at %s\n", summary.Path)
	printConsensusSummary(stdout, summary)
	return 0
}

func runConsensusSignoff(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("consensus signoff", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	review := fs.Bool("review", false, "use review consensus")
	agent := fs.String("agent", "", "participant agent ID")
	status := fs.String("status", "", "accept, reserve, reservations, or block")
	notes := fs.String("notes", "", "signoff notes")
	counter := fs.String("counter", "", "counter-proposal for block signoffs")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 || strings.TrimSpace(*agent) == "" || strings.TrimSpace(*status) == "" {
		fmt.Fprintln(stderr, "usage: parley consensus signoff [--dir DIR] [--review] --agent ID --status accept|reserve|reservations|block [--notes TEXT] [--counter TEXT] IDEA")
		return 2
	}
	summary, err := consensus.AppendSignoff(*root, fs.Arg(0), consensus.SignoffOptions{
		Review:          *review,
		Agent:           *agent,
		Status:          *status,
		Notes:           *notes,
		CounterProposal: *counter,
	})
	if err != nil {
		fmt.Fprintf(stderr, "consensus signoff failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Appended signoff for %s to %s\n", *agent, summary.Path)
	printConsensusSummary(stdout, summary)
	return 0
}

func runConsensusFinalize(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("consensus finalize", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	by := fs.String("by", "user", "final plan author")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley consensus finalize [--dir DIR] [--by AGENT] IDEA")
		return 2
	}
	path, summary, err := consensus.Finalize(*root, fs.Arg(0), consensus.FinalizeOptions{By: *by})
	if err != nil {
		fmt.Fprintf(stderr, "consensus finalize failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Finalized consensus and created %s\n", path)
	printConsensusSummary(stdout, summary)
	return 0
}

func runConsensusReopen(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("consensus reopen", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	review := fs.Bool("review", false, "use review consensus")
	reason := fs.String("reason", "", "reason for reopening consensus")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 || strings.TrimSpace(*reason) == "" {
		fmt.Fprintln(stderr, "usage: parley consensus reopen [--dir DIR] [--review] --reason TEXT IDEA")
		return 2
	}
	path, err := consensus.Reopen(*root, fs.Arg(0), consensus.ReopenOptions{
		Review: *review,
		Reason: *reason,
	})
	if err != nil {
		fmt.Fprintf(stderr, "consensus reopen failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Reopened consensus; preserved blocked draft at %s\n", path)
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

	if *runID != "" {
		target := *runID
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
	if *ideaSlug != "" {
		run, err := runstate.ResolveRun(*root, *ideaSlug)
		if err == nil {
			if *jsonOut {
				return printJSON(stdout, run, stderr)
			}
			printRunDetail(stdout, run)
			printConsensusIfPresent(stdout, *root, *ideaSlug)
			return 0
		}
		idea, ok := findIdeaStatus(status, *ideaSlug)
		if !ok {
			fmt.Fprintf(stderr, "status failed: %v\n", err)
			return 1
		}
		if *jsonOut {
			payload := struct {
				Idea      protocol.IdeaStatus `json:"idea"`
				Consensus *consensus.Summary  `json:"consensus,omitempty"`
				Review    *consensus.Summary  `json:"review_consensus,omitempty"`
			}{Idea: idea}
			if summary, err := consensus.Status(*root, *ideaSlug, false); err == nil {
				payload.Consensus = &summary
			}
			if summary, err := consensus.Status(*root, *ideaSlug, true); err == nil {
				payload.Review = &summary
			}
			return printJSON(stdout, payload, stderr)
		}
		printIdeaDetail(stdout, idea)
		printConsensusIfPresent(stdout, *root, *ideaSlug)
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
			fmt.Fprintf(stdout, "  %s  status=%s  participants=%s%s\n", idea.Slug, idea.Status, strings.Join(idea.Participants, ","), consensusTriageLabel(*root, idea.Slug))
		}
	}
	printRunsOverview(stdout, runs, 10)
	return 0
}

func runSessions(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: parley sessions list [--json]\n       parley sessions inspect [--dir DIR] [--json] RUN_ID")
		return 2
	}
	switch args[0] {
	case "list":
		return runSessionsList(args[1:], stdout, stderr)
	case "inspect":
		return runSessionsInspect(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown sessions command: %s\n", args[0])
		return 2
	}
}

func runSessionsList(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("sessions list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOut := fs.Bool("json", false, "print JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "usage: parley sessions list [--json]")
		return 2
	}
	sessionStore, err := sessionstore.Default()
	if err != nil {
		fmt.Fprintf(stderr, "sessions list failed: %v\n", err)
		return 1
	}
	sessions, err := sessionStore.List()
	if err != nil {
		fmt.Fprintf(stderr, "sessions list failed: %v\n", err)
		return 1
	}
	if *jsonOut {
		payload := struct {
			IndexPath string                 `json:"index_path"`
			Sessions  []sessionstore.Session `json:"sessions"`
		}{IndexPath: sessionStore.Path(), Sessions: sessions}
		return printJSON(stdout, payload, stderr)
	}

	fmt.Fprintf(stdout, "Session index: %s\n", sessionStore.Path())
	if len(sessions) == 0 {
		fmt.Fprintln(stdout, "Sessions: none")
		return 0
	}
	fmt.Fprintln(stdout, "Sessions:")
	for _, session := range sessions {
		fmt.Fprintf(stdout, "  %s  idea=%s  workspace=%s  started=%s  status=%s\n",
			session.RunID,
			valueOr(session.IdeaSlug, "unknown"),
			session.WorkspaceRoot,
			sessionStartedTimestamp(session),
			sessionStatus(session),
		)
		if len(session.Participants) > 0 {
			fmt.Fprintf(stdout, "      participants=%s\n", strings.Join(session.Participants, ","))
		}
	}
	return 0
}

func runSessionsInspect(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("sessions inspect", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory used to disambiguate or inspect an unindexed run")
	jsonOut := fs.Bool("json", false, "print JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley sessions inspect [--dir DIR] [--json] RUN_ID")
		return 2
	}

	sessionStore, err := sessionstore.Default()
	if err != nil {
		fmt.Fprintf(stderr, "sessions inspect failed: %v\n", err)
		return 1
	}
	filterRoot := ""
	fs.Visit(func(flag *flag.Flag) {
		if flag.Name == "dir" {
			filterRoot = *root
		}
	})
	session, err := resolveIndexedSession(sessionStore, fs.Arg(0), filterRoot, *root)
	if err != nil {
		fmt.Fprintf(stderr, "sessions inspect failed: %v\n", err)
		return 1
	}
	payload := inspectSession(sessionStore.Path(), session)
	if *jsonOut {
		return printJSON(stdout, payload, stderr)
	}
	printSessionDetail(stdout, payload)
	return 0
}

type sessionInspectPayload struct {
	IndexPath     string                `json:"index_path"`
	Session       sessionstore.Session  `json:"session"`
	Manifest      *runmanifest.Manifest `json:"manifest,omitempty"`
	ManifestPath  string                `json:"manifest_path,omitempty"`
	ManifestError string                `json:"manifest_error,omitempty"`
	Run           *runstate.RunSummary  `json:"run,omitempty"`
	RunError      string                `json:"run_error,omitempty"`
}

func resolveIndexedSession(sessionStore sessionstore.Store, target, filterRoot, fallbackRoot string) (sessionstore.Session, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return sessionstore.Session{}, fmt.Errorf("run ID is required")
	}
	var rootAbs string
	if strings.TrimSpace(filterRoot) != "" {
		var err error
		rootAbs, err = filepath.Abs(filterRoot)
		if err != nil {
			return sessionstore.Session{}, err
		}
	}

	sessions, err := sessionStore.List()
	if err != nil {
		return sessionstore.Session{}, err
	}
	var matches []sessionstore.Session
	for _, session := range sessions {
		if session.RunID != target {
			continue
		}
		if rootAbs != "" && session.WorkspaceRoot != rootAbs {
			continue
		}
		matches = append(matches, session)
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		var roots []string
		for _, match := range matches {
			roots = append(roots, match.WorkspaceRoot)
		}
		return sessionstore.Session{}, fmt.Errorf("run %s exists in multiple workspaces; retry with --dir (%s)", target, strings.Join(roots, ", "))
	}
	if strings.TrimSpace(fallbackRoot) == "" {
		return sessionstore.Session{}, fmt.Errorf("run %s is not in the local session index", target)
	}

	fallbackAbs, err := filepath.Abs(fallbackRoot)
	if err != nil {
		return sessionstore.Session{}, err
	}
	run, err := runstate.ResolveRun(fallbackAbs, target)
	if err != nil {
		if rootAbs == "" {
			return sessionstore.Session{}, fmt.Errorf("run %s is not in the local session index and was not found in workspace %s: %w", target, fallbackAbs, err)
		}
		return sessionstore.Session{}, fmt.Errorf("run %s was not found for workspace %s: %w", target, fallbackAbs, err)
	}
	return sessionstore.Session{
		WorkspaceRoot: fallbackAbs,
		RunID:         run.RunID,
		IdeaSlug:      run.IdeaSlug,
		Task:          run.Task,
		Participants:  append([]string(nil), run.Participants...),
		UpdatedAt:     time.Now().UTC(),
		LastEventAt:   run.LastEventAt,
		Terminal:      run.Terminal,
	}, nil
}

func inspectSession(indexPath string, session sessionstore.Session) sessionInspectPayload {
	payload := sessionInspectPayload{
		IndexPath:    indexPath,
		Session:      session,
		ManifestPath: runmanifest.Path(session.WorkspaceRoot, session.RunID),
	}
	manifest, err := runmanifest.Load(session.WorkspaceRoot, session.RunID)
	if err == nil {
		payload.Manifest = &manifest
	} else if !errors.Is(err, os.ErrNotExist) {
		payload.ManifestError = err.Error()
	}
	run, err := runstate.LoadRun(session.WorkspaceRoot, session.RunID)
	if err == nil {
		payload.Run = &run
	} else {
		payload.RunError = err.Error()
	}
	return payload
}

func printSessionDetail(stdout io.Writer, payload sessionInspectPayload) {
	session := payload.Session
	fmt.Fprintf(stdout, "Session: %s\n", session.RunID)
	fmt.Fprintf(stdout, "Index: %s\n", payload.IndexPath)
	fmt.Fprintf(stdout, "Workspace: %s\n", session.WorkspaceRoot)
	fmt.Fprintf(stdout, "Idea: %s\n", valueOr(session.IdeaSlug, "unknown"))
	if session.Task != "" {
		fmt.Fprintf(stdout, "Task: %s\n", session.Task)
	}
	if len(session.Participants) > 0 {
		fmt.Fprintf(stdout, "Participants: %s\n", strings.Join(session.Participants, ","))
	}
	fmt.Fprintf(stdout, "Status: %s\n", sessionStatus(session))
	fmt.Fprintf(stdout, "Last event: %s\n", sessionTimestamp(session))
	if payload.Manifest != nil {
		fmt.Fprintf(stdout, "Manifest: %s\n", payload.ManifestPath)
		fmt.Fprintf(stdout, "Manifest schema: %d\n", payload.Manifest.SchemaVersion)
		if payload.Manifest.Status != "" {
			fmt.Fprintf(stdout, "Manifest status: %s\n", payload.Manifest.Status)
		}
		if payload.Manifest.Mode != "" {
			fmt.Fprintf(stdout, "Mode: %s\n", payload.Manifest.Mode)
		}
		if payload.Manifest.Transport != "" {
			fmt.Fprintf(stdout, "Transport: %s\n", payload.Manifest.Transport)
		}
		if !payload.Manifest.CreatedAt.IsZero() {
			fmt.Fprintf(stdout, "Created: %s\n", payload.Manifest.CreatedAt.Format(time.RFC3339))
		}
	} else if payload.ManifestError != "" {
		fmt.Fprintf(stdout, "Manifest: error: %s\n", payload.ManifestError)
	} else {
		fmt.Fprintf(stdout, "Manifest: missing (legacy run; expected at %s)\n", payload.ManifestPath)
	}
	if payload.RunError != "" {
		fmt.Fprintf(stdout, "Run: error: %s\n", payload.RunError)
		return
	}
	if payload.Run != nil {
		fmt.Fprintln(stdout)
		printRunDetail(stdout, *payload.Run)
	}
}

func sessionStatus(session sessionstore.Session) string {
	if manifest, err := runmanifest.Load(session.WorkspaceRoot, session.RunID); err == nil && manifest.Status != "" {
		return manifest.Status
	}
	if session.Terminal {
		return "terminal"
	}
	return "active"
}

func sessionStartedTimestamp(session sessionstore.Session) string {
	if session.CreatedAt.IsZero() {
		return "-"
	}
	return session.CreatedAt.Format(time.RFC3339)
}

func sessionTimestamp(session sessionstore.Session) string {
	value := session.LastEventAt
	if value.IsZero() {
		value = session.UpdatedAt
	}
	if value.IsZero() {
		value = session.CreatedAt
	}
	if value.IsZero() {
		return "-"
	}
	return value.Format(time.RFC3339)
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
	if err := resumePendingConsensusSignoffs(*root, run, stdout); err != nil {
		fmt.Fprintf(stderr, "resume failed: %v\n", err)
		return 1
	}
	if refreshed, err := runstate.ResolveRun(*root, fs.Arg(0)); err == nil {
		run = refreshed
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
	resumeKill, resumeLiveness := reattachSeams(run.RunDir)
	if err := tui.RunLive(tui.LiveOptions{
		Status:           status,
		Idea:             idea,
		Participants:     run.Participants,
		RunID:            run.RunID,
		RunDir:           run.RunDir,
		Resume:           true,
		KillAgent:        resumeKill,     // durable kill across the restart
		Liveness:         resumeLiveness, // RUN vs STALE badge
		Root:             *root,
		ReattachKill:     func(runID string) tui.KillAgentFunc { k, _ := reattachSeams(runDirFor(*root, runID)); return k },
		ReattachLiveness: func(runID string) tui.LivenessFunc { _, l := reattachSeams(runDirFor(*root, runID)); return l },
	}); err != nil {
		fmt.Fprintf(stderr, "resume tui failed: %v\n", err)
		return 1
	}
	return 0
}

func runContinue(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("continue", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	jsonOut := fs.Bool("json", false, "print unstable JSON output")
	auto := fs.Bool("auto", false, "execute the next action (auto-drive) instead of only printing it")
	noImplement := fs.Bool("no-implement", false, "stop the auto-driver at FINAL.md (skip code-writing implementation/fix-up phases)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley continue [--dir DIR] [--json] [--auto] RUN_OR_IDEA")
		return 2
	}

	run, err := runstate.ResolveRun(*root, fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "continue failed: %v\n", err)
		return 1
	}
	if *auto {
		return continueAuto(ctx, *root, run, *noImplement, stdout, stderr)
	}
	steers, _ := steer.List(run.RunDir)
	if *jsonOut {
		payload := struct {
			Run     runstate.RunSummary  `json:"run"`
			Actions []runplan.NextAction `json:"actions"`
			Steers  []steer.Queued       `json:"steers,omitempty"`
		}{Run: run, Actions: run.NextActions, Steers: steers}
		return printJSON(stdout, payload, stderr)
	}

	fmt.Fprintf(stdout, "Run: %s\n", run.RunID)
	fmt.Fprintf(stdout, "Idea: %s\n", valueOr(run.IdeaSlug, "unknown"))
	fmt.Fprintf(stdout, "State: %s\n", displayRunState(run))
	printNextActions(stdout, run)
	printQueuedSteers(stdout, steers)
	return 0
}

// continueAuto reconstructs the driver for an existing run and ticks it forward
// (consensus D9): `parley continue --auto` EXECUTES the next action instead of
// only printing it. Reuses the same gates/adapters as `parley run --auto`.
func continueAuto(ctx context.Context, root string, run runstate.RunSummary, noImplement bool, stdout, stderr io.Writer) int {
	if run.IdeaSlug == "" || run.IdeaSlug == "unknown" {
		fmt.Fprintln(stderr, "continue --auto failed: run has no idea slug")
		return 1
	}
	discovered, err := discoverConfigured(ctx, root)
	if err != nil {
		fmt.Fprintf(stderr, "agent config failed: %v\n", err)
		return 1
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", run.IdeaSlug)
	idea := protocol.IdeaStatus{Slug: run.IdeaSlug, Path: ideaDir, Participants: run.Participants}
	runOpts := runner.Options{
		Root:   root,
		RunID:  run.RunID,
		Idea:   idea,
		Agents: discovered,
		Store:  store.New(run.RunDir),
	}
	d := driver.New(driver.Config{
		IdeaDir:           ideaDir,
		IdeaSlug:          run.IdeaSlug,
		Participants:      run.Participants,
		RunDir:            run.RunDir,
		Root:              root,
		Events:            runOpts.Store,
		CrossReviewRounds: driver.ReadCrossReviewRounds(ideaDir),
		Auto:              true,
		Consensus:         newDriverConsensusOps(root, run.IdeaSlug, ideaDir, run.Participants, discovered, stdout),
		Impl:              newDriverImplOps(runOpts, root, run.IdeaSlug, ideaDir, run.Participants, stdout),
		AutoImplement:     driver.ReadAutoImplement(ideaDir) && !noImplement,
		Out:               stdout,
	}, driver.NewRunnerAdapter(runOpts))
	if err := d.Run(ctx); err != nil {
		fmt.Fprintf(stderr, "continue --auto: %v\n", err)
		return 1
	}
	if updated, err := runstate.LoadRun(root, run.RunID); err == nil {
		registerWorkspaceSessions(root, []runstate.RunSummary{updated})
	}
	return 0
}

// runSteer queues a follow-up / steering instruction for an agent (or the whole
// deck) on an existing run. It records the steer as durable events; it does not
// execute delivery (the queued attempt runs when the agent next executes).
func runSteer(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("steer", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	agent := fs.String("agent", "", "target agent ID (omit to steer the whole deck)")
	jsonOut := fs.Bool("json", false, "print unstable JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() < 2 {
		fmt.Fprintln(stderr, "usage: parley steer [--dir DIR] [--agent AGENT] [--json] RUN_OR_IDEA -- TEXT...")
		return 2
	}

	run, err := runstate.ResolveRun(*root, fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "steer failed: %v\n", err)
		return 1
	}
	// The run id is a positional, so Go's flag parser stops before any trailing
	// "--"; accept it as an optional separator before the text.
	rest := fs.Args()[1:]
	if len(rest) > 0 && rest[0] == "--" {
		rest = rest[1:]
	}
	text := strings.TrimSpace(strings.Join(rest, " "))
	if text == "" {
		fmt.Fprintln(stderr, "steer failed: instruction text is required")
		fmt.Fprintln(stderr, "usage: parley steer [--dir DIR] [--agent AGENT] [--json] RUN_OR_IDEA -- TEXT...")
		return 2
	}
	agentID := strings.TrimSpace(*agent)
	target := steer.TargetDeck
	if agentID != "" {
		target = steer.TargetAgent
		if !contains(run.Participants, agentID) {
			fmt.Fprintf(stderr, "warning: %q is not a participant of run %s; queuing anyway\n", agentID, run.RunID)
		}
	}
	result, err := steer.Submit(run.RunDir, steer.Request{
		Target:    target,
		Agent:     agentID,
		Text:      text,
		CreatedBy: "cli",
		SegmentID: segmentForAgent(run, agentID),
	}, time.Now().UTC())
	if err != nil {
		fmt.Fprintf(stderr, "steer failed: %v\n", err)
		return 1
	}
	if *jsonOut {
		return printJSON(stdout, result, stderr)
	}
	if target == steer.TargetAgent {
		fmt.Fprintf(stdout, "Recorded %s for %s (queued; auto-exec is not wired up yet).\n", result.ID, agentID)
	} else {
		fmt.Fprintf(stdout, "Recorded %s for the deck (queued; auto-exec is not wired up yet).\n", result.ID)
	}
	return 0
}

func segmentForAgent(run runstate.RunSummary, agentID string) string {
	if agentID == "" {
		return ""
	}
	for _, a := range run.State.Agents {
		if a.ID == agentID {
			return a.Segment
		}
	}
	return ""
}

func printQueuedSteers(stdout io.Writer, steers []steer.Queued) {
	if len(steers) == 0 {
		return
	}
	fmt.Fprintf(stdout, "Queued steers (%d):\n", len(steers))
	for _, s := range steers {
		target := s.Target
		if s.Agent != "" {
			target = s.Agent
		}
		fmt.Fprintf(stdout, "  %s  %-8s  %-7s  %s\n", s.ID, target, s.Status, truncateForList(s.Text))
	}
}

func truncateForList(text string) string {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
	if len(text) > 60 {
		return text[:57] + "..."
	}
	return text
}

func resumePendingConsensusSignoffs(root string, run runstate.RunSummary, stdout io.Writer) error {
	if run.Mode != "consensus-signoff" {
		return nil
	}
	events, err := store.New(run.RunDir).Load()
	if err != nil {
		return err
	}
	review := false
	for _, event := range events {
		if event.Type != "run.created" {
			continue
		}
		if value, ok := event.Data["review"].(bool); ok {
			review = value
		}
	}
	summary, err := consensus.Status(root, run.IdeaSlug, review)
	if err != nil {
		return err
	}
	if len(summary.Errors) > 0 {
		return fmt.Errorf("target consensus is malformed: %s", strings.Join(summary.Errors, "; "))
	}
	completed := map[string]bool{}
	var pending []store.Event
	for _, event := range events {
		agent := eventAgent(event)
		if agent == "" {
			continue
		}
		switch event.Type {
		case "agent.handoff.completed":
			completed[agent] = true
		case "agent.handoff.pending", "agent.handoff.started":
			pending = append(pending, event)
		}
	}
	signed := map[string]bool{}
	for _, signoff := range summary.Signoffs {
		signed[signoff.Agent] = true
	}
	runStore := store.New(run.RunDir)
	for _, event := range pending {
		agent := eventAgent(event)
		if completed[agent] || !signed[agent] {
			continue
		}
		if err := validateResumedConsensusHandoff(event, summary); err != nil {
			return err
		}
		artifact := ""
		if value, ok := event.Data["artifact"].(string); ok {
			artifact = value
		}
		if err := runStore.Append(store.Event{
			Time: time.Now().UTC(),
			Type: "agent.handoff.completed",
			Data: map[string]any{
				"agent":       agent,
				"artifact":    artifact,
				"launch_mode": event.Data["launch_mode"],
			},
		}); err != nil {
			return err
		}
		completed[agent] = true
		fmt.Fprintf(stdout, "Validated pending signoff from %s.\n", agent)
	}
	return nil
}

func validateResumedConsensusHandoff(event store.Event, summary consensus.Summary) error {
	agent := eventAgent(event)
	artifact, _ := event.Data["artifact"].(string)
	if agent == "" || artifact == "" {
		return fmt.Errorf("pending handoff event is missing agent or artifact")
	}
	beforeLen, ok := eventDataInt(event.Data, "before_len")
	if !ok || beforeLen < 0 {
		return fmt.Errorf("pending handoff for %s is missing before_len", agent)
	}
	beforeHash, _ := event.Data["before_sha256"].(string)
	if beforeHash == "" {
		return fmt.Errorf("pending handoff for %s is missing before_sha256", agent)
	}
	data, err := os.ReadFile(artifact)
	if err != nil {
		return err
	}
	if len(data) < beforeLen {
		return fmt.Errorf("%s changed existing consensus content before %s signoff", agent, agent)
	}
	beforeRaw := string(data[:beforeLen])
	afterRaw := string(data)
	if got := sha256Hex(beforeRaw); got != beforeHash {
		return fmt.Errorf("%s changed existing consensus content before append-only suffix", agent)
	}
	if err := validateAppendOnlyContent(beforeRaw, afterRaw, agent); err != nil {
		return err
	}
	var found bool
	for _, signoff := range summary.Signoffs {
		if signoff.Agent != agent {
			continue
		}
		found = true
		status, err := consensus.CanonicalStatus(signoff.Status)
		if err != nil {
			return err
		}
		if status == consensus.StatusBlock {
			return fmt.Errorf("%s appended BLOCK signoff", agent)
		}
	}
	if !found {
		return fmt.Errorf("%s did not append a signoff", agent)
	}
	return nil
}

func eventDataInt(data map[string]any, key string) (int, bool) {
	switch value := data[key].(type) {
	case int:
		return value, true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}

func eventAgent(event store.Event) string {
	if event.Data == nil {
		return ""
	}
	value, _ := event.Data["agent"].(string)
	return value
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

func printIdeaDetail(stdout io.Writer, idea protocol.IdeaStatus) {
	fmt.Fprintf(stdout, "Idea: %s\n", idea.Slug)
	fmt.Fprintf(stdout, "Status: %s\n", idea.Status)
	fmt.Fprintf(stdout, "Participants: %s\n", strings.Join(idea.Participants, ","))
}

func printConsensusIfPresent(stdout io.Writer, root, ideaSlug string) {
	if summary, err := consensus.Status(root, ideaSlug, false); err == nil {
		printConsensusSummary(stdout, summary)
	} else if !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(stdout, "Consensus: error: %v\n", err)
	}
	if summary, err := consensus.Status(root, ideaSlug, true); err == nil {
		printConsensusSummary(stdout, summary)
	} else if !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(stdout, "Review consensus: error: %v\n", err)
	}
}

func printConsensusSummary(stdout io.Writer, summary consensus.Summary) {
	label := "Consensus"
	if summary.Review {
		label = "Review consensus"
	}
	fmt.Fprintf(stdout, "%s: %s\n", label, summary.Triage)
	fmt.Fprintf(stdout, "Path: %s\n", summary.Path)
	if len(summary.Missing) > 0 {
		fmt.Fprintf(stdout, "Missing signoffs: %s\n", strings.Join(summary.Missing, ","))
	}
	if len(summary.Errors) > 0 {
		fmt.Fprintln(stdout, "Errors:")
		for _, errText := range summary.Errors {
			fmt.Fprintf(stdout, "  %s\n", errText)
		}
	}
	if len(summary.Signoffs) == 0 {
		fmt.Fprintln(stdout, "Signoffs: none")
		return
	}
	fmt.Fprintln(stdout, "Signoffs:")
	for _, signoff := range summary.Signoffs {
		fmt.Fprintf(stdout, "  %-10s %s", signoff.Agent, signoff.Status)
		if signoff.Notes != "" {
			fmt.Fprintf(stdout, " — %s", signoff.Notes)
		}
		fmt.Fprintln(stdout)
	}
}

func consensusTriageLabel(root, ideaSlug string) string {
	summary, err := consensus.Status(root, ideaSlug, false)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		}
		return "  consensus=error"
	}
	if len(summary.Errors) > 0 {
		return "  consensus=error"
	}
	return "  consensus=" + summary.Triage
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
	if len(run.NextActions) > 0 {
		if command := runaction.Command(run.NextActions[0], run.RunID, run.IdeaSlug); command != "" {
			return command
		}
		return fmt.Sprintf("parley continue %s", run.RunID)
	}
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

func printNextActions(stdout io.Writer, run runstate.RunSummary) {
	if len(run.NextActions) == 0 {
		fmt.Fprintln(stdout, "Next actions: none")
		return
	}
	fmt.Fprintf(stdout, "Recommended: %s\n", actionSummary(run.NextActions[0]))
	if command := runaction.Command(run.NextActions[0], run.RunID, run.IdeaSlug); command != "" {
		fmt.Fprintf(stdout, "Command: %s\n", command)
	}
	fmt.Fprintln(stdout, "Next actions:")
	for _, action := range run.NextActions {
		fmt.Fprintf(stdout, "  %s  kind=%s  risk=%s", action.ID, action.Kind, valueOr(action.Risk, "-"))
		if action.AgentID != "" {
			fmt.Fprintf(stdout, "  agent=%s", action.AgentID)
		}
		if action.ArtifactPath != "" {
			fmt.Fprintf(stdout, "  artifact=%s", action.ArtifactPath)
		}
		if action.RequiresYes {
			fmt.Fprint(stdout, "  requires=yes")
		}
		fmt.Fprintf(stdout, "\n      %s\n", actionSummary(action))
		if command := runaction.Command(action, run.RunID, run.IdeaSlug); command != "" {
			fmt.Fprintf(stdout, "      command: %s\n", command)
		}
	}
}

func actionSummary(action runplan.NextAction) string {
	return valueOr(action.Summary, action.Kind)
}

func findIdeaStatus(status protocol.WorkspaceStatus, slug string) (protocol.IdeaStatus, bool) {
	for _, idea := range status.Ideas {
		if idea.Slug == slug {
			return idea, true
		}
	}
	return protocol.IdeaStatus{}, false
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
	auto := fs.Bool("auto", true, "auto-drive the protocol past round-01 (default; use --no-auto to disable)")
	noAuto := fs.Bool("no-auto", false, "disable auto-drive: stop after round-01 and advance the protocol manually")
	noImplement := fs.Bool("no-implement", false, "stop the auto-driver at FINAL.md (skip code-writing implementation/fix-up phases)")
	noPreflight := fs.Bool("no-preflight", false, "skip the pre-idea readiness check (CI escape)")
	noPing := fs.Bool("no-ping", false, "preflight presence-only: skip the hosted-PONG roster ping (faster; for CI)")
	participantsFlag := fs.String("participants", "", "comma-separated agent IDs to run")
	yes := fs.Bool("yes", false, "launch selected agents without interactive confirmation")
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	// Auto-drive is on by default; --no-auto is the explicit opt-out.
	if *noAuto {
		*auto = false
	}
	task := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if task == "" {
		fmt.Fprintln(stderr, "usage: parley run [--no-tui] [--no-auto] [--participants AGENTS] [--yes] TASK")
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
		fmt.Fprintln(stderr, "no installed agents found; run `parley agents list` to inspect configuration")
		return 1
	}
	// Pre-idea readiness check (preflight §9.0): runs AFTER discovery and BEFORE
	// runcontrol.Create so a pending gate stops here with no half-open idea. It
	// evaluates the EXACT selected participant set (so --participants cannot bypass
	// the §1 non-solo hard-stop). Unattended (--auto, no TTY) hard-stops on a gate
	// and never reads stdin; attended prints the gate + confirm command and stops.
	// --no-preflight is the CI escape. Confirmed (--yes) exclusions are recorded in
	// the created idea.
	var preflightExcluded []string
	if !*noPreflight {
		code, excluded, stop := runTaskPreflight(ctx, *root, discovered, participants, attendedRun(*auto, *yes), *noPing, *yes, stdout, stderr)
		if stop {
			return code
		}
		preflightExcluded = excluded
	}
	if !*auto && !*yes && !confirmLaunch(os.Stdin, stdout, participants) {
		fmt.Fprintln(stdout, "No run started. Use `--yes` or `--auto` to launch without an interactive confirmation prompt.")
		return 0
	}
	created, err := runcontrol.Create(runcontrol.CreateOptions{
		Root:         *root,
		Task:         task,
		Participants: participants,
		Excluded:     preflightExcluded,
		Discovered:   discovered,
		Auto:         *auto,
	})
	if err != nil {
		fmt.Fprintf(stderr, "run create failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Created idea %s and run %s\n", created.Idea.Slug, created.RunID)
	fmt.Fprintf(stdout, "Starting round-01 with participants: %s\n", strings.Join(participants, ", "))
	runOpts := created.RunOptions
	if *noTUI {
		runCtx, cancelRun := context.WithCancel(ctx)
		defer cancelRun()
		if *auto {
			runcontrol.StartAutoAnswerer(runCtx, created.RunDir)
		}
		results := runner.RunRoundOne(runCtx, runOpts)
		failed := printRunResults(stdout, results)
		if run, err := runstate.LoadRun(*root, created.RunID); err == nil {
			registerWorkspaceSessions(*root, []runstate.RunSummary{run})
		}
		// Auto-advance past round-01 under --auto (any transport as of 1.27.0;
		// the driver advances canonical artifacts, not PR/MR branches). Round-01
		// must have succeeded. --no-auto keeps the one-shot behavior.
		if *auto && !failed {
			d := driver.New(driver.Config{
				IdeaDir:           created.Idea.Path,
				IdeaSlug:          created.Idea.Slug,
				Participants:      created.Idea.Participants,
				RunDir:            created.RunDir,
				Root:              *root,
				Events:            runOpts.Store,
				CrossReviewRounds: driver.ReadCrossReviewRounds(created.Idea.Path),
				Auto:              *auto,
				Consensus:         newDriverConsensusOps(*root, created.Idea.Slug, created.Idea.Path, created.Idea.Participants, discovered, stdout),
				Impl:              newDriverImplOps(runOpts, *root, created.Idea.Slug, created.Idea.Path, created.Idea.Participants, stdout),
				AutoImplement:     driver.ReadAutoImplement(created.Idea.Path) && !*noImplement,
				Out:               stdout,
			}, driver.NewRunnerAdapter(runOpts))
			if err := d.Run(runCtx); err != nil {
				fmt.Fprintf(stderr, "driver: %v\n", err)
			}
			if run, err := runstate.LoadRun(*root, created.RunID); err == nil {
				registerWorkspaceSessions(*root, []runstate.RunSummary{run})
			}
		}
		cancelRun()
		if failed {
			return 1
		}
		return 0
	}

	runCtx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()
	handle := runner.RunRoundOneAsync(runCtx, runOpts)
	// startAutoDrive advances the protocol past round-01 while the live TUI shows
	// it (cross-review → consensus → finalize → opted-in implementation). It is
	// idempotent (sync.Once): the auto default calls it now; the `/run` command
	// lets a --no-auto run kick it on demand. The driver emits run.phase/agent
	// events the TUI renders; its output is io.Discard so it never corrupts the
	// render, and quitting the TUI cancels runCtx, which stops the driver.
	var driveOnce sync.Once
	startAutoDrive := func() {
		driveOnce.Do(func() {
			runcontrol.StartAutoAnswerer(runCtx, created.RunDir)
			go func() {
				select {
				case <-runCtx.Done():
					return
				case <-handle.Done():
				}
				if anyRunFailed(handle.Results()) {
					return
				}
				d := driver.New(driver.Config{
					IdeaDir:           created.Idea.Path,
					IdeaSlug:          created.Idea.Slug,
					Participants:      created.Idea.Participants,
					RunDir:            created.RunDir,
					Root:              *root,
					Events:            runOpts.Store,
					CrossReviewRounds: driver.ReadCrossReviewRounds(created.Idea.Path),
					Auto:              true,
					Consensus:         newDriverConsensusOps(*root, created.Idea.Slug, created.Idea.Path, created.Idea.Participants, discovered, io.Discard),
					Impl:              newDriverImplOps(runOpts, *root, created.Idea.Slug, created.Idea.Path, created.Idea.Participants, io.Discard),
					AutoImplement:     driver.ReadAutoImplement(created.Idea.Path) && !*noImplement,
					Out:               io.Discard,
				}, driver.NewRunnerAdapter(runOpts))
				if err := d.Run(runCtx); err != nil {
					_ = runOpts.Store.Append(store.Event{Time: time.Now().UTC(), Type: "driver.error", Data: map[string]any{"error": err.Error()}})
				}
				if run, err := runstate.LoadRun(*root, created.RunID); err == nil {
					registerWorkspaceSessions(*root, []runstate.RunSummary{run})
				}
			}()
		})
	}
	if *auto {
		startAutoDrive()
	}
	workspaceStatus.Ideas = []protocol.IdeaStatus{created.Idea}
	steerFn, killFn, livenessFn := liveSteerKillSeams(runCtx, handle)
	if err := tui.RunLive(tui.LiveOptions{
		Status:           workspaceStatus,
		Idea:             created.Idea,
		Participants:     participants,
		RunID:            created.RunID,
		RunDir:           handle.RunDir,
		Done:             handle.Done(),
		Cancel:           cancelRun,
		SubmitSteer:      steerFn,
		KillAgent:        killFn,
		Liveness:         livenessFn,
		ReattachKill:     func(runID string) tui.KillAgentFunc { k, _ := reattachSeams(runDirFor(*root, runID)); return k },
		ReattachLiveness: func(runID string) tui.LivenessFunc { _, l := reattachSeams(runDirFor(*root, runID)); return l },
		StartAutoDrive:   startAutoDrive,
		Root:             *root,
		// `parley run` watches one run; start new ideas from `parley tui` (N).
	}); err != nil {
		cancelRun()
		results := handle.Wait()
		printRunResults(stdout, results)
		fmt.Fprintf(stderr, "tui failed: %v\n", err)
		return 1
	}
	results := handle.Wait()
	if run, err := runstate.LoadRun(*root, created.RunID); err == nil {
		registerWorkspaceSessions(*root, []runstate.RunSummary{run})
	}
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
			if err := tui.RunInit(root, results); err != nil {
				fmt.Fprintf(stderr, "tui failed: %v\n", err)
				return 1
			}
			// Re-read after the init wizard quits. If the user initialized the
			// workspace, fall through to the unified Home view; if they quit
			// without initializing, there is nothing more to show.
			status, err = protocol.ReadWorkspaceStatus(root)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return 0
				}
				fmt.Fprintf(stderr, "tui failed: %v\n", err)
				return 1
			}
		} else {
			fmt.Fprintf(stderr, "tui failed: %v\n", err)
			return 1
		}
	}

	runs, err := runstate.ListRuns(root)
	if err != nil {
		fmt.Fprintf(stderr, "tui failed: %v\n", err)
		return 1
	}
	registerWorkspaceSessions(root, runs)

	// Unified TUI: open to Home; N launches a new run (all available agents)
	// through the same runner path used by `parley run`. Detaching the TUI
	// (/quit, esc) does NOT cancel in-flight N-launched runs — only the active
	// run's Cancel (ctrl+c inside the TUI) does. On exit we wait for any
	// still-running launched runs so they finish and record their sessions
	// instead of being abandoned when the process returns.
	var reaper launchReaper
	if err := tui.RunLive(tui.LiveOptions{
		Home:   true,
		Root:   root,
		Status: status,
		Start:  newLaunchFunc(ctx, root, results, &reaper),
		// Opening an existing run (/open) is observational, but it can still kill
		// that run's agents via their persisted process identity, and show RUN/STALE.
		ReattachKill:     func(runID string) tui.KillAgentFunc { k, _ := reattachSeams(runDirFor(root, runID)); return k },
		ReattachLiveness: func(runID string) tui.LivenessFunc { _, l := reattachSeams(runDirFor(root, runID)); return l },
	}); err != nil {
		fmt.Fprintf(stderr, "tui failed: %v\n", err)
		return 1
	}
	// The parent context is still live here, so a real ctrl+c (SIGINT, now that
	// the TUI released the terminal) still aborts a run we are waiting on.
	reaper.waitForActive(stdout)
	return 0
}

// launchReaper waits for TUI-launched runs that are still in flight when the TUI
// detaches, so they complete and record their sessions rather than being killed
// when the command returns and cancels the parent context.
type launchReaper struct {
	wg       sync.WaitGroup
	inFlight atomic.Int64
}

// track runs fn (the wait+register reap) in the background and counts it as
// in-flight until it returns.
func (r *launchReaper) track(fn func()) {
	r.wg.Add(1)
	r.inFlight.Add(1)
	go func() {
		defer r.wg.Done()
		defer r.inFlight.Add(-1)
		fn()
	}()
}

// waitForActive blocks until every tracked run has finished, printing a hint
// first if any are still running so the now-restored terminal explains the wait.
func (r *launchReaper) waitForActive(out io.Writer) {
	if n := r.inFlight.Load(); n > 0 {
		fmt.Fprintf(out, "Waiting for %d in-progress run(s) to finish (ctrl+c to abort)...\n", n)
	}
	r.wg.Wait()
}

// liveSteerKillSeams builds the injected SubmitSteer/KillAgent seams that let the
// TUI execute a steer (a fresh single-agent attempt whose stdout is the reply)
// and kill one agent, against a live runner.Handle. internal/tui never imports
// internal/runner; these plain func seams are the only bridge.
func liveSteerKillSeams(runCtx context.Context, handle *runner.Handle) (tui.SteerFunc, tui.KillAgentFunc, tui.LivenessFunc) {
	runDir := handle.RunDir
	submit := func(req tui.SteerRequest) (tui.SteerResult, error) {
		// Always keep the durable audit trail (steer.requested/delivered) first.
		rec, err := steer.Submit(runDir, steer.Request{
			Target:    req.Target,
			Agent:     req.AgentID,
			Text:      req.Text,
			CreatedBy: "tui",
		}, time.Now().UTC())
		if err != nil {
			return tui.SteerResult{}, err
		}
		out := tui.SteerResult{ID: rec.ID, Status: "recorded"}
		// Only agent-targeted steers execute; deck steers stay record-only for now.
		if req.Target == steer.TargetAgent {
			att, aerr := handle.RunSteerAttempt(runCtx, runner.SteerAttemptRequest{
				AgentID: req.AgentID,
				Text:    req.Text,
				SteerID: rec.ID,
			})
			if aerr != nil {
				return out, aerr
			}
			if !att.Accepted {
				out.Status = "rejected"
				return out, fmt.Errorf("%s", att.Message)
			}
			out.Status = att.Status
			out.SegmentID = att.SegmentID
			out.StdoutPath = att.StdoutPath
		}
		return out, nil
	}
	kill := killOutcome(func(agentID string) runner.KillResult { return handle.KillAgent(agentID) })
	liveness := func(agentID string) string { return runner.AgentLivenessAt(runDir, agentID) }
	return submit, kill, liveness
}

// killOutcome adapts a KillResult-returning kill into the TUI seam: a failed
// operation returns a Go error (shown red); success returns its outcome message.
func killOutcome(do func(agentID string) runner.KillResult) tui.KillAgentFunc {
	return func(agentID string) (string, error) {
		res := do(agentID)
		if res.Failed {
			return "", fmt.Errorf("%s", res.Message)
		}
		if res.Message != "" {
			return res.Message, nil
		}
		if res.Killed {
			return "killed " + agentID, nil
		}
		return "", nil
	}
}

// runDirFor is the on-disk directory for a run id (mirrors runner.Handle.RunDir).
func runDirFor(root, runID string) string {
	return filepath.Join(root, protocol.DeckDir, "runs", runID)
}

// reattachSeams builds the durable kill + liveness seams for a run that has no
// live in-process handle (resume / open): they operate purely on the run dir's
// persisted process identities.
func reattachSeams(runDir string) (tui.KillAgentFunc, tui.LivenessFunc) {
	kill := killOutcome(func(agentID string) runner.KillResult { return runner.DurableKillAt(runDir, agentID) })
	liveness := func(agentID string) string { return runner.AgentLivenessAt(runDir, agentID) }
	return kill, liveness
}

// newLaunchFunc returns a tui.LaunchFunc that starts a new round-01 run with all
// available agents through the existing runner path and returns a handle to
// attach to (no parallel engine). The run is reaped through the reaper so its
// session is recorded and the command waits for it on exit; it is canceled only
// via the returned Cancel (ctrl+c on the attached run), never by TUI detach.
func newLaunchFunc(ctx context.Context, root string, discovered []agents.Discovery, reaper *launchReaper) tui.LaunchFunc {
	return func(req tui.LaunchRequest) (tui.LaunchResult, error) {
		participants := installedAgentIDs(discovered)
		if len(participants) == 0 {
			return tui.LaunchResult{}, fmt.Errorf("no installed agents found")
		}
		created, err := runcontrol.Create(runcontrol.CreateOptions{
			Root:         root,
			Task:         req.Task,
			Participants: participants,
			Discovered:   discovered,
		})
		if err != nil {
			return tui.LaunchResult{}, err
		}
		runCtx, cancelRun := context.WithCancel(ctx)
		handle := runner.RunRoundOneAsync(runCtx, created.RunOptions)
		reaper.track(func() {
			handle.Wait()
			if run, err := runstate.LoadRun(root, created.RunID); err == nil {
				registerWorkspaceSessions(root, []runstate.RunSummary{run})
			}
		})
		steerFn, killFn, livenessFn := liveSteerKillSeams(runCtx, handle)
		return tui.LaunchResult{
			Idea:         created.Idea,
			Participants: participants,
			RunID:        created.RunID,
			RunDir:       handle.RunDir,
			Done:         handle.Done(),
			Cancel:       cancelRun,
			SubmitSteer:  steerFn,
			KillAgent:    killFn,
			Liveness:     livenessFn,
		}, nil
	}
}

func registerWorkspaceSessions(root string, runs []runstate.RunSummary) {
	sessionStore, err := sessionstore.Default()
	if err != nil {
		return
	}
	now := time.Now().UTC()
	for _, run := range runs {
		if run.RunID == "" || run.Error != "" {
			continue
		}
		createdAt := run.LastEventAt
		if createdAt.IsZero() {
			createdAt = now
		}
		_ = sessionStore.Upsert(sessionstore.Session{
			WorkspaceRoot: root,
			RunID:         run.RunID,
			IdeaSlug:      run.IdeaSlug,
			Task:          run.Task,
			Participants:  append([]string(nil), run.Participants...),
			CreatedAt:     createdAt,
			UpdatedAt:     now,
			LastEventAt:   run.LastEventAt,
			Terminal:      run.Terminal,
		})
	}
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
		if result.Found && result.ID != "gemini" {
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

// anyRunFailed reports whether round-01 should block auto-drive: no participants,
// or any non-successful result (mirrors printRunResults' failure verdict without
// printing, for the TUI auto-drive goroutine).
func anyRunFailed(results []runner.Result) bool {
	if len(results) == 0 {
		return true
	}
	for _, result := range results {
		if !result.Success() {
			return true
		}
	}
	return false
}

func printRunResults(stdout io.Writer, results []runner.Result) bool {
	if len(results) == 0 {
		fmt.Fprintln(stdout, "No participants were selected for round-01.")
		return true
	}

	failed := false
	for _, result := range results {
		switch {
		case result.Warning != "" && result.ExitError == "":
			fmt.Fprintf(stdout, "  %-8s warning: %s\n", result.AgentID, result.Warning)
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
