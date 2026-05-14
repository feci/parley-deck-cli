---
agent: claude
idea: tui-init-bootstrap
round: 1
date: 2026-05-14
---

## Summary

Allow `parley tui --dir <path>` to launch in an uninitialized directory. Instead of refusing on startup, the TUI detects the missing workspace, renders a dedicated setup screen that shows the resolved target path, and waits for an explicit user confirmation. On confirmation it invokes the same `protocol.InitWorkspace` semantics used by `parley init`, then transitions the running TUI session into the normal dashboard via the existing `protocol.ReadWorkspaceStatus` path. No new commands, no implicit initialization, no changes to `parley init` or the protocol.

## Proposed approach

1. Lift the early "workspace must exist" guard out of the TUI entry path so startup no longer hard-fails on a missing `parley-deck/COOPERATION.md`.
   - Verify: `parley tui --dir <empty-dir>` returns a running TUI process instead of an exit-with-error.
2. Introduce a small bootstrap state in the TUI model with two states: `needsInit` and `initialized`. At startup, resolve the target dir, run `protocol.ReadWorkspaceStatus`, and pick the state from the result.
   - Verify: unit test on the model constructor — given a fake status reader returning "not initialized", model is in `needsInit`; given "initialized", model is in `initialized` and identical to today's behavior.
3. Render a setup view for `needsInit`: show the absolute target path, a one-line explanation, and a single confirm action (e.g. `Enter` to initialize, `q`/`Esc` to quit). No other dashboard widgets in this state.
   - Verify: golden/snapshot test of the setup view; ensure the path and confirm hint are visible.
4. On confirm, call `protocol.InitWorkspace(targetDir)` exactly the same way the `parley init` command path calls it (same args, same error surface). On success, re-run `protocol.ReadWorkspaceStatus` and transition the model to `initialized`, reusing the existing dashboard mount path so the normal TUI takes over in-process.
   - Verify: integration-style test using a temp dir — drive the model with a synthetic "confirm" message, assert `InitWorkspace` was called once with the resolved path, and that the post-init state matches the dashboard's initial state when launched against an already-initialized dir.
5. On init failure, stay in `needsInit`, surface the error inline (one-line error region under the prompt), and keep the confirm action available for retry. Do not auto-retry.
   - Verify: test with an `InitWorkspace` stub that returns an error — model stays in `needsInit`, error string is rendered, no transition occurs.
6. Leave the already-initialized path byte-for-byte equivalent: if `ReadWorkspaceStatus` reports initialized at startup, skip the setup view entirely and mount the dashboard exactly as today.
   - Verify: existing TUI tests for the initialized path continue to pass without modification.

Scope is intentionally one screen, one keybind, two states. No wizard, no transport selection, no agent config — those are explicit non-goals.

## Concerns / open questions

- Confirm action UX: a single `Enter`-to-init keybind is the simplest thing that works, but it can be triggered accidentally. A `y`/`n` prompt is slightly safer but adds a keypress. I'd default to `Enter` to confirm + `Esc`/`q` to quit, matching typical Bubble Tea / TUI confirmation flows, and revisit only if reviewers prefer `y/n`.
- Init-in-progress feedback: `InitWorkspace` is fast in the common case, but if it ever blocks (e.g. slow filesystem, future hook), we should at minimum render a "Initializing…" line and disable the confirm key while the call is in flight. Worth doing now to avoid double-confirm bugs.
- Path resolution semantics: the setup view must display the same absolute path that `InitWorkspace` will write to. We should reuse whatever path-resolution helper `parley init` uses (relative → absolute, `~` expansion, symlink handling) rather than re-implementing it in the TUI, to avoid a mismatch between displayed and actual target.
- Permission / pre-existing partial state: if the target dir is unwritable, or `parley-deck/` exists but is incomplete (e.g. directory present, `COOPERATION.md` missing), what does `InitWorkspace` do today? We should not change that behavior — just surface the resulting error verbatim. Worth confirming the current contract before writing the failure-path test.
- Transport defaults: the prompt requires preserving the active `github-pr` transport. I'm assuming `InitWorkspace` already picks that up from the same source `parley init` uses (config / flag) and we should not introduce a TUI-side override. Calling out so reviewers can confirm the call site matches `parley init` exactly.
- Test layering: prefer driving the model directly (no terminal emulation) for state transitions, and reserve any end-to-end TUI test for the happy-path smoke check. Open to reviewer guidance on existing test conventions in this repo.

## Risks

- Divergence from `parley init`: if the TUI calls `InitWorkspace` with subtly different args (path normalization, flags), users could end up with two slightly different workspace shapes depending on entry point. Mitigation: extract or reuse the exact call site used by the `init` command — don't re-derive arguments in the TUI layer.
- Accidental initialization: a stray keypress on the setup screen could initialize a workspace in the wrong directory. Mitigation: always render the resolved absolute path prominently above the confirm action; require an explicit confirm key (not just any key); quit on `Esc`/`q`.
- Initialized-path regression: refactoring the TUI entry to support two startup states risks breaking the dashboard's existing init sequence. Mitigation: keep the `initialized` branch on the exact same code path as today (mount the same root model/component), and run existing TUI tests unchanged as the regression gate.
- Re-entrancy on retry: if init fails and the user confirms again, we must not leave behind half-written files or double-initialize on a partially-succeeded run. Mitigation: re-run `ReadWorkspaceStatus` before each `InitWorkspace` call; if it now reports initialized, transition to the dashboard instead of re-initializing.
- Scope creep into a wizard: it's tempting to add transport choice, agent selection, or a "what is Parley Deck" panel to the setup screen. Explicit non-goal — keep this round to "ask, init, switch."
