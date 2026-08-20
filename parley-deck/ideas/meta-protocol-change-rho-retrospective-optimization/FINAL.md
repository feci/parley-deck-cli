---
idea: meta-protocol-change-rho-retrospective-optimization
status: final
drafted-by: claude
date: 2026-06-16
participants: [claude, codex, agy, hermes]
consensus: 4/4 ACCEPT
design_only: true
awaiting: owner-review-before-implementation
---

## What this is

A **design proposal** (not an implementation) for adopting **Retrospective
Harness Optimization (RHO**, arXiv:2606.05922) into parley-deck. Produced by a
full Parley Deck deliberation (claude + codex + agy + hermes; round-01 →
round-02 → consensus, 4/4). Per the owner's instruction, this FINAL.md is for
**review before any code or protocol text is written**. Full research briefing:
`reference/rho-research.md`.

## The core idea (one paragraph)

RHO improves an agent's *harness* (skills/tools/instructions/workflows) from
**past trajectories only**, with no labels and no external grader, and accepts a
change only if it doesn't regress. parley-deck already owns *stronger* versions of
RHO's three single-model steps: genuine multi-agent disagreement instead of
self-consistency, and consensus + all-participant signoff + human-approved
protocol changes instead of single-model self-preference. **So we adopt RHO's
retrospective loop and its strict no-regression discipline, drop its single-model
judgments in favor of our quorum, and never auto-apply anything.**

## Proposal in brief

**A. Protocol amendment** — a new, small top-level section "Retrospective
optimization" in both `COOPERATION.md` copies. It is normative and says: a retro
pass mines prior Parley Deck history to *propose* improvements; its output is
**advisory input only and applies nothing**; proposals enter as a normal idea
(Phases 0–8), and any protocol-text change goes through a meta-protocol-change
idea + human approval + drift-guard lockstep. It defines the editable-surface
layers, the acceptance gate, the audit/provenance/hygiene rules, and the
prohibited direct edits.

**B. Tooling** — an optional, **read-only-by-default** `parley retro` CLI with
staged subcommands `scan` / `select` / `diagnose` / `propose`. It mines our own
structured artifacts, builds a deterministic "diverse, hard cases" coreset, emits
a diagnosis, and at most **scaffolds a single new `ideas/<slug>/00-prompt.md`** —
the seed of a normal idea. It never writes any participant/consensus/review/
final/protocol/skill/memory file.

## The decisions (consensus D1–D7)

1. **New top-level protocol section** "Retrospective optimization" — advisory
   input only; never auto-applies.
2. **Layered harness** (what a retro pass may touch):
   - *Protocol* (both `COOPERATION.md`) → meta-protocol-change idea only, human
     approval, drift-guard lockstep. Never a direct retro edit.
   - *Runtime/shared — "Repository Instruction Files"* (repo `CLAUDE.md`,
     `AGENTS.md`, checked-in `SKILL.md`, CLI behavior, helper scripts, docs) →
     ordinary idea + full review gate.
   - *Local — "Agent Local Memory"* (`~/.claude`, `~/.codex`, caches, ignored
     config) → retro **reports only**, never canonical.
   - *Evidence corpus*: structured artifacts primary; raw session JSONL
     secondary, quarantined/off-by-default, provenance-recorded.
3. **Acceptance gate** = multi-agent consensus + all-participant signoff + human
   approval (protocol/shared harness) + strict no-regression (drift guard green,
   suite green, clean re-review). Self-preference is an advisory note, never the
   gate.
4. **Coreset v1** = deterministic ranking over structured artifacts by
   failure-density signals (review cycles, NOT-FIXED/recurring/dismissed-recurring
   findings, escalations, fix-up churn, blocked/abandoned, drift-guard failures,
   watchdog/`agent.failed`); diversity by failure-mode type. DPP/embeddings later.
5. **Tooling v1** = `parley retro scan/select/diagnose/propose`, read-only
   default; `propose` writes only a new `00-prompt.md` at an explicit non-existing
   slug, fail-if-exists, nothing else.
6. **Normative guardrails** = full audit (retro pass is an idea); provenance for
   selected + excluded sources; adversarial-trajectory hygiene + exclusion of
   compromised/out-of-scope evidence; reversibility via idea branch + git;
   multi-agent diagnosis (each agent diagnoses the coreset independently in the
   retro idea's round-01).
7. **Cut from v1**: fast-track for "minor" edits, dedicated held-out corpus, live
   re-rollout, best-of-N candidate generation, DPP/embeddings, auto-application of
   any diff, persistent quarantine registry.

## How the design neutralizes RHO's stated risks

| RHO risk (from the paper) | Our mitigation in this design |
|---|---|
| Self-preference amplifies one model's mistaken judgment | Acceptance is multi-agent consensus + 4/4 signoff + human approval; self-preference is advisory only (D3) |
| Adversarial trajectories entrench bad behavior | Structured artifacts primary; raw JSONL quarantined; compromised/out-of-scope sources excluded + provenance-recorded (D2, D6) |
| Persistent edits from model judgments | Nothing auto-applies; every change is a normal/meta idea with human approval for protocol/harness (D1, D3) |
| Selection noise / coreset over-fit | Strict no-regression gate (drift guard + suite + re-review) is the held-out check; diversity by failure-mode type (D3, D4) |
| Needs resettable environments | No live re-rollout in v1 — diagnose from existing artifacts (D7) |

## If approved — the implementation path (two separately-approvable ideas)

1. **Protocol amendment idea** (meta-protocol-change): write the new section into
   both `COOPERATION.md` copies in lockstep (drift guard enforces it), run the
   normal review/signoff, ship.
2. **Tooling idea** (`parley retro`): implement the read-only staged CLI +
   deterministic coreset + tests; ship through the normal lifecycle.

These are independent so you can approve one, both, or neither.

## Next step

**Owner review.** No implementation has started and none will until you approve.
Tell us: (a) proceed with both follow-up ideas, (b) protocol-only, (c)
tooling-only, (d) changes to the design first, or (e) hold.
