package app

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"parley-deck-cli/internal/repomap"
)

func runContext(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printContextUsage(stderr)
		return 2
	}
	switch args[0] {
	case "repo-map":
		return runContextRepoMap(args[1:], stdout, stderr)
	default:
		printContextUsage(stderr)
		return 2
	}
}

func printContextUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: parley context repo-map [--dir DIR] [--format markdown|json] [--max-files N]")
}

func runContextRepoMap(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("context repo-map", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory to map")
	format := fs.String("format", "markdown", "output format: markdown|json")
	maxFiles := fs.Int("max-files", 1000, "maximum number of files to include")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		printContextUsage(stderr)
		return 2
	}

	m, err := repomap.Build(repomap.Options{Root: *root, MaxFiles: *maxFiles})
	if err != nil {
		fmt.Fprintf(stderr, "context repo-map failed: %v\n", err)
		return 1
	}
	switch strings.ToLower(*format) {
	case "markdown", "md":
		if err := repomap.RenderMarkdown(m, stdout); err != nil {
			fmt.Fprintf(stderr, "context repo-map failed: %v\n", err)
			return 1
		}
	case "json":
		if err := repomap.RenderJSON(m, stdout); err != nil {
			fmt.Fprintf(stderr, "context repo-map failed: %v\n", err)
			return 1
		}
	default:
		fmt.Fprintf(stderr, "invalid format %q; expected markdown or json\n", *format)
		return 2
	}
	return 0
}
