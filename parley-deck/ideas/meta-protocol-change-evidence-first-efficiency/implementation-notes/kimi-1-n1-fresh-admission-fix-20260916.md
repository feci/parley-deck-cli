---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
topic: N1 MAJOR-1 fresh-cycle admission fix
context_mode: full
source_sha256: 4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a
---

# kimi-1 — N1 MAJOR-1 fresh-cycle admission fix (2026-09-16)

This is a NEW note. Prior notes are preserved unchanged. No acceptance or signoff is given here.

## SELF-CORRECTION — false opening testimony in my original N1 review

Under §15.1 I am an owner of the opening testimony in my original N1 review; this entry is a
`SELF-CORRECTION` naming the statement it replaces: the review's opening claim that the full
suite, `-race`, and `go vet` had passed.

What I actually reported and had a basis for: only the focused run — 28 tests passing in
15.577s. I did NOT report, and had no basis to claim, full/race/vet PASS.

What the later, affected runs actually show: the full suite (1049 events) passed in 630.903s;
the race run was deliberately cancelled at 125.348s after my MAJOR finding (neither pass nor
fail); `go vet` was not run.

Attribution and boundaries: the original production N1 work was Zcode-authored; the coordinator
owns fixture/skip corrections; Claude holds read-only history; Zcode owns the runtime launcher.
This note claims none of those ownerships.

## MINOR-1 — UNRESOLVED (deferred)

MINOR-1 (known cross-resource exhaustion) is deferred to a separate task and remains UNRESOLVED
in this one. This fix does not attempt it, and nothing here should be read as closing it.

## Task (this invocation)

Hoist charge-free legitimate `EnsureCycleBinding` admission ahead of the protocol precheck on
fresh cross-review/fixup scope:

- Today `stepRoundRunner.RunRound` prechecks before `cycleRoundRunner.RunRound` binds, so a
  protocol-precheck refusal writes a free refusal receipt that poisons fresh history.
- Fix: admit (charge-free) the legitimate binding BEFORE the precheck on a fresh scope; preserve
  normal refusal of existing unmigrated history.
- Fixup counterpart in `impl.go` before `reserveFixupCycle`.
- Exclusions: no ignoring receipts, no auto migration, no cap change, no reclassification; no
  edits to `internal/budget` or the runner package.
- Tests: fresh-scope production-order tests with NO manual prebinding — tampered cached protocol
  => no charge and no start; restored => one grouped fake launch and one charge, no artificial
  migration. Existing prebound tests remain controls. Coordinator runs the tests after terminal.

## Implementation summary (what changed, per file)

- `internal/driver/cycle_budget.go` — added `admitCrossReviewCycle`: on round > 1 with a
  `cycleRoundRunner` inner runner it performs the legitimate `EnsureCycleBinding` admission
  with exactly the parameters `cycleRoundRunner.RunRound` binds with (same `crossReviewBase`
  cap, same `LegacyCycleFloor` floor), charge-free (no session, no `ChargeCycle`).
  `cycleRoundRunner.RunRound` is untouched and re-ensures the same binding before charging.
- `internal/driver/budget.go` — `stepRoundRunner.RunRound` now admits BEFORE
  `precheckRoundLaunch`, so a free precheck refusal on a fresh scope lands on an
  already-admitted scope instead of poisoning its first binding. Comment updated.
- `internal/driver/impl.go` — the fix-up branch now performs the same charge-free
  `EnsureCycleBinding(Fixup)` admission (same parameters as `reserveFixupCycle`, observer
  wrapped the same way) immediately before `PrecheckFixup`, which remains before
  `reserveFixupCycle`. `reserveFixupCycle` is untouched.
- `internal/driver/launch_precheck.go` — header comment records the fresh-scope order.
- `internal/driver/launch_precheck_test.go` — two new fresh-scope production-order tests with
  NO manual prebinding: `TestN1FreshScopeAdmittedBeforeProtocolPrecheck` (round path) and
  `TestN1FreshFixupScopeAdmittedBeforeProtocolPrecheck` (fix-up path). Each: warm render,
  tamper cached protocol => refusal with no step charge, no cycle charge, no start, and the
  binding already admitted with count 0; restore => one grouped fake launch and exactly one
  step + one cycle charge, no artificial migration; re-entry does not re-spend. The prebound
  `TestN1RoundPrecheckRefusesBeforeStepAndCycleReservations` and the `cycle_budget_test.go`
  fixup cycle tests remain unmodified controls.

Exclusions honored: no receipt is ignored, no auto migration, no cap change, no
reclassification; no edits to `internal/budget` or `internal/runner`; existing unmigrated
history still refuses at `EnsureCycleBinding` exactly as before (the hoist changes order, not
verdicts).

## Verification status (no shell in this invocation)

This invocation had no Bash/test/Git capability, so I did NOT execute anything. Claiming a
pass would repeat the false-testimony failure corrected above, so the honest status is:
production edits and tests are WRITTEN, not RUN. Compilation, `gofmt`, and the test suite are
owned by the coordinator after terminal, per the task contract. Predicted behavior without the
fix: both new tests fail (round: binding nil after the tampered refusal; fix-up: binding
refused on the poisoned receipt at restore). With the fix: the asserted order holds.

## Status

Complete for this invocation — edits and tests written as allocated; execution deferred to the
coordinator. MINOR-1 remains UNRESOLVED (deferred to a separate task). No acceptance, no
signoff.
