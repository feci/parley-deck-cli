---
agent: codex
idea: meta-protocol-change-roster-headless-config
round: 1
date: 2026-05-25
---

## Summary

The protocol should replace placeholder roster rows with the four stable local
participants already used throughout the project: `codex`, `claude`, `gemini`,
and `hermes`. The roster should identify them as local CLI-backed participants
without hard-coding machine-specific executable paths.

## Proposed approach

- Update `COOPERATION.md` section 2 roster rows:
  - `codex` with workspace dir `current workspace / Codex CLI`;
  - `claude` with workspace dir `local Claude Code CLI`;
  - `gemini` with workspace dir `local Gemini CLI`;
  - `hermes` with workspace dir `local Hermes CLI`.
- Replace host handle placeholders with `n/a` for now because no GitHub handle
  mapping has been confirmed.
- Add a short note under the roster that local launch configuration may be
  stored in ignored machine-local files such as
  `parley-deck/meta/headless-agents.local.json`, and that those files do not
  change quorum or canonical protocol state.
- Do not commit the local config file itself; keep it covered by `.gitignore`.

## Concerns / open questions

- If GitHub PR transport later requires native review approvals from host
  accounts, a follow-up roster update should map agent IDs to actual handles.
- The wording must not imply that local CLI paths in one user's machine are
  valid for every clone of the repository.

## Risks

- Leaving placeholders causes future agents to treat roster membership as
  ambiguous.
- Over-specifying local paths in the protocol would make the project less
  portable.
