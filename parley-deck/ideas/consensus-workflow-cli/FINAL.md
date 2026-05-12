---
idea: consensus-workflow-cli
status: final
author: codex
consensus-date: 2026-05-12
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

### Goal

Implement the next shippable `parley-deck-cli` slice: deterministic CLI support for design and review consensus workflows.

The slice lets users move from round files to validated consensus and finalization without hand-crafting every protocol step. It must keep canonical source of truth in `parley-deck/` files and preserve the rule that each participant owns its own signoff.

### Scope

Add a focused `internal/consensus` package plus CLI commands under:

```text
parley consensus ...
```

The slice covers:

- consensus template scaffolding;
- signoff parsing and append;
- validation and triage;
- design finalization scaffolding;
- blocked-consensus reopen;
- review consensus support through `--review`;
- status integration for idea detail views.

It does not run agents automatically. Automated `request-signoffs` is the next follow-up after these deterministic primitives land.

### Consensus package

Add `internal/consensus` with APIs equivalent to:

- `Schema` for design vs review consensus:
  - design file: `parley-deck/ideas/<slug>/consensus.md`;
  - review file: `parley-deck/ideas/<slug>/review/consensus.md`.
- `Draft(...)` to write an empty consensus template.
- `Parse(...)` to parse frontmatter, sections, and signoff blocks.
- `Validate(...)` to check participants, duplicate signoffs, required notes/counter-proposals, and triage.
- `AppendSignoff(...)` to append one canonical signoff block.
- `Finalize(...)` to create `FINAL.md` for design consensus and update `00-prompt.md` only after the file write succeeds.
- `Reopen(...)` to preserve a blocked consensus file under a numbered aborted filename and move the idea back to discussion.

Triage values:

- `ready`: all participants signed ✅;
- `reserved`: all participants signed, at least one 🟡, no ❌;
- `blocked`: at least one ❌;
- `partial`: missing signoffs;
- `malformed`: parse or validation error.

Canonical status rendering:

- CLI `accept` -> `✅ ACCEPT`;
- CLI `reserve` or `reservations` -> `🟡 ACCEPT-WITH-RESERVATIONS`;
- CLI `block` -> `❌ BLOCK`.

Notes are required for 🟡 and ❌. Counter-proposal is required for ❌.

### CLI commands

Add these commands:

```text
parley consensus status [--dir DIR] [--review] [--json] IDEA
parley consensus draft [--dir DIR] [--review] [--round N] [--by AGENT] IDEA
parley consensus signoff [--dir DIR] [--review] --agent ID --status accept|reserve|reservations|block [--notes TEXT] [--counter TEXT] IDEA
parley consensus finalize [--dir DIR] [--by AGENT] IDEA
parley consensus reopen [--dir DIR] IDEA --reason TEXT
```

Command behavior:

- `status` prints triage, participant signoff table, missing signoffs, malformed blocks, and next action. `--json` may be a small unstable developer schema.
- `draft` writes a template only. It must not synthesize the consensus prose. It updates design idea status to `consensus`.
- `signoff` appends one block for the specified participant. It rejects unknown agents and duplicate signoffs unless a later explicit replace flow is added.
- `finalize` is design-only. It succeeds for `ready`; it may succeed for `reserved` only when reservations are visibly captured in open items. It creates `FINAL.md` as a separate scaffold and then sets `00-prompt.md` to `final`.
- `reopen` is valid only for `blocked`. It records the reason, preserves the old consensus, and returns the idea to discussion. It does not create the next round directory.
- `--review` switches paths and section labels to review consensus. Review finalization is not `FINAL.md`; review consensus feeds the existing fix-up cycle.

### File templates

Design consensus template:

- frontmatter with `idea`, `drafted-by`, and `date`;
- `## Agreed decisions`;
- `## Agreed trade-offs`;
- `## Open items deferred to implementation`;
- `## Signoffs` plus the append-only warning comment.

Review consensus template:

- frontmatter with `idea`, `cycle`, `drafted-by`, and `date`;
- `## Agreed fixes`;
- `## Deferred follow-ups`;
- `## Dismissed findings`;
- `## Signoffs` plus the same append-only warning comment.

`FINAL.md` scaffold:

- frontmatter with `idea`, `status: final`, `author`, `consensus-date`, and participants;
- links to `consensus.md`, `round-01/`, and later rounds;
- explicit sections for goal, scope, implementation details, tests, non-goals, and verification.

### Status and events

When `parley status --idea <slug>` sees a design or review consensus file, it should include consensus triage and participant signoff state.

If the existing event store is available in the command context, state-changing commands emit:

- `consensus.drafted`;
- `consensus.signed`;
- `consensus.finalized`;
- `consensus.reopened`.

The filesystem remains canonical if events and files disagree.

### Transport handling

The active project transport is `github-pr`.

The implementation should not submit GitHub reviews through the API in this slice. It may print the expected native review mirror:

- ✅ -> Approve;
- 🟡 -> Approve with reservation comment;
- ❌ -> Request changes with counter-proposal.

The canonical signoff block in `consensus.md` or `review/consensus.md` wins on conflict.

### Explicit non-goals

Do not implement in this slice:

- `parley consensus request-signoffs`;
- autonomous model calls;
- generated consensus summaries;
- GitHub API review submission;
- GitLab MR automation;
- TUI consensus panel;
- stale-signoff content hashes;
- strong git/GitHub identity enforcement;
- release packaging.

### Tests

Add deterministic tests for:

- design consensus draft scaffolding;
- review consensus draft scaffolding;
- participant extraction from `00-prompt.md`;
- parsing valid signoff blocks;
- rejecting duplicate signoffs;
- rejecting unknown participants;
- requiring notes for 🟡 and ❌;
- requiring counter-proposal for ❌;
- triage states: `ready`, `reserved`, `blocked`, `partial`, `malformed`;
- `signoff` CLI output and file append;
- `status` CLI plain output;
- `status --json`;
- `finalize` creating `FINAL.md` and updating `00-prompt.md` only after success;
- `finalize` refusing `partial`, `blocked`, and malformed consensus;
- `reopen` preserving blocked consensus and restoring discussion status;
- review consensus paths selected by `--review`;
- `parley status --idea` surfacing consensus triage.

### Verification

Implementation is complete when:

- `go test ./...` passes;
- manual `parley consensus status consensus-workflow-cli` reports the expected state for this idea;
- manual `parley consensus draft` works on a fixture idea;
- manual `parley consensus signoff` appends one canonical block on a fixture idea;
- manual `parley consensus finalize` creates a separate `FINAL.md` on a fixture idea;
- manual `parley consensus status --review` works against a fixture review consensus.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
