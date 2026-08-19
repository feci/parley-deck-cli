---
idea: cognee-mechanism-mining
status: final
drafted-by: claude-1
date: 2026-08-14
track: standard
participants: [claude-1, codex-1, hermes-1, kimi-1]
rounds: 2
signoffs: [claude-1 ✅, codex-1 🟡, hermes-1 ✅, kimi-1 🟡]
supersedes: n/a
---

# FINAL — cognee lends us nothing; it sent us back to a design we already ratified and switched off

## The answer

**No cognee mechanism is adopted.** Four participants examined roughly a dozen mechanisms across two
independent rounds. Zero survive as imports.

That is not a disappointing result, because of *why*:

> cognee manages a knowledge graph's internal consistency. Parley Deck's measured cost is that its
> normative text is too large to read and its history is re-transmitted every round. **No participant
> found a cognee mechanism addressing either.** Where the two systems do meet, we are ahead: cognee's
> documentation describes no contradiction, supersession, or conflict-resolution mechanism at all —
> the hardest thing a deliberation protocol does, and the thing we specified a year of incidents ago.

The exercise produced value by a different route than intended. Being asked *"what does cognee have
that you don't?"* forced a re-read of our own record, and the answer found there was: **a design we
already ratified, built, and deleted.**

## What we found instead

### The one cognee-shaped thing in this deck is unimplemented, not missing

`protocol-read-cost-regression/FINAL.md:53-100` specifies owner-namespaced IDs, exact scoped
propositions never paraphrased, SHA-256 provenance, `OPEN|RESOLVED|DEFERRED|SUPERSEDED`, append-only
transitions, owner-only objection disposition, terminal tombstones, and **forced `DISPUTED` on
contradictory verdicts**.

That is entity deduplication *with* contradiction handling — strictly more than cognee documents. It
was implemented and deleted in v1.43.1.

Raised by @codex-1 in round 1, adopted by @kimi-1 in round 2, concurred by @hermes-1 and @claude-1.
@codex-1's instruction governs the successor: **do not open a fresh cognee-inspired design.**

### `forget` is four operations, not one

One vendor verb produced four incompatible mappings across four participants — transmission
(@claude-1), lifecycle (@codex-1), artifact versioning (@hermes-1), subtractive maintenance
(@kimi-1). @codex-1 named the category error: underdetermination is evidence *against* an analogy,
not for it. **This FINAL uses Parley operation names. cognee's verb is recorded only as the prompt
that triggered inspection.** Importing it as one mechanism would reproduce the very terminology
fragmentation this idea diagnosed.

### A live defect in our own protocol, found while looking for cognee's

The non-facilitator-artifact rule is stated three times and the statements have diverged:

| Locator | Quantifier |
| --- | --- |
| `COOPERATION.md:78-84` (§1) | plural, unquantified |
| `COOPERATION.md:876` (§9 item 6) | **at least one** |
| `COOPERATION.md:430` (Phase 4) | **every** active; a missing one is "a blocker" |

Verified independently by @kimi-1, @claude-1 and @hermes-1. @hermes-1's reading is adopted over
@claude-1's sharper original: this is a **progressive gate** (start needs ≥1, close needs all), so a
clarity defect rather than a formal contradiction — §1 never declares itself a floor, so a reader can
take "artifacts exist" as sufficient for claiming completion and be surprised at Phase 4.

## Outputs — three successors, none a cognee import

**S1 — Revive the ledger** (successor to `protocol-read-cost-regression`, implementation-scoped).
Fixed window: `00-prompt.md` + round N−1 in full + the running consensus draft, everything older
addressable by path; the **ratified ledger contract verbatim** as the carried-forward list. The only
surviving piece of @claude-1's R1 is its admissibility argument — *a fixed window selects nothing*,
so completeness is readable rather than trusted, which is precisely what the v1.43.1 machinery could
not offer. **Two blocking preconditions:** (a) a v1.43.1 failure analysis drawn from the closed
review record and code history, not from any participant's memory — **nobody in either round could
answer what actually caused that deletion**; (b) a shadow-mode `context.receipt`, proven under
mutation tests at the final serialization boundary, audit-only and **never a gate**.

**S2 — Protocol-text maintenance** (§7 successor, targets net-negative). Every protocol-change
proposal states its net-byte delta **×3 lockstep copies**, or the convention is evadable by editing
one. Consolidation passes only behind a **mechanical before/after rule inventory** — not a reviewer's
reading, since the §13 retro prose-match false-positive is standing proof that reading-based gates
fail here. **First worked example: the §1/§9/Phase-4 divergence above** — demonstrating the norm on a
defect the norm would have caught. The removal-candidates duty is **dropped**, not gated: all four
participants independently rejected non-use as a deletion criterion, because a rare-but-severe rule
(the squash-merge ban) is exactly where non-use proves nothing.

**S3 — deferred, not an output of this idea.** `PRIMARY claims carry their producing command or
path` (§15.2). Its author withdrew it, @codex-1 asked for it to be dropped from this idea, and
@kimi-1 declined to restore it while affirming its evidence — seven consensus BLOCKs this month, none
about code, all traceable to retyped derived text. It may return as an ordinary Parley-native idea.

## Dropped, with reasons

- **Ontology grounding / RDF-OWL / fuzzy 80% matching.** @hermes-1's line, adopted by all: our
  generated-view pattern *prevents* duplicates by construction; fuzzy dedup is post-hoc repair for a
  system that permits them to form. Strictly weaker than ours.
- **Truth-Subspace Reranking / feedback edge weights.** Rejected independently by three participants
  in round 1; @claude-1 conceded plainly in round 2 without a counter-proposal. Decisive arguments:
  correlated past errors become self-reinforcing "truth" and dissent becomes less retrievable
  (@codex-1); permutation makes burial *harder to detect*, not safer — a selector with plausible
  deniability (@kimi-1).
- **NodeSets / Datasets / Principals / Tenants / ACL** — single-tenant deck, no such problem.
- **Temporal Cognify / `TEMPORAL` search** — frozen `FINAL.md` + `<slug>-v2` + changelog + git +
  `DISPUTED` already cover time and supersession.
- **Memory Provenance visualization, ECL, `improve`, graph/hybrid/triplet search** — no measured
  problem; several occupy the forbidden context-selector position.
- **`memify` over playbooks** — `ls parley-deck/playbooks/` → no such directory. Tooling for a corpus
  of zero.
- **Canonical term registry** (@kimi-1's own) and **terminology-drift diagnostic** (@codex-1's own) —
  both withdrawn by their authors; the per-instance authority+generated-view remedy suffices, and a
  registry risks becoming the second authority it treats.
- **`superseded-by:` on a frozen FINAL** (@hermes-1's own, withdrawn). @codex-1's objection is
  decisive: `COOPERATION.md:424` freezes an invalidated FINAL and opens `<slug>-v2`; mutating it
  would violate immutability and create a second direction of authority.

## Unresolved — DISPUTED, carried forward undecided

**Whether the fix-up-budget statements are duplication or a documented cross-reference structure.**
@kimi-1 reads `:229` / `:267` / `:661-673` as one rule stated three times; @hermes-1 reads them as
definition + value + enforcement, noting LE-5's loop budget (steps/wall-clock/cost) is a *different*
budget from fix-up cycles. Both examined the same lines. Per §15.3 a conflict is never resolved by
count, and neither produced evidence decisive over the other. **This is S2's first adjudication task,
ahead of any edit.** The goal-done pair (`:268` / `:675-680`) is in the same state; only @kimi-1
examined it.

Separately agreed: the "printed cap of 2 ran 15 cycles" incident (`skills-cli-install-path`) is an
**enforcement** defect, not a duplication defect, and is already addressed by the ratified
phase-packet/fix-up-budget FINAL. It is not evidence for S2.

## Process failure recorded against this idea

A participant's round-01 artifact was **silently overwritten mid-idea by a second live process of the
same participant**, and the facilitator drafted consensus from the destroyed version. Three
attributions were wrong as a result. All three were caught — by the participant reviewing what had
been attributed to it at signoff, not by any automated check.

The artifact trail showed exactly one file per participant per round at every point. **A silent
overwrite is invisible to every check this deck runs.** Full account and remediation guidance:
`inbox/claude-1-to-all_cognee-mechanism-mining_artifact-race.md`; corrections at the end of
`consensus.md`.

## Provenance

Every cognee claim in this idea is **SECONDARY** — vendor documentation and vendor blog, fetched
2026-08-14, quoted in `00-prompt.md`. None was upgraded to a demonstrated property, and **no decision
here depends on any cognee capability being real.** All three successors rest on PRIMARY evidence from
this repository. If every cognee claim turned out false, this FINAL stands unchanged.
