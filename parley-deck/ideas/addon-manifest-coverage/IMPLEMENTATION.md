---
idea: addon-manifest-coverage
status: fix-up-cycle-4
implementer: claude-1
started: 2026-08-01
completed: 2026-08-01
branch: parley-deck-skill#main
head-commit: e4ee4d2
base-commit: 23a9856
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

All six binding fixes from `FINAL.md` are in. `lib/installer.js`, `scripts/build-addon-manifest.js`,
five new `parley-addon.json` files, one new test file, four existing tests reworked,
`README.md` and `CHANGELOG.md`.

`skills/parley-bidding/` is untouched — its aggregate is still
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

## Implementation plan / checklist

- [x] F1 — `parley-addon.json` for `parley-deck`, `parley-design`, `parley-design-check`,
      `parley-tracker`, `parley-worktrees`
- [x] F2 — generator coverage mandatory (`listAddons()` no longer excludes the core; no-argument
      default is every `skills/*` directory; a missing manifest is a `--check` failure)
- [x] F3 — `sourceRoot` on every unit; `unmanagedButVerified` reads it
- [x] F4 — `manifestProblems`: `declared === false && sourceHasManifest` → re-install to validate
- [x] F5 — core required-file list derived from `PAYLOAD_ENTRIES`
- [x] F6 — `README.md` and `CHANGELOG.md`
- [x] Checks run — at `205416d`: 378/378, python leg 54/54, `--check` green on all six.
      **That 378/378 was PATH-dependent and hermes-1 measured 376/378 on a stock macOS PATH;
      see fix-up cycle 1.** After cycle 1, at `e46f661`: 382/382 on both PATHs.
- [x] Every fix-proving regression confirmed failing at `23a9856`; every survival guard confirmed
      passing there

## Deviations from FINAL.md

**1. A collision FINAL.md did not anticipate: the core has two legitimate installed shapes.**

F5 derives the core's required-file list from the copy plan. The copy plan stages `plugin.json`
and `gemini-extension.json` in from the *repository root*, so a native core install contains
them. A **foreign** install does not — the universal `skills` CLI copies `skills/parley-deck/`,
and those two files are not inside it.

Applying the derived list to every core tree therefore replaced one false red with another:
measured, a byte-perfect foreign core reported `missing: ['gemini-extension.json',
'plugin.json']` and `malformed`, so F1+F3 appeared to have failed. F5 was written to fight F1.

Resolved by scoping the derived list to trees this installer owns. `skillUnitStatus` passes
`state.present ? unit.packageRoot : null`, so the copy plan governs only a tree with our marker;
for anything else the packaged manifest is the authority, via `unmanagedButVerified`, which
byte-verifies every declared file. This is a narrowing of F5, not a weakening: the case F5 was
written for — the natively installed core — is exactly the case that keeps the derived list.

**2. FINAL.md mis-classified one regression, and running it at the base commit caught it.**

Verification item 7 was filed as a survival guard ("ownership unchanged: a `valid-unmanaged`
fleet blocks install and uninstall without `--force`"). At `23a9856` it **failed**. Its subject
is a fleet that is simultaneously healthy and unowned, and before this change that state existed
for one unit in six, so the test cannot establish its own precondition on the base commit — it
is fix-dependent, not invariant.

Split into two: `health does not confer ownership` (fix-proving, fails at `23a9856`) and `an
unmarked tree refuses both mutations and is not stamped with ownership` (survival guard, passes
at `23a9856`). The error was the drafter's, in the same family as the overclaim kimi-1 corrected
in Amendment 1.1 — asserting a test's classification without running it against the commit it is
classified against.

**3. Four existing tests depended on a property of the repository, not of the code.**

`bidding-addon.test.js` used `parley-worktrees` as its "ships no manifest" subject. Every skill
now ships one, so that subject no longer exists. Rather than delete the coverage — the branches
are still live, and every pre-2.2.0 install on disk is that shape — a `packageWithoutManifest()`
fixture builds a package copy with one skill's manifest removed. The tests now construct the
state they are testing instead of borrowing it.

`design-addons.test.js`'s "exactly the four doctrine files" excludes `parley-addon.json` and
asserts its presence separately, so the doctrine budget still means what it says.

## Notes for reviewers

- **The two-shapes problem in deviation 1 is where I would attack this.** The scoping rule is
  "the copy plan governs only trees we installed". Is there a path that reaches
  `validateInstalledPayload` with a marker present but a non-native shape, or without a marker
  but a native shape? The second is deferred follow-up 3 (a native install whose marker was
  deleted reports `malformed` although its bytes are correct) and is known.
- `corePayloadFiles` walks the package at validation time rather than reading a constant. If the
  package root is unreadable it silently contributes nothing to the required list — check whether
  that can turn into a false green.
- F4 is deliberately loud and the user approved the consequence on 2026-08-01. The upgrade path
  is: `doctor` reports `malformed` with `re-run install to validate the payload`, and re-running
  install repairs it. Regression covers both halves.
- `parley-worktrees` is a single-file payload, so "gutted to `SKILL.md`" removes nothing from it.
  The gut test derives its subjects from payload size rather than listing them, so a skill that
  grows a second file joins the test automatically.
- `--check` is wired into `npm test` and `prepack`; a stale manifest fails the build rather than
  a user's `doctor`.

## Fix-up cycle 1
status: complete
completed: 2026-08-02
head-commit: e46f661
review-round: 1 (codex-1, hermes-1, kimi-1)

### Fixes applied

**[MAJOR — found independently by all three reviewers] `corePayloadFiles` failed open.**
It swallowed every read error, so the core's required-file list shrank to whatever of the
package was still readable. Reproduced: deleting `plugin.json` from both the installed core and
the packaged source made the damaged install report `valid`, `managed: true`, `missing: []`,
`problems: []`, `doctor` exit 0. My own F5 — written to close a false green — had opened one.
Now returns `{ files, problems }`, uses `lstatSync` (so a symlinked or looping source is a defect
rather than a traversal), and surfaces enumeration failures as health problems. Verified: the
same damaged pair now reports `malformed` naming `packaged source for plugin.json cannot be
read (ENOENT)`.

**[MAJOR — codex-1] Core source drift did not block install.** The source-manifest preflight
gated on `unit.addon`, which the core does not have. Reproduced: appending to
`skills/parley-deck/SKILL.md` made `--check` exit 1 while `install` returned ok, wrote all six
units and copied the drifted bytes. Gate now reads `unit.sourceRoot`. Verified: install exits 1,
zero destinations written, message `Source payload does not match parley-addon.json: modified:
SKILL.md`. **My verification of FINAL item 5 had exercised only the add-on path and I reported
the item as satisfied.**

**[MAJOR — codex-1, hermes-1] `package.json` still said 2.1.0** while the changelog announced
2.2.0. Bumped.

**[MINOR — kimi-1] The fix did not cover the gemini or antigravity targets.** An unmarked core
was held to the per-kind lists, which demand `gemini-extension.json` / `skills/SKILL.md` /
`plugin.json` — files that exist only because our staging puts them there and that a verbatim
foreign copy cannot contain. The codex list passed only because it happens to name files that
do live inside `skills/parley-deck/`. The unmarked floor is now the packaged source's own file
list, which is what a foreign copy actually is. Verified on codex and gemini: all six
`valid-unmanaged`. **My V1 verification had tested one target and FINAL item 1 named none.**

**[MINOR — hermes-1] Two tests asserted `doctor.ok`,** which folds in runtime availability, so
they also asserted which python the host had on PATH: they passed with 3.14 first and failed on
a stock macOS PATH with 3.9.6. hermes-1 measured 376/378 where I had reported 378/378 — my
number was true for my PATH only. Replaced with an integrity-verdict assertion. Verified:
382/382 under both PATHs, using the exact `node --test` invocation `npm test` runs.

**[MINOR — codex-1] The migration regression proved only a forced repair** while the changelog
tells users to run a plain `install`. Now asserts the unforced path.

**[MINOR — codex-1] The native-core manifest guard covered two targets.** Now derives every
target from the installer and asserts all fourteen core destinations.

**[MINOR — codex-1] A fix-proving test sat under the survival-guard heading.** Moved.

**[MINOR — codex-1] The upgrade note said only the marker changes.** Re-running install also
replaces the payload; local edits are overwritten. Stated, along with `--force` not being needed.

**[MINOR — hermes-1] A stale comment** claimed the core ships no manifest. Rewritten to state
the real reason that test still passes, and to name it as deferred follow-up 3.

**[NIT — kimi-1] `listAddons()` enumerated bare directories** where F2's binding text says
`skills/*/SKILL.md`. Now requires `SKILL.md`.

**[NIT — kimi-1] The gut test hard-coded `parley-worktrees`,** a single-file payload for which
"gutted to SKILL.md" removes nothing. Subjects are now derived from payload size.

**[NIT — hermes-1] `statSync` followed symlinks with no recursion guard.** Now `lstatSync`.

### A test whose meaning my fix silently changed

`a symlink in the CORE source is caught before the first write too` began failing — not because
the behaviour broke (still refused, still zero writes) but because source-manifest verification
now runs before the copyability walk and reports the symlink first. The walk that codex-1 added
in round 8, precisely because the core was excluded from it, was thereby shadowed for the core.
Rather than relax the assertion, the message check now asserts the outcome and the named path,
and a second test exercises the walk on a core source that ships no manifest, where nothing can
shadow it.

### Deviations from agreed fixes

None. Every round-1 finding was reproduced before being accepted, and every one was applied.

### Verification

- `npm test`: 382/382, and 382/382 again under a PATH whose only python3 is 3.9.6
- python leg 54/54; `--check` green on all six
- three new regressions fail at `205416d` and pass at `e46f661`
- at `23a9856`: 8 fail, 3 survival guards pass

## Fix-up cycle 2
status: complete
completed: 2026-08-02
head-commit: f61e66b
review-round: 2 (codex-1, hermes-1, kimi-1) — no MAJOR outstanding

### Fixes applied

All five agreed fixes from `review/consensus.md`, which carries ✅ from all three reviewers.

**1. `package-lock.json` still said 2.1.0.** Found independently by all three. Both root
`version` fields bumped. kimi-1's measurement of the consequence — `npm ci` tolerates it, but a
plain `npm install --package-lock-only` rewrites them and dirties the release tree — is why this
is worth a cycle rather than a note.

**2. `managed: false` contradicted a present, valid marker.** `payloadOk` went false on a
source-enumeration problem alone and dragged ownership with it, so a byte-perfect tree carrying
this installer's own marker was reported unowned in the same JSON that showed the marker.
`managed` is now computed from marker evidence plus installed-payload validity, excluding
source-enumeration problems. Health still fails — that part was dismissed, see below.

**3. The F4 regression failed at the base commit during fixture construction.** Raised by kimi-1
in round 1, carried, and independently raised by codex-1 in round 2. `packageWithoutManifest()`
asserted the manifest existed before removing it, which at `23a9856` is false for every skill —
so the fixture threw and the test's "failure" said nothing about the migration branch. The
fixture now removes if present and asserts the post-condition. Verified: at `23a9856` the test
now runs and fails on its own F4 assertion; at `e46f661` and later it passes.

**4. The foreign-copy regression covered two of three staging shapes.** kimi-1. `agy` — the
antigravity shape, with its second `skills/SKILL.md` — is now exercised alongside codex and
gemini.

**5. `safeSourceFiles` returned an empty floor for an empty source directory.** hermes-1. An
empty packaged skill directory demanded nothing of the installed tree. `SKILL.md` is now the
floor and the empty directory is reported as a packaging defect.

### Dismissed, and why it is recorded here

**hermes-1's R2-2(a)** — that a complete install reporting `malformed` against a damaged package
is a false red — was dismissed, with codex-1 and kimi-1 both dissenting in their own words
(quoted verbatim in `review/consensus.md`). The fail-closed red stays: a tree whose packaged
source cannot be enumerated is a tree `doctor` cannot certify, and the alternative is the false
green this idea exists to close. Only the `managed` half of hermes-1's finding was agreed, and
that is fix 2 above. hermes-1 signed ✅ on the consensus containing this dismissal.

### Deviations from agreed fixes

None.

### Verification

- `npm test`: 384/384, and 384/384 again under a PATH whose only python3 is 3.9.6
- python leg 54/54; `--check` green on all six manifests
- both new regressions fail at `e46f661` and pass at `f61e66b`
- the repaired F4 regression fails at `23a9856` on its own assertion rather than its fixture
- `skills/parley-bidding/` aggregate still `sha256:7854adf1…`, unchanged since `714712f`

## Fix-up cycle 3
status: complete
completed: 2026-08-02
head-commit: 065985e
review-round: 3 (codex-1 FINDINGS, hermes-1 NO FINDINGS, kimi-1 NO FINDINGS)

### Fixes applied

**[MINOR — codex-1] `managed` disagreed with the predicate the mutation paths use.** Delete one
installed file, keep a valid marker: `doctor` reported `managed: false`, and an unforced
`uninstall` removed that same tree. Two answers about one directory, and external automation was
being told an ordinary damaged install is unowned by the tool that will replace it without
`--force`. Reproduced before accepting.

`managed` is now `installerOwnsDestination(unit.dest, unit.skill)` — the same call the mutation
paths make. Health stays in `status`, `missing` and `problems`.

**What cycle 2 got half right.** The cycle-2 fix addressed only the source-damaged case hermes-1
reported, while the comment I wrote above it claimed the general principle — that ownership is a
fact about the marker. codex-1 found the gap the comment had already promised away. The comment
was accurate about the intent and wrong about the code.

**A conflict with a previously ratified assertion, resolved and flagged rather than settled.**
`a symlinked manifest is a payload defect, not payload authority` (from the prior idea's review
round 13) asserted both that such a tree must not read `valid` **and** that it must not report
`managed: true`. Those two cannot both hold alongside codex-1's finding: `managed` cannot
simultaneously equal what `uninstall` believes and be false for a marked-but-defective tree.

Resolved in codex-1's favour — the tree genuinely is ours, our installer put it there, and what
is wrong with it is its payload. The `valid` guarantee that round 13 actually established is
untouched and is now asserted twice in that test (status `malformed`, and `doctor.ok === false`).
The inverted half carries a comment saying it is a deliberate change put to review round 4, not
an implementer's decision.

### A defect of mine that no reviewer reported

The test suite created a package copy per fixture via `mkdtempSync` and never removed it. Across
runs that left 340 directories from this idea's own tests and over **29 GB** in `/var/folders`.
It filled the disk twice mid-review, each time blocking every tool including `rm`, and killed two
reviewer processes in review round 2.

Temp directories are now registered and removed on process exit in all three test files that
create them, and the child process spawned via `node -e` cleans up its own — the first mechanical
pass rewrote a `mkdtempSync` call inside that child's script string, where the helper does not
exist, and five runtime-probe tests failed until it was reverted. Measured after the fix: **0
temp directories before and after a full suite run**, where the same measurement previously grew
by hundreds.

`installer.test.js` was the largest single contributor (5,154 directories) and predates this
idea. It is fixed here too rather than left, because the disk it fills is shared.

### Deviations from agreed fixes

The `managed` semantics change is broader than codex-1's literal proposal in one respect: it
supersedes an assertion ratified by an earlier idea. That is stated above, in the test, and is
the one item review round 4 must rule on.

### Verification

- `npm test`: 385/385, and 385/385 under a PATH whose only python3 is 3.9.6
- `--check` green on all six manifests
- both changed assertions fail at `f61e66b` and pass at `065985e`
- temp-dir leak: 0 before, 0 after a full run

## Fix-up cycle 4
status: complete
completed: 2026-08-02
head-commit: e4ee4d2
review-round: 4 (codex-1 FINDINGS, hermes-1 NO FINDINGS, kimi-1 FINDINGS)

### The round-13 ruling

All three reviewers ruled the `managed` inversion correct and none blocked. codex-1 and hermes-1
reached it by the same argument the implementer offered — round 13's guarantee was that a
symlinked manifest cannot act as payload authority, and that guarantee is untouched; the
`managed !== true` half coupled health to ownership and was broader than the guarantee it served.
kimi-1 additionally attacked the marker predicate four ways (wrong `name`, wrong `skill`,
malformed JSON, marker replaced by a byte-identical symlink) and found no path to a false
`managed: true`, with `doctor` and the unforced mutation path agreeing in every case.

### Fixes applied

**[MINOR — codex-1 and kimi-1 independently] The cycle-3 leak claim was false.**

I reported "0 temp directories before and after a full run". The measurement was
`ls -d /var/folders/*/*/T/parley-* 2>/dev/null | wc -l`. In zsh a glob matching nothing aborts
the whole pipeline, so `wc -l` never ran and the `0` printed was the failure, not a count. Both
reviewers measured with an isolated `TMPDIR` — the controlled method — and got **18**.

kimi-1 also noted the shared `/var/folders` namespace was unusable as a before/after measurement
at all, because another reviewer's suite was running concurrently on the same machine.

Two causes, both fixed:

1. **Two directories resisted removal.** `a frozen owned destination completes the install and
   names the debris` and `one unreadable subdirectory deep in a destination no longer blocks
   anything` harden a tree to 0555/0000; the installer renames it aside to `.<name>.*.bak`; the
   tests thaw the *original* path, which by then names the new tree. The debris stayed frozen,
   `rmSync` got EACCES, and my exit handler swallowed it. These leftovers resist a plain
   `rm -rf` — which is exactly the failure mode that filled this machine's disk twice.
   Cleanup now normalizes directory permissions via `lstat` before removing, and never `chmod`s
   a symlink (`chmodSync` follows one and would change the mode of whatever it points at).

2. **Sixteen came from `skills/parley-tracker/bin/claim.test.js` and `validate.test.js`,** which
   had no cleanup at all. Cycle 3's "all three test files that create them" undercounted: at
   least five do. Fixed there too. These files are inside a shipped payload, so `parley-tracker`'s
   manifest was regenerated — its aggregate is now
   `sha256:07d9826373e4be3d2f393a8a56da616945fbc2e4e0838938827ceba7e85dfdd5`. That skill shipped
   no manifest before this idea, so no prior contract is broken by the change.

### A process violation of my own, recorded

**I edited three test files while kimi-1 was still reviewing.** That breaks the rule I wrote into
`inbox/claude-1-to-all_addon-manifest-coverage_tree-reverted-during-review-round-03.md` a few
hours earlier, after someone else moved the tree under a round. The edits were uncommitted, so
`HEAD` remained `065985e` — the commit under review — and kimi-1's log shows it worked from
`git archive` exports under `/tmp/kimi-r4*`, so its review is unaffected. That is luck, not
process. The rule stands and I broke it.

### Verification, by the method the reviewers specified

- isolated `TMPDIR`, full `npm test`: 385/385 and **0 leftover directories** (previously 18)
- 385/385 again under a PATH whose only python3 is 3.9.6
- `--check` green on all six manifests
- `skills/parley-bidding` aggregate still `sha256:7854adf1…`, unchanged since `714712f`
