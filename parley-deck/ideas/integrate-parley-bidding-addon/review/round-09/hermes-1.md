---
idea: integrate-parley-bidding-addon
review-round: 9
agent: hermes-1
date: 2026-07-31
reviewed-commit: dcd200e
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

### codex-1 round-8 MAJOR (F26) — destination preflight was per-target and followed symlinks — CLOSED

Both halves of the finding are fixed and independently verified at dcd200e.

**Fleet-wide preflight.** `installCommand` (lib/installer.js:622-653) now builds the
complete target × unit plan and preflights every unit before the first write. A single
blocker anywhere returns every unit as `blocked` or `skipped` with zero writes. I
reproduced codex-1's original scenario — an unmarked `parley-bidding` destination in the
last target (`aionrs`) under `--target all --include-undetected` — and confirmed that no
earlier target's `parley-deck` directory was written. The per-target preflight inside
`installTarget` (lines 1122-1145) is now a redundant defense-in-depth check; it cannot
block because the fleet-wide gate already passed the same conditions, but it is harmless
and provides a second safety net.

**lstat-aware destination checks.** `pathEntryExists` (lib/installer.js:1788-1795) replaces
`fs.existsSync` on every destination-entry check: install preflight, `installSkillUnit`,
the backup/replace path in `copyPayloadAtomically`, `skillUnitStatus` (health), and
`uninstallSkillUnit`. `lstatSync` stats the link entry itself rather than following it, so
a dangling symlink is seen as present rather than silently passing preflight and failing
with `ENOTDIR` mid-fleet. Only `ENOENT` is treated as absence; any other error (EACCES,
etc.) is treated as present, which is correct — a permission-denied entry is not "missing"
and should not invite an install that then fails.

### Cycle 11 — `skillUnitStatus` and `uninstallSkillUnit` still used `existsSync` — CLOSED

Cycle 10 converted only the install path. `skillUnitStatus` (line 1440) still called a
dangling destination symlink `missing`, and `uninstallSkillUnit` (line 1269) returned
`missing`/`ok:true` and left the link in place — including under `--force`. Both now use
`pathEntryExists`. I verified the two cycle-11 regressions FAIL against `3553f47`'s
installer and PASS at `dcd200e`:

- `health does not call a dangling destination symlink missing`: at `3553f47`,
  `core.status` was `missing` (the assertion `notEqual(core.status, "missing")` failed); at
  `dcd200e`, it is `malformed` and `result.ok` is `false`.
- `a forced uninstall removes a dangling destination symlink`: at `3553f47`, the link
  survived `--force` (the `assert.throws` found no ENOENT); at `dcd200e`, the link is
  removed.

I also verified the broader edge: a dangling symlink on an *add-on* destination (not just
the core) is correctly seen by the uninstall fleet-wide blocker check — the add-on is
`blocked`, the core is `skipped`, and the core is not removed. `doctor` sees a dangling
add-on symlink as `malformed`, not `missing`. A `--force` install over a dangling symlink
correctly replaces it with a proper directory.

### All earlier findings — remain closed

The source-side B5 hole (symlink in a manifest-free add-on or the core payload), the
ownership predicate, the marker skill-identity check, the selection recording, the B6
runtime floor, and the manifest integrity mechanism are all unchanged by cycles 10 and 11
and remain closed under the full 320-test suite. The five manifest-free skills under the
universal installer remain the recorded B3.11 follow-up, not a 2.1.0 blocker.

## New findings

None.

The remaining `fs.existsSync` calls (lines 812, 825, 1051, 1314, 1662, 1834) are all on
files *inside* an already-located root — `rewriteStagedGeminiManifest`, `validatePayload`,
`copySourcesFor` optional entries, `validateInstalledPayload`, and `readJsonFile`. In
every case a dangling symlink and a truly absent file both mean "required file not usable"
and receive the same disposition (missing-file problem or skipped optional entry). The
cycle-11 comment on `pathEntryExists` (lines 1782-1787) correctly states this boundary
rather than claiming universality. I checked each call site and agree with the boundary.

## Release judgement

Releasable as 2.1.0. The codex-1 round-8 MAJOR is closed on both halves, the cycle-11 gap
is closed, and the ordinary success paths (install, doctor, uninstall, dry-run, --no-addons,
--only, --force, all-target) are undisturbed. The deferred B3.11 follow-up (manifests for
the five remaining skills) is the first item in the consensus and is not a 2.1.0 blocker.

## What I verified

- `npm test`: 320 node tests, 0 fail. Python leg refused 3.9.6 as below the declared
  `>=3.10` floor (correct by design); ran 54/54 across seven files on python3.10.
- Manifest check: `node scripts/build-addon-manifest.js --check` — 47 files, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- `git diff ebe269e..dcd200e`: reviewed the full cycle-10/11 diff — 148 insertions, 10
  deletions across `lib/installer.js` and `test/bidding-addon.test.js`.
- `git diff 49fc3ec..dcd200e`: reviewed the full change since the round-7 tree, including
  the `test/design-addons.test.js` parameterization and backslash-continuation grammar
  assertions.
- Cycle-10 regressions: extracted both tests and ran them against `ebe269e`'s installer
  (both FAIL — earlier targets written, dangling symlink not blocked) and `dcd200e`'s
  (both PASS).
- Cycle-11 regressions: extracted both tests and ran them against `3553f47`'s installer
  (both FAIL — `missing` status, surviving symlink) and `dcd200e`'s (both PASS).
- Ordinary success paths: clean install, doctor (healthy with python3.14 on PATH),
  uninstall, dry-run, --no-addons, --only, --force over unmanaged destination, all-target
  --include-undetected — all pass against dcd200e's installer.
- Edge cases: dangling symlink on an add-on (not just core) — fleet blocker check sees
  it, `doctor` reports `malformed`, `--force` uninstall removes it, `--force` install
  replaces it with a proper directory.
- `pathEntryExists` semantics: ENOENT → false (absent), any other error → true (present),
  success → true (present, including dangling symlink). Verified by direct `lstatSync`
  behaviour and code reading.
- Remaining `fs.existsSync` calls: audited all six — each operates on a file inside an
  already-located root where dangling and absent receive the same disposition.
- `readMarkerState` (line 1748) uses the same ENOENT-only-is-absent pattern independently
  and is consistent with `pathEntryExists`.
- IMPLEMENTATION.md: cycle-10, round-9 outage, and cycle-11 records are accurate.
  `head-commit` is `dcd200e` via `git rev-parse`, confirmed.
