---
idea: integrate-parley-bidding-addon
review-round: 15
agent: codex-1
date: 2026-07-31
reviewed-commit: 26478e9
---

## Verdict

BLOCK

## Outstanding findings — closed or not

The exact round-14 reproductions are closed. Two empty runtime roots symlinked to one physical
root are refused before staging, and a fresh case-only Codex/Kimi home pair on this
case-insensitive volume is refused with zero writes. `readManifest` now applies the same
regular-file predicate as the other three exported manifest operations. The installed-unit
uninstall regression now promises zero removals when one installed unit blocks.

Those closures are not complete over the full input space, however. The new findings below
show a different physical-overlap route, an absent-unit dry-run branch left on the old result
semantics, and a manifest consumer outside the four exported operations.

The concurrent-installer question remains settled exactly as round 14 ruled it: cross-process
isolation is a recorded follow-up and does not gate 2.1.0. `CHANGELOG.md` contains the adopted
single-writer warning verbatim. None of the findings below uses a second process.

## New findings

### [MAJOR] Planned ancestor/descendant destinations can silently erase a committed unit

`aliasedDestinations` groups only identical `physicalKey` values
(`lib/installer.js:1385-1434`). It does not reject one planned destination physically inside
another. That is not merely an unusual spelling: staging creates destination parents
(`lib/installer.js:1689-1696`), so one unit can materialize another unit's destination after
planning and before commit.

I set `KIMI_CODE_HOME` to a fresh directory and `CODEX_HOME` to Kimi's planned core destination,
so the Codex core destination was `<kimi-core>/skills/parley-deck`. A single real
`install --target all --include-undetected` at `26478e9` returned:

- top-level `ok: true`;
- Codex core `action: installed`;
- Kimi core `action: replaced`;
- the reported Codex destination did not exist afterwards;
- the surviving ancestor marker recorded `target: kimi`.

The sequence is deterministic. Codex staging creates the Kimi destination as a parent, Codex
commits inside it, the later Kimi commit renames that whole ancestor to its backup, and phase-3
cleanup deletes the backup containing the supposedly successful Codex install. This is the
single-process false-success and partial-fleet state that B5 and the 2.1.0 changelog expressly
claim cannot occur.

This blocks 2.1.0. Reject physical overlap, not only equality: no planned destination may equal,
contain, or be contained by another after resolving existing ancestors. Regress both nesting
directions for install, uninstall, and dry-run before relying on staging/rollback. Also make the
case-only regression construct two targets in one plan: the current test at
`test/bidding-addon.test.js:1628-1641` only installs one generic target and then runs `doctor`
through another spelling; I copied that test to `d7ab1c3` and it passed there, so it does not
discriminate the alias fix.

### [MINOR] Uninstall dry-run and real output still disagree for absent units

The cycle-18 gate delays `remove` results, but it still records absent units as `missing` and
`ok: true` before checking `blockedAnywhere` (`lib/installer.js:1579-1613`). The gate rewrites
only `pending` units. Meanwhile the real-only preflight in `uninstallCommand`
(`lib/installer.js:664-696`) rewrites every unblocked unit, including absent ones, to `skipped`
and `ok: false` when any unit blocks.

Measured with all targets absent except a foreign unmarked Aionrs core: both commands returned
top-level `ok: false`, but dry-run reported every Codex unit `missing`, `ok: true`, with no
message, while the real command reported those same units `skipped`, `ok: false`, with
`Not attempted: another skill or target in this uninstall failed preflight.` The Aionrs add-ons
showed the same divergence beside the blocked core.

This is another arm of the round-14 dry-run finding and blocks this strict gate at MINOR. Remove
the real-only duplicate preflight or make both modes use one result-construction path, then
compare per-unit `ok`, action, and refusal message for installed, absent, and mixed fleets.

### [MINOR] The Python test runner is a fourth manifest reader outside the shared predicate

`scripts/run-python-tests.js:47-55` reads and parses `parley-addon.json` directly. It neither
calls `readManifest` nor enforces the regular-file rule; parse or read errors silently become
"no declared floor."

I replaced the manifest in a temporary archive with a symlink to a byte-identical external
file. The module consistently returned `hasManifest: false`, `readManifest.ok: false`,
`manifestFileHash: null`, and `verifyPayload.ok: false`, while
`node scripts/run-python-tests.js` still returned success with **54/54** tests. The full
`npm test` has later checks that catch the source defect, but the standalone Python gate consumes
external manifest policy and reports green, contradicting the claim that every manifest read
shares the predicate.

This blocks the strict gate at MINOR rather than being a 2.1.0 follow-up. Have the runner use
`readManifest` and fail closed when the required manifest or its Python floor is unusable.

## Release judgement

Not releasable as 2.1.0 at `26478e9`. The one indispensable release fix is to make the
single-process fleet guarantee true for physically overlapping planned destinations; this
strict gate also requires the two measured MINOR completeness gaps above to close, followed by
a fresh full-scope clean review.

## What I verified

- Read the live cooperation protocol, `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through
  cycle 18, the round-14 reviews, `review/round-09/VOID.md`, the current changelog, and the full
  current transaction/manifest code. Read both requested diff scopes; `git diff --check` is
  clean for `d7ab1c3..26478e9` and `49fc3ec..26478e9`.
- Used a code-only Graphify traversal from a temporary archive to map the manifest consumers and
  the preflight/stage/commit/revert/quarantine/cleanup call graph, then checked the mapped code
  line by line. Ordinary equal symlink aliases are gated in both mutations; `..` is normalized
  before filesystem use; no new evidence invalidates the settled cross-process disposition.
- Ran `npm test` in an archive of `26478e9`: **353/353 Node tests**, **54/54 Python tests** on
  Python 3.14, and the 47-file manifest check with aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`. Ran the Node suite
  again with `/usr/bin/python3` 3.9 first on `PATH`: **353/353**.
- Ran `npm pack --dry-run --json --ignore-scripts`: **202 files**, **48** under
  `skills/parley-bidding/`, with no `__pycache__`, `.pyc`, or `.pyo`. All four adapters validate
  and all **16** bidding JSON files parse.
- Confirmed `skills/parley-bidding/` is unchanged from `714712f`; the source comparison shows
  only the recorded rename/integration differences, dropped nested `.gitignore`, consent text,
  schema IDs, and added manifest. Static scans found no symlink/cache artifact, customer or
  credential material, unresolved placeholder, or portal/network-capability import in the
  deterministic Python tools.
- Reproduced the three new findings independently in temporary copies: nested destinations
  returned false success and lost Codex; absent-unit uninstall results differed by `ok`, action,
  and message; and the Python runner accepted the symlinked external manifest while all four
  manifest-module operations refused it.
