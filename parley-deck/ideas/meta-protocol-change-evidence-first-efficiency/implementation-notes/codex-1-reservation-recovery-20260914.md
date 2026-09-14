---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
status: implemented-review-pending
base-commit: 2ab9e82f0dd5a885717483a363de9482ae7c938a
---

# Original reservation intent and missing-row recovery

The new preparation callback runs under the common cycle guard after the source
observer validates precharge state and before Store.Reserve. It receives the
actual original request, limits, entry key, typed action input and epoch; the
reservation timestamp is supplied only by the subsequent original ledger charge.
The original intent is exclusively created and synced before charging. Attempts
pin its digest. Partial intent files remain failures, never retry permission.

The recover-reservation CLI previews and hash-confirms one exact missing charged
row. Original state, archives, charge identity and every other resolution remain
required. An intent without a charge is a separate observation and cannot apply.
Current changed worktree content is not used to reconstruct an original source.
Recovered accounting has no invented launch, terminal, after-state or independent
verdict. The outstanding attempt still blocks completion and further fixups.
This is accounting recovery, not complete helper/process/workflow recovery.

Initial focused tests passed in 29.403 seconds (budget and trajectory), including
actual process exit at precharge and postcharge boundaries, ten corrupted-evidence
cases, wrong preview, publication failure, preserved spend and unresolved execution.
Evidence: .parley-runtime/reservation-development-20260914/focused-result.json.
These tests are development evidence, not accepted final validation. The wider budget/trajectory/app integration selection passed in 172.754s.
Source inspection then corrected original/extended cycle-policy compatibility;
its dedicated tests passed in 2.526s. Full final, race, platform and shared-storage
checks remain in progress. Original preparation policies must be an exact prefix
of the current validated extension chain; later grants cannot rewrite old limits.

Independent Claude supporting source review of the preceding frozen 2ab9e82
checkpoint is running in evidence-first-reconciliation-review. It does not review
this later reservation implementation, replace a full Phase-6 review or sign for
another participant. All experiment, historical quorum and delivery gates remain.


## Accepted automated checkpoint

Source manifest a5e1f06be6cff9a60322318d2b9dcdb591749ae471ac417c68d46cff97caf1c5
covers 407 Go/module files. Full suite PASS 337.880s (32 package terminals);
six-package race PASS 387.708s; vet PASS. Windows budget/trajectory/app cross-builds
PASS (PE amd64 checked; runtime unverified). Compiled shared-volume budget,
trajectory and app selections PASS 0.556s / 16.633s / 8.412s. All ten new top-level
tests are accounted for in full/race and applicable shared selections.

Five removed-protection overlays fail at their intended assertions: original
semantic action, original policy prefix, preview digest, required intent on ordinary
reads and byte-preserving replay. Native/shared logs and unchanged source inventories
match. The historical HTML remains unchanged. Evidence and final verifier:
.parley-runtime/reservation-final-validation-20260914/final-verification.json.
No test run here is independent real-model acceptance or proves process inactivity.

The supporting Claude review of 2ab9e82 completed with six findings. Its F2 test
counterexample was exercised through two Go overlays: removing the provider class
predicate, and disabling the entire completed/failed qualification, both left the
original contradictory-evidence test passing (2.688s / 2.493s). These are retained
counterexamples, not positive validation. A subsequent owned test correction is
required; this reservation checkpoint does not claim those old lifecycle tests
already distinguish the qualification clauses. F1 abnormal terminal recovery,
F3 quorum binding, F4 exclude-source uncertainty, F5 wording and F6 guarded cost
remain recorded for disposition, without suppressing or dismissing findings.
