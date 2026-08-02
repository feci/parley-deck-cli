---
agent: hermes-1
idea: addon-manifest-coverage
review-round: 3
date: 2026-08-02
reviewed-commit: f61e66b
---

## Status of the five agreed fixes

1. **package-lock.json version** — CLOSED. Both root `version` fields read `2.2.0` at f61e66b. Ran `npm install --package-lock-only` against the clean tree: exit 0, `git status --porcelain` empty afterward. The lockfile no longer dirties the release tree on first install.

2. **managed vs marker contradiction** — CLOSED. The new formula `sourceOnlyFailure ? (ownershipOk && !unmanaged) : (payloadOk ? !unmanaged : false)` was exercised against six scenarios by execution (wrong package name, wrong skill identity, unreadable marker, no marker, source-damaged with valid marker, foreign verbatim copy). In every case `managed` agrees with `installerOwnsDestination`. The source-damaged scenario reports `managed: true, status: malformed` — ownership retained, health failed. `uninstall --only parley-deck` succeeds on that tree (confirming kimi-1's blast-radius measurement).

3. **F4 fixture discrimination** — CLOSED. `packageWithoutManifest` now removes-if-present and asserts the post-condition. Ran the F4 test at 23a9856 with the fixed fixture: it runs, and fails on its own F4 assertion (`actual: 'valid', expected: 'malformed'`) rather than dying during construction. At f61e66b it passes.

4. **Foreign-copy antigravity shape** — CLOSED. The test at `manifest-coverage.test.js:327` now iterates `[codex, gemini, agy]` with `.gemini/config/plugins` as the antigravity path. All three pass at f61e66b.

5. **safeSourceFiles empty floor** — CLOSED. An empty packaged skill directory now floors at `["SKILL.md"]` and pushes a source problem. Verified by execution: empty source + empty dest reports `malformed` with `missing: ["SKILL.md"]`; empty source + dest-has-SKILL.md reports `malformed` with the source problem; unreadable source (chmod 000) reports `malformed` with an EACCES problem. The floor is never reachable without its accompanying source problem.

## Findings on the cycle-2 code

No findings.

All three attack questions were pursued to a negative result:

**Q1 — Can `managed` report true for an unowned tree, or disagree with `installerOwnsDestination`?** No. `ownershipOk` is computed from `problems.filter(p => !sourceProblems.includes(p)).length === 0` plus `missing.length === 0`. The marker-defect problems (wrong name, wrong skill, unreadable, absent) are pushed into the local `problems` array and are not in `sourceProblems`, so they survive the filter and force `ownershipOk = false`. The only way `ownershipOk` is true is when no marker-defect problem was pushed — which requires present, readable, correct-name, correct-skill — exactly the predicate `installerOwnsDestination` evaluates. `managed` cannot be true when `installerOwnsDestination` is false. The reverse asymmetry (managed=false while owns=true) exists for manifest-problem and marker-selection-problem cases, but that is conservative: it withholds the ownership report from a tree with a metadata integrity issue, while the mutation paths still operate because they read the marker from disk. Reproduced all six scenarios; no disagreement in the dangerous direction.

**Q2 — Is the `["SKILL.md"]` floor reachable in a way that lets a defective tree pass?** No. The floor is returned only alongside a non-empty `sourceProblems` entry (either "is empty" or "cannot be enumerated"). That entry is folded into `problems` in `validateInstalledPayload`, so `problems.length > 0` and `payloadOk = false` regardless of whether `missing` is empty. The floor prevents `missing: []` from making an empty source look like nothing is required, but it cannot make the tree valid because the source problem always accompanies it. Reproduced for empty-directory, has-SKILL.md, and unreadable-source cases — all report `malformed`.

**Q3 — Does any test using `packageWithoutManifest` silently pass for the wrong reason at HEAD?** No. At f61e66b every skill ships a manifest, so the fixture finds it, removes it, and asserts the post-condition (`existsSync(manifest) === false`). If the removal failed (wrong path, etc.), the install would record a manifest aggregate rather than `manifest: false`, and the first assertion in each test (`marker.manifest === false` at line 210, line 242) would fail — not silently pass. The fixture's correctness is guarded by the tests that use it. All four call sites pass at HEAD, and the F4 regression discriminates correctly at 23a9856 (fails on the F4 assertion, not the fixture).

## Verdict

NO FINDINGS — ready to complete
