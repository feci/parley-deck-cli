---
idea: tui-protocol-visibility
author: user
created: 2026-06-11
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + design synthesis owner; implementer
  codex: Go/Bubble Tea correctness — snapshot computation, event plumbing, import cycles, render perf, tests
  agy: UX — information hierarchy, ribbon/glyph legibility, narrator noise control, degraded-state honesty
  hermes: protocol + filesystem correctness — phase derivation from disk, virtio-fs I/O budget, signoff/frontmatter parsing edge cases
transport: local-dir
cross_review_rounds: 1
status: final
---

## Problem (owner's words, translated)

Two complaints about `parley tui` during live multi-agent runs:

1. **"I only see what codex prints."** Every participant has a tab, but headless
   agents that buffer stdout until exit (agy `--print`) show a dead-looking, empty
   tab for minutes. There is no at-a-glance signal whether the other agents are
   alive, working, finished, or failed — you must switch tabs and guess.
2. **"I can't see the state of the idea — which part of the protocol I'm in."**
   Nothing in the TUI shows the Parley Deck phase (0 kickoff → 1 round-01 →
   2 cross-review → 3 consensus → 4 FINAL → 5 implement → 6 review →
   7 review-consensus → 8 fix-up/complete), the current round, who has delivered
   their artifact, who is pending, or signoff status. The status line only shows
   `round=<running|completed>`; Home shows only the raw idea status string.

## Design seed (facilitator's synthesis — critique it, do not treat it as decided)

Principle: *the transcript stays the hero; protocol awareness wraps around it.*
Five surfaces, all fed from events.jsonl (already streaming) + a cached,
low-frequency disk snapshot — **never** from run.json (best-effort since v1.22.0):

1. **Protocol ribbon** — 1 line under the tab strip on every tab:
   `▸ round-02/2 (cross-review) · delivered 1/3 · waiting: agy, hermes · next: consensus`.
   Ctrl+P expands to 3 lines (pipeline + per-agent delivery + next/reconciled-at) or hides.
2. **Tab activity glyphs** — spinner (output flowing), hourglass (running, silent),
   ✓ (artifact delivered, from `agent.finished artifact_ok=true`), ✗/STALE.
   Driven by a new `agentBuffer.lastGrowthAt` on existing tail cursors — zero new I/O.
3. **Narrator lines** — protocol events woven into transcripts as dim rule-lines
   ("── round-02 started ──", "── codex delivered round-02/codex.md ✓ ──"),
   reusing the shipped steer-weave path (1.20.0). Strict event-type allowlist;
   acp chunk chatter excluded (verbose opt-in via /narrate). For a silent buffered
   agent, the empty-transcript placeholder becomes that agent's own filtered event
   summaries (from in-memory m.events — zero I/O) + an explicit
   "buffers all output until exit; stderr above is live" hint.
4. **Protocol tab** — the existing Status tab grows a pipeline timeline (steps 0–8
   with [✓]/[▶]/[ ]), per-round delivery matrix, signoff matrix
   (consensus.Status), and a NEXT line (runplan.Plan). /protocol and /status land here.
5. **Home phase column** — per-idea phase chip + runstate.Attention badge, computed
   only inside refreshHomeRuns (never on a tick), for cold-open/detached orientation.

Data sources (verified against the code):
- Phase from `driver.Rebuild(ideaDir, maxRounds)` (driver/cursor.go) — disk-authoritative,
  7-state cursor mapped to the 9 display steps by a new pure func; `cross_review_rounds`
  via protocol.ReadFrontmatter. Never reads run.json.
- Live delivery from `agent.started/finished/failed/skipped/killed` events →
  runstate.ProjectEvents; offline/reconcile fallback via missingRoundArtifacts-style
  os.Stat loop (consensus package). Events primary, disk fallback; disagreement shows `✓?`.
- Signoffs via `consensus.Status(root, slug, review)`; NEXT via `runplan.Plan`.
- New `ProtocolSnapshot` producer (internal/tui/protosnap.go), computed async in a
  tea.Cmd, event-triggered + reconcile timer 15s (running) / 60s (done/STALE);
  **no snapshot work on the 250 ms event tick**; /refresh as manual escape hatch;
  "reconciled HH:MM:SS" honesty marker; phase never regresses unless two consecutive
  reconciles agree (virtio-fs stale-cache discipline).
- One new event type: `run.phase` emitted by the driver after cursor.Save so
  Phases 5–8 flips arrive event-driven (today they are silent file writes).
- Optional `buffers_stdout` agent flag threaded through the run.created runtime
  payload so the buffering hint is a declared fact for agy-style agents.
- Per-agent heartbeat in the agent header: state + elapsed + proc:live/proc:stale
  (LivenessFunc/procctl seam) + stdout/stderr byte counters (existing tail cursor offsets).

Status line: replace `round=<status>` with `ph=2:round-02 wait:agy,hermes` grammar
(`ph=3:consensus signoff 2/3 wait:hermes`, `BLOCKED:hermes` when blocked).

Implementation slices (single release):
1. ProtocolSnapshot + collapsed ribbon + `ph=` status line.
2. Tab glyphs + narrator lines + buffered-agent placeholder/hint + heartbeat counters.
3. Protocol tab panes (pipeline/delivery/signoffs/NEXT) + expanded ribbon +
   Ctrl+P, /protocol, /waiting, /refresh, /narrate.
4. `run.phase` driver event + `buffers_stdout` flag + Home phase column.

## Constraints

- The workspace lives on virtio-fs (weak cache coherence) — bounded, low-frequency
  disk I/O only; event-driven first. Include a tick-budget regression test
  (≤4 snapshot computes per 60 simulated event ticks without protocol events).
- Must degrade gracefully: absent run.json, STALE/detached runs, reattach via /open.
- Phase-6 waiting-set correctness: reviewers = participants minus implementer.
- English-only under parley-deck/. One release ships everything (no piecemeal).

## Questions for round-01

1. Is the ProtocolSnapshot producer (inputs, triggers, merge rules, cadence) correct
   and minimal? Any import cycles or races (runToken gating) you'd flag?
2. Is the 7-state cursor → 9-step pipeline mapping sound, incl. fix-up cycles and
   blocked/reopen? What proves step 7 vs 6 vs 8 on disk?
3. Ribbon/glyph/narrator UX: information hierarchy right? What's noise? What's missing
   for the "is agy alive" question?
4. virtio-fs I/O budget: is the 15s/60s reconcile + event-trigger design safe? Where
   would you add resilience (keep-last, two-agree, fsutil)?
5. Anything in the slices you'd cut, reorder, or add to ship as one release?
