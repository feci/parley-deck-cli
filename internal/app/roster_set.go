package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// rosterSetField is one requested change to one roster entry.
type rosterSetField struct{ key, value string }

// rosterSet patches a single [roster.<id>] block in exactly one config file.
//
// Scope names a FILE, not a vague locality: `deck` writes the COMMITTED
// parley-deck/agents.toml, never the gitignored agents.local.toml — a roster change
// invisible to the repository is how a deck silently diverges from its own history, which
// is the failure this whole change exists to end. `machine` writes ~/.parley/agents.toml.
func rosterSet(root, scope, agent string, fields []rosterSetField, dryRun, yes bool, stdout, stderr io.Writer) int {
	if strings.TrimSpace(agent) == "" {
		fmt.Fprintln(stderr, "roster set: AGENT is required")
		return 2
	}
	if len(fields) == 0 {
		fmt.Fprintln(stderr, "roster set: nothing to change — pass at least one of --adapter/--state/--model/--effort/--speed")
		return 2
	}
	target, err := rosterScopeFile(root, scope)
	if err != nil {
		fmt.Fprintf(stderr, "roster set: %v\n", err)
		return 2
	}

	existing, err := os.ReadFile(target)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(stderr, "roster set: %v\n", err)
		return 1
	}
	updated, changes, err := applyRosterBlock(string(existing), agent, fields)
	if err != nil {
		fmt.Fprintf(stderr, "roster set: %v\n", err)
		return 1
	}

	if len(changes) == 0 {
		fmt.Fprintf(stdout, "roster set: %s already matches in %s — nothing to do\n", agent, target)
		return 0
	}
	fmt.Fprintf(stdout, "%s  [roster.%s] in %s\n", previewLabel(dryRun, yes), agent, target)
	for _, c := range changes {
		fmt.Fprintf(stdout, "  %s\n", c)
	}
	// Preview is the default. A mutation happens only on an explicit --yes, so an
	// unattended invocation can never rewrite a roster by accident.
	if dryRun || !yes {
		fmt.Fprintln(stdout, "\nNothing was written. Re-run with --yes to apply.")
		return 0
	}
	if err := writeRosterFileAtomic(target, []byte(updated)); err != nil {
		fmt.Fprintf(stderr, "roster set: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "\nWrote %s\n", target)
	return 0
}

func previewLabel(dryRun, yes bool) string {
	if dryRun || !yes {
		return "would change"
	}
	return "changing"
}

// rosterScopeFile maps a scope name to the exact file it writes.
func rosterScopeFile(root, scope string) (string, error) {
	switch scope {
	case "deck", "session": // `session` is the pre-1.40 spelling, kept as a hidden alias
		return filepath.Join(root, "parley-deck", "agents.toml"), nil
	case "machine", "global":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if h := strings.TrimSpace(os.Getenv("PARLEY_HOME")); h != "" {
			home = h
		}
		return filepath.Join(home, ".parley", "agents.toml"), nil
	default:
		return "", fmt.Errorf("invalid --scope %q (want deck|machine)", scope)
	}
}

// applyRosterBlock rewrites (or appends) one [roster.<id>] block in a TOML document,
// preserving every other byte. It is deliberately line-based rather than a
// marshal/unmarshal round-trip: these files carry extensive human comments recording WHY
// a value is pinned, and a round-trip would silently delete all of them.
func applyRosterBlock(doc, agent string, fields []rosterSetField) (string, []string, error) {
	header := fmt.Sprintf("[roster.%s]", agent)
	lines := strings.Split(doc, "\n")

	start, end := -1, len(lines)
	for i, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if trimmed == header {
			start = i
			continue
		}
		if start >= 0 && strings.HasPrefix(trimmed, "[") {
			end = i
			break
		}
	}

	var changes []string
	want := map[string]string{}
	order := make([]string, 0, len(fields))
	for _, f := range fields {
		if _, seen := want[f.key]; !seen {
			order = append(order, f.key)
		}
		want[f.key] = f.value
	}

	if start < 0 {
		// New member. Appending is the only mutation that adds membership, so the
		// caller gates it behind the breaking-change confirmation.
		var b strings.Builder
		b.WriteString(strings.TrimRight(doc, "\n"))
		if strings.TrimSpace(doc) != "" {
			b.WriteString("\n\n")
		}
		b.WriteString(header + "\n")
		for _, k := range order {
			b.WriteString(fmt.Sprintf("%s = %s\n", k, tomlValue(want[k])))
			changes = append(changes, fmt.Sprintf("+ %s = %s", k, tomlValue(want[k])))
		}
		return b.String(), changes, nil
	}

	block := lines[start+1 : end]
	seen := map[string]bool{}
	for i, ln := range block {
		trimmed := strings.TrimSpace(ln)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		key, _, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		newVal, requested := want[key]
		if !requested {
			continue
		}
		seen[key] = true
		formatted := fmt.Sprintf("%s = %s", key, tomlValue(newVal))
		if trimmed == formatted {
			continue
		}
		changes = append(changes, fmt.Sprintf("- %s\n  + %s", trimmed, formatted))
		block[i] = formatted
	}
	var added []string
	for _, k := range order {
		if seen[k] {
			continue
		}
		added = append(added, fmt.Sprintf("%s = %s", k, tomlValue(want[k])))
		changes = append(changes, fmt.Sprintf("+ %s = %s", k, tomlValue(want[k])))
	}
	sort.Strings(added)

	// New keys go at the end of the block's content, before any trailing blank lines,
	// so a block does not grow a gap every time it is edited.
	tail := len(block)
	for tail > 0 && strings.TrimSpace(block[tail-1]) == "" {
		tail--
	}
	out := append([]string{}, lines[:start+1]...)
	out = append(out, block[:tail]...)
	out = append(out, added...)
	out = append(out, block[tail:]...)
	out = append(out, lines[end:]...)
	return strings.Join(out, "\n"), changes, nil
}

// tomlValue renders a value as TOML. `active` is the only boolean in the block.
func tomlValue(v string) string {
	if v == "true" || v == "false" {
		return v
	}
	return fmt.Sprintf("%q", v)
}

// writeRosterFileAtomic writes via a temp file in the same directory and renames, so a
// crash mid-write cannot leave a half-written roster behind.
func writeRosterFileAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".agents-*.toml")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}
