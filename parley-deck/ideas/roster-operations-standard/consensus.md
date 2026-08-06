---
idea: roster-operations-standard
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: standard
rounds: 2
date: 2026-08-06
status: consensus
---

## What was decided

The user asked for a stable roster table, easy local/global update, and one defined sync. The run
found that none of the three could be built on the current data, because **the roster is not one
object and the table reports values the launcher does not use**. The decisions below fix the data
first and freeze the format second — that ordering is itself a decision, and it is unanimous.

### 1. Three concepts, one answer — unanimous

- **adapter inventory** — what this machine can launch. `parley agents list`, relabelled
  *"adapter/runtime inventory — not the roster"*.
- **roster** — who is on this deck's team, with what settings. `parley roster show`. This is **the**
  answer to "what is the current agent roster?".
- **run snapshot** — what a given run actually used. Immutable, written at run creation.

`parley roster show` must appear in `parley --help` and in the docs. Today it is dispatched
(`app.go:100`) but omitted from `printUsage` (`app.go:111-144`) and absent from
`docs/cli-reference.md` and `docs/agent-runtime-configuration.md` — kimi-1, `PRIMARY`. **The command
the skill tells agents to run is undocumented and unlisted**, which is why every session invented
its own answer.

### 2. `MODEL` and `EFFORT` are effective-or-`unknown` — unanimous

Never a declaration occupying the effective cell. Divergence surfaces in `STATUS`
(`model-drift`, `effort-unknown`), never by silently printing the declared value.

This is the §15/1.39.0 `AUTO` rule generalized. All four participants independently confirmed the
current defect; three verified it in source. Measured today (`PRIMARY`, facilitator):

| adapter | declared MODEL | model in launch argv | effort in launch argv |
|---|---|---|---|
| claude | `claude-opus-5[1m]` | `claude-opus-4-8[1m]` | `--effort max` |
| codex | `gpt-5.6-sol` | *none* | *none* |
| kimi | `kimi-code/k3` | *none* | *none* |
| hermes, agy, opencode | — | matches | *none* |

**Six of seven adapters pass no effort flag at all**, so `EFFORT` is a pure declaration for every
agent but `claude`. Root cause: `applyOverride` (`runtime.go:594-596`) sets `spec.Model` and never
touches `HeadlessArgs`; `buildAgentInvocation` (`runner.go:1097-1108`) substitutes only `{root}`
and `{prompt}`.

### 3. Frozen column contract v1 — agreed except where noted in `## Verdict conflicts`

```
AGENT  ADAPTER  STATE  INSTALLED  MODEL  MODEL-FAMILY  MODEL-COMPANY  EFFORT  SPEED  AUTO  STATUS
```

Unanimous on: this order, these eleven, `DISPLAY-NAME` **removed** from the table, and the header
plus `--json` `schema_version`/`columns` being an **API** with golden tests and additive-only
change. `VERSION`, `SANDBOX`, `APPROVAL`, `TIMEOUT`, `HOME`, `BACKEND` stay in `agents list`.

`DISPLAY-NAME` is removed because it contradicts `MODEL` today: `RenderDisplayName` prefers
`ModelLabel` (`naming.go:189-191`), the built-in label is `Opus 4.8 1m` (`discover.go:227`), and
decks pin `model` without `model_label` — kimi-1, `PRIMARY`. It survives for TUI/artifact rendering,
derived from the effective row.

`STATUS` carries a closed vocabulary: `ok`, `unmapped`, `not-installed`, `model-drift`,
`model-unbound`, `effort-unknown`, `metadata-unknown`, `masked-by-env`, `legacy-roster`,
`inactive`, `stale-snapshot`.

### 4. Derivation is CLI-owned — unanimous

A versioned, tested `modelmeta` registry ships with the CLI. Peel recognized gateway prefixes
(`litellm/`, `openrouter/`) before deriving company; never infer company from the adapter; return
`unknown` + `metadata-unknown` on no match. **No deck ever hand-writes `model_family`/
`model_company`.** The release that adds or changes a rostered model updates the registry and its
golden tests in the same change.

Worked cases: `litellm/xai/grok-4.5` → company `xAI`, route `LiteLLM`. `glm-5p2` → company
`Zhipu AI`, route unknown (hermes' LiteLLM routing lives in `~/.hermes/config.yaml`, invisible to
parley — hermes-1).

### 5. Command surface — agreed except scope labels (see VC-2)

```
parley roster show   [--scope <deck|machine>] [--all] [--json] [--explain AGENT]
parley roster set    AGENT --scope <deck|machine> [--adapter A] [--state active|inactive]
                     [--model M] [--effort E] [--speed S] [--dry-run] [--yes]
parley roster sync   [--dir DIR] [--dry-run] [--yes]        # machine -> deck, ONE direction
```

Unanimous safety properties: preview is the default, `--yes` applies, `--dry-run` is
side-effect-free, writes are validated and atomic, and **`--yes` alone is refused when the change
alters membership** — a second explicit confirmation is required. Sync never copies deck values
upward. `roster init` becomes a deprecated alias. `--scope session` survives as a hidden
compatibility alias.

`roster show --all` reveals IDs mapped in config but absent from the deck roster, clearly marked —
this is what would have made the `opencode` situation visible on day one (kimi-1).

### 6. Session means an immutable run snapshot — unanimous

Not a third mutable scope. At run creation, write a secret-free roster snapshot plus a
`roster_revision` into run state; every later phase of that run uses it. `sessions inspect` reports
`stale-snapshot` when the deck roster has moved since. The single workflow becomes
**machine update → `roster sync` → new run** (codex-1).

### 7. The model-argv fix is IN SCOPE for this idea — unanimous

Freezing a contract whose headline cell is defined as "effective" requires an effective value to
exist. Mechanism (hermes-1's synthesis, adopted): built-in `HeadlessArgs` carry `{model}` and
`{effort}` placeholders resolved by the existing substitution path in `runner.go:1101-1103`, plus a
**legacy normalizer** for deck overrides that hardcode model literals in `headless_args`.

`codex` and `kimi` both accept `-m/--model` — codex-1 `PRIMARY` via `--help`, independently
re-verified by the facilitator (`codex exec --help` → `-m, --model <MODEL>`; `kimi --help` →
`-m, --model <model>`). So injection is the fix for them too, and `MODEL=unknown` is a transitional
state rather than a permanent one.

### 8. Skill/CLI boundary — unanimous

The skill **invokes `parley roster show` and reproduces its output**. It never parses §2, TOML, or
`agents list` to build a roster, and never documents a second table format. `SKILL.md`'s roster
section becomes a pointer plus the three verbs.

### 9. Migration — agreed in principle

Measured (`PRIMARY`, facilitator): 40 decks, **nine distinct §2 rosters**; 17 with no §2 roster at
all; 17 still naming retired `antigravity-1`; 3 naming `gemini-1`; 1 naming `agy-1`; `ai_prezz` has
no `hermes-1` row. This is the user's "hermes works in some sessions, not others", measured.

Agreed: existing decks keep working; `legacy-roster` in `STATUS` marks a deck whose roster predates
the contract; `roster sync --dry-run` is the documented remediation. **A fleet-wide migration is
NOT performed by this idea.**

### 10. §2 stops being the hand-edited membership store — agreed, but deferred

All four participants converged, two of them by reversing a round-1 position: §2 is the store that
drifted nine ways precisely because the protocol instructs humans to edit it by hand.
`parley-deck/agents.toml` should become the deck authority with §2 a generated view.

**This is a protocol change and is explicitly NOT ratified here.** It requires its own
`meta-protocol-change-*` idea per §7. This idea ships the CLI/skill standardization; §2 remains
authoritative for membership until that idea lands. Recorded as follow-up F1.

## Verdict conflicts

### VC-1 — Does the canonical table carry `SOURCE`? OPEN

- **kimi-1: yes** (12 columns). The winning layer is where you go to edit; naming it in the row
  makes drift actionable at a glance.
- **codex-1, hermes-1, claude-1: no.** Per-field provenance is too wide for the table and belongs
  in `roster show --explain AGENT` plus the JSON `sources` object.

**3-to-1, and §15.3 forbids resolving it that way.** The substantive argument against is that
`SOURCE` can only name the winning layer for *one* field, so a single column silently privileges
`MODEL` and misleads about `EFFORT`, `SPEED`, and `AUTO`, whose winning layers may differ.
kimi-1 should say whether that defeats the proposal or whether a `MODEL`-scoped `SOURCE` is still
worth its width. **claude-1 and hermes-1 both held kimi-1's position in round 1 and changed it**, so
this is not majority-by-default; it is two documented reversals plus codex-1's original argument.

### VC-2 — `sync` semantics: rebase or additive pin? OPEN

- **codex-1, claude-1 — rebase.** Sync *removes* deck overrides that mask machine values, so the
  deck keeps inheriting. A synced deck stays current.
- **hermes-1, kimi-1 — additive, source-aware pin.** Sync *writes* machine-sourced values into
  `parley-deck/agents.toml`, never touching a field whose winning source is already the deck. A
  synced deck becomes self-contained and reproducible.

**2-to-2, and the failure modes are opposite:** rebase risks a deck silently changing when the
machine default changes; additive-pin risks a deck going stale the day after sync. Both are real.

The unexamined question that would settle it: **does the audit trail require a deck to be
reproducible from its own files?** If a run's roster must be reconstructable from the repository
alone, additive-pin wins and rebase is inadmissible. If the run snapshot (decision 6) already
guarantees reproducibility, rebase wins and pinning is redundant staleness.

#### VC-2 measurement — `PRIMARY` (claude-1), with an ownership caveat

The drafter ran the check its own text asked for.

`runmanifest.Manifest` (`internal/runmanifest/manifest.go:28-45`) captures:
`schema_version`, `run_id`, `workspace_root`, `idea_slug`, `task`, `mode`, `transport`, `status`,
`phase`, `idea_status`, `current_round`, `active_steps`, `last_action_at`, `next_actions`,
**`participants []string`**, `created_at`, `updated_at`. `Step` carries `agent_id` and
`artifact_path`.

**The snapshot records participant IDs and nothing else.** No adapter, no model, no effort, no
speed, no launch plan. A completed run therefore does **not** record what any agent actually ran.

**What this does to VC-2 — stated, not verdicted.** codex-1's rebase position rests on the run
snapshot guaranteeing reproducibility. Today it does not. But decision 6 of this consensus —
unanimous — *adds* exactly that snapshot. So the two are coupled: **rebase is only admissible if
decision 6 ships in the same change**, and if decision 6 slips, additive-pin is the only option
that keeps a run reconstructable.

**§15.1 caveat: the drafter owns this measurement and is in the rebase camp**, so this cannot be
the resolving verdict. It is offered as evidence with its command and locators so a non-owner can
reproduce or refute it. **Resolution path for signoff:** a non-owner should verdict the
measurement, then say whether the coupling argument holds or whether reproducibility must not
depend on an unshipped feature.

### VC-3 — scope labels: `local|global` or `deck|machine`? OPEN, low stakes

codex-1 proposed `local|global`; hermes-1 and kimi-1 propose `deck|machine`, on the ground that
`local` is ambiguous between machine-local and project-local. claude-1 raised a related question
codex-1 did not answer: does `--scope local` write the committed `parley-deck/agents.toml` or the
gitignored `agents.local.toml`? **The answer must be the committed file** — an invisible gitignored
change is how a deck silently diverges from its own repository. Recorded so the decision is
explicit rather than incidental.

## Drafter position changes

claude-1 is facilitator, participant and drafter. Required by §15.5.

| # | Prior position | Source | New position | Why |
|---|---|---|---|---|
| 1 | `MODEL-COMPANY` is unknowable for gateway-routed models; print `unknown` rather than maintain a mapping | `round-01/claude-1.md` | Adopted codex-1's peel-the-route derivation | The objection assumed a flat string; peeling `litellm/` makes company unambiguous |
| 2 | Add a `SOURCE` column naming the winning layer | `round-01/claude-1.md` | Withdrawn in favour of `--explain` | codex-1's argument that provenance is per-field. Changed against the only participant who agreed with me |
| 3 | The model-argv fix might belong in a separate idea | `round-01/claude-1.md` | It must be in scope here | A frozen "effective" contract with no effective value would ship `unknown` for three adapters on day one |
| 4 | `roster sync --from machine --to deck` as a copy | `round-01/claude-1.md` | codex-1's rebase | Answers my own round-1 risk that a copy flattens deliberate pins. Now contested — see VC-2 |

Four changes, all four forced by another participant.

## Comparison & blind spots

**Correlated-agreement caveat (§15.6).** Four related models converged fast on the column contract.
Mitigating: three participants independently read `discover.go`/`runtime.go`/`runner.go` and cited
different line ranges, and each of the four reversed at least one of its own round-1 positions on
evidence — hermes-1 reversed six, kimi-1 three, codex-1 three, claude-1 four.

**What would have to be true for this to be wrong.** The whole package assumes the effective launch
plan is knowable *before* launch. If any adapter resolves its model at runtime from a source parley
cannot inspect (hermes' LiteLLM gateway config is already one such case), then `MODEL` is
`unknown` for that adapter permanently, and a contract that promises "effective" over-promises.
Nobody tested an adapter whose model is resolved server-side.

**Where nominally independent findings are one family.** The `MODEL` drift, the `EFFORT` absence,
and the `DISPLAY-NAME` contradiction are the same defect at three altitudes: the spec stores a value
in more than one place and nothing reconciles them. They were found separately by three
participants and should be fixed as one change, not three.

**Unmeasured.** VC-2's deciding question (what the run snapshot captures), and whether `SPEED` has
the same declared/effective defect as `MODEL` and `EFFORT` — nobody checked `SPEED` at all.

## Signoffs

<!-- Each participant appends its own block. Do not edit another participant's block. -->
