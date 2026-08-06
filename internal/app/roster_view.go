package app

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runmanifest"
)

// rosterViewOpts selects WHICH roster `roster show` answers about. Before these existed
// `--scope` parsed, was advertised in help, and changed nothing — the deck roster came
// back either way. A silently wrong answer is worse than a rejected flag, so scope is now
// load-bearing and every value it accepts produces a different, correct view.
type rosterViewOpts struct {
	// scope is "deck" (default) or "machine".
	scope string
	// all widens the view to every configured adapter, including ones no roster
	// declares. This is the answer to "I added opencode and it is invisible":
	// the roster legitimately excludes it, and --all is where you see it anyway.
	all bool
}

// rosterScopeFor resolves the membership+values view for a scope.
func rosterScopeFor(root, scope string) (config.RosterScope, error) {
	switch scope {
	case "", "deck":
		return config.LoadRosterScoped(root)
	case "machine":
		path := config.CentralAgentsPath()
		if path == "" {
			return config.RosterScope{}, fmt.Errorf("machine scope: cannot resolve the central config directory (set PARLEY_HOME)")
		}
		entries, err := config.RosterEntriesInFile(path)
		if err != nil {
			return config.RosterScope{}, err
		}
		out := config.RosterScope{Entries: entries, Members: map[string]bool{}, Source: path}
		for id := range entries {
			out.Members[id] = true
		}
		return out, nil
	default:
		return config.RosterScope{}, fmt.Errorf("invalid --scope %q (want deck|machine)", scope)
	}
}

// section2OnlyRows returns rows for IDs that exist ONLY in the legacy §2 table while a
// config roster is authoritative. The ratified field table decides their fate: "TOML
// wins; a §2-only ID is reported `unmapped`, never auto-added". Before this, they were
// dropped from `roster show` entirely and erased by `roster render` without a word.
func section2OnlyRows(root string, members map[string]bool) []rosterRow {
	if len(members) == 0 {
		return nil // legacy fallback already renders §2 itself
	}
	active, inactive, ok := protocol.ReadRosterIDs(root)
	if !ok {
		return nil
	}
	seen := map[string]bool{}
	var ids []string
	for id := range active {
		seen[id] = true
	}
	for id := range inactive {
		seen[id] = true
	}
	for id := range seen {
		if !members[id] {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	rows := make([]rosterRow, 0, len(ids))
	for _, id := range ids {
		r := rosterRow{Agent: id, State: "active", Model: agents.Unknown, Effort: agents.Unknown}
		if inactive[id] {
			r.State = "inactive"
			r.addStatus("inactive")
		}
		r.addStatus("unmapped")
		r.addStatus("section2-only")
		r.Note = fmt.Sprintf("declared only in the §2 table, which is no longer authoritative; add it with "+
			"`parley roster set %s --scope deck --adapter <family> --yes --confirm-breaking`, or let "+
			"`parley roster render` drop the row", id)
		rows = append(rows, r)
	}
	return rows
}

// unrosteredRows returns rows for configured adapters that no roster declares, used by
// --all. They are not members; they are shown so an operator can see that an agent they
// installed is deliberately not on this deck.
func unrosteredRows(root string, members map[string]bool, machineOnly bool) []rosterRow {
	specs, err := config.LoadAgentSpecsScoped(root, machineOnly)
	if err != nil {
		return nil
	}
	claimed := map[string]bool{}
	mapping, _ := config.LoadRosterAdaptersScoped(root, machineOnly)
	for id := range members {
		if fam := mapping[id]; fam != "" {
			claimed[fam] = true
		}
	}
	var rows []rosterRow
	for _, sp := range specs {
		if claimed[sp.ID] {
			continue
		}
		r := rosterRow{Agent: "–", Adapter: sp.ID, State: "not-in-roster"}
		r.Model, _ = sp.EffectiveModel()
		if r.Model == "" {
			r.Model = agents.Unknown
		}
		r.Effort, _ = sp.EffectiveEffort()
		if r.Effort == "" {
			r.Effort = agents.Unknown
		}
		r.Speed = sp.Speed
		meta := agents.DeriveModelMeta(r.Model)
		r.ModelFamily, r.ModelCompany = meta.Family, meta.Company
		r.addStatus("not-in-roster")
		r.Note = fmt.Sprintf("adapter %q is configured on this machine but no roster ID maps to it; "+
			"add one with `parley roster set <id> --scope deck --adapter %s --yes --confirm-breaking`", sp.ID, sp.ID)
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Adapter < rows[j].Adapter })
	return rows
}

// rosterExplain prints per-field provenance for one agent: which layer supplied each
// value and what the launch will actually use. D3 parked provenance here rather than in
// a SOURCE column, because a single column cannot describe eleven fields' origins.
func rosterExplain(root, agent string, opts rosterViewOpts, stdout, stderr io.Writer) int {
	scope, err := rosterScopeFor(root, opts.scope)
	if err != nil {
		fmt.Fprintf(stderr, "roster show: %v\n", err)
		return 1
	}
	if !scope.Members[agent] {
		fmt.Fprintf(stderr, "roster show: %q is not in this deck's roster (scope %s)\n", agent, orDeck(opts.scope))
		return 1
	}
	rows, err := resolveRoster(root, nil, opts)
	if err != nil {
		fmt.Fprintf(stderr, "roster show: %v\n", err)
		return 1
	}
	var row *rosterRow
	for i := range rows {
		if rows[i].Agent == agent {
			row = &rows[i]
			break
		}
	}
	if row == nil {
		fmt.Fprintf(stderr, "roster show: %q resolved to no row\n", agent)
		return 1
	}
	layers, err := config.RosterFieldSourcesScoped(root, agent, opts.scope == "machine")
	if err != nil {
		fmt.Fprintf(stderr, "roster show: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "%s — membership from %s%s\n\n", agent, scope.Source, inheritedSuffix(scope.Inherited))
	fmt.Fprintf(stdout, "%-14s %-24s %s\n", "FIELD", "EFFECTIVE", "SET BY")
	// A field no [roster.<id>] block sets can still come from a central or deck
	// [agents.<family>] block. Reporting "built-in default" there misattributes the value
	// that actually reaches the launch.
	specSources := config.AgentFieldSources(root, row.Adapter, opts.scope == "machine")
	// `active` is not layered — it follows the membership authority — so its provenance
	// must name that authority rather than whichever layer last wrote the key.
	if src, serr := config.RosterStateSource(root); serr == nil && src != "" {
		layers["active"] = src
	}
	show := func(field, effective string) {
		src := layers[field]
		if src == "" {
			src = specSources[field]
		}
		if src == "" {
			src = "built-in default"
		}
		fmt.Fprintf(stdout, "%-14s %-24s %s\n", field, orDash(effective), src)
	}
	show("adapter", row.Adapter)
	show("model", row.Model)
	show("effort", row.Effort)
	show("speed", row.Speed)
	show("active", row.State)
	if len(row.Status) > 0 {
		fmt.Fprintf(stdout, "\nstatus: %s\n", strings.Join(row.Status, ","))
	}
	if row.Note != "" {
		fmt.Fprintf(stdout, "note: %s\n", row.Note)
	}
	return 0
}

func orDeck(s string) string {
	if strings.TrimSpace(s) == "" {
		return "deck"
	}
	return s
}

func inheritedSuffix(inherited bool) string {
	if inherited {
		return " (INHERITED — this deck declares no roster of its own)"
	}
	return ""
}

// membersOf returns just the membership set for a scope; errors degrade to empty, which
// makes the callers treat the deck as legacy rather than crash a read-only command.
func membersOf(root, scope string) map[string]bool {
	sc, err := rosterScopeFor(root, scope)
	if err != nil {
		return map[string]bool{}
	}
	return sc.Members
}

// Roster snapshot states reported by `sessions inspect`.
const (
	rosterSnapshotCurrent = "current"
	rosterSnapshotStale   = "stale-snapshot"
	rosterSnapshotAbsent  = "no-snapshot"
	rosterSnapshotUnknown = "unknown"
)

// rosterSnapshotState compares a run's frozen roster revision with the deck's roster as
// it stands now.
//
// D6 ratified this as the audit half of the immutable-snapshot story: freezing the roster
// protects a running idea, but an operator still needs to know that the deck has since
// moved, or the protection is invisible. `stale-snapshot` was in the frozen STATUS
// vocabulary from the start with nothing able to emit it.
func rosterSnapshotState(root string, m runmanifest.Manifest) string {
	if len(m.RosterSnapshot) == 0 || strings.TrimSpace(m.RosterRevision) == "" {
		return rosterSnapshotAbsent
	}
	current, _, err := RosterSnapshot(root)
	if err != nil {
		return rosterSnapshotUnknown
	}
	if runmanifest.RosterRevisionOf(current) == m.RosterRevision {
		return rosterSnapshotCurrent
	}
	return rosterSnapshotStale
}
