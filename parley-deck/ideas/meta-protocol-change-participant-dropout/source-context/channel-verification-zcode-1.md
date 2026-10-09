---
agent: zcode-1
idea: meta-protocol-change-participant-dropout
date: 2026-10-09
role: post-close channel verification (AC14), fresh process
reviewed-at: 2026-10-09T00:11Z (UTC), delivery dir release-delivery/2026-10-09-participant-dropout
protocol-context:
  context_mode: full
  source_sha256: 091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf
  packet_sha256: 091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf
  fallback: none
---

# Channel verification — zcode-1 (AC14 delivery)

Post-close channel verification required by the owner brief, run as a fresh zcode-1
process. I was the independent reviewer for AC1–13 (review rounds 01/02 + goal check);
that evidence is closed and was not rerun. Scope here is delivery AC14 only: verify the
delivered channels against live state, treating `channel-evidence.json` and all delivery
JSON as producer claims. All verdicts below are PRIMARY provenance — commands were
executed by me against live GitHub/tap/local state on 2026-10-09 unless explicitly marked
otherwise. Closed artifacts were preserved (verified, see §8).

## 1. CLI release channel — PASS

- PR76 MERGED, merge commit `2f4c9afc0c52f78ac1b23f733a494aa4385b1695`
  (`gh pr view 76 --repo feci/parley-deck-cli --json state,mergeCommit`), merged
  2026-10-08T23:50:01Z; `branches/main` HEAD is the same SHA (merge landed on main).
- Tag `v1.52.0` is an annotated tag (`c3056c67…`) peeling to commit `2f4c9af…`
  (`gh api repos/feci/parley-deck-cli/git/tags/c3056c67…`); tag message
  "[codex-1] meta-protocol-change-participant-dropout: CLI 1.52.0".
- Live release `releases/tags/v1.52.0`: exactly 7 assets (6 binaries + `sha256.json`),
  all `state: uploaded`. Every one of the 7 live API digests equals the recorded sha256
  in `channel-evidence.json` / `github-release-verification.json`, byte counts included:
  darwin-arm64 `e9887a54…`, darwin-x64 `38ef93f1…`, linux-arm64 `5c95175f…`,
  linux-x64 `5bc66ea5…`, windows-arm64 `4acfac54…`, windows-x64 `0067d91e…`,
  sha256.json `7f8081e5…`.
- Representative downloads (curl from the release URLs, `/tmp` scratch):
  `parley-v1.52.0-darwin-arm64` → `e9887a54…` ✓, `parley-v1.52.0-linux-x64` →
  `5bc66ea5…` ✓, `sha256.json` → `7f8081e5…` ✓. `sha256.json` content lists all six
  binaries with the same digests as the live API.
- `go version -m` on both downloaded binaries: `vcs.revision=2f4c9afc0c52f78ac1b23f733a494aa4385b1695`,
  `vcs.modified=false` (clean checkout), `vcs.time=2026-10-08T23:50:01Z`; darwin build
  is `GOOS=darwin GOARCH=arm64`. The downloaded darwin-arm64 binary runs and reports
  `parley 1.52.0`. (2 of 6 binaries exercised locally; all 6 verified via live digests —
  see Limits.)

## 2. Skill release channel — PASS

- PR9 MERGED, merge commit `dbdb91972ab8e3a0883aac3e8898fa8c118de108` (= skill main
  HEAD); tag `v2.16.0` (annotated `03f7e9d1…`) peels to the same commit.
- Live release `releases/tags/v2.16.0`: 4 assets, digests match the records exactly —
  tgz `80f3172a…`, macos-arm64 `76ca7896…`, windows-arm64 `0bc03c66…`,
  windows-x64 `1cb4fe71…`. That is 7 CLI + 4 skill = 11 recorded asset hashes in
  `github-release-verification.json`, all consistent with live state.
- Hosted portable workflow 37861735927 (`release-portable.yml`, event `release`,
  head `v2.16.0`/`dbdb919…`): `conclusion: success`. Jobs: `macos-tests` on
  **macos-latest** (npm ci + npm test, success), then the `windows` job on
  **ubuntu-latest** (npm ci, npm test, `build:portable:windows`, "Upload release
  assets", all success). This matches the claim "Mac full tests then Ubuntu full
  tests/build/upload": the Windows skill assets were built and uploaded by the
  successful hosted workflow, not locally guessed.
- Downloaded the tgz: sha256 `80f3172a…` ✓ matches live digest; `package/package.json`
  reports `parley-deck-skill 2.16.0`; payload contains exactly 6 skills
  (parley-deck core + parley-bidding, parley-design, parley-design-check,
  parley-tracker, parley-worktrees).
- npm remains owner-only: `registry.npmjs.org/parley-deck-skill` dist-tags latest =
  **2.13.0** (no 2.15.0/2.16.0 published). Correct — publication is an owner action.

## 3. Homebrew channel — PASS

- Tap `feci/homebrew-parley` main HEAD = `69b962d9701c633fb700a24fa6051f8ef517161b`;
  the commit touches exactly `Formula/parley-deck-cli.rb` and
  `Formula/parley-deck-skill.rb` ("CLI 1.52.0 and skill 2.16.0").
- Live formulas: urls are the v1.52.0 / v2.16.0 source archives with sha256
  `613a8cde…` / `19dcff39…`. I downloaded both source archives from GitHub and
  recomputed: `613a8cdec1dbe750f8d2befcee0b00d717d2fe93b8390c667dd79ac82bf11e72` and
  `19dcff3921239d82e2a6fff9e06ce2ea9cd1d8057424accc284efac48b4c5f8f` — matching the
  formulas and `archives/source-archives.json`.
- Canonical checkout `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley`:
  HEAD `69b962d`, clean tree, formula file hashes `7af00662…`/`236fbea6…` equal the
  records. Active tap `feci/parley` (`/opt/homebrew/Library/Taps/feci/homebrew-parley`)
  HEAD is also `69b962d`.
- Checks: `logs/homebrew-style.log` — "2 files inspected, no offenses detected";
  `logs/homebrew-upgrade.log` — real upgrades 1.51.0→1.52.0 and 2.15.0→2.16.0;
  `logs/homebrew-test.log` — real `brew test` executions against the 1.52.0/2.16.0
  Cellars; audit/upgrade/test exit 0 recorded. I additionally ran
  **`brew audit --strict --online feci/parley/parley-deck-cli feci/parley/parley-deck-skill`
  myself: exit 0, no findings.**
- Active commands: `/opt/homebrew/bin/parley --version` → `parley 1.52.0`;
  `/opt/homebrew/bin/parley-deck-skill --version` → `2.16.0`.

## 4. WinGet channel (skill-only PR) — PASS as delivered (upstream merge out of scope)

- PR microsoft/winget-pkgs#449194: OPEN, not draft, headRefOid
  `265b8878acb3ab6921dee092ec91cab9943a1c19`, exactly three added files under
  `manifests/f/Feci/ParleyDeckSkill/2.16.0/` (installer / locale.en-US / version).
  Delivery copies of the three manifests are identical to the PR head content modulo
  CRLF line endings (diff --strip-trailing-r clean).
- Installer manifest: `InstallerSha256` values `1CB4FE71…` (x64) and `0BC03C66…`
  (arm64) equal the live GitHub release digests for the two Windows skill assets
  (case-insensitive), and both `InstallerUrl`s point at the live v2.16.0 assets.
  PackageVersion 2.16.0, ManifestVersion 1.12.0.
- `validate-winget.py` inspected: it downloads the real aka.ms 1.12.0 JSON schemas,
  validates with jsonschema, asserts CRLF + PackageIdentifier/Version, and cross-checks
  installer hashes/URLs against the workflow-produced `skill/windows-assets.json`.
  `winget-validation.json` records 3/3 schema pass — consistent with the manifests.
- Honest limits, correctly disclosed by the delivery: native `winget validate`/install
  is unavailable on macOS; no native install pass is claimed; upstream Windows
  validation/merge is pending (2.16.0 is not yet on winget-pkgs master — consistent).
- CLI WinGet held: on winget-pkgs master, `manifests/f/Feci/ParleyDeckCli` tops out at
  1.48.0 — no 1.52.0 submission exists. Matches "CLI WinGet remains held".

## 5. Runtime installs — PASS

Full independent recheck, not trusting recorded hashes:

- `runtime-verification.json` enumerates 114 checks. I recomputed sha256 of **every one
  of the 114 installed files live** and byte-compared each against the corresponding
  file extracted from the exact release tarball (`80f3172a…`): **114/114 exist,
  114/114 byte-identical to the tarball, 114/114 equal to the recorded hashes,
  0 mismatches.** 19 distinct install roots = 15 managed + the 4 required generic
  paths (`~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck`,
  `~/.config/opencode/skills/parley-deck`), 6 skills × 19 targets = 114.
- Doctor: producer logs show `doctor --target all` ok with all 15 managed targets
  "valid" plus 4 generic doctor runs ok (`logs/runtime-all-doctor.json`,
  `logs/hermes-*-doctor.json`, `logs/opencode-config-doctor.json`). I spot-checked
  myself: `parley-deck-skill doctor --target codex` → all 6 valid;
  `doctor --target generic --dest ~/.config/opencode/skills/parley-deck` → all 6 valid.
- Backup `runtime-skills-before-2.16.0.tar.gz` sha256 `f4e53b2b…` matches the recorded
  `backup_sha256`. No roster/model drift observed: `~/.parley/agents.toml` mtime
  2026-09-25 (predates delivery), deck `parley-deck/agents.toml` git-clean.

## 6. Core 2.16.0 staging — PASS (staged-not-published, owner-only)

- `shasum -a 256`: staged `~/.parley/staging/COOPERATION-2.16.0.md` =
  `8c6b95b611901d3d1460ccb8554e26adac67f58d82ea028713791116a3e65a63` ✓; base
  `COOPERATION-2.15.0.md` = `0d81fd807114e4ff67f5fa98bba73622b58c86d66a2bf6b57d1def68cd099c09` ✓;
  `protocol-hunks.json` = `055e7fd6e3e5f04205aebf16d72d04967716de947be697b2ab9bbda8f7b297e9` ✓
  (35 hunks); launch candidate file in `.parley-runtime` hash-equals the staged file.
- **Reproduction (PRIMARY):** applying the 35 sequential replacements from
  `protocol-hunks.json` to staged 2.15.0 yields **byte-identical** output to staged
  2.16.0 (sha256 `8c6b95b6…`, `reproduced == target: True`). No missing hunks.
  (Note: a naive per-hunk presence check falsely flags 10 hunks because later hunks
  rewrite text produced by earlier ones; the authoritative full-sequence application
  is exact.)
- Generic zones preserved: the deck copy of COOPERATION.md at `2f4c9af` vs pre-idea
  base `d16ee9c` — header (lines 1–43, incl. `Transport:` and `Protocol synced:`)
  identical; §2 roster identical; §0 differs only inside the core "Deck bootstrap"
  paragraph by reviewed hunk 1 (the `quota_auto_exclude` description), with the
  deck-specific transport text untouched. The CLI `internal/protocol/defaults` copy
  and the skill `references` copy at the released commits equal the staged core
  modulo their own header templates.
- Registry `~/.parley/protocol/core/` contains only `2.10.0` and `2.13.0` — no global
  publication of 2.15/2.16; latest published is 2.13.0; both remain owner actions.
  `owner-actions.json` lists 2.15 first, then 2.16 (npm 2.15, core 2.15, npm 2.16,
  core 2.16) ✓. I ran no owner commands (`npm publish` / `parley protocol publish`).

## 7. CI caveats — state honestly confirmed

- Candidate `f8f4f1f` check-runs (live): `go build & test (ubuntu-latest)` success ×2,
  `(macos-latest)` success ×2, `(windows-latest)` **failure** ×2. Linux+Mac PASS as
  claimed; the inherited Windows job has now concluded as failure — experimental
  support, and the delivery does not claim it passed. The native full test/vet/build
  pass (recorded in `full-host-verification.json`) plus my own closed AC1–13
  current-tree evidence remain the authoritative gate; I did not rerun the suites
  (per brief).
- Merged-main CI retained: run 37861593668 on `2f4c9af` was in progress at check time
  (not cancelled). Earlier main runs show the same pre-existing pattern of overall
  workflow failure driven by the experimental Windows job.
- Six cancelled runs (`superseded-ci-cancellations.json`): all six verified live as
  `completed/cancelled` (37861543258, 37861538631, 37861447536, 37861443603,
  37859998106, 37859992486). Their head commits (`ef1061c`, `9083c54`, `9b6edba`)
  each touch **zero** files outside `parley-deck/` — deck-only metadata reruns, product
  diff empty, correctly framed as a capacity action, not a test pass.

## 8. Closed artifacts preserved

All 7 entries in `closed-artifacts-integrity.json` (FINAL.md, IMPLEMENTATION.md,
review/consensus.md, review/consensus-cycle-01.md, review/round-01/zcode-1.md,
review/round-02/zcode-1.md, goal-check-zcode-1.md) recomputed and hash-matched — the
closed implementation record is untouched by delivery.

## Issues / limits

1. **Windows CI failure (experimental)** — disclosed by the delivery, verified live.
   Hosted Windows signal is red; Linux/Mac hosted plus the native full pass are green.
   Not a delivery defect; noted because "CI green" must not be claimed for Windows.
2. **No native winget validation on macOS** — schema validation is producer-scripted
   (genuine, live-schema) plus my digest/URL cross-checks; upstream Windows validation
   and merge remain pending. Disclosed by the delivery; excluded from my PASS scope.
3. **Representative-download scope** — I downloaded and exercised 3 of 7 CLI assets
   (darwin-arm64, linux-x64, sha256.json) plus the skill tgz and both Homebrew source
   archives; the remaining 4 CLI binaries are verified against live GitHub API digests
   (which equal the producer's independently-downloaded records) but not re-downloaded
   here, per the brief's "representative downloads" instruction.
4. **Doctor coverage** — producer ran 5 doctor invocations covering 19 targets; I
   re-ran 2 targets myself. The full-strength check here is the 114/114 byte-level
   tarball comparison, which subsumes doctor's file checks.
5. **`gh search` index lag** — a PR search for "Feci.ParleyDeck" returned nothing while
   PR 449194 verifiably exists; all winget conclusions rest on direct API calls, not
   search.
6. Owner-only actions outstanding (excluded from this verdict per the brief): npm
   publish 2.15.0 then 2.16.0; `parley protocol publish` 2.15.0 then 2.16.0; upstream
   winget-pkgs#449194 merge.

## Overall verdict — AC14 delivery: **PASS**

All seven delivered channels (CLI release, skill release + hosted portable workflow,
both Homebrew formulae, skill-only WinGet PR, 19-target runtime installation, core
2.16.0 staging, honest CI caveats) verified against live state with no material
inconsistency and no suppressed finding. Owner-only npm/core publication and the
upstream WinGet merge remain open owner actions and are excluded from this verdict as
the brief directs.
