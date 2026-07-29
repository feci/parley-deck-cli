---
idea: skills-cli-install-path
review-round: 01
agent: hermes-1
date: 2026-07-29
---

# Review by hermes-1 — 2026-07-29

Non-implementer adversarial review. Every gate below was re-run by me, not taken
from IMPLEMENTATION.md on trust. Host: macOS 26.5.2, Node v26.5.0, npm 11.17.0,
skills CLI 1.5.20. Repo branch: readme-skill-catalogue, HEAD f8e3a1c.

## Findings

### MAJOR — G7 cannot be run against the *published* repo before merge

FINAL.md gate G7 (as strengthened by kimi-1 in A4) specifies "an actual install
from the *published* repo into an isolated HOME". The branch has not been pushed
to origin (only `main` and feature branches are on the remote). I ran G7 against
a local rsync copy of the repo at the correct HEAD, which proves the *layout*
works, but the definition-of-done sentence "from the published repository" is
not yet satisfiable. This is inherent to pre-merge review, not a deficiency in
the implementation. The implementer's G7 result used `<repo>` as the source — a
local path, not the merged remote — so the same caveat applies to their run.

**Recommendation:** ship after merge, then re-run G7 against
`npx -y skills@latest add feci/parley-deck-skill` from the published `main` to
fully close the gate. The local-copy result is strong evidence it will pass.

### MINOR — G8 (skills update) confirmed NOT TESTED, reason is upstream

I reproduced the implementer's G8 result exactly. After a `--copy` install that
`skills list` sees (all 5 skills listed with correct source paths), `skills update`
reports "No project skills to update." The skills CLI's update mechanism does not
track local-path-installed skills for version comparison — it expects registry or
npm-sourced packages with version metadata. This is upstream bookkeeping, not a
defect in our tree. The implementer's assessment is correct and honest.

### MINOR — WinGet manifests are stale (1.0.4 through 1.3.1, package is 1.5.0)

The WinGet manifest versions in `packaging/winget/manifests/` top out at 1.3.1
while `package.json` is 1.5.0. This is a pre-existing condition carried forward
by the readme-skill-catalogue branch, not introduced by this idea. Consensus
follow-up #2 already flags `packaging/winget/README.md` as stale. No action
required for this idea's gates, but it should be addressed before the 1.5.0
release.

### NIT — pkg.assets has a redundant explicit entry

`package.json` `pkg.assets` lists both `"skills/**/*"` (line 46) and
`"skills/parley-deck/agents/**/*"` (line 47) plus
`"skills/parley-deck/references/**/*"` (line 48). The glob `skills/**/*`
already covers everything under `skills/`, making the two specific entries
redundant. This is harmless (pkg deduplicates), but the explicit
`skills/parley-deck/SKILL.md` entry (line 43) is also covered by the glob.
Not worth changing for this release; noted for cleanup.

## Gate-by-gate re-run results

### G1 — install --target all, status, doctor: PASS (reproduced)

Ran in an isolated HOME with `--include-undetected --force`:
- `install --target all --include-undetected --force`: exit 0, all 14 targets x
  5 skills = 70 "installed" lines.
- `status --target all --include-undetected`: exit 0, 70 "valid" lines. The only
  non-valid line is `project: missing protocol`, expected (this repo is not a
  Parley deck project).
- `doctor --target all --include-undetected`: exit 0, 70 "valid" lines, zero
  non-valid.

Also reproduced the implementer's "30 valid lines" figure: with real detection
(codex, claude, agy, hermes, kimi commands on PATH), `status --target all`
without `--include-undetected` yields exactly 30 valid lines (6 detected targets
x 5 skills). No regression.

### G2 — npm pack ships skills/, not addons/ or root SKILL.md: PASS (reproduced)

`npm pack --dry-run`: 153 total files. 145 entries under `skills/`, 0 under
`addons/`, 0 root `SKILL.md`. `skills/parley-deck/SKILL.md` is present. The
tarball ships `plugin.json` and `gemini-extension.json` at the root (correct —
they are repo-level manifests, not skill-internal files). No `addons/` directory,
no root `SKILL.md`.

### G3 — Windows binary build: NOT TESTED (no Windows host), partial proof

`node scripts/build-portable.js current` built a macOS arm64 binary
(`dist/parley-deck-skill-v1.5.0-macos-arm64`, 65 MB). The binary reports version
1.5.0 and successfully installs into an isolated HOME with the new layout (5
skills, correct destination shape). This proves `pkg.assets` embeds the `skills/`
tree correctly. The Windows binary cannot be built or tested on macOS. The
`pkg.assets` configuration is updated and proven by the macOS build. Old Windows
binaries in `dist/` (1.0.4, 1.1.1, 1.3.1) are from prior releases.

### G4 — contextFileName / destination shapes: PASS (reproduced, verified unchanged)

Listed the installed destinations for codex, antigravity, and gemini:

- **Codex** (`~/.codex/skills/parley-deck/`): SKILL.md, agents/, references/,
  plugin.json, gemini-extension.json, README.md, LICENSE, marker. NO bin/, lib/,
  package.json, test/, dist/. Shape unchanged from pre-move.
- **Antigravity** (`~/.gemini/config/plugins/parley-deck/`): same as codex PLUS
  `skills/SKILL.md` (the fabricated file for `plugin.json`'s `"skills"` path).
  Shape unchanged.
- **Gemini** (`~/.gemini/extensions/parley-deck/`): same as codex.
  `gemini-extension.json` at dest root with `contextFileName: "SKILL.md"`, which
  resolves to the `SKILL.md` at the dest root. Shape unchanged.

`validateInstalledPayload` (installer.js:1105-1124) is unchanged — its required
file lists per kind (addon, gemini, antigravity, default) all reference
destination-relative paths that the installer still produces. The implementer's
claim that "only PAYLOAD_ENTRIES.from moved, the destination shape did not" is
verified by reading the source and by listing actual installed destinations.

### G5 — agy plugin validate: PASS (reproduced)

`agy plugin validate ~/.gemini/config/plugins/parley-deck` (against an installed
destination): `skills: 1 processed`, `agents: 2 processed`, exit 0.

Also ran `agy plugin validate` against the repo root: `skills: 5 processed`,
`agents: skipped (not found)`, exit 0. The repo-root validation finds 5 skills
because `skills/` is now a real directory with five `SKILL.md` files. The
`agents` path in `plugin.json` (`agents/manifest.yaml`) does not exist at the
repo root (it's at `skills/parley-deck/agents/`), so agents are skipped — this
is expected, since `plugin.json` paths resolve relative to the installed
destination, not the source repo.

### G6 — npm test: PASS (reproduced)

247 pass, 0 fail, 0 skipped. Matches the implementer's claim exactly.

### G7 — real npx skills add, no --full-depth: PASS (against local copy)

Ran `HOME=<scratch> npx -y skills@latest add <local-repo-copy> --agent
claude-code --yes --copy` with NO `--full-depth`:

- "Found 5 skills" / "Installed 5 skills" — all five discovered and installed
  with no flag. Both S1 (discovery) and S3 (whole-repo-as-skill) are fixed.
- Installed core (`.claude/skills/parley-deck/`) contains exactly: `SKILL.md`,
  `agents/` (manifest.yaml, openai.yaml), `references/` (COOPERATION.md,
  WORKED_EXAMPLES.md, compatibility.json).
- Absent from installed core: `bin/`, `lib/`, `package.json`, `test/`, `dist/`,
  `scripts/`, `addons/`, `node_modules/`, `packaging/`, `package-lock.json`.
  All confirmed absent.
- `--list` also finds 5 with no `--full-depth`.

See the MAJOR finding above: the FINAL.md spec says "from the published repo",
which requires post-merge verification. The local-copy result is strong evidence
the published result will match.

### G8 — skills update: NOT TESTED (upstream limitation, reproduced)

`skills update` reports "No project skills to update" even though `skills list`
sees all 5 installed skills with correct source paths. The update mechanism does
not track local-path-installed skills. This is upstream bookkeeping, not our
layout. The implementer's assessment is correct.

### G9 — Homebrew: NOT TESTED (formula in separate tap, untestable pre-release)

The Homebrew formula lives in a separate tap repository
(`feci/homebrew-parley`), not in this repo. `packaging/homebrew/Formula/` is
intentionally empty (removed in commit 9ca2a46). `brew upgrade` and `brew test`
can only run against a published release tarball. NOT TESTED is correct.

### G10 — WinGet: NOT TESTED (winget not installed on macOS)

`winget` is not available on macOS. The manifests exist (versions 1.0.4 through
1.3.1) but are stale relative to package version 1.5.0. NOT TESTED is correct.

### G11 — gemini extensions install <url>: NOT TESTED (gemini CLI not installed)

`gemini` CLI is not installed on this host. The README honestly documents this
path as unsupported. See D-1 ruling below.

## Check #4 — per-target destination shape unchanged: VERIFIED

The implementer claims "only PAYLOAD_ENTRIES.from moved, the destination shape
did not." I verified this three ways:

1. **Source reading:** `PAYLOAD_ENTRIES` (installer.js:134-140) became
   `{from, to}` pairs. Every `from` now resolves under `CORE_SKILL_DIR`
   (`skills/parley-deck/`) or the repo root (for `plugin.json` and
   `gemini-extension.json`). Every `to` is unchanged — `SKILL.md`, `agents`,
   `references`, `plugin.json`, `gemini-extension.json` at the destination root.
2. **`validateInstalledPayload` reading:** (installer.js:1105-1124) — the
   required file lists per kind are destination-relative and unchanged. No
   source path appears in them.
3. **Actual installed destinations:** listed codex, antigravity, and gemini
   destinations after a real install. All have the expected shape. No target
   gained or lost files compared to the pre-move layout.

## Check #5 — D-1 ruling: dropping the gemini extension-URL path is the right call

The analysis is correct and the decision is sound:

- `gemini-extension.json` has `contextFileName: "SKILL.md"`, relative to the
  extension directory.
- `install --target gemini` copies `SKILL.md` to the destination root, so
  `"SKILL.md"` resolves. Verified: `~/.gemini/extensions/parley-deck/SKILL.md`
  exists after install.
- `gemini extensions install <url>` treats the **repository root** as the
  extension directory. The repo root no longer has `SKILL.md` (it moved to
  `skills/parley-deck/SKILL.md`), so `"SKILL.md"` would resolve to nothing.
- The alternative — setting `contextFileName` to
  `"skills/parley-deck/SKILL.md"` — would fix the URL path but break
  `install --target gemini`, because the installer copies `SKILL.md` to the
  destination root, not to `skills/parley-deck/SKILL.md` in the destination.
- Reshaping the installer to also create `skills/parley-deck/SKILL.md` in the
  gemini destination would change the destination shape, which the implementer
  correctly avoided.
- Restoring a root `SKILL.md` to keep the URL path would re-introduce the
  shadow rule (S1a) that causes the discovery defect this idea exists to fix.

One value cannot serve both consumers. Dropping the URL path is the correct
call. The README is honest about it: "Installing this repository through
`gemini extensions install <url>` is **not supported**: that path treats the
repository root as the extension, and the core skill no longer lives there."
The supported `--target gemini` path is verified working.

## Check #6 — testability of G3/G8/G9/G10 on this host

- **G3 (Windows binary):** NOT TESTABLE — no Windows host. Partial proof via
  macOS binary build (proves pkg.assets is correct).
- **G8 (skills update):** TESTABLE, TESTED — result: NOT TESTED (upstream
  limitation). `skills update` does not track local-path installs.
- **G9 (Homebrew):** NOT TESTABLE — formula is in a separate tap repo, points
  at a published tarball that doesn't exist yet for 1.5.0.
- **G10 (WinGet):** NOT TESTABLE — `winget` is not installed on macOS.

No gate that could be tested was left untested.

## Check #7 — IMPLEMENTATION.md claims I cannot reproduce

Every claim in IMPLEMENTATION.md that I could test, I reproduced:

- G1 (30 valid lines with real detection): reproduced exactly.
- G2 (145 skills/ entries, 0 addons/, 0 root SKILL.md): reproduced exactly.
- G5 (agy plugin validate, 1 skill processed): reproduced exactly.
- G6 (247 pass, 0 fail): reproduced exactly.
- G7 (Found 5, Installed 5, core has references/agents, no bin/lib/package.json):
  reproduced against a local copy (see MAJOR finding re: published repo).
- D-1 (contextFileName cannot serve both consumers): verified by analysis and
  by confirming `--target gemini` works while the URL path cannot.
- Head commit f8e3a1c: matches `git log --oneline -1`.

The four NOT TESTED gates (G3, G8, G9, G10) are honestly reported as such with
correct reasons. No NOT TESTED gate is reported as a pass.

## Check #8 — README panel wording vs A4/F1/F3/F4/F5

- **F1 (neither path labelled "recommended"):** PASS. No occurrence of
  "recommend" anywhere in the README. The panel says "One command, most agents"
  and our own installer section says "covers fourteen named runtimes and adds
  detection, health checks and project-metadata sync." Both are presented
  co-equal, neither is labelled recommended.
- **F3 (no agent count of our own):** PASS. The panel says "most agents" and
  "many more than this package's own installer knows about." No "70" or "75"
  anywhere. Our own installer says "fourteen named runtimes" — a count we
  measured (14 TARGETS entries verified in installer.js).
- **F4 (if the move fails/reverts, panel must say core only):** N/A — the move
  succeeded and is not reverted. The panel says "installs all five skills,"
  which G7 proves true. If the move were reverted, F4's fallback clause would
  require rewriting the panel — but it wasn't reverted.
- **F5 (skills list verification line in the panel):** PASS. `npx -y skills
  list` is in the panel's code block, matching the `doctor` line that follows
  our own install block.
- **A4 (G7 is an install, not a listing):** The panel's claim "installs all five
  skills" is backed by a real install (G7), not just a `--list`. The
  implementer's G7 ran a real install and so did mine.

No claim in the README panel that we did not measure. The "fourteen named
runtimes" count is verified. The "all five skills" claim is verified by G7.
The "many more than this package's own installer knows about" is a qualitative
comparison (the universal CLI supports many agents beyond our 14) that is
self-evidently true and makes no specific numeric claim.

## Additional observations

The git diff (94a4889..f8e3a1c) is surgical: 150 files changed, but 145 are pure
renames (0 bytes changed). Only 5 files have content changes: `README.md`,
`lib/installer.js`, `package.json`, `test/design-addons.test.js`,
`test/installer.test.js`. No skill content was modified. No test assertion was
weakened — only the paths they build and read were repointed. The test file's
synthetic package (installer.test.js:207-231) builds the new layout correctly,
staging `skills/parley-deck/` with the core payload and `plugin.json` +
`gemini-extension.json` at the package root.

The `discoverAddons` skip of `CORE_SKILL_NAME` (installer.js:771) is correct:
the core skill shares the `skills/` directory so a generic installer sees all
five as siblings, but our own installer must never offer the core as an add-on
to itself. The test "discovers packaged add-on skills" (installer.test.js:406-412)
asserts exactly 4 add-ons, not 5.

### Signoff: hermes-1 — 2026-07-29
Status: ✅ ACCEPT

The layout move, installer repointing, packaging change, and README panel are
all correct, surgical, and verified by running them. Every gate that can be
tested on this host was tested and passed or was honestly reported as NOT
TESTED with a correct reason. D-1 is the right call. The one open item (G7
against the published remote) is inherent to pre-merge review and is strong
evidence-based, not assertion-based. The implementation matches FINAL.md and
the A1-A4 amendments. I accept this for merge, with the recommendation to
re-run G7 against the published `main` immediately after merge as the final
definition-of-done check.
