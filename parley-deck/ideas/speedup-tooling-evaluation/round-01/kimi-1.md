---
agent: kimi-1
idea: speedup-tooling-evaluation
round: 1
date: 2026-08-11
---

# Summary

**The null answer is correct.** The fastest real win is rank 1 (`protocol-phase-scoped-packet`) plus
rank 3 (`protocol-fixup-budget`), both already named, both requiring §7 and nothing else. None of the
three tools moves a measured term by a mechanism this protocol can admit, and two of them occupy —
by design — the exact position the deleted frontier code was removed from: deciding what an agent
sees, without being able to prove completeness.

Q1 — measured-bottleneck fit, per tool:

| Tool | What it is | Measured term it could touch | Verdict |
|---|---|---|---|
| graphify | Local, deterministic code+doc graph; heading-level index of this repo (PRIMARY, Appendix A) | None directly. Bodies still must be read; rounds are protocol-driven | Different bottleneck (navigation, not volume or round count) |
| cognee | LLM-mediated memory layer: ingest → "cognify" → vector+graph recall, injected context (SECONDARY, vendor README) | None. This deck over-remembers (quadratic history re-send); it does not under-remember | Solves a problem we do not have |
| omniroute | OpenAI-compatible AI gateway: multi-provider routing, auto-fallback, lossy prompt compression (SECONDARY, vendor README) | Per-call token volume — but only via lossy compression or silent model swap | Right term, inadmissible mechanism |

The two measured terms are `cost of a round × number of rounds`
(`protocol-read-cost-regression/FINAL.md:15-27`, PRIMARY — quoted in Appendix C). Rank 1 cuts the
first term and is paid once per round, so it compounds; rank 3 bounds the tail of the second. A tool
that touches neither is a cost, not a win.

# Proposed approach

**Recommendation: adopt nothing. Build rank 1 and rank 3 as protocol changes under §7. Keep
graphify, confined to code-side navigation, where it is harmless and occasionally useful.**

## Q2 — graphify specifically: it is not an implementation of rank 1

I ran it against the existing graph (`graphify-out/graph.json`, built at current HEAD — PRIMARY,
Appendix A). Findings:

1. **The graph is a structural index, not a normative one.** 60% of edges are `contains`
   (13,794 of 23,135 — PRIMARY, Appendix A). There is no relation for "binds in phase P." §15.7's
   per-track binding table, §9's session-start duties, §14's human brake — applicability is
   *normative text*, and no extractor produces it. Rank 1 needs a phase→sections mapping whose
   correctness is a ratified claim; a graph can store such a mapping but cannot be its justification.
   Once the mapping is hand-authored and ratified under §7, the lookup is ~40 lines of Go in the
   runner. The graph adds nothing to the part that is hard.
2. **Retrieval by traversal demonstrably misses the normative answer.** `graphify query "what must a
   participant do in Phase 2 cross-review"` returned Go test files and an unrelated idea's round
   file; the Phase 2 section of `COOPERATION.md` was not in the output. A query naming §15.2's exact
   vocabulary started from three nodes all labeled "Provenance" and BFS-wandered into inbox and
   FINAL.md documents (PRIMARY, Appendix B). Start-node selection is keyword matching; when it
   misses, it misses silently. Silent miss is the disqualifying property — see Risks.
3. **Nodes are headings; bodies are not in the graph.** `explain` on the §15 node correctly lists
   §15.1–§15.7 with line numbers (PRIMARY, Appendix B) — a table-of-contents service. The agent
   still must `Read` the lines. The protocol file already contains this map (its Quickstart and
   heading structure); `grep '^##'` reproduces it in milliseconds with no staleness window.
4. **Three copies of the protocol are indexed** — `internal/protocol/defaults/COOPERATION.md`
   (70 nodes), `parley-deck/COOPERATION.md` (70), `parley-deck-skill/references/COOPERATION.md` (1)
   (PRIMARY, Appendix A). Traversal cannot distinguish the authoritative deck view from the embedded
   default. For a normative read, landing on the wrong copy is rule-loss, not navigation.
5. **Freshness is manual.** `built_at_commit` equals HEAD today (PRIMARY, Appendix A), but nothing
   rebuilds the graph on protocol edit; a stale index hands out wrong line numbers for load-bearing
   text. The toolchain already drifts: every invocation prints `skill is from graphify 0.8.38,
   package is 0.8.41` (PRIMARY, Appendix B).

What graphify is legitimately good for here: `path` / `affected` / `explain` on Go symbols while
*implementing* the rank-2 ledger follow-ups — human-and-agent navigation aid, no normative gating.
That use is already paid for (installed, zero new dependency) and carries near-zero risk precisely
because it never decides what an agent sees of the protocol.

## Q3 — cognee and omniroute: what they do, cost, and what we would be trusting

**cognee** (not installed; could not test — all capability claims SECONDARY from vendor README,
Appendix D): a Python memory platform. Ingested documents are chunked, embedded, and "cognified" —
an LLM generates entities, relationships, and an ontology — then `recall` auto-routes a search
strategy and its Claude Code plugin "injects relevant context on every prompt." Integration cost:
Python 3.10+ service or Docker, a vector store and graph store (SQLite/LanceDB/Kuzu embedded, or
Postgres — whose graph backend the vendor itself marks demo-only), an `LLM_API_KEY` for extraction,
MCP or per-agent-CLI hooks. For a Go CLI that drives *other agents'* CLIs, this is a sidecar stack
plus per-participant wiring.

What we would be trusting it with: an **LLM-generated paraphrase of normative text, selected by an
opaque ranking, injected into context**. Three independent disqualifiers under the deleted-frontier
standard: (a) extraction is generative — the ledger contract that survived 1.43.1 requires "the
exact scoped proposition, never a generated paraphrase" (`protocol-read-cost-regression/FINAL.md:64`,
PRIMARY); (b) recall completeness is unprovable — its own BEAM benchmark scores 0.67–0.79, i.e. the
vendor's best-case claim is *partial* recall on a memory task, and memory benchmarks say nothing
about never dropping a MUST; (c) it addresses cross-session memory, the opposite of our measured
disease, which is over-delivery of history.

**omniroute** (not installed; SECONDARY from vendor README and a third-party listing, Appendix D):
an AI gateway — one OpenAI-compatible endpoint over ~290+ providers, quota-aware auto-fallback,
format translation, and prompt compression advertised at ~15%/~30%/~50%/~75% token reduction across
four modes (vendor numbers; mechanism "Caveman" compression is lossy prompt rewriting). Structural
problem first: parley does not own the model calls — the participant CLIs (claude, codex, hermes,
kimi) do. OmniRoute would sit between each agent CLI and its provider, per participant, outside the
deck's control and outside the roster. What we would be trusting it with, if wired in anyway:
(a) **compression rewrites the protocol text before the agent sees it** — the "Never cut" list
(`FINAL.md:111-117`, PRIMARY) exists precisely because modals, negations, conditions and exceptions
are what lossy compression drops first; (b) **auto-fallback silently changes which model answered**,
breaking participant identity (§2 roster), and breaking §15.6's accounting, where unanimity among
related models is a shared prior — you cannot compute correlation if you cannot name the model;
(c) free-tier routing ships deliberation content — including anything a participant pasted — to
arbitrary third-party providers, against the deck's no-secrets posture. It touches no measured term
by an admissible path: round count is protocol-driven, and the 3.3× per-call cost is protocol-text
volume, which it can only reduce by cutting protocol text.

## Q4/Q5 — the null answer, and the ranking by (expected saving)/(rule-loss risk)

1. **rank 1 + rank 3, no new dependency — adopt.** Only option that attacks both measured terms with
   mechanisms the protocol itself governs (§7 idea → ratification). Rank 1's packet needs a
   hand-authored phase→sections mapping, a completeness argument, and the ledger's fallback shape:
   any missing/ambiguous/challenged state falls back to full read, visibly. Risk is process risk,
   which §7 exists to control.
2. **graphify, confined to code navigation — keep, do not extend.** Already installed; helps
   implement the ranked work; never gates normative reads. Saving small but real; risk near zero in
   that confinement.
3. **omniroute — reject.** Unproven saving (vendor numbers), and its two mechanisms are each
   independently inadmissible here.
4. **cognee — reject.** No measured term moves; largest trust surface; generative paraphrase of
   normative text is the frontier failure with a database attached.

**Must never be delegated to a tool:** completeness of normative reads (which sections/MUSTs an
agent is held to); objection lifecycle (an objection is live until its owner disposes of it);
issuance and provenance of verification verdicts (§15.1–15.3); verdict-conflict resolution (§15.3 —
never by counting); consensus close (§15.5 — signoffs, not judgment); §7 changes; the §14 human
brake; user rulings, carried verbatim. A tool may *point at* any of these; it may not *select,
summarize, compress, or dispose of* them.

# Concerns and open questions

- **Rank 3's expected saving is tail-only, and I want that on record.** The closed FINAL found the
  review explosion is a step change tracking roster growth, and the worst ideas kept finding fresh
  MAJORs at rounds 19–24 — a cap bounds cost without removing cause. Rank 1 is the median win;
  rank 3 is insurance. If the deck wants the median review-round count down, the lever is roster
  sizing / Phase-2 addressing scope, which is a different idea. (Basis: `FINAL.md:29-31`, PRIMARY.)
- **Rank 1's saving is unquantified until the mapping exists.** §15.7 binds on every track; §9, §14
  and the modals travel with most phases. The packet plausibly drops the appendices (§11–§13) and
  non-active phases from a ~26k-token read, but "plausibly" is RECALL until the mapping is drafted
  and measured. The 3.3× figure itself rests on n=3/arm — consensus-established, not to be
  re-litigated here, but the rank-1 design should re-measure with more replicates.
- **Steelman for graphify-as-packet-index (anticipated §15.6):** it is local, deterministic, and its
  output is auditable — unlike cognee/omniroute, one *could* pin its selection rules and diff them.
  My counter stands on point Q2.1: the hard part of rank 1 is the ratified mapping, and once that
  exists the storage mechanism is incidental. If round 2 disagrees, the decisive observation would
  be a demonstrated completeness property of traversal, which my two probes contradict.
- **Any future tool claiming a seat here** should be gated on the frontier standard, stated as a
  test: prove it never drops a live objection or an applicable rule, or fall back to full text
  visibly on any doubt; deterministic and auditable selection; raw artifacts always reachable; fails
  loud. A vendor benchmark is not such a proof.
- Ambiguity note: "omniroute" was evaluated as the OSS gateway at
  `github.com/diegosouzapw/OmniRoute` (the prominent match). If the owner meant a different
  omniroute, the Q3 analysis needs a re-run against that locator.

# Risks

- **Silence → consent amplification (the standing one).** Phase 2 rule 1: "Silence = implicit
  agreement" (`COOPERATION.md:350`, PRIMARY). Any selector that drops an objection manufactures
  consent; any selector that drops a rule manufactures compliance theater. This is not hypothetical:
  1.43.0 shipped over a MIXED round-3 review and a RESERVED signoff and was deleted in 1.43.1
  (`CHANGELOG.md:3-30`, PRIMARY — quoted in Appendix C). cognee (recall ranking + injected context)
  and omniroute (prompt compression) are both selectors by construction. graphify becomes one the
  moment its traversal gates a protocol read — which is why the confinement in Q5.2 matters.
- **Authority confusion across protocol copies.** Three indexed copies of `COOPERATION.md`
  (PRIMARY, Appendix A) plus a partially-implemented overlay regime (`COOPERATION.md:767-771`,
  PRIMARY) mean any graph-mediated read can silently bind an agent to the wrong text.
- **Dependency and drift risk even in the benign path.** graphify is already version-drifting
  (skill 0.8.38 vs package 0.8.41) and its build suppressed 58 candidate edges at extraction time
  (`producer_suppression_sites: 58`, PRIMARY, Appendix A) — acceptable for navigation, disqualifying
  for normative completeness. cognee/omniroute add a Python service stack / an always-on proxy to a
  Go CLI that currently needs neither.
- **Adoption-as-Precedent risk.** "It worked for finding functions" slides into "let it assemble the
  packet." The line must be drawn in the adopting text, not in habit: tools may assist *finding*;
  the protocol text and §7-ratified mappings decide *binding*.

---

## Appendix A — PRIMARY: graph contents (commands run 2026-08-11, this repo)

`python3` over `graphify-out/graph.json`:

```
nodes: 17346   links: 23135   hyperedges: 75   directed: False
built_at_commit: 41e6cd6ef56f0091c27bc4d486d846cd112fa197   # == `git rev-parse HEAD`
file_types: document 13,758 / code 3,069 / concept 234 / rationale 285
relations: contains 13,794; calls 4,199; references 3,912; method 318; cites 256;
           conceptually_related_to 213; implements 193; rationale_for 102; ...
COOPERATION.md copies: internal/protocol/defaults (70 nodes), parley-deck (70),
                       parley-deck-skill/references (1)
§15 node: parley_deck_cooperation_15_verification_integrity → parley-deck/COOPERATION.md L1232
```

`graphify diagnose multigraph` (read-only): `same_endpoint_collapsed_edges: 0`,
`exact_duplicate_edges: 0`, `producer_suppression_sites: 58` (example: `L1434 seen_dyn_pairs
arity=unknown`).

## Appendix B — PRIMARY: graphify behavior probes

`graphify query "what must a participant do in Phase 2 cross-review" --budget 1200` →
24+ nodes; top hits `internal/driver/phase_event_test.go`, `internal/driver/cursor.go`,
`ideas/meta-protocol-change-devx-speed/round-01/codex-1.md`; Phase 2 section of `COOPERATION.md`
absent from output.

`graphify query "provenance PRIMARY SECONDARY RECALL verification verdict" --budget 800` →
`Traversal: BFS depth=2 | Start: ['Provenance', 'Provenance', 'Provenance']`; hits are inbox and
FINAL.md heading nodes; §15.2 body not returned (bodies are not nodes).

`graphify explain "parley_deck_cooperation_15_verification_integrity"` → correct neighborhood:
`15.1 Scope…` through `15.7 Per-track binding`, `Degree: 8`.

Every invocation also printed: `warning: skill is from graphify 0.8.38, package is 0.8.41. Run
'graphify install' to update.`

## Appendix C — PRIMARY: the deleted frontier and the measured terms

`CHANGELOG.md:5-19` (1.43.1): frontier machinery deleted; @codex-1's never-withdrawn RESERVED
objection: *"Keeping unreachable safety code behind a constant invites exactly the rot the tests
claim to prevent…"*; "1.43.0 also still perturbed prompts … while delivering no speedup."

`ideas/protocol-read-cost-regression/FINAL.md:18-24`: 3.3× per-call; review rounds 1.6→5.1 (max 24);
review bytes 7.2×; design rounds flat; protocol 720→1,359 lines. `:64`: "the exact scoped
proposition, never a generated paraphrase". `:86-87`: "A dropped objection is not a lost datum, it
is agreement that was never given."

`COOPERATION.md:350`: "**Address every other active agent explicitly.** Silence = implicit
agreement."

## Appendix D — SECONDARY: vendor and third-party claims (not locally testable)

cognee — vendor README, `github.com/topoteretes/cognee` (fetched 2026-08-11): "open-source AI
memory platform"; `remember/recall/forget/improve`; "recall with auto-routing (picks best search
strategy automatically)"; Claude Code plugin "injects relevant context on every prompt"; Postgres
graph backend "currently a released as a demo feature"; BEAM benchmark 0.79 @100K / 0.67 @10M
tokens, "a directional signal rather than a definitive measure." All vendor copy: SECONDARY.

omniroute — vendor README, `github.com/diegosouzapw/OmniRoute` (fetched 2026-08-11): "one endpoint,
290+ providers … quota-aware auto-fallback … compression saves 15-95% tokens" (README hero text);
compression tiers Lite ~15% / Caveman ~30% / Aggressive ~50% / Ultra ~75% per third-party listing
(everydev.ai, 2026-05-04). Vendor/third-party copy: SECONDARY. Note the README's own 15–95% range
versus the listing's 15–75% — marketing ranges, not measurements.
