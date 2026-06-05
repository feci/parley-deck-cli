---
idea: deliberation-driver
status: complete
implementer: claude
started: 2026-06-05
completed: 2026-06-05
branch: parley-deck-cli#feature/deliberation-driver
head-commit: see-branch-tip
design-pr: n/a (local-dir transport)
implementation-pr: n/a (local-dir transport)
---

## Summary of work

Implemented the deliberation-driver **minimal first slice (FINAL D15)**: a new
`internal/driver` package that owns the next tick after round-01 and promotes a
completed round-01 → round-02 via `runner.RunRound`, wired behind
`parley run --auto --no-tui` on a local-dir workspace. This fixes the months-long
stall where `parley run` was a one-shot batch executor.

## Implementation plan / checklist

- [x] `internal/driver/cursor.go` — `Cursor` + atomic `Save`/`LoadCursor` +
      `Rebuild` (disk-derived phase, D2/D3); frontmatter/round helpers.
- [x] `internal/driver/driver.go` — `RoundRunner` seam, `Config`, `Driver`,
      `Advance` (one re-entrant tick), `roundComplete` (two-signal gate +
      reconciliation, D4), `setIdeaStatus` (writes the round-02 status nothing
      else wrote), `NewRunnerAdapter` (wraps `runner.RunRound`, Overwrite=false).
- [x] `internal/driver/loop.go` — `Run` poll loop + mandatory advisory
      `driver.lock` (PID, D10) + deadline escalation (D12).
- [x] `internal/driver/transport.go` — `EffectiveTransport` (idea-level →
      global, D8) + `ReadCrossReviewRounds` (default 1, D5).
- [x] `internal/driver/driver_test.go` — 9 table tests (D14) with a fake
      RoundRunner; no live agents.
- [x] `internal/runaction/action.go` + `internal/runplan/runplan.go` —
      `KindOpenNextRound` (D13 visibility): `parley continue` now surfaces "open
      round-02" instead of jumping a completed independent round to consensus.
      Updated `TestPlanDraftsConsensusWhenRoundArtifactsExist` (now tests the
      `cross_review_rounds: 0` bypass) + added
      `TestPlanOpensNextRoundAfterCompletedFirstRound`.
- [x] `internal/app/app.go runTask` — `--auto && local-dir` hands off to
      `driver.Run` after RunRoundOne in the `--no-tui` branch; non-auto unchanged.
- [x] Checks: `go build ./...`, `go vet ./...`, `go test ./...` — all green.
- [x] Real-run acceptance (see below).

## Tests

`internal/driver/driver_test.go`: promote round-01→round-02 (RunRound(2) called
once, status=round-02, cursor saved); no duplicate dispatch when round-02 already
complete; await on incomplete round; `round.incomplete` event blocks promotion;
reconcile missing terminal event (reconstructed `round.completed` appended);
Rebuild derives phase from disk (round/consensus/final); corrupt cursor ignored →
Rebuild recovers; surface-only when not auto+local-dir; `cross_review_rounds: 0`
bypass → consensus.

## Real-run acceptance (D15)

`parley run --auto --no-tui --participants codex,agy` on a fresh local-dir
workspace (`parley init` → COOPERATION.md transport=local-dir).

**Acceptance result (PASS):** the run wrote `round-01/{codex,agy}.md`, then the
driver opened round-02 and wrote `round-02/{codex,agy}.md` (each with `round: 2` +
`responding-to` frontmatter — passing the round≥2 gate), set `00-prompt.md`
status=`round-02`, and `events.jsonl` records **two** `round.completed` events
(round-01 AND round-02 — the follow-on event absent before this change). The
driver then detected the cross-review budget was spent and stopped cleanly at the
consensus boundary (`driver: cross-review complete at round-02; next step is
parley consensus draft …`) with no re-dispatch. Log: `/tmp/dd-acceptance-run.log`.

Note: the runner emits `round.completed`/`round.incomplete` (+ `run.segment_started`,
`round.index_written`), not a distinct `round.started`; the "second round pair" the
brief asked for is realized as the second `round.completed` + its segment-start.

## Deviations from FINAL.md

- **(RESOLVED in fix-up cycle 1, AF2)** Slice 1 initially weakened the round≥2 gate
  to `responding-to` presence because the runner prompt did not emit per-agent
  headings. Per review consensus AF2 this was reversed: `runner.BuildRoundPrompt`
  now requires a `### @<other participant>` subsection for each other participant,
  and the driver re-enforces the full D4 check (`### @<other>` heading for every
  other participant + `responding-to`). No remaining deviation from D4.

## Fix-up cycle 1 (Phase 8)

status: complete
head-commit: see-branch-tip

Applied the Phase 7 review-consensus agreed fixes (review/consensus.md):

- **AF1 — Atomic driver.lock** (`loop.go`): `acquireLock` now creates the lock
  with `os.OpenFile(O_CREATE|O_EXCL|O_WRONLY)`, refuses if held by this PID or a
  live PID, reclaims only a different dead PID, and releases only when the file
  still carries our token. New `TestAcquireLockIsExclusive` (8 goroutines race →
  exactly one holder).
- **AF2 — Restore D4 cross-review gate**: `runner.BuildRoundPrompt` emits a
  `### @<other>` subsection per other participant; `driver.validateCrossReviewBody`
  enforces it for round≥2; `terminalRoundEvent` now matches idea+round (codex OQ).
  New `TestRound02RequiresCrossReviewHeadings`; fake artifacts carry the headings.
- **AF3 — Durable escalation** (`loop.go`): generalized `escalateDeadline` into
  `escalate(cursor, topic, body)` writing a blocking inbox note; `Run` routes
  `Advance` errors (malformed event log, runner failure) through it instead of
  stderr-only.
- **AF4 — Robust per-tick transport** (`transport.go`, `driver.go`):
  `EffectiveTransport` global fallback now uses `protocol.ReadWorkspaceStatus`
  (tolerates no-backtick `**Transport:** local-dir`); `Config.AutoLocalDir` was
  replaced by `Root`+`Auto` and `Advance` re-reads the effective transport from
  disk every tick (D8). New `TestAdvanceSurfaceOnlyWhenTransportNotLocalDir`.

Checks after fix-up: `go build ./...`, `go vet ./...`, `go test ./...` — all green.
**Real-run acceptance re-verified** (`/tmp/dd-accept2`): with the AF2 prompt change,
codex and agy each wrote a valid `round-02/<id>.md` carrying a `### @<other>`
heading (codex→@agy, agy→@codex), the stricter D4 gate promoted to round-02,
status=`round-02`, and the driver stopped at the consensus boundary. The
strengthen-the-prompt-then-enforce approach works with live agents.
- **Slice scope:** only the round phase is driven (D15). Consensus/final/impl gates
  (D6/D7), the `internal/signoffs` extraction (D9), and MaxRounds reopen (D11) are
  later slices (S2–S5 in FINAL.md); `Advance` returns `surface-only` once the cursor
  leaves the round phase, so the driver never auto-drives consensus yet.

## Notes for reviewers

- Disk is authoritative; the cursor (`<runDir>/driver.json`) is a rebuildable
  cache; `driver.lock` (PID) enforces single-driver. Re-entry is idempotent
  (Overwrite=false + fileExists/event guards) — a duplicate tick or crash-restart
  does not re-dispatch a completed round (test:
  `TestAdvanceNoDuplicateAfterRound02Complete`).
- The two-signal gate reconstructs `round.completed` ONLY after every artifact
  validates (never on file presence alone); a malformed event log escalates rather
  than guessing.
- Transport is read per tick: idea-level `transport:` → COOPERATION.md global;
  auto-advance only when effective transport is local-dir AND `--auto`.
- The driver imports `internal/runner` + `internal/store` only — never
  `internal/app` (the D9 extraction is deferred with the signoff slice, which
  slice 1 does not need).

## Phase 8 complete (slice 1) — 2026-06-05

Fix-up cycle 1 re-review: codex, agy, hermes all signed ✅ ACCEPT (review/consensus.md
+ review/round-02/). Zero remaining agreed fixes. **Slice 1 (round-phase driver) is
complete**: the months-long round-01 stall is fixed and proven with live agents.
Later slices S2–S5 (consensus→final→impl→review auto-drive + internal/signoffs
extraction) remain as scoped follow-ups per FINAL.md.
