---
agent: agy
idea: tui-command-picker
round: 1
date: 2026-06-06
---

## Summary
The goal of the `tui-command-picker` is to eliminate the need for users to manually memorize and type long, timestamped run IDs, question IDs, or slug identifiers. 

This proposal evaluates the UX flow, keyboard navigation, transition states, empty states, and discoverability. The focus is to make the picker intuitive, predictable, non-disruptive to power users, and highly accessible via keyboard-only interactions.

## Proposed approach

### 1. Command Eligibility & Consistency
To maintain a predictable mental model, we must establish a consistent pattern:
- **Eligible commands:** ONLY commands requiring an identifier target will trigger a picker:
  - `/open` (for selecting from ideas + recent runs).
  - `/answer` (for selecting from open HITL questions).
- **Ineligible commands:** Commands that do not take arguments (`/status`, `/home`, `/help`, `/follow`, `/quit`) or those that require arbitrary free-form input (`/deck <text>`) will not have a picker.
- **Affordance:** Typing `/open` or `/answer` with arguments (e.g., `/open my-slug-123`) runs the command immediately. Running them bare triggers the picker.

### 2. The `/answer` Two-Step UX Flow
Answering a question requires selecting the question ID (`qid`) and typing the response (`text`). We propose leveraging the TUI's existing spatial layout rather than adding a nested/second input mode:
- **Recommendation (Auto-tab Switch & Focus):**
  1. The user types `/answer` to open the picker and selects a question.
  2. Upon selection, the TUI automatically switches the active tab to the Agent tab associated with that question.
  3. The input prompt is instantly placed in answer mode (`answer <agentID>/<qid> › `) using the existing high-contrast warning style.
  4. The user types their answer and presses `Enter` to submit.
- **Fallback (No Agent):** If a question is not associated with any agent (e.g., a deck-level question), the TUI switches to the Status/Home tab and pre-fills the input box with `/answer <qid> ` (with a trailing space and cursor positioned at the end), allowing the user to simply type the text and press Enter.

### 3. Key Navigation & Filter-as-you-Type
- **Key Bindings:**
  - `↑` / `↓` moves the selection highlight in the picker.
  - `Enter` confirms the selection and triggers the command.
  - `esc` cancels the action.
  - Printable characters filter the selection.
  - `Backspace` / `ctrl+h` edits the filter text.
- **Filter Matching:** Simple, case-insensitive substring filtering. This is highly predictable and sufficient for short lists of runs/questions.
- **Escape Key (`esc`) Behavior:**
  - If the filter text is not empty: `esc` clears the filter first, returning the list to its full state.
  - If the filter text is already empty: `esc` exits the picker mode.
  - *Rationale:* Users expect `esc` to undo/reset their search input before completely closing the dialog.
- **Single-Match Selection:** Even if a filter reduces the candidates to a single item, the user **must still press Enter** to confirm. Auto-confirming feels jarring and can lead to accidental submissions.

### 4. Empty States
If there are no items to pick (e.g., no recent runs or no open questions):
- Render a centered, muted message in the picker area:
  - For `/open`: `(no recent runs or ideas to open)`
  - For `/answer`: `(no open questions to answer)`
- Do not draw a broken/empty table layout.
- The user can press `esc` to close.

### 5. Discoverability & Wording
We should update the existing hints and help screens to advertise the new picker affordance:
- **Input row hint:**
  - Change command suggestion from `... /answer <qid> <t> /open <slug|run> ...` to:
    `commands: /help /status /follow /deck <t> /answer [pick] /open [pick] /quit`
- **Help overlay (`/help`):**
  - Update description for `/open`: `/open [pick|target]` (Open picker or specify target)
  - Update description for `/answer`: `/answer [pick|qid text]` (Open question picker or answer directly)

## Concerns / open questions
1. **Scrolling & Height Constraints:** If a user has dozens of runs/ideas, how do we constrain the picker size?
   - *Recommendation:* Limit the picker's maximum height to 8 items. If the matched list exceeds this, allow `↑`/`↓` to scroll the viewport, showing `↑` and `↓` indicators at the boundary.
2. **Stale Data:** Should the picker perform a fresh fetch of runs/questions upon opening?
   - *Recommendation:* To keep the UI fast, display cached data (`homeRuns`/`questions`) instantly, but trigger background async reads to update the picker in real-time if a run finishes or a new question arrives.

## Risks
1. **Key Event Hijacking:** When the picker is active, key bindings like `N` (new idea) or `tab`/arrows (switch tabs) must be blocked to prevent unexpected side effects. Proper scoping of key routing in `updateMain` is critical.
2. **Vertical Space Collisions:** On small terminals (under 15-20 rows), displaying the picker overlay might cover crucial agent logs or the status table. We must ensure the picker does not crash the layout if height is constrained, potentially reducing its item count dynamically.
