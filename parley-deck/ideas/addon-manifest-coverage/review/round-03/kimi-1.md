---
agent: kimi-1
idea: addon-manifest-coverage
review-round: 3
date: 2026-08-02
reviewed-commit: f61e66b
---

## Status of the five agreed fixes

All five verified by execution against temp exports of `f61e66b`, `e46f661` and `23a9856`
(`git archive`, no `node_modules` copied; a symlink served `commonmark`). Method: full node leg
at `f61e66b`; the two changed test files from `f61e66b` overlaid onto the older trees so each
claim is measured, not read.

1. **Lockfile version — CLOSED.** At `f61e66b`: `package.json` 2.2.0, `package-lock.json` root
   `version` and `packages[""].version` both 2.2.0. The `npm install --package-lock-only`
   dirty-tree consequence no longer applies.
2. **`managed` vs present, valid marker — CLOSED.** The new regression fails at `e46f661` on
   its discriminating assertion (`managed`: `false !== true`) and passes at `f61e66b`. My own
   adversarial probe (below) confirms the fix shape.
3. **F4 regression's base-commit failure mode — CLOSED.** Overlaid onto `23a9856`, the test now
   constructs its fixture and fails at `test/bidding-addon.test.js:246`
   (`status: 'valid' !== 'malformed'`) — its own F4 assertion about the migration branch, not
   the fixture's pre-condition. The overclaim is gone.
4. **Third staging shape (`agy`) — CLOSED.** The extended foreign-copy test passes at
   `f61e66b` (inside 384/384). It also passes with the new test file overlaid on `e46f661`
   (127/129 there, only the two new regressions failing), which is exactly what a
   coverage-only extension with no accompanying `lib` change should do.
5. **`safeSourceFiles` empty-source floor — CLOSED.** The new regression fails at `e46f661` on
   its discriminating assertion (problems contained only the no-marker text; nothing matching
   `/SKILL\.md|is empty/`) and passes at `f61e66b`.

Cross-checks at `f61e66b`: `node --test` 384/384; `build-addon-manifest.js --check` green on
all six; `parley-bidding` aggregate still `sha256:7854adf1…`.

## Findings on the cycle-2 code

None. The three attack surfaces, taken in turn:

**`managed = sourceOnlyFailure ? (ownershipOk && !unmanaged) : (payloadOk ? !unmanaged : false)`.**
The new branch can report `managed: true` only when `missing` is empty, `unmanaged` is false,
and `problems` minus the source-enumeration problems is empty. Every non-owned marker state —
absent, unreadable, foreign name, wrong or missing skill identity — pushes a problem that is
not in `sourceProblems`, so it survives that filter and forces `ownershipOk` false. No string
collision is possible: the source-problem texts (`packaged source …`) share no wording with
the marker-state, `markerProblem`, or `manifestProblems` texts, and `manifestProblems` only
exists for add-ons, where `sourceProblems` is always empty so the new branch never runs.
Executed: with the package damaged (plugin.json removed) and (a) the marker rewritten to
another installer or (b) the marker deleted, `managed` reports `false` both times; with the
marker intact it reports `true` while `status` stays `malformed`. It cannot report ownership
for a tree this installer does not own, and in every case I could construct it agrees with
`installerOwnsDestination` on the owned side. A residual divergence exists in the other
direction — a valid marker plus a genuinely damaged *installed* payload (missing file) still
reports `managed: false` while the mutation paths, reading the marker, would proceed — but
that is precisely the semantics the consensus ratified ("marker evidence and installed-payload
validity, excluding source-enumeration problems"), it predates cycle 2, and I measured it
(case: damaged source + installed `SKILL.md` deleted → `managed: false`). Not a finding.

**The `["SKILL.md"]` floor in `safeSourceFiles`.** Reachable two ways — `listPayloadFiles`
returns `[]` (empty or all-ignored source) or throws (unreadable) — and both branches push a
problem alongside the floor, so health goes red and the install preflight throws; the floor
only ever *adds* a requirement on top of a guaranteed failure. Its single call site is the
unmarked-core branch of `validateInstalledPayload` (marked trees take `corePayloadFiles`;
add-ons never reach it; install staging always passes `packageRoot`), where `managed` is false
regardless. The new regression covers the empty-directory case; the unreadable case was
covered in cycle 1. No path lets a defective tree pass through the floor.

**`packageWithoutManifest` without the pre-condition assertion.** All five uses
(`test/bidding-addon.test.js:206,222,239,363,1686`) build a package whose `parley-worktrees`
manifest genuinely exists at `f61e66b` and is genuinely removed; the post-condition assertion
guarantees the fixture state the tests exercise. What the dropped assertion guarded — that the
real package still ships the manifest — is independently enforced by the manifest-coverage
tests and by `--check` inside `npm test`. Measured at `23a9856`: the fixture now constructs
and the F4 test fails on its own assertion (fix 3 above). No test silently passes for the
wrong reason at HEAD.

## Verdict

NO FINDINGS — ready to complete
