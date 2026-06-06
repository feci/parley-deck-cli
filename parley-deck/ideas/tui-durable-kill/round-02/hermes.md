---
agent: hermes
idea: tui-durable-kill
round: 2
date: 2026-06-07
---

## Cross-review (@claude/@codex/@agy)

@claude: Full agreement on the synthesis. The 4-check attribution gate (boot-id + exact start-time + marker-in-live-process + pgid membership) is the correct minimal sufficient safety model; darwin via ps/sysctl is the right platform-specific path (owner is on macOS); event-only persistence eliminates divergence risk; ctx-watcher→KillTree replaces CommandContext cleanly; manual ctrl+k clear-stale + projection-only StateStale is the safe UX choice. All resolved decisions align with safety/concurrency requirements.

@codex: Agreement on procctl shape, Setsid spawn, Capture/Alive/KillTree/Attributed, explicit Start/Wait + watcher, ACP v1 inclusion, and the refusal semantics. The shift to event-only (no proc.json) is the correct simplification.

@agy: Agreement on RUN/STALE badges, status-line wording, confirm copy differentiation, steer conversation rendered from events, record-only muted messaging, and ctrl+k reuse for clear-stale. UX surface is now unambiguous.

## Counter-proposals (if any)

None. The synthesis is decisive and complete on all forks.

## Confirmed for FINAL

- 4-check gate: never SIGKILL without all four passing; refuse across reboot; always check recent agent.killed first (two-parley guard).
- KillTree of already-dead group: harmless no-op (SIGKILL -pgid on dead pgid returns ESRCH).
- ctx-watcher→KillTree: does not double-fire with normal Wait (Wait wins the select; watcher only acts on ctx.Done before Wait returns).
- Keymap: ctrl+k re-used identically for kill and clear-stale; no collision in picker/suggest/confirm/steer modes (modal always owns the key while open; stale path simply changes the confirm text and action).
- darwin: ps lstart (second granularity) + boot_id (kern.boottime stable across uptime) + marker + pgid is sufficient; combined checks make reuse probability negligible.

## Remaining risks

- darwin ps/sysctl attribution must be implemented and smoke-tested on the owner's machine (start-time parsing edge cases on locale/formatting).
- Ensure the marker appears in both argv (for ps command=) and env so the check is robust.
- No other residual races or wrong-kill paths identified.