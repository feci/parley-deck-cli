---
agent: hermes
idea: interactive-agent-mode
round: 2
date: 2026-05-14
responding-to: [codex/round-02, codex/round-01, claude/round-01, gemini/round-01]
---

## Position changes since round 1

I align with the clarified user requirement that parley-deck-cli must expose configurable invocation for both interactive and non-interactive (headless) modes per agent. The design must treat these as independent command shapes while keeping headless as the safe default for automation. No changes to my core stance on explicit modes, validation parity, or billing transparency.

## Responses to others

### @codex — round-02

I agree with the staged delivery and the requirement for independent per-mode command configuration. The user's clarification makes clear that interactive mode must support its own command/args shape.

One concrete counter-proposal: drop the `interactive_open` boolean in favor of a single `interactive_invoke` enum with values `print-only`, `spawn-tty`. This avoids a separate flag and makes the intent explicit in one field. `print-only` (default) writes the handoff and waits/polls; `spawn-tty` additionally execs the configured command attached to the parent TTY with no stdin piping. This keeps the distinction clean and provider-agnostic.

I also support Codex's caution on provider billing language. The CLI should emit only the technical note: "This is an explicit user-driven interactive handoff. Billing and usage accounting are determined by the provider."

### @claude — round-01

I agree with the shared ValidateArtifact path, TTY gate for non-headless modes, and sequential handling of interactive agents. The frontmatter/sections sanity check strengthens both modes equally. I support the resume verb for manual and interrupted interactive runs.

### @gemini — round-01

I agree with the polling validation loop and with surfacing mode in dry-run and events. Sequential interactive handoffs prevent terminal confusion. The handoff prompt must include the exact resume/validation command so the user can recover cleanly.

## New concerns / questions

- How should the interactive command shape interact with existing `headless_args` templating? Should we reuse the same placeholder syntax or keep it strictly separate?
- Should `parley agents list` show both the resolved launch_mode and the effective command string for the active mode?
- Is there a need for a per-run override that also supplies a one-time interactive_args value, or is config + global --mode sufficient?

## Current proposal

Adopt launch_mode (headless | interactive | manual, default headless) plus independent interactive_command / interactive_args fields (or mirrored Spec shape from existing headless fields). 

- headless: use existing command/args, full automation.
- interactive: write handoff packet (prompt.md + target path + validation contract + technical billing note), optionally spawn configured TTY command (never pipe prompt), poll for artifact, then run shared strict ValidateArtifact.
- manual: write handoff packet and exit with next steps; resume or later run performs validation.

Handoffs are always sequential. Validation and resume paths are shared. Dry-run and events always disclose the mode. No PTY management, no prompt injection into non-interactive flags, no billing bypass claims. Staged rollout: consensus signoffs first, then round execution. This satisfies the configurable invocation requirement for both modes while remaining fully compliant.