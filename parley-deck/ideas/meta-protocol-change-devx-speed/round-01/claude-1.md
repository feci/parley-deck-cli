---
agent: claude-1
idea: meta-protocol-change-devx-speed
round: 1
date: 2026-07-03
---

## Summary

The protocol is safe and complete but applies **one fixed 9-phase lifecycle to every
change regardless of risk or size**. A typo fix and a security-critical feature pay the
same tax: full quorum, sequential rounds, two full signoff barriers (Phase 3 + Phase 7),
a §9.0 liveness ping, 20–30 min deep-reasoning timeouts, and a fix-up loop that reruns a
*full* review every cycle. That is the root cause of "it takes long." The fix is
**conditional rigor**: keep the full track for risky work, add lighter tracks for the
common case, and restructure the 1046-line monolith so a developer can start in minutes.
None of this requires deleting the safety core — only stopping it from firing on work
that doesn't need it.

## 1. Time-sink audit (biggest wall-clock/effort per unit of safety bought)

1. **No fast path — uniform 9 phases (§4).** Phase 0→8 is mandatory shape for *any* idea.
   Trivial, reversible changes (docs, comments, a one-file rename) still imply the whole
   machinery. Highest-frequency work pays the highest ceremony.
2. **Two full-quorum signoff barriers per idea (§4 Phase 3 *and* Phase 7; §5).** Every
   participant must append a signoff at design consensus *and again* at review consensus.
   With 4 agents that is ≥8 human/agent-gated signoff events even with zero disagreement.
3. **§9.0 readiness ping before every idea.** A liveness round-trip to every rostered
   agent runs *before any work starts* — latency + token cost paid up front, per idea,
   even for a 2-line change.
4. **Deep-reasoning timeouts (agents.toml: round/review/deepReasoning = 20–30 min).** The
   slowest agent dominates each round; multiple rounds compound. There is no per-task
   "fast profile" — every round budgets for the worst case.
5. **Fix-up loop reruns full review each cycle (§4 Phase 8; strict_gate).** "Repeat until
   zero Agreed fixes," and `strict_gate` demands a *fresh full-scope* review round at the
   end. Loop-engineering Tier 4 took **4 cycles = 4 full multi-agent review rounds**.
6. **Monolithic doc (1046 lines).** The actual "how do I start" is buried; the TL;DR is at
   line 684; transport mechanics (§11), pipeline (§12), retro (§13), loop (§14) are all in
   the same file a first-timer must wade through. Cognitive load *is* wall-clock for a human.

## 2. Tiering model — three tracks, objective classifier

Route every task through a classifier, not a fixed lifecycle:

- **Track 0 — Direct.** Trivial, fully reversible, zero protocol/security/data surface
  (typo, docs, comment, single-file rename, dep bump with green tests). **No idea dir, no
  rounds, no quorum.** Do it, record a one-line note in a lightweight log. Solo allowed.
- **Track 1 — Lite (default for ordinary code/design).** Small, low-risk, reversible.
  Collapse Phase 0–4 into **one `SPEC.md`** (prompt + decision + acceptance criteria) and
  Phase 6–7 into **one review pass by a single *model-diverse* reviewer** (refutation-default
  kept). No cross-review rounds, no strict_gate, fix-up capped at **1 cycle then escalate/upgrade**.
- **Track 2 — Full (today's protocol, unchanged).** High-risk / irreversible / security /
  data-loss / protocol-change / `auto_implement` / public-API / migration / large surface.
  Full quorum, full rounds, refutation, strict_gate available.

**Objective triggers (any one trips upward; classifier picks the floor):**

| Signal | Track 0 | Track 1 | Track 2 (forced) |
|---|---|---|---|
| Reversibility | fully reversible | reversible | irreversible / migration |
| Files / LOC | 1 file, tiny | ≲5 files, bounded | large / cross-cutting |
| Security or data surface | none | none | any |
| Protocol change (§7) | no | no | **yes → always Track 2** |
| `auto_implement` / action / pipeline | no | no | **yes → always Track 2** |
| Public API / schema / migration | no | no | yes |

Rule: the author may **upgrade** a track anytime; **downgrading below the classifier's
floor requires a recorded user OK.** This is the same spirit as the existing `strict_gate`
and fusion "conditional rigor," generalized to the whole lifecycle.

## 3. DevX improvements

- **One-page QUICKSTART at the very top** (or a separate QUICKSTART.md the skill loads
  first): the 3 tracks, the classifier table, "start an idea in 5 minutes." Everything else
  becomes layered reference.
- **Progressive disclosure / layered doc.** Split the monolith: Layer 1 = quickstart +
  tracks (1 screen); Layer 2 = the phase reference; Layer 3 = transport (§11), pipeline
  (§12), retro (§13), loop (§14) as separate referenced files. Agents load only the layer
  their track needs — faster *and* cheaper on tokens.
- **Role framing.** A short "If you're a developer, you live in Track 0/1" section so a
  non-protocol-author knows the 5% of the doc that applies to them.
- **Defaults that just work.** local-dir default transport, auto-classified track, model +
  roster from central config — a developer types the task and the system scaffolds only
  what the track needs.

## 4. Speed levers (concrete mechanics)

- Make **parallel round-01 dispatch the documented default**; state the "wall-clock =
  slowest agent" model explicitly (today it reads sequential).
- **Collapse Phase 3+4** when there is no ❌ and a single round: the drafter writes FINAL
  with inline signoffs in one transaction (skip the separate consensus barrier).
- **Single-reviewer fast path** for Track 1 (drop full quorum on review; keep one
  model-diverse refutation reviewer).
- **Timeout profiles per track:** Track 1 uses the `fast` profile; only Track 2 uses the
  20–30 min `deep` budgets.
- **§9.0 ping:** skip entirely for Track 0/1; for Track 2 use a **cached liveness with a
  short TTL** instead of a fresh full-roster probe every idea.
- **Early-exit consensus:** close as soon as all-✅ arrive; don't wait for a deadline.
- Lean on the existing **auto-drive driver** to remove human round-trips between sequential
  phases (already shipped; make it the Track 1/2 default).

## 5. Modern agentic concepts → concrete edits

- **Conditional rigor / right-altitude** → the 3-track model itself (the core proposal).
- **Deterministic workflows vs model-driven control flow** (Anthropic, *Building effective
  agents*) → Track 0/1 are deterministic short pipelines; Track 2 stays model-driven
  deliberation. Choose control flow by risk, not by habit.
- **Lead-agent + subagent orchestration** (Anthropic multi-agent research) → sharpen §1.6:
  the facilitator may spawn *helper subagents* to draft/refute in parallel (not quorum
  members) to compress rounds; formalize a "lead drafts, subagents refute" fast review for
  Track 1.
- **Spec-driven development** (spec-kit / Kiro) → FINAL.md *is* the spec; for Track 1
  collapse prompt+consensus+final into one `SPEC.md` with observable acceptance criteria as
  the single downstream-of-code source of truth.
- **Plan mode** → make the "nothing is built until FINAL/SPEC" boundary explicit; Track 1 =
  lightweight plan-then-execute.
- **Context engineering / progressive disclosure** → the layered-doc restructure (§3).
- **Parallel worktrees** → reference the existing `parley-worktrees` addon as the default
  for Track 2 parallel implementers, so implementation doesn't serialize.
- **The bitter lesson / less scaffolding** → the justification for Track 0: trust the model
  for low-risk reversible work; remove ceremony that buys no safety.

## 6. What MUST stay (defend against speed-at-all-costs)

- Non-solo **independent** verification for anything risky (Track 2) and **refutation-default**
  review. The **durable audit trail**. **§14 human brake** + `auto_implement` gating stays
  Track 2 only. English-only + no-secrets. Files-canonical. One-file-per-agent, no cross-edit.
  The **always-available upgrade** to a higher track.

## 7. Prioritized shortlist

**MUST**
- M1 — 3-track conditional-rigor model + objective classifier (the central G1+G2 win).
- M2 — Track 0 "Direct": trivial reversible work needs no idea machinery.
- M3 — Restructure doc: 1-page QUICKSTART + layered reference files.
- M4 — Track 1: collapse Phase 0–4 → one SPEC.md; single model-diverse reviewer.

**SHOULD**
- S1 — Skip/cheapen §9.0 ping for Track 0/1; cached-TTL liveness for Track 2.
- S2 — Per-track timeout profiles (fast vs deep).
- S3 — Collapse Phase 3+4 when no ❌ and single round.
- S4 — Cap Track 1 fix-up at 1 cycle → escalate/upgrade.

**COULD**
- C1 — Formalize lead+subagent fast-refutation review (Track 1).
- C2 — Role-based onboarding section + `parley new` track auto-pick.
- C3 — Default Track 2 implementers to worktrees.
