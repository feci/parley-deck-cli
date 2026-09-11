---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-11
status: partial
source-manifest-sha256: c3e2507ec3c76781e1606e97a466ae34b198d38b68b5fc9f615b1a52436952c8
---

# Exact active reservation continuity

An active step/fixup/cross-review group previously rechecked policy, timing and
aggregate history without retaining the exact reservation that admitted it.
Replacing that one entry with a different ID, while keeping the same count,
could therefore permit another nested operation. Changing the original reserved
amount, timestamp or accounting start could also pass. This was reproduced at
the real ledger/session boundary in fifteen cases before the fix; the failure
log is retained. These fixtures did not launch model processes.

## Change and scope

Each active session now captures an immutable reservation receipt directly from
the snapshot returned by successful original ledger publication. The receipt
binds scope, accounting start, entry identity, kind, reservation timestamp and
known/unknown reserved exposure. Its numeric value is copied rather than aliased
to a mutable snapshot. Cycle publication returns its receipt under the existing
policy guard, so a later read cannot silently adopt a replacement as the original.
The existing public Reserve behavior, caps and failure accounting are unchanged.

Before further nested work, both session kinds recheck the exact original entry.
Same-count substitution is insufficient. Missing/changed original fields refuse
without recreating a charge or rewriting any accounting. Once a continuity
refusal is observed, it remains in that session even if old bytes reappear.
Existing policy, migration, time, cancellation and ended-session guards still
apply; failed publication remains spent where publication occurred and grants
no work.

Later separately charged operations, valid operator policy extensions, terminal
settlement and explicit unknown-cost reconciliation remain compatible with an
unchanged original reservation. A live group at its inclusive action cap may
finish under its original charge, subject to its existing lifetime clock.
Settlements remain observations rather than a replacement reservation.

## Validation and retained evidence

The exact set and hashes of all 351 Go/module files match manifest
c3e2507ec3c76781e1606e97a466ae34b198d38b68b5fc9f615b1a52436952c8.
Full Go suite PASS (146.656s). Budget/evidence/app/driver/runner race PASS
(165.232s). Scoped vet PASS (1.256s). Windows amd64 app cross-build PASS;
Windows runtime remains unverified. A compiled budget test binary passed the
new reservation cases on the isolated shared volume (5.721s).

The fifteen original negatives cover three accounting kinds and five structurally
valid ledger mutations: substitute the reservation ID, change reservation time,
turn reserved exposure unknown, change its amount, and change accounting start.
Every case initially granted the nested boundary; every corrected case refuses
and leaves the ledger bytes unchanged. Restoring original bytes cannot revive
the same refused session. Positive cases fill the actual five-action limit with
other reservations, add settlement and continue the original group without any
new charge. A separate real Store settlement/reconciliation case checks that
unknown observations remain compatible and the receipt cannot alias changed
input snapshot data. Wrong-kind and missing-entry receipt requests refuse.

Existing nested-session, ended-context, policy-extension, publication-failure,
clock, migration, driver and runner coverage also passed. Current focused logs,
the retained pre-fix failure, full/race/platform logs and checksums are under
.parley-runtime/reservation-continuity-validation-20260911 and its referenced
sibling logs. No model or real project/operator accounting action ran.

## Limits and remaining obligations

This receipt lives with the synchronous session. It is not a durable semantic
action identifier, a cross-process resume token, permission to retry ambiguous
execution, or human authentication. A consistent malicious same-UID rewrite of
all state is outside the cooperative boundary. No check can retroactively stop
work already admitted before a later corruption is observed. Corrupt history
remains an explicit refusal, not an implicit recovery decision.

Durable semantic action identity/replay, safe guard/lock-origin migration, the
opt-in independently confirmed two-patch regression trajectory, fresh independent
acceptance and all original live experiment/review gates remain open. Historical
quorum, the full-six pilot and unresolved operator/review decisions remain intact.
The model inventory remains 34 terminal attempts, 17 unknown costs, USD 46.1887585
known CLI estimates and unknown total.

The most recent verified HTML report remains the prior bd0018865946764cd65cd5cde7860bb5bfd480c9
checkpoint (629641 bytes, SHA256
06ed736c5ef7379adf1ce3c3d1adc51f5afbb34ef6869b70798f8ad28b87a0be).
Its ten report tests and 56 browser checks do not certify this new source or a
future populated experiment report. Presentation refresh remains pending while
implementation continues. No final merge, release, deployment, global install or
immutable-core publication occurred.
