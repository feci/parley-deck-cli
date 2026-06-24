---
idea: verification-honesty
drafter: claude-1
implementer: claude-1
date: 2026-06-24
status: final
spawned_from: loop-engineering-research
---

## Purpose

Implementation spec for Tier 1 (verification honesty). Closes the three inert-gate
seams. Each section names exact files + the observable acceptance criteria a reviewer
(or the driver) can check.

## LE-1 — Refutation-default review

**Files:** `internal/runner/phase58.go` (`BuildReviewPrompt`, `ValidateReviewArtifact`);
`COOPERATION.md` §Phase 6 (both copies).

- `BuildReviewPrompt`: change the posture from neutral/confirmatory to **adversarial** —
  "assume the implementation is wrong until you fail to break it; for each FINAL.md
  acceptance criterion, attempt to construct a failing case or run the relevant check;
  report 'no findings' only after stating what you tried that failed to break it."
- Add a `## Refutation attempts` section to the required review file shape.
- `ValidateReviewArtifact`: require the `## Refutation attempts` section (in addition
  to `## Findings`), so an empty-findings review must show its work.
- §Phase 6: one normative line — "Reviewers default to refutation; a 'no findings'
  review must record what was attempted under `## Refutation attempts`."

**Acceptance:** a review file lacking `## Refutation attempts` fails
`ValidateReviewArtifact`; the prompt contains the adversarial posture and the section.

## LE-3 — Model-diversity guard

**Files:** `internal/app/driver_impl.go` (`newDriverImplOps` / `OpenReviewRound`);
`internal/driver/transport.go` (new `ReadRequireModelDiversity`); `COOPERATION.md`
§Phase 6.

- When the review round opens, compare the implementer's `agents.Discovery.Model`
  against every reviewer's. If **all** reviewers share the implementer's model, emit a
  warning (stdout + an `agent.model_diversity` event).
- Opt-in escalation: 00-prompt `require_model_diversity: true` → the driver escalates
  before opening review instead of warning (checked in `advanceImpl`).
- Default = **warn** (non-breaking). No-op when models differ (today's roster).

**Acceptance:** with a same-model reviewer set, opening review prints a warning; with
`require_model_diversity: true` it escalates; differing models → silent.

## LE-2 — strict_gate enforcement

**Files:** `internal/driver/transport.go` (`ReadStrictGate`); `internal/driver/driver.go`
(`Config.StrictGate`); `internal/app/app.go` (wire `StrictGate: driver.ReadStrictGate(ideaDir)`);
`internal/driver/impl.go` (`advanceReview` close path); `internal/app/driver_impl.go`
(`ReviewStatus` parses new fields; `DraftReviewConsensus` prompt carries strict_gate);
`internal/runner/phase58.go` (`BuildReviewConsensusPrompt` adds the two fields under
strict_gate; `ValidateReviewConsensusArtifact` tolerant); `COOPERATION.md` §Phase 8.

Contract: under `strict_gate: true`, completion (`OutstandingAgreedFixes == 0`) is **not**
sufficient. The driver requires a **fresh full-scope closing review round with zero
findings of any severity**, certified by the drafter and not contradicted by a
deterministic scan:
- Review consensus gains two machine-readable fields the drafter sets under strict_gate:
  `closing_review_round: <N>` and `strict_gate_clean: <true|false>` (true only when round
  N's reviewers reported zero findings of any severity).
- `driver.ReviewStatus` carries `StrictGateClean bool` and `ClosingReviewRound int`.
- In `advanceReview`, when `OutstandingAgreedFixes == 0` AND `StrictGate`:
  - If `StrictGateClean && ClosingReviewRound == round` AND a deterministic finding-scan
    of round N's review files finds **no real finding** → `Complete()`.
  - Else → open one more fresh review round (round+1) as the strict closing round and
    continue (await → reviews → consensus). Bounded by `MaxFixupCycles` (escalate if
    exceeded) so the strict-close loop always terminates.
  - The deterministic scan (`reviewRoundHasFindings`) counts a `### [CRITICAL|MAJOR|
    MINOR|NIT] <title>` heading as a real finding only when `<title>` is non-empty and
    not the literal placeholder `<title>`. It can only **veto** a clean claim (fail
    closed), never auto-pass — the positive decision still needs the drafter's
    certification.
- Without `strict_gate`, behavior is unchanged (`OutstandingAgreedFixes == 0` completes).

**Acceptance:** a `strict_gate: true` idea does not complete on the first 0-fix round
unless that round is certified clean AND scans clean; a drafter that claims clean while
a review file has a real finding → escalates; the strict-close loop is bounded.

## LE-4 — Generalize RunChecks + tie artifact-wins

**Files:** `internal/app/driver_impl.go` (`RunChecks`); `COOPERATION.md` §Phase 4/5
(document `checks:`).

- `RunChecks` resolution order:
  1. If 00-prompt frontmatter has `checks: <command>` → run `sh -c "<command>"` in root;
     pass/fail on exit code.
  2. Else if `go.mod` exists in root → `go test ./...` (today's behavior).
  3. Else if the idea is code-writing (`auto_implement: true`) → **fail closed**:
     `(false, "no checks configured for a code-writing idea; set checks: in 00-prompt")`.
  4. Else (design-only, non-Go) → `(true, "no checks to run")`.
- Because `advanceImpl` (pre-review) and `advanceReview` (post-fix-up) both call
  `RunChecks` and escalate on failure, making it fail-closed for code ideas
  **transitively ties the "artifact-wins" fix-up override to a real check** (hermes #8):
  a fix-up that produced a valid-shaped artifact but a broken build now escalates at the
  post-fix-up gate instead of passing. Documented in IMPLEMENTATION.md.

**Acceptance:** an `auto_implement` idea with no `checks:` and no `go.mod` → RunChecks
fails closed; a `checks:` command is honored; a design-only idea still passes.

## Conditional rigor (rejected over-application, kept)
LE-1's `## Refutation attempts` requirement is universal (cheap, always-on). strict_gate
(LE-2) and RunChecks fail-closed (LE-4 step 3) scale with `strict_gate`/`auto_implement`
respectively — trivial/design-only ideas keep the lightweight close. (Consensus REJECTED
uniform rigor.)

## Tests
- `phase58_test.go`: review prompt contains refutation posture + `## Refutation attempts`;
  ValidateReviewArtifact rejects a file without it.
- `driver` tests: strict_gate forces a fresh closing round; finding-scan veto escalates;
  bounded by MaxFixupCycles; non-strict unchanged.
- `driver_impl` tests: RunChecks honors `checks:`, fails closed for code-no-checks,
  passes design-only; model-diversity warn/escalate.
- Drift guard `TestEmbeddedDefaultMatchesLiveDeck` stays green (edit BOTH COOPERATION.md).
