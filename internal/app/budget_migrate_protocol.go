package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"parley-deck-cli/internal/budget"
)

type protocolMigrationOptions struct {
	Root, Idea, IdeaPath, Kind, History, Decision, Reason, Started string
	Total, Steps, Cycles                                           int
	Wall                                                           time.Duration
	Writers, Yes                                                   bool
}

func runBudgetMigrateProtocol(ctx context.Context, verb string, o protocolMigrationOptions, visited map[string]bool, stdout, stderr io.Writer, attended bool) int {
	kind, ok := map[string]budget.Kind{"step": budget.DriverStep, "fixup": budget.Fixup, "cross-review": budget.CrossReview}[o.Kind]
	if !ok || o.Idea == "" {
		fmt.Fprintln(stderr, "protocol migration requires --idea and --kind step, fixup or cross-review")
		return 2
	}
	allowed := map[string]bool{"dir": true, "idea": true, "idea-path": true, "kind": true}
	if verb == "apply" {
		for _, name := range []string{"expected-history-sha256", "decision-id", "reason", "started-at", "total-actions", "writers-stopped", "yes"} {
			allowed[name] = true
		}
		if kind == budget.DriverStep {
			allowed["max-steps"], allowed["wall-clock"] = true, true
		} else {
			allowed["max-cycles"] = true
		}
	}
	for name := range visited {
		if !allowed[name] {
			fmt.Fprintf(stderr, "%s migration %s does not accept --%s\n", o.Kind, verb, name)
			return 2
		}
	}
	var result any
	var err error
	if verb == "inspect" {
		result, err = budget.InspectProtocolMigration(ctx, o.Root, o.Idea, kind, o.IdeaPath)
	} else {
		if !attended {
			fmt.Fprintln(stderr, "protocol migration requires an attended terminal and the operator's explicit accounting decision")
			return 2
		}
		at, parseErr := time.Parse(time.RFC3339Nano, o.Started)
		maximum := o.Cycles
		wall := time.Duration(0)
		complete := visited["max-cycles"] && maximum >= 0
		if kind == budget.DriverStep {
			maximum, wall = o.Steps, o.Wall
			complete = visited["max-steps"] && maximum >= 0 && visited["wall-clock"] && wall >= 0
		}
		if parseErr != nil || !visited["total-actions"] || o.Total < 0 || !complete || !o.Writers || !o.Yes {
			fmt.Fprintln(stderr, "protocol migration apply requires --total-actions, --started-at, --expected-history-sha256, --decision-id, --reason, --writers-stopped, --yes and explicit ceilings (--max-steps/--wall-clock or --max-cycles)")
			return 2
		}
		result, err = budget.MigrateProtocolBudget(ctx, o.Root, o.Idea, kind, budget.ProtocolMigrationRequest{ExpectedHistorySHA256: o.History, DecisionID: o.Decision, Reason: o.Reason, IdeaPath: o.IdeaPath, StartedAt: at, TotalActions: o.Total, Maximum: maximum, WallClockNS: int64(wall), WritersStopped: o.Writers})
	}
	if err != nil {
		fmt.Fprintf(stderr, "budget migrate %s %s: %v\n", o.Kind, verb, err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(result); err != nil {
		fmt.Fprintln(stderr, "protocol migration output failed; inspect or replay the exact decision before work")
		return 1
	}
	return 0
}
