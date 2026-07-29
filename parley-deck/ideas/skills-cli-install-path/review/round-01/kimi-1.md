---
idea: skills-cli-install-path
review-round: 01
agent: kimi-1
date: 2026-07-29
---

Reviewed branch `readme-skill-catalogue` at `f8e3a1c` (matches IMPLEMENTATION.md's head-commit;
worktree clean). Every check below was re-run by me on this host (macOS arm64, node/npm via
Homebrew, `skills` CLI **1.5.20** — the same version the consensus measurements used), not read
out of IMPLEMENTATION.md. Scratch evidence: `/tmp/kimi1-g7.hrrZhF` (G7), `/tmp/kimi1-g3.lb6uob`
(portable binary), `/tmp/kimi1-g3-status.txt`.

## Findings

### MAJOR-1 — the rename was not swept through shipped skill content; `addons/…` checkout-relative instructions are now false

The move updated the installer, packaging, root tests and README, but no skill's own text
followed the rename (the diff confirms the SKILL.md files were pure zero-line renames). What was
a correct path at `94a4889` is a false path at `f8e3a1c`, so this defect was **introduced by
this change**, not inherited:

- `skills/parley-design-check/SKILL.md:49` — *"Run it as `node addons/parley-design-check/bin/check.js check <paths...>` from a checkout"*. I ran it from a checkout: **`MODULE_NOT_FOUND`**. The correct path (`node skills/parley-design-check/bin/check.js`) works. The skill's own primary run instruction fails as written.
- `skills/parley-design-check/SKILL.md:372` — the documented test command `node --test "addons/parley-design-check/test/*.test.js"` matches nothing post-move.
- `skills/parley-design-check/SKILL.md:383-390` — the "Files" map lists eight `addons/parley-design-check/…` paths, all nonexistent.
- `skills/parley-tracker/SKILL.md:79-83` — the "live next to this file" map lists five `addons/parley-tracker/…` paths, all nonexistent as repo-relative references.

FINAL.md's "no change to any skill's content" non-goal was about doctrine and behaviour; it is
not a licence to ship content that the restructure itself made false. This package's own
honesty doctrine (the one `parley-design` exists to enforce) applies to run instructions in
shipped skills. The fix is a bounded path sweep, fix-up sized.

### MINOR-1 — same stale class in the tracker's filled exemplars

`skills/parley-tracker/templates/epic.md:60-71`, `story.md:11,58-75`, `subtask.md:11,49-74` and
`skills/parley-tracker/bin/validate.test.js:51` carry exemplar `files:` / `Verify:` references to
`addons/parley-tracker/…`. These are illustrative, not operative, but the templates are sold as
"filled, self-passing examples" that agents copy — copied forward, they seed tickets pointing at
paths that no longer exist. Same sweep, same commit.

### NIT-1 — the panel's "many more than this package's own installer knows about" is an unmeasured comparative

`README.md:110`. F3 bans restating upstream's agent **count** as our fact, and the panel keeps
the number out — but the comparative clause leans on the same uncounted upstream claim. We
measured discovery/install into one agent (claude-code, G7); everything beyond that is their
documentation. "…installs all five skills into whichever coding agents you have" stands on its
own; the comparison should be attributed or dropped.

### NIT-2 — G8's recorded reason is imprecise (root cause identified below)

IMPLEMENTATION.md attributes G8's no-op to "upstream bookkeeping". The precise mechanism is in
the gate table below. No deliverable is defective; the note should be corrected so the post-merge
re-test is aimed at the right thing (a **remote** install, not a local one).

No CRITICAL findings.

## Gate re-runs (all executed by me)

| Gate | My result |
|---|---|
| **G1** install/status/doctor `--target all` | **PASS, reproduced.** `install --target all --force` rc 0; `status` shows exactly the claimed 30 `valid` lines (6 detected targets × 5 skills) with the only non-valid line `project: missing protocol` (expected — this repo is not a Parley deck); `doctor --target all` rc 0, zero malformed/missing. No regression. |
| **G2** `npm pack --dry-run` | **PASS, reproduced.** 145 `skills/` entries (exactly as claimed), 0 `addons/` entries, no root `SKILL.md`. `skills/parley-deck/SKILL.md` present. |
| **G3** Windows binary | **PARTIALLY DISCHARGED — stronger than reported, still not a full pass.** I built `dist/parley-deck-skill-v1.5.0-macos-arm64` from this tree (`build:portable:current`, pkg 6.19.0) and ran **the binary itself**: `install --target all --include-undetected` into an isolated HOME → rc 0, then `status` → **70/70 skill units valid across all 14 targets, 0 malformed/missing**. This proves the rewritten `pkg.assets` embed the new layout correctly and the embedded-payload validation path works — the layout-sensitive half of G3. What remains unproven is Windows-OS execution of the .exe, which is layout-independent. `dist/` mtimes (latest v1.3.1, June) confirm the implementer did not secretly build and rightly claimed NOT TESTED. |
| **G4** legacy Gemini `contextFileName` | **PASS for the supported path, reproduced.** `~/.gemini/extensions/parley-deck/` post-install contains `SKILL.md` at its root, and the staged `gemini-extension.json` carries `contextFileName: "SKILL.md"` → resolves. `gemini` CLI is **not installed on this host** (`which gemini` empty), confirming the URL flow is untestable here. |
| **G5** `agy plugin validate` | **PASS, reproduced.** `agy plugin validate ~/.gemini/config/plugins/parley-deck` → rc 0, `✔ skills: 1 processed`, `✔ agents: 2 processed`. The destination's fabricated `skills/SKILL.md` satisfies `plugin.json`'s `"skills": ["skills/SKILL.md"]`. |
| **G6** `npm test` | **PASS, reproduced.** 247 pass, 0 fail. I diffed both changed test files (`94a4889..f8e3a1c`): only path construction changed — the claim "no test assertion was weakened" is true. |
| **G7** real install, no `--full-depth`, isolated HOME | **PASS, independently reproduced.** Cloned the branch HEAD (tracked files only — the faithful proxy for the published repo) to scratch, `HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy`: **"Found 5 skills" / "Installed 5 skills" with no flag**. Installed core contains exactly `SKILL.md`, `agents/`, `references/`; `bin`, `lib`, `package.json`, `test`, `dist` all absent. Both S1 and S3 are fixed and measured. **Scoping caveat:** `git ls-remote origin` shows only `main` — this branch is not pushed, so a true-remote G7 remains impossible pre-merge; the published-repo run is still owed at merge time. Discovery semantics operate on the unpacked tree, so residual risk is low, but it is not zero. |
| **G8** `skills update` | **NOT TESTED — confirmed correct, and root-caused.** I reproduced the exact result ("No project skills to update") in the G7 scratch HOME. The reason is not vague bookkeeping: `getProjectSkillsForUpdate` in skills@1.5.20 (`dist/cli.mjs:5837`) **skips `sourceType === "local"` by design** — `skills-lock.json` records `"sourceType": "local"` for all five skills because the source was a filesystem path, and a local copy has nothing to update *from*. G8 is therefore only meaningful post-merge against the published GitHub remote (install from `feci/parley-deck-skill`, then `skills update`). It is honestly NOT TESTED; IMPLEMENTATION.md's "upstream bookkeeping" should read "local sources are not updatable upstream". |
| **G9** Homebrew | **NOT TESTED — confirmed untestable here.** `packaging/homebrew/Formula/` in this repo is **empty** (pre-existing since May); the formula lives in the `feci/parley` tap and points at the published tarball, which does not exist for this unreleased layout. `brew` is installed on this host but there is nothing meaningful to run it against pre-release. |
| **G10** WinGet | **NOT TESTED — confirmed.** `winget` is not installed on this macOS host and cannot be. Manifests under `packaging/winget/` (1.2.0/1.3.0/1.3.1) reference GitHub-release binaries, not repo paths, so the move does not invalidate them; they are stale by version, which is release-process business, not this idea's. |
| **G11** `gemini extensions install <url>` | **DROPPED per D-1 — ruling below.** |

## Check 4 — per-target destination shapes (implementer's claim: none changed)

**Verified true, by three independent means.** (a) `validateInstalledPayload`
(`lib/installer.js:1105-1124`) is **byte-identical** to `94a4889` (diffed the extracted
function, not just the line range). (b) Only `PAYLOAD_ENTRIES.from` moved; every `.to` is
unchanged (`SKILL.md`, `agents`, `references`, `plugin.json`, `gemini-extension.json` at the
destination root; Antigravity additionally fabricates `skills/SKILL.md`). (c) I listed the
actual destinations my `install --target all --force` run wrote: **codex** `~/.codex/skills/parley-deck/` = {SKILL.md, agents, references, plugin.json, gemini-extension.json, README.md, LICENSE}; **agy** `~/.gemini/config/plugins/parley-deck/` = same + `skills/SKILL.md`; **gemini** `~/.gemini/extensions/parley-deck/` = same as codex. All three match the unchanged per-kind required lists, and `doctor` agrees at runtime.

## Check 5 — ruling on D-1: dropping the `gemini extensions install <url>` path

**Right call. It was the only honest call available in this release.**

The mechanics check out as the implementer describes: `contextFileName` is
extension-directory-relative; `--target gemini` stages `SKILL.md` at the destination root (so
`"SKILL.md"` resolves — verified above), while the URL flow treats the repository root as the
extension, where no root `SKILL.md` exists any more. One static value cannot serve both.

I considered the reshape seriously, because dropping a capability should not be the lazy win.
The viable keep-it option — `contextFileName: "skills/parley-deck/SKILL.md"` plus our installer
fabricating that nested path into the gemini destination, mirroring the Antigravity staging —
is real, but it is **untestable on this host** (`gemini` CLI absent) and A4/G11 explicitly makes
an actual run the gate; shipping an unverifiable reshape would repeat exactly the sin this idea
exists to punish. The alternatives are worse: a root `SKILL.md` symlink resurrects S1c
(symlinks discover **0**) and a real root file resurrects the S1 shadow (discovery collapses to
1) — the two defects this restructure fixed. The lost capability is narrow: `--target gemini`
writes the same destination and is verified working, and the universal installer covers
gemini-family agents through a different mechanism entirely. If a real user of the URL flow
appears, the nested-fabrication reshape is a legitimate follow-up, gated on a host that has the
CLI.

**README honesty: verified.** The old "Other channels" line advertising
`gemini extensions install https://github.com/feci/parley-deck-skill` was removed, and
`README.md:209-211` now states plainly that the path is **not supported**, why, and what to use
instead. No residual text implies the flow works.

## Check 7 — IMPLEMENTATION.md claims vs. my reproduction

Every claim I could run, I reproduced: head-commit `f8e3a1c` ✓; G1 rc 0 / 30 valid / project
line ✓; G2 145 `skills/` entries, 0 `addons/`, no root `SKILL.md` ✓; G4 dest `SKILL.md` exists ✓;
G5 "✔ skills: 1 processed" ✓; G6 247/0 ✓; G7 verbatim output ✓ (including the exact core
contents and the five absences). Every NOT TESTED is genuinely untestable or genuinely a no-op
here — none is smuggled past as a pass. The two imprecisions I found are NIT-2 (G8's reason) and
the G3 row underselling what is provable (I discharged the pkg-assets half above). Nothing is
claimed that I could not reproduce.

## Check 8 — the README panel against A4/F1/F3/F4/F5

- **F1** ✓ — neither path is labelled "recommended"; the universal panel is first and our
  installer's distinguishing verbs (`--target all`, `doctor`, plus `status`/`sync-project` in the
  reference below) are visible in the same screenful.
- **F3** ✓ in letter — "most agents" plus the link, no upstream count restated. See NIT-1 for
  the comparative clause.
- **F4** ✓ not triggered — the move was not reverted and G7 passes, so the core-only fallback
  wording is correctly absent.
- **F5** ✓ — `npx -y skills list` is in the panel, matching the `doctor` line after our own
  install block.
- **A4/G7 wording** ✓ — the panel's command (`npx -y skills add feci/parley-deck-skill`) is the
  one G7 proves, with no `--full-depth`.
- "Fourteen named runtimes" (`README.md:120`) is our own count and checks out: codex, claude,
  agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe, cursor, opencode, aionrs = 14,
  all 14 exercised by my portable-binary run.

### Signoff: kimi-1 — 2026-07-29
Status: 🟡 ACCEPT-WITH-RESERVATIONS

The architecture is right and proven: the move is atomic with the installer and packaging
changes, destination shapes are demonstrably unchanged, G7 — the point of the idea — reproduces
exactly, and every untestable gate is honestly NOT TESTED with a verifiable reason (two of which
I strengthened: G3's assets half now proven, G8 root-caused). D-1 is ruled the right call and is
honestly documented. Reservation, required in fix-up before release: **MAJOR-1** (sweep the
`addons/…` checkout-relative references in shipped skill content — the design-check skill's own
run command is proven broken) and **MINOR-1** (the tracker exemplars in the same sweep). NIT-1
(panel comparative) and NIT-2 (G8 reason) are one-line edits. None of these touches the verified
gates; all of them are this package's own honesty rule applied to itself.
