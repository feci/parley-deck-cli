---
idea: meta-protocol-change-participant-dropout
status: in-progress
implementer: codex-1
started: 2026-10-08
branch: participant-dropout
head-commit: 431d6b0748ef7c16722c09899d1cc0e609e5eccb
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
- [ ] Focused tests (AC1–AC12): paired failures/success, timeout/watchdog, structurally valid BLOCK and malformed own output, control-plane/tamper stops, restart after first/second failure, concurrency and no third launch, floor/protected/gate permutations, permanent return after opt-out/downgrade, legacy round-trip bytes, kickoff replay and surfaces. Preserve existing legacy quota tests/expectations; explicit legacy fixtures may use omitted trigger while new default expectations change.
- [ ] Current-tree checks (AC13): affected-package tests during implementation; full `go test ./... -count=1 -timeout 45m`, `go vet ./...`, `go build ./cmd/parley`, formatting/diff checks; skill full tests and generated-manifest check; protocol drift and phase 0/5/8 full/facilitator packets with unchanged guards. Record exact tree/commands/outcomes and material limits.
- [ ] Independent review: focused zcode brief with full implementation diff in both worktrees, FINAL, ACs and validator outputs, open inspection/no suppression. Reviewer owns every finding and verdict. Sign review consensus before fixes; maximum five fix-up cycles and stopping judgment. A fresh independent goal-done check cannot establish close by itself.
- [ ] Attended close: both final review-consensus signoffs, no open CRITICAL/MAJOR, current-tree independent AC evidence and brief's standing authority. No product two-reviewer-gate waiver or automatic close claim.
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

No implementation AC is yet claimed met. Design validators: `parley wait --for round` reported 2/2 valid for round 2; `consensus request-signoffs` validated Zcode's own append; `consensus status` ready; `consensus finalize --by codex-1` accepted FINAL. Product tests have not run against a new implementation.

## Outcomes & Retrospective

Pending implementation, independent review, close and release.

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
