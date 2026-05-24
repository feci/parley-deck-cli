---
idea: continuous-run-tui
implemented-by: codex
date-started: 2026-05-23
status: complete
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
- 2026-05-24: Completed review round 01 with Claude, Gemini, and Hermes.
- 2026-05-24: Applied fix-up cycle 1 for agreed review findings:
  - moved `NextAction` into shared `internal/runaction` and aliased it from planner/manifest packages;
  - passed `current_round` from run manifests into planner input so retry/draft actions target the active round instead of hardcoded `round-01`;
  - removed hardcoded `--by codex` and hardcoded draft round from generated continuation commands.
- 2026-05-24: Completed review round 02 and review consensus with ACCEPT signoffs from codex, claude, gemini, and hermes.

## Verification

- `env GOCACHE=/private/tmp/parley-go-cache go test ./...`
- `env GOCACHE=/private/tmp/parley-go-cache go run ./cmd/parley help`
- `env GOCACHE=/private/tmp/parley-go-cache go run ./cmd/parley continue --dir . 20260510T201528.433687000Z`
- `env GOCACHE=/tmp/parley-go-cache GOMODCACHE=/tmp/parley-go-mod /tmp/go1.26.3-linux-arm64/go/bin/go test ./internal/runaction ./internal/runplan ./internal/runmanifest ./internal/runstate ./internal/app`

## Review status

Review round 02 produced no remaining agreed fixes. Review consensus is ready and all participants accepted. Implementation is complete for this slice.
