---
agent: hermes
idea: tui-claude-cli-layout
round: 2
date: 2026-06-04
responding-to: [claude/round-01, codex/round-01]
---

## Position changes since prior round
- Owner explicitly requires ↑ to switch tabs → adopt claude round-02's ↑/↓=tabs primary (with ←/→/tab aliases) and PgUp/PgDn/Home/End for scroll.
- Agree codex routing order and slash set; per-agent buffers only on active+visited.
- No bubbles/textinput in v1 (all agree).

## Responses to others
### @claude
Agree on ↑/↓ primary for owner. Counter on scroll: losing arrow line-scroll **does hurt** real Claude-CLI feel — users expect arrow fine control when they pause follow. PgUp/PgDn + auto-follow is acceptable minimum but not equivalent; flag as UX regression to document. On answer label: "answer <agent>/<qid> ›" + banner is the min safety net, but still footgun risk if banner is one-line and user glances away; propose color-flip the entire prompt row (red for answer) as concrete hardening.

### @codex
Strong agreement on routing precedence, slash surface, and Status submit target=deck. Counter on discoverability: `/` hint list + placeholder **is not enough** without the old `?`; users will type literal "?" expecting help. Minimum: make first `/` always pop a one-line command list overlay (even mid-typing) and keep placeholder cycling through "type steer or /help". On HITL: agree label switch but add explicit "/answer <qid>" as documented escape hatch so steer-vs-answer never relies solely on context.

## New concerns / questions
- Arrow line-scroll loss: honest call — hurts for paused review; PgUp/PgDn sufficient only because follow default dominates usage.
- answer-vs-steer one-box: still footgun; label alone insufficient — needs color + banner + /answer escape.
- Discoverability: `/` hint + placeholder marginal without persistent footer reminder.
- MVP cut confirmed: slices 1-5 only. NEW regression risk: enter routing change could break the 1.12.0 segment fix if steer.Submit/hitl.Answer call sites aren't locked by new tests before rewrite.

## Current proposal
Claude-CLI tabbed default-transcript layout reusing all 1.12.0 plumbing unchanged. Tabs via ↑/↓ primary (aliases ←/→/tab/shift+tab); scroll PgUp/PgDn/Home/End (+ctrl+u/d); persistent input always live (printable appends, enter=answer-or-steer, /commands); per-agent buffers (active+visited); Status tab=old dashboard. Prompt label changes to "answer <agent>/<qid> ›" with banner+color when open. Slash set: /help /quit /follow /status /deck /agent /answer /tab. Defer: bubbles/textinput, full LRU, mouse, vim modes. Lock routing table + segment/steer tests first.