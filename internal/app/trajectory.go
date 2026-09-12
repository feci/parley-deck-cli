package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/trajectory"
)

func runTrajectory(ctx context.Context, args []string, out, errout io.Writer) int {
	if len(args) > 0 && args[0] == "verify" {
		return runTrajectoryVerify(ctx, args[1:], out, errout)
	}
	if len(args) > 0 && args[0] == "verify-helper" {
		return runTrajectoryVerifyHelper(ctx, args[1:], out, errout)
	}
	supported, attended := budgetAttendance()
	if len(args) > 0 && (args[0] == "reconcile" || args[0] == "continue" || args[0] == "history") {
		return runTrajectoryReconciliation(ctx, args, out, errout, supported && attended)
	}
	return runTrajectoryControl(ctx, args, out, errout, supported && attended)
}

// Only the platform attendance probe supplies the production bool. This CLI
// activates capture, never an operator extension or a verification verdict.
func runTrajectoryControl(ctx context.Context, args []string, out, errout io.Writer, attended bool) int {
	if len(args) == 0 || (args[0] != "initialize" && args[0] != "configure" && args[0] != "inspect") {
		fmt.Fprintln(errout, "usage: parley trajectory initialize --dir DIR --idea ID --max-fixups N --yes; parley trajectory configure --dir DIR --idea ID --implementer ID [--policy-sha256 SHA --cycle-policy-sha256 SHA --yes]; parley trajectory inspect --dir DIR --idea ID; parley trajectory verify --dir DIR --idea ID --verifier ID --yes")
		fmt.Fprintln(errout, "       parley trajectory reconcile --dir DIR --idea ID --run RUN [--sha256 SHA --yes]; parley trajectory history --dir DIR --idea ID; parley trajectory continue --dir DIR --idea ID [--sha256 SHA --decision-id ID --reason REASON --yes] [--acknowledge-review] [--acknowledge-inconclusive]")
		return 2
	}
	f := flag.NewFlagSet("trajectory "+args[0], flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "Git worktree root")
	idea := f.String("idea", "", "idea identity")
	implementer := f.String("implementer", "", "frozen patch author")
	expected := f.String("policy-sha256", "", "exact trajectory preview digest")
	cycle := f.String("cycle-policy-sha256", "", "exact original cycle policy digest")
	maximum := f.Int("max-fixups", 0, "initial protocol fixup ceiling (must match the selected track)")
	yes := f.Bool("yes", false, "activate this exact attended policy")
	if err := f.Parse(args[1:]); err != nil || f.NArg() != 0 || *idea == "" {
		return 2
	}
	fail := func(err error) int { fmt.Fprintf(errout, "trajectory: %v\n", err); return 1 }
	encode := func(v any) int {
		if err := json.NewEncoder(out).Encode(v); err != nil {
			return fail(err)
		}
		return 0
	}
	if args[0] == "inspect" {
		if *yes || *implementer != "" || *expected != "" || *cycle != "" || *maximum != 0 {
			return 2
		}
		s, err := trajectory.Inspect(ctx, *root, *idea)
		if err != nil {
			return fail(err)
		}
		return encode(s)
	}
	if *yes && !attended {
		fmt.Fprintln(errout, "trajectory: activation requires the attended operator control; no policy was changed")
		return 2
	}
	if args[0] == "initialize" {
		if !*yes || *maximum <= 0 || *implementer != "" || *expected != "" || *cycle != "" {
			return 2
		}
		b, err := budget.EnsureCycleBinding(ctx, *root, *idea, budget.Fixup, *maximum, 0, "", filepath.Join(*root, "parley-deck", "ideas", *idea))
		if err != nil {
			return fail(err)
		}
		s, err := b.Inspect(ctx)
		if err != nil {
			return fail(err)
		}
		return encode(s)
	}
	if *maximum != 0 || *implementer == "" {
		return 2
	}
	criteria, isList, err := driver.ReadChecksContract(filepath.Join(*root, "parley-deck", "ideas", *idea))
	if err != nil {
		return fail(err)
	}
	if !isList || len(criteria) == 0 {
		return fail(fmt.Errorf("trajectory requires the original named material checks"))
	}
	checks := make([]trajectory.Criterion, len(criteria))
	for i, c := range criteria {
		checks[i] = trajectory.Criterion{Name: c.Name, Command: c.Command}
	}
	p, cycleSHA, err := trajectory.NewPolicy(ctx, *root, *idea, *implementer, checks)
	if err != nil {
		return fail(err)
	}
	sha, err := p.SHA256()
	if err != nil {
		return fail(err)
	}
	if *yes {
		if *expected != sha || *cycle != cycleSHA {
			return fail(fmt.Errorf("trajectory or original cycle policy differs from the preview; inspect again"))
		}
		if err := trajectory.Activate(ctx, *root, cycleSHA, p); err != nil {
			return fail(err)
		}
	} else if *expected != "" || *cycle != "" {
		return 2
	}
	return encode(struct {
		Policy    trajectory.Policy `json:"policy"`
		SHA       string            `json:"policy_sha256"`
		CycleSHA  string            `json:"cycle_policy_sha256"`
		Activated bool              `json:"activated"`
	}{p, sha, cycleSHA, *yes})
}
