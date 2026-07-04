---
idea: named-roster-presets
author: user
created: 2026-07-04
track: standard
participants: [claude-1, codex-1, hermes-1]
roles:
  claude-1: facilitation + config-layering coherence
  codex-1: Go config internals (LoadAgentSpecs layering)
  hermes-1: UX of preset selection + failure modes
status: round-01
---

## Problem / idea

Inspired by Hermes Agent v0.18.0, where Mixture-of-Agents presets became first-class
selectable "models": a named ensemble you pick like any single model.

Parley analogue: today the roster for an idea is chosen ad hoc — the facilitator
hand-picks participants per idea and types them into `00-prompt.md`. There is no way
to say "use the council" or "use the fast pair" as a named, reusable selection.

Proposal to deliberate:

- **Named roster presets** in `~/.parley/agents.toml` (central) and deck
  `parley-deck/agents.toml` (override), e.g.:
  ```toml
  [rosters.council]
  participants = ["claude", "codex", "hermes", "agy"]
  [rosters.pair]
  participants = ["claude", "codex"]
  ```
- **Track-linked defaults**: optional mapping so `track: fast` defaults to
  `rosters.pair`, `deliberation` to `rosters.council`, `standard` to a middle preset —
  used only when the idea author does not list participants explicitly.
- **CLI surface**: `parley idea new <slug> --roster pair` (or equivalent existing
  entry point) expands the preset into the `participants:` list in `00-prompt.md`.
  The expanded list in `00-prompt.md` stays the single source of truth for quorum.

## Constraints

- `participants:` in `00-prompt.md` remains canonical; presets are sugar that EXPANDS
  at idea creation. Quorum rules, §2 roster, and signoff semantics unchanged.
- Config precedence must match the existing layering: built-ins → ~/.parley/agents.toml
  → deck config (deck wins).
- A preset referencing an agent not in the §2 roster or not available must fail loudly
  at expansion time, not silently shrink quorum.
- Track policy (§4.0 reviewer counts etc.) still applies AFTER expansion — a preset
  cannot bypass track minimums (e.g. non-solo).
- CLI-only change preferred; protocol text change only if genuinely needed (a §2 or
  §4.0 sentence at most — keep it additive).

## Non-goals

- No per-preset model overrides in v1 (agents keep their own model config).
- No dynamic/conditional rosters (no "pick agents by load/cost").
- No change to the deck-bootstrap roster confirmation flow.
