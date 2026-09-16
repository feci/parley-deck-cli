---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
note: n1-blocked-signoff-mirror
context_mode: full
attestation_source_sha256: 4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a
---

# kimi-1 — N1: blocked-consensus mirror in precheckSignoffLaunch + preflight doc scoping

## Task scope (bounded)
Narrow correction ONLY to claude-1's two MINOR findings in
`claude-1-n1-final-source-review-20260916.md`:

1. MINOR-1: `precheckSignoffLaunch` (`internal/app/driver_precheck.go`) omits the
   production blocked-consensus abort (`consensus_request_signoffs.go:94-96`,
   `TriageBlocked`). Mirror it BEFORE `discoverConfigured`/render/telemetry,
   preserving the actual blocked reason. Regression: review consensus frontmatter
   `blocked: false` + one `❌ BLOCK` signoff + a missing signer + invalid protocol
   → no protocol-refusal terminal, no launch; the normal blocked error wins.
   Reuse nearby production fixtures/helpers.
2. MINOR-2: narrow ONLY the doc comments of `PreflightStepCharge`
   (`internal/budget/step_session.go`) and `PreflightCycleCharge`
   (`internal/budget/cycle_session.go`) to new-reservation exhaustion; disclose
   that cached successful sessions still revalidate at the real charge and that
   late refusals stay spent. NO preflight behavior widening; NO budget code change
   beyond comments.

Bounds: 15 min / 60 steps. No Bash/tests/Git/subagents/providers/apply. No test
PASS claims from own execution. Do NOT touch the stale driver migration test,
history helper/tests, or the launcher.

## Independent review posture
Source-reading only; I ran nothing. I concur with claude-1's no-MAJOR/no-CRITICAL
verdict on the two items in scope; this note is not a signoff and not a re-audit.

## Plan
1. Add the `TriageBlocked` gate to `precheckSignoffLaunch` at the same sequence
   position production has it (after the malformed gate, before target/agent
   resolution and any render), with the exact production reason.
2. Add `TestDriverPrecheckReviewSignoffsBlockedConsensusWins` reusing
   `precheckImplOpsFixture`, `refuseAppProtocolCache`, `driverPrecheckRefusals`.
3. Scope the two budget doc comments; add the revalidation/late-spend disclosures.

## Status — COMPLETE (source candidate; coordinator runs tests)

### What changed
- `internal/app/driver_precheck.go` (`precheckSignoffLaunch`): added the
  production blocked-consensus abort — `summary.Triage == consensus.TriageBlocked`
  returns the exact production reason ("target consensus is blocked; resolve the
  BLOCK before requesting more signoffs", `errors.New`, mirroring
  `consensus_request_signoffs.go:94-96`) at the same sequence position: after the
  malformed gate, BEFORE `requestSignoffTargets` / `discoverConfigured` /
  `requestSignoffAgents` / `PrecheckProtocolLaunch` (render + telemetry). Added
  the `errors` import; the function comment now records the triage-gate mirror.
- `internal/app/driver_precheck_test.go`: new
  `TestDriverPrecheckReviewSignoffsBlockedConsensusWins` — review/consensus.md with
  frontmatter `blocked: false`, one `❌ BLOCK` signoff (builder, with Notes +
  Counter-proposal so triage is blocked, not malformed) and reviewer still
  missing, under a tampered phase-7 cached protocol. Reuses
  `precheckImplOpsFixture`, `refuseAppProtocolCache`, `driverPrecheckRefusals`.
  Asserts the exact production blocked reason (not a protocol refusal), zero
  retained records (no `protocol_context_refused` terminal, no launch), and the
  same blocked verdict after cache restore (a valid protocol does not reopen a
  blocked consensus).
- `internal/budget/step_session.go` / `internal/budget/cycle_session.go`: doc
  comments ONLY, scoped to new-reservation exhaustion. Both now disclose: a
  cached SUCCESSFUL session is replayed nil without the charge path's
  attempted-branch revalidation (step: `Current()`, runtime inspect,
  `receipt.check`, `checkTime`; cycle: `checkProtocolMigrationCharges`,
  `receipt.check`, the lost-reserved-charge guard, plus `cycleRefusalKey`
  unconsulted and unreachable from current call sites), and a refusal that only
  materializes at the real charge stays spent — a different budget charged
  between preflight and refusal is not refunded. Widening was rejected per
  claude-1's rationale: false refusals on sessions the charge path would honour
  are the more expensive error.

### Explicit non-claims
- I ran NO compiler, tests, vet or race (no Bash per launch bounds); no PASS is
  claimed from my own execution. Coordinator state at handoff is testimony only:
  combined broad compile + focused 83 passed; changed-full/race/vet still
  running.
- The stale full-run failure
  `TestProtocolMigrationDriverNestedActualChildrenAndBlockResume/4` is known
  intended-expectation drift owned by the coordinator (corrected in a separate
  immutable final snapshot preserving imported+actual charges and child counts);
  I did not touch the driver test, the history helper/tests, or the launcher.
- The two history test fixture corrections (serialized newline, obsolete
  declaration expectation) are separately recorded elsewhere; out of scope here.
- Preflight helper behavior is unchanged — comment narrowing only. No budget
  migration, cap change, refund, or atomicity claim was introduced.
- This note is not a signoff; claude-1's conditional acceptance still stands on
  the pending broad checks, which the coordinator owns.
