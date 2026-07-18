---
agent: codex-1
idea: composite-agent-naming-and-roster-reinit
round: 2
date: 2026-07-18
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

The two late requirements change two details, but not my answer to the crux:

1. Composite names should preserve version dots. The canonical generated display grammar is
   therefore a lowercase subset of `[a-zA-Z0-9_.-]`, with a model token such as `opus4.8` or
   `gpt5.5`. Dots are never permitted at a token boundary or consecutively.
2. `speed` is an independent launch axis. Its default becomes `fast`, but choosing it must
   leave the confirmed model and effort unchanged. Legacy `profiles.fast` entries that select
   Sonnet/low or another cheaper model/effort are no longer speed profiles.

My final crux position remains: **stable §2 roster ID = protocol identity; composite = derived
display name**. I strengthen my round-01 position on the namespace schism: this feature should
fix it with an explicit resolver, not merely avoid making it worse. A roster command that writes
`claude-1` while `runner.selectedAgents` recognizes only `claude` would ship a knowingly unusable
roster.

## Responses to other participants

### @claude-1

I agree with claude-1 on preserving dots, deriving the model token from an operator-confirmed
human label, using a closed effort vocabulary, allocating collision suffixes, reusing one roster
selection service from bootstrap, and making autonomous write behavior explicit per adapter.

I disagree with the frozen-slug identity. It protects one open idea, but the same logical seat
becomes `claude-opus4.8-max` in one idea and `claude-opus4.9-max` in the next. That fragments
presets, run history, TUI state, inbox addressing, signoffs, and cross-idea audit queries. It also
does not resolve the existing §2-ID/spec-ID split; it adds a third identity that must be aliased.

Concrete counter-proposal: retain `claude-1` for `participants:`, artifact paths, frontmatter,
events, run directories, and signoffs. Derive `claude-opus4.8-max` for display, and snapshot the
exact adapter/model/effort/display/speed profile in `00-prompt.md` and `run.json` when the idea
starts. This gives frozen audit evidence without turning mutable launch metadata into identity.
No historical artifact rename or alias graph is required.

I also prefer a distinct rerunnable `parley roster init` command over claude-1's mixed surface
that names both `roster init` and `parley init --reinit`. Fresh `parley init` should call the same
internal roster service once; subsequent changes use the roster verb.

### @codex-1

I agree with codex-1's stable-ID/display split, explicit `AdapterID`, immutable run-profile
snapshot, structured autonomous-write mode, fail-closed participant resolution, and warning that
metadata-only model selection is false success while `HeadlessArgs` still hard-code another
model.

I now disagree with my round-01 no-dot display grammar because the amended brief explicitly
requires natural version dots. Concrete replacement: model tokens match
`[a-z0-9]+(?:\.[a-z0-9]+)*`; the whole generated display matches the grammar below and rejects
leading/trailing dots and `..`.

I also refine my earlier flat `participant-profiles` sketch. The counter-proposal is a structured
map in `00-prompt.md` plus a typed `runmanifest.ParticipantProfile` slice. It records stable ID,
adapter ID, exact model, human model label, configured reasoning, effective effort label/source,
derived display, speed, and autonomous-write profile. Resume uses this snapshot instead of
newly reinitialized defaults.

### @hermes-1

I agree with hermes-1 that `parley roster init` is a maintenance surface distinct from deck
scaffolding, that session and machine writes match different config layers, that discovery needs
per-adapter model/effort probes with honest fallback, and that reinit must guard open ideas.

I disagree with two proposed escape hatches. First, this idea should not perform a one-time rewrite
from old spec-ID artifacts such as `claude.md` to `claude-1.md`; historical files and signoffs stay
untouched. Second, `--yes` must not turn roster maintenance into an implicit edit of open ideas'
`excluded:` fields. That is a quorum decision governed by §9.0, not an overwrite confirmation.

Concrete counter-proposal: use forward-only resolution. Legacy participants equal to an adapter
ID (`claude`) remain valid legacy identities and bind to the same-named adapter. New/reinitialized
rosters write explicit bindings such as `claude-1 -> claude`. Session reinit hard-refuses removal
or rename of an ID referenced by a nonterminal idea/run; it may mark that binding unavailable for
new ideas but retains the pinned binding until terminal. Adding a seat or changing only its future
profile is allowed. A separate, explicit roster-retire/quorum workflow can handle removal later.

Hermes-1 also allowed Claude `acceptEdits` as an autonomous choice. The amended requirement settles
this: built-in autonomous mode is `bypassPermissions`, with workspace enforcement supplied by the
launcher boundary; `acceptEdits` is not the autonomous default.

### @kimi-1

I agree with kimi-1's end-to-end trace: the runner currently keys on spec family IDs, §2 and this
idea key on stable roster IDs, `roster.go` accepts only lowercase alphanumeric/hyphen IDs, and the
signoff parser separately accepts `[A-Za-z0-9._-]`. I also agree that Claude's built-in
`acceptEdits` and Codex's `on-failure` are not the required zero-prompt defaults, that Kimi `-p`
must never be combined with `--yolo`/`--auto`, and that display data belongs beside rather than in
artifact keys.

I disagree with deferring the two-namespace fix to another idea. Component B necessarily creates
or refreshes stable roster IDs; without a resolver, `parley run --participants claude-1` still
selects nothing and cannot produce `round-01/claude-1.md`. "Contain the schism" is not an adequate
acceptance criterion for a roster initializer.

Concrete counter-proposal: add an explicit stable-ID-to-adapter binding to config and make all run,
preflight, signoff-request, pipeline, and TUI launch selection operate on bound specs. Keep the
strict §2 ID grammar unchanged; dotted composites appear only in a new display column. Vendor
branches switch from `Spec.ID` to `Spec.AdapterID`. The signoff parser's broader regex remains for
historical compatibility, but membership still compares the exact stable participant ID.

I also replace kimi-1's boolean `AutonomousWrite` assertion with a structured profile containing
mode, ordered args, scope, and enforcement. A boolean can say "yes" while the effective argv still
contains `acceptEdits` or `on-failure`; the structured profile can be validated against the final
invocation.

## Converged design

### A. Identity, adapter, and composite display

There are three deliberately separate values:

| Value | Example | Authority and uses |
|---|---|---|
| Stable identity | `claude-1` | §2 ID, `participants:`, filenames, frontmatter, signoffs, events, run/TUI keys |
| Adapter/family | `claude` | Built-in catalog, command discovery, invocation assembly, vendor-specific environment behavior |
| Derived display | `claude-opus4.8-max` | §2 annotation, runtime matrix, digest/TUI label, prompt/run profile snapshot |

Stable IDs remain in the intersection accepted everywhere:
`^[a-z0-9][a-z0-9-]*$`. Dots and underscores are not added to canonical identities even though the
consensus parser tolerates them. This keeps `internal/protocol/roster.go` authoritative and avoids
case-insensitive-filesystem surprises.

Generated composite displays are lowercase and follow:

```text
<family>-<model-token>-<effort-label>[-<collision-index>]

family          := [a-z0-9]+
model-token     := [a-z0-9]+(?:\.[a-z0-9]+)*
effort-label    := low|medium|high|xhigh|max|ultracode|clidefault
collision-index := integer >= 2
```

The generator lowercases a confirmed human model label, removes characters outside
`[a-z0-9.]`, collapses dot runs, strips leading/trailing dots, and rejects an empty token. It never
emits uppercase or underscore, although both remain within the user's outer allowed charset. The
whole-name validator independently rejects `..`, any leading/trailing dot, path separators, and a
non-round-tripping parse. Parsing is display-only and never reconstructs the exact launch model.

The required displays are:

- `claude-opus4.8-max`
- `codex-gpt5.5-xhigh`
- `hermes-glm5.2-high`
- `agy-gemini3.5flash-high`
- `kimi-k3-max`

Agy needs an explicit distinction: its confirmed independent `Reasoning` remains `cli-default`, as
required by component D, while `EffectiveEffortLabel=high` and `EffortSource=model-tier` come from
the selected `(High)` model variant. The model token omits that tier, so it appears once in the
display. This does not invent an Agy effort flag.

Collision allocation is persisted per stable ID (`display_instance`), not recomputed from discovery
order. Existing allocations survive reinit and removal; a new collision receives the lowest unused
suffix `2, 3, ...`. Thus removing the unsuffixed seat never silently renames `...-2`.

At idea creation, `00-prompt.md` snapshots each participant profile. `runmanifest.Manifest` gets the
same typed profile set under a schema bump. Current config derives displays; the snapshots freeze
what the idea/run actually used. Historical files are never renamed.

### Exact Go refactor and order

The safe implementation sequence is:

1. **Split types without changing current behavior.** In `internal/agents/discover.go`, make
   `Spec.ID` explicitly the stable protocol identity and add `AdapterID`, `DisplayName`,
   `ModelLabel`, `EffectiveEffortLabel`, `EffortSource`, `AutonomousWrite`, and typed speed-mode
   data. Introduce adapter templates keyed by `AdapterID`; initially bind legacy specs with
   `ID == AdapterID` so existing spec-ID ideas still work. Change `mergeACPCatalog` and
   `specFromACPBackend` in `internal/agents/acp_specs.go` to key templates by `AdapterID`.
2. **Move every vendor decision to `AdapterID`.** In `internal/runner/runner.go`, change
   `cleanParticipantEnv`, the Hermes environment branch in `buildAgentInvocation`, and
   `isolatedAgentHome` to inspect `agent.AdapterID`. Search/test every remaining family comparison.
   All artifact paths, validation, event keys, tracker keys, snapshots, and prompts continue using
   `agent.ID`.
3. **Add the resolver and eliminate silent omission.** Extend `internal/config/runtime.go` with an
   explicit `adapter` field. New config is keyed by stable identity:

   ```toml
   [agents.claude-1]
   adapter = "claude"
   model = "claude-opus-4-8[1m]"
   model_label = "Opus 4.8"
   reasoning = "max"
   speed = "fast"
   ```

   `LoadAgentSpecs` clones the named adapter template and applies the seat/profile overrides.
   Legacy `[agents.claude]` blocks remain adapter defaults; a legacy participant literally named
   `claude` can still be synthesized as identity `claude`. Add `ResolveParticipants` that accepts
   exact participant IDs and returns bound discoveries in the same order, failing on missing,
   duplicate, inactive, ambiguous, or unavailable bindings. Change `runner.selectedAgents` from a
   silent filter to an error-returning exact check.
4. **Make §2 and app selection use the same namespace.** Add a structured `ReadRoster` alongside
   `ReadRosterIDs` in `internal/protocol/roster.go`, plus a narrowly scoped roster-table writer. The
   ID cell keeps its current regex; adapter/display/model/effort/speed are separate columns.
   `discoverConfigured`, `installedAgentIDs`, `selectedParticipantIDs`,
   `participantDiscoveries`, `preflight.checkRoster`, signoff-request selection, pipeline launch,
   and driver consensus/implementation lookup must consume resolved stable IDs. The non-negotiable
   test is: an explicit binding `antigravity-1 -> agy` launches Agy and writes/events/signs as
   `antigravity-1`.
5. **Add display rendering and snapshots.** Put derivation/validation/collision code in one
   `internal/agents/naming.go`. Add display and exact profile fields to `runcontrol.RuntimeEventData`,
   `runmanifest.Manifest`, `driver.AgentLine`, the runtime matrix, and TUI rows. Digest/artifact
   reads remain keyed by stable ID. Resume resolves from the manifest snapshot first; reinitialized
   defaults cannot change an active run.
6. **Only then add roster writes and invocation-policy changes.** This prevents the new command
   from emitting stable IDs before the runner can resolve them.

This fixes the schism rather than broadening §2 IDs to accept display syntax. `consensus.go` needs
no regex narrowing; it should keep reading historical signoffs and continue exact membership
validation against `participants:`.

### B. Command surface and scopes

Use a distinct, rerunnable verb:

```text
parley roster init [--scope session|machine] [--from FILE] [--dry-run] [--yes] [--json]
parley roster show [--scope session|machine] [--json]
```

There is no `--reinit`; running `roster init` again is reinitialization. `parley init` keeps deck
creation/transport ownership but calls the same internal `RosterInit` service once for a fresh deck.
It must not contain a second picker. In unattended bootstrap, `--from` supplies the confirmed
selection; otherwise bootstrap reports that roster confirmation remains required.

Scope behavior is exact:

- `session` is the default only when a deck exists. It writes the full confirmed portable roster
  bindings/profile selections to `parley-deck/agents.toml`, renders §2, and regenerates only the
  roster-owned entries in `meta/headless-agents.local.json` (resolved paths/argv for the skill and
  manual orchestrators). The TOML plus §2 IDs are runtime/canonical authorities; the JSON is a
  versioned projection with a source fingerprint, not a competing config layer.
- `machine` must be explicit and writes only the managed roster/profile portion of
  `~/.parley/agents.toml`, preserving `[defaults]`, presets, unknown keys, and comments. It never
  edits a deck. Outside a deck, omitted scope is an error rather than a surprising machine-wide
  write.

Both paths discover adapter CLIs first, use adapter-specific model/effort catalogs where available,
and show configured value, `cli-default`, and exact free-text fallback when discovery cannot prove
options. Coupled choices such as Agy model+tier are presented as pairs. The service builds and
validates every output before replacement, uses managed-block/lossless writes for user TOML, shows
one diff, and makes `--yes` mean only "apply this shown diff."

Session reinit scans nonterminal ideas and runs twice (plan time and immediately before write). It
hard-refuses renaming/removing a referenced stable ID; it never edits open `participants:` or
`excluded:` fields. Model/effort/speed changes apply to future ideas, while active runs use their
snapshots. Machine reinit cannot scan all projects, which is another reason the run snapshot is
mandatory. §9.0 remains a read-only liveness check over the already resolved roster.

### C. Actually autonomous built-in write profiles

Use a first-class structured value, not a boolean hidden behind arbitrary `HeadlessArgs`:

```go
type AutonomousWriteMode struct {
    Mode             string
    Args             []string
    Scope            string // "workspace"
    ScopeEnforcement string // "cli-sandbox" or "outer-sandbox"
}
```

Factor model, reasoning, autonomous, and speed fragments out of hard-coded argv. A structured
assembler called by `buildAgentInvocation` must expand ordered list slots so Agy can still keep its
value-taking prompt last. Tests assert the final argv, not just `Spec` metadata.

Built-in/default mappings are:

| Adapter | Autonomous mapping |
|---|---|
| Claude | `--permission-mode bypassPermissions`, rooted at the workspace; require an outer workspace sandbox because `--add-dir` alone is not confinement |
| Codex | `-s workspace-write` plus non-prompting approval policy `never`; remove built-in `on-failure` |
| Hermes | `--yolo` (retain `--accept-hooks` where required), under the workspace guard |
| Agy | `--dangerously-skip-permissions` with workspace root, under the workspace guard |
| Kimi | plain headless `-p`; reject any effective argv combining it with `--yolo` or `--auto` |

Accordingly, change Claude's built-in away from `acceptEdits`, Codex away from `on-failure`, add
Hermes `--yolo`, retain Agy's mapping, and provide Kimi's headless `-p` profile in addition to its
ACP catalog capability. Cwd/`--add-dir` select a workspace but are not themselves a filesystem
sandbox; where the CLI does not provide confinement, `ScopeEnforcement=outer-sandbox` is required
and verification fails closed if the launcher cannot supply it.

The skill wording should state that every headless participant uses its mapped non-interactive
autonomous-write mode, that the mode must be confined to the assigned workspace, that Kimi's `-p`
is mutually exclusive with `--yolo`/`--auto`, and that permission automation never relaxes secret
redaction or external-data-disclosure rules.

### D. `fast` is orthogonal to effort

Change `agents.DefaultSpeed` and the central seed to `fast`; `roster init` writes `speed = "fast"`
for new selections. A speed mode may contain only an activation mechanism, timeout/pacing policy,
and notes—never `model`, `model_label`, `reasoning`, or effort. Invocation assembly applies the
already confirmed model/effort first, then the speed mechanism, and validates that the effective
model/effort did not change.

- `fast`: same model and effort, activate Claude `/fast` or an adapter's proven equivalent; if an
  installed headless adapter has no equivalent, record a documented no-op rather than substitute a
  smaller model or lower effort.
- `deep`: same model and effort, disable fast-output acceleration and use the deeper timeout/pacing
  budget for heavy ideas. It is not permission to mutate reasoning.

Legacy `balanced` normalizes to `fast` with a warning. Legacy `profiles.fast` blocks that contain
model/thinking overrides are ignored for speed selection, reported as downgrade-incompatible, and
rewritten by reinit into speed-only mappings. The confirmed efforts stay Claude=max, Codex=xhigh,
Hermes=high, Agy=cli-default (effective label high from its model tier), and Kimi=max. Speed is shown
separately in snapshots/runtime output and is never appended to the composite name.

## Verification and consensus-ready acceptance criteria

- Naming tests produce the five dotted displays exactly; reject empty tokens, `..`, leading or
  trailing dots, separators, ambiguous parses, and unstable collision reassignment.
- Resolver tests cover two stable seats on one adapter, explicit `antigravity-1 -> agy`, legacy
  `claude -> claude`, missing/ambiguous bindings, and the invariant that
  `round-NN/<stable-id>.md`, frontmatter, events, signoffs, reviews, and TUI keys never use the
  adapter or display.
- Reinit tests prove session and machine touch only their specified files, preserve user-owned TOML,
  refuse open-idea ID removal, and leave an active run on its manifest profile after defaults change.
- Invocation tests prove the exact five autonomous mappings, reject Kimi flag conflicts, and prove
  built-in Claude/Codex no longer use `acceptEdits`/`on-failure`.
- Speed tests start each adapter at `fast` and assert byte-for-byte equality of effective
  model/reasoning before and after speed switching; legacy downgrade profiles emit a warning and
  never alter launch quality.
- `go build ./...`, `go vet ./...`, `gofmt -l`, the protocol drift guard, CLI/skill config fixtures,
  and cross-repo skill fallback synchronization all stay green.

This is one coherent design: stable IDs preserve protocol continuity; explicit adapter bindings
make those IDs executable; dotted composites make profiles readable; snapshots preserve audit
truth; roster init owns all selection; autonomous mode removes write prompts inside a verified
workspace boundary; and speed improves output latency without touching quality.
