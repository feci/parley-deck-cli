---
agent: hermes-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 1
date: 2026-08-12
---

verdict: CLEAN

# Phase 6 code review — round 1

Release candidate: parley-deck-cli 1.44.0 + parley-deck-skill 2.8.0.
Slice 1 only: the two §4.0 deliberation budget cells + the inclusive-cap fix.
The phase-scoped packet is not in scope and not in this diff.

I ran the tests myself rather than trusting IMPLEMENTATION.md. All evidence
below is tagged per §15.2: PRIMARY means I observed it in the working tree,
quoted a command output, or traced it in source I read.

## Q1 — Is the inclusive-cap change correct at every boundary?

**Yes.** [PRIMARY, `git diff internal/driver/impl.go`]

The changed guard at `internal/driver/impl.go:284`:

```go
cycle := round
if cycle > d.cfg.MaxFixupCycles {
    return ActionEscalated, ...
}
```

was `cycle >= d.cfg.MaxFixupCycles`. `round` is `highestReviewRound(d.cfg.IdeaDir)`
(line 132), so `cycle` equals the current review round number. When
`OutstandingAgreedFixes > 0`, this guard decides run-fixup vs escalate.

Trace by track (MaxFixupCycles from `internal/track/track.go:149,172,165`):

| Track | MaxFixupCycles | Old `>=` published | New `>` publishes | New `>` escalates at |
|-------|---------------|--------------------|--------------------|----------------------|
| fast | 1 | 0 (cycle 1 → 1>=1 → escalate) | 1 (cycle 1 runs) | cycle 2 |
| standard | 2 | 1 (cycle 2 → 2>=2 → escalate) | 2 (cycles 1,2 run) | cycle 3 |
| deliberation | 5 | 2 (driver default 3, cycle 3 → 3>=3) | 5 (cycles 1-5 run) | cycle 6 |

[PRIMARY: values verified by `TestFixupCapIsInclusive` — 6 sub-tests, all PASS,
`go test ./internal/driver/... -run TestFixupCapIsInclusive -count=1 -v`]

The new boundary test `TestFixupCapIsInclusive` (`internal/driver/impl_test.go:432`)
covers all three tracks at the cap and one past it:

- cap=5, cycle 5 → runs fixup (true); cap=5, cycle 6 → escalates (false)
- cap=2, cycle 2 → runs fixup (true); cap=2, cycle 3 → escalates (false)
- cap=1, cycle 1 → runs fixup (true); cap=1, cycle 2 → escalates (false)

Each sub-test asserts both `ranFixup` and `action == ActionEscalated` when
escalation is expected. This is an honest boundary test, not a tautology.

**Does this change anything else that depends on that guard?** No. The guard
at line 284 is the sole entry to the fix-up path. The only other `MaxFixupCycles`
guard is the strict-gate loop at line 215 (Q2 below). No other code path reads
or branches on this comparison. [PRIMARY: `search_files` for `MaxFixupCycles`
in `internal/driver/` returns exactly two guard sites: line 215 and line 284,
plus the Config struct field and test assertions.]

## Q2 — The second `round >= d.cfg.MaxFixupCycles` guard (strict-gate, line 215)

**Leaving it as `>=` is defensible but introduces a real asymmetry. I concur
with the implementer's reasoning, with one caveat noted as a finding.**

[PRIMARY, `internal/driver/impl.go:212-217`]

```go
if !certifiedClean {
    // Not a certified-clean closing round yet → run one more fresh review
    // round. Bounded by MaxFixupCycles so the strict-close loop terminates.
    if round >= d.cfg.MaxFixupCycles {
        return ActionEscalated, c, fmt.Errorf("strict_gate: no clean closing review round after %d round(s) ...")
    }
    // ... open round+1 ...
}
```

This guard fires when `OutstandingAgreedFixes == 0` but strict_gate is on and
the current round is not certified clean. It opens a fresh closing review round.
The guard bounds that loop.

The asymmetry: with `MaxFixupCycles=5` and strict_gate=true on deliberation:

- Fix-up loop (line 284, `cycle > 5`): cycles 1-5 run, cycle 6 escalates.
- Strict-gate loop (line 215, `round >= 5`): rounds 1-4 can open a closing
  round, round 5 escalates. **One fewer round than the fix-up loop.**

[PRIMARY: `TestStrictGateBoundedByMaxFixupCycles` — unchanged, uses
MaxFixupCycles=3, rounds 1-3, StrictGateClean=false, expects escalation at
round 3. This test exercises the `>=` guard and passes because the guard was
not changed. `go test -run TestStrictGateBoundedByMaxFixupCycles -v` → PASS.]

The implementer's reasoning (IMPLEMENTATION.md "Deviation from FINAL"): the
FINAL's inclusive-boundary requirement is written about *published fix-up
cycles*, and LE-2's strict-gate termination semantics were ratified by a
different idea. Changing both under one mandate would alter a bound nobody
in this idea reviewed.

I accept this reasoning. The two loops are genuinely different: the fix-up
loop runs when there are agreed fixes to apply; the strict-gate loop runs
when there are zero fixes but the closing round isn't certified clean. They
share `MaxFixupCycles` as a bound, but the semantics of what they count are
not identical — a fix-up cycle produces an IMPLEMENTATION.md entry, a strict-
gate closing round does not. The FINAL's language ("inclusive published
cycles") targets the former.

However, the asymmetry means a strict_gate deliberation idea that reaches
round 5 with zero agreed fixes but no certified-clean round will escalate
instead of getting its fifth closing round — while a fix-up idea at the same
round with agreed fixes would still run cycle 5. This is a minor surprise but
not a correctness defect: escalation is the safe direction, and the strict-
gate loop's purpose is termination, not maximizing rounds.

### [MINOR] strict-gate loop has a different inclusive boundary than the fix-up loop

The strict-gate closing-round guard at `impl.go:215` remains `round >=
MaxFixupCycles` while the fix-up guard at `impl.go:284` is now `cycle >
MaxFixupCycles`. For the same `MaxFixupCycles` value, the strict-gate loop
gets one fewer round than the fix-up loop. This is deliberate (see
IMPLEMENTATION.md "Deviation from FINAL") and the escalation direction is
safe, but the two guards now read inconsistently against the same field.
A future reader will need the IMPLEMENTATION.md note to understand why they
differ. No change required for this release; flag for the follow-up audit
(`meta-protocol-change-track-gate-enforcement-audit`) to decide whether LE-2
should also become inclusive.

## Q3 — Does `ApplyOverrides: true` change anything beyond the two intended cells?

**No. Only MaxFixupCycles and the cross-review cap change. All other fields
are preserved.** [PRIMARY, trace through `driver.go:128-146` + `track.go:160-166`]

The deliberation policy returned by `PolicyFor` (`track.go:160-166`):

```go
return Policy{
    Track:                Deliberation,
    ApplyOverrides:       true,
    CrossReviewRounds:    -1, // keep whatever the idea configured…
    CapCrossReviewRounds: 3,  // …but never above 3 rounds after round 1
    MaxFixupCycles:       5,
}, nil
```

The override application in `driver.go:128-142`:

```go
if pol.ApplyOverrides {
    if pol.CrossReviewRounds >= 0 { cfg.CrossReviewRounds = pol.CrossReviewRounds }
    if pol.CapCrossReviewRounds > 0 && cfg.CrossReviewRounds > pol.CapCrossReviewRounds {
        cfg.CrossReviewRounds = pol.CapCrossReviewRounds
    }
    if pol.MaxFixupCycles > 0 { cfg.MaxFixupCycles = pol.MaxFixupCycles }
    cfg.MaxReviewers = pol.MaxReviewers
    if pol.MinReviewers > 0 { cfg.MinReviewers = pol.MinReviewers }
}
```

Field-by-field:

- **MaxReviewers**: `pol.MaxReviewers = 0` (zero value, not set in the Policy
  struct literal). `cfg.MaxReviewers = 0`. The old path (ApplyOverrides=false)
  also left it at 0 (driver default). 0 = "all non-implementers." **Unchanged.**

- **MinReviewers**: `pol.MinReviewers = 0` (zero value). The guard
  `if pol.MinReviewers > 0` is false, so `cfg.MinReviewers` is NOT set here.
  It falls through to `driver.go:145-146`: `if cfg.MinReviewers <= 0 {
  cfg.MinReviewers = 2 }`. The LE-11 non-solo floor of 2 is preserved.
  **Unchanged.** [PRIMARY: `TestPolicyForDeliberationBoundsBothLoopsAndKeepsTheRest`
  asserts `p.MinReviewers == 0` and `p.MaxReviewers == 0` — PASS]

- **CrossReviewRounds**: `pol.CrossReviewRounds = -1`, so the
  `>= 0` guard is false — the configured value is preserved. But
  `pol.CapCrossReviewRounds = 3 > 0`, so if the configured value exceeds 3
  it is clamped. **This is the intended cross-review cap.** [PRIMARY:
  `TestDeliberationClampsCrossReviewRoundsToThree` — configured 9, got 3.
  PASS. Also `TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops` —
  configured -1 (default→1), got 1 (not clamped, since 1 < 3). PASS.]

- **MaxFixupCycles**: `pol.MaxFixupCycles = 5 > 0`, so `cfg.MaxFixupCycles
  = 5`. **This is the intended fix-up cap.**

The test `TestPolicyForDeliberationBoundsBothLoopsAndKeepsTheRest`
(`track_test.go:164`) explicitly asserts the "rest is untouched" invariant:
`p.MaxReviewers == 0`, `p.MinReviewers == 0`, `p.CrossReviewRounds == -1`.
[PRIMARY: PASS, `go test -run TestPolicyForDeliberationBounds -v`]

The `TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops`
(`track_test.go:44`) asserts the driver-level result: `MaxReviewers=0,
MinReviewers=2, MaxFixupCycles=5, CrossReviewRounds=1`. [PRIMARY: PASS]

**No invisible override.** The zero-value fields in the Policy struct are
handled by the conditional application logic — `MinReviewers > 0` prevents
a zero from clobbering the default, and `MaxReviewers = 0` is the same as
the default. The LE-11 floor is structurally preserved by the fallthrough at
line 145-146.

## Q4 — Were the three rewritten tests honest, or was a guarantee weakened?

**Honest. Each rewritten test preserves the invariant it was checking, and
adapts the expected values to the new (correct) semantics.**

1. `TestNewDeliberationIsLegacy` → `TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops`
   (`track_test.go:44`)

   Old: asserted `MaxReviewers=0, MinReviewers=2, MaxFixupCycles=3,
   CrossReviewRounds=1` — the "legacy" deliberation = today's full lifecycle.
   New: asserts `MaxReviewers=0, MinReviewers=2, MaxFixupCycles=5,
   CrossReviewRounds=1`.

   The invariant the old test protected — "deliberation keeps the full
   lifecycle, no reviewer cap, LE-11 min 2" — is still asserted (Max=0,
   Min=2, Cross=1 preserved). Only Fixup changed from 3 to 5, which is the
   intended change. The test name now says "BoundsBothLoops" instead of
   "IsLegacy" because deliberation is no longer a no-op. **Not a weakening.**

2. `TestPolicyForDeliberationIsLegacy` → `TestPolicyForDeliberationBoundsBothLoopsAndKeepsTheRest`
   (`track_test.go:164`)

   Old: asserted `ApplyOverrides == false` ("must NOT override").
   New: asserts `ApplyOverrides == true` AND `MaxFixupCycles == 5` AND
   `CapCrossReviewRounds == 3` AND `MaxReviewers == 0, MinReviewers == 0,
   CrossReviewRounds == -1`.

   The new test is strictly **stronger** than the old: it still checks that
   the rest of the lifecycle is untouched (Max=0, Min=0, Cross=-1), but it
   also verifies the two new bounds. The old test's guarantee ("deliberation
   does not override") was encoding a bug — the whole point of this idea is
   that deliberation SHOULD override its two budget cells. **Not a weakening;
   the replaced guarantee was wrong.**

3. `TestPhaseReviewMaxFixupCyclesEscalates` (`impl_test.go:366`)

   Old: fixture created rounds 1-3, MaxFixupCycles=3, asserted escalation
   (cycle 3 == cap → `>=` → escalate).
   New: fixture creates rounds 1-4, MaxFixupCycles=3, asserts escalation.

   The test still verifies the same guarantee — "past the cap, the driver
   escalates instead of running fix-up." It moved the fixture from round 3
   to round 4 because with the inclusive `>` guard, cycle 3 is now allowed
   and cycle 4 is the first one that escalates. The assertion (`action ==
   ActionEscalated`, `fixup` not called) is identical. **Not a weakening;
   the fixture was adjusted to match the corrected boundary.**

[PRIMARY: all three tests PASS, `go test ./internal/driver/...
./internal/track/... -count=1 -v`]

## Q5 — Protocol text: three copies identical, accurate, no residual "unbounded"

**All three copies are byte-identical in the changed rows. The wording is
accurate against the code. No other part of the protocol still asserts the
old "unbounded" behaviour.**

[PRIMARY: `diff` of the §4.0 table region across all three copies — IDENTICAL]

The two changed cells in all three copies:
- Cross-review: `capped at 3 after round 1, then escalate`
- Fix-up: `cap 5 cycles; \`strict_gate\` available`

Accuracy against code:
- Fix-up "cap 5 cycles": `MaxFixupCycles=5`, guard `cycle > 5` → cycles 1-5
  run. Accurate.
- Cross-review "capped at 3 after round 1": `CapCrossReviewRounds=3`,
  `CrossReviewRounds` clamped to 3, driver guard `CurrentRound >= 1+3=4` →
  rounds 1-4 run (round 1 + 3 cross-review). Accurate. [PRIMARY: trace
  through `driver.go:271` and `driver.go:128-134`]

Residual "unbounded" search: [PRIMARY: `search_files` for `unbounded` in
both repos] — the only matches in protocol text are in CHANGELOG.md entries
describing what *was* changed (past tense, correct), and in historical idea
artifacts (`protocol-read-cost-regression` rounds/FINAL, which are records
of past deliberation, not normative text). No live protocol copy still
asserts "unbounded" for any deliberation cell.

The `§4.0` prose at line 224 still says "This table is the single
authoritative per-track gate" — which was the sentence the FINAL flagged as
aspirational when cells had no enforcing path. With both cells now enforced,
that claim is closer to true for these two cells. The follow-up audit
(`meta-protocol-change-track-gate-enforcement-audit`) is where the remaining
cells get dispositioned, per FINAL §4.

## Q6 — Is the CHANGELOG accurate and not overstated?

**Accurate. I independently verified the central claim.**

The CHANGELOG's most striking claim: "fast published no fix-up cycle at all
while its cell claimed one." [PRIMARY, verified by trace]

Old code: `cycle >= MaxFixupCycles`, fast track `MaxFixupCycles=1`
(`track.go:149`).

At round 1 (first review round) with outstanding agreed fixes:
- `cycle = round = 1`
- `1 >= 1` → TRUE → `return ActionEscalated`
- `Fixup()` is never reached.

Fast published 0 fix-up cycles. **Claim verified.**

The table in the CHANGELOG ("Actually published" column) is also verified:
- fast: 0 (shown above)
- standard: 1 (cycle 1 runs, cycle 2 → 2>=2 → escalate)
- deliberation: 2 (driver default 3, cycles 1-2 run, cycle 3 → 3>=3 → escalate)

The "Now" column:
- fast: 1 (cycle 1 → 1>1 → false → runs; cycle 2 → 2>1 → escalate)
- standard: 2 (cycles 1-2 run; cycle 3 → 3>2 → escalate)
- deliberation: 5 (cycles 1-5 run; cycle 6 → 6>5 → escalate)

All verified by `TestFixupCapIsInclusive` sub-tests. [PRIMARY: PASS]

The CHANGELOG's "Why 5" paragraph reproduces the distribution from FINAL.md
§2 verbatim. The "Not in this release" section correctly states the packet
is not started. No overstatement detected.

The skill CHANGELOG (2.8.0) is also accurate: it describes the two cell
changes, references the ratifying idea, and notes the CLI enforces both in
the same release.

## Q7 — Anything that should stop this release

**No.**

- All tests pass: `go test ./... -count=1` → all packages green. [PRIMARY]
- `go vet ./...` → clean. [PRIMARY]
- Version bumps are consistent: VERSION=1.44.0, version.go=1.44.0,
  CHANGELOG=1.44.0; skill package.json=2.8.0, package-lock.json=2.8.0,
  compatibility.json=2.8.0, CHANGELOG=2.8.0. [PRIMARY]
- Protocol drift guard (`TestEmbeddedDefaultMatchesLiveDeck`) passes — the
  three COOPERATION.md copies are in lockstep. [PRIMARY: PASS]
- The one declared deviation (strict-gate guard left as `>=`) is safe
  (escalation direction), explicitly documented, and flagged for the
  follow-up audit. It does not block this release.
- No end-to-end run was performed (IMPLEMENTATION.md "Not verified"), but
  this is the same instrumentation gap already on record and the unit tests
  exercise the changed guards directly.

## Findings

### [MINOR] strict-gate loop has a different inclusive boundary than the fix-up loop

See Q2 above. `impl.go:215` uses `round >= MaxFixupCycles` while `impl.go:284`
uses `cycle > MaxFixupCycles`. The strict-gate loop gets one fewer effective
round than the fix-up loop for the same cap. This is deliberate and safe
(escalation direction), but the inconsistency should be dispositioned in the
follow-up audit. No action required for this release.

### [NIT] IMPLEMENTATION.md line references are slightly stale

IMPLEMENTATION.md references `internal/driver/impl.go:284` and
`internal/driver/impl.go:215`. The `>` guard is at line 284 in the working
tree (confirmed), and the strict-gate `>=` guard is at line 215 (confirmed).
These are accurate as of the working tree. No action needed; noting that
line numbers will shift if other edits land before commit.

## Open questions

None. The release is correct for what it ships. The strict-gate asymmetry
(Q2) is the only substantive observation, and it is explicitly deferred to
the follow-up audit by the implementer's own deviation note.
