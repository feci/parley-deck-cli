---
idea: meta-protocol-change-devx-speed
type: meta-protocol-change
facilitator: claude-1
participants: [claude-1, codex-1, hermes-1, antigravity-1]
transport: github-pr
created: 2026-07-03
status: complete
track: deliberation
---

# Idea: Complete protocol revision — developer usability + speed

## The ask (from the deck owner)

> Do a **complete review of the whole protocol**. We need to make it **usable for
> developers and similar roles**, and we need to **make it faster** — the whole
> interaction feels like it takes too long. Propose improvements; feel free to draw
> on **modern agentic-AI concepts** that have recently landed in the agentic world.

This is a meta-protocol-change idea per §7. The target is `COOPERATION.md` as a
whole (currently 1046 lines, 15 sections, a 9-phase idea lifecycle Phases 0–8 with
per-phase rounds, review rounds, and fix-up cycles, plus §9.0 readiness ping, §12
pipeline, §13 retro, §14 outer loop).

## Two hard goals

- **G1 — Developer usability.** A developer (or PM/designer/other non-protocol-author
  role) should be able to pick this up and use it *without reading 1000 lines*. Clear
  entry point, obvious defaults, role-appropriate framing, a real quickstart, minimal
  ceremony for the common case. Today the protocol reads as authored *for* protocol
  authors.
- **G2 — Speed.** Cut the end-to-end latency of a *typical* interaction. Be concrete
  about where the wall-clock goes today: sequential independent rounds, mandatory full
  quorum, the design-consensus → finalize → implement → review-consensus → fix-up loop,
  §9.0 liveness pings, deep-reasoning timeouts (20–30 min each), multi-round refutation.

## Hard constraint — do NOT throw away the safety core

Speed must not silently delete the properties that make Parley worth using:
durable audit trail, **non-solo independent verification for risky changes**,
refutation-default review, human brake on automation (§14), no-secret rules. The
right shape is almost certainly **conditional rigor / tiering**, not "delete rounds".

## What each participant must deliver in `round-01/<agent-id>.md`

Write your own independent analysis (do not read others' round-01 first). Cover:

1. **Time-sink audit.** The top 4–6 places the current protocol spends wall-clock or
   human/agent effort out of proportion to the risk it buys. Cite concrete sections/phases.
2. **A tiering model.** Propose named tracks (e.g. a lightweight "dev / fast" track vs
   the full "deliberation" track — invent your own). Give **objective triggers** for which
   track a task lands in (size of change, risk class, reversibility, files touched,
   security surface, whether it's a protocol change, etc.). Say exactly which phases each
   track keeps, collapses, parallelizes, or makes optional.
3. **DevX improvements.** Concrete: the entry point / quickstart, defaults that "just work",
   role-based framing, what to move out of the main doc, how a first-time developer starts
   an idea in <5 minutes.
4. **Speed levers.** Specific mechanics: parallel vs sequential rounds, single-reviewer
   fast-path vs full quorum, collapsing consensus+finalize, auto-advance/driver, streaming,
   timeout tuning, when one agent is genuinely enough.
5. **Modern agentic concepts to adopt.** Map concrete recent ideas (2025–2026) to protocol
   changes — e.g. conditional rigor / right-altitude, lead-agent + subagent orchestration,
   deterministic workflows vs model-driven control flow, plan mode, spec-driven development,
   context engineering / progressive disclosure, parallel worktrees, verification/refutation
   gates, right-sized autonomy, skills, "closing the loop", the bitter lesson (less
   scaffolding). Only propose what genuinely helps G1/G2; name the concept and the exact
   protocol edit it implies.
6. **What MUST stay.** The non-negotiables you would defend against a speed-at-all-costs push.
7. **Prioritized shortlist.** MUST / SHOULD / COULD, each item one line, ordered by
   (impact on G1+G2) ÷ (cost + risk).

Keep it concrete and citeable — we will synthesize a consensus and bring the top
decisions to the owner before writing FINAL.md.
