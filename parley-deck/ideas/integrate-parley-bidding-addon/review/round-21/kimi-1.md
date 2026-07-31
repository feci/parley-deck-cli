---
idea: integrate-parley-bidding-addon
review-round: 21
agent: kimi-1
date: 2026-07-31
reviewed-commit: 381e639
---

## Verdict
BLOCK

## Outstanding findings — closed or not

Round 20's two named arms are closed, and I verified the closing rather than trusting the
record:

- **Intermediate links are expanded when entered** — closed. The cycle-24 regression test fails
  at `2b7ca3e` (I ran the `381e639` test file against a `git archive` copy of `2b7ca3e`:
  `actual: true, expected: false` at the `result.ok` assertion) and passes at `381e639`.
  0/1 → 1/1 as claimed.
- **Absolute target's root replayed as a component** — closed by arithmetic. The pin test
  passes at both commits, as labelled; see MINOR 2 on what that pin cannot cover.

All six accumulated arms are refused at `381e639`; I read each of the six tests and ran the
full suite myself.

However, round 20's root cause — **the walk was lexical where the kernel is physical** — is
not fully closed. Cycle 24 made the *interior* of `walkRawTarget` physical and left its
*entry condition* lexical. That is the new MAJOR below, so the umbrella finding from round 20
reopens in narrowed form rather than closing.

## Position on the gate
**3 — still wrong. The arm: `walkRawTarget` is seeded from a lexical path, so a link reached
through an earlier symlink with depth skew walks the wrong tree.**

`resolutionTouchpoints` (lib/installer.js:1461-1480) accumulates `logical` by string join and
never rewrites it after a link is found. `walkRawTarget` *returns* the physical landing point
— the recursion uses it internally at line 1442 (`if (landed !== null) current = landed`) —
and the outer caller discards it at line 1479. So when a destination's path contains a second
symlink, that link's relative target is walked from `path.dirname(logical)` — the **spelling**
— not from the directory the link physically sits in. Any `..` in that target climbs the wrong
tree whenever the first link's target is deeper than its spelling.

Concretely, reproduced at `381e639` (twice), in the exact image of the round-19 arm:

- `home/KM/skills/parley-deck/transient` — real, inside kimi's installed core destination
- `home/alias -> R/D` — codex's runtime root is a symlink whose target is one level deeper
  than its spelling (the depth skew)
- `home/R/D/skills -> ../../KM/skills/parley-deck/transient/../../../../away` — the skills
  container is a symlink whose target transits kimi's destination and climbs out
- codex dest `home/alias/skills/parley-deck` resolves, by the kernel, to `home/away/...`,
  **through** `home/KM/skills/parley-deck`. I proved the transit: `realpathSync` lands at
  `away` only while `transient` exists; remove it and resolution throws.

The kernel splices the inner target at `R/D` (2 ups to `home`); the walk seeds at
`home/alias` (2 ups to `/`) and records `/KM/skills/...`, which cannot exist. Same-identity is
silent (the finals differ), containment is silent (`identityChain` stats spelling prefixes;
`ino(home/KM/skills/parley-deck)` is an intermediate of the expansion, never a spelling
prefix), and the touchpoint walk records a tree under `/` that contains nothing.

**Measured at `381e639`:** `ok: true`, **14 units written, 0 blocked**. After the fleet
commit, kimi's `replaced` removed `transient`, so codex's reported destination no longer
resolves, and the payload codex wrote sits orphaned at `home/away/parley-deck` — the verbatim
harm signature of rounds 16/17/19 ("the earlier unit left dangling, files orphaned at the
chain's end"), through a seam none of those fixes touched. **Control:** identical layout with
a real directory instead of the `alias` link is refused (`ok: false`, 0 written, "resolving
one passes through the other"), so the probe discriminates on exactly the lexical seed.

This is not a new exotic class that position 2 could pension off. It is the round-19 arm —
which all three of us rated MAJOR and claude-1 spent two cycles closing — composed with one
ordinary ingredient: a symlinked runtime root pointing one level deeper than its spelling
(`~/codex` → `~/Sync/dotfiles/codex`). The gate's own doctrine since cycle 24 is "expand a
link the moment you enter it"; the outer walker violates that doctrine at its entry.

**Remedy, verified:** make the outer walker keep the same invariant as the recursion —

```js
// lib/installer.js:1479
if (target !== null) {
  const landed = walkRawTarget(logical, target, record);
  if (landed !== null) logical = landed;
}
```

On a `git archive` copy of `381e639` with only this change: **366/366 node tests still pass**,
the arm is refused (`ok: false`, 0 written), and the control still refuses. It needs a
discriminating regression test beside it — the round-19 test's layout with the link's parent
behind a depth-skewed symlink — which fails at `381e639` and passes with the remedy.

## New findings

**MAJOR — lexical seed in `resolutionTouchpoints`.** As above. `walkRawTarget`'s physical
return value is discarded by its only external caller, so every link after the first in a
destination's path is walked from the spelling. Silent gate bypass with cross-destination
write, a dead reported destination, and orphaned files; measured, not hypothetical.

**MINOR 1 — the pin test cannot pin the Windows arithmetic it names.** `walkRawTarget`
hard-codes the host `path`, so the root-not-replayed fix it carries for Windows
(`C:\target\x`, UNC) is exercisable only by inspection from the POSIX CI host. Cycle 19
learned this lesson for the same class of bug: `splitAtRoot` takes an injectable `impl`
precisely so the win32 splitting can be executed, and its comment says why a `path.win32`
self-assertion pins nothing. The new walk is the shipped Windows channel's arithmetic with no
executable coverage.

**MINOR 2 — dead scaffolding in the pin test.** `test/bidding-addon.test.js` ("an absolute raw
target does not replay its root as a component") defines `seen`/`record`, never uses them, and
disposes of them with `void seen; void record;`. Either assert the recorded probe sequence
(which, with MINOR 1's injection, is what would actually pin the Windows behavior) or delete
the scaffolding.

## Release judgement

Not releasable as 2.1.0. One thing must change: the outer walk in `resolutionTouchpoints` must
propagate `walkRawTarget`'s physical landing into `logical` (remedy above, verified 366/366
plus arm refused), with a discriminating regression test. The payload itself remains
unimpeached — unchanged since `714712f`, manifest aggregate
`sha256:7854adf1…b95a6d` verified — and I found nothing else standing between this tree and
release.

## What I verified

- **366/366 node tests** at `381e639` (own run, node v26.5.0).
- **Python leg 54/54** on Homebrew python3 3.14.6 (own run). **3.9.6 refusal is fail-closed**:
  PATH restricted to `/usr/bin`, the runner exits 1 with `python3 is 3.9, but the add-on
  declares >=3.10`.
- **Manifest check ok** — 47 files, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, matching the claim.
- **Payload stability**: `git diff 714712f..381e639` over `skills/`, manifest tooling and
  package metadata shows only `lib/addon-manifest.js` (+86) — the rounds 11-14 trust-boundary
  work, which I read in full (manifest entry rule on every read, key grammar + confinement,
  hash inside try).
- **Discriminating regression**: `git archive 2b7ca3e` + the `381e639` test file — the
  intermediate-link test **fails** (`actual: true, expected: false`); the absolute-root pin
  passes at both commits, as labelled. Archive and all temp homes deleted afterward.
- **The six accumulated arms**: each test read; all pass at `381e639`.
- **The seventh arm**: reproduced twice at `381e639` (layout, kernel-transit proof, `ok:
  true`, 14 written, dead destination, orphan tree measured); control with a physical seed is
  refused; remedy verified on an archive copy (366/366, arm refused, control still refused).
- **Reviewers'-note item 1** (`design-addons.test.js` parameterization): the cycle adds
  assertions pinning the backslash-sentinel exception; nothing weakened.
- **`IMPLEMENTATION.md` cycle-24 record** matches the shipped tree; **CHANGELOG "Known
  limits"** carries the single-writer wording verbatim as ratified in round 14.
- **Working tree untouched**: `git status --porcelain` empty after all runs; all probe homes
  and archive copies removed (`rm -rf`), nothing written under the repo.
