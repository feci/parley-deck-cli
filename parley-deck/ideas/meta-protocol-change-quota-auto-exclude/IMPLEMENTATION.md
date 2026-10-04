---
idea: meta-protocol-change-quota-auto-exclude
status: implemented
implementer: codex-1
started: 2026-10-04
branch: quota-auto-exclude
head-commit: 78ac536
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

Both stages are implemented at CLI `e8081ab`, with protocol compression at `78ac536`; the current
skill counterpart is `dc85b53`. Stage 1 applies the
owner-authorized bounded zcode stderr recognizer at kickoff. Stage 2 supplies durable membership
transitions, replay and serialization, current-member consumers, retained historical obligations, role
and review gates, and read-only pending-state views. Full independent acceptance is pending.

Host build, vet and changed-file gofmt pass. The first full HOST suite found one regression: the new
quota prose exceeded the existing facilitator packet guardrail by 150 bytes. Shortening only the new
prose fixes it: the unchanged guard now measures 69,966 bytes against 70,000. The full HOST rerun exits 0
with every package passing. Final frozen protocol/drift/packet checks and skill/deck byte comparison pass;
the final skill suite passes 399 Node and 54 Python tests, with all six manifests matching.

The last prose-only edits occurred during the full HOST rerun, while all Go code stayed frozen. Separate
final frozen protocol tests cover those edits. The separate claude-1 full-scope review completed against
CLI `78ac536` and skill `dc85b53`; prior rounds 01/02 were focused provenance feedback only. No merge,
release or installation has occurred.

Role concentration (§15.5): codex-1 organizes and implements; a separate configured claude-1 process
owns every independent review and its signoffs. The owner confirms the attended close after both
review-consensus signoffs and current-tree AC evidence. There is one non-implementer reviewer, so the
unattended two-reviewer gate is not relaxed. Transport is owner-authorized local canonical files on the
two existing quota-auto-exclude branches; no development PRs.

## Implementation plan / checklist

The pre-code plan was committed as `67b954d`; owner handoff as `c6f8b49`. This is its current checklist.

- [x] Stage 1 classifier/reset parsing, scrubbed evidence, owner-defined zcode support, diagnostic-only
  fallback adapters, artifact/later-success precedence and bare-503 provider gates.
- [x] Layered presence-aware policy, frozen scope, complete readiness batch/floor/protected-role checks,
  kickoff filtering before every initial membership consumer, report-only standalone preflight.
- [x] Stage 2 immutable history, one durable batch record, idea lifetime serialization, checked projection
  recovery, pending/integrity gates and no automatic scope widening for old/off/kickoff-only ideas.
- [x] Rebind dispatch/await/signoff/close consumers; historical-known/current-required signers; retain
  vetoes/disputes/findings, partial files, failed invocations and prior round events; validate survivors.
- [x] Protected designee/pin/drafter and review/diversity/goal-check/strict gates; deduplicated notices and
  blocking escalation; same membership/reset/pending views in status/wait/organizer brief.
- [x] Apply normative protocol hunks in all three copies, owner exception, skill guidance, changelog,
  exact skill/deck byte comparison and phase 0/5/8 packet checks.
- [x] Focused adversarial, whole-batch, fault/replay, serialization, history, consumer and race tests;
  host build/vet, changed/new-file formatting and final skill suite.
- [x] Complete full HOST `go test ./... -count=1 -timeout 45m` and preserve the result; final frozen
  protocol/packet/drift checks and full skill suite also pass.
- [ ] Separate full-scope claude-1 review, binding findings, review consensus, and at most five fix-up cycles.
- [ ] Both review-consensus signoffs, independent current-tree AC1–AC21 evidence, blocking attended-close note.
- [ ] Only after owner close confirmation: complete, merge, release channels, independent channel verification.

## Deviations from FINAL.md

- The owner's 2026-10-04 scope/reset answer below explicitly replaces unavailable native root attribution
  for zcode with the bounded stderr rule and adopts the bounded display-clock interpretation. The
  accepted risk is a subagent quota error accompanying an unrelated root failure. All other adapters
  remain diagnostic-only without native evidence. These are owner rulings, not implementer inferences.
- Incident 1 preserves readiness results but no raw stderr; its replay fixture pairs those facts with
  incident 2's preserved provider body. It is never described as a second native capture. The actual
  incident's survivor arithmetic still blocks. See the producer evidence and support table.
- Idea serialization uses OS file locks from existing dependencies, a lifetime lease and short projection
  locks, instead of directly reusing the driver's PID-file primitive named in FINAL §10. This choice and
  crash/recovery behavior require independent review; no extra dependency or off-scope singleton is added.
- The exact deferred record representation uses `quota-history/NNNNNN.json`, hash-bound batch records,
  `quota-applied/<batch-id>` receipts, and optional manifest/digest membership fields. Retained obligations
  have immutable snapshots and structured disposition paragraphs with authority, rationale and independent
  evidence. These structural checks do not establish the truth of the supplied claims. The detailed
  grammar and limitations are in `source-context/codex-1-implementation-evidence.md`.
- Historical AC1 discrepancy resolved: the skill reference now matches the source deck byte-for-byte.
  The CLI embedded bootstrap keeps its existing generic headers/tables and unchanged drift assertions.
  The skill instructs a bootstrapper to replace upstream project header/host mappings under Appendix A.
- The owner transferred organization from claude-1 to codex-1 for Phases 5–8/release; the frozen FINAL's
  old organizer assignment is superseded by the ratification and handoff, not edited in place.

## Notes for reviewers

Review the complete CLI product diff since FINAL `27e42b8`, all follow-up fixes, and the skill diff since
`a5664d8`. The producer map below is not an acceptance verdict. Attempt counterexamples against the
owner's exact zcode rule, all consumers after a transition, incomplete writers, missing/truncated history,
replay/idempotence, protected roles, historical veto/dispute/finding force and every review/close gate.
Also inspect owner-confirmed re-inclusion/withdrawal and scope-change behavior against the frozen rules.
The single-reviewer configuration still requires the attended close. The reviewer may report any issue.

## Progress

- 2026-10-04 00:40Z: owner scope answer recorded in `550fbf8`; answered scope gate archived. Configured
  codex-1 implementation child started (gpt-6-astra/max, 5400-second ceiling). The organizer owned protocol,
  skill and orchestration files; the child owned Go/tests/fixtures and its producer evidence. Commits
  were serialized after it exited. Existing worktrees were reused; no declarations or pruning.
- 2026-10-04 02:02Z: child exited 0 after 4942.3 seconds, without timeout or real provider/auth/quota
  failure. Both stages, focused/affected-package/race checks and build/vet are recorded in
  `source-context/codex-1-implementation-evidence.md`; raw final logs are retained in
  `.parley-runtime/quota-implementation/producer-checks/`.
- 2026-10-04 02:03Z onward: host build/vet, changed/new-file gofmt, exact skill/deck comparison, drift test
  and packets pass; full HOST suite is running. Both-stage protocol wording and skill suite are complete.
- Product commits: CLI `e8081ab`, skill `21f82e7` (preceded by owner amendment/snapshot `844a8b0`).
  Independent full review, consensus, fix-up if needed, owner close and release remain pending.

## Full review outcome — round-03

claude-1's separate process exited 0 after 1116.5 seconds and authored the 31,389-byte review. The
validator reports 1/1 filed-and-valid, unparsed=false. It reports 2 CRITICAL, 4 MAJOR, 5 MINOR and 4 NIT
findings, with no quota/auth failure. The current code is not ready to merge. Shared-volume creation and
locking, manual membership/return/scope paths, stderr framing and signoff handoffs require fixes. The
first review consensus proposes nine grouped fixes; its signoffs precede fix-up cycle 1. Detailed
criterion verdicts and limitations are in the unmodified reviewer artifact. Producer test passes are
not independent acceptance. No owner-close request is made at this defective checkpoint.

## Packet-size correction before full review

The first complete HOST run finished after 679 seconds; all packages except internal/app passed.
`TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` measured 70,150 bytes against its unchanged 70,000-byte
limit. The correction shortens only newly added quota prose identically across all normative copies.
No guardrail, applicability map, omission rule or test is changed. The focused guard passes; final full
Go and skill checks pass and are recorded under `.parley-runtime/quota-implementation/`. This is pre-review
implementation validation, not a Phase-8 fix-up cycle or independent acceptance.

## Decision Log

- Owner scope/reset direction governs the prior MAJOR-1/MAJOR-2 choices. Their reviewer-authored files
  remain unchanged; the next full review must evaluate the actual implementation of the ruling.
- Driver-first resumption prematurely drafted review consensus from focused feedback while implementation
  was in-progress. It was stopped before a consensus artifact existed and its context-canceled note was
  archived. The recorded configured-CLI fallback is used for the focused full-scope review brief; status,
  wait and validators remain in use. Driver gap 11 and legacy accounting migration stay outside scope (D6).
- Normative hunks are identical in all copies; the skill snapshot obeys literal AC1. No drift test,
  bootstrap assertion, test behavior, model/provider/effort setting or roster file was weakened/changed.
- Both stages are required. There is no stage-1-only completion or waiver of independent evidence.

## Surprises & Discoveries

- Child sandbox denied the default Go cache and some host budget-lock helpers. A temporary Go cache
  allowed focused checks; broader denied tests were deferred honestly to the HOST suite, not suppressed.
- An all-file gofmt inventory lists five pre-existing untouched files; all 50 changed/new Go files are
  clean. This is not a claim that the entire pre-existing repository is gofmt-clean.
- Source metadata is still 2.12.0, while installed CLI is 1.50.0 and installer/runtime core skills 2.14.0.
  Dry-run sync is recorded; metadata refresh is authorized release work, not a source-protocol overwrite.
- Root OpenViking recall and the owner-answer update were verified by readback and retrieval. The child
  could not use memory under its sandbox approval policy and relied on local authority. Graphify's current
  graph covers design documents, so code relationships were inspected directly.

## Validation evidence

Producer results only; independent criterion verdicts will be cited from claude-1's own artifact.
Full command logs are under `.parley-runtime/quota-implementation/`; exact child checks and evidence
limits are in `source-context/codex-1-implementation-evidence.md`.

| Criterion | Current producer evidence / remaining gate |
| --- | --- |
| AC1 | Three normative copies, exact skill/deck `cmp`, drift test, changelog, phase 0/5/8 packets and packet-map check pass. Full suite exposed a packet-size failure; shortened new prose makes the unchanged guard pass. Final frozen checks pass; full HOST rerun exits 0. |
| AC2 | Recorded incident-2 fixture, incident-1 labeled replay, every-record/receipt-clock and bounded display-reset tests pass. |
| AC3 | Adversarial class/reset/quoted/mixed/duplicate/truncated/success/watchdog fixtures pass. |
| AC4 | Support table: only bounded zcode is supported; unsupported/adversarial provenance tests pass. |
| AC5 | Whole-batch permutations, 4→2/3→1, duplicates, facilitator/unresolved floor cases pass. |
| AC6 | Kickoff C1 filtering, replay, automatic records and standalone report-only preflight tests pass. |
| AC7 | Shared bare-503 gate/preflight tests pass. |
| AC8 | Designee/pin/started-draft, global-default and stub/protected-batch tests pass. |
| AC9 | Historical veto remains TriageBlocked, current-only append gate, kickoff-never-known and retained-finding tests pass. |
| AC10 | Prospective review/diversity/strict gates and consumer rebinding tests pass; both complete host runs passed that gate coverage. |
| AC11 | Commit/projection fault injection, pending-action refusal, replay/dedup and stale-dispatch tests pass. |
| AC12 | Lifetime/cross-process lease, different-run/off-scope and unsettled-prior-writer tests pass. |
| AC13 | Preserved partial files/failed events and survivor-validation terminal evaluation tests pass. |
| AC14 | Immutable history, later-run discovery, contradiction and consumer-rebind tests pass. |
| AC15 | Status/wait/brief pending/applied agreement, notices and wait exit behavior tests pass. |
| AC16 | Layered/presence-aware defaults, malformed records, legacy/off/frozen scope and upgrade tests pass. |
| AC17 | Roster-free design, immutable scope and reset-hint tests pass; no timer/rejoin path added. |
| AC18 | Provider-evidence-only candidacy and adversarial fixtures pass under the explicit owner exception. |
| AC19 | Both stages present; status is implemented, with full review and attended close pending; never complete here. |
| AC20 | Host build/vet, changed-file formatting, final drift/packet checks and full skill suite pass. First full HOST suite had only the packet-size failure, corrected; full HOST rerun exits 0. See the prose-only timing qualification above. |
| AC21 | Attended close pending both independent review evidence and review-consensus signoffs, then owner answer. |

Protocol attestation read by organizer/implementer at resumption: context_mode=full; source and packet
SHA256 `8e9213bd45059069d484bd10e5ca1a1c509297039dfd8fc67d5e9ebda7590416`; no fallback reason.
The final both-stage packet-check source/packet SHA256 is
`73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e` for phases 0/5/6/8 (full, no fallback).

## User direction — scope/reset answer

The owner relay below is quoted verbatim, including its original Slovak phrases and supplied English
translations, as the controlling brief requires.

## Owner answers to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_scope-reset.md`

Relayed by the owner's Claude Code session on 2026-10-04 at about 02:35 CEST. The relay asked both questions
in Slovak and gave three options for question 1: A (build a reliable zcode channel, codex-1's
recommendation), B (defer, no supported adapter) and a pragmatic option. The selected options are below,
verbatim. A translation follows each one.

**Question 1, AC2 and zcode support.** Selected: **"Pragmaticky: stderr zcode stačí"** ("Pragmatic: zcode's
stderr is enough"). The option read (Slovak, verbatim): "Moje odporúčanie, ide o najrýchlejšiu funkčnú
cestu. Pri zcode stačí JSON `responseBody` v stderr, ak proces skončil chybou, nenechal výstup a každý
záznam o chybe je 429 „Limit Exhausted“ s `reset_at`. Je to výslovná výnimka z FINAL. Riziko: chyba
sub-agenta by mohla vyradiť agenta, ktorý zlyhal z iného dôvodu. Tlmia to minimum 2 a notifikácia."

Translation: "My recommendation, the fastest path that works. For zcode, the `responseBody` JSON in stderr is
enough when the process ended with an error, left no output, and every error record is a 429 'Limit
Exhausted' with `reset_at`. This is an explicit exception to FINAL. Risk: a sub-agent's error could exclude an
agent that failed for a different reason. The minimum of 2 and the notice soften that."

**Question 2, display time without a timezone.** Selected: **"Áno, podľa návrhu claude-1 (Recommended)"**.
This adopts the bounded interpretation exactly as your note states it. A display clock counts only when it
appears together with a machine reset value (`reset_at` or `retry_after`) in the same record and agrees with
it within one second. Missing, contradictory or display-only resets still gate.

## What this authorizes (an owner-directed deviation from FINAL; record it in IMPLEMENTATION.md)

- The zcode adapter becomes a **supported** recognizer, but only under this owner-defined evidence rule. All
  of the following must hold:
  - the zcode process exited non-zero;
  - it produced no valid artifact, and no later attempt in the batch succeeded;
  - stderr contains at least one provider error record whose `responseBody` JSON is a 429 with explicit
    exhaustion text ("Limit Exhausted", or allowance semantics as in FINAL §4.4) and a machine reset value;
  - **every** provider error record in that stderr agrees: each is that same exhaustion class, and their
    reset values agree within the tolerance. A mixed or contradictory record set gates;
  - the run ends with zcode's turn-failure line;
  - the reset clears FINAL's 60-minute threshold.
- Positive fixtures are this run's recorded zcode stderr (evidence incidents 1 and 2). Adversarial fixtures
  include a quoted 429 in assistant or tool text, a mixed 429 plus other-error stderr, a reset under 60
  minutes, a display-clock-only reset and a success after a 429.
- Everything else in FINAL is unchanged: the floor of 2, the role guards, fail-closed, the record and
  notice, both stages and claude-1's binding full review. Other adapters stay diagnostic-only unless they
  have native evidence.
- Record the deviation in `IMPLEMENTATION.md` with this note quoted, and carry the protocol wording into the
  §9.0 hunk under §7. The claude-1 review checks the deviation as implemented, not whether to make it.

## Outcomes & Retrospective

Pending full-scope review, any agreed fixes, both signoffs and attended close. Release preparation is in
`release-plan-codex-1.md`; the owner-only npm OTP and attended protocol publication remain later steps.
No completion, merge, release or channel verification is claimed.
