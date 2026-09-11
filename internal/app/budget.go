package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"parley-deck-cli/internal/budget"
)

func runBudget(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	supported, attended := budgetAttendance()
	return runBudgetPlatformControl(ctx, args, stdout, stderr, supported, attended)
}

func runBudgetPlatformControl(ctx context.Context, args []string, stdout, stderr io.Writer, supported, attended bool) int {
	if len(args) > 0 && (args[0] == "reconcile" || args[0] == "configure") && !supported {
		if args[0] == "configure" {
			fmt.Fprintln(stderr, "budget configure: attended policy activation is unavailable on this platform; no policy was activated. No unattended override exists.")
			return 2
		}
		fmt.Fprintln(stderr, "budget reconcile: attended recovery is unavailable on this platform. A ledger originating here has no supported attended recovery or migration yet; preserve all charges and stop writers. No unattended override exists.")
		return 2
	}
	return runBudgetControl(ctx, args, stdout, stderr, attended)
}

// The bool is an internal test seam. The CLI always computes it from the
// platform probe; no flag or participant-authored field supplies attendance.
// As with protocol publication, terminal presence is not human authentication.
func runBudgetControl(ctx context.Context, args []string, stdout, stderr io.Writer, attended bool) int {
	if len(args) > 0 && args[0] == "configure" {
		return runBudgetConfigure(ctx, args[1:], stdout, stderr, attended)
	}
	if len(args) == 0 || (args[0] != "inspect" && args[0] != "reconcile") {
		fmt.Fprintln(stderr, "usage: parley budget inspect|reconcile --ledger DIR --scope ID [options]; parley budget configure --dir DIR [--idea ID] --max-launches N --max-cost-micros N --wall-clock D --yes")
		return 2
	}
	fs := flag.NewFlagSet("budget "+args[0], flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("ledger", "", "exact existing ledger directory")
	scope := fs.String("scope", "", "exact existing budget scope")
	action := fs.String("action-id", "", "original charged action identity")
	decision := fs.String("decision-id", "", "unique operator decision identity; reuse only for exact replay")
	ceiling := fs.Int64("ceiling-micros", -1, "explicit conservative monetary ceiling in millionths of USD")
	reason := fs.String("reason", "", "operator rationale for the conservative ceiling")
	yes := fs.Bool("yes", false, "apply this exact attended cost decision")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	if fs.NArg() != 0 || strings.TrimSpace(*dir) == "" || strings.TrimSpace(*scope) == "" {
		fmt.Fprintln(stderr, "budget: --ledger and --scope are required; no positional arguments")
		return 2
	}
	s := budget.Store{Dir: *dir, Scope: *scope}
	if args[0] == "inspect" {
		if *action != "" || *decision != "" || *ceiling != -1 || *reason != "" || *yes {
			fmt.Fprintln(stderr, "budget inspect accepts only --ledger and --scope")
			return 2
		}
		state, err := s.Inspect(ctx)
		if err != nil {
			fmt.Fprintf(stderr, "budget inspect: %v\n", err)
			return 1
		}
		if err := json.NewEncoder(stdout).Encode(state); err != nil {
			fmt.Fprintln(stderr, "budget inspect: cannot write result")
			return 1
		}
		return 0
	}
	if !attended {
		fmt.Fprintln(stderr, "budget reconcile: refusing unattended cost adjustment; use this explicit operator control from a terminal. Participant frontmatter does not grant a budget extension.")
		return 2
	}
	if !*yes || *action == "" || *decision == "" || *ceiling < 0 || strings.TrimSpace(*reason) == "" {
		fmt.Fprintln(stderr, "budget reconcile requires --action-id, --decision-id, --ceiling-micros, --reason and --yes")
		return 2
	}
	if *ceiling == 0 {
		fmt.Fprintln(stderr, "budget reconcile: explicit zero ceiling counts this unknown observation as zero exposure; the observed cost remains unknown and this operator decision is retained.")
	}
	if _, err := s.ReconcileUnknown(ctx, *action, *decision, *ceiling, *reason); err != nil {
		fmt.Fprintf(stderr, "budget reconcile: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, "Recorded conservative operator ceiling. Observed cost stays unknown; the action stays spent.")
	return 0
}
