---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
source-commit: a223a0d
status: tested-slice-independent-acceptance-pending
---

# Persistent protocol cycle counters

This is Codex's implementation disposition, not a participant review or final
acceptance. Source a223a0d completes the persistent fixup/cross-review counter
slice left open at d0aae57. D6 and the full delivery remain incomplete.

## Executed counterexamples

A read-only git archive of d0aae57 was unpacked into a local temporary directory;
no worktree was registered or changed. Two driver-seam probes failed there:

- A failed fixup at maximum 1, followed by cursor deletion and a new run,
  dispatched another fixup. The failure printed "new run refunded the failed
  fixup: fix-up cycle 1: synthetic fix-up failure [fixup]".
- Four failed cross-review dispatches ran despite a hard cap of 3. The failure
  printed "hard cap granted 4 failed attempts, want 3".

The actual archived-source baseline log is 0.445s, exit 1, SHA256
`af616a2ca35395845b6826a6790d24e6339ab2551bd6791eeb19ba0f41353ed9`.
Both regression probes pass in the final source. The archive/probe remains at
`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-cycle-counterexample-yx5pqlcp/`.

Further executable probes during refinement exposed ignored unknown historical
phases, a missing prompt creating legacy policy, and a standalone handoff
activating policy without execution. A closing probe also observed
`complete <nil> [goal-check complete]` for five spent fixups against a frozen
maximum four when the current runtime maximum was five. The final source
refuses that closure. The initial ordinary full suite had passed before these
additional probes; that pass did not establish those boundaries. Initial and
refinement logs are retained separately from final validation.

## Final behavior

Each idea has distinct frozen fixup and Phase-2 cross-review policy/ledger scopes,
using the existing durable store and verified host-local lock. Git worktrees
share them. A known current cursor/marker floor is frozen as `carried`; failed
attempts, new runs and deleted replaceable cursors cannot reset later charges.
Other positive runs or ambiguous historical requests require explicit migration.
The relative idea path is bound to the scope; missing/corrupt/aliased policy
fields and lost ledger continuity refuse further work.

The actual driver dispatch reserves before publishing its compatibility cursor
or invoking a fixup. Ordinary cross-review and consensus BLOCK reuse the same
boundary. Supported typed manual launches and grouped protocol runners consult
the same idea ledger. One live synchronous operation shares a reservation among
nested children; an independent operation receives a distinct reservation.
Repeated, failed or cancelled execution never refunds an existing charge. An
ended session cannot authorize later work, and a nested changed policy refuses.

Existing track cells stay unchanged. In deliberation the fifth fixup and third
cross-review group are allowed, while the sixth/fourth are refused. Zero cycle
maximum forbids the operation, preserving fast-track no-cross-review policy.
Reviews/signoffs do not spend a fixup. The last allowed fixup can be verified
and closed, while over-cap state or a changed frozen maximum refuses closure.

Strict YAML parsing rejects missing prompt authority, ambiguous keys/types,
unknown tracks and existing track contradictions. Pipeline blocks use their
seeded idea prompts as well. A valid legacy prompt without a track retains
bounded legacy defaults. A stricter current track cannot silently reconcile a
saved higher maximum; a larger current track cannot increase the saved maximum.

Requested/terminal invocation records survive manual launch refusal. A cancelled
typed launch records cancellation without starting a child or charging a cycle.
Standalone unobserved handoff preparation records the handoff without activating
policy or charging work. Unknown historical work still requires reconciliation;
an unobserved request is not proof of zero outside execution.

## Validation

The final source manifest covers 311 implementation/module files. It was checked
unchanged before commit; manifest SHA256:
`f1bd9572d39925f9ef38c1840cdecc512cfcb509801c8c40f267c7b2152d8b23`.

| Executed check | Result |
| --- | --- |
| `go test -count=1 ./...` | PASS; app 104.956s, budget 11.712s, driver 11.907s, runner 31.677s; wall 108.595s |
| `go test -race -count=1 ./internal/budget ./internal/driver ./internal/runner` | PASS; 15.027s / 13.494s / 58.611s |
| `go vet ./internal/budget ./internal/driver ./internal/runner ./internal/app` | PASS |
| `GOOS=windows GOARCH=amd64 go test -c ... ./internal/app` | PASS; cross-build only, no Windows execution |
| Actual isolated shared-volume `^TestCycle` binaries | PASS; budget 1.941s, driver 1.105s, runner 1.397s wall time |

The final full log is 1405 bytes, SHA256
`f5cbd12d7da2b09eeac900fe5cc641f18c6fffdbae3dd012ee20b61010071e7d`.
Race log SHA256:
`f7a4f87635111ee63d1a4f5c18d72a5de5cd314b4f1fbf154fdc01d1ba2eef97`.

Tests include an actual Git-worktree shared cap, independent concurrent
reservation sessions, grouped children, failed real local processes, the mixed
driver-to-real-runner-to-manual boundary, cancelled invocation evidence,
missing/corrupt policy/ledger, nested close revocation, failed durable publication,
unknown historical phases, ordinary/BLOCK limits and final allowed verification.
These are local-process and driver fixtures, not independent live-model review.

Shared-volume binaries were compiled into local temporary storage and executed
with TMPDIR at `/Volumes/My Shared Files/AI_WORKSPACE/parley-cycle-fixtures-4kzeoqup`.
`git rev-parse --show-toplevel` confirmed the fixture was outside any enclosing
Git repository. Output was captured through pipes/in memory and written once
with flush/fsync; logs were not streamed to the shared volume. Fixture root is
preserved. Exact logs, commands, checksums and the archived negative probe are
copied into the integration worktree's ignored
`.parley-runtime/persistent-cycle-validation-20260911/`.

The refreshed offline report is 599647 bytes, SHA256
`ca92c4bda9170e92ce2fa02ceab1c9aee26fd013beae6823fadec956dca02357`,
built at 2026-09-11T14:27:08.773519+00:00. All 10 report tests pass.
An ego-browser checkpoint smoke check loaded this exact local file and verified
the a223a0d text, partial status, 1440x900 desktop, 390x844 and 320x720 mobile,
and 1280x540 low-height view without document-width overflow. Screenshots were
visually inspected; both charts have painted pixels. Search produced zero
matches, reset restored 673 ideas, and Tab moved from search to Project.
Print emulation displayed all five main panels and hid the navigation through
its parent aside (zero navigation bounds). This does not verify physical print
pagination or replace final populated-report QA. The dedicated ego task space
was closed after the successful smoke check; other spaces were untouched.

## Remaining obligations and limits

The initial carried-count bootstrap is narrow, not a general legacy migration.
Explicit operator migration/extension, safe guard recovery, durable semantic
action replay, monetary default mapping, canonical durable refusal publication
and independently confirmed material patch-regression trajectory remain open.
Random reservation identity and live grouping do not prove exactly-once work
across process crashes. Typed phases cannot classify arbitrary prompt text or
commands outside Parley. Same-filesystem-access processes are not authenticated
humans merely because the accounting files carry attribution.

Fresh independent source acceptance and real-model concurrency/closure remain
required. Claude's earlier owned review stays preserved; the subsequent 2194dca
review terminated at the weekly limit without an artifact. Its recovery choice,
historical quorum amendment, full-six pilot amendment and funded experiment
ceilings remain unanswered. No response, retry or replacement is inferred.

The exact packet trial, twelve-task solo/duo/full-six pilot, two blind nonauthor
graders, owned final reviews/signatures, final report browser QA and genuine
14/30-day follow-ups remain required. Inventory is unchanged: 34 actual terminal
model attempts, 17 unknown costs, USD 46.1887585 known CLI estimates and unknown
total. Local fixtures add no model calls. No final merge, release, deployment,
global install or immutable-core publication occurs in this slice.
