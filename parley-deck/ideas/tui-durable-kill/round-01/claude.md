---
agent: claude
idea: tui-durable-kill
round: 1
date: 2026-06-06
---

## Summary

Two root causes, one dominant. The "2-day running agent" is almost certainly a **stale
badge**: the 30-min timeout means the real process died long ago; parley just never wrote
a terminal event (it was killed mid-run), so the projection shows running forever. So the
fix is two-pronged: (a) **truth** — a per-agent process-liveness probe so the TUI shows
running-vs-stale-vs-dead and can clear/kill a stale badge; (b) **power** — persist PID+PGID,
spawn agents in their own process group, and make `KillAgent` work across a restart by
falling back to the persisted PID and group-killing the tree (reaping orphaned
grandchildren). Plus wire the kill seam onto resume/open (today they get none).

## Proposed approach

### 1. Spawn each agent in its own process group
In `CommandFor` (headless), set `SysProcAttr` like the ACP path already does. Unix:
`Setpgid: true, Pgid: 0` (new group = child's PID). Build-tagged so it compiles on
windows (`CREATE_NEW_PROCESS_GROUP`). One shared helper `setProcGroup(cmd)` reused by ACP
+ headless. The existing per-agent timeout-context cancel still kills the direct child;
group-kill is the new escalation that also reaps grandchildren.

### 2. Persist PID + PGID per agent
After `cmd.Start()` (need to split Start/Wait in `execAgentProcess`, or read
`cmd.Process.Pid` right after start), write a durable record:
- extend `agent.started` event data with `pid` and `pgid`, AND
- write `runs/<id>/agents/<a>/proc.json` = `{pid, pgid, started_at, boot_id, cmd}`.
Clear/supersede on the terminal event (or leave it and rely on liveness). The proc.json
is what a restarted parley reads with only the run dir.

### 3. Durable, cross-restart KillAgent
`KillAgent(agentID)`:
1. in-memory attempt present → cancel its ctx (current behavior) AND group-kill its PID.
2. else load `proc.json`; liveness-check (`Signal(0)`); if alive AND attributable
   (PID-reuse guard, see §5) → `killtree(pgid)` (SIGTERM then SIGKILL); emit `agent.killed`.
3. if the PID is dead → it was a stale badge; write a synthetic terminal event
   (`agent.failed` reason "process not running / stale") so the badge clears.

### 4. Reattach kill seam on resume/open
A new app seam `reattachKillFunc(root, runID)` that needs only the run dir (no live
`Handle`): it reads `proc.json`, probes liveness, group-kills. Pass it as
`LiveOptions.KillAgent` on the resume (≈app.go:1006) and open/workspace (≈1970) `RunLive`
calls. Now a restarted/observational TUI can kill. (A live `Handle` keeps using its
in-memory path.)

### 5. PID-reuse safety (critical)
Never SIGKILL a PID we can't attribute. Guard with: (a) a **boot id** stored in
proc.json (`/proc/sys/kernel/random/boot_id` on linux, `kern.boottime`/sysctl on darwin,
or a parley-run nonce) — if the boot id differs, the PID is from a prior boot → refuse
(the process is gone anyway). (b) Best-effort cmdline check (the live process's argv
should contain the agent binary). At minimum, refuse-across-reboot is enough since a
post-reboot PID can't be our agent. Group-kill targets the recorded PGID, which further
scopes the blast radius to our spawned tree.

### 6. Per-agent liveness display + stale reconcile
On attach (and each tick) probe `proc.json` PIDs. Add a per-agent liveness to the
projection/TUI: `running` (alive), `stale` (badge says running but PID dead/absent),
`dead`. Show it in the badge (e.g. `RUN` vs `STALE`). A reconcile pass on attach can emit
the synthetic terminal event for stale agents so the workspace is self-healing. ctrl+k on
a stale agent just clears it (no process to kill).

### 7. Steer visibility
- Make record-only unmistakable: if `SubmitSteer` is nil (observational/resume), the steer
  row/hint says "observational run — steer recorded only, reopen the live run to chat".
- Persist the conversation: the steer reply already lives in `steers/<id>/`; show a small
  per-agent "conversation" list (prior steer replies) so it survives tab switches, not
  just the transient panel.
- Diagnose the owner's case: most likely they opened the OLD (2-day) run = observational =
  record-only, so no reply. A fresh `parley run` (or Home-N) gives the live seam. Also,
  steering agy can yield an empty reply (agy headless flakiness) → show `steer.reply_failed`
  clearly. Consider: enable a real `SubmitSteer` on resumed runs by building a handle from
  the resolved run + agents discovery (heavier — decide scope in consensus).

## Concerns / open questions
- Splitting `cmd.Start()`/`cmd.Wait()` in the shared `execAgentProcess` to capture the PID
  without regressing the round path — confirm with codex.
- boot-id portability (linux vs darwin vs windows) — codex/hermes to nail the per-OS source.
- Is writing a synthetic `agent.failed` to clear a stale badge acceptable, or should stale
  be a pure projection state (no event)? I lean projection-only for read paths + an
  explicit user "clear" action that writes the event.
- Scope of steer-on-resume: real round-trip vs record-only-but-clear. I lean record-only +
  clear messaging for v1; full resumed-run steer is a follow-up.

## Risks
- Killing the wrong process via PID reuse — mitigated by boot-id + group scoping + liveness
  gate (§5); this is the highest-stakes part.
- Group-kill semantics differ unix/windows — strict build-tagged split + fake-process tests.
- Splitting Start/Wait could change error/stdout-capture timing on the hot round path —
  keep the change minimal and covered by the existing round tests.
- Over-eager stale reconcile racing a genuinely-slow-but-alive agent — gate stale on a
  failed liveness probe, never on elapsed time alone.
