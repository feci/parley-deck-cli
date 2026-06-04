---
idea: tui-interactivity-overhaul
review-cycle: 1
drafted-by: claude
date: 2026-06-04
reviewed-commit: 09b28e6
---

Synthesis of the Phase 6 review-round-01 findings from codex and hermes (see
review/round-01/{codex,hermes}.md). Agreed fixes are applied in Phase 8 fix-up
cycle 1 on branch feat/tui-interactivity-overhaul.

## Agreed fixes

### AF1 — Composer must not over-promise execution (hermes CRITICAL)
The composer/help strings ("runs on next attempt") imply the queued steer
executes, but slice 4 only records it (no consumer; slice 5 deferred). This
violates the consensus "label delivery state explicitly — no over-promising"
rule. Fix: reword the composer body, the queued-confirmation, and the help line
to state the steer is *recorded/queued* and that execution lands in a later
slice (not yet auto-run).

### AF2 — Focus buffer must be byte-bounded, not only line-bounded (codex MAJOR)
`refreshFocus` caps only by `maxFocusLines`; `readAppendedLines` does
`io.ReadAll` from the offset, so a >4 MiB single line or burst exceeds the
mandated 4 MiB cap and reallocates each tick. Fix: enforce a byte budget too —
bound the per-tick read and evict oldest lines until both `len<=maxFocusLines`
and retained bytes `<=maxFocusBytes`. Tests for a >4 MiB line and a >4 MiB burst.

### AF3 — Steering IDs must not collide across CLI/TUI (codex MAJOR)
`steer-NNNN` is a load-count-write across two separate appends; concurrent
`parley steer` / TUI submitters can mint duplicate IDs and corrupt the
projection. Fix: make IDs collision-resistant (monotonic prefix + random
suffix) so the by-id projection never merges two distinct steers; concurrency
test. True cross-process *ordering*/locking stays deferred (F4) and is noted.

### AF4 — loadFocusTail must not advance past a partial final line (codex MINOR)
`loadFocusTail` returns `offset=size` even when the tail lacks a trailing `\n`,
so a later completion shows only the appended suffix as a new line. Fix: advance
the offset only past the last newline (like `readAppendedLines`); re-read the
partial line when complete. Test: initial `partial`, append ` line\n`.

### AF5 — Untagged agent events inherit the current segment (codex MINOR)
FINAL requires "events without explicit ids inherit the latest known segment."
Fix: track `currentSegment` while projecting and apply it in `applyAgentEvent`
when the event has no `segment_id`. Test: untagged `agent.started` after a
`run.segment_started`.

### AF6 — Tag ACP agent events with segment_id (hermes MAJOR / closes a deviation)
Close the documented ACP gap: tag `agent.started`/`agent.finished` in
`internal/runner/acp.go` with `opts.SegmentID`, so the audit trail is uniform
and segment scoping does not rely solely on the reset barrier for ACP launches.

### AF7 — Integration test of the segment fix through the real load path (hermes MAJOR)
Add a `runstate.LoadRunAt` test that writes real `run.segment_started` +
`agent.*` events to `events.jsonl` and asserts the targeted agent has no stale
terminal badge after a continue segment (not only the direct `ProjectEvents`
unit tests).

### AF8 — `parley steer` arg validation (hermes MINOR)
Empty TEXT and unknown `--agent` should produce a clear error/warning, not a
confusing failure or silent queue for a non-participant. Fix: explicit empty-text
usage error; warn (non-fatal) when the agent is not in the run's participants.

### AF9 — Overview preview should fill available rows, not a hard 6 (codex NIT)
`refreshLogPreview` still uses `tailLogFile(..., 6)`. Fix: derive the preview
line budget from the layout height (the right pane is clipped anyway), getting
closer to FINAL's "as many lines as fit".

## Deferred follow-ups (not this cycle)

- F4 (carried) — cross-process *atomic* append / file-locking for `events.jsonl`
  when CLI + TUI write concurrently. AF3 removes the ID-collision corruption;
  strict global ordering/locking remains deferred and is documented.
- Slice 5 (carried) — actually executing the queued steer (`new_attempt`
  launch), live `acp_prompt`/`native_resume`, opt-in ACP thoughts, per-stream
  (stdout/stderr) scrollback tabs in the focus view.

## Dismissed findings

- hermes NIT "mode enum comment drift": the comment intentionally records what
  the enum replaced (migration context); reworded for clarity but not a defect.

## Signoffs

### claude — ACCEPT (2026-06-04)
Accept the synthesis. AF1–AF9 are correct and I will apply them in fix-up cycle 1;
the two carried deferrals (F4 cross-process locking, slice 5 execution) are the
right scope boundary. No blockers.

### codex — ACCEPT (2026-06-04)
Confirmed across review rounds 01, 02, and 03 that AF1-AF9 are resolved, including the cycle-2 `capFocusLines` fix, and zero agreed fixes remain. The carried deferrals F4 and slice 5 remain acceptable scope boundaries.

### hermes — ACCEPT (2026-06-04)
Confirmed in round-02 re-review (incl. CRITICAL AF1 composer over-promise) that AF1-AF9 are resolved and zero agreed fixes remain; honest bounded viewport post-fix-up.
