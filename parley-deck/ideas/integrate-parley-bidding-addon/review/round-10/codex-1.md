---
idea: integrate-parley-bidding-addon
review-round: 10
agent: codex-1
date: 2026-07-31
reviewed-commit: 9ed2081
---

## Verdict

BLOCK

## Outstanding findings — closed or not

- The specific round-8 and cycles 10–13 failures are closed. The installer now builds and preflights the complete target-by-unit plan, treats dangling destination entries as present on install/status/uninstall paths, rejects a non-directory destination ancestor even under `--force`, and requires write-and-search permission on the nearest existing ancestor. The four cycle 12/13 regressions are present in the passing 325-test suite.
- I did not reopen B1, B2, B4, B6, or B7. The consent boundary, exposure wording, separate Python runtime health, and channel identity remain consistent with `FINAL.md` in the full tree and packaged artifact.
- kimi-1's dangling-link discovery NIT remains a legitimate non-blocking follow-up. I agree that changing whether an unselected add-on is discovered through a dangling link is outside this release's ratified discovery semantics.
- The universal-health result for five add-ons without manifests remains the recorded follow-up: one `valid-unmanaged` and five `malformed`. It is unchanged and does not excuse either new false-green/partial-fleet path below.

## New findings

### [MAJOR] Replacement cleanup can fail after commit and leave a partial fleet

`destinationAncestorObstacle` checks only whether the destination's nearest existing ancestor can be entered and written (`lib/installer.js:1820`). It does not establish that an existing destination can be recursively removed after replacement. `copyPayloadAtomically` then renames the old destination to a backup, commits the new destination, and only afterward recursively removes the backup (`lib/installer.js:1356`). If that cleanup fails, its catch path retries the same removal and propagates the error (`lib/installer.js:1363`); it neither restores the old destination nor turns the already-committed replacement into a coherent success.

I reproduced this at exact `9ed2081` as uid 501 by placing an existing final-plan destination, `~/.aionrs/skills/parley-worktrees`, with a mode-000 nested directory and sentinel. In both a managed unforced reinstall and a foreign `--force` replacement, the complete plan preflight passed. The command reported 83 successful units and failed on unit 84 with `ENOTEMPTY`; the new unit was nevertheless installed and a hidden `.bak` tree remained. This is a predictable destination-state blocker, not an unexpected mid-copy fault, and it violates B5's no-partial-fleet guarantee.

Before release, replacement/removal feasibility must be covered by fleet-wide preflight, or the replacement transaction must make post-commit backup cleanup unable to convert a committed replacement into a failed partial fleet. A regression should put the obstacle in the last target/unit and cover both managed unforced and foreign forced replacement, asserting either zero writes on preflight failure or coherent success with no backup residue.

### [MAJOR] Removing only `markerSchema` silently disables manifest validation

`manifestProblems` treats every installer-owned marker with an absent `markerSchema` as legacy and immediately returns healthy (`lib/installer.js:1698`). The released 2.0.0 legacy shape lacked both `markerSchema` and `manifest`; a current marker that still contains its manifest but has lost only the schema field is not that legacy shape.

I reproduced the false green at exact `9ed2081`: after a clean `parley-bidding` install, I removed only `markerSchema` from the marker, left its current `manifest` object intact, and modified a manifest-covered payload file. `doctor` returned `ok: true`, status `valid`, and no problems while still displaying the marker manifest. One missing metadata field therefore downgrades a current managed install from full byte validation to no manifest validation, contrary to B3 and the stated no-silent-downgrade rule.

The legacy exemption must be limited to the exact supported legacy shape. A manifest-bearing marker without a valid schema must be unhealthy. Add a regression that removes only `markerSchema` from a current marker and verifies a non-green result, while retaining compatibility for the genuine 2.0.0 marker without either field.

## Release judgement

No. Commit `9ed2081` is not releasable as 2.1.0. The release condition is to close both remaining installer-integrity gaps: replacement cleanup must not permit a partial fleet, and a manifest-bearing schema-less marker must not silently bypass validation.

## What I verified

- Kept the source worktree clean and stable at `9ed2081ff241f36a0b5b96e930be116327ac6fdc`; older commits and destructive probes were exercised only in isolated `git archive` copies.
- Read the complete feature change from 2.0.0, both requested diffs (`dcd200e..9ed2081` and `49fc3ec..9ed2081`), `FINAL.md`, all implementation cycles, the consensus, and the relevant prior review artifacts. Round 9 was treated as void, not as a signoff.
- Ran the complete suite at the reviewed commit: 325 Node tests passed, the Python leg passed 54/54 across seven files, and the 47-file add-on manifest verified at aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- Re-ran the 325 Node tests with `/usr/bin/python3` 3.9.6 first on `PATH`; all passed. The complete default run used Homebrew Python 3.14.6.
- Ran `npm pack --dry-run --json`: prepack manifest verification passed and the package contained 202 entries, including the complete 48-file bidding add-on payload without caches, symlinks, nested `.gitignore`, or source repository metadata.
- Built the current portable binary and used it to install only `parley-bidding` into an isolated home. Portable `doctor` reported core and bidding healthy, including the expected manifest aggregate and Python runtime.
- Ran adapter validation directly (four adapters), parsed all 16 JSON files and four schemas, compared the integrated payload with its source inventory, and inspected the intentional renames/consent edits and absence of portal/network automation.
- Independently reproduced both new findings above at exact `9ed2081`, including both forced and unforced arms of the replacement failure and the manifest-bearing schema-less marker false green.
