---
agent: codex-1
idea: addon-manifest-coverage
review-round: 3
date: 2026-08-02
reviewed-commit: f61e66b
---

## Status of the five agreed fixes

1. **Lockfile version — CLOSED.** Executed a JSON assertion over both root version fields (`2.2.0`) and `npm install --package-lock-only --ignore-scripts --no-audit --no-fund` in a clean archive; it exited 0 without changing `package-lock.json`.
2. **Source damage must not revoke ownership — CLOSED as agreed.** The focused regression passed, and a direct probe returned `status: malformed`, `managed: true`, `missing: []`; an unforced core-only uninstall also succeeded.
3. **F4 regression must reach F4 at the base commit — CLOSED.** With the HEAD test copied into `23a9856`, the fixture completed and the test failed on its own behavioral assertion (`actual: valid`, `expected: malformed`), not on manifest-presence setup.
4. **All three foreign-copy staging shapes — CLOSED.** The focused regression executed Codex, Gemini, and Antigravity and passed all three as `valid-unmanaged`.
5. **Non-zero `safeSourceFiles` floor — CLOSED.** The focused empty-source regression passed; direct empty- and missing-source probes over a gutted destination both returned `malformed`, including the source problem even when `SKILL.md` satisfied the fallback file list.

## Findings on the cycle-2 code

### [MINOR] `managed` still disagrees with the mutation ownership predicate for installed-payload damage

The new branch cannot produce a false `managed: true` relative to `installerOwnsDestination`: with an absent marker plus an empty source, and separately with a foreign marker plus a damaged source, `doctor` returned `managed: false` and unforced uninstall was blocked. A present, readable marker with the expected package name and skill is required before the source-only branch can return true.

It can still disagree in the opposite direction. I installed the core, deleted only the installed `plugin.json`, and retained its valid marker. `doctor` returned `status: malformed`, `managed: false`, and `missing: ["plugin.json"]`; an unforced core-only uninstall then returned `ok: true`, proving the mutation path's `installerOwnsDestination` predicate considered the same tree owned. This is the pre-existing non-source arm retained by the cycle-2 ternary, but it conflicts with the new comments and the documented meaning of `managed` as the ownership fact. External automation can therefore be told that an ordinary damaged install is unowned even though this installer will replace or remove it without `--force`.

Use the same predicate for both answers: compute `managed` directly from `installerOwnsDestination(unit.dest, unit.skill)` (or a shared parsed-state equivalent), and keep payload health solely in `status`, `missing`, and `problems`. Add a regression for a marker-valid, payload-missing tree that asserts `managed: true` and confirms the unforced mutation path agrees.

The other two cycle-2 attacks produced no findings. `safeSourceFiles` always carries a non-empty source problem alongside its fallback, so neither an empty nor an unreadable source certified a defective tree; a marker-present core takes the stricter `corePayloadFiles` branch instead. All three `packageWithoutManifest` callers passed at HEAD, and removing the current package's `parley-worktrees` manifest made the F4 migration test fail `valid` versus `malformed`, so the relaxed setup cannot silently make that test pass for the old fixture-construction reason. The other two callers intentionally test a manifest-free package, for which the asserted post-condition is sufficient.

## Verdict

FINDINGS — listed above
