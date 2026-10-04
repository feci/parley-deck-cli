---
from: codex-1
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: implementation
blocking: yes
date: 2026-10-04
---

## Question

Two owner rulings are needed before implementation can continue. This is a scope/reset
interpretation gate, **not** the attended close and **not** a request to release.

**1. AC2 / zcode support.** Do you authorize additional, separately reviewed zcode
terminal-error channel work so AC2 remains required, or explicitly defer/amend AC2?

- **A — recommended by codex-1:** keep AC2 and both stages required; extend scope to a
  zcode-side channel/integration that preserves the original provider failure and every
  decisive reset value, bound to the root invocation. This can require a zcode-side
  change or a separately reviewed transport. Merely switching the installed CLI to
  app-server does not work. Native positive and adversarial fixtures, including gateway
  pass-through provenance where relevant, must pass before any automatic eligibility.
- **B:** explicitly defer AC2 to a linked follow-up whose precondition is that complete
  zcode-side channel. The current implementation has **no supported adapter**, so this
  option means no automatic quota exclusion in the current adapter set. Both lifecycle
  stages remain required unless you give a separate, explicit stage-1-only amendment.
- **C:** amend AC2's adapter or positive message explicitly. No alternative adapter is
  claimed provenance-verified today; the replacement still needs native evidence and
  independent review. Specify the intended replacement or authorize a proposal first.

**2. Display time without a timezone.** Do you adopt claude-1's bounded interpretation
below (recommended), or retain the strict rule and amend/defer AC2's positive example?

A display clock without a timezone is non-decisive only when ALL of these hold:

1. The **same native terminal record** has at least one complete machine reset value:
   RFC3339 with `Z`/offset, or numeric-seconds `retry_after` / `Retry-After`.
2. All decisive values agree within the existing fixed one-second tolerance, including
   any `reset after` duration measured from observation.
3. The display clock matches the machine instant at a real UTC offset from −12:00 to
   +14:00 in 15-minute steps (the recorded example corresponds to +08:00).
4. The raw display string is kept in the evidence and exclusion marker.

Missing values (including named reset keys whose values were discarded), contradictory
values, an irreconcilable display clock, or a display clock as the only reset still
**gate**. A stated reset never becomes `unknown` or enters the no-reset branch. The
current JSONL channel remains unsupported even if this interpretation is approved.

## Context

The owner authorized both stages and release after independent review and an attended
close. The controlling brief says **"Implement exactly FINAL"** and **"Its verdicts
bind"** for the separate claude-1 reviewer. FINAL AC2 requires a provenance-verified
zcode recognizer for the recorded 429 Weekly/Monthly Limit Exhausted with its 49-hour
reset. FINAL §4.5 separately says an unparseable reset gates.

The separate claude-1 reviews are canonical:

- `../ideas/meta-protocol-change-quota-auto-exclude/review/round-01/claude-1.md` (40e44d4).
- `../ideas/meta-protocol-change-quota-auto-exclude/review/round-02/claude-1.md` (3ea19c6).

Round 02 exited 0 after 1729.8 seconds; `parley wait` validates its 26442-byte artifact.
It confirms and extends both MAJOR findings. In the installed zcode:

- stderr does not bind quota text to the terminal failure;
- the JSONL business-error path discards the reset values;
- the recorded HTTP-429 path loses the exhaustion message itself;
- app-server retains less data, and traceId is shared with subagents;
- interpreting a timezone-free clock as supplemental display text is an owner policy
  decision under this FINAL, not a routine parser fix.

The implementation is saved, not complete. CLI HEAD before this gate is 3ea19c6; skill
HEAD is c0d7f58. Partial stage-1 code is in 3aa05cf. All adapters remain diagnostic-only.
Stage 2 has not started. Host build, vet, gofmt, focused tests and
`go test ./... -timeout 45m` pass; the full skill suite passed 399 Node and 54 Python
checks. These are producer checks, not independent AC acceptance. The literal AC1
whole-file comparison also remains an open, documented bootstrap-zone discrepancy.
There are no review-consensus signoffs, and no merge or release has occurred.

The review's two MINOR findings remain queued: extend the support table and stop the
reset-time regex at the opening parenthesis. They do not resolve the MAJOR scope tension.

## What I need from you

Reply with the choice for **1** (A, B, or C) and **2** (adopt the bounded interpretation,
or keep the strict rule), or provide a different explicit direction. The recommended
pair is **1A + adopt the bounded interpretation**. It preserves the requested positive
path and both delivery stages while keeping missing or contradictory evidence blocked.

I will quote the answer into the authoritative continuation record and next review
context. FINAL remains immutable; any required normative change follows §7. A linked
follow-up and explicit deviation are required if AC2 is deferred or replaced.

Your answer here does not waive full-scope claude-1 review, both review-consensus signoffs,
current-tree criterion evidence, or your later attended-close confirmation. No further
implementation or agent launch is authorized by silence.

Authority for pausing: the Parley Deck skill's Escalation To User rule says, "If
`blocking: yes`, pause the escalating agent's work for that idea." The live protocol
§4 escalation rule applies the same gate. This note records it explicitly.
