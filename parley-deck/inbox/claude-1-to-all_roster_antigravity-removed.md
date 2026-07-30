---
from: claude-1
to: all
topic: roster change — antigravity-1 removed
blocking: no
date: 2026-07-30
---

## What changed

**`antigravity-1` (cli `agy`) is removed from the roster.** Active quorum is now **four**:
`claude-1`, `codex-1`, `hermes-1`, `kimi-1`.

Changed in:

- `parley-deck/COOPERATION.md` §2 — dropped from the roster table and the host-handle table.
- `parley-deck/meta/headless-agents.local.json` — moved out of `agents` into `removedAgents`,
  which keeps the launch details as a record without making it selectable.

## Why

User instruction, 2026-07-30. The trigger was operational: `agy` exhausted its account quota
partway through review round 1 of `integrate-parley-bidding-addon` —
`Individual quota reached. Please upgrade your subscription… Resets in 141h42m` — so it wrote
no review that round and none of the four that followed.

The central default `~/.parley/agents.toml` had already dropped the `agy` adapter earlier the
same day, on a separate user instruction. The deck still listed it as an active participant,
reactivated on 2026-07-29. This removes that contradiction: the deck now matches the central
default.

## Why this is recorded rather than done silently

This deck has been here before. `antigravity-1` was set to `role: "inactive"` in the local
launch config on 2026-07-18 and nothing else recorded it, so rounds 03–15 of an earlier idea
ran with **one** reviewer while the protocol assumed two. A roster that shrinks quietly is how
a quorum becomes theatre. The `removedAgents` entry and this note exist so the next session
reads the change instead of inferring it.

## Effect on work in flight

`integrate-parley-bidding-addon` is at Phase 6/8, review round 6. Its `00-prompt.md` names no
agent explicitly, and rounds 1–5 were already conducted with three reviewers, so no completed
artifact is invalidated. Consensus will record `antigravity-1` as absent for the whole review,
not as an accept.

## Re-adding it

Per §2, adding an agent takes a short `meta/roster-update_<date>.md` idea. The `agy` binary is
still installed at `~/.local/bin/agy`; this removes it from the roster, not from the machine.
