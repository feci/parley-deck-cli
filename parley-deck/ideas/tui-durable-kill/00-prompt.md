---
idea: tui-durable-kill
author: user
created: 2026-06-06
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + TUI wiring (liveness display, kill-on-resume seam, confirm)
  codex: process lifecycle — process-group spawn, PID/PGID persistence, cross-restart/group kill, cross-platform (unix/windows), liveness; Go correctness + tests
  agy: UX — stale-vs-running display, kill affordance on resumed runs, steer reply visibility/diagnosis, messaging
  hermes: safety/concurrency + keymap — PID-reuse hazards, killing the wrong process, race vs real exit, no key collisions
transport: local-dir
cross_review_rounds: 1
status: final
---

## Problem / idea (owner's words)

1. **Durable, cross-restart agent kill.** The owner sees an agent that has shown
   "running" for ~2 days and cannot kill it from the TUI. They suspect it may not really
   be running. They want: **Parley should hold a reference to the agents' OS processes and
   be able to kill them even after parley is interrupted and restarted.**
2. **Steer reply still not visible.** When the owner types to a model they still see no
   conversation/reply. They may be in an old session / observational run and will retest
   with a fresh session — but we should diagnose and harden this.

## Current state (VERIFIED against the code — design against these facts)

- **Headless agent timeout default = 30 min** (`agents.DefaultTimeoutMS = 1_800_000`,
  `runner.timeoutForAgent`). So an agent "running for ~2 days" is NOT a live 30-min-bounded
  process — almost certainly a **stale running badge** (parley was killed before writing
  the terminal `agent.finished`/`agent.failed`, so `runstate.ProjectEvents` shows it
  running forever) and/or an **orphaned grandchild** the agent CLI spawned.
- **The kill registry is in-memory only.** `runner.Handle.active map[string]*attempt`
  (cancel funcs) lives in the running process; `KillAgent` cancels a child context. After
  a restart there is no handle → nothing to cancel. **No PID is persisted for headless
  agents.**
- **Headless agents are NOT spawned in their own process group.** `runner.CommandFor`
  builds `exec.CommandContext(ctx, …)` with no `SysProcAttr`. So canceling the context
  kills only the DIRECT child; grandchildren orphan and survive (codex flagged this as a
  v1.18 follow-up). NOTE: the ACP path already sets `Setsid` via
  `internal/acp/sysproc_unix.go` (`setSysProcAttr`), and a PID-liveness probe already
  exists at `internal/driver/proclive_unix.go` (`os.FindProcess` + `Signal(0)`, build
  tagged unix/windows) — reuse these patterns.
- **Resumed / opened runs get NO kill (or steer) seam.** `tui.RunLive` is called with
  `SubmitSteer`/`KillAgent` ONLY for freshly-launched runs (`parley run` ≈app.go:1752 and
  the Home-N `newLaunchFunc` ≈2083). The `resume` path (≈1006, `Resume:true`) and the
  workspace/open path (≈1970) pass NO seams → after a restart you literally cannot kill.
  This is the core reason the owner can't kill the 2-day agent, and likely why a steer on
  an opened old run shows no reply (observational runs are record-only).
- **Liveness is per-run, not per-process.** `runstate.deriveLiveness` returns
  `unverified` for a running+resumed run; there is no per-agent process-aliveness check,
  so the TUI cannot tell a genuinely-running agent from a stale badge.
- Steer round-trip (shipped 1.18.0): `Handle.RunSteerAttempt` spawns a fresh single-agent
  attempt; the reply tails into the agent tab. It only runs when a live `SubmitSteer` seam
  is present (TUI-launched runs) — record-only otherwise. agy's headless CLI is
  intermittently empty, so a steer to agy specifically may produce no reply.

## Proposed direction (a STARTING proposal — challenge it in round-01)

- **Spawn every agent in its own process group** (unix `Setpgid`/`Setsid`; windows
  `CREATE_NEW_PROCESS_GROUP`), in `CommandFor`/`execAgentProcess`, mirroring
  `acp/sysproc_*`. Then killing the group reaps grandchildren. Add build-tagged
  `killtree_unix.go` / `killtree_windows.go` (`syscall.Kill(-pgid, SIGKILL)` /
  `taskkill /T /F /PID`).
- **Persist each agent's PID (+PGID) durably**: carry `pid`/`pgid` in the `agent.started`
  event AND/OR a small `runs/<id>/agents/<a>/proc.json`; clear/mark on the terminal event.
- **Durable KillAgent**: in-memory cancel first; else fall back to the persisted PID/PGID
  → liveness-check (`Signal(0)`) → kill the process group; emit `agent.killed`. This must
  work from a restarted/observational TUI that has only the run dir.
- **Wire the kill seam on resume/open**: build a "reattach kill" seam (needs only the run
  dir + persisted PIDs, not a live `Handle`) and pass it to the resume/open `RunLive`
  calls so the owner can kill after a restart.
- **Per-agent liveness + stale reconcile**: on attach, probe persisted PIDs; show a real
  running/stale/dead state per agent; let the owner clear or kill a stale "running" agent.
- **Steer visibility**: make record-only vs executing unmistakable in the UI; persist +
  show the steer conversation so it survives tab switches; diagnose the owner's "no reply".

## Round-01 focus questions (answer independently)

1. **Process-group spawn.** Exact unix mechanism for headless `CommandFor` (Setpgid with
   Pgid=0 vs Setsid; interaction with the existing per-agent timeout-context kill) and the
   windows equivalent. Reuse/extend `acp/sysproc_*`? One shared helper for ACP + headless?
2. **PID/PGID persistence + lifecycle.** Where to record (event data vs proc file vs
   both), when to write/clear, and how a restarted parley enumerates a run's live agents.
3. **Durable, cross-restart kill.** The KillAgent fallback-to-PID path; liveness-gated
   SIGKILL of the group; the `agent.killed` event; how it reaches the TUI on resume/open
   (the reattach seam shape that needs only the run dir).
4. **Stale-running reconcile.** How the TUI shows running-vs-stale-vs-dead per agent, and
   how to clear a stale "running" badge that has no terminal event (write a synthetic
   terminal event? a reconcile pass on attach?). Keep the normal live path correct.
5. **PID-reuse safety (critical).** A persisted PID may be reused by an unrelated process
   after a reboot. How do we avoid killing the wrong process — record start-time/boot-id,
   verify the cmdline/our marker, scope to our process group, or refuse across reboots?
6. **Steer on resumed runs + visibility.** Can a resumed run get a working `SubmitSteer`
   (it needs agents discovery + a handle)? Worth it, or keep resumed-run steer record-only
   but make that obvious? How to surface/persist the conversation and diagnose the owner's
   "no reply" (observational? agy-empty? panel not noticed?).
7. **Cross-platform + testability.** Build tags; how to test process-group kill + liveness
   + PID-reuse guards headlessly with fake processes (mirror `proclive_*`).

## Constraints (non-negotiable)
- Reuse the existing `proclive_*` liveness probe and `acp/sysproc_*` process-group pattern;
  build-tagged `*_unix.go` / `*_windows.go` for anything syscall-specific (the repo builds
  for both). No new third-party deps.
- `internal/tui` reaches new capability only through injected `LiveOptions` func seams.
- PID-reuse safety is mandatory: never SIGKILL a PID we cannot reasonably attribute to our
  agent. Killing must reap the whole tree (group), not just the direct child.
- Do not break the normal live run path, `--no-tui`, or the driver. Keep the durable steer
  audit trail.
- Testable headlessly (fake processes / fake agent command); build/vet/test green.

## Non-goals
- No long-running supervisor/daemon, no auto-restart/resurrection of agents, no
  cross-machine kill. Not changing what a round does. No mouse.

## Deliverables (in order)
1. FINAL.md: the process-group spawn, PID/PGID persistence, durable cross-restart
   group-kill, the reattach seam + resume/open wiring, stale-running reconcile + liveness
   display, PID-reuse safety, and the steer-visibility handling.
2. Implement in `internal/{runner,app,tui,runstate}` (+ build-tagged killtree files +
   tests), build/vet/test green; a manual-smoke note. Then Phases 6-8 review.
