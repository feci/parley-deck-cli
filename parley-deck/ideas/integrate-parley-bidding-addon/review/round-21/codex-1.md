---
idea: integrate-parley-bidding-addon
review-round: 21
agent: codex-1
date: 2026-07-31
reviewed-commit: 381e639
---

## Verdict

BLOCK

## Outstanding findings — closed or not

- The round-20 intermediate-component arm is closed. When `mid` itself is encountered inside
  the raw target, `walkRawTarget` expands it before applying the later `..` components. The new
  regression fails at `2b7ca3e` and passes at `381e639`.
- The absolute Windows root arithmetic is corrected by inspection: an absolute target's parsed
  root is removed before the remaining components are walked. It is not closed by a
  discriminating regression, however; the added test is a POSIX pin that passes at both commits
  and never exercises drive or UNC roots. See the MINOR below.
- The linked-ancestor arm inside kimi-1's round-20 MAJOR is **still open**. The recursion expands
  links it enters inside `rawTarget`, but it starts a relative target at the lexical
  `path.dirname(from)`. If `from` was reached through an earlier symlink in the destination,
  the kernel starts from the already-expanded physical parent and leading `..` components walk
  a different tree. I measured the false-success install and uninstall states at `381e639`.
- The prior plain-containment, firmlink, link-entry, link-chain, simple raw-`..`, marker,
  manifest, runtime, selection, transaction, rollback, cleanup, and result-shape arms remain
  closed in the code and regressions I rechecked. Concurrent-installer isolation remains the
  unanimously recorded follow-up and was not re-litigated.
- The canonical implementation record is not current: it still ends at cycle 23 and names
  `2b7ca3e`, despite the brief saying cycle 24 is recorded. See the second MINOR below.

## Position on the gate

**3 — still wrong: the linked-ancestor / physical-parent arm remains open.**

The gate is correct for the six shapes its current tests actually exercise, but the walker is
not yet the kernel model its comments claim. Expanding a link after entering it is only half of
physical resolution: a relative link target must also begin at the physical parent reached
through all earlier path components. At `lib/installer.js:1413`, it instead begins at the
parent's spelling. A second platform mismatch at line 1414 treats a backslash as a separator on
POSIX, where it is an ordinary filename byte. Both mismatches admit a green fleet mutation that
breaks the fleet.

I do not choose position 2 for this reviewed tree. Symlinked runtime homes were deliberately
supported by earlier fixes, and the changelog currently promises fleet-wide atomicity without
excluding these configurations. Narrowing remains a possible design decision, but it would
require an explicit scope change and release-note limit; it is not what `381e639` implements or
claims.

## New findings

### [MAJOR] The raw-target walk still diverges from kernel traversal in two measured arms

`resolutionTouchpoints` keeps walking the destination's lexical spelling after it encounters an
ancestor link (`lib/installer.js:1461-1479`). Later, `walkRawTarget` anchors a relative target at
`path.dirname(from)` (`lib/installer.js:1413`) rather than at that directory's physical location.
This leaves kimi-1's linked-ancestor arm from round 20 unfixed.

I reproduced it with `linkA -> real/A`, the Codex container link at
`linkA/container/skills`, and its raw target
`../../Btree/skills/parley-deck/inner/deep`. Kimi's planned core is
`real/Btree/skills/parley-deck`; `inner` is a link inside it to an outside directory. The kernel
has already expanded `linkA`, so the first two `..` components climb from
`real/A/container` to `real`. The walker climbs from the spelling `linkA/container` to the
home directory and never records Kimi's core.

At `381e639`, fleet install returned `ok:true`, reported Codex `installed` and Kimi `replaced`,
then left Codex's reported destination absent and its payload orphaned under the outside
directory. After seeding the two units separately, fleet uninstall also returned `ok:true` and
reported both `removed`, while leaving a hidden
`.parley-deck.<pid>.<timestamp>.removing` directory behind. The install reproduction has the
same false-success result at `2b7ca3e`, confirming that cycle 24 did not close this round-20
arm.

There is a second, new path to the same result. `rawTarget.split(/[\\/]+/)` treats `\` as a
separator on POSIX. I created a real symlink whose one-component name was `mid\link` and used it
in a raw target that entered Kimi's core and walked back out with `..`. The kernel followed the
one component; the gate split it into `mid` and `link`, missed the dependency, returned
`ok:true`, reported both units successful, left Codex's destination absent, and orphaned its
payload. This configuration is unusual, but it is valid POSIX path syntax and directly
contradicts the claim that the walker records what the kernel consults.

The fix should use one physical cursor for the complete destination resolution: carry forward
the expansion of earlier ancestors, base each relative link target on the physical containing
directory, and split components according to the active platform (`/` only on POSIX). Add
discriminating install and uninstall regressions for the linked-ancestor case and a POSIX
literal-backslash regression, while retaining the cycle-24 intermediate-link arm.

### [MINOR] The Windows absolute-target blocker has no regression that exercises its fix

The test at `test/bidding-addon.test.js:1731-1750` creates and follows an ordinary POSIX
absolute symlink. It does not inject `path.win32`, construct a drive or UNC target, inspect the
walked components, or exercise `walkRawTarget`'s root slicing. Its `seen` array and `record`
callback are unused. As disclosed in the brief, the test passes at both `2b7ca3e` and
`381e639`; it is a preservation pin, not evidence that the round-20 Windows defect fails before
the fix and passes after it.

The code change looks correct for drive and UNC roots, but `FINAL.md` requires every blocker to
be closed by a test that fails without the fix, and this Windows channel still has no executing
installer CI. Extract an injectable raw-target splitter/walker, as was done for `splitAtRoot`,
and assert drive, UNC, mixed-separator, and root-relative component sequences against the real
helper; retain the POSIX public-gate pin separately.

### [MINOR] `IMPLEMENTATION.md` does not record cycle 24 or the reviewed commit

The canonical record's frontmatter still says `status: fix-up-cycle-23` and
`head-commit: 2b7ca3e`; its last cycle is 23 and its last measurement is 364 Node tests. There
is no cycle-24 section for `381e639`, its 366-test result, the one discriminating regression,
or the absolute-root pin. This contradicts the review brief and the Phase-8 requirement to
update `IMPLEMENTATION.md` before re-review. Append cycle 24, update the frontmatter, and keep
the evidence scoped accurately; the current claim that all accumulated arms are refused is
disproved by the MAJOR above.

## Release judgement

Not releasable as 2.1.0. The one blocking subsystem is destination-dependency resolution: it
must traverse raw targets from the physical parent with platform-correct component semantics,
and that behavior must be pinned by discriminating tests and reflected in the canonical
implementation record. No change to the `parley-bidding` payload is indicated.

## What I verified

- Read the live `COOPERATION.md`, `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md` through its
  actual cycle-23 ending, `review/round-09/VOID.md`, and the round-19/20 review artifacts. The
  live deck is the protocol source; the installed skill/status check reported only advisory
  metadata/reference drift, and its required sync was run dry-run only.
- Read the complete cycle-24 diff `2b7ca3e..381e639`, the full gate-era diff
  `49fc3ec..381e639`, and the implementation files against `FINAL.md`. Both requested diffs
  pass `git diff --check`; `skills/parley-bidding/` has no change after `714712f`.
- Ran `npm test` from a `git archive` of `381e639`, with the repository dependency tree linked
  read-only and all generated homes under one disposable temp root: **366/366 Node tests**,
  **54/54 Python tests** on Python 3.14, and the 47-file manifest check at
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- Re-ran the complete Node suite with `/usr/bin/python3` 3.9.6 first on `PATH`: **366/366**.
  The Python runner separately refused 3.9.6 with exit 1, as its `>=3.10` contract requires.
- Ran the two cycle-24 tests against the prior installer: the intermediate-link regression
  failed at `2b7ca3e`, while the absolute-target POSIX pin passed there; both pass at
  `381e639`.
- Reproduced the linked-ancestor false-success install at both `2b7ca3e` and `381e639`, and the
  false-success uninstall with quarantine residue at `381e639`. Reproduced the literal-POSIX-
  backslash install arm at `381e639`. All fabricated homes were removed.
- Parsed `npm pack --dry-run --ignore-scripts --json` with an isolated cache: **202 files**,
  **48** under `skills/parley-bidding/`, and no `__pycache__`, `.pyc`, or `.pyo` entry.
- Compared the read-only source and integrated payload by path and SHA-256: 48 files on each
  side; source-only `.gitignore`; integrated-only `parley-addon.json`; and the same nine
  documented content differences. Parsed all **16** JSON files and ran the validator
  successfully against all **4** platform adapters.
- Rechecked the consent text, six-skill documentation, single-writer warning, interpreter
  reporting, manifest/marker trust boundary, and deterministic-script imports. I found no
  portal-mutation capability or payload change.
- Removed the complete disposable archive and test-home tree (about 862 MB after the suite's
  immutable-debris cases). The reviewed repository remains clean at `381e639`; no file under
  `skills/parley-bidding/` was changed.
