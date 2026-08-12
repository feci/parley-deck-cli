---
agent: codex-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 6
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 6

## Summary

The new permanent test does catch the three round-4 regressions: against the combined reversion,
the reservation, AF2-equality, and corrupt-cursor subtests go red; AF2 beyond the cap stays green,
as it should. The trust-anchor scope call remains accepted and is not reopened here.

The release is still not clean. The unit/failure-domain correction did not reach all shipped
claims, and the ordinary over-cap path reports one more cycle than was actually charged. Separately,
the policy-to-driver assignment that makes the consensus-BLOCK cap real has no regression test: I
removed that one assignment in an isolated copy and both relevant packages remained green.

## Findings

### [MAJOR] The corrected unit and boundary claims are still internally contradictory, and the shipped escalation count is wrong

**[PRIMARY — WRONG] Source plus executed probe.** The main `publishedFixupCycles` comment now states
the reservation-only window accurately, but several other release-candidate statements still
describe completed or published cycles, or universal tamper safety:

- CLI `CHANGELOG.md:12`: **“`MaxFixupCycles = N` publishes cycles 1..N”**. It permits/reserves
  attempts 1..N; an errored reservation need not publish a cycle.
- Skill `CHANGELOG.md:16`: **“At either boundary the run halts”**. Equality is deliberately allowed.
  The run escalates only when another round/attempt would exceed the inclusive cap; merely reaching
  the last allowed attempt does not halt a successful transition.
- `internal/driver/cursor.go:52-53`: **“count of completed fix-up cycles.”** The cursor includes a
  reserved attempt even when `Fixup` errors or never completes.
- `internal/driver/impl_test.go:100-103`: **“The §4.0 budget counts driver-completed cycles”** and
  `:501-503`: **“The ratified unit is published cycles”**. The corrected unit is a reserved attempt.
- `internal/driver/impl_test.go:456-458`: **“`fast` (cap 1) published none at all.”** That is the
  withdrawn end-to-end `fast` claim; the test injects `auto_implement` into a route `fast` forbids.
- `internal/driver/impl_test.go:529-530`: **“every tamper direction on the fix-up budget must be
  fail-safe.”** The now-documented cursor-only window and the settled trust-anchor limits are
  counterexamples. The two cases in that test cover only the post-marker two-record state.
- `internal/driver/impl_test.go:638`: **“the cycle is spent when it runs.”** It is spent when
  reserved, before the code-writing call.

There is also a user-visible off-by-one in the ordinary escalation. `internal/driver/impl.go:305-307`
sets `cycle := published + 1` and then reports that value as the number of cycles already elapsed.
With five charged attempts and cap 5, my isolated probe printed:

```text
$ go test ./internal/driver -count=1 -run '^TestR06ProbeOverCapMessage$' -v
r06_probe_test.go:29: review still has 1 agreed fixes after 6 cycle(s) (MaxFixupCycles=5); escalating
PASS
```

Only five attempts have been charged; attempt 6 was refused before reservation. The two read-error
messages at `internal/driver/impl.go:146` and `:303` likewise still say “cycles have been published”
although the value being protected is the maximum charged/reserved count.

Concrete fix: consistently say **charged/reserved attempt** where the combined count is meant;
reserve “published/completed” for `.fixup-done` markers only. Change the skill boundary sentence to
“when another round or attempt would exceed the cap.” In the ordinary over-cap error, report the
already-charged `published`/`spent` value (5), not the refused next ordinal (6), and add an assertion
for the exact count in the escalation.

### [MAJOR] The consensus-BLOCK hard-cap wiring can disappear with the shipped tests green

**[PRIMARY] Isolated one-line reversion.** `internal/driver/driver.go:143-148` is the integration
seam that copies `Policy.CapCrossReviewRounds` into `Config.HardCrossReviewCap`. The production
guard is tested only by `TestBlockedConsensusRespectsTheHardCrossReviewCap`, which manually supplies
`HardCrossReviewCap: 3` (`internal/driver/consensus_test.go:369-374`). Conversely, the `New` tests
check the initially scheduled clamp but do not assert `d.cfg.HardCrossReviewCap` for either
`standard` or `deliberation` (`internal/driver/track_test.go:35-67`).

I copied the candidate to `/tmp/codex-r06-seam-WYhHlg/repo`, removed only:

```go
cfg.HardCrossReviewCap = pol.CapCrossReviewRounds
```

and ran:

```text
$ go test ./internal/driver ./internal/track -count=1
ok  parley-deck-cli/internal/driver
ok  parley-deck-cli/internal/track
```

That reversion restores the release-blocking behavior: explicit `standard` and `deliberation`
drivers still clamp the initially scheduled count, but a consensus BLOCK can again walk past the
printed cap. Both CHANGELOGs explicitly claim that back-edge is bounded, so this is release behavior,
not incidental plumbing.

Concrete fix: add a table-driven behavior test that constructs the driver through `New` from an
explicit `standard` and `deliberation` prompt, reaches a blocked consensus at each cap, and asserts
escalation plus no runner call. That test must not inject `HardCrossReviewCap` directly.

## Answers to the requested questions

### Q1 — do the permanent tests catch the combined regression?

Yes.

**[PRIMARY] Executed current and reverted code.** Current candidate:

```text
$ go test ./internal/driver -count=1 -run '^TestRound4FixesHavePermanentTests$' -v
PASS an_errored_fix-up_keeps_its_reservation
PASS AF2_finishes_the_last_allowed_cycle_at_equality
PASS AF2_refuses_beyond_the_cap
PASS a_corrupt_cursor_escalates,_an_absent_one_does_not
```

I then copied the candidate to `/tmp/codex-r06-S0niM6/repo` and combined exactly the three round-4
reversions: save the cursor after `Fixup`, change AF2 from `spent > cap` to `spent >= cap`, and
restore the cursor-error-swallowing `err == nil && ...` path. The same command produced:

```text
FAIL an_errored_fix-up_keeps_its_reservation
     driver.json: no such file or directory
FAIL AF2_finishes_the_last_allowed_cycle_at_equality
     fix-up budget exceeded: 2 cycle(s) recorded against MaxFixupCycles=2
PASS AF2_refuses_beyond_the_cap
FAIL a_corrupt_cursor_escalates,_an_absent_one_does_not
     a corrupt cursor must escalate rather than be read as zero
```

Thus all three reverted fixes have a red case. The beyond-cap case correctly remains green because
both `>` and `>=` reject a value already beyond the cap; equality is the discriminating assertion.
In the combined corrupt/absent subtest, the corrupt-cursor assertion is the one that goes red; the
absent-cursor behavior itself is unchanged by the reversion.

### Q2 — are the corrected claims now exact?

No. The exact remaining wrong quotations are listed in the first MAJOR. The corrected
`publishedFixupCycles` comment and the two CHANGELOG qualifications about the cursor-only window are
right, but they coexist with stale completed/published-cycle claims, a still-universal tamper-safety
claim in the tests, an inaccurate skill boundary sentence, and an operator-visible N+1 count.

### Q3 — is any other shipped behavior missing a regression-catching test?

Yes. Removing the policy-to-driver `HardCrossReviewCap` assignment leaves
`internal/driver` and `internal/track` green while re-enabling the consensus-BLOCK bypass for both
tracks. Also, the ordinary over-cap tests assert only escalation/no `Fixup`; none asserts the
already-charged count in the message, which is why the current “after 6 cycle(s)” defect passes.

I found no additional untested production branch in the round-4 reservation, AF2 comparison, or
corrupt-versus-absent cursor fixes themselves.

### Q4 — should this release ship?

No. The specific defects in what ships are: at cap 5 the CLI tells the operator that six cycles have
already occurred although only five attempts were charged; release/source/test text still
contradicts the reserved-attempt and limited-redundancy model; and the one assignment that enforces
the claimed consensus-BLOCK cap can be removed without a shipped test failing. These are narrow
fixes. The settled out-of-repository trust-anchor follow-up is not a blocker and is not reopened.

## Validation record

**[PRIMARY] Commands executed by this reviewer.** All mutations and probes were confined to the two
isolated `/tmp` copies above; the shared working trees were not edited by them.

```text
$ go test ./internal/driver -count=1 -run '^TestRound4FixesHavePermanentTests$' -v  # current: pass
$ go test ./internal/driver ./internal/track -count=1                               # current: pass
$ go test ./internal/driver -count=1 -run '^TestRound4FixesHavePermanentTests$' -v  # combined reversion: 3 red
$ go test ./internal/driver ./internal/track -count=1                               # cap-wiring reversion: pass
$ node scripts/build-addon-manifest.js --check                                      # all six manifests ok
$ npm test                                                                           # exit 0
$ git diff --check                                                                    # silent, both repositories
```
