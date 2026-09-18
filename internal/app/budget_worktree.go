package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"parley-deck-cli/internal/budget"
)

// Read-only worktree registration inventory. "inspect" is the only verb: this
// surface has no apply, attest, prune, repair or recovery path, and it takes no
// attendance gate because it changes nothing. It reports unavailable
// registrations instead of pruning them and leaves every other budget control,
// including the existing refusal on unavailable launch history, untouched.
func runBudgetWorktree(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	usage := "usage: parley budget worktree inspect --dir DIR (read-only inventory; this surface has no mutating verb)"
	if len(args) == 0 || args[0] != "inspect" {
		fmt.Fprintln(stderr, usage)
		return 2
	}
	f := flag.NewFlagSet("budget worktree inspect", flag.ContinueOnError)
	f.SetOutput(stderr)
	// --dir is the only flag this surface defines, so every other flag —
	// including any decision, expectation or apply flag — is rejected by Parse.
	root := f.String("dir", ".", "workspace to inspect; every registration of its repository is reported")
	if err := f.Parse(args[1:]); err != nil {
		return 2
	}
	if f.NArg() != 0 {
		fmt.Fprintln(stderr, "budget worktree inspect takes no positional arguments; it cannot be given an action to perform")
		return 2
	}
	inventory, err := budget.InspectWorktreeInventory(ctx, *root)
	if err != nil {
		fmt.Fprintf(stderr, "budget worktree inspect: %v\n", err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(inventory); err != nil {
		fmt.Fprintln(stderr, "budget worktree inspect: cannot write result")
		return 1
	}
	return 0
}
