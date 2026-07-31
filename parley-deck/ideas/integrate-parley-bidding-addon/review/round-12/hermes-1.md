---
idea: integrate-parley-bidding-addon
review-round: 12
agent: hermes-1
date: 2026-07-31
reviewed-commit: 5100f34
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All four round-11 findings are closed. I verified each by independent
reproduction at `5100f34` in isolated temporary homes, not by reading the
diff alone.

**CRITICAL — marker path confinement.** `markerAddonNames` now returns
`{ names, problem }` and validates every entry through `unusableAddonName`:
must be a string, not `.` or `..`, no `/` or `\`, must match
`^[A-Za-z0-9][A-Za-z0-9._-]*$`, no duplicates. Any violation sets `problem`,
which propagates to the core unit's `markerProblem` field and makes
`expectedAddonNames` return `[]` — so no unit is constructed from the bad
name at all. Behind that, `targetSkillUnits` resolves each dest with
`path.resolve(skillsDir, name)` and checks `path.dirname(dest) !== skillsRoot`
before creating a unit. I verified: the `../../outside-sentinel` sentinel
survives forced uninstall, no unit is constructed for the out-of-root name,
health reports `malformed`, and the read paths (`doctor`, `paths`) no longer
list the outside path. The install path was already safe (derives from flags,
not marker) and remains so.

**MAJOR — removability preflight.** `firstRemovalObstacle` is deleted
entirely. Uninstall is now a two-phase transaction in `removeFleetAtomically`:
phase A renames every destination in the whole plan aside (fleet-wide, not
per-target), rolling back every rename on any failure; phase B deletes the
quarantined trees, where a failure is a warning, not a fleet failure. I
measured the `uchg` arm that was the round-11 MAJOR: `chflags uchg` on
`parley-bidding/SKILL.md` in a single-target uninstall — 6/6 removed, 0
failures, 1 warning naming the quarantine debris. Fleet-wide with `uchg` on
the last target's bidding: 84/84 removed, 0 failures, 1 warning. The empty
0555 directory false positive is also gone — no longer refused. The
self-regression (per-target quarantine letting 78 units be removed before the
14th refused) is fixed: phase A spans the whole plan, and the
`blockedAnywhere` gate prevents any quarantine when any unit is blocked.

**MAJOR — legacy exemption scope.** `manifestProblems` now takes a
`sourceHasManifest` parameter. The legacy exemption (no `markerSchema` and
no `manifest`) only fires when `sourceHasManifest` is false — i.e., the
packaged source genuinely ships no manifest. `parley-bidding` (which ships a
manifest and did not ship in 2.0.0) with both fields deleted now reports
`malformed` with "install marker predates payload manifests, but this skill
ships one; re-install to validate it". `parley-worktrees` (which shipped in
2.0.0 and ships no manifest) keeps the exemption. The cycle-14 test was
rewritten to use `parley-worktrees` as the legacy subject, and a new test
confirms `parley-bidding` with both fields deleted and a tampered payload is
`malformed`.

**MINOR — hashFile EACCES.** `hashFile` is now inside a try/catch in
`verifyPayload`. A mode-000 declared file is reported as
`unreadable (EACCES): SKILL.md` in the problems list, not thrown out of
`doctorCommand`. I verified: `chmod 000` on installed `parley-bidding/SKILL.md`,
`doctor` returns `malformed` with the expected problem string.

The two recorded follow-ups remain unchanged and correctly deferred: only
`parley-bidding` ships a manifest (B3.11), and kimi-1's round-9 NIT on the
`dirExists` discovery guard.

## New findings

None.

I attacked the quarantine transaction and the confinement as the round
instructed. The areas examined, with the outcome:

**Rollback failure.** When a quarantine rename fails, the rollback loop
renames every already-quarantined unit back to its original dest. If
rollback itself fails for a unit, that unit is reported `failed` with the
aside path so an operator can find it; all other pending units get
`skipped`. No deletions occur. I traced this by monkey-patching
`fs.renameSync` in the core-name-as-addon case (two units with the same
dest): first rename succeeds, second fails with ENOENT, rollback succeeds,
fleet blocked, core preserved, zero `.removing` debris. Fail-closed.

**Name collisions.** `quarantineName` uses
`.${basename}.${process.pid}.${Date.now()}.removing`. Different basenames
get different quarantine names even in the same millisecond. Same-basename
duplicates (a corrupt marker listing the core name as an addon) produce
different timestamps, but the second rename fails with ENOENT because the
source is already gone — triggering rollback. I also tested a pre-existing
directory at a predicted quarantine name: the actual `Date.now()` at
quarantine time differed, so the rename succeeded and the uninstall
completed normally.

**Interaction with install's backup path.** Install uses `.bak` suffix,
quarantine uses `.removing` suffix — no collision. Both treat cleanup
failures as warnings after a committed operation. I verified that leftover
`.removing` debris from a previous failed uninstall does not affect a
subsequent install.

**--dry-run.** The outer preflight is skipped for dry-run (line 680), and
`removeFleetAtomically` handles dry-run at line 1418 — no quarantine, no
deletion, just a `remove` action with `dryRun: true`. I verified no
filesystem mutation occurs, no `.removing` files appear, and a corrupt
marker still fails closed in dry-run mode.

**Symlinked destinations.** The quarantine rename operates on the directory
entry (the symlink itself), not its target. I verified: a symlinked
destination is renamed aside, then `rmSync` deletes the symlink-aside, and
the real target is preserved. A symlinked skills directory also works
correctly for both install and uninstall. Dangling symlinks are handled by
the existing `pathEntryExists` (lstat-based) and are removed by forced
uninstall.

**Concurrent runs.** `quarantineName` includes `process.pid`, so two
processes get different quarantine names. If process A quarantines a
destination, process B finds it missing (`pathEntryExists` returns false)
and reports `missing` — a no-op, not a deletion or error. No partial fleet.

**Other routes from stored data to a path.** I checked every path from
mutable data to filesystem operations: the marker `addons` array (now
validated and confined), the `--only` flag (validated against discovered
addon names, which are real directory entries), the read-path discovery
loop (uses on-disk directory names, which already exist), `skillsDir`
itself (derived from hardcoded `TARGETS` definitions and `HOME` env, not
stored data), and the marker's `skill` and `target` fields (used for
comparison only, not path construction). No other route exists.

**Two observations that are not findings.** A marker listing the core skill
name as an addon (`addons: ["parley-deck"]`) passes name validation and
creates two units with the same dest. On uninstall, the quarantine
transaction correctly rolls back and blocks the fleet. On health, the
duplicate addon unit reads the core marker and reports `malformed` (core
marker has no `manifest` field, expected for core but not addon). Both
paths are fail-closed. Similarly, a case-insensitive duplicate
(`["parley-bidding", "Parley-bidding"]`) passes the string-based duplicate
check but resolves to the same path on macOS APFS; health reports the second
as `malformed`, and uninstall rolls back. Both are corrupt marker states
handled safely.

## Release judgement

Releasable as 2.1.0. All four round-11 findings are closed and verified by
independent reproduction. The quarantine transaction is sound across
rollback, name collisions, dry-run, symlinks, and concurrent runs. The
confinement has no other route from stored data to path construction. The
known limit (phase B debris invisible to `doctor`) is recorded in
IMPLEMENTATION.md. The two deferred follow-ups remain correctly recorded.

## What I verified

- **Tree discipline.** `git status --porcelain` clean before and after all
  work. HEAD stayed at `5100f34`. All probes ran in a `git archive` scratch
  copy under `$TMPDIR`. No edits under `skills/parley-bidding/`, no tree
  mutation, no scratch files in the repo root.

- **Suites at `5100f34`.** `npm test` — 338 node tests, 0 fail. Python leg
  — 54/54 across 7 files on Homebrew python3 3.14.6. Node leg re-run with
  `/usr/bin/python3` 3.9.6 on PATH: 338/338 (Python leg refuses 3.9 by
  design). Manifest check — 47 files, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`,
  unchanged since `714712f`.

- **Discrimination.** Copied the tree, checked out `12f9071`, ran the
  cycle-15 test file against the cycle-14 installer: 9 of 11 cycle-15
  tests fail, 2 pass (the cycle-14 preserved properties — "marker that kept
  its manifest but lost its schema" and "legacy marker with no schema keeps
  validating"). At `5100f34`: 11/11 pass. Consistent with the
  IMPLEMENTATION.md claim of 3/11 at `12f9071` (the third is the rewritten
  legacy-marker test, which changes subject from `parley-bidding` to
  `parley-worktrees` and passes at both because the exemption itself is
  unchanged for manifest-free skills).

- **CRITICAL reproduction.** Installed Codex target, corrupted core marker's
  `addons` to `["../../outside-sentinel"]`, created sentinel with KEEP file,
  ran forced uninstall. `ok:false`, sentinel survives, no unit constructed
  for the out-of-root name, core preserved. Also verified `doctor` and
  `paths` no longer list the outside path. Verified `--only` is safe
  (`validateAddonSelection` rejects unknown names). Verified install is safe
  (uses `selectedAddons` from flags, not marker).

- **MAJOR (uchg) reproduction.** `chflags uchg` on
  `parley-bidding/SKILL.md`, unforced uninstall: 6/6 removed, 0 failures, 1
  warning. Fleet-wide with `uchg` on last target: 84/84 removed, 0 failures,
  1 warning. Also verified empty 0555 directory is no longer refused.

- **MAJOR (legacy exemption) reproduction.** `parley-bidding` with both
  `markerSchema` and `manifest` deleted, tampered payload: `malformed` with
  "predates payload manifests, but this skill ships one". `parley-worktrees`
  with both fields deleted: `valid` (exemption preserved).

- **MINOR reproduction.** `chmod 000` on installed `parley-bidding/SKILL.md`:
  `doctor` returns `malformed` with `unreadable (EACCES): SKILL.md`, no
  thrown exception.

- **Quarantine transaction edge cases.** Traced rename sequence with
  monkey-patched `fs.renameSync` for normal uninstall (6 clean renames, 6
  clean deletes) and core-name-as-addon (first rename succeeds, second
  ENOENT, rollback succeeds, fleet blocked). Verified dry-run leaves no
  mutations. Verified symlinked destination: symlink removed, real target
  preserved. Verified symlinked skills directory: install and uninstall
  work through the symlink. Verified dangling symlink: removed by forced
  uninstall. Verified file (non-directory) destination: removed by forced
  uninstall. Verified quarantine name collision with pre-existing `.removing`
  directory: different `Date.now()` avoids collision. Verified concurrent
  simulation: second process reports `missing`, no deletion.

- **Confinement completeness.** Checked all paths from mutable data to
  filesystem operations: marker `addons` (validated + confined), `--only`
  (validated against discovered names), read-path discovery (real directory
  entries), `skillsDir` (hardcoded targets + HOME env), marker `skill`/
  `target` (comparison only). No other route.

- **Name validation edge cases.** Tested 19 entries against the validator:
  valid names (simple, with dots, dashes, underscores, digits, uppercase)
  pass; invalid names (`.`, `..`, empty, separators, leading/trailing
  spaces, leading dot, exclamation, colon) are rejected. Very long names
  (300 chars) pass validation but fail at the filesystem level on uninstall
  — the quarantine rename fails with ENAMETOOLONG, triggering rollback and
  blocking the fleet. Fail-closed.

- **Code read in full.** `lib/installer.js` (2295 lines), `lib/addon-manifest.js`
  (237 lines), the cycle-15 diff (`12f9071..5100f34`), the full diff from
  round 7 (`49fc3ec..5100f34`), IMPLEMENTATION.md cycles 10-15, FINAL.md B5
  text, all three round-11 reviews, and the design-addons.test.js changes.

- **Not running as root.** uid 501. All permission and flag arms genuinely
  exercised.
