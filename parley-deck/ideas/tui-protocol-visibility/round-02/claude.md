---
agent: claude
idea: tui-protocol-visibility
round: 2
responding-to: round-01
date: 2026-06-11
---

## Summary

Round-01 converged: every participant ACCEPTs the five-surface direction. This
round folds the corrections into a single resolved spec. The notable deltas vs
my round-01: phase detail moves into an exported driver helper (one disk pass),
snapshot ordering gets a `protoSeq`/single-in-flight gate on top of runToken,
`run.phase` is emitted from every cursor-commit branch only after a successful
Save (and Save errors stop being ignored), glyphs for unvisited tabs come from
an async stat-based growth cache (never View reads), and `/waiting` is cut.

### @codex

ACCEPT all four corrections; they are now part of the spec.

1. **PhaseDetail export.** Agreed — duplicating stat probes in tui was the
   weakest part of my round-01. Spec: `driver.RebuildDetail(ideaDir, maxRounds)
   (PhaseDetail, error)`-shaped helper returning your struct (Cursor +
   HighestReviewRound + ReviewConsensusExists + ImplementationStatus +
   FinalScaffoldReason) in one disk pass; `Rebuild` becomes a thin wrapper so
   existing callers are untouched. This also satisfies hermes's ReviewRound
   disambiguation ask.
2. **protoSeq + single in-flight.** Adopted: at most one snapshot cmd in
   flight; triggers while busy set `protoDirty` and the completion handler
   re-fires once. The msg carries {token, seq}; stale seq dropped.
3. **Separate reconcile tick** (15s/60s) rather than piggybacking the 1s tick;
   the 1s tick only animates the spinner from cached `lastGrowthAt`.
4. **run.phase from every cursor-commit branch.** Adopted, including your
   precondition: phase-changing branches stop discarding `c.Save` errors
   (driver.go:146,164; consensus.go:69; impl.go:68,103,188) and the event is
   emitted only after a successful save. Save failure → error return, no event.
   Event append itself is best-effort (reconcile is the net). Payload as you
   specified, with `source: "driver"`.
5. **Unvisited-tab glyphs.** Adopted your Update-time cache: an async cmd on a
   2s cadence stats the running agents' stdout/stderr (≤2×N stats), updating a
   per-agent growth cache; `lastGrowthAt` from advanceBuffer covers loaded
   buffers for free. View never touches the filesystem.
6. **Re-cap after weave** — adopted; and the existing steer-path append gets
   the same `capTranscriptLines` treatment (latent unbounded growth you found).
7. **Slices** — adopted your reorder: run.phase lands in slice 1; `/waiting`
   cut (the ribbon's waiting list + /protocol cover it); Home column last.

One trim on your trigger allowlist: keep `hitl.question`/`hitl.answered` out of
snapshot triggers (they do not change phase; the questions pane has its own
refresh) but keep them in the narrator allowlist — they are high-signal to a
human watching a transcript. Snapshot triggers and narrator allowlist are two
different sets; round-01 conflated them.

### @agy

ACCEPT the ribbon reorder, "Pending" wording, `[STALE]` prefix + color-flip,
the `?` suffix on disk-fallback delivery counts, reconciled-age surfacing rule
(>30s shows on the collapsed ribbon), the 3-line expanded layout
(Pipeline/Delivery/System), the Home column layout, and weaving narrator
delivery/phase lines into EVERY tab.

Counters:

1. **Glyph portability.** `⧗` (U+29D7) and `⊘` (U+2298) render inconsistently
   in common terminal fonts. Spec keeps your semantics with safer glyphs:
   pending `·` (keeps the existing tab convention, live.go:673), running-active
   = braille spinner, running-silent `◌` (U+25CC, present in the codebase
   already for "killed" headers), delivered `✓`, failed `✗`, killed `†`,
   skipped `–`, STALE `!` styled warn (single-width everywhere). Happy to
   revisit in review if `⧗` proves portable.
2. **Placeholder block without box-drawing.** Your structured content, no
   `┌─┐` frame (boxes break on narrow widths and double-width glyph fallback);
   plain indented lines with the muted style. Content adopted: buffering
   notice, live status + elapsed, stdout/stderr byte counters, last 5 of that
   agent's own event summaries.
3. **Status line.** Adopt `ph=2:xrev-r02 wait=agy,hermes` (your `wait=` form),
   but phase names stay the canonical step names in the ribbon; `xrev` only in
   the compressed status line.

### @hermes

ACCEPT your two extra reads and the merge-rule validation. Resolutions:

1. **ReviewRound disambiguation** — satisfied by codex's PhaseDetail (carries
   HighestReviewRound + ReviewConsensusExists); no separate cursor field.
2. **Extra reads**: 00-prompt `status:` is already read by Rebuild
   (cursor.go:89); review/consensus.md frontmatter (review-cycle) joins the
   PhaseDetail disk pass. Both inside the same 15s/60s cadence — no extra
   trigger.
3. **fsutil wrapping scope** — agreed: resilience wraps only the disk-fallback
   stat/read paths in the snapshot (and PhaseDetail's ReadDir), not the event
   path. Keep-last + two-consecutive-agreement stay as specified.
4. **Participants precedence** (opts → run.created → frontmatter) — confirmed
   as you stated; on disagreement the wider set wins for DISPLAY (a dropped
   agent still shows its delivered artifact) while waiting-set math uses the
   narrower live set. This is the one place display and gating diverge; it is
   deliberate and documented in protosnap.go.

## Resolved spec (delta over 00-prompt seed)

- driver: exported PhaseDetail one-pass helper; Save errors returned, not
  discarded; `run.phase` after successful save from all five commit branches.
- tui: protosnap.go consumes PhaseDetail + in-memory events; protoSeq +
  single-in-flight + protoDirty; separate reconcile tick 15s/60s; stat-based
  growth cache for unvisited tabs (2s cadence, async); narrator via
  appendProtocolEvents with re-cap; placeholder = structured lines, no box;
  glyph set per @agy with portability swaps; ribbon format per @agy; status
  line `ph=N:short wait=...`; Home column in refreshHomeRuns only.
- Slices: (1) PhaseDetail + run.phase + snapshot + collapsed ribbon + ph=
  status line; (2) glyphs + growth cache + narrator + placeholder + heartbeat;
  (3) Protocol tab + expanded ribbon + Ctrl+P + /protocol + /refresh +
  /narrate; (4) buffers_stdout flag + Home phase column. One release.

Ready for consensus. No open blockers from my side.
