---
idea: skills-cli-install-path
review-round: 06
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at
`46b5730eba6ac93b091a2a58ffce7f44823d3d70`. The implementation repository remained clean;
all mutations were made in isolated scratch clones.

All eight findings I filed in rounds 01–05 are fixed in their concrete behavior. The two
current published test commands run 159 and 35 passing tests, respectively; the manifest
reconciliation, stale-path guard, package shape, native destination shapes, G7 layout
behavior, and G8 update behavior all reproduce.

Cycle 6 fixes the prior tilde-fence miss, but its claim to close the published-command class
does not reproduce. The extractor checks only the first `node --test` occurrence on each
line. A second broken command on the same shipped Markdown line is invisible and leaves all
253 tests green. That remaining guard bypass is a MAJOR finding.

## Prior finding dispositions

### FIXED — round 01 [MAJOR] D-1 removed a reconcilable Gemini channel

The repository manifest uses `contextFileName: "skills/parley-deck/SKILL.md"` and that path
exists. A native Gemini install stages `contextFileName: "SKILL.md"` and that path exists in
the destination. All other manifest fields compare equal.

Both tests can fail independently. Changing the repository value to `SKILL.md` made the
repository-manifest test fail 0/1. Changing the staged rewrite to `BROKEN.md` made the
staged-install test fail 0/1. The README restores the URL channel and honestly says the
Gemini CLI has not been run end to end.

### FIXED — round 01 [MINOR] G1 was recorded as an unconditional pass

In a fresh scratch `HOME`, the first exact `install --target all --force` exited 0; the
following status reported five missing Gemini units and doctor exited 1. A second identical
install converged, after which all 25 detected core/add-on units were valid and doctor
exited 0. The implementation record retains this precondition and does not call the
first-pass ordering behavior a layout regression.

### FIXED — round 01 [MINOR] The installed README named checkout-only paths

The README now distinguishes checkout paths under `skills/parley-deck/` from the flat
installed paths `SKILL.md` and `references/COOPERATION.md`. Both stated shapes exist.

### FIXED — round 01 [MINOR] Universal-agent wording exceeded measured scope

The panel attributes support to the upstream installer, gives no agent count, and no longer
says “whichever coding agents you have.” The cached `skills` 1.5.20 bundle contains 75 named
agent configurations versus this installer's 14 named targets, so “a longer list” is
measured.

### FIXED — round 02 [MAJOR] Shipped instructions referenced the removed `addons/` tree

The exact `grep -rn "addons/" skills/` command returned no output. There is no `addons/`
tree, compatibility copy, or symlink, and the package contains no `addons/` entry.

The regression guard genuinely fails: reintroducing one stale path made the focused test
name `skills/parley-tracker/SKILL.md:79`, and the full suite went to 252 pass / 1 fail.

### FIXED — round 03 [MAJOR] The tracker exemplar's test command failed with zero passes

The published replacement, `node --test "skills/parley-tracker/bin/*.test.js"`, runs 35
tests with 35 pass and 0 fail. Reverting it in a scratch copy to
`node --test skills/parley-tracker/bin` made the full suite fail 252/1; the child command
reported `MODULE_NOT_FOUND`, zero passes, and one failure.

### FIXED — round 04 [MAJOR] The guard ignored backtick-fenced commands

Breaking the current fenced design-check command made the full suite fail 252/1. Breaking
the current inline tracker command independently produced the same result. Narrowing the
extractor fixture back to inline-backtick lines made its focused test fail and name the
missed fenced command.

### FIXED — round 05 [MAJOR] The guard ignored tilde-fenced commands

Adding a broken `node --test` command inside a `~~~bash` block made the full suite fail
252/1 and name the nonexistent target. The original tilde-container miss is fixed; the new
finding below is about multiple command occurrences on one line, not another container.

## Refutation attempts and gate results

### G7 — install without `--full-depth`

I copied the exact review commit and ran:

`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy`

with no `--full-depth`. The wrapper could not resolve `registry.npmjs.org` (`ENOTFOUND`), so
I do not claim the network-backed `@latest` resolution passed. Running the locally cached
`skills` 1.5.20 CLI with the same add arguments produced `Found 5 skills` and
`Installed 5 skills`. Exactly five installed `SKILL.md` roots exist. The installed core
contains `SKILL.md`, `agents/`, and `references/`; `bin/`, `lib/`, and `package.json` are
absent.

### Native installer, destination shapes, status, and doctor

The exact first install/status/doctor sequence reproduced the known clean-`HOME`
Antigravity-to-Gemini ordering described above. The second install converged with all 25
detected units valid and doctor exit 0. `agy plugin validate` against the installed
Antigravity destination passed with one skill and two agents processed.

`validateInstalledPayload` has the same SHA-256 at pre-move parent `94a4889` and at this
head. I installed both commits and compared relative path sets for the actual core
destinations; Codex, Antigravity, and Gemini each had diff exit 0. Current Codex has the
flat core payload, Antigravity adds fabricated `skills/SKILL.md`, and Gemini retains the
flat core payload with its staged manifest rewrite. No per-target destination shape changed.

### Tests, package, shipped commands, and guards

- `npm test`: 253 pass, 0 fail.
- `npm pack --dry-run --json`: 153 files; 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` roots; no `addons/`; no root `SKILL.md`.
- `node --test "skills/parley-design-check/test/*.test.js"`: 159 pass, 0 fail.
- `node --test "skills/parley-tracker/bin/*.test.js"`: 35 pass, 0 fail.
- The concrete design-check help, tracker strict-directory, tracker single-file, and tracker
  claim commands all resolved and exited 0 in scratch state.
- A single broken inline, backtick-fenced, or tilde-fenced command turns the suite red.
- A second broken inline command on a line whose first command is valid does not; see the
  finding below.

### Previously NOT TESTED gates

- **G3 Windows:** PARTIALLY TESTED. With the cached x64 and ARM64 base runtimes copied to a
  scratch pkg cache, `scripts/build-portable.js windows` built both binaries. `file`
  identifies PE32+ x86-64 and AArch64 executables. Windows execution/install remains
  NOT TESTED because this host has neither Windows nor Wine.
- **G8 `skills update`: PASS.** I installed all five from an isolated `file://` Git remote,
  committed a new core reference file to that remote, and ran `skills update --yes`. It
  reported all five updated, and the new file appeared in the installed core.
- **G9 Homebrew:** NOT TESTED against this head. Homebrew is installed, but this repository
  contains no candidate formula; `brew test` or `brew upgrade` would exercise separately
  published material rather than `46b5730`.
- **G10 WinGet:** NOT TESTED. `winget` and a Windows host are unavailable, and this
  repository contains only `packaging/winget/README.md`.
- **G11 legacy Gemini URL install:** NOT TESTED. The `gemini` executable is absent. Both
  static manifest resolutions and both independent negative tests pass, but they do not
  substitute for the end-to-end CLI gate.

### README and implementation claims

The README panel satisfies A4/F1/F3/F4/F5: neither path is labelled recommended; no numeric
agent count is claimed; `skills list` is present; the native installer and doctor remain in
the same screenful; and the all-five statement is backed by G7. The comparative “longer
list” claim is measured as 75 versus 14. F4 is not triggered because the five-skill move
passes. I found no unsupported panel claim.

The package, manifest, destination, stale-path, real-command, and five specifically measured
container-form claims in `IMPLEMENTATION.md` reproduce. The exact network-backed
`npx ...@latest` wrapper remains blocked by DNS. The cycle-6 assertion that line scanning
“closes the class by construction” does not reproduce.

## Findings

### [MAJOR] The cycle-6 guard checks only the first published command on each line

`test/design-addons.test.js:230-240` calls `line.indexOf("node --test ")` once, slices from
that first occurrence, and stops at its first closing backtick. It never resumes scanning
the remainder of the line.

In a scratch copy I added this single shipped Markdown line:

```markdown
Round-06 probe: first run `node --test "skills/parley-tracker/bin/*.test.js"`; then run `node --test definitely-missing-second-command.test.js`.
```

The first command is valid and runs 35 tests. The second is a real, exact
`node --test ` occurrence with a nonexistent target. `npm test` still exited 0 with
253 pass and 0 fail. By contrast, publishing that broken command as the only occurrence on
its line makes the suite fail. The extractor fixture also contains at most one command per
line, so it cannot expose this bypass.

Scan every occurrence on each line, advancing past the closing delimiter before searching
again. Add an extractor fixture with two inline commands on one line and a full negative
test whose first command is valid and whose second is broken. The guard should fail on the
second command.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
