package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/fsutil"
)

// rosterSetField is one requested change to one roster entry.
type rosterSetField struct{ key, value string }

// rosterSet patches a single [roster.<id>] block in exactly one config file.
//
// Scope names a FILE, not a vague locality: `deck` writes the COMMITTED
// parley-deck/agents.toml, never the gitignored agents.local.toml — a roster change
// invisible to the repository is how a deck silently diverges from its own history, which
// is the failure this whole change exists to end. `machine` writes ~/.parley/agents.toml.
func rosterSet(root, scope, agent string, fields []rosterSetField, dryRun, yes, confirmBreaking bool, stdout, stderr io.Writer) int {
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
	// Adding a member or retiring one changes WHO deliberates, and therefore who a future
	// idea's quorum is. `--yes` is the ordinary confirmation; a membership change needs a
	// second, explicit one so it can never ride along with a routine model change.
	if breaking := membershipChange(changes, blockExists(string(existing), agent), priorActiveIn(target, agent)); breaking != "" {
		if !confirmBreaking {
			fmt.Fprintf(stderr, "\nroster set: this %s — a membership change, not a settings change.\n"+
				"Re-run with --confirm-breaking as well as --yes.\n", breaking)
			return 2
		}
		fmt.Fprintf(stdout, "\n(%s — confirmed with --confirm-breaking)\n", breaking)
	}
	if err := writeRosterFileAtomic(target, []byte(updated)); err != nil {
		fmt.Fprintf(stderr, "roster set: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "\nWrote %s\n", target)
	// POST-WRITE RE-RESOLVE. A write to a lower layer can be completely masked by a
	// higher one ($PARLEY_HEADLESS_AGENT_CONFIG, agents.local.toml). Reporting "Wrote"
	// and stopping there is a false success: the file changed and the effective value
	// did not. `masked-by-env` was in the frozen STATUS vocabulary with nothing to emit
	// it; this is the emitter.
	for _, f := range fields {
		if src, masked := rosterFieldMaskedBy(root, agent, f.key, target); masked {
			fmt.Fprintf(stderr, "\nwarning: %s = %q is MASKED — %s sets it at a higher layer, so the effective value did not change.\n"+
				"  (status `masked-by-env`; see `parley roster show --explain %s`)\n", f.key, f.value, src, agent)
		}
	}
	return 0
}

// rosterFieldMaskedBy reports whether a higher config layer overrides the field just
// written to `target`, and which one.
func rosterFieldMaskedBy(root, agent, field, target string) (string, bool) {
	sources, err := config.RosterFieldSources(root, agent)
	if err != nil {
		return "", false
	}
	if field == "active" {
		// State is decided by the membership authority, not by the layer stack, so a
		// higher layer's `active` cannot mask a write to the authority. Warning here
		// claimed the opposite of what `roster show` then reported.
		// The authority depends on the scope being written: a machine-scope write is
		// governed by the machine roster, a deck-scope write by the deck's.
		authority, aerr := config.RosterStateSourceForTarget(root, target)
		if aerr != nil || authority == "" {
			return "", false
		}
		ap := config.RosterSourcePath(root, authority)
		if ap == "" {
			ap = authority
		}
		tp, _ := filepath.Abs(target)
		app, _ := filepath.Abs(ap)
		return authority, tp != app
	}
	src := sources[field]
	if src == "" {
		return "", false
	}
	// RosterFieldSources reports the LAST layer that set the field, as a DISPLAY LABEL
	// ("~/.parley/agents.toml", "parley-deck/agents.toml"). Comparing that label against
	// an absolute path made every machine-scope write claim it was masking itself.
	// Resolve the label to the path it names and compare paths.
	winner := config.RosterSourcePath(root, src)
	if winner == "" {
		return "", false
	}
	tp, _ := filepath.Abs(target)
	wp, _ := filepath.Abs(winner)
	return src, tp != wp
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
		// Ask the config loader, never reconstruct the path. PARLEY_HOME names the central
		// config DIRECTORY, not a user home, so composing $PARLEY_HOME/.parley/agents.toml
		// wrote a file no resolver reads — a machine update that reported success and
		// changed nothing.
		path := config.CentralAgentsPath()
		if path == "" {
			return "", fmt.Errorf("cannot resolve the central config directory")
		}
		return path, nil
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
	// Preserve the existing mode. os.CreateTemp makes a 0600 file, and renaming it over a
	// 0644 config silently tightened the machine file to owner-only — harmless for one
	// user, wrong on a shared or team setup, and invisible until someone else's read
	// fails. New files get 0644, the conventional mode for these configs.
	perm := os.FileMode(0o644)
	if fi, err := os.Stat(path); err == nil {
		perm = fi.Mode().Perm()
	}
	return fsutil.WriteFileAtomic(path, data, perm)
}

// membershipChange reports whether a change set alters who is in the roster, rather than
// how an existing member is configured.
//
// The gate keys on whether the BLOCK existed before, not on which field is written. Using
// "+ adapter = " as a proxy for "new member" let `roster set sneaky-9 --model k3 --yes`
// create a real member with no second confirmation: the member is only as new as its
// block, and any first write to a missing block creates one.
func membershipChange(changes []string, existed, priorActive bool) string {
	if !existed {
		return "adds a new roster member"
	}
	for _, c := range changes {
		// Gate on an actual STATE FLIP. Writing `active = true` to a block that had no
		// `active` key is a no-op — absence already means active — and demanding a
		// membership confirmation for it trains operators to pass --confirm-breaking
		// reflexively, which is how a gate stops being one.
		if strings.Contains(c, "+ active = true") && !priorActive {
			return "reactivates a retired roster member"
		}
		if strings.Contains(c, "+ active = false") && priorActive {
			return "retires a roster member"
		}
	}
	return ""
}

// priorActiveIn reports the member's state in the file being edited, before this write.
// Absence of the key (or of the block) means active, matching the field table.
func priorActiveIn(path, agent string) bool {
	entries, err := config.RosterEntriesInFile(path)
	if err != nil {
		return true
	}
	e, ok := entries[agent]
	if !ok {
		return true
	}
	return e.Active
}

// blockExists reports whether [roster.<agent>] is already declared in the file being
// edited. It is deliberately a text check against the SAME bytes applyRosterBlock
// patches, so the gate and the write can never disagree about what existed.
func blockExists(doc, agent string) bool {
	header := "[roster." + agent + "]"
	for _, line := range strings.Split(doc, "\n") {
		if strings.TrimSpace(line) == header {
			return true
		}
	}
	return false
}
