# Changelog

## 1.45.0 — 2026-08-19

### Added — the `zcode` adapter

`zcode` is a full roster adapter. Headless launch is
`zcode --prompt=<text> --mode yolo --cwd <root>`, verified against the real binary
(`zcode-app-cli 3.7.7-13` / `zcode-runtime 0.16.3`) by `parley agents verify --full`.

Two shapes in that command line are not stylistic:

- **The equals form is required.** `zcode --prompt "-leading dash"` fails with *"Option '--prompt'
  argument is ambiguous"*, so any prompt whose first character is a dash would be lost. Placeholder
  substitution happens inside the argv element, with no shell involved, so a prompt containing
  quotes, newlines or flag-lookalikes still arrives as exactly one argument.
- **`--mode yolo` is the autonomous-write mode**, and it is *not* a sandbox. `--cwd` is a working
  directory, not an enforced boundary, so the adapter declares an empty `Scope` rather than
  claiming `workspace`.

**zcode has no `--model` and no effort flag** (`--model`, `--settings` and `--max-turns` all exit
1). A new `NoModelBinding` spec bit makes that fail-closed **at the source**: `ResolveLaunchArgs`
strips a config-supplied `--model`/effort flag for such adapters, so `roster show`, `agents list`
and the resolved argv cannot disagree, and a flag the CLI rejects can never reach the launch.

### Added — `MODEL`/`EFFORT` read from the agent's own config

Three roster rows reported `unknown` for the same structural reason: the CLI has no flag for the
value, so no parley layer can bind one and the process reads its **own** config instead. Those
rows now report what that file says, under two new `STATUS` terms — `model-from-config` and
`effort-from-config`, never plain `ok`.

| Adapter | Read from | Resolves |
| --- | --- | --- |
| `zcode` | `~/.zcode/cli/config.json` → `model.main`; `~/.zcode/v2/config.json` → `reasoning.defaultVariant` | model + effort |
| `kimi` | `~/.kimi-code/config.toml` → `[thinking] effort` (only when enabled) | effort |
| `opencode` | `opencode.jsonc` → `provider.*.models.<bound-model>.options.reasoningEffort` | effort |

**This is not a loosening of the roster contract.** That contract stops a *parley-side*
declaration being shown as if the argv carried it; this reads the file the agent itself reads at
launch, which is a different source and the actual one. A `model =` written into `agents.toml` for
an adapter that cannot carry it still never reaches the cell — pinned by test. `--explain` names
the file it read and states the limitation: the file can change before launch and none of these
CLIs echo the model back, so the value is not confirmable after a run.

`kimi` and `opencode` bind their model in argv, so the resolver deliberately returns no model for
them — a config `default_model` the launch overrides must not appear.

### Fixed — three `unknown` cells that were not unknowable

- **`hermes` effort was never passed.** It is now bound in argv (`--reasoning {effort}`) instead of
  inherited silently from `config.yaml`.
- **codex effort read as `unknown`.** `EffectiveEffort` did not recognise codex's own
  `-c model_reasoning_effort=<level>` form.
- **Model metadata gaps**: `fireworks` as a gateway prefix, `inkling` as Thinking Machines Lab.

## 1.44.0 — 2026-08-12

### Fixed — the fix-up budget was counted in the wrong unit, and off by one

Two defects in the same guard, both found by review of this release candidate.

**Off by one.** `internal/driver/impl.go` bounded the loop with `cycle >= MaxFixupCycles`, so the
last allowed cycle never ran: `standard` published 1 cycle where its §4.0 cell prints "cap 2", and
`deliberation` published 2 against a silent driver default of 3. The bound is now inclusive —
`MaxFixupCycles = N` permits attempts 1..N and escalates when N+1 would start.

**Wrong unit, and then a wrong source.** The cycle number came from the review-round ordinal, but a
strict-gate round with zero agreed fixes opens the next review round without publishing a fix-up —
so rounds that produced nothing spent budget. **The unit is now a reserved fix-up attempt**: a cycle
is charged when it is reserved, before the code-writing call, so an attempt that errors or is
interrupted has still spent it. It is deliberately not "a completed cycle" — that definition let a
failing fix-up loop forever against a cap that never depleted.

Getting the *source* right took two further review rounds. Counting `## Fix-up cycle N` headings in
`IMPLEMENTATION.md` was fail-open in four ways at once (a missing file restored the whole budget; a
heading inside a code fence counted; malformed and duplicate ordinals counted; renaming one past
heading bought a cycle) — and that file is owned by the implementer the cap constrains. Counting the
driver's own `.fixup-done` markers alone only moved the editable state.

The budget is now the **maximum of two driver-authored records**: a monotonic counter in the run
cursor, carried across rebuilds, and the `.fixup-done` markers. **The cycle is reserved before the
code-writing call**, so a fix-up that errors or crashes mid-flight cannot get its cycle back — a
fix-up that breaks the build is exactly the churn the cap exists to interrupt. Once both records exist, losing
either one does not lower the count; forging either can only raise it (escalating sooner); an
unreadable cursor escalates instead of counting as zero; and the crash-recovery path consults the cap
before opening another round, while still being allowed to *finish* the last allowed cycle.

**There is one window with a single record.** Between reserving a cycle and the marker being written
— which includes every attempt that errored — the cursor is the only record, and losing it there
loses that count.

**What this is not.** These are ordinary files in the repository. The budget is robust against a
stale or deleted single record once both exist, and against an errored or crashed fix-up losing its
reservation. It is **not** a
security boundary: a participant with workspace write, a deleted run directory, a repository
rollback, or two concurrent runs of the same idea can still reduce or duplicate the count. Making
the ceiling idea-scoped, serialized across runs and anchored outside the participant-writable tree
is a named follow-up, not something this release claims.

### Added — the two `deliberation` budget cells are explicit and enforced

Ratified by `meta-protocol-change-phase-packet-and-fixup-budget` (§7, four cross-review rounds,
accepted by codex-1, hermes-1 and kimi-1).

- **Fix-up (Phase 8): 5 inclusive reserved cycles**, was the printed word "unbounded" against a
  silent driver default of 3.
- **Cross-review (Phase 2): capped at 3 rounds after round 1**, was "unbounded" against a silent
  default of 1.

Both are escalation thresholds. **Hitting one never marks work complete.** There is no severity
floor: the deck's two worst ideas produced fresh MAJORs at rounds 19–24, so "late findings are
trivial" is false here.

**The cross-review cap binds every path that opens a round**, including the consensus-BLOCK
back-edge. This also bounds `standard`'s back-edge at its own printed cap of 2, which it previously
ignored — the two tracks share the mechanism, and leaving one unbounded while bounding the other
would be incoherent. Clamping only the initially scheduled budget left that back-edge governed by `MaxRounds`
alone, so a deliberation idea capped at 3 could still run a 4th cross-review round — a cap one code
path ignores is exactly the class this release exists to close.

**Why 5.** Counting `## Fix-up cycle` headings across all 69 ideas gives
`0×17, 1×34, 2×7, 3×2, 4×3, 5×2, then 9, 14, 15, 25`. Every value above 5 is in {9, 14, 15, 25} —
nothing has ever closed in the 6–8 band, so 5, 6, 7 and 8 escalate an identical set. No evidence
separates them, so the choice fell to error asymmetry: a too-low cap costs one recorded escalation a
human can grant; a too-high cap costs another cycle of the pathology the cap exists to interrupt.

### Changed — protocol text, all three copies

`parley-deck/COOPERATION.md`, `internal/protocol/defaults/COOPERATION.md` and the
`parley-deck-skill` bundled snapshot now print `cap 5 cycles` and
`capped at 3 after round 1, then escalate`. **Text and code land in one patch**: a `standard` idea
was measured running 15 cycles against a printed cap of 2 with no recorded escalation, because a
printed cap binds only where enforcement lives.

### Not in this release

- **The phase-scoped protocol packet** — the half of that idea that reduces read cost — is **not
  started**. Its release is gated on a pre-registered experiment whose ship/refute thresholds are
  already written into the idea's `FINAL.md` and may not be changed after data exists. **This
  release makes Parley Deck correct, not faster.**
- **The ratified escalation payload** (trajectory, findings by severity, fresh-vs-relitigated,
  unresolved fixes, validation status, recommendation) is **not implemented**. Both boundaries
  escalate and halt with the counts and the cap they enforced; the structured payload is a named
  remaining piece of the same idea.
- **`fast` Phase 8 is not exercised by the driver at all** — `fast` forbids idea-level
  `auto_implement`, so its fix-up route is manual. The inclusive-boundary fix corrects the
  arithmetic for every cap value, but no claim is made here about end-to-end `fast` behaviour.

## 1.43.1 — 2026-08-11

### Removed — the dormant frontier machinery

Adopts @codex-1's counter-proposal in full. `internal/runner/frontier.go` and its tests are deleted
and `runner.go` / `phase58.go` are restored to their pre-idea form.

**Why.** 1.43.0 shipped context-compaction machinery disabled by a constant. @codex-1 signed the
review consensus RESERVED over exactly this and its objection was never withdrawn:

> "Keeping unreachable safety code behind a constant invites exactly the rot the tests claim to
> prevent and gives a later one-line enablement change unjustified confidence."

It was right. A constant-false branch is executed by no test, so "compiled" was not "verified", and
its guards had to be asserted by matching source text rather than behaviour. 1.43.0 also still
perturbed prompts — the instruction wording and the `_ledger.md` exclusion — while delivering no
speedup, which is the worst of both.

**Nothing of value is lost.** The measured diagnosis, the located code paths, the signed
carry-forward ledger contract and the enablement gate all live in
`parley-deck/ideas/protocol-read-cost-regression/`. They are the deliverable; the inert code was not.

### Process note

1.43.0 was released after a MIXED round-3 review verdict (codex-1 NOT CLEAN, hermes-1 CLEAN, kimi-1
CLEAN) on the strength of a RESERVED consensus signoff. The release gate the owner had set was
"round 3 returns CLEAN", and it did not. This release closes the objection that gate existed to
catch.

## 1.43.0 — 2026-08-11

### Why

A measured investigation into why Parley Deck felt slower over recent versions
(`ideas/protocol-read-cost-regression`, 2 design rounds + 3 review rounds + 3 fix-up cycles).

**It is not the CLI** — every command runs under a second. It is the cost of a round multiplied by
the number of rounds, and both grew:

- reading `COOPERATION.md` in full costs **3.3x** median wall clock (n=3/arm, same agent, same task);
- the protocol grew **720 -> 1,359 lines** in ten weeks, monotonic, with `MUST` 15 -> 37;
- across 76 ideas split at 2026-07-01, **review rounds 1.6 -> 5.1** (max 24) and **review bytes
  20,237 -> 146,290 (7.2x)**, while design rounds stayed flat.

Both quadratic context paths were located: `gatherPriorRounds` and `gatherReviewContext` each
re-send every prior artifact, and **the protocol never required that** — §4 Phase 2 asks only that
every active participant be addressed. The CLI was stricter than the protocol it implements.

### Added

- Frontier context selection (`internal/runner/frontier.go`) with a participant-authored
  carry-forward ledger contract, its fallback, and tests.
- `parley protocol overlay show|validate`; the deck overlay parser and the
  `parley.protocol-lock/v2` lock (`internal/protocolcore/overlay.go`, `lock.go`).

### Changed

- The cross-review prompt now **describes its own contents truthfully**, derived from the same
  constant that controls the behaviour so the two cannot drift apart.
- `_ledger.md` is excluded by every artifact walker, so it is never handed to an agent as if a
  participant had written it. **This is an active input change**, not a no-op.
- The deck lock verifies the core's bytes, not just its version label.

### NOT shipped — stated plainly

**No speedup ships.** `compactionEnabled` is a **constant set to false**: no file, environment
variable or configuration can enable compaction. The mechanism that produces the speedup could not
be shown safe across three review rounds — a marker-derived ledger was fail-open, and an objection
that silently leaves the context becomes recorded consent under Phase 2's "Silence = implicit
agreement".

It may be enabled only once `protocol-ledger-validator` is **implemented** and the enabled path has
**end-to-end and mutation coverage for G3, G5 and G6**. A placeholder follow-up does not satisfy
that gate. Flipping the constant is a source change that goes through review.

The overlay is **partially implemented**: the roster-annotation identity slot (B6) and the removal
of prose-matched zone addressing (H9) are absent. The protocol text in all three copies now says so
instead of promising a future.

### Known dissent

@codex-1 signed the review consensus **RESERVED**, objecting to carrying unreachable machinery
behind a constant: "compiled" is not "verified", and whoever flips that constant will be flipping
code no test has run end-to-end. That objection is unwithdrawn and recorded in
`review/consensus.md`; @hermes-1 and @kimi-1 signed OK.

## 1.42.1 — 2026-08-07

Two defects found by the independent deploy verification of 1.42.0, which is exactly what that
step is for.

### Fixed

- **`parley protocol --help` exited 2** with an unknown-subcommand error, so the documented way to
  discover the new command failed. It now prints the group usage and succeeds.
- This repository's own `parley-deck/meta/version.json` was left at the previous skill version, so
  the tool reported its own deck as stale.

## 1.42.0 — 2026-08-07

A global core protocol, and `parley protocol` to work with it. Designed and reviewed across two
design rounds, three consensus revisions and nine review rounds by claude-1, codex-1, hermes-1 and
kimi-1 (idea `meta-protocol-change-global-core-protocol`).

### Why

Measured across 36 decks: **eight different `deckVersion` values**, §15 present in 5 of them, the
§2 roster-authority change in 1. `COOPERATION.md` was copied into every deck and hand-editable
there, and the sync that repaired it was a one-off script, not a mechanism. Yet the drift was not
project customization: only **one** deck in 36 carried a genuine local section, and it was
governance about how the protocol is synced — content that belongs in the core.

### Added

- **A global core protocol store** at `~/.parley/protocol/core/<version>/`. Releases are
  **write-once**: a release is never edited in place, so a change by the user becomes a new version
  by construction rather than by discipline. Symlinked store components are refused on read and
  write; version strings are validated before they reach a path.
- **`parley protocol status|render|check|publish`.**
  - `render` regenerates the deck's `COOPERATION.md` from the core, preserving the six per-deck
    identity values (workspace, created, transport, the sync stamp, and both §2 tables). Preview is
    the default. It **reports what it will not carry forward**, in preview and on apply.
  - `check` reports a hand-edited or stale deck copy and **never overwrites** it, exiting non-zero
    so a script can act.
  - `publish` is attended-only: it refuses without a controlling terminal, which stops an ordinary
    agent run whose stdin is a pipe or `/dev/null`. It does not stop an agent that allocates a pty —
    the command says so itself rather than implying more.
- **A deck lock** (`parley-deck/meta/protocol-lock.yaml`). A machine that lacks the exact pinned
  release **blocks** rather than substituting its own same-named or current version.

### Changed

- **§7 gains a blast-radius clause** in all three protocol copies: a CORE change requires the meta
  idea and explicit user ratification; a DECK change is a normal idea. The section states plainly
  which parts of this design are **not yet in force** — per-idea pinning and the tamper signal are
  ratified but unimplemented — instead of describing an intended future as present fact.

### Known limits, stated rather than implied

`render`'s loss report is a **line-level diff**, not a Markdown semantic analysis. An empty report
means no line disappeared; it does not prove no meaning was lost. Read the diff before `--yes`. The
report was wrong in nine distinct ways during review before reaching this form, which is why its
claim is narrow.

### Not shipped (ratified, scheduled)

Per-idea protocol pinning, the deck overlay (local override/extension), the OS sandbox that would
make core-write prevention enforceable for parley-launched agents, and migration of the existing
36 decks. Nothing in this release claims any of them.

## 1.41.0 — 2026-08-06

Standardized roster operations, completed. 1.40.x shipped the surface; this release fixes what six
rounds of multi-agent review found underneath it. **The headline defect: a deck's declared
membership was not what ran.**

### Fixed — membership authority

- **A deck declaring N participants ran however many the machine configured.** Membership was read
  from the layered machine+deck view, so a deck listing two agents resolved to five whenever
  `~/.parley/agents.toml` listed five — and, since 1.40.1 routed participant selection through the
  same view, it would have *deliberated* with five. `roster render` then committed that inherited
  roster into `COOPERATION.md`, re-creating the drift this whole change exists to remove, inside a
  file collaborators share.
  Membership is now **the committed `parley-deck/agents.toml`**. The machine layer seeds values
  only.
- **A legacy §2 deck was overridden by the machine roster.** A deck that predates the cutover keeps
  its §2 table as its membership — reported `legacy-roster` — until it is migrated. Authority
  order: committed deck blocks → valid legacy §2 → machine roster (marked `inherited-roster`, and
  `roster render` refuses to commit it without `--adopt-inherited`).
- **A gitignored or machine-local file could retire a committed member.** `active` merged from
  every layer, so `active = false` in `agents.local.toml`, `$PARLEY_HEADLESS_AGENT_CONFIG` or the
  machine file silently dropped a member from the quorum. State now follows the layer that granted
  membership.
- **`§2`-only IDs are reported, never silently erased.** They surface as `unmapped` /
  `section2-only`, are never auto-added to the generated table, and `roster render` reports every
  row it removes — in preview and on apply.

### Fixed — run snapshots

- **Per-roster-ID pins collapsed.** The frozen map was keyed by adapter family, so two roster IDs
  sharing an adapter overwrote each other. Now keyed by agent, and applied *after* participant
  resolution — before, the freeze was computed against adapter-keyed discoveries and never reached
  the launch.
- **AUTO was reported but not pinned.** The snapshot now carries the resolved launch argv, so a
  config change cannot alter a running idea's autonomy posture, and `RosterRevisionOf` hashes it so
  the drift is detectable.
- **`sessions inspect` reports `stale-snapshot`** — the frozen vocabulary previously shipped a code
  no surface could emit.

### Fixed — the query surface

- **`--scope` parsed, was advertised, and did nothing.** It now selects the deck or machine roster,
  including values, adapter mappings and provenance.
- **`--all`** lists configured adapters no roster declares — the answer to "I installed an agent and
  it is invisible".
- **`--explain AGENT`** reports per-field provenance: which config layer set each value.
- **Text and JSON now agree.** A healthy row printed `ok` and marshalled to `null`; `display_name`
  and `note` were serialized outside the eleven frozen columns.

### Fixed — writes

- **The membership gate was bypassable.** `roster set new-9 --model X --yes` created a member
  without the second confirmation, because the gate keyed on `--adapter` rather than on whether the
  block existed. `roster init --yes` bypassed it entirely. Both fixed; the gate now fires on a real
  state change, not on writing a value a member already has.
- **`roster sync` accepted typoed `--keep` tokens** (`--keep kimi-1.modle` removed `kimi-1.model`),
  and applied against a second read without binding to the preview. Both are errors now.
- **`masked-by-env` has an emitter**: a write overridden by a higher layer is reported instead of
  returning a false success.
- **File mode regression**: writes no longer tighten a `0644` config to `0600`.

### Fixed — model arguments

- **D7's legacy normalizer shipped.** A config layer that hardcodes `--model <literal>` in
  `headless_args` no longer outranks the `model` field beside it; the literal is normalized back to
  the `{model}` placeholder.

### Fixed — protocol and docs

- Every protocol surface still instructed the reader to treat §2 as a store while calling it
  generated. Removed from all three `COOPERATION.md` copies and `SKILL.md`, with a drift assertion
  so the contradiction cannot silently return — it caught a fourth instance the moment it was added.
- `roster` gained full sections in `docs/cli-reference.md` and `docs/agent-runtime-configuration.md`
  (previously zero mentions), and all five verbs appear in `parley --help`.
- `agents list` is labelled the adapter/runtime inventory — not the roster.
- `roster migrate --yes` requires `--confirm-breaking` and skips decks with uncommitted changes.

## 1.40.1 — 2026-08-06

Fixes from the Phase 6 review of 1.40.0. **1.40.0 shipped before that review** — the user directed
the work to deploy — and the review found two CRITICALs and several MAJORs. Everything below is a
defect in 1.40.0, found by codex-1, hermes-1 and kimi-1 — all three corroborated both
CRITICALs independently.

### Fixed

- **The run snapshot was written and never read.** 1.40.0 froze each participant's launch identity
  at run creation, but `continue` still re-discovered configuration, so changing a machine default —
  or running `roster sync`, which exists to change deck values — could move a running idea onto a
  different model mid-deliberation. Continuations now consume the frozen row, with the acceptance
  test the ratified gate required. **Rebase had shipped without the half that made it safe.**
- **The authority cutover was half-done.** `roster show` read config while participant selection
  still parsed §2, so the table and the run could disagree about who was in it — the exact
  two-sources-of-truth defect this work exists to remove. All membership decisions now resolve
  through one function.
- **There was no §2 generator**, although 1.40.0 had already changed the protocol in three copies to
  describe §2 as generated. `parley roster render` generates it idempotently and preserves the
  workspace-dir and role values that only ever existed in the hand-written table.
- **`--yes` alone performed membership changes.** Adding or retiring a member alters who deliberates
  and therefore a future quorum; it now additionally requires `--confirm-breaking`.
- **Machine-scope writes went to a file nothing reads.** `PARLEY_HOME` names the central config
  *directory*, but the writer composed `$PARLEY_HOME/.parley/agents.toml` while the loader reads
  `$PARLEY_HOME/agents.toml`. A machine update reported success and changed nothing. The tests had
  encoded the wrong path and passed the defect.

### Added

- `parley roster migrate` — the fleet migration tool: per-deck inventory, dry-run by default,
  file-level backups, post-write validation with automatic rollback, skip-and-report on anything
  unclean, and a machine-readable final report.

## 1.40.0 — 2026-08-06

Standardize roster operations. "What is the current agent roster?" had no single answer: three CLI
surfaces produced three different tables, two independent stores held membership, and the table
reported a `MODEL` the launcher never passed.

Measured across 40 decks: **nine distinct rosters**, 17 with no roster at all, 17 still naming an
agent retired months earlier, and one deck missing a participant entirely — which is why the same
agent worked in some sessions and not others.

### Fixed

- **A configured model never reached the process.** The model lived in two places — `Spec.Model`
  AND a literal baked into `HeadlessArgs` — and config layers set the field without rewriting the
  args. Pinning a model changed only what was *displayed*: `claude` launched Opus 4.8 while every
  config layer said Opus 5, and **six of seven adapters passed no effort flag at all**. Built-in
  args now carry `{model}`/`{effort}` placeholders that the runner substitutes, so one value lives
  in one place. `codex`, `kimi` and `opencode` gained `-m`. An unbindable placeholder drops its
  introducing flag rather than leaving a value-taking flag dangling, which would abort the CLI.
- **`AUTO` is computed from the resolved argv**, not the raw one.
- **Retired agents rendered as full members.** `resolveRoster` discarded the inactive set, so
  marking a row inactive did nothing. `STATE` is now wired.
- **`roster` was dispatched but absent from `parley --help` and the docs** — the command the skill
  tells agents to run was undiscoverable.

### Added

- **A frozen 11-column roster contract**, identical in text and `--json`, carrying
  `schema_version` and the ordered column list so consumers can detect a contract change:
  `AGENT ADAPTER STATE INSTALLED MODEL MODEL-FAMILY MODEL-COMPANY EFFORT SPEED AUTO STATUS`.
  `MODEL` and `EFFORT` hold what the launch actually passes, or `unknown` — never a declaration
  wearing the effective cell. Divergence surfaces as `STATUS=model-drift` / `effort-unknown`.
- **`modelmeta`**, a CLI-owned derivation of model family and company. Gateway prefixes are peeled
  first, so `litellm/xai/grok-4.5` is **xAI via LiteLLM**, and an adapter never implies a company.
- **`parley roster set AGENT --scope deck|machine`** — change one member in one file. Preview by
  default. `--scope deck` writes the **committed** `parley-deck/agents.toml`, never the gitignored
  local file. `--state inactive` marks; rows are never deleted.
- **`parley roster sync`** — the single defined machine → deck reconciliation, with **rebase**
  semantics: redundant deck overrides are removed so the deck inherits, and a deliberate pin is
  never dropped silently — it is enumerated with the exact `--keep AGENT.FIELD` that retains it.
- **An immutable per-run roster snapshot** plus `roster_revision`. Runs previously recorded
  participant IDs and nothing else, and `continue` re-discovers configuration — so changing a
  machine default mid-run could silently continue it on a different model.

### Changed — protocol

- **`COOPERATION.md` §2 is now a generated, non-authoritative view**; `parley-deck/agents.toml`
  owns the roster. Decks with only the old hand-written table keep working and report
  `legacy-roster` until `parley roster sync` moves them across. See
  `parley-deck/meta/protocol-changelog.md` for the recorded, user-authorized one-off venue
  deviation from §7.

Designed and reviewed through a full Parley Deck run (`roster-operations-standard`, track
`deliberation`, four participants, two rounds, three consensus revisions). codex-1 blocked twice;
both blocks were upheld and discharged. Fifteen drafter position changes are recorded.

## v1.39.0 - 2026-08-06

**`kimi` and `opencode` are full adapters, and a silent auto-approve defect is fixed.**

Both agents existed only as one-line entries in the ACP catalog. `specFromACPBackend()` never sets
`AutonomousWrite`, so `agents list` reported `AUTO=no` for them — not a statement about the CLIs,
which both write unattended, but about parley not knowing which flag enables it.

- **Full built-in specs.** `kimi` launches as `kimi -p <prompt>`, `opencode` as
  `opencode run --auto <prompt>`. Both keep `kimi acp` / `opencode acp` as an alternative launch
  mode. Every field was probed live before being written: `kimi --auto -p …` exits 1 with
  *"Cannot combine --prompt with --auto"*, so `-p` is the only autonomous headless shape kimi has.
- **`Scope` stays empty for both.** Only `codex --sandbox workspace-write` enforces a real sandbox;
  the type forbids claiming confinement that is not enforced.

**The defect, found in review by codex-1 and wider than the change that exposed it.** A config
layer replaces `HeadlessArgs` wholesale without touching `AutonomousWrite`, so parley could declare
an autonomous mode whose enabling flag the launched command never passed — and still print
`AUTO=yes`. **`hermes` was already in that state**: its override had dropped `--yolo`.

- `AUTO` now **fails closed**: a declared mode whose enabling args are absent from the effective
  launch reports `no`, with a warning naming the missing args.
- The `headless:` line in `agents list` now shows the **resolved binary and effective argv**
  instead of the built-in label. The label is what hid the defect.

Also: kimi's notes now record that its installer does not add the binary to `PATH`, so `command`
must be set; opencode's telemetry description no longer understates its streamed output.

## v1.38.0 - 2026-08-04

**Protocol: new `§15 Verification integrity`.** Ratified by idea
`meta-protocol-change-verification-integrity` (two design rounds, four consensus revisions, four
signoff rounds, two review rounds, two fix-up cycles).

The protocol had strong rules about who *writes* which artifact and no rules about what makes a
*verification* valid. A participant could stamp `CONFIRMED` on any claim — including its own —
with no stated basis, and two contradictory verdicts had nowhere to live and no resolution rule.

- **§15.1 — scope, ownership, location.** A claim enters the regime only when someone verdicts it,
  someone challenges it, or §15 requires it. **Every participant that asserts a claim as true where
  it first appears canonically is an owner, and an owner MUST NOT verdict a claim it owns.**
  Material transcribed and explicitly marked as unverified testimony is not owned by the
  transcriber — without that branch a facilitator could never verify anything it put in a brief.
- **§15.2 — provenance.** `PRIMARY` (source located and quoted, **or a check the verifier executed
  with command, inputs and output quoted**) / `SECONDARY` (a *named* participant's non-`RECALL`
  verdict, chain acyclic and terminating in `PRIMARY`) / `RECALL` (caps at `UNVERIFIED`).
  **An untagged verdict is treated as `RECALL`** — the scheme fails closed.
- **§15.3 — conflicting verdicts.** Resolved by reviewable evidence and argument, **never by
  counting participants, including where the count is unanimous.** Provenance controls
  admissibility; it does not select the winner. Unresolved → `DISPUTED`, which may not support any
  acceptance criterion. No new file.
- **§15.4 — exemption-claim admissibility.** A claim to avoid a named obstacle needs a witness
  logically sufficient for the scoped claim. Adjectives are not witnesses.
- **§15.5 — role concentration.** Facilitator procedural calls are provisional until the signoff
  gate passes. A facilitator-drafter must publish `## Drafter position changes` with an exact prior
  quotation and source path per change.
- **§15.6 — correlated agreement.** On unanimous judgment-shaped ideas, consensus may not close
  without a steelman of the strongest alternative; a null result recording the search scope is a
  finding, not non-compliance. `deliberation` takes an assigned round artifact, `standard` a
  section inside an existing round-02 file.

**Two text fixes.** §4.0 listed round-1 independence among invariants "never dropped for speed"
while §11.A said "there is no enforcement beyond agent discipline" — the qualifier reconciles them.
§6 rule 4 now applies explicitly to scoping.

Applied to **both** `COOPERATION.md` copies (live deck and the embedded `parley init` template),
per the drift guard.

**Recorded finding about §15.5 itself.** Across four consensus revisions the drafter's own
disclosure went 8 → 13 → 21 → 23 of 23 material changes; every increment came from other
participants re-running the source comparison. §15.5 is not self-enforcing, and `FINAL.md` carries
that as an open follow-up rather than pretending otherwise.

## v1.37.0 - 2026-07-30

**Roster: `agy` is a participant again, on Gemini 3.6 Flash (High).**

- `antigravity-1` was `inactive` from 2026-07-18, when `kimi-1` took the fourth roster slot. It
  returns as a **fifth participant alongside `kimi-1`**, not as a replacement.
- The built-in `agy` spec is pinned to **`Gemini 3.6 Flash (High)`** — a generation newer than
  the `3.5 Flash (High)` it replaces. `agy models` refuses to list in a headless context
  (`Please sign in`); passing an invalid `--model` prints the valid list instead, which is how
  the current set was read: 3.6/3.5 Flash Low/Med/High and 3.1 Pro Low/High.
- The headless regression that had kept `agy` out of headless rounds — no artifact emitted for
  non-trivial `--print` prompts on 1.0.4 — is **gone on 1.1.8**, verified with a read-then-write
  prompt that produced an accurate description of code it had to open to describe.

The pin lives in four places and they must not drift: `~/.parley/agents.toml`, the deck's
`parley-deck/agents.toml`, the machine-local `headless-agents.local.json`, and the built-in Go
spec. `TestDefaultSpecsPreferAntigravityAndStrongVerifiedDefaults` catches the last one.

**Known capacity limit.** `agy` returned `Individual quota reached … Resets in 158h59m` after
two full review rounds on a large repository. It is an account quota, not a per-model limit — a
one-word prompt fails identically once it trips — so there is no cheaper fallback and no point
retrying. Budget it as a scarce reviewer, and record its absence from a round as an outage
rather than as an accept.

## v1.36.0 - 2026-07-18

**composite-agent-naming-and-roster-reinit** — designed + implemented + reviewed via a real 5-agent Parley Deck run (claude/codex/hermes/kimi active, agy quota-out; unanimous ACCEPT, then a refutation-default code review + two fix-up cycles).

- **Self-documenting agent display names** `family_model_effort` (e.g. `codex_gpt-5.6-sol_xHigh`, `claude_opus-4.8-1m_max`, `agy_gemini-3.5-flash_high`) — `_` separates the three meanings, `-` separates words, `.` keeps versions; camelCase effort (`xHigh`/`cliDefault`). Derived at render from config, never an identity, path-safe & fail-closed to parse.
- **Two-namespace schism fixed**: `agents.Spec.ID` split into a stable roster identity + an `AdapterID` family for launch/vendor dispatch; a fail-closed participant resolver (exact spec-ID → explicit `[roster.*]` mapping → hard error, no prefix heuristic, participant-id grammar validated against path traversal). The roster ID is now the identity used for artifact paths/signoffs; the driver + every app path (run selection, preflight, consensus/FINAL drafter, request-signoffs, TUI steer) resolves a `[claude-1, …]` roster.
- **`parley roster show|init`**: `show` renders the resolved roster with composite names; `init` proposes + writes the roster-ID→family `[roster.*]` mapping (idempotent against the target file, `--scope session|machine`, `--dry-run/--yes/--json`, atomic write, fail-closed on an unresolved/typoed adapter).
- **Autonomous write is first-class** (`AutonomousWrite{Mode,Args,Scope}`, `AUTO` column) and the built-in defaults are now actually autonomous. **Posture change:** built-in `claude` moved `--permission-mode acceptEdits` → `bypassPermissions` (workspace-scoped via `--add-dir {root}`); `codex` `approval_policy="on-failure"` → `"never"` (still `--sandbox workspace-write`); `hermes` gained `--yolo`.
- **`fast` is the standard speed on a separate axis** — same model + same effort, faster output (never a downgrade); central template defaults `speed = "fast"`; a guard test locks speed↛model/effort.
- Skill (`parley-deck-skill`): an **Autonomous Execution (required)** section with the per-CLI yolo mapping (incl. kimi's `-p`) + the display-name / `parley roster init` convention.

## v1.35.0 - 2026-07-04

Five features inspired by Hermes Agent v0.18.0, each designed + implemented + reviewed via a real multi-agent Parley Deck run (claude/codex/hermes/antigravity):

- **completion-contracts-evidence-ledger** (protocol): `checks:` in 00-prompt.md now accepts an optional named list of `{name, command}` criteria (a completion contract). The driver runs each criterion, writes a secret-scrubbed/truncated per-criterion table into IMPLEMENTATION.md's `## Validation evidence` section (committed each cycle), and vetoes `status: complete` while any criterion fails at HEAD (fail-closed, independent of strict_gate). Scalar/absent `checks:` unchanged. Protocol: LE-4 + Phase-5 + Phase-8.
- **named-roster-presets**: `[rosters.<name>]` + `[defaults.track_rosters]` in agents.toml; `parley run --preset NAME --track T` expands a preset into the canonical participants (validated against the §2 roster, fail-closed); `parley preset list` shows presets, source layer, and stale-member warnings.
- **tui-round-summary**: the driver emits an idempotent `round.digest` event when a round completes; the TUI Home tab renders a bounded position-map digest (per-agent one-liners, keyword mention-flags, next action) that never pushes Recent runs off-screen.
- **tui-editor-composer**: `/editor` + ctrl+e open $VISUAL/$EDITOR/vi on a 0600 temp file and drop the (multi-line) result into the composer; the existing Enter path keeps steer/answer semantics.
- **parley-learn-playbooks** (protocol): `parley learn <closed-idea-slug>` distills a COMPLETED idea into an advisory playbook under `parley-deck/playbooks/` (fail-closed write boundary: parent-symlink refusal + O_EXCL). Protocol: §13.5 (advisory, beside consults).

Protocol changed (contracts + §13.5) — both COOPERATION.md copies + the skill fallback are in sync (drift guard green).

## v1.34.0 - 2026-07-03

**Protocol progressive-disclosure layout (pure reorder).** `COOPERATION.md` now reads
core-first, reference-last: §9 (session-start checklist) was relocated to sit after §10 (TL;DR),
so the reference sections (§9, §11, §12, §13, §14 + Appendix A) all follow the core (§0–§8, §10).
A **provably content-preserving move** — every section keeps its number (all `§N` cross-references
resolve), no rule text was added, removed, or changed (the sorted-line diff is empty). Designed +
reviewed by a real multi-agent Parley (`ideas/protocol-restructure-appendices`, deliberation
track, unanimous ✅ ×3). No `core ≤200 lines` compression (that needs a separate §4 phase-split);
no `## Appendices` banner (kept the change a zero-addition move).

## v1.33.0 - 2026-07-03

**Track-aware driver — deterministic §4.0 enforcement.** The follow-up to v1.32.0's
conditional-rigor protocol text: the CLI/driver now actually routes and gates by the declared
`track:`. Designed + reviewed by a real multi-agent Parley (`ideas/track-aware-driver`,
deliberation track, unanimous ✅ ×3 review over three rounds).

- **`parley classify [--files N --loc N --security …] [--declared T] [--json]`** — a pure,
  script-checkable §4.0 classifier (deliberation-first, fail-safe: unknown/negative size is never
  `fast`); `--declared` exits 4 on an under-tier so CI can gate.
- **`track:` enforcement in the driver** (new `internal/track` package + `driver.ReadTrack`):
  `fast` runs with 1 model-diverse-required reviewer, no cross-review rounds, and a 1-cycle
  fix-up cap; explicit `standard` caps reviewers at 2, cross-review at 2, and fix-up at 2;
  `deliberation` and an absent `track:` behave **exactly as before** (backward-compatible).
- **Hard-rejects (escalate, never silently proceed):** `fast` + `auto_implement`, `fast` +
  `strict_gate`, and any non-solo config (0 independent reviewers) on an explicit track — the
  contradiction check reads idea-level intent, not the `--no-implement`-masked runtime flag.
- Refutation-default review stays structural and non-optional on every track.
- Also fixes a pre-existing driver-lock TOCTOU in `acquireLock` (surfaced during review): an
  empty/just-created lock file is now treated as held, closing a two-concurrent-holders race.
- Deferred to follow-ups: per-track timeouts, fast §9.0 ping-skip, collapsed fast consensus/FINAL,
  per-phase human gates, mid-idea upgrade.

## v1.32.0 - 2026-07-03

**Conditional-rigor tracks + developer Quickstart (DevX & speed).** Designed and reviewed by a
real multi-agent Parley deliberation (`ideas/meta-protocol-change-devx-speed`, deliberation
track, unanimous ✅ ×4 design + ✅ ×3 review) to make the protocol usable without reading 1000
lines and much faster for ordinary work — without touching the safety core.

- **`track: fast | standard | deliberation`** in `00-prompt.md` (default `standard`), chosen by
  an objective, **fail-safe** classifier (§4.0): deliberation-first, then fast, else standard;
  on any doubt/boundary → the stricter track. `deliberation` is forced by protocol change,
  security/secrets/production, data migration/irreversible, `strict_gate`, `auto_implement`,
  pipeline/action, or public-API/schema break.
- **Per-track ceremony** (§4.0 table = the single authoritative gate, overriding the
  full-lifecycle defaults in §4/§5/§9.0/§11): `fast` = cross-review skipped, collapsed
  consensus/FINAL, 1 model-diverse reviewer, ≤1 fix-up, ~5-min timeouts; `standard` = 2
  reviewers, cross-review capped at 2, ~15-min; `deliberation` = today's full lifecycle, unchanged.
- **Invariants preserved on every track:** non-solo, refutation-default review, round-1
  independence, append-only signoffs, audit trail, §14 human brake, English-only, no-secrets.
- **DevX:** a top-of-doc Quickstart (5-minute start), a "Who are you?" role table, a
  core-vs-appendix reading guide, an off-ramp ("trivial reversible work needs no Parley"), and a
  consolidated plain-English `LE-N` glossary (§4.0.1).
- Additive change; both protocol copies stay byte-identical (drift guard green) and the skill
  fallback is re-synced. Deferred to ratified follow-ups: deterministic CLI/driver enforcement
  (`track-aware-driver`) and the physical appendix relocation (`protocol-restructure-appendices`).

## v1.31.0 - 2026-06-24

**Loop engineering (LE-1..11).** A four-tier program — designed by a real multi-agent
Parley deliberation (`ideas/loop-engineering-research`) and implemented tier-by-tier with
full refutation review — that turns Parley Deck into a loop-engineering substrate with a
human-gated consensus brake.

- **Tier 1 — verification honesty (LE-1/2/3/4).** Phase 6 review is refutation-default
  (reviewers assume the change is wrong until they fail to break it) and gains a
  `## Refutation attempts` section; `strict_gate` lets the driver block a close on a
  not-certified-clean round (deterministic finding-scan veto, bounded by `MaxFixupCycles`);
  optional reviewer model-diversity signal (`require_model_diversity`); a `checks:`
  frontmatter command the driver runs before an auto-implement close (fail-closed for a
  code-writing idea with nothing to check).
- **Tier 2 — loop budgets (LE-5/6/10).** The driver enforces a per-run loop budget
  (`--max-driver-steps`, `--max-wall-clock`, `MaxCostUSD`; central `[defaults.loop]` in
  `~/.parley/agents.toml`) and escalates on breach; cost telemetry per tick; the §12.11
  monitoring watcher opens **`status: candidate`** remediation ideas (no auto-staffed quorum).
- **Tier 3 — close-decision integrity (LE-7/11).** Under `auto_implement`, an
  `ACCEPT-WITH-RESERVATIONS` triage or fewer than two independent reviewers escalates instead
  of auto-completing; a fresh non-implementer **goal-done check** (advisory, fail-open, 2-min
  bounded) runs before close.
- **Tier 4 — the outer loop (LE-8/9).** New **COOPERATION.md §14 human brake**: any
  automated/scheduled loop may discover-and-draft Phase 0/1 candidates only — never promote,
  run, implement, land/merge, finalize, edit the roster, or override consensus without a
  recorded human or full-quorum gate. New **`parley loop tick`** command (`internal/loop`):
  one-shot, scheduler-friendly, disabled-by-default; it drafts `status: candidate` idea
  prompts from a signals file and dedupes them, and never runs/pushes/merges/finalizes/staffs
  a quorum. Hardened against frontmatter injection, dedupe-digest collision, poisoned-dir
  liveness, symlink escape (at any ancestor depth), and line-break separator tricks.

## v1.30.6 - 2026-06-20

- **Seeded `[defaults]` tuning.** The `~/.parley/agents.toml` template now
  defaults `preferred_transport = "local-dir"` (local files) and trims the
  default timeouts to `signoff_ms = 600000` (10 min) and `round_ms` /
  `review_ms` / `deep_reasoning_ms = 1200000` (20 min). Existing central files
  are untouched (the template only seeds a fresh `~/.parley/agents.toml`).

## v1.30.5 - 2026-06-20

- **Central `[defaults]` policy block in `~/.parley/agents.toml`.** Beyond the
  per-agent catalog, the central config (and a deck's `parley-deck/agents.toml`)
  now carries a `[defaults]` block, merged with the same low→high precedence:
  - `speed` — applied as the global default speed for every agent spec
    (`config.LoadAgentSpecs`); a per-agent override still wins.
  - `ping_tier` — `none`/`off` opts out of the §9.0 hosted-PONG round-trip in
    `parley preflight` / `parley run` (explicit `--no-ping` still forces skip).
  - `preferred_transport` — `parley init` seeds the fresh deck's transport from
    it (`local-dir`/`github-pr`/`gitlab-mr`; unknown → `local-dir`).
  - `roster_change_policy`, `timeouts` — exposed via `config.LoadDefaults` and
    honored by the facilitator/skill.
  New `config.LoadDefaults`, `protocol.InitWorkspaceWithTransport`. The seeded
  `~/.parley/agents.toml` template includes a commented `[defaults]` block.

## v1.30.4 - 2026-06-20

- **Central per-user agent defaults (`~/.parley/agents.toml`).** A new
  user-global config lists each agent's command, model, and reasoning/effort
  level, inherited by every project. Wired into `config.LoadAgentSpecs` as the
  lowest config-override layer: built-in defaults → `~/.parley/agents.toml`
  (central) → `parley-deck/agents.toml` / `agents.local.toml` (per-project
  override) → `$PARLEY_HEADLESS_AGENT_CONFIG`. A deck overrides the central
  default; fields the deck leaves unset fall through to the central value.
  `parley init` seeds a starter `~/.parley/agents.toml` (never clobbering an
  existing one) and prints where to override per-project. `PARLEY_HOME`
  overrides the central dir (used for hermetic tests).
- **Reasoning/effort is now part of the deck-bootstrap confirmation.** §0 and
  the skill's deck-bootstrap step confirm the roster, each agent's model **and
  each agent's reasoning/effort level**; the default reasoning/effort is the
  **strongest (highest) level the agent supports**, falling back to
  `cli-default` only when it cannot be discovered. Protocol stays model- and
  reasoning-agnostic.

## v1.30.3 - 2026-06-19

- **Fix: roster & model confirmation is a deck-BOOTSTRAP gate, not per-idea/per-session.**
  Corrects 1.30.2, which placed the mandatory roster + per-agent model confirmation in
  the per-idea §9.0 readiness check. It now lives in **§0 (deck bootstrap)**: the
  confirmation fires **once, when `parley-deck/` is first created (`parley init`)** —
  not per idea, not per later session; an already-bootstrapped deck reuses the saved
  selection. §9.0 keeps only the per-idea agent **liveness** ping (it no longer
  re-selects models). Both COOPERATION.md copies; drift-guard lockstep; protocol stays
  model-agnostic. (Skill side: `parley-deck-skill` 1.3.3 moves the interactive flow to
  "Transport Selection / deck bootstrap".)

## v1.30.2 - 2026-06-19

- **Mandatory session-start roster & model confirmation.** §9.0 now states that at a
  session's first readiness check the facilitator MUST confirm the active roster and
  each agent's selected model with the user before the first idea; the user's
  persistent per-agent model choice is recorded in `meta/headless-agents.local.json`
  and reused until changed (later sessions show the saved picks for explicit
  confirmation). The protocol stays model-agnostic — it mandates the confirmation, not
  any specific model. (Both COOPERATION.md copies; drift-guard lockstep. The detailed
  interactive list-roster → confirm → list-models → pick flow lives in the
  `parley-deck-skill` SKILL.md Startup Flow / Selection Checkpoint.)

## v1.30.1 - 2026-06-19

- **Pin the `claude` participant to Opus 4.8 (1M context).** The built-in `claude`
  agent spec launched with `--model opus` — an **alias** the `claude` CLI resolves to
  "the latest opus", which on some installs/accounts landed on an older Opus (e.g.
  4.6). The spec now pins the exact model ID **`claude-opus-4-8[1m]`** (verified the
  CLI accepts it) so `parley run` always launches Opus 4.8 with the 1M-context window,
  not whatever the alias happens to resolve to. (Tradeoff: a future Opus bump must be
  re-pinned; the alias would auto-track but mis-resolved here.) The local
  `headless-agents.local.json` roster was pinned to match.

## v1.30.0 - 2026-06-19

Pre-idea readiness check (idea `meta-protocol-change-preflight-readiness`; 4-agent
deliberation, signoffs claude-1/codex-1/hermes-1, agy waived; Phase-6 review caught a
CRITICAL §1-bypass + 5 MAJORs, all fixed in fix-up cycle 1).

- **Protocol §9.0 "Pre-idea readiness check"** (both COOPERATION.md copies, drift-guard
  lockstep): at idea start the facilitator (a) checks protocol freshness —
  `source`=advisory/no-write, `consumer` additive bump=auto-sync (zone-preserving),
  breaking/unknown-role=confirm — and (b) hosted-PONG-pings the roster, gating per-idea
  exclude / re-include behind explicit user confirmation. Plus a §5 quorum-locks-at-
  Phase-0 sentence and a §7 carve-out (an upstream version sync is not a protocol change).
- **New `parley preflight` command** `[--dir][--json][--yes][--ping-timeout][--no-ping]`:
  freshness classifier + zone-preserving merge + bounded concurrent hosted-PONG probe
  (process-group-killed on timeout) + report/JSON + exit codes 0/1/2/3. Shared with the
  `parley run` pre-check, which runs **before idea creation**, defaults to hosted PONG
  (`--no-ping`/`--no-preflight` opt out), never auto-answers the new gates, and
  hard-stops unattended without reading stdin. The §1 non-solo hard-stop is evaluated on
  the exact `--participants` set; confirmed exclusions are recorded in `00-prompt.md`.
- `meta/version.json` gains `protocolRole` (`source`/`consumer`, fail-closed); `parley
  init` now writes `protocolRole: consumer`.
- Also bundles a 4-participant roster update (§2 tables → `claude-1`/`codex-1`/`hermes-1`/
  `antigravity-1`; backend map).
- Known follow-ups (deferred): roster-ID↔runtime-ID `-1` reconciliation in reports;
  preflight freshness-probe perf for source/`--no-ping`.

## v1.29.0 - 2026-06-19

Protocol: Fusion + ExecPlans inspiration (idea
`meta-protocol-change-fusion-execplans`; 4-agent deliberation, signoffs
claude/codex/hermes, agy waived on a tooling hang). Additive, **conditional-rigor**
guidance applied byte-identically to both `COOPERATION.md` copies (drift-guard
lockstep); no Go logic changed; embedded `parley init` default stays genericized.

- **`FINAL.md` gains static, self-contained design-time sections** (Phase 4): Purpose
  / user-visible outcome, Context & orientation, **Observable acceptance criteria**,
  **Idempotence & recovery**, Known risks / de-risking. `FINAL.md` stays immutable.
- **`IMPLEMENTATION.md` becomes a living execution doc** (Phase 5): Progress
  (timestamped), Decision Log, Surprises & Discoveries, Validation evidence, Outcomes
  & Retrospective — so a fresh headless agent or the auto-drive driver can resume
  **from the artifact alone**, and §13 `parley retro` gets richer evidence.
- **Advisory "Comparison & blind spots" lens** in `consensus.md` and
  `review/consensus.md` (Phase 3/7) — surfaces what *no* participant addressed. Not a
  gate; append-only signoffs remain the only gate.
- **Phase 6** reviewers may check observable acceptance criteria; severities
  (CRITICAL/MAJOR/MINOR/NIT) unchanged.
- **§13** gains a **confident-error** retro evidence signal (diagnostic only — never a
  new severity, blame label, or merge gate).
- Full living/static sections are required only for complex / `auto_implement` /
  driver-managed / pipeline ideas; trivial or design-only ideas may use `N/A`.
- Explicitly **rejected** (inspiration we did *not* adopt): confidence-by-breadth
  gates, a single-model judge with authority, hiding raw rounds behind a summary, the
  Fusion panel/recursion/cost/web-search machinery, collapsing the deck into one file,
  proceed-without-prompting autonomy across gates, and the anti-list prose maximalism.

## v1.28.1 - 2026-06-16

- **`parley retro` precision fix.** The deterministic scanner matched signal
  patterns in free text, so it false-positived on prose that merely *discussed*
  them — e.g. it flagged `rho-retro-tooling` as "blocked-or-abandoned" because its
  own review consensus quoted ``Verdict: BLOCK``. Blocker detection is now
  anchored to structure: a real `Status: ❌` signoff line, or a `## Verdict`
  heading whose leading token is `BLOCK`/`BLOCKER` (or contains ❌) — not the
  substring "block" in prose, and not a `REQUEST-CHANGES`/"no blocking issues"
  explanation. NOT-FIXED is counted only in review round files, dismissed-findings
  only in consensus files. Regression test included
  (`TestBlockerDetectionIgnoresProse`). Surfaced by dogfooding `parley retro` on
  this repo right after the 1.28.0 ship.

## v1.28.0 - 2026-06-16

Retrospective optimization (RHO adoption — two reviewed ideas,
`meta-protocol-change-rho-retrospective-optimization` + `rho-retro-tooling`):

- **Protocol §13 "Retrospective optimization"** added to both COOPERATION.md
  copies (drift-guard lockstep). A retrospective pass mines the deck's own history
  to *propose* improvements but **applies nothing**: proposals enter as a normal
  idea (protocol-text changes via a meta-protocol-change idea + human approval),
  acceptance is the normal multi-agent gate (consensus + all-participant signoff +
  no-regression), and RHO-style self-preference is a diagnostic note only. Defines
  the layered harness (protocol / runtime "Repository Instruction Files" / local
  "Agent Local Memory" / evidence corpus) and the guardrails (audit,
  adversarial-trajectory hygiene, reversibility, multi-agent diagnosis).
- **New `parley retro` command** — read-only mining of the deck's structured
  artifacts: `scan` (failure-density signals per idea), `select` (type-diverse
  "hard cases" coreset), `diagnose` (grouped report), and `propose --slug` (which
  scaffolds **only** a single new `ideas/<slug>/00-prompt.md`, fail-closed). No
  raw session transcripts and no DPP/embeddings/re-rollout in v1; deterministic.
  Inspired by RHO (arXiv:2606.05922) but replacing its single-model self-preference
  with the deck's multi-agent quorum.

## v1.27.0 - 2026-06-15

- **Auto-drive now works on every transport.** The driver's auto-advance was
  hard-gated to `local-dir`, so a `github-pr` / `gitlab-mr` run stalled at
  round-01 even with auto-drive on. The gate is now transport-independent: the
  canonical artifacts (rounds, consensus, FINAL, …) are the source of truth under
  every transport, so auto-drive advances them everywhere. Only `--auto` /
  `--no-auto` gates it now. The driver still does NOT create PR/MR branches — that
  mirroring stays a manual, ergonomic step.

## v1.26.0 - 2026-06-13

- **New TUI `/run` command.** Advance the protocol on demand from inside the live
  TUI — it kicks the auto-driver (cross-review → consensus → finalize → opted-in
  implementation) for the current run. Most useful with `--no-auto` runs; under
  the default auto-drive it is a no-op once driving has started (idempotent). The
  command appears in `/help` and slash autocomplete.

## v1.25.0 - 2026-06-13

Auto-drive is now the default.

- **`parley run` auto-drives by default.** After round-01 the protocol now
  advances automatically — cross-review rounds, consensus draft, signoff
  requests, and finalize — without you running the next step. Pass **`--no-auto`**
  to opt out (stop after round-01 and advance manually). The flipped flag also
  governs the launch prompt: a default run launches and drives unattended, while
  `--no-auto` (without `--yes`) restores the pre-launch confirmation.
- **Auto-drive now runs inside the TUI.** Previously the driver only ran on the
  `--no-tui` path, so a TUI run stalled at round-01. The driver now runs in the
  background while the live TUI shows it advancing (its output is discarded so it
  never corrupts the render; quitting the TUI stops it).
- **Code-mutation stays gated.** The implementation/fix-up phases (Phase 5–8) are
  still only auto-driven when the idea opts in via `auto_implement`; flipping the
  auto default does not auto-write code. `--no-implement` still stops the driver
  at `FINAL.md`.
- `parley continue` is unchanged: it still prints the next action by default and
  executes it only with `--auto`.

## v1.24.1 - 2026-06-13

Maintenance (idea `embedded-default-protocol-resync`, PR #47):

- **Embedded default protocol resynced** with the live deck. The `parley init`
  bootstrap template (`internal/protocol/defaults/COOPERATION.md`) gained the
  missing `## 12. Pipeline blocks & action stages` section (byte-identical to the
  live deck) and was **genericized**: header `Workspace`/`Created` are now
  placeholders and both §2 tables ship empty bodies, so a freshly `parley init`-ed
  project no longer inherits this repo's roster/workspace.
- **Anti-drift guard**: a fail-closed Go test (`TestEmbeddedDefaultMatchesLiveDeck`)
  asserts the embedded default stays in sync with `parley-deck/COOPERATION.md`
  (modulo five documented, anchored project-specific zones) and that the embedded
  bootstrap shape holds — so a protocol edit landing in only one copy now breaks
  the build. Plus `TestDefaultCooperationForInit` for the init output.
- Synced the project deck to `parley-deck-skill` 1.3.1 (§12 was already present).

## v1.24.0 - 2026-06-12

Adopted from the MIT-licensed "kindly" skill (ideas `runner-hardening-kindly` +
`meta-protocol-change-review-gate-honesty`):

- **Agent supervision**: first-output watchdog (120s, one retry), stall guard
  (30m, output-growth based), persisted `agent.heartbeat` events (60s; excluded
  from transcripts/triggers); counting writers — zero healthy-path I/O; typed
  `agent.no_first_output`/`agent.stalled` events appended BEFORE the kill.
  Config: `first_event_timeout_ms`, `stall_timeout_ms`, `heartbeat_ms`.
- **Failure classification**: `agent.failed` now carries `failure_class` +
  `recovery_hint` (rate-limit/auth/billing/overloaded/…); surfaced in the TUI
  narrator and agent headers.
- **Artifact beats exit code**: a validated artifact with an ordinary nonzero
  exit finishes with `agent_exit` instead of failing (removes the agy
  wrote-then-exit-1 flake); ACP validation now respects the run phase; fix-ups
  validate IMPLEMENTATION.md instead of trusting exit 0; `Result.Success()`.
- **Review snapshots**: Phase 6 reviewers read a disposable shared-clone
  checkout on local tmp (dirty trees become temp-index snapshot commits);
  artifacts move back via copy+fsync+rename; loud fallback events.
- **parley consult** + `parley consults list`: advisory cross-agent questions
  with durable artifacts under parley-deck/consults/ (never quorum evidence).
- Hardening: claude participants shed nested host markers; read-only git probes
  set GIT_OPTIONAL_LOCKS=0; `fsutil.AppendLine`; docs/agent-cli-mechanics.md.
- **Protocol**: Phase 6 "Review briefs and dispositions" (no-suppression),
  Phase 8 opt-in `strict_gate` + "Stopping judgment", §8 "Consults" standing;
  mirrored to the embedded default protocol.

## v1.23.0 - 2026-06-12

- Protocol visibility in the live TUI (idea `tui-protocol-visibility`):
  collapsible protocol ribbon on every tab (Ctrl+P), tab activity glyphs
  (spinner/silent/delivered/failed/STALE), woven narrator lines, a Protocol
  tab with pipeline/delivery/signoff/next panes, and a Home phase column.
- New `run.phase` event emitted by the driver after every phase-changing
  cursor commit; cursor save errors are no longer discarded.
- `driver.RebuildDetail` exports the phase evidence (review round, review
  consensus, implementation status) in one disk pass.
- Declared `buffers_stdout` agent flag (TOML + run.created runtime payload);
  silent buffered agents get a structured placeholder instead of a blank tab.
- Status line shows `ph=N:<phase> wait=<agents>` instead of `round=<status>`.

## v1.5.4 - 2026-05-27

- Treat ACP as a selectable launch mode on an existing agent instead of
  exposing duplicate `*-acp` agent IDs for Codex, Claude, and Hermes.
- Add the TUI `a` key for session-only ACP launch overrides and show ACP
  command details in the selected-agent panel.
- Add `acp_args` runtime configuration so local installs can enable ACP for
  CLIs when their concrete ACP launch args are known.
- Apply TUI launch-mode overrides to newly started runs and record effective
  launch metadata in run runtime events.

## v1.5.3 - 2026-05-27

- Add Antigravity CLI `agy` as a first-class headless agent and default
  replacement for Gemini.
- Mark Gemini as legacy compatibility while keeping existing overrides working.
- Prefer verified stronger defaults for Claude (`opus`/`max`) and Hermes
  (`xai/grok-4.3`) while keeping Antigravity model/thinking fields at
  `cli-default` until the CLI exposes flags.
- Update project and embedded protocol rosters, docs, examples, and runtime
  configuration defaults for the Antigravity migration.

## v1.5.2 - 2026-05-26

- Add TUI planner action execution and focus-aware action controls.
- Refresh dashboard and live TUI layouts with height-aware compact modes,
  two-column normal views, semantic badges, and short-terminal tests.
- Embed the full default `COOPERATION.md` protocol for workspace initialization
  while preserving local-dir bootstrap transport.
- Record concrete project roster metadata and ignore machine-local headless
  agent launch config.
- Improve `parley version --all` project status probing timeouts and fallback
  behavior.

## v1.5.1 - 2026-05-25

- Complete Parley review cycles for the continuous-run TUI planner slice and
  version-awareness project sync.
- Unify planner and manifest next-action serialization through a shared
  `internal/runaction` type.
- Make continuation planning round-aware and remove hardcoded `codex` ownership
  from generated continuation commands.
- Add `parley version --dir DIR --all` project targeting, indented JSON output,
  and cleaner missing-installer fallback errors.
