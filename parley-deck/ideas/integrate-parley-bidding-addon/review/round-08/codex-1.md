---
idea: integrate-parley-bidding-addon
review-round: 8
agent: codex-1
date: 2026-07-30
reviewed-commit: ebe269e
---

## Verdict
BLOCK

## Outstanding findings — closed or not

The three findings carried into cycles 8 and 9 are closed in the cases they named:

- The source-side B5 hole is closed for a symlink in a manifest-free add-on and for a symlink
  in the core payload. `copySourcesFor` covers the add-on root and all core payload entries,
  `firstCopyObstacle` accepts file roots, and both regressions leave the destination absent.
- D-2 now says only what the file-count and cache checks establish; it no longer calls them a
  byte-level proof that the source was untouched.
- D-3 now documents the trailing continuation-sentinel exception, and the grammar regression
  accepts that case while still refusing an embedded backslash.

The earlier health, ownership, selection, interpreter, and manifest findings remain closed
under the full suite and focused probes. The five manifest-free skills under the universal
installer remain the recorded B3.11 follow-up, not a new 2.1.0 blocker.

B5 as a whole is not closed, however: destination preflight still has two paths to a partial
fleet, described below.

## New findings

### [MAJOR] Destination preflight is neither fleet-wide nor lstat-aware, so known failures still write a partial fleet

**Where:** `lib/installer.js:606-624`, `:958-969`, `:1079-1110`, `:1126-1143`, and
`:1257-1322`.

**What:** `installCommand` calls `installTarget` sequentially. Each target is preflighted only
immediately before that target is written, so a predictable blocker in a later target is
discovered after earlier targets have already been installed. Within one target,
`preflightSkillUnit`, `installSkillUnit`, and the backup/replace path use `fs.existsSync` for
the destination. `existsSync` follows symlinks and returns false for a dangling link, so that
visible filesystem entry bypasses ownership preflight; the final rename then fails with
`ENOTDIR` after earlier units were written.

**Why it matters:** FINAL.md B5 requires every unit and destination to be preflighted before
the first write and requires zero writes for a predictable failure. Both cases below are
statically discoverable before mutation, yet both leave a partial installation. This is the
same acceptance blocker as the source-side partial-fleet defect, reached through destinations
instead.

**Evidence:** I reproduced both cases through `installer.installCommand` in an isolated copy
of `ebe269e`:

1. `--target all --include-undetected`, with an unmarked `parley-bidding` destination in the
   final `aionrs` target, returned `ok:false`; `aionrs` was refused, but all thirteen preceding
   targets from `codex` through `opencode` had already been written.
2. A default Codex install with a dangling destination symlink at the final
   `parley-worktrees` unit returned `ok:false` with `ENOTDIR`; `parley-deck`,
   `parley-bidding`, `parley-design`, `parley-design-check`, and `parley-tracker` were already
   installed with markers, while the dangling link remained.

**Fix:** Build the complete target/unit plan first and preflight every source and destination
across that plan before invoking any write. Treat only `lstat` `ENOENT` as absence, use the
same path-entry predicate in preflight and replacement, and define safe `--force` handling for
symlinks and other non-directory entries. Add regressions for a blocker in the last target and
a dangling destination at the last unit; both must assert zero installed markers/directories.
When recording that fix-up, also refresh `IMPLEMENTATION.md`'s top-level `head-commit`, which
currently still says `3634cc8` despite `status: fix-up-cycle-9`.

## Release judgement

No. Commit `ebe269e` is not releasable as 2.1.0. B5 destination preflight must become
fleet-wide and lstat-aware so both late-blocker probes produce zero writes; the ensuing
fix-up must receive another fresh full-scope strict-gate review.

## What I verified

- The reviewed repository was clean at `ebe269e15a79d0a992d51fe568e8a9cdb895fb55`
  before and after review. All execution happened in an isolated clone; the source workspace
  was read only and no Python was run there.
- I read the complete implementation diff from `a544dcd` through `ebe269e` and the complete
  cycle-8/9 diff from `49fc3ec` through `ebe269e`, including installer, manifest, test runner,
  packaging, workflows, documentation, and the bidding payload/instructions.
- `npm test` passed with **316 Node tests and 54 Python tests** on Python 3.14. The Python leg
  separately passed **54/54** on 3.10 and 3.11. The Node suite passed with system Python 3.9.6,
  while the Python leg refused 3.9.6 with the declared `>=3.10` message.
- `npm pack --dry-run --json` passed after using an isolated writable npm cache: **202 files**,
  **48** under `skills/parley-bidding/`, with prepack manifest verification and no Python cache
  artifacts.
- The current portable binary built, installed `--only parley-bidding` into an isolated home,
  and `doctor` reported both core and bidding payloads valid; the bidding marker carried
  aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- The adapter validator accepted all four shipped adapters. The Python suite covered all 16
  shipped JSON documents, the seven deterministic tools, lifecycle/recovery behavior, and the
  54-test per-file count contract.
- Source-vs-integrated comparison showed the documented shape: source `.gitignore` dropped,
  `parley-addon.json` added, and nine content-different files accounting for the rename,
  schema IDs, and consent paragraph. The integrated payload has 48 files, no symlink, and no
  `__pycache__`, `.pyc`, or `.pyo`.
- Production Python scripts contain no network/browser/process-launch dependency and the
  instruction text preserves the human approval ladder, exact-byte evidence rules,
  upload/submission separation, ambiguous-result no-retry rule, adapter ceilings, and the
  prohibition on deterministic portal mutation.
- Focused probes reconfirmed marker/payload integrity, `valid-unmanaged` anchoring, interpreter
  availability separation, source-side core/add-on copy-obstacle preflight, and the two
  destination/fleet failures reported above.
