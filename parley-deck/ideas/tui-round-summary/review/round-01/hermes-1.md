---
agent: hermes-1
idea: tui-round-summary
review-round: 1
date: 2026-07-04
reviewed-commit: 947effa
---

## Summary

The implementation faithfully realizes the ratified design. The digest is built in the
driver (the canonical completion-gate holder) and carried on an idempotent `round.digest`
event; Home is a pure renderer over consumed events. Stance is presented as keyword
*mentions*, never verdicts — the label is "mentions:" and the test
`TestRenderRoundDigestBoundedAndLabelsHints` asserts the word is present. Degraded
extraction (no `## Summary`) is tagged `[no Summary — fell back]` on the agent row, and
`AgentLine.Fell` carries the signal structurally. Completeness header `complete (X/Y)` is
present. The two deferrals — last-N scrollable history and round-2+ engagement-gap
extraction from `responding-to:` — are acceptable for v1: they are additive enrichments,
and shipping the latest-digest-only block with a hard viewport cap (rows/3, ≤10) is the
correct response to the viewport-regression risk I raised in round-01. Tests pass; `go
vet` and `gofmt -l` are clean. The information design is sound and the regression guard
holds. I found no CRITICAL or MAJOR issues. There are two MINOR findings on tag
visibility and a stale next-action edge case, plus a few NITs.

## Findings

### MINOR-1 — Degraded-extraction tag can be truncated away, hiding the honesty signal

digest.go:62 builds `ln.Position` capped at 120 chars. roundsummary.go:86 then appends
`"  [no Summary — fell back]"` to that position BEFORE truncating the combined string to
`width-18`. When the position is long (≥~50 chars on an 80-col terminal, ≥~80 on 120-col)
the tag is cut off and the user sees an untagged one-liner — the exact
"misleading one-liner with no quality tag" failure I flagged as a risk in my round-01
proposal (hermes-1.md:114). The `Fell` flag is in the struct and could drive a separate
visual marker (e.g. a leading glyph or a fixed-width tag prefix) that truncation cannot
erase. As-is the honesty signal is only visible when the position is already short — the
case that needs it least.

### MINOR-2 — `next:` says "drafting consensus" even when the consensus gate is unwired

driver.go:262 sets `nextAction = "drafting consensus"` whenever
`c.CurrentRound >= 1+d.cfg.CrossReviewRounds`, which is correct for the wired path. But
when `d.cfg.Consensus == nil` (slice-1 stop), the driver returns `ActionConsensus` and
halts — it does not draft anything; the loop prints "next step is `parley consensus
draft` (consensus auto-drive not wired)". The digest's `next:` line will nonetheless say
"drafting consensus", implying an action the driver is not taking. The user on Home sees a
misleading next-step. Minor because the unwired-gate configuration is the legacy/slice-1
path and the wired path is the default, but the label should reflect what the driver will
actually do (e.g. "consensus-ready (manual draft)" when the gate is nil).

### NIT-1 — `mentions:` line lists four flags but not which agent contributed them

roundsummary.go:65 renders aggregate counts ("2 block · 1 counter-proposal · 0 accept · 0
escalate"). This is the ratified design (aggregate counts on a trailer line), so no change
is required. Noting only that per-agent attribution would make drill-in decisions easier
and was discussed as a follow-up; the current aggregate is fine for v1's "index, not
summary" scope.

### NIT-2 — `digestView` duplicates the `RoundDigest` struct shape

roundsummary.go:16-32 mirrors driver.RoundDigest field-for-field with a comment explaining
the avoidance of a tui→driver import edge. This is a reasonable decoupling choice and the
IMPLEMENTATION.md documents the `m.events` deviation. The duplication is a maintenance
hazard if fields drift, but a test would catch a decode mismatch. Acceptable as-is.

### NIT-3 — Round-01 keyword false-positive on negation is unmitigated beyond the label

digest.go:115 counts every "block" substring including "this is *not* a blocker". The
"mentions:" framing (MINOR-1's sibling concern) is the agreed mitigation, and it is
present. No action needed; recording that the known false-positive surface is carried into
v1 by design.

### Verified (information-design checklist)

- Stance flags framed as HINTS/mentions, never verdicts: YES. Label is "mentions:"; no
  ACCEPT/counter/block verdict words appear in the renderer; test asserts "mentions:".
- Degraded extraction visibly tagged: YES, structurally (`Fell` flag) and textually
  (`[no Summary — fell back]`) — subject to the truncation gap in MINOR-1.
- Completeness header X/Y present: YES. `Round %02d digest — complete (%d/%d)`; test
  asserts `(3/3)`.
- Deferral of last-N history + engagement-gaps acceptable for v1: YES. Latest-digest-only
  with hard viewport cap (rows/3, ≤10) directly addresses the viewport-regression risk I
  raised. Both deferrals are additive and documented in IMPLEMENTATION.md.
- Idempotent emission: YES. emitRoundDigest scans prior events for a matching (idea,
  round) key and no-ops; no duplicate on re-run.
- Never blocks advancement: YES. All failures in emitRoundDigest are swallowed; missing
  artifacts yield not-present lines, not errors.
- Pure renderer (no disk scan / markdown parse on render path): YES.
  latestRoundDigest decodes the newest `round.digest` from `m.events`.
- Tests: `go test ./internal/driver ./internal/tui` green (driver 0.707s, tui cached).
  `go vet` clean, `gofmt -l` clean on all three implementation files.

## Open questions

- MINOR-1: should the `Fell` flag drive a truncation-proof marker (leading glyph or
  fixed-width prefix) so degraded extraction is always visible regardless of position
  length? I believe yes — it is the honesty signal I argued for in round-01.

- MINOR-2: is the unwired-consensus-gate path (`Consensus == nil`) reachable in the
  deployments that will show this digest, or is it purely legacy? If reachable, the
  `next:` label should distinguish "drafting consensus" (auto) from "consensus-ready,
  manual draft" (halted).

- Should the overflow note (`… N more (open the agent tabs)`) be tested? It is reachable
  when agentBudget is exhausted but has no dedicated test case; the bounded-rows test uses
  a budget of 4 with 3 agents, which does not trigger overflow.
