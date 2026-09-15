package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"time"

	"parley-deck-cli/internal/budget"
)

func runBudgetMigrationRecovery(ctx context.Context, args []string, stdout, stderr io.Writer, attended bool) int {
	if len(args) == 0 || args[0] != "inspect" && args[0] != "apply" {
		fmt.Fprintln(stderr, "usage: parley budget migrate recover inspect|apply --kind launch|step|fixup|cross-review --dir DIR [--idea ID] [options]")
		return 2
	}
	f := flag.NewFlagSet("budget migrate recover "+args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	root := f.String("dir", ".", "workspace containing the retained import")
	idea := f.String("idea", "", "idea identity; empty means auxiliary launches")
	rawKind := f.String("kind", "launch", "launch, step, fixup or cross-review")
	importHash := f.String("expected-import-sha256", "", "exact import state hash from recovery inspect")
	history := f.String("expected-history-sha256", "", "exact current history hash from recovery inspect")
	decision := f.String("decision-id", "", "unique operator recovery decision")
	reason := f.String("reason", "", "operator explanation and accounting basis")
	started := f.String("started-at", "", "original or earlier accounting epoch in RFC3339 format")
	additional := f.Int("additional-launches", -1, "total pre-telemetry attempts; cannot decrease")
	total := f.Int("total-actions", -1, "total historical protocol actions; cannot decrease")
	writers := f.Bool("writers-stopped", false, "operator confirms all affected writers are stopped")
	yes := f.Bool("yes", false, "apply this exact attended recovery decision without changing ceilings")
	if err := f.Parse(args[1:]); err != nil {
		return 2
	}
	kind, ok := map[string]budget.Kind{"launch": budget.Launch, "step": budget.DriverStep, "fixup": budget.Fixup, "cross-review": budget.CrossReview}[*rawKind]
	if !ok || f.NArg() != 0 || kind != budget.Launch && *idea == "" {
		fmt.Fprintln(stderr, "recovery requires a valid --kind and protocol --idea, with no positional arguments")
		return 2
	}
	visited := map[string]bool{}
	f.Visit(func(f *flag.Flag) { visited[f.Name] = true })
	var result any
	var err error
	if args[0] == "inspect" {
		for name := range visited {
			if name != "dir" && name != "idea" && name != "kind" {
				fmt.Fprintln(stderr, "recovery inspect accepts only --dir, --idea and --kind")
				return 2
			}
		}
		result, err = budget.InspectMigrationRecovery(ctx, *root, *idea, kind)
	} else {
		if !attended {
			fmt.Fprintln(stderr, "migration recovery requires an attended terminal and the operator's concrete decision; no participant field supplies attendance")
			return 2
		}
		count, countFlag, wrongFlag := *additional, "additional-launches", "total-actions"
		if kind != budget.Launch {
			count, countFlag, wrongFlag = *total, "total-actions", "additional-launches"
		}
		at, parseErr := time.Parse(time.RFC3339Nano, *started)
		if !*yes || !*writers || parseErr != nil || !visited[countFlag] || visited[wrongFlag] || count < 0 {
			fmt.Fprintf(stderr, "recovery apply requires --%s, --started-at, --expected-import-sha256, --expected-history-sha256, --decision-id, --reason, --writers-stopped and --yes; policy ceilings are inherited unchanged\n", countFlag)
			return 2
		}
		result, err = budget.RecoverMigration(ctx, *root, *idea, kind, budget.MigrationRecoveryRequest{ExpectedImportSHA256: *importHash, ExpectedHistorySHA256: *history, DecisionID: *decision, Reason: *reason, StartedAt: at, AccountedActions: count, WritersStopped: *writers})
	}
	if err != nil {
		fmt.Fprintf(stderr, "budget migrate recover %s: %v\n", args[0], err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(result); err != nil {
		fmt.Fprintln(stderr, "recovery output failed; inspect existing state or replay the exact decision before work")
		return 1
	}
	return 0
}
