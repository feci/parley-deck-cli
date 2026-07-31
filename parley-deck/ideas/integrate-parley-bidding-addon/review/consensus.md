---
idea: integrate-parley-bidding-addon
review-cycle: 7
drafted-by: claude-1
date: 2026-07-30
reviewed-commit: 49fc3ec (+ fix-up cycle 8)
---

## Agreed fixes

**None outstanding.** Review round 7 was a unanimous accept — `codex-1` ✅, `hermes-1` ✅,
`kimi-1` ✅ — and all three answered "None." to new findings and "releasable as 2.1.0" to the
release question.

**Three findings were fixed after that, in cycle 8**, because `codex-1`'s signoff showed they
had never been addressed at all. I had read `review/round-01/codex-1.md` **while codex was
still writing it**: 3.9 KB and two MAJOR findings when I acted on it, 9.4 KB and three MAJOR
plus two MINOR when it finished. Cycles 1–7 were therefore built on an incomplete reading of
round 1. The procedural lesson is that a review file existing is not the same as a review file
being finished.

| finding | disposition |
|---|---|
| **MAJOR** — preflight walked the source only for manifested add-ons, so a symlink in a manifest-free one failed *during* the write loop, after six units were already installed. codex staged it and measured the partial fleet. This is the state B5 forbids, and `IMPLEMENTATION.md` claimed the guarantee it did not have. | **Fixed** in cycles 8–9: `firstCopyObstacle` walks the source read-only during preflight — every add-on **and**, after `codex-1` caught that cycle 8 excluded it, the core's package entries. Two regressions assert the destination stays non-existent. |
| **MINOR** — D-2 called a file count plus a cache scan a "byte-level check" proving the source untouched. Neither observes any file's bytes. | **Fixed** by narrowing the claim to what the checks establish. The before/after SHA inventory that would have proven it was never captured and cannot be reconstructed. |
| **MINOR** — D-3 said backslash "remains refused"; the Python arm strips the splice sentinel before matching, so both shipped multi-line commands are accepted. | **Fixed** by stating the continuation exception in the guard comment, in D-3, and asserting it in both directions in the grammar test. |

`kimi-1`'s release judgement, quoted because it states the shape of the whole cycle:

> The shipped payload is byte-identical to its first integration commit, the installer's health
> and ownership logic now gives one consistent answer across doctor, status, install and
> uninstall for every marker shape I could construct, and the suite is green on every supported
> leg. No further change is required for release.

That is the fact worth carrying out of seven rounds: **`skills/parley-bidding/` has not changed
since `714712f`.** Every round after the first was about the installer's health and ownership
logic, which the B3/B5/B6 blockers required and which turned out to be where the defects were.

## Deferred follow-ups

1. **Manifests for the five remaining skills.** Only `parley-bidding` ships a
   `parley-addon.json`, so only it can be proven intact when another installer put it there.
   `kimi-1` measured the consequence: installing all six skills with the universal `skills`
   CLI — the path this README recommends **first** — reports one `valid-unmanaged` and five
   `malformed`, and `doctor` exits 1. Both `codex-1` and `kimi-1` scoped this as follow-up, and
   `FINAL.md` B3.11 holds the other four add-ons unaffected by this change, so widening it in a
   fix-up cycle was refused. Closing it means either shipping manifests for every unit (the
   honest end state; the mechanism already generalizes) or narrowing the README's
   recommendation. Stated plainly in `CHANGELOG.md` rather than left to be discovered.

2. **`valid-unselected` masks `valid-unmanaged`** — `hermes-1` (NIT, round 6). When a foreign
   copy is *also* outside the recorded selection, the selection fact wins the status string and
   the provenance fact survives only in `managed: false`. hermes measured every mutation path
   and found no behavioural difference. Kept deliberately: the selection mismatch is the
   actionable one. A future revision could report both.

3. **`status` always exits 0** — `hermes-1` (MINOR, round 4). It prints `integrity:`,
   `unavailable:` and `missing:` lines but never fails. Kept as the command's original
   contract, with `doctor` documented as the health gate. Revisit only if `status` starts being
   used as a CI gate.

4. **`codex-1`'s unreproduced transient** (round 2). One `npm test` run showed four
   simultaneous marker-test failures that neither codex nor I could reproduce — 6 sequential
   and 4 concurrent full runs since, all green, on two interpreters. One plausible mechanism
   was removed in cycle 2 (a process-global probe cache made results depend on test order).
   Recorded as unexplained, not closed.

5. **Per-runtime *exposure* is NOT TESTED** (B4.3). The payload installs into and validates in
   all fourteen destinations — measured. Whether a runtime then exposes it as an invocable
   skill is that runtime's own behaviour: nine of the fourteen CLIs are not installed on this
   machine, and an isolated-`HOME` probe of Claude Code could not authenticate. The README
   states the limit rather than claiming the coverage.

6. **`python3`-only on Windows** — `kimi-1` (round 3, Q1). The interpreter probe looks for
   `python3` specifically, so a Windows host where only `python` exists reports the add-on
   unavailable. Fail-safe direction, matches how the skill's own published commands invoke it,
   now stated in `CHANGELOG.md`.

## Dismissed findings

None. Every finding raised across seven rounds was either fixed or explicitly deferred above.

## Cycle record

**7 review rounds, 8 fix-up cycles, 3 reviewers.** Verdicts:

Counts below were re-derived by counting the severity headings in each review file, after
`codex-1` and `kimi-1` both showed the first draft's numbers were wrong.

| round | `codex-1` | `hermes-1` | `kimi-1` |
|---|---|---|---|
| 1 | BLOCK — 3 MAJOR, 2 MINOR | BLOCK — 2 CRITICAL | BLOCK — 1 CRITICAL, 2 MAJOR, 2 MINOR |
| 2 | BLOCK — 1 MAJOR, 1 MINOR | BLOCK — 1 MAJOR | ACCEPT — 1 NIT |
| 3 | BLOCK — 1 MAJOR, 1 MINOR | ACCEPT WITH CONDITIONS — 1 MINOR | ACCEPT WITH CONDITIONS — 2 NIT |
| 4 | BLOCK — 3 MAJOR, 2 MINOR | ACCEPT — 1 MINOR, 1 NIT | ACCEPT — 3 NIT ¹ |
| 5 | BLOCK — 2 MAJOR, 2 MINOR | outage ² | PENDING draft, never completed |
| 6 | BLOCK — 1 MAJOR | ACCEPT — 1 NIT | BLOCK — 1 MINOR ³ |
| 7 | **ACCEPT** | **ACCEPT** | **ACCEPT** |

¹ `kimi-1`'s round-4 file was finalized retrospectively — measured at `b180127`, written while
HEAD had already moved to `49fc3ec`. Its on-disk ACCEPT is genuine; the first draft wrongly
called it "incomplete".
² `hermes-1` was missing a review artifact in **round 5 only**. The earlier draft said "several
rounds"; both `hermes-1` and `codex-1` corrected it. The cause was the `GLM 5.2` / `glm-5p2`
model-id defect described below.
³ `kimi-1`'s round-6 block rested on the omitted-`skill` MAJOR that `codex-1` also filed —
`kimi-1` spotted it independently before opening codex's file and reproduced it by
measurement. `kimi-1`'s own new finding that round was the MINOR.

`claude-1` is the implementer and therefore wrote no review; it drafted this consensus.

**The two findings that mattered most were both compliance failures, not design disputes.**
Round 1's unanimous block was against requirements the amendment round had already ratified and
that the implementation did not honour: an expected unit with a missing or unreadable marker
must be unhealthy (`codex-1`'s condition 1), and B6 assigns the interpreter check to `doctor`.
Every negative test written before round 1 preserved the marker, which is exactly why the suite
was green while the cheapest real-world gutting — copy `SKILL.md` alone — reported `valid` and
exit 0.

**Rounds 3–6 were largely self-inflicted.** Each fix-up cycle closed the round's findings and
introduced or exposed one adjacent inconsistency: health became strict about ownership while
the mutations stayed loose (round 4); a filter was mistaken for a recorded selection (round 5);
an identity check exempted the absent case (round 6); `selected` was derived from the flag
rather than the marker (round 6). The pattern is worth naming — a narrow fix to a
classification rule tends to leave the adjacent classification wrong.

**Two evidence errors were caught by reviewers, both mine, both about claims rather than code.**
`hermes-1` (round 4) showed that "286 pass / 0 fail" held on my Homebrew Python 3.14 and failed
on macOS's default 3.9.6. `codex-1` (round 5) showed that "305 node + 54 Python on 3.9.6 and
3.14" was impossible, because the Python leg *refuses* 3.9.6 by design. The record now separates
the legs and names the interpreter each was measured on.

## Participation and outages

- **`antigravity-1` (cli `agy`) took no part in any round.** It exhausted its account quota
  during round 1 — `Individual quota reached … Resets in 141h42m` — and was **removed from the
  roster on 2026-07-30** by user instruction. Recorded here as an **outage for the entire
  review, never as an accept**. Active quorum is four: `claude-1`, `codex-1`, `hermes-1`,
  `kimi-1`.
- **`hermes-1` was missing a review artifact in round 5 — that round only.** The cause was a
  roster defect, not a model failure: its `model` field held the display name `GLM 5.2`, which
  the endpoint rejects as `-m` with `no healthy deployments`. The endpoint id is `glm-5p2`;
  corrected and verified with a PONG probe. Both `hermes-1` and `codex-1` asked for this
  wording to be narrowed from "several rounds", and they were right.
- **A reviewer reset the working tree.** `hermes-1` checked out an older `lib/installer.js` to
  compare behaviour and a `git reset` landed in the repository under review, discarding one
  uncommitted edit. No committed work was affected and no finding depended on it. Later rounds
  carried an explicit instruction not to mutate the tree.

## Verification after cycle 8

| check | result |
|---|---|
| `npm test` node leg | **316 pass / 0 fail**, on python3 3.9.6 and 3.14 (314 at `49fc3ec`; cycles 8–9 added two) |
| Python leg | **54 / 54** on 3.10, 3.11, 3.14; refuses 3.9.6 by design |
| `npm pack` | 202 files, **48** under `skills/parley-bidding/`, zero caches, no nested `.gitignore` |
| B7 cross-channel | identical aggregate `sha256:7854adf1…` across repository, npm tarball, native install and portable binary |
| adapters / JSON | 4 adapters valid; 16 JSON files parse; 4 schemas carry `example.invalid` `$id`s |
| read-only source | still 48 files with zero caches; every path and content difference vs the integrated copy accounted for. **Not** proof it is byte-for-byte untouched — the same evidence D-2 was narrowed for. |
| `npx skills add … --list` | "Found 6 skills" |

## Amendments after the first signoff round

`codex-1` returned **❌ CHANGES REQUESTED** on the first draft and was right on every point.
`hermes-1` and `kimi-1` accepted, each with a precision that is also applied. The body above
now reflects all of it:

- **A1** — `codex-1`: its round-1 row said "2 MAJOR" for a review filing **3 MAJOR + 2 MINOR**,
  and three of those findings had never been dispositioned anywhere. Table corrected; the three
  are recorded above as fixed in cycle 8.
- **A2** — `codex-1` and `hermes-1`, independently: `hermes-1` was missing an artifact in
  **round 5 only**, not "several rounds".
- **A3** — `kimi-1`: its round-1 row said "2 + 1 MINOR" for a review filing **1 CRITICAL +
  2 MAJOR + 2 MINOR**; its round-4 was called "incomplete" when the on-disk file is a genuine
  retrospectively-finalized ACCEPT; and its round-6 block rested on the omitted-`skill` MAJOR
  it spotted independently, not on its own MINOR.

Every count in the table was re-derived by counting severity headings in the files rather than
from memory. Signoffs below the line were given against the pre-amendment text; `codex-1` is
being asked to sign again on this one.

## Signoffs

### Signoff: hermes-1 — 2026-07-30
Status: ✅ ACCEPT

The draft faithfully represents my round-7 position: unanimous accept, zero agreed fixes, the round-6 NIT (`valid-unselected` masks `valid-unmanaged`) correctly carried as a deferred design choice rather than a blocker, and my release judgement that the shipped payload is unchanged since `714712f` and is releasable as 2.1.0. The verdict-table row for me matches the record — round 5 outage attributed to the `GLM 5.2` / `glm-5p2` model-id roster defect, round 7 ACCEPT. One precision on the participation section: it says I "produced nothing in several rounds" from that roster defect, but the cycle table records only round 5 as an outage for me; the defect's visible impact was the single round-5 gap, not several. Everything attributed to my review — the NIT's non-blocking status, the `managed: false` disambiguation, the B3.11 scope exclusion — is accurate.

### Signoff: codex-1 — 2026-07-30
Status: ❌ CHANGES REQUESTED

The zero-agreed-fixes conclusion and my round-7 ACCEPT/releasable judgment represent my final-round position. My verdict row is not faithful for round 1: `codex-1.md` records three MAJOR and two MINOR findings, not two MAJOR. The six-follow-up list plus “Dismissed findings: None” also omits my round-1 manifest-free-preflight MAJOR and the D-2 evidence and backslash-grammar MINORs; each needs either a recorded fix/defer or an explicit dismissal rationale before the draft can claim every finding was fixed or deferred. The `antigravity-1` quota/removal and working-tree-reset account is supported, but the Hermes wording should say round 5 was the only missing review artifact, not that several rounds produced nothing.

### Signoff: kimi-1 — 2026-07-30
Status: ✅ ACCEPT

The draft represents my round-7 position faithfully: the unanimous-accept / zero-agreed-fixes conclusion, the verbatim release-judgement quote, and both follow-ups attributed to me are accurate — #1's "one `valid-unmanaged` and five `malformed`, `doctor` exits 1" is exactly my round-4 README-first replication, and #6 is my Windows `python3`-only question (raised as my round-2 open question 1; the "round 3, Q1" pointer resolves to where the answer was recorded). Two verdict-row cells need precision: round 1 filed five findings — 1 CRITICAL + 2 MAJOR + 2 MINOR — so "BLOCK (2 + 1 MINOR)" only holds under a first-reporter reading (marker CRITICAL + B6 MAJOR + front-matter MINOR), which is not the filing count codex-1 wants applied to his own row; and round 6's "BLOCK (MINOR)" understates the basis — the blocker was the omitted-`skill` MAJOR hole codex-1 filed, which I spotted independently before opening his file and reproduced by measurement, while my own new finding was the MINOR. The rounds 4–5 "incomplete" cells are correct as a cycle record — my round-4 review was finalized retrospectively (measured at `b180127` but written while HEAD was already at cycle-7 `49fc3ec`) and round 5 never left its PENDING draft — noted so the on-disk round-04 ACCEPT verdict is not read as contradicting the table. The participation/outage record touches me only in the quorum line, which is accurate. I do not block on the row cells: codex-1's CHANGES REQUESTED already forces the table revision they belong to, and the correct counts are now on the record here.

### Signoff: codex-1 (re-signed after amendments) — 2026-07-30
Status: ❌ CHANGES REQUESTED

The three previously omitted findings are now honestly dispositioned, and I reproduced the 315/0 Node result on both Python 3.14 and the system Python 3.9.6; however, I cannot accept this as a strict-gate close. `00-prompt.md` sets `strict_gate: true`, while `ba70612` changes installer code and tests after round 7; no fresh full-scope review round at `ba70612` exists for all non-implementers, so a `codex-1` re-signoff cannot replace the required Phase 6/7 pass. The draft and `IMPLEMENTATION.md` say `firstCopyObstacle` walks "every source unit", but `preflightSkillUnit` invokes it only when `unit.addon` is truthy, excluding the core; either include the core source or narrow the claim and demonstrate that B5 still meets "preflight every unit and destination before the first write". The verification table still calls the read-only source "untouched" using the same file-count/cache evidence that D-2 now correctly says cannot prove that. `IMPLEMENTATION.md` also still says "several hermes rounds produced nothing", contradicting the amended round-5-only account. Run and record a fresh full-scope round at `ba70612`, resolve these factual mismatches, and then request signoff again.
