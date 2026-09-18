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

// configure activates a first frozen policy. It does not clear charges, extend
// a policy, or invent a legacy count. Attendance is runtime attribution, not
// proof that a human controls an otherwise writable same-UID environment.
func runBudgetConfigure(ctx context.Context, args []string, stdout, stderr io.Writer, attended bool) int {
	flags := flag.NewFlagSet("budget configure", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("dir", ".", "workspace; Git worktrees share this scope")
	idea := flags.String("idea", "", "idea identity; empty means auxiliary launches")
	launches := flags.Int("max-launches", -1, "explicit lifetime launch cap; 0 = unlimited")
	cost := flags.Int64("max-cost-micros", -1, "explicit USD microdollar cap; 0 = unlimited")
	reserve := flags.Int64("reserve-micros", -1, "conservative per-launch maximum; absent means unknown")
	wall := flags.Duration("wall-clock", -1, "explicit wall time from policy activation; 0 = unlimited")
	yes := flags.Bool("yes", false, "activate this exact attended policy")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if !attended {
		fmt.Fprintln(stderr, "budget configure requires an attended terminal; participant frontmatter cannot activate or extend a budget")
		return 2
	}
	if flags.NArg() != 0 || !*yes || *launches < 0 || *cost < 0 || *wall < 0 || *reserve < -1 || (*wall)%time.Millisecond != 0 {
		fmt.Fprintln(stderr, "usage: parley budget configure --dir DIR [--idea ID] --max-launches N --max-cost-micros N --wall-clock D [--reserve-micros N] --yes; provide every ceiling explicitly")
		return 2
	}
	policy := budget.LaunchPolicy{MaxLaunches: *launches, MaxCostMicros: *cost, WallClockMS: wall.Milliseconds()}
	if *reserve >= 0 {
		policy.ReserveMicros = reserve
	}
	binding, err := budget.ConfigureLaunchBudget(ctx, *root, *idea, policy)
	if err != nil {
		fmt.Fprintf(stderr, "budget configure: %v\n", err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(binding); err != nil {
		fmt.Fprintln(stderr, "budget configured but result output failed")
		return 1
	}
	return 0
}
