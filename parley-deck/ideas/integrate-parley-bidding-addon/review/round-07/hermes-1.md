---
idea: integrate-parley-bidding-addon
review-round: 7
agent: hermes-1
date: 2026-07-30
reviewed-commit: 49fc3ec
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

### Round-6 codex-1 MAJOR — identity check exempted an absent `skill` field — CLOSED

The guard at lib/installer.js:1355 is now `state.marker.skill !== unit.skill` with no
`!== undefined` exemption. An absent identity gets its own problem message ("carries no
skill identity, so this directory cannot be confirmed as …"), distinct from the
mismatched-identity message. I confirmed by runtime probe (P1 below): full install, delete
only `skill` from both the core and bidding markers, leave every other field intact.
Result: both units `malformed`, `managed: false`, `doctor.ok: false`, and `install --only
parley-bidding` and `uninstall` both return `ok: false`. Health and the two mutation paths
now give one answer. The shared predicate `installerOwnsDestination` (lib/installer.js:1652)
requires `state.marker.skill === skill` with no exemption, so the two consumers are
structurally aligned.

The added regression test ("a marker with no skill identity is malformed for core and
add-on alike") exercises exactly this branch — it deletes `skill`, asserts `malformed` +
`managed: false` + the "no skill identity" problem for both units, and asserts both
mutation commands refuse.

### Round-6 kimi-1 MINOR — `selected` derived from the inspection flag, not the recorded selection — CLOSED

`targetSkillUnits` (lib/installer.js:894-912) now reads the recorded selection from the
core marker for read commands (`doctor`, `status`, `paths`) and derives each add-on's
`selected` from it: `recordedSet.has(name)`. Install and uninstall keep the requested set
as the selection (they are writing it), so `recorded` is `null` for them and `selected`
stays `true`. No core marker of ours means nothing to contradict, so units stay
`selected: true`.

I confirmed by runtime probe (P2 below): full install, then `install --force --only
parley-design` so bidding is a residual on disk. Unflagged `doctor` reports
`parley-bidding: selected: false, valid-unselected, ok: false`. `doctor --only
parley-bidding` on the same tree reports `selected: false, valid-unselected, ok: false` —
the same recorded fact. The narrowing survives: the scoped read inspects only
`[parley-deck, parley-bidding]`, not the other recorded add-ons. Two reads, one answer.

The added regression test ("a filtered read reports the same recorded-selection fact as an
unfiltered one") asserts both the unflagged and scoped `selected: false` / `valid-unselected`
and the scoped unit-list narrowing.

### Round-6 hermes-1 NIT — `valid-unselected` masks `valid-unmanaged` for a foreign copy outside the selection — NOT CLOSED, not blocking

This was a NIT in round 6 and I explicitly marked it not blocking. It is unchanged at
49fc3ec: the status precedence remains `malformed > valid-unselected > valid-unmanaged >
valid`. A foreign unmanaged copy outside the recorded selection still reports
`valid-unselected` rather than `valid-unmanaged`. The `managed: false` boolean
disambiguates provenance in both cases, the code comments direct automation to use
`managed` not the status string, and the current precedence is the stricter health
outcome. The new round-7 test ("a tree that is both unselected and unmanaged reports the
selection fact first") explicitly asserts this precedence is intentional —
`valid-unselected` wins when both facts are true, with `managed: false` still carried.
This is a documented design choice, not a defect. Follow-up, not a blocker.

### Round-6 hermes-1 open question — status-string precision vs `managed` authority — answered by the implementer

The cycle-7 commit message and the new both-true test make the design intent explicit: the
selection fact is actionable and wins the status; `managed` carries the provenance fact.
This is a coherent answer to my open question. No action needed.

## New findings

None.

## Release judgement

Releasable as 2.1.0. The two round-6 blockers (codex-1 MAJOR, kimi-1 MINOR) are both
fixed, regression-tested, and verified by runtime probe at 49fc3ec. The shipped payload
`skills/parley-bidding/` has not changed since 714712f — this entire cycle was about the
installer's health and ownership logic, not the payload. The round-6 NIT is a documented
design choice, not a release blocker. The deferred B3.11 follow-up (no manifests for the
four other add-ons) is pre-existing and explicitly out of scope for this change.

## What I verified

1. **314/314 node tests pass, 0 fail, on python3 3.9.6.** The Python leg refuses by
   design ("python3 is 3.9, but the add-on declares >=3.10") and runs zero Python tests —
   exactly as the validation record states. (`PYTHONDONTWRITEBYTECODE=1 npm test`.)

2. **Working tree clean at 49fc3ec.** `git status` empty; `git diff --check c673111..49fc3ec`
   clean. The fix-up changes only `lib/installer.js` (+15) and
   `test/bidding-addon.test.js` (+55); `skills/parley-bidding/` is untouched since
   714712f (the only commit that ever touched it).

3. **P1 — omitted `skill` identity (codex-1 round-6 MAJOR), fixed.** Runtime probe: full
   install, delete `skill` from core + bidding markers. Both units `malformed`,
   `managed: false`, `doctor.ok: false`; `install --only parley-bidding` and `uninstall`
   both `ok: false`. The identity check at lib/installer.js:1355 no longer exempts
   `undefined`; `installerOwnsDestination` at lib/installer.js:1652 requires exact
   equality. One answer across health and both mutations.

4. **P2 — `selected` from the core marker (kimi-1 round-6 MINOR), fixed.** Runtime probe:
   full install then `install --force --only parley-design`. Unflagged `doctor`:
   `parley-bidding: selected: false, valid-unselected, ok: false`. `doctor --only
   parley-bidding`: `selected: false, valid-unselected, ok: false`. Same recorded fact
   from both commands. The scoped read still inspects only `[parley-deck, parley-bidding]`
   — the cycle-5 narrowing is unchanged.

5. **P3 — 2.0.0-shape markers still valid.** Runtime probe: full install, then strip
   `markerSchema` from every marker (keeping `skill`). All six units report `valid`,
   `managed: true`, empty problems. The stricter identity predicate refuses no released
   marker shape. (`doctor.ok` is false on this machine only because the bidding runtime
   probe finds python3 3.9.6 and the add-on declares >=3.10 — the designed runtime gate,
   not a marker problem; confirmed by re-running with a fresh install and no tampering.)

6. **Payload integrity.** `skills/parley-bidding/` is 48 files, zero caches
   (`__pycache__`, `.pytest_cache`, `node_modules`, `.DS_Store`). The manifest aggregate
   is `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d` across 47
   files — matching every prior round's measurement. All 16 JSON files parse. `npm pack
   --dry-run` reports 202 total files.

7. **Diff scope.** `git diff c673111..49fc3ec` touches only `lib/installer.js` and
   `test/bidding-addon.test.js`. No drive-by changes to the payload or unrelated code.
