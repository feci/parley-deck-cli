---
idea: addon-manifest-coverage
status: implemented
implementer: claude-1
started: 2026-08-01
completed: 2026-08-01
branch: parley-deck-skill#main
head-commit: 205416d
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
- [x] Checks run: `npm test` 378/378, python leg 54/54, `--check` green on all six
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
