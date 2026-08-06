package app

import (
	"context"
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

// RosterSchemaVersion versions the canonical roster contract. The column list, their
// order, and the JSON field names are an API: adding a column is an additive, documented
// change that bumps this, and golden tests pin the shape. Before it existed, three CLI
// surfaces answered "what is the roster?" with three different tables.
const RosterSchemaVersion = 1

// RosterColumns is the frozen v1 column list, in render order.
var RosterColumns = []string{
	"AGENT", "ADAPTER", "STATE", "INSTALLED", "MODEL",
	"MODEL-FAMILY", "MODEL-COMPANY", "EFFORT", "SPEED", "AUTO", "STATUS",
}

// rosterRow is one canonical roster row.
//
// MODEL and EFFORT hold the value the launch ACTUALLY passes, or "unknown" — never a
// declaration wearing the effective cell. That distinction is the whole point: the old
// table printed spec.Model while the process launched a different model entirely.
// Divergence and absence surface through Status codes rather than by quietly
// substituting the declared value.
type rosterRow struct {
	Agent        string   `json:"agent"`
	Adapter      string   `json:"adapter"`
	State        string   `json:"state"`
	Installed    bool     `json:"installed"`
	Model        string   `json:"model"`
	ModelFamily  string   `json:"model_family"`
	ModelCompany string   `json:"model_company"`
	Effort       string   `json:"effort"`
	Speed        string   `json:"speed"`
	Auto         bool     `json:"autonomous"`
	Status       []string `json:"status"`
	// Display is the derived composite name. It left the canonical table because it
	// duplicates ADAPTER+MODEL+EFFORT and could contradict MODEL; it stays here for
	// TUI and artifact rendering, derived from this row.
	Display string `json:"display_name,omitempty"`
	Note    string `json:"note,omitempty"`
}

// addStatus appends a status code, keeping the list free of duplicates.
func (r *rosterRow) addStatus(code string) {
	for _, existing := range r.Status {
		if existing == code {
			return
		}
	}
	r.Status = append(r.Status, code)
}

// statusOrOK renders the status cell; an empty list means nothing is wrong.
func (r rosterRow) statusOrOK() string {
	if len(r.Status) == 0 {
		return "ok"
	}
	return strings.Join(r.Status, ",")
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
	// INSTALLED comes from a real discovery probe: a rostered agent whose binary is not
	// on this machine is operationally absent, and the roster must say so.
	installed := map[string]bool{}
	for _, d := range agents.Discover(context.Background(), specs) {
		installed[d.ID] = d.Found
	}
	mapping, _ := config.LoadRosterAdapters(root)

	// MEMBERSHIP AUTHORITY. parley-deck/agents.toml owns the roster; §2 is a generated,
	// non-authoritative view. A deck that predates the cutover has no [roster.*] block at
	// all, so it falls back to the legacy §2 table and every row is flagged
	// `legacy-roster` — the drift this idea exists to end came from 40 decks maintaining
	// that table by hand, nine different ways.
	entries, entriesErr := config.LoadRoster(root)
	if entriesErr != nil {
		return nil, entriesErr
	}
	legacy := len(entries) == 0
	inactive := map[string]bool{}
	var ids []string
	if legacy {
		// The parser puts EVERY row in `active`, including inactive ones, and reports
		// `inactive` separately; STATE is what distinguishes them. That second map used
		// to be discarded, so a retired agent rendered as a full member.
		active, inactiveIDs, ok := protocol.ReadRosterIDs(root)
		if !ok {
			return nil, fmt.Errorf("no roster: declare [roster.<id>] in parley-deck/agents.toml (or keep a legacy §2 table in COOPERATION.md)")
		}
		inactive = inactiveIDs
		for id := range active {
			ids = append(ids, id)
		}
	} else {
		for id, e := range entries {
			ids = append(ids, id)
			if !e.Active {
				inactive[id] = true
			}
		}
	}
	sort.Strings(ids)
	rows := make([]rosterRow, 0, len(ids))
	for _, id := range ids {
		row := rosterRow{Agent: id, State: "active"}
		if legacy {
			row.addStatus("legacy-roster")
		}
		if inactive[id] {
			row.State = "inactive"
			row.addStatus("inactive")
		}
		family, mapped := mapping[id]
		if !mapped {
			family = proposeFamily(id, byFamily)
			if family != "" {
				row.Note = "unmapped — run `parley roster init`"
				row.addStatus("unmapped")
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
		row.Adapter = family
		if family == "" {
			if row.Note == "" {
				row.Note = "unresolved — add `[roster." + id + "] adapter = \"<family>\"`"
			}
			row.addStatus("unmapped")
			row.Model, row.Effort = agents.Unknown, agents.Unknown
			row.ModelFamily, row.ModelCompany = agents.Unknown, agents.Unknown
			rows = append(rows, row)
			continue
		}
		spec := byFamily[family]
		// Per-roster-ID settings beat the adapter-family default, so two roster IDs can
		// share an adapter and still run different models.
		if e, ok := entries[id]; ok {
			if e.Model != "" {
				spec.Model = e.Model
			}
			if e.Effort != "" {
				spec.Reasoning = e.Effort
			}
			if e.Speed != "" {
				spec.Speed = e.Speed
			}
		}
		row.Installed = installed[family]
		if !row.Installed {
			row.addStatus("not-installed")
		}
		// MODEL and EFFORT are the values the launch passes, or unknown. A configured
		// value the argv never carries is NOT effective and must not fill the cell.
		if model, ok := spec.EffectiveModel(); ok {
			row.Model = model
			if declared := strings.TrimSpace(spec.Model); declared != "" && declared != agents.CLIDefault && declared != model {
				row.addStatus("model-drift")
			}
		} else {
			row.Model = agents.Unknown
			row.addStatus("model-unbound")
		}
		if effort, ok := spec.EffectiveEffort(); ok {
			row.Effort = effort
		} else {
			row.Effort = agents.Unknown
			row.addStatus("effort-unknown")
		}
		meta := agents.DeriveModelMeta(row.Model)
		row.ModelFamily, row.ModelCompany = meta.Family, meta.Company
		if !meta.Known {
			row.addStatus("metadata-unknown")
		}
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
		// schema_version and the ordered column list travel WITH the rows: a consumer
		// must be able to detect a contract change rather than infer it from the keys.
		return writeJSON(stdout, stderr, map[string]any{
			"schema_version": RosterSchemaVersion,
			"columns":        RosterColumns,
			"roster":         rows,
		})
	}
	const format = "%-12s %-10s %-8s %-9s %-22s %-14s %-13s %-8s %-8s %-4s %s\n"
	fmt.Fprintf(stdout, format, "AGENT", "ADAPTER", "STATE", "INSTALLED", "MODEL",
		"MODEL-FAMILY", "MODEL-COMPANY", "EFFORT", "SPEED", "AUTO", "STATUS")
	for _, r := range rows {
		fmt.Fprintf(stdout, format,
			r.Agent, orDash(r.Adapter), r.State, yesNo(r.Installed), orDash(r.Model),
			orDash(r.ModelFamily), orDash(r.ModelCompany), orDash(r.Effort),
			orDash(r.Speed), yesNo(r.Auto), r.statusOrOK())
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
		if r.Adapter == "" {
			unresolved = append(unresolved, r.Agent)
			continue
		}
		if _, ok := existing[r.Agent]; ok {
			continue // already in the target file — idempotent
		}
		toWrite = append(toWrite, rosterMapEntry{r.Agent, r.Adapter})
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
