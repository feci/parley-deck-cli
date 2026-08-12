---
agent: codex-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 3
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 3

## Summary

The fail-open counter is not closed. **[PRIMARY — WRONG]** The new source is described as
driver-owned, but a `.fixup-done` file has neither access separation nor authenticated contents:
the participant and driver processes can create, delete, or replace it in the same workspace. A
deleted marker lowers the count; deleting the review tree restores zero; an unreadable child round
is silently omitted; and a fabricated marker in the current canonical round trips AF2 before review
completion and skips both `Fixup` and the post-fix-up check gate. The same bare file is being asked
to prove two different facts — budget consumption and crash-safe completion — without proving
either's provenance.

**[PRIMARY — CONFIRMED]** Counting only markers written after successful checks also changes the
unit materially: in an isolated cap-1 probe, two consecutive fix-ups whose checks failed both ran,
and the published count remained zero after each. A check failure halts that driver tick, so this is
not an unattended spin; it still hands the budget back on every resume and does not deliver a
finite fix-up-attempt breaker.

The `standard` BLOCK back-edge itself is correct and is now disclosed in both CHANGELOGs. Its
diagnostic count is one too high, however. The tests and reversion check are genuine; they do not
cover the state-integrity failures above.

## Findings

### [CRITICAL] The replacement marker ledger is forgeable, deletable, and still fail-open

**[PRIMARY — WRONG]** `internal/driver/impl.go:490-524` does not establish the claimed
driver-only provenance. It lists every directory whose name merely starts with `round-`, then uses
`fileExists` (`internal/driver/cursor.go:316-319`) to count any regular `.fixup-done` file. It does
not validate an exact canonical round name, a contiguous sequence, marker contents, the expected
cycle, the reviewed consensus, the checked HEAD, or which process created the file.

I copied the current module to an isolated directory under `/tmp`, added a same-package adversarial
probe only there, and ran:

```text
$ go test ./internal/driver -count=1 -run '^TestMarkerCounterAdversarialReview$' -v
missing review dir -> count=0 err=<nil>
agent-created marker in odd round dir -> count=1 err=<nil>
agent-deleted marker -> count=0 err=<nil>
marker replaced by non-regular file -> count=0 err=<nil>
unreadable child round containing marker -> count=0 err=<nil>
unreadable review dir -> err=... permission denied
PASS
```

**[PRIMARY — CONFIRMED]** The outcomes attack each requested case:

- A missing `review/` directory is expressly converted to `(0, nil)` at
  `internal/driver/impl.go:509-512`. On disk rebuild, an implementation with no remaining review
  tree returns to PhaseImpl (`internal/driver/cursor.go:159-168`); after round 1 is reopened, the
  budget starts from zero. Missing historical state therefore extends the budget.
- An unreadable `review/` root returns an error and eventually escalates. That case is correctly
  fail-closed. But an unreadable child `round-NN/`, an unreadable/missing marker, or a marker
  replaced by a directory is collapsed by `fileExists` to `false`; the function returns a lower
  count with no error. The advertised “a read error is returned, never swallowed” comment at
  `internal/driver/impl.go:504-505` is false.
- An arbitrary marker in `round-banana/` counts. Odd names inflate rather than extend the budget by
  themselves, but they prove the counter does not identify one valid receipt per actual review
  round.
- The same participant-level filesystem access used by the isolated review probe could create and
  delete the markers. “Driver-written” describes the normal writer in one branch; it is not an
  enforced ownership property.

**[PRIMARY — CONFIRMED]** The interaction with AF2 makes fabrication more serious than a false
count. `internal/driver/impl.go:132-153` checks the current round's bare marker before
`ReviewRoundComplete`, `ReviewStatus`, `Fixup`, or `RunChecks`. The existing
`TestPhaseReviewFixupMarkerSkipsRefixup` writes arbitrary contents (`"done"`) and passes while
asserting that `Fixup` was not called. Thus a workspace-writing agent can create
`review/round-NN/.fixup-done` and make the driver archive/advance as though fixes and checks had
succeeded. Conversely, deleting a genuine marker in the AF2 crash window makes recovery re-run the
fix-up and lowers the budget count.

This is release-blocking for the same reason as round 2: the constrained party can still lower or
reset the safety count. Suggested fix:

1. Separate a monotonic **attempt reservation** from the AF2 **successful completion receipt**.
   Reserve/increment before invoking code-writing `Fixup`, so crashes and failed checks consume an
   attempt rather than handing it back.
2. Put the monotonic record outside every participant's writable sandbox. A same-user, same-workspace
   bare file cannot substantiate “driver-owned.” If the launcher cannot enforce that separation,
   the release must not make an ownership claim and must use fail-closed redundant reconciliation.
3. Validate exact canonical round/cycle identity, a contiguous sequence, structured receipt
   contents, the reviewed consensus, check result and relevant HEAD. Missing, unreadable, duplicate,
   malformed, or contradictory historical state must escalate.
4. Do not let an unauthenticated budget record serve as AF2 completion proof.

### [MAJOR] Failed post-fix-up checks hand the budget back

**[PRIMARY — CONFIRMED]** The order at `internal/driver/impl.go:304-315` is `Fixup` →
`RunChecks` → write marker. `TestPhaseReviewFixupChecksFailEscalates` explicitly asserts that no
marker is written when checks fail. I then ran this isolated probe with `MaxFixupCycles=1` and the
same review consensus twice:

```text
$ go test ./internal/driver -count=1 -run '^TestFailedPostFixupChecksDoNotSpendBudgetProbe$' -v
failed attempt 1 with cap=1 -> marker count=0 err=<nil>
failed attempt 2 with cap=1 -> marker count=0 err=<nil>
PASS
```

Both calls invoked `Fixup`; both escalated only after the check failure; neither spent the cap of
one. **[PRIMARY — WRONG]** That is not the finite breaker justified in `FINAL.md:72-75`: a failed
code-mutating fix-up is at least as relevant to churn and cost as a successful one. It also differs
from the historical heading distribution used to select five: `RunFixup` validates the updated
`IMPLEMENTATION.md` before the driver's post-fix-up check (`internal/runner/phase58.go:125-141`),
so a failed check can leave a durable `## Fix-up cycle` section that the old/calibration count saw
but the new count erases.

If “published” is intentionally redefined to mean only check-passing cycles presented to the next
review round, the code still needs a separate finite attempt budget. Otherwise repeated resumes
grant unlimited failed fix-up attempts without the recorded finite extension required by FINAL.

### [MAJOR] The release notes and implementation record describe a different, stronger system

**[PRIMARY — WRONG]** The CLI `CHANGELOG.md:14-17` still says the count comes from
`## Fix-up cycle N` records in `IMPLEMENTATION.md`; the code now counts markers. Its unconditional
claim at `CHANGELOG.md:9-12` that cap N publishes 1..N and escalates at N+1 is false when markers
are missing, deleted, unreadable below the root, or withheld by failed checks. The statement at
`CHANGELOG.md:60-63` that both boundaries halt with “the counts” they enforced is also too strong:
the cross-review diagnostic is off by one, and the fix-up count is not validated.

**[PRIMARY — WRONG]** The sibling `../parley-deck-skill/CHANGELOG.md:22-24` says the printed and
enforced numbers are the same. The configured maximum is five, but the effective breaker is not
five when its count can be lowered or reset.

**[PRIMARY — WRONG]** `IMPLEMENTATION.md:208-211` says the markers are the driver's own and that a
read error escalates instead of restarting at zero. The isolated probe shows neither statement is
true as an enforcement property: participants can mutate the files, missing root state returns
zero, and child-marker stat errors are swallowed. The stale source comment at
`internal/driver/impl.go:284-285` likewise still says the cycle number is derived from published
cycles in `IMPLEMENTATION.md`.

Correct these claims after the state model is fixed; changing the prose alone is insufficient.

### [NIT] The BLOCK-cap diagnostic reports one unrun cross-review round

**[PRIMARY — CONFIRMED]** The guard correctly prevents the next round, but
`internal/driver/consensus.go:95-98` prints `next-1`. For `standard`, after allowed rounds 2 and 3,
the attempted round 4 is rejected with:

```text
standard derived HardCrossReviewCap=2 CrossReviewRounds=2
standard BLOCK after round 3 -> action=escalated
err=consensus still blocked after 3 cross-review round(s) after round 1 (§4.0 cap=2)
runnerCalls=[]
```

Only two cross-review rounds completed; the third was not opened. Use `next-2` for completed
cross-review rounds, or say that opening cross-review round `next-1` was refused.

## Answers to the requested questions

### Q1 — is the fail-open counter actually closed?

No. **[PRIMARY — CONFIRMED]** Root unreadability is closed, but root absence returns zero;
unreadable/missing child markers are silently omitted; odd round names and arbitrary marker
contents count; and participants can create/delete the same files. Deletion or removal of the
review tree extends the budget. Fabrication of the current canonical marker additionally triggers
AF2 and bypasses Fixup/checks. The new source therefore relocates the round-2 defect rather than
closing it.

### Q2 — should failed-check fix-ups count?

Yes, at least against an attempt breaker. **[PRIMARY — CONFIRMED]** The current marker definition
does under-count and hand budget back: the isolated cap-1 probe ran two failed fix-ups while the
count stayed zero. A check failure stops the current tick, which limits unattended damage, but a
resume repeats the same unspent cycle. “Successful publication” may remain a separate metric; it
cannot be the only finite limit on code-mutating attempts.

### Q3 — is `standard`'s bounded BLOCK back-edge right and disclosed?

Yes, apart from the diagnostic NIT. **[PRIMARY — CONFIRMED]** An isolated real-policy probe read
`track: standard`, derived `HardCrossReviewCap=2` and `CrossReviewRounds=2`, escalated before round
4, and made zero runner calls. This matches the existing §4.0 standard cell and is the correct
enforcement of that already-printed cap. CLI `CHANGELOG.md:33-38` and skill
`CHANGELOG.md:16-18` now both explicitly disclose the `standard` consequence.

### Q4 — independent tests and reversion check

All mutation work was confined to copied repositories under `/tmp`; I did not edit either shared
working tree.

**[PRIMARY — CONFIRMED]** Current-code checks:

```text
$ go build ./...                         exit 0
$ go vet ./...                           exit 0
$ go test ./... -count=1                 all packages pass except the environment-specific test below
--- FAIL: TestDurableKillEndToEndRealProcess
process verification failed (no recorded boot id); not killed
$ go test ./... -count=1 -skip '^TestDurableKillEndToEndRealProcess$'
all packages pass
$ go test ./internal/driver -count=1 -run \
  '^(TestFixupCapIsInclusive|TestZeroFixRoundsDoNotSpendTheFixupBudget|TestPhaseReviewFixupMarkerSkipsRefixup|TestPhaseReviewFixupChecksFailEscalates|TestBlockedConsensusRespectsTheHardCrossReviewCap)$'
ok
$ go test ./internal/protocol -count=1 -run '^TestEmbeddedDefaultMatchesLiveDeck$'
ok
```

The sole full-suite failure is the same sandbox/host boot-id limitation recorded in my round-2
review; it is outside the changed packages. I do not represent the literal full suite as green in
this environment.

**[PRIMARY — CONFIRMED]** Skill package checks in the isolated sibling copy:

```text
$ node scripts/build-addon-manifest.js --check
parley-deck: ok (6 files, sha256:bee732321d0f5279c5ef83ae308a64e8533337e46be6cf185dae13680a1363db)
$ npm test
tests 386; pass 386; fail 0
python 3.14: 54 tests OK across 7 files
all six packaged manifests ok
```

**[PRIMARY — CONFIRMED]** Reversion check: in the isolated CLI copy I replaced the marker count
with the prior `IMPLEMENTATION.md` heading parser and ran the published-count tests. The three
one-past-cap cases all turned red because they ran fix-up cycle 1 instead of escalating:

```text
FAIL cap_5:_the_6th_escalates
FAIL cap_2:_the_3rd_escalates
FAIL cap_1:_the_2nd_escalates
reversion_rc=1
```

I restored the isolated copy and verified its `impl.go` SHA-256 exactly matched the shared source,
then reran the focused suite green. The test genuinely detects reversion to prose counting; it
does not test missing/deleted/forged markers or failed-check attempt reuse.

**[PRIMARY — CONFIRMED]** `git diff --check` is clean in both shared repos, and `gofmt -d` on all
changed Go files produced no output.

### Q5 — anything still overstated?

Yes. **[PRIMARY — WRONG]** The CLI CHANGELOG names the obsolete heading source and overstates
inclusive enforcement; the skill CHANGELOG overstates equality between printed and effective
enforcement; IMPLEMENTATION overstates marker ownership and read-error handling; and the Go source
comment still names `IMPLEMENTATION.md` as the source. The standard back-edge disclosure itself is
now present and accurate. The implementation's reversion claim is also accurate; I reproduced it.

### Q6 — anything that should stop the release?

Yes. The CRITICAL state-integrity defect and the MAJOR failed-attempt under-count both stop
parley-deck-cli 1.44.0 and parley-deck-skill 2.8.0. Do not release until the counter is monotonic,
participant-inaccessible or equivalently fail-closed under reconciliation, separate from AF2
completion proof, and inclusive of failed/crashed code-writing attempts; the adversarial cases
above are permanent regression tests; and both CHANGELOGs plus IMPLEMENTATION match the resulting
semantics.

## Responses to the round-2 reviewers

### @hermes-1

Your round-2 trust-model defense no longer applies to the source now under review. The new marker
is not protected by that trust model or by filesystem permissions; I created and deleted it using
participant-level access, and the driver accepted arbitrary contents. More importantly, the prior
round's own standard was that a careless edit must not have the same effect as a malicious one.
Marker deletion and review-tree loss still do.

### @kimi-1

I agree with your round-2 diagnosis of the prose source. The in-flight fixture collision you saw
has been repaired: the current inclusive-cap fixture puts markers only in completed prior rounds,
and the focused test is green. The deeper AF2 collision remains architectural, however: the same
unauthenticated file still both increments the budget and authorizes skipping Fixup/checks.

## Open questions

None. The release blocker is code-grounded; no operator policy decision is needed to establish it.
