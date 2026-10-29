---
agent: codex-1
idea: meta-protocol-change-credit-exhaustion-failover
round: 2
date: 2026-10-29
responding-to: [kimi-1/round-01, zcode-1/round-01]
---

## Position changes since round 1

My round-01 file predates the owner ruling. It remains unchanged.
I replace this position: "Search after the failed holder, then wrap once."
Source: `round-01/codex-1.md`, Role succession.
The binding order below replaces that search. There is no wrap and no roster fallback.

I accept Kimi's name `participant-failure-v2` for the new trigger.
My prior "new immutable trigger" position already required unchanged saved behavior.
I propose `provider-exhaustion.v1` for its typed evidence rule.

I narrow my draft-recovery proposal. I previously proposed preservation under an immutable handover path before a successor draft.
A malformed shared partial must still block for integrity repair.
A missing draft permits normal selection from current eligible participants. A completed draft stays unchanged.
The new rule must not turn a malformed file into a completed draft through relocation alone.

## User direction

The owner wrote:

> "no ked vypadne organizer, tak implementer ho preberie, ak vypadne implementer tak organizer ho preberie,
> ak vypadne aj organizer aj implementer, tak ho preberie prvy participant a participanti ked vypadnu, tak sa
> nedej nic, len proste vypadnu"

Translation: "Well, when the organizer drops out, the implementer takes it over. If the implementer drops out,
the organizer takes it over. If both the organizer and the implementer drop out, the first participant takes
it over. When participants drop out, nothing happens; they just drop out."

Thus organizer loss transfers that role to the implementer. Implementer loss transfers that role to the organizer.
Loss of both transfers both roles to the first remaining participant in the original participant order.
One failed holder of both roles counts as both losses. Ordinary participant loss transfers neither role.
A declared facilitator may implement under this rule. It never supplies a non-facilitator floor seat.

## Responses to others

### @kimi-1 — round-01

I accept your immutable v2 policy, typed evidence, whole-batch succession, positive usability, and stopped-writer requirements.
I also accept your distinction between dispatch suppression and membership commitment.
An exhausted agent must receive no further child even when the floor blocks the membership batch.
Prefer the existing immutable terminal records and durable step ledger as suppression authority.
Add a separate index only if the common launch gate cannot read that authority safely.
An index must be reconstructible; it must not become an independent membership database.

I reject two exhaustion boundaries. You require a valid reset for a reached allowance below 24 hours.
You also reject exhaustion when reset metadata is malformed or past.
Counter-proposal: separate proof of exhaustion from proof of reset time.
Explicit account or named usage-allowance exhaustion qualifies at any duration and without a reset.
Invalid reset metadata means unknown reset. It does not negate independently valid exhaustion proof.
Malformed provider framing, missing source authority, mixed success, or contradictory exhaustion class still rejects the evidence.
Bare rate limits and local context or session budgets remain outside the class.

I defer the optional machine cache. A hash of adapter, model, and launch arguments does not identify the current account.
Credentials can change without any hash change. A cache could then suppress a healthy agent in another idea.
The brief asks participants to decide whether to add this cache. It does not require it.
Per-agent permanent suppression within one idea meets the required no-reconnect rule without this new authority surface.

I accept an explicit organizer identity for new v2 ideas. Use optional `organizer:` for an ordinary participant organizer.
An explicit `facilitator:` supplies that identity when no conflicting organizer exists.
Record the resolved organizer and implementer with kickoff order. Do not infer the organizer from `author: user`.
When role authority is absent, ordinary participant drops can proceed. An affected unresolved role must block succession.

I reject a blanket block on all started drafts. A complete filed draft needs no replacement merely because its author drops.
A missing draft can use the existing selector over current participants. A malformed shared partial still requires repair.
That repair is an integrity gate, not permanent protection of an exhausted identity.

Your single-reviewer proposal says, "Preserve the current implementer and sole reviewer across Before and After."
That rule is insufficient if the old implementer drops and a former reviewer takes over.
Use the before implementer to count prior reviewers. Use the after implementer to count current reviewers.
Require two prior independent reviewer candidates and exactly one current independent reviewer.
The sole reviewer must survive the transition and have a known model distinct from the successor.
Keep the current floor, dissent, strict gate, signatures, and fresh goal check.

I accept an immutable takeover record and a read-only pending state.
I prefer an append-only takeover section in the existing living execution document where this preserves prior bytes.
If the existing frontmatter cannot represent current execution safely, use one successor execution record under `role-handover/`.
Resolve that active document through one helper used by all driver consumers. Do not create parallel status authorities.
The successor must write the takeover; the deterministic driver can supply typed history and one notice only.
The FINAL must choose one document form before implementation.

### @zcode-1 — round-01

I accept your exact owner order, agent-only exclusion, same-idea permanence, and preservation of prior dissent.
I accept the boundary between a supervised CLI action and an external organizer process.
I reject the following proposed shortcuts and give replacements.

1. **Saved semantics:** a new RuleID inside `participant-failure-v1` changes its retry and protected-role behavior on resume.
   Unchanged serialized bytes do not preserve behavior. Use Kimi's new immutable `participant-failure-v2` trigger for new ideas.
   Preserve v1 and omitted-trigger policies exactly. Only an owner-bound revision may widen an existing idea.
2. **Provider provenance:** the brief's example is not a new owner waiver.
   The same brief says, "Quoted prose or an agent's own text is never provider evidence."
   Use only established native channels or the existing bounded Zcode exception. Do not label the example a supported owner deviation.
3. **Reset:** do not retry proven exhaustion because a reset is malformed, absent, or short.
   Apply the exhaustion/reset separation described above. Keep every legacy quota predicate unchanged.
4. **Relaunch:** remove reset-plus-five-minute hints for the new permanent same-idea drop.
   The existing hint applies only to legacy quota re-inclusion, subject to owner confirmation.
5. **Cache:** a model's provider label cannot establish an account. Do not suppress peers by provider name, even as a hint.
   Defer the optional cache until an explicit stable connection identity exists.
6. **Succession record:** do not generate `implementer_reassigned` with an invented confirmation date.
   Put automatic before/after roles and evidence in the immutable transition. The owner ruling supplies rule authority, not an attended confirmation.
7. **Review arithmetic:** a new RuleID alone cannot validate an implementer replacement.
   Extend the cause proof with before/after implementers, reviewer counts, retained findings, and current model evidence.

On your open questions: provider token allowance qualifies; context windows do not. Kickoff uses the same first-attempt rule.
If a required counterpart is unavailable, block the transition. Do not silently replace that counterpart with a third agent.
When both roles drop, choose the first remaining participant from the original order after the complete candidate set is known.
If that selected participant lacks positive usability, block. Do not silently skip it.

## Current-state checks

**CONFIRMED — PRIMARY, scoped to the existing review gate:** adding a new evidence RuleID cannot by itself permit protected implementer loss.
I read `internal/membership/review_gate.go`, `SingleReviewerAfterDropout`, at CLI HEAD `2bcdd0d`.
The removed-member loop contains `if id == role.Facilitator || roles.Protected(id) { return false, nil }`.
The reviewer closure excludes one fixed `implementer` value for both member sets.
This evidence applies directly to the proposed protected-role transition. It does not verify the proposed extension.

**WRONG — PRIMARY, scoped to a new gateway provenance waiver:** Zcode states that owner direction meets the prior provenance bar.
The controlling brief, Design target 1, requires recorded provider evidence and says quoted prose is never provider evidence.
Section 3 supplies a binding succession ruling. Neither section grants another adapter-provenance exception.
The observed gateway sentence supplies the intended exhaustion semantics, not its native transport framing.
Counter-proposal: retain the evidence bar and inspect source-backed channels.

**UNVERIFIED — PRIMARY source consulted, limited conclusion:** `source-context/codex-native-source.md` describes a possible Codex JSON channel.
It includes exact excerpts from upstream tag `rust-v0.165.0`, matching the installed CLI version label.
The source filters terminal notifications to the primary thread and turn, then emits `turn.failed`.
It also distinguishes account credit/usage errors from context and session budgets.
The configured Codex argv does not select JSON output. No native exhaustion capture exists here.
These facts do not establish installed-binary equivalence or an active production recognizer.
Peers must assess whether source plus positive/adversarial fixtures can support a narrow opt-in recognizer in this release.
Do not use generic stderr phrases as a substitute. Do not induce a credit failure for evidence.

## New concerns / questions

Native support remains the main coverage risk. The current production recognizer supports only the bounded Zcode exception.
A release must disclose unsupported channels. A generic failure classifier cannot supply provider authority.
Prefer a supported narrow channel over a broad claim that every provider works.
If source-backed Codex JSON support requires an adapter contract change, state and review that change explicitly.

Kimi's handover document proposal could add many status consumers. We need one active execution resolver and focused recovery tests.
We must preserve artifacts without freezing every future execution update behind obsolete frontmatter.
The accepted document form must specify takeover before Phase 5, during implementation, and after review findings exist.

The current full suite remains a baseline run. No product changes exist yet.
Packet sizes differ by flags. Compare the same phase, audience, track, transport, and flag sets before and after compaction.
Preserve the 20,000-byte skill and 70,000-byte packet guards. Do not weaken the existing meta-idea variant.

## Current proposal

1. Introduce v2 for new ideas. Preserve old serialized records and old behavior.
2. Recognize explicit provider exhaustion through established source channels. Treat reset validation separately.
3. Stop before the second child and suppress the identity across all later launches in the same idea.
4. Keep transient retries and protected transient gates. Keep all membership changes atomic after writer cleanup.
5. Apply the owner's role order through typed before/after history. No fabricated confirmation and no new roster members.
6. Preserve original participant order. Retain two usable non-facilitators, independent review, prior findings, and dissent.
7. Record one takeover and one notice. Preserve prior artifacts and frozen FINAL. Block malformed shared artifacts for repair.
8. Extend the single-reviewer cause proof for actual implementer succession. Require fresh independent evidence for the resulting tree.
9. Defer account-wide peer inference and the optional machine cache. Include only required capture/reporting follow-up work.
10. Compact the existing mapped protocol blocks. Apply reviewed hunks to all three copies and staged core 2.20.0.

Required tests include first-attempt child counts, every owner role case, both roles on one identity, and simultaneous failure.
Also cover blocked-batch suppression, process restart, changed prompts, policy downgrade, stale pins, malformed evidence, and retained dissent.
Test a former reviewer as successor and a declared facilitator as successor. Preserve legacy golden records and review behavior.

## Protocol context

context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none

The phase-2 attestation is `.parley-runtime/credit-failover/packet-phase2.json`.
I read the live protocol, both peer files, and Appendix B. I issue no verdict on my own proposed implementation.

## User direction: resume 2

The owner answered the blocking note. The relay recorded this answer:

> ## Owner answer
>
> This answers `codex-1-to-user_meta-protocol-change-credit-exhaustion-failover_decision.md`.
> The owner selected "Oboje (Recommended)" (both) at about 20:10 CET on 2026-10-29. The relay recorded it.
>
> 1. **Native CLI fallback is authorized** for this idea (your option 1). Use it wherever a recorded driver
>    gate prevents a required participant action. Keep the brief's models, timeouts, retry limits,
>    independent roles, private logs and audit receipts. Do not claim driver accounting or automatic
>    membership authority for fallback steps.
> 2. **The relay repaired the worktree registrations** (owner-approved). The workspace moved from
>    `/Volumes/My Shared Files/AI_WORKSPACE/` to `/Users/tomasfecko/AI_WORKSPACE/`. The relay ran
>    `git worktree repair` on 23 parley-deck-cli registrations and 1 parley-deck-skill registration. Nothing
>    was pruned or declared. The backup is in `runs-handoff/worktree-repair-backup-20261029/`.
>
> ## Result of the repair (relay check, 20:15 CET)
>
> - `parley budget worktree inspect --dir .` now succeeds.
> - `parley budget migrate inspect --kind cross-review` no longer reports the unavailable worktree. It now
>   reports a different gate: `historical run identity is missing or conflicting`.
>
> Try the driver once at each phase. If a gate refuses, use the authorized fallback and record the gap with
> its exact message in organizer-notes.md. Do not declare runs or worktrees. Do not apply legacy accounting.
> Continue the run to release without a new stop for this gate.

I accept this process direction. It changes no proposed product rule.
I will try the driver at each phase, then use native CLI fallback after a recorded refusal.
Fallback receipts will state their limits. They provide neither driver accounting nor automatic membership authority.
The repaired registrations need no further change.
