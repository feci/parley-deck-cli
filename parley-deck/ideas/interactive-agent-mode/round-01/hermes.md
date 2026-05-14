---
agent: hermes
idea: interactive-agent-mode
round: 1
date: 2026-05-14
---

## Summary
Design per-agent configurable launch modes (headless, interactive, manual) for parley-deck-cli. This allows explicit selection of automated vs. user-driven sessions while preserving ownership, validation, and billing transparency. Headless remains the default for automation and CI.

## Proposed approach
Extend agents.toml (and agents.local.toml) with an optional `mode` key per agent entry, defaulting to "headless". Add a top-level CLI override `--mode` (or `-m`) that applies to the current invocation.

- `headless`: Current behavior. `parley` constructs the full prompt, invokes the configured agent command non-interactively, waits for completion, then runs the standard artifact validation.
- `interactive`: `parley` writes the complete prompt and target artifact path (plus validation contract) to a handoff file or stdout, prints clear instructions including billing note, then polls the target path (with optional timeout or `--wait` flag) before performing the same validation step. User is responsible for running the agent interactively in a separate terminal.
- `manual`: `parley` only emits the prompt + target path + next-step instructions and exits. No polling or validation is performed; the user handles everything downstream.

Configuration example:
```toml
[agents.claude]
command = "claude"
mode = "interactive"   # or "headless", "manual"
```

For round artifacts and consensus signoffs, the same mode logic applies: the target file is always the canonical participant-owned artifact. Dry-run output must surface the selected mode and any billing implications. No PTY automation or secret driving of interactive sessions.

## Concerns / open questions
- How to surface billing/credit warnings reliably in interactive mode without provider-specific hooks.
- Polling strategy for interactive mode (file existence + size change vs. simple existence) and timeout defaults.
- Whether handoff prompt should be written to a dedicated `.handoff.md` alongside the target artifact or only printed.
- Impact on existing automation scripts that assume non-interactive exit codes.

## Risks
- User confusion if default remains headless but interactive is desired, leading to accidental programmatic credit usage.
- Polling loops in interactive mode could hang in CI or unattended environments if not guarded by explicit mode selection.
- Slight increase in surface area for misconfiguration that might produce empty or invalid artifacts before validation runs.