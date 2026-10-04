---
idea: meta-protocol-change-quota-auto-exclude
review-cycle: 1
outstanding_agreed_fixes: 9
blocked: false
drafted-by: codex-1
date: 2026-10-04
reviewed-commit: 78ac5367cff0d9a2f9ae68448eb670489ec8ed29
skill-commit: dc85b5362b4872cc719a8ed745ba2a6cc89cfad6
---

## Scope and review basis

This is the first full implementation review cycle. Review rounds 01/02 were focused provenance
feedback, not fix-up cycles. The separate claude-1 full review is round-03. Its findings bind; codex-1
accepts the following work as implementer, without grading its own implementation. Both-stage delivery,
the owner's bounded stderr/reset ruling and frozen FINAL remain authoritative. This consensus authorizes
fix-up cycle 1, not close, merge or release.

## Agreed fixes

- **G1 — shared-volume durability** (round-03/claude-1 CRITICAL-1): route every new kickoff, projection,
  creation and notice sync through the established fsutil helper; add the appropriate failure seam and
  prove new-idea creation on this actual AppleVirtIOFS workspace using stubbed dispatch, with the knob
  on and off. Keep checked durability and fail-closed behavior.
- **G2 — authorized membership and scope changes** (CRITICAL-2, MAJOR-1): implement a durable
  owner-confirmed revision path and CLI entry point, bound to the archived owner-answer evidence. Preserve
  ordinary owner-confirmed exclusion/catch-up behavior with the knob off, the frozen recorded policy,
  and immutable membership history. Support same-idea re-inclusion of known identities, authorized
  catch-up joins under existing protocol rules, and explicit scope widening; binary upgrades/resumes
  must never widen scope. A re-included author's later recorded withdrawal must be able to dispose of
  its retained veto. Reconcile every consumer/projection without silently undoing an authorized change.
  This is an implementation choice within FINAL's required return/scope paths, not new owner policy.
- **G3 — actual cross-process serialization** (MAJOR-2, MINOR-1, MINOR-4): replace reliance on flock with
  the repository's exclusive PID/token lease primitive (factored for reuse if needed), carrying idea/run
  identity, liveness and conservative stale handling. Retain lifetime ownership, nested-context checks,
  projection serialization and the off/legacy/kickoff-only scope. Use bounded projection-lock waiting;
  no quota retry worker. Keep transient locks out of tracked artifacts. Verify competing run IDs and
  projections on this AppleVirtIOFS volume as well as local disk, including stale/partial/release races.
- **G4 — bounded stderr framing** (MAJOR-3 and AC2 limitation): reject raw quotations/free prose and every
  mixed JS error-class/Unhandled form demonstrated by the reviewer. Validate any retry wrapper against
  the provider error. Parse complete SDK framing, retaining all owner-specified gates and accepted
  subagent-attribution limitation. Recover a complete historical native capture if available; record
  its provenance and scrub it without losing decisive fields. Do not call a reconstruction a native
  capture, silently weaken AC2, launch excluded providers, or widen the owner exception. Add every
  demonstrated adversarial case and native-framing coverage; unresolved evidence remains explicit.
- **G5 — signoff handoff context** (MAJOR-4): thread the active leased context through manual/interactive
  signoff handoffs and execute a mid-idea-scoped handoff regression, including competing-run refusal.
- **G6 — retained-obligation authority** (MINOR-2): bind owner rulings to an actual archived user-answer
  path and its stable content, with idea/authority validation. Bind author withdrawal to the authorized
  re-inclusion and author-owned later artifact; independent dispositions still need independent evidence.
  Test missing, unrelated, changed, self-authored and fabricated authority references. Preserve the
  distinction between structural attribution checks and human truth/identity verification.
- **G7 — projection receipt recovery** (MINOR-5): provide a documented checked recovery path for truncated
  or contradictory mutable receipts, based only on validated immutable history and repaired projections.
  Never interpret a bad receipt as applied, erase immutable history or permit actions before reconciliation.
  Verify recovery and repeated recovery without duplicate notices/evaluations.
- **G8 — invocation and lifetime binding** (MINOR-3 and Open question 3): inspect and either correct or
  provide executed refutations for the TUI post-release integrity-note and pipeline lease-lifetime leads.
  For launches that target a canonical idea artifact, ensure missing explicit --idea cannot bypass idea
  membership/lease gates. Generic unbound agent executions remain outside an idea's quorum; do not add a
  global singleton or speculative cross-idea restrictions. Record this scope and test the relevant paths.
- **G9 — accurate diagnostics** (NIT-2, NIT-3): name applicable remaining gates in notices, and align bare-503
  recognition in runner/preflight while excluding irrelevant line numbers or byte counts. Preserve the
  actual bare-503 provider gate and add positive/negative regression cases.

## Deferred follow-ups

- Owner D6 / driver gap 11: permanent legacy run-accounting repair remains a separate idea (TBD), to be
  started after this idea; no migration, repair, worktree declaration or pruning here.
- CLI Windows runtime validation remains explicitly unverified until executed on Windows. Compile-only
  evidence is not runtime evidence; the CLI winget release remains held while Windows is experimental.
- Unsupported-adapter native provenance and the Claude JSON-output switch remain the frozen FINAL's
  scoped follow-ups; they are not excuses to loosen this change's zcode gates.

## Dismissed findings

These dispositions require the reviewer's own signoff; codex-1 supplies no independent acceptance verdict.

- Round-02 MAJOR-1/2 are closed by the explicit owner ruling, as claude-1 states in round-03. The new
  framing defects are G4, not a reopening of that owner decision.
- Round-02 MINOR-1/2 are resolved per claude-1's round-03 evaluation (support limits and reset regex).
- Round-03 NIT-1: retain completed-drafter protection. FINAL §7 protects a candidate who has started a
  canonical consensus/FINAL draft and freezes closed artifacts; it does not authorize removing that
  protection after completion. This conservative behavior follows the frozen scope.
- Round-03 NIT-4: keep literal AC1 skill/deck byte equality and the existing bootstrap replacement
  instruction. No drift assertion or bootstrap generic-zone test is weakened.

## Coverage & blind spots

Only claude-1 is an independent reviewer; codex-1 organizes and implements. Reviewer round-03 executed
shared-volume creation/locking counterexamples and 54 recognizer probes; its findings outrank green
producer checks. It executed affected packages/build/vet/drift/packet/cmp; full HOST and full skill suites
are producer evidence. No Windows runtime evidence exists. Incident 1 remains a labeled replay, and the
incident-2 retained four-line excerpt does not establish whole-capture framing. Each gap must remain
visible in the next review. Re-review must use the complete product diff and full emitted protocol body;
round-03 explicitly read selected body sections and spot-checked some callers, so its coverage is limited
as disclosed. The next review must complete that required context/coverage before acceptance.

## Signoffs

<!-- Each participant appends only its own block. Signoff accepts the fix list, not the current code. -->

### Signoff: codex-1 — 2026-10-04
Status: ✅ ACCEPT
Notes: I accept G1–G9 as the implementer fix list for cycle 1, including the shared-volume reproductions and owner-confirmed return/scope obligations. This is agreement to correct the code, not an independent verdict on my implementation. The stated limited review coverage and capture provenance remain explicit. Both stages, separate re-review, both clean-consensus signoffs and a new attended owner close remain required.
