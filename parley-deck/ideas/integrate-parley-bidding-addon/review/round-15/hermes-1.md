---
idea: integrate-parley-bidding-addon
review-round: 15
agent: hermes-1
date: 2026-07-31
reviewed-commit: 26478e9
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All findings from round 14 are addressed at 26478e9. I verified each by reading the code, running the tests, and reproducing the edge cases.

- **Symlinked manifest — readManifest unguarded (hermes-1 NIT-1, codex-1 MINOR): CLOSED.** `readManifest` (addon-manifest.js:124-132) now calls `manifestEntryProblem` before `readFileSync`. I traced every manifest read path in both files: `hasManifest` (lstatSync, line 100), `manifestFileHash` (guarded, line 112), `readManifest` (guarded, line 129), `verifyPayload` (guarded, line 228). All four share the predicate. The test "every read of the manifest shares the regular-file rule" (test:1603) confirms all four refuse a symlinked manifest. No fifth reader or bypass exists — `rewriteStagedGeminiManifest` (installer.js:898) reads `gemini-extension.json`, not `parley-addon.json`.

- **physicalKey realpath fallback (hermes-1 MINOR-1, codex-1 MAJOR): CLOSED.** `physicalKey` (installer.js:1385) now walks to the nearest existing ancestor, keys it by `stat.dev`/`stat.ino` (statSync follows symlinks), and appends the not-yet-created tail with case-normalization on darwin/win32. I reproduced the exact scenario from my round-14 MINOR-1: `~/.codex` is a symlink to `/shared`, `skills/` does not exist yet, another target points at `/shared/skills` directly. Both keys are `dev:ino/skills/parley-deck` — identical. I also tested symlinked tails (statSync follows them, same dev/ino, empty tail) and case-only spellings on this case-insensitive volume (tail lowercased, same key). The tests "aliased runtime roots are caught before anything exists to realpath" (test:1604) and "case-only spellings of one directory are one destination" (test:1623) cover the two measured false negatives.

- **Redundant install preflight masking alias message (hermes-1 NIT-2): CLOSED.** Cycle 18 deleted the separate preflight block from `installCommand` (the 44-line removal at line 615). Install is now fully unified: both dry and real go through `installFleetAtomically`, which checks `aliased.get(unit) || preflightSkillUnit(...)` inside the function (line 1453). I confirmed message parity by running install with an alias + a separate blocker: both dry and real report the same `action` and `message` for every unit.

- **Uninstall dry-run promising removals the real command refuses (codex-1 MINOR): CLOSED.** The fleet gate in `removeFleetAtomically` (line 1601) now runs before any dry-run removal is recorded. The test "uninstall --dry-run does not promise removals the real command refuses" (test:1658) confirms `dry.ok === real.ok` and action-count parity.

- **CHANGELOG Known limits (round-14 ruling): CLOSED.** The "Known limits" heading (CHANGELOG.md:109-115) carries codex-1's release-note wording verbatim. The single-writer limit is documented; the concurrency follow-up does not gate 2.1.0.

## New findings

### MINOR-1: Uninstall dry-run and real diverge on `action` and `message` when an alias coexists with a separate blocker

**Reproduced.** `uninstallCommand` (installer.js:664-696) retains a separate preflight block guarded by `if (!context.options.dryRun)` that calls `preflightUninstallUnit` for every unit. This block does NOT check `aliasedDestinations`. It returns early — before `removeFleetAtomically` — when it finds any blocker (ownership, markerProblem). `removeFleetAtomically` (line 1565) checks aliases first, then the same preflight conditions, for both dry and real.

When a non-aliased blocker (e.g. a corrupted marker on target qwen) coexists with aliased destinations (e.g. hermes symlinked to codex):

- **Dry-run** skips the separate preflight, enters `removeFleetAtomically`, and the aliased units get `action: "blocked"`, `message: "Destination is shared by codex/parley-deck, hermes/parley-deck — they resolve to the same directory (...)"`.
- **Real** runs the separate preflight, finds the qwen blocker, and returns early with every unit as `action: "skipped"`, `message: "Not attempted: another skill or target in this uninstall failed preflight."`. The aliased units never reach `removeFleetAtomically`.

`ok` matches (both `false`). `action` and `message` do not. I reproduced this with and without `--force`, with both an ownership blocker and a markerProblem blocker. The existing dry-run parity test (test:1658) uses a single target and checks `ok` and action-count, not messages, so it does not catch this.

This is the exact same class of bug cycle 18 just fixed for install by removing the redundant separate preflight. The uninstall path's mirror-image block (lines 664-696) was not removed. The fix is the same: delete the separate preflight block and let `removeFleetAtomically` handle all preflight checks for both dry and real, exactly as `installFleetAtomically` does for install.

**Does not gate 2.1.0.** Both paths correctly refuse the operation (`ok: false`, zero deletions). The divergence is in the reported reason, not the outcome. A user running `uninstall --dry-run` to preview gets the more informative alias message; the real command gets the less informative generic skip. No data is lost or corrupted.

No other new findings.

## Release judgement

**Releasable as 2.1.0.** Nothing I found must change before release.

MINOR-1 is a cosmetic message divergence in a multi-target edge case (alias + separate blocker in the same uninstall). Both dry-run and real correctly refuse; `ok` matches. The fix is a one-block deletion parallel to the install fix already applied this cycle, but it can ship as a follow-up without risk.

The single-process transaction guarantee holds: install staging/commit/revert cleanup is sound, uninstall quarantine/rollback is sound, and the concurrency follow-up is documented under Known limits.

## What I verified

1. **Test suite.** `npm run test`: 353/353 node tests pass, 0 fail. Python leg refuses 3.9.6 by design (floor >=3.10). Manifest check: 47 files, aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, matching `714712f` exactly.

2. **Cycle 18 diff.** Read `git diff d7ab1c3..26478e9` — four files: CHANGELOG.md (+8), lib/addon-manifest.js (+8), lib/installer.js (+89/-49), test/bidding-addon.test.js (+72). The install preflight deletion, physicalKey rewrite, readManifest guard, and removeFleetAtomically dry-run gate are all present as described.

3. **physicalKey completeness.** Tested five edge cases with standalone Node scripts:
   - Symlinked ancestor with non-existent tail (my round-14 MINOR-1 scenario): same key. ✓
   - Case-only spellings on case-insensitive volume: same key. ✓
   - Symlinked tail (destination itself is a symlink): statSync follows it, same dev/ino, empty tail. ✓
   - `..` normalization: `path.resolve` normalizes before the walk. ✓
   - Two different directories: different inodes, different keys. ✓
   Bind mounts are not testable on macOS; on Linux, `statSync` returns the same `dev`/`ino` for bind-mounted directories, so the key would match. Windows path forms are covered by `CASE_INSENSITIVE_FS` and `path.resolve` backslash handling. The TOCTOU case (tail becomes existing between planning and staging) is unreachable in a single process — `aliasedDestinations` is called at the start of `installFleetAtomically`/`removeFleetAtomically`, before any write — and the concurrency follow-up covers multi-process.

4. **Manifest predicate coverage.** Traced every manifest read in both files. All four exported readers (`hasManifest`, `manifestFileHash`, `readManifest`, `verifyPayload`) now call `manifestEntryProblem` or use `lstatSync` directly. `readManifest`'s new guard (line 129) closes the gap I flagged as NIT-1 in round 14. No direct `readFileSync` of `parley-addon.json` exists outside `readManifest` itself. `rewriteStagedGeminiManifest` reads `gemini-extension.json`, not the addon manifest.

5. **Dry-run message parity.** Ran install and uninstall with alias + separate blocker scenarios:
   - Install: dry and real both go through `installFleetAtomically`. Messages and actions match for all units. ✓
   - Uninstall: dry goes through `removeFleetAtomically` (alias message); real hits the separate preflight (generic skip message). Messages and actions diverge for aliased units. ✗ (MINOR-1)

6. **Single-process transaction guarantee.** Traced install staging (1489-1505), commit (1507-1540), revert (1513-1525, 1769-1774), and cleanup (1542-1552). A staging failure cleans all staged temps. A commit failure reverts all committed units in reverse and cleans all temps. A revert failure reports the backup location and continues reverting the rest. Traced uninstall quarantine (1622-1667) and rollback (1631-1642): a quarantine failure rolls back all quarantined units; a rollback failure reports the aside location; a delete failure is a warning (destination already gone). The fleet is untouched in every failure path. Sound.

7. **CHANGELOG.** The "Known limits" heading is present with the single-writer warning. The release notes describe fleet-wide atomicity, transaction semantics, `--force` scope, marker validation, and the legacy exemption.
