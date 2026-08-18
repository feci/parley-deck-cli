---
idea: cognee-mechanism-mining
author: claude-1
created: 2026-08-14
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: final
track: standard
---

## Problem / idea

The owner asks whether **cognee** can help the Parley Deck protocol, **or what ideas we can borrow
from it**.

The first half is already answered and is **not reopened here** (see Constraints). This idea is
scoped to the second half, which no idea has ever asked: **cognee has solved problems structurally
similar to ours. Which of its mechanisms name something Parley Deck is missing, doing worse, or
doing implicitly — and which of those are worth turning into a proposal?**

The unit of output is a **ranked list of borrowable mechanisms**, each with: the cognee mechanism,
the Parley Deck gap it names, whether Parley already has an equivalent (and if so whether cognee's
is better), the concrete shape a proposal would take, and the cost of being wrong.

A mechanism that maps onto something Parley already does well is a valid finding — say so and drop
it. **Returning few or zero adoptable ideas is an acceptable outcome.** Do not manufacture parallels;
an analogy that survives only at the level of vocabulary is noise, and this deck has been bitten by
plausible-sounding structure before.

## Prior decision — binding, do not re-litigate

`ideas/speedup-tooling-evaluation/FINAL.md` (2026-08-11, signed, PRIMARY) already decided
**adoption**:

> **Adopt neither cognee, graphify-as-context-selector, nor omniroute.** […]
> | cognee | LLM-mediated memory: ingest, recall, injected context | none | *"solves a problem we do
> not have"* — this deck **over**-remembers (quadratic history re-send); it does not under-remember |

and set the standard that constrains any successor:

> Any tool that decides what an agent sees occupies the exact position of the frontier machinery
> deleted in 1.43.1, which was removed because it could not prove it never drops a participant
> objection. Under Phase 2 — *"Silence = implicit agreement"* — a dropped objection is not a lost
> datum; it is recorded consent that was never given.

**Both hold.** Do not propose installing cognee, adding it as a dependency, or routing normative
context through it. This idea mines *design ideas*, not software.

Note the asymmetry that makes this idea non-redundant: the prior verdict rejected cognee's **memory
store** because we over-remember. It says nothing about cognee's **other operations** — and one of
them (`forget`) is precisely the operation an over-remembering system lacks.

## Factual brief on cognee — shared context

Provenance: **SECONDARY** throughout — vendor documentation and vendor blog, fetched 2026-08-14 by
claude-1 from `github.com/topoteretes/cognee`, `docs.cognee.ai/llms-core.md`, and
`cognee.ai/blog/deep-dives/grounding-ai-memory`. These are the vendor's own claims about its
software; none is independently measured. Treat capability claims as **asserted, not demonstrated**.
Apache-2.0; ~30k GitHub stars.

**Pipeline — ECL (Extract, Cognify, Load).** Top-level operations: `remember`, `recall`, `forget`,
`improve`. `remember` "runs add + cognify + improve".

**`cognify`** — builds the graph: classifies documents, extracts entities and relationships via an
LLM, generates summaries, and commits graph edges plus vector embeddings.

**`memify`** — *enriches an existing graph with derived knowledge*, as a separate pass after the
graph exists. Named variants:
- *Entity Consolidation*: "Rewrite fragmented entity descriptions using LLM analysis"
- *Entity Deduplication*: "Detect near-duplicate nodes and merge into canonical ones"
- *Triplet Embeddings*, *Session Persistence*

**Ontology grounding.** The stated problem is **"graph fragmentation"**: the same concept extracted
in different linguistic forms — "car manufacturer", "automobile maker", "vehicle producer" — becomes
redundant nodes, and "the result: a fragmented, redundant knowledge base that degrades retrieval
quality and breaks cross-document reasoning." The remedy is an **optional RDF/OWL file supplied as a
reference vocabulary**, defining classes, individuals, object properties, and `rdfs:subClassOf`
hierarchies. The pipeline **fuzzy-matches extracted entities against canonical ontology terms at an
80% cutoff**, does a **BFS traversal to pull in surrounding ontology structure**, and stamps **every
node with an `ontology_valid` flag** recording whether it matched the schema. The vendor claims
consistency, enrichment and control; it makes **no quantified claim** about hallucination or
retrieval precision.

**Feedback / self-improvement.** `improve` "enriches an existing graph and bridges session memory
into it." There is a documented **feedback** flow, and **"Truth-Subspace Reranking": "Let finished
sessions reshape retrieval ordering by reranking against learned 'truth' directions."** One
secondary source states that when an agent rates a response, that feedback updates **edge weights**.
Internal weighting mechanics are not exposed in the docs.

**Grouping and access.** **NodeSets** — "Tag and group data in Cognee with NodeSets." **Datasets** —
"Organize documents, permissions, and processing with datasets", "the core unit of data in Cognee's
permission system." Principals / Users / Tenants / Roles / ACL.

**Memory Provenance.** "Visualize the ownership and data-flow story behind your memory — tenants,
users, agents, datasets, and files."

**Time.** "Time Awareness" with a "temporal mode for time-aware queries"; a *Temporal Cognify*
variant "extracts events and timestamps", queried with a `TEMPORAL` search type.

**Search types.** `GRAPH_COMPLETION`, `HYBRID_COMPLETION`, `TRIPLET_COMPLETION`, `TEMPORAL`.

**Documented gaps.** The docs expose **no mechanism for contradiction detection, conflict
resolution, or superseded facts**, and no versioning/temporal-branching mechanism beyond the search
type name. Do not assume these exist.

## Parley Deck context a participant should weigh

All PRIMARY, from this repository:

- The protocol is **104,804 bytes / 1,363 lines**, up from 49,428 ten weeks ago (+112%). §4 alone is
  **34%** of it. Four top-level sections (§12–§15) did not exist before skill 1.4.1.
- Reading `COOPERATION.md` in full costs **3.3× median wall clock** per call (n=3/arm).
- Review rounds per idea went **1.6 → 5.1** (max 24); design rounds stayed flat at 1.4 → 1.6.
- `gatherPriorRounds` / `gatherReviewContext` **re-send all prior artifacts every round**; the
  protocol never required that.
- Two ranked interventions remain **unbuilt**: rank 1 `protocol-phase-scoped-packet` (send only the
  sections a phase needs — must be built in the *instruction* layer, because the Go runner never
  reads `COOPERATION.md` at all) and, now shipped in 1.44.0, the fix-up budget.
- The deck has repeatedly produced **two names for one thing**: the roster two-namespace schism, the
  composite-naming fix, "three different tables answering one question". The standing remedy has
  each time been *one authority + a generated view*.
- §15 already defines claim provenance — `PRIMARY` / `SECONDARY` / `RECALL`, untagged reads as
  `RECALL` — and §13 already defines advisory-only retrospective optimization with §13.5 playbooks.

## Constraints

- **No new dependency, no new service, no tool in the normative path.** The output is protocol or
  tooling design, not an install.
- Anything that would change `COOPERATION.md` is a **§7 meta-protocol-change** and must be proposed
  as a *successor idea*, not written here.
- Every claim about cognee is **SECONDARY at best** (vendor self-description). Do not upgrade a
  vendor claim to a demonstrated property. Every claim about Parley Deck must be `PRIMARY` with the
  command or file that proves it, or tagged `RECALL`.
- The protocol is already the measured cost centre. **A proposal that adds net bytes to
  `COOPERATION.md` must say what it removes or why the addition pays for itself.**
- Do not re-derive the read-cost measurements; they are given above.

## Non-goals

- Installing, vendoring, or depending on cognee, or any successor evaluation of it as a tool.
- Building a knowledge graph of the deck.
- Editing `COOPERATION.md` in this idea.
- Re-opening the `speedup-tooling-evaluation` verdict.
