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
	"parley-deck-cli/internal/fsutil"
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

// rosterMappingFor loads the roster-ID -> family map for participant resolution at
// the app layer (run selection); errors degrade to nil (exact-ID matching only).
func rosterMappingFor(root string) map[string]string {
	m, _ := config.LoadRosterAdapters(root)
	return m
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
// resolveRoster builds one row per active §2 roster ID. allowedFamilies, when
// non-nil, restricts usable families to a scope catalog (machine-only for
// `--scope machine`) so a deck-only family is never proposed/blessed there.
func resolveRoster(root string, allowedFamilies map[string]bool) ([]rosterRow, error) {
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
		} else if _, ok := byFamily[family]; !ok {
			// Mapped to a family that is not configured/discovered: treat as
			// unresolved so init exits nonzero instead of "already initialized"
			// (review MAJOR — a typoed adapter would later fail the resolver).
			row.Note = fmt.Sprintf("mapped to unknown family %q — fix `[roster.%s] adapter`", family, id)
			family = ""
		}
		// Scope filter: a family outside the allowed catalog (machine-only for
		// --scope machine) is not usable in this scope (review MAJOR, codex-1).
		if family != "" && allowedFamilies != nil && !allowedFamilies[family] {
			row.Note = fmt.Sprintf("family %q is not known machine-wide — define [agents.%s] in ~/.parley/agents.toml", family, family)
			family = ""
		}
		row.Family = family
		if family == "" {
			if row.Note == "" {
				row.Note = "unresolved — add `[roster." + id + "] adapter = \"<family>\"`"
			}
			rows = append(rows, row)
			continue
		}
		spec := byFamily[family]
		row.Model = spec.Model
		row.Effort = spec.Reasoning
		row.Speed = spec.Speed
		// Fail-closed: a declared mode whose enabling args the launch does not pass is
		// not autonomous (review MAJOR, codex-1).
		row.Auto = spec.AutonomousEffective()
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
	rows, err := resolveRoster(root, nil)
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

type rosterMapEntry struct{ id, family string }

func rosterInit(root, scope string, dryRun, yes, jsonOut bool, stdout, stderr io.Writer) int {
	if scope != "session" && scope != "machine" {
		fmt.Fprintf(stderr, "roster init: invalid --scope %q (want session|machine)\n", scope)
		return 2
	}
	// Machine scope must never propose, write, or bless a deck-only family (consensus
	// §B); build a machine-only family catalog and restrict resolution to it. Session
	// scope uses the full layered catalog (allowed == nil).
	var allowed map[string]bool
	if scope == "machine" {
		c, cerr := config.MachineFamilyCatalog()
		if cerr != nil {
			fmt.Fprintf(stderr, "roster init: %v\n", cerr)
			return 1
		}
		allowed = c
	}
	rows, err := resolveRoster(root, allowed)
	if err != nil {
		fmt.Fprintf(stderr, "roster init: %v\n", err)
		return 1
	}
	target, targetLabel := rosterTargetPath(root, scope)
	// Idempotency is judged against the TARGET file, not the layered stack — an
	// inherited mapping must not suppress a write to the requested scope (review MAJOR).
	existing, err := config.RosterAdaptersInFile(target)
	if err != nil {
		fmt.Fprintf(stderr, "roster init: %v\n", err)
		return 1
	}
	// Validate the TARGET file's own existing mappings against the SCOPE catalog: a
	// broken adapter there (e.g. a typo the deck layer happens to override, or a
	// deck-only family in the machine file) must be a hard error, not a silent
	// "already initialized" while the requested scope stays broken (review MAJOR).
	known := allowed
	if known == nil {
		specs, _ := config.LoadAgentSpecs(root)
		known = make(map[string]bool, len(specs))
		for _, s := range specs {
			known[s.ID] = true
		}
	}
	for id, fam := range existing {
		if !known[fam] {
			fmt.Fprintf(stderr, "roster init: %s maps %s -> %q, which is not a known agent family — fix that [roster.%s] block\n", targetLabel, id, fam, id)
			return 1
		}
	}
	var toWrite []rosterMapEntry
	var unresolved []string
	for _, r := range rows {
		if r.Family == "" {
			unresolved = append(unresolved, r.RosterID)
			continue
		}
		if _, ok := existing[r.RosterID]; ok {
			continue // already in the target file — idempotent
		}
		toWrite = append(toWrite, rosterMapEntry{r.RosterID, r.Family})
	}
	if len(unresolved) > 0 {
		fmt.Fprintf(stderr, "roster init: cannot resolve a family for: %s\n  add `[roster.<id>] adapter = \"<family>\"` for each, then re-run.\n", strings.Join(unresolved, ", "))
		return 1
	}

	// Determine the outcome and perform any write BEFORE rendering, so --json
	// reports what actually happened rather than a proposal (review MAJOR).
	outcome := "unchanged"
	if len(toWrite) > 0 {
		switch {
		case dryRun:
			outcome = "dry-run"
		case !yes:
			outcome = "needs-confirmation"
		default:
			n, werr := writeRosterMappings(target, toWrite)
			if werr != nil {
				fmt.Fprintf(stderr, "roster init: %v\n", werr)
				return 1
			}
			if n == 0 {
				outcome = "unchanged"
			} else {
				outcome = "written"
			}
		}
	}

	if jsonOut {
		adds := make([]map[string]string, 0, len(toWrite))
		for _, e := range toWrite {
			adds = append(adds, map[string]string{"roster_id": e.id, "adapter": e.family})
		}
		code := 0
		if outcome == "needs-confirmation" {
			code = 1
		}
		_ = writeJSON(stdout, stderr, map[string]any{"scope": scope, "target": target, "outcome": outcome, "adds": adds})
		return code
	}

	switch outcome {
	case "unchanged":
		fmt.Fprintf(stdout, "Roster already initialized (%s): every §2 roster id already maps to a family.\n", targetLabel)
		return 0
	case "dry-run":
		fmt.Fprintf(stdout, "roster init (%s) would add to %s:\n", scope, targetLabel)
		for _, e := range toWrite {
			fmt.Fprintf(stdout, "  [roster.%s] adapter = %q\n", e.id, e.family)
		}
		fmt.Fprintln(stdout, "(dry-run: nothing written)")
		return 0
	case "needs-confirmation":
		fmt.Fprintf(stdout, "roster init (%s) will add %d mapping(s) to %s.\n", scope, len(toWrite), targetLabel)
		fmt.Fprintln(stderr, "refusing to write without --yes (or use --dry-run to preview)")
		return 1
	default: // written
		fmt.Fprintf(stdout, "Wrote %d mapping(s) to %s. The driver can now run this roster.\n", len(toWrite), targetLabel)
		return 0
	}
}

// writeRosterMappings appends the missing [roster.<id>] blocks to target via an
// atomic temp+rename write (fsutil.WriteFileAtomic), re-reading first and skipping
// any block already present textually so a concurrent/repeat run cannot duplicate a
// table or leave a partial file (review MAJOR). Returns the number written.
func writeRosterMappings(target string, toWrite []rosterMapEntry) (int, error) {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return 0, err
	}
	prior, err := os.ReadFile(target)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	content := string(prior)
	// Re-parse under the write to decide presence from the PARSED mapping, not a
	// substring: an empty/malformed `[roster.<id>]` block (adapter = "") must be a
	// hard error, never a silent "already initialized" (review MINOR, kimi-1). A
	// concurrently-written VALID block is skipped (idempotent).
	present, err := config.RosterAdaptersInFile(target)
	if err != nil {
		return 0, err
	}
	var b strings.Builder
	b.WriteString(content)
	if content != "" && !strings.HasSuffix(content, "\n") {
		b.WriteString("\n")
	}
	b.WriteString("\n# roster-ID -> family adapter map (parley roster init; composite-agent-naming-and-roster-reinit).\n")
	wrote := 0
	for _, e := range toWrite {
		if _, ok := present[e.id]; ok {
			continue // a valid mapping already exists (concurrent/repeat guard)
		}
		b.WriteString(fmt.Sprintf("[roster.%s]\nadapter = %q\n", e.id, e.family))
		wrote++
	}
	if wrote == 0 {
		return 0, nil
	}
	candidate := []byte(b.String())
	// Structurally validate the whole candidate BEFORE the atomic replace so an
	// existing empty/quoted-key `[roster.<id>]` block that our append would duplicate
	// is rejected (a duplicate TOML table would break every later load) — review MAJOR.
	if verr := config.ValidateAgentsConfigBytes(candidate); verr != nil {
		return 0, fmt.Errorf("refusing to write %s: candidate is invalid (likely a duplicate/malformed [roster.*] table already present): %w", target, verr)
	}
	return wrote, fsutil.WriteFileAtomic(target, candidate, 0o644)
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
