---
idea: readme-skill-catalogue
implementer: claude-1
date: 2026-07-29
status: fix-up-cycle-4
target: parley-deck-skill/README.md
head-commit: 94a4889
prior-commits: [0061dc2 initial, 4399bd4 cycle-1, 5df9d04 cycle-2, f855a6f cycle-3]
---

## What was done

`parley-deck-skill/README.md` rewritten per `FINAL.md`. One file changed; nothing else.

- **401 → 278 lines** (binding ceiling was 300).
- The hook and all five catalogue entries are the FINAL.md copy **verbatim**, including the
  five `<!-- Base: … Grafts: … -->` attribution comments.
- Section order follows C1 exactly.
- Eight prompt blocks → **three**, with a one-line note that other requests are substitutions.

## C9 truth table — every row

| Row | Done |
|---|---|
| `:9` runtime list | now "fourteen named runtimes" enumerated, "plus `generic`, a destination you point at with `--dest`" |
| `:21-23` "append-only" | section deleted; the claim does not appear anywhere |
| `:26-27` "rates confidence by agreement" | deleted |
| `:119` "any capable tier-1 model" | deleted; replaced with "plain Markdown by design" |
| `:148-176` Repository Layout | rewritten with `addons/` + its four skills, `test/`, `packaging/`, `scripts/`, `NOTICE.md`, `RELEASING.md` |
| `:239` `v1.2.1` | version string removed; the Windows line is now versionless |
| `:242` WinGet "until accepted" | replaced with `winget install Feci.ParleyDeckSkill` |
| `:371` "all discovered installed CLI agents" | now "a bounded participant set — normally two to four, including at least one non-facilitator when one is available" |
| `:397` "value should be obvious" | deleted; the section now names the artifacts to read |
| eight prompt blocks | three kept |

## Verification

```text
$ wc -l README.md
278

$ for s in "15 runtimes" "tier-1" "should be obvious" "append-only" \
           "rates confidence" "v1.2.1" "Until the WinGet"; do grep -c "$s" README.md; done
0 0 0 0 0 0 0

$ npm test
ℹ pass 247
ℹ fail 0
```

The WinGet command was not asserted from a shipped file (C10 forbids that). It was checked
against the external registry:
`gh api repos/microsoft/winget-pkgs/contents/manifests/f/Feci/ParleyDeckSkill` →
`1.0.4 … 1.4.6`.

## Deviations from FINAL.md

**D-1 — "Why this exists" was kept, not deleted.** codex-1's cut list (round-02) and my own
round-01 both marked it for deletion, and consensus C9 did not name it either way; FINAL.md's
section order omitted it. I kept a compressed two-paragraph version (11 lines) because after
cutting to 278 lines the budget was not the binding constraint, and the failure-mode list is
the argument the hook only gestures at. **This is my judgement overriding two participants'
stated preference, so it is flagged for the review round to reverse if they disagree.**

**D-2 — section renames.** `## Local Agent Contract` → `## Local agent contract`,
`## Repository Layout` → `## Repository layout`, and the relationship/provenance sections
merged under `## Related repositories, and what this one owes`. Sentence case is now
consistent across the file. FINAL.md fixed the order, not the capitalisation.

**D-3 — the `#install-update-and-remove` anchor** in the Install section is a GitHub-generated
anchor. It renders on GitHub; in a plain terminal it is inert text. Not verified on npmjs.com.

## Not done, and deliberately

- `parley-deck-cli/README.md` — out of scope (C6), recorded as a follow-up.
- `packaging/winget/README.md` still calls its manifest a draft, which C10 identified as now
  stale. It is a separate file and a separate follow-up.
- Nothing about the `npx skills` discovery defect appears in this README (C13). It is owned by
  `skills-cli-install-path`, which is running in parallel.

---

## Fix-up cycle 1 — commit 4399bd4

Review round 01: **codex-1 ❌ BLOCK** (4 MAJOR, 2 MINOR), **hermes-1 🟡 ACCEPT-WITH-RESERVATIONS**
(2 MAJOR, 3 MINOR, 2 NIT). Zero CRITICAL. Every agreed fix applied; nothing dismissed silently.

| Finding | Reviewers | Action |
|---|---|---|
| **D-1: cut "Why this exists"** | both | **Cut.** My override is withdrawn — see below. |
| **`## Status` unreported deviation, and false for the `fast` track** | both | **Cut.** codex-1 is right that `fast` collapses `consensus.md` into `FINAL.md`, so the sentence promised artifacts a conforming run need not produce. |
| Worktree entry overclaims runtime isolation | codex-1 | Rewritten: the manifest records the per-worktree env/port/database/cache overrides that isolation actually requires. Git gives a separate tree, index, `HEAD` and branch — not isolated ports. |
| Provenance + cross-repo claims unsupported | both | Feature inventory and four of five causal attributions removed. See the note below. |
| "Installing this package installs five skills" | codex-1 | → "The installer places five skills into each detected runtime." There is no npm postinstall hook. |
| Unqualified restart claim | codex-1 | → "Restart any runtime that caches skills." |
| "AionUI-style" lineage unsupported | codex-1 | removed |

### On D-1, plainly

Two participants said cut it; I kept it and argued budget. Both reviewers ruled against me,
and hermes-1's reason is the one that lands: **the binding constraint was the section order,
not the line count**, and the section restated what the hook already says. Overridden and cut.

### On the provenance finding — I was wrong about a reviewer's finding, and both caught it

Cycle 1 claimed hermes-1's MAJOR-1 was "partly inaccurate" because `Fusion`, `ExecPlans` and
`preflight` "appear in `references/COOPERATION.md`". **That was wrong, and wrong in a
self-favourable direction.** My grep was case-insensitive (`grep -ril`), so `Fusion` matched
the lowercase idea slug `meta-protocol-change-fusion-execplans` and `preflight` matched the
command `parley preflight`. Case-sensitively:

```text
$ grep -c "Fusion" references/COOPERATION.md            → 0
$ grep -c "ExecPlans" references/COOPERATION.md         → 0
$ grep -c "Preflight readiness" references/COOPERATION.md → 0
```

hermes-1's round-01 finding was **accurate exactly as filed**. codex-1 independently reached
the same conclusion ("an internal lower-case idea slug and a command named `parley preflight`
are not those attributions"). The characterisation is withdrawn. The fix was already correct;
the justification I attached to it was not.

### Rejected: MINOR-2 (heading case)

hermes-1 says `## Use Parley Deck` breaks the sentence-case convention of D-2. It does not —
sentence case capitalises the first word, and "Parley Deck" is a proper noun. `## Use Parley
Deck` is already correct and is unchanged. This is the only finding not actioned, and the
reason is stated here rather than left silent.

### Verification

```text
$ wc -l README.md          → 260   (ceiling 300)
$ npm test                 → pass 247, fail 0
$ grep -c  "15 runtimes" "tier-1" "should be obvious" "append-only" "v1.2.1" \
           "AionUI" "Installing this package"   → 0 each
```

### Still open

- `MINOR-3` — the `#install-update-and-remove` anchor is not verified on GitHub or npmjs.com.
  It is a GitHub-generated anchor and renders as inert text elsewhere. **NOT TESTED.**
- `NIT-1` — `packaging/winget/README.md` still calls its manifest a draft. Follow-up.
- Protocol prior art in `NOTICE.md`. Follow-up.

---

## Fix-up cycle 2 — commit 5df9d04

Review round 02: **codex-1 ❌ BLOCK** (2 MAJOR, 3 MINOR, 1 NIT), **hermes-1 🟡
ACCEPT-WITH-RESERVATIONS** (1 MAJOR, 2 MINOR, 1 NIT). Both round-01 MAJORs confirmed FIXED by
both reviewers. Zero CRITICAL in either round.

| Finding | Action |
|---|---|
| **The documented project-scope install command fails** | **Real defect, reproduced.** `install --scope project --target all --project .` prints *"No installed agent runtimes were detected"* and installs nothing: under project scope, detection looks for runtime-directory evidence *inside the project*, and a clean project has none. Fixed by adding `--include-undetected`, and **verified**: the old command produces the error, the new one plans 70 skill actions (14 targets x 5 skills). |
| Cross-repo location claim still unsupported | Removed. Now: "Deterministic automated orchestration is not part of it and requires separate tooling." No named external repository. |
| "The one lineage a shipped file records" too absolute | → "The protocol lineage recorded here is…" — the next sentence names two more shipped design lineages. |
| `paths` does not report every target | Now says every **detected** target, and documents `paths --target all --include-undetected` for all fourteen. Reproduced: default returns fewer entries than the long form. |
| Runtime-channel instructions stated as fact | Qualified. The restart line is now conditional; the Gemini/Antigravity lines are explicitly marked as depending on those CLIs, not on anything this package ships. The Codex `$skill-installer` sentence is **deleted** — no shipped file records that path or what portion of a five-skill repository it installs. |
| Frontmatter stale (`commit: 0061dc2`) | Updated to the cycle-2 head with the earlier commits listed. |
| **My mischaracterisation of hermes-1's MAJOR-1** | **Withdrawn.** See the corrected section above — case-insensitive grep, wrong conclusion, caught by both reviewers. |
| NIT: worktree sentence reads like a patch splice | Rewritten as two sentences that say what git does and does not give you. |

### Verification

```text
$ wc -l README.md   → 260  (ceiling 300)
$ npm test          → pass 247, fail 0
$ node bin/parley-deck-skill.js install --scope project --target all --project /tmp/x --dry-run
  → "No installed agent runtimes were detected…"        (the documented command, before)
$ …same + --include-undetected --dry-run  → 70 skill actions across 14 targets, rc 0   (after)
$ node bin/parley-deck-skill.js paths                      → detected targets only
$ node bin/parley-deck-skill.js paths --target all --include-undetected → all fourteen
```

### Still open, unchanged

`#install-update-and-remove` anchor **NOT TESTED** on GitHub/npmjs · `packaging/winget/README.md`
stale · protocol prior art belongs in `NOTICE.md` — all three are follow-ups.

---

## Fix-up cycle 3 — commit f855a6f

Review round 03: **codex-1 ❌ BLOCK** (1 MAJOR, 2 MINOR), **hermes-1 ✅ ACCEPT** (1 MINOR, 1 NIT — see the correction in cycle 4).
Zero CRITICAL in any round.

**The MAJOR is one my own cycle-2 edit introduced.** Fixing the "uncheckable runtime
instructions" finding, I wrote *"The last two lines depend on those CLIs rather than on
anything this package ships"* — and both halves were wrong. The last two lines are the npm
and Gemini lines, and the npm line runs `parley-deck-skill install`, which is exactly what
this package ships. I also asserted the Gemini CLI and `--target gemini` "write to different
directories". They do not: `lib/installer.js:40` resolves the gemini target to
`.gemini/extensions`, which is where the Gemini CLI puts extensions too. **Same destination,
two managers** — the reason not to mix them is competing management, not divergence.

codex-1 also named that sentence as the one line in the file that reads machine-made. It was:
vague backward reference, abstract contrast, written to patch a finding rather than to say
something. That is what fix-up prose degenerates into if nobody is checking.

| Finding | Action |
|---|---|
| Gemini/channel explanation factually wrong (MAJOR) | Rewritten from the installer source: named the shared destination `~/.gemini/extensions/parley-deck`, "use one or the other, never both", and separated the Antigravity guidance instead of lumping it in |
| Restart sentence still not checkable (MINOR) | The cache assertion is gone. Now purely non-factual guidance: "follow its own instructions for reloading skills" |
| Frontmatter not in Phase-8 shape; "71 installs" wrong (MINOR) | `status: fix-up-cycle-4`, `head-commit:`, `prior-commits:`. The count is **70 skill actions** (14 targets × 5 skills); 71 was the number of *output lines*, the 71st being the add-on hint |

### Verification

```text
$ grep -n "\.gemini" lib/installer.js   → :33 .gemini/config/plugins (agy), :40 .gemini/extensions (gemini)
$ wc -l README.md                       → 260   (ceiling 300)
$ npm test                              → pass 247, fail 0
```

### Still open, unchanged

`#install-update-and-remove` anchor **NOT TESTED** on GitHub/npmjs · `packaging/winget/README.md`
stale · protocol prior art belongs in `NOTICE.md`.

---

## Fix-up cycle 4 — head-commit below

Review round 04: **codex-1 🟡 ACCEPT-WITH-RESERVATIONS** (1 MINOR, 1 NIT) — up from ❌ BLOCK
in rounds 01, 02 and 03. codex-1 also confirmed: *"No factual README claim, command argument,
link, or anchor regressed in cycle 3."*

**A correction to my own audit, and it is the reason this cycle exists.** Cycle 3 recorded
round 03 as *"hermes-1 ✅ ACCEPT (0 findings)"*. That was wrong: `review/round-03/hermes-1.md`
contains a `MINOR-1` and a `NIT-1`. I read the signoff line and not the body. **An ACCEPT
signoff does not turn filed findings into zero findings**, and recording it that way claimed
a cleaner review than actually happened while silently dropping two dispositions. The cycle-3
entry above is corrected in place, and both findings are actioned here.

| Finding | Round | Action |
|---|---|---|
| Worktree prose line is 132 chars, breaking the file's wrap convention | hermes-1 r03 MINOR-1, re-filed by codex-1 r04 NIT | Rewrapped. A sweep for prose lines >100 chars now returns none — the eleven remaining long lines are HTML attribution comments, `###` headings containing links, fenced prompt text and command lines, none of which wrap |
| `--yes` is in the general flag block but only affects `sync-project` | hermes-1 r03 NIT-1 | **Verified in source** — `lib/installer.js:403,407` gate only the `sync-project` path on `options.yes`. Moved out of the shared flag list and annotated: *"sync-project only: without it, sync-project is a dry run"* |
| Cycle-3 audit reports hermes-1 as 0 findings | codex-1 r04 MINOR | Corrected in place, above |

### Verification

```text
$ awk 'length($0)>100' README.md   → 11 lines, all comments/headings/code, no prose
$ wc -l README.md                  → 261   (ceiling 300)
$ npm test                         → pass 247, fail 0
$ grep -n 'options.yes' lib/installer.js → :403, :407   (sync-project only)
```

### Findings per round

`r01: 11 · r02: 8 · r03: 5 · r04: 2` — and no BLOCK for the first time.

### Still open, unchanged

`#install-update-and-remove` anchor **NOT TESTED** on GitHub/npmjs · `packaging/winget/README.md`
stale · protocol prior art belongs in `NOTICE.md`.
