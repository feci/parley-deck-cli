---
agent: codex-1
idea: addon-manifest-coverage
review-round: 5
date: 2026-08-04
reviewed-commit: e4ee4d2
---

## Leak measurement (isolated TMPDIR)

The test-fixture leak is closed. I measured both commits from clean `git archive` exports,
with an initially empty isolated `TMPDIR`. I did not copy `node_modules`; each archive used a
read-only symlink to the repository's installed dependencies.

- Control, `065985e`: the full `npm test` command exited 0 (385/385 Node tests, 54/54 Python
  tests, all six manifest checks) and left exactly 18 fixture roots: 2
  `parley-bidding-test-*`, 8 `parley-tracker-claim-*`, and 8
  `parley-tracker-validate-*`.
- Reviewed commit, `e4ee4d2`: the same full command exited 0 with the same test/check counts and
  left 0 fixture roots. With `NODE_DISABLE_COMPILE_CACHE=1`, the literal count of every entry
  remaining in `TMPDIR` was also 0.
- On this host's Node 26.5.1 with its compile cache left at the runtime default, the literal
  count is 1: `node-compile-cache`. The count of `parley-*` fixtures is still 0. That directory
  is Node's intentional runtime cache, not a suite fixture or a regression in this fix.

This confirms the before/after result that matters: 18 leaked test roots at `065985e`, 0 at
`e4ee4d2`.

## Findings

No findings.

`forceRemove` is the same function in all three edited top-level test files. Every target passed
to it is a tracked `mkdtempSync` root. It uses `lstatSync`, does not `chmod` symlinks, and the
recursive removal unlinks rather than traverses symlinks. Using the exact function extracted
from `e4ee4d2`, I tested both a child symlink to an external sentinel and a tracked-root path
replaced by a symlink; the external directory and file retained their bytes and modes. The
tracked roots were removed. There is no reachable path in these tests that chmods or deletes
outside the tracked root.

The cleanup cannot mask a genuine test failure or change its exit code. It runs synchronously
on process exit, catches cleanup failures, and never assigns `process.exitCode`. An intentional
`node:test` failure with the exact handler still exited 1 while its hardened fixture was
removed.

Editing the two shipped `parley-tracker/bin/*.test.js` files is acceptable. They were already
part of the shipped payload; the change only adds cleanup for their own ordinary temp fixtures
and does not alter the tracker executables. Because shipped test files are manifest-covered,
regeneration was required. I independently recomputed all eight payload-file hashes and the
aggregate as
`sha256:07d9826373e4be3d2f393a8a56da616945fbc2e4e0838938827ceba7e85dfdd5`, exactly matching
`parley-addon.json`. The all-skills `--check` passed. `git grep` at `e4ee4d2` found no reference
to the old `sha256:2b178723...` aggregate and no separate tracker-aggregate assertion. The base
commit `23a9856` had no tracker manifest, so this does not revise a pre-idea published manifest
contract.

The remaining cycle-4 technical notes reproduce: the stock macOS PATH exposes Python 3.9.6 and
the Node suite passes 385/385 there; all six manifests check clean; and the `parley-bidding`
aggregate is the same at `714712f` and `e4ee4d2`
(`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`). The self-reported
timing of the uncommitted edits during round 4 is a historical process event and cannot be
reconstructed from commit contents; I could verify that `e4ee4d2` has `065985e` as its sole
parent and that the committed delta is the six files described. This is not a release concern.

## Ready to release?

Yes. I would ship `e4ee4d2` and release 2.2.0 on all five channels. No finding or unreproduced
technical claim stops release.

## Verdict

NO FINDINGS — ready to complete and release
