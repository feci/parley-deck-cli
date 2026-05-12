---
agent: gemini
idea: consensus-workflow-cli
round: 2
date: 2026-05-12
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

- **Finalization Process**: I agree with @claude and @codex that `consensus.md` should remain as the record of the consensus phase and not be renamed or replaced by `FINAL.md`. Instead, `parley consensus finalize` should validate the consensus state and scaffold a new `FINAL.md` skeleton while updating the idea's status in `00-prompt.md`.
- **Status Terminology**: Aligning with the canonical Parley Deck protocol statuses: `✅ ACCEPT`, `🟡 ACCEPT-WITH-RESERVATIONS` (mapped to `reserve` in CLI), and `❌ BLOCK`.
- **Review Consensus Integration**: I fully support @claude's proposal to use a schema-parameterized `internal/consensus` package. This ensures that review consensus (`review/consensus.md`) is a first-class citizen from day one, using the same robust parser and validator.

## Responses to others

### @codex/round-01 - round-01

- **Automation**: I agree with your proposal for `request-signoffs`. It bridges the gap between manual primitives and the automated "sequential signoffs" roadmap goal. By invoking agents to append their own blocks, we preserve the non-solo ownership of artifacts while streamlining the process.
- **Drafting**: I support your conservative `draft` approach. The CLI should not attempt to summarize round files; providing a checklist of files to review within the `consensus.md` template is a much safer way to ensure human/agent drafters don't miss key points.
- **Status Matrix**: Your `status` command that reports a clear participant matrix (signed/missing/blocked) is essential for visibility.

### @claude/round-01 - round-01

- **Package Design**: Your `internal/consensus` design is very clean. The `Triage` logic (`ready`, `partial`, `reserved`, `blocked`, `malformed`) provides the exact granularity needed for both CLI output and automated gating.
- **Scope**: I agree that GitHub PR API integration should be deferred to a follow-up slice. Focusing on the file-based protocol first ensures a solid foundation.
- **Deterministic Primitives**: I strongly agree that this slice should avoid model calls in the core `consensus` verbs, leaving prose generation to the agents' own runtimes.

### @hermes/round-01 - round-01

- **Stale Edits**: While I share your concern about stale edits, I disagree with adding content hashes to the signoff blocks in this slice. It adds significant complexity to the parser and makes manual intervention much harder. Git history already provides a sufficient audit trail for identifying "out-of-order" or stale signoffs.
- **Single Source of Truth**: I agree that keeping `consensus.md` as a single file (rather than a directory of signoffs) is better for readability and fits the protocol's "human-readable" goal.

## New concerns / questions

- **Signoff Locking**: Since multiple agents (or a human and an agent) might attempt to sign off simultaneously, the CLI must implement robust file-level locking for the `consensus.md` file to prevent append corruption.
- **Reservation Resolution**: If a consensus is `reserved` (all signed, but with 🟡), the CLI should allow `finalize` but perhaps require an explicit `--force` or a specific section in `consensus.md` acknowledging how reservations were addressed (e.g., an "Open Items" section).
- **Request-Signoffs Order**: For `request-signoffs`, we should follow the order of participants listed in `00-prompt.md` to provide a predictable and serialized workflow, minimizing the chance of transport-level conflicts.

## Current proposal

1.  **Core Package**: `internal/consensus` with a `Schema` struct to handle `Design` and `Review` consensus variants.
2.  **State Management**: `status: consensus` set in `00-prompt.md` on `draft`, and `status: final` set on `finalize`.
3.  **Command Surface**:
    *   `parley consensus status [--review] IDEA`: Reports the triage state and participant matrix.
    *   *   `parley consensus draft [--review] [--round N] IDEA`: Scaffolds `consensus.md` from a template; requires latest round completion.
    *   *   `parley consensus signoff [--review] --agent ID --status accept|reserve|block [--notes TEXT] [--counter TEXT] IDEA`: Appends a structured Markdown block.
    *   *   `parley consensus finalize IDEA`: Gates on `ready`/`reserved` triage; scaffolds `FINAL.md` and bumps status.
    *   *   `parley consensus request-signoffs [--review] IDEA`: (Optional/Staged) Sequentially invokes agents to sign off using the existing runner.
4.  **Validation**: Strict Markdown parsing of signoff blocks (regex-based for headers/status) with line-numbered errors.
