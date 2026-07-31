---
idea: integrate-parley-bidding-addon
review-round: 20
agent: codex-1
date: 2026-07-31
reviewed-commit: 2b7ca3e
---

## Verdict

BLOCK

## Outstanding findings — closed or not

- The cycle-23 root fix closes the round-19 Windows defect for `identityChain` and the outer
  destination walk in `resolutionTouchpoints`. The real exported `splitAtRoot` produces correct
  drive, UNC, relative, bare-`C:`, and mixed-separator decompositions under `path.win32`.
  `walkRawTarget`, however, performs separate root arithmetic and still breaks absolute Windows
  link targets; that is the second new MAJOR below.
- The exact simple raw-`..` target from round 19 is closed. I also measured absolute POSIX
  targets, repeated separators, `..` past the root, and a target that revisits a directory; each
  was refused by the dependency gate before a write. The broader property is not closed because
  a symlink in an intermediate target component is recorded but never expanded; that is the first
  new MAJOR below.
- The 64-hop bound is effective. An instrumented two-link cycle performed exactly 64 raw-chain
  `readlinkSync` calls, returned promptly, and failed preflight with zero writes.
- The firmlink-with-existing-inner-parent pin is present and correctly labelled as a pin. The
  Linux case-sensitive normalization issue and the stale deleted-caller comment from round 19
  are closed.
- The earlier manifest, marker, runtime, selection, staging/commit/revert,
  quarantine/rollback/cleanup, dry-run/result-shape, and documentation findings remain closed.
  Concurrent-installer isolation remains the unanimously recorded follow-up and was not
  re-litigated.

## New findings

### [MAJOR] Intermediate symlinks in a raw target are not expanded, so `..` walks the wrong tree

`walkRawTarget` joins and records each ordinary component but never asks whether that component
is itself a symlink (`lib/installer.js:1393-1406`). `resolutionTouchpoints` calls `lstatSync`
only after the entire raw target has been reduced to `landed` (`lib/installer.js:1443-1451`).
That differs from kernel path resolution: an intermediate link is expanded immediately, and
subsequent `..` components walk from the expanded target rather than from the link's lexical
spelling.

I reproduced the resulting false success against the real installer on this POSIX host with
only Codex and Kimi detected:

- Kimi core was installed at `KM/skills/parley-deck`, with `subdir/transient` inside it.
- `mid` pointed to `KM/skills/parley-deck/subdir`.
- `B/skills` pointed to the raw target
  `../mid/transient/../../../../../away`.

The Codex path resolves to `away` only while Kimi's `subdir/transient` exists. The walker records
`mid`, but it applies the five `..` components to the spelling containing `mid`; its chains never
contain Kimi core's identity. Fleet install returned `ok:true`, reported Codex `installed` and
Kimi `replaced`, then left Codex's reported destination absent and its payload orphaned at
`away/parley-deck`.

The same gap reaches uninstall. After seeding both units separately, fleet uninstall returned
`ok:true` and reported both `removed`; quarantining Kimi broke the Codex cleanup spelling, and
one hidden `.parley-deck.<pid>.<timestamp>.removing` directory remained under `away`. Thus install
and uninstall do invoke the same gate, but they inherit the same false negative.

Fix the raw component walk so every encountered symlink is expanded into the pending component
stream, absolute targets reset to their root, and `..` is applied to the expanded path. Keep one
global hop budget across nested expansions. Add discriminatory install and uninstall regressions
using the exact `mid` target above; the existing simple raw-`..`, repeated-separator, root-clamp,
revisit, and cycle tests must remain.

### [MAJOR] Absolute Windows raw targets replay their drive or UNC root as ordinary components

For an absolute target, `walkRawTarget` correctly initializes `current` to the parsed root, but
then splits and iterates the complete unstripped raw string (`lib/installer.js:1394-1404`). With
the current function compiled unchanged except for injecting `path.win32`, I measured:

- `C:\\target\\x` records `C:\\C:`, `C:\\C:\\target`, and
  `C:\\C:\\target\\x`;
- `\\\\server\\share\\dir\\x` starts at `\\\\server\\share\\` and then records
  `\\\\server\\share\\server`, duplicating the server and share.

The final `lstatSync(landed)` therefore probes a non-path and stops the chain. Absolute link
targets are normal on Windows, especially for directory junctions, so the destination-chain
root fix does not restore dependency coverage for this common aliasing arm. A junction or
absolute symlink whose raw target passes through another planned destination can again admit
the same false-success install and uninstall states the gate exists to prevent.

Strip the raw root before iterating its components without normalizing away the remaining
`.`/`..` sequence. Make that raw-root operation injectable, as `splitAtRoot` is, and pin drive,
mixed-separator, root-relative, and UNC targets through the installer's walker. A Windows
junction install/uninstall regression should exercise the complete gate when a Windows runner
is available.

## Release judgement

Not releasable as 2.1.0. The one subsystem that must change is raw-link component resolution in
the destination-dependency gate: it must expand intermediate links and preserve platform roots,
with both install and uninstall regressions. The rest of the reviewed implementation is
release-ready.

## What I verified

- Read the live cooperation protocol, `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through
  cycle 23, `review/round-09/VOID.md`, the round-18/19 records, the complete cycle-23 diff
  `2b680a2..2b7ca3e`, and the complete implementation diff `49fc3ec..2b7ca3e`. The full diff
  passes `git diff --check`.
- Read the complete current identity/touchpoint gate and traced the same
  `aliasedDestinations(plan)` call into install before staging and uninstall before quarantine.
  The simple raw-target regression blocks both paths before mutation; the intermediate-link
  reproduction above proves both paths also share the remaining gap.
- Exercised `walkRawTarget` cases requested for this round: POSIX absolute and relative targets,
  repeated separators, `..` past root, a directory revisit, an intermediate symlink, and the
  exact 64-hop cycle bound. Exercised `splitAtRoot` with win32 relative paths, `C:` without a
  separator, drive-absolute and UNC paths, and mixed `/` and `\\` separators. Injected
  `path.win32` into the current walker in memory to measure the absolute-root replay above.
- Ran the full package suite from a `git archive` of `2b7ca3e`: **364/364 Node tests**, **54/54
  Python tests** on Python 3.14, and the 47-file manifest check with aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- Parsed `npm pack --dry-run --ignore-scripts --json` from an isolated archive/cache: **202
  files**, **48** under `skills/parley-bidding/`, with no nested `.gitignore`, `__pycache__`,
  `.pyc`, or `.pyo`.
- Compared the read-only source and integrated payload by path and SHA-256: 48 files on each
  side, source-only `.gitignore`, integrated-only `parley-addon.json`, and the same nine
  documented content differences. Parsed all **16** shipped JSON files and ran the adapter
  validator successfully on all **4** platform adapters.
- Confirmed the six-skill documentation, payload-integrity boundary, operational Python
  reporting, and the verbatim single-writer warning under `CHANGELOG.md` "Known limits".
- Removed every temporary archive and fabricated home. The reviewed repository remained clean
  at `2b7ca3e`; no file under `skills/parley-bidding/` was changed.
