---
agent: codex-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 1
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — release candidate slice 1

## Summary

**[PRIMARY — CONFIRMED]** The simple outstanding-fixes comparison is numerically inclusive, and the
deliberation policy preserves the reviewer fields and LE-11 floor; the source locators and executed
checks are recorded under Q1 and Q3 below. The release is nevertheless not ready: the hard
cross-review cap is bypassed after a consensus BLOCK, strict-gate-only rounds consume the counter
that is now described as published fix-up cycles, the claimed fast-track behavior is tested through
an unreachable synthetic configuration, the required trajectory escalation payload is absent, and
the sibling skill package rejects its own stale payload manifest.

## Findings

### [CRITICAL] The deliberation cross-review cap is bypassed on the consensus BLOCK back-edge

**[PRIMARY — WRONG]** The release claim that deliberation is capped at three cross-review rounds
after round 1 and then escalates is false on the auto-driver path. `internal/track/track.go:160-165`
sets `CapCrossReviewRounds: 3`, and `internal/driver/driver.go:128-136` clamps the configured initial
round budget. After rounds 1 through 4, `advanceRound` drafts consensus
(`internal/driver/driver.go:270-284`). If that consensus is BLOCKED,
`internal/driver/consensus.go:90-122` ignores `CrossReviewRounds`, opens another Phase-2 round, and
uses only `MaxRounds`. `driver.New` defaults `MaxRounds` to 4
(`internal/driver/driver.go:95-104`), so at round 4 the calculation is `next=5`, and the guard
`next > 1 + MaxRounds` is false. Round 5 therefore runs: four cross-review rounds after round 1.
Escalation occurs only when another BLOCK tries to open round 6.

**[PRIMARY — CONFIRMED]** The current tests confirm the two disconnected behaviors rather than their
required composition:

```text
$ go test ./internal/driver -run 'TestDeliberationClampsCrossReviewRoundsToThree|TestConsensusBlocked(ReopensRound|MaxRoundsEscalates)' -count=1 -v
--- PASS: TestConsensusBlockedReopensRound
--- PASS: TestConsensusBlockedMaxRoundsEscalates
--- PASS: TestDeliberationClampsCrossReviewRoundsToThree
PASS
```

**[PRIMARY — WRONG]** This does not satisfy
`parley-deck/ideas/meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md:67-75,97-100` or
`parley-deck/ideas/meta-protocol-change-phase-packet-and-fixup-budget/consensus.md:65-76`: the cap is
a blocking escalation threshold, not merely the point at which the first consensus draft is
attempted.

Suggested fix: preserve a distinct hard cross-review cap in `driver.Config` instead of collapsing
it into the initial scheduling count. Before `advanceConsensus` reopens a BLOCKED explicit
deliberation consensus, escalate when opening the next round would exceed three rounds after round
1. Add a behavioral test starting from deliberation round 4 with a BLOCK and prove that no round-05
runner call occurs and that the blocking escalation is written.

### [CRITICAL] The 2.8.0 packaged skill payload manifest is stale and installation fails closed

**[PRIMARY — WRONG]** The sibling repository is not a releasable 2.8.0 package. The release changes
`skills/parley-deck/references/COOPERATION.md` and
`skills/parley-deck/references/compatibility.json`, but does not update the tracked
`skills/parley-deck/parley-addon.json` hashes (`parley-addon.json:7-11`). The package's own manifest
check reports:

```text
$ node scripts/build-addon-manifest.js --check
parley-deck: STALE parley-addon.json — regenerate it
```

**[PRIMARY — CONFIRMED]** This is functional, not cosmetic. My independent `npm test` run reported
`tests 386`, `pass 271`, `fail 115`; the first cascading installer error was:

```text
Source payload does not match parley-addon.json:
modified: references/COOPERATION.md; modified: references/compatibility.json
```

Suggested fix: regenerate the add-on manifest with the repository's canonical manifest command,
include the resulting `skills/parley-deck/parley-addon.json` change in the candidate, then rerun
`npm test` and the manifest check.

### [MAJOR] The “fast now publishes one cycle” case is not an actual fast-track path

**[PRIMARY — WRONG]** The changelog's fast row (`CHANGELOG.md:10-17`) and the new test name claim an
actual fast idea moved from zero to one published fix-up cycle. The test does not create a fast
idea: every subtest writes only `auto_implement: true` and directly injects
`MaxFixupCycles: tc.cap` (`internal/driver/impl_test.go:451-465`). Thus its “fast” cases are an
absent-track driver with a synthetic cap of 1.

**[PRIMARY — CONFIRMED]** The production combination is unreachable. `PolicyFor(Fast, ...)` rejects idea-level
`auto_implement` (`internal/track/track.go:142-149`), and the app derives `cfg.AutoImplement` from
that same idea field (`internal/app/app.go:1941-1946,1995-2000`). With a valid non-auto fast idea,
the new inclusive guard passes at round 1 but `internal/driver/impl.go:287-289` then escalates before
`Impl.Fixup` is called. The focused run demonstrates that both the synthetic case and the production
contradiction test pass simultaneously:

```text
$ go test ./internal/driver -run 'TestFastContradictionEscalates|TestFixupCapIsInclusive/fast' -count=1 -v
--- PASS: TestFixupCapIsInclusive/fast:_cap_1_still_publishes_one_cycle
--- PASS: TestFixupCapIsInclusive/fast:_cycle_2_escalates
--- PASS: TestFastContradictionEscalates
PASS
```

**[PRIMARY — CONFIRMED]** The old `>=` predicate itself would reject a synthetic cap-1 driver before
calling `Fixup`: `HEAD:internal/track/track.go:142-149` supplies cap 1, while
`HEAD:internal/driver/impl.go:277-289` checks `cycle >= cap` first. That proves the arithmetic
off-by-one, but it does not prove the claimed end-to-end fast-track before/after behavior.

Suggested fix: decide and test the real fast Phase-8 route. If fast fix-up is intentionally manual,
remove the claim that this driver change makes fast publish one cycle and document where the manual
cap is enforced. If the driver must orchestrate it, add a production-reachable transition that does
not violate the `fast + auto_implement` prohibition. In either case, replace the synthetic track
labels with a real `track: fast` integration test.

### [MAJOR] Strict closing rounds consume a budget now defined as published fix-up cycles

**[PRIMARY — WRONG]** Leaving the second comparison as `round >= MaxFixupCycles` is locally correct
if it is a ceiling on total review-round ordinals: round N may certify clean, but an uncertified round
N must not open N+1 (`internal/driver/impl.go:201-230`). Simply changing that comparison to `>` would
allow review round N+1 and is not the right fix.

**[PRIMARY — WRONG]** The shared counter is nevertheless inconsistent with the new contract. A strict-gate round with
zero agreed fixes can open the next review round without publishing a fix-up
(`internal/driver/impl.go:201-230`). Later, the outstanding-fixes branch sets `cycle := round`
(`internal/driver/impl.go:277-294`). For example, with cap 5, two uncertified zero-fix strict rounds
can advance to round 3 without any `Fixup` call; the first later fix is numbered 3, and round 6
escalates even if fewer than five fix-ups were actually published. Therefore the release does not
guarantee
`parley-deck/ideas/meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md:67-70`'s “5 inclusive
published cycles” under `strict_gate`.

**[PRIMARY — CONFIRMED]** The prior verification-honesty review explicitly accepted the shared
budget because its then-binding requirement was only termination bounded by `MaxFixupCycles`
(`parley-deck/ideas/verification-honesty/review/consensus.md:71-76`). That evidence does not resolve
the new conflict: this FINAL strengthens the unit from a generic termination budget to a count of
published fix-up cycles.

Suggested fix: keep the strict closing comparison's inclusive review-round boundary, but give the
strict-close retry loop its own durable counter/ceiling and derive the fix-up cap from completed
published fix-up cycles (for example, validated `## Fix-up cycle` records), not from the highest
review-round number. Add a mixed-sequence test with zero-fix strict rounds followed by fix-ups.

### [MAJOR] Neither budget boundary produces the ratified trajectory payload

**[PRIMARY — WRONG]** `consensus.md:65-76` requires the escalation payload to contain trajectory,
findings by severity, fresh-vs-relitigated status, unresolved fixes, validation status, and a
recommendation. No Go implementation reads or emits those fields:

```text
$ rg -n 'trajectory' internal --glob '*.go'
(no matches; exit 1)
```

**[PRIMARY — WRONG]** At the fix-up boundary, `internal/driver/impl.go:284-285` returns only the outstanding-fix count,
round-derived cycle, and cap. `Driver.Run` wraps that as the generic `driver-error` inbox note at
`internal/driver/loop.go:50-59,202-227`. The cross-review boundary does not currently escalate at
all, as the first finding shows. The changelog statements that either boundary raises a blocking
trajectory escalation (`CHANGELOG.md:29-33`; sibling `CHANGELOG.md:16-21`) are therefore overstated.

Suggested fix: implement a dedicated budget-escalation builder that derives and records every
ratified field, preserves the completed count across a finite extension, and has artifact-level
tests for both Phase 2 and Phase 8 boundaries. Until then, remove the trajectory and extension
claims from the release notes.

### [NIT] Two comments still say explicit deliberation applies no overrides

**[PRIMARY — WRONG]** `internal/driver/driver.go:109-112` says “absent/deliberation preserve today's
behaviour byte-for-byte,” and `internal/track/track.go:115` says `ApplyOverrides: false` represents
“absent/deliberation.” Explicit deliberation now returns `ApplyOverrides: true`
(`internal/track/track.go:150-166`). These comments should say only an absent/unknown track takes
the legacy no-override path.

## Answers to the requested questions

### Q1 — inclusive fix-up boundary

**[PRIMARY — CONFIRMED]** For a driver configuration that reaches the outstanding-fixes branch, the
new predicate is numerically correct: `cycle := round`; `cycle > N` escalates, so cycles 1 through N
reach `Impl.Fixup` and N+1 does not (`internal/driver/impl.go:277-299`). The independent boundary run
passed N/N+1 for 1, 2, and 5:

```text
$ go test ./internal/driver ./internal/track -count=1 -v -run 'TestFixupCapIsInclusive|TestPhaseReviewMaxFixupCyclesEscalates'
--- PASS: TestFixupCapIsInclusive/at_the_cap_runs_the_last_allowed_cycle
--- PASS: TestFixupCapIsInclusive/one_past_the_cap_escalates
--- PASS: TestFixupCapIsInclusive/standard:_cycle_2_of_2_runs
--- PASS: TestFixupCapIsInclusive/standard:_cycle_3_escalates
--- PASS: TestFixupCapIsInclusive/fast:_cap_1_still_publishes_one_cycle
--- PASS: TestFixupCapIsInclusive/fast:_cycle_2_escalates
--- PASS: TestPhaseReviewMaxFixupCyclesEscalates
```

**[PRIMARY — WRONG]** This does not establish the actual fast-track result, and strict-gate-only
rounds break the equivalence between `round` and published fix-up count; see the two MAJOR findings.

**[PRIMARY — CONFIRMED]** The changed guard has no other branch consumers. It globally affects any
`driver.Config.MaxFixupCycles`, including the absent-track default of 3 (now three rather than two
fix-up calls), explicit standard 2, explicit deliberation 5, and direct test/custom configs. The
strict-close guard at `impl.go:215`, pipeline guard at `internal/pipeline/review.go:72`, zero-fix close
path, check gates, and completion gates are separate and unchanged.

### Q2 — the second strict-gate guard

**[PRIMARY — CONFIRMED]** Keeping `round >= MaxFixupCycles` is the correct local comparison for its
existing total-review-round ceiling: a clean round N can close because the guard is not entered, and
an unclean round N cannot open N+1. The focused strict tests passed opening, clean completion, and
ceiling escalation.

**[PRIMARY — WRONG]** Reusing `MaxFixupCycles` and the review-round ordinal is no longer consistent
with the new published-fixup-cycle contract. The fix is separate counters, not changing this second
operator. This inconsistency can prematurely exhaust a strict idea's actual fix-up allowance and
should be fixed before release.

### Q3 — deliberation `ApplyOverrides: true`

**[PRIMARY — CONFIRMED]** No unintended production reviewer override was found.
`PolicyFor(Deliberation)` leaves `MaxReviewers=0`, `MinReviewers=0`, and
`CrossReviewRounds=-1` while setting only the cap and fix-up fields
(`internal/track/track.go:150-166`). `newDriverImplOps` caps the real reviewer list only when
`pol.MaxReviewers > 0`, so deliberation retains every non-implementer
(`internal/app/driver_impl.go:39-73`). `driver.Config.MaxReviewers` itself has no runtime consumer
outside derivation/tests.

**[PRIMARY — CONFIRMED]** `MinReviewers=0` is not copied over an existing value, and the production
zero value becomes 2 at `internal/driver/driver.go:139-147`; LE-11 still escalates when the review
consensus has fewer than two reviewers (`internal/driver/impl.go:244-254`). Explicit deliberation
with zero available non-implementers still fails at `internal/track/track.go:136-140`. The focused
non-solo and single-reviewer tests passed.

**[PRIMARY — CONFIRMED]** `CrossReviewRounds=-1` preserves the configured/default count, and values
above 3 are clamped. That is the intended cell change, although the separate BLOCK back-edge fails
to respect the hard cap as described in the CRITICAL finding.

### Q4 — rewritten tests

**[PRIMARY — CONFIRMED]** The three existing rewrites are honest. The deliberation policy tests
replace only the two superseded legacy expectations and retain assertions for no reviewer cap,
LE-11 minimum 2, and preserved configured/default cross-review count
(`internal/driver/track_test.go:44-67`; `internal/track/track_test.go:164-190`). The old
round-03 escalation fixture moved to round 4 while retaining both “must escalate” and “must not run
fixup” assertions (`internal/driver/impl_test.go:364-385`). No prior guarantee was deleted.

**[PRIMARY — CONFIRMED]** The gap is in new coverage, not dishonest rewriting: the new cross-review
test observes only the config clamp, and the new “fast” boundary cases do not instantiate fast.
Those gaps allowed the substantive findings above.

### Q5 — protocol text

**[PRIMARY — CONFIRMED]** The two changed table rows are byte-identical in all three copies. A full
diff shows the embedded default and bundled skill snapshot differ only in their expected bootstrap
header values; the live deck differs only in its project header/sync metadata and generated roster
tables. The normalized in-repo drift guard passed:

```text
$ go test ./internal/protocol -run '^TestEmbeddedDefaultMatchesLiveDeck$' -count=1 -v
--- PASS: TestEmbeddedDefaultMatchesLiveDeck
PASS
```

**[PRIMARY — CONFIRMED]** A case-insensitive search across all three protocol files found no remaining
`unbounded`, `uncapped`, or `no cap` assertion for either changed cell. The only match in that search
family is `0 means unlimited` for the separate driver-step/wall-clock/cost ceilings
(`parley-deck/COOPERATION.md:671` and corresponding template lines 662).

**[PRIMARY — WRONG]** The Phase-2 wording “then escalate” is not accurate against the consensus
BLOCK path, so identical text does not mean correct text-to-code alignment.

### Q6 — changelogs

**[PRIMARY — CONFIRMED]** Excluding this new idea's still-in-progress `IMPLEMENTATION.md`, I counted
69 historical implementation files and reproduced the published distribution exactly:

```text
17 0
34 1
 7 2
 2 3
 3 4
 2 5
 1 9
 1 14
 1 15
 1 25
```

**[PRIMARY — CONFIRMED]** The historical `>=` guard made a synthetic cap-1 auto-driver reject round
1 before `Fixup`, which is the arithmetic source of “zero.” **[PRIMARY — WRONG]** Calling that a real
fast-track before/after result, saying both deliberation cells are enforced, and claiming trajectory
payloads at both boundaries are overstated for the reasons in the findings. The sibling changelog's
statement that “the CLI enforces both cells” (`../parley-deck-skill/CHANGELOG.md:19-21`) is also false.

### Q7 — release gate

**[PRIMARY — CONFIRMED]** This release should stop for the two CRITICAL findings and the three MAJOR
behavior/contract findings. The stale package manifest alone makes 2.8.0 fail its own installer
integrity checks; the cross-review cap and trajectory contract are slice-1 requirements, not packet
work deferred to slice 2.

## Independent validation record

**[PRIMARY — CONFIRMED]** These checks passed in this working tree:

```text
go vet ./...                                                    exit 0
GOOS=darwin|linux|windows GOARCH=amd64 go build ./...           exit 0
focused driver/track boundary and LE-11 tests                   PASS
focused protocol normalized-drift test                         PASS
gofmt -d on all changed Go files                               no output
git diff --check (both repositories)                           exit 0
```

**[PRIMARY — CONFIRMED]** `go test ./... -count=1` was not fully green in this sandbox: every listed
package passed except unchanged `internal/runner`, where
`TestDurableKillEndToEndRealProcess` failed three isolated reruns because `sysctl -n kern.boottime`
returns `Operation not permitted`, leaving no boot ID. `git diff --quiet -- internal/runner
internal/procctl` returned 0. I treat that as a disclosed environment-specific qualification, not a
slice-1 defect; it does not mitigate the independently reproducible package-manifest and behavior
failures above.
