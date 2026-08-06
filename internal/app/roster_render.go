package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"parley-deck-cli/internal/config"
)

// The §2 roster table is a GENERATED VIEW of parley-deck/agents.toml, not a store.
//
// It used to be the store, maintained by hand in every project, and it drifted into nine
// different rosters across 40 decks — 17 with none at all, 17 still naming an agent
// retired months earlier. Generation is what stops that recurring: the table is rendered
// from config, and hand-edits to it are overwritten rather than obeyed.
//
// Generation MUST be idempotent — rendering twice produces byte-identical output —
// otherwise every render is a diff and the drift returns wearing a new name.

const rosterSectionHeader = "## 2. Active agents (roster)"

// renderRosterTable renders the §2 table rows for a deck. Ordering is fixed: active
// before inactive, then agent ID byte-ascending, so the output cannot depend on map
// iteration order.
func renderRosterTable(root string, carry map[string]legacyCells, adoptInherited bool) (string, []string, error) {
	scope, err := config.LoadRosterScoped(root)
	if err != nil {
		return "", nil, err
	}
	if len(scope.Members) == 0 {
		return "", nil, fmt.Errorf("no [roster.*] entries in this deck's config — nothing to render")
	}
	// A deck that declares no roster of its own shows the machine roster, but committing
	// it into COOPERATION.md would bake a machine-local accident into a shared file and
	// stale it on the next machine change. Refuse unless the operator says otherwise.
	if scope.Inherited && !adoptInherited {
		return "", nil, fmt.Errorf(
			"this deck declares no roster of its own; the %d rows shown come from %s.\n"+
				"Writing them into §2 would commit a machine-local roster into a shared file.\n"+
				"Declare the roster with `parley roster set <agent> --scope deck --adapter <family>`,\n"+
				"or re-run with --adopt-inherited to accept the inherited roster as this deck's own",
			len(scope.Members), scope.Source)
	}
	entries := scope.Entries
	ids := make([]string, 0, len(scope.Members))
	for id := range scope.Members {
		ids = append(ids, id)
	}
	// Rows present in the existing §2 but absent from the authority are REMOVED by this
	// render. The ratified field table says a §2-only ID is "reported `unmapped`, never
	// auto-added" — so it does not re-enter the generated table, but its removal must be
	// stated. Erasing project data silently is how the previous generator lost rows.
	removed := make([]string, 0)
	for id := range carry {
		if !scope.Members[id] {
			removed = append(removed, id)
		}
	}
	sort.Strings(removed)
	sort.Slice(ids, func(i, j int) bool {
		a, b := entries[ids[i]], entries[ids[j]]
		if a.Active != b.Active {
			return a.Active // active rows first
		}
		return ids[i] < ids[j]
	})

	var b strings.Builder
	b.WriteString("| Agent ID | Workspace dir | Role | State |\n")
	b.WriteString("| -------- | ------------- | ---- | ----- |\n")
	for _, id := range ids {
		e := entries[id]
		state := "active"
		if !e.Active {
			state = "inactive"
		}
		// Config wins, but a value only the old hand-written table carried is PRESERVED
		// rather than replaced with a placeholder. Generation must never lose project
		// data that no one has migrated yet.
		role, dir := e.Role, e.WorkspaceDir
		if c, ok := carry[id]; ok {
			if role == "" {
				role = c.role
			}
			if dir == "" {
				dir = c.dir
			}
		}
		if role == "" {
			role = "participant"
		}
		if dir == "" {
			dir = "–"
		}
		b.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n", id, dir, role, state))
	}
	return b.String(), removed, nil
}

// rosterRender writes the generated table into COOPERATION.md between the §2 heading and
// the next section, leaving every other byte untouched.
func rosterRender(root string, dryRun, yes, adoptInherited bool, stdout, stderr io.Writer) int {
	path0 := filepath.Join(root, "parley-deck", "COOPERATION.md")
	prior, _ := os.ReadFile(path0)
	table, removed, err := renderRosterTable(root, parseLegacyCells(string(prior)), adoptInherited)
	if err != nil {
		fmt.Fprintf(stderr, "roster render: %v\n", err)
		return 1
	}
	// Report removals before anything else, in preview AND on apply: they are the one
	// consequence of regeneration that destroys existing content.
	reportRemoved := func(w io.Writer) {
		if len(removed) == 0 {
			return
		}
		fmt.Fprintf(w, "the following §2 row(s) are NOT in this deck's roster and will be removed:\n")
		for _, id := range removed {
			fmt.Fprintf(w, "  - %s (reported `unmapped` by `parley roster show`; add it with\n"+
				"    `parley roster set %s --scope deck --adapter <family> --yes --confirm-breaking` to keep it)\n", id, id)
		}
	}
	path := path0
	doc, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "roster render: %v\n", err)
		return 1
	}
	updated, changed, err := replaceRosterSection(string(doc), table)
	if err != nil {
		fmt.Fprintf(stderr, "roster render: %v\n", err)
		return 1
	}
	if !changed {
		fmt.Fprintf(stdout, "roster render: §2 already matches the roster — nothing to do\n")
		return 0
	}
	if dryRun || !yes {
		reportRemoved(stdout)
		fmt.Fprintf(stdout, "would regenerate §2 in %s:\n\n%s\nNothing was written. Re-run with --yes to apply.\n", path, table)
		return 0
	}
	reportRemoved(stdout)
	if err := writeRosterFileAtomic(path, []byte(updated)); err != nil {
		fmt.Fprintf(stderr, "roster render: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Regenerated §2 in %s\n", path)
	return 0
}

// replaceRosterSection swaps the first markdown table after the §2 heading for the
// generated one. It deliberately replaces only the table: the surrounding prose explains
// WHY the section is generated and must survive.
func replaceRosterSection(doc, table string) (string, bool, error) {
	idx := strings.Index(doc, rosterSectionHeader)
	if idx < 0 {
		return "", false, fmt.Errorf("no %q heading found", rosterSectionHeader)
	}
	lines := strings.Split(doc, "\n")
	headerLine := -1
	for i, ln := range lines {
		if strings.TrimSpace(ln) == rosterSectionHeader {
			headerLine = i
			break
		}
	}
	if headerLine < 0 {
		return "", false, fmt.Errorf("no %q heading found", rosterSectionHeader)
	}
	start, end := -1, -1
	for i := headerLine + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "## ") {
			break // reached the next section without finding a table
		}
		if strings.HasPrefix(trimmed, "|") {
			if start < 0 {
				start = i
			}
			end = i
			continue
		}
		if start >= 0 && trimmed != "" {
			break // the table ended
		}
	}
	if start < 0 {
		return "", false, fmt.Errorf("§2 has no table to regenerate")
	}
	existing := strings.Join(lines[start:end+1], "\n") + "\n"
	if existing == table {
		return doc, false, nil
	}
	out := append([]string{}, lines[:start]...)
	out = append(out, strings.Split(strings.TrimRight(table, "\n"), "\n")...)
	out = append(out, lines[end+1:]...)
	return strings.Join(out, "\n"), true, nil
}

// legacyCells holds the render-only values the old hand-written §2 table carried.
type legacyCells struct{ dir, role string }

// parseLegacyCells reads the existing §2 table so a render can preserve workspace dir and
// role for agents whose config does not carry them yet.
func parseLegacyCells(doc string) map[string]legacyCells {
	out := map[string]legacyCells{}
	lines := strings.Split(doc, "\n")
	in := false
	for _, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if trimmed == rosterSectionHeader {
			in = true
			continue
		}
		if in && strings.HasPrefix(trimmed, "## ") {
			break
		}
		if !in || !strings.HasPrefix(trimmed, "|") {
			continue
		}
		cols := strings.Split(strings.Trim(trimmed, "|"), "|")
		if len(cols) < 3 {
			continue
		}
		id := strings.Trim(strings.TrimSpace(cols[0]), "`")
		if id == "" || id == "Agent ID" || strings.HasPrefix(id, "-") {
			continue
		}
		out[id] = legacyCells{dir: strings.TrimSpace(cols[1]), role: strings.TrimSpace(cols[2])}
	}
	return out
}
