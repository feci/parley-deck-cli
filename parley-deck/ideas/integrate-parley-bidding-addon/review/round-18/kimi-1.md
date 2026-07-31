---
idea: integrate-parley-bidding-addon
review-round: 18
agent: kimi-1
date: 2026-07-31
reviewed-commit: 64e43f9
---

## Verdict

BLOCK

## Outstanding findings — closed or not

- **Round 17 MAJOR (firmlink respelling aliases a destination).** CLOSED. The regression at
  `test/bidding-addon.test.js:1655` runs on this volume (`/private/tmp` firmlink pair present)
  and passes; verified in my own 359/359 suite run at `64e43f9`. Key-prefix equality on the
  shared inode anchor is the right close for that case.
- **Round 17 MAJOR (symlink-chain intermediate hops invisible).** CLOSED. `resolutionTouchpoints`
  (`lib/installer.js:1392`) now follows each chain with `readlink` and records every hop as a
  `physicalEntryKey`; the regression passes. Cycle termination measured separately: a symlink
  loop (`a ↔ b`) on a destination path terminates the walk (bounded at 64 hops) and the plan
  fails closed ("Destinations depend on each other"), no hang.
- **Round 16 MAJOR (destination reached through another destination).** Still closed; that
  regression is intact in the suite.
- **Round 15 MAJOR (nested destinations) — REOPENED by cycle 21.** The containment guarantee
  holds only for pairs that anchor on the same inode. See MAJOR-1.
- Round-14 settled matter (single-writer limit) not re-examined; the CHANGELOG "Known limits"
  wording is present verbatim.

## New findings

### MAJOR-1 — `dev:ino/tail` keys are prefix-comparable only when both paths anchor on the same inode; nested destinations with anchors at different depths evade both plan guards, and the fleet then reports `ok: true` while deleting the inner unit's fresh install

`physicalKey` (`lib/installer.js:1323`) keys a destination by the dev:ino of its nearest
**existing** ancestor plus the lexical tail below it, and `overlaps` (`lib/installer.js:1451`)
treats two such keys as prefix-comparable path strings. That is sound exactly when the two
paths anchor on the **same** inode — the round-17 firmlink case, where both spellings resolve
to one directory. When two physically nested destinations anchor at **different** depths —
because the inner unit's path already exists one or more components below the outer unit's
destination — the keys carry different `dev:ino` tokens and are string-incomparable, so
`overlaps` is silent. Cycle 21 *replaced* the resolved-string containment with key containment
(`git diff fa3da41..64e43f9`, `lib/installer.js` hunk at the containment arm), instead of adding
identity alongside strings. The two encodings have complementary blind spots: firmlink
respellings share an inode anchor (keys catch them, strings cannot); depth-skew nesting keeps
the lexical prefix (strings catch it, keys cannot). Only the union is complete.

Both guards that consume these keys are affected, and I measured both end-to-end against a
`git archive` copy (working tree untouched), no `--force` anywhere:

**Arm A — containment, no symlinks at all.**

1. `HOME=<tmp>` `install --target claude` — a legitimate prior install, so claude's
   destination `<H>/.claude/skills/parley-deck` exists and is installer-owned.
2. `mkdir <H>/.claude/skills/parley-deck/deep` — one pre-existing component inside it
   (pollution; nothing here requires privilege or force).
3. `HOME=<H> CODEX_HOME=<H>/.claude/skills/parley-deck/deep install --target all
   --include-undetected --json` — codex's destination becomes
   `<claude-dest>/deep/skills/parley-deck`, nested inside claude's.

At `64e43f9`, measured: **`ok: true`.** codex core plus all five add-ons reported `installed`;
claude core reported `replaced`; codex's reported destination **does not exist afterwards** —
codex's staged trees were committed inside claude's old tree (phase 2, codex first), claude's
commit then renamed that tree into its backup, and `discardBackup` deleted it. The surviving
marker names claude. This is the round-15 shape exactly: single-process false success, partial
fleet, contradicting the B5/changelog claim.

At `fa3da41`, identical scenario: `ok: false`, both units blocked — "Destination overlaps
another in this plan: codex/parley-deck and claude/parley-deck" — fleet untouched. So this is
a cycle-21 regression of a guard that worked in cycle 19/20, not a pre-existing gap.

**Arm B — dependency (touchpoints), one symlink.**

Same setup, but `<x>/skills` is a symlink to `<claude-dest>/sub/deep` (two existing levels
below claude's destination) and `CODEX_HOME=<x>`. codex's destination resolves into claude's
tree through the link. At `64e43f9`, measured: **`ok: true`**, codex `installed`, its
destination dangling afterwards (payload again deleted inside claude's discarded backup). At
`fa3da41`: blocked. The cross arm misses because a chain hop is keyed by `physicalEntryKey` —
its *parent's* identity plus its name — so a hop landing at depth ≥ 2 below the other unit's
anchor inode is incomparable with that unit's key (at depth exactly 1 the hop's entry key
carries the other unit's anchor and the check fires; the boundary is precisely anchor depth).

**Why the suite is green:** the round-15 regression (`test/bidding-addon.test.js:1769`) builds
both destinations under a fresh `tmpDir()` home, so both anchor on the *same* inode and the
tail prefix fires. The skew needs one pre-existing component below the outer destination —
which neither that test nor any other arm supplies.

**Fix (small):** containment and dependency must be the **union** of both comparisons —
physical-key prefix (closes firmlinks) **or** resolved-string prefix (closes depth skew).
`resolvedDestination` is still computed for every unit at `lib/installer.js:1490` and, since
cycle 21, never read — the string arm is one comparison away, and that now-dead `entry.path`
field is itself the tell (delete it only if the union is implemented some other way). Add both
measured scenarios as regressions, asserting not only `ok: false` but that the fleet is
untouched — Arm A's failure mode is silent deletion, which an `ok`-only assertion can miss.
Uninstall needs no separate work: `removeFleetAtomically` consults the same `aliasedDestinations`
first (`lib/installer.js:1669`), so it inherits both the gap and the fix.

### NIT-1 — tails are lowercased but not Unicode-normalized

`physicalKey`/`physicalEntryKey` fold case on darwin/win32 but do not `.normalize()`. Measured
on this volume (APFS, case- and normalization-insensitive): an NFD spelling `stat`s an
NFC-created directory, and `rename` onto the NFD spelling silently replaces the NFC directory.
Two units whose not-yet-created tails differ only in Unicode normalization would key
differently while naming one directory — the same aliasing hole the lowercase fold closed for
case, one equivalence class over. Reachability is thin: every installer-constructed tail
component is ASCII (`SAFE_SKILL_NAME` for recorded names, fixed skillDir segments, fixed skill
names), so this needs env-supplied non-ASCII, not-yet-existing components in two units' tails
in one plan. Recorded follow-up, not a blocker; a `.normalize("NFC")` beside the existing
lowercase fold closes it.

## Release judgement

Not releasable as 2.1.0. The one thing that must change: `aliasedDestinations` must compare
**both** physical-key prefixes and resolved-string prefixes (union), restoring the depth-skew
containment/dependency coverage cycle 21 dropped while keeping the firmlink coverage it added —
with the two measured scenarios above added as regressions that also assert the fleet is
untouched.

## What I verified

- Read `lib/installer.js` (2601 lines), `lib/addon-manifest.js` (309), and the full
  `fa3da41..64e43f9` diff; traced `physicalKey`, `physicalEntryKey`, `resolutionTouchpoints`,
  `overlaps`, `aliasedDestinations`, and both fleet transactions.
- Measured MAJOR-1 Arm A and Arm B against `git archive` copies of `64e43f9` (both: `ok: true`,
  inner unit `installed`, its destination absent/dangling afterwards) and `fa3da41` (both:
  `ok: false`, blocked, fleet untouched). Working tree never mutated; all scratch under `/tmp`,
  removed afterwards (including the `node_modules` symlink before `rm -rf`).
- Symlink cycle (`a ↔ b`) on a destination path: walk terminates (bounded at 64 hops), plan
  fails closed, no hang.
- Suite at `64e43f9` from the archive with the repo's `node_modules` symlinked in:
  **359 node tests, 0 fail**; python leg **54/54 on 3.14.6**; manifest check ok, 47 files,
  aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d` — all as
  recorded. (A bare archive run fails `design-addons.test.js` on missing devDependency
  `commonmark`; that is an artifact of the copy, not of the commit.)
- Volume semantics for NIT-1 measured on this machine's APFS: NFD spelling of an NFC directory
  `stat`s successfully; rename onto the NFD spelling replaces silently.
- `addon-manifest.js` swept for the full-scope pass: manifest key confinement, entry rules
  (`lstat`, symlink refusal), aggregate binding — no new findings there.
- Chain-walk coverage of uninstall's quarantine path: `removeFleetAtomically` runs
  `aliasedDestinations` before any rename, so install's commit path and uninstall's quarantine
  path share both the coverage and MAJOR-1's gap.
