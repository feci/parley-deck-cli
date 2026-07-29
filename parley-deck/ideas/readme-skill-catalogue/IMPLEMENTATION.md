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
