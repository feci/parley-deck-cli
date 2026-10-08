---
idea: meta-protocol-change-participant-dropout
status: complete
implementer: codex-1
started: 2026-10-08
branch: participant-dropout
head-commit: cb78e9f43e70cb59284a21ba6e0e8a331cb447a1
design-pr: https://github.com/feci/parley-deck-cli/pull/75
implementation-pr: https://github.com/feci/parley-deck-cli/pull/76
---

## Summary of work

The participant-failure implementation is now present in both owner worktrees. Signed design D1–D7 / AC1–AC14 was published by design PR #75 at merge 431d6b0; the original plan was committed at 1ae039b before product edits. codex-1 organizes, drafts and implements under the controlling brief; zcode-1 independently reviews and never delegates acceptance to the implementer. There is no declared pure facilitator.

Full phase-5 protocol context has source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e; fallback_reason absent. FINAL remains frozen. CLI/skill worktrees and release bases are those named in FINAL.

## Implementation plan / checklist

- [x] Compatible policy/evidence foundation (AC1, AC4, AC6–AC8): add exact optional saved trigger and immutable paired-attempt evidence; retain old JSON/hash bytes when absent; derive permanent dropped IDs from kickoff/batches; validate all revision/manual/catch-up paths against that history. Never weaken the existing reducer or gate semantics.
- [x] Bounded execution (AC1–AC3, AC6, AC9): introduce the smallest shared step helper around existing invocation telemetry and phase validators. Freeze logical identity independently of run ID or mutable prompt; serialize per-step admissions under existing membership/lease controls; record original/retry IDs and consume two total. Separate actual child failure from control-plane refusal/cancellation and unresolved writer recovery. Archive own invalid output without touching peers; valid dissent always wins. Exercise exec and ACP.
- [x] Dispatch integration (AC4, AC5, AC9, AC10): select new default only for new ideas; wire readiness batches, round/review and signoff/single-step paths to the paired evidence. Add only phase-valid evidence for undispatched usable seats, require designated/pinned implementer usability, preserve precommit CheckGates and current/known consumers.
- [x] Narrow kickoff reporting (AC5, AC11): safe missing-inbox creation, detailed stderr fallback on blocked publication, and kickoff notice publication receipt/replay. Preserve owner-edited/archived notice copies; no exactly-once delivery claim.
- [x] Protocol/skill/release metadata (AC12): apply P1–P6 to the three normative copies, short matching skill guidance, protocol changelog and usage. Version CLI 1.52.0 and skill 2.16.0 including required generated metadata. Do not import Windows, alias, D6, model or roster changes.
- [x] Focused tests (AC1–AC12): paired failures/success, timeout/watchdog, structurally valid BLOCK and malformed own output, control-plane/tamper stops, restart after first/second failure, concurrency and no third launch, floor/protected/gate permutations, permanent return after opt-out/downgrade, legacy round-trip bytes, kickoff replay and surfaces. Preserve existing legacy quota tests/expectations; explicit legacy fixtures may use omitted trigger while new default expectations change.
- [x] Current-tree checks (AC13): affected-package tests during implementation; full `go test ./... -count=1 -timeout 45m`, `go vet ./...`, `go build ./cmd/parley`, formatting/diff checks; skill full tests and generated-manifest check; protocol drift and phase 0/5/8 full/facilitator packets with unchanged guards. Record exact tree/commands/outcomes and material limits.
- [x] Independent review: focused zcode brief with full implementation diff in both worktrees, FINAL, ACs and validator outputs, open inspection/no suppression. Reviewer owns every finding and verdict. Sign review consensus before fixes; maximum five fix-up cycles and stopping judgment. A fresh independent goal-done check cannot establish close by itself.
- [x] Attended close: both final review-consensus signoffs, no open CRITICAL/MAJOR, current-tree independent AC evidence and brief's standing authority. No product two-reviewer-gate waiver or automatic close claim.
- [ ] Delivery (AC14): merge implementations, release GitHub assets, both Homebrew formulae, skill-only WinGet PR, all managed and four generic installations with hashes; stage core from staged 2.15.0 plus exact reviewed hunks; separate zcode channel verification. Finish canonical released note with exact owner npm/core commands (2.15 first if still unpublished), usage and limitations.

## Deviations from FINAL.md

None. Internal helper/record names will follow the narrowest shared execution seam, as FINAL expressly permits. Changes in scope or behavior require the existing consensus/owner path; implementation difficulty is not permission to omit an AC.

## Notes for reviewers

Priority refutation targets: hidden third attempts after restart; false eligibility from control-plane failures or quoted text; invalid-output preservation; historical hash changes; no-return bypass after policy off/downgrade or kickoff catch-up; role-only floor inflation; protected drafters; precommit reviewer/model gates; loss of kickoff notice across crash. Read raw artifacts and exact source, not this checklist as proof.

The practical precommit reviewer limit is an accepted design choice, not an unresolved finding: FINAL D3 preserves it under brief item 4. Please assess implementation against that exact rule. Narrow kickoff reporting repairs are newly promised; the previous release's waivers do not establish their pass. Native Windows remains outside demonstrated runtime support and CLI WinGet remains held.

## Progress

- 2026-10-08 21:44Z — Design merged as 431d6b0, both canonical ACCEPT blocks retained; implementation plan written (completed: signed design and plan gate; remaining: all product work, verification, review and delivery).

## Decision Log

- Keep the named participant-dropout worktree/branch for implementation after fast-forwarding to the merged design. Rationale: owner's explicit workspace instruction; avoids new worktree allocation or pruning. 2026-10-08 · codex-1.
- Use the existing saved quota policy/versioned trigger and reducer, not a second membership authority. Rationale: signed FINAL D1. 2026-10-08 · codex-1.

## Surprises & Discoveries

- Driver D6 legacy cycle accounting also refused the first typed manual round-02 launch before child start. Recorded configured-CLI fallback produced Zcode's independent round; the signoff/finalize CLI worked normally. No accounting migration or fake participant failure.

## Validation evidence

Historical pre-implementation snapshot (superseded by the implementation checks below): no implementation AC was then claimed met. Design validators: `parley wait --for round` reported 2/2 valid for round 2; `consensus request-signoffs` validated Zcode's own append; `consensus status` ready; `consensus finalize --by codex-1` accepted FINAL. Product tests had not yet run at that snapshot.

## Outcomes & Retrospective

Implementation is present; independent review, completed current-tree checks, attended close and release remain outstanding.

## Implementation snapshot — 2026-10-09

The saved optional trigger preserves omitted-trigger JSON bytes. The existing reducer/history,
revision/return guards, projections and receipts remain authoritative. New-trigger evidence binds
both supervisor observations and stable step identity; legacy recognizers remain unchanged.
The runner's small participant-step ledger uses existing immutable invocation records, a PID lease,
private structural receipts and preserved output. It never fabricates terminal telemetry during
stopped-writer recovery. Private probe logs now retain the readiness attempts too.

New ideas select participant-failure-v1. Kickoff, rounds/reviews and dispatched signoffs use the
same original-plus-one retry, including watchdog/ACP paths; protected drafts/implementers stop
through the existing gate. Valid own output (including dissent) wins; shared-signoff tampering
stops. Invalid own suffixes are privately copied before restoring their unchanged shared prefix;
interrupted recovery consults the same receipts. Read-only captured-evidence verification retains
its existing separate accounting contract; untargeted consultations are not canonical outputs.

The precommit reviewer gate remains unchanged. Tests exercise both a permitted reduction and
an auto_implement reduction refused with one independent reviewer. Permanent-return tests cover
owner revision, downgrade, manual edit, catch-up and recovery, plus historical BLOCK/findings.
Kickoff publication now shares applied receipts; blocking notices safely create inboxes and print
the decision on delivery failure. Replayed paired failures deduplicate blocking decisions across
new run IDs.

P1–P6 are applied to all three protocol copies. The first wording exceeded the unchanged
70,000-byte facilitator cap; it was compacted, not waived. Current packet tests pass. The skill's
20,000-byte core cap is also retained (current SKILL.md: 19,995 bytes). Versions are staged as
CLI 1.52.0 / skill 2.16.0; no release, install or core publication has occurred.

## Validation record (implementation checks, not independent review)

- Focused policy/membership/runner/app/telemetry/consensus checks passed. Actual local child tests
  cover two attempts, five-second kickoff delay, exec/ACP failures, timeout/watchdog sharing,
  restart/changed-run cap, concurrent admission, private partials, valid artifacts and control
  failures. Invalid-own signoffs, shared tampering and valid BLOCK each have separate fixtures.
- Owner/manual/catch-up/recovery permanent-return tests and kickoff receipt replay passed.
- Packet/drift tests pass after compacting the protocol; no applicability map or limit changed.
- `go vet ./...` and build passed on the preceding implementation tree; current-tree checks follow.
- The first full Go run was started before final integration changes and is not final acceptance.
  It is recorded as a superseded check; a fresh full host command runs on the committed candidate.
- The skill test setup first lacked the locked commonmark development dependency. `npm ci
  --ignore-scripts` restores locked dependencies; no dependency version or audit fix is introduced.
  Earlier manifest failures coincided with protocol edits during that run; final manifests are
  regenerated before the new full skill run. These attempts are not claimed as passes.

Private validator logs live in .parley-runtime/participant-dropout-launches/. Canonical summaries
will record completed commands and the reviewed commits, not private raw provider transcripts.
Independent Zcode review, final host/skill checks, fresh criterion evidence, attended close and
all delivery steps remain outstanding. The existing D6/Windows/alias limitations are not repaired.

## Review-01 candidate and validation update

CLI e4681cf5b9144ed786b269d8ab14b8d65b23d581; skill efe296c7acf13a147ab820ce6cbf8e6705b68691. Both PRs (#76 / skill #9) are ready for review. Separate measured `parley agents exec` process runs Zcode round 01 with the complete diffs, frozen FINAL, all ACs and open-ended refutation brief. The known driver phase-pointer gap is retained; no D6 migration or model/roster adjustment.

- Full host command completed exit 1: 33 package passes, three packages without tests, one failed app package. Only `TestEvidenceVerifierProductionClosure/report-persistence-failure-recovery` failed; its recovery launch was refused. The same subtest passed three focused repetitions; a diagnostic-only private overlay repeated the entire 19-case test three times and passed (156.779s), without source changes. Cause remains unestablished. Fresh full host rerun is pending, not a pass claim.
- Full skill `npm test` now passes on efe296c: all 399 Node tests, 54 Python tests and all six manifests. The prior 398/399 run's initial fleet-install failure is retained, with its focused 1/1 pass and full rerun evidence. No dependency version or source change between these attempts.
- Hosted Linux push CI passed; PR CI failed the new ACP started-exit test's exact failure-class assertion with its 150ms ceiling. Thirty focused repetitions passed locally. The hypothesis of timeout-sensitive test setup and proposed realistic started-exit ceiling are disclosed to review; no fix has been applied ahead of signed review consensus.
- Current candidate `go vet ./...` and build pass. Packet/drift and focused behavior checks passed. Native Windows remains experimental and is not claimed supported.

Raw implementer observations and exact log names: `source-context/validation-review01.md`. These are not independent review verdicts. Release preparation is `release-plan-codex-1.md`; no channel publication, global install, main merge or core staging has occurred.

## Fix-up cycle 1

status: complete
completed: 2026-10-09
head-commit: cb78e9f43e70cb59284a21ba6e0e8a331cb447a1

### Fixes applied

Both ACCEPT blocks were committed at 85ea372 before product fixes. This cycle implements all three agreed items in review/consensus.md; it is not a close claim.

- Z1: process-output validators now receive immutable invocation identity and recover retained output. Goal-check restores PASS/FAIL from its original run's stdout even when terminal telemetry survived but its structural receipt or private copy did not. Readiness similarly reconstructs PONG/classification from private per-invocation logs. Receipts bind raw output hashes; missing or changed previously-valid evidence blocks without a replacement. A nonzero/failed child cannot establish goal completion merely by writing PASS.
- Z3: replay validates retained output before skipping a control-class terminal. Valid cancelled output and integrity errors remain blocking, with no replacement/dropout. Repaired pre-dispatch refusals without valid output still allow their first real child.
- ACP fixture: started-exit now has a five-second ceiling; intentional timeout remains 150ms and watchdog two seconds. Both original/retry share the ceiling and existing cap; assertion failures print actual classes.

### Operational boundary and dispositions

Zcode explicitly withdrew Z2 in its own cycle-1 signoff after rereading frozen FINAL D2: before an idea exists the cap is within the proposed readiness batch. A stable ProbeID replay uses the same two attempts. A separate new proposal probes afresh. After creation the identity is per idea and run/input changes never replenish it. Tests cover same-batch PONG and malformed replay; CHANGELOG.md and docs/quota-membership.md visibly retain this boundary. A durable pre-idea proposal/resume/abandon lifecycle is a separate TBD follow-up. FINAL and the cap are unchanged.

Z4 remains an accepted maintenance limitation: phase-1 packet 69,963/70,000 B and skill core 19,995/20,000 B. Next additive change needs budgeted compaction (TBD). No limit or criterion waiver.

### Validation evidence

- `go test ./internal/runner ./internal/app -run 'Dropout|GoalCheck' -count=1 -timeout 10m` passes on cb78e9f's source (runner 5.191s, app 34.857s), recorded in cycle1-regressions-ready.log. Real child regression cases cover PASS/FAIL/malformed crash replay across run IDs, the pre-preservation window, changed/missing valid evidence, unsuccessful PASS output, same-ProbeID readiness, cancelled valid output with/without a receipt, preserved integrity refusals and repaired pre-dispatch refusal.
- Two intermediate regression runs exposed errors in newly written test fixtures (missing idea driving lease and an assertion about a refusal occurring before the step ledger); fixtures were corrected before the passing run. No product behavior was weakened to satisfy them.
- Early post-compaction tests still encountered internal disk-full at the host budget-lock cache. The disk subsequently had 6 GiB available without any agent cleanup. Large Go scratch/cache stays on the task-local external paths in test-storage.json. No shared cache deletion or budget migration.
- Fresh full host test/vet/build and Zcode round-02 are next. Full skill tests on unchanged efe296c remain passing. Earlier failing/invalidated full-host runs remain recorded in source-context/validation-review01.md.

### Deviations from agreed fixes

None. Source hashes and record-aware validation extend the existing small participant-step receipt; no new membership authority, retry cap, policy knob, protocol hunk or release version change.

## Cycle-1 completed native validation

At candidate f8f4f1f (product cb78e9f), fresh full host `go test ./... -json -count=1 -timeout 45m` completed exit0 in731.058s:34 packages pass, three packages without tests,3262 passing test/subtest events. Vet/build and focused regressions pass. Full skill tests/manifests remain passing on unchanged efe296c. Exact logs, hashes, environment-sensitive preceding failures and local-scratch refutation are in source-context/validation-cycle1.md. This completes the implementer's required checks; Zcode owns independent criterion evidence and final verdict. No attended close or release yet.

## Attended implementation close — 2026-10-09

Under the controlling brief's standing owner authority and frozen FINAL D7, codex-1 records the attended implementation close. The final review consensus lists zero agreed fixes and BOTH current participants independently appended ACCEPT; `parley consensus status --review` reports ready. No CRITICAL/MAJOR remains. Exactly one implementation fix-up cycle was needed (maximum five). Codex did not review its own code.

Zcode's full-scope round-02 review independently evidences AC1–AC12 and verifies AC13. Its separate fresh goal-check-zcode-1.md confirms current product identity and AC1–AC13 after renewed live-tree dropout/goal-check/readiness tests, build/vet and full-host log recount. The qualifying full-host evidence is explicitly full-host-tests-cycle1-local-result.json/jsonl (exit0 atf8f4f1f,731.058s,34 passing packages/3no-test,3262 passing test/subtest events; four built-in skips). The similarly named non-local result is FAILED environment history, not a pass. Full skill tests/manifests pass on unchanged efe296c.

Product head remains cb78e9f43e70cb59284a21ba6e0e8a331cb447a1. Later commits contain only review, evidence and close artifacts; no source change invalidates the independent checks. Z1/Z3 are fixed and independently verified, Z2 reviewer-withdrawn with its commitments landed, Z4 maintenance headroom and the recorded TBD follow-ups retained. The reviewer repaired its own exact heading; artifact validation exits0. Prior signed fix consensus is preserved byte-identically.

This is an attended close under the brief, not driver auto-close or a product reviewer-count waiver. Native Windows, real hosted error timing and cross-host/PID limits remain disclosed. AC14 is expressly post-close delivery and is tracked by release-plan-codex-1.md and the final released inbox note; this closed implementation artifact does not claim the channels already published. npm/core remain owner-only commands. The user task continues through all authorized release channels and separate Zcode channel verification.
