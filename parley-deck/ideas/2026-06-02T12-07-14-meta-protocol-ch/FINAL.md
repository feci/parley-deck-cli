---
idea: 2026-06-02T12-07-14-meta-protocol-ch
status: final
author: claude
consensus-date: 2026-06-02
participants: [codex, claude, hermes]
---

## Final plan / specification

### Goal

Extend the Parley Deck protocol so a single idea can be driven — automatically, by the available local CLI agents, with human gates where it matters — from idea through business spec, technical spec, implementation design, implementation, deployment, operations, and monitoring. This is achieved by **composing the existing, unchanged Phase 0–8 cooperation engine into blocks** under a new additive, opt-in `COOPERATION.md` **§12**.

### Scope

Ratified now (the protocol vocabulary): the §12 block model, typed stage artifacts, the driver-only execution boundary with a provider-agnostic interface, the seeding contract, durable cursor + per-effect ledger, idempotency + reconcile, supervised-first gates with non-bypassable production mutation, capability-dispatch-halts-not-degrades, linear v1, and monitoring loop-closure semantics. Implementation is **staged** (see below); nothing here changes Sections 0–11 or any existing deck.

### Implementation details

The following is the proposed `COOPERATION.md` §12, ready to apply (the drafter applies it to the canonical protocol per §7, then logs `meta/protocol-changelog.md`). NOTE: the canonical protocol text lives in `parley-deck-skill/references/COOPERATION.md`; this dogfood deck's `COOPERATION.md` is updated in lockstep.

```markdown
## 12. Pipeline blocks & action stages

§12 is additive and opt-in. If `parley-deck/pipelines/<slug>/pipeline.yaml` does not exist, the deck behaves exactly as Sections 0–11 define. Every existing idea is a valid degenerate one-block pipeline; no migration is required.

### 12.1 Block model
A pipeline is an ordered list of **blocks**. Each block is exactly one invocation of the existing cooperation engine; the pipeline layer only sequences blocks, records gates, and (for action blocks) executes approved side effects through the driver. Block kinds:
- `deliberation` — runs Phase 1–4, produces one typed stage artifact.
- `implementation` — runs Phase 5–8 unchanged, produces `IMPLEMENTATION.md` plus code state.
- `action` — runs Phase 1–4 to reach consensus on an action plan, then enters `execute`, a driver-only sub-phase that performs approved side effects.
- `watcher` — defines monitoring/alerting policy; breach handling opens a gated follow-up by default.

### 12.2 Canonical stage artifacts
Stage artifacts are normal consensus/finalize-compatible markdown with added typed frontmatter (transports and the round/consensus mechanics are unchanged). Names: `BUSINESS_SPEC.md`, `TECHNICAL_SPEC.md`, `IMPLEMENTATION_DESIGN.md`, `IMPLEMENTATION.md`, `DEPLOYMENT.md`, `RUNBOOK.md`, `MONITORING.md`. Required typed frontmatter:
`artifact_kind`, `pipeline_slug`, `block_id`, `derived_from[]`, `risk` (low|normal|high|production), `providers_required[]`, `effects_intent` (none|planned|executed).

### 12.3 Manifest (`pipelines/<slug>/pipeline.yaml`)
`schema_version`, `idea_slug`, `autonomy` (supervised|auto-left), `transport` (sticky for the whole pipeline in v1), `participants`, `blocks[]` (id, kind, stage, role_lens, input_artifacts, output_artifact, risk, provider_capabilities, gate_policy), `edges[]` (reserved for future DAG; a v1 driver MUST reject any non-linear graph as a validation error).

### 12.4 Execution boundary (hard rule)
Local CLI agents author markdown only — their own round/consensus/signoff/plan artifact. The **driver** is the sole actor that performs side effects, through a **provider-agnostic interface** using generic capability names (`deploy.preview`, `deploy.production`, `runtime.rollback`, `monitor.alert`, `issue.create`, `notify.send`). Provider integrations (e.g. Vercel, Atlassian) are optional implementations behind that interface and are never protocol dependencies.

### 12.5 Seeding contract
The driver authors block N+1's `00-prompt.md` as initiator-owned kickoff material, built from the manifest, the prior block's finalized typed artifact (`derived_from` paths), the gate decision, and the next block contract (exact `input_artifacts` and expected `output_artifact`). A `00-prompt.md` is not a participant artifact, so canonical ownership (the facilitator/driver never writes a participant's round/consensus/signoff/final content) is preserved. If any required input artifact is missing, stale, or not finalized, the block does not start.

### 12.6 Durable state and effects ledger
`pipelines/<slug>/pipeline-run.json` is a cursor/index only: current block, completed blocks, pending gate, status, timestamps, and references to effects. Each side effect is its own file `pipelines/<slug>/effects/<key>.json`, semantically append-only, transitioning `planned → dry_run_ok → executing → succeeded|failed → reconciled|abandoned`, recording `external_ref` and an appended `attempts[]`. Agents never mutate the ledger; the driver is its only writer.

### 12.7 Idempotency and reconcile
Idempotency key = `sha256(pipeline_slug | block_id | provider | capability | target | request_hash)`, where `request_hash` is over the normalized request body. The stored filename is a stable digest prefix; the full key is recorded inside. Retrying without an idempotency key is prohibited. On resume the driver MUST reconcile external state (look up `external_ref`) for any `executing`/`failed`/ambiguous effect before retrying. Where a provider cannot dry-run, this is recorded explicitly, raises risk, and blocks auto-approval; for production a human gate may approve with the limitation visible.

### 12.8 Gates and autonomy (supervised-first)
Auto-advance is permitted only INSIDE a block, and only after the normal quorum + signoff predicates pass. Block-boundary gates default to `human`. Gate files (`pipelines/<slug>/gates/<edge-id>.gate.json`) reuse the HITL question/risk model and are resolved by one central policy evaluator keyed by policy names (no embedded per-gate scripts). `risk: production` mutations (`*.production`, `*.rollback`) are non-bypassable regardless of autonomy. A per-pipeline `autonomy: auto-left` flag may auto-resolve only low-risk, non-production left-half boundaries. A decider agent, if ever configured, is a future low-risk-only hook; block-and-wait is the default for every unresolved consensus or gate conflict.

### 12.9 Capability dispatch
Before consensus, the driver validates both sides: can the active roster satisfy the required advisory `role_lens`, and can the configured provider satisfy every required action capability? If either check fails, the block STOPS before consensus and raises a gate for a roster/provider/manifest change. It must never silently degrade to a solo run or to an unqualified executor — this is the automation-safe form of the non-solo execution requirement.

### 12.10 Execute sub-phase (action blocks)
Agents produce and reach consensus on a markdown action plan. The driver may call a provider only after: the plan artifact is finalized, the boundary gate is approved, provider-capability checks pass, and a ledger record exists in `planned` or `dry_run_ok`. Execution is never an informal continuation of Phase 1–4.

### 12.11 Monitoring loop-closure
`MONITORING.md` defines signal sources, thresholds, destinations, breach fingerprints, and dedupe windows. A breach notifies and opens a human gate by default. Auto-opening a remediation idea is allowed only for predeclared low-risk breach classes and uses the same sticky transport as the pipeline; production remediation remains gated. Breaches are deduplicated by fingerprint so one ongoing breach cannot spawn duplicate ideas.

### 12.12 Compatibility
All pipeline files are optional and live under `parley-deck/pipelines/<slug>/`; `ideas/`, `inbox/`, `meta/`, `runs/` are unchanged. Existing `run.json`/manifests may gain optional `pipeline_slug`/`block_id` fields under a schema bump with zero-value defaulting; older drivers ignore unknown fields and degrade to advisory.
```

### Tests

Protocol-level acceptance (to be enforced by the implementing ideas):
- A deck with no `pipeline.yaml` produces byte-identical behavior to Sections 0–11 (regression).
- Manifest validation rejects non-linear `edges[]` in v1.
- The driver never writes a participant's round/consensus/signoff/final file (ownership test); it may write `00-prompt.md`, `gate.json`, `pipeline-run.json`, and `effects/*.json`.
- A side effect interrupted between external success and ledger write is reconciled (not double-applied) on resume.
- `risk: production` gates cannot auto-resolve under any autonomy setting.
- Capability shortfall halts before consensus rather than degrading to solo.

### Non-goals

- No code is implemented in this idea (design + consensus only).
- The round/consensus model is not replaced; it is the engine for every block.
- No hidden autonomous execution: production mutations always gate on a human by default.
- The protocol is not coupled to any single CI/CD, ticketing, or model vendor.

### Verification

Unanimous Phase-3 signoff: codex ✅, hermes ✅, claude ✅ (see `consensus.md`). agy excluded from this idea (headless print-mode produced no artifact across three attempts; see `inbox/claude-to-all_..._agy-excluded.md`). Three active non-solo participants; quorum satisfied.

Implementation is staged as follow-on ideas:
1. Linear, single-transport, supervised **deliberation** pipeline (`BUSINESS_SPEC → TECHNICAL_SPEC → IMPLEMENTATION_DESIGN`): manifest parser/validator, boundary gates, prompt seeding, typed artifacts, durable cursor + resume. (No production executor.)
2. **Round-02+/Phase 5–8 runner** so blocks can advance past round-01 automatically.
3. Dry-run-only **action** blocks: provider-capability validation, per-effect ledger, reconcile-on-resume tests.
4. One real provider-backed action (Vercel deploy) behind a mandatory production gate.
5. `watcher` definitions + breach dedupe (notify-and-gate only).

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
- Owner decisions: scope=whole-pipeline; autonomy=supervised-first; side-effects=driver-only; deploy=provider-agnostic interface.
- Roster note: ./../../inbox/claude-to-all_2026-06-02T12-07-14-meta-protocol-ch_agy-excluded.md
