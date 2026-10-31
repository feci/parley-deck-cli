---
idea: meta-protocol-change-participants-optional
author: user
created: 2026-10-31
participants: [claude-1]
implementer: claude-1
track: deliberation
auto_implement: false
require_model_diversity: false
checks: go test ./... -count=1 -timeout 45m
excluded: [codex-1 — no provider tokens (429 Too Many Requests at the relay probe; owner statement) — confirmed 2026-10-31, kimi-1 — no provider tokens (no answer to a one-word probe in 120 s; owner statement) — confirmed 2026-10-31, zcode-1 — no provider tokens (owner statement) — confirmed 2026-10-31]
status: final
---

## Problem / idea

Make participants optional. When an agent has no tokens, exclude it at once. While an organizer or an
implementer exists, the idea continues.

## Owner direction

The owner set this rule on 2026-10-31 in another session (project librade-chart-analytics, file
`librade-chart-analytics/docs/decisions/2026-10-31-parley-participants-optional.md`). A peer session relayed it:

> "ak hocijaky agent nema tokeny hned ho vyluc, pokial existuje organizer alebo implementer, tak moze parley
> pokracovat, participanti su volitelny"

Translation: "If any agent has no tokens, exclude it at once. While an organizer or an implementer exists,
Parley can continue. Participants are optional."

The relay asked whether to make this rule a protocol change. The owner selected "Áno, nová idea (Recommended)"
("Yes, a new idea"): the floor falls to the organizer or the implementer, participants become optional, and the
release is CLI 1.57.0 and skill 2.21.0 on all channels.

## User direction (solo exception)

The relay probed the quorum on 2026-10-31. codex-1 returned `429 Too Many Requests`. kimi-1 gave no answer in
120 s. The owner then wrote:

> "ani codex ani kimi ani zcode uz nemaju tokeny takze to sprav sam"

Translation: "Neither codex nor kimi nor zcode has tokens any more, so do it yourself."

This is an explicit user-authorized solo exception under §1. claude-1 (Claude Code, `claude-opus-5-5[1m]`)
designs, implements and releases alone. No independent review exists. The close records
`independent review: none`. See `inbox/claude-1-to-all_meta-protocol-change-participants-optional_solo-exception.md`.

## Constraints

- Keep saved v1 and v2 policies unchanged. Keep the no-retry exhaustion rule and the owner succession order.
- Keep the 20,000-byte skill limit and the 70,000-byte phase-1 packet guardrail.
- A missing review is never a pass (§15).

## Non-goals

- Ideas that start with one participant by choice. They still need the §1 solo exception.
- A machine exhaustion cache or shared-account inference.
