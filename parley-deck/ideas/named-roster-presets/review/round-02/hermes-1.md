---
agent: hermes-1
idea: named-roster-presets
review-round: 2
date: 2026-07-04
reviewed-commit: 8ab0e8c
responding-to: [hermes-1/review/round-01, codex-1/review/round-01]
---

## Position changes since prior review round

Two of my round-01 findings were in scope for fix-up cycle 1, and both are now
resolved at commit 8ab0e8c. I re-read the two changed sites and re-ran the tests
rather than trusting the diff summary.

- CRITICAL (provenance inside frontmatter fence): FIXED. The format string in
  `internal/protocol/workspace.go:151-171` now interpolates `provenanceBlock`
  into the body — after the closing `---\n\n` and before `## Problem / idea`
  (line 159, the `%s` immediately following `---\n\n`). The `track:` line
  (workspace.go:140-143) stays inside the fence as a real canonical key, which
  is correct. `ReadFrontmatter` (workspace.go:314-343) returns at the closing
  `---`, so the comment is never ingested. The inline comment at
  workspace.go:144-145 documents *why* it is placed below the fence, which is
  helpful. I withdraw this finding.

- MAJOR (preset list silently drops warnings when §2 unparseable): FIXED.
  `internal/app/preset.go:50-53` now prints
  `⚠ could not read the §2 roster (COOPERATION.md) — stale-member validation skipped.`
  to stdout when `ReadRosterIDs` returns `ok=false`. The per-preset warning
  block (preset.go:71-84) remains gated behind `if ok`, which is the right
  structure: the one-line notice explains *why* no per-preset flags appear, so
  the absence of ⚠ marks is no longer mistaken for a clean bill of health. I
  withdraw this finding.

Verification:
- `go test ./internal/protocol ./internal/app` → ok (protocol cached, app
  9.669s). `TestCreateIdeaFullProvenanceOutsideFence`
  (roster_test.go:57-84) asserts (a) no frontmatter key starts with `<!-` or
  contains `roster-preset`, (b) `track == "fast"`, and (c) the comment is
  present in the file body — all three hold.

No position change on the findings that were *not* in the fix-up scope (see
Updated findings / Open questions).

## Updated findings

Zero remaining on the two in-scope items above.

The following round-01 findings were NOT in scope for fix-up cycle 1 and are
unchanged at 8ab0e8c (verified by re-reading; not regressions, just not yet
addressed):

- MINOR — membership fail-open on unparseable §2 in `parley run`
  (app.go:1789-1791, `rosterIDs=nil` then continues). Still as documented in
  IMPLEMENTATION.md. This is the same behavior codex-1 flagged as their MAJOR
  ("Preset expansion fails open when the §2 roster cannot be parsed"). Still a
  divergence from consensus item 6's "no silent fallback" for the membership
  dimension. Not a regression; deferred.

- MINOR — `--track` / `--preset` / `parley preset list` absent from
  `printUsage` and the `run` usage string (app.go:142-241, 1742). Still
  undiscoverable from `parley help`. Not a regression; deferred.

- MINOR — track-default expansion prints the preset name/source but no
  override hint (app.go:1800). Consensus item 8's `track=standard → preset
  'trio'; override with --preset/--participants` hint still missing. Not a
  regression; deferred.

- NIT — `ReadRosterIDs` inactive detection is an unscoped substring match on
  the whole row (roster.go:62). Still fragile to a workspace dir containing
  "inactive". Not a regression; deferred.

- NIT — unknown-preset error lists known names but no closest-match hint and
  no layers-searched trace (roster.go:99). Consensus item 6's closest-match
  hint still absent. Not a regression; deferred.

- (From codex-1/round-01, MINOR) CLI-level preset acceptance criteria —
  `--preset`+`--participants` hard-error, track-default expansion, provenance
  in 00-prompt.md, fail-closed on unparseable §2 — are not covered by
  `internal/app` tests. Still no app-level integration tests for these.
  Not a regression; deferred.

## Open questions

1. **Run-path fail-closed on unparseable §2 (was my MINOR, codex-1's MAJOR).**
   Both round-01 reviews independently flagged that `parley run` silently skips
   membership validation when `ReadRosterIDs` returns `ok=false`. Fix-up cycle 1
   addressed the `preset list` UX surface (my MAJOR) but not the `run` path.
   Should round-02/03 hard-stop `parley run --preset ...` on `!ok` before
   `ResolveRoster`, or is the documented narrow fail-open accepted as the final
   behavior? This is the largest remaining gap versus consensus item 6.

2. **Deferred MINORs/NITs — triage for round-02.** The usage-text gaps
   (`--track`/`--preset`/`preset list` invisible in `parley help`) and the
   missing track-default override hint are both cheap, user-facing wins. Are
   they in scope for the next cycle, or is this idea considered feature-complete
   with only the safety-critical items (provenance, §2 fail-closed) eligible for
   further change?

3. **App-level integration tests (codex-1 MINOR).** The resolver/protocol
   tests are solid, but the CLI call-order invariants in `runTask`
   (`--preset`+`--participants` hard-error, track-default before preflight,
   provenance written) have no `internal/app` test guard. Should round-02 add
   these, given they are the most regression-prone integration points?
