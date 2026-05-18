---
idea: session-resume-cache-plan
review-cycle: 1
drafted-by: codex
date: 2026-05-18
reviewed-commit: a07bd53
status: user-directed-release-exception
---

## Agreed fixes

- Add `status` to the `run.json` schema before release.
- Document intentional schema field deviations in `IMPLEMENTATION.md`.
- Expand tests beyond the happy path for `sessions list`, `sessions inspect`,
  legacy runs, missing runs, and direct manifest round-trip/defaults.
- Improve CLI help and legacy-run output for the new `sessions` commands.

## Deferred follow-ups

- Malformed `sessions.json` recovery.
- Terminal run status updates.
- Native agent resume handles.
- Prompt input packs, heartbeats, stale process recovery, and TUI retry/resume
  actions.

## Dismissed findings

- None.

## Release exception

This review consensus is not a normal full-participant consensus. Claude
created `review/round-01/claude.md`, and Codex applied the agreed fixes.
Follow-up review attempts for Claude, Gemini, and Hermes were blocked by local
runtime/auth/permission issues recorded in
`parley-deck/inbox/codex-to-all_session-resume-cache-plan_release-exception.md`.

The user explicitly instructed Codex to complete all remaining release steps.
This file records that the release proceeds under direct user instruction
without fabricating missing participant signoffs.

## Signoffs

### Signoff: codex — 2026-05-18
Status: ACCEPT-WITH-RESERVATIONS
Notes: Slice 1 is implemented and tested. Release proceeds under user-directed
exception because clean non-facilitator re-review/signoff could not be obtained
from the available local CLIs.
