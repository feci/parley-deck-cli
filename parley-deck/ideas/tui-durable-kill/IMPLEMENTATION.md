---
idea: tui-durable-kill
phase: implementation
status: complete
implementer: claude
date: 2026-06-07
---

# IMPLEMENTATION — durable cross-restart agent kill (+ liveness, steer messaging)

Implements FINAL.md. New package `internal/procctl`; changes to `internal/{runner,app,tui}`.
No new deps; build-tagged for darwin/linux/windows; normal round path + `--no-tui` intact.
`internal/tui` imports neither runner nor app (only func seams).

## internal/procctl (new, build-tagged)
- `procctl.go` (portable): `Spawned{PID,PGID,BootID,ProcStart,Command,Marker,StartedAt}`,
  a swappable `probe` interface, `Capture`, `Attributed` (the 5-check gate: supported +
  boot-id + alive + exact start-time + pgid + command), `KillTreeAttributed`, `Alive`,
  `CurrentBootID`, `MarkerEnv`, `commandMatches`.
- `procctl_unix.go` (`!windows`): `SetNewProcessGroup` → `Setsid`; `KillGroup` →
  `SIGTERM -pgid` → 1.5s grace → `SIGKILL -pgid` (ESRCH = no-op); `unixAlive`/`unixPgid`.
- `procctl_darwin.go`: boot id `sysctl -n kern.boottime`; `procStart`/`command`/pgid via
  `ps -p <pid> -o lstart=|pgid=|command=` (start-time captured + verified by the same call).
- `procctl_linux.go`: `/proc` boot_id, `/proc/<pid>/stat` field 22 start-time, cmdline.
- `procctl_windows.go`: durable kill unsupported (refuse); `KillGroup` best-effort `taskkill /T /F`.

## Runner
- Exec core rewrite (`execAgentProcess`): `buildAgentInvocation` (extracted from
  `CommandFor`, which is unchanged for its other callers) → `exec.Command` →
  `SetNewProcessGroup` → marker env → `Start` → `procctl.Capture` → `onStarted(sp)` →
  one `Wait` goroutine; a ctx watcher calls `procctl.KillGroup` on cancel (reaps the whole
  tree — fixes orphan-on-timeout). Returns `(Spawned, error)`. Shared by round + steer.
- `runAgent`: emits `agent.started` enriched with `pid/pgid/boot_id/proc_start/proc_marker/
  command` via `onStarted`.
- `durablekill.go`: `KillAgentDurable(store, agentID)` — reads the latest non-terminal
  `agent.started`, refuses if a terminal already follows (two-parley/idempotency), else:
  live+attributed → `KillGroup` + `agent.killed`; dead or unattributable (reused pid) →
  clear the stale badge via a synthetic `agent.failed` (never signals the reused pid).
  `AgentLiveness(store, id)` → "live"/"stale"/""; `DurableKillAt`/`AgentLivenessAt`
  (run-dir convenience). `Handle.KillAgent` falls back to `KillAgentDurable` when the agent
  isn't in its in-memory `active` map.

## App
- `liveSteerKillSeams` now also returns a `LivenessFunc`; `kill` returns `(message, error)`.
- `reattachSeams(runDir)` builds durable kill + liveness from only the run dir.
- `runDirFor(root, runID)` helper. Wired `KillAgent`+`Liveness`+`ReattachKill`+
  `ReattachLiveness` onto ALL `tui.RunLive` sites: `parley run`, Home/`newLaunchFunc`,
  **resume** (~1006, previously had NO kill), and **open/workspace** (~1970).

## TUI
- `KillAgentFunc` now `func(agentID)(string,error)`; new `LivenessFunc` + `LiveOptions`
  `Liveness`/`ReattachKill`/`ReattachLiveness`; `LaunchResult.Liveness`; `activateRun`
  copies `Liveness`; `openRun` wires the opened run's reattach kill + liveness.
- Badge: a projected-running agent whose `Liveness=="stale"` shows **STALE** (else RUN/etc).
- `ctrl+k` works on running AND stale agents and on resumed/opened runs; the confirm copy
  differs (live → "kill … and its process tree?"; stale → "clear stale running status?");
  the outcome message (killed tree / cleared stale / verification-failed) is shown.
- Steer record-only messaging unchanged from 1.18 (observational runs already say so).

## Tests
- `procctl` (real processes on this host): self-process attributed; each tampered identity
  field (start-time/boot/pgid/command) refused; `KillGroup` reaps a backgrounded grandchild;
  dead pid never attributed.
- `runner`: `KillAgentDurable` clears a stale dead pid (synthetic `agent.failed`); no-op
  after a terminal event; `AgentLiveness` stale/absent; **real end-to-end** —
  `TestDurableKillEndToEndRealProcess` spawns `/bin/sleep` in its own group, persists
  `agent.started`, and durably kills it via the event log alone.
- `tui`: STALE badge via the liveness seam + ctrl+k "clear stale" confirm + outcome msg;
  existing kill/steer/suggest tests updated to the new seam signature.

## Verification
`go build ./...`, `go vet ./...`, `go test ./...`, `go test -race ./internal/runner
./internal/procctl` all green; cross-compiles for windows + linux.

## Real macOS smoke
`TestDurableKillEndToEndRealProcess` and the `procctl` integration tests ARE the real-OS
verification on this macOS host (real process-group spawn, grandchild reap, attribution
refusal, durable kill via the event log). The full interactive flow for the owner:
`parley run` → ctrl+k an agent (kills its tree); kill parley; `parley resume <run>` (or
`/open` in `parley tui`) → ctrl+k again → a genuinely-running agent's tree is reaped, a
stale "running" badge clears; `ps` shows no leftover.

## Fix-up cycle 1 (Phase 6 review → addressing codex + agy)
- **CRITICAL — `Attributed` fail-closed** (codex+agy): now refuses unless every facet is
  recorded AND readable from the live process AND matches (boot-id, exact start-time, pgid,
  command) + the process is its own session leader (`pgid==pid`, our Setsid signature).
  Tests `TestAttributedFailsClosedOnMissingFields`,
  `TestDurableKillRefusesUnattributableLiveProcess` (a live unattributable pid is NEVER
  signaled; no event written).
- **CRITICAL — ACP integration** (agy): `runACPAgent` now captures the durable identity
  (`procctl.CaptureByPID`) into the ACP `agent.started`, and `acp.Process.Stop/Kill` route
  through `killProcessGroup` (reaps the tree, not just the child).
- **MAJOR — darwin precise start time** (codex+agy): `procStart` now uses
  `unix.SysctlKinfoProc("kern.proc.pid")` p_starttime at **microsecond** resolution (no new
  dep — `golang.org/x/sys` was already in the module graph), eliminating the 1-second
  `ps lstart` reuse window.
- **MAJOR — `KillGroup` self-protection** (agy): refuses to signal parley's own process
  group (guards the Setsid-skipped edge).
- **MAJOR — Home stale banner** (agy): `renderHome` shows `⚠ N stale agent process(es)…`.
- **MINOR — `commandMatches`** (codex+agy): only `live==recorded || HasPrefix(live,recorded)`.
- **MINOR — refusal handling** (codex+agy): an alive-but-unattributable agent is no longer
  auto-cleared (could be our own live process the probe couldn't read); durable kill returns
  `Failed`, and the seams return a Go error so the TUI shows it red. clear-stale now only on a
  provably-dead pid.
- **NIT** (codex+agy): `agent.started` stamped at actual process start.
- **NIT** (agy): steer attempts don't persist identity (transient) — documented below.

`go build/vet/test ./...` + `-race` + windows/linux cross-compile all green. Ready for re-review.

## Deviations / scope (per FINAL)
- Liveness is a TUI-side decoration via the seam (kept `runstate.ProjectEvents` pure — no
  `StateStale` written; the STALE badge is computed from the liveness seam over a
  projected-running agent). Stale clear writes a synthetic `agent.failed` (manual, on
  ctrl+k; no auto-reconcile).
- windows durable kill refuses (live-handle kill only). Resumed-run steer stays record-only
  (full resumed steer execution + the larger steer-conversation-from-events rendering remain
  deferred follow-ups; the durable-kill ask is fully delivered).
