---
idea: tui-agent-controls
status: final
author: codex
consensus-date: 2026-05-14
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

### Goal

Make the static `parley tui` dashboard useful before adding a full workflow launcher. Users should be able to move between ideas and configured agents/models, inspect the selected agent's resolved runtime details, and preview non-headless launch-mode choices safely.

### Scope

- Add dashboard focus and selection state for ideas and agents.
- Add keyboard navigation for the static dashboard:
  - `tab` and `shift+tab` switch focus.
  - `j`/`down` and `k`/`up` move the selected row in the focused pane.
  - `h`, `i`, and `m` set the selected agent's session launch-mode override to `headless`, `interactive`, or `manual`.
  - `x` clears the selected agent's session override.
  - Existing quit keys remain unchanged.
- Render a selected-agent detail view using resolved `agents.Discovery` data.
- Keep launch-mode overrides session-local and preview-only.
- Make `initModel` transition to the real dashboard model after successful initialization so dashboard keybindings work immediately.

### Implementation details

- Extend `internal/tui/app.go`'s dashboard model with:
  - focus zone enum/string for `ideas` and `agents`;
  - selected idea and agent indexes;
  - `map[string]string` for session launch-mode overrides.
- Add helper methods to clamp selections, move selection, switch focus, get selected agent, compute effective launch mode, set overrides, and clear overrides.
- Render the agent pane as a selectable list with markers that do not rely on color.
- Render selected-agent details including installed state, version/probe error, configured launch mode, effective launch mode, model/profile/reasoning, sandbox, approval, timeout, backend, home isolation, headless command shape, interactive command shape, prompt mode, invoke strategy, and notes when present.
- The detail view must clearly label overrides as session-only.
- Do not call `exec.Command`, `runner.CommandFor`, or any process-launching path from this dashboard slice.
- After workspace initialization succeeds in `initModel.Update`, return a real dashboard `model` rather than keeping `initModel` and only delegating `View`.

### Tests

- Dashboard view renders selectable agent details for the selected agent.
- Agent navigation with `j`/`k` or arrow keys clamps at list boundaries.
- `tab` and `shift+tab` switch focus while preserving selections.
- `h`/`i`/`m` set only session overrides and do not mutate the underlying discovered agent.
- `x` clears the selected agent override.
- Fresh-init transition returns the real dashboard model so post-init dashboard keybindings work.
- Existing init and dashboard tests continue to pass.

### Non-goals

- Do not persist mode changes to `parley-deck/agents.local.toml`.
- Do not start hosted or local agent processes from the static dashboard.
- Do not build a run creation wizard or workflow launcher.
- Do not embed external CLIs or automate interactive sessions through a PTY.
- Do not change agent runtime schema or billing/accounting semantics.

### Verification

- `go test ./internal/tui`
- `go test ./...`
- Manual smoke: open `parley tui`, use `tab`, `j/k`, `h/i/m`, and `x`, and verify the selected agent detail panel changes without launching agents.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
