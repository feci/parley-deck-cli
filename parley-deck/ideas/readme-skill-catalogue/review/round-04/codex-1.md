---
idea: readme-skill-catalogue
review-round: 04
agent: codex-1
date: 2026-07-29
---

## Summary

Commit `f855a6f` fixes every factual README defect I filed through round 03. The
Gemini paragraph now describes one shared destination and two competing managers,
the cache assertion is gone, the project-scope example works in an empty project,
the provenance and worktree claims are supportable, all twelve rows of consensus C9
are resolved, and the README remains at 260 lines with no broken link or anchor.
`npm test` passes 247/247.

The README is factually ready. I found one non-blocking audit-trail defect and one
formatting NIT: `IMPLEMENTATION.md` says hermes-1 had zero round-03 findings, but
hermes-1 actually filed one MINOR and one NIT, and the 132-character worktree line
named by that review remains in the shipped file.

## Refutation attempts

- Verified the target repository is clean at `f855a6f`; each of fix-up commits
  `4399bd4`, `5df9d04`, and `f855a6f` changes only `README.md`.
- Counted only `name:` records inside `const TARGETS = [...]` in
  `lib/installer.js`: **14** — `codex`, `claude`, `agy`, `gemini`, `hermes`,
  `qwen`, `codebuddy`, `goose`, `kimi`, `droid`, `vibe`, `cursor`, `opencode`,
  `aionrs`. `generic` is separate and the parser rejects it without `--dest`.
- Ran the checked-out CLI's `--help`. All seven documented command names, all
  target values, and every printed flag exist. The add-on selectors also behave
  as described: default dry-run = 14 targets × 5 skills = 70 skill actions;
  `--no-addons` = core only; `--only parley-design-check,parley-tracker` = core
  plus exactly those two add-ons.
- Ran dry-run equivalents of every mutating npx argument set. Install-all,
  force-update, preview, one-target, generic destination, include-undetected,
  project scope, and uninstall all exited 0. `doctor --target all` exited 0.
- Ran the checker with an explicitly missing registry: exit **3**. Ran it over
  `README.md` with the shipped registry: exit **4**.
- Rechecked the tracker claims against `addons/parley-tracker/SKILL.md` and its
  binaries: the canonical-file boundary, one-way default, `mirror-owned` pull
  boundary, full gap-scan, claim refusal, neutral projections, and opt-in live
  connector are accurate.
- Rechecked the worktree claims against `addons/parley-worktrees/SKILL.md`: the
  allocation table is the lock manifest; an intersecting file set is refused or
  requires a recorded override; Git supplies a separate working tree, index,
  `HEAD`, and branch, but not separate ports, databases, or caches.
- Rechecked provenance against `SKILL.md`,
  `references/COOPERATION.md` section 13, and `NOTICE.md`. Manual facilitation,
  the RHO lineage, and the `hallmark`/`impeccable` design prior art are all
  recorded there.
- Ran `npm test`: **247 passed, 0 failed**.

## Project-scope command reproduction

Both forms were run from the same newly-created empty directory with
`--dry-run --json`:

| Form | Result |
|---|---|
| Old: `install --scope project --target all --project .` | exit **1**, zero actions: “No installed agent runtimes were detected” |
| New: `install --scope project --target all --project . --include-undetected` | exit **0**, 14 target actions and 70 skill actions |

The cycle-2 project-scope fix works. The old command still fails for the exact
reason recorded by the implementer; the new command fixes it without relying on
pre-existing runtime directories.

## Documented command results

| Documented surface | Result | Evidence |
|---|---|---|
| `winget install Feci.ParleyDeckSkill` | **NOT TESTED** | `winget` is unavailable on this Mac. The shipped manifests establish the package identifier, not execution against the live Windows registry. |
| `brew install feci/parley/parley-deck-skill` | **NOT TESTED** | I ran the exact command. The sandbox denied Homebrew's Cellar, cache, and trust-lock writes before a meaningful install test. `brew info` and `brew list` show the exact formula at stable/installed version 1.5.0, but that is not a fresh-install pass. |
| `gemini extensions install https://github.com/feci/parley-deck-skill` | **NOT TESTED** | The Gemini CLI is unavailable. The shipped manifest names `parley-deck`; the official [Gemini CLI extension reference](https://github.com/google-gemini/gemini-cli/blob/main/docs/extensions/reference.md) confirms the command syntax, the `~/.gemini/extensions` root, and the manifest-name/directory-name rule, but the command itself was not executed. |
| The `npx -y parley-deck-skill@latest ...` lines | **NOT TESTED exactly; checked-out CLI PASS** | Exact npm resolution ended with `ENOTFOUND registry.npmjs.org`. Every shown argument set passed through the checked-out 1.5.0 CLI, using dry-run for mutating forms. |
| `parley-deck-skill paths` | **PASS** | Ran exactly, exit 0, six detected targets. The documented `paths --target all --include-undetected` form also exited 0 and returned all fourteen. |

The additional documented command
`agy plugin validate ~/.gemini/config/plugins/parley-deck` passed exactly.

## C9 row-by-row verification

Every row was checked against the shipped README, not accepted from
`IMPLEMENTATION.md`:

| Consensus C9 row | Result |
|---|---|
| “README is 402 lines” was wrong; old file was 401 | **RESOLVED** — `git show 94127c2:README.md \| wc -l` returns 401; current file is 260 |
| Seven prompt blocks was the wrong count; there were eight | **FIXED** — exactly three prompt blocks remain under `## Use Parley Deck` |
| Repository layout omitted `addons/` | **FIXED** — all four add-ons plus `test/`, `packaging/`, `scripts/`, `NOTICE.md`, and `RELEASING.md` are shown and exist |
| Windows example used stale `v1.2.1` | **FIXED** — the version is absent; the WinGet command is versionless |
| Lifecycle called “append-only” | **FIXED** — the phrase is absent |
| Consensus said to “rate confidence by agreement” | **FIXED** — the claim is absent |
| Universal “tier-1 model” claim | **FIXED** — the claim is absent |
| Default said “all discovered installed CLI agents” | **FIXED** — bounded two-to-four default and non-facilitator requirement are present |
| “Value should be obvious” opinion | **FIXED** — absent |
| Opening runtime list was incomplete | **FIXED** — fourteen counted runtime targets plus separate `generic` are stated |
| The old line-9/line-186 “contradiction” accusation was itself wrong | **RESOLVED** — the current text makes one complete, internally consistent enumeration |
| Stale “WinGet manifest not yet accepted” claim | **FIXED** — absent; the real package command is present |

## Round-01 findings after fix-up cycle 3

| My round-01 finding | Status | Verification against the current file |
|---|---|---|
| D-1: cut `Why this exists` | **FIXED** | The section is absent. The override was transparent but unjustified: the binding order and repetition controlled, not spare line budget. **Ruling remains cut.** |
| Unreported `Status` deviation and false artifact promise | **FIXED** | The section and its fast-track-incompatible promise are absent. |
| Worktree entry promised isolated runtime state | **FIXED** | The current copy correctly separates Git state from ports, databases, and caches. |
| Repository and provenance claims lacked shipped evidence | **FIXED** | Unsupported cross-repository inventories and causal attributions are gone; the remaining provenance is recorded by shipped files. |
| “Installing this package installs five skills” conflated package installation with running the installer | **FIXED** | The sentence now assigns the action to the installer; selector dry-runs reproduce its stated behavior. |
| Restart and “AionUI-style” operational claims were uncheckable | **FIXED** | “AionUI” and the cache/restart assertion are gone. The remaining reload sentence is troubleshooting guidance that delegates to the runtime's own instructions. |

## Required rulings

### D-1

**Cut was correct.** The original override was not justified: it violated the
binding section order and repeated the hook. The implementer withdrew the override,
and the section remains absent.

### Heading-case rejection

The implementer's rejection of hermes-1's former heading-case finding is
**sustained**. `Use` is the first word and `Parley Deck` is a proper noun, so
`## Use Parley Deck` is sentence case.

### Cycle-1 provenance rebuttal

The cycle-2 withdrawal is complete and honest. Case-sensitive searches for
`Fusion`, `ExecPlans`, and `Preflight readiness` in
`references/COOPERATION.md` each return zero. Case-insensitive searches for
`fusion` and `execplans` match only
`meta-protocol-change-fusion-execplans`; `preflight` matches only the
`parley preflight` command. Those occurrences did not support the named
prior-art attributions. `IMPLEMENTATION.md` now admits that the grep was
case-insensitive, calls the conclusion wrong and self-favouring, and states that
hermes-1's finding was accurate exactly as filed.

### Rewritten provenance and worktree copy

Both are now true. The provenance paragraph stays within the manual-facilitation
boundary and names only lineages recorded by shipped files. The worktree paragraph
accurately distinguishes Git's isolation from runtime resources that require
manifest overrides. Its remaining defect is formatting, not truth.

## Machine-made reading check

No current sentence reads as generated marketing copy. The catalogue entries have
distinct rhythms and evidence-specific refusals; the vague cycle-2 sentence about
“the last two lines” is gone. The worktree paragraph is mechanically overlong at
line 102, but its actual sentences are precise rather than slop.

## Line count and links

- `wc -l README.md` → **260**, 40 lines below the ceiling.
- All five relative skill links resolve to shipped files.
- `#install-update-and-remove` resolves to
  `## Install, update, and remove` under GitHub's heading-slug rules.
- No dead internal link or broken anchor was found.

## Findings

### [MINOR] Cycle 3 inaccurately reports hermes-1 as having zero findings

`IMPLEMENTATION.md` says round 03 was “hermes-1 ✅ ACCEPT (0 findings).” The
canonical `review/round-03/hermes-1.md` contains a `MINOR-1` about README line
102 and a `NIT-1` about the placement of `--yes`. An ACCEPT signoff does not turn
filed findings into zero findings.

This matters because the audit trail claims a cleaner review than actually
occurred and silently drops the disposition of two findings. Correct the cycle-3
summary to “1 MINOR, 1 NIT” and either action or explicitly dismiss/defer each
item in review consensus.

### [NIT] The worktree paragraph still contains a 132-character prose line

README line 102 combines the end of the Git/runtime distinction with the start
of the use-case sentence on a 132-character line. This is the exact formatting
defect hermes-1 filed in round 03; cycle 3 did not touch it.

Rewrap the paragraph to the surrounding roughly 80–95-character convention. No
wording change is required.

## What cycle 3 broke

No factual README claim, command argument, link, or anchor regressed in cycle 3.
The cycle-3 audit summary itself introduced the inaccurate “0 findings” account
described above.

### Signoff: codex-1 — 2026-07-29
Status: 🟡 ACCEPT-WITH-RESERVATIONS
