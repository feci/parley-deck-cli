---
agent: hermes-1
idea: track-aware-driver
review-round: 2
date: 2026-07-03
reviewed-commit: 70bdc5c
---

## Summary

Fix-up cycle 1 (commit 70bdc5c) addresses my round-01 MAJOR-1 plus four of the
six round-01 MINOR/NIT findings (the remaining two were open questions, not
required fixes). I verified every claimed fix against the live source, the
fix-up diff, the test suite, and four CLI smoke runs. `go test ./...` is green
in this environment — including `internal/runner/TestDurableKillEndToEndRealProcess`,
which PASSES here (the codex-sandbox "no recorded boot id" limitation does not
reproduce outside the sandbox, confirming the standing exception). `go vet
./internal/{track,driver,app}/...` is clean.

The five round-01 findings the fix-up cycle targeted are all RESOLVED. One new
MINOR finding: the fast-track model-diversity hard gate (the MAJOR-1 fix) has no
dedicated test — the three existing `checkModelDiversity` tests all use
absent-track or explicit `require_model_diversity: true`, none sets `track: fast`
to exercise the new `if t == track.Fast { required = true }` branch.

## Verification of round-01 findings

### [MAJOR-1] fast did not force a model-diverse reviewer — RESOLVED

Claim: `checkModelDiversity` now makes model diversity a hard gate on
`track: fast` regardless of the frontmatter flag.

Verified in source: `internal/app/driver_impl.go:154-160`:
  required := driver.ReadRequireModelDiversity(o.ideaDir)
  if t, present := driver.ReadTrack(o.ideaDir); present && t == track.Fast {
      required = true
  }
When `required=true`, the function returns an escalation error (line 176-177)
instead of the warning path (line 179). The gate fires before any reviewer is
launched (called from `OpenReviewRound`), so a same-model single-reviewer fast
idea is now a hard reject, not a warning. This closes the §4.0 "model-diverse"
safety invariant my round-01 review flagged.

Refutation attempt: can the gate be bypassed? `checkModelDiversity` is the only
model-diversity check; it reads `ReadTrack` (the same `00-prompt.md` source `New`
uses) and short-circuits only when `!same` (line 151-153). For a fast idea with
all reviewers sharing the implementer's model, `same=true` and `required=true`
→ error returned → `OpenReviewRound` propagates it. No bypass. RESOLVED.

Caveat (new finding, below): the new branch is not exercised by any test.

### Explicit `track: deliberation` non-solo hard-reject — RESOLVED

Claim: `PolicyFor` moved the non-solo check above the switch, so every explicit
track (including deliberation) hits the hard-reject; absent-track stays exempt.

Verified in source: `internal/track/track.go:133-140`:
  if !present {
      return Policy{Track: Standard, ApplyOverrides: false, CrossReviewRounds: -1}, nil
  }
  if availableReviewers < 1 {
      return Policy{}, fmt.Errorf("track: %s requires at least 1 independent reviewer (non-solo, §1); none available", t)
  }
  switch t { ... }
The `!present` early-return (line 133-135) preserves the absent-track legacy
path (no non-solo error — preflight handles it). The `availableReviewers < 1`
check (line 138-140) runs BEFORE the switch, so fast, standard, AND deliberation
all hit it. The per-case `availableReviewers < 1` checks that were inside the
fast/standard cases have been removed (diff confirms deletion of lines 185-186
and 195-196). The driver `New` passes `avail := distinctNonImplementers(...)` to
`PolicyFor` and records `trackErr` on error (driver.go:121-123), and `Advance`
hard-gates on `trackErr` (first statement). So an explicit deliberation idea
with a solo roster now escalates on the first `Advance`.

Tests: `TestPolicyForDeliberationNonSolo` (track_test.go:246-255) verifies
`PolicyFor(Deliberation, true, 0, ...)` errors AND `PolicyFor(Standard, false, 0,
...)` does NOT (absent exemption). `TestExplicitDeliberationNonSoloEscalates`
(driver/track_test.go:87-95) verifies `New` sets `trackErr` for a solo
deliberation roster. `TestAbsentTrackSoloDoesNotTrackError` (track_test.go:97-106)
verifies the absent-track exemption. All pass. RESOLVED.

### Classifier: unknown/negative size no longer fast — RESOLVED

Claim: `Inputs.FilesKnown/LOCKnown` added; `Classify` requires size positively
known AND non-negative for fast; `parley classify` sets known via `fs.Visit` and
rejects negative counts with exit 2.

Verified in source:
- `internal/track/track.go:56-57`: `FilesKnown`/`LOCKnown` bool fields added.
- `internal/track/track.go:99-101`: fast condition now requires
  `in.FilesKnown && in.Files >= 0 && in.Files <= 5 && in.LOCKnown && in.LOC >= 0
  && in.LOC <= 300`. Unknown size (Known=false) or negative size fails to
  standard.
- `internal/app/classify.go:40-43`: `if *files < 0 || *loc < 0` → stderr message
  + `return 2`.
- `internal/app/classify.go:51-60`: `fs.Visit` sets `FilesKnown`/`LOCKnown` only
  for flags actually supplied on the command line.

Smoke tests (built `cmd/parley`, ran locally):
1. `classify --reversible --mechanically-verifiable` (no --files/--loc) →
   `standard`, exit 0. ✓ (unknown size is not fast)
2. `classify --files 1 --loc -1 --reversible --mechanically-verifiable` →
   `--files and --loc must be non-negative`, exit 2. ✓
3. `classify --files 2 --loc 20 --reversible --mechanically-verifiable` →
   `fast`, exit 0. ✓ (known small size is still fast)

Tests: `TestClassifyFastAndStandard` (track_test.go:56-90) now uses
`fastable := Inputs{..., FilesKnown: true, LOCKnown: true, ...}` and adds cases
for unknown size → standard, LOC-unknown → standard, negative LOC → standard.
All pass. RESOLVED.

### fast + auto_implement + --no-implement now escalates — RESOLVED

Claim: `New` now reads the IDEA-LEVEL `ReadAutoImplement(cfg.IdeaDir)` /
`ReadStrictGate(cfg.IdeaDir)` for the contradiction check, not the runtime-masked
`cfg.AutoImplement`.

Verified in source: `internal/driver/driver.go:117-121`:
  // The §4.0 contradiction check uses the IDEA-LEVEL auto_implement / strict_gate
  // (review-01 fix), not cfg.AutoImplement — the latter is masked to false by the
  // runtime --no-implement brake, which would otherwise let fast + auto_implement
  // slip past the contradiction gate.
  pol, err := track.PolicyFor(t, present, avail, ReadAutoImplement(cfg.IdeaDir), ReadStrictGate(cfg.IdeaDir))
`ReadAutoImplement` (transport.go:46-51) reads `00-prompt.md` directly and
ignores the runtime `--no-implement` flag. So a `track: fast` idea with
`auto_implement: true` in frontmatter now gets `trackErr` set at construction
time regardless of `--no-implement`, and `Advance` escalates immediately.

Refutation attempt: `TestFastContradictionEscalates` (track_test.go:64-75)
writes `track: fast` + `auto_implement: true` and passes `AutoImplement: true` +
`Auto: true` to `New`, then calls `Advance` and asserts `ActionEscalated` +
error. This test exercises the contradiction path. The fix changes `New` to read
from `IdeaDir` rather than `cfg.AutoImplement`, but the test's `Config` sets both
`AutoImplement: true` AND writes `auto_implement: true` to the 00-prompt, so both
the old and new code paths would set `trackErr`. The test does NOT specifically
isolate the `--no-implement` masking scenario (where `cfg.AutoImplement=false`
but `ReadAutoImplement(ideaDir)=true`). However, the code change is verifiable by
inspection: `New` line 121 unconditionally calls `ReadAutoImplement(cfg.IdeaDir)`
and never references `cfg.AutoImplement` for the contradiction check. The masking
path is structurally closed. RESOLVED (by code inspection; the specific
`cfg.AutoImplement=false` + frontmatter-`true` isolation test is a minor
test-coverage gap, not a blocker).

### standard caps cross-review at 2 — RESOLVED

Claim: `Policy.CapCrossReviewRounds` added; explicit standard clamps
`CrossReviewRounds` to 2.

Verified in source:
- `internal/track/track.go:119`: `CapCrossReviewRounds int` field added to
  `Policy` (0 = no cap).
- `internal/track/track.go:159`: standard case returns
  `CapCrossReviewRounds: 2`.
- `internal/driver/driver.go:131-133`:
  if pol.CapCrossReviewRounds > 0 && cfg.CrossReviewRounds > pol.CapCrossReviewRounds {
      cfg.CrossReviewRounds = pol.CapCrossReviewRounds
  }
  This clamps (not overwrites) — so if frontmatter says 1, it stays 1; if it says
  10, it clamps to 2. The clamp runs inside `if pol.ApplyOverrides`, so only
  explicit standard is affected; absent/deliberation keep `ApplyOverrides=false`
  and are untouched.

Tests: `TestPolicyForStandardCapsCrossReview` (track_test.go:257-262) verifies
`p.CapCrossReviewRounds == 2`. `TestExplicitStandardCapsCrossReview`
(driver/track_test.go:108-116) writes `track: standard` with
`CrossReviewRounds: 5` and asserts `d.cfg.CrossReviewRounds == 2`. Both pass.
RESOLVED.

### Full suite green — RESOLVED

`go test ./...` → exit 0, zero FAIL lines. `go vet
./internal/{track,driver,app}/...` → clean (no output). The
`TestDurableKillEndToEndRealProcess` test PASSES in this environment (0.13s),
confirming the codex-sandbox "no recorded boot id" limitation is environmental,
not a code defect. RESOLVED.

## New findings (if any)

### [MINOR-new-1] No test exercises the fast-track model-diversity hard gate

The MAJOR-1 fix adds `if t, present := driver.ReadTrack(o.ideaDir); present && t
== track.Fast { required = true }` to `checkModelDiversity`
(driver_impl.go:158-160). The three existing `checkModelDiversity` tests are:
- `TestCheckModelDiversityWarnsAndEmitsEvent` — uses `writePrompt(t, ideaDir, "")`
  (absent track, no `require_model_diversity`) → warns, does NOT escalate.
- `TestCheckModelDiversityEscalatesWhenRequired` — uses
  `writePrompt(t, ideaDir, "require_model_diversity: true\n")` (absent track,
  explicit flag) → escalates via the frontmatter flag, NOT the fast-track branch.
- `TestCheckModelDiversitySilentWhenDiverse` — diverse roster, no track set.

None writes `track: fast` to the 00-prompt, so the new `t == track.Fast` branch
is never exercised. A regression that breaks only that branch (e.g., someone
removes the `if t == track.Fast` check) would not be caught. The branch is
small and verifiable by inspection, but a one-line test addition
(`writePrompt(t, ideaDir, "track: fast\n")` + same-model roster + assert error)
would close the gap.

This does not block acceptance — the fix is correct and the gate is structurally
sound — but it is a test-coverage gap on the highest-risk safety invariant
(same-model rubber-stamp on a single-reviewer track).

### Round-01 MINOR findings not addressed by fix-up cycle 1 (for completeness)

These were open questions / MINOR findings the implementer chose not to address
in cycle 1. They are NOT blockers and were not required fixes:

- [MINOR-1] Boundary thresholds `Files > 15` / `LOC > 1000` (vs `>=`): UNCHANGED.
  Still `> 15` / `> 1000` (track.go:90-91). The IMPLEMENTATION.md fix-up section
  does not list this as addressed. This was my round-01 open question #1 — the
  implementer retained the "≤" reading. Acceptable as a documented judgment call.
- [MINOR-3] `TestFastNonSoloEscalates` still does not call `Advance`: UNCHANGED
  (track_test.go:82-84 checks `d.trackErr == nil` only). The new
  `TestExplicitDeliberationNonSoloEscalates` mirrors this pattern (trackErr only,
  no Advance). `TestFastContradictionEscalates` remains the only non-solo/
  contradiction test that calls `Advance`. Minor test-gap, not a safety issue —
  the `Advance` hard-gate on `trackErr` is verified by the contradiction test.
- [MINOR-4] No driver-level test for `fast + strict_gate` via `New` + `Advance`:
  PARTIALLY addressed. `TestExplicitDeliberationNonSoloEscalates` adds a
  driver-level non-solo test for deliberation, but there is still no
  `TestFastStrictGateEscalates` (fast + strict_gate → Advance → ActionEscalated)
  and no `TestStandardNonSoloEscalates`. The policy-level test
  (`track_test.go:102`) covers fast+strict_gate at the `PolicyFor` layer. Minor
  gap.
- [NIT-1, NIT-2]: Not addressed (dead MinReviewers for fast; raw "declared"
  string in JSON). Cosmetic, non-blocking.

## Signoff

Status: ✅ ACCEPT

All five round-01 findings targeted by fix-up cycle 1 are RESOLVED — verified
against live source, the diff, the test suite, and CLI smoke runs. The full
suite is green (including the runner test that is the standing sandbox
exception, which passes here). `go vet` is clean. The one new MINOR finding
(no test for the fast-track model-diversity hard gate) is a test-coverage gap on
a correct fix, not a safety defect. The unaddressed round-01 MINOR/NIT findings
were open questions or cosmetic, none blocking. The implementation meets all
five observable acceptance criteria.
