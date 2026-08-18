---
agent: claude-1
idea: cognee-mechanism-mining
round: 1
date: 2026-08-14
---

## Summary

Seven cognee mechanisms map onto something in this deck. **Two are worth proposing, one is worth
proposing narrowly, three we already have, and one is a naming win only.**

The reframe I want to put on the table: the prior FINAL rejected cognee's *memory store* on the
grounds that we over-remember. That verdict is correct and I do not reopen it — but it points at the
answer to this idea. **The borrowable thing is the operation an over-remembering system lacks.**
cognee's top-level verb list is `remember, recall, forget, improve`. Parley Deck has three of those
four. It has no `forget` at all, and no `improve` that changes anything.

## Proposed approach

Ranked. Each entry: mechanism → the gap it names → what we already have → proposal shape → cost of
being wrong.

### R1. `forget` as a first-class, rule-based operation (propose)

**Mechanism.** `forget` is a top-level operation, peer to `remember` and `recall` (SECONDARY).

**Gap it names.** Parley Deck cannot forget. `gatherPriorRounds` / `gatherReviewContext` re-send
every prior artifact every round (PRIMARY), so context grows quadratically in rounds while review
rounds themselves went 1.6 → 5.1 (PRIMARY). The protocol never required the re-send — it is runner
behaviour nobody ever specified. Every other cost lever we have found is a per-call constant; this
one is the only *super-linear* term in the system.

**What we already have.** Nothing. There is no rule about what a round's context may exclude.

**Proposal shape.** A protocol rule fixing a **bounded, mechanically-defined context window**: round
N receives `00-prompt.md`, round N−1 in full, and the running consensus draft — with all earlier
rounds *addressable by path* but not resent. Successor idea, §7.

**Why this is admissible where the deleted frontier machinery was not**, and this is the whole
argument: the 1.43.1 machinery was rejected because it *selected* context by relevance and could not
prove it never dropped an objection. A fixed window does not select. It is complete over a stated
boundary, the boundary is in the protocol text rather than in a model's judgement, and any
participant can name exactly what is outside it and re-cite it in one line. Completeness becomes a
property you can *read*, not one you have to trust.

**Cost of being wrong.** An objection raised in round 1, never restated, silently dies in round 4 —
"recorded consent that was never given". Mitigation must be part of the proposal, not bolted on: a
carried-forward open-objections list that the window always includes. If that list cannot be made
mechanical, R1 fails and should be dropped.

### R2. `memify` — a consolidation pass, separate from the change pass (propose)

**Mechanism.** cognee splits graph *construction* (`cognify`) from a later *enrichment* pass
(`memify`), whose named variants are "Entity Consolidation — rewrite fragmented entity descriptions"
and "Entity Deduplication — detect near-duplicate nodes and merge into canonical ones" (SECONDARY).

**Gap it names.** Our protocol has only ever been built, never consolidated. 49,428 → 104,804 bytes
in ten weeks, monotonic, +112% (PRIMARY); §4 accreted +19,936 B across six separate versions and is
now 34% of the file (PRIMARY). Every one of those additions passed a §7 review that asked "is this
rule correct?" — none ever asked "does this rule already exist three sections up?"

**What we already have.** §7 changes the protocol; §13 retro proposes improvements advisory-only.
Neither has a merge operation, and §13's v1 has a known prose-match false-positive.

**Proposal shape.** A §7 sub-mode — a **consolidation change** — with an acceptance gate inverted
from a normal change: it must be **byte-neutral or byte-negative**, and it must demonstrate that no
binding rule was removed, only merged. We already own the enforcement pattern: the embedded-default
drift guard, and the parse-don't-diff technique from the fleet sync.

**Cost of being wrong.** A merge quietly drops a binding rule and no test notices, because the rules
are prose. This is the same failure class as the retro prose-match false-positive, so the gate must
be mechanical (rule inventory before/after), not a reviewer's reading.

### R3. Truth-Subspace Reranking — order, never omit (propose, narrowly)

**Mechanism.** "Let finished sessions reshape retrieval ordering by reranking against learned 'truth'
directions"; feedback reportedly updates edge weights (SECONDARY, mechanics undocumented).

**Gap it names.** §13 retro output is advisory and inert — a closed idea changes nothing about what
the next idea sees.

**What we already have.** §13 + §13.5 playbooks. The distillation exists; the effect does not.

**Proposal shape.** Playbooks **reorder** what a phase packet presents first. Explicitly:
**reranking is not filtering.** Everything still ships; only the order changes. A mechanism that can
only permute a set cannot drop a member of it, which keeps it clear of the 1.43.1 standard by
construction rather than by promise.

**Cost of being wrong.** Ordering influences attention, so a bad rank buries a rule without deleting
it. Weaker than omission, but real — and the honest reason to rank this third rather than first.

### R4–R6. Already solved here (report, drop)

- **Ontology grounding / `ontology_valid`.** cognee's "graph fragmentation" — one concept, several
  names, redundant nodes — is our roster two-namespace schism and our "three tables answering one
  question", exactly. Its remedy (canonical reference vocabulary; every node stamped with whether it
  matched) is *our* remedy: one authority plus a generated view, and `parley protocol check` /
  `roster show STATUS` already stamp per-row validity from a closed vocabulary. **Independent
  convergence on a solution we already shipped twice. Confirmation, not an idea.**
- **Memory Provenance** ("tenants, users, agents, datasets, files") tracks *ownership and data-flow*.
  §15 tracks *epistemic tier*. Different axes; ours is the one that matters for a deliberation
  protocol. One narrow borrow is real and cheap — see R7.
- **Datasets / NodeSets / Principals / Tenants / ACL.** Multi-tenant access control. We have one
  user. No gap.

### R7. PRIMARY claims should carry their producer (propose, cheap)

Adjacent to cognee's lineage view rather than borrowed from it, and I flag that honestly. §15.2
tiers a claim but does not require a PRIMARY claim to name the command or path that produced it.
This deck's actual recorded failure mode — seven consensus BLOCKs this month, none about code, all
about retyped derived text — is precisely an untraceable PRIMARY. Cheap in bytes, and it is the one
change here whose need is demonstrated by our own audit trail rather than argued from analogy.

### Rejected

`GRAPH_COMPLETION` / `TRIPLET_COMPLETION` / `HYBRID_COMPLETION` search types, Temporal Cognify,
vector/graph storage backends: all presuppose a knowledge graph of the deck, which is a non-goal, and
all sit in the normative-selection position the standing standard forbids.

**NodeSets as vocabulary for the phase packet** is a naming win, not a mechanism win — with one
design point worth keeping: cognee's grouping is a **tag on the data**, authored once, not a query
evaluated per call. If rank-1 `protocol-phase-scoped-packet` is ever built, that distinction is what
separates an admissible packet (static, auditable, reviewable as text) from the inadmissible one we
already deleted.

## Concerns / open questions

1. **The strongest finding is that cognee mostly confirms us.** Three of seven mechanisms are things
   we already do. I would rather report that plainly than pad the list — and I want the other three
   participants to attack R1–R3 rather than extend them.
2. **R1 does the work.** It is the only proposal touching the measured super-linear term. R2 and R3
   are real but second-order. If only one thing survives this idea, it should be R1.
3. **Every cognee claim here is vendor self-description.** I found no independent benchmark, and the
   docs expose *no* contradiction-resolution or superseding-fact mechanism — which happens to be the
   hardest part of what a deliberation protocol does. We should not assume cognee has solved the
   thing we would most want solved.
4. **Open question for the round:** is a bounded context window with a carried-forward objections
   list actually mechanical, or does "open objection" require a judgement call that reintroduces the
   selection problem through the back door? I think it is mechanical — an objection is open until its
   raiser marks it resolved — but I hold that loosely and it is the crux of R1.

## Risks

- **Analogy inflation.** The cheapest failure here is a tidy mapping that dissolves on contact. R4–R6
  are my guard against my own list: I would rather return three ideas than seven.
- **Byte debt.** R1–R3 and R7 all add protocol text to fix a problem that *is* protocol text. Each
  successor idea must state what it removes. R2 is the only one that is byte-negative by definition.
- **Re-litigation.** R3 sits close to the deleted frontier machinery. If the reranking-vs-filtering
  line cannot be held mechanically, R3 must be dropped rather than argued.
- **Selection bias in this brief.** I chose which cognee pages to fetch, so the brief reflects my
  reading. A participant who thinks a mechanism was missed should say so.
