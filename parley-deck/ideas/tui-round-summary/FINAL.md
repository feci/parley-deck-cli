---
idea: tui-round-summary
status: final
drafter: claude-1
track: standard
date: 2026-07-04
participants: [claude-1, codex-1, hermes-1]
---

## Decision

A consolidated round digest in the TUI Home tab, built in the driver/event layer (the
only holder of the canonical round-completion predicate) and rendered as a pure view
over snapshot state. Stance is shown as keyword *flags/hints*, never verdicts. No LLM,
no protocol change, read-only over canonical files.

## Design (as ratified in consensus.md)

1. **Trigger** — after the driver's `roundComplete(round)` predicate passes, build the
   digest for that `(idea, round)`. Home NEVER decides completion from file presence.

2. **Event** — emit a presentation event `round.digest` keyed by `(idea, round)`,
   appended to the run event stream (durable, replayable). Emission is idempotent: a
   stable per-(idea,round) key; re-running the driver never duplicates.

3. **Builder** (`internal/driver`, deterministic, per-agent-degrading — never blocks
   advancement): for each already-validated round file, extract the first sentence/
   paragraph under `## Summary` (cap ~120 chars at a sentence boundary); missing
   `## Summary` ⇒ first prose paragraph tagged `[no Summary — fell back]`. Compute:
   - header `Round NN — complete (X/Y)`;
   - `flags:` raw keyword counts over {`block`/`blocker`, `counter-proposal`,
     `accept`/`agree`, `escalat`} — labeled as mentions/hints, never verdicts;
   - `next:` from driver phase state (opening round N+1 / consensus drafting / re-round);
   - round 2+: engagement gaps from the structured `responding-to:` frontmatter.

4. **Render** — extend `ProtocolSnapshot` with recent round digests; `renderHome` is a
   pure renderer (no disk scan, no markdown parse per render). History is bounded by a
   viewport budget (scrollable sub-region), NOT a fixed count, so chips/roster/runs are
   never pushed off-screen.

## Verification (done criteria)

- `go test ./internal/driver/... ./internal/tui/...` green, including: digest builder
  over well-formed file / missing-Summary fallback / no-stance-markers / unicode-long
  lines; idempotent emission (no duplicate on re-run); flags are counts not verdicts;
  Home still renders chips/roster/runs (regression guard); extraction failure for one
  agent degrades to a tagged line and does not error the driver.
- `go build ./...`, `go vet`, `gofmt -l` clean.

## Non-goals

No stance verdicts (flags only); no cross-round analytics; no LLM summaries; no change
to per-agent tabs/transcript.

## Signoffs

<!-- each participant appends its own block -->
