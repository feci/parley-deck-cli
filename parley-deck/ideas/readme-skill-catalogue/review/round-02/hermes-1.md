---
idea: readme-skill-catalogue
review-round: 02
agent: hermes-1
date: 2026-07-29
---

Adversarial re-review of claude-1's README rewrite after fix-up cycle 1 (commit
4399bd4, 278→260 lines). Every claim below was checked against a shipped file or a
live command, not eyeballed. I tried to break the claim that this is done.

## Verification log (commands actually run)

- `wc -l README.md` → 260 (ceiling 300, pass, 40 lines of margin)
- `node bin/parley-deck-skill.js --help` → all 7 commands and all flags match the README
- `node bin/parley-deck-skill.js install --target all --dry-run --yes` → works, exit 0
- `node bin/parley-deck-skill.js paths` → works, lists install dirs per target
- `node bin/parley-deck-skill.js --version` → 1.5.0 (matches package.json)
- `npm view parley-deck-skill@latest version` → 1.5.0 (npm registry live)
- `npm test` → 247 pass, 0 fail (matches IMPLEMENTATION.md)
- `curl -sI https://github.com/feci/homebrew-parley` → HTTP 200 (tap exists)
- Banned-string sweep (`15 runtimes`, `tier-1`, `should be obvious`, `append-only`,
  `v1.2.1`, `AionUI`, `Installing this package`) → all 0 occurrences
- TARGETS array counted in `lib/installer.js:13-113` → 14 named entries + generic
- Exit code 4 verified at `addons/parley-design-check/lib/engine.js:2069`
- Exit code 3 verified at `addons/parley-design-check/SKILL.md:60`
- Tracker connector boundary verified at `addons/parley-tracker/SKILL.md:381-388`
- Worktree env/port/database/cache claims verified at
  `addons/parley-worktrees/SKILL.md:45,200-201,319,484-485`
- Five catalogue link targets all resolve to real files
- Repository layout tree verified against actual directory listing

Commands NOT TESTED (cannot run on this host):
- `winget install Feci.ParleyDeckSkill` — Windows-only; PackageIdentifier confirmed
  in local manifest and C10's external registry evidence (versions 1.0.4–1.4.6)
- `brew install feci/parley/parley-deck-skill` — Homebrew not available in this
  context; tap repo `feci/homebrew-parley` confirmed to exist (HTTP 200), and
  RELEASING.md:88-108 documents the tap as the formula location
- `gemini extensions install https://github.com/feci/parley-deck-skill` — requires
  Gemini CLI; the extension mode is confirmed in `lib/installer.js:38-43` and
  `gemini-extension.json`
- `npx -y parley-deck-skill@latest install --target all` — run via local binary
  instead; npm package confirmed live at 1.5.0

## Round-01 findings — disposition after fix-up cycle 1

| Round-01 finding | Status | Verification |
|---|---|---|
| MAJOR-1 (provenance attributions uncheckable) | FIXED | The five-attribution block is gone. README:251-256 now states only the RHO lineage that COOPERATION.md:1078 actually records, plus the hallmark/impeccable prior art that NOTICE.md:13,17 records. Verified both against the shipped files. |
| MAJOR-2 (`## Status` unreported deviation) | FIXED | The `## Status` section is gone from the shipped file. `grep -n "## Status" README.md` → no hits. The file now has exactly the C1 section order (with the D-2 renames). |
| MINOR-1 (D-1: cut "Why this exists") | FIXED | "Why this exists" is gone. `grep -n "Why this exists" README.md` → no hits. The implementer withdrew the override (IMPLEMENTATION.md "On D-1, plainly") and cut it. |
| MINOR-2 (heading case: `## Use Parley Deck`) | REJECTED by implementer — ruling below | The implementer rejected this finding. The heading is unchanged at README:115. See the ruling section below. |
| MINOR-3 (`#install-update-and-remove` anchor) | NOT FIXED, correctly deferred | The anchor remains. IMPLEMENTATION.md "Still open" records it as NOT TESTED. It is a GitHub-generated anchor that renders as inert text elsewhere. Not a defect; correctly left as a known limitation. |
| NIT-1 (winget README still says "draft") | NOT FIXED, correctly deferred | Out of scope (separate file). IMPLEMENTATION.md "Still open" records it as a follow-up. |
| NIT-2 (gemini line trailing comment) | NOT FIXED | The trailing `# legacy Gemini only` comment is still at README:190. This was a usability note, not a truth defect. Acceptable to leave. |

## Ruling on the implementer's rejection of MINOR-2

The implementer rejected my MINOR-2, arguing "## Use Parley Deck" is already
sentence case because "Parley Deck" is a proper noun and "Use" is the first word.

**Ruling: the rejection is sustained.** Sentence case capitalises the first word
and proper nouns; "Use" is the first word and is correctly capitalised, and
"Parley Deck" is a proper noun. `## Use Parley Deck` is sentence case. My
round-01 finding conflated "Title Case" with "the first word is capitalised" —
the first word being capitalised is correct in sentence case. The heading is
consistent with `## Local agent contract`, `## Repository layout`, and `## Related
repositories, and what this one owes` (all sentence case). I withdraw MINOR-2.

## Ruling on the implementer's claim about MAJOR-1

The implementer's IMPLEMENTATION.md "On the provenance finding" states that my
MAJOR-1 was "partly inaccurate" because "`Fusion`, `ExecPlans` and `preflight` all
appear in `references/COOPERATION.md` (at `:1108` and `:751`)."

**Ruling: the implementer's claim about my claim is itself inaccurate.** I ran
`grep -r "Fusion\|ExecPlans"` across the entire repo: zero hits for either string.
COOPERATION.md:1108 contains the idea slug `meta-protocol-change-fusion-execplans`
(lowercase, part of a hyphenated slug), not the attribution labels "Fusion" and
"ExecPlans" that the round-01 README used. The string "preflight" does appear at
COOPERATION.md:751, but as `parley preflight` (a tooling command name), not as the
"Preflight readiness" attribution label the README used.

My MAJOR-1 stated that the four attributions "appear nowhere in the repo outside
this README." That was correct: "OpenRouter Fusion", "OpenAI ExecPlans / PLANS.md",
"kindly", and "Preflight readiness" (as attribution labels) do not appear in any
shipped file. The implementer's counter-evidence is a slug fragment, not an
attribution. However, the implementer's conclusion was correct regardless: the fix
(following codex-1's stronger version) is to keep only the RHO lineage a shipped
file records and delete the rest. The fix is right; the argument used to reach it
contained an inaccurate characterisation of my finding. This is a process finding,
not a file defect — see MAJOR-1 below.

## Ruling on DEVIATION D-1 (keep or cut "Why this exists")

**Cut.** Both reviewers (myself and codex-1) ruled against keeping it. The
implementer withdrew the override and cut it. This is resolved. The section is
absent from the shipped file. No further action needed.

## CRITICAL

None.

## MAJOR

### MAJOR-1 — IMPLEMENTATION.md mischaracterises my round-01 MAJOR-1 as "partly inaccurate"

IMPLEMENTATION.md "On the provenance finding" claims my MAJOR-1 was partly
inaccurate because "`Fusion`, `ExecPlans` and `preflight` all appear in
`references/COOPERATION.md` (at `:1108` and `:751`)."

I verified this against the shipped file. `grep -r "Fusion\|ExecPlans"` across the
entire repo returns zero hits. COOPERATION.md:1108 contains the idea slug
`meta-protocol-change-fusion-execplans` — a lowercase hyphenated slug, not the
attribution labels "Fusion" and "ExecPlans". "preflight" appears at :751 as
`parley preflight` (a tooling command), not as the "Preflight readiness" label.

My round-01 MAJOR-1 was accurate: the four attribution labels appeared nowhere in
any shipped file. The implementer's counter-claim is a slug fragment dressed up as
an attribution. The fix itself is correct (codex-1's version: keep only RHO,
delete the rest), but the implementer's justification for it misrepresents my
finding.

Why this is MAJOR and not MINOR: the review process depends on accurate
characterisations of reviewer findings. An implementer who rebutts a finding with
evidence that does not support the rebuttal, and the rebuttal is recorded as
justification for the fix, undermines the audit trail. The shipped README is
correct; the IMPLEMENTATION.md narrative around it is not. The "On the provenance
finding" section should be corrected to state that my finding was accurate and the
fix follows codex-1's stronger version, without the false counter-claim about
Fusion/ExecPlans/preflight appearing as attributions.

## MINOR

### MINOR-1 — Worktree sentence runs past the line-width of the surrounding paragraphs

README:100-103:

    Each implementer gets a sibling worktree, and the manifest records the
    per-worktree environment, port, database and cache overrides that isolating runtime
    state actually requires. Use it when two or more sessions or Phase-5 implementers work in one repository
    at once.

The final sentence ("Use it when two or more sessions or Phase-5 implementers work
in one repository at once.") is 108 characters on a single line, while every other
prose line in the file wraps at ~75–80 characters. This is a cosmetic wrap defect,
not a truth defect. The content is accurate — verified against
`addons/parley-worktrees/SKILL.md:6` ("Use when two or more agents/sessions touch
the same repo concurrently"). Fix: rewrap the sentence to match the file's line
width.

### MINOR-2 — The `--yes` flag appears in the README's flag block but not in any documented command

README:180 lists `--yes` in the flags block:
    `--force  --dry-run  --yes  --json  --include-undetected`

The `--help` output confirms `--yes` exists as a flag (it appears in the install,
doctor, status, uninstall, and paths command help, plus sync-project). However,
none of the eight documented command examples in the README use `--yes`. It is
listed but never shown in action. This is not a truth defect (the flag exists),
but a reader sees it in the flags block with no example of when to use it. Minor
documentation gap; the flag is correctly documented as existing.

## NIT

### NIT-1 — `## Use Parley Deck` is the only section heading that is an imperative verb phrase

Every other section heading is a noun phrase: "What's in the box", "Install",
"Install, update, and remove", "Local agent contract", "Transports", "Repository
layout", "Related repositories, and what this one owes", "License". "Use Parley
Deck" is an imperative verb phrase. This is a style observation, not a defect —
the heading is clear and the casing is correct (see MINOR-2 ruling above). No
action required.

## C9 truth table — row-by-row re-verification

Every row checked against the shipped README (260 lines), not against
IMPLEMENTATION.md's claim:

| Row | Fixed in file? | Evidence |
|---|---|---|
| `:9` runtime list | YES | README:147-150 says "fourteen named runtimes" + lists them + "plus `generic`, a destination you point at with `--dest`." Counted 14 in `lib/installer.js:13-113` (codex, claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe, cursor, opencode, aionrs). "15 runtimes" absent. |
| `:21-23` "append-only" | YES | `grep -c "append-only" README.md` → 0. |
| `:26-27` "rates confidence by agreement" | YES | `grep -ci "rates confidence" README.md` → 0. |
| `:119` "any capable tier-1 model" | YES | `grep -c "tier-1" README.md` → 0. Replaced with "plain Markdown by design" (README:139). |
| `:148-176` Repository Layout | YES | Rewritten tree (README:224-242) includes `addons/` with four skills, `test/`, `packaging/`, `scripts/`, `NOTICE.md`, `RELEASING.md`. All verified against actual directory listing. |
| `:239` `v1.2.1` | YES | `grep -c "v1.2.1" README.md` → 0. Windows line (188) is versionless. |
| `:242` WinGet "until accepted" | YES | Replaced with `winget install Feci.ParleyDeckSkill` (README:188). |
| `:371` "all discovered installed CLI agents" | YES | Now "a bounded participant set — normally two to four, including at least one non-facilitator when one is available" (README:208-209). |
| `:397` "value should be obvious" | YES | `grep -c "should be obvious" README.md` → 0. Section deleted. |
| Eight prompt blocks → three | YES | Three prompt blocks at README:119, 124, 129. |

All ten rows fixed in the file, not just claimed in IMPLEMENTATION.md.

## Factual sentence audit (CHECK 1)

Every factual sentence in the new README was checked against a shipped file or live
command:

- "The installer places five skills into each detected runtime" (README:19) —
  verified: installer copies 5 skill payloads (core + 4 add-ons); `--no-addons`
  and `--only` flags confirmed in `--help`.
- "fourteen named runtimes — Codex, Claude Code, Antigravity CLI (plugin mode),
  legacy Gemini CLI (extension mode), Hermes, Qwen, CodeBuddy, Goose, Kimi,
  Factory Droid, Vibe, Cursor, OpenCode and AionRS — plus `generic`" (README:147-150)
  — verified: 14 TARGETS entries in `lib/installer.js:13-113`; `generic` requires
  `--dest` per `--help`.
- "exits `3`" / "exits `4`, not `0`" (README:67, 68-69) — verified at
  `addons/parley-design-check/SKILL.md:60` (exit 3) and
  `lib/engine.js:2069` (exit 4).
- "live create/update requires an opt-in connector" (README:85-86) — verified at
  `addons/parley-tracker/SKILL.md:381-388`.
- "Sync is one-way by default, and pull reconciliation may write back only fields
  declared `mirror-owned`" (README:82-84) — verified at
  `addons/parley-tracker/SKILL.md:27,29,135`.
- "The skill defines neutral projections for Jira, Linear, GitHub Issues, GitLab,
  Trello, and plain boards" (README:84-85) — verified at
  `addons/parley-tracker/SKILL.md:3,14,23-24`.
- "The worktree-allocation table in `IMPLEMENTATION.md` is the lock manifest"
  (README:97) — verified at `addons/parley-worktrees/SKILL.md:340-341`.
- "an intersection is refused unless an explicit override is recorded" (README:99-100)
  — verified at `addons/parley-worktrees/SKILL.md:55,325-328,393-395`.
- "the manifest records the per-worktree environment, port, database and cache
  overrides that isolating runtime state actually requires" (README:100-102) —
  verified at `addons/parley-worktrees/SKILL.md:45,200-201,319,484-485`.
- "RHO (Retrospective Harness Optimization), credited in
  `references/COOPERATION.md` §13" (README:252-253) — verified: §13 is at
  COOPERATION.md:1076, RHO credited at :1078.
- "RHO's single-model self-preference is deliberately replaced by the deck's
  multi-agent quorum" (README:253-254) — verified: COOPERATION.md:1078 says
  "RHO's single-model self-preference is replaced here by the deck's multi-agent
  quorum."
- "`NOTICE.md` records `hallmark` and `impeccable` as the prior art studied for
  the design add-ons" (README:254-255) — verified at NOTICE.md:13,17.
- "manifest.yaml (neutral), openai.yaml (Codex UI metadata)" (README:237) —
  verified: `manifest.yaml` has `vendor_neutral: true`; `openai.yaml` has
  `display_name`, `short_description`, `default_prompt` (Codex UI metadata).
- Commands list (README:175): `install`, `paths`, `doctor`, `status`,
  `sync-project`, `uninstall`, `--version` — verified against `--help`: all 7
  present, exact match.
- Flag block (README:178-181): all flags verified against `--help` output. The
  `--target` value list matches exactly, character for character.

No uncheckable factual sentence found in the shipped README.

## Machine-made reading check (CHECK 5)

The README does not read as generated. The five catalogue entries retain their
distinct voices (codex-1's precision on artifact ownership, kimi-1's compression
into refusals, the "no stack trace" framing). The hook is direct. The install
section is terse and operational. The provenance section is now specific and
evidence-anchored rather than hand-waving.

No sentence reads as machine-made slop. The repeated colon-explanation pattern I
noted in round 01 (which was the closest thing to a machine tell) is gone — the
"AionUI-style local runtime registry" sentence and its colon-explanations were
removed in the fix-up. The file reads as written by someone who knows the codebase.

## Line count and links (CHECK 6)

- `wc -l README.md` → 260. Ceiling 300. Pass with 40 lines of margin.
- Five catalogue links (`./SKILL.md`, `./addons/parley-design/SKILL.md`, etc.) —
  all resolve to real files (verified with `[ -f ]`).
- One internal anchor link: `[Install, update, and remove](#install-update-and-remove)`
  at README:113 → heading `## Install, update, and remove` at README:141 → GitHub
  anchor `#install-update-and-remove`. Correct on GitHub; inert text elsewhere
  (correctly flagged as NOT TESTED in IMPLEMENTATION.md).
- No dead internal links. No broken anchors detected.

## What the fix-up broke

Nothing. The fix-up cycle removed three sections (Why this exists, Status,
AionUI lineage), rewrote the provenance and worktree sentences, and changed one
installer sentence. I checked every section that was touched and every section
adjacent to a touch. No new factual error, no broken link, no orphaned reference
was introduced. The file is internally consistent and every factual sentence
checks out.

### Signoff: hermes-1 — 2026-07-29
Status: 🟡 ACCEPT-WITH-RESERVATIONS
