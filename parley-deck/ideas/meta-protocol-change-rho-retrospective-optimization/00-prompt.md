---
idea: meta-protocol-change-rho-retrospective-optimization
author: claude
created: 2026-06-16
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + protocol-fit / mapping lens
  codex: tooling / Go-implementation feasibility lens
  agy: protocol-semantics / safety-gate lens
  hermes: risk / adversarial / self-preference-bias lens
status: round-01
design_only: true
cross_review_rounds: 1
---

## Problem / idea

The owner found **Retrospective Harness Optimization (RHO)** (arXiv:2606.05922,
code github.com/wbopan/retro-harness) and wants us to evaluate it and **propose
improvements to the parley-deck protocol** that adopt the good ideas — for the
owner to **review before any implementation**.

**Read `reference/rho-research.md` first** — it is the full briefing (method,
results, limitations/risks, the reference implementation, and a mapping table of
RHO concepts → parley-deck structures we already have).

RHO in one line: self-supervised improvement of an agent's *harness*
(skills/tools/instructions/workflows) from **past trajectories only** — no labels,
no external grader — via diverse-hard-case coreset → re-rollout → self-validation
+ self-consistency diagnosis → best-of-N candidate harnesses → accept by pairwise
self-preference **only if it does not regress**. The reference impl literally
mines `~/.claude/projects/<slug>/*.jsonl` and edits `CLAUDE.md`/`SKILL.md` — i.e.
our own infrastructure. parley-deck already owns the richest version of every RHO
ingredient (structured idea/review/consensus artifacts; genuine multi-agent
disagreement instead of single-model self-consistency; consensus+signoff instead
of single-model self-preference; human-approved protocol changes; a drift guard).

**This idea is DESIGN ONLY (Phases 0–4).** Produce a `FINAL.md` proposal; the
owner reviews it before we touch any code or protocol text. Do not implement.

## What to decide (deliverable = a concrete, minimal, reviewable proposal)

Address the seed questions in §7 of the research doc. Concretely, each
participant should take a position on:

1. **Protocol vs. tooling vs. both.** Should RHO enter as a new COOPERATION.md
   section (e.g. "Retrospective optimization"), as a `parley retro` command +
   retrospection workflow that mines our own idea/run history, or both — and
   exactly where the boundary is.
2. **What our "harness" is and what a retro pass may edit vs. must NOT.** Protocol
   text is consensus-gated and drift-guard-bound (both COOPERATION.md copies);
   SKILL.md is vendor-neutral; memory is local. Define the editable surface and
   what must route through a normal/meta idea + human approval.
3. **Acceptance gate.** Replace RHO's single-model pairwise self-preference with
   our multi-agent consensus + 4/4 signoff + no-regression (drift guard, green
   suite, re-review)? State the gate normatively.
4. **The coreset analogue.** What counts as our "diverse, hard cases" — over raw
   session jsonl, or over structured artifacts (ideas with most review cycles,
   dismissed findings, escalations, fix-up churn, blocked rounds)? How selected.
5. **Mandatory guardrails** given the paper's risks (self-preference bias,
   adversarial-trajectory entrenchment, irreversibility, selection noise): human
   approval for protocol edits, full audit, adversarial-trajectory hygiene, strict
   no-regression accept, reversibility/backup. Which become normative.

## Constraints

- DESIGN ONLY this idea; the owner reviews FINAL.md before implementation.
- Any protocol-text change is a meta-protocol-change (this idea is one) and must
  keep BOTH `parley-deck/COOPERATION.md` and `internal/protocol/defaults/
  COOPERATION.md` in lockstep (the drift guard enforces it) — but write NO
  protocol text yet; FINAL.md describes WHAT to change, not the diff.
- Keep it minimal and reviewable; separate "protocol amendment" from "new tooling"
  so the owner can approve each independently.
- Respect the paper's safety findings — do not propose anything that lets a single
  model silently entrench harness/protocol edits.
- English-only; transport github-pr.

## Non-goals

- Implementing the proposal (separate, post-review).
- Re-deriving RHO (the research doc is the briefing).
- Weight training / anything requiring labeled validation sets.
- Auto-editing the live protocol without human approval + multi-agent consensus.
