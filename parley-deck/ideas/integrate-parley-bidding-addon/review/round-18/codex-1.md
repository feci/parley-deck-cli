---
idea: integrate-parley-bidding-addon
review-round: 18
agent: codex-1
date: 2026-07-31
reviewed-commit: 64e43f9
---

## Verdict

BLOCK

## Outstanding findings — closed or not

The two exact cycle-21 regressions pass, but the round-17 MAJORs are not closed as guarantees.

- The `/private` firmlink regression is closed only while both destinations choose the same
  nearest-existing-ancestor inode as their key anchor. If an existing directory lies between
  the shared ancestor and one destination, the anchors differ and the same firmlink spelling
  again bypasses containment. I measured the original `ok:true`, Codex `installed`, Kimi
  `replaced`, Codex destination absent end state at `64e43f9` with that one change.
- The two-link `B/skills -> kimiCore/redirect -> away` regression is closed only while the
  relevant link is directly under the planned destination. A link under an existing descendant
  directory gets a key anchored on that descendant, not on the outer destination, and the
  dependency disappears. A raw link target whose component walk enters the destination and
  later exits through `..` is also normalized by `path.join` before any of those intermediate
  components are recorded. Both forms reproduce the install false success and uninstall's
  silent quarantine residue at `64e43f9`.
- Hermes's round-17 NIT about the symlink-only comment is closed; cycle 21 replaced it.
- Kimi's round-17 NIT remains open: the `pathEntryExists` comment still names deleted
  `uninstallSkillUnit` as a caller (`lib/installer.js:2342`).
- The Python fail-closed, uninstall dry-run, dead `preflightUninstallUnit`, round-15
  containment/result-shape, and earlier manifest findings remain closed.
- Concurrent-installer isolation remains the unanimously recorded follow-up. I did not
  re-litigate it.

## New findings

### [MAJOR] Scalar nearest-ancestor keys still lose physical ancestry and path traversal

`physicalKey` stops at the nearest existing entry and returns one scalar
`dev:ino/uncreated-tail` key (`lib/installer.js:1323-1346`). `overlaps` then treats string-prefix
relations between two such keys as physical containment (`lib/installer.js:1449-1453`). That
relation is valid only when both keys happen to use the same inode anchor. An existing inner
parent makes the inner key start with that parent's inode while the outer key starts with the
outer destination's inode; neither string can be a prefix of the other even though the paths
are physically nested. `physicalEntryKey` inherits the same loss by keying only the immediate
parent (`lib/installer.js:1386-1390`).

Measured against the actual installer at `64e43f9`, without `--force` or another process:

1. I installed Kimi core, created an ordinary `existing/parent/skills` directory inside it,
   and pointed `CODEX_HOME` at `existing/parent`. Fleet install returned `ok:true`; Codex was
   `installed`, Kimi was `replaced`, and Codex's reported destination was absent afterwards.
   The surviving outer marker named Kimi. This is the round-15 signature with different
   anchors.
2. I repeated that setup with the inner home spelled through
   `/System/Volumes/Data/private/...` and the outer through `/private/...`. The two spellings
   had identical `stat` dev/ino, yet the different inner anchor again admitted the plan and
   erased Codex after reporting success. Cycle 21's firmlink test does not create an existing
   inner parent, so both of its keys share an anchor and it misses this arm.
3. I put `redirect -> away` below two existing directories inside Kimi core and pointed
   `CODEX_HOME` at it. Fleet install returned `ok:true`, left Codex's reported path dangling,
   and orphaned its files under `away`. After separately seeding both managed units, fleet
   uninstall returned both `removed`, emitted no warning, and left one hidden
   `.parley-deck.*.removing` directory. The shared install/uninstall gate is reached, but
   `physicalEntryKey(redirect)` is anchored on the deep parent rather than Kimi core.
4. A neutral link with raw target
   `../KM/skills/parley-deck/transient/../../../../away` behaved identically. Resolution
   depends on `transient`, which Kimi's replacement removes, but `path.join` at
   `lib/installer.js:1434` collapses the target to `away` before touchpoints are recorded.
   Install again reported success with a dangling destination and orphan; uninstall silently
   left quarantine residue.

The tail is not a filesystem-canonical identity either. On this normalization-insensitive APFS
volume I used canonically equivalent composed/decomposed spellings (`é` versus `e` plus a
combining acute) for two missing runtime roots. `toLowerCase()` at lines 1334 and 1389 preserved
the normalization difference, so fleet install returned `ok:true`, Codex `installed`, Kimi
`replaced`; `stat` showed both reported destinations were the same inode and the surviving
marker named Kimi. The existing case-only regression does not cover canonical Unicode
equivalence.

Replace the scalar-key/prefix model with a structured, component-wise physical representation
that retains every existing ancestor identity and the remaining tail, so two paths can be
compared through any shared inode rather than only one chosen anchor. The resolution walk must
record the raw components traversed by each link target before `..` normalization, and tail
names must use the supported filesystem's canonical case/Unicode equivalence. At minimum, add
regressions for the four measured arms above in both install and uninstall where applicable.
The current direct firmlink, direct-link, sibling, and cycle tests should remain as separate
guards. A symlink cycle itself terminates under the 64-step bound and then fails preflight with
`ELOOP`; that part is sound.

### [NIT] Cycle 21 leaves a dead realpath walk and factually stale comments

`resolvedDestination` is still called for every unit and stored as `entry.path`
(`lib/installer.js:1348-1372,1490`), but cycle 21 removed every read of `entry.path`. Its comment
still says it is used for containment, and the touchpoint preamble still contrasts itself with
that now-nonexistent final-path comparison. Remove the dead computation or restore it as an
explicit secondary check with accurate scope. Also close the outstanding `uninstallSkillUnit`
comment noted above. Under this idea's strict gate, these objective dead-code/documentation
issues are findings rather than deferred style preferences.

## Release judgement

Not releasable as 2.1.0. The one thing that must change is the destination-dependency model:
replace nearest-anchor scalar prefix comparison with component-correct physical ancestry and
raw resolution-chain tracking, including filesystem-canonical tail names. The cleanup NITs
should be closed as part of that change before the required fresh full-scope clean round.

## What I verified

- Read the live cooperation protocol, `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through
  cycle 21, `review/round-09/VOID.md`, all three round-17 reviews, the cycle-21 diff
  `fa3da41..64e43f9`, and the full implementation diff `49fc3ec..64e43f9`. Both requested
  diffs pass `git diff --check`.
- Read the complete current `lib/installer.js`, `lib/addon-manifest.js`, Python runner, relevant
  command-guard changes, changelog, and identity regressions. Traced `aliasedDestinations` into
  both `installFleetAtomically` and `removeFleetAtomically`; uninstall does invoke the same
  gate before quarantine, but inherits all comparison gaps above.
- Reproduced the plain different-anchor, firmlink different-anchor, deep-link, raw-`..` link
  target, and Unicode-normalization false-success arms against the real installer. Reproduced
  both uninstall residue arms with managed destinations. All fabricated homes were under
  temporary directories and were removed.
- Built a two-link cycle and confirmed the walk completed in 48 ms, preflight reported `ELOOP`,
  and the fleet performed zero writes.
- Ran the full package suite with an isolated temporary directory: **359/359 Node tests**, then
  **54/54 Python tests** on Python 3.14, and the 47-file manifest check with aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`. With
  `/usr/bin/python3` 3.9 first on `PATH`, the Python runner refused it as designed.
- Parsed the `npm pack --dry-run --ignore-scripts` inventory: **202 files**, **48** below
  `skills/parley-bidding/`, no nested `.gitignore`, `__pycache__`, `.pyc`, or `.pyo`.
- Validated all four platform adapters directly with
  `PYTHONDONTWRITEBYTECODE=1 python3 -B`, parsed all **16** bidding JSON files, and confirmed no
  payload symlinks or Python cache artifacts.
- Compared the read-only source and integrated payload without executing the source: 48 files
  versus 48; source-only `.gitignore`, integrated-only `parley-addon.json`, and the same nine
  documented content differences. No file under `skills/parley-bidding/` was changed.
- Removed every temporary suite/home directory I created. The reviewed repository remained at
  `64e43f9` with a clean working tree before this review artifact was written.
