---
agent: hermes-1
idea: composite-agent-naming-and-roster-reinit
round: 1
date: 2026-07-18
role: reinit UX — interactive discovery then per-agent model+effort selection; session vs machine scope defaults, non-interactive/CI path, failure modes
---

## Summary

I support all three components in principle but I reject the framing that the
composite name should BE the stable artifact ID. The single most important fact
in this codebase for the crux is that **the artifact filename is already NOT the
§2 roster ID** — `internal/runner/runner.go:349` writes
`filepath.Join(opts.Idea.Path, opts.RoundLabel, agent.ID+".md")` where `agent.ID`
is the **family ID** (`claude`, `codex`, `agy`, `hermes`) from `agents.Spec.ID`,
not the `claude-1` row in the §2 table. The §2 roster IDs (`claude-1` etc.) are
consulted **only** by roster-preset validation (`config/roster.go:ResolveRoster`
→ `protocol/roster.go:ReadRosterIDs`). So today there are *already* two
identities in flight — the §2 roster ID and the runtime artifact ID — and they
do not even agree. That pre-existing divorce is the thing to fix, and the
composite name is the right vehicle for the *display* layer, not the *identity*
layer.

My position on the crux: **stable base ID + composite display name** (option
(b) in the brief), where the base ID is a short stable kebab token (`claude-1`)
and the composite (`claude-opus48-max`) is a derived, regenerable display string
carried in config, the §2 table, the TUI, and digests — never in the artifact
path or the `agent:` frontmatter key. Renaming an agent mid-idea by making the
composite the ID would orphan every `round-NN/<id>.md`, every `### Signoff:
<id>`, every `responding-to:` reference, and every `excluded:`/`reincluded:`
record — none of which is a mechanical find-replace because the composite
changes again on the next model bump.

For B, I want `parley roster` as a **distinct verb** (not `parley init
--reinit`), because `parley init` has a different job (create the deck
scaffold + transport) and overloading it conflates "first creation" with
"rebuild the roster on a live deck." The session/machine scope split maps
cleanly onto the existing config layers (`configLayers` in
`internal/config/runtime.go:127`), and reinit must be **open-idea-aware**: it
must never silently rewrite the quorum of an in-flight idea.

For C, autonomous mode is a per-CLI flag string stored as a first-class
`autonomous_write_mode` field on the spec, scoped to the workspace, with
secret-redaction preserved. The skill states the requirement generically; the
per-CLI mapping lives in config + builtin specs.

## Proposed approach

### A — Naming scheme (sanitization, parse, collision, ID-vs-display)

**Crux position: stable base ID + composite display name.**

- **Base ID** (the identity): short, stable, kebab-case, assigned once at deck
  creation, never changes for the agent's lifetime in the project. Today's
  family IDs (`claude`, `codex`, `agy`, `hermes`) are *too* short and collide
  with the vendor name; the §2 roster already tried to fix this with `claude-1`
  but the runner ignored it. I would make the base ID the **§2 roster ID**
  (`claude-1`, `hermes-1`, `agy-1`, `codex-1`, `kimi-1`) and **wire the runner
  to use it** — that is the actual bug this idea should fix. The base ID is what
  appears in `round-01/<base-id>.md`, `### Signoff: <base-id>`, `agent:
  <base-id>` frontmatter, `responding-to:`, `participants:`, and `excluded:`.

- **Composite display name** (the descriptor): `agent-model-effort`, derived
  from the configured model + effort at reinit time, stored alongside the base
  ID in config and the §2 table, shown in the TUI/digest/runtime matrix. It is
  **regenerable**: change the model → the composite changes, the base ID does
  not. It is never used as a filename or a frontmatter key.

Why not (a) composite-as-ID: a model bump (`opus48`→`opus49`) or effort change
(`high`→`max`) would rename the agent and break continuity for every open and
archived idea that references it. The risk is not theoretical — `responding-to:
[claude-opus48-max/round-01]` becomes a dangling pointer, and `### Signoff:
claude-opus49-max` in a new round no longer matches the `claude-opus48-max`
that wrote round-01. Migration would require rewriting every historical file,
which violates "never edit another agent's file" and the static-FINAL rule.
Why not (c) composite-with-stable-prefix: the prefix *is* the base ID, so (c)
collapses to (b) with extra parsing; (b) is simpler.

**Sanitization (honoring the HARD `[a-zA-Z0-9_-]` constraint):**

1. Take the human model label; lowercase it.
2. Delete every char not in `[a-z0-9]` (dots, spaces, parens, brackets,
   slashes all vanish; adjacent alphanumerics join). This is the claude-1
   proposal and I agree — it is the only rule that is both deterministic and
   readable.
3. Effort token from a **fixed closed vocabulary**:
   `low|medium|high|xhigh|max|ultracode|clidefault` (with `cli-default` →
   `clidefault`). Closed vocab is what makes the parse unambiguous and the
   collision rule tractable.
4. Join: `<agent>-<model>-<effort>`, all lowercase, single `-` between
   sections. Each section is a single hyphen-free token, so the whole name
   splits cleanly on `-`.
5. **agy special case**: it has no per-invocation effort flag; its reasoning
   tier is baked into the model label (`(High)`). Surface that tier as the
   effort token. Concretely: strip the `(High)`/`(Low)`/`(Medium)` tier
   parenthetical from the model token *and* capture it as the effort token, so
   `Gemini 3.5 Flash (High)` → model `gemini35flash`, effort `high`. This keeps
   a single source of truth (the model label) and avoids inventing an effort
   flag agy does not have.

**The 5 resulting names** (claude / codex / hermes / agy / kimi; models Opus
4.8 / GPT-5.5 / GLM 5.2 / Gemini 3.5 Flash (High) / K3; efforts max / xhigh /
high / high / max):

| Base ID     | Composite display name          |
| ----------- | ------------------------------- |
| `claude-1`  | `claude-opus48-max`             |
| `codex-1`   | `codex-gpt55-xhigh`             |
| `hermes-1`  | `hermes-glm52-high`             |
| `agy-1`     | `agy-gemini35flash-high`        |
| `kimi-1`    | `kimi-k3-max`                   |

Note: these are **display names**. The artifact paths stay
`round-01/claude-1.md`, `round-01/hermes-1.md`, etc. The §2 roster table grows
a fourth column ("Composite") so a glance at COOPERATION.md shows both.

**Parse rule** (for tooling that needs to split a composite back out — e.g. the
TUI rendering a chip):

1. Split on `-`.
2. `token[0]` = agent family.
3. If `token[-1]` is all digits → it is an instance index; effort = `token[-2]`;
   the remaining single token is the model. (Instance index only arises on
   collision, see below.)
4. Else effort = `token[-1]`; the single token between agent and effort is the
   model.
5. The model token is **non-canonical** (it is a display slug, not the real
   model id). Tooling that needs the real model id reads it from config, not
   by reversing the slug. This is important: the composite is *lossy* by
   design; it is a label, not a database key.

**Collision rule:** collisions happen when two rostered agents share the same
agent family + model + effort (e.g. two `claude` instances on the same model).
On collision, append a numeric instance suffix to the **composite only**
(`claude-opus48-max-2`); the base IDs are already unique (`claude-1`,
`claude-2`). Because the base ID is the identity, a composite collision is
cosmetic, not structural — the two agents are still distinguished by their
base IDs in all paths and signoffs.

**Where the composite lives:** `parley-deck/agents.toml` and
`~/.parley/agents.toml` get a `display_name = "claude-opus48-max"` field per
`[agents.<id>]` block (or it is derived on load from `model`+`reasoning` and
cached). The §2 roster table in COOPERATION.md grows the Composite column. The
TUI runtime matrix (`agents.PrintRuntimeMatrix` in
`internal/agents/discover.go:314`) prints the composite next to the base ID.

### B — `parley roster` command (surface, scope, reconciliation)

**Command surface — distinct verb, not `parley init --reinit`:**

```
parley roster init [--scope session|machine] [--reinit] [--from <file>] [--yes] [--json]
parley roster show  [--scope session|machine] [--json]
parley roster diff  [--scope session|machine] [--json]   # what would change vs current
```

Why a distinct verb: `parley init` (app.go:362) creates the deck scaffold
(`COOPERATION.md`, `ideas/`, `inbox/`, `meta/`) and seeds transport. That is a
*creation* job. Roster rebuild is a *maintenance* job on a live deck and can
destabilize open ideas. Overloading `init --reinit` hides that distinction and
makes the help text lie. A `roster` verb group also leaves room for `roster
add`/`roster drop`/`roster retire` later without re-shaping the CLI.

**Interactive flow (the reinit UX I own):**

1. **Discover** every agent CLI available on the machine. Today `agents.Discover`
   (discover.go:261) probes `exec.LookPath` + a version probe over
   `agents.DefaultSpecs()` + ACP catalog + user config. reinit reuses this but
   adds a **model+effort discovery** step that does not exist yet (see
   "Discovery" below). For each found agent, list `(model options, effort
   options)`.
2. **For each agent, prompt model + effort.** Seed the default from the
   existing config (the central `[agents.<id>].model/reasoning` if present,
   else the builtin spec default). The **default effort is the strongest
   (highest) level the agent supports** (mirrors the skill's bootstrap rule and
   `centralDefaultTemplate` in runtime.go:277). Present a numbered list; Enter
   accepts the default.
3. **Compute the composite** per the A sanitization rules and show a preview
   table: `base-id | composite | model | effort | source`.
4. **On confirm, write** the roster. Composite display names go into the
   agents.toml `[agents.<id>]` blocks (a `display_name` field, or derived);
   the §2 roster table in COOPERATION.md is updated with the Composite column.
5. **Open-idea guard** (see below) fires *before* any write.

**Session vs machine scope — how the writes differ:**

The scopes map exactly onto the existing `configLayers` (runtime.go:127):

- **`--scope machine`** writes `~/.parley/agents.toml` only (via
  `config.EnsureCentralDefault`'s path, but a *rewrite* not a seed). It changes
  the central default inherited by *every* project. It does NOT touch any
  deck. This is the right knob for "I got a new machine / I want my global
  defaults to reflect today's models."
- **`--scope session`** writes the **deck override layer**
  (`parley-deck/agents.toml`, and `meta/headless-agents.local.json` for
  machine-local launch flags). It overrides the central default for *this
  project only*. It also updates the §2 roster table in this deck's
  COOPERATION.md. This is the right knob for "this project uses a different
  model/effort mix than my global default."

Precedence is unchanged (central < deck < env). reinit-session writes an
**override layer**, not a full copy of the central roster: it writes only the
`[agents.<id>]` blocks that differ plus the `[defaults]` knobs the user
touched. This keeps the deck file small and the diff readable. reinit-machine
writes the full central file (it *is* the source).

Default scope when the flag is omitted: **`session`** if a deck exists in the
cwd (it is the more common, lower-blast-radius operation and matches where the
user already is), **`machine`** if no deck exists (you cannot write a session
roster for a deck that isn't there; fall back to seeding the central default,
with a hint to run `parley init` then `parley roster init --scope session`).

**Reconciliation with `parley init` bootstrap and §9.0:**

- **`parley init`**: today the Go CLI's `runInit` (app.go:362) does scaffold +
  central seed + transport, and **no** interactive roster/model/effort
  confirmation — that gate lives only in the *skill* (SKILL.md "Deck
  bootstrap"). This is the duplication risk: the skill tells the facilitator to
  run an interactive list-roster → confirm-models-and-effort → pick flow, but
  the CLI has no command for it, so the facilitator improvises. `parley roster
  init` becomes the **mechanical home** for that flow: the skill's bootstrap
  section says "run `parley roster init --scope session`" instead of
  hand-rolling the prompts. `parley init` stays the scaffold+transport step;
  at first deck creation the skill chains `parley init` → `parley roster init
  --scope session`. No duplicated confirmation, one command owns the UX.
- **§9.0 readiness ping**: preflight (preflight.go) only pings *liveness*
  per-idea. reinit is a heavier, user-initiated, deck-level operation. They do
  not overlap: preflight never selects models or effort, and reinit is not run
  per-idea. The only touchpoint is that reinit should optionally run a
  liveness probe against the newly-selected roster (reuse `hostedPONG` via
  `pingProbe`) so the user gets immediate feedback that their picks actually
  boot — but this is a post-write verification, not a gate.

**Open-idea guard (must not silently rewrite an OPEN idea's quorum):**

This is the hardest failure mode and it is not optional. Before writing, reinit
scans `parley-deck/ideas/*/00-prompt.md` for `status:` in
(`round-01`,`round-02`,…,`consensus`,`final`-but-not-archived, i.e. any
non-archived open idea). For each open idea it reads `participants:`. If the
proposed roster **drops** a participant that is in an open idea's quorum, or
**renames** its base ID, reinit **refuses by default** and prints:

```
refusing to write: open idea <slug> has participants [claude-1, hermes-1] but
the new roster drops claude-1. Confirm with --yes to record the quorum break
in each affected idea's 00-prompt.md (excluded:/reincluded:) and proceed.
```

`--yes` proceeds but **writes the break into each affected idea's
`00-prompt.md`** (an `excluded:`/`reincluded:` line, reusing the exact format
the preflight-readiness idea already ratified) — never silently. Adding a new
agent is not a break (it is not in any open idea's quorum yet); the new agent
just becomes available for future ideas. Changing an *existing* agent's model
or effort (composite changes, base ID stable) is **never** a quorum break —
that is the whole point of keeping the base ID stable — so model/effort
re-picks are always safe and never trip the guard.

**Non-interactive / CI path:**

`--from <file>` loads a pre-authored roster spec (TOML/JSON) and writes it
without prompting. `--yes` skips the open-idea guard confirmation (but still
*records* any quorum break it detected — CI does not get to silently break
quorum either; if a break is detected and `--yes` is set, it writes the
`excluded:` lines and exits 0; if a break is detected and `--yes` is not set,
it exits 3 with a JSON report of the affected ideas). `--json` emits the
preview/diff/break report machine-readably. Exit codes mirror preflight: 0
ok, 1 hard failure, 2 usage, 3 pending guard.

**Failure modes I want the design to call out explicitly:**

- **Model discovery unavailable** (a CLI exposes no `models`/`model list`):
  fall back to the configured/builtin model string, present it as the only
  option, let the user type an exact id. Record `model: cli-default` if they
  decline. Never invent model names.
- **Effort discovery unavailable**: same; default to the strongest *known*
  level for that family (from a builtin table) or `clidefault`.
- **`agy` tier parsing**: if the model label has no `(Tier)` parenthetical,
  effort falls back to `clidefault` and a warning is printed.
- **Write fails** (permissions, disk): atomic write to a temp file then rename;
  never leave a half-written agents.toml. The §2 COOPERATION.md edit is the
  riskier write — use the same allowlisted-zone approach the drift guard
  already uses (`TestEmbeddedDefaultMatchesLiveDeck`) so only the roster table
  block is touched.
- **Central write on a read-only home** (`--scope machine`, no write access to
  `~/.parley`): error early, name the path, suggest `--scope session`.
- **Two agents collide on composite**: append the numeric suffix to the
  composite only; warn the user that two rostered agents are
  indistinguishable by display name and should get distinct base IDs.

### C — Autonomous / "yolo" write mode per agent

There is no common flag across vendors, so this is a **per-agent string field**,
not a boolean. Concretely, add to `agents.Spec` (discover.go:12) and to the
`agentOverride` TOML shape (runtime.go:88):

```
AutonomousWriteMode string  // toml: "autonomous_write_mode"
AutonomousWriteArgs []string // toml: "autonomous_write_args"  (optional override)
```

`AutonomousWriteMode` is a short enum-ish token identifying the *strategy*
(`bypass-permissions`, `workspace-write`, `accept-edits`, `accept-hooks`,
`print-auto`), and `AutonomousWriteArgs` is the concrete argv the runner
appends to realize it. The builtin specs already hardcode these today
(e.g. claude `--permission-mode acceptEdits`, codex `--sandbox
workspace-write`, agy `--dangerously-skip-permissions`, hermes
`--accept-hooks`); this field just makes them first-class and named instead of
buried in `HeadlessArgs`. For kimi headless `-p` print mode, the strategy is
`print-auto` (the brief's note that `--yolo`/`--auto` are mutually exclusive
with `-p`, so `-p` *is* the autonomous mode for kimi headless) — so
`AutonomousWriteArgs = ["-p"]` and no separate yolo flag is added.

**What the skill says (generic requirement + per-CLI mapping):**

> Every headless participant MUST be invoked in its CLI's non-interactive
> auto-approve mode so it can write its own canonical artifact
> (`round-NN/<id>.md`, signoffs, review files) without a blocking permission
> prompt. The mode is workspace-scoped: it authorizes writes only within the
> deck/workspace root, never a blanket machine-wide bypass. The per-CLI
> mapping is recorded in the agent spec's `autonomous_write_mode` field and
> applied by the runner. Obvious-secret redaction is preserved regardless of
> mode: credentials, customer data, and private documents unrelated to the
> task are still not sent to external backends.

The skill then lists the mapping (claude `--permission-mode bypassPermissions`
or `acceptEdits`; codex `-s workspace-write`; hermes `--yolo`; agy
`--dangerously-skip-permissions`; kimi `-p`). The skill does not hardcode the
argv in prose — it points at the spec field as the source of truth so a
vendor flag change is a config edit, not a skill revision.

**Safety (workspace-scoped, secret-redaction preserved):**

- The runner already passes `{root}` into `HeadlessArgs` (discover.go:134, 156)
  via `config.ExpandPlaceholders` (runtime.go:315); the autonomous args are
  scoped to that root. The design must state that `autonomous_write_mode` never
  implies a broader scope than the workspace root the runner already passes.
- Secret redaction is a separate, pre-existing invariant (the skill's
  external-backend disclosure rule + LE-4's "output secret-scrubbed" for
  `checks:`). Autonomous mode does not relax it. The design should say so
  explicitly so a future reader does not assume "yolo = send everything."
- `bypassPermissions` (claude) is the one mode that *can* escape the workspace;
  the design should prefer `workspace-write` / `acceptEdits` where the CLI
  offers a narrower mode, and treat `bypassPermissions` as an opt-in escalation
  the user confirms at reinit (a prompt like "claude autonomous mode:
  workspace-write (recommended) or full bypass?").

## Concerns / open questions

1. **The runner→§2-roster gap is the real bug.** Today `agent.ID` is the family
   ID and the §2 `claude-1` is ignored by the runner. My proposal wires the
   runner to the base ID (`claude-1`). That is a **behavior change for every
   existing deck**: artifact paths would change from `round-01/claude.md` to
   `round-01/claude-1.md`. Is the group willing to take that migration
   (one-time rename of existing artifacts + frontmatter `agent:` values), or
   do we keep the family ID as the base ID and accept that the §2 `claude-1`
   stays cosmetic? I lean toward taking the migration because the composite
   name is pointless if the identity layer is still the bare family name — but
   this is the single biggest cost in the idea and needs a deliberate decision.

2. **Composite derivation vs storage.** Do we store `display_name` explicitly
   in agents.toml (durable, can drift from model+reasoning) or derive it on
   load (always fresh, but the §2 table needs a re-render after a config
   edit)? I lean derive-on-load + re-render the §2 table on `parley roster
   show`/`preflight`, and store nothing — single source of truth.

3. **agy effort tier as a parseable token.** Stripping `(High)` from the model
   and surfacing it as effort works for agy, but what happens when a future
   Gemini model label has no tier? Defaulting to `clidefault` is safe but
   makes the composite less self-describing. Acceptable, but worth flagging.

4. **`parley roster init` on a deck with no §2 table yet** (fresh `parley
   init`): should it create the §2 table rows, or expect the user to hand-edit
   COOPERATION.md first? I think reinit should create/refresh the table rows
   itself (within the allowlisted zone), otherwise the bootstrap chain is
   incomplete.

5. **kimi is not a builtin launch spec** — it is only in the ACP catalog
   (`acp_specs.go:37`). reinit's discovery must handle agents that are
   *user-configured only* (present in agents.toml but not in
   `DefaultSpecs()`), which `applyFile` already permits (runtime.go:361-381).
   The model/effort discovery step for kimi needs a documented fallback since
   there is no builtin spec to seed from.

## Risks

- **Renaming an agent mid-idea (if composite were the ID).** Orphaned
  `round-NN/<id>.md`, dangling `responding-to:`/`### Signoff:` references,
  broken `excluded:`/`reincluded:` audit trail, and a violation of the
  static-FINAL / never-edit-another-agent's-file rules when migrating. This is
  why the base ID must be stable. Severity: CRITICAL if we pick (a); avoided
  by (b).
- **Open-idea quorum break on reinit.** Dropping or renaming a participant in
  a roster that an in-flight idea depends on silently shrinks quorum (§5).
  The open-idea guard mitigates; the residual risk is a race where an idea is
  opened *during* the reinit write. Mitigation: reinit reads the idea list
  once at start and again at write time, and aborts if it changed.
- **§2 COOPERATION.md edit clobbering project-specific zones.** reinit writes
  the roster table in the allowlisted zone only; the drift guard
  (`TestEmbeddedDefaultMatchesLiveDeck`) must still pass. Risk: the reinit
  table writer and the drift-guard allowlist disagree on what the "roster
  zone" is. Mitigation: share one allowlist definition between them.
- **Model/effort discovery flakiness.** Running `<cli> models` can be slow,
  rate-limited, or require auth. reinit must not block on it — fall back to
  configured values and warn. Risk: a user ships a composite with a stale
  model slug because discovery silently failed. Mitigation: always print
  whether each agent's options came from discovery or from config fallback.
- **Autonomous mode escape.** A user picks `bypassPermissions` for claude at
  reinit and the runner launches it outside the workspace root. Risk:
  machine-wide writes. Mitigation: the design names `workspace-write` as the
  default and treats `bypassPermissions` as a confirmed escalation; the runner
  still passes `{root}` so the CLI's own scoping applies; and the skill's
  secret-redaction rule is restated as mode-independent.
- **Composite collision suffix vs base ID.** If two `claude` agents get
  `claude-opus48-max` and `claude-opus48-max-2`, a reader may assume `-2` is
  the base ID instance index. Mitigation: the §2 table always shows base ID
  and composite side by side so the distinction is never implicit.

## Concrete Go touch points

- `internal/runner/runner.go:349` — `agent.ID+".md"`: the filename
  construction. If we adopt the base-ID-as-identity position, the runner must
  receive the base ID (not the family ID) here. Today `selectedAgents`
  (runner.go:~328) keys by `agent.ID` from `Discovery`, so the wiring change
  is: resolve each selected participant to its §2 base ID before launch, and
  pass that as the artifact identity.
- `internal/agents/discover.go:12` (`Spec`) and `:92` (`Discovery`): add
  `AutonomousWriteMode` + `AutonomousWriteArgs` fields; `defaultBuiltinSpecs`
  (discover.go:105) gets the per-CLI autonomous args factored out of
  `HeadlessArgs`.
- `internal/config/runtime.go:88` (`agentOverride`): add
  `AutonomousWriteMode`/`AutonomousWriteArgs` TOML fields;
  `applyOverride` (runtime.go:411) applies them; `configLayers` (runtime.go:127)
  is the scope split reinit writes into; `EnsureCentralDefault` (runtime.go:255)
  + `centralDefaultTemplate` (runtime.go:277) get the `display_name`/autonomous
  fields in the seed.
- `internal/config/roster.go` + `internal/protocol/roster.go`: the §2 roster
  table parser/generator. reinit needs a *writer* counterpart to
  `ReadRosterIDs` that updates the roster table block within the allowlisted
  zone, adding the Composite column.
- `internal/app/app.go:362` (`runInit`) + `:59` dispatch: add `roster` case;
  `runRoster` subcommands `init`/`show`/`diff`. `selectedParticipantIDs`
  (app.go:2389) and `installedAgentIDs` (app.go:2379) become base-ID-aware.
- `internal/app/preflight.go`: reuse `pingProbe`/`hostedPONG` for the
  post-write liveness check in `parley roster init`.
- New: model/effort discovery helpers (e.g. `internal/agents/models.go`) —
  per-CLI `<cli> models` parsing with fallback; this does not exist today.
- Tests: sanitization + parse + collision (table-driven, the 5 names above);
  reinit session-vs-machine write targets; open-idea guard (drop/rename/add
  cases); composite-on-load derivation; autonomous-mode field round-trip.
  `go build ./...`, `go vet`, `gofmt -l` must stay clean (per the idea
  constraints); the drift guard `TestEmbeddedDefaultMatchesLiveDeck` must stay
  green after the §2 table writer lands.
