---
agent: zcode-1
idea: meta-protocol-change-credit-exhaustion-failover
round: 1
date: 2026-10-29
---

## Summary

I propose one new typed failure class, "provider exhaustion", that drops an agent on its first failed attempt, and a role-succession rule that implements the owner ruling of 2026-10-29 verbatim. The class reuses the shipped recognizer framework, the exclusion reducer, the immutable history and the retained-dissent rules. It adds no new policy knob and no new protocol heading, so the packet map and the saved policy bytes stay stable. The main open points are the gateway recognizer provenance, the drafter takeover detail, and the byte budget inside §9.0.

## Proposed approach

### A. Current state (verdicts, §15)

All verdicts in this section are PRIMARY. I read each file and quote it at HEAD `f3b8a4f` of this worktree.

- **Only zcode has a production quota recognizer.** `internal/telemetry/quota.go:20` returns `{"zcode", "supported-owner-deviation", ...}`; the same table lists codex, kimi and claude as `"diagnostic-only"` (`internal/telemetry/quota.go:21-23`). `ClassifyQuota` returns early for every other adapter (`internal/telemetry/quota.go:71-74`).
- **A short stated reset does not qualify today.** `internal/telemetry/quota.go:125-127` reads `if reset.Sub(in.ObservedAt) < quota.MinimumReset { e.Reason = "known reset below 60 minutes"`. `quota.MinimumReset = time.Hour` (`internal/telemetry/quota.go` companion `internal/quota/quota.go:14`).
- **A failed participant gets exactly two attempts.** `internal/quota/dropout.go:43` requires `len(f.Attempts) != 2` to fail validation, and `internal/runner/dropout.go:200` loops `for len(attempts) < 2 {`. The retry fires after `const ParticipantRetryDelay = 5 * time.Second` (`internal/runner/dropout.go:22`). So the observed gateway text "All claude accounts have exhausted their quota (reset after 5m)" (ORGANIZER-BRIEF.md, Design target 1) burns one futile retry today, then drops the agent.
- **A protected role blocks the whole batch.** `internal/quota/quota.go:197` sets `d.Block = "protected role %s failed or has unresolved usability; re-designate, record implementer_waived, or set implementer: none with owner confirmation; no exclusions applied"`. The floor block is `internal/quota/quota.go:224` (`quota floor: ... < %d`, `Floor = 2` at `internal/quota/quota.go:13`).
- **Organizer failover is out of scope today.** `parley-deck/COOPERATION.md:1013-1014` says "Organizer failover or replacement agents stay outside this rule."
- **A re-pin mechanism already exists.** Phase 5 (`parley-deck/COOPERATION.md:462`) defines `implementer_reassigned: {old-agent-id} to {new-agent-id} — {reason} — confirmed {date}`; `internal/protocol/implementer.go:237` parses it.
- **The single-reviewer exception takes frozen rule ids.** `internal/membership/review_gate.go:16-19` lists `legacyResetRule = "quota.named-reset-ge-60m.v1"`, `legacyAccountRule = "quota.account-exhausted-no-reset.v1"`, `legacyAllowanceRule = "quota.allowance-ge-24h-no-reset.v1"`, and `legacyProvenance = "zcode.stderr.owner-deviation.2026-10-04"`. Distinct snapshot models are required at `internal/membership/review_gate.go:136-137`.
- **Size budgets are nearly full.** I ran `go test ./internal/app -run TestLiveDeckFacilitatorPacket -count=1 -v`; it logged `facilitator body=68874 B ... guardrail=70000 B` (`internal/app/facilitator_packet_live_test.go:147`; guard at `:21`). `wc -c` gives 19,482 bytes for the skill core against `const CORE_BUDGET = 20000;` (`test/lean-organizer.test.js:19`) in the skill worktree. The packet-mode shadow for this launch is 78,701 bytes over 36 included blocks (`.parley-runtime/credit-failover/resume-packet-phase1.json`, `shadow`).

### B1. Exhaustion evidence class

Define one class: **provider exhaustion**. A terminal provider record states that the account's credits, tokens, quota, or a named usage allowance is exhausted, depleted, or reached. The record must come from the invocation's own terminal evidence through a recognizer with recorded provenance. Model prose, tool output and quoted text never qualify; this keeps the shipped integrity rule.

- Rule id: `participant-exhaustion.v1`. It is a RuleID inside the existing `participant-failure-v1` trigger, not a new trigger. Saved policies keep their bytes and hashes unchanged (`internal/quota/quota.go:27-52` unmarshals exactly two or three fields and freezes the trigger).
- Recognizers: keep `zcode.stderr.owner-deviation.2026-10-04` as is. Add one owner-bounded gateway recognizer for pool exhaustion wording, with the observed text recorded as a fixture. The owner's words name this exact case, so the prior bar ("not positive forms until the gateway's source establishes their cause", quota-auto-exclude FINAL §4) is met by owner direction, the same way the zcode exception was. Table status: `supported-owner-deviation`, with the limitation written beside it (`internal/telemetry/quota.go:18-25`).
- Adapters without provenance stay diagnostic-only and keep the two-attempt path.

### B2. First-attempt stop

- After the first terminal failure, classify the terminal record before any retry. The seam is the retry loop at `internal/runner/dropout.go:200`: insert the check between the failed attempt and the wait. If the record is eligible exhaustion, do not launch the second child. Do not probe the agent again in this idea. Do not reconnect.
- Evidence shape: one `FailedAttempt`, same binding discipline as today (idea, agent, logical step, invocation id, terminal class, scrubbed excerpt). Add a sibling validator next to `ValidateFailureEvidence` (`internal/quota/dropout.go:41-64`): require exactly one attempt, `RuleID == participant-exhaustion.v1`, and the recognizer provenance string.
- The drop still settles through the whole-batch evaluator (`quota.Evaluate`, `internal/quota/quota.go:147`) inside `membership.Settle` (`internal/membership/membership.go:282`). All-or-none, order-independent, one transition, one notice — unchanged.
- Kickoff readiness uses the same rule: an exhausted probe in the proposed batch is a one-attempt candidate (`internal/app/preflight.go:422` builds the members; `internal/app/consensus_request_signoffs.go:217-218` wires signoff batches).

### B3. Reset semantics

- The reset value is recorded, never a gate for this class. A stated machine reset of any length, five minutes included, qualifies. Reason: the shipped retry fires five seconds after the first failure (`internal/runner/dropout.go:22`), so any machine reset longer than that delay proves the retry cannot succeed. The owner intent is exactly this: an agent out of credits or tokens is not retried.
- The 60-minute threshold stays where it is today: the legacy predicate (`quota.named-reset-ge-60m.v1` and friends) is untouched.
- Fail closed on bad resets: an unparseable, contradictory, or mixed reset on an exhaustion record keeps the two-attempt path. Reuse `quotaReset` (`internal/telemetry/quota.go:154-253`) with its every-record agreement and timezone rules.
- With no stated reset, only explicit account exhaustion or a reached allowance of 24 hours or longer qualifies — the legacy wording tests apply unchanged.
- The relaunch hint stays "reset plus five minutes, provider estimate" (`Evidence.RelaunchHint`, `internal/quota/quota.go`).

### B4. Peer-account scope

One exhaustion record excludes only the agent whose invocation failed. The CLI cannot map agents to accounts or gateway connections reliably; credentials are opaque, and several roster adapters share one gateway (the machine roster routes kimi through omniroute and hermes through openrouter, `~/.parley/agents.toml`). Automatic peer exclusion from one agent's record would guess at that mapping. Do not guess.

The machine record below may carry a provider key. Preflight uses it only as a skip hint for agents whose configured model names the same provider. A skip is not an exclusion: the agent reports unavailable with a reason, and the floor arithmetic stays honest.

### B5. Optional machine record

Add a small machine-level skip record, optional and auditable:

- Fields: provider key, agent id, rule id, observed time, reset time, scrubbed excerpt, invocation id. No secrets, no tokens, no full logs.
- Location: the existing private quota authority pattern (`internal/quota/authority.go`, `internal/quota/receipt.go`).
- Effect: later kickoffs and preflights skip the agent's probe until reset plus five minutes, then probe afresh. It never writes `agents.toml`. It never gates membership by itself. The next idea always probes afresh.
- Absence of the record changes nothing; the one-attempt drop inside an idea works without it.

### B6. Role succession (owner ruling, binding)

The ruling, verbatim from `00-prompt.md`:

> "no ked vypadne organizer, tak implementer ho preberie, ak vypadne implementer tak organizer ho preberie, ak vypadne aj organizer aj implementer, tak ho preberie prvy participant a participanti ked vypadnu, tak sa nedej nic, len proste vypadnu"

Implement it exactly:

1. The organizer drops: the implementer takes over the organizer role.
2. The implementer drops: the organizer takes over the implementer role.
3. Both drop, or one agent held both roles (this idea's shape): the first remaining agent in `participants:` order takes over both roles. Here that successor is kimi-1 (`participants: [codex-1, kimi-1, zcode-1]`).
4. A plain participant drops: no succession. It only drops.

Details the ruling leaves open, my positions:

- **Product split.** Implementer succession is driver-enforced: on an exhaustion drop of the designee or pin, the driver writes `implementer_reassigned: <dropped> to <successor> — exhaustion succession, owner ruling 2026-10-29 — confirmed 2026-10-29` in the same transition, then dispatches to the successor. The existing parser is `internal/protocol/implementer.go:237`; the protocol anchor is `COOPERATION.md:462`. Organizer succession cannot be driver-enforced: the organizer is the driving process itself. Keep it a protocol rule plus a recorded transition; `parley continue` accepts the changed organizer identity and names the successor per the ruling.
- **Protected-role block yields for this cause.** Today `internal/quota/quota.go:197` blocks the whole batch when a protected role failed. For an exhaustion candidate with a ruling-successor, the batch applies the drop and records the succession in the same transition instead of blocking. The three-exit gate stays for every other protected-role failure. Extend `quota.Roles` with the organizer identity so `Evaluate` can derive the successor.
- **Facilitator default yields.** The rule "a declared facilitator does not implement" (Phase 5, `COOPERATION.md:462`) gets one exception line: it yields to the ruling for exhaustion succession. The 00-prompt already states this.
- **Started drafters.** An exhaustion drop removes drafter protection for that cause. If the dropped drafter is the organizer or implementer, the ruling's successor takes the draft over. If a plain participant started a draft and drops, the ruling says "nothing happens": the draft duty falls back to today's Phase 4 fallback chain over the current set. In both cases the partial draft is preserved, identified as incomplete, and never promoted because the file exists (shipped partial-artifact rules, D4 of participant-dropout).
- **No succession below the floor.** If the named successor is not usable, or the floor would break, apply nothing and escalate once with a blocking owner note.

### B7. Handover

The successor continues from the canonical artifacts and the branch state. It records the takeover in `IMPLEMENTATION.md` (Decision Log, or the phase artifact before Phase 5) and in one inbox notice. No artifact of the dropped agent is rewritten; git attribution survives. The transition record binds: dropped id, cause, ruling clause, successor id, timestamp, invocation evidence.

### B8. Dissent

An exhaustion drop never withdraws filed positions. The dropped agent's ❌, `DISPUTED` claims, and findings stay binding (`COOPERATION.md:781` already states this for exclusions; repeat it in one line inside the new block). Retained dissent of a permanently dropped author still needs an owner ruling quoted into the next artifact, or abandonment/v2. A drop of a dissenter is never a way to remove dissent.

### B9. Review independence and floor

- A successor implementer never reviews its own code. Its earlier filed findings stay filed.
- Reviewer counts are re-evaluated after succession, per track. The shipped single-reviewer exception (`COOPERATION.md:981`; `internal/membership/review_gate.go:26`) applies only when its causal proof holds. Add `participant-exhaustion.v1` to the accepted rule ids at `internal/membership/review_gate.go:16-19`, with the same evidence validation; nothing else about the exception changes.
- Model diversity still binds where required: the implementer and the sole reviewer must have known distinct snapshot models (`internal/membership/review_gate.go:136-137`). In this idea, if codex-1 drops and kimi-1 succeeds as implementer, the sole remaining reviewer is zcode-1; kimi-1 (omniroute/k3.1-pro) and zcode-1 (zai/glm-5.3) are distinct, so the exception can qualify on models — but only through the mechanical causal proof, never by assertion.
- The floor stays two distinct usable non-facilitators (`quota.Floor = 2`, `internal/quota/quota.go:13`). Succession never creates a floor seat. Below the floor: one blocking escalation, nothing applied.

### B10. Legacy policy

- The legacy quota predicate, its 60-minute threshold, its no-reset variants and the zcode owner-deviation recognizer stay untouched.
- The trigger enum stays `participant-failure-v1` only. The exhaustion class is a RuleID, so saved policy bytes, hashes and resume behavior never change (`internal/quota/quota.go:73-79` rejects unknown triggers; adding none avoids that whole class of migration).
- Knob-off ideas keep today's behavior: no automatic reduction. (The owner's words ask for the drop; the default-on path covers new ideas, which record the trigger at kickoff.)

### B11. Recovery

- Same-idea return: none. The drop is permanent for the idea, like participant-failure dropout. No timer, no polling, no automatic rejoin.
- Next idea: fresh probe. The machine record skips probes only until its recorded reset passes.
- History: the transition is immutable; read-only surfaces report it and never repair it; `quota recover` reconciles projections only. A replay of the transition is idempotent, one notice.
- Rollback of the feature: revert the CLI commits; protocol text reverts through a new meta-protocol-change idea only (§7).

### B12. Protocol hunks, files, and the byte budget

Protocol hunks, identical in all three copies (`parley-deck/COOPERATION.md`, `internal/protocol/defaults/COOPERATION.md`, skill `skills/parley-deck/references/COOPERATION.md`):

1. §9.0, inside the existing Automatic exclusion block (`COOPERATION.md:912ff`): add the exhaustion first-drop contract — class, one attempt, reset recorded never gating, no reconnect, machine record as a skip hint, succession per the ruling, retained dissent. **No new heading**, so `parley-deck/meta/packet-applicability.yaml` keeps its exact block list.
2. §9.0 history line (`COOPERATION.md:1013`): replace "Organizer failover or replacement agents stay outside this rule." with the bounded carve-out: exhaustion succession per the owner ruling; every other organizer failover stays outside.
3. Phase 5 (`COOPERATION.md:462`): one sentence — an exhaustion drop of the designee or pin succeeds per §9.0; the three-exit gate does not open for this cause; the facilitator default yields to it.
4. §5 (`COOPERATION.md:781` area): one clause — succession never withdraws dissent.
5. `meta/protocol-changelog.md`: one entry naming this idea and the ruling.

CLI files: `internal/telemetry/quota.go` and a new gateway recognizer file with fixtures; `internal/runner/dropout.go` (pre-retry check); `internal/quota/{dropout,quota,history,record,notice}.go`; `internal/quota/{authority,receipt}.go` (skip record); `internal/membership/{membership,review_gate}.go`; `internal/protocol/implementer.go` (reuse, no change expected); `internal/app/{preflight,preflight_liveness,consensus_request_signoffs,driver_impl,driver_consensus}.go`; `internal/protocol/defaults/COOPERATION.md` under the drift test.

Skill files: `skills/parley-deck/SKILL.md` (exclusion summary: the owner never re-probes an exhausted agent in the same idea; succession order in one line), `skills/parley-deck/references/ROSTER_AND_PROTOCOL.md` if it names the trigger, and the bundled COOPERATION.md.

Byte budget: the facilitator packet has 1,126 bytes of headroom (68,874 of 70,000, measured above) and SKILL.md has 518 (19,482 of 20,000). The additions must be funded by compaction inside the same always-included blocks, with no rule lost — the same discipline as review-gate-timing D5. The meta-idea packet variant is already over any 70,000 figure (78,701 B shadow for this launch); keep it net-neutral. I do not verify a compaction here; that is implementation work with its own check.

### B13. Regression tests

- `internal/telemetry/quota_test.go` and the new recognizer test: positive gateway fixture with a five-minute reset; negatives — bare 429, generic "quota exceeded", quoted model prose, tool-emitted text, missing or mixed reset, unknown body fields.
- `internal/runner/dropout_test.go`: exhaustion on attempt one launches no second child; a transient failure still gets two; a crash between classification and transition recovers without a new child.
- `internal/quota/dropout_test.go`: one-attempt validator positives and negatives (two attempts, wrong rule, missing binding).
- `internal/quota/quota_test.go`: exhaustion candidate with the designee dropped applies the drop and records succession instead of blocking; floor still blocks 3→1; knob-off applies nothing.
- `internal/membership/{dropout_test,review_gate_test}.go`: the new rule id as accepted cause; retained ❌ yields blocked, not malformed; drafter takeover and fallback.
- `internal/app/{preflight_liveness_test,preflight_test}.go`: skip-until-reset from the machine record; fresh probe after the reset; no `agents.toml` write.
- `internal/protocol` drift test; `internal/app/facilitator_packet_live_test.go` size; skill `test/lean-organizer.test.js` core budget; full `go test ./... -count=1 -timeout 45m` (the idea's `checks:`).

## Existing alternatives

§15.6(a). The mechanisms this proposal would build by hand, and the closest shipped toolchain component:

- ALT-1 exhaustion recognizer — inherited: the per-adapter recognizer framework and `refineNativeQuota` (`internal/telemetry/quota.go:95-149`) plus the zcode precedent (`internal/telemetry/quota_zcode.go`). Extend, do not rebuild.
- ALT-2 first-attempt stop — inherited: the durable step ledger and retry loop (`internal/runner/dropout.go:200`, `appendParticipantAttempt` at `:338`). Add a pre-retry classification; no new retry service.
- ALT-3 one-attempt evidence — inherited: `FailedStep`/`ValidateFailureEvidence` (`internal/quota/dropout.go:10-64`). A sibling validator, same binding fields.
- ALT-4 no new knob — inherited: the frozen trigger enum and policy unmarshal (`internal/quota/quota.go:44-79`). Constraint-forced: saved policy bytes and hashes must not change.
- ALT-5 succession record — inherited: `implementer_reassigned` (`internal/protocol/implementer.go:237`; `COOPERATION.md:462`). The ruling supplies the standing authority and the reason string.
- ALT-6 batch drop and transition — inherited: `quota.Evaluate` (`internal/quota/quota.go:147`), `membership.Settle` (`internal/membership/membership.go:282`), immutable history, notices.
- ALT-7 machine skip record — inherited shape, constraint-forced novelty: no toolchain component tracks provider outages. Reuse the private authority/receipt store pattern (`internal/quota/authority.go`, `internal/quota/receipt.go`); keep it a hint, never membership authority.
- ALT-8 single-reviewer cause — inherited: the frozen rule-id allowlist and distinct-model check (`internal/membership/review_gate.go:16-19`, `:136-137`).
- ALT-9 protocol placement — inherited: the always-included §9.0 block (`COOPERATION.md:912ff`; map entry `### 9.0 ... include: always` in `parley-deck/meta/packet-applicability.yaml`). No new heading, no map change.
- ALT-10 organizer succession enforcement — null. No CLI verb hands a live organizer process to another agent. Constraint-forced: protocol rule, recorded transition, notice, and acceptance on `continue`. Consulted: `internal/driver/loop.go` (run-local lock), `internal/app` run/continue surfaces, `COOPERATION.md:1013`.
- ALT-11 gateway recognizer provenance — null. No omniroute or gateway source is located in this worktree, so native provenance cannot be established the way the zcode deviation was from CLI behavior. Constraint-forced: an owner-bounded fixture recognizer, status `supported-owner-deviation`, with the recorded observed text and the same every-record agreement. Consulted: `internal/telemetry/quota.go:18-25`, `internal/telemetry/quota_zcode.go`, quota-auto-exclude FINAL §4, the open follow-up `quota-zcode-native-exhaustion-capture/00-prompt.md`.

## Concerns / open questions

- "Tokens" in the owner's words: I read it as provider quota tokens, not a model context window. Confirm in round 2.
- Does the exhaustion class fire on readiness probes at kickoff? I say yes, same evidence, one attempt. Confirm.
- Gateway recognizer scope: which adapters carry the pool text (kimi through omniroute, or others)? The design needs the raw terminal record and its adapter. The `quota-zcode-native-exhaustion-capture` follow-up's private scrubbed-capture retention is the natural precondition for a source-established recognizer later; the owner-bounded fixture is the bridge. State this split in FINAL.
- Successor unusable or exhausted in the same batch: apply nothing, escalate once. Confirm this reading of "below the floor".
- Machine record and the `quota-kickoff-reporting-and-alias-guidance` follow-up: the skip hint touches the same kickoff notice path. Keep that follow-up deferred; name the overlap in FINAL.
- Compaction target: which §9.0 legacy sentences can shrink without rule loss. Needs a budgeted pass like review-gate-timing D5; I do not pre-approve specific deletions here.
- The three driver gaps in `organizer-notes.md` (fixed-slug refusal, the historical facilitator declaration defect, adapter default drift) stay process facts. None blocks this design; record them if they recur.

## Risks

- **False positive.** A healthy agent could drop for a whole idea. Bounded: per-idea only, explicit exhaustion wording, provenance-gated recognizer, fail-closed on bad resets, fresh probe next idea.
- **Gateway wording drift or spoofed text.** The recognizer keys on the provider's own terminal record with recorded fixtures and every-record agreement; prose still never qualifies. An unrecognized text falls back to the two-attempt path.
- **Role concentration after succession.** One agent can hold organizer, implementer and drafter roles. Independence rules must hold: the successor never reviews its own code; reviewer counts and model diversity re-evaluate; the goal check stays with a fresh non-implementer process.
- **Byte budget.** §9.0 grows inside a nearly full always-included block. A careless edit either breaks the 70,000-byte guard or deletes a rule. Compaction must be budgeted and reviewed hunk by hunk.
- **Correlated agreement (§15.6(b)).** All three participants share one brief, one source context and one prior-art family. Unanimity here is a shared prior, not independent proof. The guards are the adversarial fixtures and the refutation-default review, not the agreement.
- **Machine record staleness.** A record past its reset could skip a healthy agent. The record expires at reset plus five minutes; a later kickoff always probes afresh.
- **Organizer succession stays process-level.** A dead organizer process still needs someone to start the successor. The product can record and accept the takeover; it cannot migrate a live process. This limit must be stated in FINAL, not hidden.
