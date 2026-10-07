---
from: codex-1
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: review-round-08
blocking: yes
date: 2026-10-08
---

## Question

Authorize accepting and deferring the three remaining findings as explicit release
limitations, using the exact disclosure proposal below and a linked follow-up,
then obtaining both final review-consensus signoffs and releasing? I recommend
this because the sole independent reviewer reports no CRITICAL or MAJOR and all
five cycle-5 fixes are verified. The alternative is to keep the change parked.
I have not opened a sixth code fix-up cycle or inferred a waiver.

## Context

The independent full-scope round-08 completed at 2026-10-07T23:48:59Z, exit 0,
1393.6s. It used plain `claude-opus-5-5[1m]`, max effort, after two provider
relaunches separated by the required 900 seconds. No model/provider/quorum changed.
`parley wait --for review --timeout 1s --json` exits 0, 1/1 filed-and-valid.

- Frozen CLI: `2705a1e850132f74aed2dbb87149df5491bfe94b`.
- Frozen skill: `99b3f3f9fee161e8e61e585ad6c8e5b1bbbd4d2b`.
- Review: [round-08/claude-1.md](../ideas/meta-protocol-change-quota-auto-exclude/review/round-08/claude-1.md),
  SHA256 `6ea871262a87731b2c5273b28128b851c7a0ad797abaf26601e4c845d404ce06`.
- Result: **0 CRITICAL, 0 MAJOR, 1 MINOR, 2 NIT**. Every G21–G25 fix and signed
  reservation is verified. No new finding needs a scope or FINAL change.

| Finding | Effect and evidence |
| --- | --- |
| R8-MINOR-1 | A missing/read-only kickoff inbox loses a floor/role blocking escalation and its decision details. The command still fails closed and creates no idea. Executed on shared and local filesystems; pre-existing stage-1 code missed by prior reviews. |
| R8-NIT-1 | A crash between the kickoff manifest write and notice publication permanently loses that notice; other exclusion surfaces remain. Source evidence only, no injected crash; pre-existing. |
| R8-NIT-2 | On an aliased deck with the policy off, a plain participants edit blocks every signer until a physical path is restored. The behavior matches the signed alias-refusal plan, but the generic disable-policy advice and disclosure are imprecise. Executed on both filesystems. |

The concrete proposed release text and follow-up scope are in
[codex-1-round08-disclosure-proposal.md](../ideas/meta-protocol-change-quota-auto-exclude/source-context/codex-1-round08-disclosure-proposal.md).
Accepting this option explicitly accepts the AC5 escalation-detail exception and
AC15 kickoff crash-window exception for this release; none becomes a PASS without
its caveat, and none is called fixed. The alias diagnostic caveat remains visible.

The review independently ran the full Go suite (616.8s, 34 packages), build, vet,
race (411 PASS), shared/local focused runs (244 PASS each), all 100 changed Go files'
format checks, packets/drift and skill tests (399 Node, 54 Python, six manifests).
All those checks passed. Windows CI actually fails and is disclosed; compilation
is not runtime verification. AC2 stays NOT MET / expressly owner-waived. The
current AC1–AC21 table is in IMPLEMENTATION.md under Validation evidence.

## Why this is blocking

This is the residual-findings boundary in the owner's newest
[finish-now direction](user-to-codex-1_meta-protocol-change-quota-auto-exclude_finish-now.md),
point 4, verbatim:

> **Fix-up.** Do cycle 4 now: the plan signoff, the implementation and the full re-review. If the cycle-4
> re-review finds new findings, fix them narrowly in cycle 5 without asking. That is the last cycle the
> protocol's deliberation cap allows. Escalate only if a CRITICAL needs a change of scope or FINAL, or
> if findings remain after cycle 5.

Point 5's pre-confirmed close is preserved. This is not a new attended-close
request: the separate point-4 condition is now true, and the final consensus
signoffs also do not yet exist. The signed file currently approves the cycle-5
plan, not the final code. IMPLEMENTATION remains `fix-up-cycle-5`; no main merge,
version bump, tag, publication, installation or actual core staging has occurred.

## What I need from you

Choose **accept/defer with the proposed disclosures and linked follow-up, then
final signoffs and release** (recommended), or **keep the change parked**.
The recommendation adds no sixth code fix-up cycle. After acceptance, claude-1
still owns its final consensus signoff and independent channel verification;
the existing two owner-only npm/core publication commands would go in the one
final released note. No such released note is written while the release is blocked.
