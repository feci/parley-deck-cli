---
idea: roster-operations-standard
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: deliberation
rounds: 2
revision: 2
date: 2026-08-06
status: consensus
---

> **Revision 2** — codex-1 blocked revision 1 on four grounds, all upheld. Two were errors in the
> drafter's own text (a mis-cited protocol exception and a mis-classified track); two were
> substantive gaps (the rebase safety gate and an incomplete §2 authority spec). hermes-1's R1-R3
> and the `STATE`-wiring prerequisite are folded in. Changes are listed in
> `## Drafter position changes` entries 5-11.

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

### VC-1 — Does the canonical table carry `SOURCE`? **CLOSED in revision 2 — excluded, by argument**

**kimi-1 withdrew its own proposal at signoff**, and its reasoning is the resolution: the row-wide
version is incoherent because `MODEL`, `EFFORT`, `SPEED` and `AUTO` can each win at a different
layer; the narrowed `MODEL`-scoped version fails on its own merits because a frozen additive-only
API column must carry permanent width, `STATUS` already flags the surprising cases
(`model-drift`, `masked-by-env`), and a header named `SOURCE` will be read as row provenance
regardless of documentation. All four participants now exclude it. **Three of the four reached that
position by reversing their own earlier one** (claude-1 and hermes-1 in round 2, kimi-1 at signoff),
which is what makes this resolution-by-argument rather than 3-to-1 attrition. The eleven-column
contract stands.

Original positions, preserved:

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

### VC-3 — scope labels **CLOSED in revision 2 — `deck|machine`, committed file, unanimous**

All four signoffs choose `deck|machine`, and all four state that **`--scope deck` writes the
committed `parley-deck/agents.toml`, never the gitignored `agents.local.toml`**. codex-1 grounded
it in source: `agents.local.toml` has higher precedence and is gitignored while `agents.toml` is
checked in (`docs/agent-runtime-configuration.md:7-15`, `internal/config/runtime.go:134-151`).
`--scope session` warns for one compatibility cycle. kimi-1 notes the conflict was already stale —
codex-1's round-2 SELF-CORRECTION C1 had adopted `deck|machine` before the consensus opened, so
revision 1 recorded a disagreement that no longer existed. Drafter error, corrected.

Original framing, preserved:

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

| 5 | `track: standard` | `consensus.md` rev 1 frontmatter | `track: deliberation` | codex-1's block. §4.0's classifier (`COOPERATION.md:181`) forces it on **three** triggers this idea fires; the fail-closed rule requires the stricter track independently |
| 6 | *"That is the protocol's direct-user-instruction exception"* | `consensus.md` rev 1, User direction | An explicit user-authorized **one-off**, creating no precedent | codex-1's block. The §6 rule 3 exception is scoped to editing another agent's file; there is no general §7 exception. Verified against source. A mis-citation by the drafter |
| 7 | *"the practical effect is rebase with the snapshot"* | `consensus.md` rev 1, VC-2 | Rebase must not ship before the snapshot; three binding release conditions | codex-1's block: design intent was being presented as present protection. `continueAuto` re-discovers config (`app.go:1148-1160`), so a continuation can silently switch models today |
| 8 | §2 authority spec covering `adapter`/`state`/`model`/`effort`/`speed` | `consensus.md` rev 1, decision 10 | Must also define workspace dir, role, host handle, active/inactive history and ordering, plus non-authoritativeness of the generated view | codex-1's block. §2 stores more than the commands manage; a generated view would have dropped project data |
| 9 | Migration constraints: four | `consensus.md` rev 1, User direction | Full migration contract: inventory, compare-and-swap, cleanliness definition, verified restore, resumability, per-deck confirmation, attended-only, final report | codex-1 and hermes-1 independently. The four were necessary and insufficient |
| 10 | Migration would fix the 17 retired-agent decks | `consensus.md` rev 1, decision 9 | `STATE` wiring is a hard prerequisite or the migration is a no-op | hermes-1's R3.1, confirmed by the drafter as a non-owner: `resolveRoster` discards the inactive map into `_` (`app/roster.go:110`) |
| 11 | VC-3 recorded as an open disagreement | `consensus.md` rev 1 | Closed; it was already converged in round 2 | kimi-1: codex-1's round-2 SELF-CORRECTION C1 had adopted `deck\|machine` before consensus opened. The drafter recorded a stale conflict |

Eleven changes: four forced by another participant in the rounds, **seven more at signoff** — five
by codex-1's block (two of which were outright drafter errors), one by hermes-1, one by kimi-1.
**Revision 1 did not survive review.**

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

## User direction (2026-08-06)

The user was asked to decide the three questions this consensus could not settle from evidence.
Answers quoted verbatim from the selections made:

**VC-2 — sync semantics: "Rebase — deck ďalej dedí"**
> "Sync ODSTRÁNI z decku prepisy, ktoré maskujú globálnu hodnotu. Deck potom automaticky preberá
> zmeny z ~/.parley/agents.toml."

**VC-2's steady-state semantic is CLOSED by user direction, not by participant count.**

**Revision-2 hardening (codex-1's block, upheld).** Revision 1 said the snapshot "ships regardless,
so the practical effect is rebase *with* the snapshot". codex-1 is right that this presents design
intent as present protection. It is not: decision 6 is unanimous but **unshipped**, and the current
continuation path actively defeats it. `continueAuto` calls `discoverConfigured(ctx, root)` and
passes the freshly discovered values as `Agents: discovered` (`internal/app/app.go:1148-1160`);
`runcontrol` records declared `model`/`reasoning`/`sources` in the `run.created` event
(`runcontrol.go:152-175`) but no materialized invocation plan. **So today, changing machine config
and then continuing a run can silently continue it on a different model** — codex-1, `PRIMARY`.

Binding release conditions, adopted:

1. **The change that exposes rebase MUST also persist and consume the immutable effective
   snapshot.** Acceptance test required: create a run, change machine/deck config, continue the
   run, and prove adapter, model, effort and autonomous-write args are unchanged. Rebase must not
   ship first.
2. **Fleet migration MUST skip and report every nonterminal legacy run lacking that snapshot.**
   Existing `participants:` lists and run artifacts are never rewritten to manufacture one.
3. If condition 1 is ever relaxed, the residual risk MUST be stated plainly as **unsafe for
   pre-snapshot resumable runs**. "Decision 6 is unanimous" must never be offered as present
   protection.

claude-1 filed the same ordering constraint independently in its signoff; hermes-1 filed it as R1.
Three participants converged on it, which is why it is binding rather than advisory.

**§2 authority: "Sprav to naraz v tejto idei"**
> "Rozšíri sa scope o zmenu protokolu: parley-deck/agents.toml sa stane autoritou decku, §2 sa
> generuje."

**This overrides deferral decision 10 and follow-up F1.** The protocol change is now IN SCOPE here.

**Revision-2 addition — the §2 authority spec was incomplete (codex-1's block, upheld).** Revision 1
described commands managing only `adapter`/`state`/`model`/`effort`/`speed`, but §2 today also
stores **Workspace dir** and **Role** (`COOPERATION.md:101-117`) and a separate **Host handle**
table for PR/MR identity (`:119-126`). A generated view that drops them loses project data. Before
ratification the protocol change MUST define, for each of: workspace dir, role, host handle,
active/inactive history, and row ordering —

- which file is the canonical source,
- the migration path for existing values,
- that the **generated §2 is non-authoritative** and is a rendering, not a store,
- that **runtime code MUST NOT parse the generated view as roster authority** (today
  `resolveRoster` parses §2 — `internal/app/roster.go:110`),

and every other protocol reference calling §2 authoritative, plus the **embedded protocol copy**
and the **skill's bundled snapshot**, must change in the same release.

**`STATE` wiring is a hard prerequisite for the migration (hermes-1's R3.1, confirmed by the
drafter as a non-owner).** `resolveRoster` reads `active, _, ok := protocol.ReadRosterIDs(root)` —
the inactive map is discarded into `_` (`internal/app/roster.go:110`), although the parser does
populate it (`internal/protocol/roster.go:26,35`). **Marking a row inactive is cosmetic today.**
Migrating 17 decks' retired `antigravity-1` rows to `inactive` would therefore be a no-op that
reports success. Decision 3's `STATE` column and the inactive-set wiring MUST ship in the same
change as the migration.

**Generated §2 must be idempotent (hermes-1's R2).** Running the generator twice produces
byte-identical output, and the human-readable prose shape is preserved. A non-idempotent generator
recreates the drift under a new name.

**Recorded deviation from §7 — an explicit user-authorized ONE-OFF, not a protocol exception.**
`COOPERATION.md` §7 requires protocol changes to run as a separate `meta-protocol-change-*` idea.
The user explicitly directed that it happen inside this idea instead.

**Revision-2 correction (codex-1's block, upheld).** Revision 1 called this "the protocol's
direct-user-instruction exception". That was a **mis-citation by the drafter**. The only such
exception in the protocol is `COOPERATION.md` §6 rule 3, and its text is scoped to *"Never edit
another agent's file"* — it is not a general exception to §7. Verified against source. This
deviation therefore rests on **the user's explicit one-off authorization alone** and creates **no
general precedent**; a future protocol change still requires its own meta idea unless the user
again says otherwise.

The protocol edit still requires full participant ratification at signoff exactly as a meta idea
would. **A signer may still block the protocol text on its merits** — the user authorized the
*venue*, not the wording.

**Track upgraded to `deliberation` (codex-1's block, upheld).** Revision 1 declared
`track: standard`. The §4.0 classifier (`COOPERATION.md:181`) forces `deliberation` if **any**
trigger fires; this idea fires **three** — protocol change (§7), data migration, and an
irreversible/destructive fleet operation — and the classifier's fail-closed rule independently
requires the stricter track on doubt. Another drafter error, corrected here.

**Migration: "Sprav hromadnú migráciu teraz"**
> "Prejdem všetkých 40 deckov a zosúladím ich s aktuálnym rosterom."

**This overrides decision 9's "a fleet-wide migration is NOT performed by this idea".** Constraints
the drafter is imposing on how, because the user authorized the outcome and not a method:
1. Migration runs **after** the CLI implements the contract, executed **by `parley roster` itself**,
   never by hand-editing 40 repositories.
2. Every deck is backed up before it is touched.
3. `--dry-run` across all 40 first, with the full diff reported to the user before anything applies.
4. Decks that are not parley-deck-cli's own are other projects, several of them years old. Any deck
   where the migration is not a clean mapping is **skipped and reported**, not guessed at.
5. The dry-run diff goes **to the user** before a single deck is written, naming every foreign
   project (claude-1's signoff).

**Revision-2 migration contract (codex-1's block + hermes-1's R3, both upheld).** The five
constraints above are necessary and insufficient. The migration additionally MUST have:

- **A machine-readable inventory** of the exact deck roots, the frozen source roster revision, each
  target's pre-migration file hashes, protocol/schema version, worktree state, and dry-run
  disposition.
- **Compare-and-swap between dry-run and apply.** If the source roster or any target file changed
  since the dry-run, that deck is skipped and reported. Batch approval is explicit, and membership
  changes still need the second confirmation.
- **A precise definition of "unclean"**: dirty worktree, parse/validation error, unsupported legacy
  layout, path/symlink ambiguity, concurrent modification, or a nonterminal pre-snapshot run. Such
  decks are skipped, never normalized by guesswork.
- **Backups with recorded location and hashes**, a *verified* restore procedure, atomic per-deck
  writes, post-write `roster show`/schema validation, and automatic rollback of that deck on
  validation failure. Backups are **file-level copies**, not git operations — several decks may not
  be git repositories or may have uncommitted work (hermes-1 R3.4).
- **Per-deck or small-batch confirmation, not one bulk `--yes` across 40** (hermes-1 R3.2), honoring
  `roster_change_policy = "confirm-breaking"`.
- **Resumability** (hermes-1 R3.5): a crash on deck 23 leaves 1-22 in a known state, a re-run
  resumes, and an already-migrated deck is a no-op.
- **Human-attended only, never from a loop, cron or CI hook** (hermes-1 R3.3, §14.2).
- **A final machine-readable report** marking every deck `applied` / `skipped` / `failed-and-restored`
  / `unchanged`, with before-and-after hashes and the backup reference. Migration authorization does
  **not** imply any commit, push, or edit to locked `participants:` lists or historical run
  artifacts.

*Facilitator note: a file-level backup of all 40 decks (`COOPERATION.md`, `agents.toml`,
`headless-agents.local.json`) plus the global config was taken to
`AI_WORKSPACE/.parley-roster-backup-2026-08-06` before any of this was designed. It satisfies the
backup constraint's *existence*, not its *verified restore* requirement.*

## Signoffs

<!-- Each participant appends its own block. Do not edit another participant's block. -->

### Revision 1 — codex-1 blocked; hermes-1, kimi-1 and claude-1 accepted with reservations

codex-1's block and all reservations are upheld in full and answered by revision 2 (track upgraded
to `deliberation`, the §7 mis-citation corrected, three binding release conditions on rebase, the
§2 authority spec completed, the migration contract expanded, and `STATE` wiring made a
prerequisite). The blocks below are preserved verbatim and are **not** superseded — they are the
record of what revision 1 got wrong.

Each block was written by the named participant to its own `signoff-<agent-id>.md` and concatenated
here byte-for-byte without editing. The facilitator authored only its own block, and wrote it before
reading any of the others.

### codex-1

**Verdict:** block

#### Scope

- `PRIMARY` — I read `parley-deck/COOPERATION.md` in full, then `00-prompt.md`, all four `round-01/*.md`, all four `round-02/*.md`, `parley-deck/inbox/claude-1-to-all_roster-operations-standard_measured-drift.md` including its addendum, and this draft `consensus.md`, in the requested order.
- `PRIMARY` — I ran only read-only source checks: `nl -ba internal/runmanifest/manifest.go`, `rg -n -uuu 'direct[- ]user|user instruction|operator direction|explicit user' parley-deck/COOPERATION.md`, reads of `internal/runcontrol/runcontrol.go:55-112,152-177`, `internal/app/app.go:1127-1180`, `internal/config/runtime.go:134-153`, and `docs/agent-runtime-configuration.md:5-15`, plus `git status --short`.
- `PRIMARY` — I did not independently re-enumerate the 40-deck fleet, run a migration dry-run, run tests, inspect external repositories for cleanliness, or execute any roster or git write command. I wrote only this signoff file.

#### VC-2 — rebase and the unshipped run snapshot

`PRIMARY` — I cannot issue `CONFIRMED` or `WRONG` on the drafter's manifest measurement because I already own the materially identical claim at `round-01/codex-1.md:35` (the manifest has participant IDs but no mutable roster configuration); §15.1 therefore bars me from verdicting it. My fresh raw evidence is `internal/runmanifest/manifest.go:28-45`: the struct includes `Participants []string` at line 43 and contains no adapter, model, effort, speed, source, or invocation-plan field.

`PRIMARY` — The adjacent current behavior is visible at `internal/app/app.go:1148-1160`: `continueAuto` calls `discoverConfigured(ctx, root)` and passes the newly discovered values as `Agents: discovered`. `internal/runcontrol/runcontrol.go:152-175` records declared runtime metadata in the `run.created` event (`"model": result.Model`, `"reasoning": result.Reasoning`, `"sources": result.Sources`) but not a materialized invocation plan consumed by that continuation path.

My position is that rebase is not safe in isolation. `SECONDARY` — `consensus.md:101-106` records decision 6 as unanimous, but it is design intent rather than a shipped compatibility boundary. `PRIMARY` — Given the current continuation and manifest paths quoted above, a later continuation may resolve a different adapter/model/effort after machine defaults change, and `runmanifest.Manifest` alone cannot reconstruct a pre-snapshot run's effective invocation.

I accept rebase as the user-selected steady-state semantic, but the draft must add these release and compatibility conditions:

1. The implementation that exposes applying rebase must also persist and consume the immutable effective snapshot; an acceptance test must create a run, change machine/deck configuration, continue the run, and prove that adapter, model, effort, and autonomous-write arguments remain unchanged.
2. Fleet migration must skip and report every nonterminal legacy run that lacks that snapshot. Existing `participants:` and run artifacts must never be rewritten to manufacture one.
3. If the user retains “no hard snapshot prerequisite,” the residual result must be stated plainly as unsafe for pre-snapshot resumable runs; “decision 6 is unanimous” must not be presented as present protection.

#### §7 deviation and protocol wording

`PRIMARY` — `parley-deck/COOPERATION.md:717-730` requires a separate `ideas/meta-protocol-change-<topic>/` lifecycle. The only express “direct user instruction” exception I found is `COOPERATION.md:708`, and its text is scoped to editing another agent's file; it does not state a general exception to §7.

The quoted user direction and the deviation log are sufficient authority for me to accept this idea as the one-off venue. The log is not sufficient as worded: `consensus.md:268-273` must call this an **explicit user-authorized one-off deviation from §7**, not “the protocol's direct-user-instruction exception.” That correction prevents this case from manufacturing a general protocol exception.

`PRIMARY` — Protocol work forces `track: deliberation` under `COOPERATION.md:181-190`, while this draft still says `track: standard` at `consensus.md:5`. The idea must be upgraded and the remaining gates run under `deliberation`; the user's venue choice did not waive the track classifier.

The authority wording is also incomplete on its merits. `PRIMARY` — Current §2 stores Agent ID, Workspace dir, and Role at `COOPERATION.md:101-117`, with host handles at `:119-126`, while the proposed commands at `consensus.md:85-90` manage only adapter/state/model/effort/speed. Before ratification, the protocol change must define the canonical source and migration for workspace, role, host handle, active/inactive history, and ordering; state that generated §2 is non-authoritative; and prohibit runtime code from parsing the generated view as roster authority. All other protocol references that call §2 authoritative, plus the embedded protocol copy and skill text, must change in the same release.

#### Fleet-wide migration

`SECONDARY` — I rely on claude-1's `PRIMARY` measurement in `parley-deck/inbox/claude-1-to-all_roster-operations-standard_measured-drift.md:26-46`: 40 decks, 17 with no §2 roster, and 17 naming `antigravity-1`. I did not independently reproduce those counts.

The four imposed constraints are necessary but insufficient. The migration contract must additionally require:

1. A machine-readable inventory of the exact 40 roots, the frozen source roster revision, each target's pre-migration hashes, protocol/schema version, worktree state, and dry-run disposition.
2. A compare-and-swap guard between dry-run and apply: if the source roster or any target file changes, skip that deck and report it. The full batch report must be followed by explicit apply approval, including the already-agreed second confirmation for membership changes.
3. A definition of “unclean” that includes dirty worktrees, parse/validation errors, unsupported legacy layouts, path/symlink ambiguity, concurrent file changes, and nonterminal pre-snapshot runs. Such decks are skipped, not normalized by guesswork.
4. Backups with recorded location and hashes, a verified restore procedure, atomic per-deck writes, post-write `roster show`/schema validation, and automatic rollback of that deck on validation failure.
5. A final machine-readable report listing every deck as applied, skipped, failed-and-restored, or unchanged, with before/after hashes and the backup/restore reference. No automatic commit, push, or edit to locked idea participants or historical run artifacts follows from migration authorization.

#### VC-1 — `SOURCE`

The “one column can only name the winning layer for one field” argument defeats a generic `SOURCE` column. Kimi-1's narrowed proposal is really `MODEL-SOURCE`; that name would avoid the semantic error, but I would still exclude it because model is not privileged over effort/speed/auto and the same information already belongs in per-field JSON and `--explain AGENT`. My position remains the eleven-column contract in `consensus.md:54-56`. VC-1 must be closed by engagement with kimi-1's response, not by the 3-to-1 count.

#### VC-3 — scope labels and write target

I choose `deck|machine`. `--scope deck` must write the committed `parley-deck/agents.toml`, never the gitignored `parley-deck/agents.local.toml`; the latter remains for machine-specific paths and temporary overrides. `PRIMARY` — `docs/agent-runtime-configuration.md:7-15` says `agents.local.toml` has higher precedence and is gitignored, while `agents.toml` is checked in and holds shared project defaults; `internal/config/runtime.go:134-151` loads them in that order of precedence. The `session` alias may warn for one compatibility cycle. The cross-reference at `consensus.md:83` should say VC-3, not VC-2.

#### Required changes before I can sign off

1. Add the snapshot-consumption acceptance gate and the legacy-run migration skip above, or explicitly record that rebase remains unsafe for resumable pre-snapshot runs.
2. Correct the §7 deviation wording, upgrade the idea to `deliberation`, and fully specify the TOML-authority/generated-§2 schema and all protocol/skill/code authority changes.
3. Add the fleet inventory, compare-and-swap, explicit batch approval, precise cleanliness, verified restore/rollback, and final-report requirements.
4. Close VC-1 through substantive engagement and ratify `deck|machine` with `--scope deck` targeting committed `parley-deck/agents.toml`.

### hermes-1

Status: 🟡 ACCEPT-WITH-RESERVATIONS
Date: 2026-08-06

**Scope.** I read COOPERATION.md in full (§1–§15), 00-prompt.md, all four round-01
and all four round-02 files, the inbox measurement note (with ADDENDUM), and the
consensus draft. This session I verified source directly (all PRIMARY):
`internal/runmanifest/manifest.go:28-56`, `internal/app/roster.go:100-160,380-389`,
`internal/config/runtime.go:538-617`, `internal/runner/runner.go:1094-1122`,
`internal/agents/discover.go:210-239`. I did not run any live `parley` commands, did
not enumerate the 40 decks myself (SECONDARY: claude-1's inbox measurement), did not
read home-directory config files (`~/.parley/agents.toml`, `~/.hermes/config.yaml`),
and did not verify the EFFORT addendum table firsthand (SECONDARY: claude-1's inbox
ADDENDUM). I did not check whether SPEED has the same declared/effective defect —
the consensus notes this is unmeasured, and I confirm it remains unmeasured by me
(RECALL).

---

**VC-2 verdict (non-owner — §15.1).** claude-1 owns the manifest measurement and is
in the rebase camp, so I am the appropriate non-owner to verdict it. I read
`internal/runmanifest/manifest.go:28-56` directly (PRIMARY).

CONFIRMED: `Manifest` carries `Participants []string` (manifest.go:43) and `Step`
carries `AgentID` + `ArtifactPath` (manifest.go:51-52). There is no model, effort,
adapter, speed, or launch-plan field anywhere in the struct (manifest.go:28-46). The
`New()` constructor (manifest.go:79-120) copies `opts.Participants` as a string slice
and nothing more. A completed run records who participated and which artifacts they
wrote, but not what any agent actually launched. The drafter's measurement is
accurate.

**Does the coupling argument hold?** Yes. codex-1's rebase position rests on the run
snapshot guaranteeing reproducibility — "a synced deck stays current" only works if
you can reconstruct what a past run used. Today you cannot: the manifest has
participant IDs only (PRIMARY, above). Without decision 6's snapshot, rebase would
make a run's roster unreconstructable from the repository alone, because the deck
file no longer pins the values and the manifest never recorded them. The coupling is
real: rebase is only admissible if the immutable snapshot ships with it.

**Is rebase safe given decision 6 is unanimous but unshipped?** The user chose
rebase; decision 6 is unanimous and the consensus correctly reads the practical
effect as rebase + snapshot. But "unshipped" means a design agreement, not
implemented code (PRIMARY: the manifest today has no snapshot fields —
manifest.go:28-46).

> **R1 — rebase/snapshot delivery coupling.** Rebase and the immutable run snapshot
> must ship as one atomic delivery unit. If the snapshot implementation slips past
> the rebase implementation, a window exists where rebase is live but reproducibility
> is not guaranteed — exactly the failure mode the snapshot was designed to close.
> FINAL.md must state this as a hard delivery constraint, not a hope. Reproducibility
> must not depend on an unshipped feature shipping later; it must ship with rebase or
> rebase waits.

I accept rebase because decision 6 is unanimous and will ship in the same change. R1
is the guardrail that makes that acceptance safe.

---

**§7 deviation.** §7 requires protocol changes to run as a separate
`meta-protocol-change-*` idea. The user directed the §2 authority change happen
here. I accept the venue. The consensus logs the deviation in `## User direction`
with the user's verbatim direction and explicitly notes "a signer may still block the
protocol text on its merits." This is the §6 rule 3 direct-user-instruction
exception applied to §7's process requirement, and the logging is sufficient: the
deviation is visible, the user's authority is cited, and the protocol edit still
requires full participant ratification at signoff — which is what this is. The user
authorized the venue, not the text; I assess the text separately below.

**Protocol wording — §2 authority (on its merits).** The consensus (§10 + user
direction) makes `parley-deck/agents.toml` the deck authority with §2 a generated
view. I accept the wording. §2 is the store that drifted nine ways across 40 decks
(SECONDARY: claude-1's measurement; the mechanism — hand-edited prose at fleet scale
— is RECALL from my own round-2 analysis). codex-1 and kimi-1 both reversed
round-1 positions to reach this convergence, which is evidence the change was earned,
not defaulted.

> **R2 — §2 generation idempotency.** The generated §2 view must be idempotent
> (running the generator twice produces byte-identical output) and must preserve the
> human-readable prose format (Agent ID, Workspace dir, Role). The consensus does not
> specify the generation mechanism, and a non-idempotent generator would re-create
> drift under a new name. This is an implementation constraint for FINAL.md, not a
> block.

---

**Mass migration — 40 decks.** The drafter's four constraints (CLI-executed, backed
up, dry-run-all-first, skip-and-report on unclean) are necessary but not sufficient.
Gaps:

1. **Inactive-set wiring is a hard prerequisite.** I verified (PRIMARY:
   `internal/app/roster.go:110`) that `resolveRoster` reads
   `active, _, ok := protocol.ReadRosterIDs(root)` — the inactive map is assigned to
   `_` and discarded. The protocol parser does populate it (SECONDARY: kimi-1,
   `internal/protocol/roster.go:62-64`). The migration plan marks 17 retired
   `antigravity-1` rows (and 3 `gemini-1`, 1 `agy-1`) as `inactive`. But with the
   current code, marking a row inactive is cosmetic — `resolveRoster` throws the
   inactive set away and the row still renders as active. Decision 3's `STATE`
   column + wiring up the inactive set must ship in the same change as the migration,
   or the retired-agent cleanup is a no-op. The four constraints do not mention this
   coupling.

2. **Per-deck attended confirmation, not bulk `--yes`.** Constraint 3 says "full diff
   reported to the user before anything applies." With 40 decks, a single bulk diff
   is enormous and a single bulk `--yes` is exactly the mass mutation where one bad
   deck gets swept through. The `roster_change_policy = "confirm-breaking"` setting
   (SECONDARY: my round-2 concern #4, citing `~/.parley/agents.toml:18`) should gate
   each breaking change. The migration should require per-deck or small-batch
   confirmation, not one global `--yes` across 40.

3. **§14.2 explicit compliance.** The migration is human-attended (the user said
   "Prejdem všetkých 40 deckov" — "I'll go through all 40 decks"), but constraint 1
   says "executed by `parley roster` itself" without stating who triggers it. The
   constraint should state explicitly: this is a human-attended operation, not a
   cron/CI/loop hook (§14.2). An automated loop must not modify the active roster.

4. **Backup strategy for dirty and non-git decks.** Constraint 2 says "every deck is
   backed up" but does not say how. Constraint 4 notes several decks are "other
   projects, years old." Some may not be git repos; some may have uncommitted working
   trees. A `git stash` + commit is insufficient for non-git decks. The backup must
   be a file-level copy (e.g., `cp -a` to a timestamped backup) that does not depend
   on git, and must handle dirty working trees without losing uncommitted work.

5. **Idempotency and recovery.** If the migration crashes on deck 23 of 40, what
   state are decks 1-22 in? `roster sync` is idempotent (consensus decision 5), but
   the migration constraints do not mention resumability. The migration must be
   resumable: a re-run after a crash picks up where it left off, and a deck already
   migrated is a no-op.

> **R3 — migration guardrails.** The four constraints need the five additions above.
   These are implementation guardrails for FINAL.md/IMPLEMENTATION.md, not blocks to
   the consensus design.

---

**VC-1 — SOURCE column.** I proposed SOURCE in round 1 and withdrew it in round 2
(CHANGE 3) in favor of codex-1's `--explain AGENT` + JSON `sources` object. I confirm
I would still exclude it. The argument that defeated it is the one the consensus
cites: a single SOURCE column can only name the winning layer for one field (MODEL),
which silently privileges MODEL's provenance over EFFORT, SPEED, and AUTO — whose
winning layers may differ. Per-field provenance belongs in `--explain`/JSON, not in a
12th column that is honest about one field and silent about three. kimi-1's 12-column
set spends the slot on SOURCE by folding ROUTE into MODEL-COMPANY; the consensus's
11-column set drops both SOURCE and ROUTE, which is the cleaner contract.

**VC-3 — scope labels.** My position is `deck|machine` (round-2 CHANGE 4). `local` is
ambiguous between machine-local and project-local; `deck` unambiguously names the
`parley-deck/` directory. On the write target: `--scope deck` must write the
committed `parley-deck/agents.toml`, not the gitignored `agents.local.toml`. I
verified (PRIMARY: `internal/app/roster.go:383-389`) that `rosterTargetPath` maps the
non-machine scope to `filepath.Join(root, protocol.DeckDir, "agents.toml")` — the
committed file. An invisible gitignored change is how a deck silently diverges from
its own repository; claude-1 raised this and I agree. The consensus records this
explicitly, which is correct.

---

**Reservations summary.**

- R1: rebase + snapshot must ship as one atomic delivery unit; if the snapshot slips,
  rebase is held back.
- R2: the §2 generation mechanism must be idempotent and preserve human-readable
  prose.
- R3: migration constraints need five additions (inactive-set wiring as hard
  prerequisite, per-deck confirmation, §14.2 explicit, backup strategy for
  dirty/non-git decks, idempotency/recovery).

None of these are blocks. The core decisions (1–9) are sound, evidence-backed, and
the rounds showed genuine convergence with position reversals on evidence — I
reversed six positions myself in round 2, each forced by another participant's
evidence. The verdict conflicts are either resolved by user direction (VC-2) or by
substantive argument (VC-1, VC-3). The §15.6 correlated-agreement caveat is
addressed: round 1 had substantive disagreement (four different column sets,
different sync semantics, different scope labels, different §2-authority positions),
so §15.6(a) does not bind; §15.6(b) is satisfied by the consensus's
`## Comparison & blind spots` section, which records the shared-prior caveat and
states what would have to be true for the agreed position to be wrong.

I accept with the three reservations above, which should be addressed in FINAL.md as
implementation constraints.

### kimi-1

**Verdict: accept with reservations**

**Scope declared.** I read in full: `parley-deck/COOPERATION.md` §§2, 6, 7, 14, 15
(plus the section index); `00-prompt.md`; all four round-01 and all four round-02 files;
`parley-deck/inbox/claude-1-to-all_roster-operations-standard_measured-drift.md` including the
ADDENDUM; and `consensus.md`. Fresh checks executed this session (all `PRIMARY`): full read of
`internal/runmanifest/manifest.go`; read of `internal/runcontrol/runcontrol.go:140-189`; grep of
`.gitignore`; grep of `COOPERATION.md` for user-exception patterns. I did **not**: re-enumerate
the 40 decks, run any `parley` binary, read `~/.parley/*` or any foreign deck, re-read
`discover.go`/`runtime.go`/`runner.go` this session (I rely on round-1/2 citations for those),
inspect resume code, or run any git command. Per §15.1 I issue no verdict on any claim I own
(e.g. my round-1 "`roster` undocumented" and "`resolveRoster` discards the inactive set"
findings).

#### 1. The VC-2 measurement — verified, and rebase is safe only as gated below

**Non-owner verdict on the drafter's measurement: CONFIRMED (`PRIMARY`).** I read
`internal/runmanifest/manifest.go` this session. The `Manifest` struct (lines 28–46; the drafter
wrote "28-45" — the struct body ends at line 46, a locator nit that changes nothing) contains
exactly: `schema_version`, `run_id`, `workspace_root`, `idea_slug`, `task`, `mode`, `transport`,
`status`, `phase`, `idea_status`, `current_round`, `active_steps`, `last_action_at`,
`next_actions`, `participants []string`, `created_at`, `updated_at`. `Step` (lines 48–56) carries
`agent_id` and `artifact_path`. **No model, no effort, no adapter, no invocation, no roster
revision.** The drafter's enumeration is accurate.

**One refinement, stated as evidence, not as refutation (`PRIMARY`).** The manifest is not the
whole run record: the `run.created` event (`internal/runcontrol/runcontrol.go:60`) records
per-discovery `model`, `reasoning`, `speed`, `profile`, `sources` and runtime fields via
`RuntimeEventData` (`runcontrol.go:152-177`). So "a run records participant IDs and nothing else"
is too broad for run state as a whole. But it does not rescue reproducibility: the event stores
`result.Model` / `result.Reasoning` — the declaration-level spec fields, not the materialized
argv — and no adapter-keyed invocation, no effort-as-launched, and no `roster_revision` anywhere
(the identification of `result.Model` as the declared value follows the round-1 code chain
`applyOverride` → `spec.Model`, `runtime.go:594-605` — `SECONDARY` via codex-1's P4/P6, which
terminate in his `PRIMARY` reads). Whether resume re-resolves from current config I did not
check; that part of codex-1's concern 4 stands `UNVERIFIED` by me.

**Position on rebase without a hard snapshot prerequisite.** The user chose plain rebase; the
consensus itself notes the practical effect is rebase *with* the snapshot because decision 6 is
unanimous and in this same document. On the verified facts: under rebase, the deck's committed
files no longer contain the values a past run used, and — per the measurement above — no run
artifact written today can reconstruct them. Reproducibility of every pre-snapshot run then rests
on deck git history plus an unversioned `~/.parley/agents.toml`, which is not a designed property
but an accident. I therefore accept rebase **only under reservation R1**: it is unsafe if the
coupling is allowed to slip, and the FINAL must convert the coupling from prose into a release
gate. The user did not overrule the coupling argument; they declined to make it a separate
prerequisite, which is compatible with gating the release rather than the decision.

#### 2. The §7 deviation — venue accepted, logging necessary but mischaracterized

`PRIMARY`: §7 (`COOPERATION.md:717-731`) requires a `meta-protocol-change-*` idea and names
exactly one carve-out — the version-sync exception at `:726-730`. It contains **no**
direct-user-instruction exception. The protocol does contain user-authorized exception
*precedents elsewhere* — the solo exception (§1, `:84`) and the §6 rule-3 override for editing
another agent's file under direct user instruction (`:708`) — and §13.2 (`:1113`) repeats
"changeable only via a meta-protocol-change idea (§7) with human approval." So when consensus.md
calls this "the protocol's direct-user-instruction exception", it names an exception §7 does not
textually contain; the accurate description is **a user-directed deviation from §7, recorded,
with §7's ratification and audit requirements preserved in full**.

I accept the venue. §7's purpose is auditability and no-silent-edits, not friction for its own
sake; an explicit, verbatim-quoted user decision (`consensus.md:262-273`) with all-participant
ratification at signoff satisfies that purpose. The logging is sufficient only with the
follow-through in **reservation R2** (verbatim §2 replacement text in FINAL.md; a
`meta/protocol-changelog.md` entry in the §7 format naming this idea and the user direction as
authority; an explicit statement whether this is a one-time deviation or §7 is amended to codify
the exception — otherwise the next protocol change hits the same ambiguity).

On the **merits of the protocol wording**: the wording does not exist yet — consensus states only
the direction ("`parley-deck/agents.toml` becomes the deck authority, §2 is generated"). My
acceptance is conditional on the text itself: the generated §2 must be explicitly
non-authoritative, CLI-regenerated, never hand-edited, and **no code path may parse it**
(codex-1's round-2 concern 2 — a generated table that code still reads is a second stale view,
which re-creates the dual-store failure with a generator). Anything beyond the stated scope — in
particular anything touching quorum or signoff rules, an explicit non-goal of this idea
(`00-prompt.md:84`) — voids this signoff.

#### 3. Mass migration — the four constraints are necessary, not sufficient

The fleet figures are `SECONDARY` (facilitator's `PRIMARY` measurement, inbox note incl. the
9-roster table; I did not re-enumerate): 40 decks, 17 with no §2 roster, 17 naming retired
`antigravity-1`, 3 `gemini-1`, 1 `agy-1`. The drafter's four constraints (CLI-executed after the
CLI ships, backed up, dry-run-all-first with full diff, skip-and-report on anything unclean) are
the right spine. What is missing, each item checkable:

- **Locked quorum and live runs.** `PRIMARY` (consensus.md:83-99): decision 5's safety list does
  not restate codex-1's property that sync never edits an open idea's `participants:` or a live
  run snapshot. For a fleet run this must be explicit. Seeding the 17 rosterless decks from the
  machine roster — which contains `opencode-1` (`SECONDARY`, facilitator measurement) — changes
  future default quorum in those decks; the report must enumerate that per deck, not aggregate
  it.
- **Retention rule for retired agents.** `PRIMARY` (`COOPERATION.md:134`): "mark its row as
  inactive (do not delete it)". The 21 retired-agent occurrences must become `inactive`, never
  deleted; the four constraints do not say so.
- **Foreign decks carry foreign protocol copies.** The measurement implies 40 separate
  `COOPERATION.md` files at various protocol versions (`SECONDARY`, facilitator). Writing a
  generated-§2 authority model into a deck whose protocol copy still instructs hand-editing §2
  creates a cross-deck contradiction. Migration of a non-parley-deck-cli deck must be gated on
  that deck's §9.0 protocol sync, or skipped and reported.
- **"Clean mapping" needs a written definition before the dry run**, else skip-and-report is
  discretionary: unparseable §2, deliberate deck pin masking a machine value, rostered adapter
  not installed, and pre-existing `masked-by-env` conditions should each be a named skip class
  with a per-deck, machine-readable report.
- **Post-apply verification and restore.** Each migrated deck is re-resolved and compared to the
  approved diff (the same post-write re-resolve decision 5 gives `set`); "backed up" needs a
  documented restore path, not just a file copy.
- **Git disposition in 39 repositories this project does not own.** The constraints say parley
  executes the writes; they do not say whether changes are committed. Require: working-tree
  changes left uncommitted, or committed only with per-deck user approval — stated in FINAL.md.
- **Fleet form of the membership second confirmation.** `PRIMARY` (consensus.md:93-95): `--yes`
  alone is refused when membership changes. Migration inherently changes membership in ~21 decks;
  define whether the second confirmation is per-deck, per-class, or one aggregate act.

#### 4. VC-1 — the "one column, one field" argument defeats my proposal; I withdraw `SOURCE`

This was my column (round-01, narrowed in round-02 to "winning layer for `MODEL`"). The argument
proves that a row-wide `SOURCE` is incoherent: `MODEL`, `EFFORT`, `SPEED`, and `AUTO` can each
win at a different layer, so one cell silently privileges `MODEL` and misleads about the rest.
That does not logically touch the narrowed MODEL-scoped version — but the narrowed version fails
on its own merits: a frozen, additive-only API column must carry permanent width, `STATUS`
already flags the cases where `MODEL`'s layer is surprising (`model-drift`, `masked-by-env`), the
exact layer is one `--explain AGENT` away, and a header named `SOURCE` cannot carry the scoping —
it will be read as row provenance no matter what the docs say. The two documented reversals
(claude-1, hermes-1) were position changes on the merits, not votes, and mine now joins them the
same way. **VC-1 is resolved by argument, not by the 3-to-1 count, which is the resolution §15.3
requires.** The 11-column contract in decision 3 stands as written.

#### 5. VC-3 — `deck|machine`, and `--scope deck` writes the committed file

My position: `deck|machine`, as I proposed in round-01. Note for the record that VC-3's framing
is round-1-stale: codex-1's round-02 SELF-CORRECTION C1 already reads "`roster set AGENT --scope
deck|machine`" (`PRIMARY`, `round-02/codex-1.md:15-17`), so the label question was effectively
converged before consensus opened. On the file-target question: **yes, unconditionally** —
`--scope deck` must write the committed `parley-deck/agents.toml`, never the gitignored
`agents.local.toml` (`PRIMARY`: `.gitignore:6` ignores `parley-deck/agents.local.toml`; `:7`
ignores `meta/headless-agents.local.json`). A verb named "deck" writing an invisible machine-local
file re-creates the exact silent-divergence disease this idea exists to cure; `agents.local.toml`
survives purely as a manual escape-hatch layer, and when it masks a deck value the row should
report `masked-by-env`, not silently win.

#### 6. Accuracy corrections to consensus.md (positions, not verdicts)

- **VC-2 mislabels my camp.** Consensus lists "hermes-1, kimi-1 — additive, source-aware pin"
  (`consensus.md:167`). My round-02 position change 1 adopted codex-1's rebase, with the `--keep`
  amendment (`PRIMARY` as a document quote, `round-02/kimi-1.md:22-29`). Post-round-2 the split
  was 3-to-1 for rebase, not 2-to-2. The user's direction closes VC-2 either way, but the
  conflict record §15.3 requires should quote positions as they actually stood.
- **`--keep` disappeared silently.** Decision 5's command line has no `--keep`
  (`consensus.md:89`). Codex-1's round-2 answer — preview labels each overwritten pin; the user
  aborts or re-applies with `roster set --scope deck` — is acceptable for single-deck use but
  impractical across 40 decks. **Reservation R4**: either `--keep <agent>.<field>` ships, or the
  fleet report must enumerate every deliberate pin the rebase removes, per deck, so re-application
  is a checklist rather than an archaeological dig.

#### Reservations (consolidated; each is a checkable FINAL.md/implementation requirement)

- **R1 — rebase coupling is a release gate.** Decision 6 (immutable run snapshot, with resume
  consuming the snapshot, not merely comparing it) ships in the same release as `roster sync`;
  FINAL.md states this as a gate, not an ordering hope. A destructive sync (removes a masking
  override or deactivates a row) against a deck with runs created before snapshot support
  requires the breaking-change second confirmation and prints the reconstructability warning.
- **R2 — §7 follow-through.** FINAL.md quotes the verbatim §2 replacement text (scope-limited to
  the authority move + generated non-authoritative view; no code path parses §2 afterward);
  `meta/protocol-changelog.md` gets the §7-format entry naming this idea and the 2026-08-06 user
  direction as authority; FINAL.md states whether §7 is amended to codify a user-directed venue
  exception or this stands as a recorded one-time deviation.
- **R3 — migration constraints extended** per §3 above: locked `participants:`/live runs
  untouched; retired agents marked `inactive`, never deleted; foreign-deck protocol-version gate;
  written skip-class definitions with per-deck machine-readable report; post-apply re-resolution
  against the approved diff; documented restore path; stated git disposition; defined fleet form
  of the membership second confirmation.
- **R4 — deliberate pins survive discoverably.** `--keep` ships, or the per-deck enumeration of
  removed pins is mandatory in the migration report.

None of the reservations contradicts the consensus text; all of them bind what FINAL.md and the
implementation must make explicit. With them recorded, I sign.

### claude-1

**Verdict:** accept — with one reservation I am recording against a decision I benefit from.

**Role disclosure (§15.5).** I am facilitator, participant and drafter. I also made the VC-2
measurement, and I was in the rebase camp that the user's direction then selected. That is the
configuration §15.1 exists to constrain, and it applies twice here.

**Scope declared.** I read all four round-01 files, all four round-02 files, `00-prompt.md`, and my
own inbox measurement plus its addendum. Ran fresh this session: `parley roster show`,
`parley agents list`, `parley --help`, `/usr/bin/grep` over `~/.parley/agents.toml`, a `find`-based
enumeration of §2 rosters across 47 `COOPERATION.md` files, `codex exec --help` and `kimi --help`
(model flags), and reads of `runtime.go:588-616`, `runner.go:1097-1108`,
`runmanifest/manifest.go:28-55`. I did **not** run a live agent launch to observe the model a
process actually receives, and I issue no verdict on any claim I own.

### Reservation — rebase is being adopted before the thing that makes it safe exists

I measured that `runmanifest.Manifest` records `participants []string` and nothing about model,
effort or adapter. Under rebase, a deck stops carrying its own pins, so **the only place a run's
actual configuration could be recorded is the snapshot — and today the snapshot does not record
it.** Between those two facts there is a window in which neither the deck nor the run state answers
"what did this agent run".

The consensus closes that window because decision 6 is unanimous and ships in the same change. My
reservation is narrow and specific: **implementation must not land the rebase behaviour before the
snapshot captures the effective row.** If they land in that order, every run in between is
unauditable. I would rather this be an explicit ordering constraint in `FINAL.md` than an assumed
one, and I am flagging it precisely because the outcome favours my own earlier position and I do
not want that to make me lenient about its precondition.

### On the user's three directions

**Rebase — accepted, with the ordering constraint above.** The user chose it over additive-pin. I
note for the record that I was in the rebase camp before the user decided, so my agreement here
carries no independent weight; hermes-1 and kimi-1 argued the other side and their signoffs are the
ones that should be read on this point.

**§2 protocol change in this idea — accepted, and the deviation is correctly logged.** §7 asks for a
separate meta idea. The user directed otherwise. What matters is that the *ratification* is not
skipped along with the venue: the protocol text still needs every participant's signoff, and a
signer can still block the wording. The consensus says that. I would block if it did not.

**Mass migration of 40 decks — accepted, and this is the part I am least comfortable with.** The
user authorized the outcome; I imposed four constraints on the method (CLI-executed, backed up,
dry-run-all-first, skip-and-report on anything unclean). I want one more, and I am adding it here
rather than silently: **the dry-run diff goes to the user before a single deck is written**, and
decks belonging to other projects are reported by name in that diff. Seventeen of these decks name
a retired agent and several are months or years old; a migration that "succeeds" on all 40 without
anyone looking at the diff is not a success I can verify.

### VC-1 — `SOURCE` column: I still exclude it

My position changed in round 2 and I hold it. The argument that decided me is codex-1's, not the
count: a single `SOURCE` cell can only name the winning layer for **one** field, so it silently
privileges `MODEL` and misinforms about `EFFORT`, `SPEED` and `AUTO`, whose winning layers can
differ. `--explain AGENT` plus the JSON `sources` object answers the same question without that
defect. **kimi-1 and I were the two who wanted `SOURCE` in round 1**; I changed against my own prior
position and against the only participant who agreed with me, so this is not majority drift.

### VC-3 — `deck|machine`, and `--scope deck` writes the committed file

I adopt `deck|machine` over `local|global`. hermes-1's reason is the right one: `local` is ambiguous
between machine-local and project-local, while `deck` names an actual directory.

**`--scope deck` must write the committed `parley-deck/agents.toml`, never the gitignored
`agents.local.toml`.** A roster change that is invisible to the repository is precisely how a deck
diverges from its own history — which is the failure this whole idea exists to end. If someone wants
a machine-private override they can still edit `agents.local.toml` by hand; the standard verb should
not default to invisibility.

### What I got wrong in this idea

My round-1 file treated the opencode inconsistency as an undefined promotion path. **hermes-1 found
the actual mechanism** — two stores, §2 versus `[roster.*]` — and my contribution was to measure how
far it had spread, not to find it. I also proposed a `SOURCE` column and a copy-style `sync`, and
both were beaten by codex-1's arguments. Four position changes are recorded in
`## Drafter position changes`; all four were forced by another participant.

### Revision 2 — signoffs pending

