---
idea: integrate-parley-bidding-addon
review-round: 11
agent: hermes-1
date: 2026-07-31
reviewed-commit: 12f9071
---

## Verdict

BLOCK

## Outstanding findings — closed or not

Cycle 14's seven regressions are closed and verified: frozen destination trees
are blocked before install and uninstall, a foreign marker in the last uninstall
target deletes nothing, a committed replacement survives a failed backup cleanup,
and a manifest-bearing marker that lost only `markerSchema` is malformed. The
corrected Python-leg wording now distinguishes the 3.14 run from the deliberate
3.9 refusal. All seven regressions pass at `12f9071`; I confirmed 0/7 at
`9ed2081` by reading the diff, and 7/7 by running them.

Round 10's three independent MAJORs — codex-1's replacement-cleanup, codex-1's
one-field marker downgrade, and my own fleet-wide-uninstall — are closed for the
specific scenarios that were filed. The broader properties they pointed at are
**not closed**, and the new findings below are the doors those properties left
open. I say this plainly because the pattern matters: each of the last five
cycles closed the exact scenario it was given and left the class standing.

The two recorded follow-ups remain follow-ups, not findings: only `parley-bidding`
ships a manifest (B3.11, first consensus follow-up), and the `dirExists` discovery
guard for dangling symlinks (kimi-1 round-9 NIT). I re-confirmed both.

## New findings

Four findings. I independently reproduced every one at `12f9071` in isolated
temporary homes before writing them. `codex-1` has already filed the same four;
this review confirms them from independent measurement, not from reading
`codex-1`'s file.

### [CRITICAL] A core marker can steer forced uninstall outside the skills directory

`markerAddonNames` returns the core marker's `addons` array without validating
its values (`lib/installer.js:947-955`). `expectedAddonNames` treats those values
as the uninstall selection, and `targetSkillUnits` passes each one directly to
`path.join(skillsDir, name)` (`lib/installer.js:1010`). With `--force`,
`preflightUninstallUnit` deliberately skips the ownership check
(`lib/installer.js:729`), so traversal components supplied by a corrupt marker
become deletion targets outside the runtime's skills directory. The new
fleet-wide preflight makes this worse in one precise sense: it validates that
the marker-derived out-of-scope tree *can be removed*, then authorizes the plan.

Reproduced at `12f9071` in an isolated home:

1. Installed the normal Codex target (6 units, ok:true).
2. Changed only the core marker's `addons` to `["../../outside-sentinel"]`.
3. Created `$HOME/outside-sentinel/KEEP`.
4. Ran `uninstallCommand` with `--force`, no explicit `--only`.

Result: `ok:true`. The core was removed, and a skill literally named
`../../outside-sentinel` was removed, with its destination resolved to
`$HOME/outside-sentinel`. The sentinel tree was deleted. `KEEP` is gone.

The traversal also reaches the read paths: `doctorCommand` and `pathsCommand`
both list `$HOME/outside-sentinel` as a skill with its resolved destination.
`statusCommand` reports it as `malformed`. The install path is safe —
`installCommand` derives expected add-ons from `selectedAddons(context)` (the
flags), not from the marker — so only read and uninstall commands are affected.

I verified this predates cycle 14: `markerAddonNames` at `49fc3ec` has the same
unvalidated return. It is not a regression cycles 10-14 introduced; it is the
door they never looked at because none of them examined what the marker's
`addons` field can contain. `--only` is safe — `validateAddonSelection`
(interjection of new content: validated against discovered addon names before
path construction, so `--only "../../foo"` throws "Unknown add-on(s)"). The
marker bypasses that validation entirely.

`--force` may override ownership at a caller-selected skill destination; it must
not let mutable marker data expand the command's path scope. Before release,
recorded selections must be validated before path construction, and every
derived destination must be confined to an exact direct child of `skillsDir`.
Invalid, unknown, non-string, duplicate, `"."`, `".."`, or
separator-containing entries must make health and mutation fail closed. A
regression must preserve an out-of-root sentinel under forced uninstall.

This blocks the 2.1.0 release.

### [MAJOR] The removability preflight is not a mirror of rmSync — a predictable partial uninstall remains

`firstRemovalObstacle` checks `R_OK | W_OK | X_OK` on every directory in the
destination tree (`lib/installer.js:1971`). This is an access-permission walk,
not a deletion-feasibility walk. It has false negatives and false positives
relative to `fs.rmSync`:

**False negative (B5 violation):** On macOS, a user-immutable (`uchg`) regular
file retains ordinary mode bits — `accessSync` passes, but `rmSync` fails. I
installed all 84 units across 14 targets, set `chflags uchg` on
`~/.aionrs/skills/parley-worktrees/SKILL.md` (the last unit of the last target),
and ran an unforced fleet uninstall. Fleet preflight passed. The command removed
**83 units**, then the last unit failed with `ENOTEMPTY` and left the flagged
file. The earlier 13 targets' cores are gone. This is the same B5 partial-fleet
result cycle 14 was meant to exclude, reachable through `chflags` rather than
`chmod`. The failure is 100% predictable and was not preflighted.

**False positive (user-facing regression):** An empty mode-0555 directory inside
a destination is blocked by `firstRemovalObstacle` (it requires `W_OK` on the
directory itself), but `rmSync` succeeds because removing an empty directory
needs write-and-search on its parent, not write inside it. I verified both
directions directly: `rmSync` on a tree with an empty 0555 dir succeeds; on a
tree with a non-empty 0555 dir it fails with `ENOTEMPTY`; on a tree with a
`uchg` file it fails with `ENOTEMPTY`. The helper is neither a sound
over-approximation nor a sound under-approximation of deletion semantics.

No shipped skill currently contains empty directories (verified by walking
`skills/parley-bidding` and `skills/parley-deck`), so the false positive is
theoretical for the current payload. But a foreign tree at a destination path
could contain one, and `--force` install would block a replacement that
`rmSync` would handle.

The fix should not continue to equate an access walk with deletion semantics.
Uninstall needs a transaction/rollback or quarantine design that cannot delete
earlier units when a later cleanup fails — the install path's post-commit
warning already demonstrates the relevant distinction between logical commit
and cleanup. Regressions are needed for an immutable file in the final unit and
a removable empty read-only directory.

This blocks the 2.1.0 release. It is the same B5 class that cycles 10-14 closed
for mode-bit obstacles, reachable through a different filesystem attribute. I
note honestly: `chflags uchg` is less common than `chmod -R a-w`, and the
previous five doors each used progressively more exotic mechanisms. But the
contract says "a predictable failure must produce zero writes," and this one is
predictable.

### [MAJOR] Deleting both marker fields still silently disables manifest validation for parley-bidding

Cycle 14 narrowed the legacy exemption: a marker that keeps `manifest` but loses
`markerSchema` is now malformed. But a marker that loses **both** fields still
takes the legacy exemption (`lib/installer.js:1792-1794`), because the exemption
fires on `markerSchema === undefined && manifest === undefined`. The released
2.0.0 marker carried neither field — but `parley-bidding` did not ship in 2.0.0.
`git ls-tree v2.0.0 skills/` lists `parley-deck`, `parley-design`,
`parley-design-check`, `parley-tracker`, `parley-worktrees`. No `parley-bidding`.

So for `parley-bidding` — the one add-on that ships a manifest — deleting both
marker fields silently downgrades from byte validation to none. I reproduced
this at `12f9071`: installed `parley-bidding`, deleted both `markerSchema` and
`manifest` from its marker, left `parley-addon.json` on disk, and appended
`"\n# TAMPERED\n"` to `SKILL.md`. `doctor` returned `ok:true`, `status:"valid"`,
`problems:[]`. The tampered payload is invisible.

The new regression at `test/bidding-addon.test.js:1448` calls this a "genuine
2.0.0 marker" by installing the current bidding add-on and deleting both fields.
But `parley-bidding` has no legitimate released schema-less marker to preserve.
Cycle 14 moved the silent downgrade from one deleted field to two; it did not
close it.

Legacy compatibility must be scoped to units that could actually have a
manifest-free legacy install. A manifest-required source such as
`parley-bidding` must treat an owned marker with no schema/manifest declaration
as malformed and require reinstall. The compatibility regression should use a
skill that really shipped in 2.0.0 (e.g. `parley-worktrees`), while a bidding
regression deletes both fields, tampers with payload bytes, and requires a
non-green result.

This blocks the 2.1.0 release.

### [MINOR] An unreadable manifest-covered file makes health throw instead of reporting malformed

`verifyPayload` catches `lstatSync` failures (`lib/addon-manifest.js:180-185`)
but the subsequent `hashFile(abs)` call at line 190 is outside any try/catch.
`hashFile` (`lib/addon-manifest.js:58-60`) does
`crypto.createHash("sha256").update(fs.readFileSync(file)).digest("hex")` with
no error handling. A mode-000 manifest-covered file passes the `lstatSync` +
`isFile` checks (lstat needs search on the parent, not read on the file), then
throws raw `EACCES` from `readFileSync`.

Reproduced at `12f9071`: after a clean install, `chmod 000` on the installed
`parley-bidding/SKILL.md`, then `doctorCommand` — it threw `EACCES` instead of
returning a target/unit result. Through the CLI wrapper this becomes a bare
one-line error and exit 1, so JSON consumers receive no health document and
cannot identify the unit as malformed.

This is MINOR rather than MAJOR because it fails non-green (the right outcome),
but it violates the function's list-returning contract and gives no structured
output. Catch read/hash failures per declared file and append an `unreadable`
integrity problem, as the function's contract promises.

## Release judgement

Not releasable as 2.1.0. Four things must change:

1. **(CRITICAL)** Validate marker-derived addon names before path construction;
   confine every destination to an exact direct child of `skillsDir`. A corrupt
   marker must not expand the deletion or inspection scope outside the skills
   directory, even under `--force`. Regression: out-of-root sentinel survives
   forced uninstall.

2. **(MAJOR)** Make uninstall's removability model match `rmSync` or use a
   transaction/rollback design that cannot delete earlier units when a later
   cleanup fails. Regression: immutable file in the final unit produces zero
   deletions.

3. **(MAJOR)** Scope the legacy marker exemption to skills that actually shipped
   in 2.0.0. A `parley-bidding` marker with no `markerSchema` and no `manifest`
   must be malformed, not legacy. Regression: two-field deletion on a bidding
   marker with a tampered payload requires non-green.

4. **(MINOR)** Catch `hashFile` failures in `verifyPayload` and report them as
   integrity problems instead of throwing.

## What I verified

- **Tree discipline**: `git status --porcelain` clean before and after all work.
  HEAD stayed at `12f9071`. All probes ran in `$TMPDIR` homes. No edits under
  `skills/parley-bidding/`, no tree mutation. Temporary test files were cleaned
  up and the tree verified clean afterwards.

- **Suites at `12f9071`**: `node --test` — 332 tests, 0 fail. Python leg —
  54/54 across 7 files under Homebrew python3 3.14.6. Manifest check — 47 files,
  aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
  Node leg re-run under `/usr/bin/python3` 3.9.6: 332/332. Python leg under 3.9.6
  refuses by design. All confirmed by running them.

- **CRITICAL reproduction**: installed Codex target, corrupted core marker's
  `addons` to `["../../outside-sentinel"]`, created `$HOME/outside-sentinel/KEEP`,
  ran forced uninstall. `ok:true`, sentinel deleted. Also confirmed traversal
  reaches `doctor` (reports outside path as a skill), `paths` (lists it), and
  `status` (reports it malformed). Confirmed `--only` is safe
  (`validateAddonSelection` rejects unknown names through `run()`). Confirmed
  install is safe (uses `selectedAddons` from flags, not marker). Confirmed
  predates cycle 14 by reading `markerAddonNames` at `49fc3ec`.

- **MAJOR (uchg) reproduction**: installed all 84 units, `chflags uchg` on
  `~/.aionrs/skills/parley-worktrees/SKILL.md`, unforced fleet uninstall.
  83 removed, 1 failed, earlier targets' cores gone. Also verified the
  false-positive direction: empty 0555 dir is blocked by `firstRemovalObstacle`
  but `rmSync` succeeds. Verified no shipped skill contains empty directories.

- **MAJOR (two-field downgrade) reproduction**: installed `parley-bidding`,
  deleted both `markerSchema` and `manifest` from marker, tampered with
  `SKILL.md`. `doctor` returned `ok:true, status:valid, problems:[]`. Verified
  `parley-bidding` did not ship in v2.0.0 (`git ls-tree v2.0.0 skills/`).

- **MINOR reproduction**: `chmod 000` on installed `parley-bidding/SKILL.md`,
  `doctorCommand` threw `EACCES`. Confirmed `hashFile` at
  `lib/addon-manifest.js:58-60` has no try/catch, and the call at line 190 is
  outside the `lstatSync` try/catch.

- **Code read in full**: `lib/installer.js` (2187 lines), `lib/addon-manifest.js`
  (225 lines), the cycle 14 diff (`9ed2081..12f9071`), the full diff from round 7
  (`49fc3ec..12f9071`), IMPLEMENTATION.md cycles 10-14, FINAL.md B5 text, all
  three round-10 reviews, and codex-1's round-11 review.

- **Not running as root**: uid 501. All permission arms genuinely exercised.
