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
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/pipeline"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
)

// launchBlockRound runs one cooperation round for a block in its engine
// workspace via the existing runner (round 1 -> RunRoundOne, later ->
// RunRound). Shared by `pipeline run-block` and `pipeline auto`.
func launchBlockRound(ctx context.Context, root, deck, slug, blockID string, participants []string, discovered []agents.Discovery, round int) []runner.Result {
	blockWS := pipeline.BlockWorkspace(deck, slug, blockID)
	runID := fmt.Sprintf("pipe-%s-r%02d-%s", blockID, round, time.Now().UTC().Format("20060102T150405.000000Z"))
	idea := protocol.IdeaStatus{Slug: slug + "__" + blockID, Path: blockWS, Participants: participants}
	opts := runner.Options{
		Root:    root,
		RunID:   runID,
		Idea:    idea,
		Agents:  discovered,
		Round:   round,
		Timeout: 30 * time.Minute,
		Store:   store.New(filepath.Join(deck, "runs", runID)),
	}
	if round <= 1 {
		return runner.RunRoundOne(ctx, opts)
	}
	return runner.RunRound(ctx, opts)
}

func printPipelineUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: parley pipeline validate MANIFEST | start [--dir DIR] MANIFEST | status [--dir DIR] SLUG | run-block [--dir DIR] [--participants IDS] [--round N] [--yes] SLUG | continue [--dir DIR] SLUG | auto [--dir DIR] [--rounds N] [--max-blocks M] [--drafter ID] [--participants IDS] [--yes] SLUG | watch [--dir DIR] [--spec FILE] --signals FILE [--once] SLUG | execute [--dir DIR] [--provider P] [--dry-run] [--json] SLUG BLOCK CAPABILITY TARGET | record-effect [--dir DIR] --status STATUS [--external-ref REF] SLUG DIGEST | gate approve|reject [--dir DIR] [--by WHO] SLUG EDGE")
}

func runPipeline(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printPipelineUsage(stderr)
		return 2
	}
	switch args[0] {
	case "validate":
		return runPipelineValidate(args[1:], stdout, stderr)
	case "start":
		return runPipelineStart(args[1:], stdout, stderr)
	case "status":
		return runPipelineStatus(args[1:], stdout, stderr)
	case "run-block":
		return runPipelineRunBlock(ctx, args[1:], stdout, stderr)
	case "continue":
		return runPipelineContinue(args[1:], stdout, stderr)
	case "auto":
		return runPipelineAuto(ctx, args[1:], stdout, stderr)
	case "watch":
		return runPipelineWatch(args[1:], stdout, stderr)
	case "execute":
		return runPipelineExecute(args[1:], stdout, stderr)
	case "record-effect":
		return runPipelineRecordEffect(args[1:], stdout, stderr)
	case "gate":
		return runPipelineGate(args[1:], stdout, stderr)
	default:
		printPipelineUsage(stderr)
		return 2
	}
}

func selectProvider(name string) (pipeline.Provider, error) {
	switch name {
	case "", "vercel":
		return pipeline.VercelProvider{}, nil
	case "noop":
		return pipeline.NoopProvider{}, nil
	default:
		return nil, fmt.Errorf("unknown provider %q (want vercel|noop)", name)
	}
}

func findBlock(m pipeline.Manifest, id string) (pipeline.Block, bool) {
	for _, b := range m.Blocks {
		if b.ID == id {
			return b, true
		}
	}
	return pipeline.Block{}, false
}

func deckDirFor(root string) string { return filepath.Join(root, protocol.DeckDir) }

func runPipelineValidate(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley pipeline validate MANIFEST")
		return 2
	}
	m, err := pipeline.ParseFile(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline validate failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "valid: pipeline %q, autonomy %s, transport %s, %d block(s)\n", m.IdeaSlug, m.Autonomy, m.Transport, len(m.Blocks))
	for i, b := range m.Blocks {
		fmt.Fprintf(stdout, "  %d. %-16s kind=%-14s stage=%s\n", i+1, b.ID, b.Kind, valueOrDash(b.Stage))
	}
	return 0
}

func runPipelineStart(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline start", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley pipeline start [--dir DIR] MANIFEST")
		return 2
	}
	data, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline start failed: %v\n", err)
		return 1
	}
	m, err := pipeline.Parse(data)
	if err != nil {
		fmt.Fprintf(stderr, "pipeline start failed: %v\n", err)
		return 1
	}
	deck := deckDirFor(*root)
	if err := os.MkdirAll(pipeline.PipelineDir(deck, m.IdeaSlug), 0o755); err != nil {
		fmt.Fprintf(stderr, "pipeline start failed: %v\n", err)
		return 1
	}
	if err := os.WriteFile(pipeline.ManifestPath(deck, m.IdeaSlug), data, 0o644); err != nil {
		fmt.Fprintf(stderr, "pipeline start failed: %v\n", err)
		return 1
	}
	now := time.Now()
	run := pipeline.NewPipelineRun(m, now)
	if _, err := pipeline.SeedBlockPrompt(deck, m, nil, m.Blocks[0], now); err != nil {
		fmt.Fprintf(stderr, "pipeline start failed: %v\n", err)
		return 1
	}
	if err := run.Save(pipeline.RunPath(deck, m.IdeaSlug)); err != nil {
		fmt.Fprintf(stderr, "pipeline start failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Started pipeline %q. First block %q seeded at %s\n", m.IdeaSlug, m.Blocks[0].ID, pipeline.BlockWorkspace(deck, m.IdeaSlug, m.Blocks[0].ID))
	fmt.Fprintf(stdout, "Run that block's cooperation rounds, then `parley pipeline continue --dir %s %s`.\n", *root, m.IdeaSlug)
	return 0
}

func runPipelineStatus(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline status", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley pipeline status [--dir DIR] SLUG")
		return 2
	}
	deck := deckDirFor(*root)
	run, err := pipeline.LoadRun(pipeline.RunPath(deck, fs.Arg(0)))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline status failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "pipeline %q: status=%s current=%s\n", run.PipelineSlug, run.Status, valueOrDash(run.CurrentBlock))
	fmt.Fprintf(stdout, "  completed: %s\n", valueOrDash(strings.Join(run.CompletedBlocks, ", ")))
	if run.PendingGate != "" {
		fmt.Fprintf(stdout, "  pending gate: %s\n", run.PendingGate)
	}
	return 0
}

func runPipelineContinue(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline continue", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley pipeline continue [--dir DIR] SLUG")
		return 2
	}
	deck := deckDirFor(*root)
	slug := fs.Arg(0)
	m, err := pipeline.ParseFile(pipeline.ManifestPath(deck, slug))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline continue failed: %v\n", err)
		return 1
	}
	run, err := pipeline.LoadRun(pipeline.RunPath(deck, slug))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline continue failed: %v\n", err)
		return 1
	}
	driver := pipeline.Driver{DeckDir: deck, Manifest: m, Now: time.Now()}
	res, err := driver.Advance(&run, blockCompleteFunc(deck, slug))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline continue failed: %v\n", err)
		return 1
	}
	if err := run.Save(pipeline.RunPath(deck, slug)); err != nil {
		fmt.Fprintf(stderr, "pipeline continue failed: %v\n", err)
		return 1
	}
	switch res.Action {
	case pipeline.ActionRunBlock:
		fmt.Fprintf(stdout, "run-block: block %q has not finished its rounds + consensus. %s\n", res.Block, res.Detail)
	case pipeline.ActionAwaitGate:
		fmt.Fprintf(stdout, "await-gate: %s\n  approve with: parley pipeline gate approve --dir %s %s %s\n", res.Detail, *root, slug, res.Edge)
	case pipeline.ActionSeededNext:
		fmt.Fprintf(stdout, "advanced: seeded next block %q (%s)\n", res.Block, res.Edge)
	case pipeline.ActionDone:
		fmt.Fprintf(stdout, "done: pipeline %q complete\n", slug)
	case pipeline.ActionRejected:
		fmt.Fprintf(stdout, "stopped: boundary gate rejected (%s)\n", res.Edge)
	}
	return 0
}

func runPipelineGate(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || (args[0] != "approve" && args[0] != "reject") {
		fmt.Fprintln(stderr, "usage: parley pipeline gate approve|reject [--dir DIR] [--by WHO] SLUG EDGE")
		return 2
	}
	approve := args[0] == "approve"
	fs := flag.NewFlagSet("pipeline gate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	by := fs.String("by", "user", "who is resolving the gate")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	if fs.NArg() != 2 {
		fmt.Fprintln(stderr, "usage: parley pipeline gate approve|reject [--dir DIR] [--by WHO] SLUG EDGE")
		return 2
	}
	deck := deckDirFor(*root)
	slug, edge := fs.Arg(0), fs.Arg(1)
	gate, ok, err := pipeline.LoadGate(deck, slug, edge)
	if err != nil {
		fmt.Fprintf(stderr, "pipeline gate failed: %v\n", err)
		return 1
	}
	if !ok {
		fmt.Fprintf(stderr, "pipeline gate failed: no gate %q for pipeline %q\n", edge, slug)
		return 1
	}
	gate.Resolve(approve, *by, time.Now())
	if err := pipeline.SaveGate(deck, gate); err != nil {
		fmt.Fprintf(stderr, "pipeline gate failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "gate %q for pipeline %q -> %s (by %s)\n", edge, slug, gate.Status, *by)
	return 0
}

// runPipelineRunBlock launches the current block's next cooperation round in
// its engine workspace using the existing runner. This is the live-engine
// wiring: round 1 uses RunRoundOne, later rounds use RunRound (cross-review).
// After rounds converge, draft+finalize consensus, then `pipeline continue`.
func runPipelineRunBlock(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline run-block", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	participantsFlag := fs.String("participants", "", "comma-separated agent IDs (default: manifest participants)")
	roundFlag := fs.Int("round", 0, "round to run (default: next pending round)")
	yes := fs.Bool("yes", false, "launch selected agents without interactive confirmation")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley pipeline run-block [--dir DIR] [--participants IDS] [--round N] [--yes] SLUG")
		return 2
	}
	deck := deckDirFor(*root)
	slug := fs.Arg(0)
	m, err := pipeline.ParseFile(pipeline.ManifestPath(deck, slug))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline run-block failed: %v\n", err)
		return 1
	}
	run, err := pipeline.LoadRun(pipeline.RunPath(deck, slug))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline run-block failed: %v\n", err)
		return 1
	}
	if run.CurrentBlock == "" {
		fmt.Fprintf(stderr, "pipeline run-block failed: pipeline %q has no current block (status %s)\n", slug, run.Status)
		return 1
	}
	blockWS := pipeline.BlockWorkspace(deck, slug, run.CurrentBlock)
	if _, err := os.Stat(filepath.Join(blockWS, "00-prompt.md")); err != nil {
		fmt.Fprintf(stderr, "pipeline run-block failed: block %q is not seeded (%v)\n", run.CurrentBlock, err)
		return 1
	}
	round := *roundFlag
	if round == 0 {
		round = nextBlockRound(blockWS)
	}

	discovered, err := discoverConfigured(ctx, *root)
	if err != nil {
		fmt.Fprintf(stderr, "agent config failed: %v\n", err)
		return 1
	}
	participantsArg := *participantsFlag
	if strings.TrimSpace(participantsArg) == "" {
		participantsArg = strings.Join(m.Participants, ",")
	}
	participants, err := selectedParticipantIDs(discovered, participantsArg)
	if err != nil {
		fmt.Fprintf(stderr, "participant selection failed: %v\n", err)
		return 1
	}
	if len(participants) == 0 {
		fmt.Fprintln(stderr, "no installed manifest participants found; capability halt (see COOPERATION.md §12.9)")
		return 1
	}
	if !*yes && !confirmLaunch(os.Stdin, stdout, participants) {
		fmt.Fprintln(stdout, "No round started. Use --yes to launch without confirmation.")
		return 0
	}

	fmt.Fprintf(stdout, "Running %s block %q round-%02d with: %s\n", slug, run.CurrentBlock, round, strings.Join(participants, ", "))
	results := launchBlockRound(ctx, *root, deck, slug, run.CurrentBlock, participants, discovered, round)
	if printRunResults(stdout, results) {
		return 1
	}
	blockIdeaSlug := slug + "__" + run.CurrentBlock
	fmt.Fprintf(stdout, "Round-%02d complete in %s. Run more rounds, then draft+finalize this block's consensus (parley consensus draft/request-signoffs/finalize on idea %s), then `parley pipeline continue --dir %s %s`. Or use `parley pipeline auto` to do all of it.\n", round, blockWS, blockIdeaSlug, *root, slug)
	return 0
}

// nextBlockRound returns the next round to run: one past the highest round that
// already has participant artifacts (an empty seeded round-01 -> 1).
func nextBlockRound(blockWS string) int {
	highest := 0
	for r := 1; r <= 50; r++ {
		dir := filepath.Join(blockWS, fmt.Sprintf("round-%02d", r))
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") && e.Name() != "_index.md" {
				highest = r
				break
			}
		}
	}
	return highest + 1
}

// runPipelineAuto is the unattended driver loop. For each deliberation block
// from the cursor it runs round-01 + cross-review rounds, drafts consensus,
// collects signoffs, and finalizes; then it advances across the boundary gate
// (auto-approving only low-risk non-production boundaries under auto-left
// autonomy, pausing at human gates). It stops at the first human gate, the
// first non-deliberation block, a non-ready consensus, or completion — never
// performing a production mutation.
func runPipelineAuto(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline auto", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	rounds := fs.Int("rounds", 1, "cross-review rounds per block in addition to round-01")
	maxBlocks := fs.Int("max-blocks", 20, "safety cap on blocks advanced in one invocation")
	drafter := fs.String("drafter", "", "consensus drafter agent ID (default: first participant)")
	participantsFlag := fs.String("participants", "", "comma-separated agent IDs (default: manifest participants)")
	yes := fs.Bool("yes", false, "launch hosted agents without confirmation")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley pipeline auto [--dir DIR] [--rounds N] [--max-blocks M] [--drafter ID] [--participants IDS] [--yes] SLUG")
		return 2
	}
	deck := deckDirFor(*root)
	slug := fs.Arg(0)

	for iter := 0; iter < *maxBlocks; iter++ {
		m, err := pipeline.ParseFile(pipeline.ManifestPath(deck, slug))
		if err != nil {
			fmt.Fprintf(stderr, "auto failed: %v\n", err)
			return 1
		}
		run, err := pipeline.LoadRun(pipeline.RunPath(deck, slug))
		if err != nil {
			fmt.Fprintf(stderr, "auto failed: %v\n", err)
			return 1
		}
		switch run.Status {
		case pipeline.StatusCompleted:
			fmt.Fprintf(stdout, "auto: pipeline %q complete.\n", slug)
			return 0
		case pipeline.StatusFailed:
			fmt.Fprintf(stdout, "auto: pipeline %q stopped (rejected gate %s).\n", slug, run.PendingGate)
			return 1
		}
		if run.CurrentBlock == "" {
			fmt.Fprintf(stderr, "auto failed: pipeline %q has no current block (status %s)\n", slug, run.Status)
			return 1
		}
		block, ok := findBlock(m, run.CurrentBlock)
		if !ok {
			fmt.Fprintf(stderr, "auto failed: current block %q not in manifest\n", run.CurrentBlock)
			return 1
		}

		done, err := blockCompleteFunc(deck, slug)(block)
		if err != nil {
			fmt.Fprintf(stderr, "auto failed: %v\n", err)
			return 1
		}
		// Action and implementation blocks are never auto-advanced past: their
		// side effects (deploy) and code-review/fix-up are explicit gated steps.
		if block.Kind == pipeline.KindAction {
			if !done {
				if code := autoDriveDeliberationBlock(ctx, *root, deck, slug, block, *participantsFlag, *drafter, *rounds, *yes, stdout, stderr); code != 0 {
					return code
				}
			}
			fmt.Fprintf(stdout, "auto: action block %q ready (plan finalized). status=needs_human_gate — run `parley pipeline execute [--dry-run] [--json] --dir %s %s %s <capability> <target>` (production-gated) + `record-effect`, then `parley pipeline continue` to advance. auto never auto-executes or advances past an action block.\n", block.ID, *root, slug, block.ID)
			return 0
		}
		if block.Kind == pipeline.KindImplementation {
			if !done {
				if code := autoDriveImplementationBlock(ctx, *root, deck, slug, block, *participantsFlag, *drafter, *yes, stdout, stderr); code != 0 {
					return code
				}
			}
			fmt.Fprintf(stdout, "auto: implementation block %q — implementation + review round-01 done. status=needs_artifact — draft review/consensus.md (machine-readable agreed_fixes) + fix-up, then `parley pipeline continue`.\n", block.ID)
			return 0
		}
		if !done {
			// deliberation or watcher: drive Phase 1-4 to a finalized artifact.
			if code := autoDriveDeliberationBlock(ctx, *root, deck, slug, block, *participantsFlag, *drafter, *rounds, *yes, stdout, stderr); code != 0 {
				return code
			}
			if block.Kind == pipeline.KindWatcher {
				fmt.Fprintf(stdout, "auto: watcher block %q spec finalized. Schedule `parley pipeline watch %s --signals <file>` (e.g. via cron) to evaluate SLOs.\n", block.ID, slug)
			}
		}

		driver := pipeline.Driver{DeckDir: deck, Manifest: m, Now: time.Now()}
		res, err := driver.Advance(&run, blockCompleteFunc(deck, slug))
		if err != nil {
			fmt.Fprintf(stderr, "auto failed: %v\n", err)
			return 1
		}
		if err := run.Save(pipeline.RunPath(deck, slug)); err != nil {
			fmt.Fprintf(stderr, "auto failed: %v\n", err)
			return 1
		}
		switch res.Action {
		case pipeline.ActionAwaitGate:
			fmt.Fprintf(stdout, "auto: paused at gate %s.\n  approve: parley pipeline gate approve --dir %s %s %s, then re-run `parley pipeline auto`.\n", res.Edge, *root, slug, res.Edge)
			return 0
		case pipeline.ActionRejected:
			fmt.Fprintf(stdout, "auto: gate %s rejected; stopped.\n", res.Edge)
			return 1
		case pipeline.ActionDone:
			fmt.Fprintf(stdout, "auto: pipeline %q complete (final block %q).\n", slug, res.Block)
			return 0
		case pipeline.ActionSeededNext:
			fmt.Fprintf(stdout, "auto: advanced to block %q.\n", res.Block)
		case pipeline.ActionRunBlock:
			fmt.Fprintf(stderr, "auto: block %q did not finalize after consensus; stopping to avoid a loop.\n", res.Block)
			return 1
		}
	}
	fmt.Fprintf(stderr, "auto: hit --max-blocks safety cap; re-run to continue.\n")
	return 1
}

// autoDriveDeliberationBlock runs the rounds + consensus + finalize for one
// deliberation block. Returns 0 on success, non-zero (an exit code) on a stop.
func autoDriveDeliberationBlock(ctx context.Context, root, deck, slug string, block pipeline.Block, participantsFlag, drafter string, rounds int, yes bool, stdout, stderr io.Writer) int {
	discovered, err := discoverConfigured(ctx, root)
	if err != nil {
		fmt.Fprintf(stderr, "auto: agent config failed: %v\n", err)
		return 1
	}
	pArg := participantsFlag
	if strings.TrimSpace(pArg) == "" {
		// block-local roster falls back to the manifest participants.
		if m, err := pipeline.ParseFile(pipeline.ManifestPath(deck, slug)); err == nil {
			pArg = strings.Join(m.Participants, ",")
		}
	}
	participants, err := selectedParticipantIDs(discovered, pArg)
	if err != nil {
		fmt.Fprintf(stderr, "auto: participant selection failed: %v\n", err)
		return 1
	}
	if len(participants) == 0 {
		fmt.Fprintln(stderr, "auto: no installed participants for this block (capability halt, §12.9)")
		return 1
	}
	by := drafter
	if by == "" {
		by = participants[0]
	}
	blockIdeaSlug := slug + "__" + block.ID

	first := nextBlockRound(pipeline.BlockWorkspace(deck, slug, block.ID))
	last := 1 + rounds
	for r := first; r <= last; r++ {
		fmt.Fprintf(stdout, "auto: block %q round-%02d (%s)\n", block.ID, r, strings.Join(participants, ", "))
		if printRunResults(stdout, launchBlockRound(ctx, root, deck, slug, block.ID, participants, discovered, r)) {
			fmt.Fprintf(stderr, "auto: round-%02d had failures; stopping.\n", r)
			return 1
		}
	}

	if _, err := consensus.Draft(root, blockIdeaSlug, consensus.DraftOptions{By: by}); err != nil {
		fmt.Fprintf(stderr, "auto: consensus draft failed: %v\n", err)
		return 1
	}
	if err := requestConsensusSignoffs(ctx, requestSignoffsOptions{Root: root, IdeaSlug: blockIdeaSlug, Yes: yes}, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "auto: request signoffs failed: %v\n", err)
		return 1
	}
	summary, err := consensus.Status(root, blockIdeaSlug, false)
	if err != nil {
		fmt.Fprintf(stderr, "auto: consensus status failed: %v\n", err)
		return 1
	}
	if summary.Triage != consensus.TriageReady && summary.Triage != consensus.TriageReserved {
		fmt.Fprintf(stdout, "auto: block %q consensus triage=%s (not finalizable); resolve reservations/blocks then re-run auto.\n", block.ID, summary.Triage)
		return 1
	}
	if _, _, err := consensus.Finalize(root, blockIdeaSlug, consensus.FinalizeOptions{By: by}); err != nil {
		fmt.Fprintf(stderr, "auto: finalize failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "auto: block %q finalized.\n", block.ID)
	return 0
}

// autoDriveImplementationBlock runs Phase 5 (implementation) and Phase 6
// review round-01 for an implementation block using the runner. It stops before
// review consensus + fix-up (human-guided). Returns 0 on success.
func autoDriveImplementationBlock(ctx context.Context, root, deck, slug string, block pipeline.Block, participantsFlag, drafter string, yes bool, stdout, stderr io.Writer) int {
	discovered, err := discoverConfigured(ctx, root)
	if err != nil {
		fmt.Fprintf(stderr, "auto: agent config failed: %v\n", err)
		return 1
	}
	pArg := participantsFlag
	if strings.TrimSpace(pArg) == "" {
		if m, err := pipeline.ParseFile(pipeline.ManifestPath(deck, slug)); err == nil {
			pArg = strings.Join(m.Participants, ",")
		}
	}
	participants, err := selectedParticipantIDs(discovered, pArg)
	if err != nil {
		fmt.Fprintf(stderr, "auto: participant selection failed: %v\n", err)
		return 1
	}
	if len(participants) == 0 {
		fmt.Fprintln(stderr, "auto: no installed participants for this block (capability halt, §12.9)")
		return 1
	}
	implementer := drafter
	if implementer == "" {
		implementer = participants[0]
	}
	blockWS := pipeline.BlockWorkspace(deck, slug, block.ID)
	idea := protocol.IdeaStatus{Slug: slug + "__" + block.ID, Path: blockWS}

	if _, statErr := os.Stat(filepath.Join(blockWS, "IMPLEMENTATION.md")); os.IsNotExist(statErr) {
		fmt.Fprintf(stdout, "auto: block %q Phase 5 implementation by %s\n", block.ID, implementer)
		ideaImpl := idea
		ideaImpl.Participants = []string{implementer}
		runID := fmt.Sprintf("pipe-%s-impl-%s", block.ID, time.Now().UTC().Format("20060102T150405.000000Z"))
		res := runner.RunImplementation(ctx, runner.Options{Root: root, RunID: runID, Idea: ideaImpl, Agents: discovered, Timeout: 30 * time.Minute, Store: store.New(filepath.Join(deck, "runs", runID))})
		if res.ExitError != "" || !res.ArtifactOK {
			fmt.Fprintf(stderr, "auto: implementation failed: %s\n", res.ExitError)
			return 1
		}
	}

	reviewers := make([]string, 0, len(participants))
	for _, p := range participants {
		if p != implementer {
			reviewers = append(reviewers, p)
		}
	}
	if len(reviewers) == 0 {
		fmt.Fprintf(stdout, "auto: block %q has no non-implementer reviewer; Phase 6 needs another agent (non-solo review).\n", block.ID)
		return 0
	}
	if _, statErr := os.Stat(filepath.Join(blockWS, "review", "round-01")); os.IsNotExist(statErr) {
		fmt.Fprintf(stdout, "auto: block %q Phase 6 review round-01 (%s)\n", block.ID, strings.Join(reviewers, ", "))
		ideaRev := idea
		ideaRev.Participants = reviewers
		runID := fmt.Sprintf("pipe-%s-rev01-%s", block.ID, time.Now().UTC().Format("20060102T150405.000000Z"))
		if printRunResults(stdout, runner.RunReviewRound(ctx, runner.Options{Root: root, RunID: runID, Idea: ideaRev, Round: 1, Agents: discovered, Timeout: 30 * time.Minute, Store: store.New(filepath.Join(deck, "runs", runID))})) {
			return 1
		}
	}
	return 0
}

// runPipelineWatch performs one monitoring evaluation pass (§12.11): it loads a
// structured watcher spec + current signal values, raises breaches, dedupes
// them against persisted records, opens a remediation idea ONLY for predeclared
// low-risk classes, and marks recovered breaches resolved. Periodic execution
// is the caller's cron/scheduler invoking this repeatedly (one pass each).
func runPipelineWatch(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline watch", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	specPath := fs.String("spec", "", "monitoring spec YAML (default: pipelines/<slug>/monitoring.yaml)")
	signalsPath := fs.String("signals", "", "JSON file of {signal: number} current values")
	fs.Bool("once", true, "evaluate a single pass (the only mode; schedule repeats via cron)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 || strings.TrimSpace(*signalsPath) == "" {
		fmt.Fprintln(stderr, "usage: parley pipeline watch [--dir DIR] [--spec FILE] --signals FILE [--once] SLUG")
		return 2
	}
	deck := deckDirFor(*root)
	slug := fs.Arg(0)
	spec := *specPath
	if spec == "" {
		spec = filepath.Join(pipeline.PipelineDir(deck, slug), "monitoring.yaml")
	}
	m, err := pipeline.LoadMonitoring(spec)
	if err != nil {
		fmt.Fprintf(stderr, "pipeline watch failed: %v\n", err)
		return 1
	}
	if m.PipelineSlug == "" {
		m.PipelineSlug = slug
	}
	values, err := pipeline.FileSignalSource{Path: *signalsPath}.Values()
	if err != nil {
		fmt.Fprintf(stderr, "pipeline watch failed: %v\n", err)
		return 1
	}
	now := time.Now()
	breaches, err := pipeline.EvaluateBreaches(m, values, now)
	if err != nil {
		fmt.Fprintf(stderr, "pipeline watch failed: %v\n", err)
		return 1
	}

	priorOpen, err := pipeline.ListOpenBreaches(deck, slug)
	if err != nil {
		fmt.Fprintf(stderr, "pipeline watch failed: %v\n", err)
		return 1
	}

	// Dedupe contract: Breach.Fingerprint() is deterministic over
	// signal+target+threshold+class, so an ongoing breach maps to the same key
	// across passes and is not re-opened.
	current := map[string]bool{}
	newCount, dedupCount, autoOpened, notified := 0, 0, 0, 0
	for _, b := range breaches {
		fp := b.Fingerprint()
		current[fp] = true
		rec, exists, err := pipeline.LoadBreachRecord(deck, slug, fp)
		if err != nil {
			fmt.Fprintf(stderr, "pipeline watch failed: %v\n", err)
			return 1
		}
		if exists && rec.Status != pipeline.BreachResolved {
			rec.LastSeen = now
			rec.Breach = b
			dedupCount++
			_ = pipeline.SaveBreachRecord(deck, rec)
			continue
		}
		rec = pipeline.BreachRecord{Fingerprint: fp, PipelineSlug: slug, Breach: b, FirstSeen: now, LastSeen: now, Status: pipeline.BreachNotified}
		newCount++
		if m.CanAutoOpen(b) {
			ideaSlug, err := openRemediationIdea(deck, slug, b, fp, now)
			if err != nil {
				fmt.Fprintf(stderr, "pipeline watch failed: %v\n", err)
				return 1
			}
			rec.Status = pipeline.BreachOpen
			rec.RemediationIdea = ideaSlug
			autoOpened++
			fmt.Fprintf(stdout, "breach %s (%s %s, class=%s): auto-opened low-risk remediation idea %s\n", fp, b.Signal, b.Threshold, b.Class, ideaSlug)
		} else {
			notified++
			fmt.Fprintf(stdout, "breach %s (%s %s, class=%s): notify+gate (no auto-remediation; not a low-risk class)\n", fp, b.Signal, b.Threshold, b.Class)
		}
		if err := pipeline.SaveBreachRecord(deck, rec); err != nil {
			fmt.Fprintf(stderr, "pipeline watch failed: %v\n", err)
			return 1
		}
	}

	resolved := 0
	for fp, rec := range priorOpen {
		if current[fp] {
			continue
		}
		rec.Status = pipeline.BreachResolved
		rec.LastSeen = now
		if err := pipeline.SaveBreachRecord(deck, rec); err != nil {
			fmt.Fprintf(stderr, "pipeline watch failed: %v\n", err)
			return 1
		}
		resolved++
		fmt.Fprintf(stdout, "breach %s resolved (signal recovered)\n", fp)
	}

	fmt.Fprintf(stdout, "watch %s: %d breach(es) — %d new (%d auto-remediation, %d notify), %d ongoing, %d resolved\n", slug, len(breaches), newCount, autoOpened, notified, dedupCount, resolved)
	return 0
}

// openRemediationIdea creates a gated remediation idea linked back to the
// pipeline and the triggering breach fingerprint.
func openRemediationIdea(deck, slug string, b pipeline.Breach, fp string, now time.Time) (string, error) {
	ideaSlug := slug + "__remediation-" + fp
	dir := filepath.Join(deck, "ideas", ideaSlug)
	if err := os.MkdirAll(filepath.Join(dir, "round-01"), 0o755); err != nil {
		return "", err
	}
	prompt := fmt.Sprintf(`---
idea: %s
author: watcher
created: %s
participants: []
status: round-01
pipeline_slug: %s
breach_fingerprint: %s
derived_from: [pipelines/%s/monitoring.yaml]
---

## Problem / idea

Auto-opened by the §12.11 monitoring watcher for pipeline %q after a predeclared low-risk SLO breach.

- signal: %s
- target: %s
- breach threshold: %s
- observed: %s
- breach class: %s (low-risk, auto-open permitted)

Remediation work itself remains gated: any production mutation must still go through an action block's non-bypassable production gate.

## Constraints
## Non-goals
`, ideaSlug, now.Format("2006-01-02"), slug, fp, slug, slug, b.Signal, b.Target, b.Threshold, b.Observed, b.Class)
	if err := os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte(prompt), 0o644); err != nil {
		return "", err
	}
	return ideaSlug, nil
}

// runPipelineExecute is the §12.10 execute sub-phase for action blocks. It
// validates preconditions (action plan finalized; production gate approved for
// production capabilities), plans the concrete provider call, and records the
// effect ledger entry. The Go CLI does NOT perform the side effect — it emits
// the ProviderCall for the driver/harness to execute via MCP, preserving the
// agents-write-markdown / driver-executes boundary (§12.4).
// ExecuteJSON is the versioned machine-readable output of `pipeline execute
// --json` (schema_version 1). An external driver/harness consumes it, performs
// the MCP call when the gate permits, then calls `record-effect`.
type ExecuteJSON struct {
	SchemaVersion  int                    `json:"schema_version"`
	Status         string                 `json:"status"` // dry_run|pending_gate|ready_for_harness
	ProviderCall   *pipeline.ProviderCall `json:"provider_call,omitempty"`
	EffectDigest   string                 `json:"effect_digest,omitempty"`
	IdempotencyKey string                 `json:"idempotency_key,omitempty"`
	Gate           executeGateJSON        `json:"gate"`
}

type executeGateJSON struct {
	Required bool   `json:"required"`
	State    string `json:"state"` // open|approved|not_required
}

func runPipelineExecute(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline execute", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	providerName := fs.String("provider", "vercel", "provider (vercel|noop)")
	dryRun := fs.Bool("dry-run", false, "plan only; never mutate")
	jsonOut := fs.Bool("json", false, "emit a machine-readable execution plan for an external harness")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 4 {
		fmt.Fprintln(stderr, "usage: parley pipeline execute [--dir DIR] [--provider P] [--dry-run] [--json] SLUG BLOCK CAPABILITY TARGET")
		return 2
	}
	deck := deckDirFor(*root)
	slug, blockID, capability, target := fs.Arg(0), fs.Arg(1), fs.Arg(2), fs.Arg(3)
	capName := pipeline.Capability(capability)

	m, err := pipeline.ParseFile(pipeline.ManifestPath(deck, slug))
	if err != nil {
		fmt.Fprintf(stderr, "pipeline execute failed: %v\n", err)
		return 1
	}
	block, ok := findBlock(m, blockID)
	if !ok {
		fmt.Fprintf(stderr, "pipeline execute failed: no block %q\n", blockID)
		return 1
	}
	if block.Kind != pipeline.KindAction {
		fmt.Fprintf(stderr, "pipeline execute failed: block %q is kind %q, only action blocks execute\n", blockID, block.Kind)
		return 1
	}
	// Precondition: the action plan must be finalized (§12.10).
	done, err := blockCompleteFunc(deck, slug)(block)
	if err != nil {
		fmt.Fprintf(stderr, "pipeline execute failed: %v\n", err)
		return 1
	}
	if !done {
		fmt.Fprintf(stderr, "pipeline execute failed: action plan for block %q is not finalized (status: final required)\n", blockID)
		return 1
	}
	prov, err := selectProvider(*providerName)
	if err != nil {
		fmt.Fprintf(stderr, "pipeline execute failed: %v\n", err)
		return 1
	}
	if !prov.Supports(capName) {
		fmt.Fprintf(stderr, "pipeline execute failed: provider %q does not support %q (capability halt, §12.9)\n", prov.Name(), capName)
		return 1
	}

	now := time.Now()
	risk := block.Risk
	if capName.IsProduction() {
		risk = pipeline.RiskProduction
	}
	reqHash := pipeline.HashRequest([]byte(target))
	key := pipeline.IdempotencyKey(slug, blockID, prov.Name(), string(capName), target, reqHash)
	digest := pipeline.KeyDigest(key)

	gateRequired := capName.IsProduction()
	gateApproved := false
	gateState := "not_required"
	if gateRequired && !*dryRun {
		edge := pipeline.EdgeID(blockID, "execute")
		g, exists, err := pipeline.LoadGate(deck, slug, edge)
		if err != nil {
			fmt.Fprintf(stderr, "pipeline execute failed: %v\n", err)
			return 1
		}
		if !exists {
			g = pipeline.NewGate(slug, blockID, "execute", pipeline.RiskProduction, "execute", now)
			if err := pipeline.SaveGate(deck, g); err != nil {
				fmt.Fprintf(stderr, "pipeline execute failed: %v\n", err)
				return 1
			}
			return emitExecutePendingGate(*jsonOut, stdout, *root, slug, edge, key, digest)
		}
		switch g.Status {
		case pipeline.GateOpen:
			return emitExecutePendingGate(*jsonOut, stdout, *root, slug, edge, key, digest)
		case pipeline.GateRejected:
			fmt.Fprintf(stderr, "pipeline execute failed: production gate %s was rejected\n", edge)
			return 1
		case pipeline.GateApproved:
			gateApproved = true
			gateState = "approved"
		}
	}

	call, err := prov.Plan(pipeline.ActionRequest{Capability: capName, Target: target, DryRun: *dryRun, GateApproved: gateApproved})
	if err != nil {
		fmt.Fprintf(stderr, "pipeline execute failed: %v\n", err)
		return 1
	}

	eff, exists, err := pipeline.LoadEffect(deck, slug, key)
	if err != nil {
		fmt.Fprintf(stderr, "pipeline execute failed: %v\n", err)
		return 1
	}
	if !exists {
		eff = pipeline.NewEffect(slug, blockID, prov.Name(), string(capName), target, reqHash, risk, now)
	}
	if eff.NeedsReconcile() {
		fmt.Fprintf(stderr, "pipeline execute blocked: effect %s is %s (external_ref %q); reconcile with record-effect before re-attempting (§12.7)\n", digest, eff.Status, eff.ExternalRef)
		return 1
	}
	status := "ready_for_harness"
	if *dryRun {
		eff.DryRunResult = "planned (dry-run)"
		eff.Advance(pipeline.EffectDryRunOK, "", "dry-run planned", now)
		status = "dry_run"
	} else {
		eff.Advance(pipeline.EffectExecuting, "", "planned for harness execution", now)
	}
	if err := pipeline.SaveEffect(deck, eff); err != nil {
		fmt.Fprintf(stderr, "pipeline execute failed: %v\n", err)
		return 1
	}

	if *jsonOut {
		out := ExecuteJSON{
			SchemaVersion:  1,
			Status:         status,
			ProviderCall:   &call,
			EffectDigest:   digest,
			IdempotencyKey: key,
			Gate:           executeGateJSON{Required: gateRequired, State: gateState},
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(stdout, string(b))
		return 0
	}

	callJSON, _ := json.MarshalIndent(call, "", "  ")
	fmt.Fprintf(stdout, "effect %s -> %s (risk=%s)\n", digest, eff.Status, risk)
	fmt.Fprintf(stdout, "ProviderCall (the driver/harness performs this via MCP; the CLI does not):\n%s\n", callJSON)
	if !*dryRun {
		fmt.Fprintf(stdout, "after the harness runs it, record the result:\n  parley pipeline record-effect --dir %s --status succeeded --external-ref <ref> %s %s\n", *root, slug, digest)
	}
	return 0
}

// emitExecutePendingGate reports that a production gate must be approved before
// execution, in JSON or human form.
func emitExecutePendingGate(jsonOut bool, stdout io.Writer, root, slug, edge, key, digest string) int {
	if jsonOut {
		out := ExecuteJSON{
			SchemaVersion:  1,
			Status:         "pending_gate",
			EffectDigest:   digest,
			IdempotencyKey: key,
			Gate:           executeGateJSON{Required: true, State: "open"},
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(stdout, string(b))
		return 0
	}
	fmt.Fprintf(stdout, "production execution requires approval (non-bypassable, §12.8).\n  approve with: parley pipeline gate approve --dir %s %s %s\n", root, slug, edge)
	return 0
}

// runPipelineRecordEffect persists the outcome of a harness-executed effect
// (or reconciles an ambiguous one) by digest.
func runPipelineRecordEffect(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pipeline record-effect", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	status := fs.String("status", "", "succeeded|failed|reconciled|abandoned")
	externalRef := fs.String("external-ref", "", "provider external reference (e.g. deployment id)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 2 {
		fmt.Fprintln(stderr, "usage: parley pipeline record-effect [--dir DIR] --status STATUS [--external-ref REF] SLUG DIGEST")
		return 2
	}
	allowed := map[string]pipeline.EffectStatus{
		"succeeded":  pipeline.EffectSucceeded,
		"failed":     pipeline.EffectFailed,
		"reconciled": pipeline.EffectReconciled,
		"abandoned":  pipeline.EffectAbandoned,
	}
	st, ok := allowed[*status]
	if !ok {
		fmt.Fprintln(stderr, "record-effect --status must be succeeded|failed|reconciled|abandoned")
		return 2
	}
	deck := deckDirFor(*root)
	slug, digest := fs.Arg(0), fs.Arg(1)
	eff, ok, err := pipeline.LoadEffectByDigest(deck, slug, digest)
	if err != nil {
		fmt.Fprintf(stderr, "record-effect failed: %v\n", err)
		return 1
	}
	if !ok {
		fmt.Fprintf(stderr, "record-effect failed: no effect %q for pipeline %q\n", digest, slug)
		return 1
	}
	eff.Advance(st, *externalRef, "recorded by harness/driver", time.Now())
	if err := pipeline.SaveEffect(deck, eff); err != nil {
		fmt.Fprintf(stderr, "record-effect failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "effect %s -> %s (external_ref %q)\n", digest, eff.Status, eff.ExternalRef)
	return 0
}

// blockCompleteFunc reports a block finished when its engine workspace holds a
// finalized artifact (FINAL.md or the block's output artifact with frontmatter
// status: final).
func blockCompleteFunc(deck, slug string) pipeline.BlockComplete {
	return func(b pipeline.Block) (bool, error) {
		ws := pipeline.BlockWorkspace(deck, slug, b.ID)
		candidates := []string{"FINAL.md"}
		if b.OutputArtifact != "" {
			candidates = append([]string{b.OutputArtifact}, candidates...)
		}
		for _, name := range candidates {
			data, err := os.ReadFile(filepath.Join(ws, name))
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return false, err
			}
			if isFinalized(string(data)) {
				return true, nil
			}
		}
		return false, nil
	}
}

func isFinalized(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == "status: final" {
			return true
		}
	}
	return false
}

func valueOrDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
