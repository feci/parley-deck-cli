---
agent: codex
idea: runner-hardening-kindly
round: 1
phase: review
date: 2026-06-12
---

## Summary

Reviewed the implementation diff against consensus D1-D12, including the runner supervision, classifier, artifact-wins, snapshot, consult, driver/app, docs, and protocol-adjacent changes. `go build ./...` passes. `go test ./internal/...` fails only at `internal/runner TestDurableKillEndToEndRealProcess` with `process verification failed (no recorded boot id); not killed`, matching the known codex seatbelt disposition.

The implementation has strong coverage for the exec happy path, D2 config parsing, heartbeat exclusions, Result.Success consumers, git optional locks, and the protocol-copy changes. I found several consensus-compliance gaps that should be fixed before this ships.

## Findings

### [MAJOR] ACP supervision does not satisfy the D1/D3 activity and attempt contract

internal/runner/acp.go:83, internal/runner/acp.go:84, internal/runner/acp.go:145, internal/runner/acp.go:160, internal/runner/acp.go:175, internal/runner/acp.go:319

D1 says ACP activity is marked on session updates and protocol events, while `agent.started` itself never satisfies first output. Here only `SessionUpdate` calls `MarkEvent`; `agent.acp.initialized`, `agent.acp.session_opened`, and `agent.acp.prompt_completed` are appended without marking activity, so a live ACP initialize/session sequence can still be classified as `no_first_output`. D3 also says `attempt_id` threads through `agent.started` and the procctl marker, and retry-once applies to `no_first_output`; ACP uses marker `runID:agentID`, omits `attempt_id` from `agent.started`, and explicitly records ACP watchdog action as failed rather than retrying.

Suggested fix: add an ACP attempt id, include it in the durable marker and `agent.started`, mark initialize/session-open/prompt-complete as ACP protocol activity, and either implement the same one-retry `no_first_output` loop for ACP or record an explicit consensus deviation before review consensus accepts it.

### [MAJOR] Snapshot move-back failures delete the recovery artifact

internal/runner/runner.go:392, internal/runner/runner.go:407, internal/runner/runner.go:608, internal/runner/runner.go:633

D9 requires `review.snapshot_artifact_move_failed` to report a recovery path and retain the snapshot. `runAgent` defers `snap.Cleanup()` immediately after snapshot creation, so if `MoveArtifactBack` fails, the event is emitted but the snapshot directory is deleted on return. The terminal failure also continues to use the snapshot `outputPath` unless publish succeeds, so the recorded artifact path can point at a soon-deleted temp tree rather than the live canonical path D9 requires.

Suggested fix: make snapshot cleanup conditional. On move-back failure, retain the snapshot directory, preserve or replace the marker so stale sweeping will not immediately erase the recovery artifact, and make terminal events report the live canonical artifact path with the snapshot path only as recovery metadata.

### [MAJOR] Attempt-1 invalid artifact move-aside can overwrite an existing recovery file

internal/runner/runner.go:445

D3 says an invalid attempt-1 artifact is moved aside to `<artifact>.attempt-1.invalid`, but never overwrites a pre-existing one. The current `os.Rename(outputPath, outputPath+".attempt-1.invalid")` does not check whether the destination exists; on Unix-like systems it can replace that file. The error is also ignored, so a failed move-aside can leave the invalid artifact in the canonical path before the retry.

Suggested fix: move aside only to a non-existing destination, using an existence check plus a unique fallback suffix or fail the retry safely. Add a test with a pre-existing `.attempt-1.invalid` file and a rename-failure case.

### [MAJOR] Phase 8 fix-ups bypass the hardened exec path

internal/runner/phase58.go:72, internal/runner/phase58.go:94, internal/runner/phase58.go:99, internal/runner/runner.go:939, internal/runner/runner.go:951

`RunFixup` still uses `CommandFor` plus direct `cmd.Run`, with plain stdout/stderr files. That bypasses the new `execAgentProcess` path that sets process groups, procctl markers, `cleanParticipantEnv`, counting writers, `waitSupervised`, watchdog events, and the shared terminal failure behavior. As a result, a Phase 8 claude implementer can still inherit host session markers, a fix-up can hang without D1 heartbeats/watchdogs, and fix-up failures are not surfaced through the same `agent.failed` payload contract described in D5.

Suggested fix: run fix-up attempts through the same supervised exec path, with `ValidateFixupArtifact` as the phase validator and the existing artifact-wins decision table. At minimum, apply `cleanParticipantEnv`, process-group kill, supervision, and the same classified terminal payload builder to `RunFixup`.

### [MINOR] Early setup failures still lack D5 classification payloads

internal/runner/runner.go:706

D5 says one payload builder is shared by exec, ACP, `failEarly`, and fix-up terminal paths, and that `agent.failed` gains `failure_class`, `recovery_hint`, `exit_code`, `signal`, and `stderr_tail_bytes` where applicable. `failEarly` still appends a bare `agent.failed` with only `error` and no classification fields. Spawn/config/path failures are exactly where actionable hints are useful.

Suggested fix: route `failEarly` through the same classified failure payload helper used by normal terminal failures, using bounded tails when logs exist and `classifyFailure` over the setup error text when they do not.

### [MINOR] Failure recovery hints are not the exact agreed strings

internal/runner/failclass.go:27, internal/runner/failclass.go:37, internal/runner/failclass.go:39, internal/runner/failclass.go:41, internal/runner/failclass.go:49, internal/runner/failclass.go:50, internal/runner/failclass.go:52

D5 adopted agy's failure-class table with exact hint strings. Several implemented hints are paraphrases, for example auth, sandbox, budget, invalid-request, no-first-output, stalled, and unknown. This is not a behavioral crash, but it breaks the promised UX contract and makes tests unable to assert the consensus table verbatim.

Suggested fix: copy the exact round-01 hint strings into `failclass.go` and extend `TestClassifyFailure` to assert representative class/hint pairs, not only non-empty hints.

### [MINOR] Consult artifacts omit required provenance fields

internal/app/consult.go:110, internal/app/consult.go:118, internal/runner/consult.go:65

D10 requires consult frontmatter to include `session_id` and `timeout_ms` as part of the canonical provenance. The artifact currently omits `session_id`, and when `--timeout` is not supplied it records `timeout_ms: 0` even though `RunConsult` uses the agent's effective timeout. That makes the durable consult artifact less traceable than the consensus-specified schema.

Suggested fix: return the effective timeout and session id, when discoverable, from `RunConsult`; write `session_id` in frontmatter even if empty/unknown, and record the effective timeout rather than the raw flag value.

## Dispositions

- Finding/disposition: `internal/runner TestDurableKillEndToEndRealProcess` fails under the codex seatbelt sandbox with "no recorded boot id". Prior disposition: dismissed as environment artifact. I concur. The failure reproduced exactly in `go test ./internal/...`, and the scope here does not change the underlying `sysctl kern.boottime` restriction. This should remain a non-blocking sandbox artifact, with normal-shell verification still required before release.

- Finding/disposition: ACP-path claude marker shedding is deferred because claude-as-ACP-participant is not in the active roster. Prior disposition: accepted trade-off. I concur for marker shedding specifically. That concurrence does not cover the separate ACP supervision and attempt-id gaps reported above; those are D1/D3 contract issues, not merely the deferred `acp.MergedEnv` cleanup.

## Verdict

ACCEPT-WITH-FIXES

The exec-path core is close, but D1/D3 ACP gaps, D9 snapshot retention, D3 move-aside safety, and Phase 8 hardening coverage need fixes before this implementation should be accepted.
