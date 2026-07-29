---
idea: readme-skill-catalogue
review-round: 03
agent: codex-1
date: 2026-07-29
---

## Summary

Commit `5df9d04` fixes the two cycle-2 blockers: the documented project-scope command now
works in an empty project, and the unsupported cross-repository orchestration claim is gone.
The worktree and provenance rewrites are true, all ten C9 rows remain fixed, the README is
260 lines, all internal links resolve, and `npm test` passes 247/247.

The implementation is not done. Fix-up cycle 2 introduced a false explanation beside the
Gemini install command, and the restart sentence remains a factual runtime claim with no
shipped-file evidence. The cycle-2 audit record also still misstates its phase metadata and
counts an informational hint as an install.

## Refutation attempts

- Verified the target repository is clean at `5df9d04`.
- Ran the root CLI and every documented subcommand with `--help`: `install`, `doctor`,
  `status`, `sync-project`, `uninstall`, and `paths` all exited 0. `--version` returned
  `1.5.0`. Every target and flag printed in the README is accepted by the parser.
- Counted only `name:` entries inside `const TARGETS = [...]` in `lib/installer.js`:
  **14**, in the README's order. `generic` is separate and rejects a missing `--dest`.
- Ran dry-run equivalents of every mutating npx argument set. Install-all, force update,
  preview, one-target, generic destination, include-undetected, project scope, and
  uninstall all exited 0. `doctor --target all` also exited 0.
- Ran the add-on selectors. Default install planned the core plus all four add-ons;
  `--no-addons` planned only `parley-deck`; `--only
  parley-design-check,parley-tracker` planned the core plus exactly those two add-ons.
- Ran the design checker with an explicitly missing registry: exit **3**. Ran it over
  `README.md` with the shipped registry: exit **4** with `UNJUDGEABLE`.
- Checked the tracker canonical-file, one-way sync, `mirror-owned`, gap-scan, claim refusal,
  neutral projection, and opt-in connector statements against
  `addons/parley-tracker/SKILL.md` and its binaries. They are accurate.
- Checked the worktree lock manifest, sibling layout, file-set intersection rule, Git
  isolation boundary, and runtime-override statements against
  `addons/parley-worktrees/SKILL.md`. The rewritten sentence is accurate.
- Checked the provenance paragraph against `SKILL.md`, `references/COOPERATION.md` §13, and
  `NOTICE.md`. The manual-facilitation boundary, RHO lineage, and design-add-on prior art are
  all supported.
- Ran `npm test`: **247 passed, 0 failed**.
- Counted **260 lines**, below the 300-line ceiling. Parsed all six Markdown links and
  generated the heading slugs: the five relative skill links and
  `#install-update-and-remove` all resolve.

## Project-scope command reproduction

Both forms were run from the same newly-created empty project with `--dry-run --json`:

| Form | Result |
|---|---|
| Old: `install --scope project --target all --project .` | exit **1**, 0 target actions: “No installed agent runtimes were detected” |
| New: `install --scope project --target all --project . --include-undetected` | exit **0**, 14 target actions and 70 skill actions |

The README fix works. The ordinary text renderer prints 71 lines only because it appends one
informational add-on hint after the 70 planned skill actions; it does not plan 71 installs.

## Documented command results

| Documented surface | Result | Evidence |
|---|---|---|
| `winget install Feci.ParleyDeckSkill` | **NOT TESTED** | `winget` is unavailable on this Mac. The shipped manifests establish the package ID, not execution against the Windows registry. |
| `brew install feci/parley/parley-deck-skill` | **NOT TESTED** | I ran the exact command. Homebrew exited 1 before a meaningful formula test because the sandbox denied Cellar, cache, and trust-lock writes. `brew list --versions` reports the already-installed `parley-deck-skill 1.5.0`, but that is not a fresh install test. |
| `gemini extensions install https://github.com/feci/parley-deck-skill` | **NOT TESTED** | `gemini` is unavailable. The official CLI documentation confirms the syntax, not execution on this host. |
| The `npx -y parley-deck-skill@latest ...` lines | **NOT TESTED exactly; local argument sets PASS** | npm resolution failed with `ENOTFOUND registry.npmjs.org`. The checked-out and installed 1.5.0 CLI accepted all shown arguments, and all safe dry-run equivalents passed. |
| `parley-deck-skill paths` | **PASS** | Exact command exited 0 and reported 6 detected targets. The documented `paths --target all --include-undetected` form exited 0 and reported all 14. |

The additional `agy plugin validate ~/.gemini/config/plugins/parley-deck` command passed
exactly.

## Round-01 findings after fix-up cycle 2

| My round-01 finding | Status | Re-verification against the current file |
|---|---|---|
| D-1: cut `Why this exists` | **FIXED** | The section is absent. The override was transparent but unjustified: binding order and repetition, not spare line budget, controlled. **Ruling remains cut.** |
| Unreported `Status` deviation and false artifact promise | **FIXED** | The section and its fast-track-incompatible promise remain absent. |
| Worktree entry promises isolated runtime state | **FIXED** | The current copy correctly separates Git's working-tree/index/`HEAD`/branch isolation from ports, databases, and caches, which need manifest overrides. |
| Repository and provenance claims lack shipped evidence | **FIXED** | Cycle 2 removed the remaining unsupported external-repository location claim. The current manual-facilitation and lineage statements trace to shipped files. |
| “Installing this package installs five skills” conflates package installation with running the installer | **FIXED** | The copy assigns the action to the installer, whose default and selector behavior I re-ran. |
| Restart and “AionUI-style” operational claims are uncheckable | **PARTIALLY FIXED** | “AionUI-style” is gone. The restart copy is conditional now, but “Some runtimes cache skills” is still a factual sentence for which no shipped file supplies evidence. |

## C9 row-by-row re-verification

Every row was checked against the current README, not against `IMPLEMENTATION.md`:

| C9 row | Result |
|---|---|
| Runtime enumeration | **FIXED** — 14 counted target names plus separate `generic`; “15 runtimes” absent |
| Lifecycle “append-only” | **FIXED** — absent |
| “rates confidence by agreement” | **FIXED** — absent |
| Universal “tier-1 model” claim | **FIXED** — absent |
| Repository layout | **FIXED** — all four add-ons, `test/`, `packaging/`, `scripts/`, `NOTICE.md`, and `RELEASING.md` exist |
| Stale `v1.2.1` | **FIXED** — absent |
| Stale WinGet “until accepted” | **FIXED** — absent; current package command present |
| “all discovered installed CLI agents” | **FIXED** — bounded two-to-four default present |
| “value should be obvious” | **FIXED** — absent |
| Eight prompt blocks | **FIXED** — exactly three under `## Use Parley Deck` |

## Required rulings

### D-1

**Cut was the correct ruling.** The implementer's original override was not justified because
the omitted section repeated the binding hook and violated the binding section order. The
override is now fully withdrawn and the section is absent.

### Heading-case rejection

The rejection of hermes-1's former MINOR-2 is **sustained**. `Use` is the first word and
`Parley Deck` is a proper noun, so `## Use Parley Deck` is sentence case. Hermes-1's
withdrawal was correct.

### Withdrawal of the cycle-1 provenance rebuttal

The withdrawal is complete and honest. Case-sensitive counts for `Fusion`, `ExecPlans`, and
`Preflight readiness` in `references/COOPERATION.md` are all 0; case-insensitive searches
each return 1 because they match the lower-case idea slug or the `parley preflight` command.
Those matches did not support the named prior-art attributions. The current
`IMPLEMENTATION.md` explicitly admits the case-insensitive grep, calls the conclusion wrong
and self-favouring, and says hermes-1's finding was accurate exactly as filed.

## Findings

### [MAJOR] Cycle 2 made the external-channel explanation factually wrong

README lines 193–196 say:

> The last two lines depend on those CLIs rather than on anything this package ships. Use
> either the Gemini extension command or `--target gemini`, not both — they write to
> different directories.

“The last two lines” are the npm and Gemini lines, not Gemini and Antigravity. The npm line
ends by running `parley-deck-skill install`, which is precisely the executable this package
ships, so the first sentence's contrast is false.

The directory explanation is also false. The shipped Gemini target resolves to
`<home>/.gemini/extensions/parley-deck` (`lib/installer.js`), and the shipped
`gemini-extension.json` names the extension `parley-deck`. The official
[Gemini CLI extension reference](https://github.com/google-gemini/gemini-cli/blob/main/docs/extensions/reference.md)
says extensions live under `<home>/.gemini/extensions` and expects the manifest name to match
the extension directory. The two mechanisms target the same extension directory; the reason
not to mix them is competing management of one destination, not different destinations.

This is exactly the kind of install-path sentence this idea was created to make factual.
Delete the vague “last two lines” sentence. State that the Gemini CLI command and
`--target gemini` are alternative managers for the same destination, then keep the
separate `--target agy` guidance for Antigravity.

### [MINOR] The restart sentence is still not checkable against a shipped file

README line 112 says, “Some runtimes cache skills; restart yours if it does not pick the
change up.” A search of every shipped file other than README finds no cache, reload, or
restart contract for any of the fourteen runtimes. The sentence may be useful operational
advice, but this idea's explicit rule is that every factual sentence must be checkable
against shipped evidence.

Remove the cache assertion, name and cite a specific runtime contract, or reduce it to
non-factual troubleshooting guidance such as consulting that runtime's reload instructions
if the change is not detected.

### [MINOR] The cycle-2 implementation audit is still internally inaccurate

`IMPLEMENTATION.md` frontmatter remains `status: ready-for-review` with a custom `commit:`
field. Phase 8 requires the top-level status to advance to `fix-up-cycle-2` and the head
commit to advance in `head-commit:`. Cycle 2 updated the SHA value but did not finish the
metadata repair raised in round 02.

The same cycle says the new project-scope form “plans 71 installs.” The JSON result is 14
targets × 5 skills = **70** planned skill actions. The human renderer has 71 output lines
only because line 71 is the explanatory add-on hint.

Update the frontmatter to the Phase-8 shape and say “70 skill actions” or “71 output lines,”
not “71 installs.”

## Machine-made reading check

The catalogue as a whole does not read as machine-made: the five entries have different
rhythms, concrete refusals, and source-specific detail. The sentence that does read as
generated is the one already filed above: “The last two lines depend on those CLIs rather
than on anything this package ships.” Its vague backward reference and abstract contrast
read like fix-up connective prose, and the referent became wrong when the preceding list
changed.

No other sentence crosses the anti-slop threshold.

## Line count and links

- `wc -l README.md` equivalent: **260**, 40 below the 300-line ceiling.
- All five relative skill links resolve to shipped files.
- `#install-update-and-remove` resolves to `## Install, update, and remove` under GitHub's
  heading slug rules.
- No dead internal link or broken anchor was found.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
