---
idea: integrate-parley-bidding-addon
review-round: 11
agent: codex-1
date: 2026-07-31
reviewed-commit: 12f9071
---

## Verdict

BLOCK

## Outstanding findings — closed or not

- The **specific seven round-10 regressions are closed**. In the complete 332-test run,
  cycle 14 blocks mode-frozen destination trees before install or uninstall, blocks a
  foreign marker in the last uninstall target before any deletion, preserves a committed
  replacement when backup cleanup fails, and rejects a manifest-bearing marker whose
  `markerSchema` alone was deleted. The two corrected Python-leg claims now distinguish the
  3.14 run from the deliberate 3.9 refusal.
- `codex-1`'s replacement-cleanup finding is closed at the transaction boundary:
  `copyPayloadAtomically` no longer reports a committed replacement as failed merely because
  post-commit backup cleanup failed. The mode-permission preflight arm is also closed.
  However, the broader B5 removability claim is **not closed**: the new probe does not match
  `rmSync` for other filesystem states, as the immutable-file reproduction below shows.
- `hermes-1`'s fleet-wide-uninstall finding is closed for its foreign-marker scenario.
  `uninstallCommand` now preflights every target and unit before the first deletion. That
  fleet gate can still approve an incomplete removability model and then delete a partial
  fleet, so the property itself remains open through a new door.
- `codex-1`'s one-field manifest downgrade is closed exactly as filed. A current marker that
  keeps `manifest` but loses `markerSchema` is malformed. The intended no-silent-downgrade
  property is **not closed**, because deleting both fields still takes the legacy exemption
  even for `parley-bidding`, which did not exist in the 2.0.0 release.
- The two recorded follow-ups remain follow-ups, not findings in this round: only
  `parley-bidding` ships a manifest, and the `dirExists` guard does not discover an
  unselected add-on through a dangling symlink.

## New findings

### [CRITICAL] A core marker can steer forced uninstall outside the skills directory

`markerAddonNames` returns the core marker's `addons` array without validating its values
(`lib/installer.js:947-955`). `expectedAddonNames` treats those values as the uninstall
selection, and `targetSkillUnits` passes each one directly to
`path.join(skillsDir, name)` (`lib/installer.js:974-1015`). With `--force`,
`preflightUninstallUnit` deliberately skips ownership, so traversal components supplied by a
corrupt marker become deletion targets outside the runtime's skills directory. The new
fleet-wide preflight makes this worse in one precise sense: it validates that the
marker-derived out-of-scope tree can be removed, then authorizes the plan.

I reproduced this at exact `12f9071` in an isolated home:

1. Installed the normal Codex target.
2. Changed only the core marker's selection to
   `addons: ["../../outside-sentinel"]`.
3. Created `$HOME/outside-sentinel/KEEP`.
4. Ran forced uninstallation with no explicit add-on selector.

The command returned `ok:true`. It reported the core removed, reported a skill literally
named `../../outside-sentinel` removed, resolved its destination to
`$HOME/outside-sentinel`, and deleted the sentinel tree.

`--force` may override ownership at a caller-selected skill destination; it must not let
mutable marker data expand the command's path scope. Before release, recorded selections must
be schema-validated before path construction and every derived destination must be confined
to an exact direct child of `skillsDir`. Invalid, unknown, non-string, duplicate, `"."`,
`".."`, or separator-containing entries must make health and mutation fail closed. A
regression must preserve an out-of-root sentinel under forced uninstall.

### [MAJOR] The removability preflight still approves a predictable 83-unit partial uninstall

`firstRemovalObstacle` checks permissions on directories but deliberately does nothing for
non-directories (`lib/installer.js:1950-1977`). On macOS, a user-immutable (`uchg`) regular
file can retain ordinary mode bits while recursive removal fails. The full-tree access walk
therefore is not a mirror of `rmSync`.

I installed all 84 units across fourteen targets, set `uchg` on
`~/.aionrs/skills/parley-worktrees/SKILL.md` in the final unit, and ran an ordinary unforced
fleet uninstall as uid 501. Fleet preflight passed. The command removed **83 units**, then
the last unit failed with `ENOTEMPTY` and left the flagged file. This is a predictable
destination state and the same B5 partial-fleet result cycle 14 was meant to exclude.

The helper is also not exact in the opposite direction. Adding an empty mode-0555 directory
to an owned destination makes preflight block it as “cannot be emptied (EACCES),” while a
direct `fs.rmSync(dest, {recursive:true, force:true})` succeeds because removing an empty
directory needs permission on its parent, not write permission inside the empty directory.

The fix should not continue to equate an access walk with deletion semantics. Uninstall needs
a transaction/rollback or quarantine design that cannot delete earlier units when a later
cleanup fails, with regressions for an immutable file in the final unit and a removable empty
read-only directory. The install path's post-commit warning already demonstrates the relevant
distinction between logical commit and cleanup.

### [MAJOR] The “genuine 2.0.0” bidding-marker regression manufactures a state that never shipped

When both `markerSchema` and `manifest` are absent, `manifestProblems` returns healthy without
looking at the manifest that is still on disk (`lib/installer.js:1783-1795`). The new
regression calls this a genuine 2.0.0 marker by installing the current bidding add-on and
deleting those two fields (`test/bidding-addon.test.js:1448-1458`). But `git ls-tree v2.0.0`
confirms that release shipped neither `parley-bidding` nor any `parley-addon.json`.
`parley-bidding` therefore has no legitimate released schema-less marker to preserve.

I installed the current add-on, deleted both marker fields, left `parley-addon.json` present,
and modified `SKILL.md`. With Python 3.14 available, `doctor` returned **`ok:true`**,
`status:"valid"`, and `problems:[]`. Cycle 14 moved the silent downgrade from one deleted
field to two; it did not close it.

Legacy compatibility must be scoped to units that could actually have a manifest-free legacy
install. A manifest-required source such as `parley-bidding` must treat an owned marker with
no schema/manifest declaration as malformed and require reinstall. The compatibility
regression should use a skill that really shipped in 2.0.0, while a bidding regression deletes
both fields, tampers with payload bytes, and requires a non-green result.

### [MINOR] An unreadable manifest-covered file makes health throw instead of reporting malformed

`verifyPayload` catches `lstatSync` failures, but the subsequent `hashFile(abs)` call is
outside a `try` (`lib/addon-manifest.js:177-192`). After a clean install I changed only the
installed bidding `SKILL.md` to mode 000 and called `doctorCommand`. It threw raw `EACCES`
instead of returning a target/unit result. Through the CLI wrapper this becomes a bare
one-line error and exit 1, so JSON consumers receive no health document and cannot identify
the unit as malformed.

Catch read/hash failures per declared file and append an `unreadable` integrity problem, as
the function's list-returning contract promises. Add health and source-preflight regressions.
This is MINOR rather than MAJOR because it fails non-green, but `strict_gate: true` means it
still prevents this round from closing.

## Release judgement

No. Commit `12f9071` is not releasable as 2.1.0. The one release gate is to make the
installer's marker-derived path scope, removal transaction, and manifest validation fail
closed for all four cases above, then obtain a fresh full-scope review with no findings.

## What I verified

- Read `FINAL.md`, the complete feature diff from `v2.0.0`, both requested diffs
  (`9ed2081..12f9071` and `49fc3ec..12f9071`), implementation cycles 8-14, the void-round
  record, and all relevant prior review findings and dispositions.
- Ran `npm test` at `12f9071`: **332 Node tests passed**, the Python leg passed **54/54**
  across seven files on Python 3.14, and the 47-file manifest matched aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- Ran the Node suite again with `/usr/bin/python3` 3.9 first on `PATH`: **332 passed, 0
  failed**. I did not claim the Python leg passed there; it refuses 3.9 by design.
- Ran `npm pack --dry-run --json` with an isolated writable npm cache after the user's
  default cache refused writes: **202 entries**, including all **48** bidding files and no
  cache artefacts or nested `.gitignore`.
- Ran the adapter validator separately for all four adapters: 4/4 valid. Parsed all 16 JSON
  files, including four schemas. The integrated tree has 48 files and the read-only source
  has 48: source-only `.gitignore`, integrated-only `parley-addon.json`, and the nine recorded
  content-different files for rename/schema/consent integration.
- Scanned the production Python tools for network, browser, subprocess, and portal-mutation
  dependencies; none are present. The full Python tests exercised the HITL, evidence,
  exact-byte, ambiguous-result, authority, pricing, adapter-ceiling, and upload/submission
  boundaries.
- Independently reproduced every new finding above in temporary homes at exact `12f9071`,
  including the **83 removed / 1 failed** immutable-file arm, the out-of-root forced deletion,
  the fully green two-field marker downgrade, the removable mode-0555 empty directory, and
  the unreadable-payload health exception.
- The tracked worktree started clean and remained unchanged at
  `12f9071239dc6bc0da3173001fa49ef0eb500dc5`. During the review a concurrent process created
  untracked `tmp_edge_test.js`; I did not read it as reviewed code, execute it deliberately,
  modify it, or remove it. No tracked file or reviewed commit moved.
