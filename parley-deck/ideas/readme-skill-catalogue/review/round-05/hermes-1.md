---
idea: readme-skill-catalogue
review-round: 05
agent: hermes-1
date: 2026-07-29
---

Adversarial re-review of claude-1's README rewrite after fix-up cycle 4 (commit
94a4889, 261 lines). Every claim below was checked against a shipped file or a
live command, not eyeballed. I tried to break the claim that this is done.

## Verification log (commands actually run)

- `wc -l README.md` → 261 (ceiling 300, pass, 39 lines of margin)
- `node bin/parley-deck-skill.js --help` → all 7 commands and all flags match
  the README's documented set; `--target` value list matches character for
  character against the installer's `targetList` (lib/installer.js:316)
- `node bin/parley-deck-skill.js --version` → 1.5.0 (matches package.json)
- `npm view parley-deck-skill@latest version` → 1.5.0 (npm registry live)
- `npm test` → 247 pass, 0 fail
- `node bin/parley-deck-skill.js install --target all --dry-run` → exit 0
- `node bin/parley-deck-skill.js install --target all --include-undetected
  --dry-run` → 70 skill actions (14 targets x 5 skills), 71 lines total
- `node bin/parley-deck-skill.js doctor --target all` → exit 0
- `node bin/parley-deck-skill.js status --target all` → exit 0
- `node bin/parley-deck-skill.js paths` → 6 detected targets (codex, claude,
  agy, gemini, hermes, kimi)
- `node bin/parley-deck-skill.js paths --target all --include-undetected` → 14
  unique targets
- OLD project-scope command: `install --scope project --target all --project
  /tmp/pd-r05-test --dry-run` → exit 1, "No installed agent runtimes were
  detected" (reproduces the defect the cycle-2 fix addressed)
- NEW project-scope command: `install --scope project --target all --project
  /tmp/pd-r05-test --include-undetected --dry-run` → exit 0, 70 skill actions
- Banned-string sweep (`15 runtimes`, `tier-1`, `should be obvious`,
  `append-only`, `rates confidence`, `v1.2.1`, `Until the WinGet`, `AionUI`,
  `Installing this package`, `Why this exists`) → all 0 occurrences
- TARGETS array counted in `lib/installer.js:13-113` → 14 named entries (codex,
  claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe,
  cursor, opencode, aionrs); `generic` is separate, requires `--dest`
  (lib/installer.js:258-259)
- Exit code 3 verified at `addons/parley-design-check/SKILL.md:60`
- Exit code 4 verified at `addons/parley-design-check/SKILL.md:61` and
  `addons/parley-design-check/lib/engine.js:2069`
- Tracker connector boundary verified at
  `addons/parley-tracker/SKILL.md:381-388`
- `mirror-owned` fields verified at `addons/parley-tracker/SKILL.md:29,135`
- Tracker projections (Jira, Linear, GitHub Issues, GitLab, Trello, plain
  boards) verified at `addons/parley-tracker/SKILL.md:3,14,23-24`
- Worktree lock manifest verified at
  `addons/parley-worktrees/SKILL.md:340-341`
- Worktree intersection refusal verified at
  `addons/parley-worktrees/SKILL.md:55,325-328,393-395`
- Worktree git isolation (own working tree, index, HEAD, branch) verified at
  `addons/parley-worktrees/SKILL.md:62-63,72`
- Worktree env/port/database/cache overrides verified at
  `addons/parley-worktrees/SKILL.md:45,201,319`
- RHO lineage verified at `references/COOPERATION.md:1076-1078` (§13)
- `hallmark` and `impeccable` prior art verified at `NOTICE.md:13,17`
- Bounded participant set verified at `SKILL.md:102,283`
- 30-minute timeout verified at `SKILL.md:288`
- Marker-only directory exclusion verified at `lib/installer.js:716-718`
- `--yes` flag scoped to sync-project only: verified at
  `lib/installer.js:403,407`; `--help` shows `[--yes]` only on sync-project
- `manifest.yaml` vendor_neutral and `openai.yaml` display_name verified at
  `agents/manifest.yaml:3`, `agents/openai.yaml:2`
- `plugin.json` and `gemini-extension.json` both exist as shipped files
- parley-design refusals (no score, no house look, no "good default aesthetic")
  verified at `addons/parley-design/SKILL.md:38,40-41`
- Bounded grafts (zero to three, winner's token file not modified) verified at
  `addons/parley-design/references/PDS.md:198,331-332,383-385`
- "no runtime, no network, no framework" verified at
  `addons/parley-design/SKILL.md:11`
- design-check "no fallback registry" verified at
  `addons/parley-design-check/SKILL.md:173`
- design-check capability declaration from detectors verified at
  `addons/parley-design-check/SKILL.md:187`
- design-check `rule-id — violation — remedy` format verified at
  `addons/parley-design-check/SKILL.md:12`
- Tracker section names verified at
  `addons/parley-tracker/SKILL.md:99,104-106`
- Tracker gap-scan and claim refusal verified at
  `addons/parley-tracker/SKILL.md:212,247,414-415`
- Standalone Windows binaries verified at `RELEASING.md:47`
- Homebrew tap confirmed: `curl -sI https://github.com/feci/homebrew-parley` →
  HTTP 200; formula at `Formula/parley-deck-skill.rb` (v1.5.0)
- WinGet PackageIdentifier confirmed in local manifest at
  `packaging/winget/manifests/f/Feci/ParleyDeckSkill/`
- Five catalogue link targets all resolve to real files
- Repository layout tree verified against actual directory listing
- Gemini shared destination: README says `~/.gemini/extensions/parley-deck`;
  `lib/installer.js:40` resolves gemini to `.gemini/extensions` — match
- Prose line length sweep: `awk 'length($0)>100'` → 11 lines, all HTML
  attribution comments, `###` headings containing links, fenced prompt text,
  and command/flag lines. No prose line exceeds 90 characters.

Commands NOT TESTED (cannot run on this host):
- `winget install Feci.ParleyDeckSkill` — Windows-only; PackageIdentifier
  confirmed in local manifest and C10's external registry evidence
- `brew install feci/parley/parley-deck-skill` — Homebrew not available in this
  context; tap repo confirmed (HTTP 200), formula confirmed at
  `Formula/parley-deck-skill.rb` (v1.5.0)
- `gemini extensions install https://github.com/feci/parley-deck-skill` —
  requires Gemini CLI; extension mode confirmed in `lib/installer.js:38-43` and
  `gemini-extension.json`
- `npx -y parley-deck-skill@latest install --target all` — run via local binary
  instead; npm package confirmed live at 1.5.0

---

## Round-01 findings — disposition after fix-up cycle 4

| Round-01 finding | Status | Verification |
|---|---|---|
| MAJOR-1 (provenance attributions uncheckable) | FIXED | The five-attribution block is gone. README:252-257 now states only the RHO lineage that COOPERATION.md:1078 records, plus the hallmark/impeccable prior art that NOTICE.md:13,17 records. Case-sensitive grep for "Fusion", "ExecPlans", "Preflight readiness" in COOPERATION.md → 0 each. |
| MAJOR-2 (`## Status` unreported deviation) | FIXED | `grep -n "## Status" README.md` → no hits. The file follows the C1 section order. |
| MINOR-1 (D-1: cut "Why this exists") | FIXED | `grep -n "Why this exists" README.md` → no hits. The implementer withdrew the override and cut it. |
| MINOR-2 (heading case: `## Use Parley Deck`) | REJECTED by implementer — ruling sustained | See ruling below. I withdrew this in round 02 and do not revive it. |
| MINOR-3 (`#install-update-and-remove` anchor) | NOT FIXED, correctly deferred | The anchor remains. IMPLEMENTATION.md "Still open" records it as NOT TESTED. GitHub-generated anchor, inert text elsewhere. Not a defect. |
| NIT-1 (winget README still says "draft") | NOT FIXED, correctly deferred | Out of scope (separate file). IMPLEMENTATION.md "Still open" records it as a follow-up. |
| NIT-2 (gemini line trailing comment) | NOT FIXED | The trailing `# legacy Gemini only` comment is still at README:191. Usability note, not a truth defect. Acceptable. |

All round-01 findings are either FIXED or correctly deferred/rejected with
stated reasons. No finding was silently dropped.

## Round-02 findings — disposition after fix-up cycle 4

| Round-02 finding | Status | Verification |
|---|---|---|
| MAJOR-1 (IMPLEMENTATION.md mischaracterises my round-01 MAJOR-1) | FIXED | The implementer withdrew the characterisation in cycle 2. See the withdrawal ruling below. |
| MINOR-1 (worktree sentence wrap defect, 108 chars) | FIXED | The worktree paragraph was rewritten in cycle 1 and rewrapped in cycle 4. The longest prose line in the worktree section is now 90 chars (README:102). No prose line in the entire file exceeds 90 chars. |
| MINOR-2 (`--yes` flag listed but never shown in action) | FIXED | `--yes` is now annotated "(sync-project only: without it, sync-project is a dry run)" at README:181, moved out of the shared flag list. Verified in source: `lib/installer.js:403,407` gate only sync-project on `options.yes`. |
| NIT-1 (`## Use Parley Deck` imperative verb phrase) | NOT FIXED, style observation | No action required. Still the only imperative heading. Not a defect. |

## Round-03 findings — disposition after fix-up cycle 4

| Round-03 finding | Status | Verification |
|---|---|---|
| MINOR-1 (worktree sentence 132 chars, wrap convention) | FIXED | Cycle 4 rewrapped the paragraph. The "Use it when two or more sessions..." clause is now split across README:102-103, max line length 85 chars. No prose line in the file exceeds 90 chars. |
| NIT-1 (`--yes` in general flag block, only affects sync-project) | FIXED | Annotated and moved out of the shared flag list (README:181). |

## Ruling on the implementer's rejection of MINOR-2 (round 01)

Sustained. I withdrew this in round 02 and do not revive it. Sentence case
capitalises the first word and proper nouns; "Use" is the first word and is
correctly capitalised, and "Parley Deck" is a proper noun. `## Use Parley Deck`
is sentence case. The heading is consistent with the D-2 renames.

## Ruling on the implementer's withdrawal of the MAJOR-1 mischaracterisation

The withdrawal is complete and honest. I verified the exact mechanism the
implementer describes in IMPLEMENTATION.md "On the provenance finding":

- Case-sensitive `grep -c "Fusion" references/COOPERATION.md` → 0
- Case-sensitive `grep -c "ExecPlans" references/COOPERATION.md` → 0
- Case-sensitive `grep -c "Preflight readiness" references/COOPERATION.md` → 0
- Case-insensitive `grep -in "fusion" references/COOPERATION.md` → :1108, the
  idea slug `meta-protocol-change-fusion-execplans` (lowercase, hyphenated)
- Case-insensitive `grep -in "preflight" references/COOPERATION.md` → :751,
  the command `parley preflight`

The implementer's corrected statement — "hermes-1's round-01 finding was
accurate exactly as filed" — is correct. The self-favourable direction was
acknowledged ("wrong, and wrong in a self-favourable direction"). The fix
(keep only RHO, delete the rest) was already correct in cycle 1; only the
justification was wrong, and it is now corrected. The withdrawal is complete,
honest, and accurately describes the mechanism. No further action needed.

## Ruling on DEVIATION D-1 (keep or cut "Why this exists")

Cut. Both reviewers ruled against keeping it. The implementer withdrew the
override and cut it in cycle 1. The section is absent from the shipped file.
Resolved.

## C9 truth table — row-by-row re-verification

Every row checked against the shipped README (261 lines), not against
IMPLEMENTATION.md's claim:

| Row | Fixed in file? | Evidence |
|---|---|---|
| `:9` runtime list | YES | README:147-150 says "fourteen named runtimes" + lists them + "plus `generic`, a destination you point at with `--dest`." Counted 14 in `lib/installer.js:13-113`. "15 runtimes" absent. |
| `:21-23` "append-only" | YES | `grep -c "append-only" README.md` → 0. |
| `:26-27` "rates confidence by agreement" | YES | `grep -ci "rates confidence" README.md` → 0. |
| `:119` "any capable tier-1 model" | YES | `grep -c "tier-1" README.md` → 0. Replaced with "plain Markdown by design" (README:139). |
| `:148-176` Repository Layout | YES | Rewritten tree (README:226-243) includes `addons/` with four skills, `test/`, `packaging/`, `scripts/`, `NOTICE.md`, `RELEASING.md`. All verified against actual directory listing. |
| `:239` `v1.2.1` | YES | `grep -c "v1.2.1" README.md` → 0. Windows line (189) is versionless. |
| `:242` WinGet "until accepted" | YES | Replaced with `winget install Feci.ParleyDeckSkill` (README:189). |
| `:371` "all discovered installed CLI agents" | YES | Now "a bounded participant set — normally two to four, including at least one non-facilitator when one is available" (README:210-211). Matches `SKILL.md:283`. |
| `:397` "value should be obvious" | YES | `grep -c "should be obvious" README.md` → 0. Section deleted. |
| Eight prompt blocks → three | YES | Three prompt blocks at README:119, 124, 129. |

All ten rows fixed in the file, not just claimed in IMPLEMENTATION.md.

## Factual sentence audit (CHECK 1)

Every factual sentence in the new README was checked against a shipped file or
live command. No uncheckable factual sentence found.

- "The installer places five skills into each detected runtime" (README:19) —
  verified: installer copies 5 skill payloads; `--no-addons` and `--only` flags
  confirmed in `--help`.
- "fourteen named runtimes" enumeration (README:147-150) — verified: 14
  TARGETS entries in `lib/installer.js:13-113`; `generic` requires `--dest`
  per `--help` and `lib/installer.js:258-259`.
- "exits `3`" / "exits `4`, not `0`" (README:67, 68-69) — verified at
  `addons/parley-design-check/SKILL.md:60` (exit 3) and
  `lib/engine.js:2069` (exit 4).
- "carries no fallback registry" (README:64) — verified at
  `addons/parley-design-check/SKILL.md:173`.
- "capability declaration is generated from its detector modules" (README:69-70)
  — verified at `addons/parley-design-check/SKILL.md:187`.
- "live create/update requires an opt-in connector" (README:85-86) — verified
  at `addons/parley-tracker/SKILL.md:381-388`.
- "Sync is one-way by default, and pull reconciliation may write back only
  fields declared `mirror-owned`" (README:82-84) — verified at
  `addons/parley-tracker/SKILL.md:29,135`.
- "neutral projections for Jira, Linear, GitHub Issues, GitLab, Trello, and
  plain boards" (README:84-85) — verified at
  `addons/parley-tracker/SKILL.md:3,14,23-24`.
- "The worktree-allocation table in `IMPLEMENTATION.md` is the lock manifest"
  (README:97) — verified at `addons/parley-worktrees/SKILL.md:340-341`.
- "an intersection is refused unless an explicit override is recorded"
  (README:99-100) — verified at
  `addons/parley-worktrees/SKILL.md:55,325-328,393-395`.
- "Git gives that worktree its own working tree, index, `HEAD` and branch — it
  does not give it its own ports, databases or caches" (README:100-102) —
  verified at `addons/parley-worktrees/SKILL.md:62-63,72` (git gives working
  tree, index, HEAD, branch) and :45,76 (env/ports/databases/caches are
  per-worktree only if inside the working tree, not git-provided).
- "RHO (Retrospective Harness Optimization), credited in
  `references/COOPERATION.md` §13" (README:252-253) — verified: §13 is at
  COOPERATION.md:1076, RHO credited at :1078.
- "RHO's single-model self-preference is deliberately replaced by the deck's
  multi-agent quorum" (README:253-254) — verified: COOPERATION.md:1078.
- "`NOTICE.md` records `hallmark` and `impeccable` as the prior art studied for
  the design add-ons" (README:254-255) — verified at NOTICE.md:13,17.
- "manifest.yaml (neutral), openai.yaml (Codex UI metadata)" (README:237) —
  verified: `manifest.yaml:3` has `vendor_neutral: true`; `openai.yaml:2` has
  `display_name`.
- "A marker-only directory created by this installer is not treated as a real
  runtime" (README:144-145) — verified at `lib/installer.js:716-718`.
- Commands list (README:175): `install`, `paths`, `doctor`, `status`,
  `sync-project`, `uninstall`, `--version` — verified against `--help`: all 7
  present, exact match.
- Flag block (README:178-182): all flags verified against `--help` output. The
  `--target` value list matches exactly, character for character.
- "standalone binaries also on GitHub releases" (README:189) — verified at
  RELEASING.md:47.
- Bounded participant set "normally two to four, including at least one
  non-facilitator" (README:210-211) — verified at `SKILL.md:283`.
- "30-minute timeout" (README:212) — verified at `SKILL.md:288`.
- "The Gemini CLI command and `--target gemini` are two managers for the *same*
  destination — `~/.gemini/extensions/parley-deck`" (README:194-195) — verified:
  `lib/installer.js:40` resolves gemini to `.gemini/extensions`. Same
  destination, two managers. Correct.
- parley-design "no numeric aesthetic score, no house look, and no 'good
  default aesthetic' guessed from the category" (README:53-54) — verified at
  `addons/parley-design/SKILL.md:38,40-41`.
- "zero to three bounded grafts may come from losing directions, but none may
  modify the winner's token file" (README:55-56) — verified at
  `addons/parley-design/references/PDS.md:198,331-332,383-385`.
- "It is markdown doctrine with no runtime, network, or framework" (README:49-50)
  — verified at `addons/parley-design/SKILL.md:11`.
- "Acceptance criteria carry audience tags, and the Definition of Done points
  back to those criteria with verification commands" (README:78-80) — verified
  at `addons/parley-tracker/SKILL.md:174,166,267-274`.

## Documented command results (CHECK 2)

| Documented surface | Result | Evidence |
|---|---|---|
| `winget install Feci.ParleyDeckSkill` | NOT TESTED | Windows-only; PackageIdentifier confirmed in local manifest and C10's external registry evidence |
| `brew install feci/parley/parley-deck-skill` | NOT TESTED | Homebrew not available; tap repo confirmed (HTTP 200), formula confirmed at `Formula/parley-deck-skill.rb` (v1.5.0) |
| `gemini extensions install …` | NOT TESTED | Requires Gemini CLI; extension mode confirmed in `lib/installer.js:38-43` and `gemini-extension.json` |
| `npx -y parley-deck-skill@latest install --target all` | NOT TESTED (npx) / PASS (local binary) | npm package confirmed live at 1.5.0; local binary `install --target all --dry-run` exits 0 |
| `parley-deck-skill paths` | PASS | Exit 0, 6 detected targets. README correctly says "every *detected* target". |
| `parley-deck-skill paths --target all --include-undetected` | PASS | Exit 0, 14 unique targets. |
| Project-scope install (old command, without `--include-undetected`) | FAIL (correctly removed) | Exit 1, "No installed agent runtimes were detected." |
| Project-scope install (new command, with `--include-undetected`) | PASS | 70 skill actions, exit 0. |
| `agy plugin validate ~/.gemini/config/plugins/parley-deck` | NOT TESTED | Requires `agy` CLI; path matches `lib/installer.js:33` skillDir. |

## Rewritten provenance section and worktree sentence — are they now true?

Yes, both.

The provenance section (README:248-257) stays within the manual-facilitation
boundary, names only the RHO lineage that COOPERATION.md:1076-1078 records, and
the hallmark/impeccable design prior art that NOTICE.md:13,17 records. No
uncheckable attribution remains. "Deterministic automated orchestration is not
part of it and requires separate tooling" is accurate — no shipped file
promises automated orchestration.

The worktree sentence (README:100-102) accurately distinguishes Git's isolation
(working tree, index, HEAD, branch) from runtime resources (ports, databases,
caches) that require manifest overrides. Verified against
`addons/parley-worktrees/SKILL.md:62-63,72,45,201,319`. The wrap defect from
rounds 02-04 is fixed: no prose line exceeds 90 characters.

## Machine-made reading check (CHECK 5)

The README does not read as generated. The five catalogue entries retain their
distinct voices. The hook is direct. The install section is terse and
operational. The provenance section is specific and evidence-anchored.

No sentence reads as machine-made slop. The cycle-3 sentence that codex-1
flagged as the one machine-made line ("The last two lines depend on those CLIs
rather than on anything this package ships") is gone, replaced by the
accurate Gemini shared-destination explanation. The worktree paragraph is now
cleanly wrapped with no stacked-clause run-on. The file reads as written by
someone who knows the codebase.

## Line count and links (CHECK 6)

- `wc -l README.md` → 261. Ceiling 300. Pass with 39 lines of margin.
- Five catalogue links (`./SKILL.md`, `./addons/parley-design/SKILL.md`, etc.)
  — all resolve to real files (verified with `[ -f ]`).
- One internal anchor link:
  `[Install, update, and remove](#install-update-and-remove)` at README:113 →
  heading `## Install, update, and remove` at README:141 → GitHub anchor
  `#install-update-and-remove`. Correct on GitHub; inert text elsewhere
  (correctly flagged as NOT TESTED in IMPLEMENTATION.md).
- No dead internal links. No broken anchors detected.
- Prose line sweep: no prose line exceeds 90 characters. The 11 lines over 100
  chars are all HTML comments, link-bearing headings, fenced prompt text, and
  command/flag lines — none wrap.

## What the fix-up broke

Nothing. Cycle 4 rewrapped the worktree paragraph and scoped `--yes` to
sync-project. I checked every section that was touched and every section
adjacent to a touch. No new factual error, no broken link, no orphaned
reference was introduced. The file is internally consistent and every factual
sentence checks out.

## Findings

### CRITICAL

None.

### MAJOR

None.

### MINOR

None.

### NIT

None.

I have no remaining findings. Every factual sentence in the README is checkable
against a shipped file or live command. Every C9 row is fixed in the file. The
round-01 MAJORs are fixed, the round-01 MINORs are fixed or correctly
deferred/rejected, the round-02/03 findings are fixed, the provenance
withdrawal is complete and honest, the project-scope install fix works (verified
by running both the old and new commands), and the README does not read as
machine-made. The file is 261 lines against a 300 ceiling with no dead links or
broken anchors.

### Signoff: hermes-1 — 2026-07-29
Status: ✅ ACCEPT
