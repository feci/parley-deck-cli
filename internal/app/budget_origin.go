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

func runBudgetOrigin(ctx context.Context, args []string, stdout, stderr io.Writer, attended bool) int {
	if len(args) == 0 || args[0] != "inspect" && args[0] != "apply" {
		fmt.Fprintln(stderr, "usage: parley budget origin inspect|apply --resource DIR --target-lock-dir DIR [--expected-sha256 HASH --decision-id ID --reason TEXT --yes]")
		return 2
	}
	f := flag.NewFlagSet("budget origin "+args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	resource := f.String("resource", "", "exact existing guard or ledger metadata directory")
	target := f.String("target-lock-dir", "", "existing host-local destination directory for the permanent lock")
	hash := f.String("expected-sha256", "", "exact origin inspection digest")
	decision := f.String("decision-id", "", "unique operator migration decision; reuse only for exact recovery/replay")
	reason := f.String("reason", "", "operator rationale for relocating the lock")
	yes := f.Bool("yes", false, "apply the exact attended origin migration")
	if e := f.Parse(args[1:]); e != nil {
		return 2
	}
	if f.NArg() != 0 || *resource == "" || *target == "" {
		fmt.Fprintln(stderr, "origin control requires --resource and --target-lock-dir, with no positional arguments")
		return 2
	}
	wait, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	var result budget.LockOriginPreview
	var err error
	if args[0] == "inspect" {
		invalid := false
		f.Visit(func(flag *flag.Flag) {
			if flag.Name != "resource" && flag.Name != "target-lock-dir" {
				invalid = true
			}
		})
		if invalid {
			fmt.Fprintln(stderr, "origin inspect accepts only --resource and --target-lock-dir")
			return 2
		}
		result, err = budget.PreviewLockOriginMigration(wait, *resource, *target)
	} else {
		if !attended || !*yes {
			fmt.Fprintln(stderr, "origin apply requires an attended terminal and --yes; participant metadata or a flag cannot supply attendance")
			return 2
		}
		result, err = budget.ApplyLockOriginMigration(wait, *resource, *target, budget.LockOriginRequest{ExpectedSHA256: *hash, DecisionID: *decision, Reason: *reason})
	}
	if err != nil {
		fmt.Fprintf(stderr, "budget origin %s: %v\n", args[0], err)
		return 1
	}
	if err = json.NewEncoder(stdout).Encode(result); err != nil {
		fmt.Fprintln(stderr, "origin output failed; preserve records and replay the exact apply decision before work")
		return 1
	}
	return 0
}
