---
idea: meta-protocol-change-participant-dropout
status: final
author: codex-1
drafted-by: codex-1
implementer: codex-1
consensus-date: 2026-10-08
participants: [codex-1, zcode-1]
---

## Purpose / user-visible outcome

A non-protected participant whose dispatched step fails twice is dropped permanently from that idea using the shipped exclusion machinery, while preserving every filed dissent and all safety gates. New ideas use the broader failure trigger by default through the existing quota_auto_exclude knob; existing saved policies retain quota-only behavior. The participant is probed again at the next idea.

The principal limitation is deliberate: a three-participant auto_implement run losing one participant still fails the existing precommit two-independent-reviewer gate. The tool records one actionable escalation and applies no reduction. Recorded owner authority can support an attended continuation, as for this idea; this feature introduces no automatic gate waiver.

## Context & orientation

Authority: source-context/ORGANIZER-BRIEF.md (SHA256 54448074f5a4dac2587ee3ed9c4e08866d84e50d42c02df7dc9b1e6af70651b5), codex-1 and zcode-1 round-01/round-02, and signed consensus.md. Kimi's two failed actual round-01 attempts are preserved; its explicit owner-authorized exclusion left these two participants. No Kimi position/signoff is invented and no Claude participant was used.

Bases: CLI d16ee9c (1.51.0) and skill 352a475 (2.15.0). CLI worktree participant-dropout and sibling participant-dropout-skill are the named owner worktrees. Transport github-pr; design PR https://github.com/feci/parley-deck-cli/pull/75. codex-1 is organizer, drafter, participant and implementer under explicit owner authority; zcode-1 supplies independent review. The proposals form one architecture family, both extending existing quota/membership machinery; model diversity does not make that common architecture independent proof.

Protocol context: full; source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e, no fallback reason. New protocol text itself is implemented only after this FINAL and the written implementation plan.

## Final plan / specification

### D1 — One reducer, one historical knob, a saved trigger version

Extend the shipped quota exclusion mechanism: the same whole-batch evaluator, leases, immutable kickoff/history, retained obligations, projections, receipts, recovery, notices and consumer refresh. No second membership system or global roster mutation.

The existing presence-aware `[defaults].quota_auto_exclude` and per-idea `quota_auto_exclude: false` remain the only boolean. Machine → deck → idea precedence stays unchanged. Its historical name now governs both triggers. New ideas record `{enabled, scope, trigger: "participant-failure-v1"}` by default, with scope `kickoff-and-mid-idea`; explicit false disables automatic reductions. Saved legacy `{enabled, scope}` records mean quota-only v1, retain their old serialized bytes and hashes, and never widen on resume or binary upgrade. Accept only the old shape and the new non-null exact trigger; reject unknown, duplicate, incomplete or malformed policy fields. Widening a saved policy requires the existing owner-bound revision and committed authority.

Legacy quota classification, including its reset predicate and owner-bounded zcode exception, stays untouched. The new trigger does not consult the quota recognizer for candidacy. Its rule id is `participant-failure.v1`, reset is unknown, and it gives no same-idea relaunch suggestion. The earlier native-positive waiver remains a historical limitation of legacy quota mode.

### D2 — Two actual attempts for the same failed step

Under the new trigger, an authorized non-protected participant's dispatched step gets the original attempt and exactly one relaunch after five seconds, with the original process ceiling unchanged. Nonzero child exit, provider error of any code (including 400/401/429/503), child crash/start failure, watchdog/hang timeout, and exit zero with missing or structurally invalid own output qualify. Existing first-output-watchdog retry uses that same second slot, never a third. Exec and ACP share the contract; no hidden session resumption or retry-to-reset worker.

A structurally valid artifact on either attempt, including a valid BLOCK/dispute/dissent, wins over an ordinary execution failure and prevents dropout for that step. Parent/operator cancellation, policy/budget/protocol/telemetry refusal, unresolved writer state and shared-artifact integrity failure never count as participant dropout evidence. A child timeout is distinct from cancellation of its parent. An invalid signoff that alters somebody else's shared blocks is preserved and stops for repair; an invalid own output may be preserved and retried without erasing other content.

The two-attempt bound is durable and keyed by idea, agent and stable logical step/artifact. Run IDs, restarts, changed prompt inputs or a new driver process do not replenish it. Both immutable invocation IDs and retry linkage survive. A started attempt with no terminal outcome remains blocked until the existing host/boot/PID stopped-writer proof resolves it; do not launch alongside it or invent terminal success. Use a small retry binding/helper at the existing execution seam, with existing telemetry wherever possible, not a general retry service. A genuinely later protocol step is a distinct step; input edits are not.

Kickoff readiness uses the same two-attempt rule within the proposed batch, before the first authoritative participant list or manifest. A standalone preflight reports observations and never applies exclusions. Undispatched readiness/setup uncertainty is not fabricated child evidence. After an idea exists the restart bound is per idea; a new idea probes everyone afresh.

### D3 — Atomic batch, usable floor and protected roles

Settle once per whole batch, after every original/retry writer stops. Apply every eligible candidate together or none; completion order cannot change the result. The fixed floor is two distinct usable participants: the usable designated/pinned implementer when present, plus at least one other usable non-organizer. Before a pin, require two usable non-organizers including any live designee. The organizer never contributes an extra floor seat. A declared facilitator never counts even when participating; an undeclared organizer must already be protected as designee/pin/drafter or must declare the existing facilitator field for product enforcement. No new organizer field.

Usability requires positive execution or validated current phase evidence, never role identity alone. Existing valid signoffs can seed a signoff batch; validated IMPLEMENTATION.md can seed a review batch's unlaunched implementer. Equivalent relevant validated evidence may seed unlaunched peers, but an unresolved failure cannot be relabeled usable. Conservative unresolved-usability escalation is acceptable. A failed designated/pinned implementer cannot be absent from the usable survivor set.

Protected roles stay protected: declared facilitator, per-idea designee, implementation pin and any started consensus/FINAL/review-consensus drafter. A protected failure/candidate prevents automatic reduction and uses the existing recorded three-exit owner gate. Protection does not itself create a usable floor slot.

Keep `membership.CheckGates` BEFORE transition commit: reviewer minimum, LE-7/LE-11, independent goal-checker, model diversity and strict-gate requirements are re-evaluated and never waived. Thus **a three-participant auto_implement idea losing one participant still blocks the reduction because it leaves only one independent reviewer, despite meeting the two-member floor**. This is the principal limit of this change, required by brief point 4. The rule will not claim that every such run can continue unattended.

A blocked batch applies nothing and produces one actionable blocking escalation for that settled decision, naming candidates, verbatim evidence, usable arithmetic and options: (a) authorize a separate eligible model-diverse reviewer process (recommended when review independence is needed), (b) record an attended evidence-backed continuation/close with reduced reviewer count, or (c) pause/abandon. A named example is a separately authorized Claude or Codex process; the tool chooses none. An owner may give standing per-idea authority for an attended path, as this brief does; that is not a new automatic product gate override. A later distinct failed batch gets its own decision: if the first reduction applied but the second would breach a gate/floor, keep the first history and block the second, one notice for each settled batch.

### D4 — Permanent per-idea dropout, retained dissent and evidence

Only new-trigger dropout is permanent for the idea. History, not `excluded:` display markers, binds the dropped ID. Every rejoin path refuses it, including `quota revise`, recovery, plain policy-off membership edits, late catch-up and changes that disable/downgrade the policy after the drop. Include kickoff-dropped IDs even though they never entered the historical quorum. Legacy quota-only exclusions keep the shipped owner-confirmed return behavior. The next idea probes afresh. An owner can abandon and open v2, but cannot use this idea's ordinary revision path to erase the terminal membership fact.

Preserve all filed artifacts, invalid partial output, prior invocation evidence, every retained veto, DISPUTED claim and finding. Excluded historical signers remain known; only current members are required new signers. A valid filed BLOCK still blocks, not a malformed unknown-author error. No dropped author may append a signoff or fabricated withdrawal. For a permanently dropped author, retained dissent requires an explicit owner ruling quoted into the next artifact (or abandonment/v2); absence never withdraws it. Independent disposition duties continue for findings and disputes.

Before committing a drop, bind both failed attempts in the same immutable transition evidence: idea/agent/logical step, invocation IDs and linkage, terminal status and exit/watchdog class, structural validator reason, UTC observation times, and verbatim decisive failure facts/excerpts with only secret redaction and visible truncation. Supervisor facts authorize the decision; participant prose, tool quotes, dissatisfaction and elapsed time alone do not. Full raw logs remain private. Existing stopped-writer and projection recovery apply unchanged.

Publish one non-blocking notice per applied batch. Owner-edited/archived notices remain owner-owned and never authorize membership. Delivery failures are visible diagnostics, not rollback or false delivery claims. Narrowly repair missing-safe-inbox handling and a detailed stderr fallback for blocked kickoff batches, and replay interrupted kickoff notice publication through the existing receipt machinery. No exactly-once delivery claim across an unavailable filesystem; the contract is one idempotent publication attempt/notice with explicit delivery diagnostics.

### D5 — Scope and exact change inventory

The normative hunks below apply to all three COOPERATION.md copies: `parley-deck/COOPERATION.md`, `internal/protocol/defaults/COOPERATION.md`, and `../participant-dropout-skill/skills/parley-deck/references/COOPERATION.md`. Preserve each existing generic/project header zone; match the normative hunks, do not overwrite independent headers.

| Hunk | Existing anchor | Required delta |
| --- | --- | --- |
| P1 | §0 defaults sentence naming quota_auto_exclude | Explain the historical knob selects automatic exclusion with a kickoff-frozen trigger; no extra knob. |
| P2 | Phase-0 quota_auto_exclude comment | False opts out of automatic quota/participant-failure exclusion; saved scope/trigger are immutable authority projections. |
| P3 | Phase-5 protected-role paragraph | Replace quota-only trigger wording by recorded automatic-exclusion failures, retaining the designee/pin three-exit gate. |
| P4 | §5 membership/return exception | Name versioned automatic exclusion, preserve history/current/known and obligations, state permanent participant-failure no-rejoin versus legacy quota owner-confirmed return. |
| P5 | §9.0 Quota auto-exclusion block | Retain and label the legacy predicate, add compact D1–D4 contract: new default, two attempts, control failures, evidence, batch/floor, protected roles, precommit gates, notices, migration and permanent return. Use existing section structure; no applicability-map expansion. |
| P6 | Existing phase-3/6/7 §9.0 references | Adjust terminology only where quota-only wording would contradict P4/P5; no change to quorum/close duties. |

CLI product files are limited to these existing seams and a bounded shared helper:

- `internal/quota/{quota,record,history,revision}.go`: compatible policy/evidence shape, version-bound validation, permanent-history query, notice and return checks.
- `internal/protocol/{quota,quota_manual}.go`: saved trigger projections, policy matching and all manual/catch-up return refusals.
- `internal/membership/{membership,notice}.go` (and existing revision seam if needed): usability/protected integration, actionable block text, history-derived permanent return guard, kickoff notice replay. Keep `gates.go` gate semantics.
- `internal/telemetry/dropout.go` plus `record.go` only for narrowly required immutable step/evidence metadata; existing quota recognizers untouched.
- `internal/runner/{runner,acp,quota,telemetry,launch}.go` and `dropout.go` shared retry helper: bound attempts, classification/validation, exec/ACP parity, preserved output and result handoff.
- `internal/app/{app,preflight,preflight_liveness,quota,quota_signoff,consensus_request_signoffs,driver_consensus,driver_impl}.go`: new-idea policy selection, kickoff probes, single-step/signoff/draft/goal-check wiring and surfaces. `wait`/`organizer`/status change only if their shared quota surface is insufficient.
- `internal/runcontrol/runcontrol.go`/`internal/runmanifest/manifest.go` only for initial policy/notice receipt plumbing. `internal/config/runtime.go` explanatory text only; reuse its presence-aware boolean.
- Focused `dropout_test.go` files in quota, telemetry, runner, app, membership, protocol and consensus, plus existing adjacent fixtures where the new default requires explicit legacy policy. Preserve existing legacy quota expectations and all integrity/packet guards. Exact internal helper names can follow the narrowest shared seam; the behavioral scope above is fixed.
- Documentation: protocol changelog, CLI usage/release notes, skill `SKILL.md` and `references/ROSTER_AND_PROTOCOL.md`, compatibility/version/release metadata and generated manifests. Versions: CLI 1.52.0, skill/core 2.16.0. No Windows/alias, D6, roster/model or broad retry-service changes.

### D7 — This run’s phase and close authority

Kimi's two failed actual round-01 attempts and the owner-authorized exclusion are preserved in kickoff/source-context/inbox. Kimi authored no position and is not imputed a signoff. Active participants are codex-1 and zcode-1. The brief explicitly permits this pair to continue and pre-authorizes attended close after both final review-consensus signoffs, no open CRITICAL/MAJOR and current-tree criterion evidence. Both participants agree this leaves no unresolved owner boundary. This is standing authority for this run, not driver auto-close or a product gate waiver. Zcode performs independent code review and a fresh goal/check or channel verification as applicable; codex never reviews itself. Maximum five fix-up cycles, and stopping judgment still applies.

Driver gaps remain recorded. CLI `consensus draft/status/signoff/finalize`, `status` and `wait` remain preferred; configured CLI fallback is used only for steps the known D6/driver gaps prevent. Never prune or declare worktrees, repair accounting, change models or add Claude to this run.

## Observable acceptance criteria

| ID | Observable requirement / independent refutation target |
| --- | --- |
| AC1 | New-trigger two-failure fixtures for 400/401/429/503, nonzero exit/crash/start failure, watchdog/timeout and missing/invalid own artifact produce one eligible batch reduction when every gate passes. The first valid artifact or second successful attempt prevents it; valid BLOCK/dissent is success. |
| AC2 | Original plus one relaunch after five seconds with the same ceiling; watchdog retry shares the slot; exec/ACP agree. Restart between failures and after second failure launches at most two total and eventually commits at most one reduction; input/run changes do not reset the step. |
| AC3 | Parent cancellation, control-plane refusal, telemetry/integrity error and unresolved live/unknown writer never authorize dropout; shared-signoff tampering is preserved and blocks. Stopped-writer recovery never invents success. |
| AC4 | Order-independent whole-batch 3→2 applies only with two positive usable non-organizer slots including required implementer; 3→1/2→1 applies nothing. A protected candidate or unusable implementer blocks; no role-only seat. |
| AC5 | Applicable reviewer/diversity/strict/goal-check gates veto before commit. In particular auto_implement 3→2 remains blocked with one independent reviewer. One actionable escalation shows candidates/evidence/arithmetic and owner options; a missing inbox is created safely, unwritable destination prints the decision diagnostic. |
| AC6 | Every committed drop binds both immutable attempts and step identity, preserves raw partial output privately and prior canonical files/history/obligations; valid historical BLOCK remains TriageBlocked, DISPUTED/finding duties remain. |
| AC7 | Same-idea rejoin refuses new-trigger dropped IDs through every automatic/manual/owner/catch-up/recovery path, even after opt-out/downgrade and for kickoff candidates. Next idea probes the ID normally. Legacy return remains unchanged. |
| AC8 | New ideas default to the new trigger; machine/deck/idea false opts out; legacy absent-trigger policy bytes/hashes and behavior remain unchanged on upgrade/resume. Unknown/null/duplicate fields fail closed; saved widening requires owner authority. |
| AC9 | Both kickoff and mid-idea round/review/signoff and relevant single-step paths follow the contract; standalone preflight never applies; unusable unobserved peers are not fabricated. Protected draft/implementer failures escalate. |
| AC10 | status, wait, organizer brief and signoff/driver consumers agree on current membership, permanent drops and pending state. Read-only calls do not repair; wait exit contract remains. |
| AC11 | Exactly one idempotent notice publication per applied batch; interrupted kickoff publication is replayed from history/receipt, owner-edited/archived notices are preserved, and unavailable delivery is diagnostic. Later-batch failures keep earlier history and get their own decision. |
| AC12 | All three reviewed normative protocol hunks agree; phase 0/5/8 full/facilitator rendering and existing packet-size/drift tests pass without weakening limits. Skill guidance and opt-out examples match implementation. |
| AC13 | Focused behavioral tests, full host `go test ./... -count=1 -timeout 45m`, build/vet and skill tests/manifests pass on the reviewed tree. Current-tree independent reviewer evidence covers AC1–AC12; no open CRITICAL/MAJOR remains at close. Inherited Windows limits are disclosed, not passed. |
| AC14 | After authorized close: CLI 1.52.0 / skill 2.16.0 GitHub releases, both Homebrew formulae, skill-only WinGet PR, every managed and four named generic installs hash-verified; separate available reviewer's channel verification. Core 2.16.0 derives from staged 2.15.0 plus exactly P1–P6, and npm/core remain owner commands. |

## Idempotence & recovery

Immutable kickoff/history remains the membership authority. Mutable prompt/manifests/markers and notices are projections; read-only surfaces never repair them. `quota recover` validates history and stopped-writer state, then reconciles interrupted projections/receipts. It cannot change membership history, waive gates, make a foreign/unknown writer safe or re-add a permanently dropped ID. Save both attempt identities and evidence before the membership commit. The retry budget follows the idea/agent/logical step across driver runs and changed inputs; never grant a third attempt. Replaying a settled step or committed batch cannot create another child or transition. A live or unresolved attempt blocks until authentic terminal/stopped-writer evidence resolves it. Preserve private logs and invalid partial output before retry; never overwrite another participant's shared content.

Block/escalation is atomic per settled batch, and an earlier valid reduction is not rolled back by a later batch's failure. Notice publication is idempotent and diagnostic on delivery failure. A missing kickoff publication receipt is replayable; owner-edited/archived copies remain owner-owned. No timer/reset-based return. Explicit owner-bound policy revisions may widen legacy policy but cannot erase new-trigger permanent dropout. Owner ruling on retained dissent or a new v2 idea are the documented exits.

Implementation and release are not complete merely because the process stops or a budget is consumed. Maximum five fix-up cycles with stopping judgment; no unresolved CRITICAL/MAJOR or stale/self-issued criterion evidence may establish close. Owner-only npm OTP/core TTY commands remain pending owner actions, not simulated publication.

## Known risks / de-risking

The most visible limitation is the preserved precommit reviewer gate in D3. Future owners may pre-authorize an attended path; this release adds no automatic reviewer-count waiver. If that is still too interruptive, a separate candidate `participant-dropout-review-gate-timing` can consider splitting hard membership gates from recorded review escalations, with its own owner decision and independent design.

Two attempts can waste one relaunch on a long quota outage; a blip exceeding five seconds or a consistently too-short timeout can still remove a healthy provider for this idea. Evidence, whole-batch floor/gates and next-idea probing bound that risk; no automatic same-idea return. Conservative recovery may require original-host action. Unknown evidence blocks instead of guessing. Extra bounded retry metadata is necessary to make the cap durable; previous LOC estimates are unverified, not a requirement to omit safeguards.

The tests and independent current-tree refutation target restart duplication, secret-safe failure evidence, malformed own output versus shared-content tampering, retained vetoes, permanent return after policy changes, positive floor evidence and precommit gate behavior. Real provider timing, cross-host/PID namespaces and inherited Windows filesystem limits remain disclosed. Prior native-positive quota recognition, aliased-deck plain-edit guidance and D6 accounting are separate follow-ups. No old waiver is claimed as a new test pass. Keep all existing packet-size/drift limits; compact the protocol text instead of weakening them.

## DISPUTED claims

None remains from design. No acceptance criterion depends on an unresolved disputed factual claim. If implementation/review finds one, preserve it and resolve by evidence or owner ruling under §15; do not count votes as verification.

## Release contract

After attended close, merge to main and release CLI 1.52.0 and skill 2.16.0 on GitHub, update both Homebrew formulae, submit the skill-only WinGet PR and retain the CLI WinGet hold while Windows is experimental. Install the released skill into every managed runtime including existing undetected targets, then each of `~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck` and `~/.config/opencode/skills/parley-deck` via `parley-deck-skill install --target generic --dest <dir> --force`. Verify every SKILL.md by hash; no roster/model changes.

Build staged core 2.16.0 from `/Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md` (SHA256 0d81fd807114e4ff67f5fa98bba73622b58c86d66a2bf6b57d1def68cd099c09) plus exactly the reviewed normative hunks, preserving generic zones. Give zcode-1 a separate short channel-verification task. Final note: `inbox/codex-1-to-user_meta-protocol-change-participant-dropout_released.md`, with channel evidence, limitations/deferred items, usage and exact owner-only `! npm publish ...` / `! parley protocol publish ...` commands. If 2.15.0 remains unpublished, list its commands first. Never bypass OTP/TTY gates.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
- Controlling brief: ./source-context/ORGANIZER-BRIEF.md
- Kimi failure evidence: ./source-context/kimi-round01-failure-evidence.md
- Organizer/driver record: ./organizer-notes.md
- Prior implementation: ../meta-protocol-change-quota-auto-exclude/FINAL.md and IMPLEMENTATION.md
