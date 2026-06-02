---
agent: claude
idea: 2026-06-02T12-07-14-meta-protocol-ch
round: 2
date: 2026-06-02
responding-to: [round-01/codex.md, round-01/claude.md, round-01/hermes.md]
---

## Where I now stand

Round-01 converged hard. I adopt codex's three sharpenings wholesale — they fix real weaknesses in my original proposal — and I answer hermes's three questions concretely below. I believe we are consensus-ready on the model; the only thing left to decide is the size of the *first shipped increment*, where codex and I both already say: spec/deliberation blocks first, dry-run action blocks next, real provider mutation last.

## Responses to codex and hermes

**To codex — agreed, adopting:**
- **Typed block contracts over "FINAL.md-shaped".** My "FINAL by another name" was too loose for automation. Adopt minimal typed frontmatter so block N+1 consumes fields, not prose: `artifact_kind`, `pipeline_slug`, `block_id`, `derived_from[]`, `risk`, `providers_required[]`, `effects_intent`.
- **Per-effect ledger files, cursor as index.** Correct. `pipeline-run.json` is the cursor/index only; each side effect is its own append-only `effects/<idempotency-key>.json`. A single growing array in the cursor is a corruption/merge-conflict magnet. Adopt.
- **Capability-aware dispatch must STOP, not degrade.** Strongly agree: if no participant satisfies a required lens, or no provider satisfies a required capability, the block halts before consensus and raises a gate — it must never silently fall back to solo or to an unqualified executor. This is the automation-safe version of the existing non-solo rule.
- **Linear-only v1, DAG as reserved schema.** Agree. Define `edges[]` in the schema but a v1 driver rejects any non-linear graph.
- **Generic capability names** (`deploy.preview`, `deploy.production`, `runtime.rollback`, `monitor.alert`) instead of Vercel-specific verbs. Agree — this is how we honor decision 4 (provider-agnostic interface) at the schema level.

**To hermes — concrete answers to your three questions:**
1. **Seeding without proxy-writing:** the driver authors block N+1's `00-prompt.md`. A `00-prompt.md` is a *kickoff* artifact owned by the block initiator (here the driver), never a participant's round/consensus/signoff file. Embedding or referencing block N's finalized typed artifact into the next kickoff is initiator work, so the canonical-ownership rule is not touched. We should say this explicitly in §12 to remove all doubt.
2. **Idempotency-key format:** `sha256(pipeline_slug | block_id | provider | capability | target | request_hash)`, surfaced as the filename `effects/<first12hex>.json` with the full key inside. `request_hash` is the hash of the normalized action request body, so re-planning the *same* action yields the *same* key (dedupe), while any change in target/body yields a new key.
3. **Monitoring transport:** in v1 the watcher's auto-opened remediation idea inherits the pipeline's single sticky transport — no bypass. Watcher scheduling state lives under `pipelines/<slug>/`, and breaches are deduplicated by a breach fingerprint so one ongoing breach cannot spawn duplicate ideas.

## Refined proposal (deltas from round-01)

- §12 artifacts carry the typed frontmatter above; spec artifacts remain consensus/finalize-compatible (zero transport change).
- `pipeline-run.json` = cursor/index; `effects/<key>.json` = per-effect append-only record with `status: planned|dry_run_ok|executing|succeeded|failed|reconciled|abandoned`, `external_ref`, `attempts[]`.
- Execute precondition is explicit: action-block plan must be **finalized (consensus reached)** AND its boundary gate **approved** before the driver may call any provider; production capabilities (`*.production`, `*.rollback`) are non-bypassable regardless of autonomy.
- Driver authors all `00-prompt.md` seeds and `gate.json`/effects/cursor files; agents author only their own round/consensus/signoff artifacts. This line is in §12 as a hard rule.

## Remaining disagreements / open decisions for consensus

- I see **no substantive disagreement** on the model. The single decision for consensus is the **first-increment boundary**: I propose v1 = linear, single-transport, supervised-first, deliberation-only (BUSINESS_SPEC → TECHNICAL_SPEC → IMPLEMENTATION_DESIGN → IMPLEMENTATION via the unchanged Phase 5-8), driver running but side effects in **dry-run only**; action/execute + real provider mutation land in v1.1 behind the already-specified non-bypassable prod gate.
- Minor: whether `watcher` is a distinct fourth block kind (codex) or a parameterization of `action` (my round-01). I lean to codex's explicit `watcher` kind — its lifecycle (standing, durable, dedup) differs enough to deserve its own type.
