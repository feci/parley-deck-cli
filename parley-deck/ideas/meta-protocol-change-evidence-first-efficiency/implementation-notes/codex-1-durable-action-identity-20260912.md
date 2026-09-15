---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-12
status: implemented-pending-independent-review
base-commit: b3cec48b4f222bce64aa49abba701ff7f72cddb9
---

# Durable action accounting identity

The runtime now persists logical operation/input identity separately from each
spent attempt. A typed action descriptor is published atomically with its charge
under the existing ledger lock. It binds scope, original entry key, declared
input, effective limits and the immutable reservation (epoch, kind, time and
reserved amount). The ledger envelope remains schema 1 with an additive typed
field; older strict readers refuse the new field. Legacy unbound entries retain
their original absence of semantic evidence.

Runner round/implementation/fixup/review-consensus calls bind a declared request
and bounded protocol-file hashes. Driver Advance binds its reconstructed cursor
and runtime controls. Nested groups retain the outer operation, manual/lower-level
launch metadata has an explicit weaker basis, and otherwise-direct sessions
identify only their policy. These are declared recipe/cursor fingerprints, not
whole-source-tree proof or authenticated participant identity.

The read-only `budget action inspect|replay` command exposes deterministic
receipts and exact original identity comparison. Its output grants no execution
and establishes no operation terminal. Original receipts survive later settlement
and conservative cost reconciliation. Duplicate IDs refuse; changed bound input
or accounting adds an explicit conflict. New attempts remain separately charged.

## Implementer-discovered counterexample and correction

The initial implementation bound input and limits but did not bind the original
reservation fields into the replay identity. The retained reproducer failed for
all five mutations: kind, reservation timestamp, ledger epoch, reserved amount,
and known-to-unknown reserve. It therefore could return a changed receipt under
the old expected identity. The correction computes `reservation_sha256` inside
the reservation transaction and validates it on read. A structurally consistent
replacement with a recomputed commitment still has a different identity, which
the original replay request rejects. Existing live-receipt mutation fixtures
recompute internal hashes to stay structurally valid and exercise the intended
continuity guard rather than only malformed-JSON rejection.

Evidence is retained under `.parley-runtime/action-identity-validation-20260912/`:
`receipt-counterexample.log` has all five intended failures; `receipt-fixed.log`
and `context-continuity.log` have the corresponding passing checks. The original
fixture failures in `focused-initial.log` are preserved. These targeted results
are not final full-source or independent participant acceptance.

The expanded real runner fixture then exposed a second integration omission:
RunFixup persisted its invocation and charge but did not return InvocationID in
its Result. Its existing launch observer now propagates that ID, with attempt
ordinal 1. The test requires the launch ledger entry key to match the actual
returned invocation and retains the failed child's settlement. The failed first
final-validation attempt remains in
`.parley-runtime/action-identity-final-validation-20260912/focused.jsonl`.

## Required remaining work

Accepted validation uses source manifest
`eadb9d9c91d4038e6732c91d9134e2a745c712da03efee12e31ac151e84b3f1e`
(392 Go/module files) and
`.parley-runtime/action-identity-final-validation-v2-20260912/`.
Native copies: `/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-action-final-v2-7gzu_ry7`.
All validation processes are terminal.

- Focused checks: PASS, 4.879s.
- Full Go suite: PASS, 220.143s, all 32 package terminals accounted for.
- Six-package race: PASS, 261.187s.
- Vet and Windows budget/app/runner/driver cross-builds: PASS. Windows runtime is
  unverified.
- Compiled shared-volume budget/app/runner/driver checks: PASS,
  2.834s / 0.432s / 0.705s / 0.542s.
- Five removed-protection overlays fail at their intended assertions: replay
  identity comparison, reservation commitment, duplicate-input conflict,
  session-input continuity and live receipt identity.
- The actual compiled CLI inspects and replays the original receipt, refuses a
  different identity and an activation flag, preserves ledger bytes and refuses
  a missing ledger without creating it.

The final verifier checks all thirteen new top-level test functions (twelve
substantive plus their separate-process harness), logs against native byte
copies and checksums, unchanged source hashes, PE amd64 build identities and the
unchanged historical report. Implementer-authored automated evidence does not
supply independent current-source participant acceptance.

This slice does not complete semantic workflow-effect recovery. Missing parent
results, unfinished helper tickets, orphan reservations, ambiguous process
liveness and unchanged-source recovery require their own evidence and controls.
New step/cycle sessions still create distinct attempt IDs. Receipt replay cannot
be substituted for those obligations, independent current-source acceptance,
actual packet/pilot experiments, participant-owned final signoffs, the final
HTML or actual-delivery-based follow-ups. No real-model call or operator migration
was performed for this slice.
