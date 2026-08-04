---
agent: kimi-1
idea: addon-manifest-coverage
review-round: 5
date: 2026-08-04
reviewed-commit: e4ee4d2
---

## Leak measurement (isolated TMPDIR)

Method, so the next reviewer can repeat it exactly:

- Export: `git archive e4ee4d2 | tar -x -C /tmp/kimi-r5-export-*` (repo untouched; no
  comparison checkout). `node_modules` was **symlinked**, not copied — `test/design-addons.test.js`
  requires `commonmark`.
- `TMPDIR=/tmp/kimi-r5-tmpdir-*` (created empty), ran full `npm test` from the export.
- Counted **all** entries under the TMPDIR before and after, not just `parley-*` globs —
  deliberately avoiding the broken-glob class of mistake that produced cycle 3's false `0`.

Result of run 1 (normal PATH, python 3.14): **385/385**, python leg 54/54, `--check` green on
all six manifests. Leftovers under the isolated TMPDIR: **0 test-created entries**. The only
content is `node-compile-cache/`, which Node v26 itself creates under `os.tmpdir()` — not the
suite, and present on any node invocation.

Run 2 (`env -i`, PATH where the only `python3` is `/usr/bin/python3` = 3.9.6): **385/385**
again, python leg skipped with `python3 is 3.9, but the add-on declares >=3.10`, `npm test`
exit 0, still **0 leftovers**.

The 18 that codex-1 and I measured in round 4 is now 0. **Confirmed, twice.**

Coverage cross-check: six `*.test.js` files call `mkdtempSync` at `e4ee4d2`. Five now have the
exit-handler cleanup; the sixth, `skills/parley-design-check/test/checker.test.js`, cleans up
per-test via `t.after()`. That is consistent with the measured 0 — no uncovered call site.

## Findings

None. Detail per assigned check:

**1. Temp leak — closed.** See above. Reproduced by the controlled method, twice, including
the 3.9.6-PATH variant. Both resistant-debris tests named in the cycle-4 notes
(`a frozen owned destination completes the install and names the debris`, `one unreadable
subdirectory deep in a destination no longer blocks anything`) exist at `e4ee4d2` and pass.

**2. `forceRemove` cannot escape the tracked tree, and cannot mask a failure.** I tested a
verbatim copy of the helper empirically, not just by reading:

- Fixture: tracked temp dir containing a 0555 `.bak` debris tree, a 0000 subdirectory, a
  symlink to an outside file, a symlink to an outside directory, and a self-referential
  symlink loop; outside sentinels at 0644/0755 with known contents.
- After `forceRemove`: tracked tree fully removed; outside file and dir **present, contents
  intact, modes unchanged**; loop not followed. `lstat` + symlink skip + recursion only
  through real directories holds; `rmSync` unlinks symlinks rather than following them.
- Exit code: spawned child processes whose exit handlers run the same pattern exited with
  their original codes (0 and 3). The handler wraps everything in try/catch, and `forceRemove`
  swallows its own `lstat`/`chmod`/`readdir` errors, so nothing can throw out of an `exit`
  listener.
- Masking: the handler runs at process exit, after the test file's results are already
  recorded, so it cannot intercept a failure — and the 0700 normalization happens after all
  assertions, so it cannot soften a test's own permission expectations either.
- `target` provenance is only `trackTemp(mkdtempSync(os.tmpdir()))`. The one theoretical
  residue is an lstat→chmod TOCTOU inside a private TMPDIR in a test exit handler; the threat
  model for that is nil and I am not raising it, even as a NIT.

**3. Editing the two shipped `skills/parley-tracker/bin/*.test.js` files — acceptable, and the
regenerated manifest is correct.**

- The edits add cleanup only (comment, `TEMP_DIRS` registry, wrapped `mkdtempSync`). The
  shipped runtime files `claim.js` / `validate.js` are untouched — their hashes are unchanged
  in the manifest diff.
- I recomputed all 8 payload file hashes from disk and the aggregate independently, using the
  generator's documented algorithm (`<path>\n<hash>\n` per file, path byte order) — result
  `sha256:07d9826373e4be3d2f393a8a56da616945fbc2e4e0838938827ceba7e85dfdd5`, exactly the value
  in `parley-addon.json` and in the cycle-4 notes. `--check` is green on all six manifests.
- Nothing else in the repo asserts the old manifest: `git grep` at `e4ee4d2` for the old
  aggregate (`2b178723…`) and both old test-file hashes (`a87f15f9…`, `102dd95b…`) returns
  zero references. The new aggregate appears only in the manifest itself; no test hard-codes
  it. Staleness going forward is guarded by `--check` in `npm test` and `prepack`.

**4. Cycle-4 notes — everything reproduced.** 385/385 on both PATHs; 0 leftovers (was 18);
`--check` green; tracker aggregate matches the notes; `parley-bidding` aggregate still
`sha256:7854adf15071…` (observed in the `--check` output); the two named debris tests exist
and pass. Nothing in the notes failed to reproduce. Given cycle 3's false `0`, I want to be
explicit about why this one is believable: the measurement is mine, from an isolated TMPDIR,
counting every entry, and it survived two full-suite runs plus my own forceRemove experiments.

## Ready to release?

Yes. The round-4 findings are all resolved with reproduction, the new code introduced to
resolve them (`forceRemove`, tracker test cleanup, manifest regeneration) checks out under
direct attack, and nothing else in the tree contradicts the new state. Nothing stops me from
shipping 2.2.0 on all five channels.

## Verdict

NO FINDINGS — ready to complete and release
