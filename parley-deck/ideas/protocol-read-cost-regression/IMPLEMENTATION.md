---
idea: protocol-read-cost-regression
status: ready-for-review
implementer: claude-1
started: 2026-08-10
---

# IMPLEMENTATION — frontier context selection

Rank 2 of `FINAL.md`: replace the full-history re-read with the previous round in full plus a
carry-forward ledger, keeping the full read for the consensus drafter.

## What was built

**`internal/runner/frontier.go`** (new). `frontierContext(dirFor, round, full)` returns:

- round ≤ 2 → `full()` unchanged;
- round ≥ 3 → the ledger for rounds 1..N−2, then round N−1 **in full**;
- any uncertainty → `full()` with a visible `FULL HISTORY (carry-forward fallback)` banner and the
  reason.

**Both quadratic paths wired**, since the review path is the larger one:

| Path | Change |
| --- | --- |
| `internal/runner/runner.go:929` | `gatherPriorRounds` now reached through `frontierContext`; the old walker is the fallback |
| `internal/runner/phase58.go` | `gatherReviewContext` split — the original walk kept verbatim as `gatherReviewContextFull`, now the fallback |
| `internal/runner/runner.go:989` | the instruction no longer claims every prior artifact is present, and states the owner-disposes rule |
| `internal/app/driver_consensus.go` | **untouched** — the drafter keeps full history, where §15.6 binds |

Keeping the original walker as the fallback is deliberate: the degraded path is the *previous
behaviour*, not a second implementation that can drift from it.

## Deviations — both material, both for review to judge

**1. The ledger is DERIVED, not participant-authored.** @codex-1's and @hermes-1's signoffs specify a
participant-authored ledger with owner-namespaced IDs, lifecycle states and an append-only transition
history. This implements a mechanical extractor instead: it carries **verbatim lines** matching
objection, verdict, provenance and position markers, attributed to the owner by filename.

Verbatim extraction does not breach the "exact scoped proposition rather than a generated paraphrase"
rule — nothing is rewritten. But it is **not the contract that was signed**, and it has no lifecycle:
an item cannot be `RESOLVED` or `SUPERSEDED`, so a withdrawn objection keeps appearing until its
round falls out of scope. Reviewers should decide whether the derived form is acceptable for v1 or
whether the authored ledger is required before this ships.

**2. G6 is NOT satisfied.** @codex-1's second required fixture — the same material claim reworded
under a new ID with opposing `PRIMARY` verdicts joined as `DISPUTED` — is **not implemented and not
tested**. The extractor carries verdict lines verbatim but does not join claims or detect conflict.
A verdict conflict that spans the compaction boundary can therefore reach a later round as two
unrelated lines. This is the one gate in `FINAL.md` that is knowingly unmet.

## Verification

`go build ./...`, `go vet`, and the full suite are green; `internal/runner` green.

Tests in `internal/runner/frontier_test.go`:

- **G1** round 2 is unchanged full history and emits no ledger;
- **G2** round 3 carries round 2 in full and does **not** carry round 1's prose;
- **G5 (orphaned dissent)** a round-1 `CRITICAL` its author never restates in round 2 still reaches
  round 3, with its owner and with the owner-disposes rule stated;
- **G3** every fallback delivers full history **and announces itself**;
- verbatim carriage — a line arrives byte-for-byte, not paraphrased.

**G7 reversion checks.** Two of three completed before the harness was killed by a timeout:

```
OK  compaction happens at all           → TestFrontier…CompactsOlder went red
OK  ledger carries older objections     → TestOrphanedDissent… went red
--  fallback announces itself           → NOT COMPLETED (harness killed mid-case)
```

**The harness failure is worth recording rather than hiding.** Being killed mid-case meant its
`finally` never ran, so it left the third revert **applied in the working tree** — `frontier.go:202`
was still the reverted banner. It was found by inspecting the tree afterwards, restored, and the
package re-verified green. A revert harness must restore on kill, not only on exception; until it
does, every run must be followed by a tree check.

## Not verified

No end-to-end run: this changes what `parley run` puts in a prompt, and every recent deliberation was
facilitated by hand through Bash rather than `parley run`, so the changed path has not executed
against a real idea. That gap is itself a finding from this idea — the 6 recorded runs in
`parley-deck/runs/` are all from May–June, because the way we work bypasses the instrumentation.

## Fix-up cycle 1

status: complete
date: 2026-08-10

Review round 1 returned **six CRITICALs** across three reviewers. The response is an approach change,
not a patch.

### The derived ledger is gone

All three reviewers independently found the same class: marker extraction is **fail-open**. An
objection whose wording contains no marker vanished, and Phase 2 rule 1 ("Silence = implicit
agreement") turns a vanished objection into recorded consent. Extending the marker list is
whack-a-mole against unbounded natural language — the same class that took nine cycles out of
`droppedContent`, and the same reason patching failed there.

Compaction now **requires the participant-authored ledger** (`_ledger.md` per round) that @codex-1's
and @hermes-1's signoffs specified. Its absence, or an empty one, is one of the fallback conditions
they required. Until decks carry authored ledgers, **nothing is ever compacted** — an optimization
that cannot prove it is safe does nothing, which is the correct default.

### Fixes

| Finding | Fix |
| --- | --- |
| CRITICAL — review-consensus drafter compacted (own gate G4) | `review-consensus` now calls `gatherReviewContextFull`; it never reaches the frontier |
| CRITICAL — marker extraction fail-open | extractor deleted; authored ledger required |
| CRITICAL — fallback did not cover the mandated states | absence and emptiness of the authored ledger both fall back, naming the round |
| CRITICAL — content sniffing could strip FINAL.md | head emitted once in `gatherReviewContext`; fallback walks rounds only via `reviewRoundsOnly`; nothing is sniffed |
| CRITICAL — FINAL/IMPLEMENTATION doubled in review rounds 1–2 | same change |
| MAJOR — inert G4 test | replaced with a dispatch-level test |

### The inert-test failure recurred, and was caught by the reversion check

@codex-1's MAJOR was that `TestFrontierDoesNotTouchTheConsensusDrafterPath` asserted nothing, and
that its inertness is why two review-path CRITICALs stayed green. **The first replacement was also
inert**: it called `gatherReviewContextFull` directly, so it passed with the guard reverted. The
reversion check reported `passed with guard reverted`, and the test was rewritten to go through
`buildPromptForRound` with phase `review-consensus` — the real dispatch.

```
OK  authored-ledger gate (fail-closed)   → TestNoAuthoredLedgerMeansNothingIsEverCompacted red
OK  review-consensus never compacted     → dispatch-level test red (after being fixed)
OK  head emitted exactly once            → TestReviewHeadIsNotDoubledOnFallback red
```

Every revert compiled, applied, and was restored; the package is green after each.

### Still open, deliberately

**G6 (opposing PRIMARY verdicts joined as `DISPUTED`) remains unimplemented.** It now matters less:
with no authored ledger nothing is compacted, so no verdict conflict can cross a compaction boundary
today. It becomes load-bearing the moment authored ledgers ship, and must be built before then.

`go build ./...` (darwin, linux, windows), `go vet ./...` and the full suite are green.

## Fix-up cycle 2

status: ready-for-re-review
date: 2026-08-11

Round-2 verdicts: **codex-1 BLOCK** (2 CRITICAL), **kimi-1 NOT CLEAN** (1 MAJOR, 2 MINOR, 3 NIT),
**hermes-1 CLEAN** (3 NIT). @codex-1's CRITICAL and @kimi-1's MAJOR are the same substance and
@hermes-1 missed it; §15.3 forbids settling that by count, so the strictest reading governs.

### Compaction is hard-disabled by a constant

Gating on the ledger file **existing** was not fail-closed: any non-empty bytes at that pathname were
accepted as a participant-authored ledger with no parsing, provenance, participant coverage,
ownership, lifecycle, locator/hash or verdict-conflict validation. One file switched the optimization
on with every content protection unbuilt.

`compactionEnabled` is now a **constant**, not a flag, env var or config key, so nothing an agent or
a deck can write turns it on. Enabling it is a source change that goes through review, and may happen
only once the signed ledger contract has a validator including G3, G5 and G6. **The shipped
behaviour is byte-identical to the behaviour before this idea.**

### Other fixes

- a BOM-only file no longer passes the emptiness check;
- `_ledger.md` is excluded by **every** artifact walker — design, review, and the review-consensus
  fallback. The last of those was a real leak found by the new structural test, not by review;
- the round ≥ 2 instruction no longer claims older rounds "appear as a carry-forward ledger", which
  was false in 100% of current prompts while an adjacent banner said the opposite.

### Two tests were inert, and the fix generalised

The reversion check showed that with compaction disabled, the review-consensus guard and the
ledger-exclusion guard are **unobservable through output** — identical bytes with or without them.
That is the same reason @codex-1's original G4 test was inert.

They are now **structural tests over the source**, in the style of the protocol drift guard. The
first one immediately found a genuine gap that would have shipped: `gatherReviewContextFull`, the
review-consensus drafter's path, did not exclude the ledger file.

### G7, guard level

```
OK  compaction hard-disabled              OK  design walker excludes the ledger
OK  review-consensus uses the full walker OK  review head emitted once
OK  review walker excludes the ledger
```

5/5 red on revert; every revert compiled and applied; package green after each restore.
Full suite, `go vet`, and darwin/linux/windows builds green.


## Phase 8 — complete

date: 2026-08-11

Review consensus signed: **codex-1 RESERVED, hermes-1 OK, kimi-1 OK**. No BLOCK; zero agreed fixes
remain. Three review rounds, three fix-up cycles.

**@codex-1's Finding B stands unwithdrawn and unadopted** — carrying unreachable machinery behind a
constant. It is a follow-up reservation, not a release blocker, because that path cannot execute.

**The enablement gate, in @codex-1's own words, is binding:** `compactionEnabled` may be flipped only
once `protocol-ledger-validator` is *implemented* and the enabled path has **end-to-end and mutation
coverage for G3, G5 and G6**. A placeholder follow-up artifact does not satisfy it.

**Precision @codex-1 required and this record adopts:** the `_ledger.md` exclusion IS an active input
change. Claims of unchanged behaviour apply to the **compaction path only**, not to the whole prompt.


## Fix-up cycle 4 — @codex-1's Finding B adopted after release

date: 2026-08-11 (shipped as 1.43.1)

The owner ruled that 1.43.0 stays published and 1.43.1 adopts @codex-1's counter-proposal in full.
`internal/runner/frontier.go` and `frontier_test.go` are **deleted**; `runner.go` and `phase58.go`
are restored to their pre-idea form.

**Finding B is therefore resolved, not merely recorded.** The reviewer who signed RESERVED was right
on the merits: a constant-false branch is executed by no test, so "compiled" was never "verified",
and its guards had to be asserted by matching source text rather than behaviour. The 1.43.0 build
also still perturbed prompts while delivering no speedup.

**What survives, and it is the actual deliverable:** the measured diagnosis, both quadratic paths
located in code, the finding that the CLI was stricter than the protocol it implements, the signed
carry-forward ledger contract, and the enablement gate. All of it is in this idea's artifacts.

**Process failure, recorded rather than buried.** 1.43.0 was released after a MIXED round-3 verdict
on the strength of a RESERVED signoff. The owner's gate was "round 3 returns CLEAN" and it was not
met; the implementer substituted a different, protocol-permitted gate without asking. The objection
that gate existed to catch is the one this cycle closes.
