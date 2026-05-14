---
agent: claude
idea: tui-agent-controls
round: 1
date: 2026-05-14
---

## Summary

The dashboard today renders agents as a static list. The smallest useful step is to make that list *navigable*, render a resolved runtime detail panel for the highlighted agent, and let the user override `launch_mode` for the current session only. Persistence to `agents.local.toml` and any process spawning stay out of scope; the dashboard remains a "read + plan" surface, and `parley run` continues to be the only thing that actually starts an agent. This keeps the slice safe (no paid hosted process can leak out of an arrow-key press) while still answering the user's real complaint: "I can't move around or see what would happen."

## Proposed approach

Three concrete pieces, all inside the existing Bubble Tea model:

1. **Focus + selection state.** Add `focusZone` (`ideas | agents`) and per-zone selection indices to the dashboard model. `tab` / `shift+tab` switches `focusZone`. Inside the focused zone, `j`/`down` and `k`/`up` move the selection, clamped to the slice length. The non-focused zone keeps its last selection but renders dimmed. Selection marker is a single leading glyph (e.g. `>` for focused, `·` for unfocused-but-selected) so it works without color.

2. **Resolved-runtime detail panel.** When `focusZone == agents`, render a right-hand panel showing the highlighted agent's resolved values: `launch_mode`, `command` (joined for display), `prompt_mode`, `invoke_strategy`, `model` / `profile`, `sandbox`, `approval`, `timeout`, `backend`, and `notes`. Resolution should reuse whatever loader `parley agents list` uses today — do *not* re-implement merging. The panel reads from an in-memory `runtimeView` struct populated once at TUI start (and re-resolved on workspace reload, if that path already exists).

3. **Session-only mode override.** Keys `H` / `I` / `M` set the highlighted agent's *displayed* `launch_mode` to `headless` / `interactive` / `manual` for this TUI session. The override lives in a `map[agentID]launchMode` on the model; the detail panel reads override-first, falling back to the resolved value. A small footer hint (`H headless · I interactive · M manual · = reset`) and an `=` reset key keep this discoverable. Crucially, the dashboard **does not act on these overrides** — it only displays them, and prints them in the footer status line as "would launch as: interactive" so the user can copy that intent into a real `parley run` invocation.

For non-headless modes, add a single explicit guard: any future "launch from dashboard" affordance is gated behind a separate code path that does not exist yet in this slice. Document this in a top-of-file comment in the dashboard model so the next contributor doesn't accidentally wire a hotkey to `exec.Command`.

Tests, all table-driven against the model's `Update`:
- navigation clamps at 0 and `len-1` in both zones
- `tab` cycles focus and preserves per-zone selection
- `H`/`I`/`M` mutate only the override map, never the underlying agent struct
- `=` clears the override
- rendered detail panel string contains the override value when set, the resolved value otherwise

## Concerns / open questions

- **Where does runtime resolution actually live?** If `parley agents list` resolves through a helper that returns a struct, reuse it. If resolution is inlined in the command, the slice needs a small extract first — but that extract should be its own commit so the TUI diff stays surgical. I'd rather pay that cost than duplicate the merge logic in the dashboard.
- **`command` rendering.** Some agents will have long argv with secrets-adjacent flags (model names, profiles). Joining with spaces is fine for display, but we should not shell-quote in a way that suggests "copy this and run it" — it's a *preview*, not a launcher. A `# preview only` prefix or a panel header that says "resolved configuration" should be enough.
- **Manual mode semantics.** The prompt lists `manual` as a launch mode but doesn't say what it means operationally. For this slice I'm treating it as "user will start the agent themselves outside parley"; the detail panel just shows the resolved command so they can copy it. If `manual` has stronger semantics elsewhere in the codebase, the detail panel copy should match.
- **Workspace-not-initialized path.** The kickoff says initialized-workspace and newly-initialized TUI behavior must keep working. If the agents list is empty (fresh init), the agents zone should render a single "no agents configured" row and ignore navigation keys rather than panicking on a nil slice. Worth a dedicated test.
- **Color vs. monochrome.** Bubble Tea on some terminals has no color. The selection marker must be a glyph, not just a foreground color, or users on basic terminals lose the focus indicator entirely.

## Risks

- **Scope creep into a launcher.** The moment we render "this is what would run", someone will want a key to actually run it. Resist that here; it's a separate PR with separate review concerns (process lifecycle, output streaming, cancellation, billing context for hosted agents). The dashboard guard comment is the cheap defense.
- **Stale resolved view.** If the user edits `agents.local.toml` while the TUI is open, the detail panel will lie. Acceptable for this slice — document it in the footer hint or accept it silently — but call it out so we don't ship a "reload" key by accident and grow the surface.
- **Override confusion.** Users may expect `H`/`I`/`M` to persist. The footer must make "session only" obvious, and the override should clear on workspace reload. Otherwise we'll get bug reports that "parley forgot my mode" after restart, which is exactly the persistence question we're deferring.
- **Key collisions.** `j`/`k` are common Vim-style nav, but if the TUI already binds them elsewhere (e.g. inside the live-run view), the dashboard binding must be scoped to the dashboard model's `Update` only. Verify before wiring.
- **Test brittleness on rendered strings.** Asserting substrings in rendered output is fine for this slice but couples tests to layout. Prefer asserting on model state where possible, and use rendered-string checks only for the detail panel's override-vs-resolved branch.
