---
agent: agy
idea: tui-agent-output
round: 2
date: 2026-06-07
---

## Cross-review (@claude/@codex/@agy)

I strongly support the synthesis and resolved decisions outlined in the round-02 facilitator document. The design successfully resolves the forks and delivers a premium, highly interactive Codex-CLI-style feel while maintaining performance and architectural safety.

*   **@claude**: The woven conversation layout is a major improvement. Appending steers (`❯ you:`) and streaming replies directly into the main scrollback history instead of swapping panels (as in 1.18) feels natural and preserves context. Additionally, adopting `/stderr` and `/artifact` as autocomplete-supported slash commands instead of custom single-character keybindings (such as `v`, `a`, or `ctrl+e`) is the correct UX choice to avoid typing collisions with the active input box.
    
    *UX Styling Clarification*: We must ensure that stripping all ANSI color codes from the log streams for v1 (to avoid breaking naive truncation and width calculations) does *not* result in a completely plain, monochrome interface. The TUI-native styling (implemented via `lipgloss` styles in [internal/tui/live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go)) must remain fully intact. Specifically, the status header (green/orange/red indicators and elapsed times), the `❯ you:` prefix, the `[err]` tag, and the dimmed styling of the stderr lines themselves must use rich, premium terminal colors. The ANSI stripping must only apply to raw data parsed from the log streams.
*   **@codex**: The dual-cursor `tailCursor` pipeline for stdout/stderr within `agentBuffer` is robust and maintains a bounded memory footprint. I agree with your decision to defer streaming flags for Claude/others to a safe parser follow-up, keeping `HeadlessArgs` unchanged in v1 to protect the runner's stdout artifact recovery fallbacks. I also endorse the lower-risk byte-capped trailing partial line strategy for live updates.
*   **@agy** (reflecting on my round-01 UX): The synthesis captures the core UX goals of a live, non-blank, and responsive interface. The status wording (`● working...`, `✓ wrote...`, `✗ failed...`) is exactly what is needed to make the tab feel alive immediately.

## Counter-proposals (if any)

I have no fundamental disagreements with the synthesis. However, I propose three minor UX refinements to ensure maximum clarity during implementation:

1.  **Visual Mode Indicator**: When the user toggles `/artifact` mode on, the viewport should display a clear, styled header (e.g., `[Viewing Artifact: <path>]` in a distinct style, such as dimmed green) so that the user immediately realizes they are viewing the generated markdown file rather than the live execution log.
2.  **Clickable/Copyable Artifact Path**: When a run completes successfully, the status header should display the relative path to the artifact (e.g., `✓ wrote round-02/agy.md`), and this path should be styled to make it easily readable and copyable (and clickable if the terminal emulator supports standard file links).
3.  **Empty Carriage Return Filtering**: If a log line becomes completely blank after CR collapse (e.g., repeated carriage returns without text or empty progress updates), the ingester should drop it instead of rendering empty lines in the scrollback, preventing vertical clutter.

## Confirmed for FINAL

1.  **Conversation Weave**: Steers and live replies are appended to the main scrollback history chronologically, preventing viewport jumps when typing.
2.  **Always-on Status Header**: Displaying state, elapsed time, and final artifact paths from `AgentState` directly below the tab strip.
3.  **Merged Streams with /stderr Toggle**: Chronological per-tick merging of stdout and stderr. Stderr is dimmed and tagged with `[err]` by default, and can be toggled off with the `/stderr` slash command to keep the log uncluttered.
4.  **CR-Aware Ingest**: Implementing the custom helper to rewrite live lines on a lone `\r` while preserving normal line breaks for `\n` and `\r\n`.
5.  **Artifact View via /artifact**: Toggling between live logs and the generated markdown artifact using a slash command.
6.  **Unchanged HeadlessArgs**: Keeping agent execution flags untouched for v1 to guarantee no regressions in stdout fallbacks or validation.

## Remaining risks

*   **Viewport Yanking on Scrollback**: Verify that when `follow` is false (user scrolled up), updates to the live partial at the bottom do not cause the viewport to scroll back down or jitter.
*   **Prefix Duplication**: Ensuring that during tick-by-tick partial reads, tracking the line-start offset is robust so that re-reading does not duplicate the stream prefix.
*   **CR Interpretation vs Stripping**: Ensure the CR-aware ingester parses carriage returns *before* ANSI stripping takes place, in case some sequences contain embedded returns or are stripped in a way that breaks sequence matching.
*   **Artifact File Verification**: Ensure that running `/artifact` on a missing or not-yet-written file fails gracefully (e.g., displaying `[Artifact not yet written...]` or `[Artifact empty]`) instead of crashing or leaking file descriptors.
