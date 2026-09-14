---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
source-base: fbeac7910be3535dfbcfbf411535ffb359efab46
status: integrated-validated-awaiting-independent-review
---

# Refuse inconsistent helper scope before consuming verification

Claude N3 identified a gap between the live helper membership and the retained
activation membership. The lower-level ticket authority checked the activation
and captured archives, while the app's early check compared the helper request
with the live prompt that supplied it. Their disagreement could be detected only
after ticket consumption and a verifier run, leaving no supported retry.

CheckHelperScope now reuses the strict reconciliation scope decoder and original
activation authority. It first binds the complete supplied captured request to
the actual charged attempt under the shared state/ledger/archive guard, then
compares the helper's ordered membership to activation and its live membership
and named criteria to that same request. The app invokes it before ticket
preparation, before helper execution and before accepting the parent result.
The previous app-only check is replaced by this shared check, rather than adding
a separate interpretation of YAML or criterion scope.

This is read-only early refusal. It creates no ticket, invocation, receipt,
resolution, schema migration or acceptance. Existing final reconciliation and
execution-authority checks remain. It does not fence an external edit after
the check or recover a previously consumed ticket. It adds a full guarded read;
the history-size/guard-contention concern in N2 remains open.

## Executed focused evidence

Three top-level tests passed in 11.538s over 421 unchanged Go/module files:

- An actual charged patch followed by widened/reordered live membership refuses
  before any ticket, verifier invocation or criterion execution. State and ledger
  remain byte-equivalent. Restoring the exact original prompt allows the same
  original charged patch to receive real independent AB/BA criterion execution.
- A three-member activation prompt with a YAML trailing comment exercises the
  existing live-parser/strict-parser disagreement. The mismatch refuses before
  any one-time verification authority is consumed. This does not rewrite the
  general workspace parser or claim that every YAML syntax is now supported.
- A fixture-issued original ticket/launch plus a coordinated changed live prompt
  and matching helper request refuses at the actual helper entrypoint before
  creating a claim, preparation, process, step or receipt record.
- Direct checks refuse widened/reordered/duplicate members, a changed original
  invocation, changed criterion command, unsupported request version and missing
  trajectory authority; the original request passes.

The helper-entrypoint fixture supplies test-owned environment attribution and
is not authentication of a human or real model. The production-path verifier is
a local fixture CLI running real criterion commands, not a participant signoff.

Evidence: .parley-runtime/helper-scope-development-20260914/.

## Frozen validation and controls

The 421-file Go/module manifest is
7474a4fdc3befca4db92b02a1a71614a9facc34b0a3202ce69497d9dd125ca51.
The focused, full, race and compiled shared-volume selections used identical
source. Full Go tests passed in 416.466s: 31 passing package terminals and the
CLI package's no-test-files skip. Six-package race passed in 465.572s; vet passed.
All three new top-level tests passed in both full and race runs.

Compiled shared-volume trajectory selection passed in 45.755s; app selection
passed in 66.566s. The app selection includes the valid production helper flow.
Windows CLI, trajectory tests and app tests cross-built successfully; PE amd64
headers were checked. Windows execution remains unverified.

Three source-overlay controls failed at their intended assertions: removing
activation-membership comparison (4.303s), original captured-request equality
(1.829s), or live scope/criteria validation (1.830s). The first is detected by
the actual app path's assertion that refusal precedes consuming verification.
No accepted test was rerun merely to correct final-verifier package bookkeeping.

Final evidence: .parley-runtime/helper-scope-final-validation-20260914/,
including final-verification.json. Native/shared terminal log hashes match,
the exact five changed Go paths match the frozen development source, and the
three participant reviews and historical September-5 HTML remain unchanged.
Commands were captured through native pipes and their logs copied after exit.

## Remaining obligations

This is an implemented response for independent review of N3, not a withdrawal
of Claude's finding. N1 precharged driver/CLI refusal prevention, N2 contention,
historical/abnormal terminal and helper/workflow recovery, full-scope participant
acceptance/signatures and the remaining experiment/delivery gates stay open.
The budget-refused-before-start candidate remains a separate native overlay.
