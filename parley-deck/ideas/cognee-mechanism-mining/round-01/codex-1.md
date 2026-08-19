---
agent: codex-1
idea: cognee-mechanism-mining
round: 1
date: 2026-08-14
---

## Summary

I would carry forward only three items: one new audit mechanism, one conditional revival of an
existing Parley design, and one low-priority retrospective diagnostic. Cognee is not the better
design in any of the three. Its useful contribution is naming operations Parley can implement with
stronger ownership, provenance, and fail-open rules.

Ranked result:

1. **ADOPT — prompt-context provenance receipts.** This is the one clearly new proposal.
2. **CONDITIONAL — explicit `forget` as owner-authorized retirement with tombstones.** Parley
   already specified the stronger mechanism and then removed its implementation; do not start a
   parallel design.
3. **ADOPT, LOW PRIORITY — a separate advisory terminology-consolidation pass in `parley retro`.**
4. **DROP — ontology grounding and `ontology_valid`.** Parley's deterministic authority/view
   pattern is better than fuzzy canonicalization.
5. **DROP — feedback-weighted or “truth-subspace” reranking.** Parley's explicit retrospective and
   ratification path is safer.
6. **DROP — temporal search as a protocol mechanism.** It does not solve contradiction or
   supersession, and Parley's actual versioning gap is already named.
7. **DROP — NodeSets, datasets/ACL, and graph/hybrid search types.** These do not survive the
   analogy to Parley's shared, files-canonical workflow.

Cognee descriptions below are **SECONDARY** vendor self-description from
`parley-deck/ideas/cognee-mechanism-mining/00-prompt.md:58-106`; none is treated as demonstrated.

## Proposed approach

### 1. ADOPT — prompt-context provenance receipts

**Cognee mechanism.** **SECONDARY:** “Memory Provenance” presents the ownership and data-flow story
behind memory across tenants, users, agents, datasets, and files
(`00-prompt.md:93-94`).

**Parley gap.** Parley records outcomes and process events, but not an exact, reviewable account of
the inputs that formed a dispatched prompt. **PRIMARY:** `store.Event` is an append-only JSONL event
shape (`internal/store/events.go:16-55`); `run.created` records task, mode, idea, participants, and
runtime (`internal/runcontrol/runcontrol.go:63-72`); `agent.started` records process/artifact data
(`internal/runner/runner.go:545-568`). The check

```text
rg -n 'Type: "[^"]*(context|prompt|receipt|manifest)' \
  internal/runner internal/app internal/runcontrol internal/hitl --glob '*.go'
```

returned only `run.manifest_deferred` and `agent.acp.prompt_completed`; the companion search for
`context receipt|context manifest|prompt receipt|prompt manifest|input sha|prompt sha` returned no
matches. That proves only the searched implementation surface, not the impossibility of a
differently named equivalent; the broader absence claim is **RECALL / UNVERIFIED**.

**Existing equivalent / is cognee better?** The ratified phase packet already requires
`sourceSha256`, a complete omission index, and visible full-fallback reasons for protocol text
(**PRIMARY:** `parley-deck/ideas/meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md:35-49`).
That is a partial equivalent and stronger than the vendor visualization claim, but it covers the
protocol packet rather than every input to the final prompt. Cognee is not better; its description
names the missing audit surface.

**Concrete proposal.** Open an ordinary tooling successor, not a protocol change, for a shadow-mode
`context.receipt` event emitted at the final prompt-serialization boundary. Each receipt should carry:

- run/segment/attempt/agent/phase;
- SHA-256 and byte length of the exact dispatched prompt;
- an ordered segment map: stable source path or input kind, source SHA-256, byte range in the final
  prompt, inclusion mode (`full`, `packet`, `retirement-ledger`, `generated`), and transformation ID;
- protocol `sourceSha256`, packet/fallback mode, and every omission/fallback reason;
- hashes rather than copied prompt bodies, so the event does not become a second secret-bearing
  transcript.

The receipt observes what was sent; it must never select what is sent. Start as audit-only. Making
it an artifact-validity or close gate would change protocol semantics and therefore requires a §7
SUCCESSOR idea (**PRIMARY:** `parley-deck/COOPERATION.md:743-786`). This proposal adds zero bytes to
`COOPERATION.md` in its shadow form.

**Cost of being wrong.** A receipt generated from intended inputs rather than the serialized prompt
creates false confidence; storing raw excerpts leaks secrets; unstable hashes make audits
irreproducible; and per-call segment records add run-log volume. Build at the final byte boundary,
secret-test it, and do not call it complete unless mutation tests prove that dropping or changing an
input changes the receipt.

### 2. CONDITIONAL — explicit `forget` as owner-authorized retirement with tombstones

**Cognee mechanism.** **SECONDARY:** `forget` is a top-level operation, but the supplied vendor brief
does not specify its deletion, ownership, audit, or conflict semantics (`00-prompt.md:60-64`).

**Parley gap.** The live runner still concatenates all earlier design artifacts and all earlier
review artifacts. **PRIMARY:** `gatherPriorRounds` walks rounds `1..N-1` and reads every Markdown
artifact (`internal/runner/runner.go:936-965`), while `gatherReviewContext` reads `FINAL.md`,
`IMPLEMENTATION.md`, and every prior review round (`internal/runner/phase58.go:276-306`). The prompt
then orders the participant to read every prior artifact (`internal/runner/runner.go:968-990`). This
is stricter than Phase 2, which requires addressing every active peer and a counter-proposal for
disagreement but does not require every historical version (**PRIMARY:**
`parley-deck/COOPERATION.md:347-352`).

**Existing equivalent / is cognee better?** Parley already designed the stronger equivalent:
owner-namespaced IDs, exact propositions, source hashes, `OPEN|RESOLVED|DEFERRED|SUPERSEDED`,
append-only transitions, owner-only objection disposition, terminal tombstones, and full-history
fallback (**PRIMARY:** `parley-deck/ideas/protocol-read-cost-regression/FINAL.md:53-100`). It was
implemented and deleted in v1.43.1, and the current cost remains
(**PRIMARY:** `parley-deck/ideas/meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md:191-196`).
Cognee contributes only the useful verb; its asserted mechanism is not better and supplies no
witness for Parley's dropped-objection hazard.

**Concrete proposal.** Do not open a fresh cognee-inspired design. If the owner wants this term
attacked again, open a successor to `protocol-read-cost-regression` that reuses its ledger contract,
adds rank 1's actual-context receipt, and begins with a failure analysis of v1.43.1. The invariant
remains: previous round in full; every older live item exact and in full; only an item's owner may
retire an objection; terminal items leave tombstones; any missing, ambiguous, challenged, or
hash-mismatched state visibly falls back to full history; consensus drafting retains full history.
Keep it implementation-scoped unless consensus chooses to make the ledger normative; that choice
would require a §7 SUCCESSOR. No net protocol bytes are justified before the implementation-scoped
experiment succeeds.

**Cost of being wrong.** This is the highest-risk candidate. Phase 2 makes silence implicit
agreement (**PRIMARY:** `COOPERATION.md:349-352`), so one wrongly retired objection becomes consent
never given. The safer error is a noisy full-history fallback. Receipt-first observability is
necessary but not sufficient: it can prove what was sent, not that a human-authored retirement was
semantically correct.

### 3. ADOPT, LOW PRIORITY — separate advisory terminology consolidation

**Cognee mechanism.** **SECONDARY:** `memify` is a distinct post-graph enrichment pass; Entity
Consolidation and Entity Deduplication rewrite fragmented descriptions and merge near-duplicates
(`00-prompt.md:67-73`).

**Parley gap.** Parley has repeatedly needed one authority plus generated views, but its current
retro implementation diagnoses operational friction rather than terminology/authority drift.
**PRIMARY:** `IdeaSignals` contains rounds, review rounds, fix-up cycles, not-fixed/dismissed counts,
escalations, blocked/abandoned state, run failures, score, and failure type
(`internal/retro/retro.go:19-35`); its classifier covers blocked, escalation, runtime failure,
fix-up/review/design churn, and low friction (`internal/retro/retro.go:124-164`).

**Existing equivalent / is cognee better?** The separation already exists: §13 makes retrospective
optimization advisory and routes any change through a normal idea; playbooks are non-canonical
(**PRIMARY:** `COOPERATION.md:1156-1191`). For concrete namespaces, Parley is stronger and
deterministic: `agents.toml` is the authority, §2 is generated, and retired rows are marked rather
than deleted (**PRIMARY:** `COOPERATION.md:101-129`); §2-only IDs are reported `unmapped` and
`section2-only`, never silently adopted (**PRIMARY:** `internal/app/roster_view.go:52-90`). Cognee's
claimed fuzzy matching and merge are not better.

**Concrete proposal.** Extend `parley retro` with an optional, read-only terminology-drift diagnostic:
collect machine identifiers, frontmatter keys, command names, and explicitly declared aliases;
report candidate “two names / two authorities” clusters with exact locators; never fuzzy-merge,
rewrite, or assign an authority. Its only write remains the existing candidate-idea scaffold. This
is an ordinary tooling idea and adds no protocol bytes or dependency.

**Cost of being wrong.** False positives waste review time and can encourage premature
standardization. An auto-merge would be much worse: it could collapse distinct concepts. Keep the
output advisory, deterministic, and locator-backed; no 80% fuzzy threshold.

### 4. DROP — ontology grounding and `ontology_valid`

**Cognee mechanism.** **SECONDARY:** an optional RDF/OWL vocabulary, fuzzy matching at an asserted
80% cutoff, BFS enrichment, and an `ontology_valid` flag (`00-prompt.md:75-85`).

**Parley gap / equivalent.** The real Parley problem is competing authorities, not lack of a general
ontology. For roster identity, authority + generated view + explicit unmapped status already solves
the structural case deterministically (**PRIMARY:** `COOPERATION.md:101-129` and
`internal/app/roster_view.go:52-90`). **RECALL / UNVERIFIED:** other protocol nouns may still drift;
rank 3 is the bounded way to look for them.

**Concrete proposal.** None. Do not add RDF/OWL, BFS expansion, fuzzy matching, or a new normative
validity flag. **Cost of being wrong:** false equivalence silently merges distinct rules, while an
ontology becomes another authority to maintain and adds bytes to the measured cost centre.

### 5. DROP — feedback-weighted and “truth-subspace” reranking

**Cognee mechanism.** **SECONDARY:** finished-session feedback allegedly changes edge weights and
reranks retrieval against learned “truth” directions; the internal weighting is not exposed
(`00-prompt.md:87-91`).

**Parley gap / equivalent.** Parley already has the safe structural equivalent: retrospective
mining is advisory, changes pass through normal multi-agent acceptance, and automated loops may
discover candidates but may not promote, implement, merge, finalize, or change the roster without
a recorded gate (**PRIMARY:** `COOPERATION.md:1156-1191,1193-1229`). The implementation also mines
explicit failure signals rather than latent weights (**PRIMARY:** `internal/retro/retro.go:19-35,124-164`).
That is better for an auditable protocol.

**Concrete proposal.** None. Feedback may remain an input to an ordinary retro diagnosis; it must
not silently reweight which normative rules, objections, or reviewers are seen. **Cost of being
wrong:** correlated past errors become self-reinforcing “truth,” dissent becomes less retrievable,
and the system recreates the forbidden context-selector position.

### 6. DROP — temporal search as a protocol mechanism

**Cognee mechanism.** **SECONDARY:** Temporal Cognify extracts events/timestamps and exposes a
`TEMPORAL` search type, while the supplied docs expose no contradiction, supersession, or temporal
branching mechanism (`00-prompt.md:96-103`).

**Parley gap / equivalent.** Parley already records dated rounds, append-only signoffs, static
`FINAL.md`, and explicit verdict conflicts (**PRIMARY:** `COOPERATION.md:319-352,379-404,1263-1309`).
Its actual temporal gap—per-idea protocol-version pinning—is already ratified but not implemented
(**PRIMARY:** `COOPERATION.md:773-780`). Cognee's asserted timestamp search does not improve that
design.

**Concrete proposal.** None from cognee. Continue the existing version-pinning successor when its
priority warrants it. **Cost of being wrong:** treating “found by time” as “valid at time” hides
supersession and conflict rather than resolving them.

### 7. DROP — NodeSets, datasets/ACL, and graph/hybrid search types

**Cognee mechanism.** **SECONDARY:** NodeSets tag/group data; datasets organize documents and form
the permission unit; several graph/hybrid/triplet search modes are exposed
(`00-prompt.md:93-101`).

**Parley gap / equivalent.** Phase/track grouping is already the ratified phase-packet design with
an omission index and fail-open fallback (**PRIMARY:**
`meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md:35-63`). Files are canonical and
artifact ownership is already per participant (**PRIMARY:** `COOPERATION.md:60-70,97-99,727-743`).
**RECALL / UNVERIFIED:** I see no product requirement in this brief for tenant-style confidentiality
inside a deck.

**Concrete proposal.** None. Do not turn phase grouping into ACL or retrieval semantics. **Cost of
being wrong:** hidden context, permission drift, and a new service/tool in the normative path—the
exact class excluded by the binding prior decision.

## Concerns / open questions

1. For rank 1, can the receipt segment map be derived from the exact serialized byte buffer rather
   than from pre-serialization intent? If not, it should remain a prompt hash only, not claim
   source-level completeness.
2. Before rank 2, what precisely caused the v1.43.1 deletion: an unprovable semantic invariant, an
   implementation defect, or both? The successor must answer from the closed review record and code
   history; **RECALL / UNVERIFIED** here.
3. Rank 3 should prove that candidate terminology drift is common enough to justify code. A dry-run
   corpus report belongs in the idea brief; otherwise drop it without implementation.
4. None of these proposals justifies adding net bytes to `COOPERATION.md`. Rank 1 and rank 3 are
   ordinary audit/advisory tooling. Rank 2 stays implementation-scoped unless a later §7 SUCCESSOR
   explicitly pays for normative text by deleting superseded wording or demonstrates a measured
   benefit large enough to justify the addition.

## Risks

- **Audit theatre:** a context receipt can be internally consistent and still omit an unregistered
  input source. Shadow comparison and mutation tests are mandatory before trusting it.
- **Consent fabrication:** semantic forgetting is uniquely dangerous because an omission is treated
  as agreement. Full-history fallback must dominate every uncertain case.
- **Authority proliferation:** an ontology, alias table, ledger, or receipt can become another source
  of truth. Each proposed artifact must be derived/audit-only unless a successor explicitly assigns
  authority.
- **Protocol growth:** copying cognee vocabulary into normative prose would worsen the measured cost
  centre. Prefer existing event logs and retro code; add no rule merely to name a mechanism.
- **Vendor-claim inflation:** cognee's claimed capabilities are SECONDARY and unmeasured here. No
  recommendation depends on their effectiveness—only on whether the structural operation exposes a
  Parley gap.
