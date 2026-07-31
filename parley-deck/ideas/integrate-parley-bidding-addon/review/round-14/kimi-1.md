---
idea: integrate-parley-bidding-addon
review-round: 14
agent: kimi-1
date: 2026-07-31
reviewed-commit: d7ab1c3
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

Three of my four round-13 findings are closed at `d7ab1c3`, each re-measured:

- **MINOR (dry-run previewed a green run the real command refused).** Closed. Dry-run now
  goes through `installFleetAtomically`, so the preflight set {recorded-selection problem,
  impossible ancestor, source-manifest mismatch, ownership} and the new alias gate are
  identical between preview and execution; only staging/commit/cleanup are omitted. The
  per-action `dryRun` flag is back in the JSON contract. I re-measured beyond the shipped
  regression: on a fresh aliased-home fixture the dry run and the real run return the same
  `ok` (both true there — see new finding 1), and the shipped test proves agreement on a
  damaged-selection tree for install *and* uninstall. Note the recorded-selection arm of my
  round-13 arm no longer blocks either mode — it is now repaired by install (next item) —
  and the two modes agree about that too.
- **MINOR (recorded-selection damage was a recovery dead-end).** Closed, with the resolution
  I preferred of the two I offered: install is exempted (its units come from discovery and
  flags, never the marker, so the block bought no path safety) *and* the message names the
  exit ("Re-run install to rewrite it."). Re-measured via the shipped regression: health
  reports the damage with the repair named, `uninstall --force` still refuses (it builds
  paths from the selection), install repairs and rewrites the marker, post-repair `doctor`
  is `valid`.
- **NIT (a recorded selection naming the core was authorized by the ownership clause).**
  Closed. `authorize` in `targetSkillUnits` rejects `parley-deck` as an add-on name with a
  dedicated problem; the shipped regression passes and the code path is as I described it.
- **NIT (install transaction cleanup paths not guarded the way uninstall's are).** **Not
  closed — carried forward, disposition unchanged.** Cycle 17 did not touch
  `commitStagedUnit` or the Phase-1/Phase-2 leftover `rmSync` loops
  (`lib/installer.js:1510-1518, 1541-1545, 1760-1774`): the intra-commit backup restore and
  the debris cleanups still run unguarded, and the composed commit-failure message still
  does not name the surviving `.bak` path the way the revert-failure message does. I still
  have no deterministic arm — every trigger needs mid-run external interference — so this
  stays a NIT follow-up, not a gate. It was not in the facilitator's round-13 fix list; I
  flag it so it is not silently dropped from the record.

codex-1's round-13 MAJORs (manifest symlink trust, aliased destinations) are closed as
measured below; the cross-process arm is ruled on in the next section.

## Ruling on concurrent-installer isolation

**Recorded follow-up. It does not gate 2.1.0.**

The finding is real — I take codex-1's measurement as given and the mechanics are plain in
the code: the stage/commit/revert transaction is per-process, and `commitStagedUnit`'s
`pathEntryExists(dest)` check and the renames around it have no mutual exclusion against a
second installer. But gating on it would be a category error, for three reasons:

1. **Reachability is operator-inflicted, transient, and self-announcing — the opposite of
   what this cycle has been fixing.** Every door closed in cycles 8–17 was reachable by a
   single actor and left a *persistent lie* in health output (a green `doctor` over a broken
   or externally-steered tree). Two interleaved installers require someone to run two
   mutations at once against overlapping skills roots; the worst end state is a wrong
   `ok:true` at race time plus possibly a missing destination or named `.bak`/`.tmp` debris
   — and the byte-verification layer this idea shipped catches the divergence on the very
   next `doctor`. Nothing about the steady state stays falsely green.
2. **The failure is self-healing by design.** Install is convergent: same package, same
   selection, same bytes, marker rewritten. "Re-run install, then run doctor" is a complete
   recovery with no data at risk — the add-on payloads are inert instruction trees, and the
   source of truth is the package, not the destination.
3. **The remedy is a design change, and round 14 of a fix-up sequence is the wrong place to
   debut one.** A lock protocol over every affected skills root, held through preflight,
   commit, rollback and cleanup, brings stale locks after crashes, lock ordering across
   roots, NFS atomicity caveats, and Windows rename semantics — each a new door of exactly
   the kind eight cycles have been closing. Shipping that unreviewed is how the next six
   rounds get spent on the lock instead. hermes-1 and I scoped this out in round 13; I hold
   that position after re-reading the transaction code at `d7ab1c3`.

For the follow-up design, the minimum acceptable mechanism when it is taken: a per-skills-root
lock directory acquired by atomic `mkdir` (no lockfile content races), holding an owner record
(pid, started-at, command), stale-lock detection via process liveness plus a generous max-age,
acquisition across multiple roots in a canonical sorted order, refusal rather than waiting on
contention, and release in `finally`.

The release notes should say, plainly: **do not run two `parley-deck-skill` install or
uninstall commands concurrently against the same runtime skills directories.** The commands
are atomic per invocation but not isolated across processes. If a concurrent run may have
happened, re-run `install` (it converges) and confirm with `doctor`.

## New findings

Two, both **recorded follow-ups; neither gates 2.1.0.**

### MINOR — the alias check misses two routes to one physical destination

`physicalKey` (`lib/installer.js:1415`) resolves `realpath` of the parent and appends the
basename. Two measured gaps at `d7ab1c3`:

- **Missing shared container.** When the shared skills container does not exist yet,
  `realpathSync` throws and the key falls back to lexical `path.resolve`, which cannot see a
  symlinked ancestor. Measured: `CODEX_HOME=R`, `KIMI_CODE_HOME=L` with `L` a symlink to
  `R`, neither `skills/` existing → `install --target all` returns **ok:true**, codex's six
  units `installed`, then kimi's six `replaced` them in the same physical directory; the
  surviving marker records `target: "kimi"`; post-hoc `doctor` is green for both runtimes;
  no debris. (Control: with `R/skills` pre-existing the alias is caught and the whole plan
  refused — the shipped regression's case.)
- **Case-insensitive filesystem.** This volume is case-insensitive APFS, and
  `fs.realpathSync` (the non-native implementation) preserves the caller's component case —
  measured: `realpathSync('/tmp/pd-exp/caseprobe')` → `/private/tmp/pd-exp/caseprobe` while
  `realpathSync.native` → `/private/tmp/pd-exp/CaseProbe`. So `CODEX_HOME=/…/CaseHome` and
  `KIMI_CODE_HOME=/…/casehome` — one physical directory — produce different keys even with
  the container present. Measured: **ok:true**, kimi's units again `replaced` codex's.

Reachability bound, checked before rating: only the env-var homes (codex, kimi) or
project-scope symlinked roots can spell one directory two ways, and only **identical-payload
target pairs** can collide through these gaps — gemini's and agy's specializations differ in
the final container component (`extensions` vs `config/plugins`), so those two alias only
when the containers already exist (caught) or dangle (blocked by
`destinationAncestorObstacle`). No cross-shape overwrite — the round-13 measured harm — is
reachable here; the residual is a success report for a configuration the tool's own doctrine
says it must refuse, plus a misattributed `marker.target`. Fix shape: resolve the nearest
*existing* ancestor with `realpathSync.native` and append the unresolved tail lexically —
the same walk `destinationAncestorObstacle` already performs. Follow-up, non-gating.

### NIT — `runtimeAvailability` reads the manifest without the regular-file predicate

Cycle 17's claim is that one predicate — regular, non-symlink file — now applies in
`hasManifest`, `manifestFileHash` and `verifyPayload`. There is a fourth reader:
`runtimeAvailability` calls `readManifest` directly (`lib/installer.js:2026`) and follows a
symlink. On the managed path that is unreachable — health only probes the runtime after the
guarded predicates have vouched for the tree — but the **legacy-marker path** (the genuine
2.0.0 shape: no `markerSchema`, no `manifest` field, for a skill whose source ships no
manifest) skips every manifest check. Measured at `d7ab1c3`: parley-tracker with a
legacy-shaped marker and `parley-addon.json` replaced by a symlink to an external file
declaring `runtime.python: ">=99.9"` → unit `valid`, `managed: true`, zero problems, and
`runtime: {ok:false, requirement:">=99.9"}` — `doctor` fails on the authority of bytes
outside the installed tree, the class cycle 17 closed everywhere else. It cannot launder a
broken payload (payload checks never read it); it only steers the availability verdict.
One-line fix: `return null` unless `hasManifest(root)`. Follow-up, non-gating.

For completeness, the other routes the prompt names were checked and hold: directory
hardlinks are not creatable on APFS/ext4 (and a hardlinked *file* manifest is bounded by the
marker content hash), and `..` inside a symlinked parent is resolved by `realpath` whenever
the parent exists — the missing-parent case is the finding above, not a separate door.

## Release judgement

**Releasable as 2.1.0.** Cycle 17's claims verified as stated (below); the round-13 severity
dispute is moot because both disputed items were fixed and both fixes re-measure clean. Three
follow-ups are recorded — the alias-key gaps (MINOR), the legacy-path runtime read (NIT), and
my carried cleanup-guard NIT — plus the concurrency follow-up with release-note wording per
the ruling above. None is the one thing that must change, because none can produce a
persistent false health verdict or a cross-shape overwrite.

## What I verified

Environment: the reviewed tree itself at `d7ab1c3` (`git status` clean before and after; the
tree was never mutated — no resets, no edits, no scratch files in the repo), node v26.5.0,
python3 3.14.6 (homebrew) and 3.9.6 (`/usr/bin/python3`). All fixtures ran under `/tmp/pd-exp`.

- **Baseline claims, all reproduced.** `node --test`: **349 tests, 0 fail**. Python leg:
  **54/54 on 3.14**; with `PATH` putting `/usr/bin/python3` first it refuses by design
  ("python3 is 3.9, but the add-on declares >=3.10", exit 1). Manifest check ok — 47 files,
  aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`,
  matching the facilitator. Only `parley-bidding` ships a manifest.
- **Cycle-17 diffs read in full** (`dd8d756..d7ab1c3`, and `49fc3ec..d7ab1c3` for full
  scope: lib, tests, CHANGELOG). CHANGELOG claims match the code I read.
- **Round-13 manifest repro, re-run**: installed manifest moved out and replaced by a
  symlink to the byte-identical file → `verifyPayload ok:false` ("symbolic link, not a
  regular file"), `doctor` `malformed`, `managed: false`, `doctor.ok false`.
- **Predicate coverage, read line by line**: every runtime read of `parley-addon.json`
  enumerated — `hasManifest` (lstat), `manifestFileHash` (guarded), `verifyPayload`
  (guarded), `manifestProblems` presence/hash arms (via the guarded three), `writeMarker`
  (behind `hasManifest` + preflight-verified source), and `runtimeAvailability` (unguarded —
  new NIT above). Hardlink and `..` routes argued and bounded as recorded above.
- **Aliasing**: the shipped symlinked-container regression passes and I re-measured its
  blocked message and zero-writes assertion; then the two new gap arms (missing container,
  case-variant) measured as described in the finding, including post-hoc health and debris
  checks.
- **Dry-run parity**: shipped regression (install + uninstall agreement on a damaged
  selection) passes; fresh-fixture dry vs real agreement measured independently; the dry-run
  result map carries the per-action `dryRun` flag; uninstall dry-run applies the marker,
  ownership and alias gates (same function as the real run).
- **Damaged-selection repair and core-as-addon refusal**: shipped regressions pass; the
  `targetSkillUnits`/`authorize` code paths confirm install derives units from discovery and
  flags only, so self-healing opens no new path from stored data.
- **Full-scope re-read of `lib/installer.js` (2506 lines) and `lib/addon-manifest.js`
  (301 lines) at `d7ab1c3`**, including the unchanged transaction machinery my round-13 NIT
  concerns — confirmed still present, hence carried rather than dropped.
