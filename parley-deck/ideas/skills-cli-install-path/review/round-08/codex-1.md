---
idea: skills-cli-install-path
review-round: 08
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at
`895d0f4df2e61035d5bbaacdd759d5d1ad069110`. The implementation repository was clean before
and after review; every mutation below was made in an isolated scratch copy.

All ten findings I filed in rounds 01–07 are fixed in their concrete behavior. Cycle 8 catches
the tab-separated broken command and correctly executes a legitimate two-target command.
However, the same guard still changes command text before execution and implements only a
partial shell lexer. I defeated it in both requested directions: one broken fenced command
manually runs zero tests but is silently repaired by the guard, while one legitimate
single-quoted Bash command manually runs 35 tests but is made to fail by the guard. I therefore
block on one new MAJOR.

## Prior finding dispositions

### FIXED — round 01 [MAJOR] D-1 removed a reconcilable Gemini channel

The repository manifest has
`contextFileName: "skills/parley-deck/SKILL.md"` and that file exists. A native Gemini install
stages `contextFileName: "SKILL.md"` and that file exists in its flat destination. Removing
`contextFileName` leaves the two parsed manifests deeply equal, so the rewrite changes no
other current field.

Both regression tests genuinely fail. Changing the repository value to `SKILL.md` made the
focused repository test exit 1 with the wrong actual value. Changing the staged rewrite to
`BROKEN.md` independently made the staged-install test exit 1 and name that value. The README
restores `gemini extensions install <repo-url>` and says that the Gemini CLI has not been run
end to end. Reconciliation remains the right design; G11 remains honestly NOT TESTED because
`gemini` is absent.

### FIXED — round 01 [MINOR] G1 was recorded as an unconditional pass

In a fresh scratch `HOME`, the first exact sequence produced:

- `install --target all --force`: exit 0;
- `status --target all`: exit 0, with all 20 initially detected units valid and five newly
  detected Gemini units missing;
- `doctor --target all`: exit 1.

The identical one-pass sequence at pre-move parent `94a4889` also leaves five Gemini units
missing and doctor at exit 1. A second exact install at the current head converges: status
reports all 25 detected units valid and doctor exits 0. This is the recorded, pre-existing
Antigravity-to-Gemini detection ordering, not a layout regression.

### FIXED — round 01 [MINOR] The installed README named checkout-only paths

The README distinguishes checkout paths under `skills/parley-deck/` from installed paths
`SKILL.md` and `references/COOPERATION.md`. Both checkout paths and both installed paths
exist in the places assigned to them.

### FIXED — round 01 [MINOR] Universal-agent wording exceeded measured scope

The panel attributes the support surface upstream, states no numeric agent count, and no
longer says “whichever coding agents you have.” As a fresh measurement, the cached `skills`
1.5.20 CLI installed all five scratch skills with `--agent '*'` and printed
`Installing to all 75 agents`; this installer's `paths --target all --include-undetected
--json` returned fourteen targets. The comparative “longer list” wording is therefore
measured.

### FIXED — round 02 [MAJOR] Shipped instructions referenced the removed `addons/` tree

The exact `grep -rn "addons/" skills/` command returned no matches. There is no `addons/`
directory, compatibility copy, or symlink, and the package contains no `addons/` entry.

The regression guard genuinely fails. Reintroducing one stale path made the focused suite
exit 1 and named `skills/parley-tracker/SKILL.md:79`.

### FIXED — round 03 [MAJOR] The tracker exemplar's test command failed with zero passes

Both exemplar occurrences now publish
`node --test "skills/parley-tracker/bin/*.test.js"`, which runs 35 tests with 35 pass and
zero fail. Reverting them in a scratch copy to `node --test skills/parley-tracker/bin` made
the full `npm test` suite exit 1 with `MODULE_NOT_FOUND`, zero child passes, and one child
failure.

### FIXED — round 04 [MAJOR] The published-command guard ignored fenced commands

The real fenced design-check command runs 159 tests with 159 pass; the real inline tracker
command runs 35 tests with 35 pass. Replacing the fenced target with a missing file made the
guard suite exit 1. Replacing the inline target independently also made it exit 1.

### FIXED — round 05 [MAJOR] The guard ignored tilde-fenced commands

Adding a missing target inside a `~~~bash` block made the guard suite exit 1 and name that
target. The current structure-independent scan catches the original tilde-container probe.

### FIXED — round 06 [MAJOR] The guard checked only the first command on each line

I repeated the exact bypass shape: one valid tracker command first and one missing command
second on the same shipped Markdown line. The guard suite now exits 1 on the second command.

### FIXED — round 07 [MAJOR] The guard missed shell whitespace and combined multiple targets

A fenced missing target separated from `--test` by a literal tab makes the guard suite exit
1. Narrowing the extractor back to literal spaces makes its fixture test exit 1 and name the
missed tab-separated command.

The legitimate two-target command
`node --test skills/parley-tracker/bin/claim.test.js
skills/parley-tracker/bin/validate.test.js` runs 35 tests manually and also leaves the guard
suite green. Both cycle-8 fixes therefore reproduce.

## Current gate evidence

### G7 — install with no `--full-depth`

I copied the exact head to scratch and ran the requested command with no `--full-depth`:

`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy`

The wrapper failed before the CLI ran because `registry.npmjs.org` could not resolve
(`ENOTFOUND`). I therefore do not claim the network-backed `@latest` resolution passed in
this environment.

Running the locally cached `skills` 1.5.20 executable with the same add arguments, the same
scratch `HOME`, and still no `--full-depth` printed `Found 5 skills` and
`Installed 5 skills`. Exactly five installed `SKILL.md` roots exist. The installed core's
top level is exactly `SKILL.md`, `agents/`, and `references/`; `bin/`, `lib/`, and
`package.json` are absent. The layout behavior required by G7 passes.

### Native destinations, validation, status, and doctor

The extracted `validateInstalledPayload` function has the same SHA-256 at parent `94a4889`
and at this head. I installed both commits and diffed actual relative path sets:

- Codex: 13 entries at each commit, empty diff;
- Antigravity: 15 entries at each commit, empty diff;
- Gemini: 13 entries at each commit, empty diff.

The current Codex destination has the flat core payload, both manifests, README, license,
and marker. Antigravity has that shape plus fabricated `skills/SKILL.md`. Gemini has the
flat shape and its staged `contextFileName: "SKILL.md"`. No per-target destination shape
changed. `agy plugin validate` against the actual scratch Antigravity destination exits 0
with one skill and two agents processed.

### Tests, package, commands, and channel gates

- `npm test`: 253 pass, 0 fail.
- `npm pack --dry-run --json`: 153 entries; 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` roots; zero `addons/` entries; no root `SKILL.md`.
- Every currently published test command runs a non-zero suite: the three Markdown
  occurrences reduce to the 159-test design-check command and the 35-test tracker command.
- **G3 Windows:** PARTIALLY TESTED. A scratch build produced x64 and ARM64 executables;
  `file` identifies them as PE32+ x86-64 and AArch64. Execution/install remains NOT TESTED
  because this host has neither Windows nor Wine.
- **G8 `skills update`: PASS.** I installed all five from a scratch `file://` Git remote,
  advanced that remote with a new core reference file, and ran `skills update --yes`. It
  reported all five updated, and the new file appeared in the installed core.
- **G9 Homebrew:** NOT TESTED against this head. Homebrew is installed, but the tap formula
  downloads the fixed published `v1.5.0` tarball; `brew test` or `brew upgrade` would test
  that release, not commit `895d0f4`.
- **G10 WinGet:** NOT TESTED. `winget` and a Windows host are absent, and this repository
  contains packaging instructions but no candidate manifest to validate.
- **G11 Gemini URL install:** NOT TESTED. The Gemini executable is absent. Both static
  resolutions and their independent negative tests pass, but they do not replace the CLI
  gate.

The README panel satisfies A4/F1/F3/F4/F5: neither route is labelled recommended; no numeric
agent count is claimed; the “longer list” comparison is measured as described above;
`skills list` is present; and the native installer plus doctor remain in the same screenful.
The all-five statement is backed by G7. F4 is not triggered because the five-skill move
passes. I found no unmeasured panel claim.

The passing package, layout, destination, manifest, path-guard, real-command, G8, and
cycle-8 probe claims in `IMPLEMENTATION.md` reproduce. The exact network-backed
`npx ...@latest` wrapper is the only claimed pass I cannot rerun, due to DNS rather than an
observed implementation failure. The initial unconditional G1 table is superseded by the
explicit correction in fix-up cycle 1.

## Findings

### [MAJOR] The cycle-8 guard still silently repairs broken commands and breaks legitimate Bash quoting

`test/design-addons.test.js:233-237` captures a command and then removes trailing
punctuation and normalizes its whitespace. Lines 315-320 implement an argv lexer that
understands only double quotes. These transformations mean the guard still does not execute
the command the shipped Markdown actually publishes.

I measured both directions in separate scratch copies:

1. I added this command inside a fenced Bash block:

   `node --test "skills/parley-tracker/bin/*.test.js".`

   Typed exactly, the shell concatenates the final period to the target. Node exits 0 after
   running **zero tests**. The guard strips the period as presumed prose punctuation, turns
   the broken target back into the valid existing glob, and the full suite remains
   **253 pass / 0 fail**. This is the exact false-green class the guard claims to prevent.

2. I added this legitimate fenced Bash command:

   `node --test 'skills/parley-tracker/bin/*.test.js'`

   Typed exactly, it runs **35 tests, 35 pass, 0 fail**. The guard's lexer leaves the single
   quotes in the argv value, Node runs zero matching tests, and `npm test` exits 1 with the
   assertion that the published command ran zero tests. The guard manufactures a failure.

This is not another Markdown-container or whitespace case. It is the unresolved mismatch
between arbitrary shell syntax and a handwritten normalizer/lexer. Fix it without another
case-by-case parser extension: capture command text without silently rewriting it, and use
a fail-closed exact-command-to-argv registry (or an actual, explicitly scoped shell) so an
unrecognized form fails as unsupported instead of being guessed. Add both probes above as
regressions: the period-suffixed command must turn the suite red, while the single-quoted
command must execute 35 tests and stay green.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
