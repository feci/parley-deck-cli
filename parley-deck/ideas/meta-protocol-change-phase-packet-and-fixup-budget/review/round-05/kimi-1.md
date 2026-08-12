---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 5
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 5 — release candidate 1.44.0 + 2.8.0

## Summary

The verdict is NOT CLEAN, and it is pinned to the candidate this round was called on: **the
fix-up-cycle-4 state** (sha256 `impl.go e6aa4c7d…`, `driver.go a88f7ddb…`, measured by me this
session; companion pins from my round-3 file: `cursor.go 0df19271…`, `consensus.go 1330db98…`,
`impl_test.go 6418965c…`, CLI `CHANGELOG.md f9247104…`).

Against that candidate: the three cycle-4 fixes are genuinely fixed — I verified each with my own
probes and reversions in an isolated copy (Q1). The scope call is the right release decision (Q2).
But the candidate is not shippable as cut, for two demonstrated, in-scope reasons:

1. **[MAJOR] The shipped test suite is provably blind to all three cycle-4 fixes.** I reverted
   each fix in an isolated copy; the full suite (all 28 packages) stayed green under each
   individual reversion **and under all three combined** (exit 0). Only out-of-tree probes catch
   them. Cycle 4 is the first cycle in this idea to change the safety boundary with zero
   regression tests and zero recorded reversion checks — on the boundary whose entire history is
   silent regressions.
2. **[MAJOR] The "corrected" limits still overstated the implementation.** The S4 CLI CHANGELOG
   said "The unit is now completed fix-up cycles" and "Losing one record does not lower the
   count" — both false under the reservation semantics cycle 4 itself introduced (the
   reservation→marker window holds exactly one record). Quotes in Q3.

Both findings converge with @codex-1's round-5 file, which I had not read when my evidence was
captured (timeline below); the same two points were found a third time in my own late-arriving
round-04 file (parallel session). Three independent derivations, same two defects.

**State of the tree (not the verdict object):** the tree moved again during this review — fix-up
cycle 5 (16:24–16:26) and cycle 6 (16:39–16:42) landed while I worked, and both findings above
are addressed in the current tree. I spot-verified that (Q5), but this file certifies only the
pinned S4 candidate; certifying the successor belongs to the round now in progress, against a
tree that holds still.

## Review-validity event (read before any finding)

**[PRIMARY — CONFIRMED]** Third consecutive round in which the candidate mutated under my review.
Timeline (mtimes; my session opened ~16:05):

```text
16:05   round-05/ exists, EMPTY; the brief describes the cycle-4 candidate
~16:06–16:12  I read IMPLEMENTATION.md (through cycle 4), the round-04 review files,
              both CHANGELOGs, and the full diff; go build/vet green
~16:12  go test ./... green on the tree (S4); isolated copy made (rsync, no .git)
~16:13–16:19  my 6 probes written and run green against S4; 3 individual reversions,
              each red against my probes for the intended reason, each restored and
              hash-re-verified
~16:20–16:23  suite-blindness demonstration: full suite green under each individual
              reversion AND under the combined 3-way reversion (exit 0)
16:22   review/round-05/codex-1.md lands — NOT read by me until after my demonstration
16:24–16:26  impl_test.go, impl.go, both CHANGELOGs, IMPLEMENTATION.md rewritten (cycle 5)
16:25   review/round-05/hermes-1.md lands
16:35   review/round-04/kimi-1.md lands (parallel session, late)
16:38   review/round-06/codex-1.md lands
16:39–16:42  cursor.go, impl.go, impl_test.go, CHANGELOG.md, consensus*.go, IMPLEMENTATION.md
             rewritten (cycle 6)
```

Everything I certify about the verdict object ran against S4 or its isolated copy before
16:24. Everything I say about cycles 5–6 is tagged as a later spot check. My footprint on both
working trees is zero: all mutations ran in `mktemp -d` copies; no git write commands; the
`git status` file-set is unchanged by me.

## Q1 — are the three fixed defects actually fixed?

**Yes — all three, verified by my own probes and reversions in the isolated copy.** Six probes
(`internal/driver/r5_probe_test.go`, same-package, isolated copy only), all PASS against S4:

- **Reservation ordering** (`impl.go:317-328` on S4): `saveCursor` precedes `Fixup`. Probe 1a: an
  errored `Fixup` escalates and leaves `FixupCyclesPublished=1` durable on disk; a subsequent
  healthy `Advance` is asked for **cycle 2**, not cycle 1 — the errored cycle never comes back.
  Probe 1b: checks failing after a successful `Fixup` likewise leave the spend persisted.
  *Reversion:* moving the save back after `Fixup` turns 1a red ("no persisted cursor after errored
  fix-up") — the exact round-4 leak.
- **AF2 at equality** (`impl.go:147` on S4, `spent > MaxFixupCycles`): probe 2a — cap 5, five
  markers, marker at the current round → AF2 archives and opens round 6, no `Fixup` re-run, no
  escalation. Probe 2b — six markers against cap 5 → escalates, no round opened. *Reversion:* `>`
  back to `>=` turns 2a red with precisely @codex-1's round-4 failure ("legitimate AF2 recovery at
  the cap was stranded … budget exceeded: 5 cycle(s) recorded against MaxFixupCycles=5").
- **Cursor read error** (`driver.go:232-248` on S4): probe 3a — a corrupt `driver.json` escalates
  ("cannot read the driver cursor … unknown fix-up budget") with zero `Impl` calls; probe 3b — a
  genuinely absent cursor proceeds as a fresh run. *Reversion:* restoring the `err == nil &&`
  swallow turns 3a red — a fix-up ran on an unknown count (`action=fixup err=<nil>`).

After each reversion I restored the files and re-matched sha256 to the working tree; the final
isolated state was hash-identical to S4 and fully green.

## Q2 — the scope call: is "correct the claims + name a follow-up" an acceptable disposition?

**Yes — as a release decision, not a preference.** The stronger property does not have to ship in
this release, and requiring it would be the error, not the caution.

- **The ratified contract is two numbers and their enforcement.** `FINAL.md` §2 ratifies "5
  inclusive published cycles" and "3 rounds after round 1", blocking escalation, never auto-close,
  and three binding implementation conditions (same-patch text/code; inclusive boundary *with a
  test*; the deliberation cross-review wiring). Nothing in it ratifies an idea-scoped ledger,
  cross-run serialization, or an out-of-repository trust anchor. "An extension … never resets the
  count" governs *extensions*, not adversarial state deletion.
- **What breaks if it does not ship: nothing ratified.** The disclosed residual is real — a
  participant with workspace write, a deleted run directory, a rollback, or two concurrent runs
  can reduce or duplicate the count — so the cap is not a security boundary. But the cap's
  ratified purpose is to interrupt pathological *churn*, and every careless path to that is
  closed in S4 (Q1, and my round-3/round-04 files' probes). A participant with workspace write
  has cheaper attacks on this release than the counter (the driver source and the tests are
  equally writable); the deck's actual trust boundary is git history plus human review — which is
  exactly how all four rounds of budget defects were caught.
- **What breaks if it DOES ship under this mandate:** the deck's file-based portability — a fresh
  clone would lose an out-of-repo count entirely, reintroducing fail-open at a different layer —
  and the scope discipline the FINAL itself enforces. Silently building an unratified trust
  anchor under a two-number mandate is the same class as silently expanding any other ratified
  boundary. And it would hold this release's live, demonstrated fixes (inclusive caps — `fast`
  published **zero** cycles against a printed 1 — and the BLOCK back-edge bypass) hostage to an
  undesigned system.

One necessary condition on "acceptable": the disposition is "correct the claims **and** name the
follow-up". In S4 the claims were only half-corrected (Q3) — the disposition was right, its S4
execution was not. `fixup-budget-trust-anchor` is the right named follow-up.

## Q3 — do the corrected CHANGELOG limits accurately describe what ships?

**Not in S4. Two sentences still overstated the implementation, both quotable:**

> "The unit is now completed fix-up cycles." (S4 CLI CHANGELOG, line 16)

False under cycle 4's own semantics: the cycle is charged at *reservation*, before the
code-writing call — an attempt that errors without writing a line is charged. "Completed"
describes the marker population only. (I flagged this tension myself before reading @codex-1's
round-5 file; he names the same sentence.)

> "Losing one record does not lower the count" (S4 CLI CHANGELOG, lines 27-28)

False in the window cycle 4 created: `impl.go:322-336` (S4) persists the reservation at :323,
calls `Fixup` at :326, and writes the marker only at :336 after checks pass. Between reservation
and marker — which includes **every errored attempt** — the cursor is the *only* record, and
losing it there loses the count. @codex-1's round-5 composition (reserve → error → delete
`driver.json` → advance re-enters `Fixup`) is correct by my source inspection [SECONDARY origin,
PRIMARY confirmation]; the S4 "What this is not" paragraph's "a deleted run directory … can still
reduce … the count" covers it only if the reader already knows which window applies, while the
unqualified "losing one record does not lower the count" asserts the opposite.

The skill CHANGELOG's "robust against accidental loss of its own records" inherited the same
window. Everything else load-bearing in S4 checked out against my probes: the max-of-two
mechanism, forge-only-raises, unreadable-cursor-escalates, AF2 consults the cap yet finishes the
last allowed cycle, the "not a security boundary" list, the two §4.0 cells, the BLOCK back-edge
disclosure for `standard`, and the "Not in this release" section.

## Q4 — anything NEW in the round-4 diff that was not there before?

The S4 delta over cycle 3 is exactly: the three fixes (driver.go switch; impl.go reservation
order and AF2 operator + comment), the two CHANGELOG disclosure paragraphs, and the
IMPLEMENTATION.md cycle-4 section. No test file was touched (`impl_test.go` mtime 15:41 predates
cycle 4). Three things deserve naming:

1. **[MAJOR — the release blocker] No permanent tests, and the suite is blind.** Demonstrated,
   not argued: with all three fixes reverted simultaneously in an isolated copy,
   `go test ./... -count=1` over the unmodified shipped suite exited 0 (all 28 packages). Each
   individual reversion likewise left the suite green; only my out-of-tree probes went red. This
   also voids cycle 4's recorded verification ("go build, go vet, full Go suite green") as
   evidence *for the fixes* — the suite passes with or without them. Cycles 1–3 each shipped
   pinning tests plus recorded isolated-copy reversion checks; cycle 4 broke that discipline on
   the fourth consecutive change to the same boundary. The AF2 `>=` defect was itself
   *introduced* by cycle 3 while fixing another defect, and nothing in the suite noticed. That
   is the inert-test class this deck has now paid for in three consecutive ideas.
2. **[MAJOR — folded into Q3] A genuinely new state:** the reservation→marker single-record
   window did not exist before cycle 4 (previously the save followed `Fixup`, so a failed attempt
   simply left no trace anywhere). The fix was correct; the S4 text describing it was not (Q3).
3. **[NITs, carried/created]** the call-site comment at S4 `impl.go:297-298` still named
   `IMPLEMENTATION.md` as the count source (my round-3 F3; @hermes-1's round-4 NIT); the
   `LoadCursor` doc comment at S4 `cursor.go:82-83` ("a missing or corrupt file is non-fatal")
   became false for the `Advance` path precisely because of cycle 4's fix; the fix-up escalation
   message printed the *refused* cycle as the count ("after 6 cycle(s)" at cap 5 — my round-3
   note; parallel round-04 file F2). All three were comment/diagnostic-level.

## Q5 — should this release ship?

**The S4 cut: no.** Two demonstrated MAJORs in the candidate as briefed — the suite cannot see
any of the three fixes regress, and the release notes still overstated the failure domain. Both
are narrow, both are named above, and neither touches the scope call.

**The tree has already answered.** While this review ran, fix-up cycle 5 landed
(`TestRound4FixesHavePermanentTests` — four subtests covering exactly my probes 1a/2a/2b/3a/3b;
the doc comments and both CHANGELOGs rewritten onto "reserved/charged" vocabulary with the
single-record window stated plainly), and cycle 6 followed (the `New`→`HardCrossReviewCap` seam
test, the remaining unit-wording sweep, and fixes for my round-3 F5 and the F2/F3/F4 NITs).
**[PRIMARY — CONFIRMED, spot checks by me on the current tree]** — full suite green (exit 0);
the new tests pass; and under the same combined three-way reversion in a fresh isolated copy,
exactly three of the four new subtests now go red (the fourth is green under both operators, as
it should be) — the shipped suite can see the regressions. Current pins at my run:
`impl.go 82d5c49d…`, `driver.go a88f7ddb…`, `cursor.go f1f45442…`, `consensus.go a389a8a9…`,
`impl_test.go a630cea9…`, CLI `CHANGELOG.md 0a4d2805…`.

That is a spot check, not a certification: the tree last moved at 16:42 and a round-06 review
already exists. The release decision now belongs to a full pass over a **frozen** successor
candidate — the same demand my round-3 and round-04 files made, still unmet as a process matter.

## Findings index

| Severity | Finding |
| --- | --- |
| MAJOR | S4 shipped the three cycle-4 safety fixes with no regression tests; the full shipped suite stays green (exit 0) under each individual reversion and under all three combined. Demonstrated in an isolated copy. |
| MAJOR | S4 CLI CHANGELOG still overstated the implementation: "The unit is now completed fix-up cycles" (line 16) and "Losing one record does not lower the count" (lines 27-28) are false under the reservation semantics cycle 4 introduced (reservation→marker window is single-record). Skill CHANGELOG inherited the window. |
| NIT | S4 `impl.go:297-298` call-site comment still names `IMPLEMENTATION.md` as the count source (carried from rounds 3-4); S4 `cursor.go:82-83` "corrupt file is non-fatal" made false by cycle 4's own fix; fix-up escalation message reports the refused cycle as the count. All comment/diagnostic-level. |

## Responses to other reviewers

**@codex-1.** Your round-5 NOT CLEAN is right on both counts, and my two MAJORs are the same two
defects, derived independently: my probes were written and my reversion demonstration was complete
before I opened your file (timeline above). Your reservation-window composition is the sharper
half of the claims finding; I confirmed it by inspection rather than rerunning it. On Q2 we agree
outright: the trust-anchor property is a different, unratified system; corrected claims plus
`fixup-budget-trust-anchor` is the right disposition; and the S4 claims were not yet corrected.

**@hermes-1.** Your round-5 CLEAN examined the same three fixes and verified them correctly — our
Q1 evidence agrees probe for probe. On Q3 we differ: you read "completed … cycles" as a minor
tension "immediately clarified" by line 26; I read it, and the unqualified "losing one record does
not lower the count", as factually wrong in the reservation window the same diff created. On this
deck's own rule — a claim about the safety mechanism must be re-checked every time the mechanism
changes — the S4 text had not been. The current tree's wording (which postdates your review)
concedes the point: it now says "deliberately not 'a completed cycle'" and states the
single-record window explicitly.

**My round-04 file (late-arriving, parallel session).** Its F1 is my MAJOR-1, its Q5 pre-cleared
the corrected texts I only spot-checked. Convergent and independent; corroboration, not
coordination.

## Validation record

All commands run by me this session; read-only in the trees, mutating only in `mktemp -d` copies
(deleted after use; Go caches under the repo's own `.gocache`/`.gomodcache`):

```text
TREE (S4, ~16:06–16:12):  go build ./... / go vet ./...        exit 0
TREE (S4, ~16:12):        go test ./... -count=1               28 packages ok, zero FAIL
ISOLATED (S4 copy):       6 probes (r5_probe_test.go)          all PASS on shipped code
  reversion: reservation moved after Fixup   → probe 1a RED (no persisted cursor)
  reversion: AF2 `>` → `>=`                  → probe 2a RED (recovery stranded at 5/5)
  reversion: cursor-error swallow restored   → probe 3a RED (fixup on unknown budget)
  each reversion individually: full suite GREEN (exit 0)
  ALL THREE COMBINED:       go test ./... -count=1             exit 0 — suite blind
  restored after each; sha256 re-matched to S4; package re-verified green
SKILL (isolated copy):    build-addon-manifest.js --check      all six ok
                          npm test                             386 pass / 0 fail; 54 python OK
TREE (current, spot check, ~16:41): full suite exit 0;
  TestRound4FixesHavePermanentTests + TestUnreadableMarkerEscalates…  PASS
ISOLATED (current copy):  combined 3-way reversion             → exactly 3 of 4 new
                          subtests RED (reservation, AF2 equality, corrupt cursor)
```

No git write commands were issued in either repository; neither working tree was modified by me.
