package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/protocol"
)

// runLearn implements `parley learn <closed-idea-slug>` (COOPERATION.md §13): distill a
// COMPLETED idea into a reusable, advisory playbook at parley-deck/playbooks/<topic>.md.
// Like `parley retro propose` it is read-only over the idea and writes exactly one new
// file, fail-closed if the target exists (Lstat symlink guard). A playbook is advisory
// (beside consults) — never quorum, never overriding protocol.
func runLearn(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("learn", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace root")
	topic := fs.String("topic", "", "playbook topic filename (default: the idea slug)")
	// Pull the first non-flag arg (the slug) out before flag parsing so
	// `parley learn <slug> --dir X` works (Go's flag stops at the first positional).
	slug := ""
	rest := make([]string, 0, len(args))
	for _, a := range args {
		if slug == "" && !strings.HasPrefix(a, "-") {
			slug = strings.TrimSpace(a)
			continue
		}
		rest = append(rest, a)
	}
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if slug == "" || fs.NArg() != 0 {
		fmt.Fprintln(stderr, "usage: parley learn <closed-idea-slug> [--topic NAME] [--dir DIR]")
		return 2
	}
	if !reSlug.MatchString(slug) {
		fmt.Fprintf(stderr, "learn: invalid idea slug %q (lowercase kebab-case)\n", slug)
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "learn: %v\n", err)
		return 1
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	if _, err := os.Stat(filepath.Join(ideaDir, "00-prompt.md")); err != nil {
		fmt.Fprintf(stderr, "learn: no idea at %s\n", ideaDir)
		return 1
	}

	// Precondition: only a COMPLETED idea is distilled.
	impl, err := os.ReadFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
	if err != nil || !isImplComplete(string(impl)) {
		fmt.Fprintf(stderr, "learn: idea %q is not complete (IMPLEMENTATION.md status must be `complete`) — only closed ideas are distilled\n", slug)
		return 1
	}

	name := strings.TrimSpace(*topic)
	if name == "" {
		name = slug
	}
	if !reSlug.MatchString(name) {
		fmt.Fprintf(stderr, "learn: invalid --topic %q (lowercase kebab-case)\n", name)
		return 2
	}
	playbook := filepath.Join(root, protocol.DeckDir, "playbooks", name+".md")
	// Fail closed if anything already exists at the target (Lstat: a symlink cannot
	// redirect the write). --refresh is a deferred follow-up.
	if _, err := os.Lstat(playbook); err == nil {
		fmt.Fprintf(stderr, "learn: %s already exists — refusing (fail-closed; use --topic or remove it first)\n", playbook)
		return 1
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(stderr, "learn: %v\n", err)
		return 1
	}
	if err := os.MkdirAll(filepath.Dir(playbook), 0o755); err != nil {
		fmt.Fprintf(stderr, "learn: %v\n", err)
		return 1
	}

	content := distillPlaybook(name, slug, ideaDir)
	if err := os.WriteFile(playbook, []byte(content), 0o644); err != nil {
		fmt.Fprintf(stderr, "learn: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Wrote advisory playbook %s (distilled from ideas/%s).\n", filepath.Join(protocol.DeckDir, "playbooks", name+".md"), slug)
	fmt.Fprintln(stdout, "Playbooks are advisory (beside consults): reference one in Phase 0 for context; it never counts as quorum or overrides the protocol.")
	return 0
}

func isImplComplete(impl string) bool {
	fm := extractFM(impl)
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "status:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "status:")) == "complete"
		}
	}
	return false
}

func extractFM(doc string) string {
	s := strings.TrimLeft(doc, " \t\n")
	if !strings.HasPrefix(s, "---") {
		return ""
	}
	s = s[3:]
	if end := strings.Index(s, "\n---"); end >= 0 {
		return s[:end]
	}
	return ""
}

// distillPlaybook builds the advisory playbook from the closed idea's artifacts. v1 is
// a deterministic skeleton: it records the proven shape (track, participants, phases,
// fix-up cycles) and points the reader at the source idea for specifics; the human
// refines the transferable prose in review before committing.
func distillPlaybook(name, slug, ideaDir string) string {
	track := frontmatterValue(filepath.Join(ideaDir, "00-prompt.md"), "track")
	participants := frontmatterValue(filepath.Join(ideaDir, "00-prompt.md"), "participants")
	fixups := strings.Count(readFile(filepath.Join(ideaDir, "IMPLEMENTATION.md")), "## Fix-up cycle")
	var b strings.Builder
	fmt.Fprintf(&b, `---
playbook: %s
distilled-from: ideas/%s
distilled: (set at commit)
status: advisory
---

## When to use

Work resembling **%s** — see ideas/%s for the concrete precedent. (Generalize this line:
what class of task does this playbook cover? Strip the idea-specific specifics.)

## Proven shape

- Track: %s
- Participants: %s
- Fix-up cycles taken: %d
- Lifecycle actually run: round-01 → cross-review (if divergent) → consensus + signoffs
  → FINAL → implement → refutation review → fix-up → complete.

## Step checklist

- [ ] Frame the idea narrowly in 00-prompt.md; set the track deliberately.
- [ ] Independent round-01 from every participant before anyone reads the others.
- [ ] Open a cross-review round only where positions genuinely diverge.
- [ ] Draft consensus; collect every participant's signoff; then FINAL.
- [ ] Implement strictly to FINAL; record deviations in IMPLEMENTATION.md.
- [ ] Refutation-default review; fix-up until zero agreed findings; then complete.

## Gotchas & fixes

(Fill from this idea's review consensus + IMPLEMENTATION.md deviations — the mistakes
that were caught and how. Keep only the transferable ones.)

## Verification pattern

(How "done" was proven — the checks/evidence used. If a completion contract was used,
note the shape here.)
`, name, slug, slug, slug, orDash(track), orDash(participants), fixups)
	return b.String()
}

func frontmatterValue(path, key string) string {
	fm := extractFM(readFile(path))
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key+":") {
			return strings.TrimSpace(strings.TrimPrefix(line, key+":"))
		}
	}
	return ""
}

func readFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}
