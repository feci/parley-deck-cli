---
idea: readme-skill-catalogue
review-round: 03
agent: hermes-1
date: 2026-07-29
---

Adversarial re-review of claude-1's README rewrite after fix-up cycle 2 (commit
5df9d04, 260 lines). Every claim below was checked against a shipped file or a
live command, not eyeballed. I tried to break the claim that this is done.

## Verification log (commands actually run)

- `wc -l README.md` → 260 (ceiling 300, pass, 40 lines of margin)
- `node bin/parley-deck-skill.js --help` → all 7 commands and all flags match
  the README's documented set; `--target` value list matches character for character
- `node bin/parley-deck-skill.js --version` → 1.5.0 (matches package.json)
- `npm view parley-deck-skill@latest version` → 1.5.0 (npm registry live)
- `npm test` → 247 pass, 0 fail (matches IMPLEMENTATION.md)
- `node bin/parley-deck-skill.js install --target all --dry-run --yes` → exit 0
- `node bin/parley-deck-skill.js paths` → 6 detected targets (codex, claude, agy,
  gemini, hermes, kimi)
- `node bin/parley-deck-skill.js paths --target all --include-undetected` → 14
  targets (all named runtimes)
- `node bin/parley-deck-skill.js install --scope project --target all --project
  /tmp/pd-test-proj --dry-run` → exit 1, "No installed agent runtimes were
  detected" (the OLD command, reproducing codex-1's round-02 MAJOR)
- `node bin/parley-deck-skill.js install --scope project --target all --project
  /tmp/pd-test-proj --include-undetected --dry-run` → 71 planned install lines,
  exit 0 (the NEW command, verifying the cycle-2 fix)
- `curl -sI https://github.com/feci/homebrew-parley` → HTTP 200 (tap exists)
- `curl -s https://raw.githubusercontent.com/feci/homebrew-parley/main/Formula/parley-deck-skill.rb`
  → formula exists, version 1.5.0
- Banned-string sweep (`15 runtimes`, `tier-1`, `should be obvious`,
  `append-only`, `rates confidence`, `v1.2.1`, `AionUI`, `Installing this
  package`) → all 0 occurrences
- `grep -c "Fusion" references/COOPERATION.md` → 0
- `grep -c "ExecPlans" references/COOPERATION.md` → 0
- `grep -c "Preflight readiness" references/COOPERATION.md` → 0
- Case-insensitive: `grep -in "fusion" COOPERATION.md` → :1108 (idea slug
  `meta-protocol-change-fusion-execplans`); `grep -in "preflight" COOPERATION.md`
  → :751 (`parley preflight` command)
- TARGETS array counted in `lib/installer.js:13-113` → 14 entries (codex, claude,
  agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe, cursor,
  opencode, aionrs)
- Exit code 3 verified at `addons/parley-design-check/SKILL.md:60`
- Exit code 4 verified at `addons/parley-design-check/lib/engine.js:2069`
- Tracker connector boundary verified at `addons/parley-tracker/SKILL.md:381-388`
- `mirror-owned` fields verified at `addons/parley-tracker/SKILL.md:27,29,135`
- Tracker projections (Jira, Linear, GitHub Issues, GitLab, Trello, plain boards)
  verified at `addons/parley-tracker/SKILL.md:3,14,23-24`
- Worktree git isolation claim verified at
  `addons/parley-worktrees/SKILL.md:62-63,72`
- Worktree env/port/database/cache overrides verified at
  `addons/parley-worktrees/SKILL.md:45,200-201,319,484-485`
- RHO lineage verified at `references/COOPERATION.md:1076-1078` (§13)
- `hallmark` and `impeccable` prior art verified at `NOTICE.md:13,17`
- Bounded participant set verified at `SKILL.md:283` ("normally 2-4 active
  participants … at least one non-facilitator")
- 30-minute timeout verified at `SKILL.md:288` ("timeout: 30 minutes per agent
  process")
- `manifest.yaml` vendor_neutral and `openai.yaml` display_name verified at
  `agents/manifest.yaml:3-4`, `agents/openai.yaml:2`
- Five catalogue link targets all resolve to real files
- Repository layout tree verified against actual directory listing
- `parley-design` refusals (no score, no house look, no "good default aesthetic",
  bounded grafts 0-3, winner's token file not modified) verified at
  `addons/parley-design/SKILL.md:38,40-41` and `references/PDS.md:198,331-332,383-385`
- `parley-design` "no runtime, no network, no framework" verified at
  `addons/parley-design/SKILL.md:11`
- `parley-design-check` capability declaration generated from detectors verified
  at `addons/parley-design-check/SKILL.md:187-189`
- `parley-design-check` "no fallback registry" verified: `SKILL.md:173` ("This
  add-on carries no registry of its own, and a test enforces that")
- `parley-design-check` `rule-id — violation — remedy` format verified at
  `addons/parley-design-check/SKILL.md:12`
- Tracker section names (`At a glance`, `[B] Business`, `[T] Technical`,
  `[A] Agent directives`) verified at `addons/parley-tracker/SKILL.md:99,104-106`
- Marker-only directory exclusion verified at `lib/installer.js:716-718`
- Standalone Windows binaries on GitHub releases verified at
  `RELEASING.md:47,53-60`
- Homebrew formula location verified at `RELEASING.md:91,106` and live at
  `feci/homebrew-parley/Formula/parley-deck-skill.rb`
- `plugin.json` (Antigravity plugin mode) and `gemini-extension.json` (legacy
  Gemini extension mode) both exist as shipped files
- IMPLEMENTATION.md frontmatter `commit: 5df9d04` matches actual cycle-2 head

Commands NOT TESTED (cannot run on this host):
- `winget install Feci.ParleyDeckSkill` — Windows-only; PackageIdentifier
  confirmed in local manifest and C10's external registry evidence (versions
  1.0.4–1.4.6)
- `brew install feci/parley/parley-deck-skill` — Homebrew not available in this
  context; tap repo `feci/homebrew-parley` confirmed to exist (HTTP 200), formula
  confirmed at `Formula/parley-deck-skill.rb` (v1.5.0), RELEASING.md:88-108
  documents the tap as the formula location
- `gemini extensions install https://github.com/feci/parley-deck-skill` — requires
  Gemini CLI; the extension mode is confirmed in `lib/installer.js:38-43` and
  `gemini-extension.json`
- `npx -y parley-deck-skill@latest install --target all` — run via local binary
  instead; npm package confirmed live at 1.5.0

## Round-01 findings — disposition after fix-up cycles 1 and 2

| Round-01 finding | Status | Verification |
|---|---|---|
| MAJOR-1 (provenance attributions uncheckable) | FIXED | The five-attribution block is gone. README:251-256 now states only the RHO lineage that COOPERATION.md:1078 actually records, plus the hallmark/impeccable prior art that NOTICE.md:13,17 records. No occurrence of "Fusion", "ExecPlans", "kindly", or "Preflight readiness" as attribution labels. Verified with case-sensitive grep: 0 hits for each in COOPERATION.md. |
| MAJOR-2 (`## Status` unreported deviation) | FIXED | `grep -n "## Status" README.md` → no hits. The file now has exactly the C1 section order (with the D-2 renames). |
| MINOR-1 (D-1: cut "Why this exists") | FIXED | `grep -n "Why this exists" README.md` → no hits. The implementer withdrew the override (IMPLEMENTATION.md "On D-1, plainly") and cut it. |
| MINOR-2 (heading case: `## Use Parley Deck`) | REJECTED by implementer — ruling sustained | See ruling below. I withdrew this in round 02 and do not revive it. |
| MINOR-3 (`#install-update-and-remove` anchor) | NOT FIXED, correctly deferred | The anchor remains. IMPLEMENTATION.md "Still open" records it as NOT TESTED. It is a GitHub-generated anchor that renders as inert text elsewhere. Not a defect. |
| NIT-1 (winget README still says "draft") | NOT FIXED, correctly deferred | Out of scope (separate file). IMPLEMENTATION.md "Still open" records it as a follow-up. |
| NIT-2 (gemini line trailing comment) | NOT FIXED | The trailing `# legacy Gemini only` comment is still at README:190. This was a usability note, not a truth defect. Acceptable to leave. |

## Round-02 findings — disposition after fix-up cycle 2

| Round-02 finding | Status | Verification |
|---|---|---|
| MAJOR-1 (IMPLEMENTATION.md mischaracterises my round-01 MAJOR-1) | FIXED | The implementer withdrew the characterisation in cycle 2. IMPLEMENTATION.md "On the provenance finding" now states: "That was wrong, and wrong in a self-favourable direction. My grep was case-insensitive (`grep -ril`)…" and "hermes-1's round-01 finding was accurate exactly as filed." I verified the withdrawal is complete and honest: case-sensitive grep for "Fusion", "ExecPlans", and "Preflight readiness" in COOPERATION.md returns 0; case-insensitive grep matches the idea slug and the `parley preflight` command exactly as the implementer describes. The withdrawal is accurate. |
| MINOR-1 (worktree sentence wrap defect) | PARTIALLY FIXED | See MINOR-1 below. |
| MINOR-2 (`--yes` flag listed but never shown in action) | NOT FIXED, correctly left | The `--yes` flag is still listed in the flag block (README:180) and no command example uses it. This was a documentation-gap observation, not a truth defect. codex-1 correctly noted in round 02 that a flag inventory promises existence, not an example. Acceptable. |
| NIT-1 (`## Use Parley Deck` is an imperative verb phrase) | NOT FIXED, style observation | No action required. Still the only imperative heading. Not a defect. |

## codex-1 round-02 findings — cross-check

| codex-1 round-02 finding | Status | Verification |
|---|---|---|
| MAJOR: project-scope install command fails | FIXED | README:165 now includes `--include-undetected`. I ran both: the old command (`install --scope project --target all --project /tmp/pd-test-proj --dry-run`) exits 1 with "No installed agent runtimes were detected"; the new command (`+ --include-undetected --dry-run`) produces 71 planned install lines, exit 0. Fix verified. |
| MAJOR: provenance cross-repo claim unsupported | FIXED | README:249 now says "Deterministic automated orchestration is not part of it and requires separate tooling." No named external repository. The "The one lineage" phrasing was changed to "The protocol lineage recorded here" (README:251). |
| MINOR: `paths` does not report every target | FIXED | README:198 now says "every *detected* target" and documents `paths --target all --include-undetected` for all fourteen (README:199). I verified: default `paths` returns 6 detected targets; `--target all --include-undetected` returns 14. |
| MINOR: runtime-channel instructions uncheckable | FIXED | The restart line is now conditional ("if it does not pick the change up", README:112). The Gemini/Antigravity lines are qualified ("depend on those CLIs rather than on anything this package ships", README:193). The `$skill-installer` sentence is deleted (grep returns 0). |
| MINOR: IMPLEMENTATION.md frontmatter stale | FIXED | Frontmatter now says `commit: 5df9d04` with earlier commits listed. Matches actual cycle-2 head. |
| NIT: worktree sentence reads like patch splice | PARTIALLY FIXED | See MINOR-1 below. |

## Ruling on the implementer's rejection of MINOR-2 (round 01)

Sustained. I withdrew this in round 02 and do not revive it. Sentence case
capitalises the first word and proper nouns; "Use" is the first word and is
correctly capitalised, and "Parley Deck" is a proper noun. `## Use Parley Deck`
is sentence case. My round-01 finding conflated "Title Case" with "the first
word is capitalised" — the first word being capitalised is correct in sentence
case. The heading is consistent with the D-2 renames.

## Ruling on the implementer's withdrawal of the MAJOR-1 mischaracterisation

The withdrawal is complete and honest. I verified the exact mechanism the
implementer describes: case-insensitive grep (`grep -il`) matches the lowercase
idea slug `meta-protocol-change-fusion-execplans` at COOPERATION.md:1108 and the
command `parley preflight` at COOPERATION.md:751. Case-sensitive grep for
"Fusion", "ExecPlans", and "Preflight readiness" returns 0 hits in
COOPERATION.md. The implementer's corrected statement — "hermes-1's round-01
finding was accurate exactly as filed" — is correct. The self-favourable
direction was acknowledged. The fix itself (keep only RHO, delete the rest) was
already correct in cycle 1; only the justification was wrong, and it is now
corrected. No further action needed.

## Ruling on DEVIATION D-1 (keep or cut "Why this exists")

Cut. Both reviewers ruled against keeping it. The implementer withdrew the
override and cut it in cycle 1. The section is absent from the shipped file.
Resolved.

## C9 truth table — row-by-row re-verification

Every row checked against the shipped README (260 lines), not against
IMPLEMENTATION.md's claim:

| Row | Fixed in file? | Evidence |
|---|---|---|
| `:9` runtime list | YES | README:147-150 says "fourteen named runtimes" + lists them + "plus `generic`, a destination you point at with `--dest`." Counted 14 in `lib/installer.js:13-113`. "15 runtimes" absent. |
| `:21-23` "append-only" | YES | `grep -c "append-only" README.md` → 0. |
| `:26-27` "rates confidence by agreement" | YES | `grep -ci "rates confidence" README.md` → 0. |
| `:119` "any capable tier-1 model" | YES | `grep -c "tier-1" README.md` → 0. Replaced with "plain Markdown by design" (README:139). |
| `:148-176` Repository Layout | YES | Rewritten tree (README:224-242) includes `addons/` with four skills, `test/`, `packaging/`, `scripts/`, `NOTICE.md`, `RELEASING.md`. All verified against actual directory listing. |
| `:239` `v1.2.1` | YES | `grep -c "v1.2.1" README.md` → 0. Windows line (188) is versionless. |
| `:242` WinGet "until accepted" | YES | Replaced with `winget install Feci.ParleyDeckSkill` (README:188). |
| `:371` "all discovered installed CLI agents" | YES | Now "a bounded participant set — normally two to four, including at least one non-facilitator when one is available" (README:209-210). Matches `SKILL.md:283`. |
| `:397` "value should be obvious" | YES | `grep -c "should be obvious" README.md` → 0. Section deleted. |
| Eight prompt blocks → three | YES | Three prompt blocks at README:119, 124, 129. |

All ten rows fixed in the file, not just claimed in IMPLEMENTATION.md.

## Factual sentence audit (CHECK 1)

Every factual sentence in the new README was checked against a shipped file or
live command. No uncheckable factual sentence found.

- "The installer places five skills into each detected runtime" (README:19) —
  verified: installer copies 5 skill payloads; `--no-addons` and `--only` flags
  confirmed in `--help`.
- "fourteen named runtimes" enumeration (README:147-150) — verified: 14 TARGETS
  entries in `lib/installer.js:13-113`; `generic` requires `--dest` per `--help`.
- "exits `3`" / "exits `4`, not `0`" (README:67, 68-69) — verified at
  `addons/parley-design-check/SKILL.md:60` and `lib/engine.js:2069`.
- "carries no fallback registry" (README:64) — verified at
  `addons/parley-design-check/SKILL.md:173` ("carries no registry of its own,
  and a test enforces that").
- "capability declaration is generated from its detector modules" (README:69-70)
  — verified at `addons/parley-design-check/SKILL.md:187-189` ("generated by
  scanning `lib/detectors`").
- "live create/update requires an opt-in connector" (README:85-86) — verified at
  `addons/parley-tracker/SKILL.md:381-388`.
- "Sync is one-way by default, and pull reconciliation may write back only
  fields declared `mirror-owned`" (README:82-84) — verified at
  `addons/parley-tracker/SKILL.md:27,29,135`.
- "neutral projections for Jira, Linear, GitHub Issues, GitLab, Trello, and
  plain boards" (README:84-85) — verified at
  `addons/parley-tracker/SKILL.md:3,14,23-24`.
- "The worktree-allocation table in `IMPLEMENTATION.md` is the lock manifest"
  (README:97) — verified at `addons/parley-worktrees/SKILL.md:340-341`.
- "an intersection is refused unless an explicit override is recorded"
  (README:99-100) — verified at `addons/parley-worktrees/SKILL.md:55,325-328,393-395`.
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
  verified: `manifest.yaml:4` has `vendor_neutral: true`; `openai.yaml:2` has
  `display_name`.
- "A marker-only directory created by this installer is not treated as a real
  runtime" (README:144-145) — verified at `lib/installer.js:716-718`.
- Commands list (README:175): `install`, `paths`, `doctor`, `status`,
  `sync-project`, `uninstall`, `--version` — verified against `--help`: all 7
  present, exact match.
- Flag block (README:178-181): all flags verified against `--help` output. The
  `--target` value list matches exactly, character for character.
- "standalone binaries also on GitHub releases" (README:188) — verified at
  RELEASING.md:47,53-60.
- Bounded participant set "normally two to four, including at least one
  non-facilitator" (README:209-210) — verified at `SKILL.md:283`.
- "30-minute timeout" (README:211) — verified at `SKILL.md:288`.

## Documented command results (CHECK 2)

| Documented surface | Result | Evidence |
|---|---|---|
| `winget install Feci.ParleyDeckSkill` | NOT TESTED | Windows-only; PackageIdentifier confirmed in local manifest and C10's external registry evidence |
| `brew install feci/parley/parley-deck-skill` | NOT TESTED | Homebrew not available in this context; tap repo confirmed (HTTP 200), formula confirmed at `Formula/parley-deck-skill.rb` (v1.5.0) |
| `gemini extensions install …` | NOT TESTED | Requires Gemini CLI; extension mode confirmed in `lib/installer.js:38-43` and `gemini-extension.json` |
| `npx -y parley-deck-skill@latest install --target all` | NOT TESTED (npx) / PASS (local binary) | npm package confirmed live at 1.5.0; local binary `install --target all --dry-run --yes` exits 0 |
| `parley-deck-skill paths` | PASS | Exit 0, 6 detected targets. README correctly says "every *detected* target". |
| `parley-deck-skill paths --target all --include-undetected` | PASS | Exit 0, 14 targets. |
| Project-scope install (old command, without `--include-undetected`) | FAIL (correctly removed) | Exit 1, "No installed agent runtimes were detected." |
| Project-scope install (new command, with `--include-undetected`) | PASS | 71 planned install lines, exit 0. |
| `agy plugin validate ~/.gemini/config/plugins/parley-deck` | NOT TESTED | Requires `agy` CLI; path matches `lib/installer.js:33` skillDir. |

## CRITICAL

None.

## MAJOR

None.

## MINOR

### MINOR-1 — Worktree sentence line 102 is 132 characters, breaking the file's wrap convention

README:100-103:

    recorded. Each implementer gets a sibling worktree. Git gives that worktree its own
    working tree, index, `HEAD` and branch — it does not give it its own ports, databases or
    caches, so the manifest records those overrides too. Use it when two or more sessions or Phase-5 implementers work in one repository
    at once.

Line 102 ("caches, so the manifest records those overrides too. Use it when two
or more sessions or Phase-5 implementers work in one repository") is 132
characters. Every other prose line in the file wraps at ~80–95 characters. This
is a cosmetic wrap defect, not a truth defect — the content is accurate (verified
against `addons/parley-worktrees/SKILL.md:6,62-63,72`).

This was raised by codex-1 in round 02 (NIT: "the fix-up left the `Use it...`
clause on a 104-character line") and by me in round 02 (MINOR-1: "108
characters"). The cycle-2 IMPLEMENTATION.md claims the NIT was "Rewritten as two
sentences that say what git does and does not give you." The git does/does-not
split was rewritten, but the "Use it when two or more sessions..." clause was
not rewrapped — it is still on one 132-character line. The implementer's claim
that this was fixed is inaccurate: the sentence structure was improved but the
wrap defect persists, and is now worse (132 chars vs the 104-108 chars reported
in round 02).

Fix: rewrap the sentence so no prose line exceeds ~95 characters, matching the
file's convention.

## NIT

### NIT-1 — The `--yes` flag is listed in the general flag block but only affects `sync-project`

README:180 lists `--yes` among the general flags:

    --force  --dry-run  --yes  --json  --include-undetected

The parser accepts `--yes` globally (`lib/installer.js:182-183`), but the flag
only has an effect in `sync-project` (`lib/installer.js:403,407`). For
`install`, `doctor`, `status`, `uninstall`, and `paths`, `--yes` is a silent
no-op. The `--help` usage lines reflect this: only `sync-project` shows `[--yes]`
in its usage. This is not a truth defect (the flag exists), but a reader seeing
it in the general flag block might expect it to auto-confirm installs. This was
my round-02 MINOR-2; codex-1 correctly noted a flag inventory promises
existence, not an example. I record it as a NIT for awareness, not as a finding
requiring action.

## Machine-made reading check (CHECK 5)

The README does not read as generated. The five catalogue entries retain their
distinct voices. The hook is direct. The install section is terse and
operational. The provenance section is now specific and evidence-anchored.

No sentence reads as machine-made slop. The closest candidate is the worktree
sentence at README:100-102, which has a stacked noun list ("ports, databases or
caches") and the "Use it when..." clause — but the content is precise and
factually supported, not slop. The wrap defect (MINOR-1) is a formatting issue,
not a voice issue.

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

Nothing new. The cycle-2 fix-up fixed the project-scope command, the cross-repo
provenance claim, the `paths` scope, and the runtime-channel qualifications. The
worktree sentence rewrite improved the content (the git does/does-not split is
accurate) but did not rewrap the "Use it when..." clause, leaving the 132-char
line (MINOR-1). No new factual error, no broken link, no orphaned reference was
introduced.

### Signoff: hermes-1 — 2026-07-29
Status: ✅ ACCEPT
