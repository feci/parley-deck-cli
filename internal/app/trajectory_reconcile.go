package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"parley-deck-cli/internal/trajectory"
)

func runTrajectoryReconciliation(ctx context.Context, args []string, out, errout io.Writer, attended bool) int {
	if len(args) == 0 {
		return 2
	}
	f := flag.NewFlagSet("trajectory "+args[0], flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "original trajectory worktree")
	idea := f.String("idea", "", "idea identity")
	run := f.String("run", "", "original verifier run identity")
	expected := f.String("sha256", "", "exact preview digest")
	id := f.String("decision-id", "", "unique attended continuation decision")
	reason := f.String("reason", "", "reason for the continuation")
	review := f.Bool("acknowledge-review", false, "explicitly continue past the reported regression review gate")
	inconclusive := f.Bool("acknowledge-inconclusive", false, "explicitly continue while retaining the inconclusive assessment")
	yes := f.Bool("yes", false, "apply this exact preview")
	if err := f.Parse(args[1:]); err != nil || f.NArg() != 0 || !trajectoryPathID(*idea) {
		return 2
	}
	fail := func(err error) int { fmt.Fprintf(errout, "trajectory: %v\n", err); return 1 }
	encode := func(value any) int {
		if err := json.NewEncoder(out).Encode(value); err != nil {
			return fail(err)
		}
		return 0
	}
	if args[0] == "history" {
		if *yes || *run != "" || *expected != "" || *id != "" || *reason != "" || *review || *inconclusive {
			return 2
		}
		h, err := trajectory.InspectHistory(ctx, *root, *idea)
		if err != nil {
			return fail(err)
		}
		return encode(h)
	}
	if args[0] == "reconcile" {
		if !trajectoryPathID(*run) || *id != "" || *reason != "" || *review || *inconclusive || (!*yes && *expected != "") {
			return 2
		}
		p, err := trajectory.PreviewReconciliation(ctx, *root, *idea, *run)
		if err != nil {
			return fail(err)
		}
		if *yes {
			if err = trajectory.Reconcile(ctx, *root, *idea, *run, *expected); err != nil {
				return fail(err)
			}
		}
		return encode(struct {
			Preview trajectory.ReconciliationPreview `json:"preview"`
			SHA256  string                           `json:"sha256"`
			Applied bool                             `json:"applied"`
		}{p, p.SHA256(), *yes})
	}
	if args[0] != "continue" || *run != "" || (!*yes && (*expected != "" || *id != "" || *reason != "" || *review || *inconclusive)) {
		return 2
	}
	if *yes && !attended {
		fmt.Fprintln(errout, "trajectory continuation requires the attended operator control; no decision was applied")
		return 2
	}
	var p trajectory.ContinuationPreview
	var err error
	if *yes {
		p, err = trajectory.Continue(ctx, *root, *idea, *expected, *id, *reason, *review, *inconclusive)
	} else {
		p, err = trajectory.PreviewContinuation(ctx, *root, *idea)
	}
	if err != nil {
		return fail(err)
	}
	return encode(struct {
		Preview trajectory.ContinuationPreview `json:"preview"`
		SHA256  string                         `json:"sha256"`
		Applied bool                           `json:"applied"`
	}{p, p.SHA256(), *yes})
}
