---
agent: codex
idea: tui-protocol-visibility
round: 1
date: 2026-06-11
---

## Summary

ACCEPT direction, with two required corrections: make `ProtocolSnapshot` an
Update-time async cache, never View-time filesystem work, and pair `runToken`
with a same-run snapshot sequence/in-flight gate. The import direction is safe,
but the 7-state cursor does not currently expose enough detail for a lossless
9-step display without either duplicate I/O or a small new driver export.

## ProtocolSnapshot Producer

ACCEPT `internal/tui/protosnap.go`. `internal/tui` already imports lower layers
only (`protocol`, `runstate`, `steer`, `store`) at `internal/tui/live.go:18`.
Adding `driver`, `consensus`, and `runplan` does not create a cycle: `driver`
imports `runner`/`store` (`internal/driver/driver.go:13`) and its phase files
import `consensus` (`internal/driver/consensus.go:10`,
`internal/driver/impl.go:11`); `consensus` imports `protocol`
(`internal/consensus/consensus.go:15`); `runplan` imports
`consensus`/`protocol` (`internal/runplan/runplan.go:11`). None import `tui`.

Inputs should be explicit value copies: root, runID, runDir, idea slug, ideaDir,
participants, maxRounds, copied events/questions or copied `runstate.RunState`,
and the previous snapshot for keep-last / anti-regression decisions. Do not call
`runstate.LoadRunAt`; the live model already has events and questions, while
`LoadRunAt` would reload events, manifest, HITL, and runplan
(`internal/runstate/runstate.go:105`, `internal/runstate/runstate.go:154`,
`internal/runstate/runstate.go:165`).

Async `tea.Cmd` is correct and should mirror the existing token pattern:
`eventsMsg`, `questionsMsg`, ticks, and `doneMsg` carry `token`
(`internal/tui/live.go:312`), and stale run messages are ignored
(`internal/tui/live.go:398`, `internal/tui/live.go:414`,
`internal/tui/live.go:423`, `internal/tui/live.go:443`). Counterproposal:
`runToken` only protects run swaps; it does not stop two snapshot commands for
the same run from landing out of order. Add `protoSeq` to the message or allow
only one in-flight snapshot and coalesce dirty triggers.

ACCEPT 15s running / 60s done-or-stale cadence, but use a separate reconcile
tick. Do not attach snapshot work to the 250 ms event loop
(`internal/tui/live.go:413`, `internal/tui/live.go:2512`) or the 1s elapsed
tick (`internal/tui/live.go:2518`) except to schedule/cache.

Event-trigger allowlist should be tight:
`run.created`, `run.phase`, `run.segment_started`, `round.completed`,
`round.incomplete`, `agent.finished`, `agent.failed`, `agent.skipped`,
`agent.killed`, `agent.fixup_finished`, `agent.fixup_failed`, `run.failed`.
Exclude ACP chatter, steer events, HITL noise, and index-writing events.
`agent.stdout_fallback` need not trigger because normal completion later emits
`agent.finished` with `artifact_ok` after validation
(`internal/runner/runner.go:400`, `internal/runner/runner.go:437`).

Export needs:

- Already exported: `driver.Rebuild`, `driver.Cursor`, `driver.Phase*`
  (`internal/driver/cursor.go:25`, `internal/driver/cursor.go:87`).
- Needs export or wrapper if TUI must avoid duplicated disk probes:
  `implementationStatus`, `highestReviewRound`, `finalScaffoldReason`
  (`internal/driver/cursor.go:124`, `internal/driver/cursor.go:130`,
  `internal/driver/consensus.go:150`).
- Do not export `consensus.missingRoundArtifacts` just for display; a tiny local
  stat loop is fine (`internal/consensus/consensus.go:486`).
- `runplan.readCrossReviewRounds` is private, but `protocol.ReadFrontmatter`
  can read `cross_review_rounds` (`internal/runplan/runplan.go:253`,
  `internal/protocol/workspace.go:236`).

## Cursor To 9-Step Pipeline

ACCEPT the broad mapping, but not the "fully derivable from current Cursor"
claim. The seven cursor states are exactly `round`, `consensus`, `final`,
`impl`, `review`, `done`, `blocked` (`internal/driver/cursor.go:28`). However,
`Rebuild` has no branch that returns `PhaseBlocked`; blocked consensus comes
from `consensus.Status(...).Triage == TriageBlocked`
(`internal/driver/cursor.go:97`, `internal/consensus/consensus.go:416`).

Mapping I would ship:

- 0 kickoff: `PhaseRound` plus `IdeaStatus == kickoff`, which `Rebuild` already
  reads from `00-prompt.md` (`internal/driver/cursor.go:89`,
  `internal/driver/cursor.go:181`).
- 1 round-01: `PhaseRound && CurrentRound == 1`.
- 2 cross-review: `PhaseRound && CurrentRound >= 2`; current round comes from
  `highestRound` (`internal/driver/cursor.go:90`,
  `internal/driver/cursor.go:147`).
- 3 consensus: `PhaseConsensus`, caused by `consensus.md`
  (`internal/driver/cursor.go:111`).
- 4 final: `PhaseFinal`, caused by a valid `FINAL.md`, with scaffold caveats
  (`internal/driver/cursor.go:106`, `internal/driver/consensus.go:150`).
- 5 implementation: `PhaseImpl`, caused by `IMPLEMENTATION.md`
  (`internal/driver/cursor.go:104`).
- 6 review: `PhaseReview` with review rounds and no root `review/consensus.md`.
- 7 review-consensus: `PhaseReview` with root `review/consensus.md`.
- 8 fix-up/complete: `PhaseDone` for `IMPLEMENTATION.md status: complete`
  (`internal/driver/cursor.go:100`); otherwise fix-up when review consensus has
  outstanding fixes or implementation status starts with `fix-up-cycle`, which
  the driver treats as review-ready (`internal/driver/impl.go:250`).

The problem is that `Cursor` collapses step 6 and 7: both
`review/consensus.md` and any review round produce `PhaseReview`
(`internal/driver/cursor.go:102`). Step 8 also needs implementation status or
review consensus parsing. Counterproposal: add an exported driver detail helper
that reuses the same disk pass:

```go
type PhaseDetail struct {
    Cursor driver.Cursor
    HighestReviewRound int
    ReviewConsensusExists bool
    ImplementationStatus string
    FinalScaffoldReason string
}
```

## Glyphs And Narrator

ACCEPT `agentBuffer.lastGrowthAt`, with one caveat. Buffers are lazy:
`ensureBuffer` creates/loads them for active or visited agents
(`internal/tui/live.go:1433`), and `refreshBuffers` advances only existing
buffers (`internal/tui/live.go:1515`). Therefore "zero new I/O" is true only
for visited tabs. If every tab must show spinner-vs-hourglass before visit, add
a slower Update-time per-agent size/growth cache; never read logs in View.

Set `lastGrowthAt = m.now` inside `advanceBuffer` when `readAppendedChunk`
returns bytes (`internal/tui/live.go:1494`). Current design already keeps View
cheap: View just renders cached model (`internal/tui/live.go:506`), while log
reads happen from Update (`internal/tui/live.go:409`,
`internal/tui/live.go:439`).

Glyph rule: spinner means running with recent stdout/stderr/steer growth;
hourglass means running and silent but process liveness is not stale; STALE means
`agentLiveness(agentID) == "stale"` (`internal/tui/live.go:1120`); terminal
glyphs come from projected `agent.*` events (`internal/runstate/runstate.go:450`).

Narrator reuse is good, but the current weave is only steer-specific:
`appendSteerEvents` handles `steer.replied` / `steer.reply_failed`
(`internal/tui/live.go:1544`). Add a separate `appendProtocolEvents` that writes
`transcriptEvent` lines (`internal/tui/live.go:146`) and uses the existing event
styling (`internal/tui/live.go:724`). Do not scan all `m.events` in View for
empty placeholders; append/cache filtered summaries in Update.

Perf nits: after narrator appends, re-run `capTranscriptLines`; the steer path
currently appends after `advanceBuffer` without re-capping
(`internal/tui/live.go:1558`, `internal/tui/live.go:1566`). If protocol lines
are added to unvisited buffers, do not accidentally start 250 ms log tailing for
all agents.

## `run.phase` Event

ACCEPT, but emit from the cursor-commit path used by `Advance`, not only the
outer `Run` loop. `Advance` is the phase transition API
(`internal/driver/driver.go:100`); `Run` is one caller
(`internal/driver/loop.go:23`).

Before adding the event, stop ignoring cursor save errors. Saves are currently
discarded in several phase-changing branches (`internal/driver/driver.go:146`,
`internal/driver/driver.go:164`, `internal/driver/consensus.go:69`,
`internal/driver/impl.go:68`, `internal/driver/impl.go:103`,
`internal/driver/impl.go:188`). Emit `run.phase` only after successful save.
The event append can be best-effort because reconcile is the fallback; the save
failure cannot be best-effort because then the event would lie.

Payload: `idea`, `run_id`, `action`, `phase`, `previous_phase`,
`current_round`, `round_label`, `idea_status`, `rounds_run`, `max_rounds`,
`source: "driver"`.

No new writer-concurrency concern if it uses `store.Store.Append`. The store has
an in-process append mutex and one-line `O_APPEND` writes
(`internal/store/events.go:26`, `internal/store/events.go:36`). Runner goroutines
already append agent/round events through the same path
(`internal/runner/runner.go:367`, `internal/runner/runner.go:441`), and the
driver has a single-driver lock (`internal/driver/loop.go:18`). Cross-process
append semantics remain the existing store contract.

## Required Tests

`internal/tui`: pure phase mapping tests for every cursor phase and
`cross_review_rounds` 0/1/2; temp-dir tests for step 6 vs 7 vs 8 with review
rounds, `review/consensus.md`, `IMPLEMENTATION.md status: fix-up-cycle-*`, and
`status: complete`; stale token and stale seq snapshot messages ignored after
`activateRun`; allowlist triggers schedule snapshots while ACP/steer/HITL chatter
does not; 60 simulated 250 ms event ticks without protocol events produce no
more than the expected reconcile count; View/render tests prove no disk reads;
`lastGrowthAt` covers chunks, no-growth, rotation/truncation; narrator lines
append once, exclude chatter, and stay capped.

`internal/driver`: every phase-changing action emits one `run.phase` after save;
`ActionAwait`, `ActionSurfaceOnly`, and errors emit none; save failure emits none
and returns an error; event append failure does not undo the saved cursor.

`internal/runstate`: `SummarizeEvent` renders useful `run.phase` text
(`internal/runstate/runstate.go:409`), and optional `buffers_stdout` runtime data
is parsed from `run.created` if the UI uses the buffering hint
(`internal/runcontrol/runcontrol.go:55`, `internal/runcontrol/runcontrol.go:149`).

`internal/runplan` / `internal/consensus`: if snapshot uses `runplan.Plan`, add
NEXT-line fixtures for completed round, partial consensus, ready consensus,
review consensus with fixes, and clean review consensus.

## Slice Plan

Ship as one release, but reorder internally:

1. Move `run.phase` into slice 1; the snapshot/ribbon should not wait for slice 4.
2. Build snapshot, mapping, collapsed ribbon, and `ph=` status line together.
3. Add glyphs, heartbeat counters, `buffers_stdout`, and buffered placeholders.
4. Add protocol tab, expanded ribbon, `/protocol`, `/refresh`, and `/narrate`.
   I would cut `/waiting` unless it is only an alias for `/protocol`.
5. Add Home phase column last, computed only in `refreshHomeRuns`
   (`internal/tui/live.go:454`), because Home rendering is cached
   (`internal/tui/live.go:956`).

Do not ship a ribbon without snapshot race gates, and do not ship glyphs that
force filesystem reads from View.
