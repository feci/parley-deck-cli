---
agent: codex-1
idea: addon-manifest-coverage
review-round: 1
date: 2026-08-01
reviewed-commit: 205416d
---

## Summary

The main foreign-install fix works, and the current full suite is green, but the implementation is not ready to release. I found three MAJOR issues: two core-specific fail-open paths contradict binding verification items, and the release metadata was not bumped. Four MINOR test/documentation defects leave ratified guarantees weaker or less accurate than stated.

## Refutation attempts

- Ran `npm test` from a temporary archive of `205416d`: 378/378 Node tests, 54/54 Python tests on the available Python 3.14, and all six manifest checks passed.
- Ran the new `test/manifest-coverage.test.js` in an exact `git archive 23a9856` copy. Five fix-dependent tests failed and the three genuine survival guards passed. The failures were the foreign fleet, mandatory coverage, marker-retained gutting, core copy-plan coverage, and healthy-unmanaged fleet; the passing guards were unmarked gutting, mutation refusal/no marker synthesis, and no native core manifest.
- Reproduced a foreign verbatim six-skill Codex copy as six `valid-unmanaged` units with `doctor` healthy. Reproduced marker-retained gutted add-ons as `malformed` and unmarked gutted copies as `malformed` with the manifest both kept and deleted.
- Reproduced F4 using a real marker-schema-2 `manifest: false` install. Current `doctor` reported the promised re-install message; an ordinary install without `--force` replaced the unit and returned it to `valid`.
- Attacked deviation 1 directly. Putting this installer's genuine marker onto the foreign core shape reported `malformed` with `plugin.json` and `gemini-extension.json` missing. Removing the marker from the native shape reported the explicitly deferred false red. An intact package root reached through a directory symlink installed and validated successfully.
- Reproduced F5 with an intact package source: deleting `plugin.json`, `agents/openai.yaml`, or `references/WORKED_EXAMPLES.md` from a native core reported each file missing. The first finding below records the failing partial-source variant.
- Installed all 14 targets with `includeUndetected: true` and confirmed that none of their native core destinations carried `parley-addon.json`.

## Findings

### [MAJOR] Package-source read failures can erase core requirements

`corePayloadFiles()` silently returns when `fs.statSync()` fails, while `listVisibleEntries()` silently converts a failed directory read to an empty list. The required-file set therefore shrinks according to whatever portion of the package source remains readable at health-check time.

Reproduced in a temporary `205416d` archive: after a native Codex install, deleting the installed core's `plugin.json` correctly produced `malformed` with `missing: ["plugin.json"]`. Deleting the package source's own `plugin.json` as well made the unchanged damaged install report `doctor.ok: true`, `status: valid`, `managed: true`, with no missing files or problems. A dangling source symlink or unreadable directory reaches the same omission branches.

This defeats F5 and the no-false-green invariant precisely when the packaged source is partial or unreadable. Make source enumeration fail closed: return source-enumeration problems alongside the file list (using `lstat`/explicit read errors), or consume a generated immutable core file list. A missing, unreadable, or symlinked copy-plan source must add a health problem, never remove destination requirements. Add the reproduced two-deletion case and an unreadable-directory case as regressions.

### [MAJOR] A stale core manifest does not block install preflight

The source-integrity preflight still runs only under `if (unit.addon && ...)`. The core now has `sourceRoot` and a manifest, but `unit.addon` remains null, so its manifest is ignored before writes.

Reproduced in a temporary `205416d` archive by appending bytes to `skills/parley-deck/SKILL.md` without regenerating its manifest. `build-addon-manifest.js --check` exited 1 with `parley-deck: STALE`, yet `installCommand` returned `ok: true`, wrote all six units, and copied the drifted core bytes. This contradicts binding verification item 5, which requires source drift to make `--check` fail and install preflight write nothing. No new regression exercises source drift in a newly manifested unit; the existing corrupt-source test covers only `parley-bidding`, which already had a manifest at the base commit.

Use `unit.sourceRoot` rather than `unit.addon.root` for source-manifest preflight on every unit. Keep staged verification shape-aware—the native core staging tree is intentionally different—but reject an invalid packaged core source before staging anything. Add a core-drift regression that fails at `23a9856` and asserts zero destination writes.

### [MAJOR] The release still identifies itself as 2.1.0

`CHANGELOG.md` declares a 2.2.0 release, but both `package.json` and the root package entry in `package-lock.json` remain 2.1.0. Reproduced by the full suite banner (`parley-deck-skill@2.1.0`) and by newly written markers, which still record version/source 2.1.0.

This leaves the FINAL definition of done unmet and prevents publishing the new release under the already-used npm version while also making status/marker evidence misidentify the implementation. Bump `package.json` and `package-lock.json` to 2.2.0, rebuild any version-bearing channel artifacts, rerun the release checks, and publish/verify the five required channels before marking the implementation complete.

### [MINOR] The migration regression proves only a forced repair

The F4 regression repairs the `manifest: false` fixture with `force: true`, although the emitted error and changelog tell users simply to re-run their normal install. I separately reproduced that an unforced reinstall works today.

The current test would remain green if a future change accidentally made the documented repair require `--force`, weakening the migration guarantee the test is meant to protect. Remove `force: true`, assert that the ordinary reinstall replaces the unit, and then assert the repaired marker contains the manifest digest.

### [MINOR] The native-core manifest guard is not asserted per target

FINAL and Amendment 1 require the native core's lack of `parley-addon.json` to be asserted per target. The committed regression checks only Codex and Gemini. I manually checked all 14 targets and the current behavior is correct, but target-specific paths—especially Antigravity's extra staged copy—remain outside the automated guard.

Install `target: all` with `includeUndetected: true` in the regression and assert the manifest is absent from every returned core destination (or iterate the complete resolved target table explicitly).

### [MINOR] A fix-proving test remains under the survival-guard heading

`health does not confer ownership: a healthy fleet is still unowned` is placed below `SURVIVAL GUARDS — each of these must pass before AND after`, while its own comment correctly calls it fix-proving and the exact base run fails it. That contradicts the load-bearing split required by Amendment 1 and the implementation note claiming the cases were separated.

Move the test above the survival-guard divider (or add a distinct fix-proving ownership subsection) so the file's structure agrees with its base-commit behavior.

### [MINOR] The upgrade note incorrectly says only the marker changes

`CHANGELOG.md` says re-running install causes “no payload changes; only the marker is refreshed.” Reproduced F4 repair replaces the installed unit and adds the newly shipped `parley-addon.json`; the manifest is part of the copied installed tree even though the pre-existing substantive files are unchanged.

This gives operators an inaccurate description of the approved migration. State instead that existing skill content is unchanged, while the new manifest is installed and the marker/tree is transactionally refreshed.

## Open questions

- The local environment exposed Python 3.14 only. Where is the required green evidence for the two supported interpreter legs named by FINAL.md?
- Are publication to the five release channels and their post-publish verification intentionally deferred until review consensus? If so, `IMPLEMENTATION.md` should keep them explicitly pending rather than treating the definition of done as complete.
