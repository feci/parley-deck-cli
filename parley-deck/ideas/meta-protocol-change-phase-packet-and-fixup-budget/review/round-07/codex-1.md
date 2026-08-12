---
agent: codex-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 7
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 7

## Summary

The integration-seam finding is fixed. The new test drives a real `00-prompt.md` track through
`New` and the consensus-BLOCK back-edge for both `deliberation` and `standard`; deleting only the
wiring assignment in an isolated copy makes both subtests fail before any runner call.

The vocabulary finding is not fixed. The candidate still describes the combined maximum of cursor
reservations and `.fixup-done` markers as cycles/published cycles in operator-visible errors, code,
tests, both CHANGELOGs, the three protocol copies, and `IMPLEMENTATION.md`. This is not merely
editorial: an executed cursor-only probe had zero completed cycles and one charged attempt, and the
shipped error reported `after 1 cycle(s)`.

There is also a fourth instance of the removal-survives-tests class. Reverting the corrected
already-charged count in the ordinary escalation from `published` to the refused `cycle` leaves all
`internal/driver` and `internal/track` tests green.

## Findings

### [MAJOR] Fix-up cycle 6 did not apply its claimed one-vocabulary correction

**[PRIMARY — WRONG] Source plus executed probe.** Fix-up cycle 6 says the combined count is always a
“charged/reserved attempt” and that “published/completed” is used only for `.fixup-done` markers.
The candidate does not satisfy that claim.

The clearest shipped defect is in the two operator-visible errors:

- `internal/driver/impl.go:152`: **“fix-up budget exceeded: %d cycle(s) recorded”**. `spent` is the
  maximum of the charged-attempt cursor and completed-cycle markers, so it need not be a cycle count.
- `internal/driver/impl.go:307`: **“after %d cycle(s)”**. The argument is the same combined count.

I drove the reservation-only window in an isolated copy: cap 1, the first `Fixup` returns an error,
the cursor records one charged attempt, no `.fixup-done` marker exists, and the next attempt is
refused. The current candidate printed:

```text
$ GOCACHE=/tmp/codex-r07.u1lmlp/gocache-vocab go test ./internal/driver -count=1 -run '^TestR07VocabProbeErroredReservation$' -v
r07_vocab_probe_test.go:36: zero completed cycles, one charged attempt: review still has 1 agreed fixes after 1 cycle(s); cycle 2 would exceed MaxFixupCycles=1; escalating
PASS
```

That is a factual error in what ships: no fix-up cycle completed. One attempt was reserved and
failed.

Other current uses of the superseded unit are:

- `internal/driver/impl.go:292-295`: **“MaxFixupCycles=N allows cycles 1..N”**; `:315-319`:
  **“RESERVE the cycle”**; and `:510-511`: **“publishedFixupCycles reports how many fix-up cycles
  have been CHARGED”**. The unexported helper and local `published` variable also retain the wrong
  unit. The persisted `FixupCyclesPublished` field is a compatibility name and its comment now
  identifies that fact; the unexported names have no such constraint.
- `internal/driver/impl_test.go:392-393`: **“the cap … counts PUBLISHED cycles”**;
  `:454-458`: **“`fast` (cap 1) published none at all … cycles 1..N run”**; `:501-503`:
  **“The ratified unit is published cycles”**; `:623-624`: **“a fix-up that errors has still spent
  its cycle”**; and `:638`: **“the cycle is spent when it runs.”** The last statement is also wrong
  about timing: the attempt is charged before `Fixup` runs.
- `internal/driver/impl_test.go:529-530` still says **“every tamper direction on the fix-up budget
  must be fail-safe. Deleting driver state must not buy a cycle”**, despite the accepted cursor-only
  window, deleted-run-directory, rollback, and concurrent-run limits. This is the round-6 universal
  claim that was supposed to be narrowed.
- CLI `CHANGELOG.md:16-18`: **“a cycle is charged when it is reserved”**; `:28-31`: **“The cycle is
  reserved before the code-writing call”**; and `:52`: **“5 inclusive reserved cycles.”**
- Skill `CHANGELOG.md:13`: **“cap 5 cycles”** and `:16-17`: **“charges a cycle when it is RESERVED.”**
- The current §4.0 cell still says **“cap 5 cycles”** in `parley-deck/COOPERATION.md:229`,
  `internal/protocol/defaults/COOPERATION.md:220`, and the skill snapshot at
  `skills/parley-deck/references/COOPERATION.md:220`.
- `IMPLEMENTATION.md` still contains unqualified current-sounding uses: `:72` **“published fix-up
  cycles”** as the governing unit; `:142-148` **“the ratified unit … published cycles”** and
  **“rewritten onto the same unit”**; `:279-291` **“The cycle is spent the moment it runs”**; and
  `:317-318` **“The cycle is now reserved.”** Its initial `:38` `fast` publication claim is later
  withdrawn at `:152-155`, but it remains contradictory without an inline supersession marker.

Uses of “published/completed” that refer only to `.fixup-done` marker fixtures are correct and are
not included above. Likewise, the frozen `FINAL.md` records the original published-cycle design; it
should not be edited, but `IMPLEMENTATION.md` must identify the charged-attempt implementation as a
review-mandated deviation instead of continuing to alternate units.

Concrete fix: use **charged/reserved attempt** for every combined count, error and current budget
statement; use **published/completed cycle** only for `.fixup-done` markers. Rename the unexported
helper/local variable, retain the persisted field name only with its explicit legacy-name note,
change the three protocol cells and both release notes to “5 charged attempts,” narrow the tamper
comment, remove the withdrawn `fast` assertion, and mark superseded historical passages in
`IMPLEMENTATION.md` explicitly rather than flattening the audit history.

### [MAJOR] The corrected escalation count can still be removed with the shipped tests green

**[PRIMARY] Isolated one-line reversion.** Round 6 required an exact assertion that the ordinary
over-cap error reports the already-charged count, not the refused next ordinal. No test asserts that
message. In a fresh isolated candidate copy I changed only the error argument at
`internal/driver/impl.go:307`:

```diff
- ..., published, cycle, d.cfg.MaxFixupCycles)
+ ..., cycle, cycle, d.cfg.MaxFixupCycles)
```

That restores “after 6 cycle(s)” at cap 5. The shipped tests for both affected packages remained
green:

```text
$ GOCACHE=/tmp/codex-r07.u1lmlp/gocache-message go test ./internal/driver ./internal/track -count=1
ok  parley-deck-cli/internal/driver  1.027s
ok  parley-deck-cli/internal/track   0.476s
```

This is the fourth occurrence of the same assurance defect: behavior named in the release record
can disappear while its owning package stays green. Add a behavioral assertion for the exact
already-charged count. The reservation-only probe above is the stronger fixture because it also
forces the corrected unit: zero completed cycles, one charged attempt, attempt 2 refused.

## Answers to the requested questions

### Q1 — are both round-6 findings fixed?

No. The seam finding is fixed; the vocabulary finding is not.

**[PRIMARY — CONFIRMED] Seam mutation.** Current candidate:

```text
$ go test ./internal/driver -count=1 -run '^TestTrackWiresTheHardCrossReviewCapThroughNew$' -v
PASS deliberation
PASS standard
```

I copied the candidate to `/tmp/codex-r07.u1lmlp/repo`, deleted only:

```go
cfg.HardCrossReviewCap = pol.CapCrossReviewRounds
```

and reran the test with an isolated Go cache. Both subtests failed exactly at the seam:

```text
track deliberation: HardCrossReviewCap=0, want 3 — the policy is not wired into the driver
track standard: HardCrossReviewCap=0, want 2 — the policy is not wired into the driver
FAIL parley-deck-cli/internal/driver
```

The broader mutated command failed in `internal/driver` for those two cases while
`internal/track` stayed green, which is the expected localization. No runner call is reached.

### Q2 — is the vocabulary consistent everywhere?

No. The exact remaining quotations are listed in the first MAJOR. The two most consequential are
the shipped errors **“%d cycle(s) recorded”** and **“after %d cycle(s)”**, because the executed
cursor-only case proves those counts can contain an errored reservation and therefore are not
completed/published cycles.

### Q3 — is there a fourth removal-survives-tests case?

Yes. Replacing the correct already-charged argument with the refused next ordinal in the ordinary
escalation leaves `go test ./internal/driver ./internal/track -count=1` green. The mutation and
output are in the second MAJOR.

### Q4 — ship or not?

Do not ship this candidate. The specific defect in what ships is that its operator output and
release/protocol documentation label a combined charged-attempt count as completed/published
cycles. The executable counterexample is zero completed cycles being reported as “after 1
cycle(s).” The correct N-vs-N+1 diagnostic also still has no regression-catching assertion.

The settled trust-anchor scope call remains accepted and is not a blocker here. The packet slice is
also not used as a reason to reject this budget-only release.

## Validation record

**[PRIMARY] Commands executed by this reviewer.** All probes and mutations were confined to
`/tmp/codex-r07.u1lmlp`; the shared working trees were not modified by them.

```text
$ go test ./internal/driver ./internal/track -count=1
ok  parley-deck-cli/internal/driver
ok  parley-deck-cli/internal/track

$ go vet ./...
PASS (silent)

$ npm test  # ../parley-deck-skill
386 Node tests passed; 54 Python tests passed; all six add-on manifests OK

$ git diff --check  # CLI repository and skill repository
PASS (silent in both)
```

`go test ./... -count=1` passed every package except the unrelated existing
`internal/runner.TestDurableKillEndToEndRealProcess`, which failed because this environment supplied
no recorded boot id. The same runner failure occurred in the isolated message-mutation copy; both
candidate-relevant packages passed there, so I do not attribute that environmental failure to this
release diff.
