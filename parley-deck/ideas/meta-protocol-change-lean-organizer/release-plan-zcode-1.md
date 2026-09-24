---
author: zcode-1
idea: meta-protocol-change-lean-organizer
date: 2026-09-24
phase: release-plan (auxiliary; corrected at dispatch — see "Corrections at dispatch" below)
status: corrected-at-dispatch-2026-09-24
---

# Release plan — CLI 1.49.0 / skill 2.13.0 (metadata prep by zcode-1; release OPERATIONS by the organizer)

**Roles (corrected 2026-09-24 at dispatch, post-Phase-8 close):** zcode-1 (implementer,
participant) executes ONLY the release metadata preparation — §2 version/CHANGELOG bumps,
manifest consistency, and participant-side preparation checks — in the paired
lean-organizer worktrees, committing each worktree separately. The organizer (codex-1)
performs the release operations — §3–§9 integration, tags, builds, GitHub releases, npm,
Homebrew, winget — after zcode-1's report; the organizer does not verify code (verdicts
already sit with the reviewers and the commissioned LE-7 check). The owner ALONE performs
the attended core publish (`parley protocol publish`); it is owner-only, currently
PENDING, and is NOT a prerequisite for releasing the channels below.

## Corrections at dispatch (2026-09-24, post-Phase-8 close — made by zcode-1)

1. **Roles/header:** the original header said "executed by zcode-1 on organizer
   dispatch". Corrected as above: zcode-1 prepares metadata only; the organizer
   publishes; the owner alone attends the core publish.
2. **§1.2 core-publish framing:** removed "the release is reported complete only when
   both are done" — channel release does not wait on the owner's core publish. The
   publish is reported honestly as a pending owner-only action.
3. **§2a test timeout:** `go test ./... -count=1` → `-timeout 45m` added (the reviewed
   G3 convention; the slowest package measured 626.9 s locally, over Go's 600 s default).
4. **§13 gate (4):** winget completion no longer gated on external PR merge — open PRs
   are reported honestly as open; catalog merge latency is outside our control.
5. **Inventory freshness:** §0 was verified mid-review at CLI HEAD `3c97f44`; both
   branches have since advanced (fix-up cycles + records; CLI now `3d57c59`, skill
   `b06a65a`) and the fast-forward claim was re-verified at dispatch for both repos.

Everything below was verified read-only on 2026-09-24 against the live repositories, the winget
catalog, npm, and the local checkouts. No product code, version file, branch ref, release, or
external PR was touched while preparing this plan; no commits were made.

## 0. Verified inventory (evidence, not assumptions)

| Item | Verified value | Evidence |
| --- | --- | --- |
| CLI worktree / branch | `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer`, branch `lean-organizer`, HEAD `3c97f44` | `git worktree list`, `git log` |
| CLI remote | `https://github.com/feci/parley-deck-cli.git` (`origin`) | `git remote -v` |
| CLI main | `origin/main` = `b4831d6` "[codex-1] release: finish evidence-first CLI 1.48.0"; `lean-organizer` = `origin/main` + 7 idea commits → **fast-forward integration possible, history preserved** | `git merge-base --is-ancestor origin/main HEAD` → true |
| Skill worktree / branch | `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill`, branch `lean-organizer`, HEAD `0f513f6` = `origin/main` (`d1e57d5`, "Release 2.12.1") + 1 commit → fast-forward | same |
| Skill remote | `https://github.com/feci/parley-deck-skill` (`origin`) | `git remote -v` |
| CLI version locations | `VERSION` (`1.48.0`), `internal/app/version.go` (`const version = "1.48.0"`), `CHANGELOG.md` top section | `git grep -n "1\.48\.0"` |
| Skill version locations | `package.json` (`2.12.1`), `package-lock.json` (2 entries), `skills/parley-deck/references/compatibility.json` (`skillVersion`), `CHANGELOG.md` | `git grep -n "2\.12\.1"` |
| CLI build entrypoint | `./cmd/parley`; formula build line: `go build -trimpath -ldflags "-s -w" -o bin/parley ./cmd/parley` | `ls cmd/`, tap formula |
| CLI release assets (convention) | 6 binaries on tag `v1.48.0`: `parley-v1.48.0-{darwin-arm64,darwin-x64,linux-arm64,linux-x64,windows-arm64.exe,windows-x64.exe}`; title `Parley Deck CLI 1.48.0`; body mirrors CHANGELOG | `gh release view v1.48.0` |
| Skill release assets | 5 portable binaries `parley-deck-skill-v2.12.1-{macos-arm64,macos-x64,linux-x64,windows-x64.exe,windows-arm64.exe}`; title `Parley Deck Skill 2.12.1` | `gh release view v2.12.1` |
| Skill portable build | `npm run build:portable` = `node scripts/build-portable.js all` → builds **all 5** locally via `@yao-pkg/pkg`; output `dist/parley-deck-skill-v<ver>-<suffix>` | `scripts/build-portable.js`, `package.json` scripts |
| Homebrew tap | `https://github.com/feci/homebrew-parley.git`; local checkout `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley`, on `main`, **1 commit behind origin** (local at 1.47.0/2.11.0, remote `f20b78f` at 1.48.0/2.12.1) → `git pull --ff-only` first | `git status -sb`, `git log origin/main` |
| BOTH formula paths | `Formula/parley-deck-cli.rb` (Go source build from release tarball) and `Formula/parley-deck-skill.rb` (node, `libexec` install + bin symlink + post_install shebang restore) — formulae exist **only** under `Formula/` | `ls Formula/`, `git show origin/main:Formula/*.rb` |
| Winget application IDs | `Feci.ParleyDeckCli` (catalog latest 1.48.0) and `Feci.ParleyDeckSkill` (catalog latest 2.12.1) in `microsoft/winget-pkgs` under `manifests/f/Feci/<App>/<ver>/` — 3 YAMLs each: installer, locale.en-US, version | `gh api repos/microsoft/winget-pkgs/contents/manifests/f/Feci` |
| Winget fork | `feci/winget-pkgs` (fork, default `master`, last pushed 2026-09-18) | `gh api repos/feci/winget-pkgs` |
| Approved two-PR pattern | PR #438660 branch `release/parley-deck-cli-1.48.0` "Update: Feci.ParleyDeckCli to 1.48.0" + PR #438661 branch `release/parley-deck-skill-2.12.1` "Update: Feci.ParleyDeckSkill to 2.12.1" (both MERGED); one application per PR | `gh pr list --author feci` |
| npm package | `parley-deck-skill`, dist latest `2.12.1`; identity `npm whoami` → `tomasfecko` (verified read-only today; token never printed) | `npm view`, `npm whoami` |
| Skill release automation | `.github/workflows/release-portable.yml`: on **release published** → macOS job `npm test` → ubuntu job `npm test` + `build:portable:windows` + `gh release upload "$TAG" dist/* --clobber` (CI replaces the Windows exes — **take winget hashes from the FINAL release assets only**) | workflow file |
| Skill prepack gate | `prepack` = `npm run manifest:check` (refuses stale add-on manifests); `npm test` needs python3 ≥ 3.10 | `package.json`, `RELEASING.md` |
| Skill repo guide | `RELEASING.md` (preflight, npm, tag, portable, winget, homebrew) — this plan follows it | file at skill root |
| CLI automation | `.github/workflows/tests.yml` matrix (ubuntu/windows/macos go build+test on push) — **no CLI release-asset automation, no goreleaser**: the 6 binaries are hand-built and hand-uploaded | `ls .github/workflows/` |
| Local installs today | `/opt/homebrew/bin/parley` → 1.48.0; `/opt/homebrew/bin/parley-deck-skill` → 2.12.1 (both brew-managed) | `which`, `--version` |
| Skill install targets | codex, claude, agy, gemini, hermes, qwen (`doctor --target all --json`) | run today |

## 1. Prerequisites (before any step below)

1. Phase 8 complete; claude-1 and kimi-1 review verdicts adjudicated; any fix-up commits are already
   on `lean-organizer` (both repos) and re-verified. This plan then applies on top of the final HEAD.
2. Owner-only parallel action, NOT a channel gate: the attended `parley protocol publish` of the
   combined staged core (`~/.parley/staging/COOPERATION-2.13.0.md`, built from the core TEMPLATE =
   2.10.0 + 2.11.0 hunks + 1.48.0 hunks + this idea's hunks) is performed by the owner ALONE (per
   owner inbox `user-to-codex-1_..._core-publish.md`). Channel publication below neither depends on
   nor waits for it; the pending core publish is reported honestly alongside the channel states.
3. Release date D (YYYY-MM-DD) fixed at dispatch; used in CHANGELOG entries and release bodies.

## 2. Version + CHANGELOG prep (zcode-1, in the lean-organizer worktrees — no dev PRs)

### 2a. CLI → 1.49.0 (worktree `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer`)

Edit exactly:
- `VERSION` → `1.49.0`
- `internal/app/version.go` → `const version = "1.49.0"`
- `CHANGELOG.md` → new top section `## 1.49.0 — <D>` summarizing the A–D feature set (wait/digest,
  audience packet + organizer brief, slim SKILL core, handoff/usage ingest, kimi telemetry fix) plus a
  limitations paragraph, mirroring the 1.48.0 section style.

Then:

```bash
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer"
go build ./... && go vet ./... && go test ./... -count=1 -timeout 45m
git add VERSION internal/app/version.go CHANGELOG.md
git commit -m "[zcode-1] meta-protocol-change-lean-organizer: release CLI 1.49.0"
```

### 2b. Skill → 2.13.0 (worktree `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill`)

Edit exactly:
- `package.json` → `"version": "2.13.0"`
- `package-lock.json` → regenerate rather than hand-edit: `npm install --package-lock-only`
- `skills/parley-deck/references/compatibility.json` → `"skillVersion": "2.13.0"`
- `CHANGELOG.md` → new top section `## 2.13.0 — <D>`.

Then (RELEASING.md preflight):

```bash
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill"
npm test                                  # node + python3 (>=3.10) + add-on manifest check
npm pack --dry-run                        # prepack re-verifies add-on payload manifests
npm run build:portable:current            # smoke-build only; full build in step 4
node bin/parley-deck-skill.js install --target all --dry-run
git add package.json package-lock.json skills/parley-deck/references/compatibility.json CHANGELOG.md
git commit -m "[zcode-1] meta-protocol-change-lean-organizer: release skill 2.13.0"
```

If any add-on payload changed during review fix-ups, run `npm run manifest:addons` first and commit
the regenerated manifests with the version bump.

## 3. Integrate to main (direct, history-preserving) and tag

Both branches are strictly ahead of `origin/main`, so a push is a fast-forward — no merge commits,
no development PRs (owner-authorized for this run).

```bash
# CLI
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer"
git push origin lean-organizer:main
git tag v1.49.0
git push origin v1.49.0

# Skill
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill"
git push origin lean-organizer:main
git tag v2.13.0
git push origin v2.13.0
```

Note: the plain `main` worktrees (`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli`,
`.../parley-deck-skill`) stay where they are during the release; fast-forward them afterwards
(`git -C <main-worktree> fetch origin && git -C <main-worktree> merge --ff-only origin/main`) so the
next run baselines correctly.

## 4. Build release artifacts

### 4a. CLI — 6 binaries (hand-built; no automation exists)

From the CLI worktree at `v1.49.0` (file names follow the v1.48.0 convention — `x64`, not `amd64`):

```bash
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer"
git checkout v1.49.0   # worktree HEAD already is this commit after step 3
mkdir -p dist && rm -f dist/*
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o dist/parley-v1.49.0-darwin-arm64    ./cmd/parley
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/parley-v1.49.0-darwin-x64     ./cmd/parley
CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/parley-v1.49.0-linux-x64      ./cmd/parley
CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o dist/parley-v1.49.0-linux-arm64    ./cmd/parley
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/parley-v1.49.0-windows-x64.exe   ./cmd/parley
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o dist/parley-v1.49.0-windows-arm64.exe ./cmd/parley
shasum -a 256 dist/parley-v1.49.0-* | tee /tmp/cli-1.49.0-sha256.txt
```

### 4b. Skill — npm tarball + 5 portable binaries

```bash
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill"
npm ci
npm test
npm pack                                              # → parley-deck-skill-2.13.0.tgz (prepack gates)
shasum -a 256 parley-deck-skill-2.13.0.tgz            # record; this EXACT tarball is published in step 6
npm run build:portable                                # all 5 → dist/parley-deck-skill-v2.13.0-*
shasum -a 256 dist/* | tee /tmp/skill-2.13.0-sha256.txt
```

## 5. GitHub releases

```bash
# CLI (body = the 1.49.0 CHANGELOG section, mirroring v1.48.0 style)
gh release create v1.49.0 --repo feci/parley-deck-cli \
  --title "Parley Deck CLI 1.49.0" --notes-file /tmp/cli-1.49.0-notes.md \
  dist/parley-v1.49.0-*

# Skill — release publication AUTO-TRIGGERS .github/workflows/release-portable.yml:
# macOS npm test → ubuntu npm test + windows rebuild + `gh release upload --clobber`.
gh release create v2.13.0 --repo feci/parley-deck-skill \
  --title "Parley Deck Skill 2.13.0" --notes-file /tmp/skill-2.13.0-notes.md \
  dist/parley-deck-skill-v2.13.0-*
```

Cross-link the two releases in each body (v1.48.0 did: "Companion installer: …"). Wait for the
skill workflow to finish, then confirm the Windows assets were clobbered by CI and re-digest:

```bash
gh run watch --repo feci/parley-deck-skill   # or: gh run list --repo feci/parley-deck-skill --limit 1
gh release view v2.13.0 --repo feci/parley-deck-skill --json assets --jq '.assets[] | [.name,.digest] | @tsv'
gh release view v1.49.0 --repo feci/parley-deck-cli   --json assets --jq '.assets[] | [.name,.digest] | @tsv'
```

## 6. npm — publish the EXACT packed tarball, browser disabled

`npm whoami` → `tomasfecko` verified 2026-09-24 (read-only). If npm prints a web-verification URL,
do NOT open any browser: write it to
`parley-deck/inbox/zcode-1-to-user_meta-protocol-change-lean-organizer_npm-web-verification.md` and
wait for the owner. (ego-browser exists but the owner instruction for this run is escalation, not
agent browsing; Chrome is never driven.)

```bash
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill"
BROWSER=/usr/bin/true npm publish ./parley-deck-skill-2.13.0.tgz --access public
npm view parley-deck-skill@2.13.0 version dist.integrity dist.tarball
```

`BROWSER=/usr/bin/true` makes any npm attempt to open a URL a no-op instead of launching a browser.

## 7. Homebrew — BOTH formulae, one commit

```bash
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley"
git pull --ff-only origin main          # local checkout is 1 commit behind (verified)

curl -fsSL -o /tmp/parley-deck-cli-v1.49.0.tar.gz   https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.49.0.tar.gz
curl -fsSL -o /tmp/parley-deck-skill-v2.13.0.tar.gz https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.13.0.tar.gz
shasum -a 256 /tmp/parley-deck-cli-v1.49.0.tar.gz /tmp/parley-deck-skill-v2.13.0.tar.gz
```

Edit exactly `url` + `sha256` in BOTH files (nothing else — the install blocks and the skill's
post-install shebang restore stay as-is):
- `Formula/parley-deck-cli.rb` → `.../archive/refs/tags/v1.49.0.tar.gz` + new sha256
- `Formula/parley-deck-skill.rb` → `.../archive/refs/tags/v2.13.0.tar.gz` + new sha256

```bash
git add Formula/parley-deck-cli.rb Formula/parley-deck-skill.rb
git commit -m "parley-deck-cli 1.49.0, parley-deck-skill 2.13.0"   # tap commit convention (verified)
git push origin main
```

## 8. Winget — two PRs against microsoft/winget-pkgs (one application per PR)

Source of truth for hashes: the FINAL GitHub release digests from step 5 (skill Windows exes were
rebuilt by CI). Winget manifests use UPPERCASE hex sha256 (verified in 1.48.0 manifest).

Per application, starting from the current catalog dir (which is also in the synced fork):
`manifests/f/Feci/ParleyDeckCli/1.48.0/` and `manifests/f/Feci/ParleyDeckSkill/2.12.1/`, each with
3 YAMLs (installer / locale.en-US / version).

```bash
# fork checkout (feci/winget-pkgs exists; default branch master)
cd <winget-pkgs fork checkout>
git checkout master && git pull
# sync fork with microsoft/winget-pkgs master first (gh repo sync feci/winget-pkgs --source microsoft/winget-pkgs)

# ---- PR 1: CLI ----
git checkout -b release/parley-deck-cli-1.49.0 origin/master
mkdir -p manifests/f/Feci/ParleyDeckCli/1.49.0
# copy the 3 YAMLs from .../1.48.0/, then set in ALL of them:
#   PackageVersion: 1.49.0
#   installer: InstallerUrl x64 →   https://github.com/feci/parley-deck-cli/releases/download/v1.49.0/parley-v1.49.0-windows-x64.exe
#            InstallerSha256 x64 →  <digest from step 5>
#   installer: arm64 → .../parley-v1.49.0-windows-arm64.exe + its digest
#   locale:   ReleaseNotesUrl → https://github.com/feci/parley-deck-cli/releases/tag/v1.49.0
#   (rename files to Feci.ParleyDeckCli.<kind>.yaml in the 1.49.0 dir)
git add manifests/f/Feci/ParleyDeckCli/1.49.0
git commit -m "Update: Feci.ParleyDeckCli to 1.49.0"
git push -u origin release/parley-deck-cli-1.49.0
gh pr create --repo microsoft/winget-pkgs --title "Update: Feci.ParleyDeckCli to 1.49.0" \
  --body "Version bump to 1.49.0. Installer hashes taken from the published GitHub release assets."

# ---- PR 2: Skill (separate branch off master, NOT stacked on PR 1) ----
git checkout master
git checkout -b release/parley-deck-skill-2.13.0 origin/master
mkdir -p manifests/f/Feci/ParleyDeckSkill/2.13.0
# copy 3 YAMLs from .../2.12.1/; PackageVersion: 2.13.0 everywhere;
# InstallerUrl x64   → https://github.com/feci/parley-deck-skill/releases/download/v2.13.0/parley-deck-skill-v2.13.0-windows-x64.exe
#              arm64 → .../parley-deck-skill-v2.13.0-windows-arm64.exe   (CI-clobbered digests!)
# locale: ReleaseNotesUrl → v2.13.0 tag; refresh the Description sentence if it names versions.
git add manifests/f/Feci/ParleyDeckSkill/2.13.0
git commit -m "Update: Feci.ParleyDeckSkill to 2.13.0"
git push -u origin release/parley-deck-skill-2.13.0
gh pr create --repo microsoft/winget-pkgs --title "Update: Feci.ParleyDeckSkill to 2.13.0" \
  --body "Version bump to 2.13.0. Installer hashes taken from the published GitHub release assets."
```

Local `winget validate` needs Windows (unavailable on this Mac, as in previous releases); the
catalog's Azure Pipeline validation runs on the PRs — watch both until merged.

## 9. Local system update + skill install into all sessions

```bash
brew update
brew upgrade feci/parley/parley-deck-cli feci/parley/parley-deck-skill
readlink -f "$(which parley-deck-skill)"   # must resolve into /opt/homebrew/Cellar/parley-deck-skill/2.13.0/...
parley --version                            # parley 1.49.0
parley-deck-skill --version                 # 2.13.0

parley-deck-skill install --target all --force
# verify EVERY runtime SKILL.md by content hash against the released source:
shasum -a 256 "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill/skills/parley-deck/SKILL.md"
parley-deck-skill doctor --target all --json   # per-target installed paths
# for each reported runtime skill dir: shasum -a 256 <path>/parley-deck/SKILL.md  → must equal the source hash
```

## 10. Independent post-deploy audit (participant — claude-1 or kimi-1, NOT zcode-1)

Per the brief, a participant who did not execute the release verifies every channel against the real
artifacts and files findings before the release is reported complete. Exact endpoints:

| Channel | Check | Command / URL |
| --- | --- | --- |
| Homebrew CLI | url + sha256 match the tag archive | `curl -fsSL https://raw.githubusercontent.com/feci/homebrew-parley/main/Formula/parley-deck-cli.rb`; re-hash `https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.49.0.tar.gz` |
| Homebrew skill | url + sha256 match the tag archive | same for `Formula/parley-deck-skill.rb` vs `.../parley-deck-skill/archive/refs/tags/v2.13.0.tar.gz` |
| Homebrew install | Cellar resolution + versions | `readlink -f "$(which parley-deck-skill)"`, `parley --version`, `parley-deck-skill --version` |
| npm | version, integrity, exact tarball | `npm view parley-deck-skill@2.13.0 dist.integrity dist.tarball`; `curl -fsSL <tarball> | shasum -a 512`-equivalent integrity compare vs published pack |
| GitHub CLI | 6 assets + digests | `gh release view v1.49.0 --repo feci/parley-deck-cli --json assets --jq '.assets[] | [.name,.digest]'`; spot-download 2 assets and re-hash |
| GitHub skill | 5 assets + digests (Windows = CI-built) | `gh release view v2.13.0 --repo feci/parley-deck-skill ...` |
| winget | both PRs state + manifest hashes | `gh pr view <n> --repo microsoft/winget-pkgs` (branches `release/parley-deck-cli-1.49.0`, `release/parley-deck-skill-2.13.0`); compare `InstallerSha256` vs release digests |
| Installed skill | content hash per runtime | source SKILL.md sha256 == every installed runtime SKILL.md sha256 (step 9) |
| Smoke | installed binaries work | isolated `PARLEY_HOME` `parley run` probe with `--no-tui --no-auto --yes` (pattern from `codex-1-to-all_release-1.47.0_channel-verification.md`) |

## 11. Existing release automation found (prefer over hand-editing)

- Skill `release-portable.yml`: builds and uploads Windows assets on release publish (auto) — do NOT
  hand-upload skill Windows exes as final; let CI clobber, then take digests.
- Skill `prepack`/`npm test` manifest gates — regenerate with `npm run manifest:addons`, never edit.
- CLI `tests.yml` runs the 3-OS go matrix on every push including the release push — treat its green
  run on `main` as part of the release evidence.
- Everything else (CLI binaries, GitHub releases, npm publish, both formulae, winget manifests) is
  manual today; the commands above are the exact manual sequence.

## 12. Unknowns / assumptions (explicit)

1. Final CHANGELOG/release-notes wording for 1.49.0/2.13.0 is drafted at dispatch, after review
   verdicts — reviews may still add fix-up commits that belong in the notes.
2. npm may still demand web verification despite the working token (`whoami` OK today); the
   escalation path in §6 is mandatory then. No browser is opened by the agent either way.
3. `git push origin lean-organizer:main` assumes no branch protection blocking direct pushes — the
   previous two releases pushed main directly, so this is established practice, not a guarantee.
4. Local Windows `winget validate` is impossible here; PR pipeline validation substitutes (as in the
   1.46.0–2.12.1 releases, all merged).
5. The attended core publish (`parley protocol publish` of combined 2.13.0 core) is the owner's
   action; this plan covers only the CLI/skill channels and records the core as a parallel
   prerequisite (§1.2), not a channel step.
6. Skill `package-lock.json` bump is delegated to `npm install --package-lock-only` (avoids
   hand-editing the two version fields).

## 13. Genuine release blockers (none found in the channel machinery)

- No structural blocker: every channel, script, formula path, winget ID, fork, and convention was
  verified to exist today, and both branches fast-forward onto `origin/main`.
- Gates, in order: (1) Phase 8 completion with review findings resolved — DONE 2026-09-24 (closure
  commit on the CLI branch; zero-fix closing consensus signed, LE-7 PASS); (2) no channel gate on the
  core publish — the owner's attended publish of the combined 2.13.0 core is a parallel owner-only
  action, currently pending, and is reported as such; (3) npm web-verification escalation if
  triggered — publication of the npm channel waits for the owner in that case; (4) winget channel
  state is reported honestly: two PRs opened by us is the completable action; their merge by the
  external catalog maintainers is outside our control and is NOT claimed as a completion gate —
  open PRs are reported as open until they merge.
