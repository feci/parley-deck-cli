---
idea: integrate-parley-bidding-addon
review-round: 14
agent: hermes-1
date: 2026-07-31
reviewed-commit: d7ab1c3
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All findings from round 13 are addressed at d7ab1c3. I verified each fix by reading the code and running the tests.

- **Symlinked manifest (codex-1 MAJOR, hermes-1 MINOR): CLOSED.** `manifestEntryProblem` (addon-manifest.js:206) checks `lstatSync` for symbolic link and regular file, and is called by `verifyPayload` (line 220), `manifestFileHash` (line 112), and `hasManifest` (line 102 uses `lstatSync` directly). The test "a symlinked manifest is a payload defect, not payload authority" reproduces the exact attack (rename manifest out, symlink back) and confirms `verifyPayload` returns `ok:false`, `doctor` reports `malformed`, and `managed` is not `true`. I confirmed the predicate is applied in every manifest-reading path that gates payload integrity (see New findings for one narrow exception).

- **Aliased physical destinations (codex-1 MAJOR, hermes-1 MINOR): CLOSED.** `physicalKey` (installer.js:1415) resolves `path.dirname(dest)` with `realpathSync` and appends the basename. `aliasedDestinations` (line 1426) builds a Map keyed by physical key and blocks every unit sharing one. The test "two targets resolving to one physical directory refuse the whole plan" confirms `ok:false` with zero writes. Applied in both `installFleetAtomically` (line 1465) and `removeFleetAtomically` (line 1580). I reproduced the realpath fallback gap (see New findings).

- **Dry-run parity (codex-1 MINOR, hermes-1, kimi-1): CLOSED.** The dry-run and real install paths are unified through `installFleetAtomically`, which runs `aliasedDestinations` + `preflightSkillUnit` for both. The per-action `dryRun` flag is restored (line 674, 1498). Uninstall dry-run is unified through `removeFleetAtomically`. The test "install --dry-run answers exactly what the real install would" verifies `dry.ok === real.ok` for both install and uninstall with a damaged marker. I confirmed the only remaining difference is structural redundancy, not a predicate gap (see New findings).

- **Repairable damaged selection (kimi-1 MINOR): CLOSED.** `preflightSkillUnit` (line 1217) no longer blocks on `unit.markerProblem` for install. The comment explains: install's units come from discovery and flags, never from the marker. `writeMarker` overwrites the damaged marker with the correct selection. Health still reports the damage (`skillUnitStatus` line 1899), and uninstall still refuses (`removeFleetAtomically` line 1588, `preflightUninstallUnit` line 770). The test "a damaged recorded selection is repairable, and the message says how" confirms the full cycle: `doctor` malformed → uninstall blocked → install ok → `doctor` valid. The message now says "Re-run install to rewrite it" (lines 1046, 1054).

- **Core skill as add-on (kimi-1 NIT): CLOSED.** `authorize` (line 1131) refuses `name === SKILL_NAME` when `fromMarker` is true. The test "a recorded selection naming the core skill is refused" confirms `doctor` reports malformed with the matching message.

## Ruling on concurrent-installer isolation

**Recorded follow-up. Does not gate 2.1.0.**

The vulnerability is real: two interleaved installer processes targeting the same skills directory can overwrite each other's committed units, and one process's rollback can destroy the other's already-committed core — while both return `ok: true`. I traced the mechanism: `commitStagedUnit` renames `temp → dest` with no cross-process coordination, and `revertStagedUnit` renames `dest → temp` then `backup → dest`, clobbering whatever a concurrent process placed at `dest` between those two renames. The temp/backup/quarantine names already include `process.pid` + `Date.now()`, so temp-file collisions are prevented — but the commit and rollback operations on the shared destination path are not.

My reasoning for not gating:

1. **No normal workflow produces this.** `npx -y parley-deck-skill install` is sequential. CI pipelines run one install at a time. Two users sharing a home directory and running install simultaneously is an operator error, not a supported mode.

2. **A lock protocol is a design change, not a fix.** It introduces stale locks (crash recovery), NFS semantics (`flock` is advisory and non-functional across some network filesystems), and a new failure mode (lock acquisition timeout) that must itself be handled gracefully. Implementing this at round 14 of a fix-up sequence — as the brief's author argues — is the wrong moment for a new mechanism.

3. **The scope was already settled by two reviewers.** hermes-1 and kimi-1 both scoped it out in round 13. codex-1's finding is accurate but the remedy is disproportionate to the risk for this release.

4. **The module disclaims tamper resistance.** `addon-manifest.js:11` states "this is defect detection, not tamper resistance." Anyone who can write the skills directory can interfere with an install. Concurrent installers are the same threat model, just from a second copy of the same tool.

**What the release notes should tell users:**

> Do not run two `parley-deck-skill install` (or `uninstall`) commands concurrently against the same skills directory. The installer does not hold a cross-process lock; concurrent runs can overwrite each other's committed units and both report success. Run installations sequentially.

## New findings

Three findings, all below the blocking threshold. None gates 2.1.0.

### NIT-1: `readManifest` in `runtimeAvailability` is not guarded by `manifestEntryProblem`

`runtimeAvailability` (installer.js:2026) calls `addonManifest.readManifest(root)` directly on an installed destination. `readManifest` uses `fs.readFileSync` (addon-manifest.js:127), which follows symlinks. In the normal path this is safe because `skillUnitStatus` only calls `runtimeAvailability` when `ok` is true, which requires `manifestProblems` to have returned `[]`, which requires `verifyPayload` to have passed — and `verifyPayload` calls `manifestEntryProblem` first.

The narrow gap: the **legacy marker exemption** (manifestProblems, line 2153). When a marker has no `markerSchema` and the source ships no manifest (`sourceHasManifest === false`), `manifestProblems` returns `[]` without calling `verifyPayload`. If the destination then carries a symlinked manifest, `readManifest` in `runtimeAvailability` follows the link and reads the `runtime` field from an external file.

Impact: the runtime availability check (python version probe) reads from an external file. Payload integrity is not affected — the legacy exemption skips integrity validation entirely, so this is not a false green on the mechanism that justifies shipping. No add-on currently shipped lacks a manifest (`parley-bidding` ships one, so `sourceHasManifest === true` for it), so this path is unreachable for the payload in this release. It is within the module's stated "defect detection, not tamper resistance" scope.

Fix if desired: have `runtimeAvailability` call `manifestEntryProblem` before `readManifest`, or have `readManifest` itself call it as a guard.

### MINOR-1: `physicalKey` realpath fallback does not resolve symlinks in existing ancestors

I reproduced this with a standalone Node script. When `path.dirname(dest)` does not exist, `realpathSync` throws `ENOENT` and the fallback uses `path.resolve`, which normalizes `.` and `..` but does **not** resolve symlinks in existing ancestors.

Scenario: `~/.codex` is a symlink to `/shared`, `~/.codex/skills` does not exist yet, and another target's skills container is `/shared/skills` (also not yet created). `physicalKey` for the codex target returns `~/.codex/skills/parley-deck` (literal, symlink unresolved); `physicalKey` for the other target returns `/shared/skills/parley-deck`. Different keys, same physical destination — alias not detected.

This is the same class of bug cycle 17 fixed, with a more exotic trigger: the symlink must be at an ancestor of the skills directory, not at the skills directory itself, and the skills subdirectory must not yet exist. The tested scenario (containers themselves symlinked to the same existing directory) is caught correctly because `realpathSync` succeeds on existing paths.

Fix: in the catch block, walk up to the first existing ancestor, `realpathSync` that, then append the remaining non-existent path components.

### NIT-2: Redundant first preflight in `installCommand` can mask alias message

`installCommand` (lines 622–654) retains a preflight block guarded by `if (!context.options.dryRun)` that calls `preflightSkillUnit` for every unit. This block does **not** check `aliasedDestinations`. `installFleetAtomically` (line 1465) checks both `aliasedDestinations` and `preflightSkillUnit` for all runs.

For a real install with both an alias and another blocker: the first preflight may catch the other blocker and return early, never reporting the alias. For a dry run: `installFleetAtomically` checks aliases first (`aliased.get(unit) || preflightSkillUnit`), so the alias is reported. The `ok` values match (both `false`), but the reported `message` can differ. The test checks `dry.ok === real.ok`, not message equality.

This is not a functional gap — the install is correctly refused in both cases. The first preflight block is now redundant (installFleetAtomically does the same checks plus aliases) and could be removed for clarity, but its presence does not produce a wrong answer.

## Release judgement

**Releasable as 2.1.0.** Nothing I found must change before release.

The three new findings are all below the blocking threshold: two NITs that are within the module's stated scope limitations or produce no wrong answer, and one MINOR that is the same class of bug already fixed, with a more exotic trigger. The concurrent-installer isolation question is a recorded follow-up with a clear release-note warning.

## What I verified

1. **Test suite.** `npm run test`: 349/349 node tests pass, 0 fail. Python leg: 54/54 on `/opt/homebrew/bin/python3` 3.14.6; refuses `/usr/bin/python3` 3.9.6 by design (floor `>=3.10`). Manifest check: `parley-bidding: ok (47 files, sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d)`, matching `714712f` exactly.

2. **Cycle 17 diff.** Read the full `git diff dd8d756..d7ab1c3` — three files: `lib/addon-manifest.js` (+30), `lib/installer.js` (+149/-64), `test/bidding-addon.test.js` (+110). Also read `git diff 49fc3ec..d7ab1c3` for the full since-round-7 scope and the CHANGELOG changes.

3. **Manifest predicate coverage.** Traced every manifest read path: `hasManifest` (lstatSync, line 102), `manifestEntryProblem` (lstatSync, lines 206–217), `manifestFileHash` (guarded by manifestEntryProblem, line 112), `verifyPayload` (guarded by manifestEntryProblem, line 220), `readManifest` (readFileSync, line 127 — unguarded, called directly at installer.js:1826 and 2026). Confirmed installer.js:1826 (`writeMarker`) reads source/staged temp (trusted, copyRecursive refuses symlinks). installer.js:2026 (`runtimeAvailability`) reads destination — guarded by verifyPayload in the normal path, unguarded only in the legacy exemption path (NIT-1).

4. **Aliased destinations.** Read `physicalKey` (line 1415) and `aliasedDestinations` (line 1426). Confirmed both `installFleetAtomically` (line 1465) and `removeFleetAtomically` (line 1580) apply the check. Reproduced the realpath fallback gap with a standalone Node script creating a symlinked ancestor with non-existent skills subdir (MINOR-1). Confirmed the tested case (containers themselves symlinked to same existing dir) is caught.

5. **Dry-run parity.** Traced both install and uninstall paths. Confirmed dry-run and real install both go through `installFleetAtomically` (line 662), which runs `aliasedDestinations` + `preflightSkillUnit` for both. Confirmed uninstall dry-run and real both go through `removeFleetAtomically` (line 742). Confirmed per-action `dryRun` flag is present (line 674, 1498). Identified the redundant first preflight block (lines 622–654) that runs only for real installs and does not check aliases (NIT-2).

6. **Self-healing of damaged selection.** Traced `targetSkillUnits` (line 1095): for install, `fromMarker` is false (line 1123–1125), `expectedAddonNames` ignores `markerAddonNames` (line 1078), so the damaged marker never steers install's path scope. `preflightSkillUnit` (line 1217) no longer blocks on `markerProblem`. `writeMarker` (line 1799) overwrites the marker with the correct selection from `selectedAddons`. Confirmed health still reports the damage (`skillUnitStatus` line 1899) and uninstall still refuses (`removeFleetAtomically` line 1588, `preflightUninstallUnit` line 770).

7. **Core-as-add-on refusal.** Confirmed `authorize` (line 1131) checks `name === SKILL_NAME` only when `fromMarker` is true (non-install commands with no explicit flags). Confirmed install never produces a core-as-addon unit (`discoverAddons` skips `CORE_SKILL_NAME`, line 968; `fromMarker` is false for install).

8. **Concurrent installer isolation.** Traced `commitStagedUnit` (line 1760), `revertStagedUnit` (line 1778), and `quarantineName` (line 803). Confirmed temp/backup/quarantine names include `process.pid` + `Date.now()` (lines 806, 1701, 1702), preventing temp-file collisions but not commit/rollback races on the shared destination. Confirmed no lock/flock/mutex mechanism exists (searched for `lock`, `flock`, `mutex` — zero matches).

9. **CHANGELOG.** Read the full CHANGELOG diff. The release notes describe fleet-wide atomicity, transaction semantics, `--force` scope, marker validation, and the legacy exemption. They do not mention concurrent installers — the release-note warning I recommend above should be added.
