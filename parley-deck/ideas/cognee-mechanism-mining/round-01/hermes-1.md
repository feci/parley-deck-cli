---
agent: hermes-1
idea: cognee-mechanism-mining
round: 1
date: 2026-08-14
---

## Summary

I examined cognee's pipeline operations (`remember`, `recall`, `forget`, `improve`, `cognify`, `memify`), its ontology-grounding mechanism, its NodeSets/Datasets grouping, its feedback/edge-weight reranking, and its temporal search, against the Parley Deck protocol and runner as they actually exist. The result is a short list: one mechanism names a real gap that is currently unaddressed, two name gaps that are already covered by something stronger, and the rest collapse on inspection.

The one candidate worth proposing is cognee's `forget` operation — not as a memory store, but as the design shape for a protocol-level "obsolete artifact" or "superseded rule" marker. Parley Deck has SUPERSEDED lifecycle states in the proposed objection ledger (PRIMARY: `protocol-read-cost-regression/FINAL.md`, the ledger contract section) but no mechanism for marking a closed idea's FINAL.md or an older protocol version as superseded when a newer idea contradicts it. The protocol has `meta/protocol-changelog.md` (PRIMARY: `parley-deck/meta/protocol-changelog.md`, verified) but that is an append-only log, not a discoverable "this finding is stale" flag. However, even this candidate is marginal: the cost of adding it may not pay for itself given the protocol is already 105,382 bytes (PRIMARY: `wc -c parley-deck/COOPERATION.md`), and the changelog plus git history already serve as a discovery path.

## Proposed approach

### Ranked list of cognee mechanisms examined

**Rank 1 — `forget` as a "superseded finding" marker (marginal candidate)**

- Cognee mechanism: `forget` is a top-level operation that removes or de-prioritizes graph nodes. SECONDARY — vendor docs, not independently measured. The 00-prompt.md brief itself notes the asymmetry: "one of them (`forget`) is precisely the operation an over-remembering system lacks."
- Parley Deck gap it names: When a new idea's FINAL contradicts an older idea's FINAL (e.g. `speedup-tooling-evaluation` overturning assumptions from earlier ideas), there is no machine-discoverable marker on the older artifact saying "this finding was superseded by idea X." The `meta/protocol-changelog.md` logs protocol changes but does not mark individual closed ideas as superseded. A reader encountering an old FINAL.md has no in-band signal that a later idea overruled it.
- Does Parley already have an equivalent? Partially. The objection ledger (specified in `protocol-read-cost-regression/FINAL.md`, not yet shipped) carries `SUPERSEDED` as a lifecycle state for individual ledger items. And `meta/protocol-changelog.md` (PRIMARY: verified, file exists) records protocol-level supersessions. But neither covers the case of one closed idea's design conclusion being overturned by a later idea that did not change protocol text. That gap is real but narrow.
- Concrete shape a proposal would take: A §7 successor idea proposing that any new idea whose FINAL explicitly contradicts a prior idea's FINAL MUST append a `superseded-by:` field to the older idea's FINAL.md frontmatter (one line, no body change). This is not a new artifact — it is a one-field annotation on an existing file, governed by the same §7 process. The annotation is discoverable by `grep`, requires no new tool, and adds ~30 bytes per superseded idea rather than per protocol revision.
- Cost of being wrong: Low. The worst case is a stale `superseded-by:` annotation that should have been removed but wasn't — which is the same failure mode as the existing changelog, just per-artifact. The risk of not doing it is also low: a reader mis-cites an overturned FINAL, but this is already mitigated by git history and the changelog. The honest assessment is that this solves a minor discoverability problem, not a structural one.
- Verdict: **Marginal. Worth raising in round 2 only if another participant independently identifies the same gap. Otherwise drop.**

**Rank 2 — `memify` Entity Deduplication / "two names for one thing" (DROP — Parley already solved this better)**

- Cognee mechanism: `memify` detects near-duplicate nodes and merges them into canonical ones. SECONDARY — vendor self-description.
- Parley Deck gap it names: The deck's recurring "two names for one thing" pathology (roster two-namespace schism, composite-naming fix, "three different tables answering one question").
- Does Parley already have an equivalent? Yes, and it is structurally stronger. The remedy the deck arrived at each time is "one authority + a generated view": `agents.toml` is the roster authority, §2 is a generated view (PRIMARY: `parley-deck/COOPERATION.md:103-108`). `parley roster show` is the single canonical answer. This is not a fuzzy-match dedup pass — it is a single-source-of-truth architecture that makes duplication impossible by construction. Cognee's fuzzy dedup is a post-hoc repair for a system that allows duplicates to form; Parley's generated-view pattern prevents them from forming.
- Verdict: **DROP. Parley's mechanism is better because it prevents the disease rather than treating the symptom.**

**Rank 3 — Ontology grounding / RDF-OWL reference vocabulary (DROP — wrong shape)**

- Cognee mechanism: An optional ontology file supplies a canonical vocabulary; entities are fuzzy-matched against it at 80% cutoff; every node gets an `ontology_valid` flag. SECONDARY.
- Parley Deck gap it names: Protocol vocabulary drift — 37 MUSTs, 15 MUST NOTs, 65 sections (PRIMARY: `grep -c`), across 105,382 bytes. Terms like "canonical," "authority," "consensus" could fragment across sections.
- Does Parley already have an equivalent? The protocol itself is the controlled vocabulary. §15 defines a fixed provenance vocabulary (`PRIMARY`/`SECONDARY`/`RECALL`), §6 defines fixed signoff vocabulary (`✅/🟡/❌`), §4.0.1 defines fixed `LE-N` rule tags, and severity tags are fixed (`CRITICAL/MAJOR/MINOR/NIT`, PRIMARY: `parley-deck/COOPERATION.md:525`). These are not fuzzy-matched — they are normatively defined terms with fixed spellings. An ontology file with fuzzy matching would add a layer of indirection without adding precision. The protocol's terms are already canonical by definition.
- Verdict: **DROP. The analogy survives only at the vocabulary level ("both have terms"). Parley's terms are normatively fixed, not extracted-and-matched.**

**Rank 4 — `improve` / feedback-driven edge-weight reranking (DROP — inapplicable)**

- Cognee mechanism: Finished sessions reshape retrieval ordering by reranking against learned "truth" directions; feedback updates edge weights. SECONDARY.
- Parley Deck gap it names: Nothing actionable. Parley Deck does not have a retrieval system whose ordering could be reranked. The protocol's analogues — §13 retrospective optimization and §13.5 playbooks (PRIMARY: `parley-deck/COOPERATION.md:1159-1191`) — already mine completed ideas for improvement proposals, but they produce advisory input that enters through the normal lifecycle, not weight updates. The deck does not rank or weight its own rules.
- Verdict: **DROP. No retrieval-ordering surface exists to rerank. §13 already covers the advisory mining of past experience.**

**Rank 5 — NodeSets/Datasets/ACL grouping (DROP — no access-control problem)**

- Cognee mechanism: Tag and group data; datasets are the core unit of the permission system with principals/users/tenants/roles/ACL. SECONDARY.
- Parley Deck gap it names: None. Parley Deck has no multi-tenant access control problem. The roster (§2) governs who participates; the directory layout (§3) governs where artifacts live; the transport (§11) governs how they are published. There is no shared memory store that needs partitioning by principal.
- Verdict: **DROP. Solves a problem the deck does not have.**

**Rank 6 — Temporal Cognify / time-aware queries (DROP — no temporal query surface)**

- Cognee mechanism: Extracts events and timestamps; a `TEMPORAL` search type queries time-aware data. SECONDARY.
- Parley Deck gap it names: None that cognee's mechanism addresses. Parley Deck's temporal dimension is git history (every artifact is a committed file with a timestamp) and `meta/protocol-changelog.md` (append-only, dated). There is no query engine that would benefit from a `TEMPORAL` search type. If anything, the deck's problem is the opposite: too much history is re-sent per round (`gatherPriorRounds` concatenates all prior artifacts, PRIMARY: `internal/runner/runner.go:936-966`), not too little temporal querying.
- Verdict: **DROP. The temporal dimension is already handled by git + changelog. Adding a query type implies a query engine that does not exist.**

## Concerns / open questions

1. Is the `superseded-by:` annotation (rank 1) actually worth a §7 successor idea, or is it a tooling-level feature (a `parley supersede <old-idea> <new-idea>` command that writes the frontmatter field) that needs no protocol text change at all? If it is tooling-only, it does not touch COOPERATION.md and avoids the byte-cost constraint entirely. But that makes it an implementation idea, not a design idea, and may be too small to justify the multi-agent overhead.

2. The 00-prompt.md brief says cognee's `forget` is "precisely the operation an over-remembering system lacks." I agree with the framing but question whether Parley Deck's over-remembering problem (quadratic history re-send in `gatherPriorRounds`/`gatherReviewContext`) is actually a *forgetting* problem or a *transmission* problem. The `protocol-read-cost-regression` FINAL already diagnosed it as transmission: the protocol requires "address every other active agent explicitly" (PRIMARY: `parley-deck/COOPERATION.md:350`), not "re-read every historical artifact." The fix is sending less per round, not deleting old artifacts. Cognee's `forget` deletes from a store; Parley Deck's problem is over-transmission from a filesystem. These are different operations.

3. I could not find any mechanism in cognee that addresses the protocol's actual measured cost centre (105,382 bytes, +112% in ten weeks, 3.3× read cost). Cognee's mechanisms are all about managing a knowledge graph's internal consistency; Parley Deck's cost is in the *protocol text itself* being too large to read efficiently. The phase-scoped packet (already proposed, unbuilt) is the direct remedy and has no cognee analogue.

## Risks

- **Manufactured parallel risk.** The strongest finding here is "one marginal candidate, five drops." I have been careful to reject analogies that survive only at the vocabulary level (ontology grounding, temporal queries, NodeSets). The deck's history of being "damaged by plausible-sounding structure" (00-prompt.md) is the direct warning. My rank-1 candidate is deliberately narrow and honestly labeled marginal.

- **Byte-cost risk on rank 1.** If the `superseded-by:` annotation is proposed as a §7 protocol change (adding a rule that supersession MUST be marked), it adds net bytes to COOPERATION.md. The addition would need to justify itself: it adds one normative sentence and one frontmatter field, and it removes nothing. The justification would be "discoverability of overturned conclusions without git archaeology," which is weak. The tooling-only path (a `parley supersede` command, no protocol text change) avoids this entirely but is a small idea.

- **Prior-decision boundary.** The `speedup-tooling-evaluation` FINAL rejected cognee's *memory store* because the deck over-remembers. My rank-1 candidate does not re-open that verdict — it borrows the *shape* of `forget` (a supersession marker), not the store. But it is close enough to the boundary that I want other participants to check whether I have accidentally re-litigated the adoption decision. I do not believe I have: I am not proposing installing cognee, depending on it, or routing context through it. I am proposing a one-field frontmatter annotation inspired by a concept.

- **Secondary-source ceiling.** Every claim about cognee in this file is SECONDARY (vendor docs via 00-prompt.md). I have not upgraded any vendor claim to a demonstrated property. Cognee's `forget` may work differently than described; the ontology grounding may not achieve what the vendor claims. This does not affect my analysis because I am borrowing design *shapes*, not capability claims — but it means I cannot say "cognee's forget works, therefore ours will." The proposal stands or falls on its own protocol-level merit.
