---
agent: hermes
idea: tui-claude-cli-layout
round: 1
date: 2026-06-04
---

## Summary
Attack the proposed Claude-CLI layout from real usage parity. The 1.12.0 live.go already has liveMode (lines 42-53), focus viewport reads (loadFocusTail/readAppendedLines/capFocusLines at 84-92, bounded 20000 lines/4MiB), steer.Submit (steer.go:5-9), and hitl questions. The redesign must preserve these exactly; any change to key routing or input will break the segment fix and queued steers.

## Proposed approach
Claude Code / Codex CLI parity: top tabs (agents + Status), main = active transcript (reuse focus), bottom = persistent input + status line. Tabs via ←/→ (non-printable, never conflict). Input always live: printable chars type; Enter = steer active agent (steer.go) or answer hitl if open for that agent. Slash-commands (/help, /status, /follow, /deck, /quit) replace all single-letter globals from live.go:48-53. Per-agent buffers keep scroll/follow state so tab switch is instant (no reload). Status tab = old dashboard verbatim. Status line shows run/idea/segment + q:N + follow + queued.

## Concerns / open questions
- KEY-ROUTING trap (live.go always-on input): letters can't be hotkeys. Slash-commands feel most Claude-CLI-authentic; empty-input heuristic or mode-toggle breaks when user expects to type immediately. Modifier keys (ctrl+letter) feel alien in Claude CLI. Minimum that feels right: arrows+PgUp/PgDn+Tab for nav, / for commands, esc clears input, ctrl+c cancels.
- Tabs vs multiplex: showing ONE transcript at a time (tabs) helps focus ("see what THIS agent generates") but owner loses simultaneous visibility; split would require gold-plating and regress narrow-terminal handling in layout.go:26-31. Argue for tabs + quick ←/→ switch.
- Arrow keys: users instinctively expect arrows to move input cursor, not tabs. Conflict is real. Proposed: arrows ALWAYS = tabs (steers are short, backspace/ctrl+u suffice); Tab/shift+tab as safe alternates. Instinct wins over purity.
- HITL answer vs steer in same box: confusing if not disambiguated. Banner above status line ("? question — Enter submits answer") + q:N counter makes it honest; Enter routes based on hitl state for active agent only.
- Minimum viable: cut everything beyond tabs+transcript+input+slash routing. No per-agent LRU, no compact-mode special cases, no new events. Gold-plating to avoid: vim modes, mouse support, configurable hotkeys.
- Regression risk on 1.12.0: segment fix, steer queue, help overlay must survive. Keep steer pkg and runstate untouched; liveMode extended only for activeTab (not modeCompose/modeAnswerQuestion collapse).

## Risks
- UX breakage if arrows steal cursor (users type then can't edit). Mitigate by documenting "arrows switch tabs; input is append-only".
- Discoverability of /commands drops without old ? key. Input placeholder + `/` hint list required.
- Memory: N*4MiB buffers (live.go:38-39) ok for 2-4 agents but must cap explicitly.
- Narrow terminals: tab strip truncates (layout.go:26); test 80-col edge.