---
agent: hermes-1
idea: addon-manifest-coverage
review-round: 5
date: 2026-08-04
reviewed-commit: e4ee4d2
---

## Leak measurement (isolated TMPDIR)

Method: extracted e4ee4d2 via `git archive` into a temp dir, ran `npm ci`, then ran the full
suite (`node --test && node scripts/run-python-tests.js && node scripts/build-addon-manifest.js
--check`) with `TMPDIR` pointed at an isolated mktemp directory. Counted leftovers with `find`
(not a zsh glob pipeline — the bug that made the cycle-3 claim false).

Three runs, all isolated TMPDIR:

Run 1 — system python3 (3.9.6, below the >=3.10 floor, python leg skipped by the runner):
  385/385 JS, 0 parley-* dirs after, 0 non-cache leftovers.

Run 2 — python3.10 on PATH (symlinked into a temp bin dir so `python3` resolves to 3.10):
  385/385 JS, 54/54 python, --check green on all six, 0 parley-* dirs, 0 non-cache leftovers.

Run 3 — per-file isolation (installer, bidding-addon, manifest-coverage, tracker claim+validate
  each run individually with their own isolated TMPDIR): every file 0 leftovers.

The only thing remaining in the isolated TMPDIR after a full run is Node.js's own
`node-compile-cache/` directory — a runtime artifact, not a test fixture. It is created by
Node itself, not by any test code, and is harmless.

The 18-directory leak from round 4 is closed. I confirm 0.

## Findings

No findings.

### Check 1 — temp leak

Closed. Measured three times with an isolated TMPDIR (the method the round-4 reviewers
specified and the cycle-3 measurement did not use). 0 leftover test directories in every run.
The measurement command I used was `find "$ISOLATED_TMP" -maxdepth 4 -type d -name 'parley-*'`
plus a broader `find ... -not -path '*/node-compile-cache*'` sweep — no zsh glob pipeline.

### Check 2 — forceRemove

`forceRemove` is added to the three `test/` files (installer, bidding-addon, manifest-coverage).
It is NOT in the two shipped `skills/parley-tracker/bin/*.test.js` files — those use plain
`rmSync` because they never harden directories (no `chmod` calls, confirmed by grep).

Can it chmod or delete anything outside the tracked temp directory? No.
  - `forceRemove` is called only from `process.on("exit")`, iterating over `TEMP_DIRS`.
  - Every entry in `TEMP_DIRS` is the return of `trackTemp(fs.mkdtempSync(...))` — a freshly
    created temp directory.
  - `relax()` walks the tree using `readdirSync` (returns names, not resolved paths) and
    `path.join(abs, entry)` — it stays within the tree.
  - `lstatSync` does not follow symlinks. `stat.isSymbolicLink()` → return: symlinks are
    skipped entirely, never chmod-ed. I verified empirically: a symlink inside the temp tree
    pointing to an external directory was skipped, and the external directory's mode was
    unchanged (755, not 700).
  - `rmSync(target, { recursive: true, force: true })` with `target` always a tracked temp
    dir. Node's `rmSync` does not follow symlinks — it removes the link itself. Verified
    empirically: a symlink to an external directory was removed; the external target survived.

Can it mask a genuine test failure or change the exit code? No.
  - It runs in `process.on("exit")`, after all tests have completed and their pass/fail
    status is recorded.
  - Every call is wrapped in `try { forceRemove(dir); } catch (_error) {}` — any throw is
    swallowed. (I confirmed that an uncaught throw in an exit handler does set exit code 1,
    but the try/catch prevents that path.)
  - It cannot change a test's result or prevent an assertion from running.

### Check 3 — shipped payload test files and manifest regeneration

The two `skills/parley-tracker/bin/*.test.js` files were edited to add temp-dir tracking
(TEMP_DIRS array, trackTemp helper, process.on("exit") cleanup). The changes are cleanup-only
— no test logic, assertions, or expectations were altered.

The regenerated manifest is correct:
  - `node scripts/build-addon-manifest.js --check` passes on all six skills.
  - Regenerating `parley-tracker`'s manifest produces a byte-identical file (diff confirmed).
  - The new aggregate is `sha256:07d98263...`, matching IMPLEMENTATION.md.
  - The old aggregate `sha256:2b178723...` appears nowhere else in the repo (search confirmed).
  - `parley-bidding`'s aggregate is still `sha256:7854adf1...`, unchanged since 714712f, as
    claimed.

Is it acceptable? Yes. The parley-tracker skill shipped no manifest before this idea, so there
is no prior contract being broken. The test files are part of the shipped payload and their
hashes are correctly reflected in the manifest. The `--check` step in `npm test` and
`prepack` ensures a stale manifest fails the build.

### Check 4 — cycle-4 notes reproducibility

Every claim in the cycle-4 notes that I can verify, I verified:

  - "385/385 and 0 leftover directories" — reproduced (385/385, 0 leftovers, isolated TMPDIR).
  - "385/385 under a PATH whose only python3 is 3.9.6" — reproduced (run 1 above).
  - "--check green on all six manifests" — reproduced.
  - "parley-bidding aggregate still 7854adf1..." — reproduced.
  - "16 came from tracker tests" — the tracker tests now have cleanup; individually they
    leave 0 leftovers. The count of 16 is not something I can retroactively verify (the leak
    is fixed), but the fix is confirmed working.
  - "Two directories resisted removal" (hardened trees renamed aside) — the forceRemove fix
    addresses this; the installer and bidding-addon tests (which create hardened trees) leave
    0 leftovers.
  - The child-process `node -e` cleanup — confirmed working (bidding-addon tests, which
    spawn children, leave 0 leftovers).
  - "The first mechanical pass rewrote a mkdtempSync inside the child's script string" —
    the child script at bidding-addon.test.js line ~711 uses its own `fs.rmSync(h, ...)`
    rather than trackTemp, which is correct since the child cannot see the parent's helpers.
    All 13 runtime-probe tests pass.

Nothing in the cycle-4 notes is unreproducible.

## Ready to release?

Yes. The temp leak is closed (0 leftovers, measured three ways with an isolated TMPDIR).
`forceRemove` is confined to tracked temp directories, cannot escape via symlinks, and cannot
mask test failures or alter exit codes. The regenerated manifest is correct and consistent.
All 385 JS tests and 54 python tests pass. `--check` is green on all six manifests. Version
numbers are 2.2.0 across `package.json` and `package-lock.json`. I see nothing that stops
shipping.

## Verdict

NO FINDINGS — ready to complete and release
