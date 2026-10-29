---
agent: codex-1
idea: meta-protocol-change-credit-exhaustion-failover
round: 1
date: 2026-10-29
---

## Summary

Add a versioned provider-exhaustion trigger to the existing membership system. A qualified first failure must prevent another launch for that agent.
Resolve affected roles from the remaining participants in a deterministic order. Preserve independent review, historical dissent, and the usable-participant floor.

This file precedes my access to peer round-01 files. I organize and implement under the owner's brief.

## Proposed approach

### Evidence and retry

Define provider exhaustion as an explicit provider or gateway report that the account's credits, tokens, quota, or allowance are exhausted.
A bare 429, overload, timeout, or transport error does not establish exhaustion. Context-window limits and a per-response token cap also do not qualify.
A short reset does not justify another attempt. The example `All claude accounts have exhausted their quota (reset after 5m)` qualifies only through a proven provider channel.

Record a new immutable trigger for new ideas. Preserve all saved trigger meanings and hashes. An owner-bound revision can upgrade an existing idea.
Run the exhaustion predicate before the five-second retry branch. On qualified evidence, record one terminal attempt and stop dispatch to that identity.
The durable step ledger must prevent a second attempt after restart, prompt changes, and driver re-entry. Valid completed output retains precedence.
Transient failures retain two attempts. Cancellation, unresolved writers, control-plane refusals, and corrupted evidence do not authorize a drop.

Keep the legacy 60-minute predicate for saved quota-only policy. It must not limit the new account-exhaustion class.
Separate exhaustion recognition from reset validation. An absent reset still permits an idea drop. A contradictory reset must never authorize a cache expiry.

### Provider scope

Start with the failed agent identity. An adapter name or model label does not identify an account or connection.
Do not exclude peers solely because they use the same provider brand. They may use different accounts.
I favor deferring a machine cache unless existing configuration supplies a stable, non-secret connection identity with proven scope.
This release can satisfy no reconnection within the current idea through immutable permanent-drop history.
If peers support a cache, require explicit connection keys, bounded records, reset-aware expiry, and no secrets. Unknown reset must not create a permanent machine outage.

### Role succession

Use the kickoff participant order. Search after the failed holder, then wrap once. Use roster order only for an explicitly eligible current member absent from that order.
Do not silently add roster members to the idea. Do not add a `succession:` field unless a concrete case requires it.

Resolve all lost roles as one durable transition after stopped-writer proof. Skip dropped, unusable, and role-ineligible members.
The transition must bind each old holder, new holder, role, artifact, and evidence record. Re-entry must produce the same successor without another launch.
A successor implementer continues from the branch and canonical documents. It records the handover in IMPLEMENTATION and the single batch notice.
The new implementer cannot review any code it wrote. Track previous implementers when reviewer eligibility needs that history.

A started drafter can leave a partial canonical file. Preserve that file's exact bytes under an immutable handover path before a successor draft.
Never overwrite its signoffs, findings, or dissent. A frozen FINAL stays frozen; later execution changes belong in IMPLEMENTATION.
Consensus recovery must preserve existing append-only signoffs and require signoff on changed decisions. A takeover does not imply acceptance.

The driver can replace its dispatched role holder. It cannot replace the external operating-system process that organizes a manual run.
For an external organizer, the protocol specifies the ordered successor and handover. State this enforcement boundary clearly.
A declared pure facilitator still cannot supply one of the two usable non-facilitator seats.

### Review gate and floor

A failed implementer can turn one former reviewer into its successor. The old single-reviewer predicate requires the same implementer before and after.
Extend the cause check for this specific recorded succession. Require at least two independent reviewer candidates before and one eligible reviewer after.
The after-set must match current membership. Every removal and role transfer needs typed exhaustion evidence and deterministic eligibility.
Require known different snapshot models for a sole reviewer. Retain current signers, independent goal checks, reservations, strict review, and current-tree evidence.

If no eligible successor exists, stop with one blocking owner note. Do the same below two usable non-facilitators.
Exhaustion must never erase a veto or unresolved finding. Retained dissent needs the existing owner ruling or a new idea.

### Scope and size

Use the existing reducer, step ledger, membership history, notice receipts, and driver role resolution.
Compact repeated protocol prose before new text. Keep all three normative copies equal and preserve their project header zones.
Record exact hunks for §§0, 4 Phase 0/4/5/6/7/8, 5, and 9.0. Include only sections whose current wording conflicts.
Measure the normal and meta-idea phase-1 packets. Keep the unchanged caps and reduce the current meta overage where possible.
Compact SKILL and update its existing references. Do not move a binding rule into an unreferenced document.

## Existing alternatives

- ALT-1: Extend `internal/quota.Evaluate` and immutable `History`; do not add another membership database. Constraint-forced by existing authority and recovery rules.
- ALT-2: Extend `internal/runner.RunParticipantStep`; do not add a retry service. It already binds attempts to idea, agent, and logical step.
- ALT-3: Reuse `internal/membership.SingleReviewerAfterDropout`; extend its cause proof for implementer succession. Preserve the existing failure and quota cases.
- ALT-4: Reuse `internal/protocol.ResolveImplementerChain` with a validated succession override. Do not rewrite owner designations as fabricated owner confirmations.
- ALT-5: Keep per-idea identity scope and defer a machine circuit breaker. Stable account identity is not established by model names.
- ALT-6: Reuse the Zcode stderr framing parser where its established evidence contract applies. Add each other adapter only with located native framing and adversarial fixtures.
- ALT-7: Preserve manual organizer handover as a protocol duty. Do not create a process supervisor merely to replace this external session.

These are design choices. Their acceptance depends on peer analysis and concrete implementation evidence.

## Concerns / open questions

The shipped recognizer supports only the bounded Zcode owner exception. Other adapters remain diagnostic-only in `internal/telemetry/QuotaSupport`.
Broad phrase matching over assistant output would violate the brief. Locate native provider channels before claiming adapter support.
Do not claim a complete native Zcode exhaustion capture from the existing partial excerpts or source-derived fixtures.
Assess the capture follow-up's minimum needed portion. Preserve its unfulfilled native-capture requirement if no ordinary failure supplies the missing evidence.

The startup command cannot create the required named idea. The measured adapter dispatch also ignored the roster's model values.
The task-local launch file restores the approved models. These gaps do not justify unrelated product work.
The standalone preflight reports a defect in another idea. That defect does not alter this idea's explicit participant list.

## Risks

False provider attribution could drop a healthy dissenter. Require typed origin, failed invocation facts, and no valid-output precedence conflict.
Concurrent writers could corrupt a takeover. Keep the existing lease and stopped-writer requirements.
A successor may inherit partial code with no valid review. Require fresh review of the complete resulting implementation.
A machine cache without account identity could disable unrelated agents. Prefer the smallest auditable scope until identity is established.

## Protocol context

context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none

I read the complete live authority and Appendix B. No verification verdict about my own proposal is claimed.
