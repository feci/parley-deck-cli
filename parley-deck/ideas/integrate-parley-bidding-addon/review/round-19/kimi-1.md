---
idea: integrate-parley-bidding-addon
review-round: 19
agent: kimi-1
date: 2026-07-31
reviewed-commit: 2b680a2
---

## Verdict
BLOCK

## Outstanding findings — closed or not

All three round-18 regressions are closed at `2b680a2`, measured by me this round, not trusted
from the record:

- **Nesting with skewed anchors** — closed. The new test
  (`test/bidding-addon.test.js:1655`) pins it and passes; I additionally re-measured the shape
  function-level in both orders.
- **Symlink buried below a subdirectory** — closed. Pinned at
  `test/bidding-addon.test.js:1684`, passes.
- **Firmlink respelling** — closed, and I went past the pinned arm (see "What I verified",
  arm 1): the existing test leaves the inner tail nonexistent, so I measured the harder shape —
  respelled outer with the inner side's parent existing under both spellings — in all three
  spelling/order combinations. The chain blocks every one (`blocked=2` each).

The cycle-22 choice of the chain over my round-18 union proposal is justified, and I withdraw
the proposal rather than argue it: my arm-1 measurement is exactly the arm the union loses, and
the chain covers it. The reasoning in `IMPLEMENTATION.md` survives contact with the code.

## New findings

### MAJOR — the identity-chain gate is void on win32 (cycle-22 regression)

`identityChain` (`lib/installer.js:1332`) and `resolutionTouchpoints` (`lib/installer.js:1369`)
rebuild each component path as `path.join(path.sep, part)` starting from `path.sep`. That is
correct only because on POSIX `path.parse(resolved).root === path.sep`. On win32,
`path.resolve("C:\\a\\b")` splits into `["C:", "a", "b"]` and `path.win32.join("\\", "C:")`
produces `"\\C:"` — a UNC-shaped non-path. Every `statSync`/`lstatSync` the gate attempts
throws, so:

- no chain ever carries a `dev:ino` anchor — identity degrades to case-folded text;
- `resolutionTouchpoints` records **zero** touchpoints — the pass-through check (rounds 16-18)
  is absent;
- junction-aliased destinations (the common no-admin aliasing mechanism on Windows, and round
  13's exact case) compare as different text and are **not blocked**;
- only textual equality/prefix survives — below even cycle 19's realpath model, which on win32
  resolves junctions.

Measured, not conjectured: I ran the gate functions extracted verbatim from `2b680a2` with
`path.win32` injected and a fake fs that answers stat for real win32 paths (`C:\Users\...`).
The gate queried only `stat:\C:`, `stat:\C:\Users`, … — never one valid path; the chain for an
existing directory came back `["\\c:", "\\c:\\users", ...]` with no anchor; touchpoints `[]`;
two aliased destinations (`C:\Users\t\skills\parley-deck` vs `C:\RuntimeA\skills\parley-deck`)
produced `blocked=0`. The arithmetic is platform-deterministic, so on a real Windows box every
stat fails identically regardless of what exists. POSIX control in the same harness anchors
`/tmp` on `dev:ino` correctly.

This is a regression cycle 22 introduced: cycles 19-21 stat'd/`realpath`'d the **full resolved
path** (walk-up via `path.dirname`), which works on win32. It is invisible to the project's own
gates: CI is ubuntu-only (`.github/workflows/test.yml`), and the platform constants are read at
module load so no test fakes win32. But Windows is a shipped platform —
`dist/*windows-x64.exe` / `*-windows-arm64.exe`, `packaging/winget/`, README's
`winget install Feci.ParleyDeckSkill`. On those builds the plan gate this idea spent rounds
13-18 constructing silently does not exist: aliased containers, pass-through links, and buried
links all return to their pre-fix false-success behavior.

Fix direction (implementer's call, small): accumulate from `path.parse(resolved).root` with the
root stripped from the component list, in both functions, and pin it — either a unit test that
drives the path arithmetic through `path.win32` or a Windows CI leg. Without a pin, a third
model swap can re-break a platform nobody's tests run on.

### MINOR — the arm that decided the model has no regression test

The firmlink respelling **with an existing inner parent** — codex-1's round-18 arm that killed
both cycle 21 and my union — is not pinned anywhere. The existing firmlink test
(`test/bidding-addon.test.js:1706`) leaves `skills/` below the respelled home nonexistent, so
both spellings anchor on the same inode: the scalar-friendly shape that cycle 21's model already
passed. A future refactor to the union (or any scalar+firmlink hybrid) would run the full suite
green while reopening the arm. I measured the arm at `2b680a2` and it is handled (three of three
shapes blocked) — the code is right, the pin is missing. This is the same pattern codex-1 named
in round 18 ("the test that let it through"); that makes the fifth/sixth occurrence in this
idea. One test, cheap insurance.

### NIT — `canonicalSegment` NFC-folds on case-sensitive filesystems too

On Linux/ext4, NFC and NFD spellings are distinct objects, but two nonexistent destinations
differing only in normalization form would be judged one destination and refused. Fails safe,
requires deliberately adversarial input; noted for completeness, no action needed.

## Answers to the attacks the round posed (measured, no findings)

- **EACCES vs absent.** Robust. `stat` needs search permission on a path's *ancestors*, not on
  the final component — so a mode-000 directory still anchors the chain by its own `dev:ino`,
  through a symlink spelling as well as a direct one (measured: identical chains for both
  spellings of a subtree below a `chmod 000` parent). Every spelling fails identically below the
  anchor and name-keys off the same `dev:ino`, so consistency is preserved. End-to-end, an
  install with a destination below the unreadable ancestor fails closed
  (`destinationAncestorObstacle`/staging EACCES): `ok:false`, nothing written.
- **Same dev/ino for genuinely different objects across mounts.** The collision direction is a
  false-positive block — annoying, never unsafe. The dangerous reverse (one object visible under
  two mounts with distinct `st_dev`, e.g. a double-mounted NFS export) is undetectable by any
  client-side identity model, realpath included; that is residual physics, not a cycle-22
  defect.
- **`entryChain` with the parent itself a link.** Robust: `stat` follows, so the entry is
  located by the resolved parent's identity, byte-identical whether the parent is spelled
  through the link or directly; and a buried link's touchpoints do contain the enclosing
  destination's identity (both measured). Hardlinked directories are not creatable on the
  supported filesystems. Unicode normalization forms are NFC-folded and, for existing
  components, irrelevant — identity there is `dev:ino` (measured: NFD-created, NFC-spelled gives
  identical chains; nonexistent NFC/NFD pair blocked as one destination).
- **Uninstall quarantine vs install commit.** Identical coverage: `aliasedDestinations` runs
  before any mutation in both `installFleetAtomically` (`lib/installer.js:1489`) and
  `removeFleetAtomically` (`lib/installer.js:1604`), with the same fleet gate. Measured
  end-to-end: seed an install, alias a second runtime's home at it, `uninstall --target all` →
  `ok:false`, both blocked, the seeded destination untouched — nothing quarantined.
- **Rest of the single-process guarantee.** Staging can materialize another unit's *parents*
  (`mkdir -p`), but name-keyed tails anchor on the nearest existing identity at plan time and
  textual non-nesting cannot become nesting via `mkdir`, so no collision appears between gate
  and commit inside one process. Commit/revert/quarantine/cleanup failure directions all roll
  back or degrade to warnings with debris paths reported; I found no new hole in this sweep.
  (The concurrent-installer case remains the recorded round-14 follow-up — not re-litigated.)

## Release judgement

Not releasable as 2.1.0. One thing must change: make the chain's component-path accumulation
platform-correct (root-anchored, not `path.sep`-anchored) in `identityChain` and
`resolutionTouchpoints`, with a regression pin that exercises the win32 arithmetic. The POSIX
gate at `2b680a2` is, by everything I could measure, correct and complete against every arm this
idea has accumulated — the first cycle I can say that of — but the same binary ships for Windows
and there the gate is absent.

## What I verified

- **Facilitator's claims, all reproduced:** 361 node tests, 0 fail (`node --test`); python leg
  54/54 across 7 files on python3 3.14.6; a 3.9.6-first PATH refuses by design
  ("python3 is 3.9, but the add-on declares >=3.10"); manifest check ok — 47 files, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- **Cycle-22 diff read in full** (`64e43f9..2b680a2`), plus the full-range stat
  (`49fc3ec..2b680a2`) for scope; `IMPLEMENTATION.md` cycles 10-22 and `review/round-09/VOID.md`
  read. The three "Notes for reviewers" spots checked: the `manifestProblems` legacy exemption
  is properly scoped to sources shipping no manifest; the marker binds aggregate + manifest
  sha256, so the self-consistent swap is caught as designed; the `design-addons.test.js` delta
  since the round-7 tree is additive assertions, no weakened parameterization.
- **Measured arms** (scratch outside the repo, gate functions extracted verbatim from
  `2b680a2`, end-to-end arms via the real `installer` module; scratch removed after):
  1. firmlink respelling with existing inner parent — blocked in all 3 spelling/order shapes;
  2. EACCES middle component — chains identical across spellings; end-to-end install fails
     closed, nothing written;
  3. unicode NFC/NFD — nonexistent pair blocked as shared; existing NFD dir spelled NFC gives
     identical chains; children below normalization-equal parents blocked;
  4. `entryChain` through a parent link — identical keys across spellings; buried link's
     touchpoints contain the enclosing destination's identity;
  5. uninstall aliasing parity — blocked end-to-end, nothing quarantined;
  6. win32 simulation (`path.win32` + fake fs answering real win32 paths) — gate queries only
     `\C:`-prefixed non-paths, no anchor, zero touchpoints, junction aliasing unblocked; POSIX
     control anchors correctly. Arms 1-5 pass; arm 6 is the MAJOR above.
- **Platform story for the MAJOR:** CI ubuntu-only; Windows binaries + winget packaging shipped;
  cycles 19-21 stat'd full paths and worked on win32; cycle 22's component-rebuild is what
  introduced the `path.sep` assumption.
