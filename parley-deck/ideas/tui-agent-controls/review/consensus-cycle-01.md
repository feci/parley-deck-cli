---
idea: tui-agent-controls
cycle: 1
drafted-by: codex
date: 2026-05-14
reviewed-commit: f4315d4
---

## Agreed fixes

- Remove the stale `initModel.View()` post-init branch that manually constructs an under-initialized dashboard model.
- Render a headless command line even when an agent only has base `Commands` configured and no explicit `HeadlessMode` or `HeadlessArgs`.
- Avoid trailing spaces in rendered interactive command lines.
- Use an `unknown` backend fallback instead of `cli-default` when `ExternalBackend` is empty.
- Extend footer copy to mention arrow keys, `shift+tab`, `esc`, and `ctrl+c`.
- Add focused tests for arrow-key navigation, `h` and `m` mode overrides, and no-op mode keys while the ideas pane has focus.

## Deferred follow-ups

- Add selected-idea details or actions in a future dashboard slice.
- Add persistent `agents.local.toml` editing in a separate config-writing slice.

## Dismissed findings

- `shift+tab` using the same two-pane toggle as `tab` is acceptable for the current two-focus-zone dashboard.
- Agent pane selection markers already distinguish focused selected rows (`>`) from unfocused selected rows (`.`); no fix required.
- The effective mode line already adds `(session only)` when an override exists; no repeated label is required.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-14
Status: ✅ ACCEPT
Notes: Review consensus captures the agreed minor fix-up items before merge.

### Signoff: claude — 2026-05-14
Status: ✅ ACCEPT
Notes: Agreed fixes cover all my round-01 findings (stale init branch, headless line symmetry, trailing-space, footer copy, and arrow/`h`/`m`/focus-gated test gaps). Deferred and dismissed items match my read of FINAL scope.

### Signoff: gemini — 2026-05-14
Status: ✅ ACCEPT
Notes: The agreed fixes and tests address the primary risks of UI stale state and navigation predictability. The focus on rendering symmetry for headless lines and the addition of specific mode-override tests provide high confidence for this first interactive slice.

### Signoff: hermes — 2026-05-14
Status: ✅ ACCEPT
Notes: All agreed fixes, deferred follow-ups, and dismissed findings precisely match the round-01 review scope and the original prompt constraints for a minimal interactive slice. Ready to merge.
