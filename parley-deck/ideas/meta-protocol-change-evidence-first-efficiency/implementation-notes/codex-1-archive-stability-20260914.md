---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
source-base: 83eec56fe5c585c6aebac407e54f3e6c2bf79b6c
status: implemented-awaiting-independent-review
---

# One strict full reread for timestamp-only archive transitions

The original shared 16 MiB history probe safely refused a snapshot after all
archive/canonical/tree hashes and link checks passed. Inode 8559299, size
16,782,848 bytes and private mode 0600 matched; mtime moved by 66,835 ns.
The original failure, one passing diagnostic and PASS/FAIL/PASS diagnostic
repetitions remain under .parley-runtime/history-cost-probe-20260914/.
The underlying host cause is unproven. No old failure is relabeled successful.

## Implemented boundary

The first descriptor stays open. Initial, opened, ended descriptor and named-path
file identity, mode and size must match after complete hashes and link validation.
A timestamp-only transition permits one full strict reread of the same original
file. That pass checks all bytes and original metadata and forbids another
timestamp transition. Replacement, other metadata/content change or any reread
failure refuses. There is no sleep, skipped hash or unbounded retry. Both reads
respect the caller context. Archive encodings and content addresses stay fixed.

The first pass's selected/extracted regular bytes are already bound to expected
complete hashes. The second pass receives no destination/member, so it cannot
append duplicate selected content or overwrite extraction. The first pass returns
success and creates final symlinks only after that strict check succeeds.
RestoreSnapshot's existing private-directory removal on failure remains intact.
This does not authenticate against a coordinated same-UID attacker or prove
inactivity of any model/helper/criterion descendants.

## Validation plan and prototype evidence

The earlier native-only candidate passed two focused tests/eight scenarios in
4.590s and one actual shared 16 MiB/five-attempt history probe in 58.957s. Six
predicate-removal controls and two after-hash boundary controls failed at their
intended assertions; the protected after-hash boundary variant passed. Those are
prototype evidence only. They are retained under archive-stability-prototype-20260914.

Integrated snapshot tests passed in 10.355s. Six regenerated single-predicate
controls and both after-hash boundary controls failed at the intended assertions.
The protected boundary Go test passed in 1.619s; its wrapper initially rejected
Go's whitespace after ok. That wrapper failure is retained and its output parser
was corrected without rerunning accepted tests. Final frozen integrated-source validation passed as recorded below. Deterministic boundary
injection remains outside shipped production source. Do not present prototype
measurements as a current-source benchmark or participant-owned acceptance.

F1 abnormal/old-precharge recovery, consumed/incomplete helper tickets and
workflow effects, the F6 full size/history envelope, fresh independent reviews,
real launch/concurrency/closure evidence, original quorum/signatures, frozen
packet/pilot experiments, final offline HTML/ego-browser QA and delivery-based
follow-ups remain open. No live model, release, roster or global-core change is
part of this correction.

## Frozen integrated-source validation

The 415 Go/module files are pinned by manifest b579964cc39ded6b23c8fc41e93a74e3cb1acc39332e6fd6237f1fad1f1e3397.
Full Go suite PASS (357.367s, all 32 package terminals); six-package race
PASS (408.239s: budget/app/driver/runner/trajectory/evidence); vet PASS.
Windows production CLI and trajectory/app test cross-builds passed; PE amd64
headers checked. Runtime is unverified. The evidence test binary's pre-existing
syscall.Mkfifo compile failure remains retained in earlier F4 evidence.

Compiled shared trajectory snapshot/quorum selection PASS (62.476s)
and selected trajectory app tests PASS (77.576s). Actual shared
16 MiB/five-charged-attempt history fixture PASS (31.16s), using
test-fixture overlays only; integrated production source has no overlay.
Five-attempt Inspect samples: 781.839, 840.537, 862.23 ms. This is warm-read,
single-machine evidence, not the full 128-attempt/256 MiB envelope or cold cache.
No claim of a 30-second violation follows from it.

The final verifier checks source hashes, all package terminals, both new tests
and eight scenarios, eight intended negative controls plus the protected boundary
pass, fixture-only probe provenance, native/shared output equality, Windows
headers and unchanged participant review/historical HTML. Original shared
timestamp failures and both original F3 race failures remain retained as failures.
Evidence: .parley-runtime/archive-stability-final-validation-20260914/ and
.parley-runtime/archive-stability-development-20260914/.
