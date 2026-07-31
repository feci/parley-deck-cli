---
idea: integrate-parley-bidding-addon
review-round: 24
agent: kimi-1
date: 2026-07-31
reviewed-commit: e274eb8
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

None open. Codex-1's two round-23 findings — both cycle-26 self-inflicted — are closed in
cycle 27, and I re-measured both closures on this tree rather than taking the record's word:

1. **Asymmetric injection — closed.** `rawTargetArithmetic` now takes its separator semantics
   from `impl.sep` alone (`lib/installer.js:1414`). I extracted `lib/` + `package.json` from
   `7e8ccec` via `git archive` into a temp dir, redefined `process.platform` to `"win32"`, and
   called the old helper with `path.posix` injected on `"we\\ird/x"`: it returned
   `["we","ird","x"]` — the host overrode the injection, exactly the arm codex-1 exercised. The
   same probe against `e274eb8` is green in the suite. The regression fails at `7e8ccec` and
   passes at `e274eb8`, as recorded. The new test asserts both directions and restores
   `process.platform` in a `finally`.
2. **The wrong edit — closed.** The three duplicated tests each survive exactly once
   (`test/bidding-addon.test.js:1697`, `:1746`, `:1784`); the dead absolute-target test appears
   zero times. Its coverage is not lost: the drive/UNC/relative anchoring assertions live on in
   the surviving arithmetic regression. Test count arithmetic is consistent: 372 (cycle 26)
   minus 3 duplicates minus 1 dead = 368, which is what the runner reports.

The false cycle-26 sentence is corrected in place in `IMPLEMENTATION.md` with an explicit note
that it was false — the right way to handle it, and the pattern being on the record is
appropriate.

## Position on the gate

1 (correct).

Cycle 27 does not change the gate's shape. On any host, `path.sep === "\\"` holds exactly when
the host is Windows, so the production call in `walkRawTarget` (which passes the host `path`)
computes identical separators before and after the change; only injected-impl semantics moved,
and only tests inject. The gate itself is untouched. My position from rounds 22 and 23 stands,
and the round-22 refusal to narrow (position 2) stands with it.

## New findings

None. I looked at the cycle-27 change adversarially: the `finally` restore is correct,
`impl.sep` discriminates `path.win32` from `path.posix` exactly, and no production path
reaches the helper with anything but the host `path`. A hypothetical custom impl with `sep`
`"\\"` but POSIX `parse` would mis-split — but nothing injects such a thing; that is a
thought experiment, not a finding.

Carried follow-up, unchanged and non-blocking: concurrent-installer isolation (recorded,
unanimous round-14 ruling).

## Release judgement

Releasable as 2.1.0. The payload has not moved since `714712f` (empty `git diff
714712f..e274eb8 -- skills/`, measured), the aggregate manifest digest is bit-identical to
every measurement since that commit, every leg of `npm test` is green, and the only defects
found in the last two rounds were in the repair process itself, not in what ships. In round 22
I said that with the round-21 arm measured shut, nothing stood between the tree and release.
Cycle 27 measured shut the last thing that appeared after that. Nothing does.

## What I verified

All commands run in `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill` at
`e274eb8` (`git status` clean, no stray files; no working-tree mutation at any point):

- `npm test` (background task, full output read): **368 tests, 0 fail**; python leg **54/54
  across 7 files on python 3.14**; `parley-bidding: ok (47 files,
  sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d)` — matches the
  record, unchanged since `714712f`.
- `git diff 714712f..e274eb8 --stat -- skills/` — **empty**; payload untouched since the
  manifest hash was first recorded. Full-scope diff since `49fc3ec` is confined to the known
  reviewed files: `CHANGELOG.md`, `lib/addon-manifest.js`, `lib/installer.js`,
  `scripts/run-python-tests.js`, and the three test files.
- `git diff 7e8ccec..e274eb8` — read in full: the `impl.sep` fix, the merged
  platform-redefinition regression, and the removal of exactly three duplicate tests plus the
  one dead test (156 lines deleted, 18 added).
- Grep of the test file: each surviving test name occurs once; the dead test's name occurs
  zero times.
- Discrimination probe (my own construction): `git archive 7e8ccec lib package.json` into a
  `mktemp -d` dir, redefined `process.platform` to `"win32"`, called the old
  `rawTargetArithmetic("/a/link", "we\\ird/x", path.posix)` — returned `["we","ird","x"]`
  instead of `["we\\ird","x"]`, so the new regression fails at `7e8ccec`; it passes at
  `e274eb8` via the suite. Temp dir removed with `rm -rf` (confirmed); no checkout, no
  bytecode written.
- 3.9.6 refusal: `PATH="/usr/bin:$PATH" node scripts/run-python-tests.js` with
  `/usr/bin/python3` 3.9.6 → `python tests: python3 is 3.9, but the add-on declares >=3.10`,
  exit 1. Fail-closed, by design.
- `lib/installer.js:1395-1534` read in context: `rawTargetArithmetic`, `walkRawTarget`
  (immediate expansion, landing-point continuation, depth cap), `resolutionTouchpoints`, and
  `aliasedDestinations` are coherent with the cycle-27 change; production separators are
  host-identical before and after.
- `IMPLEMENTATION.md` cycles 26–27 read: the false sentence is corrected in place with the
  note, and the cycle-27 measured numbers match what I measured here.
- `CHANGELOG.md` 2.1.0 entry read: matches shipped behavior (fleet-wide atomicity, manifest
  integrity as defect detection, python leg fail-closed).
