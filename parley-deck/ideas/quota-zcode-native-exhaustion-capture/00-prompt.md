---
idea: quota-zcode-native-exhaustion-capture
author: codex-1
facilitator: claude-1
created: 2026-10-06
track: deliberation
participants: [codex-1, kimi-1, zcode-1]
auto_implement: false
status: open
---

## Problem / idea

Follow-up required by the owner's round05-answer Q2 for
`meta-protocol-change-quota-auto-exclude`. That release accepts native-positive AC2 as
unmet: zcode auto-exclusion may not fire on real native output; unrecognized failures
fall back to the owner-confirmed path. R5-MAJOR-2 is accepted and deferred, not fixed.

## Precondition and scope

Grammar work waits for a complete scrubbed capture of the next real native zcode
exhaustion, with adapter/version, invocation identity, start and terminal-receipt times,
nonzero exit, output/artifact facts, and hashes/provenance. Do not spend quota to
manufacture an error, run a provider probe now, or weaken any predicate without evidence.
Source-derived SDK fixtures and partial native excerpts cannot satisfy this precondition.

First assess and implement, under its own reviewed plan, private scrubbed failed-zcode
stderr retention where needed so the next ordinary failure becomes usable evidence.
Specify private location/permissions, credential and payload scrubbing, truncation
signaling, bounds, retention/cleanup, and opt-in/owner controls before treating a saved
stream as safe complete evidence. This capture facility is explicitly deferred from the
narrow cycle-3 repair. Existing logs are not certified to meet that contract.

After the precondition: compare real framing to the strict recognizer; propose only
observed benign grammar extensions (including actual stack frames and inspect
placeholders if observed). Preserve every-record agreement, machine reset, terminal
60-minute threshold, nonzero exit, no valid artifact/later success, failure-line framing,
quoted/tool spoofing negatives, survivor floor, role guards and owner notice. Any change
to the owner's evidence rule needs an explicit owner decision.

## Acceptance direction

A reviewer must reproduce the real positive from the scrubbed capture and all existing
adversarial negatives, distinguishing native, source-derived and synthetic evidence.
The follow-up resolves the accepted limitation only after its own design, implementation,
independent review and release. The parent release never claims this work shipped.

## Sources and status

- `../meta-protocol-change-quota-auto-exclude/review/round-05/claude-1.md`, R5-MAJOR-2.
- `../../inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_round05-answer.md`, Q2.
- Parent `IMPLEMENTATION.md`, the AC2 waiver and immutable earlier evidence.

Backlog only: no readiness probe, participant invocation or active deliberation has been
started. Listed facilitator/quorum are the owner's current defaults for a future new run;
readiness and effective config are checked at that launch. They do not alter the parent
idea's two-agent quorum or authorize participation in it.
