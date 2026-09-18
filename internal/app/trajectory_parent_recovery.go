package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"parley-deck-cli/internal/trajectory"
)

// Deterministic recovery publishes retained observations only. It does not
// activate a policy or issue the separately attended continuation decision.
func runTrajectoryParentRecovery(ctx context.Context, args []string, out, errout io.Writer) int {
	f := flag.NewFlagSet("trajectory recover-parent", flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "original trajectory worktree")
	idea := f.String("idea", "", "idea identity")
	run := f.String("run", "", "original independent verifier run")
	expected := f.String("sha256", "", "exact original recovery preview digest")
	yes := f.Bool("yes", false, "publish this exact derived observation without executing work")
	if err := f.Parse(args); err != nil || f.NArg() != 0 || !trajectoryPathID(*idea) || !trajectoryPathID(*run) || (!*yes && *expected != "") || (*yes && *expected == "") {
		return 2
	}
	var p trajectory.ParentRecoveryPreview
	var err error
	if *yes {
		p, err = trajectory.RecoverParent(ctx, *root, *idea, *run, *expected)
	} else {
		p, err = trajectory.PreviewParentRecovery(ctx, *root, *idea, *run)
	}
	if err != nil {
		fmt.Fprintf(errout, "trajectory parent recovery: %v\n", err)
		return 1
	}
	if err = json.NewEncoder(out).Encode(struct {
		Preview trajectory.ParentRecoveryPreview `json:"preview"`
		SHA256  string                           `json:"sha256"`
		Applied bool                             `json:"applied"`
	}{p, p.SHA256(), *yes}); err != nil {
		fmt.Fprintln(errout, "cannot write parent recovery result; preserve the record and repeat the exact request")
		return 1
	}
	return 0
}
