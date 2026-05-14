---
idea: interactive-agent-mode
status: final
author: codex
consensus-date: 2026-05-14
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

`parley-deck-cli` will support explicit per-agent launch modes so an agent can be run as headless automation or handled through a user-driven interactive/manual handoff.

Launch modes:

- `headless`: current behavior. `parley` invokes the configured command/args, waits for completion, and validates the expected artifact.
- `interactive`: `parley` writes a handoff packet, shows the target artifact and command information, optionally launches a real user-driven terminal command, waits/polls for the artifact, then validates it.
- `manual`: `parley` writes the handoff packet and exits with next steps. A later resume/validate command checks the artifact.

Configuration:

- Add `launch_mode = "headless" | "interactive" | "manual"`, defaulting to `headless`.
- Keep existing headless command fields intact.
- Add independently configurable interactive command fields matching the current `agents.Spec` style:
  - `interactive_mode`: human-readable descriptor.
  - `interactive_command`: optional command override; defaults to the first resolved agent command.
  - `interactive_args`: argv for interactive launch.
  - `interactive_prompt_mode`: `none | file | arg`; `stdin` is forbidden for interactive mode.
  - `interactive_invoke`: `print-only | spawn-tty`; default `print-only`.
  - `interactive_timeout_ms`, `interactive_poll_ms`, `interactive_notes`.

Interactive safety rules:

- `spawn-tty` attaches the configured command to the user's terminal.
- `spawn-tty` must not pipe the task prompt through stdin, pass it through a provider headless flag, allocate or manage a pseudo-terminal, scrape output, or drive the process programmatically.
- Core CLI text stays provider-agnostic: it describes headless mode as programmatic execution and interactive mode as user-driven handoff. Provider/account-specific notes may be configured locally.
- Interactive handoffs run sequentially within a run. Headless agents may still run in parallel.

Validation and resume:

- Add a shared validation path used by headless, interactive, and manual/resume flows.
- Round artifact validation checks frontmatter identity and required sections.
- Consensus/review signoff validation reuses the append-only signoff validator and canonical status rules.
- Interactive/manual completion must not rely on file existence alone.
- `parley resume <run-id>` validates pending artifacts and records completion events. It may rerun selected headless agents that never started.

Visibility:

- `agents list`, runtime matrix output, dry-runs, launch summaries, and events show the resolved launch mode and relevant command shape.
- Handoff packets are written under `parley-deck/runs/<run-id>/agents/<agent>/` when possible, otherwise under a timestamped run directory.

Implementation order:

1. Schema/config plumbing, handoff packet writer, shared validation, resume basics, mode visibility, and consensus signoff support.
2. Mixed-mode round execution using the same primitives.
3. TUI polish and complete `spawn-tty` signal handling if not finished in the first two slices.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
- Design PR: https://github.com/feci/parley-deck-cli/pull/13
