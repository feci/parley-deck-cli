---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
source-base: d5d57e14d5ed5a2b68876070cd0e9f5f6b72ef77
status: integrated-source-validation-passed
---

# Preserve existing-ticket stop authority when historical evidence disappears

StopCapturedVerification now retains its original policy, complete published
charged ledger, structural state and exact ticket/request/launch/process
attribution while omitting historical archive, reservation-intent content and
resolution-content reads.
An already issued ticket can therefore be stopped after losing its baseline
archive or an earlier resolved invocation's terminal. The operation grants
neither another execution nor acceptance of incomplete evidence.

## Implementation boundary

internal/trajectory/state.go adds private withStateControl through the same
guarded state-authority reader. Only the false/stop branch of
withVerificationAuthority selects it. Execution, ticket preparation/reservation,
reuse, reconciliation and acceptance continue to use full source/result checks.
Both routes inspect the published accounting ledger and validate structural
transitions under the existing resource guard. The stop route omits
checkSourceSnapshots, including its separate retained-intent content validation;
the original intent hash remains part of the exact ticket-bound attempt. No policy/state/archive format
changes, charge refunds, synthesized outcomes or general recovery grants exist.

The stop continues to bind the original launch and preserve its first stop.json.
An existing process is signalled only with its original boot/start/group/command
attribution. Missing or changed structural authority still refuses. A missing
helper claim can establish that the stop won before this protocol could issue a
new claim, but cannot certify arbitrary descendant inactivity. Older incomplete
process identities and escaped sessions still need explicit recovery treatment.

## Direct checks and falsification

Four new top-level tests in verification_cleanup_test.go exercise real fixture
charges and original authority:

- Registered criterion execution remains stoppable after baseline deletion.
  A changed start identity must specifically report process start time mismatch
  while the actual fixture-owned process remains alive. Restoring test-owned
  identity permits cleanup; the stopped journal cannot supply accepted evidence.
- A real prior unchanged process is reconciled and continued, then a second
  real changed attempt prepares/reserves its ticket. Removing the first
  invocation's terminal blocks ordinary Inspect but still permits the second
  ticket's stop and byte-identical stop replay, without a helper claim.
- A stop before helper execution survives baseline deletion and exact replay;
  later execution still refuses and no helper claim is created.
- Coordinated test-owned state/request/launch charge rewrites cannot replace
  the original published reservation. The exact ledger predicate must refuse
  and no stop is published.

All four plus four surrounding control tests passed in 26.040 seconds. Four
integrated-source negative controls then failed at their intended assertions:
old archive-dependent route (2.199s), old result-dependent route (2.805s), removed
original-ledger validation (2.117s), and a no-signal mutant pretending process
attribution succeeded (4.524s). The last mutant never signals an unknown process.
Native overlays leave integrated production files unchanged; test-owned cleanup
retains custody of the actual fixture child.

Evidence: .parley-runtime/cleanup-authority-development-20260914/ and the
native directory recorded there. The earlier six-test prototype (23.451s), its
separate ledger test (2.403s) and two prototype negatives retain their own source
provenance in .parley-runtime/cleanup-authority-prototype-20260914/.

## Validation in progress and retained failure

The first full-suite attempt failed in 53.801 seconds: Go compilation/linking and
17 configuration tests reported native disk exhaustion. Its complete native bytes
and result are retained as initial-disk-full-full.jsonl and
initial-disk-full-full-result.json in
.parley-runtime/cleanup-authority-final-validation-20260914/. Source hashes still
matched. That attempt is not a passing suite. Only rebuildable Go compiler/test
cache was selected for clearing; native fixtures, logs and participant worktrees
remain retained. The completed retry and other planned checks will be recorded
separately after termination.

## Remaining audit limits

This removes one unnecessary cleanup dependency. It does not reconstruct
historical or abnormal launch terminals, recover a consumed/incomplete helper,
establish process quiescence or repair workflow effects. The separate full-size
history witness now reproduces both Inspect and Run.Finish deadline failures;
this stop-only correction does not alter terminal publication or resolve F6.
Current-source participant review/acceptance, live launch/concurrency/closure
checks, packet/full-six pilot experiments and delivery follow-ups remain open.
The earlier Claude-owned source review and historical HTML are not write targets.


## Completed checkpoint validation

The 417-file Go/module manifest is 19fc844af60bcc3795139ed9bbeb15e0942431f7233b2a2e12669bfb536f36a3.
The unchanged-source full retry passed in 359.222s with all
32 package terminals; six-package race passed in 422.705s;
vet passed. Windows production CLI and trajectory test cross-builds passed;
PE amd64 headers are checked and Windows runtime remains unverified. The
compiled shared cleanup/verification selection passed in
36.021s. All eight focused top-level tests
are present in the full, race and shared results. Four intended negative
controls, exact manifests, native/shared logs, the unchanged old participant
review and historical HTML were checked.

The original disk-exhaustion failure is preserved separately. Clearing only
rebuildable Go cache increased available native storage from 171,974,656 to
8,801,476,608 bytes before the retry. This is an environment recovery, not a
source correction or a reclassification of that failed attempt.

The later native-only terminal prototype binds Begin's exact state in its
runtime handle and changes only the prospective terminal-recording authority
route. Its four focused cases pass, a removed-state-binding control fails at
the intended assertion, and the retained full-size publication fixture completes
in 4.020197291s. That candidate is not part of this source manifest and does not
resolve the remaining F6 full-history-read or historical-recovery obligations.
