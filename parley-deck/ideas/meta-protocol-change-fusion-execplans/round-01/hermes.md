---
agent: hermes
idea: meta-protocol-change-fusion-execplans
round: 1
date: 2026-06-18
---

# Hermes Round-01 Analysis: OpenRouter Fusion + OpenAI ExecPlans for parley-deck

## Independent Summary

parley-deck is already a mature multi-agent deliberation + artifact protocol with strong invariants: append-only signoffs, one-file-per-agent-per-round, canonical artifacts (consensus.md, FINAL.md, IMPLEMENTATION.md), graded review severities (CRITICAL/MAJOR/MINOR/NIT), §13 parley retro, and human-gated auto-drive. Many Fusion/ExecPlan ideas have direct or stronger analogs here. I evaluate only concepts that could improve the existing structure without violating invariants.

## Per-Concept Verdicts

### 1. Judge "compare-don't-merge" structured analysis (Fusion)
- **Verdict: ADAPT**
- parley-deck's consensus.md is a drafted synthesis. The explicit blind-spots dimension ("issues no model addressed") is a genuine gap; our reviews catch errors in what was written but not systematic omissions across agents. Adding a lightweight "Blind Spots" subsection to consensus.md (populated during Phase 2 cross-review) would surface this without replacing the human-readable synthesis.
- Rationale: Improves coverage without ceremony. Fits naturally into existing review flow.
- Rough cost: Low (one new subsection, no new files or agents).

### 2. Confidence-by-breadth-of-agreement (Fusion)
- **Verdict: REJECT**
- parley-deck uses binary ✅ ACCEPT / ❌ BLOCKER signoffs that are append-only and agent-specific. Graded confidence would require either numeric scores (new ceremony) or post-hoc counting that conflicts with our one-file-per-agent model. Severities already encode strength of concern.
- Rationale: Net-negative; adds overhead without improving the human-gate or retro evidence.

### 3. Synthesis-as-distinct-value (Fusion)
- **Verdict: REJECT**
- The dedicated judge role is already approximated by the facilitator drafting consensus.md after independent round-01/02 files. Adding another synthesis agent would duplicate effort and risk recursion (explicitly avoided in Fusion). Our append-only signoffs already separate drafting from ratification.
- Rationale: Redundant with current Phase 2 cross-review + consensus step.

### 4. Negative-weight / "confidently wrong" lens (Fusion)
- **Verdict: ADAPT (minor)**
- Our CRITICAL severity already punishes high-impact confident errors. §13 retro already mines for systemic issues. A small explicit "Confident Error" tag in review comments could feed retro better, but only if it emerges naturally from existing reviews.
- Rationale: Low-value standalone; worth a lightweight convention if it reduces retro mining effort.

### 5. ExecPlan self-containment for FINAL.md / IMPLEMENTATION.md
- **Verdict: ADAPT**
- Making FINAL.md and IMPLEMENTATION.md explicitly self-contained living documents (so a fresh headless agent or the auto-drive driver can resume from the artifact alone) directly addresses the cross-invocation resumption the driver needs. This aligns with ExecPlan's core goal without adding new files.
- Rationale: High overlap with auto-drive needs; strengthens the existing handoff artifacts.
- Rough cost: Medium (discipline to keep sections complete; no new structure).

### 6. Living sections (Progress, Surprises & Discoveries, Decision Log, Outcomes & Retrospective) → §13 retro
- **Verdict: ADAPT**
- Embedding Decision Log and Surprises & Discoveries directly into IMPLEMENTATION.md would provide richer, timestamped evidence for the parley retro than current mining. Progress checklist with ISO timestamps is already partially present in our workflow.
- Rationale: Improves retro signal quality with minimal new overhead.
- Rough cost: Low (add two subsections to IMPLEMENTATION.md template).

### 7. Idempotence & Recovery section → auto_implement safety
- **Verdict: ADAPT**
- An explicit "Idempotence & Recovery" section in FINAL.md would harden the gated auto-drive phase (clean-tree, no-land). This maps cleanly to existing safety checks without new mechanisms.
- Rationale: Directly strengthens the most fragile part of the driver loop.
- Rough cost: Low (one new required section).

### 8. Behavior-focused acceptance criteria
- **Verdict: ADAPT**
- Requiring observable acceptance criteria in FINAL.md (test commands + expected outputs) that both Phase 6 review and the driver must verify would reduce "compiles but does nothing" risk. Already partially implied by review severities.
- Rationale: Improves validation discipline; fits existing artifact.
- Rough cost: Low.

### 9. Cost/overhead realism
- Fusion's 4-5× multiplier is unacceptable on top of parley-deck's existing multi-agent tax. ExecPlan-style living sections add maintenance overhead but are bounded to two artifacts and pay for themselves in driver reliability and retro quality.

## Prioritized Top-3 Recommendations

1. **ADAPT: Self-containment + living sections in FINAL.md / IMPLEMENTATION.md** (highest leverage)
   - Rationale: Directly improves auto-drive resumption and §13 retro evidence quality. Highest alignment with parley-deck's core handoff problem.
   - Cost: Medium (ongoing discipline) but bounded.

2. **ADAPT: Blind-spots subsection in consensus.md + lightweight idempotence section in FINAL.md**
   - Rationale: Fills the single clear coverage gap from Fusion while hardening driver safety.
   - Cost: Low.

3. **ADAPT: Behavior-focused acceptance criteria convention**
   - Rationale: Low-cost improvement to validation that compounds with the above.
   - Cost: Low.

## What to Explicitly NOT Adopt

- Any form of parallel model panel or dedicated judge agent (conflicts with one-file-per-agent, append-only signoffs, and human gate).
- Numeric confidence scoring or breadth-of-agreement weighting (breaks binary signoff model).
- Full Fusion pipeline or recursion-protection machinery (unnecessary overhead; parley-deck is already single-level and human-gated).
- New artifact types or living-document requirements outside FINAL.md and IMPLEMENTATION.md (would violate canonical artifact discipline).

These recommendations stay strictly within design hypotheses for human review. No protocol text changes proposed.