---
idea: tui-init-bootstrap
drafted-by: codex
date: 2026-05-14
---

## Agreed decisions

- `parley tui --dir <path>` may start before the target directory contains a Parley Deck workspace.
- Only the `parley tui` path should gain bootstrap behavior. `parley init`, `parley run`, `parley resume`, and other non-TUI commands keep their existing initialization requirements.
- When `protocol.ReadWorkspaceStatus(root)` reports a missing workspace for `parley tui`, the TUI should show a setup screen instead of exiting.
- The setup screen must show the target workspace path, provide a clear explicit initialize action, and retain a quit path.
- Initialization must reuse existing `protocol.InitWorkspace(root)` semantics, then reload with `protocol.ReadWorkspaceStatus(root)` and transition to the normal dashboard.
- Initialization failures should stay inside the setup screen, show an inline error, and allow retry without automatically re-running initialization.
- The already-initialized path should keep using the existing dashboard behavior.

## Agreed trade-offs

- Use a narrow setup screen, not a broader onboarding wizard.
- Prefer a simple explicit key such as `i` and/or `enter` for initialization, with `q`, `esc`, and `ctrl+c` for quit.
- Reuse already-discovered agent information after initialization unless implementation reveals a concrete need to rerun discovery.
- Keep tests at the model/state level where possible because full Bubble Tea terminal tests are brittle.

## Open items deferred to implementation

- Decide whether the setup screen displays the raw `--dir` value or an absolute path. It must not display a path different from the path passed to `InitWorkspace`.
- Choose the exact key hint text and error line copy during implementation.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-14
Status: ✅ ACCEPT
Notes: Consensus matches round-01 convergence; implementation should stay scoped to parley tui bootstrap.

### Signoff: claude — 2026-05-14
Status: ✅ ACCEPT
Notes: Aligns with my round-01 proposal — two-state model, explicit confirm, reuse InitWorkspace/ReadWorkspaceStatus, retry-in-place on failure, initialized path unchanged. Deferred items (raw vs absolute path display, exact key/error copy) are appropriate to settle in implementation as long as the displayed path equals the path passed to InitWorkspace.

### Signoff: gemini — 2026-05-14
Status: ✅ ACCEPT
Notes: The consensus captures the core requirements of my round-01 proposal: deferring the workspace check, providing a clear setup screen with an explicit 'i' initialization action, and reusing existing protocol semantics. The trade-offs and deferred items are reasonable and focus the scope on a clean TUI bootstrap experience.

### Signoff: hermes — 2026-05-14
Status: ✅ ACCEPT
Notes: Consensus accurately reflects the bootstrap requirements for parley tui while preserving all constraints and non-goals from the prompt; implementation scope is appropriately narrow.
