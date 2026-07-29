---
idea: skills-cli-install-path
review-round: 03
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at
`bddbf1af908997a0fa6ed5976038cd2ad4301723`. The Gemini reconciliation, all five prior
findings, the stale-path sweep, and its new regression guard reproduce. Native destination
shapes remain unchanged, `npm test` and packaging pass, and G8 is testable on this host and
passes.

One shipped tracker exemplar still contains a broken test command. It claims to verify a
non-zero test suite but instead fails with `MODULE_NOT_FOUND` and zero passing tests. That
remaining documentation/verification defect should be fixed before merge.

## Prior finding dispositions

### FIXED — round 01 [MAJOR] D-1 removed a reconcilable Gemini channel

The implementer withdrew D-1 and implemented the proposed two-consumer layout:

- Repository `gemini-extension.json` has
  `contextFileName: "skills/parley-deck/SKILL.md"`, and that path exists in the repository.
- A native Gemini install stages `contextFileName: "SKILL.md"`, and that path exists in the
  installed destination.
- Removing `contextFileName` from both parsed manifests leaves semantically identical
  objects; the rewrite did not alter the other manifest fields.
- Changing the repository value to `SKILL.md` in a scratch copy made the focused repository
  test fail with exit 1 and show the wrong actual/expected values.
- Changing the staged rewrite to `BROKEN.md` in a separate scratch copy made the focused
  staged-install test fail with exit 1 and show the wrong actual/expected values.

Both tests can fail and guard their respective value. The README restores
`gemini extensions install <repo-url>` and says plainly that the Gemini CLI has not been run
end to end. `gemini` is absent on this host, so G11 remains NOT TESTED rather than being
misreported as a pass. Reconciling the two consumers was the right decision.

### FIXED — round 01 [MINOR] G1 was recorded as an unconditional pass

The fix-up record now qualifies the clean-HOME ordering behavior. Fresh isolated runs
reproduced it:

- At `bddbf1a`, the first exact `install --target all --force` exits 0, `status --target all`
  exits 0 with five Gemini units missing, and `doctor --target all` exits 1.
- The same sequence at pre-move parent `94a4889` gives the same exit codes and missing Gemini
  units. It is not a layout regression.
- A second identical install detects and installs Gemini; status and doctor then report all
  25 core/add-on units valid and exit 0.

### FIXED — round 01 [MINOR] Installed README named checkout-only paths

The README now distinguishes repository checkout paths under `skills/parley-deck/` from the
flat `SKILL.md` and `references/COOPERATION.md` paths in native destinations. Both shapes
exist in the places the README assigns them.

### FIXED — round 01 [MINOR] Universal-agent wording exceeded measured scope

The panel now scopes support to coding agents supported by the upstream installer, attributes
that support surface upstream, and states no numeric agent count. It does not say
“whichever coding agents you have.”

### FIXED — round 02 [MAJOR] Shipped instructions still referenced the removed `addons/` tree

`grep -rn "addons/" skills/` returns no matches. There is no repository `addons/` directory,
no compatibility tree or symlink, and the package contains no `addons/` entry.

The cycle-3 guard is genuine. In a scratch copy I changed
`skills/parley-tracker/SKILL.md:79` back to an `addons/parley-tracker` path. `npm test`
then exited 1 with 250 passes and one failure, and its assertion named
`skills/parley-tracker/SKILL.md:79`. The unmodified tree returns 251 passes.

The concrete repository-relative commands swept in cycle 2 were also exercised:

- `node --test "skills/parley-design-check/test/*.test.js"` runs 159 tests, all passing.
- The design-check CLI path resolves and its help command exits 0.
- `node skills/parley-tracker/bin/validate.js --strict --dir
  skills/parley-tracker/templates` validates all three exemplars.
- In a constructed `tickets/` tree, the epic exemplar's strict-directory and single-file
  validation commands both pass.
- The story exemplar's `claim` command succeeds in a scratch copy and writes
  `status: in-progress` plus the requested assignee.

The remaining tracker test command is a separate finding below; it is not a stale
`addons/` reference and does not invalidate the fixed disposition above.

## Refutation attempts and gate results

### G7 — actual install without `--full-depth`

I copied `bddbf1a` to scratch and ran the exact requested spelling:

`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy`

with no `--full-depth`. The npm wrapper failed before the CLI ran because
`registry.npmjs.org` could not resolve (`ENOTFOUND`). I therefore do not claim that the
network-backed `@latest` resolution passed.

Using the locally cached `skills` 1.5.20 executable—the same version used by the ratified
design evidence—with the same add arguments and no `--full-depth` produced `Found 5 skills`
and `Installed 5 skills`. Assertions found exactly five installed `SKILL.md` roots. The core
contained `SKILL.md`, `agents/`, and `references/`, while `bin/`, `lib/`, and `package.json`
were absent. The layout behavior passes; only the unavailable registry wrapper remains
unreproduced.

### Native destination shapes and validation

`validateInstalledPayload` is unchanged from `94a4889` and still validates
destination-relative paths. After the second isolated install, the actual relative path sets
were:

- Codex: marker, `LICENSE`, `README.md`, `SKILL.md`, `agents/` with both manifests, both
  repository manifests, and `references/` with all three references.
- Antigravity: the Codex set plus the fabricated `skills/SKILL.md`.
- Gemini: the Codex set, with staged `contextFileName: "SKILL.md"`.

Each of those three path sets is exactly identical to the corresponding install from
`94a4889` (`diff` exit 0). `agy plugin validate` against the installed Antigravity
destination exits 0 with one skill and two agents processed.

### Tests and package

- `npm test`: 251 passed, 0 failed.
- `npm pack --dry-run --json`: 153 files, 145 under `skills/`, exactly five
  `skills/<name>/SKILL.md` roots, no `addons/` entries, and no root `SKILL.md`.

### Previously NOT TESTED gates

- **G3 Windows:** PARTIALLY TESTED. The normal build script could not replace a stale,
  read-only pkg cache entry in this sandbox. Pointing each build explicitly at the existing
  cached x64 and ARM64 Node runtimes produced both v1.5.0 executables; `file` identifies them
  as PE32+ x86-64 and PE32+ Aarch64. Windows execution/install remains NOT TESTED because
  neither Windows nor Wine is available.
- **G8 `skills update`: PASS.** I installed all five skills from an isolated
  `file://...source.git` repository, committed a new core reference file to that source, and
  ran `skills update --yes`. It reported all five updated, and the new reference appeared in
  the installed core. A Git remote makes this gate testable before publication; a plain
  filesystem install still records `sourceType: local` and is intentionally not updateable.
- **G9 Homebrew:** NOT TESTED against `bddbf1a`. Homebrew and version 1.5.0 are installed, but
  the tapped formula downloads tag `v1.5.0` at `978fce0`, not this review head, and this
  repository's `packaging/homebrew/Formula/` is empty. `brew test` additionally could not run
  in this sandbox because Homebrew attempted network and vendor-cache writes. An upgrade
  would test different source and mutate the host, so it was not treated as branch evidence.
- **G10 WinGet:** NOT TESTED. `winget` is unavailable on this macOS host; the checked-in
  manifests also target earlier release versions, not `bddbf1a`.

### README panel and implementation claims

The panel satisfies A4/F1/F3/F4/F5: neither path is labelled recommended, no numeric
universal-agent count is claimed, `skills list` is present, the native install/doctor path is
immediately visible, and the all-five statement is backed by the cached G7 run. F4's
core-only fallback wording is not triggered because the restructure and five-skill install
remain present.

The current passing claims in `IMPLEMENTATION.md` reproduce, including cycle 3's negative
guard proof, except for the explicitly separated network-backed `npx ...@latest` resolution.
The initial G1 table is not independently true in a clean HOME, but fix-up cycle 1 already
supersedes it with the reproducible precondition and pre-existing ordering explanation. G8 is
stronger than the implementation record: it is testable with a Git remote and passes.

## Findings

### [MAJOR] A shipped tracker exemplar's test command fails and runs zero passing tests

`skills/parley-tracker/templates/subtask.md:68` and `:74` claim the subtask's tests are
verified by:

`node --test skills/parley-tracker/bin`

I ran that exact command from the repository root. Node v26.5.0 treats the directory as the
test entrypoint, fails with `MODULE_NOT_FOUND`, and reports one failed harness entry with
zero passing tests. This is a filled exemplar whose acceptance criterion and Definition of
Done tell implementers that the command proves tests are green; it currently proves nothing
and cannot be checked off honestly.

Replace both occurrences with a command that names the test files, for example
`node --test "skills/parley-tracker/bin/*.test.js"`. I ran that form and it executes 35
tests, all passing. Add a focused regression test that runs the exact command published by
the exemplar and asserts both exit 0 and a non-zero test count; the `addons/` text guard
cannot detect this class of broken verification command.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
