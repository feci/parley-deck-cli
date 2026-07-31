---
idea: integrate-parley-bidding-addon
review-round: 17
agent: codex-1
date: 2026-07-31
reviewed-commit: fa3da41
---

## Verdict

BLOCK

## Outstanding findings — closed or not

The exact round-16 findings are closed at `fa3da41`.

- The ordinary-symlink dependency is refused before either mutation. I reproduced the round-16
  install shape and independently built the corresponding uninstall state: both units were
  `blocked`, both managed destinations remained present, and uninstall created no quarantine
  residue. The check is shared by `installFleetAtomically` and `removeFleetAtomically`.
- The Python runner now fails closed on both malformed interpreter output and a missing or
  malformed `runtime.python` floor. The shipped floor still runs 54/54 tests under Python 3.14
  and refuses Python 3.9.6 by design.
- `preflightUninstallUnit` is gone, and `uninstall --dry-run` now carries `dryRun: true` on the
  per-action object as well as at the top level and per skill.
- The round-15 containment, uninstall single-result-path, manifest-reader, and case-only-test
  findings remain closed.

The broader single-process destination guarantee is not closed. The new MAJOR uses an APFS
firmlink already provided by macOS, not a concurrent writer, and recreates the same false
success the ordinary-symlink fix was intended to prevent. The round-14 concurrent-installer
follow-up remains settled and is unrelated.

## New findings

### [MAJOR] Firmlink aliases bypass the symlink-only dependency check

`resolutionTouchpoints` records only components for which `lstatSync(...).isSymbolicLink()` is
true (`lib/installer.js:1382-1408`). That is not the complete set of namespace redirects on a
supported host. On macOS, APFS firmlinks expose paths such as `/private/...` and
`/System/Volumes/Data/private/...` as two directory spellings of the same objects: `lstat`
reports directories, `stat` gives matching device/inode identities, and `realpath` preserves
the two spellings.

The planner already computes an identity-plus-tail key that can see the shared physical
ancestor (`physicalKey`, lines 1323-1346), but it groups those keys by exact equality only
(`aliasedDestinations`, lines 1417-1438). The later containment pass compares the two preserved
`realpath` strings (lines 1447-1493), while the new touchpoint pass sees no symlink. A physical
ancestor/descendant pair reached through the two firmlink spellings therefore passes every
check.

I reproduced this at `fa3da41` with no `--force`:

- Kimi's home was a directory below `/private/tmp`; its planned core was the outer destination.
- `CODEX_HOME` used `/System/Volumes/Data/private/tmp/...` and pointed at that planned Kimi core,
  making Codex's planned core physically nested inside it but textually unrelated.
- `install --target all --include-undetected --no-addons` returned top-level `ok: true`;
  Codex reported `installed` and Kimi reported `replaced`.
- After the command, Codex's reported destination was absent and Kimi's surviving marker named
  `target: kimi`.

The mechanism is deterministic: Codex staging materializes Kimi's destination, Codex commits
inside it first, and Kimi's later commit renames and ultimately deletes that whole ancestor.
This is the cycle-19 false-success/partial-fleet state through a non-symlink alias, on the
default macOS filesystem topology. It directly disproves the cycle-20 comment that a symlink is
the only component able to redirect resolution out of the lexical tree.

Fix the planner by comparing physical ancestor chains, not only equal endpoint keys plus
symlink entries. A minimum regression should exercise both `/private` and
`/System/Volumes/Data/private` spellings in both nesting orders and require a preflight refusal
with zero writes. One viable implementation is to retain structured device/inode anchors and
tail segments for each destination and reject equal or prefix-related tails under the same
anchor; a complete implementation should use component identity chains so it also works when
the nearest existing ancestor differs between the two paths. Keep the ordinary-symlink
touchpoint test as a separate regression.

## Release judgement

Not releasable as 2.1.0 at `fa3da41`. The one thing that must change is destination dependency
detection: it must reject physical ancestor/descendant plans reached through APFS firmlink
aliases before staging, so a single invocation cannot report a unit installed and then erase
it with a later commit.

## What I verified

- Read the live cooperation protocol, `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through
  cycle 20, `review/round-09/VOID.md`, all three round-16 reviews, and both requested diff
  scopes. `git diff --check` is clean for `a49d68f..fa3da41` and `49fc3ec..fa3da41`.
- Audited planning, preflight, staging, commit, reverse rollback, backup cleanup, uninstall
  quarantine, quarantine rollback, and cleanup in `lib/installer.js`. The ordinary-symlink
  dependency check reaches uninstall as well as install; I reproduced uninstall with both
  managed destinations and confirmed zero renames, deletions, or residue.
- Attacked the narrowing specifically. Hard-linked regular files do not redirect directory
  resolution, and ordinary userland directory hard links are not constructible here. A
  symlink whose target lies in a later destination is caught by final containment; a symlink
  entry stored in that destination is caught by touchpoints. This process cannot make a path
  component become a symlink after planning because source symlinks are refused and staging
  creates only real directories/files before rename. A nested mount point remains lexically
  contained, but I did not claim unmeasured dual-mount behavior as a separate finding. The
  measured APFS firmlink arm is sufficient to show that “only symlinks” is incomplete.
- Ran `npm test` from a `git archive` of `fa3da41` using the repository's existing read-only
  dependencies: **357/357 Node tests**, **54/54 Python tests** on Python 3.14.6, and the
  47-file manifest check with aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
  Ran the Node suite again with `/usr/bin/python3` 3.9.6 first on `PATH`: **357/357**.
- Ran `npm pack --dry-run`; its prepack manifest check fired. Parsed an `--ignore-scripts`
  inventory separately: **202 files**, **48** under `skills/parley-bidding/`, with no
  `__pycache__`, `.pyc`, or `.pyo`.
- Validated all four adapters directly with
  `PYTHONDONTWRITEBYTECODE=1 python3 -B`, parsed all **16** bidding JSON files, and confirmed no
  payload symlinks or Python cache artifacts. The integrated payload is unchanged from
  `714712f`; source comparison remains 48 files versus 48, with the recorded dropped
  `.gitignore`, added manifest, and nine changed integration files.
- All temporary archives and homes were removed. The reviewed repository remained clean, and
  no file under `skills/parley-bidding/` was changed.
