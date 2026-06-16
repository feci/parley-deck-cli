---
idea: meta-protocol-change-rho-retrospective-optimization
drafted-by: claude
date: 2026-06-16
round: consensus
participants: [claude, codex, agy, hermes]
design_only: true
---

## Decisions (all 4/4 — design only; nothing is implemented by this idea)

All four participants converged across round-02 with no open disagreement. The
proposal **adopts RHO's retrospective loop and its strict no-regression
discipline, but replaces RHO's single-model self-validation/self-consistency/
self-preference with parley-deck's existing multi-agent quorum, and never
auto-applies anything.**

**D1 — New top-level protocol section "Retrospective optimization."** Small and
normative; NOT a subsection of §7. It states: a retro pass mines prior Parley
Deck history to *propose* improvements; its output is **advisory input only and
applies nothing**. Proposals enter as a normal idea (Phases 0–8); protocol-text
changes require a meta-protocol-change idea + human approval + drift-guard
lockstep across both `COOPERATION.md` copies. The section defines eligibility,
evidence, acceptance, audit, and the prohibited direct edits — not tool schemas.

**D2 — Layered harness definition** (the editable-surface boundary):
- **Protocol harness** — both `COOPERATION.md` copies. Meta-protocol-change idea
  only; human approval; drift-guard lockstep. No direct retro edits, ever.
- **Runtime / shared harness — "Repository Instruction Files"**: tracked, shared
  files (repo `CLAUDE.md`, `AGENTS.md`, checked-in `SKILL.md`, CLI behavior,
  helper scripts, docs). Changes route through an ordinary idea + full review
  gate (a meta idea if they alter protocol semantics).
- **Local harness — "Agent Local Memory"**: operator-local state (`~/.claude`,
  `~/.codex`, caches, ignored headless config). Retro may **report observations
  only**; it must never make them canonical or infer protocol rules from one
  operator's local setup.
- **Evidence corpus**: structured Parley artifacts (`ideas/*` rounds, `review/`,
  `consensus.md`, `FINAL.md`, `IMPLEMENTATION.md`, run event logs) are primary;
  raw session JSONL is **secondary, quarantined/off-by-default, provenance-
  recorded** when explicitly included.

**D3 — Acceptance gate.** A retro-proposed change is accepted only by: multi-agent
consensus + **all-participant signoff** + **human approval** (for protocol /
shared-harness changes) + **strict no-regression** (drift guard green, relevant
test suite green, clean multi-agent re-review). RHO-style self-preference may be
**attached as a diagnostic note only** — it is never the gate. (Directly answers
the paper's #4 risk: single-model self-preference amplifying mistaken judgment.)

**D4 — Coreset v1 (the "diverse, hard cases" analogue).** Deterministic ranking
over structured artifacts using failure-density signals already recorded: number
of review cycles, NOT-FIXED / recurring findings, dismissed-but-recurring
findings, escalations (`inbox/*to-user*`), fix-up churn, blocked/abandoned work,
drift-guard failures, watchdog / `agent.failed` events. Diversity by failure-mode
*type* (protocol ambiguity, transport mechanics, role/quorum confusion, artifact
ownership collisions, review-gate weakness, tooling friction, safety/adversarial).
DPP/embedding-based diversity is a later enhancement, not v1.

**D5 — Tooling spec `parley retro` (v1).** Optional CLI, **read-only by default**,
staged subcommands `scan` / `select` / `diagnose` / `propose`. `propose` may write
**only a single new `ideas/<slug>/00-prompt.md`** at an explicit, non-existing
slug and must **fail if the target exists**; it must never write rounds,
consensus, review, final, implementation, protocol, skill, memory, or
helper-script files. (Exact schemas/flags live in the tooling/implementation
idea, not the protocol text.)

**D6 — Normative guardrails.** Full audit trail (a retro pass IS an idea);
provenance for both selected and excluded sources; adversarial-trajectory
hygiene; exclusion of compromised or out-of-scope evidence (recorded in the run's
coreset/provenance manifest); reversibility via an idea branch + git history;
**multi-agent diagnosis** — the retro idea's round-01 has each agent diagnose the
coreset independently (our genuine, less-biased analogue of RHO's
self-consistency), applied at diagnosis, not only at acceptance.

**D7 — Explicitly cut from v1** (to keep it minimal, safe, reviewable): fast-track
semantics for "minor" edits (the layer decides the path instead); a dedicated
held-out corpus (the no-regression gate is the held-out check); live re-rollout
(RHO's G=3 — assumes resettable envs, costly); best-of-N candidate-harness
generation; DPP/embedding selection; auto-application of any diff; and a
persistent quarantine registry (v1 records exclusions per-run; a registry is a
follow-up only if repeated runs prove it needed).

## Implementation split (for the owner's review — NOT executed by this idea)

Two independent, separately-approvable follow-up ideas, so the owner can accept
protocol and tooling on their own merits:
1. **Protocol amendment** (meta-protocol-change): add the §"Retrospective
   optimization" section to BOTH `COOPERATION.md` copies (drift-guard lockstep).
2. **Tooling** (`parley retro`, ordinary idea): the read-only `scan/select/
   diagnose/propose` CLI per D5, with deterministic coreset (D4) and tests.

## Deferred follow-ups

DPP/embedding coreset; held-out corpus; live re-rollout + best-of-N (sandboxed
worktrees, no network, resettable fixtures); persistent quarantine registry;
raw-JSONL ingestion beyond quarantined secondary evidence.

## Signoffs

<!-- Each participant APPENDS their own signoff block. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-16
Status: ✅ ACCEPT
Notes: I drafted this; it faithfully captures the 4/4 round-02 convergence (RHO loop adopted, single-model judgments replaced by our quorum, never auto-applies, minimal v1). Design only — for owner review before any implementation.

### Signoff: codex — 2026-06-16
Status: ✅ ACCEPT
Notes: This matches my converged proposal: advisory retro mining, layered harness gates, structured evidence first, and no v1 auto-application.

### Signoff: hermes — 2026-06-16
Status: ✅ ACCEPT
Notes: Consensus.md accurately reflects the 4/4 convergence: RHO loop adopted, single-model steps replaced by multi-agent quorum at both diagnosis and acceptance, strict no-regression, never auto-apply, protocol text in new top-level section, tooling limited to 00-prompt.md only.

### Signoff: agy — 2026-06-16
Status: ✅ ACCEPT
Notes: I accept; this draft faithfully captures the layered harness definitions, independent multi-agent diagnosis, and deterministic failure-density coreset for a minimal and safe v1.
