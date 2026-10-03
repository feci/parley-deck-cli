---
agent: codex-1
idea: meta-protocol-change-quota-auto-exclude
round: 1
date: 2026-10-03
---

## Summary
Add a narrow, recorded exception to the quorum lock for positively identified provider quota exhaustion, both before kickoff and at quiescent boundaries during an idea. Keep at least two non-facilitator participants, preserve dissent and review obligations, and stop when exclusion cannot satisfy those constraints; never change the machine roster or automatically replace an implementer.

This is a design proposal, not an exclusion authorization for this run. Supporting mid-idea recovery requires more than a classifier regex and an `excluded:` annotation: participant replay, existing signoffs, and the round-completion gate must agree.

Protocol attestation: `context_mode=full`; `source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`; `packet_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`; `fallback_reason` absent. The shadow audit was not applied. Transport is `github-pr`, with the kickoff's explicit local-files/no-push/no-merge override. Both required source-context documents and the raw readiness JSON were read in full; no other participant's round-01 answer was read. OpenViking tools are unavailable in this launch, so this analysis uses local sources.

## Proposed approach

### Evidence at HEAD

CLI HEAD inspected: `04b22e265cc8bb96975bd59cddaf11b737d4cde9`. The inspected classifier, telemetry, preflight, app and runstate files were byte-compared with `git show HEAD:<path>` and matched. Skill HEAD: `a5664d803f1fb6156ef95ff5921dcf13a3dd2031`. Locators below refer to those commits unless marked machine-local. Source-context incident reports remain attributed testimony about historical executions; replaying their text does not independently establish their cause.

| Brief claim / relevant mechanism | PRIMARY evidence and scoped result |
| --- | --- |
| Exclusion and re-inclusion require confirmation; exclusions are temporary. | `parley-deck/COOPERATION.md:880–889`: “requires **explicit user confirmation**”; “the agent stays in the §2 roster and is re-probed at the next idea”; re-inclusion “also” requires confirmation. The brief's protocol description is CONFIRMED. |
| Quorum locks after Phase 0. | Same file, `745–754`: “a mid-idea unavailability does not silently shrink quorum”; the inactive-agent route requires “> 2 rounds”, a deadline and an inbox ping. Two-participant ideas require both signoffs. CONFIRMED. |
| The floor does not replace reviewer requirements. | Same file, `234–235,262–263`: fast has one reviewer, standard two, deliberation all non-implementers; standard degrades to one with two participants. Phase 8, “Close-decision integrity”, additionally refuses auto-completion with “fewer than two independent reviewers”. CONFIRMED as protocol requirements, not a claim that every implementation path enforces them. |
| The implementer chain has protected assignments. | Same file, `449–451`: the pin precedes designation, global default and fallback; a per-idea designee's unavailability gates, while a global-default designee falls through. “Designations are validated against `participants:` only — nothing reads `excluded:`”. CONFIRMED. |
| Preflight provider failures cannot currently be waived with `--yes`. | `internal/app/preflight_liveness.go:33–36`: “Blocking but never an exclusion”; `internal/app/preflight.go:404–429`: provider/ambiguous gates “are NOT auto-excluded and NOT waivable into an exclusion by --yes”. CONFIRMED. A bare 503 is a separate classifier gap, below. |
| Where `--yes` records exclusions. | `internal/app/preflight.go:147,419–420`: flag text is “confirm excluding unavailable agents (records the exclusion)”; the report appends `"%s — %s — confirmed %s"`. `app.go:1922–1943` passes this list to `runcontrol.Create`; `internal/protocol/workspace.go:187–191` emits `"excluded: " + line`. Standalone preflight returns the report; this path writes the kickoff only during run creation. |
| Recording an exclusion is not participant removal. | `workspace.go:224,232` still uses `strings.Join(participants, ", ")` and `Participants: participants`; `internal/runcontrol/runcontrol.go:70,87` records `opts.Participants` in both the creation event and manifest. New implementation obligation: filter the selected set once and use it for every creation/launch consumer. |
| Continue uses frozen identities, but editing only the manifest will not fix membership. | `internal/app/app.go:1230–1245`: `applyRosterSnapshotToParticipants(run.Participants, … m.RosterSnapshot …)` and `Participants: run.Participants`. `internal/runstate/runstate.go:130–138` takes participants from `run.created`, after applying manifest defaults. The frozen-roster claim is CONFIRMED; replay must also consume the new exclusion record. |
| An already-failed round remains blocked. | `internal/driver/driver.go:464–468`: `case "round.incomplete": return false, nil // authoritative block`. Removing a participant from a slice alone cannot resume this round. |
| Historical signoffs can be retained without requiring absent signers. | `internal/consensus/consensus.go:523–525` already separates who may sign (`known`) from who is awaited (`required`). Lines `553–590` reject unknown signers and still aggregate `StatusBlock` into `TriageBlocked`. Reuse this distinction; do not delete departed signers or their blocks. |
| Machine policy is confirm-breaking. | Machine-local `/Users/tomasfecko/.parley/agents.toml:19`: `roster_change_policy = "confirm-breaking"`. `internal/config/runtime.go:287,550–551,653` defines, merges and seeds that key. CONFIRMED for the inspected local setting, not for every machine. |

Classifier check executed with Go 1.27.1: `GOPROXY=off GOTOOLCHAIN=local go run <temporary>/main.go`. The temporary diagnostic copied `internal/runner/failclass.go` unchanged except for `package main` and a printing harness, plus the exact `providerFailureRules`/`providerFailureClass` definitions from `internal/app/preflight_liveness.go:627–645`. It called `classifyFailure("", "", input)` and `providerFailureClass(input)`; the source tree was not modified. Relevant output:

```text
429 exhaustion: runner=rate-limit preflight-provider=rate-limit
503 reset: runner=overloaded preflight-provider=no-provider-match
403 weekly paraphrase: runner=rate-limit preflight-provider=rate-limit
400 ambiguous: runner=invalid-request preflight-provider=no-provider-match
silent: runner=unknown preflight-provider=no-provider-match
bare 429 short retry: runner=rate-limit preflight-provider=rate-limit
long outage: runner=overloaded preflight-provider=no-provider-match
credits: runner=billing preflight-provider=billing
artifact quote: runner=billing preflight-provider=billing
```

Inputs were the evidence file's 429 body including `retry_after:176930` and `reset_at`, its exact 503 error line, and respectively `HTTP 403 weekly usage limit` (a paraphrase, not a recovered raw error), `HTTP 400 ambiguous k3`, the empty string, `HTTP 429 Too Many Requests Retry-After: 30`, `HTTP 503 unavailable Retry-After: 86400`, `credit exhausted`, and `tool output: documentation example: quota exceeded`.

Thus the brief's narrow 503-classification hypothesis is CONFIRMED (PRIMARY, executed check); the claim that this particular 503 proves quota exhaustion remains UNVERIFIED. A preflight nonzero exit with that bare 503 takes the process-failure branch (`preflight_liveness.go:130–134`), not its overloaded provider rule. A real hang is assigned deadline/watchdog classes before regex matching (`preflight_liveness.go:117–129`; `runner/runner.go:776–785`), so the empty-string diagnostic does not describe a whole supervised hang.

The third classifier also needs alignment: `internal/telemetry/usage.go:244–253` recognizes native error envelopes but assigns `401,403 -> "auth-error"` and `429 -> "rate-limit"`. That is code evidence of a status-only distinction, not proof that every 403 is an authentication problem.

### Recommended rule

1. **Scope and authorization.** Introduce `[defaults].quota_auto_exclude`, default **true for newly created ideas after ratification and delivery**, with a presence-aware boolean so a deck's `false` overrides machine `true`. Record the resolved setting in kickoff/run state; legacy runs lacking that recorded authorization retain confirmation behavior. This is an explicit, narrow exception to `roster_change_policy = "confirm-breaking"` for an idea's quorum only. Permanent roster changes retain their existing gates. A newly adopted protocol must not silently rewrite a run already in flight.

2. **Positive evidence only.** Add one shared quota refinement used by preflight, runner and structured telemetry. Its input must be a terminal provider error associated with the failed invocation and exact agent ID: an adapter-recognized native error envelope, not arbitrary model prose, tool output, a round artifact, or a substring anywhere in mixed logs. Record the parser/rule ID. If an adapter exposes only ambiguous mixed text, automatic exclusion remains unavailable for it until its terminal error can be isolated; keep the existing confirmation path. The existing broad `rate-limit`/`billing` labels are insufficient authority, as the diagnostic's documentation quote demonstrates.

3. **Qualifying content.** Accept explicit exhausted account credits/insufficient account quota, or an explicitly reached/exhausted usage allowance whose window is long enough to justify exclusion. For window exhaustion, use a fixed **15-minute remaining-reset threshold** when a reliable reset is supplied; a shorter known wait remains transient. Without a reset, explicit hourly-or-longer allowance language, such as a reached weekly usage limit, qualifies with `reset_at=unknown`. Generic `quota exceeded`, bare 429, context-window/token-length errors, local loop budgets, auth/payment-setup errors and generic 5xx do not qualify without the required account/allowance semantics. A long `Retry-After` alone cannot turn an outage into quota evidence. Fifteen minutes is a proposed policy choice, not a provider guarantee.

4. **503 decision.** A 503 carrying explicit quota/allowance exhaustion must be recognized before the generic overload rule. The supplied “Unavailable (reset after 5h 51m 11s)” lacks that meaning by itself. Keep it on the human path unless a located provider contract establishes that this exact native error denotes quota exhaustion; a reset duration is not such a contract. Do not add a generic `reset after => quota` regex. Likewise, a disabled gateway's 400 or silent hang never qualifies from secondhand dashboard knowledge.

5. **Batch and floor.** Let `P` be the current distinct participant IDs and `E` all positively classified exhausted participants from the completed dispatch/probe batch. Evaluate `P' = P − E` once. At least **two usable non-facilitator IDs** must remain; subtract the declared facilitator even with `facilitator_participates: true`, and do not count missing/unresolved agents as usable. If the floor fails, apply **none** of the batch, stop and escalate with the whole candidate set. Three non-facilitators losing two stop regardless of arrival order; four losing two may continue if the other gates pass. No sequential “exclude until two are left” algorithm. A facilitator's own failure stops orchestration; exclusion does not appoint a replacement. The existing §1 solo exception remains an explicit human decision, never an exit from this automatic rule.

6. **Mid-idea boundary.** Apply after the current dispatch batch has settled and before any next dispatch, signoff evaluation or close. Do not change quorum while an affected invocation can still write. Preserve complete and partial artifacts, mark a failed invocation as failed, and never synthesize its acceptance. Reconcile the failed round against the new membership: all remaining artifacts must still validate, and a new terminal round evaluation must explicitly cite the exclusion revision instead of ignoring `round.incomplete`. Closed FINAL/IMPLEMENTATION artifacts stay frozen.

7. **Roles and verification.** Exclusion is not reassignment. A live per-idea implementer designation or an implementation pin belonging to a candidate blocks automatic application pending the existing reassignment/waiver/`none` route; preserve partial code and the old pin until a recorded handoff. Before implementation begins, an unavailable global-default implementer may fall through exactly as today. Do not automatically transfer an existing consensus/FINAL draft to another writer; use existing eligible volunteer/fallback rules where applicable, otherwise gate the handoff. Recompute required reviewers and model diversity from the remaining eligible participants. The floor of two does not waive the auto-implementation two-reviewer gate, strict review requirements, an independent goal checker, or `require_model_diversity`. If those cannot be satisfied, stop rather than downgrade the track or silently recruit a new agent.

8. **Dissent survives.** Keep an original/historical participant set for artifact identity and `known` signers, and the reduced set for future dispatch and `required` signers. Existing ✅/🟡/❌ blocks remain append-only; an excluded agent's ❌ or unresolved material `DISPUTED` claim still blocks the affected decision. Exclusion itself cannot resolve it. Carry every open CRITICAL/MAJOR finding into the next review consensus, with an explicit disposition and independent evidence for a fix; strict-gate findings of every severity remain subject to that gate. If no valid resolution can be reached without the absent author, ask the owner rather than count votes. Historical reviews remain evidence, but an unavailable identity cannot stand in for an available reviewer needed for a new round.

9. **Return after reset.** The exclusion lasts for this idea. Store reset time as a scheduling hint, not a re-inclusion timer; no same-idea automatic rejoin and no repeated quota probes. At the next idea, probe the unchanged roster afresh and include ready agents through normal kickoff selection. Same-idea re-inclusion still requires recorded owner confirmation and the existing catch-up rules. This prevents quota resets from changing quorum during signoff.

### Recording, recovery and surfaces

Reuse kickoff `excluded:` records with an honest marker, for example `agent-x — quota-exhausted — automatic 2026-10-03T12:00:00Z — evidence <invocation-id>`. Never write `confirmed` for an automatic decision. Remove the ID from current `participants:` at the same logical transition; a historical marker alone is insufficient.

Keep one bounded structured exclusion record per affected invocation in the existing run manifest/event machinery: original/effective participants, affected IDs, rule ID, UTC observation/reset times, sanitized decisive error excerpt and its source locator, policy setting, phase/round, and application state. Preserve `RosterSnapshot` identities; do not re-resolve models or delete historical entries. Bind a stable transition ID to the run, membership revision and invocation IDs so resume cannot apply it twice or recycle an old error after a later success.

Reuse the existing lock/atomic-write primitives, with orchestration serialized for the **idea**, not merely independent locks for two run IDs. Record a pending transition before changing canonical membership, then mirror the prompt/manifest/event and mark it applied. Before any launch or close, finish a matching interrupted transition idempotently; incompatible participant state or missing evidence stops. Single-file rename does not make these multiple files transactional. All driver/app adapters that captured the old participant list must be rebuilt or refreshed before dispatch. Read-only `status`, `wait` and organizer-brief readers report pending inconsistency rather than repairing it.

Send one deduplicated non-blocking owner inbox notice after a successful exclusion: who, decisive reason, reset if known, survivors, and any remaining gates. A floor/role/integrity failure gets a blocking notice instead. `parley status`, `parley wait` and `parley organizer brief` show automatic exclusions and the effective quorum; `wait` continues for a successful recorded reduction and uses its existing degradation behavior for a new blocking failure. Standalone preflight reports classified candidates and evidence; without a target idea/complete role context it cannot certify that an automatic quorum change is safe. The run kickoff performs that decision for its exact proposed participants.

### Incident consequences under this proposal

| Supplied observation | Result if reproduced as current terminal evidence |
| --- | --- |
| Readiness JSON: Kimi deadline; Zcode `provider-failure:rate-limit` | Neither coarse row alone authorizes exclusion. Capture the decisive terminal error; Kimi stays unresolved. This run also cannot satisfy the proposed non-facilitator floor after reduction; its explicit owner exception and current protocol govern. |
| Zcode 429 “Weekly/Monthly Limit Exhausted”, long reset | Qualifies if captured as native provider failure; exclude only when the batch floor and role guards pass. A pasted report is not a live failure record. |
| Kimi silent probe | No automatic exclusion; silence does not establish exhaustion. |
| Codex 503 “Unavailable (reset after …)” | No automatic exclusion on the supplied evidence. A future proven quota form could qualify, but an affected implementer pin or organizer still prevents automatic handoff. |
| Kimi 400 with disabled-connection explanation | No automatic exclusion. The console error does not establish the dashboard explanation. |
| Kimi 403 weekly usage limit | A corresponding native error explicitly reporting a reached weekly allowance qualifies; the supplied prose summary alone does not. A current implementer pin still requires handoff. |

The supplied Zcode `reset_at` is `2026-10-04T22:14:57.003Z`; do not copy the unzoned “06:14:57” rendering as Europe/Berlin local time. Parse UTC timestamps with `time.Parse(time.RFC3339Nano, …)`, normalize durations from the observation time, and reject contradictory/unparseable reset evidence as automatic authority. Preserve the original string for diagnosis.

### Proposed change footprint and acceptance checks

Protocol hunks: §0 defaults (new narrowly scoped policy); Phase 0 template (resolved policy); §5 and §9.0 (exception, floor, batch behavior, return); Phase 5 (quota does not override designation/pin); Phase 3/6/7–8 (historical signers/findings survive; reviewer/close gates remain); §9 session checklist (notices and effective membership). Mirror the same rule into `parley-deck/COOPERATION.md`, `internal/protocol/defaults/COOPERATION.md`, and the skill worktree's `skills/parley-deck/references/COOPERATION.md`; record the eventual ratified change in `parley-deck/meta/protocol-changelog.md`. This round changes none of them.

Implementation seams, proposed rather than a claimed finished patch:

- Shared quota evidence/refinement: a small `internal/telemetry/quota.go` helper, `internal/telemetry/usage.go`, `internal/runner/failclass.go`, `internal/runner/telemetry.go`, and the result/terminal paths in `internal/runner/runner.go`, `phase58.go` and `acp.go`. Keep generic diagnostic classes; only the typed refinement can authorize exclusion.
- Policy and kickoff: `internal/config/runtime.go`, `internal/app/preflight_liveness.go`, `internal/app/preflight.go`, `internal/app/app.go`, `internal/runcontrol/runcontrol.go`, `internal/protocol/workspace.go`.
- Boundary/replay: a small `internal/driver/quota.go`, the loop/round gates in `internal/driver/loop.go` and `driver.go`, the app's `driver_consensus.go`/`driver_impl.go` adapters, `internal/runmanifest/manifest.go`, `internal/runstate/runstate.go`, and `internal/consensus/consensus.go` for historical versus required membership.
- Visibility: `internal/driver/phasedigest.go`, `internal/app/organizer.go`, and existing status/wait rendering in `app.go`/`wait.go`. Skill changes: `skills/parley-deck/SKILL.md` Startup Flow/Driver-First Operation, plus `references/ROSTER_AND_PROTOCOL.md`; require the CLI evidence path rather than facilitator intuition.

Tests should extend `internal/runner/hardening_test.go`, `internal/telemetry/usage_test.go`, `internal/app/preflight_liveness_test.go`, `preflight_test.go`, `internal/config/runtime_test.go`, `internal/runcontrol/runcontrol_test.go`, `internal/runstate/runstate_test.go`, `internal/runmanifest/manifest_test.go`, `internal/consensus/consensus_test.go`, `internal/app/wait_test.go`, `organizer_test.go`, and add focused `quota_test.go` files beside the new helpers. Concrete acceptance cases:

1. The diagnostic positives/negatives above, native versus quoted/tool-output errors, 403/503 explicit-quota variants, unknown/short/conflicting resets, and truncated/missing evidence. Successful invocations and later-successful retries cannot exclude an agent.
2. `4→2` succeeds, `3→1` changes no membership, duplicate IDs do not inflate counts, a facilitator cannot fill the floor, and every permutation of simultaneous failures gives the same result.
3. A kickoff's prompt, first dispatch, creation event and manifest agree; restart never resurrects an excluded agent from `run.created` or frozen identities. An old `round.incomplete` is only superseded after valid survivor artifacts and the recorded transition.
4. Crash/fault injection at each transition write, replay twice, competing continuations, stale invocation evidence and partial artifacts. No dispatch from mixed membership state and no overwritten participant content.
5. Pinned/designated implementer and partial drafter cases gate; global-default fallback remains available; an excluded dissenter's BLOCK/DISPUTED claim and CRITICAL/MAJOR finding survive. Reviewer/model-diversity and strict/goal-check gates are not weakened.
6. Config `false` overrides `true`, malformed config fails closed, legacy runs stay on confirmation, no same-idea timer rejoin, fresh next-idea probing, and notices/status/wait/brief agree without duplicate notices.

Run affected Go packages and `internal/protocol/drift_test.go`'s `TestEmbeddedDefaultMatchesLiveDeck`; compare the third protocol copy and run the skill repository's relevant installer/structure tests (`test/installer.test.js`, `test/lean-organizer.test.js`). These are future implementation checks, not tests claimed run in this design round. Do not promise a three-file patch: the classifier addition is small, but correct mid-idea resume touches several existing consumers. Leave out quota polling services, billing/gateway APIs, automatic purchases/re-enabling, provider/model switching, permanent roster changes, configurable floors, and a general reassignment engine.

## Existing alternatives

These are enumerated components and existing replacements, not an open-ended search. “Constraint-forced” describes the required behavior; “inherited” describes a replaceable implementation choice.

| ID / hand-built component | Closest shipped mechanism and locator | Recommendation / constraint status |
| --- | --- | --- |
| ALT-1: quota detector | `runner.classifyFailure` (`internal/runner/failclass.go:58`), `app.providerFailureClass` (`preflight_liveness.go:638`), `telemetry.Collector.consume` (`usage.go:229`) | Reuse the native-error parsing and add one shared refinement. A distinct, trustworthy exhaustion predicate is constraint-forced; today's separate broad regex tables are inherited and cannot authorize exclusion as-is. |
| ALT-2: reset parser and cooldown timer | Go `encoding/json`, `time.Parse(time.RFC3339Nano, …)`, `time.ParseDuration`; existing terminal timestamp metadata | Adopt stdlib parsing, no timer subsystem. Evidence normalization is constraint-forced; the fixed 15-minute cutoff is a proposed policy choice. Automatic same-idea return is rejected. |
| ALT-3: quorum reducer / roster editor | `parley preflight --yes`, `preflightReport.Excluded` (`preflight.go:122`), `protocol.CreateIdeaFull` (`workspace.go:178`); permanent alternative `parley roster set <id> --scope machine --state inactive` | Reuse exclusion recording; add actual filtering and replay. Per-idea scope and minimum two non-facilitators are constraint-forced. Reject permanent roster mutation and using `--yes` as a blanket quota waiver. |
| ALT-4: transaction/journal/lock | `store.Store.Append`, `runmanifest.Write/Load`, `driver.acquireLock` (`loop.go:362`), `fsutil.WriteFileAtomic` (`fsutil.go:83`) | Reuse these primitives and existing run records, with an idea-scoped orchestration boundary. Serialization/recoverability are constraint-forced for mid-idea mutation; a new database or daemon is unnecessary. Existing run-local locking alone is not sufficient. |
| ALT-5: departed-signer handling | `consensus.validateDocumentAwaiting(known, required, …)` (`consensus.go:525`) | Adopt the existing distinction. Historical identity and preserved vetoes are constraint-forced; do not invent a second signoff format or rewrite blocks. |
| ALT-6: failure retries / manual pause | Runner retries once only for `no_first_output` (`runner.go:510–515`); §5 async/deadline route and blocking inbox escalation | Keep them for their existing purposes, not quota diagnosis. No retry storm or new sleep-until-reset worker. Waiting preserves membership but cannot meet the requested automatic recovery by itself. |
| ALT-7: preflight-only exclusion | Existing `runTaskPreflight` before `runcontrol.Create` (`app.go:1913–1926`) | This is the smallest alternative, but reject it as the complete solution: the supplied incidents include mid-idea failures. If implementation must be split, label preflight-only as a partial delivery, not fulfillment of the whole proposal. |

## Concerns / open questions

- The decisive design choice is strict provenance over broader coverage. Neither long 429 backoff nor “503 reset after” establishes exhaustion without allowance semantics. Provider-specific proof could change that recommendation; the copied hypothesis cannot.
- Default-on for new ideas and the fixed 15-minute threshold are recommendations for consensus and eventual owner ratification. No human answer is required to produce this independent round. Keeping existing runs unchanged is part of that recommendation, not an assertion that protocol version pinning has already shipped.
- The owner explicitly authorized this two-process/two-participant design run. The stricter proposed non-facilitator floor is prospective; it must not retroactively invalidate or restaff this idea.
- Full mid-idea support has a broader file footprint than the owner's phrase “small change” might suggest. The proposal limits behavioral scope and reuses primitives; it does not hide the resume/signoff work behind a one-line classifier change. If consensus chooses a staged delivery, FINAL should state exactly which incidents remain blocked after stage one.

## Risks

- **False exclusion:** broad regexes can match quoted examples or a long outage. Restrict authority to a failed invocation's recognized native provider error and explicit quota semantics; uncertainty retains the human gate.
- **Dissent erased by a membership filter:** keep historical signers and all findings, and preserve the original BLOCK/DISPUTED obligations. Exclusion changes availability, not the truth of a claim.
- **Resume resurrects an agent or advances a failed round:** reconcile creation events, canonical participants, manifest state and terminal-round evidence under one recorded transition; test interrupted writes and concurrent continuations.
- **Progress is still blocked:** the floor, implementer ownership, reviewer diversity or an unclassifiable provider may legitimately stop a run. This is the cost of a narrow automatic rule; claiming that it solves every supplied incident would be unsupported.
