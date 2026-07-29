---
idea: skills-cli-install-path
implementer: claude-1
date: 2026-07-29
status: ready-for-review
target: parley-deck-skill
head-commit: f8e3a1c
prior-commits: [951d7a5 the move + installer + panel]
---

## What was done

The layout move from `FINAL.md`, the installer repointing, and the README panel.

```
skills/parley-deck/{SKILL.md, references/, agents/}   ← was the repository root
skills/parley-design/  parley-design-check/  parley-tracker/  parley-worktrees/
                                                      ← was addons/*
```

`addons/` is gone. A **rename**, not a copy — no second skill tree, so no drift guard is
needed (S4b). `plugin.json` and `gemini-extension.json` stayed at the repository root per A3,
with the installer's *source* path for them set explicitly.

### Installer (`lib/installer.js`)

`CORE_SKILL_DIR` and `CORE_SKILL_NAME` introduced. `REQUIRED_PAYLOAD_FILES` and
`COMPATIBILITY_FILE` resolve under it. `PAYLOAD_ENTRIES` / `OPTIONAL_PAYLOAD_ENTRIES` became
`{from, to}` pairs — **the source moved, the destination shape did not**, so an installed skill
still has `SKILL.md` at its root and every per-target validator is unchanged. `ADDONS_DIR`
`"addons"` → `"skills"`, with `discoverAddons` skipping `parley-deck` so the core is never
offered as an add-on. `packagedProtocolPath` and the core `skillSha256` repointed. The
Antigravity staging that fabricates `skills/SKILL.md` in the destination now reads from the
new source.

### Packaging

`package.json` `files` ships `skills/` and no longer `addons/` or a root `SKILL.md`.
`pkg.assets` (which the standalone binaries embed) rewritten the same way.

### Tests

`test/design-addons.test.js` repointed. `test/installer.test.js`'s synthetic package now
builds the new layout. **No test assertion was weakened** — only the paths they build and read.

## Gates

| Gate | Result |
|---|---|
| **G1** `install --target all`, `status --target all` | **PASS** — rc 0; 30 `valid` lines; the only non-valid line is `project: missing protocol`, which is expected (this repo is not a Parley deck) |
| **G2** `npm pack` ships `skills/`, not `addons/` or a root `SKILL.md` | **PASS** — 145 `skills/` entries, `skills/parley-deck/SKILL.md` present, 0 `addons/` entries, 0 root `SKILL.md` |
| **G3** Windows binary build + install | **NOT TESTED** — no Windows host. `pkg.assets` is updated but unproven |
| **G4** legacy Gemini `contextFileName` | **PASS for the supported path** — `install --target gemini` rc 0 and `~/.gemini/extensions/parley-deck/SKILL.md` exists, so `contextFileName: "SKILL.md"` still resolves. See the deviation below for the unsupported path |
| **G5** `agy plugin validate` | **PASS** — `✔ skills: 1 processed` |
| **G6** `npm test` | **PASS** — 247 pass, 0 fail |
| **G7** real install, no `--full-depth`, isolated `HOME` | **PASS, and it is the point of the idea.** `Found 5 skills` / `Installed 5 skills` with **no flag**. Installed core contains exactly `SKILL.md`, `agents`, `references`; `bin`, `lib`, `package.json`, `test`, `dist` are all **absent**. Both the discovery defect (S1) and the whole-repo-as-skill defect (S3) are fixed and *measured*, not asserted |
| **G8** `skills update` | **NOT TESTED** — `update` reported "No project skills to update" after a `--copy` install, so the run proved nothing about our layout. The reason is upstream bookkeeping, not our tree |
| **G9** Homebrew `brew upgrade` / `brew test` | **NOT TESTED** — the formula points at a published tarball; untestable before the release |
| **G10** WinGet `winget validate` / `winget install` | **NOT TESTED** — `winget` is not installed on this host (macOS) |
| **G11** `gemini extensions install <url>` | **NOT TESTED and now unsupported** — see below |

G7 verbatim:

```text
$ HOME=/tmp/g7/home npx -y skills@latest add <repo> --agent claude-code --yes --copy
◇  Found 5 skills
◇  Installed 5 skills
     ✓ parley-deck  ✓ parley-design  ✓ parley-design-check  ✓ parley-tracker  ✓ parley-worktrees
$ ls .claude/skills/parley-deck/     → SKILL.md  agents  references
$ test -e .claude/skills/parley-deck/{bin,lib,package.json,test,dist}  → all absent
```

## Deviation D-1 — the legacy Gemini extension-URL path is dropped, not fixed

kimi-1 called `gemini-extension.json`'s `contextFileName` "the sharpest single edge", and it
is. The value is relative to the extension directory, and the two consumers now disagree:

- `install --target gemini` copies `SKILL.md` to the destination root → `"SKILL.md"` resolves.
- `gemini extensions install <repo-url>` treats the **repository** as the extension → there is
  no root `SKILL.md` any more → it would resolve to nothing.

**One value cannot satisfy both.** I could not test the second path (`gemini` is not installed
here), and shipping a documented command that my own analysis says is broken is exactly what
`parley-design`'s honesty rule forbids. So the README line for it is **removed**, and the
README now states plainly that installing this repository through
`gemini extensions install <url>` is not supported and why. The supported `--target gemini`
path is verified working.

This is a **capability we lost**, not a defect we fixed, and it is the real price of the
restructure. It is written here rather than left for someone to discover.

## What reviewers should attack

1. The four `NOT TESTED` gates. None may be reported as a pass.
2. D-1 — is dropping the extension-URL path the right call, or should the restructure be
   reshaped to keep it?
3. Whether any per-target *destination* shape changed. I claim none did, because only
   `PAYLOAD_ENTRIES.from` moved. `validateInstalledPayload` is untouched — verify that.
4. The README panel's wording against F1/F3/F4/F5.
