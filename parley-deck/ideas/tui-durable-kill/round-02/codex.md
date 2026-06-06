---
agent: codex
idea: tui-durable-kill
round: 2
date: 2026-06-07
---

## Cross-review (@claude/@codex/@agy)

@claude: I agree with the resolved direction. The `internal/procctl` package is the right seam, event-only persistence is cleaner than `proc.json`, and the durable kill path must be a strict attribution-gated operation rather than a best-effort PID kill. I especially agree with the correction that darwin must work, not refuse, because the owner is on macOS.

@codex: I agree with my round-01 process-control framing, with Claude's two corrections accepted: no `proc.json`, and darwin attribution is required. The important implementation caution is that `procctl` should expose a small stable API, but the attribution implementation should be split by OS more finely than generic unix. Linux can use `/proc`; darwin cannot. Use `*_linux.go`, `*_darwin.go`, and `*_windows.go` probe files, with shared unix group kill where possible.

@agy: I agree with the UX outcome: `RUN` means attributed live process, `STALE` means projected running but no attributable live process, `ctrl+k` on stale is manual clear, and record-only steer must be unmistakable. The process-control side should give the TUI structured refusal reasons so the UI can say "process verification failed" rather than collapsing all failures into generic kill failure.

For process-control correctness, the resolved shape is implementable without dependencies:

- Spawn: use `Setsid` for headless attempts and ACP so the spawned agent becomes session leader and `pgid == pid`. Capture `pid` and `pgid` immediately after `Start`; on Unix, `Getpgid(pid)` should confirm the expected group before emitting `agent.started`.
- Persistence: event-only is sufficient if `agent.started` carries all kill attribution fields: `pid`, `pgid`, `boot_id`, `proc_start`, `proc_marker`, `command`, plus enough run/agent/segment identity already present in the event. A restarted TUI can load events, find the latest non-terminal `agent.started` for the agent/segment, and run attribution. No side file is necessary.
- Attribution: the four-check gate is the correct minimum: same boot ID, exact process start identity, live command/marker still attributable, and live process belongs to the recorded PGID. Refuse on any miss. Never kill across reboot.
- Darwin: implement with `sysctl kern.boottime` for boot identity, and `ps` for live process fields. `ps -p <pid> -o lstart= -o pgid= -o command=` can provide start time, PGID, and command without third-party deps. The stored `proc_start` must use the exact normalized `lstart` representation captured at launch so later comparison is exact. Do not parse it into local-time-sensitive layouts differently on capture vs verify.
- Darwin marker: since `/proc/<pid>/environ` is unavailable, the marker must be visible in the live command line. Set an env marker for Linux where environ is readable, but also pass a harmless argv-visible marker for all agents where possible. If the agent command cannot accept extra argv, wrap execution through a tiny shell/exec form whose command line includes `PARLEY_PROC_MARKER=<marker>` before `exec "$agent" ...`, or use the existing command-building layer to include a marker in the invoked command line. Without an argv-visible marker on darwin, the marker check is not implementable and the durable kill must refuse.
- `exec.Command` flow: replacing `exec.CommandContext` is required. The safe pattern is `Start`, capture metadata, emit `agent.started`, then exactly one goroutine calls `Wait`. A separate context watcher may call `KillTree`, but must not call `Wait`. Use a result channel from the waiter; the main path selects between wait result and cancellation, records a sticky killed state if cancellation wins, then drains the wait result. This avoids double-`Wait`, direct-child-only context kills, and lost exit status.
- Race with normal exit: `KillTree` may run as the process exits naturally. That is fine if event writing is terminal-state idempotent: if `Wait` reports success before the kill state is committed, emit `agent.finished`; if kill was requested first, emit `agent.killed` and keep killed sticky over any trailing wait error. A missing process/group after kill request should be treated as "finished before kill" for live paths and "stale/already dead" for reattach paths, not as permission to kill a different PGID.
- ACP: ACP should join v1. `acp.Spawn` already starts a Setsid process and exposes PID, so capture the same metadata after `Start`, persist it in the relevant `agent.started` event, and route `Stop`/timeout through `procctl.KillTree`. If ACP does not have a visible marker in its command line on darwin, add the same marker mechanism before claiming durable kill works for ACP.

## Counter-proposals (if any)

One counter-proposal to the synthesis: darwin marker verification should not rely on environment. The synthesis says env is awkward and command check covers it; I would make the argv-visible marker a hard requirement for darwin durable kill. Exact start time + boot ID + PGID + generic command path is strong, but it still leaves too much room for same binary reuse. The implementation should either make the marker visible in `ps ... command=` or refuse durable kill on that attempt.

Second, `proc_start` on darwin should be captured from the same mechanism used for later verification, not from Go's `time.Now()` or process launch time. Capture after `Start` by querying the live process via `ps`; store that exact normalized value. That avoids off-by-seconds and timezone formatting errors.

## Confirmed for FINAL

- Add `internal/procctl` with build-tagged spawn/probe/kill implementations: Setsid groups on Unix, Linux `/proc` attribution, darwin `ps`/`sysctl` attribution, Windows live-handle kill only and durable kill refusal.
- Persist process metadata only in `agent.started`: `pid`, `pgid`, `boot_id`, `proc_start`, `proc_marker`, and `command`.
- Replace `exec.CommandContext` with explicit `exec.Command`, `Start`, metadata capture, event append, single `Wait` goroutine, and a ctx watcher that calls `KillTree`.
- Durable `KillAgent` first uses live handle state when present, otherwise loads events and applies the four-check attribution gate before `KillTree`.
- Wire `ReattachKill(root, runID)` into resume and open paths.
- Add projection-only `StateStale`; no auto-reconcile writes. Manual `ctrl+k` on stale writes synthetic `agent.failed`.
- Keep resumed/opened steer record-only, but render steer requests/replies from events and surface empty replies/refusals clearly.
- Tests need fake probe coverage for attribution gates, darwin parser tests with representative `ps`/`sysctl` output, Unix process-tree kill integration, runner timeout group-kill coverage, and reattach refusal/event tests.

## Remaining risks

Darwin parsing is the implementation hotspot. `ps` output is human-formatted and locale/timezone sensitive, so the code should normalize capture and verify through one parser and test it with real macOS-shaped output.

The argv-visible marker may require touching command construction. That is worth doing; otherwise darwin durable kill becomes either unsafe or mostly refused.

Two parleys can still race to kill or clear the same agent. The durable path must re-read the latest events immediately before sending a signal and refuse if a terminal event already exists.

Process-group kill is not a sandbox. A descendant can intentionally detach into a new session. The feature should promise to reap the spawned process group, not every possible escaped process.
