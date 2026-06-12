---
idea: tui-protocol-visibility
drafted-by: claude
date: 2026-06-11
---

## Agreed decisions

Source: round-01 + round-02 (all four participants converged; codex round-02
carries the final API shapes).

- **D1 — Five surfaces, transcript stays the hero.** (1) collapsible protocol
  ribbon on every tab; (2) tab activity glyphs; (3) narrator lines woven into
  transcripts; (4) Protocol tab grown out of the existing Status tab;
  (5) Home phase column. All fed from in-memory events + a cached async disk
  snapshot. run.json is never read.
- **D2 — driver.PhaseDetail.** Export `PhaseDetail{Cursor, HighestReviewRound,
  ReviewConsensusExists, ImplementationStatus, FinalScaffoldReason}` and
  `RebuildDetail(ideaDir, maxRounds) (PhaseDetail, error)`; `Rebuild` stays as a
  compatibility wrapper. No `Cursor.ReviewRound` field; no exporting of private
  helpers. Missing artifacts are zero values; unexpected stat/read errors return
  non-nil error with partial detail.
- **D3 — Cursor→9-step mapping** (pure func, golden-tested): step 0 =
  PhaseRound && IdeaStatus==kickoff; 1 = PhaseRound && CurrentRound==1; 2 =
  PhaseRound && CurrentRound>=2; 3 = PhaseConsensus; 4 = PhaseFinal; 5 =
  PhaseImpl; 6 = PhaseReview && !ReviewConsensusExists; 7 = PhaseReview &&
  ReviewConsensusExists; 8 = PhaseDone or ImplementationStatus has prefix
  "fix-up-cycle". Blocked display comes from consensus.Status Triage==blocked
  (Rebuild never returns PhaseBlocked); renders "blocked → reopening round-NN".
- **D4 — run.phase event.** Emitted from every phase-changing Advance branch
  (driver.go:146,164; consensus.go:69,117; impl.go:68,103,130,188,226 — via one
  commit helper) ONLY after a successful c.Save. Save errors are returned (no
  longer discarded) → ActionEscalated. cursor.Save's MkdirAll switches to
  fsutil.MkdirAllResilient so spurious virtio-fs failures do not become new
  escalations (consistent with the 1.21/1.22 hardening). The event append itself
  is best-effort — disk reconcile is the safety net. Payload: idea, run_id,
  action, phase, previous_phase, current_round, round_label, idea_status,
  rounds_run, max_rounds, source:"driver".
- **D5 — ProtocolSnapshot producer** (internal/tui/protosnap.go).
  `BuildProtocolSnapshot(in ProtocolSnapshotInput) (ProtocolSnapshot, error)`
  with explicit value copies {Root, RunID, RunDir, IdeaSlug, IdeaDir,
  Participants, MaxRounds, CrossReviewRounds, Events, Questions, State,
  Previous, Now}. Never calls runstate.LoadRunAt. Uses RebuildDetail +
  in-memory events; consensus.Status (design or review schema per phase) only
  when phase >= consensus; runplan.Plan for the NEXT line. Async tea.Cmd;
  results gated by runToken + protoSeq; at most one in-flight, protoDirty
  coalesces re-triggers. Separate reconcile tick: 15s attached+running, 60s
  done/STALE/detached. Triggers: allowlisted event batch (D6), doneMsg,
  Protocol-tab switch, /refresh. No snapshot work on the 250ms event tick or
  the 1s elapsed tick. On error: keep previous snapshot + reconcile_error
  marker; phase regression only after two consecutive agreeing reconciles.
  Delivery merge: events primary, disk fallback, `?` marker on disagreement.
  Participants precedence opts → run.created payload → 00-prompt frontmatter;
  display uses the wider set, waiting-set math uses the live set. Phase-6+
  waiting sets = participants minus implementer. fsutil-style bounded retries
  wrap only the snapshot's local disk-fallback stat/read helpers.
- **D6 — Snapshot trigger allowlist** (state-changing only): run.created,
  run.phase, run.segment_started, agent.started, agent.finished, agent.failed,
  agent.skipped, agent.killed, agent.fixup_finished, agent.fixup_failed,
  round.completed, round.incomplete, run.failed. ACP/steer/HITL/index events
  excluded.
- **D7 — Narrator allowlist** (a different, display-oriented set): run.created,
  run.phase, run.segment_started, agent.started/finished/failed/skipped/killed,
  round.completed/incomplete, hitl.question, hitl.answered, run.failed,
  run.manifest_deferred. Steer reply markers stay target-tab only. Mechanics:
  append to loaded buffers + a global ring of the last 32 narrator events
  replayed once when an unvisited buffer first loads (per-buffer narrator seq
  prevents duplicates); re-run capTranscriptLines after appends (and fix the
  existing steer-path append that skips the cap). Text via
  runstate.SummarizeEvent on the transcriptEvent stream. /narrate cycles
  off → protocol (default) → verbose (ACP tool calls included).
- **D8 — Glyphs (tab strip, conservative set):** pending ○, running-active
  braille spinner (ASCII fallback *), running-silent ·, delivered ✓ (only with
  artifact_ok), failed ✗, killed x, skipped -, STALE overlay !. Wide glyphs
  (⚠, ⏸) are allowed only in ribbon/placeholder/Home text copy, never in tabs
  (byte-vs-cell truncation risk, app.go:144-152, live.go:612-617). Spinner
  source: agentBuffer.lastGrowthAt set in advanceBuffer for loaded buffers,
  plus an async stat-based growth cache for unvisited tabs (2s cadence,
  <=2 stats per running agent). View never reads the filesystem.
- **D9 — Ribbon.** Collapsed: `◆ Ph 2: Cross-Review (R02) · Delivered 1/3 ·
  Waiting: agy, hermes · Next: consensus  ⌃P`. Degraded: `[STALE]` prefix +
  warn color; `?` suffix on counts in disk-fallback mode; reconciled-age shows
  on the collapsed ribbon once it exceeds 30s. Expanded (Ctrl+P cycles
  collapsed→expanded→hidden), 3 lines: Pipeline (Kick ✓ ── R01 ✓ ── XRev ▶ ──
  Cons ── Final ── Impl ── Revw ── RCon ── Fixp), Delivery (per-agent state +
  time), System (Next + Reconciled Ns ago + source).
- **D10 — Status line.** `round=<status>` is replaced by `ph=2:xrev-r02
  wait=agy,hermes` (compressed phase names only here; falls back to the old
  displayRoundStatus until the first snapshot lands).
- **D11 — Protocol tab.** The Status tab gains PIPELINE (9 steps [✓][▶][ ]),
  DELIVERY (current round matrix), SIGNOFFS (consensus.Status, design or
  review per phase), and NEXT (runplan.Plan) panes above the existing
  agent/event/question panes; /protocol is an alias of /status. The snapshot
  must NOT depend on the review-cycle frontmatter field (naming is not
  normalized across template/validator — deferred follow-up).
- **D12 — buffers_stdout.** `BuffersStdout bool` on agents.Spec/Discovery
  (TOML `buffers_stdout`), default true for agy-style --print agents; included
  in the run.created runtime payload. The silent-agent placeholder uses the
  declared flag, with the heuristic (running && stdout 0B && elapsed>30s) as
  fallback. Placeholder content (frameless, indented): buffering notice, live
  status + elapsed + process liveness, stdout/stderr byte counters, last 5 of
  the agent's own narrator-allowlisted event summaries. No invented tool
  telemetry.
- **D13 — Heartbeat header.** renderAgentStatusHeader gains proc:live /
  proc:stale (LivenessFunc seam), stdout/stderr byte counters (existing tail
  cursor offsets), and the artifact target path.
- **D14 — Home phase column.** Per-idea phase chip + runstate.Attention badge,
  computed only inside refreshHomeRuns (never on a tick).
- **D15 — Commands/keys.** Ctrl+P (ribbon), /protocol, /refresh (manual
  reconcile), /narrate. /waiting is cut.
- **D16 — One release; internal order:** (1) driver PhaseDetail + save-error
  handling + run.phase; (2) snapshot + collapsed ribbon + ph= status line;
  (3) glyphs + growth cache + heartbeat + buffers_stdout; (4) narrator +
  placeholder; (5) Protocol tab + expanded ribbon + commands; (6) Home column.
- **D17 — Tests.** Golden 9-step mapping incl. scaffold FINAL.md and
  fix-up-cycle-N; merge-rule table (finished+missing, present+no-event,
  killed-then-resegment); tick-budget regression (<=4 snapshot computes across
  60 simulated 250ms event ticks without protocol events); stale token/seq
  dropped after activateRun; allowlist trigger tests; View-reads-no-disk test;
  lastGrowthAt growth/rotation; narrator dedup + cap; driver: every
  phase-changing action emits exactly one run.phase after save, await/surface/
  error branches emit none, save failure emits none and returns an error;
  teatest fixtures: no-run.json, reattached STALE, injected ReadDir ENOENT;
  race detector over snapshot vs kill/relaunch.

## Agreed trade-offs

- Unvisited-tab glyphs lag up to ~2s (stat-cache cadence) — accepted to keep
  virtio-fs I/O bounded; loaded buffers update instantly.
- Conservative tab glyph set over expressive Unicode — portability beats looks
  in the tab strip; wide surfaces may use richer copy.
- Narrator replay ring is bounded at 32 — an unvisited buffer may miss older
  protocol lines; acceptable, the Protocol tab is the complete record.
- run.phase event append is best-effort — a lost event is healed by the 15s/60s
  reconcile; in exchange, the driver never lies about an unsaved cursor.

## Open items deferred to implementation

- Exact ribbon truncation behavior at very narrow widths (graceful clipping).
- impl.checks event (checks-gate visibility) — follow-up idea.
- review-cycle frontmatter naming normalization (template writes `cycle`,
  validator requires `outstanding_agreed_fixes` only) — follow-up idea.
- ACP tool-call lines in the silent-agent placeholder — follow-up.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-11
Status: ✅ ACCEPT
Notes: Drafted from converged round-02 positions; D4 folds codex's save-error correction with fsutil hardening so escalations stay meaningful on virtio-fs.

### Signoff: codex — 2026-06-11
Status: ✅ ACCEPT
Notes: Consensus matches my PhaseDetail/snapshot/run.phase/glyph/narrator positions; fsutil save hardening is compatible implementation detail.

### Signoff: agy — 2026-06-11
Status: ✅ ACCEPT
Notes: The converged positions on glyph portability, placeholder layout, status grammar, and ribbon styling are faithfully recorded.

### Signoff: hermes — 2026-06-12
Status: ✅ ACCEPT
Notes: Late append after CLI outage on 2026-06-12; tooling exception recorded in inbox/; round-02 position was ACCEPT with no blockers; shipped in 1.23.0 with all-ACCEPT reviews.
