---
agent: hermes-1
idea: speedup-tooling-evaluation
round: 1
date: 2026-08-11
---

## Summary

The fastest real win is rank 1 + rank 3 with no new dependency. That is my
recommendation, and I argue it below. None of the three tools — graphify, cognee,
omniroute — touches the two measured terms (per-call protocol read cost, review
round churn) in a way that survives the standard the project already set: the
deleted frontier machinery (v1.43.1) was removed because it could silently drop
a participant objection, and under Phase 2 rule 1 "Silence = implicit agreement"
a dropped objection becomes recorded consent. Any tool that decides what an agent
sees occupies that same position.

The single most important finding I made this round is PRIMARY and structural:
the Go runner that builds prompts — `BuildRoundOnePrompt`, `BuildRoundPrompt`,
`gatherPriorRounds`, `gatherReviewContext`, `buildConsensusDraftPrompt` — never
reads COOPERATION.md at all (zero references in `internal/runner/runner.go`,
`internal/runner/phase58.go`, `internal/app/driver_consensus.go`). The 3.3×
per-call wall-clock cost is not a Go function call; it is the agent reading the
1,372-line COOPERATION.md into its context window as instructed by the
parley-deck skill ("Always read parley-deck/COOPERATION.md first"). This means
the per-call bottleneck lives in the skill's instruction to read the whole
protocol, and in the agent's compliance — not in a code path a tool can
intercept from outside.

## Proposed approach

Rank the options by (expected saving) / (risk of losing a rule that catches real
errors), then name what must never be delegated.

### The ranking

1. **rank 3 — `protocol-fixup-budget` (needs §7).** Highest expected saving,
   lowest rule-weakening risk. The measured tail is review rounds reaching 24
   with fresh MAJORs at rounds 19–24 (PRIMARY: FINAL.md L29–31). An unbounded
   fix-up cap does not cause the tail but enables it. A bounded cap directly
   reduces the `rounds × per-call tax` product — the per-call tax is paid again
   inside every extra cycle (PRIMARY: 00-prompt.md L32–33). The risk to
   rule-catching is low because the fix is a ceiling on cycles, not a filter on
   content: no objection, verdict, or §15 provenance datum is dropped. The
   worst case is that a real finding arrives after the cap and must be deferred
   to a follow-up idea — which is a scheduling cost, not a consent fabrication.

2. **rank 1 — `protocol-phase-scoped-packet` (needs §7).** Second-highest
   saving, moderate rule-weakening risk if implemented naively, manageable if
   the fallback contract from FINAL.md is carried over. The per-call protocol
   read is 3.3× median wall clock (PRIMARY: 00-prompt.md L24–25), and the
   protocol grew 720→1,372 lines monotonically in ten weeks (PRIMARY:
   00-prompt.md L28). Sending only the sections a phase needs directly attacks
   this term. The risk: deciding which sections a phase "needs" is exactly the
   filtering position. The protection is the same one FINAL.md already
   specified for the deleted frontier — any missing, invalid, ambiguous, or
   challenged state falls back to full history, visibly in the prompt (PRIMARY:
   FINAL.md L79–83, gates G2–G5). With those gates, rank 1 is safe. Without
   them, it is the frontier machine rebuilt under a new name.

3. **graphify — keep for codebase navigation, do not use for protocol
   delivery.** graphify is installed (v0.8.41, PRIMARY: 00-prompt.md L49) and
   already has a built graph for this repo (PRIMARY: I ran `graphify diagnose
   multigraph` — 17,346 nodes, 23,135 edges, 89% EXTRACTED, 11% INFERRED at
   avg confidence 0.81, built from commit 41e6cd6). Its edge types are
   `contains` (13,794), `calls` (4,199), `references` (3,912),
   `conceptually_related_to` (213), `cites` (256) — it is an AST-extraction
   code graph, not a semantic-retrieval or compression layer. It costs zero
   API tokens (PRIMARY: GRAPH_REPORT.md — "Token cost: 0 input · 0 output").
   It is useful for what CLAUDE.md already uses it for: answering codebase
   questions without raw grep. But it does not address either measured term:
   - It does not reduce per-call protocol read cost. `graphify query` returns
     node names and source locators (PRIMARY: I ran `graphify query "protocol
     read cost review rounds fix-up budget scoped packet"` — it returned 23
     nodes as `NODE <title> [src=<path> loc=<line>]` plus `EDGE` lines, not
     section text). An agent still has to read the actual file content to use
     a rule. The graph gives pointers, not packets.
   - It does not reduce review round count. It has no opinion on fix-up caps.
   - Could a graph of the PROTOCOL let an agent load only what a phase needs?
     Theoretically, but only if someone builds a layer that selects sections
     from the graph and hands them to the agent — which is rank 1 with graphify
     as the selector, which is the dangerous position. The graph itself is
     inert; the selector is the risk. And the graph's edges are structural
     (`contains`, `calls`), not normative — a `contains` edge from "Phase 2"
     to "Silence = implicit agreement" does not encode that the latter is a
     MUST that cannot be dropped. A graph traversal that reaches a normative
     rule without reading its full modal context (MUST, MUST NOT, unless,
     except) loses the conditions and exceptions that §15 and the Phase 2
     rules carry. That is the same failure that killed the frontier code.

4. **cognee — do not adopt for protocol delivery.** cognee is an open-source
   Python "AI memory platform for agents" (SECONDARY: cognee.ai landing page,
   GitHub README topoteretes/cognee). API: `remember()`, `recall()`,
   `forget()`, `improve()` (SECONDARY: README quickstart). It ingests data,
   builds a knowledge graph with vector embeddings + graph reasoning, and
   provides recall with "auto-routing (picks best search strategy
   automatically)" (SECONDARY: README). Integration cost: a Python dependency
   (`pip install cognee`), a graph store, a vector store, an LLM API key for
   its extraction pipeline, and Docker for the MCP server (SECONDARY: README).
   What we would be trusting it with: the recall layer decides what an agent
   sees from the protocol and prior rounds. That is exactly the position the
   deleted frontier code occupied — "a memory or routing layer that decides
   what an agent sees" (PRIMARY: 00-prompt.md L71–73). cognee's `recall` does
   auto-routing and "picks best search strategy automatically" — that is a
   selector we cannot prove never drops an objection. Under §15.2 a vendor
   claim is SECONDARY at best and cannot be tagged PRIMARY without a stable
   locator and quoted output; I have not run cognee (not installed, PRIMARY:
   00-prompt.md L50), so I cannot verify its recall completeness. The standard
   from the deleted frontier was not "it probably works" — it was "prove it
   never drops an objection, or do not ship it." cognee meets that standard
   no better than the frontier code did, and it adds a dependency and an LLM
   extraction cost the frontier code did not have.

5. **omniroute — does not address the bottleneck at all.** omniroute is an AI
   gateway/router, not a memory or retrieval layer (SECONDARY: GitHub
   diegosouzapw/OmniRoute, npm registry). It aggregates 290+ AI providers
   behind one endpoint with auto-fallback and "RTK + Caveman stacked
   compression saves 15–95% tokens" (SECONDARY: README). It is a model-routing
   proxy. The measured bottleneck is not model latency, token cost of model
   calls, or provider availability — it is protocol loading and review round
   count. omniroute touches none of those. Its compression claim is about
   compressing context sent to models, which is a different axis from the
   rank-5 `compression-experiment` already deferred in FINAL.md, and that
   experiment was itself deferred pending @hermes-1's 10-rule falsification
   test which predicts failure on at least 3 of 10 (PRIMARY: FINAL.md L108–109).
   Adopting omniroute would add a routing dependency for a problem it does not
   solve.

### What must never be delegated to a tool

The list FINAL.md already named (PRIMARY: FINAL.md L113–117): the authoritative
protocol; applicable modals, negations, conditions and exceptions; round-1
independent proposals; every live objection or finding; provenance and
verdict-conflict data; the §15.6 correlated-agreement audit; FINAL.md, the
current diff, acceptance criteria and check results; explicit user rulings;
ownership, no-secrets and the §14 human brake; and on-demand access to raw
historical artifacts. No tool — graphify, cognee, omniroute, or any future
selector — may sit between these and the agent. A tool may point to them; it
may not choose which ones the agent sees.

## Concerns and open questions

1. **Where exactly does the per-call read cost live?** I established it is the
   skill instruction ("read COOPERATION.md first"), not the Go runner. But
   rank 1 (`protocol-phase-scoped-packet`) needs §7 — does the §7 change target
   the skill's instruction, the protocol's own structure, or a future runner
   path that does embed protocol sections? The prompt's note that rank 1
   "Needs §7" suggests a protocol change, but the protocol is already
   sectioned (§0–§15). If the fix is "the skill instructs the agent to read
   only the relevant sections," that is a skill change, not a §7 protocol
   change — unless the protocol itself must encode the phase→section mapping
   as a normative rule, which would make it a §7 change. This matters for
   implementation scope.

2. **Can graphify's graph encode normative modals?** The graph edges are
   `contains`, `calls`, `references` — structural, not modal. A phase-scoped
   packet built on graphify would need a modal-aware selector that understands
   MUST vs SHOULD vs MAY, and the conditions/exceptions on each rule. If
   graphify cannot encode that, it cannot safely select protocol sections. I
   did not test whether graphify can be extended with custom edge types — that
   is a graphify feature question, not something I can verify from the built
   graph alone.

3. **cognee's recall completeness is unverified.** I cannot test it locally
   (not installed). The vendor claims auto-routing recall, but §15.2 says a
   vendor claim without a stable locator is SECONDARY, and a material claim
   reaching FINAL.md with only RECALL support MUST remain UNVERIFIED. If
   someone wants to propose cognee seriously, they must install it, ingest the
   protocol, and run the G5 fixture from FINAL.md (orphaned dissent: a
   round-1 minority objection not restated in round 2 must survive recall or
   trigger fallback). Until that test passes, cognee is UNVERIFIED for this
   use.

## Risks

1. **The frontier failure repeats.** The highest-severity risk is not that a
   tool fails to speed things up — it is that a tool is adopted for speed and
   silently drops a rule that catches real errors. The project already paid
   this cost: the frontier optimization was built, reviewed three times, and
   deleted in 1.43.1 (PRIMARY: 00-prompt.md L40–42). Any selector — graphify-
   based, cognee-based, or hand-rolled — that stands between the protocol and
   the agent must meet the G5 standard: an orphaned objection must reach the
   next round or trigger visible fallback. A tool that cannot prove this is a
   cost, not a win.

2. **Rank 1 implemented without the fallback contract is the frontier machine
   rebuilt.** The protection is not the packet; it is the fallback. FINAL.md
   gates G2–G5 specify the fallback (missing/invalid/ambiguous/challenged →
   full history, visibly in the prompt). If rank 1 ships without those gates,
   it will be reviewed and deleted a second time, and the review cost will
   itself contribute to the review-round churn it was meant to reduce.

3. **Adding a dependency for a problem it does not solve is a net negative.**
   cognee adds a Python stack, a graph store, a vector store, an LLM API key,
   and Docker. omniroute adds a routing proxy. Neither touches the measured
   terms. The cost is not just the dependency — it is the maintenance surface,
   the review burden when the tool is proposed, and the risk that the tool is
   later found to occupy the dangerous position and must be deleted (another
   review cycle, another round of churn).

4. **The null answer is possibly correct and possibly the only answer that
   moves the numbers without adding risk.** rank 1 + rank 3 are already named,
   already ranked, and already understood by the participants who signed off
   on FINAL.md. They need §7 protocol changes, but they do not need a new
   tool, a new dependency, or a new selector. The honest read of the evidence
   is that the protocol grew (720→1,372 lines) and the fix-up cap is unbounded
   — and the fix for both is a protocol change, not a tool adoption.
