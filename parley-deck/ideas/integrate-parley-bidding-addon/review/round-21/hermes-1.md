---
idea: integrate-parley-bidding-addon
review-round: 21
agent: hermes-1
date: 2026-07-31
reviewed-commit: 381e639
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

All findings from rounds 7 through 20 are closed. I checked each accumulated arm
against the code at 381e639 and its test:

- Raw-target `..` (round 19, codex-1 MAJOR): `walkRawTarget` walks one component
  at a time and applies `..` only after recording the entry. Fixed.
- Intermediate link expansion (round 20, unanimous): `walkRawTarget` now lstat's
  each component after joining, and if it is a symlink, reads the target and
  recurses before continuing. The `current` variable is updated to the landed
  expansion, so subsequent `..` steps back through the real tree. Fixed.
- Absolute root replayed as component (round 20, codex-1 + hermes-1): the root
  is now the starting point only; parts are sliced from `root.length`, so the
  root is not re-entered. I verified with `path.win32` that `C:\target\x`
  produces root `C:\` and parts `['target','x']`, and a UNC path
  `\\server\share\dir\x` produces root `\\server\share\` and parts `['dir','x']`.
  Fixed.
- Symlink entry inside a destination (round 17, hermes-1 MAJOR): recorded by
  `resolutionTouchpoints` outer loop. Fixed.
- Chain hop inside a destination (round 17, hermes-1 MAJOR): the recursive
  `walkRawTarget` follows chain hops — each expansion that lands on a symlink
  recurses again, up to depth 32. The old outer hop loop is gone, replaced by
  recursion. Fixed.
- Firmlink respelling (round 17, codex-1 MAJOR): `identityChain` uses `statSync`
  (follows links) and keys on `dev:ino`, so `/private/x` and
  `/System/Volumes/Data/private/x` share an identity. Fixed.
- Plain nesting with skewed anchors (round 18, unanimous): `identityChain`
  builds a full ancestry chain of physical identities; containment is checked by
  chain inclusion, not scalar prefix. Fixed.

No open findings remain from prior rounds.

## Position on the gate

**1 — correct as it stands.**

The gate models three independent failure modes and checks them in the right
order:

1. **Crossing** (`touchpoints`): one destination's resolution passes through
   the other's tree. This is the case the identities alone cannot express — a
   symlink whose target walks into another destination and back out. The
   touchpoint set is built by `resolutionTouchpoints`, which walks the dest path
   component by component, and for each symlink encountered, records the entry
   and calls `walkRawTarget` on its raw target. `walkRawTarget` now models the
   kernel: it expands intermediate links the moment they are entered, so `..`
   after an expansion steps through the real tree, not the spelling. This is the
   arm that rounds 19 and 20 fixed, and it is now correct.

2. **Same identity** (`a.identity === b.identity`): two destinations resolve to
   the same physical object. Keyed by `dev:ino` of the final existing ancestor,
   which handles firmlink respellings and aliased runtime roots.

3. **Containment** (`a.chain.includes(b.identity)`): one destination is inside
   the other. Keyed by full ancestry chain, which handles skewed anchors where
   both sides have different nearest existing ancestors.

I considered position 2 (narrow it). The argument for narrowing is that symlink
chains, firmlinks, and Windows junctions require a deliberately unusual
configuration. But the gate does not penalize those configurations — it refuses
a *plan* where two destinations collide through them. A user with such a
configuration who points two targets at independent paths will not be refused.
The gate only acts when the collision is real, and the six arms are the real
ways two destinations can collide on a filesystem that has symlinks. Moving them
to "Known limits" would mean shipping a gate that silently allows orphaned files
and dangling destinations for configurations the gate can detect. That is the
class of bug this idea has spent fourteen cycles closing. The gate is the right
shape.

## New findings

**None.**

I examined the cycle-24 diff in full, the full installer (2604 lines), the
addon-manifest module (309 lines), the Python test runner, the SKILL.md, the
complete test file (2351 lines, 111 bidding-addon tests), and the CHANGELOG. I
verified the payload has not changed since 714712f (zero diff under
`git diff 714712f..381e639 -- skills/parley-bidding/`).

One stale comment block (lines 1387-1392 in `lib/installer.js`) describes the
old lexical walk that was replaced. It is superseded by the new comment at
1393-1402. This is cosmetic and does not affect correctness — not a finding
worth another cycle.

## Release judgement

Releasable as 2.1.0. The payload is unchanged since 714712f and no round has
found a defect in it. The destination-collision gate has been the sole focus of
cycles 10-24, and the cycle-24 fix closes the last arm: the resolution walk is
now physical where the kernel is physical. All 366 node tests pass, 54 Python
tests pass under 3.14.6, the manifest check is ok (47 files, unchanged
aggregate), and the single discriminating regression fails at 2b7ca3e and passes
at 381e639.

## What I verified

1. **Test suite — full run.** `npm test` at 381e639: 366 node tests, 0 fail.
   Python leg: 54/54 under python3.14 (3.14.6); refuses 3.9.6 by design (the
   add-on declares `>=3.10`). Manifest check: 47 files,
   `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`,
   unchanged since 714712f.

2. **Cycle-24 diff.** `git diff 2b7ca3e..381e639`: 2 files, +105/-22.
   `walkRawTarget` rewritten to recurse with immediate symlink expansion;
   absolute root no longer replayed as component; outer hop loop in
   `resolutionTouchpoints` removed (replaced by recursion). Two new tests added.

3. **Test discrimination — verified against old commit.** Extracted the
   cycle-24 tree at 2b7ca3e via `git archive` to a temporary directory, injected
   the two new tests, and ran them. Result: "an intermediate link in a raw
   target is expanded before `..` is applied" FAILS (ok:true where false is
   expected — plan accepted, files would be orphaned). "an absolute raw target
   does not replay its root as a component" PASSES (correctly labelled as a pin
   — it asserts a property cycle 24 preserves, not one it introduces). Temp
   directory cleaned up.

4. **walkRawTarget logic trace.** Traced the intermediate-link test scenario
   component by component. The old lexical walk recorded `home/mid`,
   `home/mid/transient`, and `/away`, but `entryChain(home/mid/transient)`
   resolves through the symlink only to `subdir`'s `dev:ino` — not to
   `kimiCore`'s. So the crosses check did not fire and the plan was accepted.
   The new physical walk expands `home/mid` immediately (it is a symlink to
   `kimiCore/subdir`), recurses, and records `kimiCore` directly. The crosses
   check fires and the plan is refused.

5. **resolutionTouchpoints outer loop interaction.** The outer loop does not
   update `logical` to the expansion after a symlink — it continues walking the
   original spelling. This is correct: `lstatSync` on `path.join(logical, part)`
   resolves symlink prefixes transparently through the OS. `walkRawTarget`
   handles target expansion (including intermediate symlinks); the outer loop
   handles dest-path components. Together they cover all touchpoints.

6. **Edge cases.** `..` at root: `path.dirname('/')` returns `'/'` — no crash,
   no loop. Circular symlinks: depth > 32 returns null, caller does not update
   `current`, walk continues from the symlink's own path (already recorded).
   Symlink to `.`: resolves to containing directory. Symlink to `..`: resolves
   to grandparent. All correct.

7. **Windows root extraction.** Verified with `path.win32`: `C:\target\x`
   produces root `C:\` and parts `['target','x']` (new code) vs parts
   `['C:','target','x']` (old code, which probed `C:\C:\target`). UNC path
   `\\server\share\dir\x` produces root `\\server\share\` and parts
   `['dir','x']`. Both correct.

8. **Full-scope review.** Payload unchanged since 714712f (zero diff).
   SKILL.md (172 lines): well-structured, HITL classification E0-E8, lifecycle
   states, deterministic tools — no issues. addon-manifest.js (309 lines):
   lstat rule on every manifest read, key confinement, aggregate digest
   verification, undeclared file detection — no issues. run-python-tests.js
   (182 lines): per-file expected counts, fails on missing interpreter, fails
   on unparseable version, reads floor through shared parser — no issues.
   CHANGELOG: "Known limits" section carries the concurrent-installer warning
   verbatim per the round-14 ruling. design-addons.test.js: backslash
   continuation comment clarified (round 1, codex-1 MINOR) — no issues.

9. **Working tree integrity.** No mutations. `git status` clean before and
   after. Temp directory created for old-commit comparison was removed.
