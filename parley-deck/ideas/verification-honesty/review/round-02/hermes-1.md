---
agent: hermes-1
idea: verification-honesty
review-round: 2
date: 2026-06-24
---

## Summary

Round-02 re-review of the fix-up diff b6a0624..HEAD (2dd5782). The implementer applied
all eight agreed fixes from review/consensus.md. My round-01 MAJOR (F1 — finding-scan
fail-OPEN on read errors) is now fail-closed, and F2/F3/F5/F7 hold. I attempted to break
the strict_gate close logic from a new angle and could not find a remaining fail-open path.
Zero findings; the fixes hold.

Build: `go build ./...` clean. Tests: `go test -count=1 ./...` all green (driver, app,
runner, and 21 other packages).

## Refutation attempts

F1 (my round-01 MAJOR) — fail-closed verification.
  reviewRoundHasFindings (internal/driver/impl.go:312-336) now:
    - ReadDir error → `return !errors.Is(err, fs.ErrNotExist)`. An absent dir
      (ErrNotExist) does not veto (nothing to scan — ReviewRoundComplete guards
      presence). Any OTHER ReadDir error (unreadable dir, path-is-a-file, permission)
      returns true → veto → escalate. Confirmed by TestReviewRoundHasFindingsFailsClosed,
      which creates round-05 as a regular file and asserts reviewRoundHasFindings returns
      true. The test passed under my run.
    - Per-file ReadFile error → `return true` (impl.go:329). A file ReadDir confirmed
      present but now unreadable (e.g. mode 000) vetoes. Logic is correct and fail-closed.
      No dedicated test exercises this exact branch (the F1 test only covers ReadDir),
      but the code path is unambiguous: ReadFile can only fail after ReadDir listed the
      entry, and any such failure → true. This is a coverage gap, not a logic defect (see
      Open questions).

F2 — case/whitespace-tolerant severity tag. scanHasRealFinding (impl.go:342-366) now
  uppercases the tag via strings.ToUpper(strings.TrimSpace(rest[1:closeIdx])) and tolerates
  spacing after ### via strings.TrimPrefix + TrimSpace. Verified: "### [critical]",
  "###   [MAJOR]", "### [CrItIcAl]" all match. TestScanHasRealFinding covers lowercase and
  extra-spacing cases; passed.

F3 — empty-title heading counts. The title check is now `strings.TrimSpace(rest[closeIdx+1:])
  != "<title>"` — an empty title (finding text on the next line) returns true; only the
  literal placeholder <title> is ignored. TestScanHasRealFinding "empty title on line"
  case passed.

F5 — heading-vs-substring validation. ValidateReviewArtifact (phase58.go:413-441) now uses
  hasHeadingLine (exact trimmed-line match for "## Findings") and hasNonEmptySection (a
  "## Refutation attempts" heading followed by ≥1 non-blank line before the next "## "
  heading). A substring mention in prose and an empty section both reject. Tests in
  phase58_le_test.go cover both; passed.

F7 — strict_gate close-field absence escalates. DraftReviewConsensus (app/driver_impl.go:
  263-289) now reads strict_gate, and under strict requires non-empty strict_gate_clean and
  closing_review_round in the drafted consensus frontmatter; absent → error returned →
  advanceReview escalates (DraftReviewConsensus error → "draft review consensus: %w"). This
  prevents the spin-to-MaxFixupCycles. Logic traced end-to-end; correct.

Strict_gate close-logic fail-open hunt (the assigned refutation target).
  The close path (impl.go:193-234) when OutstandingAgreedFixes==0 && StrictGate:
    certifiedClean := rs.StrictGateClean && rs.ClosingReviewRound == round
    if certifiedClean && reviewRoundHasFindings(ideaDir, round) → escalate (veto)
    if !certifiedClean → open another round (or escalate at bound)
    else → Complete
  I traced every input that feeds certifiedClean and the scan for a fail-open:
    - ClosingReviewRound parsing (app/driver_impl.go:308-313): a malformed/non-integer/
      negative value silently leaves closingRound=0. Since round≥1, certifiedClean becomes
      false → opens another round → fail-closed. Good.
    - StrictGateClean parsing (app/driver_impl.go:307): defaults false on absent/bad →
      certifiedClean false → fail-closed.
    - Scan scope (impl.go:324): excludes _index.md, consensus.md, subdirs, non-.md. The
      veto target is raw reviewer findings (e.g. codex.md, hermes-1.md), which are exactly
      the non-excluded .md files. A finding hidden only in consensus.md is a drafter-
      managed file, not a raw reviewer finding — out of veto scope by design.
    - roundLabel consistency: reviewRoundHasFindings (driver pkg) uses roundLabel =
      round-%02d (cursor.go:259); the app layer writes reviewer files to roundDirLabel =
      round-%02d (app/driver_impl.go:379). Same format → scan reads the right directory.
    - The only way to reach Complete is certifiedClean==true AND the scan found no real
      finding in any readable reviewer .md. A ReadDir or ReadFile error now vetoes (F1).
      I could not construct a path where a real finding evades the scan AND Complete is
      reached. The close logic is fail-closed.

## Findings

None. All agreed fixes (F1-F8) are correctly implemented, build clean, and pass their
tests. The strict_gate close logic has no remaining fail-open path I could find.

## Open questions

- The per-file ReadFile-error branch of F1 (impl.go:328-329) has no dedicated test. The
  existing TestReviewRoundHasFindingsFailsClosed covers only the ReadDir non-NotExist
  path. The logic is correct (return true on any ReadFile failure for a listed file), but
  a test that chmods a reviewer file to 000 and asserts a veto would close the coverage
  gap. Not a defect — flagging for completeness only.
