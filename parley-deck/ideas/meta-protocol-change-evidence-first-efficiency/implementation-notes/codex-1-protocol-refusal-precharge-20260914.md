---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
source-base: db50dad638fc227c6d8d41e9de3ebef057c3d65c
status: implemented-awaiting-independent-review
---

# Record known protocol refusal before fresh reservations

The protocol renderer can refuse before any model process is permitted to start.
Previously beginProtocolLaunch still called beginLaunch first, which could spend
a new fixup/step/launch reservation or consume a prepared verification ticket.
The later refusal then left charged trajectory work unresolved (one F1 class).

The launch boundary now records that already-known refusal request and its actual
unstarted terminal before preparing a cycle or making any fresh reservation.
Request metadata/observation setup is shared with ordinary beginLaunch; normal
launches retain their existing cycle, ticket and budget order. The refused context
is never emitted. Recording failure still returns failure and grants no execution.

Existing precharges stay spent and byte-identical, including absent launch and
terminal fields. They remain unresolved; no refund, historical terminal,
after-source, helper receipt or independent pass is invented. A previously
unconsumed verifier ticket remains available after restoring valid authority;
an already consumed or incomplete ticket is not recovered by this change.

## Verification scope

Three new tests exercise the real measured boundary: a fresh refusal spends no
cycle/step/launch budget and later ordinary launch does; an existing charge stays
unchanged and cannot satisfy completion; and a prepared verifier ticket remains
unconsumed until a valid launch. The fresh/precharged fixtures alter only ignored
protocol-cache bytes, preserving the material baseline so a later source mismatch
cannot mask removal of this gate. The verifier fixture restores exact source
bytes and original mode.

The native-only prototype passed the strengthened focused set in 21.723s, with
three intended removed-gate failures. Its earlier fixtures failed by expecting
an unreserved launch ledger to exist and by restoring the wrong source mode.
Those failures remain under .parley-runtime/protocol-refusal-prototype-20260914/;
the fixtures were corrected without weakening production guards. Integrated focused tests passed in 19.951s, including all three new tests
and five existing launch/protocol tests. All three regenerated removed-gate
controls failed at their intended assertions (2.337s, 1.669s, 1.742s). Frozen
full/race/vet/cross-build/shared validation passed as recorded below. Prototype evidence does not establish integrated-source acceptance.

This is prospective prevention of one known pre-spawn refusal class only.
Timeout/cancellation/watchdog/signal/provider failure after charging, historical
precharges, missing terminal/source publication, consumed/incomplete helper
tickets, workflow effects and descendant custody remain open. So do F6 full
history/size evidence, current participant-owned acceptance/signatures, real
launch/concurrency/closure evidence, frozen packet/pilot experiments, final HTML
and delivery-based follow-ups. The whole audit remains incomplete.

## Frozen integrated-source validation

Manifest 8938e8612351bd7efbbf85df589d0f4fdd541cff178daf9b42f8f6dec195abc2 pins 416 Go/module files. Full Go suite
PASS 383.281s (all 32 package terminals); six-package race PASS 406.722s; vet PASS.
Windows production CLI and runner test cross-builds passed; PE amd64 headers
checked, runtime unverified. Compiled shared runner selection PASS 29.665s.
Both full/race and shared runs include all three new tests and the five surrounding
protocol/launch tests. Three intended removed-gate failures, exact source hashes,
native/shared logs and unchanged old participant review/historical HTML are checked.
Evidence: .parley-runtime/protocol-refusal-final-validation-20260914/ and
.parley-runtime/protocol-refusal-development-20260914/.

The separate current-source supporting Claude review failed with HTTP 429/session
limit, no artifact and USD 0 CLI estimate; it supplies no independent acceptance.
Live inventory now 36 terminal attempts, 17 unknown costs, USD 54.424894 known estimates;
total unknown. The original idea quorum and packet/full-six pilot remain frozen.
