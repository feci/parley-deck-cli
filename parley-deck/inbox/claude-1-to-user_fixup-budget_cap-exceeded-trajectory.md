---
from: claude-1
to: user
idea: meta-protocol-change-phase-packet-and-fixup-budget
date: 2026-08-12
kind: budget-escalation
status: awaiting-user-decision
---

# Fix-up budget exceeded — trajectory for a human decision

This idea ratified a **5-cycle** fix-up cap for `deliberation`. Its own implementation is at
**cycle 6, with round 7 returning NOT CLEAN**. Under the rule this idea exists to enforce, that is a
blocking escalation, not a reason to start cycle 7.

## Trajectory

| Cycle | Highest severity | Substance | Findings on new or unchanged code |
| --- | --- | --- | --- |
| 1 | CRITICAL ×2 | skill package not installable; cross-review cap bypassed on the BLOCK back-edge | unchanged |
| 2 | CRITICAL | budget derived from implementer-owned prose, fail-open four ways | **new** (cycle 1's fix) |
| 3 | CRITICAL ×2 | defect relocated to `.fixup-done`; AF2 bypassed the cap entirely | **new** (cycle 2's fix) |
| 4 | CRITICAL | reservation taken after the code-writing call; AF2 boundary off by one | **new** (cycle 3's fix) |
| 5 | MAJOR ×2 | the three cycle-4 fixes shipped with no regression test; claims inverted | **new** (cycle 4's fix) |
| 6 | MAJOR ×2 | integration seam untested; seven statements on the superseded unit | **new** (cycle 5's fix) |
| 7 | MAJOR ×2 | cycle 6's own "one vocabulary" claim not fully applied; the corrected count untested | **new** (cycle 6's fix) |

**Severity is falling and the count is not.** Every cycle from 2 onward found a defect in the
*previous cycle's fix* — the exact churn signature the cap was ratified to interrupt.

## Verification status

`go build` (darwin/linux/windows), `go vet`, the full Go suite, the skill `npm test` and its
manifest check are **green**. Two round-7 verdicts (@hermes-1, @kimi-1) were still running when this
was written.

## What is actually left

@codex-1's two round-7 MAJORs are bounded and specified:

1. Two operator-visible messages still say "cycle(s)" for a count that is the maximum of charged
   attempts and completed cycles. With one errored attempt and zero completed cycles it prints
   "after 1 cycle(s)" — factually wrong output.
2. That corrected count has no assertion, so it can regress silently.

Neither is a behaviour defect in the budget itself. Both are naming and test-coverage.

## Options

1. **Authorise cycle 7** with a new explicit finite ceiling. The remaining work is renaming a unit in
   two strings and adding one assertion — but that is what cycles 5 and 6 looked like too.
2. **Ship as-is**, accepting a wrong noun in two escalation messages, and open a follow-up.
3. **Ship the skill 2.8.0 protocol text only** and move the CLI budget work to its own idea.
4. **Abandon/defer** the budget slice.

## Recommendation

**Option 1, capped at one cycle**, with the release blocked if cycle 7 produces another finding in
its own fix. The remaining defects are real and cheap, but the pattern is strong enough that another
open-ended run is not justified.

Recorded rather than decided: **the implementer must not extend its own budget.**
