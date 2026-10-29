---
idea: meta-protocol-change-credit-exhaustion-failover
status: final
author: codex-1
implementer: codex-1
consensus-date: 2026-10-29
participants: [codex-1, kimi-1, zcode-1]
---

## Final plan / specification

All three current participants filed ACCEPT. This FINAL freezes decisions D1–D12.
Codex-1 organizes, drafts, and implements. Kimi-1 and Zcode-1 independently review and goal-check.

### D1 — Freeze a new policy generation

New ideas use trigger `participant-failure-v2`, with typed exhaustion rule `provider-exhaustion.v1`.
Keep the existing presence-aware `quota_auto_exclude` setting and frozen scope.
Saved v1 and omitted-trigger policies retain their behavior, serialized bytes, and hashes.
Only an owner-bound revision may widen an existing policy. A binary upgrade or resume never widens it.
A widening revision freezes the required role/order snapshot without rewriting the old kickoff record.

Reuse the existing reducer, immutable history, terminal records, and step ledger. Add no second membership database.
Ordinary participant failure remains independent of quota recognition.
Its original-plus-one retry, five-second delay, frozen ceiling, and protected-role gates remain unchanged.
Do not introduce automatic retries for code-mutating steps that currently run once.

### D2 — Separate exhaustion from reset data

Exhaustion requires failed execution and explicit account credit, token, quota, or named usage-allowance exhaustion.
Use only an adapter's reviewed native evidence contract. Prose, quotes, tool output, and `is_error` supply no provider authority.
Bare 429/5xx, auth, overload, watchdogs, context limits, and local response/session budgets do not qualify.
Malformed, mixed, contradictory, or truncated class/framing evidence rejects.

Reset data has a separate validity result. Short, absent, past, malformed, or contradictory resets become unknown.
They never cause a retry of otherwise proven exhaustion. The legacy quota predicate and 60-minute threshold remain unchanged.

Validate completed output first. Valid output wins, including BLOCK or FAIL.
An already in-flight later success also wins. Do not launch another attempt to seek success after qualified exhaustion.

### D3 — Stop once and suppress this identity in this idea

After the first qualifying failure, persist the bound invocation and validation result before return.
Start no retry, probe, or reconnection for that identity in this idea, including protected roles.
Common launch admission enforces suppression across steps, prompts, runs, restarts, recovery, policy changes, and blocked membership batches.
Suppression alone does not change membership.

A crash between terminal recording and output validation creates a pending admission gate.
Reconcile the existing output without another child. Missing or corrupt proof blocks for repair.
Old valid artifacts must not mask a failed new implementation or fix-up step.

Scope stays one agent and one idea. Provider/model/gateway labels prove no shared account.
The next idea probes afresh. Defer peer inference and machine caches.
The CLI governs its next supervised child. It cannot undo SDK retries before terminal exhaustion became observable.

### D4 — Apply the owner's succession order to the complete batch

Freeze original participant order and organizer identity at the v2 authority boundary.
Optional `organizer:` names an ordinary participant organizer; otherwise use the declared `facilitator:`.
Conflicting fields block. `author: user` supplies no identity.

The initial implementer may remain unresolved. Before the first role transition, resolve the existing pin/designation/default/FINAL chain.
Bind its result and source in immutable history. Validated role history governs later transitions.
Malformed designations, conflicting pins, or an unknown required counterpart block.
An ordinary participant loss may proceed when no unresolved role affects it.

Apply the binding owner order after all losses settle:

- Organizer loss transfers organization to the implementer.
- Implementer loss transfers implementation to the organizer.
- Loss of both transfers both roles to the first remaining participant in original order.
- One identity with both roles counts as both losses.
- Ordinary participant loss transfers neither role.

The selected counterpart or first participant needs positive usability.
Never skip an unusable selection, wrap, infer an identity, or choose a substitute.
Only proven v2 exhaustion removes protected roles automatically. Protected transient failures retain existing gates.

### D5 — Permit only the frozen facilitator exception

A frozen external organizer's exhaustion permits a role-only transition without membership deletion.
If the frozen declared facilitator takes implementation, append that identity to current membership in the same transition.
It becomes a known and required signer. It supplies no floor seat or independent review.
No other automatic membership addition or new roster identity is permitted.
The declared facilitator remains a permanent floor/reviewer exclusion, including after participation or implementation takeover.

Require two distinct positively usable non-facilitators after the whole batch.
Role labels alone supply no usability. Validate phase artifacts or execution evidence for survivors.
Below the floor, or without an eligible selected successor, apply no membership transition and escalate once.
Exhaustion suppression still holds if that transition blocks.

### D6 — Preserve one active successor execution record

Use one successor-owned `role-handover/<transition-id>.md` per role transition.
This record is the brief's successor-owned phase artifact. It records takeover and replaces the prior record as current execution authority.
Do not rewrite the dropped author's artifact.

Bind transition, predecessor, successor, affected roles, phase, base path/hash, branch, HEAD, and date.
Use ordinary `status:` and `implementer:` keys. Add `## Takeover` and `## Resume state`.
Record the next action and canonical resume paths.
Before Phase 5, `base: none` is valid. This must not imply review-ready implementation.
When implementation starts, the same record must satisfy every ordinary Phase 5 field and required section.
It must pass resolver-aware implementation and fix-up validators.

Later handovers bind the previous active record by exact hash.
Preserve every predecessor byte. Inherit progress, cycle count, validation history, and findings through the base chain.
Never compact or merge old handover files.

One checked `ResolveActiveExecution` resolver serves every current-state reader and writer.
It selects only the pointer from validated typed role history. Filename order and directory listings supply no authority.
Missing, corrupt, contradictory, or changed bases block without fallback.
Route driver, runner, checks, evidence, status, wait, briefs, review, goal checks, fix-up, close, and other execution consumers through it.

Commit typed `role_transition.v1` history before successor dispatch.
Until the successor's handover validates, admit only that selected author's handover action.
Emit one idempotent takeover notice. Do not fabricate an attended `implementer_reassigned` confirmation.
An external organizer process must follow the protocol; this release does not create an external process supervisor.

### D7 — Preserve shared drafts and all filed positions

Complete shared drafts remain unchanged, signed or unsigned. Current participants continue append-only signoffs on that same draft.
Missing drafts use the current eligible selector.

Automatic relocation requires a structurally valid incomplete unsigned draft.
The parser must accept frontmatter and every present structural element. Idea and author identities must match.
Only absent required content or sections may explain incompleteness.
No ACCEPT, BLOCK, reservation, or other signoff may exist.
Malformed identity/frontmatter, truncated constructs, contradictory syntax, or any signoff blocks for integrity repair.

Preserve eligible bytes exactly once at a recorded archive path, without overwrite.
Bind original path, author, hash, and archive before successor authorship at the canonical path.
No old signoff transfers to the new draft. Participant rounds and reviews never move.
Frozen FINAL and closed execution records stay frozen.

Retain every filed dissent, BLOCK, DISPUTED claim, and finding.
A dropped author cannot withdraw them. Use the existing quoted owner-ruling or abandonment/v2 path for retained dissent.

### D8 — Preserve independent review and causal exception proof

A successor implementer never reviews its own code. Its earlier findings remain binding.
Count Before reviewers without the before-implementer and frozen facilitator.
Count After reviewers without the after-implementer and frozen facilitator.
The single-reviewer exception requires at least two Before reviewers and exactly one After reviewer.
After must equal current membership. The sole reviewer must belong to both sets; the successor must belong to After.
Only D5 permits that successor's absence from Before.

Validate typed evidence for every simultaneous loss, including an external role holder.
Keep legacy and paired-failure proof constraints. Do not require the lost set to contain only the before-implementer.
Require known distinct snapshot models for implementer and sole reviewer, even when diversity is disabled.
Re-derive the same proof at precommit, review, review consensus, goal check, and close.

Only the numeric reviewer minimum changes. Keep all other signoff, floor, strict-gate, reservation, dissent, and current-tree evidence duties.
A fresh independent reviewer process performs the goal check. Valid FAIL remains final.

### D9 — Add narrow source-backed Codex JSON support

Accept only explicit effective Codex JSON argv, a captured primary-thread/turn terminal `turn.failed` envelope, and failed execution.
The reviewed source is `rust-v0.165.0` in `source-context/codex-native-source.md`.
Bind the discovered local CLI version to the invocation. Enable this recognizer only for exact version `0.165.0`.
Unknown or different versions remain diagnostic-only, with a support notice. A future version needs its own reviewed contract.
A local version check does not contact the provider and must not create an exhaustion re-probe.

Accept only these exact `error.message` strings:

```text
Your workspace is out of credits. Add credits to continue.
Your workspace is out of credits. Ask your workspace owner to refill in order to continue.
You hit your spend cap set in your workspace. Increase your spend cap to continue.
You hit your spend cap set by the owner of your workspace. Ask an owner to increase your spend cap to continue.
```

Use provenance `codex.exec-json.turn-failed.v1` and label support source-backed.
The event carries text. Never claim runtime access to `RateLimitReachedType` or installed-binary equivalence.
Reject added detail, generic quota text, mixed terminal outcomes, malformed JSON, assistant quotes, and tool output.
Negative fixtures include the source's usage-limit, quota-exceeded, UsageNotIncluded, context-window, rollout-token, and bare turn-failed forms.

Keep default Codex text, Kimi, and Claude diagnostic-only. Do not silently change default argv or persistent adapter configuration.
The opt-in launch contract explicitly selects JSON and must receive ordinary smoke verification.
`source-context/codex-json-smoke-receipt.json` records that ordinary smoke at version 0.165.0.
It arrived after the peer round-3 launches. Both signoff briefs must expose it.
It covers ordinary argv/envelope behavior only. No exhaustion was induced or captured.

### D10 — Keep Zcode grammar and add private diagnostic capture

Keep the existing bounded Zcode framing grammar. V2 changes exhaustion/reset evaluation within that grammar.
Retain its native-positive limitation and existing owner-deviation provenance. No new gateway waiver or grammar expansion exists.

Capture defaults off. Only owner-set local machine configuration may enable it.
Use `[defaults].capture_zcode_failures = true` in the machine configuration only.
Deck/idea settings, prompts, environment variables, and dispatch flags cannot enable capture.
Bind capture state to invocation metadata and expose it in operator status.
Never send capture state to a provider.

Capture only an ordinary failed Zcode invocation's stderr after termination and validation.
A valid completed artifact suppresses capture. Never probe, retry, or contact a provider for capture.
Use `~/.parley/captures/zcode/`, outside worktrees. Directories use 0700; files use 0600 and exclusive creation.
Refuse symlinks and unsafe existing permissions. Scrub before persistence; uncertain scrubbing writes nothing.
Remove credentials, authorization, cookies, prompts, request payloads, and unrelated response payloads.
Keep only bounded provider-error framing and review fields. Mark removals with `[REDACTED]` and caps with `[TRUNCATED]`.
Keep the terminal tail when a cap applies. Do not silently alter non-secret framing.

Use a 256 KiB per-file cap, a 2 MiB total cap, and seven-day retention.
On each use, delete expired files first, then oldest files to fit the total cap.
Record cleanup counts and time only, without content, in a bounded private cleanup receipt.
Provide `parley quota capture status` and `parley quota capture purge` for the owner.
Directory deletion also provides immediate cleanup. Disabling stops new captures without deleting retained evidence.
No timed deletion is promised while the facility is idle.

No recognizer or membership reducer reads the capture store in this release.
A capture supplies material for later review, not native-evidence authority.
Truncation or redaction of required framing prevents a complete-native-evidence claim.
The open Zcode follow-up retains its provenance and completeness precondition.

### D11 — Apply corrected exact protocol hunks within fixed limits

Use `source-context/consensus-protocol-hunks-proposed.json` as the sole signed hunk set.
SHA256: `f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`.
The corrected set replaces the round-3 set. It restores the non-solo advisory-lens protection and independent participant-failure basis.
It routes all current execution updates through the active record and retains the full implementation schema.
The readable full patch and delta from round 3 accompany it.

Apply exactly the signed hunk bytes to the live CLI protocol, bundled CLI default, and skill reference.
Build staged core 2.20.0 from staged core 2.19.0 plus exactly those hunks.
Preserve each file's project/header zone. Do not publish global core.

Keep the 20,000-byte SKILL.md and 70,000-byte packet limits.
Add the actual phase-1 meta packet as a hard test: deliberation, github-pr, facilitator, `protocol_change`.
Use the existing test's source display fields and fail unexpected fallback. Retain the unflagged case.
The separate same-request design measurement records 68,442 unflagged and 69,874 meta bytes.
These measurements do not certify semantic preservation.

Disclose two compaction changes: remove the historical 40-deck rationale sentence; correct the §10 label from escalation to summary.
Restore the omitted non-solo qualifier instead of retaining that trim.
Compact required template descriptions without changing their fields, headings, conditional duties, or gates.

### D12 — Deliver the full authorized release

Release CLI 1.56.0 and skill/core 2.20.0 after independent review, fix-up, and fresh goal check.
Use the allocated CLI and skill worktrees. Follow the controlling brief's retry limits and maximum five fix-up cycles.
Merge design and implementation through merge commits. Publish GitHub releases and both Homebrew formulae.
Submit a skill-only WinGet PR. CLI Windows stays experimental and CLI WinGet stays held.

Install all managed skill runtimes and the four named generic destinations.
Run the requested copy inventory and verify SKILL.md hashes. Preserve historical sources and backups.
A separate Kimi or Zcode process verifies release channels.

npm publication and `parley protocol publish` remain owner-only. Do not execute or bypass those gates.
The relay owns the owner-publication script. The released inbox note lists pending owner commands in version order 2.15.0 through 2.20.0.
Include channel evidence, deferred scope, and usage. Native fallback receipts are not driver accounting or membership authority.


## Purpose / user-visible outcome

A reviewed native account-exhaustion failure ends that agent's work in the current idea after one failed attempt.
The CLI does not reconnect it. The owner's succession order selects the next role holder when the floor and gates permit.
Transient failures retain their shipped behavior. Existing ideas retain their saved policy semantics.

## Context & orientation

Work in `/Users/tomasfecko/AI_WORKSPACE/parley-deck/worktrees/credit-failover` for the CLI.
Use its sibling `credit-failover-skill` for the skill package. Both allocated branches are `credit-failover`.
Do not prune or declare worktrees. Never change credentials or gateway settings.

The current product bases are CLI 1.55.0 and skill/core 2.19.0.
The shipped quota reducer, immutable kickoff/history, terminal telemetry, and participant-step ledger supply the extension points.
Preserve legacy serialization and use optional fields for new typed authority.
Resolve role state before dispatch. Commit authority before its projections and handover dispatch.
The active execution resolver must not call a reader that recursively invokes it.
Route the Phase 5 pin read and generated Phase 0 template comment through D6.
The listed legacy chain describes initial fallback; validated handover history takes priority.

Implementation areas and files:

- `internal/quota/quota.go`, `record.go`, `history.go`, `revision.go`, and `dropout.go`; new exhaustion and succession schema/validation files.
- `internal/telemetry/quota.go`, `quota_zcode.go`, `quota_framing.go`, `record.go`, and `usage.go`; new source-backed exhaustion recognizer and capture helpers.
- `internal/runner/dropout.go`, `telemetry.go`, `runner.go`, `phase58.go`, `quota.go`, `action_identity.go`, `launch_precheck.go`, and `consult.go`; new handover dispatch helpers.
- `internal/membership/membership.go`, `gates.go`, `review_gate.go`, revision and notice paths; typed whole-batch role settlement and recovery.
- `internal/protocol/quota.go`, `implementer.go`, `facilitator.go`, workspace and retained-obligation readers; new checked active execution resolver.
- `internal/app/app.go`, `preflight.go`, `driver_impl.go`, `driver_consensus.go`, `quota_consumers.go`, and `quota_revision.go`; new capture status/purge command.
- `internal/app/driver_checks.go`, `driver_evidence.go`, `evidence_table.go`, `evidence_verify.go`, `wait.go`, `organizer.go`, `learn.go`, and `pipeline_cmd.go`.
- `internal/driver/cursor.go`, `impl.go`, `phasedigest.go`, and action dispatch; `internal/consensus/consensus.go`; `internal/evidence/report_guard.go`.
- `internal/tui/protosnap.go`, `internal/retro/retro.go`, `internal/runcontrol/runcontrol.go`, `internal/runmanifest/manifest.go`, and `internal/config/runtime.go`.
- The live and bundled CLI `COOPERATION.md` copies, `internal/app/facilitator_packet_live_test.go`, CLI version/release metadata, and targeted tests beside changed packages.
- Skill `skills/parley-deck/SKILL.md`, its `references/COOPERATION.md`, `ROSTER_AND_PROTOCOL.md`, headless invocation reference, artifact templates, and targeted tests.
- Skill package/lock, installer/manifest compatibility metadata, release metadata, and generated payload manifests when their source changes.

Audit all remaining fixed `IMPLEMENTATION.md` references. Keep literal initial/history references only when their purpose requires them.
The exact protocol hunk set is D11's hash-bound JSON, not an ad hoc edit list.
Record any necessary file-list adjustment in the living execution plan. It must serve these decisions and acceptance criteria.

## Observable acceptance criteria

| ID | Observable result | Required evidence |
| --- | --- | --- |
| AC-01 | Actual new-idea paths freeze v2. Saved v1 and omitted-trigger records retain behavior, bytes, and hashes. Only owner-bound revision widens policy. | Legacy goldens; new run/readiness creation tests; owner-revision tests. |
| AC-02 | Only reviewed native framing and exact contracts yield exhaustion. Reset uncertainty never negates valid exhaustion. Valid completed output or an in-flight later success wins. | Positive/reset/adversarial fixtures for Zcode and exact-version explicit-JSON Codex; diagnostic fallback fixtures. |
| AC-03 | A qualified first failure launches exactly one child. No retry delay, probe, or reconnect follows. Suppression survives changed steps, prompts, runs, restart, policy downgrade, and blocked batches. | Child-count lifecycle tests through common launch admission and replay. |
| AC-04 | Terminal-before-validation recovery admits no new child and validates existing output. Old valid implementation/fix-up artifacts cannot excuse a new failed step. | Fault injection around terminal and validation writes; stale-artifact tests. |
| AC-05 | Transient participant failures retain two total attempts, five-second spacing, and the original ceiling. Protected/code-mutating paths retain their prior retry gates. | Existing tests plus v2 transient parity, control-refusal, cancellation, and protected-step cases. |
| AC-06 | Every owner-order case selects the required identity after the whole loss set settles. Unknown or unusable counterparts block. The selector never skips or substitutes. | Table-driven ordinary, organizer, implementer, both, same-holder, simultaneous-loss, and original-order cases. |
| AC-07 | The frozen external facilitator alone may enter membership through implementation takeover. It supplies no floor/review seat. External organizer loss can be role-only. | Membership/history tests, positive usability proof, fixed-floor cases, and invalid-addition rejection. |
| AC-08 | Each transition has one selected handover. Predecessor bytes and base hashes remain fixed. Current readers/writers follow only validated history; missing/corrupt authority blocks. | All-phase handover-chain tests, filename-order adversaries, base corruption, pending-action admission, and durable-write fault injection. |
| AC-09 | Complete drafts stay. Only valid incomplete unsigned drafts move byte-exactly once. Signed, malformed, contradictory, or truncated drafts block. | Separate consensus/FINAL/review-consensus recovery cases; unchanged round/review bytes. |
| AC-10 | Succession preserves every filed obligation. A successor never reviews its code. The single-reviewer exception proves the complete before/after cause and distinct snapshot models. | Retained BLOCK/finding tests; reviewer-to-implementer cases; simultaneous/external losses; fresh goal-check eligibility. |
| AC-11 | Codex support is inert without explicit JSON and exact discovered version 0.165.0. Exact messages pass; source near-misses and model/tool quotes fail. | Effective argv/version binding tests, source fixtures, support output, and the ordinary smoke receipt. |
| AC-12 | Capture is off unless local machine owner configuration enables it. Only failed validated Zcode steps create private scrubbed bounded files. Status, purge, markers, quotas, and cleanup obey D10. | Default-off and source-precedence tests; permissions/symlink tests; scrub/tail/size/age/cleanup cases; no capture-store authority. |
| AC-13 | The three protocol copies receive exactly the signed hunks. Staged core uses staged 2.19.0 plus those hunks. Both phase-1 packets stay within 70,000 bytes; SKILL.md stays within 20,000 bytes. | Exact hunk/copy checks, unchanged header zones, hard meta and unflagged tests, package size tests. |
| AC-14 | CLI 1.56.0 and skill/core 2.20.0 explain adapter limits, saved-policy behavior, role handover, capture controls, and account-alias scope. | Help/reference/template review, package/version checks, full Go and skill test suites. |
| AC-15 | Independent review closes with no open agreed fixes, a fresh independent goal PASS, and current-tree AC evidence. Authorized channels and runtime installs match released bytes. | Kimi/Zcode artifacts, review consensus, goal receipt, merge/release hashes, Homebrew and skill-only WinGet evidence, runtime hashes, and separate channel review. |

Use the prompt's full check command: `go test ./... -count=1 -timeout 45m`.
Run focused package tests during implementation. Add targeted race checks for changed concurrent admission or settlement paths.
Run the skill package's `npm test`, package checks, and native portable build.
Use synthetic failure processes for exhaustion tests. Never induce account exhaustion or make a provider call for evidence.

## Idempotence & recovery

Immutable kickoff, owner revisions, role transitions, terminal attempts, and validation receipts are authority.
Mutable prompt/run/status projections are recoverable views. Never reconstruct authority from excluded markers, filenames, or elapsed time.
A committed transition with incomplete projections blocks ordinary work until checked recovery finishes.
A selected pending handover permits only that successor's handover action.
A durable exhaustion result suppresses further admission even when a membership batch blocks or policy later changes.

Reconcile existing output after terminal/validation interruption. Do not buy another child attempt through a new run or prompt.
Retain predecessor execution files and unsigned-partial archive bytes with exact hashes.
Retry publication or notices only through existing identity-bound receipts. Preserve owner-edited or archived notices.
Do not overwrite a release, signature, frozen FINAL, closed execution record, or existing staged core with different bytes.

Before release, verify reviewed source commits and a clean build source. Bind assets to those commits and record hosted hashes.
Verify hosted Windows skill hashes after the release workflow finishes, because it can replace uploaded portable assets.
Install the reviewed release payload, then check runtime hashes. Preserve source archives and backups.
The owner alone performs npm and global-core publication. Report pending commands instead of invoking those gates.

## Known risks / de-risking

All three proposals extend one existing reducer/ledger architecture. Nominally separate proposals are one family.
Model agreement does not independently prove source correctness, complete call-site coverage, or recovery safety.
Try to break native framing, restart admission, handover authority, draft integrity, and review-cause proof through fixtures and independent review.

Source-backed Codex support does not prove binary/source equivalence or captured native exhaustion.
Zcode retains its native-positive limitation. Unsupported adapters may still take their ordinary retry path.
The CLI cannot prevent SDK-internal retries before it observes terminal exhaustion.
Exact version binding may disable Codex support after an upgrade until a new contract is reviewed.

The capture facility is diagnostic only. Lost required framing prevents later completeness claims.
Cleanup occurs on use, not through an idle background service.
Historical execution chains cost one file per succession; preserve them without compaction.

## DISPUTED claims

The peers differ on whether the old 16-hunk proposal lost a rule or only repeated wording.
The consensus preserves their exact statements and provenance treatment.
The corrected delivery restores the disputed clauses and replaces the stale active-record text.
No decision or acceptance criterion depends on the old claim that all 16 hunks preserved every rule.
No vote or organizer verdict resolves that historical interpretation.

## References

- Consensus: `./consensus.md`; signoffs are the phase gate.
- Independent rounds: `./round-01/`, `./round-02/`, and `./round-03/`.
- Controlling brief: `./source-context/ORGANIZER-BRIEF.md`.
- Owner succession and fallback answers: the two named `user-to-codex-1` inbox records, quoted in the prompt/round 2.
- Exact corrected hunks, full diff, correction delta, and measurement: `./source-context/consensus-*`.
- Native source and ordinary smoke: `./source-context/codex-native-source.md`, `./source-context/codex-json-smoke-receipt.json`.
- Prior quota/dropout/review-gate/driver-unstall FINALs and the two scoped quota follow-ups listed in the prompt.

## Protocol context

context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none
