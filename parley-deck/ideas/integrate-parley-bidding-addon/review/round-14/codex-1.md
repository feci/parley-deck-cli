---
idea: integrate-parley-bidding-addon
review-round: 14
agent: codex-1
date: 2026-07-31
reviewed-commit: d7ab1c3
---

## Verdict

BLOCK

Cycle 17 closes the measured `doctor`/`managed` manifest false green, makes a damaged recorded
selection self-healing through install, and gives install dry-run the real read-only plan. It
does not yet close the physical-destination or dry-run findings in full. The alias key has two
reproduced false negatives, including one on this Mac's case-insensitive filesystem, and
uninstall dry-run still describes removals the real fleet gate will skip.

## Outstanding findings — closed or not

### [MAJOR] Round-13 physical-destination uniqueness: NOT CLOSED

`physicalKey` tries `realpathSync(dirname(dest))` and falls back to the unresolved logical
parent when that call fails (`lib/installer.js:1415-1422`). That is not a physical key in two
important cases:

1. **A symlinked runtime root whose `skills/` directory does not exist yet.** I symlinked both
   `.codex` and `.hermes` to one empty physical runtime root, then ran the real install with
   `--target all --include-undetected`. Both `realpathSync(.../skills)` calls failed before
   staging, so the fallback keys differed. Install returned `ok: true`; Codex reported six
   `installed` units, Hermes reported the same six as `replaced`, and the one physical core's
   final marker said `target: hermes`. This is a natural first-install form of the configuration
   cycle 17 intends to refuse. A later one-target uninstall can remove the shared install from
   the other runtime.
2. **Case-equivalent paths on a case-insensitive filesystem.** On this case-insensitive APFS
   volume, `realpathSync` preserves the caller's spelling rather than canonicalizing case. I set
   `CODEX_HOME` to an existing `RuntimeHome` and `KIMI_CODE_HOME` to the case-only spelling
   `runtimehome`. Both resolved to the same physical `skills/` directory, but their string keys
   differed. The real fleet install again returned `ok: true`; Codex reported `installed`, Kimi
   `replaced`, and the physical core's final marker said `target: kimi`.

The current regression pre-creates the exact aliased parent, so it exercises only the case where
`realpathSync(parent)` succeeds and returns identical strings. The minimum release fix is to
identify the nearest existing ancestor by filesystem identity (for example `dev`/`ino`) and
append the unresolved tail, with case-equivalence handled for case-insensitive volumes, then
reject every repeated physical unit key before either dry-run reporting or staging. Regress both
an empty pair of symlinked runtime roots and case-only home overrides. Directory hardlinks are
not a normal supported filesystem operation; existing-parent inode identity also covers the
relevant hardlink/bind-alias class better than a realpath string.

### [MINOR] Round-13 manifest predicate: operational arm CLOSED, `readManifest` arm NOT CLOSED

The installed-health result is fixed: `hasManifest`, `manifestFileHash`, and `verifyPayload`
now reject a symlink at `parley-addon.json`; the original external-manifest reproduction reports
`malformed` and no longer confers `managed: true`.

But the exported parser still calls `readFileSync` directly without the new entry predicate
(`lib/addon-manifest.js:124-130`). With a symlink to a byte-identical external manifest I
measured:

- `hasManifest: false`
- `manifestFileHash: null`
- `verifyPayload.ok: false`
- **`readManifest.ok: true`**

That is the exact disagreement my round-13 fix request included. Current stable health happens
to validate before `runtimeAvailability` calls the parser, so I rate the residual MINOR rather
than the prior MAJOR. Make `readManifest` enforce `manifestEntryProblem` itself so every read,
including future callers and the current check-then-read call sites, shares the same rule.

### [MINOR] Round-13 dry-run fidelity: install CLOSED, uninstall NOT CLOSED

Install dry-run now runs `installFleetAtomically`'s read-only planning and preserves the
per-action `dryRun` flag. Its preflight result matches the real install in the cases I checked.

Uninstall still has two different fleet gates. The real command's early preflight is guarded by
`!dryRun` (`lib/installer.js:702-734`) and marks every otherwise-removable unit `skipped` if any
unit blocks. The dry path enters `removeFleetAtomically`, records each good unit as `remove`
immediately, and discovers `blockedAnywhere` only afterwards (`lib/installer.js:1580-1629`);
those already-recorded results are never changed.

I installed all six Codex units, changed only `parley-bidding`'s marker to belong to a foreign
installer, and compared dry and real uninstall. Both top-level results were `ok: false`, but the
dry run reported core plus four other add-ons as `remove`/`ok: true`, while the identical real
command reported all five as `skipped`/`ok: false`. The cycle-17 regression asserts only the two
top-level `ok` booleans and uses a damaged selection that collapses the plan to the blocked core,
so it cannot see this mismatch. Build all dry-run candidates first, apply the same fleet-wide
block, and only then label otherwise-removable units.

### Round-13 damaged recorded selection: CLOSED

Install units are derived from discovery and explicit flags, not from the damaged marker. The
real install and install dry-run therefore proceed safely, the real install rewrites the
selection, health names the damage and repair, and uninstall continues to refuse the marker it
would use to build paths. I found no new stored-data-to-path route in this self-healing branch.

### Round-13 recorded selection naming the core: CLOSED

Read and uninstall paths now reject `parley-deck` as an add-on name before the ownership clause
can authorize the duplicate unit. The regression passes.

### Earlier findings and recorded follow-ups

All findings closed through round 12 remain closed under the attacks repeated here. The prior
recorded follow-ups remain non-blocking on their existing dispositions: manifests for the
remaining units/B3.11, the `dirExists` discovery NIT, quarantine debris not visible to
`doctor`, and residual-disposal failures that leave named debris after the destination has been
successfully replaced or removed.

## Ruling on concurrent-installer isolation

**Recorded follow-up; cross-process locking does not gate 2.1.0.**

The reproduced race is real: without exclusive ownership, one process's rollback can move a
second process's committed directory aside after that second process has logically succeeded.
But the ratified design requires an atomic selected fleet for one invocation; it does not claim
multi-process serializability. A portable lock protocol across several skills roots, including
crash recovery, stale ownership, and network-filesystem semantics, is a new subsystem rather
than the minimum correction to cycle 17's in-process alias check. Adding it in this fix-up would
expand the design and its failure surface late in the release.

The 2.1.0 release notes must state the boundary plainly, adjacent to the fleet-atomicity claim:

> Installer mutations are single-writer in 2.1.0. Do not run two `install`/`uninstall`
> commands, or another skill manager targeting any of the same skills roots, at the same time.
> Wait for one command to finish before starting the next. Concurrent processes are not
> isolated; an overlapping rollback can invalidate a command that already reported success.
> After any suspected overlap, serialize further commands, run `doctor`, and reinstall the
> intended selection.

The follow-up should design deterministic multi-root locking held from before preflight through
commit, rollback and cleanup, with an explicit crash/stale-lock recovery contract. That design
work is not a condition of this release once the warning above ships.

## New findings

None beyond the incomplete closures above.

## Release judgement

Not releasable as 2.1.0 at `d7ab1c3`. The release-blocking code change is to make the
same-process physical-destination set genuinely unique when parents are absent and when path
spellings are case-equivalent; the current implementation still returns a false success for one
physical fleet represented twice. Because this idea has `strict_gate: true`, the two residual
MINOR discrepancies in `readManifest` and uninstall dry-run must also be fixed or explicitly
closed by the operator, and the single-writer warning above must be added to the release notes.

## What I verified

- Read the live cooperation protocol, `00-prompt.md`, `FINAL.md`, consensus, the implementation
  record through cycle 17, all round-13 reviews, the complete current `lib/installer.js` and
  `lib/addon-manifest.js`, the cycle-17 diff `dd8d756..d7ab1c3`, and the full round-7 delta
  `49fc3ec..d7ab1c3`. Both requested ranges pass `git diff --check`.
- Ran the full Node suite at `d7ab1c3`: **349 tests, 349 pass, 0 fail** with Python 3.14 first
  on `PATH`, and again with `/usr/bin/python3` 3.9.6 first on an otherwise complete `PATH`.
- Ran the Python leg under its cache-safe runner: **54/54** on Python 3.14 across the seven
  declared files. The manifest check reports **47 files** and aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- Ran `npm pack --dry-run --json` with an isolated temporary npm cache: **202 entries**, **48**
  under `skills/parley-bidding/`, and no `__pycache__`, `.pyc`, or `.pyo`. Parsed all **16**
  bidding JSON files and ran the adapter validator successfully against all four adapter files.
- Confirmed `skills/parley-bidding/` is unchanged from `714712f`, contains no symlink or Python
  cache artifact, and the five cycle-17 regressions pass as part of the full suite.
- Reproduced the two physical-key false negatives, the uninstall per-unit dry-run discrepancy,
  and the direct `readManifest` symlink acceptance in isolated temporary homes. The reviewed
  repository remained clean at `d7ab1c3`; no file under `skills/parley-bidding/` was modified.
- For the routes named in the brief: `path.join`/`path.resolve` remove lexical `..` before a unit
  is planned; ordinary directory hardlinks are unavailable here and hardlinked regular
  destination entries do not alias the directory created by rename. The reproduced missing-
  parent and case-insensitive routes are sufficient to disprove the current `realpath` key even
  without those two arms.
