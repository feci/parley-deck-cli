package app

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/loop"
	"parley-deck-cli/internal/protocol"
)

// runLoop dispatches `parley loop <sub>`. LE-9: the only subcommand is `tick`, a
// one-shot, scheduler-friendly discovery pass. It is governed by COOPERATION.md §14
// (the human brake): it drafts `status: candidate` ideas only — never runs, pushes,
// merges, finalizes, or staffs a quorum.
func runLoop(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: parley loop tick [--dir DIR] [--signals PATH] [--enable] [--json]")
		return 2
	}
	switch args[0] {
	case "tick":
		return runLoopTick(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown loop subcommand: %s\n", args[0])
		fmt.Fprintln(stderr, "usage: parley loop tick [--dir DIR] [--signals PATH] [--enable] [--json]")
		return 2
	}
}

func runLoopTick(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("loop tick", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace root")
	signalsPath := fs.String("signals", "", "signals JSON path (default: <deck>/loop/signals.json)")
	enable := fs.Bool("enable", false, "force-enable this one-off tick (still candidate-only)")
	jsonOut := fs.Bool("json", false, "JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "loop tick: %v\n", err)
		return 1
	}
	deck := filepath.Join(root, protocol.DeckDir)

	cfg, err := loop.ReadConfig(deck)
	if err != nil {
		fmt.Fprintf(stderr, "loop tick: %v\n", err)
		return 1
	}
	if *enable {
		// A human explicitly running `--enable` is the human gate for THIS tick.
		// It still only drafts candidates (§14); it never promotes/runs/pushes.
		cfg.Enabled = true
	}

	// AF3: short-circuit the disabled case BEFORE reading signals, so a disabled
	// tick is fully inert and cron-safe even when the signals file is malformed.
	if !cfg.Enabled {
		if *jsonOut {
			return printJSON(stdout, loop.TickResult{Enabled: false}, stderr)
		}
		fmt.Fprintln(stdout, "loop tick: disabled (set loop/config.json {\"enabled\": true} or pass --enable). Wrote nothing.")
		return 0
	}

	sigPath := *signalsPath
	if sigPath == "" {
		sigPath = loop.SignalsPath(deck)
	}
	signals, err := loop.ReadSignals(sigPath)
	if err != nil {
		fmt.Fprintf(stderr, "loop tick: %v\n", err)
		return 1
	}

	res, err := loop.Tick(deck, cfg, signals, time.Now().UTC())
	if err != nil {
		fmt.Fprintf(stderr, "loop tick: %v\n", err)
		return 1
	}

	if *jsonOut {
		return printJSON(stdout, res, stderr)
	}
	fmt.Fprintf(stdout, "loop tick: %d candidate(s) drafted, %d skipped (already present), %d rejected (unknown source).\n",
		len(res.Created), len(res.Skipped), len(res.Rejected))
	for _, s := range res.Created {
		fmt.Fprintf(stdout, "  + %s (status: candidate — promote manually to deliberate)\n", s)
	}
	for _, s := range res.Skipped {
		fmt.Fprintf(stdout, "  = %s (exists)\n", s)
	}
	for _, s := range res.Rejected {
		fmt.Fprintf(stdout, "  ! %s (rejected — source not in commit|ci|issue|manual)\n", s)
	}
	return 0
}
