package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"parley-deck-cli/internal/budget"
)

// Accounting replay is read-only on every platform. It never activates an
// attempt or substitutes an operator/participant completion decision.
func runBudgetAction(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "inspect" && args[0] != "replay" {
		fmt.Fprintln(stderr, "usage: parley budget action inspect|replay --ledger DIR --scope SCOPE [--entry-key KEY --expected-identity-sha256 HASH]")
		return 2
	}
	f := flag.NewFlagSet("budget action "+args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	dir := f.String("ledger", "", "existing ledger directory")
	scope := f.String("scope", "", "exact existing accounting scope")
	entry := f.String("entry-key", "", "hashed original ledger entry identity")
	expected := f.String("expected-identity-sha256", "", "exact identity digest returned by inspection")
	if e := f.Parse(args[1:]); e != nil {
		return 2
	}
	if f.NArg() != 0 || *dir == "" || *scope == "" {
		fmt.Fprintln(stderr, "action control requires --ledger and --scope, with no positional arguments")
		return 2
	}
	var result any
	var err error
	s := budget.Store{Dir: *dir, Scope: *scope}
	if args[0] == "inspect" {
		if *expected != "" {
			fmt.Fprintln(stderr, "inspect does not accept an expected replay identity")
			return 2
		}
		rows, e := s.InspectActions(ctx)
		err = e
		result = rows
		if err == nil && *entry != "" {
			found := false
			for _, row := range rows {
				if row.EntryKey == *entry {
					result = row
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("action entry is missing; absence does not authorize a new execution")
			}
		}
	} else {
		if *entry == "" || *expected == "" {
			fmt.Fprintln(stderr, "replay requires --entry-key and --expected-identity-sha256; it returns accounting only")
			return 2
		}
		result, err = s.ReplayAction(ctx, *entry, *expected)
	}
	if err != nil {
		fmt.Fprintf(stderr, "budget action %s: %v\n", args[0], err)
		return 1
	}
	if e := json.NewEncoder(stdout).Encode(result); e != nil {
		fmt.Fprintln(stderr, "cannot write action accounting result; repeating this read grants no execution")
		return 1
	}
	return 0
}
