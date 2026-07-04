---
agent: claude-1
idea: tui-round-summary
round: 1
date: 2026-07-04
---

## Summary

A pure-presentation feature: a deterministic, mechanical digest of a completed round
rendered in the Home tab. No LLM, no new protocol artifacts, read-only over the
canonical round files. Derive it from the same round-completion signal the driver
already emits, and render through the existing Home/narrator infrastructure from the
protocol-visibility work.

## Proposed approach

### 1. Data model — `RoundDigest`
```go
type AgentLine struct {
    Agent    string
    Position string   // extracted, capped ~120 chars
    Stance   string   // ACCEPT | BLOCK | COUNTER | "" (unknown)
}
type RoundDigest struct {
    Idea       string
    Round      int
    Lines      []AgentLine
    Accept     int
    Counter    int
    Block      int
    NextAction string   // from driver state: "round-03 opening" | "consensus drafting"
}
```

### 2. Extraction (mechanical, robust to imperfect files)
- Position: from a round file, take the `## Summary` section's first non-empty
  paragraph; fall back to the first prose paragraph after frontmatter; cap length.
- Stance: scan for signoff/status markers already used in the protocol
  (`Status: ✅ ACCEPT` / `🟡 ACCEPT-WITH-RESERVATIONS` / `❌ BLOCK`) and
  `responding-to` counter-proposal cues; unknown ⇒ blank, never guessed.
- All parsing tolerant: missing sections ⇒ fallback, never a crash or empty tab.

### 3. Trigger
- Hook the existing driver round-completion detection (the same "all expected
  round-NN files exist" check that opens the next round). On that event, build the
  digest and push it to the Home model via the existing event/narrator channel.
- No polling in the view; the digest is state pushed on the event, matching the
  current run.phase event pattern.

### 4. Render (Home tab)
- A "Round digest" panel below the existing chips/roster/runs — MUST NOT displace
  them. Show the latest digest expanded; keep the last N (e.g. 3) collapsed/scrollable
  for catch-up.
- Convergence line: `3 ACCEPT · 1 COUNTER · 0 BLOCK → consensus drafting`.

### 5. Tests
- Table tests for the extractor over: well-formed file, missing Summary, no stance
  markers, unicode/long lines. Golden-render test for the panel string. All offline.

Footprint: a new `internal/tui/roundsummary.go` (extractor + model) + a render hook in
the Home view + driver event wiring + tests. Standard track, ~2 reviewers.

## Concerns / open questions

1. Does the driver already surface a round-complete event to the TUI, or only
   run.phase? If only phase, add a thin `RoundComplete{idea, round}` event rather than
   re-deriving completion in the view (keep detection in the driver, single source).
2. Home tab real estate: with chips + roster + runs + digest, small terminals get
   cramped. Proposal: digest panel is collapsible and scrolls within its own region.
3. Stance detection depends on where signoff markers live (round files vs consensus).
   In cross-review rounds, ACCEPT/BLOCK may only appear at consensus; for plain rounds
   the digest may legitimately show blank stances (that is honest, not a bug).

## Risks

- Coupling to file format: if round-file section headers drift, extraction degrades —
  mitigated by fallback-to-first-paragraph and "unknown" stance, never a hard failure.
- Home tab regression — mitigated by additive panel + explicit test that existing
  content still renders.
- Scope creep to LLM summaries — explicit non-goal; v1 is mechanical only.
