---
idea: agent-runtime-config
author: codex
created: 2026-05-11
participants: [codex, claude, gemini, hermes]
status: round-01
---

## Problem / idea

Design the next `parley-deck-cli` slice: explicit agent runtime configuration and user guidance for adding cooperating agents.

The CLI must stop treating agent setup as hidden environment state. It should expect and surface the effective runtime configuration needed for reliable cooperation: stable agent ID, CLI path, headless invocation, narrow workspace-write behavior, sandbox/approval behavior, model/reasoning/profile defaults, timeout policy, isolated-home requirements, and external-backend disclosure.

Codex participation has a concrete recommended profile from the user:

- `sandbox_mode = "workspace-write"`
- `approval_policy = "on-failure"`
- Git write capability should be verified directly in the target repository.
- If a Git write command fails because of sandbox or permission restrictions, the exact failed command should be retried through the runtime approval path.
- The CLI and docs should not recommend silently moving work to an alternate work directory without explicit user approval.

This slice should turn the current documentation note in `docs/agent-runtime-configuration.md` into an implementable CLI/product plan.

## Constraints

- Active Parley Deck transport is `github-pr`; canonical files remain under `parley-deck/`.
- Keep the implementation scope small enough for one design PR and one implementation PR.
- Preserve existing agent support for Codex, Claude, Gemini, and Hermes.
- Do not invent model names, reasoning levels, or CLI flags that discovery cannot prove.
- Prefer explicit `cli-default` values over guessed model or reasoning settings.
- Keep credentials, customer data, and unrelated private documents behind explicit user consent.
- Every file under `parley-deck/` and every PR comment/review summary must be in English.

## Non-goals

- Do not implement M6 consensus/auto-mode policy in this slice.
- Do not implement `parley resume` in this slice unless the design concludes a small prerequisite is unavoidable.
- Do not implement GitHub/GitLab transport automation in this slice.
- Do not implement release packaging, npm, Homebrew, or CI release workflows in this slice.
- Do not require users to commit local machine config containing private paths or credentials.
