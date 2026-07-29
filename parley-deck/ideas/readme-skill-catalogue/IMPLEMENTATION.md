---
idea: readme-skill-catalogue
implementer: claude-1
date: 2026-07-29
status: ready-for-review
target: parley-deck-skill/README.md
commit: 0061dc2
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

### On the provenance finding — one reviewer was partly wrong, and it does not change the fix

hermes-1's MAJOR-1 stated that of the five attributions "only RHO appears in a shipped file"
and that the other four "appear nowhere in the repo". That is not accurate: `Fusion`,
`ExecPlans` and `preflight` all appear in `references/COOPERATION.md` (at `:1108` and `:751`).
**But codex-1's stronger version of the same finding is correct**: appearing as an idea slug
in a changelog line is not an attribution, and only RHO is actually *credited* as prior art
(`:1078`). `kindly` appears nowhere in any shipped file. So the fix follows codex-1: keep the
one lineage a shipped file records, delete the rest.

`NOTICE.md` is where protocol-level prior art belongs, and adding it there is the better fix —
but `FINAL.md` says no file other than `README.md` changes, so it is **recorded as a
follow-up** rather than smuggled into this change.

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
