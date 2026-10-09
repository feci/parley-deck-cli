---
agent: codex-1
idea: meta-protocol-change-participant-dropout-review-gate-timing
round: 1
date: 2026-10-09
---

## Summary

Extend the existing membership gates with one evidence-derived exception: when validated automatic exclusion history proves that losing reviewers caused the count to become exactly one, one known model-diverse non-implementer may review and perform a fresh-process goal check. Keep every substantive review/strict/criterion gate intact. Tighten participant-failure streaming watchdog bounds with the existing supervisor and its existing two-attempt ledger; never create a second retry or membership system.

## Proposed approach

R1-1: Add a small shared decision at the existing membership gate seam, used prospectively before committing a settled eligible batch and again from validated immutable history at review/open/close. Require the current implementer and exactly one eligible reviewer, distinct known configured model IDs, a previously proposed/current set containing at least two independent reviewers, and a typed failure/exhaustion removal that exactly accounts for the reduction. No waiver from a marker, prose, unrelated old dropout, role edit, arbitrary owner membership edit, count alone or corrupt/pending history. Include kickoff transitions and the legacy typed quota-exclusion trigger; saved trigger semantics stay unchanged. An explicitly owner-authorized attended exclusion without typed automatic history keeps the attended path, as this dogfood run does.

R1-2: Exemption changes only the numerical threshold. Require all current participants' ACCEPT signoffs for deliberation review consensus; retain excluded signers' BLOCKs, disputes and findings. Reservations still stop unattended close. A strict idea still needs a fresh full-scope zero-findings closing round and subsequent clean consensus; one eligible reviewer supplies that round. Current-tree independent AC evidence remains required. A fresh one-shot process of the remaining reviewer may goal-check, never the implementer or an additional quorum identity. The current code's goal check already chooses an independent reviewer/drafter; record process freshness and keep missing/failed/inconclusive verdicts fail-closed.

R1-3: Use supervisor facts, not elapsed-time inference. For streaming participant-failure steps use first-output default 120s and stall default 300s, bounded by the existing effective per-track/process ceiling (fast 300s, standard 900s, deliberation 1800s where applicable). Existing shorter explicit settings remain effective; explicit disabled/long overrides and buffered transports need an explicit documented decision: my preference is preserve the configured transport-disable semantics and report that only the hard ceiling applies there, rather than killing a healthy buffered agent. First-output/stall/timeout count as the existing first attempt; after five seconds the same ceiling applies to the sole retry. Cancellation/control-plane failure/integrity are excluded, and valid own dissent wins. No third launch after restart. FINAL must specify whether the new defaults are limited to participant-failure mode, and tests pin that boundary. My preference: scope them to the new trigger, leaving legacy/default supervision intact.

R1-4: Batch settlement stays atomic after all writers stop. Existing current/known membership and artifact validation should already permit survivors to complete the interrupted round/consensus/signoff batch; use tests to expose any stale expected-set copies and repair only those sites. Excluding a valid dissenting author must remain impossible, and a failure in a later step never disposes of an earlier filed BLOCK. Floor/protected-role failure yields the existing one decision, never retry-until-success.

R1-5: Compact normative repetition in §4 and §9.0 without removing conditions; annotate exact old/new substitutions for all three protocol copies and staged core reconstruction. Target meaningful headroom (several KB for the phase-1 packet, at least 1KB for core SKILL.md) below the existing 70,000/20,000 limits. Change no applicability-map scope or limits to hide growth. Keep the implementation bounded to gates/history qualification, caller threshold wiring, timing defaults, stale-batch fixes only if demonstrated, focused tests and matching guidance/version artifacts.

## Evidence and scope (owner propositions awaiting independent verdict)

Read at CLI HEAD 85c5a8f. These are my located implementation observations, not self-issued §15 verification verdicts; Zcode should independently test them.

- internal/membership/gates.go: CheckGates sets minimum=2 under auto, then errors with `review/LE-7/LE-11 gate`; internal/membership/membership.go: Settle calls it before CommitBatch. This is the shipped precommit stall, corroborating the prior released note.
- internal/driver/impl.go: under AutoImplement, ReviewerCount < MinReviewers escalates; reserved consensus also escalates. GoalCheck runs before close, followed by scope/evidence validation. Merely changing the membership gate would leave close blocked.
- internal/track/track.go: deliberation leaves MinReviewers unset (driver default 2), while standard's ordinary two-participant degradation already exists. Do not change the general track policy to one reviewer.
- internal/runner/supervision.go: first-output default=120,000ms; stall=1,800,000ms clamped to hardTimeout-1s; BuffersStdout disables both soft timers. Thus the reported seven-minute no-output hang cannot be explained as a claim that a 120s enabled supervisor ran; inspect actual probe/transport wiring and avoid claiming prior anecdotes prove measured product behavior.
- internal/runner/dropout.go: RunParticipantStep already binds retries by idea/agent/logical step and waits five seconds before attempt two; reuse it.
- internal/quota/{record,history,revision}.go: kickoff transitions, validated automatic batches and explicit owner revisions already differ structurally. Reason checks can be typed and history-backed without a new policy schema.

## Existing alternatives

- ALT-1, membership reduction: internal/membership.Settle plus internal/quota.History already ship atomic batch, durable history, receipts and floor/protection. Adopt; constraint-forced by compatibility and scope. Reject a new roster manager.
- ALT-2, single reviewer: internal/track.PolicyFor already degrades standard for two participants. Reject using that unconditionally for auto_implement/deliberation because it erases the recorded-cause requirement. Add an evidence-derived exception at existing gates.
- ALT-3, timing/retry: internal/runner/supervision.go and RunParticipantStep already ship activity tracking, watchdog classes, cleanup and durable two-attempt retry. Adopt. Reject a timer that infers dropout directly or a new retry service.
- ALT-4, independence: separate existing GoalCheck process plus current-tree completion evidence, not a second invented reviewer identity. Adopt fresh-process reuse; model diversity versus the implementer is the hard lower bound. The same reviewer can share conceptual blind spots, disclosed below.
- ALT-5, manual continuation: prior participant-dropout FINAL D7 and this controlling brief supply attended authority. Keep only as fallback for historical/manual evidence that cannot prove the automatic exception; do not ask for it again in this run.
- ALT-6, size: existing drift/packet-size/skill-size tests. Adopt semantic compaction and unchanged limits; reject larger limits or omitted normative sections.

## Concerns / open questions

Precisely distinguish a causal exclusion from an unrelated old transition; avoid admitting an unrelated owner removal after a legitimate dropout. Confirm whether default tightening should honor explicit disabling and whether first-output startup text should count as activity (current code says yes). Buffered transports cannot support a truthful no-output deadline earlier than their hard limit; document this rather than pretending heartbeat proves provider progress. The preflight command has no participant filter and hits an unrelated historical facilitator declaration; recorded measured launches are necessary for this named idea.

## Risks

One reviewer plus a fresh instance shares model/architecture priors; actual adversarial AC evidence, retained dissent and strict scope discipline matter more than agreement. Short stall defaults may kill legitimate deep reasoning on a streaming-but-quiet tool; honor disclosed override/buffer behavior. History corruption must fail closed. Too broad a causal predicate silently relaxes unrelated ideas, so negative counterexamples are primary acceptance tests. Full host tests should use local tmp storage: earlier shared-mount execution failures are not product proof. Windows remains explicitly experimental.

## User direction and this run

Read the full copied controlling brief. It appoints codex-1 organizer/implementer and zcode-1 reviewer, allows excluding still-broken Kimi, and pre-authorizes owner-attended close after no open CRITICAL/MAJOR, both ACCEPT blocks, fresh Zcode goal-check PASS and current-tree AC evidence. No open owner decision is asserted at this stage. We will use the brief's exact quotes in the final authority record and never claim the old driver's auto-close gate passed.

## Protocol attestation

context_mode: full; source_sha256: 091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf; packet_sha256: 091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf; fallback_reason: absent. Complete live protocol read. This file was authored before reading any other round-01 artifact.
