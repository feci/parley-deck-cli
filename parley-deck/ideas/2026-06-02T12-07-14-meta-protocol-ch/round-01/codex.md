---
agent: codex
idea: 2026-06-02T12-07-14-meta-protocol-ch
round: 1
date: 2026-06-02
---

## Summary

The proposed architecture is mostly sound: reusing the existing Phase 0-8 engine as the unit of composition preserves the current consensus mechanics, file ownership model, and non-solo discipline. The critical boundary is also right: agents deliberate and write markdown; the driver alone performs side effects through typed provider interfaces after gates.

The main weakness is that the proposal treats several very different things as "blocks" without enough type discipline. A deliberation block, implementation block, deployment action, operations change, and monitoring watcher have different inputs, outputs, failure modes, and replay semantics. The protocol should keep the round/consensus engine unchanged, but add strict typed block contracts around it.

My recommendation is to ratify the whole pipeline, but ship it behind a small, conservative v1: linear manifests only, one sticky transport, supervised-first gates at every boundary, provider interfaces with dry-run support, and append-only effects records. DAGs, per-block transports, decider agents, and automatic remediation loops should be defined as future-compatible fields, not required in the first implementation.

## Proposed approach

Add `COOPERATION.md` §12, "Pipeline blocks & action stages", as an additive opt-in layer:

```text
12. Pipeline blocks & action stages

12.1 Scope
When `parley-deck/pipelines/<slug>/pipeline.yaml` exists, the idea may be executed as a pipeline. If no pipeline manifest exists, the deck behaves exactly as Sections 0-11 define. Each pipeline block invokes the existing cooperation engine; the pipeline layer only sequences blocks, records gates, and executes approved side effects through the driver.

12.2 Block model
A pipeline is an ordered list of blocks. v1 pipelines are linear; future manifests may declare DAG edges, but a v1 driver must reject a non-linear graph unless it explicitly supports it.

Block kinds:
- `deliberation`: runs Phase 1-4 and produces one consensus artifact.
- `implementation`: runs Phase 5-8 and produces `IMPLEMENTATION.md` plus the implementation branch/worktree state.
- `action`: runs Phase 1-4 to approve an action plan, then enters `execute`, a driver-only sub-phase that performs approved side effects.
- `watcher`: installs or updates monitoring definitions and defines breach handling policy; any remediation work starts as a new gated root idea unless predeclared low-risk policy allows auto-open.

12.3 Canonical pipeline artifacts
The default pipeline artifacts are:
- `ideas/<slug>/BUSINESS_SPEC.md`
- `ideas/<slug>/TECHNICAL_SPEC.md`
- `ideas/<slug>/IMPLEMENTATION_DESIGN.md`
- `ideas/<slug>/IMPLEMENTATION.md`
- `ideas/<slug>/DEPLOYMENT.md`
- `ideas/<slug>/RUNBOOK.md`
- `ideas/<slug>/MONITORING.md`

Each stage artifact uses normal `FINAL.md`-style frontmatter plus optional typed fields:
- `artifact_kind`: business-spec | technical-spec | implementation-design | implementation | deployment | runbook | monitoring
- `pipeline_slug`: string
- `block_id`: string
- `derived_from`: list of prior artifact paths
- `risk`: low | normal | high | production
- `providers_required`: list of provider capability names
- `effects_intent`: none | planned | executed

12.4 Manifest
`parley-deck/pipelines/<slug>/pipeline.yaml` declares:
- `schema_version`
- `idea_slug`
- `autonomy`: supervised-first | auto-left | custom
- `transport`: local-dir | github-pr | gitlab-mr
- `participants`
- `blocks[]`: id, kind, stage, role_lens, input_artifacts, output_artifact, risk, provider_capabilities, gate_policy
- `edges[]`: from, to, gate_id

The transport is sticky for the whole pipeline in v1.

12.5 Durable run state
`parley-deck/pipelines/<slug>/pipeline-run.json` is the durable cursor. It records:
- `schema_version`
- `pipeline_slug`
- `status`
- `current_block`
- `completed_blocks[]`
- `pending_gate`
- `action_ledger[]`
- `created_at`, `updated_at`
- optional driver metadata

The driver must be able to resume using only the manifest, canonical artifacts, gate files, effects files, and `pipeline-run.json`.

12.6 Gates
Every edge between blocks has a typed gate file at `parley-deck/pipelines/<slug>/gates/<edge-id>.gate.json`.

Gate schema:
- `id`
- `pipeline_slug`
- `edge`
- `from_block`
- `to_block`
- `risk`: low | normal | high | production
- `status`: open | approved | rejected | expired | superseded
- `prompt`
- `details`
- `default_answer`
- `policy`
- `approved_by`
- `created_at`
- `answered_at`
- `decision_artifacts[]`

Supervised-first means the driver pauses at every block boundary until the gate is approved. `auto-left` may auto-approve non-production gates from idea through implementation only when policy explicitly allows it. Production-mutating gates are non-bypassable and always require human approval.

12.7 Execute sub-phase
In an action block, agents produce and sign off on a markdown plan. The driver then performs `execute` only after the matching gate is approved. Agents never invoke deploy, ops, ticketing, monitoring, or production mutation tools directly as protocol work.

12.8 Effects ledger
Every side effect is recorded before execution under `parley-deck/pipelines/<slug>/effects/<idempotency-key>.json`.

Effects schema:
- `idempotency_key`
- `pipeline_slug`
- `block_id`
- `action_type`
- `provider`
- `target`
- `risk`
- `status`: planned | dry_run_ok | executing | succeeded | failed | reconciled | abandoned
- `gate_id`
- `request_hash`
- `dry_run_result`
- `external_ref`
- `attempts[]`
- `created_at`, `updated_at`

On resume, the driver must reconcile any `executing`, `failed`, or externally ambiguous effect before retrying. Retrying without an idempotency key is prohibited.

12.9 Monitoring loop closure
Monitoring artifacts define SLOs, signal sources, alert destinations, and breach classes. A breach opens a human gate by default. Auto-opening a remediation idea is allowed only for predeclared low-risk breach classes; production remediation remains gated.

12.10 Compatibility
All pipeline files are optional. Existing ideas, rounds, consensus files, reviews, and implementations remain valid without migration. Existing `run.json` or manifests may gain optional pipeline fields, but old drivers may ignore them.
```

Minimal CLI capabilities required:

- Manifest parser and validator for `pipeline.yaml`, initially linear-only.
- Block runner that can invoke Phase 1-4, Phase 5-8, and round-02+ continuation for a named block.
- Durable executor that reads and writes `pipeline-run.json` and resumes from disk.
- Typed gate primitive using the existing HITL question risk model, with policy evaluation.
- Effects ledger writer plus reconcile-on-resume for ambiguous side effects.
- Provider-neutral action layer: deploy, runtime, ops, monitoring, notification, and issue/remediation providers behind interfaces.
- Capability-aware dispatch that matches block requirements to available local CLI agents and driver providers.
- Seeding layer that converts one block's typed output into the next block's `00-prompt` without changing participant file ownership.

Safety model:

- All mutating actions require a gate; production mutation always requires human approval.
- Agents write markdown only; the driver executes side effects.
- Every side effect has a stable idempotency key derived from pipeline slug, block id, provider, target, action type, and request hash.
- Dry-run is mandatory when a provider supports it; unsupported dry-run must be explicit in the action plan and raises risk.
- The ledger is append-only in semantics even if represented as one JSON file per effect.
- Resume first reconciles external state, then decides whether to retry, mark succeeded, or ask for human judgment.
- Provider integrations must be optional and replaceable; Vercel can be the first deploy provider, not a protocol dependency.

Migration and compatibility:

- No migration is required for existing decks.
- An existing idea is a degenerate one-block pipeline only when a manifest is added.
- Existing Phase 0-8 commands continue to work unchanged.
- New CLI commands should sit beside current commands, for example `pipeline validate`, `pipeline start`, `pipeline continue`, `pipeline gate`, `pipeline execute`, and `pipeline reconcile`.
- Existing run files may add optional fields but must not require old artifacts to be rewritten.

Smallest first increment:

Ship a supervised, linear, local-dir pipeline with three deliberation blocks: business spec, technical spec, and implementation design. Add manifest validation, gate files, durable cursor, and prompt seeding. Do not ship real deployment yet. Then add a dry-run-only action block and effects ledger. Only after that, enable one real provider-backed deployment action behind a mandatory production gate.

## Concerns / open questions

The term "FINAL.md-shaped artifact" is useful for reuse but too vague for automation. Each stage needs a small typed contract so the next block can consume it without scraping prose.

The proposed action block "Phase 1-4 plan + execute" needs a clear consensus boundary. Execution should not start from a draft plan or partial agreement; it starts only after the action plan has consensus and the gate is approved.

The effects ledger must be per-effect, not only an array in `pipeline-run.json`. A central cursor is too easy to corrupt, merge-conflict, or rewrite. Keep `pipeline-run.json` as an index/cursor and store each effect as its own idempotency-keyed file.

DAG support is premature for v1. Define enough schema to avoid painting the protocol into a corner, but require linear execution until the executor can prove deterministic gate ordering and resume behavior.

The "decider agent" default is risky. I would allow it only as a future policy hook for low-risk gates and never as a default behavior. Human block-and-wait should be the default for all unresolved consensus or gate conflicts.

Monitoring needs clearer ownership. A standing watcher block implies ongoing state outside a single cooperation round. The protocol should define whether the driver owns watcher scheduling, whether alert state is durable, and how duplicate breach ideas are suppressed.

Capability-aware dispatch needs a failure mode. If no local CLI participant can satisfy a required lens or provider capability, the block should stop before consensus rather than silently degrading to solo or unqualified execution.

## Risks

- Automation may blur deliberation and execution. Mitigation: keep consensus, gate approval, and execute as separate recorded states.
- Side effects may duplicate after restart. Mitigation: idempotency keys, per-effect ledgers, and mandatory reconcile before retry.
- Typed schemas may become too heavy and slow adoption. Mitigation: require only minimal frontmatter fields and keep rich detail in markdown.
- Provider interfaces may leak Vercel-specific assumptions. Mitigation: name generic capabilities such as `deploy.preview`, `deploy.production`, `runtime.rollback`, and `monitor.alert`.
- Monitoring loop closure may create noisy remediation ideas. Mitigation: default notify-and-gate, deduplicate by breach fingerprint, and allow auto-open only for declared low-risk breach classes.
- Auto-advance could weaken non-solo consensus. Mitigation: auto-advance only after normal quorum and signoff have completed; it must not replace participant artifacts.
- Pipeline state may conflict with existing transport mechanics. Mitigation: make pipeline state additive and let canonical stage artifacts remain under `ideas/<slug>/`.
- The first increment could be too ambitious if it includes production deployment. Mitigation: ship typed spec blocks and dry-run action blocks before any real provider mutation.
