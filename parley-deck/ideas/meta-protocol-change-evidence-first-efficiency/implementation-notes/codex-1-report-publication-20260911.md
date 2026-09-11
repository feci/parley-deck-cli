---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: implementation checkpoint and executed counterexamples
not-a-signoff: true
---

# Serialized evidence publication

All cooperative report writers now share a pinned host-local resource guard
using the already-tested budget lock primitive, with no monetary accounting
side effect. Git administration holds per-report origins outside the tested
worktree. Non-Git check fixtures initialize stable runtime metadata before
hashing; no new digest exclusion is introduced. Aliases share a guard; an
established origin/cache change refuses. This is cooperative same-host
serialization, not same-UID authentication or distributed locking.

SaveIfUnchanged compares the frozen original bytes and replaces the report
within one transaction. The helper retains the guard through receipt writing.
The parent accepts that precise report and receipt under the guard and saves
their digests in the actual persistent run store. Complete requires the same
acceptance, request and receipt before its current-tree gate and final status
write, all under the shared guard. Delayed Save/check cycles cannot overwrite
evidence for status complete. Partial/failed receipt publication is not accepted.

Report and table replacement now use synchronized unique sibling staging and
a directory barrier. Completion status validation still rejects malformed YAML;
the pure transition returns its validated prior status alongside a no-op error
so a delayed writer can recognize complete without reparsing it differently.

The helper now retains actual measured before/after tree values explicitly.
Receipt criteria reconcile by name, and persisted JSON comparison replaces
time.Time implementation-detail equality. The runtime ignore prerequisite is
checked before creating verification files, tracked runtime files and symlink
directories refuse, and Windows verification explicitly refuses the unsupported
POSIX command contract. Windows cross-build is not a Windows runtime result.

## Executed probes

- Two actual test processes snapshot the same report and race CAS: exactly one
  wins, one receives ErrReportChanged. No winner is overwritten.
- Ordinary Save waits for the same publication guard, then refuses after the
  simulated guarded completion. Escaped writers and alias/cancel bypass refuse.
- Built production CLI + actual independent helper/test subprocesses: report
  persistence fails after a real criterion rerun; no acceptance is written.
  Restoring permissions and starting fresh checks/helper completes successfully.
- Actual receipt persistence failure retains the attested report but cannot
  close. Fresh checks and a new helper recover without rewriting the failed
  frozen attempt or pretending that its missing receipt existed.
- Replacing a still-valid report after driver acceptance, or removing the
  accepted receipt, vetoes Complete with the precise expected reason.
- Missing runtime ignore and a misleading probe-only ignore refuse before
  runtime files or an agent invocation are created.

Initial tests exposed a diagnostic mismatch for guard creation failure and a
no-op status error that obscured the finalized-report diagnosis. Those were
corrected. Focused evidence/app probes pass (evidence 0.647s, app 28.478s).
Subsequent full evidence/budget race checks pass (4.489s/6.350s), actual
shared-volume process/alias checks pass (0.635s), vet passes, and Windows app
cross-build passes. Windows runtime remains untested.

The initial full Go test command exited zero, but its shared-volume streamed
log later showed a 1403-byte size with a 133-byte readable prefix and then
1270 NUL bytes. It is retained as a damaged log, not complete validation output.
A fresh full check captures output in memory and writes local and shared copies
once with fsync; its result will be recorded separately. No causal claim about
that filesystem behavior is made from this observation alone.

All source results remain implementer-run and await fresh independent review.
No live real-model verifier closure or complete six-recommendation acceptance
is claimed. Automatic persistent budget policy/scope/action/time wiring, exact
experiments, final own reviews/signatures and pending amendment decisions remain.
