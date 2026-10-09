package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"parley-deck-cli/internal/budget"
)

func runBudgetLegacy(ctx context.Context, args []string, stdout, stderr io.Writer, attended bool) int {
	if len(args) == 0 || args[0] != "inspect" && args[0] != "apply" {
		fmt.Fprintln(stderr, "usage: parley budget legacy inspect|apply --dir DIR --run parley-deck/runs/NAME [options]")
		return 2
	}
	f := flag.NewFlagSet("budget legacy "+args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	root := f.String("dir", ".", "workspace; all visible linked worktrees are checked")
	path := f.String("run", "", "exact canonical relative legacy run directory")
	preview := f.String("expected-preview-sha256", "", "exact read-only preview digest")
	decision := f.String("decision-id", "", "unique operator decision ID; replay only identically")
	reason := f.String("reason", "", "operator explanation; no historical zero is asserted")
	writers := f.Bool("writers-stopped", false, "operator confirms all writers are stopped")
	unknown := f.Bool("acknowledge-unknown-history", false, "acknowledge that declared history remains unknown, never zero")
	yes := f.Bool("yes", false, "apply this exact attended declaration")
	if err := f.Parse(args[1:]); err != nil {
		return 2
	}
	if f.NArg() != 0 || *path == "" {
		fmt.Fprintln(stderr, "legacy declaration requires --run and no positional arguments")
		return 2
	}
	var result any
	var err error
	if args[0] == "inspect" {
		invalid := false
		f.Visit(func(f *flag.Flag) {
			if f.Name != "dir" && f.Name != "run" {
				invalid = true
			}
		})
		if invalid {
			fmt.Fprintln(stderr, "legacy inspect accepts only --dir and --run")
			return 2
		}
		result, err = budget.InspectLegacyDeclaration(ctx, *root, *path)
	} else {
		if !attended {
			fmt.Fprintln(stderr, "legacy apply requires an attended terminal and the operator's concrete declaration; no unattended override exists")
			return 2
		}
		if !*yes || !*writers || !*unknown {
			fmt.Fprintln(stderr, "legacy apply requires --yes, --writers-stopped and --acknowledge-unknown-history")
			return 2
		}
		result, err = budget.ApplyLegacyDeclaration(ctx, *root, budget.LegacyDeclarationRequest{Path: *path, ExpectedPreviewSHA256: *preview, DecisionID: *decision, Reason: *reason, WritersStopped: *writers, AcknowledgeUnknown: *unknown})
	}
	if err != nil {
		fmt.Fprintf(stderr, "budget legacy %s: %v\n", args[0], err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(result); err != nil {
		fmt.Fprintln(stderr, "legacy output failed; preserve state and replay the exact decision before work")
		return 1
	}
	return 0
}
