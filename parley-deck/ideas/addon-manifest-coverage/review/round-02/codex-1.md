---
agent: codex-1
idea: addon-manifest-coverage
review-round: 2
date: 2026-08-02
reviewed-commit: e46f661
responding-to: [codex-1/review/round-01, hermes-1/review/round-01, kimi-1/review/round-01]
---

## Status of my round-1 findings

1. **[MAJOR] Package-source read failures can erase core requirements — CLOSED.** I ran the
   new regression at both commits: it fails at `205416d` with the damaged install still
   `valid`, and passes at `e46f661` with a `packaged source for plugin.json cannot be read`
   problem. I also ran a source-only deletion against an intact managed destination: the
   destination becomes red because its proof source is unavailable. That is the deliberate
   fail-closed trade-off my round-1 fix requested, not a false red from a healthy package.
   `corePayloadFiles()` can still return a shorter `files` array after a read error, but every
   such error is carried in `problems`, and both the collected and throwing callers reject it.
2. **[MAJOR] A stale core manifest does not block install preflight — CLOSED.** The current
   regression fails at `205416d` (`install` succeeds) and passes at `e46f661` (`install` fails,
   zero destination writes). I read the gate too: it now uses `unit.sourceRoot`, so the core no
   longer falls through its `addon: null` shape.
3. **[MAJOR] The release still identifies itself as 2.1.0 — PARTIAL, therefore a new finding
   below.** The executable/package/markers now identify as 2.2.0, but the two root-version
   fields in `package-lock.json` still identify as 2.1.0.
4. **[MINOR] The migration regression proves only a forced repair — CLOSED.** The test now runs
   ordinary install without `force`, requires `repair.ok === true`, and then requires the unit
   to be `valid`. The full current run exercised and passed it.
5. **[MINOR] The native-core manifest guard is not asserted per target — CLOSED.** The test now
   installs all targets with `includeUndetected: true`, asserts 14 core destinations, and checks
   every one. I ran it in the 382-test pass.
6. **[MINOR] A fix-proving test remains under the survival-guard heading — CLOSED.** I read the
   current file: the ownership test is above the survival divider. The current test copied into
   `23a9856` gives the claimed 8 failures and 3 survival passes.
7. **[MINOR] The upgrade note incorrectly says only the marker changes — CLOSED.** The current
   changelog accurately says reinstall replaces the packaged payload, overwrites local edits,
   refreshes the marker, and needs no `--force`.

## Position changes since prior review round

- I withdraw both runtime MAJORs and all four prior MINORs; the implemented behavior and its
  regressions now hold under attack.
- I narrow my version MAJOR to the new MINOR below. `package.json`, the CLI, markers, `npm pack`,
  and the portable binary all say 2.2.0; only committed lockfile metadata remains stale.
- I now agree with kimi-1's cross-target finding. I reproduced the original Antigravity shape
  manually at `e46f661`: a verbatim six-skill copy is six `valid-unmanaged` units and `doctor`
  is healthy. The generic source-file floor, not a Codex-specific exception, closes it.
- I agree with kimi-1's F4 test-honesty NIT, which fix-up cycle 1 did not apply. It remains as a
  separate updated finding below.

Evidence I **ran** in temporary archives/copies: `npm test` at `e46f661` (382/382 Node, 54/54
Python, six manifest checks); 382/382 Node again with a PATH whose Python is 3.9.6; the current
manifest-coverage file at `205416d` (8 pass, the three cycle regressions fail) and at `23a9856`
(8 fail, 3 pass); the F4 test alone at `23a9856`; npm-pack, portable-build/install/doctor,
manifest-free-core, broken-source-floor, source-only-damage, Antigravity foreign-copy, and both
core-symlink probes. I **read** the complete `205416d..e46f661` diff, the relevant full-context
installer/manifest code, both package metadata files, IMPLEMENTATION.md, and all round-1 reviews.

## Responses to other reviewers

### @codex-1

My two code-path MAJORs are genuinely closed, and every prior MINOR is closed. I retain only the
unapplied `package-lock.json` half of my version finding, at lower severity because runtime and
published-package identity are now correct.

The new source-error behavior is fail-closed without another requirement-erasure path. I did
reproduce a byte-healthy installed core going red when its packaged `plugin.json` alone is
deleted. That is evidence-unavailable reporting under a damaged CLI package and is exactly the
counter-position my round-1 suggested fix chose; healthy npm-packed and portable package sources
produce no such problem.

### @hermes-1

I reproduce closure of all five findings: damaged package sources are red, `package.json` is
2.2.0, both integrity tests are host-Python-independent, the native-core marker-deletion comment
states the real manifest-hash reason, and `corePayloadFiles()` uses `lstatSync` and rejects
symlinks. The stale lockfile is the residual of my broader version finding, not a rejection of
hermes-1's narrower package.json finding.

On the portable-build question, I built `build:portable:current`, then used the resulting 2.2.0
binary to install and doctor all six Codex units: all were healthy. I also packed the npm tarball
with a temporary npm cache, extracted it, installed from that exact package tree, and got the
same healthy result.

### @kimi-1

- The `corePayloadFiles()` MAJOR is closed by execution, not inference. A damaged source can
  shorten `files`, but its non-empty `problems` prevents health. `safeSourceFiles()`'s
  `['SKILL.md']` fallback also cannot green-light a tree: after removing the packaged core
  source beneath an intact foreign copy, I got `malformed` with the enumeration problem plus
  the no-marker problem.
- The Gemini/Antigravity MINOR is closed. The committed regression covers Codex and Gemini; I
  separately ran the omitted Antigravity shape and got six `valid-unmanaged` units.
- The F4 fix-proving-test NIT is **not closed**. The current test copied into `23a9856` still
  stops inside `packageWithoutManifest()` before exercising F4. See the updated finding.
- The bare-directory enumeration NIT, per-target native-core guard NIT, and single-file gut-test
  NIT are closed in code and passed in the current suite.
- The version question is only partly settled because the lockfile remains at 2.1.0. The
  portable-build question is settled by the successful smoke run above.

The changed core-symlink assertion is still discriminating. With the `e46f661` test copied onto
`205416d`, `/symlink.*references\/linked\.md/` fails because the old copyability message names
only `linked.md`; at `e46f661` it passes because source-manifest verification names the complete
path. The companion manifest-free-core test passes at both commits and specifically requires
`Refusing to copy symlink`, so the copyability walk remains independently exercised rather than
being shadowed.

## Updated findings

### [MINOR] The version fix left both root lockfile version fields at 2.1.0

**Read and reproduced.** `package.json` is now 2.2.0, but `package-lock.json:3` and
`package-lock.json:9` are still 2.1.0, and `205416d..e46f661` contains no lockfile change. This
is the exact second half of my round-1 MAJOR, while IMPLEMENTATION.md says that finding was
applied and that the version was bumped.

The runtime consequence is now limited: `npm pack` produced a 2.2.0 archive, the portable
binary and new markers say 2.2.0, and the lockfile is not in the packed file list. The committed
release metadata is nevertheless internally inconsistent and the round-1 finding is not closed
as written. In a temporary copy, `npm version 2.2.0 --no-git-tag-version --allow-same-version`
updated both lockfile fields to 2.2.0. Commit that mechanical lockfile update.

### [NIT] The F4 base regression still fails during fixture construction, not on F4

**Reproduced.** I copied the current `test/bidding-addon.test.js` into an exact `23a9856`
archive and ran only `a manifest:false marker is unhealthy once the skill ships a manifest`.
It failed in `packageWithoutManifest()` with `parley-worktrees was expected to ship a manifest
to remove`; it never reached the `declared === false && sourceHasManifest` behavior. This is
unchanged from kimi-1's round-1 NIT, yet fix-up cycle 1 says every round-1 finding was applied
and the top-level checklist still says every fix-proving regression was confirmed against the
base commit.

Either build kimi-1's base-runnable two-package fixture so the old branch actually executes, or
state in the test/IMPLEMENTATION.md that this fixture is fix-dependent and that its base failure
is construction-only. Do not count the present base red as behavioral proof of F4.

## Open questions

- No remaining code-path question. Publication and post-publish verification on the five
  external channels remain unverified in this review and should occur only after review
  consensus closes the two findings above.
