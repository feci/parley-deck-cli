---
idea: readme-skill-catalogue
review-round: 01
agent: hermes-1
date: 2026-07-29
---

Adversarial review of claude-1's README rewrite (commit 0061dc2, 401→278 lines).
Every claim below was checked against a shipped file or a live command, not eyeballed.

## Verification log (commands actually run)

- `wc -l README.md` → 278 (ceiling 300, pass)
- `node bin/parley-deck-skill.js --help` → all 7 commands and all 10 flags match the README's documented set
- `node bin/parley-deck-skill.js paths --target all` → works, lists install dirs per target
- `node bin/parley-deck-skill.js install --target all --dry-run` → works, exit 0
- `node bin/parley-deck-skill.js --version` → 1.5.0 (matches package.json)
- `npm view parley-deck-skill@latest version` → 1.5.0 (npm registry live)
- `npm test` → 247 pass, 0 fail (matches IMPLEMENTATION.md)
- `curl -sI https://github.com/feci/homebrew-parley` → HTTP 200 (tap exists)
- Banned-string sweep (`15 runtimes`, `tier-1`, `should be obvious`, `append-only`, `rates confidence`, `v1.2.1`, `Until the WinGet`) → all 0 occurrences
- TARGETS array counted in `lib/installer.js:13-113` → 14 entries (codex, claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe, cursor, opencode, aionrs)
- Exit code 4 verified at `addons/parley-design-check/lib/engine.js:2069`
- Exit code 3 verified at `addons/parley-design-check/SKILL.md:60`
- Tracker connector boundary verified at `addons/parley-tracker/SKILL.md:381-388`
- Bounded participant set verified at `SKILL.md:283` and `references/COOPERATION.md:20`
- Worktree lock manifest verified at `addons/parley-worktrees/SKILL.md:340-341`
- Five catalogue link targets all resolve to real files

Commands NOT TESTED (cannot run on this host):
- `winget install Feci.ParleyDeckSkill` — Windows-only; PackageIdentifier confirmed in local manifest and C10's external registry evidence (versions 1.0.4–1.4.6)
- `brew install feci/parley/parley-deck-skill` — Homebrew not available on this macOS in this context; tap repo `feci/homebrew-parley` confirmed to exist (HTTP 200), and RELEASING.md documents the tap as the formula location
- `gemini extensions install https://github.com/feci/parley-deck-skill` — requires Gemini CLI; the extension mode is confirmed in `lib/installer.js:38-43` and `gemini-extension.json`
- `npx -y parley-deck-skill@latest install --target all` — run via local binary instead; npm package confirmed live at 1.5.0

---

## CRITICAL

None.

## MAJOR

### MAJOR-1 — Protocol prior-art attributions are uncheckable against any shipped file

README lines 262-265 list five protocol-level prior-art attributions:

> **OpenRouter Fusion** → the compare-not-merge consensus lens. **OpenAI ExecPlans / PLANS.md** → a resume-from-the-document `FINAL.md` and a living `IMPLEMENTATION.md`. **RHO** → advisory, quorum-gated retrospective optimization. **kindly** → strict gates, stopping judgment, no-suppression dispositions, artifact-wins. **Preflight readiness** → protocol-freshness and roster liveness before each idea.

Of these five, only **RHO** appears in a shipped file (`references/COOPERATION.md:1078`). "OpenRouter Fusion", "OpenAI ExecPlans / PLANS.md", "kindly", and "Preflight readiness" appear **nowhere in the repo** outside this README. I ran `grep -r` across all `.md`, `.yaml`, and `.json` files (excluding `node_modules`): zero hits for any of the four. `NOTICE.md` records prior art only for the design add-ons (`parley-design` and `parley-design-check`), not for the protocol itself.

The README's own sentence at line 266 says "NOTICE.md records the prior art studied for the design add-ons" — which is true. But the five protocol attributions above it are not in NOTICE.md or anywhere else. The idea's entire purpose (stated in FINAL.md line 14: "every claim in C9's truth table fixed" and the Definition of Done: "no invented number anywhere") is that factual sentences are checkable. These four attributions are not checkable.

These claims were carried forward from the pre-rewrite README (commit 66a94ff, line 35-39), so the implementer did not invent them. But the C9 truth audit did not cover this section, and the rewrite preserved uncheckable claims in a README whose raison d'être is checkability.

Fix: either add the four missing attributions to NOTICE.md (or a new "Protocol prior art" section there), or remove the four uncheckable lines from the README and keep only the RHO reference that COOPERATION.md already supports. The RHO attribution itself is safe — it traces to COOPERATION.md:1078.

### MAJOR-2 — `## Status` section is not in the FINAL.md binding section order and was not flagged as a deviation

FINAL.md C1 defines a binding 10-section order. The README ships 12 sections. Two are not in the order:

1. `## Why this exists` (line 140) — self-flagged as D-1 in IMPLEMENTATION.md.
2. `## Status` (line 270) — **not flagged anywhere in IMPLEMENTATION.md**.

The `## Status` section sits between "Related repositories" (C1 item 9) and "License" (C1 item 10). It was carried forward from the old README (line 393 in commit 66a94ff). The implementer's deviation list covers D-1 (Why this exists), D-2 (section renames), and D-3 (anchor), but is silent on Status.

The section's factual content is accurate (the `ideas/<slug>/` file list matches `SKILL.md:90`), so this is not a truth defect. It is a process defect: a binding section order was silently extended by one section the implementer chose not to report.

Fix: either cut `## Status` (its content can fold into the hook or the related-repositories section), or add it to the deviation list in IMPLEMENTATION.md with a justification.

## MINOR

### MINOR-1 — D-1 override ruling: "Why this exists" should be cut

The implementer kept "Why this exists" (11 lines, 142-149) despite codex-1's round-02 cut list deleting it, the implementer's own round-01 marking it for deletion, and FINAL.md's section order omitting it. The justification in IMPLEMENTATION.md D-1 is: "after cutting to 278 lines the budget was not the binding constraint, and the failure-mode list is the argument the hook only gestures at."

Ruling: **cut it.** The hook already states the failure mode ("One model playing four reviewers is still one model") and the proof ("files in your repository that you can read, diff, and resume"). The "Why this exists" section rephrases the same argument in a less compressed form — its first sentence ("Multi-agent workflows fail in predictable ways: one agent anchors the rest before they form their own view…") is the hook expanded, and its second paragraph ("Parley Deck turns the conversation into project artifacts…") is the catalogue's lead-in restated. The implementer overrode two participants' stated preferences on a section that duplicates content already present in the binding copy. The override is not justified: the budget argument is irrelevant when the binding section order omitted it, and "the hook only gestures at" the failure mode is a reading the hook itself contradicts (it names the failure mode in its first sentence).

### MINOR-2 — `## Use Parley Deck` heading uses Title Case, inconsistent with the file's sentence-case renames

The implementer's D-2 renamed headings to sentence case (`## Local agent contract`, `## Repository layout`, `## Related repositories, and what this one owes`). But `## Use Parley Deck` (line 114) retains Title Case — "Parley Deck" is a proper noun, so that part is fine, but "Use" is not. The section is `## Use Parley Deck` in both FINAL.md and the shipped file. This is not a defect against FINAL.md (which used the same casing), but it is internally inconsistent with the D-2 sentence-case convention the implementer applied elsewhere. If D-2 was worth doing, "Use" should have been lowercased to "use" — or D-2 should not have been done at all.

### MINOR-3 — The `#install-update-and-remove` anchor was not verified on GitHub

D-3 notes the anchor is GitHub-generated and "not verified on npmjs.com." The heading `## Install, update, and remove` generates the GitHub anchor `#install-update-and-remove` (lowercase, punctuation stripped, spaces → hyphens), which matches the link at line 112. This is correct on GitHub. The implementer's flag is honest. No fix needed; recorded for completeness.

## NIT

### NIT-1 — `packaging/winget/README.md` still says "draft manifest"

IMPLEMENTATION.md line 75-76 notes this as a deliberate follow-up, not a defect in this change. The winget README says "This directory contains a draft manifest for submitting…" while the package is published (C10). This is a pre-existing staleness in a different file, correctly out of scope. No action required for this idea.

### NIT-2 — README says "Other channels:" with a code block mixing four install methods

Line 194-201 presents `brew`, `winget`, `npm install -g`, and `gemini extensions install` in a single code block under "Other channels:". The `brew` line chains `&& parley-deck-skill install --target all`, the `winget` line has a comment, and the `npm` line chains `&& parley-deck-skill install`. This is readable but the `gemini` line has a trailing comment (`# legacy Gemini only`) that could be mistaken for part of the command if copy-pasted without reading. Not a truth defect; a minor usability note.

---

## C9 truth table — row-by-row verification

Every row checked against the shipped README, not against IMPLEMENTATION.md's claim:

| Row | Fixed in file? | Evidence |
|---|---|---|
| `:9` runtime list | YES | README:157-160 says "fourteen named runtimes" + lists them + "plus `generic`, a destination you point at with `--dest`." Counted 14 in `lib/installer.js:13-113`. "15 runtimes" absent. |
| `:21-23` "append-only" | YES | No occurrence of "append-only" anywhere in README. `grep -c "append-only" README.md` → 0. |
| `:26-27` "rates confidence by agreement" | YES | No occurrence. `grep -ci "rates confidence" README.md` → 0. |
| `:119` "any capable tier-1 model" | YES | No occurrence of "tier-1". Replaced with "plain Markdown by design" (line 138). |
| `:148-176` Repository Layout | YES | Rewritten tree (lines 234-252) includes `addons/` with four skills, `test/`, `packaging/`, `scripts/`, `NOTICE.md`, `RELEASING.md`. All verified against actual directory listing. |
| `:239` `v1.2.1` | YES | No occurrence of `v1.2.1`. Windows line (198) is versionless. |
| `:242` WinGet "until accepted" | YES | Replaced with `winget install Feci.ParleyDeckSkill` (line 198). PackageIdentifier confirmed in local manifest. C10 external registry evidence sustained. |
| `:371` "all discovered installed CLI agents" | YES | Now "a bounded participant set — normally two to four, including at least one non-facilitator when one is available" (line 219). Matches `SKILL.md:283`. |
| `:397` "value should be obvious" | YES | No occurrence. Section deleted. |
| Eight prompt blocks → three | YES | Three prompt blocks at lines 118, 123, 128. The other two `text` blocks (187, 234) are the flags listing and the directory tree, not prompt blocks. |

All ten rows fixed in the file, not just claimed in IMPLEMENTATION.md.

## D-1 ruling (explicit)

**Cut "Why this exists."** The implementer's override of two participants' stated preferences is not justified. The section duplicates the hook's failure-mode statement and the catalogue's lead-in. FINAL.md's binding section order omitted it, and the implementer's own round-01 marked it for deletion. The line-budget argument is irrelevant: the binding constraint was the section order, not the line count. See MINOR-1 above.

## Machine-made reading check

The README does not read as generated. The five catalogue entries have distinct voices (codex-1's precision on artifact ownership, kimi-1's compression into refusals, the "no stack trace" framing). The hook is direct. The install section is terse and operational.

One sentence that reads slightly as generated: line 153, "The installer uses an AionUI-style local runtime registry: it checks known user-level agent directories and CLI commands, then installs into the runtimes it detects." The colon-introduced explanation pattern ("X: it does Y") appears three times in the install section (lines 153, 157, 209) and is the closest thing to a machine tell. It is not slop — it is compressed technical prose — but the repeated colon-explanation rhythm is a mild structural sameness. Not a finding against the anti-slop skill; noted for awareness.

## Line count and links

- `wc -l README.md` → 278. Ceiling 300. Pass with 22 lines of margin.
- Five catalogue links (`./SKILL.md`, `./addons/parley-design/SKILL.md`, etc.) — all resolve to real files.
- One internal anchor link: `[Install, update, and remove](#install-update-and-remove)` at line 112 → heading `## Install, update, and remove` at line 151 → GitHub anchor `#install-update-and-remove`. Correct.
- No dead internal links. No broken anchors detected.

### Signoff: hermes-1 — 2026-07-29
Status: 🟡 ACCEPT-WITH-RESERVATIONS
