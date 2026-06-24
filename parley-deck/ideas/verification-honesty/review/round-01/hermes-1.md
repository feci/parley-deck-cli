---
agent: hermes-1
idea: verification-honesty
review-round: 1
date: 2026-06-24
---

## Summary

I reviewed FINAL.md (the spec) against the implementation diff 1962cb5..HEAD in
refutation mode, assuming the code was wrong until I failed to break it. I read the
full `advanceReview`/`advanceImpl` control flow (internal/driver/impl.go), the
RunChecks/model-diversity logic (internal/app/driver_impl.go), the prompt/validator
changes (internal/runner/phase58.go), both COOPERATION.md copies, and the new tests;
I also ran `go build ./...`, `go vet ./...`, `go test -count=1 ./internal/driver/
./internal/runner/ ./internal/app/`, and the drift guard — all green, and the two
COOPERATION.md edits are byte-identical around the LE regions.

Verdict: the three inert-gate seams are honestly closed for the intended/common path.
The strict_gate logic genuinely gates `Complete()`: under `strict_gate`, completion is
reachable only when `certifiedClean := StrictGateClean && ClosingReviewRound == round`
is true AND `reviewRoundHasFindings` returns false; otherwise a fresh round opens
(bounded by MaxFixupCycles, which `driver.New` defaults to 3). A drafter cannot
trivially fake past the finding-scan veto for standard-format findings — the
`ClosingReviewRound == round` check defeats pointing at a stale clean round, and the
scan reads the actual reviewer files (consensus.md is correctly excluded). RunChecks
fail-closed (LE-4) is correctly ordered (checks -> go.mod -> fail-closed -> pass) and
the post-fix-up `RunChecks` gate in `advanceReview` does tie the artifact-wins override
to a real check transitively.

I did break one thing: the finding-scan veto fails OPEN on directory/file read errors,
which contradicts the spec's "can only veto (fail closed), never auto-pass" invariant.
That is the only finding that touches a stated invariant; the rest are minor gaps.

## Refutation attempts

LE-1 (refutation-default review):
- Tried to pass `ValidateReviewArtifact` with `## Refutation attempts` present but
  EMPTY (just the header, no content). Result: it PASSES — the check is
  `strings.Contains(data, "## Refutation attempts")` (phase58.go:434-437), presence
  only. A rubber-stamping reviewer satisfies LE-1 with an empty section. See Findings
  (MINOR).
- Tried a review missing the section entirely: rejected (confirmed by code + the
  TestValidateReviewArtifactRequiresRefutation test). Holds.
- Confirmed `BuildReviewPrompt` emits the adversarial posture and the section in the
  required shape (TestBuildReviewPromptRefutationDefault). Holds.

LE-2 (strict_gate enforcement):
- Tried a drafter lying `strict_gate_clean: true` + `closing_review_round: <round>`
  while the round's review file contains `### [MAJOR] real bug`. Result: escalates
  (TestStrictGateVetoesUncleanCertification + traced impl.go:198-200). The drafter
  cannot trivially fake past the veto for standard-format findings. Holds.
- Tried stale-round evasion: drafter sets `closing_review_round: 1` while the current
  round is 2 and round-02 has a finding. `certifiedClean = true && (1 == 2)` = false
  -> opens round 3, does NOT complete (fail-closed). Holds.
- Tried completing on the first 0-fix round without certification: `!certifiedClean`
  -> opens a fresh round, no Complete() (TestStrictGateOpensFreshClosingRound). Holds.
- Tried a certified-clean round that scans clean: completes
  (TestStrictGateCompletesWhenCertifiedClean). Matches "does not complete on the first
  0-fix round UNLESS certified clean AND scans clean." Holds.
- Tried finding shapes the scan does NOT match: `### [Major] x` (wrong case),
  `### [MAJOR]` (empty title), `### [MAJOR] <title>` (unreplaced placeholder), and a
  finding written in prose with no heading. `scanHasRealFinding` (impl.go:334-355)
  misses all of them. This is BY DESIGN per FINAL.md ("counts ... only when <title> is
  non-empty and not the literal placeholder <title>"), and the severity tags are
  specified uppercase, so case-sensitivity is correct. Not a bug, but it bounds the
  veto's strength to reviewers who follow the heading format.
- Tried an `os.ReadDir` failure on the round directory while `certifiedClean` is true.
  `reviewRoundHasFindings` returns `false` on ANY ReadDir error (impl.go:312-314) ->
  no veto -> falls through to `Complete()`. FAIL-OPEN. See Findings (MAJOR). The same
  fail-open applies per-file at `os.ReadFile` (impl.go:321-323: `continue` on error).
- Tried exhausting the bound: `round >= MaxFixupCycles` -> escalates
  (TestStrictGateBoundedByMaxFixupCycles; `New` defaults MaxFixupCycles to 3 so no
  zero-value trap). Holds — the strict-close loop terminates.
- Tried a reviewer whose agent ID is literally `consensus`: its file `consensus.md` is
  skipped by the scan (impl.go:318) but still counted by `ReviewRoundComplete`. See
  Findings (NIT).

LE-3 (model-diversity guard):
- Tried same-model roster: `reviewersShareImplementerModel` returns `(true, model)`
  (TestReviewersShareImplementerModel). Holds at the helper level.
- Tried unknown implementer model and differing models: guard is silent (no false
  positive). Holds.
- Tried to find a test that exercises the OBSERVABLE acceptance behavior — "opening
  review prints a warning" and "require_model_diversity: true escalates" inside
  `OpenReviewRound` (driver_impl.go:147-155). No such test exists; only the pure
  helper is unit-tested. See Findings (MINOR).

LE-4 (RunChecks generalize + artifact-wins tie):
- Tried `auto_implement: true`, no go.mod, no checks -> fail closed
  (TestRunChecksFailsClosedForCodeIdeaWithoutChecks). Holds.
- Tried `checks: "exit 7"` -> fails; `checks: "exit 0"` -> passes
  (TestRunChecksHonorsChecksCommand). Holds.
- Tried design-only (no auto_implement, no go.mod, no checks) -> passes
  (TestRunChecksDesignOnlyPasses). Holds.
- Tried to break the artifact-wins tie: traced `advanceReview` impl.go:246-253 —
  after `Fixup` it calls `RunChecks` and escalates on failure. Combined with
  fail-closed for code-no-checks, a fix-up that wrote a valid-shaped artifact but
  cannot be verified escalates. Holds (the transitive tie claimed in IMPLEMENTATION.md
  is real).
- Tried `checks:` with shell metacharacters: run via `sh -c <checks>` in the workspace
  root. Metacharacters are intended (sh -c is the point); `checks` is author-controlled
  frontmatter, not untrusted input. Not a finding.

Drift guard: `TestEmbeddedDefaultMatchesLiveDeck` green; both COOPERATION.md copies
carry identical LE-1/3/2/4 edits. Holds.

## Findings

### [MAJOR] Finding-scan veto fails OPEN on read errors, contradicting the spec's "fail closed, never auto-pass"

`reviewRoundHasFindings` (internal/driver/impl.go:310-330) returns `false` on ANY
`os.ReadDir` error, and skips any reviewer file whose `os.ReadFile` errors
(impl.go:321-323). When `certifiedClean` is true, a `false` from the scan means NO
veto, so `advanceReview` falls through to `Complete()` (impl.go:198-221). Thus, if the
round directory or a reviewer file is unreadable at scan time, a drafter's
`strict_gate_clean: true` claim passes WITHOUT the deterministic scan actually
verifying it.

Why it matters: FINAL.md LE-2 states the scan "can only veto a clean claim (fail
closed), never auto-pass." A read-error path that auto-passes is exactly the class of
fail-open/false-green seam this idea exists to eliminate. The codebase already has the
correct fail-closed pattern for the same operation — `highestReviewRoundErr`
(cursor.go:199-256) surfaces non-`NotExist` ReadDir errors rather than silently
treating them as "no rounds." Practical trigger probability is low (the directory is
known to exist because `ReviewRoundComplete` just confirmed reviewer files are present,
and the round-label functions match: `roundDirLabel` == `roundLabel` == `round-%02d`),
but for a verification gate the default on "I could not scan" must be "I cannot certify
clean," not "clean."

Concrete fix: mirror the existing fail-closed pattern. On a non-`fs.ErrNotExist`
`ReadDir` error, return `true` (veto -> escalate) — or, cleaner, change
`reviewRoundHasFindings` to `(bool, error)` and have `advanceReview` escalate on the
error. At minimum:
```go
entries, err := os.ReadDir(dir)
if err != nil {
    if errors.Is(err, fs.ErrNotExist) {
        return false // absent dir -> nothing to veto (ReviewRoundComplete guards this)
    }
    return true // unreadable but should exist -> fail closed
}
```
Apply the same fail-closed default to the per-file `os.ReadFile` error (veto rather
than `continue`), since a file `ReviewRoundComplete` confirmed present becoming
unreadable should not silently drop its findings from the scan.

### [MINOR] `## Refutation attempts` validation is presence-only, undermining LE-1's "must show its work" intent

`ValidateReviewArtifact` (internal/runner/phase58.go:434-437) accepts the section via
`strings.Contains(data, "## Refutation attempts")`. A review containing only the bare
header `## Refutation attempts` with zero content passes. LE-1's stated goal is that
"an empty-findings review must show its work," but the gate only proves the header
exists, not that work was recorded. This is consistent with the existing presence-only
`## Findings` check (so not a regression), and it meets the letter of FINAL.md's
acceptance criterion ("a review file lacking `## Refutation attempts` fails"), but a
rubber-stamping reviewer can satisfy it trivially with an empty section.

Why it matters: LE-1 is the cheap, universal anti-rubber-stamp measure; an empty
section that passes validation gives false assurance that refutation was attempted.

Concrete fix: require at least one non-heading, non-blank line under the section (best-
effort heuristic, same spirit as the severity-heading scan), e.g. scan for a line after
`## Refutation attempts` and before the next `## ` that has non-whitespace content. If
that is deemed too strict, document the presence-only limitation explicitly so reviewers
know the work-recording is advisory, not enforced.

### [MINOR] LE-3 model-diversity warning/escalation in `OpenReviewRound` is not tested

The acceptance criterion for LE-3 is observable: "with a same-model reviewer set,
opening review prints a warning; with `require_model_diversity: true` it escalates;
differing models -> silent." `TestReviewersShareImplementerModel` unit-tests the pure
helper `reviewersShareImplementerModel`, but no test exercises `OpenReviewRound`
itself (internal/app/driver_impl.go:147-155) — neither the stdout warning path nor the
`require_model_diversity` error-return path that `advanceImpl`/`advanceReview` escalate
on. The helper logic is correct, but the integration that the acceptance criterion
names is unverified.

Why it matters: a regression that, say, inverts the `same` condition or drops the
`ReadRequireModelDiversity` branch would not be caught.

Concrete fix: add a test that constructs a `driverImplOps` with a same-model roster and
asserts (a) `OpenReviewRound` writes the WARNING line to its `out` writer, and (b) with
`require_model_diversity: true` in the 00-prompt, `OpenReviewRound` returns a non-nil
error. A differing-model roster asserts silence (no warning text, no error).

### [NIT] `reviewRoundHasFindings` excludes `consensus.md`/`_index.md` by filename, creating a reviewer-id asymmetry

The scan skips any entry named exactly `consensus.md` or `_index.md`
(impl.go:318). `ReviewRoundComplete` (driver_impl.go:218-229) does NOT skip them — it
validates `<reviewer>.md` for every reviewer. So a reviewer whose agent ID is the
literal string `consensus` would have its file counted toward round-completeness but
excluded from the finding scan, letting a finding in that file evade the veto.

Why it matters: practically nil (no agent is named `consensus` or `_index`), but it is
a real asymmetry between the two passes over the same directory.

Concrete fix: exclude only the archived `consensus.md` by path relative to the round
dir is already what the code intends; the risk is the bare-name match. Either accept
this (documented) or have the scan skip a file only when it is the archived consensus
artifact (e.g. the archived file is always at `review/round-NN/consensus.md`, which is
unavoidable), and rely on agent IDs never colliding with reserved names. Low priority.

## Open questions

- strict_gate + no `require_model_diversity` + a same-model roster: the strict closing
  round can still be rubber-stamped "no findings" by same-model checkers, after which
  the drafter certifies it clean and the (empty) scan agrees -> Complete(). strict_gate
  guarantees a fresh zero-finding round but cannot detect a rubber-stamp; that is LE-3's
  job, and LE-3 defaults to WARN, not gate. So an idea that opts into strict_gate but
  NOT require_model_diversity retains a residual false-green path via same-model
  reviewers. Is that residual accepted, or should strict_gate imply at least a
  model-diversity warning that names the risk on the closing round?

- strict_gate and the fix-up loop share the single `MaxFixupCycles` budget (both use
  `round >= MaxFixupCycles` / `cycle := round`). With the default of 3, an idea that
  spent round 1 on a fix-up gets only ~1-2 strict closing attempts before escalating.
  This satisfies FINAL.md's "bounded by MaxFixupCycles" (termination is guaranteed),
  but the closing-round budget can be consumed by fix-ups. Should the strict-close
  bound be a separate, smaller counter, or is the shared budget the intended
  trade-off?

- The verification-honesty idea's own `parley-deck/ideas/verification-honesty/00-prompt.md`
  sets `checks:` but does NOT set `strict_gate: true` on itself. The idea that
  introduces strict_gate does not dogfood it. This is an authoring choice, not a code
  defect, but worth confirming it is intentional (e.g. strict_gate deemed unnecessary
  for a protocol-only change with no runtime state of its own to verify beyond the
  test suite already covered by `checks:`).
