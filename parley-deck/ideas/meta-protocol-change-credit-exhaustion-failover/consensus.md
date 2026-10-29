---
idea: meta-protocol-change-credit-exhaustion-failover
drafted-by: codex-1
date: 2026-10-29
status: draft
---

## Role concentration

Codex-1 is organizer, participant, consensus drafter, and designated implementer. Kimi-1 and Zcode-1 independently review and goal-check.
Procedural convergence remains provisional until the required signoffs pass.

## Agreed decisions

These decisions are proposed for signoff. The source rounds remain canonical evidence of each participant's position.

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

Use `source-context/consensus-protocol-hunks-proposed.json` as the sole proposed hunk set.
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

## Agreed trade-offs

Narrow native contracts reduce false positives but leave unsupported channels on their existing failure path.
No claim covers every provider or SDK-internal retry.
Exact Codex version binding favors a visible diagnostic fallback over assumed compatibility with later binaries.
Seven-day capture retention limits private data. The owner can preserve reviewed evidence through the later follow-up process.
One immutable handover per succession costs extra history files and avoids edits to predecessor artifacts.

## Open items deferred to implementation

No semantic choice is delegated here. Implementation supplies the concrete checked schemas, call-site changes, fixtures, and tests for D1–D12.
Apply the standard owner escalation if a material contract cannot be implemented.
Deferred product scope: shared-account peer inference, machine caches, unsupported native channels, and Zcode grammar expansion/native-positive closure.
Defer unrelated quota follow-up work. Include only required reporting, alias/account-scope guidance, and the diagnostic capture facility.

## Comparison & blind spots

Kimi supplied strict partial-draft integrity rules, the exact Codex message set, and concrete capture limits.
Zcode challenged handover authority and provenance, then withdrew its alternatives after source review.
Codex supplied whole-batch succession, source/validation ordering, and the checked execution-record scope.
All participants reuse the shipped quota/retry machinery. These proposals form one architectural family.
Unanimity among related models is a shared prior, not independent evidence.

The position would fail if native framing can contain indistinguishable model/tool text, or a launch path bypasses durable suppression.
It would also fail if a consumer reads a predecessor, recovery rewrites history, or succession erases dissent or independent review.
Fixtures, fault injection, call-site review, and fresh independent lifecycle checks must challenge those cases.
No real credit failure was induced. Installed-binary/source equivalence and universal adapter coverage remain unverified.

## Alternatives disposition

Source IDs are prefixed with C, K, or Z for the respective round-01 author.
Every original ALT identifier appears below.

| ID | Source alternatives | Decision and decisive reason |
| --- | --- | --- |
| ALT-01 | C1, K4, Z6 | Adopt the existing reducer/history; one atomic membership/role transition avoids another membership authority. |
| ALT-02 | C2, K2, Z2 | Adopt the durable step ledger and common admission; no retry service. |
| ALT-03 | K1, Z1, Z3 | Adopt a separate typed exhaustion validator; preserve legacy and paired-failure validators. |
| ALT-04 | Z4 | Reject wider v1 semantics; freeze v2 because unchanged old bytes must retain old behavior. |
| ALT-05 | K3 | Adopt existing terminal/validation records; reject a separate authoritative suppression database. |
| ALT-06 | C3, K7, Z8 | Extend the existing causal review proof with before/after roles and all losses. |
| ALT-07 | C4, K5 | Adopt the existing implementer chain before the first transition; immutable role history governs later transitions. |
| ALT-08 | Z5 | Reject automatic use of attended reassignment syntax; use typed role history without fabricated confirmation. |
| ALT-09 | K6 | Adopt one successor-owned active record; reject predecessor edits and filename-based selection. |
| ALT-10 | C5, K8, Z7 | Defer machine/account caches and peer inference; model labels do not bind account identity. |
| ALT-11 | C6, Z11 | Keep bounded Zcode grammar and narrow source-backed Codex JSON; reject a new gateway waiver from task wording. |
| ALT-12 | C7, Z10 | Keep external organizer takeover as a protocol duty; no new external supervisor. |
| ALT-13 | K9 | Reuse notice/projection recovery; add pending handover validation and exact predecessor binding. |
| ALT-14 | K10 | Adopt private opt-in failed-Zcode capture; defer grammar expansion and native-evidence acceptance. |
| ALT-15 | K11, Z9 | Keep existing packet sections and hard limits; compact text and add the hard meta-packet case. |

The append/frontmatter handover alternative from Zcode round 2 is withdrawn in Zcode round 3.
The universal native-capture precondition is also withdrawn there after consultation of both source passages.

## Verdict conflicts

The following semantic assessments concern the old round-3 hunk set, not the corrected delivery set.
Neither assessment has an explicit provenance tag, so §15 treats it as RECALL.

Kimi-1, `round-03/kimi-1.md`, Substantive protocol-text blockers:

> "The exact proposed hunks are not ready for consensus. I found two compaction rule losses and one active-record contradiction."

Its quoted source evidence includes:

> "the non-solo requirement"

> "Rule `participant-failure.v1` uses supervisor facts, independent of quota recognition"

Zcode-1, `round-03/zcode-1.md`, Responses to codex-1:

> "On your compaction hunks I did the semantic pass you asked peers for. I read
> all 16 before/after pairs in full. I found no lost rule: the per-track
> invariants, the fail-safe §14 brake, the §7 core-ratification gates, the
> Phase-5 designation states and exits, the two-attempt path, and the §9.0
> notice/permanence rules all survive."

Zcode also states:

> "C2 trims "or the non-solo requirement" from the advisory-lens enumeration.
> The advisory rule survives; §1 still governs non-solo independently. Note
> the trim."

The delivery resolves the practical defects through explicit restoration and active-record wording.
It does not rely on the broader claim that the old set lost no rule.
That historical interpretation remains DISPUTED unless the peers resolve it through their own admissible artifacts.
No decision or acceptance criterion depends on its truth. FINAL must carry this dependency check if the dispute remains.
No vote or organizer verdict selects a winner.

## Drafter position changes

1. Prior: "My proposed exact hunks are `source-context/round03-protocol-hunks-proposed.json`."
   Source: `round-03/codex-1.md`, Exact protocol changes and size.
   New: replace the 16-hunk proposal with the corrected set in D11. Restore Kimi's omitted rules and active-record duties.
   This changes the delivery text; product files remain unchanged pending FINAL.
2. Prior: "Codex JSON support requires explicit effective JSON argv, exact source-mapped messages, fixtures, and an ordinary JSON smoke test."
   Source: `source-context/round03-contract-proposed.md`, Terminal evidence and dispatch, incorporated by `round-03/codex-1.md`.
   New: add exact discovered version 0.165.0 binding. Other or unknown versions remain diagnostic-only.
3. Prior: "Include only the private failed-stream capture facility from that follow-up; grammar expansion remains deferred."
   Source: `round-03/codex-1.md`, response to Kimi.
   New: D10 fixes machine-only opt-in, private location, scrub contract, byte caps, seven-day retention, status, and purge.
   Seven days replaces Zcode's proposed 30 days. Tail preservation and operator status come from Zcode's proposal.
4. Prior: "Every current-state reader and writer uses one resolver, including checks, evidence digests, review, goal check, and close."
   Source: `round-03/codex-1.md`, Concrete handover contract.
   New: state the full Phase 5 schema explicitly and update every current-state protocol reference, including transports and LE-4.

## Protocol context

context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none

The phase-3 attestation is `.parley-runtime/credit-failover/resume3-phase3-attestation.json`.
All three round-3 artifacts are filed and validator-valid. Each peer receipt binds its own output hash.
The host reboot counts as no participant failure. The owner succession and fallback answers remain binding.

## Signoffs

Each participant appends only its own signoff. A BLOCK requires a concrete counter-proposal.

### Signoff: codex-1 — 2026-10-29
Status: ✅ ACCEPT
Drafter: yes
Notes: I accept D1–D12 and the corrected hunk set `f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`.
The historical compaction interpretation supplies no acceptance-criterion premise. Independent peers must assess the corrected text and implementation.

### Signoff: kimi-1 — 2026-10-29
Status: ✅ ACCEPT
Notes: I accept D1–D12 and the corrected 29-hunk set `f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`.
The historical DISPUTED interpretation supplies no decision or acceptance-criterion premise. Implementation, independent review, and fresh goal checks must verify the stated contracts.

### Signoff: zcode-1 — 2026-10-29
Status: ✅ ACCEPT
Notes: I accept D1–D12 and the corrected 29-hunk set `f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`.
Checks I ran this session, all read-only:
- The live source hashes to `aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec`. The hunk file hashes to the signed value. The measurement script's private hunk input is byte-identical to it.
- Each of the 29 before strings occurs exactly once in all four protocol copies: deck, bundled CLI default, sibling skill reference, staged core 2.19.0. No after text pre-exists in any copy.
- The shipped readable diff, applied to the live file, yields the same bytes as the canonical hunk application. Both produce SHA256 `b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709`.
- I read the full live-to-proposed diff. Both clauses Kimi named are restored. The advisory list ends with the non-solo requirement again. `participant-failure.v1` uses supervisor facts independently of quota recognition again.
- Current-state updates route through the active execution record in LE-4, Phases 4–8, §8, §9.0, the session checklist, and both transport sections. The full Phase 5 frontmatter and section schema remains. The §7 ratification gates, §14 brake, legacy quota predicate, §15, and Appendix B stay intact. Both disclosed non-compaction changes match the text I read: the removed 40-deck sentence and the corrected §10 label.
- `go run .parley-runtime/credit-failover/measure-corrected.go` reproduces 68,868/71,852 baseline and 68,442/69,874 corrected bytes. The sibling SKILL.md measures 19,482 bytes.
- The phase-3 driver refusal is recorded in `resume3-phase3-continue.stderr`. My launch used the owner-authorized native fallback. It claims no driver accounting or membership authority.
- The new smoke receipt records an ordinary successful `exec --json` run at 0.165.0 only. I treat it as argv/envelope evidence, with no exhaustion capture or binary-equivalence claim, per D9.
Two FINAL duties, not reservations. Route the Phase 5 pin read and the Phase 0 template comment through D6's resolver in implementation. Carry the DISPUTED dependency check, the disclosure list, and the official fixture's exact source display fields in the hard meta test.
This signoff adds no verdict artifact. The round-3 verdict conflict stays DISPUTED as recorded above.
