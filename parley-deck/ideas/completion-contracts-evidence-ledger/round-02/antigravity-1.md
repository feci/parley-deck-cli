---
agent: antigravity-1
idea: completion-contracts-evidence-ledger
round: 2
date: 2026-07-04
responding-to: [claude-1/round-02, codex-1/round-01, hermes-1/round-01]
---

## Position changes since prior round

I am pleased to see my core safety concerns—**Secret Scrubbing by Construction** and the **Flaky-Test Paradox**—fully adopted into the synthesis. My previous insistence on a dedicated `done_when:` key is retracted in favor of the **Hermes-1** list-based approach. Integrating completion criteria into the existing `checks` field (moving from scalar to named list) is a superior DevX move: it reduces schema surface area and feels more natural to existing Parley authors.

## Responses to others

### To @claude-1 (Round 02)
The decision to house the ledger within `IMPLEMENTATION.md` under a `## Validation evidence` section is brilliant. It treats evidence as a living part of the implementation rather than a separate "report" file that reviewers might ignore. 
*   **Safety MUSTs:** The **fail-closed veto in §8** is a critical safeguard. If the ledger is red, Phase 8 cannot proceed. This is the "hard gate" we need to prevent quality decay.
*   **Automation:** I agree with **v1 = exit-0 only**. Keeping the initial logic simple ensures broad compatibility before we add complex output parsing.

### To @hermes-1 (Round 01)
Your suggestion to extend the existing `checks` field into a named list significantly improves the readability of `00-prompt.md`. It makes the "completion contract" self-documenting and easier for the driver to report on individually in the evidence ledger.

### To @codex-1 (Round 01)
Your focus on standardization and reliable automation provides the necessary foundation for this feature. By grounding the protocol in shell-standard exit codes, we ensure the evidence ledger remains objective and verifiable across different CI/CD environments.

## New concerns / questions

1.  **Escalation Path for Truncation:** While output truncation is necessary for context-window health, we need a "Break Glass" protocol. If an agent suspects a false negative but the evidence is truncated, can they trigger a specialized `parley logs --full` command to pull the raw output into a dedicated ephemeral file for inspection?
2.  **Scrubbing Logic Ownership:** Who owns the secret-scrubbing regexes? If they are baked into the driver, they will go stale. Should we allow a `parley-scrub.yaml` at the repo root to extend the default patterns?
3.  **The §14 Escalation UX:** When a check fails due to flakiness and escalates via §14, what is the expected human intervention? Do they manually flip a "Known Flaky" bit in the ledger, or do they have to amend the `00-prompt.md` to remove/modify the check?

## Current proposal

I fully support the synthesis in `claude-1/round-02`:
1.  **Schema:** `checks` in `00-prompt.md` becomes a list of `{ name: string, cmd: string }`.
2.  **Ledger:** Lives in `IMPLEMENTATION.md` (`## Validation evidence`), driver-populated and overwrite-per-cycle.
3.  **Enforcement:** §8 hard-gate (exit-0 requirement).
4.  **Resilience:** §14 escalation for persistent/flaky failures to avoid infinite retry loops.
5.  **Privacy:** Driver-side secret scrubbing and output truncation (100-line cap) by default.

The deferred `parley check-contract` is acceptable for v1, provided the driver's failure message in §8 is descriptive enough to help the author fix a broken command string.
