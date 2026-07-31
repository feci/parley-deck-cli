---
idea: integrate-parley-bidding-addon
review-round: 15
agent: kimi-1
date: 2026-07-31
reviewed-commit: 26478e9
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

All four round-14 items are closed at `26478e9`; each was re-measured, not trusted.

1. **Destination identity (`physicalKey`) — closed.** I attacked the dev/ino-of-nearest-existing-ancestor
   construction with six shapes (probe details below): symlinked env-homes to one empty root (the round-14
   false negative, re-run through the `CODEX_HOME`/`KIMI_CODE_HOME` path rather than the `.codex`/`.hermes`
   path the regression uses), case-only spellings of an existing home, case-only spellings of a home that
   does not exist anywhere yet, a lexical `..` composed with a symlinked tail, a dangling symlink home vs
   the direct not-yet-existing path, and `skills/` itself as a symlink to a shared container. The first
   four and the sixth are caught with the "resolve to the same directory" refusal and zero writes. The
   dangling-symlink case is *not* recognized as an alias — the two spellings correctly produce different
   keys, because nothing physical exists to share — and it dies fail-closed in
   `destinationAncestorObstacle` ("broken link or unreadable") with the fleet gate skipping the rest.
   Nothing written in any shape. Within the single-process model I found no remaining route: staging only
   ever `mkdir`s real directories, which cannot manufacture an alias the planning-time walk did not see,
   and every `..` is consumed lexically by `path.resolve` before the walk.
2. **Manifest rule on every read — closed.** `readManifest` now applies `manifestEntryProblem` before
   parsing, so all four runtime readers (`hasManifest`, `manifestFileHash`, `verifyPayload`,
   `readManifest`) share the lstat predicate; the cycle-18 regression asserting all four refuse a symlink
   is in the suite and passes. I audited for a fourth/fifth consumer: no read of `parley-addon.json` in
   `lib/` or `bin/` bypasses the four. The only direct read anywhere is
   `scripts/run-python-tests.js:49`, which `JSON.parse`s the *packaged source* manifest to pick an
   interpreter — release-time, trusted origin, fails loud on garbage, and the release gate
   (`build-addon-manifest --check`) goes through `readManifest`. Not a trust-boundary gap.
3. **Uninstall dry-run promised removals — closed for the round-14 shape.** The gate now runs before any
   unit is recorded. The added regression passes, and my own clean-case and alias-blocked parity probes
   (O3, O4b) show dry and real agreeing per-unit, verb-normalized. See the new finding for the one
   residual divergence, which is a different mechanism.
4. **Redundant install preflight (hermes-1 NIT-2) — closed on the install side.** The block is gone from
   `installCommand`; install dry and real now share `installFleetAtomically` end to end, and my probe M
   (one foreign-marked target in a 14-target plan) shows identical per-unit records between dry and real,
   modulo the intended `install`/`installed` verbs.

The scope question stays settled. The "Known limits" note in `CHANGELOG.md` matches the adopted wording
verbatim, and I found no evidence the *single-process* guarantee is broken (transaction audit below).

## New findings

**MINOR — uninstall kept the preflight whose install twin cycle 18 deleted, and it still splits dry-run
from real on one unit label.** `uninstallCommand` (lib/installer.js:664-696) still runs its separate
fleet preflight with its own early-return builder. That builder flattens *every* non-blocked unit to
`{ok:false, action:"skipped", message:"Not attempted: another skill or target in this uninstall failed
preflight."}` — including units whose destination simply does not exist. Dry-run skips that block and goes
straight to `removeFleetAtomically`, which records a missing unit as `{ok:true, action:"missing"}`.
Measured at `26478e9`, same state twice: with `parley-bidding` foreign-marked (blocked) and
`parley-tracker` deleted (missing), dry reports tracker `ok:true/"missing"` while the real command
reports the same unit `ok:false/"skipped"` with the preflight message; same shape when the two units are
on different targets. The preflight is now a strict subset of `removeFleetAtomically`'s own checks
(marker problem, existence, ownership — identical strings; aliasing lives only in
`removeFleetAtomically`), so it buys nothing and costs exactly the divergence cycle 18 removed from
install for exactly this reason.

Why this does **not** gate 2.1.0: top-level `ok` and the exit code agree in every measured shape; the
blocked unit's own record is identical between dry and real; no mutation is promised that will not happen
(the divergent unit is a no-op in both runs); and the direction is never phantom success — measured, dry
is never stricter than real per unit. The round-14 blocker was dry *promising five deletions* the real
command refused; this is a label disagreement about a directory that is not there. Recorded follow-up:
delete the redundant preflight block (or build its early return from `removeFleetAtomically`'s
dispositions), mirroring the cycle-18 install change — which also fixes the real command's mildly
misleading attribution, since a missing unit is not "not attempted because of the preflight", it is just
absent.

## Release judgement

Releasable as 2.1.0. The one new finding is a recorded follow-up, stated above with its mechanism and its
measured blast radius; it does not touch the exit-code contract, the mutation set, or the single-process
transaction guarantee.

## What I verified

- Full-suite re-run at `26478e9`: **353 node tests, 0 fail**; python leg 54/54 on 3.14; manifest check ok
  (47 files, aggregate `sha256:7854adf1…b95a6d`, unchanged since `714712f`). Working tree untouched
  (`git status` clean before and after; all probes wrote only to `os.tmpdir()`).
- `physicalKey`, six empirical attack shapes through `install --target all --include-undetected` with
  `CODEX_HOME`/`KIMI_CODE_HOME` as the aliasing vehicle: symlinked homes to an empty shared root
  (caught), case-only existing homes (caught on this case-insensitive volume), case-only not-yet-existing
  homes (caught), `sub/../alias` lexical-`..` composition (caught), `skills/` as a symlink to a shared
  container (caught), dangling-symlink home vs direct missing path (not an alias — correctly — and
  blocked fail-closed by the ancestor walk; zero writes). Reasoned but not measurable on this machine,
  and left as theoretical residuals rather than findings: two mounts of one physical share with distinct
  `st_dev` (bindfs/FUSE, dual-mounted NFS, Windows UNC-vs-drive-letter) would still produce two keys for
  one directory; exploiting that needs two runtime homes on two different mounts of the same share.
- Manifest predicate: source audit of every `parley-addon.json` consumer across `lib/`, `bin/`,
  `scripts/` (results above); the cycle-18 four-reader regression passes.
- Dry-run parity, measured per-unit on identical states: install with one blocked target (identical
  records, verb-normalized); clean uninstall (identical); alias-blocked uninstall (identical after
  normalizing the dev/ino the alias message prints); preflight-blocked uninstall with a missing unit
  (the one divergence, reported above); reverse direction checked (dry never reports `ok:false` where
  real reports `ok:true`).
- Single-process transaction guarantee, measured: a `uchg`-frozen destination makes the fleet commit
  fail mid-phase-2 — every earlier commit reverted, both pre-existing installs byte-valid afterwards
  (`doctor` green per target, payload *and* runtime arms), zero `.tmp`/`.bak` debris. A read-only skills
  parent makes quarantine fail mid-phase-A — full rollback, nothing deleted, no `.removing` debris.
  Backup-cleanup and quarantine-cleanup failures remain warnings, not failures, by construction; the
  commit-revert double-failure path requires the parent to change writability mid-sequence, which the
  single-writer model excludes.
- Round-15 scope discipline: I did not re-open the concurrency ruling, and I did not read codex-1's
  round-15 file before writing this one.
