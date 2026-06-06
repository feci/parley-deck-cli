---
agent: claude
idea: tui-durable-kill
round: 2
date: 2026-06-07
---

## Cross-review

Strong convergence: a build-tagged process-control layer (Setsid process groups +
persisted PID/PGID + attribution + tree-kill), a durable `KillAgent` that falls back to
persisted data and is wired onto resume/open, a per-agent liveness display, manual
"clear stale", and resumed-steer kept record-only with much clearer messaging. I'm
resolving the few forks and making one **critical correction (darwin must work, not
refuse — the owner is on macOS).**

### @codex
Adopting almost all of your design — it's the backbone:
- **`internal/procctl` package**, `Setsid` (not Setpgid) shared by headless + ACP,
  `Capture`/`Alive`/`KillTree`/`CurrentBootID`. **Adopt.**
- **`exec.CommandContext` is dangerous here** — replace it: the exec core uses
  `exec.Command` + explicit `Start`/`Wait` + a context-watcher goroutine that calls
  `procctl.KillTree` on `ctx.Done()` (so a timeout/cancel reaps the whole group, not just
  the direct child). **Adopt — this also fixes orphan-on-timeout.**
- **Attribution is mandatory before any cross-restart SIGKILL** (boot-id + start-time +
  marker + pgid); refuse if any check fails; never across reboot. **Adopt.**
- **ACP joins v1**: capture metadata + route its stop through `KillTree` (it already
  `Setsid`s and exposes the PID — low effort). **Adopt.**
- **Two CORRECTIONS to your proposal:**
  1. **darwin/macOS MUST be a first-class attribution target, NOT "refuse".** The owner
     is on macOS — refusing durable kill on darwin defeats the entire feature for them.
     Implement darwin attribution with no deps via shelling out: process start-time +
     command via `ps -p <pid> -o lstart=,command=`, boot id via
     `sysctl -n kern.boottime`. (linux: `/proc/<pid>/stat` field 22 + `/proc/<pid>/cmdline`
     + `/proc/sys/kernel/random/boot_id`. windows: live-handle kill only; durable kill
     refuses — windows isn't the owner's platform.) procctl is build-tagged
     `*_darwin.go` / `*_linux.go` / `*_windows.go` (or `_unix` + a darwin/linux probe split).
  2. **Persistence is EVENT-ONLY (no proc.json).** You proposed both; hermes is right that
     a second file is duplication and you yourself flag event/file divergence as a trap.
     Put `pid,pgid,boot_id,proc_start,proc_marker,command` in the `agent.started` event;
     a restarted/observational parley reads them via `store.Load` (the run dir's
     events.jsonl) → the projection. One source of truth, no divergence.

### @hermes
Adopting your safety model as the gate:
- **Minimal-sufficient attribution = all four checks** (boot-id match, exact start-time,
  marker still in the live process, process is in the recorded pgid) before
  `kill(-pgid, SIGKILL)`; refuse + diagnostic otherwise; SIGKILL never without all four.
  **Adopt verbatim.** Marker: set it in BOTH the env (`PARLEY_RUN_ID`/`PARLEY_AGENT_ID`/
  `PARLEY_PROC_MARKER`) and rely on the recorded argv/command for the check (env is
  awkward to read cross-process on darwin; the `ps … command=` / `/proc/cmdline` check
  covers it).
- **Synthetic terminal event ONLY on explicit user action** (never automatic) — matches
  agy and me; overrides codex's auto-reconcile-on-attach. **Adopt.**
- **No new keys**: durable kill reuses the exact `ctrl+k` + confirm modal; on a stale
  agent `ctrl+k` becomes "clear stale". Durable path must also check for a recent
  `agent.killed` before acting (two-parley guard). **Adopt.**

### @agy
Adopting your UX wholesale:
- Per-agent liveness badges **RUN / STALE / (FIN/ERR/KILL)**; status-line wording for
  stale (`<agent> STALE (process lost) · ctrl+k to clear`). **Adopt** (new
  `runstate.StateStale`, projection-only for display).
- Confirm copy that differs for active vs stale; success/already-dead/can't-attribute
  messages; "finished before it could be killed" graceful case; PID-reuse refusal demotes
  `ctrl+k` to "clear stale badge". **Adopt.**
- **Steer conversation rendered from events into the agent transcript** (steer query block
  + reply streamed below), surviving tab switches and restarts; empty reply →
  visible `[agent returned an empty reply]` / `steer.reply_failed`. **Adopt.**
- Record-only (observational/resume) steer row in muted style + explicit
  "records steer (no live execution)" hint. **Adopt.**
- Short-terminal compaction + a "N stale agents" banner. **Adopt** (banner is nice-to-have).

## Resolved decisions (for FINAL)
1. **procctl** (build-tagged): `Setsid` spawn; `Capture(cmd,marker)→Spawned{pid,pgid,
   boot_id,start_time,command,marker,started_at}`; `Alive`; `KillTree` (SIGTERM→grace→
   SIGKILL on `-pgid`); `CurrentBootID`; `Attributed(spawned)` (the 4-check gate). darwin
   via `ps`/`sysctl`, linux via `/proc`, windows = refuse durable kill.
2. **Persistence: event-only.** `agent.started` carries pid/pgid/boot_id/proc_start/marker/
   command. No proc.json.
3. **Exec core** owns cancellation: `exec.Command`+`Start`+`Capture`+emit `agent.started`+
   `Wait`, with a ctx-watcher → `KillTree`. Shared by round + steer; ACP captures metadata
   + uses `KillTree` in its stop path.
4. **Durable KillAgent**: live `active` cancel + group-kill if in-memory; else load from
   events → `Attributed` gate → `KillTree` → one `agent.killed`. Dead-but-projected-running
   → no kill (it's stale). `ReattachKill(root,runID)` seam wired to resume (app.go ~1006)
   and open (via `LiveOptions.ReattachKill` factory → `openRun`); live handle keeps its path.
5. **Liveness display**: projection-aware `StateStale` (proc dead but no terminal event) +
   liveness `live/stale/idle/unverified`; **no automatic event writes**. `ctrl+k` on a
   stale agent → confirm → synthetic `agent.failed` ("stale process cleared by user").
6. **Attribution gate** = hermes's 4 checks; refuse + diagnostic on any miss; never across
   reboot; check for a recent `agent.killed` first.
7. **Steer**: resumed/observational = record-only with explicit muted messaging; conversation
   rendered from events into the agent transcript; empty reply surfaced. Full resumed-run
   steer execution = deferred follow-up.
8. **Tests**: procctl with a fake `Probe` (Alive/BootID/Start/Cmdline/KillGroup) for the
   reuse guards; a unix integration test (spawn a parent that backgrounds a child →
   `KillTree` reaps both); reattach test (attributed→kill+event, mismatch/boot-diff→refuse,
   no event); runstate stale-projection test; runner test (proc metadata in agent.started,
   timeout→group kill).

## Risks (carried to FINAL)
PID/PGID reuse = catastrophic → the 4-check gate, refuse-on-doubt, group scoping. darwin
attribution via `ps`/`sysctl` must be correct (owner's platform) — test it. `exec.Command`
+ watcher must not regress the round path's stdout/error timing. Two parleys on one run dir
→ check recent `agent.killed` + first-writer-wins. Killed state stays sticky over a trailing
`agent.failed` (existing rule).
