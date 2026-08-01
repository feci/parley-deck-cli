---
agent: hermes-1
idea: addon-manifest-coverage
review-round: 1
date: 2026-08-01
reviewed-commit: 205416d
---

## Summary

The implementation correctly addresses all six binding fixes (F1–F6). I verified each fix-proving test fails at base 23a9856 (5 of 5) and each survival guard passes there (2 of 3; the third — "health does not confer ownership" — was correctly identified by the implementer as mis-filed in FINAL.md, and is genuinely fix-dependent). F4's migration branch fires correctly, produces the right message, and is repairable by re-running install — confirmed with both intact and gutted trees. The reworked existing tests preserve their original discriminations: the manifest:false test, the legacy-marker exemption test, and the design-addons doctrine-budget test all still test what they tested before, using fixture packages instead of borrowed repository state. The `--check` passes on all six manifests. The source-anchoring (F3) works: a foreign copy of all six skills reports six `valid-unmanaged`, `managed: false`.

Two issues block the definition of done, and a test portability bug causes 2 of 8 new tests to fail on this machine.

## Findings

### [MAJOR] corePayloadFiles can return an empty required list, re-opening the F5 false green

`corePayloadFiles(packageRoot, kind)` walks each `PAYLOAD_ENTRIES` source via `fs.statSync(abs)`. If `statSync` throws — nonexistent path, unreadable file (chmod 000), dangling symlink — the walk returns early and that entry contributes nothing to the required list. When the entire packageRoot is nonexistent, every entry throws and the function returns `[]`, making `required = []`.

I confirmed by execution: a native codex install, core gutted to only SKILL.md + marker, doctor run with `packageRoot` pointed at a nonexistent path → reports `valid`, `managed: true`, `missing: []`, `problems: []`. I also confirmed a partial case: package with `skills/parley-deck/agents/` replaced by a dangling symlink, `agents/openai.yaml` deleted from the installed tree → reports `valid`, `missing: []`. This is the exact false green F5 was designed to close, re-entering through the back door.

In normal CLI usage `PACKAGE_ROOT = path.resolve(__dirname, "..")` is always valid, so this is a latent defect rather than a live bug in the default path. But the implementer flagged it in IMPLEMENTATION.md notes ("If the package root is unreadable it silently contributes nothing to the required list — check whether that can turn into a false green"), and the answer is: yes, it can.

Suggested fix: if `corePayloadFiles` returns an empty array when `packageRoot` is truthy, treat it as an error — either throw (matching `validatePayload`'s behavior for a broken package) or emit a problem. Alternatively, compute the required list once at install time and record it in the marker so validation never depends on the package root being readable at doctor time.

### [MAJOR] Version not bumped

FINAL.md's definition of done requires "version bumped." `package.json` still says `2.1.0`. The CHANGELOG has a `## 2.2.0 — 2026-08-01` entry. The diff `23a9856..205416d -- package.json` is empty — the file was not touched.

Suggested fix: bump `package.json` version to `2.2.0`.

### [MINOR] Two new tests fail on hosts with python3 < 3.10

The tests "a verbatim foreign copy of every packaged skill is valid-unmanaged" and "health does not confer ownership" both assert `result.ok === true`. But `doctor.ok` folds in runtime availability (`doctorCommand` at line 381: `(!skill.runtime || skill.runtime.ok)`), and `parley-bidding` declares `runtime.python: ">=3.10"`. The test comment (line 51–53) anticipates this and uses `process.env.PATH` to find python3, but on this machine that PATH finds python 3.9.6, so `parley-bidding.runtime.ok` is false and `doctor.ok` is false.

I measured: `node --test test/manifest-coverage.test.js` → 6 pass, 2 fail. Full suite: 376/378. The status assertions (all six `valid-unmanaged`) are correct and pass; only the `result.ok` assertion fails. The IMPLEMENTATION.md claim of "npm test 378/378" was not reproducible on this machine.

Suggested fix: assert on status fields only (not `result.ok`) in these two tests, or filter `parley-bidding` out of the `ok` check, or resolve a python3 >= 3.10 in the test environment. The other 6 tests in the file are unaffected because they don't assert `result.ok === true`.

### [MINOR] Stale comment in bidding-addon.test.js line 856

The test "the core unit never becomes valid-unmanaged" says the core "ships no manifest and there is nothing to verify against." The core now ships a manifest (`skills/parley-deck/parley-addon.json` was added by F1). The test still passes — I confirmed the actual reason is that the installed core doesn't carry the manifest (it's not in `PAYLOAD_ENTRIES`), so `unmanagedButVerified` fails on the manifest hash mismatch (`manifestFileHash(source) !== manifestFileHash(dest)` because `dest` has no manifest file). The comment should reflect this.

Suggested fix: update the comment to say the installed core doesn't carry the manifest (the copy plan doesn't include it), so `unmanagedButVerified` fails on hash mismatch.

### [NIT] corePayloadFiles uses statSync (follows symlinks) with no recursion guard

`corePayloadFiles` uses `fs.statSync` which follows symlinks, and recurses into directories. A symlink loop in the package would cause infinite recursion. Not a realistic concern for a well-formed npm package, but `copyRecursive` (line 1975) explicitly checks `lstatSync` and refuses symlinks — `corePayloadFiles` doesn't share that guard.

Suggested fix: use `lstatSync` or add a visited-set, or accept the asymmetry with `copyRecursive` as a non-issue since the package is trusted.

## Open questions

1. Is the `corePayloadFiles` empty-list degradation reachable in any real CLI usage path? `PACKAGE_ROOT` is `path.resolve(__dirname, "..")` which should always exist. But portable builds (`process.pkg`) or unusual install locations could potentially make entries unreadable. Should this be guarded now or recorded as a follow-up?

2. Was the 378/378 claim measured on a machine with python3 >= 3.10? If so, the test portability issue should still be addressed for contributors on older systems.

3. The `corePayloadFiles` function is called at validation time (doctor) rather than reading a constant or a value recorded at install time. Was computing-at-validate-time chosen deliberately over recording-in-marker? Recording in the marker would eliminate the package-root readability dependency entirely, but would also mean the required list doesn't update if the package is upgraded without re-installing — which is exactly the F4 scenario.
