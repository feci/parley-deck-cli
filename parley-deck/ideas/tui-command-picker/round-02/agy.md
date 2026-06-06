---
agent: agy
idea: tui-command-picker
round: 2
date: 2026-06-06
---

## Cross-review (@claude / @codex / @hermes)

### @claude
The round-02 synthesis has successfully consolidated the group's views. I agree with the consensus to use a single, reusable `pickerState` on `liveModel` with a `pickerKind` discriminator and without function callbacks (ensuring the model remains a plain, testable value). 

I appreciate the adoption of my round-01 UX proposals:
1. Viewport height cap (max ~8 rows) with scroll boundaries and `↑` / `↓` scrolling indicators.
2. The specific empty-state wordings (`(no recent runs or ideas to open)` and `(no open questions to answer)`).
3. The `/help` overlay and bottom input-line hint updates to improve picker discoverability.

### @codex
I support the design details you outlined:
1. Placing the `if m.picker.Active { return m.updatePicker(msg) }` block at the absolute top of `updateMain` (immediately after `ctrl+c`). This guarantees the picker consumes keys before any tab navigation or steering logic can intercept them.
2. Keeping the explicit `&& !m.picker.Active` check inside the `N` key handler as a clear documentation of our state machine invariants.
3. Defining explicit reset states where cancellation and execution paths clear `composing`, `answerQID`, `picker.Active`, and `inputText` consistently.

### @hermes
I agree with your test criteria:
1. Re-entrancy verification (verifying that opening the picker from a run tab and hitting escape preserves the run attachment and tab focus).
2. Special-casing `N` typing so that it behaves as normal text filtering inside the picker, rather than starting a new idea.
3. Ignoring secondary scrolling and paging keys (`PgUp`, `PgDn`, `ctrl+u`, `ctrl+d`) while the picker is active to prevent visual jitter on the background transcript buffer.

---

## Response to overrules (esc, /answer)

### 1. Two-step esc
*   **The Overrule:** Single-esc-cancels the picker outright, rather than having the first press clear a non-empty filter and the second press close the picker.
*   **My Position:** **ACCEPT**. 
*   **UX Rationale:** In a terminal TUI, consistency is paramount. Escape has a well-established mental model of "close/back out immediately." Since candidate lists for `/open` and `/answer` are relatively short, losing a typed filter query is cheap, and correcting typos can be easily managed via backspace. Adopting single-esc-cancels avoids confusing the user with contextual escape behaviors and keeps command-line interactions snappy.

### 2. /answer flow
*   **The Overrule:** Reusing the existing `composing` state + `answerQID` field (rendering a high-contrast `answer <qid> › ` label) rather than performing an auto-tab switch to the question's agent tab.
*   **My Position:** **ACCEPT**.
*   **UX Rationale:** Auto-switching tabs adds state coupling (e.g. tracking the previous tab for rollback, resolving questions that are deck-level or aren't mapped to active agent buffers). Reusing the compose line is simpler, robust, and keeps the user's focus exactly where they are already typing. The dedicated warning label `answer <qid> › ` is sufficient to prevent accidental inputs.

---

## Confirmed for FINAL

The final specification must codify these UX and behavioral invariants:
*   **Triggering:** Only bare `/open` and `/answer` open the picker. If arguments are passed (e.g., `/open slug-123`), the command executes immediately without opening the picker.
*   **Filter Confirmations:** Single-match filtering still requires an explicit `Enter` to confirm. Auto-confirming on a single match is jarring and prohibited.
*   **Viewport Constraints:** The picker is capped at a maximum of 8 visible items. If the filtered list is larger, scrolling with viewport boundary indicators (`↑` and `↓`) must be rendered to guide the user.
*   **Empty State:** If a bare command is triggered but has zero candidate items (no questions or no runs), the picker does not open and an error is set in `inputErr`. If the picker opens and filtering results in zero matches, show the empty-state label.
*   **State Reset:** All transition boundaries (confirming, cancelling, run activation) must wipe out picker and compose-related variables: `composing = false`, `answerQID = ""`, `picker.Active = false`, `inputText = ""`, `inputErr = ""`.
*   **Verification:** Ensure unit tests cover: re-entrancy from active run tabs, key-routing invariants (active tab stays unchanged while typing in picker), case-insensitive substring matches, and filter text clamping.

---

## Remaining UX risks

1.  **Dynamic Height Clamping on Tiny Terminals:**
    On terminals with very small screen heights (e.g., height < 15 lines), a fixed 8-row picker plus tab strips, status bars, and input rows will cause rendering overflow or view crashes. 
    *   *Mitigation:* The maximum picker height must be dynamically capped to `transcriptHeight()` (minimum 3 lines) to ensure the picker adapts gracefully to tight vertical terminal spaces without breaking the layout.
2.  **Stale Data During Active Selection:**
    Because candidates are pulled from the cached `homeRuns` and `questions` lists, a run completion or new question tick could arrive while the user has the picker open. 
    *   *Mitigation:* Ensure that background `eventsMsg` and `questionsMsg` updates rebuild the picker items list reactively without resetting the user's active filter text (`picker.Filter`) or destructively resetting the cursor.
3.  **Visual Focus Muting for Tabs:**
    While the picker is active, steering and tab switching via arrows is disabled, but the tab strip remains visible. This can confuse users into thinking they can still interact with the tabs.
    *   *Mitigation:* The rendering method `renderTabStrip` should visually dim or mute all tabs (using `mutedStyle`) while `picker.Active` is true, drawing focus exclusively to the command picker overlay.
