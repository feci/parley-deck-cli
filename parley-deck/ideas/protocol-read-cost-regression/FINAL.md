---
idea: protocol-read-cost-regression
status: final
drafted-by: claude-1
date: 2026-08-10
participants: [claude-1, codex-1, hermes-1, kimi-1]
consensus: accepted — RESERVED by codex-1, hermes-1, kimi-1; OK by claude-1
track: standard
---

# FINAL — why Parley Deck got slower, and the one fix to build now

## Diagnosis

Not the CLI: every command runs under a second. The felt slowdown is **the cost of a round times the
number of rounds, and both grew**.

```
per call   : reading COOPERATION.md in full costs 3.3x median wall clock (n=3/arm)
per idea   : review rounds 1.6 -> 5.1 (max 24); review bytes 20,237 -> 146,290 (7.2x)
             design rounds 1.4 -> 1.6 (flat)
protocol   : 720 -> 1,359 lines, ~12,300 -> ~26,100 tokens in ten weeks, monotonic
             MUST 15 -> 37, MUST NOT 6 -> 15
```

@codex-1's reconciliation: for a single response, protocol loading dominates; for a whole idea,
review churn dominates. **The per-call tax is paid again inside every extra cycle.**

The review explosion is a step change tracking roster growth (@kimi-1), not gradual protocol growth.
The unbounded `deliberation` fix-up cap enables the tail but does not drive it — the worst ideas kept
finding **fresh MAJORs at rounds 19–24**, so no severity floor would have stopped them.

## The change to build (rank 2 of 5) — no §7 protocol change required

**Two runtime paths embed history quadratically. Both are stricter than the protocol they
implement.**

| Path | What it does today | Where |
| --- | --- | --- |
| `gatherPriorRounds` | concatenates every participant artifact from rounds 1..N-1 | `internal/runner/runner.go:938`, instruction at `:989` |
| `gatherReviewContext` | embeds `FINAL.md` + `IMPLEMENTATION.md` + every prior review round, per reviewer call | `internal/runner/phase58.go:278-299` |

Phase 2 requires only: *address every other active agent explicitly; disagreement requires a
counter-proposal; continue until nobody has new substantive objections.* It never requires re-reading
every historical artifact.

**Not touched:** `buildConsensusDraftPrompt` (`internal/app/driver_consensus.go:113`) keeps its
full-history read, because that is where §15.6's correlated-agreement duty binds.

**Scope note:** the review path is the *larger* target, since the measured 7.2× growth is in review.
It was found by @kimi-1 during signoff and no round named it.

### What each round receives

- **Round 2** — round 1 in full.
- **Round N ≥ 3** — every active participant's round N−1 artifact **in full**, plus the ledger state.
- **Review round N ≥ 2** — the same shape: previous review round in full, plus ledger.
- **Consensus drafter** — unchanged, full history.

### The ledger contract

Owned by @codex-1's and @hermes-1's signoff specifications, not by any phrase in `consensus.md`.
Each item carries: an immutable owner-namespaced **ID**; **kind** (position, objection/counter-
proposal, material claim, verification verdict, user ruling, exemption witness); the **exact scoped
proposition, never a generated paraphrase**; author/claim-owner/verifier identities; introduced and
current source path with stable locator and SHA-256; lifecycle `OPEN|RESOLVED|DEFERRED|SUPERSEDED`;
and an append-only transition history with actor, round, reason and resolution locator.

For material claims: materiality, all claim owners, each verdict's exact
`CONFIRMED|WRONG|UNVERIFIED` state, verifier, `PRIMARY|SECONDARY|RECALL` tag and decisive evidence.
Contradictory verdicts on one claim force **`DISPUTED`**; neither the compiler nor another
participant may invent a resolution. User rulings are carried **verbatim** with `owner: user`.

**An objection stays live until its own owner disposes of it.** Terminal items keep tombstones and
never silently disappear.

### Binding boundary (@codex-1, adopted)

The ledger is an **implementation-scoped context optimization**. It must not become an
artifact-validity or consensus rule — if it does, the change **does** require §7. Any missing,
invalid, ambiguous or challenged ledger state, and any citation or provenance challenge against an
older item, **falls back to full history**, visibly in the prompt. Raw artifacts stay directly
available.

Why this is not belt-and-braces: Phase 2 rule 1 is **"Silence = implicit agreement."** The protocol
converts an omission into consent. A dropped objection is not a lost datum, it is agreement that was
never given.

## Gates

- **G1** Round 2 output is unchanged from today (round 1 in full).
- **G2** Round N≥3 contains round N−1 in full and no round ≤N−2 body, unless a fallback fired.
- **G3** Every fallback trigger — missing, invalid, ambiguous, challenged, unresolved hash, verdict
  conflict not marked `DISPUTED` — selects full history **and says so in the prompt**.
- **G4** `buildConsensusDraftPrompt` still receives full history.
- **G5** Fixture: **orphaned dissent.** A round-1 minority objection not restated in round 2, whose
  ledger entry is omitted or closed by a non-owner, must reach round 3 or trigger fallback. It must
  never silently vanish.
- **G6** Fixture: **claim-ID forking.** The same material claim reworded under a new ID with opposing
  `PRIMARY` verdicts in different rounds must join as `DISPUTED` or trigger fallback.
- **G7** Every fix verified by reverting it and confirming the test goes red, with the revert
  required to compile and to actually apply.

## Deferred, with names

`protocol-phase-scoped-packet` (rank 1 — needs §7). `protocol-fixup-budget` (rank 3 — needs §7).
`round-size-alarms` (rank 4 — alarm, never truncate). `compression-experiment` (rank 5 — confined to
already-resolved analysis; run @hermes-1's 10-rule falsification test first, which predicts failure
on at least 3 of 10).

## Never cut

The authoritative protocol; applicable modals, negations, conditions and exceptions; round-1
independent proposals; every live objection or finding; provenance and verdict-conflict data; the
§15.6 audit; `FINAL.md`, the current diff, acceptance criteria and check results; explicit user
rulings; ownership, no-secrets and the §14 human brake; and on-demand access to raw historical
artifacts.

## How this was decided

Two rounds, four participants, three RESERVED signoffs each independently re-verifying the code
claims rather than inheriting them. The drafter withdrew three positions: that the slowdown was
"read cost, not execution time" (a false separation — context processing *is* execution); that the
unbounded fix-up cap was the driver; and an A/B result reported after a single replicate that two
further replicates reversed.

The drafter framed the original diagnosis and the measurements supporting it, and it did not survive
contact with the other participants. What replaced it is better and was found by attacking it.
