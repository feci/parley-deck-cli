package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
)

func runQuota(ctx context.Context, args []string, out, errs io.Writer) int {
	if len(args) == 0 || (args[0] != "revise" && args[0] != "recover") {
		fmt.Fprintln(errs, "usage: parley quota revise|recover --dir DIR --idea SLUG --run RUN [--request FILE]")
		return 2
	}
	fs := flag.NewFlagSet("quota "+args[0], flag.ContinueOnError)
	fs.SetOutput(errs)
	root := fs.String("dir", ".", "workspace directory")
	slug := fs.String("idea", "", "idea slug")
	run := fs.String("run", "", "existing idea-bound run")
	request := fs.String("request", "", "JSON revision request with committed owner authority")
	if fs.Parse(args[1:]) != nil {
		return 2
	}
	idPattern := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
	if !idPattern.MatchString(*slug) || !idPattern.MatchString(*run) || fs.NArg() != 0 {
		fmt.Fprintln(errs, "quota requires a single idea slug and existing run identity")
		return 2
	}
	dir := filepath.Join(*root, protocol.DeckDir, "ideas", *slug)
	if args[0] == "recover" {
		scoped, release, err := membership.Acquire(ctx, dir, *run)
		if err != nil {
			fmt.Fprintln(errs, err)
			return 1
		}
		defer release()
		h, err := membership.Before(scoped, *root, dir, *run)
		if err != nil {
			fmt.Fprintln(errs, err)
			return 1
		}
		if h == nil {
			fmt.Fprintln(errs, "idea has no recorded quota history")
			return 1
		}
		fmt.Fprintf(out, "Reconciled %s at revision %d; current participants: %s\n", *slug, h.Revision, strings.Join(h.Current, ", "))
		return 0
	}
	f, err := os.Open(*request)
	if err != nil {
		fmt.Fprintln(errs, err)
		return 1
	}
	defer f.Close()
	var r struct {
		Participants []string          `json:"participants"`
		Policy       quota.Policy      `json:"policy"`
		Authority    quota.Authority   `json:"authority"`
		Catchup      map[string]string `json:"catchup,omitempty"`
	}
	raw, err := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if err != nil || len(raw) > 1<<20 {
		fmt.Fprintln(errs, "revision request unreadable or too large")
		return 1
	}
	if err = quota.DecodeRevisionRequest(raw, &r); err != nil {
		fmt.Fprintln(errs, err)
		return 1
	}
	// Require the caller's independently reviewable hashes; never manufacture owner
	// authority or infer a requested scope from the current binary's defaults.
	if err = quota.ValidateAuthority(*root, *slug, r.Authority); err != nil {
		fmt.Fprintln(errs, err)
		return 1
	}
	b, err := membership.Revise(ctx, *root, dir, *run, r.Participants, r.Policy, r.Authority, r.Catchup)
	if err != nil {
		fmt.Fprintln(errs, err)
		return 1
	}
	fmt.Fprintf(out, "Recorded owner revision %s; current participants: %s; policy enabled=%t scope=%s\n", b.ID, strings.Join(b.Decision.After, ", "), b.Policy.Enabled, b.Policy.Scope)
	return 0
}
