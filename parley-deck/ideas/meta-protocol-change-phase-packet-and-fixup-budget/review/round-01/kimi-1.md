---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 1
date: 2026-08-12
---
verdict: CLEAN

# Review round 1 — parley-deck-cli 1.44.0 + parley-deck-skill 2.8.0 (SLICE 1 only)

Independent verification per §15.1: I ran every check myself against the working trees of
both repos; nothing below relies on IMPLEMENTATION.md's self-report. Scope judged: the two
§4.0 budget cells + the inclusive-cap fix. The packet (slice 2) is correctly absent from the
diff and I did not review it.

Verification environment (PRIMARY): `go version go1.26.5 darwin/arm64`,
`GOCACHE=$PWD/.gocache GOMODCACHE=$PWD/.gomodcache`.

## Q1 — Is the inclusive cap correct at every boundary?

Yes. PRIMARY:

- `internal/driver/impl.go:283-286` (working tree): `cycle := round; if cycle >
  d.cfg.MaxFixupCycles { escalate }`. `round` is `highestReviewRound(ideaDir)`
  (`internal/driver/cursor.go:194-197`); fix-up cycle *k* runs when round *k*'s consensus
  carries outstanding agreed fixes, then opens round *k+1*. So cap N publishes cycles 1..N
  and escalates exactly when cycle N+1 would start.
- I ran the boundary suite myself:

  ```
  $ go test ./internal/driver/ -count=1 -v -run 'TestFixupCapIsInclusive|...'
  --- PASS: TestFixupCapIsInclusive (6 subcases, all run)
  --- PASS: TestPhaseReviewMaxFixupCyclesEscalates
  ok  parley-deck-cli/internal/driver
  ```

  Subcases cover cap 5 at 5/6, cap 2 at 2/3, cap 1 at 1/2 — i.e. fast(1): cycle 1 runs,
  cycle 2 escalates; standard(2): 1..2 run, 3 escalates; deliberation(5): 1..5 run, 6
  escalates. All three boundaries are correct.
- Anything else depending on the guard: `MaxFixupCycles` has exactly two consumers —
  the fix-up guard (`impl.go:284`, changed) and the strict-gate guard (`impl.go:215`,
  unchanged, see Q2). Repo-wide grep for `MaxFixupCycles` shows no other enforcement site
  (the `internal/app/driver_impl.go:315` hit is a comment). The AF2 crash-recovery path
  (`impl.go:139-153`) re-opens round N+1 without re-checking the cap, but that completes a
  transition already authorized before `Fixup` ran — not a bypass.
- Two side observations, neither blocking:
  - **(NIT)** The escalation message at `impl.go:285` now reads "after %d cycle(s)" with
    `cycle = N+1` at the moment of escalation, when N fix-up cycles actually ran. The old
    message was off by one in the other direction. Substance (escalate, with counts) is
    right; the count in the string is not.
  - **(MINOR)** Absent-track legacy ideas (no `track:` in 00-prompt) use the driver default
    3 (`driver.go:103-104`) and silently move from 2 to 3 published fix-up cycles under the
    inclusive guard. This is the ratified semantics — FINAL binding condition 2 explicitly
    pre-registered boundary tests "at 5/6 and 3/4", and `TestPhaseReviewMaxFixupCyclesEscalates`
    pins the 3/4 case — but the CHANGELOG table lists only fast/standard/deliberation, so
    this behaviour change is invisible to a reader. One line in the CHANGELOG would close it.

## Q2 — The second `round >= MaxFixupCycles` guard (strict-gate closing loop)

Leaving it as `>=` is right. PRIMARY: `internal/driver/impl.go:212-217` — the strict-close
loop runs only when `OutstandingAgreedFixes == 0`; it publishes **review rounds**, not
fix-up cycles (no `Fixup` call exists on that path). The FINAL's inclusive-boundary
requirement is written about *published fix-up cycles* (FINAL.md §2, binding condition 2),
a different unit. Flipping this guard to `>` would grant strict ideas N+1 closing rounds —
a *loosening* of LE-2 termination that no round of this idea reviewed. The pre-existing
`TestStrictGateBoundedByMaxFixupCycles` pins the `>=` behaviour at `round == cap` and passes
(I ran it).

Two caveats the record should carry:

- **(MINOR)** The deviation's reasoning — "changing both under one mandate would alter a
  bound nobody in this idea reviewed" — is half-true. The operator was preserved, but the
  bound itself **did** change for deliberation, from 3 to 5 effective closing rounds,
  because both loops share the one knob. That numeric change is a legitimate consequence of
  the ratified cell (LE-2 ratified "bounded by MaxFixupCycles", the knob, not a number —
  and the protocol prose at `parley-deck/COOPERATION.md:644` says exactly that), but the
  deviation text should be read as "operator preserved", not "bound untouched".
- Resulting asymmetry (observation, no action): a non-strict deliberation idea can reach
  review round 6 via five fix-ups; a strict deliberation idea escalates at round 5
  not-clean. Directionally consistent (strict is meant to be tighter) and it terminates,
  escalates, never auto-closes. It will not "bite". Both guards belong in the follow-up
  audit's enumeration (`meta-protocol-change-track-gate-enforcement-audit`).

## Q3 — `ApplyOverrides: true` for deliberation: anything beyond the two cells?

No unintended override found. PRIMARY, traced through `driver.New` (`driver.go:128-147`)
field by field:

- `CrossReviewRounds: -1` → not applied (`pol.CrossReviewRounds >= 0` is false); the
  idea-configured value survives. `CapCrossReviewRounds: 3` → clamp fires only when the
  configured value exceeds 3 (`driver.go:132-134`); the default 1 (set at `driver.go:100-102`
  before policy application) is untouched. Pinned by
  `TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops` (configured 1 preserved) and
  `TestDeliberationClampsCrossReviewRoundsToThree` (configured 9 → 3). Both pass (I ran
  them).
- `MaxReviewers`: assigned unconditionally at `driver.go:138`, but the policy value is 0
  ("all non-implementers"), matching prior effective behaviour — and no caller passes a
  competing value: all three `driver.New` call sites (`internal/app/app.go:1202-1219`,
  `1934-1951`, `1988-2005`) set `CrossReviewRounds` only. The reviewer trim at
  `internal/app/driver_impl.go:59` keys on `pol.MaxReviewers > 0`, still 0 for deliberation
  → no trim, unchanged.
- `MinReviewers: 0` → not applied (`pol.MinReviewers > 0` false) → falls through to the
  default 2 (`driver.go:145-147`). **The LE-11 floor of 2 is preserved.**
- LE-11 non-solo floor: `PolicyFor` still hard-errors for deliberation at
  `availableReviewers < 1` (`track.go:137-140`), escalated via `trackErr` on the first
  `Advance` (`driver.go:222-224`). `TestPolicyForDeliberationNonSolo` passes (I ran it).
- **(NIT)** Comment drift: `track.go:115` ("false = reproduce today's behaviour
  (absent/deliberation)"), `driver.go:110-111` ("absent/deliberation preserve today's
  behaviour byte-for-byte") and `track.go:151` ("Deliberation == today's full lifecycle")
  now describe the old no-op policy. The new comment block at `track.go:154-159` explains
  the change, so the file argues with itself. Cosmetic, but these are exactly the comments
  a future reader will trust.

## Q4 — The three rewritten tests: honest or weakened?

Honest. PRIMARY (diff read in full, tests re-run):

1. `TestPhaseReviewMaxFixupCyclesEscalates` — fixture moved round-03 → round-04, assertions
   unchanged (escalated action; `must not run fixup past MaxFixupCycles`). With an
   inclusive cap of 3, round-03's fixes *should* run fix-up cycle 3, so the old fixture no
   longer exercises escalation; round-04 is the ratified boundary. The guarantee "past the
   cap escalates, no fix-up runs" is preserved at the correct edge.
2. `TestNewDeliberationIsLegacy` → `TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops`
   — still pins `MaxReviewers=0`, `MinReviewers=2`, and a configured `CrossReviewRounds=1`
   preserved; changes only the ratified cell (fix-up 3 → 5).
3. `TestPolicyForDeliberationIsLegacy` →
   `TestPolicyForDeliberationBoundsBothLoopsAndKeepsTheRest` — flips the `ApplyOverrides`
   assertion (that no-op was the bug being fixed) and *adds* pins on the untouched fields
   (`MaxReviewers=0`, `MinReviewers=0`, `CrossReviewRounds=-1`). Strengthened, not weakened.

The new `TestFixupCapIsInclusive` also adds run-at-boundary assertions the old suite never
had (previously only escalation-at-boundary was tested). Net coverage increased.

## Q5 — Protocol text: three copies, wording, residual "unbounded"

- The two changed cells are **byte-identical in all three copies**. PRIMARY: the three
  `git diff` hunks carry the same old→new lines, and
  `diff <(sed -n '200,240p' internal/protocol/defaults/COOPERATION.md) <(sed -n '200,240p'
  ../parley-deck-skill/skills/parley-deck/references/COOPERATION.md)` on the §4.0 region
  prints nothing. The whole-file diff between embedded default and skill copy shows only
  the two pre-existing bootstrap header lines (Transport/Created placeholders). Deck ↔
  embedded identity outside the five allowlisted zones is machine-proven by
  `TestEmbeddedDefaultMatchesLiveDeck` — PASS (I ran it).
- Wording vs code: "cap 5 cycles" ↔ `MaxFixupCycles: 5` + inclusive guard
  (`impl.go:284`). "capped at 3 after round 1, then escalate" ↔ `CapCrossReviewRounds: 3`
  clamp (`driver.go:132-134`) plus the round budget `CurrentRound >= 1+CrossReviewRounds`
  (`driver.go:262,271`). Accurate, with one nuance: for cross-review the driver *clamps*
  rather than escalates — "then escalate" binds the agent-driven flow (same shape as
  standard's pre-existing "capped at 2, then escalate/upgrade"). **(NIT)** Pre-existing
  shape, not introduced here.
- Residual old behaviour: none. `[Uu]nbounded` appears in **zero** locations in any of the
  three copies (grep, all three files). The only `unlimited` hit is the LE-5 loop-ceiling
  prose (`0 means unlimited`, `COOPERATION.md:671`), a different feature, correct. §4 Phase
  8 prose (`COOPERATION.md:644,646-664`) — "strict-close loop is bounded by the fix-up
  budget", "escalation thresholds, not close criteria" — matches the code exactly. No
  other part of the protocol asserts the old behaviour.

## Q6 — CHANGELOG accuracy, esp. "fast published ZERO"

Accurate and not overstated. Independent verification, all PRIMARY:

- **fast = 0**: `git show HEAD:internal/driver/impl.go` has the old
  `if cycle >= d.cfg.MaxFixupCycles` guard, and `git show HEAD:internal/track/track.go`
  has fast at `MaxFixupCycles: 1`. Trace: round-01 consensus with outstanding fixes →
  `cycle=1 >= 1` → escalate *before* any `Fixup` call. Zero published. (Belt and braces:
  a real fast idea can never run `Fixup` anyway — fast + auto_implement is rejected at
  `track.go:143-145` and `!AutoImplement` escalates at `impl.go:287-289`.) The claim stands.
- standard "printed 2, published 1" and deliberation "unbounded vs default 3, published 2"
  verify by the same trace (old guard escalates at cycle 2 and cycle 3 respectively).
- **"a standard idea was measured running 15 cycles"**: independently confirmed —
  `parley-deck/ideas/skills-cli-install-path/00-prompt.md` carries `track: standard` and
  its IMPLEMENTATION.md contains 15 `^## Fix-up cycle` headings with no escalation
  mentioned. The claim's mechanism (a printed cap binds only where enforcement lives) is
  consistent with the code: under auto-drive an explicit standard idea would have been
  capped, so that loop was agent-driven.
- **Distribution**: my own recount over the current deck
  (`grep -c '^## Fix-up cycle' parley-deck/ideas/*/IMPLEMENTATION.md`) gives
  `0×36, 1×34, 2×7, 3×2, 4×3, 5×2, then 9, 14, 15, 25`. Every non-zero bucket matches the
  CHANGELOG exactly; the 0-bucket grew 17 → 36 because the deck grew from 69 to 88 ideas
  with IMPLEMENTATION.md since the FINAL's measurement. The load-bearing claim — nothing
  has ever closed in the 6–8 band — reproduces.
- **(NIT)** Skill CHANGELOG: "the printed number and the enforced number are the same
  number" is scoped to "both cells" and is true for them; a reader could overgeneralize —
  absent-track ideas still run driver default 3 against the standard column's printed 2
  (known, and explicitly deferred to the follow-up audit by FINAL §4). Accurate as scoped.

## Q7 — Should anything stop this release?

No blocker found. Build, vet, and the full test suite are green on my own run
(`go build ./... && go vet ./... && go test ./... -count=1` → `BUILD_VET_OK`, all 27
packages `ok`, including driver, track, and the protocol drift guard). Versions are
consistent: `VERSION` 1.44.0, `internal/app/version.go` 1.44.0, skill `package.json` 2.8.0,
`compatibility.json` 2.8.0.

Advisories only:

- **(MINOR) Commit hygiene.** The working tree mixes in other ideas' process artifacts:
  `parley-deck/ideas/protocol-read-cost-regression/` (IMPLEMENTATION.md +59, review/
  consensus.md +1374, untracked `review/round-04/`) and untracked
  `parley-deck/ideas/speedup-tooling-evaluation/`. The code/text diff itself is exactly
  slice 1; if the release commit sweeps the whole tree, two unrelated ideas' artifacts land
  in it. Commit selectively.
- **(NIT)** Escalation-message off-by-one (Q1), stale comments (Q3), cross-review
  "then escalate" vs clamp (Q5), skill-CHANGELOG generalization risk (Q6).
- The declared "no end-to-end run" gap (IMPLEMENTATION.md) is honest and acceptable for a
  one-operator change with boundary tests at every track; it restates a known
  instrumentation gap rather than hiding a new one.

## Verdict

**CLEAN.** Slice 1 ships what FINAL ratified — the two deliberation cells (5 inclusive
fix-up cycles, cross-review clamped at 3 after round 1), the inclusive-cap fix that also
corrects the printed-vs-enforced divergence on every track, and identical protocol text in
all three copies — with honest tests, an accurate CHANGELOG, and the packet correctly out
of scope. The findings above are 3× MINOR + 4× NIT; none blocks the release.

### Findings index

| # | Severity | Finding |
| --- | --- | --- |
| F1 | MINOR | Absent-track legacy ideas silently move 2 → 3 published fix-up cycles; invisible in the CHANGELOG table (Q1) |
| F2 | MINOR | Deviation framing: the strict-gate `>=` preserves the operator but deliberation's strict-close bound still moved 3 → 5 via the shared knob (Q2) |
| F3 | MINOR | Working tree carries two other ideas' artifacts; the release commit should be selective (Q7) |
| F4 | NIT | Escalation message reports "after N+1 cycle(s)" when N fix-ups ran (`impl.go:285`) (Q1) |
| F5 | NIT | Stale comments: `track.go:115`, `track.go:151`, `driver.go:110-111` still describe deliberation as a no-op policy (Q3) |
| F6 | NIT | Cross-review cell says "then escalate"; the driver clamps (pre-existing shape, same as standard) (Q5) |
| F7 | NIT | Skill CHANGELOG's "same number" claim is scoped-true but overgeneralizable to absent-track ideas (Q6) |
