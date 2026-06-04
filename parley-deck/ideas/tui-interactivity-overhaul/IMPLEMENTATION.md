---
idea: tui-interactivity-overhaul
status: implementing
implementer: claude
started: 2026-06-04
completed: 2026-06-04
branch: parley-deck-cli#feat/tui-interactivity-overhaul
head-commit: pending
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
- [ ] **Slice 2 — per-agent focus viewport** (`bubbles/viewport`, offset reads,
      follow, bounded scrollback 20k lines / 4 MiB, `g`/`G`/`f`).
- [ ] **Slice 3 — view-state machine + keymap + help overlay** (`overview |
      agentDetail | compose | answerQuestion | help`; `?` help; `a` preserved).
- [ ] **Slice 4 — steering composer + `steer.*` events + `parley steer` CLI**
      (`SubmitSteering`, queued `new_attempt` default, no gate bypass).

- [ ] Checks to run: `go build ./...`, `go vet ./...`, `go test ./...` per slice.
- [ ] Review notes: see Deviations below.

## Deviations from FINAL.md

- **Segment-event tagging on the ACP launch path is deferred.** Slice 1 tags the
  standard subprocess events (`runAgent`, fixup). ACP agents (`acp.go`) still get
  correct `State` AND `Segment` because the `run.segment_started` reset barrier
  sets `Segment`, and `applyAgentEvent` only overrides it when an event carries
  `segment_id`. So the projection is correct for ACP without tagging its events;
  explicit ACP-event tagging is left for slice 5 (per-attempt history), per the
  simplicity-first rule. The reset-barrier — not per-event tagging — is what
  fixes the badge.

## Notes for reviewers

- Old, unsegmented runs emit no `run.segment_started`, so `ProjectEvents`
  behavior is unchanged for them (verified by `TestProjectEventsLegacyUnsegmentedUnchanged`
  and the untouched existing runstate tests).
- `nextSegmentID` counts prior segment events via `store.Load()`; round-runs are
  sequential per run so the monotonic count is stable. Cross-process append
  safety (deferred follow-up F4) is unchanged by this slice.
