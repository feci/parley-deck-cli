---
agent: hermes
idea: meta-protocol-change-rho-retrospective-optimization
round: 1
date: 2026-06-16
---

## Summary
RHO's core loop (hard-case mining + diagnosis + candidate generation + strict no-regression) is valuable, but its single-model self-preference and self-consistency steps are exactly the vectors for self-preference bias and adversarial-trajectory entrenchment the paper itself flags. Parley-deck already counters these via multi-agent disagreement and consensus+signoff. Proposal: tooling that surfaces candidate improvements as ordinary ideas only; protocol that makes any acceptance route through the existing quorum gate. No auto-apply, no single-model override of protocol or harness.

## Proposed approach
**Q1 — both, boundary explicit.** Tooling mines history and emits a draft `00-prompt.md` (a "retro idea"). Protocol adds one short section stating retro output is non-authoritative input that must still pass Phases 0-8 and, for protocol text, the meta-protocol-change path + drift guard. Tooling never writes protocol artifacts or SKILL.md directly.

**Q2 — harness surface.** A retro pass may propose edits to: local SKILL.md guidance, helper scripts, docs, auto-memory, and tooling. It must never touch: either copy of COOPERATION.md, quorum/signoff/ownership semantics, or the drift guard itself. Any such change is a meta-protocol-change idea requiring human approval.

**Q3 — acceptance gate (normative).** Replace RHO's pairwise self-preference entirely with our existing gate: multi-agent consensus + 4/4 signoff + explicit no-regression (drift guard green + full test suite green + re-review). State: "No retrospective harness or protocol edit is applied without multi-agent consensus and human approval for protocol text." This directly mitigates paper risk #4 (self-preference amplifying one model's mistaken judgment).

**Q4 — coreset analogue.** Select over structured artifacts ranked by failure-density signals already recorded: number of review cycles, NOT-FIXED re-reviews, dismissed findings, escalations, fix-up churn, blocked rounds, watchdog failures. Apply DPP-style diversity on idea fingerprints (topic + files touched). Raw session jsonl is secondary; never mine runs flagged as compromised.

**Q5 — mandatory guardrails (normative).** (a) human approval for any protocol edit (already §7); (b) full audit trail — the retro pass itself is an idea; (c) adversarial-trajectory hygiene — exclude any run marked compromised or containing external steers; (d) strict no-regression accept only; (e) reversibility — all changes land on an idea branch with git history. Diagnosis stage itself must be multi-agent to avoid re-introducing single-model bias.

## Concerns / open questions
- Cost of re-rollout: RHO's G=3 assumption does not hold for one-shot or irreversible tasks common in our workflow; v1 should diagnose from existing artifacts only.
- Over-fitting risk: mining our own artifacts to edit our own protocol can silently reinforce house style. The multi-agent + human gate is the only real counter.
- Where the protocol text lives: must be a new top-level section, not buried inside §7, so the drift guard and quorum rules remain clearly separate.

## Risks
- Self-preference bias amplification: even with multi-agent diagnosis, a dominant model voice across rounds could still entrench suboptimal preferences; the 4/4 signoff and human approval are non-negotiable.
- Adversarial trajectory poisoning: any trajectory containing injected steers or external content could bias the coreset and downstream proposals; hygiene filter must be mandatory and auditable.
- Selection noise + false acceptance: self-preference (even multi-agent) is noisy; without the strict no-regression gate we risk accepting a harness that looks better on the coreset but regresses on held-out cases. The drift guard + re-review is the minimum safeguard.