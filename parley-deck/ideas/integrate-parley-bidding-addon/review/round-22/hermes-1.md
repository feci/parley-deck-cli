---
idea: integrate-parley-bidding-addon
review-round: 22
agent: hermes-1
date: 2026-07-31
reviewed-commit: b1f43e4
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All findings from rounds 7 through 21 are closed. I checked each accumulated
arm against the code at b1f43e4 and its test:

- **Linked-ancestor / physical-parent arm** (round 21, codex-1 MAJOR + kimi-1
  MAJOR — the same arm, the one that split round 21): CLOSED. Two fixes in
  cycle 25, working together:

  1. `resolutionTouchpoints` now propagates `walkRawTarget`'s return value
     into `logical` (lib/installer.js:1495-1496). The old code discarded it,
     so after expanding an ancestor link, a later link's relative target was
     walked from the spelling. kimi-1 gave this remedy exactly; it is applied
     verbatim.

  2. `walkRawTarget` now anchors a relative target at
     `fs.realpathSync(path.dirname(from))` instead of `path.dirname(from)`
     (lib/installer.js:1420-1425). The lexical anchor meant that once an
     earlier ancestor link had been expanded, every `..` in a relative target
     climbed a different tree than the kernel climbs. codex-1's related point
     — a relative target is anchored at the link's physical parent — is
     applied. The `realpathSync` call is in a try/catch that falls back to the
     lexical parent if the physical parent does not exist yet, so a not-yet-
     created destination does not crash the walk.

  The discriminating regression ("a link reached through an earlier link is
  walked from where it physically sits") FAILS at 381e639 and PASSES at
  b1f43e4. I verified this by extracting the 381e639 tree via `git archive`,
  injecting the b1f43e4 test file, and running the single test: `actual:
  true, expected: false` at the `result.ok` assertion. At b1f43e4 the same
  test passes. See "What I verified" for the commands.

- **Backslash separator on POSIX** (round 21, codex-1): CLOSED.
  `split(/[\\/]+/)` treated `\` as a separator on every platform; on POSIX it
  is an ordinary byte in a filename. Now platform-conditional:
  `process.platform === "win32" ? /[\\/]+/ : /\/+/` (lib/installer.js:1412).
  The pin test ("a backslash inside a POSIX link target is a filename byte,
  not a separator") passes at both 381e639 and b1f43e4. The brief and the test
  body both label it a pin, not a proof — the parent component is recorded
  before the torn name, so the arm is caught by accident at 381e639. I
  confirmed this labelling is honest: the test passes at both commits.

- **Intermediate link expansion** (round 20, unanimous): `walkRawTarget`
  lstat's each component after joining, and if it is a symlink, reads the
  target and recurses before continuing. `current` is updated to the landed
  expansion. Fixed. The regression fails at 2b7ca3e and passes at 381e639 and
  b1f43e4.

- **Absolute root replayed as component** (round 20, codex-1 + hermes-1):
  root is the starting point only; parts are sliced from `root.length`. Fixed.
  The pin test passes at both commits, as labelled.

- **Raw-target `..`** (round 19, codex-1 MAJOR): `walkRawTarget` walks one
  component at a time and applies `..` only after recording the entry. Fixed.

- **Firmlink respelling** (round 17, codex-1 MAJOR): `identityChain` uses
  `statSync` (follows links) and keys on `dev:ino`. Fixed.

- **Plain nesting with skewed anchors** (round 18, unanimous): `identityChain`
  builds a full ancestry chain of physical identities; containment is checked
  by chain inclusion, not scalar prefix. Fixed.

- **Symlink entry inside a destination** (round 17, hermes-1 MAJOR): recorded
  by `resolutionTouchpoints` outer loop. Fixed.

- **Symlink chain whose middle hop sits inside a destination** (round 17):
  the recursive `walkRawTarget` follows chain hops, up to depth 32. Fixed.

- **Concurrent-installer isolation** (round 14, unanimous ruling): recorded
  follow-up, warning verbatim in CHANGELOG.md under "Known limits". Not
  re-litigated.

No open findings remain from prior rounds.

## Position on the gate

**1 — correct as it stands.**

Round 21's split has closed. Two of three reviewers examined the anchor and
found the arm; I examined the expansion and judged it correct. Cycle 25 fixes
the anchor — the part I did not examine closely enough — and the expansion I
vouched for still works. The walk is now physical where the kernel is
physical, at both the entry condition (physical parent via `realpathSync`)
and the expansion (immediate symlink recursion with `current = landed`).

The gate models three independent failure modes and checks them in the right
order:

1. **Crossing** (`touchpoints`): one destination's resolution passes through
   the other's tree. `resolutionTouchpoints` walks the dest path component by
   component; for each symlink, it records the entry and calls `walkRawTarget`
   on its raw target. `walkRawTarget` now models the kernel: it anchors
   relative targets at the physical parent, expands intermediate links the
   moment they are entered, and returns the physical landing point — which
   the caller propagates into `logical`, so a later link's relative target
   is also walked from the right tree. This is the arm cycles 19-25 closed.

2. **Same identity** (`a.identity === b.identity`): two destinations resolve
   to the same physical object. Keyed by `dev:ino` of the final existing
   ancestor, which handles firmlink respellings and aliased runtime roots.

3. **Containment** (`a.chain.includes(b.identity)`): one destination is inside
   the other. Keyed by full ancestry chain, which handles skewed anchors.

I considered position 2 (narrow it) and I accept codex-1's round-21 argument
against it: symlinked runtime homes were deliberately supported by earlier
fixes in this idea, and CHANGELOG.md promises "Installation and removal are
atomic across the whole fleet" without excluding them. Narrowing would
withdraw a documented promise rather than trim an untouched edge. If the
project narrows later it must rewrite that promise in the same change. The
gate is the right shape.

## New findings

**None.**

I examined the cycle-25 diff in full (2 files, +108/-4), the full installer
(2621 lines), the addon-manifest module, the Python test runner, the complete
test file (113 bidding-addon tests), the design-addons test, and the
CHANGELOG. I verified the payload has not changed since 714712f (zero diff
under `git diff 714712f..b1f43e4 -- skills/parley-bidding/`).

The stale comment block I noted at round 21 (the old lexical-walk description
at lines 1387-1392) has been replaced by an accurate comment describing the
physical walk and the cycle-25 fixes. No cosmetic issues remain.

## Release judgement

Releasable as 2.1.0. The payload is unchanged since 714712f and no round has
found a defect in it. The destination-collision gate has been the sole focus
of cycles 10-25, and the cycle-25 fix closes the last arm: the resolution
walk is now physical where the kernel is physical, at both the entry
condition and the expansion. All 368 node tests pass, 54 Python tests pass
under 3.14.6, the manifest check is ok (47 files, unchanged aggregate), and
the one discriminating regression fails at 381e639 and passes at b1f43e4.
The backslash pin is honestly labelled and passes at both commits.

## What I verified

1. **Test suite — full run.** `npm test` at b1f43e4: 368 node tests, 0 fail
   (node v26.5.0, 43.9s). Python leg: 54/54 under `/opt/homebrew/bin/python3`
   3.14.6; refuses 3.9.6 by design (the add-on declares `>=3.10`). Manifest
   check: 47 files,
   `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`,
   unchanged since 714712f.

2. **Cycle-25 diff.** `git diff 381e639..b1f43e4`: 2 files, +108/-4.
   `walkRawTarget` anchors relative targets at `realpathSync(lexicalParent)`
   with a lexical fallback; `resolutionTouchpoints` propagates `walkRawTarget`'s
   return value into `logical`; separators are platform-conditional. Two new
   tests added.

3. **Test discrimination — verified against old commit.** Extracted the
   381e639 tree via `git archive 381e639 | tar -x` to a temporary directory,
   installed deps, injected the b1f43e4 test file, and ran both new tests:
   - "a link reached through an earlier link is walked from where it physically
     sits" — FAILS at 381e639 (`actual: true, expected: false` at the
     `result.ok` assertion). PASSES at b1f43e4. Discriminating.
   - "a backslash inside a POSIX link target is a filename byte, not a
     separator" — PASSES at both commits. Labelled as a pin in the test body;
     the labelling is honest.

4. **Old tree full bidding-addon suite.** The 381e639 archive with the
   b1f43e4 test file injected: 113 tests, 112 pass, 1 fail — exactly the one
   discriminating test. Confirms no other test was broken by the cycle-25
   changes.

5. **Payload stability.** `git diff 714712f..b1f43e4 -- skills/parley-bidding/`
   produces zero output. The payload is unchanged since 714712f. The only
   post-714712f changes outside the gate are `lib/addon-manifest.js` (+86,
   rounds 11-14 trust-boundary work), `scripts/run-python-tests.js` (+28/-10,
   rounds 15-16 fail-closed interpreter gate), and `test/design-addons.test.js`
   (+22, round 1 backslash comment clarification) — all previously reviewed.

6. **Edge cases.** `realpathSync` failure (parent does not exist yet): falls
   back to lexical parent, walk continues, no crash. `..` at root:
   `path.dirname('/')` returns `'/'` — no loop, no crash. Circular symlinks:
   depth > 32 returns null; caller does not update `current`/`logical`; the
   symlink entry was already recorded before the call. POSIX separators:
   backslash is not split (`'a\\b'.split(/\/+/)` → `['a\\b']`); forward slash
   is. Non-existent target paths: walk returns the computed path, subsequent
   lstat fails gracefully, no crash.

7. **Logical trace of the discriminating test.** The test layout:
   `linkA -> real/A`, Codex container at `linkA/container`, skills symlink at
   `real/A/container/skills` with target `../../Btree/skills/parley-deck/inner/deep`,
   `inner` inside kimi's core is a symlink to `outside`. The kernel expands
   `linkA` first, so `..` climbs from `real/A/container` to `real`, then
   `Btree/skills/parley-deck` enters kimi's core. At 381e639, `walkRawTarget`
   anchored at `path.dirname(from)` (lexical), so `..` climbed from the
   spelling through `linkA` to `home`, and `home/Btree` does not exist —
   kimi's core was never recorded, crossing was not detected, `ok: true`. At
   b1f43e4, `realpathSync` resolves the anchor to `real/A/container`, `..`
   climbs to `real`, `Btree` enters kimi's core, the crossing is detected,
   `ok: false`, 0 units written.

8. **Package check.** `npm pack --dry-run --ignore-scripts --json`: 202 files,
   48 under `skills/parley-bidding/`, no `.pyc`/`.pyo`/`__pycache__`, no `.git`
   entries.

9. **IMPLEMENTATION.md.** Frontmatter: `status: fix-up-cycle-25`,
   `head-commit: b1f43e4`. Cycle-25 section is fully written and matches the
   shipped tree: the arm, the separator, the position on the gate, the pin
   labelling, the discrimination table, and the measured results (368 node
   tests, 54 Python tests, manifest aggregate).

10. **Working tree integrity.** No mutations. `git status --porcelain` empty
    before and after. All temporary directories and scripts (`git archive`
    copy at `/tmp/parley-r22.I8FhTQ`, trace/edge scripts at `/tmp/*.js`)
    removed with `rm -rf`. Nothing written under the repo.
