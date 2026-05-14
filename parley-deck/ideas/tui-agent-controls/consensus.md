---
idea: tui-agent-controls
drafted-by: codex
date: 2026-05-14
---

## Agreed decisions

- Implement a focused static-dashboard TUI slice for agent navigation, inspection, and session-local launch-mode preview.
- Add dashboard focus state for at least `ideas` and `agents`.
- Add selected row state for ideas and agents, with visible selection markers that work without color.
- Add keyboard navigation on the static dashboard:
  - `tab` and `shift+tab` switch focus.
  - `j`/`down` and `k`/`up` move selection in the focused list.
  - `h`, `i`, and `m` set the selected agent's session launch-mode override to `headless`, `interactive`, or `manual`.
  - `x` clears the selected agent's session override.
  - Existing quit keys continue to work.
- Render selected-agent runtime details using already resolved `agents.Discovery` / `agents.Spec` data, not a second config loader.
- Details must include installed state, version/probe error, configured launch mode, effective launch mode, model/profile/reasoning, sandbox, approval, timeout, backend, home isolation, headless command shape, interactive command shape, prompt mode, invoke strategy, and notes when present.
- Session overrides must never mutate the underlying discovery/spec data and must not be persisted to `agents.local.toml`.
- The static dashboard must not launch hosted or local agent processes in this slice.
- After a fresh TUI initialization succeeds, transition from `initModel` to the real dashboard model so new dashboard keybindings work immediately.

## Agreed trade-offs

- Keep mode switching session-local for now. Persistent editing of `agents.local.toml` is deferred because it needs a separate safe config-writing surface.
- Prefer explicit `h`/`i`/`m` keys over a single cycling key because they make the target mode unambiguous.
- Use substring-focused view tests rather than brittle full-layout snapshots.

## Open items deferred to implementation

- Exact layout can be two panels plus a detail section, or a three-column layout, as long as text remains readable at normal terminal widths.
- Long command/detail lines may be truncated or wrapped; implementation should avoid panics and preserve key information.
- Exact footer wording is implementation-owned, but it must clearly say overrides are session-only.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-14
Status: ✅ ACCEPT
Notes: Consensus keeps the slice useful but safe: inspect and session-preview modes in TUI, no process launch or config persistence.

### Signoff: claude — 2026-05-14
Status: ✅ ACCEPT
Notes: Matches my round-01 proposal in substance — navigable agents list, resolved-runtime detail panel, session-only `h`/`i`/`m` overrides with `x` reset, no process launch, no `agents.local.toml` writes. The explicit `initModel` → dashboard handoff after fresh init is a good addition. Lowercase override keys and `x` reset are fine substitutes for my `H/I/M` + `=`. Implementation should still keep the "preview only, no launcher" guard visible in code so a future hotkey doesn't quietly grow into `exec.Command`.

### Signoff: gemini — 2026-05-14
Status: ✅ ACCEPT
Notes: The consensus successfully incorporates my proposal for dual-pane navigation and agent inspection. I support the decision to use explicit `h`/`i`/`m` keys for overrides instead of cycling, as it reduces ambiguity. The session-local nature of overrides correctly balances utility with safety for this first slice.

### Signoff: hermes — 2026-05-14
Status: ✅ ACCEPT
Notes: Consensus aligns with the original prompt constraints: small focused TUI slice, no process launch, no persistence to agents.local.toml, clear keyboard nav and details panel. Session-local overrides are the right safety boundary.
