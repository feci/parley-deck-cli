---
idea: meta-protocol-change-quota-auto-exclude
artifact: release-channel-verification
author: claude-1
model: claude-opus-5-5[1m] (plain id), max effort
date: 2026-10-08
role: separate independent release-channel verifier (not a review cycle, not a consensus signoff)
---

# Release-channel verification — claude-1

## Disposition

**Agent-controlled delivery is complete. No channel blocker needs action.**

- Channel findings: 0 CRITICAL, 0 MAJOR, 0 MINOR. There are 2 NITs, both pre-existing and outside this release's diff, plus several stated limits.
- Published and verified:
  - CLI v1.51.0 and skill v2.15.0 (main, tags, GitHub releases, all 11 assets).
  - Both Homebrew formulae.
  - All 15 local runtime installations.
- Pending upstream: the skill-only WinGet PR. 7 of 10 validations plus the CLA pass; 3 are still running or queued.
- Pending owner-only: npm publish and global core publish. Both are staged and verified, and neither has been executed.
- Pending organizer step: the single final `…_released.md` inbox note does not exist yet.
- CI is **not all-platform green**. CLI release-commit run 37706766781 concluded `failure`: Linux and macOS passed, Windows failed.

This snapshot was taken between 2026-10-08 00:57Z and 01:10Z. GitHub releases are mutable (they are not GitHub-immutable releases), so these results hold only at that time.

## Context and attestation

- **Invocation.** This is `channel-verification-resume10-timeout-1` (2400 s ceiling). The first attempt timed out at 1800 s without a terminal API error.
  - Raw downloads left by that attempt were re-hashed before any use.
  - Every verdict below comes from this process's own commands. Fresh outputs are under `.parley-runtime/claude1-release-channels-resume10/r2/`.
- **Packet.** Command: `parley protocol packet --phase 8 --track deliberation --flag protocol_change --idea meta-protocol-change-quota-auto-exclude --json`.

  | Field | Value |
  |---|---|
  | `context_mode` | `full` |
  | `source_sha256` | `73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e` (matches the expected value) |
  | `packet_sha256` | `73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e` |
  | `fallback_reason` | absent / none. Full context is the default mode, not a fallback. |
  | Transport and role | `github-pr`; source role `source` (live source file) |
  | Body | `.parley-runtime/protocol-packets/full-phase8-deliberation-73613f95….md`, 125,862 bytes |

  - The emitted full body was read from that path, re-hashed, and is equal to `parley-deck/COOPERATION.md`.
  - The shadow audit packet (`38ebdb5e…`, 40 blocks included, 29 omitted) was not used.
- **Binding inputs read.**
  - The finish-now note (points 3, 5, 6, 7) and the round08-answer note.
  - Review consensus, via `parley consensus status --review --json`: triage `ready`, codex-1 ✅ ACCEPT, claude-1 ✅ ACCEPT.
  - IMPLEMENTATION's "Attended implementation close — 2026-10-08": status complete; AC2 NOT MET / owner-waived; AC5 and AC15 NOT MET in full / accepted-deferred; no cycle 6.
- **Closed-artifact hashes.** These are unchanged and equal to HEAD:
  - IMPLEMENTATION `fc9bbbc6…`
  - `review/consensus.md` `4f9fde28…`
  - `review/round-08/claude-1.md` `6ea87126…`
- **Go suite.** Not re-run, because no channel-related reason required it.

## 1. Source refs and GitHub releases — PUBLISHED, VERIFIED

**Commands:** `git ls-remote` on all three repos, `gh api …/git/ref/tags/*`, `gh release view`, `gh api …/releases/tags/*` (uploader and timestamps), a fresh `gh release download` of all assets with `shasum -a 256`, `go version -m`, and `file`.

**Refs**

| Repo | `main` / HEAD | Tag object → commit |
|---|---|---|
| parley-deck-cli | `343aab21a7d95083767260c01c66d5e69f557cd6` | v1.51.0 `cbb41ba5…` → `343aab21…` |
| parley-deck-skill | `352a475c4697d43e7582d00cc789bb6462867ee5` | v2.15.0 `eb410244…` → `352a475c…` |

**Releases.** CLI v1.51.0 was published at 00:17:09Z and skill v2.15.0 at 00:17:15Z. Neither is a draft or prerelease. Each is its repo's `releases/latest`.

**Assets.** All 11 match: fresh-download SHA256, GitHub server digest, delivery copy and `channel-evidence.json` agree for every asset. `sha256.json` (`2aeafa8d…`) lists all six CLI binaries with matching hashes.

| Asset | SHA256 |
|---|---|
| parley-v1.51.0-darwin-arm64 | e8d15ea0bcc5a3433464098844e86d169415610b03ea9c648837a221bc8fe829 |
| parley-v1.51.0-darwin-x64 | fc7d07ff17437d6bbb8012a2baf73baa89ddf2a992f2d75daa0e76574e6a9473 |
| parley-v1.51.0-linux-arm64 | 4230dd163c9b6a4c16143908a79e353889c6253ff7e77bed3325c39e6234ac04 |
| parley-v1.51.0-linux-x64 | 36e733904b4a09405cf47f6fee474aad41f6da2120c2aa091a7deb687448413f |
| parley-v1.51.0-windows-arm64.exe | b949b0f10eeb7a5fbc4f2a06a03c61ed67805abcc1baed09c2ab2727f4f524c1 |
| parley-v1.51.0-windows-x64.exe | c92c51465ebbb404c046abe00f3ee45d0b7e126b0031ce3852157f1ddf898611 |
| sha256.json | 2aeafa8d9386296fbc2a38bdb9c22636ffa7cbdbf175f4059fd77d5381ea97ef |
| parley-deck-skill-2.15.0.tgz | 64c274a86d0b4f2e136decbb1cc2810a18de89191192c6b9e680c5ec9af89b2a |
| parley-deck-skill-v2.15.0-macos-arm64 | 3ee27e552e695cb84be06be4e209c5419f7aaa5c173e7de1691948dcb64750df |
| parley-deck-skill-v2.15.0-windows-arm64.exe | 5e158ee8cb8da4fba2f4288a477325d0bc211b62aa30f918723faba60f2235a6 |
| parley-deck-skill-v2.15.0-windows-x64.exe | f02cdf6903f56a68420f0a141f48c6eccb03a07aab8379c22ce792435b6f662d |

**CLI provenance.** `go version -m` reports the same build settings for all six binaries:
- go1.27.1, path `parley-deck-cli/cmd/parley`, module version v1.51.0
- `vcs.revision=343aab21…`, `vcs.time=2026-10-08T00:14:30Z`, **`vcs.modified=false`**
- `-trimpath=true`, `CGO_ENABLED=0`, and GOOS/GOARCH matching each asset name

`-ldflags` is not recorded because Go omits it under `-trimpath`. This agrees with the argv in `build-results.json`. The binaries were built locally and uploaded by `feci` at 00:17:05–08Z; the CLI repo has no release workflow, only `tests.yml`.

**Skill Windows EXE provenance.**
- Both EXEs were uploaded by **`github-actions[bot]`**. They were created at 00:18:51Z and last updated at 00:18:54–55Z.
- That is inside run 37706962158's "Upload release assets" step (00:18:50Z, `gh release upload "$TAG" dist/* --clobber`). The run concluded `success` at 00:18:58Z.
- Neither file has been modified since. Their digests equal the WinGet `InstallerSha256` values. They were never replaced with local bytes.
- The tgz and the macOS portable were uploaded by `feci` (local build) at 00:17:10–14Z.

**Formats** (checked with `file`):
- CLI: Mach-O arm64 and x86_64, statically linked ELF aarch64 and x86-64, PE32+ console Aarch64 and x86-64.
- Skill: Mach-O arm64, PE32+ console Aarch64 and x86-64.

**Native smoke tests** (isolated `HOME` under a temp directory):
- CLI darwin-arm64: `version` printed `parley 1.51.0`. `init --dir` created a deck, and `protocol status` ran.
- CLI darwin-x64: `bad CPU type in executable` because Rosetta is not installed, so it was not executed.
- Skill macOS portable:
  - `--version` printed `2.15.0`.
  - `install --target generic --dest <tmp>` exited 0, and `doctor --json` returned ok.
  - Core plus 5 add-ons are valid at 2.15.0, and SKILL.md is `76765df8…`.
  - Each add-on's manifest aggregate equals the aggregate printed by the CI `npm test` run.

**npm tarball vs tag source.** All 210 tarball files are byte-identical to the GitHub source archive of v2.15.0 (archive SHA256 `08b17760…`).

**CI**

| Run | Trigger | Result |
|---|---|---|
| CLI 37706766781 | push main @343aab2 | **failure**: ubuntu success, macOS success (00:32:10Z), **windows failure** |
| Skill 37706769960 | push main | success (Python 3.10 and 3.13 jobs) |
| Skill portable 37706962158 | release v2.15.0 | success (`macos-tests` and `windows` jobs) |

**Limits**
- No Linux or Windows binary was executed, and darwin-x64 was not executed.
- The macOS skill portable is a local build. Its provenance was checked by what it installs (identical payload output), not by a reproducible build.

## 2. Homebrew — PUBLISHED, VERIFIED

**Commit**
- `feci/homebrew-parley` main = HEAD = `320a1d70ee59b9eb47a1c37d2f38360eba313809`, with single parent `9855366a…`.
- It touches exactly `Formula/parley-deck-cli.rb` and `Formula/parley-deck-skill.rb`. Remote blobs `31554d47…` and `734c4e25…` equal the local diff.
- The diff changes only the `url` and `sha256` lines. Install logic, the `post_install_steps` shebang restoration and both `test` blocks are unchanged.

**Local copies.** The source checkout and the installed tap are both at `320a1d70…` and clean. The formula files hash to `966b1d49…` (CLI) and `2d0e7b60…` (skill) in both places.

**Formula sha256 vs fresh GitHub archives**
- `fabbc864…` = the v1.51.0 archive
- `08b17760…` = the v2.15.0 archive

**Installed state**
- Kegs: parley-deck-cli 1.50.0 and 1.51.0; parley-deck-skill 2.14.0 and 2.15.0.
- Links point to 1.51.0 and 2.15.0. `INSTALL_RECEIPT.json` names tap `feci/parley` for both.
- `parley --version` prints 1.51.0, and `parley-deck-skill --version` prints 2.15.0.

**My rechecks** (`HOMEBREW_NO_AUTO_UPDATE=1`, at about 01:08Z):
- `brew style`: 2 files, no offenses.
- `brew audit --strict --online feci/parley/…`: exit 0, no problems.
- `brew test` on the CLI formula (version and help): exit 0.
- `brew test` on the skill formula (version, install, doctor): exit 0.

**Skill installer operation and manifest integrity** (Homebrew binary, isolated `HOME`, `install --target codex --yes`, then `doctor --target codex --json`):
- Result is ok. Core and all 5 add-ons are valid, managed and at 2.15.0.
- Add-on manifest aggregates equal the CI `npm test` aggregates:

  | Add-on | Aggregate |
  |---|---|
  | bidding | `7854adf1…` |
  | design | `fca17c70…` |
  | design-check | `2c65c148…` |
  | tracker | `07d98263…` |
  | worktrees | `ca446d59…` |

- SKILL.md is `76765df8…`.
- No script under `libexec/skills` has a shebang pointing into a Homebrew prefix, so the restoration works. For example, `parley-tracker/bin/claim.js` starts with `#!/usr/bin/env node`.

See NIT-1 and NIT-2 below.

## 3. WinGet (skill only) — PR OPEN, UPSTREAM PENDING

**The PR**
- https://github.com/microsoft/winget-pkgs/pull/448514, "Update: Feci.ParleyDeckSkill to 2.15.0".
- OPEN and not a draft. Head `feci:feci-skill-2.15.0` @ `cf64f71bb02761b28c722c180e5085b51d64b769`, base `master`.
- One commit on parent `8f767551…`, adding exactly three files under `manifests/f/Feci/ParleyDeckSkill/2.15.0/`: installer (15 lines), locale en-US (27 lines), version (6 lines).

**Manifest content**
- Git blob ids at the PR head equal the local sparse clone: `c4e76884…`, `9edc996e…`, `6d2ba117…`.
- Content is equal after CRLF normalization. The CRLF comes from the local checkout; the repo uses `*.yaml text=auto`.
- `ManifestVersion` and `$schema` are 1.12.0 in all three files. `PackageVersion` is 2.15.0, `InstallerType` is portable, and `Commands` is `parley-deck-skill`.
- The x64 and arm64 `InstallerUrl` values point at the v2.15.0 release assets.
- The `InstallerSha256` values (`F02CDF69…`, `5E158EE8…`) equal the live digests of the workflow-built EXEs.

**Upstream checks at 01:08Z**

| Check | State |
|---|---|
| 01 Pull Request Validation | SUCCESS |
| 02 Manifest Validation | SUCCESS |
| 03 URLs Validation | SUCCESS |
| 04 URL Domain Validation | SUCCESS |
| 05 Manifest Policy Validation | SUCCESS |
| 06 Catalog Content Verification | SUCCESS |
| 07 Installers Scan | SUCCESS |
| license/cla | SUCCESS |
| **08 Installation Validation** | **IN_PROGRESS** |
| **09 Installer Metadata Validation** | **QUEUED** |
| **10 Validation Completed** | **QUEUED** |

- The agentic-workflow jobs show SKIPPED.
- `mergeStateStatus` is BLOCKED and `reviewDecision` is REVIEW_REQUIRED. The only label is `New-Manifest`, and there are no reviews.

**Context and the CLI hold**
- Upstream master's latest Feci.ParleyDeckSkill is 2.14.0, merged in #441054 through the same pipeline.
- Feci.ParleyDeckCli's latest upstream is 1.48.0. No CLI manifest for 1.49 or later, including 1.51.0, has been submitted, so the CLI WinGet hold is respected.

**Limits**
- No local `winget validate` or install was run.
- The skill EXEs have never been executed on Windows. The portable workflow's `windows` job runs on **ubuntu-latest**: it cross-builds with `pkg`, and its `npm test` runs on Linux.
- The only native Windows install check is upstream step 08, which is still running.
- I did not re-run the local schema validation. I relied on upstream step 02 passing and on inspecting the fields myself.

## 4. Runtime installations — 15/15 VERIFIED

**Doctor.** `parley-deck-skill doctor --target all --include-undetected --json` (installed 2.15.0) returned ok with 15 targets, all `valid`.
- Every core marker shows version 2.15.0 and source `npm:parley-deck-skill@2.15.0`.
- Each target has 5 of 5 add-ons valid at 2.15.0.

**SKILL.md hashes.** I hashed each file myself: 15 of 15 equal `76765df83428a23288d5d1d09f0fc7b40d53197a832c0f3f51e11a3a3687a599`.

| Target | Path |
|---|---|
| codex | `~/.codex/skills/parley-deck` |
| claude | `~/.claude/skills/parley-deck` |
| agy | `~/.gemini/config/plugins/parley-deck` |
| gemini | `~/.gemini/extensions/parley-deck` |
| hermes | `~/.hermes/skills/parley-deck` |
| qwen | `~/.qwen/skills/parley-deck` |
| codebuddy | `~/.codebuddy/skills/parley-deck` |
| goose | `~/.goose/skills/parley-deck` |
| kimi | `~/.kimi-code/skills/parley-deck` |
| droid | `~/.factory/skills/parley-deck` |
| vibe | `~/.vibe/skills/parley-deck` |
| cursor | `~/.cursor/skills/parley-deck` |
| opencode | `~/.opencode/skills/parley-deck` |
| aionrs | `~/.aionrs/skills/parley-deck` |
| zcode | `~/.zcode/skills/parley-deck` |

**Byte-level comparison** (install marker excluded):
- **All 75 add-on trees are byte-identical to the npm tarball payload.**
- **All 15 core trees are byte-identical to fresh isolated installs made by the released 2.15.0 portable binary.** The reference was a generic-layout install for 13 targets, and target-specific fresh installs for agy (14 files) and gemini (13 files).
- The bundled `references/COOPERATION.md` hashes to `73613f95…`. That equals CLI v1.51.0's `parley-deck/COOPERATION.md` and the skill tag's reference copy.

**Scope note.** Only the end state was verified. The initial `--target all` run reached 8 detected runtimes. The later all-15 `--include-undetected` run is recorded in the delivery `runtime-verification.json` and in `logs/runtime-install-all15.json`.

## 5. npm — STAGED, OWNER-ONLY, NOT PUBLISHED

**Registry**
- `dist-tags.latest` is **2.13.0**.
- The version list lacks 2.14.0 and 2.15.0 (and 2.11.0).
- `npm view parley-deck-skill@2.15.0` returns E404.

**Tarball**
- SHA256 `64c274a8…`, equal to the GitHub asset.
- sha1 `4c247f42…` and integrity `sha512-ZcHrXsd24…FCepg==`, both equal to `package-evidence.json`.
- `package.json` has name `parley-deck-skill` and version 2.15.0. All 210 files match the tag archive.

**Dry run.** `npm publish --dry-run --access public '<tgz>'` from an isolated cwd exited 0. It printed "Publishing to https://registry.npmjs.org/ with tag latest and public access (dry-run)", and its name, version, shasum, integrity and 210-file count all match. Nothing was written to the registry.

**Owner command.** The path exists, and the quoting is correct for a path with spaces:
`! npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude/skill/parley-deck-skill-2.15.0.tgz' --access public`

**Effect until the owner publishes.** The skill's help text recommends `npx -y parley-deck-skill@latest install`, which currently installs 2.13.0. The skill release notes say npm is staged, so this is expected.

## 6. Global core — STAGED, OWNER-ONLY, NOT PUBLISHED

**Current state**
- `~/.parley/staging/COOPERATION-2.15.0.md` has SHA256 `0d81fd807114e4ff67f5fa98bba73622b58c86d66a2bf6b57d1def68cd099c09`. That equals `core-staging.json` and the preview.
- `parley protocol status` shows only 2.10.0 and 2.13.0 installed in the core store.

**Reproduction.** In a temp git repo, I applied `core-reviewed.patch` (`35072b16…`) to staged 2.14.0 (`51476d69…`) with `git apply`. It applied cleanly, with no fuzz. The result is **byte-identical** to `0d81fd80…`.

**Patch provenance**
- The hunks are identical to `git diff 35a6a3d 94f16aa`. Those blobs are `internal/protocol/defaults/COOPERATION.md` at the v1.50.0 and v1.51.0 tags.
- Staged 2.14.0 and 2.15.0 differ by 117 lines (109 added, 8 removed) and nothing else.

**Base chain.** The unpublished 2.14.0 base is a faithful template:
- Staged 2.13.0 equals published core 2.13.0 (`fc907e59…`).
- Published 2.13.0 plus the v1.49.1→v1.50.0 embedded-default diff is byte-identical to staged 2.14.0.

**Template zones**
- Staged 2.15.0 differs from v1.51.0's embedded default only on lines 5–6 (`Transport:` and `Created:`).
- The generic placeholders `<workspace-name>`, `<transport-choice>` and `<YYYY-MM-DD>` are preserved, on the same three lines as in 2.14.0.

**Owner command.** It was not executed. Its syntax matches CLI 1.51.0's usage (`publish --version V --from FILE`, attended, needs a TTY):
`! parley protocol publish --version 2.15.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md`

## 7. Release notes and changelogs

**Bodies.** The live GitHub release bodies are byte-equal to the delivery `cli-release-notes.md` and `skill-release-notes.md`.

**R8 disclosures.** The three disclosures from `source-context/codex-1-round08-disclosure-proposal.md` (R8-MINOR-1, R8-NIT-1, R8-NIT-2) appear verbatim, comparing with whitespace normalized, in four places:
- the CLI release body
- the skill release body
- `CHANGELOG.md` at v1.51.0 (section `## 1.51.0 — 2026-10-08`)
- `CHANGELOG.md` at v2.15.0 (section `## 2.15.0 — 2026-10-08`)

**Follow-up links** (both return HTTP 200):
- `quota-kickoff-reporting-and-alias-guidance` (`blob/v1.51.0/…/00-prompt.md`), linked from both bodies and both changelogs.
- The AC2 follow-up, `quota-zcode-native-exhaustion-capture`.

**CLI release body states:**
- AC2 is NOT MET and owner-waived (round05-answer Q2); R5-MAJOR-2 is deferred, not fixed.
- Windows CLI is known broken for new idea creation: directory sync fails with "Access is denied". Scoped driving and signing, imports, revisions and transitions are also affected.
- The Windows test suite fails, including the quota directory-sync errors and the Unix-only `syscall.Mkfifo` test.
- Windows assets are experimental, and CLI WinGet is held.
- **D6** (legacy run accounting) is a separate follow-up.
- AC5 and AC15 are unmet, with the explicit exceptions.

**Skill release body states:**
- The AC2 waiver and its follow-up.
- Windows is broken, and "Windows CLI CI is failing, not passing or merely untested".
- CLI WinGet is held.
- The three R8 disclosures and the follow-up link.
- npm and core are staged.
- D6 appears only in the CLI notes, which suits a CLI-scoped item.

**CI claims.** Neither body claims all-platform green. "Linux and macOS candidate CI … pass" matches the release-commit run.

**Windows CI evidence**
- Job 113082869607's log contains `undefined: syscall.Mkfifo` (`internal/evidence/tree_report_test.go:97`) and 142 "Access is denied" lines, including the quota catch-up `sync` failures.
- Totals: 18 packages FAIL, 15 ok, 159 top-level failing tests.
- Windows CI was already failing at v1.50.0 (run 36110596716: 14 packages FAIL, 16 ok), so Windows' experimental status predates this release.
- The disclosure says the suite fails "including" those causes. That is accurate, though not a complete list.

## Findings

No CRITICAL, MAJOR or MINOR channel findings. No channel blocker needs action.

- **NIT-1: stale duplicate formula names (local machine, pre-existing; the published channel is fine).**
  - A local tap that is not a git repo, `feci/release-validation-20260918`, still carries `parley-deck-cli` 1.48.0 and `parley-deck-skill` 2.12.1.
  - Because the names collide, the unqualified `brew info parley-deck-cli` fails with "Formulae found in multiple taps".
  - The installed kegs come from `feci/parley`, and fully qualified commands work.
  - The owner or organizer may remove it with `brew untap feci/release-validation-20260918`. I did not change it.
- **NIT-2: Homebrew installs omit LICENSE and README.md from the core skill (packaging, pre-existing since 2.14.0 or earlier; not a regression).**
  - Homebrew moves LICENSE, README.md, NOTICE.md and CHANGELOG.md to the keg root.
  - As a result, `parley-deck-skill install` run from a Homebrew keg puts the core skill in place without LICENSE and README.md. I confirmed the same behavior with the 2.14.0 keg.
  - npm and portable installs do include both files, and doctor still reports `valid`.
  - This is a candidate for a later formula change.
- **INFO**
  - The Windows skill EXEs are cross-built on Linux. Their only native Windows execution is upstream WinGet step 08, which is still running.
  - The npm registry never received 2.14.0, and the default `npx …@latest` command installs 2.13.0 until the owner publishes.
  - The old kegs (CLI 1.50.0, skill 2.14.0) are still installed. That is harmless; running `brew cleanup` is the owner's choice.

## Pending operations (not done by me)

**Owner-only**
1. `! npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-08-quota-auto-exclude/skill/parley-deck-skill-2.15.0.tgz' --access public` (needs an OTP).
2. `! parley protocol publish --version 2.15.0 --from /Users/tomasfecko/.parley/staging/COOPERATION-2.15.0.md` (attended, needs a TTY).

**Upstream**

3. microsoft/winget-pkgs#448514: step 08 is in progress and steps 09–10 are queued. Moderator approval and the merge/publish pipeline follow. No agent action is needed unless a validation fails.

**Organizer (agent-controlled, still open)**

4. Write the single final note `parley-deck/inbox/codex-1-to-user_meta-protocol-change-quota-auto-exclude_released.md` with the two exact owner commands. Finish-now point 6 requires it, and it did not exist at 01:08Z.

**Held by design**
- CLI WinGet.
- Windows CLI stays experimental.

## Side effects and scratch

- **Scratch locations**
  - `.parley-runtime/claude1-release-channels-resume10/r2/`: fresh downloads, extracted archives, logs and JSON.
  - Temp directories with an isolated `HOME`: `/tmp/c1smoke.*`, `c1skill.*`, `c1brew.*`, `c1brew214.*`, `c1core.*`, `c1chain.*`, `c1npm.*`, `c1tgt.*`.
- **One unintended side effect, reverted**
  - My first Homebrew recheck set `HOMEBREW_NO_INSTALL_FROM_API=1`, and the audit started cloning `homebrew/core`.
  - I stopped my own process. git removed the partial clone.
  - I then used `rmdir` to remove the empty `/opt/homebrew/Library/Taps/homebrew` directory it had created (born 01:06:48Z, empty).
  - The `brew tap` list is back to its original four entries: altfins-com/tap, anomalyco/tap, feci/parley, feci/release-validation-20260918.
  - The recheck re-ran in API mode and passed.
- **Not done:** no commits or pushes, no publication, no installs into real runtime paths, no edits to the roster, settings, release assets, IMPLEMENTATION or reviews, and no external comments. This report is not a consensus signoff.
