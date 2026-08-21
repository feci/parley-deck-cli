package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

const EnvAgentConfig = "PARLEY_HEADLESS_AGENT_CONFIG"

// EnvParleyHome overrides the user-global config dir (default ~/.parley). It
// lets tests point the central-default layer at a scratch directory.
const EnvParleyHome = "PARLEY_HOME"

type fileConfig struct {
	Defaults *globalDefaults           `toml:"defaults"`
	Agents   map[string]agentOverride  `toml:"agents"`
	Rosters  map[string]rosterOverride `toml:"rosters"`
	Roster   map[string]rosterAdapter  `toml:"roster"`
}

// rosterAdapter is one [roster.<roster-id>] block: the family/adapter a roster ID
// resolves to (composite-agent-naming-and-roster-reinit). NOTE the singular table
// name `[roster.*]` (the ID->family map) is distinct from the plural `[rosters.*]`
// participant presets above.
// rosterAdapter is one [roster.<id>] block. It began as an adapter mapping only; the
// remaining fields make parley-deck/agents.toml the deck's roster AUTHORITY, so
// membership and per-agent settings stop living in hand-edited §2 prose. Only ID,
// Adapter, Active, Model, Effort and Speed are runtime-semantic; WorkspaceDir, Role and
// HostHandle are render-only and are carried verbatim from the legacy §2 table.
//
// Absent Active means true: existing blocks that carry only `adapter` keep working.
type rosterAdapter struct {
	Adapter      string `toml:"adapter"`
	Active       *bool  `toml:"active"`
	Model        string `toml:"model"`
	Effort       string `toml:"effort"`
	Speed        string `toml:"speed"`
	WorkspaceDir string `toml:"workspace_dir"`
	Role         string `toml:"role"`
	HostHandle   string `toml:"host_handle"`
}

// RosterEntry is the resolved, layered view of one [roster.<id>] block.
type RosterEntry struct {
	ID           string
	Adapter      string
	Active       bool
	Model        string
	Effort       string
	Speed        string
	WorkspaceDir string
	Role         string
	HostHandle   string
}

// LoadRoster returns the layered roster entries keyed by roster ID, lowest layer first
// so a deck overrides the machine default field by field. An empty result means no
// config layer declares a roster at all — the caller then falls back to the legacy §2
// table and reports `legacy-roster`.
func LoadRoster(root string) (map[string]RosterEntry, error) {
	out := map[string]RosterEntry{}
	for _, item := range configLayers(root) {
		data, err := os.ReadFile(item.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		var cfg fileConfig
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return out, fmt.Errorf("%s: %w", item.path, err)
		}
		for id, ra := range cfg.Roster {
			out[id] = mergeRosterEntry(out[id], id, ra)
		}
	}
	return out, nil
}

// RosterScope is the layered roster split into the two things that were previously
// conflated: WHO is on the deck, and WHAT VALUES each member launches with.
//
// LoadRoster unions membership across every layer, which meant a deck declaring two
// members resolved to five whenever the machine file listed five — and, because
// `roster render` writes §2 from the same view, that inherited membership got committed
// into COOPERATION.md, where it went stale on the next machine change. That is the exact
// drift this change exists to end, re-created inside committed files.
//
// The rule: the DECK FILE owns membership. The machine layer seeds values only — the same
// rebase model `roster sync` already implements. A deck that declares no roster of its own
// (neither a [roster.*] block nor a valid legacy §2 table) may still display the machine
// roster, but every such row is marked Inherited so no surface can mistake it for a
// deck-level decision.
type RosterScope struct {
	// Entries holds the fully layered values for every member, keyed by roster ID.
	Entries map[string]RosterEntry
	// Members is the authoritative membership set.
	Members map[string]bool
	// Inherited is true when Members came from the machine layer because no deck
	// layer declared a roster. Callers that write committed files MUST refuse to
	// persist an inherited roster without an explicit operator flag.
	Inherited bool
	// Source names the layer that decided membership, for display and diagnostics.
	Source string
	// Legacy is true when membership came from the pre-cutover §2 table. Values still
	// layer normally; the rows are reported `legacy-roster` until the deck is migrated.
	Legacy bool
}

// LoadRosterScoped returns the roster with membership and values separated. See
// RosterScope for the rule and the reason it exists.
func LoadRosterScoped(root string) (RosterScope, error) {
	out := RosterScope{Entries: map[string]RosterEntry{}, Members: map[string]bool{}}
	deckMembers := map[string]bool{}
	machineMembers := map[string]bool{}
	deckSource, machineSource := "", ""
	// ACTIVE STATE FOLLOWS THE AUTHORITY, not the layer stack. Membership IDs were gated
	// to the committed deck file, but `active` still merged from every layer — so
	// `[roster.claude-1] active = false` in the gitignored agents.local.toml, the env
	// config, or the machine file could quietly drop a committed member from the quorum
	// (or revive one the deck retired) using a file collaborators never see. Retiring a
	// member is a membership change; it belongs to the same record that grants membership.
	deckActive := map[string]bool{}
	machineActive := map[string]bool{}
	for _, item := range configLayers(root) {
		ids, entries, err := rosterLayer(item.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		// Values merge across every layer, lowest first, exactly as before.
		for id, ra := range entries {
			out.Entries[id] = mergeRosterEntry(out.Entries[id], id, ra)
		}
		if len(ids) == 0 {
			continue
		}
		if item.machine || !item.membership {
			// Only the machine layer is tracked as a membership candidate; the gitignored
			// agents.local.toml and $PARLEY_HEADLESS_AGENT_CONFIG are per-machine value
			// layers. Letting them add members would put a deck's quorum in a file its
			// collaborators never see.
			if item.machine {
				for _, id := range ids {
					machineMembers[id] = true
					machineActive[id] = entries[id].Active == nil || *entries[id].Active
				}
				if machineSource == "" {
					machineSource = item.source
				}
			}
			continue
		}
		for _, id := range ids {
			deckMembers[id] = true
			deckActive[id] = entries[id].Active == nil || *entries[id].Active
		}
		if deckSource == "" {
			deckSource = item.source
		}
	}
	// AUTHORITY ORDER, decided before any value layering:
	//   1. committed deck blocks (parley-deck/agents.toml)
	//   2. else a VALID legacy §2 table — a deck that predates the cutover keeps its own
	//      membership until it is migrated; the machine roster must not be inherited over
	//      a roster the deck actually declares, merely because it declares it in prose
	//   3. else the machine roster, explicitly marked Inherited
	// Step 2 was ratified and then omitted: LoadRosterScoped knew only about TOML, so any
	// machine roster silently outranked a legacy deck's four declared members.
	if len(deckMembers) > 0 {
		out.Members, out.Source = deckMembers, deckSource
		out.applyAuthorityState(deckActive)
		return out, nil
	}
	if active, inactive, ok := protocol.ReadRosterIDs(root); ok {
		out.Members = map[string]bool{}
		legacyActive := map[string]bool{}
		for id := range active {
			out.Members[id] = true
			legacyActive[id] = true
		}
		for id := range inactive {
			out.Members[id] = true
			legacyActive[id] = false
		}
		out.Source, out.Legacy = "COOPERATION.md §2", true
		out.applyAuthorityState(legacyActive)
		return out, nil
	}
	if len(machineMembers) > 0 {
		out.Members, out.Source, out.Inherited = machineMembers, machineSource, true
		out.applyAuthorityState(machineActive)
	}
	return out, nil
}

// applyAuthorityState forces each member's Active to the value the AUTHORITY layer
// declared, discarding any `active` a value-only layer merged in.
func (r RosterScope) applyAuthorityState(active map[string]bool) {
	for id := range r.Members {
		e := r.Entries[id]
		if e.ID == "" {
			e = RosterEntry{ID: id}
		}
		state, declared := active[id]
		e.Active = !declared || state
		r.Entries[id] = e
	}
}

// rosterLayer reads one config file and returns the roster IDs it declares plus their
// raw blocks. A missing file yields os.ErrNotExist for the caller to skip.
func rosterLayer(path string) ([]string, map[string]rosterAdapter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var cfg fileConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	ids := make([]string, 0, len(cfg.Roster))
	for id := range cfg.Roster {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, cfg.Roster, nil
}

// mergeRosterEntry applies one layer's block over the accumulated entry. Extracted so
// LoadRoster and LoadRosterScoped cannot drift apart in how a field wins.
func mergeRosterEntry(e RosterEntry, id string, ra rosterAdapter) RosterEntry {
	if e.ID == "" {
		e = RosterEntry{ID: id, Active: true}
	}
	if v := strings.TrimSpace(ra.Adapter); v != "" {
		e.Adapter = v
	}
	if ra.Active != nil {
		e.Active = *ra.Active
	}
	if v := strings.TrimSpace(ra.Model); v != "" {
		e.Model = v
	}
	if v := strings.TrimSpace(ra.Effort); v != "" {
		e.Effort = v
	}
	if v := strings.TrimSpace(ra.Speed); v != "" {
		e.Speed = v
	}
	if v := strings.TrimSpace(ra.WorkspaceDir); v != "" {
		e.WorkspaceDir = v
	}
	if v := strings.TrimSpace(ra.Role); v != "" {
		e.Role = v
	}
	if v := strings.TrimSpace(ra.HostHandle); v != "" {
		e.HostHandle = v
	}
	return e
}

// rosterOverride is one [rosters.<slug>] block — a named participant preset
// (named-roster-presets). Members are §2 canonical roster IDs (e.g. "claude-1").
type rosterOverride struct {
	Participants []string `toml:"participants"`
}

// globalDefaults is the optional [defaults] block of a layered config file —
// non-agent policy knobs that apply project-wide.
type globalDefaults struct {
	Speed              string            `toml:"speed"`
	PingTier           string            `toml:"ping_tier"`
	PreferredTransport string            `toml:"preferred_transport"`
	RosterChangePolicy string            `toml:"roster_change_policy"`
	Timeouts           *timeoutsBlock    `toml:"timeouts"`
	Loop               *loopBlock        `toml:"loop"`
	TrackRosters       map[string]string `toml:"track_rosters"`
}

type timeoutsBlock struct {
	SignoffMS       int `toml:"signoff_ms"`
	RoundMS         int `toml:"round_ms"`
	ReviewMS        int `toml:"review_ms"`
	DeepReasoningMS int `toml:"deep_reasoning_ms"`
}

// loopBlock is the [defaults.loop] policy: explicit auto-drive loop ceilings (LE-5).
// A breach escalates and halts; it never marks an idea complete. 0 = unlimited.
// Fields are pointers so a deliberate `= 0` at a higher layer overrides a lower
// layer's seeded value (presence-aware merge); absence (nil) falls through (review
// fix F-T2-1).
type loopBlock struct {
	MaxDriverSteps *int     `toml:"max_driver_steps"`
	MaxWallClockMS *int     `toml:"max_wall_clock_ms"`
	MaxCostUSD     *float64 `toml:"max_cost_usd"`
}

// CentralDefaults are the merged [defaults] policy knobs across the layered
// config files (central ~/.parley/agents.toml first, then the project deck).
// A zero-value field means "not set" — the consumer supplies its own fallback.
type CentralDefaults struct {
	Speed              string
	PingTier           string
	PreferredTransport string
	RosterChangePolicy string
	SignoffMS          int
	RoundMS            int
	ReviewMS           int
	DeepReasoningMS    int
	// Loop ceilings (LE-5); 0 = unlimited.
	MaxDriverSteps int
	MaxWallClockMS int
	MaxCostUSD     float64
}

type configLayer struct {
	path     string
	source   string
	optional bool
	// membership marks the ONE layer that may declare who is on the deck: the committed
	// parley-deck/agents.toml. Every other layer supplies values only.
	membership bool
	// machine marks the user-global ~/.parley layer. Deck layers own MEMBERSHIP;
	// the machine layer only seeds VALUES for members the deck already declares.
	// See LoadRosterScoped for why the distinction exists.
	machine bool
}

type agentOverride struct {
	Command               string            `toml:"command"`
	Path                  string            `toml:"path"`
	Commands              []string          `toml:"commands"`
	VersionArgs           []string          `toml:"version_args"`
	LaunchMode            string            `toml:"launch_mode"`
	HeadlessMode          string            `toml:"headless_mode"`
	HeadlessArgs          []string          `toml:"headless_args"`
	ACPArgs               []string          `toml:"acp_args"`
	InteractiveMode       string            `toml:"interactive_mode"`
	InteractiveCommand    string            `toml:"interactive_command"`
	InteractiveArgs       []string          `toml:"interactive_args"`
	InteractivePromptMode string            `toml:"interactive_prompt_mode"`
	InteractiveInvoke     string            `toml:"interactive_invoke"`
	InteractiveTimeoutMS  int               `toml:"interactive_timeout_ms"`
	InteractivePollMS     int               `toml:"interactive_poll_ms"`
	InteractiveNotes      string            `toml:"interactive_notes"`
	PromptMode            string            `toml:"prompt_mode"`
	SandboxMode           string            `toml:"sandbox_mode"`
	ApprovalPolicy        string            `toml:"approval_policy"`
	Model                 string            `toml:"model"`
	ModelLabel            string            `toml:"model_label"`
	Reasoning             string            `toml:"reasoning"`
	Profile               string            `toml:"profile"`
	Speed                 string            `toml:"speed"`
	TimeoutMS             int               `toml:"timeout_ms"`
	FirstEventTimeoutMS   *int              `toml:"first_event_timeout_ms"`
	StallTimeoutMS        *int              `toml:"stall_timeout_ms"`
	HeartbeatMS           *int              `toml:"heartbeat_ms"`
	IsolateHome           *bool             `toml:"isolate_home"`
	BuffersStdout         *bool             `toml:"buffers_stdout"`
	IsolatedHomeEnv       map[string]string `toml:"isolated_home_env"`
	ExternalBackend       string            `toml:"external_backend"`
	Telemetry             string            `toml:"telemetry"`
	Notes                 string            `toml:"notes"`
}

// configLayers lists the agent/defaults config files in low-to-high precedence:
// the user-global central default (~/.parley/agents.toml) seeds every project,
// the project deck files override it, and $PARLEY_HEADLESS_AGENT_CONFIG wins.
func configLayers(root string) []configLayer {
	deck := filepath.Join(root, protocol.DeckDir)
	layers := []configLayer{}
	if central := CentralAgentsPath(); central != "" {
		layers = append(layers, configLayer{path: central, source: "~/.parley/agents.toml", optional: true, machine: true})
	}
	layers = append(layers,
		configLayer{path: filepath.Join(deck, "agents.toml"), source: "parley-deck/agents.toml", optional: true, membership: true},
		configLayer{path: filepath.Join(deck, "agents.local.toml"), source: "parley-deck/agents.local.toml", optional: true},
	)
	if envPath := strings.TrimSpace(os.Getenv(EnvAgentConfig)); envPath != "" {
		if !filepath.IsAbs(envPath) {
			envPath = filepath.Join(root, envPath)
		}
		layers = append(layers, configLayer{path: envPath, source: EnvAgentConfig + ":" + envPath, optional: false})
	}
	return layers
}

func LoadAgentSpecs(root string) ([]agents.Spec, error) {
	specs := cloneSpecs(agents.DefaultSpecs())
	for _, item := range configLayers(root) {
		var err error
		specs, err = applyFile(root, specs, item.path, item.source, item.optional)
		if err != nil {
			return nil, err
		}
	}
	return specs, nil
}

// LoadDefaults merges the [defaults] policy block across the layered config
// files (central first, deck overrides). Missing files are skipped; later
// non-empty values win. Consumers (preflight ping tier, init transport) read
// the result and fall back to their own default on a zero value.
func LoadDefaults(root string) (CentralDefaults, error) {
	var out CentralDefaults
	for _, item := range configLayers(root) {
		data, err := os.ReadFile(item.path)
		if err != nil {
			if item.optional && errors.Is(err, os.ErrNotExist) {
				continue
			}
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		var cfg fileConfig
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return out, fmt.Errorf("%s: %w", item.path, err)
		}
		if cfg.Defaults != nil {
			mergeDefaults(&out, cfg.Defaults)
		}
	}
	return out, nil
}

// LoadRosterAdapters merges `[roster.<id>] adapter = "<family>"` across the layered
// config files (central < deck < local < env), later layers winning per id. The
// result feeds the participant resolver so a roster whose IDs are `claude-1`, … can
// be resolved to families fail-closed (idea composite-agent-naming-and-roster-reinit).
func LoadRosterAdapters(root string) (map[string]string, error) {
	out := map[string]string{}
	for _, item := range configLayers(root) {
		data, err := os.ReadFile(item.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		var cfg fileConfig
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return out, fmt.Errorf("%s: %w", item.path, err)
		}
		for id, ra := range cfg.Roster {
			if fam := strings.TrimSpace(ra.Adapter); fam != "" {
				out[id] = fam
			}
		}
	}
	return out, nil
}

// RosterAdaptersInFile parses only the `[roster.*]` mappings in a single config
// file (not the layered stack), so `parley roster init` can decide idempotency
// against the exact target file it is about to write rather than an inherited
// layer (review MAJOR, codex-1). A missing file yields an empty map, nil error.
func RosterAdaptersInFile(path string) (map[string]string, error) {
	out := map[string]string{}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return out, nil
		}
		return out, err
	}
	var cfg fileConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return out, fmt.Errorf("%s: %w", path, err)
	}
	for id, ra := range cfg.Roster {
		if fam := strings.TrimSpace(ra.Adapter); fam != "" {
			out[id] = fam
		}
	}
	return out, nil
}

// MachineFamilyCatalog returns the set of agent families known machine-wide —
// built-in specs plus the central ~/.parley/agents.toml [agents.*] keys — with NO
// deck/local/env layer. `parley roster init --scope machine` uses it so a deck-only
// custom family is never proposed for, written to, or blessed in the central file
// (consensus §B "never copies deck values up"; review MAJOR, codex-1).
func MachineFamilyCatalog() (map[string]bool, error) {
	out := map[string]bool{}
	for _, s := range agents.DefaultSpecs() {
		out[s.ID] = true
	}
	central := CentralAgentsPath()
	if central == "" {
		return out, nil
	}
	data, err := os.ReadFile(central)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return out, nil
		}
		return out, err
	}
	var cfg fileConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return out, fmt.Errorf("%s: %w", central, err)
	}
	for id := range cfg.Agents {
		out[id] = true
	}
	return out, nil
}

// ValidateAgentsConfigBytes reports whether data parses as a valid agents.toml
// document. `parley roster init` calls it on the candidate BEFORE the atomic
// replace, so a candidate that would install a duplicate or malformed `[roster.*]`
// table (e.g. against an existing quoted-key empty block) is rejected instead of
// silently breaking every later config load (review MAJOR, codex-1).
func ValidateAgentsConfigBytes(data []byte) error {
	var cfg fileConfig
	return toml.Unmarshal(data, &cfg)
}

func mergeDefaults(out *CentralDefaults, gd *globalDefaults) {
	if s := strings.TrimSpace(gd.Speed); s != "" {
		out.Speed = s
	}
	if s := strings.TrimSpace(gd.PingTier); s != "" {
		out.PingTier = s
	}
	if s := strings.TrimSpace(gd.PreferredTransport); s != "" {
		out.PreferredTransport = s
	}
	if s := strings.TrimSpace(gd.RosterChangePolicy); s != "" {
		out.RosterChangePolicy = s
	}
	if gd.Timeouts != nil {
		if gd.Timeouts.SignoffMS > 0 {
			out.SignoffMS = gd.Timeouts.SignoffMS
		}
		if gd.Timeouts.RoundMS > 0 {
			out.RoundMS = gd.Timeouts.RoundMS
		}
		if gd.Timeouts.ReviewMS > 0 {
			out.ReviewMS = gd.Timeouts.ReviewMS
		}
		if gd.Timeouts.DeepReasoningMS > 0 {
			out.DeepReasoningMS = gd.Timeouts.DeepReasoningMS
		}
	}
	if gd.Loop != nil {
		// Presence-aware (F-T2-1): a deliberate `= 0` overrides a lower layer's seed
		// to unlimited; absence (nil) falls through.
		if gd.Loop.MaxDriverSteps != nil {
			out.MaxDriverSteps = *gd.Loop.MaxDriverSteps
		}
		if gd.Loop.MaxWallClockMS != nil {
			out.MaxWallClockMS = *gd.Loop.MaxWallClockMS
		}
		if gd.Loop.MaxCostUSD != nil {
			out.MaxCostUSD = *gd.Loop.MaxCostUSD
		}
	}
}

// CentralHome returns the user-global Parley config directory (~/.parley),
// honoring the PARLEY_HOME override. It returns "" when the home directory
// cannot be resolved, so the central-default layer is simply skipped.
func CentralHome() string {
	if v := strings.TrimSpace(os.Getenv(EnvParleyHome)); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return ""
	}
	return filepath.Join(home, ".parley")
}

// CentralAgentsPath returns ~/.parley/agents.toml, or "" if home is unresolved.
func CentralAgentsPath() string {
	home := CentralHome()
	if home == "" {
		return ""
	}
	return filepath.Join(home, "agents.toml")
}

// EnsureCentralDefault writes a starter ~/.parley/agents.toml that lists the
// built-in agents with their default model and reasoning, if the file does not
// already exist. It is a no-op when the file is present or home is unresolved.
// Returns the ensured path ("" when skipped).
func EnsureCentralDefault() (string, error) {
	path := CentralAgentsPath()
	if path == "" {
		return "", nil
	}
	if _, err := os.Stat(path); err == nil {
		return path, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(centralDefaultTemplate()), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// centralDefaultTemplate renders a starter central config from the built-in
// specs: one [agents.<id>] block per agent with model + reasoning. Default
// reasoning should be the strongest level each agent supports; edit to taste.
func centralDefaultTemplate() string {
	var b strings.Builder
	b.WriteString("# ~/.parley/agents.toml — central per-user Parley Deck defaults.\n")
	b.WriteString("# Lists each agent's default model and reasoning/effort level. Used by\n")
	b.WriteString("# every project unless a deck overrides it in parley-deck/agents.toml.\n")
	b.WriteString("# Prefer an exact model id over a vendor \"latest\" alias, and keep the\n")
	b.WriteString("# strongest (highest) reasoning level each agent supports.\n\n")
	b.WriteString("# Project-wide policy defaults; a deck's parley-deck/agents.toml overrides them.\n")
	b.WriteString("[defaults]\n")
	b.WriteString("speed = \"fast\"                             # fast output at the SAME model+effort (Claude Code /fast), NOT a downgrade; a separate axis from reasoning. Use \"deep\" per idea for heavy work.\n")
	b.WriteString("ping_tier = \"hosted-pong\"                 # §9.0 roster liveness ping before each idea (or \"none\")\n")
	b.WriteString("preferred_transport = \"local-dir\"          # parley init default transport (local-dir|github-pr|gitlab-mr)\n")
	b.WriteString("roster_change_policy = \"confirm-breaking\"  # auto-add new agents; user confirms drops/breaking changes\n\n")
	b.WriteString("[defaults.timeouts]\n")
	b.WriteString("signoff_ms = 600000          # 10 min\n")
	b.WriteString("round_ms = 1200000           # 20 min\n")
	b.WriteString("review_ms = 1200000          # 20 min\n")
	b.WriteString("deep_reasoning_ms = 1200000  # 20 min\n\n")
	b.WriteString("[defaults.loop]              # auto-drive loop ceilings (LE-5); breach escalates, never completes. 0 = unlimited.\n")
	b.WriteString("max_driver_steps = 200       # total progress steps before escalation (generous safety net)\n")
	b.WriteString("max_wall_clock_ms = 7200000  # 2 h total run budget (distinct from the per-tick 30 min round deadline)\n")
	b.WriteString("max_cost_usd = 0             # 0 = unlimited; best-effort, telemetry-gated (LE-6)\n\n")
	for _, spec := range agents.DefaultSpecs() {
		model := spec.Model
		if strings.TrimSpace(model) == "" {
			model = agents.CLIDefault
		}
		reasoning := spec.Reasoning
		if strings.TrimSpace(reasoning) == "" {
			reasoning = agents.CLIDefault
		}
		b.WriteString(fmt.Sprintf("[agents.%s]\n", spec.ID))
		b.WriteString(fmt.Sprintf("model = %q\n", model))
		b.WriteString(fmt.Sprintf("reasoning = %q\n\n", reasoning))
	}
	return b.String()
}

func ExpandPlaceholders(value, root, tempdir string) string {
	replacements := map[string]string{
		"{root}":    root,
		"{deck}":    filepath.Join(root, protocol.DeckDir),
		"{tempdir}": tempdir,
	}
	for from, to := range replacements {
		value = strings.ReplaceAll(value, from, to)
	}
	return value
}

func applyFile(root string, specs []agents.Spec, path, source string, optional bool) ([]agents.Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if optional && errors.Is(err, os.ErrNotExist) {
			return specs, nil
		}
		return nil, err
	}

	var cfg fileConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	// Global [defaults].speed applies to every spec (overriding the built-in
	// per-agent default); a per-agent override in this same file still wins below.
	if cfg.Defaults != nil && strings.TrimSpace(cfg.Defaults.Speed) != "" {
		for i := range specs {
			specs[i].Speed = cfg.Defaults.Speed
			if specs[i].Sources == nil {
				specs[i].Sources = map[string]string{}
			}
			specs[i].Sources["speed"] = source + ":defaults"
		}
	}

	if len(cfg.Agents) == 0 {
		return specs, nil
	}

	byID := map[string]int{}
	for i, spec := range specs {
		byID[spec.ID] = i
	}
	for id, override := range cfg.Agents {
		index, ok := byID[id]
		if !ok {
			specs = append(specs, agents.Spec{
				ID:                    id,
				PromptMode:            agents.PromptStdin,
				LaunchMode:            agents.LaunchHeadless,
				InteractivePromptMode: agents.InteractivePromptNone,
				InteractiveInvoke:     agents.InteractiveInvokePrintOnly,
				InteractivePollMS:     agents.DefaultInteractivePollMS,
				Model:                 agents.CLIDefault,
				Reasoning:             agents.CLIDefault,
				Profile:               agents.CLIDefault,
				Speed:                 agents.DefaultSpeed,
				TimeoutMS:             agents.DefaultTimeoutMS,
				ExternalBackend:       agents.ExternalUnknown,
				Sources:               configDefaultSources(source),
			})
			index = len(specs) - 1
			byID[id] = index
		}
		spec := specs[index]
		spec = applyOverride(root, spec, override, source)
		specs[index] = spec
	}

	return specs, nil
}

func configDefaultSources(source string) map[string]string {
	sources := map[string]string{}
	for _, field := range []string{
		"id",
		"launch_mode",
		"prompt_mode",
		"interactive_prompt_mode",
		"interactive_invoke",
		"interactive_poll_ms",
		"model",
		"reasoning",
		"profile",
		"speed",
		"timeout_ms",
		"external_backend",
	} {
		sources[field] = source + ":default"
	}
	return sources
}

func applyOverride(root string, spec agents.Spec, override agentOverride, source string) agents.Spec {
	tempdir := os.TempDir()
	if spec.Sources == nil {
		spec.Sources = map[string]string{}
	}

	if override.Path != "" {
		spec.Commands = []string{ExpandPlaceholders(override.Path, root, tempdir)}
		spec.Sources["commands"] = source
	} else if override.Command != "" {
		spec.Commands = []string{ExpandPlaceholders(override.Command, root, tempdir)}
		spec.Sources["commands"] = source
	} else if len(override.Commands) > 0 {
		spec.Commands = expandSlice(override.Commands, root, tempdir)
		spec.Sources["commands"] = source
	}
	if len(override.VersionArgs) > 0 {
		spec.VersionArgs = expandSlice(override.VersionArgs, root, tempdir)
		spec.Sources["version_args"] = source
	}
	if override.LaunchMode != "" {
		spec.LaunchMode = override.LaunchMode
		spec.Sources["launch_mode"] = source
	}
	if override.HeadlessMode != "" {
		spec.HeadlessMode = override.HeadlessMode
		spec.Sources["headless_mode"] = source
	}
	if len(override.HeadlessArgs) > 0 {
		args := expandSlice(override.HeadlessArgs, root, tempdir)
		// LEGACY NORMALIZER (D7's second half). A config layer replaces headless_args
		// wholesale, so an override that spells out `--model <literal>` silently outranks
		// the `model` field beside it — the exact declared-vs-effective split this change
		// exists to remove, just relocated into the operator's own file. Rewrite the
		// literal back to the placeholder so the declared field wins; record it so the
		// rewrite is visible rather than magic.
		if normalized, changed := agents.NormalizeLegacyModelArgs(args); changed {
			args = normalized
			if spec.Sources == nil {
				spec.Sources = map[string]string{}
			}
			spec.Sources["headless_args_normalized"] = source
		}
		spec.HeadlessArgs = args
		spec.Sources["headless_args"] = source
	}
	if override.ACPArgs != nil {
		spec.ACPArgs = expandSlice(override.ACPArgs, root, tempdir)
		spec.Sources["acp_args"] = source
	}
	if override.InteractiveMode != "" {
		spec.InteractiveMode = override.InteractiveMode
		spec.Sources["interactive_mode"] = source
	}
	if override.InteractiveCommand != "" {
		spec.InteractiveCommand = ExpandPlaceholders(override.InteractiveCommand, root, tempdir)
		spec.Sources["interactive_command"] = source
	}
	if len(override.InteractiveArgs) > 0 {
		spec.InteractiveArgs = expandSlice(override.InteractiveArgs, root, tempdir)
		spec.Sources["interactive_args"] = source
	}
	if override.InteractivePromptMode != "" {
		spec.InteractivePromptMode = override.InteractivePromptMode
		spec.Sources["interactive_prompt_mode"] = source
	}
	if override.InteractiveInvoke != "" {
		spec.InteractiveInvoke = override.InteractiveInvoke
		spec.Sources["interactive_invoke"] = source
	}
	if override.InteractiveTimeoutMS > 0 {
		spec.InteractiveTimeoutMS = override.InteractiveTimeoutMS
		spec.Sources["interactive_timeout_ms"] = source
	}
	if override.InteractivePollMS > 0 {
		spec.InteractivePollMS = override.InteractivePollMS
		spec.Sources["interactive_poll_ms"] = source
	}
	if override.InteractiveNotes != "" {
		spec.InteractiveNotes = override.InteractiveNotes
		spec.Sources["interactive_notes"] = source
	}
	if override.PromptMode != "" {
		spec.PromptMode = agents.PromptMode(override.PromptMode)
		spec.Sources["prompt_mode"] = source
	}
	if override.SandboxMode != "" {
		spec.SandboxMode = override.SandboxMode
		spec.Sources["sandbox_mode"] = source
	}
	if override.ApprovalPolicy != "" {
		spec.ApprovalPolicy = override.ApprovalPolicy
		spec.Sources["approval_policy"] = source
	}
	if override.Model != "" {
		spec.Model = override.Model
		spec.Sources["model"] = source
	}
	if override.ModelLabel != "" {
		spec.ModelLabel = override.ModelLabel
		spec.Sources["model_label"] = source
	}
	if override.Reasoning != "" {
		spec.Reasoning = override.Reasoning
		spec.Sources["reasoning"] = source
	}
	if override.Profile != "" {
		spec.Profile = override.Profile
		spec.Sources["profile"] = source
	}
	if override.Speed != "" {
		spec.Speed = override.Speed
		spec.Sources["speed"] = source
	}
	if override.TimeoutMS > 0 {
		spec.TimeoutMS = override.TimeoutMS
		spec.Sources["timeout_ms"] = source
	}
	// Supervision knobs are pointer-typed so an explicit `0` (disable) is
	// distinguishable from "not set"; 0 maps to -1 (disabled) on the Spec.
	applySupervisionMS := func(value *int, field *int, key string) {
		if value == nil {
			return
		}
		if *value == 0 {
			*field = -1
		} else {
			*field = *value
		}
		spec.Sources[key] = source
	}
	applySupervisionMS(override.FirstEventTimeoutMS, &spec.FirstEventTimeoutMS, "first_event_timeout_ms")
	applySupervisionMS(override.StallTimeoutMS, &spec.StallTimeoutMS, "stall_timeout_ms")
	applySupervisionMS(override.HeartbeatMS, &spec.HeartbeatMS, "heartbeat_ms")
	if override.IsolateHome != nil {
		spec.IsolateHome = *override.IsolateHome
		spec.Sources["isolate_home"] = source
	}
	if override.BuffersStdout != nil {
		spec.BuffersStdout = *override.BuffersStdout
		spec.Sources["buffers_stdout"] = source
	}
	if len(override.IsolatedHomeEnv) > 0 {
		spec.IsolatedHomeEnv = map[string]string{}
		for key, value := range override.IsolatedHomeEnv {
			spec.IsolatedHomeEnv[key] = ExpandPlaceholders(value, root, "{tempdir}")
		}
		spec.Sources["isolated_home_env"] = source
	}
	if override.ExternalBackend != "" {
		spec.ExternalBackend = override.ExternalBackend
		spec.Sources["external_backend"] = source
	}
	if override.Telemetry != "" {
		spec.Telemetry = override.Telemetry
		spec.Sources["telemetry"] = source
	}
	if override.Notes != "" {
		spec.Notes = override.Notes
		spec.Sources["notes"] = source
	}

	return spec
}

func expandSlice(values []string, root, tempdir string) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = ExpandPlaceholders(value, root, tempdir)
	}
	return out
}

func cloneSpecs(specs []agents.Spec) []agents.Spec {
	out := make([]agents.Spec, len(specs))
	for i, spec := range specs {
		out[i] = spec
		out[i].Commands = append([]string(nil), spec.Commands...)
		out[i].VersionArgs = append([]string(nil), spec.VersionArgs...)
		out[i].HeadlessArgs = append([]string(nil), spec.HeadlessArgs...)
		out[i].ACPArgs = cloneOptionalStringSlice(spec.ACPArgs)
		out[i].InteractiveArgs = append([]string(nil), spec.InteractiveArgs...)
		if spec.IsolatedHomeEnv != nil {
			out[i].IsolatedHomeEnv = map[string]string{}
			for key, value := range spec.IsolatedHomeEnv {
				out[i].IsolatedHomeEnv[key] = value
			}
		}
		if spec.Sources != nil {
			out[i].Sources = map[string]string{}
			for key, value := range spec.Sources {
				out[i].Sources[key] = value
			}
		}
	}
	return out
}

func cloneOptionalStringSlice(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string{}, values...)
}

// RosterEntriesInFile parses the [roster.*] blocks of ONE file rather than the layered
// stack. `roster sync` needs to compare a deck's own declarations against the machine's,
// and a layered read would already have merged them.
func RosterEntriesInFile(path string) (map[string]RosterEntry, error) {
	out := map[string]RosterEntry{}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return out, nil
		}
		return out, err
	}
	var cfg fileConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return out, fmt.Errorf("%s: %w", path, err)
	}
	for id, ra := range cfg.Roster {
		e := RosterEntry{
			ID:           id,
			Adapter:      strings.TrimSpace(ra.Adapter),
			Active:       true,
			Model:        strings.TrimSpace(ra.Model),
			Effort:       strings.TrimSpace(ra.Effort),
			Speed:        strings.TrimSpace(ra.Speed),
			WorkspaceDir: strings.TrimSpace(ra.WorkspaceDir),
			Role:         strings.TrimSpace(ra.Role),
			HostHandle:   strings.TrimSpace(ra.HostHandle),
		}
		if ra.Active != nil {
			e.Active = *ra.Active
		}
		out[id] = e
	}
	return out, nil
}

// RosterFieldSources reports which config layer last set each [roster.<id>] field, for
// `roster show --explain`. D3 parked per-field provenance here rather than in a SOURCE
// column: one column cannot describe eleven fields whose values come from different
// files. An absent field means no layer set it (the built-in default applies).
func RosterFieldSources(root, id string) (map[string]string, error) {
	out := map[string]string{}
	for _, item := range configLayers(root) {
		_, entries, err := rosterLayer(item.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		ra, ok := entries[id]
		if !ok {
			continue
		}
		if strings.TrimSpace(ra.Adapter) != "" {
			out["adapter"] = item.source
		}
		if ra.Active != nil {
			out["active"] = item.source
		}
		if strings.TrimSpace(ra.Model) != "" {
			out["model"] = item.source
		}
		if strings.TrimSpace(ra.Effort) != "" {
			out["effort"] = item.source
		}
		if strings.TrimSpace(ra.Speed) != "" {
			out["speed"] = item.source
		}
	}
	return out, nil
}

// RosterSourcePath resolves a layer's display label back to the file it names, so callers
// comparing "which layer won" against a path they just wrote compare paths with paths.
func RosterSourcePath(root, source string) string {
	for _, item := range configLayers(root) {
		if item.source == source {
			return item.path
		}
	}
	return ""
}

// machineOnlyLayers narrows the layer stack to the user-global config, for `--scope
// machine`. Without it, a machine-scope query answered with deck values.
func machineOnlyLayers(layers []configLayer) []configLayer {
	out := make([]configLayer, 0, 1)
	for _, item := range layers {
		if item.machine {
			out = append(out, item)
		}
	}
	return out
}

// LoadAgentSpecsScoped is LoadAgentSpecs, optionally restricted to the machine layer.
func LoadAgentSpecsScoped(root string, machineOnly bool) ([]agents.Spec, error) {
	layers := configLayers(root)
	if machineOnly {
		layers = machineOnlyLayers(layers)
	}
	specs := cloneSpecs(agents.DefaultSpecs())
	for _, item := range layers {
		var err error
		specs, err = applyFile(root, specs, item.path, item.source, item.optional)
		if err != nil {
			return nil, err
		}
	}
	return specs, nil
}

// LoadRosterAdaptersScoped is LoadRosterAdapters, optionally restricted to the machine layer.
func LoadRosterAdaptersScoped(root string, machineOnly bool) (map[string]string, error) {
	layers := configLayers(root)
	if machineOnly {
		layers = machineOnlyLayers(layers)
	}
	out := map[string]string{}
	for _, item := range layers {
		_, entries, err := rosterLayer(item.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		for id, ra := range entries {
			if v := strings.TrimSpace(ra.Adapter); v != "" {
				out[id] = v
			}
		}
	}
	return out, nil
}

// RosterFieldSourcesScoped is RosterFieldSources, optionally restricted to the machine layer.
func RosterFieldSourcesScoped(root, id string, machineOnly bool) (map[string]string, error) {
	layers := configLayers(root)
	if machineOnly {
		layers = machineOnlyLayers(layers)
	}
	out := map[string]string{}
	for _, item := range layers {
		_, entries, err := rosterLayer(item.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		ra, ok := entries[id]
		if !ok {
			continue
		}
		for field, set := range map[string]bool{
			"adapter": strings.TrimSpace(ra.Adapter) != "",
			"active":  ra.Active != nil,
			"model":   strings.TrimSpace(ra.Model) != "",
			"effort":  strings.TrimSpace(ra.Effort) != "",
			"speed":   strings.TrimSpace(ra.Speed) != "",
		} {
			if set {
				out[field] = item.source
			}
		}
	}
	return out, nil
}

// AgentFieldSources reports which layer last set each [agents.<family>] field, so
// `roster show --explain` can attribute a value that reaches the launch through the agent
// block rather than through a [roster.<id>] block.
func AgentFieldSources(root, family string, machineOnly bool) map[string]string {
	layers := configLayers(root)
	if machineOnly {
		layers = machineOnlyLayers(layers)
	}
	out := map[string]string{}
	for _, item := range layers {
		data, err := os.ReadFile(item.path)
		if err != nil {
			continue
		}
		var cfg fileConfig
		if err := toml.Unmarshal(data, &cfg); err != nil {
			continue
		}
		ov, ok := cfg.Agents[family]
		if !ok {
			continue
		}
		if strings.TrimSpace(ov.Model) != "" {
			out["model"] = item.source
		}
		if strings.TrimSpace(ov.Reasoning) != "" {
			out["effort"] = item.source
		}
		if strings.TrimSpace(ov.Speed) != "" {
			out["speed"] = item.source
		}
	}
	return out
}

// RosterStateSource names the layer that decides a member's active/inactive state — the
// same layer that granted membership. `active` is deliberately NOT layered, so reporting
// its provenance from the general layer stack made `--explain` and the masking warning
// contradict the state `roster show` actually resolved.
func RosterStateSource(root string) (string, error) {
	scope, err := LoadRosterScoped(root)
	if err != nil {
		return "", err
	}
	return scope.Source, nil
}

// RosterStateSourceForTarget names the layer that decides active/inactive state for the
// scope the given file belongs to. A machine-scope write is governed by the machine
// roster; a deck-scope write by the deck's authority chain.
func RosterStateSourceForTarget(root, target string) (string, error) {
	if central := CentralAgentsPath(); central != "" {
		a, _ := filepath.Abs(central)
		b, _ := filepath.Abs(target)
		if a == b {
			return "~/.parley/agents.toml", nil
		}
	}
	scope, err := LoadRosterScoped(root)
	if err != nil {
		return "", err
	}
	return scope.Source, nil
}

// RosterMaskedFields reports the roster fields of `id` that a LOWER config layer declares and a
// HIGHER layer overrides with a different value — i.e. the fields for which someone's written
// declaration never reaches the launch.
//
// claude-1/F3, MAJOR, confirmed by @zcode-1 and @kimi-1 and never dispositioned: `masked-by-env`
// was in the closed STATUS vocabulary (`docs/cli-reference.md`) and in the docs, but nothing ever
// put it in a STATUS cell — `roster set` printed it as a one-off stderr warning after a write and
// `roster show` could not report it at all. A status you cannot observe is not a status.
//
// Only fields whose value actually CHANGES are reported. Two layers agreeing is not masking, and
// reporting it as such is the false-positive that a previous fix already had to undo for
// machine-scope writes.
func RosterMaskedFields(root, id string) ([]string, error) {
	seen := map[string]string{}
	masked := map[string]bool{}
	for _, item := range configLayers(root) {
		_, entries, err := rosterLayer(item.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		ra, ok := entries[id]
		if !ok {
			continue
		}
		for field, value := range map[string]string{
			"adapter": strings.TrimSpace(ra.Adapter),
			"model":   strings.TrimSpace(ra.Model),
			"effort":  strings.TrimSpace(ra.Effort),
			"speed":   strings.TrimSpace(ra.Speed),
		} {
			if value == "" {
				continue
			}
			if prev, had := seen[field]; had && prev != value {
				masked[field] = true
			}
			seen[field] = value
		}
	}
	out := make([]string, 0, len(masked))
	for field := range masked {
		out = append(out, field)
	}
	sort.Strings(out)
	return out, nil
}
