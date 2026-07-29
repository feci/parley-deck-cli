---
idea: readme-skill-catalogue
review-round: 01
agent: codex-1
date: 2026-07-29
---

## Summary

Commit `0061dc2` is not ready to accept. The five-skill catalogue is substantially accurate,
all ten binding C9 repairs are present, the binding catalogue copy matches `FINAL.md`
byte-for-byte, and the mechanical checks pass. The remaining defects matter because factual
restraint and a non-machine-made voice are the idea's central acceptance criteria, not
optional polish.

## Refutation attempts

- Counted `name:` entries inside the `TARGETS` array in `lib/installer.js` with a script:
  **14**, in the same order and with the same names as README lines 157–160. `generic` is
  separately validated by the CLI as requiring `--dest`.
- Ran the shipped CLI's root and subcommand `--help`. Every subcommand and flag printed at
  README lines 162–192 exists. Dry-run equivalents of the mutating install, force-update,
  single-target, project-scope, generic-destination, include-undetected, and uninstall
  examples all returned `ok: true`.
- Exercised the add-on selectors in dry-run mode. The default result was the core plus all
  four add-ons; `--no-addons` returned only `parley-deck`; `--only
  parley-design-check,parley-tracker` returned the core plus exactly those two add-ons.
- Ran `parley-deck-skill doctor --target all`; all five installed skills were valid in every
  detected runtime. Ran `parley-deck-skill paths` exactly; it exited 0 and listed each
  detected target's core and add-on destinations.
- Ran the checker with an explicitly missing registry: exit **3**. Ran it over an input it
  could not judge with a valid registry: exit **4** and `UNJUDGEABLE`. README lines 67–70 are
  correct.
- Checked the tracker copy against `addons/parley-tracker/SKILL.md` and its two binaries.
  The canonical-file/mirror split, one-way default, `mirror-owned` pull boundary, gap-scan,
  claim refusal, and opt-in live connector boundary are all stated correctly.
- Checked the worktree copy against `addons/parley-worktrees/SKILL.md`; one absolute runtime
  isolation claim failed, as reported below.
- Ran the complete C9 forbidden-string and required-layout scan. Every C9 row is fixed:
  three prompt fences remain; the stale phrases and version are absent; the bounded roster
  wording is present; and the layout includes all required add-ons, `test/`, `NOTICE.md`,
  and `RELEASING.md`.
- Ran `npm test`: **247 passed, 0 failed**. Ran `npm pack --dry-run` with an isolated cache:
  the tarball contains the core `SKILL.md` and all four add-on `SKILL.md` files.
- Counted **278 lines**, below the 300-line ceiling. Parsed every Markdown link and generated
  heading slug: all five relative skill links and `#install-update-and-remove` resolve; no
  dead internal link or broken anchor was found.

## Documented command results

| Documented surface | Result | Evidence |
|---|---|---|
| `winget install Feci.ParleyDeckSkill` | **NOT TESTED** | `winget` is unavailable on this Mac. The shipped manifests establish the package ID but cannot execute the registry install. |
| `brew install feci/parley/parley-deck-skill` | **NOT TESTED** exactly | `brew info` found the exact formula, stable 1.5.0, already installed. Even `brew install --dry-run` attempted an out-of-workspace trust lock and was denied. The chained skill installer passed in dry-run mode. |
| `gemini extensions install https://github.com/feci/parley-deck-skill` | **NOT TESTED** | `gemini` is unavailable on this Mac. |
| The `npx -y parley-deck-skill@latest ...` lines | **NOT TESTED** exactly | npm registry resolution failed with `ENOTFOUND`. The installed 1.5.0 binary accepted every shown flag; `doctor` passed exactly and every mutating example passed with `--dry-run`. |
| `parley-deck-skill paths` | **PASS** | Ran exactly; exit 0 and valid paths returned. |

The additional `agy plugin validate ~/.gemini/config/plugins/parley-deck` command passed
exactly. The global npm install was not run because it is a machine mutation; `npm install
--help` confirms the shown `-g` syntax.

## Findings

### [MAJOR] D-1 is not justified: cut `Why this exists`

README lines 140–149 violate the binding section order in `FINAL.md` C1. More importantly,
the implementer's reason for overriding both reviewers is only that 278 lines left room.
The two cut proposals were based on repetition and voice, not merely the ceiling: the hook,
core entry, and ownership copy already make this argument.

The retained opening sentence is also the clearest sentence in the README that reads as
generated:

> Multi-agent workflows fail in predictable ways: one agent anchors the rest before they
> form their own view, disagreements dissolve into a long chat history, implementation
> starts before there is real consensus, reviews are informal and unowned, and vendor
> assumptions leak into the workflow.

It is a stock thesis followed by five evenly abstract failure clauses, and it repeats claims
the reader has just read. **Ruling: cut the section.** The override was transparent, but it
was not justified.

### [MAJOR] An unreported `Status` deviation violates C1 and makes a false artifact promise

README lines 270–274 retain a `## Status` section even though the binding order proceeds
directly from repository relationship/provenance to License. Unlike D-1, this deviation is
not recorded in `IMPLEMENTATION.md`; the implementation instead claims that the section
order follows C1 exactly.

The retained sentence is also not true for every shipped track:

> `ideas/<slug>/` will contain ... the consensus with its signoffs, `FINAL.md`, and
> `IMPLEMENTATION.md`.

`references/COOPERATION.md` defines the `fast` track as a collapsed `FINAL.md` with embedded
signoffs, not a separate `consensus.md`. A discussion that has reached design finalization
also has not necessarily entered Phase 5 and produced `IMPLEMENTATION.md`. Delete the Status
section, as the round-02 cut list required.

### [MAJOR] The worktree entry promises runtime isolation that git does not provide

README line 100 says, “Each implementer gets a sibling worktree and isolated runtime state.”
The shipped worktree skill is more careful: it guarantees a separate working tree, index,
`HEAD`, and branch, but says runtime files are per-worktree only when they live inside that
tree (`addons/parley-worktrees/SKILL.md:61-78`). Ports, database names, and cache directories
need explicit per-worktree overrides (`:195-202`, `:484-486`), while refs under the common
`.git` remain shared.

Replace the absolute claim with the actual contract: each implementer gets a sibling
worktree, and the manifest records the per-worktree environment, port, database, and cache
overrides needed to isolate runtime state.

### [MAJOR] Repository and provenance claims have no shipped-file evidence

README lines 256–267 say what the separate `parley-deck` and `parley-deck-cli` repositories
contain and attribute specific protocol features to OpenRouter Fusion, OpenAI ExecPlans,
RHO, kindly, and Preflight readiness. The RHO attribution is supported by
`references/COOPERATION.md`; the cross-repository feature inventory and the other causal
mappings are not established by a shipped file in this repository. `SKILL.md` supports only
the narrower manual-facilitation boundary. `NOTICE.md` records Hallmark and Impeccable as
prior art for the design add-ons; it does not record most of the sources named in this README
paragraph.

This directly fails the requirement that every factual sentence be checkable against a
shipped file. Keep the manual-facilitation statement supported by `SKILL.md` and the narrow
`NOTICE.md` statement, or remove the unsupported feature inventory and causal attributions.

### [MINOR] “Installing this package” conflates package installation with running the installer

README line 19 says, “Installing this package installs five skills.” `package.json` has no
install/postinstall hook. `npm install -g` or Homebrew installs the package and executable;
the subsequent `parley-deck-skill install` command places the five skills into agent runtime
directories. The later channel examples correctly show both steps, which makes the earlier
sentence avoidably misleading. Change it to “The installer installs five skills by
default.”

### [MINOR] Two operational claims are not checkable against shipped files

README lines 111 and 153 respectively require a runtime restart “so it reloads `SKILL.md`”
and call the detection model “AionUI-style.” No shipped file establishes a universal reload
requirement across fourteen runtimes, and no shipped file establishes the AionUI lineage.
The installer code does establish the behavior that follows the attribution.

Qualify the first sentence (“Restart runtimes that cache skills”) and remove
“AionUI-style,” leaving the directly checkable directory/CLI/marker behavior.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
