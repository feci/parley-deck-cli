---
idea: integrate-parley-bidding-addon
review-round: 17
agent: hermes-1
date: 2026-07-31
reviewed-commit: fa3da41
---

## Verdict

BLOCK

## Outstanding findings — closed or not

Round 16 codex-1 MAJOR (symlink inside a managed dest, another runtime pointing through it, later commit destroys the link, earlier unit's dest goes dangling): **partially closed.** Cycle 20's `resolutionTouchpoints` catches the case where the symlink *entry* is inside a planned destination. It does not catch the case where the symlink *entry* is outside both destinations but an intermediate hop in the symlink's resolution chain passes through a planned destination. That remaining case is the MAJOR below — same defect class, same measured outcome, unguarded attack surface.

Round 16 codex-1 MINOR / kimi-1 MINOR (Python gate fails open on unparseable version and malformed runtime.python): **closed.** `resolveInterpreter` now regex-parses the version output and calls `fail` on no match; `declaredPythonFloor` returning null is now an explicit `fail` instead of a skipped comparison. Verified: python3 3.9.6 is rejected with a clear message; 3.14.6 passes 54/54.

Round 16 hermes-1 NIT (uninstall --dry-run missing per-action dryRun flag): **closed.** `uninstallCommand` now spreads `...(core.dryRun ? { dryRun: true } : {})`, matching install's shape. Verified: `dry.actions[0].dryRun === true`, and the dest is not touched.

Dead `preflightUninstallUnit`: **closed.** No remaining references in the tree.

## New findings

### [MAJOR] resolutionTouchpoints records only symlink entries, not their targets — a symlink chain whose intermediate hop is inside a planned destination is not caught

`resolutionTouchpoints` (lib/installer.js:1382-1409) walks the logical path component by component, calling `lstatSync` on each. When it finds a symlink, it records the symlink's *entry location* — `realpathSync(parent) + basename` — and continues to the next logical component. It never examines what the symlink points at. The `crosses` check in `aliasedDestinations` (lib/installer.js:1463-1465) tests whether any touchpoint overlaps the *other* unit's resolved destination. So the only symlinks that can trigger a refusal are those whose entry is lexically inside another planned destination.

That misses the case where the symlink entry is outside both destinations but the symlink's target — or an intermediate hop in a resolution chain — is inside one. Replacing that destination destroys the intermediate hop, breaking the other unit's resolution after it has already committed.

**Reproduction (run against fa3da41, actual installer, not a logic mock):**

1. Install kimi core at `$KIMI_CODE_HOME/skills/parley-deck`.
2. Inside that managed tree, create a `redirect` symlink pointing to an unrelated directory `$HOME/away`.
3. Set `CODEX_HOME` to a directory whose `skills` entry is a symlink pointing at `kimiCore/redirect`. The resolution chain is: `B/skills` → `kimiCore/redirect` → `$HOME/away`. The intermediate hop `kimiCore/redirect` is inside kimi's destination.
4. Run `install --target all --include-undetected --no-addons` with both `KIMI_CODE_HOME` and `CODEX_HOME` set.

Measured result: `ok: true`, codex `installed`, kimi `replaced`, codex's dest path `B/skills/parley-deck` no longer resolves (lstat fails — the chain is broken because kimi's commit replaced the tree containing `redirect`), codex's files orphaned at `$HOME/away/parley-deck` with a `target: codex` marker. The transaction reported success for a unit whose destination is now dangling.

**Why neither check fires:**
- `resolutionTouchpoints` for codex records `B/skills` (the symlink entry). `B/skills` is not inside kimi's resolved path. The intermediate hop `kimiCore/redirect` is never recorded because the function only walks the *logical* path (`B/skills/parley-deck`), not the *physical* resolution chain.
- `resolvedDestination` for codex gives `$HOME/away/parley-deck`; for kimi it gives `kimiCore`. No overlap, so containment is silent.

**Uninstall is affected identically.** With both targets managed and in the same uninstall plan, codex is quarantined first (renamed to `$HOME/away/.parley-deck.*.removing` through the symlink chain), then kimi is quarantined (kimiCore renamed aside, destroying `redirect`). Phase B's `rmSync(aside, {force:true})` for codex treats the now-broken path as successful cleanup — `removed`, no warning — while the quarantined tree remains at `$HOME/away/.parley-deck.*.removing`. Measured and confirmed.

**Why this is the same class as codex-1's round 16 MAJOR, not a new class:** both involve a managed tree containing a symlink, another runtime's resolution passing through that symlink, and a later commit/quarantine destroying it. The only difference is the *direction* of the indirection: codex-1's case had `CODEX_HOME` set to the symlink itself (entry inside the dest); this case has `CODEX_HOME` set to a directory whose `skills` subdirectory is a symlink pointing at something inside the dest (entry outside, target inside). The defect state is identical: `ok: true`, orphan tree, dangling dest.

**What the fix would need:** `resolutionTouchpoints` must also record the physical locations a symlink's resolution passes through, not just the entry. For each symlink found in the logical walk, `readlinkSync` the target, resolve it against the parent's realpath, and record that target path (with its own parent realpath-normalized) as an additional touchpoint. For chains, repeat until a non-symlink is reached. Then the `crosses` check would see `kimiCore/redirect` as a touchpoint of codex's resolution, recognize it overlaps kimi's resolved path, and reject the plan.

**Release judgement:** this must be fixed before 2.1.0. It is the same single-process false-success / partial-fleet state that B5 and the changelog claim cannot occur, reachable with a symlink configuration that is no more contrived than the one cycle 20 already guards against.

### [NIT] resolutionTouchpoints comment says "only symlinks matter" without stating the scope limitation

The comment at lib/installer.js:1383-1386 explains why plain components are excluded (shared-parent false positives) and asserts "a symlink is the only component that makes resolution depend on something a different unit might own." That is true for the symlink's *entry*, but the function does not examine the symlink's *target*, so the assertion over-covers. The comment should state that only symlink entries in the logical path are recorded, not intermediate targets in the resolution chain — the gap the MAJOR above describes.

## Release judgement

Not releasable as 2.1.0. The one thing that must change: `resolutionTouchpoints` must record the physical target (and chain intermediates) of each symlink it encounters, not just the entry, so that `aliasedDestinations` rejects a plan where one destination's resolution chain passes through another planned destination regardless of whether the symlink entry or an intermediate hop is inside it.

## What I verified

- `git diff a49d68f..fa3da41` (cycle 20) and `git diff 49fc3ec..fa3da41` (full since round 7): read every changed line in lib/installer.js, scripts/run-python-tests.js, and test/bidding-addon.test.js.
- 357 node tests pass (`node --test`), 0 fail.
- Python leg: 54/54 on python3 3.14.6; refuses 3.9.6 by design with a clear message.
- Manifest check: 47 files, sha256:7854adf1...68b95a6d, unchanged since 714712f.
- `preflightUninstallUnit`: no remaining references in the tree (dead code fully removed).
- Uninstall `--dry-run` per-action flag: verified `dry.actions[0].dryRun === true` and dest not touched.
- Python gate fail-closed: verified both the unparseable-version path and the null-floor path call `fail` (exit 1).
- **The MAJOR was reproduced against the actual installer** (require('../lib/installer')), not a logic mock. Two independent runs:
  - Install path: kimi installed first, then `--target all` with kimi+codex. Result: `ok: true`, codex `installed`, kimi `replaced`, codex dest dangling, orphan at `away/parley-deck`.
  - Uninstall path: both managed, `uninstall --target all`. Result: `ok: true`, both `removed`, no warning, orphan quarantined tree at `away/.parley-deck.*.removing`.
- Edge cases considered and ruled out as gaps:
  - Hardlinked directories: refused on macOS (EPERM); on Linux requires root. Even if possible, `physicalKey` (dev:ino) catches same-inode dests.
  - Firmlinks: transparent to `realpathSync`, not visible to `lstatSync` as symlinks, targets are fixed.
  - Mount points: `lstatSync` sees a directory, not a symlink; cannot be destroyed by rename (EBUSY).
  - Component becoming a symlink after planning: staging uses `mkdirSync` (creates dirs) and `copyRecursive` (throws on symlinks); dest untouched until commit.
  - `/var` → `/private/var` system symlink on macOS: produces a touchpoint at `/var` but it does not overlap any real destination; no false positive.
- The `aliasedDestinations` check is called by both `installFleetAtomically` (line 1511) and `removeFleetAtomically` (line 1626), so the touchpoints check does cover uninstall's quarantine path — but with the same gap as install.
- Commit/revert transaction integrity: the three-phase staging → commit → housekeeping flow with fleet-wide rollback on any commit failure is sound for independent destinations. The gap is not in the rollback mechanics but in the planning-time check that admits a dependent plan.
