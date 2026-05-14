---
agent: gemini
idea: interactive-agent-mode
review-round: 1
date: 2026-05-14
reviewed-commit: 348b693
---

## Summary
The implementation follows the `FINAL.md` design and `IMPLEMENTATION.md` staged plan. It successfully adds per-agent `launch_mode` configuration, handoff packet generation, and support for interactive and manual consensus signoffs. The addition of stricter Round 1 artifact validation improves protocol compliance.

## Findings
### [MAJOR] Signal handling for `spawn-tty`
As noted in the `IMPLEMENTATION.md`, process-group signal handling for `spawn-tty` is deferred. In an interactive environment, if the user interrupts the `parley` command (e.g., via Ctrl-C), the spawned agent process might not receive the interrupt signal correctly or could become orphaned. This can lead to inconsistent state or leaked processes if the interactive agent is not designed to handle its parent's termination independently.

### [MINOR] Manual mode exit code semantics
The `parley consensus request-signoffs` command exits with code 0 even when signoffs are in `manual` (pending) mode. For automation and CI pipelines, exit code 0 usually indicates that the requested operation has been fully completed. Since `manual` mode leaves the task in a pending state, a distinct non-zero exit code (or a specific "pending" status code) would be more appropriate to signal that human intervention is still required before the run can be considered successful.

### [MINOR] `ValidateArtifact` Round 2+ compatibility
The `ValidateArtifact` function in `internal/runner/validation.go` currently enforces Round 1 specific sections (`## Summary`, `## Proposed approach`, `## Concerns / open questions`, `## Risks`) for any artifact it validates, regardless of the `round` parameter. While later rounds are currently deferred, this implementation will cause validation failures once Round 2+ is introduced, as those rounds have different required sections defined in `FINAL.md`. The validation logic should eventually be branched or parameterized by the round number.

### [NIT] `ReadFrontmatter` parser robustness
The `ReadFrontmatter` parser in `internal/protocol/workspace.go` is a simple line-by-line splitter. It does not support more complex YAML/TOML features like multi-line strings or nested structures that agents might use in their frontmatter. While it serves the current requirements, a more robust parser (like the one used for `agents.toml`) would be more resilient as the protocol evolves.

### [NIT] Hardcoded poll interval for interactive mode
The poll interval for checking signoffs in interactive mode is hardcoded to 2 seconds (`DefaultInteractivePollMS`). While this is configurable per agent in the TOML, the default might feel slightly slow for users expecting immediate feedback after they manually append a signoff.

## Open questions
- Is there a plan to unify the exit code behavior for all "pending" states (e.g., round handoffs vs. consensus handoffs) in a future slice?
- Will the `resume` command eventually support re-running headless agents that failed or were skipped in a previous attempt?
