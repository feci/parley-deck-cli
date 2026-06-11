---
idea: tui-protocol-visibility
agent: claude
status: implemented
date: 2026-06-12
---

## What was implemented (per FINAL.md / consensus D1-D17)

**Driver (slice 1).**
- `driver.PhaseDetail` + `RebuildDetail(ideaDir, maxRounds) (PhaseDetail, error)`
  (cursor.go); `Rebuild` is now a thin wrapper. `highestRound`/`highestReviewRound`
  gained error-surfacing variants (non-NotExist ReadDir errors propagate as the
  partial-detail error).
- `commitCursor` helper (driver.go): saves the cursor via the `saveCursor` seam
  (errors RETURNED → ActionEscalated at every one of the 9 phase-changing commit
  sites in driver.go/consensus.go/impl.go) and, only after a successful save,
  appends the best-effort `run.phase` event with the D4 payload.
- `cursor.Save` uses `fsutil.MkdirAllResilient` (virtio-fs hardening).
- `runstate.SummarizeEvent` gained `run.phase` and `agent.finished` ("wrote
  <idea-relative artifact>") cases; new `runstate.Outcome` export so the snapshot
  builder can reuse the terminal/outcome projection without `LoadRunAt`.

**Snapshot (slice 2).** `internal/tui/protosnap.go`:
- `BuildProtocolSnapshot(ProtocolSnapshotInput)` — value-copy input, never
  `LoadRunAt`; one `RebuildDetail` pass + 00-prompt frontmatter +
  `consensus.Status` (design schema at steps 3-4, review at 7+) + `runplan.Plan`
  NEXT; pure `displayStep` cursor→9-step mapping; delivery merge events-primary/
  disk-fallback with `Unvalidated` rows and Phase-6 waiting = participants minus
  implementer (implementer read from IMPLEMENTATION.md frontmatter);
  keep-last-on-error; two-consecutive-agreement step-regression guard.
- live.go integration: `protoMsg`/`protoTickMsg` gated by runToken + protoSeq,
  single in-flight + protoDirty coalescing; reconcile tick 15s/60s
  (`protoInterval`); triggers = D6 allowlist batch check, doneMsg, /protocol tab,
  /refresh. Zero snapshot work on the 250ms/1s ticks (regression-tested).

**Surfaces (slices 3-6).**
- Tab glyphs (D8): `agentGlyph` — ○ pending, braille spinner flowing, · silent,
  ✓ delivered, ✗ failed, x killed, - skipped, ! STALE. `agentBuffer.lastGrowthAt`
  (set in advanceBuffer) + 2s `statGrowthCmd` stat cache for unvisited tabs
  (`growthPaths` skips loaded buffers). `shortState` became orphaned and was
  removed with its test (replaced by a glyph test).
- Ribbon (D9): `renderRibbon`/`renderRibbonExpanded`, Ctrl+P cycles
  collapsed→expanded→hidden; `[STALE]` prefix, `?` disk-fallback marker,
  reconciled-age >30s on the collapsed line; transcriptHeight accounts for the
  ribbon rows.
- Status line (D10): `statusPhaseSegment` — `ph=2:xrev-r02 wait=agy,hermes`,
  legacy `round=` fallback until the first snapshot.
- Narrator (D7): `appendProtocolEvents` + 32-entry replay ring with per-buffer
  seq dedup; re-cap after appends (and the pre-existing steer-path append now
  re-caps too); `/narrate` cycles protocol→verbose→off; verbose weaves the
  agent's own `agent.acp.*` (chunks excluded) into its tab only.
- Buffered-agent placeholder (D12): `renderSilentPlaceholder` — declared
  `buffers_stdout` (new field on agents.Spec, TOML `buffers_stdout`, default
  true for agy, threaded through run.created runtime payload and cached via
  `noteRuntimeFlags`) with the 30s/0B heuristic fallback; liveness, byte
  counters, last 5 own event summaries.
- Heartbeat (D13): running header shows proc:live/proc:stale, stdout/stderr byte
  counters, artifact target.
- Protocol tab (D11): `renderProtocolPanes` (PIPELINE 0-8, DELIVERY, SIGNOFFS,
  NEXT) above the existing panes; `/protocol` aliases `/status`. No dependency
  on review-cycle frontmatter.
- Home (D14): `ideaPhaseChip` per idea computed only in `refreshHomeRuns`.

**Version:** 1.23.0 (version.go + VERSION + CHANGELOG entry).

## Tests (D17)

- driver: `phase_event_test.go` — run.phase emission on promoted and
  consensus-drafted commits with full payload; none on await/surface-only; save
  failure (seam) escalates and emits nothing; `RebuildDetail` review-evidence
  golden + `Rebuild` wrapper equivalence.
- tui: `protosnap_test.go` — 10-case displayStep golden table; delivery
  merge-rule matrix over a real temp idea dir (disk-only → delivered? +
  DiskFallback, event-only → delivered?, neither → waiting); keep-last on an
  unreadable review dir; regression needs two agreeing reconciles; tick-budget
  (60 event + 60 elapsed ticks → zero snapshot schedules) + trigger/dirty
  coalescing + stale-seq drop; `ph=` status grammar; narrator weave/replay/
  dedup/off; spinner-vs-silent-vs-stale glyph matrix.
- Full suite green, `go test -race` green on internal/tui + internal/driver.

## Deviations from FINAL.md

- None functional. Notes: (a) the `?` disagreement marker renders as `✓?` in
  delivery cells and as a `?` suffix on the ribbon count, per agy's round-01
  string; (b) `narrateProtocol` is the zero-value default so a fresh model
  narrates without initialization; (c) gofmt applied to new files only —
  live.go/app.go/procctl were already unformatted on main and were left alone.

## Sandbox note for reviewers

`internal/runner TestDurableKillEndToEndRealProcess` fails under codex's
seatbelt sandbox ("no recorded boot id" — sysctl kern.boottime is restricted).
It passes in a normal shell; this is the known recurring sandbox artifact, not
a defect of this change.
