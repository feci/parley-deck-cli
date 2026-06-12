---
agent: hermes
idea: meta-protocol-change-review-gate-honesty
round: 1
date: 2026-06-12
---

## P6 collision with Phase 7

Phase 7 triage records "dismissed" findings. A later reviewer must be able to state whether it concurs with a dismissal made after its own review. The "weigh openly" rule already satisfies this: the disposition travels in the brief, the reviewer records concurrence or dissent in its round file, and the operator's final call is the only thing that actually closes the item. No protocol change needed beyond the normative sentence.

## P7 — strict_gate definition of "full scope"

"Full scope" must mean the complete implementation diff at the time of the fresh review pass (all files changed since the design FINAL plus any fix-up commits up to that point). A verification pass on a subset is never a closing pass. The flag belongs in 00-prompt.md frontmatter and is immutable once set for an idea (removal would require a separate meta-protocol idea). MaxFixupCycles still bounds the loop; strict_gate simply changes the exit condition from "zero agreed fixes" to "zero findings of any severity on a fresh full-scope pass".

## Consult-standing sentence

The sentence "the protocol idea (sibling) defines its standing" belongs in this amendment under a new short subsection "Consult standing" in Phase 8 (or a new §12). It is not an implementation detail.

## Precise fixes

Add after the Phase 7 paragraph:

> Review briefs may contain dispositions of known findings (rebuttals, accepted trade-offs, recorded artifacts). The reviewer must state in its review whether it concurs with each disposition. No brief may instruct a reviewer to suppress or ignore any class of finding.

Add a new subsection under Phase 8:

> ### Strict gate
> When `strict_gate: true` is present in 00-prompt.md, Phase 8 closes only on a fresh full-scope review pass whose verdict is zero findings of any severity. Fix-up verification passes converge the gate but never close it. Operator rulings remain the only other closure mechanism. Trajectory guidance (converging / churning / blocked) applies to all gates.

## Position

Both points are sound. The two normative paragraphs above plus the strict_gate subsection are the minimal sufficient amendment. No other COOPERATION.md sections require change.