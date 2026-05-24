---
idea: continuous-run-tui
implemented-by: codex
date-started: 2026-05-23
status: slice-1-implemented-awaiting-review
---

# Implementation

## Scope

Implement the first slice from `FINAL.md`:

- durable run action planner;
- conservative `run.json` continuation snapshot fields;
- `parley continue` planning command;
- planner data available to TUI/workspace run summaries.

## Protocol exception

Implementation proceeds under the user-authorized exception recorded in:

- `parley-deck/inbox/codex-to-all_continuous-run-tui_round-01-agent-blockers.md`
- `parley-deck/inbox/user-to-all_continuous-run-tui_solo-exception.md`

## Log

- 2026-05-23: Started slice 1 implementation.
- 2026-05-23: Added `internal/runplan` with planner-derived next actions for HITL answers, missing/failed participant retry, consensus drafting, signoff requests, finalization, and inspect fallbacks.
- 2026-05-23: Extended run summaries with `next_actions` and run manifests with conservative continuation snapshot fields: phase, idea status, current round, active steps, last action, and cached next actions.
- 2026-05-23: Added `parley continue [--dir DIR] [--json] RUN_OR_IDEA` as the read-only planning entrypoint for continuing an existing run or idea.
- 2026-05-23: Surfaced planner actions in the workspace TUI actions panel while preserving the existing questions panel.
- 2026-05-23: Updated README and CLI reference.

## Verification

- `env GOCACHE=/private/tmp/parley-go-cache go test ./...`
- `env GOCACHE=/private/tmp/parley-go-cache go run ./cmd/parley help`
- `env GOCACHE=/private/tmp/parley-go-cache go run ./cmd/parley continue --dir . 20260510T201528.433687000Z`

## Review status

Peer-agent review remains blocked by the same local agent auth/model/sandbox issues recorded for the design phase. This slice is implemented and verified locally, but not cleanly multi-agent reviewed.
