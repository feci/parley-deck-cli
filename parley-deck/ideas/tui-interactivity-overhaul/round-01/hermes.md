---
agent: hermes
idea: tui-interactivity-overhaul
round: 1
date: 2026-06-04
---

## Summary
Current TUI (live.go:576-594 tailLogFile, ProjectEvents in runstate) is a passive 6-line stdout tail + global-last-state projection. Lacks Codex/Claude-Code parity on: scrollable per-agent transcript, segment-scoped agent state (sticky [FINISHED] post-resume), mid-run composer for steering, and discoverable keymap. Minimum viable interactive set = focus viewport + follow + segment reset + typed intents; gold-plating (thoughts.log, fsnotify) deferred.

## Proposed approach
- **Per-agent focus (live.go)**: ViewAgentFocus state; bubbles/viewport on full stdout.log (not 4KB tail). enter/esc toggle; f= follow (auto G on new bytes), manual scroll drops it. Matches `less +F` feel but no live "thoughts" interleaving (acp.go:207-212 discards thoughtBuf).
- **Segment state (runstate/runstate.go:49-60 AgentState)**: Key state by (agent, currentRound/segment). Emit round.started on continue to reset pending; ProjectEvents scopes lookup so prior finished never leaks. Fixes sticky badge without touching events.jsonl append contract (store/events.go).
- **Composer (app.go:39-50)**: i opens bar producing {intent: answer|steer|new-round}. Reconciles one-shot: default = fresh seeded round (runner.go:231 skip-on-artifact); resume only if CLI supports (--resume) + explicit opt-in. Emits user.steer event, mirrors to round inbox; never bypasses HITL gates.
- **Keymap / model**: Overview (table) -> AgentFocus (viewport) -> Compose. ? overlay; tab cycle; esc back; q/ctrl-c unchanged. Preserves --no-tui path. No change to runner/acp contract.

Trade-off: full log read uses offset-increment (cheap) vs ring-buffer cap for 10MB+ runs. Minimum viable stops at slices 1-2 (visibility + badge); steering semantics require codex/claude input.

## Concerns / open questions
- Resume vs one-shot: does steering mean hidden continuity (Codex does this) or explicit new-round seed? Latter safer for protocol; former only for claude --resume.
- Segment key: round-NN label or new cursor? Ambiguous on continue/resume paths.
- Does focus view need stderr+thoughts tabs or stdout-only? Parity suggests stdout primary.
- Help overlay discoverability: ? must be persistent footer hint or users never find f/g keys.

## Risks
- State regression on resume/skip projection if segment scoping not isolated.
- Composer becomes gate bypass if intents not routed through existing ActionRunner.
- Latency: 250ms poll (live.go) + full-log on every focus tick feels laggy vs Codex instant; fsnotify required for real parity but adds complexity.
- Scope: 4 problems tempt 5 slices; ship visibility+state first.