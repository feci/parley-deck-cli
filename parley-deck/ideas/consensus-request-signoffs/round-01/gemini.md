---
agent: gemini
idea: consensus-request-signoffs
round: 1
date: 2026-05-13
---

## Summary
The goal is to implement the `parley consensus request-signoffs` command to automate the collection of consensus or review signoffs from configured headless agents. This command will identify missing participants from an idea's consensus document, sequentially invoke each agent with a targeted request to review and sign off, and verify the integrity of the consensus file after each step. This slice prioritizes sequential execution and explicit human confirmation before engaging external/hosted agents.

## Proposed approach
- **Command implementation:** Add `request-signoffs` to `internal/app/app.go` under the `consensus` command group. Support `--review`, `--participants`, `--yes`, and `--dry-run` flags.
- **Missing participant resolution:**
    - Use `internal/consensus.Status` to parse the target consensus file (`consensus.md` or `review/consensus.md`).
    - Resolve the list of participants who have not yet signed.
    - If `--participants` is provided, filter or validate against this list.
- **Sequential Runner:**
    - Iterate through target agents one by one.
    - For each agent:
        - Check if it is a hosted/external agent (e.g., Gemini, Claude, Codex) and require `--yes` unless `--dry-run` is active.
        - **Signoff Prompt Generation:** Create a specialized prompt for the agent that includes:
            - The full content of the current consensus document.
            - The original idea prompt (`00-prompt.md`) and any relevant round artifacts.
            - Explicit instructions to append a canonical signoff block (Status, Notes, Counter-proposal) directly to the consensus file.
        - **Agent Execution:** Invoke the agent using the existing `runner.CommandFor` infrastructure, ensuring `{root}` is passed so the agent can access and modify the consensus file.
        - **Post-invocation Validation:** Immediately after the agent exits, call `internal/consensus.Status` again.
            - Verify that the agent appended exactly one valid signoff block.
            - Stop the entire sequence if any error occurs, if the signoff is malformed, or if the agent records a `❌ BLOCK`.
- **Dry-run:** Print the intended execution plan, including the list of agents, the order of invocation, and the target file path.

## Concerns / open questions
- **Append logic vs. Full Rewrite:** Should we encourage agents to only append, or are they likely to rewrite the whole file? The `internal/consensus` parser is robust, but a full rewrite might accidentally lose other agents' signoffs if the agent doesn't have the latest content. We should emphasize "append only" in the prompt.
- **Agent Concurrency:** While the CLI executes agents sequentially, there's no locking mechanism. We should warn the user not to manually edit the consensus file while the command is running.
- **Hosted Agent Identification:** We should ensure the "hosted" check is consistent with `parley run` and uses the `ExternalBackend` field from discovery.

## Risks
- **Parsing Failures:** If an agent appends a malformed signoff, it might make the consensus file unreadable for subsequent steps. The "verify after each invocation" strategy is critical here.
- **Cost Accumulation:** Running multiple hosted agents in sequence can quickly consume tokens/quota. Users should be aware of this via the `--yes` confirmation.
- **"BLOCK" Interruption:** A single agent blocking can stop a long-running signoff request process, which is correct for consensus but might be frustrating if it happens at the end of a long queue.
