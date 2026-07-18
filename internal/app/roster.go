package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
)

// Inject the roster-ID -> family mapping loader into the runner so every run path
// resolves roster-ID participants (claude-1, …) without threading it through each
// Options literal (idea composite-agent-naming-and-roster-reinit).
func init() {
	runner.RosterMappingLoader = func(root string) map[string]string {
		m, _ := config.LoadRosterAdapters(root)
		return m
	}
}

// familyAliases maps a roster-ID stem that is NOT a family prefix to its family.
// (antigravity-1 -> agy breaks a plain prefix match.) Used only for the init-time
// PROPOSAL, never for silent runtime resolution.
var familyAliases = map[string]string{"antigravity": "agy"}

var rosterInstanceSuffix = regexp.MustCompile(`-[0-9]+$`)

// runRoster implements `parley roster show|init` (component B): show renders the
// resolved roster with composite display names; init writes the roster-ID -> family
// `[roster.*]` mapping so the resolver can run a deck whose §2 roster is claude-1, …
func runRoster(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: parley roster <show|init> [--scope session|machine] [--dir DIR] [--dry-run] [--yes] [--json]")
		return 2
	}
	sub := args[0]
	rest := args[1:]
	fs := flag.NewFlagSet("roster", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace directory")
	scope := fs.String("scope", "session", "session (deck) or machine (~/.parley)")
	dryRun := fs.Bool("dry-run", false, "print what would change; write nothing")
	yes := fs.Bool("yes", false, "write without confirmation")
	jsonOut := fs.Bool("json", false, "machine-readable output")
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "roster: %v\n", err)
		return 1
	}
	switch sub {
	case "show":
		return rosterShow(root, *jsonOut, stdout, stderr)
	case "init":
		return rosterInit(root, *scope, *dryRun, *yes, *jsonOut, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "roster: unknown subcommand %q (want show|init)\n", sub)
		return 2
	}
}

type rosterRow struct {
	RosterID string `json:"roster_id"`
	Family   string `json:"family"`
	Display  string `json:"display_name"`
	Model    string `json:"model"`
	Effort   string `json:"effort"`
	Speed    string `json:"speed"`
	Auto     bool   `json:"autonomous"`
	Note     string `json:"note,omitempty"`
}

// resolveRoster builds one row per active §2 roster ID: family from the explicit
// mapping (else an init-time proposal), the matching spec, and the derived display
// name. A row whose family cannot be resolved is flagged (unmapped).
func resolveRoster(root string) ([]rosterRow, error) {
	specs, err := config.LoadAgentSpecs(root)
	if err != nil {
		return nil, err
	}
	byFamily := map[string]agents.Spec{}
	families := make([]string, 0, len(specs))
	for _, s := range specs {
		byFamily[s.ID] = s
		families = append(families, s.ID)
	}
	mapping, _ := config.LoadRosterAdapters(root)
	active, _, ok := protocol.ReadRosterIDs(root)
	if !ok {
		return nil, fmt.Errorf("could not read the §2 roster (COOPERATION.md)")
	}
	ids := make([]string, 0, len(active))
	for id := range active {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	rows := make([]rosterRow, 0, len(ids))
	for _, id := range ids {
		row := rosterRow{RosterID: id}
		family, mapped := mapping[id]
		if !mapped {
			family = proposeFamily(id, byFamily)
			if family != "" {
				row.Note = "unmapped — run `parley roster init`"
			}
		}
		row.Family = family
		if family == "" {
			row.Note = "unresolved — add `[roster." + id + "] adapter = \"<family>\"`"
			rows = append(rows, row)
			continue
		}
		spec := byFamily[family]
		row.Model = spec.Model
		row.Effort = spec.Reasoning
		row.Speed = spec.Speed
		row.Auto = spec.AutonomousWrite.Declared()
		if name, derr := agents.RenderDisplayName(family, spec); derr == nil {
			row.Display = name
		} else {
			row.Display = id
			if row.Note == "" {
				row.Note = "display fell back to the roster id: " + derr.Error()
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// proposeFamily proposes a family for a roster ID at init time (never at runtime):
// exact family match, then the id with a trailing -N stripped, then an alias.
func proposeFamily(id string, byFamily map[string]agents.Spec) string {
	if _, ok := byFamily[id]; ok {
		return id
	}
	stem := rosterInstanceSuffix.ReplaceAllString(id, "")
	if _, ok := byFamily[stem]; ok {
		return stem
	}
	if fam, ok := familyAliases[stem]; ok {
		if _, ok := byFamily[fam]; ok {
			return fam
		}
	}
	return ""
}

func rosterShow(root string, jsonOut bool, stdout, stderr io.Writer) int {
	rows, err := resolveRoster(root)
	if err != nil {
		fmt.Fprintf(stderr, "roster show: %v\n", err)
		return 1
	}
	if jsonOut {
		return writeJSON(stdout, stderr, map[string]any{"roster": rows})
	}
	fmt.Fprintf(stdout, "%-16s %-10s %-30s %-22s %-10s %-8s %s\n", "ROSTER-ID", "FAMILY", "DISPLAY-NAME", "MODEL", "EFFORT", "SPEED", "AUTO")
	for _, r := range rows {
		auto := "no"
		if r.Auto {
			auto = "yes"
		}
		fmt.Fprintf(stdout, "%-16s %-10s %-30s %-22s %-10s %-8s %s\n",
			r.RosterID, orDash(r.Family), orDash(r.Display), orDash(r.Model), orDash(r.Effort), orDash(r.Speed), auto)
		if r.Note != "" {
			fmt.Fprintf(stdout, "  ⚠ %s\n", r.Note)
		}
	}
	return 0
}

func rosterInit(root, scope string, dryRun, yes, jsonOut bool, stdout, stderr io.Writer) int {
	rows, err := resolveRoster(root)
	if err != nil {
		fmt.Fprintf(stderr, "roster init: %v\n", err)
		return 1
	}
	existing, _ := config.LoadRosterAdapters(root)
	// Propose the mapping entries that are missing; fail closed on any unresolved id.
	type entry struct{ id, family string }
	var toWrite []entry
	var unresolved []string
	for _, r := range rows {
		if r.Family == "" {
			unresolved = append(unresolved, r.RosterID)
			continue
		}
		if _, ok := existing[r.RosterID]; ok {
			continue // already mapped — idempotent
		}
		toWrite = append(toWrite, entry{r.RosterID, r.Family})
	}
	if len(unresolved) > 0 {
		fmt.Fprintf(stderr, "roster init: cannot resolve a family for: %s\n  add `[roster.<id>] adapter = \"<family>\"` for each, then re-run.\n", strings.Join(unresolved, ", "))
		return 1
	}
	target, targetLabel := rosterTargetPath(root, scope)

	if jsonOut {
		changes := make([]map[string]string, 0, len(toWrite))
		for _, e := range toWrite {
			changes = append(changes, map[string]string{"roster_id": e.id, "adapter": e.family})
		}
		return writeJSON(stdout, stderr, map[string]any{"scope": scope, "target": target, "dry_run": dryRun, "adds": changes})
	}

	if len(toWrite) == 0 {
		fmt.Fprintf(stdout, "Roster already initialized (%s): every §2 roster id already maps to a family.\n", targetLabel)
		return 0
	}
	fmt.Fprintf(stdout, "roster init (%s) will add to %s:\n", scope, targetLabel)
	for _, e := range toWrite {
		fmt.Fprintf(stdout, "  [roster.%s] adapter = %q\n", e.id, e.family)
	}
	if dryRun {
		fmt.Fprintln(stdout, "(dry-run: nothing written)")
		return 0
	}
	if !yes {
		fmt.Fprintln(stderr, "refusing to write without --yes (or use --dry-run to preview)")
		return 1
	}
	var b strings.Builder
	b.WriteString("\n# roster-ID -> family adapter map (parley roster init; composite-agent-naming-and-roster-reinit).\n")
	for _, e := range toWrite {
		b.WriteString(fmt.Sprintf("[roster.%s]\nadapter = %q\n", e.id, e.family))
	}
	if err := appendToFile(target, b.String()); err != nil {
		fmt.Fprintf(stderr, "roster init: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Wrote %d mapping(s) to %s. The driver can now run this roster.\n", len(toWrite), targetLabel)
	return 0
}

func rosterTargetPath(root, scope string) (path, label string) {
	if scope == "machine" {
		p := config.CentralAgentsPath()
		return p, "~/.parley/agents.toml"
	}
	return filepath.Join(root, protocol.DeckDir, "agents.toml"), "parley-deck/agents.toml"
}

func writeJSON(stdout, stderr io.Writer, v any) int {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "json: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, string(b))
	return 0
}

// appendToFile creates or appends to a config file (0644), creating the parent dir.
func appendToFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
