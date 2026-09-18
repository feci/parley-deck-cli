package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"parley-deck-cli/internal/trajectory"
)

func runTrajectoryReservationRecovery(ctx context.Context, args []string, out, errout io.Writer) int {
	f := flag.NewFlagSet("trajectory recover-reservation", flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "original trajectory worktree")
	idea := f.String("idea", "", "original idea")
	entry := f.String("entry", "", "original hashed reservation entry key")
	expected := f.String("sha256", "", "exact recovery preview digest")
	yes := f.Bool("yes", false, "restore only the original charged observation, without executing work")
	if err := f.Parse(args); err != nil || f.NArg() != 0 || !trajectoryPathID(*idea) || len(*entry) != 64 || (!*yes && *expected != "") || (*yes && *expected == "") {
		return 2
	}
	var p trajectory.ReservationRecoveryPreview
	var err error
	if *yes {
		p, err = trajectory.RecoverReservation(ctx, *root, *idea, *entry, *expected)
	} else {
		p, err = trajectory.PreviewReservationRecovery(ctx, *root, *idea, *entry)
	}
	if err != nil {
		fmt.Fprintf(errout, "trajectory reservation recovery: %v\n", err)
		return 1
	}
	if err = json.NewEncoder(out).Encode(struct {
		Preview trajectory.ReservationRecoveryPreview `json:"preview"`
		SHA256  string                                `json:"sha256"`
		Applied bool                                  `json:"applied"`
	}{p, p.SHA256(), *yes}); err != nil {
		fmt.Fprintln(errout, "cannot write reservation recovery result; preserve history and replay the exact request")
		return 1
	}
	return 0
}
