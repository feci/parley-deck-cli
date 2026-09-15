---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-11
status: partial
reviewed-source: 0cdc56e
---

# Publication review dispositions and launch policy attachment

Claude's own `claude-1-publication-review-20260911.md` is preserved unchanged
at fbb3c31. It reviews 0cdc56e by source inspection, not executed tests. Its
SHA256 is d6d50d66c764ee0acb5a2354277f6e5981e8f29283c660e531e284e988407042.
This is an implementer's disposition, not independent acceptance or Phase 7
consensus. Every original severity and unresolved finding remains visible.

## Dispositions of every finding

| Original finding | Current disposition and evidence |
| --- | --- |
| CRITICAL: deleting cursor, checks and report removes closure gate | Corrected in current source, pending independent acceptance. The separately persisted original checks witness is pinned under the report guard before Driver.Advance starts work and before direct named RunChecks hashes/executes criteria. ObserveChecksContract consults it independently of the replaceable cursor/report. The actual direct Complete and Driver.Advance triple-deletion probes both reproduced false closure before the correction and refuse afterward. Complete deletion of Git/runtime administration remains outside this cooperative boundary. |
| MAJOR: lost cache or changed host denies evidence publication | OPEN recovery limitation. The refusal is real. The suggested unconditional inode recreation is rejected for a specific correctness reason explained below; no supported migration or automatic recovery is claimed. |
| MAJOR: CRLF and nested status cause refusal | Corrected in source, pending independent acceptance. The byte-preserving transform accepts CRLF delimiters/status while recognizing only an unindented top-level normalized status; nested fields remain bound content. The combined old-source test failed first on nested LF status, before its CRLF iteration; no separately executed pre-fix CRLF failure is claimed. Both forms pass after correction. Non-normalized top-level status still intentionally refuses because the authorized transition binds exact bytes. |
| MAJOR: refused verification has no committed canonical record | OPEN. Ignored runtime requests, invocation records and helper receipts remain locally retained and do not establish commit durability. A canonical bounded refusal publication path is still required; it must not introduce an excluded unbound hiding place. |
| MINOR: before/after digests called tautological | Reasoned disagreement, pending reviewer response. Both values are independently measured before being accepted against the frozen request. A successful attestation necessarily contains equal values because a changed tree has already been refused. This does not make either measurement fabricated. The saved result documents successful validated observations; it does not promise that a failed mismatch is encoded as a successful receipt. An extra identical hash merely to make a downstream assertion reachable is not added. Failure-observation publication remains the separate open finding above. |
| MINOR: Git/ignore/non-Git diagnostics conflated | Git absence and non-Git worktrees now get separate diagnostics. The ignore prerequisite still runs before runtime artifacts are created. Non-Git independent closure remains unsupported, not silently accepted. |
| MINOR: non-Git origin affects source digest | Explicit limitation retained. Direct non-Git checks initialize stable local origin/witness metadata before hashing; their source identity includes host-specific origin content. No cross-host comparable digest or independent non-Git closure is claimed. No broad exclusion was added. |
| MINOR: unanchored replaceSection and mutable evidence table | Splicing corrected in source: actual level-two line headings outside backtick/tilde fenced blocks are recognized; inline examples and level-three subheadings are preserved. The generated human table remains excluded and is not integrity-certified independently of typed EVIDENCE.json. A table-binding requirement remains unresolved; do not cite the mutable table alone as proof. |
| NIT: reentrant Save deadlocks | Exported Save and WithReportWriter now document the non-reentrant contract; callbacks use their supplied writer. Existing cancellation/escaped-writer tests remain relevant. No reentrant lock is claimed. |
| NIT: shallow copy for transition comparison | The expected report is deep-copied through JSON before authorization. Records/maps cannot share the mutable comparison structure. |
| NIT: positional provenance restoration | Original criterion provenance is indexed and restored by name; unknown names refuse. Existing exact original/report reconciliation remains in force. |
| NIT: permissive report Load | Load now uses a bounded regular-file stable-open reader, unknown-field refusal and trailing-value refusal. This is not authentication of arbitrary same-UID edits. |
| NIT: Windows Complete diagnostic | Named-criterion Complete now explicitly refuses with the POSIX-only execution contract before trying to interpret absent acceptance. Cross-build is not Windows runtime verification. |

### Why cache-lock recreation cannot be accepted as proposed

Compare-and-save is a read, compare and atomic replacement inside one shared
lock. Atomic rename alone is not an atomic filesystem compare-and-swap. If
writer A holds the old unlinked inode, recreating the path permits writer B
to take a different inode. Both can compare the same original bytes before
either replacement; both can then publish different reports/receipts. A's
late write can also race B's completion. No money needs to be stored for this
exclusion requirement to apply. Correct recovery must establish quiescence
and one shared lock identity; it cannot infer that an old descriptor is gone
from a missing pathname.

Claude's distinct-mount-alias observation is also retained. Symlink and case
aliases are coordinated; separate mount paths for the same remote directory
are not established equivalent and are not certified by the current guard.
No distributed or hostile same-UID protection is claimed.

## New precharged launch-policy binding

The attended `parley budget configure` command requires explicit launch,
cost and wall-clock ceilings (zero means unlimited), plus a conservative
reservation whenever a monetary cap is configured. Its initial policy is
frozen; exact replay preserves the original activation clock and charges.
Git worktrees resolve one common-administration scope, preserving nested
deck prefixes. Non-Git roots have local scope only. Missing/corrupt policy in
an existing scope refuses instead of silently becoming unlimited.

The common runner boundary reads this persisted binding for every actual
process launch; an explicit programmatic policy cannot override it. An
unobserved handoff remains non-process work. The six actual fixture surfaces
are manual, round, consult, probe, terminal and ACP. Missing configuration
still means no configured launch policy; this is not a claim that all
existing driver/loop policy is now enforced.

A durable ledger-established witness is persisted after the charged ledger
and before granting work. Deleting only ledger.json was reproduced to
reset a one-launch cap on the prior source; the corrected store refuses.
Witness publication failure retains the conservative charge without a grant.
The witness does not reconstruct lost charges or resist deletion of all history.

### Fresh historical-accounting counterexamples

Before correction, TestLaunchBindingRefusesIncompleteHistory reproduced
three fresh-policy grants over unknowable history: an invocation directory
with no requested.json, an invocations path that was a file, and a symlink
to different history. The old glob silently skipped the missing request,
ignored directory-reading errors and followed the alias.

Configuration now enumerates real directories with error propagation, requires
a readable bounded request for each recorded invocation, and refuses aliased
runtime/history paths. Refusal occurs before partial policy activation.
Concurrent distinct first configurations produce one frozen winner.
The unsupported-platform configure diagnostic now describes activation rather
than incorrectly labeling it reconciliation.

The recovered pending test session 66526 completed successfully: budget 0.958s,
app 30.392s and runner 2.046s. After the additional history correction, selected
binding/configure/six-surface tests pass: budget 1.127s, app 0.408s, runner 1.796s.
Full, race, shared-volume, vet and cross-build results are recorded separately
after completion; older passing runs do not certify newer source.

### Full validation and corrected integration failures

The first full run failed on an old CRLF-refusal test and on a canceled launch
being classified as budget_refused. The old refusal fixture now tests exact
CRLF transformation and inverse bytes. Context cancellation remains cancellation
through policy loading/reservation; it is not reported as exhausted budget.
The failed 2229-byte captured log is retained at
.parley-runtime/continuation-full-captured-20260911.log, SHA256
a7c99e2369e294d34eea9b242048d9176882dde28c9492f93af538fc53476d0f.

Nested decks absent from another accessible worktree are now handled separately
from an unavailable worktree. An existing policy load needs only its shared
scope, not a new scan of historical worktrees. First activation still refuses
an unavailable historical worktree. A dedicated fixture proves nested-prefix
separation and inheritance when the corresponding deck is created later.

The corrected full go test -count=1 ./... passes: app 92.263s, budget 4.117s,
driver 5.357s, evidence 3.038s and runner 26.011s. Its 1402-byte captured log
has SHA256 0411e1e3ad1b24eaff7792d9936c65cbbd450a9a737f32515393ea4679177c33:
.parley-runtime/continuation-full-corrected-captured-20260911.log.
The independent local capture is at
/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-final-binding-checks-8ugdgbqr/full.log.

Budget/evidence/runner race suites pass (6.059s/4.247s/44.763s); scoped vet and
Windows app test cross-build pass. Windows runtime remains untested. Logs are
under /var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-scoped-final-wowsyr0b/.

Shared-volume runtime fixtures pass with a locally built budget.test binary,
using the evaluation workspace outside any enclosing Git worktree. Earlier
attempts that placed fixture roots inside this integration repository inherited
its actual Git administration and correctly encountered unavailable historical
scratch worktrees; these are retained failures, not passing non-Git fixtures.
The passing shared log is
/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-shared-isolated-bin-ik1rjpp7/shared.log.

## Coverage and unfinished obligations

Claude did not inspect readiness, driver/impl.go, Windows locking or the new
test files in this review. Its open questions about those paths remain
coverage obligations. The original witness now runs at Advance entry, before
dispatch; completed ideas rebuild PhaseDone and do not need to replace their
report. Those source observations are not substituted for independent review.

Persistent driver-step/fixup/cross-review limits, legacy migration, explicit
operator extensions and the two-confirmed-regression trajectory remain
unfinished. Live-model verifier closure/concurrency evidence, complete live
launch coverage, exact packet trial, full-six pilot, final independent review
and participant-owned signatures are still required. The historical quorum
and frozen experiment are unchanged; pending user decisions are not approvals.
