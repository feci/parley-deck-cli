---
agent: hermes
idea: 2026-06-02T12-07-14-meta-protocol-ch
round: 2
date: 2026-06-02
responding-to: [round-01/codex.md, round-01/claude.md, round-01/hermes.md]
---

## Where I now stand

Round-01 shows strong convergence on the additive §12 spine, typed artifacts, per-effect ledger, non-bypassable prod gates, and driver-only execution. Codex's typed contracts and linear v1 sharpen the proposal; Claude's emphasis on dry-run-first increment and seeding ownership clarification is the right safety-first cut. My three round-01 questions now have concrete answers from the round-02 artifacts, moving us to consensus.

## Responses to codex and claude

**To codex — agree on all points:**
- Adopt typed frontmatter (`artifact_kind`, `derived_from[]`, `providers_required[]`, `effects_intent`) over vague "FINAL.md-shaped".
- Per-effect `effects/<idempotency-key>.json` files + cursor index only; reject central array.
- Capability dispatch must halt (not degrade) when no matching participant/provider.
- Linear-only v1 with reserved DAG schema; generic capability names from day one.
- Smallest increment: deliberation blocks first, dry-run action next, real mutation last.

**To claude — agree and refine:**
- Explicit consensus boundary for execute: finalized plan artifact + approved gate required before driver acts.
- Seeding contract: driver authors `00-prompt.md` as initiator-owned kickoff (not participant file), embedding prior typed artifact.
- Idempotency: `sha256(pipeline_slug|block_id|provider|capability|target|request_hash)`.
- Monitoring: inherits sticky transport; auto-open only low-risk with dedup fingerprint.

Trade-off: driver now owns more surface (seeding, ledger, reconcile) but keeps agents strictly markdown-only, preserving ownership model.

## Refined proposal (deltas from round-01)

- Add explicit seeding rule in §12: driver writes block N+1 `00-prompt.md` from block N's finalized typed artifact; this is kickoff, not consensus artifact.
- Concrete idempotency key: `sha256(pipeline_slug|block_id|provider|capability|target|request_hash)` as filename prefix + full record.
- Monitoring watcher: new root idea stays under same transport; no bypass allowed in v1.
- First increment locked: linear supervised deliberation blocks (BUSINESS_SPEC→TECHNICAL_SPEC→IMPLEMENTATION_DESIGN→IMPLEMENTATION) + dry-run action only.

## Remaining disagreements

None on core model. Single open decision: exact first-increment scope (I align with codex/claude: deliberation + dry-run only). Watcher as separate kind vs parameterized action remains minor; codex's explicit kind is preferable for lifecycle differences.