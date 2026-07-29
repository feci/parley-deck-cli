---
idea: skills-cli-install-path
review-round: 07
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at
`c64263649f443ebeef973a33c9fd527d8f6c5478`. The repository was clean and all mutations were
made in isolated scratch clones.

All nine findings I filed in rounds 01–06 are fixed in their concrete behavior. The exact
round-06 two-commands-on-one-line bypass now turns the suite red, as do the earlier stale
path, exemplar, inline, backtick-fenced, tilde-fenced, indented, and four-backtick probes.
The install layout, Gemini reconciliation, native destination shapes, package contents, G7,
and G8 also reproduce.

I found one new MAJOR in the cycle-7 guard. Its global match fixes the positional bypass but
still does not recognize ordinary shell whitespace, and its runner treats the entire command
tail as one path rather than as command arguments. A broken command separated with a tab is
invisible and leaves all 253 tests green; conversely, a valid two-file `node --test` command
runs 35 tests when typed normally but turns the suite red because the guard combines both
paths into one nonexistent path. The claim that the guard covers every published
`node --test` command therefore still does not reproduce.

## Prior finding dispositions

### FIXED — round 01 [MAJOR] D-1 removed a reconcilable Gemini channel

The repository manifest uses
`contextFileName: "skills/parley-deck/SKILL.md"` and that path exists. A native Gemini
install stages `contextFileName: "SKILL.md"` and that path exists in its destination.
Removing `contextFileName` leaves the two parsed manifests deeply equal, so the rewrite
changes no other field.

Both tests genuinely fail independently. Changing the repository value to `SKILL.md` made
the focused repository test fail 0/1 with the expected and actual values. Changing the
staged rewrite to `BROKEN.md` made the focused staged-install test fail 0/1. Reconciliation
was the right design, and the README restores the URL channel while honestly disclosing that
the Gemini CLI has not been run end to end.

### FIXED — round 01 [MINOR] G1 was recorded as an unconditional pass

In a fresh scratch `HOME`, the exact first sequence produced:

- `install --target all --force`: exit 0;
- `status --target all`: exit 0, with five Gemini units missing;
- `doctor --target all`: exit 1.

The same one-pass sequence at pre-move parent `94a4889` also leaves five Gemini units missing
and doctor at exit 1. A second exact install detects and installs Gemini; status then reports
all 25 detected core/add-on units valid and doctor exits 0. The fix-up record correctly
qualifies this pre-existing Antigravity-to-Gemini detection ordering.

### FIXED — round 01 [MINOR] The installed README named checkout-only paths

The README distinguishes repository paths under `skills/parley-deck/` from installed paths
`SKILL.md` and `references/COOPERATION.md`. Both stated shapes exist.

### FIXED — round 01 [MINOR] Universal-agent wording exceeded measured scope

The panel attributes support to the upstream installer, states no numeric agent count, and no
longer says “whichever coding agents you have.” The cached `skills` 1.5.20 CLI listed 75
valid agents versus this package's fourteen named native runtimes, so “a longer list” is
measured.

### FIXED — round 02 [MAJOR] Shipped instructions referenced the removed `addons/` tree

The exact `grep -rn "addons/" skills/` command returned no output. There is no `addons/`
directory, compatibility tree, or symlink, and the package contains no `addons/` entry.

The regression guard genuinely fails. Replacing one current path with
`addons/parley-tracker` made the full suite fail 252/1 and named
`skills/parley-tracker/SKILL.md:79`.

### FIXED — round 03 [MAJOR] The tracker exemplar's test command failed with zero passes

Both exemplar occurrences now publish
`node --test "skills/parley-tracker/bin/*.test.js"`, which runs 35 tests with 35 pass and
zero fail. Reverting them to `node --test skills/parley-tracker/bin` made the full suite fail
252/1 with `MODULE_NOT_FOUND`, zero child passes, and one child failure.

### FIXED — round 04 [MAJOR] The guard ignored backtick-fenced commands

Breaking the existing fenced design-check command made the full suite fail 252/1. Breaking
the existing inline tracker command independently also made it fail 252/1. Narrowing the
extractor back to inline spans made its fixture fail and reduced published-command discovery
to one.

### FIXED — round 05 [MAJOR] The guard ignored tilde-fenced commands

Adding a broken command in a `~~~` block made the full suite fail 252/1 and named the missing
target. Broken commands in an indented block and a four-backtick fence independently produced
the same result.

### FIXED — round 06 [MAJOR] The guard checked only the first command on each line

I repeated the exact round-06 probe: one valid tracker command first and one broken command
second on the same shipped Markdown line. The full suite now fails 252/1 on the second
command. The extended fixture also asserts that both same-line commands are extracted.

## Refutation attempts and gate results

### G7 — real install without `--full-depth`

I copied the exact review commit and ran:

`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy`

with no `--full-depth`. The npm wrapper failed before the CLI ran because
`registry.npmjs.org` could not resolve (`ENOTFOUND`), so I do not claim the network-backed
`@latest` resolution passed in this session.

Using the locally cached `skills` 1.5.20 executable with the same add arguments and still no
`--full-depth` produced `Found 5 skills` and `Installed 5 skills`. Exactly five installed
`SKILL.md` roots were present. The installed core contained only `SKILL.md`, `agents/`, and
`references/` at its top level; `bin/`, `lib/`, and `package.json` were absent.

### Native installer, status, doctor, and destination shapes

The first exact install/status/doctor sequence reproduces the known ordering described above;
the second install converges with all 25 detected units valid. `agy plugin validate` against
the installed Antigravity destination exits 0 with one skill and two agents processed.

The extracted `validateInstalledPayload` function has the same SHA-256 at `94a4889` and
`c642636`. I installed both commits into separate scratch homes and diffed relative path
sets; Codex, Antigravity, and Gemini each had an empty diff. Current core destinations are:

- Codex: marker, `LICENSE`, `README.md`, root `SKILL.md`, `agents/`, both manifests, and
  `references/`.
- Antigravity: the Codex shape plus fabricated `skills/SKILL.md`.
- Gemini: the Codex shape, with staged `contextFileName: "SKILL.md"`.

No per-target installed destination shape changed.

### Tests, package, shipped commands, and regression guards

- `npm test`: 253 pass, 0 fail.
- `npm pack --dry-run --json`: 153 files; 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` roots; no `addons/`; no root `SKILL.md`.
- `node --test "skills/parley-design-check/test/*.test.js"`: 159 pass, 0 fail.
- `node --test "skills/parley-tracker/bin/*.test.js"`: 35 pass, 0 fail.
- The design-check help command, tracker strict-directory validation over all three
  exemplars, tracker single-file validation, and tracker claim command all resolved and
  exited 0.
- The two Gemini manifest tests fail independently when their respective values are broken.
- The stale-path guard fails with file and line.
- Reverted exemplar, broken fenced, broken inline, broken tilde, broken indented, broken
  four-backtick, and the round-06 same-line probes all turn the suite red.

The cycle-7 guard can nevertheless still be defeated or made to misread valid syntax; see
the finding below.

### Previously NOT TESTED gates

- **G3 Windows:** PARTIALLY TESTED. With cached x64 and ARM64 Node runtimes,
  `scripts/build-portable.js windows` built both v1.5.0 executables. `file` identifies them
  as PE32+ x86-64 and AArch64. Windows execution/install remains NOT TESTED because this host
  has neither Windows nor Wine.
- **G8 `skills update`: PASS.** I installed all five from an isolated `file://` Git remote,
  advanced that remote with a new core reference file, and ran `skills update --yes`. It
  reported all five updated, and the new reference appeared in the installed core.
- **G9 Homebrew:** NOT TESTED against `c642636`. Homebrew is installed, but this repository
  has no candidate formula. The tapped formula downloads published tag `v1.5.0`, not this
  review commit, so `brew test` or `brew upgrade` would test different source.
- **G10 WinGet:** NOT TESTED. `winget` and a Windows host are unavailable; this repository
  contains only `packaging/winget/README.md`.
- **G11 legacy Gemini URL install:** NOT TESTED. The `gemini` executable is absent. Both
  static manifest resolutions and both independent negative tests pass, but they do not
  substitute for the end-to-end CLI gate.

### README panel and implementation claims

The panel satisfies A4/F1/F3/F4/F5: neither path is labelled recommended; no numeric
universal-agent count is claimed; the comparative “longer list” is measured; `skills list`
is present; and the native installer plus doctor remain in the same screenful. The all-five
claim is backed by G7. F4 is not triggered because the five-skill move passes. The Gemini
paragraph is honest about its remaining end-to-end unknown.

The package, destination, manifest, stale-path, current-command, and specific negative-probe
claims in `IMPLEMENTATION.md` reproduce. The network-backed `npx ...@latest` wrapper remains
blocked by DNS, and the initial G1 table is only true under the later recorded precondition.
The cycle-7 claim that the global match covers every published command does not reproduce.

## Findings

### [MAJOR] The cycle-7 guard still misses and mis-executes valid shell command syntax

`test/design-addons.test.js:233` matches only the literal ASCII sequence
`node --test `, while lines 305 and 313 treat everything after it as one test target and pass
that combined string as a single argument. The global match fixes “second occurrence on the
same line,” but it does not model the command syntax the test claims to execute.

I measured both failure modes in separate scratch copies:

1. I published a fenced command with a literal tab between `--test` and
   `definitely-missing-tab-command.test.js`. A shell treats the tab as ordinary argument
   whitespace; typing the command exits non-zero and names the missing file. The guard did
   not discover it, and the full suite remained **253 pass / 0 fail**.
2. I published
   `node --test skills/parley-tracker/bin/claim.test.js skills/parley-tracker/bin/validate.test.js`.
   Typed normally, this legitimate two-target command runs **35 tests, 35 pass, 0 fail**.
   The guard instead invoked Node with one combined target containing both paths and made the
   full suite fail **252/1** with “Could not find
   `claim.test.js validate.test.js`.”

This is not another Markdown-container variant. It is a mismatch between the claimed shell
command surface and a handwritten extractor/runner that supports only one exact whitespace
form and one target argument. Either define and enforce an explicitly narrow published
command grammar and scope the test/implementation claims to it, or recognize shell
whitespace and execute the captured command with its real argument boundaries. Add negative
coverage for a tab-separated broken target and positive coverage for multiple targets (plus
options or trailing comments if those are intended to be supported).

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
