---
idea: integrate-parley-bidding-addon
review-round: 19
agent: hermes-1
date: 2026-07-31
reviewed-commit: 2b680a2
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

**Round 18 hermes-1 MAJOR (containment regression: scalar nearest-ancestor keys lose nesting when two sides anchor on different existing ancestors): CLOSED.** The `identityChain` model replaces the scalar `physicalKey` with a full per-component chain. Every component's identity is recorded root-first: existing components contribute `dev:ino`, not-yet-created tails contribute the nearest existing identity plus remaining names. Containment is `chain.includes(otherIdentity)` — an exact element match on the array, not a string prefix test. I measured the exact round-18 reproduction (outer exists, inner has its own existing parent under outer) end-to-end through the actual installer at 2b680a2: `ok: false`, both blocked, zero units written. The skewed-anchor case is caught because the outer destination's identity appears as an element in the inner destination's chain regardless of how many existing intermediates lie between them.

**Round 18 codex-1 MAJOR (firmlink respelling with an existing inner parent — the arm kimi-1's union proposal missed): CLOSED.** Measured end-to-end: installed kimi at `/private/tmp/X/KM`, then ran `--target all` with `CODEX_HOME` set to `/System/Volumes/Data/private/tmp/X/KM/skills/parley-deck` (the firmlink respelling of kimi's existing core, which has an existing `skills/` subdirectory inside it). Result: `ok: false`, blocked, zero units written. The chain model catches this because both firmlink spellings of the core directory stat to the same `dev:ino`, so both chains share that anchor element. The existing inner parent (`skills/`) gets its own `dev:ino` in the chain, but the outer destination's identity is still present as an earlier chain element — `chain.includes(outerIdentity)` is true.

**Round 18 codex-1 MAJOR (symlink buried below an existing subdirectory): CLOSED.** `entryChain` now calls `identityChain(path.dirname(entry))` which walks the full ancestry, not just the immediate parent. Measured end-to-end: a symlink at `kimiCore/one/two/redirect` (two directories deep inside kimi's managed tree) is correctly located inside kimi's destination. The plan is refused.

**Round 18 kimi-1 NIT-1 (tails lowercased but not Unicode-normalized): CLOSED.** `canonicalSegment` (line 1327) does `name.normalize("NFC")` before the case fold. I verified on this volume (APFS) that an NFD spelling and an NFC spelling of the same name produce identical `canonicalSegment` output, and that `statSync` resolves both spellings to the same `dev:ino`.

**Round 17/18 NIT (stale comment on `pathEntryExists` referencing deleted `uninstallSkillUnit`): still open.** Line 2277 of `lib/installer.js` still reads "`skillUnitStatus`, and `uninstallSkillUnit`." The function `uninstallSkillUnit` does not exist (confirmed: zero `function uninstallSkillUnit` definitions, one comment reference). `skillUnitStatus` does exist (line 1904). No behavioral impact. See NIT-1 below.

**All earlier settled matters (concurrent-installer isolation as recorded follow-up, CHANGELOG "Known limits" wording, Python fail-closed gate, uninstall dry-run parity, dead `preflightUninstallUnit`, manifest key confinement, manifest entry rule, `insidePayload` check):** verified present and unchanged. Not re-litigated.

## New findings

### NIT-1 — stale comment references deleted `uninstallSkillUnit`

`lib/installer.js` line 2277: the comment above `pathEntryExists` names `uninstallSkillUnit` as one of its callers. That function was deleted (replaced by `removeFleetAtomically`) and no definition exists in the file. `skillUnitStatus`, the other named caller, does exist. This is the same NIT carried from rounds 17 and 18. No behavioral impact — it is a comment, not code. Fix: remove "`and `uninstallSkillUnit`" from the comment, or replace with "`removeFleetAtomically`".

## Release judgement

Releasable as 2.1.0. The identity-chain model is the correct basis for the plan gate: it subsumes all three predecessor models (realpath strings, scalar dev:ino keys, and their union) and survives every adversarial arm I could construct or find in the record. The stale comment (NIT-1) is cosmetic and does not gate the release.

## What I verified

**Suite:** 361 node tests, 0 fail (npm test, this machine). Python gate refuses 3.9.6 by design (add-on declares >=3.10) — fail-closed, correct. Manifest check: 47 files, aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since 714712f. design-addons.test.js: 14 tests, 0 fail.

**All three regression arms, end-to-end through the actual installer at 2b680a2 (not logic mocks):**

1. Nesting with skewed anchors (outer exists, inner has its own existing parent): `ok: false`, 0 units written. PASS.
2. Symlink buried below an existing subdirectory (two directories deep): `ok: false`, nothing written to the symlink target. PASS.
3. Firmlink respelling with an existing inner parent (the arm kimi-1's union missed): `ok: false`, 0 units written. PASS.
4. Symlink chain whose middle hop sits inside a destination (round-17 shape): `ok: false`, dependency named in the refusal message. PASS.

**Adversarial cases from the review brief, tested directly:**

- **EACCES on a component (unreadable, not absent):** `identityChain` uses `statSync`, which succeeds on a mode-000 directory's own entry (you can stat it from its parent) but throws EACCES on its children. The child becomes a pending name-based segment, but every ancestor up to the lock point still contributes real `dev:ino` values to the chain. Containment (`chain.includes`) still works because the ancestor chain is intact. End-to-end test with a mode-000 directory inside kimi's tree: the nesting was caught by the aliased check, not deferred to preflight. Additionally, `destinationAncestorObstacle` in preflight would block any install whose parent is not writable/searchable, so EACCES cannot produce a false negative that reaches the commit phase. Siblings under a locked directory (two non-overlapping dests) correctly do NOT trigger a false positive — their pending tails differ.

- **Same dev:ino for genuinely different objects across mounts:** if two different directories report the same `dev:ino` (FUSE, NFS, bind mounts), `identityChain` treats them as the same object. This produces false positives (over-blocking: "Destination is shared" or "overlaps"), never false negatives. Over-blocking is fail-closed — safe. I could not construct a case where dev:ino collision causes an unsafe plan to be accepted.

- **`entryChain` when the parent is itself a symlink:** `identityChain` uses `statSync` (follows symlinks), so a symlink parent's identity is the target's `dev:ino`, not the symlink's own location. This is correct: the child physically exists in the target directory (the kernel resolves through the symlink), so keying by the target's ancestry locates it where it actually sits. The symlink parent itself is recorded as a separate touchpoint by `resolutionTouchpoints` (which uses `lstatSync`), so the dependency on the symlink is still captured.

- **Hardlinked directories:** same `dev:ino` → same anchor in the chain. Two hardlinks to one directory get the same identity; the equality check (`a.identity === b.identity`) catches them. No issue.

- **Unicode normalization:** `canonicalSegment` NFC-normalizes before case-folding. Verified on APFS: NFC and NFD spellings of the same name produce identical `canonicalSegment` output and `statSync` resolves both to the same `dev:ino`. For existing components, identity is `dev:ino` (name-independent). For non-existing components, identity is the NFC-normalized name. Both sides match. kimi-1's round-18 NIT-1 is closed by this.

- **Uninstall quarantine path parity:** `removeFleetAtomically` (line 1604) calls the same `aliasedDestinations` gate as `installFleetAtomically` (line 1489). End-to-end: uninstall with two nested destinations (kimi core and codex home set inside kimi's tree) returns `ok: false`, both blocked with "Destination overlaps another in this plan." The gate is identical; the quarantine rename path is covered.

- **Single-process guarantee (staging, commit, revert, quarantine, cleanup):** the identity-chain model is used ONLY in `aliasedDestinations` (the plan gate). The actual mutations — `copyPayloadAtomically`, `commitStagedUnit`, `revertStagedUnit`, `quarantineName`, `discardBackup`, `removeFleetAtomically` — all use direct path operations (`path.dirname`, `path.basename`, `fs.renameSync`, `fs.rmSync`). No identity-chain computation during mutations. The gate ensures no two units in the plan conflict before any mutation begins; the mutations themselves are atomic renames within the same parent directory. Correct design.

- **`path.join` normalization of `..` in symlink targets:** `resolutionTouchpoints` resolves relative symlink targets with `path.join(path.dirname(hop), target)`, which normalizes `..` components before the chain walk. This matches kernel resolution: a symlink target with `..` resolves to its final physical path, so no false dependencies are introduced and no real dependencies are hidden.

- **`chain.includes()` semantics:** `Array.includes()` checks exact element equality, not substring matching. `dev:ino_X/foo` is not "included" in an array containing `dev:ino_X/foobar`. I verified this directly: siblings with names that are string prefixes of each other do not trigger false positives.

**Code reviewed:** `lib/installer.js` (2536 lines — full read of identityChain, entryChain, resolutionTouchpoints, aliasedDestinations, installFleetAtomically, removeFleetAtomically, copyPayloadAtomically, commitStagedUnit, revertStagedUnit, discardBackup, quarantineName, preflightSkillUnit, destinationAncestorObstacle, pathEntryExists, installerOwnsDestination), `lib/addon-manifest.js` (full diff 49fc3ec..2b680a2), `scripts/run-python-tests.js`, `test/bidding-addon.test.js` (new tests at lines 1655-1736), `test/design-addons.test.js` (diff and test run). Diffs: 64e43f9..2b680a2 (cycle 22 only) and 49fc3ec..2b680a2 (full since round 7).

**Disk hygiene:** all scratch directories created under `/tmp` were removed with `fs.rmSync(home, {recursive: true, force: true})` in finally blocks. No repo mutations. No temporary checkouts or archives needed.
