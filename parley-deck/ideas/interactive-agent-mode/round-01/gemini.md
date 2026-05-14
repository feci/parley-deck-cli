---
agent: gemini
idea: interactive-agent-mode
round: 1
date: 2026-05-14
---

## Summary

The `interactive-agent-mode` slice introduces explicit launch modes for agents in `parley-deck-cli`. By differentiating between `headless`, `interactive`, and `manual` modes, the CLI can better support diverse agent capabilities and provider billing models (specifically Claude's upcoming Agent SDK changes). `headless` remains the default for automation, while `interactive` and `manual` provide high-signal handoffs for human-driven agent sessions. The goal is to preserve the Parley Deck ownership rules while making the handoff between automated orchestration and manual execution safe and auditable.

## Proposed approach

1.  **Launch Mode Configuration:**
    Extend the `agents.Spec` struct and the `agents.toml` schema to include a `launch_mode` field.
    -   `headless`: (Default) Standard programmatic execution via `headless_args`.
    -   `interactive`: Parley prepares the prompt and environment, displays a clear handoff instruction, and waits (polls) for the agent to write the target artifact before proceeding with validation.
    -   `manual`: Parley prepares the environment/prompt, prints instructions, and exits. This is the "preparation-only" mode.

2.  **Interactive Handoff Workflow:**
    When a participant is triggered in `interactive` mode:
    -   **Preparation:** Parley writes the specific task prompt to a file (e.g., `runs/{run_id}/agents/{agent_id}/prompt.md`).
    -   **Handoff:** Parley prints a high-visibility instruction to the terminal (e.g., `[HANDOFF] Agent 'claude' requires interactive session. Task: claude -p <path_to_prompt>`).
    -   **Observation:** The runner enters a polling loop, checking for the creation/modification of the `outputPath`. 
    -   **Validation Loop:** Once the artifact is detected, Parley runs the standard protocol validation. If validation fails (e.g., missing required sections), Parley reports the errors to the terminal and resumes polling, allowing the user to correct the artifact in the interactive session.

3.  **Command-level Overrides:**
    Add a global `--mode` flag (e.g., `parley run --mode claude=interactive`) to override configured defaults. A shorthand `--interactive` could set all participants to interactive mode for a specific run.

4.  **Auditability & Visibility:**
    -   Include the launch mode in the `agent.started` and `agent.completed` events in `events.jsonl`.
    -   Surface the mode in `parley run --dry-run` and the initial launch summary so users aren't surprised by a pause or an exit.

5.  **Resuming from Manual/Interactive:**
    Leverage existing "skip if artifact exists" logic. If a `manual` run was performed, a subsequent standard run will detect the artifact, skip invocation, and perform the final validation/event logging.

## Concerns / open questions

-   **Sequential vs. Parallel:** Interactive handoffs should almost certainly be sequential to avoid terminal confusion, even if the runner usually supports parallel execution.
-   **Terminal multiplexing:** Should `parley` attempt to launch the interactive command itself (e.g., `exec.Command` with `Stdin/Stdout/Stderr` tied to the parent TTY)? This would be more seamless but might conflict with the "no PTY manager" constraint.
-   **Handoff Prompt format:** Should we provide a raw text prompt, or a pre-formatted command string that the user can copy-paste?
-   **Timeouts:** `headless` timeouts should not apply to `interactive` sessions. We might need a separate `interactive_timeout_ms` or simply disable it.

## Risks

-   **Polling Efficiency:** Excessive polling could impact system performance, while slow polling might feel laggy to the user.
-   **User Friction:** If the handoff requires many manual steps (find file, run command, wait), users might find it too cumbersome compared to `headless`.
-   **Validation Strictness:** If the protocol validation is too pedantic, users in `interactive` mode might find it frustrating to repeatedly fix minor formatting issues to satisfy the "wait" loop.
-   **Security:** Ensure that the temporary prompt files created for handoffs are properly cleaned up or secured, as they may contain sensitive project context.
