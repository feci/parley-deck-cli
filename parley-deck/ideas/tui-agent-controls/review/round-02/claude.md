---
agent: claude
idea: tui-agent-controls
review-round: 2
date: 2026-05-14
reviewed-commit: f5d5e45
responding-to: [claude/review/round-01, gemini/review/round-01, hermes/review/round-01]
---

## Position changes since prior review round

All five MINOR/NIT findings I raised in round-01 have been addressed in fix-up cycle 1. I have no new findings against the implementation at f5d5e45.

Verification against the agreed fixes in `review/consensus.md`:

- **Stale `initModel.View()` post-init branch removed.** `internal/tui/app.go:187-218` no longer contains the `if m.status != nil { return model{...}.View() }` shortcut. The only path that yields a dashboard view is `initModel.Update`'s success branch at `internal/tui/app.go:180-182`, which constructs via `newModel` and preserves width. Fix applied.
- **Symmetric headless command rendering.** `headlessCommandLine` at `internal/tui/app.go:442-453` now falls through `HeadlessMode` → `HeadlessArgs` → `Commands[0]` → `agents.CLIDefault`, so an agent configured only via `Commands` (e.g. the `claude` fixture) still renders `headless: claude` and parity with the `interactive:` line is restored. Fix applied.
- **No trailing space on interactive line.** `interactiveCommandLine` at `internal/tui/app.go:455-461` only appends the joined args when `len(InteractiveArgs) > 0`, eliminating the dangling space when args are empty. `TestDashboardRendersFallbackCommandDetails` at `internal/tui/app_test.go:141-143` explicitly asserts `"interactive: claude \n"` is absent. Fix applied.
- **Backend fallback is `unknown`.** `backendOrUnknown` at `internal/tui/app.go:435-440` returns `agents.ExternalUnknown` ("unknown") for empty `ExternalBackend`, no longer leaking `cli-default` through `valueOrDefault`. Asserted by `internal/tui/app_test.go:134` (`backend: unknown`). Fix applied.
- **Footer extended.** `renderFooter` at `internal/tui/app.go:328` now reads `tab/shift+tab focus  j/k/up/down select  h/i/m set agent mode  x clear mode  q/esc/ctrl+c quit`, covering every key the dashboard `Update` handler accepts. Fix applied.
- **Test coverage for arrow keys, `h`/`m`, and focus gating.** `TestDashboardAgentNavigationClamps` at `internal/tui/app_test.go:163-172` drives `tea.KeyUp` and `tea.KeyDown`. `TestDashboardLaunchModeOverridesAreSessionOnly` at `internal/tui/app_test.go:208-217` exercises `m` → `LaunchManual` and `h` → `LaunchHeadless` paths in addition to the prior `i` path. `TestDashboardModeKeysNoopOutsideAgentFocus` at `internal/tui/app_test.go:241-254` switches focus to ideas and asserts that `i` does not create an override. Fixes applied.

## Responses to other reviewers

### @gemini

- **Finding 1 (headless rendering symmetry):** Same as my MINOR. Addressed by `headlessCommandLine` in `internal/tui/app.go:442-453`.
- **Finding 2 (redundant `initModel.View` branch):** Same as my MINOR. Addressed by the removal verified above.
- **Finding 3 (`cli-default` backend fallback):** Same as my NIT. Addressed by `backendOrUnknown` using `agents.ExternalUnknown`.
- **Finding 4 (`shift+tab` does not reverse):** Dismissed in `review/consensus.md` ("acceptable for the current two-focus-zone dashboard"). With two zones, forward and reverse are identical; revisit if a third zone lands.
- **Finding 5 (idea selection has no detail):** Listed under deferred follow-ups in `review/consensus.md` ("Add selected-idea details or actions in a future dashboard slice"). Out of scope for this slice per FINAL.md.

### @hermes

- **Finding (agent pane focus markers):** Dismissed in `review/consensus.md` — the `selectionMarker` helper at `internal/tui/app.go:331-339` already renders `>` for focused-selected and `.` for unfocused-selected, so the focus signal is glyph-based and survives monochrome terminals. Confirmed in the current code.
- **Finding (session-only label repetition):** Dismissed in `review/consensus.md`. The `effective:` line at `internal/tui/app.go:295-298` appends `(session only)` whenever an override exists, which is the single canonical place to surface that qualifier; repeating it is unnecessary noise.

## Updated findings

None. The implementation at f5d5e45 satisfies FINAL.md, applies every agreed fix from `review/consensus.md`, and the test suite at `internal/tui/app_test.go` now covers the previously-missing arrow-key, `h`/`m`, and focus-gated cases. The dashboard still does not import or reach any process-launching path (`exec.Command`, `runner.CommandFor`), and overrides remain a session-only `map[string]string` populated through `setSelectedAgentMode`/`clearSelectedAgentMode`. Ready to merge.

## Open questions

None.
