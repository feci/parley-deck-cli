package app

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
)

// runPreset implements `parley preset list` (named-roster-presets): print the roster
// presets available across the layered config, their winning source layer, and a flag
// for any preset that references an agent missing from or inactive in the §2 roster.
func runPreset(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("preset", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	// Accept a leading "list" subcommand in any position: strip it before flag
	// parsing so `parley preset list --dir X` works (Go's flag stops at the first
	// positional, which would otherwise swallow flags placed after "list").
	rest := make([]string, 0, len(args))
	sawSub := false
	for _, a := range args {
		if !sawSub && a == "list" {
			sawSub = true
			continue
		}
		rest = append(rest, a)
	}
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(stderr, "usage: parley preset list [--dir DIR]")
		return 2
	}

	rc, err := config.LoadRosterPresets(*root)
	if err != nil {
		fmt.Fprintf(stderr, "roster presets failed: %v\n", err)
		return 1
	}
	if len(rc.Presets) == 0 {
		fmt.Fprintln(stdout, "No roster presets defined. Add a [rosters.<name>] block to ~/.parley/agents.toml or parley-deck/agents.toml.")
		return 0
	}
	rosterIDs, inactive, ok := protocol.ReadRosterIDs(*root)

	// Reverse the track-default map for a "fits track" hint.
	fits := map[string][]string{}
	for track, preset := range rc.TrackDefault {
		fits[preset] = append(fits[preset], track)
	}

	names := make([]string, 0, len(rc.Presets))
	for n := range rc.Presets {
		names = append(names, n)
	}
	sort.Strings(names)

	fmt.Fprintln(stdout, "Roster presets:")
	for _, n := range names {
		p := rc.Presets[n]
		warn := ""
		if ok {
			var bad []string
			for _, id := range p.Participants {
				switch {
				case !rosterIDs[id]:
					bad = append(bad, id+" (not in §2)")
				case inactive[id]:
					bad = append(bad, id+" (inactive)")
				}
			}
			if len(bad) > 0 {
				warn = "  ⚠ " + strings.Join(bad, ", ")
			}
		}
		track := ""
		if t := fits[n]; len(t) > 0 {
			sort.Strings(t)
			track = "  [track: " + strings.Join(t, ",") + "]"
		}
		fmt.Fprintf(stdout, "  %-14s %d agents  (%s)%s%s\n",
			n, len(p.Participants), p.Source, track, warn)
		fmt.Fprintf(stdout, "                 %s\n", strings.Join(p.Participants, ", "))
	}
	return 0
}
