---
agent: claude
idea: tui-init-bootstrap
round: 1
date: 2026-05-14
reviewed-commit: 20a568c
implementation-pr: https://github.com/feci/parley-deck-cli/pull/16
---

## Findings

No blocking findings.

Non-blocking observations:

- `internal/app/app.go:1196` (`runTUIViewWithDiscovery`): the `os.ErrNotExist` branch now routes to `tui.RunInit(root, results)` instead of returning the "run `parley init` first" message, matching FINAL. Other error classes still surface through `stderr` with a non-zero exit, which preserves prior behavior.
- `internal/tui/app.go:107` (`initModel.Update`): the `i` / `enter` handler is gated by `m.status != nil || m.initializing`, satisfying the FINAL requirement to ignore further initialization keys while one is in flight (and after success). Quit keys (`q`, `esc`, `ctrl+c`) remain wired in both states.
- `internal/tui/app.go:138` (`initModel.View`): on success it re-renders the dashboard via `model{status: *m.status, agents: m.agents, width: m.width}.View()`, so the existing initialized dashboard path is reused unchanged. Width is propagated through `tea.WindowSizeMsg` on `initModel.Update`, so the dashboard reflects later resizes too.
- `internal/tui/app.go:175` (`initWorkspaceCmd`): performs `protocol.InitWorkspace` followed by `protocol.ReadWorkspaceStatus`, surfacing either error back through `initWorkspaceMsg.err`. On failure `errText` is set and `initializing` clears, leaving the setup screen interactive for retry — matches FINAL.
- `internal/tui/app_test.go` covers (a) the setup view including path / init keys / quit keys, (b) the happy-path transition where `COOPERATION.md` is created and the dashboard renders, and (c) the failure-stays-on-setup case using a pre-existing `parley-deck` file to force `MkdirAll` to fail. These are the three test cases enumerated in FINAL.
- `internal/config/runtime.go:53` (`LoadAgentSpecs`) treats both `agents.toml` and `agents.local.toml` as optional, so the pre-init `discoverConfigured` call inside `runTUIView` does not error on a missing workspace; it returns the default specs, keeping the setup screen reachable.
- IMPLEMENTATION.md, FINAL.md, and the head commit (`20a568c` on `impl-tui-init-bootstrap`) declare no deviations and the diff is scoped to the TUI bootstrap path and its tests; no incidental refactoring detected.

## Verification

- Read `parley-deck/ideas/tui-init-bootstrap/FINAL.md` and `IMPLEMENTATION.md`; cross-checked each behavioral requirement (setup view content, explicit init/quit keys, reuse of `protocol.InitWorkspace` + `ReadWorkspaceStatus`, dashboard reuse on success, inline retry on failure, duplicate-init guard, non-goals untouched) against the implementation in `internal/tui/app.go` and `internal/app/app.go`.
- Read the test file `internal/tui/app_test.go` and confirmed each of the three FINAL-required tests is present and asserts the right surface (path + key hints, `COOPERATION.md` creation + dashboard markers, error string + workspace still unreadable).
- Sandbox prevented running `go test ./...` locally; relied on `IMPLEMENTATION.md` and codex round-01 review reports that `go test ./internal/tui` and `go test ./...` pass for `20a568c`.
- No code edits were made.

## Residual risk

- Agent discovery (`discoverConfigured` → version probes) still runs before `ReadWorkspaceStatus`, so a slow / missing agent binary can delay the setup screen even on a fresh workspace. This was acknowledged in codex round-01 and is consistent with FINAL's scope; addressing it would touch the initialized path too.
- The dashboard is rendered from within `initModel.View()` after success, but the active `tea.Model` remains `initModel`. Today the user-visible surface is fine because the only interactive keys on the dashboard are quit keys, which `initModel.Update` handles identically. If the normal dashboard later gains keybindings, this delegation will need to switch to returning the real `model` from `Update` rather than only its `View`.
- Init failure surfaces the raw `error.Error()` string (e.g., `mkdir .../parley-deck: not a directory`). Acceptable for an internal tool, but worth noting for future UX polish.
- The setup screen captures the agents slice at construction time and does not re-run discovery after a successful init. For the current flow this matches the initialized-dashboard contract, but any future change that makes agent discovery depend on workspace contents would need to refresh `m.agents` after the init transition.
