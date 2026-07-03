---
agent: antigravity-1
idea: meta-protocol-change-devx-speed
round: 1
date: 2026-07-03
---

# Round 01 Analysis — antigravity-1

The Parley Deck protocol is a robust safety framework, but its current form suffers from "monolithic ceremony"—applying the same high-friction process to every change regardless of risk. To achieve the goals of developer usability (G1) and speed (G2), we must transition from a flat protocol to a **multi-tiered, context-aware orchestration engine**.

## 1. Time-sink Audit
The following areas represent the largest wall-clock "taxes" in the current protocol:

1.  **Sequential Round dependencies (§4 Phase 1-2):** Even though Round 1 is independent, Round 2 cannot start until the *slowest* participant finishes Round 1. This "lowest common denominator" speed limits the whole deck.
2.  **Consensus/Finalization split (§4 Phase 3-4):** The gap between reaching consensus and the initiator drafting `FINAL.md` is a manual handoff that often stalls. Drafting a separate `FINAL.md` when `consensus.md` already contains the "Agreed decisions" is redundant for ~80% of tasks.
3.  **Mandatory Quorum Liveness (§9.0):** Probing every rostered agent before every idea adds 1-2 minutes of overhead to the "inner loop" of a developer just trying to get a quick sign-off.
4.  **Review Loop Redundancy (§4 Phase 6-8):** The implementation review cycle mirrors the design cycle exactly. For small changes, this "double deliberation" is the primary reason the protocol feels "too long."
5.  **Append-only Signoff Bottleneck (§4 Phase 3/7):** Requiring agents to append signoffs sequentially (to avoid git conflicts) turns a parallel agreement into a sequential queue.

## 2. Tiering Model: The "Speed Tracks"
I propose three tracks based on **Impact x Reversibility**:

| Track | Name | Objective Triggers | Process Deltas |
| :--- | :--- | :--- | :--- |
| **T1** | **Fast-Track** | Single file, <50 LOC, non-logic (docs/styles), or `revertible: high`. | Collapse Design Phase. 1 Reviewer. Auto-merge on Approval. |
| **T2** | **Standard** | Multi-file features, logic changes, no security/API surface. | Parallelized rounds. Consensus = Final. |
| **T3** | **Deliberation** | Protocol changes (§7), Security-critical, API-breaking, `reversibility: low`. | Full §4 Phase 0-8. Strict Gate (§4 Phase 8). |

**Track Transition Rules:** An idea defaults to T2. Any participant can "escalate" a T1 to T2 or T2 to T3 in Round 1 by citing a specific risk.

## 3. DevX: The "5-Minute" Entry Point
To make the protocol usable for developers (G1):

*   **The "Protocol Header" Quickstart:** Move the 1000 lines of `COOPERATION.md` below a `## Quickstart` section that fits on one screen.
*   **Default Roster:** If `participants:` is omitted in `00-prompt.md`, default to `[facilitator, random_peer]`.
*   **The `parley run` abstraction:** Developers shouldn't handle `round-01/` directories manually. A single command should scaffold, invoke participants, and present a diff-based consensus.
*   **Role-Appropriate Views:** Create a "Reviewer's Cheat Sheet" and an "Implementer's Guide" so non-experts only read what applies to their current phase.

## 4. Speed Levers
*   **Parallel Consensus:** Instead of waiting for a draft, agents write their "Signoff-if-no-changes" directly into Round 2.
*   **Optimistic Finalization:** For T1/T2, the `consensus.md` draft *is* the `FINAL.md`. If no blocks occur within N minutes, it is considered final.
*   **Quorum Waiving:** For T1, allow a "Quorum of 2" (Implementer + 1 Reviewer) even if the roster is larger.
*   **Driver-Led Implementation:** Let the facilitator agent auto-drive Phase 5-8 if the plan is "Green" (no findings).

## 5. Modern Agentic Concepts
*   **Spec-Driven Development:** Treat `FINAL.md` as an executable contract. If implementation deviates, the "Implementation Review" becomes a "Spec-Violation Check."
*   **Context Engineering:** Instead of "Read the whole repo," use **Progressive Disclosure**. The facilitator provides only the relevant "context-slices" to participants.
*   **Lead-Agent + Subagent Orchestration:** One "Lead" agent (facilitator) does the heavy lifting; "Checker" agents (participants) only verify specific constraints. This reduces "Thinking" costs significantly.
*   **Verification Gates:** Replace human-like "Review Prose" with **automated verification commands** (`checks:` in `00-prompt.md`, line 159). If the check passes, the "Minor/Nit" review is skipped.

## 6. What MUST Stay
*   **Non-Solo Requirement (§1):** The "Two-Keys" safety property is non-negotiable.
*   **Durable Audit Trail (§1):** Files must remain canonical. PR comments are ephemeral; the `parley-deck/` history must survive.
*   **Refutation-Default Review (§4 Phase 6):** Reviewers must continue to assume the code is broken until proven otherwise.

## 7. Prioritized Shortlist
1.  **MUST:** Implement **Tiering Model** (Impact: High / Cost: Low).
2.  **MUST:** Parallelize **Consensus/Finalization** (Impact: High / Cost: Low).
3.  **SHOULD:** Add **Quickstart/Role-based docs** (Impact: Medium / Cost: Medium).
4.  **SHOULD:** Enable **Quorum Scaling** (fewer agents for small tasks) (Impact: High / Cost: Medium).
5.  **COULD:** Integrate **Automated Verification Gates** (Impact: High / Cost: High).
