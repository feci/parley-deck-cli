---
idea: integrate-parley-bidding-addon
review-round: 20
agent: hermes-1
date: 2026-07-31
reviewed-commit: 2b7ca3e
---

## Verdict

BLOCK

## Outstanding findings — closed or not

**Round 19 hermes-1 NIT-1 (stale comment referencing deleted `uninstallSkillUnit`): CLOSED.**
Cycle 23 replaced "and `uninstallSkillUnit`" with "and the removal transaction" at
line 2316. Zero occurrences of `uninstallSkillUnit` remain in the file. No behavioral
impact, as expected.

**Round 19 kimi-1 NIT (NFC fold on case-sensitive filesystems): CLOSED.**
`canonicalSegment` (line 1330-1332) now guards with `if (!CASE_INSENSITIVE_FS) return name`
before normalizing. On Linux (case-sensitive ext4/xfs), the name is returned unchanged;
on macOS/Windows, it is NFC-normalized and lowercased. This matches the filesystem's own
behavior. Correct.

**Round 19 codex-1 MAJOR (raw link target normalized before dependencies seen): CLOSED.**
`walkRawTarget` (line 1393-1407) now records each entry via `record(current)` before
applying `..`, so a target like `../KM/skills/parley-deck/transient/../../../../away`
records `transient` (the dependency inside another destination) before the `..` components
collapse the path. The regression at test line 1697 exercises this and passes. I confirmed
the recorded intermediate paths include every component the kernel must consult. The simple
raw-`..` case is fixed.

**Round 19 codex-1/kimi-1 MAJOR (chain walk starts at `path.sep` instead of platform root):
PARTIALLY CLOSED.**
`splitAtRoot` (line 1343-1348) now anchors on `impl.parse(resolved).root`, and both
`identityChain` and `resolutionTouchpoints` use it. The regression at test line 1729
exercises drive-absolute, UNC, and POSIX roots through the injected `path.win32` and
passes. However, `walkRawTarget` performs its own root extraction for absolute targets
(line 1394-1396) and still splits the raw target string — see MAJOR-2 below.

**Round 19 kimi-1 MINOR (firmlink respelling with existing inner parent — the pin arm):
CLOSED as a pin.** The regression at test line 1655 is present, correctly labelled as a
pin rather than evidence, and passes at both `2b680a2` and `2b7ca3e`.

**All earlier settled matters** (concurrent-installer isolation as recorded follow-up,
CHANGELOG "Known limits" wording verbatim at line 106-112, Python fail-closed gate,
uninstall dry-run parity, dead `preflightUninstallUnit`, manifest key confinement,
manifest entry rule, `insidePayload` check, fleet-wide transactional install/uninstall,
`pathEntryExists` using `lstatSync`, ownership predicate shared across health/install/
uninstall): verified present and unchanged. Not re-litigated.

## New findings

### [MAJOR-1] Intermediate symlinks in a raw target are not expanded, so `..` walks the wrong tree

`walkRawTarget` (lib/installer.js:1393-1407) joins and records each ordinary component
but never asks whether that component is itself a symlink. `resolutionTouchpoints`
calls `lstatSync` only after the entire raw target has been reduced to `landed`
(line 1443-1450). This differs from kernel path resolution: an intermediate link is
expanded immediately, and subsequent `..` components walk from the expanded target
rather than from the link's lexical spelling.

I reproduced the false success end-to-end through the real installer on this POSIX host:

- Installed kimi at `KM/skills/parley-deck`, created `subdir/transient` inside it.
- `mid` -> `KM/skills/parley-deck/subdir` (a symlink).
- `B/skills` -> raw target `../mid/transient/../../../../../away`.

Kernel resolution: `mid` expands to `KM/skills/parley-deck/subdir`, then `transient`
is `KM/skills/parley-deck/subdir/transient`, and five `..` components walk back to
`home`, so the final landing is `home/away`. The dependency on `kimiCore` (which holds
`subdir`) is real — remove kimi and the symlink chain breaks.

`walkRawTarget` does not expand `mid`. It records `home/mid` as a literal path, joins
`transient` to get `home/mid/transient`, then applies five `..` components:
`home/mid` -> `home` -> `/` (root, clamped) -> `/` -> `/`. The final landing is `/away`,
a non-existent path. The recorded touchpoints never include kimi core's identity.

Measured result: fleet install returned `ok:true`, codex `installed`, kimi `replaced`.
The payload was written to `home/away/parley-deck` (kernel resolution through the symlink),
but `B/skills/parley-deck` (the reported destination) was absent — exactly the orphaned-
payload, absent-destination state the gate exists to prevent.

The same false negative reaches uninstall. After seeding both units and running fleet
uninstall with `--force`, the gate did not block the plan; both units reported
`removed`/`missing`, and `away/parley-deck` survived as debris.

This is not the same as the cycle-23 raw-`..` fix. That fix ensured `..` is applied
after each entry is recorded; this is about what `..` is applied *to* when an
intermediate component is itself a symlink. The fix requires expanding intermediate
symlinks into the component stream as they are encountered, applying `..` to the
expanded path, with one global hop budget across nested expansions. Discriminating
install and uninstall regressions using the `mid` target above are needed; the existing
simple raw-`..`, repeated-separator, root-clamp, and cycle tests must remain.

### [MAJOR-2] Absolute Windows raw targets replay their drive or UNC root as ordinary components

`walkRawTarget` correctly initializes `current` to `path.parse(path.resolve(rawTarget)).root`
for absolute targets (line 1394-1395), but then splits and iterates the complete unstripped
raw target string (line 1397). On Windows, the root components survive the split and are
joined back onto the already-extracted root, doubling them.

I measured this by simulating `walkRawTarget` with `path.win32` injected:

- `C:\Users\a` produces records `C:\C:`, `C:\C:\Users`, `C:\C:\Users\a` — none are real paths.
- `\\server\share\dir` starts at `\\server\share\`, then records `\\server\share\server`,
  `\\server\share\server\share`, `\\server\share\server\share\dir` — server and share doubled.
- Forward-slash variants (`C:/Users/a`) are equally affected: the regex `/[\\/]+/` consumes
  the separators, but `C:` survives as a component.

On POSIX, the root `/` is pure separator, so it is consumed by the split and the bug does
not manifest. On Windows, the root includes named components (drive letter, server, share)
that survive the split.

The final `lstatSync(landed)` probes a non-existent path, the chain breaks, and `entryChain`
produces spelling-derived identities for the doubled paths. These never match a real
destination identity, so the dependency gate silently degrades — the exact failure mode
cycle 23 fixed in `splitAtRoot` for `identityChain` and the outer `resolutionTouchpoints`
walk. Absolute link targets are normal on Windows (directory junctions, `mklink`), so this
is a realistic aliasing arm on a shipped channel that CI never executes.

The fix: for absolute targets, resolve the target, extract the root, slice it off, and
split the rest — exactly as `splitAtRoot` does. The raw-root operation should be injectable
(as `splitAtRoot` is) so drive, UNC, mixed-separator, and root-relative targets can be
pinned from a POSIX host.

(Note: `walkRawTarget` is not exported and not injectable. `splitAtRoot` was made
injectable and exported for exactly this kind of Windows-from-POSIX testing. The same
treatment is needed here.)

## Release judgement

Not releasable as 2.1.0. The one subsystem that must change is raw-link component
resolution in `walkRawTarget`: it must expand intermediate symlinks (MAJOR-1) and
preserve platform roots for absolute targets (MAJOR-2). Both produce false negatives
in the destination-dependency gate — the gate returns `ok:true` for plans where one
destination's resolution passes through another, which is the exact condition the gate
exists to refuse. Both reach install and uninstall identically (both call
`aliasedDestinations`). The rest of the reviewed implementation is release-ready.

## What I verified

**Suite:** 364 node tests, 0 fail (`node --test`). Python leg: 54/54 on Python 3.14
(via `PATH=/opt/homebrew/bin:$PATH node scripts/run-python-tests.js`); correctly refuses
3.9.6 by design (add-on declares `>=3.10`). Manifest check: 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged
since `714712f`. `verifyPayload` returns `ok:true`.

**Cycle-23 regressions (all pass at 2b7ca3e):**
- Raw link target's intermediate directories (line 1697): the `transient` dependency
  is recorded before `..` collapses. PASS.
- Chain walk starts at platform root (line 1729): `splitAtRoot` with injected `path.win32`
  produces correct drive, UNC, and POSIX roots. PASS.
- Firmlink respelling with existing inner parent (line 1655, pin): passes at both
  `2b680a2` and `2b7ca3e`. PASS (as pin).

**Edge cases probed directly (walkRawTarget, replicated from source):**
- Absolute POSIX target (`/opt/data`): records `/opt`, `/opt/data`. Correct — root `/`
  is consumed by the split.
- Repeated separators (`target//a///b`): regex `/[\\/]+/` collapses them. Correct.
- `..` past the root (`/link` + `../../../../away`): `path.dirname("/") === "/"` on POSIX,
  clamps at root, records `/away`. Correct.
- Chain that revisits a directory (`sub/../sub/../sub`): records each visit. Correct.
- 64-hop bound: circular chain (A->B->A with absolute targets) performed exactly 64
  iterations and returned. No infinite loop. Defensive bound is adequate.
- Absolute Windows target (`C:\Users\a` with `path.win32`): records `C:\C:`,
  `C:\C:\Users`, `C:\C:\Users\a` — **WRONG**, root doubled. MAJOR-2.
- UNC Windows target (`\\server\share\dir` with `path.win32`): records
  `\\server\share\server`, etc. — **WRONG**, server/share doubled. MAJOR-2.
- Intermediate symlink in raw target: `mid` not expanded, `..` walks wrong tree.
  MAJOR-1.

**splitAtRoot edge cases (injected `path.win32`):**
- Relative input: `win.resolve` anchors to cwd, root is `C:\` (on this host's cwd).
  Expected — `splitAtRoot` is defined on absolute resolved paths.
- `C:` without separator: `win.isAbsolute("C:")` is false (drive-relative). Treated as
  relative, resolves to cwd. Correct for drive-relative paths.
- Mixed `/` and `\`: `win.resolve` normalizes separators. Correct.
- Drive-absolute: root `C:\`, parts stripped. Correct.
- UNC: root `\\server\share\`, parts stripped. Correct.
- POSIX unchanged: root `/`, parts as expected. Correct.

**Uninstall quarantine path parity:** `installFleetAtomically` (line 1528) and
`removeFleetAtomically` (line 1643) both call `aliasedDestinations(plan)` with identical
plan structure. The alias gate is the same function. Post-gate checks differ
appropriately (install: source validation via `preflightSkillUnit`; uninstall: marker
problem, existence, ownership). Both paths inherit the same false negatives from
MAJOR-1 and MAJOR-2 — confirmed by reproducing the intermediate-symlink false success
on both install and uninstall.

**Code reviewed:** `lib/installer.js` (2577 lines — full read of `splitAtRoot`,
`identityChain`, `entryChain`, `walkRawTarget`, `resolutionTouchpoints`,
`aliasedDestinations`, `installFleetAtomically`, `removeFleetAtomically`,
`preflightSkillUnit`, `destinationAncestorObstacle`, `pathEntryExists`,
`installerOwnsDestination`, `quarantineName`, `canonicalSegment`,
`copyPayloadAtomically`, `commitStagedUnit`, `revertStagedUnit`, `discardBackup`).
`test/bidding-addon.test.js` (new tests at lines 1655-1753). `scripts/run-python-tests.js`.
Diffs: `2b680a2..2b7ca3e` (cycle 23 only) and `49fc3ec..2b7ca3e` (full since round 7).

**Disk hygiene:** all scratch directories under `os.tmpdir()` were removed with
`fs.rmSync(home, {recursive: true, force: true})` in finally blocks. No repo mutations.
Working tree clean at `2b7ca3e` (`git status` clean). No temporary checkouts or archives
needed.
