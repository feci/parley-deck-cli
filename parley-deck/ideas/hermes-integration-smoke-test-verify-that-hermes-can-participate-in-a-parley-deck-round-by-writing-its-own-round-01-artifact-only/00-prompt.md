---
idea: hermes-integration-smoke-test-verify-that-hermes-can-participate-in-a-parley-deck-round-by-writing-its-own-round-01-artifact-only
author: user
created: 2026-05-10
participants: [hermes]
status: abandoned
---

## Problem / idea

Hermes integration smoke test: verify that Hermes can participate in a Parley Deck round by writing its own round-01 artifact only.

## Constraints

- Local-directory transport for this initial run.
- Human-in-the-loop mode by default unless the run was started with auto mode.

## Non-goals

- Do not make unrelated repository changes.

## Outcome

Abandoned after the initial smoke run showed that Hermes needed an isolated writable `HERMES_HOME` for logs in this sandbox. The follow-up smoke idea with suffix `-2` verifies the fixed adapter.
