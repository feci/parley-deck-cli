---
idea: integrate-parley-bidding-addon
review-round: 19
agent: codex-1
date: 2026-07-31
reviewed-commit: 2b680a2
---

## Verdict

BLOCK

## Outstanding findings — closed or not

- The cycle-21 nearest-anchor regression is closed. Existing inner parents no longer hide
  ordinary nesting, and the skewed-anchor regression passes at `2b680a2`.
- The buried-link regression is closed. `entryChain` retains the full parent ancestry; I also
  put a link below a link-parent inside the Kimi destination, and the plan was blocked before
  any write.
- The firmlink and canonical-normalization arms are closed on this host. The `/private` test
  passes, and composed/decomposed spellings of one missing runtime root were refused with zero
  writes on this normalization-insensitive APFS volume.
- **The raw-`..` symlink-target arm from my round-18 MAJOR is not closed.** Cycle 22 changed
  destination and entry chains, but `resolutionTouchpoints` still normalizes the target at
  `lib/installer.js:1403` before recording it. The measured false-success state survives in
  both install and uninstall; details are below.
- The dead `resolvedDestination`/`physicalKey` machinery and its stale containment comments
  are gone. The other round-18 NIT is still open: the `pathEntryExists` comment names deleted
  `uninstallSkillUnit` as a caller at `lib/installer.js:2277`.
- The earlier manifest, marker, Python, dry-run/result-shape, staging, commit/revert,
  quarantine/rollback, and cleanup findings remain closed. Concurrent-installer isolation is
  the recorded follow-up and was not re-litigated.

### [MAJOR] Raw symlink targets are still normalized before their dependencies are recorded

`resolutionTouchpoints` reads a link target and immediately computes
`path.join(path.dirname(hop), target)` (`lib/installer.js:1397-1407`). `path.join` collapses
`name/..` pairs. `record` and `entryChain` therefore see only the normalized endpoint, not every
directory entry the kernel must consult while resolving the raw target.

I reproduced the round-18 arm against `2b680a2` with a link at `B/skills` whose raw target was
`../KM/skills/parley-deck/transient/../../../../away`. The normalized endpoint is `away`, but
the link works only while `transient` exists inside the managed Kimi core.

- Fleet install returned top-level `ok:true`, reported Codex `installed` and Kimi `replaced`,
  then left the reported Codex destination absent and its payload orphaned at
  `away/parley-deck`.
- After seeding both managed units separately, fleet uninstall returned top-level `ok:true`,
  reported both `removed`, emitted no warning, and left
  `away/.parley-deck.<pid>.<timestamp>.removing` behind. Quarantining Kimi removes
  `transient`; the Codex cleanup path then becomes dangling, and forceful `rmSync` treats it as
  already gone.

This is the same single-process false-success property that B5 and the 2.1.0 changelog say is
impossible. The fix must resolve link targets component by component and retain each entry
actually consulted before applying `..`; passing only a normalized absolute path into
`entryChain` cannot recover the discarded dependency. Add discriminatory install and uninstall
regressions using the raw target above.

### [NIT] The destination-entry comment still names a deleted caller

`lib/installer.js:2276-2281` says every destination-entry check goes through
`pathEntryExists`, naming `uninstallSkillUnit` as one of the callers. That function was removed
in cycle 20; uninstall now checks entries inside `removeFleetAtomically`. This was outstanding
in round 18 and cycle 22 did not change it. Under the strict gate, the factual comment should be
corrected rather than carried into the release.

## New findings

### [MAJOR] The chain walkers discard Windows drive and UNC roots

Both new walkers split an absolute path and restart from `path.sep`:

- `identityChain`: `lib/installer.js:1332-1339`
- `resolutionTouchpoints`: `lib/installer.js:1372-1385`

That is valid for POSIX `/`, but not for Windows roots. Using Node's own `path.win32` semantics,
the exact loop transforms `C:\Users\a` into probes `\C:`, `\C:\Users`, and
`\C:\Users\a`, rather than `C:\`, `C:\Users`, and `C:\Users\a`. A UNC path
`\\server\share\dir\x` starts at `\server` and loses the `\\server\share\` root as well.
Those are not prefixes of the destination the installer will mutate.

Consequently `statSync`/`lstatSync` cannot obtain the destination's physical component
identities on Windows. The chain falls back to spelling-derived synthetic values, and the
touchpoint walk sees no junctions or symlinks. Two differently-spelled runtime roots that are
junctions into one skills container can therefore pass this gate; actual staging and commit use
the original valid paths, so the later unit can replace the earlier unit while both are reported
successful—the aliased-root failure this gate exists to prevent.

Windows is a supported release channel (`winget` and a Windows portable binary are documented
in `README.md:247`), while `.github/workflows/test.yml:18` runs the suite only on Ubuntu and the
Windows release job cross-builds rather than executes the installer. Start the walk at
`path.parse(resolved).root`, split only the remainder, and add Windows drive-root and UNC/junction
regressions. A pure injected path/filesystem test can cover the root construction, but an actual
Windows install test should cover junction identity and the zero-write result.

## Release judgement

Not releasable as 2.1.0. The one subsystem that must change is destination-dependency
resolution: it must preserve raw symlink traversal and walk from the real platform root on
Windows, with install and uninstall regressions for both. The stale caller comment must also be
removed before the strict gate can close.

## What I verified

- Read the live cooperation protocol, `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through
  cycle 22, `review/round-09/VOID.md`, the round-18 record, the cycle-22 diff
  `64e43f9..2b680a2`, and the complete implementation diff `49fc3ec..2b680a2`. Both diffs pass
  `git diff --check`.
- Read the complete current identity/touchpoint gate and traced it into
  `installFleetAtomically` before staging and `removeFleetAtomically` before quarantine. The
  uninstall path invokes the same gate, but therefore inherits both comparison gaps above.
- Reproduced the raw-`..` install orphan and uninstall quarantine residue against the real
  installer with no concurrent writer and no `--force` on install. All fabricated homes were
  below `/private/tmp` and were removed.
- Made one destination component inaccessible with mode `000`. Install failed preflight and
  wrote no other target. For uninstall I made the later Kimi ancestry inaccessible after both
  units were managed: Codex was quarantined first, Kimi failed with `EACCES`, Codex was restored,
  both destinations remained intact, and no `.removing` residue remained.
- Exercised `entryChain` with a link whose parent path was itself a link inside another managed
  destination; the plan was blocked and the link target stayed empty. This APFS host refuses a
  directory hardlink with `EPERM`; on a filesystem that permits one, equal `dev:ino` makes the
  gate conservatively block it. Likewise, genuinely different objects that report the same
  `dev:ino` can cause a false refusal, not a write-admitting false negative.
- Confirmed NFC-equivalent composed/decomposed names alias on this volume and are refused before
  writes. Existing components use `dev:ino`; missing tails use NFC names. On a
  normalization-sensitive filesystem this unconditional NFC choice can conservatively refuse
  two distinct missing paths, but it cannot admit a partial-fleet mutation.
- Used `path.win32` to evaluate the current root reconstruction for drive-letter and UNC paths,
  confirmed the malformed probe sequences above, and confirmed the repository has no Windows
  execution job despite shipping Windows artifacts.
- Ran the full package suite in an isolated temporary directory: **361/361 Node tests**, then
  **54/54 Python tests** on Python 3.14, and the 47-file manifest check with aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`. The suite's
  immutable-debris regressions initially resisted removal; I cleared their flags/permissions
  and removed the entire isolated directory.
- Parsed `npm pack --dry-run --ignore-scripts --json` with an isolated npm cache: **202 files**,
  **48** under `skills/parley-bidding/`, and no nested `.gitignore`, `__pycache__`, `.pyc`, or
  `.pyo`. The isolated cache was removed.
- Verified every review-created temporary path was gone and the reviewed repository remained
  at `2b680a2` with a clean working tree before this completed artifact was written. No file
  under `skills/parley-bidding/` was changed.
