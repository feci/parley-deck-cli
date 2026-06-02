---
agent: claude
idea: 2026-06-02T12-07-14-meta-protocol-ch
round: 1
date: 2026-06-02
---

## Summary

I support the composable-pipeline-manifest spine but want to cut its v1 surface hard. The single highest-value, lowest-risk move is **chained deliberation blocks** (business-spec → technical-spec → impl-design → implementation), which needs *no* protocol change at all — it is just N existing ideas linked by a manifest and a driver that seeds each block's `00-prompt.md` from the previous block's FINAL-shaped artifact. The right half (deploy/ops/monitor `action` blocks with a side-effecting `execute` sub-phase) is where the genuine protocol and safety novelty lives, and it should be specified now but shipped behind a flag, because it introduces the only truly new and dangerous concept: agents reasoning about actions the *driver* then executes against production. We should ratify the whole §12 model in this idea but stage the implementation so the executor and effects-ledger are proven on no-op/dry-run before any real mutation.

## Proposed approach

COOPERATION.md §12 "Pipeline blocks & action stages":

- **Block** = one invocation of the existing engine. Three kinds: `deliberation` (Phase 0-4 → a FINAL-shaped artifact), `impl` (Phase 5-8, unchanged), `action` (Phase 1-4 plan + one `execute` sub-phase).
- **Artifacts** are FINAL.md by another name with identical frontmatter plus `stage:` and optional `derived-from:`. This is the key simplification — transports, consensus, signoff, finalize need zero change. Stage names: BUSINESS_SPEC / TECHNICAL_SPEC / IMPLEMENTATION_DESIGN / IMPLEMENTATION / DEPLOYMENT / RUNBOOK / MONITORING.
- **Manifest** `pipelines/<slug>/pipeline.yaml`: ordered blocks, each with `kind`, `roster`, `role_lens`, `consumes`, `produces`, and an outgoing `gate: auto|human|policy`.
- **Driver state** `pipelines/<slug>/pipeline-run.json`: durable cursor + `action_ledger[]`. The driver is the only writer.
- **Gate** `gates/<edge-id>.gate.json`: reuses the existing `hitl.Question` risk model. `risk:high` (all prod mutations) is non-bypassable — it can never auto-resolve even in auto mode.
- **Effects ledger** `effects/<idempotency-key>.json`: append-only; records `dry_run` result and `external_ref`. On resume the driver MUST reconcile (look up `external_ref`) before re-attempting.
- **Execution boundary**: agents write markdown only; the driver executes side effects via a provider-agnostic run-action interface (Vercel is the first impl). This belongs in §12 as a hard rule, not a convention.

Auto-advance is allowed only *inside* a block (the round/consensus predicates are already machine-checkable). Block boundaries follow the manifest gate; "supervised-first" means the default manifest sets every boundary to `human` and a per-pipeline `autonomy: auto` later relaxes the left half only.

## Concerns / open questions

- **Executor correctness is the whole ballgame.** Promoting `runplan.Plan` from advisor to executor + durable resume + idempotent side effects is a distributed-systems problem. Idempotency keys reduce but never eliminate double-apply; reconcile-on-resume is mandatory and must be specified, not assumed. I'd gate the entire right half behind a dry-run-only mode until reconcile is tested against a provider that can crash mid-call.
- **`execute` sub-phase is a real protocol expansion**, not "Phase 1-4 plus a bit". It is the first time a Parley artifact has a side effect. It needs its own consensus precondition (the plan artifact must be finalized AND a human gate cleared) before the driver may act.
- **Capability-aware dispatch**: the manifest must cleanly separate "who deliberates" (agents) from "who executes" (driver). If this leaks, an agent could be asked to deploy, which breaks the safety model.
- **Monitoring auto-open loop** risks a feedback storm (breach → idea → deploy → breach). Default must be notify-and-gate; auto-open only for pre-declared low-risk classes, with a rate limit.
- **Slug/roster per block**: do specs and SRE share one roster? Different stages want different lenses; but multiple rosters multiply credential scope. I lean single roster + role-lens for v1.

## Risks

- **Over-build**: shipping the executor, effects-ledger, run-action layer, and §12 all at once is a large surface that could destabilize the proven Phase 0-8 engine. Mitigation: additive-only, the engine stays byte-for-byte; pipeline files live under a new dir; a deck with no manifest is unchanged.
- **False sense of automation**: "automatic" must never mean "unsupervised production mutation". The non-bypassable prod gate is the one rule we cannot compromise.
- **Reconcile gaps**: an action that succeeds externally but crashes before ledger write is the classic failure. Without mandatory reconcile this double-deploys. This is the top correctness risk.
- **Vendor lock creep**: wiring Vercel/Atlassian directly instead of behind an interface would violate vendor-neutrality. The interface must exist from day one even with one implementation.
- Smallest first increment I would ship: chained *deliberation* blocks only (specs → design → impl), driver in dry-run/advisory mode, no `action` blocks — proves the manifest + driver + boundary-gates with zero production risk.
