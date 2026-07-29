---
idea: readme-skill-catalogue
review-round: 02
agent: codex-1
date: 2026-07-29
---

## Summary

Commit `4399bd4` is not done. The fix-up correctly cut `Why this exists` and `Status`,
narrowed the worktree copy, removed the AionUI attribution, and fixed the package/installer
conflation. All ten C9 rows are fixed, the catalogue claims I retested are supported, all
links resolve, the README is 260 lines, and `npm test` passes 247/247.

The completion claim still fails on two central requirements: one documented install command
exits 1 in a clean project, and the provenance section still assigns automated orchestration
to external repositories without evidence in any shipped file. The exact `paths` command also
returns detected targets, not “every target” as its sentence claims.

## Refutation attempts

- Counted `name:` entries only inside `const TARGETS = [...]` in `lib/installer.js` with a
  Node script: **14**, in README order (`codex`, `claude`, `agy`, `gemini`, `hermes`, `qwen`,
  `codebuddy`, `goose`, `kimi`, `droid`, `vibe`, `cursor`, `opencode`, `aionrs`). `generic`
  is separate and requires `--dest`.
- Ran the root and every subcommand `--help`. All seven documented commands, all fourteen
  target values plus `generic`, and all listed flags exist. The help text confirms
  `--no-addons` means core only and `--only` means core plus the named add-ons.
- Ran dry-run forms of every mutating npx argument set against the shipped 1.5.0 source.
  Install-all, force, preview, one-target, generic destination, include-undetected, and
  uninstall all exited 0. The documented project-scope form exited **1**; adding
  `--include-undetected` made it exit 0 and enumerate all fourteen target paths.
- Ran the add-on selectors. Default install planned five skills per detected runtime;
  `--no-addons` planned only `parley-deck`; `--only
  parley-design-check,parley-tracker` planned the core plus exactly those two.
- Ran `parley-deck-skill doctor --target all`: exit 0, with all five installed skills valid
  for every detected runtime.
- Ran the design checker with an explicitly missing registry: exit **3**. Ran it over
  `README.md` with the shipped registry: exit **4**, with `UNJUDGEABLE`.
- Checked the tracker connector, one-way sync, `mirror-owned`, gap-scan, and claim boundaries
  against `addons/parley-tracker/SKILL.md` and its tests. The README copy is accurate.
- Checked the revised worktree copy against `addons/parley-worktrees/SKILL.md`, including the
  sibling layout, lock manifest, intersection rule, and per-worktree runtime overrides. The
  factual overclaim from round 01 is gone.
- Searched every shipped file except README for `parley-deck-cli`, the separate
  `parley-deck` server, and an assignment of automated orchestration to either repository:
  no supporting statement exists. `SKILL.md` supports only the narrower manual-facilitation
  boundary.
- Verified the provenance dispute lexically and semantically. The lower-case fragment
  `fusion-execplans` occurs only inside the internal idea slug
  `meta-protocol-change-fusion-execplans`; `preflight` occurs as the local command
  `parley preflight`. Neither is evidence for the removed OpenRouter Fusion, OpenAI
  ExecPlans/PLANS.md, or “Preflight readiness” prior-art attributions.
- Ran `npm test`: **247 passed, 0 failed**.
- Counted **260 lines**, below the 300-line ceiling. Parsed all Markdown links and generated
  heading slugs: five relative skill links and `#install-update-and-remove` resolve. There is
  no dead internal link or broken anchor.

## Round-01 findings after fix-up cycle 1

| My round-01 finding | Status | Re-verification |
|---|---|---|
| D-1: cut `Why this exists` | **FIXED** | The section is absent. The original override was transparent but unjustified: section order and repetition, not spare line budget, were decisive. **Ruling remains cut.** |
| Unreported `Status` deviation and false artifact promise | **FIXED** | The section and its fast-track-incompatible promise are absent. |
| Worktree entry promises isolated runtime state | **FIXED** | The new copy distinguishes Git isolation from the env/port/database/cache overrides runtime isolation needs. |
| Repository and provenance claims lack shipped evidence | **PARTIALLY FIXED** | The feature inventory and four unsupported causal attributions are gone. The new sentence still claims automated orchestration “lives in” two external repositories, which no shipped file establishes. |
| “Installing this package installs five skills” conflates package and skill installation | **FIXED** | The copy now correctly assigns the five-skill action to the installer. Default and selector behavior were retested. |
| Restart and AionUI-style operational claims are uncheckable | **PARTIALLY FIXED** | “AionUI-style” is gone. “Restart any runtime that caches skills” remains an uncheckable cross-runtime instruction; no shipped file defines those cache/restart semantics. |

## Response to hermes-1/round-02

- I agree with hermes-1’s withdrawal of the heading-case finding and sustain the
  implementer’s rejection. `Use` is the first word and `Parley Deck` is a proper noun;
  `## Use Parley Deck` is sentence case.
- I agree with hermes-1 on the substance of the IMPLEMENTATION.md provenance rebuttal:
  a lower-case idea-slug fragment and a local `parley preflight` command do not make the
  removed named prior-art attributions appear in shipped files. I rate the inaccurate
  characterization **MINOR**, not MAJOR, because the implementer applied the right README
  fix and the error is now confined to the audit narrative.
- I disagree with “No uncheckable factual sentence found” and “What the fix-up broke:
  Nothing.” The project-scope command exits 1 with the shipped CLI in a clean checkout, and
  plain `paths` returned 6 detected targets rather than all 14. The worktree fix also
  introduced the file’s only conspicuously overlong prose line.
- I do not accept hermes-1’s `--yes` documentation gap as a finding. A compact flag inventory
  promises that a flag exists, not that every flag has an example. `--help` establishes the
  flag and its command context.
- I agree that the worktree sentence should be rewrapped. I also find its “that isolating
  runtime state actually requires” splice and the provenance phrase “The one lineage a
  shipped file records” read like review-cycle patch language rather than finished README
  prose.

## Documented command results

| Documented surface | Result | Evidence |
|---|---|---|
| `winget install Feci.ParleyDeckSkill` | **NOT TESTED** | `winget` is unavailable on this Mac. Shipped manifests establish the package ID, not a live registry install. |
| `brew install feci/parley/parley-deck-skill` | **NOT TESTED** | I ran the exact command, but the sandbox denied Homebrew’s Cellar/cache/trust-lock writes before a meaningful install test. `brew list --versions` reports `parley-deck-skill 1.5.0`, and its installed binary runs, but that is not an exact fresh-install test. |
| `gemini extensions install https://github.com/feci/parley-deck-skill` | **NOT TESTED** | `gemini` is unavailable. `gemini-extension.json` proves repository shape, not execution of the external CLI. |
| The `npx -y parley-deck-skill@latest ...` lines | **FAIL / PARTIALLY TESTED** | Exact npm resolution was unavailable (`ENOTFOUND registry.npmjs.org`), so the npx transport is NOT TESTED exactly. The installed and checked-out 1.5.0 CLI accepted every flag. All dry-run examples passed except `install --scope project --target all --project .`, which exited 1 in a clean project. |
| `parley-deck-skill paths` | **PASS execution; FAIL described scope** | Ran exactly, exit 0. It listed six detected targets. `paths --target all --include-undetected` listed all fourteen, so “every target” is false for the documented plain command. |

The additional `agy plugin validate ~/.gemini/config/plugins/parley-deck` command passed
exactly.

## C9 truth table

Every row was checked against the shipped README rather than IMPLEMENTATION.md:

| C9 row | Result |
|---|---|
| Runtime enumeration | **FIXED** — fourteen counted names plus `generic`; “15 runtimes” absent |
| Lifecycle “append-only” | **FIXED** — absent |
| “rates confidence by agreement” | **FIXED** — absent |
| Universal “tier-1 model” claim | **FIXED** — absent |
| Repository layout | **FIXED** — four add-ons, `test/`, `packaging/`, `scripts/`, `NOTICE.md`, and `RELEASING.md` present |
| Stale `v1.2.1` | **FIXED** — absent |
| Stale WinGet “until accepted” | **FIXED** — absent; current command present |
| “all discovered installed CLI agents” | **FIXED** — bounded two-to-four default present |
| “value should be obvious” | **FIXED** — absent |
| Eight prompt blocks | **FIXED** — exactly three in `Use Parley Deck` |

## Findings

### [MAJOR] The documented project-scope install command fails in a clean project

README line 165 presents this as the project-scope member of three working examples:

```bash
npx -y parley-deck-skill@latest install --scope project --target all --project .
```

Running the shipped CLI with the same arguments plus `--dry-run --json` exits 1:

```text
No installed agent runtimes were detected. Use --target all --include-undetected
or --target generic --dest <path>.
```

This is deterministic, not a network or sandbox limitation. Under project scope,
`isRuntimeDetected()` checks for pre-existing runtime-directory evidence inside the project;
an ordinary clean project has none. The command succeeds only after such evidence exists or
when `--include-undetected` is added.

The README’s central promise is copy-pasteable, fact-checked installation. Change the example
to include `--include-undetected`, or document a specific target such as
`--scope project --target codex --project .`.

### [MAJOR] The provenance repair still makes an unsupported cross-repository claim

README lines 246–249 correctly say this skill performs manual facilitation and is not itself
deterministic orchestration. The sentence then exceeds the shipped evidence:

> Deterministic automated orchestration is not part of it — that lives in the separate
> `parley-deck` server and `parley-deck-cli` repositories.

No shipped file identifies those repositories or assigns that responsibility to them.
`SKILL.md` says only that this skill “is not a deterministic A2A facilitator service by
itself.” This is the same evidentiary boundary as round 01: a true negative about this
package does not prove a positive inventory of external repositories.

Keep the supported first clause and delete the location claim, or state only that automated
orchestration requires separate tooling. Also change “The one lineage a shipped file
records” to “The protocol lineage recorded here”: the next sentence itself names two more
shipped design lineages, so “one lineage” is needlessly absolute.

### [MINOR] `paths` does not report “every target” without `--include-undetected`

README line 198 says:

> Run `parley-deck-skill paths` for the install directory of every target.

The exact command exited 0 but returned only the six detected targets on this machine:
Codex, Claude, Antigravity, Gemini, Hermes, and Kimi. The source defaults `paths` to
`--target auto`; it does not enumerate undetected targets. With
`--target all --include-undetected`, it returned all fourteen.

Say “every detected target,” or document the longer command when all supported paths are
intended.

### [MINOR] Several runtime-channel instructions remain uncheckable against shipped files

The following operational claims have no shipped-file contract:

- “Restart any runtime that caches skills.”
- “Use either the Gemini extension command or `--target gemini`, not both.”
- “Codex users can also use the built-in `$skill-installer` with the repository URL, then
  restart Codex.”

The first still does not establish which runtimes cache skills or that restart is their
required reload mechanism. The second depends on external Gemini CLI naming/install
behavior, and that CLI was unavailable. The third depends on an external Codex system skill;
no file shipped in this package records that path or what portion of this five-skill
repository it installs.

These may be useful advice, but the idea explicitly requires every factual instruction to be
checkable. Qualify them as runtime-dependent possibilities, link to shipped evidence, or
remove them.

### [MINOR] The fix-up audit is not internally current

`IMPLEMENTATION.md` frontmatter still says `status: ready-for-review` and `commit: 0061dc2`;
the Phase-8 section names `4399bd4`, but the protocol requires the top-level status and head
commit to advance for a fix-up cycle. A reader or tool consuming frontmatter sees the
superseded implementation.

The same section says hermes-1 was partly inaccurate because `Fusion`, `ExecPlans`, and
`preflight` “appear.” Hermes-1’s finding concerned the named prior-art attributions. An
internal lower-case idea slug and a command named `parley preflight` are not those
attributions. Correct the frontmatter and remove that characterization; the sound statement
is already present immediately afterward: occurrence is not attribution, and only RHO is
credited at protocol level.

### [NIT] The worktree fix reads like a patch splice

This sentence is factually supportable but is the clearest machine-edited line in the
README:

> Each implementer gets a sibling worktree, and the manifest records the per-worktree
> environment, port, database and cache overrides that isolating runtime state actually
> requires. Use it when two or more sessions or Phase-5 implementers work in one repository
> at once.

The stacked noun list and “actually requires” sound like a review rebuttal pasted into
catalogue copy, and the fix-up left the `Use it...` clause on a 104-character line while the
surrounding prose is wrapped consistently. Split it into two direct sentences and rewrap it.
The provenance phrase “The one lineage a shipped file records” has the same audit-generated
voice; the MAJOR provenance fix above supplies the cleaner wording.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
