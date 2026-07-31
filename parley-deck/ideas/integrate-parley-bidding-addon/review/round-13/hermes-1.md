---
idea: integrate-parley-bidding-addon
review-round: 13
agent: hermes-1
date: 2026-07-31
reviewed-commit: dd8d756
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All five round-12 findings are closed at `dd8d756`. I verified each by
independent reproduction in isolated temporary homes, not by reading the
diff alone.

**Round-12 codex-1 CRITICAL — unknown marker name deletes an unrelated
sibling: CLOSED.** `targetSkillUnits` now defines `authorize(name, dest)`
which runs when names come from the marker (`fromMarker` guard at
`lib/installer.js:1124-1126`). A recorded name is accepted only when the
package ships that add-on (`discovered.has(name)`) or the destination
carries this installer's marker claiming that identity
(`installerOwnsDestination(dest, name)`). I reproduced the attack: installed
codex, created `unrelated-sentinel` with a KEEP file, set the core marker's
`addons` to `["unrelated-sentinel"]`, ran `uninstall --force`. Result:
`ok:false`, sentinel survives, core survives, message matches
`/does not ship and does not own/`. The guard regression (dropped add-on
still uninstallable via its own marker) also passes.

**Round-12 codex-1 MAJOR — non-array `addons` container types: CLOSED.**
`markerAddonNames` now checks `marker.addons === undefined || marker.addons
=== false` for core-only, and any other non-array value produces a problem
(`lib/installer.js:1041-1048`). I tested all five shapes codex-1 named
(`"parley-bidding"`, `true`, `null`, `{}`, `42`): each reports `malformed`
in doctor with `/neither a list nor false/`, and `uninstall --force` returns
`ok:false` with the core preserved. The absent-field and explicit-`false`
guards still report `valid`.

**Round-12 codex-1 MAJOR — manifest keys escape the payload: CLOSED.**
`unusableManifestKey` rejects backslashes, absolute paths, drive letters,
empty segments, and `.`/`..` segments (`lib/addon-manifest.js:174-185`).
`insidePayload` checks that the resolved path is a strict descendant of the
payload root (`lib/addon-manifest.js:189-193`). I tested four escape keys
(`../outside-sentinel`, `../parley-deck/SKILL.md`, `/etc/passwd`,
`a/../../b`), each with a correct digest and recomputed aggregate: all
return `ok:false` with `/unusable manifest entry|escapes the payload/`.

**Round-12 kimi-1 MAJOR — install leaves a partial fleet after commit
failure: CLOSED for single-process operation.** `installFleetAtomically`
(`lib/installer.js:1458-1553`) stages every unit before committing any,
and a commit failure at any unit reverts every earlier commit by rename
within the same parent. I reproduced the regression scenario: installed to
all 14 targets, `chflags uchg` on `aionrs`'s core dest, re-installed to
all. Result: `ok:false`, 0 units with action `installed`/`replaced`,
untouched target byte-identical, no `.tmp`/`.bak` debris in any skills
directory. The committed-then-reverted units report `skipped` (see below).

**Round-12 kimi-1 NIT — discrimination count: CLOSED.** The implementation
record separates four proof-of-fix regressions from two over-correction
guards.

**Recorded follow-ups: still open and non-blocking on their previously
agreed terms.** B3.11 (only `parley-bidding` ships a manifest), the
`dirExists` discovery-guard NIT, quarantine debris invisible to `doctor`,
and residual-disposal arms (`uappnd`, delete-denying ACLs) producing debris
rather than a partial fleet.

## New findings

None that block the 2.1.0 release. I reproduced all three of codex-1's
round-13 findings; my assessment of each follows, with the severity I
assign and why.

### [MINOR] Symlinked manifest trusted as payload authority (codex-1 MAJOR)

Reproduced. After a normal install I replaced `parley-bidding/parley-addon.json`
with a symlink to a byte-identical external file. `verifyPayload` returned
`ok:true`; `doctor` reported `valid` and `managed`. Health depends on bytes
outside the destination directory.

The payload walker (`listPayloadFiles`) uses `lstatSync` and refuses
symlinks for every payload file, but the manifest file is explicitly
excluded by `isIgnored` (`lib/addon-manifest.js:28`) and is read through
`readFileSync` (follows symlinks) in `readManifest`, `manifestFileHash`,
and `hasManifest` (which uses `statSync`). This is a consistency gap: the
key confinement cycle 16 added guards the keys the manifest supplies, but
not the manifest file that supplies them.

I rate this MINOR rather than MAJOR for three reasons. First, it is
pre-existing: `hasManifest` has used `statSync` since `714712f` and the
manifest has been excluded from the symlink check since the same commit.
Cycle 16 did not open this door. Second, the module's own scope statement
says "this is defect detection, not tamper resistance" and "anyone who can
rewrite the payload can rewrite the manifest beside it"
(`lib/addon-manifest.js:12-13`). Replacing the manifest with a symlink
requires write access to the destination directory, which is the same
access level the threat model already disclaims. Third, a modification to
the external file IS caught by the marker's `sha256` comparison
(`manifestProblems` line 2177) — I verified this. The gap is that a
byte-identical external file is trusted, not that a modified one escapes.

Recording as a follow-up: the manifest file should be lstat-checked and
required to be a regular file, for consistency with the payload walker.
This closes the last file in a destination directory that is read as truth
rather than input.

### [MINOR] Aliased physical destinations (codex-1 MAJOR)

Reproduced. I symlinked `.codex` to `.hermes` and ran `install --target all
--include-undetected`. Both targets resolved to the same physical skills
directory. Hermes (later in the TARGETS array) replaced codex's install;
the physical destination carries hermes's marker. `doctor` reports both as
`valid` because ownership checks `marker.name` and `marker.skill`, not
`marker.target`. Uninstalling codex deletes the shared physical directory,
making hermes `missing`.

I also reproduced the different-payload alias: `.gemini/config/plugins`
and `.gemini/extensions` both symlinked to the same directory. Agy
(antigravity kind) committed first, Gemini (gemini kind) committed second
and overwrote Agy's payload. Doctor correctly reports Agy as `malformed`
(missing `skills/SKILL.md`) and Gemini as `valid`.

I rate this MINOR rather than MAJOR. It is pre-existing: at cycle 9
(`ebe269e`), install was `targets.map(installTarget)`, which had the same
alias behaviour — each target independently writes to its dest, and aliased
dests overwrite. Cycle 16's fleet transaction did not introduce this. The
scenario requires the user to have symlinked their runtime configuration
directories so that two different runtimes resolve to the same physical
skills directory, which is an unusual system configuration. For same-kind
targets (codex/hermes), the payload is byte-identical and only the
informational `target` field differs. For different-kind targets
(agy/gemini), doctor catches the mismatch as `malformed`. The uninstall
cross-deletion is a natural consequence of sharing a physical directory —
if two targets alias, deleting for one necessarily affects the other.

The fix codex-1 proposes (canonicalize destinations, reject aliases, acquire
locks) is reasonable for a future hardening cycle but is not a release gate
for 2.1.0, which does not claim transaction isolation across aliased
destinations or concurrent processes.

### [MINOR] Dry-run bypasses fleet-wide preflight (codex-1 MINOR)

Independently reproduced. I installed a valid core-only codex target,
corrupted the marker's `addons` to `{}`, and ran `install --dry-run`. It
returned `ok:true` with action `replace`. The real install returned
`ok:false` with `blocked` and the unusable-selection message. The same
discrepancy exists for uninstall: a corrupted marker in one target causes
real uninstall to skip all targets, but dry-run uninstall reports the other
targets as `remove` (ok:true).

This is pre-existing: the fleet-wide preflight was added in cycle 10
(`3553f47`) and was gated on `!dryRun` from the start. The dry-run path
falls through to `installTarget`/`removeFleetAtomically` which do per-target
preflight, not the fleet-wide gate. Dry-run does not write anything, so no
partial fleet results. The issue is fidelity: dry-run cannot be relied on
to predict whether the real command will succeed.

Recording as a follow-up: run the same read-only fleet preflight for dry-run
and report the same blocked/skipped plan; only staging and commit should be
omitted.

### Observations that are not findings

**Reverted units report "Not attempted".** When a commit failure triggers
revert, successfully reverted units are not in `results`, so `skipRest`
gives them action `skipped` with message "Not attempted: another skill or
target in this install could not be committed." They WERE attempted —
staged, committed, then reverted. The message is misleading but the fleet
state is correct: I verified the reverted unit's marker timestamp is the
original, not the new install's, and no staging debris remains. Cosmetic,
not a finding.

**No other route from stored data to path/verdict/mutation.** I traced
every field of the install marker and the manifest:

- `marker.addons` — validated by `unusableAddonName`, confined to direct
  child of skills dir, authorized by `authorize()`. Closed.
- `marker.manifest.aggregate` / `marker.manifest.sha256` — strings compared
  for equality, not used to construct paths.
- `marker.skill` — compared for equality in `installerOwnsDestination`, not
  used to construct paths.
- `marker.target` — informational only, not read back for any path or
  ownership decision.
- `marker.scope` — informational only, not read back.
- `marker.markerSchema` — integer compared, not a path.
- `manifest.files` keys — validated by `unusableManifestKey` and confined by
  `insidePayload`. Closed.
- `manifest.runtime.python` — parsed by anchored regex `/^>=\s*(\d+)\.(\d+)$/`,
  used to probe `python3` on PATH, not a path.
- `manifest.aggregate` — string compared for equality.

No other stored value reaches a filesystem path, a health verdict, or a
mutation.

## Release judgement

Releasable as 2.1.0. Cycle 16's four claimed fixes are correctly implemented
and verified by independent reproduction. The install transaction is sound
for single-process operation: staging failures leave nothing, commit
failures revert every earlier commit by rename within the same parent, and
backup disposal failures are warnings. The three findings codex-1 raised
are all pre-existing (none was introduced by cycle 16), all are MINOR
rather than MAJOR under the module's stated threat model, and none creates
a route from stored data to an uncontrolled filesystem mutation. They are
recorded as follow-ups for a future hardening cycle, not release gates.

## What I verified

- **Tree discipline.** `git status --porcelain` clean before and after all
  work. HEAD stayed at `dd8d756`. All probes ran in isolated temporary
  homes under `$TMPDIR`. No edits under `skills/parley-bidding/`, no tree
  mutation, no scratch files in the repo root.

- **Suites at `dd8d756`.** `node --test` — 344 tests, 0 fail. Manifest
  check — 47 files, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`,
  unchanged since `714712f`.

- **Code read in full.** `lib/installer.js` (2483 lines),
  `lib/addon-manifest.js` (273 lines), the cycle-16 diff
  (`5100f34..dd8d756`), the full round-7 delta (`49fc3ec..dd8d756`), all
  three round-12 reviews, and codex-1's round-13 review.

- **Round-12 CRITICAL reproduction.** Installed codex, created
  `unrelated-sentinel` with KEEP file, set core marker `addons` to
  `["unrelated-sentinel"]`, ran `uninstall --force`. `ok:false`, sentinel
  survives, core survives, message matches `/does not ship and does not
  own/`. Guard regression (dropped add-on uninstallable via own marker)
  also passes.

- **Round-12 container-types reproduction.** Tested all five non-array
  shapes (`"parley-bidding"`, `true`, `null`, `{}`, `42`): each reports
  `malformed` in doctor with `/neither a list nor false/`, and
  `uninstall --force` returns `ok:false` with core preserved. Absent-field
  and explicit-`false` both still `valid`.

- **Round-12 manifest-key reproduction.** Tested four escape keys
  (`../outside-sentinel`, `../parley-deck/SKILL.md`, `/etc/passwd`,
  `a/../../b`), each with correct digest and recomputed aggregate: all
  return `ok:false` with `/unusable manifest entry|escapes the payload/`.
  Also tested null-byte key (`"SKILL.md\0../outside"`) — rejected as
  `missing` (no escape, no crash).

- **Round-12 install-fleet reproduction.** Installed to all 14 targets,
  `chflags uchg` on `aionrs` core dest, re-installed. `ok:false`, 0 units
  with action `installed`/`replaced`, untouched target byte-identical, no
  staging debris. Reverted units' marker timestamps are original, confirming
  full rollback.

- **Install transaction edge cases.** Verified staging temp names are
  unique within a fleet (different basenames for different skills, different
  parents for different targets). Verified `commitStagedUnit` error path
  restores backup when temp→dest rename fails. Verified `revertStagedUnit`
  for both replaced and fresh-install cases. Verified `discardBackup`
  failure produces a warning, not a failure. Verified no `.tmp`/`.bak`
  debris after successful fleet install.

- **Symlinked destinations.** Install into a symlinked skills directory:
  works, skill accessible through both symlink and real path. Install
  `--force` on a symlinked destination: symlink replaced with real
  directory. Uninstall of symlinked destination: symlink removed, real
  target preserved (from round-12 verification, still holds).

- **Dry-run discrepancy.** Confirmed for both install and uninstall: dry-run
  reports per-target results, real command applies fleet-wide gate. Tested
  with corrupted marker (`addons: {}`): dry-run `ok:true` action `replace`,
  real install `ok:false` action `blocked`. Pre-existing since cycle 10.

- **Symlinked manifest.** Reproduced: byte-identical external manifest
  reports `valid` and `managed`. Modified external manifest caught by
  marker `sha256` comparison (`malformed`, "has been replaced since
  installation"). Pre-existing since `714712f`.

- **Aliased destinations.** Reproduced two cases: `.codex`→`.hermes` symlink
  (same payload kind, marker `target` field differs, both report `valid`,
  uninstall cross-deletes); `.gemini/config/plugins` and `.gemini/extensions`
  both symlinked to same dir (different payload kinds, gemini overwrites
  agy, doctor catches agy as `malformed`). Pre-existing.

- **Authorization on read-only commands.** Verified `authorize()` runs for
  `doctor`/`status` (not just `uninstall`): a marker claiming a non-shipped,
  non-owned add-on makes the core `malformed` in doctor. A marker claiming a
  shipped add-on that is absent reports the add-on as `missing`. A marker
  claiming a non-shipped add-on that carries our marker (dropped-add-on
  case) reports both as `valid`. All correct.

- **Stored-data completeness.** Traced every marker field (`addons`, `skill`,
  `target`, `scope`, `markerSchema`, `manifest.aggregate`, `manifest.sha256`)
  and every manifest field (`files` keys, `aggregate`, `runtime.python`).
  No route to an uncontrolled filesystem path, health verdict, or mutation
  beyond what cycle 16 already addressed.

- **Not running as root.** uid 501. All permission and flag arms genuinely
  exercised.
