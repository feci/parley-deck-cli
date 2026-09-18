---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
review-kind: implementation-candidate-with-local-refutation
---

## Publication interruption in the attended recovery candidate

The new Kimi app path completes its real replacement runner/helper before
publishing parent-recovered.json. If that publication fails, the bound
invocation has been consumed. Repeating relaunch is refused, the old failed
parent-result.json is deliberately retained, and the existing recover-parent
route does not accept that failure. The exposed library PreviewRecoveredParent /
PublishRecoveredParent can derive and publish the complete observation, but the
candidate has no CLI route to invoke them without another attempted relaunch.
The Kimi note treats this as future ergonomics; this coordinator position is
that it is a recovery gap under the original Idempotence & recovery contract.
No participant-owned note is changed.

A native app-entrypoint control ran the actual fake verifier and built helper,
obstructed parent-recovered.json with a directory, observed one real replacement
start and one launch charge, removed only the test obstruction and requested
the dedicated evidence publication route. The unchanged candidate failed that
control in 15.658 seconds because the route is absent. Logs are preserved at
managed-continuation-20260915/app-recovery-publication-control/before.*; this is
local software validation, not a provider or treatment invocation.

## Minimal corrective candidate

Add trajectory recover-verifier-parent --dir DIR --idea ID --run RUN with a
read-only preview and exact --sha256 SHA --yes apply. Reuse the existing fully
revalidated library preview/publication APIs. Do not call a runner, reset a
budget, overwrite the original failure, add a general retry or weaken the
missing-parent path. Apply replay retains the same bytes. Correct the relaunch
output-write diagnostic to name the recoverable observation route.

The test also checks refusal of duplicate relaunch and the unrelated missing-
parent route, no preview write, stale-digest refusal, exact publish/replay, no
extra process or charge, original failure/reservation/accounting preservation,
and subsequent explicit reconciliation. It is running against the correction;
this note does not yet assert a passing outcome or independent acceptance.

Separately, the completed Kimi focused candidate run had 87 passing test/subtest
events: telemetry 0.459s, trajectory 201.143s, runner 30.108s. Its app package
failed only because the rewritten-terminal adversary expected a later error,
while fresh reconciliation correctly refused with the earlier requested/
terminal lifecycle mismatch. The corrected candidate updates that expected
message, retaining mandatory refusal and every state assertion. Prior failed
source/logs and the original candidate remain unchanged.


## Executed correction check

The corrective app command and rewritten-terminal adversary passed in 27.234s:
three test/subtest pass events, no failures, source manifest unchanged. Both the
original failing control and passing correction are retained. The correction
still refuses duplicate relaunch and ordinary missing-parent recovery, then
publishes the retained complete observation with no new charge or process.
Independent Kimi and Zcode source reviews are pending. A full-suite attempt is
running; coordinator fixture omissions (VERSION and live protocol/map) already
mean that attempt cannot establish an all-pass result. Exact missing static
fixtures are prepared separately; no active source was silently repaired.
