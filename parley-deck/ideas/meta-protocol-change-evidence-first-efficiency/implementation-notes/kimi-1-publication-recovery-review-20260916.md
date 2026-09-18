---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
reviewed-source: /private/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-recovery-publication-pmwo7ia9/source
review-kind: supporting-source-review (not Phase-6 acceptance; no signoff; no full-audit claim)
---

# kimi-1 supporting source review — recover-verifier-parent correction

Review of Codex's new bits only: the `recover-verifier-parent` dispatch
(`internal/app/trajectory.go:23-25`, usage `:53`), `runTrajectoryVerifierParentRecovery`
and the corrected relaunch diagnostic (`trajectory_verifier_recovery.go:315-349`, `:243`),
the new actual-process regression (`trajectory_verifier_publication_test.go`), and one
corrected expected error (`trajectory_verifier_recovery_test.go:556`). My own
attended-recovery layer is not self-verdicted; the pre-existing library
preview/publish mechanics are cited as located evidence, not re-approved by me. I
executed nothing: runtime outcomes are SECONDARY (coordinator execution), supported
by preserved records I read directly.

## The disagreement — recovery gap, not ergonomics

Codex (codex-1-recovered-parent-publication-20260916.md:10-19): a successful
replacement whose parent-recovered publication fails is stranded at the CLI —
invocation consumed, relaunch refused, the retained failed parent refused by
ordinary recover-parent — while the library preview/publish has no CLI route. My
candidate note treated that route as future ergonomics (as codex's note :17-18
characterizes it; my original note stays unchanged). On located evidence I concur:

- Route absence is mechanical fact: the before-control log shows the request
  failing exit 2 with a usage list lacking recover-verifier-parent
  (before.stdout.log lines 4-9; before-result.json `"exit": 1`, 15.658s). PRIMARY.
- The consumed replacement cannot relaunch: the test drives the real replacement
  to completion, then asserts duplicate relaunch is refused (publication_test.go:46-49).
  PRIMARY (test source; outcome SECONDARY).
- Ordinary recover-parent cannot reinterpret the retained failure:
  `compatibleUnpublishedParent` (parent_recovery.go:48-54) admits only
  `FailureStage == ""` or `"parent-publication"`; the retained `"launch"` failure
  is rejected at `:113-121`. PRIMARY.
- Under the design's own recovery contract (original retained byte-for-byte,
  invocation consumed, no refund/retry), publishing the retained observation is
  completion of the recovery, not cosmetic ergonomics. The factual clauses carry
  the tags above; the contract classification is my position change.

## Refutation attempts against the correction (LE-1)

R1 wrong/missing lineage — fails. `newRecoveredParentPreview` requires
`effective.Recovered` (parent_recovery.go:503-505) and pins the original to the
exact pre-start budget refusal (`:517-524`: RefusedInvocationID,
RefusedTerminalSHA256, empty receipt, nil assessment, `"launch"`, pending);
retained-record replay revalidates the same lineage live (`:457-488`).
R2 stale digest — fails. Apply re-derives the preview from live state and refuses
on mismatch (`:619-622`); replay refuses `expected != retained.SHA256` (`:601-603`);
`validHash` gates format (`:582-584`); the test exercises a well-formed wrong
digest (publication_test.go:76-81).
R3 partial publication — fails. Publish is stage-fsync-rename with an O_EXCL
random tmp and deferred cleanup (`:551-574`); readers only Lstat
`parent-recovered.json`, so a crash leaves at worst an ignored tmp; the directory
obstruction fails the rename without clobbering.
R4 replay after later state change — fails, with a note. Replay returns the
immutable retained preview after live lineage revalidation; `StateSHA256` is
publication-pinned (format-checked only, `:459`); post-reconciliation replay is a
validated fsync-only no-op (`:604-608`).
R5 scope/config change after replacement — considered, rejected. The route
launches nothing; its binding is ticket/recovery/invocation lineage, not the live
roster; re-running scope checks would let a later config change strand
already-produced evidence. Omission is correct.
R6 accidental work/charge — fails. The handler calls only
`PreviewRecoveredParent`/`PublishRecoveredParent` (`:331-335`); no runner, budget
store or telemetry launch. The test asserts ledger and verifier-start records
byte-unchanged (`:56-59`, `:100-111`) and pins original retention (`:43`, `:112`).

## Findings

### [NIT] F1 — silent semantic usage failure
`trajectory_verifier_recovery.go:326-328` returns 2 with no usage text for
positional-arg/sha256-pairing failures, unlike the sibling handlers that print
usage (`:71-75`, `:214-218`). One-line fix: print the route's usage line (text
already at `trajectory.go:53`) before `return 2`.

### [NIT] F2 — refusal reasons unpinned in the new regression
The duplicate-relaunch (`:46-49`), missing-parent (`:50-52`) and pre-publication
reconciliation (`:53-55`) assertions check only that an error occurs; a
wrong-reason refusal still passes. The same package's adversary table pins exact
messages (`:516-556`); pinning substrings would match that strength.

### [NOTE] Residual coverage boundary (no change requested)
A corrupt non-directory `parent-recovered.json` refuses closed via the
`readRecoveredParent` error path (`:592-614`): the route repairs absence, not
corruption. Atomic publish makes corruption external-only, so operator escalation
is the intended handling.

### [NOTE] Validation status
The focused after-control passed per the coordinator record (after-result.json:
`"exit": 0`, `"seconds": 27.234`, `"pass_events": 3`, `"source_unchanged": true`)
— SECONDARY. Full suite/race/vet runs separately, not yet complete (launch brief;
codex-1-managed-continuation-20260915.md). No passing claim from test source
alone; no acceptance asserted.

## Corrected expected error and diagnostic — verified

Expected error — before: `"differs from the complete validated lineage"`
(before-source trajectory_verifier_recovery_test.go:556; the later
`validateRecoveredParent` error, parent_recovery.go:466-470). After: `"original
requested lifecycle differs from the verifier terminal"` (snapshot `:556`),
verbatim the reconcile.go:319 guard; the planted `AttemptOrdinal = 99` mutation
mismatches `requested.Metadata` vs `terminal.Metadata` at reconcile.go:318, firing
the earlier exact guard before lineage comparison. Mandatory refusal and all state
assertions retained. PRIMARY. Diagnostic — `:243` now reads "cannot write verifier
relaunch result; inspect recover-verifier-parent for the retained replacement
evidence" (before: "...preserve the record and repeat the exact request",
before-source `:243`). Accurate: the encode at `:242` runs only after
`PublishRecoveredParent` (`:308-312`) already published; the named route's replay
then validates and returns the retained record. PRIMARY.

## Existing alternatives — evaluated, not rewritten

No new library machinery: the route reuses the exact
`PreviewRecoveredParent`/`PublishRecoveredParent` the relaunch already invokes
(`:308-312` vs `:331-335`). The missing-parent recovery is a confirmed
non-alternative (refusal mechanics above); the reviewed delta is the four named
files, and the bundle preserves before-copies of exactly those four.

## Scoped support and remaining items

Scoped support: the four-file correction as reviewed; I concur with the
coordinator position for the source-located reasons above. Remaining: optional
F1/F2; await the separate full suite/race/vet before any completion claim;
independent review of the whole candidate remains owed. No signature, no
full-audit acceptance, no pricing or experiment-treatment content.
