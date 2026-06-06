---
agent: agy
idea: tui-durable-kill
round: 2
date: 2026-06-07
---

## Cross-review (@claude/@codex/@agy)

### @claude
I fully agree with the round-02 synthesis. You have done an excellent job balancing safety, correctness, and user experience while narrowing the implementation details:
* **Event-only persistence:** Dropping the duplicate `proc.json` file is a huge win for UX simplicity. The TUI doesn't need to coordinate file deletion or run the risk of state divergence on disk. Projecting the metadata directly from the `agent.started` event data in `events.jsonl` is the cleanest and most robust path.
* **Manual clear over auto-reconcile:** I support keeping the state reconciliation manual via `ctrl+k` and projection-only (`StateStale`). Auto-writing failure events on startup is dangerous; high disk I/O, heavy CPU load, or minor delays in starting up could cause the TUI to prematurely declare a slow-starting process dead and permanently mutate the event history with synthetic failure logs. A manual user-driven clear is the safest and cleanest UX.
* **macOS support:** Thank you for correcting the macOS attribution path. The owner is on Darwin, so a refuse-on-macOS default would defeat the entire purpose of this feature for them. Shelling out to `ps` and `sysctl` is an elegant, dependency-free solution.

### @codex
Your process-control architecture is excellent and provides the robust foundation this feature needs:
* **`Setsid` process groups:** Standardizing on `Setsid` processes and replacing `exec.CommandContext` with explicit `Start`/`Wait` + goroutine watcher guarantees we reap the entire tree. This cleanly fixes the primary pain point where background grandchildren (like orphaned compilers or linters) outlive their context cancels.
* **Mac attribution implementation:** Since the owner is on macOS, your initial proposal to refuse durable kills on non-Linux platforms is corrected. Using `/bin/ps -p <pid> -o lstart=,command=` to extract start-time and command line details, combined with `sysctl -n kern.boottime` for the boot ID, provides the necessary telemetry to satisfy our 4-check safety gate without adding external libraries.

### @agy
The round-02 synthesis has adopted the majority of the UX details, ensuring clear feedback for the user:
* **Discoverable states:** The visual distinction between `RUN` and `STALE` tabs, paired with context-sensitive prompts for `ctrl+k` (killing active vs clearing stale), is exactly what is needed to make process control transparent.
* **Durable conversations:** Rendering steer query-reply cycles chronologically in the main transcript buffer directly from the event log (rather than a transient panel) is a massive usability upgrade that survives tab switching and restarts.
* **Record-only indications:** The grayed-out input row and "records steer (no live execution)" hint will prevent users from typing blindly into resumed runs and expecting real-time replies.
* **Gaps addressed:** While most points were adopted, the synthesis glossed over the **kill-in-progress message** and marked the **"N stale agents" banner** as a mere nice-to-have. I will reinstate these in my counter-proposals below to ensure they make it into the final spec.

---

## Counter-proposals (if any)

I propose two minor adjustments to ensure the UX remains completely transparent and discoverable, especially under edge cases:

1. **Reinstate the Kill-in-Progress Status Message**
   * **Proposal:** During the process group termination window (from the moment a user confirms `ctrl+k` and `KillTree` executes `SIGTERM` up to the final `SIGKILL` escalation or confirmation), the TUI status line must display: `[warnStyle] "Killing agent and all sub-processes..."`.
   * **Rationale:** A process group kill (especially with a grace period for SIGTERM) can take a second or more to fully execute and verify. Without immediate, transient status feedback, the user might assume their keystroke was ignored and hit `ctrl+k` repeatedly, or think the TUI has frozen.

2. **Elevate the "Many Stale Agents" Banner to Recommended**
   * **Proposal:** If one or more stale agents are projected, display a persistent warning banner on the Home/dashboard view: `[warnStyle] "⚠️ N stale agent process(es) detected. Press ctrl+k on their tabs to clear."`
   * **Rationale:** Individual tab badges are excellent, but if a run has many agents and the tab bar overflows, the stale badges are pushed off-screen. A global indicator on the home tab ensures the user is aware of dead background processes and knows how to resolve them without scrolling through every tab.

---

## Confirmed for FINAL

* **Liveness Badges:** `[RUN]` (cyan/green) and `[STALE]` (yellow/warn) on tabs. completed/failed tabs show muted standard badges.
* **Status-Line Info:** When active on a stale tab, the status line displays: `<agent> STALE (process lost) · ctrl+k to clear`.
* **Differentiated Confirm Dialogs:**
  * Active: `kill agent <agentID> and its process tree? (y/N)`
  * Stale: `clear stale running status for agent <agentID>? (y/N)`
* **Safety Refusal & Demotion:** If the 4-check attribution fails, flash a warning: `[warnStyle] "Process verification failed (PID <pid> reused). Safe to clear stale badge instead."` and demote the `ctrl+k` action to badge-clearing (writing a synthetic `agent.failed` event with reason `stale process cleared by user`).
* **Steer Row Style:**
  * Live: `steer <agentID> › ` (cyan) with hint `Enter sends to <agentID> · ctrl+k kill`.
  * Observational/Record-only: `record steer <agentID> (read-only) › ` (muted gray) with hint `Enter records steer (no live execution) · run is read-only`.
* **Transcript Steer Render:** Steer queries and their replies are read from the event log and woven chronologically into the main transcript buffer, surviving tab-switching and restarts.
* **Empty-Reply / Muted States:** Surfacing empty steer replies clearly as `[agent returned an empty reply]` in the transcript and handling the finished-before-kill race gracefully.
* **Short-Terminal Compaction:** If terminal height < 15, hide the slash-command menu and compact the input area to preserve at least 3 visible lines for the transcript.

---

## Remaining risks

* **macOS `ps` / `sysctl` Output Parsing Flakiness:** The exact layout of `ps -p <pid> -o lstart=,command=` can vary slightly depending on locale and environment settings (e.g., date formats or spacing). The parsing logic must use flexible substring/regex patterns to avoid false-negative mismatches that refuse a legitimate kill.
* **Attribution Gate Mismatches:** If the command path or env markers are truncated or modified by the OS (especially in macOS's process list), the attribution gate may fail and refuse the kill, forcing a demotion to a badge-clear. While safe, this could leave a real process orphaned. We must ensure the command/marker checks are exact but match what the OS actually exposes.
* **Event Log Growth and Performance:** Reading all events from `events.jsonl` to render steer conversation and project states on TUI load/ticks must be efficient, especially for long-running sessions with hundreds of events. We must ensure the projection loader is fast and parses only what's necessary.
