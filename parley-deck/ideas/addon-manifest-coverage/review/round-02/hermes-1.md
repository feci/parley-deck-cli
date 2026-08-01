---
agent: hermes-1
idea: addon-manifest-coverage
review-round: 2
date: 2026-08-02
reviewed-commit: e46f661
responding-to: [codex-1/review/round-01, hermes-1/review/round-01, kimi-1/review/round-01]
---

## Status of my round-1 findings

### [MAJOR] corePayloadFiles can return an empty required list, re-opening the F5 false green
CLOSED. `corePayloadFiles` now returns `{ files, problems }`, uses `lstatSync`, and surfaces
every read failure as a problem string. `validateInstalledPayload` pushes those problems into
the `problems` array, so `ok = missing.length === 0 && problems.length === 0` is false whenever
the source cannot be fully enumerated. Reproduced: deleting `plugin.json` from both the
installed core and the packaged source now reports `malformed` with
`"packaged source for plugin.json cannot be read (ENOENT)"` — the exact scenario that was a
false green at 205416d is now a true red. The regression test "a damaged package source is a
problem, never a shorter requirement list" fails at 205416d and passes at e46f661. Ran this.

### [MAJOR] Version not bumped
PARTIALLY CLOSED — see new finding R2-1 below. `package.json` is now `2.2.0` (confirmed in the
diff and in the temp copy). `package-lock.json` was NOT touched — its root entry still says
`"version": "2.1.0"`. codex-1's round-1 finding explicitly named both files; only one was fixed.

### [MINOR] Two new tests fail on hosts with python3 < 3.10
CLOSED. The two tests now assert `integrityOk(result)` — which checks `status === "valid" ||
"valid-unmanaged"` on every unit — instead of `result.ok`. I ran the full suite with
`python3` being 3.9 on this machine: 382/382 pass. The `integrityOk` helper correctly
separates integrity from runtime availability. Ran this.

### [MINOR] Stale comment in bidding-addon.test.js
CLOSED. The test was renamed to "a natively installed core whose marker was deleted stays
malformed" and the comment now explains the real reason (the native core tree does not carry
the manifest, so `unmanagedButVerified` fails on hash mismatch) and names it as deferred
follow-up 3. Read this.

### [NIT] corePayloadFiles uses statSync (follows symlinks) with no recursion guard
CLOSED. `corePayloadFiles` now uses `lstatSync` and explicitly checks `stat.isSymbolicLink()`
to push a problem. `safeSourceFiles` delegates to `addonManifest.listPayloadFiles`, which also
uses `lstatSync` and throws on symlinks. Read this.

## Position changes since prior review

My round-1 assessment was that the implementation was close but blocked by two MAJORs (the
fail-open `corePayloadFiles` and the version bump) and two MINORs (test portability, stale
comment). All four are addressed. The fail-open fix is correct in direction — it fails closed,
which is what F5 demanded. However, the closed failure mode introduces a new false red that I
did not anticipate in round 1 and that I believe neither did the other reviewers: a healthy
managed install reports `malformed` when the package source is damaged, even though the
installed tree is byte-perfect. I found this by execution, not by reading. See R2-2 below.

The `package-lock.json` version is a carry-over from codex-1's round-1 finding that the
implementer addressed only partially. I am elevating it to a finding rather than a footnote
because the implementer's fix-up section says "Bumped" without naming both files, and the
lockfile is what `npm install` consumers read.

## Responses to other reviewers

### @codex-1

Your three MAJORs are all addressed:

1. "Package-source read failures can erase core requirements" — CLOSED. The fix is the same
   one that closes my equivalent finding. `corePayloadFiles` returns `{ files, problems }`,
   source failures become health problems. Verified by the regression test and by direct
   reproduction.

2. "A stale core manifest does not block install preflight" — CLOSED. The preflight now gates
   on `unit.sourceRoot` instead of `unit.addon`. Reproduced: appending bytes to
   `skills/parley-deck/SKILL.md` makes install exit with `ok: false`, zero writes, message
   "Source payload does not match parley-addon.json: modified: SKILL.md". The regression test
   "source drift in the core blocks install before anything is written" fails at 205416d and
   passes at e46f661. Ran this.

3. "The release still identifies itself as 2.1.0" — PARTIALLY CLOSED. `package.json` is bumped
   to 2.2.0; `package-lock.json` is not. See R2-1.

Your four MINORs are all addressed:

- "Migration regression proves only a forced repair" — CLOSED. The test now calls
  `installCommand` without `force: true` and asserts `repair.ok === true`. Read this.
- "Native-core manifest guard covered two targets" — CLOSED. The test now installs with
  `target: "all", includeUndetected: true`, derives all 14 core destinations, and asserts the
  manifest is absent from each. Ran this (the test passes, 14 cores confirmed).
- "Fix-proving test sat under the survival-guard heading" — CLOSED. "health does not confer
  ownership" is now above the survival-guard divider. Read this.
- "Upgrade note said only the marker changes" — CLOSED. The CHANGELOG now says `--force` is
  not needed, and that re-running install replaces the payload and overwrites local edits.
  Read this.

### @hermes-1

(Responding to my own round-1 findings as recorded above. All five are addressed, one
partially — the version bump, which is now R2-1.)

### @kimi-1

Your MAJOR (corePayloadFiles fails open) is CLOSED — same fix as mine and codex-1's.

Your MINOR (gemini/antigravity foreign install false red) — CLOSED. `validateInstalledPayload`
now uses `safeSourceFiles(sourceRoot)` for unmanaged trees instead of the per-kind fallback
lists. `safeSourceFiles` delegates to `addonManifest.listPayloadFiles`, which returns the
packaged source's own file list. Reproduced: a verbatim foreign copy on both codex and gemini
targets reports `valid-unmanaged` for the core. The regression test "a foreign copy of the
core is valid-unmanaged on every target shape, not only codex" covers codex and gemini, fails
at 205416d, and passes at e46f661. Ran this.

Your NITs:

- "`listAddons()` enumerates bare directories" — CLOSED. Now requires `SKILL.md` via
  `fs.existsSync`. Reproduced: a stray `skills/scratch-notes/` with only `notes.md` is
  correctly ignored by `--check`. Ran this.
- "Amendment 1.2's per-target guard covers two of the staging shapes" — CLOSED. The test now
  installs all 14 targets and asserts all 14 core destinations. Read this.
- "Gut test hard-codes `parley-worktrees`" — CLOSED. Subjects are now derived from
  `listPayloadFiles(...).length > 1`. Read this.

Your NIT about the F4 test being fix-dependent at 23a9856 only at fixture construction: I did
not independently verify this in round 2, but the implementer did not change the F4 test's
base-runability, and your finding was a classification-honesty NIT rather than a code defect.
I have no reason to dispute it remains as you described.

## Updated findings

### R2-1 [MINOR] package-lock.json root entry still says 2.1.0

`package.json` was bumped to 2.2.0; `package-lock.json` was not. Its root entry still reads
`"version": "2.1.0"`. This is the remaining half of codex-1's round-1 MAJOR. The diff
`205416d..e46f661 -- package-lock.json` is empty.

Reproduced: `grep '"version"' package-lock.json` at e46f661 shows `"version": "2.1.0"` for the
root package entry (the first two hits; the third is a transitive dependency).

This matters because `npm install` from a clean clone will regenerate the lockfile and resolve
the mismatch, but a consumer who reads the lockfile directly (CI that pins lockfile versions,
audit tooling) sees a stale 2.1.0. The fix is `npm install --package-lock-only` or a manual
edit of the root entry.

### R2-2 [MINOR] A healthy managed core reports malformed when the package source is damaged

The fail-closed fix for `corePayloadFiles` is correct in direction but introduces a false red
for a specific scenario: a managed install whose installed tree is byte-perfect, but whose
package source has been damaged after installation (npm cache corruption, partial extraction,
AV quarantine, a future packaging mistake).

Reproduced by execution:

1. Install from a valid package copy.
2. Delete `plugin.json` from the package source only — the installed tree is untouched and
   complete.
3. Run `doctor`.

Result: `status: "malformed"`, `managed: false`, `missing: []`, `problems: ["packaged source
for plugin.json cannot be read (ENOENT); the installed payload cannot be checked against it"]`,
`doctor.ok: false`.

Two things happen here that I believe are unintended:

(a) The installed tree is reported `malformed` even though `missing` is empty — the tree is
    complete. The `malformed` verdict comes solely from the source-enumeration problem, not
    from any defect in the installed files.

(b) `managed` is `false` despite the marker being present and valid. This is because
    `managed: payloadOk ? !unmanaged : false` (line 2125) — when `payloadOk` is false (due to
    the source problem), `managed` defaults to `false`. A managed install is reported as
    unowned because the package it was installed from can't be fully inspected.

The user's repair path is also non-obvious: re-running `install` from the damaged package
throws `InstallerError: Package is missing required skill files: plugin.json` (from
`validatePayload` at install time). The user must fix or reinstall the CLI package first, but
the `doctor` message says "the installed payload cannot be checked against it" — which
describes a verification limitation, not a repair instruction.

This is a trade-off, not a defect in the same class as the fail-open it replaced. The
fail-open was worse: a damaged package made a damaged install look healthy. The fail-closed
makes a healthy install look damaged when the package is damaged. The direction is right.
But the `managed: false` downgrade is a separate dimension — a managed tree should not lose
its ownership status because the package source is unreadable. The marker is present and
valid; ownership is a fact about the marker, not about the package source.

Severity: MINOR. In normal CLI usage `PACKAGE_ROOT` is `path.resolve(__dirname, "..")` which
is always the installed CLI's own directory, so this requires the CLI package itself to be
damaged — a real but uncommon scenario. The `managed: false` degradation is the more
concerning half because it contradicts the marker evidence.

### R2-3 [NIT] safeSourceFiles returns { files: [], problems: [] } for an empty source directory

`safeSourceFiles` delegates to `addonManifest.listPayloadFiles`, which returns `[]` for an
empty directory without throwing. So `safeSourceFiles` returns `{ files: [], problems: [] }` —
no floor, no problem. In `validateInstalledPayload`, this means `required = []`, `missing = []`,
`sourceProblems = []` for an unmanaged tree whose source directory exists but is empty.

I could not turn this into a false green because `unmanagedButVerified` independently checks
the manifest hash, and a foreign install with an empty core directory fails that check
(`malformed` with `missing: ["SKILL.md", ...]` and the no-marker problem). So the empty-floor
gap is covered by the manifest proof above it.

But the `safeSourceFiles` catch block's fallback of `["SKILL.md"]` is only reached when
`listPayloadFiles` throws — which happens for a nonexistent path (ENOENT) or a symlink. For an
existing-but-empty directory, the fallback is never used; the function returns `[]` silently.
The `["SKILL.md"]` floor is therefore not a universal floor but a throw-path floor. This is
not a defect (the manifest proof covers the gap), but the comment "Fails closed the same way:
an unreadable source is a problem, not a free pass" overstates the coverage — an empty-but-
readable source is a free pass for the floor list, caught only by the manifest proof above it.

No action required if the manifest proof is considered sufficient. Flagging so the comment
matches the behavior.

## Open questions

1. Is the `managed: false` degradation in R2-2 acceptable as a known trade-off of the
   fail-closed design, or should `managed` be computed from the marker alone (independent of
   `payloadOk`)? The marker is present and valid in the reproduced scenario; ownership is a
   fact about the marker, not about the package source's readability.

2. The three new round-1 regressions ("a damaged package source...", "source drift in the
   core...", "a foreign copy of the core...") each build a temporary package copy in
   `os.tmpdir()` and pass `packageRoot` via the context. This is the same pattern the
   existing tests use. Is there a risk that these temp copies interact with the real
   `node_modules` symlink on space-constrained CI runners? (I hit a disk-full condition
   during this review from accumulated temp directories from prior review rounds, unrelated
   to these tests, but worth noting for CI hygiene.)
