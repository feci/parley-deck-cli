---
idea: consensus-request-signoffs
status: final
author: codex
consensus-date: 2026-05-13
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

### Goal

Add a focused CLI command that requests missing design or review consensus signoffs from configured headless agents while preserving Parley Deck ownership rules: the command facilitates invocation and validation, but each agent appends its own canonical signoff.

### Scope

- Command:

  ```text
  parley consensus request-signoffs [--dir DIR] [--review] [--participants IDS] [--yes] [--dry-run] IDEA
  ```

- Design consensus target: `parley-deck/ideas/<IDEA>/consensus.md`.
- Review consensus target with `--review`: `parley-deck/ideas/<IDEA>/review/consensus.md`.
- Participants are loaded from the idea's `00-prompt.md` frontmatter.
- Signoff state and validation use `internal/consensus`.
- Agent launch behavior reuses existing runner configuration and command construction where practical.

### Implementation details

- Target resolution:
  - fail clearly if the idea, target consensus file, or participant frontmatter is missing;
  - support `--review` by switching only the target consensus path.
- Participant resolution:
  - if `--participants` is provided, treat it as the exact ordered target set;
  - fail on unknown participants or participants that already have a signoff;
  - if `--participants` is omitted, select all currently missing signoffs in `00-prompt.md` participant order.
- Preflight:
  - validate every selected participant has a configured/installed runner entry before launching any agent;
  - classify each selected backend as local or non-local/hosted using existing runtime metadata;
  - require `--yes` for any selected non-local/hosted backend unless `--dry-run` is set.
- Dry run:
  - print target file, current consensus status, selected participants, backend types, launch order, command/runtime preview, and whether `--yes` would be required;
  - perform no invocation and write no files.
- Invocation:
  - invoke agents sequentially;
  - stream child stdout/stderr to the parent process;
  - pass absolute paths in the prompt for the consensus file, `00-prompt.md`, the relevant discussion artifacts, and `COOPERATION.md`;
  - instruct the agent to append exactly one signoff block for its own ID and to edit no other file or line.
- Post-invocation validation:
  - re-parse the consensus file after each agent before launching the next;
  - verify exactly one new canonical signoff exists for the expected participant;
  - reject duplicate or malformed signoff blocks;
  - reject changes to other participants' signoff blocks;
  - stop on `BLOCK`;
  - stop if the child process exits non-zero, even if a valid signoff was appended.
- Failure behavior:
  - preserve partial progress on disk;
  - print which signoffs succeeded before failure;
  - make rerun behavior clear: rerunning without `--participants` selects remaining missing signoffs.

### Tests

- Command wiring and help text include `request-signoffs`.
- Design vs review target path resolution.
- Participant selection from `--participants`.
- Default missing-signoff selection in `00-prompt.md` order.
- Unknown participant rejection.
- Already-signed participant rejection.
- Missing or unconfigured runner rejection before any invocation.
- Hosted/non-local `--yes` gate.
- `--dry-run` output and no side effects.
- Happy-path multi-agent run using fake CLIs.
- Invocation failure stops the loop and preserves earlier signoffs.
- Missing signoff after successful child exit is an error.
- Duplicate or malformed signoff is an error.
- `BLOCK` signoff stops the loop.
- Valid signoff plus non-zero child exit still returns failure.
- Review consensus mode uses `review/consensus.md`.

### Non-goals

- Full autonomous phase orchestration.
- Native GitHub review API mirroring.
- GitLab MR automation.
- Cross-process locking.
- Generated consensus prose or recommendation text.
- Durable run-log directory unless existing runner plumbing provides it cheaply.
- Release packaging.

### Verification

- `go test ./...`
- Manual smoke checks with fake local agent commands:
  - dry-run for design consensus;
  - dry-run for review consensus;
  - successful sequential append with two fake participants;
  - failure after first fake participant to confirm partial progress and resumability.
- Confirm `parley consensus status <IDEA>` reports `ready` after all requested signoffs accept.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
