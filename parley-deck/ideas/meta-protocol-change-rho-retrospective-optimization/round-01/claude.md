---
agent: claude
idea: meta-protocol-change-rho-retrospective-optimization
round: 1
date: 2026-06-16
---

## Summary

RHO is a strong fit, but parley-deck already owns *stronger* versions of RHO's two
riskiest single-model steps (self-consistency, self-preference) via genuine
multi-agent disagreement and consensus+signoff. So the proposal should **adopt
RHO's retrospective LOOP and its no-regression discipline, but replace its
single-model judgments with our quorum, and never auto-apply**. The cleanest shape
is: tooling that *mines history and drafts a retro idea*, + a thin protocol section
that says *retro findings are just another idea and are accepted only through the
normal gate*.

## Proposed approach

**Q1 — both, with a hard boundary.** Tooling produces the *input*; protocol
governs *acceptance*.
- Tooling: `parley retro` + a `retrospection` workflow that mines our own history
  (structured idea/review artifacts first, session `*.jsonl` second), builds a
  diverse-hard coreset, runs the diagnosis, and **drafts a normal idea's
  `00-prompt.md`** (a "retro idea") with the candidate improvement directions.
- Protocol: a short new section "Retrospective optimization" stating that a retro
  pass is non-authoritative input; its proposals enter ONLY as a normal idea
  (Phases 0–8) and, for protocol text, a meta-protocol-change idea. No auto-apply.

**Q2 — editable surface.** A retro pass may *propose* edits to: SKILL guidance,
`CLAUDE.md`, helper scripts, tooling, docs, and the auto-memory (local). It may
*never* directly edit: `COOPERATION.md` (both copies, drift-guard-bound) or
anything touching quorum/ownership/signoff semantics — those require a
meta-protocol-change idea + human approval. The retro tool only *writes a
00-prompt*, never a protocol artifact owned by another agent.

**Q3 — acceptance gate (normative).** Replace RHO's single-model pairwise
self-preference with our existing, stronger gate: a retro-proposed change is
accepted only via idea consensus + **all-participant signoff** + **no-regression**
(drift guard green, full `go test ./...` green, clean re-review). State plainly:
"No retrospective edit is applied without multi-agent consensus and, for protocol
text, human approval." This directly answers the paper's #4 risk
(self-preference amplifying one model's mistaken judgment).

**Q4 — the coreset analogue.** Select over **structured artifacts**, which are
richer than raw jsonl: rank by failure-density signals we already record —
number of review cycles, NOT-FIXED re-reviews, dismissed findings, escalations
(`inbox/*to-user*`), fix-up churn, blocked rounds, watchdog/`agent.failed` events,
`driver.error`. Diversity over idea "fingerprints" (topic + files/areas touched),
DPP-style. Raw `*.jsonl` is a secondary signal for steps that never produced a
clean artifact.

**Q5 — mandatory guardrails (normative).** (a) human approval for any protocol
edit (already true via §7 — restate); (b) full audit — the retro pass *is* an
idea, so it is audited by construction; (c) adversarial-trajectory hygiene — never
distill from runs marked compromised or from adversarial steers; (d) strict
no-regression accept; (e) reversibility — all edits land on an idea branch with
git history + the drift guard, never a silent in-place rewrite.

**Bonus — diagnosis should also be multi-agent.** RHO's self-validation/
self-consistency run on one model. To avoid re-introducing single-model bias at
the *diagnosis* stage, the retro idea's round-01 should have each agent
independently diagnose the coreset; consensus over diagnoses is our analogue of
"self-consistency", and it is genuinely independent.

## Concerns / open questions

- Scope creep: the mining/coreset tooling is real engineering. I'd phase it —
  protocol section + a minimal `parley retro` that drafts a 00-prompt from simple
  artifact-derived signals first; DPP/embeddings later. Keep THIS proposal to the
  design; let implementation be its own reviewed idea.
- Determinism/cost: re-rollout (RHO's G=3) is expensive and assumes resettable
  envs. Do we even need re-rollout, or is our existing artifact history enough
  signal? I lean: skip live re-rollout v1; diagnose from existing artifacts.
- Where does the protocol section live — a new top-level `## N` or a sub-section
  of §7 (Changing this protocol)? It interacts with §7 and §8 (inbox).

## Risks

- **Self-referential over-fitting**: mining our own transcripts to change our own
  protocol can entrench house style. Mitigant: multi-agent diagnosis + human
  approval; treat retro output as hypotheses, not findings.
- **Adversarial/steer poisoning**: steers and external content in trajectories
  could bias proposals; hygiene rule needed (exclude flagged runs).
- **Tooling becoming a second source of truth**: the retro tool must only emit a
  `00-prompt` (input), never canonical artifacts — otherwise it bypasses the
  protocol it is supposed to feed.
