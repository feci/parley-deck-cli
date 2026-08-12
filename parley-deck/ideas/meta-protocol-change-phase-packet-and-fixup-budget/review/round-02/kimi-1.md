---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 2
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review, round 2 — release candidate 1.44.0 + 2.8.0

## Summary

**[PRIMARY — CONFIRMED]** All six round-1 findings are fixed or honestly withdrawn, and I
re-produced every check myself: the full Go suite is green (27/27 packages, including
`internal/runner`, which failed only in codex-1's sandbox), the skill package now installs
(`npm test` 386/386 vs 271/115 in round 1), both reversion claims go red exactly as claimed, and
the 69-idea distribution behind "why 5" reproduces digit-for-digit. Against the fix-up-cycle-1
snapshot this candidate was close to shippable.

**[PRIMARY — CONFIRMED]** Two things stop the release anyway. First, a new MAJOR introduced by the
round-1 fix itself: the fix-up budget is now counted from `## Fix-up cycle` headings in the
implementer-owned `IMPLEMENTATION.md` — an unvalidated, author-controlled source that is fail-open
in four distinct ways, so the ratified "never resets the count" guarantee is not delivered (Q3).
Second, a review-validity event: **the working tree mutated while this review was running.** At
15:07 CEST `internal/driver/impl.go` and `impl_test.go` were rewritten (fix-up cycle 2, replacing
the heading count with driver-owned `.fixup-done` markers), and the tree is currently RED. My
green/reversion evidence covers the pre-mutation snapshot (hashes below), not the current tree.

Disclosure: in round 1 I returned CLEAN without having examined the packaged skill or the
consensus BLOCK back-edge — the two places the CRITICALs were. Both were examined first this time.

## Review-validity event (read before any finding)

**[PRIMARY — CONFIRMED]** Timeline (file mtimes, `ls -lT`, and my own shell history; my session
opened 14:44 CEST):

```text
15:00  review/round-02/hermes-1.md written
15:03  review/round-02/ directory created
15:05  review/round-02/codex-1.md written
15:07:01  internal/driver/impl.go rewritten      (NOT by me — see below)
15:07:18  internal/driver/impl_test.go rewritten (NOT by me)
15:07:52  my byte-exact restore of internal/driver/consensus.go after my reversion check
```

I did not read either round-02 review before writing this one; my findings are independent. The
new `publishedFixupCycles` comment in the current `impl.go:490-505` says "Review round-02 showed
that derivation was fail-open in four ways at once" — i.e. the implementer is already applying a
fix-up cycle 2 in response to codex-1, and it enumerates the same four fail-open modes I
independently derive under Q3 below. Convergence, not contamination.

**[PRIMARY — CONFIRMED]** My own footprint on the tree is zero. My reversion checks (Q4) were
performed on copies restored byte-exactly: sha256 before/after `77142dfe…` for `impl.go`,
`57c8e626…` for `consensus.go`; backups deleted afterward. The current `impl.go` /
`impl_test.go` are entirely the implementer's new versions — I never wrote over them after
15:07:01. One caveat I cannot exclude: if the implementer wrote to `consensus.go` inside the
~50-second window before my 15:07:52 restore, my `cp` would have overwritten it. Its current
content is byte-identical to what I reviewed (guard present, diff matches), so nothing looks lost,
but the implementer should re-confirm `consensus.go` matches their intent.

**[PRIMARY — CONFIRMED]** Current tree state: `go build ./...` and `gofmt` clean, but

```text
$ go test ./internal/driver -count=1 -run TestFixupCapIsInclusive
--- FAIL: TestFixupCapIsInclusive/cap_5:_the_5th_published_cycle_is_allowed  ranFixup=false want true (action=fixup err=<nil>)
--- FAIL: TestFixupCapIsInclusive/cap_5:_the_6th_escalates                   action=fixup want escalated
--- FAIL: TestFixupCapIsInclusive/cap_2:_the_2nd_is_allowed                  ranFixup=false want true
--- FAIL: TestFixupCapIsInclusive/cap_2:_the_3rd_escalates                   action=fixup want escalated
--- FAIL: TestFixupCapIsInclusive/cap_1:_the_2nd_escalates                   action=fixup want escalated
```

**[PRIMARY — CONFIRMED]** Cause, from reading the new code: the rewritten fixture
(`writeImplWithCycles`) plants a `.fixup-done` marker in every round 1..N *including the highest
review round*, which trips the pre-existing AF2 crash-idempotency branch (`impl.go`, the
`fixupMarker` check at the top of `advanceReview`): the driver then finishes the transition
without calling `Impl.Fixup` and without escalating. The new counter and the AF2 crash marker are
the same file, and at the current round those two meanings collide. Whether the finished fix-up
cycle 2 resolves this cleanly is for the next review — I have not verified it and certify nothing
about it.

## Findings

### [MAJOR] The fix-up budget is counted from author-controlled prose, fail-open in four ways

**[PRIMARY — WRONG]** Against the reviewed snapshot, the release claim "`MaxFixupCycles = N`
publishes cycles 1..N and escalates when N+1 would start" is not guaranteed, because the count it
compares against can be silently lowered or reset by the party the budget constrains. The reviewed
`publishedFixupCycles` (pre-mutation `impl.go`, sha256 `77142dfe…`) was:

```go
b, err := os.ReadFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
if err != nil { return 0 }
for _, line := range strings.Split(string(b), "\n") {
    if strings.HasPrefix(strings.TrimSpace(line), "## Fix-up cycle") { n++ }
}
```

**[PRIMARY — CONFIRMED]** Each mode, from the code itself:

1. **Missing or unreadable file → count 0 → the whole budget restarts.** `err → return 0`, so
   cycle becomes 1 regardless of history. Fail-open in the unsafe direction. (Mostly unreachable
   in production because Phase-8 ideas carry an `IMPLEMENTATION.md`, but "mostly" is doing the
   work here: any read error — permissions, transient I/O — is indistinguishable from "no cycles
   published".)
2. **A heading inside a code fence counts.** Line-prefix match with `TrimSpace`, no fence
   tracking; a documented example like ```` ## Fix-up cycle 9 ```` inside a ``` block spends real
   budget. Safe direction (early escalation), but false.
3. **The count is a heading COUNT, not an ordinal max.** Duplicate `## Fix-up cycle 2` headings
   count twice (over-count, safe direction); renumbering, merging or deleting a heading lowers the
   count and **buys another cycle** (unsafe direction). Careless tidying of `IMPLEMENTATION.md` —
   e.g. folding two cycles into one section — silently extends the very budget the cap exists to
   enforce.
4. **The source file is owned by the constrained party.** The implementer writes
   `IMPLEMENTATION.md`; the driver has its own tamper-evident record of completed cycles (the
   `.fixup-done` markers it writes itself after `Fixup` + the post-fix-up check gate) and does not
   consult it. `FINAL.md` requires that an extension "never resets the count"; this source makes a
   reset one careless edit away, with no in-loop signal.

**[PRIMARY — CONFIRMED]** One mitigating fact: the cap of 5 was calibrated with exactly this naive
count (`rg -c '^## Fix-up cycle'` over 69 ideas — reproduced below), so counter and calibration
are at least consistent; a stricter parser would not obviously improve the calibration. That
argues the *number* is fine; it does not fix the fail-open directions.

Suggested fix: exactly the direction the in-flight fix-up cycle 2 takes — count the driver's own
`.fixup-done` markers and escalate (never restart at zero) when the count cannot be determined.
When that lands it must be re-verified end to end, including the AF2 collision above, and the
"Wrong unit" paragraph of `CHANGELOG.md` must be rewritten, since it currently says the count
comes from `IMPLEMENTATION.md` headings.

### [MINOR] Neither changelog discloses that `standard`'s BLOCK back-edge is now bounded too

**[PRIMARY — CONFIRMED]** `track.PolicyFor(Standard)` has always carried `CapCrossReviewRounds:
2`; this release additionally wires it into `HardCrossReviewCap` (`driver.go:143-148`), so the
consensus-BLOCK back-edge (`consensus.go:95-99`) now stops a `standard` idea at round-03 where a
repeatedly-BLOCKed one previously could reach round-05 (`MaxRounds` default 4). The CLI
changelog's back-edge paragraph is framed entirely around `deliberation`; the sibling changelog
says nothing about it either. This is a real behavior change on the **default track**, beyond the
two ratified cells — correct in substance (Q2), but users of `standard` deserve one sentence in
the release notes. The implementer's own scope note in `IMPLEMENTATION.md` ("Fix-up cycle 1")
asked reviewers to confirm or reject exactly this.

### [NIT] The cross-review escalation message promises "the trajectory" that this release says does not exist

**[PRIMARY — CONFIRMED]** `consensus.go:97` escalates with "…escalating for human review of the
trajectory", while `rg -n trajectory internal --glob '*.go'` finds only that string and both
changelogs explicitly list the structured trajectory payload as **not implemented**. The human
who receives this escalation gets counts and the cap, not a trajectory. Reword the message, or
ship the payload.

## Answers to the requested questions

### Q1 — is each round-1 finding fixed, or relocated?

**[PRIMARY — CONFIRMED]** All six addressed; none relocated to another layer:

1. **CRITICAL (BLOCK back-edge bypass) — FIXED.** `consensus.go:91-99` escalates at
   `next > 1+HardCrossReviewCap` before `RunRound`; `driver.go:143-148` wires it from the policy.
   `TestBlockedConsensusRespectsTheHardCrossReviewCap` asserts escalation, names the cap, and
   asserts zero runner calls. I reverted the guard myself: the bypass reproduces (`action=reopened`,
   round 5 would open) — Q4.
2. **CRITICAL (stale skill manifest) — FIXED.** `node scripts/build-addon-manifest.js --check`
   exits 0 with `parley-deck: ok` and an aggregate matching `parley-addon.json`; `npm test`
   reports `tests 386, pass 386, fail 0` (round 1: 271 pass / 115 fail). `compatibility.json`
   2.8.0 tracks `package.json` 2.8.0.
3. **MAJOR (fast "0→1" claim) — FIXED BY WITHDRAWAL.** `TestFixupCapIsInclusive` no longer uses
   track labels — its cases are cap-labeled and vary published cycles. The changelog's "Not in
   this release" states `fast` Phase-8 is manual. `PolicyFor(Fast)` still rejects idea-level
   `auto_implement` (`track.go:143-145`), so the claimed route was indeed unreachable.
4. **MAJOR (strict rounds spending the fix-up budget) — FIXED IN UNIT, new MAJOR on the source.**
   `TestZeroFixRoundsDoNotSpendTheFixupBudget` passes and goes red on revert (Q4). The strict-gate
   closing guard (`round >= MaxFixupCycles`) is deliberately untouched, which round-1 accepted.
   But the chosen source for "published" is the author-controlled heading count — Finding 1 above.
5. **MAJOR (trajectory payload absent) — FIXED BY WITHDRAWAL**, which codex-1's own suggested fix
   explicitly permitted. Both changelogs now list the structured payload as not implemented; only
   the escalation-message word remains (Finding 3).
6. **NIT (stale comments) — FIXED.** `driver.go:119-123` and `track.go:115` now say only an
   absent/unknown track takes the legacy no-override path.

### Q2 — `HardCrossReviewCap` now also bounds `standard`'s back-edge at 2, not separately ratified

**Right in substance, with a ratification caveat — not, alone, release-blocking.**

**[PRIMARY — CONFIRMED]** The cap does bound every path that opens a Phase-2 round: the initially
scheduled budget (`advanceRound`, `driver.go:284-305`) is clamped to the cap, and the
consensus-BLOCK back-edge (`consensus.go:90-130`) is now hard-capped. Those are the only two
round-opening paths — `Rebuild` reconstructs from disk and opens nothing. Absent-track ideas keep
`HardCrossReviewCap=0` and the legacy `MaxRounds` bound, byte-for-byte backward compatible.

**[PRIMARY — CONFIRMED]** For `standard` this changes worst case from round-05 to round-03 —
exactly what its printed §4.0 cell ("capped at 2, then escalate/upgrade") always said. That is
enforcement of an **already-printed** cap, not a new policy cell, and it is the same text-vs-code
class this release exists to close (FINAL binding condition 1: "Printed caps bind only where
enforcement lives"). Bounding `deliberation`'s back-edge while leaving `standard`'s open under
`MaxRounds` would be incoherent.

**[RECALL]** FINAL's "no further §4.0 cell edits until the audit" binds text edits, not
enforcement alignment. Recommendation: the owner should explicitly confirm the `standard`
consequence (the implementer asked for exactly that), the changelogs should name it (Finding 2),
and the general per-track enforcement question remains with the follow-up audit idea as planned.

### Q3 — is counting `## Fix-up cycle N` from IMPLEMENTATION.md robust?

**No — fail-open in four ways.** See Finding 1 for the full enumeration and code citation:
missing/unreadable file restarts the budget at zero; fence-enclosed headings count; the count is a
heading count rather than an ordinal maximum, so duplicates over-count (safe) and renumbering or
deletion under-counts (**extends the budget**, unsafe); and the file is owned by the implementer
whose work the budget bounds, with no cross-check against the driver's own `.fixup-done` markers.
A careless edit can extend the budget with no in-loop signal; a malicious one can do it
deliberately. The one thing in its favor is calibration consistency (the 5 was measured with the
same naive count). The in-flight fix-up cycle 2 moves the count to driver-owned markers and
escalates on unknown count — the right shape, currently unverified and red (see the validity
event).

### Q4 — does the rewritten `TestFixupCapIsInclusive` actually go red on revert?

**[PRIMARY — CONFIRMED]** Verified by my own hands on the pre-mutation snapshot, not trusted:

- **Unit revert** (`cycle := publishedFixupCycles(d.cfg.IdeaDir) + 1` → `cycle := round`, operator
  kept `>`): red on exactly the three "escalates" subtests (`cap=5 published=5`, `cap=2
  published=2`, `cap=1 published=1` all ran `Fixup` instead of escalating) — matching the
  implementer's claim — and `TestZeroFixRoundsDoNotSpendTheFixupBudget` also went red (round 4,
  cap 2, zero published → escalated instead of running cycle 1). The three "allowed" cases stayed
  green because the review-round ordinal was 1 in all six cases: the rewritten test genuinely
  measures published cycles and cannot be gamed by round directories.
- **Operator revert** (`>` → `>=`, unit kept): red on exactly the three "at the cap is allowed"
  subtests. The test is sensitive to both halves of the fix under the new unit.
- **Hard-cap revert** (guard removed from `consensus.go`):
  `TestBlockedConsensusRespectsTheHardCrossReviewCap` red — `action=reopened`, no escalation,
  reproducing round-1 CRITICAL-1.
- **Restoration**: sha256 identical before/after (`77142dfe…` impl.go, `57c8e626…` consensus.go),
  package green after restore, backups deleted.

**[PRIMARY — CONFIRMED]** Caveat: the implementer's 15:07 rewrite of this same test is currently
red for a different reason (the AF2 marker collision in the validity event). My Q4 evidence
certifies the snapshot's test; the next round must re-run reversion checks on the final counter.

### Q5 — are the two CHANGELOGs accurate and not overstated?

**[PRIMARY — CONFIRMED]** Yes for the snapshot, with Finding 2 (omission) and Finding 3 (wording)
as the only residuals:

- The off-by-one narrative checks out arithmetically: old `cycle >= cap` with `cycle := round`
  publishes `standard` 1 cycle (printed cap 2) and `deliberation` 2 (silent default 3;
  `transport.go:34-42` default `cross_review_rounds` is 1). The inclusive bound now publishes
  1..N and escalates at N+1.
- The "Wrong unit" paragraph matches the reviewed code.
- "The cross-review cap binds every path that opens a round" — verified (Q2); both opening paths
  bounded, absent-track legacy preserved.
- "Why 5" — I re-ran the count over the 69 other ideas'
  `IMPLEMENTATION.md` files (`grep -c '^## Fix-up cycle'`): `0×17, 1×34, 2×7, 3×2, 4×3, 5×2, then
  9, 14, 15, 25`. Digit-for-digit the published distribution; nothing closed in the 6–8 band.
- "Not in this release": no protocol-packet code exists (`rg packet` finds only the pre-existing
  signoff `HandoffPacket` and the idea slug); the trajectory payload does not exist; the `fast`
  route is manual (`track.go:143-145` still rejects `fast + auto_implement`). All three true.
- Sibling: the two changed §4.0 cells are byte-identical across all three protocol copies, and
  "the CLI enforces both cells… the printed number and the enforced number are the same number" is
  now true for `deliberation` (5/5 inclusive published; 3-round hard cap) — with the Q2 note that
  `standard`'s printed 2 is now *also* enforced on the back-edge, which goes unmentioned.
- Conditional: if the in-flight marker-based counter lands, the "Wrong unit" paragraph ("derived
  from the `## Fix-up cycle N` records… in `IMPLEMENTATION.md`") becomes stale and must be
  rewritten in the same patch.

### Q6 — does the skill package now install?

**[PRIMARY — CONFIRMED]** Yes, run by me in the sibling working tree:

```text
$ node scripts/build-addon-manifest.js --check      → exit 0
parley-deck: ok (6 files, sha256:bee732321d0f…)     # matches parley-addon.json aggregate
$ npm test                                          → tests 386, pass 386, fail 0
```

(Round 1: `STALE parley-addon.json`, 271 pass / 115 fail.) The Python suite (54 tests) and all six
skill manifests are also green.

### Q7 — should anything stop this release?

**Yes, two things.**

1. **Finding 1 (MAJOR)** against the reviewed snapshot: the safety-critical count of this release
   came from author-controlled, unvalidated prose with four fail-open modes, so the ratified
   "5 inclusive published cycles, never reset" was not actually delivered.
2. **The candidate no longer exists in the tree.** A fix-up cycle 2 is visibly in flight
   (15:07 rewrite), the driver package is currently red, and whatever lands has not been reviewed
   by anyone. The release must wait until that work is finished, the suite and the skill package
   are green again, the changelogs are updated to match it, and a fresh review pass (with new
   reversion checks) covers the final counter.

Also to be recorded, not blocking: explicit owner confirmation of the `standard` back-edge
consequence (Q2), and the two wording fixes (Findings 2–3).

## Independent validation record

**[PRIMARY — CONFIRMED]** Run by me in this working tree (snapshot unless noted):

```text
go build ./... ; go vet ./...                                  exit 0
GOOS=linux|windows GOARCH=amd64 go build ./...                 exit 0
go test ./... -count=1                                         27/27 ok (incl. internal/runner)
focused driver/track/protocol tests (7 named tests, all subtests) PASS
reversion: unit revert → 3 subtests + TestZeroFixRounds… red   confirmed, restored byte-exact
reversion: operator revert → 3 "allowed" subtests red          confirmed, restored byte-exact
reversion: hard-cap guard removal → TestBlocked… red           confirmed, restored byte-exact
69-idea fix-up-cycle distribution                              reproduced digit-for-digit
node scripts/build-addon-manifest.js --check                   exit 0 (all six skills ok)
npm test (sibling)                                             386 pass / 0 fail
git diff --check (both repos)                                  exit 0
gofmt -l internal/driver internal/track                        clean
```

**[PRIMARY — CONFIRMED]** Current (post-15:07) tree: builds and `gofmt` clean; `go test
./internal/driver -run TestFixupCapIsInclusive` fails 5 subtests for the AF2-collision reason
above. Not certified.
