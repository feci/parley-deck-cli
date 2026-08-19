---
agent: kimi-1
idea: cognee-mechanism-mining
round: 1
date: 2026-08-14
---

## Summary

I mined the cognee brief for mechanisms, checked each against what Parley Deck already does, and
judge **two worth a successor proposal**, ranked:

1. **`forget` as a first-class operation → a protocol sunset/consolidation duty inside the existing
   §13 retro pass.** This is the only cognee mechanism that names a gap Parley has *measured* and
   has *no* remedy for: the protocol is accretion-only (104,804 bytes, +112% in ten weeks — PRIMARY,
   per `00-prompt.md`). Byte-negative by construction.
2. **Ontology grounding → a canonical term registry with a generated conformance check.** The
   two-names-for-one-thing failure is recurring and documented (PRIMARY, `00-prompt.md`: roster
   schism, composite-naming fix, "three different tables answering one question"). Real but
   second-ranked, because the remedy risks reproducing the disease it treats.

Everything else in cognee either maps onto machinery Parley already runs (and I say which section
proves it), is already rejected by the binding `speedup-tooling-evaluation` principle, or is an
analogy that survives only at the level of vocabulary. Full drop-list with reasons is in
*Proposed approach*. I explicitly do **not** propose: installing anything, routing context through
any selector, or touching `COOPERATION.md` in this idea — both proposals below are *successor-idea
shapes*, per §7.

## Proposed approach

### Adoptable, ranked

#### Rank 1 — Borrow `forget`: protocol sunset as a standing retro duty

- **Cognee mechanism (SECONDARY, vendor docs):** `forget` is a top-level operation alongside
  `remember`/`recall`/`improve`; `memify` additionally *rewrites and merges* fragmented store
  content in place ("Entity Consolidation", "Entity Deduplication") rather than only appending.
  Capability asserted, not demonstrated.
- **Parley Deck gap it names (PRIMARY):** the protocol only grows. 104,804 bytes / 1,363 lines, up
  from 49,428 ten weeks ago; §12–§15 did not exist before skill 1.4.1; §13's own footer shows pure
  accretion ("ratified by…, amended by…, extended by…" — COOPERATION.md §13, line ~1191). Reading
  the file costs 3.3× median wall clock. Nothing anywhere in §0–§15 or the ratification footers
  describes removing or consolidating protocol text. The existing remedies are all *recall-side*:
  the Quickstart ("You do not need to read all of this", COOPERATION.md line ~14) and the unbuilt
  rank-1 `protocol-phase-scoped-packet` scope what an agent *reads*; nothing scopes what the deck
  *keeps*. The brief states this asymmetry itself: the prior verdict killed cognee's memory store
  because we over-remember, and `forget` is precisely the operation an over-remembering system
  lacks.
- **Existing equivalent?** Partial. §13 retro passes (COOPERATION.md §13.1–13.3) may already
  *propose* any improvement through the normal gate, so removal is not forbidden — but it is not
  anyone's duty, nothing measures non-use, and ten weeks of +112% show the optional channel does
  not fire. Cognee's is better in exactly one respect: deletion is a named, first-class operation,
  not an emergent possibility.
- **Proposal shape (successor §7 meta-protocol-change idea, NOT written here):** two small edits,
  byte-accounted. (a) In §13.1 or §13.4 add one duty: every retro pass MUST output a *removal
  candidates* list — protocol passages with evidence of non-use from the evidence corpus (sections
  never cited by any round/consensus/FINAL in the corpus window), each candidate entering the
  normal lifecycle like any retro proposal. (b) In §7 add one convention: every protocol-change
  proposal states its net-byte delta, and net-positive proposals name what they displace. Stated
  cost: ~300–400 bytes added. What it removes: unbounded; the mechanism exists to delete
  kilobytes, and it pays for itself the first time it removes one dead subsection. The `memify`
  insight folds in here: consolidation rewrites *in place* (merge §11.B/§11.C deltas, collapse
  superseded amendment footers) rather than appending §16.
- **Cost of being wrong:** deleting load-bearing text. Bounded by construction — every removal
  passes the same multi-agent consensus + signoff + human approval gate as any §7 change, and git
  history makes every deletion reversible. The realistic failure mode is milder: the duty fires
  and finds nothing, costing one retro-pass section of output. Acceptable.

#### Rank 2 — Borrow ontology grounding: canonical term registry + generated check

- **Cognee mechanism (SECONDARY, vendor docs):** an optional RDF/OWL reference vocabulary;
  extracted entities are fuzzy-matched against canonical terms (80% cutoff), surrounding structure
  is pulled in by BFS, and every node is stamped `ontology_valid`. The vendor claims consistency
  gains; it makes **no quantified claim** — asserted, not demonstrated.
- **Parley Deck gap it names (PRIMARY, `00-prompt.md`):** the deck has repeatedly produced two
  names for one thing — roster two-namespace schism, composite-naming fix, "three different tables
  answering one question". The standing remedy each time was *one authority + a generated view*,
  e.g. Appendix A step 4: `agents.toml` is authoritative, the §2 table is generated, "Do NOT
  hand-edit" (COOPERATION.md line ~1094).
- **Existing equivalent?** Yes, per-instance, and it works — but it is *reactive*: each schism was
  fixed only after it damaged a round. What no instance provides is a general registry of protocol
  vocabulary (track names, artifact kinds, phase names, verdict labels, frontmatter keys) that new
  text is checked against *before* a schism forms. Whether cognee's version is "better" is
  UNVERIFIED in both directions — vendor claim on one side, no measurement on ours.
- **Proposal shape:** minimal form needs **zero** COOPERATION.md bytes: a machine-readable
  `parley-deck/meta/glossary.yaml` (the single authority) plus an advisory CLI lint that reports
  canonical artifacts using near-miss variants of defined terms (exact/prefix matching, not an
  80% fuzzy embed — we have no embedding infrastructure and adding one violates the no-dependency
  constraint). Only if the deck later wants review to *reference* lint output does a §7 successor
  add one sentence. Follows the proven pattern: one authority, generated views, advisory tooling —
  the lint is not in the normative path and gates nothing.
- **Cost of being wrong:** the registry itself becomes a second authority — the exact disease it
  treats — if it rots or diverges from the text it describes. Also honest: the per-instance remedy
  may already be sufficient, and the marginal value of a general registry is unproven (RECALL —
  no measurement exists either way). That uncertainty, plus the maintenance cost, is why this
  ranks below rank 1, which has a measured gap and no maintenance surface.

### Mined and dropped (valid findings, no proposal)

- **Feedback / edge weights / "Truth-Subspace Reranking" (SECONDARY).** Automatic learned
  reweighting of what an agent sees next is the exact position of the frontier machinery deleted
  in 1.43.1 — rejected by the binding `speedup-tooling-evaluation` principle (a selector that
  cannot prove it never drops an objection turns a dropped objection into fabricated consent under
  Phase 2). The *advisory* form of the same idea already exists: §13 retro, §13.2 confident-error
  signals, §13.5 playbooks (COOPERATION.md lines ~1159–1191). Both forms covered. Drop.
- **`memify` as a standalone item.** The consolidate-in-place insight is real but is subsumed into
  rank 1; as a separate mechanism it maps onto §13.5 playbooks, which already distill closed ideas
  — Parley equivalent exists. Drop.
- **NodeSets / Datasets / ACL (SECONDARY).** Parley's grouping unit is `ideas/<slug>/` +
  participant file ownership (§6) with access control at the transport layer (branch protection,
  §11.B). No gap named; the parallel holds only at vocabulary level. Noise. Drop.
- **Memory Provenance visualization (SECONDARY).** Parley has §15.2 provenance tags with mandatory
  locators, owner / no-self-verdict rules (§15.1), and git history over every canonical file.
  Cognee's is a graph visualization; for a file-based deck the equivalent exists and the marginal
  picture is not shown to be better. Drop.
- **Time awareness / temporal search (SECONDARY).** Artifacts carry frontmatter dates,
  `meta/protocol-changelog.md` exists, git gives time. No gap named. Drop.
- **Contradiction detection.** Notable inversion: the cognee docs expose **no** contradiction /
  superseded-fact mechanism (per the brief's documented gaps), while Parley has §15.3 conflicting
  verdicts, `DISPUTED` handling, and §15.6 adversarial alternatives (COOPERATION.md lines
  ~1293–1360). Here the borrowing direction fails — Parley is ahead. No proposal.
- **Selective retrieval / search types (`GRAPH_COMPLETION` etc., SECONDARY).** The one real
  structural pressure — cognee *selects*, Parley *re-sends everything*
  (`gatherPriorRounds`/`gatherReviewContext`, PRIMARY per `00-prompt.md`) — is already named,
  ranked, and decided in this deck: rank-1 unbuilt intervention `protocol-phase-scoped-packet`,
  deterministic and instruction-layer. Cognee adds nothing to it beyond confirmation that the deck
  already gave itself. Building that beats anything this mining can import. No new proposal.

## Concerns / open questions

- Rank 1's "evidence of non-use" needs a defensible operationalization before its successor idea:
  citation-count over the corpus window is cheap but confuses *load-bearing-but-rarely-cited* text
  (§9 session-start) with dead text. The successor idea must define the metric, or reviewers will
  relitigate it per candidate.
- Rank 2's value is genuinely uncertain (RECALL, stated above). If cross-review has evidence the
  per-instance remedy suffices, I will downgrade it to the drop list without objection.
- All cognee claims herein are SECONDARY vendor self-description; none was independently verified,
  and none of my proposals depends on any cognee capability claim being *true* — only on the
  mechanism being *named*, which the brief supplies.
- I did not independently re-measure the Parley figures; they are taken as PRIMARY from
  `00-prompt.md`'s context list, with section-level corroboration from my own read of
  COOPERATION.md (Quickstart, §13, §15, Appendix A).

## Risks

- **Rank 1 deletes something load-bearing.** Mitigated: full §7 gate + git reversibility; the duty
  proposes, humans dispose. Residual risk accepted as lower than the measured cost of unchecked
  growth.
- **Rank 2 becomes a second authority / maintenance sink.** Mitigated by the one-authority +
  generated-view pattern and advisory-only enforcement; if the registry cannot be kept generated
  and single-sourced, it should not exist.
- **Both proposals could be used to justify protocol growth** ("the sunset duty needs §16 to
  explain itself"). Guard: the byte-accounting convention in rank 1(b) applies to rank 1's own
  successor idea first.
- **False-positive risk in mining itself:** I may have dropped a mechanism another participant
  correctly promotes. That is what round-1 independence and cross-review are for; my drop list
  states checkable reasons, not vibes.
