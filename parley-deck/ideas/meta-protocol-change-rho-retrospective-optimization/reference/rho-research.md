# Research briefing: Retrospective Harness Optimization (RHO)

Prepared by claude (facilitator) for the Parley Deck deliberation. Source of
truth is the paper + repo; this is a faithful synthesis, not the protocol.

## 1. What RHO is (one paragraph)

**Retrospective Harness Optimization (RHO)** is a self-supervised method that
improves an LLM agent's **harness** — its skills, tools, instructions, and
workflows — using **only past trajectories**, with **no ground-truth labels and
no external grader**. In one retrospective pass it mines hard cases from history,
re-solves them, diagnoses what went wrong using the agent's own introspection,
proposes a few edited harnesses, and keeps the best one by the agent's own
pairwise preference. Headline result: one round lifts SWE-Bench Pro pass rate
**59% → 78%** without any external grading.

- Paper: arXiv:2606.05922 — v1 "Retrospective Harness Optimization: Improving LLM
  Agents via Self-Preference over Trajectory Rollouts"; v2 "Evolving Agents in
  the Dark: Retrospective Harness Optimization via Self-Preference".
- Code: github.com/wbopan/retro-harness (MIT).

## 2. The method, end to end

Three sequential stages over a dataset of past trajectories:

**Stage 1 — Coreset selection (hard + diverse).**
- An LLM judge scores each past trajectory's difficulty `r ∈ [0,10]` and writes a
  textual "fingerprint" of its structural challenges.
- A **Determinantal Point Process (DPP)** selects `k=10` tasks balancing
  difficulty and diversity: kernel `K = diag(r̃)·S·diag(r̃)`, where `S` is the
  cosine-similarity matrix of fingerprint embeddings and `θ=0.7` trades off
  difficulty vs. coverage. DPP(θ=0.7) beat pure-difficulty, pure-diversity, and
  random sampling (Fig 5).

**Stage 2 — Group rollout + diagnosis (the label-free signal).**
- Re-solve each coreset task `G=3` times in parallel under the baseline harness
  `h₀`, fixing the first trajectory as the reference.
- Two introspective signals, no ground truth:
  - **Self-validation**: the agent inspects a single trajectory against task
    requirements + environment feedback, flagging wrong tool calls, false
    assumptions, premature stops.
  - **Self-consistency**: the agent compares the parallel trajectories; divergent
    plans/answers flag high uncertainty.
- Merged into a structured instruction set `I` with per-task improvement
  directions and severity weights.
- **Ablation (Table 4)**: full RHO 0.78; without self-consistency 0.56; without
  self-validation 0.70; raw trajectories without diagnosis 0.60. **Both signals
  are essential** — explicit diagnosis >> dumping raw trajectories.

**Stage 3 — Best-of-N harness proposal + self-preference selection.**
- Generate `N=3` candidate harnesses via an `optimize(h₀, I)` operator — a code
  agent with write access to a fresh copy of the harness (a directory of files).
- It can **add/modify/remove executable tools, skills, and instruction files**
  (filesystem-level), e.g. add a `check_build_and_lint` tool, or a skill noting
  "Go toolchain resides at a non-standard location."
- Re-solve the coreset with each candidate; score by **pairwise self-preference**:
  `Sⱼ = mean over coreset of rank(task, τ_candidate, τ_baseline)`, rank ∈ [-10,10],
  judged by the agent comparing trajectories on correctness/efficiency/approach.
- **Accept the top candidate only if `Sⱼ > 0`; otherwise keep `h₀`** (strict
  no-regression gate).

**Single-pass by design**: one retrospective pass over unlabeled history — no
iterative validation-scored loop, no weight training. Vs. the "Meta-Harness"
baseline it reaches 0.78 at 1.0× compute where Meta-Harness needs 3.1× compute
*and* labels.

## 3. Results & behavioral findings

| Domain | Vanilla | RHO | Δ |
|---|---|---|---|
| SWE-Bench Pro (software eng) | 0.59 | 0.78 | +0.19 |
| Terminal-Bench 2 (CLI) | 0.71 | 0.76 | +0.05 |
| GAIA-2 (knowledge work) | 0.29 | 0.37 | +0.08 |

- Gains concentrate on **long-horizon tasks** (>50 steps); the optimized harness
  shifts the agent toward **verification** actions.
- Best-of-N robustness (Table 3): selection reliably avoids the worst candidate
  but does **not** perfectly pick the test-best one (self-preference is noisy).

## 4. Limitations & risks (verbatim emphasis from the paper)

1. **Needs resettable environments** — group rollout replays tasks several times;
   "one-shot or irreversible tasks [are] outside the setting RHO targets."
2. **Harness-mediated competence assumption** — only helps where an editable
   harness of prompts/skills/tools mediates a meaningful part of competence.
3. **Trusts past trajectories as the only input** — in open environments these
   "can embed adversarial content injected mid-task," and edits distilled from
   compromised trajectories "could entrench such behavior."
4. **Self-preference can amplify mistakes** — "RHO modifies persistent agent
   behavior from model-generated judgments. This can amplify mistaken
   preferences, unsafe procedures, or biased behavioral rules if the evaluator
   prefers them." (self-preference bias / reward-hacking risk.)
5. **Selection noise** — the chosen harness beats the worst but isn't always the
   test-best.

**Paper's own mitigations**: keep **full audit logs**, **require human approval
for sensitive harness edits**, use **domain-specific safety checks** before
applying an accepted harness to high-impact tasks.

## 5. Reference implementation (github.com/wbopan/retro-harness)

- Harness = a collection of markdown + scripts, per target:
  **Claude Code** = `CLAUDE.md` + auto-memory dir + helper scripts;
  **Codex CLI** = `AGENTS.md` + `.agents/skills/*/SKILL.md`.
- Mines trajectories from **`~/.claude/projects/<slug>/*.jsonl`** (incl. worktree
  sessions) and **`~/.codex/sessions/YYYY/MM/DD/rollout-*.jsonl`** (filtered by
  project dir).
- Loop `src/rho/loop.py` = select → rollout → propose; pluggable `selection/`,
  `strategies/`, `orchestrators/`, `stores/`.
- Ships a **`/retrospection` dynamic workflow** (`.claude/workflows/retrospection.js`)
  and a stdlib `codex/retrospection.py`.
- CLI: `rho evolve` (one round), `rho solve`, `rho ui` (run browser).
- Knobs: `k` (coreset, 10), `n` (candidates, 3), `probes` (self-preference probes, 4),
  `theta` (0.7), `maxSessions`, `apply` (modify live files vs. stage only).
- **Safety in the impl**: "winner applied only if mean score is positive, with a
  full backup first"; backups under `~/.claude/rho-runs/<ts>-<project>/backup/`;
  rich run artifacts (prompts, diagnoses, candidates, diffs, held-out reports).

## 6. Why this is a near-perfect fit for parley-deck (mapping)

| RHO concept | parley-deck equivalent we already have |
|---|---|
| **Harness** (editable skills/tools/instructions) | `SKILL.md`, **`COOPERATION.md` protocol**, the embedded default, `CLAUDE.md`, the auto-memory, helper scripts, the `parley` CLI |
| **Past trajectories** | `~/.claude/projects/-…-parley-deck/*.jsonl` **plus** our durable, *structured* artifacts: `ideas/*/round-*`, `consensus.md`, `review/`, `FINAL.md`, `IMPLEMENTATION.md`, `runs/*` event logs, signoffs |
| **Self-validation** | the per-agent review files (CRITICAL/MAJOR/MINOR/NIT) + RunChecks |
| **Self-consistency** | **multi-agent disagreement is native to us** — independent round-01, cross-review, consensus signoffs already surface divergence (arguably stronger than RHO's G=3 same-agent reruns) |
| **Self-preference / no-regression gate** | our **consensus + 4/4 signoff** gate, the drift guard test, "accept only if it doesn't regress" discipline |
| **Human approval for sensitive edits** | already mandatory for protocol changes (meta-protocol-change ideas, §7) |
| **Single-pass, audited** | our ideas ARE the audit trail; a retro pass would itself be a Parley Deck idea |

Notable: RHO's self-consistency relies on re-running ONE agent G times; parley-deck
gets consistency signal "for free" from genuinely independent multi-agent rounds.
And RHO's biggest risk — self-preference bias entrenching one model's mistaken
judgment — is exactly what parley-deck's multi-agent quorum + human-approved
protocol changes are designed to counter.

## 7. Seed questions for the deliberation (not conclusions)

These are prompts for the agents, deliberately open:

1. **Should RHO enter parley-deck as protocol, as tooling, or both?** e.g. a new
   `## N. Retrospective optimization` protocol section, and/or a `parley retro`
   command + a retrospection workflow that mines our own idea/run history.
2. **What is OUR "harness" that a retro pass may edit, and what is off-limits?**
   (Protocol text is consensus-gated; the drift guard binds both COOPERATION.md
   copies. SKILL.md is vendor-neutral. Memory is local. Where can a retro pass
   propose edits, and what must route through a normal/meta idea?)
3. **Replace RHO's single-model self-preference with our multi-agent quorum?**
   Does our consensus/signoff already give a stronger, less-biased acceptance
   gate than pairwise self-preference — and should the protocol say so?
4. **What is the no-regression gate for a protocol/harness edit?** (drift guard,
   green suite, re-review; "accept only if it doesn't regress" formalized.)
5. **Coreset over what?** Our richest signal is structured idea/review artifacts,
   not just raw jsonl. What's the parley-deck analogue of "diverse coreset of hard
   cases" — e.g. ideas with the most review cycles, dismissed findings,
   escalations, fix-up churn, blocked rounds?
6. **Mandatory guardrails** given the paper's risks: human approval for protocol
   edits (already true), full audit (already true), adversarial-trajectory
   hygiene, and the strict no-regression accept gate. What must be normative?
7. **Scope discipline** — keep the proposal minimal and reviewable; separate
   "protocol amendment" from "new tooling" so the user can approve each.

## 8. Sources

- Paper (v2, full): https://arxiv.org/html/2606.05922v2
- Paper (abstract): https://arxiv.org/abs/2606.05922 · https://huggingface.co/papers/2606.05922
- Code: https://github.com/wbopan/retro-harness
- Discussion: https://www.alphaxiv.org/abs/2606.05922 (403 at fetch time; listed for the record)
