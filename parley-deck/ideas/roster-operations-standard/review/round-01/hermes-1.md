---
agent: hermes-1
idea: roster-operations-standard
review-round: 1
date: 2026-08-06
reviewed-commit: 203f73b (v1.40.0)
track: deliberation
---

## Scope

I am hermes-1, a participant in this idea (rounds 1-2, signoff revision 3) and now a
reviewer. I was NOT the implementer (claude-1 implemented). This is a deliberation-track
review; all non-implementers review.

Evidence provenance tags per §15: PRIMARY = I read the source or ran the tool this
session; SECONDARY = I rely on another participant's measurement; RECALL = from my own
prior round analysis, not re-verified this session.

Fresh checks this session (all PRIMARY unless noted):
- `go build ./...` (exit 0), `go test ./internal/app/ ./internal/agents/ ./internal/config/ ./internal/protocol/ ./internal/runmanifest/` (all ok).
- Read in full: internal/agents/launchargs.go, internal/agents/modelmeta.go, internal/app/roster.go (621 lines), internal/app/roster_set.go, internal/app/roster_sync.go, internal/runmanifest/manifest.go, internal/config/runtime.go (819 lines), internal/protocol/drift_test.go.
- Read: internal/app/app.go:111-145 (--help), :826-1001 (sessions inspect), :1131-1193 (continueAuto), :1840-1899 (run creation snapshot).
- `/usr/bin/grep` and `rg -uuu` searches for: §2 generator functions, legacy normalizer, stale-snapshot, --explain, --all, membership confirmation.
- `git log --oneline v1.39.0..v1.40.0` and `git show --stat` on commits 3b94f85, 203f73b.
- Built the binary (`go build -o /tmp/parley-review ./cmd/parley/`) and ran `roster show` on an empty deck and an isolated deck (PARLEY_HOME set to a nonexistent dir).
- Diffed the live §2 table (parley-deck/COOPERATION.md) against the embedded copy (internal/protocol/defaults/COOPERATION.md).
- Read the skill's COOPERATION.md (../parley-deck-skill/skills/parley-deck/references/COOPERATION.md) §2 region.
- Read parley-deck/meta/protocol-changelog.md in full.
- Checked ~/.parley/agents.toml for the D7 interim workaround override.

I did NOT: enumerate the 40-deck fleet (SECONDARY: claude-1's measurement), run a live
agent launch, run `roster sync` against a real deck, or inspect foreign deck repositories.

## Summary

The implementation lands the core of D1-D4, D7, D9 and the protocol text changes well.
The frozen 11-column contract, modelmeta, the {model}/{effort} resolver, STATE wiring,
roster set/sync with rebase semantics, and the protocol-changelog entry are all present
and tested. The build is green and the test suite passes.

However, there are two CRITICAL findings, both of which strike at binding release gates:

1. **G1 is half-satisfied.** The snapshot is written at run creation but `continueAuto`
   does NOT consume it — it still calls `discoverConfigured` and passes the fresh result
   as `Agents: discovered`. The acceptance test FINAL requires does not exist. This is
   the exact hazard G1 was written to close.

2. **G4 is not satisfied.** There is no §2 generator. The protocol text says §2 is "a
   generated, non-authoritative view" and "will be overwritten on the next render," but
   no code renders it. The live §2 table still carries hand-written rows while the
   embedded copy has an empty table — they diverge, and the drift guard masks it by
   emptying both bodies during comparison.

There are also two MAJOR findings: the D5 membership-change second confirmation is not
implemented, and the D7 legacy normalizer is absent. Details below.

## Findings

### [CRITICAL] G1 — snapshot is written but NOT consumed; the continuation acceptance test is absent

FINAL.md G1 (line 88-92): "The change that exposes rebase MUST also persist **and
consume** the immutable effective row, with an acceptance test that creates a run,
changes machine/deck config, continues the run, and proves adapter/model/effort/auto-args
unchanged."

The persist half is done. `RosterSnapshot` (roster.go:600) builds the snapshot at run
creation; app.go:1860-1874 stores `RosterSnapshot` and `RosterRevision` into the manifest
via `runcontrol.Create`. The manifest struct carries both fields (manifest.go:54-57).

The consume half is NOT done. `continueAuto` (app.go:1147-1193) is the continuation path
FINAL explicitly identifies as the hazard ("Today `continueAuto` re-discovers config
(internal/app/app.go:1148-1160)"). It still does exactly that:

```go
discovered, err := discoverConfigured(ctx, root)   // app.go:1152
...
runOpts := runner.Options{
    ...
    Agents: discovered,   // app.go:1163 — fresh discovery, NOT the snapshot
}
```

The manifest's `RosterSnapshot` is loaded nowhere in this function. The snapshot is
written and never read by the continuation path. Changing `~/.parley/agents.toml` or
`parley-deck/agents.toml` mid-run and then running `parley continue --auto` still
silently continues on whatever config is current — the precise failure mode G1 was
created to prevent.

The acceptance test is also absent. I searched every test file for `RosterSnapshot`,
`RosterRevision`, `roster_snapshot`, `roster_revision`, and `snapshot` in
roster_test.go, roster_sync_test.go, and config/roster_test.go (PRIMARY, `/usr/bin/grep`
and `rg -uuu`): zero hits. No test creates a run, changes config, continues, and asserts
invariance. FINAL requires this test as a gate, not a suggestion.

Why this matters: the user chose rebase (VC-2) on the condition that the snapshot
guarantees reproducibility. My own signoff R1 (consensus.md:598-607) states "Rebase and
the immutable run snapshot must ship as one atomic delivery unit" and "Reproducibility
must not depend on an unshipped feature shipping later; it must ship with rebase or
rebase waits." The snapshot persistence without consumption is exactly the half-shipped
state R1 forbids. The rebase semantics in `roster sync` are live, but the safety
precondition that makes rebase safe for resumable runs is not.

Suggested fix: `continueAuto` must load the manifest's `RosterSnapshot` and, when
present, build `runner.Options.Agents` from the frozen entries rather than from
`discoverConfigured`. When the snapshot is absent (legacy pre-snapshot run), fall back to
the current discovery path and emit the "unsafe for pre-snapshot resumable runs" warning
FINAL requires. Then add the acceptance test: create a run, mutate config, continue,
assert the agents slice matches the snapshot.

### [CRITICAL] G4 — there is no §2 generator; §2 is not a generated view

FINAL.md G4 (line 104): "Generated §2 is idempotent. Two runs of the generator produce
byte-identical output."

consensus.md (line 376-382): "Runtime code MUST NOT parse the generated §2 as roster
authority. Today `resolveRoster` does (internal/app/roster.go:110); that call site is
removed in the same change."

There is no §2 generator. I searched exhaustively for any function that generates,
renders, or writes the §2 roster table (PRIMARY, `/usr/bin/grep -rn` for
`func.*[Gg]enerat.*[Ss]ection`, `func.*[Rr]ender.*[Ss]ection`, `GenerateSection2`,
`WriteSection2`, `UpdateSection2`, `SyncSection2`, `regen`, `rerender`,
`render.*cooperat`, `generate.*cooperat`): zero matches (exit 1).

The protocol text in all three COOPERATION.md copies now says: "the table below is a
generated, human-readable view and is NOT authoritative. Do not hand-edit it ... the edit
will not take effect and will be overwritten on the next render." But there is no "next
render" — no command or code path regenerates the table from `agents.toml`.

Consequences:

1. The live §2 table (parley-deck/COOPERATION.md:129-134) still carries four hand-written
   rows (claude-1, codex-1, hermes-1, kimi-1). The embedded copy
   (internal/protocol/defaults/COOPERATION.md) has an EMPTY table body (zero rows). They
   diverge. I confirmed this with a direct `diff` of the §2 regions (PRIMARY).

2. The drift guard (drift_test.go) masks this divergence: `normalizeProtocol` (line
   122-124) empties BOTH table bodies before comparison, and `assertEmptyTableBody` (line
   163-178) only asserts the EMBEDDED copy is empty — it does NOT check the live deck.
   So a live deck with stale hand-written rows passes the drift test silently. This is
   the exact "generated table that code still reads is a second stale view" failure
   kimi-1 warned about (consensus.md:807).

3. `resolveRoster` still parses the legacy §2 table via `protocol.ReadRosterIDs` when no
   `[roster.*]` entries exist (roster.go:243). The consensus requirement "that call site
   is removed" is not met — it is the fallback path. This is acceptable for legacy decks
   (the `legacy-roster` STATUS), but the protocol text says §2 "is NOT authoritative" and
   "will be overwritten," which is false when no generator exists.

Why this matters: the entire D9 decision is that §2 "becomes a generated, non-
authoritative view." Without a generator, §2 is still the hand-edited store it always
was — the protocol text now lies about its own state. A user reading the new §2 header
will believe editing it is futile, but for a legacy deck (17 of 40 have no `[roster.*]`
at all), the §2 table is STILL the only membership source the CLI reads.

Suggested fix: ship a `parley roster render` (or fold it into `roster sync`/`roster set`)
that reads `agents.toml` and writes the §2 table deterministically (active before
inactive, agent ID byte-ascending, per the consensus ordering rule). Pin it with an
idempotency test (two runs → byte-identical). Then the protocol text's "overwritten on
the next render" becomes true.

### [MAJOR] D5 — `--yes` alone is not refused for membership changes

FINAL.md D5 (line 64-65): "Preview by default; `--yes` applies; `--yes` alone is
**refused** for membership changes."
consensus.md (line 107-108): "`--yes` alone is refused when the change alters membership
— a second explicit confirmation is required."

`roster set` with `--yes` adds a new member with no second confirmation. I tested this
directly (PRIMARY): `parley roster set newagent-1 --scope deck --adapter claude --model
foo --yes` wrote a new `[roster.newagent-1]` block and exited 0. The code at
roster_set.go:126-139 appends a new block on `--yes` with no membership-specific gate.

The comment at roster_set.go:127-128 says "the caller gates it behind the breaking-change
confirmation," but the caller (`runRoster`, roster.go:112-135) does no such gating — it
passes `*yes` straight through to `rosterSet`.

This is a safety property the consensus made unanimous and explicit. A single `--yes`
adding a new agent to the roster is exactly the "mass mutation where one bad deck gets
swept through" risk hermes-1's R3.2 raised for the fleet migration, but at the
single-deck level.

Suggested fix: detect whether `applyRosterBlock` will create a new `[roster.<id>]` header
(start < 0 in the current logic). If so, require either a second flag (e.g. `--confirm-
membership`) or an interactive prompt; refuse on `--yes` alone with a message naming the
new member.

### [MAJOR] D7 — the legacy normalizer is absent

FINAL.md D7 (line 71-73): "built-in `HeadlessArgs` carry `{model}` and `{effort}`
placeholders ... plus a **legacy normalizer** for configs that hardcode a model literal in
`headless_args`."
FINAL.md (line 152-153): "This override MUST be removed when D7 lands — a wholesale
`headless_args` override is exactly how `hermes` silently lost `--yolo`. The removal is
part of Stage 1's legacy normalizer work."

The {model}/{effort} placeholders are in the built-in specs (discover.go:202, 225, 250,
298, 329, 361) and `ResolveLaunchArgs` (launchargs.go:57) substitutes them correctly. This
is the first half of D7 and it works.

The legacy normalizer is absent. I searched for `normaliz.*model`, `legacyModel`,
`legacyNorm`, `normalizeHeadlessArgs`, `normalizeArgs`, `legacyHeadless`, `baked.*model`,
`hardcod.*model.*headless`, `model.*literal.*headless` (PRIMARY, `/usr/bin/grep -rn` over
internal/): zero matches (exit 1).

`applyOverride` (runtime.go:625-627) still does `spec.HeadlessArgs =
expandSlice(override.HeadlessArgs, root, tempdir)` — a wholesale replacement. A config
layer that overrides `headless_args` with a literal model (e.g. `["exec", "--model",
"claude-opus-4-8[1m]"]`) still replaces the built-in args that now carry `{model}`, and
there is no normalizer to detect this and inject the placeholder. The model would be
correct but `{model}` would be gone — and if the user later changes the config `model`
field, the argv would NOT update, which is the original defect D7 exists to fix.

The interim workaround is still in `~/.parley/agents.toml` (PRIMARY, `/usr/bin/grep`):
line 38-40 carry a `headless_args` override pinning `claude-opus-5[1m]`, with a comment
saying "headless_args override is exactly how hermes silently lost --yolo." FINAL says
this "MUST be removed when D7 lands." The removal is a user-environment action, not a
code change, so I cannot verify it was done — but the normalizer that would make the
removal safe is not in the code.

Why this matters: without the normalizer, a deck that still carries a `headless_args`
override (like the current machine config) does not benefit from the {model} fix. The
override replaces the placeholder, so `EffectiveModel` may still report the baked-in
literal rather than the config `model` field. The root cause D7 was designed to close
remains open for any such deck.

Suggested fix: in `applyOverride`, when `override.HeadlessArgs` is set, detect whether it
contains a model flag (`--model`/`-m`) with a literal value and either (a) replace the
literal with `{model}` so the config field drives it, or (b) warn that a hardcoded
headless_args bypasses the model field. At minimum, document that the override must be
removed manually and that the normalizer is deferred.

### [MINOR] D5 — `--explain AGENT` and `--all` flags are not implemented

FINAL.md D5 (line 59): `parley roster show [--scope deck|machine] [--all] [--json]
[--explain AGENT]`
consensus.md (line 100): same.
consensus.md (line 113): "`roster show --all` reveals IDs mapped in config but absent
from the deck roster, clearly marked — this is what would have made the `opencode`
situation visible on day one (kimi-1)."

Neither `--explain` nor `--all` is implemented. `rosterUsage` (roster.go:57) shows only
`[--scope deck|machine] [--dir DIR] [--json]` for `show`. The flag parser (roster.go:87-
100) defines no `explain` or `all` flag. I confirmed with `/usr/bin/grep -n "explain\|
--all" internal/app/roster.go` (exit 1, no matches).

`--all` was kimi-1's specific ask for making the opencode situation visible. Without it,
a deck whose machine config maps an agent absent from the deck roster has no single
command that surfaces it.

This is MINOR rather than MAJOR because the core `roster show` works and the 11-column
contract is intact; these are documented flags that did not ship.

### [MINOR] D6 — `sessions inspect` does not report `stale-snapshot`

FINAL.md D6 (line 68-69): "`sessions inspect` reports `stale-snapshot`."
consensus.md (line 119): "`sessions inspect` reports `stale-snapshot` when the deck
roster has moved since."

The manifest carries `RosterRevision` (manifest.go:57), and `RosterSnapshot` is stored at
run creation. But `inspectSession` (app.go:938-957) loads the manifest and never compares
its `RosterRevision` against the deck's current roster. `printSessionDetail`
(app.go:959-1001) prints manifest status/mode/transport/created but no `stale-snapshot`
line. I searched for `stale.snapshot` and `StaleSnapshot` in internal/app/ (PRIMARY):
only the manifest comment at manifest.go:56 mentions it.

The `stale-snapshot` STATUS code is in the D3 vocabulary (FINAL.md line 51) and the skill
(hermes SKILL.md roster section), but it is never produced by any code path.

Suggested fix: in `inspectSession`, after loading the manifest, call `RosterSnapshot` on
`session.WorkspaceRoot`, compute the current revision, compare to
`manifest.RosterRevision`, and set a `StaleSnapshot bool` on the payload. Print it in
`printSessionDetail`.

### [MINOR] The three COOPERATION.md copies are not identical in §2

G3 (FINAL.md line 99-102) requires the authority cutover to be atomic across "Three
COOPERATION.md copies, per the standing drift guard."

The live deck copy (parley-deck/COOPERATION.md) §2 table has 4 rows (claude-1, codex-1,
hermes-1, kimi-1) plus the host-handle table with 4 rows. The embedded copy
(internal/protocol/defaults/COOPERATION.md) §2 table is EMPTY (0 rows) and its host-handle
table is EMPTY. The skill copy (../parley-deck-skill/.../COOPERATION.md) §2 table is also
EMPTY (0 rows).

The drift test passes because `normalizeProtocol` empties both bodies before comparing
and `assertEmptyTableBody` only checks the embedded copy. This is not a true divergence
in protocol RULE text — the §2 table is a project-specific zone the guard intentionally
normalizes. But the live deck carrying rows the embedded/skill copies do not is a
consistency gap: a `parley init` from the embedded copy produces a deck with an empty §2,
while the live deck has a populated one. If a generator existed (G4), this would be
self-healing. Without one, the live deck's rows are hand-maintained prose that the
protocol now says is non-authoritative.

This is MINOR because the drift guard is working as designed (the §2 table is an
allowlisted zone), but it underscores the G4 gap.

### [NIT] D1 — `agents list` is not relabelled in help or output

FINAL.md D1 (line 34): "`parley agents list` = adapter/runtime inventory (relabelled)."
consensus.md (line 35): 'relabelled "adapter/runtime inventory — not the roster".'

`runAgentsList` (app.go:422) still uses `flag.NewFlagSet("agents list", ...)` and
`agents.PrintRuntimeMatrix` with no relabelling. The `--help` text (app.go:121) still
says `parley agents list [--dir DIR]` with no "adapter/runtime inventory" framing. I did
not find any string "adapter/runtime inventory" or "not the roster" in the agents list
output path (PRIMARY, `/usr/bin/grep`).

The relabelling may be in `PrintRuntimeMatrix` output text; I did not read that function
this session. If the output header was changed, this NIT is moot. But the help text and
flag set name are unchanged, so the command surface does not reflect the relabelling.

### [NIT] Docs do not mention `roster`

FINAL.md D1 (line 36): "`roster` must appear in `parley --help` and the docs."
consensus.md (line 42-44): "absent from docs/cli-reference.md and
docs/agent-runtime-configuration.md — kimi-1, PRIMARY."

`roster` now appears in `parley --help` (app.go:123-124, confirmed PRIMARY). But
`/usr/bin/grep -n "roster" docs/cli-reference.md docs/agent-runtime-configuration.md`
returned zero matches (exit 0, no output). The docs gap kimi-1 identified is not closed.

## Decision-by-decision landing assessment

- **D1 (three concepts, one answer):** PARTIALLY lands. `roster show` works and is in
  `--help`. `agents list` is not relabelled (NIT). Docs do not mention `roster` (NIT).
  The run snapshot exists but is not consumed (CRITICAL G1).
- **D2 (MODEL/EFFORT effective-or-unknown):** LANDS. `EffectiveModel`/`EffectiveEffort`
  (launchargs.go:97-122) read the resolved argv, not the declaration. `model-drift`,
  `model-unbound`, `effort-unknown` STATUS codes are produced (roster.go:321-335).
  PRIMARY.
- **D3 (frozen 11-column contract):** LANDS. `RosterColumns` (roster.go:153-156) matches
  the contract exactly. `RosterSchemaVersion = 1`. JSON carries `schema_version` and
  `columns` (roster.go:385-389). Golden test in roster_test.go pins the header. PRIMARY.
- **D4 (modelmeta CLI-owned):** LANDS. `DeriveModelMeta` (modelmeta.go:70) peels
  gateways, never infers company from adapter, returns `unknown`/`metadata-unknown` on no
  match. Golden test covers plain, qualified, gateway-routed, and unresolvable cases.
  PRIMARY.
- **D5 (three verbs):** PARTIALLY lands. `show`/`set`/`sync` work. `--keep` works (tested
  in roster_sync_test.go). Preview is default. `--scope deck` writes committed
  `agents.toml` (tested). BUT `--yes` alone is not refused for membership (MAJOR),
  `--explain` and `--all` are absent (MINOR).
- **D6 (session = immutable run snapshot):** PARTIALLY lands. Snapshot is written at run
  creation. `stale-snapshot` is NOT reported by `sessions inspect` (MINOR). Snapshot is
  NOT consumed by `continueAuto` (CRITICAL G1).
- **D7 (model-argv fix + legacy normalizer):** PARTIALLY lands. `{model}`/`{effort}`
  placeholders in built-in specs + `ResolveLaunchArgs` works. Legacy normalizer is absent
  (MAJOR). The interim workaround override is still in `~/.parley/agents.toml`.
- **D8 (skill/CLI boundary):** LANDS. The skill's SKILL.md roster section
  (../parley-deck-skill/.../SKILL.md:254-295) invokes `parley roster show` and reproduces
  its output; it does not parse §2 or TOML. PRIMARY.
- **D9 (§2 generated view):** DOES NOT LAND. No generator exists (CRITICAL G4). Protocol
  text claims §2 is generated, but it is still hand-maintained. `agents.toml` is the
  authority for decks with `[roster.*]`; legacy decks still fall back to §2 parsing.

## Gate-by-gate assessment

- **G1 (rebase gated on snapshot consumption + test):** NOT SATISFIED. Snapshot is
  written but not consumed; the acceptance test does not exist. CRITICAL.
- **G2 (STATE wiring):** SATISFIED. `resolveRoster` (roster.go:237-269) now reads the
  inactive map from both `LoadRoster` (entries with `Active=false`) and the legacy
  `ReadRosterIDs` inactive map. Inactive rows get `State: "inactive"` and
  `addStatus("inactive")`. The discard-into-`_` bug is fixed. PRIMARY.
- **G3 (atomic authority cutover, 3 copies + skill):** PARTIALLY SATISFIED. All three
  COOPERATION.md copies have the new §2 authority text. The protocol-changelog entry
  exists. The skill SKILL.md is updated. BUT the copies diverge in §2 table content
  (live has rows, embedded/skill are empty), and there is no generator to reconcile them
  (ties to G4). The drift guard masks the divergence. MINOR (standing alone) but it
  enables the G4 failure.
- **G4 (idempotent §2 generator):** NOT SATISFIED. No generator exists. CRITICAL.
- **G5 (protocol-changelog entry):** SATISFIED. parley-deck/meta/protocol-changelog.md
  has a dated entry naming the idea, the change, the venue deviation, and the one-off
  authorization. PRIMARY.

## The legacy fallback

A deck with no `[roster.*]` but a §2 table falls back to `protocol.ReadRosterIDs`
(roster.go:243) and reports `legacy-roster` on every row (roster.go:264). This is correct
and safe — it matches consensus.md:146-147 and FINAL.md:124-125.

A deck with NEITHER (no `[roster.*]` and no §2 table): `LoadRoster` returns the machine-
level entries if `~/.parley/agents.toml` has `[roster.*]` blocks. I confirmed this
(PRIMARY): an empty deck with `PARLEY_HOME` pointing at the real home shows the machine
roster. With `PARLEY_HOME` set to a nonexistent dir, `roster show` errors: "no roster:
declare [roster.<id>] in parley-deck/agents.toml (or keep a legacy §2 table in
COOPERATION.md)" (exit 1). This is safe — the error is clear and actionable. The only
subtlety is that a deck with no deck-level roster silently inherits the machine one,
which may surprise users who expect an empty deck to have no roster. This is not a
finding; it is the designed layering behavior.

## `roster sync` removes deliberate pins by default

The preview/`--keep` protection is well-designed for the single-deck case:
- Preview is default (roster_sync.go:121: `if dryRun || !yes`).
- Deliberate pins are enumerated with the exact `--keep AGENT.FIELD` needed to retain them
  (roster_sync.go:108-112).
- `--keep` exempts a pin (tested in TestRosterSyncKeepExemptsAPin).
- The report states how many redundant overrides and pins were removed (roster_sync.go:
  139).

The concern: `roster sync --yes` removes deliberate pins (not just redundant ones) from a
committed file, and there is no second confirmation for pin removal (unlike the D5
membership case, which at least has the requirement even if unimplemented). A pin is a
deliberate deck divergence; removing it silently changes the deck's behavior. The preview
shows the pins, but `--yes` bypasses the human. For a single attended invocation this is
acceptable. For a fleet operation (Stage 4, not yet shipped), this needs the per-deck
confirmation hermes-1's R3.2 and kimi-1's R3 require. Since Stage 4 is not in this
release, this is not a finding against 1.40.0, but the single-deck path sets the
precedent that `--yes` removes pins without a second gate — worth noting for the fleet
migration.

## Anything shipped that FINAL did not decide, or decided and not shipped

Shipped but not decided by FINAL: nothing I found. Every code change traces to a FINAL
decision or gate.

Decided by FINAL but not shipped:
1. G1 snapshot consumption + acceptance test (CRITICAL).
2. G4 §2 generator (CRITICAL).
3. D5 `--yes`-alone-refused for membership (MAJOR).
4. D7 legacy normalizer (MAJOR).
5. D5 `--explain` and `--all` flags (MINOR).
6. D6 `stale-snapshot` reporting in `sessions inspect` (MINOR).
7. D1 `agents list` relabelling (NIT).
8. D1 docs mention of `roster` (NIT).

## Open questions

1. The implementer's commit message for 3b94f85 says "§2 becomes a generated view
   (protocol, 3 copies)." Was the generator intended to ship in this commit and dropped,
   or was it always deferred? The protocol text claims generation is live. If it was
   deferred, the protocol text should say "will be generated" or the generator should
   ship before the text claims it.

2. Is the `continueAuto` consumption gap a known deferral or an oversight? The snapshot
   struct and write path are clearly built; the absence of the read path in the one
   function FINAL names suggests it was planned but not completed.

3. For the fleet migration (Stage 4, separate attended operation): will the G1
   consumption gap affect migration safety? If migration runs `roster sync` on decks with
   nonterminal pre-snapshot runs, and those runs are later continued, the continuation
   will use fresh config, not the snapshot — because there is no snapshot to consume.
   FINAL's G1 release condition 2 says "Fleet migration MUST skip and report every
   nonterminal legacy run lacking that snapshot." That skip logic is Stage 4 and not yet
   shipped, but it depends on G1 consumption working for new runs, which it does not yet.
