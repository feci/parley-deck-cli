---
agent: zcode-1
idea: meta-protocol-change-ste100-output-format
role: independent channel verifier
date: 2026-10-13
cli-merge-commit: 88dc6c5da139d2deb039bda4b605acfbe42663b4
skill-merge-commit: f557c5e679b0c1a03f57a885f41ae54f68348b2a
reviewed-product-cli: 99d04b3ecd0740fb834e70ec109c9384c7b7cd4d
reviewed-product-skill: 92e3dc78da2f1c55d35c905cabe57fc69ede0670
process: fresh headless channel check, separate from all reviews, signoffs, and the source goal check
---

## Scope and freshness

I am zcode-1, the independent channel verifier. This process is fresh. It is separate from every source review, signoff, and the source goal check. Codex-1 organized and implemented. I implemented no part and published nothing.

I read in full: `source-context/ORGANIZER-BRIEF.md`, `FINAL.md`, `IMPLEMENTATION.md`, `review/consensus.md`, `source-context/goal-check-zcode-1.md`, and `source-context/windows-verification-zcode-1.md`. I read the delivery records in `release-delivery/2026-10-13-ste100`: `merge-commits.json`, `cli-build-verification.json`, `skill/payload-verification.json`, `hosted-assets-verification.json`, `homebrew-publication.json`, `homebrew-verification.json`, `winget-validation.json`, `winget-submission.json`, `runtime-verification.json`, `owner-publication-status.json`, `d6-preview.json`, `source-close-receipt.json`, `windows-final-comparison.json`, and `release-handoff-draft.md`.

I changed no product file, signed artifact, release asset, manifest, managed runtime, credential, or gateway setting. I created this file only. I did not publish, merge, install, apply legacy history, declare worktrees, or prune.

## Protocol context attestation

The launcher supplied one private prompt file. I hashed it before any task action.

```text
/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100/.parley-runtime/invocations/54e2db04-8c86-4cc8-ad05-ff2f901b8ea5/headless-prompt-2918760535.md
SHA256 fc9f24ceb84b79981a3230f0295c0addb096098653591cc85378462116acd063
```

The hash matched the launcher value. I then read the complete file. Its attestation, verbatim:

```text
Protocol context attestation: {"context_mode":"full","source_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec","packet_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec"}
Shadow packet audit: {"packet_sha256":"9d8b2219d8a8a949761c9f46deb359f6471f38aecd165c9cbb7ce45527203f1f","packet_bytes":81966,"source_bytes":130042,"included_blocks":37,"omitted_blocks":33}
```

Mode `full`. Renderer reachable. No fallback. `Transport: github-pr`. The shadow audit line is an unapplied diagnostic.

## Tree identity

CLI worktree HEAD: `118d0c5e646eb8b49af1b77a984bae38116bb7ee`, clean. Commits after the reviewed `99d04b3` are the main merge `88dc6c5` and idea-evidence commits. `git diff --name-only 99d04b3..HEAD` lists paths under `parley-deck/` only. No product file changed.

Skill worktree HEAD: `f557c5e679b0c1a03f57a885f41ae54f68348b2a`, the main merge commit, clean. My hashes: `SKILL.md` `64d5236baa0aece207becaa9def8e6fba24776f148dd2289281efa0855341df1`; `references/STE100.md` `f9eab1935e9681bf13360140cd3673f59598b66e64d94cb13b0ac6e35449ab63`. Both match the hosted package bytes below. Verdict: CONFIRMED (PRIMARY).

## DC1 — merges, tags, releases, assets, Windows workflow

Attempted refutation: I queried live GitHub state myself and downloaded decisive hosted bytes.

Merges and tags:

```text
gh pr view 84 --repo feci/parley-deck-cli   → state MERGED, mergeCommit 88dc6c5da139d2deb039bda4b605acfbe42663b4, base main
gh pr view 12 --repo feci/parley-deck-skill → state MERGED, mergeCommit f557c5e679b0c1a03f57a885f41ae54f68348b2a, base main
git commit 88dc6c5 → 2 parents: d51e1c7…, 4a9ec8d…   message "Merge pull request #84 from feci/ste100"
git commit f557c5e → 2 parents: 0549544…, 92e3dc7…   message "Merge pull request #12 from feci/ste100"
refs/tags/v1.55.0 → commit 88dc6c5da139d2deb039bda4b605acfbe42663b4 ; branches/main head = 88dc6c5…
refs/tags/v2.19.0 → commit f557c5e679b0c1a03f57a885f41ae54f68348b2a ; branches/main head = f557c5e…
```

Both PRs used real merge commits. Both tags resolve to the merge commits. Both merge commits are the live `main` heads.

Live releases: CLI `v1.55.0` and skill `v2.19.0` are full releases. `draft: false`, `prerelease: false`. The CLI release carries seven assets: six binaries plus `sha256.json`. The skill release carries four assets. Every live GitHub `digest` equals the recorded value in `hosted-assets-verification.json` and `cli-build-verification.json`.

I downloaded hosted bytes and hashed them myself:

```text
skill-2.19.0.tgz              f07102c862908b3f6da6fb65f6a01f948fd3b83ba6fafa3d014c51145dce21e4
cli darwin-arm64              56f8b6277879defae544ee1272fc6e1fbda8ba2fc410f69463a142f5356af966
cli sha256.json               b1ffab59a7f01bb319372512f989887a5eedc810eb49ed523ce178fd54d5e3ed
skill windows-x64.exe         2bb927dfee907777e1bf8e8af369f62b9176f9ff403fc5e896730a8b136cfadf
skill windows-arm64.exe       74088d91b2758ba0fbbdae885b73461b9f28293f4f655d4c224c3bde3b9c2bb1
```

Each hash equals its live GitHub digest. The hosted `sha256.json` content lists all six CLI digests and matches my measurements.

Clean built CLI revision: `go version -m` on the downloaded hosted darwin-arm64 binary reports `vcs.revision=88dc6c5da139d2deb039bda4b605acfbe42663b4`, `vcs.modified=false`, `-trimpath=true`. The binary prints `parley 1.55.0`. The hosted binary records the clean merge revision.

Skill package content: the hosted tarball holds 211 files and `skills/parley-deck/references/STE100.md`. Its `package.json` says 2.19.0. All six packaged `SKILL.md` hashes equal the payload table, and `STE100.md` equals the reviewed source byte for byte.

Skill Windows workflow evidence:

```text
run 38415283057 "Release portable binaries", event release, head f557c5e…, completed/success
  job macos-tests (macos-latest): npm ci, npm test — success
  job windows (runs-on: ubuntu-latest): npm ci, npm test, npm run build:portable:windows, upload — success
```

The workflow file sets `runs-on: ubuntu-latest` for the `windows` job. The Windows assets are an Ubuntu cross-build. This is not native Windows validation. Verdict: CONFIRMED (PRIMARY).

## DC2 — Homebrew

Attempted refutation: I fetched the live tap, downloaded both source archives, and re-hashed the full installed payload.

Tap commit `8ef3b97bcf97b9d710fa5a26305f010bb9c33e2c` exists. Message: `[codex-1] meta-protocol-change-ste100-output-format: update both release formulae`. It changes exactly two formula files, two additions and two deletions each. The live formula URLs and hashes:

```text
parley-deck-cli.rb   url …/parley-deck-cli/archive/refs/tags/v1.55.0.tar.gz   sha256 2c611512ca87a471a815ac5c3dd345eca3ac334b5742c27263e2c8d1e22dff9b
parley-deck-skill.rb url …/parley-deck-skill/archive/refs/tags/v2.19.0.tar.gz sha256 60de2608664d4fb17e6e1dfd6f09638d746cdbb162ffc2eee66915f2016c2349
```

I downloaded both tag archives. My hashes equal both formula values.

Local install: `brew list --versions` reports `parley-deck-cli 1.55.0` and `parley-deck-skill 2.19.0`. The installed CLI prints `parley 1.55.0`.

Recorded exits: `homebrew-verification.json` records exit 0 for both `brew test` runs and `brew audit --strict` (SECONDARY, named record with logs). My own live checks above agree with its versions, tap commit, and bytes.

Payload re-hash: I hashed every file under `/opt/homebrew/opt/parley-deck-skill/libexec`. All 202 files of the `skills/` subtree are byte-identical to the npm package. The only differing file anywhere is `bin/parley-deck-skill.js`. Its first line changed from `#!/usr/bin/env node` to `#!/opt/homebrew/opt/node/bin/node`. The bytes after line one hash equal: `78dc0e3fbe32d874cb47b9c681d1a0ae0af91d1cae8e173a2c52a32f4222123b`. This is the standard Homebrew shebang rewrite of the launcher script. No skill byte differs. The installed prefix also holds repo files such as `packaging/` because the formula builds from the source archive. Extra files are not changed skill bytes. Verdict: CONFIRMED (PRIMARY).

## DC3 — skill WinGet PR

Attempted refutation: I fetched the live PR and its head bytes, then compared them with my own asset hashes.

Live state:

```text
PR 451758 microsoft/winget-pkgs "Update: Feci.ParleyDeckSkill to 2.19.0"
state OPEN, mergedAt null — upstream has NOT merged it; review remains open
changedFiles 3, additions 48, deletions 0
manifests/f/Feci/ParleyDeckSkill/2.19.0/Feci.ParleyDeckSkill.installer.yaml      +15
manifests/f/Feci/ParleyDeckSkill/2.19.0/Feci.ParleyDeckSkill.locale.en-US.yaml   +27
manifests/f/Feci/ParleyDeckSkill/2.19.0/Feci.ParleyDeckSkill.yaml                +6
```

Exactly three added manifests exist. I downloaded each file from the PR head `2e805b796d3e9981daf39ea3377b141c50e7309c`. My hashes equal the local validation record:

```text
installer.yaml      a2665087667dd18710f2f0f0acfc774788e406eb09741940a2b5f82cfc7c0c3a
locale.en-US.yaml   804f80e3278b3133d611396bc8ff573a682f00f552c3af02d1cdf3836ced0b47
version yaml        92c4697347f539d60d443b81dc0c71c6e3bac60a45b724415b8d99c31afa44c7
```

Every line ends with CRLF: 15/15, 27/27, and 6/6 lines. Fields check: `PackageIdentifier: Feci.ParleyDeckSkill`, `PackageVersion: 2.19.0`, `ManifestVersion: 1.12.0`, schemas 1.12.0 for installer, defaultLocale, and version. `InstallerType: portable`, `Commands: parley-deck-skill`.

Both installer URLs point at the final hosted Windows assets. Both `InstallerSha256` values equal the hashes of the bytes I downloaded myself: x64 `2BB927DF…CFADF`, arm64 `74088D91…C2BB1`.

`winget` is not installed on this Mac. Native `winget validate` and `install` are unavailable here. Upstream validation remains open.

CLI WinGet held: the `feci` PR history in winget-pkgs shows no CLI PR after `Feci.ParleyDeckCli to 1.48.0` (2026-09-22). The skill PR references only skill assets. It publishes no CLI. Verdict: CONFIRMED (PRIMARY).

## DC4 — managed runtimes and generic destinations

Attempted refutation: I built my own inventory, hashed every `SKILL.md`, read every receipt, and ran doctor myself.

I inventoried 19 targets: 15 managed — codex, claude, agy (`~/.gemini/config/plugins`), gemini (`~/.gemini/extensions`), hermes, qwen, codebuddy, goose, kimi (`~/.kimi-code`), droid (`~/.factory`), vibe, cursor, opencode, aionrs, zcode — plus the four generic destinations `~/.hermes/profiles/{ldx,librade,testprofile}/skills` and `~/.config/opencode/skills`.

I hashed all 114 `SKILL.md` files myself. Every hash equals its packaged skill value:

```text
parley-bidding       aad7d7d8d43259893626a4593c69c3d90b8bf2a96b90dc723067fa83a04d9540
parley-deck          64d5236baa0aece207becaa9def8e6fba24776f148dd2289281efa0855341df1
parley-design        9ccf8e9d03cc3655fc30c29714b390f79aa4fa90a0ca47125cbdf8d00b25d1bd
parley-design-check  337aafd6ec126bdbbee55c7718f5090a8f35f22b4bff55f45f14d71c48a42b1a
parley-tracker       25b7a6b2226833b1e8291864685cd0343de635f5b989b23799c40c98d5741198
parley-worktrees     91f94c4017edfea19f61ad298131d45cb74c4e56d9661eef1f3ebca2383fc3a9
```

Result: 19 targets × 6 skills = 114 files, zero mismatches. This count matches the expected 114 across 19 targets.

Receipts: each of the 19 `parley-deck` directories carries `.parley-deck-skill-install.json`. Each says `version: 2.19.0`, `source: npm:parley-deck-skill@2.19.0`, the correct target id, `markerSchema: 2`, the five add-ons, and `installedAt` near 2026-10-13T14:42:4xZ. The four generic destinations report `target: generic`. A copied file alone would not produce these per-target receipts.

Doctor: I used the task-local Python 3.14 PATH `/opt/homebrew/opt/python@3.14/libexec/bin` (`python3` → 3.14.8). I changed no global Python configuration; the PATH export lived in my shell only.

```text
parley-deck-skill doctor --target all --json --include-undetected → overall ok: true
  15 targets, all status=valid, detected=true, missing=[] (undetected ones included)
parley-deck-skill doctor --target generic --dest <each of the four skill dirs> → ok: true, valid, version 2.19.0
```

My first generic probe passed the skills root instead of the skill directory and returned `malformed`. That was my invocation error. The documented destination is `…/skills/parley-deck`; the corrected probe passes. Verdict: CONFIRMED (PRIMARY).

## DC5 — npm registry, cores, staged core delta, D6 preview

Attempted refutation: I queried the live registry and installed cores, re-hashed every referenced artifact, and rebuilt core 2.19.0.

Live npm registry (`npm view parley-deck-skill`): latest `2.13.0`; the published list ends `2.12.1, 2.13.0`. Versions 2.15.0 through 2.19.0 are absent. Installed global cores: `~/.parley/protocol/core/` holds `2.10.0` and `2.13.0` only. Both facts match `owner-publication-status.json`. Only the absent 2.15.0–2.19.0 set belongs in the pending owner actions, as the brief bounds it.

Observation, not a defect against the brief: the registry also lacks 2.14.0. The brief scopes the pending commands to 2.15.0–2.19.0. The staged `COOPERATION-2.14.0.md` is mode 444, unlike the others. I added no version to the pending set and removed none.

I re-hashed every referenced local artifact. All ten hashes match the record:

```text
2.15.0 tgz 64c274a8…89b2a   core 0d81fd80…9c09
2.16.0 tgz 80f3172a…7db6   core 8c6b95b6…a63
2.17.0 tgz 40b6eb42…e3ca   core 6073c311…7957
2.18.0 tgz 2d1e7901…eaf0   core 9ead4475…c512
2.19.0 tgz f07102c8…21e4   core 16323e04…fe3b
```

Core delta: I rebuilt core 2.19.0 from `~/.parley/staging/COOPERATION-2.18.0.md` plus the signed payload `source-context/protocol-hunks.json`. My script required exactly one match of the old rule 6 and no existing Appendix B. The result equals `COOPERATION-2.19.0.md` byte for byte. My hash: `16323e04ffb076db3b2019169010033c6668199ebab0ca6896e79924e74bfe3b`. The new rule 6 and the appendix each occur exactly once. Core 2.19.0 equals staged 2.18.0 plus only H1 and H2.

Owner-note order: the draft lists the D6 `inspect` and `apply` commands first. It then lists `npm publish` and `parley protocol publish` per version, in order 2.15.0, 2.16.0, 2.17.0, 2.18.0, 2.19.0. I executed none of them.

D6 preview, read-only: `d6-preview.json` records 25 roots, `decisions: []`, `history: unknown-history`, manifest SHA `7283348e…`, preview SHA `d78dfd90…`. I counted the 25 roots myself. The run directory `parley-deck/runs/20260510T194003Z` holds `events.jsonl` only, so I cannot recompute the manifest digest from it alone. I searched `~/.parley` for the decision id `driver-unstall-20260510-smoke-v1` and found no applied declaration. The decision remains unapplied. I emulated no writer-stop claim and no attended terminal.

Windows material limits, first-party: runs `38406035407` and `38406028243` at `99d04b3` are `completed/failure`. Both Linux and macOS jobs passed at the same commit, as the goal check records. CLI Windows stays experimental. CLI WinGet stays held. The skill WinGet PR remains open upstream. npm and core publication stays pending. The ASD summary is a working summary, not full certification. Verdict: CONFIRMED (PRIMARY).

## DC6 — draft handoff completeness

The draft `release-handoff-draft.md` contains all four required elements. It carries a channel-evidence table with merge commits and asset hashes. It carries the ordered owner commands: D6 first, then npm and core per version. It carries deferred items. It carries usage limits: measured invocation counts only, native token costs marked unavailable, cumulative counters marked not marginal.

I cross-checked its numeric claims against my own measurements: both merge commits, the tarball hash, the core 2.19.0 hash and delta claim, the 114-file count, the 15+4 target split, npm latest 2.13.0, installed cores, the 25-root preview, and the workflow evidence. Every checked claim matches. The final canonical inbox note and the verified shared-memory write follow my result. I claim neither as complete. Verdict: CONFIRMED (PRIMARY).

## Findings and observations

No mismatch defeats a delivery criterion. Material observations:

- The npm registry also lacks 2.14.0. The brief scopes pending actions to 2.15.0–2.19.0, and the draft complies. I report the fact without widening the pending set.
- Homebrew rewrote the launcher shebang in `bin/parley-deck-skill.js`. No skill payload byte differs. The check criterion holds.
- The formula prefix holds extra repo files from the source archive build. All 202 skill files match the package.
- The skill WinGet PR remains open. Upstream controls the merge.
- My first generic doctor probe used the wrong path and failed. The corrected, documented path passes. The runtime state is valid.

## Coverage limits

- I ran on macOS only. Native Windows validation does not exist in this evidence. The skill Windows assets are an Ubuntu cross-build. CLI Windows remains experimental.
- `winget` is unavailable on this Mac. No native winget validate or install claim is made.
- The recorded `brew test` and `brew audit` exits are SECONDARY evidence from the named record. I verified the tap, formulas, installed versions, and payload bytes myself.
- The D6 manifest digest is not recomputable from the run directory alone. The preview record and the absent decision are my evidence.
- File delivery proves transport, not that a model read a file.
- Configured model identity is not native observed-model proof.
- The staged core was verified by hash and reconstruction. No release channel carried it yet; its publication stays an owner command.
- Doctor checks needed the task-local Python 3.14 PATH. The system Python 3.9.6 cannot run the bidding add-on checks.

## Commands run

```text
shasum -a 256 <private prompt file>
gh pr view 84 --repo feci/parley-deck-cli / pr 12 --repo feci/parley-deck-skill --json …
gh api repos/…/git/ref/tags/v1.55.0 / v2.19.0 ; gh api repos/…/git/commits/<merge> ; branches/main
gh api repos/feci/parley-deck-cli/releases/tags/v1.55.0 / parley-deck-skill/releases/tags/v2.19.0
curl -L -o <skill tgz, cli darwin-arm64, sha256.json, both skill windows exe, both formula source archives> ; shasum -a 256
go version -m <downloaded cli binary> ; <binary> --version
tar -tzf / -xzf <skill tgz> ; shasum -a 256 <six packaged SKILL.md, STE100.md>
gh run view 38415283057 --repo feci/parley-deck-skill --json jobs ; gh api …/contents/.github/workflows/release-portable.yml
gh api repos/feci/homebrew-parley/commits/8ef3b97… / contents/Formula/<both>.rb
brew list --versions feci/parley/parley-deck-cli feci/parley/parley-deck-skill
python3 <full libexec payload re-hash vs package; skills subtree compare>
gh pr view 451758 --repo microsoft/winget-pkgs --json … ; gh search prs --repo microsoft/winget-pkgs --author feci
curl raw.githubusercontent.com/winget-pkgs/2e805b7…/<three manifests> ; shasum ; grep -c $'\r'
python3 <114 SKILL.md inventory over 19 targets; receipt reader>
PATH=/opt/homebrew/opt/python@3.14/libexec/bin:$PATH parley-deck-skill doctor --target all --json --include-undetected
PATH=… parley-deck-skill doctor --target generic --dest <each documented skill dir> --json
npm view parley-deck-skill versions --json ; npm view parley-deck-skill version
python3 <re-hash five tarballs and five staged cores; reconstruct core 2.19.0 from 2.18.0 + H1/H2>
grep -rl "driver-unstall-20260510-smoke-v1" ~/.parley ; find ~/.parley -name "*legacy*" -newer <staged 2.18.0>
git -C ste100 log/diff --name-only 99d04b3..HEAD ; git -C ste100-skill rev-parse HEAD ; shasum <SKILL.md, STE100.md>
gh run view 38406035407 / 38406028243 --repo feci/parley-deck-cli --json headSha,status,conclusion
```

CHANNEL-CHECK: PASS
