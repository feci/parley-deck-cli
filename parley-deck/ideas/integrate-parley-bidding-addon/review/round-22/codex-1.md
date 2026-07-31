---
idea: integrate-parley-bidding-addon
review-round: 22
agent: codex-1
date: 2026-07-31
reviewed-commit: b1f43e4
---

## Verdict

BLOCK

## Outstanding findings — closed or not

- **Closed — the linked-ancestor / physical-parent MAJOR from round 21.** The outer
  `resolutionTouchpoints` walk now propagates the landing returned by `walkRawTarget`, and a
  relative target is anchored at the link's physical parent. I ran the cycle-25 regression
  against both trees: with the `b1f43e4` test file and the `381e639` installer it fails at
  `true !== false`; at `b1f43e4` it passes and records zero writes. The fix closes the arm that
  split round 21.
- **Closed — POSIX backslash semantics.** `walkRawTarget` now treats `\\` as a separator only
  on Windows. The test is accurately labelled as a pin rather than discrimination: it passes
  at both `381e639` and `b1f43e4` because the earlier implementation happened to record the
  containing component before tearing the name.
- **Closed — canonical implementation record.** `IMPLEMENTATION.md` now records cycles 24 and
  25 and names `b1f43e4` in its frontmatter.
- **Still open — [MINOR] the Windows absolute-target blocker has no test that exercises its
  fix.** This is the carried codex-1/kimi-1 round-21 finding, not a new finding. I copied the
  `b1f43e4` test file into a `git archive` of the known-broken `2b7ca3e` tree and ran:

  ```text
  node --test --test-name-pattern='an absolute raw target does not replay its root as a component' test/bidding-addon.test.js
  ```

  Observed state: **1/1 passes at `2b7ca3e`, and 1/1 passes at `b1f43e4`**. The test creates
  only a POSIX absolute link; it does not exercise a drive root, UNC root, or the Windows
  component walk. Its `seen` array and `record` callback remain unused and are discarded with
  `void`. The current arithmetic is correct by inspection, but `FINAL.md` requires every
  blocker to have a regression that fails without the fix, and the strict gate does not close
  over an outstanding maintainability finding. The user cost is on a shipped Windows channel:
  CI would remain green if the known drive/UNC false-success behavior returned, allowing a
  Windows fleet operation to report success while leaving a destination dangling and its files
  orphaned.

## Position on the gate

**1 — correct.**

The current gate has the right three arms: resolution crossing, identical physical identity,
and physical containment. Cycle 25 completes the crossing arm by carrying one physical cursor
through the whole destination and anchoring each relative link target at its physical parent.
All seven accumulated collision regressions pass. Narrowing would withdraw the documented
fleet-wide guarantee for symlinked runtime homes, so I do not choose position 2. The remaining
block is executable Windows coverage, not a defect in the gate's current shape.

## New findings

None.

## Release judgement

Not yet releasable as 2.1.0 under this idea's strict gate. One thing must change: make the
actual raw-target root/component arithmetic injectable (or extract it into an injectable
helper), then add drive and UNC assertions against that real helper which fail at `2b7ca3e`
and pass at the release tree. Use or remove the dead `seen`/`record` scaffolding in the same
change. No payload change is indicated.

## What I verified

- Read the live protocol, `00-prompt.md`, `FINAL.md`, the complete `IMPLEMENTATION.md` through
  cycle 25, `review/round-09/VOID.md`, and the round-20/21 review artifacts. The project is the
  protocol source; installer/reference drift is advisory and I made no sync write.
- Read the complete `381e639..b1f43e4` diff and the full six-file
  `49fc3ec..b1f43e4` diff. Both pass `git diff --check`. I traced both mutation paths through
  the destination gate, the install/rollback transaction, quarantine removal, marker trust
  boundary, manifest parser, Python runner, command guard, and release documentation.
- Built a disposable deterministic Graphify AST map of the six changed code files: **211
  nodes, 418 edges**. It confirmed that `installCommand` and `uninstallCommand` both reach
  `aliasedDestinations` through their fleet transaction functions. The semantic documentation
  pass failed with a connection error and was not used as evidence.
- Ran `npm test` from a `git archive` of `b1f43e4`, with the repository dependency tree linked
  read-only and all generated state under the disposable archive: **368 Node tests, 0 fail**;
  **54/54 Python tests** on Python 3.14; manifest check **47 files**, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- Re-ran the Node leg with `/usr/bin/python3` 3.9 first on `PATH`: all **368** passed. The
  Python runner then refused 3.9 as designed with exit 1 and `python3 is 3.9, but the add-on
  declares >=3.10`.
- Ran the seven accumulated collision-arm tests as one selection: **7/7 pass**. The cycle-25
  linked-ancestor test fails against `381e639` and passes at `b1f43e4`. The POSIX-backslash pin
  passes at both. The Windows absolute-target pin also passes at both the known-broken
  `2b7ca3e` installer and the current tree, which is why its prior MINOR remains open.
- Parsed `npm pack --dry-run --ignore-scripts --json`: **202 files**, **48** under
  `skills/parley-bidding/`, and no `__pycache__`, `.pyc`, or `.pyo` entry.
- Verified `skills/parley-bidding/` has no diff from `714712f`, no symlink, and no generated
  cache. README states six skills / five add-ons; deterministic Python scripts import no
  network or browser library.
- Removed every disposable graph, suite, and old-commit archive, including about **573 MB** of
  immutable-debris fixtures after clearing their test flags. The implementation repository was
  clean before and after; no file under `skills/parley-bidding/` was changed.
