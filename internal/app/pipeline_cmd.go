package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/pipeline"
	"parley-deck-cli/internal/protocol"
)

func printPipelineUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: parley pipeline validate MANIFEST | start [--dir DIR] MANIFEST | status [--dir DIR] SLUG | continue [--dir DIR] SLUG | gate approve|reject [--dir DIR] [--by WHO] SLUG EDGE")
}

func runPipeline(args []string, stdout, stderr io.Writer) int {
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
	case "continue":
		return runPipelineContinue(args[1:], stdout, stderr)
	case "gate":
		return runPipelineGate(args[1:], stdout, stderr)
	default:
		printPipelineUsage(stderr)
		return 2
	}
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
