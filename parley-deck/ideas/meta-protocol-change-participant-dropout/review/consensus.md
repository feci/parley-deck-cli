---
idea: meta-protocol-change-participant-dropout
review-cycle: 1
outstanding_agreed_fixes: 3
blocked: false
drafted-by: codex-1
date: 2026-10-09
reviewed-commit: e4681cf5b9144ed786b269d8ab14b8d65b23d581
---

## Scope and phase

This is the Phase-7 fix plan, not final close. The controlling brief and live Phase 7→8 require both signoffs BEFORE fixes; review/round-01/zcode-1.md's concluding suggestion to wait for fixes is handled through that normal sequence. Codex-1 organizes, participates, drafts and implements under the owner brief; Zcode-1 alone owns the independent code-review verdict. Skill commit efe296c7acf13a147ab820ce6cbf8e6705b68691 remains in scope.

## Agreed fixes

1. **Z1 / MAJOR — durable goal-check validation.** Replace the volatile res.Answer validator with a replay-capable validator bound to the attempt's retained logs. Rehydrate the actual answer on replay across changed run IDs, preserve FAIL as valid dissent, and never turn a refused/cancelled/integrity outcome into a completion pass. Missing or changed previously-valid evidence stops rather than authorizing another child. Audit readiness's similar volatile observation closure and require its replay validation to read retained per-invocation probe output too. Test the terminal-before-validation-receipt crash window for valid PASS/FAIL/PONG and malformed output, using actual child logs, plus replay with a different run ID and no extra child.
2. **Z3 / MINOR — do not dispatch past valid output on a control-class replay.** Recognize the receipt/validated output before the generic control-record skip. Preserve control/integrity refusal as a blocking error (not dropout and not a completion pass), and do not launch a replacement over valid output. Retain repaired pre-dispatch refusal behavior where no valid artifact exists. Tests cover cancelled valid output and integrity refusal separately.
3. **Producer CI observation — distinguish ACP exit from deliberately short timeout.** The new started-exit subcase shares a 150ms ceiling intended for the timeout subcase; PR Linux CI observed its exact class assertion fail while the push Linux run and 30 focused local repetitions passed. Give the started-exit case a realistic ceiling while preserving the separate timeout/watchdog cases, two-attempt cap, and equal original/retry ceilings. Print actual class values on assertion failure. This fixes a fragile new fixture without changing runtime failure classification or weakening expectations.

## Disputed finding / proposed disposition for Z2

Zcode's PRIMARY probe measured four children across two complete `checkRoster` calls with separately generated kickoff IDs. Codex accepts that observation and disputes the conclusion that they are the same idea/batch. FINAL D2 expressly distinguishes the boundary:

> “Kickoff readiness uses the same two-attempt rule within the proposed batch, before the first authoritative participant list or manifest. A standalone preflight reports observations and never applies exclusions. Undispatched readiness/setup uncertainty is not fabricated child evidence. After an idea exists the restart bound is per idea; a new idea probes everyone afresh.”

The same language is in the signed design consensus. Each new `parley run` proposes a new idea; no idea or resumable run exists before preflight succeeds. This probe completed one proposed batch and invoked a second, so its two distinct IDs exercise next-proposal probing, not resume of an existing idea. Mid-idea resume uses the saved slug and cannot replenish the ledger. A genuinely same readiness batch with a supplied stable ProbeID already consumes the same ledger; add an explicit replay test, including the durable PONG receipt repair in item 1, so that contract is demonstrated.

**Proposed disposition:** withdraw Z2 as an implementation defect under FINAL's explicit pre-idea/idea distinction; retain this operational boundary visibly in IMPLEMENTATION and release notes. A product feature that resumes an uncreated proposal across separate `parley run` commands would need a durable proposal identity plus explicit abandon/new-proposal semantics to distinguish “restart” from “new idea”; FINAL did not specify that lifecycle. Such a feature is a separate follow-up (TBD), not a silent relaxation of the existing per-idea cap. This is a rebuttal for the reviewer to judge, not an implementer-issued withdrawal. Zcode should BLOCK with its counter-proposal if the quoted boundary does not resolve its concern; no finding is suppressed.

## Deferred follow-ups

- **Z4 / MINOR — headroom warning:** measured phase-1 body 69,963/70,000 B and SKILL.md 19,995/20,000 B pass unchanged caps. Accept the present bounded size as a known maintenance limitation; the next additive protocol change must budget compaction (follow-up TBD). No limit is raised and no current AC is waived.
- A persistent pre-idea proposal/resumption lifecycle, if desired, is TBD under the Z2 proposed disposition. No shared membership authority or global roster change is introduced here.

## Dismissed findings

None is unilaterally dismissed. Z2's proposed disposition requires Zcode's explicit concurrence or a new review round. Z1/Z3 remain open until implemented and independently verified. Z4 is a documented maintenance risk, not a present cap breach.

## Validation prerequisites and updated observations

Full source log references are in source-context/validation-review01.md. The review accidentally names `review-persistence-failure-recovery`; the actual failed case is `report-persistence-failure-recovery`. Do not edit the reviewer's artifact; it can self-correct in its next artifact.

- Skill efe296c full `npm test` rerun now passed 399 Node tests, 54 Python tests and all manifests; prior failure retained. macOS portable build/version/install/doctor smoke passed, as did all 19 installation dry runs.
- Current CLI vet/build passed; an empty redirected vet log is expected on success. The completed process exit was zero. No standalone empty file proves completion.
- The first full host run failed one recovery case; isolated three-repeat and full 19-case diagnostic-overlay three-repeat checks passed. Cause remains unestablished.
- The second full host run later failed because the internal disk filled: Go linker and TempDir errors explicitly report “no space left on device”. This is not passing evidence. Fresh tests will use task-local TMPDIR/GOTMPDIR/GOCACHE on the workspace volume, which has about 385 GiB free; filesystem permission enforcement was probed successfully. No user files or shared caches are deleted. External-storage build passes.
- Fix-up cycle 1 has not started. Fresh full host command, focused regressions, vet/build, protocol checks, independent current-tree AC1–AC12 evidence and a zero-fix final review consensus are mandatory before attended close. AC14 delivery is after close; npm/core remain owner-only.

## Coverage & blind spots

One independent reviewer remains under the brief's explicit exclusion/close authority. Zcode ran the seven changed CLI packages and scratch probes, inspected both complete diffs and all normative hunks. It did not run all 37 packages itself, and AC10 was spot-checked; the fresh verification brief must cover these limits. Real hosted error timing, cross-host/PID namespaces and native Windows remain outside demonstrated coverage. The pending storage-limited run cannot replace required evidence.

## Drafter position changes

Codex has no review-round artifact because it is the implementer. Since the candidate's IMPLEMENTATION snapshot, codex accepts Z1's durable-validator defect and Z3's replay gap, and adds the ACP fixture repair based on actual CI evidence. It disputes Z2's same-batch inference using the quoted frozen FINAL boundary; no FINAL text is changed. All procedural dispositions remain provisional until signoff.

## Signoffs

<!-- Each participant appends only its own signoff. -->

### Signoff: codex-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: I accept the three fix items and validation prerequisites as the implementer, with Z2 subject to the explicit reviewer concurrence on the frozen FINAL boundary. This is a fix-plan signoff, not code self-review or close. No fixes before both ACCEPT blocks.
