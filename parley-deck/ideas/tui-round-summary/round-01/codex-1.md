---
agent: codex-1
idea: tui-round-summary
round: 1
date: 2026-07-04
---

## Summary

Round completion is not just "all files exist" today. The driver gate in `internal/driver/driver.go` requires the round directory, one valid artifact per configured participant, extra cross-review evidence for round 2+, and a terminal `round.completed` event; if valid artifacts exist but the terminal event is missing, the driver reconstructs `round.completed`. A terminal `round.incomplete` remains authoritative and blocks completion.

The ordinary `round.completed` / `round.incomplete` event is emitted by the runner after the parallel agents finish. The driver emits `run.phase` only for cursor/phase-changing commits such as promotion to the next round, consensus drafting, finalization, implementation, review, fix-up, or completion. The TUI already consumes both families: `round.completed` triggers protocol snapshot refresh and narrator lines, while `run.phase` explains the phase transition.

## Proposed approach

Build the digest in the driver/event-state layer, not inside `renderHome`. The driver is the only place that currently has the canonical completion predicate, including artifact validation and the terminal-event reconciliation rule. If Home scans files and decides "complete" independently, it will eventually diverge from the advancement gate.

Concretely, after `advanceRound` observes `roundComplete(round) == true`, have the driver call a deterministic round-digest builder for that specific idea and round. The builder should read the already-validated round files, extract each participant's `## Summary` section with the requested fallback, compute simple mechanical signals, and include the driver's expected next action for that completed round. The result can be appended to the run event stream as a presentation event such as `round.digest` keyed by `idea` and `round`, or as a strictly idempotent enrichment path tied to `round.completed`.

The TUI should then render cached digest data, preferably by extending `ProtocolSnapshot` with recent round digests and having Home render that field. `renderHome` should stay a pure renderer over model/snapshot state: no validation, no markdown parsing, no per-render disk scan. This preserves one source of truth while still keeping the feature read-only over canonical protocol files.

## Concerns / open questions

The runner's terminal `round.completed` event is earlier and simpler than the driver's validated completion gate. A digest emitted only by the runner would miss the stricter round 2+ cross-review checks and the driver's "what happens next" action, so it should not be the source of the Home digest.

Digest emission must be idempotent. Re-running the driver against an already-complete round should not append duplicate Home history entries for the same `idea` plus `round`; the event writer should skip an existing digest or use a stable event key in the model.

There is a sequencing question: a digest can describe "round N complete" as soon as `roundComplete` succeeds, but "what happens next" is only actual after the driver action succeeds. I would record an `expected_next` before the action only if necessary, or emit the digest after the successful action with the actual driver action included.

## Risks

If the digest is built in the Home view, malformed or merely present files may be summarized as completed even when the driver would still await or block the round.

If digest extraction failures are treated as driver failures, a display feature could block protocol progress. Extraction should degrade per agent and still let advancement continue.

If the TUI stores digest history only in memory, reopened runs and returning users lose the catch-up value. The run event log is the better history source because it is already what recent runs, narrator replay, and protocol snapshots consume.
