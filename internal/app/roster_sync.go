package app

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"parley-deck-cli/internal/config"
)

// rosterSync reconciles a deck with the machine roster in exactly ONE direction:
// machine -> deck. It never copies deck values upward.
//
// The semantics are REBASE, not copy-down: sync REMOVES deck-level roster fields whose
// value merely restates the machine value, so the deck goes back to inheriting. A
// point-in-time copy would freeze the deck the day after it was written, and every later
// central improvement would stop at its door.
//
// A field the deck sets DIFFERENTLY is a deliberate pin and is never touched silently.
// Deliberate pins that sync would drop are enumerated in the preview and the report, and
// `--keep AGENT.FIELD` exempts one explicitly, so re-applying a pin is a checklist rather
// than an archaeological dig.
func rosterSync(root string, keep []string, dryRun, yes bool, stdout, stderr io.Writer) int {
	machine, err := rosterScopeFile(root, "machine")
	if err != nil {
		fmt.Fprintf(stderr, "roster sync: %v\n", err)
		return 1
	}
	deckFile, err := rosterScopeFile(root, "deck")
	if err != nil {
		fmt.Fprintf(stderr, "roster sync: %v\n", err)
		return 1
	}

	machineEntries, err := config.RosterEntriesInFile(machine)
	if err != nil {
		fmt.Fprintf(stderr, "roster sync: reading %s: %v\n", machine, err)
		return 1
	}
	if len(machineEntries) == 0 {
		fmt.Fprintf(stderr, "roster sync: %s declares no [roster.*] entries — nothing to inherit from\n", machine)
		return 1
	}
	deckEntries, err := config.RosterEntriesInFile(deckFile)
	if err != nil {
		fmt.Fprintf(stderr, "roster sync: reading %s: %v\n", deckFile, err)
		return 1
	}

	keepSet := map[string]bool{}
	for _, k := range keep {
		keepSet[strings.ToLower(strings.TrimSpace(k))] = true
	}
	// An unmatched --keep token is almost always a typo, and a typo used to be silent:
	// `--keep kimi-1.modle --yes` protected nothing and removed kimi-1.model anyway. A
	// keep flag is a statement of intent about a specific field, so it must name one.
	usedKeep := map[string]bool{}

	type drop struct{ agent, field, deckVal, machineVal string }
	var redundant, pins, kept []drop

	ids := make([]string, 0, len(deckEntries))
	for id := range deckEntries {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		d := deckEntries[id]
		m, inMachine := machineEntries[id]
		if !inMachine {
			continue // deck-only member: sync has nothing to say about it
		}
		for _, f := range []struct{ name, deckVal, machineVal string }{
			{"adapter", d.Adapter, m.Adapter},
			{"model", d.Model, m.Model},
			{"effort", d.Effort, m.Effort},
			{"speed", d.Speed, m.Speed},
		} {
			if f.deckVal == "" {
				continue // already inheriting
			}
			rec := drop{id, f.name, f.deckVal, f.machineVal}
			switch {
			case keepSet[strings.ToLower(id+"."+f.name)]:
				usedKeep[strings.ToLower(id+"."+f.name)] = true
				kept = append(kept, rec)
			case f.deckVal == f.machineVal:
				redundant = append(redundant, rec)
			default:
				pins = append(pins, rec)
			}
		}
	}

	if unmatched := unmatchedKeeps(keepSet, usedKeep); len(unmatched) > 0 {
		fmt.Fprintf(stderr, "roster sync: --keep names %d field(s) this deck does not override:\n", len(unmatched))
		for _, k := range unmatched {
			fmt.Fprintf(stderr, "  - %s\n", k)
		}
		fmt.Fprintln(stderr, "Nothing was written. Fix the spelling (AGENT.FIELD, e.g. kimi-1.model) or drop the flag.")
		return 2
	}

	if len(redundant) == 0 && len(pins) == 0 {
		fmt.Fprintf(stdout, "roster sync: %s already inherits from %s — nothing to do\n", deckFile, machine)
		return 0
	}

	fmt.Fprintf(stdout, "roster sync (machine -> deck, rebase)\n  machine: %s\n  deck:    %s\n\n", machine, deckFile)
	if len(redundant) > 0 {
		fmt.Fprintln(stdout, "Redundant deck overrides — removing these makes the deck inherit:")
		for _, r := range redundant {
			fmt.Fprintf(stdout, "  - [roster.%s] %s = %q  (same as machine)\n", r.agent, r.field, r.deckVal)
		}
	}
	if len(pins) > 0 {
		// Never silent. A deliberate pin is the whole point of the layering, so removing
		// one has to be a decision the operator can see and reverse.
		fmt.Fprintln(stdout, "\nDELIBERATE PINS this rebase would remove (deck differs from machine):")
		for _, r := range pins {
			fmt.Fprintf(stdout, "  ! [roster.%s] %s = %q  (machine has %q)  — keep with --keep %s.%s\n",
				r.agent, r.field, r.deckVal, r.machineVal, r.agent, r.field)
		}
	}
	if len(kept) > 0 {
		fmt.Fprintln(stdout, "\nExempted by --keep (left untouched):")
		for _, r := range kept {
			fmt.Fprintf(stdout, "  = [roster.%s] %s = %q\n", r.agent, r.field, r.deckVal)
		}
	}

	if dryRun || !yes {
		fmt.Fprintln(stdout, "\nNothing was written. Re-run with --yes to apply.")
		return 0
	}

	doc, err := os.ReadFile(deckFile)
	if err != nil {
		fmt.Fprintf(stderr, "roster sync: %v\n", err)
		return 1
	}
	// BIND APPLY TO THE PREVIEW. Drops were computed from an earlier read; deleting from
	// a second read without checking the values still match means an edit landing between
	// the two reads is silently discarded, atomic rename or not. Refuse instead.
	current, cerr := config.RosterEntriesInFile(deckFile)
	if cerr != nil {
		fmt.Fprintf(stderr, "roster sync: re-reading %s: %v\n", deckFile, cerr)
		return 1
	}
	for _, r := range append(append([]drop{}, redundant...), pins...) {
		if got := rosterFieldValue(current[r.agent], r.field); got != r.deckVal {
			fmt.Fprintf(stderr, "roster sync: %s changed since the preview ([roster.%s] %s is now %q, was %q).\n"+
				"Nothing was written. Re-run to preview against the current file.\n",
				deckFile, r.agent, r.field, got, r.deckVal)
			return 1
		}
	}
	updated := string(doc)
	for _, r := range append(append([]drop{}, redundant...), pins...) {
		updated = removeRosterField(updated, r.agent, r.field)
	}
	if err := writeRosterFileAtomic(deckFile, []byte(updated)); err != nil {
		fmt.Fprintf(stderr, "roster sync: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "\nWrote %s — %d redundant override(s) and %d deliberate pin(s) removed; the deck now inherits.\n",
		deckFile, len(redundant), len(pins))
	return 0
}

// removeRosterField deletes one key from one [roster.<id>] block, leaving every other
// byte — including comments — untouched.
func removeRosterField(doc, agent, field string) string {
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
	if start < 0 {
		return doc
	}
	out := make([]string, 0, len(lines))
	out = append(out, lines[:start+1]...)
	for _, ln := range lines[start+1 : end] {
		trimmed := strings.TrimSpace(ln)
		if key, _, ok := strings.Cut(trimmed, "="); ok && strings.TrimSpace(key) == field {
			continue
		}
		out = append(out, ln)
	}
	out = append(out, lines[end:]...)
	return strings.Join(out, "\n")
}

// unmatchedKeeps returns the --keep tokens that matched no deck override, sorted.
func unmatchedKeeps(keepSet, used map[string]bool) []string {
	var out []string
	for k := range keepSet {
		if k != "" && !used[k] {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

// rosterFieldValue reads one named field off a roster entry, for preview/apply binding.
func rosterFieldValue(e config.RosterEntry, field string) string {
	switch field {
	case "adapter":
		return e.Adapter
	case "model":
		return e.Model
	case "effort":
		return e.Effort
	case "speed":
		return e.Speed
	}
	return ""
}
