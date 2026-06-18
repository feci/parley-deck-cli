---
agent: agy
idea: meta-protocol-change-fusion-execplans
round: 1
date: 2026-06-18
---

# Inspiration from OpenRouter Fusion + OpenAI ExecPlans for the parley-deck Protocol

## Summary
This analysis evaluates how concepts from OpenRouter Fusion (multi-model panel deliberation and judge-based structured synthesis) and OpenAI ExecPlans (living, self-contained execution documents) can enhance the `parley-deck` protocol. While `parley-deck` already possesses a more robust multi-agent dialogue and consensus mechanism, incorporating ExecPlan-style living sections into `IMPLEMENTATION.md` and behavior-focused validation into `FINAL.md` will significantly benefit the auto-drive driver's resumption capabilities. Conversely, we reject Fusion's single-model judge synthesis and graded confidence models to preserve `parley-deck`'s core invariants: strict veto-based consensus and independent multi-agent verification.

---

## Concepts and Verdicts

### 1. OpenRouter Fusion Concepts

#### 1.1 Judge "Compare-Don't-Merge" & Structured Analysis
*   **Concept:** A judge model evaluates panel responses and structures the output into consensus, contradictions, partial coverage, unique insights, and blind spots (issues no model addressed).
*   **Verdict:** **ADAPT**
*   **Rationale:** `parley-deck` resolves contradictions during dialogue rounds and drafts consensus in `consensus.md`. However, `consensus.md` currently lacks a structured format for identifying what was *omitted* or *contested*. Adapting the **"Blind Spots"** and **"Contradictions / Alternative Perspectives"** dimensions into `consensus.md` will strengthen the design audit trail.
*   **Cost/Benefit:** Low Cost (simple markdown template update). High Benefit (forces agents to explicitly search for omissions and document rejected alternatives).

#### 1.2 Confidence-by-Breadth-of-Agreement
*   **Concept:** Graded confidence based on how many models in a panel agree on a specific point.
*   **Verdict:** **REJECT**
*   **Rationale:** `parley-deck` relies on a strict veto-based consensus mechanism where any active participant's `❌ BLOCK` halts progression until addressed. Introducing graded confidence or majority-rule thresholds would weaken the safety invariants of the protocol, potentially allowing a majority of agents to override a minority agent's critical warning.
*   **Cost/Benefit:** Net-Negative.

#### 1.3 Synthesis-as-Distinct-Value
*   **Concept:** A dedicated judge step adds synthesis value even when fusing a single model with itself.
*   **Verdict:** **REJECT**
*   **Rationale:** In `parley-deck`, synthesis is performed by the initiator (or a volunteer) drafting `consensus.md` and `FINAL.md`. Because all participants must review, sign off, or block the drafted consensus, bias is already mitigated through multi-agent critique. Mandating a separate "judge/synthesizer" role/phase would introduce unnecessary coordination overhead.
*   **Cost/Benefit:** Low Benefit / High Ceremony.

#### 1.4 DRACO "Confidently Wrong" Penalty
*   **Concept:** Penalizing models that confidently state wrong things to improve benchmark hygiene.
*   **Verdict:** **ADAPT** (into Section 13 Retrospective)
*   **Rationale:** Currently, `parley-deck` does not track when an agent acts as a "false-alarm blocker" (i.e., raising a `CRITICAL` or `MAJOR` review finding that is subsequently dismissed as invalid). By adapting this concept into the `parley retro` (Section 13) analysis, we can identify agents with high "confident-wrongness" rates to refine their prompts or adjust roster selection.
*   **Cost/Benefit:** Low Cost. Medium Benefit.

---

### 2. OpenAI ExecPlans (PLANS.md) Concepts

#### 2.1 Living Execution Tracking in IMPLEMENTATION.md
*   **Concept:** A living design/implementation document containing a timestamped Progress log, Surprises & Discoveries, and a Decision Log, enabling stateless resumption.
*   **Verdict:** **ADAPT**
*   **Rationale:** In `parley-deck`, `FINAL.md` remains static to preserve the consensus snapshot, while `IMPLEMENTATION.md` tracks implementation progress. Currently, `IMPLEMENTATION.md` is updated at milestones, but lacks the granular tracking needed for a stateless agent or the **auto-drive driver** to resume execution seamlessly after a timeout or crash. Adding a living `## Progress` log (with ISO timestamps), `## Surprises & Discoveries` (with evidence), and a local `## Decision Log` to `IMPLEMENTATION.md` will solve this.
*   **Cost/Benefit:** Medium Cost (requires updating the implementer's update cycle logic). High Benefit (essential for auto-drive driver robustness and resumption).

#### 2.2 Self-Containment of Context
*   **Concept:** Embedding all required signatures, code snippets, and dependencies directly in the document so no external context is required.
*   **Verdict:** **ADOPT** (as a best-practice guideline for `FINAL.md`)
*   **Rationale:** This aligns with conflict-avoidance rule 6.4 (copying snippets due to workspace boundaries). Formalizing this as a requirement for `FINAL.md` ensures that the implementer (who may run in a separate sandboxed workspace) has all necessary interfaces and constraints on hand.
*   **Cost/Benefit:** Low Cost. High Benefit.

#### 2.3 Behavior-Focused Validation & Acceptance
*   **Concept:** Mandating exact verification commands and expected outputs (e.g., HTTP response codes, CLI outputs) rather than simple compilation checks.
*   **Verdict:** **ADOPT**
*   **Rationale:** Currently, `FINAL.md` does not strictly require observable acceptance criteria, and `IMPLEMENTATION.md` only has a generic "Checks to run" section. Mandating a `## Validation and Acceptance` section in `FINAL.md` allows both human/agent reviewers in Phase 6 and the auto-drive driver to test the implementation against an objective, execution-based standard.
*   **Cost/Benefit:** Low Cost (template update). High Benefit (enables automated verification gates).

#### 2.4 Idempotence & Recovery Plan
*   **Concept:** Explicitly documenting how to safely run and retry commands without causing drift, and how to roll back changes.
*   **Verdict:** **ADAPT**
*   **Rationale:** While git handles recovery for pure code changes, ideas involving infrastructure, databases, or API side-effects require deliberate recovery planning. Requiring an `## Idempotence and Recovery` section in `FINAL.md` for high-risk changes will prevent the driver from corrupting external environments during retries.
*   **Cost/Benefit:** Low Cost. Medium Benefit.

---

## Prioritized Recommendations (Top 3)

### 1. Mandate Behavior-Focused Validation in `FINAL.md`
*   **Action:** Add a required `## Validation and Acceptance` section to the `FINAL.md` schema, containing runnable verification commands and expected outcomes.
*   **Rationale:** This bridges the gap between design and testing, giving reviewers and the auto-drive driver a clear, non-subjective gate to verify implementation success.
*   **Rough Cost:** Low (template schema update).

### 2. Transform `IMPLEMENTATION.md` into a Living ExecPlan
*   **Action:** Restructure `IMPLEMENTATION.md` to include living sections: `## Progress` (timestamped checklist), `## Surprises & Discoveries`, and `## Decision Log`.
*   **Rationale:** Gives stateless implementers and the auto-drive driver the necessary state history to resume work after failures, while generating structured evidence for retrospectives.
*   **Rough Cost:** Medium (updates to Phase 5–8 agent prompts to maintain the document state continuously).

### 3. Track "Confidently Wrong" Trajectories in Section 13 Retrospectives
*   **Action:** Update the `parley retro` specification to parse `review/consensus.md` and measure the frequency of dismissed `CRITICAL`/`MAJOR` findings per agent.
*   **Rationale:** Helps operators identify hallucinating or overly pedantic agents that slow down the consensus process without adding value.
*   **Rough Cost:** Low (retrospective analysis script update).

---

## What to Explicitly NOT Adopt

1.  **Do NOT make `FINAL.md` a Living Document:** In ExecPlans, the design document is updated continuously. In `parley-deck`, `FINAL.md` must remain static once consensus is reached. Changing `FINAL.md` during execution destroys the audit trail of what was actually agreed upon. Any deviations must be logged in `IMPLEMENTATION.md` instead.
2.  **Do NOT adopt Graded/Majority Consensus:** We must preserve the binary append-only signoff with veto power. Majority-based or confidence-based consensus risks ignoring correct warnings from minority models.
3.  **Do NOT adopt Single-Model Judge Synthesis:** The facilitator or driver must never act as a single-model judge to bypass multi-agent rounds. Synthesis must remain a collaborative draft reviewed by all participants.
