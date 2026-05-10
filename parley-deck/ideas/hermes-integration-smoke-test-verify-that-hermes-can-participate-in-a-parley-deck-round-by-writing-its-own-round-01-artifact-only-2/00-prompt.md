---
idea: hermes-integration-smoke-test-verify-that-hermes-can-participate-in-a-parley-deck-round-by-writing-its-own-round-01-artifact-only-2
author: user
created: 2026-05-10
participants: [hermes]
status: final
---

## Problem / idea

Hermes integration smoke test: verify that Hermes can participate in a Parley Deck round by writing its own round-01 artifact only.

## Constraints

- Local-directory transport for this initial run.
- Human-in-the-loop mode by default unless the run was started with auto mode.

## Non-goals

- Do not make unrelated repository changes.

## Outcome

Hermes successfully participated in a Parley Deck round and wrote its own canonical artifact at `round-01/hermes.md`.
