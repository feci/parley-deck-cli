---
idea: integrate-parley-bidding-addon
review-round: 23
agent: codex-1
date: 2026-07-31
reviewed-commit: 7e8ccec
---

## Verdict

BLOCK

## Outstanding findings — closed or not

The round-22 Windows-coverage finding is closed. Cycle 26 extracts and exports
`rawTargetArithmetic`, and the new regression calls that real helper with `path.win32`; its drive
and UNC assertions pass at `7e8ccec`. I overlaid the current test file on `git archive` trees for
both named older commits and ran the exact helper test. It exits 1 at `b1f43e4` and `2b7ca3e`
because the helper does not exist there, and passes at `7e8ccec`. The coverage seam I requested
now exists.

No prior collision-arm finding remains open. Two new cycle-26 findings below keep the strict gate
open.

## Position on the gate

**1 (correct).** The production gate still has the right shape: physical identity, physical
containment, and resolution-crossing checks before mutation. Cycle 26 does not expose another
escape from those arms. Narrowing it would still contradict the fleet-wide atomicity promise for
the deliberately supported symlinked-runtime-home layouts.

## New findings

### [MAJOR] The injected path semantics are overridden by the host platform

`lib/installer.js:1410` chooses Windows separators when either the injected implementation is
`path.win32` **or** the ambient host is Windows:

```js
const separators = impl === path.win32 || process.platform === "win32" ? /[\\/]+/ : /\/+/;
```

That makes the helper asymmetric: win32 injection works on POSIX, but the `path.posix` arm added
by this same regression fails on Windows. I exercised the exact branch on `7e8ccec` by setting
the process platform to the value Node supplies on Windows and calling the exported helper:

```text
node <<'NODE'
const assert = require('node:assert/strict');
const path = require('node:path');
const installer = require('./lib/installer');
Object.defineProperty(process, 'platform', { value: 'win32' });
const got = installer.rawTargetArithmetic('/a/link', 'we\\ird/x', path.posix, (d) => d);
console.log(JSON.stringify({ platform: process.platform, start: got.start, parts: got.parts }));
assert.deepEqual(got.parts, ['we\\ird', 'x']);
NODE
```

Observed: exit 1, with
`{"platform":"win32","start":"/a","parts":["we","ird","x"]}`. The expected parts were
`["we\\ird","x"]`.

This costs a real Windows user a red `npm test` on the supported Windows channel: the new test
unconditionally runs its `path.posix` assertion, but the ambient `process.platform === "win32"`
tears the backslash despite the injected implementation. It also means the helper's stated
injection contract depends on which host executes it. Make `impl` authoritative for separator
semantics (for example, select from `impl === path.win32` or `impl.sep`) instead of consulting
the ambient platform, then run this branch and the full suite again.

### [NIT] Cycle 26 appended three existing tests a second time and did not remove the dead test

The cycle diff is `+160/-0` in `test/bidding-addon.test.js`. This command on `7e8ccec`:

```text
rg -o '^test\("[^"]+' test/bidding-addon.test.js | sort | uniq -cd
```

reports these three names twice:

```text
2 test("a backslash inside a POSIX link target is a filename byte, not a separator
2 test("a link reached through an earlier link is walked from where it physically sits
2 test("an intermediate link in a raw target is expanded before `..` is applied
```

They are the blocks at lines 1697/1857, 1746/1906, and 1784/1944. The full test output executes
each pair and therefore rises from 368 tests to 372 even though cycle 26 adds only one real new
case; without the copies the count would be 369. The old absolute-target pin also remains at
line 1978, including unused `seen`/`record` values silenced with `void` at lines 1996-1997.
`IMPLEMENTATION.md` factually says that scaffolding and its test are gone.

The user cost is small but concrete: every `npm test` repeats three filesystem-heavy setups,
duplicates diagnostics, inflates the suite count, and leaves the canonical implementation
record false. Under this idea's strict gate, duplicate dead test code and a factual record error
remain findings. Delete the second copies, remove or replace the obsolete pin/scaffolding as
recorded, and correct the cycle-26 record.

## Release judgement

Not releasable as 2.1.0 on this strict gate. Correct cycle 26 so the injected implementation
fully determines path semantics and the test change contains only the intended new regression,
then obtain a fresh full-scope clean review.

## What I verified

- `git rev-parse HEAD` is `7e8ccec6d9049dd7f2bca798d99898d5dd594d4c`; both repositories
  were clean before the review file was written.
- Read the live protocol, `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through cycle 26, all
  round-22 reviews, the full `b1f43e4..7e8ccec` diff, and the six-file
  `49fc3ec..7e8ccec` full-scope diff. `git diff --check` passes.
- `npm test`: **372 Node tests, 372 pass, 0 fail**; Python **54/54** on 3.14; manifest **47
  files**, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- With `/usr/bin/python3` first on `PATH`, the Python runner exits 1 with
  `python3 is 3.9, but the add-on declares >=3.10`, as designed.
- Current helper regression: passes at `7e8ccec`; with the current test overlaid on `git archive`
  copies, fails at both `b1f43e4` and `2b7ca3e` with
  `TypeError: installer.rawTargetArithmetic is not a function`.
- Simulated the exact Windows host branch of the pure helper and observed the failing POSIX arm
  reported above. Independently counted the duplicate test declarations and saw both copies of
  all three execute in the full-suite output.
- `npm pack --dry-run --ignore-scripts --json`, using a disposable cache after the machine's
  default npm cache refused access: **202 files**, **48** under `skills/parley-bidding/`, no
  `.pyc`, `.pyo`, or `__pycache__` entry.
- `git diff 714712f..7e8ccec -- skills/parley-bidding/ bin` is empty. The payload remains
  unchanged; it contains no symlink or generated Python cache.
- Removed both archive trees and the disposable npm cache. No `parley-r23*` temporary directory
  remains.
