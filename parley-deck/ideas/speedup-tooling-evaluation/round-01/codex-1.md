---
agent: codex-1
idea: speedup-tooling-evaluation
round: 1
date: 2026-08-11
---

## Summary

The honest answer is the null answer: implement the two protocol changes without adding a retrieval or routing dependency. Use a deterministic, phase-scoped protocol packet to attack per-call protocol-loading cost, and enforce a finite fix-up budget that escalates rather than closes to attack per-idea review churn. They address the two measured terms directly. Cognee and graphify are retrieval systems; OmniRoute is primarily a model gateway/router with optional prompt compression. None reduces review round count. Retrieval or compression can reduce tokens only by taking authority over what an agent sees, which is the exact failure position of the deleted frontier optimization.

The measurements are accepted rather than re-derived: full protocol loading costs 3.3x median wall clock; review rounds rose from 1.6 to 5.1 with a maximum of 24; review bytes rose 7.2x; design rounds stayed flat. **PRIMARY evidence:** the task authority at `parley-deck/ideas/speedup-tooling-evaluation/00-prompt.md:16-42`.

My independent ranking by expected saving divided by the risk of losing a rule that catches a real error is:

| Rank | Option | Measured term moved | Expected saving | Risk of losing a rule/objection | Decision |
|---:|---|---|---|---|---|
| 1 | Finite fix-up budget with fail-closed escalation | Review round count, therefore repeated per-call tax and review bytes | Very high on the long tail | Very low: it stops work; it does not declare success | Implement/finish under §7 |
| 2 | Deterministic phase-scoped protocol packet | Protocol token volume and per-call latency | High and repeated on every agent call | Low only after the completeness gates below; otherwise high | Implement under §7 |
| 3 | graphify as an advisory code/protocol discovery and packet-audit tool | Human/agent exploration time only | Low to moderate locally | Low while its output is non-authoritative | Keep installed; do not put in the prompt-selection path |
| 4 | Cognee for non-canonical playbooks or optional recall | Retrieval quality, not measured round count | Low/unknown for this workload | Low only if never used to select normative or canonical context; plus substantial operational cost | Do not adopt now |
| 5 | OmniRoute as a pass-through gateway with compression off and exact model pinned | Provider availability/cost/possibly upstream latency, none of the established bottleneck | Near zero against the measured terms | Moderate operational/model-identity risk | Do not adopt for speedup |
| 6 | graphify or Cognee selecting normative context; OmniRoute applying semantic compression to Parley prompts | Token volume | Potentially material | Unbounded: a missed rule or objection can become recorded consent | Reject |

The recommended deployment is ranks 1 and 2 together, with no new dependency. Rank 1 limits how many times the tax is paid; rank 2 reduces the tax each time. The tools remain optional diagnostics outside the trusted path.

## Proposed approach

### 1. Make rank 1 a compiler, not retrieval

`protocol-phase-scoped-packet` should compile exact protocol text from an explicit, versioned manifest. It should not ask a graph, embedding search, LLM, memory layer, or heuristic which passages are relevant.

The minimum safe contract is:

1. Assign stable IDs and applicability metadata to normative units. Maintain an explicit `phase × track × transport -> rule IDs` manifest as part of the same write-once protocol release.
2. Include an invariant safety kernel in every packet: roster/quorum and file ownership, objection/no-suppression and “Silence = implicit agreement,” escalation, phase-transition/signoff rules, §7 change authority, and §15 verification integrity. Include those exact bytes, not summaries.
3. Fail closed to the full protocol if the phase/track/transport is unknown, the protocol hash is unexpected, a normative unit is unclassified, an ID is duplicated/missing, or composition cannot be proven complete.
4. Emit the protocol version/hash, included rule IDs, source line locators, and manifest version in the packet so an agent and reviewer can audit what was delivered.
5. Add a deterministic coverage check: every normative unit belongs either to the always-on kernel or at least one applicable packet. A protocol edit that adds or moves a normative unit must fail CI until the manifest is updated.
6. Shadow the packet against full-protocol runs before enabling it. Seed regression cases with an objection, a provenance conflict, a protocol-change request, and a strict-gate failure. A missing unit or different phase decision is a blocker, not an acceptable retrieval miss.

This packet may reduce protocol context. It must not also reduce canonical deliberation context. Those are different trust boundaries. **PRIMARY evidence:** `gatherPriorRounds` reads every markdown artifact from all earlier design rounds (`internal/runner/runner.go:936-965`), while `gatherReviewContext` reads `FINAL.md`, `IMPLEMENTATION.md`, and every prior review artifact (`internal/runner/phase58.go:276-306`). The prompt requires each participant to address every other participant (`internal/runner/runner.go:968-990`), and the protocol makes silence agreement (`parley-deck/COOPERATION.md:347-352`). Therefore, no phase-packet implementation should opportunistically summarize, graph-retrieve, top-k rank, or drop prior participant material.

If prior-artifact volume is addressed later, require a separately ratified, append-only objection/disposition ledger with deterministic completeness checks and participant-visible source locators. Until that exists and survives adversarial review, resend the canonical material. The protocol packet is the fast, bounded win; objection compression is not bundled into it.

### 2. Bound fix-up churn by escalation, never by consent

Set a finite deliberation fix-up budget (three cycles is a reasonable initial ceiling because the current driver already uses that value), allow an operator to grant an explicit recorded extension, and otherwise stop with a durable trajectory report. The cap must never auto-dismiss findings, manufacture a clean review, or mark the implementation complete.

The desired behavior already has a strong implementation shape in the current checkout: `MaxFixupCycles` defaults to 3 (`internal/driver/driver.go:62-67,95-105`); outstanding fixes at the ceiling return an escalation (`internal/driver/impl.go:277-281`); and the live protocol says the budget is an escalation threshold, not a close criterion (`parley-deck/COOPERATION.md:646-664`). **PRIMARY evidence:** direct reads at those stable repository locators. The open discrepancy with the task brief is recorded below; it does not create a case for a dependency.

### 3. Keep graphify advisory

Graphify is not an implementation of rank 1. It is useful for finding symbols and candidate relationships, but its traversal has no normative completeness proof.

**PRIMARY evidence — commands executed against the existing `graphify-out/graph.json`:**

```text
graphify query "gatherPriorRounds gatherReviewContext protocol loading review rounds" --budget 2400
  -> found gatherPriorRounds() at internal/runner/runner.go L938
  -> found gatherReviewContext() at internal/runner/phase58.go L278

graphify path "Verification integrity" "gatherReviewContext()"
  -> warning: source match was ambiguous
  -> No path found

graphify path "Changing this protocol" "gatherPriorRounds()"
  -> No path found
```

A broader natural-language query for the integration question started from generic nodes such as `COULD`, `without()`, and `Context` and surfaced unrelated TUI/inbox material before the runner context. Exact symbol names worked; semantic completeness did not. Several returned call edges were marked `INFERRED`. `graphify diagnose multigraph` found no same-endpoint edge collapse in the current post-build graph, but it also warned that raw producer loss must be measured earlier. That diagnostic says something useful about stored-edge shape; it does not prove that every normative rule was extracted, linked, matched, or traversed.

Safe uses are generating a candidate phase manifest for human review, locating code touched by a protocol rule, or auditing whether expected rule IDs are disconnected. Unsafe uses are selecting the packet, deciding that no other rule applies, resolving provenance, or selecting prior objections.

### 4. Do not integrate Cognee for this problem

Cognee is a memory/retrieval stack. Its vendor documentation describes relational storage for documents/provenance, vector storage for semantic similarity, and graph storage for relationships; its `cognify` pipeline performs chunking, LLM-based entity extraction, relationship construction, embeddings, summarization, and indexing. Its current recommended operations expose remember/recall/improve/forget and an MCP integration. **SECONDARY evidence only (vendor-controlled, not locally tested):** [Cognee architecture](https://docs.cognee.ai/core-concepts/architecture), [cognify pipeline](https://docs.cognee.ai/api-reference/cognify/cognify), [MCP tools](https://docs.cognee.ai/cognee-mcp/mcp-tools), and the vendor's pinned [v1.4.0 package metadata](https://github.com/topoteretes/cognee/blob/v1.4.0/pyproject.toml). The README/landing-page performance and reliability claims are vendor copy and are not evidence for this repository.

Integration would require a separate Python 3.10–3.14 environment, Cognee and its transitive database/LLM stack, an MCP or service boundary from the Go CLI, ingestion and update hooks, three-store lifecycle/pruning, an LLM and embedding provider (or local models), provenance back-links, and a repository-specific recall evaluation. The pinned vendor package metadata lists a large base dependency surface including LiteLLM, Instructor, SQLAlchemy, LanceDB, FastAPI, and graph libraries. **SECONDARY evidence only:** the pinned vendor metadata above and [vendor installation guide](https://docs.cognee.ai/getting-started/installation).

The trust grant would be larger than the benefit: protocol and artifacts go through chunking/extraction/embedding, are persisted in multiple indexes, and recall/top-k logic decides what reaches the agent. Cloud providers would also receive content unless both generation and embeddings are configured locally; the vendor docs explicitly warn that configuring only one can leave the other on the OpenAI default. **SECONDARY evidence only:** [vendor LLM-provider configuration](https://docs.cognee.ai/setup-configuration/llm-providers) and [local setup](https://docs.cognee.ai/guides/local-setup). This may improve recall, but it neither caps review rounds nor proves complete rule delivery. Since graphify is already installed for advisory retrieval, Cognee duplicates a non-bottleneck capability at materially higher integration cost.

### 5. Do not place OmniRoute in Parley Deck's semantic path

OmniRoute is a local OpenAI-compatible gateway that routes across provider/model connections, supports fallbacks and multiple selection strategies, and offers a prompt-compression pipeline. **SECONDARY evidence only (vendor repository, not locally installed/tested):** the vendor's pinned [release/v3.8.49 README](https://github.com/diegosouzapw/OmniRoute/blob/release/v3.8.49/README.md) and pinned [Auto-Combo documentation](https://github.com/diegosouzapw/OmniRoute/blob/release/v3.8.49/docs/routing/AUTO-COMBO.md). Installation is advertised as a global npm package or Docker service; integration then repoints each agent's base URL, installs endpoint/provider credentials, pins route policies, and adds service health/configuration/testing. Those are vendor claims and instructions, not verified local behavior.

Routing/fallback may help rate limits, provider outages, cost, or upstream latency. The measured evidence says the CLI is not the bottleneck and does not establish provider availability as the bottleneck. Auto-routing also conflicts with the audit requirement that a recorded roster/model describes the deliberation that actually happened: a “fast,” “cheap,” or fallback route may silently substitute model/company/capability.

Compression is worse. Vendor documentation says the gateway can remove/condense content, compress tool/file results, summarize or age history, and apply heuristic pruning before the provider. **SECONDARY evidence only:** the mutable vendor [compression guide](https://github.com/diegosouzapw/OmniRoute/wiki/Compression-Guide); its savings and preservation claims are vendor copy. Any non-byte-preserving mode occupies the deleted frontier optimizer's trust position. Even a mode marketed as “lite” or “safe” is inadmissible for normative or canonical Parley context until this project independently proves exact preservation of every rule, objection, source locator, structured field, and negation. At that point, a byte-preserving transport has no semantic compression win.

If OmniRoute is ever used for unrelated availability reasons, Parley calls should pin the exact model/provider, disable semantic compression and memory, log the resolved route, and fail rather than fall back to an unapproved model. That configuration does not solve this speed task.

### 6. Never delegate these decisions

No retrieval, memory, graph, router, compressor, or summarizer may decide:

- which normative protocol rules an agent receives;
- which canonical participant artifacts, objections, findings, or rebuttals are included;
- whether an objection is resolved, withdrawn, duplicated, stale, or safe to omit;
- whether silence counts as agreement or whether quorum/signoff is satisfied;
- provenance ownership, admissibility, verdict conflicts, or the §15 dependency check;
- protocol-change authority under §7;
- roster/model identity or mid-run fallback;
- whether an idea/review is complete.

Tools may propose or audit. Deterministic code over canonical bytes must enforce inclusion and veto claims of completeness. Humans and protocol signoffs retain the close decision.

## Concerns and open questions

1. **The task brief and current checkout disagree about rank 3.** The brief says `protocol-fixup-budget` is unbuilt and deliberation is unbounded (`00-prompt.md:34-37`). Current code defaults `MaxFixupCycles` to 3 and tests deliberation at 3 (`internal/driver/driver.go:62-67,95-105`; `internal/driver/track_test.go:44-51`), while the live protocol already describes a configured budget and fail-closed escalation (`COOPERATION.md:646-664`). These are **PRIMARY evidence** from stable repository locators. I did not re-derive the historical measurement and do not resolve the discrepancy here. Before implementation, determine whether the missing work is protocol ratification, manual/non-driver enforcement, configuration propagation, or stale task wording. The tool-adoption conclusion is unchanged.
2. What is the smallest always-on safety kernel? It should be chosen from incident history, not token targets. §7 and §15 must be included whole initially; fragmenting either before stable rule IDs exist recreates the retrieval problem.
3. How will normative units be identified without prose heuristics? A parser that merely searches `MUST` is insufficient because requirements also use imperative prose and cross-section composition.
4. How will overlays and protocol-version drift compose? The packet compiler needs the effective composed protocol and must fail to full read while per-idea pinning remains explicitly not in force (`COOPERATION.md:773-780`, **PRIMARY evidence**).
5. What independent acceptance test would justify retiring full-read shadow mode? Token savings alone are not sufficient. The gate must exercise known past failures and demonstrate packet completeness mechanically.
6. `graphify` warned that its installed skill text is v0.8.38 while the package is v0.8.41. This is **PRIMARY command output**, but it does not establish that the graph is stale. It is nevertheless another reason not to make the graph a normative runtime dependency without a pinned tool/index contract.
7. Cognee and OmniRoute were not installed or executed. Every capability, dependency, security, performance, and savings statement about them above is deliberately tagged `SECONDARY`; no vendor benchmark should enter an acceptance criterion for this idea.

## Risks

- **False consent is the dominant risk.** One omitted objection can invert the audit record because silence is agreement. Expected token savings cannot compensate for that failure mode.
- **A phase manifest can become stale.** Stable IDs, total normative-unit coverage, version hashes, and fail-to-full-read behavior are mandatory; review by taste is insufficient.
- **A fix-up cap can become a pressure-to-close mechanism.** The only safe terminal action at the ceiling is escalation/stop. An operator extension must be explicit and durable.
- **Protocol packets can accidentally absorb artifact compression.** Keep protocol selection and deliberation-history selection as separately ratified designs and trust boundaries.
- **Routing can destroy model provenance.** Transparent fallback changes the participant actually used and can correlate supposedly independent evidence.
- **Memory/indexes can retain sensitive or superseded content.** Multi-store deletion, tenant boundaries, provider disclosure, stale embeddings, and re-indexing become new operational obligations without touching the measured bottleneck.
- **Vendor “safe,” “lossless,” “private,” and savings claims are not project evidence.** They require local, adversarial, version-pinned tests before even non-canonical use; they are inadequate for normative context selection.
- **The null plan can still fail if implemented semantically.** “No new dependency” is not itself safe. Rank 1 is acceptable only as deterministic exact-text composition with a complete manifest and fail-closed fallback.
