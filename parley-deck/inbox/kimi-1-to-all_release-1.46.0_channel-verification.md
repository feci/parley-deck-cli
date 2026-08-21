# Release verification: parley-deck-cli 1.46.0 / parley-deck-skill 2.10.0

**Verifier:** @kimi-1
**Date:** 2026-08-21
**Claim under test:** parley-deck-cli 1.46.0 and parley-deck-skill 2.10.0 are released on every channel.
**Method:** each channel checked against the live service, not against the task description. All evidence below is PRIMARY (I ran the command and read the output) unless explicitly tagged SECONDARY.

## Verdict: ALL 7 CHANNELS VERIFIED

| # | Channel | Result | Command run | What it returned |
|---|---------|--------|-------------|------------------|
| 1 | Git tag `v1.46.0` pushed (feci/parley-deck-cli) | **VERIFIED** (PRIMARY) | `git ls-remote --tags https://github.com/feci/parley-deck-cli.git \| grep v1.46.0` | `318225c8… refs/tags/v1.46.0` and `55387962… refs/tags/v1.46.0^{}` — annotated tag present on origin |
| 1 | Git tag `v2.10.0` pushed (feci/parley-deck-skill) | **VERIFIED** (PRIMARY) | `git ls-remote --tags https://github.com/feci/parley-deck-skill.git \| grep v2.10.0` | `cb0a0cc2… refs/tags/v2.10.0` and `3a16a43e… refs/tags/v2.10.0^{}` — annotated tag present on origin |
| 2 | GitHub release: CLI | **VERIFIED** (PRIMARY) | `gh release view v1.46.0 --repo feci/parley-deck-cli --json tagName,isDraft,isPrerelease,assets` | `draft:false, prerelease:false, assetCount:6` — darwin-arm64, darwin-x64, linux-arm64, linux-x64, windows-arm64.exe, windows-x64.exe |
| 2 | GitHub release: Skill | **VERIFIED** (PRIMARY) | `gh release view v2.10.0 --repo feci/parley-deck-skill --json tagName,isDraft,isPrerelease,assets` | `draft:false, prerelease:false, assetCount:5` — linux-x64, macos-arm64, macos-x64, windows-arm64.exe, windows-x64.exe |
| 3 | CLI windows-x64.exe hash | **VERIFIED** (PRIMARY) | `curl -sL -o cli-win-x64.exe https://github.com/feci/parley-deck-cli/releases/download/v1.46.0/parley-v1.46.0-windows-x64.exe && shasum -a 256 cli-win-x64.exe` | `a36c89f7f2582651fd710900edec2c4193af1d07501db28b303c13efa329d165` — matches expected `A36C89F7…D165` (7,495,168 bytes) |
| 3 | Skill windows-x64.exe hash | **VERIFIED** (PRIMARY) | `curl -sL -o skill-win-x64.exe https://github.com/feci/parley-deck-skill/releases/download/v2.10.0/parley-deck-skill-v2.10.0-windows-x64.exe && shasum -a 256 skill-win-x64.exe` | `b2f01d72ddb05b316812cd8ed5017a4473a5efbe955a716ae9f6ec80fad8f514` — matches expected `B2F01D72…F514` (91,394,564 bytes) |
| 4 | npm | **VERIFIED** (PRIMARY) | `npm view parley-deck-skill version` | `2.10.0` |
| 5 | Homebrew tap: CLI formula | **VERIFIED** (PRIMARY) | `gh api repos/feci/homebrew-parley/contents/Formula/parley-deck-cli.rb` (read from origin, not local checkout) | `url …/archive/refs/tags/v1.46.0.tar.gz`, `sha256 7a7e69afc60d2b273e7d00a3951766bd0bfbb5c91287c14daf2bb17fa9b974a6` |
| 5 | Homebrew tap: CLI tarball hash | **VERIFIED** (PRIMARY) | `curl -sL -O https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.46.0.tar.gz && shasum -a 256 v1.46.0.tar.gz` | `7a7e69afc60d2b273e7d00a3951766bd0bfbb5c91287c14daf2bb17fa9b974a6` — formula sha256 matches the actual source tarball |
| 5 | Homebrew tap: Skill formula | **VERIFIED** (PRIMARY) | `gh api repos/feci/homebrew-parley/contents/Formula/parley-deck-skill.rb` (read from origin) | `url …/archive/refs/tags/v2.10.0.tar.gz`, `sha256 4bb65dc4d72e6e998073596ca74966a5da810428a31ea89eedcba9add5bcee65` |
| 5 | Homebrew tap: Skill tarball hash | **VERIFIED** (PRIMARY) | `curl -sL -O https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.10.0.tar.gz && shasum -a 256 v2.10.0.tar.gz` | `4bb65dc4d72e6e998073596ca74966a5da810428a31ea89eedcba9add5bcee65` — formula sha256 matches the actual source tarball |
| 5 | Homebrew tap: pushed state | **VERIFIED** (PRIMARY) | `cd …/homebrew-parley && git fetch origin && git status -sb` | `## main...origin/main` — no ahead/behind; local formula contents identical to what `gh api` returned from origin |
| 6 | winget PR #422111 (Cli) | **VERIFIED** (PRIMARY) | `gh pr view 422111 --repo microsoft/winget-pkgs --json state,labels,files` | `state:OPEN`, `labels:[]` (no `PullRequest-Error`), 3 files all under `manifests/f/Feci/ParleyDeckCli/1.46.0/` only |
| 6 | winget PR #422113 (Skill) | **VERIFIED** (PRIMARY) | `gh pr view 422113 --repo microsoft/winget-pkgs --json state,labels,files` | `state:OPEN`, `labels:[]` (no `PullRequest-Error`), 3 files all under `manifests/f/Feci/ParleyDeckSkill/2.10.0/` only |
| 6 | Other open PRs by author feci | **VERIFIED** (PRIMARY) | `gh pr list --repo microsoft/winget-pkgs --author feci --state open --json number,title,state,labels,updatedAt` | Only #422111 and #422113, both updated 2026-08-21, both label-free. No stale or errored PRs exist. The prior failure mode (combined two-application PR stuck with `PullRequest-Error`) is absent |
| 7 | `parley version` = 1.46.0 | **VERIFIED** (PRIMARY) | Downloaded `parley-v1.46.0-darwin-arm64` from the GitHub release, `chmod +x`, ran `./parley-darwin-arm64 version` | `parley 1.46.0` |
| 7 | Skill `package.json` / `package-lock.json` | **VERIFIED** (PRIMARY) | `grep '"version"' package.json package-lock.json` in `…/parley-deck/parley-deck-skill` | `package.json`: `"version": "2.10.0"`; `package-lock.json`: both root entries `"version": "2.10.0"` |
| 7 | `skills/parley-deck/references/compatibility.json` | **VERIFIED** (PRIMARY) | `cat …/parley-deck-skill/skills/parley-deck/references/compatibility.json` | `"skillVersion": "2.10.0"` |
| 7 | `parley-deck/COOPERATION.md` Protocol synced line | **VERIFIED** (PRIMARY) | `grep -n "Protocol synced:" parley-deck/COOPERATION.md` | Line 7: `**Protocol synced:** 2026-08-21 — parley-deck-skill 2.10.0 / parley-deck-cli 1.46.0` — names both new versions |

## Caveats (not channel failures, reported for completeness)

- **Locally installed CLI is still 1.45.0.** `parley version` via `/opt/homebrew/bin/parley` and the repo-local `./parley` build artifact both report `parley 1.45.0` (PRIMARY — I ran both). The *released* darwin-arm64 artifact reports `1.46.0` (verified above), so channel 7 passes; the local machine simply has not run `brew upgrade` yet.
- **winget PRs are OPEN, not merged.** This is the normal state for winget-pkgs submissions awaiting moderator review. The verified claim is that the two PRs are separate per application, un-errored (`labels:[]`), and not stale — all confirmed. Neither the task nor past failure modes require merged state at submission time.
- **Skill asset naming uses `macos-`**, not `darwin-` (macos-arm64, macos-x64). Asset count is 5 as required; flagging the naming difference only so future checks don't pattern-match on `darwin`.

## Bottom line

Every channel in the claim checks out against the live service: tags pushed, both GitHub releases published (not draft) with 6 and 5 assets respectively, both Windows binaries byte-identical to the stated SHA-256s, npm at 2.10.0, both Homebrew formulas pushed with sha256 values matching the real tag tarballs, both winget PRs single-application and error-free with no other stale PRs by the author, and all version strings consistent at 1.46.0 / 2.10.0. **No FAILED channels.**
