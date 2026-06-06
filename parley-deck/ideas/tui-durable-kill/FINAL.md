---
idea: tui-durable-kill
phase: final
status: final
drafter: claude
implementer: claude
date: 2026-06-07
participants: [claude, codex, agy, hermes]
supersedes: consensus.md
---

# FINAL — durable, cross-restart, process-group agent kill (+ liveness, steer visibility)

Ratified by claude, codex, agy, hermes (all ✅ ACCEPT in consensus.md). Build spec for
Phase 5. Owner is on **macOS (darwin)** — darwin is first-class. No new deps; the normal
round path, `--no-tui`, and the driver stay intact.

## 1. `internal/procctl` (new package, build-tagged)
- `procctl.go` (portable): types + the OS-independent orchestration.
  ```go
  type Spawned struct { PID, PGID int; BootID, ProcStart, Command, Marker string; StartedAt time.Time }
  func SetNewProcessGroup(cmd *exec.Cmd)            // Setsid (unix); no-op refuse-marker (windows)
  func Capture(cmd *exec.Cmd, marker string) Spawned // after Start: pid,pgid,bootid,proc_start,command
  func Alive(s Spawned) bool                         // pid alive (Signal 0 / FindProcess)
  func Attributed(s Spawned) (bool, string)          // the 4-check gate; reason on refusal
  func KillTree(s Spawned) error                     // SIGTERM -pgid → 1.5s grace → SIGKILL -pgid
  func CurrentBootID() string
  func MarkerEnv(runID, agentID, marker string) []string // PARLEY_RUN_ID/AGENT_ID/PROC_MARKER
  ```
- `procctl_linux.go`: bootid `/proc/sys/kernel/random/boot_id`; proc_start `/proc/<pid>/stat`
  field 22 (jiffies; store raw); command `/proc/<pid>/cmdline`; pgid `syscall.Getpgid`.
- `procctl_darwin.go`: bootid `sysctl -n kern.boottime`; live fields one call
  `ps -p <pid> -o lstart= -o pgid= -o command=`; **`ProcStart` captured at launch via the
  SAME `ps -o lstart=` call** so capture==verify byte-for-byte.
- `procctl_unix.go` (`//go:build !windows`): `SetNewProcessGroup` → `SysProcAttr{Setsid:true}`;
  `KillTree` via `syscall.Kill(-pgid, SIGTERM/SIGKILL)`; ESRCH = already dead = nil.
- `procctl_windows.go`: `SetNewProcessGroup` no-op; `Attributed` → `(false,"durable kill
  unsupported on windows")`; `KillTree` → `taskkill /T /F /PID` best-effort; `CurrentBootID` "".
- **Attribution gate** (`Attributed`): ALL must pass else refuse — (1) `BootID==CurrentBootID()`
  (refuse across reboot), (2) `Alive`, (3) live `proc_start == s.ProcStart` exactly, (4) live
  `pgid == s.PGID`, (5) live `command` has the recorded `s.Command` prefix. Reason string on
  any miss.
- Fake-probe seam for tests: an unexported `probe` interface (Alive/BootID/ProcStart/
  Command/Pgid/killGroup) with the real OS impl as default and a swappable var for tests.

## 2. Runner exec-core rewrite (`internal/runner/runner.go`, `steer.go`)
- `execAgentProcess` (shared by round + steer): build via `CommandFor` (drop the ctx from
  `CommandFor`; it builds the cmd only), `procctl.SetNewProcessGroup(cmd)`, set env
  `+= MarkerEnv(...)`, wire stdio, `cmd.Start()`, `sp := procctl.Capture(cmd, marker)`,
  return `sp` to the caller so it enriches `agent.started`. Then: ONE goroutine `Wait`s
  (result→chan); a watcher selects `ctx.Done()` → `procctl.KillTree(sp)`. Returns
  `(Spawned, error)` with sticky-killed signalled via the existing tracker.
- `runAgent`: emit `agent.started` enriched with `pid,pgid,boot_id,proc_start,proc_marker,
  command` from the captured `Spawned`; register the attempt (cancel via the watcher path);
  on timeout/cancel the group is reaped.
- `Handle.RunSteerAttempt`/`runSteerAgent`: same exec core (steer attempts also group-killed).
- ACP (`acp.go`/`acp/spawn.go`): capture the same metadata after `Start`, enrich its
  `agent.started`, route `Stop`/timeout through `procctl.KillTree`.

## 3. Durable KillAgent + reattach (`internal/runner/steer.go`, `internal/app/app.go`)
- `runner.KillAgentDurable(store, root, runID, agentID) KillResult`: load events
  (`store.Load`), find the latest non-terminal `agent.started` for the agent; if a terminal
  `agent.killed`/`finished`/`failed` already follows → refuse (two-parley/idempotency);
  build `Spawned` from the event; `Attributed` gate → `KillTree` → append ONE `agent.killed`
  (`source:"reattach"`). If not `Alive` → return `{Killed:false, Reason:"stale"}` (no kill,
  no event — the TUI offers clear-stale). On `Attributed` refusal → `{Killed:false,
  Reason:<reason>}`.
- `Handle.KillAgent`: live `active` path as today (cancel → watcher group-kills); fall back
  to `KillAgentDurable` when not in-memory.
- App: `runner.ReattachKill(root, runID) tui.KillAgentFunc` wrapping `KillAgentDurable`.
  Wire `KillAgent: runner.ReattachKill(*root, run.RunID)` on the **resume** RunLive
  (~app.go:1006). Add `LiveOptions.ReattachKill func(runID string) KillAgentFunc`; the
  **open** path (`internal/tui/live.go openRun`) sets `KillAgent: m.opts.ReattachKill(runID)`;
  app supplies the factory on the workspace/Home RunLive (~1970) and parley-run.

## 4. runstate liveness (`internal/runstate/runstate.go`)
- `StateStale = "stale"`. `AgentState` already carries killed; add proc fields read from
  `agent.started` (`PID,PGID,BootID,ProcStart,Command`) for the gate + display.
- Projection: `agent.started` → running (as today) but also store proc fields. Stale is
  **computed at read time** by the TUI (process-aware liveness), NOT written during
  projection. Add `agent.failed` reason passthrough already exists; killed sticky exists.
- `deriveLiveness`: when proc fields exist and projected running → probe `procctl.Alive`/
  `Attributed`: live→`live`, dead/unattributable→`stale`; else existing `unverified/idle`.
  (TUI calls a helper; ProjectEvents stays pure/process-free for testability — the
  liveness probe is a TUI-side decoration over the projection.)

## 5. TUI (`internal/tui/live.go`, styles in `app.go`)
- Per-agent liveness: compute `live/stale` for running agents via `procctl` (throttled, e.g.
  on tab switch + the existing tick); badge `RUN` (cyan) vs `STALE` (warn); `shortState`
  handles `stateStale`→`STALE`; `stateBadge` warn.
- `ctrl+k` (existing modal): if agent attributed-live → `kill agent <id> and its process
  tree? (y/N)` → `opts.KillAgent`; if stale/unattributable → `clear stale running status
  for <id>? (y/N)` → write synthetic `agent.failed` ("stale process cleared by user") via a
  new `opts.ClearStale func(agentID) error` seam (or fold into KillAgent returning a
  "stale" result that the TUI turns into the clear action). Kill-in-progress status
  `Killing agent and all sub-processes…`. Refusal → flash "process verification failed …
  safe to clear stale" + demote to clear-stale. Finished-before-kill graceful.
- Home banner `⚠️ N stale agent process(es) — ctrl+k on their tabs to clear` when any stale.
- Steer: resumed/observational row muted `record steer <id> (read-only) › ` + hint; render
  the steer conversation from events (`steer.requested/reply_started/replied/reply_failed`)
  into the agent transcript (survives tab switch/restart); empty reply → `[agent returned
  an empty reply]`. (Keep the live reply panel for live runs.)
- No new keys; short-terminal compaction preserved.

## 6. Tests
Per consensus §H: procctl fake-probe gate tests (attributed→kill; boot/start/pgid/command
mismatch→refuse) + darwin parser tests with real `ps`/`sysctl`-shaped strings; unix
KillTree-reaps-tree integration (`!windows`); runner metadata-in-`agent.started` + timeout
group-kill + no-double-Wait; reattach attributed→`agent.killed` / mismatch→refuse /
terminal-present→refuse; runstate stale projection + sticky killed.

## Acceptance
`go build ./... && go vet ./... && go test ./...` green (repo-local caches);
`go test -race ./internal/runner ./internal/procctl` green; **a real macOS smoke**: spawn
an agent, kill it via TUI; restart parley, `parley resume`/open the run, kill again; verify
the process tree is gone (`ps`), a stale badge clears. Manual-smoke note in IMPLEMENTATION.md.
Then Phases 6-8 review to zero agreed fixes.
