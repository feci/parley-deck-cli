---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-11
status: partial
source-manifest-sha256: 6f60552b591b1c23ff448e3c6b94b5672ca5e6164701f32028ba82a643077e84
---

# Patch-regression trajectory execution core

FINAL D6/AC-B2 requires two consecutive independently confirmed material
patch-induced regressions to trigger a trajectory review. No implementation
existed for that rule at the preceding source checkpoint. The new core implements
paired execution and deterministic sequence assessment. It does not yet expose
or enforce an opt-in runtime policy; AC-B2 remains incomplete.

## Implemented behavior

`internal/trajectory.Freeze` binds separate clean Git baseline and patched
snapshots: complete commit identities, actual source-tree digests, exact binary
patch digest, ordered request ancestry, material criterion command hashes and
selected non-implementer verifier identity. The baseline must be an ancestor of
the changed patch. Empty/unchanged patches, nested roots, self-attribution,
duplicate criteria, invalid hashes and excessive inputs refuse. Raw commands,
environment, diff bodies and absolute workspace paths are not in the request.

`Verify` executes the frozen commands through the existing criterion machinery
on both snapshots in AB/BA order. It rechecks both snapshot identities and the
patch around every execution. Returned observations contain actual execution
records and before/after digests; a failed run returns its partial observation
with an error rather than filling missing runs or expected digests. The caller
must retain those observations on failure.

`evidence.RunCriterionDetailed` is an additive execution API that retains the
existing RunCriterion and persisted report shapes. Its typed completeness flag
distinguishes an observed failing test from capture overflow, malformed structured
output, process interruption and cancellation. A failing test in a truncated
stream cannot confirm a patch regression. No diagnostic-prose parsing supplies
that distinction. Ordinary completion attestation remains pass-only.

`Assess` derives a regression from consistent structured pass-to-fail executions
of an original material criterion. There is no supplied pass/fail vote. Opaque,
missing, zero-case, skipped, contradictory, incomplete, self-attributed, wrong-tree
or wrong-command records cannot provide confirmation. Inconsistent AB/BA results,
existing failures and reduced passing coverage remain inconclusive. Aggregate
failure counts cannot establish that the same individual failures persisted.
One confirmed material regression is sufficient even if a different criterion
remains unresolved; unresolved evidence never certifies a clean patch.

`Evaluate` takes the complete frozen expected patch list separately from the
observations. It checks idea, criterion scope, commit/tree continuity, request
hash ancestry, sequence and duplicate identities. Missing observations and
replayed patches refuse. A clean intervening patch resets the consecutive count;
an inconclusive patch cannot supply another confirmation or certify a clean
interval. The first two-confirmation review requirement remains recorded even
if later input includes a clean patch. The API neither authorizes continuing
after that trigger nor claims the implementation complete.

## Verification

The exact set and hashes of all 356 Go/module files match manifest
6f60552b591b1c23ff448e3c6b94b5672ca5e6164701f32028ba82a643077e84.
Full Go suite PASS (124.571s). Trajectory/evidence/budget/app/driver/runner race
PASS (146.321s). Scoped vet PASS (1.237s). Windows amd64 trajectory and app
cross-builds PASS; Windows runtime remains unverified. The compiled trajectory
test binary passed all new cases on an isolated shared volume (60.063s).

Executable fixtures create actual Git histories and separate snapshots, run
actual shell criterion checks, and verify: two different regressing patches;
preserved review requirement after a later clean patch; a clean intervening
patch; unchanged prior failures; actual AB/BA execution order and instability;
source drift; opaque, zero-case and malformed output; exact command scope and
cancelled verification. Mutations of retained fixture observations exercise
coverage gaps, replay, self-attribution, incorrect identity/tree/hash/counts,
interrupted execution and contradictory verdicts. Detailed-execution fixtures
include complete pass/failure, masked test failure, malformed output, a
five-million-byte overflowing stream and cancellation.

Initial builds failed because telemetry.SafeLabel returns a string pointer,
not a string. That type-use error was corrected before trajectory runtime
validation. Both failed-build logs are retained. The missing trajectory API was
source-observed; no executed pre-fix behavioral trial is claimed. All subsequent
focused, full, race and compiled shared-volume runs passed. Logs, manifests and
checksums are under .parley-runtime/trajectory-execution-validation-20260911 and
its referenced sibling logs. These are synthetic executable fixtures with
asserted fixture identities, not independent model acceptance or live pilot calls.

## Required next integration

Persist an explicit opt-in policy and complete expected patch inventory bound to
the real shared fixup charges. Missing policy/history must never become a guessed
empty list. Capture every opted-in patch attempt, including interrupted/failed
attempts, with its exact pre/post source snapshots. Independently select and
invoke a non-implementer verifier through the existing instrumented runner;
freeze its request and verify the actual invocation, helper receipt, execution
records and publication before accepting an observation.

Wire the resulting decision into all relevant driver/manual fixup and completion
paths, preserve the review requirement across resume and provide durable refusal
and explicit recovery behavior. Source snapshots of failed/dirty patch work need
an explicit mechanism; this core only accepts clean pinned Git snapshots and
does not silently commit, discard or pretend to verify dirty state. Full AC-B2
acceptance requires that integrated runtime behavior and fresh independent review.

## Boundaries and remaining scope

The current API trusts the orchestrator to supply the full real expected patch
list and original material acceptance commands. A supplied empty list proves
nothing about real patch history. Runtime labels and mutually consistent hashes
do not authenticate an agent or human against a malicious same-UID caller.
The frozen commands must actually implement the intended acceptance criteria;
paired observations do not establish universal causation, test adequacy or an
independent model invocation. Persisted opt-in policy, runtime wiring, escalation
publication and recovery are not implemented by this source checkpoint.

All other unfinished FINAL requirements remain: safe guard/lock-origin migration,
durable semantic action identity/replay, fresh independent acceptance, complete
real launch coverage, exact packet experiment, frozen 12-task solo/duo/full-six
pilot, blind nonauthor grading, participant-owned final reviews/signatures,
populated-report QA and actual delivery-based follow-up obligations. Historical
quorum/pilot, funding/ceilings and Claude review recovery decisions are unanswered.

Model inventory remains 34 terminal attempts, 17 unknown costs, USD 46.1887585
known CLI estimates and unknown total. No model or real operator accounting
action ran. The latest browser-verified report still represents bd00188; it does
not yet include this core or the preceding reservation-continuity change. No
final merge, release, deployment, global install or immutable-core publication
is claimed.
