---
idea: skills-cli-install-path
review-round: 01
agent: codex-1
date: 2026-07-29
---

## Summary

Reviewed `parley-deck-skill` branch `readme-skill-catalogue` at `f8e3a1c`. The layout move,
native destination shapes, npm package, tests, universal installation, add-on discovery, and
portable assets work. I block because D-1 removes an explicitly protected install channel
without consensus or user approval, and the stated technical obstacle has a concrete
compatibility-preserving solution.

## Refutation attempts

### G7 — universal installer

- I copied `f8e3a1c` to a scratch source and attempted the exact
  `HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy` command
  without `--full-depth`. The npm wrapper could not resolve `registry.npmjs.org`
  (`ENOTFOUND`) in this sandbox.
- I then ran the locally cached `skills` 1.5.20 CLI—the same version named in the design
  evidence—with the identical `add <path> --agent claude-code --yes --copy` arguments.
  It printed `Found 5 skills` and `Installed 5 skills`.
- Five installed `SKILL.md` files were present. The installed core contained `SKILL.md`,
  `agents/`, and `references/`; assertions confirmed that `bin/`, `lib/`, and
  `package.json` were absent. I therefore reproduce the substantive G7 behavior, while the
  network-backed `@latest` resolution itself remains unavailable in this session.

### Native installer, status, doctor, and destination shapes

- In a clean scratch HOME, the first
  `node bin/parley-deck-skill.js install --target all --force` installed the four runtimes
  detected at command start: Codex, Claude, Antigravity, and Hermes. Creating
  Antigravity's `~/.gemini` tree then made the following `status --target all` detect Gemini
  as missing, and `doctor --target all` exited 1.
- The same one-pass failure occurs at parent commit `94a4889`; it is a pre-existing
  detection-order defect, not a regression from this layout move. A second install detects
  and installs Gemini, after which status and doctor report all 25 core/add-on units valid.
- `validateInstalledPayload` is unchanged. I compared relative path sets from installations
  of `94a4889` and `f8e3a1c`; Codex, Antigravity, and Gemini were identical. Actual core
  destination listings at `f8e3a1c` were:
  - Codex: marker, `SKILL.md`, `agents/`, `references/`, both manifests, README, and license.
  - Antigravity: the Codex shape plus the required fabricated `skills/SKILL.md`.
  - Gemini: the Codex shape, with root `SKILL.md` matching
    `gemini-extension.json.contextFileName`.
- `agy plugin validate` against the scratch Antigravity destination passed with
  `skills: 1 processed` and `agents: 2 processed`.

### Tests, package, portable assets, and deferred gates

- `npm test`: 247 pass, 0 fail.
- `npm pack --dry-run --json`: 153 files, including 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` entries; zero `addons/` entries; no root `SKILL.md`.
- The current-host portable binary built, installed from its embedded assets, and passed
  `doctor`; this additionally exercises `package.json` `pkg.assets`.
- G3 is partially testable on this Mac. With cached `pkg` runtimes, both Windows x64 and
  ARM64 PE binaries built successfully. Installation/execution remains NOT TESTED because
  neither Windows nor Wine is available.
- G8 is testable and passes. I installed all five from a scratch `file://` Git remote,
  committed a new core reference file to that remote, and ran `skills update --yes`.
  The CLI reported all five updated and the new file appeared in the installed core.
- G9 remains NOT TESTED. Homebrew is installed, but this repository contains no candidate
  formula and the external formula points at a published release, so `brew upgrade` or
  `brew test` would test old published material rather than this branch.
- G10 remains NOT TESTED because `winget` is unavailable on macOS.
- G11 remains NOT TESTED because Gemini CLI is not installed. It must not be converted from
  an untested gate into an unsupported capability without the ratified decision being
  changed.

## Findings

### [MAJOR] D-1 removes a protected Gemini channel even though the two manifest consumers can be reconciled

The original constraint says existing legacy Gemini installation must keep working, A4 adds
an actual G11 install gate, and FINAL.md requires every existing channel to be run. D-1
instead removes the README command and declares the channel unsupported. Honest disclosure
of the removal at `README.md:209-211` is preferable to concealing it, but it does not
authorize divergence from the accepted design.

The claim that one `contextFileName` cannot serve both layouts is too narrow. Gemini's current
extension manager joins `contextFileName` to the extension root and accepts nested paths that
remain inside it; it also loads agent skills from the extension's `skills/` directory
([official source](https://github.com/google-gemini/gemini-cli/blob/main/packages/cli/src/config/extension-manager.ts)).
Set the repository manifest to `skills/parley-deck/SKILL.md`, then have the native installer
rewrite the staged Gemini manifest to `SKILL.md` for its unchanged flat destination. This
keeps a single canonical skill tree, preserves every native destination path, and serves both
consumers. Add focused tests for both staged manifest values and run G4 plus an actual G11.
Alternatively, obtain an explicit user-approved capability removal and amend the ratified
decision before merge.

### [MINOR] G1 is recorded as an unconditional pass but is environment-dependent

IMPLEMENTATION.md says the only non-valid status line was the missing project protocol. In a
clean HOME, the exact requested install → status → doctor sequence produces five missing
Gemini units and a non-zero doctor result because the Antigravity install itself creates the
evidence that detects Gemini. The parent commit proves this is not a layout regression, but
the claimed evidence is not reproducible as written. Record the preconditions, or make the
gate deterministic with `--include-undetected`; track the Antigravity-to-Gemini detection
ordering as a separate pre-existing defect.

### [MINOR] The README copied into native installs contains paths that do not exist there

`README.md:154-155` now tells readers to attach
`skills/parley-deck/{SKILL.md,references/COOPERATION.md}`. Those paths are correct in a
repository checkout, but the README is also copied into every native core destination, where
the unchanged shape contains root `SKILL.md` and `references/COOPERATION.md`; the documented
nested paths are absent. Distinguish checkout paths from installed-destination paths, or
stage a destination-specific README.

### [MINOR] “Whichever coding agents you have” exceeds the measured and supported scope

The panel satisfies F1, F4, and F5, avoids an exact universal-agent count per F3, and
`skills list` did show all five after G7. However, `README.md:109-110` says the command
installs into “whichever coding agents you have.” The implementation measured Claude Code,
not every possible agent, and the upstream CLI supports a finite named set. Replace this
with “supported coding agents it detects” or similarly bounded wording; that preserves the
“most agents” message without making an unmeasured universal claim.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
