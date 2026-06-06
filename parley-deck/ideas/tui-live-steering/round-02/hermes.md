---
agent: hermes
idea: tui-live-steering
round: 2
date: 2026-06-06
---

## Cross-review

### @claude
I accept the override to a dedicated non-modal `suggest` sub-mode (own `suggest`/`suggestIndex` state, never touching `pickerState`). The modal picker clears `inputText` on open and owns typing as a filter; autocomplete must preserve live editable slash input, so the separation is the correct call. The synthesis resolves every fork I raised in round-01.

### @codex
Adopt your `Handle` design, `KillAgent`/`RunSteerAttempt` signatures, per-steer dirs, `steer.reply_*` events, and the `activateRun` copy requirement verbatim. The segment-ID mutex guard for steers and the explicit scope note (kill/steer only on the async `Handle` path used by `parley tui`; fix-up/impl phases out of v1) close the races I flagged.

### @agy
Adopt the slim floating menu, `ctrl+k` + `y/N` confirm, cyan `steer <id> › ` prefix, inline reply append with divider, and spinner. Your UX details plus the synthesis refinements produce a collision-free, discoverable flow.

## Counter-proposals
None. The dedicated `suggest` mode, conditional-Tab rule (Tab/shift+tab only when input not slash-prefixed; ←/→ always switch), modal `confirmKillAgentID` (blocks every other key), depth-1 queue with rejection on second steer, kill-vs-finish race guard (registry check before cancel + "first writer wins"), handle-locked segment ID, and late post-run steer allowance are all unambiguous and race-free under the keymap + concurrency lens.

## Confirmed for FINAL
- Full key table remains collision-free with dedicated `suggest` sub-mode: normal mode, picker, suggest (only Tab/↑/↓/Enter/esc/backspace), confirm-kill, and answer-compose are cleanly partitioned.
- Conditional-Tab, modal confirm-kill, depth-1 queue, segment locking, and post-run steer rules require no further clarification.
- Implementation cautions only: (1) `suggest` state must be cleared on any non-slash input change or Esc; (2) the confirm sub-mode must be the highest priority interceptor in `updateMain`; (3) per-steer dir creation must be atomic under the handle mutex before spawning the attempt goroutine.

## Remaining risks
None at the contract level. The only residual implementation risk is ensuring `activateRun` copy test actually exercises a Home-launched run (easy to miss in unit tests). All other races (kill/finish, dual-steer, segment ID, log clobber) are closed by the locked registry + per-steer paths.