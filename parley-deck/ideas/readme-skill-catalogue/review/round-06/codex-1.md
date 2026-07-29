---
idea: readme-skill-catalogue
review-round: 06
agent: codex-1
date: 2026-07-29
---

## Summary

Adversarial re-review of commit `94a4889` found no remaining defect. The README is 261
lines, every consensus C9 row is resolved in the file, all six findings from my round-01
review are fixed, the cycle-4 changes address both round-04 findings, and no factual claim,
CLI argument, link, anchor, or prose-wrap convention regressed.

The target worktree contained unrelated unstaged work that moves the skill payloads into a
new layout. I did not attribute that later work to this idea: source and test checks were run
from an isolated archive of the requested review commit, while external command probes used
the installed tools.

## Refutation attempts

- Counted `name:` records only inside `const TARGETS = [...]` at `94a4889`: **14** —
  `codex`, `claude`, `agy`, `gemini`, `hermes`, `qwen`, `codebuddy`, `goose`, `kimi`,
  `droid`, `vibe`, `cursor`, `opencode`, `aionrs`. `generic` is separate and rejects a
  missing `--dest`.
- Ran the root CLI and every subcommand with `--help`. All documented commands, target
  values, and flags exist. `--yes` appears only on the `sync-project` usage line, matching
  its implementation at `lib/installer.js:403,407`.
- Exercised add-on selection over all fourteen targets in dry-run JSON mode: default =
  **70** skill actions and all five skills; `--no-addons` = **14** actions and only
  `parley-deck`; `--only parley-design-check,parley-tracker` = **42** actions and the core
  plus exactly those two add-ons.
- Ran safe equivalents of every mutating installer example. Install-all, force update,
  preview, one target, generic destination, include-undetected, and uninstall all exited
  zero. `doctor --target all` also exited zero.
- Ran both project-scope forms from the same empty directory. The old command without
  `--include-undetected` exited **1** with zero actions. The documented replacement exited
  **0**, planned 14 target records and 70 skill actions, and required no pre-existing
  runtime directory. The cycle-2 fix works.
- Ran the design checker with an explicitly missing registry: exit **3**. Ran it over
  `README.md` with the shipped registry: exit **4** and `UNJUDGEABLE`.
- Rechecked the tracker boundary against its shipped skill: canonical markdown, one-way
  default sync, `mirror-owned` pull reconciliation, complete gap reporting, claim refusal,
  neutral projections, and the opt-in live connector are all stated accurately.
- Rechecked the worktree copy: the allocation table is the lock manifest; intersecting file
  sets are refused or require a recorded override; Git supplies a separate working tree,
  index, `HEAD`, and branch, but not separate ports, databases, or caches. The manifest
  carries those per-worktree overrides.
- Rechecked provenance against `SKILL.md`, `references/COOPERATION.md` section 13, and
  `NOTICE.md`. Manual facilitation, RHO lineage, and the `hallmark`/`impeccable` design
  prior art are all recorded by those shipped files.
- Ran `npm test` from the commit snapshot: **247 passed, 0 failed**.

## Documented command results

| Documented surface | Result | Evidence |
|---|---|---|
| `winget install Feci.ParleyDeckSkill` | **NOT TESTED** | Ran exactly; this Mac has no `winget` (`127`, command not found). The shipped manifests do establish `Feci.ParleyDeckSkill`. |
| `brew install feci/parley/parley-deck-skill` | **NOT TESTED** | Ran exactly; the sandbox denied Cellar, cache, and trust-lock writes before Homebrew could perform a meaningful install test. |
| `gemini extensions install https://github.com/feci/parley-deck-skill` | **NOT TESTED** | Ran exactly; this host has no `gemini` CLI (`127`, command not found). The shipped extension manifest and installer support the documented extension mode and destination. |
| The `npx -y parley-deck-skill@latest ...` lines | **NOT TESTED exactly; checked-out CLI PASS** | `npx` package resolution failed with `ENOTFOUND registry.npmjs.org`. Every shown argument set passed against the 1.5.0 CLI from `94a4889`, with dry-run added to mutating forms. |
| `parley-deck-skill paths` | **PASS** | Ran exactly; exit zero and five valid skill destinations for each of six detected targets. The long form with `--target all --include-undetected` also exited zero and returned all fourteen targets. |

The additional `agy plugin validate ~/.gemini/config/plugins/parley-deck` command passed
exactly. The global npm install was not executed because it is a machine mutation;
`npm install --help` confirms the documented `-g` syntax.

## Round-01 findings after fix-up cycle 4

| My round-01 finding | Status | Verification against `94a4889` |
|---|---|---|
| D-1: cut `Why this exists` | **FIXED** | The section is absent. The original override was transparent but unjustified: the binding order and repetition controlled, not spare line budget. **Ruling remains cut.** |
| Unreported `Status` deviation and false artifact promise | **FIXED** | The section and its fast-track-incompatible promise are absent. |
| Worktree entry promised isolated runtime state | **FIXED** | The current copy separates Git state from ports, databases, and caches, and assigns runtime overrides to the manifest. |
| Repository and provenance claims lacked shipped evidence | **FIXED** | Unsupported cross-repository inventories and causal attributions are gone; every remaining provenance claim is recorded by a shipped file. |
| “Installing this package installs five skills” conflated package installation with running the installer | **FIXED** | The sentence assigns the action to the installer; selector dry-runs reproduce the stated default and overrides. |
| Restart and “AionUI-style” operational claims were uncheckable | **FIXED** | Both assertions are gone. The remaining reload sentence delegates to the runtime's own instructions and makes no cache claim. |

No round-01 finding is partially fixed or not fixed.

## Required rulings

### D-1

**Cut was correct.** Keeping `Why this exists` overrode two reviewers, violated the binding
section order, and repeated the hook. The implementer withdrew that override in cycle 1, and
the section remains absent.

### Heading-case rejection

The implementer's rejection of hermes-1's heading-case finding is **sustained**. `Use` is
the first word and `Parley Deck` is a proper noun, so `## Use Parley Deck` is already
sentence case.

### Cycle-1 provenance rebuttal and cycle-2 withdrawal

The original claim that hermes-1's MAJOR-1 was partly inaccurate was wrong.
Case-sensitive searches in `references/COOPERATION.md` return zero for `Fusion`,
`ExecPlans`, and `Preflight readiness`. Case-insensitive searches find only the internal
slug `meta-protocol-change-fusion-execplans` and the command `parley preflight`; neither
supports the removed attributions.

The withdrawal is complete and honest. `IMPLEMENTATION.md` identifies the case-insensitive
grep, calls the conclusion wrong and self-favouring, states that hermes-1's finding was
accurate exactly as filed, and preserves the already-correct README fix.

### Rewritten provenance and worktree copy

Both are true. Provenance stays within the shipped manual-facilitation, RHO, and
design-prior-art evidence. The worktree sentence accurately distinguishes Git's isolated
state from runtime resources requiring explicit overrides.

### Fix-up cycle 4

Both prior findings were actioned. The worktree paragraph has no prose line over 100
characters, and `--yes` is separated from the shared flags and correctly annotated as
`sync-project`-only. The cycle-3 audit now records hermes-1's one MINOR and one NIT rather
than claiming zero findings. Cycle 4 broke nothing in the reviewed commit.

## Consensus C9 — row-by-row verification

| C9 row | Result in the shipped README |
|---|---|
| Old README called 402 lines | **RESOLVED** — the pre-rewrite file at `94127c2` is 401 lines; the current README makes no false historical claim. |
| Seven prompt blocks was the wrong count; there were eight | **FIXED** — the old file has eight and exactly three remain under `## Use Parley Deck`. |
| Repository layout omitted `addons/` | **FIXED** — all four add-ons, `test/`, `packaging/`, `scripts/`, `NOTICE.md`, and `RELEASING.md` are listed and exist at the reviewed commit. |
| Windows example used stale `v1.2.1` | **FIXED** — the version is absent and the WinGet command is versionless. |
| Lifecycle called “append-only” | **FIXED** — the phrase is absent. |
| Consensus said to “rate confidence by agreement” | **FIXED** — the claim is absent. |
| Universal “tier-1 model” claim | **FIXED** — the claim is absent. |
| Default said “all discovered installed CLI agents” | **FIXED** — the README states the bounded two-to-four default and non-facilitator requirement. |
| “Value should be obvious” opinion | **FIXED** — absent. |
| Opening runtime list was incomplete | **FIXED** — all fourteen counted targets are enumerated, with `generic` correctly separate. |
| Old line-9/line-186 “contradiction” accusation was itself wrong | **RESOLVED** — the current file contains one complete, internally consistent enumeration. |
| Stale “WinGet manifest not yet accepted” claim | **FIXED** — absent; the published package identifier is used in the command. |

Every C9 row is fixed or resolved in the file, not merely claimed fixed in
`IMPLEMENTATION.md`.

## Machine-made reading check

No current sentence reads as generated marketing copy. The entries vary in rhythm and make
specific, checkable refusals; the vague cycle-2 channel sentence is gone. I do not name a
sentence here because inventing a “machine-made” example after failing to find one would
violate the review's refutation discipline.

## Line count and links

- `wc -l README.md` equivalent at `94a4889`: **261**, 39 lines below the ceiling.
- All five relative catalogue links resolve to files in the reviewed commit.
- `#install-update-and-remove` resolves to `## Install, update, and remove` under GitHub's
  heading-slug rules.
- Exactly three prompt fences remain under `## Use Parley Deck`.
- No dead internal link, broken anchor, or prose line over 100 characters was found.

## Findings

None.

### Signoff: codex-1 — 2026-07-29
Status: ✅ ACCEPT
