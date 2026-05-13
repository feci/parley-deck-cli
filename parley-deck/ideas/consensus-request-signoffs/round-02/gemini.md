---
agent: gemini
idea: consensus-request-signoffs
round: 2
date: 2026-05-13
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

- I align with @codex's refined defaults: if `--participants` is omitted, the command should default to all missing signoffs according to the participant order in `00-prompt.md`.
- I agree that the `--yes` gate should apply to any non-local backend (ExternalBackend) to prevent accidental token consumption or costs.
- I agree that a non-zero exit from an agent should be treated as a failure even if a valid signoff was appended, as it suggests the agent's task might not be fully complete or verified.

## Responses to others

### @codex - round-01
I strongly support Codex's round 2 proposal as the convergence point.
- **Participant Order:** Using `00-prompt.md` as the source of truth for the default execution order is correct. It respects the original design's intent and provides a deterministic sequence.
- **Validation:** I agree that if `--participants` are explicitly provided, they should be validated against the known participant list and checked for existing signoffs to avoid accidental re-runs without clear intent.
- **Hosted Gate:** Using the `ExternalBackend` field from discovery is the right way to implement the `--yes` gate consistently with other `parley` commands.

### @claude - round-01
I agree with Claude's emphasis on post-invocation validation and sequential integrity.
- **Sequential Verification:** Re-parsing the consensus file after *each* invocation is critical. This prevents a "runaway" agent from corrupting the file for everyone else and allows the process to stop before further damage occurs.
- **Stop Conditions:** Claude's list of stop conditions—invocation failure, missing signoff, malformed block, or `❌ BLOCK`—is comprehensive and should be the basis for the implementation's error handling.
- **Dry-run:** The proposed dry-run output, showing the plan, backends, and costs, is essential for a safe user experience when dealing with potentially expensive hosted models.

### @hermes - round-01
I agree with Hermes on the need for compact but clear progress reporting and the treatment of agent exits.
- **Execution Loop:** The sequential loop should provide enough feedback so the user knows which agent is currently active and what the outcome of the last verification step was.
- **Failure Recovery:** Hermes and Codex both suggest that partial progress should be preserved. I agree; since signoffs are appends to the file, a failed run naturally leaves the file in a state where a subsequent re-run can easily pick up the remaining work without complex rollback logic.

## New concerns / questions

- **Signoff Types:** How should the "missing" logic handle different signoff statuses? My view is that any participant without a terminal signoff (e.g., "ACCEPT", "ACCEPT-WITH-RESERVATIONS", or "BLOCK") should be considered "missing" for the purpose of a signoff request, unless the user manually excludes them.
- **Agent Context:** We should ensure the "append only" instruction is prominent in the prompt generated for the agent. If an agent tries to be "helpful" by refactoring the whole consensus file, it might introduce formatting conflicts or lose metadata.

## Current proposal

Implement `parley consensus request-signoffs` as a subcommand that automates the collection of signoffs:
- **Target Resolution:** Resolve the target file path based on `IDEA` and the `--review` flag.
- **Participant Selection:** If `--participants` is omitted, default to all participants missing a terminal signoff, in the order defined by `00-prompt.md`.
- **Safety Gate:** Require `--yes` if any selected agent uses a non-local/external backend.
- **Dry-run:** Provide a summary of the target file, participants to be invoked, their backend types, and the expected order.
- **Execution:** Sequentially invoke each agent using the existing `runner` infrastructure.
- **Verification:** After each agent exits, re-parse the consensus file. Verify exactly one new canonical signoff for that agent was added and no `❌ BLOCK` was introduced.
- **Error Handling:** Stop immediately on process failure, missing signoff, malformed/duplicate blocks, or `❌ BLOCK`. Report partial progress and exit non-zero.
