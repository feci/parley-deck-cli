---
agent: agy
idea: tui-live-steering
round: 1
date: 2026-06-06
---

## Summary

This independent analysis evaluates the three requested TUI features—Slash-command autocomplete, single-agent process termination (kill), and steer round-trip execution—focusing strictly on user experience correctness, interface consistency, discoverability, and edge-case handling. 

From a UX lens, these features must feel premium, predictable, and cohesive. The autocomplete menu must assist typing without hijacking the input flow, the kill command must be explicit and safe from accidental triggers, and the steer round-trip must provide immediate visual feedback, transforming the TUI from a passive tail viewer into an active steering console.

---

## Proposed approach

### 1. Slash-Command Autocomplete UX
*   **Menu Presentation:** Instead of reusing the full-screen [pickerState](file:///Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L117-L126) (which intercepts all keyboard events and is designed for heavy searching/filtering), we will implement a **slim floating suggestion menu** rendered directly above the persistent input row.
    *   This menu only appears when `inputText` begins with `/` and contains no spaces (e.g., `/o` but not `/open `).
    *   It lists matching slash commands (`/help`, `/status`, `/follow`, `/deck`, `/answer`, `/open`, `/home`, `/quit`).
    *   The menu is compact (capped at 4 visible rows with scroll indicators if needed) and does not disrupt the main transcript layout.
*   **Interaction & Tab Behavior:**
    *   **Normal Typing:** Normal typing is never intercepted by the autocomplete menu; characters are appended to `inputText` and filter the suggestions in real-time.
    *   **Tab Key:** Pressing **Tab** completes the longest common prefix of the matching commands. If there is a single match (e.g. `/op` -> `/open`), it completes the full command name and appends a trailing space (e.g., `/open `). Subsequent presses of Tab cycle through the matching list if multiple options remain.
    *   **Arrow Keys:** While suggestions are visible, `up`/`down` arrow keys navigate the selection in the autocomplete menu without changing the text cursor position.
    *   **Enter Key:** Pressing **Enter** completes the input field with the highlighted command. If the command takes no arguments (e.g. `/status`, `/follow`, `/home`, `/quit`, `/help`), it executes immediately. If it takes arguments (e.g. `/open`, `/answer`), it populates the command and automatically transitions into the appropriate picker mode (e.g. `/open ` immediately launches the recent run picker).
*   **Argument Hand-off:** Typing `/open ` or `/answer ` with a space immediately triggers the existing `pickerState` with the respective `pickerOpen` or `pickerAnswer` kinds. This bridge ensures a smooth transition from text-entry autocomplete to visual candidate selection.

### 2. Per-Agent Process Termination (Kill) UX
*   **Trigger Key:** We propose `ctrl+k` (Standard CLI "Kill" shortcut) instead of `K` or other single-letter hotkeys. Since the bottom input row in [live.go](file:///Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L641-L693) is "always typeable", hotkeys like `K` are prone to accidental invocation and conflict with normal message composition.
*   **Safety Confirm:** Pressing `ctrl+k` on an active agent tab (only when that agent's state is `RUN` or `STEER`) prompts the user with a confirmation prompt in the input row:
    `Kill agent <agent_id>? (y/N)`
    The input row label is color-flipped to red/orange (`warnStyle`) to signal danger. The user must type `y` or `Y` to confirm; any other key cancels the action and restores normal input.
*   **Visual Representation of Terminated State:**
    *   The badge in the tab strip changes from `RUN` to `KILL` (or `KILLED`) in a dimmed, muted red style.
    *   A system event message is appended to the agent's transcript:
        `--- AGENT TERMINATED BY OWNER AT 15:34:22 ---`
    *   The runner registry cancels the specific agent's context, freeing it to be re-run or steered. The rest of the deliberation run keeps going unaffected.

### 3. Steer Round-Trip UX
*   **Targeting Indicators:** The bottom input row must make the recipient of the steer message clear. The label changes based on the active tab:
    *   On `Home` tab: `› ` (no-run state, typing is disabled or prompts for run/idea launch).
    *   On `Status` tab: `steer deck › ` (submits a deck-level steer).
    *   On `agent:<id>` tab: `steer <id> › ` (submits an agent-specific steer).
*   **Active Processing State:**
    *   When the user submits a steer message, the input row resets, and a temporary spinner appears in the active agent's tab badge: `agy ⠋`.
    *   A bold visual banner is appended to the active agent's transcript view to signal the steer is active:
        ```text
        ┌────────────────────────────────────────────────────────┐
        │  STEER INPUT: <user text>                              │
        │  Sent at 15:34:02 (Segment: segment-0002)              │
        └────────────────────────────────────────────────────────┘
        ```
    *   While the steer re-invocation process is running, a live indicator displays at the bottom of the transcript pane: `⠋ agy is replying...`
*   **Reply Surface:**
    *   We reject creating a separate tab or dedicated chat pane for replies. Splitting the interface dilutes focus.
    *   Instead, the steer attempt's output (stdout of the re-invoked process) appends directly into the active agent's transcript tab, underneath the steer input banner. This preserves a chronologically unified history of the agent's output and subsequent manual steering.

---

## Responses to focus questions

### 1. Steer execution model (the crux)
*   **Parallel vs. Sequential:** We must enforce a **single-active-process-per-agent** invariant. Running multiple parallel attempts for the same agent will clobber its directories and write-logs, leading to race conditions.
*   **Execution Strategy:**
    *   If the agent is **idle** (pending, finished, failed, or killed), the steer re-invocation runs **immediately**.
    *   If the agent is **currently running** its main round attempt:
        *   The steer is **queued**.
        *   The agent's transcript tab displays: `[Steer queued: waiting for active round attempt to complete...]`
        *   Once the active process finishes, the queued steer execution automatically launches.
        *   *Alternatively*, the user can press `ctrl+k` to terminate the active attempt immediately and force the steer to run.
*   **Context Payload:** The steer attempt re-invokes the agent's CLI with:
    *   The user's steer message as the primary prompt input.
    *   The agent's latest written artifact (from the active round directory) and the tail of its stdout log (up to 4KB) appended to the instruction context to maintain continuity.
*   **Durability:** The execution writes to a dedicated log file (`<agentDir>/steer-<segmentID>.log`). The TUI's tail reader merges this log into the active tab buffer. The execution status is recorded via `steer.requested`, `steer.delivered`, and a new `steer.replied` event in `events.jsonl` to maintain the audit trail.

### 2. Per-agent kill mechanism
*   **Registry & Signal:** The runner holds a registry of context cancel functions keyed by agent ID. Calling `KillAgent(agentID)` invokes the cancel function, sending an interrupt/kill signal to the process.
*   **State & Continuation:** The killed agent transitions to `KILLED` status, emitting an `agent.killed` event. The overall run's coordinator ignores this failure as a blocker, permitting other agents to finish. The killed agent is then treated as idle, allowing the user to steer or re-run it manually.
*   **Safety Invariant:** A per-agent cancellation context child ensures that triggering a kill never bubbles up to cancel the global run context.

### 3. Autocomplete UX + mechanism
*   **Dedicated Floating Menu:** We choose a dedicated, floating autocomplete menu over the `pickerState`. The autocomplete menu does not capture keyboard input (except `up`/`down` for selection, `Tab` for completion, and `Enter` to confirm), allowing normal typing to filter suggestions inline.
*   **Tab Semantics:** Tab completes the longest common prefix. If there's a unique match, it completes it and adds a space. If multiple matches remain, subsequent Tab presses cycle through them.
*   **Argument Hand-off:** Typing a trailing space after `/open` or `/answer` immediately opens the full-screen picker.

### 4. Reply surfacing
*   **Inline Appending:** The reply stdout is appended directly to the agent's transcript buffer in the active tab. 
*   **Visual Dividers:** Markdown-style box borders isolate the steer query, followed by the live-tailed stdout of the steer process.
*   **Liveness Indicator:** While the steer runs, the tab badge displays a loading spinner, and the transcript bottom shows `[agent is replying...]`.

### 5. Concurrency & safety
*   **Mutex / Segment Locking:** A mutex or state-lock per agent directory prevents concurrent execution. If a process is active, any new execution request (like a steer) is queued.
*   **Steering Ended Runs:** If the overall deliberation run has completed, the user can still steer any agent. The TUI re-enables input, spawns the steer attempt, and updates the agent's tab status to show it is active.

### 6. Seams & testability
*   **Decoupled Seams:** We define abstract hook signatures in `LiveOptions` to keep `internal/tui` decoupled from `internal/runner` and `internal/app`:
    ```go
    type LiveOptions struct {
        // ... existing options ...
        KillAgent   func(agentID string) error
        SubmitSteer func(agentID string, text string) (<-chan string, error) // Returns stdout channel or log path to tail
    }
    ```
*   **Headless Testability:** We can test the autocomplete logic, key interceptors, and confirm transitions inside `live_test.go` using bubbletea's headless test program (`tea.NewProgram` in test mode), sending key sequences and asserting the model's text buffers and view outputs.

---

## Concerns / open questions

1.  **Tab Collision on Autocomplete:** When the autocomplete menu is active, `up`/`down` arrows select menu items, but what about `left`/`right` or `tab`? Since `tab` and `shift+tab` normally switch agent tabs, we must ensure that while the autocomplete menu is open, `tab` is hijacked exclusively for command completion. Once the menu closes (e.g. by typing a space or backspacing the `/`), `tab` reverts to switching tabs.
2.  **Handling Steer Errors:** If a steer re-invocation fails immediately (e.g., model API rate-limit or missing executable), how is that error surfaced? We should display a transient error message in red text in the input row (using `inputErr`) and append a system failure notice to the agent's transcript.
3.  **Terminal Height Constraints:** On extremely short terminals (height < 10), the autocomplete menu, input row, and status lines might consume all vertical space, leaving no room for the transcript. We should automatically suppress the floating autocomplete menu if `transcriptHeight() < 3` and fallback to simple inline text hints.

---

## Risks

*   **Accidental Steer Submissions:** Unlike answer mode, which has an explicit orange safety color-flip, normal typing on an agent tab defaults to sending a steer message. There is a risk of a user typing a message, thinking they are in another mode, and hitting Enter, triggering a costly API call. We must make the input prefix `steer agy › ` visually distinct (e.g., using a cool cyan color) to remind the user of their active context.
*   **Agent Directory Lock Contention:** If the runner does not properly release directory locks or file descriptors on agent termination, a subsequent steer execution will fail due to write permissions or locked state files. Codex must ensure clean state release in the runner's per-agent process control.
*   **Orphaned Processes:** If the user kills an agent via `ctrl+k` and exits the TUI quickly, there is a risk of leaving an orphaned CLI process running on the host system. The runner must ensure that killing an agent sends a robust process-group kill (`syscall.SIGKILL` or equivalent) to terminate any child processes spawned by the agent's CLI.
