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

func runBudgetMigrate(ctx context.Context, args []string, stdout, stderr io.Writer, attended bool) int {
	if len(args) == 0 || args[0] != "inspect" && args[0] != "apply" {
		fmt.Fprintln(stderr, "usage: parley budget migrate inspect|apply --kind launch --dir DIR [--idea ID] [options]")
		return 2
	}
	f := flag.NewFlagSet("budget migrate "+args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	root := f.String("dir", ".", "workspace; all linked worktrees are inventoried")
	idea := f.String("idea", "", "idea identity; empty means auxiliary launches")
	kind := f.String("kind", "launch", "accounting kind; launch is currently supported")
	history := f.String("expected-history-sha256", "", "exact inventory hash from migration inspect")
	decision := f.String("decision-id", "", "unique operator migration decision")
	reason := f.String("reason", "", "operator explanation and legacy accounting basis")
	started := f.String("started-at", "", "original accounting epoch in RFC3339 format; no later than observed history")
	additional := f.Int("additional-launches", -1, "explicit attempts preceding invocation telemetry; 0 is an explicit assertion")
	writers := f.Bool("writers-stopped", false, "operator confirms all writers are stopped for migration")
	max := f.Int("max-launches", -1, "absolute future lifetime launch ceiling including imported attempts; 0 = unlimited")
	cost := f.Int64("max-cost-micros", -1, "absolute monetary ceiling in millionths of USD; 0 = unlimited")
	reserve := f.Int64("reserve-micros", -1, "conservative reservation per future launch; absent = unknown")
	wall := f.Duration("wall-clock", -1, "absolute lifetime from started-at, including prior work; 0 = unlimited")
	yes := f.Bool("yes", false, "apply this exact attended migration decision")
	if err := f.Parse(args[1:]); err != nil {
		return 2
	}
	if f.NArg() != 0 || *kind != "launch" {
		fmt.Fprintln(stderr, "migration accepts no positional arguments; only launch accounting is currently supported")
		return 2
	}
	visited := map[string]bool{}
	f.Visit(func(f *flag.Flag) { visited[f.Name] = true })
	var result any
	var err error
	if args[0] == "inspect" {
		for name := range visited {
			if name != "kind" && name != "dir" && name != "idea" {
				fmt.Fprintln(stderr, "migration inspect accepts only --kind, --dir and --idea")
				return 2
			}
		}
		result, err = budget.InspectLaunchMigration(ctx, *root, *idea)
	} else {
		if !attended {
			fmt.Fprintln(stderr, "launch migration requires an attended terminal and the operator's concrete decision; no participant field supplies attendance")
			return 2
		}
		at, parseErr := time.Parse(time.RFC3339Nano, *started)
		if !*yes || !*writers || !visited["additional-launches"] || *additional < 0 || !visited["max-launches"] || *max < 0 || !visited["max-cost-micros"] || *cost < 0 || !visited["wall-clock"] || *wall < 0 || *wall%time.Millisecond != 0 || visited["reserve-micros"] && *reserve < 0 || parseErr != nil {
			fmt.Fprintln(stderr, "migration apply requires explicit --additional-launches, --started-at, --max-launches, --max-cost-micros, --wall-clock, --expected-history-sha256, --decision-id, --reason, --writers-stopped and --yes")
			return 2
		}
		p := budget.LaunchPolicy{MaxLaunches: *max, MaxCostMicros: *cost, WallClockMS: wall.Milliseconds()}
		if *reserve >= 0 {
			p.ReserveMicros = reserve
		}
		result, err = budget.MigrateLaunchBudget(ctx, *root, *idea, budget.LaunchMigrationRequest{ExpectedHistorySHA256: *history, DecisionID: *decision, Reason: *reason, StartedAt: at, AdditionalLaunches: *additional, WritersStopped: *writers, Policy: p})
	}
	if err != nil {
		fmt.Fprintf(stderr, "budget migrate %s: %v\n", args[0], err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(result); err != nil {
		fmt.Fprintln(stderr, "migration output failed; inspect existing state or replay the exact decision before work")
		return 1
	}
	return 0
}
