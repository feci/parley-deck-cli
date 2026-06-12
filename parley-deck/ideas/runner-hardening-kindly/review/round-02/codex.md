---
agent: codex
idea: runner-hardening-kindly
round: 2
phase: review
date: 2026-06-12
---

## Fix verification

- Fix 1 (ACP activity marks, attempt_id threading, ACP retry-once): VERIFIED. `internal/runner/acp.go` marks initialize/session-open/prompt-complete as activity, threads `attempt_id` through the procctl marker, `agent.started`, heartbeats, watchdogs, and terminal payloads, and `internal/runner/runner.go` now wraps ACP runs in the same retry-once loop for `no_first_output`.
- Fix 2 (snapshot cleanup and live terminal artifact path): NOT-FIXED. Conditional cleanup and `Abandon()` are implemented, and the exec finalizer assigns the live path even on publish failure, but ACP still does not: `internal/runner/runner.go` returns `liveOutput` on move-back error, while `internal/runner/acp.go` only assigns `result.OutputPath = livePath` on the success branch. An ACP review snapshot publish failure can therefore still emit a terminal event with the snapshot artifact path.
- Fix 3 (moveAsideInvalidArtifact safety and test): NOT-FIXED. The implementation now chooses a unique destination and removes the invalid artifact after a rename error, but the agreed rename-failure test case is missing; `TestMoveAsideInvalidArtifact` only exercises the pre-existing `.attempt-1.invalid` case.
- Fix 4 (RunFixup hardened path): VERIFIED. `RunFixup` now calls `execAgentProcess` with a `runID:agentID:fixup` marker, supervised wait, counting writers, participant env cleanup, watchdog events, classified `agent.fixup_failed`, and no retry loop.
- Fix 5 (failEarly classified payload): VERIFIED. `failEarly` now classifies setup failures and writes `failure_class` plus `recovery_hint` on `agent.failed`.
- Fix 6 (consult provenance): VERIFIED. Consult artifacts and the JSONL ledger now include `session_id`, and `timeout_ms` is written from `ConsultResult.EffectiveTimeout`, the effective timeout enforced by `RunConsult`.
- Fix 7 (class/hint contract tests): VERIFIED. `TestClassifyFailure` now locks exact class-to-hint pairs for the implemented provider taxonomy and watchdog classes.

## New findings

None beyond the NOT-FIXED agreed items above.

## Dispositions

- `go build ./...` passed.
- `go test ./internal/runner/ ./internal/app/ ./internal/fsutil/` failed only in `TestDurableKillEndToEndRealProcess` with `process verification failed (no recorded boot id); not killed`; `internal/app` and `internal/fsutil` passed.
- I concur with the prior sandbox-artifact disposition for `TestDurableKillEndToEndRealProcess`. The reproduced failure is the same restricted `kern.boottime`/boot-id attribution path previously recorded for the codex seatbelt sandbox, not a new fix-up regression.
- Targeted runner hardening tests passed: `TestClassifyFailure`, `TestMoveAsideInvalidArtifact`, `TestNoFirstOutputWatchdogRetriesThenFails`, `TestStallGuardKillsAfterFirstOutput`, `TestHeartbeatEventsEmitted`, and `TestReviewSnapshotLifecycle`.

## Verdict

ACCEPT-WITH-FIXES
