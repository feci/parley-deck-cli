---
agent: hermes
idea: consensus-request-signoffs
round: 1
date: 2026-05-13
---

## Summary
Implement a focused CLI command `parley consensus request-signoffs` that resolves missing consensus or review signoffs for a given idea by sequentially invoking configured headless agents. It reuses the existing `internal/consensus` parser and status logic, operates strictly on files under `parley-deck/`, supports both design and review consensus via the `--review` flag, and enforces safety through explicit participant lists, `--yes` confirmation, and `--dry-run` inspection. Sequential execution with immediate failure on missing, duplicate, or blocking signoffs keeps the slice minimal while providing a practical automation step after the deterministic parser has landed.

## Proposed approach
- Resolve the target consensus file (`consensus.md` or `review/consensus.md`) for the supplied IDEA using existing path conventions and the `internal/consensus` module.
- Determine required participants from the `--participants IDS` flag or fall back to a safe default of agents still missing signoffs (never auto-include all).
- For `--dry-run`, emit a concise summary of intended participants, execution order, launch commands, and target file without any invocations.
- When not in dry-run mode, require `--yes` (or interactive confirmation) before launching any hosted or external agent; local agents may proceed with lighter gating.
- Invoke agents sequentially: construct the appropriate runtime command for each, capture its append of a canonical signoff block, then immediately re-parse and verify the result.
- Halt on any invocation failure, absent signoff, duplicate/malformed block, or detection of a `❌ BLOCK` status.
- Keep all logic file-system centric; no GitHub review API, no cross-process locking, and no generated prose in this slice.
- Reuse existing config loading, command builders, and agent profiles to minimize new surface area.

## Concerns / open questions
- How to surface per-agent timeouts or retry behavior without expanding scope beyond the current constraints?
- Should the default missing-signoff list be computed before or after the `--participants` filter, and how do we communicate that choice to the user?
- Is there a minimal safe way to handle an agent that writes a signoff but then exits with a non-zero code?
- Will the sequential loop need any form of progress reporting, or is the existing parser output sufficient for observability?

## Risks
- Sequential execution could leave the consensus file in a partially-signed state if an external agent hangs or crashes mid-invocation; mitigation is immediate verification after each call and explicit stop on error.
- Reliance on `--yes` / `--dry-run` is only as strong as user discipline; a future auto mode is explicitly out of scope for this slice.
- Any drift between the runtime command construction and the actual agent entrypoints could produce malformed signoffs that the parser then rejects, requiring tight reuse of existing builders.