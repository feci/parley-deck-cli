package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"parley-deck-cli/internal/budget"
)

func runBudgetPolicy(ctx context.Context, surface string, args []string, stdout, stderr io.Writer, attended bool) int {
	if len(args) == 0 || args[0] != "inspect" && args[0] != "extend" {
		fmt.Fprintf(stderr, "usage: parley budget %s inspect|extend --dir DIR [--idea ID] [options]\n", surface)
		return 2
	}
	f := flag.NewFlagSet("budget "+surface+" "+args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	root := f.String("dir", ".", "workspace; Git worktrees share this policy scope")
	idea := f.String("idea", "", "idea identity; required for steps, empty for auxiliary launches")
	launches := f.Int("max-launches", -1, "absolute launch ceiling; unchanged zero retains an unlimited axis")
	steps := f.Int("max-steps", -1, "absolute driver-step ceiling; unchanged zero retains an unlimited axis")
	cost := f.Int64("max-cost-micros", -1, "absolute monetary ceiling; reservation amount is preserved")
	wall := f.Duration("wall-clock", -1, "absolute lifetime ceiling from original activation, not extra time")
	expected := f.String("expected-policy-sha256", "", "canonical policy hash returned by inspect")
	decision := f.String("decision-id", "", "unique operator decision ID; exact replay is idempotent")
	reason := f.String("reason", "", "reason for this finite extension")
	yes := f.Bool("yes", false, "record this exact attended decision")
	if err := f.Parse(args[1:]); err != nil {
		return 2
	}
	if f.NArg() != 0 || surface == "step" && strings.TrimSpace(*idea) == "" {
		fmt.Fprintln(stderr, "budget policy: no positional arguments; step policy requires --idea")
		return 2
	}
	kind := budget.Launch
	if surface == "step" {
		kind = budget.DriverStep
	}
	visited := map[string]bool{}
	f.Visit(func(f *flag.Flag) { visited[f.Name] = true })
	var status budget.PolicyStatus
	var err error
	if args[0] == "inspect" {
		for name := range visited {
			if name != "dir" && name != "idea" {
				fmt.Fprintln(stderr, "policy inspect accepts only --dir and --idea")
				return 2
			}
		}
		status, err = budget.InspectRuntimeBudget(ctx, *root, *idea, kind)
	} else {
		if !attended {
			fmt.Fprintln(stderr, "policy extension requires an attended terminal and explicit operator decision; participant frontmatter is not a grant")
			return 2
		}
		valid := *yes && *expected != "" && *decision != "" && strings.TrimSpace(*reason) != "" && visited["wall-clock"] && *wall >= 0
		ceilings := budget.PolicyCeilings{WallClockNS: int64(*wall)}
		if kind == budget.Launch {
			valid = valid && visited["max-launches"] && *launches >= 0 && visited["max-cost-micros"] && *cost >= 0 && !visited["max-steps"] && *wall%time.Millisecond == 0
			ceilings.Actions, ceilings.CostMicros = *launches, *cost
		} else {
			valid = valid && visited["max-steps"] && *steps >= 0 && !visited["max-launches"] && !visited["max-cost-micros"]
			ceilings.Actions = *steps
		}
		if !valid {
			fmt.Fprintln(stderr, "policy extend requires every absolute ceiling for its kind, --wall-clock, --expected-policy-sha256, --decision-id, --reason and --yes; at least one finite axis must increase")
			return 2
		}
		status, err = budget.ExtendRuntimeBudget(ctx, *root, *idea, kind, budget.PolicyExtensionRequest{DecisionID: *decision, ExpectedPolicySHA256: *expected, Reason: *reason, Ceilings: ceilings})
	}
	if err != nil {
		fmt.Fprintf(stderr, "budget %s %s: %v\n", surface, args[0], err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(status); err != nil {
		fmt.Fprintln(stderr, "policy result output failed; inspect or replay the exact decision before any work")
		return 1
	}
	return 0
}
