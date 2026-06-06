---
agent: hermes
idea: tui-durable-kill
round: 1
date: 2026-06-06
---

## Summary
PID-reuse is the highest-stakes hazard: a persisted PID can be reused by an unrelated process (post-reboot or PID wrap). Minimal sufficient guard is: store boot-id + process start-time (from /proc or equivalent) + cmdline marker containing the run/agent ID; on kill, re-verify all three + process-group membership before any SIGKILL. Refuse the kill (and surface "stale/unverifiable") if any check fails or if boot-id differs. Never kill across reboots. This is decisive and minimal.

Concurrency races (liveness probe vs natural exit, two parleys on same run dir, synthetic terminal event vs real one, SIGTERM-then-SIGKILL timing) are handled by first-writer-wins on the active map under mu (already in steer.go), liveness-gated fallback only, and atomic event appends. Keymap remains clean: no new single-key collisions with ctrl+k (confirm), picker navigation, Enter (steer/answer), or scroll.

## Proposed approach
Minimal safe attribution rule (decisive):
- Persist with every `agent.started`: boot_id (from `uname -r` or syscall), start_time (unix timestamp or /proc/pid/stat field 22), cmdline (full argv with embedded "parley-run:<runid>-agent:<agentid>" marker), pgid.
- On durable KillAgent fallback: (1) boot_id must match current boot, (2) start_time must match exactly, (3) cmdline must still contain the marker, (4) the live process must be in the recorded pgid. Only then `kill(-pgid, SIGKILL)`. If any fails, refuse and emit diagnostic event; never fall back to blind PID kill.
- This is sufficient; anything less (e.g. PID+pgid alone) allows wrong-process kill after reboot or PID reuse. Must NEVER happen: SIGKILL without all four checks passing.

Process-group spawn: reuse/extend acp/sysproc_unix.go pattern (Setpgid + Pgid=0 for new group) into a shared internal/runner/sysproc.go (build-tagged). Same for windows CREATE_NEW_PROCESS_GROUP. execAgentProcess timeout ctx still works because group kill is explicit fallback only.

Persistence: write pid/pgid/boot_id/start_time/cmdline into the `agent.started` event Data (already done for other fields); optional small proc.json is unnecessary duplication. Clear on terminal event. Restarted parley enumerates via store.Load + filter running agents.

Durable kill seam: LiveOptions.KillAgent already injected; on resume/open supply a reattach func that reads the run dir events, probes liveness per the rule above, and does group kill + appends `agent.killed`. Synthetic terminal reconcile only on explicit "clear stale" user action (never automatic).

Key table delta (only additions, no collisions):
- Existing: ctrl+k → confirm-kill modal, y/N in modal, /kill picker, Enter on steer/answer, arrows/scroll.
- New (resume path only): none. Kill affordance re-uses the exact same ctrl+k + modal path; the injected KillAgent func simply becomes the durable version. No new keys.

Concurrency: liveness probe (Signal(0)) is racy with natural exit but harmless (kill of already-dead pgid is noop). Two parleys: first-writer-wins on active map prevents double-cancel; durable path is read-only until explicit kill. SIGTERM (ctx cancel) then SIGKILL only after timeout or explicit user kill.

## Concerns / open questions
- How expensive is boot_id + start_time verification on every TUI tick for very large runs (100+ agents)? Negligible.
- Should the marker be in argv or env? Argv is visible in ps; env slightly harder to spoof but argv is already used by ACP path.

## Risks
- If boot_id retrieval is not portable or fails on some unix, the guard degrades to "refuse all cross-restart kills" (safe but reduces utility). Mitigation: make boot_id optional with explicit "unverifiable" state shown in TUI.
- Double-kill across two parleys: already prevented by first-writer-wins + killed flag in attempt struct. Durable path must also check for recent `agent.killed` event before acting.