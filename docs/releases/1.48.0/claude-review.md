---
reviewer: claude (release reviewer; not a historical audit signer)
date: 2026-09-18
candidate: CLI head 646254cd39d1261ccaeaf0e92eaaf5fae04c095b plus release edits (1.48.0); skill 2.12.0; tap as checked out
context_mode: full
source_sha256: 4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a
packet_sha256: 4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a
fallback_reason: (none; full mode)
---

# Release review — parley-deck-cli 1.48.0 / parley-deck-skill 2.12.0

This is one reviewer's own review. It is not a signoff, a review consensus, or a ratification of
any disposition, and it does not stand in for the empirical audit. I did not run any tests, builds
or commands. Tests are being run separately, and nothing below claims they pass. No experiment,
pilot or packet-trial results exist in my inputs, and I treat them as absent. Historical releases
1.47.0 and 2.11.0 are unknown history to me (see N2).

## Summary

The release edits are small and consistent with each other. The CLI bootstrap protocol and the
skill's portable snapshot are now byte-identical (both diffs end at blob `c12e523`). The same
LE-7/LE-11 text and the same §9 step 1 appear in the attested live source-role deck. The
`parley protocol packet` flags in SKILL.md match the CLI's flag set, and the dispatch is wired.
The R3 timing window exists exactly as disclosed. It fails closed, and the release notes disclose
it accurately.

I found no CRITICAL or MAJOR issues. I found three MINOR issues and four NITs. The main open item
before announcing the release is the Homebrew tap (m2), which still points at 1.47.0 / 2.11.0.

**Recommendation: GO.** This is conditional on the separately run test suites being green and on
the tap formulas being bumped after the tags are pushed (m2). R3 should stay a disclosed NIT and
should not block the release.

## Inspected scope

- `evidence/review-context.json` and its `body_path` (full attested live source, 1381 lines, read
  completely).
- `evidence/cli-release-diff.patch`: CHANGELOG, VERSION, `internal/app/version.go`,
  `internal/protocol/defaults/COOPERATION.md`.
- `evidence/skill-release-diff.patch`: CHANGELOG, package(-lock).json, SKILL.md,
  parley-addon.json, `references/COOPERATION.md`, `references/compatibility.json`.
- The release notes, `merges.log`, `main-to-release.txt` and `source-scope.txt`.
- `skill/RELEASING.md`, `skill/package.json`, `skill/scripts/build-portable.js` and
  `skill/.github/workflows/release-portable.yml`.
- From `skill/lib/installer.js`: `syncProjectCommand` / `projectStatus` (421–532) and the grep
  map of install/uninstall/atomic-copy entry points. The installer/bin code is **not changed** by
  this release, so I did not re-review the atomic install/remove paths line by line.
- The bytecode guard in `skill/scripts/run-python-tests.js`.
- `tap/Formula/parley-deck-cli.rb` and `tap/Formula/parley-deck-skill.rb`.
- For R3: `implementation-notes/codex-1-r3-disposition-20260916.md`,
  `internal/app/trajectory_verify.go:157-278`, `internal/trajectory/verification.go:296-385`
  and `internal/trajectory/helper_scope.go`.
- `internal/app/protocol_packet.go`, plus greps of `internal/app/protocol.go` and
  `internal/protocol/{workspace.go,drift_test.go,initheader_test.go}`.

**Not inspected:** the ~92k-line integration merge beyond the files above; the Go test bodies;
the add-on manifest hash values (not recomputed, because I ran no commands);
`.gitignore` coverage of `.parley-runtime/`; and the WinGet flow beyond its directory listing.

## Refutation attempts

1. **Protocol/skill drift.** I checked the CLI defaults, the skill snapshot and the live deck for
   diverging close-integrity or §9 text. All three carry the same LE-7/LE-11 bullet and Phase 8
   "Close-decision integrity" paragraph (live body lines 275 and 682–698). They also carry the
   same §9 step 1 (live line 883). Both diffs end at blob `c12e523`, so the CLI and skill files
   are byte-identical. No drift found. `PRIMARY` (diff index lines; attested body).
2. **Skill command vs. the real CLI.** I compared SKILL.md's
   `parley protocol packet --dir --phase --track --idea --json [--flag] [--optimize]` with
   `protocol_packet.go:58-67`. Every flag exists, and `--flag` is repeatable with the same four
   values. The dispatch is present at `internal/app/protocol.go:41-42`. A refusal exits 1, and
   with `--json` it still emits `context_mode: refused` (`protocol_packet.go:87-93`). This matches
   the SKILL's rule that a refusal stops the launch. No mismatch found. `PRIMARY`.
3. **Stale version strings.** I grepped the CLI README/docs/cmd/app and the skill
   README/plugin.json/gemini-extension.json/SKILL.md/agents for 1.47.0, 1.46.0, 2.11.0 and
   2.10.0. There were no hits. The only stale versions are in the tap (m2) and the CHANGELOG gaps
   (N2). `PRIMARY`.
4. **Could the header change in the packaged protocol reach a deck through the installer?**
   `sync-project` writes only `meta/version.json` (`installer.js:459-460`). `projectStatus`
   only hashes the packaged file (`installer.js:495-523`). The installer never writes the
   packaged `COOPERATION.md` into a deck, and `parley init` uses the CLI's own defaults. No leak
   through tooling, but see m1 for manual bootstrap. `PRIMARY`.
5. **npm whitelist.** `files` covers `skills/`, `bin/`, `lib/`, README, LICENSE, NOTICE,
   `plugin.json` and `gemini-extension.json`. It excludes `scripts/`, `test/` and `packaging/`.
   The whitelist is unchanged, and it still covers everything this release touched under
   `skills/parley-deck/`: SKILL.md, the references and parley-addon.json. Stray bytecode is
   guarded against in two ways. `run-python-tests.js:142-145` uses `-B` with
   `PYTHONDONTWRITEBYTECODE`, and `addon-manifest.js:279` rejects undeclared files. No finding.
6. **Portable build / cross-platform.** CI builds on `ubuntu-latest` with Node 24, so the
   `pkg.cmd`-without-shell path in `build-portable.js:9,57` is not exercised in CI. See N4.
   Asset names derive from `package.json` 2.12.0. No release-specific regression.
7. **R3.** See the dedicated section below.

## Findings

### [MINOR] m1 — Portable snapshot header now carries CLI substitution anchors instead of placeholders
`skill/skills/parley-deck/references/COOPERATION.md:3-6` changed from
`<transport-choice> (pick one of local-dir | github-pr | gitlab-mr …)` / `<YYYY-MM-DD>` to
`` `github-pr` `` / `` `<date> — created by parley init` ``. Those strings are the CLI's
`strings.Replace` anchors (`cli/internal/protocol/workspace.go:111,138-142`;
`drift_test.go:79-81`). They are not meant as reader-facing content.

The skill still calls this file a "portability/bootstrap reference", and Appendix A says "Copy
this file… fill in the header". A manual adopter would now get a transport pre-set to
`github-pr` and a false "created by parley init" marker. The workspace field stays a placeholder,
so the header is also internally inconsistent. No tool path writes this file into a deck
(refutation 4), and SKILL.md now forbids using it as launch authority, so the impact is limited to
manual bootstrap.

**Fix:** keep byte-identity everywhere below the header, and have the skill's sync step
re-template lines 5–6. Alternatively, the CHANGELOG could state that the snapshot header shows
CLI anchors. This can be deferred if disclosed.

### [MINOR] m2 — Tap formulas are not bumped (distribution surface incomplete)
`tap/Formula/parley-deck-cli.rb:4-5` still points at `v1.47.0` and `parley-deck-skill.rb:4-5`
still points at `v2.11.0`. That is expected before tagging, because the sha256 must come from the
GitHub tag tarball (`RELEASING.md:106-121`). Until both are updated, Homebrew users stay on the
old versions, and the release is not complete on this surface.

**Fix:** after pushing `v1.48.0` and `v2.12.0`, update `url`/`sha256` from the final tarballs and
run the RELEASING.md Homebrew checks. The existing `post_install` shebang restore and the
install/doctor test should keep working, because the changed files have no shebangs.

### [MINOR] m3 — `compatibility.json` does not express the new CLI dependency
SKILL.md now tells official launches to call `parley protocol packet`, a command that is new in
1.48.0. `references/compatibility.json` records no CLI floor or feature hint. Only
`skillVersion` changed. The SKILL covers older CLIs through the disclosed `full-fallback`, and
the policy there is deliberately "no lockstep", so this is not a blocker. However, `doctor`/`status`
cannot tell a facilitator why every launch shows `full-fallback`.

**Fix (optional):** add an advisory `recommendedCli: ">=1.48.0"` field, or a doctor hint.

### [NIT] N1 — Stale comment in `protocol_packet.go`
`cli/internal/app/protocol_packet.go:16-22` says the dispatch line has not landed and that the
command "is reachable only through this function". The dispatch exists at
`internal/app/protocol.go:41-42`. Under a strict gate, a misleading comment counts as a finding.
**Fix:** delete the paragraph.

### [NIT] N2 — CHANGELOGs skip the versions currently shipped
The CLI CHANGELOG goes from 1.48.0 straight to 1.46.0, but `VERSION` was 1.47.0 and the tap
ships `v1.47.0`. The skill CHANGELOG goes from 2.12.0 straight to 2.10.0, but package.json was
2.11.0 and the tap ships `v2.11.0`. I do not know what 1.47.0 and 2.11.0 contained, and it should
not be reconstructed from memory.

**Fix:** add an explicit "1.47.0 / 2.11.0 — not recorded in this changelog" line, or leave the gap
and state it in the release notes.

### [NIT] N3 — `fallback_reason` is omitted from full-mode JSON
SKILL.md and §9 step 1 say to record `fallback_reason`, but the JSON attestation omits the key in
`full` mode (`review-context.json` has no such key). A strict consumer may treat the missing key
as a malformed attestation. **Fix:** in SKILL.md, say "record `fallback_reason` (absent/empty in
`full` mode)", or always emit the key.

### [NIT] N4 — Pre-existing: `build-portable.js` on Windows hosts
`build-portable.js:9,57` spawns `pkg.cmd` without `shell: true`. Current Node releases reject that
with EINVAL. The script also exits 1 without printing `result.error`. CI uses Linux, so the
release is unaffected, but the "Manual build" in `RELEASING.md:59-64` would fail on a Windows
host. This is not introduced by this release.

## R3 — release blocker or disclosed NIT?

**Verdict: disclosed NIT is appropriate; not a release blocker.** `PRIMARY`, based on the
source locators below.

The window exists as disclosed:

- `trajectory_verify.go:182` runs `checkTrajectoryHelperScope` outside any guard.
- The protocol precheck follows at `:212`.
- `trajectory.PrepareCapturedVerification` at `:217` writes the ticket under `withState`
  (`verification.go:310-329`). Its guarded callback checks only
  `checkCapturedActivationQuorum` (`:315`). It does not receive `req.Participants`/`Criteria`,
  and it does not run `checkReconciliationScope`, which only `CheckHelperScope` does
  (`helper_scope.go:34-46`).

A concurrent participant or criteria edit between `:182` and `:217` can therefore consume a
ticket.

It fails closed:

- The helper rechecks scope before executing (`trajectory_verify.go:344`).
- Acceptance rechecks it before reading the receipt (`:257`).
- `withVerificationAuthority` rechecks the prepared ticket's authority (`verification.go:360-385`).

The consequence is an extra retained ticket and spent attempt plus a later refusal. That is a
cost and diagnostics issue, not a path to unauthorized acceptance. The CLI CHANGELOG
"Release scope and limitations" paragraph states this accurately.

The disposition note itself says that the four-person review has not ratified it. My concurrence
is one reviewer's opinion and does not supply that ratification. I agree with the suggested
future fix: pass the full `HelperRequest` into preparation, recheck scope before the ticket
write, and add a barrier-controlled test. It narrows the window but cannot fence file edits
outside the guard.

Out of scope for R3 and not raised: between the recheck at `:257` and `ReadCapturedVerification`
at `:260`, there is a similar unguarded gap. I did not trace whether later close consumers
recheck scope, so I make no claim either way.

## Open questions

1. Is `.parley-runtime/` ignored in consumer repositories? SKILL.md now has every facilitator run
   the renderer against `<project-root>`, which writes bodies under
   `<project-root>/.parley-runtime/protocol-packets/`. I did not verify the ignore rules outside
   the CLI repo.
2. Were the parley-addon.json hashes regenerated by `npm run manifest:addons` rather than by
   hand? The separately run `npm test` / `prepack --check` is the evidence for this. I did not
   recompute them.
3. The latest WinGet manifest in `packaging/winget/` is 1.3.1. Is WinGet still a maintained
   distribution surface? If it is, it has been stale since before this release.

---

## Follow-up: final-diff assessment (2026-09-18)

**Scope:** `evidence/final-release-edits.patch` only. I also inspected:
- the modified build script's resolution target: `skill/node_modules/@yao-pkg/pkg/package.json`,
  `lib-es5/bin.js` and `skill/package-lock.json:652-654`;
- `skill/skills/parley-deck/references/COOPERATION.md:1-8`;
- CLI `internal/protocol/defaults/COOPERATION.md:268,681`;
- a skill-wide grep for `recommendedCli`;
- the `IMPLEMENTATION.md` frontmatter.

I ran no tests, builds or commands. The operator reports 391 Node and 54 Python tests passing and
the Go suite still running. That is their report, and I do not certify it. The brief does not
state the manifest `--check` result separately, and I did not recompute the new hashes.
Structurally, the three changed hashes in `parley-addon.json` match the three changed payload
files: SKILL.md, `references/COOPERATION.md` and `compatibility.json`. The unchanged files keep
their hashes.

### Disposition of my earlier findings

- **m1 — resolved.** `references/COOPERATION.md:3-6` has the portable placeholders again. The
  skill diff now contains only the three body hunks, and their text is identical to the CLI body
  change. The CLI defaults still carry the fail-closed text (`:268`, `:681`). I compared hunks
  and did not run a byte diff. On that basis, the snapshot and the CLI defaults now differ only in
  the header lines 5–6, by design.
- **m2 — open, pending post-tag.** This is still a release-completeness gate. I cannot verify it
  until the archives are published.
- **m3 — resolved as advisory.** `recommendedCli: ">=1.48.0"` is at `compatibility.json:11`. No
  code in lib, bin or tests reads it; it only appears as raw data under
  `packaged.compatibilityManifest` in `status` output (`installer.js:514-520`). That is acceptable
  under the no-lockstep policy. A `doctor` hint remains an optional follow-up.
- **N1 — resolved.** The hunk removes only the stale paragraph.
- **N2 — resolved.** Both CHANGELOGs now carry 1.47.0 / 2.11.0 entries that link the published
  release pages. They say the detail was not recorded, rather than reconstructing it. I could not
  verify the `2026-08-29` dates or the link targets because I have no network access, so those
  rest on the operator's word.
- **N3 — resolved.** Minor caveat, not a finding: in `refused` mode the same field carries the
  refusal reason (`protocol_packet.go:87-88`). This is harmless because a refusal stops the
  launch.
- **N4 — resolved.** `build-portable.js` now runs `process.execPath` with
  `require.resolve("@yao-pkg/pkg/lib-es5/bin.js")`, so there is no `.cmd` file, no shell and no
  EINVAL. I confirmed the following from source:
  - the lockfile pins `@yao-pkg/pkg` 6.19.0;
  - its `package.json` declares `"bin": {"pkg": "lib-es5/bin.js"}` and has no `exports` map, so
    the subpath resolves;
  - `bin.js` calls `exec(process.argv.slice(2))` on load and exits 2 on error, so arguments and
    failure status propagate;
  - `result.error` is now printed, and a null status still exits 1.

  The operator reports rerunning the portable builds; I did not verify that.
- **Open question 3 (WinGet) — answered by operator report.** The 1.47.0 and 2.11.0 submissions
  were merged, and the final manifests will use remote assets, which matches `RELEASING.md:83`.
  I have not verified this.

### New observations on the final edits

#### [NIT] F1 — The build script depends on pkg's internal file layout
`skill/scripts/build-portable.js` hard-codes `lib-es5/bin.js`. This release is safe because the
lockfile pins 6.19.0. However, a future bump within `^6.19.0` that moves the file would make
`require.resolve` fail at startup. That failure would be loud, not silent.

**Optional hardening:** resolve `@yao-pkg/pkg/package.json` and join its `bin.pkg` path.

#### [NIT] F2 — The release-preparation note paraphrases the user's instruction
The new `IMPLEMENTATION.md` section "User-directed release preparation" (patch lines 92–106) was
added by the file's owner (`implementer: codex-1`, frontmatter line 4), so ownership is fine. Its
content matches my own brief and keeps the experiments and audit open. However, it paraphrases
the user's instruction instead of quoting it with a date and medium. The instruction "supersedes
the earlier no-merge/no-release scope restriction". §4 requires user answers to be quoted
verbatim; applying that standard here would give a stronger audit record.

The note also commits a machine-local `/private/tmp/...` checkout path, which will go stale.
Neither issue affects the release artifacts.

### Updated recommendation

**GO. The final edits add no new CRITICAL, MAJOR or MINOR issues.** This remains conditional on:
- the Go suite finishing green (operator-run);
- the manifest check being part of the reported passing run;
- the tap being updated after tagging (m2).

R3 remains a disclosed NIT and does not block the release. My concurrence is still one
reviewer's opinion, not a ratification.

---

## Follow-up 2: two final corrections (2026-09-18)

**Scope:** only `skill/skills/parley-deck/SKILL.md:8-60` and the fixture in
`cli/internal/app/trajectory_verifier_recovery_test.go` around `canonicalRoot` (lines 88–130,
plus every use of `runBase()`). I ran no commands. The operator reports that
`TestAppVerifierRecoveryFullLifecycle` passes and that the full suite and race checks are
rerunning with a 30-minute timeout. Those are their results, and I do not certify them.

### SKILL.md — accepted

- The Core Rule (`SKILL.md:12`) now points to **Required Protocol Context** and "the resolved
  authority". This removes the contradiction with the renderer-first flow. Nothing in lines
  8–60 still tells the facilitator to read `parley-deck/COOPERATION.md` first.
- Lines 46–48 say that a reachable renderer failing without an attestation stops the launch.
  This matches the CLI:
  - `BuildProtocolContext` errors exit 1 with no attestation (`protocol_packet.go:83-86`);
  - usage errors exit 2 (`:68-74`).

  It also matches §9 step 1: a launch with no attestation is unresolved, and only an unreachable
  renderer allows `full-fallback`. The rule fails in the safe direction.
- The manifest was regenerated along with the edit. The SKILL.md hash at `parley-addon.json:4`
  changed from `a2b02871…` to `c7aa31a2…`, and the aggregate changed to `9f0d4eb7…`. I did not
  recompute these; the manifest `--check` is the evidence.

#### [NIT] G1 — No stated way to tell "older CLI" from "reachable failure"
The text now separates two cases:
- the renderer is unreachable (for example, an older CLI without the command), which leads to
  `full-fallback`;
- a reachable renderer fails without an attestation, which stops the launch.

To the caller, both look the same: a nonzero exit and no attestation. If a facilitator treats an
authority or I/O error as "older CLI", they would fall back to reading unattested text, which is
exactly what the new sentence forbids.

**Suggested fix:** state the test explicitly. The renderer is unreachable only if `parley` is not
on PATH or `parley version` is below `recommendedCli` (1.48.0). Any failure from a 1.48.0+ CLI
stops the launch. This does not block the release. Line 49 also wraps short in mid-sentence,
which is cosmetic because it still renders as one paragraph.

### Test fixture — accepted; a test-only fix is the right scope

- **Production canonicalizes paths before persisting them.** `verifyTrajectoryWithAgent`
  resolves the root with `filepath.Abs` and `filepath.EvalSymlinks`
  (`trajectory_verify.go:162-164`) before building `RequestPath` (`:196`).
  `withVerificationAuthority` also rejects a `ticket.Root` that is not canonical
  (`verification.go:351-352`).
- **The old expectation was the bug, not production.** On macOS, `t.TempDir()` returns a path
  under `/var/…`, while production persists `/private/var/…`, so only the test's string
  comparison was wrong.
- **The fix covers every affected comparison.** The fixture now builds the expected `base` from
  `EvalSymlinks(root)` (test lines 91–97). The only string comparison against a persisted
  absolute path (`:108`) uses that value.
- **`runBase()` is unaffected.** It still uses the original root (`:128-130`), but every use is
  file I/O (`:146, 343, 444, 483, 528, 539`), where both spellings reach the same file.
- **Production canonicalization stays under test.** Later calls still pass the uncanonicalized
  `f.root` into production APIs.
- **Portability.** On a path with no symlinks, such as a typical Linux `/tmp`, `EvalSymlinks`
  returns the path unchanged. On Windows it applies the same normalization production applies,
  so the two stay aligned. No assertion is weakened, and no production code changed.

### Updated recommendation

**Unchanged: GO.** This remains conditional on:
- the rerun Go suite and race checks finishing green (operator-run);
- the manifest check passing;
- the tap being updated after tagging (m2).

R3 remains a disclosed NIT.
