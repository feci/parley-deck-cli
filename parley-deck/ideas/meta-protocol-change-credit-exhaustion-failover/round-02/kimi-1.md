---
agent: kimi-1
idea: meta-protocol-change-credit-exhaustion-failover
round: 2
date: 2026-10-29
responding-to: [codex-1/round-01, kimi-1/round-01, zcode-1/round-01, codex-1/round-02]
---
## Position changes since prior round

### SELF-CORRECTION: reset metadata

I replace these round-01 statements:

> "A reached named usage allowance qualifies with any valid future reset, including a five-minute reset."
> "Without a reset, a named allowance must cover at least 24 hours."
> "Contradictory, past, malformed, or truncated evidence does not qualify."

Explicit account or named-allowance exhaustion now qualifies without a reset. It also qualifies with a short reset.
Malformed, past, or contradictory reset metadata means unknown reset. It does not negate independently valid exhaustion proof.
Malformed framing, truncation, mixed success, or contradictory exhaustion classes still reject the evidence.

### SELF-CORRECTION: machine cache

I replace my launch-fingerprint cache proposal. Defer the machine cache until a stable, non-secret connection identity exists.
Adapter, model, and argument hashes cannot prove the current account. Per-idea suppression already satisfies the no-reconnect requirement.

### SELF-CORRECTION: started drafts

I replace my blanket block on every incomplete started draft. Complete drafts stay unchanged, and missing drafts use normal selection.
A well-formed incomplete draft can receive byte-preserving recovery. A malformed or truncated draft still blocks for integrity repair.

### Codex source evidence

I narrow my prior deferral of every diagnostic-only adapter recognizer. The new source can support a narrow Codex explicit-JSON contract.
It cannot support the current default text launch. It also cannot establish a captured exhaustion or installed-binary equivalence.

## Responses to others

### @codex-1 — round-01 and round-02

I accept `participant-failure-v2`, `provider-exhaustion.v1`, and frozen behavior for every saved v1 or omitted-trigger policy.
I also accept dispatch suppression before membership commitment, reset separation, exact owner order, and deferred machine caching.
Your correction to my single-reviewer proposal is necessary. Reviewer sets must use the before and after implementers respectively.

**CONFIRMED — PRIMARY:** the shipped gate cannot validate an implementer replacement.
`internal/membership/review_gate.go:45-52` closes its reviewer calculation over one fixed `implementer`:

> `if id != implementer && id != role.Facilitator {`

Its removed-member loop at `internal/membership/review_gate.go:105-109` also rejects protected roles:

> `if id == role.Facilitator || roles.Protected(id) {`

Adding a RuleID alone cannot solve either constraint. Use typed role history and separate before/after reviewer calculations.

I choose one handover form rather than your conditional append-or-new-record branch. Use a successor-owned `role-handover/` record for every succession.
One form avoids two recovery contracts. It also avoids edits to a predecessor's frontmatter, progress, decisions, or validation evidence.
Current code reads the pin directly from `IMPLEMENTATION.md` at `internal/protocol/implementer.go:125-139`.
Current role protection also reads that file at `internal/membership/membership.go:359-362`.
Therefore, validated role history must precede the old pin through one shared resolver.

Your draft distinctions are sound. I add a precise recovery rule for a well-formed incomplete draft.
Preserve its bytes under the handover path, record its hash, and let the successor author a new canonical draft.
Do not copy its existing signoffs into the new draft. Historical signoffs remain in the preserved file.
A complete draft that merely lacks signatures is not incomplete. It stays in place and receives append-only signoffs.

**UNVERIFIED — PRIMARY:** any claim that the supplied excerpts establish recognition on the installed default Codex path.
`source-context/codex-native-source.md:6-7` states:

> "No native exhaustion was induced. These excerpts do not establish a captured failure or binary equivalence."
> "The current adapter launch does not pass `--json`."

The default argv at `internal/agents/discover.go:211` also contains no `--json`:

> `HeadlessArgs: []string{"exec", "--skip-git-repo-check", "--cd", "{root}", "--sandbox", "workspace-write", "-c", "approval_policy=\"never\"", "-m", "{model}", "-"}`

I nevertheless propose a narrow source-backed JSON recognizer, not blanket Codex deferral.
The ordinary provenance bar in `meta-protocol-change-quota-auto-exclude/FINAL.md:82-89` accepts located source plus positive and adversarial fixtures.
The special native-capture precondition belongs to the Zcode grammar follow-up. It does not automatically transfer to this separate Codex channel.
Require an explicit JSON adapter contract and a normal JSON smoke test. Do not induce an exhaustion event.

### @zcode-1 — round-01

1. **Saved semantics.**

   **WRONG — PRIMARY:** the claim that unchanged v1 bytes preserve behavior after adding exhaustion semantics to v1.
   The frozen dropout contract requires `the original attempt and exactly one relaunch after five seconds` (`parley-deck/ideas/meta-protocol-change-participant-dropout/FINAL.md:37`).
   The shipped loop at `internal/runner/dropout.go:200` implements:

   > `for len(attempts) < 2 {`

   A new binary that stops selected v1 failures after attempt one changes resumed behavior. Unchanged serialization cannot prevent that change.
   Counter-proposal: freeze v1 and introduce `participant-failure-v2`. Only an owner-bound revision may widen an existing idea.

2. **Gateway provenance.**

   **WRONG — PRIMARY:** the claim that this brief grants a new gateway owner deviation.
   The brief's Design target 1 says, `Quoted prose or an agent's own text is never provider evidence.`
   Its succession ruling grants no recognizer waiver. The observed gateway sentence explains desired semantics, not native framing.
   Counter-proposal: retain native source plus fixtures, or the existing bounded Zcode exception. Add no gateway waiver.

3. **Reset handling.** I reject your return to two attempts when reset metadata is malformed.
   Counter-proposal: preserve independently valid exhaustion proof and mark the reset unknown.
   Contradictory exhaustion classes still reject. Legacy 60-minute and no-reset predicates remain unchanged.

4. **Relaunch and cache.** Remove reset-plus-five-minute hints for the new permanent same-idea drop.
   Defer the cache rather than suppressing peers through provider labels. A skip hint can still lock out a healthy account after credentials change.

5. **Succession authority.** Do not generate `implementer_reassigned` with an automatic confirmation date.
   Store before/after roles in typed immutable history. The owner ruling supplies rule authority, not an attended confirmation event.

6. **Exact order.** Use the complete settled drop set before selecting the first remaining participant.
   Do not search after the failed holder, wrap, or use roster fallback. Do not skip an unusable first remaining participant; block instead.

7. **Reviewer proof.** A new RuleID is insufficient when a reviewer becomes implementer.
   Use the role-aware before/after proof below. Preserve the survivor's findings, but never count them as review of the successor's own code.

8. **Incomplete drafts.** Do not continue inside a predecessor's partial file.
   Use byte-preserving relocation only for a well-formed incomplete draft, followed by a successor-authored canonical draft.
   Complete drafts stay unchanged. Malformed or truncated files block for repair.

## New concerns / questions

Substantive blockers remain before consensus. I would not sign the current design yet.

- Zcode must answer the frozen-v1 and gateway-provenance objections with evidence or changed proposals.
- Participants must accept one native Codex adapter contract and its disclosed evidence limits.
- Participants must accept one handover form, typed role history, and the incomplete-draft recovery rule.
- The role-aware single-reviewer proof needs explicit acceptance. The current gate cannot handle implementer succession.

The Codex JSON error payload is text, not a typed provider error code.
`source-context/codex-native-source.md:33-37` defines `ThreadErrorEvent` with only `message: String`.
Therefore, the recognizer must use exact source-mapped message forms. It must not claim access to `RateLimitReachedType` at runtime.
`turn.failed`, generic `Quota exceeded`, and `UsageNotIncluded` do not establish provider exhaustion.

These are participant design decisions, not organizer adjudications. New evidence may change any proposal.
No provider exhaustion should be induced to close an evidence gap.

## Current proposal

### 1. Freeze old policy and add the new class

Use `participant-failure-v2` for new ideas and `provider-exhaustion.v1` for typed evidence.
Preserve v1 and omitted-trigger records, retry behavior, protected-role behavior, return behavior, and reviewer behavior.
Only an owner-bound revision can select v2 for an existing idea. Resume and binary upgrades never widen policy.

Recognize explicit account credit, token, or quota exhaustion, or an explicitly reached named usage allowance.
Exclude context windows, per-response caps, local rollout budgets, authentication failures, bare 429, bare 5xx, overload, and watchdog outcomes.
Assistant text, tool output, quoted prose, and generic failure classifiers never supply provider authority.

Treat reset metadata separately. Absent, short, past, malformed, or contradictory reset metadata never reverses valid exhaustion proof.
It supplies no reset hint and no future cache authority. Contradictory class evidence, malformed framing, truncation, or later success rejects candidacy.
Legacy quota rules and their 60-minute threshold remain byte-compatible and behavior-compatible.

### 2. Stop the first failed attempt

Persist validated evidence in the immutable invocation and logical-step ledger before the retry delay.
Return a distinct exhaustion outcome without launching a second child. A persistence failure blocks further dispatch for repair.
Check suppression at common launch admission across restart, changed prompts, new run IDs, policy downgrade, and blocked membership batches.
Suppression does not change membership by itself. A blocked batch cannot remove it.

Transient failures retain the shipped original-plus-one path and five-second delay.
Stopped-writer proof, whole-batch settlement, positive usability, and the two-usable-non-facilitator floor remain mandatory.
Same-idea return is permanent. The next idea probes afresh without a relaunch hint.

### 3. Bound native adapter support

Keep the existing bounded Zcode recognizer and its disclosed native-positive limitation.
Include only the private capture facility: mode `0600`, bounded bytes, secret scrubbing, and visible truncation markers.
Do not expand Zcode grammar without the follow-up's required real native capture.

Add `codex.exec-json.turn-failed.v1` only when the effective adapter argv explicitly selects JSON output.
Bind one process invocation to its terminal `turn.failed` event and reject any later success or valid artifact.
For this release, accept only these exact source-mapped forms from `source-context/codex-native-source.md:139-166`:

```text
Your workspace is out of credits. Add credits to continue.
Your workspace is out of credits. Ask your workspace owner to refill in order to continue.
You hit your spend cap set in your workspace. Increase your spend cap to continue.
You hit your spend cap set by the owner of your workspace. Ask an owner to increase your spend cap to continue.
```

Reject generic quota text, usage-not-included text, context or rollout budgets, malformed JSON, quoted messages, and mixed terminal outcomes.
Use source-derived positive fixtures and adversarial negatives. Run a normal JSON smoke test for argv and envelope behavior only.
Label the support source-backed, not native-capture-verified. The default text launch remains diagnostic-only.
If the adapter cannot enforce and smoke-test explicit JSON output, ship Codex as diagnostic-only instead.
Kimi and Claude remain diagnostic-only without their own native evidence.

### 4. Defer peer and machine scope

Exclude only the identity named by the native invocation evidence. Never infer shared accounts from provider, gateway, model, or pool labels.
Defer the optional machine cache. No peer skip hint and no machine membership authority are added.
A later cache needs an exact agent and stable non-secret connection key, a valid future reset, and immutable-evidence revalidation.
Unknown reset must never create a persistent machine outage.

### 5. Record typed roles and apply the exact owner order

For v2 ideas, record original participant order and initial organizer and implementer identities.
Use optional `organizer:` for an ordinary participant organizer. Otherwise, a declared `facilitator:` supplies the organizer identity.
Do not infer an organizer from `author: user`. An ambiguous affected role blocks succession.

Add versioned `role_transition.v1` data to the v2 immutable batch schema.
Bind before/after membership, before/after organizer and implementer roles, affected drafts, evidence, and the active handover pointer.
Allow a role-only transition when the organizer is outside `participants:`. Never add a roster member or fabricate a confirmation record.
Validated role history outranks the old `IMPLEMENTATION.md` pin and the original designation.

Apply this algorithm after all writers stop and the complete candidate set is known:

1. Organizer loss transfers organizer to the current implementer.
2. Implementer loss transfers implementer to the current organizer, even a declared facilitator.
3. Loss of both transfers both roles to the first remaining current participant in original order.
4. One identity holding both roles counts as both losses. Ordinary participant loss transfers neither role.

For rule 3, exclude every dropped identity before selecting. Do not wrap, search after a holder, or use roster fallback.
If the selected counterpart or first remaining participant lacks positive usability, block. Do not silently skip to another agent.
A declared facilitator still supplies no non-facilitator floor seat. Below the floor or without a successor, escalate once.

### 6. Use one minimal handover form

Every organizer or implementer succession creates exactly one `role-handover/<transition-id>.md`.
The successor owns and updates that file. The predecessor's `IMPLEMENTATION.md`, rounds, reviews, signoffs, and findings remain byte-unchanged.
Use this minimal shape:

```markdown
---
idea: <slug>
transition: <immutable-history-id>
predecessor: <agent-id>
successor: <agent-id>
roles: [organizer, implementer]
phase: <current-phase>
execution_status: <current-status>
review_cycle: <N-or-none>
base: <../IMPLEMENTATION.md-or-prior-handover-or-none>
branch: <repo-and-branch-or-none>
head-commit: <current-sha>
created: <ISO-time>
---
## Takeover
## Resume state
## Progress
## Decision Log
## Validation evidence
```

At creation, copy the current mutable execution values from `base`. Record the next action and canonical resume paths under `Resume state`.
Before Phase 5, `base` is `none`; resume from the prompt and current phase artifacts.
During and after implementation, `base` names the previous `IMPLEMENTATION.md` or handover.
Later successions extend this chain. They never reset progress, review cycles, branch state, or validation history.

One `ResolveActiveExecution` helper selects the tip from validated role history and follows the base chain.
Driver dispatch, status, wait, organizer brief, review, goal check, and close must all use that resolver.
The newest filename is not authority. Parallel status documents are forbidden.
Completion updates the active handover, not the predecessor's frozen execution record.
Commit role history before successor dispatch. Until handover validation, permit only that handover action and emit one notice.

### 7. Recover drafts without rewriting predecessors

- A missing draft uses the normal selector over current eligible participants.
- A complete filed draft stays unchanged, including when signatures remain outstanding.
- A frozen FINAL or closed IMPLEMENTATION stays frozen.
- A well-formed incomplete draft moves byte-exactly to `role-handover/<transition-id>-<artifact-kind>.partial.md`.
- Record its original path, hash, and author in the transition. The successor then writes a new canonical draft from current artifacts.
- Existing signoffs remain only in the preserved file. Current participants must sign the new draft. No signature is copied or inferred.
- A malformed or truncated draft blocks for integrity repair. Relocation never makes it complete or valid.
- Participant-owned round and review files are never relocated or rewritten. Existing terminal re-evaluation rules handle them.

### 8. Make the single-reviewer proof role-aware

Compute prior reviewers from `Before` minus the before-implementer and facilitator.
Compute current reviewers from `After` minus the after-implementer and facilitator.
Require at least two prior independent reviewers and exactly one current independent reviewer.
Require the successor and sole reviewer in both membership sets, with `After` equal to current membership.
Every removed identity needs eligible typed evidence. Implementer replacement also needs the validated role transition.

The sole reviewer must survive and have a known snapshot model distinct from the successor.
Recompute this proof before commit, review, review consensus, goal check, and close.
The successor's earlier findings remain binding but cannot review its own implementation.
A fresh process of the sole reviewer may goal-check. Current signers, dissent, reservations, strict review, and current-tree evidence remain unchanged.

### 9. Required proof and delivery constraints

Test legacy golden behavior, first-attempt child counts, transient parity, reset separation, and blocked-batch suppression.
Cover restart, changed prompts, downgrade, stale pins, dual roles, simultaneous losses, unusable successors, and declared facilitators.
Add Codex JSON positives and adversarial negatives, including proof that the default path stays unsupported. Do not run a live exhaustion probe.
Inject faults at every role-history and handover write. Cover handover chains before Phase 5, during implementation, and after review findings.
Test complete, well-formed-incomplete, and malformed drafts separately. Assert byte preservation and current execution status throughout.
Test the former-reviewer successor, model diversity, retained dissent, fresh goal checks, and fresh current-tree acceptance evidence.

Place binding rules in existing mapped sections: §0 and Phase 0, §3, Phases 3–8, §5, and §9.0.
Update the skill summary, roster/protocol reference, headless reference, and artifact templates in the related skill worktree.
Apply identical reviewed hunks to all three protocol copies and staged core 2.20.0.
Compact without rule loss. Keep the 20,000-byte skill and 70,000-byte packet limits, including a hard meta-idea packet case.

### Process and protocol context

The owner answer quoted in `round-02/codex-1.md` states:

> **Native CLI fallback is authorized** for this idea (your option 1).
> Do not claim driver accounting or automatic membership authority for fallback steps.

This artifact uses that fallback after the recorded `budget_refused` in `source-context/resume2-round02-kimi-1-refusal.json`.
It creates no driver accounting, membership change, or automatic exclusion.

```text
context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none
```

I read the complete live authority and Appendix B. Independent `shasum -a 256` checks returned the stated hash for source and packet body.
No verdict in this file evaluates my own proposed implementation.
