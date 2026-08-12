---
agent: hermes-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 4
date: 2026-08-12
---
verdict: CLEAN

# Phase 6 review — round 4

## Summary

Fix-up cycle 3 closes the fail-open defect that cycles 1 and 2 relocated. The
budget is now the maximum of two independent driver-authored sources — a
monotonic counter in the run cursor (outside the idea directory) and the
`.fixup-done` markers — and a cycle is spent when it runs, not when it passes.
The AF2 crash-recovery branch now consults the cap before opening another
round, and only exact `round-<digits>` directories count.

I attacked the budget from every angle the round-04 brief asked for: two runs
sharing an idea, a deleted run directory, a deleted review tree, the crash
window between Fixup and cursor persist, the crash window between cursor
persist and RunChecks, and the simultaneous deletion of both sources. I ran
three reversion checks in an isolated copy of the module. The design holds on
every path except the theoretical floor where an attacker deletes BOTH
independent sources at once, which is a property of any on-disk counter
without an external authority and was already the case under every prior
design in this idea.

One NIT remains: a stale code comment at impl.go:293-298 still names
IMPLEMENTATION.md as the source. It does not affect behaviour.

## Q1 — Attack the budget again

I traced every path the brief asked about, built adversarial probes in an
isolated copy at /tmp/hermes1-r04-iso2, and ran them. Every probe was in the
isolated copy only — the working tree was never edited for this review.

### Two runs sharing an idea

[PRIMARY — CONFIRMED, isolated probe TestTwoRunsSharingIdea] Run A spends 2
fix-up cycles (markers in rounds 1-2, cursor with FixupCyclesPublished=2).
Run B starts with the same idea directory but its own empty run directory.
Run B's cursor has FixupCyclesPublished=0, but the markers hold 2, so
publishedFixupCycles = max(0, 2) = 2. Cycle 3 > cap 2 → escalate. Run B
cannot re-spend the budget.

This is the key property of the max-of-two design: a second run cannot reset
the count by starting fresh, because the markers — written by the driver, in
the idea's review directory — survive across runs.

### Deleted run directory

[PRIMARY — CONFIRMED, isolated probe TestDeletedRunDirectory] Deleting the
entire run directory (and thus driver.json) does not buy a cycle. The cursor
is gone, but the markers remain. publishedFixupCycles = max(0, markers) =
markers. With 2 markers and cap=2, cycle 3 > 2 → escalate.

### Rolled-back git checkout

A `git checkout` that removes the `.fixup-done` markers (because they are
uncommitted or on a different branch) would lower the marker count. But the
cursor lives in the run directory (RunDir/driver.json), which is outside the
idea directory and typically outside the git tree entirely (under
`~/.parley/runs/` or a project-level runs/ dir). A git checkout of the idea
directory does not touch the run directory.

If the run directory IS inside the git tree and the checkout removes it too,
then both sources are lost — the theoretical floor discussed below.

The more realistic scenario: a reviewer rolls back the idea directory to
remove markers from a prior round. The cursor still holds the count via the
carry-forward in Advance (driver.go:235-237). The count cannot drop.

### The window between Fixup returning and the cursor being persisted

[PRIMARY — CONFIRMED, isolated probe TestCrashBetweenFixupAndCursorSave] The
order at impl.go:313-321 is: Fixup runs → c.FixupCyclesPublished = cycle →
saveCursor → RunChecks → write marker → archive → open next round.

If a crash happens between Fixup (313) and saveCursor (321): neither the
cursor nor the marker records the spend. On re-entry,
publishedFixupCycles returns the OLD count. The same cycle number is
recomputed, and Fixup runs again.

This is CORRECT crash recovery, not budget extension: the cycle was never
durably spent. No round was opened, no marker was written, no cursor was
persisted. Re-running the same fix-up is the right thing to do. The fix-up's
own idempotency is a separate concern — the budget correctly does not count a
cycle that was never persisted.

[PRIMARY — CONFIRMED, isolated probe TestCrashAfterCursorSaveBeforeChecks]
If a crash happens after saveCursor (321) but before RunChecks (326): the
cursor HAS FixupCyclesPublished=cycle, but no marker exists and checks never
ran. On re-entry, the carry-forward in Advance restores
FixupCyclesPublished from the persisted cursor. publishedFixupCycles =
max(cursor, markers) = cursor. The next cycle (published+1) escalates if it
exceeds the cap.

This is the "spent when it runs" rule working correctly: the cycle was spent
at Fixup time, persisted before the check gate, and an escalation below
cannot un-spend it. A fix-up that crashed after running but before checks
still consumed its cycle — which is right, because the fix-up DID run and
potentially mutated code.

### Concurrent runs

Two concurrent driver processes sharing the same run directory would race on
cursor writes. The cursor save is atomic (tmp + rename, cursor.go:63-79), so
a mid-write crash cannot corrupt it. But two processes could both read the
old cursor, both run Fixup, and both write — the second write wins. This is
the same concurrency hazard that exists for every cursor write in the driver,
not specific to the fix-up budget. The driver assumes a single process per
run (the lock is at the run level, not shown in this diff). I did not find a
new concurrency vulnerability specific to the budget.

### The theoretical floor: both sources deleted

[PRIMARY — CONFIRMED, isolated probe TestBothMarkersAndCursorDeleted] If an
attacker deletes ALL `.fixup-done` markers AND the run directory (or
driver.json), publishedFixupCycles returns (0, nil) and the budget restarts
at zero. This is the floor of the max-of-two approach: two independent
sources are harder to tamper with than one, but both are on-disk state
accessible to a process with workspace write.

This is not a regression from any prior design in this idea. The prose source
(cycle 1) was editable by the implementer. The marker-only source (cycle 2)
was editable by any participant. The max-of-two (cycle 3) requires deleting
state in two separate directories — the idea directory and the run directory
— which is a strictly higher bar. No on-disk counter can be tamper-proof
without an external authority (signed ledger, OS-enforced permissions), and
the FINAL does not require one.

### Is there ANY remaining path that lowers the count or opens a fix-up past the cap?

No path that I could find, other than the simultaneous deletion of both
sources (the theoretical floor). Every individual-source deletion is caught
by the other source. Every forge only raises the count. Every crash window
either correctly re-runs an un-persisted cycle or correctly counts a spent
one. The AF2 branch is now gated. The exact `round-<digits>` filter prevents
oddly-named directories from inflating or confusing the count.

## Q2 — "Spent when it runs, not when it passes"

I am persuaded. I change my position from round 3.

In round 3, I argued that a fix-up whose checks failed should not spend the
budget because it "was not published — it escalated." @codex-1 argued the
opposite: a fix-up that ran and broke the build is exactly the churn the cap
exists to interrupt, and handing the budget back lets a failing fix-up loop
forever against a ceiling that never depletes.

On reflection, @codex-1 is right on the merits, and the implementation is
correct. The key insight is that "published" was my word, not the FINAL's.
The FINAL says "5 inclusive published cycles" but the intent — stated in the
same paragraph — is a "finite fix-up-attempt breaker" that interrupts
pathological churn. A fix-up that runs, mutates code, and breaks the build IS
churn. If it does not spend the budget, a fix-up that always fails can run
indefinitely: crash → re-enter → Fixup → checks fail → escalate → human
re-enters → Fixup → checks fail → escalate → ... The cap never depletes
because no cycle is ever "published" under my old definition.

The implementation at impl.go:316-323 persists `c.FixupCyclesPublished =
cycle` immediately after Fixup returns, BEFORE RunChecks. If RunChecks fails
and escalates, the count is already persisted. On re-entry, the carry-forward
in Advance (driver.go:235-237) restores the persisted value. The cycle
cannot be un-spent.

This is the right design. I was wrong in round 3.

## Q3 — Is the AF2 gate correct, or does it strand a legitimate crash recovery?

The AF2 gate is correct and does not strand legitimate crash recovery.

The AF2 branch (impl.go:140-162) fires when a `.fixup-done` marker exists in
the current review round — meaning Fixup + RunChecks already succeeded for
this round, but the driver crashed before opening the next one. The branch
archives the consensus and opens round+1 without re-running Fixup.

The budget gate added at impl.go:145-149 checks `spent >= MaxFixupCycles`
BEFORE the archive-and-open. If the budget is exhausted, it escalates instead
of opening another round.

This is correct because:

1. The marker means the fix-up already ran and checks passed. The cycle was
   already spent (the cursor was persisted at impl.go:320-321 before the
   marker was written at 331). Opening another round would start a NEW
   review cycle, which needs a NEW fix-up budget slot. If the budget is
   exhausted, there is no slot.

2. The crash recovery that AF2 provides — "don't re-run Fixup" — is
   preserved. The gate does not remove the marker or re-run Fixup. It simply
   refuses to open another round when the budget is spent, and escalates
   instead. The human can grant a recorded finite extension, which is exactly
   what the FINAL says should happen at a budget boundary.

3. A legitimate crash recovery that should be allowed to finish is one where
   the fix-up succeeded but the round transition didn't complete. That
   transition is: archive consensus → open next round. If the budget is
   exhausted, there IS no "finish" — the next round would be past the cap.
   Escalating is the correct finish.

[PRIMARY — CONFIRMED, isolated probe TestCrashAfterFixupBeforeCursorPersist]
I verified that after a successful fix-up (cursor persisted, marker written,
round opened), a re-entry with checks failing correctly escalates with
"fix-up budget exhausted" and does NOT re-run fixup.

## Q4 — Tests and reversion checks

All tests ran in the working tree [PRIMARY — CONFIRMED, terminal output this
session]:

```
go build ./...                                    exit 0
go vet ./...                                      exit 0
go test ./... -count=1                            all 28 packages PASS
node scripts/build-addon-manifest.js --check      all 6 skills ok
npm test (parley-deck-skill)                      386 pass, 0 fail
TestEmbeddedDefaultMatchesLiveDeck                PASS
```

Focused tests, all PASS:

```
TestFixupCapIsInclusive                           PASS (all 6 subtests)
TestZeroFixRoundsDoNotSpendTheFixupBudget         PASS
TestFixupBudgetIsTamperFailSafe                   PASS (both subtests)
TestOnlyExactRoundDirsCountAsPublishedCycles      PASS
TestBlockedConsensusRespectsTheHardCrossReviewCap PASS
TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops PASS
TestDeliberationClampsCrossReviewRoundsToThree    PASS
```

### Reversion checks — in an ISOLATED COPY, not the working tree

Per the binding rule and the round-2 process failure record, I copied the Go
module to /tmp/hermes1-r04-iso2 (rsync, excluding .git and .gomodcache) and
performed all reversion checks there. The working tree was never edited for
these checks. After all checks, the isolated copy's impl.go, driver.go, and
cursor.go SHA-256 hashes exactly matched the working tree's, confirming no
contamination.

**Reversion 1: monotonic carry-forward removed.**
Replaced the carry-forward in Advance (driver.go:235-237) with a no-op.
[PRIMARY — CONFIRMED] TestFixupBudgetIsTamperFailSafe subtest "deleting the
markers does not lower the count when the cursor holds it" went RED:
"wiping markers bought a cycle past the cap: action=fixup err=<nil>". The
subtest "a forged marker in the current round cannot drive AF2 past the cap"
still passed (it does not depend on the carry-forward). Restored, re-verified
green.

**Reversion 2: AF2 budget gate removed.**
Removed the budget check at impl.go:145-149 (the `publishedFixupCycles` call
and cap comparison on the AF2 branch).
[PRIMARY — CONFIRMED] TestFixupBudgetIsTamperFailSafe subtest "a forged
marker in the current round cannot drive AF2 past the cap" went RED:
"AF2 opened another round past the cap: action=fixup err=<nil>
calls=[open-review]". The marker-deletion subtest still passed (it does not
depend on the AF2 gate). Restored, re-verified green.

**Reversion 3: exact round-dir matching reverted to HasPrefix.**
Replaced `isRoundDirName(e.Name())` with `strings.HasPrefix(e.Name(),
"round-")` at impl.go:549.
[PRIMARY — CONFIRMED] TestOnlyExactRoundDirsCountAsPublishedCycles went RED:
"counted 5 published cycles, want 2 — only round-01 and round-02 are real
rounds". The HasPrefix test counts round-backup, round-x, rounds-03, and
round- as valid rounds. Restored, re-verified green.

### Adversarial probes — also in the isolated copy

I wrote six additional adversarial probes in the isolated copy to test the
Q1 scenarios. All passed:

```
TestTwoRunsSharingIdea                           PASS
TestDeletedRunDirectory                          PASS
TestBothMarkersAndCursorDeleted                  PASS (documents the floor)
TestCrashAfterFixupBeforeCursorPersist           PASS
TestCrashBetweenFixupAndCursorSave               PASS
TestCrashAfterCursorSaveBeforeChecks             PASS
```

The isolated copy and all probe files were deleted after verification. The
working tree's `git diff --stat` is unchanged.

## Q5 — Are both CHANGELOGs and IMPLEMENTATION.md now accurate?

### CLI CHANGELOG.md

[PRIMARY — CONFIRMED] The CLI CHANGELOG has been fully rewritten for cycle 3.
Lines 14-30 now accurately describe:

- The wrong unit and wrong source history (round ordinal → IMPLEMENTATION.md
  headings → .fixup-done markers alone → max-of-two).
- The current design: "the maximum of two independent driver-authored
  sources: a monotonic counter in the run cursor, carried forward across
  rebuilds and written inside the fix-up transaction, and the .fixup-done
  markers."
- "Deleting either cannot lower the count; forging either can only raise it,
  which escalates sooner."
- "A fix-up spends its cycle the moment it runs, not when it passes."
- "An unreadable count escalates instead of restarting the budget at zero."
- "The crash-recovery path now consults the cap before opening another round."

The "Why 5" paragraph (CHANGELOG.md:53-57) references `## Fix-up cycle`
headings across all 69 ideas — this is a historical analysis used to
calibrate the cap value, not a claim about the enforcement source. It remains
accurate as motivation.

The "Not in this release" section (CHANGELOG.md:67-79) accurately states the
packet is not started, the trajectory payload is not implemented, and `fast`
Phase 8 is not exercised by the driver.

### Skill CHANGELOG.md

[PRIMARY — CONFIRMED] The skill CHANGELOG accurately describes the protocol
text changes (the two budget cells), the cross-review ceiling binding the
BLOCK back-edge, and the trajectory payload not being implemented. The claim
"the printed number and the enforced number are the same number" (line 24)
is accurate for the max-of-two design — the enforced number is the max of the
cursor and markers, which matches the printed cap as long as neither source
is tampered with. The CHANGELOG does not claim tamper-proofness.

### IMPLEMENTATION.md

[PRIMARY — CONFIRMED] The fix-up cycle 3 section accurately describes:

- The max-of-two-sources design and the tamper table.
- The AF2 bypass fix (my round-3 CRITICAL).
- The "spent when it runs" rule and the persist-before-check-gate ordering.
- The exact round-dir matching fix (my round-3 MAJOR).
- The CHANGELOG rewrite (my round-3 MAJOR).
- The carry-forward fix found by the first tamper test failing.
- The reversion checks run in an isolated copy.

### [NIT] Stale code comment at impl.go:293-298

The comment at impl.go:293-298 still says:

> The cycle number is derived from the fix-up cycles actually PUBLISHED in
> IMPLEMENTATION.md, not from the review-round ordinal.

This is stale. The cycle number is now derived from `publishedFixupCycles`,
which returns the maximum of the cursor counter and the `.fixup-done`
markers — NOT from IMPLEMENTATION.md headings. This was flagged in my round-3
review (MAJOR) and @codex-1's round-3 review, and the IMPLEMENTATION.md
fix-up cycle 3 section mentions the CHANGELOG rewrite but does not mention
fixing this comment. The comment was not touched in cycle 3.

It is a NIT, not a MAJOR: the code below the comment is correct, and the
comment's claim about the unit (published cycles, not round ordinal) is
still right. Only the source named in the comment is wrong. The
`publishedFixupCycles` function's own doc comment (impl.go:507-525) is
correct and describes the max-of-two design accurately.

## Q6 — Anything that should stop this release?

No.

The fail-open defect that spanned three fix-up cycles is closed. The budget
is the maximum of two independent driver-authored sources, making every
single-source tamper fail-safe. The AF2 bypass is gated. The "spent when it
runs" rule prevents a failing fix-up from looping forever. The exact
round-dir matching prevents oddly-named directories from confusing the count.
Both CHANGELOGs and IMPLEMENTATION.md accurately describe what ships.

The one remaining issue is a NIT (stale code comment at impl.go:293-298)
that does not affect behaviour and should not block a release.

The theoretical floor — deleting both independent sources simultaneously
resets the budget — is a property of any on-disk counter without an external
authority, was present in every prior design in this idea, and is not a
regression. The max-of-two design raises the bar from "edit one file in the
idea directory" to "delete state in two separate directories."

I am CLEAN on this release.

## Findings index

| Severity | Finding |
| --- | --- |
| NIT | The code comment at impl.go:293-298 still names IMPLEMENTATION.md as the source of the cycle count; the actual source is the max of the cursor counter and .fixup-done markers. The comment's claim about the unit (published cycles, not round ordinal) is still correct; only the named source is stale. |

## Responses to other reviewers

### @codex-1

Your round-3 CRITICAL (marker forgeability + AF2 bypass) and MAJOR (failed
fix-up handing the budget back) are both fixed and verified. On Q2, you were
right and I was wrong — I say so plainly in Q2 above. The "spent when it
runs" rule is the correct design, and the implementation persists the count
before the check gate so an escalation cannot un-spend it. The stale comment
NIT at impl.go:293-298 is the same one you flagged at impl.go:284-285 in
round 3; it was not fixed in cycle 3 but does not affect behaviour.
