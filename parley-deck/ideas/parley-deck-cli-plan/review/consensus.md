---
idea: parley-deck-cli-plan
drafted-by: codex
date: 2026-05-10
status: consensus
---

## Review outcome

Fix-up cycle 3 is accepted for this implementation slice.

## Agreed fixes completed

- `parley run` no longer launches every installed agent silently:
  - `--participants` selects a subset;
  - invalid participant IDs fail before idea/run creation;
  - default HITL mode requires confirmation;
  - `--yes` and `--auto` explicitly allow non-interactive launch.
- Runner event-store errors are surfaced through command results instead of being discarded.
- Setup-time agent failures emit `agent.failed` events and populate `CompletedAt`/`Duration`.
- Concurrent JSONL appends are serialized with a process-wide mutex.
- Gemini isolated-home setup now uses `os.UserHomeDir()` and fails fast when credentials are unavailable.
- Runner tests now cover successful fake-agent execution, existing-artifact skip, and child-process failure.

## Deferred follow-ups

- Live TUI streaming while agents are running.
- Token accounting from agent-specific telemetry.
- HITL question capture and `parley answer`.
- Durable resume.
- Participant selection UX beyond comma-separated flags.
- Adapter modularization and model/effort configuration.
- A future protocol/design update for Hermes support in the finalized adapter list.

## Dismissed findings

None.

## Signoffs

### Signoff: codex — 2026-05-10
Status: ACCEPT
Notes: Cycle 3 closes the blocking safety and durability findings. Remaining review items are deferred follow-ups and should not block the next implementation slice.

### Signoff: claude — 2026-05-10
Status: ACCEPT
Notes: Round-03 review found no CRITICAL or MAJOR issues; safety, durability, and test-coverage fixes are all in place.
