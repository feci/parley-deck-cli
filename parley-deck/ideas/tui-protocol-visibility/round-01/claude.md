---
agent: claude
idea: tui-protocol-visibility
round: 1
date: 2026-06-11
---

## Summary

The transcript stays the hero; protocol awareness wraps around it in five surfaces:
a one-line protocol ribbon on every tab, activity glyphs in the tab strip, narrator
lines woven into transcripts, a Protocol tab grown out of the existing Status tab,
and a Home phase column. Everything is fed from events.jsonl (already streamed to
the TUI) plus a cached, low-frequency disk snapshot; run.json is never read. Exactly
one new event type (`run.phase`) is needed; every other datum already exists in the
codebase. This round-01 is the full design; I prepared it by surveying the actual
code paths cited below.

## 1. Why the user only sees codex

Every participant already gets a tab (live.go:1320-1321) and three tailed streams
(stdout/stderr/steer, live.go:486-490). The failure is presentational:

- agy-style headless agents buffer stdout until exit, so their tab renders the bare
  "no output yet" placeholder (live.go:696-698) for minutes — indistinguishable from
  a hung agent.
- The tab strip shows state badges (RUN/FIN/ERR + STALE, live.go:593-600) but nothing
  distinguishes "running and producing output" from "running and silent".
- Nothing tells you another agent delivered while you watch the current tab.

## 2. Why the protocol state is invisible

The status line shows only `round=<running|completed>` (live.go:866-868, RoundStatus
from runstate.go:46). Home shows the raw idea status string (live.go:969). There is
no phase pipeline, no per-participant delivery, no signoff state. Yet the inputs all
exist: `driver.Rebuild` derives the phase cursor from disk (driver/cursor.go),
`runstate.ProjectEvents` projects per-agent state from events (runstate.go:325),
`consensus.Status` parses signoff blocks, `runplan.Plan` computes next actions
(already wired into RunSummary.NextActions, runstate.go:165-175).

## 3. The five surfaces

### 3.1 Protocol ribbon (every tab, collapsible)

```
 Home │ codex ✓ │ agy ⠙ │ hermes ⏸ │ Protocol                       run 1f3c9a · tui-aw
─────────────────────────────────────────────────────────────────────────────────────────
 ▸ round-02/2 (cross-review) · delivered 1/3 · waiting: agy, hermes · next: consensus ⌃P
```

Ctrl+P cycles collapsed → expanded → hidden. Expanded adds the pipeline strip
(`kickoff ✓ ▸ round-01 ✓ ▸ [round-02] ▸ consensus ▸ …`), a per-agent delivery line
with timestamps/sizes, and `next: … · reconciled HH:MM:SS`. Once the snapshot age
exceeds a threshold, the `reconciled …` marker surfaces on the collapsed ribbon too —
staleness honesty must not hide behind Ctrl+P, precisely when a detached/STALE run
needs it.

### 3.2 Tab activity glyphs

Spinner ⠙ = output grew <5s ago; hourglass ⏸ = running but silent; ✓ = artifact
delivered (`agent.finished` with `artifact_ok=true`, runner.go:437-447 — distinct
from FIN-with-bad-artifact); ✗/⊘/† = failed/skipped/killed; `!` = STALE overlay
(existing liveness seam live.go:1122-1127). Implemented by adding `lastGrowthAt` to
`agentBuffer`, set when `advanceBuffer` (live.go:1482-1513) moves any tail cursor;
animated by the existing 1s elapsed tick (live.go:2518-2522). Zero new ticks, zero
new I/O.

### 3.3 Narrator lines + buffered-agent placeholder

Protocol events woven into transcripts as dim rule-lines, reusing the steer-weave
mechanism shipped in 1.20.0. Allowlist: run.created, run.segment_started,
agent.started/finished/failed/skipped/killed, round.completed/incomplete,
round.index_written/failed, steer.*, hitl.question/answered, agent.fixup_*,
run.failed, run.manifest_deferred, and the new run.phase. Excluded: agent.acp.*
chunk chatter (opt-in via /narrate verbose). Text via runstate.SummarizeEvent
(already exported, used for RunState.Recent).

For a running-but-silent agent the empty placeholder becomes useful: render that
agent's own filtered event summaries (from in-memory m.events — zero extra I/O) plus
an explicit hint: `⠙ no stdout yet — agy buffers all output until exit; stderr above
is live`. Heuristic v1: State==Running && stdoutOffset==0 && elapsed>30s; upgraded
to a declared fact by the `buffers_stdout` flag (§5).

### 3.4 Protocol tab

The existing Status tab (statusTabID, live.go:132; renderStatusTab live.go:815-825)
grows three panes above the current agent/event/question panes:

```
 PIPELINE  tui-aw                                              phase: review-consensus
   [✓] 0 kickoff          00-prompt.md                         06-10 18:02
   [✓] 1 round-01         3/3 delivered                        06-10 19:14
   [✓] 2 round-02         3/3 delivered · round.completed      06-11 12:51
   [✓] 3 consensus        triage: ready (3 accept)             06-11 13:20
   [✓] 4 FINAL            FINAL.md                             06-11 13:24
   [✓] 5 implement        IMPLEMENTATION.md                    06-11 14:40
   [✓] 6 review round-01  2/2 reviewers delivered              06-11 15:12
   [▶] 7 review-consensus triage: partial · waiting: hermes
   [ ] 8 fix-up/complete  —
 DELIVERY (current round)                SIGNOFFS (consensus.Status)
 NEXT (runplan.Plan first action)
```

`/protocol` and `/status` both land here (alias).

### 3.5 Home phase column

Each idea row gains a phase chip + attention badge: `tui-aw  ph7 review-consensus ·
wait:hermes  ⚠action`. Computed only inside refreshHomeRuns (live.go:455-466) — on
demand, never on a tick. This closes the cold-open/detached orientation gap without
a Home rework.

## 4. ProtocolSnapshot — single producer

New file `internal/tui/protosnap.go`:

- **Inputs:** idea dir via `driver.Rebuild(ideaDir, maxRounds)` (7-state cursor;
  maxRounds from `cross_review_rounds` via the public protocol.ReadFrontmatter —
  not runplan's unexported reader); events already in m.events (ProjectEvents);
  participants from m.opts, falling back to the run.created payload
  (runcontrol.go:55-67) and finally 00-prompt.md frontmatter. Never run.json.
- **Cursor→display mapping** (pure func, unit-tested): PhaseRound+round-01→step 1;
  PhaseRound+N≥2→step 2; PhaseConsensus→3; PhaseFinal→4; PhaseImpl→5; PhaseReview→6,
  or 7 once review/consensus.md exists; step 8 when IMPLEMENTATION.md frontmatter
  status is fix-up-cycle-N or agent.fixup_* events exist; PhaseBlocked renders
  `blocked → reopening round-NN`, never a guess.
- **Delivery merge rule:** events primary, disk fallback. Events-finished+file-missing
  (or the reverse) renders `✓?` "on disk, unvalidated" — never silently trust a bare
  stat on virtio-fs. Phase-6 waiting set = participants minus implementer (mirrors
  consensus.Status(review=true) semantics).
- **Triggers:** (a) an events batch containing an allowlisted protocol type
  (in-memory check, zero I/O); (b) doneMsg; (c) reconcile timer — 15s attached &
  running, 60s done/STALE/detached; (d) switching to the Protocol tab; (e) /refresh
  (manual escape hatch for stale-cache lag). Computed async in a tea.Cmd, results
  gated by runToken exactly like eventsMsg (live.go:312-317). **No work on the
  250 ms event tick.**
- **Resilience:** transient ENOENT/EIO → keep previous snapshot + "reconcile
  retrying"; the phase never regresses unless two consecutive reconciles agree.
- **Cost per refresh:** ~1-2 ReadDir + 3-4 stats + ≤2 ReadFile; consensus.Status
  parsed only when phase ≥ consensus.

## 5. New signals (exactly one event + one flag)

1. **run.phase** — in the driver advance loop immediately after cursor.Save
   (driver/driver.go:164), append `{phase, current_round, impl_status, reason}` to
   events.jsonl. Today, Phases 5-8 transitions are silent file writes; this makes
   them event-driven and demotes disk reconcile to a safety net.
2. **buffers_stdout** — per-agent flag on the roster definition (internal/agents),
   threaded through the run.created runtime payload (runcontrol.go:55-67,
   RuntimeEventData) so reattached runs see it without roster access. Default true
   for agy-style --print agents.
3. Signoff recompute trigger needs no new event: recompute when an agent.finished
   artifact path ends in consensus.md; reconcile timer covers human edits.
4. Checks-gate result (driver impl phase) stays out of scope: narrator line +
   runplan NEXT covers v1; a proper impl.checks event is deferred follow-up.

## 6. Status line grammar

Replace `round=<status>` (live.go:866-868) with a compact phase grammar:
`ph=2:round-02 wait:agy,hermes` · `ph=3:consensus signoff 2/3 wait:hermes` ·
`ph=3:consensus BLOCKED:hermes` · `ph=8:fix-up cycle-2`. Falls back to
displayRoundStatus until the first snapshot lands.

## 7. Heartbeat header

renderAgentStatusHeader (live.go:739-774) gains: proc:live/proc:stale from the
LivenessFunc/procctl seam (process truth, distinct from the 10-min event-age STALE),
stdout/stderr byte counters (existing tailCursor offsets — render numbers we already
track), and the artifact target path.

## 8. Slices (one release)

1. protosnap.go + collapsed ribbon + ph= status line — works manifest-less and on
   STALE runs; highest value alone.
2. Glyphs + narrator + buffered-agent placeholder/hint + heartbeat counters.
3. Protocol tab panes + expanded ribbon + Ctrl+P, /protocol, /waiting, /refresh,
   /narrate.
4. run.phase driver event + buffers_stdout flag + Home phase column.

## 9. Risks & tests

- Half-written artifacts mid-stat → events-primary merge + `✓?` marker.
- Phase flapping on virtio-fs → two-consecutive-agreement + keep-last on errors.
- Narrator spam → strict allowlist; chunk events excluded; verbose opt-in.
- Async snapshot racing run switch/kill → runToken gating; race detector run.
- Tick-budget regression test: ≤4 snapshot computes across 60 simulated 250 ms event
  ticks with no protocol events — turns virtio-fs discipline into an enforced
  invariant. Plus: golden snapshot tests over synthetic idea dirs for all 9 steps;
  merge-rule table tests (finished+missing, present+no-event, killed-then-resegment);
  teatest fixtures for no-run.json, reattached STALE, injected ReadDir ENOENT;
  virtio-fs convergence fixture (round-02/ appears late; assert convergence within
  one reconcile); false-STALE liveness test via the procctl seam.

## Open questions for the other participants

- codex: import direction tui→driver/consensus/runplan — any cycle risk? Is the
  runToken gating sufficient for the async snapshot, or do we need a generation
  counter per idea too?
- agy: is the collapsed-ribbon line the right information order? What would you cut
  from the expanded view to keep 3 lines honest at 100 cols?
- hermes: does the cursor→9-step mapping hold for blocked→reopen and for
  cross_review_rounds>1? Where exactly should fsutil-style resilience wrap reads?
