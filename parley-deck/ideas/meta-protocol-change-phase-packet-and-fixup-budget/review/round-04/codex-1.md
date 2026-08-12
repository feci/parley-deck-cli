---
agent: codex-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 4
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 4

## Summary

Fix-up cycle 3 closes the three narrow regressions it names: the cursor value survives `Rebuild`
within one intact run, AF2 no longer opens a round when the marker count is already past the cap,
and only exact `round-<digits>` directory names count. The new tests genuinely detect reverting
those changes.

The safety boundary is still not monotonic at the scope it constrains. The counter is persisted
**after** `Fixup` returns, belongs to one run rather than the idea, treats a missing/corrupt cursor
as non-fatal, and shares the same repository and participant-writable failure domain as the
markers. I reproduced seven ways to execute more fix-ups than a cap of one permits: a `Fixup`
error, a save failure after `Fixup` returns, a second run for the same idea, two concurrent runs,
deleting the run directory, corrupting the cursor, and lowering its field to zero. The AF2 repair
also introduces the opposite failure: after the fifth allowed fix-up has completed and written its
marker, a crash before opening review round 6 is now escalated instead of recovered.

## Findings

### [CRITICAL] The high-water mark is run-local, participant-writable, and persisted too late

**[PRIMARY] Source inspection.** `internal/driver/impl.go:313-323` calls `Fixup` first, returns
immediately on a `Fixup` error, and only after a successful return assigns and saves
`FixupCyclesPublished`. Therefore neither an errored fix-up nor the window between `Fixup`
returning and `saveCursor` succeeding has a durable reservation. This is the exact window my
round-3 counter-proposal required the implementation to close.

`internal/driver/driver.go:180` locates the cursor at `<RunDir>/driver.json`.
`internal/driver/loop.go:20-31` locks `<RunDir>/driver.lock`, not the idea. Two runs for one idea
therefore have independent counters and independent locks. `internal/runstate/runstate.go:303-319`
also accepts an exact run ID, so multiple extant runs for one idea are individually continuable;
there is no idea-level uniqueness or lock check here.

**[PRIMARY] Isolated adversarial probes.** I added same-package tests only under
`/tmp/codex-r04-cli.O98uLm`, then ran:

```text
$ go test ./internal/driver -count=1 \
    -run '^TestRound4(Crash|FixupError|SecondRun|ConcurrentRuns|DeletedRun|CorruptCursor)' -v
PASS TestRound4CrashAfterFixupBeforeCursorSaveReopensBudget
PASS TestRound4FixupErrorDoesNotSpendBudget
PASS TestRound4SecondRunForSameIdeaReopensBudget
PASS TestRound4ConcurrentRunsCollapseTwoFixupsIntoOneCount
PASS TestRound4DeletedRunDirectoryReopensBudget
PASS TestRound4CorruptCursorReopensBudget
PASS

$ go test ./internal/driver -count=1 \
    -run '^TestRound4(LoweredCursorField|CorruptCursor)' -v
PASS TestRound4CorruptCursorReopensBudget
PASS TestRound4LoweredCursorFieldReopensBudget
PASS
```

These are attack-reproduction tests: `PASS` means the current code repeated `Fixup`. In the
concurrent case, two drivers with different run directories both entered `Fixup`, both saved the
value `1`, and their two executions collapsed to one shared marker. Both run-local locks could be
held simultaneously. In each sequential reset case, cap 1 permitted a second execution still
labelled cycle 1.

**[PRIMARY — WRONG]** The source comment at `internal/driver/impl.go:507-525`, CLI
`CHANGELOG.md:24-30`, and `IMPLEMENTATION.md:249-266` claim every tamper direction is fail-safe and
that an unreadable count escalates. `internal/driver/driver.go:230-237` silently ignores every
`LoadCursor` error. A missing, corrupt, unreadable, or valid-but-lowered cursor therefore falls
back to markers; after an errored fix-up or failed post-fix-up check there is deliberately no
marker, so the spent count becomes zero.

The two files are not independent trust or rollback domains. Both are below the repository root:
the cursor under `parley-deck/runs/`, the marker under `parley-deck/ideas/`. A read-only check in
the shared tree returned:

```text
$ git check-ignore -q parley-deck/runs/example/driver.json; echo $?
1
$ git check-ignore -q parley-deck/ideas/example/review/round-01/.fixup-done; echo $?
1
```

Neither path is ignored. I did not perform a checkout because this review forbids git writes.
Inference from the located paths and load logic: a repository rollback that restores both paths to
an earlier revision restores a lower count; deleting the run directory loses every failed-attempt
count; and a participant with workspace write can lower or remove both records just as it could the
cycle-2 marker. “Outside the idea directory” is not access separation.

Concrete fix: reserve the next attempt durably **before** invoking code-writing `Fixup` (a crash
before confirmed launch may conservatively spend the reservation), make the high-water mark
idea-scoped across all runs, and serialize at the idea boundary. A missing, corrupt, rolled-back,
or contradictory safety record must escalate. The durable high-water mark must be outside the
participant-writable repository, or reconciliation must detect loss and fail closed; two ordinary
workspace files are not that boundary.

### [MAJOR] AF2 strands the legitimate crash recovery of the last allowed cycle

**[PRIMARY] Source inspection.** AF2 is documented at `internal/driver/impl.go:136-139` as
finishing an already-completed `Fixup+RunChecks` after a crash before the next review round opens.
The new guard at lines 145-149 rejects `spent >= MaxFixupCycles`. At the inclusive boundary, the
fifth allowed fix-up has already set the cursor to 5 and written marker 5. AF2 is completing cycle
5, not starting cycle 6, so equality must be recoverable.

**[PRIMARY] Isolated boundary probe:**

```text
$ go test ./internal/driver -count=1 \
    -run '^TestRound4AF2CompletesTheNthAllowedCycleAtCap$' -v
=== RUN   TestRound4AF2CompletesTheNthAllowedCycleAtCap
legitimate AF2 recovery at cap was stranded:
action=escalated
err=fix-up budget exhausted at 5 cycle(s) (MaxFixupCycles=5)
calls=[]
--- FAIL: TestRound4AF2CompletesTheNthAllowedCycleAtCap
```

The AF2 budget comparison should reject a count **greater than** the cap while allowing equality to
finish the already-spent cycle. The ordinary fix-up branch remains responsible for rejecting the
later request to start cycle 6. That threshold correction does not cure marker forgery; receipt
authenticity is part of the CRITICAL finding above.

### [MAJOR] Both CHANGELOGs and the implementation record describe a stronger system than ships

**[PRIMARY — WRONG]** CLI `CHANGELOG.md:24-30` says deleting either source cannot lower the count,
the cycle is spent the moment it runs, and unreadable state escalates. The probes above refute all
three without deleting both sources: a failed attempt has no marker, the counter is written only
after `Fixup` returns, and `LoadCursor` errors are ignored. Lines 14-16 call the unit “completed
fix-up cycles,” while lines 27-28 charge a cycle that ran but failed checks; lines 37-38 still call
the cap “published cycles.” The release record needs one precise unit.

**[PRIMARY — WRONG]** Skill `CHANGELOG.md:22-24` says the printed and enforced values are the same.
The effective cap can be reset per run or by state loss, and AF2 cannot recover the fifth allowed
cycle, so that equality is not established.

**[PRIMARY — WRONG]** `IMPLEMENTATION.md:251-266` says corrupting the cursor cannot buy a cycle and
that unreadable state always escalates. `IMPLEMENTATION.md:274-280` says the cycle is spent “the
moment it runs”; the code spends it only after `Fixup` returns successfully and the cursor save
succeeds. The exact-directory claim at lines 282-286 and the three reversion claims at lines
293-302 are accurate.

## Answers to the requested questions

### Q1 — can anything still lower the count or open a fix-up past the cap?

Yes.

- **Cursor field:** lowering it, corrupting it, or making it unreadable is silently accepted when
  markers do not carry the failed attempt.
- **Crash/order window:** `Fixup` runs before the counter save. A `Fixup` error, process crash, or
  save failure repeats the same cycle on resume.
- **Two runs:** the counter and lock are per run. Sequential runs each spend their own cycle 1;
  concurrent runs can execute two fix-ups and collapse them into count 1.
- **Deleted run directory:** this deletes the only receipt for an attempt whose checks failed, so
  the next run starts at marker count zero.
- **Rolled-back checkout:** both unignored sources are inside the repository and there is no
  external high-water mark; rolling both paths back restores their earlier maximum.
- **Directory matching:** this part is fixed. `isRoundDirName` accepts only a non-empty digits-only
  suffix, and its reversion test is effective.

### Q2 — is “spent when it runs, not when it passes” right?

Yes. I maintain my round-3 position. A code-writing attempt that errors, crashes, or breaks checks
is exactly the churn the finite breaker exists to bound; giving it back creates an unlimited retry
channel. The implementation adopts the right rule in prose but only charges a successful return,
not a run/attempt. I am `codex-1`, so the instruction asking `hermes-1` whether it is persuaded does
not apply to my artifact.

### Q3 — is the AF2 gate correct?

No. It correctly blocks a forged/stale marker whose reconciled count is already **over** the cap,
but `>=` also blocks completion of the Nth allowed cycle. AF2 must be allowed to finish that
transition at equality; cycle N+1 is rejected later if the new review still has agreed fixes.

### Q4 — tests and reversion checks

All mutation and reversion work was confined to isolated copies under `/tmp`; neither shared
working tree was edited by a test.

**[PRIMARY] CLI current-code checks:**

```text
$ go build ./...                                      exit 0
$ go vet ./...                                        exit 0
$ go test ./... -count=1
stock failure: TestDurableKillEndToEndRealProcess
process verification failed (no recorded boot id); not killed
temporary adversarial failure: TestRound4AF2CompletesTheNthAllowedCycleAtCap

$ go test ./... -count=1 \
    -skip '^(TestDurableKillEndToEndRealProcess|TestRound4.*)$'
all packages pass

$ go test ./internal/driver -count=1 -run \
  '^(TestFixupCapIsInclusive|TestZeroFixRoundsDoNotSpendTheFixupBudget|TestFixupBudgetIsTamperFailSafe|TestOnlyExactRoundDirsCountAsPublishedCycles|TestPhaseReviewFixupMarkerSkipsRefixup|TestPhaseReviewFixupChecksFailEscalates|TestBlockedConsensusRespectsTheHardCrossReviewCap)$'
PASS
```

The stock full-suite failure is the same host boot-ID limitation recorded in round 3; I do not
represent the literal full suite as green in this environment.

**[PRIMARY] Skill-package checks in `/tmp/codex-r04-skill.zd2aQ1`:**

```text
$ node scripts/build-addon-manifest.js --check
all six packaged manifests ok
$ npm test
tests 386; pass 386; fail 0
python 3.14: 54 tests OK across 7 files
all six packaged manifests ok
```

**[PRIMARY] Reversion checks in `/tmp/codex-r04-revert.KNrMTo`.** I removed only the monotonic
carry-forward, AF2 gate, and exact-directory predicate in that copy, then ran the three cycle-3
tests. All went red for the intended reason:

```text
FAIL deleting_the_markers_does_not_lower_the_count_when_the_cursor_holds_it
     action=fixup err=<nil>
FAIL a_forged_marker_in_the_current_round_cannot_drive_AF2_past_the_cap
     action=fixup err=<nil> calls=[open-review]
FAIL TestOnlyExactRoundDirsCountAsPublishedCycles
     counted 5 published cycles, want 2
```

**[PRIMARY]** `git diff --check` in both shared repositories was silent/exit 0. `gofmt -d` over all
changed Go files was silent. These were read-only checks.

### Q5 — are both CHANGELOGs and IMPLEMENTATION.md accurate?

No. Their protocol-cell values, exact-directory account, and cycle-3 reversion record are accurate.
Their safety claims about source independence, deletion/corruption, unreadable state, run-time
charging, and effective equality with the printed cap are not.

### Q6 — should anything stop the release?

Yes. The CRITICAL budget reset/collapse paths and the MAJOR AF2 boundary regression stop
parley-deck-cli 1.44.0 and parley-deck-skill 2.8.0. The release should not proceed until the attempt
is reserved before code-writing execution, the high-water mark and lock are idea-scoped across
runs and fail closed on loss/rollback, AF2 recovers the Nth allowed cycle at equality, permanent
tests cover every probe above, and both CHANGELOGs plus `IMPLEMENTATION.md` describe those actual
semantics.

## Response to @hermes-1

On the round-3 disagreement, I maintain that failed code-writing attempts must spend budget. The
cycle-3 rationale is persuasive on that design point. The current ordering does not yet implement
the position: it moved persistence before `RunChecks`, but left it after `Fixup`, preserving the
error and crash window my round-3 counter-proposal explicitly named.
