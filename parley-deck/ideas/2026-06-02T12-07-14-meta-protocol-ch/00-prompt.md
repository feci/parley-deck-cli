---
idea: 2026-06-02T12-07-14-meta-protocol-ch
author: user
created: 2026-06-02
participants: [codex, claude, hermes]
status: final
---

## Problem / idea

Meta-protocol-change: evolve Parley Deck into a full automatic idea-to-monitoring pipeline.

We want to extend the Parley Deck cooperation protocol so that, from a single raw idea, the deck can drive the work all the way through to a deployed, operated, and monitored system — automatically, using the available local CLI agents, while preserving everything that already works. The target stages are: idea -> business specification -> technical specification -> implementation design -> implementation -> deployment -> operations -> monitoring. This idea runs as a §7 meta-protocol-change: its FINAL.md must propose the concrete COOPERATION.md additions (a new §12) plus the new canonical artifacts and CLI capabilities. Do NOT edit COOPERATION.md in this idea; only design and reach consensus on the change.

You are being asked for independent Phase-1 analysis. A proposed architecture is given below as the starting point. Critique it hard, surface flaws and risks, and propose concrete protocol text and artifact schemas. Disagreement and counter-proposals are welcome; do not rubber-stamp.

=== AGREED DIRECTION (locked by the project owner) ===
1. Scope: ratify the WHOLE pipeline at once — specifications through monitoring, including the side-effecting deploy/ops/monitor stages, the effects-ledger, and the new execute sub-phase.
2. Autonomy: "supervised-first" — fully automatic INSIDE each block (rounds -> consensus -> finalize auto-advance), but the driver PAUSES at every block boundary for human approval until the system is trusted; a per-pipeline flag later flips the left half (idea..implementation) to auto. Anything that mutates production always gates regardless of autonomy level.
3. Side-effect execution boundary: local CLI agents only WRITE markdown (their own artifact). All side-effecting actions (deploy, ops changes, monitor setup) are executed by the driver/harness via MCP tools behind an explicit gate — agents deliberate the plan, the driver acts.
4. Deploy provider: a provider-agnostic deploy/runtime interface from day one (Vercel is the first concrete implementation behind it).

=== PROPOSED ARCHITECTURE (the starting point to critique) ===
Composable Pipeline Manifest (spine) + staged gate-object/effects-ledger discipline. Keep today's Phase 0-8 rounds+consensus engine UNCHANGED and wrap it in a thin composition layer.

- New artifact pipelines/<slug>/pipeline.yaml lists ordered/DAG BLOCKS. Each block is exactly one invocation of today's engine:
  - deliberation block = Phase 0-4, produces a FINAL.md-shaped artifact under a stage-specific name;
  - impl+review block = Phase 5-8 (unchanged), produces IMPLEMENTATION.md + branch;
  - action block = Phase 1-4 plan + ONE new side-effecting `execute` sub-phase, produces a plan artifact + an execution ledger section.
- Stage -> artifact mapping (all spec artifacts are FINAL.md by another filename, identical frontmatter, so transports/consensus need zero change):
  - business-spec -> ideas/<slug>/BUSINESS_SPEC.md (role-lens: product-analyst)
  - technical-spec -> ideas/<slug>/TECHNICAL_SPEC.md (derived-from: BUSINESS_SPEC.md; role-lens: architect)
  - implementation-design -> ideas/<slug>/IMPLEMENTATION_DESIGN.md (role-lens: tech-lead)
  - implementation -> ideas/<slug>/IMPLEMENTATION.md + code branch (Phase 5-8 unchanged)
  - deployment -> ideas/<slug>/DEPLOYMENT.md + "## Execution Ledger" (role-lens: SRE/release)
  - operations -> ideas/<slug>/RUNBOOK.md (role-lens: SRE/ops)
  - monitoring -> ideas/<slug>/MONITORING.md (SLOs + watch: spec; standing watcher block; role-lens: observability)
- New non-idea artifacts: pipelines/<slug>/pipeline-run.json (durable cursor + action_ledger[]), pipelines/<slug>/gates/<edge-id>.gate.json (typed promote-to-next-block approval reusing the existing HITL hitl.Question risk model), pipelines/<slug>/effects/<idempotency-key>.json (side-effect ledger: external_ref + dry-run result).
- Durable driver: promote today's read-only continue/runplan.Plan into an executor that loads pipeline-run.json, runs the recommended action, seeds block N+1's 00-prompt from block N's typed output, and persists the cursor. Must survive restart purely from pipeline-run.json.
- Safety: every production-mutating action runs only behind a gate.json and is recorded with an idempotency key in an append-only effects ledger so a restart can never double-deploy; on resume, reconcile any action that may have succeeded externally before re-attempting.
- Monitoring closes the loop: on SLO breach the watcher opens a new root idea (gating on remediation).
- Backward compatibility: fully additive and opt-in. A deck with no pipeline.yaml behaves exactly as today; an existing idea is a valid degenerate one-block pipeline; pipeline files live under a new parley-deck/pipelines/ dir; existing run.json/manifest gain only optional fields.

=== DEFAULTS for the remaining open points (override if you disagree, with reasons) ===
- Deadlock on strict unanimity: block-and-wait for the human by default; an always-available decider agent may break ties ONLY for explicitly low-risk, non-prod block boundaries via policy.
- Transport: single sticky transport for the whole pipeline in v1 (revisit per-block transport later).
- Monitoring loop-closure: notify-and-gate for all SLO breaches initially; auto-open a remediation idea only for pre-declared low-risk breach classes.

=== WHAT TO DELIVER IN YOUR ROUND-01 ARTIFACT ===
- Your independent assessment of the proposed architecture: what is sound, what is wrong or risky, what is missing.
- Concrete proposal for COOPERATION.md §12 "Pipeline blocks & action stages": the phase/block model, the canonical artifacts and their frontmatter, the gate and effects-ledger semantics, and how auto-advance + HITL gates work.
- The minimal CLI capability set required (manifest parser, round-02+/Phase 5-8 runner, durable executor, run-action layer for deploy/ops/monitor, typed gate primitive + policy, effects ledger + reconcile-on-resume, capability-aware dispatch).
- The safety model for side-effecting actions (idempotency, dry-run, reconcile, non-bypassable prod gates).
- The migration/backward-compatibility story for existing decks and ideas.
- Your top risks and the smallest first increment you would ship.

=== CONSTRAINTS ===
- Preserve Parley Deck's vendor-neutral design; do not make any single vendor (Claude, Vercel, Atlassian, OpenAI) a required dependency. Provider integrations sit behind interfaces.
- Preserve canonical file ownership: each participant writes only its own protocol artifact; the facilitator must not proxy-write another participant's content.
- Preserve the non-solo execution requirement and the consensus/signoff mechanics.
- Changes must be additive and backward-compatible: existing Phase 0-8 decks and ideas keep working with zero migration.
- Local CLI agents write markdown only; side-effecting actions are driver/harness-executed behind gates.
- All files under parley-deck/ are English-only.
- Keep the change implementable incrementally; identify what needs no protocol change vs what genuinely requires this meta-change.

=== NON-GOALS ===
- Do not implement code in this idea; this is a design + consensus idea (Phase 0-4).
- Do not replace the round/consensus model; reuse it as the engine for every block.
- Do not add hidden autonomous execution: production mutations always gate on a human by default.
- Do not couple the protocol to one CI/CD, ticketing, or model vendor.

## Constraints

- Local-directory transport for this initial run.
- Human-in-the-loop mode by default unless the run was started with auto mode.

## Non-goals

- Do not make unrelated repository changes.
