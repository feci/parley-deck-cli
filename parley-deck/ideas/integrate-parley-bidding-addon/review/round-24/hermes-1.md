---
idea: integrate-parley-bidding-addon
review-round: 24
agent: hermes-1
date: 2026-07-31
reviewed-commit: e274eb8
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All findings from prior rounds are closed. The two cycle-26 errors that codex-1 blocked on in round 23 have been repaired in cycle 27 and verified independently below.

1. **Asymmetric injection (codex-1 round 23, MAJOR):** `impl === path.win32 || process.platform === "win32"` allowed win32 semantics to be injected on POSIX but let the host platform override POSIX semantics on Windows — the exact arm the cycle-26 regression had added. Fixed: separators now come from `impl.sep === "\\"` alone (lib/installer.js:1414). The test at test/bidding-addon.test.js:1858-1867 redefines `process.platform` to `"win32"` and asserts both directions: `path.posix` injection still yields `["we\\ird","x"]` (backslash is a filename byte), and `path.win32` injection still yields `["t","x"]`. I reproduced the regression at 7e8ccec (posix parts become `["we","ird","x"]` under a redefined win32 host) and confirmed the fix at e274eb8 (posix parts stay `["we\\ird","x"]`). Closed.

2. **Duplicated tests and dead absolute-target test (codex-1 round 23, MAJOR):** The cycle-26 edit left the old block in place — three tests duplicated, and the dead "an absolute raw target does not replay its root as a component" test still present. All removed. I confirmed: `grep -n 'an absolute raw target does not replay' test/bidding-addon.test.js` returns no matches; `grep '^test(' | awk | sort | uniq -d` returns no duplicate titles; test count is 109 definitions in bidding-addon.test.js, 368 in the full suite. Closed.

No other outstanding findings. The settled items (concurrent-installer isolation as a recorded follow-up; gate narrowing refused in round 22) are not re-litigated.

## Position on the gate

1 (correct). The gate correctly refuses all seven accumulated collision arms and admits the payload. No round has found a defect in the payload itself across twenty-three rounds, and this round is no exception.

## New findings

None.

## Release judgement

Releasable as 2.1.0. The payload has not changed since 714712f (zero diff in `skills/parley-bidding/`), the manifest is intact (47 files, aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`), and the cycle-27 fix is the last repair needed. Nothing remains that must change before release.

## What I verified

**Node test suite:** `npm test` — 368 tests, 0 fail, 0 skip. (npm run test, full suite)

**Python test leg:** `PATH="/opt/homebrew/bin:$PATH" PYTHONDONTWRITEBYTECODE=1 node scripts/run-python-tests.js` — 54/54 across 7 files under python3 3.14.6. The leg correctly refuses 3.9.6 by design (the add-on declares `>=3.10`).

**Cycle-27 arithmetic fix (the core verification):**
- At 7e8ccec (old tree, copied via `git archive` to /tmp): redefining `process.platform` to `"win32"` and calling `rawTargetArithmetic("/a/link", "we\\ird/x", path.posix, ...)` yields `["we","ird","x"]` — the bug. Backslash is split as a separator on POSIX because the host platform overrode the injection.
- At e274eb8 (current tree): same call yields `["we\\ird","x"]` — correct. Backslash stays a filename byte because semantics come from `impl.sep` alone. The win32 direction (`rawTargetArithmetic("C:\\a\\link", "C:\\t\\x", path.win32)`) yields `["t","x"]` — also correct. Both directions asserted in the test with the host platform redefined.
- Temp copy cleaned up (`rm -rf /tmp/pd-r24-old`).

**Test hygiene:** 109 test definitions in bidding-addon.test.js, no duplicate titles. The dead "an absolute raw target does not replay its root as a component" test is gone. The three earlier-round tests ("a link reached through an earlier link", "a backslash inside a POSIX link target", "an intermediate link in a raw target") remain — they are legitimate, not duplicates.

**Manifest integrity:** `verifyPayload('skills/parley-bidding')` returns `ok: true, problems: []`. 47 files, aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since 714712f (`git diff 714712f..e274eb8 -- skills/parley-bidding/` is empty).

**Payload stability:** `git diff 714712f..e274eb8 --stat -- skills/parley-bidding/` — no output. The payload has not been touched in 13 fix-up cycles.

**No .pyc contamination:** `find skills/parley-bidding -name '*.pyc' -o -name '__pycache__'` — no matches. `PYTHONDONTWRITEBYTECODE=1` and `-B` are enforced in the test runner.

**Working tree:** `git status --short` — clean. No mutations were made during this review.
