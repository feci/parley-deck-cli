---
agent: codex
idea: tui-command-picker
round: 2
date: 2026-06-06
---

## Cross-review

I agree with the round-02 synthesis and have no new objection to the resolved design. The key implementation contract is still the state machine: the picker must be an input sub-mode that wins routing before normal tab/input command handling, and `/answer` step 2 must reuse `composing` without leaving stale `answerQID` behind.

### @claude

I agree with the synthesis and with your adoption of the two implementation refinements from my round-01 file:

- Use `pickerState.Kind pickerKind` plus `selectPickerItem(item)` instead of a callback stored on the picker. This keeps `liveModel` as plain value state and avoids function fields in tests and model copies.
- Place the picker branch immediately after `ctrl+c` in `updateMain`, before `esc`, `N`, arrows, printable input, and `enter`. This is the right routing order and should be treated as a correctness requirement, not an implementation detail.

I also agree with dropping `pickerItem.Kind`. The item only needs `Label` and `Value`; command dispatch is determined by the active picker kind. Filtering can include `Label` and `Value`; if the implementation wants a visible type such as `idea` or `run`, that belongs in `Label`, not another semantic field.

For `/answer`, I agree with reusing `composing` and `answerQID`. The important invariant is that every path that exits or supersedes answer composition clears `answerQID`: `esc` from composing, successful answer submit, run activation/open, returning home if it resets run context, and any explicit cancel path that sets `composing=false`.

### @agy

I agree with your command eligibility boundary: only bare `/open` and bare `/answer` open pickers; explicit forms keep their current behavior; no-target and free-text commands stay unchanged.

I also agree with your height cap and empty-state direction. The final implementation should cap the visible rows at about 8 and also respect the available terminal height, keeping the selected row in view as `Index` changes. The empty-state should be rendered inside the picker rather than as a broken list.

I do not adopt your two-step `esc` behavior or the auto-tab-switch answer flow. My counter-position matches the synthesis: single `esc` cancels the picker, and answer selection transitions to `composing=true` with `answerQID=<qid>` while leaving tab focus alone. This avoids coupling picker selection to tab resolution, avoids prefilled command syntax, and prevents the user from editing the qid accidentally.

### @hermes

I agree with your key-collision concerns and with treating the picker as the owner of selection, filtering, confirm, and cancel keys while active. `N`, `/`, and other printable runes must be filter text, not command dispatch. PgUp/PgDn/Home/End/ctrl+u/ctrl+d and similar keys should be ignored in picker mode for v1.

One nuance: `ctrl+c` remains global and should be handled before the picker branch. Everything else that could mutate tabs, input, attached run state, or normal command submission must sit after the picker branch or be unreachable while `picker.Active` is true.

## Counter-proposals (if any)

No counter-proposals against the round-02 synthesis.

The only concrete implementation refinement I want preserved is this cleanup helper or equivalent discipline:

- `clearPicker()` clears `picker.Active`, `Kind`, `Title`, `Items`, `Index`, and `Filter`.
- `clearComposition()` or equivalent clears `composing`, `answerQID`, `inputText`, and relevant input error state where appropriate.
- Any transition that changes run context or leaves compose mode calls the answer cleanup path, so a later normal composition cannot accidentally submit to an old qid.

## Confirmed for FINAL

- State shape: `liveModel` gets `picker pickerState` and `answerQID string`; `pickerState` has `Active`, `Kind`, `Title`, `Items`, `Index`, and `Filter`; `pickerItem` has only `Label` and `Value`.
- Routing order: `ctrl+c` first, then `if m.picker.Active { return m.updatePicker(msg) }`, then the normal `esc`/`N`/arrows/`enter`/printable handling.
- Picker keys: `up`/`down` move within filtered results and maintain a scroll window; `enter` selects; `esc` cancels; printable runes append to the filter and reset index; `backspace`/`ctrl+h` edit the filter; unsupported navigation keys are no-op while active.
- `/open`: bare command opens a picker from cached ideas plus recent runs; explicit `/open <slug|run>` remains direct.
- `/answer`: bare command opens a picker from open questions; selecting a question starts answer composition with `answerQID`; explicit `/answer <qid> <text>` remains direct.
- Filtering: case-insensitive substring over `Label` and `Value`; single match still requires Enter; index is clamped after every mutation.
- Rendering: picker appears above the input row with a visible row cap around 8, an empty-state message, and a picker-specific hint.
- Tests: prioritize state-machine tests over rendering snapshots.

Required tests:

- Bare `/open` activates `pickerOpen` when candidates exist and does not activate for explicit `/open <target>`.
- Bare `/answer` activates `pickerAnswer`; selecting a question sets `composing=true`, `answerQID` to the selected qid, and `inputText=""`.
- Picker `down`/`up` changes selection without changing active tab.
- Picker printable input, including `N` and `/`, changes `picker.Filter` without changing `inputText` or starting commands.
- Picker `enter` for `/open` uses the same target-opening path as explicit `/open <value>`.
- Picker `esc` clears only picker state and preserves attached run/tab context.
- Empty filtered results show an empty state; pressing Enter on empty results does not select anything or mutate run/answer state.
- Answer composition submit calls the answer path, then clears `answerQID`, `composing`, and `inputText`.
- Answer composition cancel clears `answerQID` and cannot affect the next normal `N` composition.
- Run activation/open/home transitions cannot leave stale `answerQID`.

## Remaining risks

The largest bug risk is placing the picker branch too late in `updateMain`. If that happens, the feature will render but arrows will switch tabs, `N` may start a new idea, and Enter may submit the wrong thing.

The second risk is stale `answerQID`. Reusing `composing` is correct only if every cancel, submit, and run-context transition clears answer state consistently.

The third risk is a scroll/window off-by-one when the filtered list shrinks. The implementation should clamp `Index` before selection and before rendering, and calculate the visible window from the clamped index.
