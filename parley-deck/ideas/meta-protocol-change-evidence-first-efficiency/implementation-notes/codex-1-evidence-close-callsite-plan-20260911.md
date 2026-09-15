---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: integration call-site plan
not-a-signoff: true
---

# Actual verifier execution at closure — source boundary

Directly inspected Kimi's owned driver_evidence.go, driver_checks.go and
internal/evidence/evidence.go at 89a4305 plus the uncommitted parser/count fixes;
also the current integration driver_impl.go and internal/driver/impl.go.
No new closure command or runtime attestation is implemented by this plan.

## Existing behavior and exact gap

The production close branch in internal/driver/impl.go refreshes RunChecks for
list-form contracts, then optionally invokes GoalCheck for auto/strict ideas,
and finally calls Complete. Kimi's runChecksContract records the implementer's
actual criterion executions and tree digest before/after; EvidenceCloseGate can
validate a persisted independent rerun, but nothing calls AttestExecution on an
actual verifier execution in production. A second supplied verifier name cannot
fill that gap. GoalCheck's independent agent and exact verdict parsing also do
not substitute for typed criterion evidence. List-form contracts need the typed
gate on every closure, not only auto/strict.

## Integration requirements

1. Preserve the existing non-implementer reviewer/drafter selection. Refuse an
   absent/self/unavailable verifier. Keep the signed participant identity; a
   helper is not a new participant or a vote.
2. Give that actual verifier invocation an explicit runtime verification
   command/helper which resolves the original named contract, executes every
   criterion through RunCriterion, retains before/after tree digests and exact
   command hashes, and calls AttestExecution with those executions. Merely
   running the same commands in the facilitator and changing the actor string
   is not an independent agent execution path.
3. The helper must consume current frozen commands from the contract, not
   arbitrary shell text supplied by model output or a pasted report. Agent
   launch context supplies attribution; it is not human authentication or
   protection against a malicious same-UID process. Do not invent that claim.
4. Keep the report/tree/extra IMPLEMENTATION remainder bindings through the
   entire verifier run, including evidence commit changes. Changes in code,
   contract scope or command identity after the run invalidate closure. Runtime
   artifacts stay ignored; exact evidence exclusions remain defined.
5. Refuse writes/failed commands, skipped or zero cases, malformed streams,
   unknown shell semantics, self/stale/partial scope and absent helper execution.
   Preserve failed attempt evidence; don't replace it with a fabricated pass.
6. Wire the resulting EvidenceCloseGate immediately before Complete after the
   last operation that can alter the tree. An ordinary GOAL-CHECK: PASS cannot
   override its refusal. Re-entry must not reuse stale attestation silently.
7. Regression cases must exercise the real production close call path and a
   real fixture verifier process invoking the helper, not just a report with
   two manually assigned names. Cover a successful distinct verifier, self,
   missing invocation/helper, stale/changed tree, masked package failure, skip,
   partial contract and failed evidence persistence.

The shared app/driver integration is Codex-owned. Kimi owns the evidence package
and its existing app helper files; any required API refinement goes back to its
owner. Before adding a new app command/source file, record its non-overlapping
claim in IMPLEMENTATION.md. All budget/launch policy requirements still compose
with this path. This plan is not an AC-E1/E2 completion claim.
