---
idea: close-integrity
drafter: claude-1
implementer: claude-1
date: 2026-06-24
status: final
spawned_from: loop-engineering-research
---

## Purpose

Implementation spec for Tier 3. Hardens the auto-driver's completion decision so a clean
`outstanding_agreed_fixes == 0` is necessary but not sufficient under auto-drive.

## LE-11 — HITL-fatigue guardrails

**Files:** `internal/driver/impl.go` (`advanceReview` close branch + `ReviewStatus`);
`internal/app/driver_impl.go` (`ReviewStatus` fills `ReviewerCount`); `COOPERATION.md`
§Phase 6/8.

- `driver.ReviewStatus` gains `ReviewerCount int` (set by `driverImplOps.ReviewStatus` to
  `len(o.reviewers)`).
- In `advanceReview`, in the `OutstandingAgreedFixes == 0` close branch, **under
  `AutoImplement`** (the code-writing, higher-risk mode), before `Complete()`:
  - if `rs.Summary.Triage == TriageReserved` → escalate (ACCEPT-WITH-RESERVATIONS needs a
    human to read the reservations; do not silently auto-close).
  - if `rs.ReviewerCount < 2` → escalate (a single checker is too weak for unattended
    completion; add a reviewer or sign off manually).
- A non-`auto_implement` (design-only) idea keeps the lighter close (conditional rigor).

**Acceptance:** an `auto_implement` idea with `TriageReserved` escalates instead of
completing; with `< 2` reviewers it escalates; with `≥ 2` reviewers + `Ready` it proceeds.

## LE-7 — Goal-done gate

**Files:** `internal/driver/impl.go` (`ImplOps.GoalCheck` + `advanceReview`);
`internal/app/driver_impl.go` (`GoalCheck` via `runner.RunConsult`);
`internal/runner/consult.go` (a goal-check prompt builder); `COOPERATION.md` §Phase 4.

- New `ImplOps` method `GoalCheck(ctx) (ok bool, detail string)`.
- The app adapter runs a **fresh non-implementer** agent (the review drafter) via the
  existing `RunConsult` execution path with a goal-check prompt: it reads FINAL.md +
  IMPLEMENTATION.md and ends with a structured verdict line `GOAL-CHECK: PASS` or
  `GOAL-CHECK: FAIL` + the unmet criteria.
  - Parse: a clear `FAIL` → `(false, reasons)`; a clear `PASS` → `(true, "")`.
  - **Checker-failure is advisory (fail-open):** if the consult errors or the verdict is
    ambiguous, return `(true, "advisory: inconclusive")` and print a warning — the gate is
    defense-in-depth on top of the already-passed review consensus, NOT the sole gate, so a
    broken checker must not block a clean idea. Only a confident FAIL escalates.
- In `advanceReview`, in the close branch, **under `AutoImplement || StrictGate`**, after
  the LE-11 guards and before `Complete()`: call `GoalCheck`; `!ok` → escalate.
- Runs **once before close** (not every tick) — cost/HITL-fatigue conscious (hermes's
  Adapt, not full Adopt).

**Acceptance:** a confident `GOAL-CHECK: FAIL` escalates before completion; a `PASS`
completes; a checker error/timeout warns and proceeds (does not block).

## Non-goals / deferred
- **LE-11a batch/rate-limit driver-opened HITL questions** — deferred. Per-idea fatigue is
  already bounded (the driver halts on escalation rather than flooding); cross-idea
  question batching belongs to the Tier-4 outer loop (`standing-loop-watch-mode`), where
  many ideas open. Documented, not built here.
- The goal-check verdict is an LLM judgment (fresh, different-agent), not a proof — it is a
  defense-in-depth signal, fail-open on its own failure.

## Tests
- `driver`: `TriageReserved` + auto → escalate; `<2` reviewers + auto → escalate; `≥2` +
  Ready + GoalCheck pass → complete; GoalCheck fail → escalate; non-auto design idea with
  1 reviewer still completes.
- `app`: `GoalCheck` parses PASS/FAIL/ambiguous from a consult answer (fail-open on
  ambiguous).
- Drift guard green (edit BOTH COOPERATION.md copies).
