---
idea: skills-cli-install-path
review-round: 05
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at
`4f7fd32`. The current files, native destination shapes, Gemini manifest
reconciliation, package contents, shipped commands, and the cycle-3 through
cycle-5 negative tests were checked by running them.

Cycle 5 fixes the concrete backtick-fence miss from round 04: broken commands in
new triple-backtick fenced and inline forms both turn the suite red, and narrowing
the extractor to inline spans makes its fixture test fail. One blocking gap
remains. The extractor still ignores the other standard Markdown fence marker,
`~~~`; a newly published broken command in a tilde-fenced block leaves all 253
tests green. The guard therefore remains gameable and still does not enforce its
claim to run every published `node --test` command.

## Prior finding dispositions

### FIXED — round 01 [MAJOR] D-1 removed a reconcilable Gemini channel

D-1 remains withdrawn. Both consumers resolve:

- Repository `gemini-extension.json` uses
  `contextFileName: "skills/parley-deck/SKILL.md"`, and that file exists in the
  repository.
- A native Gemini install stages `contextFileName: "SKILL.md"`, and that file
  exists in the destination.
- Removing `contextFileName` from the two parsed manifests leaves equal objects;
  the staged rewrite changes no other field.

The two regression tests genuinely fail independently. Changing the repository
value to `SKILL.md` makes the repository-manifest test fail while the staged test
passes. Changing the rewrite to `BROKEN.md` makes the staged test fail while the
repository test passes. Reconciliation is the right design; G11 remains an
honestly disclosed end-to-end unknown because `gemini` is absent on this host.

### FIXED — round 01 [MINOR] G1 was recorded as an unconditional pass

The fix-up record retains the clean-HOME precondition and ordering disclosure. In
a new isolated HOME, the exact first sequence produced:

- `install --target all --force`: exit 0;
- `status --target all`: exit 0, with five Gemini units missing;
- `doctor --target all`: exit 1.

The Antigravity install creates the `.gemini` evidence that causes the later
commands to detect legacy Gemini. A second identical install converges:
`status` reports all 25 detected core/add-on units valid and `doctor` exits 0.
This is the previously demonstrated pre-move ordering defect, not a regression
from the layout change.

### FIXED — round 01 [MINOR] The installed README named checkout-only paths

The README distinguishes repository paths under `skills/parley-deck/` from the
flat installed paths `SKILL.md` and `references/COOPERATION.md`. Both stated
shapes exist.

### FIXED — round 01 [MINOR] Universal-agent wording exceeded measured scope

The panel attributes the supported-agent surface to the upstream installer,
states no agent count, and no longer says “whichever coding agents you have.”
The cached upstream CLI contains 75 named agent configurations versus 14 native
targets in this installer, so “a longer list” is measured. The five-skill source
claim is backed by G7.

### FIXED — round 02 [MAJOR] Shipped instructions referenced the removed `addons/` tree

`grep -rn "addons/" skills/` returns no matches. There is no `addons/`
directory, compatibility tree, or symlink, and the package contains no
`addons/` entry.

The regression guard is genuine. Replacing one
`skills/parley-tracker` path with `addons/parley-tracker` in a scratch copy makes
`npm test` fail with 252 passes and one failure, explicitly naming
`skills/parley-tracker/SKILL.md:79`.

### FIXED — round 03 [MAJOR] The tracker exemplar's test command failed with zero passes

Both exemplar occurrences publish
`node --test "skills/parley-tracker/bin/*.test.js"`. It runs 35 tests, all
passing. Reverting the exemplar to `node --test skills/parley-tracker/bin` makes
the suite fail with 252 passes and one failure; the child reports
`MODULE_NOT_FOUND`, zero passes, and one failure.

The other published command,
`node --test "skills/parley-design-check/test/*.test.js"`, runs 159 tests, all
passing. The design-check help command exits 0, and strict validation of all
three tracker templates passes.

### PARTIALLY FIXED — round 04 [MAJOR] The guard ignored fenced `node --test` commands

The exact original miss is fixed:

- A newly added broken triple-backtick fenced command makes `npm test` fail
  252/1.
- A newly added broken inline command makes `npm test` fail 252/1.
- Returning before fenced-block extraction makes the focused extractor fixture
  fail 0/1 and show both missing fenced fixture commands.
- The unmodified suite finds both current unique published commands and passes.

It is only partially fixed because the extractor recognizes backtick fences but
not tilde fences. The remaining failure is detailed below.

## Refutation attempts and gate results

### G7 — install without `--full-depth`

I copied `4f7fd32` to scratch and ran the requested
`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy`
spelling with no `--full-depth`. The wrapper failed before the CLI ran because
`registry.npmjs.org` could not resolve (`ENOTFOUND`), so the network-backed
`@latest` resolution is not reproducible in this environment.

Using the locally cached `skills` 1.5.20 CLI with identical `add` arguments and
still no `--full-depth` produced `Found 5 skills` and `Installed 5 skills`.
Exactly five installed `SKILL.md` roots exist. The core contains exactly
`SKILL.md`, `agents/`, and `references/`; `bin/`, `lib/`, and `package.json` are
absent.

### Native destinations, validation, and Antigravity

`validateInstalledPayload` remains destination-relative. After the converged
native install, the actual core destination path sets are:

- Codex: 13 entries.
- Antigravity: 15 entries, including fabricated `skills/SKILL.md`.
- Gemini: 13 entries, with staged `contextFileName: "SKILL.md"`.

Each set is identical to the corresponding install from pre-move parent
`94a4889`. `agy plugin validate` against the installed Antigravity destination
exits 0 with one skill and two agents processed. No per-target destination shape
changed.

### Tests, package, and README

- `npm test`: 253 passed, 0 failed.
- `npm pack --dry-run --json`, using a writable scratch npm cache: 153 files,
  145 under `skills/`, exactly five `skills/<name>/SKILL.md` roots, no
  `addons/`, and no root `SKILL.md`.
- The README panel is first under Install, labels neither path recommended,
  states no numeric universal-agent count, includes `skills list`, and keeps the
  native installer plus `doctor` in the same screenful. F4 is not triggered
  because the five-skill move passes. I found no unsupported panel claim under
  A4/F1/F3/F4/F5.

The implementation's package, destination, manifest, stale-path, exemplar, and
current-command pass claims reproduce. The network-backed `npx ...@latest`
wrapper remains blocked by DNS, and the initial G1 table is only true under the
later recorded precondition. The cycle-5 claim that every fenced form is covered
does not reproduce.

### Previously NOT TESTED gates

- **G3 Windows:** PARTIALLY TESTED. In a scratch copy, using the two cached base
  runtimes, `scripts/build-portable.js windows` builds both executables. `file`
  identifies PE32+ x86-64 and AArch64 binaries. Windows execution/install
  remains NOT TESTED because this host has neither Windows nor Wine.
- **G8 `skills update`: PASS.** I installed all five from an isolated
  `file://` Git remote, advanced that remote with a new core reference file, and
  ran `skills update --yes`. It reported all five updated, and the new file
  appeared in the installed core.
- **G9 Homebrew:** NOT TESTED against `4f7fd32`. Homebrew is present, but the
  tap formula downloads tag `v1.5.0` at `978fce0`, not this review head. A
  `brew test` or `brew upgrade` would test published material rather than this
  branch.
- **G10 WinGet:** NOT TESTED. `winget` and a Windows execution environment are
  unavailable, and this repository contains only `packaging/winget/README.md`.
- **G11 legacy Gemini URL install:** NOT TESTED. `gemini` is absent. The two
  manifest resolutions and their independent negative tests pass, but they do
  not substitute for the end-to-end CLI gate.

## Findings

### [MAJOR] The cycle-5 extractor still ignores tilde-fenced published commands

`test/design-addons.test.js:230` recognizes only triple-backtick fenced blocks.
Tilde fences are also fenced Markdown blocks. In a scratch copy, I inserted this
after the valid `parley-worktrees` frontmatter and heading:

```markdown
~~~bash
node --test definitely-missing-tilde.test.js
~~~
```

`npm test` still exited 0 with 253 passes and no mention of the missing test
file. The equivalent new command inside a triple-backtick fence makes the suite
fail 252/1. Thus the guard remains gameable by a normal published fenced form,
and both its test name (“every `node --test` command”) and the cycle-5
implementation claim are broader than what it enforces.

Extend fenced extraction to matching backtick **and tilde** fence markers, and
add valid and broken tilde-fenced commands to the extractor fixture. Prove the
new broken tilde-fenced form turns the full suite red.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
