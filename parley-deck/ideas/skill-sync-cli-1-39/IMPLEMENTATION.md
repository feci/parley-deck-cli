---
idea: skill-sync-cli-1-39
implementer: claude-1
target-repo: parley-deck-skill
released-as: 2.4.0
date: 2026-08-06
status: ready-for-review
---

# IMPLEMENTATION — skill 2.4.0

Implements D1-D5 from `FINAL.md`. Target repository is **`parley-deck-skill`**, not
`parley-deck-cli`; the idea lives in the CLI deck because that is where the deck is.

## What landed

| decision | files |
|---|---|
| D1 `opencode` row | `skills/parley-deck/SKILL.md` (Autonomous Execution table) |
| D2 `:251` replacement | `skills/parley-deck/SKILL.md` |
| D3 manual/CLI branch split | `skills/parley-deck/SKILL.md` (Generic CLI Invocation Contract), `skills/parley-deck/references/WORKED_EXAMPLES.md` |
| D5 `writeModeArgs` removal + migration rule | `skills/parley-deck/SKILL.md` (Headless Agent Configuration + manual branch), `references/WORKED_EXAMPLES.md` |
| D4 version bump + guard | `package.json`, `package-lock.json`, `references/compatibility.json`, `scripts/build-addon-manifest.js`, `test/manifest-coverage.test.js` |
| release notes | `CHANGELOG.md` |
| regenerated | `skills/parley-deck/parley-addon.json` (payload hash follows the edited files) |

`skills/parley-deck/references/COOPERATION.md` is **untouched**, per D6.

## Deviations from FINAL

**One, in D4's implementation shape.** FINAL says "add exactly one equality assertion … to the
existing Node test harness" and separately requires closing the `prepack` half by extending
`prepack` **or** by putting the check inside `build-addon-manifest.js`. Read literally those two
could be satisfied by two independent copies of the comparison. Implemented instead as **one**
comparison with two callers:

- `versionSyncProblem()` in `scripts/build-addon-manifest.js`, called from `main()` — so `prepack`,
  which runs that script, fails on drift.
- `test/manifest-coverage.test.js` requires and asserts the same exported function — so `npm test`
  reports it by name.

`main()` is now guarded by `require.main === module` so the script can be required without
executing. This is a shape difference, not a scope change: no new script, job, or checker exists,
and both gates are closed.

## Verification

All measured, not asserted.

**The guard fires on the real defect (fix-proving).** Before the bump, with `skillVersion` at
`1.4.3` and the package at `2.3.0`:

```
$ node -e "…versionSyncProblem()"
compatibility.json skillVersion is "1.4.3" but package.json version is "2.3.0" — update …
```

**`prepack` gates on drift.** With `skillVersion` deliberately desynced to `2.3.0`:

```
$ npm run prepack
build-addon-manifest: compatibility.json skillVersion is "2.3.0" but package.json version is "2.4.0" …
```
non-zero exit. Restored afterwards; the desync also (correctly) failed the payload-integrity tests,
since `compatibility.json` is part of the hashed payload.

**Full suite green after the change:** `npm test` → **386 node tests pass, 0 fail**; 54 python
tests across 7 files; all six add-on manifests `ok`.

**Scope of the `writeModeArgs` removal was enumerated with `find`, not `grep`** — per the FINAL
follow-up 2, a negative result from the facilitator's `grep` is not evidence of absence. Only two
files in the repository ever mentioned it (`SKILL.md`, `references/WORKED_EXAMPLES.md`); both are
updated.

## Not done here

- No deck outside this repository was edited. The 9 exclusive decks are follow-up 1 in `FINAL.md`;
  D5's migration rule is the instruction, not an automated migration.
- No release performed yet — awaiting `hermes-1`'s revision-2 signoff (see `consensus.md`).
