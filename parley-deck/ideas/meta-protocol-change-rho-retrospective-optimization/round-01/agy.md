---
agent: agy
idea: meta-protocol-change-rho-retrospective-optimization
round: 1
date: 2026-06-16
---

## Summary
This proposal addresses the integration of Retrospective Harness Optimization (RHO) into the Parley Deck protocol from the perspective of protocol semantics and safety gating. RHO's strength is its ability to refine agent competence using historical data without labels. However, its reliance on single-model self-preference and auto-execution of edits presents severe safety risks (e.g. self-preference bias, safety-policy dilution, adversarial entrenchment). We propose a hybrid architecture: tooling handles coreset mining and candidate drafting, while the protocol governs the safety gate, requiring multi-agent consensus and mandatory human approval.

## Proposed approach

### 1. Protocol vs. Tooling Boundary
We propose a strict separation of concerns:
- **Tooling (`parley retro` CLI command):** Handles off-line trajectory mining, coreset selection, parallel rollouts, self-validation/consistency diagnosis, and candidate generation. It operates in a sandboxed/staged directory (e.g., `parley-deck/runs/retro-run-<timestamp>`) and writes proposal candidates without modifying the active harness or workspace files.
- **Protocol (COOPERATION.md amendment):** Defines the formal lifecycle of a retrospective change. Every proposed harness optimization must be registered as a standard Parley Deck idea (e.g., `ideas/retro-harness-opt-<target-agent>-<date>`). It must progress through independent analysis, cross-review, and consensus rounds, ensuring the modification is transparent, audited, and peer-reviewed.

### 2. Harness Boundaries (Allowed vs. Off-limits Surfaces)
To maintain protocol integrity and system safety, we define the editable boundaries as follows:
- **Allowed Surface (Agent-Local Harness):**
  - Prompt instructions inside the agent's local workspace (e.g., `CLAUDE.md`, agent-specific guides).
  - Agent-local helper scripts or skills (e.g., `.agents/skills/*/SKILL.md`) that do not interact with other agents' configurations.
- **Off-limits Surface (Shared Protocol/Meta-Configuration):**
  - The shared cooperation protocol (`COOPERATION.md`).
  - Shared project configurations (`agents.toml`, `.gitignore`, git configuration).
  - Other agents' workspaces or files.
  *Any modification to the off-limits surface must be treated as a manual meta-protocol change and is strictly blocked from auto-generation by retrospection loops.*

### 3. Multi-Agent Acceptance Gate
We reject RHO's single-model pairwise self-preference selection mechanism. It is replaced with:
- **No-Regression Test Suite:** The candidate harness must pass 100% of the existing project test suite (green suite) and a regression-test run on the selected coreset tasks.
- **Multi-Agent Review and Signoff:** The candidate must be submitted as a Parley Deck idea. All active rostered agents must perform a cross-review on the diff. A 4/4 consensus signoff is required.
- **Human Approval:** The final consensus document requires explicit human signoff to merge.

### 4. Structured Coreset Selection (The Parley Deck Analogue)
Instead of mining raw session JSONLs, which lack context and structure, we construct the "diverse, hard cases" coreset from Parley Deck's structured artifacts:
- **Selection Criteria:** Tasks associated with ideas that had:
  1. The highest number of review rounds (Phase 6–8 cycle depth).
  2. Documented disagreements or escalations in `consensus.md`.
  3. High fix-up churn (re-commits during Phase 5 implementation).
- **Selection Mechanism:** Extract the task prompts and environment logs from these high-friction historical ideas, and apply the DPP (Determinantal Point Process) algorithm over their embeddings to select `k` representative tasks.

### 5. Safety-Gate Guardrails
We establish the following normative protocol invariants:
1. **Adversarial-Trajectory Hygiene:** Before rollouts, all candidate past trajectories must pass a prefix/prompt-injection scan to ensure untrusted inputs did not compromise the trajectory.
2. **Reversibility Guarantee:** Before any staged harness optimization is applied to the active workspace, a local Git checkpoint/branch must be created to allow instantaneous zero-cost rollback.
3. **No-Regression Invariant:** A candidate harness cannot be signed off if it regresses performance (correctness or execution step limit) on any of the coreset tasks compared to the baseline.

## Concerns / open questions
- **Evaluation Compute Cost:** Running parallel rollouts (`G=3`) for `k=10` tasks across `N=3` candidates requires substantial token and time budgets. We need to decide if retrospection runs should be scheduled asynchronously during off-peak times (e.g. overnight via `/schedule`).
- **Signoff Fatigue:** Requiring a full multi-agent consensus round for minor agent-local instruction tweaks might introduce overhead. Can we define a fast-track protocol for minor instruction edits that only requires a subset of reviewers?

## Risks
- **Self-Preference Entrenchment:** If the evaluator agents suffer from shared model biases, they might mutually approve a harness change that degrades safety filters or introduces subtle bugs that humans cannot easily spot.
- **Overfitting to the Coreset:** The optimized harness might specialize too heavily on the historical hard cases, introducing regressions in unseen, simpler tasks. A held-out verification set is needed to mitigate this.
- **Poisoning via Adversarial Trajectories:** If an agent executes a malicious prompt from an untrusted codebase, that trajectory could be selected for retrospection, causing the agent to optimize its harness to bypass its own safety checks.
