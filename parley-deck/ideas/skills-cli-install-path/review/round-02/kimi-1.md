---
idea: skills-cli-install-path
review-round: 02
agent: kimi-1
date: 2026-07-29
---

Re-reviewed branch `readme-skill-catalogue` at `a05bac7` (matches IMPLEMENTATION.md's
head-commit; worktree clean before and after my run — I did not edit the repo). Every check
below was re-run by me on this host (macOS arm64, node v26.5.0, `skills` CLI **1.5.20** — the
same version as round-01 and the consensus measurements), not read out of IMPLEMENTATION.md.
Scratch evidence: `/tmp/kimi1-g7r2.l54F3y` (G7), `/tmp/kimi1-r2-status.txt` (G1 status).
Mutation testing was done in a `git archive` scratch copy, since deleted.

## Round-01 findings — verified disposition

### MAJOR-1 (stale `addons/…` references in shipped skill content) — **FIXED**

- `grep -rn "addons/" skills/` → **0 matches**, reproduced.
- The fix-up commit `a05bac7` touches exactly the six files I flagged, no more.
- The provably broken command now resolves: `node skills/parley-design-check/bin/check.js
  check skills/parley-design-check/SKILL.md` runs from a checkout (rc 4 = `UNJUDGEABLE`,
  which is the documented exit code for a run without a contract — the tool works;
  round-01's `MODULE_NOT_FOUND` is gone).
- The documented test command `node --test "skills/parley-design-check/test/*.test.js"` →
  **159 pass, 0 fail**.
- The "Files" map in `parley-design-check/SKILL.md` and the "live next to this file" map in
  `parley-tracker/SKILL.md` now list `skills/…` paths that exist.

### MINOR-1 (tracker exemplars) — **FIXED**

`templates/epic.md`, `story.md`, `subtask.md` and `bin/validate.test.js` now carry
`skills/parley-tracker/…` in their `files:` frontmatter and `Verify:` commands. Re-ran the
property the exemplars exist to have: `node skills/parley-tracker/bin/validate.js --strict
--dir skills/parley-tracker/templates` → **All 3 ticket(s) passed**, rc 0.

### NIT-1 (panel comparative unmeasured) — **FIXED**

`README.md:110-111` now reads "the coding agents it supports — a longer list than this
package's own installer covers, and theirs to state, not ours". My round-01 ask was
"attributed or dropped"; it is attributed. The comparison no longer stands as our
measurement.

### NIT-2 (G8's recorded reason imprecise) — **FIXED**

IMPLEMENTATION.md's fix-up-cycle-2 section now carries the precise mechanism in substance:
`skills` skips `sourceType === "local"` by design, `skills-lock.json` records `"local"` for a
filesystem-path source, and G8 is therefore only meaningful post-merge against the published
remote. That is the root cause I identified, and the re-test is now aimed at the right thing.

## D-1 withdrawal — verified by running and by mutation

codex-1's reconciliation is implemented and **both values resolve**:

- **Repository consumer:** `gemini-extension.json` carries `contextFileName:
  "skills/parley-deck/SKILL.md"`, and that file exists in the repository (verified).
- **Destination consumer:** a real `node bin/parley-deck-skill.js install --target gemini
  --force` (rc 0) stages `~/.gemini/extensions/parley-deck/gemini-extension.json` with
  `contextFileName: "SKILL.md"`, and `SKILL.md` exists at that destination's root (verified).
- The rewrite (`rewriteStagedGeminiManifest`, `lib/installer.js:755-766`) is narrow —
  parse, set one field, write back — is called only for `target.kind === "gemini"`
  (`:1031-1033`), and runs before `validateInstalledPayload` on the staged tree. The diff
  `f8e3a1c..a05bac7` on `lib/installer.js` is purely additive (new function + call site);
  `validateInstalledPayload` is untouched (0 hits in the diff).

**The two new tests can fail — mutation-tested, not trusted:**

- M1 — reverted the repo manifest to `"SKILL.md"` in a scratch copy → *"the repository
  gemini manifest points at the skill's repository path"* **fails** (46 pass / 1 fail).
- M2 — commented out the rewrite line in `lib/installer.js` → *"a staged gemini install
  rewrites contextFileName to the flat destination shape"* **fails** (46 pass / 1 fail).
- Scratch-copy baseline before mutation: 249 pass / 0 fail, so the failures are attributable
  to the mutations, nothing else.

Also relevant: the staged-install test runs with `packageRoot: root` (the real repository),
so it exercises the real manifest end to end — no fixture divergence between what the test
stages and what ships.

**README honesty on the restored channel: verified.** `README.md:210-219` documents
`gemini extensions install <repo-url>` again, explains the two-consumer split, warns to use
one manager or the other never both, and states plainly: **"we have not been able to run the
Gemini CLI to confirm it end to end."** G11 therefore remains honestly NOT TESTED in the
document itself, not only in IMPLEMENTATION.md. That is the right way to restore a channel
one cannot execute.

## Gate re-runs at a05bac7 (all executed by me)

| Gate | My result |
|---|---|
| **G1** install/status/doctor `--target all` | **PASS, reproduced.** `install --target all --force` rc 0; `status` shows exactly 30 `valid` lines, the only non-valid line being `project: missing protocol` (expected — not a Parley deck); `doctor --target all` rc 0 with **0** non-valid lines. No regression from the fix-ups. |
| **G2** `npm pack --dry-run` | **PASS, reproduced.** 145 `skills/` entries, 0 `addons/` entries, no root `SKILL.md`. |
| **G5** `agy plugin validate` | **PASS, reproduced.** rc 0, `✔ skills: 1 processed`, `✔ agents: 2 processed` — against the destination that now carries the *unrewritten* repository manifest, so the new manifest value does not trouble the Antigravity validator. |
| **G6** `npm test` | **PASS, reproduced.** 249 pass, 0 fail — both in the working tree and in a pristine `git archive HEAD` scratch copy (no `node_modules` needed). The test diff `f8e3a1c..a05bac7` is 31 insertions, 0 deletions: two tests added, **no assertion touched**. |
| **G7** real install, no `--full-depth`, isolated HOME | **PASS, reproduced at a05bac7.** From a `git archive` of HEAD: `Found 5 skills` / `Installed 5 skills` (parley-deck, parley-design, parley-design-check, parley-tracker, parley-worktrees) **with no flag**. Installed core contains exactly `SKILL.md`, `agents/`, `references/` (incl. `COOPERATION.md` and `compatibility.json`, per S4a); `bin`, `lib`, `package.json`, `test`, `dist` all **absent**. S1 and S3 remain fixed and measured. The round-01 scoping caveat stands unchanged: the branch is not pushed, so the true-remote G7 against `feci/parley-deck-skill` is still owed at merge time. |
| **G3 / G8 / G9 / G10 / G11** | **Still NOT TESTED, correctly not claimed as passes.** G3: the fix-ups did not touch `package.json`/`pkg.assets` (empty diff), so my round-01 partial discharge (macOS portable binary, 70/70 valid) still covers the layout-sensitive half; Windows-OS execution remains unproven. G8: correctly re-aimed at the post-merge remote (NIT-2 fix). G9/G10: untestable on this host for reasons independent of this change. G11: manifest reconciled + two guarding tests + README caveat; no Gemini CLI here to run it. |

## Check 4 — per-target destination shapes

File sets are **unchanged** from round-01, verified by listing the destinations my
`install --target all --force` run wrote: codex `~/.codex/skills/parley-deck/` = {SKILL.md,
agents, references, plugin.json, gemini-extension.json, README.md, LICENSE}; agy
`~/.gemini/config/plugins/parley-deck/` = same + fabricated `skills/SKILL.md`; gemini
`~/.gemini/extensions/parley-deck/` = same as codex. `validateInstalledPayload` untouched
(verified by diff, above); `doctor` agrees at runtime.

One content change the "shapes unchanged" claim does not cover — filed as NIT-1 below.

## Findings this round

No CRITICAL, no MAJOR, no MINOR.

### NIT-1 (new, informational) — the codex/agy destinations now carry a gemini manifest that does not resolve where it sits

The rewrite is applied only for `target.kind === "gemini"`, so the codex and agy
destinations stage the repository manifest unchanged: their `gemini-extension.json` now says
`contextFileName: "skills/parley-deck/SKILL.md"`, a path that does not exist inside those
destinations (pre-move it was `"SKILL.md"`, which did). Impact today is **zero**: nothing
consumes `gemini-extension.json` in a codex or agy destination, `agy plugin validate` passes,
`doctor` passes, and the one real consumer (the Gemini CLI, reading
`~/.gemini/extensions/parley-deck`) gets the rewritten value. I file it because check 4
asks precisely whether any installed destination changed, and the honest answer is: the file
*set* did not, one file's *content* did. If uniformity is ever wanted, applying the rewrite
for every kind makes all staged manifests self-consistent at near-zero cost; leaving it is
also defensible, since the file is cargo outside gemini-kind. No action required for release.

## Check 7 — IMPLEMENTATION.md claims vs. my reproduction

Every claim I could run, I reproduced at `a05bac7`: head-commit ✓; the cycle-2 verification
block verbatim (`grep addons/ skills/` → 0 ✓; tracker validator 3/3 ✓; check.js present and
matching its SKILL.md ✓; `npm test` 249/0 ✓); the cycle-1 claims (gemini install rc 0 with
staged `"SKILL.md"` ✓; repo manifest value with existing target ✓; G7 re-run ✓). Every NOT
TESTED is genuinely untestable here and none is smuggled past as a pass. The claim "no test
assertion was weakened" extends through the fix-ups: the only test changes are two additions
and both are mutation-proven to bite. Nothing is claimed that I could not reproduce.

## Check 8 — the README panel against A4/F1/F3/F4/F5

- **F1** ✓ — neither path labelled "recommended"; universal panel first, our installer's
  distinguishing verbs in the same screenful.
- **F3** ✓ — no upstream count restated; the comparative is now attributed (round-01 NIT-1,
  fixed).
- **F4** ✓ not triggered — no revert; and the round-01 "not supported" wording is correctly
  gone, replaced by the restored-channel text with its untested caveat.
- **F5** ✓ — `npx -y skills list` remains in the panel, matching the `doctor` line after our
  own install block.
- **A4/G7 wording** ✓ — the panel's command is the one G7 proves, no `--full-depth`.

### Signoff: kimi-1 — 2026-07-29
Status: ✅ ACCEPT

All four of my round-01 findings are **FIXED**, each verified by running the thing rather
than reading the claim: the sweep is complete (0 stale references, and the skill's own
documented commands execute), the exemplars self-pass under `--strict`, the panel comparative
is attributed, and G8 is root-caused. The D-1 withdrawal is sound: both manifest consumers
resolve against their respective roots, the reconciliation diff is additive and narrow, and
the two guarding tests are mutation-proven to fail when the reconciliation is broken. G7 —
the point of the idea — reproduces exactly at `a05bac7` with no `--full-depth`. The only
remainder is one zero-impact NIT (codex/agy carry the repo-shaped manifest as cargo) and the
standing caveat that true-remote G7 and G8 are owed at merge time, both honestly recorded.
Nothing left to block on.
