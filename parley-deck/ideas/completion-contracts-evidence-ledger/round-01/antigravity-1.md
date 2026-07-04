---
agent: antigravity-1
idea: completion-contracts-evidence-ledger
round: 1
date: 2026-07-04
---

## Summary

The proposal to introduce `done_when:` completion contracts and a `review/evidence.md` ledger into the Parley Deck protocol addresses a critical trust gap: the reliance on an implementer's prose claim of success. From a DevX perspective, this shifts the protocol from "trust-but-verify" (manual verification) to "evidence-first" (grounded verification). For a developer running Parley in their own project, this provides a definitive receipt of quality that reduces the review burden while increasing confidence in model-driven changes.

## Proposed approach

1.  **Contract-Driven Kickoff (`done_when:`):**
    *   The `done_when:` block in `00-prompt.md` should be the primary entry point for defining quality gates. 
    *   **DevX Recommendation:** The driver should offer a `parley check-contract` command to validate that the commands listed in the contract are actually executable in the local environment before the agents start work.

2.  **Automated Evidence Capture (`review/evidence.md`):**
    *   The ledger should be written by the driver (via `parley fix-up`) to ensure objectivity.
    *   **Structure:** Each entry should include a "Verification Context" (e.g., git hash, timestamp) and a "Result Matrix" (command, exit code, duration, and a safe, truncated output digest).
    *   **Privacy by Construction:** The driver must automatically scrub common secret patterns (tokens, keys) from the output digest before writing to the ledger.

3.  **Phase 8 Enforcement:**
    *   The driver MUST prevent marking Phase 8 as complete if a `done_when:` contract exists but the latest ledger entry shows failures.
    *   This enforces a "hard gate" that prevents rubber-stamping broken code, which is especially valuable in high-velocity projects.

## Concerns / open questions

*   **Config Drift:** If the project's build system changes (e.g., moving from `npm` to `pnpm`) during a long-lived idea, the `done_when:` contract might become stale. Should the contract allow for "dynamic resolution" of commands?
*   **Context Limit Sensitivity:** Recording full output in `review/evidence.md` could quickly bloat the file. The "output digest" needs a strictly defined size limit (e.g., first 50 lines + last 50 lines + error count) to remain readable for both humans and agents.
*   **The "Flaky Test" Paradox:** How does the driver handle flaky tests in the `done_when:` contract? If a check is non-deterministic, it could block the idea indefinitely.

## Risks

*   **Contract Brittle-ness:** Over-specifying completion criteria at Phase 0 can lead to frustration if the implementation path requires a different verification strategy. There should be a protocol-approved way to amend the `done_when:` block via a round-NN update.
*   **The "Prose vs. Evidence" Conflict:** If an agent provides prose explaining why a failed check in the ledger is "actually fine," it could confuse the consensus process. The protocol must clearly state that evidence failures require a contract amendment or a code fix, not just a prose rebuttal.
*   **Trust Model Escalation:** While the prompt assumes the same trust level as `RunChecks`, a `done_when:` contract might entice authors to include more complex (and potentially risky) commands. The driver should warn when a contract contains commands not previously seen in the project's `RunChecks` history.
