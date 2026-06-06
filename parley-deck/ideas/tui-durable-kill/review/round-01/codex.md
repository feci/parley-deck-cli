---
agent: codex
idea: tui-durable-kill
phase: review
round: 1
date: 2026-06-07
---

## Summary
REQUEST-CHANGES: the process-group spawn/watcher wiring is mostly sound, but the durable attribution gate can approve a PID with missing persisted identity fields, which violates the headline safety requirement.

## Findings

CRITICAL, `internal/procctl/procctl.go` `Attributed`: missing boot/start/pgid/command fields are treated as optional, so a persisted `agent.started` with only `pid` can pass attribution when that PID is alive. `Capture` can produce exactly that partial identity if `sysctl`, `ps`, `/proc`, or `Getpgid` fails; in this macOS sandbox, `go test ./internal/procctl` failed because `Capture` recorded empty `BootID`, `ProcStart`, and `Command`. On a resumed run, `KillAgentDurable` would then call `KillGroup` for an alive PID without the required boot-id, exact start-time, pgid, and command proof. Fix: make `Attributed` fail closed unless durable kill is supported and all required recorded fields are present (`BootID != ""`, `ProcStart != ""`, `PGID > 0`, `Command != ""`), and unless the live probe can read each corresponding value. `Capture` should either return a complete durable identity or mark it unusable so no `agent.started` can later be kill-attributed from partial data.

MAJOR, `internal/procctl/procctl_darwin.go` `procStart`: Darwin attribution uses `ps -o lstart`, which has only one-second wall-clock resolution. That does make capture and verify byte-for-byte comparable, but it is not by itself a unique process-start identity. A reused PID for the same command/process-group shape inside the same second can satisfy boot, lstart, pgid, and command checks and be killed as unrelated work. The spec's safety bar says "never SIGKILL a reused/unrelated PID"; this probe does not actually prove that on macOS. Fix: use a higher-resolution Darwin process start identity, for example a sysctl/libproc-backed start time with microseconds/nanoseconds, and store that raw value instead of `ps lstart`.

MINOR, `internal/procctl/procctl.go` `commandMatches`: the command check is symmetric prefix matching, while FINAL.md specifies that the live command must have the recorded command prefix. Accepting `strings.HasPrefix(recorded, live)` weakens the check if a captured or tampered recorded command is longer than the live probe output. Fix: make the direction explicit and only accept `live == recorded` or `strings.HasPrefix(live, recorded)` unless there is a documented platform truncation case covered by tests.

MINOR, `internal/runner/durablekill.go` `KillAgentDurable`: attribution refusal immediately appends a synthetic `agent.failed` and clears the stale badge. FINAL.md says attribution refusal should return a non-kill reason, and the TUI should tell the user verification failed and offer clear-stale as the non-signaling action. The current behavior is safe from an OS-signal perspective, but it collapses "verification refused" and "user chose clear stale" into one action. Fix: return a refusal result without writing a terminal event, or add an explicit clear-stale seam/action that writes `agent.failed` only after the UI has moved into the clear-stale confirmation path.

NIT, `internal/runner/runner.go` `runAgent`: `agent.started` uses the `now` captured before directory setup, prompt construction, command start, and `procctl.Capture`. This can make event time precede the actual process start by a noticeable amount under slow setup. Fix: stamp `agent.started` with `time.Now().UTC()` inside `onStarted`.

Correct points: live in-memory kill and the exec watcher call ungated `KillGroup` only for processes this process just spawned and captured, which is the right trust boundary. `execAgentProcess` has a single `cmd.Wait` goroutine and drains it on cancellation, so I do not see a double-Wait regression. The resume/open seams keep `internal/tui` decoupled from `runner`/`app`, and `KillAgentFunc`'s `(string,error)` result is handled coherently. Windows durable kill refuses through `Attributed`, which is the safe cross-restart behavior.

Test note: the existing real-process procctl tests are useful, but they currently do not cover the most important refusal case: a partial persisted identity with a live PID must refuse and must not signal. Add a fake-probe or real-path test for missing boot/start/pgid/command fields, plus a durable-path test proving an attribution refusal does not call `KillGroup`.

## Verdict
REQUEST-CHANGES with blocking item: `Attributed` must fail closed on incomplete persisted identity before this can be accepted.
