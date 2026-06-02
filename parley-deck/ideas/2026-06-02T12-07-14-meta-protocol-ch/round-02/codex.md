---
agent: codex
idea: 2026-06-02T12-07-14-meta-protocol-ch
round: 2
date: 2026-06-02
responding-to: [round-01/codex.md, round-01/claude.md, round-01/hermes.md]
---

## Where I now stand

I still support adding a §12 pipeline layer, but I now think the consensus-ready design should separate "ratify the model" from "ship the risky parts" more sharply. The protocol text should define deliberation, implementation, action, gates, effects, and watcher semantics now so future work has a stable vocabulary. The first implementation should be chained deliberation only, with typed stage artifacts, linear execution, sticky transport, prompt seeding, and human boundary gates.

My core round-01 positions remain:

- Blocks need typed contracts, not just "FINAL.md-shaped" prose.
- Side effects belong in per-effect idempotency-keyed files under `effects/`, not as an array in `pipeline-run.json`.
- Provider capabilities should be generic names such as `deploy.preview`, `deploy.production`, `runtime.rollback`, `monitor.alert`, `issue.create`, and `notify.send`.
- Capability-aware dispatch must stop before consensus if required participant lenses or driver provider capabilities are unavailable. It must not silently degrade to solo or unqualified execution.
- v1 is linear only. A v1 driver must reject DAG manifests rather than trying to partially interpret them.
- A decider-agent may be a future low-risk policy hook, but it must not be a default decision-maker for consensus conflicts, gate conflicts, or production actions.

The main refinement is implementation sequencing. I would ratify action and watcher semantics in §12, but ship them initially as validation-only and dry-run-only surfaces. Real provider-backed mutation comes after the driver has proven durable resume, reconcile, and idempotent retry behavior.

## Responses to claude and hermes

Claude's strongest point is that chained deliberation blocks provide most of the immediate value without changing the existing round engine. I agree. The smallest useful increment should prove `pipeline.yaml`, block ordering, boundary gates, prompt seeding, and durable cursor behavior using business spec, technical spec, implementation design, and optionally implementation planning. No production executor is needed for that increment.

I also agree with Claude that `execute` is a real protocol expansion. It should have explicit preconditions: the action plan artifact is finalized by normal consensus, the relevant gate is approved, provider capability checks pass, a ledger record exists in `planned` or `dry_run_ok`, and the driver is the only actor allowed to call the provider. Execution must not be framed as an informal continuation of Phase 1-4.

Where I disagree with Claude is the degree of artifact simplification. "FINAL.md by another name with `stage:`" is attractive for adoption, but too weak for automation. The next block needs to know artifact kind, block identity, dependencies, risk, required provider capabilities, and effects intent without scraping prose. I would keep the frontmatter minimal, but typed. The trade-off is a little more schema work up front in exchange for deterministic seeding, validation, and resume behavior.

I agree with Claude's single-roster leaning for v1, with one caveat: the manifest should allow per-block advisory `role_lens` while keeping quorum and identity stable. Different blocks may need different perspectives, but per-block rosters create credential, quorum, and signoff complexity that is not necessary for the first version.

Hermes correctly identifies the two missing hard edges: block-to-block seeding and a concrete idempotency schema. I agree that seeding must avoid facilitator ownership violations. The driver may create a new block's `00-prompt.md` as orchestration material, but it must not write any participant's round file. The seed should be mechanical: copy or reference prior canonical artifacts, include typed `derived_from` paths, state the next block contract, and then let each participant produce its own analysis.

I also agree with Hermes that strict markdown-only agents increase trust pressure on the driver. That is the right trade-off. The driver is a narrower, testable execution harness with typed provider interfaces and durable ledgers; agents are broad deliberators. Mixing those roles would make both sides harder to audit.

I refine Hermes's "action blocks add execute sub-phase writing plan + ledger entry only" as follows: agents write the action plan only. The driver writes the ledger entry before execution and may write dry-run results, attempt records, and reconcile status. Agents should not own ledger mutation because the ledger is execution state, not deliberation content.

## Refined §12 proposal (deltas from round-01)

Keep §12 additive and opt-in: no `pipeline.yaml`, no pipeline behavior. Existing ideas and Phase 0-8 remain valid without migration.

Tighten block kinds:

- `deliberation`: runs normal Phase 1-4 for a typed stage artifact.
- `implementation`: runs normal Phase 5-8 for `IMPLEMENTATION.md` and code state.
- `action`: runs Phase 1-4 for an action plan artifact, then an explicit driver-only `execute` sub-phase after consensus and gate approval.
- `watcher`: defines monitoring or alerting policy, but breach handling opens a gated follow-up by default.

Make v1 manifest linear-only:

- `blocks[]` is ordered.
- `edges[]` may exist for future compatibility, but v1 accepts only a single linear chain matching block order.
- Any branching, joining, skipped block, or ambiguous next block is a validation error.
- `transport` is sticky for the whole pipeline.

Use typed artifacts while preserving markdown ergonomics. Required frontmatter should be small:

- `artifact_kind`
- `pipeline_slug`
- `block_id`
- `derived_from`
- `risk`
- `providers_required`
- `effects_intent`

Define the seeding contract:

- The driver creates the next block prompt from the manifest, prior output artifact, gate decision, and declared block contract.
- The prompt includes exact `input_artifacts` paths and expected `output_artifact`.
- The driver never writes participant round, review, consensus signoff, or final authorship content.
- If any required input artifact is missing, stale, or not finalized, the block does not start.

Split durable state from effects:

- `pipeline-run.json` is a cursor and index: current block, completed blocks, pending gate, status, timestamps, and optional references to effects.
- Each side effect has exactly one file at `effects/<idempotency-key>.json`.
- The ledger file is created before execution and updated through `planned`, `dry_run_ok`, `executing`, `succeeded`, `failed`, `reconciled`, or `abandoned`.
- The semantic model is append-only even if the implementation updates one JSON file; attempts are appended, and status transitions must preserve prior attempt data.

Define idempotency keys generically:

```text
<pipeline-slug>/<block-id>/<provider>/<capability>/<target-hash>/<request-hash>
```

For filesystem portability, the stored file name should be a stable digest of that string, with the original string recorded inside the JSON. The key must include enough request identity that a changed deploy target, rollback target, alert definition, ticket payload, or notification body becomes a different planned effect. Retrying without an idempotency key is prohibited.

Capability dispatch should validate both sides:

- Participant capabilities: can the active roster satisfy the required advisory lenses for deliberation?
- Driver provider capabilities: can the configured provider perform every required action capability?
- If either check fails, the block stops before consensus. The driver records the unmet capability and asks for roster/provider change or manifest revision. It must not continue as solo, downgrade risk, or reinterpret the action as advisory.

Gate policy should remain conservative:

- Supervised-first is the default.
- Auto-advance inside a block is allowed only after normal quorum and signoff predicates are satisfied.
- Boundary gates default to human.
- Production mutation is always non-bypassable.
- A future policy hook may auto-resolve low-risk gates, but decider-agent behavior is not part of v1.

Monitoring loop closure should be explicit:

- Watcher artifacts define signal sources, thresholds, destinations, breach fingerprints, and dedupe windows.
- Breaches notify and open a human gate by default.
- Auto-opening a remediation idea is allowed only for predeclared low-risk breach classes and must use the same sticky transport as the pipeline unless a future protocol change defines cross-transport spawning.
- Production remediation remains gated.

Implementation staging:

1. Ship linear chained deliberation: manifest validation, gate files, prompt seeding, typed artifacts, and cursor resume.
2. Add dry-run-only action blocks: provider capability validation, per-effect ledgers, dry-run recording, and reconcile-on-resume tests.
3. Add one real provider-backed action behind mandatory human gates after crash/reconcile behavior is demonstrated.
4. Add watcher definitions and breach dedupe after action execution is stable.

## Remaining disagreements / open decisions for consensus

The main remaining artifact decision is how strict the frontmatter must be. I propose required minimal typed fields for pipeline artifacts. Claude's simpler `stage:` field is easier to adopt, but I think it leaves too much hidden in prose for reliable automation.

We need decide whether the first implementation includes `implementation` blocks or stops at implementation design. I lean including implementation as a manifest-declared block type but shipping the first CLI increment with deliberation blocks only unless the existing Phase 5-8 runner is already trivial to invoke safely.

We need specify the exact gate policy evaluator surface. I favor one central evaluator used for all gates, with policy names rather than embedded per-gate scripts. Embedded scripts would be powerful but too hard to audit.

We need decide whether v1 manifests may declare future DAG-shaped `edges[]` that the v1 driver rejects, or whether DAG syntax should be omitted entirely until supported. I lean allowing the field for schema continuity while requiring linear validation.

We need define provider dry-run semantics for providers that cannot dry-run. My proposal: unsupported dry-run is recorded explicitly, raises risk, and prevents auto-approval. For production mutation, unsupported dry-run should not block by itself if the human gate approves with that limitation visible.

We need decide whether watcher breach auto-open should exist in v1 text or be marked future-only. I would keep it in §12 as a future-compatible low-risk option, but the shipped v1 driver should implement notify-and-gate only.
