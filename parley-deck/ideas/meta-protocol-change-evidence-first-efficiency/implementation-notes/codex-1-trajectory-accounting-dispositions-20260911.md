---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
phase: implementation
date: 2026-09-11
---

# Persistent trajectory capture and refusal enforcement

## Outcome and remaining acceptance

The shared fixup policy can explicitly require durable trajectory capture. The
actual driver/manual reservation boundary records its exact charge and before
source observation; the instrumented process boundary binds that attempt to
its real invocation and retains its actual terminal/post-state observations.
A failed or interrupted attempt stays spent and pending. Missing coverage,
changed charge identity, repeated launch on the same charge and pending history
refuse another fixup and completion. Linked worktrees and new processes consult
the same retained state. This advances FINAL D6/AC-B2 but does not complete it.

Every captured attempt currently remains pending. The paired AB/BA execution
core at 26a6b434dc9433597e17853a86c05d2e3a2caaaa is retained, but independent
verifier invocation, helper receipts, durable execution publication and explicit
review/disposition recovery are not yet connected to this capture state. There
is no caller-supplied acceptance flag and no reset/override command. No actual
project policy was activated as part of this work.

## Durable authority and boundaries

- `parley trajectory initialize` creates/reuses the original fixup accounting;
  the selected finite ceiling must match the protocol track. `configure` previews
  original named material criteria, implementer and current clean source, then
  requires the exact trajectory/cycle hashes and attended activation. `inspect`
  exposes retained state. The CLI attendance seam in tests is synthetic.
- The policy SHA is a required field in the shared cycle policy. Existing readers
  reject the unknown field rather than silently running an opted-in policy.
  Activation requires uncharged, unimported accounting without prior extensions.
  Existing history must not be silently treated as a new empty experiment.
- Capture initialization is published before its cycle-policy reference. A
  retained initialization without its reference refuses new work and needs exact
  activation replay. Conflicting initialization is not overwritten.
- Budget callbacks run under the existing common cycle resource guard. A failed
  before-observation grants nothing; a failed post-reservation publication leaves
  its charge spent. Later reads compare the complete ledger entry set with the
  attempt inventory and each original reservation's scope, epoch, key, kind,
  timestamp and exposure. An orphan charge is unknown history, not zero attempts.
- The runner binds the actual launch to the selected implementer and its unique
  instrumented invocation ID. It refuses a second invocation on that reservation.
  Readiness/review/signoff work retains its existing separate budget semantics.
- Source observations bracket the code-tree digest with actual Git observations.
  They retain commit, source/status digests and clean/dirty state, never raw
  diff/status output, environment or commands. Post-state capture runs even when
  the process was cancelled or failed. An unavailable digest is explicitly
  recorded without filling in the expected digest.
- A retained live session/run also pins the required trajectory policy. Its
  disappearance or replacement before launch/terminal publication cannot be
  reinterpreted as disabled observation. Terminal capture failure is a runtime
  integrity failure even when the process itself exits zero.
- Both the deterministic driver close branch and application completion writer
  refuse pending trajectory evidence. Separate finite budget extensions preserve
  the trajectory reference and cannot resolve its independent evidence gate.

These are cooperative runtime controls. Runtime callbacks, process labels,
terminal attendance and consistent same-UID files are not human authentication.
Removing every durable copy of an accounting origin remains a broader origin/
migration problem; this slice does not claim external protection against total
state erasure. Active handles and retained policy references do detect loss.

## Refutation and validation scope

Executable fixtures exercise a real shell command that changes source and exits
7, an interrupted sleeping process after a source write, and a child that moves
its accounting directory then exits zero. Manual `RunMeasured`, grouped `RunFixup`,
the driver reservation boundary and both completion paths are covered. Tests
also cover orphan charges, removed state/reference, empty history, duplicated
JSON, same-count entry substitution, changed reservation time, changed baseline,
wrong implementer, duplicate invocation, unavailable post-state, exact activation
replay, extension compatibility, and a separate process using a linked worktree.
Tests preserve dirty files instead of committing, stashing or discarding them.

The initial live-handle implementation reproduced two incorrect successful
returns after the whole referenced accounting directory was displaced: one
before launch and one during terminal publication. The retained negative log is
`.parley-runtime/trajectory-accounting-live-negative-20260911.log`. The corrected
cases refuse, and the real runner child additionally produces a terminal
`trajectory_failure` instead of successful capture.

Earlier fixture corrections are retained: an in-tree valid symlink was wrongly
chosen to model an unavailable source (changed to a dangling link), the grouped
agent fixture lacked its discovery `Found` flag, and the driver fixture needed
its own Git initialization. These were test-fixture mistakes, not executed
pre-fix product regressions. The first broad suite/race/platform pass remains in
`.parley-runtime/trajectory-accounting-validation-20260911/` and is not used to
certify the subsequently corrected live-handle source.

Current-source validation belongs in
`.parley-runtime/trajectory-accounting-final-validation-20260911/` and is bound
to the exact 364-file Go/module manifest
`5bc5352d75ff5f5bc5e8c1da088510199ae97595c7a8ac55039fa4103020e8e2`.
The final result/checksum files record actual completed checks. Executed process
fixtures and asserted identities do not constitute independent model acceptance,
real operator approval or live experimental results. Windows runtime remains
unverified even if its cross-build passes.

## Next required integration

1. Retain reconstructible, exact before/after source snapshots in private runtime
   storage, including failed/interrupted dirty output. Current capture retains
   identities/digests and leaves actual files untouched; it does not yet archive
   dirty file contents for later reconstruction. Unavailable/lost source stays
   unresolved rather than becoming a clean or verified attempt.
2. Bind each retained charge to a frozen request and independently selected
   verifier invocation through the instrumented runner. Have the selected agent
   actually invoke the exact helper; parent acceptance must match identity,
   process marker, request and retained execution receipt, including failures.
3. Publish complete observations durably, derive the expected patch list from
   the full charge inventory, and feed the existing AB/BA assessment/trajectory
   evaluator. A clean intervening patch and the first retained two-regression
   trigger must behave exactly as FINAL requires.
4. Enforce required review and explicit, durable recovery/disposition across
   manual/driver/resume paths without replaying charged model work, dropping
   quorum or granting completion. Fresh independent review remains required.

The six audit recommendations, unchanged full-six/packet live experiments,
participant-owned reviews/signatures and actual-delivery follow-up obligations
remain the full objective. Latest browser-verified HTML still represents bd00188;
this source checkpoint is not presented as refreshed report QA. Actual model
inventory remains 34 terminal attempts, 17 unknown costs, USD 46.1887585 known
CLI estimates and unknown total. OpenViking tools were unavailable; no shared
memory persistence is claimed.

### Validation capture anomaly

The first broad run of the corrected source exited 1 after 156.489 seconds,
but its retained shared-volume `full.log` contains 1,697 NUL bytes out of 1,830,
with no surviving failure diagnostic after the first three package lines. The
same source's first race run exited 0 after 175.973 seconds, but its retained
`race.log` contains 133 NUL bytes out of 275. Neither sparse log is used as complete
current-source evidence. The actual failed exit is retained and its cause must
not be invented. Platform/build/shared-fixture logs contain no NUL bytes.

Replacement full/race runs use `go test -json`, capture output through a pipe,
and retain native temporary log copies before publishing their complete bytes
into the shared runtime directory. Source remains unchanged at the final manifest
above. Final verification checks complete package terminal events, command exit,
no failed test events, NUL absence, and exact native/shared log digests. A later
successful run does not retroactively explain the earlier nonzero result.

The completed replacement runs passed: full suite 152.744s with all 32 expected
package terminal events; scoped race 160.730s with all six packages. Every new
top-level test has a passed event in the full run. Native and shared JSON bytes
match their retained digests and contain no NUL bytes. Vet passed in 1.100s;
Windows amd64 trajectory/app cross-builds passed; compiled shared-volume
trajectory/runner/driver/app fixtures passed in 24.074/10.900/1.852/6.234s.
The exact source set and evidence hashes are verified in final-verification.json
and checksums.json. The earlier unexplained nonzero run remains a limitation of
this validation record; these passes are not independent source acceptance.
