# hermes-1 implementation handoff — D7 liveness observation + BuffersStdout supervision

- idea: meta-protocol-change-evidence-first-efficiency
- agent: hermes-1
- branch: feature/meta-protocol-change-evidence-first-efficiency/hermes-1
- code commit: 4b95baffc21baae1a071473d88160aa3cfce0e8b
- date: 2026-09-05
- status: implementation committed; awaiting codex-1 preflight call-site glue + non-owner review

## Owned slice (D7 — "Liveness Is Observation, Not Diagnosis From Silence")

Exact files (from IMPLEMENTATION.md allocation row 54):

- internal/app/preflight.go (modified)
- internal/app/preflight_test.go (modified)
- internal/app/preflight_liveness.go (new)
- internal/app/preflight_liveness_test.go (new)
- internal/runner/supervision.go (modified)
- internal/runner/buffered_test.go (new)

(preflight_hash_test.go is in my allocation row but was already present and unchanged on the
base branch; my diff does not touch it.)

## What changed

1. Typed readiness observation (preflight_liveness.go). The probe seam `probeFunc`
   changed from `func(...) (bool, string)` to `func(...) readinessObservation`. The
   hosted-PONG probe now returns one of seven typed classes instead of a bool+reason:
   ready, malformed-reply, process-exited-empty, deadline-no-output,
   deadline-after-output, provider-failure, process-failure. Only `ready` (an exact
   PONG assistant response extracted from a recognized envelope) marks Available.
   classifyReadiness is pure and unit-tested. An echoed instruction, bullet, fence,
   or malformed JSON containing "PONG" is malformed-reply, never ready. A JSON error
   envelope is a failure even when its nested subtype claims success.

2. Gate construction (preflight.go). A non-ready entry now raises one of three gates
   by class:
   - missing / empty class / plain non-provider process failure → gateExcludeAgent
     (existing explicit operator exclusion; `--yes` still records the exclusion);
   - provider-failure → gateProviderFailure (blocking, never auto-excluded, never
     waivable by `--yes`);
   - malformed/empty/deadline → gateResolveReadiness (blocking, never auto-excluded,
     never waivable by `--yes`; the resolve gate advertises `parley preflight --dir <root>`
     with no `--yes`).
   The §1 non-solo hard-stop (`< 2` would-be participants → exit 1) is unchanged.

3. Provider/process split (preflight_liveness.go). A minimal provider-side classifier
   (rate-limit / auth / billing / overloaded / model-not-found) mirrors
   internal/runner/failclass.go and only feeds the provider-vs-process split. The
   runner's full classifier stays authoritative for run-level classification.

4. Soft-watchdog suppression for declared-buffered transports (runner/supervision.go).
   `supervisionForAgent` now honors the pre-existing `agents.Spec.BuffersStdout`
   (discover.go:60, config runtime `buffers_stdout`): when declared, both
   FirstEventTimeout and StallTimeout are zeroed. The hard per-agent timeout and the
   caller-owned process-group cleanup remain, and the heartbeat (which never counts
   as activity) stays on.

No other changes. I did NOT touch runner.CommandFor, any telemetry path, Kimi's
evidence code, the shared IMPLEMENTATION.md narrative, or any file outside the list
above.

## Actual tests (all passing)

Focused app package (real fake-child fixtures via hostedPONG + pure classifier table):

- TestClassifyReadiness (16 sub-cases, pure classifyReadiness table)
- TestClassifyReadinessProviderSubclass
- TestMalformedReadinessIsResolveGateNotExclusion
- TestEmptyAndDeadlineAreResolveGatesNotExclusion (process-exited-empty / deadline-no-output / deadline-after-output)
- TestProviderFailureIsBlockingDistinctGate
- TestProcessFailureKeepsExplicitExclusion
- TestHostedPONGRealChildFixtures (10 sub-cases incl. never-ending and partial-then-hang kill+reap timing)

Focused runner package (real child processes through RunRoundOne):

- TestSupervisionForAgentBuffersStdoutDisablesSoftGuards (buffered disables both soft guards, keeps heartbeat; non-buffered keeps both)
- TestBufferedQuietLateSuccessNotKilledBySoftGuard (quiet 900ms then exit-0 succeeds)
- TestBufferedNeverEndingChildBoundedByHardTimeout (fails with "timeout", NOT "no_first_output")
- TestNonBufferedSilentChildStillKilledBySoftGuard (contrast: "no_first_output" still fires without the declared buffer)

Commands run (exit 0 each):

- go build ./...
- go vet ./internal/app/ ./internal/runner/
- go test ./internal/app/ (ok, 66s)
- go test ./internal/runner/ (ok, 15s)
- focused -v runs of the tests above, all PASS.

A background `go test ./...` was still running at handoff time; see "Current limitations".

## Current limitations

- No concrete test failures were found; nothing to fix in this invocation.
- Full suite result (completed after the first handoff write): `go test ./...`
  exited 0 — every package `ok`, including internal/app (~103s) and
  internal/runner (~18s) with my changes present. This is additional signal; the
  focused runs above remain the per-slice evidence.
- No real hosted-agent PONG was exercised (only fake-child fixtures). Real per-agent
  PONG behavior is the reviewers'/integration's job, not claimed here.
- BuffersStdout is consumed as the pre-existing declared bool; the TUI's tri-state
  heuristic (declared / 30s-0B heuristic) is unrelated and untouched.
- This slice alone does not satisfy AC-L1's "unchanged quorum and cleanup are
  verified" against the whole tree; that is integration + independent review, not
  self-verified here.

## Exact integration requirements

1. probeFunc seam signature changed (app-package internal). No cross-package consumer
   exists other than the app test fakes already updated. Integration must recompile
   internal/app; nothing else imports readinessObservation or probeFunc.

2. runner.CommandFor call site (internal/app/preflight.go:847) is UNCHANGED and still
   uses the current `(*exec.Cmd, func(), error)` interface. Codex is adding launch
   telemetry and its CommandFor may become a tracked wrapper; I deliberately did NOT
   edit or accommodate that here. Codex applies the small preflight call-site glue
   after my commit/handover. My hostedPONG only uses cmd.Dir / procctl.SetNewProcessGroup /
   cmd.Stdout / cmd.Stderr / cmd.Stdin / cmd.Start / cmd.Wait — all standard *exec.Cmd
   surface, so a tracked wrapper that still returns *exec.Cmd should slot in with only
   the call-site change Codex owns.

3. agents.Spec.BuffersStdout is pre-existing (discover.go:60), mapped from TOML
   `buffers_stdout` (config/runtime.go:897). I consume it; I did not add it. No config
   or schema change is required from my slice.

4. Non-owner review required: a non-hermes-1 participant must review this diff against
   AC-L1 (quiet late success vs hard timeout vs partial output vs malformed/echo/empty
   vs auth/provider/process errors; cleanup; unchanged quorum). I have not self-issued
   any acceptance.

## Evidence provenance

All statements above are from executed commands in this worktree on
feature/meta-protocol-change-evidence-first-efficiency/hermes-1 at commit
4b95baffc21baae1a071473d88160aa3cfce0e8b. Test names/outcomes are copied from real
`go test -v` output; no test, source check, or measurement is invented.
