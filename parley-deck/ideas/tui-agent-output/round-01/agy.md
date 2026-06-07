---
agent: agy
idea: tui-agent-output
round: 1
date: 2026-06-07
---

## Summary

The current agent tabs in the Parley TUI are frequently blank during runs because they only display `stdout`, whereas most agents emit their live progress to `stderr` and only write their final output (the artifact) at the very end. Furthermore, carriage returns (`\r`) used by streaming CLIs to update progress lines are stripped.

This proposal transforms the agent tab into a **live, scrollable conversation history** that mirrors the premium, interactive feel of the Codex CLI. We merge `stdout` and `stderr` into a single chronological feed, styling `stderr` to distinguish it from main output. We introduce real-time in-place carriage-return (`\r`) handling to allow streaming progress bars to update in place. We also integrate a toggleable **Artifact Preview** so users can immediately view "what they wrote" (the markdown protocol file), alongside a persistent **Status Header** so the tab is never blank.

---

## Proposed approach

### 1. The Conversation Layout (Steers + Replies + Streams)
Instead of the current behavior where the live steer reply panel (1.18) replaces the entire transcript history, we weave steers, stdout/stderr streams, and completed replies into a single scrollable timeline:
* **Steer Inception**: When a user inputs a steer, the TUI appends a styled block into the transcript:
  ```
  ❯ Steer: Why is the build failing?
  ```
* **Streaming Response**: Immediately below the steer block, the agent's live stdout and stderr streams begin printing. If the agent writes progress lines (e.g., using `\r`), the last line updates dynamically in place.
* **Finalization**: When the steer execution completes, a final boundary is drawn:
  ```
  [steer reply complete — 12s]
  ```
This chronological layout keeps the full context of the discussion visible and scrollable at all times.

### 2. Always-On Status Header
To ensure the tab is never blank even before the agent writes anything, we place a persistent, high-visibility status line immediately below the tab strip (at the top of the transcript area):
* **Working**: `● claude working... (elapsed: 1m12s)` (styled in warm orange/yellow).
* **Success**: `✓ claude wrote round-01/claude.md (finished in 2m30s)` (styled in vibrant green).
* **Failed**: `✗ claude failed: exit status 1 (elapsed: 45s)` (styled in red).
* **Killed/Stale**: `◌ claude killed` (styled in dimmed grey).

### 3. stderr Visibility and Merging
* **Chronological Merging**: We merge `stdout` and `stderr` into the same chronological list of lines. When the TUI polls the run directory (every 250ms), any new data appended to `stdout.log` or `stderr.log` is processed. Lines are stored in a structured list:
  ```go
  type TranscriptLine struct {
      Text    string
      IsError bool
      Time    time.Time
  }
  ```
* **Styling**: `stderr` lines are styled using a dimmed, low-contrast color (such as dark grey or a low-intensity purple/blue) and prefixed with a subtle `[err]` tag.
* **Stderr Toggle**: Users can press `ctrl+e` (or type `/stderr`) to toggle `stderr` visibility. When hidden, the renderer simply filters out lines where `IsError` is true.

### 4. In-Place Carriage-Return (`\r`) Handling
To achieve the Codex-CLI "potom sa to prepíše" (rewriting in place) behavior:
* **The `collapseCR` Algorithm**: For any string read from the logs, we simulate a physical terminal carriage return. When a `\r` is encountered, the write pointer resets to the start of the current line, overwriting previous characters but preserving trailing text if the replacement is shorter (unless cleared by space padding).
* **Handling Trailing Partial Lines**: Since progress lines frequently end in `\r` rather than `\n`, the tailing pipeline must not wait for a newline to display them. We modify the file reader to read the trailing partial line (the block after the last `\n`). We do not advance the persistent read offset past this partial line (meaning we will re-read it next tick), but we pass it to the renderer to display as the active, live-updating bottom line.
* **ANSI Handling**: We strip cursor control sequences (e.g., `\x1b[A` to move up) to prevent TUI rendering corruption, but we **preserve color ANSI escape codes** so that the agent's progress bars and styled text render correctly.

### 5. Artifact Preview Pane
The ultimate output of an agent's run is its Parley Deck protocol file.
* **Toggle Key**: The user can press `v` (or `a` for Artifact) to toggle the tab's main viewport between **Live Logs** (the combined stdout/stderr transcript) and the **Artifact View** (the markdown file itself).
* **Path Resolution**: The path is dynamically resolved based on the active run's metadata (e.g., `parley-deck/ideas/<slug>/round-<NN>/<agentID>.md` or `IMPLEMENTATION.md`).
* **Incremental Rendering**: If the file is being written to, we reload it on tick (guided by `os.Stat` modification time changes) and render the current markdown structure.

---

## Concerns / open questions

1. **Stdout/Stderr Interleaving Synchronization**:
   Since `stdout.log` and `stderr.log` are separate files, reading them sequentially at 250ms intervals might cause small chunks to appear slightly out of order if both streams are written to simultaneously. However, because agents typically write verbose progress to `stderr` and only write final results to `stdout` at the end, tick-based merging should be more than sufficient for a clean UX.
2. **Stderr Volume Control**:
   Some CLI tools are exceptionally noisy on `stderr` (e.g. dumping raw download progress bar frames). If this happens, it can overwhelm the transcript scrollback. Implementing a sensible default limit on repeating characters or consecutive `stderr` progress writes may be necessary if the `/stderr` hide toggle is not enough.

---

## Risks

1. **Regressions in Artifact Capture**:
   If we change the headless execution flags of agents (such as switching `claude` to `--output-format stream-json`), we risk corrupting the output captured in `stdout.log`, which is used by the runner's fallback mechanism to save the final artifact.
   * *Mitigation*: We should prioritize rendering `stderr` and in-place `\r` updates first. Changing agent CLI flags is high-risk; we should keep the execution commands exactly as they are today and let the TUI focus on rendering what is already written to `stderr.log`.
2. **Terminal Jitter During Scrolling**:
   If the user scrolls up to read earlier history, live in-place updates happening at the bottom of the transcript could cause the screen to jitter or force-scroll back to the bottom.
   * *Mitigation*: When `follow` is off (scrolled up), the viewport position must remain strictly anchored. New lines or in-place updates must be appended below the viewport without shifting the visible range.
