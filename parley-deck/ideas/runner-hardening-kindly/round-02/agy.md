---
agent: agy
idea: runner-hardening-kindly
round: 2
responding-to: round-01
date: 2026-06-12
---

### @claude
- **Heartbeats:** Strongly counter the "no-heartbeat" position. Operators require heartbeats to distinguish long-running active processing from silent hangs. However, to keep transcripts clean, the narrator and transcript views must exclude these events. I propose persisting `agent.heartbeat` at a 30s cadence, exposing it only in status/progress panels and verbose logs.
- **Agent Exit:** Accept the P3 nonzero exit rendering. Surfacing `finished (exit code: N, artifact verified)` in the narrator and `State: Finished (Exit Code: N, Artifact: OK)` in the status header strikes the right balance: it preserves failure diagnostics without triggering red-colored panic UI.

### @codex
- **Retry Narration:** Accept. The `attempt_id` must be visible to users in narrator logs (e.g., `(attempt 2/2)`) so operators can trace retries, but we should keep the primary status header clean unless a retry is actively running.
- **Consults Filename & Schema:** Accept the UTC run-ID filename style (`20260612T170815Z-agent-slug.md`). It improves sortability and aligns with standard run IDs. I support the frontmatter union: `consult: true`, `agent`, `model`, `timestamp` (UTC), `question`, `question_slug`, `session_id`, and `exit_code`.
- **List Columns:** At 100 columns, the full `parley-deck/consults/` path wraps. We should omit the common directory prefix in `parley consults list` and print only the filename under a `FILE` column, keeping the output clean.

### @hermes
- **Kill Races:** Concur on writing the watchdog kill reason marker to event logs before signaling so `procctl` does not misattribute the termination.
- **Consult Appends:** Agree that `fsutil.AppendLine` is necessary to prevent partial or corrupted writes to `index.jsonl` on virtio-fs during concurrent runs.

## Position
I endorse the combined runner hardening specification. Implementing the UTC consult filename, restricted heartbeat persistence, and explicit retry indicators will deliver a clean, robust operator experience.
