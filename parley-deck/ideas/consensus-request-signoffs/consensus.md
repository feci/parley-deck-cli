---
idea: consensus-request-signoffs
drafted-by: codex
date: 2026-05-13
---

## Agreed decisions

- Implement a focused command:

  ```text
  parley consensus request-signoffs [--dir DIR] [--review] [--participants IDS] [--yes] [--dry-run] IDEA
  ```

- Resolve the target consensus file from the idea slug:
  - default: `parley-deck/ideas/<IDEA>/consensus.md`;
  - with `--review`: `parley-deck/ideas/<IDEA>/review/consensus.md`.
- Load idea participant order from `00-prompt.md` frontmatter and current signoff state from `internal/consensus`.
- Resolve target participants as follows:
  - with `--participants`, treat the provided IDs as the ordered target set and fail on unknown or already-signed participants;
  - without `--participants`, default to all currently missing signoffs in `00-prompt.md` participant order.
- Validate every selected participant has a configured/installed runner entry before launching any agent.
- In `--dry-run`, print target file, current status summary, selected participants, backend types, launch order, command/runtime preview, and whether `--yes` would be required. Do not invoke agents or modify files.
- Without `--dry-run`, require `--yes` before launching any selected hosted or non-local backend.
- Invoke selected agents sequentially through the existing runner command construction path where practical.
- Generate a narrow signoff prompt that uses absolute paths and instructs the agent to append exactly one canonical signoff block for its own ID without editing any other file or line.
- After each invocation, re-parse the consensus file and verify:
  - exactly one new canonical signoff exists for the expected participant;
  - no duplicate or malformed signoff block was introduced;
  - no other participant signoff block changed;
  - the signoff status is not `❌ BLOCK`;
  - the child process exited successfully.
- Stop immediately on child process failure, missing signoff, malformed/duplicate signoff, changed foreign block, or `❌ BLOCK`.
- Preserve partial progress on disk. Rerunning without `--participants` should select any remaining missing signoffs.
- Keep native GitHub review API mirroring out of scope for this slice.

## Agreed trade-offs

- The command optimizes for correctness and auditability over speed: selected agents run sequentially, not in parallel.
- The safe default may include hosted agents, but the `--yes` gate prevents accidental launches.
- A child that appends a valid signoff but exits non-zero still causes the command to stop and return failure. The valid file change remains on disk for inspection/resume.
- Durable per-agent run logs are deferred. The command should stream child output and print precise diagnostics; a run-log directory can be added later if the existing runner already provides it cheaply.
- Cross-process locking is deferred. The help text should warn users not to run two signoff request commands against the same consensus file concurrently.

## Open items deferred to implementation

- Decide whether participant frontmatter parsing belongs in `internal/consensus` or a small helper near the command.
- Confirm the exact runtime metadata field used to classify local vs hosted/non-local backends.
- Extend the consensus parser API only if current status details cannot distinguish malformed file, duplicate block, missing block, accepted reservation, and block status.
- Keep prompt construction testable without invoking real CLIs.
- Cover the agreed behavior with focused fake-CLI tests, including dry-run, hosted gate, review path, participant selection, already-signed rejection, missing runner, malformed/duplicate append, block stop, non-zero-after-append, and a happy-path multi-agent run.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-13
Status: ✅ ACCEPT
Notes: Ready for finalization after the remaining participants append their own signoffs.
