---
agent: codex
idea: tui-durable-kill
round: 1
date: 2026-06-06
---

## Summary

The durable kill design should be implemented as process-control infrastructure, not as a TUI-only patch. The current root cause is split across three layers: `internal/runner/runner.go:704-729` creates a context-bound direct child with no process group; `execAgentProcess` at `internal/runner/runner.go:677-701` hides `Start`/`Wait`, so it cannot persist the PID; and the resumed/opened TUI paths at `internal/app/app.go:1006-1013` and `internal/tui/live.go:1504-1519` have no kill seam. A two-day "running" badge is almost certainly stale because `agents.DefaultTimeoutMS` is 30 minutes (`internal/agents/discover.go:55-59`), but a group-kill design is still needed for orphaned descendants and for real in-flight resumed runs.

My recommendation: add a small `internal/procctl` package with build-tagged spawn metadata, liveness attribution, and tree kill. Runner uses it for headless attempts and preferably ACP too. Every attempt writes `runs/<runID>/agents/<agentID>/proc.json` plus `pid`/`pgid`/identity fields in `agent.started`. Live `KillAgent` first cancels in memory, then uses the same proc file fallback. Resume/open wire a reattached kill seam that needs only `root` and `runID`. Stale badge cleanup should be an explicit synthetic terminal event after a liveness/identity check proves the recorded process is gone or unactionable.

## Proposed approach (concrete signatures/seams/build-tags)

1. Process-group spawn

Use one helper shared by headless and ACP:

```go
// internal/procctl/spawn.go
type Spawned struct {
    PID int `json:"pid"`
    PGID int `json:"pgid,omitempty"`
    BootID string `json:"boot_id,omitempty"`
    StartTime string `json:"start_time,omitempty"`
    Command string `json:"command"`
    Args []string `json:"args,omitempty"`
    CWD string `json:"cwd,omitempty"`
    Marker string `json:"marker"`
    StartedAt time.Time `json:"started_at"`
}

func SetNewProcessGroup(cmd *exec.Cmd)
func Capture(cmd *exec.Cmd, marker string) (Spawned, error)
func Alive(sp Spawned) (bool, error)
func KillTree(ctx context.Context, sp Spawned) error
func CurrentBootID() (string, error)
```

Build tags:

- `internal/procctl/sysproc_unix.go`: set `cmd.SysProcAttr.Setsid = true`, matching `internal/acp/sysproc_unix.go:10-18`. I prefer `Setsid` over `Setpgid` for this repo because ACP already uses it, it guarantees a new session and group where `pgid == pid`, and it avoids terminal job-control inheritance. `Setpgid: true, Pgid: 0` is also viable, but mixing mechanisms increases edge cases.
- `internal/procctl/sysproc_windows.go`: set `cmd.SysProcAttr.CreationFlags |= windows.CREATE_NEW_PROCESS_GROUP` if available via `syscall` constants, or define the constant locally as `0x00000200` to avoid a dependency. Windows still needs `taskkill /T` for descendants; CTRL events are not reliable for detached CLI tools.
- `internal/procctl/killtree_unix.go`: verify attribution, send `SIGTERM` to `-pgid`, wait a short grace period, then `SIGKILL` to `-pgid` if still alive. If `pgid == 0`, refuse group kill and fall back only if attribution says direct-child-only is safe.
- `internal/procctl/killtree_windows.go`: verify attribution, then run `taskkill /T /F /PID <pid>` with a bounded context. Without third-party deps this is the pragmatic tree-kill path.

Change `execAgentProcess` so it can persist metadata between `Start` and `Wait`:

```go
type ExecAttemptOptions struct {
    Root string
    RunDir string
    Agent agents.Discovery
    SegmentID string
    Kind string // "round" | "steer"
    SteerID string
    Prompt string
    StdoutPath string
    StderrPath string
    Store store.Store
}

func execAgentProcess(ctx context.Context, opts ExecAttemptOptions) (procctl.Spawned, error)
```

Implementation: build `cmd` via `CommandFor`, call `procctl.SetNewProcessGroup(cmd)` before `Start`, wire stdio as today, call `cmd.Start`, call `procctl.Capture`, write `proc.json`, emit or enrich `agent.started`, then `cmd.Wait`. Do not rely on `exec.CommandContext` alone for timeout cleanup; with a process group, context cancellation must trigger `procctl.KillTree` instead of only `cmd.Process.Kill`. The clean way is to use `exec.Command` plus a goroutine/select: `Wait` wins normally; `ctx.Done()` calls `KillTree`, then waits.

ACP should be brought along or explicitly scoped as phase 2. `acp.Spawn` already calls `cmd.Start` at `internal/acp/spawn.go:65-70` and exposes `PID` at `internal/acp/spawn.go:96-102`, so the metadata hook is easy. But `Process.Stop` and `Kill` at `internal/acp/spawn.go:104-139` are direct-child-only; they should call the same process-control helper after attribution is captured.

2. PID/PGID persistence and lifecycle

Persist both an event snapshot and a file:

- Event fields on `agent.started`: `pid`, `pgid`, `boot_id`, `proc_start_time`, `proc_marker`, `command`, `launch`, `proc_file`.
- File: `runs/<runID>/agents/<agentID>/proc.json`, with the full `procctl.Spawned` plus `run_id`, `agent`, `segment_id`, `kind`, `steer_id`, `stdout`, `stderr`, and `terminal_event`.

The file is the restart kill source of truth because a restarted TUI has the run dir and agent dirs. The event is for projection/audit and for old tooling that only reads `events.jsonl`. On `agent.finished`, `agent.failed`, `agent.killed`, or `agent.reconciled`, update `proc.json` with `terminal_event`, `completed_at`, and `exit_error`; do not delete it, because deletion removes audit and weakens stale diagnosis. A restarted parley enumerates candidates by reading `runs/<id>/agents/*/proc.json`, filtering files whose segment/kind has no terminal event or whose projected state is still `running`.

3. Durable cross-restart KillAgent and seams

Change runner kill to one durable function:

```go
type KillOptions struct {
    Root string
    RunID string
    RunDir string
    Store store.Store
    Handle *Handle // optional live handle
}

func KillAgent(ctx context.Context, opts KillOptions, agentID string) (KillResult, error)
func ReattachKillAgent(root, runID string) tui.KillAgentFunc
```

`Handle.KillAgent` can delegate to this with `Handle: h`. Algorithm:

- If a live `active[agentID]` exists (`internal/runner/steer.go:85-113`), mark killed, call the attempt cancel, and let the context watcher group-kill. This preserves the current race behavior and avoids duplicate events.
- If no in-memory attempt exists, load `proc.json`.
- Verify it is not terminal and matches the requested run/agent.
- Run attribution: same boot ID, same process start time, marker/cmdline/cwd match.
- If alive and attributed, call `procctl.KillTree`.
- Emit one `agent.killed` event with `agent`, `segment_id`, `kind`, `steer_id`, `pid`, `pgid`, `source: "reattach"` or `"live"`.
- If not alive but projected as running, emit `agent.reconciled` or `agent.failed` synthetic event, not `agent.killed`.

Wire the seam on resume at `internal/app/app.go:1006-1013`:

```go
KillAgent: runner.ReattachKillAgent(*root, run.RunID),
```

For Home/open, add a seam factory to `LiveOptions` rather than baking runner into `internal/tui`:

```go
type LiveOptions struct {
    ...
    ReattachKill func(runID string) KillAgentFunc
}
```

Then `openRun` (`internal/tui/live.go:1504-1519`) passes `KillAgent: m.opts.ReattachKill(run.RunID)`. The existing `Start` path already returns live seams at `internal/app/app.go:2090-2100`.

4. Stale-running reconcile

Projection alone is not enough because it leaves no audit trail and the next load can show the same stale state. Add a terminal synthetic event:

```json
{"type":"agent.reconciled","data":{"agent":"codex","segment_id":"segment-0001","state":"stale_dead","reason":"recorded process is not alive","pid":123,"pgid":123}}
```

`runstate.ProjectEvents` currently handles `agent.started`, terminal events, and `agent.killed` at `internal/runstate/runstate.go:343-361`, then maps `agent.started` to `StateRunning` at `internal/runstate/runstate.go:449-454`. Add `StateStale` or map reconciled stale to `StateFailed` with `Reason`. I prefer adding `StateStale = "stale"` so the UI can distinguish "the process is gone and the event log was incomplete" from a real agent failure. `deriveLiveness` at `internal/runstate/runstate.go:499-506` should become process-aware only when proc files exist: `live`, `stale`, `idle`, `unverified`.

Reconcile should run on attach/list for projected `running` agents only, and it should be conservative: write `agent.reconciled` only when proc identity is valid for this boot and `Alive` is false. If the proc file is missing or attribution cannot be checked, show `running/unverified` or `stale/unverified` in projection, but do not write a terminal event automatically.

5. PID-reuse safety

This is the main correctness trap. `Signal(0)` only says "some process exists"; it says nothing about identity. The durable kill path must refuse if it cannot prove identity.

Record:

- boot ID: Linux `/proc/sys/kernel/random/boot_id`; Darwin `sysctl kern.boottime` normalized to seconds/nanos; Windows can use a boot-time value derived from system uptime or WMI is unavailable without deps, so record a conservative boot/session token and refuse if it cannot be verified.
- process start time: Linux `/proc/<pid>/stat` field 22; Darwin `ps -o lstart= -p <pid>` or `kinfo_proc` is better but harder without deps; Windows can use PowerShell/CIM only if acceptable, otherwise treat as unverifiable and refuse durable kill.
- cmdline/cwd attribution: store resolved `agent.Path`, sanitized args shape, cwd root, and an environment marker like `PARLEY_RUN_ID`, `PARLEY_AGENT_ID`, `PARLEY_SEGMENT_ID`, `PARLEY_PROC_MARKER`. On Linux verify `/proc/<pid>/cmdline` and `/proc/<pid>/environ` when readable. If environment is unreadable but start time + boot ID + cwd + command all match, allow; if any critical field is missing, refuse.

Across reboot: refuse. Across missing start-time support: refuse durable kill, but allow live in-memory cancellation because that handle owns the child it started. Never send `SIGKILL -pgid` when `pgid` cannot be attributed; a reused process group is worse than a stale badge.

6. Steer on resumed runs and visibility

Do not try to make resumed steer executable in the first durable-kill implementation. A real `SubmitSteer` requires rediscovering agents, reconstructing `runner.Options` including idea status, task, participants, timeout, store, segment semantics, and then safely queuing against possibly live proc files. That is feasible, but it is not required to kill stale agents and it creates new concurrency/race surface.

Recommended scope:

- Resume/open get durable `KillAgent`.
- Resume/open steer remains record-only, but the UI text should be explicit. The fallback message at `internal/tui/live.go:1550-1562` is already close; it should say "recorded only; no live agent process is attached to execute a reply" rather than implying the user merely opened the wrong view.
- Persist steer requests as today, and keep the live-run reply panel behavior for `SubmitSteer` at `internal/tui/live.go:1522-1548`.
- Add diagnostics for `steer.reply_failed` where `stdout`/`stderr` paths are surfaced. For agy-empty, an empty stdout should be a visible `steer.reply_failed: no reply produced`, not a silent panel.

After durable kill is stable, a `runner.ReattachSteer(root, runID, discovered)` seam can be considered. It should only execute if the run has no active round attempt for that agent or if the proc file proves the old attempt is dead and reconciled first.

7. Cross-platform and testability

Tests should avoid killing real unrelated processes and should isolate build-tagged behavior:

- Unit-test `procctl` with a fake `Probe` interface for `Alive`, `BootID`, `ProcessStart`, `Cmdline`, `Environ`, and `KillGroup`. This covers PID reuse guards without OS-specific flakiness.
- Unix integration test: spawn a tiny shell script that starts a background child (`sleep 60`) and then sleeps. Verify `KillTree` on the parent's PGID kills both. Mark with `//go:build !windows`.
- Windows integration test: spawn `cmd /c start /b ...` or a small Go helper process and verify `taskkill /T /F /PID` removes the child. Mark with `//go:build windows`; skip if `taskkill` is unavailable.
- Runner test: fake agent command writes its PID/PGID and sleeps; assert `proc.json` exists after `agent.started`, timeout cancellation calls group kill, and terminal event updates proc lifecycle.
- Reattach test: create a temp run dir with `events.jsonl` and `proc.json`; fake liveness says alive and attributed; call `ReattachKillAgent`; assert `agent.killed` is appended. Repeat with boot mismatch/start-time mismatch and assert no kill command and a refusal error.
- Runstate test: `agent.started` with no terminal event projects running; adding `agent.reconciled` projects stale/terminal and prevents the infinite running badge.

## Concerns / open questions

The biggest design choice is whether ACP joins v1. I think it should at least capture metadata and use group kill in `Process.Stop`, because ACP already creates a separate session on Unix but currently kills only the direct child after timeout (`internal/acp/spawn.go:116-122`). If schedule forces a cut, document that durable kill initially covers headless launch mode only and do not pretend ACP is safe.

Darwin and Windows PID attribution are weaker without dependencies. The safe default is refusal, not best-effort kill. This may disappoint users on those platforms, but it preserves the non-negotiable guarantee that Parley must not kill an unrelated reused PID.

The TUI should not auto-clear every old running badge on list. Automatic reconcile is acceptable only when proc identity is valid and dead. Missing proc metadata from old runs should stay "unverified/stale" with an explicit clear action, because there is no safe process attribution.

`exec.CommandContext` is subtly dangerous here. If it remains in use, Go may kill the direct child before the group-kill cleanup runs. The implementation should own cancellation explicitly with `exec.Command`, `Start`, `Wait`, and a context watcher that calls `KillTree`.

## Risks

PID reuse and process-group reuse are the catastrophic risks. A stale `proc.json` with only `pid`/`pgid` is not enough; boot ID and start time must be mandatory for durable kill, and missing attribution must produce a refusal.

Race with normal exit is likely. The kill path must append only one terminal event, and `agent.finished`/`agent.failed` should not overwrite sticky killed state, following the current sticky behavior in `runstate.applyAgentEvent` at `internal/runstate/runstate.go:458-470`.

Grandchildren can escape if an agent deliberately calls `setsid` or uses platform-specific detachment. Process-group kill is the right default and fixes common CLI child trees, but it is not a sandbox or job object. On Windows, `taskkill /T` is best-effort without a dependency; if future requirements demand stronger guarantees, job objects would be the next step.

Event/file divergence is another trap. `proc.json` and events must be written in a deterministic order: start process, capture metadata, write proc file, append `agent.started`; on terminal, append terminal event and update proc file. If terminal file update fails, the event log remains authoritative for projection, while the stale proc file can be ignored because it references a terminal segment.
