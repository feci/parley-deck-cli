---
idea: meta-protocol-change-roster-headless-config
author: codex
created: 2026-05-25
participants: [codex, claude, gemini, hermes]
roles:
  codex: protocol owner and local workspace facilitator
  claude: protocol safety and resumability reviewer
  gemini: cross-runtime configuration reviewer
  hermes: local runtime and long-running agent reviewer
status: final
---

## Problem / idea

`parley-deck/COOPERATION.md` still contains placeholder roster rows even though
the project has repeatedly used stable agent IDs `codex`, `claude`, `gemini`,
and `hermes`. The project also now has verified local headless launch settings
for those agents in ignored machine-local configuration.

Decide whether to update the protocol roster to list the active local agents
and how to mention local headless configuration without making machine-specific
paths canonical project state.

## Constraints

- Do not treat machine-local CLI paths as portable project state.
- Keep `parley-deck/meta/headless-agents.local.json` ignored and local.
- Preserve the current `github-pr` transport.
- Keep protocol content vendor-neutral and clear that host handles are
  currently unavailable or not mapped.

## Non-goals

- Do not redesign the full protocol.
- Do not add new quorum participants beyond `codex`, `claude`, `gemini`, and
  `hermes`.
- Do not require GitHub account mapping in this slice.
