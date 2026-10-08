---
idea: meta-protocol-change-participant-dropout
review-cycle: 2
outstanding_agreed_fixes: 0
blocked: false
drafted-by: codex-1
date: 2026-10-09
reviewed-commit: f8f4f1f17bc99832aeb46c0a8ee44c6331325d1b
---

## Scope and attended-close authority

This is the final zero-fix review consensus after fix-up cycle 1 of the five-cycle maximum. Codex-1 organizes and implements under source-context/ORGANIZER-BRIEF.md and frozen FINAL D7; Zcode-1 is the independent reviewer. Kimi's two real failed attempts and owner-authorized exclusion remain recorded; no solo exception or product gate waiver is invented. Both current participants must append their own ACCEPT block before the pre-authorized attended close.

Reviewed CLI product commit cb78e9f43e70cb59284a21ba6e0e8a331cb447a1 is unchanged at candidate f8f4f1f17bc99832aeb46c0a8ee44c6331325d1b and subsequent deck-only validation commits. Skill remains efe296c7acf13a147ab820ce6cbf8e6705b68691. The signed cycle-1 fix plan is retained byte-identically in review/consensus-cycle-01.md (SHA256 4db5a60edc68c89a3194328200e8c973363a7bb53e74d9618a8071c88a7a07b3); both signoffs were committed before implementation. FINAL is unchanged.

## Agreed fixes

None. Independent review/round-02/zcode-1.md requests zero agreed fixes after full-scope inspection, current-tree refutation and four additional adversarial probes. No CRITICAL or MAJOR remains open.

- Z1, formerly MAJOR: fixed by durable invocation-aware goal-check/readiness validation and retained-output digest checks. Reviewer independently verified PASS/FAIL/PONG/malformed crash replay, before-preservation fallback, changed-run cap, failed PASS not completing, and changed/missing/symlinked evidence refusals.
- Z3, formerly MINOR: fixed by validity-before-control-skip replay. Cancelled valid output and integrity failures remain blocking without replacement or dropout; repaired pre-dispatch refusal without valid output can dispatch its first child.
- ACP fixture: agreed five-second started-exit ceiling is implemented; timeout/watchdog cases and original-plus-one cap unchanged. Three independent repetitions pass.

## Deferred follow-ups

- Z4 / MINOR: protocol/skill headroom remains an accepted maintenance limitation; phase-1 body 69,963/70,000 bytes and SKILL.md 19,995/20,000 bytes. Next additive change must budget compaction (TBD). No cap or AC waived.
- Persistent pre-idea proposal/resume/abandon lifecycle: TBD. Stable-ProbeID same-batch replay holds the two-attempt cap; separate new proposals probe afresh; created ideas retain a durable per-idea cap. The distinction is visible in IMPLEMENTATION, CHANGELOG and docs/quota-membership.md.
- Reviewer observation: inherited fixtures relying on internal-disk flock/git semantics or unquoted scratch paths may merit separate test hygiene (TBD). They are not findings against this product diff.
- Existing D6 accounting, native-positive legacy quota recognition, aliased-deck plain-edit guidance and Windows portability follow-ups remain separate. AC14 delivery follows attended close; npm/core publication stays owner-only.

## Dismissed or withdrawn findings

Z2 was withdrawn by Zcode in its own cycle-1 signoff under the frozen FINAL D2 pre-idea/created-idea boundary and reaffirmed in review round 02 after verifying all conditional commitments. It is not an implementer dismissal. Z1/Z3 are fixed and independently verified, not dismissed. No other finding is silently suppressed.

## Current-tree criterion evidence and exact full-host result

Independent evidence for AC1–AC12 is review/round-02/zcode-1.md, including expanded AC10 read-only surface tests, seven changed-package runs, ACP repetitions, and four additional crash/evidence probes. Zcode also independently ran vet/build and the full unchanged skill suite: 399 Node tests,54 Python tests and all six manifests.

The qualifying native full-host evidence is explicitly **full-host-tests-cycle1-local-result.json** and **full-host-tests-cycle1-local.jsonl**, under .parley-runtime/participant-dropout-launches/: exit0 at f8f4f1f,731.058s,34 passing packages plus3 packages without tests,3262 passing test/subtest events. The raw log SHA256 and commands are in source-context/validation-cycle1.md. Zcode independently recounted the log and verified the result. This completes AC13 along with passing focused regressions, vet/build and skill checks.

**full-host-tests-cycle1-result.json** and **full-host-tests-cycle1.jsonl** record the preceding FAILED external-scratch run (exit1), retained as environment-failure history and never claimed as a pass. Earlier review-01 failures also remain retained. Local TMPDIR/GOTMPDIR plus external GOCACHE yielded the qualifying complete pass without source edits or selective test exclusions (four built-in skips remain reported). No shared cache deletion or accounting migration.

Separate measured Zcode invocation eddf724d-2365-4cd3-a09f-e7494525026f completed exit0 in232.593s and authored goal-check-zcode-1.md (SHA256 e5c3aeb364d2494ffcfdde94ad1d9c875319f140996cb9d197e379c247bbacdf). It independently verified current product identity, recounted the qualifying full-host log and reran live-tree dropout/goal-check/readiness checks plus build/vet. Its explicit implementation AC1–AC13 PASS remains conditional on both final signoffs. The reviewer corrected its own exact Refutation attempts heading, and parley wait now reports its round-02 artifact valid and exits0. This is attended assessment of the canonical goal-check artifact, not a claim that a driver auto-close parser ran. Textual PASS alone cannot establish close.

## Coverage & blind spots

One independent model-diverse reviewer remains under the brief's explicit authority. Zcode reviewed both complete product diffs and used its own archived product-identical tree to keep live files untouched. Real provider 429/503/auth timing, cross-host/PID namespaces and native Windows are not demonstrated by local supervisor fixtures. Windows remains experimental with CLI WinGet held. Existing review/diversity gates are unchanged in the product; this run's attended authority does not become an automatic override.

## Drafter position changes

Codex accepts Zcode's cycle-1 verification and all dispositions above. There are no product changes after cb78e9f and no fresh design change. The reviewer-owned exact-heading repair is a structural artifact correction, not a suppressed finding or another implementation cycle. Zero fixes is the reviewer's independent conclusion; codex signs implementation completion and the evidence/authority chain, not its own code review.

## Signoffs

<!-- Each current participant appends only its own block. -->

### Signoff: codex-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: I accept the zero-fix final consensus as implementer and organizer under the controlling brief. Product cb78e9f is unchanged; Zcode independently verified AC1–AC13 and completed the fresh goal check. The qualifying full-host result is explicitly the -local exit0 pair; failed runs remain failed history. Z1/Z3 fixed, Z2 withdrawn by reviewer with commitments met, Z4 and stated TBD follow-ups retained. This is implementation/evidence/authority acceptance, not self-review. Attended close follows the other participant own ACCEPT block; AC14 delivery then proceeds.

### Signoff: zcode-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: I accept as the independent reviewer on the basis of my own current-tree work in review/round-02/zcode-1.md (full product diff, four adversarial probes, focused suites, skill suite) and PRIMARY re-verification at signoff time: consensus-cycle-01.md and goal-check-zcode-1.md SHA256s match this file's claims, full-host-tests-cycle1-local-result.json is exit 0 at f8f4f1f while the exit-1 run stays retained as environment-failure history — satisfying my round-02 open question 1 — and HEAD 9083c54 is deck-only over f8f4f1f with product unchanged since cb78e9f. Zero agreed fixes is my own conclusion; Z1/Z3 fixed and verified, Z2 withdrawn by me and the disposition stands, Z4 deferral and TBD follow-ups retained. Remaining disclosed limits (local supervisor fixtures, cross-host/PID, native Windows) are FINAL-recorded, not waivers.
