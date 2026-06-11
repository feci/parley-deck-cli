---
idea: tui-protocol-visibility
status: final
drafted-by: claude
date: 2026-06-12
---

## Summary

Make the live TUI answer two questions at all times: "are the other agents alive
and working?" and "where in the Parley Deck protocol is this idea?" Five
surfaces — a collapsible protocol ribbon, tab activity glyphs, narrator lines,
a Protocol tab grown from the Status tab, and a Home phase column — all fed
from in-memory events plus a cached async disk snapshot. run.json is never
read. Consensus: claude/codex/agy ✅ ACCEPT; hermes signed round-02 "ACCEPT,
no blockers" but its CLI failed at the signoff append (recorded inbox
exception, signoff outstanding).

## Final plan / specification

The authoritative decision record is consensus.md D1-D17; this section restates
the build contract.

**Driver (D2-D4).** Export `PhaseDetail{Cursor, HighestReviewRound,
ReviewConsensusExists, ImplementationStatus, FinalScaffoldReason}` and
`RebuildDetail(ideaDir, maxRounds) (PhaseDetail, error)`; keep `Rebuild` as a
compatibility wrapper; export no private helpers and add no Cursor fields.
Add one `commitCursor` helper used by every phase-changing Advance branch
(driver.go:146,164; consensus.go:69,117; impl.go:68,103,130,188,226): it saves
the cursor (errors now RETURNED → ActionEscalated, no longer discarded), and
only after a successful save appends a best-effort `run.phase` event {idea,
run_id, action, phase, previous_phase, current_round, round_label, idea_status,
rounds_run, max_rounds, source:"driver"}. cursor.Save's MkdirAll switches to
fsutil.MkdirAllResilient.

**Snapshot (D5-D6).** New internal/tui/protosnap.go with
`BuildProtocolSnapshot(in ProtocolSnapshotInput) (ProtocolSnapshot, error)`
taking explicit value copies (Root, RunID, RunDir, IdeaSlug, IdeaDir,
Participants, MaxRounds, CrossReviewRounds, Events, Questions, State, Previous,
Now); never calls runstate.LoadRunAt. Inputs: driver.RebuildDetail + in-memory
events; consensus.Status only at phase 3+/7+ (design/review schema per phase);
runplan.Plan for the NEXT line. Pure cursor-to-9-step mapping per consensus D3.
Async tea.Cmd gated by runToken + protoSeq, single in-flight with a dirty-flag
coalescer; separate reconcile tick 15s (attached and running) / 60s
(done, STALE, or detached); triggers: allowlisted event batch (D6 list),
doneMsg, Protocol-tab switch, /refresh. Zero snapshot work on the 250ms event
tick or the 1s elapsed tick. Keep-last on error with a reconcile_error marker;
phase regression requires two consecutive agreeing reconciles. Delivery merge:
events primary, disk fallback, question-mark marker on disagreement; Phase-6+
waiting sets exclude the implementer.

**Surfaces (D8-D14).** Tab glyphs: pending ○, active braille spinner (fallback
*), silent ·, delivered ✓ (artifact_ok only), failed ✗, killed x, skipped -,
STALE !; wide glyphs only in ribbon/placeholder/Home copy. Spinner from
agentBuffer.lastGrowthAt plus an async 2s stat growth cache for unvisited tabs;
View never reads the filesystem. Ribbon collapsed: "◆ Ph 2: Cross-Review (R02)
· Delivered 1/3 · Waiting: agy, hermes · Next: consensus"; [STALE] prefix,
disk-fallback question-mark suffix, reconciled-age surfaces when older than
30s; Ctrl+P cycles collapsed/expanded/hidden with the 3-line
Pipeline/Delivery/System expansion. Status line: `ph=2:xrev-r02
wait=agy,hermes` replacing round=. Narrator (D7): separate display allowlist,
appended to loaded buffers + a 32-event replay ring with per-buffer dedup,
re-capped after append (also fixes the steer-path missing cap);
/narrate cycles off/protocol/verbose. Buffered-agent placeholder: frameless
block with buffering notice (declared `buffers_stdout` flag on
agents.Spec/TOML/run.created runtime, default true for agy; heuristic
fallback), elapsed + liveness, byte counters, last 5 own event summaries.
Heartbeat header: proc:live/proc:stale + byte counters + artifact path.
Protocol tab: PIPELINE/DELIVERY/SIGNOFFS/NEXT panes above the existing panes;
/protocol aliases /status; no dependency on the review-cycle frontmatter
field. Home: per-idea phase chip + Attention badge, computed only in
refreshHomeRuns.

**Tests (D17).** Golden 9-step mapping; merge-rule table; tick-budget
regression (max 4 snapshot computes across 60 simulated event ticks); stale
token/seq drops; allowlist trigger tests; View-no-disk-reads;
lastGrowthAt growth/rotation; narrator dedup + cap; driver emission matrix
(one run.phase per phase-changing action after save; none on await/surface/
error; save failure returns error and emits nothing); race detector over
snapshot vs kill/relaunch.

## Implementation slices (one release: 1.23.0)

1. driver PhaseDetail + save-error handling + run.phase.
2. Snapshot + collapsed ribbon + ph= status line.
3. Glyphs + growth cache + heartbeat + buffers_stdout.
4. Narrator + silent-agent placeholder.
5. Protocol tab + expanded ribbon + Ctrl+P//protocol//refresh//narrate.
6. Home phase column.

## Implementer

claude (FINAL drafter; default per protocol Phase 5). Reviewers: codex, agy,
hermes (if recovered).

## Deferred follow-ups

- impl.checks event (checks-gate visibility).
- review-cycle frontmatter naming normalization.
- ACP tool-call lines in the silent-agent placeholder.
- hermes signoff append (outstanding; tooling exception in inbox).
