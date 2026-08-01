---
idea: addon-manifest-coverage
status: fix-up-cycle-1
implementer: claude-1
started: 2026-08-01
completed: 2026-08-01
branch: parley-deck-skill#main
head-commit: e46f661
base-commit: 23a9856
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

All six binding fixes from `FINAL.md` are in. `lib/installer.js`, `scripts/build-addon-manifest.js`,
five new `parley-addon.json` files, one new test file, four existing tests reworked,
`README.md` and `CHANGELOG.md`.

`skills/parley-bidding/` is untouched — its aggregate is still
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

## Implementation plan / checklist

- [x] F1 — `parley-addon.json` for `parley-deck`, `parley-design`, `parley-design-check`,
      `parley-tracker`, `parley-worktrees`
- [x] F2 — generator coverage mandatory (`listAddons()` no longer excludes the core; no-argument
      default is every `skills/*` directory; a missing manifest is a `--check` failure)
- [x] F3 — `sourceRoot` on every unit; `unmanagedButVerified` reads it
- [x] F4 — `manifestProblems`: `declared === false && sourceHasManifest` → re-install to validate
- [x] F5 — core required-file list derived from `PAYLOAD_ENTRIES`
- [x] F6 — `README.md` and `CHANGELOG.md`
- [x] Checks run — at `205416d`: 378/378, python leg 54/54, `--check` green on all six.
      **That 378/378 was PATH-dependent and hermes-1 measured 376/378 on a stock macOS PATH;
      see fix-up cycle 1.** After cycle 1, at `e46f661`: 382/382 on both PATHs.
- [x] Every fix-proving regression confirmed failing at `23a9856`; every survival guard confirmed
      passing there

## Deviations from FINAL.md

**1. A collision FINAL.md did not anticipate: the core has two legitimate installed shapes.**

F5 derives the core's required-file list from the copy plan. The copy plan stages `plugin.json`
and `gemini-extension.json` in from the *repository root*, so a native core install contains
them. A **foreign** install does not — the universal `skills` CLI copies `skills/parley-deck/`,
and those two files are not inside it.

Applying the derived list to every core tree therefore replaced one false red with another:
measured, a byte-perfect foreign core reported `missing: ['gemini-extension.json',
'plugin.json']` and `malformed`, so F1+F3 appeared to have failed. F5 was written to fight F1.

Resolved by scoping the derived list to trees this installer owns. `skillUnitStatus` passes
`state.present ? unit.packageRoot : null`, so the copy plan governs only a tree with our marker;
for anything else the packaged manifest is the authority, via `unmanagedButVerified`, which
byte-verifies every declared file. This is a narrowing of F5, not a weakening: the case F5 was
written for — the natively installed core — is exactly the case that keeps the derived list.

**2. FINAL.md mis-classified one regression, and running it at the base commit caught it.**

Verification item 7 was filed as a survival guard ("ownership unchanged: a `valid-unmanaged`
fleet blocks install and uninstall without `--force`"). At `23a9856` it **failed**. Its subject
is a fleet that is simultaneously healthy and unowned, and before this change that state existed
for one unit in six, so the test cannot establish its own precondition on the base commit — it
is fix-dependent, not invariant.

Split into two: `health does not confer ownership` (fix-proving, fails at `23a9856`) and `an
unmarked tree refuses both mutations and is not stamped with ownership` (survival guard, passes
at `23a9856`). The error was the drafter's, in the same family as the overclaim kimi-1 corrected
in Amendment 1.1 — asserting a test's classification without running it against the commit it is
classified against.

**3. Four existing tests depended on a property of the repository, not of the code.**

`bidding-addon.test.js` used `parley-worktrees` as its "ships no manifest" subject. Every skill
now ships one, so that subject no longer exists. Rather than delete the coverage — the branches
are still live, and every pre-2.2.0 install on disk is that shape — a `packageWithoutManifest()`
fixture builds a package copy with one skill's manifest removed. The tests now construct the
state they are testing instead of borrowing it.

`design-addons.test.js`'s "exactly the four doctrine files" excludes `parley-addon.json` and
asserts its presence separately, so the doctrine budget still means what it says.

## Notes for reviewers

- **The two-shapes problem in deviation 1 is where I would attack this.** The scoping rule is
  "the copy plan governs only trees we installed". Is there a path that reaches
  `validateInstalledPayload` with a marker present but a non-native shape, or without a marker
  but a native shape? The second is deferred follow-up 3 (a native install whose marker was
  deleted reports `malformed` although its bytes are correct) and is known.
- `corePayloadFiles` walks the package at validation time rather than reading a constant. If the
  package root is unreadable it silently contributes nothing to the required list — check whether
  that can turn into a false green.
- F4 is deliberately loud and the user approved the consequence on 2026-08-01. The upgrade path
  is: `doctor` reports `malformed` with `re-run install to validate the payload`, and re-running
  install repairs it. Regression covers both halves.
- `parley-worktrees` is a single-file payload, so "gutted to `SKILL.md`" removes nothing from it.
  The gut test derives its subjects from payload size rather than listing them, so a skill that
  grows a second file joins the test automatically.
- `--check` is wired into `npm test` and `prepack`; a stale manifest fails the build rather than
  a user's `doctor`.

## Fix-up cycle 1
status: complete
completed: 2026-08-02
head-commit: e46f661
review-round: 1 (codex-1, hermes-1, kimi-1)

### Fixes applied

**[MAJOR — found independently by all three reviewers] `corePayloadFiles` failed open.**
It swallowed every read error, so the core's required-file list shrank to whatever of the
package was still readable. Reproduced: deleting `plugin.json` from both the installed core and
the packaged source made the damaged install report `valid`, `managed: true`, `missing: []`,
`problems: []`, `doctor` exit 0. My own F5 — written to close a false green — had opened one.
Now returns `{ files, problems }`, uses `lstatSync` (so a symlinked or looping source is a defect
rather than a traversal), and surfaces enumeration failures as health problems. Verified: the
same damaged pair now reports `malformed` naming `packaged source for plugin.json cannot be
read (ENOENT)`.

**[MAJOR — codex-1] Core source drift did not block install.** The source-manifest preflight
gated on `unit.addon`, which the core does not have. Reproduced: appending to
`skills/parley-deck/SKILL.md` made `--check` exit 1 while `install` returned ok, wrote all six
units and copied the drifted bytes. Gate now reads `unit.sourceRoot`. Verified: install exits 1,
zero destinations written, message `Source payload does not match parley-addon.json: modified:
SKILL.md`. **My verification of FINAL item 5 had exercised only the add-on path and I reported
the item as satisfied.**

**[MAJOR — codex-1, hermes-1] `package.json` still said 2.1.0** while the changelog announced
2.2.0. Bumped.

**[MINOR — kimi-1] The fix did not cover the gemini or antigravity targets.** An unmarked core
was held to the per-kind lists, which demand `gemini-extension.json` / `skills/SKILL.md` /
`plugin.json` — files that exist only because our staging puts them there and that a verbatim
foreign copy cannot contain. The codex list passed only because it happens to name files that
do live inside `skills/parley-deck/`. The unmarked floor is now the packaged source's own file
list, which is what a foreign copy actually is. Verified on codex and gemini: all six
`valid-unmanaged`. **My V1 verification had tested one target and FINAL item 1 named none.**

**[MINOR — hermes-1] Two tests asserted `doctor.ok`,** which folds in runtime availability, so
they also asserted which python the host had on PATH: they passed with 3.14 first and failed on
a stock macOS PATH with 3.9.6. hermes-1 measured 376/378 where I had reported 378/378 — my
number was true for my PATH only. Replaced with an integrity-verdict assertion. Verified:
382/382 under both PATHs, using the exact `node --test` invocation `npm test` runs.

**[MINOR — codex-1] The migration regression proved only a forced repair** while the changelog
tells users to run a plain `install`. Now asserts the unforced path.

**[MINOR — codex-1] The native-core manifest guard covered two targets.** Now derives every
target from the installer and asserts all fourteen core destinations.

**[MINOR — codex-1] A fix-proving test sat under the survival-guard heading.** Moved.

**[MINOR — codex-1] The upgrade note said only the marker changes.** Re-running install also
replaces the payload; local edits are overwritten. Stated, along with `--force` not being needed.

**[MINOR — hermes-1] A stale comment** claimed the core ships no manifest. Rewritten to state
the real reason that test still passes, and to name it as deferred follow-up 3.

**[NIT — kimi-1] `listAddons()` enumerated bare directories** where F2's binding text says
`skills/*/SKILL.md`. Now requires `SKILL.md`.

**[NIT — kimi-1] The gut test hard-coded `parley-worktrees`,** a single-file payload for which
"gutted to SKILL.md" removes nothing. Subjects are now derived from payload size.

**[NIT — hermes-1] `statSync` followed symlinks with no recursion guard.** Now `lstatSync`.

### A test whose meaning my fix silently changed

`a symlink in the CORE source is caught before the first write too` began failing — not because
the behaviour broke (still refused, still zero writes) but because source-manifest verification
now runs before the copyability walk and reports the symlink first. The walk that codex-1 added
in round 8, precisely because the core was excluded from it, was thereby shadowed for the core.
Rather than relax the assertion, the message check now asserts the outcome and the named path,
and a second test exercises the walk on a core source that ships no manifest, where nothing can
shadow it.

### Deviations from agreed fixes

None. Every round-1 finding was reproduced before being accepted, and every one was applied.

### Verification

- `npm test`: 382/382, and 382/382 again under a PATH whose only python3 is 3.9.6
- python leg 54/54; `--check` green on all six
- three new regressions fail at `205416d` and pass at `e46f661`
- at `23a9856`: 8 fail, 3 survival guards pass
