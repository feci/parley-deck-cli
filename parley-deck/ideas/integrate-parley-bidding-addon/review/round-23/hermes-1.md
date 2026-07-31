---
idea: integrate-parley-bidding-addon
review-round: 23
agent: hermes-1
date: 2026-07-31
reviewed-commit: 7e8ccec
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All findings from rounds 7 through 22 are closed. The one outstanding block
from round 22 — codex-1's MINOR for executable Windows drive/UNC coverage — is
closed by cycle 26. I verified each:

- **Windows drive/UNC coverage gap** (round 22, codex-1 BLOCK — "executable
  Windows coverage, not a defect in the gate's current shape"): CLOSED.
  `rawTargetArithmetic(from, rawTarget, impl, resolveParent)` is extracted
  from `walkRawTarget` and exported (lib/installer.js:1409-1419, 2627). The
  regression test "a raw link target's root arithmetic is platform-correct,
  drive and UNC alike" (test/bidding-addon.test.js:1818) injects `path.win32`
  and asserts four arms: a drive-absolute target does not replay `C:\` as a
  component; a UNC target does not duplicate server and share; a relative
  target anchors on the resolved parent with the resolver observed being
  consulted; a POSIX target keeps a backslash as an ordinary filename byte.

  I verified the test FAILS at both commits it names. Extracted each via
  `git archive` to a temporary directory, installed deps, injected the
  7e8ccec test file, and ran:
  `node --test --test-name-pattern="root arithmetic is platform-correct" test/bidding-addon.test.js`
  - At b1f43e4: 1 fail — `TypeError: installer.rawTargetArithmetic is not a
    function` (the logic was inline in `walkRawTarget`, not exported).
  - At 2b7ca3e: 1 fail — same `TypeError` (the logic was also inline, and
    used `path.resolve(rawTarget).root` with unconditional `/[\\/]+/` split,
    which on a POSIX host treats `C:\target\x` as relative and splits the
    drive letter into a component).

  At 7e8ccec the test passes. The discrimination is genuine: the function
  does not exist at either prior commit, and the inline logic at 2b7ca3e
  was materially wrong for Windows paths on a POSIX host (confirmed by
  tracing: `path.isAbsolute('C:\\target\\x')` → false on POSIX, so the
  drive letter became a path component rather than being stripped as root).

- **Linked-ancestor / physical-parent arm** (round 21, codex-1 + kimi-1):
  CLOSED in cycle 25, re-verified this round. The regression test "a link
  reached through an earlier link is walked from where it physically sits"
  passes at 7e8ccec.

- **Intermediate link expansion** (round 20, unanimous): CLOSED. The
  regression test "an intermediate link in a raw target is expanded before
  `..` is applied" passes at 7e8ccec.

- **POSIX backslash separator** (round 21, codex-1): CLOSED. The pin test
  "a backslash inside a POSIX link target is a filename byte, not a
  separator" passes at 7e8ccec. Labelled as a pin, not a proof — the
  labelling is honest (I confirmed this in round 22 and it remains so).

- **All prior rounds (7-19)**: CLOSED, as recorded in my round-22 review
  and the consensus. No new evidence emerged this round to reopen any.

- **Concurrent-installer isolation** (round 14, unanimous ruling): recorded
  follow-up, warning verbatim in CHANGELOG.md "Known limits" (lines 106-111).
  Not re-litigated.

No open findings remain from prior rounds.

## Position on the gate

**1 — correct as it stands.**

I held position 1 in round 22 and nothing has changed the gate's shape.
Cycle 26 is purely additive: it extracts existing logic into an injectable
helper and adds a regression that exercises the Windows arithmetic from a
POSIX host. The gate's three arms — crossing (touchpoints), same identity
(dev:ino), containment (chain inclusion) — are unchanged. The `walkRawTarget`
function now delegates to `rawTargetArithmetic` for root splitting and
relative-target anchoring, and the inline logic that was duplicated between
`splitAtRoot` and `walkRawTarget` is consolidated into one injectable
function.

I considered position 2 (narrow it) and refuse it for the same reason I gave
in round 22 and that codex-1 and kimi-1 both stated: CHANGELOG.md promises
"Installation and removal are atomic across the whole fleet" (line 49)
without excluding symlinked runtime homes, which earlier cycles deliberately
supported. Narrowing would withdraw a documented promise.

## New findings

**None.**

One observation, not a finding: the "an absolute raw target does not replay
its root as a component" test (test/bidding-addon.test.js:1978) still
carries the dead `const seen = []; const record = (p) => seen.push(p); ...
void seen; void record;` scaffolding that codex-1 flagged in round 22.
IMPLEMENTATION.md line 1848 says this scaffolding is "gone with the test it
belonged to" — that claim is inaccurate; the test still exists and still
carries the unused variables. This is cosmetic dead code in a test, not a
functional defect. It costs a real user nothing: the test passes, exercises
the POSIX absolute path through the public gate, and the Windows coverage
gap it failed to cover is now closed by the new `rawTargetArithmetic` test.
I do not raise it as a finding because it is not a defect in the shipped
code or the gate, and the brief instructs that "I have not exhausted every
filesystem exotic" is not a blocker — neither is an unused variable in a
test.

## Release judgement

Releasable as 2.1.0. The payload is unchanged since 714712f (zero diff under
`git diff 714712f..7e8ccec -- skills/parley-bidding/`). No round has found a
defect in it. The destination-collision gate has been the sole focus of
cycles 10-26, and cycle 26 closes the last outstanding block: the Windows
drive/UNC arithmetic is now extracted, injectable, and asserted with a
regression that fails at the two commits it names. All 372 node tests pass,
54 Python tests pass under 3.14.6, the manifest check is ok (47 files,
aggregate unchanged since 714712f), and npm pack is clean (202 files, 48
under parley-bidding, zero caches). Nothing remains that must change for
release.

## What I verified

1. **Test suite — full run.** `npm test` at 7e8ccec with
   `PATH="/opt/homebrew/bin:$PATH"` (python3 3.14.6): 372 node tests, 0 fail
   (52.0s). Python leg: 54/54 across 7 files on 3.14. Manifest check: 47
   files, `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.

2. **Python 3.9.6 refusal.** The test suite's node leg passes under
   `/usr/bin/python3` 3.9.6 (372/0). The Python leg refuses 3.9.6 by design
   (the add-on declares `>=3.10`). Verified consistent with prior rounds.

3. **Cycle-26 diff.** `git diff b1f43e4..7e8ccec`: 2 files, +184/-13.
   `lib/installer.js`: `rawTargetArithmetic` extracted and exported;
   `walkRawTarget` refactored to call it. `test/bidding-addon.test.js`: one
   new regression test ("root arithmetic is platform-correct") plus three
   tests carried forward from cycle 25 (linked-ancestor, backslash pin,
   intermediate link expansion). The diff is purely additive in the test
   file — no test was removed or weakened.

4. **Regression discrimination — verified against both named commits.**
   Extracted b1f43e4 and 2b7ca3e via `git archive` to
   `/tmp/parley-r23.0J4KAZ/{b1f43e4,2b7ca3e}`, installed deps, injected the
   7e8ccec test file, and ran the root-arithmetic test:
   - At b1f43e4: 1 fail — `TypeError: installer.rawTargetArithmetic is not
     a function`.
   - At 2b7ca3e: 1 fail — same TypeError.
   - At 7e8ccec: 1 pass.
   The function does not exist at either prior commit. At 2b7ca3e the inline
   logic was also materially wrong: `path.isAbsolute('C:\target\x')` returns
   false on POSIX, so the drive letter became a path component, and
   `rawTarget.split(/[\\/]+/)` tore UNC server/share into ordinary
   components. Confirmed by tracing both paths through `path.posix`.

5. **rawTargetArithmetic edge cases.** Traced the function with `path.win32`
   and `path.posix` for: drive-only root (`C:\` → start `C:\`, parts `[]`),
   UNC root only, relative with `.`/`..` (preserved as parts, handled by the
   walk loop), empty target, separator-only target, POSIX backslash-in-
   filename, POSIX absolute, and Windows path with forward slashes. All
   behave correctly: root is stripped before splitting, `.`/`..` are
   preserved for the walk loop, and POSIX backslashes stay as filename bytes.

6. **walkRawTarget integration.** `walkRawTarget` at line 1430 calls
   `rawTargetArithmetic(from, rawTarget, path, resolver)` with the real
   `path` impl and a `realpathSync` resolver with lexical fallback. The
   test calls it with `path.win32` to simulate Windows on a POSIX host —
   the injectable design is correct. `resolutionTouchpoints` (line 1472)
   continues to use `splitAtRoot` for the destination path, which is a
   different concern (destination root/components vs. raw link target
   root/components). Both are now injectable.

7. **All seven collision-arm tests.** Ran the seven accumulated collision
   regressions as a single selection: 7/7 pass.

8. **Payload stability.** `git diff 714712f..7e8ccec -- skills/parley-bidding/`
   produces zero output. The payload is unchanged since 714712f.

9. **Broader diff since round-7 tree.** `git diff 49fc3ec..7e8ccec --stat`:
   6 files changed (+2571/-167). All changes are in the installer, test
   suite, addon-manifest module, Python test runner, and CHANGELOG — no
   payload changes.

10. **npm pack.** `npm pack --dry-run --ignore-scripts --json`: 202 files,
    48 under `skills/parley-bidding/`, zero `__pycache__`/`.pyc`/`.pyo`,
    no `.git` entries.

11. **IMPLEMENTATION.md.** Frontmatter: `status: fix-up-cycle-26`,
    `head-commit: 7e8ccec`. Cycle-26 section (lines 1825-1888) is fully
    written and matches the shipped tree. The one inaccurate claim — that
    the dead `seen`/`record` scaffolding is "gone" (line 1848) — is noted
    above as a non-blocking observation; the scaffolding is still present
    in the absolute-target test at line 1982.

12. **CHANGELOG.md.** "Known limits" section (lines 106-111) covers
    concurrent-installer isolation verbatim. Atomicity claim at line 49
    does not exclude symlinked runtime homes — consistent with the refusal
    of position 2.

13. **Working tree integrity.** No mutations. `git status --porcelain`
    empty before and after. Temporary directory
    `/tmp/parley-r23.0J4KAZ` removed with `rm -rf`. Nothing written under
    the repo.
