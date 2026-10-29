---
agent: kimi-1
idea: meta-protocol-change-credit-exhaustion-failover
round: 1
date: 2026-10-29
---
## Summary

I propose a new `participant-failure-v2` policy for new ideas. It stops recognized provider exhaustion after the first attempt and applies the owner's exact role succession. Existing policies keep their current retry, reset, return, and reviewer behavior. Peer exclusion should remain agent-scoped until evidence can prove a shared account.

## Proposed approach

### Current-state evidence

- **CONFIRMED — PRIMARY:** the generic participant path permits two attempts before it emits dropout evidence. `internal/runner/dropout.go:200` says `for len(attempts) < 2 {`. Lines 261–267 construct evidence only after that loop.
- **CONFIRMED — PRIMARY:** protected candidates block instead of receiving a successor. `internal/quota/quota.go:217–220` returns `protected role %s; re-designate, record implementer_waived, or set implementer: none with owner confirmation`.
- **CONFIRMED — PRIMARY:** the legacy recognizer rejects a stated reset below one hour. `internal/telemetry/quota.go:124–129` tests `reset.Sub(in.ObservedAt) < quota.MinimumReset` before assigning `quota.named-reset-ge-60m.v1`.
- **CONFIRMED — PRIMARY:** native authority remains narrow. `internal/telemetry/quota.go:20–23` marks codex, kimi, and claude `diagnostic-only`; only zcode has the bounded owner deviation.
- **CONFIRMED — PRIMARY:** current role data lacks an organizer field. `internal/quota/quota.go:120–126` contains only `Facilitator`, `Designee`, `PinnedImplementer`, and `StartedDrafters`.
- **CONFIRMED — PRIMARY:** size headroom is small. `wc -c skills/parley-deck/SKILL.md` returned `19482`. `go test ./internal/app -run 'TestLiveDeckFacilitator(PacketNamedSetsAndGuardrail|AcrossPhases)$' -count=1 -v` returned `facilitator body=68874 B` and `guardrail=70000 B`.

These verdicts describe HEAD `f3b8a4fc0d453c0b78011924e5267a38a369bf59`. They do not verify my proposal.

### 1. Version the policy and the exhaustion class

Use `participant-failure-v2` for ideas created after this release. Keep `participant-failure-v1` and omitted-trigger quota policies byte-compatible and behavior-compatible. An owner-bound revision may select v2 explicitly. Resume and binary upgrade never widen a saved policy.

Add one typed `provider-exhaustion.v1` evidence class. Build it only from established native provider channels. Never derive it from the generic regex classifier in `internal/runner/failclass.go`, assistant text, tool output, or an `is_error` bit.

Apply these semantics:

- Explicit account credit, token, or quota exhaustion qualifies with or without a reset.
- A reached named usage allowance qualifies with any valid future reset, including a five-minute reset.
- Without a reset, a named allowance must cover at least 24 hours. Retain this conservative boundary for short throttles.
- A valid artifact or an already in-flight later success wins. Do not launch another attempt after recognized exhaustion.
- Bare 429, bare 5xx, overload, authentication errors, generic quota wording, timers without allowance semantics, and watchdog outcomes do not qualify.
- Contradictory, past, malformed, or truncated evidence does not qualify. Every decisive provider record must agree.

The 60-minute rule remains unchanged for legacy quota policy and its return hints. It does not delay the new same-idea stop. A short reset affects only the optional machine cache duration.

Split native capture from policy evaluation in `internal/telemetry/quota.go` and `quota_zcode.go`. Preserve the existing legacy result separately from the new v2 result. This prevents the new reset rule from changing legacy evidence.

Include the follow-up's private capture facility for failed zcode streams. Require mode `0600`, bounded bytes, secret scrubbing, and visible truncation markers. Keep raw streams out of canonical artifacts. Do not relax zcode grammar without a genuine native capture that meets the existing follow-up precondition. Do not manufacture that capture through a probe.

### 2. Stop the first attempt and suppress reconnection

Persist validated v2 evidence in an immutable idea-local `provider-exhaustion/<invocation-id>.json` record before returning from the failed step. This record suppresses dispatch; it does not change membership by itself. A write failure blocks further dispatch for repair because retry suppression is not durable.

Check this evidence in `RunParticipantStep` immediately after terminal persistence and before the retry delay. Return a distinct exhaustion outcome without a second child. Add the same check to the legacy exec, ACP, readiness, and signoff retry seams. Protected roles keep their existing transient-failure behavior.

Check the suppression ledger at the common launch admission boundary. A later run, changed prompt, or new run ID must not reconnect that agent in this idea. A committed exclusion or a blocked batch must not remove this suppression.

For a transient failure, retain the shipped original-plus-one behavior and five-second delay. For an unsupported adapter, retain that path and the existing diagnostic-only limitation.

CLI enforcement covers invocations under its supervisor. An external organizer must supply equivalent native evidence through its supervisor. A note, timeout, or organizer prose cannot authorize automatic organizer failover.

### 3. Keep peer scope narrow and add an optional machine cache

Exclude only the agent named by the native evidence in this release. Do not infer a shared account from adapter, model, gateway, or pool wording. The current roster snapshot has no account binding, and the zcode exception accepts subagent-source ambiguity.

Do not probe peers to test that inference. Each peer requires its own qualifying evidence. Account-wide exclusion needs a later design with a native account identifier and an exact non-secret roster binding.

Add a small machine cache only for a known future reset. Key it by agent ID and a hash of the effective adapter, model, and launch arguments. Suggested location: `~/.parley/provider-exhaustion/v1/<agent-fingerprint>.json`.

Store only schema version, agent, fingerprint, class, observation time, reset time, source idea, source invocation, and evidence digest. Store no endpoint, credential, account secret, raw log, or launch argument text. Revalidate the linked immutable evidence before each cache use.

Before kickoff probes, skip only an exact fingerprint match with an unexpired record. Apply the normal floor, role, and succession gates to the resulting candidate. At or after reset, probe afresh. Never rejoin the same idea automatically.

Do not create a cache record for an unknown reset. A changed fingerprint, corrupt record, or missing source evidence disables the cache and emits a diagnostic. Normal readiness then applies. A later valid probe may expire the cache without changing historical records.

### 4. Implement the owner's succession atomically

Record an initial role snapshot for v2 ideas. Add optional `organizer:` frontmatter for an undeclared organizer. `facilitator:` supplies the organizer when declared. Automatic role failover requires an unambiguous recorded organizer; do not infer one from authorship or participant order.

Extend the existing immutable batch with optional role before/after state. Allow a role-only transition when the organizer is outside `participants:`. Keep membership and role changes in one checked-durable commit. Preserve old records under the existing strict decoder.

Evaluate the whole batch before selecting a successor. Use kickoff `participants:` order, never the sorted helper used for set comparison. Apply the binding ruling exactly:

1. Organizer loss: the implementer takes the organizer role.
2. Implementer loss: the organizer takes the implementer role.
3. Loss of both roles: the first remaining eligible participant takes both roles.
4. Ordinary participant loss: no succession occurs.

One agent with both roles counts as both losses. A declared facilitator may implement under rule 2 despite the usual pure-organizer restriction. It still supplies no non-facilitator floor seat.

Eligibility requires current membership where applicable, positive usability evidence, and no matching exhaustion suppression. Evaluate every simultaneous drop before choosing the first remaining participant. If no eligible successor exists, apply no membership or role transition and file one blocking escalation.

Keep an unstarted drafter selection over current roles. Keep a started consensus, FINAL, or review-consensus draft attributed to its author. Do not complete or rewrite it automatically. Its incomplete state blocks the affected phase and receives one owner escalation.

### 5. Preserve independence, dissent, and the floor

Require two distinct usable non-facilitators after the complete batch. Require positive evidence for each successor; a role name alone supplies no seat. Preserve the all-or-none membership and role decision. First-attempt suppression remains even when that decision blocks.

Resolve the current implementer from validated role history before the old file pin. A valid succession must not trigger the existing pin-versus-designation conflict. An old `IMPLEMENTATION.md` pin must never resurrect a dropped implementer.

Recompute reviewers after every succession. The new implementer never reviews or goal-checks its own code. Its earlier findings remain binding, but they are not independent review of its new implementation. A review-consensus drafter must remain a non-implementer.

Extend `SingleReviewerAfterDropout` to accept `provider-exhaustion.v1` only with valid role before/after state and exact typed evidence. Preserve the current implementer and sole reviewer across Before and After. Require known distinct snapshot models. A manual membership edit still invalidates the cause.

Retain every filed BLOCK, DISPUTED claim, finding, and signoff. A dropped author cannot append or withdraw them. Resolve retained dissent only through the existing owner-ruling or v2-idea paths. Exclusion never supplies a disposition.

Require fresh current-tree acceptance evidence and a fresh non-implementer goal check after handover. Keep strict-gate, reservation, signer, diversity, and close duties unchanged.

### 6. Make handover and recovery explicit

Use a new canonical `role-handover/<transition-id>.md` as the active execution record after succession. The successor writes it as its first action. It names dropped roles, new roles, evidence, branch, HEAD, phase, and canonical resume paths.

Keep every prior artifact byte-unchanged. The handover carries later plan, progress, validation, and completion state instead of rewriting the dropped implementer's `IMPLEMENTATION.md`. The latest validated handover supersedes the old execution record for driving and close.

Commit membership and role history before the handover launch. Mark reconciliation pending until that artifact validates. Permit only the successor's handover action during this state. Read-only status, wait, and organizer brief must show the pending state without repair.

A crash before, during, or after handover must replay without a duplicate transition or a second ordinary dispatch. A new exhaustion during handover forms a new settled batch. Missing, contradictory, or truncated authority blocks recovery.

### 7. Exact change inventory

Apply the normative hunks identically to:

- `parley-deck/COOPERATION.md`;
- `internal/protocol/defaults/COOPERATION.md`;
- `skills/parley-deck/references/COOPERATION.md` in the skill worktree;
- staged core 2.20.0, derived only from staged 2.19.0 plus the reviewed hunks.

Protocol locations:

- §0 and Phase 0: record v2 and optional `organizer:` without a new boolean.
- §3: add the conditional handover artifact.
- Phases 3–4: define unstarted drafter fallback and the started-draft block.
- Phase 5: define role resolution, facilitator succession, and the active handover record.
- Phases 6–8: define reviewer independence, current signers, goal checks, and close evidence after succession.
- §5: define permanent provider suppression, role history, and retained dissent.
- §9.0: define evidence, reset, cache, first-attempt stop, succession, floor, legacy policy, and recovery.
- §9 checklist, §10, and §11: add only required notice, role-chain, and publication cross-references.
- `parley-deck/meta/protocol-changelog.md`: add this idea's ratified entry.

Keep `parley-deck/meta/packet-applicability.yaml` unchanged. Place all binding rules inside existing mapped sections.

CLI locations:

- Evidence and telemetry: `internal/telemetry/{quota,quota_zcode,quota_framing,record,dropout}.go` and a focused `exhaustion.go`.
- Policy and immutable records: `internal/quota/{quota,history,record,dropout,revision,notice}.go` plus focused exhaustion, succession, and machine-cache files.
- Retry and admission: `internal/runner/{dropout,runner,telemetry,launch,acp,consult,quota}.go`.
- Membership and gates: `internal/membership/{membership,gates,review_gate,revision,notice}.go` plus `succession.go`.
- Role consumers: `internal/protocol/{implementer,facilitator,quota,quota_manual,workspace}.go` and a focused organizer-role file.
- Driver and application: `internal/driver/{driver,impl,consensus,phasedigest}.go`; `internal/app/{app,preflight,preflight_liveness,driver_impl,driver_consensus,consensus_request_signoffs,quota_signoff,quota,organizer,wait}.go`.
- Creation and snapshots: `internal/runcontrol/runcontrol.go` and `internal/runmanifest/manifest.go`.
- Guidance and release metadata: `docs/{quota-membership,agent-cli-mechanics,agent-runtime-configuration,cli-reference}.md`, `README.md`, `CHANGELOG.md`, `VERSION`, and `internal/app/version.go`.

Skill locations:

- `skills/parley-deck/SKILL.md`: a compact rule summary and organizer handover duties.
- `references/ROSTER_AND_PROTOCOL.md`: policy, cache, succession, dissent, and recovery detail.
- `references/HEADLESS_LAUNCH.md`: first-attempt stop and native-evidence limits.
- `references/ARTIFACT_TEMPLATES.md`: the handover record shape.
- `references/COOPERATION.md`: the identical normative hunks.
- Version and compatibility metadata required by the existing release tests.

### 8. Regression and size contract

Add or extend these tests:

- `internal/telemetry/{exhaustion,quota,quota_zcode,quota_framing}_test.go`: short resets, no-reset classes, 24-hour boundaries, malformed evidence, spoofing, truncation, success precedence, and legacy equivalence.
- `internal/runner/{dropout,exhaustion,acp,review_gate_timing}_test.go`: no second child, transient two-attempt parity, protected-role preservation, restart suppression, and handover-only recovery.
- `internal/quota/{exhaustion,succession,machine,history}_test.go`: every owner role case, simultaneous batches, ordered succession, role-only records, cache expiry, and strict replay.
- `internal/membership/{succession,review_gate,gates,dropout}_test.go`: floor arithmetic, positive usability, no self-review, retained dissent, model diversity, and the causal one-reviewer exception.
- `internal/protocol/{implementer,facilitator,quota}_test.go`: no stale-pin return, facilitator rule 2, projection integrity, and legacy prompt compatibility.
- `internal/app/{dropout,preflight,preflight_liveness,driver_impl,driver_consensus,quota_surfaces}_test.go`: real dispatch suppression, cache-backed kickoff, handover recovery, blocking notices, and consistent read surfaces.
- `internal/app/facilitator_packet_live_test.go`: retain the unflagged guard and add a hard `protocol_change` variant below 70,000 bytes.
- Skill `test/lean-organizer.test.js`: retain the 20,000-byte guard and assert the new guidance and protocol parity.

Cover a blocked batch with no eligible successor. Prove that suppression survives it. Cover both roles on one agent, an external organizer, a declared facilitator, and a successor who was previously a reviewer. Inject a fault at every history and handover write.

Run `go test ./... -count=1 -timeout 45m`, build, vet, packet checks, and the full skill suite on the reviewed tree. Independent reviewers must supply current-tree evidence, not producer summaries.

Compact §9.0 by removing repeated trigger, floor, and recovery wording without dropping a rule. Compact the skill's automatic-exclusion section and move operational detail to its references. Do not raise either cap or remove required packet content.

The current unflagged packet leaves 1,126 bytes. The brief reports the meta-idea variant at 71,934 bytes, so compaction must first recover at least 1,934 bytes. Set a reviewed target below 69,000 bytes for both phase-1 variants to retain margin.

### 9. Deferred parts of the open follow-ups

Keep native zcode grammar expansion and native-positive acceptance in `quota-zcode-native-exhaustion-capture` until its evidence precondition occurs. Include only the private capture facility and existing strict forms here.

Keep alias-path guidance in `quota-kickoff-reporting-and-alias-guidance`. Include only the required kickoff escalation details and notice recovery for the new role transition. Do not reopen unrelated Windows, accounting, or roster work.

Defer account-wide peer exclusion and native recognizers for diagnostic-only adapters. Each needs separate evidence and an adapter-specific design.

## Existing alternatives

| ID | Hand-built component | Closest existing component and locator | Disposition |
| --- | --- | --- | --- |
| ALT-1 | Typed exhaustion and reset predicate | `internal/telemetry/quota.go:95` `refineNativeQuota`; `quota_zcode.go:34` `classifyZcodeQuota` | **Inherited.** Split native observation from the legacy predicate. |
| ALT-2 | First-attempt stop and launch suppression | `internal/runner/dropout.go:99` `RunParticipantStep`; `runner/telemetry.go:315` terminal quota evidence | **Inherited plus constraint-forced gate.** The gate must precede every retry path. |
| ALT-3 | Durable suppression evidence | Immutable invocation records in `internal/telemetry/record.go:93`; durable writes in `internal/quota/history.go:196` | **Inherited persistence; constraint-forced idea-local record.** Batch history alone cannot cover a blocked transition. |
| ALT-4 | Atomic role and membership transition | `internal/membership/membership.go:282` `Settle`; `internal/quota/history.go:31` `Batch` | **Inherited reducer; constraint-forced role extension.** Current batches admit only current-member candidates. |
| ALT-5 | Organizer and implementer role resolution | `internal/protocol/implementer.go:147` `ResolveImplementerChain`; `internal/app/driver_impl.go:57` | **Inherited chain; constraint-forced role snapshot.** Current roles lack an organizer identity. |
| ALT-6 | Handover artifact and resume pointer | `internal/driver/phasehandoff.go:247`; `internal/runner/handoff.go` | **Constraint-forced.** Those records are advisory or invocation-scoped, not canonical role authority. |
| ALT-7 | Review floor and single-reviewer derivation | `internal/membership/gates.go:22`; `review_gate.go:26` | **Inherited.** Extend typed evidence and role history; add no second gate. |
| ALT-8 | Optional machine cache | `internal/runmanifest/manifest.go:69` launch identity; `internal/quota/history.go:196` durable write | **Constraint-forced small cache.** Do not misuse roster files as an exhaustion store. |
| ALT-9 | Notice and projection recovery | `internal/membership/membership.go:25` `Before`; `membership/notice.go:56` kickoff receipts | **Inherited.** Add role pending state and handover validation only. |
| ALT-10 | Native zcode capture | `internal/runner/dropout.go:443` private file preservation; `internal/telemetry/quota_framing.go` | **Inherited capture seams.** New grammar remains evidence-gated. |
| ALT-11 | Packet and skill compaction | `internal/app/facilitator_packet_live_test.go:60`; `test/lean-organizer.test.js:81` | **Inherited guards.** Add the meta-idea case; never change the limits. |

## Concerns / open questions

- I recommend agent-only peer scope for this release. Cross-review should test whether the owner's wording requires native account-wide scope now.
- I recommend blocking incomplete started drafts rather than transferring them. Cross-review should test whether this creates unacceptable failover stalls.
- Automatic organizer failover needs a recorded organizer and supervised native evidence. A fully external organizer cannot receive the same CLI guarantee.
- The machine cache is optional and reset-bound. Its source-evidence revalidation must remain mandatory to avoid a persistent unverified lockout.

## Risks

- Provider wording can change or describe a short throttle. Native channel gates, conservative no-reset rules, and adversarial fixtures limit false positives.
- Role history adds recovery states beyond the current reducer. Atomic records, ordered role fixtures, and fault injection are necessary before release.
- Diagnostic-only adapters cannot meet the owner's immediate-stop goal with certainty. FINAL must disclose this limit and must not substitute generic failure text.
- A successor can destroy review independence when it was previously a reviewer. Current-role recomputation and retained-finding tests must block that path.
- A stale machine record can suppress a healthy agent or leak connection data. Exact fingerprints, reset expiry, source revalidation, and minimal fields reduce that risk.
- Protocol additions can exceed both byte caps. The FINAL needs measured compaction and a meta-idea guard before implementation review.
- A handover schema can accidentally resurrect old pins or rewrite history. Role-history precedence and byte-preservation tests must cover every consumer.
