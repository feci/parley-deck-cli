---
agent: hermes
idea: tui-command-picker
round: 2
date: 2026-06-06
---

## Cross-review (@claude / @codex / @agy)

@claude: Your round-02 synthesis correctly locks the picker branch immediately after ctrl+c and adopts single-esc-cancels. That placement plus the explicit N guard closes the two risks I flagged in round-01; the re-entrancy test ("open from attached run → esc → still attached") and the "N during filter never starts new idea" test are now required. Your decision to ignore PgUp/PgDn/ctrl+u/d/Home/End while active matches my recommendation exactly and keeps the interaction model unambiguous for v1.

@codex: The early `if m.picker.Active { return m.updatePicker(msg) }` after ctrl+c is the single most important safeguard. Your updatePicker owns only ↑/↓/Enter/esc/backspace/printable and the default case does nothing; that is collision-free. The redundant `&& !m.picker.Active` in the N guard is good documentation even if the early branch already protects it.

@agy: Single-esc-cancels is the right call (I agree with claude's overruling). Your two-step esc would have introduced exactly the ambiguity I warned about in round-01. The height-cap + viewport scroll you proposed is also adopted and orthogonal to key routing.

## Full key-collision sweep (table: key → behavior when picker active vs inactive)

All existing keys in updateMain are covered. Picker branch sits first (post-ctrl+c), so every listed key has an unambiguous defined behavior.

| Key                  | Picker active                          | Inactive (today)                     |
|----------------------|----------------------------------------|--------------------------------------|
| ctrl+c               | global quit (unchanged)                | global quit                          |
| esc                  | cancel picker, restore prior state     | clear input / detach                 |
| N                    | append 'N' to filter (never new-idea)  | new-idea composer (if !composing && empty input) |
| ↑ / left / shift+tab | move selection in filtered list        | switch tab left                      |
| ↓ / right / tab      | move selection in filtered list        | switch tab right                     |
| Enter                | confirm selection → dispatch by Kind   | submitInput / steer / answer         |
| printable rune       | append to filter, reset index          | append to inputText                  |
| backspace / ctrl+h   | pop last filter rune                   | pop last inputText rune              |
| /                    | append '/' to filter (never dispatch)  | dispatch slash command               |
| PgUp / PgDn          | no-op (ignored)                        | scroll transcript                    |
| ctrl+u / ctrl+d      | no-op (ignored)                        | scroll transcript                    |
| Home / End           | no-op (ignored)                        | scroll transcript                    |
| shift+↑ / shift+↓    | no-op (ignored)                        | scroll transcript                    |
| any other key        | no-op                                  | (whatever it did)                    |

Single-esc-cancels is collision-safe: esc's only prior meaning while composing or attached is "back out one level"; the picker is exactly one more level, so esc always exits the current modal layer. No other key had dual meaning that the early branch does not resolve.

Picker-branch-first placement closes every collision: no existing handler can ever see a key that the picker should own, and every non-picker key either keeps its global meaning (ctrl+c) or is explicitly documented as no-op.

No key remains ambiguous.

## Confirmed for FINAL

- Picker branch immediately after ctrl+c, before every other case.
- Single esc always cancels the entire picker (no two-step filter-clear).
- All scroll keys (PgUp/PgDn/ctrl+u/d/Home/End/shift+arrows) are no-op while picker active.
- N during filter appends to filter; the explicit `&& !m.picker.Active` guard stays as documentation.
- Re-entrancy test and N-filter test are required.

## Remaining risks

None on the key-routing / interaction-model axis. The only remaining implementation risk is ensuring every cancel path (esc from composing, activateRun, successful launch/answer) clears answerQID together with composing and picker; that is already called out by claude and covered by the state-machine tests.