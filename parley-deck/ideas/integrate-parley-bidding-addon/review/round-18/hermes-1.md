---
idea: integrate-parley-bidding-addon
review-round: 18
agent: hermes-1
date: 2026-07-31
reviewed-commit: 64e43f9
---

## Verdict

BLOCK

## Outstanding findings — closed or not

Round 17 hermes-1 MAJOR (symlink chains: touchpoints recorded only the entry, never intermediate hops): **closed.** `resolutionTouchpoints` now follows each symlink chain with `readlinkSync` and records every hop as a `physicalEntryKey`. The regression test "a symlink chain whose middle hop sits inside a destination refuses the plan" reproduces the exact round-17 shape (B/skills -> kimiCore/redirect -> away) and the plan is refused. Verified by running the test against the actual installer at 64e43f9.

Round 17 codex-1 MAJOR / kimi-1 MAJOR (firmlink respelling: realpath preserves both spellings, containment and touchpoints both silent): **closed.** Containment and touchpoint comparisons now run on `physicalKey` (dev:ino + tail) with prefix semantics, so both firmlink spellings anchor on the same inode. The regression test "a firmlink respelling of one directory is one destination" passes and correctly skips on volumes without firmlinks.

Round 17 kimi-1 NIT-1 (stale comment on `pathEntryExists` referencing deleted `uninstallSkillUnit`): **still open.** Line 2342 of lib/installer.js still reads "`skillUnitStatus`, and `uninstallSkillUnit`." That function was deleted in cycle 19. No behavioral impact. See NIT-1 below.

Round 16 codex-1 MINOR / kimi-1 MINOR (Python gate fails open): **closed.** Verified at 64e43f9: 54/54 on 3.14.6, refuses 3.9.6 by design.

Round 16 hermes-1 NIT (uninstall --dry-run per-action flag): **closed.** Verified.

Dead `preflightUninstallUnit`: **closed.** No references remain.

## New findings

### [MAJOR] Containment regression: physicalKey anchors on the nearest existing ancestor, so two physically nested destinations with different existing intermediates get different dev:ino prefixes and the prefix-based overlaps() test fails to detect the nesting

Cycle 21 switched the containment check from `resolvedDestination` (realpath string of nearest existing ancestor + full tail) to `physicalKey` (dev:ino of nearest existing ancestor + tail). The purpose was to close firmlink respellings, which it does. But `physicalKey` replaces the entire common prefix with a single dev:ino anchor, and when two nested paths anchor on **different existing directories**, their keys start with different dev:ino values. The `overlaps()` prefix test then fails even though one path is physically inside the other.

The old `resolvedDestination` preserved the full path string from the common realpath ancestor, so the nesting was always visible as a string prefix regardless of which intermediate directories existed.

**Reproduction (run against 64e43f9, actual installer, not a logic mock):**

1. Install kimi first: `install --target kimi --include-undetected --no-addons` with `KIMI_CODE_HOME=$HOME/K`. This creates `$HOME/K/skills/parley-deck` with real subdirectories including `references/`, `agents/`, etc.
2. Set `CODEX_HOME=$HOME/K/skills/parley-deck/references` — an existing real subdirectory inside kimi's managed tree.
3. Run `install --target all --include-undetected --no-addons` with both `KIMI_CODE_HOME` and `CODEX_HOME` set.

Measured result at 64e43f9: `ok: true`, codex `installed`, kimi `replaced`, codex's destination (`$HOME/K/skills/parley-deck/references/skills/parley-deck`) **does not exist on disk**, surviving marker names `target: kimi`. This is the single-process false-success / partial-fleet state that B5 and the changelog say cannot occur.

Measured result at fa3da41 (cycle 20): `ok: false`, both `blocked` with "Destination overlaps another in this plan — one contains the other." The old `resolvedDestination` containment caught it.

**Why the keys miss:**

- kimi dest = `$HOME/K/skills/parley-deck` — exists as a real dir, so `physicalKey` anchors on it: key = `dev:ino_pd`
- codex dest = `$HOME/K/skills/parley-deck/references/skills/parley-deck` — `references/` exists as a real dir (created by the kimi install), so `physicalKey` anchors on `references/`: key = `dev:ino_ref/skills/parley-deck`
- `overlaps(dev:ino_ref/skills/parley-deck, dev:ino_pd)` — `dev:ino_ref` != `dev:ino_pd`, no prefix match, returns false.

The touchpoint check also misses: there are no symlinks in the path, so `resolutionTouchpoints` returns an empty array and `crosses` is false.

**Why this is a regression, not a pre-existing gap:** fa3da41 caught this exact scenario. The switch from `resolvedDestination` to `physicalKey` for containment broke it. The existing "nested inside another destination" test passes because it does not pre-create the outer destination — both paths anchor on the same tmpdir, the keys share a prefix, and containment fires. The regression requires an intermediate existing directory inside the outer dest, which the test does not exercise.

**Uninstall is affected identically.** `removeFleetAtomically` calls the same `aliasedDestinations` gate. With both targets managed, `uninstall --target all` returns `ok: true`, both `removed`. In practice the quarantine rename of the outer (kimi) moves the inner's (codex's) tree aside with it, and phase B's `rmSync` deletes everything recursively — no debris, no warning, but the plan should have been refused.

**What the fix would need:** containment must work when two paths anchor on different existing ancestors but one is still physically nested inside the other. The cleanest approach is to walk the full ancestor chain (by `physicalKey` at each level) for each destination and check whether any ancestor of one equals the other's key — rather than relying on a single anchor and a prefix test on the tail. Alternatively, retain `resolvedDestination` for the containment check (where it was correct) and use `physicalKey` only for the equality and firmlink cases it was introduced to close.

**Release judgement:** this must be fixed before 2.1.0. It is a regression from fa3da41 on the same invariant — "No planned destination may equal, contain, or be contained by another" — that cycles 15 through 20 were built to enforce. The false-success state is the one the changelog and B5 explicitly say cannot occur.

### [NIT-1] pathEntryExists comment references deleted uninstallSkillUnit

lib/installer.js:2342 still reads "`skillUnitStatus`, and `uninstallSkillUnit`." That function was deleted in cycle 19. kimi-1 flagged this as NIT-1 in round 17; it remains. No behavioral impact.

### [NIT-2] resolvedDestination is dead code within aliasedDestinations

Line 1490 computes `entry.path = resolvedDestination(unit.dest)`, but `entry.path` is never read in the comparison logic — the crosses check uses `entry.key` and `entry.touchpoints`, and the containment check uses `entry.key`. `resolvedDestination` (24 lines) and its `entry.path` assignment are now dead within this function. The comment on `resolvedDestination` at line 1349 says "Used for CONTAINMENT, which a dev/ino key cannot express" — but containment is now expressed with dev/ino keys, contradicting the comment. No behavioral impact.

## Release judgement

Not releasable as 2.1.0. The one thing that must change: the containment check in `aliasedDestinations` must detect physically nested destinations when the two paths anchor on different existing ancestors. Cycle 21's switch from `resolvedDestination` to `physicalKey` for containment closed the firmlink case but opened a regression on the plain-directory nesting case that fa3da41 caught. The invariant "No planned destination may contain or be contained by another" is false as shipped, with the same single-process false-success state that blocked the last three rounds.

## What I verified

- Read the full cycle-21 diff (`git diff fa3da41..64e43f9`) and the complete `aliasedDestinations`, `physicalKey`, `physicalEntryKey`, `resolutionTouchpoints`, `resolvedDestination`, `overlaps`, `installFleetAtomically`, `removeFleetAtomically`, `quarantineName`, `copyPayloadAtomically`, `commitStagedUnit`, `revertStagedUnit`, `discardBackup`, `pathEntryExists`, `destinationAncestorObstacle`, `installerOwnsDestination`, and `targetSkillUnits` functions in lib/installer.js.
- Read all three round-17 reviews (hermes-1, codex-1, kimi-1) and the IMPLEMENTATION.md cycle-21 record.
- 359 node tests pass (`node --test`), 0 fail, from a `git archive` copy of 64e43f9.
- Python leg: 54/54 on python3 3.14.6; refuses 3.9.6 by design.
- Manifest check: 47 files, aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since 714712f.
- **The MAJOR was reproduced against the actual installer** (`require('../lib/installer')`), not a logic mock. Three independent runs from a `git archive` copy:
  - Install path: kimi installed first (creating real subdirs), then `--target all` with `CODEX_HOME` set inside kimi's `references/` subdir. Result: `ok: true`, codex `installed`, kimi `replaced`, codex dest absent, surviving marker `target: kimi`. Confirmed at 64e43f9.
  - Same scenario at fa3da41 (cycle 20): `ok: false`, both `blocked` with "Destination overlaps another in this plan." Confirmed regression.
  - Same scenario at 714712f: `ok: true`, codex `installed`, kimi `replaced`, codex dest absent. The containment check did not exist at 714712f (added in cycle 15), so this was never caught until cycle 20, and cycle 21 broke it.
  - Uninstall path: both targets managed (kimi installed, codex re-installed alone after the first false-success), `uninstall --target all`. Result: `ok: true`, both `removed`, no debris. The gate should have refused the plan.
- Chain walk cycle termination: tested with a symlink cycle (A/skills -> B/skills -> A/skills) at both the top level (caught by `destinationAncestorObstacle`) and inside a managed tree (caught by the `seen` set in `resolutionTouchpoints`). Neither hangs; both terminate within the 64-step limit.
- `physicalEntryKey` case handling: lowercases the base name on `CASE_INSENSITIVE_FS` (darwin/win32). Correct. Unicode normalization differences in the tail name would produce different keys, but this is the same inherent limitation the old `resolvedDestination` had and requires deliberately constructing files with different normalizations — not a new issue.
- `physicalEntryKey` when parent is reached through a link: `physicalKey(path.dirname(resolved))` uses `statSync` which follows the link, returning the target's dev:ino. The entry key is then target_dev:ino + "/" + name. This is correct — the entry physically lives at the target's location.
- Stale comment NIT-1: confirmed `uninstallSkillUnit` is referenced at line 2342 but does not exist in the tree.
- Dead code NIT-2: confirmed `entry.path` is computed at line 1490 but never read in `aliasedDestinations`. Grepped the function body for `.path` — zero uses.
- All temporary archives and probe homes removed. Repo remained clean throughout: `git status --porcelain` empty before and after.
