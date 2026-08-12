---
agent: hermes-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 5
date: 2026-08-12
---
verdict: CLEAN

# Phase 6 review — round 5

## Summary

Fix-up cycle 4 fixed three defects @codex-1 proved in round 4 (reservation
ordering, AF2 boundary at equality, cursor-error path) and declined a fourth on
scope, correcting the release claims instead. I verified all three fixes by
isolated reversion checks — each fix reverted, the test confirmed red, restored,
re-verified green. The scope call is the right release decision: the stronger
property @codex-1's CRITICAL asks for is a different system from the two numbers
this idea ratified, building it silently would repeat the error the FINAL warns
about, and the corrected CHANGELOGs state the limits plainly.

## Q1 — Are the three fixed defects actually fixed?

Yes. All three, verified by isolated reversion checks in /tmp/hermes1-r05-iso.

### Fix 1 — Reservation ordering (impl.go:317-325)

[PRIMARY — CONFIRMED, source inspection + isolated reversion check]

The cycle is now reserved BEFORE the code-writing call. The code at
impl.go:322-325 sets `c.FixupCyclesPublished = cycle` and calls `saveCursor`
BEFORE `d.cfg.Impl.Fixup(ctx, cycle)` at line 326. A Fixup that errors, or a
crash in the window between reservation and Fixup, cannot get the cycle back.

Reversion check: I moved the reservation back to after Fixup (the old ordering)
in the isolated copy. My probe `TestR05ReservationBeforeFixupOnError` wraps
fakeImpl with a Fixup that always returns an error, then checks the cursor after
Advance. With the old ordering: the cursor file was never written (Fixup errored
before saveCursor), so `LoadCursor` returned "no such file" — the cycle was
never reserved. RED. With the fixed ordering: `FixupCyclesPublished=1` was
persisted even though Fixup errored. GREEN. The reservation is before the
code-writing call.

### Fix 2 — AF2 boundary at equality (impl.go:147)

[PRIMARY — CONFIRMED, source inspection + isolated reversion check]

The AF2 gate at impl.go:147 now uses strictly greater: `spent >
d.cfg.MaxFixupCycles`. At equality (spent == cap), AF2 is allowed to finish the
Nth allowed cycle — the one whose budget is already spent — instead of
escalating. Starting the next cycle (N+1) is still refused by the ordinary
branch at impl.go:308 (`cycle > d.cfg.MaxFixupCycles`).

Reversion check: I changed `>` back to `>=` in the isolated copy. My probe
`TestR05AF2AtEqualityCompletesLastCycle` sets up 5 published cycles (markers in
rounds 1-5, cap=5), a marker in round 5 (triggering AF2), and a cursor with
FixupCyclesPublished=5. With `>=`: AF2 escalated with "fix-up budget exceeded:
5 cycle(s) recorded against MaxFixupCycles=5" — stranding the legitimate crash
recovery of the 5th (last allowed) cycle. RED. With `>`: AF2 completed
(action=fixup, calls=[open-review]) — the 5th cycle's transition finished.
GREEN.

### Fix 3 — Cursor read error (driver.go:235-248)

[PRIMARY — CONFIRMED, source inspection + isolated reversion check]

The carry-forward in Advance now distinguishes three cases:
- `err == nil`: carry forward monotonically (existing behavior, preserved).
- `os.IsNotExist(err)`: a fresh run, markers are the only source (correct).
- default (exists but unreadable): **escalates** with "cannot read the driver
  cursor...refusing to act on an unknown fix-up budget".

The old code was `if err == nil && prev.FixupCyclesPublished > c.FixupCyclesPublished`
— the `&&` meant every LoadCursor error (corrupt, unreadable, permission denied)
was silently swallowed, and the cursor value was treated as zero.

Reversion check: I reverted the carry-forward to the old `err == nil &&` form in
the isolated copy. My probe `TestR05CorruptCursorAfterErroredFixupRestartsBudget`
sets up the exact scenario: zero markers (a fix-up errored before writing one),
a corrupted cursor file (invalid JSON), cap=1. With the old code: LoadCursor
failed silently, fell back to markers (0), publishedFixupCycles returned 0, the
budget restarted at zero, and a fix-up ran (action=fixup). RED — the exact
fail-open @codex-1 proved. With the fixed code: Advance escalated with "cannot
read the driver cursor...refusing to act on an unknown fix-up budget: invalid
character 't' looking for beginning of object key string". GREEN.

### All reversion checks: isolated copy verified clean

After all three reversion checks, I compared the SHA-256 hashes of impl.go,
driver.go, cursor.go, and consensus.go between the isolated copy and the working
tree. All four matched exactly. The isolated copy was deleted.

### Full test suite

[PRIMARY — CONFIRMED, terminal output this session]

```
go build ./...                     BUILD_OK
go vet ./...                       VET_OK
go test ./... -count=1             all 28 packages PASS
node scripts/build-addon-manifest.js --check   all 6 skills ok
npm test (parley-deck-skill)       386 pass, 0 fail
```

Focused tests, all PASS:

```
TestFixupCapIsInclusive                           PASS (6 subtests)
TestZeroFixRoundsDoNotSpendTheFixupBudget         PASS
TestFixupBudgetIsTamperFailSafe                   PASS (2 subtests)
TestOnlyExactRoundDirsCountAsPublishedCycles      PASS
TestBlockedConsensusRespectsTheHardCrossReviewCap PASS
TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops PASS
TestDeliberationClampsCrossReviewRoundsToThree    PASS
TestPhaseReviewMaxFixupCyclesEscalates            PASS
```

## Q2 — The scope call

"Correct the claims and name a follow-up" is the correct release decision. The
stronger property must NOT ship before this release.

@codex-1's CRITICAL requires the ceiling to be idea-scoped, serialized across
concurrent runs, and anchored outside the participant-writable repository. Its
probes are correct: two runs of one idea keep independent cursors and independent
locks, a deleted run directory loses the counter, and a repository rollback
restores a lower one. I confirmed all of these in my round-4 review and did not
dispute them.

But the disposition is a release decision, and on the release merits:

1. The idea ratified two numbers and their enforcement — "cap 5 cycles" and
   "capped at 3 after round 1, then escalate" — plus the inclusive boundary
   fix. It did not ratify an idea-scoped, serialized, out-of-repo trust anchor.
   @codex-1's CRITICAL asks for a system that is categorically different from
   what this idea's FINAL defines.

2. An out-of-repo anchor breaks the deck's file-based portability. A fresh clone
   loses the count entirely — the exact fail-open @codex-1 objects to, just
   moved to a different layer. The deck's design assumes a clone is a complete
   workspace; an external anchor silently violates that.

3. Building the stronger system silently under this mandate is the error the
   FINAL itself warns about: "a number that is a safety boundary must not be
   authored by the party it constrains" — but the scope of what this idea
   authorized is also a boundary, and silently expanding it is the same class
   of mistake.

4. The corrected CHANGELOGs state the limits plainly. The CLI CHANGELOG
   (lines 32-37) says: "The budget is robust against accidental loss, a stale
   or deleted single record, and an errored or crashed fix-up. It is not a
   security boundary: a participant with workspace write, a deleted run
   directory, a repository rollback, or two concurrent runs of the same idea
   can still reduce or duplicate the count." The skill CHANGELOG (lines 17-18)
   says the same and points to the CLI for exact limits. A user reading either
   knows exactly what they have and what they do not.

Nothing breaks if the stronger property does not ship in this release. The
release ships correct enforcement of the two ratified numbers, with limits
honestly stated. The stronger property is a named follow-up
(`fixup-budget-trust-anchor`), not a prerequisite for the two numbers this idea
was created to enforce. Requiring it to ship would block a release that fixes
real, verified, agreed-upon defects (the three from cycle 4, plus the three from
cycle 3, plus the three from cycle 2, plus the two from cycle 1) behind a system
nobody ratified.

## Q3 — Do the corrected CHANGELOG limits accurately describe what ships?

Yes. I read both CHANGELOGs line by line against the actual code.

### CLI CHANGELOG.md

- Lines 9-12 (off-by-one): accurate. The guard was `>=`, now `>`, inclusive.
- Lines 14-16 (wrong unit): accurate. The unit was the review-round ordinal, now
  completed fix-up cycles. The word "completed" on line 16 has a minor tension
  with "reserved before the code-writing call" on line 26 — a cycle is now spent
  when it runs, not when it completes — but line 26 immediately clarifies the
  actual semantics, and the "What this is not" section (lines 32-37) makes the
  limits fully clear. Not wrong, not misleading.
- Lines 18-22 (wrong source history): accurate. The four fail-open ways of
  counting IMPLEMENTATION.md headings are the ones @codex-1 proved in round 2.
- Lines 24-30 (current design): accurate. "Maximum of two driver-authored
  records", "reserved before the code-writing call", "Losing one record does not
  lower the count", "forging one can only raise it", "an unreadable cursor
  escalates instead of counting as zero", "the crash-recovery path consults the
  cap before opening another round — while still being allowed to finish the
  last allowed cycle." All of these match the code I verified.
- Lines 32-37 (what this is not): accurate and correctly scoped. "Not a security
  boundary: a participant with workspace write, a deleted run directory, a
  repository rollback, or two concurrent runs of the same idea can still reduce
  or duplicate the count." This is the exact set of residual limits I
  independently confirmed in my round-4 review.
- Lines 44-45 (the two cells): accurate. "5 inclusive published cycles" and
  "capped at 3 after round 1, then escalate."
- Lines 53-58 (cross-review cap binds every path): accurate. The
  HardCrossReviewCap is checked on the BLOCK back-edge at consensus.go:92-96.
- Lines 74-86 (not in this release): accurate. The packet is not started, the
  trajectory payload is not implemented, and `fast` Phase 8 is not exercised.

### Skill CHANGELOG.md

- Lines 11-14 (the two cells): accurate.
- Lines 16-18 (enforcement limits): accurate. "robust against accidental loss
  of its own records but is not a security boundary against a participant that
  edits the workspace — see parley-deck-cli 1.44.0 for the exact limits."
- Lines 19-20 (cross-review ceiling binds BLOCK back-edge): accurate.
- Lines 21-22 (trajectory payload not implemented): accurate.
- Lines 24-26 (ratification + "the printed number and the enforced number are
  the same number"): This is accurate under the corrected limits. The printed
  number is the §4.0 table value (5, 3). The enforced number is the
  max-of-two-sources count, which matches the printed cap as long as neither
  source is tampered with. The CHANGELOG does not claim tamper-proofness — the
  preceding lines (17-18) explicitly disclaim it. "The same number" means "the
  text and the code agree about what the cap is," not "the cap is
  cryptographically enforced." A reader who gets to line 26 has already read
  the disclaimer on lines 17-18.

Nothing is still wrong.

## Q4 — Anything NEW in the round-4 diff that was not there before

No new defects. The round-4 diff contains exactly the three fixes (reservation
ordering, AF2 boundary, cursor-error path) and the two CHANGELOG corrections.
No new code paths, no new functions, no new tests beyond those covering the
three fixes. The `publishedFixupCycles` function, `markedFixupCycles`,
`isRoundDirName`, and the carry-forward in Advance are all from cycle 3 and
unchanged. The AF2 gate comment was updated to explain the `>` vs `>=` decision.
The stale comment at impl.go:297-298 (naming IMPLEMENTATION.md as the source)
is the same NIT I flagged in round 4 and is still present — it is a comment,
not behaviour.

## Q5 — Should this release ship?

Yes.

The three defects @codex-1 proved in round 4 are fixed and verified by isolated
reversion checks. The scope call — correct the claims, name the follow-up — is
the right release decision: the stronger property is a different system from
what this idea ratified, building it silently would repeat the mistake the FINAL
warns about, and the CHANGELOGs now state the limits honestly. Both CHANGELOGs
accurately describe what ships, neither overstated nor understated.

The one remaining issue is the NIT from round 4: the code comment at
impl.go:297-298 still names IMPLEMENTATION.md as the source of the cycle count,
when the actual source is the max of the cursor counter and .fixup-done markers.
It does not affect behaviour and should not block a release.

## Findings index

| Severity | Finding |
| --- | --- |
| NIT | The code comment at impl.go:297-298 still names IMPLEMENTATION.md as the source of the cycle count; the actual source is the max of the cursor counter and .fixup-done markers. The comment's claim about the unit (published cycles, not round ordinal) is still correct; only the named source is stale. Carried from round 4. |

## Responses to other reviewers

### @codex-1

Your round-4 CRITICAL has three parts. The three concrete defects you proved
(reservation taken too late, AF2 rejecting equality, cursor read error silently
ignored) are all fixed — I verified each by isolated reversion check. Your
fourth requirement (idea-scoped, serialized, out-of-repo trust anchor) is the
scope question, and on the release merits the implementer's disposition is
correct: correct the claims, name the follow-up. The stronger property is a
different system from the two numbers this idea ratified. Building it silently
would be the exact error the FINAL warns about — a party silently expanding the
scope of what was ratified. The CHANGELOGs now say plainly what this is and what
it is not, so a user knows the limits without reading the source.

On the AF2 boundary: your MAJOR was right, and the fix adopts your exact
threshold. `>` at AF2, `>` at the ordinary branch. Equality finishes the Nth
cycle; N+1 is refused.
