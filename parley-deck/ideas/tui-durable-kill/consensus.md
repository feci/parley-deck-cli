---
idea: tui-durable-kill
phase: consensus
drafter: claude
date: 2026-06-07
participants: [claude, codex, agy, hermes]
---

## Consensus

After round-01 (4 independent positions) and round-02 (cross-review), all four
participants agree on the design below. No open disagreements. codex's process-control
framing is the backbone; hermes's safety gate is the law; agy owns the UX. Two
facilitator resolutions are recorded inline (darwin = first-class, not refuse; the marker
question).

### Goal
Make agent kill **durable across a parley restart**: persist each agent's process
identity, spawn agents in their own process group so a kill reaps the whole tree, and let
a restarted/observational TUI kill via the persisted data with strict PID-reuse safety —
plus a truthful per-agent liveness badge (RUN/STALE), manual "clear stale", and clearer
steer visibility. **The owner is on macOS (darwin); darwin must fully work.**

### A. `internal/procctl` (build-tagged process control)
- `Setsid` spawn (session leader, `pgid == pid`) shared by headless + ACP. After `Start`,
  confirm the group via `Getpgid(pid)` before emitting `agent.started`.
- API (stable, small):
  `type Spawned{PID,PGID int; BootID,ProcStart,Command,Marker string; StartedAt time.Time}`;
  `SetNewProcessGroup(cmd)`, `Capture(cmd, marker) Spawned`, `Alive(Spawned) bool`,
  `KillTree(ctx, Spawned) error` (SIGTERM → short grace → SIGKILL on `-pgid`),
  `CurrentBootID() string`, `Attributed(Spawned) (bool, reason string)` (the gate).
- Build-tag split by OS for the probe/attribution (linux `/proc`, darwin `ps`/`sysctl`):
  `procctl_linux.go`, `procctl_darwin.go`, `procctl_windows.go`, plus a unix group-kill
  file. **darwin**: boot id via `sysctl -n kern.boottime`; live fields via
  `ps -p <pid> -o lstart= -o pgid= -o command=`. **`ProcStart` is captured at launch via
  the SAME `ps -o lstart=` mechanism used to verify it** (codex) — never Go `time.Now()` —
  so capture and verify compare byte-identical normalized strings. **windows**: durable
  kill refuses (live-handle kill only); not the owner's platform.

### B. Attribution gate (the safety law — hermes)
A cross-restart `KillTree` proceeds ONLY if ALL pass; otherwise REFUSE with a structured
reason (so the UI can say "process verification failed"):
1. **boot-id** matches the current boot (refuse across reboot — the PID can't be ours);
2. **proc_start** matches exactly (the decisive PID-reuse guard: a reused PID always has a
   later start time);
3. **pgid membership** — the live process is in the recorded PGID;
4. **command match** — live `command` matches the recorded agent command/path.
An env marker (`PARLEY_RUN_ID`/`PARLEY_AGENT_ID`/`PARLEY_PROC_MARKER`) is set on every
spawn as linux-readable defense-in-depth (via `/proc/<pid>/environ`); it is NOT required
on darwin (where environ isn't cross-process readable) because checks 1-4 already make
mis-attribution effectively impossible. **No shell wrapper** is added around agent spawns
(it would not survive `exec`, and complicates stdio on the hot path); the start-time +
pgid + command + boot-id combination is the agreed sufficient gate. Before signaling,
re-read the latest events and refuse if a terminal event already exists (two-parley guard).
SIGKILL of an already-dead group is a harmless no-op (ESRCH).

### C. Persistence (event-only)
`agent.started` carries `pid, pgid, boot_id, proc_start, proc_marker, command`. **No
proc.json** (avoids the file/event divergence trap). A restarted/observational parley
loads `events.jsonl` (`store.Load`), finds the latest non-terminal `agent.started` per
agent/segment, and runs the gate.

### D. Exec core (owns cancellation — codex)
Replace `exec.CommandContext`: `CommandFor` (no ctx) → `SetNewProcessGroup` → `Start` →
`Capture` → emit `agent.started` → exactly ONE goroutine calls `Wait` (result via channel).
A ctx-watcher calls `procctl.KillTree` on `ctx.Done()` but never `Wait`. The main path
selects wait-result vs cancellation; on cancellation it records sticky killed state then
drains the wait. So a timeout/cancel reaps the whole group (fixes orphan-on-timeout), no
double-Wait, no lost exit status. Shared by round + steer. ACP captures the same metadata
after its `Start` and routes `Stop`/timeout through `KillTree`.

### E. Durable KillAgent + reattach seam
`KillAgent`: if a live `active[agentID]` exists → cancel its ctx (current behavior) which
group-kills via the watcher; else load from events → gate (§B) → `KillTree` → emit ONE
`agent.killed`. A PID that is dead-but-projected-running is NOT killed (it's stale).
`ReattachKill(root, runID) KillAgentFunc` needs only the run dir; wired onto **resume**
(app.go ~1006) and **open** (via a `LiveOptions.ReattachKill` factory → `openRun`). The
live `Handle` keeps its in-memory path. Terminal events are idempotent; killed stays
sticky over a trailing `agent.failed` (existing rule).

### F. Liveness display + manual clear (agy)
- New `runstate.StateStale` (projection-only: projected running but no attributable live
  process); `deriveLiveness` becomes process-aware when proc metadata exists
  (`live/stale/idle/unverified`). **No automatic event writes** on attach.
- Badges: `RUN` (cyan/green, attributed live) vs `STALE` (warn). Stale status line:
  `<agent> STALE (process lost) · ctrl+k to clear`.
- `ctrl+k` (existing key + modal, no new keys): active → `kill agent <id> and its process
  tree? (y/N)`; stale → `clear stale running status for <id>? (y/N)` → writes a synthetic
  `agent.failed` ("stale process cleared by user"). On attribution failure, flash
  "Process verification failed (PID reused). Safe to clear stale badge instead." and
  demote `ctrl+k` to clear-stale.
- Kill-in-progress status `Killing agent and all sub-processes…` during the SIGTERM→grace
  window (agy). "Finished before it could be killed" graceful case. A Home-tab banner
  `⚠️ N stale agent process(es) — ctrl+k on their tabs to clear` when any stale (agy,
  recommended). Short-terminal compaction preserved.

### G. Steer visibility (record-only on resumed; conversation persisted)
Resumed/observational steer stays **record-only** with an explicit muted row
`record steer <id> (read-only) › ` + hint "records steer (no live execution)". The steer
**conversation is rendered from events** (`steer.requested`/`steer.reply_started`/
`replied`/`reply_failed`) woven chronologically into the agent transcript, surviving tab
switches and restarts; an empty reply shows `[agent returned an empty reply]`. Full
resumed-run steer execution = deferred follow-up. (Diagnosis of the owner's "no reply":
they opened an observational/old run = record-only, and/or steered agy which can return
empty — both now surfaced clearly.)

### H. Tests (no real unrelated processes killed)
- `procctl` with a fake `Probe` (Alive/BootID/ProcStart/Command/KillGroup) for the gate:
  attributed→kill; boot mismatch / start-time mismatch / pgid mismatch / command mismatch
  → refuse. darwin parser tests with representative `ps`/`sysctl` output strings.
- Unix integration: spawn a parent that backgrounds a child → `KillTree` reaps both
  (`//go:build !windows`).
- runner: fake agent records pid/pgid, `agent.started` carries metadata, ctx timeout →
  group kill, terminal event updates lifecycle, no double-Wait.
- reattach: temp run dir + events; attributed→`agent.killed`; mismatch/boot-diff→refuse,
  no event; terminal-already-present→refuse.
- runstate: `agent.started` with no terminal → running; proc-dead → `StateStale`;
  synthetic `agent.failed` clears it; killed sticky over failed.

### Non-goals
No supervisor/daemon, no auto-restart, no cross-machine kill, no sandbox/job-object
guarantee (we reap the spawned process group, not deliberately-detached escapees), no full
resumed-run steer execution in v1. windows durable kill refuses.

## Signoffs

(each participant appends its own block below — append-only, do not edit others')

### claude — ✅ ACCEPT
2026-06-07. As drafter I accept. This locks the round-02 resolution: procctl (Setsid +
build-tagged darwin `ps`/`sysctl` + linux `/proc`), event-only persistence, the 4-check
attribution gate (boot-id + exact-start-time + pgid + command; no shell wrapper), the
exec-core cancellation rewrite, durable KillAgent + ReattachKill on resume/open,
projection-only StateStale + manual ctrl+k clear, and record-only resumed-steer with the
conversation rendered from events. darwin is first-class. Ready for FINAL + implementation.

### codex — ✅ ACCEPT
I accept the drafted consensus as ready for FINAL + implementation. The decisive
start-time attribution gate, event-only persistence, exec-core cancellation rewrite,
resume/open reattach kill behavior, projection-only stale state, and record-only resumed
steer are coherent and sufficient for the durable kill scope, with darwin first-class.

### agy — ✅ ACCEPT
I accept this consensus, which fully aligns with the UX goals for a clear process-aware liveness display and manual clear action. The integration of `ctrl+k` for both killing active process trees and clearing stale badges is elegant and prevents UI clutter. The fallback feedback on verification failure provides excellent clarity for the user, and darwin support is solid. This is ready for FINAL and implementation.

### hermes — ✅ ACCEPT
I accept the consensus. The 4-check attribution gate (boot-id + exact start-time + pgid + command), event-only persistence via agent.started, exec-core rewrite with KillTree, durable KillAgent/ReattachKill on resume/open, projection-only StateStale + ctrl+k, and record-only resumed steer from events are complete and safe. darwin first-class. Ready for FINAL + implementation.

