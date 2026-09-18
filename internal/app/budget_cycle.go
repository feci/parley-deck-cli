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

func runBudgetCycle(ctx context.Context, args []string, stdout, stderr io.Writer, attended bool) int {
	if len(args) == 0 || args[0] != "inspect" && args[0] != "extend" {
		fmt.Fprintln(stderr, "usage: parley budget cycle inspect|extend --dir DIR --idea ID --kind fixup|cross-review [options]")
		return 2
	}
	f := flag.NewFlagSet("budget cycle "+args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	root := f.String("dir", ".", "workspace; Git worktrees share this idea scope")
	idea := f.String("idea", "", "existing idea identity")
	kind := f.String("kind", "", "fixup or cross-review")
	maximum := f.Int("max-cycles", -1, "new finite absolute maximum, not additional cycles")
	expected := f.String("expected-policy-sha256", "", "exact canonical policy digest from cycle inspect")
	decision := f.String("decision-id", "", "unique operator decision identity; exact replay is idempotent")
	reason := f.String("reason", "", "operator reason for this finite grant")
	yes := f.Bool("yes", false, "record this exact attended finite extension")
	if err := f.Parse(args[1:]); err != nil {
		return 2
	}
	if f.NArg() != 0 || strings.TrimSpace(*idea) == "" || (*kind != string(budget.Fixup) && *kind != string(budget.CrossReview)) {
		fmt.Fprintln(stderr, "budget cycle requires --idea and --kind fixup|cross-review, with no positional arguments")
		return 2
	}
	var status budget.CycleStatus
	var err error
	if args[0] == "inspect" {
		if *maximum != -1 || *expected != "" || *decision != "" || *reason != "" || *yes {
			fmt.Fprintln(stderr, "budget cycle inspect accepts only --dir, --idea and --kind")
			return 2
		}
		status, err = budget.InspectCycleBudget(ctx, *root, *idea, budget.Kind(*kind))
	} else {
		if !attended {
			fmt.Fprintln(stderr, "budget cycle extend requires an attended terminal and an explicit operator decision; participant frontmatter is not a grant")
			return 2
		}
		if !*yes || *maximum <= 0 || *expected == "" || *decision == "" || strings.TrimSpace(*reason) == "" {
			fmt.Fprintln(stderr, "budget cycle extend requires --max-cycles N --expected-policy-sha256 HASH --decision-id ID --reason TEXT --yes; maximum is finite and absolute")
			return 2
		}
		status, err = budget.ExtendCycleBudget(ctx, *root, *idea, budget.Kind(*kind), budget.CycleExtensionRequest{DecisionID: *decision, ExpectedPolicySHA256: *expected, Reason: *reason, Maximum: *maximum})
	}
	if err != nil {
		fmt.Fprintf(stderr, "budget cycle %s: %v\n", args[0], err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(status); err != nil {
		fmt.Fprintln(stderr, "budget cycle result output failed; inspect policy before retrying any work")
		return 1
	}
	return 0
}
