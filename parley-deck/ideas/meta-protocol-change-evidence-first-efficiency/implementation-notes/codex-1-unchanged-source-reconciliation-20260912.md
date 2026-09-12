---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-12
base-commit: 0585deb665375fed4b16d73ff0013027b22f03b7
status: implemented-independent-review-pending
---

# Unchanged-source reconciliation

## Delivered behavior

A normally exited charged attempt with identical captured material source can
now be recorded through `trajectory reconcile-unchanged --sequence N`, followed
by its exact `--sha256 SHA --yes`. This is deterministic observation publication
with no attendance gate. The existing attended continuation control is unchanged.
The original attempt and charge remain retained; a retry is not granted.

The new optional `unchanged` branch in ReconciliationPreview binds policy,
complete original attempt, archived scope, invocation and requested/started/
terminal hashes. The attempt digest covers both original sources and archives,
including different commit/status metadata when material bytes are identical.
No independent verifier run, parent result, recovery result or acceptance verdict
is invented. CapturedRequest still refuses unchanged material source as a new
independent patch comparison.

The original quorum and normalized named commands are read from the retained
before-source archive, checked against frozen policy criteria and compared with
the current original contract. A bounded selected-member reader uses the same
full inspector as snapshot restoration, validating every member, canonical
footer, archive digest, complete source digest, links and file stability before
returning any selected bytes. Missing, symlink and oversized scope members refuse.
Normal parent reconciliation uses the same extracted scope parser and preserves
its exact original command/quorum checks.

Matching lifecycle evidence must contain an actual process start, observed exit
and coherent original identities/times. Both zero success and nonzero ordinary
process failure are supported. Missing starts/exits, signal exits, unobserved
handoffs, timeout/cancellation/provider failures and unavailable capture refuse.
The driver may charge before requested telemetry; the measured launch may emit
requested telemetry before charging. Both original orders remain supported when
the request and charge precede the actual start.

All original criteria remain unresolved in an inconclusive assessment. An
unchanged observation neither increments nor resets the material regression
streak. Only a newly confirmed regression advances RequiredReviewSequence; an
unchanged attempt cannot reopen an acknowledged review trigger. Existing
inconclusive acknowledgment is still required for continuation. The latest
unchanged outcome cannot satisfy RequireResolved, even after acknowledgment.

Preview/apply/replay run under the existing cycle guard and rederive original
evidence on every read. Historical state hash normalization is only for comparing
with the exact originally saved preview. Concurrent applies converge on one
resolution; output/publication recovery and later legitimate history preserve
the original preview and record bytes. Replay completes file/directory durability
barriers without rewriting the retained state. Parent recovery dispatches typed
unchanged evidence while revalidating the remaining historical resolutions.

## Compatibility and trust limits

State remains version 3, with additive `unchanged` provenance. Older strict
readers refuse the field; they must not write this state. Frozen version-2 policy
and existing parent resolution encodings remain unchanged. No participant-owned
artifact or signature was edited and no quorum or experiment amendment was made.

Hashes and process metadata detect retained-evidence mismatches, not hostile
same-UID fabrication or authenticated model/human identity. A normal parent exit
does not prove that escaped descendants are inactive. This control records
material equality and observed lifecycle, not general inactivity, correctness,
health, deployment or independent acceptance of this implementation.

## Execution checks

Development trajectory and actual-CLI checks returned exit zero, with native and
shared copies in `.parley-runtime/unchanged-development-20260912/`. Final validation
uses the frozen Go/module inventory in
`.parley-runtime/unchanged-final-validation-20260912/source-manifest.json`.
Final automated validation on 401 Go/module files returned exit zero: focused 155.765s, full 297.826s (32 package terminals), six-package race 352.626s, vet, Windows trajectory/app cross-builds and compiled shared-volume trajectory/app selections (45.861s / 64.312s). Eleven new top-level scenarios are present in focused/full/race and the applicable shared runs. Six removed-protection overlays returned nonzero at their intended assertions. The final verifier checked frozen source/test inventories, native/shared log equality, PE amd64 build headers and the unchanged historical HTML. Windows runtime and independent participant acceptance remain unverified.

Eleven new top-level scenarios cover actual exited processes, precharged and
request-first ordering, empty commits, failed attempts, unchanged charges across
later history, exact preview/replay, missing or contradictory lifecycle, original
archived scope versus rewritten current scope, publication interruption before/
after the state write, rehashed state corruption, complete-archive selection,
real CLI output failure, two concurrent production CLI apply processes and a
five-attempt regression/unchanged/correction history with actual verifier helpers.
The helper-history scenario also reconstructs a lost older parent while typed
unchanged resolutions are present. Those fixture helpers are actual local
processes; they are not independent real-model acceptance of the delivered source.

Six isolated source overlays remove regression-streak preservation, review-trigger
guarding, whole-source archive checking, frozen original command checking,
current contract checking and pinned lifecycle comparison. Each selected test
returned nonzero at its intended assertion. Overlays do not modify the source
under validation and are retained with their exact logs.

## Remaining whole-goal work

Consumed/incomplete helper-ticket, orphan-reservation and other workflow-effect
recovery remain required. The inspected next-source map is
`.parley-runtime/helper-ticket-recovery-next-20260912.md`: persist original intent
before the charge/observer gap, retain all failures and charges and establish
process scope before authorizing new work. Missing terminal evidence or a dead
leader alone cannot prove inactivity or make a retry free.

Fresh participant-owned current-source acceptance, all real launch surfaces,
independent real-model concurrency/closure, the exact phase-1/6 packet trial,
twelve real tasks in solo/duo/full-six arms with rotation and blind nonauthor
grading, final owned reviews/signatures, populated offline HTML with ego-browser
QA and actual-delivery-based 14/30-day follow-up remain obligations. The live
experiments have not run. Resource ceilings and historical roster/pilot amendment
decisions remain unanswered. The historical HTML and assessment stay frozen.
