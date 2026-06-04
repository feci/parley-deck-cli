---
idea: tui-interactivity-overhaul
status: complete
implementer: claude
started: 2026-06-04
completed: 2026-06-04
branch: parley-deck-cli#feat/tui-interactivity-overhaul
head-commit: see-branch-tip
design-pr: https://github.com/feci/parley-deck-cli/pull/33
implementation-pr: https://github.com/feci/parley-deck-cli/pull/33
---

## Summary of work

Implementing the FINAL.md committed scope (slices 1–4) for the interactive
live-run TUI redesign. Each slice builds + the full `go test ./...` suite stays
green before moving on.

## Implementation plan / checklist

- [x] **Slice 1 — segment projection + sticky `[FINISHED]` fix**
  - [x] `runstate`: `AgentState.Segment`; `run.segment_started` resets target
        agents in `ProjectEvents`; `applyAgentEvent` records `segment_id`;
        `SummarizeEvent` case.
  - [x] `runner`: `Options.SegmentID`; `appendSegmentStarted`/`nextSegmentID`
        (monotonic `segment-NNNN`)/`segmentReason`/`agentIDs`; emit at start of
        `RunRoundOne` (covers RunRound + RunReviewRound) and in
        `RunImplementation`/`RunFixup`(retry)/`RunReviewConsensus`; tag
        `agent.started|finished|failed|skipped`.
  - [x] Tests: `internal/runstate/segment_test.go` (legacy-unsegmented,
        unstick-finished, continue→running, non-targeted-untouched, skip,
        retry, unknown-target no-op); updated `TestRunRoundOneSkipsExistingArtifact`
        for the leading segment event.
  - [x] Build + vet + `go test ./...` green.
- [x] **Slice 2 — per-agent focus viewport**
  - [x] `enter`/`o` opens an `agentDetail` focus view; `esc` returns; `tab`
        cycles the focused agent.
  - [x] Offset-incremental reads (`loadFocusTail` + `readAppendedLines`) over the
        focused agent's full stdout log; bounded scrollback (20k lines / 4 MiB)
        with a truncation marker; ANSI stripped.
  - [x] Follow mode (`f`, default on; manual scroll disables it); `g`/`G`
        top/bottom; `j`/`k` line scroll; `pgup`/`pgdn` page; reload on log
        truncation/rotation (new segment).
  - [x] Tests: focus render/exit, follow+scroll, bounded-lines cap, incremental
        append. Full `go test ./...` green.
  - [x] **Deviation:** implemented a lightweight in-house viewport instead of
        adding the `bubbles/viewport` dependency (not currently in go.mod) —
        simpler, dependency-free, sufficient for line-oriented logs. Stream tabs
        (stdout/stderr/thoughts) deferred; the focus view shows stdout (the
        agent's working output); stderr stays in the overview log preview.
- [x] **Slice 3 — view-state machine + keymap + help overlay**
  - [x] Replaced the overloaded `answerMode`/`focus` booleans with a single
        `mode` enum (`modeOverview|modeAgentDetail|modeCompose|modeAnswerQuestion|
        modeHelp`); `modeCompose` reserved for Slice 4.
  - [x] `?` opens a help overlay (full keymap); `esc`/`?`/`q`/`enter` dismiss it;
        `a` remains the sole HITL answer key (preserves `hitl-tui-questions`).
  - [x] Footer advertises `? help` (overview + resume).
  - [x] Tests: help-overlay toggle; migrated focus/answer tests to the `mode`
        enum. Full `go test ./...` green.
- [x] **Slice 4 — steering composer + `steer.*` events + `parley steer` CLI**
  - [x] New `internal/steer` package: `Request`/`Result`/`Queued`, `Submit`
        (records `steer.requested` + `steer.delivered` with monotonic
        `steer-NNNN` ids, queued `new_attempt` mode), `List` projection.
  - [x] TUI composer (`modeCompose`): `i` steers the selected agent, `I` the
        deck; type + enter queues via `steer.Submit`, esc cancels; queued
        confirmation shown in the footer; help + footer updated.
  - [x] `parley steer [--dir D] [--agent A] [--json] RUN -- TEXT...` CLI;
        `continue` surfaces queued steers (text + `--json` `steers`).
  - [x] `runstate.SummarizeEvent` cases so `steer.*` read well in the events
        pane. Tests: steer package, TUI composer (queue + cancel), CLI + continue.
        Full `go test ./...` green.
  - [x] **Deviation:** the composer persists the steer intent directly to the
        run event log (mirroring how the live TUI already records HITL answers),
        rather than through a `SubmitSteering` callback. Recording intent has no
        side effect and bypasses no gate. The driver-owned execution boundary
        (FINAL D5) applies to slice 5, where the queued attempt is actually
        launched / delivered live and goes through the gates.

- [x] Checks run per slice: `go build ./...`, `go vet ./...`, `go test ./...` all green.
- [x] Review notes: see Deviations above and Notes for reviewers below.

## Deviations from FINAL.md

- The steering composer persists intent directly to the run event log (mirroring
  the live TUI's existing HITL-answer write) rather than via a `SubmitSteering`
  callback; recording intent has no side effect and bypasses no gate. The
  driver-owned execution boundary (FINAL D5) applies to slice 5, when the queued
  attempt is actually launched / delivered live through the gates.
- Per-stream (stdout/stderr) scrollback tabs in the focus view are deferred
  (slice 5); the focus view shows stdout, stderr stays in the overview preview.
- (Resolved in fix-up cycle 1) ACP segment tagging — now applied (AF6).

## Fix-up cycle 1
status: complete
completed: 2026-06-04
head-commit: see-branch-tip

Applied all agreed fixes from review/consensus.md (AF1–AF9):

### Fixes applied
- AF1 — composer/help/CLI wording no longer over-promises execution: steers are
  "recorded / queued (auto-exec not wired yet)".
- AF2 — focus buffer is byte-bounded as well as line-bounded: `readAppendedLines`
  caps the per-tick read window to `maxFocusBytes`, `capFocusLines` evicts oldest
  by both bytes and lines. Tests for an oversized line and the byte budget.
- AF3 — steer ids are collision-resistant (`steer-NNNN-<rand>`); the by-id
  projection never merges two distinct steers under concurrent CLI/TUI submits.
- AF4 — `loadFocusTail`/`readAppendedLines` advance the offset only past complete
  lines (shared `completeLinesFrom`), so a partial line is merged, not fragmented.
- AF5 — untagged `agent.*` events inherit the current segment in `ProjectEvents`.
- AF6 — ACP `agent.started`/terminal events are tagged with `segment_id` (closes
  the prior deviation).
- AF7 — added a `LoadRunAt` integration test asserting the badge unsticks through
  the real events.jsonl path.
- AF8 — `parley steer` validates empty text and warns on a non-participant agent.
- AF9 — overview preview fills available rows (height-derived) instead of a hard 6.

### Deviations from agreed fixes
None. Carried deferrals: F4 (cross-process atomic append/locking) and slice 5
(executing queued steers, live delivery, opt-in thoughts, per-stream tabs).

## Fix-up cycle 2
status: complete
completed: 2026-06-04
head-commit: see-branch-tip

### Fixes applied
- AF2 (round-02 follow-up, codex MAJOR) — `capFocusLines` now head-truncates a
  single retained line that still exceeds `maxFocusBytes` (UTF-8-safe via
  `strings.ToValidUTF8`), so the byte invariant holds even for one oversized
  line. Added `TestCapFocusLinesTruncatesSingleOversizedLine`.

### Deviations from agreed fixes
None.

## Notes for reviewers

- Old, unsegmented runs emit no `run.segment_started`, so `ProjectEvents`
  behavior is unchanged for them (verified by `TestProjectEventsLegacyUnsegmentedUnchanged`
  and the untouched existing runstate tests).
- `nextSegmentID` counts prior segment events via `store.Load()`; round-runs are
  sequential per run so the monotonic count is stable. Cross-process append
  safety (deferred follow-up F4) is unchanged by this slice.
