---
idea: speedup-tooling-evaluation
author: claude-1
created: 2026-08-11
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: final
track: standard
---

## The task

The owner reports Parley Deck is unbearably slow and asks whether **cognee**, **graphify** or
**omniroute** — or anything comparable — can make it faster and more efficient. Evaluate them in
detail against the measured bottleneck, and recommend.

## The measured bottleneck — already established, do not re-derive

From the closed idea `protocol-read-cost-regression` (2 design rounds, 3 review rounds, signed
consensus). All PRIMARY, reproducible from this repository:

```
NOT the CLI: --version 0.00s, roster show 0.29s, preflight 0.07s, skill status 0.91s

per call : reading COOPERATION.md in full costs 3.3x median wall clock
           (n=3/arm, same agent, same question, same output length)
per idea : review rounds 1.6 -> 5.1 (max 24), review bytes 20,237 -> 146,290 (7.2x)
           design rounds 1.4 -> 1.6 (FLAT)
protocol : 720 -> 1,359 lines in ten weeks, monotonic; MUST 22 -> 37 at §15 (skill 2.3.0, 2026-08-05)
```

**The reconciliation that matters:** for a single response, protocol loading dominates; for a whole
idea, review churn dominates. The per-call tax is paid again inside every extra cycle.

**Two ranked interventions are already named and unbuilt:**
- rank 1 `protocol-phase-scoped-packet` — send only the sections a phase needs. Needs §7.
- rank 3 `protocol-fixup-budget` — `deliberation` currently has an UNBOUNDED fix-up cap, which is
  why 24-round reviews exist. Needs §7.

Also established: `gatherPriorRounds` and `gatherReviewContext` re-send all prior artifacts, and the
protocol never required that. A frontier optimization for this was built, reviewed three times, and
**deleted in 1.43.1** because it could silently drop a participant objection, and Phase 2 rule 1
reads "Silence = implicit agreement" — a dropped objection becomes recorded consent.

## What you must evaluate

**Availability, checked before opening this idea (PRIMARY):**

```
graphify   /Users/tomasfecko/.local/bin/graphify     INSTALLED (v0.8.41)
cognee     not installed
omniroute  not installed
```

`graphify` already has a built graph for this repository under `graphify-out/`. Its commands include
`path "A" "B"`, `explain "X"`, `diagnose multigraph`, and per `CLAUDE.md` also `query`, `update`.

For **cognee** and **omniroute** you cannot test locally. Research them (web, docs) and say plainly
which claims you verified and which are vendor copy. Tag provenance per §15: a claim taken from a
project's own README is `SECONDARY` at best.

Answer:

1. **Does each tool address the MEASURED bottleneck, or a different one?** Be specific: retrieval
   quality, token volume, round count, latency, or none of them. A tool that improves retrieval does
   not reduce review round count.
2. **graphify specifically.** It is installed and already indexes this repo. Could a graph of the
   PROTOCOL let an agent load only what a phase needs — i.e. is it an implementation of rank 1? What
   breaks if a normative rule is reached by graph traversal rather than read in full? Remember §15
   and §7 are load-bearing and were each bought with a real failure.
3. **cognee / omniroute.** What do they actually do, what would integration cost, and what would we
   be trusting them with? Note that a memory or routing layer that decides what an agent sees is
   exactly the position the deleted frontier code occupied — and it was deleted for being unable to
   prove it never drops an objection. Apply the same standard.
4. **The honest null answer.** Is the fastest real win simply rank 1 + rank 3 with no new dependency?
   Say so if you think so. Adding a tool that does not touch the measured terms is a cost, not a win.
5. **Rank every option** by (expected saving) / (risk of losing a rule that catches real errors),
   and name what must never be delegated to a tool.

## Constraints

- Read-only: no edits, no git write commands. You MAY run `graphify` read-only commands and web
  research.
- English only. Redact obvious secrets.
- §15 provenance: `PRIMARY` needs a stable locator or quoted command output; untagged reads as
  `RECALL`. Vendor documentation is `SECONDARY` — say so.
- Do not recommend a tool because it is interesting. Recommend it because it moves one of the
  measured numbers above, or say it does not.
