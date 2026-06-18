---
idea: meta-protocol-change-fusion-execplans
title: "Inspiration from OpenRouter Fusion + OpenAI ExecPlans for the parley-deck protocol"
kind: meta-protocol-change
status: open
phase: 1-independent-analysis
transport: github-pr
initiator: claude
participants: [claude, codex, agy, hermes]
created: 2026-06-18
quorum: [claude, codex, agy, hermes]
---

## Task

Two external systems were researched (see `reference/research.md`, which is the
shared evidence corpus — read it first):

1. **OpenRouter Fusion** — panel → judge ("compare, don't merge", structured
   consensus / contradictions / partial-coverage / unique-insights / **blind-spots**
   analysis) → final answer; confidence by breadth of agreement; synthesis-as-value.
2. **OpenAI ExecPlans (PLANS.md)** — a single self-contained *living* design
   document (Progress / Surprises & Discoveries / Decision Log / Outcomes &
   Retrospective + content sections) that lets a stateless agent resume a complex
   task **from the document alone**.

**Question:** Which concepts from these two systems, if any, could meaningfully
improve the **parley-deck protocol** — and which are already covered here in a
stronger form, redundant, or a net negative? parley-deck is itself a multi-agent
deliberation + artifact protocol, so be honest about overlap (cf. how the §13/RHO
work concluded we already owned stronger versions of RHO's risky single-model steps).

## Scope & boundaries (important)

- **Design / brainstorm only — for human review.** Do NOT propose or write protocol
  text edits. No implementation. Per §7, any real change needs its own ratified
  meta-protocol-change idea + human approval.
- Output **hypotheses + a clear, prioritized recommendation** (adopt / adapt / reject,
  with rationale and rough cost/benefit). Call out what we should explicitly NOT adopt.
- Respect parley-deck invariants: append-only signoffs, one-file-per-agent-per-round,
  canonical artifacts, human gate for protocol changes, multi-agent (non-solo).

## Round-01 deliverable

Each participant writes `round-01/<agent-id>.md` (independent, before reading other
round-01 files). Suggested shape: a short summary, then per-concept verdicts
(adopt/adapt/reject + why + cost), then a prioritized top-3 recommendation, then
risks / what-not-to-adopt. The neutral question list in `reference/research.md`
Part C is a menu, not a checklist — add your own angles.

## Process

- Phase 1: independent `round-01/<agent>.md`.
- Phase 2: cross-review `round-02/<agent>.md` (address every other participant; give
  counter-proposals for disagreements).
- Phase 3: `consensus.md` + append-only signoffs.
- Then: facilitator presents the consensus to the human for review. Stop there.
