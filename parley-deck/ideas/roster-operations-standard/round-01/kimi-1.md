---
agent: kimi-1
idea: roster-operations-standard
round: 1
date: 2026-08-06
---

## Summary

The roster problem is not a missing table — it is that "roster" currently means three
different things stored in five places, and no command owns the answer. My proposal:
one canonical roster table with a frozen column contract (`parley roster show`, unhidden
and fixed), one update verb (`parley roster set`), one sync verb (`parley roster sync`,
deck ← machine), a structural fix so a configured `model` actually reaches the launch
argv, and an explicit CLI/skill boundary in which the skill only ever prints CLI output
and never maintains its own roster store.

Evidence base (all commands run 2026-08-06 against parley 1.39.0 at `/opt/homebrew/bin/parley`;
absence claims below were made with `rg -uuu` or per-file `/usr/bin/grep`, never an
ignore-honoring bare grep):

- **Three CLI surfaces, three row sets.** `parley roster show` prints 4 rows
  (`claude-1, codex-1, hermes-1, kimi-1`) with columns
  `ROSTER-ID FAMILY DISPLAY-NAME MODEL EFFORT SPEED AUTO` (PRIMARY: command output).
  `parley agents list` prints 19 rows with 12 different columns (PRIMARY: command output).
  `parley preflight` builds a third roster table (`rosterEntry{RosterID, Runtime,
  Version, Available, Reason}`, `internal/app/preflight.go:95-101`). `roster` is not in
  `parley --help` at all (PRIMARY: `parley --help`; dispatch exists at
  `internal/app/app.go:100` but `printUsage` omits it, `app.go:111-144`), and it is
  undocumented in `docs/cli-reference.md` and `docs/agent-runtime-configuration.md`
  (PRIMARY: per-file `grep -ci roster` → 0 in both).
- **The roster row set is the §2 table of COOPERATION.md, regex-parsed.**
  `resolveRoster` reads §2 via `protocol.ReadRosterIDs` (`internal/app/roster.go:110`,
  parser at `internal/protocol/roster.go:26-70`) and joins each ID through the
  `[roster.<id>] adapter` map (`internal/config/runtime.go:200-221`) into the family
  specs. That is why opencode — a full adapter since 1.39.0 (PRIMARY:
  `CHANGELOG.md:3-8`, dated 2026-08-06) with a mapping in `~/.parley/agents.toml:131-132`
  — does not appear: it is in no deck's §2 (PRIMARY: this repo's `COOPERATION.md` §2
  table has exactly 4 rows).
- **There is no session-scoped roster store; `--scope session` is nominal.**
  `rosterTargetPath` maps `session` → `parley-deck/agents.toml` and `machine` →
  `~/.parley/agents.toml` (`internal/app/roster.go:383-389`). `rg -uuu 'session'` over
  `internal/config/` finds zero matches. "Session" is a misleading label for the deck
  layer.
- **The skill maintains a sixth, divergent truth.** `parley-deck/meta/headless-agents.local.json`
  (updatedAt 2026-07-30) has no opencode entry, carries a kimi note stale since 1.36
  ("classifies kimi as an ACP backend (AUTO=no)", line 161), and lists claude-1
  headlessArgs (`--output-format json --no-session-persistence`, lines 18-25) that do
  not match the CLI's effective args (`--output-format text`, no such flags — PRIMARY:
  `parley agents list` headless line for claude). The CLI never reads this file
  (PRIMARY: `rg -uuu headless-agents` over `internal/` hits only protocol-defaults
  prose, no loader).

## Proposed approach

### 1. The canonical table — `parley roster show`, one contract

Fixed columns, identical in text and `--json`, additive-only across releases (a new
column is a changelogged change):

| Column | Source | Why it earns its place |
|---|---|---|
| `AGENT` | §2 roster ID (`ReadRosterIDs`) | The identity used in artifact paths and signoffs. |
| `FAMILY` | `[roster.<id>] adapter` | Without it nothing else is interpretable; closes the catalog/membership ambiguity. |
| `MODEL` | **launch argv when a model flag is present, else declared** (see §3) | The user's #1 question; must be effective-or-flagged, never silently declared. |
| `MODEL-FAMILY` | derived (see §4) | User-requested; grouping key. |
| `MODEL-COMPANY` | derived (see §4) | User-requested; answers "how many independent vendors do we actually have" (§15.6 relevance). |
| `EFFORT` | configured `reasoning`, with provenance marker | User-requested; same declared/effective hazard as MODEL. |
| `SPEED` | configured `speed` | User-requested; already layered. |
| `AUTO` | `AutonomousEffective()` (fail-closed, `discover.go:138-140`) | User-requested; 1.39.0 already fixed this one — keep the pattern. |
| `INSTALLED` | discovery probe (`Found`) | A rostered agent that is not installed is operationally absent; cheap (already probed). |
| `SOURCE` | `spec.Sources["model"]` (winning layer path) | Makes drift visible in one glance; the layer is where you edit, so the table must name it. |

Deliberately excluded (stay in `agents list` diagnostics): `VERSION`, `SANDBOX`,
`APPROVAL`, `TIMEOUT`, `HOME`, `BACKEND` — machine-capability facts, not roster facts.
`DISPLAY-NAME` is dropped from the table: it is a lossy compression of `FAMILY`+`MODEL`+`EFFORT`
and is today the source of a self-contradicting row — `roster show --json` reports
`display_name: "claude_opus-4.8-1m_max"` next to `model: "claude-opus-5[1m]"`
(PRIMARY: command output), because `RenderDisplayName` prefers `ModelLabel`
(`internal/agents/naming.go:189-191`), the built-in label is `Opus 4.8 1m`
(`discover.go:227`), and the deck pins only `model`, not `model_label`
(PRIMARY: `parley-deck/agents.toml:14-20` has no `model_label`). The composite name
stays for TUI/artifact rendering; the canonical table shows components instead.

**Two commands stay two commands.** `roster show` = "who is on this deck's team"
(§2 members; a mapped-but-not-in-§2 entry like today's `opencode-1` would appear only
under a `--all` flag, clearly marked `not in §2`). `agents list` = "what can this
machine run" (full adapter catalog + launch diagnostics). The defect is not the
existence of both; it is that neither is documented as the answer to "show me the
roster" and the hidden one is the roster. Fix: unhide `roster`, document both, and
have the skill name `roster show` as the only roster answer.

### 2. Update — `parley roster set`

```
parley roster set <roster-id> [--model M] [--effort E] [--speed S] [--timeout-ms N]
                  [--scope deck|machine] [--dry-run] [--yes]
```

- Resolves roster-id → family, then writes `[agents.<family>]` fields into the target
  layer file only (`deck` → `parley-deck/agents.toml`, `machine` → `~/.parley/agents.toml`).
- Safety properties, all reusable from existing code: atomic temp+rename
  (`fsutil.WriteFileAtomic`, already used by `writeRosterMappings`,
  `internal/app/roster.go:380`), whole-candidate TOML validation before replace
  (`config.ValidateAgentsConfigBytes`, `roster.go:377`), re-read-under-write guard,
  effort validated against `EffortVocabulary` (`naming.go:32`), unknown family in the
  target scope refused (reuse `MachineFamilyCatalog` for machine scope,
  `runtime.go:253-277`), idempotent (equal value → no-op), and a field-level
  before→after report naming the previous winning source. Default is preview +
  refusal, same pattern as `roster init` (`roster.go:325-328`).
- Rename the scope values to `deck|machine`; keep `session` as a hidden alias. The
  current name describes a store that does not exist (evidence above).

### 3. Fix the declared-vs-effective model split at the root

Today `applyOverride` sets `spec.Model` and never touches `HeadlessArgs`
(`internal/config/runtime.go:594-597`), while several built-ins embed a model literal
inside `HeadlessArgs` (claude: `--model claude-opus-4-8[1m]`, `discover.go:219`).
Result, reproduced from the central layer alone (`parley agents list --dir /tmp/parley-probe`,
PRIMARY): MODEL says `claude-opus-5[1m]` (source `~/.parley/agents.toml`), the launch
runs `--model claude-opus-4-8[1m]`. The deck pin is silently ineffective at launch —
the exact disease 1.39.0 cured for AUTO, one field over. Codex and kimi carry no model
flag in argv at all (`discover.go:196`, `:323`), so their MODEL column is pure
declaration about the CLI's own config file.

Proposal:

- Give each spec a per-family model-flag descriptor (e.g. `ModelFlag{Args: ["--model"],
  Style: "value"}`; codex `-c model="…"`, opencode `-m`, hermes `--model`, agy `--model`,
  kimi none) and **build the argv from `Spec.Model` at launch**, stripping embedded
  model literals from built-in `HeadlessArgs`. One source of truth; config overrides
  then genuinely pin the launch. Trade-off: this changes real launch behavior the day
  it ships (decks that pinned `claude-opus-5[1m]` will suddenly *get* 5 instead of
  silent 4.8) — it must ship with a changelog entry and a preflight probe, not quietly.
- Until a family has a descriptor, the table must not lie: MODEL renders the declared
  value with a `†` marker ("declared; not present in argv — effective value lives in
  the CLI's own config"), exactly mirroring the existing AUTO warning pattern
  (`discover.go:538-541`). When declared ≠ argv-carried, render
  `claude-opus-4-8[1m] ⚠ declared claude-opus-5[1m]` — effective first, divergence
  flagged, per the idea's constraint.
- EFFORT gets the same honesty treatment, but note the limit: for kimi, hermes, and
  opencode the effort lives in each CLI's own config, outside parley's reach
  (PRIMARY: `~/.kimi-code/config.toml:5` `effort = "max"`;
  `~/.hermes/config.yaml:29` `reasoning_effort: high`;
  `~/.config/opencode/opencode.jsonc:18` `"reasoningEffort": "xhigh"`). Parley cannot
  compute effective effort there without reading other programs' config formats. v1:
  show configured effort with a provenance marker; defer live probing to
  `agents verify --full`. Do not pretend otherwise.

### 4. Deriving `MODEL-FAMILY` / `MODEL-COMPANY`

Derivation lives in the CLI (one built-in table, versioned with releases — satisfying
"derived, not hand-maintained per deck"):

1. Strip a known **gateway prefix** (`litellm`, `openrouter`, …) and remember it as the
   route. `litellm/xai/grok-4.5` → route `litellm`, remainder `xai/grok-4.5`. The
   gateway is transport, never the company.
2. If a vendor namespace remains (`xai/grok-4.5`), map it via the company table
   (`xai` → xAI) and take the model series as family (`grok`).
3. Else map by model-id prefix against the same table: `glm-5p2` → Zhipu AI / `glm`
   (hermes' id has no namespace — PRIMARY: deck pin `parley-deck/agents.toml:62`);
   `claude-opus-5` → Anthropic / `claude-opus`; `gpt-5.6-sol` → OpenAI / `gpt`;
   `kimi-code/k3` → Moonshot AI / `kimi`; `Gemini 3.6 Flash (High)` → Google /
   `gemini` (label path, already sanitized by `SanitizeSection`).
4. Unknown → literal `unknown` plus a note line. Fail-closed, never a silent guess —
  same posture as `naming.go`'s reject-don't-repair design. Optional
  `[agents.x] model_company = …` escape hatch, tracked in `Sources`.

Optionally surface the route as ` (via litellm)` on MODEL-COMPANY rather than a column.
Trade-off: the built-in table needs a CLI release for new vendors; the override hatch
and the loud `unknown` keep that from being a blocker.

### 5. Sync — `parley roster sync` (deck ← machine), the single defined way

Mechanically, sync is: for every active §2 member, materialize its effective
machine-layer values into `parley-deck/agents.toml` — the `[roster.<id>]` mapping plus
chosen `[agents.<family>]` pins (model, effort, speed) — so the deck is self-contained
and immune to silent central edits. Semantics:

- **Additive and source-aware.** A field whose winning source is already the deck is
  never touched (the `Sources` map already tracks this per field, `runtime.go:596`);
  only machine-sourced values are pinned down.
- **Never destructive by default.** Removals/breaking changes require `--prune --yes`
  and are listed first; honor `roster_change_policy = "confirm-breaking"`
  (`~/.parley/agents.toml:18`, loaded via `LoadDefaults`, `runtime.go:172-194`).
- Same write safety as `roster set`; same preview-then-`--yes` gate; idempotent.
- `roster init` stays as the minimal mapping bootstrap; document it as the subset of
  sync that only writes missing `[roster.*]` mappings. (Merging init into
  `sync --mappings-only` is viable; I keep them separate for compatibility — open for
  round 2.)
- Detection, not automation: `parley preflight` gains a roster-drift section (unmapped
  IDs, machine-sourced values the deck has not pinned, declared≠effective models) that
  *names* `roster sync` / `roster set`. The skill's session-start habit becomes
  `parley roster sync --dry-run` — one defined operation, not a per-session invention.

### 6. Why adding `opencode` left sessions inconsistent — and the fix

Timeline, from file evidence:

- On/before 2026-07-30, the central hermes override lost `--yolo`:
  `~/.parley/agents.toml.bak-2026-07-30:43` already reads
  `headless_args = ["--oneshot", "{prompt}", "--model", "glm-5p2", "--accept-hooks"]`
  (PRIMARY; every later backup — `.bak-2026-08-05`, `.bak-2026-08-05-pre-opencode`,
  `.bak-2026-08-06-pre-autofix` — has the same stripped form; the strip's actual date
  precedes my earliest snapshot, RECALL beyond "on or before 2026-07-30"). A
  `headless_args` override replaces the built-in wholesale (`runtime.go:542-545`), so
  omitting one flag strips it. The deck override had lost `--yolo` the same way until
  today (PRIMARY: `git log -p` — commit `087a32a`, 2026-08-06, adds `--yolo` to the
  deck's hermes `headless_args`).
- Meanwhile the skill's own store kept `--yolo` for hermes-1
  (PRIMARY: `meta/headless-agents.local.json:104-111`). So two launch truths coexisted:
  a session launched through the CLI config stack got hermes **without** `--yolo`
  (headless hermes then blocks on permission prompts — effectively broken), while a
  session driven from the skill's JSON store got `--yolo` and worked. That per-path
  divergence is the mechanism behind "hermes works in some sessions, not others":
  **which store the session was launched from decided the outcome**. I did not inspect
  historical session logs, so attributing any specific past failure to this mechanism
  is RECALL; the store divergence and the argv difference themselves are PRIMARY.
- 2026-08-05, opencode was added **by hand, to the central file only**
  (PRIMARY: `~/.parley/agents.toml:84-112` `[agents.opencode]` and `:131-132`
  `[roster.opencode-1]`; the file's own comment says "added by hand 2026-08-05").
  Nothing updated any deck §2, any deck `agents.toml` (PRIMARY: this repo's
  `parley-deck/agents.toml`, read in full — no `[roster.*]` at all), or the skill
  store (no opencode entry). Hence: adapter catalog says yes, roster membership says
  nothing, skill store says nothing — three answers.
- The only written "procedure" for this was a hand-maintained comment that had become a
  stale trap: `~/.parley/agents.toml:118-122` still warns "DO NOT run
  `parley roster init` — kimi is AUTO=no … it would DROP `[roster.kimi-1]` and re-add
  `[roster.antigravity-1]`". On 1.39.0 init only ever appends missing mappings and has
  no delete path (PRIMARY: `roster.go:261-299`, `writeRosterMappings` `roster.go:339-381`);
  the deck config comment documents a reproduction showing the predicted drop did not
  occur even on 1.38.0 (SECONDARY: `parley-deck/agents.toml:29-38`).
- The repair landed today with 1.39.0's `kimi-opencode-full-adapters` idea:
  `--yolo` restored in both TOML layers and `--auto` added for opencode
  (PRIMARY: `parley-deck/ideas/kimi-opencode-full-adapters/IMPLEMENTATION.md:110-112`;
  diff of `.bak-2026-08-06-pre-autofix` vs current central file shows exactly those two
  changes; post-fix probe in an empty dir resolves hermes from central **with**
  `--yolo`, AUTO=yes).

What remains to do (the standard this idea should ratify):

1. **Membership is a protocol act.** `opencode-1` joins via a §2 roster-update
   (COOPERATION.md §2 path), then `roster sync` writes the deck mapping. Provide
   `parley roster add <family> [--as <id>]` to write the `[roster.<id>]` mapping and
   print the exact §2 row to ratify — the CLI must never regex-edit COOPERATION.md
   itself (quorum consequences; the parser is already a regex,
   `protocol/roster.go:17-20`, and writing prose by regex is how we got here).
2. **Delete the stale trap** in `~/.parley/agents.toml:118-122` and replace it with a
   pointer to `roster set` / `roster sync`. Machine-local edit, no protocol weight.
3. **Demote `meta/headless-agents.local.json`** from launch truth to (at most) a
   derived cache written *by* the CLI, or delete it from the skill flow entirely.
   Its divergence from the TOML stack is one half of the hermes mechanism.

### 7. Skill/CLI boundary

- **CLI owns**: the column contract and its stability, family/company derivation, the
  layering and all writes (validated, atomic), the effective-vs-declared computation
  and markers, drift detection in preflight.
- **Skill owns**: *when* to call. Session start: `parley roster sync --dry-run` (report
  only) and answer "what is the roster" exclusively with verbatim `parley roster show`
  output — the constraint "the skill must not describe a second format" made literal.
  Bootstrap picks are persisted through `roster set`, not hand-edited files.
- **Fix the skill's current over-promise**: SKILL.md:257 says
  `parley roster init` "lets model/effort be pinned per agent" — it writes only
  adapter mappings (PRIMARY: `roster.go:261-299`). Its layer description
  (SKILL.md:143-150) also names `headless-agents.local.json` as a per-project override
  the CLI does not read. Both must be corrected in the same release as the CLI change
  (the constraint exists; the channel exists — `parley version --all` already reports
  `compatibility: warning`, PRIMARY).

### 8. Backwards compatibility

No existing deck needs migration: layering is untouched, sync is opt-in and additive,
`--scope session` keeps working as an alias, the JSON store is already invisible to
the CLI so demoting it changes nothing mechanically. The only behavior change is
model-injection making pinned models actually launch — changelogged, probed by
preflight, and it is the bug fix this idea exists to deliver.

## Concerns / open questions

- **§2 as the membership store.** The canonical roster answer is built by regex-parsing
  a prose table (`protocol/roster.go`), and `resolveRoster` currently ignores the
  `inactive` set it asks for (`roster.go:110` discards it — an inactive row would
  still render as a normal row; PRIMARY by code read, behavior not exercised). Options:
  keep §2 canonical and have `roster show` mark/exclude inactive rows; or move
  membership to TOML and generate §2. The first preserves the human-ratification gate;
  the second breaks "no new config file / no protocol change". I favor the first.
- **`roster init` vs `roster sync` overlap** — fold or keep, round 2.
- **Sync freezes values.** A synced deck stops inheriting central improvements (the
  flip side of immunity). The preflight drift report must show *both* directions, or
  we trade silent breakage for silent staleness. Real trade-off, not free.
- **EFFORT provenance markers** may read as noise to users who just want the table;
  the marker legend must be one line, and `--json` should carry provenance as fields,
  not glyphs.
- **Built-in family/company table** needs a release per new vendor; the override hatch
  covers the gap, but the table's maintenance ownership should be named in the skill
  docs.
- The `/tmp` probe also surfaced that `agy` — removed centrally by user instruction —
  still lists as installed via its built-in spec. Harmless for `roster show` (not in
  §2), but `agents list` has no notion of "removed"; worth one sentence in its help
  text rather than new mechanics.

## Risks

- **Model injection is a live behavior change.** Decks pinning `claude-opus-5[1m]`
  will start actually launching it; if a backend rejects the pin, failures replace
  silent downgrade. Mitigation: changelog, preflight probe, and the ⚠ marker visible
  one release *before* injection ships, so users see the divergence first.
- **Regex-adjacent writes.** All new writes go to TOML with parse-validate-atomic
  replace; the one prose document (§2) stays print-only for the CLI. Residual risk:
  facilitators editing §2 by hand anyway — the skill text must route them through
  `roster add` + ratification.
- **Scope rename** breaks muscle memory/scripts; hidden alias `session` mitigates.
- **Skill/CLI skew** re-creates this exact idea next quarter; the
  ships-in-the-same-release constraint plus the existing `version --all`
  compatibility warning is the mitigation — but someone must treat that warning as
  blocking, which today is convention, not enforcement.
- **Claim hygiene.** The load-bearing facts above are tagged PRIMARY with locators;
  the one inference (specific past session failures caused by the hermes store
  divergence) is marked RECALL and should not enter `consensus.md` as established
  without a non-owner verdict. I issue no verdicts on my own claims here.
