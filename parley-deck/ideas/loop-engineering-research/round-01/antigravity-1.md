---
agent: antigravity-1
idea: loop-engineering-research
round: 1
date: 2026-06-22
---

## Summary

From the perspective of **Risk & Guardrails**, loop engineering represents a double-edged sword. While it scales generation throughput, it risks introducing catastrophic compounding failure modes if automated without rigid constraints. As Parley transitions from human-driven turn-by-turn prompts to autonomous loops, we must enforce structural circuit breakers against **runaway costs**, **human-in-the-loop (HITL) fatigue**, **comprehension debt**, and **cognitive surrender**. 

Our core stance is that **Parley must refuse to automate final code merges, roster modifications, budget overrides, and safety-critical escalations**. We propose adding hard token/cost budget gates to the [driver](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/driver.go), formalizing an adversarial refutation mode for Phase 6 reviewers, and restricting event-triggered outer loops to kickoff and initial drafting.

---

## Proposed approach

To safely integrate loop engineering principles into the Parley Deck protocol ([COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md)) and the [parley-deck-cli](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli) engine, we recommend a "fail-closed" architectural approach:

1. **Hard Operational Budgets**: The [driver.Config](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/driver.go#L41-L63) must enforce explicit limits on dollar cost, wall-clock duration, and iterations. If any budget is exceeded, the loop must write a blocking escalation and halt.
2. **True Model Disjunction**: The protocol must enforce that the checker (Phase 6 reviewer) utilizes a different model instance (and preferably from a different family) than the maker (Phase 5 implementer) to eliminate collusive confirmation bias.
3. **Restricted Trigger Boundaries**: Automated cron or event triggers (e.g., CI failures or issue creation) may only kickoff an idea (Phase 0) and generate independent analyses (Phase 1). Moving to implementation (Phase 5) or completing a review (Phase 7) must require an explicit human command.
4. **Adversarial Refutation Prompting**: Reviewers must be prompted with explicit instructions to find fault, assuming the implementation is incorrect until proven otherwise.

### Prioritized Recommendations

| # | Recommendation | Target File / Struct / Section | Loop Engineering Principle | Effort | Risk | Call | Spin-off Idea Slug |
|---|---|---|---|---|---|---|---|
| 1 | **Hard Resource Budgets** | [driver.Config](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/driver.go#L41) & [CreateOptions](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runcontrol/runcontrol.go#L17) | Stopping conditions / budget limits | Medium | Low | **Adopt** | `driver-hard-resource-budgets` |
| 2 | **Adversarial Refutation Review & Model Disjunction** | [COOPERATION.md §4.Phase 6](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md#L353) & [runner.Options](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/runner.go#L22) | Maker/Checker separation | Low | Low | **Adopt** | `adversarial-refutation-checker` |
| 3 | **Restricted Trigger Boundaries** | [COOPERATION.md §11](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md#L663) & [driver.Run](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/loop.go#L23) | Automations / Trigger boundaries | High | Med | **Adapt** | `trigger-boundary-confinement` |
| 4 | **HITL Batching & Risk Gating** | [hitl.Question](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/hitl/hitl.go#L30) | Fatigue mitigation | Med | Low | **Adopt** | `hitl-batching-and-gating` |
| 5 | **Cross-Run State Ledger** | `parley-deck/STATE.md` (new) | Durable State | Low | Med | **Reject** | N/A |

### Rationale for Top Items

*   **Hard Resource Budgets (Adopt)**: Runaway execution in automated loops can quickly exhaust token quotas and pile up hundreds of dollars in API charges. We must introduce `MaxCostUSD` and `MaxDuration` fields into the driver, tracked in the session database.
*   **Adversarial Refutation (Adopt)**: The current protocol relies on social constraints for Phase 6 review. If the same agent runs implementation and review under the same model, we observe high rates of false-green reviews. Hardcoding model disjunction ensures fresh, critical eyes on the diff.
*   **Restricted Trigger Boundaries (Adapt)**: Completely rejecting automation limits throughput, but naive automated loops that merge code autonomously are unsafe. By restricting triggers to Phase 0/1 kickoff, we let agents asynchronously prepare proposals, while keeping the implementation and landing controls strictly manual.
*   **Cross-Run State Ledger (Reject)**: A centralized `STATE.md` tracking cross-run backlogs introduces a shared mutable file that increases context window bloat and write conflict risks. The existing [parley-tracker](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/COOPERATION.md) integration is already sufficient for mapping external tickets; we should avoid creating duplicate on-disk state.

---

## Concerns / open questions

1.  **Approval Fatigue Threshold**: How do we measure the point where human operators suffer "approval fatigue" and begin blindly signing off? Can the CLI engine detect high-frequency human inputs or consecutive approvals without diff inspections, and dynamically increase verification requirements?
2.  **Comprehension Debt Accrual**: When multiple agents run parallel iterations in independent worktrees, the volume of changes can outpace human reading speed. How can Parley measure and visualize "comprehension debt" before allowing a merge?
3.  **MCP Sandbox Leakage**: If MCP plugins have write or run access in the main repository, how do we prevent malicious or broken agent code from executing dangerous side-effects during automated runs?

---

## Risks

The introduction of an automated outer loop exposes Parley to several critical risks:

### 1. Runaway Cost
In a fully automated outer loop, failure is often silent. An infinite loop can occur if:
*   A test is flaky, causing the implementer and reviewers to repeatedly patch, review, and break code.
*   A prompt template contains a formatting contradiction that triggers continuous validation failures.
Without hard stopping conditions (e.g. tracking token count and dollar spend inside [store.Event](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/store)), a single unresolved bug could execute continuously, compounding API token spend.

### 2. HITL Approval Fatigue
When the [hitl.Store](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/hitl/hitl.go#L43) queue is flooded with requests for execution permission, confirmation boxes, and code reviews, the human gatekeeper's vigilance degrades. Approving tool calls or reviewing files becomes a mechanical, low-friction gesture, effectively rendering the human-in-the-loop gate useless.

### 3. Comprehension Debt
If agents run automated fix-up rounds through the driver, code can be refactored across several cycles without a human actively following the rationale. The developer inherits a codebase where the *what* might pass tests, but the *why* is buried in deep agent transcripts, leading to architectural drift and loss of code ownership.

### 4. Cognitive Surrender
Relying on loops to handle both generation and correction encourages developers to stop thinking critically about edge cases. If developers believe "the loop will find and fix it," they surrender their role as the primary designer, shifting their task to blindly reviewing automated diffs.

### What Parley Must Refuse to Automate
To preserve the safety and integrity of the repository, Parley must structurally refuse to automate:
*   **PR Merging and Releases**: The final step of merging a design branch or code implementation branch into the main integration branch must never be automated by the driver or any MCP connector.
*   **Roster Modification**: Roster updates ([COOPERATION.md §2](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md#L70)) must remain manual; automated loops must not add new agent participants or change existing agent models/thinking modes.
*   **Consensus Overrides**: The consensus phase must never auto-complete on behalf of a human participant, nor bypass a blocking review without a recorded human operator ruling.
