---
idea: skills-cli-install-path
review-round: 02
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at `085799e`.
Fix-up cycle 1 resolves all four findings from my round-01 review, including the
Gemini manifest reconciliation and regression tests. I nevertheless block this
revision because the layout move left shipped, executable instructions pointing
at the removed `addons/` tree. The primary design-check command now fails, and
its documented test command exits successfully after running zero tests.

## Round-01 finding dispositions

### FIXED — [MAJOR] D-1 removed a reconcilable Gemini channel

The deviation is withdrawn and the proposed two-consumer layout is implemented:

- Repository `gemini-extension.json` has
  `contextFileName: "skills/parley-deck/SKILL.md"`, and that file exists in the
  repository.
- A native Gemini install rewrites only that field to `SKILL.md`, and
  `<destination>/SKILL.md` exists. The other manifest fields remained
  semantically identical.
- I broke each side independently in scratch copies. Changing the repository
  value to `SKILL.md` made the repository-manifest test fail with exit 1.
  Changing the staged rewrite to `BROKEN.md` made the staged-install test fail
  with exit 1. These tests can fail and guard the intended reconciliation.
- The README restores `gemini extensions install <repo-url>` and explicitly
  says that the Gemini CLI has not been run end to end. That is honest.

### FIXED — [MINOR] G1 was recorded as an unconditional pass

The fix-up now states the precondition and records the Antigravity-to-Gemini
detection ordering as a follow-up. In a clean scratch HOME I reproduced the
exact sequence:

- first `install --target all --force`: exit 0;
- `status --target all`: exit 0, with five Gemini units missing;
- `doctor --target all`: exit 1.

I ran the same sequence at parent `94a4889` and got the same result, including
the missing Gemini units and doctor exit 1. This is not a regression. A second
exact install detects Gemini; status and doctor then report 25/25 units valid.

### FIXED — [MINOR] The installed README named checkout-only paths

The README now distinguishes
`skills/parley-deck/{SKILL.md,references/COOPERATION.md}` in a checkout from
`SKILL.md` and `references/COOPERATION.md` in an installed destination.

### FIXED — [MINOR] “Whichever coding agents you have” exceeded measured scope

The panel now scopes the claim to coding agents supported by the upstream tool,
attributes that support surface to upstream, and states no count. The all-five
claim remains backed by the real G7 install.

## Refutation attempts and gate results

### G7 — real install without `--full-depth`

I copied commit `085799e` to scratch and ran the exact
`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy`
command with no `--full-depth`. The npm wrapper could not resolve
`registry.npmjs.org` (`ENOTFOUND`) before the CLI ran. Using the locally cached
`skills` 1.5.20 CLI—the same version used by the design evidence—with otherwise
identical arguments produced `Found 5 skills` and `Installed 5 skills`.

Five installed `SKILL.md` files were present. The core contained `SKILL.md`,
`agents/`, and `references/`; `bin/`, `lib/`, `package.json`, `test/`, and
`dist/` were absent. The layout behavior passes; the network-backed `@latest`
wrapper itself is not reproducible in this sandbox.

### Native destinations, validation, and Antigravity

After the second exact install described above, actual core destinations were:

- Codex: marker, `SKILL.md`, `agents/`, `references/`, both manifests, README,
  and license.
- Antigravity: the Codex shape plus fabricated `skills/SKILL.md`.
- Gemini: the Codex shape, with staged `contextFileName: "SKILL.md"`.

Relative path sets for Codex, Antigravity, and Gemini were byte-for-byte
identical to installs made from parent `94a4889` (13, 15, and 13 entries,
respectively). `validateInstalledPayload` still validates destination-relative
paths. `agy plugin validate` passed with one skill and two agents processed.

### Tests and package

- `npm test`: 249 passed, 0 failed.
- `npm pack --dry-run --json`: 153 entries; 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` roots; zero `addons/` entries; no root `SKILL.md`.
- The initial pack attempt hit an unrelated unwritable npm-cache error; rerunning
  with a scratch npm cache produced the result above.

### Previously untested platform/update gates

- G3: both Windows x64 and ARM64 binaries cross-built successfully and were
  identified as PE32+ executables. Windows execution/install remains NOT TESTED
  because neither Windows nor Wine is available.
- G8: PASS. I installed all five skills from an isolated `file://` Git remote,
  advanced that scratch remote with a new core reference file, and ran
  `skills update --yes`. It reported all five updated, and the new file appeared
  in the installed core.
- G9: NOT TESTED against this revision. Homebrew is installed, but its formula
  consumes the published tag tarball, not unmerged commit `085799e`; running
  `brew upgrade` or `brew test` would test different source.
- G10: NOT TESTED because `winget` is unavailable on macOS.
- G11: NOT TESTED because the Gemini CLI is not installed. Both static manifest
  resolutions pass, but they do not substitute for the end-to-end CLI gate.

The pass claims in the fix-up reproduce except for the network-backed
`npx ...@latest` resolution noted above. G8 is stronger than the implementer's
record: it is testable with a Git remote and passes.

### README panel against A4/F1/F3/F4/F5

The panel has no “recommended” label, makes no numeric agent-count claim,
attributes the supported-agent surface upstream, includes `skills list`, and
claims all five only after a real install found and installed all five. F4 is
not triggered because the move remains in place. I found no remaining
unsupported claim in the panel.

## Responses to other round-01 reviewers

### @hermes-1

I agree that local G7 is strong evidence and the published-remote run remains a
post-merge gate. I no longer agree with the round-01 conclusion that the Gemini
URL channel had to be dropped: commit `085799e` implements and verifies the
two-manifest reconciliation without changing installed destination shape.

### @kimi-1

I independently reproduced the stale `addons/` references reported in the
late round-01 review. The finding below is based on my own commands at
`085799e`. The prior conclusion that one Gemini manifest value could not serve
both consumers is superseded by the verified staged rewrite.

## Findings

### [MAJOR] The layout rename left shipped run and test instructions pointing at the removed `addons/` tree

There are 34 `addons/` references across six files under `skills/`, although
the package deliberately ships no `addons/` directory. This is not merely a
historical label:

- `skills/parley-design-check/SKILL.md:49` tells checkout users to run
  `node addons/parley-design-check/bin/check.js ...`. I ran that path and got
  `MODULE_NOT_FOUND` (exit 1); the corresponding `skills/...` command works.
- Its documented test command at line 372 exits 0 while running **zero tests**
  because the removed path matches nothing. The corrected `skills/...` command
  runs 159 tests. This creates a false-green verification path.
- Lines 383–390 list eight nonexistent checkout paths.
- `skills/parley-tracker/SKILL.md:79-83` and the three shipped ticket templates
  continue to prescribe `addons/parley-tracker/...` paths; executing one of the
  template validator commands fails with `MODULE_NOT_FOUND`.

`npm test` does not detect this because the moved files themselves are valid;
the defect is in their shipped instructions. Update checkout-relative paths to
`skills/<name>/...`, state installed-directory-relative commands separately
where necessary (for example `bin/check.js`), and add a regression assertion
that shipped skill instructions/templates contain no live `addons/` paths.
Do not restore an `addons/` compatibility tree, which would violate the agreed
single-tree/package gates.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
