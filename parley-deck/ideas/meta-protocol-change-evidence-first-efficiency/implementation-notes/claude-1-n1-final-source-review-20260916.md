---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: bounded independent source review (N1 deltas)
---

# N1 final source review — claude-1

## Scope and method

Source-reading only, against the immutable bundle `.parley-runtime/n1-review-source/`
(`n1.diff` plus targeted production files from the same bundle; never the live
top-level tree). No commands, tests, git, subagents or providers were run. All paths
below are bundle-relative; line numbers are the bundle's.

Reviewed: the 15 N1 Go deltas. The `internal/budget/worktree_inventory.go` delta is a
comment-only restatement already scoped-reviewed elsewhere and I did not re-review it.
Priority per the brief was launch behaviour. Coordinator run status (compile 11.704s,
68 focused events 18.131s PASS, full run still going, race/vet pending) is testimony to
me; I make no claim about any broad check, and nothing below rests on a test result.

**Verdict: no CRITICAL, no MAJOR. Two MINOR findings. Scope acceptance conditional on
the pending full/race/vet run.**

---

## Findings

### [MINOR] `precheckSignoffLaunch` omits the real path's blocked-consensus abort

`internal/app/driver_precheck.go:91-105` reproduces `requestConsensusSignoffs`'
`consensus.Status` + `summary.Errors` + `requestSignoffTargets` gates, but not the
gate between them: `internal/app/consensus_request_signoffs.go:94-96` returns
`"target consensus is blocked; resolve the BLOCK before requesting more signoffs"`
*before* any agent is resolved or launched.

Consequence: in that state the precheck still resolves agents and renders, so on a
known-invalid protocol it records a `protocol_context_refused` terminal attributed to
a signer for a launch production never performs, and the driver escalates with
`"protocol context refused"` instead of the real blocked-consensus reason. This is the
same false-refusal class the implementation deliberately avoids elsewhere — see the
artifact-skip mirror at `internal/runner/launch_precheck.go:66` and its guard test
`TestDriverPrecheckImplementationSkipsExistingArtifact`.

Counterexample (review path, narrow but reachable): a `review/consensus.md` whose
frontmatter is `blocked: false` but which carries one `❌ BLOCK` signoff and still has
other missing signers. `internal/driver/impl.go:208` gates on `rs.Blocked`
(frontmatter), which is a different signal from `summary.Triage == TriageBlocked`
(signoff-derived), so `impl.go:212` calls `RequestReviewSignoffs` → `stepImplOps`
(`internal/driver/budget.go:138`) → `PrecheckReviewSignoffs` → the render above. The
consensus path is not reachable this way: `internal/driver/consensus.go:49-95` routes
`TriageBlocked` to reopen and only calls `RequestSignoffs` on `TriagePartial`.

No budget is spent and nothing is authorized either way, which is why this is MINOR
and not MAJOR. Suggested fix: mirror the triage gate in `precheckSignoffLaunch` (return
the same error, or `nil`, before `discoverConfigured`).

### [MINOR] Preflight helper coverage is narrower than its doc framing claims

Both helpers open with "reports whether the charge this context would attempt is
already KNOWN to refuse". The implemented coverage is narrower than that reads:

- `PreflightStepCharge` (`internal/budget/step_session.go:144-146`) returns the cached
  `s.err` on an attempted session — so `nil` when the prior charge succeeded. `ChargeStep`
  in that same branch (`step_session.go:81-100`) re-validates `Current()`, the runtime
  inspect, `receipt.check` and `checkTime`, any of which can turn that session into a
  refusal.
- `PreflightCycleCharge` (`internal/budget/cycle_session.go:152-154`) likewise skips
  `ChargeCycle`'s attempted-branch revalidation (`cycle_session.go:95-113`:
  `checkProtocolMigrationCharges`, `receipt.check`, lost-reserved-charge), and never
  consults `cycleRefusalKey`, which `ChargeCycle` checks first (`cycle_session.go:75`).

Concrete consequence in `reserveBudget`: a grouped launch whose step session already
charged, but whose step lifetime wall clock has since expired, passes
`PreflightStepCharge` (`internal/runner/launch_budget.go:125`), `ChargeCycle` then
spends a protocol cycle (`launch_budget.go:164`), and only `ChargeStep`
(`launch_budget.go:186`) refuses via `checkTime`
(`internal/budget/step_binding.go:190`). The cycle stays spent.

This is **residual narrowness, not a behaviour regression** — the pre-N1 order
(ChargeCycle → JoinStepSession → ChargeStep) produced the same spend. I recommend
scoping the two doc comments to *new-reservation exhaustion* rather than widening the
helpers: widening the attempted-session branch would risk false refusals on a session
the charge path would in fact honour, which is the more expensive error here. The
`cycleRefusalKey` gap is unreachable from the current N1 call sites
(`internal/driver/cycle_budget.go:135`, `:89` carry no deferred refusal).

---

## What I checked and found sound (so the residual above is not over-read)

**MAJOR-1 (fresh-scope admission before precheck) — mechanism confirmed, coverage
complete.** The poisoning path is exactly
`internal/budget/cycle_history.go:198-201`: `refuseUnmigratedCycles` scans
`.parley-runtime/invocations/*/requested.json` and, for a record whose
`metadata.idea` equals the idea and whose phase maps to the same cycle kind, refuses
the *first* binding permanently. A free precheck refusal writes such a record via
`retainProtocolRefusal` → `recordLaunchRequest`
(`internal/runner/protocol_context.go:147-158`), so the hoist is necessary, not
cosmetic.

Coverage is complete because only two phases map to a cycle kind — `fixup` and
`round-NN` for N≥2 (`cycle_history.go:17-24`) — and both now admit charge-free first:
`internal/driver/budget.go:53` (`admitCrossReviewCycle`, mirroring
`cycleRoundRunner.RunRound`'s exact `crossReviewBase`/`LegacyCycleFloor`/
`EnsureCycleBinding` parameters, `cycle_budget.go:62-69` vs `:99-107`) and
`internal/driver/impl.go:363-366` (mirroring `reserveFixupCycle`'s observer and
`charged` floor, `cycle_budget.go:122-123`). Every other N1 precheck phase —
`implementation`, `review`, `review-consensus`, `consensus`, `goal-check` — is in
`knownNonCyclePhase` (`cycle_history.go:33-37`), and the drafter precheck's
`LaunchInfo{RunID: "one-shot"}` (`driver_precheck.go:78`) records `idea: ""`
(`telemetry.Metadata.Idea` is a plain non-`omitempty` string,
`internal/telemetry/record.go:33`), so `cycle_history.go:192` skips it. Its
`phase: "unspecified"` (`internal/runner/telemetry.go:74-75`) is therefore inert.

I also checked the symmetric **step**-scope poisoning
(`internal/budget/step_history.go:75`, which refuses any idea-scoped invocation phase
other than `preflight`/`round-01`): it is closed **by construction, not by this
change** — `withStepBudget` calls `EnsureStepBinding` at
`internal/driver/budget.go:22`, before it installs the precheck wrappers at `:32-38`,
so no precheck can run before the step scope is bound.

Admission is genuinely hoisted, not forgiven: `EnsureCycleBinding` still runs
`refuseUnmigratedCycles` on first binding (`cycle_binding.go:161`, `:187`) and still
enforces the frozen ceiling (`:145-147`), so existing unmigrated history refuses at
admission with the same verdict, only earlier in the sequence.

**MINOR-1 exhaustion comparison is exact.** `PreflightCycleCharge`'s
`b.Count(state) >= b.Policy.Maximum` (`cycle_session.go:160`) matches the ledger
authority: `cycleLimits` sets `Actions[kind] = Maximum - Carried`
(`cycle_intent.go:27-34`) and the ledger refuses at `count >= ceiling`
(`ledger.go:174`), where `Count = Carried + entries` (`cycle_binding.go:207-218`) —
so the two agree on every value, including `Maximum == 0`, which `cycleLimits` turns
into `Denied` and the preflight refuses via `0 >= 0`. `PreflightStepCharge` delegates
to `StepBinding.Check` (`step_binding.go:157-177`), which is exactly what
`OpenStepSession` runs (`step_session.go:49-51`), so a preflight refusal in
`reserveBudget` implies the subsequent `JoinStepSession` would have refused anyway.

**Mirroring against the real launches.** Implementation is first-only with the real
skip predicate (`phase58.go:34-39` and `runner.go:406` vs `launch_precheck.go:56-57,66`);
an all-artifacts-present state is a genuine no-op (the loop falls through to `return nil`).
Review consensus forces `Overwrite = true` exactly as `RunReviewConsensus` does
(`phase58.go:361-369` vs `launch_precheck.go:43-48`). Cross-review normalization matches
`RunRound`'s (`runner.go:925-933`), and the phase computed by `protocolLaunchPhase`
(`protocol_context.go:161-170`) is the same `round-NN` the real launch carries. Signoff
phase strings are identical to the real `signoffContext`
(`consensus_request_signoffs.go:875-883`), and the driver passes no `ModeOverrides`
(`driver_consensus.go:58-63`, `driver_impl.go:460-466`), so skipping
`applyLaunchModeOverrides` in the precheck is inert. Goal-check mirrors
`GoalCheck`'s checker resolution and `LaunchInfo` field-for-field
(`driver_impl.go:385-398` vs `driver_precheck.go:45-57`), and correctly leaves the
no-checker and unresolved-checker branches as `nil` since neither launches.

**`reserveBudget` reordering is read-only and consistent with existing practice.** The
block hoisted above the charges (`launch_budget.go:76-112`) is `LoadLaunchBinding`,
`config.LoadDefaults`, `RequireMonetaryBinding` and a JSON policy comparison — the
driver already runs exactly this `LoadLaunchBinding` + `RequireMonetaryBinding` pair
before any binding or charge at `internal/driver/budget.go:14-17`, which is good
evidence the hoist is safe rather than an assumption about "Load" naming. Error
precedence changes (a config/monetary error now reports instead of an exhaustion
error), but in every such case the new order charges strictly less.

**No refund or atomicity claim was introduced.** Nothing in the diff refunds; the
monetary `Reserve` is deliberately left after cycle/step admission
(`launch_budget.go:162-172`), so a monetary-ceiling refusal still spends a cycle and a
step. That is disclosed in the header comment (`:64-70`) and is outside the stated
MINOR-1 scope — I flag it as a known residual rather than a finding.

**Existing final-at-cap gates preserved.** The inclusive `charged >= maximum` escalation
(`impl.go:348-350`), the `strict_gate`/named-checks/verifier gates (`impl.go:300-319`)
and the §4.0 cross-review cap on the blocked-consensus reopen path
(`consensus.go:97-111`) are untouched; the hoisted `EnsureCycleBinding` sits after the
cap check, so it cannot admit a scope past its ceiling.

**One consistency note, not a finding.** `roundRunnerAdapter.PrecheckRound` for round 1
would compute artifact paths under `round-01/` while `protocolLaunchPhase` yields
`round-02`, because `runner.PrecheckRound` raises `Round` but leaves a pre-set
`RoundLabel` (`launch_precheck.go:12-17`). The production `roundRunnerAdapter.RunRound`
has the identical normalization (`driver.go:600-605` into `runner.go:925-931`), so the
precheck faithfully mirrors the real adapter, and the driver never dispatches round 1
through this seam (`driver.go:317-357` requires `roundComplete(CurrentRound)` first, so
`next >= 2`). I mention it only so a later reader does not mistake it for a divergence.

---

## Acceptance

I accept the N1 deltas at source level, **conditional on** the pending broad checks
(full test run, `-race`, `go vet`) coming back clean — I ran none of them and make no
claim about them. The two MINOR items above are not merge blockers in my view; MINOR-1
(blocked-consensus mirror) is the one I would fix, as it is a small, local change to a
helper whose entire contract is faithful mirroring.

I reviewed source only, on a bounded budget, and did not review the history helper or
launcher tests owned by other agents. This note is not a signoff and is not a full
audit of the branch.
