---
idea: addon-bidding-scope
author: claude-1
created: 2026-08-17
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: final
track: standard
---

## Problem / idea

The owner is reviewing the five companion add-ons in `parley-deck-skill` one at a time, ordered by
evidence of use, asking what should change. Reconnaissance on **`parley-bidding`** returned a
question that is not "how do we improve it".

**Does a jurisdiction-bound, Python-carrying business vertical belong in this package at all — and
if it does, should it install by default?**

The unit of output is a recommendation among (at least) these, with the cost of being wrong for each:

- **KEEP** — default-on, as today.
- **FLIP THE DEFAULT** — stays in the package, opt-in rather than opt-out.
- **SPLIT** — its own package, depending on `parley-deck`.
- **CUT** — removed from this package entirely.
- something none of us thought of.

**A recommendation to keep it exactly as it is, is a legitimate outcome** provided it answers the
evidence below. Do not manufacture a problem, and do not treat "it took a lot of work" as a reason
either to keep or to remove it.

## Reconnaissance — all PRIMARY, verified 2026-08-17, re-runnable

**It is categorically unlike the other four add-ons.** The other four are process tools with no
domain and no jurisdiction. `parley-bidding` is a German/EU public-procurement vertical:

- `assets/jurisdiction-profiles/de.json`, `references/jurisdiction-de.md`
- platform adapters `cosinex-vmp.dtvp.json`, `cosinex-vmp.nrw.json`, `subreport-elvis.json`
- discovery sources `ted.json`, `service-bund.json`, `oeffentlichevergabe.json`

**It is the only add-on with a runtime dependency.** Seven Python scripts
(`manifest.py`, `bid_state.py`, `adapter_validate.py`, `init_bid_workspace.py`, `common.py`,
`release_lint.py`, `completeness_lint.py`) plus a Python test suite. Meanwhile:

- `parley-design`'s own description advertises **"zero runtime dependencies"**;
- `parley-design-check`'s advertises **"no runtime dependencies and no network access"**.

**It installs by default, including on upgrade — and the README already says this is a problem:**

> **`parley-bidding` installs by default, including on upgrade.** A routine `install --force` places
> a procurement-portal skill into every runtime this installer covers — including runtimes belonging
> to people who never asked for a bidding tool. What expands is *availability*, not permission:
> every gate in the skill still binds, it performs no portal action without an action-specific human
> approval, and it never handles credentials. Use `--no-addons`, or `--only` without it, to leave it
> out.

That paragraph is a **documented workaround for a default nobody opted into**. This deck has now
recorded three instances of the class *a printed rule binds only where enforcement lives* (the
printed fix-up cap of 2 that ran 15 cycles; the review-round-1 independence property that exists only
in the runner; `COOPERATION.md:531`'s cross-reviewer obligation at 7% compliance across 348 files).
**A README warning is not a gate.** Whether this is a fourth instance is a question for you, not a
premise — argue it either way.

**Install footprint** (`du -sk skills/*`, KB):

| add-on | KB | share |
| --- | ---: | ---: |
| `parley-design-check` | 615 | 50% |
| `parley-bidding` | 252 | 21% |
| `parley-deck` (core) | 169 | 14% |
| `parley-tracker` | 95 | 8% |
| `parley-design` | 67 | 5% |
| `parley-worktrees` | 27 | 2% |
| **total** | **1224** | |

**71% of what lands on every runtime is the two add-ons with the weakest evidence of use.**

**Evidence of use — stated carefully, because absence of evidence is not evidence of absence.**
A `find` over the workspace (max depth 5–6, `node_modules` excluded) located:

- `parley-worktrees` — four live `wt-*` trees. **Used.**
- `parley-tracker` — `tickets/` in five projects. **Used.**
- `parley-design` / `parley-design-check` — `design-system/` directories exist in two projects, but
  a grep for `PDS/1.0` / `parley-design` markers in them returned nothing, so they are probably not
  this add-on's output. **No evidence either way.**
- `parley-bidding` — real tender material exists (`BYTE/`, `IHK_PFALZ/`), and both projects have
  their own `parley-deck/` decks. But `BYTE/software-bidding` turned out to be **a copy of the skill
  itself** (its `SKILL.md` was the only grep hit), not output. **No artifact in this skill's own
  shape was found.** My search was depth-limited; treat this as absence of evidence.

**It was the most expensive thing ever built in this deck: 25 recorded fix-up cycles**
(`ideas/integrate-parley-bidding-addon/IMPLEMENTATION.md`), against a `standard`-track printed cap of
2. Do not read that as an argument in either direction without saying why.

**One thing it demonstrably bought the whole package.** Its integration forced the add-on manifest
question, and the resolution was generic: `validateInstalledPayload` now runs
`manifestProblems(root, sourceHasManifest)` (`lib/installer.js:2388`) rather than a name-keyed rule,
and **all six skills now ship `parley-addon.json`**. The blocking fork was escalated, the user chose
a third option not in `FINAL.md`, and because that design was unreviewed it went to a targeted
`round-03/` where all three participants accepted it with conditions. That is the protocol working.

## Constraints

- **This is not a protocol change.** `COOPERATION.md` is not touched by this idea.
- No new dependency, no new service.
- Provenance (§15.2): every claim about this repository needs a PRIMARY locator or a command you
  ran, or a `RECALL` tag. **Do not assert a negative without running the check** — and note that my
  own usage searches were depth-limited, so if you want to claim "never used", widen the search
  yourself and say how.
- Judge the *decision*, not the sunk cost. 25 fix-up cycles is a fact about how it was built, not
  evidence about whether it should stay.
- If you recommend removal or a default flip, state the **migration cost** concretely: who breaks,
  what the upgrade path is, and what the installed markers on existing runtimes do.

## Non-goals

- Improving the bidding workflow itself, its adapters, or its jurisdiction coverage. If the answer
  is KEEP, feature work is a separate successor idea.
- Deciding anything about `parley-design-check`'s 615 KB — that add-on gets its own recon and its
  own idea. Note it only where it changes the argument about `parley-bidding`.
- Re-litigating the manifest fork. It is closed and it worked.
