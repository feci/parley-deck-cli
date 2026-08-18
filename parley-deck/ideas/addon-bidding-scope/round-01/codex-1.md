---
agent: codex-1
idea: addon-bidding-scope
round: 1
date: 2026-08-17
---

## Summary

Recommend **FLIP THE DEFAULT**: keep `parley-bidding` in this distribution as an optional companion, but do not install it into new runtimes unless the operator selects it. A vendor-neutral package may carry an explicitly optional jurisdiction/business vertical; that does not make the vertical part of the neutral core or justify adding it to every agent's routing surface. The package already presents one core plus five companions (`../parley-deck-skill/README.md:18-29`), while the bidding manifest and skill make the specialization concrete through German discovery sources, jurisdiction data, and portal adapters (`../parley-deck-skill/skills/parley-bidding/parley-addon.json:10-17`; `../parley-deck-skill/skills/parley-bidding/SKILL.md:24-35`).

The cost if this recommendation is wrong is activation friction: a user who expected bidding to be present must opt in, and a model cannot select the skill until that happens. A conservative marker-aware migration can prevent existing users from losing it. That cost is smaller and more reversible than keeping the wrong default: an unrequested, broadly described procurement skill remains eligible for model routing, and on a Python-less runtime it advertises a workflow whose deterministic commands cannot run.

I exclude the recorded 25 fix-up cycles from the decision (`parley-deck/ideas/addon-bidding-scope/00-prompt.md`, “Reconnaissance”). They describe construction history, not present fit, demand, or migration risk.

## Proposed approach

### 1. Keep the add-on, change the runtime default

The add-on should remain packaged for now because the available usage search is insufficient evidence for `CUT`, and package co-location retains the existing common installer and manifest validation path (`../parley-deck-skill/README.md:18-29`, `:79-82`, `:200-241`). `SPLIT` would add release/discovery and stale-install migration costs without evidence that an independent cadence or audience is needed. Revisit a split if the vertical gains its own release cadence, maintainers, or multiple jurisdiction packs; the current decision need not assume those future conditions.

New installs should select the neutral core and default companions but omit `parley-bidding`. The current implementation does the opposite mechanically: absent `--no-addons` or `--only`, `selectedAddons()` returns every discovered add-on (`../parley-deck-skill/lib/installer.js:890-901`), and the core marker records that full selection (`../parley-deck-skill/lib/installer.js:1946-1964`). This is therefore an installer-default change, not a README-only change.

The CLI also needs an additive opt-in such as `--with parley-bidding` and a paired explicit removal such as `--without parley-bidding`. Current `--only` means core plus only the named add-ons (`../parley-deck-skill/lib/installer.js:963-988`), so asking a bidding user to preserve all other defaults requires enumerating the complete add-on set. Current uninstall planning also always begins with the core unit (`../parley-deck-skill/lib/installer.js:991-1003`, `:649-669`), so it is not a safe per-add-on removal interface.

### 2. Treat availability as an agent-routing effect

The README defence is narrowly sound but insufficient. Installation does not itself grant credentials or mutate a portal. The skill requires action-specific approvals (`../parley-deck-skill/skills/parley-bidding/SKILL.md:55-72`), says agents do not operate the portal (`:84-86`), and its bundled scripts do not perform portal actions (`:147-157`). In that narrow authority sense, availability is not permission.

It is not sound as a defence of the default in the agent-runtime model stipulated by this question. The frontmatter description is itself a broad routing advertisement: it matches discovery, requirements, suppliers, security, contracts, pricing, releases, staging, and most software-delivery categories (`../parley-deck-skill/skills/parley-bidding/SKILL.md:1-3`). Making that description available changes which instruction set a model may load and how it frames a tender-like task before any portal permission is relevant. Permission gates bound later effects; they do not neutralize selection, context consumption, German/EU assumptions, or an unusable tool path.

The `[!IMPORTANT]` is **not a fourth instance** of “a printed rule binds only where enforcement lives.” It is disclosure, not a rule purporting to enforce non-installation: the real default and opt-out selection live in installer code (`../parley-deck-skill/README.md:31-37`; `../parley-deck-skill/lib/installer.js:890-901`). Calling the notice itself a failed gate is pattern-matching. The valid criticism is narrower: prose about permission answers a different question from whether an unrequested skill should enter an agent's routing surface. A warning can explain that product decision, but cannot make the surface expansion neutral.

### 3. Migrate existing installations without guessing intent

Do not implement the flip by merely removing bidding from `selectedAddons()` and rewriting the marker. Existing core markers store an add-on name list, but not whether each name came from the old default or an explicit user choice (`../parley-deck-skill/lib/installer.js:904-962`, `:1958-1964`). The installer cannot reconstruct intent.

I reproduced the naïve migration by installing all add-ons, then running the current equivalent of the proposed selection, excluding only bidding:

```text
parley-deck-skill install --target generic --dest <tmp>/parley-deck
parley-deck-skill install --target generic --dest <tmp>/parley-deck --force \
  --only parley-design,parley-design-check,parley-tracker,parley-worktrees
parley-deck-skill doctor --target generic --dest <tmp>/parley-deck --json
```

Relevant output from the run:

```text
bidding_dir_exists=yes
bidding_skill_hash_before=aad7d7d8d43259893626a4593c69c3d90b8bf2a96b90dc723067fa83a04d9540
bidding_skill_hash_after=aad7d7d8d43259893626a4593c69c3d90b8bf2a96b90dc723067fa83a04d9540
core_marker_addons=["parley-design","parley-design-check","parley-tracker","parley-worktrees"]
bidding_marker_skill=parley-bidding
bidding_marker_version=2.8.0
doctor_rc=1
selected=false
status=valid-unselected
problem="installed but not part of the recorded selection: remove the directory, or re-run install including it"
```

That matches the code: an excluding install writes only selected units and does not remove the old directory; an unfiltered read later discovers the on-disk unselected add-on (`../parley-deck-skill/lib/installer.js:1085-1119`) and reports `valid-unselected` (`:2113-2147`). Thus a naïve flip neither removes the installed payload nor keeps upgrades healthy; it leaves a stale skill that the runtime may continue to expose and makes `doctor` fail.

Concrete upgrade path:

1. Fresh runtime: default install omits bidding; `--with parley-bidding` installs it and records explicit selection provenance.
2. Existing managed runtime whose core marker lists bidding: preserve and update it on the first post-flip unqualified upgrade, mark it as a migrated legacy selection, and print a one-time choice. `--with parley-bidding` affirms it; `--without parley-bidding` transactionally removes the managed sibling and rewrites the core marker.
3. Existing runtime whose marker already omits bidding: leave it omitted.
4. Existing unmanaged/foreign copy: do not delete it; the owning installer or operator must remove it. The package can validate an intact unmarked manifest copy as `valid-unmanaged` (`../parley-deck-skill/lib/installer.js:2151-2182`), but that is not deletion authority.

This grandfathering means the flip is immediate for fresh installs and conservative for existing installs. It avoids breaking real bidding users and avoids pretending the old marker proves consent. The unavoidable cost is that some old, unrequested copies remain available until the operator makes the new explicit choice.

### 4. Python-less runtime outcome

The manifest requires Python `>=3.10` (`../parley-deck-skill/skills/parley-bidding/parley-addon.json:1-5`), while the documented freeze workflow directly invokes `python3` (`../parley-deck-skill/skills/parley-bidding/SKILL.md:90-121`). I ran:

```text
rg -uuu -n -i 'python3|python 3|interpreter|fallback|without python|no python|pip|venv|uv|runtime requirement|command not found' \
  ../parley-deck-skill/skills/parley-bidding/{SKILL.md,references,scripts}
```

The relevant results were the five `python3` commands in `SKILL.md`, the seven `#!/usr/bin/env python3` script shebangs, and unrelated uses of “fallback” in bid-domain content. No interpreter-free workflow or substitute implementation appeared in `SKILL.md` or `references/` within that search.

I then installed into a generic temporary runtime with `PATH=/definitely-no-python`. The install exited 0 with `installOk: true`; `doctor` exited 1 and reported the bidding payload as byte-valid but operationally unavailable:

```json
{
  "biddingStatus": "valid",
  "biddingRuntime": {
    "ok": false,
    "requirement": ">=3.10",
    "detail": "python3 is not available, but this skill requires >=3.10"
  }
}
```

That is the implemented distinction (`../parley-deck-skill/lib/installer.js:364-390`, `:2185-2213`). But when I ran the first documented manifest command under the same PATH, the shell exited 127 with `/bin/sh: python3: command not found`. Therefore an agent that follows the skill on such a runtime can read and analyze, but cannot execute the canonical workspace/state/release toolchain; there is no documented fallback. `doctor` detects the condition only if it is run and heeded, while installation still succeeds and may expose the skill.

### 5. Widened usage evidence

This evidence affects `CUT`/`SPLIT`, not the default flip. I searched the full `/Volumes/My Shared Files/AI_WORKSPACE` tree without a depth limit:

```text
find '/Volumes/My Shared Files/AI_WORKSPACE' -type f -name 'bid-state.json' \
  ! -path '*/.git/*' ! -path '*/node_modules/*' \
  ! -path '*/parley-deck/ideas/addon-bidding-scope/round-01/*' -print
```

It returned no path. I then ran:

```text
rg -uuu -l '"event"\s*:\s*"workspace-initialized"' \
  '/Volumes/My Shared Files/AI_WORKSPACE' \
  -g '!**/.git/**' -g '!**/node_modules/**' \
  -g '!**/parley-deck/ideas/addon-bidding-scope/round-01/**'
```

It returned only `/Volumes/My Shared Files/AI_WORKSPACE/BYTE/software-bidding/scripts/init_bid_workspace.py`. Finally, bypassing ignore-aware grep behavior:

```text
/usr/bin/grep -RIl --exclude-dir=.git --exclude-dir=node_modules --exclude-dir=round-01 \
  '"event": "workspace-initialized"' \
  '/Volumes/My Shared Files/AI_WORKSPACE/BYTE' \
  '/Volumes/My Shared Files/AI_WORKSPACE/IHK_PFALZ'
```

It returned the same source script only. Within this search scope, I found no initialized workspace in the add-on's exact `bid-state.json` shape and no generated `workspace-initialized` state record. That remains absence of evidence: work could live outside this root, have been renamed, deleted, or performed manually. It supports avoiding `KEEP` as an adoption-based claim, but is not strong enough to justify `CUT` or `SPLIT`.

## Concerns / open questions

- The README explicitly says that runtime exposure as an invocable skill was not tested (`../parley-deck-skill/README.md:74-77`). Runtime variance changes blast radius, but not the default principle: a package should not rely on some runtimes hiding an installed vertical.
- The migration needs a defined duration for the “legacy selection” warning. I prefer preserving until an explicit `--with`/`--without` decision; automatically deleting later would still guess intent from inadequate markers.
- A future split criterion should be observable: independent maintainership/release cadence, multiple jurisdiction modules, or demand for installing bidding without the protocol bundle. Current repository evidence establishes none of these.
- If adding `--with`/`--without` is rejected, the minimum viable flip must document the full `--only` list and explicitly remove the stale bidding directory; changing the default alone is not a complete migration.

## Risks

- **Wrong flip:** procurement users face discoverability and opt-in friction, and an automation may assume bidding exists on a fresh runtime. The marker-preserving upgrade path limits breakage to fresh/unaffirmed installs.
- **Naïve implementation:** rewriting the core marker while leaving the sibling directory produces a stale, potentially still model-visible skill and a failing health check, as the migration run demonstrated.
- **Grandfathering:** old runtimes that never wanted bidding keep it until they choose removal. This is a deliberate safety trade-off because the current marker cannot distinguish consent from old default inheritance.
- **Over-reading usage evidence:** the widened search still covers one workspace, not every installation or past run. It should influence migration conservatism, not be stated as “never used.”
- **Premature split or cut:** removing the packaged source without an explicit cleanup path can strand old installed directories outside normal discovery/health handling, while real users take a larger installation and release-coordination break than the evidence presently warrants.
