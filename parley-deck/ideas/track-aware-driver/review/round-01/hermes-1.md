---
agent: hermes-1
idea: track-aware-driver
review-round: 1
date: 2026-07-03
reviewed-commit: a9a0ff4
---

## Summary

The implementation delivers the MVP (Slice 1 + Slice 2): a pure `internal/track`
classifier + policy deriver, driver-internal track derivation in `New`, a
hard-gate on `Advance` for contradictions/non-solo, reviewer truncation in
`newDriverImplOps`, and the `parley classify` CLI command. `go test ./...` is
green and `go vet ./internal/{track,driver,app}/...` is clean (both re-run
during this review). The backward-compat invariant (absent-track = today) holds
because `ReadTrack` returns `(Standard, false)` for absent/unknown, and
`PolicyFor` with `present=false` sets `ApplyOverrides=false`, leaving all knobs
untouched. The refutation-default validator (`## Refutation attempts` in
`ValidateReviewArtifact`) is not gated on track and runs for every review
artifact, so criterion 4 is structurally preserved.

The deviations (driver-internal derivation, deferred timeout, template absent)
are acceptable per the FINAL.md scope and do not break any acceptance criterion.
I found no CRITICAL issue — the hard-rejects fire and halt, the classifier is
fail-safe, and non-solo is enforced at construction time. I found one MAJOR
gap (§4.0's "model-diverse" requirement for fast is not enforced), several
MINOR test-coverage gaps, and boundary interpretation questions where the
code's strict `>15`/`>1000` thresholds may under-tier vs the spec's "fail
closed to the stricter track on boundary doubt" principle.

## Refutation attempts (what you tried per criterion + result)

### Criterion 1 — classify implements §4.0 verbatim + fail-safe

Tried: all boundary values files=5/6/15/16, loc=300/301/1000/1001 with
--reversible --mechanically-verifiable. Result: files=5/loc=300 → fast (✓);
files=6 → standard (✓); files=15 → standard (✓); files=16 → deliberation (✓);
loc=301 → standard (✓); loc=1000 → standard (✓); loc=1001 → deliberation (✓).
All pass. Tried fail-safe: --reversible without --mechanically-verifiable →
standard (✓); --mechanically-verifiable without --reversible → standard (✓).
Tried every deliberation trigger with small/fast-looking inputs (files=1,
loc=1, reversible, mech-verifiable + --security / --protocol-change /
--auto-implement / --strict-gate / --pipeline / --api-break / --schema-break /
--irreversible / --data-migration) → all return deliberation (✓). The
classifier is correct and fail-safe.

Attempted refutation on boundary interpretation: §4.0 lines 185-190 say "on
any doubt or boundary case — the 6-14-file band, the ~300-~1000 LOC gap —
fail closed to the stricter track." The code uses `Files > 15` (16+) and
`LOC > 1000` (1001+) for the deliberation triggers, so files=15 and loc=1000
go to standard, not deliberation. The spec's "~15" and "~1000" are
approximate, but the fail-safe principle argues for `>= 15` / `>= 1000`. This
is a judgment-call boundary, not a clear violation — see Findings [MINOR-1].

### Criterion 2 — per-track Config values

Tried: `TestNewAppliesFastTrack` verifies fast → MaxReviewers=1, MinReviewers=1,
CrossReviewRounds=0, MaxFixupCycles=1 (✓). `TestNewExplicitStandardAppliesCaps`
verifies standard → MaxReviewers=2, MinReviewers=2, MaxFixupCycles=2 (✓).
`TestNewDeliberationIsLegacy` verifies deliberation → MaxReviewers=0,
MinReviewers=2, MaxFixupCycles=3, CrossReviewRounds=1 (✓).
`TestNewAbsentTrackIsLegacy` verifies absent → same as deliberation (✓).

Attempted refutation on CrossReviewRounds for standard: §4.0 table says
"capped at 2" for standard cross-review rounds. The code sets
`CrossReviewRounds: -1` (leave as configured) for standard, so if frontmatter
says `cross_review_rounds: 10`, standard gets 10, not 2. The FINAL.md
criterion 2 uses "·" for CrossReviewRounds (meaning "leave as-is"), so this is
by FINAL.md design — but it does NOT implement the §4.0 table's "capped at 2."
See Findings [MINOR-2].

Attempted refutation on MinReviewers vs MaxReviewers consistency: for fast
(1/1), standard (2/2 or 2/1-degraded), deliberation (0/2), absent (0/2) — no
case has MinReviewers > MaxReviewers. The two-participant degradation for
standard (MinReviewers=1, MaxReviewers=2) is safe: 1 available reviewer → no
truncation (1 > 2 is false), LE-11 requires ≥1. Consistent. Not broken.

### Criterion 3 — hard-rejects fire and halt

Tried: `TestFastContradictionEscalates` verifies fast + auto_implement →
trackErr set, Advance returns ActionEscalated + error (✓).
`TestFastNonSoloEscalates` verifies fast with 1 participant → trackErr set
(✓ for trackErr, but does NOT call Advance — see Findings [MINOR-3]).

Attempted refutation: can the hard-reject be bypassed? The `trackErr` check is
the first statement in `Advance` (driver.go:214), before the auto-drive gate
and before any phase dispatch. `loop.Run` calls `Advance` in a loop. There is
no other entry point to the phase advances. All `Impl` calls are inside
`advanceFinal`/`advanceImpl`/`advanceReview`, which are only reached through
`Advance`. So if `trackErr != nil`, no `Impl` method is ever invoked. Cannot
be bypassed. Not broken.

Attempted refutation on fast + strict_gate: `PolicyFor(Fast, true, ..., false,
true)` returns an error (track.go:133). `track_test.go` tests this at the
policy level (line 102). But there is NO driver-level test that constructs
`New` with a `track: fast` + `strict_gate: true` 00-prompt and verifies
Advance escalates. See Findings [MINOR-4].

### Criterion 4 — refutation non-optional on every track

Tried: searched for any track-gated skip of `OpenReviewRound`,
`DraftReviewConsensus`, or `ValidateReviewArtifact`. Result: `advanceImpl`
always calls `OpenReviewRound(ctx, 1)` when `Impl` is wired — no track check.
`advanceReview` always calls `OpenReviewRound` for re-opens — no track check.
`ValidateReviewArtifact` (runner/phase58.go:437) requires a non-empty
`## Refutation attempts` section for every review artifact — no track check.
The `CrossReviewRounds=0` override for fast shortens the cross-review phase
but does not skip the review phase (Phase 6-8). Review always runs when `Impl`
is wired. Not broken.

### Criterion 5 — backward-compat, full suite green

Tried: `go test ./...` → all packages green. `go vet
./internal/{track,driver,app}/...` → clean. The existing driver tests use
`setupIdea` which writes 00-prompt.md WITHOUT a `track:` field, so `ReadTrack`
returns `(Standard, false)` → `ApplyOverrides=false` → all knobs untouched.
The `loop_budget_test.go` tests build Config without `IdeaDir` → track
derivation skipped entirely. Backward-compat holds.

Attempted refutation on the `cfg.IdeaDir` side effect: `New` now reads
`cfg.IdeaDir` to derive track policy. Any path that sets `IdeaDir` to a dir
whose 00-prompt.md contains `track: fast`/`track: standard` will get
overridden knobs. The three production Config sites (app.go:1156, 1829, 1883)
all set `IdeaDir` to the real idea dir — so a production idea with `track:
standard` WILL get MaxReviewers=2, MaxFixupCycles=2 (not today's ∞/3). This is
the designed opt-in. No existing test has `track:` in 00-prompt.md, so no
existing test changes behavior. Not broken for the current suite, but the
side effect is real and new — see Findings [MINOR-5].

### Deviations review

- **D2 (driver-internal derivation):** Acceptable. `New` is the single
  construction chokepoint — stronger than three app-layer call sites. Same
  observable result. No criterion broken.
- **Per-track timeout deferred:** Acceptable. FINAL.md criteria 1-5 do not
  mention timeout. The §4.0 table includes it, but the FINAL.md scoped it out.
  Not a safety item.
- **Template left absent:** Acceptable. Keeps absent = today (criterion 5).
  Auto-defaulting `track: standard` into every new idea would silently change
  new-idea behavior. The opt-in path (`parley classify` + explicit `track:`)
  is preserved.

## Findings

### [MAJOR-1] §4.0 "model-diverse" requirement for fast is not enforced

§4.0 table line 199: fast reviewers = "1 (model-diverse)." The implementation
does not set `RequireModelDiversity=true` for `track: fast`. A fast idea with
a single reviewer sharing the implementer's model gets only a WARNING (not a
hard reject) from `checkModelDiversity` (driver_impl.go:149-175), unless
`require_model_diversity: true` is explicitly set in frontmatter. The hermes-1
round-01 analysis (line 318-322) recommended auto-setting this for fast. The
FINAL.md criterion 2 does not list model-diversity, so this is within FINAL.md
scope, but it leaves a §4.0 safety invariant unenforced: fast's single
reviewer is supposed to be model-diverse, and a same-model single reviewer is
exactly the rubber-stamp scenario LE-3 warns about — compounded by fast having
only one reviewer.

Concrete fix: in `PolicyFor`'s fast case (track.go:128-138), add a
`RequireModelDiversity bool` field to `Policy`, set it `true` for fast, and
apply it in `New` (or in `newDriverImplOps`'s `checkModelDiversity` call).
Alternatively, document this as a deferred slice in IMPLEMENTATION.md.

### [MINOR-1] Boundary thresholds may under-tier vs §4.0 fail-safe principle

`Classify` uses `Files > 15` and `LOC > 1000` for the deliberation triggers.
At files=15 or loc=1000, the classifier returns standard, not deliberation.
§4.0 lines 185-190 say "on any doubt or boundary case — fail closed to the
stricter track." The spec's "~15" and "~1000" are approximate, but a strict
fail-safe reading would use `>= 15` / `>= 1000`. The current thresholds treat
the exact boundary as standard, which is the weaker track. This is a judgment
call, not a clear violation — but in REFUTATION-DEFAULT mode, the boundary
should fail to the stricter side.

Concrete fix: change `in.Files > 15` to `in.Files >= 15` and `in.LOC > 1000`
to `in.LOC >= 1000`, and update the test expectations accordingly. Or document
that 15/1000 are intentionally standard (the "≤" reading of the spec).

### [MINOR-2] Standard track does not cap CrossReviewRounds at 2 (§4.0 table)

§4.0 table line 197: standard cross-review rounds = "capped at 2, then
escalate/upgrade." The code sets `CrossReviewRounds: -1` (leave as configured)
for standard. If an idea declares `track: standard` with
`cross_review_rounds: 10`, the driver runs 10 cross-review rounds, not 2.
FINAL.md criterion 2 uses "·" for CrossReviewRounds, making this a deliberate
MVP scoping decision — but it does not implement the §4.0 table's cap.

Concrete fix: either set `CrossReviewRounds: 2` for standard in `PolicyFor`
(overriding the frontmatter value), or document in IMPLEMENTATION.md that the
cross-review cap is deferred and standard respects the frontmatter value.

### [MINOR-3] TestFastNonSoloEscalates does not verify Advance escalation

`TestFastNonSoloEscalates` (track_test.go:77-85) checks `d.trackErr == nil`
but does NOT call `d.Advance()` to verify `ActionEscalated` + error. The
contradiction test (`TestFastContradictionEscalates`) does call Advance. The
non-solo path uses the same `trackErr` check in `Advance`, so it is
structurally identical — but the test gap means a future regression that
breaks only the non-solo error path (e.g., a nil-trackErr bug specific to
non-solo) would not be caught.

Concrete fix: add `act, _, err := d.Advance(context.Background())` and
`if act != ActionEscalated || err == nil { t.Errorf(...) }` to
`TestFastNonSoloEscalates`, mirroring `TestFastContradictionEscalates`.

### [MINOR-4] No driver-level test for fast + strict_gate or standard + non-solo

`fast + strict_gate` contradiction is tested at the policy level
(`track_test.go:102`) but not at the driver level (no `New` + `Advance`
test). `standard + non-solo` (0 available reviewers) is tested at the policy
level (`track_test.go:130`) but not at the driver level. Only `fast +
auto_implement` has a driver-level escalation test.

Concrete fix: add `TestFastStrictGateEscalates` and
`TestStandardNonSoloEscalates` to `track_test.go`, mirroring
`TestFastContradictionEscalates` (write 00-prompt with `track: fast` +
`strict_gate: true` / `track: standard` with 1 participant, call Advance,
assert ActionEscalated).

### [MINOR-5] Double PolicyFor derivation (New + newDriverImplOps) with auto_implement inconsistency

`driver.New` calls `PolicyFor(t, present, avail, cfg.AutoImplement,
cfg.StrictGate)` using `cfg.AutoImplement`, which is
`driver.ReadAutoImplement(ideaDir) && !noImplement` (app.go:1167).
`newDriverImplOps` calls `PolicyFor(t, present, len(reviewers),
driver.ReadAutoImplement(ideaDir), driver.ReadStrictGate(ideaDir))` using the
raw frontmatter value (ignoring `--no-implement`). When `--no-implement` is
passed for a `track: fast` idea with `auto_implement: true` in frontmatter:
`New` sees `AutoImplement=false` (no contradiction → fast policy applied, no
trackErr), but `newDriverImplOps` sees `auto_implement=true` (contradiction →
truncation skipped). The result: fast track runs with all reviewers instead of
1. This is fail-safe (more reviewers, not fewer), and the idea won't
auto-implement anyway (`--no-implement`), so the contradiction is moot. But
it's a behavioral inconsistency from two independent derivations reading
different sources for the same signal.

Concrete fix: pass `cfg.AutoImplement` (the already-resolved value) to
`newDriverImplOps` instead of re-reading `ReadAutoImplement(ideaDir)` in the
truncation path, or move the truncation into `New`-controlled state.

### [NIT-1] MinReviewers for fast is effectively dead config at runtime

`MinReviewers=1` for fast is never checked at runtime because the LE-11 guard
(impl.go:240) is gated on `d.cfg.AutoImplement`, and fast + auto_implement is
a contradiction (always rejected). So every valid fast idea has
`AutoImplement=false`, meaning the `MinReviewers` guard never fires for fast.
Non-solo for fast is enforced solely at construction time in `PolicyFor`
(`availableReviewers < 1`). This is not a bug — the construction-time check is
stronger (fires before any advance) — but `MinReviewers=1` for fast gives a
false impression of runtime enforcement. Consider documenting that fast's
non-solo is construction-time only, or removing the runtime guard's
`AutoImplement` gating for fast.

### [NIT-2] classify --json "declared" field shows raw string, not normalized track

`classify.go:85`: `out["declared"] = *declared` stores the raw user input
(e.g., "FAST"), not the normalized track. The validation uses the normalized
value via `track.Normalize(*declared)`. This is cosmetic — the `valid` and
`message` fields are correct — but the JSON consumer sees the unnormalized
string.

## Open questions

1. Is the boundary at files=15 / loc=1000 intentionally standard (the "≤"
   reading), or should it be deliberation per the fail-safe principle? The
   spec's "~15" / "~1000" is ambiguous, but REFUTATION-DEFAULT favors the
   stricter side. (MINOR-1)

2. Should the standard track cap CrossReviewRounds at 2 per the §4.0 table, or
   is the FINAL.md "·" (leave-as-is) the final word for this MVP? If the
   latter, should IMPLEMENTATION.md explicitly note that the §4.0 cross-review
   cap is deferred? (MINOR-2)

3. Should fast auto-set `RequireModelDiversity=true` as hermes-1 round-01
   recommended, or is the warning-only path acceptable for the MVP? This
   determines whether a same-model single-reviewer fast idea is a hard reject
   or just a warning. (MAJOR-1)

4. The double-derivation of `PolicyFor` (in `New` and in `newDriverImplOps`)
   is wasteful and creates the `--no-implement` inconsistency (MINOR-5). Was
   this chosen to keep `newDriverImplOps` self-contained (not depending on
   `driver.Config`), or is it an oversight? Could the truncation read
   `MaxReviewers` from the already-constructed `Driver` instead?
