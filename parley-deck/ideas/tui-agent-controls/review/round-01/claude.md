---
agent: claude
idea: tui-agent-controls
review-round: 1
date: 2026-05-14
reviewed-commit: f4315d4
---

## Summary

The implementation matches FINAL.md scope. Dashboard focus and selection state, `tab`/`shift+tab`, `j/k`+arrow navigation, `h/i/m`/`x` session overrides, the resolved-runtime detail panel, and the `initModel` → dashboard transition are all present in `internal/tui/app.go`. No `exec.Command` / `runner.CommandFor` paths are reachable from the dashboard, overrides are kept in an in-memory `map[string]string`, and the detail line is explicitly tagged `(session only)`. Unit coverage in `internal/tui/app_test.go` exercises rendering, clamping, focus preservation, override set/clear, and the post-init transition. Findings below are quality and test-coverage gaps rather than spec violations.

## Findings

### [MINOR] Unreachable `initModel.View()` post-init branch constructs an under-initialized `model`

`internal/tui/app.go:188-190` keeps the legacy branch:

```go
if m.status != nil {
    return model{status: *m.status, agents: m.agents, width: m.width}.View()
}
```

After FINAL's requirement to "return a real dashboard model" was implemented (`internal/tui/app.go:180-182`), this branch is dead — `Update` swaps to the dashboard model before any subsequent `View()` call. The bigger concern is that the struct literal bypasses `newModel`, so `focus` is `""`, `launchOverrides` is `nil`, and `clampSelections` is never called. If a future change accidentally routes through this path again, the dashboard would silently be in a broken state (no focus zone matches, override map writes would still work because Go allocates on assignment, but the marker logic would render `"."` for everything).

**Fix:** Delete the `if m.status != nil { ... }` branch in `initModel.View()`. The success path already returns a fully-constructed dashboard model via `newModel`.

### [MINOR] Test gaps for spec-required keybindings: arrow keys, `h`/`m`, and focus gating

FINAL explicitly lists `j`/`down` and `k`/`up` as navigation, and `h`/`i`/`m` as the three override keys, but `internal/tui/app_test.go` only covers `j`/`k` (no `tea.KeyDown`/`tea.KeyUp`) and only the `i` (interactive) override path. There is also no test asserting that `setSelectedAgentMode` / `clearSelectedAgentMode` are silently ignored when `m.focus == focusIdeas` — a regression that flipped the focus check to `==` would set overrides while the user is navigating ideas, and nothing would catch it.

**Fix:** Add cases to `TestDashboardAgentNavigationClamps` (or a sibling) that drive `tea.KeyMsg{Type: tea.KeyDown}` and `tea.KeyUp`. Extend `TestDashboardLaunchModeOverridesAreSessionOnly` (or add a peer test) to cover `h` → `LaunchHeadless`, `m` → `LaunchManual`, and a "press `i` while focus is ideas → no override created" assertion.

### [MINOR] `renderAgentDetails` "headless" line is omitted when only `HeadlessArgs`/`HeadlessMode` are unset but `Commands` are present

`internal/tui/app.go:319-323` renders the headless line only if `HeadlessMode != ""` or `len(HeadlessArgs) > 0`. The parallel `interactive` line at `:324` uses `agents.InteractiveCommandOrDefault(agent.Spec)` so it always shows the base command. For an agent configured purely via `Commands` (no explicit headless mode/args), the detail panel hides the headless command shape entirely while still showing `effective: headless` in the mode line — a confusing asymmetry given FINAL calls out "headless command shape" as a required detail field.

**Fix:** When neither `HeadlessMode` nor `HeadlessArgs` are set, render `headless: <Commands[0]>` (or `cli-default` if commands are empty) so the detail panel is symmetric with the interactive line.

### [NIT] `interactive:` line has a trailing space when `InteractiveArgs` is empty

`internal/tui/app.go:324`:

```go
b.WriteString(fmt.Sprintf("interactive: %s %s\n", valueOrDefault(...), strings.Join(agent.InteractiveArgs, " ")))
```

When `InteractiveArgs == nil` (e.g., the `codex` fixture), the joined string is empty and the rendered line becomes `interactive: codex ` with a dangling space before the newline. Cosmetic only.

**Fix:** Guard the args portion, e.g. `cmd := valueOrDefault(...); if len(args) > 0 { cmd += " " + strings.Join(args, " ") }`, then write `cmd + "\n"`.

### [NIT] Footer help string omits arrow keys, `shift+tab`, `esc`, and `ctrl+c`

`renderFooter` (`internal/tui/app.go:336`) advertises `tab focus  j/k select  h/i/m set agent mode  x clear mode  q quit`, but the Update handler also accepts `down`/`up`, `shift+tab`, `esc`, and `ctrl+c`. Discoverability suffers for users who reach for the arrow keys or `Esc` first.

**Fix:** Either extend the footer (e.g., `tab/shift+tab focus  j/k/↑/↓ select … q/esc quit`) or add a `?` help overlay; the former is sufficient for this slice.

## Open questions

- Should pressing `h`/`i`/`m`/`x` while focus is on the ideas pane surface a brief footer hint ("switch focus to agents first") instead of being a silent no-op? The current behavior is correct per FINAL but offers no feedback.
- FINAL is explicit that overrides do not persist and do not feed `parley run`. Is there a follow-up slice planned that will read `launchOverrides` for actual launches, and if so, should this dashboard already namespace the map (e.g., by idea slug) to avoid bleeding overrides across run contexts?
