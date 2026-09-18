---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-11
status: partial
source-manifest-sha256: bcfafa65f3194b3511c7f00b1f49503b736d528d7f4251edf0b508f3ec419782
---

# Changed inactive import recovery disposition

The launch and protocol importers previously retained a failed, inactive import
when history changed before activation, but had no supported way to reconcile
that retained decision with the new history. This control now covers launch,
step, fixup and cross-review imports. It is implementer evidence, not independent
acceptance, Phase 7 consensus or completion of the six underlying audit areas.

## Current behavior

`parley budget migrate recover inspect` returns a read-only preview with exact
logical import-state/current-history hashes, original policy, retained charges,
accounting lower bounds, source metadata, import directory and the last exact
recovery request. `recover apply` requires actual platform attendance detection,
stopped-writer confirmation and explicit decision/count/epoch fields. A participant
field cannot supply attendance. Attendance itself is not human authentication.

Original migration.json and original policy ceilings remain unchanged. A bounded
hash-linked migration-recovery.json journal retains every recovery decision,
observed import-file hashes, typed inventory and reconstructed accounting state.
Launch recovery unions unique observed IDs with all previous charges. Existing
observations, unknown costs and entry timestamps cannot change; a conflicting
terminal record refuses even when its amount is unchanged. Previously accounted
charges survive individual source disappearance without recreating missing raw
files. Unknown/unavailable roots or malformed current history still refuse.

Anonymous pre-telemetry operator totals and protocol-action totals cannot decrease
below prior assertions or observed floors. Anonymous counts are not guessed to
match newly discovered IDs; absent a proven mapping they are conservatively kept.
Those assertions are distinct from actual unique model invocation counts. An epoch
can move earlier to include discovered work, never later. This cannot grant new
ceilings or reset elapsed lifetime; later policy extensions use their own control.

Apply serializes before observing mutable publication files. Both preview hashes
are checked under the existing guard. The journal precedes ledger completion;
only an original or journalled checkpoint may be advanced. Unknown ledger changes
cannot be overwritten. The final recovered-activation marker binds the original
import and complete journal. Older readers reject its different shape. Removing
or truncating recovery authority invalidates activation rather than resetting it.

A stopped write can replay the exact last request, including its original hashes.
If history changes again, a new decision and fresh preview are required; all prior
decisions and charges remain. Exact active replay preserves later charges and
separately approved extensions. A new decision cannot replace an active import.
Original import replay can read a fully recovered active state but cannot bypass
an inactive pending recovery. The journal is bounded at 32 decisions and 16 MiB,
with existing history/source/action bounds retained. No raw prompts, commands,
environment or source bodies are copied into the journal.

## Verification

All 349 Go/module files match manifest
bcfafa65f3194b3511c7f00b1f49503b736d528d7f4251edf0b508f3ec419782.
Full Go suite PASS (143.697s). Budget/evidence/app/driver/runner race PASS
(159.952s). Scoped vet PASS (0.763s). Windows amd64 app cross-build PASS;
Windows runtime remains unverified. An isolated compiled budget binary passed
all recovery fixtures on the shared volume (17.996s), outside real Git ancestry.

Tests exercise all four accounting kinds; count/epoch/observation preservation;
unchanged original import/policy; before/after failures at journal, ledger, policy
and activation publication; exact replay; changed-again history; retained decision
chains; later charges and independent policy extensions; active replacement,
unattended/missing/invalid flags; corruption, missing ledger/witness and conflicting
observations. Separate OS processes meet at an explicit start barrier: identical
requests all succeed with one retained decision; conflicting requests yield one
success and the expected active-replacement refusals, not arbitrary child errors.

Nine actual compiled CLI probes exercised read-only inspection, an isolated initial
fixture import, stale original replay refusal, recovery preview, unattended refusal,
PTY-shaped fixture recovery, exact active replay, conflicting active replacement
refusal and active status readback. Original import/policy bytes and replay bytes
were checked. These synthetic fixture controls do not establish a real human
approval, real project migration or model invocation. No real model ran.

The first build used a nonexistent CyclePolicy method; it was corrected to the
existing originalCyclePolicy helper before runtime testing. The log is retained.
The missing recovery path was source-observed; no executed pre-fix recovery result
is invented. All current focused, process, CLI, full and race runs passed. Detailed
logs, source manifest, executable probes and checksums remain under
.parley-runtime/import-recovery-validation-20260911 and referenced sibling logs.

## Limits and remaining requirements

Stopped-writer assertions and content hashes do not coordinate uncooperative
writers, authenticate a human or protect against consistent same-UID fabrication.
Old writers must not use recovered state. Origin/cache migration remains a separate
unfinished control; this command uses the already verified original lock identity.
Conflicting prior terminal observations and unrecognized ledger snapshots stay
unresolved rather than being guessed away. Budget recovery is not semantic
exactly-once execution and grants no evidence acceptance or participant signature.

Fresh independent acceptance, safe lock-origin migration, durable semantic action
identity/replay, independently confirmed two-patch regression trajectory, full
real-model launch/concurrency/closure coverage, the exact phase 1/6 packet trial,
the frozen 12-task solo/duo/full-six pilot with equal ceilings and rotation, two
blind nonauthor graders, participant-owned final reviews/signatures, final populated
report QA and actual elapsed 14/30-day follow-ups remain binding. Historical quorum,
pilot/funding and Claude review recovery choices remain unanswered. Model inventory
is unchanged: 34 terminal attempts, 17 unknown costs, USD 46.1887585 known CLI
estimates and unknown total. No real project accounting decision was executed.

## Delivery verification checkpoint

Source bd0018865946764cd65cd5cde7860bb5bfd480c9 is represented in the exact
629641-byte offline report. Ten report tests and all 56 ego-browser assertions
passed at 1440×900, 1280×540, 390×844 and 320×720. All sections/documents,
expanded current recovery progress, retained negative history and compile failure,
partial/not-run gates, painted charts, search/reset/pagination, actual ArrowRight
navigation and print-media visibility were verified. Normal media and desktop
metrics were restored. Dedicated cleanup of task space 33 returned done:true.

Full-file SHA256: 06ed736c5ef7379adf1ce3c3d1adc51f5afbb34ef6869b70798f8ad28b87a0be.
Embedded payload SHA256: 2c01b0d2b78666bf30c6549848922688133e3416c0dac13cb320c4a774e5e685.
Browser payload hashing and local full-file hashes bracket QA; file-scheme source
fetch is unsupported. The report was not rebuilt after QA. The exact previous
guard report remains archived. Checks, observations, cleanup and checksums are in
evaluation delivery/2026-09-05/browser-qa/import-recovery-20260911 and the adjacent
browser-import-recovery-20260911.json. Fresh screenshots, physical click placement
and physical print/PDF pagination are unverified. This is partial-report QA only;
all independent/live FINAL gates remain open. The exact 349-file source manifest
and retained validation evidence hashes still match. No model or real operator
accounting action was added by this checkpoint.
