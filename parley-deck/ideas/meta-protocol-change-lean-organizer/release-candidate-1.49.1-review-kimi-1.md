---
idea: meta-protocol-change-lean-organizer
author: kimi-1
role: independent reviewer
artifact: release-candidate-1.49.1-review-kimi-1
date: 2026-09-24
subject: REVIEWABLE CLI 1.49.1 candidate metadata at 481fb65 (commits 3002782 + 481fb65, report release-candidate-1.49.1-zcode-1.md)
context: U1/U2 source already approved by both reviewers; hosted Linux/macOS acceptance PASS at 9134c7a (claude-1, run 36009912946); 519951e comment-only, token-identity verified by claude-1 — taken as given, not re-derived
---

# Independent review — CLI 1.49.1 release-candidate metadata — kimi-1

## VERDICT: **PASS** — candidate metadata is review-ready; release completion remains owner-gated (below)

Focused metadata checks only (per dispatch): no full suites, no publication steps.
Scope of check: version fields, changelog factual accuracy, drift vs `519951e`,
the two § 8 record corrections, tag/closed-artifact/skill integrity.

## Checks performed (all pass)

1. **Version fields 1.49.1.** `VERSION` = `1.49.1`; `internal/app/version.go`
   `const version = "1.49.1"`; CHANGELOG top heading `## 1.49.1 — 2026-09-24`;
   `go run ./cmd/parley --version` → `parley 1.49.1`. `gofmt` clean,
   `go vet ./internal/app/` ok, `TestVersionCommandPrintsSemanticVersion` PASS.
2. **No feature drift vs 519951e.** `git diff 519951e..HEAD` = exactly `VERSION`
   (1 line), `version.go` (1 line), `CHANGELOG.md` (+65/−0 at 3002782), the
   report (new). `3002782..481fb65` = only CHANGELOG (18 lines reworded, hunks
   at :20 and :41–43 — both inside the 1.49.1 section; 1.49.0 entry starts at
   :75) + the report. Zero product/test/workflow/skill change.
3. **1.49.0 entry untouched.** `git diff v1.49.0..HEAD -- CHANGELOG.md` has
   **0 deletion lines** — pure insertion above the frozen 1.49.0 section.
4. **Killed-unreaped/read-error correction is accurate.** CHANGELOG and report
   § 3 now state: reaped pid's ENOENT exits at once; killed-but-unreaped pid's
   cmdline reads empty with **no error** and consumes the full 100 ms
   `cmdlinePublishBound` before failing closed exactly as pre-poll —
   cancellation is bounded, not immediate. This matches the accepted C-1/K-1
   evidence. The stale phrase "dead processes are never polled" survives only
   as a quoted original inside § 8 — verified by grep, no live instance.
5. **Windows choice left to owner; no invented mandatory no-assets policy.**
   Changelog + report § 3/§ 6 mirror the organizer scope inbox
   (`codex-1-to-user_..._scope.md`) exactly: owner chooses
   defer-and-label-experimental/unvalidated **or** a separate reviewed
   Windows-portability track; neither choice inferred; omitting Windows assets
   is itself a deviation the owner must explicitly select; publication held
   pending the decision as a process hold, not a technical prohibition. The
   withdrawn "Windows assets must not ship on this evidence" obligation
   appears only as a § 8 quote.
6. **Changelog factual limits honest.** Every load-bearing figure
   cross-checks against the record: U1 0-of-16384 (run 35987916696); U2
   100 ms/1 ms bound, retry-zero-only, first-non-empty-wins, facets/strings
   byte-identical, darwin/Windows untouched; Windows failure figures
   (`syscall.Mkfifo` build failure, 14 fail / 16 ok, acp D1 drain guard
   firing) match the acceptance report § 6 verbatim; ubuntu 32 packages
   (31 ok + 1 no-test-files) with no skips in previously-failing tests and
   macOS green at run 36009912946 @ `9134c7a`; ~1.8× margin over 10,400
   measured spawns (claude-1's numbers, credited); ~0.7 % local reproduction;
   one-green-sample honesty; silent-exhaustion caveat + raise-the-bound
   operational rule. No invented claims found.
7. **Tag v1.49.0 unchanged; no v1.49.1 tag.** `git rev-parse v1.49.0` =
   `06e563e` locally and at `git ls-remote --tags origin`; no `v1.49.1` tag
   exists locally or remotely.
8. **Closed artifacts and skill unchanged.** `FINAL.md`, `IMPLEMENTATION.md`,
   `consensus.md`, `review/` — 0 diff lines across `v1.49.0..HEAD`. Skill
   worktree clean; `v2.13.0` = `8161e5e` (separate artifact, untouched).
9. **Working tree as declared.** Only the pre-existing untracked
   peer/organizer records (acceptance report, ledger, inbox, `runs/`,
   `organizer-usage.md`) — matching the report § 7 NOT-done list.

## Findings

None. No corrections requested.

## Remaining owner-dependent conditions (not review findings — release gates)

1. **Owner**: Windows scope decision per the scope inbox — defer with Windows
   labelled experimental/unvalidated, or authorize a separate reviewed
   Windows-portability track; omitting Windows assets requires the owner's
   explicit selection. Publication is held on this.
2. **Owner**: npm verification/publish of the separate skill 2.13.0 tarball.
3. **Organizer**: stage release assets from a pinned commit (`9134c7a` as
   tested, or push comment-only `519951e` first — token-identity verified).
4. **Organizer/owner only**: tag `v1.49.1`, channel builds, GitHub release,
   distribution. `v1.49.0`/`06e563e` never moves.

This review is candidate-readiness only: no source/metadata edits, no
commits/pushes/tags, no builds of distribution assets, no installs. It is not
release completion.
