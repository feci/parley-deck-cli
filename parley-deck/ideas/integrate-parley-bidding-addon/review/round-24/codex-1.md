---
idea: integrate-parley-bidding-addon
review-round: 24
agent: codex-1
date: 2026-07-31
reviewed-commit: e274eb8
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

Both round-23 findings are closed.

- **Asymmetric injected path semantics — closed.** `rawTargetArithmetic` now derives separator
  semantics solely from `impl.sep`. The surviving regression temporarily makes the host report
  `win32` and asserts both directions: injected `path.posix` keeps `we\ird` as one filename
  component, while injected `path.win32` still splits a Windows path correctly. With the current
  test overlaid on a `git archive` of `7e8ccec`, the test fails at the POSIX assertion with
  actual `['we', 'ird', 'x']`; at `e274eb8` it passes.
- **Duplicated/dead tests and the false cycle-26 record — closed.** Each of the three collision
  tests occurs once, the obsolete absolute-target test is gone, and the unused scaffolding it
  carried is gone with it. The cycle-26 entry now explicitly records that its original removal
  claim was false; cycle 27 records the actual cleanup. The suite count is consequently 368,
  exactly four below cycle 26's 372.

All earlier findings remain closed. Concurrent-installer isolation remains the unanimously
recorded follow-up from round 14 and is not an outstanding fix for this release.

## Position on the gate

**1 (correct).** Both mutation paths reach the same fleet-wide collision gate before a write or
deletion. The gate retains the required physical-identity, containment, and resolution-crossing
arms, and the accumulated collision regressions remain green. Narrowing it would still
contradict the documented fleet-wide atomicity guarantee for supported symlinked runtime homes.

## New findings

None.

## Release judgement

Releasable as 2.1.0. No change is required before release. The payload remains byte-identical to
`714712f`, the manifest aggregate is unchanged, the cycle-27 repair discriminates against the
round-23 regression, and the fresh full-scope review found no issue of any severity.

## What I verified

- Confirmed HEAD is `e274eb841722355a462999b3bf67497af34c2648`. Read the live protocol,
  `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through cycle 27, all round-23 reviews, and the
  already-complete Hermes/Kimi round-24 artifacts after finishing my independent review.
- Read `git diff 7e8ccec..e274eb8` and the six-file full-scope
  `git diff 49fc3ec..e274eb8`; `git diff --check` passes. A disposable Graphify AST map of the
  five changed JavaScript files contained 178 nodes and 330 edges and confirmed that install and
  uninstall both reach `aliasedDestinations`, which reaches
  `resolutionTouchpoints -> walkRawTarget -> rawTargetArithmetic` before their mutation phases.
- `npm test`: **368 Node tests, 368 pass, 0 fail**; Python **54/54** across seven files on
  Python 3.14.6; manifest check **47 files**, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- Re-ran the complete Node leg with `/usr/bin/python3` 3.9.6 first on `PATH`: **368 pass, 0
  fail**. The Python runner then refuses 3.9.6 with exit 1 and
  `python3 is 3.9, but the add-on declares >=3.10`, as designed.
- Ran a 14-test superset of the accumulated destination identity, nesting, firmlink, symlink
  chain, raw-target, POSIX-backslash, drive-root, and UNC-root regressions: **14/14 pass**.
- Overlaid the current bidding test on a `git archive` of `7e8ccec` and ran the real
  root-arithmetic regression: **1 fail** at the old commit with the exact host-overrides-injection
  assertion; the same test passes at `e274eb8`.
- Counted **109** test declarations in `test/bidding-addon.test.js`, with no duplicate titles.
  The obsolete absolute-target title has zero matches.
- `npm pack --dry-run --json` ran `prepack` and reported the same 47-file manifest aggregate;
  the package contains **202 files**, **48** under `skills/parley-bidding/`, and no `.pyc`,
  `.pyo`, or `__pycache__` entry.
- `git diff 714712f..e274eb8 -- skills/parley-bidding bin` is empty. A fresh read-only
  source comparison found 48 source files and 48 integrated files: `.gitignore` is the sole
  dropped path, `parley-addon.json` the sole added path, and the documented nine files differ;
  neither tree contains a symlink or Python cache.
- The adapter validator accepts all four platform profiles. All 16 shipped JSON files parse;
  all four schemas declare draft 2020-12 and the ratified `example.invalid/parley-bidding`
  identities. The deterministic Python scripts import no network or browser library, and the
  payload scan found no bundled credential, BYTE, or customer material.
- Removed the old-commit archive, Graphify archive, and disposable npm caches. The implementation
  repository remains clean; I wrote only this review artifact in the coordination repository.
