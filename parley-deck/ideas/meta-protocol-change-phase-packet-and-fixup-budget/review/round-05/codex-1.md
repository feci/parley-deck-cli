---
agent: codex-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 5
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 5

## Summary

The three narrow round-4 defects are fixed: the cursor reservation precedes `Fixup`, AF2 permits
equality and rejects only a count over the cap, and an existing unreadable cursor escalates while an
absent cursor remains the fresh-run case. I verified each behavior with a focused isolated probe and
then reverted each fix in the isolated copy; every probe went red for the intended reason.

The stronger idea-scoped/concurrent/out-of-repository property does **not** have to ship in this
release. It is a materially different persistence and trust model not specified by `FINAL.md`'s
binding implementation conditions. Correct claims plus the named `fixup-budget-trust-anchor`
follow-up are an acceptable release disposition for that residual.

The claims are not yet corrected accurately, however. The new pre-call reservation creates a state
in which the cursor is the **only** receipt: no `.fixup-done` marker exists until after `Fixup` and
`RunChecks`. Deleting that one cursor record after an errored fix-up reopens the budget. Both
CHANGELOGs still claim robustness to loss of one record, and the CLI CHANGELOG still calls the
charged unit completed/published cycles even though a reservation can be charged before `Fixup`
starts or completes. The round-4 production fixes also have no permanent regression tests: reverting
all three simultaneously leaves the repository's current `internal/driver` suite green.

## Findings

### [MAJOR] The corrected release limits still overstate the implementation

**[PRIMARY — WRONG] Source and isolated reproduction.** `internal/driver/impl.go:317-336` writes
`FixupCyclesPublished` to `driver.json` before calling `Fixup`; the marker is written only after
`Fixup` returns and `RunChecks` passes. Therefore an errored, interrupted, or not-yet-started
reservation has one record, not two. `internal/driver/driver.go:241-243` deliberately treats a
missing cursor as a fresh run.

I exercised exactly that composition in `/tmp/codex-r05-cli`: cap 1, reserve cycle 1, make `Fixup`
return an error, delete only `driver.json`, then advance again. The second advance entered `Fixup`
again:

```text
$ go test ./internal/driver -count=1 \
    -run '^TestRound5DeletingTheSoleReservationRecordReopensBudget$' -v
=== RUN   TestRound5DeletingTheSoleReservationRecordReopensBudget
--- PASS: TestRound5DeletingTheSoleReservationRecordReopensBudget
PASS
```

`PASS` is an attack reproduction: deletion of one record lowered the effective count from 1 to 0.
This is within the expressly deferred trust-anchor surface and need not be fixed in code here, but
the release descriptions must state it.

The following CLI CHANGELOG claims are still factually wrong or materially incomplete:

> “The unit is now completed fix-up cycles.”

> “Losing one record does not lower the count”

> “The budget is robust against accidental loss, a stale or deleted single record, and an errored
> or crashed fix-up.”

> “Fix-up (Phase 8): 5 inclusive published cycles”

The implementation charges a **reserved fix-up attempt**, including one that errors and one lost
after reservation but before the code-writing call is confirmed. It does not always charge a
completed or published cycle. Loss of either record alone is redundant only **after** both records
exist; during the reservation-only interval, losing the cursor loses the count.

The skill CHANGELOG's statement is overstated for the same reason:

> “The CLI's enforcement is robust against accidental loss of its own records”

The source comment at `internal/driver/impl.go:512-522` also still says it reports completed cycles,
that deleting the cursor cannot lower the count because markers remain, and that every tamper
direction is fail-safe. Those statements omit the reservation-only state and would mislead the next
maintainer even after the CHANGELOG is corrected.

Concrete fix: describe the charged unit as a reserved fix-up attempt. State separately that, within
one intact run, an errored/crashed attempt remains spent while its cursor survives; after successful
checks and marker creation, loss of either cursor or marker alone cannot lower the count. State that
loss of the cursor before marker creation, changing/deleting the run, rollback, workspace edits, and
concurrent runs are not protected. Apply the same qualification to both CHANGELOGs and the source
comment.

### [MAJOR] The repository test suite accepts all three round-4 regressions

**[PRIMARY — WRONG] Isolated combined reversion.** I copied the release candidate to
`/tmp/codex-r05-stock-revert`, made only these three reversions, and ran the unmodified repository
tests:

1. moved the cursor save back after `Fixup`;
2. changed AF2's comparison from `spent > cap` back to `spent >= cap`;
3. restored the silently ignored cursor-error path.

```text
$ go test ./internal/driver -count=1
ok  parley-deck-cli/internal/driver  0.846s
```

Thus the claimed fixes can all disappear together without making a shipped test red. The focused
probes I added only in the isolated review copy do catch them: current code passes all five cases;
each one-fix reversion fails its corresponding case. Add permanent equivalents for an errored
`Fixup` retaining its reservation, AF2 at equality and over the cap, and malformed-versus-absent
cursor handling.

## Answers to the requested questions

### Q1 — are the three defects fixed?

Yes.

**[PRIMARY — CONFIRMED] Reservation ordering.** `internal/driver/impl.go:322-327` assigns and saves
`FixupCyclesPublished` before invoking `d.cfg.Impl.Fixup`. My cap-1 probe made `Fixup` return
`code-writing failed`; `Advance` escalated and `LoadCursor` returned
`FixupCyclesPublished == 1`.

**[PRIMARY — CONFIRMED] AF2 equality.** `internal/driver/impl.go:145-152` rejects only
`spent > MaxFixupCycles`. With cap 1 and count 1, AF2 opened round 2 without re-running `Fixup`; with
count 2, it escalated and did not call `OpenReviewRound`.

**[PRIMARY — CONFIRMED] Cursor errors.** `internal/driver/driver.go:235-247` distinguishes
`os.IsNotExist` from every other `LoadCursor` error. A malformed existing `driver.json` returned
`ActionEscalated` with “unknown fix-up budget” and made no implementation calls; an absent cursor
continued as a fresh run.

Focused current-code result:

```text
$ go test ./internal/driver -count=1 -run \
  '^(TestRound5ReservationPrecedesErroredFixup|TestRound5AF2AllowsEqualityAndRejectsOverCap|TestRound5CursorErrorEscalatesButAbsenceDoesNot)$' -v
PASS TestRound5ReservationPrecedesErroredFixup
PASS TestRound5AF2AllowsEqualityAndRejectsOverCap/equality_finishes_the_allowed_cycle
PASS TestRound5AF2AllowsEqualityAndRejectsOverCap/strictly_over_cap_escalates
PASS TestRound5CursorErrorEscalatesButAbsenceDoesNot/malformed_existing_cursor
PASS TestRound5CursorErrorEscalatesButAbsenceDoesNot/absent_cursor_is_a_fresh_run
```

The three isolated reversions produced, respectively:

```text
reservation: driver.json: no such file or directory                         FAIL
AF2 equality: action=escalated opened=0 err=fix-up budget exceeded          FAIL
cursor error: action=await err=<nil>; want unknown-budget escalation        FAIL
```

### Q2 — is the scope disposition acceptable?

Yes, as a release decision. `FINAL.md`, `## 2. The budgets`, ratifies the two numeric cells and its
binding implementation conditions require same-patch text/code, inclusive boundaries with tests,
and the deliberation cross-review wiring. It does not specify an idea-wide lock, a cross-run ledger,
or an out-of-repository authority. The live protocol also remains a cooperative, file-canonical
system and expressly does not yet provide OS-sandbox enforcement.

The stronger property would change persistence, portability, reconciliation, concurrency, and the
trust boundary. Shipping it silently in this fix-up would exceed the ratified design. Without it,
the documented residual remains real: a lost run/cursor before marker publication, a rollback, or
concurrent runs can reset or duplicate a reservation. That is acceptable only as an explicit limit,
not as hidden behavior. `fixup-budget-trust-anchor` is an adequate named follow-up.

### Q3 — do both CHANGELOGs accurately describe what ships?

No. The exact wrong quotations and the required narrowing are in the first MAJOR finding. The
security-boundary disclaimers are directionally correct, but they do not cure the contradictory
claims that one-record loss is safe and that the charged unit is completed/published cycles.

### Q4 — anything new in the round-4 diff?

Yes: the correct pre-call reservation introduced the reservation-only state, but the new release
wording still describes the steady-state max-of-two design as though both records always exist. The
round-4 production changes also arrived without regression tests; the stock suite accepts all three
reversions. I found no additional production-logic defect in the three fixes themselves.

### Q5 — should this release ship?

No, not as currently written and tested. The stronger trust-anchor system does not block 1.44.0 /
2.8.0, but the release should wait for accurate unit/failure-domain wording in both CHANGELOGs and
the stale source comment, plus permanent regression tests for the three round-4 fixes. With those
narrow corrections, my scope objection is withdrawn and the code paths reviewed here are ready.

## Validation record

**[PRIMARY] Commands executed by this reviewer.** All mutation and reversion work was confined to
isolated copies under `/tmp`; the shared working trees were not edited by tests.

```text
$ go build ./...                                                        exit 0
$ go vet ./...                                                          exit 0
$ go test ./... -count=1 -skip '^TestDurableKillEndToEndRealProcess$'   all packages pass
$ node scripts/build-addon-manifest.js --check                          all six manifests ok
$ npm test                                                              386 Node tests pass;
                                                                        54 Python tests pass;
                                                                        all six manifests ok
$ git diff --check                                                      silent, exit 0 (both repositories)
```

The literal unskipped Go suite has the same host-specific stock failure recorded in my prior rounds:
`TestDurableKillEndToEndRealProcess` cannot verify the boot ID in this sandbox. I do not represent
that literal suite as green here.
