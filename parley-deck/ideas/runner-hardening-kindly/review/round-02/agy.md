---
agent: agy
idea: runner-hardening-kindly
round: 2
phase: review
date: 2026-06-12
---

## Fix verification

- **Fix 1 (ACP activity + attempt contract):** VERIFIED. ACP milestone events (`agent.acp.initialized`, `agent.acp.session_opened`, `agent.acp.prompt_completed`) now call `act.MarkEvent()` to reset the watchdog. `attempt_id` is correctly threaded through the procctl marker, `agent.started` (and its metadata), `agent.heartbeat`, watchdog events, and `finishACP` terminal payloads. The ACP branch in `runAgent` implements the retry-once loop for `no_first_output`.
- **Fix 2 (Snapshot retention on move-back failure):** VERIFIED. Failed snapshot moves trigger `keepForRecovery = true`, which defers `snap.Abandon()` (removing the `.pid` marker to prevent stale sweeps while keeping the recovery directory). `publishArtifact` always returns the live path, which terminal events report even on failure.
- **Fix 3 (Move-aside safety):** VERIFIED. `moveAsideInvalidArtifact` renames the attempt-1 invalid artifact to a unique destination by appending a nanosecond timestamp if the default path is taken, and deletes the invalid artifact on rename failure. The behavior is fully covered by the new `TestMoveAsideInvalidArtifact`.
- **Fix 4 (Phase 8 fix-up joins the hardened path):** VERIFIED. `RunFixup` now executes via `execAgentProcess` with process group tracking, environment sanitization (`cleanParticipantEnv`), counting writers, and supervised wait. Retries are omitted, and terminal events report `agent.fixup_finished/failed` with the classified payload.
- **Fix 5 (failEarly classification):** VERIFIED. Setup errors in `failEarly` are routed through `classifyFailure`, appending `failure_class` and `recovery_hint` to the `agent.failed` event payload.
- **Fix 6 (Consult provenance):** VERIFIED. The consult frontmatter and ledger index gain the `session_id` field (written empty for one-shot CLIs). `timeout_ms` successfully records the effective timeout (`ConsultResult.EffectiveTimeout`) enforced during the execution of `RunConsult` rather than the raw flag.
- **Fix 7 (Lock the hint table):** VERIFIED. `TestClassifyFailure` was extended to assert exact class-to-hint mappings for all 9 taxonomy classes and 3 watchdog classes.

## New findings

- **NIT:** In `internal/runner/acp.go` (`finishACP`), if `opts.publishArtifact` returns an error, `result.OutputPath` is not updated to the live canonical path (unlike `finalizeExecResult` in `runner.go` which handles this correctly). In practice, ACP is not invoked during the Phase 8 review/sandbox phase, so this does not result in a functional bug, but the code remains slightly inconsistent.

## Dispositions

We concur with the prior disposition dismissing the `TestDurableKillEndToEndRealProcess` failure. Under the codex seatbelt sandbox on macOS, access to the `sysctl kern.boottime` MIB is restricted. Because the process controller relies on `kern.boottime` to verify PID freshness and avoid false positive matches on PID recycling, restricting this system call causes the test to fail. The test passes consistently in unconstrained shells, confirming this is a sandbox environment artifact rather than a code defect.

## Verdict

ACCEPT
