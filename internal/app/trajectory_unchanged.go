package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"parley-deck-cli/internal/trajectory"
)

func runTrajectoryUnchanged(ctx context.Context, args []string, out, errout io.Writer) int {
	f := flag.NewFlagSet("trajectory reconcile-unchanged", flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "original trajectory worktree")
	idea := f.String("idea", "", "idea identity")
	sequence := f.Int("sequence", 0, "original charged attempt sequence")
	expected := f.String("sha256", "", "exact original unchanged observation preview digest")
	yes := f.Bool("yes", false, "publish this exact observation without executing work")
	if err := f.Parse(args); err != nil || f.NArg() != 0 || !trajectoryPathID(*idea) || *sequence < 1 || *sequence > trajectory.MaxPatches || (!*yes && *expected != "") || (*yes && *expected == "") {
		return 2
	}
	var p trajectory.ReconciliationPreview
	var err error
	if *yes {
		p, err = trajectory.ReconcileUnchanged(ctx, *root, *idea, *sequence, *expected)
	} else {
		p, err = trajectory.PreviewUnchanged(ctx, *root, *idea, *sequence)
	}
	if err != nil {
		fmt.Fprintf(errout, "trajectory unchanged reconciliation: %v\n", err)
		return 1
	}
	if err = json.NewEncoder(out).Encode(struct {
		Preview trajectory.ReconciliationPreview `json:"preview"`
		SHA256  string                           `json:"sha256"`
		Applied bool                             `json:"applied"`
	}{p, p.SHA256(), *yes}); err != nil {
		fmt.Fprintln(errout, "cannot write unchanged reconciliation result; preserve history and repeat the exact request")
		return 1
	}
	return 0
}
