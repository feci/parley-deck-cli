package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/retro"
)

// runRetro implements `parley retro` (COOPERATION.md §13): read-only mining of the
// deck's own structured history to PROPOSE improvements. scan/select/diagnose are
// strictly read-only; propose may write only a single new ideas/<slug>/00-prompt.md
// (fail-if-exists) and nothing else. It never edits the protocol or any harness
// file and never writes another agent's artifact.
func runRetro(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: parley retro <scan|select|diagnose|propose> [--dir DIR] [--k N] [--json] [--slug SLUG]")
		return 2
	}
	sub := args[0]
	fs := flag.NewFlagSet("retro", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace root")
	k := fs.Int("k", 10, "coreset size (select/diagnose/propose)")
	jsonOut := fs.Bool("json", false, "JSON output (scan/select)")
	slug := fs.String("slug", "", "explicit, non-existing idea slug (propose)")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "retro: %v\n", err)
		return 1
	}

	signals, err := retro.Scan(root)
	if err != nil {
		fmt.Fprintf(stderr, "retro scan: %v\n", err)
		return 1
	}

	switch sub {
	case "scan":
		if *jsonOut {
			return printJSON(stdout, signals, stderr)
		}
		printSignals(stdout, signals)
		return 0
	case "select":
		coreset := retro.Select(signals, *k)
		if *jsonOut {
			return printJSON(stdout, coreset, stderr)
		}
		fmt.Fprintf(stdout, "Coreset (%d hard, type-diverse cases):\n\n", len(coreset))
		printSignals(stdout, coreset)
		return 0
	case "diagnose":
		fmt.Fprint(stdout, retro.Diagnose(retro.Select(signals, *k)))
		return 0
	case "propose":
		return retroPropose(root, *slug, retro.Select(signals, *k), stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown retro subcommand: %s (want scan|select|diagnose|propose)\n", sub)
		return 2
	}
}

func printSignals(w io.Writer, signals []retro.IdeaSignals) {
	if len(signals) == 0 {
		fmt.Fprintln(w, "No ideas found under parley-deck/ideas/.")
		return
	}
	fmt.Fprintf(w, "%-7s %-13s %s\n", "SCORE", "FAILURE-TYPE", "IDEA")
	for _, s := range signals {
		fmt.Fprintf(w, "%-7.1f %-13s %s\n", s.Score, s.FailureType, s.Slug)
	}
}

// retroPropose scaffolds a single new ideas/<slug>/00-prompt.md seeded from the
// diagnosis. It writes nothing else and fails closed if the target exists.
func retroPropose(root, slug string, coreset []retro.IdeaSignals, stdout, stderr io.Writer) int {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		fmt.Fprintln(stderr, "retro propose: --slug is required (an explicit, non-existing idea slug)")
		return 2
	}
	if slug != filepath.Base(slug) || strings.HasPrefix(slug, ".") {
		fmt.Fprintf(stderr, "retro propose: invalid slug %q\n", slug)
		return 2
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	promptPath := filepath.Join(ideaDir, "00-prompt.md")
	if _, err := os.Stat(promptPath); err == nil {
		fmt.Fprintf(stderr, "retro propose: %s already exists — refusing to overwrite (fail-closed)\n", promptPath)
		return 1
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(stderr, "retro propose: %v\n", err)
		return 1
	}
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "retro propose: %v\n", err)
		return 1
	}
	if err := os.WriteFile(promptPath, []byte(retroPromptBody(slug, coreset)), 0o644); err != nil {
		fmt.Fprintf(stderr, "retro propose: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Wrote %s\n", promptPath)
	fmt.Fprintln(stdout, "This is advisory retro input (§13): run it as a normal Parley Deck idea; nothing is applied automatically.")
	return 0
}

func retroPromptBody(slug string, coreset []retro.IdeaSignals) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("idea: " + slug + "\n")
	b.WriteString("author: claude\n")
	b.WriteString("created: <fill: date>\n")
	b.WriteString("participants: [claude, codex, agy, hermes]\n")
	b.WriteString("status: round-01\n")
	b.WriteString("drafted-by: parley retro\n")
	b.WriteString("---\n\n")
	b.WriteString("## Problem / idea\n\n")
	b.WriteString("Retrospective optimization pass (COOPERATION.md §13). This kickoff was drafted")
	b.WriteString(" by `parley retro` from the deck's own structured history and is **advisory")
	b.WriteString(" input only** — diagnose the cases below independently in round-01, then")
	b.WriteString(" decide which (if any) improvements to propose. Nothing is applied")
	b.WriteString(" automatically; any change goes through the normal gate (and a")
	b.WriteString(" meta-protocol-change idea for protocol text).\n\n")
	b.WriteString(retro.Diagnose(coreset))
	b.WriteString("\n## Constraints\n\n")
	b.WriteString("- Treat the diagnosis as hypotheses, not findings. Each participant diagnoses")
	b.WriteString(" the coreset independently (§13.4 multi-agent diagnosis).\n")
	b.WriteString("- Acceptance is the normal gate: consensus + all-participant signoff + human")
	b.WriteString(" approval for protocol/shared-harness + no-regression (§13.3).\n\n")
	b.WriteString("## Non-goals\n\n")
	b.WriteString("- Auto-applying any edit. Editing the protocol or any harness file directly")
	b.WriteString(" from this pass.\n")
	return b.String()
}
