---
idea: track-aware-driver
author: claude-1
created: 2026-07-03
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1]
status: proposed
---

## Problem / idea

Follow-up ratified by `meta-protocol-change-devx-speed` (2026-07-03). That idea added the
`track: fast | standard | deliberation` protocol text (§4.0) but deliberately deferred the
**deterministic CLI/driver enforcement** of tracks. This idea implements it:

- A **classifier** (command/library) that maps objective inputs (files touched, LOC,
  reversibility, security/data surface, `strict_gate`/`auto_implement`, pipeline/action,
  API/schema break) to a `track`, honoring the §4.0 fail-closed ordering (deliberation-first,
  then fast, else standard; on doubt → stricter).
- `parley init` / `parley run` **templating** of `track:` into `00-prompt.md` and per-track
  **timeout seeding** (5 / 15 / 30 min) into `agents.toml`.
- **Driver auto-advance + validation gates**: per-track quorum/reviewer counts, cross-review
  caps, collapsed consensus/FINAL for fast, fix-up caps; and hard **rejection** of any run
  that would drop a §4.0 all-track invariant (non-solo, refutation-default).

## Constraints
- Enforce, never weaken, the §4.0 invariants and the §14 human brake.
- `deliberation` behavior must remain exactly today's full lifecycle.

## Non-goals
- Physical protocol-doc reorganization (that is `protocol-restructure-appendices`).
