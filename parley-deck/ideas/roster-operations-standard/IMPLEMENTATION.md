

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

**Tests.** Nine new tests covering A1 (membership, inherited marking, render refusal), A2, A3,
A4/A5, A7, A8 and A15. Full suite green; `go vet` clean.
