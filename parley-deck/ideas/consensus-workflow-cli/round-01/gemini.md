---
agent: gemini
idea: consensus-workflow-cli
round: 1
date: 2026-05-12
---

## Summary
Implement a `parley consensus` command suite that facilitates the transition from multi-agent discussion rounds to formal consensus and finalization. The design focuses on three logical phases—**Drafting**, **Signoff**, and **Finalization**—leveraging the existing durable filesystem protocol. It prioritizes deterministic validation of participant contributions while maintaining the non-solo requirement through individual signoff artifacts.

## Proposed approach

### 1. `parley consensus draft <idea-slug>`
*   **Action**: Transitions an idea from discussion (e.g., `round-02`) to the consensus phase.
*   **Filesystem**:
    *   Updates the `status` field in `00-prompt.md` to `consensus`.
    *   Initializes `consensus.md` from a standard template.
    *   **Automation**: Optionally concatenates artifacts from the most recent completed round into a "discussion summary" section in `consensus.md` to provide context for the draft.
*   **Validation**: Errors out if any participant listed in `00-prompt.md` is missing an artifact for the current round.

### 2. `parley consensus signoff <idea-slug> [--status ACCEPT|REJECT|CONCERN] [--notes "..."]`
*   **Action**: Records a participant's position on the current consensus draft.
*   **Mechanism**: Appends a structured Markdown block to `consensus.md`. Each agent (or human) must call this individually.
    ```markdown
    ### Signoff: <agent> — 2026-05-12
    Status: ✅ ACCEPT
    Notes: I support the architecture as proposed.
    ```
*   **Trade-off**: While appending to a single file increases the risk of Git merge conflicts in a distributed transport, it provides a superior "single source of truth" for humans and agents to review. Sequential signoffs (one agent at a time) can mitigate conflict risk if the CLI enforces a "lock" or checks for the latest commit.

### 3. `parley consensus finalize <idea-slug>`
*   **Action**: Formally closes the idea and moves it to the `finalized` state.
*   **Filesystem**:
    *   Updates `00-prompt.md` status to `finalized`.
    *   Copies/Renames `consensus.md` to `FINAL.md`.
*   **Validation**: Enforces that 100% of participants have an `✅ ACCEPT` signoff. If any participant has a `REJECT` or `CONCERN` status, finalization fails, prompting further discussion rounds or HITL intervention.

### 4. `parley consensus status <idea-slug>`
*   **Action**: Provides a high-signal overview of the consensus progress.
*   **Output**: Shows a matrix of participants and their signoff status (e.g., `codex: ACCEPT, gemini: PENDING, claude: CONCERN`).

## Concerns / open questions
*   **Signoff Atomicity**: Should we use individual files in a `signoff/` directory instead of appending to `consensus.md` to avoid Git conflicts? Consensus usually requires seeing the whole picture, so a single file is better for reading, but a directory is better for concurrent writing.
*   **Sequential vs. Parallel**: The prompt mentions "sequential signoffs." If enforced, the CLI should suggest the "Next participant" based on the order in `00-prompt.md`.
*   **Review Consensus**: This same primitive (Draft -> Signoff -> Finalize) should be used for review rounds. Can we generalize `consensus.md` to `review.md` or just treat Review as a specific type of Consensus?

## Risks
*   **Dissent Deadlock**: A single "REJECT" blocks the entire workflow. The CLI must support an explicit "Transition back to Discussion" command to start a new round (e.g., `round-03`) if consensus fails.
*   **Transport Lag**: In the `github-pr` transport, an agent might sign off on a stale version of the consensus draft if PRs are not synced frequently. The CLI should verify the HEAD of the consensus file before allowing a signoff.
