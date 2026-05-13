---
agent: hermes
idea: consensus-request-signoffs
round: 2
date: 2026-05-13
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

Round-01 proposals already converge on the core requirements: sequential execution, `--dry-run` inspection, `--yes` gating for hosted/non-local backends, reuse of `internal/consensus` and runner command construction, append-only signoff prompts, and mandatory post-invocation validation. My position has not changed; I now endorse the explicit defaults proposed in codex/round-02 as the minimal safe behavior that keeps the slice narrow while remaining usable.

## Responses to others

### @codex - round-01

I agree with the round-01 framing and the round-02 refinements. Selecting all missing signoffs in `00-prompt.md` participant order when `--participants` is omitted, combined with the strict `--yes` gate for any non-local backend, is the right balance. Requiring explicit failure on already-signed participants when the flag is supplied also prevents silent no-ops. The post-append verification contract (exactly one canonical signoff, stop on BLOCK or malformed) matches my round-01 proposal exactly.

### @claude - round-01

I agree with the detailed success criteria and the separation of selection, gating, and validation phases. The emphasis on preserving partial progress on disk and surfacing it clearly in error messages is important. I support inheriting the existing command-construction path rather than re-implementing argv logic. One minor clarification: when `--participants` is supplied, we should reject any ID that has already signed off in the current file (as codex now states) rather than silently skipping; this keeps behavior deterministic and auditable.

### @gemini - round-01

I agree that the signoff prompt must explicitly instruct "append only" and must pass absolute paths to both the consensus file and the source `00-prompt.md`. Streaming child output live while still performing the immediate re-parse check after exit is the correct observability model. The warning about concurrent manual edits is useful documentation but does not require code changes in this slice.

## New concerns / questions

- The participant order must be read deterministically from the `participants:` frontmatter array in `00-prompt.md`; if no helper exists yet, a small, tested frontmatter parser scoped to this field is acceptable.
- The consensus parser must expose a clear predicate for "exactly one well-formed signoff block for this participant" so the verification step can reuse it without duplicating logic.
- Error messages on partial runs must explicitly list which signoffs succeeded before the failure so the user can resume by re-running without `--participants`.

## Current proposal

The slice is ready for consensus. Implement:

```text
parley consensus request-signoffs [--dir DIR] [--review] [--participants IDS] [--yes] [--dry-run] IDEA
```

Behavior:
- Resolve target file: `consensus.md` or `review/consensus.md` via `--review`.
- Load idea participants from `00-prompt.md` frontmatter and current signoff state via `internal/consensus`.
- Resolve targets: `--participants` (ordered, fail on unknown or already-signed) or default to all missing signoffs in idea order.
- Validate every target has a configured runner entry.
- `--dry-run`: print target file, current status, selected participants with backend types, launch order, and whether `--yes` would be required.
- Without `--dry-run`, require `--yes` for any non-local backend.
- Invoke sequentially using the existing runner command builder.
- After each invocation: re-parse, assert exactly one new canonical signoff for that participant, reject duplicates/malformed blocks, and stop on `❌ BLOCK`.
- Halt on child failure, missing signoff, malformed/duplicate, or BLOCK. Preserve partial progress; make resumption via re-run explicit.
- Cover all selection, gating, validation, review-path, and failure-mode cases with focused tests that use fake CLIs.

This meets every constraint in 00-prompt.md and resolves the remaining design choices.