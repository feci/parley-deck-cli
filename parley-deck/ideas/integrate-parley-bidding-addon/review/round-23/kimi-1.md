---
idea: integrate-parley-bidding-addon
review-round: 23
agent: kimi-1
date: 2026-07-31
reviewed-commit: 7e8ccec
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

The one outstanding block — codex-1's round-22 coverage block — is **closed**, and I measured
the closure rather than trusting the description. `splitAtRoot` was made injectable in cycle 23
while `walkRawTarget` kept a second, non-injectable copy of the same root arithmetic, so the
drive and UNC arms of the defect cycle 24 fixed had nothing executable asserting them. Cycle 26
extracts `rawTargetArithmetic(from, rawTarget, impl, resolveParent)`, exports it
(`lib/installer.js:1409`), makes `walkRawTarget` delegate to it, and adds a regression that
asserts against the **real exported helper** with `path.win32`: a drive-absolute target does not
replay `C:\` as a component; a UNC target does not duplicate server and share; a relative target
anchors on the resolved parent with the resolver observed being consulted; a POSIX target keeps
a backslash as an ordinary filename byte.

I confirmed the regression discriminates: against `git archive` copies of `b1f43e4` and
`2b7ca3e` with the current test file dropped in, it fails at both (exit 1, 1 fail; at `b1f43e4`
the failure is `TypeError: installer.rawTargetArithmetic is not a function` — the coverage
genuinely did not exist there, which is exactly what the block said). At `7e8ccec` it passes.

I also checked the extraction is behavior-preserving, not just plausible: the production call
passes `impl = path`, so `impl === path.win32 || process.platform === "win32"` reduces to the
old `process.platform === "win32"` condition, and the injected resolver is the same
`realpathSync`-with-lexical-fallback the inline code had. The cycle-26 payload diff is the
extraction, the export, and tests — nothing else.

No other outstanding findings of mine remain open. The concurrent-installer isolation follow-up
is recorded, the warning sits verbatim in `CHANGELOG.md` under "Known limits" — settled in round
14, not re-opened.

## Position on the gate

**1 (correct).** Unchanged from round 22, and re-checked against the text rather than carried by
inertia. The gate checks the complete target x unit plan for physical-identity collisions before
the first write or deletion and refuses the whole plan on any predictable failure — aliased
runtime roots, firmlink respellings, destinations nested through symlink chains, raw link
targets walked physically. `CHANGELOG.md:49` promises exactly that ("atomic across the whole
fleet ... a predictable failure anywhere produces zero writes and zero deletions") and names no
exclusion for symlinked runtime homes — support earlier cycles deliberately built ("a
destination parent that is a symlink to a real directory is not an obstacle"). Position 2 stays
refused on the same ground as round 22: narrowing the gate would either break that supported
case or make the written promise false.

## New findings

One cosmetic observation, stated for the record and explicitly **not** a block: cycle 26
re-added three already-existing tests verbatim — "a link reached through an earlier link is
walked from where it physically sits" (`test/bidding-addon.test.js:1697`, duplicated at 1857),
"a backslash inside a POSIX link target is a filename byte" (1746, duplicated at 1906), and "an
intermediate link in a raw target is expanded before `..` is applied" (1784, duplicated at
1944). I diffed each pair: byte-identical. Both copies run and both pass; `node:test` permits
duplicate names. Cost to a real user: zero. Cost to the suite: three redundant tests (372
counted, 369 unique) and a copy/paste wart in a file whose comments elsewhere hold themselves to
a higher standard. A one-line cleanup, before or after release; it does not gate anything.

Otherwise: None. I am not holding back a residual "filesystem exotics" worry — the arms every
round actually named are now each pinned by an executable assertion, and the one arm that could
not be made to discriminate (the POSIX backslash tear) is honestly labelled a pin in the test
itself rather than counted as a proof.

## Release judgement

**Releasable as 2.1.0.** The payload has not changed since `714712f` (zero diff in `skills/` and
`bin/`; manifest aggregate unchanged), no round has found a defect in it, all three reviewers
hold position 1 on the gate, and the single outstanding block is closed with a regression that
fails at the two commits it names. Nothing real remains.

## What I verified

- `git diff b1f43e4..7e8ccec` read in full: `lib/installer.js` (extraction + export only) and
  `test/bidding-addon.test.js` (+160, no deletions). Behavior parity of the extraction confirmed
  by inspection of `lib/installer.js:1409-1470`.
- `node --test` at `7e8ccec` (node v26.5.0): **372 tests, 372 pass, 0 fail**.
- `node scripts/run-python-tests.js`: **54/54 on python 3.14.6**. Re-run with
  `PATH="/usr/bin:$PATH"`: refuses 3.9.6 with `python3 is 3.9, but the add-on declares >=3.10`,
  exit 1 — fail-closed by design.
- `node scripts/build-addon-manifest.js --check`: **ok — 47 files,
  sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d**, matching the
  claimed aggregate.
- `git diff 714712f..7e8ccec -- skills/ bin/`: zero lines — payload unchanged.
- Regression discrimination: `git archive b1f43e4` / `git archive 2b7ca3e` into `mktemp -d`
  dirs, current test file copied in, `node --test --test-name-pattern="root arithmetic"` —
  both exit 1 with 1 fail (`installer.rawTargetArithmetic is not a function` at `b1f43e4`);
  the same command passes at `7e8ccec`. Temp dirs removed with `rm -rf`; `/tmp/pd-r23-*`
  verified absent afterwards.
- Duplicate-test claim measured, not eyeballed: `diff` of each of the three line ranges —
  all identical.
- `IMPLEMENTATION.md` records through "Fix-up cycle 26" with a measured-after section;
  `CHANGELOG.md` carries the single-writer limit under "Known limits" and the fleet-wide
  atomicity promise at line 49.
- Working tree untouched: `git status --short` clean before and after; no edits under
  `skills/parley-bidding/`, no scratch files in the repo; archive comparisons done in /tmp.
