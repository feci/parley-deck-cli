---
idea: antigravity-agent-migration
author: codex
created: 2026-05-27
participants: [codex, claude, agy, hermes]
roles:
  codex: CLI implementation, release, and packaging coordination
  claude: protocol safety, migration compatibility, and documentation review
  agy: Antigravity CLI behavior, skill/install target, and Gemini replacement lens
  hermes: operational defaults, headless launch reliability, and release risk review
status: final
---

## Problem / idea

Gemini CLI is being phased out for this project and should be replaced by
Antigravity CLI, whose command is `agy`. Review and implement the migration
across `parley-deck-cli` and `parley-deck-skill`.

The work should:

- inspect all install scripts and current references to Gemini;
- update the Parley Deck skill installer/runtime targets so Antigravity is a
  first-class install target;
- update the CLI's built-in/default agent configuration so new workflows prefer
  `agy` over `gemini`;
- review `parley init`, project defaults, and `COOPERATION.md` roster defaults;
- ensure agent launches use the strongest practical model and reasoning/thinking
  settings where the target CLI exposes those controls, while keeping those
  settings configurable;
- keep fast/smart defaults where possible without hiding the selected model,
  thinking/profile, timeout, or backend disclosure;
- update documentation, tests, release notes, and packaging metadata;
- prepare release/deployment for npm, GitHub releases, Homebrew, and WinGet
  where those channels exist for the affected packages.

## Constraints

- Active transport is `github-pr`; canonical artifacts remain under
  `parley-deck/`.
- Treat `agy` as the user-approved replacement participant for this workflow
  even though the live roster still lists `gemini`.
- Do not edit historical closed idea artifacts just to remove old Gemini
  mentions. Migrate active defaults, docs, installer targets, generated
  bootstrap files, examples, tests, and release metadata.
- Do not invent model or thinking flag names. If a CLI does not expose model or
  thinking controls through help/config discovery, record `cli-default` and make
  local overrides available.
- Keep legacy Gemini support only if needed for backwards compatibility, and
  clearly label it as legacy/deprecated rather than a default participant.
- Release work must be verified before publishing.

## Non-goals

- Do not build a new orchestration service.
- Do not rewrite old audit-trail files from completed ideas.
- Do not claim WinGet deployment is complete unless manifests are updated and a
  PR or equivalent publish path is actually created.
