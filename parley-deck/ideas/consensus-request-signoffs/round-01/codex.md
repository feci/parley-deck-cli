---
agent: codex
idea: consensus-request-signoffs
round: 1
date: 2026-05-13
---

## Summary

The first slice should implement a deterministic facilitator command that wraps existing agent runtime invocation and consensus validation without adding new consensus semantics. The safest MVP is a sequential `request-signoffs` command with `--dry-run`, explicit confirmation, strong post-invocation verification, and no GitHub review mirroring yet.

## Proposed approach

Add a new subcommand:

```text
parley consensus request-signoffs [--dir DIR] [--review] [--participants IDS] [--yes] [--dry-run] IDEA
```

Behavior:

- Resolve the consensus summary with `internal/consensus.Status`.
- Determine target participants:
  - if `--participants` is set, use that ordered list;
  - otherwise use currently missing signoffs from the consensus summary, in `00-prompt.md` participant order.
- Validate every selected participant is configured, installed, and not already signed.
- In `--dry-run`, print:
  - target file;
  - current triage;
  - selected participant order;
  - per-agent command/runtime summary;
  - whether `--yes` would be required.
- Without `--dry-run`, require `--yes` when any selected agent has `ExternalBackend != local`.
- Invoke agents sequentially, not through the parallel round runner.
- Give each agent a narrow prompt:
  - read `COOPERATION.md`;
  - read the target consensus file;
  - append exactly one signoff block for its own agent ID;
  - do not edit any other file;
  - use canonical status values;
  - return a short summary.
- After each invocation, reload `internal/consensus.Status` and verify:
  - the expected agent now has exactly one signoff;
  - the file is not malformed;
  - if the signoff is `❌ BLOCK`, stop and return a non-zero or clearly blocked status.
- Print a compact progress table.

Implementation shape:

- Add a consensus-signoff prompt builder, probably in `internal/runner` or a small new helper package if it should not pollute round runner code.
- Reuse `runner.CommandFor` so isolated homes, headless args, prompt mode, and configured runtime behavior stay consistent.
- Do not reuse `RunRoundOne` directly because it assumes parallel execution and round artifact paths.
- Add app-level tests using fake CLIs that append signoff blocks.
- Add deterministic tests for participant selection, dry-run output, already-signed rejection, missing-agent failure, block-stop behavior, and review consensus path selection.

## Concerns / open questions

- Confirmation policy should be strict. I prefer `--yes` for any hosted backend, even when `--participants` contains only one agent, because this command can spend model tokens.
- The command should not try to decide whether an agent should accept or block. It only asks the agent to evaluate and append.
- We need a clear failure contract: if agent exits zero but does not append its block, the command must fail and point to stdout/stderr logs if captured.
- We should decide whether logs go under a new run directory, a consensus-specific run directory, or only stdout/stderr on the parent process. I prefer a run directory because it already gives us durable logs.

## Risks

- Accidentally launching hosted agents without a clear gate would violate user expectations.
- Sharing the same consensus file across subprocesses can create append races if we ever parallelize; sequential execution is a hard requirement for this slice.
- Agents may create the signoff in the wrong workspace path if prompts are ambiguous. The command should pass absolute paths and verify exact target file changes.
- If the command writes signoffs itself to compensate for agent failures, it would violate Parley Deck ownership. That must be forbidden.
