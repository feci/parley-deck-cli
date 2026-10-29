---
from: user
to: codex-1
idea: meta-protocol-change-credit-exhaustion-failover
phase: round-02
blocking: no
date: 2026-10-29
---

## Owner answer

This answers `codex-1-to-user_meta-protocol-change-credit-exhaustion-failover_decision.md`.
The owner selected "Oboje (Recommended)" (both) at about 20:10 CET on 2026-10-29. The relay recorded it.

1. **Native CLI fallback is authorized** for this idea (your option 1). Use it wherever a recorded driver
   gate prevents a required participant action. Keep the brief's models, timeouts, retry limits,
   independent roles, private logs and audit receipts. Do not claim driver accounting or automatic
   membership authority for fallback steps.
2. **The relay repaired the worktree registrations** (owner-approved). The workspace moved from
   `/Volumes/My Shared Files/AI_WORKSPACE/` to `/Users/tomasfecko/AI_WORKSPACE/`. The relay ran
   `git worktree repair` on 23 parley-deck-cli registrations and 1 parley-deck-skill registration. Nothing
   was pruned or declared. The backup is in `runs-handoff/worktree-repair-backup-20261029/`.

## Result of the repair (relay check, 20:15 CET)

- `parley budget worktree inspect --dir .` now succeeds.
- `parley budget migrate inspect --kind cross-review` no longer reports the unavailable worktree. It now
  reports a different gate: `historical run identity is missing or conflicting`.

Try the driver once at each phase. If a gate refuses, use the authorized fallback and record the gap with
its exact message in organizer-notes.md. Do not declare runs or worktrees. Do not apply legacy accounting.
Continue the run to release without a new stop for this gate.
