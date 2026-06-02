---
idea: 2026-06-02T12-07-14-meta-protocol-ch
drafted-by: claude
date: 2026-06-02
---

## Agreed decisions

The protocol gains an additive, opt-in **§12 "Pipeline blocks & action stages"**. A deck with no `pipeline.yaml` behaves exactly as Sections 0–11 define today; an existing idea is a valid degenerate one-block pipeline. Existing Phase 0–8 mechanics are unchanged.

1. **Block model.** A pipeline is an ordered list of **blocks**; each block is one invocation of the existing cooperation engine. Four kinds:
   - `deliberation` — runs Phase 1–4, produces one typed stage artifact.
   - `implementation` — runs Phase 5–8 unchanged, produces `IMPLEMENTATION.md` + code state.
   - `action` — runs Phase 1–4 for an action plan, then an explicit **driver-only `execute` sub-phase** after consensus + gate approval.
   - `watcher` — defines monitoring/alerting policy; breach handling opens a gated follow-up by default.
2. **Typed artifacts (not "FINAL.md by another name").** Stage artifacts keep markdown ergonomics but carry minimal **typed frontmatter**: `artifact_kind`, `pipeline_slug`, `block_id`, `derived_from[]`, `risk`, `providers_required[]`, `effects_intent`. Canonical names: `BUSINESS_SPEC.md`, `TECHNICAL_SPEC.md`, `IMPLEMENTATION_DESIGN.md`, `IMPLEMENTATION.md`, `DEPLOYMENT.md`, `RUNBOOK.md`, `MONITORING.md`. They remain consensus/finalize/transport-compatible (zero transport change).
3. **Execution boundary (hard rule).** Local CLI agents author markdown only (their own round/consensus/signoff/plan artifact). The **driver** is the sole actor that performs side effects, via a **provider-agnostic interface** with generic capability names (`deploy.preview`, `deploy.production`, `runtime.rollback`, `monitor.alert`, `issue.create`, `notify.send`). Vercel/Atlassian are first implementations behind that interface, never protocol dependencies.
4. **Seeding contract.** The driver authors block N+1's `00-prompt.md` as initiator-owned **kickoff** material from the manifest + block N's finalized typed artifact (`derived_from` paths + the next block contract). This is not a participant artifact, so canonical ownership is preserved. If any required input artifact is missing/stale/not-finalized, the block does not start.
5. **Durable state split.** `pipelines/<slug>/pipeline-run.json` is a **cursor/index only** (current block, completed blocks, pending gate, status, timestamps, effect references). Each side effect is its **own** file `pipelines/<slug>/effects/<key>.json`, semantically append-only, transitioning `planned → dry_run_ok → executing → succeeded|failed → reconciled|abandoned`, recording `external_ref` and appended `attempts[]`. Agents never mutate the ledger.
6. **Idempotency key.** `sha256(pipeline_slug | block_id | provider | capability | target | request_hash)`; stored filename is a stable digest prefix, full key recorded inside. `request_hash` hashes the normalized request body so a changed target/body becomes a new planned effect. **Retrying without an idempotency key is prohibited; resume must reconcile external state before retry.**
7. **Gates & autonomy (supervised-first).** Auto-advance is allowed **inside** a block only after normal quorum + signoff predicates pass. Block-boundary gates default to **human**. `risk: production` mutations (`*.production`, `*.rollback`) are **non-bypassable** regardless of autonomy. A per-pipeline `autonomy` flag may later auto-resolve only **low-risk, non-prod** left-half boundaries. Gate files reuse the existing HITL question/risk model.
8. **Capability dispatch halts, never degrades.** If no active participant satisfies a required advisory lens, or no provider satisfies a required action capability, the block **stops before consensus** and raises a gate. It must never silently fall back to solo or to an unqualified executor (the automation-safe form of the non-solo rule).
9. **Linear v1.** `blocks[]` is ordered; `edges[]` exists in the schema for future DAG continuity but a v1 driver **rejects any non-linear graph** as a validation error. Transport is **sticky** for the whole pipeline in v1.
10. **Monitoring loop-closure.** Watcher artifacts define signals, thresholds, destinations, breach fingerprints, dedupe windows. Breaches **notify + open a human gate** by default; auto-opening a remediation idea is allowed only for predeclared low-risk breach classes and uses the same sticky transport. Production remediation stays gated.
11. **Gate policy evaluator.** One **central evaluator** keyed by policy **names** (no embedded per-gate scripts), for auditability.

## Agreed trade-offs

- Strict "agents write markdown only" concentrates execution surface in the driver → higher harness-trust, but a narrow, testable, auditable executor and a clean ownership model. Accepted.
- Typed frontmatter costs a little schema work up front → bought deterministic seeding, validation, and resume. Accepted over the looser `stage:`-only option.
- Ratify the whole §12 vocabulary now (per owner decision: whole pipeline at once), but **stage the implementation**: define action/watcher semantics now, ship them initially validation-/dry-run-only, real provider mutation last.

## Open items deferred to implementation

- **First CLI increment** = linear, single-transport, supervised, **deliberation blocks** (BUSINESS_SPEC → TECHNICAL_SPEC → IMPLEMENTATION_DESIGN), proving manifest validation, boundary gates, prompt seeding, typed artifacts, durable cursor/resume. `implementation` is a declared block type; including it in the first increment depends on the existing Phase 5–8 runner being trivially safe to invoke. Then: dry-run-only `action` blocks (provider capability validation, per-effect ledger, reconcile-on-resume tests) → one real provider-backed action behind a mandatory human gate → `watcher` definitions + breach dedupe.
- **Dry-run-unsupported providers:** record explicitly, raise risk, block auto-approval; for production, a human gate may approve with the limitation visible.
- **`watcher` breach auto-open** stays in §12 as a future-compatible low-risk option; the shipped v1 driver implements **notify-and-gate only**.
- **Decider agent** is a future low-risk-only policy hook; block-and-wait remains the default for every unresolved consensus/gate conflict.
- Concrete `pipeline.yaml` schema fields, gate `policy` name catalog, and `MONITORING.md` `watch:` spec are implementation-idea scope.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-06-02
Status: ✅ ACCEPT
Notes: Preserves backward compatibility, keeps side effects centralized in the driver, v1 scope narrow enough to implement safely. Deferred items are implementation-scoped, not blockers for ratifying the protocol vocabulary.

### Signoff: hermes — 2026-06-02
Status: ✅ ACCEPT
Notes: §12 cleanly extends the existing engine with typed stages, driver-only execution, idempotent effect tracking; backward-compatible; markdown-only agent contract preserved. High-risk gates and capability-dispatch halts non-bypassable. Staging deliberation-first is pragmatic. No reservations.

### Signoff: claude — 2026-06-02
Status: ✅ ACCEPT
Notes: Faithfully captures round-02 convergence; adopts typed contracts, per-effect ledger, halt-not-degrade dispatch, non-bypassable prod gates, supervised-first. Whole vocabulary ratified now, implementation staged dry-run-first.
