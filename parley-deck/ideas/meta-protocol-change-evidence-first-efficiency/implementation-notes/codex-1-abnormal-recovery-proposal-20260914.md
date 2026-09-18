---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
status: proposal-not-implemented
source-base: 83d0475fd81acc7a33a66b2a2f1588f1b870aa04
---

# Remaining F1 recovery: distinguish an observed failure from permission to resume

Claude F1 identifies a real operator gap: an unchanged attempt that timed out,
was cancelled/signalled, hit a classified provider error or was refused after a
charge cannot use either unchanged reconciliation or changed-patch verification.
Its charge remains spent and subsequent fixups/continuation remain blocked.
The prospective refusal and terminal-retention changes prevent specific new
losses; they do not resolve those historical or abnormal attempts.

## Candidate bounded control

A new attended abnormal-observation control could retain an explicitly
inconclusive assessment while binding the original policy, charged attempt,
ordered activation quorum, unchanged source archives and exact available
request/start/terminal bytes. It would need a distinct typed evidence variant;
it must never masquerade as an ordinarily exited unchanged process or a verifier.
A failure assessment preserves regression history and all unresolved criteria.
No refund, clean outcome, helper receipt or new fixup is implicit.

The apply must use a preview tied to the complete current state and exact
lifecycle evidence, an explicit operator decision ID/reason, and exact replay
semantics. Missing or contradictory expected records remain unresolved. A
recorded unstarted refusal must be distinguished from an absent start record;
a terminal's own recorded start/PID/status fields and original request must agree.
Changed material source continues to require independent patch verification.

Recording such an observation must be separate from permitting continuation.
An operator's claim that writers have been stopped is testimony, not process
proof. If the protocol accepts that attended authority for continuation, it must
record its author, scope, exact original invocation, remaining custody unknowns
and affected checkout/workspaces explicitly. It cannot be inferred from a generic
`--yes`, an elapsed timeout, an empty PID search or a vanished process leader.
An unattended resume must still refuse without the required authority.

## Custody and helper boundary still to settle

Existing helper claim v2 retains PID and time, not full boot/start/group identity.
Registered criterion groups have stronger identities and a stop fence, but a
leader's absence does not prove its group or escaped descendants inactive.
Completed step journals do not supply a missing helper receipt. A stopped or
partial helper cannot be relabeled a successful verifier.

A future helper retry would require its own invocation, genuinely spent launch
reservation and preserved original ticket/attempt linkage. Before executing it,
its control must address what the original model/helper/criterion processes can
still affect. New process identity and containment evidence can improve future
custody; it cannot be retroactively manufactured for historical records. A
separate operator-containment decision may be necessary for irreducible old
unknowns. This proposal does not choose that authority by silence.

## Review questions and falsifiers

1. Is an attended inconclusive observation useful independently of continuation,
   or should the UI present one staged preview with separate observation and
   custody/continuation authorities?
2. Which exact custody evidence or explicitly scoped human assertion may permit
   continuation after a fully retained timeout/provider failure? Which historical
   missing-start/helper cases must remain blocked?
3. What versioned evidence variant preserves old-reader fail-closed behavior,
   exact original hashes, regression streaks and all current closure gates?
4. Which production interruptions must refute the proposal: unstarted refusal,
   signal/timeout with live descendants, PID reuse, missing/changed record,
   stopped helper with completed steps but no receipt, concurrent recovery,
   changed preview and exhausted remaining budget?

The source seams are unchanged.go (lifecycle qualification), continuation.go
(inconclusive/review gates), reconcile.go (typed evidence and rederivation),
state.go (complete ledger/transition authority), verification.go (one-use helper
journal and stop), and procctl (attributed group signalling). These are locations
for review, not ownership claims to edit them. No new recovery command, schema,
execution permission, participant waiver or pilot amendment is implemented here.
