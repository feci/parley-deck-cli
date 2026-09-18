---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
source-commit: 321aef9
status: tested-slice-independent-acceptance-pending
---

# Finite operator extensions for protocol cycles

Source 321aef9 adds the actual operator control required to extend an existing
fixup or Phase-2 cross-review ceiling. This is Codex's implementation disposition,
not another participant's review or all-D6 completion. No real idea received a
grant, and no participant model was invoked in this slice.

## Behavior and authority

`parley budget cycle inspect --dir DIR --idea ID --kind fixup|cross-review`
is read-only. Its JSON includes the policy, canonical policy SHA256, current
spent count, original activation time and ledger path. It does not initialize
an absent scope or its lock.

`parley budget cycle extend` requires that same scope/kind, an attended control,
`--expected-policy-sha256`, `--decision-id`, `--max-cycles`, `--reason` and `--yes`.
The maximum is finite and absolute; it must exceed the saved maximum and actual
spent count. A skipped phase with original maximum zero cannot be enabled.
Participant frontmatter and runtime flag changes remain unable to issue grants.

An extension atomically publishes policy schema v2 with the unchanged original
maximum and carried count plus an ordered hash-linked decision history. Each
decision retains its time, previous policy hash, finite maximum, reason and
spent count observed while applying it. The charged ledger, action identities
and original activation time are unchanged. A unique decision ID is scoped to
this idea/kind policy. Exact replay is idempotent even after later grants;
conflicting reuse or a stale preview for a new decision refuses.

Reservation and extension publication share the existing resource guard, then
use the existing ledger boundary. A cached binding reloads the current valid
policy before reserving. Concurrent decisions against one preview admit only
one grant. Publication failure before replacement leaves the prior policy;
failure after actual replacement can leave the recorded grant. Inspecting or
replaying the exact decision recovers without duplicate grants or refunded work.

Driver, ordinary/BLOCK cross-review and supported typed manual runners actually
consume the saved effective ceiling. They retain the original track settings,
non-solo rules, strict closing-review loop and independent evidence gates.
An extended final allowed fixup can close after its existing checks; the grant
does not waive those checks. Live nested children continue sharing their one
cycle when a valid extension is published between children. Distinct operations
still reserve distinctly; this is not exactly-once semantic action replay.

V1 policies remain readable. Extended policies require v2 support; old readers
refuse the unknown fields/version instead of ignoring the operator history.
The history is bounded to 128 decisions and a 1 MiB policy, with 128-byte IDs
and 1024-byte reasons. Zero/unlimited and an integer maximum lacking room for
the following round ordinal are refused. Canonical hashes address policy
meaning rather than JSON whitespace. They are not human authentication or
protection against a same-filesystem-access process rewriting an entire history.

## Executed validation and correction

The first whole Go and combined race runs found one regression: the exhausted
fixup diagnostic said "another attempt" instead of naming "attempt 2" after
one charged attempt. The existing
`TestOverCapEscalationReportsChargedAttemptsNotCycles` failed in both runs.
The code now reports charged count and refused ordinal separately, using a
representable unsigned value for the diagnostic rather than overflowing `int`.
The failure logs remain retained; passing targeted cycle checks did not conceal
that existing reporting obligation.

Only `internal/driver/impl.go` differs between the initial and final source
manifests. The final 315-file implementation/module manifest has SHA256
`6005a0aeb40aa6df9391815e4ed3d10d83415650dcb571713a28854b52a018b2`.
It was verified unchanged before committing the source.

| Check | Observed result |
| --- | --- |
| Final `go test -count=1 ./...` | PASS; app 88.994s, budget 12.375s, driver 9.884s, runner 30.044s; wall 92.425s |
| Budget and runner race in the first combined run | PASS at unchanged package source; 19.136s and 60.219s; that command as a whole failed the driver diagnostic test |
| Corrected `go test -race -count=1 ./internal/driver` | PASS; 12.294s, wall 13.899s |
| Scoped vet, plus corrected driver vet | PASS |
| Final Windows amd64 app cross-build | PASS; Windows execution untested |
| Shared-volume cycle tests | PASS; budget 3.008s, corrected driver 1.611s, runner 1.673s, app control 0.444s wall time |
| Actual CLI through pipes without a terminal | Exit 2 with attended-control refusal; no runtime state created |

The final whole-suite log is 1403 bytes, SHA256
`9d9e57493624088f5665bbac2e8e7e115a2ca414e861178095fe728b321c7854`.
Corrected driver race log SHA256:
`4c900c91e82600b1bfe0ebd107ede00824d65711eeb352d24e994f5e614a7d8f`.
Corrected shared-driver log SHA256:
`8f4280459cc7142eb936f3fdd3c71db6067fd15f8f4defc7169f7b4ba9019f64`.

Tests cover exact/conflicting/stale decision replay, simultaneous operator
decisions, unchanged ledger bytes and clock, old cached bindings, Git-worktree
sharing, missing/corrupt/aliased grant fields and hash chains, lost prior charges,
publication failure on both sides of replacement, invalid/overflowing inputs,
clock regression, cancellation, skipped-phase refusal, nested grouping, actual
manual child execution/refusal, resumed driver dispatch, BLOCK beyond the
original cap and verification of the last extended fixup. No real operator
approval, authentication or independent model execution is inferred from fixtures.

Test binaries were compiled to local temporary storage and executed with TMPDIR
at `/Volumes/My Shared Files/AI_WORKSPACE/parley-cycle-extension-fixtures-jyzidcz3`.
A Git probe confirmed the root was outside an enclosing repository. Output
used pipes/in-memory capture and one flushed/fsynced write, not streaming to the
shared volume. The actual compiled CLI without a terminal printed:
"budget cycle extend requires an attended terminal and an explicit operator
decision; participant frontmatter is not a grant". That fixture granted nothing.

Initial failures, final results, source manifests, commands and checksums are
retained in the integration worktree's ignored
`.parley-runtime/cycle-extension-validation-20260911/`; original local storage is
`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-cycle-extension-validation-py4o0dtv/`.

## Remaining scope

Legacy accounting migration, launch/driver-step extensions, lock-origin recovery,
durable semantic action replay, monetary default mapping, canonical refusal
publication and independently confirmed material patch-regression trajectory
remain required D6 work. The finite cycle control does not solve those by itself.

Fresh independent source acceptance, real-model concurrency and closure, exact
phase-packet experiments, the twelve-task solo/duo/full-six pilot with frozen
equal ceilings, two blind nonauthor graders, owned final reviews/signatures,
final populated-report QA and actual 14/30-day observations remain open. The
historical quorum, pending Claude recovery, proposed roster/pilot amendments
and funding decisions remain unchanged and unanswered. There is no new model
call: inventory remains 34 terminal attempts, 17 unknown costs, USD 46.1887585
known CLI estimates and unknown total. No final merge, release, deployment,
global install or immutable-core publication is performed.

## Report publication validation

The September 11 continuation verified all 315 current implementation/module
files against the tested source manifest and all 19 retained evidence files
against checksums.json. No source change was needed for this publication step.
The refreshed English offline report passes all ten existing report tests.
It is 603127 bytes, built at 2026-09-11T15:04:26.345153+00:00, SHA256
`6efb779f8d56d86d9ef49dfc6dc7d525784d0b75f1aec97ee3d1698868eda1a1`.
Earlier browser QA addresses the previous report hash and does not establish
exact-file browser verification for this content refresh.
