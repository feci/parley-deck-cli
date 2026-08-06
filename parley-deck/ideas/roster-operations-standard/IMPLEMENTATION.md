

## Fix-up cycle 2 — review consensus A1–A16 (2026-08-06)

Review consensus rev 4 accepted by hermes-1, codex-1, kimi-1. All sixteen agreed fixes landed;
DF-1..DF-6 remain deferred as recorded.

**A1 (CRITICAL) — deck membership is the deck file.** `config.LoadRosterScoped` splits the layered
roster into *membership* (deck file only) and *values* (all layers, machine seeding deck). Before
this, `LoadRoster` unioned membership across layers, so a deck declaring two participants resolved
to five and — since fix-up cycle 1 routed participant selection through the same view — would have
**run** five. `RosterMembership`, `resolveRoster` and `renderRosterTable` all consume the scoped
view. A deck declaring no roster of its own still displays the machine roster, with every row
marked `inherited-roster`, and `roster render` refuses to commit it without `--adopt-inherited`.
Legacy fallback per codex-1's condition: "no roster of its own" means neither a deck `[roster.*]`
block nor a valid legacy §2 table.

**A2 — §2-only IDs are reported, never silently erased.** `section2OnlyRows` surfaces them with
`unmapped` + `section2-only`; per the ratified field table they are never auto-added to the
generated §2, and `roster render` now reports every row it removes, in preview and on apply.

**A3 — membership gate keys on block existence.** `membershipChange` takes `existed bool` from
`blockExists`, checked against the same bytes the patcher edits. Any first write to a missing block
is a membership change regardless of which field is written.

**A4/A5 — snapshot freezes per roster ID and pins auto-args.** The frozen map is keyed by `Agent`
(adapter is a fallback for pre-ID runs), and `RosterSnapshotEntry.LaunchArgs` carries the resolved
argv so a continuation cannot change autonomy posture under a running idea.

**A6 — D5 grammar implemented.** `--all`, `--explain AGENT` and a load-bearing `--scope
deck|machine`. `roster init` accepts the current spelling and prints a deprecation notice; the
hidden `session` alias is no longer the visible default in flag help.

**A7 — one JSON/text contract.** An empty status list renders `["ok"]` rather than `null`, matching
the text table; the payload carries `scope`. Golden test asserts text and JSON together.

**A8 — legacy normalizer (D7's second half).** `agents.NormalizeLegacyModelArgs` rewrites a
hardcoded model/effort literal in a config-supplied `headless_args` back to `{model}`/`{effort}`, so
the declared field wins; the rewrite is recorded in `Sources["headless_args_normalized"]`.

**A9 — §2-as-a-store instructions removed from all four surfaces** (three protocol copies +
SKILL.md), with `TestNoSection2AsAStoreInstructions` pinning the exact phrases so the contradiction
cannot silently return.

**A10 — `sessions inspect` reports `stale-snapshot`** via `rosterSnapshotState`, comparing the
frozen revision against the deck's current roster.

**A11 — changelog entry reformatted to the §7 template** (`Idea:` as a path, `Drafted by:`,
`Summary:`), keeping the substantive prose beneath it.

**A12 — `masked-by-env` has an emitter.** `roster set` re-resolves after writing and warns when a
higher layer overrides the value just written, instead of reporting a false success.

**A13 — discoverability.** All five verbs in `parley --help`; `agents list` relabelled
"adapter/runtime inventory — NOT the roster"; both named docs gained full roster sections
(previously zero mentions of `roster`).

**A14 — skill corrected (2.5.1).** `roster sync` no longer described as the legacy migration path,
because it moves nothing across; `migrate`/`render`/`set` are named instead. Membership-is-the-deck
rule, the new status codes and the new flags documented.

**A15 — sync hardening.** Unmatched `--keep` tokens are a hard error, and apply is bound to the
previewed old-values so an edit between preview and apply is refused rather than discarded.
`--drop-pins` deliberately not adopted (recorded in the consensus).

**A16 — assorted.** File-mode regression fixed (`writeRosterFileAtomic` preserves the target's mode
instead of leaving 0600); continuation warns when a manifest is unreadable; stale `roster init`
guidance replaced; `modelmeta`'s broad `"k"` rule moved after `"kimi"`; reactivation no longer
reports itself as a retirement.

**DF-1 interim guard.** `roster migrate --yes` now requires `--confirm-breaking`, and a deck with
uncommitted changes is skipped and reported, so a second unattended fleet run cannot happen while
the full migration contract is outstanding.

**Tests.** Ten new tests covering A1 (membership, inherited marking, render refusal), A2, A3,
A4/A5, A7, A8 and A15. `go vet` clean; `go test ./...` green on this machine. Note for reviewers:
`internal/runner`'s `TestDurableKillEndToEndRealProcess` is environment-dependent (it needs a
readable boot id) and fails under some sandboxes; it predates this cycle and `internal/runner` is
untouched by it.

## Fix-up cycle 3 — re-review round 2 findings (2026-08-06)

Re-review round 2 returned FINDINGS from all three reviewers. Every finding is fixed below.

**[CRITICAL, all three reviewers] A1's legacy-fallback clarification was never implemented.**
The consensus adopted codex-1's condition — "no roster of its own" means neither a deck `[roster.*]`
block **nor** a valid legacy §2 table — and cycle 2 recorded it as landed while
`LoadRosterScoped` knew nothing about §2. A legacy deck declaring four members inherited the
machine's five instead, quorum included. `LoadRosterScoped` now decides authority in an explicit
order: committed deck blocks → valid legacy §2 (`Legacy: true`, rows reported `legacy-roster`) →
machine roster (`Inherited: true`). Same defect class as cycle 1: documented as done, not done.

**Membership narrowed to the committed deck file.** codex-1 noted every non-machine layer was
pooled as deck membership, so `agents.local.toml` or `$PARLEY_HEADLESS_AGENT_CONFIG` could add
members. Only `parley-deck/agents.toml` carries `membership: true`; the rest supply values.

**[MAJOR] A3 — `roster init --yes` bypassed the membership gate.** It writes `[roster.*]` blocks,
so it now requires `--confirm-breaking` too. Two existing tests encoded the old behavior and were
updated.

**[MAJOR] A5 — per-ID pins still collapsed at the real call site.** `applyRosterSnapshot` keys by
roster ID, but `continueAuto` applied it to adapter-keyed discoveries, so every lookup took the
adapter fallback. `applyRosterSnapshotToParticipants` resolves participants first, then freezes.
The cycle-2 test proved the function; the new test proves the boundary.

**[MAJOR] A6 — machine scope was still layered.** Specs, adapter mappings and field provenance now
have `*Scoped` variants; `--scope machine` no longer reports deck-only values.

**[MAJOR] A7 — `display_name`/`note` are no longer serialized.** They are outside the frozen eleven
columns; keeping them made "the same table in text and JSON" false.

**[MAJOR] A9 — the drift guard covered two of four surfaces.** It now also reads the bundled skill
protocol and `SKILL.md`, and its banned-phrase list gained the stale `roster sync moves it
over/across` claim — which immediately failed on all three protocol copies (hermes-1's MAJOR), now
replaced with the real migration path.

**[MAJOR] A10 — the run revision ignored frozen launch args**, so autonomy drift reported
`current`. `RosterRevisionOf` hashes `LaunchArgs`.

**[MAJOR/MINOR] A12 — `masked-by-env` false-positived on every machine-scope write**, comparing a
display label against an absolute path. `config.RosterSourcePath` resolves the label first.

**[MINOR] A16 residuals** — unmapped guidance names `--confirm-breaking`; `modelmeta`'s broad `k`
prefix replaced by `kimiCodename` (k + digit); kimi-1 added to the 1.40.1 CHANGELOG attribution.

**[NIT] Handoff accuracy** — the test-count and "full suite green" claims are corrected, with the
environment-dependent `internal/runner` failure named explicitly.

**Tests.** Six more: legacy-§2-beats-machine, snapshot-survives-participant-resolution,
revision-covers-launch-args, machine-scope-not-masked, init-requires-confirm-breaking, plus the
extended drift guard.

## Fix-up cycle 4 — re-review round 3 findings (2026-08-06)

Round 3: hermes-1 CLEAN; codex-1 and kimi-1 FINDINGS. All fixed.

**[CRITICAL, codex-1] Value-only layers could retire or revive a committed deck member.**
Cycle 3 gated which layers may add membership IDs but left `active` merging from every layer, so
`[roster.claude-1] active = false` in the gitignored `agents.local.toml`, the env config, or the
machine file silently dropped a committed member from the quorum — using a file collaborators
never see. `RosterScope.applyAuthorityState` now forces each member's state to the value declared
by the layer that granted membership (deck file, legacy §2, or machine when inherited).

**[MAJOR, codex-1 R3-M1 / kimi-1 N-1] A regression introduced by cycle 3.** The A5 fix appended
the frozen discovery instead of replacing it. The runner resolves a participant by FIRST exact ID
match, so for a participant spelled exactly like its adapter family (`claude`, not `claude-1`) the
live unfrozen record still won: the freeze existed in memory and never reached the launch.
Composite IDs hid it because they do not collide with the adapter ID they resolve through. The
helper now replaces in place, and the test resolves through `agents.ResolveParticipant` and asserts
the record the runner would actually launch.

**[MAJOR, codex-1 R3-M2] `--scope machine` still leaked.** `--all` used the layered spec/mapping
loaders, and `--explain` reported "built-in default" for values that reach the launch through an
`[agents.<family>]` block. Both are scope-threaded now, and `config.AgentFieldSources` attributes
agent-block fields.

**[MAJOR, codex-1 R3-M3] The four-surface guard overclaimed.** It now states explicitly when the
sibling skill checkout was not readable and that those surfaces are therefore not enforced by that
run, instead of implying four-surface coverage.

**[MINOR, codex-1 R3-m1] The JSON contract had no reversion-sensitive test.** The previous test
only asserted status was non-null, so re-adding `display_name`/`note` would have passed. The new
test asserts the exact eleven-key set in both directions.

**[MINOR, both] Skill 2.5.1** — corrected on all four surfaces in this cycle; the release itself
is performed at Phase 8 close (see cycle 5).

**Tests.** Four more: value-layers-cannot-change-membership-state, freeze-reaches-bare-family-
participants, exact-frozen-JSON-columns, plus the corrected guard note.

## Fix-up cycle 5 — re-review round 4 findings (2026-08-06)

Round 4: all three reviewers verified every round-03 finding fixed at the real call site, and
hermes-1 confirmed each new test FAILS when its fix is reverted (no tautologies). Two findings
remained.

**[MAJOR, codex-1 R4-M1] The diagnostics contradicted the result they described.** Cycle 4 made
`active` follow the membership authority, but `roster show --explain` and the `masked-by-env`
warning still read it from the general layer stack — so `--explain` attributed a member's
effective state to an env file asserting the opposite, and `roster set --state` warned that a
write to the authority was masked when it was not. `config.RosterStateSource` names the deciding
layer, and both surfaces use it. Fixing the behavior without fixing what reports it leaves two
sources of truth, which is the defect this whole idea exists to remove.

**Gate correctness, found while verifying the above.** `roster set claude-1 --state active` on an
already-active member demanded `--confirm-breaking`: `membershipChange` fired on the presence of
`+ active = true` rather than on a state change, and absence of the key already means active. A
gate that fires on no-ops trains operators to pass the confirmation reflexively, which is how a
gate stops being one. It now compares the prior state in the file being edited.

**[MINOR, all three] Skill 2.5.1** is released as part of the Phase 8 close.

**Tests.** Two more: active-provenance-and-masking-follow-the-authority, and
membership-gate-ignores-no-op-state-writes.

## Fix-up cycle 6 — re-review round 5 findings (2026-08-06)

Round 5: hermes-1 CLEAN; codex-1 (MAJOR) and kimi-1 (MINOR) found the same defect — cycle 5's own
regression.

**Cycle 5's state-provenance override ignored scope.** `config.RosterStateSource(root)` always
resolved the DECK authority, so `roster show --scope machine --explain` attributed `active` to
`parley-deck/agents.toml` while its own membership header named the machine file — the same output
contradicting itself. `rosterExplain` now uses the scope it already resolved, and the masking check
uses `RosterStateSourceForTarget`, which recognizes a machine-scope write as governed by the
machine roster.

This is the second cycle in a row where fixing a reporting surface introduced a new inconsistency
in that surface. The test added here asserts the two halves of one output agree with each other,
rather than checking either alone.

**Tests.** One more: active-provenance-is-scope-aware, which fails if provenance and the membership
header of the same command disagree.
