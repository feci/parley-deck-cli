---
idea: skills-cli-install-path
review-round: 13
agent: codex-1
date: 2026-07-29
---

## Summary

Reviewed clean branch `readme-skill-catalogue` at `82507b5`. The layout move, package
contents, native destination shapes, Gemini reconciliation, README corrections, and the
concrete regression probes from rounds 01–12 reproduce. `npm test` is green at 253/253 on the
unmodified tree.

Cycle 13 fixes both exact round-12 probes, but it does not close the continuation class. A
backslash continuation placed between `node` and `--test` makes a broken published command
invisible to the guard: the copied command exits 1 while the full suite remains green at
253/253. That is one remaining MAJOR finding.

## Prior finding dispositions

| Filed finding | Status | Current verification |
|---|---|---|
| Round 01 MAJOR — D-1 dropped a reconcilable Gemini channel | **FIXED** | Repository `contextFileName` is `skills/parley-deck/SKILL.md` and resolves in the repository. A native Gemini install rewrites it to `SKILL.md`, which resolves in the flat destination; the other three manifest fields remain equal. Mutating the repository value to `SKILL.md` made its focused test fail, and mutating the staged rewrite to `BROKEN.md` independently made the staged-install test fail. |
| Round 01 MINOR — G1 was recorded as an unconditional pass | **FIXED** | The fix-up record now states the precondition and the pre-existing detection-order defect. In a fresh HOME, first `install --target all --force`, status, and doctor returned 0/0/1 because Antigravity made Gemini detectable only after target resolution. A second identical install made all 25 units valid and status/doctor returned 0/0. |
| Round 01 MINOR — installed README documented checkout-only paths | **FIXED** | README now distinguishes repository paths under `skills/parley-deck/` from root `SKILL.md` and `references/` in native destinations. Both actual shapes exist. |
| Round 01 MINOR — “whichever coding agents you have” exceeded the measured scope | **FIXED** | The panel says “the coding agents it supports,” attributes the support surface upstream, and asserts no agent count. |
| Round 02 MAJOR — shipped instructions retained `addons/` paths | **FIXED** | Exact `grep -rn "addons/" skills/` returned no matches; no `addons/` tree or root `SKILL.md` exists. Reintroducing one stale path made the guard fail and named `skills/parley-tracker/SKILL.md:79`. |
| Round 03 MAJOR — tracker exemplar published a zero-pass directory command | **FIXED** | Both published test commands run real suites: design-check 159/159 and tracker 35/35. Reverting the exemplar to `node --test skills/parley-tracker/bin` made the guard fail with zero passes. |
| Round 04 MAJOR — guard ignored fenced commands | **FIXED** | A broken command added to a backtick-fenced block made the guard fail. A separate broken inline command also failed. |
| Round 05 MAJOR — extractor ignored tilde fences and its fixture could not prove coverage | **FIXED** | A broken command added to a tilde-fenced block made the guard fail. Narrowing the extractor to literal-space detection made its fixture fail on `node --test<TAB>...`. |
| Round 06 MAJOR — only the first command on one line was checked | **FIXED** | The executed extractor fixture asserts both commands from its same-line pair, and the current whole-line compound probe is refused. |
| Round 07 MAJOR — tab whitespace was missed and two targets became one argument | **FIXED** | The fixture asserts the tab form. A real two-target command ran 35/35, and command execution now goes through the shell rather than one combined argv value. |
| Round 08 MAJOR — punctuation was silently repaired and single quotes were misread | **FIXED** | The single-quoted glob ran 35/35. The trailing-period command ran zero tests, and publishing that exact form made the guard fail with “ran zero tests.” |
| Round 09 MAJOR — surrounding shell context was discarded | **FIXED** | The current fixture captures and refuses the environment-prefix, `cd ... &&`, and command-substitution forms as whole unsupported units. |
| Round 09 MINOR — README contradicted itself about agent detection | **FIXED** | The native comparison now claims only the measured health checks and metadata sync as additions. |
| Round 09 NIT — claimed-deleted lexer remained | **FIXED** | The handwritten argv lexer and stale normalisation code remain absent. |
| Round 10 MAJOR — fenced substitution was repaired and double-backtick spans were misread | **FIXED** | A fenced substitution and the fake-closing-fence variant both made the guard fail by named refusal. A legitimate double-backtick command stayed green and ran a non-zero suite. |
| Round 11 MAJOR — handwritten fence state mishandled fake closes and blockquotes | **FIXED** | The exact fake-close substitution failed by refusal; a valid command in a blockquoted fence stayed green. |
| Round 12 MAJOR — the span-or-line discriminator dropped executable command text | **PARTIALLY FIXED** | The exact compound span-plus-shell probe and the continuation-after-target probe now turn the guard red, while the blockquote and double-backtick valid forms remain green. The new finding below shows that a continuation before `--test` still bypasses the same physical-line discriminator. |

## Refutation attempts and gate results

### G7 — real universal install, no `--full-depth`

I copied tracked HEAD to a scratch source and ran the exact
`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy` spelling,
without `--full-depth`. The wrapper could not resolve `registry.npmjs.org` (`ENOTFOUND`), so
the network-backed `@latest` resolution is the one passing step I cannot reproduce here.

Using the locally cached `skills` 1.5.20 executable—the version used by the ratified design
evidence—with the same arguments and no `--full-depth` found and installed five skills.
`skills list` listed all five. The installed core contained `SKILL.md`, `agents/`, and
`references/`; `bin/`, `lib/`, and `package.json` were absent.

### Native installer, validators, and destination shapes

- The exact first install/status/doctor sequence has only the documented clean-HOME ordering
  behavior: return codes 0/0/1 with Gemini missing. The second identical install produced
  25 valid core/add-on units and return codes 0/0/0. This is the previously recorded
  detection-order issue, not a layout regression.
- `validateInstalledPayload` is textually identical to parent `94a4889`.
- Relative path-set diffs between parent `94a4889` and current HEAD were empty for actual
  Codex, Antigravity, and Gemini core destinations. Codex and Gemini have the marker,
  README, license, root `SKILL.md`, both manifests, `agents/`, and `references/`;
  Antigravity adds only its established fabricated `skills/SKILL.md`.
- `agy plugin validate` against the current staged Antigravity destination passed with one
  skill and two agents processed.

### Tests, package, shipped commands, and deferred gates

- `npm test`: 253 pass, 0 fail.
- `npm pack --dry-run --json`: 153 files; 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` roots; no `addons/` entry and no root `SKILL.md`.
- The design-check help command exits 0. The tracker strict validator accepts all three
  filled templates, its single-epic command passes, and the story `claim` command succeeds
  in a scratch copy and writes `status: in-progress` plus `assignee: me`.
- The stale-path guard, reverted exemplar guard, inline/fenced/tilde command checks,
  narrowed-extractor fixture, and both Gemini reconciliation tests all turned red under
  their respective scratch mutations.
- G3 Windows execution/install remains **NOT TESTED**: this is macOS and neither Windows nor
  Wine is available.
- G8 is testable and **PASS**. I installed all five from a scratch `file://` Git remote,
  pushed a new core reference file, ran `skills update --yes`, saw all five reported
  updated, and found the new file in the installed core.
- G9 remains **NOT TESTED** for this branch. Homebrew exists, but the checked-in formula
  directory is empty and the external formula consumes an already-published release rather
  than `82507b5`.
- G10 remains **NOT TESTED** because WinGet and Wine are unavailable and this repository has
  no candidate manifest to validate.
- G11 remains **NOT TESTED** because Gemini CLI is unavailable. The README says so. Both
  static manifest consumers resolve and their mutation tests are genuine.

### README panel against A4/F1/F3/F4/F5

The panel is first, neither path is labelled recommended, no agent count is asserted,
upstream owns the supported-agent claim, and `skills list` is present and listed all five in
the scratch install. F4's core-only fallback is not triggered because the no-flag install
delivered all five. The native installer's measured differentiators are health checks and
metadata sync. I found no remaining unmeasured first-party claim in the panel.

The passing implementation claims reproduce except for the DNS-blocked
`npx ...@latest` wrapper explicitly separated above. The initial unconditional G1 wording is
superseded by fix-up cycle 1's reproducible qualification.

## Findings

### [MAJOR] A continuation before `--test` still bypasses the published-command guard

`publishedTestCommands` filters each physical line with `/node\s+--test/` at
`test/design-addons.test.js:231` before it checks for a terminal backslash at lines 235–239.
Cycle 13 therefore sees a continuation only when that first physical line already contains
both tokens.

I appended this fenced command to a shipped Markdown file in a scratch copy:

```bash
node \
  --test skills/parley-worktrees/round13-definitely-missing.test.js
```

Copied verbatim through `/bin/sh`, it exited 1 with “Could not find” the deliberate missing
test. The full `npm test` suite nevertheless stayed green at 253 pass / 0 fail. The first
line contains `node` and the continuation but not `--test`; the second contains `--test` but
not `node`, so neither line reaches the continuation refusal or the strict grammar.

This is the same false-green class as round 12 and disproves the cycle-13 implication that a
line-ending continuation now becomes a whole-line refused unit. Move continuation grouping
before the `node\s+--test` prefilter: form logical shell units first, perform backslash-newline
splicing for detection across boundaries, and refuse the complete raw group whenever the
spliced unit is a `node --test` command. Add the exact probe above plus splits inside `node`
and `--test` as regressions so the fix covers logical-line continuation rather than this one
token boundary.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
