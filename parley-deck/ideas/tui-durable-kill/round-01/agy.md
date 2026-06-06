---
agent: agy
idea: tui-durable-kill
round: 1
date: 2026-06-06
---

## Summary

As the **UX correctness and edge cases** owner, my goal is to ensure the TUI provides an unambiguous, reassuring, and highly visible experience for process control and agent steering. The current issues—where an owner is misled by a "running" badge for days, has no way to terminate processes on resumed runs, and types blindly to steering inputs without receiving replies or understanding why—stem from a mismatch between the UI's visual state and the underlying process state.

This proposal focuses on:
1. **Per-Agent Liveness Visualization**: Transitioning the tab badges and status line from blind log-based states to active process-probed states (`RUN` vs `STALE` vs `DEAD`).
2. **Reattach & Kill Seam UX**: Providing clean kill controls on restarted or resumed runs, with explicit confirmation and success/failure messaging.
3. **Unmistakable Steering Modes**: Visually distinguishing live-executing runs from read-only observational runs in the input row, and ensuring the steer conversation is durably rendered in the agent's main transcript.
4. **Resilient Edge-Case Handling**: Addressing short terminals, PID-reuse safety alerts, and empty steer replies.

---

## Proposed approach

### 1. Process-Group Spawn (UX aspect)
Spawning agents in their own process groups (unix `Setpgid` / windows `CREATE_NEW_PROCESS_GROUP`) ensures that when a user triggers a kill, we reap the entire tree (including orphaned grandchildren). 
* **UX Feedback**: The TUI status line should display a brief, reassuring message during the kill phase, such as: `[warnStyle] "Killing agent and all sub-processes..."`. Once done, it transitions to a clean terminal state message. This avoids leaving zombie processes running in the background while the UI reports termination.

### 2. PID/PGID Persistence + Lifecycle (UX aspect)
We will carry the `pid` and `pgid` in the `agent.started` event data and/or a small `runs/<id>/agents/<a>/proc.json` file. 
* **UX Enumeration**: When the TUI starts or attaches to a run, it immediately reads the event log/process directory to extract the PIDs. The TUI does not require a persistent daemon; instead, it uses a background tick loop (or on-demand refresh when switching tabs) to verify the process list. If a PID is found in the proc file but is dead in the OS, it transitions the UX state from `RUN` to `STALE`.

### 3. Durable, Cross-Restart Kill
Resumed and opened runs currently get no `KillAgent` seam in [internal/tui/live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L44-L70). We will inject a stateless "reattach kill" seam that operates using only the run directory and the persisted PID/PGID.
* **Confirmation Dialog**: When pressing `ctrl+k` on a running or stale agent tab, the TUI prompts the user:
  * For active: `kill agent <agentID> and its process tree? (y/N)`
  * For stale: `clear stale running status for agent <agentID>? (y/N)`
* **Termination UX Messages**:
  * **Found-and-killed**: `[okStyle] "Successfully terminated agent <agentID> (PID <pid>) and all sub-processes."`
  * **Already-dead**: `[warnStyle] "Agent <agentID> was already dead. Cleared stale badge."`
  * **Can't-attribute (PID reuse/verification failure)**: `[warnStyle] "Process verification failed (PID <pid> reused). Safe to clear stale badge instead."`

### 4. Stale-Running Reconcile
Currently, [internal/runstate/runstate.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runstate/runstate.go#L499-L506) derives liveness per-run based only on event states, leading to stale badges.
* **UX Liveness Indicators**:
  * **Active Running**: `[RUN]` (cyan/green badge in tab, e.g. `claude RUN`)
  * **Stale Running**: `[STALE]` (yellow/warn badge in tab, e.g. `claude STALE`).
  * **Dead/Finished**: `[FIN]` / `[ERR]` / `[KILL]` (muted/standard badges).
* **Status Line Wording**:
  * Normal: `claude RUNNING 02:45`
  * Stale: `claude STALE (process lost/dead) · Press ctrl+k to clear`
* **Reconciliation Action**: Pressing `ctrl+k` on a `[STALE]` tab opens a confirmation dialog. Upon confirmation (`y`), it writes a synthetic `agent.failed` event to `events.jsonl` containing the error message `stale process cleared by user`. This durably reconciles the log on disk so the state updates to `[ERR]` (or `[KILL]`), preventing the run from hanging indefinitely in future sessions.

### 5. PID-Reuse Safety (Critical UX Safeguard)
If a PID was reused by an unrelated system process, signaling it could crash important user applications.
* **UX Safeguard**: We check the process start time, boot ID, or cmdline. If the verification fails, the TUI **refuses** to send `SIGKILL`.
* **UX Warning**: The TUI status line flashes: `[warnStyle] "Refused to kill PID <pid>: process verification failed (PID was likely reused). Safe to clear stale badge instead."`
* **Fallback Behavior**: The action of `ctrl+k` automatically demotes from "Kill Process" to "Clear Stale Badge" (writing the synthetic terminal event to clean up the TUI, without touching the OS process).

### 6. Steer on Resumed Runs & Visibility
To resolve the owner's "no conversation/reply visible" issue:
1. **Unmistakable Mode Hinting**:
   * For **Live executing** runs: The input line prefix is cyan `steer <agentID> › `, with the hint `Enter sends to <agentID> · ctrl+k kill`.
   * For **Record-only/Observational** runs: The input line prefix is gray/muted `record steer <agentID> (read-only) › `, with the hint `Enter records steer (no live execution) · run is read-only`.
2. **Conversation Persistence**:
   * Instead of a transient steer reply panel that is cleared by `esc` or tab switching, we will read the steer history from `events.jsonl` (`steer.submitted`, `steer.replied`, and `steer.reply_failed` events) and render the conversations directly in the main transcript scrollback of the agent's tab. 
   * When a steer is submitted, we print a clear block in the scrollback:
     ```
     ┌── Steer Sent by User ────────────────────────────────────────────────────────┐
     │ > [user query text]                                                          │
     └──────────────────────────────────────────────────────────────────────────────┘
     ```
     The agent's reply streams directly below this block in the transcript, surviving tab switches and TUI restarts.
3. **Empty Reply Hardening**:
   * If an agent (`agy` or otherwise) returns an empty reply, the TUI displays a distinct indicator: `[agent returned an empty reply]` in `warnStyle` instead of leaving a blank/silent space, giving the owner immediate feedback.

### 7. Edge Cases
* **Short Terminals**: If the terminal height is extremely small (< 15 rows), we suppress the slash-command suggestion menu and compact the input/status lines to ensure the active transcript retains at least 3 visible rows.
* **Many Stale Agents**: The tab strip dynamically handles overflow by displaying scrolling indicators and `…+N` markers. A global status banner or home tab notification will alert the user if multiple stale agents are detected.
* **Killing an Agent that Just Finished**: If the user confirms a kill but the agent finishes naturally right before the signal is delivered, the TUI handles the error gracefully, updates the state to `[FIN]`, and displays: `[okStyle] "Agent <agentID> finished before it could be killed."`

---

## Concerns / open questions

1. **Auto-Reconciliation on Startup**: Should the TUI automatically reconcile stale agents on startup/attach by writing `agent.failed` events to the log, or should it always require a manual `ctrl+k` "clear stale badge" action? 
   * *Recommendation*: A manual `ctrl+k` confirmation is safer to avoid polluting the event log with synthetic failures if the process is temporarily unresponsive but still alive (e.g. swap storms). However, a global banner warning the user about stale agents should be shown.
2. **Steer Event History rendering**: Can we easily parse `events.jsonl` for older steer attempts and merge them chronologically with the stdout tail? 
   * *Recommendation*: Yes. The runner already records steer metadata. We should read these events and weave them into the scrollback buffer so they render as formatted message bubbles.

---

## Risks

* **Accidental OS Kill**: If the PID verification mechanism is too lenient, there is a risk of killing an unrelated system process. We must enforce strict checking (start time or cmdline marker) and default to safe badge clearing if there is any doubt.
* **UI Clutter in Small Terminals**: Formatting steer blocks, warnings, and confirm modals takes vertical space. We need rigorous height checks in `renderTabbed` to prevent text overlap in tiny windows.
