---
idea: tui-init-bootstrap
status: final
author: codex
consensus-date: 2026-05-14
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

### Goal

Allow `parley tui --dir <path>` to be the first user-facing entry point for a Parley Deck workspace. If the workspace is not initialized, the TUI asks the user to initialize it, performs the same initialization as `parley init`, and then shows the normal dashboard.

### Scope

- Change only the `parley tui` bootstrap path.
- Add a minimal setup screen for missing workspaces.
- Reuse existing `protocol.InitWorkspace` and `protocol.ReadWorkspaceStatus`.
- Keep the current initialized dashboard path unchanged.
- Add focused tests for setup rendering, confirmation, successful transition, and failure retry behavior.

### Implementation details

- In `internal/app/app.go`, update `runTUIViewWithDiscovery` so `os.ErrNotExist` from `protocol.ReadWorkspaceStatus(root)` calls a TUI bootstrap runner instead of printing `no parley-deck workspace found; run parley init first`.
- Add a setup model or setup state in `internal/tui/app.go`.
- The setup view should show the target workspace path, explain that Parley Deck is not initialized there, and expose explicit keys such as `i`/`enter` to initialize and `q`/`esc`/`ctrl+c` to quit.
- On initialization keypress, call `protocol.InitWorkspace(root)`, then `protocol.ReadWorkspaceStatus(root)`.
- On success, render the existing normal TUI dashboard using the loaded `WorkspaceStatus` and discovered agent list.
- On failure, keep the setup view active, show the error inline, and allow retry.
- Guard against duplicate initialization by ignoring further initialization keys while an initialization command is in flight.

### Tests

- Unit-test that the setup view appears for an uninitialized target and includes the target path plus initialization and quit keys.
- Unit-test the confirmation path with a temp directory: after the init command runs, `parley-deck/COOPERATION.md` exists and the model renders normal dashboard content.
- Unit-test initialization failure: the model stays on the setup screen and renders the error.
- Run existing TUI and app tests to verify initialized-workspace behavior is unchanged.

### Non-goals

- Do not add transport selection, agent selection, or first-run onboarding beyond the initialization prompt.
- Do not auto-initialize without a user keypress.
- Do not change `parley init`, `parley run`, `parley resume`, or workspace layout behavior.
- Do not change agent runtime configuration behavior.

### Verification

- `go test ./...`
- Manual smoke: `parley tui --dir "$(mktemp -d)"` opens the setup screen; pressing the initialize key creates the workspace and shows the dashboard.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
