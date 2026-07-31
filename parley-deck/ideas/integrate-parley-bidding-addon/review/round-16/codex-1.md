---
idea: integrate-parley-bidding-addon
review-round: 16
agent: codex-1
date: 2026-07-31
reviewed-commit: a49d68f
---

## Verdict

BLOCK

## Outstanding findings — closed or not

The exact round-15 regressions are closed. Simple physical ancestor/descendant destinations
are refused in both directions; dry-run and real uninstall now construct every planning-time
disposition through `removeFleetAtomically`; the Python runner uses `readManifest`; and the
case-only test now puts two targets into one plan.

The release-level no-partial-fleet guarantee is not closed over the full path-resolution space.
The MAJOR below changes one destination's relationship to another during this same process,
after the planning-time comparison. The uninstall parity fix itself is complete for every
planning disposition I traced: alias, marker and ownership blockers match; absent units remain
`missing`; pending units become `skipped` behind a blocker; and the only successful verb change
is the intended `remove` to `removed` transition.

I found no fifth direct reader of `parley-addon.json`: installer, generator and Python runner
all reach the shared module, and the generator's byte comparison is guarded by `hasManifest`.
The MINOR below is instead a consumer of the manifest's runtime policy, plus the interpreter's
reported version, that still fails open after the file-level reader was unified.

Concurrent-installer isolation remains the unanimously recorded follow-up from round 14. The
MAJOR below uses one invocation, one process and no concurrent writer, so it does not re-open or
depend on that disposition.

## New findings

### [MAJOR] A later commit can invalidate an earlier destination's symlink resolution

`resolvedDestination` records only the final realpath
(`lib/installer.js:1374-1396`), and `aliasedDestinations` compares those final strings only at
planning time (`lib/installer.js:1434-1457`). It does not record that resolving one destination
may traverse a symlink stored inside another planned destination. Replacing that other
destination then changes the first destination's meaning after it has committed.

I first installed a managed Kimi core at
`$KIMI_CODE_HOME/skills/parley-deck`, added an ordinary `redirect` symlink inside that managed
tree pointing to a separate directory, and set `CODEX_HOME` to the symlink. The planned Codex
destination therefore resolved outside the Kimi destination, so neither equality nor physical
containment fired. With `--target all --include-undetected --no-addons`, without `--force`, the
single command returned:

- top-level `ok: true`;
- Codex core `installed` and Kimi core `replaced`;
- the reported Codex destination absent after the command;
- an orphaned physical Codex tree at the symlink's former target, with marker `target: codex`;
- the surviving Kimi destination with marker `target: kimi`.

Codex commits first. The later Kimi commit renames and replaces the directory that held the
`redirect` symlink, so Codex's logical destination disappears while the transaction still
reports success. This is the same single-process false-success/partial-fleet state B5 and the
2.1.0 changelog claim cannot occur.

The removal transaction has the same blind spot. With both destinations managed, it first
quarantined Codex through the symlink and later quarantined the Kimi tree containing that link.
Phase B then called `rmSync(..., {force:true})` through the now-broken logical path, treated the
missing path as successful cleanup, returned top-level `ok: true` with both units `removed` and
no warning, but left `.parley-deck.<pid>.<timestamp>.removing` under the external target.

The fix must make destination resolution stable for the whole transaction, not only compare
final paths once. For example, snapshot and revalidate every destination's resolved path after
staging and again after all commits/quarantines but before cleanup, rolling the fleet back in
reverse on drift; also reject a plan when one destination's resolution chain traverses another
planned destination. Regress the vulnerable order (dependent target first, symlink-containing
ancestor later) for both install and uninstall.

### [MINOR] The Python gate fails open on malformed runtime policy and version output

The runner now reaches `readManifest`, but `declaredPythonFloor` turns any present
`runtime.python` value outside `^>=\s*(\d+)\.(\d+)$` into `null`, which means no floor
(`scripts/run-python-tests.js:53-61`). In a temporary `a49d68f` archive I changed only the
manifest value from `">=3.10"` to `"3.10"`. Both `npm run test:python` and
`npm run manifest:check` exited 0; `readManifest` and `verifyPayload` also returned `ok: true`.
The installed `doctor` path, by contrast, calls that syntax unsupported
(`lib/installer.js:2040-2049`).

The interpreter output has a second fail-open parser: `resolveInterpreter` splits on `.` and
uses `Number` without validating the complete output (`scripts/run-python-tests.js:63-80`). I
put a `python3` shim first on `PATH` that printed `4.not-a-version` for the version probe and
delegated test files to `/usr/bin/python3` 3.9.6. The runner exited 0 after all 54 tests and
printed `python 4.NaN: 54 tests OK`, while the health probe's anchored parser correctly rejects
this class of output.

The current shipped manifest is correct, and the full Node suite also hard-codes its expected
floor, so this is MINOR rather than a current payload failure. It is still a false-green in two
commands whose stated job is to enforce that policy. Validate a present runtime requirement in
the shared manifest parser (absence may remain optional), require a fully anchored interpreter
version, and add negative regressions for both malformed values.

### [NIT] The removed uninstall preflight left its helper dead

Cycle 19 deleted the only call site of `preflightUninstallUnit`, but the function remains at
`lib/installer.js:689-719`. `rg` finds only its definition; it is not exported. Delete the dead
duplicate so the single-path result construction is also true in the code structure. Under this
idea's strict gate, objective dead code remains a finding.

## Release judgement

Not releasable as 2.1.0 at `a49d68f`. The one indispensable release change is to make the
single-process fleet transaction detect and roll back destination-resolution relationships
changed by its own commits or quarantines; a planning-time final-path containment check is not
enough. The strict gate also requires the measured MINOR and NIT to close before a clean round.

## What I verified

- Read the live cooperation protocol, `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through
  cycle 19, all three round-15 reviews, the relevant inbox records, and both requested diffs.
  `git diff --check` is clean for `26478e9..a49d68f` and `49fc3ec..a49d68f`.
- Built a code-only Graphify map from an `a49d68f` archive under `/tmp` (155 nodes, 301 edges),
  used it to trace planner/mutation and manifest-consumer relationships, then verified every
  relevant edge directly in source. The repository graph and working tree were not modified.
- Ran `npm test` from an exact archive: **355/355 Node tests**, **54/54 Python tests** on Python
  3.14, and the 47-file manifest check with aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`. Ran the Node suite
  again with `/usr/bin/python3` 3.9.6 first on `PATH`: **355/355**.
- Ran `npm pack --dry-run --json --ignore-scripts` with an isolated npm cache: **202 files**,
  **48** under `skills/parley-bidding/`, and no `__pycache__`, `.pyc` or `.pyo`. All four
  adapters validate, all **16** bidding JSON files parse, and direct Python checks left no cache.
- Confirmed `skills/parley-bidding/` has no diff from `714712f` and no git symlinks. The
  read-only source comparison is still exactly 48 source files versus 48 integrated files: one
  dropped `.gitignore`, one added manifest and the recorded nine content differences; the
  source contains no generated Python cache.
- Audited every `parley-addon.json` reference outside the payload and tests. There is no fifth
  direct reader bypassing the regular-file predicate. Reproduced the malformed runtime-policy
  and malformed interpreter-output arms independently in temporary archives/shims.
- Reproduced both destination-resolution arms in fresh temporary homes: install returned false
  success with an absent logical Codex install, and uninstall silently retained a quarantine
  directory. Both are single-process, use supported `CODEX_HOME`/`KIMI_CODE_HOME` inputs, and
  require no concurrent writer; the install arm also requires no `--force`.
- The reviewed repository remained clean throughout; no file under `skills/parley-bidding/`
  was edited.
