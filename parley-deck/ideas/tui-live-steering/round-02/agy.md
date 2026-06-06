---
agent: agy
idea: tui-live-steering
round: 2
date: 2026-06-06
---

## Cross-review

### @claude
* **Autocomplete & Sub-Mode:** I fully agree with your synthesis to implement a dedicated, non-modal `suggest` sub-mode rather than reusing the modal `pickerState`. Keeping `inputText` editable while showing autocomplete options solves the input hijacking problem.
* **Tab Semantics Refinement:** I accept your decision to drop Tab-cycling and use Tab only for longest-common-prefix (LCP) completion. Using ↑/↓ + Enter to pick from suggestions is standard and avoids conflicting with the trailing-space picker hand-off.
* **Force-Steer Decoupling:** I accept the queue-not-force model. Keeping the steer queue independent of immediate process termination is architecturally cleaner. The user can manually use `ctrl+k` to terminate an active round run if they wish to expedite a queued steer.
* **Lost UX Details:** Your synthesis omitted my round-01 concerns regarding short-terminal height constraints and explicit TUI error surfacing for steer/kill failures. I have reintroduced these as counter-proposals below to ensure they are captured before finalization.

### @codex
* **Handle-based Runner Registry:** I support your execution design. Having a mutex-guarded `active` attempts registry and the `RunSteerAttempt` + `KillAgent` operations on `Handle` is correct and robust.
* **Per-Steer Directories:** I agree that steer attempts must write to separate directories (`runs/<runID>/agents/<id>/steers/<steerID>/`) rather than clobbering the main agent log. This prevents log truncation and preserves the audit trail.
* **Seams Copying:** Your catch that `activateRun` must copy `SubmitSteer` and `KillAgent` onto the active model is critical. Without this, TUI runs launched from the Home tab would silently lose steering capabilities.
* **Bypassing Artifact Validation:** I agree that steer attempts must not trigger standard round artifact checks or fallbacks, as steers are transient re-invocations.

### @hermes
* **Conditional-Tab Keymap:** I support the conditional-Tab logic: Tab/shift+tab only switch tabs when the input row is not slash-prefixed. Left/right arrows always switch tabs, ensuring muscle memory navigation remains intact.
* **Modal Kill Confirmation:** I agree with your `confirmKillAgentID` modal state design. It must block all other keys (including Enter and alphabetical keys) to prevent accidental steer submissions or cursor movement during confirmation.
* **Suggest Mode Overrides:** I agree with overriding your `pickerSuggest` suggestion in favor of a dedicated non-modal `suggest` state, as it keeps typing non-blocking.

### @agy
* **UX Lens Review:** As `agy`, I confirm that the autocomplete flow (slim menu, Tab = LCP, ↑/↓ + Enter, trailing-space hand-off to the picker) is highly discoverable, non-intrusive, and correct. 
* **Accidental Steer Guard:** The synthesis's adoption of my proposal for a cyan-colored `steer <id> › ` input prefix provides an excellent visual distinction to prevent accidental submissions.
* **Steer Reply Surfacing:** Surfacing the reply inline in the same tab behind a divider with a replying spinner is the most cohesive layout. However, I want to make sure we specify the exact divider style and tail styling (gray/faint) so it is clear in the scroll history.

## Counter-proposals

### 1. Surfacing Steer and Kill Errors
The synthesis lacks specifications on how the TUI handles execution or submission errors. I propose:
* **Synchronous Errors:** If `opts.SubmitSteer` or `opts.KillAgent` fails immediately (e.g. queue full, run already ended), display a red error message in the input row for 3 seconds or until the next keypress.
* **Asynchronous Errors:** If a steer attempt fails during execution (e.g., API failure, rate limits), append a clear warning block to the agent's transcript tab:
  ```text
  ┌────────────────────────────────────────────────────────┐
  │  [!] STEER REPLY FAILED (Segment: segment-NNNN)        │
  │  Error: <error details or exit code>                   │
  └────────────────────────────────────────────────────────┘
  ```

### 2. Short-Terminal Height Constraints
To prevent rendering and UI layout breaks on small terminals:
* If the terminal height is too small (`height < 10` or `transcriptHeight() < 3`), the floating suggest menu must be suppressed. Instead, show a simple inline hint in the status row like `(Type / for command suggestions)`.

### 3. Styling of Steer Replies in Transcript
To keep the agent's chronological log clean and readable:
* Render the divider with the query text: `── steer <steerID>: "<truncated user query>" ──`.
* Render the stdout of the steer attempt in a slightly dimmed style (e.g., `faintStyle` or a distinct gray) to differentiate it from the main round logs.

### 4. Busy Agent Queue Hint
Since a steer on a busy agent is queued, we must let the user know they can force it:
* When a steer is queued, display `[Steer queued. To run immediately, cancel the active attempt with ctrl+k]` in the status row or transcript warning.

## Confirmed for FINAL

1. **Suggest Menu:** Non-modal overlay above the input row, visible only when `inputText` starts with `/` and contains no spaces. Tab completes longest common prefix (LCP). If a single match remains, completes command + space (if arguments are expected).
2. **Tab Navigation:** Tab and shift+tab switch tabs *only* when the input is not slash-prefixed. Left/right arrows always switch tabs.
3. **Kill Agent:** `ctrl+k` on a running/steering agent tab triggers a modal `confirmKillAgentID` state. The input line changes to warning colors and prompts `Kill agent <agent_id>? (y/N)`. If confirmed, calls `opts.KillAgent(id)` and registers an `agent.killed` event. The overall run context is unaffected.
4. **Steer Attempt:** Steers write to `<agentDir>/steers/<steerID>/` and tail their stdout directly into the target agent's transcript tab behind a divider. They bypass round artifact validation. Only one steer attempt can run or be queued per agent (depth-1 queue, second is rejected).
5. **Seams Mapping:** `SubmitSteer` and `KillAgent` hooks copied inside `activateRun` for both `LiveOptions` and `LaunchResult`.

## Remaining risks

* **Subprocess Orphans:** Context cancellation via `exec.CommandContext` might leave orphaned child processes spawned by the agent CLI. We should document process group signaling as a fast-follow task.
* **Segment ID Race Condition:** Concurrent steer attempts might try to claim the same segment ID. We must ensure `appendSegmentStarted` is protected by a handle-level mutex.
* **Tab-switching Muscle Memory:** Users who type `/` and expect Tab to switch tabs immediately will experience a shift in behavior. Clear placeholder/hint text updates must be implemented.
