---
idea: meta-protocol-change-lean-organizer
author: kimi-1
role: independent release-channel verifier
artifact: release-1.49.1-audit-kimi-1
date: 2026-09-24
subject: CLI 1.49.1 release channels — reviewed commit/tag target 54e07981740d495f08cd90842168d9228c58c90f
basis: owner decision parley-deck/inbox/user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md; owner scope handoff codex-1-to-user_meta-protocol-change-lean-organizer_done.md; own wording review release-wording-1.49.1-review-kimi-1.md (PASS at 54e0798); prior channel audit release-audit-kimi-1.md (1.49.0 / skill 2.13.0)
scope: release-only resume — no A–D/code review, no product changes, no repeat full suites, no Windows re-test (owner authorized labelled assets; known Windows failures are not an unexpected finding)
audit-downloads: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-24-lean-organizer/release-1.49.1/audit-kimi-1/
---

# Independent release-channel audit — CLI 1.49.1 — kimi-1

## OVERALL VERDICT: **PASS**

All channels verified independently against real artifacts (organizer assertions not
trusted). Findings: none blocking; one pending item (npm skill publish) which per the
owner decision is **not a CLI gate**.

## 1. Git / GitHub CLI release — PASS

- **Tags resolve exactly.** Local: `v1.49.1` = `54e07981740d495f08cd90842168d9228c58c90f`,
  `v1.49.0` = `06e563e8b1fe8094132149b83e14da7be9aae51e`. `git ls-remote origin` shows
  identical values for both tags; **previous tag unchanged**. Both are lightweight tags
  (objecttype commit). Remote `refs/heads/main` = `54e0798…` — the reviewed commit **is**
  the current tip of remote main (integrated; no descendants yet).
- **Evidence source checkout clean.** `release-1.49.1/source` HEAD = `54e0798…`,
  `git status --porcelain` empty.
- **Release object** (api.github.com/repos/feci/parley-deck-cli/releases/tags/v1.49.1):
  tag v1.49.1, "Parley Deck CLI 1.49.1", draft=false, prerelease=false, published
  2026-09-24T21:18:37Z.
- **Six assets, names/platforms correct** (darwin-arm64, darwin-x64, linux-arm64,
  linux-x64, windows-arm64.exe, windows-x64.exe). Independently downloaded all six into
  `audit-kimi-1/`; **three-way SHA256 match** (staged `assets/`+`sha256.json` = GitHub API
  `digest` = my download):

  | asset | sha256 |
  |---|---|
  | parley-v1.49.1-darwin-arm64 | `70545193472b05df484d5bab339c3ea14f7f3079f3a452f109185a78cf7f56f5` |
  | parley-v1.49.1-darwin-x64 | `5007e4cc2169bd4407520b371309ec4c9cdf9c67ae0124b09c0c6d20761f6a9d` |
  | parley-v1.49.1-linux-arm64 | `f489bcbab967817f789a83e1dabcff66ee49c06ed8c0b806412d80038c383c7c` |
  | parley-v1.49.1-linux-x64 | `502cef4ba1f54c16555cb4c318a274d292621de14505d0e5b99d96313b71ef15` |
  | parley-v1.49.1-windows-arm64.exe | `2dcc2c60c5078e8ce1333b445a61b6465bd64a8620067d6c051b441f8dbbe76b` |
  | parley-v1.49.1-windows-x64.exe | `b6cf7644d6770af446acb5d957eb156befc3271720ab160fc1ae0ca61739cd41` |

- **Embedded Go VCS provenance** (`go version -m` on all six downloads):
  `vcs.revision=54e07981740d495f08cd90842168d9228c58c90f`, `vcs.modified=false`,
  `-trimpath=true`, `CGO_ENABLED=0`, GOOS/GOARCH matching each asset name,
  `vcs.time=2026-09-24T21:12:53Z`. Unmodified, built from the reviewed commit.
- **Notes** byte-identical (`diff`) to the reviewed
  `release-1.49.1/release-notes.md`; they carry the Windows
  **experimental/unvalidated** label and the **CLI winget hold** in text (notes lines
  47–51). **Asset labels**: both Windows assets explicitly labelled
  `Windows ARM64 — experimental/unvalidated` / `Windows x64 — experimental/unvalidated`;
  the four macOS/Linux assets unlabelled. Windows binaries not executed — per owner
  authorization the labelled assets ship as-is.

## 2. Homebrew tap — PASS

- Remote tap `feci/homebrew-parley` HEAD = `52b09f39c8f4a98612dbf052b1b417bd2f043eb2`
  (expected), dated 2026-09-24T21:19:22Z, "release CLI 1.49.1 in Homebrew".
- **Parent diff** (`697901e` → `52b09f3`): touches **only** `Formula/parley-deck-cli.rb`,
  +2/−2 (url v1.49.0→v1.49.1, sha256 updated). **Skill formula unchanged at 2.13.0.**
- CLI formula url = `https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.49.1.tar.gz`,
  sha256 `6c0760611a4c4242418015c70e7de1508eaacbaf763a80c27a5a549a1b7a31f4`.
  **Independent download of that archive hashes to exactly that value**, and also matches
  the retained evidence copy `release-1.49.1/parley-deck-cli-v1.49.1.tar.gz`
  (top-level dir `parley-deck-cli-1.49.1/`).
- Skill formula url = `…/parley-deck-skill/archive/refs/tags/v2.13.0.tar.gz`,
  sha256 `0d585386219fe1939bd37fd27ba697b25f5b932ae7415657a16833970e44ff44`;
  independent download of that real archive matches exactly.

## 3. Local install — PASS

- `which parley` = `/opt/homebrew/bin/parley` → symlink →
  `/opt/homebrew/Cellar/parley-deck-cli/1.49.1/bin/parley`; `parley version` =
  `parley 1.49.1`. `brew list --versions`: `parley-deck-cli 1.49.1`.
- Skill installer unchanged: `/opt/homebrew/bin/parley-deck-skill` →
  `/opt/homebrew/Cellar/parley-deck-skill/2.13.0/…`; `--version` = `2.13.0`;
  `brew list --versions`: `parley-deck-skill 2.13.0`.

## 4. winget (microsoft/winget-pkgs) — PASS (hold confirmed)

Search scope (GitHub issue/PR search API, audit time 2026-09-24T21:2xZ):
`repo:microsoft/winget-pkgs` with `parley-deck-cli in:title type:pr` (0),
`Feci.ParleyDeckCli+1.49 type:pr` (0), `Feci.ParleyDeckCli type:pr` (35 total, all
closed, latest = #438660 "Update: Feci.ParleyDeckCli to 1.48.0" — historical, merged,
unrelated), and broad `parley type:pr` (97; no Feci.ParleyDeckCli 1.49.x).

- **No CLI 1.49.0 or 1.49.1 PR submitted or open** — the owner-mandated CLI winget hold
  is in effect. No winget mutation performed.
- Skill PR [#440360](https://github.com/microsoft/winget-pkgs/pull/440360)
  "Update: Feci.ParleyDeckSkill to 2.13.0": state closed, **merged**
  2026-09-24T11:16:23Z (reported as current state only).

## 5. npm (parley-deck-skill) — PENDING, not a CLI gate

- `npm view parley-deck-skill` at 2026-09-24T21:27:50Z: **latest = 2.12.1**,
  dist.integrity `sha512-45s2sww7d7LVtKgIJ9jaOyHvfl1f70mk1r8/luw+4B+8eU2hT5Wwx5L7HVWQLvmtdENMi0Hp4JczEtISuPGmug==`,
  dist.tarball `https://registry.npmjs.org/parley-deck-skill/-/parley-deck-skill-2.12.1.tgz`.
- 2.13.0 is **not yet published**; per the owner decision this is explicitly **pending
  with the owner** (owner-facing Claude session handles login/publish) and is **NOT a CLI
  gate**. No login or publish action started or attempted. Since latest is still 2.12.1,
  no 2.13.0 registry tarball exists to download/match.

## 6. Skill GitHub 2.13.0 — PASS (unchanged)

- Tag `v2.13.0` → commit `8161e5e9d3ecdece25b7ae225e2304d3159be79b` (same immutable tag
  as the prior audit); release "Parley Deck Skill 2.13.0", draft=false, published
  2026-09-24T09:52:31Z, five assets (linux-x64, macos-arm64, macos-x64,
  windows-arm64.exe, windows-x64.exe). **Prior audit release-audit-kimi-1.md remains
  applicable**; nothing rebuilt or modified.

## Findings

None. All 1.49.1 CLI channels verified; the only open item is the owner-held npm skill
publish (2.12.1 latest), which is out of the CLI gate by owner decision.

## Limitations / not done

- No A–D or product-code review, no source test re-runs, no full protocol/skill re-reads.
- Windows binaries hash- and provenance-verified but not executed; known Windows failures
  are covered by the owner-authorized experimental/unvalidated label, not re-litigated.
- Homebrew bottle contents not rebuilt/compared (formula builds from source at install);
  local install verified by version + Cellar path instead.
- npm 2.13.0 tarball match not possible (not published); exact audited tarball remains
  retained under `release-delivery/2026-09-24-lean-organizer` for the owner session.
- Remote main currently equals the reviewed commit; a later records-only descendant would
  not invalidate this audit.
- This is an auxiliary release audit, not a protocol round; no commits/pushes/tags/
  publishing/logins/installs performed, no done-file written, others' artifacts untouched.
