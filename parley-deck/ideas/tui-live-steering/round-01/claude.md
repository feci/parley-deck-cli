---
agent: claude
idea: tui-live-steering
round: 1
date: 2026-06-06
---

## Summary

Three features, three sizes. Feature 1 (autocomplete) is a contained TUI addition.
Feature 2 (kill an agent) is a small runner+app change: a per-agent cancel registry +
one injected seam. Feature 3 (steer→reply) is the real work: because agents are
one-shot, a steer must spawn a fresh **single-agent attempt** and stream its reply into
the agent's tab. The unifying principle: keep `internal/tui` decoupled — every new
capability arrives as an injected function on `LiveOptions` (exactly like `Start`/`Cancel`
today), and the runner owns process lifecycle.

## Proposed approach

### Feature 1 — slash-command autocomplete (pure TUI)
Reuse the **picker** machinery with a third kind `pickerCommand`, but trigger it
*inline* rather than only on a bare command:
- As soon as `inputText` starts with `/`, compute the matching commands (prefix match on
  `commandSpecs` — a new table `{name, args, help}` that also feeds `/help`) and render a
  compact suggestion list above the input row (NOT the full modal picker; a slimmer
  variant so normal typing/Enter still works).
- **Tab**: if exactly one match → complete it (+ trailing space if it takes args); if
  several → complete the longest common prefix; if the completed command takes a target
  (`/open`, `/answer`) and has no arg yet, pressing Tab again (or Enter) drops into the
  existing bare-command picker. ↑/↓ move the highlighted suggestion, Enter accepts it.
- This is a *suggestion* sub-mode, distinct from the modal `picker`: it never swallows
  printable keys (you keep typing the command), it only owns Tab + (when shown) ↑/↓.
  Esc/empty input dismisses it. One source of truth `commandSpecs` keeps `/help`, the
  suggestions, and the hints in sync.

### Feature 2 — kill one agent (runner + app + TUI)
- Runner: in `runAgent`, create a per-agent `ctx, cancel := context.WithCancel(parent)`
  and register `cancel` in a mutex-guarded `map[agentID]context.CancelFunc` on the run
  handle; deregister on exit. Add `KillAgent(agentID string)` that looks up and calls the
  cancel, and records an `agent.killed` event so the projector lands the agent in a
  terminal `killed`/`failed` state (distinct badge). The round's `WaitGroup` still
  completes because the killed `runAgent` returns on ctx cancel.
- App: expose `KillAgent` on the run handle and pass it onto `LiveOptions.KillAgent
  func(agentID string) error`.
- TUI: on an agent tab, a key (proposed `ctrl+k`) asks "kill <agent>? y/n" (a tiny
  confirm sub-mode so it's never accidental), then calls `opts.KillAgent`. Run-wide
  ctrl+c is unchanged.
- The driver must treat `killed` like a terminal outcome (not "still running") so
  auto-drive doesn't wait forever — but auto-drive isn't active in interactive `parley
  tui`, so this is mostly a projector/state concern.

### Feature 3 — steer→agent→reply (the crux)
Model: **a steer to an agent spawns a fresh single-agent attempt; its stdout is the
reply, surfaced in that agent's tab.** Concretely:
- New runner entry `RunSteerAttempt(ctx, opts, agentID, steerText) (Result, error)`:
  opens a new segment (`steer` round label), builds the prompt = a steer template
  (steer text + the agent's latest artifact for the open idea + a short transcript tail
  for context), invokes via `CommandFor`/`runAgent` into a fresh per-attempt dir, and
  captures stdout. Emits `steer.requested` (kept) → `agent.started` → `agent.finished`
  for the attempt → a `steer.replied` event carrying the reply path.
- App seam: `LiveOptions.Steer func(agentID, text string) (replyStdoutPath string, done
  <-chan struct{}, err error)` (mirrors `Start`). `submitSteer` calls it instead of only
  recording; the returned stdout path is tailed into a transcript buffer for that agent.
- Reply surfacing: render the steer attempt in the SAME agent tab. Simplest faithful
  option: the agent tab shows the attempt's stdout buffer (a clear "── steer reply ──"
  divider), with a "agy is replying…" status while the attempt runs. Keep the original
  round stdout too (switchable) — but v1 can just append the steer attempt as the newest
  content in that agent's buffer.
- Concurrency: if the agent is currently running its round, **queue** the steer and run
  it when that attempt finishes (one attempt per agent at a time; a per-agent steer queue
  of depth 1 is enough). If the run already ended, a steer still runs as a standalone
  attempt (it just re-invokes the agent with context) — useful for "ask agy a follow-up".

## Concerns / open questions
- **Scope/risk**: Feature 3 touches the runner (the riskiest code). I lean toward keeping
  the steer attempt path fully separate from the normal round goroutines (its own segment
  + dir) so it cannot corrupt round artifacts. Confirm with codex.
- **What context to feed the steer** — too much context = slow/expensive; too little =
  useless reply. Proposal: steer text + the agent's own latest artifact + last ~50
  transcript lines. agy/codex weigh in.
- **Reply rendering** — same-tab append vs a dedicated conversation pane. I prefer
  same-tab (least new surface); agy decides the UX.
- **Kill key choice** — `ctrl+k` vs `K`-on-agent-tab; confirm no collision (hermes).

## Risks
- A steer attempt clobbering the agent's round artifact if it writes to the same path —
  mitigate with a separate per-attempt dir; never reuse the round artifact path.
- Per-agent kill racing normal completion (cancel after the process already exited) —
  guard with the registry + a "first writer wins" terminal-state rule.
- Autocomplete suggestion mode fighting the modal picker or normal Enter — keep it a
  non-modal overlay that only claims Tab and (when visible) ↑/↓.
- Decoupling: if the steer/kill seams leak runner types into `internal/tui`, the layering
  breaks — keep them plain func fields on `LiveOptions`.
