---
idea: tui-round-summary
status: fix-up-cycle-1
implementer: claude-1
started: 2026-07-04
completed: 2026-07-04
branch: parley-deck-cli#round-summary-design
head-commit: f87bf93
design-pr: https://github.com/feci/parley-deck-cli/pull/71
implementation-pr: same
---

## Summary of work

A consolidated round digest in the TUI Home tab, built in the driver (the holder of the
canonical round-completion predicate) and carried on a `round.digest` event; Home is a
pure renderer. Deterministic, LLM-free; stance shown as keyword mentions, never verdicts.

## Implementation plan / checklist

- [x] `internal/driver/digest.go` (new): pure `BuildRoundDigest` — per-agent `## Summary`
      first-sentence extraction (cap 120, sentence-boundary), fallback to first prose
      paragraph tagged degraded, X/Y completeness, keyword `stanceFlags` (block/counter/
      accept/escalate) as HINTS. Never errors (missing artifact → not-present line).
- [x] `internal/driver/driver.go`: `emitRoundDigest` appends an idempotent `round.digest`
      event (JSON blob) once per (idea, round) at the round-complete point in
      `advanceRound`, before deciding the next action; failure is swallowed (a display
      feature never blocks advancement).
- [x] `internal/tui/roundsummary.go` (new): pure renderer over already-consumed events —
      `latestRoundDigest` decodes the newest `round.digest`; `renderRoundDigest` formats a
      hard-bounded block (total lines ≤ maxRows) with a `mentions:` (not verdicts) line
      and a `next:` line.
- [x] `internal/tui/live.go` `renderHome`: renders the latest digest in a viewport-bounded
      sub-block (rows/3, ≤10) BEFORE "Recent runs" so it never pushes chips/roster/runs
      off-screen (the regression guard).
- [x] Tests: `internal/driver/digest_test.go` (well-formed, missing-Summary fallback,
      missing agent X/Y, long-line cap, no-double-count-blocker), `internal/tui/
      roundsummary_test.go` (decode, bounded rows, mentions-not-verdicts, zero-rows).
- [x] Checks: `go build ./...`, `go vet`, `gofmt -l` clean; `go test ./internal/driver
      ./internal/tui` green.

## Deviations from FINAL.md

- Render carries the digest via `m.events` (the events the TUI already consumes) rather
  than extending `ProtocolSnapshot`. Same "single source of truth" outcome (the driver
  emits, Home renders) with a smaller surface and no snapshot-pipeline change. Documented
  deviation.
- History: v1 renders the LATEST digest only (bounded block), not a scrollable last-N.
  This is the viewport-safe minimum that satisfies "catch up without tab-flipping";
  older digests remain in the event log and per-agent tabs. Last-N scroll is a deferred
  follow-up (avoids the Home-overflow regression the consensus flagged as the top risk).
- Round 2+ engagement-gap extraction from `responding-to:` is deferred with the last-N
  history (both are additive digest enrichments); the flags line + completeness + next
  deliver the core position map for every round.

## Notes for reviewers

- `emitRoundDigest` idempotency scans prior events for a matching (idea, round)
  `round.digest`; confirm the key is right and re-runs never duplicate.
- `stanceFlags` counts "block" (which subsumes "blocker") once — verified by
  `TestStanceFlagsNoDoubleCountBlocker`.
- The render is hard-capped to maxRows (header + agent rows + trailer), verified by
  `TestRenderRoundDigestBoundedAndLabelsHints` with a tiny budget.

## Fix-up cycle 1
status: complete
completed: 2026-07-04

### Fixes applied (review round-01)
- [MAJOR, codex-1] renderHome now RESERVES the rows the sections below the digest need
  (Recent runs header + runs + hint) and gives the digest only the genuine leftover, so
  it can never push Recent runs off-screen. Regression test
  `TestRenderHomeReservesRunsBelowDigest` (short terminal + 3 ideas asserts "Recent runs"
  still present).
- [MINOR, codex-1] `extractPosition` cap is now rune-safe (`[]rune`) — no UTF-8 rune
  splitting.
- [MINOR-1, hermes-1] the degraded-extraction marker is now a leading `~` PREFIX (before
  the position) so truncation can never erase the honesty signal.
- [MINOR-2, hermes-1] the digest `next:` line now says "consensus-ready (manual draft)"
  when the consensus gate is unwired (driver halts, does not draft), instead of the
  misleading "drafting consensus".
- [NIT-1, hermes-1] aggregate `mentions:` line without per-agent attribution is the
  ratified design — no change.

### Deviations from agreed fixes
None.
