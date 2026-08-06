---
agent: codex-1
idea: roster-operations-standard
round: 1
date: 2026-08-06
---

## Summary

The standard should make `parley roster show` the only answer to “what is the current agent roster?”, backed by one versioned row type and one effective-invocation resolver. `parley agents list` should remain, but be documented and named as an adapter/runtime inventory rather than a roster. The skill should execute the roster command and reproduce its table; it should not parse configuration or construct a competing table.

I propose `roster update --scope local|global` for deliberate changes and a single-direction `roster sync` whose only meaning is global → local. “Session” should cease to mean “deck”: a running Parley run should have an immutable, auditable roster snapshot, while local roster changes apply to later runs. This separates persistence from execution and prevents a mid-run global edit from silently changing later participants.

The factual statements below are claims I introduce and own. Their `PRIMARY` / `RECALL` tags state provenance, not a verification verdict; I do not assign `CONFIRMED`, `WRONG`, or an equivalent verdict to them.

## Proposed approach

### 1. Findings from the current implementation

**Two row sets and two meanings**

- `[PRIMARY P1]` I executed `parley --version`, `parley roster show`, `parley roster show --json`, and `parley agents list` from the repository root. The installed executable reported `parley 1.39.0`. `roster show` emitted four §2 identities (`claude-1`, `codex-1`, `hermes-1`, `kimi-1`) with seven columns, while `agents list` emitted eighteen adapter IDs with twelve different columns. The latter included `opencode`; the former did not.
- `[PRIMARY P2]` `internal/app/roster.go:98-118` obtains roster rows from `protocol.ReadRosterIDs(root)`, which parses the first §2 table in `COOPERATION.md` (`internal/protocol/roster.go:11-69`). In contrast, `agents list` loads and discovers the complete spec catalog, then `PrintRuntimeMatrix` iterates every discovery (`internal/agents/discover.go:483-565`). Thus adapter availability and protocol roster membership are separate sets in the current code.
- `[PRIMARY P3]` The top-level dispatcher accepts `roster` at `internal/app/app.go:100-101`, but the main usage and command descriptions at `internal/app/app.go:119-217` omit it. `parley --help` likewise omitted `roster` during the check. The command exists but is not presented as the canonical roster surface.

**Declared model versus launched model**

- `[PRIMARY P4]` The current built-in Claude spec contains `--model claude-opus-4-8[1m]` in `HeadlessArgs` and separately stores `Model: "claude-opus-4-8[1m]"` (`internal/agents/discover.go:213-228`). The effective project config changes only the `Model` field to `claude-opus-5[1m]`. `applyOverride` replaces `HeadlessArgs` only when `headless_args` is supplied and updates `Model` independently (`internal/config/runtime.go:542-545,594-605`). The runner substitutes only `{root}` and `{prompt}` in the resolved argument vector (`internal/runner/runner.go:1094-1121`). In the executed `agents list`, the Claude row declared `MODEL=claude-opus-5[1m]` while its printed effective argv retained `--model claude-opus-4-8[1m]`.
- `[PRIMARY P5]` The executed runtime matrix declared `codex MODEL=gpt-5.6-sol` although its effective argv contained no `-m/--model`, and declared `kimi MODEL=kimi-code/k3` although its argv was only `kimi -p {prompt}`. The installed CLI help showed model flags for all relevant adapters: Codex `-m/--model`; Claude `--model` and `--effort`; Hermes `--model` and `--reasoning`; Kimi `-m/--model`; OpenCode `-m/--model` and provider-specific `--variant`.
- `[PRIMARY P6]` `resolveRoster` copies `spec.Model` and `spec.Reasoning` directly into the roster row (`internal/app/roster.go:149-155`) rather than deriving them from the invocation built by the runner. Its `MODEL` and `EFFORT` cells are therefore declarations under the present design. `AUTO` already takes the stronger approach by calling `AutonomousEffective()`, which checks that the enabling arguments actually occur in `HeadlessArgs` (`internal/agents/discover.go:111-139`). Model and effort should use the same effective-value principle.

**The current “session” scope is deck scope**

- `[PRIMARY P7]` `roster init` accepts `--scope session|machine`, but `rosterTargetPath` maps `session` to `<root>/parley-deck/agents.toml` and `machine` to `~/.parley/agents.toml` (`internal/app/roster.go:213-235,383-389`). The former persists for the project and is shared by later runs; it is not scoped to one run or process.
- `[PRIMARY P8]` I enumerated Go files with `find internal cmd -type f -name '*.go' -exec /usr/bin/grep ...`, rather than an ignore-sensitive grep. The roster/session references led to `internal/app/roster.go`, the config loaders, the runner mapping loader, and tests. The actual session index stores workspace, run ID, idea, task, participant IDs, and timestamps (`internal/sessionstore/sessionstore.go:18-28`); the run manifest stores participant IDs but no mutable roster configuration (`internal/runmanifest/manifest.go:28-46`). This is evidence for treating `--scope session` as a nominal/misleading label, not as a distinct roster store.
- `[PRIMARY P9]` A new run holds the discovered specs in `runner.Options.Agents` (`internal/runcontrol/runcontrol.go:99-113`) and a live run continues to use that in-memory value. `run.created` records model/reasoning/source metadata (`internal/runcontrol/runcontrol.go:152-177`), but the manifest has no roster revision and the recorded runtime does not contain the complete materialized invocation. A long-lived process can therefore keep a launch configuration that predates a config edit.

**Why yesterday’s OpenCode change produced inconsistent sessions, including Hermes**

- `[PRIMARY P10]` The current global file `~/.parley/agents.toml` has `[agents.opencode]` and `[roster.opencode-1]`, but this repository’s §2 roster has only the four rows printed by `roster show`, and its `parley-deck/agents.toml` has no OpenCode block. Because membership is read from §2 while adapter inventory is read from layered specs, adding a global adapter/mapping did not add a project roster member.
- `[PRIMARY P11]` The dated global backup `~/.parley/agents.toml.bak-2026-08-06-pre-autofix:38-44` shows a Hermes `headless_args` replacement without `--yolo`. This repository’s project override includes `--yolo` (`parley-deck/agents.toml:61-66`). Since a higher config layer replaces the argument slice wholesale (`internal/config/runtime.go:542-545`), decks with the project override produced an autonomous Hermes invocation while decks inheriting the then-global argument list did not. The current global file now contains `--yolo`, but a process that already captured the earlier discovery remains on its captured value.
- `[PRIMARY P12]` The local session bootstrap text also conflates the two surfaces: `~/.parley/SESSION_PROMPT_EXISTING.txt:3` calls `parley agents list` “roster, models, reasoning, [defaults]”, even though the command prints the adapter catalog. `parley-deck-skill status --target all --project . --json` additionally reported installer `2.4.0` while the detected runtime skill copies, including Codex, Claude, Hermes, Kimi, and OpenCode, were marked `2.3.0`. These are additional ways two newly opened agent sessions can receive different roster instructions even when they share the same global TOML.

Relevant external snippets copied here so later participants do not need access to the user-home files:

```toml
# Current ~/.parley/agents.toml (abridged)
[agents.opencode]
model = "litellm/xai/grok-4.5"
headless_args = ["run", "--auto", "-m", "litellm/xai/grok-4.5", "{prompt}"]

[roster.opencode-1]
adapter = "opencode"

# ~/.parley/agents.toml.bak-2026-08-06-pre-autofix
[agents.hermes]
model = "glm-5p2"
headless_args = ["--oneshot", "{prompt}", "--model", "glm-5p2", "--accept-hooks"]
```

```text
# ~/.parley/SESSION_PROMPT_EXISTING.txt:3 (abridged)
Load the effective agent settings: `parley agents list --dir .` (...) — roster, models, reasoning, [defaults].
```

`[PRIMARY P13]` The OpenCode adapter did not itself alter Hermes. The inconsistency followed from updating only some independently interpreted layers: global adapter config, global roster mapping, project §2 membership, project argument overrides, installed skill copies, and already-running process snapshots. The fix is to give those layers named responsibilities and to route all roster operations through one resolver.

### 2. Canonical roster table contract

Use this exact ordered column set for schema version 1:

| Column | Meaning and source | Why it is load-bearing |
|---|---|---|
| `AGENT` | Stable roster ID, such as `hermes-1`; local membership comes from §2, global membership from `[roster.<id>]`. | Artifact identity and quorum use this value; it must not be confused with an adapter or display label. |
| `ADAPTER` | Resolved launch adapter from layered `[roster.<id>].adapter`. | Explains which CLI contract materializes the invocation; Hermes/OpenCode are adapters, not model companies. |
| `STATE` | `active` or `inactive`; local state mirrors the §2 inactive marker, global state uses `[roster.<id>].active` with existing entries defaulting to active. | Preserves historical identities without treating retired rows as participants. |
| `INSTALLED` | `yes` or `no`, from the resolved explicit command/path lookup. | A configured roster member may be unavailable on the current machine. |
| `MODEL` | Exact **effective** model reference from the same materialized invocation plan the runner will execute; `unknown` when the adapter cannot establish one. | This is the operational truth users need. A declaration that is not delivered to the process must not occupy this cell. |
| `MODEL-FAMILY` | Deterministic lineage derived from `MODEL`, for example `Claude Opus`, `GPT`, `GLM`, `Kimi K`, or `Grok`; `unknown` if unresolved. | Supports diversity and capability reasoning without grouping by CLI vendor. |
| `MODEL-COMPANY` | Underlying model maker derived from `MODEL` after gateway prefixes are removed; `unknown` if unresolved. | Separates the company that made the model from the CLI and routing gateway. |
| `ROUTE` | Outermost routing/provider layer, such as `direct`, `default`, `LiteLLM`, or another recognized gateway; `unknown` when the model reference and adapter config do not establish it. | Gateway-routed models otherwise make `MODEL-COMPANY` ambiguous. |
| `EFFORT` | Effective per-invocation effort/reasoning value; `unknown` when only a declaration exists and the adapter cannot bind or inspect it. | Prevents the table from repeating the declared/effective model defect for reasoning. |
| `SPEED` | Resolved Parley scheduling/output profile from the layered roster policy. | Speed is deliberately independent of model reasoning and should remain visible as a separate axis. |
| `AUTO` | `yes` only when the materialized invocation enables the adapter’s autonomous-write contract; otherwise `no`. | This is required to know whether the agent can create its own artifact unattended. |
| `STATUS` | `ok` or stable comma-separated issue codes such as `unmapped`, `not-installed`, `model-drift`, `effort-unknown`, `auto-off`, `masked-by-env`, `metadata-unknown`, or `stale-snapshot`. | One compact place flags divergence without replacing the effective cells with declarations. |

`DISPLAY-NAME` should leave the canonical table. It duplicates adapter/model/effort, is not identity, and is currently capable of contradicting `MODEL`. It may remain a derived JSON field during a compatibility period, but it must be derived from the effective row rather than `ModelLabel` plus declarations.

The text header and order are an API. `roster show --json` should return `schema_version`, `scope`, `roster_revision`, `columns` in that order, and typed `rows`. A column addition or semantic change requires a documented schema-version change and golden tests. `unknown` and `no` must be explicit; empty cells conceal uncertainty.

Detailed per-field provenance is useful but too wide for the canonical table. Add `parley roster show --explain AGENT` and include a `sources` object in JSON, with the winning layer and declared/effective values for each field. This is where a user sees, for example, that project `agents.toml` masks a global value.

### 3. Model, family, company, and route derivation

Create a CLI-owned `modelmeta` resolver. Deck files must not gain `model_family` or `model_company` fields.

1. Start with the effective model reference from the materialized invocation, never `Spec.Model` alone.
2. Parse the reference into outer route segments and an inner producer/model segment. Recognized gateways such as `litellm` or `openrouter` populate `ROUTE` and are peeled before company derivation.
3. Apply a versioned, tested registry of producer aliases and model-prefix patterns. It should be data/code shipped with the CLI, not copied into every deck. Unknown strings return `unknown` plus `metadata-unknown`; the resolver must not infer company from `ADAPTER`.
4. Normalize only for metadata display. Preserve the exact effective reference byte-for-byte in `MODEL`.

Under the proposed registry, `glm-5p2` should group as family `GLM` and company `Zhipu AI`; because its reference is unqualified, route remains `default` or `unknown` unless the Hermes adapter can read an explicit provider. `litellm/xai/grok-4.5` should peel `litellm` into `ROUTE=LiteLLM`, then group the inner model as `MODEL-FAMILY=Grok` and `MODEL-COMPANY=xAI`. The adapter remains `opencode`. This rule avoids the two common category errors: calling Hermes the model company, or calling LiteLLM the model company.

The registry will require maintenance as model naming changes. That trade-off is preferable to per-deck hand-maintained metadata, and `unknown` makes registry lag visible instead of guessing.

### 4. One effective-invocation resolver

Introduce one internal `RosterView`/`RosterRow` domain service used by `roster show`, preflight, run selection, run snapshots, and the TUI. It should call the same `BuildInvocationPlan` used by the runner.

Each built-in adapter should declare structured bindings for model and effort rather than embedding mutable literals in arbitrary `HeadlessArgs`. Conceptually:

```go
type ValueBinding struct {
    Kind      string // flag, config-override, native-config, unsupported
    Flag      string
    ConfigKey string
}

type InvocationPlan struct {
    Path           string
    Args           []string
    ModelEffective string
    EffortEffective string
    Issues         []string
}
```

For the installed adapters, materialization can use Claude `--model/--effort`, Hermes `--model/--reasoning`, Codex `--model` plus its config-override mechanism for reasoning, Kimi `--model`, and OpenCode `--model/--variant` where the configured variant is accepted. If an adapter cannot bind or inspect effort, `EFFORT` must be `unknown` with `effort-unknown`, even if TOML contains a desired value.

Backward compatibility needs an in-memory legacy normalizer. Existing `headless_args` remain accepted. For a known adapter, the planner should identify its model/effort flags, replace stale values with the resolved structured value, or inject missing flags. It should report `legacy-argv-normalized` in explanation output but not require users to edit existing decks. For an unknown/custom adapter without binding metadata, leave its argv unchanged and report the effective model as `unknown`; do not claim the declaration was launched.

The immediate regression tests should cover:

- a project model override changing Claude’s actual `--model` argument;
- Codex and Kimi receiving explicit configured model arguments;
- Hermes preserving `--yolo` while changing model/reasoning;
- OpenCode preserving `--auto`, model routing, and prompt position;
- a custom adapter with an unbound declaration yielding `MODEL=unknown` rather than a false value;
- `AUTO` being computed from the final materialized argv, not from pre-materialization declarations.

### 5. Command surface

Make these commands visible in top-level help:

```text
parley roster show [--scope local|global] [--dir DIR] [--json] [--explain AGENT]

parley roster update AGENT --scope local|global
    [--adapter ADAPTER] [--state active|inactive]
    [--model MODEL] [--effort EFFORT] [--speed SPEED]
    [--dry-run] [--yes] [--json]

parley roster sync [--dir DIR] [--dry-run] [--yes] [--json]
```

Semantics:

- `show` defaults to `local` when a deck exists and `global` otherwise. Both scopes use the same columns and JSON schema.
- `update --scope global` edits roster-managed fields in `~/.parley/agents.toml`. `update --scope local` edits the project’s roster-managed fields and, for state changes, the §2 roster row. It never edits an open idea’s `participants:` list.
- Extend existing `[roster.<id>]` blocks with optional `active`, `model`, `effort`, and `speed`. Missing fields inherit. This supports two roster IDs using the same adapter with different models and avoids adding another configuration file. Existing blocks containing only `adapter` continue to load.
- Resolution precedence should be field-wise and low-to-high: built-in adapter → global `[agents]` → global `[roster.ID]` → project `[agents]` → project `[roster.ID]` → `agents.local.toml` equivalents → environment config. Roster-ID fields override adapter-family defaults only within the same or higher layer.
- `sync` has exactly one direction and no `--from/--to` alternatives: global → the deck named by `--dir`. It makes the local active/inactive roster and roster-managed adapter/model/effort/speed values inherit the global roster. Local-only active rows become inactive rather than being deleted; historical §2 rows remain. Non-roster runtime settings such as command path, sandbox, approval, timeout, and isolated-home policy are preserved.
- The sync operation should remove local roster-managed overrides that would mask the global values rather than copying a point-in-time value. The existing config layering then keeps the local deck synced until it intentionally introduces a new local override. This makes “sync” a rebase onto global defaults, not a one-time copy whose provenance is immediately lost.

Examples:

```sh
# Preview, then change only this deck.
parley roster update hermes-1 --scope local --model glm-5p2 --effort high
parley roster update hermes-1 --scope local --model glm-5p2 --effort high --yes

# Preview, then change the machine-wide default.
parley roster update opencode-1 --scope global --adapter opencode \
  --state active --model litellm/xai/grok-4.5 --effort xhigh
parley roster update opencode-1 --scope global --adapter opencode \
  --state active --model litellm/xai/grok-4.5 --effort xhigh --yes

# The one global-to-local operation.
parley roster sync --dir .
parley roster sync --dir . --yes
```

Safety properties:

- Preview is the default whenever a mutation would occur; `--yes` is required to apply. `--dry-run` is retained for scripts and must be side-effect-free.
- Parse and validate every candidate before writing. Reject an unknown adapter, an active roster with no resolvable command, an invalid effort for a known closed vocabulary, or a requested model the adapter cannot bind without falsely reporting it effective.
- Acquire config/roster locks, preserve unrelated TOML and project-specific §2 text, write with atomic temp-and-rename operations, and roll back earlier files if a later write fails. Inactive entries and adapter mappings are retained, so a state change does not destroy historical resolution.
- Re-resolve after the write and compare the requested fields with the effective row. If `$PARLEY_HEADLESS_AGENT_CONFIG` or another higher layer masks the target, report `masked-by-env` and do not claim the requested effective change.
- Sync never copies local values upward. Any global mutation requires `roster update --scope global` explicitly.
- Sync may change defaults for future ideas, but never mutates locked `participants:` lists or a live run snapshot. A membership removal/addition is shown prominently in the preview because it changes future default quorum.

Deprecate the labels `session` and `machine` in favor of `local` and `global`. Accept `--scope session` as an alias for `local` and `--scope machine` as an alias for `global` for at least one compatibility cycle, with a warning. Keep `roster init` only as a deprecated mapping-repair/bootstrap alias; the skill must stop presenting it as update or sync.

### 6. Running-session semantics

Do not add a mutable session roster store. At run creation, write a secret-free immutable roster snapshot and `roster_revision` hash into the existing `run.json`/`run.created` state. The snapshot should include the canonical row values, source paths/field provenance, executable path, and a sanitized materialized invocation template (no prompt and no credentials).

The live process and all later phases of that run use the same snapshot. `sessions inspect` compares its revision with the current local roster and reports `stale-snapshot` when they differ. A user who wants the new roster starts a new run; an explicit future “rebase run” operation would need separate design because changing participant/runtime identity mid-idea can compromise the audit trail.

This gives “session” a useful, precise meaning—an immutable execution snapshot—without making it a third update scope. The single synchronization workflow is global update → local `roster sync` → new run.

### 7. Keep both commands, with a hard boundary

`parley roster show` and `parley agents list` should both remain:

- `roster show` is the operational/protocol roster: only roster identities, canonical columns, active/inactive state, effective launch values, and health issues. It is the answer to every natural-language roster question.
- `agents list` is the adapter inventory: every built-in/configured adapter, including uninstalled ACP backends, with diagnostic launch/config columns. Rename its help description to “adapter inventory (not the roster)” and optionally add `agents inspect` as a clearer alias, but keep `agents list` for compatibility.

The distinction should be visible in help and headings. `agents list` should never be called a roster by the skill, session bootstrap text, docs, or TUI.

### 8. CLI/skill boundary

The CLI enforces:

- membership/state resolution, config precedence, update/sync mutations, and locks;
- invocation materialization and declared/effective drift detection;
- model metadata derivation;
- canonical text/JSON schema, issue codes, and roster revision;
- immutable per-run snapshotting.

The Parley Deck skill says:

1. For “show/list/current roster”, run `parley roster show --dir <root>` and reproduce its stdout verbatim. If structured handling is necessary, consume `--json` but render exactly the CLI-provided `columns` order. Never reconstruct rows from `COOPERATION.md`, `agents.toml`, or `agents list`.
2. At session start, run the side-effect-free `parley roster sync --dir . --dry-run --json`. If it reports drift, show the plan and request approval before `--yes`; do not hand-edit files or silently change future quorum.
3. For a requested change, preview `roster update` or `roster sync`, show the canonical before/after tables, and apply only after the user authorizes it.
4. Record that open idea participants and live run snapshots remain unchanged.

Ship the skill text, CLI behavior, top-level help, session bootstrap templates, and golden schema tests in the same release. The CLI remains the authority even when an installed skill copy lags; a stale skill can still call the current command instead of carrying stale roster facts.

### 9. Migration and acceptance checks

No manual migration should be required:

- Existing §2 rows remain the local membership source until the user runs sync/update.
- Existing `[roster.ID] adapter` blocks load with `active=true` where global state is needed and inherit model/effort/speed.
- Existing `[agents.ADAPTER]` fields remain adapter defaults.
- Known built-in legacy model/effort flags are normalized in memory into the structured invocation plan; unknown custom adapters continue to run with visible `unknown` effective metadata.
- The old text table may be available behind `--legacy` for one release if scripts depend on it, while JSON schema v1 becomes the supported machine interface.

Acceptance should include golden and integration tests proving:

- the canonical column list is byte-identical in CLI text, CLI JSON `columns`, TUI roster view, and skill fixture;
- `roster show` returns four rows for the current §2 roster while `agents list` is explicitly labeled adapter inventory;
- a global OpenCode addition appears in a local deck only after the documented sync and then uses the expected effective `-m` model reference;
- a local override survives global update until sync, and sync removes the masking roster-managed override while preserving command/sandbox/timeout fields;
- a Hermes project without a local `headless_args` workaround still materializes `--yolo`, `--model glm-5p2`, and `--reasoning high` from structured fields;
- current runs retain their snapshot after global/local changes, while new runs receive the new revision;
- inactive rows remain resolvable historically but are excluded from default participants;
- gateway model metadata separates adapter, route, family, and company;
- masked environment overrides and unknown custom model bindings are visible and never reported as successful effective updates.

## Concerns / open questions

1. **Effort support is not uniform.** Kimi’s inspected help exposes a model flag but no effort flag; OpenCode’s `--variant` is provider-specific. The implementation must decide whether an adapter-specific native-config reader is sufficiently reliable. Until it is, `EFFORT=unknown` is more honest than displaying a desired TOML value.
2. **Exact-sync versus local exceptions.** I favor exact global inheritance for roster-managed fields because a merge policy would make “sync” ambiguous. This deliberately deactivates local-only active rows (without deleting them) and removes local model/effort/speed pins after preview. If the product instead preserves local exceptions, it needs an explicit origin/override display and a different verb such as `rebase`; it should not call a partial merge “synced”.
3. **Protocol table editing.** Local membership remains canonical in §2, so `update`/`sync` needs a narrow table editor that preserves workspace, role, and host-handle information. New IDs need deterministic defaults or explicit flags for those protocol-only fields. The CLI should not regenerate the entire cooperation document.
4. **Company registry naming.** The canonical display names (`Zhipu AI` versus another legal/brand name, `Moonshot AI` versus Kimi branding) need a small reviewed registry policy. Unknown should be acceptable and visible; per-deck overrides should not be allowed merely to improve cosmetics.
5. **Text width.** Twelve columns are wide. Stable semantics matter more than fitting 80 columns; the normal formatter can truncate display values with full values in JSON/`--explain`, or offer a vertical narrow-terminal view while retaining the same named fields.
6. **Scope vocabulary.** “Local” could mean machine-local to some users. Help must say `local = this deck/project` and `global = this user’s ~/.parley defaults` every time it previews a mutation.

## Risks

- A resolver that merely parses argv after the existing config merge would detect drift but would not fix it. The runner and table must consume the same structured invocation plan, or the project will recreate the `AUTO`/`MODEL` split under new names.
- Automatically rewriting model flags changes behavior in decks that intentionally duplicated contradictory `model` and `headless_args` values. The preview/explanation should expose normalization; the precedence rule must state that structured roster model wins for known adapters.
- Cross-file local mutation cannot be made perfectly crash-atomic with independent renames. Locking, prevalidation, rollback, retaining mappings, and writing the §2 membership change last reduce the failure surface; recovery should be idempotent and `roster sync --dry-run` should diagnose partial application.
- A model metadata registry can become stale. Returning `unknown` with a stable issue code is safer than inferring from adapter names or silently carrying obsolete company data.
- Stale installed skills will continue to exist in practice. The design remains robust only if skills delegate roster truth to the CLI and version drift is a warning, not an alternate roster source.
- Syncing membership affects defaults for future quorum. Even though locked ideas are preserved, every add/deactivate must be conspicuous in the preview and require explicit application.
